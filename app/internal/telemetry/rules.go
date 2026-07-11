package telemetry

import (
	"fmt"
	"math"
	"sync"
)

const defaultEventDurationMS int64 = 300

type RuleEngine struct {
	mu     sync.RWMutex
	nextID int
	active map[string]*activeEvent
	events []DetectedEvent
	config RuleConfig
}

type activeEvent struct {
	eventType string
	startMS   int64
	endMS     int64
	emitted   bool
	index     int
	evidence  map[string]float64
}

type ruleEvaluation struct {
	eventType        string
	condition        bool
	segment          string
	minDurationMS    int64
	evidence         map[string]float64
	severity         string
	suggestedActions []SuggestedAction
}

func NewRuleEngine() *RuleEngine {
	return &RuleEngine{
		active: make(map[string]*activeEvent),
		config: DefaultRuleConfig(),
	}
}

func (e *RuleEngine) Reset() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.nextID = 0
	e.active = make(map[string]*activeEvent)
	e.events = nil
}

func (e *RuleEngine) SetConfig(config RuleConfig) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.config = NormalizeRuleConfig(config)
	e.active = make(map[string]*activeEvent)
}

func (e *RuleEngine) Events() []DetectedEvent {
	e.mu.RLock()
	defer e.mu.RUnlock()
	out := make([]DetectedEvent, len(e.events))
	copy(out, e.events)
	for i := range out {
		out[i].Evidence = cloneEvidence(out[i].Evidence)
		out[i].SuggestedActions = append([]SuggestedAction(nil), out[i].SuggestedActions...)
	}
	return out
}

func (e *RuleEngine) Observe(frame NormalizedTelemetry) []DetectedEvent {
	e.mu.Lock()
	defer e.mu.Unlock()

	evals := evaluateRules(frame, e.config)
	activeNow := make(map[string]bool, len(evals))
	var emitted []DetectedEvent
	for _, eval := range evals {
		if !eval.condition {
			continue
		}
		activeNow[eval.eventType] = true
		active := e.active[eval.eventType]
		if active == nil {
			active = &activeEvent{
				eventType: eval.eventType,
				startMS:   frame.TimeMS,
				endMS:     frame.TimeMS,
				index:     -1,
				evidence:  cloneEvidence(eval.evidence),
			}
			e.active[eval.eventType] = active
		}
		active.endMS = frame.TimeMS
		mergeEvidencePeaks(active.evidence, eval.evidence)
		duration := active.endMS - active.startMS
		if duration < 0 {
			duration = 0
		}
		if active.emitted {
			event := &e.events[active.index]
			event.EndMS = active.endMS
			event.DurationMS = duration
			event.Severity = maxSeverity(event.Severity, eval.severity)
			event.Evidence = cloneEvidence(active.evidence)
			continue
		}
		if duration >= eval.minDurationMS {
			e.nextID++
			event := DetectedEvent{
				ID:               fmt.Sprintf("evt-%d", e.nextID),
				Type:             eval.eventType,
				Severity:         eval.severity,
				StartMS:          active.startMS,
				EndMS:            active.endMS,
				DurationMS:       duration,
				Segment:          eval.segment,
				Evidence:         cloneEvidence(active.evidence),
				SuggestedActions: append([]SuggestedAction(nil), eval.suggestedActions...),
			}
			e.events = append(e.events, event)
			active.index = len(e.events) - 1
			active.emitted = true
			emitted = append(emitted, event)
		}
	}

	for eventType := range e.active {
		if !activeNow[eventType] {
			delete(e.active, eventType)
		}
	}

	return emitted
}

