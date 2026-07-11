package storage

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"

	"fh6worker/internal/telemetry"

	"github.com/xuri/excelize/v2"
)

const (
	roadDecisionReady          = "ready"
	roadDecisionNoSymptom      = "no_matching_symptom"
	roadDecisionProfileUnbound = "profile_unbound"
	roadDecisionRollback       = "rollback_recommended"
	roadDecisionInsufficient   = "insufficient_data"
	roadDecisionKnowledgeError = "knowledge_error"
)

type RoadTuningKnowledge struct {
	Symptoms []RoadSymptomCard
	Actions  []RoadActionCard
}

func (s *Store) ReloadTuningKnowledge() error {
	path := defaultRoadTuningKnowledgePath()
	knowledge, err := LoadRoadTuningKnowledge(path)
	status := RoadTuningKnowledgeStatus{
		LoadedAt:   nowText(),
		SourcePath: path,
	}
	if err != nil {
		status.LastError = err.Error()
		status.UsingFallback = true
		knowledge = fallbackRoadTuningKnowledge()
	}
	status.SymptomCount = len(knowledge.Symptoms)
	status.ActionCount = len(knowledge.Actions)

	s.knowledgeMu.Lock()
	s.roadKnowledge = knowledge
	s.knowledgeStatus = status
	s.knowledgeMu.Unlock()
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) RoadTuningKnowledgeStatus() RoadTuningKnowledgeStatus {
	s.knowledgeMu.RLock()
	defer s.knowledgeMu.RUnlock()
	return s.knowledgeStatus
}

func (s *Store) roadTuningKnowledge() (*RoadTuningKnowledge, RoadTuningKnowledgeStatus) {
	s.knowledgeMu.RLock()
	knowledge := s.roadKnowledge
	status := s.knowledgeStatus
	s.knowledgeMu.RUnlock()
	if knowledge != nil {
		return knowledge, status
	}
	_ = s.ReloadTuningKnowledge()
	s.knowledgeMu.RLock()
	defer s.knowledgeMu.RUnlock()
	return s.roadKnowledge, s.knowledgeStatus
}

func defaultRoadTuningKnowledgePath() string {
	candidates := []string{
		filepath.Join("docs", "FH6Data.xlsx"),
		filepath.Join("..", "docs", "FH6Data.xlsx"),
	}
	if _, file, _, ok := runtime.Caller(0); ok {
		candidates = append(candidates, filepath.Clean(filepath.Join(filepath.Dir(file), "..", "..", "..", "docs", "FH6Data.xlsx")))
	}
	for _, candidate := range candidates {
		if _, err := filepath.Abs(candidate); err == nil {
			if fileExists(candidate) {
				return candidate
			}
		}
	}
	return filepath.Join("docs", "FH6Data.xlsx")
}

