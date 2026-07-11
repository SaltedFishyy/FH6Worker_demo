package storage

import (
	"database/sql"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"fh6worker/internal/telemetry"
)

func TestMigrateIsIdempotent(t *testing.T) {
	store := openTestStore(t)
	if err := store.Migrate(); err != nil {
		t.Fatalf("second migrate failed: %v", err)
	}
}

func TestProfessionalPipelineConfigPersists(t *testing.T) {
	store := openTestStore(t)
	defaults, err := store.GetProfessionalPipelineConfig()
	if err != nil {
		t.Fatalf("get default professional config: %v", err)
	}
	if defaults.DetectorID != DetectorTireLabProblems || defaults.DecisionerID != DecisionerTire || defaults.InterpreterID != InterpreterRoadDocsV12 {
		t.Fatalf("unexpected default professional config: %+v", defaults)
	}
	saved, err := store.SaveProfessionalPipelineConfig(ProfessionalPipelineConfig{
		DetectorID:    DetectorTireLabProblems,
		DecisionerID:  DecisionerTire,
		InterpreterID: InterpreterTireRepair,
	})
	if err != nil {
		t.Fatalf("save professional config: %v", err)
	}
	loaded, err := store.GetProfessionalPipelineConfig()
	if err != nil {
		t.Fatalf("reload professional config: %v", err)
	}
	if loaded != saved {
		t.Fatalf("config did not persist: got %+v want %+v", loaded, saved)
	}
}

func TestRecommendedCarCRUD(t *testing.T) {
	store := openTestStore(t)
	created, err := store.SaveRecommendedCar(RecommendedCarInput{
		Name:         "2018 Honda Civic Type R",
		UseCase:      "Road",
		PI:           700,
		CarClass:     "b",
		Drivetrain:   "fwd",
		TireCompound: "sport",
		TuneCode:     "927164038",
		Tags:         []string{"front", "", " stable "},
		Reason:       "front drive stable",
	})
	if err != nil {
		t.Fatalf("save recommended car: %v", err)
	}
	if created.CarClass != "B" || created.Drivetrain != "FWD" || len(created.Tags) != 2 {
		t.Fatalf("created recommended car = %#v", created)
	}
	if created.ID != "honda-civic-type-r-2018-road-b700-927164038" || created.UseCaseLabel == "" || created.TireCompoundLabel == "" {
		t.Fatalf("created generated identity/labels = %#v", created)
	}
	if created.WeightKG != 0 || created.FrontWeightPct != 0 {
		t.Fatalf("optional metadata should default to zero: %#v", created)
	}
	if _, err := store.SaveRecommendedCar(RecommendedCarInput{
		ID:             created.ID,
		Name:           "2018 Honda Civic Type R Updated",
		UseCase:        "Road",
		PI:             700,
		CarClass:       "B",
		Drivetrain:     "FWD",
		TireCompound:   "sport",
		WeightKG:       1390,
		FrontWeightPct: 61,
		TuneCode:       "111222333",
		Tags:           []string{"front"},
	}); err != nil {
		t.Fatalf("update recommended car: %v", err)
	}
	cars, err := store.ListRecommendedCars()
	if err != nil {
		t.Fatalf("list recommended cars: %v", err)
	}
	if len(cars) != 1 || cars[0].Name != "2018 Honda Civic Type R Updated" || cars[0].TuneCode != "111222333" {
		t.Fatalf("listed recommended cars = %#v", cars)
	}
	if err := store.DeleteRecommendedCar(cars[0].ID); err != nil {
		t.Fatalf("delete recommended car: %v", err)
	}
	if got := countRows(t, store, "recommended_car"); got != 0 {
		t.Fatalf("recommended_car rows = %d, want 0", got)
	}
}

func TestTuneHarvestStorageCRUD(t *testing.T) {
	store := openTestStore(t)
	if _, err := store.SaveFH6Cars([]FH6CarInput{{
		Year:              2020,
		Make:              "BMW",
		Model:             "M2 Competition Coup茅",
		Aliases:           []string{"M2 Comp"},
		BasePI:            718,
		DrivetrainDefault: "rwd",
		Source:            "test",
		SourceRef:         "bmw-m2",
	}}); err != nil {
		t.Fatalf("save fh6 cars: %v", err)
	}
	cars, err := store.ListFH6Cars()
	if err != nil {
		t.Fatalf("list fh6 cars: %v", err)
	}
	if len(cars) != 1 || cars[0].CarID == "" || cars[0].DrivetrainDefault != "RWD" || len(cars[0].Aliases) == 0 {
		t.Fatalf("cars = %#v", cars)
	}
	run, err := store.CreateTuneHarvestRun(TuneHarvestRunInput{Sources: []string{"jsr", "codmunity"}, DryRun: false})
	if err != nil {
		t.Fatalf("create run: %v", err)
	}
	candidate, err := store.UpsertTuneHarvestCandidate(TuneHarvestCandidateInput{
		RunID:        run.ID,
		Source:       "jsr_chronic_sheet",
		RawKey:       "row:1",
		ShareCode:    "123 456 789",
		Year:         2020,
		Make:         "BMW",
		Model:        "M2 Comp",
		CarName:      "2020 BMW M2 Comp",
		MatchedCarID: cars[0].CarID,
		MatchScore:   0.93,
		UseCase:      "Road",
		CarClass:     "A",
		PI:           700,
		Drivetrain:   "awd",
		RawJSON:      `{"ok":true}`,
	})
	if err != nil {
		t.Fatalf("upsert candidate: %v", err)
	}
	if candidate.ShareCode != "123456789" || candidate.Drivetrain != "AWD" || candidate.Status != TuneHarvestCandidatePending {
		t.Fatalf("candidate = %#v", candidate)
	}
	listed, err := store.ListTuneHarvestCandidates(TuneHarvestCandidatePending, 10)
	if err != nil {
		t.Fatalf("list candidates: %v", err)
	}
	if len(listed) != 1 || listed[0].ID != candidate.ID {
		t.Fatalf("listed = %#v", listed)
	}
	searched, err := store.SearchTuneHarvestCandidates(TuneHarvestCandidatePending, "123-456-789", 10)
	if err != nil {
		t.Fatalf("search candidate by formatted code: %v", err)
	}
	if len(searched) != 1 || searched[0].ID != candidate.ID {
		t.Fatalf("searched by code = %#v", searched)
	}
	searched, err = store.SearchTuneHarvestCandidates(TuneHarvestCandidatePending, "BMW AWD", 10)
	if err != nil {
		t.Fatalf("search candidate by terms: %v", err)
	}
	if len(searched) != 1 || searched[0].ID != candidate.ID {
		t.Fatalf("searched by terms = %#v", searched)
	}
	searched, err = store.SearchTuneHarvestCandidates(TuneHarvestCandidatePending, "ForzaFire", 10)
	if err != nil {
		t.Fatalf("search missing candidate: %v", err)
	}
	if len(searched) != 0 {
		t.Fatalf("searched missing = %#v", searched)
	}
	duplicate, err := store.UpsertTuneHarvestCandidate(TuneHarvestCandidateInput{
		RunID:     run.ID,
		Source:    "codmunity",
		RawKey:    "codmunity:bmw-m2:123456789:road:a",
		ShareCode: "123-456-789",
		Year:      2020,
		Make:      "BMW",
		Model:     "M2 Competition",
		CarName:   "2020 BMW M2 Competition",
		Tuner:     "Ghost",
		TuneName:  "Road Meta",
		RawJSON:   `{"duplicate":true}`,
	})
	if err != nil {
		t.Fatalf("upsert duplicate share code: %v", err)
	}
	if duplicate.ID != candidate.ID {
		t.Fatalf("duplicate id = %d, want existing %d", duplicate.ID, candidate.ID)
	}
	if duplicate.Tuner != "Ghost" {
		t.Fatalf("duplicate merge did not fill tuner: %#v", duplicate)
	}
	searched, err = store.SearchTuneHarvestCandidates("all", "123456789", 10)
	if err != nil {
		t.Fatalf("search duplicate share code: %v", err)
	}
	if len(searched) != 1 || searched[0].ID != candidate.ID {
		t.Fatalf("duplicate search = %#v", searched)
	}
	updated, err := store.UpdateTuneHarvestCandidateStatus(candidate.ID, TuneHarvestCandidateImported, "")
	if err != nil {
		t.Fatalf("update candidate status: %v", err)
	}
	if updated.Status != TuneHarvestCandidateImported {
		t.Fatalf("updated = %#v", updated)
	}
	finished, err := store.FinishTuneHarvestRun(run.ID, TuneHarvestRunComplete, "", 1, 1, 0, 1, 0)
	if err != nil {
		t.Fatalf("finish run: %v", err)
	}
	if finished.Status != TuneHarvestRunComplete || finished.SavedCount != 1 {
		t.Fatalf("finished = %#v", finished)
	}
	cleared, err := store.ClearTuneHarvestCandidates()
	if err != nil {
		t.Fatalf("clear tune harvest candidates: %v", err)
	}
	if cleared != 1 {
		t.Fatalf("cleared = %d, want 1", cleared)
	}
	listed, err = store.ListTuneHarvestCandidates("all", 10)
	if err != nil {
		t.Fatalf("list after clear: %v", err)
	}
	if len(listed) != 0 {
		t.Fatalf("listed after clear = %#v", listed)
	}
}

func TestTuneHarvestCandidateListCanReturnMoreThanThreeHundred(t *testing.T) {
	store := openTestStore(t)
	for index := 0; index < 325; index++ {
		if _, err := store.UpsertTuneHarvestCandidate(TuneHarvestCandidateInput{
			Source:    "codmunity",
			RawKey:    fmt.Sprintf("codmunity:%d", index),
			ShareCode: fmt.Sprintf("800%06d", index),
			CarName:   "2020 BMW M2 Competition",
			RawJSON:   `{"ok":true}`,
		}); err != nil {
			t.Fatalf("upsert candidate %d: %v", index, err)
		}
	}
	listed, err := store.SearchTuneHarvestCandidates("all", "", 500)
	if err != nil {
		t.Fatalf("list candidates: %v", err)
	}
	if len(listed) != 325 {
		t.Fatalf("listed = %d, want 325", len(listed))
	}
}

func TestRecommendedCarRecordValidationAndIDMigration(t *testing.T) {
	store := openTestStore(t)
	first, err := store.SaveRecommendedCarRecord(RecommendedCarInput{
		Name:         "2018 Honda Civic Type R",
		UseCase:      "Road",
		PI:           700,
		CarClass:     "B",
		Drivetrain:   "FWD",
		TireCompound: "sport",
		TuneCode:     "927164038",
	}, "")
	if err != nil {
		t.Fatalf("save first recommended car: %v", err)
	}
	if _, err := store.SaveRecommendedCarRecord(RecommendedCarInput{
		Name:         "2021 Porsche 911 GT3",
		UseCase:      "Road",
		PI:           800,
		CarClass:     "S1",
		Drivetrain:   "AWD",
		TireCompound: "sport",
		TuneCode:     "927-164-038",
	}, ""); err == nil {
		t.Fatal("expected duplicate tuneCode to fail")
	}
	updated, err := store.SaveRecommendedCarRecord(RecommendedCarInput{
		Name:         "2018 Honda Civic Type R Revised",
		UseCase:      "Road",
		PI:           700,
		CarClass:     "B",
		Drivetrain:   "FWD",
		TireCompound: "sport",
		TuneCode:     first.TuneCode,
	}, first.ID)
	if err != nil {
		t.Fatalf("edit with same tuneCode: %v", err)
	}
	if updated.ID == first.ID {
		t.Fatalf("expected ID to change after name edit: %q", updated.ID)
	}
	if _, err := store.GetRecommendedCar(first.ID); err == nil {
		t.Fatalf("old recommended car ID %q still exists", first.ID)
	}
	if updated.CreatedAt != first.CreatedAt || updated.UpdatedAt == first.UpdatedAt {
		t.Fatalf("timestamps after ID migration = created %q/%q updated %q/%q", first.CreatedAt, updated.CreatedAt, first.UpdatedAt, updated.UpdatedAt)
	}
	second, err := store.SaveRecommendedCarRecord(RecommendedCarInput{
		Name:         "2021 Porsche 911 GT3",
		UseCase:      "Road",
		PI:           800,
		CarClass:     "S1",
		Drivetrain:   "AWD",
		TireCompound: "sport",
		TuneCode:     "413829605",
	}, "")
	if err != nil {
		t.Fatalf("save second recommended car: %v", err)
	}
	sameCarDifferentCode, err := store.SaveRecommendedCarRecord(RecommendedCarInput{
		Name:         second.Name,
		UseCase:      second.UseCase,
		PI:           second.PI,
		CarClass:     second.CarClass,
		Drivetrain:   second.Drivetrain,
		TireCompound: second.TireCompound,
		TuneCode:     "999999999",
	}, "")
	if err != nil {
		t.Fatalf("expected same vehicle with different tuneCode to save: %v", err)
	}
	if sameCarDifferentCode.ID == second.ID {
		t.Fatalf("same vehicle with different tuneCode reused ID %q", second.ID)
	}
	if _, err := store.SaveRecommendedCarRecord(RecommendedCarInput{
		Name:         second.Name,
		UseCase:      second.UseCase,
		PI:           801,
		CarClass:     second.CarClass,
		Drivetrain:   second.Drivetrain,
		TireCompound: second.TireCompound,
		TuneCode:     "888888888",
	}, ""); err != nil {
		t.Fatalf("expected copy with changed PI and tuneCode to save: %v", err)
	}
}

func TestDeleteAllRecommendedCars(t *testing.T) {
	store := openTestStore(t)
	for _, input := range []RecommendedCarInput{
		{
			Name:         "2018 Honda Civic Type R",
			UseCase:      "Road",
			PI:           700,
			CarClass:     "B",
			Drivetrain:   "FWD",
			TireCompound: "sport",
			TuneCode:     "927164038",
		},
		{
			Name:         "2021 Porsche 911 GT3",
			UseCase:      "Road",
			PI:           800,
			CarClass:     "S1",
			Drivetrain:   "AWD",
			TireCompound: "sport",
			TuneCode:     "413829605",
		},
	} {
		if _, err := store.SaveRecommendedCarRecord(input, ""); err != nil {
			t.Fatalf("save recommended car fixture: %v", err)
		}
	}
	deleted, err := store.DeleteAllRecommendedCars()
	if err != nil {
		t.Fatalf("delete all recommended cars: %v", err)
	}
	if deleted != 2 {
		t.Fatalf("deleted rows = %d, want 2", deleted)
	}
	if got := countRows(t, store, "recommended_car"); got != 0 {
		t.Fatalf("recommended_car rows = %d, want 0", got)
	}
	deleted, err = store.DeleteAllRecommendedCars()
	if err != nil {
		t.Fatalf("delete all recommended cars on empty table: %v", err)
	}
	if deleted != 0 {
		t.Fatalf("deleted rows on empty table = %d, want 0", deleted)
	}
}

func TestCleanupLegacySessionsClearsSessionDataOnly(t *testing.T) {
	store := openTestStore(t)
	recordingPath := filepath.Join(t.TempDir(), "legacy.fh6udp")
	if err := os.WriteFile(recordingPath, []byte("legacy"), 0644); err != nil {
		t.Fatalf("write recording fixture: %v", err)
	}
	profile, err := store.CreateTuneProfile(TuneProfileInput{CarName: "Cleanup Car", UseCase: "Road"})
	if err != nil {
		t.Fatalf("create profile: %v", err)
	}
	session, err := store.CreateTelemetrySession(SessionStartInput{
		TuneProfileID: &profile.ID,
		SessionName:   "Legacy Session",
		StartedAt:     "2026-05-18T10:00:00Z",
		RecordingPath: recordingPath,
	})
	if err != nil {
		t.Fatalf("create session: %v", err)
	}
	if _, err := store.db.Exec(`INSERT INTO detected_event(session_id, event_type, severity, evidence_json, created_at) VALUES(?, 'x', 'low', '{}', ?)`, session.ID, nowText()); err != nil {
		t.Fatalf("insert event: %v", err)
	}
	if _, err := store.db.Exec(`INSERT INTO telemetry_sample_agg(session_id, timestamp_ms, speed_kmh) VALUES(?, 100, 120)`, session.ID); err != nil {
		t.Fatalf("insert sample: %v", err)
	}
	if _, err := store.db.Exec(`INSERT INTO benchmark_track(name, polyline_json, created_at, updated_at) VALUES('Track', '[]', ?, ?)`, nowText(), nowText()); err != nil {
		t.Fatalf("insert track: %v", err)
	}
	trackID, err := lastInsertID(store.db)
	if err != nil {
		t.Fatalf("track id: %v", err)
	}
	if _, err := store.db.Exec(`INSERT INTO benchmark_run(session_id, track_id, created_at) VALUES(?, ?, ?)`, session.ID, trackID, nowText()); err != nil {
		t.Fatalf("insert benchmark run: %v", err)
	}
	if _, err := store.db.Exec(`INSERT INTO track_baseline_run(track_id, avg_speed_kmh, max_speed_kmh, created_at) VALUES(?, 100, 130, ?)`, trackID, nowText()); err != nil {
		t.Fatalf("insert track baseline: %v", err)
	}
	if _, err := store.db.Exec(`INSERT INTO tune_change_log(tune_profile_id, session_id, changed_at, change_reason, change_json) VALUES(?, ?, ?, 'test', '{}')`, profile.ID, session.ID, nowText()); err != nil {
		t.Fatalf("insert snapshot: %v", err)
	}
	if err := store.CleanupLegacySessions(); err != nil {
		t.Fatalf("cleanup legacy sessions: %v", err)
	}
	for table, want := range map[string]int{
		"telemetry_session":    0,
		"detected_event":       0,
		"telemetry_sample_agg": 0,
		"benchmark_run":        0,
		"tune_profile":         1,
		"benchmark_track":      1,
		"track_baseline_run":   1,
	} {
		if got := countRows(t, store, table); got != want {
			t.Fatalf("%s rows = %d, want %d", table, got, want)
		}
	}
	var sessionID any
	if err := store.db.QueryRow(`SELECT session_id FROM tune_change_log LIMIT 1`).Scan(&sessionID); err != nil {
		t.Fatalf("read snapshot session id: %v", err)
	}
	if sessionID != nil {
		t.Fatalf("snapshot session_id was not cleared: %v", sessionID)
	}
	if _, err := os.Stat(recordingPath); !os.IsNotExist(err) {
		t.Fatalf("recording file was not deleted, stat err=%v", err)
	}
}

