package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
)

const (
	tunePlanStatusReady        = "ready"
	tunePlanStatusNoActions    = "no_actions"
	tunePlanStatusUnbound      = "profile_unbound"
	tunePlanStatusMismatch     = "vehicle_mismatch"
	tunePlanStatusCannotVerify = "cannot_verify_vehicle"
)

func (s *Store) GetTunePlanDraft(sessionID int64) (*TunePlanDraft, error) {
	session, err := s.GetTelemetrySession(sessionID)
	if err != nil {
		return nil, err
	}
	draft := &TunePlanDraft{SessionID: sessionID, Status: tunePlanStatusUnbound, Summary: "bind_tune_profile_first"}
	if session.TuneProfileID == nil {
		return draft, nil
	}
	draft.TuneProfileID = session.TuneProfileID
	profile, err := s.GetTuneProfile(*session.TuneProfileID)
	if err != nil {
		return nil, err
	}
	status := tunePlanVehicleStatus(*session, *profile)
	summary, summaryErr := s.GetSessionIssueSummary(sessionID)
	if decision, err := s.GetRoadTuningDecision(sessionID); err == nil && decision != nil && len(decision.Actions) > 0 {
		draft.Status = status
		draft.Summary = decision.Reason
		for index, action := range decision.Actions {
			draft.Actions = append(draft.Actions, draftActionsForRoadDecisionAction(index, action, profile, status)...)
		}
		if summaryErr == nil && summary != nil {
			appendGearPowerDraftActions(draft, summary.GearPower, profile, status)
		}
		s.applyRetestGuardToDraft(draft, *session, profile, status)
		finalizeTunePlanDraft(draft)
		if status == tunePlanStatusReady && len(draft.Actions) == 0 {
			draft.Status = tunePlanStatusNoActions
			draft.Summary = "no_applicable_actions"
		}
		return draft, nil
	}
	if summaryErr != nil {
		return nil, summaryErr
	}
	draft.Status = status
	draft.Summary = summary.WholeCarPlan.Summary
	draft.Conflicts = append([]TuningConflict(nil), summary.WholeCarPlan.Conflicts...)
	for _, action := range summary.WholeCarPlan.Actions {
		draft.Actions = append(draft.Actions, draftActionsForWholeCarAdjustment(action, profile, status)...)
	}
	appendGearPowerDraftActions(draft, summary.GearPower, profile, status)
	s.applyRetestGuardToDraft(draft, *session, profile, status)
	finalizeTunePlanDraft(draft)
	if status == tunePlanStatusReady && len(draft.Actions) == 0 {
		draft.Status = tunePlanStatusNoActions
		draft.Summary = "no_applicable_actions"
	}
	return draft, nil
}

func appendGearPowerDraftActions(draft *TunePlanDraft, gear GearPowerDiagnostic, profile *TuneProfile, status string) {
	if draft == nil || len(gear.RecommendedActions) == 0 {
		return
	}
	existing := map[string]struct{}{}
	for _, action := range draft.Actions {
		existing[action.FieldKey+"|"+action.Direction] = struct{}{}
	}
	for _, action := range gear.RecommendedActions {
		adjustment := WholeCarAdjustment{
			Priority:   len(draft.Actions),
			Family:     "gearing_acceleration",
			Source:     "gear_power_diagnostic",
			Confidence: gear.Confidence,
			Category:   action.Category,
			Item:       action.Item,
			Direction:  action.Direction,
			Amount:     action.Amount,
			Reason:     action.Reason,
			Evidence:   cloneEvidenceMap(gear.Evidence),
		}
		for _, candidate := range draftActionsForWholeCarAdjustment(adjustment, profile, status) {
			key := candidate.FieldKey + "|" + candidate.Direction
			if _, ok := existing[key]; ok {
				continue
			}
			candidate.ID = fmt.Sprintf("gearpower:%s", candidate.ID)
			draft.Actions = append(draft.Actions, candidate)
			existing[key] = struct{}{}
		}
	}
}

func draftActionsForRoadDecisionAction(priority int, action RoadTuningDecisionAction, profile *TuneProfile, status string) []TunePlanDraftAction {
	adjustment := WholeCarAdjustment{
		Priority:   priority,
		Family:     action.Family,
		Source:     action.Source,
		Confidence: action.Confidence,
		Category:   action.Category,
		Item:       action.Item,
		Direction:  action.Direction,
		Amount:     action.Amount,
		Reason:     action.Reason,
		Evidence:   action.Evidence,
	}
	drafts := draftActionsForWholeCarAdjustment(adjustment, profile, status)
	for i := range drafts {
		drafts[i].ID = fmt.Sprintf("road:%s:%s", action.ID, drafts[i].FieldKey)
		if !action.CanAutoApply {
			drafts[i].CanApply = false
			if action.BlockedReason != "" {
				drafts[i].BlockedReason = action.BlockedReason
			} else {
				drafts[i].BlockedReason = "manual_review_required"
			}
		}
	}
	return drafts
}

