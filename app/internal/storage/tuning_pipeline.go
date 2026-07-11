package storage

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"

	"fh6worker/internal/telemetry"
)

const (
	TuningSourceTireLabCurrent = "tire_lab_current"
	TuningSourceSession        = "telemetry_session"

	DetectorLegacyRoadEvents = "legacy_road_events_v1"
	DetectorTireLabProblems  = "tire_lab_problem_groups_v1"

	DecisionerLegacyRoad = "legacy_road_decision_v1"
	DecisionerTire       = "tire_problem_decision_v1"

	InterpreterLegacyAdvice = "legacy_advice_planner_v1"
	InterpreterTireRepair   = "tire_repair_explainer_v1"
	InterpreterRoadDocsV12  = "road_baseline_docs_v12_interpreter_v1"

	pipelineStatusReady        = "ready"
	pipelineStatusNoData       = "no_data"
	pipelineStatusIncompatible = "incompatible"

	professionalPipelineConfigKey = "professional_pipeline_config"
)

type ProfessionalPipelineConfig struct {
	DetectorID    string `json:"detectorId"`
	DecisionerID  string `json:"decisionerId"`
	InterpreterID string `json:"interpreterId"`
}

type ProfessionalTuningDiagnostic struct {
	Status    string                     `json:"status"`
	UpdatedAt string                     `json:"updatedAt"`
	Config    ProfessionalPipelineConfig `json:"config"`
	Pipeline  *TuningPipelineRunResult   `json:"pipeline,omitempty"`
	Warnings  []string                   `json:"warnings"`
}

type TuningPipelineCatalog struct {
	SourceTypes         []TuningPipelineComponent   `json:"sourceTypes"`
	Detectors           []TuningPipelineComponent   `json:"detectors"`
	Decisioners         []TuningPipelineComponent   `json:"decisioners"`
	Interpreters        []TuningPipelineComponent   `json:"interpreters"`
	DefaultCombinations []TuningPipelineCombination `json:"defaultCombinations"`
}

type TuningPipelineComponent struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	Description    string   `json:"description"`
	SourceTypes    []string `json:"sourceTypes"`
	CompatibleWith []string `json:"compatibleWith"`
	Tags           []string `json:"tags"`
}

type TuningPipelineCombination struct {
	SourceType    string `json:"sourceType"`
	DetectorID    string `json:"detectorId"`
	DecisionerID  string `json:"decisionerId"`
	InterpreterID string `json:"interpreterId"`
	Description   string `json:"description"`
}

type TuningPipelineRunInput struct {
	SourceType    string `json:"sourceType"`
	SessionID     int64  `json:"sessionId,omitempty"`
	DetectorID    string `json:"detectorId"`
	DecisionerID  string `json:"decisionerId"`
	InterpreterID string `json:"interpreterId"`
}

type TuningPipelineRunResult struct {
	Status        string                      `json:"status"`
	UpdatedAt     string                      `json:"updatedAt"`
	SourceSummary TuningPipelineSourceSummary `json:"sourceSummary"`
	ProblemSet    TuningProblemSet            `json:"problemSet"`
	DecisionSet   TuningDecisionSet           `json:"decisionSet"`
	AdviceSet     TuningAdviceSet             `json:"adviceSet"`
	Warnings      []string                    `json:"warnings"`
}

type TuningPipelineSourceSummary struct {
	SourceType  string                 `json:"sourceType"`
	SessionID   int64                  `json:"sessionId,omitempty"`
	SampleCount int                    `json:"sampleCount"`
	EventCount  int                    `json:"eventCount"`
	Vehicle     SessionVehicleSnapshot `json:"vehicle"`
	GameMode    string                 `json:"gameMode"`
	DriverMode  string                 `json:"driverMode"`
	Label       string                 `json:"label"`
}

type TuningProblemSet struct {
	DetectorID string          `json:"detectorId"`
	Status     string          `json:"status"`
	Problems   []TuningProblem `json:"problems"`
	Warnings   []string        `json:"warnings"`
}

type TuningProblem struct {
	ID            string             `json:"id"`
	SourceID      string             `json:"sourceId"`
	Family        string             `json:"family"`
	Type          string             `json:"type"`
	Phase         string             `json:"phase"`
	OperationTags []string           `json:"operationTags"`
	LimitedAxle   string             `json:"limitedAxle"`
	LimitedWheels []string           `json:"limitedWheels"`
	Severity      string             `json:"severity"`
	Confidence    string             `json:"confidence"`
	RiskLevel     string             `json:"riskLevel"`
	Count         int                `json:"count"`
	DurationMS    int64              `json:"durationMs"`
	Summary       string             `json:"summary"`
	Reason        string             `json:"reason"`
	Evidence      map[string]float64 `json:"evidence"`
}

type TuningDecisionSet struct {
	DecisionerID string           `json:"decisionerId"`
	Status       string           `json:"status"`
	Decisions    []TuningDecision `json:"decisions"`
	Warnings     []string         `json:"warnings"`
}

