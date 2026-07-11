package storage

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"fh6worker/internal/telemetry"
)

const (
	tireIssueEvalStepMS     = int64(300)
	tireIssueMergeGapMS     = int64(700)
	tireIssueMinDurationMS  = int64(700)
	tireIssueMinSegmentHits = 2
)

func defaultTireIssueAnalysis() TireIssueAnalysis {
	return TireIssueAnalysis{
		Status:    tireModelStatusNoData,
		UpdatedAt: nowText(),
		Segments:  []TireIssueSegment{},
		Groups:    []TireIssueGroup{},
		Warnings:  []string{},
	}
}

func BuildTireIssueAnalysis(samples []telemetry.NormalizedTelemetry) TireIssueAnalysis {
	analysis := defaultTireIssueAnalysis()
	if len(samples) == 0 {
		analysis.Warnings = append(analysis.Warnings, "tire_issue_no_samples")
		return analysis
	}
	ordered := sortedTelemetrySamples(samples)
	if len(ordered) == 0 {
		analysis.Warnings = append(analysis.Warnings, "tire_issue_no_samples")
		return analysis
	}
	analysis.Status = tireModelStatusReady
	analysis.SampleCount = len(ordered)
	analysis.WindowMS = ordered[len(ordered)-1].TimeMS - ordered[0].TimeMS
	if analysis.WindowMS < 0 {
		analysis.WindowMS = 0
	}

	raw := make([]TireIssueSegment, 0)
	lastEval := int64(-1 << 62)
	for i, sample := range ordered {
		if i != len(ordered)-1 && sample.TimeMS-lastEval < tireIssueEvalStepMS {
			continue
		}
		lastEval = sample.TimeMS
		window := tireIssueTrendWindowEndingAt(ordered, i)
		if len(window) < 2 {
			continue
		}
		diag := buildTireIssueWindowDiagnostic(window, ordered)
		raw = append(raw, tireIssueCandidates(diag, window)...)
	}
	segments := mergeTireIssueSegments(raw)
	analysis.Segments = segments
	analysis.Groups = groupTireIssueSegments(segments)
	analysis.SegmentCount = len(analysis.Segments)
	analysis.GroupCount = len(analysis.Groups)
	if len(analysis.Groups) == 0 {
		analysis.Warnings = append(analysis.Warnings, "tire_issue_no_groups")
	}
	return analysis
}

func buildTireIssueWindowDiagnostic(window []telemetry.NormalizedTelemetry, reference []telemetry.NormalizedTelemetry) TireModelDiagnostic {
	diag := TireModelDiagnostic{
		Status:      tireModelStatusReady,
		UpdatedAt:   nowText(),
		Confidence:  tireModelConfidence(len(window), tireIssueWindowMS(window)),
		GameMode:    quickSummarizeGameMode(window),
		PhaseDetail: BuildTirePhaseDiagnosticWithReference(window, reference),
		DataQuality: defaultTireDataQuality(),
		GripLimit:   defaultTireGripLimit(),
		LimitType:   "balanced",
		Warnings:    []string{},
		Evidence:    map[string]float64{},
		Vehicle:     quickVehicleSnapshot(window, &window[len(window)-1]),
	}
	diag.Phase = diag.PhaseDetail.CurrentPhase
	diag.SampleCount = len(window)
	diag.WindowMS = tireIssueWindowMS(window)
	fl := buildTireWheel("front_left", window, func(frame telemetry.NormalizedTelemetry) telemetry.NormalizedWheelTelemetry { return frame.WheelFL })
	fr := buildTireWheel("front_right", window, func(frame telemetry.NormalizedTelemetry) telemetry.NormalizedWheelTelemetry { return frame.WheelFR })
	rl := buildTireWheel("rear_left", window, func(frame telemetry.NormalizedTelemetry) telemetry.NormalizedWheelTelemetry { return frame.WheelRL })
	rr := buildTireWheel("rear_right", window, func(frame telemetry.NormalizedTelemetry) telemetry.NormalizedWheelTelemetry { return frame.WheelRR })
	diag.Wheels = []TireWheelDiagnostic{fl, fr, rl, rr}
	diag.FrontAxle = buildTireAxle("front", fl, fr)
	diag.RearAxle = buildTireAxle("rear", rl, rr)
	diag.LeftRight = buildTireSideBalance(fl, fr, rl, rr)
	diag.GForce = buildGForceDiagnostic(window)
	avgThrottle, avgBrake, avgHandBrake, avgSteer, avgSpeed := tireModelInputs(window)
	diag.Evidence = map[string]float64{
		"front_combined_slip_p90": diag.FrontAxle.CombinedSlipP90,
		"rear_combined_slip_p90":  diag.RearAxle.CombinedSlipP90,
		"front_slip_ratio_p90":    diag.FrontAxle.SlipRatioP90,
		"rear_slip_ratio_p90":     diag.RearAxle.SlipRatioP90,
		"front_slip_angle_p90":    diag.FrontAxle.SlipAngleP90,
		"rear_slip_angle_p90":     diag.RearAxle.SlipAngleP90,
		"front_tire_temp_avg":     diag.FrontAxle.TireTempAvg,
		"rear_tire_temp_avg":      diag.RearAxle.TireTempAvg,
		"front_suspension_max":    diag.FrontAxle.SuspensionTravelMax,
		"rear_suspension_max":     diag.RearAxle.SuspensionTravelMax,
		"left_right_slip_delta":   diag.LeftRight.Delta,
		"avg_throttle":            avgThrottle,
		"avg_brake":               avgBrake,
		"avg_handbrake":           avgHandBrake,
		"avg_steer":               avgSteer,
		"avg_speed_kmh":           avgSpeed,
		"peak_total_g":            diag.GForce.PeakTotalG,
	}
	dynamicWindow := tireModelDynamicSamples(window)
	diag.Evidence["dynamic_sample_count"] = float64(len(dynamicWindow))
	diag.DataQuality = buildTireDataQuality(window, dynamicWindow, &window[len(window)-1], diag)
	if tireModelCurrentStationary(&window[len(window)-1], avgSpeed, avgThrottle) {
		markTireModelStationary(&diag)
		return diag
	}
	appendTireModelRisks(&diag)
	if len(dynamicWindow) < tireModelMinDynamic {
		markTireModelNoDynamicLoad(&diag)
		return diag
	}
	limitDiag := buildTireLimitDiagnostic(dynamicWindow)
	diag.LimitType, diag.Summary, diag.Explanation = classifyTireLimit(limitDiag, avgThrottle, avgBrake, avgHandBrake)
	diag.GripLimit = buildTireGripLimit(limitDiag, diag.LeftRight, avgThrottle, avgBrake, avgHandBrake, diag.DataQuality)
	return diag
}

