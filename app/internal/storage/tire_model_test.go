package storage

import (
	"testing"

	"fh6worker/internal/telemetry"
)

func TestTireModelFrontLimited(t *testing.T) {
	samples := tireModelSamples(func(frame *telemetry.NormalizedTelemetry) {
		frame.Steer01 = 0.45
		frame.WheelFL.CombinedSlip = 1.08
		frame.WheelFR.CombinedSlip = 1.02
		frame.WheelFL.SlipAngle = 0.75
		frame.WheelFR.SlipAngle = 0.70
		frame.WheelRL.CombinedSlip = 0.35
		frame.WheelRR.CombinedSlip = 0.36
	})
	diag := BuildTireModelDiagnostic(samples, nil)
	if diag.LimitType != "front_limited" {
		t.Fatalf("limit type = %q, want front_limited: %#v", diag.LimitType, diag.Evidence)
	}
	if diag.FrontAxle.LimitScore <= diag.RearAxle.LimitScore {
		t.Fatalf("front score %.3f <= rear score %.3f", diag.FrontAxle.LimitScore, diag.RearAxle.LimitScore)
	}
	if diag.GripLimit.Type != "lateral_limit" || diag.GripLimit.LimitedAxle != "front" {
		t.Fatalf("grip limit = %#v, want front lateral limit", diag.GripLimit)
	}
}

func TestTireModelRearTractionLimited(t *testing.T) {
	samples := tireModelSamples(func(frame *telemetry.NormalizedTelemetry) {
		frame.Throttle01 = 0.82
		frame.Steer01 = 0.18
		frame.WheelRL.CombinedSlip = 0.95
		frame.WheelRR.CombinedSlip = 0.98
		frame.WheelRL.SlipRatio = 0.62
		frame.WheelRR.SlipRatio = 0.66
		frame.WheelFL.CombinedSlip = 0.33
		frame.WheelFR.CombinedSlip = 0.34
	})
	diag := BuildTireModelDiagnostic(samples, nil)
	if diag.LimitType != "traction_limited" || diag.Summary != "tire_model_rear_traction_limited" {
		t.Fatalf("limit/summary = %q/%q, want rear traction limit", diag.LimitType, diag.Summary)
	}
	if diag.GripLimit.Type != "traction_limit" || diag.GripLimit.LimitedAxle != "rear" {
		t.Fatalf("grip limit = %#v, want rear traction limit", diag.GripLimit)
	}
}

func TestTireModelHandbrakeMarksDynamicRearLimit(t *testing.T) {
	samples := tireModelSamples(func(frame *telemetry.NormalizedTelemetry) {
		frame.Throttle01 = 0
		frame.HandBrake01 = 0.85
		frame.Steer01 = 0.18
		frame.WheelRL.CombinedSlip = 0.95
		frame.WheelRR.CombinedSlip = 0.98
		frame.WheelRL.SlipAngle = 0.65
		frame.WheelRR.SlipAngle = 0.66
		frame.WheelFL.CombinedSlip = 0.30
		frame.WheelFR.CombinedSlip = 0.31
	})
	diag := BuildTireModelDiagnostic(samples, nil)
	if diag.Phase != "handbrake" {
		t.Fatalf("phase = %q, want handbrake", diag.Phase)
	}
	if diag.LimitType != "rear_limited" || diag.Summary != "tire_model_rear_handbrake_limited" {
		t.Fatalf("limit/summary = %q/%q, want handbrake rear limit", diag.LimitType, diag.Summary)
	}
}

func TestTireModelPhaseWindowConstant(t *testing.T) {
	if tireModelPhaseWindowMS != 800 {
		t.Fatalf("phase window = %d, want 800", tireModelPhaseWindowMS)
	}
}

func TestTireModelPhaseDetectsLaunch(t *testing.T) {
	samples := tirePhaseSamples(func(i int, frame *telemetry.NormalizedTelemetry) {
		frame.Brake01 = 0
		frame.Throttle01 = 0.86
		frame.Steer01 = 0.03
		frame.SpeedKmh = 8 + float64(i)*0.42
		frame.AccelerationZ = standardGravity * 0.45
	})
	diag := BuildTireModelDiagnostic(samples, nil)
	if diag.Phase != "launch" || diag.PhaseDetail.CurrentPhase != "launch" {
		t.Fatalf("phase = %q/%q, want launch: %#v", diag.Phase, diag.PhaseDetail.CurrentPhase, diag.PhaseDetail.Scores)
	}
}

func TestTireModelPhaseDetectsStraightBraking(t *testing.T) {
	samples := tirePhaseSamples(func(i int, frame *telemetry.NormalizedTelemetry) {
		frame.Brake01 = 0.82
		frame.Throttle01 = 0
		frame.Steer01 = 0.03
		frame.SpeedKmh = 150 - float64(i)*0.35
		frame.AccelerationZ = -standardGravity * 0.55
	})
	diag := BuildTireModelDiagnostic(samples, nil)
	if diag.Phase != "braking" || diag.PhaseDetail.CurrentPhase != "braking" {
		t.Fatalf("phase = %q/%q, want braking: %#v", diag.Phase, diag.PhaseDetail.CurrentPhase, diag.PhaseDetail.Scores)
	}
}

func TestTireModelPhaseDetectsCornerEntry(t *testing.T) {
	samples := tirePhaseSamples(func(i int, frame *telemetry.NormalizedTelemetry) {
		frame.Brake01 = 0.42
		frame.Throttle01 = 0.05
		frame.Steer01 = 0.16 + float64(i)*0.006
		frame.SpeedKmh = 135 - float64(i)*0.22
		frame.AccelerationX = standardGravity * 0.35
		frame.AccelerationZ = -standardGravity * 0.35
	})
	diag := BuildTireModelDiagnostic(samples, nil)
	if diag.Phase != "corner_entry" {
		t.Fatalf("phase = %q, want corner_entry: %#v", diag.Phase, diag.PhaseDetail)
	}
}

