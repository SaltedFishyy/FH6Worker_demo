package storage

import (
	"fmt"
	"math"
	"sort"
	"strings"
)

const autoBaselineMinConfidence = 0.65

func (s *Store) ListTrackProfiles() ([]TrackProfile, error) {
	tracks, err := s.ListBenchmarkTracks()
	if err != nil {
		return nil, err
	}
	profiles := make([]TrackProfile, 0, len(tracks))
	for _, track := range tracks {
		profile, err := s.buildTrackProfile(track)
		if err != nil {
			return nil, err
		}
		profiles = append(profiles, *profile)
	}
	return profiles, nil
}

func (s *Store) GetTrackProfile(trackID int64) (*TrackProfile, error) {
	track, err := s.GetBenchmarkTrack(trackID)
	if err != nil {
		return nil, err
	}
	return s.buildTrackProfile(*track)
}

func (s *Store) buildTrackProfile(track BenchmarkTrack) (*TrackProfile, error) {
	bestContexts, err := s.listBenchmarkRunContexts(track.ID, 500, false)
	if err != nil {
		return nil, err
	}
	recentContexts, err := s.listBenchmarkRunContexts(track.ID, 50, true)
	if err != nil {
		return nil, err
	}
	baselineRuns, err := s.ListTrackBaselineRuns(track.ID, 500)
	if err != nil {
		return nil, err
	}

	profile := &TrackProfile{
		Track:             track,
		VehicleReferences: buildTrackVehicleReferences(bestContexts, baselineRuns),
		RecentRuns:        recentContexts,
		Warnings:          []string{},
	}

	groups := map[string]*TrackAutoBaseline{}
	for _, context := range bestContexts {
		if !isAutoBaselineRun(context) {
			continue
		}
		key := trackVehicleMapKey(context.Vehicle)
		if key == "" {
			profile.Warnings = appendUniqueString(profile.Warnings, "baseline_vehicle_identity_missing")
			continue
		}
		group := groups[key]
		if group == nil {
			copyContext := context
			group = &TrackAutoBaseline{
				Vehicle:    context.Vehicle,
				BestRun:    copyContext,
				RecentRuns: []TrackRunContext{},
			}
			groups[key] = group
		}
		group.RunCount++
		if len(group.RecentRuns) < 5 {
			group.RecentRuns = append(group.RecentRuns, context)
		}
		if betterBaselineRun(context.Run, group.BestRun.Run) {
			group.BestRun = context
		}
	}

	profile.AutoBaselines = make([]TrackAutoBaseline, 0, len(groups))
	for _, group := range groups {
		profile.AutoBaselines = append(profile.AutoBaselines, *group)
	}
	sort.Slice(profile.AutoBaselines, func(i, j int) bool {
		left := profile.AutoBaselines[i].BestRun.Run.DurationMS
		right := profile.AutoBaselines[j].BestRun.Run.DurationMS
		if left == right {
			return trackVehicleMapKey(profile.AutoBaselines[i].Vehicle) < trackVehicleMapKey(profile.AutoBaselines[j].Vehicle)
		}
		return left < right
	})

	if len(profile.AutoBaselines) == 0 && !hasAutoTrackBaseline(profile.VehicleReferences) {
		profile.Warnings = appendUniqueString(profile.Warnings, "no_auto_baseline")
	}

	return profile, nil
}

func (s *Store) RenameBenchmarkTrack(trackID int64, name string) (*BenchmarkTrack, error) {
	name = strings.TrimSpace(name)
	if trackID <= 0 {
		return nil, fmt.Errorf("benchmark track id is required")
	}
	if name == "" {
		return nil, fmt.Errorf("benchmark track name is required")
	}
	if _, err := s.db.Exec(`UPDATE benchmark_track SET name = ?, updated_at = ? WHERE id = ?`, name, nowText(), trackID); err != nil {
		return nil, err
	}
	return s.GetBenchmarkTrack(trackID)
}

