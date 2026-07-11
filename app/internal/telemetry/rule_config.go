package telemetry

type RuleConfig struct {
	Events map[string]RuleEventConfig `json:"events"`
}

type RuleEventConfig struct {
	Enabled       bool               `json:"enabled"`
	MinDurationMS int64              `json:"minDurationMs"`
	Thresholds    map[string]float64 `json:"thresholds"`
}

func DefaultRuleConfig() RuleConfig {
	return RuleConfig{Events: map[string]RuleEventConfig{
		"launch_wheelspin": {
			Enabled:       true,
			MinDurationMS: defaultEventDurationMS,
			Thresholds: map[string]float64{
				"throttleMin": 0.85, "speedMax": 60, "slipMin": 1.0, "severityMedium": 1.1, "severityHigh": 1.5,
			},
		},
		"launch_bog_down": {
			Enabled:       true,
			MinDurationMS: defaultEventDurationMS,
			Thresholds: map[string]float64{
				"throttleMin": 0.85, "speedMax": 50, "slipMax": 0.3, "rpmMax": 0.35, "severityMedium": 0.12, "severityHigh": 0.22,
			},
		},
		"short_gear": {
			Enabled:       true,
			MinDurationMS: defaultEventDurationMS,
			Thresholds: map[string]float64{
				"rpmMin": 0.95, "speedMax": 80, "throttleMin": 0.5, "severityMedium": 0.97, "severityHigh": 0.99,
			},
		},
		"long_gear_bog_down": {
			Enabled:       false,
			MinDurationMS: defaultEventDurationMS,
			Thresholds: map[string]float64{
				"gearMin": 2, "speedMin": 70, "speedMax": 180, "throttleMin": 0.75, "rpmMax": 0.45, "slipMax": 0.35, "severityMedium": 0.12, "severityHigh": 0.22,
			},
		},
		"top_speed_limited_by_gearing": {
			Enabled:       false,
			MinDurationMS: defaultEventDurationMS,
			Thresholds: map[string]float64{
				"gearMin": 4, "speedMin": 180, "throttleMin": 0.75, "rpmMin": 0.96, "slipMax": 0.60, "severityMedium": 0.97, "severityHigh": 0.99,
			},
		},
		"front_brake_lockup": {
			Enabled:       true,
			MinDurationMS: defaultEventDurationMS,
			Thresholds: map[string]float64{
				"brakeMin": 0.70, "speedMin": 80, "slipDeltaMin": 0.3, "frontCombinedMin": 1.0, "severityMedium": 1.1, "severityHigh": 1.4,
			},
		},
		"rear_brake_lockup": {
			Enabled:       true,
			MinDurationMS: defaultEventDurationMS,
			Thresholds: map[string]float64{
				"brakeMin": 0.70, "speedMin": 60, "slipDeltaMin": 0.3, "rearCombinedMin": 1.0, "yawSeverityFactor": 0.25, "severityMedium": 1.15, "severityHigh": 1.5,
			},
		},
		"corner_entry_understeer": {
			Enabled:       true,
			MinDurationMS: defaultEventDurationMS,
			Thresholds: map[string]float64{
				"speedMin": 60, "steerAbsMin": 0.35, "slipDeltaMin": 0.25, "throttleMax": 0.45, "brakeMin": 0.15, "severityMedium": 0.4, "severityHigh": 0.7,
			},
		},
		"mid_corner_understeer": {
			Enabled:       false,
			MinDurationMS: defaultEventDurationMS,
			Thresholds: map[string]float64{
				"speedMin": 70, "steerAbsMin": 0.35, "throttleMin": 0.20, "throttleMax": 0.65, "brakeMax": 0.15, "slipDeltaMin": 0.30, "severityMedium": 0.45, "severityHigh": 0.75,
			},
		},
		"corner_exit_oversteer": {
			Enabled:       true,
			MinDurationMS: defaultEventDurationMS,
			Thresholds: map[string]float64{
				"speedMin": 0, "throttleMin": 0.70, "steerAbsMin": 0.20, "slipDeltaMin": 0.25, "yawSeverityFactor": 0.2, "severityMedium": 0.45, "severityHigh": 0.75,
			},
		},
		"power_understeer": {
			Enabled:       false,
			MinDurationMS: defaultEventDurationMS,
			Thresholds: map[string]float64{
				"speedMin": 55, "throttleMin": 0.65, "steerAbsMin": 0.25, "slipDeltaMin": 0.30, "severityMedium": 0.45, "severityHigh": 0.75,
			},
		},
		"snap_oversteer": {
			Enabled:       false,
			MinDurationMS: 250,
			Thresholds: map[string]float64{
				"speedMin": 70, "steerAbsMin": 0.15, "slipDeltaMin": 0.35, "yawRateMin": 0.75, "severityMedium": 0.65, "severityHigh": 1.0,
			},
		},
		"high_speed_four_wheel_slide": {
			Enabled:       true,
			MinDurationMS: defaultEventDurationMS,
			Thresholds: map[string]float64{
				"speedMin": 160, "frontCombinedMin": 0.8, "rearCombinedMin": 0.8, "severityMedium": 0.95, "severityHigh": 1.2,
			},
		},
		"tire_overheat": {
			Enabled:       false,
			MinDurationMS: 700,
			Thresholds: map[string]float64{
				"speedMin": 60, "tempMin": 112, "severityMedium": 118, "severityHigh": 125,
			},
		},
		"tire_temp_imbalance": {
			Enabled:       false,
			MinDurationMS: 700,
			Thresholds: map[string]float64{
				"speedMin": 60, "tempDeltaMin": 18, "severityMedium": 22, "severityHigh": 30,
			},
		},
		"suspension_bottom_out": {
			Enabled:       true,
			MinDurationMS: 100,
			Thresholds: map[string]float64{
				"suspensionMin": 0.95, "severityMedium": 0.97, "severityHigh": 0.995,
			},
		},
	}}
}