func fileExists(path string) bool {
	if strings.TrimSpace(path) == "" {
		return false
	}
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func LoadRoadTuningKnowledge(path string) (*RoadTuningKnowledge, error) {
	if strings.TrimSpace(path) == "" {
		return nil, errors.New("road tuning knowledge path is required")
	}
	file, err := excelize.OpenFile(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	symptoms, err := readRoadSymptomCards(file)
	if err != nil {
		return nil, err
	}
	actions, err := readRoadActionCards(file)
	if err != nil {
		return nil, err
	}
	if len(symptoms) == 0 {
		return nil, errors.New("公路症状卡 has no enabled symptoms")
	}
	if len(actions) == 0 {
		return nil, errors.New("公路建议动作 has no actions")
	}
	return &RoadTuningKnowledge{Symptoms: symptoms, Actions: actions}, nil
}

func readRoadSymptomCards(file *excelize.File) ([]RoadSymptomCard, error) {
	rows, err := sheetRows(file, "公路症状卡")
	if err != nil {
		return nil, err
	}
	required := []string{"症状ID", "阶段", "症状", "主因", "适用事件", "优先级", "启用"}
	indices, err := headerIndex(rows, required)
	if err != nil {
		return nil, fmt.Errorf("公路症状卡: %w", err)
	}
	out := []RoadSymptomCard{}
	for _, row := range rows[1:] {
		id := cell(row, indices["症状ID"])
		if id == "" {
			continue
		}
		enabled := parseBool(cell(row, indices["启用"]))
		if !enabled {
			continue
		}
		priority := parseInt(cell(row, indices["优先级"]))
		sourcePriority := priority
		if index, ok := indices["来源优先级"]; ok {
			if parsed := parseInt(cell(row, index)); parsed > 0 {
				sourcePriority = parsed
			}
		}
		out = append(out, RoadSymptomCard{
			ID:             id,
			Phase:          cell(row, indices["阶段"]),
			Symptom:        cell(row, indices["症状"]),
			PrimaryCause:   cell(row, indices["主因"]),
			EventTypes:     splitCSV(cell(row, indices["适用事件"])),
			EvidenceRule:   cell(row, optionalIndex(indices, "证据条件")),
			Priority:       priority,
			Source:         cell(row, optionalIndex(indices, "来源")),
			Enabled:        true,
			SourcePriority: sourcePriority,
		})
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Priority != out[j].Priority {
			return out[i].Priority > out[j].Priority
		}
		return out[i].ID < out[j].ID
	})
	return out, nil
}

func readRoadActionCards(file *excelize.File) ([]RoadActionCard, error) {
	rows, err := sheetRows(file, "公路建议动作")
	if err != nil {
		return nil, err
	}
	required := []string{"症状ID", "排序", "类别", "调校项", "字段Key", "方向", "起始幅度", "单位", "传动限制", "可自动应用", "说明", "来源", "来源优先级"}
	indices, err := headerIndex(rows, required)
	if err != nil {
		return nil, fmt.Errorf("公路建议动作: %w", err)
	}
	out := []RoadActionCard{}
	for _, row := range rows[1:] {
		symptomID := cell(row, indices["症状ID"])
		item := cell(row, indices["调校项"])
		if symptomID == "" || item == "" {
			continue
		}
		out = append(out, RoadActionCard{
			SymptomID:             symptomID,
			Rank:                  parseInt(cell(row, indices["排序"])),
			Category:              cell(row, indices["类别"]),
			Item:                  item,
			FieldKey:              cell(row, indices["字段Key"]),
			Direction:             strings.ToLower(cell(row, indices["方向"])),
			Amount:                parseFloat(cell(row, indices["起始幅度"])),
			Unit:                  cell(row, indices["单位"]),
			DrivetrainRestriction: cell(row, indices["传动限制"]),
			CanAutoApply:          parseBool(cell(row, indices["可自动应用"])),
			Description:           cell(row, indices["说明"]),
			Source:                cell(row, indices["来源"]),
			SourcePriority:        parseInt(cell(row, indices["来源优先级"])),
		})
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].SymptomID != out[j].SymptomID {
			return out[i].SymptomID < out[j].SymptomID
		}
		if out[i].Rank != out[j].Rank {
			return out[i].Rank < out[j].Rank
		}
		return out[i].SourcePriority > out[j].SourcePriority
	})
	return out, nil
}

func sheetRows(file *excelize.File, sheetName string) ([][]string, error) {
	rows, err := file.GetRows(sheetName)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, fmt.Errorf("%s is empty", sheetName)
	}
	return rows, nil
}

func headerIndex(rows [][]string, required []string) (map[string]int, error) {
	if len(rows) == 0 {
		return nil, errors.New("sheet is empty")
	}
	indices := map[string]int{}
	for index, name := range rows[0] {
		indices[strings.TrimSpace(name)] = index
	}
	missing := []string{}
	for _, name := range required {
		if _, ok := indices[name]; !ok {
			missing = append(missing, name)
		}
	}
	if len(missing) > 0 {
		return nil, fmt.Errorf("missing required columns: %s", strings.Join(missing, ", "))
	}
	return indices, nil
}