func TestTuneProfileCRUDDuplicateAndActive(t *testing.T) {
	store := openTestStore(t)
	pressure := 28.5
	pi := int64(900)
	carOrdinal := int64(123456)
	carCategory := int64(12)
	numCylinders := int64(6)
	created, err := store.CreateTuneProfile(TuneProfileInput{
		CarName:           "BMW M3",
		CarOrdinal:        &carOrdinal,
		CarCategory:       &carCategory,
		CarClass:          "S1",
		PI:                &pi,
		Drivetrain:        "RWD",
		NumCylinders:      &numCylinders,
		UseCase:           "Road",
		VersionName:       "Base",
		FrontTirePressure: &pressure,
	})
	if err != nil {
		t.Fatalf("create profile: %v", err)
	}
	if created.ID == 0 || created.FrontTirePressure == nil || *created.FrontTirePressure != pressure {
		t.Fatalf("created profile = %#v", created)
	}
	if created.CarOrdinal == nil || *created.CarOrdinal != carOrdinal || created.CarCategory == nil || *created.CarCategory != carCategory || created.NumCylinders == nil || *created.NumCylinders != numCylinders {
		t.Fatalf("created vehicle metadata = %#v", created)
	}

	if err := store.SetActiveTuneProfile(created.ID); err != nil {
		t.Fatalf("set active: %v", err)
	}
	active, err := store.GetActiveTuneProfile()
	if err != nil {
		t.Fatalf("get active: %v", err)
	}
	if active == nil || active.ID != created.ID {
		t.Fatalf("active = %#v, want id %d", active, created.ID)
	}

	updatedInput := created.ToInput()
	updatedInput.CarName = "BMW M3 Updated"
	updated, err := store.UpdateTuneProfile(created.ID, updatedInput)
	if err != nil {
		t.Fatalf("update profile: %v", err)
	}
	if updated.CarName != "BMW M3 Updated" {
		t.Fatalf("updated car name = %q", updated.CarName)
	}

	duplicate, err := store.DuplicateTuneProfile(created.ID, "Race")
	if err != nil {
		t.Fatalf("duplicate profile: %v", err)
	}
	if duplicate.ID == created.ID || duplicate.VersionName != "Race" {
		t.Fatalf("duplicate = %#v", duplicate)
	}

	profiles, err := store.ListTuneProfiles()
	if err != nil {
		t.Fatalf("list profiles: %v", err)
	}
	if len(profiles) != 2 {
		t.Fatalf("profile count = %d, want 2", len(profiles))
	}

	if err := store.DeleteTuneProfile(created.ID); err != nil {
		t.Fatalf("delete active profile: %v", err)
	}
	active, err = store.GetActiveTuneProfile()
	if err != nil {
		t.Fatalf("get active after delete: %v", err)
	}
	if active != nil {
		t.Fatalf("active after delete = %#v, want nil", active)
	}
}

func TestTuneProfilePowerToWeightCalculated(t *testing.T) {
	store := openTestStore(t)
	power := 320.0
	torque := 540.0
	weight := 1425.0
	frontWeight := 52.5
	peakTorqueRPM := 4200.0
	peakPowerRPM := 6800.0
	redlineRPM := 7200.0
	created, err := store.CreateTuneProfile(TuneProfileInput{
		CarName:        "Power Car",
		PowerKW:        &power,
		TorqueNM:       &torque,
		WeightKG:       &weight,
		FrontWeightPct: &frontWeight,
		PeakTorqueRPM:  &peakTorqueRPM,
		PeakPowerRPM:   &peakPowerRPM,
		RedlineRPM:     &redlineRPM,
	})
	if err != nil {
		t.Fatalf("create power profile: %v", err)
	}
	if created.PowerToWeightKWPerKG == nil || math.Abs(*created.PowerToWeightKWPerKG-0.2246) > 0.00001 {
		t.Fatalf("power to weight = %#v, want 0.2246", created.PowerToWeightKWPerKG)
	}
	if created.PeakTorqueRPM == nil || *created.PeakTorqueRPM != peakTorqueRPM || created.PeakPowerRPM == nil || *created.PeakPowerRPM != peakPowerRPM || created.RedlineRPM == nil || *created.RedlineRPM != redlineRPM {
		t.Fatalf("rpm band fields = %#v/%#v/%#v, want persisted values", created.PeakTorqueRPM, created.PeakPowerRPM, created.RedlineRPM)
	}
	weight = 1600
	peakPowerRPM = 7000
	input := created.ToInput()
	input.WeightKG = &weight
	input.PeakPowerRPM = &peakPowerRPM
	updated, err := store.UpdateTuneProfile(created.ID, input)
	if err != nil {
		t.Fatalf("update power profile: %v", err)
	}
	if updated.PowerToWeightKWPerKG == nil || math.Abs(*updated.PowerToWeightKWPerKG-0.2) > 0.00001 {
		t.Fatalf("updated power to weight = %#v, want 0.2000", updated.PowerToWeightKWPerKG)
	}
	if updated.PeakPowerRPM == nil || *updated.PeakPowerRPM != peakPowerRPM {
		t.Fatalf("updated peak power rpm = %#v, want %f", updated.PeakPowerRPM, peakPowerRPM)
	}
}

func TestTuneProfileSnapshotsUpdatePruneAndRestore(t *testing.T) {
	store := openTestStore(t)
	pressure := 28.0
	profile, err := store.CreateTuneProfile(TuneProfileInput{
		CarName:           "Snapshot Car",
		CarClass:          "A",
		UseCase:           "Road",
		FrontTirePressure: &pressure,
	})
	if err != nil {
		t.Fatalf("create profile: %v", err)
	}

	values := []float64{28.5, 29.0, 29.5, 30.0, 30.5, 31.0}
	names := []string{"v1", "v2", "v3", "v4", "v5", "v6"}
	for i, value := range values {
		input := profile.ToInput()
		input.FrontTirePressure = &value
		input.VersionName = names[i]
		profile, err = store.UpdateTuneProfile(profile.ID, input)
		if err != nil {
			t.Fatalf("update profile %d: %v", i, err)
		}
	}

	snapshots, err := store.ListTuneProfileSnapshots(profile.ID)
	if err != nil {
		t.Fatalf("list snapshots: %v", err)
	}
	if len(snapshots) != 5 {
		t.Fatalf("snapshot count = %d, want 5", len(snapshots))
	}
	if snapshots[0].Before.FrontTirePressure == nil || snapshots[0].After.FrontTirePressure == nil {
		t.Fatalf("snapshot missing before/after pressure: %#v", snapshots[0])
	}
	if *snapshots[0].Before.FrontTirePressure != 30.5 || *snapshots[0].After.FrontTirePressure != 31.0 {
		t.Fatalf("latest snapshot before/after = %.1f/%.1f, want 30.5/31.0", *snapshots[0].Before.FrontTirePressure, *snapshots[0].After.FrontTirePressure)
	}
	if !containsString(snapshots[0].ChangedFields, "frontTirePressure") {
		t.Fatalf("changed fields = %#v, want frontTirePressure", snapshots[0].ChangedFields)
	}

	restoreTarget := snapshots[len(snapshots)-1]
	if restoreTarget.Before.FrontTirePressure == nil {
		t.Fatalf("restore target missing before value: %#v", restoreTarget)
	}
	wantPressure := *restoreTarget.Before.FrontTirePressure
	restored, err := store.RestoreTuneProfileSnapshot(restoreTarget.ID)
	if err != nil {
		t.Fatalf("restore snapshot: %v", err)
	}
	if restored.FrontTirePressure == nil || *restored.FrontTirePressure != wantPressure {
		t.Fatalf("restored pressure = %#v, want %.1f", restored.FrontTirePressure, wantPressure)
	}
	snapshots, err = store.ListTuneProfileSnapshots(profile.ID)
	if err != nil {
		t.Fatalf("list snapshots after restore: %v", err)
	}
	if len(snapshots) != 5 {
		t.Fatalf("snapshot count after restore = %d, want 5", len(snapshots))
	}
	if snapshots[0].ChangeReason != "restore" {
		t.Fatalf("latest snapshot reason = %q, want restore", snapshots[0].ChangeReason)
	}
}

func TestCreateTelemetrySessionStoresTuneProfileSnapshot(t *testing.T) {
	store := openTestStore(t)
	profile, err := store.CreateTuneProfile(TuneProfileInput{CarName: "Snapshot Car", CarClass: "A", UseCase: "Road"})
	if err != nil {
		t.Fatalf("create profile: %v", err)
	}
	snapshotJSON, err := TuneProfileSnapshotJSON(profile)
	if err != nil {
		t.Fatalf("snapshot json: %v", err)
	}
	session, err := store.CreateTelemetrySession(SessionStartInput{TuneProfileID: &profile.ID, TuneSnapshotJSON: snapshotJSON, SessionName: "Snapshot Session", StartedAt: "2026-05-18T10:00:00Z"})
	if err != nil {
		t.Fatalf("create session: %v", err)
	}
	if !strings.Contains(session.TuneSnapshotJSON, "Snapshot Car") {
		t.Fatalf("session tune snapshot json = %q", session.TuneSnapshotJSON)
	}
	parsed, err := ParseTuneProfileSnapshotJSON(session.TuneSnapshotJSON)
	if err != nil {
		t.Fatalf("parse session snapshot: %v", err)
	}
	if parsed == nil || parsed.ID != profile.ID || parsed.CarName != profile.CarName {
		t.Fatalf("parsed snapshot = %#v, want profile %#v", parsed, profile)
	}
}

func TestTuneAdjustmentExplanationsForAction(t *testing.T) {
	explanations := TuneAdjustmentExplanationsForAction("brake_balance", 0)
	if len(explanations) != 1 || explanations[0].Detail == "" {
		t.Fatalf("brake balance explanations = %#v", explanations)
	}
	gear := TuneAdjustmentExplanationsForAction("current_gear", 3)
	if len(gear) != 1 || gear[0].Detail == "" {
		t.Fatalf("current gear explanations = %#v", gear)
	}
}

func TestSessionFinalizeStoresEvents(t *testing.T) {
	store := openTestStore(t)
	profile, err := store.CreateTuneProfile(TuneProfileInput{CarName: "Test Car"})
	if err != nil {
		t.Fatalf("create profile: %v", err)
	}
	session, err := store.CreateTelemetrySession(SessionStartInput{TuneProfileID: &profile.ID, SessionName: "Test Session", StartedAt: "2026-05-18T10:00:00Z"})
	if err != nil {
		t.Fatalf("create session: %v", err)
	}
	avg := 123.4
	max := 188.8
	carOrdinal := int64(123456)
	carPI := int64(860)
	numCylinders := int64(6)
	_, err = store.FinalizeTelemetrySession(SessionFinalizeInput{
		SessionID:          session.ID,
		EndedAt:            "2026-05-18T10:01:00Z",
		DurationMS:         60000,
		AvgSpeedKmh:        &avg,
		MaxSpeedKmh:        &max,
		RecordingPackets:   2,
		RecordingBytes:     700,
		RecordingTruncated: true,
		GameMode:           telemetry.GameModeRace,
		SessionVehicleSnapshot: SessionVehicleSnapshot{
			CarOrdinal:   &carOrdinal,
			CarClass:     "S1",
			CarPI:        &carPI,
			Drivetrain:   "AWD",
			NumCylinders: &numCylinders,
		},
	}, []telemetry.DetectedEvent{
		{
			Type:             "corner_exit_oversteer",
			Severity:         "high",
			Segment:          "corner_exit",
			StartMS:          1000,
			EndMS:            1600,
			DurationMS:       600,
			Evidence:         map[string]float64{"rear_combined_slip": 1.4},
			SuggestedActions: []telemetry.SuggestedAction{{Priority: 0, Category: "differential", Item: "rear_diff_accel", Direction: "decrease", Amount: "3%-5%", Reason: "reduce power oversteer"}},
		},
	}, []telemetry.NormalizedTelemetry{
		{TimeMS: 1000, IsRaceOn: true, GameMode: telemetry.GameModeRace, SpeedKmh: 120, SpeedFieldKmh: 118, VelocitySpeedKmh: 120, SpeedSource: "velocity", Rpm: 5100, RpmRatio: 0.7, Gear: 3, Throttle01: 0.8, Brake01: 0.1, Steer01: -0.2, FrontCombinedSlipAvg: 0.5, RearCombinedSlipAvg: 0.7, PositionX: 10, PositionY: 2, PositionZ: 20, DistanceTraveled: 33, BestLap: 88.5, LastLap: 91.2, CurrentLap: 12.5, CurrentRaceTime: 123.4, LapNumber: 2, RacePosition: 1, DrivingLine01: 0.35, CarOrdinal: 123456, CarClass: "S1", CarPI: 860, Drivetrain: "AWD", NumCylinders: 6},
	})
	if err != nil {
		t.Fatalf("finalize session: %v", err)
	}

	savedSession, err := store.GetTelemetrySession(session.ID)
	if err != nil {
		t.Fatalf("get session: %v", err)
	}
	if savedSession.EventCount != 1 || savedSession.SampleCount != 1 || savedSession.AvgSpeedKmh == nil || *savedSession.AvgSpeedKmh != avg || savedSession.RecordingPackets != 2 || savedSession.RecordingBytes != 700 || !savedSession.RecordingTruncated {
		t.Fatalf("saved session = %#v", savedSession)
	}
	if savedSession.CarOrdinal == nil || *savedSession.CarOrdinal != carOrdinal || savedSession.CarClass != "S1" || savedSession.CarPI == nil || *savedSession.CarPI != carPI || savedSession.Drivetrain != "AWD" || savedSession.NumCylinders == nil || *savedSession.NumCylinders != numCylinders {
		t.Fatalf("saved vehicle snapshot = %#v", savedSession)
	}
	if savedSession.GameMode != telemetry.GameModeRace {
		t.Fatalf("saved game mode = %q", savedSession.GameMode)
	}

	events, err := store.GetSessionEvents(session.ID)
	if err != nil {
		t.Fatalf("get events: %v", err)
	}
	if len(events) != 1 || events[0].Type != "corner_exit_oversteer" || events[0].Evidence["rear_combined_slip"] != 1.4 {
		t.Fatalf("events = %#v", events)
	}

	samples, err := store.GetSessionTelemetrySamples(session.ID, 10)
	if err != nil {
		t.Fatalf("get samples: %v", err)
	}
	if len(samples) != 1 || samples[0].SpeedKmh != 120 || samples[0].SpeedSource != "velocity" || samples[0].CarOrdinal != 123456 || samples[0].GameMode != telemetry.GameModeRace || !samples[0].IsRaceOn {
		t.Fatalf("samples = %#v", samples)
	}
	if samples[0].PositionX != 10 || samples[0].PositionY != 2 || samples[0].PositionZ != 20 || samples[0].DistanceTraveled != 33 || samples[0].BestLap != 88.5 || samples[0].LapNumber != 2 || samples[0].RacePosition != 1 || samples[0].DrivingLine01 != 0.35 {
		t.Fatalf("sample route fields = %#v", samples[0])
	}
}

func TestGetSessionEventsNormalizesHistoricalGearSuggestionDirections(t *testing.T) {
	store := openTestStore(t)
	session, err := store.CreateTelemetrySession(SessionStartInput{SessionName: "Historical Suggestions", StartedAt: "2026-05-18T10:00:00Z"})
	if err != nil {
		t.Fatalf("create session: %v", err)
	}
	_, err = store.FinalizeTelemetrySession(SessionFinalizeInput{
		SessionID:  session.ID,
		EndedAt:    "2026-05-18T10:01:00Z",
		DurationMS: 60000,
	}, []telemetry.DetectedEvent{
		{
			Type:       "launch_wheelspin",
			Severity:   "medium",
			Segment:    "launch",
			DurationMS: 500,
			Evidence:   map[string]float64{"gear": 1, "rear_slip_ratio": 1.3},
			SuggestedActions: []telemetry.SuggestedAction{
				{Priority: 0, Category: "gearing", Item: "gear_1", Direction: "increase", Amount: "5%-10%", Reason: "reduce wheel torque during launch"},
			},
		},
	}, nil)
	if err != nil {
		t.Fatalf("finalize session: %v", err)
	}
	events, err := store.GetSessionEvents(session.ID)
	if err != nil {
		t.Fatalf("get events: %v", err)
	}
	action := events[0].SuggestedActions[0]
	if action.Direction != "decrease" || action.Amount != "3%-5%" {
		t.Fatalf("normalized action = %#v, want decrease 3%%-5%%", action)
	}
}

