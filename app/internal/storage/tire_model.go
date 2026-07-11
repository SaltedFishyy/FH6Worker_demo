package storage

import (
	"math"
	"sort"
	"strings"

	"fh6worker/internal/telemetry"
)

const (
	tireModelStatusNoData = "no_data"
	tireModelStatusReady  = "ready"

	tireModelWindowMS      = int64(15000)
	tireModelPhaseWindowMS = int64(800)
	tireModelTrendWindowMS = int64(3000)
	tireModelMinSamples    = 10
	tireModelMaxSamples    = 250
	tireModelMinDynamic    = 8
	tireModelSlipWarn      = 0.75
	tireModelSlipLimit     = 1.00
	tireModelSlipHigh      = 0.25
	tireModelRatioWarn     = 0.35
	tireModelDynamicKmh    = 15.0
	tireModelAxleDelta     = 0.15
	tireModelSideDelta     = 0.25
	tireModelHotTemp       = 110.0
	tireModelTempDelta     = 15.0
	tireModelBottomOut     = 0.95
	tireModelParkedKmh     = 5.0
	standardGravity        = 9.80665
	camberMinSamples       = 8

	powerToTireHighThrottle = 0.65
	powerToTireMinSamples   = 6
	powerToTireLowRPM       = 0.45
	powerToTireHighRPM      = 0.90
	powerToTireSlipWarn     = 0.35
	powerToTireSlipLimit    = 0.45

	brakeToTireBrakeSampleMin  = 0.15
	brakeToTireBrakeMin        = 0.35
	brakeToTireHandbrakeSample = 0.10
	brakeToTireHandbrakeMin    = 0.20
	brakeToTireMinSamples      = 6
	brakeToTireSlipWarn        = 0.35
	brakeToTireCombinedWarn    = 0.75
	brakeToTireAxleDelta       = 0.12
	brakeToTireSteerMin        = 0.14
	brakeToTireWeakDecelG      = 0.18
)

func BuildTireModelDiagnostic(samples []telemetry.NormalizedTelemetry, current *telemetry.NormalizedTelemetry) TireModelDiagnostic {
	diag := TireModelDiagnostic{
		Status:     tireModelStatusNoData,
		UpdatedAt:  nowText(),
		Confidence: quickConfidenceLow,
		Phase:      "unknown",
		PhaseDetail: TirePhaseDiagnostic{
			CurrentPhase:   "unknown",
			SecondaryPhase: "unknown",
			StablePhase:    "unknown",
			PhaseStability: "low_confidence",
			Confidence:     quickConfidenceLow,
			Scores:         map[string]float64{},
			Evidence:       map[string]float64{},
		},
		DataQuality:   defaultTireDataQuality(),
		GripLimit:     defaultTireGripLimit(),
		LimitType:     "unknown",
		Warnings:      []string{},
		IssueAnalysis: defaultTireIssueAnalysis(),
		IssueAdvice:   defaultTireIssueAdvice(),
		Hints:         []TireModelHint{},
		Evidence:      map[string]float64{},
		Vehicle:       quickVehicleSnapshot(samples, current),
	}
	if len(samples) == 0 {
		diag.Summary = "tire_model_no_data"
		diag.Explanation = "tire_model_waiting_for_samples"
		diag.Warnings = append(diag.Warnings, "tire_model_no_data")
		return diag
	}

	ordered := sortedTelemetrySamples(samples)
	window := tireModelWindow(ordered)
	if len(window) == 0 {
		window = ordered
	}
	diag.Status = tireModelStatusReady
	diag.SampleCount = len(window)
	if len(window) >= 2 {
		diag.WindowMS = window[len(window)-1].TimeMS - window[0].TimeMS
		if diag.WindowMS < 0 {
			diag.WindowMS = 0
		}
	}
	diag.GameMode = quickSummarizeGameMode(window)
	diag.PhaseDetail = BuildTirePhaseDiagnosticWithReference(tireModelPhaseWindow(ordered), window)
	diag.Phase = diag.PhaseDetail.CurrentPhase
	diag.Confidence = tireModelConfidence(len(window), diag.WindowMS)
	if len(window) < tireModelMinSamples {
		diag.Warnings = append(diag.Warnings, "sample_insufficient")
	}

	fl := buildTireWheel("front_left", window, func(frame telemetry.NormalizedTelemetry) telemetry.NormalizedWheelTelemetry { return frame.WheelFL })
	fr := buildTireWheel("front_right", window, func(frame telemetry.NormalizedTelemetry) telemetry.NormalizedWheelTelemetry { return frame.WheelFR })
	rl := buildTireWheel("rear_left", window, func(frame telemetry.NormalizedTelemetry) telemetry.NormalizedWheelTelemetry { return frame.WheelRL })
	rr := buildTireWheel("rear_right", window, func(frame telemetry.NormalizedTelemetry) telemetry.NormalizedWheelTelemetry { return frame.WheelRR })
	diag.Wheels = []TireWheelDiagnostic{fl, fr, rl, rr}
	diag.FrontAxle = buildTireAxle("front", fl, fr)
	diag.RearAxle = buildTireAxle("rear", rl, rr)
	diag.LeftRight = buildTireSideBalance(fl, fr, rl, rr)
	diag.GForce = buildGForceDiagnostic(window)
	diag.Camber = buildCamberInference(window, diag.GForce)
	diag.PowerToTire = BuildPowerToTireDiagnostic(window, current)
	diag.BrakeToTire = BuildBrakeToTireDiagnostic(window)

	avgThrottle, avgBrake, avgHandBrake, avgSteer, avgSpeed := tireModelInputs(window)
	diag.Evidence = map[string]float64{
		"front_combined_slip_avg":  diag.FrontAxle.CombinedSlipAvg,
		"rear_combined_slip_avg":   diag.RearAxle.CombinedSlipAvg,
		"front_combined_slip_max":  diag.FrontAxle.CombinedSlipMax,
		"rear_combined_slip_max":   diag.RearAxle.CombinedSlipMax,
		"front_combined_slip_p90":  diag.FrontAxle.CombinedSlipP90,
		"rear_combined_slip_p90":   diag.RearAxle.CombinedSlipP90,
		"front_combined_high_pct":  diag.FrontAxle.CombinedSlipHighPct,
		"rear_combined_high_pct":   diag.RearAxle.CombinedSlipHighPct,
		"front_slip_ratio_avg":     diag.FrontAxle.SlipRatioAvg,
		"rear_slip_ratio_avg":      diag.RearAxle.SlipRatioAvg,
		"front_slip_ratio_p90":     diag.FrontAxle.SlipRatioP90,
		"rear_slip_ratio_p90":      diag.RearAxle.SlipRatioP90,
		"front_slip_angle_avg":     diag.FrontAxle.SlipAngleAvg,
		"rear_slip_angle_avg":      diag.RearAxle.SlipAngleAvg,
		"front_slip_angle_p90":     diag.FrontAxle.SlipAngleP90,
		"rear_slip_angle_p90":      diag.RearAxle.SlipAngleP90,
		"front_tire_temp_avg":      diag.FrontAxle.TireTempAvg,
		"rear_tire_temp_avg":       diag.RearAxle.TireTempAvg,
		"front_suspension_max":     diag.FrontAxle.SuspensionTravelMax,
		"rear_suspension_max":      diag.RearAxle.SuspensionTravelMax,
		"left_right_slip_delta":    diag.LeftRight.Delta,
		"avg_throttle":             avgThrottle,
		"avg_brake":                avgBrake,
		"avg_handbrake":            avgHandBrake,
		"avg_steer":                avgSteer,
		"avg_speed_kmh":            avgSpeed,
		"current_total_g":          diag.GForce.CurrentTotalG,
		"peak_total_g":             diag.GForce.PeakTotalG,
		"peak_abs_x_g":             diag.GForce.PeakAbsXG,
		"peak_abs_y_g":             diag.GForce.PeakAbsYG,
		"peak_abs_z_g":             diag.GForce.PeakAbsZG,
		"camber_cornering_samples": diag.Camber.Evidence["cornering_samples"],
		"camber_front_slip_angle":  diag.Camber.Evidence["front_slip_angle_avg"],
		"camber_rear_slip_angle":   diag.Camber.Evidence["rear_slip_angle_avg"],
	}
	dynamicWindow := tireModelDynamicSamples(window)
	diag.Evidence["dynamic_sample_count"] = float64(len(dynamicWindow))
	diag.DataQuality = buildTireDataQuality(window, dynamicWindow, current, diag)
	if tireModelCurrentStationary(current, avgSpeed, avgThrottle) {
		markTireModelStationary(&diag)
		diag.IssueAnalysis = BuildTireIssueAnalysis(ordered)
		diag.IssueAdvice = BuildTireIssueAdviceFromAnalysis(diag.IssueAnalysis)
		return diag
	}
	appendTireModelRisks(&diag)
	if len(dynamicWindow) < tireModelMinDynamic {
		markTireModelNoDynamicLoad(&diag)
		diag.IssueAnalysis = BuildTireIssueAnalysis(ordered)
		diag.IssueAdvice = BuildTireIssueAdviceFromAnalysis(diag.IssueAnalysis)
		return diag
	}
	limitDiag := buildTireLimitDiagnostic(dynamicWindow)
	diag.LimitType, diag.Summary, diag.Explanation = classifyTireLimit(limitDiag, avgThrottle, avgBrake, avgHandBrake)
	diag.GripLimit = buildTireGripLimit(limitDiag, diag.LeftRight, avgThrottle, avgBrake, avgHandBrake, diag.DataQuality)
	diag.Hints = tireModelHints(diag)
	diag.IssueAnalysis = BuildTireIssueAnalysis(ordered)
	diag.IssueAdvice = BuildTireIssueAdviceFromAnalysis(diag.IssueAnalysis)
	return diag
}

func buildGForceDiagnostic(samples []telemetry.NormalizedTelemetry) GForceDiagnostic {
	diag := GForceDiagnostic{
		Source:      "AccelerationX/AccelerationY/AccelerationZ",
		AxisMapping: "raw_packet_axes_unverified",
	}
	if len(samples) == 0 {
		return diag
	}
	var sumAbsX, sumAbsY, sumAbsZ, sumTotal float64
	for _, sample := range samples {
		x := sample.AccelerationX / standardGravity
		y := sample.AccelerationY / standardGravity
		z := sample.AccelerationZ / standardGravity
		total := math.Sqrt(x*x + y*y + z*z)
		diag.Series = append(diag.Series, GForcePoint{
			TimeMS: sample.TimeMS,
			XG:     x,
			YG:     y,
			ZG:     z,
			TotalG: total,
		})
		sumAbsX += math.Abs(x)
		sumAbsY += math.Abs(y)
		sumAbsZ += math.Abs(z)
		sumTotal += total
		diag.PeakAbsXG = math.Max(diag.PeakAbsXG, math.Abs(x))
		diag.PeakAbsYG = math.Max(diag.PeakAbsYG, math.Abs(y))
		diag.PeakAbsZG = math.Max(diag.PeakAbsZG, math.Abs(z))
		diag.PeakTotalG = math.Max(diag.PeakTotalG, total)
	}
	last := samples[len(samples)-1]
	diag.CurrentXG = last.AccelerationX / standardGravity
	diag.CurrentYG = last.AccelerationY / standardGravity
	diag.CurrentZG = last.AccelerationZ / standardGravity
	diag.CurrentTotalG = math.Sqrt(diag.CurrentXG*diag.CurrentXG + diag.CurrentYG*diag.CurrentYG + diag.CurrentZG*diag.CurrentZG)
	n := float64(len(samples))
	diag.AvgAbsXG = sumAbsX / n
	diag.AvgAbsYG = sumAbsY / n
	diag.AvgAbsZG = sumAbsZ / n
	diag.AvgTotalG = sumTotal / n
	diag.DominantAxis = dominantGAxis(diag.PeakAbsXG, diag.PeakAbsYG, diag.PeakAbsZG)
	return diag
}