func tireIssueTrendWindowEndingAt(samples []telemetry.NormalizedTelemetry, end int) []telemetry.NormalizedTelemetry {
	if end < 0 || end >= len(samples) {
		return nil
	}
	latest := samples[end].TimeMS
	start := latest - tireModelTrendWindowMS
	first := end
	for first > 0 && samples[first-1].TimeMS >= start {
		first--
	}
	return samples[first : end+1]
}

func tireIssueWindowMS(samples []telemetry.NormalizedTelemetry) int64 {
	if len(samples) < 2 {
		return 0
	}
	out := samples[len(samples)-1].TimeMS - samples[0].TimeMS
	if out < 0 {
		return 0
	}
	return out
}

func tireIssueCandidates(diag TireModelDiagnostic, window []telemetry.NormalizedTelemetry) []TireIssueSegment {
	candidates := make([]TireIssueSegment, 0, 3)
	if len(window) == 0 {
		return candidates
	}
	ops := tireOperationTags(diag.PhaseDetail.Evidence)
	driftSource := ""
	if diag.PhaseDetail.CurrentPhase == "drift" || diag.PhaseDetail.CurrentPhase == "handbrake" || diag.PhaseDetail.SecondaryPhase == "drift" {
		driftSource = tireDriftSource(diag.PhaseDetail.Evidence)
	}
	if diag.GripLimit.Type != "" && diag.GripLimit.Type != "no_limit_detected" {
		candidates = append(candidates, newTireIssueSegment(
			tireIssueTypeFromGrip(diag.GripLimit.Type),
			diag.PhaseDetail.CurrentPhase,
			ops,
			driftSource,
			diag.GripLimit.Type,
			diag.GripLimit.LimitedAxle,
			diag.GripLimit.LimitedWheels,
			diag.DataQuality.Status,
			diag.GripLimit.Confidence,
			diag.GripLimit.Reason,
			mergeFloatMaps(diag.PhaseDetail.Evidence, diag.GripLimit.Evidence),
			window,
		))
	}
	for _, warning := range diag.Warnings {
		switch warning {
		case "thermal_risk":
			candidates = append(candidates, newTireIssueSegment("thermal_risk", diag.PhaseDetail.CurrentPhase, ops, driftSource, "risk", "all", nil, diag.DataQuality.Status, diag.DataQuality.Confidence, "thermal_risk", diag.Evidence, window))
		case "platform_risk":
			candidates = append(candidates, newTireIssueSegment("platform_risk", diag.PhaseDetail.CurrentPhase, ops, driftSource, "risk", tirePlatformRiskAxle(diag), nil, diag.DataQuality.Status, diag.DataQuality.Confidence, "platform_risk", diag.Evidence, window))
		case "left_right_imbalance":
			candidates = append(candidates, newTireIssueSegment("left_right_imbalance", diag.PhaseDetail.CurrentPhase, ops, driftSource, "risk", "left_right", nil, diag.DataQuality.Status, diag.DataQuality.Confidence, "left_right_imbalance", diag.Evidence, window))
		}
	}
	if len(candidates) == 0 && diag.DataQuality.Status == quickConfidenceInvalid {
		candidates = append(candidates, newTireIssueSegment("data_insufficient", diag.PhaseDetail.CurrentPhase, ops, driftSource, "data_quality", "none", nil, diag.DataQuality.Status, diag.DataQuality.Confidence, strings.Join(diag.DataQuality.Reasons, ","), diag.DataQuality.Evidence, window))
	}
	return candidates
}