func containsString(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

func TestTestConditionDefaultsAndSessionStorage(t *testing.T) {
	store := openTestStore(t)
	defaults, err := store.GetTestConditionDefaults()
	if err != nil {
		t.Fatalf("get default conditions: %v", err)
	}
	if defaults.DriverMode != "unknown" || defaults.BrakeAssist != "unknown" || defaults.LaunchControl != "unknown" {
		t.Fatalf("defaults = %#v", defaults)
	}

	saved, err := store.SaveTestConditionDefaults(TestConditions{
		DriverMode:       "PLAYER",
		BrakeAssist:      "abs_on",
		SteeringAssist:   "simulation",
		TractionControl:  "off",
		StabilityControl: "invalid",
		Shifting:         "manual",
		LaunchControl:    "on",
	})
	if err != nil {
		t.Fatalf("save default conditions: %v", err)
	}
	if saved.DriverMode != "player" || saved.StabilityControl != "unknown" || saved.LaunchControl != "on" {
		t.Fatalf("saved normalized conditions = %#v", saved)
	}
	loaded, err := store.GetTestConditionDefaults()
	if err != nil {
		t.Fatalf("load default conditions: %v", err)
	}
	if loaded != saved {
		t.Fatalf("loaded defaults = %#v, want %#v", loaded, saved)
	}

	session, err := store.CreateTelemetrySession(SessionStartInput{
		SessionName: "Conditioned",
		StartedAt:   "2026-05-18T10:00:00Z",
		TestConditions: TestConditions{
			DriverMode:       "auto",
			BrakeAssist:      "abs_off",
			SteeringAssist:   "standard",
			TractionControl:  "on",
			StabilityControl: "off",
			Shifting:         "automatic",
			LaunchControl:    "off",
		},
	})
	if err != nil {
		t.Fatalf("create conditioned session: %v", err)
	}
	if session.DriverMode != "unknown" || session.BrakeAssist != "abs_off" || session.SteeringAssist != "standard" || session.TractionControl != "on" || session.StabilityControl != "off" || session.Shifting != "automatic" || session.LaunchControl != "off" {
		t.Fatalf("session conditions = %#v", session)
	}
}

func TestListTuneProfilesForVehicle(t *testing.T) {
	store := openTestStore(t)
	ordinal := int64(5001)
	otherOrdinal := int64(5002)
	for _, input := range []TuneProfileInput{
		{CarName: "Road", CarOrdinal: &ordinal, CarClass: "S1"},
		{CarName: "Rally", CarOrdinal: &ordinal, CarClass: "S1"},
		{CarName: "Lower Class", CarOrdinal: &ordinal, CarClass: "A"},
		{CarName: "Other", CarOrdinal: &otherOrdinal, CarClass: "S1"},
	} {
		if _, err := store.CreateTuneProfile(input); err != nil {
			t.Fatalf("create profile: %v", err)
		}
	}
	matches, err := store.ListTuneProfilesForVehicle(ordinal, "S1")
	if err != nil {
		t.Fatalf("list candidates: %v", err)
	}
	if len(matches) != 2 {
		t.Fatalf("matches = %#v, want two S1 profiles for ordinal", matches)
	}
}

func TestBindTelemetrySessionTuneProfileRequiresVehicleMatch(t *testing.T) {
	store := openTestStore(t)
	ordinal := int64(7001)
	otherOrdinal := int64(7002)
	matchingA, err := store.CreateTuneProfile(TuneProfileInput{CarName: "Road", CarOrdinal: &ordinal, CarClass: "S1"})
	if err != nil {
		t.Fatalf("create matching profile: %v", err)
	}
	matchingB, err := store.CreateTuneProfile(TuneProfileInput{CarName: "Rally", CarOrdinal: &ordinal, CarClass: "s1"})
	if err != nil {
		t.Fatalf("create second matching profile: %v", err)
	}
	classMismatch, err := store.CreateTuneProfile(TuneProfileInput{CarName: "Class Mismatch", CarOrdinal: &ordinal, CarClass: "A"})
	if err != nil {
		t.Fatalf("create class mismatch profile: %v", err)
	}
	ordinalMismatch, err := store.CreateTuneProfile(TuneProfileInput{CarName: "Ordinal Mismatch", CarOrdinal: &otherOrdinal, CarClass: "S1"})
	if err != nil {
		t.Fatalf("create ordinal mismatch profile: %v", err)
	}
	missingIdentity, err := store.CreateTuneProfile(TuneProfileInput{CarName: "Incomplete"})
	if err != nil {
		t.Fatalf("create incomplete profile: %v", err)
	}

	session, err := store.CreateTelemetrySession(SessionStartInput{
		SessionName: "Unbound",
		StartedAt:   "2026-05-18T10:00:00Z",
		SessionVehicleSnapshot: SessionVehicleSnapshot{
			CarOrdinal: &ordinal,
			CarClass:   "S1",
		},
	})
	if err != nil {
		t.Fatalf("create session: %v", err)
	}

	bound, err := store.BindTelemetrySessionTuneProfile(session.ID, matchingA.ID)
	if err != nil {
		t.Fatalf("bind matching profile: %v", err)
	}
	if bound.TuneProfileID == nil || *bound.TuneProfileID != matchingA.ID || bound.TuneName != "Road" {
		t.Fatalf("bound session = %#v", bound)
	}

	rebound, err := store.BindTelemetrySessionTuneProfile(session.ID, matchingB.ID)
	if err != nil {
		t.Fatalf("rebind matching profile: %v", err)
	}
	if rebound.TuneProfileID == nil || *rebound.TuneProfileID != matchingB.ID || rebound.TuneName != "Rally" {
		t.Fatalf("rebound session = %#v", rebound)
	}

	for _, profile := range []*TuneProfile{classMismatch, ordinalMismatch, missingIdentity} {
		if _, err := store.BindTelemetrySessionTuneProfile(session.ID, profile.ID); err == nil {
			t.Fatalf("expected binding %q to fail", profile.CarName)
		}
	}
	unchanged, err := store.GetTelemetrySession(session.ID)
	if err != nil {
		t.Fatalf("get unchanged session: %v", err)
	}
	if unchanged.TuneProfileID == nil || *unchanged.TuneProfileID != matchingB.ID {
		t.Fatalf("session changed after failed bind = %#v", unchanged)
	}

	noSnapshot, err := store.CreateTelemetrySession(SessionStartInput{SessionName: "No Snapshot", StartedAt: "2026-05-18T10:02:00Z"})
	if err != nil {
		t.Fatalf("create no snapshot session: %v", err)
	}
	if _, err := store.BindTelemetrySessionTuneProfile(noSnapshot.ID, matchingA.ID); err == nil {
		t.Fatal("expected session without vehicle snapshot to fail")
	}
}

func TestBenchmarkTrackFromSessionAndRunAnalysis(t *testing.T) {
	store := openTestStore(t)
	session, err := store.CreateTelemetrySession(SessionStartInput{SessionName: "Route Source", StartedAt: "2026-05-18T10:00:00Z", GameMode: telemetry.GameModeFreeRoam})
	if err != nil {
		t.Fatalf("create session: %v", err)
	}
	_, err = store.FinalizeTelemetrySession(SessionFinalizeInput{
		SessionID:  session.ID,
		EndedAt:    "2026-05-18T10:03:00Z",
		DurationMS: 180000,
		GameMode:   telemetry.GameModeFreeRoam,
	}, nil, straightRouteSamples(0, 100, 11))
	if err != nil {
		t.Fatalf("finalize source session: %v", err)
	}

	track, err := store.CreateBenchmarkTrackFromSession(session.ID, "Festival Sprint")
	if err != nil {
		t.Fatalf("create benchmark track: %v", err)
	}
	if track.ID == 0 || track.Name != "Festival Sprint" || track.RouteLengthMeters < 95 || len(track.Polyline) < 2 || !track.HasDrivingLine {
		t.Fatalf("track = %#v", track)
	}

	runs, err := store.AnalyzeSessionBenchmarkRuns(session.ID)
	if err != nil {
		t.Fatalf("analyze benchmark runs: %v", err)
	}
	if len(runs) != 1 || !runs[0].Valid || runs[0].TrackID != track.ID || runs[0].DurationMS <= 0 || runs[0].Confidence < 0.7 {
		t.Fatalf("runs = %#v", runs)
	}

	trackRuns, err := store.ListBenchmarkRuns(track.ID, 10)
	if err != nil {
		t.Fatalf("list track runs: %v", err)
	}
	if len(trackRuns) != 1 || trackRuns[0].SessionID != session.ID {
		t.Fatalf("track runs = %#v", trackRuns)
	}
}

func TestBenchmarkCircuitExtractionUsesSingleLapFromTwoLapSession(t *testing.T) {
	store := openTestStore(t)
	session, err := store.CreateTelemetrySession(SessionStartInput{SessionName: "Two Lap Circuit", StartedAt: "2026-05-18T10:10:00Z", GameMode: telemetry.GameModeFreeRoam})
	if err != nil {
		t.Fatalf("create session: %v", err)
	}
	samples := circuitRouteSamples(75, 2, 72)
	if _, err := store.FinalizeTelemetrySession(SessionFinalizeInput{SessionID: session.ID, EndedAt: "2026-05-18T10:18:00Z", DurationMS: 480000, GameMode: telemetry.GameModeFreeRoam}, nil, samples); err != nil {
		t.Fatalf("finalize circuit: %v", err)
	}

	track, err := store.ExtractBenchmarkTrackFromSession(BenchmarkTrackExtractionInput{
		SessionID:      session.ID,
		Name:           "Two Lap Extract",
		TrackType:      benchmarkTrackTypeCircuit,
		ExtractionMode: benchmarkExtractionFirstLap,
	})
	if err != nil {
		t.Fatalf("extract circuit: %v", err)
	}
	wantLapLength := 2 * math.Pi * 75
	if track.TrackType != benchmarkTrackTypeCircuit || track.LapCountObserved < 2 {
		t.Fatalf("track type/laps = %q/%d", track.TrackType, track.LapCountObserved)
	}
	if track.RouteLengthMeters < wantLapLength*0.85 || track.RouteLengthMeters > wantLapLength*1.15 {
		t.Fatalf("route length = %.1f, want around %.1f", track.RouteLengthMeters, wantLapLength)
	}
	if track.RouteLengthMeters > wantLapLength*1.5 {
		t.Fatalf("route length includes more than one lap: %.1f", track.RouteLengthMeters)
	}

	runs, err := store.AnalyzeSessionBenchmarkRuns(session.ID)
	if err != nil {
		t.Fatalf("analyze circuit: %v", err)
	}
	if len(runs) < 2 {
		t.Fatalf("runs = %#v, want two circuit laps", runs)
	}
	for _, run := range runs {
		if !run.Valid || run.TrackID != track.ID || run.DurationMS <= 0 {
			t.Fatalf("invalid circuit run = %#v", run)
		}
		if run.RouteProgress01 == nil || *run.RouteProgress01 < 0.82 {
			t.Fatalf("route progress = %#v", run.RouteProgress01)
		}
		if run.GeometryLengthMeters == nil || math.Abs(*run.GeometryLengthMeters-track.RouteLengthMeters) > track.RouteLengthMeters*0.20 {
			t.Fatalf("geometry length = %#v, track length %.1f", run.GeometryLengthMeters, track.RouteLengthMeters)
		}
		if run.DistanceTraveledDeltaMeters == nil || !strings.Contains(run.WarningFlags, "distance_traveled_mismatch") {
			t.Fatalf("distance warning not set on run = %#v", run)
		}
	}
}

func TestBenchmarkRunAnalysisRejectsReverseAndIncompleteRoute(t *testing.T) {
	store := openTestStore(t)
	track, err := store.CreateBenchmarkTrack(BenchmarkTrackInput{
		Name:           "Straight",
		SourceMode:     telemetry.GameModeFreeRoam,
		StartRadius:    10,
		EndRadius:      10,
		Polyline:       []BenchmarkPoint{{X: 0, Z: 0}, {X: 25, Z: 0}, {X: 50, Z: 0}, {X: 75, Z: 0}, {X: 100, Z: 0}},
		HasDrivingLine: true,
	})
	if err != nil {
		t.Fatalf("create track: %v", err)
	}

	reverse, err := store.CreateTelemetrySession(SessionStartInput{SessionName: "Reverse", StartedAt: "2026-05-18T10:05:00Z"})
	if err != nil {
		t.Fatalf("create reverse session: %v", err)
	}
	if _, err := store.FinalizeTelemetrySession(SessionFinalizeInput{SessionID: reverse.ID, EndedAt: "2026-05-18T10:06:00Z", DurationMS: 60000, GameMode: telemetry.GameModeFreeRoam}, nil, straightRouteSamples(100, 0, 11)); err != nil {
		t.Fatalf("finalize reverse: %v", err)
	}
	runs, err := store.AnalyzeSessionBenchmarkRuns(reverse.ID)
	if err != nil {
		t.Fatalf("analyze reverse: %v", err)
	}
	if len(runs) != 0 {
		t.Fatalf("reverse runs = %#v, want none for track %#v", runs, track)
	}

	incomplete, err := store.CreateTelemetrySession(SessionStartInput{SessionName: "Incomplete", StartedAt: "2026-05-18T10:07:00Z"})
	if err != nil {
		t.Fatalf("create incomplete session: %v", err)
	}
	if _, err := store.FinalizeTelemetrySession(SessionFinalizeInput{SessionID: incomplete.ID, EndedAt: "2026-05-18T10:08:00Z", DurationMS: 60000, GameMode: telemetry.GameModeFreeRoam}, nil, straightRouteSamples(0, 60, 7)); err != nil {
		t.Fatalf("finalize incomplete: %v", err)
	}
	runs, err = store.AnalyzeSessionBenchmarkRuns(incomplete.ID)
	if err != nil {
		t.Fatalf("analyze incomplete: %v", err)
	}
	if len(runs) != 0 {
		t.Fatalf("incomplete runs = %#v, want none", runs)
	}
}

func TestBenchmarkListsReturnEmptySlices(t *testing.T) {
	store := openTestStore(t)
	session, err := store.CreateTelemetrySession(SessionStartInput{SessionName: "No Tracks", StartedAt: "2026-05-18T10:09:00Z"})
	if err != nil {
		t.Fatalf("create session: %v", err)
	}

	tracks, err := store.ListBenchmarkTracks()
	if err != nil {
		t.Fatalf("list tracks: %v", err)
	}
	if tracks == nil || len(tracks) != 0 {
		t.Fatalf("tracks = %#v, want non-nil empty slice", tracks)
	}

	runs, err := store.AnalyzeSessionBenchmarkRuns(session.ID)
	if err != nil {
		t.Fatalf("analyze runs: %v", err)
	}
	if runs == nil || len(runs) != 0 {
		t.Fatalf("runs = %#v, want non-nil empty slice", runs)
	}
}

func TestTrackProfileGroupsAutoBaselinesByVehicle(t *testing.T) {
	store := openTestStore(t)
	track, err := store.CreateBenchmarkTrack(BenchmarkTrackInput{
		Name:              "Road Baseline",
		SourceMode:        telemetry.GameModeRace,
		TrackType:         benchmarkTrackTypeSprint,
		StartRadius:       10,
		EndRadius:         10,
		RouteLengthMeters: 1000,
		Polyline:          []BenchmarkPoint{{X: 0, Z: 0}, {X: 1000, Z: 0}},
		HasDrivingLine:    true,
	})
	if err != nil {
		t.Fatalf("create track: %v", err)
	}

	autoSlow := createTrackProfileSession(t, store, "Auto Slow", 42, "A", 700, "AWD")
	autoFast := createTrackProfileSession(t, store, "Auto Fast", 42, "A", 700, "AWD")
	autoOtherPI := createTrackProfileSession(t, store, "Auto Other PI", 42, "A", 701, "AWD")
	player := createTrackProfileSession(t, store, "Player", 42, "A", 700, "AWD")
	invalid := createTrackProfileSession(t, store, "Invalid Auto", 42, "A", 700, "AWD")

	insertTrackProfileRun(t, store, track.ID, autoSlow.ID, 62000, driverModeAuto, 0.86, true)
	insertTrackProfileRun(t, store, track.ID, autoFast.ID, 58000, driverModeAuto, 0.92, true)
	insertTrackProfileRun(t, store, track.ID, autoOtherPI.ID, 64000, driverModeAuto, 0.9, true)
	insertTrackProfileRun(t, store, track.ID, player.ID, 52000, driverModePlayer, 0.9, true)
	insertTrackProfileRun(t, store, track.ID, invalid.ID, 50000, driverModeAuto, 0.95, false)

	profile, err := store.GetTrackProfile(track.ID)
	if err != nil {
		t.Fatalf("get track profile: %v", err)
	}
	if len(profile.AutoBaselines) != 2 {
		t.Fatalf("auto baselines = %#v, want 2 vehicle groups", profile.AutoBaselines)
	}
	first := profile.AutoBaselines[0]
	if first.BestRun.Run.SessionID != autoFast.ID || first.BestRun.Run.DurationMS != 58000 {
		t.Fatalf("best baseline = %#v, want auto fast", first.BestRun.Run)
	}
	if first.RunCount != 2 {
		t.Fatalf("run count = %d, want 2", first.RunCount)
	}
	if first.Vehicle.CarOrdinal == nil || *first.Vehicle.CarOrdinal != 42 || first.Vehicle.CarPI == nil || *first.Vehicle.CarPI != 700 {
		t.Fatalf("vehicle key = %#v", first.Vehicle)
	}
	if len(profile.VehicleReferences) != 2 || profile.VehicleReferences[0].BestAutoBaseline == nil {
		t.Fatalf("vehicle references = %#v, want grouped vehicle references with auto baseline", profile.VehicleReferences)
	}
}

func TestTrackProfileNoAutoBaselineWarning(t *testing.T) {
	store := openTestStore(t)
	track, err := store.CreateBenchmarkTrack(BenchmarkTrackInput{
		Name:           "No Baseline",
		SourceMode:     telemetry.GameModeRace,
		StartRadius:    10,
		EndRadius:      10,
		Polyline:       []BenchmarkPoint{{X: 0, Z: 0}, {X: 100, Z: 0}},
		HasDrivingLine: true,
	})
	if err != nil {
		t.Fatalf("create track: %v", err)
	}
	player := createTrackProfileSession(t, store, "Player Only", 43, "S1", 900, "RWD")
	insertTrackProfileRun(t, store, track.ID, player.ID, 60000, driverModePlayer, 0.9, true)

	profiles, err := store.ListTrackProfiles()
	if err != nil {
		t.Fatalf("list track profiles: %v", err)
	}
	if len(profiles) != 1 {
		t.Fatalf("profiles = %#v", profiles)
	}
	if len(profiles[0].AutoBaselines) != 0 || !containsString(profiles[0].Warnings, "no_auto_baseline") {
		t.Fatalf("profile = %#v, want empty baseline warning", profiles[0])
	}
}

func TestBenchmarkCircuitExtractionUsesLapNumberBoundaries(t *testing.T) {
	store := openTestStore(t)
	session, err := store.CreateTelemetrySession(SessionStartInput{SessionName: "Offset Lap Circuit", StartedAt: "2026-05-18T10:20:00Z", GameMode: telemetry.GameModeRace})
	if err != nil {
		t.Fatalf("create session: %v", err)
	}
	radius := 80.0
	samples := circuitRouteSamplesWithLapNumbers(radius, 3, 72, 18)
	if _, err := store.FinalizeTelemetrySession(SessionFinalizeInput{SessionID: session.ID, EndedAt: "2026-05-18T10:28:00Z", DurationMS: 480000, GameMode: telemetry.GameModeRace}, nil, samples); err != nil {
		t.Fatalf("finalize circuit: %v", err)
	}

	track, err := store.ExtractBenchmarkTrackFromSession(BenchmarkTrackExtractionInput{
		SessionID:      session.ID,
		Name:           "Lap Boundary Extract",
		TrackType:      benchmarkTrackTypeCircuit,
		ExtractionMode: benchmarkExtractionFirstLap,
	})
	if err != nil {
		t.Fatalf("extract circuit: %v", err)
	}
	startLine := BenchmarkPoint{X: radius, Y: 0, Z: 0}
	if distanceXZ(track.Start, startLine) > 5 {
		t.Fatalf("track start = %#v, want near lap boundary %#v", track.Start, startLine)
	}
	captureStart := BenchmarkPoint{X: samples[0].PositionX, Y: samples[0].PositionY, Z: samples[0].PositionZ}
	if distanceXZ(track.Start, captureStart) < 20 {
		t.Fatalf("track start used capture start instead of lap boundary: start=%#v capture=%#v", track.Start, captureStart)
	}
	if track.LapCountObserved != 2 {
		t.Fatalf("lap count observed = %d, want 2 complete boundary-to-boundary laps", track.LapCountObserved)
	}
	wantLapLength := 2 * math.Pi * radius
	if track.RouteLengthMeters < wantLapLength*0.85 || track.RouteLengthMeters > wantLapLength*1.15 {
		t.Fatalf("route length = %.1f, want around %.1f", track.RouteLengthMeters, wantLapLength)
	}
}

func TestFindSimilarBenchmarkTracksAndMergePreservesRuns(t *testing.T) {
	store := openTestStore(t)
	track, err := store.CreateBenchmarkTrack(BenchmarkTrackInput{
		Name:           "Original Route",
		SourceMode:     telemetry.GameModeRace,
		TrackType:      benchmarkTrackTypeSprint,
		StartRadius:    10,
		EndRadius:      10,
		Polyline:       []BenchmarkPoint{{X: 0, Z: 0}, {X: 100, Z: 0}, {X: 200, Z: 0}},
		HasDrivingLine: true,
	})
	if err != nil {
		t.Fatalf("create track: %v", err)
	}
	session := createTrackProfileSession(t, store, "Auto Baseline", 44, "A", 701, "AWD")
	insertTrackProfileRun(t, store, track.ID, session.ID, 60000, driverModeAuto, 0.9, true)

	input := BenchmarkTrackInput{
		Name:           "Updated Route",
		SourceMode:     telemetry.GameModeRace,
		TrackType:      benchmarkTrackTypeSprint,
		StartRadius:    10,
		EndRadius:      10,
		Polyline:       []BenchmarkPoint{{X: 2, Z: 1}, {X: 102, Z: 1}, {X: 202, Z: 1}},
		HasDrivingLine: true,
	}
	candidates, err := store.FindSimilarBenchmarkTracks(input)
	if err != nil {
		t.Fatalf("find similar: %v", err)
	}
	if len(candidates) != 1 || candidates[0].Track.ID != track.ID {
		t.Fatalf("candidates = %#v, want original route", candidates)
	}
	updated, err := store.MergeBenchmarkTrackInput(track.ID, input)
	if err != nil {
		t.Fatalf("merge track: %v", err)
	}
	if updated.Name != "Updated Route" || updated.RouteLengthMeters <= 0 {
		t.Fatalf("updated track = %#v", updated)
	}
	runs, err := store.ListBenchmarkRuns(track.ID, 10)
	if err != nil {
		t.Fatalf("list runs: %v", err)
	}
	if len(runs) != 1 || runs[0].SessionID != session.ID {
		t.Fatalf("runs after merge = %#v, want preserved baseline run", runs)
	}
}

func TestFindSimilarBenchmarkTracksRejectsDifferentRoute(t *testing.T) {
	store := openTestStore(t)
	if _, err := store.CreateBenchmarkTrack(BenchmarkTrackInput{
		Name:           "North Route",
		SourceMode:     telemetry.GameModeRace,
		TrackType:      benchmarkTrackTypeSprint,
		StartRadius:    10,
		EndRadius:      10,
		Polyline:       []BenchmarkPoint{{X: 0, Z: 0}, {X: 100, Z: 0}, {X: 200, Z: 0}},
		HasDrivingLine: true,
	}); err != nil {
		t.Fatalf("create track: %v", err)
	}
	candidates, err := store.FindSimilarBenchmarkTracks(BenchmarkTrackInput{
		Name:           "Different Route",
		SourceMode:     telemetry.GameModeRace,
		TrackType:      benchmarkTrackTypeSprint,
		StartRadius:    10,
		EndRadius:      10,
		Polyline:       []BenchmarkPoint{{X: 1000, Z: 1000}, {X: 1200, Z: 1200}, {X: 1400, Z: 1000}},
		HasDrivingLine: true,
	})
	if err != nil {
		t.Fatalf("find similar: %v", err)
	}
	if len(candidates) != 0 {
		t.Fatalf("candidates = %#v, want no match for different route", candidates)
	}
}

func TestFindSimilarBenchmarkTracksRejectsReverseSprint(t *testing.T) {
	store := openTestStore(t)
	if _, err := store.CreateBenchmarkTrack(BenchmarkTrackInput{
		Name:           "Forward Sprint",
		SourceMode:     telemetry.GameModeRace,
		TrackType:      benchmarkTrackTypeSprint,
		Polyline:       []BenchmarkPoint{{X: 0, Z: 0}, {X: 100, Z: 0}, {X: 200, Z: 0}},
		HasDrivingLine: true,
	}); err != nil {
		t.Fatalf("create track: %v", err)
	}
	candidates, err := store.FindSimilarBenchmarkTracks(BenchmarkTrackInput{
		Name:           "Reverse Sprint",
		SourceMode:     telemetry.GameModeRace,
		TrackType:      benchmarkTrackTypeSprint,
		Polyline:       []BenchmarkPoint{{X: 200, Z: 0}, {X: 100, Z: 0}, {X: 0, Z: 0}},
		HasDrivingLine: true,
	})
	if err != nil {
		t.Fatalf("find similar: %v", err)
	}
	if len(candidates) != 0 {
		t.Fatalf("candidates = %#v, want reverse route rejected", candidates)
	}
}

func TestFindSimilarBenchmarkTracksMatchesCircuitOffset(t *testing.T) {
	store := openTestStore(t)
	track, err := store.CreateBenchmarkTrack(BenchmarkTrackInput{
		Name:           "Circuit",
		SourceMode:     telemetry.GameModeRace,
		TrackType:      benchmarkTrackTypeCircuit,
		Polyline:       circleBenchmarkPoints(100, 96, 0),
		HasDrivingLine: true,
	})
	if err != nil {
		t.Fatalf("create circuit: %v", err)
	}
	candidates, err := store.FindSimilarBenchmarkTracks(BenchmarkTrackInput{
		Name:           "Offset Circuit",
		SourceMode:     telemetry.GameModeRace,
		TrackType:      benchmarkTrackTypeCircuit,
		Polyline:       circleBenchmarkPoints(100, 96, 24),
		HasDrivingLine: true,
	})
	if err != nil {
		t.Fatalf("find similar: %v", err)
	}
	if len(candidates) != 1 || candidates[0].Track.ID != track.ID || candidates[0].MatchLevel == "" {
		t.Fatalf("candidates = %#v, want offset circuit match", candidates)
	}
}

func TestFindSimilarBenchmarkTracksRejectsReverseCircuit(t *testing.T) {
	store := openTestStore(t)
	if _, err := store.CreateBenchmarkTrack(BenchmarkTrackInput{
		Name:           "Forward Circuit",
		SourceMode:     telemetry.GameModeRace,
		TrackType:      benchmarkTrackTypeCircuit,
		Polyline:       circleBenchmarkPoints(100, 96, 0),
		HasDrivingLine: true,
	}); err != nil {
		t.Fatalf("create circuit: %v", err)
	}
	reverse := circleBenchmarkPoints(100, 96, 0)
	for i, j := 0, len(reverse)-1; i < j; i, j = i+1, j-1 {
		reverse[i], reverse[j] = reverse[j], reverse[i]
	}
	candidates, err := store.FindSimilarBenchmarkTracks(BenchmarkTrackInput{
		Name:           "Reverse Circuit",
		SourceMode:     telemetry.GameModeRace,
		TrackType:      benchmarkTrackTypeCircuit,
		Polyline:       reverse,
		HasDrivingLine: true,
	})
	if err != nil {
		t.Fatalf("find similar: %v", err)
	}
	if len(candidates) != 0 {
		t.Fatalf("candidates = %#v, want reverse circuit rejected", candidates)
	}
}

func TestTrackBaselineCaptureDoesNotCreateSession(t *testing.T) {
	store := openTestStore(t)
	track, err := store.CreateBenchmarkTrack(BenchmarkTrackInput{
		Name:           "Baseline Sprint",
		SourceMode:     telemetry.GameModeRace,
		TrackType:      benchmarkTrackTypeSprint,
		Polyline:       []BenchmarkPoint{{X: 0, Z: 0}, {X: 100, Z: 0}},
		HasDrivingLine: true,
	})
	if err != nil {
		t.Fatalf("create track: %v", err)
	}
	samples := straightRouteSamples(0, 100, 21)
	for i := range samples {
		samples[i].GameMode = telemetry.GameModeRace
		samples[i].IsRaceOn = true
		samples[i].RacePosition = 1
		samples[i].CarPI = 700
		samples[i].Drivetrain = "AWD"
	}
	run, err := store.SaveTrackBaselineCapture(track.ID, samples, nil)
	if err != nil {
		t.Fatalf("save baseline: %v", err)
	}
	if run.ID == 0 || run.TrackID != track.ID || !run.Valid || run.Vehicle.CarOrdinal == nil || *run.Vehicle.CarOrdinal != 9001 {
		t.Fatalf("baseline run = %#v", run)
	}
	sessions, err := store.ListTelemetrySessions(10)
	if err != nil {
		t.Fatalf("list sessions: %v", err)
	}
	if len(sessions) != 0 {
		t.Fatalf("sessions = %#v, want no session from baseline capture", sessions)
	}
	profile, err := store.GetTrackProfile(track.ID)
	if err != nil {
		t.Fatalf("get profile: %v", err)
	}
	if len(profile.VehicleReferences) != 1 || profile.VehicleReferences[0].BaselineRunCount != 1 {
		t.Fatalf("vehicle references = %#v", profile.VehicleReferences)
	}
}

func TestTrackBaselineCaptureAutoMatchesStrongExistingTrack(t *testing.T) {
	store := openTestStore(t)
	track, err := store.CreateBenchmarkTrack(BenchmarkTrackInput{
		Name:           "Existing Sprint",
		SourceMode:     telemetry.GameModeRace,
		TrackType:      benchmarkTrackTypeSprint,
		Polyline:       []BenchmarkPoint{{X: 0, Z: 0}, {X: 100, Z: 0}},
		HasDrivingLine: true,
	})
	if err != nil {
		t.Fatalf("create track: %v", err)
	}
	samples := straightRouteSamples(0, 100, 21)
	for i := range samples {
		samples[i].GameMode = telemetry.GameModeRace
		samples[i].IsRaceOn = true
		samples[i].CarPI = 700
		samples[i].Drivetrain = "AWD"
	}
	result, err := store.SaveTrackBaselineCaptureAuto(0, "Auto Baseline", benchmarkTrackTypeAuto, samples, nil)
	if err != nil {
		t.Fatalf("save baseline auto: %v", err)
	}
	if result.Action != trackBaselineSaveMatchedExisting || result.Track.ID != track.ID || result.Baseline.TrackID != track.ID {
		t.Fatalf("result = %#v, want matched existing track", result)
	}
	tracks, err := store.ListBenchmarkTracks()
	if err != nil {
		t.Fatalf("list tracks: %v", err)
	}
	if len(tracks) != 1 {
		t.Fatalf("tracks = %#v, want no duplicate track", tracks)
	}
}

func TestTrackBaselineCaptureAutoCreatesTrackWhenNoStrongMatch(t *testing.T) {
	store := openTestStore(t)
	if _, err := store.CreateBenchmarkTrack(BenchmarkTrackInput{
		Name:           "North Sprint",
		SourceMode:     telemetry.GameModeRace,
		TrackType:      benchmarkTrackTypeSprint,
		Polyline:       []BenchmarkPoint{{X: 0, Z: 0}, {X: 100, Z: 0}},
		HasDrivingLine: true,
	}); err != nil {
		t.Fatalf("create track: %v", err)
	}
	samples := straightRouteSamples(1000, 1120, 25)
	for i := range samples {
		samples[i].GameMode = telemetry.GameModeRace
		samples[i].IsRaceOn = true
		samples[i].CarPI = 701
		samples[i].Drivetrain = "RWD"
	}
	result, err := store.SaveTrackBaselineCaptureAuto(0, "South Sprint", benchmarkTrackTypeAuto, samples, nil)
	if err != nil {
		t.Fatalf("save baseline auto: %v", err)
	}
	if result.Action != trackBaselineSaveCreatedTrack || result.Track.ID == 0 || result.Baseline.TrackID != result.Track.ID {
		t.Fatalf("result = %#v, want created track", result)
	}
	tracks, err := store.ListBenchmarkTracks()
	if err != nil {
		t.Fatalf("list tracks: %v", err)
	}
	if len(tracks) != 2 {
		t.Fatalf("tracks = %#v, want original plus created track", tracks)
	}
	sessions, err := store.ListTelemetrySessions(10)
	if err != nil {
		t.Fatalf("list sessions: %v", err)
	}
	if len(sessions) != 0 {
		t.Fatalf("sessions = %#v, want no session from auto baseline capture", sessions)
	}
}

func TestTrackBaselineCaptureAutoDoesNotMergeMediumCandidate(t *testing.T) {
	store := openTestStore(t)
	if _, err := store.CreateBenchmarkTrack(BenchmarkTrackInput{
		Name:           "Medium Existing",
		SourceMode:     telemetry.GameModeRace,
		TrackType:      benchmarkTrackTypeSprint,
		Polyline:       []BenchmarkPoint{{X: 0, Z: 0}, {X: 1000, Z: 0}},
		HasDrivingLine: true,
	}); err != nil {
		t.Fatalf("create track: %v", err)
	}
	samples := straightRouteSamples(0, 1070, 80)
	for i := range samples {
		samples[i].GameMode = telemetry.GameModeRace
		samples[i].IsRaceOn = true
		samples[i].CarPI = 700
		samples[i].Drivetrain = "AWD"
	}
	result, err := store.SaveTrackBaselineCaptureAuto(0, "Medium New", benchmarkTrackTypeSprint, samples, nil)
	if err != nil {
		t.Fatalf("save baseline auto: %v", err)
	}
	if result.Action != trackBaselineSaveCreatedTrack {
		t.Fatalf("action = %q, want created track for medium match", result.Action)
	}
	tracks, err := store.ListBenchmarkTracks()
	if err != nil {
		t.Fatalf("list tracks: %v", err)
	}
	if len(tracks) != 2 {
		t.Fatalf("tracks = %#v, want medium candidate preserved plus new track", tracks)
	}
}

func TestTrackBaselineCaptureAutoSelectedCircuitRejectsPartialLap(t *testing.T) {
	store := openTestStore(t)
	track, err := store.CreateBenchmarkTrack(BenchmarkTrackInput{
		Name:           "Highway Circuit",
		SourceMode:     telemetry.GameModeRace,
		TrackType:      benchmarkTrackTypeCircuit,
		Polyline:       circleBenchmarkPoints(100, 96, 0),
		HasDrivingLine: true,
	})
	if err != nil {
		t.Fatalf("create circuit: %v", err)
	}
	partial := circuitRouteSamplesWithLapNumbers(100, 1, 96, 72)
	partial = partial[:25]
	for i := range partial {
		partial[i].GameMode = telemetry.GameModeRace
		partial[i].IsRaceOn = true
		partial[i].CarPI = 800
		partial[i].Drivetrain = "AWD"
	}
	if _, err := store.SaveTrackBaselineCaptureAuto(track.ID, "Highway Circuit", benchmarkTrackTypeCircuit, partial, nil); err == nil {
		t.Fatalf("save partial selected circuit baseline succeeded, want error")
	}
	tracks, err := store.ListBenchmarkTracks()
	if err != nil {
		t.Fatalf("list tracks: %v", err)
	}
	if len(tracks) != 1 {
		t.Fatalf("tracks = %#v, want no partial duplicate circuit/sprint", tracks)
	}
}

func TestEvaluateRoadSessionMatchesAutoBaselineAndGoodFit(t *testing.T) {
	store := openTestStore(t)
	auto := createRoadEvalSession(t, store, "Auto", "auto", straightRouteSamples(0, 100, 11), nil)
	if _, err := store.CreateBenchmarkTrackFromSession(auto.ID, "Road Sprint"); err != nil {
		t.Fatalf("create benchmark track: %v", err)
	}
	if _, err := store.AnalyzeSessionBenchmarkRuns(auto.ID); err != nil {
		t.Fatalf("analyze auto baseline: %v", err)
	}
	player := createRoadEvalSession(t, store, "Player", "player", straightRouteSamples(0, 100, 8), nil)

	evaluation, err := store.EvaluateRoadSession(player.ID)
	if err != nil {
		t.Fatalf("evaluate road session: %v", err)
	}
	if evaluation.BaselineStatus != roadBaselineMatched || evaluation.BaselineRun == nil || evaluation.BaselineSession == nil {
		t.Fatalf("baseline = %q run=%#v session=%#v", evaluation.BaselineStatus, evaluation.BaselineRun, evaluation.BaselineSession)
	}
	if evaluation.OverallVerdict != roadVerdictGoodFit {
		t.Fatalf("verdict = %q, want good fit; eval=%#v", evaluation.OverallVerdict, evaluation)
	}
	if evaluation.PlayerFitScore <= 70 || evaluation.RiskScore > 35 {
		t.Fatalf("scores = fit %.1f risk %.1f", evaluation.PlayerFitScore, evaluation.RiskScore)
	}
}

func TestEvaluateRoadSessionFastButRisky(t *testing.T) {
	store := openTestStore(t)
	auto := createRoadEvalSession(t, store, "Auto", "auto", straightRouteSamples(0, 100, 11), nil)
	if _, err := store.CreateBenchmarkTrackFromSession(auto.ID, "Road Sprint"); err != nil {
		t.Fatalf("create benchmark track: %v", err)
	}
	events := []telemetry.DetectedEvent{
		{Type: "corner_exit_oversteer", Severity: "high", Segment: "corner_exit", StartMS: 1000, EndMS: 2400, DurationMS: 1400},
		{Type: "high_speed_four_wheel_slide", Severity: "high", Segment: "high_speed_corner", StartMS: 3000, EndMS: 4300, DurationMS: 1300},
	}
	player := createRoadEvalSession(t, store, "Risky", "player", straightRouteSamples(0, 100, 8), events)

	evaluation, err := store.EvaluateRoadSession(player.ID)
	if err != nil {
		t.Fatalf("evaluate risky session: %v", err)
	}
	if evaluation.OverallVerdict != roadVerdictFastButRisky {
		t.Fatalf("verdict = %q, risk %.1f", evaluation.OverallVerdict, evaluation.RiskScore)
	}
	if len(evaluation.Attributions) == 0 || evaluation.Attributions[0].Type != roadAttributionStyleFitIssue {
		t.Fatalf("attributions = %#v", evaluation.Attributions)
	}
}

func TestEvaluateRoadSessionPaperFastNotFit(t *testing.T) {
	store := openTestStore(t)
	auto := createRoadEvalSession(t, store, "Auto Fast", "auto", straightRouteSamples(0, 100, 8), nil)
	if _, err := store.CreateBenchmarkTrackFromSession(auto.ID, "Road Sprint"); err != nil {
		t.Fatalf("create benchmark track: %v", err)
	}
	events := []telemetry.DetectedEvent{
		{Type: "corner_entry_understeer", Severity: "high", Segment: "corner_entry", StartMS: 2000, EndMS: 3400, DurationMS: 1400},
		{Type: "front_brake_lockup", Severity: "medium", Segment: "braking", StartMS: 4200, EndMS: 5400, DurationMS: 1200},
	}
	player := createRoadEvalSession(t, store, "Slow Player", "player", straightRouteSamples(0, 100, 12), events)

	evaluation, err := store.EvaluateRoadSession(player.ID)
	if err != nil {
		t.Fatalf("evaluate slow player: %v", err)
	}
	if evaluation.OverallVerdict != roadVerdictPaperFastNotFit {
		t.Fatalf("verdict = %q, want paper fast not fit; risk %.1f", evaluation.OverallVerdict, evaluation.RiskScore)
	}
}

func TestEvaluateRoadSessionMissingAutoBaseline(t *testing.T) {
	store := openTestStore(t)
	player := createRoadEvalSession(t, store, "No Baseline", "player", straightRouteSamples(0, 100, 9), nil)
	if _, err := store.CreateBenchmarkTrackFromSession(player.ID, "Road Sprint"); err != nil {
		t.Fatalf("create benchmark track: %v", err)
	}
	evaluation, err := store.EvaluateRoadSession(player.ID)
	if err != nil {
		t.Fatalf("evaluate missing baseline: %v", err)
	}
	if evaluation.BaselineStatus != roadBaselineMissingAuto || evaluation.OverallVerdict != roadVerdictInsufficientData {
		t.Fatalf("eval = %#v", evaluation)
	}
}

func TestCarIdentityBinding(t *testing.T) {
	store := openTestStore(t)
	ordinal := int64(42)
	if _, err := store.CreateTuneProfile(TuneProfileInput{CarName: "Known Car", CarOrdinal: &ordinal}); err != nil {
		t.Fatalf("create profile: %v", err)
	}
	name, err := store.ResolveCarNameByOrdinal(ordinal)
	if err != nil {
		t.Fatalf("resolve car name: %v", err)
	}
	if name != "Known Car" {
		t.Fatalf("resolved name = %q", name)
	}
}

func TestRuleThresholdProfileCRUDResetAndMatch(t *testing.T) {
	store := openTestStore(t)
	profiles, err := store.ListRuleThresholdProfiles()
	if err != nil {
		t.Fatalf("list default rule profiles: %v", err)
	}
	var defaultProfile *RuleThresholdProfile
	var roadProfile *RuleThresholdProfile
	for i := range profiles {
		if profiles[i].IsDefault {
			defaultProfile = &profiles[i]
		}
		if profiles[i].Name == "Road Racing" && profiles[i].UseCase == "Road" {
			roadProfile = &profiles[i]
		}
	}
	if defaultProfile == nil || roadProfile == nil {
		t.Fatalf("rule profiles = %#v, want default and Road Racing", profiles)
	}
	if err := store.DeleteRuleThresholdProfile(defaultProfile.ID); err == nil {
		t.Fatal("expected default profile delete to fail")
	}
	roadConfig, err := parseRuleConfigJSON(roadProfile.ConfigJSON)
	if err != nil {
		t.Fatalf("parse road config: %v", err)
	}
	if !roadConfig.Events["long_gear_bog_down"].Enabled || !roadConfig.Events["top_speed_limited_by_gearing"].Enabled {
		t.Fatalf("road config missing road gearing events: %#v", roadConfig.Events)
	}

	drivetrain := "AWD"
	carClass := "S1"
	useCase := "Road"
	created, err := store.CreateRuleThresholdProfile(RuleThresholdProfileInput{
		Name:       "S1 AWD Road",
		CarClass:   carClass,
		Drivetrain: drivetrain,
		UseCase:    useCase,
		ConfigJSON: defaultProfile.ConfigJSON,
	})
	if err != nil {
		t.Fatalf("create rule profile: %v", err)
	}
	if created.ID == 0 || created.IsDefault {
		t.Fatalf("created rule profile = %#v", created)
	}

	profile, err := store.CreateTuneProfile(TuneProfileInput{CarName: "Match", CarClass: carClass, Drivetrain: drivetrain, UseCase: useCase})
	if err != nil {
		t.Fatalf("create tune: %v", err)
	}
	matched, config, err := store.MatchRuleThresholdProfile(profile)
	if err != nil {
		t.Fatalf("match rule profile: %v", err)
	}
	if matched == nil || matched.ID != created.ID || !config.Events["launch_wheelspin"].Enabled {
		t.Fatalf("matched=%#v config=%#v", matched, config)
	}

	created.ConfigJSON = "{}"
	reset, err := store.ResetRuleThresholdProfile(created.ID)
	if err != nil {
		t.Fatalf("reset rule profile: %v", err)
	}
	resetConfig, err := parseRuleConfigJSON(reset.ConfigJSON)
	if err != nil {
		t.Fatalf("parse reset config: %v", err)
	}
	if !resetConfig.Events["long_gear_bog_down"].Enabled || !strings.Contains(reset.ConfigJSON, "launch_wheelspin") {
		t.Fatalf("reset config = %s", reset.ConfigJSON)
	}
	if err := store.DeleteRuleThresholdProfile(created.ID); err != nil {
		t.Fatalf("delete rule profile: %v", err)
	}
	if err := store.DeleteRuleThresholdProfile(roadProfile.ID); err != nil {
		t.Fatalf("delete Road Racing profile: %v", err)
	}
	if err := store.Migrate(); err != nil {
		t.Fatalf("migrate after deleting Road Racing: %v", err)
	}
	afterDelete, err := store.ListRuleThresholdProfiles()
	if err != nil {
		t.Fatalf("list after Road Racing delete: %v", err)
	}
	for _, profile := range afterDelete {
		if profile.Name == "Road Racing" && profile.UseCase == "Road" {
			t.Fatalf("Road Racing profile was recreated after delete: %#v", afterDelete)
		}
	}
}

func TestSessionIssueSummaryMergesAndComparesBaseline(t *testing.T) {
	store := openTestStore(t)
	profile, err := store.CreateTuneProfile(TuneProfileInput{CarName: "Road Car", CarClass: "A", UseCase: "Road"})
	if err != nil {
		t.Fatalf("create profile: %v", err)
	}
	previous := createIssueSummarySession(t, store, profile.ID, "Previous", "2026-05-18T10:00:00Z", []telemetry.DetectedEvent{
		testEvent("p1", "corner_entry_understeer", "high", 0, 900),
		testEvent("p2", "corner_entry_understeer", "medium", 1200, 2200),
		testEvent("p3", "corner_entry_understeer", "medium", 2400, 3400),
	})
	current := createIssueSummarySession(t, store, profile.ID, "Current", "2026-05-18T11:00:00Z", []telemetry.DetectedEvent{
		testEvent("c1", "corner_entry_understeer", "medium", 0, 800),
		testEvent("c2", "corner_entry_understeer", "medium", 1200, 1700),
	})
	_ = previous

	summary, err := store.GetSessionIssueSummary(current.ID)
	if err != nil {
		t.Fatalf("issue summary: %v", err)
	}
	if summary.BaselineSession == nil || summary.BaselineSession.ID != previous.ID {
		t.Fatalf("baseline = %#v, want previous session", summary.BaselineSession)
	}
	if len(summary.Groups) != 1 {
		t.Fatalf("group count = %d, want 1: %#v", len(summary.Groups), summary.Groups)
	}
	group := summary.Groups[0]
	if group.Family != "corner_entry_balance" || group.EventCount != 2 {
		t.Fatalf("group = %#v, want merged corner entry group with 2 events", group)
	}
	if group.Comparison != "improved" {
		t.Fatalf("comparison = %q, want improved", group.Comparison)
	}
}

func TestSessionIssueSummaryAddsGearTelemetryAndTuneComparisons(t *testing.T) {
	store := openTestStore(t)
	gear3 := 1.30
	profile, err := store.CreateTuneProfile(TuneProfileInput{CarName: "Gear Car", CarClass: "A", UseCase: "Road", Gear3: &gear3})
	if err != nil {
		t.Fatalf("create profile: %v", err)
	}
	previousSamples := gearComparisonSamples(3, 82)
	previous := createIssueSummarySessionWithSamples(t, store, profile.ID, "Previous Gear", "2099-05-18T10:00:00Z", nil, previousSamples)
	_ = previous

	updated, err := store.GetTuneProfile(profile.ID)
	if err != nil {
		t.Fatalf("get profile: %v", err)
	}
	input := updated.ToInput()
	gear3Next := 1.38
	input.Gear3 = &gear3Next
	if _, err := store.UpdateTuneProfile(profile.ID, input); err != nil {
		t.Fatalf("update profile gear: %v", err)
	}
	currentSamples := gearComparisonSamples(3, 96)
	current := createIssueSummarySessionWithSamples(t, store, profile.ID, "Current Gear", "2099-05-18T11:00:00Z", nil, currentSamples)

	summary, err := store.GetSessionIssueSummary(current.ID)
	if err != nil {
		t.Fatalf("issue summary: %v", err)
	}
	telemetryComparison := gearComparisonByType(summary.GearPower.Comparisons, "session_telemetry")
	if telemetryComparison == nil || telemetryComparison.Status != "ready" || len(telemetryComparison.Rows) == 0 {
		t.Fatalf("telemetry comparison = %#v", telemetryComparison)
	}
	if telemetryComparison.Rows[0].Gear != 3 || telemetryComparison.Rows[0].SpeedMaxDeltaKmh <= 0 {
		t.Fatalf("telemetry row = %#v, want gear 3 speed improvement", telemetryComparison.Rows[0])
	}
	tuneComparison := gearComparisonByType(summary.GearPower.Comparisons, "tune_settings")
	if tuneComparison == nil || tuneComparison.Status != "ready" || len(tuneComparison.Rows) == 0 {
		t.Fatalf("tune comparison = %#v", tuneComparison)
	}
	if tuneComparison.Rows[0].Item != "gear_3" || tuneComparison.Rows[0].DeltaValue <= 0 {
		t.Fatalf("tune row = %#v, want gear_3 increase", tuneComparison.Rows[0])
	}
}

func TestGearPowerDiagnosticDetectsLongGear(t *testing.T) {
	samples := make([]telemetry.NormalizedTelemetry, 0, 12)
	for i := 0; i < 12; i++ {
		samples = append(samples, telemetry.NormalizedTelemetry{
			TimeMS:     int64(i) * 100,
			Gear:       3,
			SpeedKmh:   85 + float64(i),
			Throttle01: 0.9,
			RpmRatio:   0.46,
		})
	}
	diag := BuildGearPowerDiagnostic(samples, nil, nil)
	if diag.Status != "ok" || len(diag.Gears) != 1 {
		t.Fatalf("diagnostic = %#v, want one gear", diag)
	}
	if diag.Gears[0].Finding != "too_long" {
		t.Fatalf("gear finding = %q, want too_long", diag.Gears[0].Finding)
	}
	if len(diag.RecommendedActions) == 0 || diag.RecommendedActions[0].Item != "gear_3" || diag.RecommendedActions[0].Direction != "increase" {
		t.Fatalf("recommended actions = %#v, want gear_3 increase", diag.RecommendedActions)
	}
}

func TestGearPowerDiagnosticUsesProfilePowerBand(t *testing.T) {
	peakTorqueRPM := 4200.0
	peakPowerRPM := 6800.0
	redlineRPM := 7200.0
	gear3 := 1.35
	profile := &TuneProfile{PeakTorqueRPM: &peakTorqueRPM, PeakPowerRPM: &peakPowerRPM, RedlineRPM: &redlineRPM, Gear3: &gear3}
	samples := make([]telemetry.NormalizedTelemetry, 0, 12)
	for i := 0; i < 12; i++ {
		samples = append(samples, telemetry.NormalizedTelemetry{
			TimeMS:       int64(i) * 100,
			Gear:         3,
			SpeedKmh:     75 + float64(i),
			Throttle01:   0.9,
			Rpm:          3100,
			RpmRatio:     0.62,
			EngineMaxRpm: 7200,
		})
	}
	diag := BuildGearPowerDiagnostic(samples, nil, profile)
	if diag.PowerBandSource != "profile_power_band" || diag.PowerBandStartRPM != peakTorqueRPM {
		t.Fatalf("power band = %#v, want profile power band", diag)
	}
	if diag.Gears[0].Finding != "too_long" {
		t.Fatalf("gear finding = %q, want too_long below profile band", diag.Gears[0].Finding)
	}
	if diag.Gears[0].InPowerBandRpmMax != 0 {
		t.Fatalf("in-band rpm range = %.0f-%.0f, want empty range below profile band", diag.Gears[0].InPowerBandRpmMin, diag.Gears[0].InPowerBandRpmMax)
	}
	if diag.RecommendedActions[0].Amount != "0.08" {
		t.Fatalf("action amount = %q, want severity-scaled 0.08", diag.RecommendedActions[0].Amount)
	}
}

func TestGearPowerDiagnosticReportsInPowerBandRange(t *testing.T) {
	peakTorqueRPM := 4200.0
	peakPowerRPM := 6800.0
	redlineRPM := 7200.0
	gear3 := 1.35
	profile := &TuneProfile{PeakTorqueRPM: &peakTorqueRPM, PeakPowerRPM: &peakPowerRPM, RedlineRPM: &redlineRPM, Gear3: &gear3}
	samples := make([]telemetry.NormalizedTelemetry, 0, 12)
	for i := 0; i < 12; i++ {
		samples = append(samples, telemetry.NormalizedTelemetry{
			TimeMS:       int64(i) * 100,
			Gear:         3,
			SpeedKmh:     75 + float64(i),
			Throttle01:   0.9,
			Rpm:          4300 + float64(i*100),
			RpmRatio:     0.62 + float64(i)*0.01,
			EngineMaxRpm: 7200,
		})
	}
	diag := BuildGearPowerDiagnostic(samples, nil, profile)
	if len(diag.Gears) != 1 {
		t.Fatalf("gears = %#v, want one gear", diag.Gears)
	}
	if diag.Gears[0].InPowerBandRpmMin != 4300 || diag.Gears[0].InPowerBandRpmMax != 5400 {
		t.Fatalf("in-band rpm range = %.0f-%.0f, want 4300-5400", diag.Gears[0].InPowerBandRpmMin, diag.Gears[0].InPowerBandRpmMax)
	}
}

func TestGearPowerDiagnosticDetectsShortGearNearRedline(t *testing.T) {
	peakTorqueRPM := 4200.0
	peakPowerRPM := 6800.0
	redlineRPM := 7200.0
	gear4 := 1.05
	profile := &TuneProfile{PeakTorqueRPM: &peakTorqueRPM, PeakPowerRPM: &peakPowerRPM, RedlineRPM: &redlineRPM, Gear4: &gear4}
	samples := make([]telemetry.NormalizedTelemetry, 0, 12)
	for i := 0; i < 12; i++ {
		samples = append(samples, telemetry.NormalizedTelemetry{
			TimeMS:       int64(i) * 100,
			Gear:         4,
			SpeedKmh:     115 + float64(i),
			Throttle01:   0.92,
			Rpm:          7100,
			RpmRatio:     0.86,
			EngineMaxRpm: 7200,
		})
	}
	diag := BuildGearPowerDiagnostic(samples, nil, profile)
	if diag.Gears[0].Finding != "too_short" {
		t.Fatalf("gear finding = %q, want too_short near redline", diag.Gears[0].Finding)
	}
	if len(diag.RecommendedActions) == 0 || diag.RecommendedActions[0].Item != "gear_4" || diag.RecommendedActions[0].Direction != "decrease" {
		t.Fatalf("recommended actions = %#v, want gear_4 decrease", diag.RecommendedActions)
	}
}

func TestGearPowerDiagnosticSkipsLockedProfileGears(t *testing.T) {
	samples := make([]telemetry.NormalizedTelemetry, 0, 12)
	for i := 0; i < 12; i++ {
		samples = append(samples, telemetry.NormalizedTelemetry{
			TimeMS:     int64(i) * 100,
			Gear:       9,
			SpeedKmh:   130 + float64(i),
			Throttle01: 0.9,
			RpmRatio:   0.46,
		})
	}
	profile := &TuneProfile{}
	diag := BuildGearPowerDiagnostic(samples, nil, profile)
	if diag.Status != "insufficient_data" || diag.Summary != "no_unlocked_gear_samples" || len(diag.Gears) != 0 || len(diag.RecommendedActions) != 0 {
		t.Fatalf("diagnostic = %#v, want locked gear skipped", diag)
	}

	gear9 := 1.0
	profile.Gear9 = &gear9
	diag = BuildGearPowerDiagnostic(samples, nil, profile)
	if diag.Status != "ok" || len(diag.Gears) != 1 || diag.Gears[0].Gear != 9 {
		t.Fatalf("diagnostic = %#v, want unlocked gear 9", diag)
	}
	if len(diag.RecommendedActions) == 0 || diag.RecommendedActions[0].Item != "gear_9" || diag.RecommendedActions[0].Direction != "increase" {
		t.Fatalf("recommended actions = %#v, want gear_9 increase", diag.RecommendedActions)
	}
}

func TestGearPowerDiagnosticExplainsInsufficientHighLoad(t *testing.T) {
	samples := make([]telemetry.NormalizedTelemetry, 0, 12)
	for i := 0; i < 12; i++ {
		samples = append(samples, telemetry.NormalizedTelemetry{
			TimeMS:     int64(i) * 100,
			Gear:       3,
			SpeedKmh:   80 + float64(i),
			Throttle01: 0.32,
			RpmRatio:   0.72,
		})
	}
	diag := BuildGearPowerDiagnostic(samples, nil, nil)
	if diag.Status != "ok" || diag.Summary != "not_enough_high_load" {
		t.Fatalf("diagnostic = %#v, want not_enough_high_load", diag)
	}
	if diag.Evidence["power_band_total_samples"] != 12 || diag.Evidence["power_band_high_load_samples"] != 0 {
		t.Fatalf("evidence = %#v, want sample and high-load counts", diag.Evidence)
	}
}

func TestGearPowerDiagnosticDetectsTractionLimitedPower(t *testing.T) {
	samples := make([]telemetry.NormalizedTelemetry, 0, 12)
	for i := 0; i < 12; i++ {
		samples = append(samples, telemetry.NormalizedTelemetry{
			TimeMS:              int64(i) * 100,
			Gear:                1,
			SpeedKmh:            35 + float64(i),
			Throttle01:          0.95,
			RpmRatio:            0.78,
			Drivetrain:          "RWD",
			RearCombinedSlipAvg: 1.35,
		})
	}
	diag := BuildGearPowerDiagnostic(samples, nil, nil)
	if diag.Gears[0].Finding != "traction_limited" {
		t.Fatalf("gear finding = %q, want traction_limited", diag.Gears[0].Finding)
	}
	if diag.Summary != "traction_limited_power" {
		t.Fatalf("summary = %q, want traction_limited_power", diag.Summary)
	}
	if len(diag.RecommendedActions) == 0 || diag.RecommendedActions[0].Item != "gear_1" || diag.RecommendedActions[0].Direction != "decrease" {
		t.Fatalf("recommended actions = %#v, want gear_1 decrease", diag.RecommendedActions)
	}
}

func TestGearPowerDiagnosticUsesDrivenAxleDiffAccel(t *testing.T) {
	samples := make([]telemetry.NormalizedTelemetry, 0, 24)
	for gear := 2; gear <= 3; gear++ {
		for i := 0; i < 12; i++ {
			samples = append(samples, telemetry.NormalizedTelemetry{
				TimeMS:              int64(len(samples)) * 100,
				Gear:                gear,
				SpeedKmh:            55 + float64(i),
				Throttle01:          0.95,
				RpmRatio:            0.78,
				Drivetrain:          "RWD",
				RearCombinedSlipAvg: 1.35,
			})
		}
	}
	diag := BuildGearPowerDiagnostic(samples, nil, nil)
	rearDiffAccelCount := 0
	for _, action := range diag.RecommendedActions {
		if action.Item == "drive_diff_accel" {
			t.Fatalf("recommended actions = %#v, want no generic drive_diff_accel action", diag.RecommendedActions)
		}
		if action.Item == "rear_diff_accel" {
			rearDiffAccelCount++
		}
	}
	if rearDiffAccelCount != 1 {
		t.Fatalf("recommended actions = %#v, want exactly one rear_diff_accel action for RWD", diag.RecommendedActions)
	}
	if diag.Evidence["power_band_target_min"] != powerBandRpmMin || diag.Evidence["power_band_target_max"] != powerBandRpmMax {
		t.Fatalf("evidence = %#v, want target power band", diag.Evidence)
	}
}

func TestGearPowerDiagnosticSkipsDiffAccelWhenDrivetrainUnknown(t *testing.T) {
	samples := make([]telemetry.NormalizedTelemetry, 0, 24)
	for gear := 2; gear <= 3; gear++ {
		for i := 0; i < 12; i++ {
			samples = append(samples, telemetry.NormalizedTelemetry{
				TimeMS:              int64(len(samples)) * 100,
				Gear:                gear,
				SpeedKmh:            55 + float64(i),
				Throttle01:          0.95,
				RpmRatio:            0.78,
				RearCombinedSlipAvg: 1.35,
			})
		}
	}
	diag := BuildGearPowerDiagnostic(samples, nil, nil)
	for _, action := range diag.RecommendedActions {
		if action.Item == "front_diff_accel" || action.Item == "rear_diff_accel" || action.Item == "drive_diff_accel" {
			t.Fatalf("recommended actions = %#v, want no diff accel action when drivetrain is unknown", diag.RecommendedActions)
		}
	}
	if diag.Summary != "traction_limited_power" {
		t.Fatalf("summary = %q, want traction_limited_power without automatic diff action", diag.Summary)
	}
}

func TestGearPowerDiagnosticAWDDiffAccelUsesPrimaryAxle(t *testing.T) {
	samples := make([]telemetry.NormalizedTelemetry, 0, 24)
	for gear := 2; gear <= 3; gear++ {
		for i := 0; i < 12; i++ {
			samples = append(samples, telemetry.NormalizedTelemetry{
				TimeMS:               int64(len(samples)) * 100,
				Gear:                 gear,
				SpeedKmh:             55 + float64(i),
				Throttle01:           0.95,
				RpmRatio:             0.78,
				Drivetrain:           "AWD",
				FrontCombinedSlipAvg: 0.55,
				RearCombinedSlipAvg:  1.35,
			})
		}
	}
	diag := BuildGearPowerDiagnostic(samples, nil, nil)
	foundRear := false
	for _, action := range diag.RecommendedActions {
		if action.Item == "rear_diff_accel" && action.Amount == "2" {
			foundRear = true
		}
		if action.Item == "front_diff_accel" {
			t.Fatalf("recommended actions = %#v, want rear-only primary action", diag.RecommendedActions)
		}
	}
	if !foundRear {
		t.Fatalf("recommended actions = %#v, want rear_diff_accel primary action", diag.RecommendedActions)
	}
}

func TestGearPowerDiagnosticAWDDiffAccelUsesPrimaryAndSecondarySteps(t *testing.T) {
	samples := make([]telemetry.NormalizedTelemetry, 0, 24)
	for gear := 2; gear <= 3; gear++ {
		for i := 0; i < 12; i++ {
			samples = append(samples, telemetry.NormalizedTelemetry{
				TimeMS:               int64(len(samples)) * 100,
				Gear:                 gear,
				SpeedKmh:             55 + float64(i),
				Throttle01:           0.95,
				RpmRatio:             0.78,
				Drivetrain:           "AWD",
				FrontCombinedSlipAvg: 1.30,
				RearCombinedSlipAvg:  1.35,
			})
		}
	}
	diag := BuildGearPowerDiagnostic(samples, nil, nil)
	amounts := map[string]string{}
	for _, action := range diag.RecommendedActions {
		if action.Item == "front_diff_accel" || action.Item == "rear_diff_accel" {
			amounts[action.Item] = action.Amount
		}
	}
	if amounts["rear_diff_accel"] != "2" || amounts["front_diff_accel"] != "1" {
		t.Fatalf("recommended actions = %#v, want rear primary 2 and front secondary 1", diag.RecommendedActions)
	}
}

func TestGearPowerDiagnosticKeepsDistinctGearActions(t *testing.T) {
	samples := make([]telemetry.NormalizedTelemetry, 0, 24)
	for gear := 2; gear <= 3; gear++ {
		for i := 0; i < 12; i++ {
			samples = append(samples, telemetry.NormalizedTelemetry{
				TimeMS:     int64(len(samples)) * 100,
				Gear:       gear,
				SpeedKmh:   70 + float64(i),
				Throttle01: 0.9,
				RpmRatio:   0.46,
			})
		}
	}
	diag := BuildGearPowerDiagnostic(samples, nil, nil)
	found := map[string]bool{}
	for _, action := range diag.RecommendedActions {
		found[action.Item+"_"+action.Direction] = true
	}
	if !found["gear_2_increase"] || !found["gear_3_increase"] {
		t.Fatalf("recommended actions = %#v, want distinct gear_2 and gear_3 increase actions", diag.RecommendedActions)
	}
}

func TestGearPowerDiagnosticGlobalLongGearsUsesFinalDrive(t *testing.T) {
	samples := gearPowerFindingSamples([]int{2, 3, 4}, 0.46, 70, 0)
	diag := BuildGearPowerDiagnostic(samples, nil, nil)
	if diag.StrategyMode != "global_too_long" || diag.GlobalGearIssueCount != 3 || diag.UsableGearCount != 3 {
		t.Fatalf("diagnostic = %#v, want global too-long strategy", diag)
	}
	if !hasSuggestedAction(diag.RecommendedActions, "final_drive", "increase") {
		t.Fatalf("recommended actions = %#v, want final drive increase", diag.RecommendedActions)
	}
	if hasSuggestedAction(diag.RecommendedActions, "gear_2", "increase") || hasSuggestedAction(diag.RecommendedActions, "gear_3", "increase") || hasSuggestedAction(diag.RecommendedActions, "gear_4", "increase") {
		t.Fatalf("recommended actions = %#v, want no individual gear increases for global issue", diag.RecommendedActions)
	}
}

func TestGearPowerDiagnosticGlobalShortGearsUsesFinalDrive(t *testing.T) {
	samples := gearPowerFindingSamples([]int{2, 3, 4}, 0.96, 70, 0)
	diag := BuildGearPowerDiagnostic(samples, nil, nil)
	if diag.StrategyMode != "global_too_short" || diag.GlobalGearIssueCount != 3 || diag.UsableGearCount != 3 {
		t.Fatalf("diagnostic = %#v, want global too-short strategy", diag)
	}
	if !hasSuggestedAction(diag.RecommendedActions, "final_drive", "decrease") {
		t.Fatalf("recommended actions = %#v, want final drive decrease", diag.RecommendedActions)
	}
	if hasSuggestedAction(diag.RecommendedActions, "gear_2", "decrease") || hasSuggestedAction(diag.RecommendedActions, "gear_3", "decrease") || hasSuggestedAction(diag.RecommendedActions, "gear_4", "decrease") {
		t.Fatalf("recommended actions = %#v, want no individual gear decreases for global issue", diag.RecommendedActions)
	}
}

func TestGearPowerDiagnosticTractionLimitedDoesNotUseGlobalFinalDrive(t *testing.T) {
	samples := gearPowerFindingSamples([]int{1, 2, 3}, 0.78, 35, 1.35)
	for i := range samples {
		samples[i].Drivetrain = "RWD"
	}
	diag := BuildGearPowerDiagnostic(samples, nil, nil)
	if diag.StrategyMode != "traction_limited_low_gears" {
		t.Fatalf("strategy = %q, want traction_limited_low_gears", diag.StrategyMode)
	}
	if hasSuggestedAction(diag.RecommendedActions, "final_drive", "increase") || hasSuggestedAction(diag.RecommendedActions, "final_drive", "decrease") {
		t.Fatalf("recommended actions = %#v, want no final drive action when traction-limited", diag.RecommendedActions)
	}
}

func TestGearPowerDiagnosticTopSpeedProtectionBeatsGlobalLongGears(t *testing.T) {
	samples := gearPowerFindingSamples([]int{2, 3, 4}, 0.46, 75, 0)
	for i := 0; i < 12; i++ {
		samples = append(samples, telemetry.NormalizedTelemetry{
			TimeMS:     int64(len(samples)) * 100,
			Gear:       5,
			SpeedKmh:   180 + float64(i),
			Throttle01: 0.9,
			RpmRatio:   0.98,
		})
	}
	diag := BuildGearPowerDiagnostic(samples, nil, nil)
	if diag.StrategyMode != "top_speed_limited" || diag.TopSpeedFinding != "top_speed_limited_by_gearing" {
		t.Fatalf("diagnostic = %#v, want top-speed strategy", diag)
	}
	if !hasSuggestedAction(diag.RecommendedActions, "final_drive", "decrease") {
		t.Fatalf("recommended actions = %#v, want final drive decrease for top-speed headroom", diag.RecommendedActions)
	}
	if hasSuggestedAction(diag.RecommendedActions, "final_drive", "increase") {
		t.Fatalf("recommended actions = %#v, want no final drive increase near redline", diag.RecommendedActions)
	}
}

func TestDetectDriverModeAutoPlayerAndFreeRoam(t *testing.T) {
	autoSamples := make([]telemetry.NormalizedTelemetry, 0, 80)
	for i := 0; i < 80; i++ {
		autoSamples = append(autoSamples, telemetry.NormalizedTelemetry{
			TimeMS:        int64(i) * 100,
			IsRaceOn:      true,
			GameMode:      telemetry.GameModeRace,
			SpeedKmh:      90,
			Rpm:           5000,
			CarOrdinal:    42,
			DrivingLine01: 0.4,
			Throttle01:    0.7,
			Steer01:       0.1,
		})
	}
	auto := DetectDriverMode(autoSamples, nil, telemetry.GameModeRace)
	if auto.Mode != "auto" || auto.Confidence < 0.7 {
		t.Fatalf("auto detection = %#v, want confident auto", auto)
	}

	playerSamples := append([]telemetry.NormalizedTelemetry(nil), autoSamples...)
	for i := range playerSamples {
		playerSamples[i].RearCombinedSlipAvg = 1.4
		playerSamples[i].Steer01 = math.Sin(float64(i)) * 0.5
	}
	player := DetectDriverMode(playerSamples, []telemetry.DetectedEvent{{Type: "corner_exit_oversteer", Severity: "high"}}, telemetry.GameModeRace)
	if player.Mode != "player" || player.Confidence < 0.7 {
		t.Fatalf("player detection = %#v, want confident player", player)
	}

	freeRoamSamples := append([]telemetry.NormalizedTelemetry(nil), autoSamples...)
	for i := range freeRoamSamples {
		freeRoamSamples[i].IsRaceOn = false
		freeRoamSamples[i].GameMode = telemetry.GameModeFreeRoam
	}
	freeRoam := DetectDriverMode(freeRoamSamples, nil, telemetry.GameModeFreeRoam)
	if freeRoam.Mode != "player" {
		t.Fatalf("free roam detection = %#v, want player", freeRoam)
	}
}

func TestWholeCarPlanResolvesTirePressureConflict(t *testing.T) {
	groups := []SessionIssueGroup{
		{
			Family:          "tire_temperature_stability",
			Severity:        "high",
			EventCount:      4,
			TotalDurationMS: 6000,
			PrimaryActions: []telemetry.SuggestedAction{
				{Priority: 0, Category: "tire", Item: "tire_pressure", Direction: "increase", Amount: "0.03 BAR (鈮?.5 PSI)", Reason: "stabilize tire temperature and contact patch"},
			},
			Evidence: map[string]IssueEvidence{"front_tire_temp": {Avg: 180, Count: 4}},
		},
		{
			Family:          "launch_traction",
			Severity:        "medium",
			EventCount:      1,
			TotalDurationMS: 800,
			PrimaryActions: []telemetry.SuggestedAction{
				{Priority: 0, Category: "tire", Item: "drive_tire_pressure", Direction: "decrease", Amount: "0.03 BAR (鈮?.5 PSI)", Reason: "increase launch traction"},
			},
			Evidence: map[string]IssueEvidence{"rear_slip_ratio": {Avg: 1.4, Count: 1}},
		},
	}
	plan := BuildWholeCarTuningPlan(groups, GearPowerDiagnostic{}, nil)
	if len(plan.Actions) == 0 {
		t.Fatalf("plan has no actions: %#v", plan)
	}
	if plan.Actions[0].Item != "tire_pressure" || plan.Actions[0].Direction != "increase" {
		t.Fatalf("first action = %#v, want tire_pressure increase", plan.Actions[0])
	}
	if len(plan.Conflicts) == 0 {
		t.Fatalf("expected conflict to be recorded: %#v", plan)
	}
}

func TestCornerEntryStrategyIncludesRearArbAlternative(t *testing.T) {
	groups := BuildSessionIssueGroups([]telemetry.DetectedEvent{
		testEvent("entry-1", "corner_entry_understeer", "high", 0, 1200),
	}, nil)
	applyIssueStrategies(groups, nil)
	if len(groups) != 1 {
		t.Fatalf("group count = %d, want 1", len(groups))
	}
	items := map[string]bool{}
	for _, action := range groups[0].PrimaryActions {
		items[action.Item+"/"+action.Direction] = true
	}
	if !items["front_arb/decrease"] {
		t.Fatalf("actions = %#v, want front_arb decrease", groups[0].PrimaryActions)
	}
	if !items["rear_arb/increase"] {
		t.Fatalf("actions = %#v, want rear_arb increase as rotation alternative", groups[0].PrimaryActions)
	}
}

func TestTunePlanDraftBlocksManualReviewActions(t *testing.T) {
	camber := -1.2
	profile := &TuneProfile{FrontCamber: &camber}
	actions := draftActionsForWholeCarAdjustment(WholeCarAdjustment{
		Category:  "alignment",
		Item:      "front_camber",
		Direction: "check",
		Amount:    "slightly more negative",
		Reason:    "improve front tire contact in cornering",
	}, profile, "ready")
	if len(actions) != 1 {
		t.Fatalf("actions = %#v, want one action", actions)
	}
	if actions[0].CanApply || actions[0].BlockedReason != "manual_review_required" {
		t.Fatalf("action = %#v, want manual review block", actions[0])
	}
}

func TestRetestStatusWeightsPerformanceAgainstIssueNoise(t *testing.T) {
	status := retestStatusFromMetrics([]RetestMetric{
		{Key: "issue_score", Status: "worsened"},
		{Key: "event_count", Status: "worsened"},
		{Key: "event_duration_ms", Status: "worsened"},
		{Key: "avg_speed_kmh", Status: "improved"},
		{Key: "best_run_duration_ms", Status: "improved"},
	})
	if status != "improved" {
		t.Fatalf("status = %q, want improved when speed and best segment improved strongly", status)
	}
}

func TestQuickDiagnosticComparesCurrentAndPreviousLap(t *testing.T) {
	samples := quickLapSamples(1, 0, 10, 105)
	samples = append(samples, quickLapSamples(2, 10000, 10, 112)...)
	events := []telemetry.DetectedEvent{
		testEvent("quick-prev", "corner_entry_understeer", "high", 1200, 2200),
		testEvent("quick-current", "corner_entry_understeer", "medium", 11200, 11600),
	}
	diag := BuildQuickDiagnostic(samples, events, &samples[len(samples)-1], nil)
	if diag.Status != "ready" || diag.ComparisonStatus != "lap_comparison" {
		t.Fatalf("diagnostic = %#v, want lap comparison", diag)
	}
	if diag.Comparability.Confidence != "high" || diag.Comparability.SameVehicleClass != "yes" || diag.Comparability.SameTrackContext != "yes" {
		t.Fatalf("comparability = %#v, want high same vehicle and track", diag.Comparability)
	}
	if diag.CurrentLap == nil || diag.PreviousLap == nil || diag.CurrentLap.LapNumber != 2 || diag.PreviousLap.LapNumber != 1 {
		t.Fatalf("laps = current %#v previous %#v", diag.CurrentLap, diag.PreviousLap)
	}
	if len(diag.Groups) == 0 || diag.Groups[0].Comparison != issueCompareImproved {
		t.Fatalf("groups = %#v, want current lap improved vs previous", diag.Groups)
	}
}

func TestQuickDiagnosticRejectsChangedVehicleForLapComparison(t *testing.T) {
	samples := quickLapSamples(1, 0, 10, 105)
	changed := quickLapSamples(2, 10000, 10, 112)
	for i := range changed {
		changed[i].CarOrdinal = 9900
	}
	samples = append(samples, changed...)
	diag := BuildQuickDiagnostic(samples, nil, &samples[len(samples)-1], nil)
	if diag.ComparisonStatus != "rolling_window_only" || diag.Comparability.Confidence != "invalid" {
		t.Fatalf("diagnostic = %#v, want invalid rolling window", diag)
	}
	if diag.Comparability.SameVehicleClass != "no" || !containsString(diag.Comparability.Warnings, "quick_vehicle_or_class_changed") {
		t.Fatalf("comparability = %#v, want vehicle changed warning", diag.Comparability)
	}
}

func TestQuickDiagnosticRejectsChangedClassForLapComparison(t *testing.T) {
	samples := quickLapSamples(1, 0, 10, 105)
	changed := quickLapSamples(2, 10000, 10, 112)
	for i := range changed {
		changed[i].CarClass = "S1"
	}
	samples = append(samples, changed...)
	diag := BuildQuickDiagnostic(samples, nil, &samples[len(samples)-1], nil)
	if diag.ComparisonStatus != "rolling_window_only" || diag.Comparability.Confidence != "invalid" || diag.Comparability.SameVehicleClass != "no" {
		t.Fatalf("comparability = %#v, want invalid class mismatch", diag.Comparability)
	}
}

func TestQuickDiagnosticRejectsRaceTimeResetForLapComparison(t *testing.T) {
	samples := quickLapSamples(1, 0, 10, 105)
	reset := quickLapSamples(2, 10000, 10, 112)
	for i := range reset {
		reset[i].CurrentRaceTime = float64(i + 1)
	}
	samples = append(samples, reset...)
	diag := BuildQuickDiagnostic(samples, nil, &samples[len(samples)-1], nil)
	if diag.ComparisonStatus != "rolling_window_only" || diag.Comparability.Confidence != "low" {
		t.Fatalf("diagnostic = %#v, want low-confidence rolling window", diag)
	}
	if !containsString(diag.Comparability.Warnings, "quick_race_time_reset") {
		t.Fatalf("warnings = %#v, want race time reset", diag.Comparability.Warnings)
	}
}

func TestQuickDiagnosticFallsBackToRollingWindowAndMissingFields(t *testing.T) {
	samples := make([]telemetry.NormalizedTelemetry, 0, 12)
	for i := 0; i < 12; i++ {
		samples = append(samples, telemetry.NormalizedTelemetry{
			TimeMS:              int64(i) * 100,
			GameMode:            telemetry.GameModeFreeRoam,
			CarOrdinal:          7711,
			CarClass:            "A",
			Drivetrain:          "RWD",
			Gear:                1,
			SpeedKmh:            35 + float64(i),
			Throttle01:          0.95,
			RpmRatio:            0.78,
			RearCombinedSlipAvg: 1.35,
		})
	}
	events := []telemetry.DetectedEvent{testEvent("quick-roll", "launch_wheelspin", "high", 0, 900)}
	diag := BuildQuickDiagnostic(samples, events, &samples[len(samples)-1], nil)
	if diag.Status != "ready" || diag.ComparisonStatus != "rolling_window_only" {
		t.Fatalf("diagnostic = %#v, want rolling window", diag)
	}
	if diag.Comparability.Confidence != "low" || !containsString(diag.Comparability.Warnings, "quick_non_race_track_unknown") {
		t.Fatalf("comparability = %#v, want low confidence non-race warning", diag.Comparability)
	}
	if len(diag.Suggestions) == 0 || diag.Suggestions[0].CanApply {
		t.Fatalf("suggestions = %#v, want directional non-applicable suggestions", diag.Suggestions)
	}
	if diag.Suggestions[0].Amount != quickDirectionOnlyAmount || diag.Suggestions[0].TrustLevel == "" || diag.Suggestions[0].NextStep == "" {
		t.Fatalf("suggestion = %#v, want direction-only advice with trust and next step", diag.Suggestions[0])
	}
	if !containsString(diag.Suggestions[0].MissingInputs, "tune_profile") {
		t.Fatalf("suggestion missing inputs = %#v, want tune_profile", diag.Suggestions[0].MissingInputs)
	}
	if diag.GearPower.Status != "ok" || diag.GearPower.Summary != "traction_limited_power" {
		t.Fatalf("gear power = %#v, want quick rolling gear diagnostic without race mode", diag.GearPower)
	}
	if len(diag.GearPower.RecommendedActions) == 0 {
		t.Fatalf("gear power actions empty, want directional quick gear advice")
	}
	if len(diag.MissingProfileFields) == 0 {
		t.Fatalf("missing fields empty, want required tune fields for concrete values")
	}
}

func TestTunePlanDraftApplyCreatesSnapshotWithSession(t *testing.T) {
	store := openTestStore(t)
	carOrdinal := int64(7711)
	frontARB := 20.0
	profile, err := store.CreateTuneProfile(TuneProfileInput{
		CarName:    "Road Draft Car",
		CarOrdinal: &carOrdinal,
		CarClass:   "A",
		UseCase:    "Road",
		FrontARB:   &frontARB,
	})
	if err != nil {
		t.Fatalf("create profile: %v", err)
	}
	session := createIssueSummarySession(t, store, profile.ID, "Draft Session", "2026-05-18T12:00:00Z", []telemetry.DetectedEvent{
		testEvent("draft-1", "corner_entry_understeer", "high", 0, 1200),
	})

	draft, err := store.GetTunePlanDraft(session.ID)
	if err != nil {
		t.Fatalf("get draft: %v", err)
	}
	if draft.Status != "ready" || len(draft.Actions) == 0 {
		t.Fatalf("draft = %#v, want ready actions", draft)
	}
	var frontARBAction *TunePlanDraftAction
	for i := range draft.Actions {
		if draft.Actions[i].FieldKey == "frontArb" {
			frontARBAction = &draft.Actions[i]
			break
		}
	}
	if frontARBAction == nil || !frontARBAction.CanApply || frontARBAction.TargetValue == nil {
		t.Fatalf("front ARB action = %#v, want applicable", frontARBAction)
	}
	if _, err := store.ApplyTunePlanDraft(TunePlanApplyInput{SessionID: session.ID, SelectedActionIDs: []string{"forged"}}); err == nil {
		t.Fatal("expected forged action id to fail")
	}
	result, err := store.ApplyTunePlanDraft(TunePlanApplyInput{SessionID: session.ID, SelectedActionIDs: []string{frontARBAction.ID}})
	if err != nil {
		t.Fatalf("apply draft: %v", err)
	}
	if result.Profile.FrontARB == nil || math.Abs(*result.Profile.FrontARB-*frontARBAction.TargetValue) > 0.001 {
		t.Fatalf("updated front ARB = %#v, want %.2f", result.Profile.FrontARB, *frontARBAction.TargetValue)
	}
	snapshots, err := store.ListTuneProfileSnapshots(profile.ID)
	if err != nil {
		t.Fatalf("list snapshots: %v", err)
	}
	if len(snapshots) == 0 || snapshots[0].ChangeReason != "tune_plan_apply" || snapshots[0].SessionID == nil || *snapshots[0].SessionID != session.ID {
		t.Fatalf("latest snapshot = %#v, want tune_plan_apply with session", snapshots)
	}
}

func TestTunePlanDraftIncludesGearPowerDiffActions(t *testing.T) {
	store := openTestStore(t)
	carOrdinal := int64(7711)
	gear1 := 4.10
	gear2 := 2.30
	rearDiff := 48.0
	profile, err := store.CreateTuneProfile(TuneProfileInput{
		CarName:       "Gear Draft Car",
		CarOrdinal:    &carOrdinal,
		CarClass:      "A",
		UseCase:       "Road",
		Drivetrain:    "RWD",
		Gear1:         &gear1,
		Gear2:         &gear2,
		RearDiffAccel: &rearDiff,
	})
	if err != nil {
		t.Fatalf("create profile: %v", err)
	}
	samples := make([]telemetry.NormalizedTelemetry, 0, 24)
	for gear := 1; gear <= 2; gear++ {
		for i := 0; i < 12; i++ {
			samples = append(samples, telemetry.NormalizedTelemetry{
				TimeMS:              int64(len(samples)) * 100,
				GameMode:            telemetry.GameModeFreeRoam,
				CarOrdinal:          int(carOrdinal),
				CarClass:            "A",
				Drivetrain:          "RWD",
				Gear:                gear,
				SpeedKmh:            35 + float64(len(samples)),
				Throttle01:          0.95,
				RpmRatio:            0.78,
				RearCombinedSlipAvg: 1.35,
			})
		}
	}
	session := createIssueSummarySessionWithSamples(t, store, profile.ID, "Gear Draft Session", "2026-05-18T12:05:00Z", nil, samples)
	draft, err := store.GetTunePlanDraft(session.ID)
	if err != nil {
		t.Fatalf("get draft: %v", err)
	}
	var rearDiffAction *TunePlanDraftAction
	for i := range draft.Actions {
		if draft.Actions[i].FieldKey == "rearDiffAccel" && draft.Actions[i].Direction == "decrease" {
			rearDiffAction = &draft.Actions[i]
			break
		}
	}
	if rearDiffAction == nil || !rearDiffAction.CanApply || rearDiffAction.TargetValue == nil {
		t.Fatalf("draft actions = %#v, want applicable rear diff accel action", draft.Actions)
	}
	if len(draft.Actions) > 3 || hasDraftFieldDirectionConflict(draft.Actions) {
		t.Fatalf("draft actions = %#v, want at most three non-conflicting advice actions", draft.Actions)
	}
	if math.Abs(*rearDiffAction.TargetValue-46.0) > 0.001 {
		t.Fatalf("rear diff target = %.2f, want 46.00", *rearDiffAction.TargetValue)
	}
	result, err := store.ApplyTunePlanDraft(TunePlanApplyInput{SessionID: session.ID, SelectedActionIDs: []string{rearDiffAction.ID}})
	if err != nil {
		t.Fatalf("apply draft: %v", err)
	}
	if result.Profile.RearDiffAccel == nil || math.Abs(*result.Profile.RearDiffAccel-46.0) > 0.001 {
		t.Fatalf("updated rear diff = %#v, want 46.00", result.Profile.RearDiffAccel)
	}
}

func TestRoadTuningKnowledgeLoadsFromWorkbook(t *testing.T) {
	knowledge, err := LoadRoadTuningKnowledge(defaultRoadTuningKnowledgePath())
	if err != nil {
		t.Fatalf("load road tuning knowledge: %v", err)
	}
	if len(knowledge.Symptoms) < 8 || len(knowledge.Actions) < 20 {
		t.Fatalf("knowledge counts = %d symptoms / %d actions, want workbook-backed road rules", len(knowledge.Symptoms), len(knowledge.Actions))
	}
	found := false
	for _, action := range knowledge.Actions {
		if action.SymptomID == "road_exit_oversteer" && action.Item == "rear_diff_accel" && action.Direction == "decrease" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("road_exit_oversteer rear diff action missing: %#v", knowledge.Actions)
	}
}

func TestRoadTuningDecisionUsesStageSymptomsAndRearActions(t *testing.T) {
	store := openTestStore(t)
	carOrdinal := int64(7711)
	frontARB := 20.0
	rearARB := 20.0
	profile, err := store.CreateTuneProfile(TuneProfileInput{
		CarName:    "Road Decision Car",
		CarOrdinal: &carOrdinal,
		CarClass:   "A",
		UseCase:    "Road",
		Drivetrain: "RWD",
		FrontARB:   &frontARB,
		RearARB:    &rearARB,
	})
	if err != nil {
		t.Fatalf("create profile: %v", err)
	}
	session := createIssueSummarySession(t, store, profile.ID, "Entry Understeer", "2026-05-18T12:10:00Z", []telemetry.DetectedEvent{
		testEvent("entry-1", "corner_entry_understeer", "high", 0, 1200),
		testEvent("entry-2", "corner_entry_understeer", "medium", 1500, 2500),
	})
	decision, err := store.GetRoadTuningDecision(session.ID)
	if err != nil {
		t.Fatalf("road decision: %v", err)
	}
	if decision.SymptomID != "road_entry_understeer_mid_speed" || len(decision.Actions) == 0 || len(decision.Actions) > 3 {
		t.Fatalf("decision = %#v, want mid-speed entry understeer with max 3 actions", decision)
	}
	if !roadDecisionHasAction(decision, "rear_arb", "increase") {
		t.Fatalf("actions = %#v, want rear ARB increase support action", decision.Actions)
	}
	if hasRoadDecisionDirectionConflict(decision.Actions) {
		t.Fatalf("actions have direction conflict: %#v", decision.Actions)
	}
}

func TestRoadTuningDecisionEntryUndersteerSpeedBands(t *testing.T) {
	cases := []struct {
		name  string
		speed float64
		want  string
	}{
		{name: "low", speed: 70, want: "road_entry_understeer_low_speed"},
		{name: "mid", speed: 120, want: "road_entry_understeer_mid_speed"},
		{name: "high", speed: 180, want: "road_entry_understeer_high_speed"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			store := openTestStore(t)
			carOrdinal := int64(7711)
			frontARB := 20.0
			rearARB := 20.0
			frontAero := 80.0
			rearAero := 100.0
			profile, err := store.CreateTuneProfile(TuneProfileInput{
				CarName:    "Speed Band Car",
				CarOrdinal: &carOrdinal,
				CarClass:   "A",
				UseCase:    "Road",
				Drivetrain: "RWD",
				FrontARB:   &frontARB,
				RearARB:    &rearARB,
				FrontAero:  &frontAero,
				RearAero:   &rearAero,
			})
			if err != nil {
				t.Fatalf("create profile: %v", err)
			}
			event := testEvent("entry-"+tc.name, "corner_entry_understeer", "high", 0, 1200)
			event.Evidence["speed_kmh"] = tc.speed
			session := createIssueSummarySession(t, store, profile.ID, "Entry Understeer "+tc.name, "2026-05-18T13:00:00Z", []telemetry.DetectedEvent{event})
			decision, err := store.GetRoadTuningDecision(session.ID)
			if err != nil {
				t.Fatalf("road decision: %v", err)
			}
			if decision.SymptomID != tc.want {
				t.Fatalf("symptom = %q, want %q: %#v", decision.SymptomID, tc.want, decision)
			}
			if decision.Evidence["speed_band"] == 0 || decision.Evidence["speed_avg_kmh"] == 0 || decision.Evidence["speed_max_kmh"] == 0 {
				t.Fatalf("decision evidence missing speed band stats: %#v", decision.Evidence)
			}
		})
	}
}