type TuningDecision struct {
	ID              string             `json:"id"`
	ProblemID       string             `json:"problemId"`
	Phase           string             `json:"phase"`
	PrimaryCause    string             `json:"primaryCause"`
	ShouldTune      bool               `json:"shouldTune"`
	Confidence      string             `json:"confidence"`
	Rationale       string             `json:"rationale"`
	DocumentContext string             `json:"documentContext"`
	Evidence        map[string]float64 `json:"evidence"`
}

type TuningAdviceSet struct {
	InterpreterID   string         `json:"interpreterId"`
	Status          string         `json:"status"`
	Advice          []TuningAdvice `json:"advice"`
	DocumentSources []string       `json:"documentSources"`
	Warnings        []string       `json:"warnings"`
}

type TuningAdvice struct {
	ID              string             `json:"id"`
	DecisionID      string             `json:"decisionId"`
	ProblemID       string             `json:"problemId"`
	Layer           string             `json:"layer"`
	Category        string             `json:"category"`
	Scope           string             `json:"scope"`
	Direction       string             `json:"direction"`
	RelatedFields   []string           `json:"relatedFields"`
	Rationale       string             `json:"rationale"`
	VerifyEvidence  []string           `json:"verifyEvidence"`
	TrustLevel      string             `json:"trustLevel"`
	MissingInputs   []string           `json:"missingInputs"`
	ConflictReason  string             `json:"conflictReason"`
	CanApply        bool               `json:"canApply"`
	DocumentSources []string           `json:"documentSources"`
	Evidence        map[string]float64 `json:"evidence"`
}

type tuningPipelineContext struct {
	input        TuningPipelineRunInput
	samples      []telemetry.NormalizedTelemetry
	session      *TelemetrySession
	summary      *SessionIssueSummary
	roadDecision *RoadTuningDecision
	tireAnalysis TireIssueAnalysis
	tireAdvice   TireIssueAdvice
}

func ListTuningModelPipelines() TuningPipelineCatalog {
	return TuningPipelineCatalog{
		SourceTypes: []TuningPipelineComponent{
			{
				ID:          TuningSourceTireLabCurrent,
				Name:        "Tire Lab current window",
				Description: "Use the current in-memory Tire Lab telemetry window.",
				Tags:        []string{"memory", "experimental"},
			},
			{
				ID:          TuningSourceSession,
				Name:        "Saved telemetry session",
				Description: "Use a saved expert telemetry session and persisted detected events.",
				Tags:        []string{"session", "read_only"},
			},
		},
		Detectors: []TuningPipelineComponent{
			{
				ID:          DetectorLegacyRoadEvents,
				Name:        "Legacy road events v1",
				Description: "Convert saved road-rule detected events into a unified problem set.",
				SourceTypes: []string{TuningSourceSession},
				Tags:        []string{"road", "legacy"},
			},
			{
				ID:          DetectorTireLabProblems,
				Name:        "Tire Lab problem groups v1",
				Description: "Convert Tire Lab issue groups into a unified problem set.",
				SourceTypes: []string{TuningSourceTireLabCurrent},
				Tags:        []string{"tire", "experimental"},
			},
		},
		Decisioners: []TuningPipelineComponent{
			{
				ID:             DecisionerLegacyRoad,
				Name:           "Legacy road decision v1",
				Description:    "Use the existing road tuning decision model for saved sessions.",
				CompatibleWith: []string{DetectorLegacyRoadEvents},
				Tags:           []string{"road", "legacy"},
			},
			{
				ID:             DecisionerTire,
				Name:           "Tire problem decision v1",
				Description:    "Translate Tire Lab problem groups into tune/no-tune decisions.",
				CompatibleWith: []string{DetectorTireLabProblems},
				Tags:           []string{"tire", "experimental"},
			},
		},
		Interpreters: []TuningPipelineComponent{
			{
				ID:             InterpreterLegacyAdvice,
				Name:           "Legacy advice planner v1",
				Description:    "Adapt current road decision actions or whole-car plan actions.",
				CompatibleWith: []string{DecisionerLegacyRoad},
				Tags:           []string{"road", "legacy"},
			},
			{
				ID:             InterpreterTireRepair,
				Name:           "Tire repair explainer v1",
				Description:    "Explain Tire Lab issue repair directions without numeric write values.",
				CompatibleWith: []string{DecisionerTire},
				Tags:           []string{"tire", "experimental"},
			},
			{
				ID:             InterpreterRoadDocsV12,
				Name:           "Road baseline docs v1.2 interpreter",
				Description:    "Explain directions using the road static baseline docs v1.0-v1.2 vocabulary.",
				CompatibleWith: []string{DecisionerLegacyRoad, DecisionerTire},
				Tags:           []string{"road", "docs_v1_2", "explain_only"},
			},
		},
		DefaultCombinations: []TuningPipelineCombination{
			{
				SourceType:    TuningSourceTireLabCurrent,
				DetectorID:    DetectorTireLabProblems,
				DecisionerID:  DecisionerTire,
				InterpreterID: InterpreterRoadDocsV12,
				Description:   "Experimental Tire Lab problem analysis explained through the Docs v1.2 road baseline vocabulary.",
			},
			{
				SourceType:    TuningSourceSession,
				DetectorID:    DetectorLegacyRoadEvents,
				DecisionerID:  DecisionerLegacyRoad,
				InterpreterID: InterpreterRoadDocsV12,
				Description:   "Saved session road events explained through the Docs v1.2 road baseline vocabulary.",
			},
		},
	}
}

