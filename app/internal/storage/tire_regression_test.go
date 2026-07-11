package storage

import (
	"testing"

	"fh6worker/internal/telemetry"
)

func TestTireRegressionSampleRoundTrip(t *testing.T) {
	dir := t.TempDir()
	samples := tireModelSamples(func(frame *telemetry.NormalizedTelemetry) {
		frame.Steer01 = 0.45
		frame.WheelFL.CombinedSlip = 1.08
		frame.WheelFR.CombinedSlip = 1.02
		frame.WheelFL.SlipAngle = 0.75
		frame.WheelFR.SlipAngle = 0.70
		frame.WheelRL.CombinedSlip = 0.35
		frame.WheelRR.CombinedSlip = 0.36
	})
	sample, err := SaveTireRegressionSample(dir, TireRegressionSampleInput{
		Name:          "Front limit",
		Scenario:      "corner_entry_understeer",
		WindowSeconds: 10,
	}, samples, nil)
	if err != nil {
		t.Fatalf("save sample: %v", err)
	}
	if sample.ID == "" || sample.SampleCount == 0 || len(sample.Samples) == 0 {
		t.Fatalf("sample metadata not populated: %#v", sample)
	}
	if sample.Snapshot.IssueAnalysis.GroupCount == 0 {
		t.Fatalf("sample issue analysis missing groups: %#v", sample.Snapshot.IssueAnalysis)
	}
	if sample.Expected.RequiredGripTypes[0] != "lateral_limit" || sample.Expected.AllowedAxles[0] != "front" {
		t.Fatalf("expected draft = %#v, want front lateral limit", sample.Expected)
	}
	summaries, err := ListTireRegressionSamples(dir)
	if err != nil {
		t.Fatalf("list samples: %v", err)
	}
	if len(summaries) != 1 || summaries[0].ID != sample.ID {
		t.Fatalf("summaries = %#v", summaries)
	}
	result, err := RunTireRegressionSample(dir, sample.ID)
	if err != nil {
		t.Fatalf("run sample: %v", err)
	}
	if !result.Passed {
		t.Fatalf("result = %#v, want passed", result)
	}
	if result.Actual.IssueAnalysis.GroupCount == 0 {
		t.Fatalf("actual issue analysis missing groups: %#v", result.Actual.IssueAnalysis)
	}
}

func TestTireRegressionExpectationFailureAndDelete(t *testing.T) {
	dir := t.TempDir()
	samples := tireModelSamples(func(frame *telemetry.NormalizedTelemetry) {
		frame.Throttle01 = 0.82
		frame.WheelRL.CombinedSlip = 0.95
		frame.WheelRR.CombinedSlip = 0.98
		frame.WheelRL.SlipRatio = 0.62
		frame.WheelRR.SlipRatio = 0.66
	})
	sample, err := SaveTireRegressionSample(dir, TireRegressionSampleInput{Name: "Rear traction"}, samples, nil)
	if err != nil {
		t.Fatalf("save sample: %v", err)
	}
	badExpected := sample.Expected
	badExpected.ForbiddenGripTypes = []string{"traction_limit"}
	if err := UpdateTireRegressionSampleExpectation(dir, sample.ID, badExpected); err != nil {
		t.Fatalf("update expectation: %v", err)
	}
	result, err := RunTireRegressionSample(dir, sample.ID)
	if err != nil {
		t.Fatalf("run sample: %v", err)
	}
	if result.Passed || !containsString(result.Failures, "forbidden_grip_detected") {
		t.Fatalf("result = %#v, want forbidden grip failure", result)
	}
	if err := DeleteTireRegressionSample(dir, sample.ID); err != nil {
		t.Fatalf("delete sample: %v", err)
	}
	summaries, err := ListTireRegressionSamples(dir)
	if err != nil {
		t.Fatalf("list after delete: %v", err)
	}
	if len(summaries) != 0 {
		t.Fatalf("summaries after delete = %#v", summaries)
	}
}