func TestRoadTuningDecisionEntryUndersteerHighSpeedByMaxSpeed(t *testing.T) {
	store := openTestStore(t)
	carOrdinal := int64(7711)
	frontAero := 80.0
	rearAero := 100.0
	profile, err := store.CreateTuneProfile(TuneProfileInput{
		CarName:    "High Speed Max Car",
		CarOrdinal: &carOrdinal,
		CarClass:   "A",
		UseCase:    "Road",
		Drivetrain: "RWD",
		FrontAero:  &frontAero,
		RearAero:   &rearAero,
	})
	if err != nil {
		t.Fatalf("create profile: %v", err)
	}
	event := testEvent("entry-high-max", "corner_entry_understeer", "high", 0, 400)
	event.Evidence["speed_kmh"] = 120
	samples := []telemetry.NormalizedTelemetry{
		{TimeMS: 0, SpeedKmh: 130, CarOrdinal: int(carOrdinal), CarClass: "A"},
		{TimeMS: 100, SpeedKmh: 145, CarOrdinal: int(carOrdinal), CarClass: "A"},
		{TimeMS: 200, SpeedKmh: 150, CarOrdinal: int(carOrdinal), CarClass: "A"},
		{TimeMS: 300, SpeedKmh: 180, CarOrdinal: int(carOrdinal), CarClass: "A"},
	}
	session := createIssueSummarySessionWithSamples(t, store, profile.ID, "Entry High By Max", "2026-05-18T13:10:00Z", []telemetry.DetectedEvent{event}, samples)
	decision, err := store.GetRoadTuningDecision(session.ID)
	if err != nil {
		t.Fatalf("road decision: %v", err)
	}
	if decision.SymptomID != "road_entry_understeer_high_speed" {
		t.Fatalf("symptom = %q, want high-speed by max speed: %#v", decision.SymptomID, decision)
	}
	if !roadDecisionHasAction(decision, "front_and_rear_aero", "increase") {
		t.Fatalf("actions = %#v, want high-speed aero action", decision.Actions)
	}
}