func buildCamberInference(samples []telemetry.NormalizedTelemetry, gForce GForceDiagnostic) CamberInference {
	result := CamberInference{
		Status:      "insufficient_data",
		Confidence:  quickConfidenceLow,
		FrontState:  "unknown",
		RearState:   "unknown",
		Summary:     "camber_inference_insufficient",
		Explanation: "camber_inference_needs_cornering",
		Warnings:    []string{"camber_inference_no_three_point_temps"},
		Hints:       []TireModelHint{},
		Evidence:    map[string]float64{},
	}
	stats := camberStats{}
	for _, sample := range samples {
		if math.Abs(sample.Steer01) < 0.18 || sample.SpeedKmh < 55 {
			continue
		}
		if sample.Brake01 > 0.35 || sample.Throttle01 > 0.75 {
			continue
		}
		stats.add(sample)
	}
	result.Evidence = stats.evidence()
	if stats.count < camberMinSamples {
		result.Evidence["cornering_samples"] = float64(stats.count)
		return result
	}
	result.Status = "ready"
	result.Confidence = quickConfidenceMedium
	if stats.count >= 24 && gForce.PeakTotalG >= 0.7 {
		result.Confidence = quickConfidenceHigh
	}
	result.FrontState = inferCamberAxleState(stats.frontSlipAngleAvg(), stats.frontCombinedAvg(), stats.frontTempAvg(), stats.frontSuspensionMax())
	result.RearState = inferCamberAxleState(stats.rearSlipAngleAvg(), stats.rearCombinedAvg(), stats.rearTempAvg(), stats.rearSuspensionMax())
	result.Summary = camberSummary(result.FrontState, result.RearState)
	result.Explanation = camberExplanation(result.FrontState, result.RearState)
	result.Hints = camberHints(result)
	return result
}

func BuildPowerToTireDiagnostic(samples []telemetry.NormalizedTelemetry, current *telemetry.NormalizedTelemetry) PowerToTireDiagnostic {
	diag := PowerToTireDiagnostic{
		Status:      "no_data",
		Summary:     "power_to_tire_no_data",
		Explanation: "power_to_tire_waiting_for_samples",
		Confidence:  quickConfidenceLow,
		DrivenAxle:  "unknown",
		Evidence:    map[string]float64{},
	}
	if len(samples) == 0 {
		return diag
	}
	ordered := sortedTelemetrySamples(samples)
	diag.Status = "insufficient_data"
	diag.Summary = "power_to_tire_insufficient"
	diag.Explanation = "power_to_tire_need_high_throttle"
	diag.SampleCount = len(ordered)
	diag.Drivetrain = powerToTireDrivetrain(ordered, current)
	diag.DrivenAxle = powerToTireDrivenAxle(diag.Drivetrain)
	last := ordered[len(ordered)-1]
	if current != nil {
		last = *current
	}
	diag.CurrentPowerKW = powerKW(last.Power)
	diag.CurrentTorqueNM = last.Torque
	diag.CurrentRPM = last.Rpm
	diag.CurrentRPMRatio = last.RpmRatio
	diag.CurrentGear = last.Gear

	stats := powerToTireStats{}
	for i, sample := range ordered {
		accelMps2 := 0.0
		if i > 0 {
			prev := ordered[i-1]
			dt := float64(sample.TimeMS-prev.TimeMS) / 1000
			if dt > 0 && dt <= 1.5 {
				accelMps2 = ((sample.SpeedKmh - prev.SpeedKmh) / 3.6) / dt
			}
		}
		stats.add(sample, accelMps2, diag.DrivenAxle)
	}
	stats.apply(&diag)
	if diag.HighThrottleSampleCount < powerToTireMinSamples {
		diag.Evidence["high_throttle_sample_count"] = float64(diag.HighThrottleSampleCount)
		return diag
	}
	if diag.AverageThrottle < powerToTireHighThrottle {
		diag.Summary = "power_to_tire_low_throttle"
		diag.Explanation = "power_to_tire_low_throttle_explanation"
		return diag
	}
	if !diag.PowerSignalAvailable {
		diag.Summary = "power_to_tire_power_signal_unavailable"
		diag.Explanation = "power_to_tire_power_signal_unavailable_explanation"
		return diag
	}

	diag.Status = "ready"
	diag.Confidence = quickConfidenceMedium
	if diag.HighThrottleSampleCount >= 18 {
		diag.Confidence = quickConfidenceHigh
	}
	switch {
	case diag.DrivenSlipRatioP90 >= powerToTireSlipLimit && diag.AverageAccelMps2 < 2.0:
		diag.Summary = "traction_over_power"
		diag.Explanation = "power_to_tire_traction_over_power_explanation"
		diag.TractionLimited = true
	case diag.AverageRPMRatio < powerToTireLowRPM && diag.DrivenSlipRatioP90 < 0.25 && diag.AverageAccelMps2 < 1.2:
		diag.Summary = "rpm_below_useful_range"
		diag.Explanation = "power_to_tire_rpm_below_explanation"
	case diag.RPMHighHighThrottlePct >= 0.55 && diag.AverageAccelMps2 < 1.8:
		diag.Summary = "rpm_too_high_or_gear_short"
		diag.Explanation = "power_to_tire_rpm_high_explanation"
	case (diag.AveragePowerKW >= 60 || math.Abs(diag.AverageTorqueNM) >= 120) && diag.AverageAccelMps2 < 0.8:
		diag.Summary = "power_not_reaching_ground"
		diag.Explanation = "power_to_tire_not_reaching_ground_explanation"
	default:
		diag.Summary = "power_landing_ok"
		diag.Explanation = "power_to_tire_ok_explanation"
	}
	diag.Evidence["traction_limited"] = boolFloat(diag.TractionLimited)
	return diag
}

func BuildBrakeToTireDiagnostic(samples []telemetry.NormalizedTelemetry) BrakeToTireDiagnostic {
	diag := BrakeToTireDiagnostic{
		Status:      "no_data",
		Summary:     "brake_to_tire_no_data",
		Explanation: "brake_to_tire_waiting_for_samples",
		Confidence:  quickConfidenceLow,
		Evidence:    map[string]float64{},
	}
	if len(samples) == 0 {
		return diag
	}
	ordered := sortedTelemetrySamples(samples)
	diag.Status = "insufficient_data"
	diag.Summary = "brake_to_tire_insufficient"
	diag.Explanation = "brake_to_tire_need_brake_samples"
	diag.SampleCount = len(ordered)

	stats := brakeToTireStats{}
	for i, sample := range ordered {
		decelMps2 := 0.0
		if i > 0 {
			prev := ordered[i-1]
			dt := float64(sample.TimeMS-prev.TimeMS) / 1000
			if dt > 0 && dt <= 1.5 {
				decelMps2 = ((prev.SpeedKmh - sample.SpeedKmh) / 3.6) / dt
				if decelMps2 < 0 {
					decelMps2 = 0
				}
			}
		}
		stats.add(sample, decelMps2)
	}
	stats.apply(&diag)
	if diag.PeakBrake < brakeToTireBrakeMin && diag.PeakHandBrake < brakeToTireHandbrakeMin {
		diag.Summary = "brake_to_tire_low_brake"
		diag.Explanation = "brake_to_tire_low_brake_explanation"
		return diag
	}
	if diag.BrakeSampleCount < brakeToTireMinSamples {
		return diag
	}

	diag.Status = "ready"
	diag.Confidence = quickConfidenceMedium
	if diag.BrakeSampleCount >= 18 {
		diag.Confidence = quickConfidenceHigh
	}
	switch {
	case diag.HandbrakeActive && diag.RearSlipRatioP90 >= brakeToTireSlipWarn:
		diag.Summary = "handbrake_rear_slide"
		diag.Explanation = "brake_to_tire_handbrake_rear_slide_explanation"
	case diag.TrailBraking &&
		diag.FrontCombinedSlipP90 >= brakeToTireCombinedWarn &&
		diag.FrontCombinedSlipP90 >= diag.RearCombinedSlipP90+brakeToTireAxleDelta:
		diag.Summary = "trail_brake_front_overload"
		diag.Explanation = "brake_to_tire_trail_brake_front_overload_explanation"
	case diag.FrontSlipRatioP90 >= brakeToTireSlipWarn && diag.FrontSlipRatioP90 >= diag.RearSlipRatioP90+brakeToTireAxleDelta:
		diag.Summary = "front_brake_lock_tendency"
		diag.Explanation = "brake_to_tire_front_lock_explanation"
	case diag.RearSlipRatioP90 >= brakeToTireSlipWarn && diag.RearSlipRatioP90 >= diag.FrontSlipRatioP90+brakeToTireAxleDelta:
		diag.Summary = "rear_brake_lock_tendency"
		diag.Explanation = "brake_to_tire_rear_lock_explanation"
	case diag.AverageBrake >= 0.55 &&
		diag.AverageDecelG < brakeToTireWeakDecelG &&
		math.Max(diag.FrontSlipRatioP90, diag.RearSlipRatioP90) < brakeToTireSlipWarn:
		diag.Summary = "brake_not_slowing_effectively"
		diag.Explanation = "brake_to_tire_not_slowing_effectively_explanation"
	default:
		diag.Summary = "brake_landing_ok"
		diag.Explanation = "brake_to_tire_ok_explanation"
	}
	return diag
}

type brakeToTireStats struct {
	activeCount                             int
	brakeSum, brakePeak                     float64
	handBrakeSum, handBrakePeak             float64
	speedSum, steerSum                      float64
	decelSum, decelPeakG                    float64
	planeGSum, planeGPeak                   float64
	firstSpeed, lastSpeed                   float64
	frontSlipValues, rearSlipValues         []float64
	frontCombinedValues, rearCombinedValues []float64
}

