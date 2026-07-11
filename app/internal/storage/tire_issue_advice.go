package storage

import (
	"fmt"
	"sort"
	"strings"

	"fh6worker/internal/telemetry"
)

const (
	tireAdviceStatusNoData = "no_data"
	tireAdviceStatusReady  = "ready"
	tireAdviceStatusEmpty  = "no_actionable_issues"

	tireAdviceLayerPrimary     = "primary"
	tireAdviceLayerAlternative = "alternative"
	tireAdviceLayerCheck       = "check"
	tireAdviceLayerObserve     = "observe"
)

func defaultTireIssueAdvice() TireIssueAdvice {
	return TireIssueAdvice{
		Status:          tireAdviceStatusNoData,
		UpdatedAt:       nowText(),
		Confidence:      quickConfidenceLow,
		PriorityActions: []TireIssueAdviceAction{},
		Groups:          []TireIssueAdviceGroup{},
		Warnings:        []string{},
	}
}

func BuildTireIssueAdvice(samples []telemetry.NormalizedTelemetry) TireIssueAdvice {
	return BuildTireIssueAdviceFromAnalysis(BuildTireIssueAnalysis(samples))
}

func BuildTireIssueAdviceFromAnalysis(analysis TireIssueAnalysis) TireIssueAdvice {
	advice := defaultTireIssueAdvice()
	advice.BasedOnIssueUpdatedAt = analysis.UpdatedAt
	advice.IssueGroupCount = len(analysis.Groups)
	if len(analysis.Warnings) > 0 {
		advice.Warnings = append(advice.Warnings, analysis.Warnings...)
	}
	if analysis.Status == tireModelStatusNoData || analysis.Status == "" {
		advice.Warnings = append(advice.Warnings, "tire_advice_no_issue_analysis")
		return advice
	}
	if len(analysis.Groups) == 0 {
		advice.Status = tireAdviceStatusEmpty
		advice.Warnings = append(advice.Warnings, "tire_advice_no_issue_groups")
		return advice
	}

	advice.Status = tireAdviceStatusReady
	bestConfidence := quickConfidenceLow
	for _, issueGroup := range analysis.Groups {
		groupAdvice := buildTireIssueAdviceGroup(issueGroup)
		advice.Groups = append(advice.Groups, groupAdvice)
		if confidenceRank(groupAdvice.Confidence) > confidenceRank(bestConfidence) {
			bestConfidence = groupAdvice.Confidence
		}
	}
	sort.SliceStable(advice.Groups, func(i, j int) bool {
		if advice.Groups[i].Priority == advice.Groups[j].Priority {
			return advice.Groups[i].IssueGroupID < advice.Groups[j].IssueGroupID
		}
		return advice.Groups[i].Priority < advice.Groups[j].Priority
	})
	advice.Confidence = bestConfidence
	advice.PriorityActions = tireIssuePriorityActions(advice.Groups, 3)
	return advice
}

func buildTireIssueAdviceGroup(group TireIssueGroup) TireIssueAdviceGroup {
	out := TireIssueAdviceGroup{
		IssueGroupID:  group.ID,
		IssueType:     group.Type,
		Phase:         group.Phase,
		OperationTags: append([]string(nil), group.OperationTags...),
		LimitedAxle:   normalizeTireAdviceAxle(group.LimitedAxle),
		DriftSource:   group.DriftSource,
		PrimaryCause:  tireIssuePrimaryCause(group),
		ShouldTune:    tireIssueShouldTune(group),
		Priority:      tireIssueAdvicePriority(group),
		Confidence:    group.Confidence,
		Evidence:      cloneFloatMap(group.RepresentativeEvidence),
		Actions:       []TireIssueAdviceAction{},
	}
	out.Actions = tireIssueAdviceActions(group, out.ShouldTune)
	if len(out.Actions) > 2 {
		out.Actions = out.Actions[:2]
	}
	return out
}

