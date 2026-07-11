package advisor

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"

	"fh6worker/internal/storage"
	"fh6worker/internal/telemetry"
)

func GenerateTuningReport(session storage.TelemetrySession, profile *storage.TuneProfile, events []telemetry.DetectedEvent, language string, evaluations ...*storage.RoadSessionEvaluation) string {
	var evaluation *storage.RoadSessionEvaluation
	if len(evaluations) > 0 {
		evaluation = evaluations[0]
	}
	if language == "zh" {
		return generateZH(session, profile, events, evaluation)
	}
	return generateEN(session, profile, events, evaluation)
}

func GenerateTuningReportWithIssueSummary(session storage.TelemetrySession, profile *storage.TuneProfile, events []telemetry.DetectedEvent, language string, evaluation *storage.RoadSessionEvaluation, summary *storage.SessionIssueSummary) string {
	return GenerateTuningReportWithRoadDecision(session, profile, events, language, evaluation, summary, nil)
}

func GenerateTuningReportWithRoadDecision(session storage.TelemetrySession, profile *storage.TuneProfile, events []telemetry.DetectedEvent, language string, evaluation *storage.RoadSessionEvaluation, summary *storage.SessionIssueSummary, decision *storage.RoadTuningDecision) string {
	if summary == nil || len(summary.Groups) == 0 {
		return GenerateTuningReport(session, profile, events, language, evaluation)
	}
	if language == "zh" {
		return generateIssueSummaryWithDecisionZH(session, profile, evaluation, summary, decision)
	}
	return generateIssueSummaryWithDecisionEN(session, profile, evaluation, summary, decision)
}

func generateIssueSummaryWithDecisionZH(session storage.TelemetrySession, profile *storage.TuneProfile, evaluation *storage.RoadSessionEvaluation, summary *storage.SessionIssueSummary, decision *storage.RoadTuningDecision) string {
	var b strings.Builder
	b.WriteString("# FH6 本地调校建议报告\n\n")
	writeSessionZH(&b, session, profile)
	writeRoadEvaluationZH(&b, evaluation)
	writeRoadTuningDecisionZH(&b, decision, profile)
	b.WriteString("## 一、合并问题项\n\n")
	for _, group := range summary.Groups {
		fmt.Fprintf(&b, "- %s：%s，出现 %d 次，累计 %s，对比上轮：%s。\n", issueFamilyLabelZH(group.Family), severityZH(group.Severity), group.EventCount, duration(group.TotalDurationMS), issueComparisonZH(group.Comparison))
	}
	writeIssueSummaryEvidenceZH(&b, summary.Groups)
	b.WriteString("\n## 下一轮复测方法\n\n")
	b.WriteString("使用同一车辆、同一赛道、同一驾驶方式复测。每轮只应用已确认的 1-3 个动作，观察平均速度、风险分、问题组次数/持续时间是否改善。\n")
	b.WriteString("\n## 风险与回退\n\n")
	if decision != nil && decision.RollbackRecommended {
		b.WriteString("本轮模型建议优先回退或反向微调最近相关修改，再重新测试。\n")
	} else if len(summary.RecentChangeFields) > 0 {
		fmt.Fprintf(&b, "本轮检测到最近修改项：%s。若平均速度、风险或问题持续时间恶化，优先回退相关修改的一半，再重新测试。\n", strings.Join(summary.RecentChangeFields, " / "))
	} else {
		b.WriteString("若出现反向问题，先回退最近一次改动的一半，再重新测试。\n")
	}
	return b.String()
}

func generateIssueSummaryWithDecisionEN(session storage.TelemetrySession, profile *storage.TuneProfile, evaluation *storage.RoadSessionEvaluation, summary *storage.SessionIssueSummary, decision *storage.RoadTuningDecision) string {
	var b strings.Builder
	b.WriteString("# FH6 Local Tuning Report\n\n")
	writeSessionEN(&b, session, profile)
	writeRoadEvaluationEN(&b, evaluation)
	writeRoadTuningDecisionEN(&b, decision, profile)
	b.WriteString("## 1. Merged Problems\n\n")
	for _, group := range summary.Groups {
		fmt.Fprintf(&b, "- %s: %s severity, %d events, %s total, previous run: %s.\n", issueFamilyLabelEN(group.Family), group.Severity, group.EventCount, duration(group.TotalDurationMS), issueComparisonEN(group.Comparison))
	}
	writeIssueSummaryEvidenceEN(&b, summary.Groups)
	b.WriteString("\n## Next Test\n\n")
	b.WriteString("Repeat the same car, track, and driver mode. Apply only the confirmed 1-3 actions, then check average speed, risk, and issue count/duration together.\n")
	b.WriteString("\n## Risk And Rollback\n\n")
	if decision != nil && decision.RollbackRecommended {
		b.WriteString("The model recommends rollback or reverse fine-tuning of the recent related change before adding new changes.\n")
	} else if len(summary.RecentChangeFields) > 0 {
		fmt.Fprintf(&b, "Recent changed fields: %s. If average speed, risk, or issue duration gets worse, roll back half of the related change and retest.\n", strings.Join(summary.RecentChangeFields, " / "))
	} else {
		b.WriteString("If the opposite problem appears, roll back half of the last change and test again.\n")
	}
	return b.String()
}

func generateIssueSummaryZH(session storage.TelemetrySession, profile *storage.TuneProfile, evaluation *storage.RoadSessionEvaluation, summary *storage.SessionIssueSummary) string {
	var b strings.Builder
	b.WriteString("# FH6 本地调校建议报告\n\n")
	writeSessionZH(&b, session, profile)
	writeRoadEvaluationZH(&b, evaluation)
	writeWholeCarPlanZH(&b, summary, profile)
	b.WriteString("## 一、合并问题项\n\n")
	for _, group := range summary.Groups {
		fmt.Fprintf(&b, "- %s：%s，出现 %d 次，累计 %s，对比上轮：%s。\n", issueFamilyLabelZH(group.Family), severityZH(group.Severity), group.EventCount, duration(group.TotalDurationMS), issueComparisonZH(group.Comparison))
	}
	writeIssueSummaryActionsZH(&b, summary.Groups, profile)
	writeIssueSummaryEvidenceZH(&b, summary.Groups)
	b.WriteString("\n## 四、下一轮测试方法\n\n")
	b.WriteString("使用同一车辆、同一赛道、同一驾驶方式重复测试。每轮优先修改 1-3 个主问题对应项目，并观察问题组次数、持续时间和严重度是否下降。\n")
	b.WriteString("\n## 五、风险与回退\n\n")
	if len(summary.RecentChangeFields) > 0 {
		fmt.Fprintf(&b, "本轮检测到最近修改项：%s。若问题恶化，优先回退相关修改的一半，再重新测试。\n", strings.Join(summary.RecentChangeFields, " / "))
	} else {
		b.WriteString("若出现反向问题，先回退最近一次改动的一半，再重新测试。\n")
	}
	return b.String()
}

func generateIssueSummaryEN(session storage.TelemetrySession, profile *storage.TuneProfile, evaluation *storage.RoadSessionEvaluation, summary *storage.SessionIssueSummary) string {
	var b strings.Builder
	b.WriteString("# FH6 Local Tuning Report\n\n")
	writeSessionEN(&b, session, profile)
	writeRoadEvaluationEN(&b, evaluation)
	writeWholeCarPlanEN(&b, summary, profile)
	b.WriteString("## 1. Merged Problems\n\n")
	for _, group := range summary.Groups {
		fmt.Fprintf(&b, "- %s: %s severity, %d events, %s total, previous run: %s.\n", issueFamilyLabelEN(group.Family), group.Severity, group.EventCount, duration(group.TotalDurationMS), issueComparisonEN(group.Comparison))
	}
	writeIssueSummaryActionsEN(&b, summary.Groups, profile)
	writeIssueSummaryEvidenceEN(&b, summary.Groups)
	b.WriteString("\n## 4. Next Test\n\n")
	b.WriteString("Repeat the same car, track, and driver mode. Change only 1-3 main items per run, then check whether issue group count, duration, and severity improve.\n")
	b.WriteString("\n## 5. Risk And Rollback\n\n")
	if len(summary.RecentChangeFields) > 0 {
		fmt.Fprintf(&b, "Recent changed fields: %s. If the issue gets worse, roll back half of the related change and retest.\n", strings.Join(summary.RecentChangeFields, " / "))
	} else {
		b.WriteString("If the opposite problem appears, roll back half of the last change and test again.\n")
	}
	return b.String()
}