func RoadRacingRuleConfig() RuleConfig {
	config := cloneRuleConfig(DefaultRuleConfig())
	setRuleEvent(config, "launch_wheelspin", true, 500, map[string]float64{
		"speedMax": 50, "slipMin": 1.1,
	})
	setRuleEvent(config, "launch_bog_down", true, 500, map[string]float64{
		"speedMax": 45, "rpmMax": 0.33,
	})
	setRuleEvent(config, "short_gear", true, 400, map[string]float64{
		"speedMax": 125, "throttleMin": 0.65,
	})
	setRuleEvent(config, "long_gear_bog_down", true, 400, map[string]float64{
		"gearMin": 2, "speedMin": 75, "speedMax": 190, "throttleMin": 0.75, "rpmMax": 0.43, "slipMax": 0.35, "severityMedium": 0.12, "severityHigh": 0.22,
	})
	setRuleEvent(config, "top_speed_limited_by_gearing", true, 500, map[string]float64{
		"gearMin": 4, "speedMin": 185, "throttleMin": 0.75, "rpmMin": 0.965, "slipMax": 0.60, "severityMedium": 0.975, "severityHigh": 0.99,
	})
	setRuleEvent(config, "front_brake_lockup", true, defaultEventDurationMS, map[string]float64{
		"speedMin": 100, "slipDeltaMin": 0.35, "frontCombinedMin": 1.1,
	})
	setRuleEvent(config, "rear_brake_lockup", true, defaultEventDurationMS, map[string]float64{
		"speedMin": 90, "slipDeltaMin": 0.35, "rearCombinedMin": 1.1,
	})
	setRuleEvent(config, "corner_entry_understeer", true, defaultEventDurationMS, map[string]float64{
		"speedMin": 80, "steerAbsMin": 0.30, "slipDeltaMin": 0.30, "throttleMax": 0.25, "brakeMin": 0.10,
	})
	setRuleEvent(config, "mid_corner_understeer", true, 500, map[string]float64{
		"speedMin": 75, "steerAbsMin": 0.32, "slipDeltaMin": 0.30,
	})
	setRuleEvent(config, "corner_exit_oversteer", true, defaultEventDurationMS, map[string]float64{
		"speedMin": 70, "steerAbsMin": 0.18, "slipDeltaMin": 0.30,
	})
	setRuleEvent(config, "power_understeer", true, 450, map[string]float64{
		"speedMin": 60, "throttleMin": 0.70, "steerAbsMin": 0.25, "slipDeltaMin": 0.32,
	})
	setRuleEvent(config, "snap_oversteer", true, 250, map[string]float64{
		"speedMin": 75, "yawRateMin": 0.80, "slipDeltaMin": 0.38,
	})
	setRuleEvent(config, "high_speed_four_wheel_slide", true, defaultEventDurationMS, map[string]float64{
		"speedMin": 185, "frontCombinedMin": 0.95, "rearCombinedMin": 0.95, "severityMedium": 1.1, "severityHigh": 1.35,
	})
	setRuleEvent(config, "tire_overheat", true, 900, map[string]float64{
		"speedMin": 70, "tempMin": 112, "severityMedium": 118, "severityHigh": 125,
	})
	setRuleEvent(config, "tire_temp_imbalance", true, 900, map[string]float64{
		"speedMin": 70, "tempDeltaMin": 18, "severityMedium": 22, "severityHigh": 30,
	})
	return NormalizeRuleConfig(config)
}