func TestEvidenceRuleMatchesNumericComparisonsAndIgnoresUnknownTags(t *testing.T) {
	evidence := map[string]float64{"speed_avg_kmh": 120, "speed_max_kmh": 140, "speed_band": 2}
	if !evidenceRuleMatches("speed_avg_kmh>=90;speed_avg_kmh<160;unknown_tag;front_slip>rear_slip", evidence) {
		t.Fatal("expected numeric rule with unknown tags to match")
	}
	if evidenceRuleMatches("speed_avg_kmh<90", evidence) {
		t.Fatal("expected numeric rule to reject low-speed condition")
	}
	if !evidenceRuleMatches("missing_numeric>2", evidence) {
		t.Fatal("expected missing numeric evidence to be ignored")
	}
}

func TestRoadTuningDecisionPowerOversteerPrioritizesDiffAndGearing(t *testing.T) {
	store := openTestStore(t)
	carOrdinal := int64(7711)
	rearDiff := 55.0
	gear2 := 2.40
	profile, err := store.CreateTuneProfile(TuneProfileInput{
		CarName:       "Power Oversteer Car",
		CarOrdinal:    &carOrdinal,
		CarClass:      "A",
		UseCase:       "Road",
		Drivetrain:    "RWD",
		RearDiffAccel: &rearDiff,
		Gear2:         &gear2,
	})
	if err != nil {
		t.Fatalf("create profile: %v", err)
	}
	event := testEvent("exit-1", "corner_exit_oversteer", "high", 0, 1600)
	event.Segment = "corner_exit"
	event.Evidence["gear"] = 2
	session := createIssueSummarySession(t, store, profile.ID, "Exit Oversteer", "2026-05-18T12:20:00Z", []telemetry.DetectedEvent{event})
	decision, err := store.GetRoadTuningDecision(session.ID)
	if err != nil {
		t.Fatalf("road decision: %v", err)
	}
	if decision.SymptomID != "road_exit_oversteer" {
		t.Fatalf("symptom = %q, want road_exit_oversteer: %#v", decision.SymptomID, decision)
	}
	if !roadDecisionHasAction(decision, "rear_diff_accel", "decrease") {
		t.Fatalf("actions = %#v, want rear diff accel decrease", decision.Actions)
	}
	draft, err := store.GetTunePlanDraft(session.ID)
	if err != nil {
		t.Fatalf("draft: %v", err)
	}
	if !draftHasField(draft, "rearDiffAccel") {
		t.Fatalf("draft actions = %#v, want rearDiffAccel concrete action", draft.Actions)
	}
}