func (s *brakeToTireStats) add(sample telemetry.NormalizedTelemetry, decelMps2 float64) {
	s.brakePeak = math.Max(s.brakePeak, sample.Brake01)
	s.handBrakePeak = math.Max(s.handBrakePeak, sample.HandBrake01)
	if sample.Brake01 < brakeToTireBrakeSampleMin && sample.HandBrake01 < brakeToTireHandbrakeSample {
		return
	}
	if s.activeCount == 0 {
		s.firstSpeed = sample.SpeedKmh
	}
	s.activeCount++
	s.lastSpeed = sample.SpeedKmh
	s.brakeSum += sample.Brake01
	s.handBrakeSum += sample.HandBrake01
	s.speedSum += sample.SpeedKmh
	s.steerSum += math.Abs(sample.Steer01)
	s.decelSum += decelMps2
	s.decelPeakG = math.Max(s.decelPeakG, decelMps2/standardGravity)
	planeG := math.Sqrt(sample.AccelerationX*sample.AccelerationX+sample.AccelerationZ*sample.AccelerationZ) / standardGravity
	if !math.IsNaN(planeG) && !math.IsInf(planeG, 0) && planeG > 0 {
		s.planeGSum += planeG
		s.planeGPeak = math.Max(s.planeGPeak, planeG)
	}
	frontSlip := (math.Abs(sample.WheelFL.SlipRatio) + math.Abs(sample.WheelFR.SlipRatio)) / 2
	rearSlip := (math.Abs(sample.WheelRL.SlipRatio) + math.Abs(sample.WheelRR.SlipRatio)) / 2
	frontCombined := (math.Abs(sample.WheelFL.CombinedSlip) + math.Abs(sample.WheelFR.CombinedSlip)) / 2
	rearCombined := (math.Abs(sample.WheelRL.CombinedSlip) + math.Abs(sample.WheelRR.CombinedSlip)) / 2
	s.frontSlipValues = append(s.frontSlipValues, frontSlip)
	s.rearSlipValues = append(s.rearSlipValues, rearSlip)
	s.frontCombinedValues = append(s.frontCombinedValues, frontCombined)
	s.rearCombinedValues = append(s.rearCombinedValues, rearCombined)
}

func (s brakeToTireStats) apply(diag *BrakeToTireDiagnostic) {
	diag.BrakeSampleCount = s.activeCount
	diag.PeakBrake = s.brakePeak
	diag.PeakHandBrake = s.handBrakePeak
	if s.activeCount == 0 {
		return
	}
	n := float64(s.activeCount)
	diag.AverageBrake = s.brakeSum / n
	diag.AverageHandBrake = s.handBrakeSum / n
	diag.AverageSpeedKmh = s.speedSum / n
	diag.SpeedDeltaKmh = s.lastSpeed - s.firstSpeed
	diag.AverageSteer = s.steerSum / n
	diag.AverageDecelMps2 = s.decelSum / n
	diag.AverageDecelG = diag.AverageDecelMps2 / standardGravity
	diag.PeakDecelG = s.decelPeakG
	diag.AveragePlaneG = s.planeGSum / n
	diag.PeakPlaneG = s.planeGPeak
	diag.FrontSlipRatioP90 = percentile(s.frontSlipValues, 0.90)
	diag.RearSlipRatioP90 = percentile(s.rearSlipValues, 0.90)
	diag.FrontCombinedSlipP90 = percentile(s.frontCombinedValues, 0.90)
	diag.RearCombinedSlipP90 = percentile(s.rearCombinedValues, 0.90)
	diag.FrontRearSlipDelta = diag.FrontSlipRatioP90 - diag.RearSlipRatioP90
	diag.TrailBraking = diag.AverageBrake >= brakeToTireBrakeMin && diag.AverageSteer >= brakeToTireSteerMin
	diag.HandbrakeActive = diag.PeakHandBrake >= brakeToTireHandbrakeMin
	diag.Evidence = map[string]float64{
		"brake_sample_count":      float64(s.activeCount),
		"average_brake":           diag.AverageBrake,
		"peak_brake":              diag.PeakBrake,
		"average_handbrake":       diag.AverageHandBrake,
		"peak_handbrake":          diag.PeakHandBrake,
		"average_speed_kmh":       diag.AverageSpeedKmh,
		"speed_delta_kmh":         diag.SpeedDeltaKmh,
		"average_steer":           diag.AverageSteer,
		"average_decel_mps2":      diag.AverageDecelMps2,
		"average_decel_g":         diag.AverageDecelG,
		"peak_decel_g":            diag.PeakDecelG,
		"average_plane_g":         diag.AveragePlaneG,
		"peak_plane_g":            diag.PeakPlaneG,
		"front_slip_ratio_p90":    diag.FrontSlipRatioP90,
		"rear_slip_ratio_p90":     diag.RearSlipRatioP90,
		"front_combined_slip_p90": diag.FrontCombinedSlipP90,
		"rear_combined_slip_p90":  diag.RearCombinedSlipP90,
		"front_rear_slip_delta":   diag.FrontRearSlipDelta,
		"trail_braking":           boolFloat(diag.TrailBraking),
		"handbrake_active":        boolFloat(diag.HandbrakeActive),
	}
}

type powerToTireStats struct {
	highCount                 int
	powerSumKW, powerMaxKW    float64
	torqueSum, torqueMax      float64
	rpmSum, rpmRatioSum       float64
	throttleSum, speedSum     float64
	accelSum, accelMaxG       float64
	gForceSum                 float64
	firstSpeed, lastSpeed     float64
	frontSlipValues           []float64
	rearSlipValues            []float64
	drivenSlipValues          []float64
	lowRPMCount, highRPMCount int
	powerSignalCount          int
}

func (s *powerToTireStats) add(sample telemetry.NormalizedTelemetry, accelMps2 float64, drivenAxle string) {
	if sample.Throttle01 < powerToTireHighThrottle {
		return
	}
	if s.highCount == 0 {
		s.firstSpeed = sample.SpeedKmh
	}
	s.highCount++
	s.lastSpeed = sample.SpeedKmh
	power := powerKW(sample.Power)
	torque := sample.Torque
	s.powerSumKW += power
	s.powerMaxKW = math.Max(s.powerMaxKW, math.Abs(power))
	s.torqueSum += torque
	s.torqueMax = math.Max(s.torqueMax, math.Abs(torque))
	s.rpmSum += sample.Rpm
	s.rpmRatioSum += sample.RpmRatio
	s.throttleSum += sample.Throttle01
	s.speedSum += sample.SpeedKmh
	if accelMps2 > 0 {
		s.accelSum += accelMps2
		s.accelMaxG = math.Max(s.accelMaxG, accelMps2/standardGravity)
	}
	gPlane := math.Sqrt(sample.AccelerationX*sample.AccelerationX+sample.AccelerationZ*sample.AccelerationZ) / standardGravity
	if !math.IsNaN(gPlane) && !math.IsInf(gPlane, 0) && gPlane > 0 {
		s.gForceSum += gPlane
		s.accelMaxG = math.Max(s.accelMaxG, gPlane)
	}
	frontSlip := (math.Abs(sample.WheelFL.SlipRatio) + math.Abs(sample.WheelFR.SlipRatio)) / 2
	rearSlip := (math.Abs(sample.WheelRL.SlipRatio) + math.Abs(sample.WheelRR.SlipRatio)) / 2
	s.frontSlipValues = append(s.frontSlipValues, frontSlip)
	s.rearSlipValues = append(s.rearSlipValues, rearSlip)
	s.drivenSlipValues = append(s.drivenSlipValues, drivenSlipRatio(frontSlip, rearSlip, drivenAxle))
	if sample.RpmRatio > 0 && sample.RpmRatio < powerToTireLowRPM {
		s.lowRPMCount++
	}
	if sample.RpmRatio >= powerToTireHighRPM {
		s.highRPMCount++
	}
	if math.Abs(sample.Power) > 10000 || math.Abs(sample.Torque) > 50 {
		s.powerSignalCount++
	}
}

func (s powerToTireStats) apply(diag *PowerToTireDiagnostic) {
	diag.HighThrottleSampleCount = s.highCount
	if s.highCount == 0 {
		return
	}
	n := float64(s.highCount)
	diag.AveragePowerKW = s.powerSumKW / n
	diag.MaxPowerKW = s.powerMaxKW
	diag.AverageTorqueNM = s.torqueSum / n
	diag.MaxTorqueNM = s.torqueMax
	diag.AverageRPM = s.rpmSum / n
	diag.AverageRPMRatio = s.rpmRatioSum / n
	diag.AverageThrottle = s.throttleSum / n
	diag.AverageSpeedKmh = s.speedSum / n
	diag.SpeedDeltaKmh = s.lastSpeed - s.firstSpeed
	diag.AverageAccelMps2 = s.accelSum / n
	if s.gForceSum > 0 {
		diag.AverageAccelG = s.gForceSum / n
	} else {
		diag.AverageAccelG = diag.AverageAccelMps2 / standardGravity
	}
	diag.PeakAccelG = s.accelMaxG
	diag.FrontSlipRatioP90 = percentile(s.frontSlipValues, 0.90)
	diag.RearSlipRatioP90 = percentile(s.rearSlipValues, 0.90)
	diag.DrivenSlipRatioP90 = percentile(s.drivenSlipValues, 0.90)
	diag.DrivenSlipRatioHighPct = thresholdPercent(s.drivenSlipValues, powerToTireSlipWarn)
	diag.RPMLowHighThrottlePct = float64(s.lowRPMCount) / n
	diag.RPMHighHighThrottlePct = float64(s.highRPMCount) / n
	diag.PowerSignalAvailable = s.powerSignalCount >= maxInt(2, s.highCount/3)
	diag.Evidence = map[string]float64{
		"high_throttle_sample_count": float64(s.highCount),
		"average_power_kw":           diag.AveragePowerKW,
		"max_power_kw":               diag.MaxPowerKW,
		"average_torque_nm":          diag.AverageTorqueNM,
		"max_torque_nm":              diag.MaxTorqueNM,
		"average_rpm_ratio":          diag.AverageRPMRatio,
		"average_accel_mps2":         diag.AverageAccelMps2,
		"average_accel_g":            diag.AverageAccelG,
		"peak_accel_g":               diag.PeakAccelG,
		"front_slip_ratio_p90":       diag.FrontSlipRatioP90,
		"rear_slip_ratio_p90":        diag.RearSlipRatioP90,
		"driven_slip_ratio_p90":      diag.DrivenSlipRatioP90,
		"driven_slip_ratio_high_pct": diag.DrivenSlipRatioHighPct,
		"rpm_low_high_throttle_pct":  diag.RPMLowHighThrottlePct,
		"rpm_high_high_throttle_pct": diag.RPMHighHighThrottlePct,
		"power_signal_available":     boolFloat(diag.PowerSignalAvailable),
		"speed_delta_kmh":            diag.SpeedDeltaKmh,
	}
}

func powerKW(power float64) float64 {
	return power / 1000
}

func powerToTireDrivetrain(samples []telemetry.NormalizedTelemetry, current *telemetry.NormalizedTelemetry) string {
	if current != nil && strings.TrimSpace(current.Drivetrain) != "" {
		return strings.ToUpper(strings.TrimSpace(current.Drivetrain))
	}
	for i := len(samples) - 1; i >= 0; i-- {
		if strings.TrimSpace(samples[i].Drivetrain) != "" {
			return strings.ToUpper(strings.TrimSpace(samples[i].Drivetrain))
		}
	}
	return "UNKNOWN"
}

func powerToTireDrivenAxle(drivetrain string) string {
	switch strings.ToUpper(strings.TrimSpace(drivetrain)) {
	case "FWD":
		return "front"
	case "RWD":
		return "rear"
	case "AWD", "4WD":
		return "all"
	default:
		return "unknown"
	}
}