func DefaultProfessionalPipelineConfig() ProfessionalPipelineConfig {
	return ProfessionalPipelineConfig{
		DetectorID:    DetectorTireLabProblems,
		DecisionerID:  DecisionerTire,
		InterpreterID: InterpreterRoadDocsV12,
	}
}

func NormalizeProfessionalPipelineConfig(input ProfessionalPipelineConfig) ProfessionalPipelineConfig {
	defaults := DefaultProfessionalPipelineConfig()
	if strings.TrimSpace(input.DetectorID) == "" {
		input.DetectorID = defaults.DetectorID
	}
	if strings.TrimSpace(input.DecisionerID) == "" {
		input.DecisionerID = defaults.DecisionerID
	}
	if strings.TrimSpace(input.InterpreterID) == "" {
		input.InterpreterID = defaults.InterpreterID
	}
	return input
}

func (s *Store) GetProfessionalPipelineConfig() (ProfessionalPipelineConfig, error) {
	var value string
	err := s.db.QueryRow(`SELECT value FROM app_setting WHERE key = ?`, professionalPipelineConfigKey).Scan(&value)
	if errors.Is(err, sql.ErrNoRows) {
		return DefaultProfessionalPipelineConfig(), nil
	}
	if err != nil {
		return ProfessionalPipelineConfig{}, err
	}
	var config ProfessionalPipelineConfig
	if err := json.Unmarshal([]byte(value), &config); err != nil {
		return DefaultProfessionalPipelineConfig(), nil
	}
	return NormalizeProfessionalPipelineConfig(config), nil
}

func (s *Store) SaveProfessionalPipelineConfig(input ProfessionalPipelineConfig) (ProfessionalPipelineConfig, error) {
	normalized := NormalizeProfessionalPipelineConfig(input)
	if warnings := ValidateTuningPipelineCombination(TuningSourceTireLabCurrent, normalized.DetectorID, normalized.DecisionerID, normalized.InterpreterID); len(warnings) > 0 {
		// Incompatible combinations are allowed for experimentation, but unknown IDs are not.
		for _, warning := range warnings {
			if strings.HasPrefix(warning, "unknown_") {
				return ProfessionalPipelineConfig{}, errors.New(warning)
			}
		}
	}
	payload, err := json.Marshal(normalized)
	if err != nil {
		return ProfessionalPipelineConfig{}, err
	}
	_, err = s.db.Exec(`INSERT INTO app_setting(key, value) VALUES(?, ?)
		ON CONFLICT(key) DO UPDATE SET value = excluded.value`, professionalPipelineConfigKey, string(payload))
	return normalized, err
}

func ValidateTuningPipelineCombination(sourceType, detectorID, decisionerID, interpreterID string) []string {
	catalog := ListTuningModelPipelines()
	warnings := []string{}
	detector, ok := findPipelineComponent(catalog.Detectors, detectorID)
	if !ok {
		warnings = append(warnings, "unknown_detector:"+detectorID)
	}
	decisioner, ok := findPipelineComponent(catalog.Decisioners, decisionerID)
	if !ok {
		warnings = append(warnings, "unknown_decisioner:"+decisionerID)
	}
	interpreter, ok := findPipelineComponent(catalog.Interpreters, interpreterID)
	if !ok {
		warnings = append(warnings, "unknown_interpreter:"+interpreterID)
	}
	if detector.ID != "" && len(detector.SourceTypes) > 0 && !pipelineStringSliceContains(detector.SourceTypes, sourceType) {
		warnings = append(warnings, "incompatible_detector_source:"+detectorID+"_requires_"+strings.Join(detector.SourceTypes, "_or_"))
	}
	if decisioner.ID != "" && len(decisioner.CompatibleWith) > 0 && !pipelineStringSliceContains(decisioner.CompatibleWith, detectorID) {
		warnings = append(warnings, "incompatible_decisioner_detector:"+decisionerID+"_with_"+detectorID)
	}
	if interpreter.ID != "" && len(interpreter.CompatibleWith) > 0 && !pipelineStringSliceContains(interpreter.CompatibleWith, decisionerID) {
		warnings = append(warnings, "incompatible_interpreter_decisioner:"+interpreterID+"_with_"+decisionerID)
	}
	return warnings
}

func findPipelineComponent(components []TuningPipelineComponent, id string) (TuningPipelineComponent, bool) {
	for _, component := range components {
		if component.ID == id {
			return component, true
		}
	}
	return TuningPipelineComponent{}, false
}

func pipelineStringSliceContains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