func tireIssueShouldTune(group TireIssueGroup) bool {
	if group.Type == "data_insufficient" || group.DataQuality == "invalid" || group.Confidence == quickConfidenceInvalid {
		return false
	}
	if group.Confidence == quickConfidenceLow {
		return false
	}
	if group.DriftSource == "handbrake_initiated" || group.DriftSource == "scandinavian_flick" {
		return false
	}
	return true
}

func tireIssuePrimaryCause(group TireIssueGroup) string {
	if group.Type == "data_insufficient" || group.DataQuality == "invalid" {
		return "data_not_reliable"
	}
	if group.DriftSource == "handbrake_initiated" {
		return "driver_handbrake_drift"
	}
	if group.DriftSource == "scandinavian_flick" {
		return "driver_weight_transfer_drift"
	}
	switch group.Type {
	case "lateral_limit":
		if group.LimitedAxle == "front" {
			if group.Phase == "high_speed_corner" {
				return "front_high_speed_lateral_limit"
			}
			return "front_mechanical_lateral_limit"
		}
		if group.LimitedAxle == "rear" {
			return "rear_lateral_stability_limit"
		}
		return "four_wheel_lateral_limit"
	case "traction_limit":
		if containsStringValue(group.OperationTags, "throttle_on") || group.Phase == "corner_exit" || group.Phase == "launch" || group.Phase == "straight_power" {
			return "drive_torque_exceeds_tire_grip"
		}
		return "driven_wheel_longitudinal_slip"
	case "braking_limit":
		if group.LimitedAxle == "rear" {
			return "rear_brake_or_decel_instability"
		}
		return "front_brake_overload"
	case "combined_limit":
		return "combined_longitudinal_lateral_overload"
	case "platform_risk":
		return "platform_travel_or_load_risk"
	case "thermal_risk":
		return "tire_temperature_risk"
	case "left_right_imbalance":
		return "left_right_signal_or_load_imbalance"
	default:
		return "unknown_tire_issue"
	}
}

func tireIssueAdviceActions(group TireIssueGroup, shouldTune bool) []TireIssueAdviceAction {
	if !shouldTune {
		return []TireIssueAdviceAction{tireIssueObservationAction(group)}
	}
	switch group.Type {
	case "lateral_limit":
		return tireLateralAdviceActions(group)
	case "traction_limit":
		return tireTractionAdviceActions(group)
	case "braking_limit":
		return tireBrakingAdviceActions(group)
	case "combined_limit":
		return tireCombinedAdviceActions(group)
	case "platform_risk":
		return []TireIssueAdviceAction{
			newTireIssueAdviceAction(group, tireAdviceLayerPrimary, "spring_damping", normalizeTireAdviceAxle(group.LimitedAxle), "check_platform", []string{"frontSpring", "rearSpring", "frontRebound", "rearRebound", "frontBump", "rearBump", "frontRideHeight", "rearRideHeight"}, "platform_risk_check_travel_and_damping", []string{"front_suspension_max", "rear_suspension_max", "front_combined_slip_p90", "rear_combined_slip_p90"}, group.Confidence, true),
		}
	case "thermal_risk":
		return []TireIssueAdviceAction{
			newTireIssueAdviceAction(group, tireAdviceLayerPrimary, "tire_pressure", normalizeThermalScope(group.LimitedAxle), "check_temperature_window", []string{"frontTirePressure", "rearTirePressure"}, "thermal_risk_check_pressure_and_slip", []string{"front_tire_temp_avg", "rear_tire_temp_avg", "front_combined_slip_p90", "rear_combined_slip_p90"}, group.Confidence, true),
		}
	case "left_right_imbalance":
		return []TireIssueAdviceAction{
			newTireIssueAdviceAction(group, tireAdviceLayerCheck, "data_quality", "left_right", "check_left_right", nil, "left_right_imbalance_verify_route_and_sensor", []string{"left_right_slip_delta"}, group.Confidence, false),
		}
	default:
		return []TireIssueAdviceAction{tireIssueObservationAction(group)}
	}
}