func NormalizeRuleConfig(config RuleConfig) RuleConfig {
	defaults := DefaultRuleConfig()
	if config.Events == nil {
		return defaults
	}
	for eventType, defaultEvent := range defaults.Events {
		event, ok := config.Events[eventType]
		if !ok {
			config.Events[eventType] = defaultEvent
			continue
		}
		if event.MinDurationMS <= 0 {
			event.MinDurationMS = defaultEvent.MinDurationMS
		}
		if event.Thresholds == nil {
			event.Thresholds = map[string]float64{}
		}
		for key, value := range defaultEvent.Thresholds {
			if _, ok := event.Thresholds[key]; !ok {
				event.Thresholds[key] = value
			}
		}
		config.Events[eventType] = event
	}
	return config
}

func cloneRuleConfig(config RuleConfig) RuleConfig {
	out := RuleConfig{Events: make(map[string]RuleEventConfig, len(config.Events))}
	for eventType, event := range config.Events {
		next := RuleEventConfig{
			Enabled:       event.Enabled,
			MinDurationMS: event.MinDurationMS,
			Thresholds:    make(map[string]float64, len(event.Thresholds)),
		}
		for key, value := range event.Thresholds {
			next.Thresholds[key] = value
		}
		out.Events[eventType] = next
	}
	return out
}

func setRuleEvent(config RuleConfig, eventType string, enabled bool, minDurationMS int64, thresholds map[string]float64) {
	event := config.Events[eventType]
	event.Enabled = enabled
	if minDurationMS > 0 {
		event.MinDurationMS = minDurationMS
	}
	if event.Thresholds == nil {
		event.Thresholds = map[string]float64{}
	}
	for key, value := range thresholds {
		event.Thresholds[key] = value
	}
	config.Events[eventType] = event
}

func (c RuleConfig) event(eventType string) RuleEventConfig {
	normalized := NormalizeRuleConfig(c)
	return normalized.Events[eventType]
}

func (c RuleConfig) threshold(eventType string, key string) float64 {
	event := c.event(eventType)
	return event.Thresholds[key]
}

func (c RuleConfig) minDuration(eventType string) int64 {
	return c.event(eventType).MinDurationMS
}

func (c RuleConfig) enabled(eventType string) bool {
	return c.event(eventType).Enabled
}
