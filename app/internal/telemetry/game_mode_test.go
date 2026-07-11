package telemetry

import (
	"testing"
	"time"
)

func TestGameModeTrackerClassifiesMenuFreeRoamAndRace(t *testing.T) {
	tracker := NewGameModeTracker()
	base := time.Unix(0, 0)

	if got := tracker.Observe(NormalizedTelemetry{}, base); got != GameModeMenu {
		t.Fatalf("menu mode = %q", got)
	}

	tracker.Reset()
	freeRoam := NormalizedTelemetry{TimeMS: 1000, IsRaceOn: false, CarOrdinal: 42, CarClass: "S1"}
	if got := tracker.Observe(freeRoam, base); got != GameModeUnknown {
		t.Fatalf("initial stationary free roam = %q, want unknown during debounce", got)
	}
	freeRoam.TimeMS = 2600
	if got := tracker.Observe(freeRoam, base.Add(1600*time.Millisecond)); got != GameModeFreeRoam {
		t.Fatalf("debounced free roam = %q", got)
	}

	race := freeRoam
	race.TimeMS = 2700
	race.RacePosition = 3
	if got := tracker.Observe(race, base.Add(1700*time.Millisecond)); got != GameModeRace {
		t.Fatalf("race mode = %q", got)
	}
}

func TestGameModeTrackerDetectsFreeRoamActivityWithoutRaceFlag(t *testing.T) {
	tracker := NewGameModeTracker()
	frame := NormalizedTelemetry{
		TimeMS:     1000,
		IsRaceOn:   false,
		CarOrdinal: 7,
		CarClass:   "A",
		SpeedKmh:   12,
	}
	if got := tracker.Observe(frame, time.Unix(0, 0)); got != GameModeFreeRoam {
		t.Fatalf("active free roam = %q", got)
	}
}

func TestGameModeTrackerDoesNotTreatDistanceTraveledAsRaceEvidence(t *testing.T) {
	tracker := NewGameModeTracker()
	frame := NormalizedTelemetry{
		TimeMS:           1000,
		IsRaceOn:         true,
		CarOrdinal:       3629,
		CarClass:         "S1",
		CarPI:            797,
		SpeedKmh:         35,
		DistanceTraveled: 5740.6,
	}
	if got := tracker.Observe(frame, time.Unix(0, 0)); got != GameModeFreeRoam {
		t.Fatalf("distance-only telemetry mode = %q, want free roam", got)
	}
}

func TestGameModeTrackerStillUsesRaceEvidence(t *testing.T) {
	tracker := NewGameModeTracker()
	frame := NormalizedTelemetry{
		TimeMS:       1000,
		IsRaceOn:     true,
		CarOrdinal:   3629,
		CarClass:     "S1",
		CarPI:        797,
		SpeedKmh:     35,
		CurrentLap:   12.5,
		LapNumber:    1,
		RacePosition: 3,
	}
	if got := tracker.Observe(frame, time.Unix(0, 0)); got != GameModeRace {
		t.Fatalf("race telemetry mode = %q, want race", got)
	}
}

func TestGameModeTrackerDebouncesRaceExitWhenDrivingLineCrossesZero(t *testing.T) {
	tracker := NewGameModeTracker()
	race := NormalizedTelemetry{
		TimeMS:           1000,
		CarOrdinal:       3629,
		CarClass:         "S1",
		CarPI:            797,
		SpeedKmh:         140,
		DistanceTraveled: 100,
		DrivingLine01:    0.3,
	}
	if got := tracker.Observe(race, time.Unix(0, 0)); got != GameModeRace {
		t.Fatalf("race mode = %q", got)
	}

	lineGap := race
	lineGap.TimeMS = 1100
	lineGap.DistanceTraveled = 110
	lineGap.DrivingLine01 = 0
	if got := tracker.Observe(lineGap, time.Unix(0, 0)); got != GameModeRace {
		t.Fatalf("brief driving-line gap mode = %q, want race", got)
	}

	lineGap.TimeMS = 5200
	lineGap.DistanceTraveled = 180
	if got := tracker.Observe(lineGap, time.Unix(0, 0)); got != GameModeFreeRoam {
		t.Fatalf("sustained race-exit mode = %q, want free roam", got)
	}
}

func TestGameModeTrackerUsesDistanceResetAsRaceExitCue(t *testing.T) {
	tracker := NewGameModeTracker()
	race := NormalizedTelemetry{
		TimeMS:           1000,
		CarOrdinal:       3629,
		CarClass:         "S1",
		CarPI:            797,
		SpeedKmh:         35,
		DistanceTraveled: 777,
		DrivingLine01:    -0.9,
	}
	if got := tracker.Observe(race, time.Unix(0, 0)); got != GameModeRace {
		t.Fatalf("race mode = %q", got)
	}

	exit := race
	exit.TimeMS = 1100
	exit.DistanceTraveled = 0
	exit.DrivingLine01 = 0
	if got := tracker.Observe(exit, time.Unix(0, 0)); got != GameModeFreeRoam {
		t.Fatalf("distance-reset race exit mode = %q, want free roam", got)
	}
}

func TestGameModeTrackerResetsOnClockRollback(t *testing.T) {
	tracker := NewGameModeTracker()
	race := NormalizedTelemetry{TimeMS: 5000, CarOrdinal: 9, CarClass: "S2", RacePosition: 1}
	if got := tracker.Observe(race, time.Unix(0, 0)); got != GameModeRace {
		t.Fatalf("race mode = %q", got)
	}
	older := NormalizedTelemetry{TimeMS: 1000, CarOrdinal: 9, CarClass: "S2"}
	if got := tracker.Observe(older, time.Unix(0, 0)); got != GameModeUnknown {
		t.Fatalf("rollback mode = %q, want fresh unknown free roam candidate", got)
	}
}