func TestRoadTuningDecisionFitUsesTelemetryRetestOnly(t *testing.T) {
	store := openTestStore(t)
	carOrdinal := int64(7711)
	frontARB := 20.0
	profile, err := store.CreateTuneProfile(TuneProfileInput{CarName: "Retest Car", CarOrdinal: &carOrdinal, CarClass: "A", UseCase: "Road", FrontARB: &frontARB})
	if err != nil {
		t.Fatalf("create profile: %v", err)
	}
	_ = createIssueSummarySession(t, store, profile.ID, "Previous Retest", "2026-05-18T10:00:00Z", []telemetry.DetectedEvent{
		testEvent("fb-1", "corner_entry_understeer", "high", 0, 1400),
		testEvent("fb-2", "corner_entry_understeer", "high", 1800, 3200),
	})
	current := createIssueSummarySession(t, store, profile.ID, "Current Retest", "2026-05-18T11:00:00Z", []telemetry.DetectedEvent{
		testEvent("fb-3", "corner_entry_understeer", "medium", 0, 500),
	})
	decision, err := store.GetRoadTuningDecision(current.ID)
	if err != nil {
		t.Fatalf("road decision: %v", err)
	}
	if decision.FitVerdict != "improved" {
		t.Fatalf("fit verdict = %q, want improved", decision.FitVerdict)
	}
}