func (s *Store) FindSimilarBenchmarkTracks(input BenchmarkTrackInput) ([]TrackMergeCandidate, error) {
	normalized, err := normalizeBenchmarkTrackInput(input)
	if err != nil {
		return nil, err
	}
	tracks, err := s.ListBenchmarkTracks()
	if err != nil {
		return nil, err
	}
	candidates := make([]TrackMergeCandidate, 0)
	for _, track := range tracks {
		candidate, ok := similarBenchmarkTrackCandidate(normalized, track)
		if ok {
			candidates = append(candidates, candidate)
		}
	}
	sort.Slice(candidates, func(i, j int) bool {
		leftRank := trackMatchLevelRank(candidates[i].MatchLevel)
		rightRank := trackMatchLevelRank(candidates[j].MatchLevel)
		if leftRank != rightRank {
			return leftRank > rightRank
		}
		if candidates[i].RouteFitScore == candidates[j].RouteFitScore {
			return candidates[i].LengthErrorPct < candidates[j].LengthErrorPct
		}
		return candidates[i].RouteFitScore > candidates[j].RouteFitScore
	})
	return candidates, nil
}

func trackMatchLevelRank(value string) int {
	switch value {
	case "strong":
		return 2
	case "medium":
		return 1
	default:
		return 0
	}
}

func (s *Store) MergeBenchmarkTrackInput(trackID int64, input BenchmarkTrackInput) (*BenchmarkTrack, error) {
	return s.UpdateBenchmarkTrack(trackID, input)
}

func (s *Store) listBenchmarkRunContexts(trackID int64, limit int, recent bool) ([]TrackRunContext, error) {
	if limit <= 0 || limit > 500 {
		limit = 100
	}
	orderBy := ` ORDER BY br.duration_ms ASC, br.confidence DESC, br.created_at DESC, br.id DESC`
	if recent {
		orderBy = ` ORDER BY br.created_at DESC, br.id DESC`
	}
	rows, err := s.db.Query(benchmarkRunSelectSQL+` WHERE br.track_id = ?`+orderBy+` LIMIT ?`, trackID, limit)
	if err != nil {
		return nil, err
	}

	runs := []BenchmarkRun{}
	for rows.Next() {
		run, err := scanBenchmarkRun(rows)
		if err != nil {
			return nil, err
		}
		runs = append(runs, run)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}

	contexts := make([]TrackRunContext, 0, len(runs))
	for _, run := range runs {
		session, err := s.GetTelemetrySession(run.SessionID)
		if err != nil {
			return nil, err
		}
		contexts = append(contexts, TrackRunContext{
			Run:     run,
			Session: *session,
			Vehicle: trackVehicleKeyFromSession(*session),
		})
	}
	return contexts, nil
}

func isAutoBaselineRun(context TrackRunContext) bool {
	return context.Run.Valid &&
		NormalizeDriverMode(context.Run.DriverMode) == driverModeAuto &&
		context.Run.DriverModeConfidence >= autoBaselineMinConfidence
}