func evaluateRules(f NormalizedTelemetry, config RuleConfig) []ruleEvaluation {
	frontCombined := f.FrontCombinedSlipAvg
	rearCombined := f.RearCombinedSlipAvg
	maxCombined := math.Max(frontCombined, rearCombined)
	frontSlip := f.FrontSlipRatioAvg
	rearSlip := f.RearSlipRatioAvg
	maxSlip := math.Max(frontSlip, rearSlip)
	maxSuspension := max4(f.WheelFL.SuspensionTravel, f.WheelFR.SuspensionTravel, f.WheelRL.SuspensionTravel, f.WheelRR.SuspensionTravel)
	yawAbs := math.Abs(f.YawRate)
	steerAbs := math.Abs(f.Steer01)
	understeerDelta := frontCombined - rearCombined
	cornerOperation := cornerOperationStateCode(f)
	entryUndersteerCondition := f.SpeedKmh > config.threshold("corner_entry_understeer", "speedMin") &&
		steerAbs > config.threshold("corner_entry_understeer", "steerAbsMin") &&
		understeerDelta > config.threshold("corner_entry_understeer", "slipDeltaMin") &&
		(f.Throttle01 < config.threshold("corner_entry_understeer", "throttleMax") || f.Brake01 > config.threshold("corner_entry_understeer", "brakeMin"))
	powerUndersteerCondition := f.SpeedKmh > config.threshold("power_understeer", "speedMin") &&
		f.Throttle01 > config.threshold("power_understeer", "throttleMin") &&
		steerAbs > config.threshold("power_understeer", "steerAbsMin") &&
		understeerDelta > config.threshold("power_understeer", "slipDeltaMin")
	sustainedUndersteerCondition := f.SpeedKmh > config.threshold("mid_corner_understeer", "speedMin") &&
		steerAbs > config.threshold("mid_corner_understeer", "steerAbsMin") &&
		f.Throttle01 < config.threshold("mid_corner_understeer", "throttleMax") &&
		f.Brake01 < config.threshold("mid_corner_understeer", "brakeMax") &&
		understeerDelta > config.threshold("mid_corner_understeer", "slipDeltaMin") &&
		!entryUndersteerCondition &&
		!powerUndersteerCondition

	evals := []ruleEvaluation{
		{
			eventType:     "launch_wheelspin",
			condition:     f.Gear == 1 && f.Throttle01 > config.threshold("launch_wheelspin", "throttleMin") && f.SpeedKmh < config.threshold("launch_wheelspin", "speedMax") && maxSlip > config.threshold("launch_wheelspin", "slipMin"),
			segment:       "launch",
			minDurationMS: config.minDuration("launch_wheelspin"),
			evidence: evidence(f,
				"speed_kmh", f.SpeedKmh,
				"gear", float64(f.Gear),
				"throttle", f.Throttle01,
				"front_slip_ratio", frontSlip,
				"rear_slip_ratio", rearSlip,
				"rpm_ratio", f.RpmRatio,
			),
			severity: severity(maxSlip, config.threshold("launch_wheelspin", "severityMedium"), config.threshold("launch_wheelspin", "severityHigh")),
			suggestedActions: []SuggestedAction{
				action(0, "differential", "drive_diff_accel", "decrease", "3%-5%", "reduce driven-wheel slip"),
				action(1, "gearing", "gear_1", "decrease", "3%-5%", "reduce wheel torque during launch"),
				action(2, "tire", "drive_tire_pressure", "decrease", "0.03 BAR (≈0.5 PSI)", "increase launch traction"),
			},
		},
		{
			eventType:     "launch_bog_down",
			condition:     f.Gear == 1 && f.Throttle01 > config.threshold("launch_bog_down", "throttleMin") && f.SpeedKmh < config.threshold("launch_bog_down", "speedMax") && maxSlip < config.threshold("launch_bog_down", "slipMax") && f.RpmRatio < config.threshold("launch_bog_down", "rpmMax"),
			segment:       "launch",
			minDurationMS: config.minDuration("launch_bog_down"),
			evidence: evidence(f,
				"speed_kmh", f.SpeedKmh,
				"gear", float64(f.Gear),
				"throttle", f.Throttle01,
				"max_slip_ratio", maxSlip,
				"rpm_ratio", f.RpmRatio,
			),
			severity: severity(config.threshold("launch_bog_down", "rpmMax")-f.RpmRatio, config.threshold("launch_bog_down", "severityMedium"), config.threshold("launch_bog_down", "severityHigh")),
			suggestedActions: []SuggestedAction{
				action(0, "gearing", "gear_1", "increase", "3%-5%", "help the engine stay in the power band"),
				action(1, "gearing", "final_drive", "increase", "0.10-0.25", "shorten launch gearing"),
			},
		},
		{
			eventType:     "short_gear",
			condition:     f.RpmRatio > config.threshold("short_gear", "rpmMin") && f.SpeedKmh < config.threshold("short_gear", "speedMax") && f.Throttle01 > config.threshold("short_gear", "throttleMin"),
			segment:       "acceleration",
			minDurationMS: config.minDuration("short_gear"),
			evidence: evidence(f,
				"speed_kmh", f.SpeedKmh,
				"gear", float64(f.Gear),
				"throttle", f.Throttle01,
				"rpm_ratio", f.RpmRatio,
				"max_slip_ratio", maxSlip,
			),
			severity: severity(f.RpmRatio, config.threshold("short_gear", "severityMedium"), config.threshold("short_gear", "severityHigh")),
			suggestedActions: []SuggestedAction{
				action(0, "gearing", "current_gear", "decrease", "3%-8%", "avoid hitting the top of the gear too early"),
				action(1, "gearing", "final_drive", "decrease", "0.10-0.30", "lengthen all gears if multiple gears are short"),
			},
		},
		{
			eventType:     "long_gear_bog_down",
			condition:     float64(f.Gear) >= config.threshold("long_gear_bog_down", "gearMin") && f.SpeedKmh > config.threshold("long_gear_bog_down", "speedMin") && f.SpeedKmh < config.threshold("long_gear_bog_down", "speedMax") && f.Throttle01 > config.threshold("long_gear_bog_down", "throttleMin") && f.RpmRatio < config.threshold("long_gear_bog_down", "rpmMax") && maxSlip < config.threshold("long_gear_bog_down", "slipMax"),
			segment:       "acceleration",
			minDurationMS: config.minDuration("long_gear_bog_down"),
			evidence: evidence(f,
				"speed_kmh", f.SpeedKmh,
				"gear", float64(f.Gear),
				"throttle", f.Throttle01,
				"rpm_ratio", f.RpmRatio,
				"max_slip_ratio", maxSlip,
			),
			severity: severity(config.threshold("long_gear_bog_down", "rpmMax")-f.RpmRatio, config.threshold("long_gear_bog_down", "severityMedium"), config.threshold("long_gear_bog_down", "severityHigh")),
			suggestedActions: []SuggestedAction{
				action(0, "gearing", "current_gear", "increase", "3%-6%", "help the engine stay in the power band"),
				action(1, "gearing", "final_drive", "increase", "0.05-0.15", "shorten road acceleration gearing"),
				action(2, "differential", "drive_diff_accel", "check", "one small step", "verify traction is not limiting exit drive"),
			},
		},
		{
			eventType:     "top_speed_limited_by_gearing",
			condition:     float64(f.Gear) >= config.threshold("top_speed_limited_by_gearing", "gearMin") && f.SpeedKmh > config.threshold("top_speed_limited_by_gearing", "speedMin") && f.Throttle01 > config.threshold("top_speed_limited_by_gearing", "throttleMin") && f.RpmRatio > config.threshold("top_speed_limited_by_gearing", "rpmMin") && maxSlip < config.threshold("top_speed_limited_by_gearing", "slipMax"),
			segment:       "acceleration",
			minDurationMS: config.minDuration("top_speed_limited_by_gearing"),
			evidence: evidence(f,
				"speed_kmh", f.SpeedKmh,
				"gear", float64(f.Gear),
				"throttle", f.Throttle01,
				"rpm_ratio", f.RpmRatio,
				"max_slip_ratio", maxSlip,
			),
			severity: severity(f.RpmRatio, config.threshold("top_speed_limited_by_gearing", "severityMedium"), config.threshold("top_speed_limited_by_gearing", "severityHigh")),
			suggestedActions: []SuggestedAction{
				action(0, "gearing", "current_gear", "decrease", "3%-6%", "avoid hitting the top of the gear too early"),
				action(1, "gearing", "final_drive", "decrease", "0.05-0.15", "increase top speed headroom"),
				action(2, "aero", "front_and_rear_aero", "check", "one small step", "verify aero drag is not limiting top speed"),
			},
		},
		{
			eventType:     "front_brake_lockup",
			condition:     f.Brake01 > config.threshold("front_brake_lockup", "brakeMin") && f.SpeedKmh > config.threshold("front_brake_lockup", "speedMin") && frontCombined > rearCombined+config.threshold("front_brake_lockup", "slipDeltaMin") && frontCombined > config.threshold("front_brake_lockup", "frontCombinedMin"),
			segment:       "braking",
			minDurationMS: config.minDuration("front_brake_lockup"),
			evidence: evidence(f,
				"speed_kmh", f.SpeedKmh,
				"brake", f.Brake01,
				"front_combined_slip", frontCombined,
				"rear_combined_slip", rearCombined,
			),
			severity: severity(frontCombined, config.threshold("front_brake_lockup", "severityMedium"), config.threshold("front_brake_lockup", "severityHigh")),
			suggestedActions: []SuggestedAction{
				action(0, "brake", "brake_balance", "decrease", "1%-2% rearward", "reduce front lockup tendency"),
				action(1, "brake", "brake_pressure", "decrease", "3%-5%", "make threshold braking easier"),
			},
		},
		{
			eventType:     "rear_brake_lockup",
			condition:     f.Brake01 > config.threshold("rear_brake_lockup", "brakeMin") && f.SpeedKmh > config.threshold("rear_brake_lockup", "speedMin") && rearCombined > frontCombined+config.threshold("rear_brake_lockup", "slipDeltaMin") && rearCombined > config.threshold("rear_brake_lockup", "rearCombinedMin"),
			segment:       "braking",
			minDurationMS: config.minDuration("rear_brake_lockup"),
			evidence: evidence(f,
				"speed_kmh", f.SpeedKmh,
				"brake", f.Brake01,
				"front_combined_slip", frontCombined,
				"rear_combined_slip", rearCombined,
				"yaw_rate_abs", yawAbs,
			),
			severity: severity(rearCombined+yawAbs*config.threshold("rear_brake_lockup", "yawSeverityFactor"), config.threshold("rear_brake_lockup", "severityMedium"), config.threshold("rear_brake_lockup", "severityHigh")),
			suggestedActions: []SuggestedAction{
				action(0, "brake", "brake_balance", "increase", "1%-2% forward", "reduce rear lockup tendency"),
				action(1, "differential", "rear_diff_decel", "increase", "2%-3%", "stabilize the rear axle while braking"),
				action(2, "suspension", "rear_rebound", "decrease", "0.3-0.5", "improve rear compliance under braking"),
			},
		},
		{
			eventType:     "corner_entry_understeer",
			condition:     entryUndersteerCondition,
			segment:       "corner_entry",
			minDurationMS: config.minDuration("corner_entry_understeer"),
			evidence: evidence(f,
				"speed_kmh", f.SpeedKmh,
				"steer_abs", steerAbs,
				"throttle", f.Throttle01,
				"brake", f.Brake01,
				"front_combined_slip", frontCombined,
				"rear_combined_slip", rearCombined,
				"slip_delta", understeerDelta,
				"corner_operation_state", cornerOperation,
			),
			severity: severity(understeerDelta, config.threshold("corner_entry_understeer", "severityMedium"), config.threshold("corner_entry_understeer", "severityHigh")),
			suggestedActions: []SuggestedAction{
				action(0, "suspension", "front_arb", "decrease", "3%-5%", "increase front grip on entry"),
				action(1, "suspension", "front_rebound", "decrease", "0.3-0.5", "let the front tires load more smoothly"),
				action(2, "alignment", "front_camber", "check", "slightly more negative", "improve front tire contact in cornering"),
			},
		},
		{
			eventType:     "mid_corner_understeer",
			condition:     sustainedUndersteerCondition,
			segment:       "mid_corner",
			minDurationMS: config.minDuration("mid_corner_understeer"),
			evidence: evidence(f,
				"speed_kmh", f.SpeedKmh,
				"steer_abs", steerAbs,
				"throttle", f.Throttle01,
				"brake", f.Brake01,
				"front_combined_slip", frontCombined,
				"rear_combined_slip", rearCombined,
				"slip_delta", understeerDelta,
				"corner_operation_state", cornerOperation,
			),
			severity: severity(understeerDelta, config.threshold("mid_corner_understeer", "severityMedium"), config.threshold("mid_corner_understeer", "severityHigh")),
			suggestedActions: []SuggestedAction{
				action(0, "suspension", "front_arb", "decrease", "0.5-1.0", "increase front grip on steady cornering"),
				action(1, "suspension", "rear_arb", "increase", "0.5-1.0", "rotate the car more in steady cornering"),
				action(2, "alignment", "front_camber", "check", "slightly more negative", "improve front tire contact in cornering"),
			},
		},
		{
			eventType:     "corner_exit_oversteer",
			condition:     f.SpeedKmh > config.threshold("corner_exit_oversteer", "speedMin") && f.Throttle01 > config.threshold("corner_exit_oversteer", "throttleMin") && steerAbs > config.threshold("corner_exit_oversteer", "steerAbsMin") && rearCombined > frontCombined+config.threshold("corner_exit_oversteer", "slipDeltaMin"),
			segment:       "corner_exit",
			minDurationMS: config.minDuration("corner_exit_oversteer"),
			evidence: evidence(f,
				"speed_kmh", f.SpeedKmh,
				"gear", float64(f.Gear),
				"steer_abs", steerAbs,
				"throttle", f.Throttle01,
				"front_combined_slip", frontCombined,
				"rear_combined_slip", rearCombined,
				"yaw_rate_abs", yawAbs,
			),
			severity: severity(rearCombined-frontCombined+yawAbs*config.threshold("corner_exit_oversteer", "yawSeverityFactor"), config.threshold("corner_exit_oversteer", "severityMedium"), config.threshold("corner_exit_oversteer", "severityHigh")),
			suggestedActions: []SuggestedAction{
				action(0, "differential", "rear_diff_accel", "decrease", "3%-5%", "reduce power oversteer"),
				action(1, "gearing", "current_gear", "decrease", "3%-5%", "reduce wheel torque on exit"),
				action(2, "suspension", "rear_arb", "decrease", "3%-5%", "increase rear grip"),
			},
		},
		{
			eventType:     "power_understeer",
			condition:     powerUndersteerCondition,
			segment:       "corner_exit",
			minDurationMS: config.minDuration("power_understeer"),
			evidence: evidence(f,
				"speed_kmh", f.SpeedKmh,
				"steer_abs", steerAbs,
				"throttle", f.Throttle01,
				"brake", f.Brake01,
				"front_combined_slip", frontCombined,
				"rear_combined_slip", rearCombined,
				"slip_delta", understeerDelta,
				"corner_operation_state", cornerOperation,
			),
			severity: severity(understeerDelta, config.threshold("power_understeer", "severityMedium"), config.threshold("power_understeer", "severityHigh")),
			suggestedActions: []SuggestedAction{
				action(0, "differential", "drive_diff_accel", "decrease", "2%-4%", "reduce power-on understeer"),
				action(1, "suspension", "front_arb", "decrease", "0.5-1.0", "increase front grip under power"),
				action(2, "gearing", "current_gear", "decrease", "2%-3%", "reduce wheel torque on exit"),
			},
		},
		{
			eventType:     "snap_oversteer",
			condition:     f.SpeedKmh > config.threshold("snap_oversteer", "speedMin") && steerAbs > config.threshold("snap_oversteer", "steerAbsMin") && rearCombined > frontCombined+config.threshold("snap_oversteer", "slipDeltaMin") && yawAbs > config.threshold("snap_oversteer", "yawRateMin"),
			segment:       "cornering",
			minDurationMS: config.minDuration("snap_oversteer"),
			evidence: evidence(f,
				"speed_kmh", f.SpeedKmh,
				"steer_abs", steerAbs,
				"front_combined_slip", frontCombined,
				"rear_combined_slip", rearCombined,
				"yaw_rate_abs", yawAbs,
			),
			severity: severity(rearCombined-frontCombined+yawAbs*0.25, config.threshold("snap_oversteer", "severityMedium"), config.threshold("snap_oversteer", "severityHigh")),
			suggestedActions: []SuggestedAction{
				action(0, "suspension", "rear_rebound", "decrease", "0.3-0.5", "make rear response less abrupt"),
				action(1, "suspension", "rear_arb", "decrease", "0.5-1.0", "increase rear grip"),
				action(2, "differential", "rear_diff_decel", "increase", "2%-3%", "stabilize the rear axle while off throttle"),
			},
		},
		{
			eventType:     "high_speed_four_wheel_slide",
			condition:     f.SpeedKmh > config.threshold("high_speed_four_wheel_slide", "speedMin") && frontCombined > config.threshold("high_speed_four_wheel_slide", "frontCombinedMin") && rearCombined > config.threshold("high_speed_four_wheel_slide", "rearCombinedMin"),
			segment:       "high_speed_corner",
			minDurationMS: config.minDuration("high_speed_four_wheel_slide"),
			evidence: evidence(f,
				"speed_kmh", f.SpeedKmh,
				"front_combined_slip", frontCombined,
				"rear_combined_slip", rearCombined,
				"yaw_rate_abs", yawAbs,
			),
			severity: severity(maxCombined, config.threshold("high_speed_four_wheel_slide", "severityMedium"), config.threshold("high_speed_four_wheel_slide", "severityHigh")),
			suggestedActions: []SuggestedAction{
				action(0, "aero", "front_and_rear_aero", "increase", "one small step", "increase high-speed grip"),
				action(1, "tire", "tire_pressure", "decrease", "0.03 BAR (≈0.5 PSI)", "increase tire contact patch"),
				action(2, "suspension", "ride_height", "check", "avoid bottoming", "prevent aero and suspension instability"),
			},
		},
		{
			eventType:     "tire_overheat",
			condition:     f.SpeedKmh > config.threshold("tire_overheat", "speedMin") && math.Max(f.TireTempFrontAvg, f.TireTempRearAvg) > config.threshold("tire_overheat", "tempMin"),
			segment:       "tire",
			minDurationMS: config.minDuration("tire_overheat"),
			evidence: evidence(f,
				"speed_kmh", f.SpeedKmh,
				"front_tire_temp", f.TireTempFrontAvg,
				"rear_tire_temp", f.TireTempRearAvg,
				"front_combined_slip", frontCombined,
				"rear_combined_slip", rearCombined,
			),
			severity: severity(math.Max(f.TireTempFrontAvg, f.TireTempRearAvg), config.threshold("tire_overheat", "severityMedium"), config.threshold("tire_overheat", "severityHigh")),
			suggestedActions: []SuggestedAction{
				action(0, "tire", "tire_pressure", "decrease", "0.03 BAR (≈0.5 PSI)", "reduce tire overheating tendency"),
				action(1, "alignment", "front_camber", "check", "slightly more negative", "improve front tire contact in cornering"),
				action(2, "suspension", "front_arb", "decrease", "0.5-1.0", "reduce sustained tire scrub"),
			},
		},
		{
			eventType:     "tire_temp_imbalance",
			condition:     f.SpeedKmh > config.threshold("tire_temp_imbalance", "speedMin") && math.Abs(f.TireTempFrontAvg-f.TireTempRearAvg) > config.threshold("tire_temp_imbalance", "tempDeltaMin"),
			segment:       "tire",
			minDurationMS: config.minDuration("tire_temp_imbalance"),
			evidence: evidence(f,
				"speed_kmh", f.SpeedKmh,
				"front_tire_temp", f.TireTempFrontAvg,
				"rear_tire_temp", f.TireTempRearAvg,
				"tire_temp_delta", math.Abs(f.TireTempFrontAvg-f.TireTempRearAvg),
			),
			severity: severity(math.Abs(f.TireTempFrontAvg-f.TireTempRearAvg), config.threshold("tire_temp_imbalance", "severityMedium"), config.threshold("tire_temp_imbalance", "severityHigh")),
			suggestedActions: []SuggestedAction{
				action(0, "tire", "tire_pressure", "check", "0.02-0.03 BAR", "balance front and rear tire temperatures"),
				action(1, "suspension", "front_arb", "check", "one small step", "rebalance axle load transfer"),
				action(2, "suspension", "rear_arb", "check", "one small step", "rebalance axle load transfer"),
			},
		},
		{
			eventType:     "suspension_bottom_out",
			condition:     maxSuspension > config.threshold("suspension_bottom_out", "suspensionMin"),
			segment:       "suspension",
			minDurationMS: config.minDuration("suspension_bottom_out"),
			evidence: evidence(f,
				"speed_kmh", f.SpeedKmh,
				"max_suspension_travel", maxSuspension,
				"front_suspension", f.SuspensionFrontAvg,
				"rear_suspension", f.SuspensionRearAvg,
				"pitch_rate_abs", math.Abs(f.PitchRate),
				"roll_rate_abs", math.Abs(f.RollRate),
			),
			severity: severity(maxSuspension, config.threshold("suspension_bottom_out", "severityMedium"), config.threshold("suspension_bottom_out", "severityHigh")),
			suggestedActions: []SuggestedAction{
				action(0, "suspension", "ride_height", "increase", "one small step", "restore suspension travel"),
				action(1, "suspension", "spring_rate", "increase", "3%-5%", "reduce bottoming frequency"),
				action(2, "damping", "bump", "increase", "0.2-0.4", "support compression on impacts"),
			},
		},
	}
	filtered := evals[:0]
	for _, eval := range evals {
		if config.enabled(eval.eventType) {
			filtered = append(filtered, eval)
		}
	}
	return filtered
}