func (s *Store) ApplyTunePlanDraft(input TunePlanApplyInput) (*TunePlanApplyResult, error) {
	if input.SessionID <= 0 {
		return nil, errors.New("session id is required")
	}
	if len(input.SelectedActionIDs) == 0 {
		return nil, errors.New("at least one tune plan action is required")
	}
	session, err := s.GetTelemetrySession(input.SessionID)
	if err != nil {
		return nil, err
	}
	if session.TuneProfileID == nil {
		return nil, errors.New("telemetry session has no bound tune profile")
	}
	profile, err := s.GetTuneProfile(*session.TuneProfileID)
	if err != nil {
		return nil, err
	}
	if status := tunePlanVehicleStatus(*session, *profile); status != tunePlanStatusReady {
		return nil, fmt.Errorf("cannot apply tune plan: %s", status)
	}
	draft, err := s.GetTunePlanDraft(input.SessionID)
	if err != nil {
		return nil, err
	}
	actionsByID := map[string]TunePlanDraftAction{}
	for _, action := range draft.Actions {
		actionsByID[action.ID] = action
	}
	selected := make([]TunePlanDraftAction, 0, len(input.SelectedActionIDs))
	inputProfile := profile.ToInput()
	changedSet := map[string]bool{}
	for _, id := range input.SelectedActionIDs {
		action, ok := actionsByID[id]
		if !ok {
			return nil, fmt.Errorf("tune plan action %q is not valid for this session", id)
		}
		if !action.CanApply || action.TargetValue == nil {
			return nil, fmt.Errorf("tune plan action %q cannot be applied: %s", id, action.BlockedReason)
		}
		if !setTuneProfileInputFloat(&inputProfile, action.FieldKey, *action.TargetValue) {
			return nil, fmt.Errorf("tune plan action %q targets an unsupported field", id)
		}
		changedSet[action.FieldKey] = true
		selected = append(selected, action)
	}
	updated, err := s.updateTuneProfileWithSession(profile.ID, inputProfile, "tune_plan_apply", &input.SessionID)
	if err != nil {
		return nil, err
	}
	changed := make([]string, 0, len(changedSet))
	for field := range changedSet {
		changed = append(changed, field)
	}
	sort.Strings(changed)
	return &TunePlanApplyResult{Profile: *updated, AppliedActions: selected, ChangedFields: changed}, nil
}