func TestTireModelPhaseDetectsSustainedCornering(t *testing.T) {
	samples := tirePhaseSamples(func(i int, frame *telemetry.NormalizedTelemetry) {
		frame.Brake01 = 0.03
		frame.Throttle01 = 0.24
		frame.Steer01 = 0.36
		frame.SpeedKmh = 115
		frame.AccelerationX = standardGravity * 0.72
		frame.AccelerationZ = 0
	})
	diag := BuildTireModelDiagnostic(samples, nil)
	if diag.Phase != "mid_speed_corner" {
		t.Fatalf("phase = %q, want mid_speed_corner: %#v", diag.Phase, diag.PhaseDetail)
	}
}

func TestTireModelPhaseDetectsCornerExit(t *testing.T) {
	samples := tirePhaseSamples(func(i int, frame *telemetry.NormalizedTelemetry) {
		frame.Brake01 = 0
		frame.Throttle01 = 0.56 + float64(i)*0.008
		frame.Steer01 = 0.30 - float64(i)*0.004
		frame.SpeedKmh = 75 + float64(i)*0.32
		frame.AccelerationX = standardGravity * 0.35
		frame.AccelerationZ = standardGravity * 0.25
	})
	diag := BuildTireModelDiagnostic(samples, nil)
	if diag.Phase != "corner_exit" {
		t.Fatalf("phase = %q, want corner_exit: %#v", diag.Phase, diag.PhaseDetail)
	}
}

func TestTireModelPhaseDetectsHighSpeedCorner(t *testing.T) {
	samples := tirePhaseSamples(func(i int, frame *telemetry.NormalizedTelemetry) {
		frame.Brake01 = 0
		frame.Throttle01 = 0.42
		frame.Steer01 = 0.30
		frame.SpeedKmh = 182
		frame.AccelerationX = standardGravity * 0.74
		frame.AccelerationZ = 0
	})
	diag := BuildTireModelDiagnostic(samples, nil)
	if diag.Phase != "high_speed_corner" {
		t.Fatalf("phase = %q, want high_speed_corner: %#v", diag.Phase, diag.PhaseDetail)
	}
}

func TestTireModelPhaseSpeedBandsUseDynamicReference(t *testing.T) {
	reference := make([]telemetry.NormalizedTelemetry, 0, 40)
	for i := 0; i < 40; i++ {
		frame := telemetry.NormalizedTelemetry{
			TimeMS:        int64(i * 100),
			GameMode:      telemetry.GameModeRace,
			SpeedKmh:      45 + float64(i)*4.0,
			Throttle01:    0.30,
			Brake01:       0,
			Steer01:       0.25,
			AccelerationX: standardGravity * 0.58,
			WheelFL:       baseTireWheel(),
			WheelFR:       baseTireWheel(),
			WheelRL:       baseTireWheel(),
			WheelRR:       baseTireWheel(),
		}
		reference = append(reference, frame)
	}
	cases := []struct {
		name  string
		ratio float64
		want  string
		band  float64
	}{
		{name: "low", ratio: 0.35, want: "low_speed_corner", band: 1},
		{name: "mid", ratio: 0.55, want: "mid_speed_corner", band: 2},
		{name: "high", ratio: 0.80, want: "high_speed_corner", band: 3},
	}
	referenceKmh, confidence := tirePhaseSpeedReference(reference)
	if confidence < 0.5 {
		t.Fatalf("reference confidence = %.1f, want high", confidence)
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			phase := make([]telemetry.NormalizedTelemetry, 0, 12)
			speed := referenceKmh * tc.ratio
			for i := 0; i < 12; i++ {
				frame := reference[len(reference)-12+i]
				frame.TimeMS = int64(i * 100)
				frame.SpeedKmh = speed
				frame.Throttle01 = 0.26
				frame.Brake01 = 0
				frame.Steer01 = 0.34
				frame.AccelerationX = standardGravity * 0.66
				frame.AccelerationZ = 0
				phase = append(phase, frame)
			}
			diag := BuildTirePhaseDiagnosticWithReference(phase, reference)
			if diag.CurrentPhase != tc.want {
				t.Fatalf("phase = %q, want %q: scores=%#v evidence=%#v", diag.CurrentPhase, tc.want, diag.Scores, diag.Evidence)
			}
			if diag.Evidence["speed_band"] != tc.band {
				t.Fatalf("speed band = %.0f, want %.0f: %#v", diag.Evidence["speed_band"], tc.band, diag.Evidence)
			}
		})
	}
}

func TestTireModelPhaseSpeedBandFallsBackWhenReferenceWeak(t *testing.T) {
	samples := tirePhaseSamples(func(i int, frame *telemetry.NormalizedTelemetry) {
		frame.Brake01 = 0
		frame.Throttle01 = 0.24
		frame.Steer01 = 0.34
		frame.SpeedKmh = 70
		frame.AccelerationX = standardGravity * 0.58
	})
	diag := BuildTirePhaseDiagnostic(samples)
	if diag.CurrentPhase != "low_speed_corner" {
		t.Fatalf("phase = %q, want low_speed_corner: %#v", diag.CurrentPhase, diag)
	}
	if diag.Evidence["speed_band_confidence"] != 0 {
		t.Fatalf("speed band confidence = %.1f, want low fallback", diag.Evidence["speed_band_confidence"])
	}
}