func writeRoadTuningDecisionZH(b *strings.Builder, decision *storage.RoadTuningDecision, profile *storage.TuneProfile) {
	b.WriteString("## 公路调校决策\n\n")
	if decision == nil || decision.Status == "" || decision.Status == "no_matching_symptom" || decision.Status == "insufficient_data" {
		b.WriteString("- 当前数据不足，暂不生成明确调校决策。请绑定调校档案，并在同一车辆、赛道、驾驶方式下复测。\n\n")
		return
	}
	fmt.Fprintf(b, "- 主问题：%s（%s）\n", nonEmpty(decision.Symptom, decision.SymptomID), decision.Phase)
	fmt.Fprintf(b, "- 主因：%s\n", nonEmpty(decision.PrimaryCause, decision.Reason))
	fmt.Fprintf(b, "- 置信度：%s\n", decision.Confidence)
	fmt.Fprintf(b, "- 复测判断：%s\n", roadFitVerdictZH(decision.FitVerdict))
	if decision.RollbackRecommended {
		b.WriteString("- 回退提示：上轮修改后结果变差，优先回退或反向微调相关项目。\n")
	}
	if profile == nil {
		b.WriteString("- 当前未绑定调校档案，无法给出可直接应用的具体目标值。\n")
	}
	if len(decision.Actions) == 0 {
		b.WriteString("- 暂无可执行动作，请查看证据并继续采样。\n\n")
		return
	}
	b.WriteString("\n### 建议动作（最多 3 项）\n\n")
	for _, action := range decision.Actions {
		ranked := rankedAction{SuggestedAction: telemetry.SuggestedAction{
			Category:  action.Category,
			Item:      action.Item,
			Direction: action.Direction,
			Amount:    action.Amount,
			Reason:    action.Reason,
		}, Evidence: action.Evidence}
		if profile != nil && action.CanAutoApply {
			fmt.Fprintf(b, "- %s / %s：%s。原因：%s\n", categoryZH(action.Category), itemZH(action.Item), actionPlanZH(ranked, profile), action.Reason)
		} else {
			fmt.Fprintf(b, "- %s / %s：%s %s。原因：%s\n", categoryZH(action.Category), itemZH(action.Item), directionZH(action.Direction), amountZH(action.Amount), action.Reason)
		}
	}
	if len(decision.RetestFocus) > 0 {
		fmt.Fprintf(b, "\n复测观察点：%s。\n", strings.Join(decision.RetestFocus, " / "))
	}
	b.WriteString("\n")
}

func writeRoadTuningDecisionEN(b *strings.Builder, decision *storage.RoadTuningDecision, profile *storage.TuneProfile) {
	b.WriteString("## Road Tuning Decision\n\n")
	if decision == nil || decision.Status == "" || decision.Status == "no_matching_symptom" || decision.Status == "insufficient_data" {
		b.WriteString("- Not enough data for a clear tuning decision yet. Bind a tune profile and retest the same car, track, and driver mode.\n\n")
		return
	}
	fmt.Fprintf(b, "- Primary issue: %s (%s)\n", nonEmpty(decision.Symptom, decision.SymptomID), decision.Phase)
	fmt.Fprintf(b, "- Primary cause: %s\n", nonEmpty(decision.PrimaryCause, decision.Reason))
	fmt.Fprintf(b, "- Confidence: %s\n", decision.Confidence)
	fmt.Fprintf(b, "- Retest verdict: %s\n", roadFitVerdictEN(decision.FitVerdict))
	if decision.RollbackRecommended {
		b.WriteString("- Rollback: the last change made the result worse, so roll back or reverse fine-tune the related item first.\n")
	}
	if profile == nil {
		b.WriteString("- No tune profile is bound, so concrete target values cannot be generated.\n")
	}
	if len(decision.Actions) == 0 {
		b.WriteString("- No actionable tuning move yet. Review evidence and collect more comparable data.\n\n")
		return
	}
	b.WriteString("\n### Suggested Actions (max 3)\n\n")
	for _, action := range decision.Actions {
		ranked := rankedAction{SuggestedAction: telemetry.SuggestedAction{
			Category:  action.Category,
			Item:      action.Item,
			Direction: action.Direction,
			Amount:    action.Amount,
			Reason:    action.Reason,
		}, Evidence: action.Evidence}
		if profile != nil && action.CanAutoApply {
			fmt.Fprintf(b, "- %s / %s: %s. Reason: %s\n", title(action.Category), title(action.Item), actionPlanEN(ranked, profile), action.Reason)
		} else {
			fmt.Fprintf(b, "- %s / %s: %s %s. Reason: %s\n", title(action.Category), title(action.Item), title(action.Direction), action.Amount, action.Reason)
		}
	}
	if len(decision.RetestFocus) > 0 {
		fmt.Fprintf(b, "\nRetest focus: %s.\n", strings.Join(decision.RetestFocus, " / "))
	}
	b.WriteString("\n")
}

func nonEmpty(value string, fallback string) string {
	if strings.TrimSpace(value) != "" {
		return value
	}
	return fallback
}

func roadFitVerdictZH(value string) string {
	switch value {
	case "improved":
		return "遥测改善"
	case "worsened":
		return "遥测恶化"
	case "unchanged":
		return "基本持平"
	default:
		return "未知"
	}
}

func roadFitVerdictEN(value string) string {
	switch value {
	case "improved":
		return "telemetry improved"
	case "worsened":
		return "telemetry worsened"
	case "unchanged":
		return "unchanged"
	default:
		return "unknown"
	}
}

func writeIssueSummaryActionsZH(b *strings.Builder, groups []storage.SessionIssueGroup, profile *storage.TuneProfile) {
	b.WriteString("\n## 二、优先调校建议\n\n")
	limit := 3
	if len(groups) < limit {
		limit = len(groups)
	}
	for i := 0; i < limit; i++ {
		group := groups[i]
		fmt.Fprintf(b, "### %s\n\n", issueFamilyLabelZH(group.Family))
		if len(group.RelatedRecentChanges) > 0 {
			fmt.Fprintf(b, "- 与最近修改可能相关：%s\n", strings.Join(group.RelatedRecentChanges, " / "))
		}
		if group.AdjustmentStrategy != "" {
			fmt.Fprintf(b, "- 调整策略：%s\n", issueAdjustmentStrategyZH(group.AdjustmentStrategy))
		}
		if group.FeedbackDirective != "" {
			fmt.Fprintf(b, "- 反馈判断：%s\n", issueFeedbackDirectiveZH(group.FeedbackDirective))
		}
		if profile == nil {
			b.WriteString("- 需绑定调校档案后生成具体调整量。\n")
			continue
		}
		if len(group.PrimaryActions) == 0 {
			b.WriteString("- 暂无明确调校建议，建议继续采集同条件数据。\n")
			continue
		}
		for _, action := range group.PrimaryActions {
			ranked := rankedAction{SuggestedAction: action, Evidence: representativeEvidence(group)}
			fmt.Fprintf(b, "- %s / %s：%s\n", categoryZH(action.Category), itemZH(action.Item), actionPlanZH(ranked, profile))
			if note := actionExplanationNote(ranked); note != "" {
				fmt.Fprintf(b, "  调校说明：%s\n", note)
			}
		}
	}
}

func writeIssueSummaryActionsEN(b *strings.Builder, groups []storage.SessionIssueGroup, profile *storage.TuneProfile) {
	b.WriteString("\n## 2. Priority Adjustments\n\n")
	limit := 3
	if len(groups) < limit {
		limit = len(groups)
	}
	for i := 0; i < limit; i++ {
		group := groups[i]
		fmt.Fprintf(b, "### %s\n\n", issueFamilyLabelEN(group.Family))
		if len(group.RelatedRecentChanges) > 0 {
			fmt.Fprintf(b, "- May relate to recent changes: %s\n", strings.Join(group.RelatedRecentChanges, " / "))
		}
		if group.AdjustmentStrategy != "" {
			fmt.Fprintf(b, "- Adjustment strategy: %s\n", issueAdjustmentStrategyEN(group.AdjustmentStrategy))
		}
		if group.FeedbackDirective != "" {
			fmt.Fprintf(b, "- Feedback: %s\n", issueFeedbackDirectiveEN(group.FeedbackDirective))
		}
		if profile == nil {
			b.WriteString("- Bind a tune profile to generate concrete adjustment values.\n")
			continue
		}
		if len(group.PrimaryActions) == 0 {
			b.WriteString("- No clear tuning action yet. Capture more comparable data.\n")
			continue
		}
		for _, action := range group.PrimaryActions {
			ranked := rankedAction{SuggestedAction: action, Evidence: representativeEvidence(group)}
			fmt.Fprintf(b, "- %s / %s: %s\n", title(action.Category), title(action.Item), actionPlanEN(ranked, profile))
			if note := actionExplanationNote(ranked); note != "" {
				fmt.Fprintf(b, "  Tuning note: %s\n", note)
			}
		}
	}
}

func writeIssueSummaryEvidenceZH(b *strings.Builder, groups []storage.SessionIssueGroup) {
	b.WriteString("\n## 三、关键证据\n\n")
	for _, group := range groups {
		fmt.Fprintf(b, "### %s\n\n", issueFamilyLabelZH(group.Family))
		keys := sortedEvidenceKeys(group.Evidence)
		for _, key := range keys {
			stat := group.Evidence[key]
			fmt.Fprintf(b, "- %s：平均 %.2f，范围 %.2f - %.2f\n", labelZH(key), stat.Avg, stat.Min, stat.Max)
		}
	}
}

func writeIssueSummaryEvidenceEN(b *strings.Builder, groups []storage.SessionIssueGroup) {
	b.WriteString("\n## 3. Key Evidence\n\n")
	for _, group := range groups {
		fmt.Fprintf(b, "### %s\n\n", issueFamilyLabelEN(group.Family))
		keys := sortedEvidenceKeys(group.Evidence)
		for _, key := range keys {
			stat := group.Evidence[key]
			fmt.Fprintf(b, "- %s: avg %.2f, range %.2f - %.2f\n", labelEN(key), stat.Avg, stat.Min, stat.Max)
		}
	}
}

