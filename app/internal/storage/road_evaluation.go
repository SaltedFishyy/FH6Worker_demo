package storage

import (
	"math"
	"sort"
	"strings"

	"fh6worker/internal/telemetry"
)

const (
	roadVerdictGoodFit          = "good_fit"
	roadVerdictFastButRisky     = "fast_but_risky"
	roadVerdictPaperFastNotFit  = "paper_fast_not_fit"
	roadVerdictNeedsTuning      = "needs_tuning"
	roadVerdictInsufficientData = "insufficient_data"

	roadBaselineMatched            = "matched_auto_baseline"
	roadBaselineSelfAuto           = "self_auto_baseline"
	roadBaselineMissingAuto        = "missing_auto_baseline"
	roadBaselineMissingVehicle     = "missing_vehicle_identity"
	roadBaselineNoValidStandardRun = "no_valid_standard_run"
	roadBaselineNoStandardTrack    = "no_standard_track"
	roadAttributionTuneIssue       = "tune_issue"
	roadAttributionStyleFitIssue   = "style_fit_issue"
	roadAttributionDriverExecution = "driver_execution_issue"
	roadAttributionDataGap         = "data_gap"
)

func (s *Store) EvaluateRoadSession(sessionID int64) (*RoadSessionEvaluation, error) {
	session, err := s.GetTelemetrySession(sessionID)
	if err != nil {
		return nil, err
	}

	runs, err := s.listBenchmarkRunsForSession(sessionID)
	if err != nil {
		return nil, err
	}
	if len(runs) == 0 {
		runs, err = s.AnalyzeSessionBenchmarkRuns(sessionID)
		if err != nil {
			return nil, err
		}
	}
	events, err := s.GetSessionEvents(sessionID)
	if err != nil {
		return nil, err
	}

	evaluation := &RoadSessionEvaluation{
		Session:        *session,
		BaselineStatus: roadBaselineNoValidStandardRun,
		OverallVerdict: roadVerdictInsufficientData,
		Notes:          []string{},
	}

	best := bestValidRoadRun(runs)
	if best == nil {
		evaluation.Attributions = append(evaluation.Attributions, RoadEvaluationAttribution{
			Type:     roadAttributionDataGap,
			Priority: 1,
			Message:  roadBaselineNoValidStandardRun,
		})
		evaluation.Notes = append(evaluation.Notes, roadBaselineNoValidStandardRun)
		return evaluation, nil
	}
	evaluation.BestRun = best

	if track, err := s.GetBenchmarkTrack(best.TrackID); err == nil {
		evaluation.Track = track
	} else {
		evaluation.BaselineStatus = roadBaselineNoStandardTrack
		evaluation.Notes = append(evaluation.Notes, roadBaselineNoStandardTrack)
	}

	evaluation.RiskScore = roadRiskScore(*best, events)
	baselineRun, baselineSession, baselineStatus, err := s.findRoadAutoBaseline(*session, *best)
	if err != nil {
		return nil, err
	}
	evaluation.BaselineStatus = baselineStatus
	evaluation.BaselineRun = baselineRun
	evaluation.BaselineSession = baselineSession
	evaluation.PaperPerformanceScore = roadPaperScore(baselineRun, best)
	evaluation.PlayerFitScore = roadPlayerFitScore(*session, *best, baselineRun, evaluation.RiskScore)
	evaluation.OverallVerdict = roadVerdict(*session, *best, baselineRun, evaluation.RiskScore)
	evaluation.Attributions = roadAttributions(events, *best, evaluation.BaselineStatus, evaluation.RiskScore, roadPlayerIsFaster(*best, baselineRun))
	evaluation.Notes = append(evaluation.Notes, roadEvaluationNotes(evaluation)...)
	return evaluation, nil
}

