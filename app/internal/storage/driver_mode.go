package storage

import (
	"math"
	"strings"

	"fh6worker/internal/telemetry"
)

const (
	driverModeUnknown = "unknown"
	driverModePlayer  = "player"
	driverModeAuto    = "auto"
)

func DetectDriverMode(samples []telemetry.NormalizedTelemetry, events []telemetry.DetectedEvent, sessionGameMode string) DriverModeDetection {
	detection := DriverModeDetection{
		Mode:       driverModeUnknown,
		Confidence: 0.2,
		Summary:    "insufficient_samples",
		Evidence:   map[string]float64{},
	}
	if len(samples) == 0 {
		return detection
	}

	totalVehicleSamples := 0
	raceSamples := 0
	drivingLineSamples := 0
	highSlipSamples := 0
	steerDeltaSum := 0.0
	inputDeltaSum := 0.0
	prevSet := false
	prevSteer := 0.0
	prevThrottle := 0.0
	prevBrake := 0.0

	for _, sample := range samples {
		if !sampleLooksDrivable(sample) {
			continue
		}
		totalVehicleSamples++
		isRace := telemetry.NormalizeGameMode(sample.GameMode) == telemetry.GameModeRace || sample.IsRaceOn
		if isRace {
			raceSamples++
			if math.Abs(sample.DrivingLine01) > 0.03 || math.Abs(sample.AIBrakeDifference01) > 0.03 {
				drivingLineSamples++
			}
			if math.Max(math.Abs(sample.FrontCombinedSlipAvg), math.Abs(sample.RearCombinedSlipAvg)) >= 1.15 {
				highSlipSamples++
			}
			if prevSet {
				steerDeltaSum += math.Abs(sample.Steer01 - prevSteer)
				inputDeltaSum += math.Abs(sample.Throttle01-prevThrottle) + math.Abs(sample.Brake01-prevBrake)
			}
			prevSet = true
			prevSteer = sample.Steer01
			prevThrottle = sample.Throttle01
			prevBrake = sample.Brake01
		}
	}

	detection.Evidence["sample_count"] = float64(len(samples))
	detection.Evidence["vehicle_sample_count"] = float64(totalVehicleSamples)
	detection.Evidence["race_sample_count"] = float64(raceSamples)
	if totalVehicleSamples < 20 {
		return detection
	}

	raceRatio := float64(raceSamples) / float64(totalVehicleSamples)
	lineCoverage := safeDriverDiv(float64(drivingLineSamples), float64(raceSamples))
	slipRatio := safeDriverDiv(float64(highSlipSamples), float64(raceSamples))
	steerVolatility := safeDriverDiv(steerDeltaSum, float64(maxInt(raceSamples-1, 1)))
	inputVolatility := safeDriverDiv(inputDeltaSum, float64(maxInt(raceSamples-1, 1)))
	handlingEvents := driverHandlingEventCount(events)

	detection.Evidence["race_ratio"] = raceRatio
	detection.Evidence["driving_line_coverage"] = lineCoverage
	detection.Evidence["high_slip_ratio"] = slipRatio
	detection.Evidence["steer_volatility"] = steerVolatility
	detection.Evidence["input_volatility"] = inputVolatility
	detection.Evidence["handling_event_count"] = float64(handlingEvents)

	gameMode := telemetry.NormalizeGameMode(sessionGameMode)
	if gameMode != telemetry.GameModeRace && raceRatio < 0.6 {
		detection.Mode = driverModePlayer
		detection.Confidence = 0.55
		detection.Summary = "valid_free_roam_or_non_race_telemetry"
		return detection
	}

	if raceSamples >= 50 &&
		lineCoverage >= 0.65 &&
		slipRatio <= 0.015 &&
		handlingEvents <= 1 &&
		steerVolatility <= 0.08 &&
		inputVolatility <= 0.18 {
		detection.Mode = driverModeAuto
		detection.Confidence = clampDriverConfidence(0.72 + (lineCoverage-0.65)*0.25 - slipRatio*2 - float64(handlingEvents)*0.04)
		detection.Summary = "race_line_following_low_slip"
		return detection
	}

	detection.Mode = driverModePlayer
	detection.Confidence = 0.7
	detection.Summary = "valid_driver_telemetry_not_auto_like"
	if handlingEvents > 0 || slipRatio > 0.015 || steerVolatility > 0.08 {
		detection.Confidence = 0.82
		detection.Summary = "player_input_or_handling_events_detected"
	}
	return detection
}

func NormalizeDriverMode(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case driverModeAuto:
		return driverModeAuto
	case driverModePlayer:
		return driverModePlayer
	default:
		return driverModeUnknown
	}
}

func sampleLooksDrivable(sample telemetry.NormalizedTelemetry) bool {
	return sample.CarOrdinal > 0 && (sample.SpeedKmh > 3 || sample.Rpm > sample.EngineIdleRpm+100 || sample.Throttle01 > 0.05)
}

func driverHandlingEventCount(events []telemetry.DetectedEvent) int {
	count := 0
	for _, event := range events {
		switch strings.TrimSpace(event.Type) {
		case "corner_entry_understeer", "mid_corner_understeer", "power_understeer", "corner_exit_oversteer", "snap_oversteer", "high_speed_four_wheel_slide":
			count++
		}
	}
	return count
}

func safeDriverDiv(value float64, divisor float64) float64 {
	if divisor <= 0 {
		return 0
	}
	return value / divisor
}

func clampDriverConfidence(value float64) float64 {
	if value < 0.1 {
		return 0.1
	}
	if value > 0.98 {
		return 0.98
	}
	return value
}

func maxInt(left int, right int) int {
	if left > right {
		return left
	}
	return right
}
