package storage

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"fh6worker/internal/telemetry"
)

const (
	trackBaselineSaveMatchedExisting = "matched_existing"
	trackBaselineSaveCreatedTrack    = "created_track"
)

const trackBaselineSelectSQL = `SELECT id, track_id, COALESCE(start_ms, 0), COALESCE(end_ms, 0),
	COALESCE(duration_ms, 0), COALESCE(confidence, 0), avg_speed_kmh, max_speed_kmh,
	route_progress_01, geometry_length_meters, track_length_error_pct, distance_traveled_delta_meters,
	current_race_time_delta_seconds, avg_lateral_error_meters, max_lateral_error_meters, COALESCE(warning_flags, ''),
	COALESCE(event_count, 0), COALESCE(driver_mode, 'unknown'), COALESCE(driver_mode_confidence, 0), COALESCE(driver_mode_evidence_json, ''), COALESCE(valid, 0),
	car_ordinal, COALESCE(car_class, ''), car_pi, COALESCE(drivetrain, ''), COALESCE(game_mode, 'unknown'), COALESCE(created_at, '')
	FROM track_baseline_run`

func (s *Store) ListTrackBaselineRuns(trackID int64, limit int) ([]TrackBaselineRun, error) {
	if limit <= 0 || limit > 500 {
		limit = 100
	}
	rows, err := s.db.Query(trackBaselineSelectSQL+` WHERE track_id = ? ORDER BY valid DESC, duration_ms ASC, confidence DESC, created_at DESC, id DESC LIMIT ?`, trackID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanTrackBaselineRuns(rows)
}

func (s *Store) DeleteTrackBaselineRun(id int64) error {
	if id <= 0 {
		return errors.New("track baseline id is required")
	}
	_, err := s.db.Exec(`DELETE FROM track_baseline_run WHERE id = ?`, id)
	return err
}

func (s *Store) SaveTrackBaselineCapture(trackID int64, samples []telemetry.NormalizedTelemetry, events []telemetry.DetectedEvent) (*TrackBaselineRun, error) {
	track, err := s.GetBenchmarkTrack(trackID)
	if err != nil {
		return nil, err
	}
	baseline, err := s.buildTrackBaselineRun(track.ID, *track, samples, events)
	if err != nil {
		return nil, err
	}
	return s.insertTrackBaselineRun(baseline)
}

func (s *Store) SaveTrackBaselineCaptureAuto(preferredTrackID int64, name string, trackType string, samples []telemetry.NormalizedTelemetry, events []telemetry.DetectedEvent) (*TrackBaselineSaveResult, error) {
	if len(samples) < 2 {
		return nil, errors.New("track baseline requires telemetry samples")
	}
	if preferredTrackID > 0 {
		track, err := s.GetBenchmarkTrack(preferredTrackID)
		if err != nil {
			return nil, err
		}
		baseline, err := s.buildTrackBaselineRun(track.ID, *track, samples, events)
		if err != nil {
			return nil, fmt.Errorf("selected track baseline pass was not detected: %w", err)
		}
		saved, err := s.insertTrackBaselineRun(baseline)
		if err != nil {
			return nil, err
		}
		return &TrackBaselineSaveResult{
			Track:    *track,
			Baseline: *saved,
			Action:   trackBaselineSaveMatchedExisting,
		}, nil
	}
	name = strings.TrimSpace(name)
	if name == "" {
		name = "Baseline Track"
	}
	sourceMode := dominantSamplesGameMode(samples)
	input, err := buildBenchmarkTrackInputFromSamples(name, sourceMode, samples, BenchmarkTrackExtractionInput{
		Name:           name,
		TrackType:      trackType,
		ExtractionMode: benchmarkExtractionAutoBestLap,
	})
	if err != nil {
		return nil, err
	}
	candidates, err := s.FindSimilarBenchmarkTracks(input)
	if err != nil {
		return nil, err
	}
	for _, candidate := range candidates {
		if candidate.MatchLevel != "strong" {
			continue
		}
		baseline, err := s.buildTrackBaselineRun(candidate.Track.ID, candidate.Track, samples, events)
		if err != nil {
			return nil, err
		}
		saved, err := s.insertTrackBaselineRun(baseline)
		if err != nil {
			return nil, err
		}
		candidateCopy := candidate
		return &TrackBaselineSaveResult{
			Track:          candidate.Track,
			Baseline:       *saved,
			Action:         trackBaselineSaveMatchedExisting,
			MatchCandidate: &candidateCopy,
		}, nil
	}

	normalized, err := normalizeBenchmarkTrackInput(input)
	if err != nil {
		return nil, err
	}
	if _, err := s.buildTrackBaselineRun(0, BenchmarkTrack{BenchmarkTrackInput: normalized}, samples, events); err != nil {
		return nil, err
	}
	created, err := s.CreateBenchmarkTrack(input)
	if err != nil {
		return nil, err
	}
	saved, err := s.SaveTrackBaselineCapture(created.ID, samples, events)
	if err != nil {
		return nil, err
	}
	return &TrackBaselineSaveResult{
		Track:    *created,
		Baseline: *saved,
		Action:   trackBaselineSaveCreatedTrack,
	}, nil
}

func (s *Store) buildTrackBaselineRun(trackID int64, track BenchmarkTrack, samples []telemetry.NormalizedTelemetry, events []telemetry.DetectedEvent) (TrackBaselineRun, error) {
	if len(samples) < 2 {
		return TrackBaselineRun{}, errors.New("track baseline requires telemetry samples")
	}
	runs := analyzeBenchmarkRuns(0, track, samples)
	if len(runs) == 0 {
		return TrackBaselineRun{}, errors.New("no valid track baseline pass detected")
	}
	best := runs[0]
	for _, run := range runs[1:] {
		if betterBaselineRun(run, best) {
			best = run
		}
	}
	window := samplesBetween(samples, best.StartMS, best.EndMS)
	if len(window) == 0 {
		window = samples
	}
	detection := DetectDriverMode(window, eventsBetween(events, best.StartMS, best.EndMS), dominantSamplesGameMode(window))
	best.DriverMode = NormalizeDriverMode(detection.Mode)
	best.DriverModeConfidence = detection.Confidence
	if raw, err := json.Marshal(detection); err == nil {
		best.DriverModeEvidenceJSON = string(raw)
	}
	vehicle := trackVehicleKeyFromSamples(window)
	if trackVehicleMapKey(vehicle) == "" {
		vehicle = trackVehicleKeyFromSamples(samples)
	}
	gameMode := dominantSamplesGameMode(window)
	if telemetry.NormalizeGameMode(gameMode) == telemetry.GameModeUnknown {
		gameMode = dominantSamplesGameMode(samples)
	}
	baseline := TrackBaselineRun{
		TrackID:                     trackID,
		Vehicle:                     vehicle,
		StartMS:                     best.StartMS,
		EndMS:                       best.EndMS,
		DurationMS:                  best.DurationMS,
		Confidence:                  best.Confidence,
		AvgSpeedKmh:                 best.AvgSpeedKmh,
		MaxSpeedKmh:                 best.MaxSpeedKmh,
		RouteProgress01:             best.RouteProgress01,
		GeometryLengthMeters:        best.GeometryLengthMeters,
		TrackLengthErrorPct:         best.TrackLengthErrorPct,
		DistanceTraveledDeltaMeters: best.DistanceTraveledDeltaMeters,
		CurrentRaceTimeDeltaSeconds: best.CurrentRaceTimeDeltaSeconds,
		AvgLateralErrorMeters:       best.AvgLateralErrorMeters,
		MaxLateralErrorMeters:       best.MaxLateralErrorMeters,
		WarningFlags:                best.WarningFlags,
		EventCount:                  int64(len(eventsBetween(events, best.StartMS, best.EndMS))),
		DriverMode:                  best.DriverMode,
		DriverModeConfidence:        best.DriverModeConfidence,
		DriverModeEvidenceJSON:      best.DriverModeEvidenceJSON,
		Valid:                       best.Valid,
		GameMode:                    telemetry.NormalizeGameMode(gameMode),
	}
	return baseline, nil
}

func (s *Store) insertTrackBaselineRun(run TrackBaselineRun) (*TrackBaselineRun, error) {
	now := nowText()
	result, err := s.db.Exec(`INSERT INTO track_baseline_run (
		track_id, start_ms, end_ms, duration_ms, confidence, avg_speed_kmh, max_speed_kmh,
		route_progress_01, geometry_length_meters, track_length_error_pct, distance_traveled_delta_meters,
		current_race_time_delta_seconds, avg_lateral_error_meters, max_lateral_error_meters, warning_flags,
		event_count, driver_mode, driver_mode_confidence, driver_mode_evidence_json, valid,
		car_ordinal, car_class, car_pi, drivetrain, game_mode, created_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		run.TrackID, run.StartMS, run.EndMS, run.DurationMS, run.Confidence,
		nullableFloat(run.AvgSpeedKmh), nullableFloat(run.MaxSpeedKmh),
		nullableFloat(run.RouteProgress01), nullableFloat(run.GeometryLengthMeters), nullableFloat(run.TrackLengthErrorPct), nullableFloat(run.DistanceTraveledDeltaMeters),
		nullableFloat(run.CurrentRaceTimeDeltaSeconds), nullableFloat(run.AvgLateralErrorMeters), nullableFloat(run.MaxLateralErrorMeters), run.WarningFlags,
		run.EventCount, run.DriverMode, run.DriverModeConfidence, run.DriverModeEvidenceJSON, boolInt(run.Valid),
		nullableInt(run.Vehicle.CarOrdinal), strings.TrimSpace(run.Vehicle.CarClass), nullableInt(run.Vehicle.CarPI), strings.TrimSpace(run.Vehicle.Drivetrain), telemetry.NormalizeGameMode(run.GameMode), now,
	)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	row := s.db.QueryRow(trackBaselineSelectSQL+` WHERE id = ?`, id)
	saved, err := scanTrackBaselineRun(row)
	if err != nil {
		return nil, err
	}
	return &saved, nil
}

func scanTrackBaselineRuns(rows *sql.Rows) ([]TrackBaselineRun, error) {
	runs := []TrackBaselineRun{}
	for rows.Next() {
		run, err := scanTrackBaselineRun(rows)
		if err != nil {
			return nil, err
		}
		runs = append(runs, run)
	}
	return runs, rows.Err()
}

func scanTrackBaselineRun(row scanner) (TrackBaselineRun, error) {
	var run TrackBaselineRun
	var avgSpeed, maxSpeed, routeProgress, geometryLength, trackLengthError, distanceDelta, raceTimeDelta, avgLateral, maxLateral, driverModeConfidence sql.NullFloat64
	var carOrdinal, carPI sql.NullInt64
	var carClass, drivetrain sql.NullString
	var valid int
	if err := row.Scan(
		&run.ID, &run.TrackID, &run.StartMS, &run.EndMS, &run.DurationMS, &run.Confidence,
		&avgSpeed, &maxSpeed, &routeProgress, &geometryLength, &trackLengthError, &distanceDelta,
		&raceTimeDelta, &avgLateral, &maxLateral, &run.WarningFlags,
		&run.EventCount, &run.DriverMode, &driverModeConfidence, &run.DriverModeEvidenceJSON, &valid,
		&carOrdinal, &carClass, &carPI, &drivetrain, &run.GameMode, &run.CreatedAt,
	); err != nil {
		return TrackBaselineRun{}, err
	}
	run.AvgSpeedKmh = floatPtr(avgSpeed)
	run.MaxSpeedKmh = floatPtr(maxSpeed)
	run.RouteProgress01 = floatPtr(routeProgress)
	run.GeometryLengthMeters = floatPtr(geometryLength)
	run.TrackLengthErrorPct = floatPtr(trackLengthError)
	run.DistanceTraveledDeltaMeters = floatPtr(distanceDelta)
	run.CurrentRaceTimeDeltaSeconds = floatPtr(raceTimeDelta)
	run.AvgLateralErrorMeters = floatPtr(avgLateral)
	run.MaxLateralErrorMeters = floatPtr(maxLateral)
	run.DriverMode = NormalizeDriverMode(run.DriverMode)
	run.DriverModeConfidence = floatFromNull(driverModeConfidence)
	run.Valid = valid != 0
	run.GameMode = telemetry.NormalizeGameMode(run.GameMode)
	run.Vehicle = TrackVehicleKey{
		CarOrdinal: intPtr(carOrdinal),
		CarClass:   strings.TrimSpace(carClass.String),
		CarPI:      intPtr(carPI),
		Drivetrain: strings.TrimSpace(drivetrain.String),
	}
	run.Vehicle.Label = formatTrackVehicleLabel(run.Vehicle)
	return run, nil
}

func trackVehicleKeyFromSamples(samples []telemetry.NormalizedTelemetry) TrackVehicleKey {
	for _, sample := range samples {
		if sample.CarOrdinal <= 0 {
			continue
		}
		ordinal := int64(sample.CarOrdinal)
		var pi *int64
		if sample.CarPI > 0 {
			value := int64(sample.CarPI)
			pi = &value
		}
		key := TrackVehicleKey{
			CarOrdinal: &ordinal,
			CarClass:   strings.TrimSpace(sample.CarClass),
			CarPI:      pi,
			Drivetrain: strings.TrimSpace(sample.Drivetrain),
		}
		key.Label = formatTrackVehicleLabel(key)
		return key
	}
	return TrackVehicleKey{Label: "--"}
}

func dominantSamplesGameMode(samples []telemetry.NormalizedTelemetry) string {
	counts := map[string]int{}
	for _, sample := range samples {
		mode := telemetry.NormalizeGameMode(sample.GameMode)
		if mode != telemetry.GameModeUnknown {
			counts[mode]++
		}
	}
	return dominantGameMode(counts)
}