func optionalIndex(indices map[string]int, key string) int {
	if index, ok := indices[key]; ok {
		return index
	}
	return -1
}

func cell(row []string, index int) string {
	if index < 0 || index >= len(row) {
		return ""
	}
	return strings.TrimSpace(row[index])
}

func splitCSV(value string) []string {
	out := []string{}
	for _, part := range strings.Split(value, ",") {
		part = strings.TrimSpace(part)
		if part != "" {
			out = append(out, part)
		}
	}
	return out
}

func parseBool(value string) bool {
	value = strings.ToLower(strings.TrimSpace(value))
	return value == "true" || value == "1" || value == "yes" || value == "y" || value == "启用"
}

func parseInt(value string) int {
	parsed, _ := strconv.Atoi(strings.TrimSpace(value))
	return parsed
}

func parseFloat(value string) float64 {
	parsed, _ := strconv.ParseFloat(strings.TrimSpace(value), 64)
	return parsed
}

func (s *Store) GetRoadTuningDecision(sessionID int64) (*RoadTuningDecision, error) {
	session, err := s.GetTelemetrySession(sessionID)
	if err != nil {
		return nil, err
	}
	knowledge, status := s.roadTuningKnowledge()
	decision := &RoadTuningDecision{
		SessionID:       sessionID,
		Status:          roadDecisionInsufficient,
		Reason:          "insufficient_session_data",
		Evidence:        map[string]float64{},
		RetestFocus:     []string{"same_car", "same_track", "same_driver_mode"},
		KnowledgeStatus: status,
	}
	if knowledge == nil || len(knowledge.Symptoms) == 0 {
		decision.Status = roadDecisionKnowledgeError
		decision.Reason = "road_tuning_knowledge_unavailable"
		return decision, nil
	}
	summary, err := s.GetSessionIssueSummary(sessionID)
	if err != nil {
		return nil, err
	}
	if len(summary.Groups) == 0 && summary.GearPower.Status != "ok" {
		decision.Status = roadDecisionNoSymptom
		decision.Reason = "no_issue_group"
		return decision, nil
	}
	samples, _ := s.GetSessionTelemetrySamples(sessionID, 10000)
	profile := s.issueTuneProfile(*session)
	retest, _ := s.GetRetestEvaluation(sessionID)
	candidate := chooseRoadSymptomCandidate(*summary, knowledge.Symptoms, samples)
	if candidate.card.ID == "" {
		if gearDecision := roadGearDecision(*summary, status); gearDecision != nil {
			return gearDecision, nil
		}
		decision.Status = roadDecisionNoSymptom
		decision.Reason = "no_matching_symptom_card"
		return decision, nil
	}
	evidence := cloneEvidenceMap(candidate.evidence)
	if _, ok := evidence["gear"]; !ok {
		if gear := dominantGearForWindow(samples, candidate.group.FirstStartMS, candidate.group.LastEndMS); gear > 0 {
			evidence["gear"] = float64(gear)
		}
	}
	confidence := roadDecisionConfidence(candidate.group, candidate.score)
	decision.Status = roadDecisionReady
	decision.SymptomID = candidate.card.ID
	decision.Phase = candidate.card.Phase
	decision.Symptom = candidate.card.Symptom
	decision.PrimaryCause = candidate.card.PrimaryCause
	decision.Confidence = confidence
	decision.FitVerdict = roadFitVerdict(retest)
	decision.RelatedIssueGroup = &candidate.group
	decision.Evidence = evidence
	decision.Reason = "matched_symptom_card"
	if retest != nil && retest.Status == "worsened" && len(candidate.group.RelatedRecentChanges) > 0 {
		recentDeltas := s.recentChangeDeltasForSession(*session)
		rollback := rollbackActionsForGroup(candidate.group, recentDeltas)
		if len(rollback) > 0 {
			decision.Status = roadDecisionRollback
			decision.RollbackRecommended = true
			decision.Reason = "previous_run_worsened"
			decision.Actions = append(decision.Actions, decisionActionsFromSuggested(rollback, candidate.group.Family, evidence, "rollback", confidence)...)
		}
	}
	if len(decision.Actions) < 3 {
		actions := knowledge.actionsForSymptom(candidate.card.ID)
		decision.Actions = append(decision.Actions, roadDecisionActions(actions, candidate.group.Family, evidence, session.Drivetrain, profile, confidence, 3-len(decision.Actions))...)
	}
	decision.Actions = finalizeRoadDecisionActions(decision.Actions, 3)
	if len(decision.Actions) == 0 {
		decision.Status = roadDecisionNoSymptom
		decision.Reason = "no_usable_actions_for_symptom"
	}
	decision.RetestFocus = retestFocusForDecision(decision)
	return decision, nil
}