func TestTireModelPhaseDetectsStraightPower(t *testing.T) {
	samples := tirePhaseSamples(func(i int, frame *telemetry.NormalizedTelemetry) {
		frame.Brake01 = 0
		frame.Throttle01 = 0.88
		frame.Steer01 = 0.03
		frame.SpeedKmh = 70 + float64(i)*0.36
		frame.AccelerationZ = standardGravity * 0.42
	})
	diag := BuildTireModelDiagnostic(samples, nil)
	if diag.Phase != "straight_power" {
		t.Fatalf("phase = %q, want straight_power: %#v", diag.Phase, diag.PhaseDetail)
	}
}

func TestTireModelPhaseUsesShortWindowForRecentState(t *testing.T) {
	samples := make([]telemetry.NormalizedTelemetry, 0, 180)
	for i := 0; i < 180; i++ {
		frame := telemetry.NormalizedTelemetry{
			TimeMS:     int64(i * 100),
			GameMode:   telemetry.GameModeRace,
			SpeedKmh:   80 + float64(i)*0.10,
			Throttle01: 0.82,
			Brake01:    0,
			Steer01:    0.03,
			CarOrdinal: 1001,
			CarClass:   "A",
			CarPI:      800,
			Drivetrain: "RWD",
			WheelFL:    baseTireWheel(),
			WheelFR:    baseTireWheel(),
			WheelRL:    baseTireWheel(),
			WheelRR:    baseTireWheel(),
		}
		if i >= 155 {
			frame.Throttle01 = 0
			frame.Brake01 = 0.86
			frame.SpeedKmh = 120 - float64(i-155)*0.35
			frame.AccelerationZ = -standardGravity * 0.55
		}
		samples = append(samples, frame)
	}
	diag := BuildTireModelDiagnostic(samples, nil)
	if diag.Phase != "braking" {
		t.Fatalf("phase = %q, want recent braking despite long earlier acceleration: %#v", diag.Phase, diag.PhaseDetail)
	}
	if diag.PhaseDetail.WindowMS > tireModelPhaseWindowMS {
		t.Fatalf("phase window = %d, want <= %d", diag.PhaseDetail.WindowMS, tireModelPhaseWindowMS)
	}
}

func TestTireModelPhaseStabilityTransition(t *testing.T) {
	samples := tirePhaseSamples(func(i int, frame *telemetry.NormalizedTelemetry) {
		frame.Brake01 = 0
		frame.Throttle01 = 0.64
		frame.Steer01 = 0.20
		frame.SpeedKmh = 95 + float64(i)*0.12
		frame.AccelerationX = standardGravity * 0.30
		frame.AccelerationZ = standardGravity * 0.20
	})
	diag := BuildTireModelDiagnostic(samples, nil)
	if diag.PhaseDetail.PhaseStability == "stable" {
		t.Fatalf("phase stability = stable, want transition or low confidence: %#v", diag.PhaseDetail)
	}
}

func TestTireModelFourWheelLimited(t *testing.T) {
	samples := tireModelSamples(func(frame *telemetry.NormalizedTelemetry) {
		frame.Steer01 = 0.5
		frame.WheelFL.CombinedSlip = 1.02
		frame.WheelFR.CombinedSlip = 1.01
		frame.WheelRL.CombinedSlip = 0.98
		frame.WheelRR.CombinedSlip = 1.00
	})
	diag := BuildTireModelDiagnostic(samples, nil)
	if diag.LimitType != "four_wheel_limited" {
		t.Fatalf("limit type = %q, want four_wheel_limited", diag.LimitType)
	}
	if diag.GripLimit.Type != "combined_limit" {
		t.Fatalf("grip limit = %#v, want combined_limit", diag.GripLimit)
	}
}

func TestTireModelThermalRiskDoesNotBecomeLimit(t *testing.T) {
	samples := tireModelSamples(func(frame *telemetry.NormalizedTelemetry) {
		frame.WheelFL.TireTemp = 114
		frame.WheelFR.TireTemp = 113
		frame.WheelRL.TireTemp = 90
		frame.WheelRR.TireTemp = 91
	})
	diag := BuildTireModelDiagnostic(samples, nil)
	if diag.LimitType != "balanced" {
		t.Fatalf("limit type = %q, want balanced", diag.LimitType)
	}
	if diag.GripLimit.Type != "no_limit_detected" {
		t.Fatalf("grip limit = %#v, want no_limit_detected", diag.GripLimit)
	}
	if !containsString(diag.Warnings, "thermal_risk") {
		t.Fatalf("warnings = %#v, want thermal_risk", diag.Warnings)
	}
}

func TestTireModelPlatformRiskDoesNotBecomeLimit(t *testing.T) {
	samples := tireModelSamples(func(frame *telemetry.NormalizedTelemetry) {
		frame.WheelFL.SuspensionTravel = 0.97
		frame.WheelFR.SuspensionTravel = 0.96
		frame.WheelFL.CombinedSlip = 0.62
		frame.WheelFR.CombinedSlip = 0.63
	})
	diag := BuildTireModelDiagnostic(samples, nil)
	if diag.LimitType != "balanced" {
		t.Fatalf("limit type = %q, want balanced", diag.LimitType)
	}
	if diag.GripLimit.Type != "no_limit_detected" {
		t.Fatalf("grip limit = %#v, want no_limit_detected", diag.GripLimit)
	}
	if !containsString(diag.Warnings, "platform_risk") {
		t.Fatalf("warnings = %#v, want platform_risk", diag.Warnings)
	}
	if diag.Wheels[0].SuspensionOffsetPctMax < 96.9 || diag.Wheels[0].SuspensionOffsetPctMax > 97.1 {
		t.Fatalf("front-left suspension offset max = %.2f, want about 97%%", diag.Wheels[0].SuspensionOffsetPctMax)
	}
	if diag.FrontAxle.SuspensionOffsetPctMax < 96.9 || diag.FrontAxle.SuspensionOffsetPctMax > 97.1 {
		t.Fatalf("front axle suspension offset max = %.2f, want about 97%%", diag.FrontAxle.SuspensionOffsetPctMax)
	}
}