func (s *Store) RunTuningModelPipeline(input TuningPipelineRunInput, tireLabSamples []telemetry.NormalizedTelemetry) (*TuningPipelineRunResult, error) {
	input = normalizePipelineInput(input)
	result := &TuningPipelineRunResult{
		Status:    pipelineStatusReady,
		UpdatedAt: nowText(),
		ProblemSet: TuningProblemSet{
			DetectorID: input.DetectorID,
			Status:     pipelineStatusNoData,
		},
		DecisionSet: TuningDecisionSet{
			DecisionerID: input.DecisionerID,
			Status:       pipelineStatusNoData,
		},
		AdviceSet: TuningAdviceSet{
			InterpreterID: input.InterpreterID,
			Status:        pipelineStatusNoData,
		},
	}
	ctx := &tuningPipelineContext{input: input}
	if err := s.preparePipelineSource(ctx, tireLabSamples, result); err != nil {
		return nil, err
	}
	problems, warnings := s.runPipelineDetector(ctx)
	result.ProblemSet = problems
	result.Warnings = append(result.Warnings, warnings...)

	decisions, warnings := s.runPipelineDecisioner(ctx, result.ProblemSet)
	result.DecisionSet = decisions
	result.Warnings = append(result.Warnings, warnings...)

	advice, warnings := s.runPipelineInterpreter(ctx, result.DecisionSet)
	result.AdviceSet = advice
	result.Warnings = append(result.Warnings, warnings...)

	result.Warnings = append(result.Warnings, result.ProblemSet.Warnings...)
	result.Warnings = append(result.Warnings, result.DecisionSet.Warnings...)
	result.Warnings = append(result.Warnings, result.AdviceSet.Warnings...)
	if len(result.ProblemSet.Problems) == 0 && len(result.DecisionSet.Decisions) == 0 && len(result.AdviceSet.Advice) == 0 {
		result.Status = pipelineStatusNoData
	}
	for _, warning := range result.Warnings {
		if strings.Contains(warning, "incompatible") {
			result.Status = pipelineStatusIncompatible
			break
		}
	}
	return result, nil
}

func normalizePipelineInput(input TuningPipelineRunInput) TuningPipelineRunInput {
	if input.SourceType == "" {
		input.SourceType = TuningSourceTireLabCurrent
	}
	if input.DetectorID == "" {
		if input.SourceType == TuningSourceSession {
			input.DetectorID = DetectorLegacyRoadEvents
		} else {
			input.DetectorID = DetectorTireLabProblems
		}
	}
	if input.DecisionerID == "" {
		if input.DetectorID == DetectorLegacyRoadEvents {
			input.DecisionerID = DecisionerLegacyRoad
		} else {
			input.DecisionerID = DecisionerTire
		}
	}
	if input.InterpreterID == "" {
		input.InterpreterID = InterpreterRoadDocsV12
	}
	return input
}

func (s *Store) preparePipelineSource(ctx *tuningPipelineContext, tireLabSamples []telemetry.NormalizedTelemetry, result *TuningPipelineRunResult) error {
	switch ctx.input.SourceType {
	case TuningSourceTireLabCurrent:
		ctx.samples = append([]telemetry.NormalizedTelemetry(nil), tireLabSamples...)
		result.SourceSummary = summarizePipelineSamples(TuningSourceTireLabCurrent, ctx.samples, 0)
		return nil
	case TuningSourceSession:
		if ctx.input.SessionID <= 0 {
			return errors.New("sessionId is required for telemetry_session source")
		}
		session, err := s.GetTelemetrySession(ctx.input.SessionID)
		if err != nil {
			return err
		}
		ctx.session = session
		samples, err := s.GetSessionTelemetrySamples(ctx.input.SessionID, 20000)
		if err != nil {
			return err
		}
		ctx.samples = samples
		result.SourceSummary = summarizePipelineSession(*session, len(samples))
		return nil
	default:
		return fmt.Errorf("unknown pipeline sourceType %q", ctx.input.SourceType)
	}
}

func summarizePipelineSamples(source string, samples []telemetry.NormalizedTelemetry, sessionID int64) TuningPipelineSourceSummary {
	summary := TuningPipelineSourceSummary{
		SourceType:  source,
		SessionID:   sessionID,
		SampleCount: len(samples),
		Label:       source,
	}
	if len(samples) == 0 {
		return summary
	}
	current := samples[len(samples)-1]
	ordinal := int64(current.CarOrdinal)
	pi := int64(current.CarPI)
	cylinders := int64(current.NumCylinders)
	summary.Vehicle = SessionVehicleSnapshot{
		CarOrdinal:   &ordinal,
		CarClass:     current.CarClass,
		CarPI:        &pi,
		Drivetrain:   current.Drivetrain,
		NumCylinders: &cylinders,
	}
	summary.GameMode = current.GameMode
	summary.DriverMode = "current_window"
	return summary
}

func summarizePipelineSession(session TelemetrySession, sampleCount int) TuningPipelineSourceSummary {
	return TuningPipelineSourceSummary{
		SourceType:  TuningSourceSession,
		SessionID:   session.ID,
		SampleCount: sampleCount,
		EventCount:  int(session.EventCount),
		GameMode:    session.GameMode,
		DriverMode:  session.DriverMode,
		Label:       session.SessionName,
		Vehicle: SessionVehicleSnapshot{
			CarOrdinal:   session.CarOrdinal,
			CarClass:     session.CarClass,
			CarPI:        session.CarPI,
			Drivetrain:   session.Drivetrain,
			NumCylinders: session.NumCylinders,
		},
	}
}