type roadSymptomCandidate struct {
	card     RoadSymptomCard
	group    SessionIssueGroup
	evidence map[string]float64
	score    float64
}

func chooseRoadSymptomCandidate(summary SessionIssueSummary, cards []RoadSymptomCard, samples []telemetry.NormalizedTelemetry) roadSymptomCandidate {
	best := roadSymptomCandidate{}
	for _, card := range cards {
		for _, group := range summary.Groups {
			evidence := roadIssueEvidence(group, samples)
			if !cardMatchesGroup(card, group, evidence) {
				continue
			}
			score := float64(card.Priority) + issueGroupScore(group)
			if group.PrioritizeTuning {
				score += 20
			}
			if group.Comparison == issueCompareWorsened {
				score += 15
			}
			if score > best.score {
				best = roadSymptomCandidate{card: card, group: group, evidence: evidence, score: score}
			}
		}
	}
	return best
}

func cardMatchesGroup(card RoadSymptomCard, group SessionIssueGroup, evidence map[string]float64) bool {
	eventSet := map[string]bool{}
	for _, eventType := range group.EventTypes {
		eventSet[eventType] = true
	}
	matchedEvent := false
	for _, eventType := range card.EventTypes {
		if eventSet[eventType] {
			matchedEvent = true
			break
		}
	}
	return matchedEvent && evidenceRuleMatches(card.EvidenceRule, evidence)
}

func roadIssueEvidence(group SessionIssueGroup, samples []telemetry.NormalizedTelemetry) map[string]float64 {
	out := representativeIssueEvidence(group)
	minSpeed, avgSpeed, maxSpeed, ok := speedStatsForIssueWindow(group, samples)
	if !ok {
		if stat, exists := group.Evidence["speed_kmh"]; exists && stat.Count > 0 {
			minSpeed, avgSpeed, maxSpeed, ok = stat.Min, stat.Avg, stat.Max, true
		}
	}
	if ok {
		out["speed_kmh"] = avgSpeed
		out["speed_min_kmh"] = minSpeed
		out["speed_avg_kmh"] = avgSpeed
		out["speed_max_kmh"] = maxSpeed
		out["speed_band"] = float64(roadSpeedBand(avgSpeed, maxSpeed))
	}
	return out
}

func speedStatsForIssueWindow(group SessionIssueGroup, samples []telemetry.NormalizedTelemetry) (float64, float64, float64, bool) {
	if len(samples) == 0 {
		return 0, 0, 0, false
	}
	startMS := group.FirstStartMS
	endMS := group.LastEndMS
	if endMS < startMS {
		endMS = startMS
	}
	count := 0
	sum := 0.0
	minSpeed := 0.0
	maxSpeed := 0.0
	for _, sample := range samples {
		if sample.TimeMS < startMS || sample.TimeMS > endMS {
			continue
		}
		if count == 0 || sample.SpeedKmh < minSpeed {
			minSpeed = sample.SpeedKmh
		}
		if count == 0 || sample.SpeedKmh > maxSpeed {
			maxSpeed = sample.SpeedKmh
		}
		sum += sample.SpeedKmh
		count++
	}
	if count == 0 {
		return 0, 0, 0, false
	}
	return minSpeed, sum / float64(count), maxSpeed, true
}

