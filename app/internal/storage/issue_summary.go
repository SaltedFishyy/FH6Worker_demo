package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"

	"fh6worker/internal/telemetry"
)

const (
	issueCompareImproved    = "improved"
	issueCompareWorsened    = "worsened"
	issueCompareUnchanged   = "unchanged"
	issueCompareUnavailable = "unavailable"

	powerBandRpmMin               = 0.55
	powerBandRpmMax               = 0.90
	lowRpmHighLoadThreshold       = 0.25
	highRpmHighLoadThreshold      = 0.20
	tractionLimitedPowerThreshold = 0.25
	axleTractionDiffThreshold     = 0.20
	axleTractionDominanceMargin   = 0.10
	gearGlobalIssueMinRatio       = 0.60
	gearGlobalIssueMinCount       = 3
)

func (s *Store) GetSessionIssueSummary(sessionID int64) (*SessionIssueSummary, error) {
	session, err := s.GetTelemetrySession(sessionID)
	if err != nil {
		return nil, err
	}
	events, err := s.GetSessionEvents(sessionID)
	if err != nil {
		return nil, err
	}
	recentChanges := s.recentChangeFieldsForSession(*session)
	recentChangeDeltas := s.recentChangeDeltasForSession(*session)
	groups := BuildSessionIssueGroups(events, recentChanges)
	summary := &SessionIssueSummary{
		SessionID:          sessionID,
		BaselineStatus:     issueCompareUnavailable,
		RecentChangeFields: recentChanges,
		Groups:             groups,
	}
	baseline, status, err := s.findIssueBaselineSession(*session)
	if err != nil {
		return nil, err
	}
	summary.BaselineStatus = status
	if baseline == nil {
		applyIssueComparisons(summary.Groups, nil)
		applyIssueStrategies(summary.Groups, recentChangeDeltas)
		if err := s.applyPowerAndWholeCarPlan(summary, *session, nil); err != nil {
			return nil, err
		}
		return summary, nil
	}
	baselineEvents, err := s.GetSessionEvents(baseline.ID)
	if err != nil {
		return nil, err
	}
	baselineGroups := BuildSessionIssueGroups(baselineEvents, nil)
	summary.BaselineSession = baseline
	applyIssueComparisons(summary.Groups, baselineGroups)
	applyIssueStrategies(summary.Groups, recentChangeDeltas)
	if err := s.applyPowerAndWholeCarPlan(summary, *session, baseline); err != nil {
		return nil, err
	}
	return summary, nil
}

func (s *Store) applyPowerAndWholeCarPlan(summary *SessionIssueSummary, session TelemetrySession, baseline *TelemetrySession) error {
	samples, err := s.GetSessionTelemetrySamples(summary.SessionID, 10000)
	if err != nil {
		return err
	}
	profile := s.issueTuneProfile(session)
	summary.GearPower = BuildGearPowerDiagnostic(samples, summary.Groups, profile)
	summary.GearPower.Comparisons = s.buildGearPowerComparisons(session, baseline, summary.GearPower)
	summary.WholeCarPlan = BuildWholeCarTuningPlan(summary.Groups, summary.GearPower, profile)
	return nil
}

func BuildSessionIssueGroups(events []telemetry.DetectedEvent, recentChanges []string) []SessionIssueGroup {
	builders := map[string]*issueGroupBuilder{}
	for _, event := range events {
		family := issueFamilyForEvent(event.Type)
		builder := builders[family]
		if builder == nil {
			builder = &issueGroupBuilder{
				group: SessionIssueGroup{
					ID:           family,
					Family:       family,
					Severity:     event.Severity,
					Segment:      event.Segment,
					FirstStartMS: event.StartMS,
					LastEndMS:    event.EndMS,
					Evidence:     map[string]IssueEvidence{},
				},
				evidenceSums: map[string]float64{},
				eventTypes:   map[string]bool{},
				actionKeys:   map[string]bool{},
			}
			builders[family] = builder
		}
		builder.add(event)
	}
	groups := make([]SessionIssueGroup, 0, len(builders))
	for _, builder := range builders {
		group := builder.finish(recentChanges)
		groups = append(groups, group)
	}
	sort.SliceStable(groups, func(i, j int) bool {
		if severityRank(groups[i].Severity) != severityRank(groups[j].Severity) {
			return severityRank(groups[i].Severity) > severityRank(groups[j].Severity)
		}
		if groups[i].EventCount != groups[j].EventCount {
			return groups[i].EventCount > groups[j].EventCount
		}
		if groups[i].TotalDurationMS != groups[j].TotalDurationMS {
			return groups[i].TotalDurationMS > groups[j].TotalDurationMS
		}
		return groups[i].Family < groups[j].Family
	})
	return groups
}

type issueGroupBuilder struct {
	group        SessionIssueGroup
	evidenceSums map[string]float64
	eventTypes   map[string]bool
	actionKeys   map[string]bool
}

func (b *issueGroupBuilder) add(event telemetry.DetectedEvent) {
	if b.group.EventCount == 0 || severityRank(event.Severity) > severityRank(b.group.Severity) {
		b.group.Severity = event.Severity
	}
	if event.StartMS < b.group.FirstStartMS {
		b.group.FirstStartMS = event.StartMS
	}
	if event.EndMS > b.group.LastEndMS {
		b.group.LastEndMS = event.EndMS
	}
	if b.group.Segment == "" {
		b.group.Segment = event.Segment
	}
	b.group.EventCount++
	b.group.TotalDurationMS += event.DurationMS
	b.group.EventIDs = append(b.group.EventIDs, event.ID)
	b.group.Events = append(b.group.Events, event)
	if !b.eventTypes[event.Type] {
		b.eventTypes[event.Type] = true
		b.group.EventTypes = append(b.group.EventTypes, event.Type)
	}
	for key, value := range event.Evidence {
		stat := b.group.Evidence[key]
		if stat.Count == 0 || value < stat.Min {
			stat.Min = value
		}
		if stat.Count == 0 || value > stat.Max {
			stat.Max = value
		}
		stat.Count++
		b.evidenceSums[key] += value
		stat.Avg = b.evidenceSums[key] / float64(stat.Count)
		b.group.Evidence[key] = stat
	}
	for _, action := range event.SuggestedActions {
		key := strings.Join([]string{action.Category, action.Item, action.Direction}, "/")
		if b.actionKeys[key] {
			continue
		}
		b.actionKeys[key] = true
		b.group.PrimaryActions = append(b.group.PrimaryActions, action)
	}
}

func (b *issueGroupBuilder) finish(recentChanges []string) SessionIssueGroup {
	sort.Strings(b.group.EventTypes)
	sort.SliceStable(b.group.PrimaryActions, func(i, j int) bool {
		return b.group.PrimaryActions[i].Priority < b.group.PrimaryActions[j].Priority
	})
	if len(b.group.PrimaryActions) > 4 {
		b.group.PrimaryActions = b.group.PrimaryActions[:4]
	}
	related := relatedRecentChangesForActions(b.group.PrimaryActions, recentChanges)
	b.group.RelatedRecentChanges = related
	b.group.PrioritizeTuning = shouldPrioritizeIssueFamily(b.group.Family, b.group.Severity, b.group.EventCount, len(related) > 0)
	b.group.Comparison = issueCompareUnavailable
	return b.group
}

func applyIssueComparisons(groups []SessionIssueGroup, baseline []SessionIssueGroup) {
	baselineByFamily := map[string]SessionIssueGroup{}
	for _, group := range baseline {
		baselineByFamily[group.Family] = group
	}
	for i := range groups {
		previous, ok := baselineByFamily[groups[i].Family]
		if !ok {
			groups[i].Comparison = issueCompareUnavailable
			continue
		}
		groups[i].BaselineEventCount = previous.EventCount
		groups[i].BaselineTotalDurationMS = previous.TotalDurationMS
		currentScore := issueGroupScore(groups[i])
		previousScore := issueGroupScore(previous)
		switch {
		case currentScore < previousScore*0.85:
			groups[i].Comparison = issueCompareImproved
		case currentScore > previousScore*1.15:
			groups[i].Comparison = issueCompareWorsened
		default:
			groups[i].Comparison = issueCompareUnchanged
		}
	}
}

func issueGroupScore(group SessionIssueGroup) float64 {
	return float64(group.EventCount)*10 + float64(group.TotalDurationMS)/1000*4 + float64(severityRank(group.Severity))*8
}

func applyIssueStrategies(groups []SessionIssueGroup, recentDeltas map[string]float64) {
	for i := range groups {
		groups[i].AdjustmentStrategy = adjustmentStrategyForGroup(groups[i])
		groups[i].FeedbackDirective = feedbackDirectiveForGroup(groups[i])
		if groups[i].FeedbackDirective == "rollback_related_changes" {
			if actions := rollbackActionsForGroup(groups[i], recentDeltas); len(actions) > 0 {
				groups[i].PrimaryActions = actions
			}
			continue
		}
		groups[i].PrimaryActions = scaledStrategyActions(groups[i])
	}
}

type powerBandWindow struct {
	startRPM float64
	endRPM   float64
	redline  float64
	source   string
}