func drivenSlipRatio(frontSlip, rearSlip float64, drivenAxle string) float64 {
	switch drivenAxle {
	case "front":
		return frontSlip
	case "rear":
		return rearSlip
	case "all":
		return math.Max(frontSlip, rearSlip)
	default:
		return math.Max(frontSlip, rearSlip)
	}
}

func boolFloat(value bool) float64 {
	if value {
		return 1
	}
	return 0
}

type camberStats struct {
	count                               int
	frontSlipAngleSum, rearSlipAngleSum float64
	frontCombinedSum, rearCombinedSum   float64
	frontTempSum, rearTempSum           float64
	frontSuspMax, rearSuspMax           float64
	avgSpeedSum                         float64
}

func (s *camberStats) add(sample telemetry.NormalizedTelemetry) {
	s.count++
	s.frontSlipAngleSum += (math.Abs(sample.WheelFL.SlipAngle) + math.Abs(sample.WheelFR.SlipAngle)) / 2
	s.rearSlipAngleSum += (math.Abs(sample.WheelRL.SlipAngle) + math.Abs(sample.WheelRR.SlipAngle)) / 2
	s.frontCombinedSum += (math.Abs(sample.WheelFL.CombinedSlip) + math.Abs(sample.WheelFR.CombinedSlip)) / 2
	s.rearCombinedSum += (math.Abs(sample.WheelRL.CombinedSlip) + math.Abs(sample.WheelRR.CombinedSlip)) / 2
	s.frontTempSum += (sample.WheelFL.TireTemp + sample.WheelFR.TireTemp) / 2
	s.rearTempSum += (sample.WheelRL.TireTemp + sample.WheelRR.TireTemp) / 2
	s.frontSuspMax = math.Max(s.frontSuspMax, math.Max(math.Abs(sample.WheelFL.SuspensionTravel), math.Abs(sample.WheelFR.SuspensionTravel)))
	s.rearSuspMax = math.Max(s.rearSuspMax, math.Max(math.Abs(sample.WheelRL.SuspensionTravel), math.Abs(sample.WheelRR.SuspensionTravel)))
	s.avgSpeedSum += sample.SpeedKmh
}

func (s camberStats) frontSlipAngleAvg() float64 {
	if s.count == 0 {
		return 0
	}
	return s.frontSlipAngleSum / float64(s.count)
}

func (s camberStats) rearSlipAngleAvg() float64 {
	if s.count == 0 {
		return 0
	}
	return s.rearSlipAngleSum / float64(s.count)
}

func (s camberStats) frontCombinedAvg() float64 {
	if s.count == 0 {
		return 0
	}
	return s.frontCombinedSum / float64(s.count)
}

func (s camberStats) rearCombinedAvg() float64 {
	if s.count == 0 {
		return 0
	}
	return s.rearCombinedSum / float64(s.count)
}

func (s camberStats) frontTempAvg() float64 {
	if s.count == 0 {
		return 0
	}
	return s.frontTempSum / float64(s.count)
}

func (s camberStats) rearTempAvg() float64 {
	if s.count == 0 {
		return 0
	}
	return s.rearTempSum / float64(s.count)
}

func (s camberStats) frontSuspensionMax() float64 { return s.frontSuspMax }
func (s camberStats) rearSuspensionMax() float64  { return s.rearSuspMax }

func (s camberStats) evidence() map[string]float64 {
	evidence := map[string]float64{"cornering_samples": float64(s.count)}
	if s.count == 0 {
		return evidence
	}
	n := float64(s.count)
	evidence["front_slip_angle_avg"] = s.frontSlipAngleAvg()
	evidence["rear_slip_angle_avg"] = s.rearSlipAngleAvg()
	evidence["front_combined_slip_avg"] = s.frontCombinedAvg()
	evidence["rear_combined_slip_avg"] = s.rearCombinedAvg()
	evidence["front_tire_temp_avg"] = s.frontTempAvg()
	evidence["rear_tire_temp_avg"] = s.rearTempAvg()
	evidence["front_suspension_max"] = s.frontSuspMax
	evidence["rear_suspension_max"] = s.rearSuspMax
	evidence["avg_speed_kmh"] = s.avgSpeedSum / n
	return evidence
}

func inferCamberAxleState(slipAngle, combinedSlip, temp, suspensionMax float64) string {
	if suspensionMax >= tireModelBottomOut {
		return "platform_limited"
	}
	if temp >= tireModelHotTemp {
		return "thermal_limited"
	}
	if slipAngle >= 0.55 && combinedSlip >= 0.70 {
		return "likely_needs_more_negative"
	}
	if slipAngle <= 0.18 && combinedSlip <= 0.45 {
		return "stable"
	}
	return "monitor"
}

func camberSummary(frontState, rearState string) string {
	if frontState == "likely_needs_more_negative" && rearState == "likely_needs_more_negative" {
		return "camber_inference_both_axles_need_more_negative"
	}
	if frontState == "likely_needs_more_negative" {
		return "camber_inference_front_needs_more_negative"
	}
	if rearState == "likely_needs_more_negative" {
		return "camber_inference_rear_needs_more_negative"
	}
	if frontState == "platform_limited" || rearState == "platform_limited" {
		return "camber_inference_platform_first"
	}
	if frontState == "thermal_limited" || rearState == "thermal_limited" {
		return "camber_inference_temperature_first"
	}
	return "camber_inference_monitor"
}

func camberExplanation(frontState, rearState string) string {
	if frontState == "likely_needs_more_negative" || rearState == "likely_needs_more_negative" {
		return "camber_inference_slip_angle_explanation"
	}
	if frontState == "platform_limited" || rearState == "platform_limited" {
		return "camber_inference_platform_explanation"
	}
	if frontState == "thermal_limited" || rearState == "thermal_limited" {
		return "camber_inference_temperature_explanation"
	}
	return "camber_inference_monitor_explanation"
}

func camberHints(inference CamberInference) []TireModelHint {
	hints := []TireModelHint{}
	if inference.FrontState == "likely_needs_more_negative" {
		hints = append(hints, TireModelHint{
			Code:      "front_camber_check",
			Severity:  "warning",
			Direction: "consider_more_negative_front_camber",
			Reason:    inference.Explanation,
		})
	}
	if inference.RearState == "likely_needs_more_negative" {
		hints = append(hints, TireModelHint{
			Code:      "rear_camber_check",
			Severity:  "warning",
			Direction: "consider_more_negative_rear_camber",
			Reason:    inference.Explanation,
		})
	}
	if len(hints) == 0 {
		hints = append(hints, TireModelHint{
			Code:      "camber_observe",
			Severity:  "stable",
			Direction: "use_camber_inference_as_low_confidence_evidence",
			Reason:    inference.Explanation,
		})
	}
	return hints
}

func dominantGAxis(x, y, z float64) string {
	if x >= y && x >= z {
		return "x"
	}
	if y >= x && y >= z {
		return "y"
	}
	return "z"
}

func defaultTireDataQuality() TireDataQuality {
	return TireDataQuality{
		Status:     "invalid",
		Confidence: quickConfidenceLow,
		Reasons:    []string{"tire_data_no_samples"},
		Evidence:   map[string]float64{},
	}
}

func defaultTireGripLimit() TireGripLimit {
	return TireGripLimit{
		Type:          "no_limit_detected",
		LimitedAxle:   "none",
		LimitedWheels: []string{},
		Confidence:    quickConfidenceLow,
		Reason:        "tire_grip_no_dynamic_limit",
		Evidence:      map[string]float64{},
	}
}

func buildTireDataQuality(samples []telemetry.NormalizedTelemetry, dynamicSamples []telemetry.NormalizedTelemetry, current *telemetry.NormalizedTelemetry, diag TireModelDiagnostic) TireDataQuality {
	quality := TireDataQuality{
		Status:             "valid",
		Confidence:         quickConfidenceHigh,
		SampleCount:        len(samples),
		DynamicSampleCount: len(dynamicSamples),
		SpeedSignal:        "ok",
		GForceSignal:       "ok",
		SlipSignal:         "ok",
		InputSignal:        "ok",
		Reasons:            []string{},
		Evidence:           map[string]float64{},
	}
	if len(samples) == 0 {
		return defaultTireDataQuality()
	}
	var maxSpeed, throttle, brake, handbrake, steer, maxSlipActivity float64
	for _, sample := range samples {
		maxSpeed = math.Max(maxSpeed, sample.SpeedKmh)
		throttle += sample.Throttle01
		brake += sample.Brake01
		handbrake += sample.HandBrake01
		steer += math.Abs(sample.Steer01)
		maxSlipActivity = math.Max(maxSlipActivity, tireFrameSlipActivity(sample))
	}
	n := float64(len(samples))
	avgSpeed := diag.Evidence["avg_speed_kmh"]
	avgInput := (throttle + brake + handbrake + steer) / n
	quality.Evidence = map[string]float64{
		"avg_speed_kmh":        avgSpeed,
		"max_speed_kmh":        maxSpeed,
		"peak_total_g":         diag.GForce.PeakTotalG,
		"avg_input":            avgInput,
		"max_slip_activity":    maxSlipActivity,
		"sample_count":         float64(len(samples)),
		"dynamic_sample_count": float64(len(dynamicSamples)),
	}
	if diag.GameMode == telemetry.GameModeMenu {
		quality.Status = "invalid"
		quality.Confidence = quickConfidenceLow
		quality.Reasons = append(quality.Reasons, "tire_data_menu_or_no_vehicle")
	}
	if len(samples) < tireModelMinSamples {
		if quality.Status != "invalid" {
			quality.Status = "low_confidence"
		}
		quality.Confidence = quickConfidenceLow
		quality.Reasons = append(quality.Reasons, "tire_data_sample_insufficient")
	}
	if len(dynamicSamples) < tireModelMinDynamic {
		if quality.Status != "invalid" {
			quality.Status = "low_confidence"
		}
		quality.Confidence = quickConfidenceLow
		quality.Reasons = append(quality.Reasons, "tire_data_dynamic_sample_insufficient")
	}
	if tireModelCurrentStationary(current, avgSpeed, throttle/n) {
		if quality.Status != "invalid" {
			quality.Status = "low_confidence"
		}
		quality.Confidence = quickConfidenceLow
		quality.Reasons = append(quality.Reasons, "tire_data_stationary")
	}
	if maxSpeed < tireModelDynamicKmh {
		quality.SpeedSignal = "low"
		quality.Reasons = append(quality.Reasons, "tire_data_speed_low")
	}
	if diag.GForce.PeakTotalG < 0.08 {
		quality.GForceSignal = "flat"
		quality.Reasons = append(quality.Reasons, "tire_data_g_force_flat")
	}
	if maxSlipActivity <= 0 {
		quality.SlipSignal = "flat"
		if quality.Status != "invalid" {
			quality.Status = "low_confidence"
		}
		quality.Confidence = quickConfidenceLow
		quality.Reasons = append(quality.Reasons, "tire_data_slip_signal_flat")
	}
	if avgInput < 0.02 && maxSpeed < tireModelDynamicKmh {
		quality.InputSignal = "low"
		quality.Reasons = append(quality.Reasons, "tire_data_input_low")
	}
	if quality.Status == "valid" && len(dynamicSamples) < 20 {
		quality.Confidence = quickConfidenceMedium
	}
	return quality
}

