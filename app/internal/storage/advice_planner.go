package storage

import (
	"sort"
	"strings"

	"fh6worker/internal/telemetry"
)

const (
	adviceLayerRollback    = "rollback"
	adviceLayerPrimary     = "primary"
	adviceLayerPowertrain  = "powertrain"
	adviceLayerSupport     = "support"
	adviceLayerAlternative = "alternative"

	quickDirectionOnlyAmount            = "direction_only"
	quickSuggestionMaxPerFamily         = 2
	quickSuggestionNextStepBindProfile  = "bind_tune_profile_for_values"
	quickSuggestionNextStepFillFields   = "fill_or_unlock_tune_fields"
	quickSuggestionNextStepCollectPower = "collect_power_samples"
	quickSuggestionNextStepExpertMode   = "use_expert_for_concrete_values"
)

func BuildQuickAdviceSuggestions(groups []SessionIssueGroup, gear GearPowerDiagnostic, profile *TuneProfile) []QuickSuggestion {
	plan := BuildWholeCarTuningPlan(groups, gear, profile)
	actions := finalizeWholeCarAdviceActions(plan.Actions, 0)
	out := make([]QuickSuggestion, 0, len(actions))
	perFamily := map[string]int{}
	for _, action := range actions {
		if perFamily[action.Family] >= quickSuggestionMaxPerFamily {
			continue
		}
		fields := quickFieldKeysForAction(action, profile)
		missing := quickAdviceMissingInputs(action, fields, profile)
		out = append(out, QuickSuggestion{
			Family:        action.Family,
			Source:        action.Source,
			Confidence:    action.Confidence,
			TrustLevel:    quickAdviceTrustLevel(action, missing),
			AdviceLayer:   adviceLayerForWholeCarAction(action),
			Category:      action.Category,
			Item:          action.Item,
			Direction:     action.Direction,
			Amount:        quickDirectionOnlyAmount,
			Reason:        action.Reason,
			Rationale:     adviceRationale(action.Source, action.Reason),
			NextStep:      quickAdviceNextStep(action, missing),
			FieldKeys:     fields,
			MissingInputs: missing,
			CanApply:      false,
			BlockedReason: "quick_mode_directional_only",
		})
		perFamily[action.Family]++
		if len(out) >= quickSuggestionLimit {
			break
		}
	}
	return out
}

func QuickMissingFieldsFromSuggestions(suggestions []QuickSuggestion) []string {
	missing := map[string]bool{}
	for _, suggestion := range suggestions {
		for _, field := range suggestion.FieldKeys {
			if stringSliceContains(suggestion.MissingInputs, "tune_profile") || stringSliceContains(suggestion.MissingInputs, "current_tune_value") {
				missing[field] = true
			}
		}
	}
	out := make([]string, 0, len(missing))
	for field := range missing {
		out = append(out, field)
	}
	sort.Strings(out)
	return out
}

func finalizeWholeCarAdviceActions(actions []WholeCarAdjustment, limit int) []WholeCarAdjustment {
	if len(actions) == 0 {
		return nil
	}
	selected := map[string]WholeCarAdjustment{}
	order := []string{}
	for _, action := range actions {
		action.Source = strings.TrimSpace(action.Source)
		if action.Source == "" {
			action.Source = "issue_group"
		}
		key := actionConflictKey(telemetry.SuggestedAction{Category: action.Category, Item: action.Item, Direction: action.Direction, Amount: action.Amount}, action.Evidence)
		if key == "" {
			key = action.Category + "/" + action.Item
		}
		existing, ok := selected[key]
		if !ok {
			selected[key] = action
			order = append(order, key)
			continue
		}
		if betterWholeCarAdvice(action, existing) {
			selected[key] = action
		}
	}
	out := make([]WholeCarAdjustment, 0, len(order))
	for _, key := range order {
		out = append(out, selected[key])
	}
	sort.SliceStable(out, func(i, j int) bool {
		if adviceLayerRank(adviceLayerForWholeCarAction(out[i])) != adviceLayerRank(adviceLayerForWholeCarAction(out[j])) {
			return adviceLayerRank(adviceLayerForWholeCarAction(out[i])) < adviceLayerRank(adviceLayerForWholeCarAction(out[j]))
		}
		if adviceConfidenceRank(out[i].Confidence) != adviceConfidenceRank(out[j].Confidence) {
			return adviceConfidenceRank(out[i].Confidence) > adviceConfidenceRank(out[j].Confidence)
		}
		return out[i].Priority < out[j].Priority
	})
	if limit > 0 && len(out) > limit {
		out = out[:limit]
	}
	for i := range out {
		out[i].Priority = i
	}
	return out
}