func TestTireModelNoDynamicLoad(t *testing.T) {
	samples := tireModelSamples(func(frame *telemetry.NormalizedTelemetry) {
		frame.SpeedKmh = 6
		frame.Throttle01 = 0
		frame.Brake01 = 0
		frame.Steer01 = 0
	})
	diag := BuildTireModelDiagnostic(samples, nil)
	if diag.LimitType != "no_dynamic_load" {
		t.Fatalf("limit type = %q, want no_dynamic_load", diag.LimitType)
	}
	if diag.DataQuality.Status == "valid" || diag.GripLimit.Type != "no_limit_detected" {
		t.Fatalf("quality/grip = %#v/%#v, want degraded no limit", diag.DataQuality, diag.GripLimit)
	}
}

func TestTireModelFlatSlipSignalDegradesQuality(t *testing.T) {
	samples := tireModelSamples(func(frame *telemetry.NormalizedTelemetry) {
		frame.SpeedKmh = 90
		frame.Throttle01 = 0.45
		frame.Steer01 = 0.20
		frame.AccelerationX = standardGravity * 0.25
		frame.WheelFL = telemetry.NormalizedWheelTelemetry{}
		frame.WheelFR = telemetry.NormalizedWheelTelemetry{}
		frame.WheelRL = telemetry.NormalizedWheelTelemetry{}
		frame.WheelRR = telemetry.NormalizedWheelTelemetry{}
	})
	diag := BuildTireModelDiagnostic(samples, nil)
	if diag.DataQuality.Status != "low_confidence" || diag.DataQuality.SlipSignal != "flat" {
		t.Fatalf("data quality = %#v, want low confidence flat slip", diag.DataQuality)
	}
	if diag.GripLimit.Type != "no_limit_detected" {
		t.Fatalf("grip limit = %#v, want no_limit_detected with flat slip signal", diag.GripLimit)
	}
}

func TestTireModelSingleSlipSpikeDoesNotTriggerLimit(t *testing.T) {
	samples := tireModelSamples(func(frame *telemetry.NormalizedTelemetry) {})
	samples[len(samples)/2].WheelFL.CombinedSlip = 1.25
	samples[len(samples)/2].WheelFR.CombinedSlip = 1.20
	diag := BuildTireModelDiagnostic(samples, nil)
	if diag.LimitType != "balanced" {
		t.Fatalf("limit type = %q, want balanced for single-frame spike", diag.LimitType)
	}
	if diag.FrontAxle.GripState != "stable" {
		t.Fatalf("front grip state = %q, want stable for single-frame spike", diag.FrontAxle.GripState)
	}
}

func TestTireModelStationaryDoesNotWarn(t *testing.T) {
	samples := tireModelSamples(func(frame *telemetry.NormalizedTelemetry) {
		frame.SpeedKmh = 0
		frame.Throttle01 = 0
		frame.Brake01 = 1
		frame.Steer01 = 0
		frame.WheelFL.TireTemp = 118
		frame.WheelFR.TireTemp = 118
		frame.WheelFL.SuspensionTravel = 0.99
		frame.WheelFR.SuspensionTravel = 0.99
	})
	diag := BuildTireModelDiagnostic(samples, nil)
	if diag.LimitType != "stationary" {
		t.Fatalf("limit type = %q, want stationary", diag.LimitType)
	}
	if diag.DataQuality.Status == "valid" || diag.GripLimit.Type != "no_limit_detected" {
		t.Fatalf("quality/grip = %#v/%#v, want stationary no strong limit", diag.DataQuality, diag.GripLimit)
	}
	if len(diag.Warnings) != 0 {
		t.Fatalf("warnings = %#v, want none while stationary", diag.Warnings)
	}
	if diag.FrontAxle.GripState != "stable" || diag.Wheels[0].GripState != "stable" {
		t.Fatalf("grip state front=%q wheel=%q, want stable", diag.FrontAxle.GripState, diag.Wheels[0].GripState)
	}
}

func TestTireModelLeftRightImbalanceWarning(t *testing.T) {
	samples := tireModelSamples(func(frame *telemetry.NormalizedTelemetry) {
		frame.WheelFL.CombinedSlip = 0.95
		frame.WheelRL.CombinedSlip = 0.88
		frame.WheelFR.CombinedSlip = 0.25
		frame.WheelRR.CombinedSlip = 0.22
	})
	diag := BuildTireModelDiagnostic(samples, nil)
	if diag.LeftRight.State != "imbalanced" || !containsString(diag.Warnings, "left_right_imbalance") {
		t.Fatalf("left/right = %#v warnings=%#v, want imbalance warning", diag.LeftRight, diag.Warnings)
	}
}

func TestTireModelGForceDiagnostic(t *testing.T) {
	samples := tireModelSamples(func(frame *telemetry.NormalizedTelemetry) {
		frame.AccelerationX = standardGravity
		frame.AccelerationY = standardGravity * 0.5
		frame.AccelerationZ = 0
	})
	samples[len(samples)-1].AccelerationX = standardGravity * 1.2
	diag := BuildTireModelDiagnostic(samples, nil)
	if diag.GForce.CurrentXG < 1.19 || diag.GForce.CurrentXG > 1.21 {
		t.Fatalf("current x g = %.3f, want about 1.2", diag.GForce.CurrentXG)
	}
	if diag.GForce.PeakAbsXG < 1.19 || diag.GForce.PeakAbsXG > 1.21 {
		t.Fatalf("peak x g = %.3f, want about 1.2", diag.GForce.PeakAbsXG)
	}
	if diag.GForce.DominantAxis != "x" {
		t.Fatalf("dominant axis = %q, want x", diag.GForce.DominantAxis)
	}
	if !containsString(diag.Warnings, "g_force_axis_mapping_unverified") {
		t.Fatalf("warnings = %#v, want axis mapping warning", diag.Warnings)
	}
}