func writeWholeCarPlanZH(b *strings.Builder, summary *storage.SessionIssueSummary, profile *storage.TuneProfile) {
	if summary == nil || len(summary.WholeCarPlan.Actions) == 0 {
		return
	}
	plan := summary.WholeCarPlan
	b.WriteString("## 整车调校方案\n\n")
	fmt.Fprintf(b, "- 策略：%s\n", wholeCarStrategyZH(plan.Strategy))
	fmt.Fprintf(b, "- 可信度：%s\n", wholeCarConfidenceZH(plan.Confidence))
	if summary.GearPower.Summary != "" {
		fmt.Fprintf(b, "- 齿比动力诊断：%s", gearFindingZH(summary.GearPower.Summary))
		if summary.GearPower.LaunchFinding != "" {
			fmt.Fprintf(b, " / %s", gearFindingZH(summary.GearPower.LaunchFinding))
		}
		if summary.GearPower.TopSpeedFinding != "" {
			fmt.Fprintf(b, " / %s", gearFindingZH(summary.GearPower.TopSpeedFinding))
		}
		if summary.GearPower.PowerToWeightKWPerKG > 0 {
			fmt.Fprintf(b, " / 功率比重 %.4f kW/kg", summary.GearPower.PowerToWeightKWPerKG)
		}
		if summary.GearPower.TractionLimitedPercent > 0 {
			fmt.Fprintf(b, " / 牵引受限 %.0f%%", summary.GearPower.TractionLimitedPercent*100)
		}
		b.WriteString("\n")
	}
	for _, action := range plan.Actions {
		ranked := rankedAction{SuggestedAction: telemetry.SuggestedAction{
			Priority:  action.Priority,
			Category:  action.Category,
			Item:      action.Item,
			Direction: action.Direction,
			Amount:    action.Amount,
			Reason:    action.Reason,
		}, Evidence: action.Evidence}
		if profile == nil {
			fmt.Fprintf(b, "- %s / %s：%s %s\n", categoryZH(action.Category), itemZH(action.Item), directionZH(action.Direction), amountZH(action.Amount))
			continue
		}
		fmt.Fprintf(b, "- %s / %s：%s\n", categoryZH(action.Category), itemZH(action.Item), actionPlanZH(ranked, profile))
	}
	if len(plan.Conflicts) > 0 {
		b.WriteString("- 已解决冲突：\n")
		for _, conflict := range plan.Conflicts {
			fmt.Fprintf(b, "  - %s：保留 %s，移除 %s\n", conflict.Key, conflict.KeptItem, conflict.DroppedItem)
		}
	}
	b.WriteString("\n")
}

func writeWholeCarPlanEN(b *strings.Builder, summary *storage.SessionIssueSummary, profile *storage.TuneProfile) {
	if summary == nil || len(summary.WholeCarPlan.Actions) == 0 {
		return
	}
	plan := summary.WholeCarPlan
	b.WriteString("## Whole-car Tuning Plan\n\n")
	fmt.Fprintf(b, "- Strategy: %s\n", wholeCarStrategyEN(plan.Strategy))
	fmt.Fprintf(b, "- Confidence: %s\n", wholeCarConfidenceEN(plan.Confidence))
	if summary.GearPower.Summary != "" {
		fmt.Fprintf(b, "- Gear power diagnostic: %s", gearFindingEN(summary.GearPower.Summary))
		if summary.GearPower.LaunchFinding != "" {
			fmt.Fprintf(b, " / %s", gearFindingEN(summary.GearPower.LaunchFinding))
		}
		if summary.GearPower.TopSpeedFinding != "" {
			fmt.Fprintf(b, " / %s", gearFindingEN(summary.GearPower.TopSpeedFinding))
		}
		if summary.GearPower.PowerToWeightKWPerKG > 0 {
			fmt.Fprintf(b, " / power-to-weight %.4f kW/kg", summary.GearPower.PowerToWeightKWPerKG)
		}
		if summary.GearPower.TractionLimitedPercent > 0 {
			fmt.Fprintf(b, " / traction-limited %.0f%%", summary.GearPower.TractionLimitedPercent*100)
		}
		b.WriteString("\n")
	}
	for _, action := range plan.Actions {
		ranked := rankedAction{SuggestedAction: telemetry.SuggestedAction{
			Priority:  action.Priority,
			Category:  action.Category,
			Item:      action.Item,
			Direction: action.Direction,
			Amount:    action.Amount,
			Reason:    action.Reason,
		}, Evidence: action.Evidence}
		if profile == nil {
			fmt.Fprintf(b, "- %s / %s: %s %s\n", title(action.Category), title(action.Item), title(action.Direction), action.Amount)
			continue
		}
		fmt.Fprintf(b, "- %s / %s: %s\n", title(action.Category), title(action.Item), actionPlanEN(ranked, profile))
	}
	if len(plan.Conflicts) > 0 {
		b.WriteString("- Resolved conflicts:\n")
		for _, conflict := range plan.Conflicts {
			fmt.Fprintf(b, "  - %s: kept %s, dropped %s\n", conflict.Key, conflict.KeptItem, conflict.DroppedItem)
		}
	}
	b.WriteString("\n")
}

func representativeEvidence(group storage.SessionIssueGroup) map[string]float64 {
	out := make(map[string]float64, len(group.Evidence))
	for key, stat := range group.Evidence {
		out[key] = stat.Avg
	}
	return out
}