func roadSpeedBand(avgSpeed float64, maxSpeed float64) int {
	switch {
	case avgSpeed >= 160 || maxSpeed >= 175:
		return 3
	case avgSpeed >= 90:
		return 2
	default:
		return 1
	}
}

func evidenceRuleMatches(rule string, evidence map[string]float64) bool {
	for _, part := range strings.Split(rule, ";") {
		if !evidenceRuleClauseMatches(strings.TrimSpace(part), evidence) {
			return false
		}
	}
	return true
}

func evidenceRuleClauseMatches(clause string, evidence map[string]float64) bool {
	if clause == "" {
		return true
	}
	for _, op := range []string{">=", "<=", ">", "<"} {
		if index := strings.Index(clause, op); index >= 0 {
			key := strings.TrimSpace(clause[:index])
			target, err := strconv.ParseFloat(strings.TrimSpace(clause[index+len(op):]), 64)
			if err != nil {
				return true
			}
			value, ok := roadEvidenceValue(evidence, key)
			if !ok {
				return true
			}
			switch op {
			case ">=":
				return value >= target
			case "<=":
				return value <= target
			case ">":
				return value > target
			case "<":
				return value < target
			}
		}
	}
	return true
}

func roadEvidenceValue(evidence map[string]float64, key string) (float64, bool) {
	normalized := strings.ToLower(strings.TrimSpace(key))
	if value, ok := evidence[normalized]; ok {
		return value, true
	}
	switch normalized {
	case "speed":
		if value, ok := evidence["speed_avg_kmh"]; ok {
			return value, true
		}
		return evidence["speed_kmh"], evidence["speed_kmh"] != 0
	case "front_slip", "front_slip_high":
		value, ok := evidence["front_combined_slip"]
		return value, ok
	case "rear_slip", "rear_slip_high", "driven_slip_high":
		value, ok := evidence["rear_combined_slip"]
		return value, ok
	default:
		return 0, false
	}
}

func roadDecisionConfidence(group SessionIssueGroup, score float64) string {
	switch {
	case group.Severity == "high" || group.EventCount >= 3 || score >= 150:
		return "high"
	case group.EventCount >= 2 || group.TotalDurationMS >= 2500:
		return "medium"
	default:
		return "low"
	}
}

func roadFitVerdict(retest *RetestEvaluation) string {
	if retest == nil || retest.Status == "insufficient_data" {
		return "unknown"
	}
	return retest.Status
}

func (k *RoadTuningKnowledge) actionsForSymptom(symptomID string) []RoadActionCard {
	out := []RoadActionCard{}
	for _, action := range k.Actions {
		if action.SymptomID == symptomID {
			out = append(out, action)
		}
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Rank != out[j].Rank {
			return out[i].Rank < out[j].Rank
		}
		return out[i].SourcePriority > out[j].SourcePriority
	})
	return out
}

func roadDecisionActions(actions []RoadActionCard, family string, evidence map[string]float64, sessionDrivetrain string, profile *TuneProfile, confidence string, limit int) []RoadTuningDecisionAction {
	out := []RoadTuningDecisionAction{}
	for _, action := range actions {
		if limit > 0 && len(out) >= limit {
			break
		}
		if !actionMatchesDrivetrain(action.DrivetrainRestriction, sessionDrivetrain, profile) {
			continue
		}
		role := "alternative"
		if len(out) == 0 {
			role = "primary"
		} else if len(out) == 1 {
			role = "support"
		}
		canApply := action.CanAutoApply && action.Direction != "check"
		blocked := ""
		if !canApply {
			blocked = "manual_review_required"
		}
		out = append(out, RoadTuningDecisionAction{
			ID:            fmt.Sprintf("%s:%d:%s:%s", action.SymptomID, action.Rank, action.Item, action.Direction),
			Role:          role,
			Family:        family,
			Source:        action.Source,
			Confidence:    confidence,
			TrustLevel:    trustLevelFromConfidence(confidence),
			AdviceLayer:   adviceLayerForSource(action.Source, family),
			Category:      action.Category,
			Item:          action.Item,
			FieldKey:      action.FieldKey,
			Direction:     action.Direction,
			Amount:        formatRoadAmount(action.Amount),
			Unit:          action.Unit,
			Reason:        action.Description,
			Rationale:     adviceRationale(action.Source, action.Description),
			CanAutoApply:  canApply,
			BlockedReason: blocked,
			Evidence:      cloneEvidenceMap(evidence),
		})
	}
	return out
}