func finalizeRoadDecisionActions(actions []RoadTuningDecisionAction, limit int) []RoadTuningDecisionAction {
	if len(actions) == 0 {
		return nil
	}
	for i := range actions {
		completeRoadDecisionAction(&actions[i])
	}
	actions = resolveRoadActionConflicts(actions)
	sort.SliceStable(actions, func(i, j int) bool {
		if adviceLayerRank(actions[i].AdviceLayer) != adviceLayerRank(actions[j].AdviceLayer) {
			return adviceLayerRank(actions[i].AdviceLayer) < adviceLayerRank(actions[j].AdviceLayer)
		}
		if actions[i].CanAutoApply != actions[j].CanAutoApply {
			return actions[i].CanAutoApply
		}
		return adviceConfidenceRank(actions[i].Confidence) > adviceConfidenceRank(actions[j].Confidence)
	})
	if limit > 0 && len(actions) > limit {
		actions = actions[:limit]
	}
	for i := range actions {
		switch i {
		case 0:
			actions[i].Role = "primary"
		case 1:
			actions[i].Role = "support"
		default:
			actions[i].Role = "alternative"
		}
	}
	return actions
}

func finalizeTunePlanDraft(draft *TunePlanDraft) {
	if draft == nil || len(draft.Actions) == 0 {
		return
	}
	for i := range draft.Actions {
		completeTunePlanDraftAction(&draft.Actions[i])
	}
	sort.SliceStable(draft.Actions, func(i, j int) bool {
		return betterDraftAction(draft.Actions[i], draft.Actions[j])
	})
	selected := make([]TunePlanDraftAction, 0, len(draft.Actions))
	selectedByKey := map[string]int{}
	for _, action := range draft.Actions {
		key := draftActionConflictKey(action)
		if key == "" {
			selected = append(selected, action)
			continue
		}
		index, ok := selectedByKey[key]
		if !ok {
			selectedByKey[key] = len(selected)
			selected = append(selected, action)
			continue
		}
		existing := selected[index]
		reason := "duplicate_action_removed"
		if existing.Direction != action.Direction {
			reason = "same_field_direction_conflict"
		}
		if betterDraftAction(action, existing) {
			action.ConflictReason = reason
			selected[index] = action
			draft.Conflicts = append(draft.Conflicts, TuningConflict{
				Key:         key,
				KeptItem:    action.Item + "/" + action.Direction,
				DroppedItem: existing.Item + "/" + existing.Direction,
				Reason:      reason,
			})
			continue
		}
		selected[index].ConflictReason = reason
		draft.Conflicts = append(draft.Conflicts, TuningConflict{
			Key:         key,
			KeptItem:    existing.Item + "/" + existing.Direction,
			DroppedItem: action.Item + "/" + action.Direction,
			Reason:      reason,
		})
	}
	draft.Actions = limitDraftAdviceActions(selected, 3)
}

func limitDraftAdviceActions(actions []TunePlanDraftAction, limit int) []TunePlanDraftAction {
	if limit <= 0 || len(actions) <= limit {
		return actions
	}
	return actions[:limit]
}

func completeRoadDecisionAction(action *RoadTuningDecisionAction) {
	if action == nil {
		return
	}
	if action.AdviceLayer == "" {
		action.AdviceLayer = adviceLayerForSource(action.Source, action.Family)
	}
	if action.TrustLevel == "" {
		action.TrustLevel = trustLevelFromConfidence(action.Confidence)
	}
	if action.Rationale == "" {
		action.Rationale = adviceRationale(action.Source, action.Reason)
	}
}

func completeTunePlanDraftAction(action *TunePlanDraftAction) {
	if action == nil {
		return
	}
	if action.AdviceLayer == "" {
		action.AdviceLayer = adviceLayerForSource(action.Source, action.Family)
	}
	if action.Rationale == "" {
		action.Rationale = adviceRationale(action.Source, action.Reason)
	}
	if action.TrustLevel == "" {
		action.TrustLevel = trustLevelFromConfidence(action.Confidence)
	}
}

func adviceLayerForWholeCarAction(action WholeCarAdjustment) string {
	return adviceLayerForSource(action.Source, action.Family)
}

func adviceLayerForSource(source string, family string) string {
	source = strings.ToLower(strings.TrimSpace(source))
	family = strings.ToLower(strings.TrimSpace(family))
	switch {
	case source == "retest_guard" || source == "rollback" || family == "rollback":
		return adviceLayerRollback
	case source == "gear_power_diagnostic" || family == "gearing_acceleration":
		return adviceLayerPowertrain
	case source == "issue_group" || source == "road_tuning_model" || source == "workbook":
		return adviceLayerPrimary
	default:
		return adviceLayerSupport
	}
}