func sortedEvidenceKeys(evidence map[string]storage.IssueEvidence) []string {
	keys := make([]string, 0, len(evidence))
	for key := range evidence {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	if len(keys) > 6 {
		return keys[:6]
	}
	return keys
}

func issueFamilyLabelZH(family string) string {
	switch family {
	case "launch_traction":
		return "起步牵引"
	case "gearing_acceleration":
		return "齿比 / 加速"
	case "brake_balance":
		return "刹车平衡"
	case "corner_entry_balance":
		return "入弯平衡"
	case "mid_corner_balance":
		return "持续过弯平衡"
	case "corner_exit_power":
		return "出弯动力"
	case "suspension_platform":
		return "悬挂平台"
	case "tire_temperature_stability":
		return "轮胎温度 / 稳定性"
	default:
		return family
	}
}

func issueFamilyLabelEN(family string) string {
	switch family {
	case "launch_traction":
		return "Launch traction"
	case "gearing_acceleration":
		return "Gearing / acceleration"
	case "brake_balance":
		return "Brake balance"
	case "corner_entry_balance":
		return "Corner entry balance"
	case "mid_corner_balance":
		return "Sustained cornering balance"
	case "corner_exit_power":
		return "Corner exit power"
	case "suspension_platform":
		return "Suspension platform"
	case "tire_temperature_stability":
		return "Tire temperature / stability"
	default:
		return title(family)
	}
}

func issueComparisonZH(value string) string {
	switch value {
	case "improved":
		return "改善"
	case "worsened":
		return "恶化"
	case "unchanged":
		return "持平"
	default:
		return "无法比较"
	}
}

func issueComparisonEN(value string) string {
	switch value {
	case "improved":
		return "improved"
	case "worsened":
		return "worsened"
	case "unchanged":
		return "unchanged"
	default:
		return "unavailable"
	}
}

func issueAdjustmentStrategyZH(value string) string {
	switch value {
	case "rollback_first":
		return "先回退相关修改"
	case "coarse_combination":
		return "大步组合调整"
	case "medium_combination":
		return "中等组合调整"
	case "fine_tune":
		return "小步微调"
	default:
		return value
	}
}

func issueAdjustmentStrategyEN(value string) string {
	switch value {
	case "rollback_first":
		return "rollback related changes first"
	case "coarse_combination":
		return "coarse combined adjustment"
	case "medium_combination":
		return "medium combined adjustment"
	case "fine_tune":
		return "fine tune"
	default:
		return value
	}
}

func issueFeedbackDirectiveZH(value string) string {
	switch value {
	case "rollback_related_changes":
		return "相关的上次修改后问题变差，先回退部分修改，不继续同方向加码"
	case "keep_direction_then_fine_tune":
		return "当前方向已有改善，后续改用更小步进微调"
	case "avoid_more_same_direction":
		return "问题没有明确改善，先避免继续同方向调整，复测确认"
	default:
		return value
	}
}

func issueFeedbackDirectiveEN(value string) string {
	switch value {
	case "rollback_related_changes":
		return "the last related change made this worse; roll part of it back before adding more"
	case "keep_direction_then_fine_tune":
		return "the direction improved the result; continue with smaller steps"
	case "avoid_more_same_direction":
		return "the issue did not clearly improve; avoid adding more in the same direction until retested"
	default:
		return value
	}
}

func wholeCarStrategyZH(value string) string {
	switch value {
	case "rollback_first":
		return "优先回退"
	case "coarse_whole_car":
		return "整车大步调整"
	case "targeted_whole_car":
		return "整车定向调整"
	default:
		return value
	}
}

func wholeCarStrategyEN(value string) string {
	switch value {
	case "rollback_first":
		return "rollback first"
	case "coarse_whole_car":
		return "coarse whole-car pass"
	case "targeted_whole_car":
		return "targeted whole-car pass"
	default:
		return value
	}
}

func wholeCarConfidenceZH(value string) string {
	switch value {
	case "high":
		return "高"
	case "medium":
		return "中"
	case "low":
		return "低"
	case "needs_profile":
		return "需要绑定调校档案"
	default:
		return value
	}
}

func wholeCarConfidenceEN(value string) string {
	switch value {
	case "high":
		return "high"
	case "medium":
		return "medium"
	case "low":
		return "low"
	case "needs_profile":
		return "needs tune profile"
	default:
		return value
	}
}

func gearFindingZH(value string) string {
	switch value {
	case "not_enough_samples":
		return "样本不足"
	case "gearing_window_ok":
		return "齿比动力区间正常"
	case "gearing_adjustment_needed":
		return "需要调整齿比动力区间"
	case "traction_limited_power":
		return "动力输出受牵引限制"
	case "too_long":
		return "负载下齿比偏长"
	case "too_short":
		return "负载下齿比偏短"
	case "traction_limited":
		return "牵引受限"
	case "top_speed_limited_by_gearing":
		return "极速受齿比限制"
	case "top_speed_bog_down":
		return "高挡齿比偏长"
	case "top_speed_ok":
		return "高速齿比正常"
	case "launch_wheelspin":
		return "起步打滑"
	case "launch_bog_down":
		return "起步憋转"
	default:
		return value
	}
}

func gearFindingEN(value string) string {
	switch value {
	case "not_enough_samples":
		return "not enough samples"
	case "gearing_window_ok":
		return "gearing window OK"
	case "gearing_adjustment_needed":
		return "gearing adjustment needed"
	case "traction_limited_power":
		return "traction-limited power"
	case "too_long":
		return "too long under load"
	case "too_short":
		return "too short under load"
	case "traction_limited":
		return "traction limited"
	case "top_speed_limited_by_gearing":
		return "top speed limited by gearing"
	case "top_speed_bog_down":
		return "top gear too long"
	case "top_speed_ok":
		return "top speed gearing OK"
	case "launch_wheelspin":
		return "launch wheelspin"
	case "launch_bog_down":
		return "launch bog down"
	default:
		return value
	}
}

func generateZH(session storage.TelemetrySession, profile *storage.TuneProfile, events []telemetry.DetectedEvent, evaluation *storage.RoadSessionEvaluation) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# FH6 本地调校建议报告\n\n")
	writeSessionZH(&b, session, profile)
	writeRoadEvaluationZH(&b, evaluation)
	if len(events) == 0 {
		b.WriteString("## 一、检测结果\n\n当前会话未检测到明确调校事件。建议继续使用相同路线和驾驶方式采集更多数据。\n\n")
		b.WriteString("## 二、下一轮测试方法\n\n重复当前测试路线 2-3 次，确认是否能稳定复现推头、甩尾、抱死、起步打滑或悬挂触底。\n")
		return b.String()
	}

	b.WriteString("## 一、检测结果\n\n")
	for _, event := range events {
		fmt.Fprintf(&b, "- %s：%s，持续 %s，阶段 %s。\n", eventLabelZH(event.Type), severityZH(event.Severity), duration(event.DurationMS), segmentZH(event.Segment))
	}
	b.WriteString("\n## 二、数据证据\n\n")
	for _, event := range events {
		fmt.Fprintf(&b, "### %s\n\n", eventLabelZH(event.Type))
		writeEvidence(&b, event.Evidence, labelZH)
	}

	b.WriteString("## 三、问题成因\n\n")
	for _, event := range events {
		fmt.Fprintf(&b, "- %s：%s\n", eventLabelZH(event.Type), causeZH(event.Type))
	}

	primary, secondary := rankedActions(events)
	b.WriteString("\n## 四、优先调整建议\n\n")
	writeActionsZH(&b, primary, profile)
	if len(secondary) > 0 {
		b.WriteString("\n## 五、暂不优先调整的项目\n\n")
		writeActionsZH(&b, secondary, profile)
	}

	b.WriteString("\n## 六、下一轮测试方法\n\n")
	b.WriteString("使用相同车辆、相同路段和相近入弯/出弯方式重复 3 次测试，优先观察本次最高严重度事件是否减少，Slip/CombinedSlip 峰值是否下降。\n")
	b.WriteString("\n## 七、风险与回退\n\n")
	b.WriteString("每轮最多修改 1-3 个主项。如果出现反向问题，例如由甩尾变成推头、由前抱死变成后抱死，应先回退最近一次改动的一半，再重新测试。\n")
	return b.String()
}

func generateEN(session storage.TelemetrySession, profile *storage.TuneProfile, events []telemetry.DetectedEvent, evaluation *storage.RoadSessionEvaluation) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# FH6 Local Tuning Report\n\n")
	writeSessionEN(&b, session, profile)
	writeRoadEvaluationEN(&b, evaluation)
	if len(events) == 0 {
		b.WriteString("## 1. Detected Results\n\nNo clear tuning events were detected in this session. Repeat the same route or test segment to build evidence.\n\n")
		b.WriteString("## 2. Next Test\n\nRepeat the same route 2-3 times and watch for understeer, oversteer, brake lockup, launch wheelspin, or suspension bottoming.\n")
		return b.String()
	}

	b.WriteString("## 1. Detected Results\n\n")
	for _, event := range events {
		fmt.Fprintf(&b, "- %s: %s severity, %s, %s.\n", eventLabelEN(event.Type), event.Severity, duration(event.DurationMS), segmentEN(event.Segment))
	}
	b.WriteString("\n## 2. Evidence\n\n")
	for _, event := range events {
		fmt.Fprintf(&b, "### %s\n\n", eventLabelEN(event.Type))
		writeEvidence(&b, event.Evidence, labelEN)
	}

	b.WriteString("## 3. Likely Cause\n\n")
	for _, event := range events {
		fmt.Fprintf(&b, "- %s: %s\n", eventLabelEN(event.Type), causeEN(event.Type))
	}

	primary, secondary := rankedActions(events)
	b.WriteString("\n## 4. Priority Adjustments\n\n")
	writeActionsEN(&b, primary, profile)
	if len(secondary) > 0 {
		b.WriteString("\n## 5. Lower Priority Items\n\n")
		writeActionsEN(&b, secondary, profile)
	}

	b.WriteString("\n## 6. Next Test\n\n")
	b.WriteString("Repeat the same route 3 times with similar entry and exit inputs. Watch whether the highest-severity event count drops and whether Slip/CombinedSlip peaks improve.\n")
	b.WriteString("\n## 7. Risk And Rollback\n\n")
	b.WriteString("Change no more than 1-3 main items per run. If the opposite problem appears, roll back half of the last change and test again.\n")
	return b.String()
}

func writeSessionZH(b *strings.Builder, session storage.TelemetrySession, profile *storage.TuneProfile) {
	b.WriteString("## 会话概览\n\n")
	if profile != nil {
		fmt.Fprintf(b, "- 当前档案：%s", profile.CarName)
		if profile.VersionName != "" {
			fmt.Fprintf(b, " / %s", profile.VersionName)
		}
		fmt.Fprintln(b)
		if profile.CarClass != "" || profile.Drivetrain != "" || profile.UseCase != "" {
			fmt.Fprintf(b, "- 车辆信息：%s %s %s\n", profile.CarClass, profile.Drivetrain, profile.UseCase)
		}
	} else {
		b.WriteString("- 当前档案：未关联\n")
	}
	fmt.Fprintf(b, "- 测试模式：%s\n", sessionModeLabelZH(session.GameMode))
	fmt.Fprintf(b, "- 测试条件：%s\n", testConditionsZH(session))
	if storage.TestConditionsContainUnknown(storage.SessionTestConditions(session)) {
		b.WriteString("- 可信度提示：辅助设置存在未知项，相关驾驶输入和会话对比只能作为参考。\n")
	}
	fmt.Fprintf(b, "- 会话：%s\n", fallback(session.SessionName, fmt.Sprintf("#%d", session.ID)))
	if session.AvgSpeedKmh != nil {
		fmt.Fprintf(b, "- 平均速度：%.1f km/h\n", *session.AvgSpeedKmh)
	}
	if session.MaxSpeedKmh != nil {
		fmt.Fprintf(b, "- 最高速度：%.1f km/h\n", *session.MaxSpeedKmh)
	}
	b.WriteString("\n")
}

func writeSessionEN(b *strings.Builder, session storage.TelemetrySession, profile *storage.TuneProfile) {
	b.WriteString("## Session Overview\n\n")
	if profile != nil {
		fmt.Fprintf(b, "- Tune profile: %s", profile.CarName)
		if profile.VersionName != "" {
			fmt.Fprintf(b, " / %s", profile.VersionName)
		}
		fmt.Fprintln(b)
		if profile.CarClass != "" || profile.Drivetrain != "" || profile.UseCase != "" {
			fmt.Fprintf(b, "- Vehicle: %s %s %s\n", profile.CarClass, profile.Drivetrain, profile.UseCase)
		}
	} else {
		b.WriteString("- Tune profile: none\n")
	}
	fmt.Fprintf(b, "- Test mode: %s\n", sessionModeLabelEN(session.GameMode))
	fmt.Fprintf(b, "- Test conditions: %s\n", testConditionsEN(session))
	if storage.TestConditionsContainUnknown(storage.SessionTestConditions(session)) {
		b.WriteString("- Confidence note: some assist settings are unknown, so input-based conclusions and session comparisons are lower confidence.\n")
	}
	fmt.Fprintf(b, "- Session: %s\n", fallback(session.SessionName, fmt.Sprintf("#%d", session.ID)))
	if session.AvgSpeedKmh != nil {
		fmt.Fprintf(b, "- Average speed: %.1f km/h\n", *session.AvgSpeedKmh)
	}
	if session.MaxSpeedKmh != nil {
		fmt.Fprintf(b, "- Max speed: %.1f km/h\n", *session.MaxSpeedKmh)
	}
	b.WriteString("\n")
}

