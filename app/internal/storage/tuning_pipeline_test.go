package storage

import (
	"strings"
	"testing"
)

func TestListTuningModelPipelinesCatalog(t *testing.T) {
	catalog := ListTuningModelPipelines()
	if len(catalog.Detectors) != 2 {
		t.Fatalf("expected 2 detectors, got %d", len(catalog.Detectors))
	}
	if len(catalog.Decisioners) != 2 {
		t.Fatalf("expected 2 decisioners, got %d", len(catalog.Decisioners))
	}
	if len(catalog.Interpreters) != 3 {
		t.Fatalf("expected 3 interpreters, got %d", len(catalog.Interpreters))
	}
	if !catalogHasComponent(catalog.Interpreters, InterpreterRoadDocsV12) {
		t.Fatalf("expected docs v1.2 interpreter in catalog")
	}
}

func TestRunTuningModelPipelineIncompatibleCombinationWarns(t *testing.T) {
	var store Store
	result, err := store.RunTuningModelPipeline(TuningPipelineRunInput{
		SourceType:    TuningSourceTireLabCurrent,
		DetectorID:    DetectorLegacyRoadEvents,
		DecisionerID:  DecisionerLegacyRoad,
		InterpreterID: InterpreterLegacyAdvice,
	}, nil)
	if err != nil {
		t.Fatalf("expected warning result, got error %v", err)
	}
	if result.Status != pipelineStatusIncompatible {
		t.Fatalf("expected incompatible status, got %s", result.Status)
	}
	if !containsWarning(result.Warnings, "incompatible_detector_source") {
		t.Fatalf("expected detector compatibility warning, got %#v", result.Warnings)
	}
}

func TestTireLabProblemGroupsConvertToUnifiedProblems(t *testing.T) {
	group := testTireIssueGroup("traction_limit", "corner_exit", "rear", quickConfidenceHigh)
	group.OperationTags = []string{"throttle_on", "speed_rising"}
	problems := problemsFromTireAnalysis(testTireAdviceAnalysis(group))
	if len(problems) != 1 {
		t.Fatalf("expected one problem, got %d", len(problems))
	}
	problem := problems[0]
	if problem.ID != "tire:"+group.ID || problem.Type != "traction_limit" || problem.LimitedAxle != "rear" {
		t.Fatalf("unexpected problem mapping: %+v", problem)
	}
	if problem.Evidence["avg_throttle"] == 0 {
		t.Fatalf("expected representative evidence to be preserved")
	}
}

func TestDocsV12InterpreterOutputsExplainOnlyAdvice(t *testing.T) {
	decisions := []TuningDecision{
		{
			ID:           "decision:tire:traction",
			ProblemID:    "tire:traction",
			Phase:        "corner_exit",
			PrimaryCause: "drive_torque_exceeds_tire_grip",
			ShouldTune:   true,
			Confidence:   quickConfidenceHigh,
			Evidence:     map[string]float64{"rear_slip_ratio_p90": 1.1},
		},
	}
	advice := docsV12InterpreterAdvice(decisions)
	if len(advice) != 1 {
		t.Fatalf("expected one advice, got %d", len(advice))
	}
	if advice[0].CanApply {
		t.Fatalf("docs v1.2 interpreter must not generate directly applicable values")
	}
	if advice[0].Category != "power_to_tire" {
		t.Fatalf("expected power-to-tire category, got %+v", advice[0])
	}
	if len(advice[0].DocumentSources) != 3 {
		t.Fatalf("expected document sources, got %#v", advice[0].DocumentSources)
	}
	if !containsStringValue(advice[0].RelatedFields, "finalDrive") || !containsStringValue(advice[0].RelatedFields, "rearDiffAccel") {
		t.Fatalf("expected gearing and diff fields, got %#v", advice[0].RelatedFields)
	}
}

func catalogHasComponent(items []TuningPipelineComponent, id string) bool {
	for _, item := range items {
		if item.ID == id {
			return true
		}
	}
	return false
}

func containsWarning(warnings []string, fragment string) bool {
	for _, warning := range warnings {
		if strings.Contains(warning, fragment) {
			return true
		}
	}
	return false
}