func TestPowerToTireTractionOverPower(t *testing.T) {
	samples := powerToTireSamples("RWD", func(i int, frame *telemetry.NormalizedTelemetry) {
		frame.Throttle01 = 0.9
		frame.SpeedKmh = 80 + float64(i)*0.18
		frame.Power = 320000
		frame.Torque = 520
		frame.Rpm = 5200
		frame.RpmRatio = 0.72
		frame.WheelRL.SlipRatio = 0.62
		frame.WheelRR.SlipRatio = 0.66
	})
	diag := BuildTireModelDiagnostic(samples, nil)
	if diag.PowerToTire.Summary != "traction_over_power" {
		t.Fatalf("summary = %q, want traction_over_power: %#v", diag.PowerToTire.Summary, diag.PowerToTire.Evidence)
	}
	if diag.PowerToTire.DrivenAxle != "rear" || !diag.PowerToTire.TractionLimited {
		t.Fatalf("driven/limited = %q/%v, want rear traction limited", diag.PowerToTire.DrivenAxle, diag.PowerToTire.TractionLimited)
	}
}

func TestPowerToTireRPMBelowUsefulRange(t *testing.T) {
	samples := powerToTireSamples("RWD", func(i int, frame *telemetry.NormalizedTelemetry) {
		frame.Throttle01 = 0.88
		frame.SpeedKmh = 70 + float64(i)*0.10
		frame.Power = 90000
		frame.Torque = 260
		frame.Rpm = 2800
		frame.RpmRatio = 0.38
		frame.WheelRL.SlipRatio = 0.08
		frame.WheelRR.SlipRatio = 0.08
	})
	diag := BuildTireModelDiagnostic(samples, nil)
	if diag.PowerToTire.Summary != "rpm_below_useful_range" {
		t.Fatalf("summary = %q, want rpm_below_useful_range: %#v", diag.PowerToTire.Summary, diag.PowerToTire.Evidence)
	}
}

func TestPowerToTireRPMTooHighOrGearShort(t *testing.T) {
	samples := powerToTireSamples("RWD", func(i int, frame *telemetry.NormalizedTelemetry) {
		frame.Throttle01 = 0.9
		frame.SpeedKmh = 110 + float64(i)*0.14
		frame.Power = 210000
		frame.Torque = 330
		frame.Rpm = 7100
		frame.RpmRatio = 0.94
	})
	diag := BuildTireModelDiagnostic(samples, nil)
	if diag.PowerToTire.Summary != "rpm_too_high_or_gear_short" {
		t.Fatalf("summary = %q, want rpm_too_high_or_gear_short: %#v", diag.PowerToTire.Summary, diag.PowerToTire.Evidence)
	}
}

func TestPowerToTireZeroPowerSignalDoesNotMisreport(t *testing.T) {
	samples := powerToTireSamples("RWD", func(i int, frame *telemetry.NormalizedTelemetry) {
		frame.Throttle01 = 0.9
		frame.SpeedKmh = 90 + float64(i)*0.1
		frame.Power = 0
		frame.Torque = 0
		frame.Rpm = 5200
		frame.RpmRatio = 0.7
		frame.WheelRL.SlipRatio = 0.6
		frame.WheelRR.SlipRatio = 0.6
	})
	diag := BuildTireModelDiagnostic(samples, nil)
	if diag.PowerToTire.Summary != "power_to_tire_power_signal_unavailable" {
		t.Fatalf("summary = %q, want power signal unavailable", diag.PowerToTire.Summary)
	}
	if diag.PowerToTire.TractionLimited {
		t.Fatalf("traction limited = true, want false without power signal")
	}
}

func TestPowerToTireDrivenAxleByDrivetrain(t *testing.T) {
	cases := []struct {
		drivetrain string
		wantAxle   string
		frontSlip  float64
		rearSlip   float64
		wantSlip   float64
	}{
		{drivetrain: "FWD", wantAxle: "front", frontSlip: 0.62, rearSlip: 0.08, wantSlip: 0.62},
		{drivetrain: "RWD", wantAxle: "rear", frontSlip: 0.08, rearSlip: 0.62, wantSlip: 0.62},
		{drivetrain: "AWD", wantAxle: "all", frontSlip: 0.42, rearSlip: 0.62, wantSlip: 0.62},
	}
	for _, tc := range cases {
		samples := powerToTireSamples(tc.drivetrain, func(i int, frame *telemetry.NormalizedTelemetry) {
			frame.Throttle01 = 0.9
			frame.SpeedKmh = 80 + float64(i)*0.18
			frame.Power = 260000
			frame.Torque = 440
			frame.RpmRatio = 0.7
			frame.WheelFL.SlipRatio = tc.frontSlip
			frame.WheelFR.SlipRatio = tc.frontSlip
			frame.WheelRL.SlipRatio = tc.rearSlip
			frame.WheelRR.SlipRatio = tc.rearSlip
		})
		diag := BuildTireModelDiagnostic(samples, nil).PowerToTire
		if diag.DrivenAxle != tc.wantAxle {
			t.Fatalf("%s driven axle = %q, want %q", tc.drivetrain, diag.DrivenAxle, tc.wantAxle)
		}
		if diag.DrivenSlipRatioP90 < tc.wantSlip-0.01 || diag.DrivenSlipRatioP90 > tc.wantSlip+0.01 {
			t.Fatalf("%s driven slip p90 = %.2f, want %.2f", tc.drivetrain, diag.DrivenSlipRatioP90, tc.wantSlip)
		}
	}
}