func writeRoadEvaluationZH(b *strings.Builder, evaluation *storage.RoadSessionEvaluation) {
	if evaluation == nil {
		return
	}
	b.WriteString("## 公路赛车评估\n\n")
	fmt.Fprintf(b, "- 综合结论：%s\n", roadVerdictLabelZH(evaluation.OverallVerdict))
	fmt.Fprintf(b, "- 纸面性能：%.0f / 100\n", evaluation.PaperPerformanceScore)
	fmt.Fprintf(b, "- 玩家适配：%.0f / 100\n", evaluation.PlayerFitScore)
	fmt.Fprintf(b, "- 失控风险：%.0f / 100\n", evaluation.RiskScore)
	fmt.Fprintf(b, "- 自动驾驶基线：%s\n", roadBaselineLabelZH(evaluation.BaselineStatus))
	if evaluation.BestRun != nil {
		fmt.Fprintf(b, "- 当前最佳标准赛段：%s，%s，置信度 %.0f%%\n", fallback(evaluation.BestRun.TrackName, fmt.Sprintf("#%d", evaluation.BestRun.TrackID)), duration(evaluation.BestRun.DurationMS), evaluation.BestRun.Confidence*100)
	}
	if evaluation.BaselineRun != nil {
		fmt.Fprintf(b, "- 匹配基线：%s，%s\n", fallback(evaluation.BaselineRun.TrackName, fmt.Sprintf("#%d", evaluation.BaselineRun.TrackID)), duration(evaluation.BaselineRun.DurationMS))
	}
	if evaluation.BaselineStatus == "missing_auto_baseline" {
		b.WriteString("- 当前只可评估玩家表现，不能判断纸面基线。建议用自动驾驶在同车、同级、同赛道跑一次作为基准。\n")
	}
	if len(evaluation.Attributions) > 0 {
		b.WriteString("- 问题归因：\n")
		for _, item := range evaluation.Attributions {
			label := roadAttributionLabelZH(item.Type)
			if item.EventType != "" {
				fmt.Fprintf(b, "  - %s / %s：%d 次", label, eventLabelZH(item.EventType), item.Count)
			} else {
				fmt.Fprintf(b, "  - %s：%s", label, roadBaselineLabelZH(item.Message))
			}
			if item.PrioritizeTuning {
				b.WriteString("，建议优先检查调校")
			}
			b.WriteString("\n")
		}
	}
	b.WriteString("\n")
}

func writeRoadEvaluationEN(b *strings.Builder, evaluation *storage.RoadSessionEvaluation) {
	if evaluation == nil {
		return
	}
	b.WriteString("## Road Racing Evaluation\n\n")
	fmt.Fprintf(b, "- Verdict: %s\n", roadVerdictLabelEN(evaluation.OverallVerdict))
	fmt.Fprintf(b, "- Paper performance: %.0f / 100\n", evaluation.PaperPerformanceScore)
	fmt.Fprintf(b, "- Player fit: %.0f / 100\n", evaluation.PlayerFitScore)
	fmt.Fprintf(b, "- Risk: %.0f / 100\n", evaluation.RiskScore)
	fmt.Fprintf(b, "- Auto baseline: %s\n", roadBaselineLabelEN(evaluation.BaselineStatus))
	if evaluation.BestRun != nil {
		fmt.Fprintf(b, "- Best standard segment: %s, %s, %.0f%% confidence\n", fallback(evaluation.BestRun.TrackName, fmt.Sprintf("#%d", evaluation.BestRun.TrackID)), duration(evaluation.BestRun.DurationMS), evaluation.BestRun.Confidence*100)
	}
	if evaluation.BaselineRun != nil {
		fmt.Fprintf(b, "- Matched baseline: %s, %s\n", fallback(evaluation.BaselineRun.TrackName, fmt.Sprintf("#%d", evaluation.BaselineRun.TrackID)), duration(evaluation.BaselineRun.DurationMS))
	}
	if evaluation.BaselineStatus == "missing_auto_baseline" {
		b.WriteString("- Only player performance can be evaluated right now; paper baseline cannot be judged without an auto-driver run on the same car, class, and standard segment.\n")
	}
	if len(evaluation.Attributions) > 0 {
		b.WriteString("- Attribution:\n")
		for _, item := range evaluation.Attributions {
			label := roadAttributionLabelEN(item.Type)
			if item.EventType != "" {
				fmt.Fprintf(b, "  - %s / %s: %d occurrence(s)", label, eventLabelEN(item.EventType), item.Count)
			} else {
				fmt.Fprintf(b, "  - %s: %s", label, roadBaselineLabelEN(item.Message))
			}
			if item.PrioritizeTuning {
				b.WriteString(", prioritize tuning review")
			}
			b.WriteString("\n")
		}
	}
	b.WriteString("\n")
}

func writeEvidence(b *strings.Builder, evidence map[string]float64, label func(string) string) {
	keys := make([]string, 0, len(evidence))
	for key := range evidence {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		fmt.Fprintf(b, "- %s：%.2f\n", label(key), evidence[key])
	}
	b.WriteString("\n")
}

func writeActionsZH(b *strings.Builder, actions []rankedAction, profile *storage.TuneProfile) {
	if len(actions) == 0 {
		b.WriteString("- 暂无明确调整项。\n")
		return
	}
	for i, action := range actions {
		fmt.Fprintf(b, "%d. %s / %s：%s。%s\n", i+1, categoryZH(action.Category), itemZH(action.Item), actionPlanZH(action, profile), reasonZH(action.Reason))
		if note := actionExplanationNote(action); note != "" {
			fmt.Fprintf(b, "   - 调校说明：%s\n", note)
		}
	}
}

func writeActionsEN(b *strings.Builder, actions []rankedAction, profile *storage.TuneProfile) {
	if len(actions) == 0 {
		b.WriteString("- No clear adjustment item yet.\n")
		return
	}
	for i, action := range actions {
		fmt.Fprintf(b, "%d. %s / %s: %s. %s\n", i+1, title(action.Category), title(action.Item), actionPlanEN(action, profile), sentence(action.Reason))
		if note := actionExplanationNote(action); note != "" {
			fmt.Fprintf(b, "   - Tuning note: %s\n", note)
		}
	}
}

func sessionModeLabelZH(mode string) string {
	switch telemetry.NormalizeGameMode(mode) {
	case telemetry.GameModeRace:
		return "比赛模式"
	case telemetry.GameModeFreeRoam:
		return "漫游模式"
	case telemetry.GameModeMixed:
		return "混合模式"
	case telemetry.GameModeMenu:
		return "菜单 / 过场"
	default:
		return "未知"
	}
}

func sessionModeLabelEN(mode string) string {
	switch telemetry.NormalizeGameMode(mode) {
	case telemetry.GameModeRace:
		return "Race"
	case telemetry.GameModeFreeRoam:
		return "Free roam"
	case telemetry.GameModeMixed:
		return "Mixed"
	case telemetry.GameModeMenu:
		return "Menu / transition"
	default:
		return "Unknown"
	}
}

func testConditionsZH(session storage.TelemetrySession) string {
	conditions := storage.SessionTestConditions(session)
	return strings.Join([]string{
		"驾驶=" + conditionLabelZH("driverMode", conditions.DriverMode),
		"刹车=" + conditionLabelZH("brakeAssist", conditions.BrakeAssist),
		"转向=" + conditionLabelZH("steeringAssist", conditions.SteeringAssist),
		"TCS=" + conditionLabelZH("toggle", conditions.TractionControl),
		"STM=" + conditionLabelZH("toggle", conditions.StabilityControl),
		"换挡=" + conditionLabelZH("shifting", conditions.Shifting),
		"起跑控制=" + conditionLabelZH("toggle", conditions.LaunchControl),
	}, " / ")
}

func testConditionsEN(session storage.TelemetrySession) string {
	conditions := storage.SessionTestConditions(session)
	return strings.Join([]string{
		"Driver=" + conditionLabelEN("driverMode", conditions.DriverMode),
		"Brake=" + conditionLabelEN("brakeAssist", conditions.BrakeAssist),
		"Steering=" + conditionLabelEN("steeringAssist", conditions.SteeringAssist),
		"TCS=" + conditionLabelEN("toggle", conditions.TractionControl),
		"STM=" + conditionLabelEN("toggle", conditions.StabilityControl),
		"Shifting=" + conditionLabelEN("shifting", conditions.Shifting),
		"Launch control=" + conditionLabelEN("toggle", conditions.LaunchControl),
	}, " / ")
}

