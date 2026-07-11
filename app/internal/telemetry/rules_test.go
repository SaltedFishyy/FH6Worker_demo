package telemetry

import "testing"

func TestRuleEngineDetectsMVP2Events(t *testing.T) {
	tests := []struct {
		name      string
		eventType string
		frame     NormalizedTelemetry
	}{
		{
			name:      "launch wheelspin",
			eventType: "launch_wheelspin",
			frame: baseRuleFrame(NormalizedTelemetry{
				Gear:              1,
				SpeedKmh:          35,
				Throttle01:        0.92,
				RpmRatio:          0.75,
				RearSlipRatioAvg:  1.35,
				FrontSlipRatioAvg: 0.25,
			}),
		},
		{
			name:      "launch bog down",
			eventType: "launch_bog_down",
			frame: baseRuleFrame(NormalizedTelemetry{
				Gear:              1,
				SpeedKmh:          24,
				Throttle01:        0.95,
				RpmRatio:          0.20,
				FrontSlipRatioAvg: 0.08,
				RearSlipRatioAvg:  0.10,
			}),
		},
		{
			name:      "short gear",
			eventType: "short_gear",
			frame: baseRuleFrame(NormalizedTelemetry{
				Gear:       2,
				SpeedKmh:   72,
				Throttle01: 0.70,
				RpmRatio:   0.98,
			}),
		},
		{
			name:      "front brake lockup",
			eventType: "front_brake_lockup",
			frame: baseRuleFrame(NormalizedTelemetry{
				SpeedKmh:             112,
				Brake01:              0.88,
				FrontCombinedSlipAvg: 1.35,
				RearCombinedSlipAvg:  0.55,
				FrontSlipRatioAvg:    0.95,
				RearSlipRatioAvg:     0.20,
			}),
		},
		{
			name:      "rear brake lockup",
			eventType: "rear_brake_lockup",
			frame: baseRuleFrame(NormalizedTelemetry{
				SpeedKmh:             94,
				Brake01:              0.86,
				FrontCombinedSlipAvg: 0.45,
				RearCombinedSlipAvg:  1.45,
				YawRate:              0.6,
			}),
		},
		{
			name:      "corner entry understeer",
			eventType: "corner_entry_understeer",
			frame: baseRuleFrame(NormalizedTelemetry{
				SpeedKmh:             86,
				Steer01:              0.55,
				Throttle01:           0.18,
				Brake01:              0.22,
				FrontCombinedSlipAvg: 1.05,
				RearCombinedSlipAvg:  0.45,
			}),
		},
		{
			name:      "corner exit oversteer",
			eventType: "corner_exit_oversteer",
			frame: baseRuleFrame(NormalizedTelemetry{
				SpeedKmh:             92,
				Gear:                 3,
				Steer01:              -0.35,
				Throttle01:           0.86,
				FrontCombinedSlipAvg: 0.45,
				RearCombinedSlipAvg:  1.15,
				YawRate:              0.7,
			}),
		},
		{
			name:      "high speed four wheel slide",
			eventType: "high_speed_four_wheel_slide",
			frame: baseRuleFrame(NormalizedTelemetry{
				SpeedKmh:             182,
				FrontCombinedSlipAvg: 1.05,
				RearCombinedSlipAvg:  1.00,
				YawRate:              0.5,
			}),
		},
		{
			name:      "suspension bottom out",
			eventType: "suspension_bottom_out",
			frame: baseRuleFrame(NormalizedTelemetry{
				SpeedKmh:           145,
				SuspensionFrontAvg: 0.98,
				SuspensionRearAvg:  0.74,
				WheelFL:            NormalizedWheelTelemetry{SuspensionTravel: 0.98},
				WheelFR:            NormalizedWheelTelemetry{SuspensionTravel: 0.96},
				WheelRL:            NormalizedWheelTelemetry{SuspensionTravel: 0.74},
				WheelRR:            NormalizedWheelTelemetry{SuspensionTravel: 0.72},
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine := NewRuleEngine()
			observeRepeated(engine, tt.frame, 4)
			events := engine.Events()
			if len(events) != 1 {
				t.Fatalf("event count = %d, want 1; events = %#v", len(events), events)
			}
			if events[0].Type != tt.eventType {
				t.Fatalf("event type = %q, want %q", events[0].Type, tt.eventType)
			}
			if events[0].DurationMS <= 0 {
				t.Fatalf("duration = %d, want positive", events[0].DurationMS)
			}
			if len(events[0].Evidence) == 0 {
				t.Fatal("expected evidence")
			}
			if len(events[0].SuggestedActions) == 0 {
				t.Fatal("expected suggested actions")
			}
		})
	}
}

func TestRuleEngineGearingSuggestionDirectionsMatchFHConvention(t *testing.T) {
	tests := []struct {
		name      string
		eventType string
		frame     NormalizedTelemetry
		config    RuleConfig
		item      string
		direction string
	}{
		{
			name:      "launch wheelspin uses taller first gear",
			eventType: "launch_wheelspin",
			frame: baseRuleFrame(NormalizedTelemetry{
				Gear:              1,
				SpeedKmh:          35,
				Throttle01:        0.92,
				RpmRatio:          0.75,
				RearSlipRatioAvg:  1.35,
				FrontSlipRatioAvg: 0.25,
			}),
			item:      "gear_1",
			direction: "decrease",
		},
		{
			name:      "launch bog down uses shorter first gear",
			eventType: "launch_bog_down",
			frame: baseRuleFrame(NormalizedTelemetry{
				Gear:              1,
				SpeedKmh:          24,
				Throttle01:        0.95,
				RpmRatio:          0.20,
				FrontSlipRatioAvg: 0.08,
				RearSlipRatioAvg:  0.10,
			}),
			item:      "gear_1",
			direction: "increase",
		},
		{
			name:      "short gear lowers the current gear ratio",
			eventType: "short_gear",
			frame: baseRuleFrame(NormalizedTelemetry{
				Gear:       2,
				SpeedKmh:   72,
				Throttle01: 0.70,
				RpmRatio:   0.98,
			}),
			item:      "current_gear",
			direction: "decrease",
		},
		{
			name:      "long gear raises the current gear ratio",
			eventType: "long_gear_bog_down",
			frame:     longGearBogDownFrame(),
			config:    RoadRacingRuleConfig(),
			item:      "current_gear",
			direction: "increase",
		},
		{
			name:      "top speed limited lowers the current gear ratio",
			eventType: "top_speed_limited_by_gearing",
			frame:     topSpeedLimitedByGearingFrame(),
			config:    RoadRacingRuleConfig(),
			item:      "current_gear",
			direction: "decrease",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine := NewRuleEngine()
			if tt.config.Events != nil {
				engine.SetConfig(tt.config)
			}
			observeRepeated(engine, tt.frame, 6)
			events := engine.Events()
			if len(events) == 0 || events[0].Type != tt.eventType {
				t.Fatalf("events = %#v, want first event %q", events, tt.eventType)
			}
			action := findSuggestedAction(events[0], tt.item)
			if action == nil {
				t.Fatalf("event actions = %#v, want item %q", events[0].SuggestedActions, tt.item)
			}
			if action.Direction != tt.direction {
				t.Fatalf("%s direction = %q, want %q", tt.item, action.Direction, tt.direction)
			}
		})
	}
}