func buildTrackVehicleReferences(contexts []TrackRunContext, baselines []TrackBaselineRun) []TrackVehicleReference {
	groups := map[string]*TrackVehicleReference{}
	speedSums := map[string]float64{}
	speedCounts := map[string]int{}
	for _, context := range contexts {
		if !context.Run.Valid {
			continue
		}
		key := trackVehicleMapKey(context.Vehicle)
		if key == "" {
			continue
		}
		group := groups[key]
		if group == nil {
			group = &TrackVehicleReference{
				Vehicle:            context.Vehicle,
				RecentRuns:         []TrackRunContext{},
				RecentBaselineRuns: []TrackBaselineRun{},
			}
			groups[key] = group
		}
		group.ValidRunCount++
		group.EventCount += context.Run.EventCount
		if len(group.RecentRuns) < 5 {
			group.RecentRuns = append(group.RecentRuns, context)
		}
		if context.Run.AvgSpeedKmh != nil {
			speedSums[key] += *context.Run.AvgSpeedKmh
			speedCounts[key]++
		}
		if context.Run.MaxSpeedKmh != nil && (group.MaxSpeedKmh == nil || *context.Run.MaxSpeedKmh > *group.MaxSpeedKmh) {
			value := *context.Run.MaxSpeedKmh
			group.MaxSpeedKmh = &value
		}
		if isAutoBaselineRun(context) {
			group.AutoRunCount++
			if group.BestAutoBaseline == nil || betterBaselineRun(context.Run, group.BestAutoBaseline.Run) {
				copyContext := context
				group.BestAutoBaseline = &copyContext
			}
		}
	}
	for _, baseline := range baselines {
		if !baseline.Valid {
			continue
		}
		key := trackVehicleMapKey(baseline.Vehicle)
		if key == "" {
			continue
		}
		group := groups[key]
		if group == nil {
			group = &TrackVehicleReference{
				Vehicle:            baseline.Vehicle,
				RecentRuns:         []TrackRunContext{},
				RecentBaselineRuns: []TrackBaselineRun{},
			}
			groups[key] = group
		}
		group.ValidRunCount++
		group.BaselineRunCount++
		group.EventCount += baseline.EventCount
		if len(group.RecentBaselineRuns) < 5 {
			group.RecentBaselineRuns = append(group.RecentBaselineRuns, baseline)
		}
		if baseline.AvgSpeedKmh != nil {
			speedSums[key] += *baseline.AvgSpeedKmh
			speedCounts[key]++
		}
		if baseline.MaxSpeedKmh != nil && (group.MaxSpeedKmh == nil || *baseline.MaxSpeedKmh > *group.MaxSpeedKmh) {
			value := *baseline.MaxSpeedKmh
			group.MaxSpeedKmh = &value
		}
		if isAutoTrackBaselineRun(baseline) {
			group.AutoRunCount++
			if group.BestTrackBaseline == nil || betterBaselineRun(benchmarkRunFromTrackBaseline(baseline), benchmarkRunFromTrackBaseline(*group.BestTrackBaseline)) {
				copyRun := baseline
				group.BestTrackBaseline = &copyRun
			}
		}
	}
	references := make([]TrackVehicleReference, 0, len(groups))
	for key, group := range groups {
		if speedCounts[key] > 0 {
			avg := speedSums[key] / float64(speedCounts[key])
			group.AvgSpeedKmh = &avg
		}
		references = append(references, *group)
	}
	sort.Slice(references, func(i, j int) bool {
		leftDuration, leftOK := trackReferenceBestDuration(references[i])
		rightDuration, rightOK := trackReferenceBestDuration(references[j])
		if leftOK && rightOK && leftDuration != rightDuration {
			return leftDuration < rightDuration
		}
		if leftOK && !rightOK {
			return true
		}
		if !leftOK && rightOK {
			return false
		}
		return trackVehicleMapKey(references[i].Vehicle) < trackVehicleMapKey(references[j].Vehicle)
	})
	return references
}

func hasAutoTrackBaseline(references []TrackVehicleReference) bool {
	for _, reference := range references {
		if reference.BestTrackBaseline != nil {
			return true
		}
	}
	return false
}

func isAutoTrackBaselineRun(run TrackBaselineRun) bool {
	return run.Valid &&
		NormalizeDriverMode(run.DriverMode) == driverModeAuto &&
		run.DriverModeConfidence >= autoBaselineMinConfidence
}