func BuildGearPowerDiagnostic(samples []telemetry.NormalizedTelemetry, groups []SessionIssueGroup, profile *TuneProfile) GearPowerDiagnostic {
	diag := GearPowerDiagnostic{
		Status:   "insufficient_data",
		Summary:  "not_enough_samples",
		Evidence: map[string]float64{},
	}
	applyProfilePowerToDiagnostic(&diag, profile)
	if len(samples) == 0 {
		return diag
	}
	orderedSamples := sortedTelemetrySamples(samples)
	diag.Evidence["power_band_total_samples"] = float64(len(orderedSamples))
	powerBand := resolvePowerBandWindow(orderedSamples, profile)
	diag.PowerBandStartRPM = powerBand.startRPM
	diag.PowerBandEndRPM = powerBand.endRPM
	diag.RedlineRPM = math.Max(diag.RedlineRPM, powerBand.redline)
	diag.PowerBandSource = powerBand.source
	drivetrain := gearPowerDrivetrain(profile, orderedSamples)
	builders := map[int]*gearPowerBuilder{}
	for _, sample := range orderedSamples {
		if sample.Gear < 1 || sample.Gear > 10 || !profileGearUnlocked(profile, sample.Gear) {
			continue
		}
		builder := builders[sample.Gear]
		if builder == nil {
			builder = &gearPowerBuilder{gear: sample.Gear}
			builders[sample.Gear] = builder
		}
		builder.add(sample, powerBand)
	}
	if len(builders) == 0 {
		if profile != nil {
			diag.Summary = "no_unlocked_gear_samples"
		}
		return diag
	}
	diag.Status = "ok"
	tractionLimitedGears := make([]int, 0)
	longGearBands := make([]GearPowerBand, 0)
	shortGearBands := make([]GearPowerBand, 0)
	diffEvidence := gearPowerDiffEvidence{}
	totalSamples := 0
	totalHighLoadSamples := 0
	usableGearCount := 0
	for gear := 1; gear <= 10; gear++ {
		if !profileGearUnlocked(profile, gear) {
			continue
		}
		builder := builders[gear]
		if builder == nil {
			continue
		}
		band := builder.finish(powerBand)
		applyGearShiftEstimate(&band, profile)
		totalSamples += band.SampleCount
		totalHighLoadSamples += band.HighLoadSampleCount
		if band.HighLoadSampleCount >= 4 {
			usableGearCount++
		}
		if band.Finding == "too_long" {
			longGearBands = append(longGearBands, band)
		}
		if band.Finding == "too_short" {
			shortGearBands = append(shortGearBands, band)
		}
		if band.Finding == "traction_limited" {
			tractionLimitedGears = append(tractionLimitedGears, gear)
			diffEvidence.addBand(band)
			if gear == 1 {
				diag.RecommendedActions = append(diag.RecommendedActions, actionTemplate(len(diag.RecommendedActions), "gearing", "gear_1", "decrease", "0.05", "reduce wheel torque before changing chassis balance"))
			}
		}
		diag.Gears = append(diag.Gears, band)
		diag.Evidence[fmt.Sprintf("gear_%d_rpm_avg", gear)] = band.RpmRatioAvg
		diag.Evidence[fmt.Sprintf("gear_%d_rpm_abs_avg", gear)] = band.RpmAvg
		diag.Evidence[fmt.Sprintf("gear_%d_speed_max", gear)] = band.SpeedMaxKmh
		diag.Evidence[fmt.Sprintf("gear_%d_accel_avg_mps2", gear)] = band.AccelAvgMps2
		diag.Evidence[fmt.Sprintf("gear_%d_in_power_band_pct", gear)] = band.InPowerBandPercent
		diag.Evidence[fmt.Sprintf("gear_%d_traction_limited_pct", gear)] = band.TractionLimitedPercent
		diag.Evidence[fmt.Sprintf("gear_%d_front_slip_avg", gear)] = band.FrontSlipAvg
		diag.Evidence[fmt.Sprintf("gear_%d_rear_slip_avg", gear)] = band.RearSlipAvg
		diag.Evidence[fmt.Sprintf("gear_%d_front_traction_limited_pct", gear)] = band.FrontTractionLimitedPct
		diag.Evidence[fmt.Sprintf("gear_%d_rear_traction_limited_pct", gear)] = band.RearTractionLimitedPct
		diag.LowRpmHighLoadPercent = math.Max(diag.LowRpmHighLoadPercent, band.LowRpmHighLoadPercent)
		diag.HighRpmHighLoadPercent = math.Max(diag.HighRpmHighLoadPercent, band.HighRpmHighLoadPercent)
		diag.TractionLimitedPercent = math.Max(diag.TractionLimitedPercent, band.TractionLimitedPercent)
	}
	diag.UsableGearCount = usableGearCount
	diffEvidence.writeTo(diag.Evidence)
	diag.LaunchFinding = launchFindingForGroups(groups)
	diag.TopSpeedFinding = topSpeedFinding(diag.Gears)
	applyGearPowerStrategy(&diag, longGearBands, shortGearBands, tractionLimitedGears)
	if shouldRecommendDriveDiffAccel(tractionLimitedGears) {
		for _, recommendation := range accelDiffRecommendationsForDrivetrain(drivetrain, diffEvidence) {
			diag.RecommendedActions = append(diag.RecommendedActions, actionTemplate(len(diag.RecommendedActions), "differential", recommendation.item, "decrease", recommendation.amount, recommendation.reason))
		}
	}
	diag.RecommendedActions = dedupeGearPowerActions(diag.RecommendedActions, diag.Evidence)
	if totalHighLoadSamples < 4 || usableGearCount == 0 {
		diag.Summary = "not_enough_high_load"
	} else if len(diag.RecommendedActions) == 0 {
		diag.Summary = "gearing_window_ok"
	} else {
		diag.Summary = "gearing_adjustment_needed"
	}
	if diag.TractionLimitedPercent >= tractionLimitedPowerThreshold {
		diag.Summary = "traction_limited_power"
	}
	diag.Confidence = gearPowerConfidence(diag.PowerBandSource, totalHighLoadSamples, usableGearCount)
	diag.Evidence["power_band_target_min"] = powerBandRpmMin
	diag.Evidence["power_band_target_max"] = powerBandRpmMax
	diag.Evidence["power_band_start_rpm"] = diag.PowerBandStartRPM
	diag.Evidence["power_band_end_rpm"] = diag.PowerBandEndRPM
	diag.Evidence["redline_rpm"] = diag.RedlineRPM
	diag.Evidence["power_band_total_samples"] = float64(totalSamples)
	diag.Evidence["power_band_high_load_samples"] = float64(totalHighLoadSamples)
	diag.Evidence["power_band_usable_gears"] = float64(usableGearCount)
	diag.Evidence["gear_strategy_issue_count"] = float64(diag.GlobalGearIssueCount)
	diag.Evidence["gear_strategy_issue_ratio"] = diag.GlobalGearIssueRatio
	diag.Evidence["traction_limited_percent"] = diag.TractionLimitedPercent
	diag.Evidence["low_rpm_high_load_percent"] = diag.LowRpmHighLoadPercent
	diag.Evidence["high_rpm_high_load_percent"] = diag.HighRpmHighLoadPercent
	return diag
}

type gearPowerStrategyCandidate struct {
	mode      string
	direction string
	bands     []GearPowerBand
	ratio     float64
}

func applyGearPowerStrategy(diag *GearPowerDiagnostic, longBands []GearPowerBand, shortBands []GearPowerBand, tractionLimitedGears []int) {
	if diag == nil {
		return
	}
	if diag.TopSpeedFinding == "top_speed_limited_by_gearing" {
		diag.StrategyMode = "top_speed_limited"
		diag.GlobalGearIssueCount = 1
		diag.GlobalGearIssueRatio = safeDiv(1, float64(diag.UsableGearCount))
		diag.RecommendedActions = append(diag.RecommendedActions, actionTemplate(len(diag.RecommendedActions), "gearing", "final_drive", "decrease", "0.08", "increase top speed headroom"))
		return
	}
	if candidate, ok := dominantGlobalGearCandidate(longBands, shortBands, diag.UsableGearCount); ok {
		diag.StrategyMode = candidate.mode
		diag.GlobalGearIssueCount = len(candidate.bands)
		diag.GlobalGearIssueRatio = candidate.ratio
		diag.RecommendedActions = append(diag.RecommendedActions, actionTemplate(len(diag.RecommendedActions), "gearing", "final_drive", candidate.direction, finalDriveAdjustmentAmount(candidate.bands, candidate.ratio), gearStrategyReason(candidate.mode)))
		return
	}
	if len(tractionLimitedGears) > 0 {
		diag.StrategyMode = "traction_limited_low_gears"
		diag.GlobalGearIssueCount = len(tractionLimitedGears)
		diag.GlobalGearIssueRatio = safeDiv(float64(len(tractionLimitedGears)), float64(diag.UsableGearCount))
		return
	}
	if diag.TopSpeedFinding == "top_speed_bog_down" {
		diag.StrategyMode = "top_speed_limited"
		diag.GlobalGearIssueCount = 1
		diag.GlobalGearIssueRatio = safeDiv(1, float64(diag.UsableGearCount))
		diag.RecommendedActions = append(diag.RecommendedActions, actionTemplate(len(diag.RecommendedActions), "gearing", "final_drive", "increase", "0.06", "help the engine stay in the power band"))
		return
	}
	appendSingleGearStrategyActions(diag, longBands, shortBands)
}

func dominantGlobalGearCandidate(longBands []GearPowerBand, shortBands []GearPowerBand, usableGearCount int) (gearPowerStrategyCandidate, bool) {
	if usableGearCount < gearGlobalIssueMinCount {
		return gearPowerStrategyCandidate{}, false
	}
	longRatio := safeDiv(float64(len(longBands)), float64(usableGearCount))
	shortRatio := safeDiv(float64(len(shortBands)), float64(usableGearCount))
	longGlobal := len(longBands) >= gearGlobalIssueMinCount && longRatio >= gearGlobalIssueMinRatio
	shortGlobal := len(shortBands) >= gearGlobalIssueMinCount && shortRatio >= gearGlobalIssueMinRatio
	if !longGlobal && !shortGlobal {
		return gearPowerStrategyCandidate{}, false
	}
	if longGlobal && (!shortGlobal || gearIssueDominates(longBands, longRatio, shortBands, shortRatio)) {
		return gearPowerStrategyCandidate{mode: "global_too_long", direction: "increase", bands: longBands, ratio: longRatio}, true
	}
	if shortGlobal {
		return gearPowerStrategyCandidate{mode: "global_too_short", direction: "decrease", bands: shortBands, ratio: shortRatio}, true
	}
	return gearPowerStrategyCandidate{}, false
}

func gearIssueDominates(left []GearPowerBand, leftRatio float64, right []GearPowerBand, rightRatio float64) bool {
	if len(left) != len(right) {
		return len(left) > len(right)
	}
	if math.Abs(leftRatio-rightRatio) > 0.001 {
		return leftRatio > rightRatio
	}
	return gearIssueSeverity(left) >= gearIssueSeverity(right)
}

func appendSingleGearStrategyActions(diag *GearPowerDiagnostic, longBands []GearPowerBand, shortBands []GearPowerBand) {
	if len(longBands) == 0 && len(shortBands) == 0 {
		diag.StrategyMode = "gearing_window_ok"
		return
	}
	diag.GlobalGearIssueCount = len(longBands) + len(shortBands)
	diag.GlobalGearIssueRatio = safeDiv(float64(diag.GlobalGearIssueCount), float64(diag.UsableGearCount))
	switch {
	case len(longBands) > 0 && len(shortBands) == 0:
		diag.StrategyMode = "single_gear_too_long"
	case len(shortBands) > 0 && len(longBands) == 0:
		diag.StrategyMode = "single_gear_too_short"
	default:
		diag.StrategyMode = "single_gear_mixed"
	}
	for _, band := range longBands {
		diag.RecommendedActions = append(diag.RecommendedActions, actionTemplate(len(diag.RecommendedActions), "gearing", gearActionItem(band.Gear), "increase", gearAdjustmentAmount(band), "help the engine stay in the power band"))
	}
	for _, band := range shortBands {
		diag.RecommendedActions = append(diag.RecommendedActions, actionTemplate(len(diag.RecommendedActions), "gearing", gearActionItem(band.Gear), "decrease", gearAdjustmentAmount(band), "avoid hitting the top of the gear too early"))
	}
}

func finalDriveAdjustmentAmount(bands []GearPowerBand, ratio float64) string {
	severity := gearIssueSeverity(bands)
	if severity >= 0.50 || ratio >= 0.80 {
		return "0.10"
	}
	if severity >= 0.35 || ratio >= 0.70 {
		return "0.06"
	}
	return "0.03"
}