func tireLateralAdviceActions(group TireIssueGroup) []TireIssueAdviceAction {
	axle := normalizeTireAdviceAxle(group.LimitedAxle)
	if axle == "front" {
		if group.Phase == "high_speed_corner" {
			return []TireIssueAdviceAction{
				newTireIssueAdviceAction(group, tireAdviceLayerPrimary, "aero_platform", "front", "increase_front_high_speed_support", []string{"frontAero", "frontRideHeight", "frontSpring", "frontRebound", "frontBump"}, "front_high_speed_lateral_limit_prioritize_platform", []string{"front_combined_slip_p90", "front_slip_angle_p90", "avg_speed_kmh", "front_suspension_max"}, group.Confidence, true),
				newTireIssueAdviceAction(group, tireAdviceLayerAlternative, "antiroll", "front", "increase_front_mechanical_grip", []string{"frontArb", "frontCamber", "frontTirePressure"}, "front_high_speed_lateral_limit_secondary_mechanical_grip", []string{"front_slip_angle_p90", "front_tire_temp_avg"}, group.Confidence, true),
			}
		}
		return []TireIssueAdviceAction{
			newTireIssueAdviceAction(group, tireAdviceLayerPrimary, "antiroll", "front", "increase_front_mechanical_grip", []string{"frontArb", "frontCamber", "frontToe", "frontTirePressure"}, "front_lateral_limit_prioritize_mechanical_grip", []string{"front_combined_slip_p90", "front_slip_angle_p90", "avg_steer", "avg_speed_kmh"}, group.Confidence, true),
			newTireIssueAdviceAction(group, tireAdviceLayerAlternative, "alignment", "front", "check_front_contact_patch", []string{"frontCamber", "frontToe", "caster"}, "front_lateral_limit_verify_alignment", []string{"front_slip_angle_p90", "front_tire_temp_avg"}, group.Confidence, true),
		}
	}
	if axle == "rear" {
		return []TireIssueAdviceAction{
			newTireIssueAdviceAction(group, tireAdviceLayerPrimary, "spring_damping", "rear", "increase_rear_stability", []string{"rearArb", "rearSpring", "rearRebound", "rearBump", "rearTirePressure"}, "rear_lateral_limit_prioritize_stability", []string{"rear_combined_slip_p90", "rear_slip_angle_p90", "avg_speed_kmh"}, group.Confidence, true),
			newTireIssueAdviceAction(group, tireAdviceLayerAlternative, "alignment", "rear", "check_rear_contact_patch", []string{"rearCamber", "rearToe"}, "rear_lateral_limit_verify_alignment", []string{"rear_slip_angle_p90", "rear_tire_temp_avg"}, group.Confidence, true),
		}
	}
	return []TireIssueAdviceAction{
		newTireIssueAdviceAction(group, tireAdviceLayerPrimary, "spring_damping", "all", "reduce_four_wheel_lateral_load", []string{"frontArb", "rearArb", "frontSpring", "rearSpring", "frontRebound", "rearRebound"}, "four_wheel_lateral_limit_reduce_platform_load", []string{"front_combined_slip_p90", "rear_combined_slip_p90", "peak_total_g"}, group.Confidence, true),
		newTireIssueAdviceAction(group, tireAdviceLayerAlternative, "driver_input", "all", "reduce_overlap_input", nil, "four_wheel_lateral_limit_verify_driver_overlap", []string{"avg_throttle", "avg_brake", "avg_steer"}, group.Confidence, false),
	}
}