func (s *Store) GetRetestEvaluation(sessionID int64) (*RetestEvaluation, error) {
	session, err := s.GetTelemetrySession(sessionID)
	if err != nil {
		return nil, err
	}
	currentSummary, err := s.GetSessionIssueSummary(sessionID)
	if err != nil {
		return nil, err
	}
	result := &RetestEvaluation{SessionID: sessionID, Status: "insufficient_data", Summary: "missing_comparison_baseline", Confidence: "low", BaselineReason: currentSummary.BaselineStatus}
	if snapshot := s.recentSnapshotForSession(*session); snapshot != nil {
		result.ChangedFields = append([]string(nil), snapshot.ChangedFields...)
		result.ChangeSourceSessionID = snapshot.SessionID
	}
	if currentSummary.BaselineSession == nil {
		return result, nil
	}
	baseline := currentSummary.BaselineSession
	result.BaselineSession = baseline
	result.Confidence = retestConfidence(currentSummary.BaselineStatus, s.hasComparableBenchmark(sessionID, baseline.ID))
	currentScore := issueGroupScoreSum(currentSummary.Groups)
	baselineEvents, err := s.GetSessionEvents(baseline.ID)
	if err != nil {
		return nil, err
	}
	baselineGroups := BuildSessionIssueGroups(baselineEvents, nil)
	baselineScore := issueGroupScoreSum(baselineGroups)
	result.Metrics = append(result.Metrics, retestMetric("issue_score", currentScore, baselineScore, "lower"))
	result.Metrics = append(result.Metrics, retestMetric("event_count", float64(session.EventCount), float64(baseline.EventCount), "lower"))
	result.Metrics = append(result.Metrics, retestMetric("event_duration_ms", issueGroupDurationSum(currentSummary.Groups), issueGroupDurationSum(baselineGroups), "lower"))
	if session.AvgSpeedKmh != nil && baseline.AvgSpeedKmh != nil {
		result.Metrics = append(result.Metrics, retestMetric("avg_speed_kmh", *session.AvgSpeedKmh, *baseline.AvgSpeedKmh, "higher"))
	}
	if session.MaxSpeedKmh != nil && baseline.MaxSpeedKmh != nil {
		result.Metrics = append(result.Metrics, retestMetric("max_speed_kmh", *session.MaxSpeedKmh, *baseline.MaxSpeedKmh, "higher"))
	}
	if currentBest, okCurrent := s.bestBenchmarkDuration(sessionID); okCurrent {
		if baselineBest, okBaseline := s.bestBenchmarkDuration(baseline.ID); okBaseline {
			result.Metrics = append(result.Metrics, retestMetric("best_run_duration_ms", currentBest, baselineBest, "lower"))
		}
	}
	if currentRisk, okCurrent := s.roadRiskScore(sessionID); okCurrent {
		if baselineRisk, okBaseline := s.roadRiskScore(baseline.ID); okBaseline {
			result.Metrics = append(result.Metrics, retestMetric("risk_score", currentRisk, baselineRisk, "lower"))
		}
	}
	currentTires := s.avgTireTemps(sessionID)
	baselineTires := s.avgTireTemps(baseline.ID)
	if currentTires.ok && baselineTires.ok {
		result.Metrics = append(result.Metrics, retestMetric("front_tire_temp", currentTires.front, baselineTires.front, "target"))
		result.Metrics = append(result.Metrics, retestMetric("rear_tire_temp", currentTires.rear, baselineTires.rear, "target"))
	}
	baselineSummary, err := s.GetSessionIssueSummary(baseline.ID)
	if err == nil {
		result.Metrics = append(result.Metrics, retestMetric("gear_problem_count", float64(gearProblemCount(currentSummary.GearPower)), float64(gearProblemCount(baselineSummary.GearPower)), "lower"))
	}
	result.Status = retestStatusFromMetrics(result.Metrics)
	result.Summary = result.Status
	result.MetricSummary = retestMetricSummary(result.Metrics)
	if profile := s.currentTuneProfileForSession(*session); profile != nil {
		result.RollbackActions = s.rollbackDraftActionsForRetest(*session, profile, result, tunePlanVehicleStatus(*session, *profile))
	}
	return result, nil
}

func (s *Store) currentTuneProfileForSession(session TelemetrySession) *TuneProfile {
	if session.TuneProfileID != nil {
		if profile, err := s.GetTuneProfile(*session.TuneProfileID); err == nil && profile != nil {
			return profile
		}
	}
	return s.issueTuneProfile(session)
}

func (s *Store) applyRetestGuardToDraft(draft *TunePlanDraft, session TelemetrySession, profile *TuneProfile, status string) {
	if draft == nil || profile == nil || status != tunePlanStatusReady {
		return
	}
	evaluation, err := s.GetRetestEvaluation(session.ID)
	if err != nil || evaluation == nil {
		return
	}
	if evaluation.Confidence == "low" || evaluation.Status == "insufficient_data" {
		for i := range draft.Actions {
			draft.Actions[i].TrustLevel = lowerTrustLevel(draft.Actions[i].TrustLevel)
			draft.Actions[i].TrustReasons = appendUniqueString(draft.Actions[i].TrustReasons, "low_retest_confidence")
		}
	}
	rollback := s.rollbackDraftActionsForRetest(session, profile, evaluation, status)
	if evaluation.Status != "worsened" || len(rollback) == 0 {
		return
	}
	for i := range draft.Actions {
		draft.Actions[i].CanApply = false
		draft.Actions[i].BlockedReason = "rollback_first"
		draft.Actions[i].RetestGuard = "rollback_first"
		draft.Actions[i].TrustLevel = "low"
		draft.Actions[i].TrustReasons = appendUniqueString(draft.Actions[i].TrustReasons, "retest_worsened")
	}
	draft.Actions = append(rollback, draft.Actions...)
	draft.Summary = "rollback_before_more_changes"
}