func TestRoadTuningDecisionCheckActionsAreNotAutoApplicable(t *testing.T) {
	store := openTestStore(t)
	carOrdinal := int64(7711)
	frontPressure := 2.0
	rearPressure := 2.0
	profile, err := store.CreateTuneProfile(TuneProfileInput{
		CarName:           "Tire Temp Car",
		CarOrdinal:        &carOrdinal,
		CarClass:          "A",
		UseCase:           "Road",
		FrontTirePressure: &frontPressure,
		RearTirePressure:  &rearPressure,
	})
	if err != nil {
		t.Fatalf("create profile: %v", err)
	}
	event := testEvent("tire-1", "tire_overheat", "medium", 0, 2000)
	event.Segment = "mid_corner"
	session := createIssueSummarySession(t, store, profile.ID, "Tire Temp", "2026-05-18T12:40:00Z", []telemetry.DetectedEvent{event})
	decision, err := store.GetRoadTuningDecision(session.ID)
	if err != nil {
		t.Fatalf("road decision: %v", err)
	}
	if len(decision.Actions) == 0 || decision.Actions[0].CanAutoApply {
		t.Fatalf("decision actions = %#v, want first check action not auto-applicable", decision.Actions)
	}
	draft, err := store.GetTunePlanDraft(session.ID)
	if err != nil {
		t.Fatalf("draft: %v", err)
	}
	if len(draft.Actions) == 0 || draft.Actions[0].CanApply || draft.Actions[0].BlockedReason != "manual_review_required" {
		t.Fatalf("draft actions = %#v, want manual review block", draft.Actions)
	}
}

func TestTunePlanDraftRejectsVehicleMismatch(t *testing.T) {
	store := openTestStore(t)
	carOrdinal := int64(9999)
	frontARB := 20.0
	profile, err := store.CreateTuneProfile(TuneProfileInput{
		CarName:    "Mismatch Car",
		CarOrdinal: &carOrdinal,
		CarClass:   "A",
		UseCase:    "Road",
		FrontARB:   &frontARB,
	})
	if err != nil {
		t.Fatalf("create profile: %v", err)
	}
	session := createIssueSummarySession(t, store, profile.ID, "Mismatch Session", "2026-05-18T12:30:00Z", []telemetry.DetectedEvent{
		testEvent("mismatch-1", "corner_entry_understeer", "high", 0, 1200),
	})
	draft, err := store.GetTunePlanDraft(session.ID)
	if err != nil {
		t.Fatalf("get draft: %v", err)
	}
	if draft.Status != "vehicle_mismatch" {
		t.Fatalf("draft status = %q, want vehicle_mismatch", draft.Status)
	}
	if _, err := store.ApplyTunePlanDraft(TunePlanApplyInput{SessionID: session.ID, SelectedActionIDs: []string{"anything"}}); err == nil {
		t.Fatal("expected mismatch apply to fail")
	}
}

func TestRetestEvaluationComparesPreviousSession(t *testing.T) {
	store := openTestStore(t)
	carOrdinal := int64(7711)
	profile, err := store.CreateTuneProfile(TuneProfileInput{CarName: "Retest Car", CarOrdinal: &carOrdinal, CarClass: "A", UseCase: "Road"})
	if err != nil {
		t.Fatalf("create profile: %v", err)
	}
	previous := createIssueSummarySession(t, store, profile.ID, "Previous Retest", "2026-05-18T10:00:00Z", []telemetry.DetectedEvent{
		testEvent("r1", "corner_entry_understeer", "high", 0, 1000),
		testEvent("r2", "corner_entry_understeer", "medium", 1500, 2600),
		testEvent("r3", "corner_entry_understeer", "medium", 3000, 4100),
	})
	current := createIssueSummarySession(t, store, profile.ID, "Current Retest", "2026-05-18T11:00:00Z", []telemetry.DetectedEvent{
		testEvent("r4", "corner_entry_understeer", "medium", 0, 600),
	})

	evaluation, err := store.GetRetestEvaluation(current.ID)
	if err != nil {
		t.Fatalf("retest evaluation: %v", err)
	}
	if evaluation.BaselineSession == nil || evaluation.BaselineSession.ID != previous.ID {
		t.Fatalf("baseline = %#v, want previous", evaluation.BaselineSession)
	}
	if evaluation.Status != "improved" {
		t.Fatalf("status = %q, metrics=%#v, want improved", evaluation.Status, evaluation.Metrics)
	}
}

func TestRetestEvaluationWorsenedCreatesRollbackDraft(t *testing.T) {
	store := openTestStore(t)
	carOrdinal := int64(7711)
	frontARB := 30.0
	profile, err := store.CreateTuneProfile(TuneProfileInput{CarName: "Rollback Car", CarOrdinal: &carOrdinal, CarClass: "A", UseCase: "Road", FrontARB: &frontARB})
	if err != nil {
		t.Fatalf("create profile: %v", err)
	}
	baseline := createRetestSpeedSession(t, store, profile.ID, "Before Tune Apply", "2099-05-18T10:00:00Z", 120)
	input := profile.ToInput()
	changedARB := 20.0
	input.FrontARB = &changedARB
	if _, err := store.updateTuneProfileWithSession(profile.ID, input, "tune_plan_apply", &baseline.ID); err != nil {
		t.Fatalf("apply tune profile update: %v", err)
	}
	current := createRetestSpeedSession(t, store, profile.ID, "After Tune Apply", "2099-05-18T11:00:00Z", 100)

	evaluation, err := store.GetRetestEvaluation(current.ID)
	if err != nil {
		t.Fatalf("retest evaluation: %v", err)
	}
	if evaluation.Status != "worsened" {
		t.Fatalf("status = %q, metrics=%#v, want worsened", evaluation.Status, evaluation.Metrics)
	}
	if len(evaluation.RollbackActions) == 0 || evaluation.RollbackActions[0].FieldKey != "frontArb" {
		t.Fatalf("rollback actions = %#v, want frontArb rollback", evaluation.RollbackActions)
	}
	if evaluation.RollbackActions[0].TargetValue == nil || math.Abs(*evaluation.RollbackActions[0].TargetValue-frontARB) > 0.001 {
		t.Fatalf("rollback target = %#v, want %.1f", evaluation.RollbackActions[0].TargetValue, frontARB)
	}

	draft, err := store.GetTunePlanDraft(current.ID)
	if err != nil {
		t.Fatalf("tune plan draft: %v", err)
	}
	if len(draft.Actions) == 0 || draft.Actions[0].Source != "retest_guard" || draft.Actions[0].FieldKey != "frontArb" || !draft.Actions[0].CanApply {
		t.Fatalf("draft actions = %#v, want applicable rollback first", draft.Actions)
	}
	if draft.Actions[0].AdviceLayer != adviceLayerRollback || draft.Actions[0].Rationale == "" {
		t.Fatalf("rollback action = %#v, want rollback layer and rationale", draft.Actions[0])
	}
	for _, action := range draft.Actions[1:] {
		if action.CanApply {
			t.Fatalf("draft actions = %#v, want non-rollback actions guarded after worsened retest", draft.Actions)
		}
	}
}

