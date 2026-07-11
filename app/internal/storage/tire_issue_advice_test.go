package storage

import "testing"

func TestBuildTireIssueAdviceNoGroups(t *testing.T) {
	analysis := TireIssueAnalysis{Status: tireModelStatusReady, UpdatedAt: "now", Groups: []TireIssueGroup{}}
	advice := BuildTireIssueAdviceFromAnalysis(analysis)
	if advice.Status != tireAdviceStatusEmpty {
		t.Fatalf("expected empty status, got %s", advice.Status)
	}
	if len(advice.PriorityActions) != 0 {
		t.Fatalf("expected no priority actions, got %d", len(advice.PriorityActions))
	}
}

func TestTireIssueAdviceFrontLateralLowSpeedUsesMechanicalGrip(t *testing.T) {
	group := testTireIssueGroup("lateral_limit", "low_speed_corner", "front", quickConfidenceHigh)
	advice := BuildTireIssueAdviceFromAnalysis(testTireAdviceAnalysis(group))
	if len(advice.Groups) != 1 || len(advice.Groups[0].Actions) == 0 {
		t.Fatalf("expected group actions")
	}
	action := advice.Groups[0].Actions[0]
	if action.Category != "antiroll" || action.Direction != "increase_front_mechanical_grip" {
		t.Fatalf("expected front mechanical grip primary action, got %+v", action)
	}
	if action.Category == "aero_platform" {
		t.Fatalf("low speed front lateral limit should not prioritize aero")
	}
}

func TestTireIssueAdviceFrontLateralHighSpeedUsesPlatform(t *testing.T) {
	group := testTireIssueGroup("lateral_limit", "high_speed_corner", "front", quickConfidenceHigh)
	advice := BuildTireIssueAdviceFromAnalysis(testTireAdviceAnalysis(group))
	action := advice.Groups[0].Actions[0]
	if action.Category != "aero_platform" || action.Direction != "increase_front_high_speed_support" {
		t.Fatalf("expected high-speed platform primary action, got %+v", action)
	}
}

func TestTireIssueAdviceRearTractionUsesDiffAndGearing(t *testing.T) {
	group := testTireIssueGroup("traction_limit", "corner_exit", "rear", quickConfidenceHigh)
	group.OperationTags = []string{"throttle_on", "speed_rising"}
	advice := BuildTireIssueAdviceFromAnalysis(testTireAdviceAnalysis(group))
	actions := advice.Groups[0].Actions
	if len(actions) < 2 {
		t.Fatalf("expected primary and alternative actions, got %d", len(actions))
	}
	if actions[0].Category != "differential" || actions[0].Scope != "rear" {
		t.Fatalf("expected rear diff primary action, got %+v", actions[0])
	}
	if actions[1].Category != "gearing" {
		t.Fatalf("expected gearing alternative, got %+v", actions[1])
	}
}

func TestTireIssueAdviceFrontBrakingUsesBrakeDirection(t *testing.T) {
	group := testTireIssueGroup("braking_limit", "braking", "front", quickConfidenceHigh)
	group.OperationTags = []string{"heavy_brake", "speed_falling"}
	advice := BuildTireIssueAdviceFromAnalysis(testTireAdviceAnalysis(group))
	action := advice.Groups[0].Actions[0]
	if action.Category != "brake" || action.Direction != "move_brake_balance_rearward" {
		t.Fatalf("expected front brake balance action, got %+v", action)
	}
}

func TestTireIssueAdviceHandbrakeDriftDoesNotTune(t *testing.T) {
	group := testTireIssueGroup("traction_limit", "drift", "rear", quickConfidenceHigh)
	group.OperationTags = []string{"handbrake_active"}
	group.DriftSource = "handbrake_initiated"
	advice := BuildTireIssueAdviceFromAnalysis(testTireAdviceAnalysis(group))
	action := advice.Groups[0].Actions[0]
	if action.TuneRecommended {
		t.Fatalf("handbrake drift should not recommend tuning: %+v", action)
	}
	if action.Category != "driver_input" || action.Direction != "avoid_tuning" {
		t.Fatalf("expected driver behavior action, got %+v", action)
	}
}

func TestTireIssueAdvicePriorityLimitAndSnapshot(t *testing.T) {
	groups := []TireIssueGroup{
		testTireIssueGroup("lateral_limit", "low_speed_corner", "front", quickConfidenceHigh),
		testTireIssueGroup("traction_limit", "corner_exit", "rear", quickConfidenceHigh),
		testTireIssueGroup("braking_limit", "braking", "front", quickConfidenceHigh),
		testTireIssueGroup("platform_risk", "high_speed_corner", "front", quickConfidenceMedium),
	}
	advice := BuildTireIssueAdviceFromAnalysis(testTireAdviceAnalysis(groups...))
	if len(advice.PriorityActions) != 3 {
		t.Fatalf("expected 3 priority actions, got %d", len(advice.PriorityActions))
	}
	snapshot := BuildTireDiagnosticSnapshot(TireModelDiagnostic{IssueAnalysis: testTireAdviceAnalysis(groups...), IssueAdvice: advice})
	if snapshot.IssueAdvice.Status != tireAdviceStatusReady || len(snapshot.IssueAdvice.PriorityActions) != 3 {
		t.Fatalf("expected snapshot issue advice to be preserved, got %+v", snapshot.IssueAdvice)
	}
}

func testTireAdviceAnalysis(groups ...TireIssueGroup) TireIssueAnalysis {
	return TireIssueAnalysis{
		Status:     tireModelStatusReady,
		UpdatedAt:  "2026-05-22T00:00:00Z",
		GroupCount: len(groups),
		Groups:     groups,
	}
}

func testTireIssueGroup(issueType, phase, axle, confidence string) TireIssueGroup {
	return TireIssueGroup{
		ID:              "grp-test",
		Type:            issueType,
		Phase:           phase,
		OperationTags:   []string{},
		LimitType:       issueType,
		LimitedAxle:     axle,
		Count:           2,
		TotalDurationMS: 1800,
		SpeedMinKmh:     60,
		SpeedMaxKmh:     140,
		SpeedAvgKmh:     95,
		Confidence:      confidence,
		DataQuality:     "valid",
		RiskLevel:       "high",
		RepresentativeEvidence: map[string]float64{
			"front_combined_slip_p90": 1.2,
			"rear_combined_slip_p90":  0.55,
			"front_slip_ratio_p90":    0.2,
			"rear_slip_ratio_p90":     0.4,
			"avg_speed_kmh":           95,
			"avg_throttle":            0.7,
			"avg_brake":               0.4,
		},
	}
}