func (s *Store) runPipelineDetector(ctx *tuningPipelineContext) (TuningProblemSet, []string) {
	out := TuningProblemSet{
		DetectorID: ctx.input.DetectorID,
		Status:     pipelineStatusNoData,
	}
	switch ctx.input.DetectorID {
	case DetectorTireLabProblems:
		if ctx.input.SourceType != TuningSourceTireLabCurrent {
			out.Warnings = append(out.Warnings, "incompatible_detector_source:tire_lab_problem_groups_v1_requires_tire_lab_current")
			return out, nil
		}
		ctx.tireAnalysis = BuildTireIssueAnalysis(ctx.samples)
		out.Problems = problemsFromTireAnalysis(ctx.tireAnalysis)
	case DetectorLegacyRoadEvents:
		if ctx.input.SourceType != TuningSourceSession {
			out.Warnings = append(out.Warnings, "incompatible_detector_source:legacy_road_events_v1_requires_telemetry_session")
			return out, nil
		}
		summary, err := s.GetSessionIssueSummary(ctx.input.SessionID)
		if err != nil {
			out.Warnings = append(out.Warnings, "legacy_road_events_error:"+err.Error())
			return out, nil
		}
		ctx.summary = summary
		out.Problems = problemsFromSessionSummary(*summary)
	default:
		out.Warnings = append(out.Warnings, "unknown_detector:"+ctx.input.DetectorID)
	}
	if len(out.Problems) > 0 {
		out.Status = pipelineStatusReady
	}
	return out, nil
}

func problemsFromTireAnalysis(analysis TireIssueAnalysis) []TuningProblem {
	problems := make([]TuningProblem, 0, len(analysis.Groups))
	for _, group := range analysis.Groups {
		problems = append(problems, TuningProblem{
			ID:            "tire:" + group.ID,
			SourceID:      group.ID,
			Family:        group.Type,
			Type:          group.Type,
			Phase:         group.Phase,
			OperationTags: append([]string(nil), group.OperationTags...),
			LimitedAxle:   group.LimitedAxle,
			LimitedWheels: append([]string(nil), group.LimitedWheels...),
			Severity:      group.RiskLevel,
			Confidence:    group.Confidence,
			RiskLevel:     group.RiskLevel,
			Count:         group.Count,
			DurationMS:    group.TotalDurationMS,
			Summary:       group.Type,
			Reason:        group.Reason,
			Evidence:      cloneFloatMap(group.RepresentativeEvidence),
		})
	}
	return problems
}

func problemsFromSessionSummary(summary SessionIssueSummary) []TuningProblem {
	problems := make([]TuningProblem, 0, len(summary.Groups))
	for _, group := range summary.Groups {
		problems = append(problems, TuningProblem{
			ID:         "legacy:" + group.ID,
			SourceID:   group.ID,
			Family:     group.Family,
			Type:       strings.Join(group.EventTypes, ","),
			Phase:      group.Segment,
			Severity:   group.Severity,
			Confidence: quickConfidenceMedium,
			RiskLevel:  group.Severity,
			Count:      group.EventCount,
			DurationMS: group.TotalDurationMS,
			Summary:    group.Family,
			Reason:     group.AdjustmentStrategy,
			Evidence:   issueEvidenceAverageMap(group.Evidence),
		})
	}
	return problems
}

func issueEvidenceAverageMap(input map[string]IssueEvidence) map[string]float64 {
	out := map[string]float64{}
	for key, stat := range input {
		out[key] = stat.Avg
		out[key+"_min"] = stat.Min
		out[key+"_max"] = stat.Max
	}
	return out
}

func (s *Store) runPipelineDecisioner(ctx *tuningPipelineContext, problems TuningProblemSet) (TuningDecisionSet, []string) {
	out := TuningDecisionSet{
		DecisionerID: ctx.input.DecisionerID,
		Status:       pipelineStatusNoData,
	}
	if len(problems.Problems) == 0 {
		out.Warnings = append(out.Warnings, "no_problems_from_detector")
		return out, nil
	}
	switch ctx.input.DecisionerID {
	case DecisionerTire:
		if ctx.input.DetectorID != DetectorTireLabProblems {
			out.Warnings = append(out.Warnings, "incompatible_decisioner_detector:tire_problem_decision_v1_requires_tire_lab_problem_groups_v1")
			return out, nil
		}
		out.Decisions = tireProblemDecisions(problems.Problems)
	case DecisionerLegacyRoad:
		if ctx.input.DetectorID != DetectorLegacyRoadEvents {
			out.Warnings = append(out.Warnings, "incompatible_decisioner_detector:legacy_road_decision_v1_requires_legacy_road_events_v1")
			return out, nil
		}
		decision, err := s.GetRoadTuningDecision(ctx.input.SessionID)
		if err != nil {
			out.Warnings = append(out.Warnings, "legacy_road_decision_error:"+err.Error())
			out.Decisions = legacyProblemFallbackDecisions(problems.Problems)
			break
		}
		ctx.roadDecision = decision
		out.Decisions = []TuningDecision{decisionFromRoadDecision(*decision)}
	default:
		out.Warnings = append(out.Warnings, "unknown_decisioner:"+ctx.input.DecisionerID)
	}
	if len(out.Decisions) > 0 {
		out.Status = pipelineStatusReady
	}
	return out, nil
}

