package telemetry

import (
	"math"
	"strings"
	"time"
)

const (
	GameModeUnknown  = "unknown"
	GameModeMenu     = "menu"
	GameModeFreeRoam = "free_roam"
	GameModeRace     = "race"
	GameModeMixed    = "mixed"
)

type rawGameMode string

type GameModeTracker struct {
	stable             string
	candidate          rawGameMode
	candidateStartedAt time.Time
	candidateFrame     *NormalizedTelemetry
	lastFrameAt        time.Time
	lastFrame          *NormalizedTelemetry
}

func NewGameModeTracker() GameModeTracker {
	return GameModeTracker{stable: GameModeUnknown}
}

func (t *GameModeTracker) Reset() {
	*t = NewGameModeTracker()
}

func (t *GameModeTracker) Observe(frame NormalizedTelemetry, receivedAt time.Time) string {
	if strings.TrimSpace(t.stable) == "" {
		t.stable = GameModeUnknown
	}
	frameAt := gameModeFrameTime(frame, receivedAt)
	if !t.lastFrameAt.IsZero() && frameAt.Add(time.Second).Before(t.lastFrameAt) {
		t.Reset()
	}

	previousFrame := t.lastFrame
	rawMode := RawGameMode(frame)
	hasActivity := hasDrivingActivity(frame, previousFrame)
	hasRaceExit := hasRaceExitCue(frame, previousFrame)
	t.lastFrame = frameCopy(frame)
	t.lastFrameAt = frameAt

	if rawMode == GameModeMenu || rawMode == GameModeRace {
		t.stable = rawMode
		t.clearCandidate()
		return t.stable
	}

	if t.stable == GameModeFreeRoam {
		t.clearCandidate()
		return GameModeFreeRoam
	}

	if t.candidate != GameModeFreeRoam {
		t.candidate = GameModeFreeRoam
		t.candidateStartedAt = frameAt
		t.candidateFrame = frameCopy(frame)
	}

	confirmAfter := 1500 * time.Millisecond
	if t.stable != GameModeUnknown {
		confirmAfter = 4 * time.Second
	}
	if t.stable == GameModeRace {
		if hasRaceExit || frameAt.Sub(t.candidateStartedAt) >= confirmAfter {
			t.stable = GameModeFreeRoam
			t.clearCandidate()
			return GameModeFreeRoam
		}
	} else if hasActivity || hasWorldPositionProgress(frame, t.candidateFrame) || frameAt.Sub(t.candidateStartedAt) >= confirmAfter {
		t.stable = GameModeFreeRoam
		t.clearCandidate()
		return GameModeFreeRoam
	}

	if t.stable == GameModeRace || t.stable == GameModeMenu {
		return t.stable
	}
	return GameModeUnknown
}

func (t *GameModeTracker) clearCandidate() {
	t.candidate = ""
	t.candidateStartedAt = time.Time{}
	t.candidateFrame = nil
}

func RawGameMode(frame NormalizedTelemetry) string {
	if !hasVehicleTelemetry(frame) {
		return GameModeMenu
	}
	if hasCompetitionTelemetry(frame) {
		return GameModeRace
	}
	return GameModeFreeRoam
}

func NormalizeGameMode(value string) string {
	switch strings.TrimSpace(value) {
	case GameModeMenu:
		return GameModeMenu
	case GameModeFreeRoam:
		return GameModeFreeRoam
	case GameModeRace:
		return GameModeRace
	case GameModeMixed:
		return GameModeMixed
	default:
		return GameModeUnknown
	}
}

func gameModeFrameTime(frame NormalizedTelemetry, receivedAt time.Time) time.Time {
	if frame.TimeMS > 0 {
		return time.Unix(0, 0).Add(time.Duration(frame.TimeMS) * time.Millisecond)
	}
	if !receivedAt.IsZero() {
		return receivedAt
	}
	if parsed, err := time.Parse(time.RFC3339Nano, frame.ReceivedAt); err == nil {
		return parsed
	}
	return time.Now()
}

func hasVehicleTelemetry(frame NormalizedTelemetry) bool {
	return frame.CarOrdinal > 0 || strings.TrimSpace(frame.CarClass) != "" || frame.CarPI > 0
}

func hasCompetitionTelemetry(frame NormalizedTelemetry) bool {
	return frame.RacePosition > 0 ||
		frame.LapNumber > 0 ||
		frame.CurrentLap > 0.05 ||
		frame.BestLap > 0.05 ||
		frame.LastLap > 0.05 ||
		(math.Abs(frame.DrivingLine01) > 0.05 && frame.SpeedKmh > 15)
}

func hasDrivingActivity(frame NormalizedTelemetry, previous *NormalizedTelemetry) bool {
	if frame.SpeedKmh > 5 || frame.Throttle01 > 0.05 || math.Abs(frame.Steer01) > 0.12 {
		return true
	}
	if previous == nil || frame.CarOrdinal != previous.CarOrdinal || frame.CarClass != previous.CarClass {
		return false
	}
	return worldDistance(frame, *previous) > 1.5 && frame.SpeedKmh > 1
}

func hasWorldPositionProgress(frame NormalizedTelemetry, candidateStart *NormalizedTelemetry) bool {
	if candidateStart == nil || frame.CarOrdinal != candidateStart.CarOrdinal || frame.CarClass != candidateStart.CarClass {
		return false
	}
	return worldDistance(frame, *candidateStart) > 6 && frame.SpeedKmh > 0.5
}

func hasRaceExitCue(frame NormalizedTelemetry, previous *NormalizedTelemetry) bool {
	if previous == nil || hasCompetitionTelemetry(frame) || math.Abs(frame.DrivingLine01) > 0.02 {
		return false
	}
	return previous.DistanceTraveled > 50 && frame.DistanceTraveled >= -1 && previous.DistanceTraveled-frame.DistanceTraveled > 50
}

func worldDistance(a NormalizedTelemetry, b NormalizedTelemetry) float64 {
	dx := a.PositionX - b.PositionX
	dy := a.PositionY - b.PositionY
	dz := a.PositionZ - b.PositionZ
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

func frameCopy(frame NormalizedTelemetry) *NormalizedTelemetry {
	out := frame
	return &out
}