func action(priority int, category, item, direction, amount, reason string) SuggestedAction {
	return SuggestedAction{
		Priority:  priority,
		Category:  category,
		Item:      item,
		Direction: direction,
		Amount:    amount,
		Reason:    reason,
	}
}

func NormalizeSuggestedActions(eventType string, actions []SuggestedAction) []SuggestedAction {
	out := make([]SuggestedAction, len(actions))
	for i, action := range actions {
		out[i] = normalizeSuggestedAction(eventType, action)
	}
	return out
}

func cornerOperationStateCode(f NormalizedTelemetry) float64 {
	switch {
	case f.Brake01 >= 0.12:
		return 1
	case f.Throttle01 <= 0.12:
		return 2
	case f.Throttle01 < 0.70:
		return 3
	default:
		return 4
	}
}

func normalizeSuggestedAction(eventType string, action SuggestedAction) SuggestedAction {
	switch eventType {
	case "launch_wheelspin":
		if action.Item == "gear_1" {
			action.Direction = "decrease"
			if action.Amount == "5%-10%" {
				action.Amount = "3%-5%"
			}
		}
	case "launch_bog_down":
		if action.Item == "gear_1" {
			action.Direction = "increase"
			if action.Amount == "5%-10%" {
				action.Amount = "3%-5%"
			}
		}
	case "short_gear", "top_speed_limited_by_gearing":
		if action.Item == "current_gear" {
			action.Direction = "decrease"
		}
	case "long_gear_bog_down":
		if action.Item == "current_gear" {
			action.Direction = "increase"
		}
	case "corner_exit_oversteer":
		if action.Item == "current_gear" && action.Reason == "reduce wheel torque on exit" {
			action.Direction = "decrease"
		}
	}
	return action
}

func evidence(_ NormalizedTelemetry, pairs ...any) map[string]float64 {
	out := make(map[string]float64, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		key, ok := pairs[i].(string)
		if !ok {
			continue
		}
		switch value := pairs[i+1].(type) {
		case float64:
			out[key] = value
		case int:
			out[key] = float64(value)
		}
	}
	return out
}

func cloneEvidence(in map[string]float64) map[string]float64 {
	if in == nil {
		return nil
	}
	out := make(map[string]float64, len(in))
	for key, value := range in {
		out[key] = value
	}
	return out
}

func mergeEvidencePeaks(dst, src map[string]float64) {
	for key, value := range src {
		if old, ok := dst[key]; !ok || math.Abs(value) > math.Abs(old) {
			dst[key] = value
		}
	}
}

func severity(value, medium, high float64) string {
	if value >= high {
		return "high"
	}
	if value >= medium {
		return "medium"
	}
	return "low"
}

func maxSeverity(a, b string) string {
	rank := map[string]int{"low": 0, "medium": 1, "high": 2}
	if rank[b] > rank[a] {
		return b
	}
	return a
}

func max4(a, b, c, d float64) float64 {
	return math.Max(math.Max(a, b), math.Max(c, d))
}