func tireTractionAdviceActions(group TireIssueGroup) []TireIssueAdviceAction {
	axle := normalizeTireAdviceAxle(group.LimitedAxle)
	if axle == "front" {
		return []TireIssueAdviceAction{
			newTireIssueAdviceAction(group, tireAdviceLayerPrimary, "differential", "front", "reduce_drive_lock", []string{"frontDiffAccel"}, "front_traction_limit_reduce_accel_diff", []string{"front_slip_ratio_p90", "avg_throttle", "avg_speed_kmh"}, group.Confidence, true),
			newTireIssueAdviceAction(group, tireAdviceLayerAlternative, "gearing", "front", "reduce_wheel_torque", []string{"finalDrive", "gear1", "gear2", "gear3"}, "front_traction_limit_check_gearing", []string{"front_slip_ratio_p90", "avg_throttle"}, group.Confidence, true),
		}
	}
	if axle == "rear" {
		return []TireIssueAdviceAction{
			newTireIssueAdviceAction(group, tireAdviceLayerPrimary, "differential", "rear", "reduce_drive_lock", []string{"rearDiffAccel"}, "rear_traction_limit_reduce_accel_diff", []string{"rear_slip_ratio_p90", "avg_throttle", "avg_speed_kmh"}, group.Confidence, true),
			newTireIssueAdviceAction(group, tireAdviceLayerAlternative, "gearing", "rear", "reduce_wheel_torque", []string{"finalDrive", "gear1", "gear2", "gear3"}, "rear_traction_limit_check_gearing", []string{"rear_slip_ratio_p90", "avg_throttle"}, group.Confidence, true),
		}
	}
	return []TireIssueAdviceAction{
		newTireIssueAdviceAction(group, tireAdviceLayerPrimary, "differential", "driven", "reduce_drive_lock", []string{"frontDiffAccel", "rearDiffAccel", "centerDiffBalance"}, "driven_traction_limit_balance_diff", []string{"front_slip_ratio_p90", "rear_slip_ratio_p90", "avg_throttle"}, group.Confidence, true),
		newTireIssueAdviceAction(group, tireAdviceLayerAlternative, "gearing", "driven", "reduce_wheel_torque", []string{"finalDrive", "gear1", "gear2", "gear3"}, "driven_traction_limit_check_gearing", []string{"front_slip_ratio_p90", "rear_slip_ratio_p90", "avg_throttle"}, group.Confidence, true),
	}
}

func tireBrakingAdviceActions(group TireIssueGroup) []TireIssueAdviceAction {
	axle := normalizeTireAdviceAxle(group.LimitedAxle)
	if axle == "rear" {
		return []TireIssueAdviceAction{
			newTireIssueAdviceAction(group, tireAdviceLayerPrimary, "brake", "rear", "move_brake_balance_forward", []string{"brakeBalance", "brakePressure"}, "rear_braking_limit_move_balance_forward", []string{"rear_slip_ratio_p90", "avg_brake", "avg_speed_kmh"}, group.Confidence, true),
			newTireIssueAdviceAction(group, tireAdviceLayerAlternative, "differential", "rear", "check_decel_lock", []string{"rearDiffDecel"}, "rear_braking_limit_check_decel_diff", []string{"rear_slip_ratio_p90", "avg_brake"}, group.Confidence, true),
		}
	}
	return []TireIssueAdviceAction{
		newTireIssueAdviceAction(group, tireAdviceLayerPrimary, "brake", "front", "move_brake_balance_rearward", []string{"brakeBalance", "brakePressure"}, "front_braking_limit_move_balance_rearward", []string{"front_slip_ratio_p90", "avg_brake", "avg_speed_kmh"}, group.Confidence, true),
		newTireIssueAdviceAction(group, tireAdviceLayerAlternative, "spring_damping", "front", "check_front_brake_platform", []string{"frontSpring", "frontRebound", "frontBump", "frontRideHeight"}, "front_braking_limit_check_platform", []string{"front_suspension_max", "front_combined_slip_p90"}, group.Confidence, true),
	}
}