func TestRuleEngineDebouncesContinuousEvents(t *testing.T) {
	engine := NewRuleEngine()
	frame := baseRuleFrame(NormalizedTelemetry{
		Gear:              1,
		SpeedKmh:          38,
		Throttle01:        0.94,
		RpmRatio:          0.76,
		FrontSlipRatioAvg: 0.30,
		RearSlipRatioAvg:  1.45,
	})

	observeRepeated(engine, frame, 8)
	if got := len(engine.Events()); got != 1 {
		t.Fatalf("continuous event count = %d, want 1", got)
	}

	clearFrame := frame
	clearFrame.TimeMS = 900
	clearFrame.Throttle01 = 0
	clearFrame.RearSlipRatioAvg = 0
	engine.Observe(clearFrame)

	frame.TimeMS = 1000
	observeRepeated(engine, frame, 4)
	if got := len(engine.Events()); got != 2 {
		t.Fatalf("second event count = %d, want 2", got)
	}
}

func findSuggestedAction(event DetectedEvent, item string) *SuggestedAction {
	for i := range event.SuggestedActions {
		if event.SuggestedActions[i].Item == item {
			return &event.SuggestedActions[i]
		}
	}
	return nil
}

func TestRoadRacingGearingEventsAreDisabledByDefault(t *testing.T) {
	engine := NewRuleEngine()
	for _, frame := range []NormalizedTelemetry{longGearBogDownFrame(), topSpeedLimitedByGearingFrame()} {
		observeRepeated(engine, frame, 5)
	}
	if got := len(engine.Events()); got != 0 {
		t.Fatalf("default road gearing event count = %d, want 0; events = %#v", got, engine.Events())
	}
}