func gearIssueSeverity(bands []GearPowerBand) float64 {
	severity := 0.0
	for _, band := range bands {
		severity = math.Max(severity, math.Max(band.BelowPowerBandPercent, band.AbovePowerBandPercent))
	}
	return severity
}

func gearStrategyReason(mode string) string {
	switch mode {
	case "global_too_long":
		return "most unlocked gears are below the power band, adjust final drive before individual gears"
	case "global_too_short":
		return "most unlocked gears are above the power band, adjust final drive before individual gears"
	default:
		return "adjust final drive before individual gears"
	}
}

func (s *Store) buildGearPowerComparisons(session TelemetrySession, baseline *TelemetrySession, current GearPowerDiagnostic) []GearPowerComparison {
	comparisons := []GearPowerComparison{
		{
			Type:   "session_telemetry",
			Status: "missing_baseline",
		},
		s.gearSettingComparisonForSession(session),
	}
	if baseline != nil {
		baselineSamples, err := s.GetSessionTelemetrySamples(baseline.ID, 10000)
		if err == nil && len(baselineSamples) > 0 {
			baselineDiag := BuildGearPowerDiagnostic(baselineSamples, nil, s.issueTuneProfile(*baseline))
			comparisons[0] = buildGearTelemetryComparison(baseline.ID, baselineDiag, current)
		}
	}
	return comparisons
}

func buildGearTelemetryComparison(baselineSessionID int64, baseline GearPowerDiagnostic, current GearPowerDiagnostic) GearPowerComparison {
	comparison := GearPowerComparison{
		Type:              "session_telemetry",
		Status:            "ready",
		BaselineSessionID: baselineSessionID,
	}
	beforeByGear := map[int]GearPowerBand{}
	afterByGear := map[int]GearPowerBand{}
	for _, band := range baseline.Gears {
		beforeByGear[band.Gear] = band
	}
	for _, band := range current.Gears {
		afterByGear[band.Gear] = band
	}
	gears := map[int]bool{}
	for gear := range beforeByGear {
		gears[gear] = true
	}
	for gear := range afterByGear {
		gears[gear] = true
	}
	for gear := 1; gear <= 10; gear++ {
		if !gears[gear] {
			continue
		}
		before := beforeByGear[gear]
		after := afterByGear[gear]
		row := GearPowerComparisonRow{
			Item:                   gearActionItem(gear),
			Gear:                   gear,
			BeforeSpeedMaxKmh:      before.SpeedMaxKmh,
			AfterSpeedMaxKmh:       after.SpeedMaxKmh,
			SpeedMaxDeltaKmh:       after.SpeedMaxKmh - before.SpeedMaxKmh,
			BeforeInPowerBandPct:   before.InPowerBandPercent,
			AfterInPowerBandPct:    after.InPowerBandPercent,
			InPowerBandDeltaPct:    after.InPowerBandPercent - before.InPowerBandPercent,
			BeforeTractionLimitPct: before.TractionLimitedPercent,
			AfterTractionLimitPct:  after.TractionLimitedPercent,
			TractionLimitDeltaPct:  after.TractionLimitedPercent - before.TractionLimitedPercent,
			BeforeFinding:          before.Finding,
			AfterFinding:           after.Finding,
		}
		comparison.Rows = append(comparison.Rows, row)
	}
	if len(comparison.Rows) == 0 {
		comparison.Status = "no_matching_gears"
	}
	return comparison
}

func (s *Store) gearSettingComparisonForSession(session TelemetrySession) GearPowerComparison {
	comparison := GearPowerComparison{
		Type:   "tune_settings",
		Status: "no_changed_gears",
	}
	if session.TuneProfileID == nil {
		comparison.Status = "profile_unbound"
		return comparison
	}
	snapshot := s.recentSnapshotForSession(session)
	if snapshot == nil {
		return comparison
	}
	for _, field := range snapshot.ChangedFields {
		if !isGearTuneField(field) {
			continue
		}
		before, okBefore := tuneProfileNumericField(snapshot.Before, field)
		after, okAfter := tuneProfileNumericField(snapshot.After, field)
		if !okBefore || !okAfter || before == after {
			continue
		}
		comparison.Rows = append(comparison.Rows, GearPowerComparisonRow{
			Item:        gearTuneItemForField(field),
			Gear:        gearNumberForTuneField(field),
			BeforeValue: before,
			AfterValue:  after,
			DeltaValue:  after - before,
		})
	}
	if len(comparison.Rows) > 0 {
		comparison.Status = "ready"
	}
	return comparison
}

func gearPowerDrivetrain(profile *TuneProfile, samples []telemetry.NormalizedTelemetry) string {
	if profile != nil && strings.TrimSpace(profile.Drivetrain) != "" {
		return strings.ToUpper(strings.TrimSpace(profile.Drivetrain))
	}
	for _, sample := range samples {
		if strings.TrimSpace(sample.Drivetrain) != "" {
			return strings.ToUpper(strings.TrimSpace(sample.Drivetrain))
		}
	}
	return ""
}

type accelDiffRecommendation struct {
	item   string
	amount string
	reason string
}

type gearPowerDiffEvidence struct {
	count            int
	frontSlipMax     float64
	rearSlipMax      float64
	frontTractionMax float64
	rearTractionMax  float64
}

func (e *gearPowerDiffEvidence) addBand(band GearPowerBand) {
	if e == nil {
		return
	}
	e.count++
	e.frontSlipMax = math.Max(e.frontSlipMax, band.FrontSlipAvg)
	e.rearSlipMax = math.Max(e.rearSlipMax, band.RearSlipAvg)
	e.frontTractionMax = math.Max(e.frontTractionMax, band.FrontTractionLimitedPct)
	e.rearTractionMax = math.Max(e.rearTractionMax, band.RearTractionLimitedPct)
}

func (e gearPowerDiffEvidence) writeTo(target map[string]float64) {
	if target == nil {
		return
	}
	target["diff_evidence_gears"] = float64(e.count)
	target["front_driven_slip_avg"] = e.frontSlipMax
	target["rear_driven_slip_avg"] = e.rearSlipMax
	target["front_traction_limited_percent"] = e.frontTractionMax
	target["rear_traction_limited_percent"] = e.rearTractionMax
}

func accelDiffRecommendationsForDrivetrain(drivetrain string, evidence gearPowerDiffEvidence) []accelDiffRecommendation {
	switch strings.ToUpper(strings.TrimSpace(drivetrain)) {
	case "FWD":
		if evidence.frontTractionMax < axleTractionDiffThreshold && evidence.frontSlipMax < 1.05 {
			return nil
		}
		return []accelDiffRecommendation{{item: "front_diff_accel", amount: "2", reason: "front-driven traction is limiting power delivery"}}
	case "RWD":
		if evidence.rearTractionMax < axleTractionDiffThreshold && evidence.rearSlipMax < 1.05 {
			return nil
		}
		return []accelDiffRecommendation{{item: "rear_diff_accel", amount: "2", reason: "rear-driven traction is limiting power delivery"}}
	case "AWD":
		return awdAccelDiffRecommendations(evidence)
	default:
		return nil
	}
}

func awdAccelDiffRecommendations(evidence gearPowerDiffEvidence) []accelDiffRecommendation {
	frontActive := evidence.frontTractionMax >= axleTractionDiffThreshold || evidence.frontSlipMax >= 1.05
	rearActive := evidence.rearTractionMax >= axleTractionDiffThreshold || evidence.rearSlipMax >= 1.05
	if !frontActive && !rearActive {
		return nil
	}
	if rearActive && (!frontActive || evidence.rearTractionMax-evidence.frontTractionMax >= axleTractionDominanceMargin || evidence.rearSlipMax-evidence.frontSlipMax >= 0.15) {
		return []accelDiffRecommendation{{item: "rear_diff_accel", amount: "2", reason: "rear axle is the primary traction limit under throttle"}}
	}
	if frontActive && (!rearActive || evidence.frontTractionMax-evidence.rearTractionMax >= axleTractionDominanceMargin || evidence.frontSlipMax-evidence.rearSlipMax >= 0.15) {
		return []accelDiffRecommendation{{item: "front_diff_accel", amount: "2", reason: "front axle is the primary traction limit under throttle"}}
	}
	if evidence.rearSlipMax >= evidence.frontSlipMax {
		return []accelDiffRecommendation{
			{item: "rear_diff_accel", amount: "2", reason: "rear axle is the stronger traction limit under throttle"},
			{item: "front_diff_accel", amount: "1", reason: "front axle also shows traction limit, use a smaller secondary step"},
		}
	}
	return []accelDiffRecommendation{
		{item: "front_diff_accel", amount: "2", reason: "front axle is the stronger traction limit under throttle"},
		{item: "rear_diff_accel", amount: "1", reason: "rear axle also shows traction limit, use a smaller secondary step"},
	}
}

func sortedTelemetrySamples(samples []telemetry.NormalizedTelemetry) []telemetry.NormalizedTelemetry {
	out := append([]telemetry.NormalizedTelemetry(nil), samples...)
	sort.SliceStable(out, func(i, j int) bool {
		return out[i].TimeMS < out[j].TimeMS
	})
	return out
}

func resolvePowerBandWindow(samples []telemetry.NormalizedTelemetry, profile *TuneProfile) powerBandWindow {
	if profile != nil && profile.PeakTorqueRPM != nil && profile.PeakPowerRPM != nil && profile.RedlineRPM != nil &&
		*profile.PeakTorqueRPM > 0 && *profile.PeakPowerRPM > 0 && *profile.RedlineRPM > 0 {
		low := math.Min(*profile.PeakTorqueRPM, *profile.PeakPowerRPM)
		high := math.Max(*profile.PeakTorqueRPM, *profile.PeakPowerRPM)
		end := math.Min(*profile.RedlineRPM, high*1.03)
		if end <= low {
			end = math.Min(*profile.RedlineRPM, low*1.12)
		}
		return powerBandWindow{startRPM: low, endRPM: end, redline: *profile.RedlineRPM, source: "profile_power_band"}
	}
	maxRPM := 0.0
	for _, sample := range samples {
		if sample.EngineMaxRpm > maxRPM {
			maxRPM = sample.EngineMaxRpm
		}
	}
	if maxRPM > 1000 {
		return powerBandWindow{startRPM: maxRPM * powerBandRpmMin, endRPM: maxRPM * powerBandRpmMax, redline: maxRPM, source: "telemetry_engine_max_rpm"}
	}
	return powerBandWindow{source: "rpm_ratio_fallback"}
}

func gearPowerConfidence(source string, highLoadSamples int, usableGears int) string {
	if source == "profile_power_band" && highLoadSamples >= 20 && usableGears >= 2 {
		return "high"
	}
	if highLoadSamples >= 10 && usableGears >= 1 {
		return "medium"
	}
	return "low"
}

func gearAdjustmentAmount(band GearPowerBand) string {
	severity := math.Max(band.BelowPowerBandPercent, band.AbovePowerBandPercent)
	if severity >= 0.50 {
		return "0.08"
	}
	if severity >= 0.35 {
		return "0.05"
	}
	return "0.03"
}