func (s *Store) CompareRoadEvaluations(leftID int64, rightID int64) (*RoadEvaluationComparison, error) {
	left, err := s.EvaluateRoadSession(leftID)
	if err != nil {
		return nil, err
	}
	right, err := s.EvaluateRoadSession(rightID)
	if err != nil {
		return nil, err
	}
	comparison := &RoadEvaluationComparison{
		Left:  *left,
		Right: *right,
		Metrics: []SessionComparisonMetric{
			roadEvalMetric("paper_performance_score", "Paper performance", "pt", left.PaperPerformanceScore, right.PaperPerformanceScore, true),
			roadEvalMetric("player_fit_score", "Player fit", "pt", left.PlayerFitScore, right.PlayerFitScore, true),
			roadEvalMetric("risk_score", "Risk", "pt", left.RiskScore, right.RiskScore, false),
		},
		Notes: []string{},
	}
	if left.BestRun != nil && right.BestRun != nil {
		comparison.Metrics = append(comparison.Metrics, roadEvalMetric("best_run_duration", "Best standard segment", "ms", float64(left.BestRun.DurationMS), float64(right.BestRun.DurationMS), false))
	}
	if left.BaselineStatus != right.BaselineStatus {
		comparison.Notes = append(comparison.Notes, "baseline_status_mismatch")
	}
	if right.OverallVerdict == roadVerdictGoodFit || (right.PlayerFitScore > left.PlayerFitScore && right.RiskScore <= left.RiskScore) {
		comparison.Verdict = "right_improved"
	} else if left.OverallVerdict == roadVerdictGoodFit || (left.PlayerFitScore > right.PlayerFitScore && left.RiskScore <= right.RiskScore) {
		comparison.Verdict = "left_better"
	} else {
		comparison.Verdict = "mixed"
	}
	return comparison, nil
}

func (s *Store) findRoadAutoBaseline(session TelemetrySession, best BenchmarkRun) (*BenchmarkRun, *TelemetrySession, string, error) {
	driverMode := NormalizeTestConditions(TestConditions{DriverMode: session.DriverMode}).DriverMode
	if driverMode == "auto" {
		runCopy := best
		sessionCopy := session
		return &runCopy, &sessionCopy, roadBaselineSelfAuto, nil
	}
	if session.CarOrdinal == nil || *session.CarOrdinal <= 0 || strings.TrimSpace(session.CarClass) == "" {
		return nil, nil, roadBaselineMissingVehicle, nil
	}

	sessions, err := s.ListTelemetrySessions(500)
	if err != nil {
		return nil, nil, "", err
	}
	type candidate struct {
		run      BenchmarkRun
		session  TelemetrySession
		sameTune bool
	}
	candidates := make([]candidate, 0, 8)
	for _, possible := range sessions {
		if possible.ID == session.ID || possible.CarOrdinal == nil || *possible.CarOrdinal != *session.CarOrdinal || !strings.EqualFold(strings.TrimSpace(possible.CarClass), strings.TrimSpace(session.CarClass)) {
			continue
		}
		if NormalizeTestConditions(TestConditions{DriverMode: possible.DriverMode}).DriverMode != "auto" {
			continue
		}
		runs, err := s.listBenchmarkRunsForSession(possible.ID)
		if err != nil {
			return nil, nil, "", err
		}
		if len(runs) == 0 {
			runs, err = s.AnalyzeSessionBenchmarkRuns(possible.ID)
			if err != nil {
				return nil, nil, "", err
			}
		}
		for _, run := range runs {
			if run.TrackID != best.TrackID || !run.Valid || run.DurationMS <= 0 {
				continue
			}
			candidates = append(candidates, candidate{
				run:      run,
				session:  possible,
				sameTune: sameTuneProfile(session.TuneProfileID, possible.TuneProfileID),
			})
		}
	}
	if len(candidates) == 0 {
		return nil, nil, roadBaselineMissingAuto, nil
	}
	sort.SliceStable(candidates, func(i, j int) bool {
		if candidates[i].sameTune != candidates[j].sameTune {
			return candidates[i].sameTune
		}
		if candidates[i].run.DurationMS != candidates[j].run.DurationMS {
			return candidates[i].run.DurationMS < candidates[j].run.DurationMS
		}
		return candidates[i].run.Confidence > candidates[j].run.Confidence
	})
	runCopy := candidates[0].run
	sessionCopy := candidates[0].session
	return &runCopy, &sessionCopy, roadBaselineMatched, nil
}

func bestValidRoadRun(runs []BenchmarkRun) *BenchmarkRun {
	valid := make([]BenchmarkRun, 0, len(runs))
	for _, run := range runs {
		if run.Valid && run.DurationMS > 0 {
			valid = append(valid, run)
		}
	}
	if len(valid) == 0 {
		return nil
	}
	sort.SliceStable(valid, func(i, j int) bool {
		if valid[i].DurationMS != valid[j].DurationMS {
			return valid[i].DurationMS < valid[j].DurationMS
		}
		return valid[i].Confidence > valid[j].Confidence
	})
	return &valid[0]
}