func (s *Store) rollbackDraftActionsForRetest(session TelemetrySession, profile *TuneProfile, evaluation *RetestEvaluation, status string) []TunePlanDraftAction {
	if evaluation == nil || evaluation.Status != "worsened" || profile == nil || status != tunePlanStatusReady {
		return nil
	}
	snapshot := s.recentSnapshotForSession(session)
	if snapshot == nil || snapshot.ChangeReason != "tune_plan_apply" || len(snapshot.ChangedFields) == 0 {
		return nil
	}
	restoreFull := retestWeightedScore(evaluation.Metrics) <= -2.5
	reason := "half_reverse_retest_worsened"
	if restoreFull {
		reason = "rollback_retest_worsened"
	}
	out := make([]TunePlanDraftAction, 0, len(snapshot.ChangedFields))
	for _, field := range snapshot.ChangedFields {
		before, okBefore := tuneProfileNumericField(snapshot.Before, field)
		after, okAfter := tuneProfileNumericField(snapshot.After, field)
		currentPtr := tuneProfileFloatPointer(*profile, field)
		if !okBefore || !okAfter || currentPtr == nil {
			continue
		}
		target := before
		if !restoreFull {
			target = *currentPtr - (after-before)/2
		}
		step := tunePlanFieldStep(field)
		target = roundTunePlanStep(target, step)
		delta := roundTunePlanStep(target-*currentPtr, step)
		if delta == 0 {
			continue
		}
		direction := "decrease"
		if delta > 0 {
			direction = "increase"
		}
		current := *currentPtr
		targetCopy := target
		deltaCopy := delta
		out = append(out, TunePlanDraftAction{
			ID:            fmt.Sprintf("rollback:%d:%s", snapshot.ID, field),
			Family:        "rollback",
			Source:        "retest_guard",
			Confidence:    evaluation.Confidence,
			AdviceLayer:   adviceLayerRollback,
			TrustLevel:    "high",
			TrustReasons:  []string{"retest_worsened", "recent_tune_plan_apply"},
			Rationale:     reason,
			Category:      "rollback",
			Item:          field,
			FieldKey:      field,
			Direction:     direction,
			Reason:        reason,
			CurrentValue:  &current,
			TargetValue:   &targetCopy,
			Delta:         &deltaCopy,
			Unit:          tunePlanFieldUnit(field),
			Step:          step,
			CanApply:      true,
			RetestGuard:   "rollback_first",
			BlockedReason: "",
		})
	}
	return out
}

func draftActionsForWholeCarAdjustment(action WholeCarAdjustment, profile *TuneProfile, status string) []TunePlanDraftAction {
	fields := tunePlanFields(action, profile)
	out := make([]TunePlanDraftAction, 0, len(fields))
	for index, field := range fields {
		draft := TunePlanDraftAction{
			ID:          fmt.Sprintf("%d:%s:%s:%s:%s", action.Priority, action.Category, action.Item, action.Direction, field.key),
			Family:      action.Family,
			Source:      action.Source,
			Confidence:  action.Confidence,
			AdviceLayer: adviceLayerForWholeCarAction(action),
			TrustLevel:  actionTrustLevel(action, status, field.value != nil),
			Rationale:   adviceRationale(action.Source, action.Reason),
			Category:    action.Category,
			Item:        action.Item,
			FieldKey:    field.key,
			Direction:   action.Direction,
			Reason:      action.Reason,
			Unit:        field.unit,
			Step:        field.step,
			CanApply:    false,
		}
		draft.TrustReasons = actionTrustReasons(action, status, field.value != nil)
		draft.MissingInputs = actionMissingInputs(action, field.value != nil)
		if len(fields) > 1 {
			draft.ID = fmt.Sprintf("%s:%d", draft.ID, index)
		}
		switch status {
		case tunePlanStatusReady:
		default:
			draft.BlockedReason = status
			out = append(out, draft)
			continue
		}
		if strings.EqualFold(strings.TrimSpace(action.Direction), "check") {
			draft.BlockedReason = "manual_review_required"
			out = append(out, draft)
			continue
		}
		if field.value == nil {
			draft.BlockedReason = "field_locked_or_blank"
			out = append(out, draft)
			continue
		}
		delta, ok := tunePlanActionDelta(action.Direction, action.Amount, *field.value, field.step, field.unit)
		if !ok || delta == 0 {
			draft.BlockedReason = "no_numeric_adjustment"
			out = append(out, draft)
			continue
		}
		target := roundTunePlanStep(*field.value+delta, field.step)
		delta = roundTunePlanStep(target-*field.value, field.step)
		if delta == 0 {
			draft.BlockedReason = "no_change"
			out = append(out, draft)
			continue
		}
		current := *field.value
		draft.CurrentValue = &current
		draft.TargetValue = &target
		draft.Delta = &delta
		draft.CanApply = true
		out = append(out, draft)
	}
	return out
}

type tunePlanField struct {
	key   string
	unit  string
	step  float64
	value *float64
}