func shouldRecommendDriveDiffAccel(gears []int) bool {
	for _, gear := range gears {
		if gear > 1 {
			return true
		}
	}
	return false
}

func dedupeGearPowerActions(actions []telemetry.SuggestedAction, evidence map[string]float64) []telemetry.SuggestedAction {
	if len(actions) < 2 {
		return actions
	}
	out := make([]telemetry.SuggestedAction, 0, len(actions))
	seen := map[string]struct{}{}
	for _, action := range actions {
		key := actionConflictKey(action, evidence)
		if key == "" {
			key = action.Category + "/" + action.Item
		}
		key += "/" + action.Direction
		if _, exists := seen[key]; exists {
			continue
		}
		action.Priority = len(out)
		out = append(out, action)
		seen[key] = struct{}{}
	}
	return out
}

type gearPowerBuilder struct {
	gear                 int
	count                int
	highLoadCount        int
	lowRpmHighLoadCount  int
	highRpmHighLoadCount int
	inPowerBandCount     int
	tractionLimitedCount int
	speedMin             float64
	speedMax             float64
	speedSum             float64
	rpmAbsMin            float64
	rpmAbsMax            float64
	rpmAbsSum            float64
	rpmRatioMin          float64
	rpmRatioMax          float64
	rpmRatioSum          float64
	inBandRpmMin         float64
	inBandRpmMax         float64
	inBandRatioMin       float64
	inBandRatioMax       float64
	throttleSum          float64
	accelCount           int
	accelSum             float64
	accelMax             float64
	frontSlipHighLoadSum float64
	rearSlipHighLoadSum  float64
	frontTractionCount   int
	rearTractionCount    int
	lastSample           *telemetry.NormalizedTelemetry
}

func (b *gearPowerBuilder) add(sample telemetry.NormalizedTelemetry, powerBand powerBandWindow) {
	if b.count == 0 || sample.SpeedKmh < b.speedMin {
		b.speedMin = sample.SpeedKmh
	}
	if b.count == 0 || sample.SpeedKmh > b.speedMax {
		b.speedMax = sample.SpeedKmh
	}
	if b.count == 0 || sample.Rpm < b.rpmAbsMin {
		b.rpmAbsMin = sample.Rpm
	}
	if b.count == 0 || sample.Rpm > b.rpmAbsMax {
		b.rpmAbsMax = sample.Rpm
	}
	if b.count == 0 || sample.RpmRatio < b.rpmRatioMin {
		b.rpmRatioMin = sample.RpmRatio
	}
	if b.count == 0 || sample.RpmRatio > b.rpmRatioMax {
		b.rpmRatioMax = sample.RpmRatio
	}
	b.count++
	b.speedSum += sample.SpeedKmh
	b.rpmAbsSum += sample.Rpm
	b.rpmRatioSum += sample.RpmRatio
	b.throttleSum += sample.Throttle01
	if b.lastSample != nil && sample.Gear == b.lastSample.Gear {
		dt := float64(sample.TimeMS-b.lastSample.TimeMS) / 1000
		if dt > 0 && dt <= 2 && sample.Throttle01 >= 0.65 && sample.Brake01 < 0.2 {
			accel := ((sample.SpeedKmh - b.lastSample.SpeedKmh) / 3.6) / dt
			b.accelCount++
			b.accelSum += accel
			if b.accelCount == 1 || accel > b.accelMax {
				b.accelMax = accel
			}
		}
	}
	sampleCopy := sample
	b.lastSample = &sampleCopy
	if sample.Throttle01 >= 0.65 {
		b.highLoadCount++
		frontSlip, rearSlip := sampleDrivenAxleSlip(sample)
		b.frontSlipHighLoadSum += frontSlip
		b.rearSlipHighLoadSum += rearSlip
		if frontSlip >= 1.05 {
			b.frontTractionCount++
		}
		if rearSlip >= 1.05 {
			b.rearTractionCount++
		}
		tractionLimited := sampleTractionLimited(sample)
		if tractionLimited {
			b.tractionLimitedCount++
		}
		below := sampleBelowPowerBand(sample, powerBand)
		above := sampleAbovePowerBand(sample, powerBand)
		if !below && !above {
			b.inPowerBandCount++
			b.addInPowerBandRange(sample)
		}
		if below && sample.SpeedKmh > 25 && !tractionLimited {
			b.lowRpmHighLoadCount++
		}
		if above {
			b.highRpmHighLoadCount++
		}
	}
}

func (b *gearPowerBuilder) addInPowerBandRange(sample telemetry.NormalizedTelemetry) {
	if sample.Rpm > 0 {
		if b.inPowerBandCount == 1 || b.inBandRpmMin == 0 || sample.Rpm < b.inBandRpmMin {
			b.inBandRpmMin = sample.Rpm
		}
		if sample.Rpm > b.inBandRpmMax {
			b.inBandRpmMax = sample.Rpm
		}
	}
	if sample.RpmRatio > 0 {
		if b.inPowerBandCount == 1 || b.inBandRatioMin == 0 || sample.RpmRatio < b.inBandRatioMin {
			b.inBandRatioMin = sample.RpmRatio
		}
		if sample.RpmRatio > b.inBandRatioMax {
			b.inBandRatioMax = sample.RpmRatio
		}
	}
}

func (b *gearPowerBuilder) finish(powerBand powerBandWindow) GearPowerBand {
	band := GearPowerBand{
		Gear:                b.gear,
		SampleCount:         b.count,
		HighLoadSampleCount: b.highLoadCount,
		SpeedMinKmh:         b.speedMin,
		SpeedMaxKmh:         b.speedMax,
		SpeedAvgKmh:         safeDiv(b.speedSum, float64(b.count)),
		RpmMin:              b.rpmAbsMin,
		RpmMax:              b.rpmAbsMax,
		RpmAvg:              safeDiv(b.rpmAbsSum, float64(b.count)),
		RpmRatioMin:         b.rpmRatioMin,
		RpmRatioMax:         b.rpmRatioMax,
		RpmRatioAvg:         safeDiv(b.rpmRatioSum, float64(b.count)),
		InPowerBandRpmMin:   b.inBandRpmMin,
		InPowerBandRpmMax:   b.inBandRpmMax,
		InPowerBandRatioMin: b.inBandRatioMin,
		InPowerBandRatioMax: b.inBandRatioMax,
		ThrottleAvg:         safeDiv(b.throttleSum, float64(b.count)),
		AccelAvgMps2:        safeDiv(b.accelSum, float64(b.accelCount)),
		AccelMaxMps2:        b.accelMax,
		Finding:             "ok",
	}
	if band.RpmAvg > 0 {
		band.SpeedPer1000RpmKmh = band.SpeedAvgKmh / (band.RpmAvg / 1000)
	}
	if b.highLoadCount > 0 {
		band.LowRpmHighLoadPercent = float64(b.lowRpmHighLoadCount) / float64(b.highLoadCount)
		band.HighRpmHighLoadPercent = float64(b.highRpmHighLoadCount) / float64(b.highLoadCount)
		band.BelowPowerBandPercent = band.LowRpmHighLoadPercent
		band.AbovePowerBandPercent = band.HighRpmHighLoadPercent
		band.InPowerBandPercent = float64(b.inPowerBandCount) / float64(b.highLoadCount)
		band.TractionLimitedPercent = float64(b.tractionLimitedCount) / float64(b.highLoadCount)
		band.FrontSlipAvg = b.frontSlipHighLoadSum / float64(b.highLoadCount)
		band.RearSlipAvg = b.rearSlipHighLoadSum / float64(b.highLoadCount)
		band.FrontTractionLimitedPct = float64(b.frontTractionCount) / float64(b.highLoadCount)
		band.RearTractionLimitedPct = float64(b.rearTractionCount) / float64(b.highLoadCount)
	}
	if b.count >= 8 && b.highLoadCount >= 4 {
		switch {
		case band.TractionLimitedPercent >= tractionLimitedPowerThreshold:
			band.Finding = "traction_limited"
		case band.BelowPowerBandPercent >= lowRpmHighLoadThreshold || (band.BelowPowerBandPercent >= 0.15 && band.AccelAvgMps2 < 0.4):
			band.Finding = "too_long"
		case band.AbovePowerBandPercent >= highRpmHighLoadThreshold || rpmNearRedline(band, powerBand):
			band.Finding = "too_short"
		}
	}
	return band
}

func sampleBelowPowerBand(sample telemetry.NormalizedTelemetry, powerBand powerBandWindow) bool {
	if powerBand.startRPM > 0 && sample.Rpm > 0 {
		return sample.Rpm < powerBand.startRPM
	}
	return sample.RpmRatio < powerBandRpmMin
}

func sampleAbovePowerBand(sample telemetry.NormalizedTelemetry, powerBand powerBandWindow) bool {
	if powerBand.endRPM > 0 && sample.Rpm > 0 {
		return sample.Rpm > powerBand.endRPM
	}
	return sample.RpmRatio > powerBandRpmMax
}

func rpmNearRedline(band GearPowerBand, powerBand powerBandWindow) bool {
	if powerBand.redline > 0 && band.RpmMax > 0 {
		return band.RpmMax >= powerBand.redline*0.97
	}
	return band.RpmRatioMax > 0.97
}

func applyGearShiftEstimate(band *GearPowerBand, profile *TuneProfile) {
	if band == nil || profile == nil {
		return
	}
	current := tuneProfileFloatPointer(*profile, fmt.Sprintf("gear%d", band.Gear))
	next := tuneProfileFloatPointer(*profile, fmt.Sprintf("gear%d", band.Gear+1))
	if current == nil || next == nil || *current <= 0 || *next <= 0 {
		return
	}
	rpmBefore := band.RpmMax
	if rpmBefore <= 0 {
		return
	}
	after := rpmBefore * (*next / *current)
	if after <= 0 {
		return
	}
	band.ShiftAfterRPM = after
	band.ShiftDropRPM = math.Max(0, rpmBefore-after)
}

func sampleTractionLimited(sample telemetry.NormalizedTelemetry) bool {
	front, rear := sampleDrivenAxleSlip(sample)
	switch strings.ToUpper(strings.TrimSpace(sample.Drivetrain)) {
	case "FWD":
		return front >= 1.05
	case "RWD":
		return rear >= 1.05
	default:
		return math.Max(front, rear) >= 1.15
	}
}

func sampleDrivenAxleSlip(sample telemetry.NormalizedTelemetry) (float64, float64) {
	front := math.Max(math.Abs(sample.FrontCombinedSlipAvg), math.Abs(sample.FrontSlipRatioAvg))
	rear := math.Max(math.Abs(sample.RearCombinedSlipAvg), math.Abs(sample.RearSlipRatioAvg))
	return front, rear
}