func newTireIssueSegment(issueType, phase string, ops []string, driftSource, limitType, axle string, wheels []string, dataQuality, confidence, reason string, evidence map[string]float64, window []telemetry.NormalizedTelemetry) TireIssueSegment {
	startMS := window[0].TimeMS
	endMS := window[len(window)-1].TimeMS
	minSpeed, maxSpeed, avgSpeed := tireIssueSpeedStats(window)
	if phase == "" {
		phase = "unknown"
	}
	if axle == "" {
		axle = "none"
	}
	return TireIssueSegment{
		Type:          issueType,
		Phase:         phase,
		OperationTags: uniqueSortedStrings(ops),
		DriftSource:   driftSource,
		LimitType:     limitType,
		LimitedAxle:   axle,
		LimitedWheels: uniqueSortedStrings(wheels),
		StartMS:       startMS,
		EndMS:         endMS,
		DurationMS:    maxInt64(0, endMS-startMS),
		SampleCount:   len(window),
		SpeedMinKmh:   minSpeed,
		SpeedMaxKmh:   maxSpeed,
		SpeedAvgKmh:   avgSpeed,
		Confidence:    confidence,
		DataQuality:   dataQuality,
		RiskLevel:     tireIssueRiskLevel(issueType, confidence),
		Evidence:      cloneFloatMap(evidence),
		Reason:        reason,
	}
}

func tireOperationTags(e map[string]float64) []string {
	tags := make([]string, 0, 6)
	throttle := e["avg_throttle"]
	throttleDelta := e["throttle_delta"]
	brake := math.Max(e["avg_brake"], e["peak_brake"])
	steer := e["avg_steer"]
	steerDelta := e["steer_delta"]
	speedDelta := e["speed_delta_kmh"]
	if throttle >= 0.45 || throttleDelta > 0.12 {
		tags = append(tags, "throttle_on")
	} else if throttleDelta < -0.10 || throttle < 0.12 {
		tags = append(tags, "throttle_lift")
	} else if math.Abs(throttleDelta) <= 0.08 {
		tags = append(tags, "throttle_steady")
	}
	if brake >= 0.45 {
		tags = append(tags, "heavy_brake")
	} else if brake >= 0.12 {
		tags = append(tags, "light_brake")
	}
	if e["peak_handbrake"] >= 0.20 {
		tags = append(tags, "handbrake_active")
	}
	if math.Abs(steerDelta) > 0.10 {
		if steerDelta > 0 {
			tags = append(tags, "steer_increasing")
		} else {
			tags = append(tags, "steer_unwinding")
		}
	} else if steer >= 0.16 {
		tags = append(tags, "steer_holding")
	}
	if speedDelta > 2 {
		tags = append(tags, "speed_rising")
	} else if speedDelta < -2 {
		tags = append(tags, "speed_falling")
	}
	return uniqueSortedStrings(tags)
}

func tireDriftSource(e map[string]float64) string {
	rearCombined := e["rear_combined_slip"]
	rearRatio := e["rear_slip_ratio"]
	rearAngle := e["rear_slip_angle"]
	frontCombined := e["front_combined_slip"]
	if e["peak_handbrake"] >= 0.20 && rearCombined >= 0.55 {
		return "handbrake_initiated"
	}
	if e["avg_throttle"] >= 0.45 && (rearRatio >= 0.30 || rearCombined >= frontCombined+0.20) {
		return "power_oversteer"
	}
	if e["steer_sign_change"] > 0 && rearCombined >= 0.50 {
		return "scandinavian_flick"
	}
	if e["throttle_delta"] <= -0.10 && (rearAngle >= 0.40 || rearCombined >= frontCombined+0.18) {
		return "lift_off_oversteer"
	}
	return "unknown_oversteer"
}