func decisionActionsFromSuggested(actions []telemetry.SuggestedAction, family string, evidence map[string]float64, source string, confidence string) []RoadTuningDecisionAction {
	out := make([]RoadTuningDecisionAction, 0, len(actions))
	for i, action := range actions {
		role := "alternative"
		if i == 0 {
			role = "primary"
		} else if i == 1 {
			role = "support"
		}
		out = append(out, RoadTuningDecisionAction{
			ID:           fmt.Sprintf("%s:%d:%s:%s", source, i, action.Item, action.Direction),
			Role:         role,
			Family:       family,
			Source:       source,
			Confidence:   confidence,
			TrustLevel:   trustLevelFromConfidence(confidence),
			AdviceLayer:  adviceLayerForSource(source, family),
			Category:     action.Category,
			Item:         action.Item,
			Direction:    action.Direction,
			Amount:       action.Amount,
			Reason:       action.Reason,
			Rationale:    adviceRationale(source, action.Reason),
			CanAutoApply: action.Direction != "check",
			Evidence:     cloneEvidenceMap(evidence),
		})
	}
	return out
}

func actionMatchesDrivetrain(restriction string, sessionDrivetrain string, profile *TuneProfile) bool {
	restriction = strings.TrimSpace(strings.ToUpper(restriction))
	if restriction == "" || restriction == "ANY" {
		return true
	}
	drivetrain := strings.ToUpper(strings.TrimSpace(sessionDrivetrain))
	if drivetrain == "" && profile != nil {
		drivetrain = strings.ToUpper(strings.TrimSpace(profile.Drivetrain))
	}
	if drivetrain == "" {
		return true
	}
	for _, part := range strings.FieldsFunc(restriction, func(r rune) bool { return r == '/' || r == ',' || r == ';' }) {
		if strings.TrimSpace(part) == drivetrain {
			return true
		}
	}
	return false
}

func resolveRoadActionConflicts(actions []RoadTuningDecisionAction) []RoadTuningDecisionAction {
	selected := map[string]RoadTuningDecisionAction{}
	order := []string{}
	for _, action := range actions {
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
		if existing.Direction != action.Direction {
			continue
		}
		if !existing.CanAutoApply && action.CanAutoApply {
			selected[key] = action
		}
	}
	out := []RoadTuningDecisionAction{}
	for _, key := range order {
		out = append(out, selected[key])
	}
	return out
}

func limitRoadDecisionActions(actions []RoadTuningDecisionAction, limit int) []RoadTuningDecisionAction {
	if len(actions) <= limit {
		return actions
	}
	return actions[:limit]
}

func retestFocusForDecision(decision *RoadTuningDecision) []string {
	if decision == nil {
		return nil
	}
	focus := []string{"same_car", "same_track", "same_driver_mode"}
	if decision.SymptomID != "" {
		focus = append(focus, decision.SymptomID)
	}
	if decision.RollbackRecommended {
		focus = append(focus, "verify_rollback_before_new_changes")
	}
	return focus
}

func dominantGearForWindow(samples []telemetry.NormalizedTelemetry, startMS int64, endMS int64) int {
	counts := map[int]int{}
	for _, sample := range samples {
		if sample.TimeMS < startMS || sample.TimeMS > endMS {
			continue
		}
		if sample.Gear < 1 || sample.Gear > 10 {
			continue
		}
		counts[sample.Gear]++
	}
	bestGear, bestCount := 0, 0
	for gear, count := range counts {
		if count > bestCount {
			bestGear, bestCount = gear, count
		}
	}
	return bestGear
}