func TestBrakeToTireFrontBrakeLockTendency(t *testing.T) {
	samples := brakeToTireSamples(func(i int, frame *telemetry.NormalizedTelemetry) {
		frame.Brake01 = 0.82
		frame.SpeedKmh = 135 - float64(i)*0.20
		frame.WheelFL.SlipRatio = 0.58
		frame.WheelFR.SlipRatio = 0.56
		frame.WheelRL.SlipRatio = 0.08
		frame.WheelRR.SlipRatio = 0.08
	})
	model := BuildTireModelDiagnostic(samples, nil)
	diag := model.BrakeToTire
	if diag.Summary != "front_brake_lock_tendency" {
		t.Fatalf("summary = %q, want front_brake_lock_tendency: %#v", diag.Summary, diag.Evidence)
	}
	if diag.FrontSlipRatioP90 <= diag.RearSlipRatioP90 {
		t.Fatalf("front/rear slip = %.2f/%.2f, want front higher", diag.FrontSlipRatioP90, diag.RearSlipRatioP90)
	}
	if model.GripLimit.Type != "braking_limit" || model.GripLimit.LimitedAxle != "front" {
		t.Fatalf("grip limit = %#v, want front braking limit", model.GripLimit)
	}
}

func TestBrakeToTireRearBrakeLockTendency(t *testing.T) {
	samples := brakeToTireSamples(func(i int, frame *telemetry.NormalizedTelemetry) {
		frame.Brake01 = 0.80
		frame.SpeedKmh = 130 - float64(i)*0.22
		frame.WheelFL.SlipRatio = 0.08
		frame.WheelFR.SlipRatio = 0.08
		frame.WheelRL.SlipRatio = 0.56
		frame.WheelRR.SlipRatio = 0.58
	})
	diag := BuildTireModelDiagnostic(samples, nil).BrakeToTire
	if diag.Summary != "rear_brake_lock_tendency" {
		t.Fatalf("summary = %q, want rear_brake_lock_tendency: %#v", diag.Summary, diag.Evidence)
	}
}

func TestBrakeToTireTrailBrakeFrontOverload(t *testing.T) {
	samples := brakeToTireSamples(func(i int, frame *telemetry.NormalizedTelemetry) {
		frame.Brake01 = 0.62
		frame.Steer01 = 0.36
		frame.SpeedKmh = 125 - float64(i)*0.18
		frame.WheelFL.CombinedSlip = 0.95
		frame.WheelFR.CombinedSlip = 0.92
		frame.WheelRL.CombinedSlip = 0.35
		frame.WheelRR.CombinedSlip = 0.35
		frame.WheelFL.SlipRatio = 0.16
		frame.WheelFR.SlipRatio = 0.16
	})
	diag := BuildTireModelDiagnostic(samples, nil).BrakeToTire
	if diag.Summary != "trail_brake_front_overload" {
		t.Fatalf("summary = %q, want trail_brake_front_overload: %#v", diag.Summary, diag.Evidence)
	}
	if !diag.TrailBraking {
		t.Fatalf("trail braking = false, want true")
	}
}

func TestBrakeToTireHandbrakeRearSlide(t *testing.T) {
	samples := brakeToTireSamples(func(i int, frame *telemetry.NormalizedTelemetry) {
		frame.Brake01 = 0
		frame.HandBrake01 = 0.85
		frame.SpeedKmh = 90 - float64(i)*0.16
		frame.WheelRL.SlipRatio = 0.62
		frame.WheelRR.SlipRatio = 0.64
	})
	diag := BuildTireModelDiagnostic(samples, nil).BrakeToTire
	if diag.Summary != "handbrake_rear_slide" {
		t.Fatalf("summary = %q, want handbrake_rear_slide: %#v", diag.Summary, diag.Evidence)
	}
	if !diag.HandbrakeActive {
		t.Fatalf("handbrake active = false, want true")
	}
}

func TestBrakeToTireNotSlowingEffectively(t *testing.T) {
	samples := brakeToTireSamples(func(i int, frame *telemetry.NormalizedTelemetry) {
		frame.Brake01 = 0.78
		frame.SpeedKmh = 130 - float64(i)*0.02
		frame.WheelFL.SlipRatio = 0.12
		frame.WheelFR.SlipRatio = 0.12
		frame.WheelRL.SlipRatio = 0.12
		frame.WheelRR.SlipRatio = 0.12
	})
	diag := BuildTireModelDiagnostic(samples, nil).BrakeToTire
	if diag.Summary != "brake_not_slowing_effectively" {
		t.Fatalf("summary = %q, want brake_not_slowing_effectively: %#v", diag.Summary, diag.Evidence)
	}
}

func TestBrakeToTireLowBrakeAndInsufficientSamples(t *testing.T) {
	lowBrake := brakeToTireSamples(func(i int, frame *telemetry.NormalizedTelemetry) {
		frame.Brake01 = 0.05
		frame.SpeedKmh = 120
	})
	lowDiag := BuildTireModelDiagnostic(lowBrake, nil).BrakeToTire
	if lowDiag.Summary != "brake_to_tire_low_brake" {
		t.Fatalf("low brake summary = %q, want brake_to_tire_low_brake", lowDiag.Summary)
	}

	insufficient := brakeToTireSamples(func(i int, frame *telemetry.NormalizedTelemetry) {
		frame.Brake01 = 0
		if i < 3 {
			frame.Brake01 = 0.8
			frame.SpeedKmh = 120 - float64(i)*0.3
		}
	})
	insufficientDiag := BuildTireModelDiagnostic(insufficient, nil).BrakeToTire
	if insufficientDiag.Summary != "brake_to_tire_insufficient" {
		t.Fatalf("insufficient summary = %q, want brake_to_tire_insufficient", insufficientDiag.Summary)
	}
}