func benchmarkRunFromTrackBaseline(run TrackBaselineRun) BenchmarkRun {
	return BenchmarkRun{
		ID:                          run.ID,
		TrackID:                     run.TrackID,
		StartMS:                     run.StartMS,
		EndMS:                       run.EndMS,
		DurationMS:                  run.DurationMS,
		Confidence:                  run.Confidence,
		AvgSpeedKmh:                 run.AvgSpeedKmh,
		MaxSpeedKmh:                 run.MaxSpeedKmh,
		RouteProgress01:             run.RouteProgress01,
		GeometryLengthMeters:        run.GeometryLengthMeters,
		TrackLengthErrorPct:         run.TrackLengthErrorPct,
		DistanceTraveledDeltaMeters: run.DistanceTraveledDeltaMeters,
		CurrentRaceTimeDeltaSeconds: run.CurrentRaceTimeDeltaSeconds,
		AvgLateralErrorMeters:       run.AvgLateralErrorMeters,
		MaxLateralErrorMeters:       run.MaxLateralErrorMeters,
		WarningFlags:                run.WarningFlags,
		EventCount:                  run.EventCount,
		DriverMode:                  run.DriverMode,
		DriverModeConfidence:        run.DriverModeConfidence,
		DriverModeEvidenceJSON:      run.DriverModeEvidenceJSON,
		Valid:                       run.Valid,
		CreatedAt:                   run.CreatedAt,
	}
}

func trackReferenceBestDuration(reference TrackVehicleReference) (int64, bool) {
	if reference.BestTrackBaseline != nil {
		return reference.BestTrackBaseline.DurationMS, true
	}
	if reference.BestAutoBaseline != nil {
		return reference.BestAutoBaseline.Run.DurationMS, true
	}
	return 0, false
}

func similarBenchmarkTrackCandidate(input BenchmarkTrackInput, existing BenchmarkTrack) (TrackMergeCandidate, bool) {
	if len(input.Polyline) < 2 || len(existing.Polyline) < 2 {
		return TrackMergeCandidate{}, false
	}
	if input.TrackType == "" || existing.TrackType == "" || input.TrackType != existing.TrackType {
		return TrackMergeCandidate{}, false
	}
	inputLength := input.RouteLengthMeters
	if inputLength <= 0 {
		inputLength = routeLength(input.Polyline)
	}
	existingLength := existing.RouteLengthMeters
	if existingLength <= 0 {
		existingLength = routeLength(existing.Polyline)
	}
	if inputLength <= 0 || existingLength <= 0 {
		return TrackMergeCandidate{}, false
	}
	lengthErrorPct := math.Abs(inputLength-existingLength) / math.Max(inputLength, existingLength) * 100
	startDistance := distanceXZ(input.Start, existing.Start)
	endDistance := distanceXZ(input.End, existing.End)
	routeFit := benchmarkRouteFit(input.TrackType, input.Polyline, existing.Polyline)
	if routeFit.ReverseMatched || !routeFit.DirectionMatched {
		return TrackMergeCandidate{}, false
	}
	matchLevel := ""
	switch input.TrackType {
	case benchmarkTrackTypeSprint:
		if lengthErrorPct <= 5 && routeFit.AvgErrorMeters <= 30 && routeFit.P90ErrorMeters <= 60 && startDistance <= 60 && endDistance <= 60 {
			matchLevel = "strong"
		} else if lengthErrorPct <= 8 && routeFit.AvgErrorMeters <= 45 && routeFit.P90ErrorMeters <= 90 && startDistance <= 100 && endDistance <= 100 {
			matchLevel = "medium"
		}
	case benchmarkTrackTypeCircuit:
		if lengthErrorPct <= 5 && routeFit.AvgErrorMeters <= 30 && routeFit.P90ErrorMeters <= 60 {
			matchLevel = "strong"
		} else if lengthErrorPct <= 8 && routeFit.AvgErrorMeters <= 45 && routeFit.P90ErrorMeters <= 90 {
			matchLevel = "medium"
		}
	}
	if matchLevel == "" {
		return TrackMergeCandidate{}, false
	}
	reason := "route_fit_medium"
	if matchLevel == "strong" {
		reason = "route_fit_strong"
	}
	return TrackMergeCandidate{
		Track:                  existing,
		MatchLevel:             matchLevel,
		LengthErrorPct:         lengthErrorPct,
		StartDistanceMeters:    startDistance,
		EndDistanceMeters:      endDistance,
		ShapeSimilarity:        routeFit.Score,
		RouteFitAvgErrorMeters: routeFit.AvgErrorMeters,
		RouteFitP90ErrorMeters: routeFit.P90ErrorMeters,
		RouteFitScore:          routeFit.Score,
		DirectionMatched:       routeFit.DirectionMatched,
		ReverseMatched:         routeFit.ReverseMatched,
		Reason:                 reason,
	}, true
}