func tireCombinedAdviceActions(group TireIssueGroup) []TireIssueAdviceAction {
	if containsStringValue(group.OperationTags, "heavy_brake") || containsStringValue(group.OperationTags, "light_brake") {
		return []TireIssueAdviceAction{
			newTireIssueAdviceAction(group, tireAdviceLayerPrimary, "driver_input", "all", "reduce_overlap_input", nil, "combined_limit_trail_brake_reduce_overlap", []string{"avg_brake", "avg_steer", "front_combined_slip_p90", "rear_combined_slip_p90"}, group.Confidence, false),
			newTireIssueAdviceAction(group, tireAdviceLayerAlternative, "brake", normalizeTireAdviceAxle(group.LimitedAxle), "check_brake_balance", []string{"brakeBalance", "brakePressure"}, "combined_limit_verify_brake_balance", []string{"avg_brake", "front_slip_ratio_p90", "rear_slip_ratio_p90"}, group.Confidence, true),
		}
	}
	if containsStringValue(group.OperationTags, "throttle_on") {
		return []TireIssueAdviceAction{
			newTireIssueAdviceAction(group, tireAdviceLayerPrimary, "driver_input", "all", "reduce_overlap_input", nil, "combined_limit_corner_exit_reduce_overlap", []string{"avg_throttle", "avg_steer", "front_combined_slip_p90", "rear_combined_slip_p90"}, group.Confidence, false),
			newTireIssueAdviceAction(group, tireAdviceLayerAlternative, "differential", normalizeTireAdviceAxle(group.LimitedAxle), "reduce_drive_lock", []string{"frontDiffAccel", "rearDiffAccel", "centerDiffBalance"}, "combined_limit_verify_drive_lock", []string{"front_slip_ratio_p90", "rear_slip_ratio_p90", "avg_throttle"}, group.Confidence, true),
		}
	}
	return []TireIssueAdviceAction{
		newTireIssueAdviceAction(group, tireAdviceLayerPrimary, "driver_input", "all", "reduce_overlap_input", nil, "combined_limit_reduce_combined_load", []string{"avg_throttle", "avg_brake", "avg_steer"}, group.Confidence, false),
		newTireIssueAdviceAction(group, tireAdviceLayerAlternative, "spring_damping", "all", "check_platform", []string{"frontArb", "rearArb", "frontSpring", "rearSpring", "frontRebound", "rearRebound"}, "combined_limit_verify_platform", []string{"front_combined_slip_p90", "rear_combined_slip_p90", "peak_total_g"}, group.Confidence, true),
	}
}

func tireIssueObservationAction(group TireIssueGroup) TireIssueAdviceAction {
	switch {
	case group.Type == "data_insufficient" || group.DataQuality == "invalid" || group.Confidence == quickConfidenceInvalid:
		return newTireIssueAdviceAction(group, tireAdviceLayerObserve, "data_quality", "none", "continue_sampling", nil, "data_quality_continue_sampling", []string{"dynamic_sample_count", "avg_speed_kmh", "peak_total_g"}, quickConfidenceLow, false)
	case group.DriftSource == "handbrake_initiated":
		return newTireIssueAdviceAction(group, tireAdviceLayerObserve, "driver_input", "rear", "avoid_tuning", nil, "handbrake_drift_driver_behavior", []string{"avg_handbrake", "rear_combined_slip_p90", "avg_speed_kmh"}, group.Confidence, false)
	case group.DriftSource == "scandinavian_flick":
		return newTireIssueAdviceAction(group, tireAdviceLayerObserve, "driver_input", "rear", "avoid_tuning", nil, "scandinavian_flick_driver_behavior", []string{"steer_sign_change", "rear_combined_slip_p90"}, group.Confidence, false)
	case group.Confidence == quickConfidenceLow:
		return newTireIssueAdviceAction(group, tireAdviceLayerCheck, "data_quality", normalizeTireAdviceAxle(group.LimitedAxle), "continue_sampling", nil, "low_confidence_verify_before_tuning", []string{"avg_speed_kmh", "front_combined_slip_p90", "rear_combined_slip_p90"}, quickConfidenceLow, false)
	default:
		return newTireIssueAdviceAction(group, tireAdviceLayerObserve, "data_quality", normalizeTireAdviceAxle(group.LimitedAxle), "continue_sampling", nil, "unknown_issue_continue_sampling", []string{"avg_speed_kmh"}, quickConfidenceLow, false)
	}
}