func applyProfilePowerToDiagnostic(diag *GearPowerDiagnostic, profile *TuneProfile) {
	if diag == nil || profile == nil {
		return
	}
	if profile.PowerKW != nil {
		diag.PowerKW = *profile.PowerKW
		diag.Evidence["profile_power_kw"] = *profile.PowerKW
	}
	if profile.TorqueNM != nil {
		diag.TorqueNM = *profile.TorqueNM
		diag.Evidence["profile_torque_nm"] = *profile.TorqueNM
	}
	if profile.WeightKG != nil {
		diag.WeightKG = *profile.WeightKG
		diag.Evidence["profile_weight_kg"] = *profile.WeightKG
	}
	if profile.FrontWeightPct != nil {
		diag.FrontWeightPct = *profile.FrontWeightPct
		diag.Evidence["profile_front_weight_pct"] = *profile.FrontWeightPct
	}
	if profile.PowerToWeightKWPerKG != nil {
		diag.PowerToWeightKWPerKG = *profile.PowerToWeightKWPerKG
		diag.PowerToWeightBand = powerToWeightBand(*profile.PowerToWeightKWPerKG)
		diag.Evidence["profile_power_to_weight_kw_per_kg"] = *profile.PowerToWeightKWPerKG
	}
	if profile.PeakTorqueRPM != nil {
		diag.PeakTorqueRPM = *profile.PeakTorqueRPM
		diag.Evidence["profile_peak_torque_rpm"] = *profile.PeakTorqueRPM
	}
	if profile.PeakPowerRPM != nil {
		diag.PeakPowerRPM = *profile.PeakPowerRPM
		diag.Evidence["profile_peak_power_rpm"] = *profile.PeakPowerRPM
	}
	if profile.RedlineRPM != nil {
		diag.RedlineRPM = *profile.RedlineRPM
		diag.Evidence["profile_redline_rpm"] = *profile.RedlineRPM
	}
}

func powerToWeightBand(value float64) string {
	switch {
	case value <= 0:
		return ""
	case value < 0.14:
		return "low"
	case value < 0.22:
		return "medium"
	case value < 0.32:
		return "high"
	default:
		return "extreme"
	}
}

func topSpeedFinding(bands []GearPowerBand) string {
	var highest *GearPowerBand
	for i := range bands {
		if bands[i].HighLoadSampleCount == 0 {
			continue
		}
		if highest == nil || bands[i].Gear > highest.Gear {
			highest = &bands[i]
		}
	}
	if highest == nil || highest.SpeedMaxKmh < 160 {
		return ""
	}
	if highest.RpmRatioMax > 0.95 {
		return "top_speed_limited_by_gearing"
	}
	if highest.RpmRatioAvg < 0.58 && highest.ThrottleAvg > 0.65 {
		return "top_speed_bog_down"
	}
	return "top_speed_ok"
}

func launchFindingForGroups(groups []SessionIssueGroup) string {
	for _, group := range groups {
		if group.Family != "launch_traction" {
			continue
		}
		for _, eventType := range group.EventTypes {
			if eventType == "launch_wheelspin" {
				return "launch_wheelspin"
			}
			if eventType == "launch_bog_down" {
				return "launch_bog_down"
			}
		}
	}
	return ""
}

func gearActionItem(gear int) string {
	if gear < 1 || gear > 10 {
		return "current_gear"
	}
	return fmt.Sprintf("gear_%d", gear)
}

func profileGearUnlocked(profile *TuneProfile, gear int) bool {
	if gear < 1 || gear > 10 {
		return false
	}
	if profile == nil {
		return true
	}
	return tuneProfileFloatPointer(*profile, fmt.Sprintf("gear%d", gear)) != nil
}

func safeDiv(value float64, divisor float64) float64 {
	if divisor == 0 {
		return 0
	}
	return value / divisor
}

func BuildWholeCarTuningPlan(groups []SessionIssueGroup, gear GearPowerDiagnostic, profile *TuneProfile) WholeCarTuningPlan {
	plan := WholeCarTuningPlan{
		Strategy:   wholeCarStrategy(groups),
		Confidence: "medium",
		Summary:    "whole_car_template",
	}
	if profile == nil {
		plan.Confidence = "needs_profile"
		plan.Notes = append(plan.Notes, "bind_tune_profile_for_concrete_values")
	}
	candidates := wholeCarCandidates(groups, gear)
	selected := map[string]wholeCarCandidate{}
	for _, candidate := range candidates {
		key := actionConflictKey(candidate.action, candidate.evidence)
		if key == "" {
			key = candidate.action.Category + "/" + candidate.action.Item
		}
		existing, ok := selected[key]
		if !ok {
			selected[key] = candidate
			continue
		}
		kept, dropped := chooseWholeCarCandidate(existing, candidate, groups)
		selected[key] = kept
		if kept.action.Item != dropped.action.Item || kept.action.Direction != dropped.action.Direction {
			plan.Conflicts = append(plan.Conflicts, TuningConflict{
				Key:         key,
				KeptItem:    kept.action.Item + "/" + kept.action.Direction,
				DroppedItem: dropped.action.Item + "/" + dropped.action.Direction,
				Reason:      "resolved_by_whole_car_priority",
			})
		}
	}
	out := make([]wholeCarCandidate, 0, len(selected))
	for _, candidate := range selected {
		out = append(out, candidate)
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].score != out[j].score {
			return out[i].score > out[j].score
		}
		return out[i].action.Priority < out[j].action.Priority
	})
	if len(out) > 5 {
		out = out[:5]
	}
	for index, candidate := range out {
		action := candidate.action
		plan.Actions = append(plan.Actions, WholeCarAdjustment{
			Priority:   index,
			Family:     candidate.family,
			Source:     candidate.source,
			Confidence: candidate.confidence,
			Category:   action.Category,
			Item:       action.Item,
			Direction:  action.Direction,
			Amount:     action.Amount,
			Reason:     action.Reason,
			Evidence:   cloneEvidenceMap(candidate.evidence),
		})
	}
	if len(plan.Actions) == 0 {
		plan.Confidence = "low"
		plan.Summary = "no_clear_whole_car_action"
	}
	if hasRollbackGroup(groups) {
		plan.Strategy = "rollback_first"
		plan.Summary = "rollback_before_more_changes"
	}
	return plan
}

type wholeCarCandidate struct {
	action     telemetry.SuggestedAction
	family     string
	source     string
	confidence string
	evidence   map[string]float64
	score      float64
}

func wholeCarCandidates(groups []SessionIssueGroup, gear GearPowerDiagnostic) []wholeCarCandidate {
	candidates := []wholeCarCandidate{}
	for _, group := range groups {
		evidence := representativeIssueEvidence(group)
		groupScore := issueGroupScore(group)
		if group.PrioritizeTuning {
			groupScore += 20
		}
		for _, action := range group.PrimaryActions {
			candidates = append(candidates, wholeCarCandidate{
				action:     action,
				family:     group.Family,
				source:     "issue_group",
				confidence: confidenceForGroup(group),
				evidence:   evidence,
				score:      groupScore - float64(action.Priority)*3,
			})
		}
	}
	for _, action := range gear.RecommendedActions {
		evidence := cloneEvidenceMap(gear.Evidence)
		if gearFromAction(action.Item) > 0 {
			evidence["gear"] = float64(gearFromAction(action.Item))
		}
		candidates = append(candidates, wholeCarCandidate{
			action:     action,
			family:     "gearing_acceleration",
			source:     "gear_power_diagnostic",
			confidence: "medium",
			evidence:   evidence,
			score:      70 - float64(action.Priority)*2,
		})
	}
	sort.SliceStable(candidates, func(i, j int) bool {
		if candidates[i].score != candidates[j].score {
			return candidates[i].score > candidates[j].score
		}
		return candidates[i].action.Priority < candidates[j].action.Priority
	})
	return candidates
}

func chooseWholeCarCandidate(left wholeCarCandidate, right wholeCarCandidate, groups []SessionIssueGroup) (wholeCarCandidate, wholeCarCandidate) {
	if actionConflictKey(left.action, left.evidence) == "tire_pressure" && actionConflictKey(right.action, right.evidence) == "tire_pressure" {
		if hasIssueFamily(groups, "tire_temperature_stability") {
			if right.action.Direction == "increase" {
				return right, left
			}
			if left.action.Direction == "increase" {
				return left, right
			}
		}
	}
	if right.score > left.score {
		return right, left
	}
	return left, right
}

func actionConflictKey(action telemetry.SuggestedAction, evidence map[string]float64) string {
	switch action.Item {
	case "tire_pressure", "drive_tire_pressure", "front_tire_pressure", "rear_tire_pressure":
		return "tire_pressure"
	case "drive_diff_accel":
		return "drive_diff_accel"
	case "front_diff_accel", "rear_diff_accel":
		return action.Item
	case "drive_diff_decel":
		return "drive_diff_decel"
	case "front_diff_decel", "rear_diff_decel":
		return action.Item
	case "current_gear":
		gear := int(evidence["gear"])
		if gear >= 1 && gear <= 10 {
			return fmt.Sprintf("gear_%d", gear)
		}
		return "current_gear"
	case "gear_1", "gear_2", "gear_3", "gear_4", "gear_5", "gear_6", "gear_7", "gear_8", "gear_9", "gear_10":
		return action.Item
	default:
		fields := tuneFieldsForAction(action.Item)
		if len(fields) == 0 {
			return action.Item
		}
		return strings.Join(fields, "+")
	}
}

func representativeIssueEvidence(group SessionIssueGroup) map[string]float64 {
	out := map[string]float64{}
	for key, stat := range group.Evidence {
		out[key] = stat.Avg
	}
	if gear := gearFromGroup(group); gear > 0 {
		out["gear"] = float64(gear)
	}
	return out
}

func gearFromGroup(group SessionIssueGroup) int {
	if stat, ok := group.Evidence["gear"]; ok {
		gear := int(stat.Avg + 0.5)
		if gear >= 1 && gear <= 10 {
			return gear
		}
	}
	for _, event := range group.Events {
		if gear := int(event.Evidence["gear"] + 0.5); gear >= 1 && gear <= 10 {
			return gear
		}
	}
	return 0
}

func gearFromAction(item string) int {
	if strings.HasPrefix(item, "gear_") {
		gear, _ := strconv.Atoi(strings.TrimPrefix(item, "gear_"))
		if gear >= 1 && gear <= 10 {
			return gear
		}
	}
	return 0
}

func isGearTuneField(field string) bool {
	if field == "finalDrive" {
		return true
	}
	return gearNumberForTuneField(field) > 0
}

func gearTuneItemForField(field string) string {
	if field == "finalDrive" {
		return "final_drive"
	}
	if gear := gearNumberForTuneField(field); gear > 0 {
		return gearActionItem(gear)
	}
	return field
}

func gearNumberForTuneField(field string) int {
	if !strings.HasPrefix(field, "gear") {
		return 0
	}
	gear, _ := strconv.Atoi(strings.TrimPrefix(field, "gear"))
	if gear >= 1 && gear <= 10 {
		return gear
	}
	return 0
}

func confidenceForGroup(group SessionIssueGroup) string {
	if group.Comparison == issueCompareWorsened || group.Comparison == issueCompareImproved {
		return "high"
	}
	if group.EventCount >= 3 {
		return "medium"
	}
	return "low"
}

