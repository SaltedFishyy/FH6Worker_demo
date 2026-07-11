package advisor

import (
	"strings"
	"testing"

	"fh6worker/internal/storage"
	"fh6worker/internal/telemetry"
)

func TestGenerateReportWithNoEvents(t *testing.T) {
	report := GenerateTuningReport(storage.TelemetrySession{ID: 1, SessionName: "Empty"}, nil, nil, "zh")
	if !strings.Contains(report, "未检测到明确调校事件") {
		t.Fatalf("report did not describe no-event state:\n%s", report)
	}
}

func TestGenerateReportWithEvents(t *testing.T) {
	rearDiffAccel := 60.0
	gear2 := 2.20
	profile := &storage.TuneProfile{CarName: "BMW M3", CarClass: "S1", Drivetrain: "RWD", UseCase: "Road", RearDiffAccel: &rearDiffAccel, Gear2: &gear2}
	events := []telemetry.DetectedEvent{
		{
			Type:       "corner_exit_oversteer",
			Severity:   "high",
			Segment:    "corner_exit",
			DurationMS: 700,
			Evidence:   map[string]float64{"rear_combined_slip": 1.42, "speed_kmh": 68, "gear": 2},
			SuggestedActions: []telemetry.SuggestedAction{
				{Priority: 0, Category: "differential", Item: "rear_diff_accel", Direction: "decrease", Amount: "3%-5%", Reason: "reduce power oversteer"},
				{Priority: 1, Category: "gearing", Item: "current_gear", Direction: "decrease", Amount: "3%-5%", Reason: "reduce wheel torque on exit"},
			},
		},
	}

	report := GenerateTuningReport(storage.TelemetrySession{ID: 1, SessionName: "Run"}, profile, events, "zh")
	for _, want := range []string{"出弯甩尾", "后轮综合滑移", "后差速加速", "调校说明", "下一轮测试方法"} {
		if !strings.Contains(report, want) {
			t.Fatalf("report missing %q:\n%s", want, report)
		}
	}
	for _, want := range []string{"60 -> 57", "2.20 -> 2.13"} {
		if !strings.Contains(report, want) {
			t.Fatalf("report missing concrete adjustment %q:\n%s", want, report)
		}
	}
	if strings.Contains(report, "3%-5%") {
		t.Fatalf("report should prefer concrete tune deltas when profile values exist:\n%s", report)
	}
}

func TestGenerateReportLimitsPrimaryActions(t *testing.T) {
	event := telemetry.DetectedEvent{
		Type:     "launch_wheelspin",
		Severity: "high",
		Segment:  "launch",
		Evidence: map[string]float64{"rear_slip_ratio": 1.5},
	}
	for i := 0; i < 5; i++ {
		event.SuggestedActions = append(event.SuggestedActions, telemetry.SuggestedAction{
			Priority:  i,
			Category:  "gearing",
			Item:      "gear_" + string(rune('1'+i)),
			Direction: "decrease",
			Amount:    "1%",
			Reason:    "reduce wheel torque during launch",
		})
	}
	report := GenerateTuningReport(storage.TelemetrySession{ID: 1}, nil, []telemetry.DetectedEvent{event}, "en")
	if !strings.Contains(report, "Lower Priority Items") {
		t.Fatalf("expected lower-priority section:\n%s", report)
	}
}

func TestGenerateReportIncludesTestConditions(t *testing.T) {
	report := GenerateTuningReport(storage.TelemetrySession{
		ID:               1,
		SessionName:      "Conditions",
		DriverMode:       "auto",
		BrakeAssist:      "abs_on",
		SteeringAssist:   "simulation",
		TractionControl:  "off",
		StabilityControl: "off",
		Shifting:         "manual",
		LaunchControl:    "off",
	}, nil, nil, "en")
	for _, want := range []string{"Test conditions", "Driver=Auto driver", "Brake=ABS on", "Steering=Simulation"} {
		if !strings.Contains(report, want) {
			t.Fatalf("report missing %q:\n%s", want, report)
		}
	}

	unknownReport := GenerateTuningReport(storage.TelemetrySession{ID: 2, SessionName: "Unknown"}, nil, nil, "en")
	if !strings.Contains(unknownReport, "lower confidence") {
		t.Fatalf("expected unknown-condition confidence note:\n%s", unknownReport)
	}
}