func tireFrameSlipActivity(sample telemetry.NormalizedTelemetry) float64 {
	wheels := []telemetry.NormalizedWheelTelemetry{sample.WheelFL, sample.WheelFR, sample.WheelRL, sample.WheelRR}
	out := 0.0
	for _, wheel := range wheels {
		out = math.Max(out, math.Abs(wheel.CombinedSlip))
		out = math.Max(out, math.Abs(wheel.SlipRatio))
		out = math.Max(out, math.Abs(wheel.SlipAngle))
	}
	return out
}

func tireModelStationary(avgSpeed, avgThrottle float64) bool {
	return avgSpeed < tireModelParkedKmh && avgThrottle < 0.08
}

func tireModelCurrentStationary(current *telemetry.NormalizedTelemetry, avgSpeed, avgThrottle float64) bool {
	if current != nil {
		return tireModelStationary(current.SpeedKmh, current.Throttle01)
	}
	return tireModelStationary(avgSpeed, avgThrottle)
}

func markTireModelStationary(diag *TireModelDiagnostic) {
	diag.Phase = "stationary"
	diag.PhaseDetail.CurrentPhase = "stationary"
	diag.PhaseDetail.SecondaryPhase = "unknown"
	diag.PhaseDetail.StablePhase = "stationary"
	diag.PhaseDetail.PhaseStability = "low_confidence"
	diag.PhaseDetail.Confidence = quickConfidenceLow
	diag.Confidence = quickConfidenceLow
	diag.LimitType = "stationary"
	diag.GripLimit = defaultTireGripLimit()
	diag.GripLimit.Reason = "tire_grip_stationary"
	diag.Summary = "tire_model_stationary"
	diag.Explanation = "tire_model_stationary_explanation"
	diag.LeftRight.State = "balanced"
	for i := range diag.Wheels {
		diag.Wheels[i].GripState = "stable"
	}
	diag.FrontAxle.GripState = "stable"
	diag.RearAxle.GripState = "stable"
	diag.Hints = tireModelHints(*diag)
}

func markTireModelNoDynamicLoad(diag *TireModelDiagnostic) {
	diag.Confidence = quickConfidenceLow
	diag.LimitType = "no_dynamic_load"
	diag.GripLimit = defaultTireGripLimit()
	diag.GripLimit.Reason = "tire_grip_no_dynamic_load"
	diag.Summary = "tire_model_no_dynamic_load"
	diag.Explanation = "tire_model_no_dynamic_load_explanation"
	for i := range diag.Wheels {
		diag.Wheels[i].GripState = "stable"
	}
	diag.FrontAxle.GripState = "stable"
	diag.RearAxle.GripState = "stable"
	diag.Hints = tireModelHints(*diag)
}

func appendTireModelRisks(diag *TireModelDiagnostic) {
	if diag.FrontAxle.TireTempMax >= tireModelHotTemp || diag.RearAxle.TireTempMax >= tireModelHotTemp || math.Abs(diag.FrontAxle.TireTempAvg-diag.RearAxle.TireTempAvg) >= tireModelTempDelta {
		diag.Warnings = append(diag.Warnings, "thermal_risk")
	}
	if diag.FrontAxle.SuspensionTravelMax >= tireModelBottomOut || diag.RearAxle.SuspensionTravelMax >= tireModelBottomOut {
		diag.Warnings = append(diag.Warnings, "platform_risk")
	}
	if diag.LeftRight.State == "imbalanced" {
		diag.Warnings = append(diag.Warnings, "left_right_imbalance")
	}
	if diag.GForce.AxisMapping == "raw_packet_axes_unverified" {
		diag.Warnings = append(diag.Warnings, "g_force_axis_mapping_unverified")
	}
	diag.Warnings = append(diag.Warnings, diag.Camber.Warnings...)
}

func tireModelDynamicSamples(samples []telemetry.NormalizedTelemetry) []telemetry.NormalizedTelemetry {
	out := make([]telemetry.NormalizedTelemetry, 0, len(samples))
	for _, sample := range samples {
		totalG := math.Sqrt(sample.AccelerationX*sample.AccelerationX+sample.AccelerationY*sample.AccelerationY+sample.AccelerationZ*sample.AccelerationZ) / standardGravity
		if sample.SpeedKmh >= tireModelDynamicKmh ||
			sample.Throttle01 >= 0.18 ||
			sample.Brake01 >= 0.15 ||
			sample.HandBrake01 >= 0.10 ||
			math.Abs(sample.Steer01) >= 0.12 ||
			totalG >= 0.20 {
			out = append(out, sample)
		}
	}
	return out
}

func buildTireLimitDiagnostic(samples []telemetry.NormalizedTelemetry) TireModelDiagnostic {
	fl := buildTireWheel("front_left", samples, func(frame telemetry.NormalizedTelemetry) telemetry.NormalizedWheelTelemetry { return frame.WheelFL })
	fr := buildTireWheel("front_right", samples, func(frame telemetry.NormalizedTelemetry) telemetry.NormalizedWheelTelemetry { return frame.WheelFR })
	rl := buildTireWheel("rear_left", samples, func(frame telemetry.NormalizedTelemetry) telemetry.NormalizedWheelTelemetry { return frame.WheelRL })
	rr := buildTireWheel("rear_right", samples, func(frame telemetry.NormalizedTelemetry) telemetry.NormalizedWheelTelemetry { return frame.WheelRR })
	front := buildTireAxle("front", fl, fr)
	rear := buildTireAxle("rear", rl, rr)
	return TireModelDiagnostic{Wheels: []TireWheelDiagnostic{fl, fr, rl, rr}, FrontAxle: front, RearAxle: rear, LeftRight: buildTireSideBalance(fl, fr, rl, rr)}
}

func tireModelWindow(samples []telemetry.NormalizedTelemetry) []telemetry.NormalizedTelemetry {
	if len(samples) == 0 {
		return nil
	}
	latest := samples[len(samples)-1].TimeMS
	start := latest - tireModelWindowMS
	out := make([]telemetry.NormalizedTelemetry, 0, len(samples))
	for _, sample := range samples {
		if sample.TimeMS >= start {
			out = append(out, sample)
		}
	}
	if len(out) > tireModelMaxSamples {
		out = out[len(out)-tireModelMaxSamples:]
	}
	return out
}

func tireModelPhaseWindow(samples []telemetry.NormalizedTelemetry) []telemetry.NormalizedTelemetry {
	return tireModelCustomWindow(samples, tireModelPhaseWindowMS)
}

func tireModelCustomWindow(samples []telemetry.NormalizedTelemetry, windowMS int64) []telemetry.NormalizedTelemetry {
	if len(samples) == 0 {
		return nil
	}
	latest := samples[len(samples)-1].TimeMS
	start := latest - windowMS
	out := make([]telemetry.NormalizedTelemetry, 0, len(samples))
	for _, sample := range samples {
		if sample.TimeMS >= start {
			out = append(out, sample)
		}
	}
	if len(out) == 0 {
		return samples
	}
	return out
}

func buildTireWheel(position string, samples []telemetry.NormalizedTelemetry, pick func(telemetry.NormalizedTelemetry) telemetry.NormalizedWheelTelemetry) TireWheelDiagnostic {
	stat := tireWheelAccumulator{}
	for _, sample := range samples {
		stat.add(pick(sample))
	}
	wheel := stat.frame(position)
	wheel.GripState = tireGripState(wheel.CombinedSlipAvg, wheel.CombinedSlipP90, wheel.CombinedSlipHighPct)
	return wheel
}

type tireWheelAccumulator struct {
	count                    float64
	combinedSum, combinedMax float64
	ratioSum, ratioMax       float64
	angleSum, angleMax       float64
	tempSum, tempMax         float64
	suspSum, suspMax         float64
	suspMSum, suspMMax       float64
	combinedValues           []float64
	ratioValues              []float64
	angleValues              []float64
}

func (a *tireWheelAccumulator) add(w telemetry.NormalizedWheelTelemetry) {
	a.count++
	combined := math.Abs(w.CombinedSlip)
	ratio := math.Abs(w.SlipRatio)
	angle := math.Abs(w.SlipAngle)
	susp := math.Abs(w.SuspensionTravel)
	suspM := math.Abs(w.SuspensionTravelMeters)
	a.combinedSum += combined
	a.ratioSum += ratio
	a.angleSum += angle
	a.combinedValues = append(a.combinedValues, combined)
	a.ratioValues = append(a.ratioValues, ratio)
	a.angleValues = append(a.angleValues, angle)
	a.tempSum += w.TireTemp
	a.suspSum += susp
	a.suspMSum += suspM
	a.combinedMax = math.Max(a.combinedMax, combined)
	a.ratioMax = math.Max(a.ratioMax, ratio)
	a.angleMax = math.Max(a.angleMax, angle)
	a.tempMax = math.Max(a.tempMax, w.TireTemp)
	a.suspMax = math.Max(a.suspMax, susp)
	a.suspMMax = math.Max(a.suspMMax, suspM)
}

func (a tireWheelAccumulator) frame(position string) TireWheelDiagnostic {
	if a.count <= 0 {
		return TireWheelDiagnostic{Position: position, GripState: "unknown"}
	}
	return TireWheelDiagnostic{
		Position:               position,
		CombinedSlipAvg:        a.combinedSum / a.count,
		CombinedSlipMax:        a.combinedMax,
		CombinedSlipP90:        percentile(a.combinedValues, 0.90),
		CombinedSlipHighPct:    thresholdPercent(a.combinedValues, tireModelSlipWarn),
		SlipRatioAvg:           a.ratioSum / a.count,
		SlipRatioMax:           a.ratioMax,
		SlipRatioP90:           percentile(a.ratioValues, 0.90),
		SlipRatioHighPct:       thresholdPercent(a.ratioValues, tireModelRatioWarn),
		SlipAngleAvg:           a.angleSum / a.count,
		SlipAngleMax:           a.angleMax,
		SlipAngleP90:           percentile(a.angleValues, 0.90),
		TireTempAvg:            a.tempSum / a.count,
		TireTempMax:            a.tempMax,
		SuspensionTravelAvg:    a.suspSum / a.count,
		SuspensionTravelMax:    a.suspMax,
		SuspensionOffsetPctAvg: (a.suspSum / a.count) * 100,
		SuspensionOffsetPctMax: a.suspMax * 100,
		SuspensionTravelMAvg:   a.suspMSum / a.count,
		SuspensionTravelMMax:   a.suspMMax,
	}
}