func benchmarkShapeSimilarity(left []BenchmarkPoint, right []BenchmarkPoint) float64 {
	return benchmarkRouteFit(benchmarkTrackTypeSprint, left, right).Score
}

type benchmarkRouteFitResult struct {
	AvgErrorMeters   float64
	P90ErrorMeters   float64
	Score            float64
	DirectionMatched bool
	ReverseMatched   bool
}

func benchmarkRouteFit(trackType string, left []BenchmarkPoint, right []BenchmarkPoint) benchmarkRouteFitResult {
	sampleCount := 64
	if trackType == benchmarkTrackTypeCircuit {
		sampleCount = 96
	}
	leftSamples := resampleBenchmarkPolyline(left, sampleCount)
	rightSamples := resampleBenchmarkPolyline(right, sampleCount)
	if len(leftSamples) == 0 || len(rightSamples) == 0 || len(leftSamples) != len(rightSamples) {
		return benchmarkRouteFitResult{}
	}
	reversed := reverseBenchmarkPoints(rightSamples)
	var forwardAvg, forwardP90 float64
	var reverseAvg, reverseP90 float64
	if trackType == benchmarkTrackTypeCircuit {
		forwardAvg, forwardP90 = bestCyclicRouteFit(leftSamples, rightSamples)
		reverseAvg, reverseP90 = bestCyclicRouteFit(leftSamples, reversed)
	} else {
		forwardAvg, forwardP90 = routeFitErrors(leftSamples, rightSamples)
		reverseAvg, reverseP90 = routeFitErrors(leftSamples, reversed)
	}
	scale := math.Max(80, math.Min(routeLength(left), routeLength(right))*0.08)
	if scale <= 0 {
		return benchmarkRouteFitResult{}
	}
	score := clamp01(1 - forwardAvg/scale)
	return benchmarkRouteFitResult{
		AvgErrorMeters:   forwardAvg,
		P90ErrorMeters:   forwardP90,
		Score:            score,
		DirectionMatched: forwardAvg <= reverseAvg,
		ReverseMatched:   reverseAvg < forwardAvg && reverseP90 <= forwardP90,
	}
}

func routeFitErrors(left []BenchmarkPoint, right []BenchmarkPoint) (float64, float64) {
	if len(left) == 0 || len(left) != len(right) {
		return math.Inf(1), math.Inf(1)
	}
	values := make([]float64, 0, len(left))
	total := 0.0
	for i := range left {
		value := distanceXZ(left[i], right[i])
		total += value
		values = append(values, value)
	}
	sort.Float64s(values)
	p90Index := int(math.Ceil(float64(len(values))*0.9)) - 1
	if p90Index < 0 {
		p90Index = 0
	}
	if p90Index >= len(values) {
		p90Index = len(values) - 1
	}
	return total / float64(len(left)), values[p90Index]
}

func bestCyclicRouteFit(left []BenchmarkPoint, right []BenchmarkPoint) (float64, float64) {
	if len(left) == 0 || len(left) != len(right) {
		return math.Inf(1), math.Inf(1)
	}
	bestAvg := math.Inf(1)
	bestP90 := math.Inf(1)
	for offset := 0; offset < len(right); offset++ {
		shifted := make([]BenchmarkPoint, len(right))
		for i := range right {
			shifted[i] = right[(i+offset)%len(right)]
		}
		avg, p90 := routeFitErrors(left, shifted)
		if avg < bestAvg {
			bestAvg = avg
			bestP90 = p90
		}
	}
	return bestAvg, bestP90
}