func tireProblemDecisions(problems []TuningProblem) []TuningDecision {
	out := make([]TuningDecision, 0, len(problems))
	for _, problem := range problems {
		cause := tirePipelinePrimaryCause(problem)
		out = append(out, TuningDecision{
			ID:              "decision:" + problem.ID,
			ProblemID:       problem.ID,
			Phase:           problem.Phase,
			PrimaryCause:    cause,
			ShouldTune:      tirePipelineShouldTune(problem),
			Confidence:      problem.Confidence,
			Rationale:       "tire_problem_group_decision",
			DocumentContext: "tire_lab_problem_groups_v1",
			Evidence:        cloneFloatMap(problem.Evidence),
		})
	}
	return out
}

func tirePipelinePrimaryCause(problem TuningProblem) string {
	if problem.Type == "lateral_limit" {
		if problem.LimitedAxle == "front" {
			if problem.Phase == "high_speed_corner" {
				return "front_high_speed_lateral_limit"
			}
			return "front_mechanical_lateral_limit"
		}
		if problem.LimitedAxle == "rear" {
			return "rear_lateral_stability_limit"
		}
		return "four_wheel_lateral_limit"
	}
	if problem.Type == "traction_limit" {
		return "drive_torque_exceeds_tire_grip"
	}
	if problem.Type == "braking_limit" {
		if problem.LimitedAxle == "rear" {
			return "rear_brake_or_decel_instability"
		}
		return "front_brake_overload"
	}
	if problem.Type == "platform_risk" {
		return "platform_travel_or_load_risk"
	}
	if problem.Type == "thermal_risk" {
		return "tire_temperature_risk"
	}
	if problem.Type == "combined_limit" {
		return "combined_longitudinal_lateral_overload"
	}
	return problem.Type
}

func tirePipelineShouldTune(problem TuningProblem) bool {
	if problem.Confidence == quickConfidenceInvalid || problem.Confidence == quickConfidenceLow || problem.Type == "data_insufficient" {
		return false
	}
	if containsStringValue(problem.OperationTags, "handbrake_active") && problem.Phase == "drift" {
		return false
	}
	return true
}

func decisionFromRoadDecision(decision RoadTuningDecision) TuningDecision {
	return TuningDecision{
		ID:              "decision:road:" + decision.SymptomID,
		ProblemID:       problemIDFromRoadDecision(decision),
		Phase:           decision.Phase,
		PrimaryCause:    decision.PrimaryCause,
		ShouldTune:      decision.Status == roadDecisionReady || decision.Status == roadDecisionRollback,
		Confidence:      decision.Confidence,
		Rationale:       decision.Reason,
		DocumentContext: "legacy_road_decision_v1",
		Evidence:        cloneFloatMap(decision.Evidence),
	}
}

func problemIDFromRoadDecision(decision RoadTuningDecision) string {
	if decision.RelatedIssueGroup != nil {
		return "legacy:" + decision.RelatedIssueGroup.ID
	}
	if decision.SymptomID != "" {
		return "legacy:" + decision.SymptomID
	}
	return "legacy:road_decision"
}

func legacyProblemFallbackDecisions(problems []TuningProblem) []TuningDecision {
	out := make([]TuningDecision, 0, len(problems))
	for _, problem := range problems {
		out = append(out, TuningDecision{
			ID:              "decision:" + problem.ID,
			ProblemID:       problem.ID,
			Phase:           problem.Phase,
			PrimaryCause:    problem.Reason,
			ShouldTune:      true,
			Confidence:      problem.Confidence,
			Rationale:       "legacy_problem_fallback",
			DocumentContext: "legacy_road_events_v1",
			Evidence:        cloneFloatMap(problem.Evidence),
		})
	}
	return out
}

func (s *Store) runPipelineInterpreter(ctx *tuningPipelineContext, decisions TuningDecisionSet) (TuningAdviceSet, []string) {
	out := TuningAdviceSet{
		InterpreterID: ctx.input.InterpreterID,
		Status:        pipelineStatusNoData,
	}
	if len(decisions.Decisions) == 0 {
		out.Warnings = append(out.Warnings, "no_decisions_from_decisioner")
		return out, nil
	}
	switch ctx.input.InterpreterID {
	case InterpreterTireRepair:
		if ctx.input.DecisionerID != DecisionerTire {
			out.Warnings = append(out.Warnings, "incompatible_interpreter_decisioner:tire_repair_explainer_v1_requires_tire_problem_decision_v1")
			return out, nil
		}
		if ctx.tireAnalysis.Status == "" {
			ctx.tireAnalysis = BuildTireIssueAnalysis(ctx.samples)
		}
		ctx.tireAdvice = BuildTireIssueAdviceFromAnalysis(ctx.tireAnalysis)
		out.Advice = adviceFromTireIssueAdvice(ctx.tireAdvice)
	case InterpreterLegacyAdvice:
		if ctx.input.DecisionerID != DecisionerLegacyRoad {
			out.Warnings = append(out.Warnings, "incompatible_interpreter_decisioner:legacy_advice_planner_v1_requires_legacy_road_decision_v1")
			return out, nil
		}
		out.Advice = s.legacyPipelineAdvice(ctx)
	case InterpreterRoadDocsV12:
		out.DocumentSources = roadDocsV12Sources()
		out.Advice = docsV12InterpreterAdvice(decisions.Decisions)
	default:
		out.Warnings = append(out.Warnings, "unknown_interpreter:"+ctx.input.InterpreterID)
	}
	if len(out.Advice) > 0 {
		out.Status = pipelineStatusReady
	}
	return out, nil
}