func TestTireModelCamberInferenceFrontNeedsMoreNegative(t *testing.T) {
	samples := tireModelSamples(func(frame *telemetry.NormalizedTelemetry) {
		frame.SpeedKmh = 120
		frame.Steer01 = 0.42
		frame.Throttle01 = 0.35
		frame.WheelFL.SlipAngle = 0.68
		frame.WheelFR.SlipAngle = 0.64
		frame.WheelFL.CombinedSlip = 0.78
		frame.WheelFR.CombinedSlip = 0.76
		frame.WheelRL.SlipAngle = 0.24
		frame.WheelRR.SlipAngle = 0.22
		frame.WheelRL.CombinedSlip = 0.38
		frame.WheelRR.CombinedSlip = 0.36
		frame.AccelerationX = standardGravity * 0.2
		frame.AccelerationY = standardGravity * 0.85
	})
	diag := BuildTireModelDiagnostic(samples, nil)
	if diag.Camber.Status != "ready" {
		t.Fatalf("camber status = %q, want ready", diag.Camber.Status)
	}
	if diag.Camber.FrontState != "likely_needs_more_negative" {
		t.Fatalf("front camber state = %q, want likely_needs_more_negative", diag.Camber.FrontState)
	}
	if diag.Camber.Summary != "camber_inference_front_needs_more_negative" {
		t.Fatalf("camber summary = %q", diag.Camber.Summary)
	}
}

func TestTireModelCamberInferenceInsufficientWithoutCornering(t *testing.T) {
	samples := tireModelSamples(func(frame *telemetry.NormalizedTelemetry) {
		frame.Steer01 = 0.04
		frame.SpeedKmh = 120
	})
	diag := BuildTireModelDiagnostic(samples, nil)
	if diag.Camber.Status != "insufficient_data" {
		t.Fatalf("camber status = %q, want insufficient_data", diag.Camber.Status)
	}
	if !containsString(diag.Camber.Warnings, "camber_inference_no_three_point_temps") {
		t.Fatalf("camber warnings = %#v", diag.Camber.Warnings)
	}
}

func TestTireModelPhaseDetectsStraightDecelAndLightBraking(t *testing.T) {
	decel := tirePhaseSamples(func(i int, frame *telemetry.NormalizedTelemetry) {
		frame.Throttle01 = 0.04
		frame.Brake01 = 0.02
		frame.Steer01 = 0.03
		frame.SpeedKmh = 130 - float64(i)*0.30
		frame.AccelerationZ = -standardGravity * 0.28
	})
	decelDiag := BuildTireModelDiagnostic(decel, nil)
	if decelDiag.Phase != "straight_decel" {
		t.Fatalf("phase = %q, want straight_decel: %#v", decelDiag.Phase, decelDiag.PhaseDetail.Scores)
	}

	lightBrake := tirePhaseSamples(func(i int, frame *telemetry.NormalizedTelemetry) {
		frame.Throttle01 = 0
		frame.Brake01 = 0.22
		frame.Steer01 = 0.04
		frame.SpeedKmh = 130 - float64(i)*0.24
		frame.AccelerationZ = -standardGravity * 0.25
	})
	lightBrakeDiag := BuildTireModelDiagnostic(lightBrake, nil)
	if lightBrakeDiag.Phase != "light_braking" {
		t.Fatalf("phase = %q, want light_braking: %#v", lightBrakeDiag.Phase, lightBrakeDiag.PhaseDetail.Scores)
	}
}

func TestTireModelPhaseDetectsDriftSources(t *testing.T) {
	handbrake := tirePhaseSamples(func(i int, frame *telemetry.NormalizedTelemetry) {
		frame.HandBrake01 = 0.75
		frame.Throttle01 = 0.05
		frame.Steer01 = 0.36
		frame.SpeedKmh = 70
		frame.WheelRL.CombinedSlip = 0.92
		frame.WheelRR.CombinedSlip = 0.94
		frame.WheelRL.SlipAngle = 0.68
		frame.WheelRR.SlipAngle = 0.70
	})
	handbrakeDiag := BuildTireModelDiagnostic(handbrake, nil)
	if handbrakeDiag.Phase != "handbrake" && handbrakeDiag.Phase != "drift" {
		t.Fatalf("phase = %q, want handbrake or drift: %#v", handbrakeDiag.Phase, handbrakeDiag.PhaseDetail.Scores)
	}
	if got := tireDriftSource(handbrakeDiag.PhaseDetail.Evidence); got != "handbrake_initiated" {
		t.Fatalf("drift source = %q, want handbrake_initiated", got)
	}

	power := tirePhaseSamples(func(i int, frame *telemetry.NormalizedTelemetry) {
		frame.Throttle01 = 0.82
		frame.Steer01 = 0.34
		frame.SpeedKmh = 80 + float64(i)*0.12
		frame.WheelRL.CombinedSlip = 0.94
		frame.WheelRR.CombinedSlip = 0.96
		frame.WheelRL.SlipRatio = 0.52
		frame.WheelRR.SlipRatio = 0.54
		frame.WheelFL.CombinedSlip = 0.35
		frame.WheelFR.CombinedSlip = 0.36
	})
	powerDiag := BuildTireModelDiagnostic(power, nil)
	if powerDiag.Phase != "drift" {
		t.Fatalf("phase = %q, want drift: %#v", powerDiag.Phase, powerDiag.PhaseDetail.Scores)
	}
	if got := tireDriftSource(powerDiag.PhaseDetail.Evidence); got != "power_oversteer" {
		t.Fatalf("drift source = %q, want power_oversteer", got)
	}
}