func wholeCarStrategy(groups []SessionIssueGroup) string {
	if hasRollbackGroup(groups) {
		return "rollback_first"
	}
	if len(groups) >= 3 {
		return "coarse_whole_car"
	}
	return "targeted_whole_car"
}

func hasRollbackGroup(groups []SessionIssueGroup) bool {
	for _, group := range groups {
		if group.FeedbackDirective == "rollback_related_changes" {
			return true
		}
	}
	return false
}

func hasIssueFamily(groups []SessionIssueGroup, family string) bool {
	for _, group := range groups {
		if group.Family == family {
			return true
		}
	}
	return false
}

func cloneEvidenceMap(input map[string]float64) map[string]float64 {
	if len(input) == 0 {
		return nil
	}
	out := make(map[string]float64, len(input))
	for key, value := range input {
		out[key] = value
	}
	return out
}

func adjustmentStrategyForGroup(group SessionIssueGroup) string {
	if group.Comparison == issueCompareWorsened && len(group.RelatedRecentChanges) > 0 {
		return "rollback_first"
	}
	if group.EventCount >= 8 || group.TotalDurationMS >= 10000 {
		return "coarse_combination"
	}
	if group.EventCount >= 3 || group.TotalDurationMS >= 3000 {
		return "medium_combination"
	}
	return "fine_tune"
}

func feedbackDirectiveForGroup(group SessionIssueGroup) string {
	if group.Comparison == issueCompareWorsened && len(group.RelatedRecentChanges) > 0 {
		return "rollback_related_changes"
	}
	if group.Comparison == issueCompareImproved {
		return "keep_direction_then_fine_tune"
	}
	if len(group.RelatedRecentChanges) > 0 && group.Comparison == issueCompareUnchanged {
		return "avoid_more_same_direction"
	}
	return ""
}

func scaledStrategyActions(group SessionIssueGroup) []telemetry.SuggestedAction {
	actions := mergeStrategyActions(group.PrimaryActions, strategyTemplateActionsForGroup(group))
	factor := 1.0
	switch group.AdjustmentStrategy {
	case "coarse_combination":
		factor = 2.0
	case "medium_combination":
		factor = 1.5
	}
	for i := range actions {
		actions[i].Priority = i
		actions[i].Amount = scaleActionAmount(actions[i].Amount, factor)
	}
	if len(actions) > 4 {
		return actions[:4]
	}
	return actions
}

func mergeStrategyActions(preferred []telemetry.SuggestedAction, fallback []telemetry.SuggestedAction) []telemetry.SuggestedAction {
	seen := map[string]bool{}
	out := make([]telemetry.SuggestedAction, 0, len(preferred)+len(fallback))
	for _, action := range append(preferred, fallback...) {
		key := strings.Join([]string{action.Category, action.Item, action.Direction}, "/")
		if seen[key] {
			continue
		}
		seen[key] = true
		out = append(out, action)
	}
	return out
}

func strategyTemplateActionsForGroup(group SessionIssueGroup) []telemetry.SuggestedAction {
	switch group.Family {
	case "brake_balance":
		if issueGroupHasEvent(group, "rear_brake_lockup") {
			return []telemetry.SuggestedAction{
				actionTemplate(0, "brake", "brake_balance", "increase", "1", "reduce rear lockup tendency"),
				actionTemplate(1, "differential", "rear_diff_decel", "increase", "2", "stabilize the rear axle while braking"),
				actionTemplate(2, "damping", "rear_rebound", "decrease", "0.3", "improve rear compliance under braking"),
			}
		}
	case "corner_entry_balance":
		return []telemetry.SuggestedAction{
			actionTemplate(0, "suspension", "front_arb", "decrease", "0.6", "increase front grip on entry"),
			actionTemplate(1, "suspension", "rear_arb", "increase", "0.5", "rotate the car more in steady cornering"),
			actionTemplate(2, "damping", "front_rebound", "decrease", "0.3", "let the front tires load more smoothly"),
			actionTemplate(3, "brake", "brake_balance", "decrease", "1", "reduce front lockup tendency"),
		}
	case "mid_corner_balance":
		if issueGroupHasEvent(group, "high_speed_four_wheel_slide") {
			return []telemetry.SuggestedAction{
				actionTemplate(0, "aero", "front_and_rear_aero", "increase", "2", "increase high-speed grip"),
				actionTemplate(1, "suspension", "ride_height", "check", "avoid bottoming", "prevent aero and suspension instability"),
				actionTemplate(2, "suspension", "front_arb", "decrease", "0.4", "reduce sustained tire scrub"),
				actionTemplate(3, "suspension", "rear_arb", "increase", "0.4", "rotate the car more in steady cornering"),
			}
		}
		return []telemetry.SuggestedAction{
			actionTemplate(0, "suspension", "front_arb", "decrease", "0.6", "increase front grip on steady cornering"),
			actionTemplate(1, "suspension", "rear_arb", "increase", "0.6", "rotate the car more in steady cornering"),
			actionTemplate(2, "alignment", "front_camber", "check", "slightly more negative", "improve front tire contact in cornering"),
			actionTemplate(3, "aero", "front_and_rear_aero", "increase", "2", "increase high-speed grip"),
		}
	case "corner_exit_power":
		if issueGroupHasEvent(group, "snap_oversteer") {
			return []telemetry.SuggestedAction{
				actionTemplate(0, "damping", "rear_rebound", "decrease", "0.3", "make rear response less abrupt"),
				actionTemplate(1, "suspension", "rear_arb", "decrease", "0.5", "increase rear grip"),
				actionTemplate(2, "differential", "rear_diff_decel", "increase", "2", "stabilize the rear axle while off throttle"),
			}
		}
		if issueGroupHasEvent(group, "power_understeer") {
			return []telemetry.SuggestedAction{
				actionTemplate(0, "differential", "drive_diff_accel", "decrease", "3", "reduce power-on understeer"),
				actionTemplate(1, "gearing", "current_gear", "decrease", "0.05", "reduce wheel torque on exit"),
				actionTemplate(2, "suspension", "front_arb", "decrease", "0.5", "increase front grip under power"),
				actionTemplate(3, "suspension", "rear_arb", "increase", "0.5", "rotate the car more in steady cornering"),
			}
		}
		return []telemetry.SuggestedAction{
			actionTemplate(0, "differential", "drive_diff_accel", "decrease", "3", "reduce power oversteer"),
			actionTemplate(1, "gearing", "current_gear", "decrease", "0.05", "reduce wheel torque on exit"),
			actionTemplate(2, "suspension", "rear_arb", "decrease", "0.5", "increase rear grip"),
			actionTemplate(3, "damping", "rear_rebound", "decrease", "0.3", "make rear response less abrupt"),
		}
	}
	return strategyTemplateActions(group.Family)
}

func issueGroupHasEvent(group SessionIssueGroup, eventType string) bool {
	for _, item := range group.EventTypes {
		if item == eventType {
			return true
		}
	}
	return false
}

func strategyTemplateActions(family string) []telemetry.SuggestedAction {
	switch family {
	case "launch_traction":
		return []telemetry.SuggestedAction{
			actionTemplate(0, "differential", "drive_diff_accel", "decrease", "4", "reduce driven-wheel slip"),
			actionTemplate(1, "gearing", "gear_1", "decrease", "0.08", "reduce wheel torque during launch"),
			actionTemplate(2, "tire", "drive_tire_pressure", "decrease", "0.03 BAR (≈0.5 PSI)", "increase launch traction"),
		}
	case "gearing_acceleration":
		return []telemetry.SuggestedAction{
			actionTemplate(0, "gearing", "current_gear", "increase", "0.05", "help the engine stay in the power band"),
			actionTemplate(1, "gearing", "final_drive", "increase", "0.08", "shorten road acceleration gearing"),
		}
	case "brake_balance":
		return []telemetry.SuggestedAction{
			actionTemplate(0, "brake", "brake_balance", "decrease", "1", "reduce front lockup tendency"),
			actionTemplate(1, "brake", "brake_pressure", "decrease", "2", "make threshold braking easier"),
		}
	case "corner_entry_balance":
		return []telemetry.SuggestedAction{
			actionTemplate(0, "suspension", "front_arb", "decrease", "0.6", "increase front grip on entry"),
			actionTemplate(1, "suspension", "rear_arb", "increase", "0.5", "rotate the car more in steady cornering"),
			actionTemplate(2, "damping", "front_rebound", "decrease", "0.3", "let the front tires load more smoothly"),
			actionTemplate(3, "brake", "brake_balance", "decrease", "1", "reduce front lockup tendency"),
		}
	case "mid_corner_balance":
		return []telemetry.SuggestedAction{
			actionTemplate(0, "suspension", "front_arb", "decrease", "0.6", "increase front grip on entry"),
			actionTemplate(1, "suspension", "rear_arb", "increase", "0.6", "rotate the car more in steady cornering"),
			actionTemplate(2, "tire", "tire_pressure", "increase", "0.03 BAR (≈0.5 PSI)", "stabilize tire temperature and contact patch"),
			actionTemplate(3, "aero", "front_and_rear_aero", "increase", "2", "increase high-speed grip"),
		}
	case "corner_exit_power":
		return []telemetry.SuggestedAction{
			actionTemplate(0, "differential", "drive_diff_accel", "decrease", "3", "reduce power oversteer"),
			actionTemplate(1, "gearing", "current_gear", "decrease", "0.05", "reduce wheel torque on exit"),
			actionTemplate(2, "suspension", "rear_arb", "decrease", "0.5", "increase rear grip"),
		}
	case "suspension_platform":
		return []telemetry.SuggestedAction{
			actionTemplate(0, "suspension", "ride_height", "increase", "0.2", "restore suspension travel"),
			actionTemplate(1, "suspension", "spring_rate", "increase", "0.3", "reduce bottoming frequency"),
			actionTemplate(2, "damping", "bump", "increase", "0.2", "support compression on impacts"),
		}
	case "tire_temperature_stability":
		return []telemetry.SuggestedAction{
			actionTemplate(0, "tire", "tire_pressure", "increase", "0.03 BAR (≈0.5 PSI)", "stabilize tire temperature and contact patch"),
			actionTemplate(1, "alignment", "front_camber", "check", "slightly more negative", "improve front tire contact in cornering"),
			actionTemplate(2, "suspension", "front_arb", "decrease", "0.4", "increase front grip on entry"),
		}
	default:
		return nil
	}
}

func actionTemplate(priority int, category string, item string, direction string, amount string, reason string) telemetry.SuggestedAction {
	return telemetry.SuggestedAction{Priority: priority, Category: category, Item: item, Direction: direction, Amount: amount, Reason: reason}
}

func scaleActionAmount(amount string, factor float64) string {
	if factor <= 1 {
		return amount
	}
	if amount == "slightly more negative" {
		return amount
	}
	spans := numberSpans(amount)
	if len(spans) == 0 {
		return amount
	}
	out := amount
	for i := len(spans) - 1; i >= 0; i-- {
		span := spans[i]
		replacement := formatScaledAmount(span.value * factor)
		out = out[:span.start] + replacement + out[span.end:]
	}
	return out
}