func reverseBenchmarkPoints(points []BenchmarkPoint) []BenchmarkPoint {
	out := make([]BenchmarkPoint, len(points))
	for i := range points {
		out[i] = points[len(points)-1-i]
	}
	return out
}

func resampleBenchmarkPolyline(points []BenchmarkPoint, count int) []BenchmarkPoint {
	if len(points) < 2 || count <= 0 {
		return nil
	}
	totalLength := routeLength(points)
	if totalLength <= 0 {
		return nil
	}
	out := make([]BenchmarkPoint, 0, count)
	segmentIndex := 1
	traveled := 0.0
	for i := 0; i < count; i++ {
		target := totalLength * float64(i) / float64(count-1)
		for segmentIndex < len(points)-1 && traveled+distanceXZ(points[segmentIndex-1], points[segmentIndex]) < target {
			traveled += distanceXZ(points[segmentIndex-1], points[segmentIndex])
			segmentIndex++
		}
		a := points[segmentIndex-1]
		b := points[segmentIndex]
		segmentLength := distanceXZ(a, b)
		t := 0.0
		if segmentLength > 0 {
			t = (target - traveled) / segmentLength
		}
		if t < 0 {
			t = 0
		}
		if t > 1 {
			t = 1
		}
		out = append(out, BenchmarkPoint{
			X: a.X + (b.X-a.X)*t,
			Y: a.Y + (b.Y-a.Y)*t,
			Z: a.Z + (b.Z-a.Z)*t,
		})
	}
	return out
}

func betterBaselineRun(candidate BenchmarkRun, current BenchmarkRun) bool {
	if current.ID == 0 {
		return true
	}
	if candidate.DurationMS != current.DurationMS {
		return candidate.DurationMS < current.DurationMS
	}
	if candidate.Confidence != current.Confidence {
		return candidate.Confidence > current.Confidence
	}
	return candidate.CreatedAt > current.CreatedAt
}

func trackVehicleKeyFromSession(session TelemetrySession) TrackVehicleKey {
	key := TrackVehicleKey{
		CarOrdinal: session.CarOrdinal,
		CarClass:   strings.TrimSpace(session.CarClass),
		CarPI:      session.CarPI,
		Drivetrain: strings.TrimSpace(session.Drivetrain),
	}
	key.Label = formatTrackVehicleLabel(key)
	return key
}

func trackVehicleMapKey(vehicle TrackVehicleKey) string {
	if vehicle.CarOrdinal == nil || *vehicle.CarOrdinal <= 0 || strings.TrimSpace(vehicle.CarClass) == "" {
		return ""
	}
	pi := int64(0)
	if vehicle.CarPI != nil {
		pi = *vehicle.CarPI
	}
	return fmt.Sprintf("%d|%s|%d", *vehicle.CarOrdinal, strings.ToUpper(strings.TrimSpace(vehicle.CarClass)), pi)
}

func formatTrackVehicleLabel(vehicle TrackVehicleKey) string {
	parts := []string{}
	if vehicle.CarOrdinal != nil && *vehicle.CarOrdinal > 0 {
		parts = append(parts, fmt.Sprintf("ID %d", *vehicle.CarOrdinal))
	}
	if strings.TrimSpace(vehicle.CarClass) != "" {
		parts = append(parts, strings.TrimSpace(vehicle.CarClass))
	}
	if vehicle.CarPI != nil && *vehicle.CarPI > 0 {
		parts = append(parts, fmt.Sprintf("PI %d", *vehicle.CarPI))
	}
	if strings.TrimSpace(vehicle.Drivetrain) != "" {
		parts = append(parts, strings.TrimSpace(vehicle.Drivetrain))
	}
	if len(parts) == 0 {
		return "--"
	}
	return strings.Join(parts, " / ")
}