func tunePlanFields(action WholeCarAdjustment, profile *TuneProfile) []tunePlanField {
	if profile == nil {
		return nil
	}
	pick := func(keys ...string) []tunePlanField {
		out := make([]tunePlanField, 0, len(keys))
		for _, key := range keys {
			out = append(out, tunePlanField{key: key, unit: tunePlanFieldUnit(key), step: tunePlanFieldStep(key), value: tuneProfileFloatPointer(*profile, key)})
		}
		return out
	}
	driven := func(front string, rear string) []tunePlanField {
		switch strings.ToUpper(strings.TrimSpace(profile.Drivetrain)) {
		case "FWD":
			return pick(front)
		case "RWD":
			return pick(rear)
		default:
			return pick(front, rear)
		}
	}
	switch action.Item {
	case "front_tire_pressure":
		return pick("frontTirePressure")
	case "rear_tire_pressure":
		return pick("rearTirePressure")
	case "gear_1", "gear_2", "gear_3", "gear_4", "gear_5", "gear_6", "gear_7", "gear_8", "gear_9", "gear_10":
		return pick("gear" + strings.TrimPrefix(action.Item, "gear_"))
	case "current_gear":
		gear := int(action.Evidence["gear"] + 0.5)
		if gear < 1 || gear > 10 {
			return nil
		}
		return pick(fmt.Sprintf("gear%d", gear))
	case "final_drive":
		return pick("finalDrive")
	case "brake_balance":
		return pick("brakeBalance")
	case "brake_pressure":
		return pick("brakePressure")
	case "front_diff_accel":
		return pick("frontDiffAccel")
	case "front_diff_decel":
		return pick("frontDiffDecel")
	case "rear_diff_accel":
		return pick("rearDiffAccel")
	case "rear_diff_decel":
		return pick("rearDiffDecel")
	case "drive_diff_accel":
		return driven("frontDiffAccel", "rearDiffAccel")
	case "drive_diff_decel":
		return driven("frontDiffDecel", "rearDiffDecel")
	case "drive_tire_pressure":
		return driven("frontTirePressure", "rearTirePressure")
	case "tire_pressure":
		return pick("frontTirePressure", "rearTirePressure")
	case "front_arb":
		return pick("frontArb")
	case "rear_arb":
		return pick("rearArb")
	case "front_rebound":
		return pick("frontRebound")
	case "rear_rebound":
		return pick("rearRebound")
	case "front_camber":
		return pick("frontCamber")
	case "front_and_rear_aero":
		return pick("frontAero", "rearAero")
	case "ride_height":
		return pick("frontRideHeight", "rearRideHeight")
	case "spring_rate":
		return pick("frontSpring", "rearSpring")
	case "bump":
		return pick("frontBump", "rearBump")
	default:
		return nil
	}
}

func tunePlanFieldUnit(key string) string {
	switch key {
	case "frontTirePressure", "rearTirePressure":
		return "BAR"
	case "frontCamber", "rearCamber", "frontToe", "rearToe", "caster":
		return "deg"
	case "frontSpring", "rearSpring":
		return "kgf/mm"
	case "frontRideHeight", "rearRideHeight":
		return "cm"
	case "frontAero", "rearAero":
		return "kgf"
	case "brakeBalance", "brakePressure", "frontDiffAccel", "frontDiffDecel", "rearDiffAccel", "rearDiffDecel", "centerDiffBalance":
		return "%"
	default:
		return ""
	}
}

func tunePlanFieldStep(key string) float64 {
	switch key {
	case "finalDrive", "gear1", "gear2", "gear3", "gear4", "gear5", "gear6", "gear7", "gear8", "gear9", "gear10":
		return 0.01
	case "frontTirePressure", "rearTirePressure":
		return 0.01
	case "frontArb", "rearArb":
		return 0.1
	case "frontAero", "rearAero", "brakeBalance", "brakePressure", "frontDiffAccel", "frontDiffDecel", "rearDiffAccel", "rearDiffDecel", "centerDiffBalance":
		return 1
	default:
		return 0.1
	}
}

func actionTrustLevel(action WholeCarAdjustment, status string, hasCurrentValue bool) string {
	if status != tunePlanStatusReady || !hasCurrentValue {
		return "blocked"
	}
	switch strings.ToLower(strings.TrimSpace(action.Confidence)) {
	case "high":
		return "high"
	case "low", "needs_profile":
		return "low"
	default:
		return "medium"
	}
}

func actionTrustReasons(action WholeCarAdjustment, status string, hasCurrentValue bool) []string {
	reasons := []string{}
	if status != tunePlanStatusReady {
		reasons = append(reasons, status)
	}
	if !hasCurrentValue {
		reasons = append(reasons, "missing_current_value")
	}
	if strings.TrimSpace(action.Confidence) != "" {
		reasons = append(reasons, "model_confidence_"+strings.ToLower(strings.TrimSpace(action.Confidence)))
	}
	if strings.TrimSpace(action.Source) != "" {
		reasons = append(reasons, "source_"+strings.ToLower(strings.TrimSpace(action.Source)))
	}
	return reasons
}