func buildTireAxle(name string, left TireWheelDiagnostic, right TireWheelDiagnostic) TireAxleDiagnostic {
	axle := TireAxleDiagnostic{
		Name:                   name,
		CombinedSlipAvg:        (left.CombinedSlipAvg + right.CombinedSlipAvg) / 2,
		CombinedSlipMax:        math.Max(left.CombinedSlipMax, right.CombinedSlipMax),
		CombinedSlipP90:        math.Max(left.CombinedSlipP90, right.CombinedSlipP90),
		CombinedSlipHighPct:    math.Max(left.CombinedSlipHighPct, right.CombinedSlipHighPct),
		SlipRatioAvg:           (left.SlipRatioAvg + right.SlipRatioAvg) / 2,
		SlipRatioMax:           math.Max(left.SlipRatioMax, right.SlipRatioMax),
		SlipRatioP90:           math.Max(left.SlipRatioP90, right.SlipRatioP90),
		SlipRatioHighPct:       math.Max(left.SlipRatioHighPct, right.SlipRatioHighPct),
		SlipAngleAvg:           (left.SlipAngleAvg + right.SlipAngleAvg) / 2,
		SlipAngleMax:           math.Max(left.SlipAngleMax, right.SlipAngleMax),
		SlipAngleP90:           math.Max(left.SlipAngleP90, right.SlipAngleP90),
		TireTempAvg:            (left.TireTempAvg + right.TireTempAvg) / 2,
		TireTempMax:            math.Max(left.TireTempMax, right.TireTempMax),
		SuspensionTravelAvg:    (left.SuspensionTravelAvg + right.SuspensionTravelAvg) / 2,
		SuspensionTravelMax:    math.Max(left.SuspensionTravelMax, right.SuspensionTravelMax),
		SuspensionOffsetPctAvg: (left.SuspensionOffsetPctAvg + right.SuspensionOffsetPctAvg) / 2,
		SuspensionOffsetPctMax: math.Max(left.SuspensionOffsetPctMax, right.SuspensionOffsetPctMax),
	}
	axle.LimitScore = axle.CombinedSlipAvg*0.4 + axle.CombinedSlipP90*0.6
	axle.GripState = tireGripState(axle.CombinedSlipAvg, axle.CombinedSlipP90, axle.CombinedSlipHighPct)
	return axle
}

func buildTireSideBalance(fl, fr, rl, rr TireWheelDiagnostic) TireSideBalance {
	left := (fl.CombinedSlipAvg + rl.CombinedSlipAvg) / 2
	right := (fr.CombinedSlipAvg + rr.CombinedSlipAvg) / 2
	delta := math.Abs(left - right)
	state := "balanced"
	if delta >= tireModelSideDelta {
		state = "imbalanced"
	}
	return TireSideBalance{LeftCombinedSlipAvg: left, RightCombinedSlipAvg: right, Delta: delta, State: state}
}

func tireGripState(avg, p90, highPct float64) string {
	switch {
	case p90 >= tireModelSlipLimit || highPct >= 0.45:
		return "limit"
	case avg >= tireModelSlipWarn || p90 >= tireModelSlipWarn || highPct >= tireModelSlipHigh:
		return "warning"
	case avg > 0 || p90 > 0:
		return "stable"
	default:
		return "unknown"
	}
}

func percentile(values []float64, q float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sorted := append([]float64(nil), values...)
	sort.Float64s(sorted)
	if len(sorted) == 1 {
		return sorted[0]
	}
	if q <= 0 {
		return sorted[0]
	}
	if q >= 1 {
		return sorted[len(sorted)-1]
	}
	index := int(math.Ceil(q*float64(len(sorted)))) - 1
	if index < 0 {
		index = 0
	}
	if index >= len(sorted) {
		index = len(sorted) - 1
	}
	return sorted[index]
}

func thresholdPercent(values []float64, threshold float64) float64 {
	if len(values) == 0 {
		return 0
	}
	count := 0
	for _, value := range values {
		if value >= threshold {
			count++
		}
	}
	return float64(count) / float64(len(values))
}

func tireModelInputs(samples []telemetry.NormalizedTelemetry) (throttle, brake, handbrake, steer, speed float64) {
	if len(samples) == 0 {
		return 0, 0, 0, 0, 0
	}
	for _, sample := range samples {
		throttle += sample.Throttle01
		brake += sample.Brake01
		handbrake += sample.HandBrake01
		steer += math.Abs(sample.Steer01)
		speed += sample.SpeedKmh
	}
	n := float64(len(samples))
	return throttle / n, brake / n, handbrake / n, steer / n, speed / n
}

func BuildTirePhaseDiagnostic(samples []telemetry.NormalizedTelemetry) TirePhaseDiagnostic {
	return BuildTirePhaseDiagnosticWithReference(samples, samples)
}

func BuildTirePhaseDiagnosticWithReference(samples []telemetry.NormalizedTelemetry, referenceSamples []telemetry.NormalizedTelemetry) TirePhaseDiagnostic {
	result := TirePhaseDiagnostic{
		CurrentPhase:   "unknown",
		SecondaryPhase: "unknown",
		StablePhase:    "unknown",
		PhaseStability: "low_confidence",
		Confidence:     quickConfidenceLow,
		Scores:         map[string]float64{},
		Evidence:       map[string]float64{},
	}
	if len(samples) == 0 {
		return result
	}
	ordered := sortedTelemetrySamples(samples)
	reference := sortedTelemetrySamples(referenceSamples)
	if len(reference) == 0 {
		reference = ordered
	}
	stats := buildTirePhaseStats(ordered, reference)
	result.SampleCount = stats.count
	result.WindowMS = stats.windowMS
	result.Evidence = stats.evidence()
	result.Scores = stats.scores()
	best, second, bestScore := pickTirePhase(result.Scores)
	secondScore := result.Scores[second]
	result.ScoreMargin = bestScore - secondScore
	result.StablePhase = tireStablePhase(reference)
	result.PhaseStability = tirePhaseStability(best, result.StablePhase, result.ScoreMargin, stats.count, stats.windowMS)
	if bestScore < 0.35 {
		result.CurrentPhase = "unknown"
		result.SecondaryPhase = second
		result.Confidence = quickConfidenceLow
		return result
	}
	result.CurrentPhase = best
	result.SecondaryPhase = second
	result.Confidence = tirePhaseConfidence(bestScore, stats.count, stats.windowMS)
	return result
}

func tireStablePhase(samples []telemetry.NormalizedTelemetry) string {
	window := tireModelCustomWindow(samples, tireModelTrendWindowMS)
	if len(window) == 0 {
		return "unknown"
	}
	stats := buildTirePhaseStats(window, samples)
	scores := stats.scores()
	best, _, bestScore := pickTirePhase(scores)
	if bestScore < 0.35 {
		return "unknown"
	}
	return best
}

func tirePhaseStability(current string, stable string, margin float64, sampleCount int, windowMS int64) string {
	if current == "unknown" || sampleCount < 4 || windowMS < 300 || margin < 0.05 {
		return "low_confidence"
	}
	if stable != "" && stable != "unknown" && stable != current {
		return "transition"
	}
	if margin < 0.12 {
		return "transition"
	}
	return "stable"
}

type tirePhaseStats struct {
	count                             int
	windowMS                          int64
	speedSum, throttleSum, brakeSum   float64
	handbrakeSum, steerSum, planeGSum float64
	lateralGSum, decelSum, accelSum   float64
	frontCombinedSum, rearCombinedSum float64
	frontRatioSum, rearRatioSum       float64
	frontAngleSum, rearAngleSum       float64
	peakPlaneG, peakLateralG          float64
	peakDecelG, peakAccelG            float64
	peakBrake, peakHandbrake          float64
	firstSpeed, lastSpeed             float64
	firstThrottle, lastThrottle       float64
	firstSteer, lastSteer             float64
	firstSteerSigned, lastSteerSigned float64
	speedReferenceKmh                 float64
	speedBand                         float64
	speedBandConfidence               float64
}

func buildTirePhaseStats(samples []telemetry.NormalizedTelemetry, referenceSamples []telemetry.NormalizedTelemetry) tirePhaseStats {
	stats := tirePhaseStats{}
	for i, sample := range samples {
		if stats.count == 0 {
			stats.firstSpeed = sample.SpeedKmh
			stats.firstThrottle = sample.Throttle01
			stats.firstSteer = math.Abs(sample.Steer01)
			stats.firstSteerSigned = sample.Steer01
		}
		stats.count++
		stats.lastSpeed = sample.SpeedKmh
		stats.lastThrottle = sample.Throttle01
		stats.lastSteer = math.Abs(sample.Steer01)
		stats.lastSteerSigned = sample.Steer01
		stats.speedSum += sample.SpeedKmh
		stats.throttleSum += sample.Throttle01
		stats.brakeSum += sample.Brake01
		stats.handbrakeSum += sample.HandBrake01
		stats.steerSum += math.Abs(sample.Steer01)
		stats.peakBrake = math.Max(stats.peakBrake, sample.Brake01)
		stats.peakHandbrake = math.Max(stats.peakHandbrake, sample.HandBrake01)
		stats.frontCombinedSum += (math.Abs(sample.WheelFL.CombinedSlip) + math.Abs(sample.WheelFR.CombinedSlip)) / 2
		stats.rearCombinedSum += (math.Abs(sample.WheelRL.CombinedSlip) + math.Abs(sample.WheelRR.CombinedSlip)) / 2
		stats.frontRatioSum += (math.Abs(sample.WheelFL.SlipRatio) + math.Abs(sample.WheelFR.SlipRatio)) / 2
		stats.rearRatioSum += (math.Abs(sample.WheelRL.SlipRatio) + math.Abs(sample.WheelRR.SlipRatio)) / 2
		stats.frontAngleSum += (math.Abs(sample.WheelFL.SlipAngle) + math.Abs(sample.WheelFR.SlipAngle)) / 2
		stats.rearAngleSum += (math.Abs(sample.WheelRL.SlipAngle) + math.Abs(sample.WheelRR.SlipAngle)) / 2
		planeG := math.Sqrt(sample.AccelerationX*sample.AccelerationX+sample.AccelerationZ*sample.AccelerationZ) / standardGravity
		if !math.IsNaN(planeG) && !math.IsInf(planeG, 0) {
			stats.planeGSum += planeG
			stats.peakPlaneG = math.Max(stats.peakPlaneG, planeG)
		}
		lateralG := math.Abs(sample.AccelerationX) / standardGravity
		if !math.IsNaN(lateralG) && !math.IsInf(lateralG, 0) {
			stats.lateralGSum += lateralG
			stats.peakLateralG = math.Max(stats.peakLateralG, lateralG)
		}
		if i > 0 {
			prev := samples[i-1]
			dt := float64(sample.TimeMS-prev.TimeMS) / 1000
			if dt > 0 && dt <= 1.5 {
				accel := ((sample.SpeedKmh - prev.SpeedKmh) / 3.6) / dt
				if accel > 0 {
					stats.accelSum += accel
					stats.peakAccelG = math.Max(stats.peakAccelG, accel/standardGravity)
				}
				if accel < 0 {
					decel := math.Abs(accel)
					stats.decelSum += decel
					stats.peakDecelG = math.Max(stats.peakDecelG, decel/standardGravity)
				}
			}
		}
	}
	if len(samples) >= 2 {
		stats.windowMS = samples[len(samples)-1].TimeMS - samples[0].TimeMS
		if stats.windowMS < 0 {
			stats.windowMS = 0
		}
	}
	stats.speedReferenceKmh, stats.speedBandConfidence = tirePhaseSpeedReference(referenceSamples)
	if stats.count > 0 {
		stats.speedBand = tirePhaseSpeedBand(stats.speedSum/float64(stats.count), stats.speedReferenceKmh, stats.speedBandConfidence)
	}
	return stats
}