func roadGearDecision(summary SessionIssueSummary, status RoadTuningKnowledgeStatus) *RoadTuningDecision {
	if len(summary.GearPower.RecommendedActions) == 0 {
		return nil
	}
	decision := &RoadTuningDecision{
		SessionID:       summary.SessionID,
		Status:          roadDecisionReady,
		SymptomID:       "road_gearing_power",
		Phase:           "power",
		Symptom:         "齿比动力窗口异常",
		PrimaryCause:    "发动机未稳定处于有效动力区",
		Confidence:      "medium",
		FitVerdict:      "unknown",
		Reason:          "gear_power_diagnostic",
		Evidence:        cloneEvidenceMap(summary.GearPower.Evidence),
		RetestFocus:     []string{"same_car", "same_track", "same_driver_mode", "gear_power_window"},
		KnowledgeStatus: status,
	}
	decision.Actions = decisionActionsFromSuggested(summary.GearPower.RecommendedActions, "gearing_acceleration", summary.GearPower.Evidence, "gear_power_diagnostic", "medium")
	decision.Actions = finalizeRoadDecisionActions(decision.Actions, 3)
	return decision
}

func formatRoadAmount(value float64) string {
	if value == 0 {
		return "0"
	}
	if value == float64(int64(value)) {
		return fmt.Sprintf("%.0f", value)
	}
	return strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.2f", value), "0"), ".")
}