func actionMissingInputs(action WholeCarAdjustment, hasCurrentValue bool) []string {
	missing := []string{}
	if !hasCurrentValue {
		missing = append(missing, "current_tune_value")
	}
	if action.Source == "gear_power_diagnostic" {
		if action.Evidence["power_band_start_rpm"] <= 0 || action.Evidence["power_band_end_rpm"] <= 0 {
			missing = append(missing, "profile_power_band")
		}
		if action.Evidence["power_band_high_load_samples"] < 4 {
			missing = append(missing, "high_load_samples")
		}
	}
	return missing
}

func lowerTrustLevel(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "high":
		return "medium"
	case "blocked":
		return "blocked"
	default:
		return "low"
	}
}

func appendUniqueString(values []string, value string) []string {
	if strings.TrimSpace(value) == "" {
		return values
	}
	for _, existing := range values {
		if existing == value {
			return values
		}
	}
	return append(values, value)
}

func tunePlanActionDelta(direction string, amount string, current float64, step float64, unit string) (float64, bool) {
	normalizedAmount := strings.ToLower(strings.TrimSpace(amount))
	normalizedDirection := strings.ToLower(strings.TrimSpace(direction))
	magnitude := 0.0
	switch normalizedAmount {
	case "one small step", "avoid bottoming", "0.5 psi":
		magnitude = step
	case "slightly more negative":
		return -step, true
	default:
		if strings.Contains(normalizedAmount, "%") {
			base := firstTunePlanNumber(normalizedAmount)
			if base <= 0 {
				return 0, false
			}
			if unit == "%" {
				magnitude = base
			} else {
				magnitude = math.Max(step, math.Abs(current)*base/100)
			}
		} else {
			magnitude = firstTunePlanNumber(normalizedAmount)
		}
	}
	if magnitude <= 0 {
		return 0, false
	}
	magnitude = math.Max(step, roundTunePlanStep(magnitude, step))
	switch normalizedDirection {
	case "decrease":
		return -magnitude, true
	case "increase":
		return magnitude, true
	case "check":
		if strings.Contains(normalizedAmount, "bottom") || strings.Contains(normalizedAmount, "small step") {
			return magnitude, true
		}
		return 0, false
	default:
		return 0, false
	}
}

func firstTunePlanNumber(value string) float64 {
	for _, span := range numberSpans(value) {
		return span.value
	}
	return 0
}

func roundTunePlanStep(value float64, step float64) float64 {
	if step <= 0 {
		return value
	}
	decimals := 2
	if step >= 1 {
		decimals = 0
	} else if step >= 0.1 {
		decimals = 1
	}
	factor := math.Pow10(decimals)
	return math.Round(math.Round(value/step)*step*factor) / factor
}

func tunePlanVehicleStatus(session TelemetrySession, profile TuneProfile) string {
	if session.TuneProfileID == nil {
		return tunePlanStatusUnbound
	}
	if session.CarOrdinal == nil || profile.CarOrdinal == nil || strings.TrimSpace(session.CarClass) == "" || strings.TrimSpace(profile.CarClass) == "" {
		return tunePlanStatusCannotVerify
	}
	if *session.CarOrdinal != *profile.CarOrdinal || !strings.EqualFold(strings.TrimSpace(session.CarClass), strings.TrimSpace(profile.CarClass)) {
		return tunePlanStatusMismatch
	}
	return tunePlanStatusReady
}

func setTuneProfileInputFloat(input *TuneProfileInput, key string, value float64) bool {
	value = roundTunePlanStep(value, tunePlanFieldStep(key))
	switch key {
	case "frontTirePressure":
		input.FrontTirePressure = &value
	case "rearTirePressure":
		input.RearTirePressure = &value
	case "finalDrive":
		input.FinalDrive = &value
	case "gear1":
		input.Gear1 = &value
	case "gear2":
		input.Gear2 = &value
	case "gear3":
		input.Gear3 = &value
	case "gear4":
		input.Gear4 = &value
	case "gear5":
		input.Gear5 = &value
	case "gear6":
		input.Gear6 = &value
	case "gear7":
		input.Gear7 = &value
	case "gear8":
		input.Gear8 = &value
	case "gear9":
		input.Gear9 = &value
	case "gear10":
		input.Gear10 = &value
	case "frontCamber":
		input.FrontCamber = &value
	case "rearCamber":
		input.RearCamber = &value
	case "frontToe":
		input.FrontToe = &value
	case "rearToe":
		input.RearToe = &value
	case "caster":
		input.Caster = &value
	case "frontArb":
		input.FrontARB = &value
	case "rearArb":
		input.RearARB = &value
	case "frontSpring":
		input.FrontSpring = &value
	case "rearSpring":
		input.RearSpring = &value
	case "frontRideHeight":
		input.FrontRideHeight = &value
	case "rearRideHeight":
		input.RearRideHeight = &value
	case "frontRebound":
		input.FrontRebound = &value
	case "rearRebound":
		input.RearRebound = &value
	case "frontBump":
		input.FrontBump = &value
	case "rearBump":
		input.RearBump = &value
	case "frontAero":
		input.FrontAero = &value
	case "rearAero":
		input.RearAero = &value
	case "aeroBalance":
		input.AeroBalance = &value
	case "brakeBalance":
		input.BrakeBalance = &value
	case "brakePressure":
		input.BrakePressure = &value
	case "frontDiffAccel":
		input.FrontDiffAccel = &value
	case "frontDiffDecel":
		input.FrontDiffDecel = &value
	case "rearDiffAccel":
		input.RearDiffAccel = &value
	case "rearDiffDecel":
		input.RearDiffDecel = &value
	case "centerDiffBalance":
		input.CenterDiffBalance = &value
	default:
		return false
	}
	return true
}