func (s tirePhaseStats) evidence() map[string]float64 {
	if s.count == 0 {
		return map[string]float64{}
	}
	n := float64(s.count)
	return map[string]float64{
		"avg_speed_kmh":         s.speedSum / n,
		"speed_delta_kmh":       s.lastSpeed - s.firstSpeed,
		"avg_throttle":          s.throttleSum / n,
		"throttle_delta":        s.lastThrottle - s.firstThrottle,
		"avg_brake":             s.brakeSum / n,
		"peak_brake":            s.peakBrake,
		"avg_handbrake":         s.handbrakeSum / n,
		"peak_handbrake":        s.peakHandbrake,
		"avg_steer":             s.steerSum / n,
		"steer_delta":           s.lastSteer - s.firstSteer,
		"avg_plane_g":           s.planeGSum / n,
		"peak_plane_g":          s.peakPlaneG,
		"avg_lateral_g":         s.lateralGSum / n,
		"peak_lateral_g":        s.peakLateralG,
		"front_combined_slip":   s.frontCombinedSum / n,
		"rear_combined_slip":    s.rearCombinedSum / n,
		"front_slip_ratio":      s.frontRatioSum / n,
		"rear_slip_ratio":       s.rearRatioSum / n,
		"front_slip_angle":      s.frontAngleSum / n,
		"rear_slip_angle":       s.rearAngleSum / n,
		"steer_sign_change":     boolFloat(s.firstSteerSigned*s.lastSteerSigned < -0.04),
		"avg_decel_g":           (s.decelSum / n) / standardGravity,
		"peak_decel_g":          s.peakDecelG,
		"avg_accel_g":           (s.accelSum / n) / standardGravity,
		"peak_accel_g":          s.peakAccelG,
		"speed_reference_kmh":   s.speedReferenceKmh,
		"speed_band":            s.speedBand,
		"speed_band_confidence": s.speedBandConfidence,
		"window_ms":             float64(s.windowMS),
		"sample_count":          float64(s.count),
	}
}

func (s tirePhaseStats) scores() map[string]float64 {
	e := s.evidence()
	speed := e["avg_speed_kmh"]
	speedDelta := e["speed_delta_kmh"]
	throttle := e["avg_throttle"]
	throttleDelta := e["throttle_delta"]
	brake := e["avg_brake"]
	handbrake := e["peak_handbrake"]
	steer := e["avg_steer"]
	steerDelta := e["steer_delta"]
	lateralG := e["avg_lateral_g"]
	decelG := e["avg_decel_g"]
	accelG := e["avg_accel_g"]
	speedBand := e["speed_band"]
	frontCombined := e["front_combined_slip"]
	rearCombined := e["rear_combined_slip"]
	rearRatio := e["rear_slip_ratio"]
	rearAngle := e["rear_slip_angle"]
	stationary := 0.0
	if speed < tireModelParkedKmh && throttle < 0.08 && brake < 0.08 && handbrake < 0.08 && e["peak_plane_g"] < 0.12 {
		stationary = 1
	}
	peakBrake := e["peak_brake"]
	heavyBrakeIntent := math.Max(clamp01(brake/0.45), clamp01(peakBrake/0.70))
	maxBrake := math.Max(brake, peakBrake)
	heavyBrakeCommit := clamp01((maxBrake - 0.35) / 0.35)
	lightBrakeIntent := clamp01((maxBrake - 0.10) / 0.28)
	coastIntent := clamp01((0.32 - throttle) / 0.32)
	steerIntent := clamp01(steer / 0.28)
	steerRise := clamp01(steerDelta / 0.18)
	steerUnwind := clamp01(-steerDelta / 0.18)
	decelIntent := math.Max(clamp01(decelG/0.28), clamp01(-speedDelta/18.0))
	accelIntent := math.Max(clamp01(accelG/0.22), clamp01(speedDelta/16.0))
	throttleIntent := math.Max(clamp01(throttle/0.65), clamp01(throttleDelta/0.30))
	throttleLift := clamp01(-throttleDelta / 0.25)
	straightIntent := clamp01((0.16 - steer) / 0.16)
	corneringLoad := math.Max(steerIntent, clamp01(lateralG/0.45))
	inputSteady := 1 - clamp01((math.Abs(throttleDelta)+math.Abs(steerDelta)+brake)/0.75)
	noDominantInput := 1 - clamp01((brake+math.Max(throttle-0.45, 0)+math.Abs(throttleDelta)+math.Abs(steerDelta))/0.80)
	cornerBase := clamp01(corneringLoad*0.45+clamp01(lateralG/0.55)*0.25+inputSteady*0.15+noDominantInput*0.15) * clamp01(noDominantInput)
	cornerEntry := clamp01(math.Max(heavyBrakeIntent, coastIntent)*0.35+steerIntent*0.25+decelIntent*0.25+steerRise*0.15) * clamp01(steer/0.18+steerRise*0.25)
	launchSpeedIntent := clamp01((55 - speed) / 45)
	launchIntent := clamp01((throttleIntent*0.45 + accelIntent*0.35 + straightIntent*0.20) * launchSpeedIntent)
	straightPower := clamp01(throttleIntent*0.48+accelIntent*0.32+straightIntent*0.20) * (1 - launchSpeedIntent*0.55)
	straightDecel := clamp01(straightIntent*0.42+decelIntent*0.38+coastIntent*0.20) * (1 - clamp01(brake/0.18)*0.75)
	rearOverFront := clamp01((rearCombined - frontCombined) / 0.35)
	rearSlipIntent := math.Max(clamp01(rearCombined/0.75), math.Max(clamp01(rearRatio/0.35), clamp01(rearAngle/0.50)))
	handbrakeDrift := clamp01(handbrake/0.35) * math.Max(steerIntent, rearSlipIntent)
	powerOversteer := clamp01(throttleIntent*0.45+rearSlipIntent*0.40+rearOverFront*0.15) * steerIntent
	liftOffOversteer := clamp01(throttleLift*0.45+rearSlipIntent*0.35+decelIntent*0.20) * steerIntent
	flickOversteer := clamp01(e["steer_sign_change"]*0.45+math.Abs(steerDelta)/0.35*0.25+rearOverFront*0.30) * steerIntent
	driftIntent := clamp01(math.Max(math.Max(handbrakeDrift, powerOversteer), math.Max(liftOffOversteer, flickOversteer)))

	return map[string]float64{
		"stationary":          stationary,
		"handbrake":           clamp01(handbrake/0.35 + clamp01(speed/20)*0.15),
		"launch":              launchIntent,
		"straight_decel":      straightDecel,
		"light_braking":       clamp01(lightBrakeIntent*0.55+decelIntent*0.35+straightIntent*0.20) * (1 - clamp01((maxBrake-0.55)/0.25)),
		"braking":             clamp01(heavyBrakeCommit*0.58 + decelIntent*0.27 + straightIntent*0.15),
		"corner_entry":        clamp01(cornerEntry),
		"low_speed_corner":    clamp01(cornerBase * tirePhaseBandWeight(speedBand, 1)),
		"mid_speed_corner":    clamp01(cornerBase * tirePhaseBandWeight(speedBand, 2)),
		"high_speed_corner":   clamp01(cornerBase * tirePhaseBandWeight(speedBand, 3)),
		"drift":               driftIntent,
		"sustained_cornering": clamp01(cornerBase * 0.65),
		"corner_exit":         clamp01(steerIntent*0.28 + throttleIntent*0.34 + accelIntent*0.25 + steerUnwind*0.13),
		"straight_power":      clamp01(straightPower),
	}
}

func tirePhaseSpeedReference(samples []telemetry.NormalizedTelemetry) (referenceKmh float64, confidence float64) {
	if len(samples) == 0 {
		return 0, 0
	}
	speeds := make([]float64, 0, len(samples))
	minSpeed := math.MaxFloat64
	maxSpeed := 0.0
	for _, sample := range samples {
		if sample.SpeedKmh < tireModelDynamicKmh {
			continue
		}
		speeds = append(speeds, sample.SpeedKmh)
		minSpeed = math.Min(minSpeed, sample.SpeedKmh)
		maxSpeed = math.Max(maxSpeed, sample.SpeedKmh)
	}
	if len(speeds) == 0 {
		return 0, 0
	}
	p90 := percentile(speeds, 0.90)
	referenceKmh = math.Max(p90, maxSpeed*0.90)
	if len(speeds) >= 16 && referenceKmh >= 80 && maxSpeed-minSpeed >= 25 {
		confidence = 1
	}
	return referenceKmh, confidence
}

func tirePhaseSpeedBand(speedKmh, referenceKmh, confidence float64) float64 {
	if confidence >= 0.5 && referenceKmh > 0 {
		ratio := speedKmh / referenceKmh
		switch {
		case ratio < 0.45:
			return 1
		case ratio < 0.70:
			return 2
		default:
			return 3
		}
	}
	switch {
	case speedKmh < 90:
		return 1
	case speedKmh < 160:
		return 2
	default:
		return 3
	}
}

func tirePhaseBandWeight(speedBand float64, target float64) float64 {
	if math.Round(speedBand) == target {
		return 1
	}
	return 0
}

func pickTirePhase(scores map[string]float64) (best, second string, bestScore float64) {
	order := []string{"stationary", "handbrake", "drift", "launch", "braking", "light_braking", "corner_entry", "corner_exit", "high_speed_corner", "mid_speed_corner", "low_speed_corner", "straight_power", "straight_decel", "sustained_cornering"}
	best = "unknown"
	second = "unknown"
	secondScore := -1.0
	bestScore = -1
	for _, phase := range order {
		score := scores[phase]
		if score > bestScore {
			second, secondScore = best, bestScore
			best, bestScore = phase, score
			continue
		}
		if score > secondScore {
			second, secondScore = phase, score
		}
	}
	if best == "stationary" && bestScore >= 0.85 {
		return best, second, bestScore
	}
	if best == "handbrake" && bestScore >= 0.55 {
		return best, second, bestScore
	}
	if bestScore < 0.35 {
		return "unknown", second, bestScore
	}
	return best, second, bestScore
}

func tirePhaseConfidence(score float64, sampleCount int, windowMS int64) string {
	if score >= 0.70 && sampleCount >= 8 && windowMS >= 1000 {
		return quickConfidenceHigh
	}
	if score >= 0.48 && sampleCount >= 4 {
		return quickConfidenceMedium
	}
	return quickConfidenceLow
}

func tireModelPhase(samples []telemetry.NormalizedTelemetry) string {
	return BuildTirePhaseDiagnostic(samples).CurrentPhase
}

func tireModelConfidence(sampleCount int, windowMS int64) string {
	if sampleCount >= 40 && windowMS >= 4000 {
		return quickConfidenceHigh
	}
	if sampleCount >= tireModelMinSamples {
		return quickConfidenceMedium
	}
	return quickConfidenceLow
}