func conditionLabelZH(kind string, value string) string {
	switch kind + ":" + strings.TrimSpace(value) {
	case "driverMode:player":
		return "玩家"
	case "driverMode:auto":
		return "自动驾驶"
	case "brakeAssist:assisted":
		return "辅助"
	case "brakeAssist:abs_on":
		return "ABS 开"
	case "brakeAssist:abs_off":
		return "ABS 关"
	case "steeringAssist:auto":
		return "自动转向"
	case "steeringAssist:assisted":
		return "辅助"
	case "steeringAssist:standard":
		return "标准"
	case "steeringAssist:simulation":
		return "拟真"
	case "toggle:on":
		return "开"
	case "toggle:off":
		return "关"
	case "shifting:automatic":
		return "自动"
	case "shifting:manual":
		return "手动"
	default:
		return "未知"
	}
}

func conditionLabelEN(kind string, value string) string {
	switch kind + ":" + strings.TrimSpace(value) {
	case "driverMode:player":
		return "Player"
	case "driverMode:auto":
		return "Auto driver"
	case "brakeAssist:assisted":
		return "Assisted"
	case "brakeAssist:abs_on":
		return "ABS on"
	case "brakeAssist:abs_off":
		return "ABS off"
	case "steeringAssist:auto":
		return "Auto steering"
	case "steeringAssist:assisted":
		return "Assisted"
	case "steeringAssist:standard":
		return "Standard"
	case "steeringAssist:simulation":
		return "Simulation"
	case "toggle:on":
		return "On"
	case "toggle:off":
		return "Off"
	case "shifting:automatic":
		return "Automatic"
	case "shifting:manual":
		return "Manual"
	default:
		return "Unknown"
	}
}

func roadVerdictLabelZH(value string) string {
	return mapLookup(value, map[string]string{
		"good_fit":           "好车：快且可控",
		"fast_but_risky":     "速度快但风险高",
		"paper_fast_not_fit": "纸面快但不适配",
		"needs_tuning":       "需要继续调校",
		"insufficient_data":  "数据不足",
	})
}

func roadVerdictLabelEN(value string) string {
	return mapLookup(value, map[string]string{
		"good_fit":           "Good fit",
		"fast_but_risky":     "Fast but risky",
		"paper_fast_not_fit": "Paper-fast but not a fit",
		"needs_tuning":       "Needs tuning",
		"insufficient_data":  "Insufficient data",
	})
}

func roadBaselineLabelZH(value string) string {
	return mapLookup(value, map[string]string{
		"matched_auto_baseline":      "已匹配自动驾驶基线",
		"self_auto_baseline":         "本会话是自动驾驶基线",
		"missing_auto_baseline":      "缺少自动驾驶基线",
		"missing_vehicle_identity":   "缺少车辆身份",
		"no_valid_standard_run":      "没有有效标准赛段",
		"no_standard_track":          "没有标准赛道",
		"benchmark_run_has_warnings": "赛段存在诊断警告",
	})
}

func roadBaselineLabelEN(value string) string {
	return mapLookup(value, map[string]string{
		"matched_auto_baseline":      "Matched auto baseline",
		"self_auto_baseline":         "This session is the auto baseline",
		"missing_auto_baseline":      "Missing auto baseline",
		"missing_vehicle_identity":   "Missing vehicle identity",
		"no_valid_standard_run":      "No valid standard segment",
		"no_standard_track":          "No standard track",
		"benchmark_run_has_warnings": "Benchmark segment has warnings",
	})
}

func roadAttributionLabelZH(value string) string {
	return mapLookup(value, map[string]string{
		"tune_issue":             "调校问题",
		"style_fit_issue":        "驾驶风格适配",
		"driver_execution_issue": "驾驶执行",
		"data_gap":               "数据缺口",
	})
}

func roadAttributionLabelEN(value string) string {
	return mapLookup(value, map[string]string{
		"tune_issue":             "Tune issue",
		"style_fit_issue":        "Driving style fit",
		"driver_execution_issue": "Driver execution",
		"data_gap":               "Data gap",
	})
}

type rankedAction struct {
	telemetry.SuggestedAction
	Evidence map[string]float64
}

func rankedActions(events []telemetry.DetectedEvent) ([]rankedAction, []rankedAction) {
	all := make([]rankedAction, 0)
	seen := map[string]bool{}
	for _, event := range events {
		for _, action := range event.SuggestedActions {
			key := action.Category + "/" + action.Item + "/" + action.Direction
			if seen[key] {
				continue
			}
			seen[key] = true
			all = append(all, rankedAction{SuggestedAction: action, Evidence: event.Evidence})
		}
	}
	sort.SliceStable(all, func(i, j int) bool {
		return all[i].Priority < all[j].Priority
	})
	if len(all) <= 3 {
		return all, nil
	}
	return all[:3], all[3:]
}

type tuneFieldRef struct {
	key     string
	labelZH string
	labelEN string
	unit    string
	step    float64
	value   *float64
}

func actionPlanZH(action rankedAction, profile *storage.TuneProfile) string {
	if adjustments := concreteAdjustments(action, profile); len(adjustments) > 0 {
		parts := make([]string, 0, len(adjustments))
		for _, adjustment := range adjustments {
			parts = append(parts, fmt.Sprintf("%s：%s -> %s，%s %s", adjustment.labelZH, formatTuneNumber(adjustment.current, adjustment.step), formatTuneNumber(adjustment.target, adjustment.step), directionDeltaZH(adjustment.delta, action), formatTuneNumber(abs(adjustment.delta), adjustment.step)))
		}
		return strings.Join(parts, "；")
	}
	return strings.TrimSpace(directionZH(action.Direction) + " " + amountZH(action.Amount))
}

func actionPlanEN(action rankedAction, profile *storage.TuneProfile) string {
	if adjustments := concreteAdjustments(action, profile); len(adjustments) > 0 {
		parts := make([]string, 0, len(adjustments))
		for _, adjustment := range adjustments {
			parts = append(parts, fmt.Sprintf("%s: %s -> %s, %s %s", adjustment.labelEN, formatTuneNumber(adjustment.current, adjustment.step), formatTuneNumber(adjustment.target, adjustment.step), directionDeltaEN(adjustment.delta, action), formatTuneNumber(abs(adjustment.delta), adjustment.step)))
		}
		return strings.Join(parts, "; ")
	}
	return strings.TrimSpace(title(action.Direction) + " " + action.Amount)
}

func actionExplanationNote(action rankedAction) string {
	gear := int(action.Evidence["gear"])
	explanations := storage.TuneAdjustmentExplanationsForAction(action.Item, gear)
	if len(explanations) == 0 {
		return ""
	}
	seen := map[string]bool{}
	notes := make([]string, 0, len(explanations))
	for _, explanation := range explanations {
		if strings.TrimSpace(explanation.Description) == "" || seen[explanation.Description] {
			continue
		}
		seen[explanation.Description] = true
		notes = append(notes, explanation.Description)
	}
	return strings.Join(notes, "；")
}

func concreteAdjustments(action rankedAction, profile *storage.TuneProfile) []tuneAdjustment {
	if profile == nil {
		return nil
	}
	fields := actionTuneFields(action, profile)
	adjustments := make([]tuneAdjustment, 0, len(fields))
	for _, field := range fields {
		if field.value == nil {
			continue
		}
		delta, ok := actionDelta(action, *field.value, field.step, field.unit)
		if !ok || delta == 0 {
			continue
		}
		adjustments = append(adjustments, tuneAdjustment{
			labelZH: field.labelZH,
			labelEN: field.labelEN,
			current: *field.value,
			target:  roundToStep(*field.value+delta, field.step),
			delta:   roundToStep(delta, field.step),
			step:    field.step,
		})
	}
	return adjustments
}

type tuneAdjustment struct {
	labelZH string
	labelEN string
	current float64
	target  float64
	delta   float64
	step    float64
}

func actionTuneFields(action rankedAction, profile *storage.TuneProfile) []tuneFieldRef {
	fields := tuneFieldMap(profile)
	switch action.Item {
	case "front_tire_pressure":
		return pickFields(fields, "frontTirePressure")
	case "rear_tire_pressure":
		return pickFields(fields, "rearTirePressure")
	case "gear_1":
		return pickFields(fields, "gear1")
	case "gear_2":
		return pickFields(fields, "gear2")
	case "gear_3":
		return pickFields(fields, "gear3")
	case "gear_4":
		return pickFields(fields, "gear4")
	case "gear_5":
		return pickFields(fields, "gear5")
	case "gear_6":
		return pickFields(fields, "gear6")
	case "gear_7":
		return pickFields(fields, "gear7")
	case "gear_8":
		return pickFields(fields, "gear8")
	case "gear_9":
		return pickFields(fields, "gear9")
	case "gear_10":
		return pickFields(fields, "gear10")
	case "current_gear":
		gear := int(action.Evidence["gear"])
		if gear < 1 || gear > 10 {
			return nil
		}
		return pickFields(fields, fmt.Sprintf("gear%d", gear))
	case "final_drive":
		return pickFields(fields, "finalDrive")
	case "brake_balance":
		return pickFields(fields, "brakeBalance")
	case "brake_pressure":
		return pickFields(fields, "brakePressure")
	case "rear_diff_accel":
		return pickFields(fields, "rearDiffAccel")
	case "rear_diff_decel":
		return pickFields(fields, "rearDiffDecel")
	case "front_diff_accel":
		return pickFields(fields, "frontDiffAccel")
	case "front_diff_decel":
		return pickFields(fields, "frontDiffDecel")
	case "drive_diff_accel":
		return drivetrainFields(profile.Drivetrain, fields, "frontDiffAccel", "rearDiffAccel")
	case "drive_tire_pressure":
		return drivetrainFields(profile.Drivetrain, fields, "frontTirePressure", "rearTirePressure")
	case "tire_pressure":
		return pickFields(fields, "frontTirePressure", "rearTirePressure")
	case "front_arb":
		return pickFields(fields, "frontArb")
	case "rear_arb":
		return pickFields(fields, "rearArb")
	case "front_rebound":
		return pickFields(fields, "frontRebound")
	case "rear_rebound":
		return pickFields(fields, "rearRebound")
	case "front_camber":
		return pickFields(fields, "frontCamber")
	case "front_and_rear_aero":
		return pickFields(fields, "frontAero", "rearAero")
	case "ride_height":
		return pickFields(fields, "frontRideHeight", "rearRideHeight")
	case "spring_rate":
		return pickFields(fields, "frontSpring", "rearSpring")
	case "bump":
		return pickFields(fields, "frontBump", "rearBump")
	default:
		return nil
	}
}