func issueGroupScoreSum(groups []SessionIssueGroup) float64 {
	total := 0.0
	for _, group := range groups {
		total += issueGroupScore(group)
	}
	return total
}

func issueGroupDurationSum(groups []SessionIssueGroup) float64 {
	total := int64(0)
	for _, group := range groups {
		total += group.TotalDurationMS
	}
	return float64(total)
}

func (s *Store) bestBenchmarkDuration(sessionID int64) (float64, bool) {
	runs, err := s.listBenchmarkRunsForSession(sessionID)
	if err != nil {
		return 0, false
	}
	for _, run := range runs {
		if run.Valid && run.DurationMS > 0 {
			return float64(run.DurationMS), true
		}
	}
	return 0, false
}

func (s *Store) roadRiskScore(sessionID int64) (float64, bool) {
	evaluation, err := s.EvaluateRoadSession(sessionID)
	if err != nil || evaluation == nil {
		return 0, false
	}
	return evaluation.RiskScore, true
}

type tireTempAverage struct {
	front float64
	rear  float64
	ok    bool
}

func (s *Store) avgTireTemps(sessionID int64) tireTempAverage {
	samples, err := s.GetSessionTelemetrySamples(sessionID, 10000)
	if err != nil || len(samples) == 0 {
		return tireTempAverage{}
	}
	front := 0.0
	rear := 0.0
	count := 0.0
	for _, sample := range samples {
		if sample.TireTempFrontAvg <= 0 && sample.TireTempRearAvg <= 0 {
			continue
		}
		front += sample.TireTempFrontAvg
		rear += sample.TireTempRearAvg
		count++
	}
	if count == 0 {
		return tireTempAverage{}
	}
	return tireTempAverage{front: front / count, rear: rear / count, ok: true}
}

func gearProblemCount(diag GearPowerDiagnostic) int {
	count := 0
	for _, gear := range diag.Gears {
		if gear.Finding != "" && gear.Finding != "ok" {
			count++
		}
	}
	if diag.TopSpeedFinding != "" && diag.TopSpeedFinding != "top_speed_ok" {
		count++
	}
	if diag.LaunchFinding != "" {
		count++
	}
	return count
}

func retestMetric(key string, current float64, baseline float64, direction string) RetestMetric {
	delta := current - baseline
	return RetestMetric{Key: key, Current: current, Baseline: baseline, Delta: delta, Direction: direction, Status: retestMetricStatus(key, current, baseline, direction)}
}

func retestMetricStatus(key string, current float64, baseline float64, direction string) string {
	if baseline == 0 {
		if current == 0 {
			return "unchanged"
		}
		if direction == "lower" {
			return "worsened"
		}
		return "improved"
	}
	ratio := current / baseline
	threshold := retestMetricThreshold(key)
	switch direction {
	case "lower":
		if ratio < 1-threshold {
			return "improved"
		}
		if ratio > 1+threshold {
			return "worsened"
		}
	case "higher":
		if ratio > 1+threshold {
			return "improved"
		}
		if ratio < 1-threshold {
			return "worsened"
		}
	case "target":
		if math.Abs(current-baseline) <= 3 {
			return "unchanged"
		}
		if current < baseline {
			return "improved"
		}
		return "worsened"
	}
	return "unchanged"
}

func retestMetricThreshold(key string) float64 {
	switch key {
	case "avg_speed_kmh", "max_speed_kmh":
		return 0.02
	case "best_run_duration_ms":
		return 0.015
	case "risk_score":
		return 0.05
	default:
		return 0.1
	}
}

func retestStatusFromMetrics(metrics []RetestMetric) string {
	score := retestWeightedScore(metrics)
	if score >= 1.5 {
		return "improved"
	}
	if score <= -1.5 {
		return "worsened"
	}
	return "unchanged"
}