type numberSpan struct {
	value      float64
	start, end int
}

func numberSpans(value string) []numberSpan {
	spans := []numberSpan{}
	start := -1
	end := -1
	flush := func() {
		if start < 0 || end <= start {
			return
		}
		parsed, _ := strconv.ParseFloat(value[start:end], 64)
		if parsed > 0 {
			spans = append(spans, numberSpan{value: parsed, start: start, end: end})
		}
		start = -1
		end = -1
	}
	for index, char := range value {
		if (char >= '0' && char <= '9') || char == '.' {
			if start < 0 {
				start = index
			}
			end = index + 1
			continue
		}
		if start >= 0 {
			flush()
		}
	}
	flush()
	return spans
}

func formatScaledAmount(value float64) string {
	if value >= 10 || value == float64(int64(value)) {
		return fmt.Sprintf("%.0f", value)
	}
	return strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.2f", value), "0"), ".")
}

func rollbackActionsForGroup(group SessionIssueGroup, recentDeltas map[string]float64) []telemetry.SuggestedAction {
	out := []telemetry.SuggestedAction{}
	for _, field := range group.RelatedRecentChanges {
		delta := recentDeltas[field]
		if delta == 0 {
			continue
		}
		action, ok := rollbackActionForField(field, delta)
		if !ok {
			continue
		}
		action.Priority = len(out)
		out = append(out, action)
		if len(out) >= 4 {
			break
		}
	}
	return out
}

func rollbackActionForField(field string, delta float64) (telemetry.SuggestedAction, bool) {
	category, item, ok := actionTargetForTuneField(field)
	if !ok {
		return telemetry.SuggestedAction{}, false
	}
	direction := "decrease"
	if delta < 0 {
		direction = "increase"
	}
	return actionTemplate(0, category, item, direction, formatScaledAmount(absFloat(delta)/2), "rollback half of the last related change"), true
}

func actionTargetForTuneField(field string) (string, string, bool) {
	switch field {
	case "frontTirePressure":
		return "tire", "front_tire_pressure", true
	case "rearTirePressure":
		return "tire", "rear_tire_pressure", true
	case "finalDrive":
		return "gearing", "final_drive", true
	case "gear1":
		return "gearing", "gear_1", true
	case "gear2", "gear3", "gear4", "gear5", "gear6", "gear7", "gear8", "gear9", "gear10":
		return "gearing", "gear_" + strings.TrimPrefix(field, "gear"), true
	case "frontCamber":
		return "alignment", "front_camber", true
	case "frontArb":
		return "suspension", "front_arb", true
	case "rearArb":
		return "suspension", "rear_arb", true
	case "frontRebound":
		return "damping", "front_rebound", true
	case "rearRebound":
		return "damping", "rear_rebound", true
	case "frontBump", "rearBump":
		return "damping", "bump", true
	case "frontSpring", "rearSpring":
		return "suspension", "spring_rate", true
	case "frontRideHeight", "rearRideHeight":
		return "suspension", "ride_height", true
	case "frontAero", "rearAero":
		return "aero", "front_and_rear_aero", true
	case "brakeBalance":
		return "brake", "brake_balance", true
	case "brakePressure":
		return "brake", "brake_pressure", true
	case "rearDiffAccel":
		return "differential", "rear_diff_accel", true
	case "rearDiffDecel":
		return "differential", "rear_diff_decel", true
	case "frontDiffAccel":
		return "differential", "front_diff_accel", true
	case "frontDiffDecel":
		return "differential", "front_diff_decel", true
	default:
		return "", "", false
	}
}

func absFloat(value float64) float64 {
	if value < 0 {
		return -value
	}
	return value
}

func issueFamilyForEvent(eventType string) string {
	switch strings.TrimSpace(eventType) {
	case "launch_wheelspin", "launch_bog_down":
		return "launch_traction"
	case "short_gear", "long_gear_bog_down", "top_speed_limited_by_gearing":
		return "gearing_acceleration"
	case "front_brake_lockup", "rear_brake_lockup":
		return "brake_balance"
	case "corner_entry_understeer":
		return "corner_entry_balance"
	case "mid_corner_understeer", "high_speed_four_wheel_slide":
		return "mid_corner_balance"
	case "corner_exit_oversteer", "power_understeer", "snap_oversteer":
		return "corner_exit_power"
	case "suspension_bottom_out":
		return "suspension_platform"
	case "tire_overheat", "tire_temp_imbalance":
		return "tire_temperature_stability"
	default:
		return "driver_execution"
	}
}

func shouldPrioritizeIssueFamily(family string, severity string, count int, relatedRecentChange bool) bool {
	if severity == "high" || relatedRecentChange {
		return true
	}
	switch family {
	case "launch_traction", "gearing_acceleration", "brake_balance", "suspension_platform":
		return count >= 1
	default:
		return count >= 2
	}
}

func relatedRecentChangesForActions(actions []telemetry.SuggestedAction, recentChanges []string) []string {
	if len(actions) == 0 || len(recentChanges) == 0 {
		return nil
	}
	recent := map[string]bool{}
	for _, field := range recentChanges {
		recent[field] = true
	}
	relatedSet := map[string]bool{}
	for _, action := range actions {
		for _, field := range tuneFieldsForAction(action.Item) {
			if recent[field] {
				relatedSet[field] = true
			}
		}
	}
	related := make([]string, 0, len(relatedSet))
	for field := range relatedSet {
		related = append(related, field)
	}
	sort.Strings(related)
	return related
}

func tuneFieldsForAction(item string) []string {
	switch item {
	case "front_tire_pressure":
		return []string{"frontTirePressure"}
	case "rear_tire_pressure":
		return []string{"rearTirePressure"}
	case "gear_1":
		return []string{"gear1"}
	case "gear_2":
		return []string{"gear2"}
	case "gear_3":
		return []string{"gear3"}
	case "gear_4":
		return []string{"gear4"}
	case "gear_5":
		return []string{"gear5"}
	case "gear_6":
		return []string{"gear6"}
	case "gear_7":
		return []string{"gear7"}
	case "gear_8":
		return []string{"gear8"}
	case "gear_9":
		return []string{"gear9"}
	case "gear_10":
		return []string{"gear10"}
	case "current_gear":
		return []string{"gear1", "gear2", "gear3", "gear4", "gear5", "gear6", "gear7", "gear8", "gear9", "gear10"}
	case "final_drive":
		return []string{"finalDrive"}
	case "brake_balance":
		return []string{"brakeBalance"}
	case "brake_pressure":
		return []string{"brakePressure"}
	case "rear_diff_accel":
		return []string{"rearDiffAccel"}
	case "rear_diff_decel":
		return []string{"rearDiffDecel"}
	case "front_diff_accel":
		return []string{"frontDiffAccel"}
	case "front_diff_decel":
		return []string{"frontDiffDecel"}
	case "drive_diff_accel":
		return []string{"frontDiffAccel", "rearDiffAccel"}
	case "drive_diff_decel":
		return []string{"frontDiffDecel", "rearDiffDecel"}
	case "drive_tire_pressure", "tire_pressure":
		return []string{"frontTirePressure", "rearTirePressure"}
	case "front_arb":
		return []string{"frontArb"}
	case "rear_arb":
		return []string{"rearArb"}
	case "front_rebound":
		return []string{"frontRebound"}
	case "rear_rebound":
		return []string{"rearRebound"}
	case "front_camber":
		return []string{"frontCamber"}
	case "front_and_rear_aero":
		return []string{"frontAero", "rearAero"}
	case "ride_height":
		return []string{"frontRideHeight", "rearRideHeight"}
	case "spring_rate":
		return []string{"frontSpring", "rearSpring"}
	case "bump":
		return []string{"frontBump", "rearBump"}
	default:
		return nil
	}
}

func (s *Store) findIssueBaselineSession(session TelemetrySession) (*TelemetrySession, string, error) {
	sessions, err := s.ListTelemetrySessions(500)
	if err != nil {
		return nil, issueCompareUnavailable, err
	}
	driverMode := NormalizeTestConditions(TestConditions{DriverMode: session.DriverMode}).DriverMode
	trackName := strings.TrimSpace(session.TrackName)
	startedAt := strings.TrimSpace(session.StartedAt)
	sessionUseCase := s.useCaseForSession(session)
	for _, candidate := range sessions {
		if candidate.ID == session.ID || !sessionStartedBefore(candidate.StartedAt, startedAt) {
			continue
		}
		if !sameTuneProfile(candidate.TuneProfileID, session.TuneProfileID) {
			continue
		}
		if strings.TrimSpace(candidate.TrackName) != trackName {
			continue
		}
		if NormalizeTestConditions(TestConditions{DriverMode: candidate.DriverMode}).DriverMode != driverMode {
			continue
		}
		return &candidate, "matched_profile_track_driver", nil
	}
	for _, candidate := range sessions {
		if candidate.ID == session.ID || !sessionStartedBefore(candidate.StartedAt, startedAt) {
			continue
		}
		if !sameVehicleSnapshot(candidate, session) {
			continue
		}
		if NormalizeTestConditions(TestConditions{DriverMode: candidate.DriverMode}).DriverMode != driverMode {
			continue
		}
		if sessionUseCase != "" {
			candidateUseCase := s.useCaseForSession(candidate)
			if candidateUseCase != "" && !strings.EqualFold(candidateUseCase, sessionUseCase) {
				continue
			}
		}
		return &candidate, "matched_vehicle_class_usecase_driver", nil
	}
	return nil, "missing_comparison_baseline", nil
}

func sessionStartedBefore(candidate string, current string) bool {
	if strings.TrimSpace(current) == "" {
		return true
	}
	if strings.TrimSpace(candidate) == "" {
		return false
	}
	return candidate < current
}

func sameVehicleSnapshot(left TelemetrySession, right TelemetrySession) bool {
	return left.CarOrdinal != nil && right.CarOrdinal != nil && *left.CarOrdinal == *right.CarOrdinal && strings.EqualFold(strings.TrimSpace(left.CarClass), strings.TrimSpace(right.CarClass))
}

func (s *Store) useCaseForSession(session TelemetrySession) string {
	if profile, err := ParseTuneProfileSnapshotJSON(session.TuneSnapshotJSON); err == nil && profile != nil {
		return strings.TrimSpace(profile.UseCase)
	}
	if session.TuneProfileID != nil {
		if profile, err := s.GetTuneProfile(*session.TuneProfileID); err == nil && profile != nil {
			return strings.TrimSpace(profile.UseCase)
		}
	}
	return ""
}

func (s *Store) issueTuneProfile(session TelemetrySession) *TuneProfile {
	if profile, err := ParseTuneProfileSnapshotJSON(session.TuneSnapshotJSON); err == nil && profile != nil {
		return profile
	}
	if session.TuneProfileID != nil {
		if profile, err := s.GetTuneProfile(*session.TuneProfileID); err == nil && profile != nil {
			return profile
		}
	}
	return nil
}