func adviceRationale(source string, reason string) string {
	reason = strings.TrimSpace(reason)
	if reason != "" {
		return reason
	}
	switch strings.ToLower(strings.TrimSpace(source)) {
	case "gear_power_diagnostic":
		return "source_gear_power_diagnostic"
	case "retest_guard":
		return "rollback_retest_worsened"
	case "road_tuning_model":
		return "source_road_tuning_model"
	default:
		return "source_issue_group"
	}
}

func quickAdviceMissingInputs(action WholeCarAdjustment, fields []string, profile *TuneProfile) []string {
	missing := []string{}
	if profile == nil {
		missing = append(missing, "tune_profile")
	} else {
		for _, field := range fields {
			if tuneProfileFloatPointer(*profile, field) == nil {
				missing = appendUniqueString(missing, "current_tune_value")
				break
			}
		}
	}
	if action.Source == "gear_power_diagnostic" {
		if action.Evidence["power_band_start_rpm"] <= 0 || action.Evidence["power_band_end_rpm"] <= 0 {
			missing = appendUniqueString(missing, "profile_power_band")
		}
		if action.Evidence["power_band_high_load_samples"] < 4 {
			missing = appendUniqueString(missing, "high_load_samples")
		}
	}
	return missing
}

func quickAdviceTrustLevel(action WholeCarAdjustment, missing []string) string {
	if len(missing) > 0 {
		return "low"
	}
	return trustLevelFromConfidence(action.Confidence)
}

func quickAdviceNextStep(action WholeCarAdjustment, missing []string) string {
	if stringSliceContains(missing, "tune_profile") {
		return quickSuggestionNextStepBindProfile
	}
	if stringSliceContains(missing, "current_tune_value") {
		return quickSuggestionNextStepFillFields
	}
	if action.Source == "gear_power_diagnostic" && stringSliceContains(missing, "high_load_samples") {
		return quickSuggestionNextStepCollectPower
	}
	return quickSuggestionNextStepExpertMode
}

func betterWholeCarAdvice(left WholeCarAdjustment, right WholeCarAdjustment) bool {
	if adviceLayerRank(adviceLayerForWholeCarAction(left)) != adviceLayerRank(adviceLayerForWholeCarAction(right)) {
		return adviceLayerRank(adviceLayerForWholeCarAction(left)) < adviceLayerRank(adviceLayerForWholeCarAction(right))
	}
	if adviceConfidenceRank(left.Confidence) != adviceConfidenceRank(right.Confidence) {
		return adviceConfidenceRank(left.Confidence) > adviceConfidenceRank(right.Confidence)
	}
	return left.Priority < right.Priority
}

func betterDraftAction(left TunePlanDraftAction, right TunePlanDraftAction) bool {
	if adviceLayerRank(left.AdviceLayer) != adviceLayerRank(right.AdviceLayer) {
		return adviceLayerRank(left.AdviceLayer) < adviceLayerRank(right.AdviceLayer)
	}
	if adviceTrustRank(left.TrustLevel) != adviceTrustRank(right.TrustLevel) {
		return adviceTrustRank(left.TrustLevel) > adviceTrustRank(right.TrustLevel)
	}
	return left.ID < right.ID
}

func draftActionConflictKey(action TunePlanDraftAction) string {
	if strings.TrimSpace(action.FieldKey) != "" {
		return strings.TrimSpace(action.FieldKey)
	}
	if strings.TrimSpace(action.Category) != "" || strings.TrimSpace(action.Item) != "" {
		return strings.TrimSpace(action.Category) + "/" + strings.TrimSpace(action.Item)
	}
	return ""
}

func adviceLayerRank(layer string) int {
	switch strings.ToLower(strings.TrimSpace(layer)) {
	case adviceLayerRollback:
		return 0
	case adviceLayerPrimary:
		return 1
	case adviceLayerPowertrain:
		return 2
	case adviceLayerSupport:
		return 3
	case adviceLayerAlternative:
		return 4
	default:
		return 5
	}
}

func adviceConfidenceRank(confidence string) int {
	switch strings.ToLower(strings.TrimSpace(confidence)) {
	case "high":
		return 3
	case "medium":
		return 2
	case "low", "needs_profile":
		return 1
	default:
		return 0
	}
}

func adviceTrustRank(trust string) int {
	switch strings.ToLower(strings.TrimSpace(trust)) {
	case "high":
		return 4
	case "medium":
		return 3
	case "low":
		return 2
	case "blocked":
		return 1
	default:
		return 0
	}
}

func trustLevelFromConfidence(confidence string) string {
	switch strings.ToLower(strings.TrimSpace(confidence)) {
	case "high":
		return "high"
	case "low", "needs_profile":
		return "low"
	default:
		return "medium"
	}
}

func stringSliceContains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}