func retestWeightedScore(metrics []RetestMetric) float64 {
	score := 0.0
	for _, metric := range metrics {
		weight := retestMetricWeight(metric.Key)
		switch metric.Status {
		case "improved":
			score += weight
		case "worsened":
			score -= weight
		}
	}
	return score
}

func retestMetricSummary(metrics []RetestMetric) []string {
	out := make([]string, 0, 3)
	for _, status := range []string{"worsened", "improved"} {
		for _, metric := range metrics {
			if metric.Status == status {
				out = append(out, metric.Key+":"+status)
				if len(out) >= 3 {
					return out
				}
			}
		}
	}
	if len(out) == 0 {
		out = append(out, "metrics:unchanged")
	}
	return out
}

func retestConfidence(baselineReason string, hasBenchmark bool) string {
	switch baselineReason {
	case "matched_profile_track_driver":
		if hasBenchmark {
			return "high"
		}
		return "medium"
	case "matched_vehicle_class_usecase_driver":
		return "low"
	default:
		return "low"
	}
}

func (s *Store) hasComparableBenchmark(currentSessionID int64, baselineSessionID int64) bool {
	_, currentOK := s.bestBenchmarkDuration(currentSessionID)
	_, baselineOK := s.bestBenchmarkDuration(baselineSessionID)
	return currentOK && baselineOK
}

func retestMetricWeight(key string) float64 {
	switch key {
	case "best_run_duration_ms":
		return 3
	case "avg_speed_kmh":
		return 2.5
	case "risk_score":
		return 2
	case "issue_score":
		return 1.5
	case "event_count", "event_duration_ms", "gear_problem_count":
		return 1
	case "max_speed_kmh":
		return 0.75
	case "front_tire_temp", "rear_tire_temp":
		return 0.5
	default:
		return 1
	}
}

func (s *Store) updateTuneProfileWithSession(id int64, input TuneProfileInput, reason string, sessionID *int64) (*TuneProfile, error) {
	if strings.TrimSpace(input.CarName) == "" {
		return nil, errors.New("car name is required")
	}
	input = normalizeTuneProfilePower(input)
	before, err := s.GetTuneProfile(id)
	if err != nil {
		return nil, err
	}
	now := nowText()
	args := append(profileUpdateArgs(input, now), id)
	result, err := s.db.Exec(`UPDATE tune_profile SET
		car_name = ?, car_ordinal = ?, car_category = ?, car_class = ?, pi = ?, drivetrain = ?, num_cylinders = ?, use_case = ?, version_name = ?,
		power_kw = ?, torque_nm = ?, weight_kg = ?, front_weight_pct = ?, power_to_weight_kw_per_kg = ?, peak_torque_rpm = ?, peak_power_rpm = ?, redline_rpm = ?, updated_at = ?,
		front_tire_pressure = ?, rear_tire_pressure = ?, final_drive = ?, gear_1 = ?, gear_2 = ?, gear_3 = ?, gear_4 = ?, gear_5 = ?, gear_6 = ?, gear_7 = ?, gear_8 = ?, gear_9 = ?, gear_10 = ?,
		front_camber = ?, rear_camber = ?, front_toe = ?, rear_toe = ?, caster = ?, front_arb = ?, rear_arb = ?,
		front_spring = ?, rear_spring = ?, front_ride_height = ?, rear_ride_height = ?,
		front_rebound = ?, rear_rebound = ?, front_bump = ?, rear_bump = ?,
		front_aero = ?, rear_aero = ?, aero_balance = ?, brake_balance = ?, brake_pressure = ?,
		front_diff_accel = ?, front_diff_decel = ?, rear_diff_accel = ?, rear_diff_decel = ?, center_diff_balance = ?, notes = ?
		WHERE id = ?`, args...)
	if err != nil {
		return nil, err
	}
	if changed, _ := result.RowsAffected(); changed == 0 {
		return nil, sql.ErrNoRows
	}
	if err := s.upsertCarIdentity(input); err != nil {
		return nil, err
	}
	after, err := s.GetTuneProfile(id)
	if err != nil {
		return nil, err
	}
	changedFields, err := changedTuneProfileFields(before, after)
	if err != nil {
		return nil, err
	}
	if len(changedFields) > 0 {
		if err := s.insertTuneProfileSnapshot(before, after, reason, sessionID, changedFields); err != nil {
			return nil, err
		}
		if err := s.pruneTuneProfileSnapshots(id, 5); err != nil {
			return nil, err
		}
	}
	return after, nil
}

func parseTunePlanActionID(id string) (int, string) {
	parts := strings.Split(id, ":")
	if len(parts) < 5 {
		return 0, ""
	}
	priority, _ := strconv.Atoi(parts[0])
	return priority, parts[4]
}