func fallbackRoadTuningKnowledge() *RoadTuningKnowledge {
	symptoms := []RoadSymptomCard{
		{ID: "road_entry_understeer_low_speed", Phase: "corner_entry", Symptom: "low-speed corner-entry understeer", PrimaryCause: "low-speed front mechanical grip or decel stability is insufficient", EventTypes: []string{"corner_entry_understeer"}, EvidenceRule: "speed_band<2", Priority: 94, Enabled: true, Source: "fallback", SourcePriority: 50},
		{ID: "road_entry_understeer_mid_speed", Phase: "corner_entry", Symptom: "mid-speed corner-entry understeer", PrimaryCause: "mid-speed chassis lateral-load balance is insufficient", EventTypes: []string{"corner_entry_understeer"}, EvidenceRule: "speed_band>=2;speed_band<3", Priority: 93, Enabled: true, Source: "fallback", SourcePriority: 50},
		{ID: "road_entry_understeer_high_speed", Phase: "corner_entry", Symptom: "high-speed corner-entry understeer", PrimaryCause: "high-speed front aero grip or platform stability is insufficient", EventTypes: []string{"corner_entry_understeer"}, EvidenceRule: "speed_band>=3", Priority: 95, Enabled: true, Source: "fallback", SourcePriority: 50},
		{ID: "road_entry_understeer", Phase: "corner_entry", Symptom: "入弯转向不足", PrimaryCause: "前轴入弯响应不足", EventTypes: []string{"corner_entry_understeer"}, Priority: 90, Enabled: true, Source: "fallback", SourcePriority: 50},
		{ID: "road_exit_oversteer", Phase: "corner_exit", Symptom: "出弯给油甩尾", PrimaryCause: "后轴动力释放过冲", EventTypes: []string{"corner_exit_oversteer", "snap_oversteer"}, Priority: 90, Enabled: true, Source: "fallback", SourcePriority: 50},
		{ID: "road_bottom_out", Phase: "platform", Symptom: "悬挂触底", PrimaryCause: "平台支撑不足", EventTypes: []string{"suspension_bottom_out"}, Priority: 80, Enabled: true, Source: "fallback", SourcePriority: 50},
	}
	actions := []RoadActionCard{
		{SymptomID: "road_entry_understeer_low_speed", Rank: 1, Category: "antiroll", Item: "front_arb", FieldKey: "frontArb", Direction: "decrease", Amount: 0.5, CanAutoApply: true, Description: "soften front ARB for low-speed front grip", Source: "fallback", SourcePriority: 50},
		{SymptomID: "road_entry_understeer_low_speed", Rank: 2, Category: "antiroll", Item: "rear_arb", FieldKey: "rearArb", Direction: "increase", Amount: 0.4, CanAutoApply: true, Description: "add low-speed rotation", Source: "fallback", SourcePriority: 50},
		{SymptomID: "road_entry_understeer_low_speed", Rank: 3, Category: "differential", Item: "drive_diff_decel", FieldKey: "driveDiffDecel", Direction: "decrease", Amount: 2, Unit: "%", CanAutoApply: true, Description: "reduce entry push from decel locking", Source: "fallback", SourcePriority: 45},
		{SymptomID: "road_entry_understeer_mid_speed", Rank: 1, Category: "antiroll", Item: "front_arb", FieldKey: "frontArb", Direction: "decrease", Amount: 0.5, CanAutoApply: true, Description: "soften front ARB", Source: "fallback", SourcePriority: 50},
		{SymptomID: "road_entry_understeer_mid_speed", Rank: 2, Category: "antiroll", Item: "rear_arb", FieldKey: "rearArb", Direction: "increase", Amount: 0.5, CanAutoApply: true, Description: "help rotation in mid-speed entry", Source: "fallback", SourcePriority: 50},
		{SymptomID: "road_entry_understeer_mid_speed", Rank: 3, Category: "damping", Item: "front_rebound", FieldKey: "frontRebound", Direction: "decrease", Amount: 0.3, CanAutoApply: true, Description: "let the front load more smoothly", Source: "fallback", SourcePriority: 45},
		{SymptomID: "road_entry_understeer_high_speed", Rank: 1, Category: "aero", Item: "front_and_rear_aero", FieldKey: "frontAero", Direction: "increase", Amount: 1, Unit: "kgf", CanAutoApply: true, Description: "increase high-speed grip", Source: "fallback", SourcePriority: 50},
		{SymptomID: "road_entry_understeer_high_speed", Rank: 2, Category: "springs", Item: "ride_height", FieldKey: "rideHeight", Direction: "check", Amount: 0.1, Unit: "cm", CanAutoApply: false, Description: "verify platform and bottoming before chassis changes", Source: "fallback", SourcePriority: 45},
		{SymptomID: "road_entry_understeer_high_speed", Rank: 3, Category: "antiroll", Item: "front_arb", FieldKey: "frontArb", Direction: "decrease", Amount: 0.3, CanAutoApply: true, Description: "small front ARB change after platform check", Source: "fallback", SourcePriority: 40},
		{SymptomID: "road_entry_understeer", Rank: 1, Category: "antiroll", Item: "front_arb", FieldKey: "frontArb", Direction: "decrease", Amount: 0.6, CanAutoApply: true, Description: "soften front ARB", Source: "fallback", SourcePriority: 50},
		{SymptomID: "road_entry_understeer", Rank: 2, Category: "antiroll", Item: "rear_arb", FieldKey: "rearArb", Direction: "increase", Amount: 0.5, CanAutoApply: true, Description: "stiffen rear ARB", Source: "fallback", SourcePriority: 50},
		{SymptomID: "road_exit_oversteer", Rank: 1, Category: "differential", Item: "rear_diff_accel", FieldKey: "rearDiffAccel", Direction: "decrease", Amount: 3, Unit: "%", DrivetrainRestriction: "RWD/AWD", CanAutoApply: true, Description: "soften rear acceleration diff", Source: "fallback", SourcePriority: 50},
		{SymptomID: "road_bottom_out", Rank: 1, Category: "springs", Item: "ride_height", FieldKey: "rideHeight", Direction: "increase", Amount: 0.2, Unit: "cm", CanAutoApply: true, Description: "raise ride height", Source: "fallback", SourcePriority: 50},
	}
	return &RoadTuningKnowledge{Symptoms: symptoms, Actions: actions}
}