func TestAnalyzeRoadStrategySessionsLimitAndDistribution(t *testing.T) {
	store := openTestStore(t)
	session := createIssueSummarySession(t, store, 0, "Wheelspin", "2026-05-18T10:00:00Z", nil)
	samples := launchWheelspinSamples()
	_, err := store.FinalizeTelemetrySession(SessionFinalizeInput{
		SessionID:  session.ID,
		EndedAt:    "2026-05-18T10:01:00Z",
		DurationMS: samples[len(samples)-1].TimeMS,
		GameMode:   telemetry.GameModeFreeRoam,
	}, nil, samples)
	if err != nil {
		t.Fatalf("finalize samples: %v", err)
	}
	templates, err := store.ListStrategyTemplates()
	if err != nil {
		t.Fatalf("list strategy templates: %v", err)
	}
	var templateID int64
	for _, template := range templates {
		if template.Name == "Road Racing" {
			templateID = template.ID
			break
		}
	}
	if templateID == 0 {
		t.Fatalf("Road Racing strategy template not found: %#v", templates)
	}
	analysis, err := store.AnalyzeRoadStrategySessions([]int64{session.ID}, templateID)
	if err != nil {
		t.Fatalf("analyze strategy: %v", err)
	}
	if analysis.TotalEvents == 0 || len(analysis.EventDistribution) == 0 {
		t.Fatalf("analysis did not detect events: %#v", analysis)
	}
	if _, err := store.AnalyzeRoadStrategySessions([]int64{1, 2, 3, 4, 5, 6}, templateID); err == nil {
		t.Fatal("expected more than 5 sessions to fail")
	}
}

func createRoadEvalSession(t *testing.T, store *Store, name string, driverMode string, samples []telemetry.NormalizedTelemetry, events []telemetry.DetectedEvent) *TelemetrySession {
	t.Helper()
	carOrdinal := int64(9001)
	session, err := store.CreateTelemetrySession(SessionStartInput{
		SessionName: name,
		StartedAt:   "2026-05-18T10:00:00Z",
		GameMode:    telemetry.GameModeFreeRoam,
		TestConditions: TestConditions{
			DriverMode: driverMode,
		},
		SessionVehicleSnapshot: SessionVehicleSnapshot{
			CarOrdinal: &carOrdinal,
			CarClass:   "S1",
		},
	})
	if err != nil {
		t.Fatalf("create %s session: %v", name, err)
	}
	end := int64(0)
	if len(samples) > 0 {
		end = samples[len(samples)-1].TimeMS
	}
	_, err = store.FinalizeTelemetrySession(SessionFinalizeInput{
		SessionID:  session.ID,
		EndedAt:    "2026-05-18T10:03:00Z",
		DurationMS: end,
		GameMode:   telemetry.GameModeFreeRoam,
		DriverModeDetection: DriverModeDetection{
			Mode:       driverMode,
			Confidence: 0.9,
			Summary:    "test_fixture",
			Evidence:   map[string]float64{},
		},
		SessionVehicleSnapshot: SessionVehicleSnapshot{
			CarOrdinal: &carOrdinal,
			CarClass:   "S1",
		},
	}, events, samples)
	if err != nil {
		t.Fatalf("finalize %s session: %v", name, err)
	}
	loaded, err := store.GetTelemetrySession(session.ID)
	if err != nil {
		t.Fatalf("reload %s session: %v", name, err)
	}
	return loaded
}

func createTrackProfileSession(t *testing.T, store *Store, name string, carOrdinal int64, carClass string, carPI int64, drivetrain string) *TelemetrySession {
	t.Helper()
	session, err := store.CreateTelemetrySession(SessionStartInput{
		SessionName: name,
		StartedAt:   "2026-05-18T10:00:00Z",
		GameMode:    telemetry.GameModeRace,
		SessionVehicleSnapshot: SessionVehicleSnapshot{
			CarOrdinal: &carOrdinal,
			CarClass:   carClass,
			CarPI:      &carPI,
			Drivetrain: drivetrain,
		},
	})
	if err != nil {
		t.Fatalf("create %s session: %v", name, err)
	}
	loaded, err := store.FinalizeTelemetrySession(SessionFinalizeInput{
		SessionID:  session.ID,
		EndedAt:    "2026-05-18T10:02:00Z",
		DurationMS: 120000,
		GameMode:   telemetry.GameModeRace,
		SessionVehicleSnapshot: SessionVehicleSnapshot{
			CarOrdinal: &carOrdinal,
			CarClass:   carClass,
			CarPI:      &carPI,
			Drivetrain: drivetrain,
		},
		DriverModeDetection: DriverModeDetection{
			Mode:       driverModePlayer,
			Confidence: 0.8,
			Summary:    "test_fixture",
			Evidence:   map[string]float64{},
		},
	}, nil, nil)
	if err != nil {
		t.Fatalf("finalize %s session: %v", name, err)
	}
	return loaded
}

func insertTrackProfileRun(t *testing.T, store *Store, trackID int64, sessionID int64, durationMS int64, driverMode string, confidence float64, valid bool) {
	t.Helper()
	avgSpeed := 120.0
	maxSpeed := 180.0
	_, err := store.insertBenchmarkRun(BenchmarkRun{
		SessionID:              sessionID,
		TrackID:                trackID,
		StartMS:                0,
		EndMS:                  durationMS,
		DurationMS:             durationMS,
		Confidence:             confidence,
		AvgSpeedKmh:            &avgSpeed,
		MaxSpeedKmh:            &maxSpeed,
		RouteProgress01:        floatPtrFromValue(1),
		DriverMode:             driverMode,
		DriverModeConfidence:   confidence,
		DriverModeEvidenceJSON: "{}",
		Valid:                  valid,
	})
	if err != nil {
		t.Fatalf("insert benchmark run: %v", err)
	}
}

func floatPtrFromValue(value float64) *float64 {
	return &value
}

func straightRouteSamples(from float64, to float64, count int) []telemetry.NormalizedTelemetry {
	samples := make([]telemetry.NormalizedTelemetry, 0, count)
	step := 0.0
	if count > 1 {
		step = (to - from) / float64(count-1)
	}
	for i := 0; i < count; i++ {
		x := from + step*float64(i)
		samples = append(samples, telemetry.NormalizedTelemetry{
			TimeMS:        int64(i) * 1000,
			IsRaceOn:      false,
			GameMode:      telemetry.GameModeFreeRoam,
			SpeedKmh:      80,
			PositionX:     x,
			PositionY:     0,
			PositionZ:     0,
			DrivingLine01: 0.25,
			CarOrdinal:    9001,
			CarClass:      "S1",
		})
	}
	return samples
}

func circuitRouteSamples(radius float64, laps int, pointsPerLap int) []telemetry.NormalizedTelemetry {
	total := laps*pointsPerLap + 1
	samples := make([]telemetry.NormalizedTelemetry, 0, total)
	for i := 0; i < total; i++ {
		angle := (float64(i) / float64(pointsPerLap)) * 2 * math.Pi
		x := radius * math.Cos(angle)
		z := radius * math.Sin(angle)
		samples = append(samples, telemetry.NormalizedTelemetry{
			TimeMS:           int64(i) * 1000,
			IsRaceOn:         false,
			GameMode:         telemetry.GameModeFreeRoam,
			SpeedKmh:         90,
			PositionX:        x,
			PositionY:        0,
			PositionZ:        z,
			DistanceTraveled: float64(i) * (2 * math.Pi * radius / float64(pointsPerLap)) * 1.25,
			CurrentRaceTime:  float64(i),
			DrivingLine01:    0.4,
			CarOrdinal:       9001,
			CarClass:         "S1",
		})
	}
	return samples
}

func circuitRouteSamplesWithLapNumbers(radius float64, laps int, pointsPerLap int, startOffset int) []telemetry.NormalizedTelemetry {
	total := laps*pointsPerLap + 1
	samples := make([]telemetry.NormalizedTelemetry, 0, total)
	for i := 0; i < total; i++ {
		absoluteIndex := startOffset + i
		angle := (float64(absoluteIndex) / float64(pointsPerLap)) * 2 * math.Pi
		x := radius * math.Cos(angle)
		z := radius * math.Sin(angle)
		samples = append(samples, telemetry.NormalizedTelemetry{
			TimeMS:           int64(i) * 1000,
			IsRaceOn:         true,
			GameMode:         telemetry.GameModeRace,
			SpeedKmh:         95,
			PositionX:        x,
			PositionY:        0,
			PositionZ:        z,
			DistanceTraveled: float64(i) * (2 * math.Pi * radius / float64(pointsPerLap)),
			CurrentRaceTime:  float64(i),
			LapNumber:        absoluteIndex/pointsPerLap + 1,
			RacePosition:     1,
			DrivingLine01:    0.5,
			CarOrdinal:       9001,
			CarClass:         "S1",
		})
	}
	return samples
}

func circleBenchmarkPoints(radius float64, pointsPerLap int, startOffset int) []BenchmarkPoint {
	points := make([]BenchmarkPoint, 0, pointsPerLap+1)
	for i := 0; i <= pointsPerLap; i++ {
		angle := (float64(startOffset+i) / float64(pointsPerLap)) * 2 * math.Pi
		points = append(points, BenchmarkPoint{
			X: radius * math.Cos(angle),
			Y: 0,
			Z: radius * math.Sin(angle),
		})
	}
	return points
}

func createIssueSummarySession(t *testing.T, store *Store, profileID int64, name string, startedAt string, events []telemetry.DetectedEvent) *TelemetrySession {
	return createIssueSummarySessionWithSamples(t, store, profileID, name, startedAt, events, nil)
}

func createIssueSummarySessionWithSamples(t *testing.T, store *Store, profileID int64, name string, startedAt string, events []telemetry.DetectedEvent, samples []telemetry.NormalizedTelemetry) *TelemetrySession {
	t.Helper()
	var tuneProfileID *int64
	if profileID > 0 {
		tuneProfileID = &profileID
	}
	carOrdinal := int64(7711)
	session, err := store.CreateTelemetrySession(SessionStartInput{
		TuneProfileID: tuneProfileID,
		SessionName:   name,
		TrackName:     "Road Test",
		StartedAt:     startedAt,
		GameMode:      telemetry.GameModeFreeRoam,
		TestConditions: TestConditions{
			DriverMode: "player",
		},
		SessionVehicleSnapshot: SessionVehicleSnapshot{
			CarOrdinal: &carOrdinal,
			CarClass:   "A",
		},
	})
	if err != nil {
		t.Fatalf("create issue session: %v", err)
	}
	_, err = store.FinalizeTelemetrySession(SessionFinalizeInput{
		SessionID:  session.ID,
		EndedAt:    startedAt,
		DurationMS: 60000,
		GameMode:   telemetry.GameModeFreeRoam,
		SessionVehicleSnapshot: SessionVehicleSnapshot{
			CarOrdinal: &carOrdinal,
			CarClass:   "A",
		},
	}, events, samples)
	if err != nil {
		t.Fatalf("finalize issue session: %v", err)
	}
	loaded, err := store.GetTelemetrySession(session.ID)
	if err != nil {
		t.Fatalf("reload issue session: %v", err)
	}
	return loaded
}

func createRetestSpeedSession(t *testing.T, store *Store, profileID int64, name string, startedAt string, avgSpeed float64) *TelemetrySession {
	t.Helper()
	carOrdinal := int64(7711)
	session, err := store.CreateTelemetrySession(SessionStartInput{
		TuneProfileID: &profileID,
		SessionName:   name,
		TrackName:     "Road Test",
		StartedAt:     startedAt,
		GameMode:      telemetry.GameModeFreeRoam,
		TestConditions: TestConditions{
			DriverMode: "player",
		},
		SessionVehicleSnapshot: SessionVehicleSnapshot{
			CarOrdinal: &carOrdinal,
			CarClass:   "A",
		},
	})
	if err != nil {
		t.Fatalf("create retest session: %v", err)
	}
	maxSpeed := avgSpeed + 5
	_, err = store.FinalizeTelemetrySession(SessionFinalizeInput{
		SessionID:   session.ID,
		EndedAt:     startedAt,
		DurationMS:  60000,
		AvgSpeedKmh: &avgSpeed,
		MaxSpeedKmh: &maxSpeed,
		GameMode:    telemetry.GameModeFreeRoam,
		SessionVehicleSnapshot: SessionVehicleSnapshot{
			CarOrdinal: &carOrdinal,
			CarClass:   "A",
		},
	}, nil, retestSpeedSamples(avgSpeed))
	if err != nil {
		t.Fatalf("finalize retest session: %v", err)
	}
	loaded, err := store.GetTelemetrySession(session.ID)
	if err != nil {
		t.Fatalf("reload retest session: %v", err)
	}
	return loaded
}

func gearComparisonSamples(gear int, startSpeed float64) []telemetry.NormalizedTelemetry {
	samples := make([]telemetry.NormalizedTelemetry, 0, 12)
	for i := 0; i < 12; i++ {
		samples = append(samples, telemetry.NormalizedTelemetry{
			TimeMS:     int64(i) * 100,
			GameMode:   telemetry.GameModeFreeRoam,
			Gear:       gear,
			SpeedKmh:   startSpeed + float64(i),
			Throttle01: 0.9,
			RpmRatio:   0.72,
			CarOrdinal: 7711,
			CarClass:   "A",
		})
	}
	return samples
}

func gearComparisonByType(comparisons []GearPowerComparison, comparisonType string) *GearPowerComparison {
	for i := range comparisons {
		if comparisons[i].Type == comparisonType {
			return &comparisons[i]
		}
	}
	return nil
}

func testEvent(id string, eventType string, severity string, startMS int64, endMS int64) telemetry.DetectedEvent {
	return telemetry.DetectedEvent{
		ID:         id,
		Type:       eventType,
		Severity:   severity,
		StartMS:    startMS,
		EndMS:      endMS,
		DurationMS: endMS - startMS,
		Segment:    "corner_entry",
		Evidence: map[string]float64{
			"speed_kmh":           95,
			"front_combined_slip": 1.2,
			"rear_combined_slip":  0.5,
		},
		SuggestedActions: []telemetry.SuggestedAction{
			{Priority: 0, Category: "suspension", Item: "front_arb", Direction: "decrease", Amount: "0.5-1.0", Reason: "increase front grip on entry"},
		},
	}
}

func quickLapSamples(lapNumber int, startMS int64, count int, speed float64) []telemetry.NormalizedTelemetry {
	samples := make([]telemetry.NormalizedTelemetry, 0, count)
	for i := 0; i < count; i++ {
		samples = append(samples, telemetry.NormalizedTelemetry{
			TimeMS:          startMS + int64(i)*1000,
			GameMode:        telemetry.GameModeRace,
			IsRaceOn:        true,
			LapNumber:       lapNumber,
			CurrentLap:      float64(i + 1),
			CurrentRaceTime: float64(startMS)/1000 + float64(i+1),
			CarOrdinal:      7711,
			CarClass:        "A",
			Drivetrain:      "RWD",
			Gear:            3,
			SpeedKmh:        speed,
			Throttle01:      0.7,
			RpmRatio:        0.7,
			DrivingLine01:   0.4,
		})
	}
	return samples
}

func launchWheelspinSamples() []telemetry.NormalizedTelemetry {
	samples := make([]telemetry.NormalizedTelemetry, 0, 6)
	for i := 0; i < 6; i++ {
		samples = append(samples, telemetry.NormalizedTelemetry{
			TimeMS:           int64(i) * 200,
			GameMode:         telemetry.GameModeFreeRoam,
			Gear:             1,
			SpeedKmh:         35,
			Throttle01:       0.95,
			RpmRatio:         0.75,
			RearSlipRatioAvg: 1.35,
			CarClass:         "A",
			CarOrdinal:       7711,
		})
	}
	return samples
}

func gearPowerFindingSamples(gears []int, rpmRatio float64, baseSpeed float64, rearSlip float64) []telemetry.NormalizedTelemetry {
	samples := make([]telemetry.NormalizedTelemetry, 0, len(gears)*12)
	for _, gear := range gears {
		for i := 0; i < 12; i++ {
			samples = append(samples, telemetry.NormalizedTelemetry{
				TimeMS:              int64(len(samples)) * 100,
				Gear:                gear,
				SpeedKmh:            baseSpeed + float64(gear*8) + float64(i),
				Throttle01:          0.9,
				RpmRatio:            rpmRatio,
				RearCombinedSlipAvg: rearSlip,
			})
		}
	}
	return samples
}

func retestSpeedSamples(speed float64) []telemetry.NormalizedTelemetry {
	samples := make([]telemetry.NormalizedTelemetry, 0, 20)
	for i := 0; i < 20; i++ {
		samples = append(samples, telemetry.NormalizedTelemetry{
			TimeMS:     int64(i) * 1000,
			SpeedKmh:   speed,
			Throttle01: 0.7,
			RpmRatio:   0.7,
			Gear:       4,
			CarClass:   "A",
		})
	}
	return samples
}

func hasSuggestedAction(actions []telemetry.SuggestedAction, item string, direction string) bool {
	for _, action := range actions {
		if action.Item == item && action.Direction == direction {
			return true
		}
	}
	return false
}

func roadDecisionHasAction(decision *RoadTuningDecision, item string, direction string) bool {
	if decision == nil {
		return false
	}
	for _, action := range decision.Actions {
		if action.Item == item && action.Direction == direction {
			return true
		}
	}
	return false
}

func hasRoadDecisionDirectionConflict(actions []RoadTuningDecisionAction) bool {
	seen := map[string]string{}
	for _, action := range actions {
		key := actionConflictKey(telemetry.SuggestedAction{Category: action.Category, Item: action.Item, Direction: action.Direction, Amount: action.Amount}, action.Evidence)
		if key == "" {
			key = action.Item
		}
		if existing, ok := seen[key]; ok && existing != action.Direction {
			return true
		}
		seen[key] = action.Direction
	}
	return false
}

func draftHasField(draft *TunePlanDraft, fieldKey string) bool {
	if draft == nil {
		return false
	}
	for _, action := range draft.Actions {
		if action.FieldKey == fieldKey {
			return true
		}
	}
	return false
}

func hasDraftFieldDirectionConflict(actions []TunePlanDraftAction) bool {
	seen := map[string]string{}
	for _, action := range actions {
		if action.FieldKey == "" {
			continue
		}
		if existing, ok := seen[action.FieldKey]; ok && existing != action.Direction {
			return true
		}
		seen[action.FieldKey] = action.Direction
	}
	return false
}

func lastInsertID(db interface {
	QueryRow(query string, args ...any) *sql.Row
}) (int64, error) {
	var id int64
	err := db.QueryRow(`SELECT last_insert_rowid()`).Scan(&id)
	return id, err
}

func countRows(t *testing.T, store *Store, table string) int {
	t.Helper()
	var count int
	if err := store.db.QueryRow(`SELECT COUNT(*) FROM ` + table).Scan(&count); err != nil {
		t.Fatalf("count %s: %v", table, err)
	}
	return count
}

func openTestStore(t *testing.T) *Store {
	t.Helper()
	store, err := Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("open test store: %v", err)
	}
	t.Cleanup(func() {
		if err := store.Close(); err != nil && !strings.Contains(err.Error(), "closed") {
			t.Fatalf("close store: %v", err)
		}
	})
	return store
}