func newTireIssueAdviceAction(group TireIssueGroup, layer, category, scope, direction string, fields []string, rationale string, verify []string, confidence string, tuneRecommended bool) TireIssueAdviceAction {
	if confidence == "" {
		confidence = quickConfidenceLow
	}
	scope = normalizeTireAdviceAxle(scope)
	return TireIssueAdviceAction{
		ID:              fmt.Sprintf("%s-%s-%s-%s", group.ID, layer, category, direction),
		IssueGroupID:    group.ID,
		Layer:           layer,
		Category:        category,
		Scope:           scope,
		Direction:       direction,
		RelatedFields:   uniqueSortedStrings(fields),
		Rationale:       rationale,
		VerifyEvidence:  uniqueSortedStrings(verify),
		Confidence:      confidence,
		MissingInputs:   []string{},
		TuneRecommended: tuneRecommended,
	}
}

func tireIssuePriorityActions(groups []TireIssueAdviceGroup, limit int) []TireIssueAdviceAction {
	type scoredAction struct {
		action TireIssueAdviceAction
		score  int
	}
	scored := make([]scoredAction, 0)
	for _, group := range groups {
		for _, action := range group.Actions {
			score := group.Priority + tireAdviceLayerPriority(action.Layer)
			if !action.TuneRecommended {
				score += 12
			}
			scored = append(scored, scoredAction{action: action, score: score})
		}
	}
	sort.SliceStable(scored, func(i, j int) bool {
		if scored[i].score == scored[j].score {
			return scored[i].action.ID < scored[j].action.ID
		}
		return scored[i].score < scored[j].score
	})
	seen := map[string]struct{}{}
	out := make([]TireIssueAdviceAction, 0, limit)
	for _, item := range scored {
		key := strings.Join([]string{item.action.Category, item.action.Scope, item.action.Direction}, "|")
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, item.action)
		if len(out) >= limit {
			break
		}
	}
	return out
}

func tireIssueAdvicePriority(group TireIssueGroup) int {
	priority := 60
	switch group.RiskLevel {
	case "high":
		priority = 10
	case "medium":
		priority = 30
	case "low":
		priority = 50
	}
	priority -= confidenceRank(group.Confidence) * 2
	switch group.Type {
	case "data_insufficient", "left_right_imbalance":
		priority += 15
	case "platform_risk", "thermal_risk":
		priority += 5
	}
	if group.DriftSource == "handbrake_initiated" || group.DriftSource == "scandinavian_flick" {
		priority += 20
	}
	if priority < 1 {
		return 1
	}
	return priority
}

func tireAdviceLayerPriority(layer string) int {
	switch layer {
	case tireAdviceLayerPrimary:
		return 0
	case tireAdviceLayerAlternative:
		return 5
	case tireAdviceLayerCheck:
		return 8
	default:
		return 12
	}
}

func normalizeTireAdviceAxle(value string) string {
	switch strings.TrimSpace(value) {
	case "front", "rear", "all", "driven", "left_right", "none":
		return strings.TrimSpace(value)
	case "":
		return "none"
	default:
		return strings.TrimSpace(value)
	}
}

func normalizeThermalScope(axle string) string {
	scope := normalizeTireAdviceAxle(axle)
	if scope == "none" || scope == "left_right" {
		return "all"
	}
	return scope
}