func actionDelta(action rankedAction, current float64, step float64, unit string) (float64, bool) {
	amount := strings.ToLower(strings.TrimSpace(action.Amount))
	direction := strings.ToLower(strings.TrimSpace(action.Direction))
	magnitude := 0.0
	switch amount {
	case "one small step", "avoid bottoming":
		magnitude = step
	case "slightly more negative":
		magnitude = step
		return -magnitude, true
	case "0.5 psi":
		magnitude = step
	default:
		if strings.Contains(amount, "%") {
			base := firstNumber(amount)
			if base <= 0 {
				return 0, false
			}
			if unit == "%" {
				magnitude = base
			} else {
				magnitude = mathMax(step, abs(current)*base/100)
			}
		} else {
			base := firstNumber(amount)
			if base <= 0 {
				return 0, false
			}
			magnitude = base
		}
	}
	magnitude = mathMax(step, roundToStep(magnitude, step))
	switch direction {
	case "decrease":
		return -magnitude, true
	case "increase":
		return magnitude, true
	case "check":
		if strings.Contains(amount, "negative") {
			return -magnitude, true
		}
		if strings.Contains(amount, "bottom") || strings.Contains(amount, "small step") {
			return magnitude, true
		}
		return 0, false
	default:
		return 0, false
	}
}

func tuneFieldMap(profile *storage.TuneProfile) map[string]tuneFieldRef {
	return map[string]tuneFieldRef{
		"frontTirePressure": {key: "frontTirePressure", labelZH: "前胎压", labelEN: "Front tire pressure", unit: "BAR", step: 0.01, value: profile.FrontTirePressure},
		"rearTirePressure":  {key: "rearTirePressure", labelZH: "后胎压", labelEN: "Rear tire pressure", unit: "BAR", step: 0.01, value: profile.RearTirePressure},
		"finalDrive":        {key: "finalDrive", labelZH: "终传比", labelEN: "Final drive", step: 0.01, value: profile.FinalDrive},
		"gear1":             {key: "gear1", labelZH: "1 挡齿比", labelEN: "1st gear", step: 0.01, value: profile.Gear1},
		"gear2":             {key: "gear2", labelZH: "2 挡齿比", labelEN: "2nd gear", step: 0.01, value: profile.Gear2},
		"gear3":             {key: "gear3", labelZH: "3 挡齿比", labelEN: "3rd gear", step: 0.01, value: profile.Gear3},
		"gear4":             {key: "gear4", labelZH: "4 挡齿比", labelEN: "4th gear", step: 0.01, value: profile.Gear4},
		"gear5":             {key: "gear5", labelZH: "5 挡齿比", labelEN: "5th gear", step: 0.01, value: profile.Gear5},
		"gear6":             {key: "gear6", labelZH: "6 挡齿比", labelEN: "6th gear", step: 0.01, value: profile.Gear6},
		"gear7":             {key: "gear7", labelZH: "7 挡齿比", labelEN: "7th gear", step: 0.01, value: profile.Gear7},
		"gear8":             {key: "gear8", labelZH: "8 挡齿比", labelEN: "8th gear", step: 0.01, value: profile.Gear8},
		"gear9":             {key: "gear9", labelZH: "9 挡齿比", labelEN: "9th gear", step: 0.01, value: profile.Gear9},
		"gear10":            {key: "gear10", labelZH: "10 挡齿比", labelEN: "10th gear", step: 0.01, value: profile.Gear10},
		"frontCamber":       {key: "frontCamber", labelZH: "前轮外倾角", labelEN: "Front camber", step: 0.1, value: profile.FrontCamber},
		"frontArb":          {key: "frontArb", labelZH: "前防倾杆", labelEN: "Front ARB", step: 0.1, value: profile.FrontARB},
		"rearArb":           {key: "rearArb", labelZH: "后防倾杆", labelEN: "Rear ARB", step: 0.1, value: profile.RearARB},
		"frontSpring":       {key: "frontSpring", labelZH: "前弹簧", labelEN: "Front spring", step: 0.1, value: profile.FrontSpring},
		"rearSpring":        {key: "rearSpring", labelZH: "后弹簧", labelEN: "Rear spring", step: 0.1, value: profile.RearSpring},
		"frontRideHeight":   {key: "frontRideHeight", labelZH: "前车高", labelEN: "Front ride height", step: 0.1, value: profile.FrontRideHeight},
		"rearRideHeight":    {key: "rearRideHeight", labelZH: "后车高", labelEN: "Rear ride height", step: 0.1, value: profile.RearRideHeight},
		"frontRebound":      {key: "frontRebound", labelZH: "前回弹阻尼", labelEN: "Front rebound", step: 0.1, value: profile.FrontRebound},
		"rearRebound":       {key: "rearRebound", labelZH: "后回弹阻尼", labelEN: "Rear rebound", step: 0.1, value: profile.RearRebound},
		"frontBump":         {key: "frontBump", labelZH: "前压缩阻尼", labelEN: "Front bump", step: 0.1, value: profile.FrontBump},
		"rearBump":          {key: "rearBump", labelZH: "后压缩阻尼", labelEN: "Rear bump", step: 0.1, value: profile.RearBump},
		"frontAero":         {key: "frontAero", labelZH: "前下压力", labelEN: "Front aero", step: 1, value: profile.FrontAero},
		"rearAero":          {key: "rearAero", labelZH: "后下压力", labelEN: "Rear aero", step: 1, value: profile.RearAero},
		"brakeBalance":      {key: "brakeBalance", labelZH: "刹车平衡", labelEN: "Brake balance", unit: "%", step: 1, value: profile.BrakeBalance},
		"brakePressure":     {key: "brakePressure", labelZH: "刹车压力", labelEN: "Brake pressure", unit: "%", step: 1, value: profile.BrakePressure},
		"frontDiffAccel":    {key: "frontDiffAccel", labelZH: "前差速加速", labelEN: "Front diff accel", unit: "%", step: 1, value: profile.FrontDiffAccel},
		"frontDiffDecel":    {key: "frontDiffDecel", labelZH: "前差速减速", labelEN: "Front diff decel", unit: "%", step: 1, value: profile.FrontDiffDecel},
		"rearDiffAccel":     {key: "rearDiffAccel", labelZH: "后差速加速", labelEN: "Rear diff accel", unit: "%", step: 1, value: profile.RearDiffAccel},
		"rearDiffDecel":     {key: "rearDiffDecel", labelZH: "后差速减速", labelEN: "Rear diff decel", unit: "%", step: 1, value: profile.RearDiffDecel},
	}
}

func pickFields(fields map[string]tuneFieldRef, keys ...string) []tuneFieldRef {
	out := make([]tuneFieldRef, 0, len(keys))
	for _, key := range keys {
		if field, ok := fields[key]; ok {
			out = append(out, field)
		}
	}
	return out
}

func drivetrainFields(drivetrain string, fields map[string]tuneFieldRef, frontKey string, rearKey string) []tuneFieldRef {
	switch strings.ToUpper(strings.TrimSpace(drivetrain)) {
	case "FWD":
		return pickFields(fields, frontKey)
	case "RWD":
		return pickFields(fields, rearKey)
	default:
		return pickFields(fields, frontKey, rearKey)
	}
}

func directionDeltaZH(delta float64, action rankedAction) string {
	if action.Amount == "slightly more negative" {
		return "增加负外倾"
	}
	if delta < 0 {
		return "降低"
	}
	return "增加"
}

func directionDeltaEN(delta float64, action rankedAction) string {
	if action.Amount == "slightly more negative" {
		return "more negative by"
	}
	if delta < 0 {
		return "decrease by"
	}
	return "increase by"
}

func firstNumber(value string) float64 {
	start := -1
	end := -1
	for index, char := range value {
		if (char >= '0' && char <= '9') || char == '.' {
			if start < 0 {
				start = index
			}
			end = index + 1
			continue
		}
		if start >= 0 {
			break
		}
	}
	if start < 0 || end <= start {
		return 0
	}
	parsed, _ := strconv.ParseFloat(value[start:end], 64)
	return parsed
}

func roundToStep(value float64, step float64) float64 {
	if step <= 0 {
		return value
	}
	return math.Round(value/step) * step
}