func TestTireIssueAnalysisAggregatesTractionSegments(t *testing.T) {
	samples := tireModelSamples(func(frame *telemetry.NormalizedTelemetry) {
		frame.Throttle01 = 0.84
		frame.Steer01 = 0.18
		frame.WheelRL.CombinedSlip = 0.95
		frame.WheelRR.CombinedSlip = 0.98
		frame.WheelRL.SlipRatio = 0.62
		frame.WheelRR.SlipRatio = 0.66
		frame.WheelFL.CombinedSlip = 0.33
		frame.WheelFR.CombinedSlip = 0.34
	})
	analysis := BuildTireIssueAnalysis(samples)
	if len(analysis.Groups) == 0 {
		t.Fatalf("groups empty, want traction issue: %#v", analysis)
	}
	group := analysis.Groups[0]
	if group.Type != "traction_limit" || group.LimitedAxle != "rear" {
		t.Fatalf("group = %#v, want rear traction limit", group)
	}
	if group.Count < 1 || group.TotalDurationMS <= 0 {
		t.Fatalf("group count/duration invalid: %#v", group)
	}
}

func TestTireIssueAnalysisIgnoresSingleSlipSpike(t *testing.T) {
	samples := tireModelSamples(func(frame *telemetry.NormalizedTelemetry) {})
	samples[len(samples)/2].WheelFL.CombinedSlip = 1.35
	samples[len(samples)/2].WheelFR.CombinedSlip = 1.30
	analysis := BuildTireIssueAnalysis(samples)
	for _, group := range analysis.Groups {
		if group.Type == "lateral_limit" || group.Type == "combined_limit" {
			t.Fatalf("unexpected limit group for single-frame spike: %#v", group)
		}
	}
}

func TestTireIssueAnalysisThermalRiskIsRiskGroup(t *testing.T) {
	samples := tireModelSamples(func(frame *telemetry.NormalizedTelemetry) {
		frame.WheelFL.TireTemp = 114
		frame.WheelFR.TireTemp = 113
		frame.WheelRL.TireTemp = 90
		frame.WheelRR.TireTemp = 91
	})
	analysis := BuildTireIssueAnalysis(samples)
	if len(analysis.Groups) == 0 || analysis.Groups[0].Type != "thermal_risk" {
		t.Fatalf("groups = %#v, want thermal risk group", analysis.Groups)
	}
	if analysis.Groups[0].LimitType != "risk" {
		t.Fatalf("limit type = %q, want risk", analysis.Groups[0].LimitType)
	}
}

func tireModelSamples(mutator func(*telemetry.NormalizedTelemetry)) []telemetry.NormalizedTelemetry {
	samples := make([]telemetry.NormalizedTelemetry, 0, 40)
	for i := 0; i < 40; i++ {
		frame := telemetry.NormalizedTelemetry{
			TimeMS:     int64(i * 100),
			GameMode:   telemetry.GameModeRace,
			SpeedKmh:   120,
			Throttle01: 0.2,
			Brake01:    0,
			Steer01:    0.1,
			CarOrdinal: 1001,
			CarClass:   "A",
			CarPI:      800,
			Drivetrain: "RWD",
			WheelFL:    baseTireWheel(),
			WheelFR:    baseTireWheel(),
			WheelRL:    baseTireWheel(),
			WheelRR:    baseTireWheel(),
		}
		mutator(&frame)
		samples = append(samples, frame)
	}
	return samples
}

func baseTireWheel() telemetry.NormalizedWheelTelemetry {
	return telemetry.NormalizedWheelTelemetry{
		CombinedSlip:     0.25,
		SlipRatio:        0.12,
		SlipAngle:        0.16,
		TireTemp:         88,
		SuspensionTravel: 0.42,
	}
}

func tirePhaseSamples(mutator func(int, *telemetry.NormalizedTelemetry)) []telemetry.NormalizedTelemetry {
	return tireModelSamples(func(frame *telemetry.NormalizedTelemetry) {
		index := int(frame.TimeMS / 100)
		frame.Throttle01 = 0.2
		frame.Brake01 = 0
		frame.HandBrake01 = 0
		frame.Steer01 = 0.05
		frame.SpeedKmh = 100
		frame.AccelerationX = 0
		frame.AccelerationZ = 0
		mutator(index, frame)
	})
}

func powerToTireSamples(drivetrain string, mutator func(int, *telemetry.NormalizedTelemetry)) []telemetry.NormalizedTelemetry {
	return tireModelSamples(func(frame *telemetry.NormalizedTelemetry) {
		index := int(frame.TimeMS / 100)
		frame.Drivetrain = drivetrain
		frame.Gear = 3
		frame.Rpm = 5200
		frame.RpmRatio = 0.7
		frame.EngineMaxRpm = 7600
		frame.Power = 220000
		frame.Torque = 420
		frame.Throttle01 = 0.8
		frame.Steer01 = 0.02
		frame.WheelFL.SlipRatio = 0.08
		frame.WheelFR.SlipRatio = 0.08
		frame.WheelRL.SlipRatio = 0.08
		frame.WheelRR.SlipRatio = 0.08
		mutator(index, frame)
	})
}

func brakeToTireSamples(mutator func(int, *telemetry.NormalizedTelemetry)) []telemetry.NormalizedTelemetry {
	return tireModelSamples(func(frame *telemetry.NormalizedTelemetry) {
		index := int(frame.TimeMS / 100)
		frame.Throttle01 = 0
		frame.Brake01 = 0.7
		frame.HandBrake01 = 0
		frame.Steer01 = 0.04
		frame.AccelerationX = 0
		frame.AccelerationZ = -standardGravity * 0.35
		frame.WheelFL.SlipRatio = 0.08
		frame.WheelFR.SlipRatio = 0.08
		frame.WheelRL.SlipRatio = 0.08
		frame.WheelRR.SlipRatio = 0.08
		frame.WheelFL.CombinedSlip = 0.28
		frame.WheelFR.CombinedSlip = 0.28
		frame.WheelRL.CombinedSlip = 0.28
		frame.WheelRR.CombinedSlip = 0.28
		mutator(index, frame)
	})
}