func tireIssueTypeFromGrip(gripType string) string {
	switch gripType {
	case "lateral_limit":
		return "lateral_limit"
	case "traction_limit":
		return "traction_limit"
	case "braking_limit":
		return "braking_limit"
	case "combined_limit", "balanced_near_limit":
		return "combined_limit"
	default:
		return "combined_limit"
	}
}

func tirePlatformRiskAxle(diag TireModelDiagnostic) string {
	switch {
	case diag.FrontAxle.SuspensionTravelMax >= tireModelBottomOut && diag.RearAxle.SuspensionTravelMax >= tireModelBottomOut:
		return "all"
	case diag.FrontAxle.SuspensionTravelMax >= tireModelBottomOut:
		return "front"
	case diag.RearAxle.SuspensionTravelMax >= tireModelBottomOut:
		return "rear"
	default:
		return "none"
	}
}

func mergeTireIssueSegments(raw []TireIssueSegment) []TireIssueSegment {
	if len(raw) == 0 {
		return []TireIssueSegment{}
	}
	sort.SliceStable(raw, func(i, j int) bool {
		ki := tireIssueSegmentKey(raw[i])
		kj := tireIssueSegmentKey(raw[j])
		if ki == kj {
			return raw[i].StartMS < raw[j].StartMS
		}
		return ki < kj
	})
	merged := make([]TireIssueSegment, 0, len(raw))
	var current TireIssueSegment
	currentKey := ""
	hits := 0
	flush := func() {
		if current.Type == "" {
			return
		}
		current.DurationMS = maxInt64(0, current.EndMS-current.StartMS)
		if current.DurationMS >= tireIssueMinDurationMS && hits >= tireIssueMinSegmentHits {
			current.ID = fmt.Sprintf("seg-%03d", len(merged)+1)
			merged = append(merged, current)
		}
	}
	for _, seg := range raw {
		key := tireIssueSegmentKey(seg)
		if current.Type == "" {
			current = seg
			currentKey = key
			hits = 1
			continue
		}
		if key == currentKey && seg.StartMS <= current.EndMS+tireIssueMergeGapMS {
			current = mergeTireIssueSegment(current, seg)
			hits++
			continue
		}
		flush()
		current = seg
		currentKey = key
		hits = 1
	}
	flush()
	sort.SliceStable(merged, func(i, j int) bool { return merged[i].StartMS < merged[j].StartMS })
	return merged
}

func mergeTireIssueSegment(a, b TireIssueSegment) TireIssueSegment {
	if b.EndMS > a.EndMS {
		a.EndMS = b.EndMS
	}
	if b.StartMS < a.StartMS {
		a.StartMS = b.StartMS
	}
	a.DurationMS = maxInt64(0, a.EndMS-a.StartMS)
	a.SampleCount += b.SampleCount
	a.SpeedMinKmh = math.Min(a.SpeedMinKmh, b.SpeedMinKmh)
	a.SpeedMaxKmh = math.Max(a.SpeedMaxKmh, b.SpeedMaxKmh)
	a.SpeedAvgKmh = (a.SpeedAvgKmh + b.SpeedAvgKmh) / 2
	if confidenceRank(b.Confidence) > confidenceRank(a.Confidence) {
		a.Confidence = b.Confidence
	}
	if dataQualityRank(b.DataQuality) < dataQualityRank(a.DataQuality) {
		a.DataQuality = b.DataQuality
	}
	a.Evidence = mergeFloatMaps(a.Evidence, b.Evidence)
	return a
}