func roadPaperScore(baseline *BenchmarkRun, fallback *BenchmarkRun) float64 {
	run := baseline
	if run == nil {
		run = fallback
	}
	if run == nil {
		return 0
	}
	score := 35 + roadClamp01(run.Confidence)*25
	if run.RouteProgress01 != nil {
		score += roadClamp01((*run.RouteProgress01-0.75)/0.25) * 20
	} else {
		score += 8
	}
	if run.MaxSpeedKmh != nil {
		score += roadClamp01(*run.MaxSpeedKmh/250) * 10
	}
	score -= math.Min(float64(run.EventCount)*2.5, 12)
	score -= math.Min(float64(len(splitWarningFlags(run.WarningFlags)))*4, 16)
	return roadClamp(score, 0, 100)
}

func roadRiskScore(run BenchmarkRun, events []telemetry.DetectedEvent) float64 {
	score := 0.0
	for _, event := range events {
		switch strings.TrimSpace(event.Severity) {
		case "high":
			score += 18
		case "medium":
			score += 10
		default:
			score += 5
		}
		if event.DurationMS > 1000 {
			score += 4
		}
	}
	for _, flag := range splitWarningFlags(run.WarningFlags) {
		switch flag {
		case "route_deviation":
			score += 20
		case "route_progress_low", "geometry_length_mismatch":
			score += 15
		default:
			score += 8
		}
	}
	if run.AvgLateralErrorMeters != nil {
		score += roadClamp(*run.AvgLateralErrorMeters/60, 0, 1) * 16
	}
	if run.MaxLateralErrorMeters != nil {
		score += roadClamp(*run.MaxLateralErrorMeters/120, 0, 1) * 12
	}
	return roadClamp(score, 0, 100)
}

func roadPlayerFitScore(session TelemetrySession, best BenchmarkRun, baseline *BenchmarkRun, risk float64) float64 {
	driverMode := NormalizeTestConditions(TestConditions{DriverMode: session.DriverMode}).DriverMode
	if driverMode != "player" || baseline == nil || baseline.DurationMS <= 0 || best.DurationMS <= 0 {
		return 0
	}
	timeGainPct := (float64(baseline.DurationMS-best.DurationMS) / float64(baseline.DurationMS)) * 100
	score := 55 + roadClamp(timeGainPct, -15, 15)*2
	score += (100 - risk) * 0.25
	if best.RouteProgress01 != nil {
		score += roadClamp((*best.RouteProgress01-0.85)/0.15, 0, 1) * 8
	}
	return roadClamp(score, 0, 100)
}

func roadVerdict(session TelemetrySession, best BenchmarkRun, baseline *BenchmarkRun, risk float64) string {
	driverMode := NormalizeTestConditions(TestConditions{DriverMode: session.DriverMode}).DriverMode
	if best.DurationMS <= 0 || driverMode == "unknown" || driverMode == "auto" || baseline == nil || baseline.DurationMS <= 0 {
		return roadVerdictInsufficientData
	}
	playerFaster := float64(best.DurationMS) <= float64(baseline.DurationMS)*1.02
	playerClearlySlower := float64(best.DurationMS) > float64(baseline.DurationMS)*1.05
	if playerFaster && risk <= 35 {
		return roadVerdictGoodFit
	}
	if float64(best.DurationMS) < float64(baseline.DurationMS) && risk > 35 {
		return roadVerdictFastButRisky
	}
	if playerClearlySlower && risk >= 30 {
		return roadVerdictPaperFastNotFit
	}
	return roadVerdictNeedsTuning
}

func roadPlayerIsFaster(best BenchmarkRun, baseline *BenchmarkRun) bool {
	return baseline != nil && baseline.DurationMS > 0 && best.DurationMS > 0 && best.DurationMS < baseline.DurationMS
}