func adviceFromTireIssueAdvice(input TireIssueAdvice) []TuningAdvice {
	out := []TuningAdvice{}
	for _, action := range input.PriorityActions {
		out = append(out, TuningAdvice{
			ID:             "advice:" + action.ID,
			DecisionID:     "decision:tire:" + action.IssueGroupID,
			ProblemID:      "tire:" + action.IssueGroupID,
			Layer:          action.Layer,
			Category:       action.Category,
			Scope:          action.Scope,
			Direction:      action.Direction,
			RelatedFields:  append([]string(nil), action.RelatedFields...),
			Rationale:      action.Rationale,
			VerifyEvidence: append([]string(nil), action.VerifyEvidence...),
			TrustLevel:     action.Confidence,
			MissingInputs:  append([]string(nil), action.MissingInputs...),
			ConflictReason: action.ConflictReason,
			CanApply:       false,
		})
	}
	return out
}

func (s *Store) legacyPipelineAdvice(ctx *tuningPipelineContext) []TuningAdvice {
	if ctx.roadDecision == nil && ctx.input.SessionID > 0 {
		if decision, err := s.GetRoadTuningDecision(ctx.input.SessionID); err == nil {
			ctx.roadDecision = decision
		}
	}
	out := []TuningAdvice{}
	if ctx.roadDecision != nil {
		for _, action := range ctx.roadDecision.Actions {
			out = append(out, TuningAdvice{
				ID:             "advice:" + action.ID,
				DecisionID:     "decision:road:" + ctx.roadDecision.SymptomID,
				ProblemID:      problemIDFromRoadDecision(*ctx.roadDecision),
				Layer:          action.Role,
				Category:       action.Category,
				Scope:          action.Family,
				Direction:      action.Direction,
				RelatedFields:  fieldList(action.FieldKey),
				Rationale:      firstNonEmpty(action.Rationale, action.Reason),
				VerifyEvidence: sortedEvidenceKeys(action.Evidence),
				TrustLevel:     firstNonEmpty(action.TrustLevel, action.Confidence),
				ConflictReason: action.ConflictReason,
				CanApply:       false,
				Evidence:       cloneFloatMap(action.Evidence),
			})
		}
	}
	if len(out) > 0 || ctx.summary == nil {
		return out
	}
	for _, action := range ctx.summary.WholeCarPlan.Actions {
		out = append(out, TuningAdvice{
			ID:             "advice:whole_car:" + action.Family + ":" + action.Item,
			DecisionID:     "decision:legacy:" + action.Family,
			ProblemID:      "legacy:" + action.Family,
			Layer:          "primary",
			Category:       action.Category,
			Scope:          action.Family,
			Direction:      action.Direction,
			RelatedFields:  nil,
			Rationale:      action.Reason,
			VerifyEvidence: sortedEvidenceKeys(action.Evidence),
			TrustLevel:     action.Confidence,
			CanApply:       false,
			Evidence:       cloneFloatMap(action.Evidence),
		})
	}
	return out
}

func docsV12InterpreterAdvice(decisions []TuningDecision) []TuningAdvice {
	out := []TuningAdvice{}
	for _, decision := range decisions {
		spec := docsV12SpecForDecision(decision)
		if spec.category == "" {
			continue
		}
		out = append(out, TuningAdvice{
			ID:              "advice:docs_v12:" + sanitizeID(decision.ID),
			DecisionID:      decision.ID,
			ProblemID:       decision.ProblemID,
			Layer:           "explain_only",
			Category:        spec.category,
			Scope:           spec.scope,
			Direction:       spec.direction,
			RelatedFields:   append([]string(nil), spec.fields...),
			Rationale:       spec.rationale,
			VerifyEvidence:  append([]string(nil), spec.verifyEvidence...),
			TrustLevel:      firstNonEmpty(decision.Confidence, quickConfidenceMedium),
			MissingInputs:   []string{},
			ConflictReason:  "docs_v12_interpreter_outputs_no_write_values",
			CanApply:        false,
			DocumentSources: roadDocsV12Sources(),
			Evidence:        cloneFloatMap(decision.Evidence),
		})
	}
	return out
}

type docsV12AdviceSpec struct {
	category       string
	scope          string
	direction      string
	fields         []string
	rationale      string
	verifyEvidence []string
}