func groupTireIssueSegments(segments []TireIssueSegment) []TireIssueGroup {
	if len(segments) == 0 {
		return []TireIssueGroup{}
	}
	byKey := map[string]*TireIssueGroup{}
	for _, seg := range segments {
		key := tireIssueSegmentKey(seg)
		group := byKey[key]
		if group == nil {
			group = &TireIssueGroup{
				ID:                     fmt.Sprintf("grp-%03d", len(byKey)+1),
				Type:                   seg.Type,
				Phase:                  seg.Phase,
				OperationTags:          append([]string(nil), seg.OperationTags...),
				DriftSource:            seg.DriftSource,
				LimitType:              seg.LimitType,
				LimitedAxle:            seg.LimitedAxle,
				LimitedWheels:          append([]string(nil), seg.LimitedWheels...),
				SpeedMinKmh:            seg.SpeedMinKmh,
				SpeedMaxKmh:            seg.SpeedMaxKmh,
				Confidence:             seg.Confidence,
				DataQuality:            seg.DataQuality,
				RiskLevel:              seg.RiskLevel,
				RepresentativeEvidence: cloneFloatMap(seg.Evidence),
				Reason:                 seg.Reason,
			}
			byKey[key] = group
		}
		group.Count++
		group.TotalDurationMS += seg.DurationMS
		group.SpeedMinKmh = math.Min(group.SpeedMinKmh, seg.SpeedMinKmh)
		group.SpeedMaxKmh = math.Max(group.SpeedMaxKmh, seg.SpeedMaxKmh)
		group.SpeedAvgKmh += seg.SpeedAvgKmh
		group.SegmentIDs = append(group.SegmentIDs, seg.ID)
		group.LimitedWheels = uniqueSortedStrings(append(group.LimitedWheels, seg.LimitedWheels...))
		if confidenceRank(seg.Confidence) > confidenceRank(group.Confidence) {
			group.Confidence = seg.Confidence
			group.RepresentativeEvidence = cloneFloatMap(seg.Evidence)
			group.Reason = seg.Reason
		}
		if dataQualityRank(seg.DataQuality) < dataQualityRank(group.DataQuality) {
			group.DataQuality = seg.DataQuality
		}
	}
	groups := make([]TireIssueGroup, 0, len(byKey))
	for _, group := range byKey {
		if group.Count > 0 {
			group.SpeedAvgKmh /= float64(group.Count)
		}
		groups = append(groups, *group)
	}
	sort.SliceStable(groups, func(i, j int) bool {
		if groups[i].RiskLevel == groups[j].RiskLevel {
			if groups[i].TotalDurationMS == groups[j].TotalDurationMS {
				return groups[i].Count > groups[j].Count
			}
			return groups[i].TotalDurationMS > groups[j].TotalDurationMS
		}
		return riskRank(groups[i].RiskLevel) > riskRank(groups[j].RiskLevel)
	})
	for i := range groups {
		groups[i].ID = fmt.Sprintf("grp-%03d", i+1)
	}
	return groups
}

func tireIssueSegmentKey(seg TireIssueSegment) string {
	return strings.Join([]string{
		seg.Type,
		seg.Phase,
		strings.Join(uniqueSortedStrings(seg.OperationTags), "+"),
		seg.DriftSource,
		seg.LimitType,
		seg.LimitedAxle,
		strings.Join(uniqueSortedStrings(seg.LimitedWheels), "+"),
	}, "|")
}

func tireIssueSpeedStats(samples []telemetry.NormalizedTelemetry) (minSpeed, maxSpeed, avgSpeed float64) {
	if len(samples) == 0 {
		return 0, 0, 0
	}
	minSpeed = math.MaxFloat64
	for _, sample := range samples {
		minSpeed = math.Min(minSpeed, sample.SpeedKmh)
		maxSpeed = math.Max(maxSpeed, sample.SpeedKmh)
		avgSpeed += sample.SpeedKmh
	}
	avgSpeed /= float64(len(samples))
	return minSpeed, maxSpeed, avgSpeed
}

func tireIssueRiskLevel(issueType, confidence string) string {
	switch issueType {
	case "lateral_limit", "traction_limit", "braking_limit", "combined_limit":
		if confidence == quickConfidenceHigh || confidence == quickConfidenceMedium {
			return "high"
		}
		return "medium"
	case "platform_risk", "thermal_risk":
		return "medium"
	case "left_right_imbalance", "data_insufficient":
		return "low"
	default:
		return "low"
	}
}

func mergeFloatMaps(a, b map[string]float64) map[string]float64 {
	out := cloneFloatMap(a)
	for key, value := range b {
		if existing, ok := out[key]; !ok || math.Abs(value) > math.Abs(existing) {
			out[key] = value
		}
	}
	return out
}

func uniqueSortedStrings(values []string) []string {
	if len(values) == 0 {
		return []string{}
	}
	seen := map[string]struct{}{}
	out := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		out = append(out, value)
	}
	sort.Strings(out)
	return out
}

func confidenceRank(value string) int {
	switch value {
	case quickConfidenceHigh:
		return 3
	case quickConfidenceMedium:
		return 2
	case quickConfidenceLow:
		return 1
	default:
		return 0
	}
}

func dataQualityRank(value string) int {
	switch value {
	case "valid":
		return 3
	case "low_confidence":
		return 2
	case "invalid":
		return 1
	default:
		return 0
	}
}

func riskRank(value string) int {
	switch value {
	case "high":
		return 3
	case "medium":
		return 2
	case "low":
		return 1
	default:
		return 0
	}
}

func maxInt64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}