func roadAttributions(events []telemetry.DetectedEvent, run BenchmarkRun, baselineStatus string, risk float64, playerFaster bool) []RoadEvaluationAttribution {
	type aggregate struct {
		count    int
		severity string
	}
	eventsByType := map[string]aggregate{}
	for _, event := range events {
		eventType := strings.TrimSpace(event.Type)
		if eventType == "" {
			continue
		}
		agg := eventsByType[eventType]
		agg.count++
		agg.severity = maxSeverity(agg.severity, event.Severity)
		eventsByType[eventType] = agg
	}
	attributions := make([]RoadEvaluationAttribution, 0, len(eventsByType)+4)
	for eventType, agg := range eventsByType {
		attrType, prioritize := roadEventAttribution(eventType, risk, playerFaster)
		attributions = append(attributions, RoadEvaluationAttribution{
			Type:             attrType,
			EventType:        eventType,
			Count:            agg.count,
			Severity:         agg.severity,
			Priority:         roadAttributionPriority(attrType, agg.severity, prioritize),
			Message:          "event_pattern",
			PrioritizeTuning: prioritize,
		})
	}
	for _, flag := range splitWarningFlags(run.WarningFlags) {
		attributions = append(attributions, RoadEvaluationAttribution{
			Type:             roadAttributionDriverExecution,
			EventType:        flag,
			Count:            1,
			Priority:         3,
			Message:          flag,
			PrioritizeTuning: false,
		})
	}
	switch baselineStatus {
	case roadBaselineMissingAuto, roadBaselineMissingVehicle, roadBaselineNoValidStandardRun, roadBaselineNoStandardTrack:
		attributions = append(attributions, RoadEvaluationAttribution{
			Type:             roadAttributionDataGap,
			Count:            1,
			Priority:         1,
			Message:          baselineStatus,
			PrioritizeTuning: false,
		})
	}
	sort.SliceStable(attributions, func(i, j int) bool {
		if attributions[i].Priority != attributions[j].Priority {
			return attributions[i].Priority < attributions[j].Priority
		}
		if attributions[i].Count != attributions[j].Count {
			return attributions[i].Count > attributions[j].Count
		}
		return attributions[i].EventType < attributions[j].EventType
	})
	return attributions
}

func roadEventAttribution(eventType string, risk float64, playerFaster bool) (string, bool) {
	switch eventType {
	case "front_brake_lockup", "rear_brake_lockup", "suspension_bottom_out", "short_gear", "long_gear_bog_down", "top_speed_limited_by_gearing", "launch_bog_down", "launch_wheelspin":
		return roadAttributionTuneIssue, true
	case "corner_entry_understeer", "corner_exit_oversteer", "high_speed_four_wheel_slide":
		prioritize := risk >= 55 || !playerFaster
		return roadAttributionStyleFitIssue, prioritize
	default:
		return roadAttributionDriverExecution, false
	}
}

func roadAttributionPriority(attrType string, severity string, prioritize bool) int {
	if attrType == roadAttributionDataGap {
		return 1
	}
	if prioritize && severity == "high" {
		return 1
	}
	if prioritize {
		return 2
	}
	if attrType == roadAttributionDriverExecution {
		return 4
	}
	return 3
}

func roadEvaluationNotes(e *RoadSessionEvaluation) []string {
	notes := make([]string, 0, 4)
	if e.BaselineStatus == roadBaselineMissingAuto {
		notes = append(notes, roadBaselineMissingAuto)
	}
	if e.BaselineStatus == roadBaselineMissingVehicle {
		notes = append(notes, roadBaselineMissingVehicle)
	}
	if e.OverallVerdict == roadVerdictInsufficientData {
		notes = append(notes, roadVerdictInsufficientData)
	}
	if e.BestRun != nil && e.BestRun.WarningFlags != "" {
		notes = append(notes, "benchmark_run_has_warnings")
	}
	return notes
}

func sameTuneProfile(left *int64, right *int64) bool {
	return left != nil && right != nil && *left == *right
}

func splitWarningFlags(value string) []string {
	parts := strings.Split(value, ",")
	flags := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			flags = append(flags, part)
		}
	}
	return flags
}

func maxSeverity(left string, right string) string {
	if severityRank(right) > severityRank(left) {
		return right
	}
	return left
}

func severityRank(value string) int {
	switch strings.TrimSpace(value) {
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

func roadEvalMetric(key, label, unit string, left, right float64, higherIsBetter bool) SessionComparisonMetric {
	return SessionComparisonMetric{Key: key, Label: label, Unit: unit, Left: left, Right: right, Delta: right - left, HigherIsBetter: higherIsBetter}
}

func roadClamp01(value float64) float64 {
	return roadClamp(value, 0, 1)
}

func roadClamp(value float64, min float64, max float64) float64 {
	if math.IsNaN(value) || math.IsInf(value, 0) {
		return min
	}
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