func docsV12SpecForDecision(decision TuningDecision) docsV12AdviceSpec {
	text := strings.ToLower(strings.Join([]string{decision.PrimaryCause, decision.Phase, decision.Rationale}, " "))
	switch {
	case strings.Contains(text, "traction") || strings.Contains(text, "drive_torque") || strings.Contains(text, "power"):
		return docsV12AdviceSpec{
			category:       "power_to_tire",
			scope:          "driven_wheels",
			direction:      "reduce_wheel_torque_or_drive_lock",
			fields:         []string{"finalDrive", "gear1", "gear2", "gear3", "frontDiffAccel", "rearDiffAccel", "centerDiffBalance"},
			rationale:      "Docs v1.2 treats gearing and differential as Forza slider/display levers; verify whether torque delivery exceeds tire traction before changing chassis balance.",
			verifyEvidence: []string{"driven_wheel_slip_ratio", "throttle", "rpm", "speed_gain", "gear"},
		}
	case strings.Contains(text, "brake") || strings.Contains(text, "decel"):
		return docsV12AdviceSpec{
			category:       "brake",
			scope:          "front_rear_balance",
			direction:      "verify_brake_balance_pressure_and_decel_lock",
			fields:         []string{"brakeBalance", "brakePressure", "frontDiffDecel", "rearDiffDecel"},
			rationale:      "Docs v1.2 uses conservative road brake baselines by drivetrain; explain braking issues through balance, pressure, and decel lock before changing unrelated systems.",
			verifyEvidence: []string{"brake", "deceleration_g", "front_slip_ratio", "rear_slip_ratio"},
		}
	case strings.Contains(text, "high_speed") || strings.Contains(text, "platform") || strings.Contains(text, "aero"):
		return docsV12AdviceSpec{
			category:       "platform_aero",
			scope:          "front_rear_platform",
			direction:      "verify_platform_then_use_tiered_ride_height_or_aero",
			fields:         []string{"frontRideHeight", "rearRideHeight", "frontAero", "rearAero", "frontSpring", "rearSpring", "frontRebound", "rearRebound", "frontBump", "rearBump"},
			rationale:      "Docs v1.2 keeps ride height and aero as low/medium/high tier explanations, not precise write values; use them only after verifying platform evidence.",
			verifyEvidence: []string{"speed_band", "suspension_offset", "combined_slip", "g_force"},
		}
	case strings.Contains(text, "thermal") || strings.Contains(text, "temperature"):
		return docsV12AdviceSpec{
			category:       "tire_pressure",
			scope:          "front_rear_tires",
			direction:      "verify_pressure_window_and_slip_heat",
			fields:         []string{"frontTirePressure", "rearTirePressure"},
			rationale:      "Docs v1.4 uses BAR as the primary tire-pressure unit; thermal issues need pressure, slip, and temperature evidence before any 0.02-0.03 BAR correction is trusted.",
			verifyEvidence: []string{"front_tire_temp", "rear_tire_temp", "combined_slip"},
		}
	case strings.Contains(text, "lateral") || strings.Contains(text, "understeer") || strings.Contains(text, "oversteer"):
		return docsV12AdviceSpec{
			category:       "mechanical_grip",
			scope:          "front_rear_balance",
			direction:      "rebalance_tire_contact_and_load_transfer",
			fields:         []string{"frontTirePressure", "rearTirePressure", "frontCamber", "rearCamber", "frontToe", "rearToe", "caster", "frontArb", "rearArb", "frontSpring", "rearSpring", "frontRebound", "rearRebound"},
			rationale:      "Docs v1.2/v1.4 make tire pressure, alignment, anti-roll bars, springs, and damping explicit Forza display/sliders; use BAR-first tire pressure only when slip and heat evidence support it.",
			verifyEvidence: []string{"front_combined_slip", "rear_combined_slip", "slip_angle", "speed_band", "steer"},
		}
	default:
		return docsV12AdviceSpec{
			category:       "observe",
			scope:          "vehicle",
			direction:      "collect_more_evidence_before_tuning",
			fields:         nil,
			rationale:      "Docs v1.2 interpreter did not find a specific road baseline lever for this decision; keep it as an observation.",
			verifyEvidence: []string{"sample_quality", "phase", "speed", "inputs"},
		}
	}
}

func roadDocsV12Sources() []string {
	return []string{
		"FH6_调校基线生成器_开发功能说明_v1.0.md",
		"FH6_调校基线生成器_弹簧公式修正开发修改文档_v1.1.md",
		"FH6_调校基线生成器_参数口径统一修正开发修改文档_v1.2.md",
	}
}

func fieldList(fieldKey string) []string {
	if strings.TrimSpace(fieldKey) == "" {
		return nil
	}
	return []string{fieldKey}
}

func sortedEvidenceKeys(input map[string]float64) []string {
	keys := make([]string, 0, len(input))
	for key := range input {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func sanitizeID(input string) string {
	replacer := strings.NewReplacer(":", "_", "/", "_", " ", "_", ",", "_")
	return replacer.Replace(input)
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}