func formatTuneNumber(value float64, step float64) string {
	decimals := 2
	if step >= 1 {
		decimals = 0
	} else if step >= 0.1 {
		decimals = 1
	}
	return fmt.Sprintf("%.*f", decimals, value)
}

func abs(value float64) float64 {
	if value < 0 {
		return -value
	}
	return value
}

func mathMax(left float64, right float64) float64 {
	if left > right {
		return left
	}
	return right
}

func fallback(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}

func duration(ms int64) string {
	if ms < 1000 {
		return fmt.Sprintf("%d ms", ms)
	}
	return fmt.Sprintf("%.1f s", float64(ms)/1000)
}

func title(value string) string {
	value = strings.ReplaceAll(value, "_", " ")
	return strings.Title(value)
}

func sentence(value string) string {
	if value == "" {
		return ""
	}
	return strings.ToUpper(value[:1]) + value[1:]
}

func eventLabelZH(value string) string {
	return mapLookup(value, map[string]string{
		"launch_wheelspin": "起步打滑", "launch_bog_down": "起步憋转", "short_gear": "挡位过短",
		"long_gear_bog_down": "长齿比憋转", "top_speed_limited_by_gearing": "极速受齿比限制",
		"front_brake_lockup": "前轮抱死", "rear_brake_lockup": "后轮抱死", "corner_entry_understeer": "入弯推头",
		"corner_exit_oversteer": "出弯甩尾", "high_speed_four_wheel_slide": "高速四轮侧滑", "suspension_bottom_out": "悬挂触底",
	})
}

func eventLabelEN(value string) string { return title(value) }

func severityZH(value string) string {
	return mapLookup(value, map[string]string{"low": "低严重度", "medium": "中严重度", "high": "高严重度"})
}

func segmentZH(value string) string {
	return mapLookup(value, map[string]string{"launch": "起步", "acceleration": "加速", "braking": "制动", "corner_entry": "入弯", "corner_exit": "出弯", "high_speed_corner": "高速弯", "suspension": "悬挂"})
}

func segmentEN(value string) string { return title(value) }

func labelZH(value string) string {
	return mapLookup(value, map[string]string{
		"speed_kmh": "速度", "gear": "挡位", "throttle": "油门", "brake": "刹车", "steer_abs": "方向输入",
		"front_slip_ratio": "前轮滑移率", "rear_slip_ratio": "后轮滑移率", "max_slip_ratio": "最大滑移率", "rpm_ratio": "转速比例",
		"front_combined_slip": "前轮综合滑移", "rear_combined_slip": "后轮综合滑移", "yaw_rate_abs": "横摆角速度",
		"slip_delta": "前后滑移差", "corner_operation_state": "过弯操作状态",
		"max_suspension_travel": "最大悬挂行程", "front_suspension": "前悬挂行程", "rear_suspension": "后悬挂行程",
		"pitch_rate_abs": "俯仰角速度", "roll_rate_abs": "侧倾角速度",
	})
}

func labelEN(value string) string { return title(value) }

func categoryZH(value string) string {
	return mapLookup(value, map[string]string{"aero": "空气动力", "alignment": "定位", "brake": "刹车", "damping": "阻尼", "differential": "差速器", "gearing": "齿比", "suspension": "悬挂", "tire": "轮胎"})
}

func itemZH(value string) string {
	return mapLookup(value, map[string]string{
		"brake_balance": "刹车平衡", "brake_pressure": "刹车压力", "bump": "压缩阻尼", "current_gear": "当前挡位齿比",
		"drive_diff_accel": "驱动轮加速差速", "drive_tire_pressure": "驱动轮胎压", "final_drive": "终传比",
		"front_and_rear_aero": "前后下压力", "front_arb": "前防倾杆", "front_camber": "前轮外倾角", "front_rebound": "前回弹阻尼",
		"gear_1": "1 挡齿比", "rear_arb": "后防倾杆", "rear_diff_accel": "后差速加速", "rear_diff_decel": "后差速减速",
		"rear_rebound": "后回弹阻尼", "ride_height": "车身高度", "spring_rate": "弹簧硬度", "tire_pressure": "胎压",
	})
}

func directionZH(value string) string {
	return mapLookup(value, map[string]string{"check": "检查", "decrease": "降低", "increase": "提高"})
}

func amountZH(value string) string {
	return mapLookup(value, map[string]string{"one small step": "一小格", "slightly more negative": "略微增加负外倾", "avoid bottoming": "避免触底", "1%-2% rearward": "向后 1%-2%", "1%-2% forward": "向前 1%-2%"})
}

func reasonZH(value string) string {
	return mapLookup(value, map[string]string{
		"avoid hitting the top of the gear too early":    "避免过早顶到当前挡位转速上限",
		"help the engine stay in the power band":         "帮助发动机保持在动力区间",
		"increase front grip on entry":                   "提高入弯阶段前轮抓地",
		"increase high-speed grip":                       "提高高速抓地",
		"increase launch traction":                       "提高起步牵引力",
		"increase rear grip":                             "提高后轮抓地",
		"increase tire contact patch":                    "增加轮胎接地面积",
		"improve front tire contact in cornering":        "改善持续过弯时的前轮接地状态",
		"improve rear compliance under braking":          "改善制动时后轴贴服性",
		"lengthen all gears if multiple gears are short": "如果多个挡位都偏短，整体拉长齿比",
		"let the front tires load more smoothly":         "让前轮载荷转移更平顺",
		"make threshold braking easier":                  "降低临界刹车控制难度",
		"prevent aero and suspension instability":        "避免下压力和悬挂状态不稳定",
		"reduce bottoming frequency":                     "减少触底频率",
		"reduce driven-wheel slip":                       "降低驱动轮滑移",
		"reduce front lockup tendency":                   "降低前轮抱死倾向",
		"reduce power oversteer":                         "降低动力甩尾倾向",
		"reduce rear lockup tendency":                    "降低后轮抱死倾向",
		"reduce wheel torque during launch":              "降低起步时轮上扭矩",
		"reduce wheel torque on exit":                    "降低出弯时轮上扭矩",
		"restore suspension travel":                      "恢复悬挂可用行程",
		"shorten launch gearing":                         "缩短起步齿比",
		"shorten road acceleration gearing":              "缩短公路加速齿比",
		"stabilize the rear axle while braking":          "稳定制动时的后轴",
		"support compression on impacts":                 "提高冲击压缩阶段支撑",
		"increase top speed headroom":                    "增加极速转速余量",
		"verify aero drag is not limiting top speed":     "确认空阻没有限制极速",
		"verify traction is not limiting exit drive":     "确认牵引力没有限制出弯加速",
	})
}

func causeZH(value string) string {
	return mapLookup(value, map[string]string{
		"launch_wheelspin":             "起步阶段轮上扭矩超过驱动轮抓地，优先检查一挡、终传比、驱动轮胎压和加速差速。",
		"launch_bog_down":              "起步转速被压低且没有明显打滑，齿比可能偏长或发动机未保持在动力区间。",
		"short_gear":                   "当前挡位过早接近红线，可能限制加速延展性。",
		"long_gear_bog_down":           "公路加速或出弯阶段转速偏低且没有明显打滑，当前挡位或终传可能偏长。",
		"top_speed_limited_by_gearing": "高速段大油门时过早接近红线，极速可能被高挡或终传齿比限制。",
		"front_brake_lockup":           "前轴制动负担过高或前轮抓地不足，导致前轮先进入滑移。",
		"rear_brake_lockup":            "后轴制动或减速差速稳定性不足，车尾在重刹时更容易失稳。",
		"corner_entry_understeer":      "入弯前轮滑移明显高于后轮，前轴抓地或载荷转移不足。",
		"corner_exit_oversteer":        "出弯给油时后轮先失抓，优先看后差速加速、当前挡位和后轴机械抓地。",
		"high_speed_four_wheel_slide":  "高速前后轴同时滑移，可能需要提升高速抓地、空力或检查悬挂平台。",
		"suspension_bottom_out":        "悬挂行程接近满行程，可能导致抓地突然丢失或车身平台不稳定。",
	})
}

func causeEN(value string) string {
	return mapLookup(value, map[string]string{
		"launch_wheelspin":             "Driven-wheel torque is exceeding available grip during launch.",
		"launch_bog_down":              "The engine is falling out of the power band without useful wheelspin.",
		"short_gear":                   "The current gear reaches the top of the rev range too early.",
		"long_gear_bog_down":           "The engine is below the power band under road acceleration without meaningful wheelspin.",
		"top_speed_limited_by_gearing": "The car is reaching the top of the rev range too early in high-speed sections.",
		"front_brake_lockup":           "Front brake load or front tire demand is too high under braking.",
		"rear_brake_lockup":            "Rear braking stability is weak, often from brake balance, rear decel diff, or rear compliance.",
		"corner_entry_understeer":      "Front axle slip is much higher than rear axle slip on entry.",
		"corner_exit_oversteer":        "Rear tires lose grip first under throttle on exit.",
		"high_speed_four_wheel_slide":  "Both axles are sliding at high speed, pointing to high-speed grip or platform stability.",
		"suspension_bottom_out":        "Suspension travel is nearly exhausted, which can destabilize the car.",
	})
}

func mapLookup(value string, labels map[string]string) string {
	if label, ok := labels[value]; ok {
		return label
	}
	return title(value)
}