func (s *Store) recentChangeFieldsForSession(session TelemetrySession) []string {
	snapshot := s.recentSnapshotForSession(session)
	if snapshot == nil {
		return nil
	}
	return append([]string(nil), snapshot.ChangedFields...)
}

func (s *Store) recentChangeDeltasForSession(session TelemetrySession) map[string]float64 {
	snapshot := s.recentSnapshotForSession(session)
	if snapshot == nil {
		return nil
	}
	deltas := map[string]float64{}
	for _, field := range snapshot.ChangedFields {
		before, okBefore := tuneProfileNumericField(snapshot.Before, field)
		after, okAfter := tuneProfileNumericField(snapshot.After, field)
		if okBefore && okAfter && before != after {
			deltas[field] = after - before
		}
	}
	return deltas
}

func (s *Store) recentSnapshotForSession(session TelemetrySession) *TuneProfileSnapshot {
	if session.TuneProfileID == nil {
		return nil
	}
	query := `SELECT id, tune_profile_id, session_id, changed_at, COALESCE(change_reason, ''), change_json
		FROM tune_change_log WHERE tune_profile_id = ?`
	args := []any{*session.TuneProfileID}
	if strings.TrimSpace(session.StartedAt) != "" {
		query += ` AND changed_at <= ?`
		args = append(args, session.StartedAt)
	}
	query += ` ORDER BY changed_at DESC, id DESC LIMIT 1`
	snapshot, err := scanTuneProfileSnapshot(s.db.QueryRow(query, args...))
	if err != nil {
		return nil
	}
	return &snapshot
}

func tuneProfileNumericField(profile TuneProfile, field string) (float64, bool) {
	value := tuneProfileFloatPointer(profile, field)
	if value == nil {
		return 0, false
	}
	return *value, true
}

func tuneProfileFloatPointer(profile TuneProfile, field string) *float64 {
	switch field {
	case "powerKW":
		return profile.PowerKW
	case "torqueNM":
		return profile.TorqueNM
	case "weightKG":
		return profile.WeightKG
	case "frontWeightPct":
		return profile.FrontWeightPct
	case "powerToWeightKWPerKG":
		return profile.PowerToWeightKWPerKG
	case "peakTorqueRPM":
		return profile.PeakTorqueRPM
	case "peakPowerRPM":
		return profile.PeakPowerRPM
	case "redlineRPM":
		return profile.RedlineRPM
	case "frontTirePressure":
		return profile.FrontTirePressure
	case "rearTirePressure":
		return profile.RearTirePressure
	case "finalDrive":
		return profile.FinalDrive
	case "gear1":
		return profile.Gear1
	case "gear2":
		return profile.Gear2
	case "gear3":
		return profile.Gear3
	case "gear4":
		return profile.Gear4
	case "gear5":
		return profile.Gear5
	case "gear6":
		return profile.Gear6
	case "gear7":
		return profile.Gear7
	case "gear8":
		return profile.Gear8
	case "gear9":
		return profile.Gear9
	case "gear10":
		return profile.Gear10
	case "frontCamber":
		return profile.FrontCamber
	case "rearCamber":
		return profile.RearCamber
	case "frontToe":
		return profile.FrontToe
	case "rearToe":
		return profile.RearToe
	case "caster":
		return profile.Caster
	case "frontArb":
		return profile.FrontARB
	case "rearArb":
		return profile.RearARB
	case "frontSpring":
		return profile.FrontSpring
	case "rearSpring":
		return profile.RearSpring
	case "frontRideHeight":
		return profile.FrontRideHeight
	case "rearRideHeight":
		return profile.RearRideHeight
	case "frontRebound":
		return profile.FrontRebound
	case "rearRebound":
		return profile.RearRebound
	case "frontBump":
		return profile.FrontBump
	case "rearBump":
		return profile.RearBump
	case "frontAero":
		return profile.FrontAero
	case "rearAero":
		return profile.RearAero
	case "aeroBalance":
		return profile.AeroBalance
	case "brakeBalance":
		return profile.BrakeBalance
	case "brakePressure":
		return profile.BrakePressure
	case "frontDiffAccel":
		return profile.FrontDiffAccel
	case "frontDiffDecel":
		return profile.FrontDiffDecel
	case "rearDiffAccel":
		return profile.RearDiffAccel
	case "rearDiffDecel":
		return profile.RearDiffDecel
	case "centerDiffBalance":
		return profile.CenterDiffBalance
	default:
		return nil
	}
}

func (s *Store) ListStrategyTemplates() ([]StrategyTemplate, error) {
	profiles, err := s.ListRuleThresholdProfiles()
	if err != nil {
		return nil, err
	}
	templates := make([]StrategyTemplate, 0, len(profiles))
	for _, profile := range profiles {
		config, err := parseRuleConfigJSON(profile.ConfigJSON)
		if err != nil {
			return nil, err
		}
		templates = append(templates, strategyTemplateFromRuleProfile(profile, config))
	}
	return templates, nil
}

func (s *Store) AnalyzeRoadStrategySessions(sessionIDs []int64, strategyTemplateID int64) (*RoadStrategyAnalysis, error) {
	if len(sessionIDs) == 0 {
		return nil, errors.New("at least one telemetry session is required")
	}
	if len(sessionIDs) > 5 {
		return nil, errors.New("strategy analysis supports at most 5 sessions")
	}
	profile, err := s.GetRuleThresholdProfile(strategyTemplateID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("strategy template not found")
		}
		return nil, err
	}
	config, err := parseRuleConfigJSON(profile.ConfigJSON)
	if err != nil {
		return nil, err
	}
	template := strategyTemplateFromRuleProfile(*profile, config)
	analysis := &RoadStrategyAnalysis{
		Template:     template,
		SessionIDs:   append([]int64(nil), sessionIDs...),
		SessionCount: len(sessionIDs),
	}
	eventCounts := map[string]int{}
	eventSeverity := map[string]string{}
	familyCounts := map[string]int{}
	familySessions := map[string]map[int64]bool{}
	familySeverity := map[string]string{}
	for _, sessionID := range sessionIDs {
		samples, err := s.GetSessionTelemetrySamples(sessionID, 10000)
		if err != nil {
			return nil, err
		}
		if len(samples) == 0 {
			analysis.Hints = append(analysis.Hints, StrategyAnalysisHint{Level: "warning", Message: "session_has_no_samples"})
			continue
		}
		engine := telemetry.NewRuleEngine()
		engine.SetConfig(config)
		for _, sample := range samples {
			engine.Observe(sample)
		}
		events := engine.Events()
		analysis.TotalEvents += len(events)
		sessionFamilies := map[string]bool{}
		for _, event := range events {
			eventCounts[event.Type]++
			eventSeverity[event.Type] = maxSeverity(eventSeverity[event.Type], event.Severity)
			family := issueFamilyForEvent(event.Type)
			familyCounts[family]++
			familySeverity[family] = maxSeverity(familySeverity[family], event.Severity)
			sessionFamilies[family] = true
		}
		for family := range sessionFamilies {
			if familySessions[family] == nil {
				familySessions[family] = map[int64]bool{}
			}
			familySessions[family][sessionID] = true
		}
	}
	for eventType, count := range eventCounts {
		analysis.EventDistribution = append(analysis.EventDistribution, StrategyEventDistribution{Type: eventType, Count: count, Severity: eventSeverity[eventType]})
	}
	sort.SliceStable(analysis.EventDistribution, func(i, j int) bool {
		if analysis.EventDistribution[i].Count != analysis.EventDistribution[j].Count {
			return analysis.EventDistribution[i].Count > analysis.EventDistribution[j].Count
		}
		return analysis.EventDistribution[i].Type < analysis.EventDistribution[j].Type
	})
	for family, count := range familyCounts {
		sessionCount := len(familySessions[family])
		recommendation := "keep_current_thresholds"
		if sessionCount >= 3 && count >= sessionCount*3 {
			recommendation = "increase_adjustment_step_not_detection_threshold"
		} else if count == 0 {
			recommendation = "consider_lowering_threshold_if_driver_feedback_confirms"
		} else if sessionCount >= 2 {
			recommendation = "repeated_issue_confirmed"
		}
		analysis.IssueGroups = append(analysis.IssueGroups, StrategyIssueAggregate{
			Family:         family,
			EventCount:     count,
			SessionCount:   sessionCount,
			Severity:       familySeverity[family],
			Recommendation: recommendation,
		})
		if recommendation != "keep_current_thresholds" {
			analysis.Hints = append(analysis.Hints, StrategyAnalysisHint{Level: "info", Message: recommendation, Family: family})
		}
	}
	sort.SliceStable(analysis.IssueGroups, func(i, j int) bool {
		if analysis.IssueGroups[i].SessionCount != analysis.IssueGroups[j].SessionCount {
			return analysis.IssueGroups[i].SessionCount > analysis.IssueGroups[j].SessionCount
		}
		if analysis.IssueGroups[i].EventCount != analysis.IssueGroups[j].EventCount {
			return analysis.IssueGroups[i].EventCount > analysis.IssueGroups[j].EventCount
		}
		return analysis.IssueGroups[i].Family < analysis.IssueGroups[j].Family
	})
	if analysis.TotalEvents == 0 {
		analysis.Hints = append(analysis.Hints, StrategyAnalysisHint{Level: "warning", Message: "no_events_matched_selected_strategy"})
	}
	if len(analysis.EventDistribution) > 0 && analysis.EventDistribution[0].Count > len(sessionIDs)*8 {
		analysis.Hints = append(analysis.Hints, StrategyAnalysisHint{Level: "warning", Message: "possible_overmatching", EventType: analysis.EventDistribution[0].Type})
	}
	return analysis, nil
}

func strategyTemplateFromRuleProfile(profile RuleThresholdProfile, config telemetry.RuleConfig) StrategyTemplate {
	enabled := 0
	for _, event := range config.Events {
		if event.Enabled {
			enabled++
		}
	}
	return StrategyTemplate{
		ID:                profile.ID,
		Name:              profile.Name,
		CarClass:          profile.CarClass,
		Drivetrain:        profile.Drivetrain,
		UseCase:           profile.UseCase,
		GameMode:          profile.GameMode,
		IsDefault:         profile.IsDefault,
		EnabledEventCount: enabled,
		TotalEventCount:   len(config.Events),
		UpdatedAt:         profile.UpdatedAt,
	}
}

func issueFamilyLabel(family string) string {
	switch family {
	case "launch_traction":
		return "launch traction"
	case "gearing_acceleration":
		return "gearing and acceleration"
	case "brake_balance":
		return "brake balance"
	case "corner_entry_balance":
		return "corner entry balance"
	case "mid_corner_balance":
		return "sustained cornering balance"
	case "corner_exit_power":
		return "corner exit power"
	case "suspension_platform":
		return "suspension platform"
	case "tire_temperature_stability":
		return "tire temperature and stability"
	default:
		return fmt.Sprintf("%s", family)
	}
}