func TestRoadRacingRuleConfigDetectsGearingEvents(t *testing.T) {
	engine := NewRuleEngine()
	engine.SetConfig(RoadRacingRuleConfig())

	observeRepeated(engine, longGearBogDownFrame(), 5)
	clear := longGearBogDownFrame()
	clear.TimeMS = 700
	clear.Throttle01 = 0
	engine.Observe(clear)
	observeRepeated(engine, topSpeedLimitedByGearingFrame(), 6)

	events := engine.Events()
	if len(events) != 2 {
		t.Fatalf("road gearing event count = %d, want 2; events = %#v", len(events), events)
	}
	if events[0].Type != "long_gear_bog_down" || events[1].Type != "top_speed_limited_by_gearing" {
		t.Fatalf("road gearing events = %#v", events)
	}
}

func TestRoadRacingRuleConfigDetectsAdditionalRoadEvents(t *testing.T) {
	tests := []struct {
		name      string
		eventType string
		frame     NormalizedTelemetry
	}{
		{
			name:      "mid corner understeer",
			eventType: "mid_corner_understeer",
			frame: baseRuleFrame(NormalizedTelemetry{
				SpeedKmh:             88,
				Steer01:              0.45,
				Throttle01:           0.42,
				Brake01:              0.02,
				FrontCombinedSlipAvg: 1.05,
				RearCombinedSlipAvg:  0.45,
			}),
		},
		{
			name:      "power understeer",
			eventType: "power_understeer",
			frame: baseRuleFrame(NormalizedTelemetry{
				SpeedKmh:             78,
				Steer01:              0.34,
				Throttle01:           0.82,
				FrontCombinedSlipAvg: 1.12,
				RearCombinedSlipAvg:  0.52,
			}),
		},
		{
			name:      "snap oversteer",
			eventType: "snap_oversteer",
			frame: baseRuleFrame(NormalizedTelemetry{
				SpeedKmh:             92,
				Steer01:              -0.26,
				Throttle01:           0.18,
				FrontCombinedSlipAvg: 0.45,
				RearCombinedSlipAvg:  1.15,
				YawRate:              0.95,
			}),
		},
		{
			name:      "tire overheat",
			eventType: "tire_overheat",
			frame: baseRuleFrame(NormalizedTelemetry{
				SpeedKmh:         125,
				TireTempFrontAvg: 120,
				TireTempRearAvg:  112,
			}),
		},
		{
			name:      "tire temp imbalance",
			eventType: "tire_temp_imbalance",
			frame: baseRuleFrame(NormalizedTelemetry{
				SpeedKmh:         118,
				TireTempFrontAvg: 105,
				TireTempRearAvg:  130,
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine := NewRuleEngine()
			engine.SetConfig(RoadRacingRuleConfig())
			observeRepeated(engine, tt.frame, 12)
			events := engine.Events()
			if len(events) == 0 {
				t.Fatalf("expected %s event", tt.eventType)
			}
			found := false
			for _, event := range events {
				if event.Type == tt.eventType {
					found = true
					break
				}
			}
			if !found {
				t.Fatalf("events = %#v, want %s", events, tt.eventType)
			}
		})
	}
}

func TestRoadRacingRuleConfigSeparatesUndersteerPhases(t *testing.T) {
	tests := []struct {
		name        string
		frame       NormalizedTelemetry
		want        string
		mustNotHave []string
	}{
		{
			name: "entry understeer on brake or coast",
			frame: baseRuleFrame(NormalizedTelemetry{
				SpeedKmh:             92,
				Steer01:              0.46,
				Throttle01:           0.04,
				Brake01:              0.03,
				FrontCombinedSlipAvg: 1.12,
				RearCombinedSlipAvg:  0.55,
			}),
			want:        "corner_entry_understeer",
			mustNotHave: []string{"mid_corner_understeer", "power_understeer"},
		},
		{
			name: "sustained corner understeer on maintenance throttle",
			frame: baseRuleFrame(NormalizedTelemetry{
				SpeedKmh:             95,
				Steer01:              0.48,
				Throttle01:           0.52,
				Brake01:              0.03,
				FrontCombinedSlipAvg: 1.12,
				RearCombinedSlipAvg:  0.55,
			}),
			want:        "mid_corner_understeer",
			mustNotHave: []string{"corner_entry_understeer", "power_understeer"},
		},
		{
			name: "power understeer on heavy throttle",
			frame: baseRuleFrame(NormalizedTelemetry{
				SpeedKmh:             95,
				Steer01:              0.48,
				Throttle01:           0.86,
				Brake01:              0.01,
				FrontCombinedSlipAvg: 1.16,
				RearCombinedSlipAvg:  0.55,
			}),
			want:        "power_understeer",
			mustNotHave: []string{"corner_entry_understeer", "mid_corner_understeer"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine := NewRuleEngine()
			engine.SetConfig(RoadRacingRuleConfig())
			observeRepeated(engine, tt.frame, 12)
			events := engine.Events()
			if !hasEventType(events, tt.want) {
				t.Fatalf("events = %#v, want %s", events, tt.want)
			}
			for _, unwanted := range tt.mustNotHave {
				if hasEventType(events, unwanted) {
					t.Fatalf("events = %#v, did not want %s", events, unwanted)
				}
			}
			for _, event := range events {
				if event.Type == tt.want && event.Evidence["slip_delta"] <= 0 {
					t.Fatalf("event evidence = %#v, want slip_delta", event.Evidence)
				}
			}
		})
	}
}

func TestRoadRacingRuleConfigDiffersFromDefault(t *testing.T) {
	defaults := DefaultRuleConfig()
	road := RoadRacingRuleConfig()
	if defaults.enabled("long_gear_bog_down") || defaults.enabled("top_speed_limited_by_gearing") {
		t.Fatal("road-only gearing events should be disabled in default config")
	}
	if !road.enabled("long_gear_bog_down") || !road.enabled("top_speed_limited_by_gearing") {
		t.Fatal("road-only gearing events should be enabled in road config")
	}
	for _, eventType := range []string{"high_speed_four_wheel_slide", "corner_entry_understeer", "corner_exit_oversteer"} {
		if road.threshold(eventType, "speedMin") <= defaults.threshold(eventType, "speedMin") {
			t.Fatalf("%s road speedMin = %.2f, default = %.2f", eventType, road.threshold(eventType, "speedMin"), defaults.threshold(eventType, "speedMin"))
		}
	}
}

func baseRuleFrame(overrides NormalizedTelemetry) NormalizedTelemetry {
	frame := NormalizedTelemetry{
		IsRaceOn:   true,
		SpeedKmh:   90,
		RpmRatio:   0.55,
		Gear:       3,
		Throttle01: 0.2,
		Brake01:    0,
		Steer01:    0,
	}

	if overrides.IsRaceOn {
		frame.IsRaceOn = overrides.IsRaceOn
	}
	if overrides.SpeedKmh != 0 {
		frame.SpeedKmh = overrides.SpeedKmh
	}
	if overrides.RpmRatio != 0 {
		frame.RpmRatio = overrides.RpmRatio
	}
	if overrides.Gear != 0 {
		frame.Gear = overrides.Gear
	}
	if overrides.Throttle01 != 0 {
		frame.Throttle01 = overrides.Throttle01
	}
	if overrides.Brake01 != 0 {
		frame.Brake01 = overrides.Brake01
	}
	if overrides.Steer01 != 0 {
		frame.Steer01 = overrides.Steer01
	}
	frame.FrontSlipRatioAvg = overrides.FrontSlipRatioAvg
	frame.RearSlipRatioAvg = overrides.RearSlipRatioAvg
	frame.FrontCombinedSlipAvg = overrides.FrontCombinedSlipAvg
	frame.RearCombinedSlipAvg = overrides.RearCombinedSlipAvg
	frame.TireTempFrontAvg = overrides.TireTempFrontAvg
	frame.TireTempRearAvg = overrides.TireTempRearAvg
	frame.SuspensionFrontAvg = overrides.SuspensionFrontAvg
	frame.SuspensionRearAvg = overrides.SuspensionRearAvg
	frame.YawRate = overrides.YawRate
	frame.PitchRate = overrides.PitchRate
	frame.RollRate = overrides.RollRate
	frame.WheelFL = overrides.WheelFL
	frame.WheelFR = overrides.WheelFR
	frame.WheelRL = overrides.WheelRL
	frame.WheelRR = overrides.WheelRR

	return frame
}

func hasEventType(events []DetectedEvent, eventType string) bool {
	for _, event := range events {
		if event.Type == eventType {
			return true
		}
	}
	return false
}

func longGearBogDownFrame() NormalizedTelemetry {
	return baseRuleFrame(NormalizedTelemetry{
		Gear:              4,
		SpeedKmh:          118,
		Throttle01:        0.92,
		RpmRatio:          0.30,
		FrontSlipRatioAvg: 0.08,
		RearSlipRatioAvg:  0.12,
	})
}

func topSpeedLimitedByGearingFrame() NormalizedTelemetry {
	return baseRuleFrame(NormalizedTelemetry{
		Gear:              6,
		SpeedKmh:          232,
		Throttle01:        0.94,
		RpmRatio:          0.985,
		FrontSlipRatioAvg: 0.05,
		RearSlipRatioAvg:  0.10,
	})
}

func observeRepeated(engine *RuleEngine, frame NormalizedTelemetry, count int) {
	for i := 0; i < count; i++ {
		next := frame
		next.TimeMS = int64(i * 100)
		engine.Observe(next)
	}
}