func classifyTireLimit(diag TireModelDiagnostic, throttle, brake, handbrake float64) (limitType, summary, explanation string) {
	front := diag.FrontAxle
	rear := diag.RearAxle
	frontLimited := front.LimitScore >= tireModelSlipWarn || front.CombinedSlipP90 >= tireModelSlipWarn || front.CombinedSlipHighPct >= tireModelSlipHigh
	rearLimited := rear.LimitScore >= tireModelSlipWarn || rear.CombinedSlipP90 >= tireModelSlipWarn || rear.CombinedSlipHighPct >= tireModelSlipHigh
	if frontLimited && rearLimited && math.Abs(front.LimitScore-rear.LimitScore) <= tireModelAxleDelta {
		return "four_wheel_limited", "tire_model_four_wheel_limited", "tire_model_four_wheel_explanation"
	}
	frontTraction := front.SlipRatioP90 >= tireModelRatioWarn || front.SlipRatioHighPct >= tireModelSlipHigh
	rearTraction := rear.SlipRatioP90 >= tireModelRatioWarn || rear.SlipRatioHighPct >= tireModelSlipHigh
	if throttle >= 0.35 && (frontTraction || rearTraction) && math.Max(front.SlipRatioP90, rear.SlipRatioP90) >= math.Max(front.SlipAngleP90, rear.SlipAngleP90)*0.9 {
		if rear.SlipRatioP90 >= front.SlipRatioP90+tireModelAxleDelta || (rearTraction && !frontTraction) {
			return "traction_limited", "tire_model_rear_traction_limited", "tire_model_rear_traction_explanation"
		}
		if front.SlipRatioP90 >= rear.SlipRatioP90+tireModelAxleDelta || (frontTraction && !rearTraction) {
			return "traction_limited", "tire_model_front_traction_limited", "tire_model_front_traction_explanation"
		}
		return "traction_limited", "tire_model_drive_traction_limited", "tire_model_drive_traction_explanation"
	}
	if front.LimitScore >= rear.LimitScore+tireModelAxleDelta && frontLimited {
		if brake >= 0.25 {
			return "front_limited", "tire_model_front_brake_limited", "tire_model_front_brake_explanation"
		}
		return "front_limited", "tire_model_front_limited", "tire_model_front_explanation"
	}
	if rear.LimitScore >= front.LimitScore+tireModelAxleDelta && rearLimited {
		if handbrake >= 0.20 {
			return "rear_limited", "tire_model_rear_handbrake_limited", "tire_model_rear_handbrake_explanation"
		}
		if throttle >= 0.35 {
			return "rear_limited", "tire_model_rear_power_limited", "tire_model_rear_power_explanation"
		}
		return "rear_limited", "tire_model_rear_limited", "tire_model_rear_explanation"
	}
	if frontLimited || rearLimited {
		return "balanced_near_limit", "tire_model_balanced_near_limit", "tire_model_balanced_near_limit_explanation"
	}
	return "balanced", "tire_model_balanced", "tire_model_balanced_explanation"
}

func buildTireGripLimit(diag TireModelDiagnostic, leftRight TireSideBalance, throttle, brake, handbrake float64, quality TireDataQuality) TireGripLimit {
	front := diag.FrontAxle
	rear := diag.RearAxle
	limit := defaultTireGripLimit()
	limit.Confidence = quality.Confidence
	limit.LeftRightDelta = leftRight.Delta
	limit.FrontRearDelta = front.LimitScore - rear.LimitScore
	limit.DrivenDelta = math.Max(front.SlipRatioP90, rear.SlipRatioP90) - math.Min(front.SlipRatioP90, rear.SlipRatioP90)
	limit.Evidence = map[string]float64{
		"front_combined_p90":      front.CombinedSlipP90,
		"rear_combined_p90":       rear.CombinedSlipP90,
		"front_slip_ratio_p90":    front.SlipRatioP90,
		"rear_slip_ratio_p90":     rear.SlipRatioP90,
		"front_slip_angle_p90":    front.SlipAngleP90,
		"rear_slip_angle_p90":     rear.SlipAngleP90,
		"front_limit_score":       front.LimitScore,
		"rear_limit_score":        rear.LimitScore,
		"front_rear_delta":        limit.FrontRearDelta,
		"left_right_delta":        limit.LeftRightDelta,
		"driven_slip_ratio_delta": limit.DrivenDelta,
	}
	if quality.Status == "invalid" {
		limit.Reason = "tire_grip_data_invalid"
		return limit
	}
	frontCombined := front.LimitScore >= tireModelSlipWarn || front.CombinedSlipP90 >= tireModelSlipWarn || front.CombinedSlipHighPct >= tireModelSlipHigh
	rearCombined := rear.LimitScore >= tireModelSlipWarn || rear.CombinedSlipP90 >= tireModelSlipWarn || rear.CombinedSlipHighPct >= tireModelSlipHigh
	frontLateral := front.SlipAngleP90 >= tireModelRatioWarn && front.SlipAngleP90 >= front.SlipRatioP90*1.05 && frontCombined
	rearLateral := rear.SlipAngleP90 >= tireModelRatioWarn && rear.SlipAngleP90 >= rear.SlipRatioP90*1.05 && rearCombined
	frontTraction := front.SlipRatioP90 >= tireModelRatioWarn || front.SlipRatioHighPct >= tireModelSlipHigh
	rearTraction := rear.SlipRatioP90 >= tireModelRatioWarn || rear.SlipRatioHighPct >= tireModelSlipHigh
	switch {
	case brake >= 0.25 && (frontTraction || rearTraction):
		limit.Type = "braking_limit"
		limit.PrimaryEvidence = "slip_ratio_p90"
		if front.SlipRatioP90 >= rear.SlipRatioP90+tireModelAxleDelta || (frontTraction && !rearTraction) {
			limit.LimitedAxle = "front"
			limit.LimitedWheels = []string{"front_left", "front_right"}
		} else if rear.SlipRatioP90 >= front.SlipRatioP90+tireModelAxleDelta || (rearTraction && !frontTraction) {
			limit.LimitedAxle = "rear"
			limit.LimitedWheels = []string{"rear_left", "rear_right"}
		} else {
			limit.LimitedAxle = "both"
			limit.LimitedWheels = []string{"front_left", "front_right", "rear_left", "rear_right"}
		}
		limit.Reason = "tire_grip_braking_slip"
	case throttle >= 0.35 && (frontTraction || rearTraction):
		limit.Type = "traction_limit"
		limit.PrimaryEvidence = "slip_ratio_p90"
		if rear.SlipRatioP90 >= front.SlipRatioP90+tireModelAxleDelta || (rearTraction && !frontTraction) {
			limit.LimitedAxle = "rear"
			limit.LimitedWheels = []string{"rear_left", "rear_right"}
		} else if front.SlipRatioP90 >= rear.SlipRatioP90+tireModelAxleDelta || (frontTraction && !rearTraction) {
			limit.LimitedAxle = "front"
			limit.LimitedWheels = []string{"front_left", "front_right"}
		} else {
			limit.LimitedAxle = "driven"
			limit.LimitedWheels = []string{"front_left", "front_right", "rear_left", "rear_right"}
		}
		limit.Reason = "tire_grip_power_slip"
	case frontCombined && rearCombined && math.Abs(front.LimitScore-rear.LimitScore) <= tireModelAxleDelta:
		limit.Type = "combined_limit"
		limit.LimitedAxle = "both"
		limit.LimitedWheels = []string{"front_left", "front_right", "rear_left", "rear_right"}
		limit.PrimaryEvidence = "combined_slip_p90"
		limit.Reason = "tire_grip_four_wheel_combined"
	case frontLateral || rearLateral:
		limit.Type = "lateral_limit"
		limit.PrimaryEvidence = "slip_angle_p90"
		if front.LimitScore >= rear.LimitScore+tireModelAxleDelta || (frontLateral && !rearLateral) {
			limit.LimitedAxle = "front"
			limit.LimitedWheels = []string{"front_left", "front_right"}
		} else if rear.LimitScore >= front.LimitScore+tireModelAxleDelta || (rearLateral && !frontLateral) {
			limit.LimitedAxle = "rear"
			limit.LimitedWheels = []string{"rear_left", "rear_right"}
		} else {
			limit.LimitedAxle = "both"
			limit.LimitedWheels = []string{"front_left", "front_right", "rear_left", "rear_right"}
		}
		limit.Reason = "tire_grip_lateral_slip"
	case frontCombined || rearCombined:
		limit.Type = "balanced_near_limit"
		limit.PrimaryEvidence = "combined_slip_p90"
		if front.LimitScore >= rear.LimitScore+tireModelAxleDelta {
			limit.LimitedAxle = "front"
			limit.LimitedWheels = []string{"front_left", "front_right"}
		} else if rear.LimitScore >= front.LimitScore+tireModelAxleDelta {
			limit.LimitedAxle = "rear"
			limit.LimitedWheels = []string{"rear_left", "rear_right"}
		} else {
			limit.LimitedAxle = "both"
			limit.LimitedWheels = []string{"front_left", "front_right", "rear_left", "rear_right"}
		}
		limit.Reason = "tire_grip_near_limit"
	default:
		limit.Type = "no_limit_detected"
		limit.LimitedAxle = "none"
		limit.Reason = "tire_grip_no_dynamic_limit"
	}
	if quality.Status == "low_confidence" && limit.Type != "no_limit_detected" {
		limit.Confidence = quickConfidenceLow
	}
	if handbrake >= 0.20 && limit.Type == "braking_limit" {
		limit.Reason = "tire_grip_handbrake_slip"
	}
	return limit
}

func tireModelHints(diag TireModelDiagnostic) []TireModelHint {
	switch diag.LimitType {
	case "stationary":
		return []TireModelHint{{
			Code:      "observe",
			Severity:  "stable",
			Direction: "collect_moving_tire_load_samples",
			Reason:    diag.Explanation,
		}}
	case "no_dynamic_load":
		return []TireModelHint{{
			Code:      "observe",
			Severity:  "stable",
			Direction: "collect_moving_tire_load_samples",
			Reason:    diag.Explanation,
		}}
	case "front_limited":
		return []TireModelHint{{
			Code:      "front_axle_grip",
			Severity:  diag.FrontAxle.GripState,
			Direction: "improve_front_grip_or_reduce_entry_load",
			Reason:    diag.Explanation,
		}}
	case "rear_limited":
		return []TireModelHint{{
			Code:      "rear_axle_grip",
			Severity:  diag.RearAxle.GripState,
			Direction: "improve_rear_grip_or_smooth_rotation",
			Reason:    diag.Explanation,
		}}
	case "traction_limited":
		return []TireModelHint{{
			Code:      "driven_tire_traction",
			Severity:  "warning",
			Direction: "reduce_wheel_torque_or_improve_driven_tire_grip",
			Reason:    diag.Explanation,
		}}
	case "four_wheel_limited":
		return []TireModelHint{{
			Code:      "whole_car_grip",
			Severity:  "limit",
			Direction: "reduce_speed_or_increase_total_grip",
			Reason:    diag.Explanation,
		}}
	case "thermal_limited":
		return []TireModelHint{{
			Code:      "tire_temperature",
			Severity:  "warning",
			Direction: "bring_tires_back_to_temperature_window",
			Reason:    diag.Explanation,
		}}
	case "platform_limited":
		return []TireModelHint{{
			Code:      "platform_stability",
			Severity:  "warning",
			Direction: "restore_suspension_travel_and_platform_control",
			Reason:    diag.Explanation,
		}}
	default:
		return []TireModelHint{{
			Code:      "observe",
			Severity:  "stable",
			Direction: "collect_more_representative_corner_and_power_samples",
			Reason:    diag.Explanation,
		}}
	}
}
