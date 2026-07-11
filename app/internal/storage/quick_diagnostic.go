package storage

import (
	"fmt"
	"sort"
	"strings"

	"fh6worker/internal/telemetry"
)

const (
	quickStatusNoData      = "no_data"
	quickStatusReady       = "ready"
	quickComparisonLap     = "lap_comparison"
	quickComparisonRolling = "rolling_window_only"
	quickSuggestionLimit   = 5
	quickLapMinSamples     = 5

	quickCompareYes     = "yes"
	quickCompareNo      = "no"
	quickCompareUnknown = "unknown"

	quickConfidenceHigh    = "high"
	quickConfidenceMedium  = "medium"
	quickConfidenceLow     = "low"
	quickConfidenceInvalid = "invalid"
)

func BuildQuickDiagnostic(samples []telemetry.NormalizedTelemetry, events []telemetry.DetectedEvent, current *telemetry.NormalizedTelemetry, profile *TuneProfile) QuickDiagnostic {
	diag := QuickDiagnostic{
		Status:           quickStatusNoData,
		ComparisonStatus: quickComparisonRolling,
		UpdatedAt:        nowText(),
		SampleCount:      len(samples),
		EventCount:       len(events),
		Vehicle:          quickVehicleSnapshot(samples, current),
		Comparability:    quickDefaultComparability(),
	}
	if len(samples) == 0 {
		return diag
	}
	ordered := sortedTelemetrySamples(samples)
	diag.Status = quickStatusReady
	diag.GameMode = quickSummarizeGameMode(ordered)
	if diag.GameMode == telemetry.GameModeUnknown && current != nil {
		diag.GameMode = telemetry.NormalizeGameMode(current.GameMode)
	}
	detection := DetectDriverMode(ordered, events, diag.GameMode)
	diag.DriverMode = detection.Mode
	diag.DriverModeConfidence = detection.Confidence

	lapContext, ok := quickLapContext(ordered, events, diag.GameMode)
	diag.Comparability = lapContext.comparability
	if ok {
		diag.ComparisonStatus = quickComparisonLap
		diag.CurrentLap = &lapContext.current
		diag.PreviousLap = &lapContext.previous
		groups := BuildSessionIssueGroups(lapContext.currentEvents, nil)
		baselineGroups := BuildSessionIssueGroups(lapContext.previousEvents, nil)
		applyIssueComparisons(groups, baselineGroups)
		applyIssueStrategies(groups, nil)
		diag.Groups = groups
		diag.GearPower = BuildGearPowerDiagnostic(lapContext.currentSamples, groups, profile)
	} else {
		diag.ComparisonStatus = quickComparisonRolling
		groups := BuildSessionIssueGroups(events, nil)
		applyIssueComparisons(groups, nil)
		applyIssueStrategies(groups, nil)
		diag.Groups = groups
		diag.GearPower = BuildGearPowerDiagnostic(ordered, groups, profile)
	}
	diag.Suggestions = BuildQuickAdviceSuggestions(diag.Groups, diag.GearPower, profile)
	diag.MissingProfileFields = QuickMissingFieldsFromSuggestions(diag.Suggestions)
	return diag
}

type quickLapData struct {
	lapNumber int
	samples   []telemetry.NormalizedTelemetry
}

type quickLapAnalysis struct {
	current         QuickLapSummary
	previous        QuickLapSummary
	currentSamples  []telemetry.NormalizedTelemetry
	previousSamples []telemetry.NormalizedTelemetry
	currentEvents   []telemetry.DetectedEvent
	previousEvents  []telemetry.DetectedEvent
	comparability   QuickComparability
}

func quickLapContext(samples []telemetry.NormalizedTelemetry, events []telemetry.DetectedEvent, gameMode string) (quickLapAnalysis, bool) {
	comparability := quickDefaultComparability()
	comparability.CurrentVehicle = quickVehicleSnapshot(samples, nil)
	if gameMode != telemetry.GameModeRace {
		comparability.Confidence = quickConfidenceLow
		comparability.Warnings = append(comparability.Warnings, "quick_non_race_track_unknown")
		return quickLapAnalysis{comparability: comparability}, false
	}
	lapsByNumber := map[int][]telemetry.NormalizedTelemetry{}
	for _, sample := range samples {
		if sample.LapNumber <= 0 {
			continue
		}
		lapsByNumber[sample.LapNumber] = append(lapsByNumber[sample.LapNumber], sample)
	}
	if len(lapsByNumber) < 2 {
		comparability.Confidence = quickConfidenceLow
		comparability.Warnings = append(comparability.Warnings, "quick_lap_data_insufficient")
		return quickLapAnalysis{comparability: comparability}, false
	}
	lapNumbers := make([]int, 0, len(lapsByNumber))
	for lap := range lapsByNumber {
		if len(lapsByNumber[lap]) >= quickLapMinSamples {
			lapNumbers = append(lapNumbers, lap)
		}
	}
	sort.Ints(lapNumbers)
	if len(lapNumbers) < 2 {
		comparability.Confidence = quickConfidenceLow
		comparability.Warnings = append(comparability.Warnings, "quick_lap_data_insufficient")
		return quickLapAnalysis{comparability: comparability}, false
	}
	currentLapNumber := lapNumbers[len(lapNumbers)-1]
	previousLapNumber := lapNumbers[len(lapNumbers)-2]
	currentSamples := sortedTelemetrySamples(lapsByNumber[currentLapNumber])
	previousSamples := sortedTelemetrySamples(lapsByNumber[previousLapNumber])
	comparability.BaselineVehicle = quickVehicleSnapshot(previousSamples, nil)
	comparability.CurrentVehicle = quickVehicleSnapshot(currentSamples, nil)
	comparability.SameVehicleClass = quickSameVehicleClass(comparability.BaselineVehicle, comparability.CurrentVehicle)
	if comparability.SameVehicleClass == quickCompareNo {
		comparability.Confidence = quickConfidenceInvalid
		comparability.Warnings = append(comparability.Warnings, "quick_vehicle_or_class_changed")
		return quickLapAnalysis{comparability: comparability}, false
	}
	if comparability.SameVehicleClass == quickCompareUnknown {
		comparability.Confidence = quickConfidenceLow
		comparability.Warnings = append(comparability.Warnings, "quick_vehicle_class_unknown")
		return quickLapAnalysis{comparability: comparability}, false
	}
	trackOK, trackWarnings := quickTrackContextComparable(previousSamples, currentSamples)
	comparability.Warnings = append(comparability.Warnings, trackWarnings...)
	if !trackOK {
		comparability.SameTrackContext = quickCompareUnknown
		comparability.Confidence = quickConfidenceLow
		if len(trackWarnings) == 0 {
			comparability.Warnings = append(comparability.Warnings, "quick_track_context_unknown")
		}
		return quickLapAnalysis{comparability: comparability}, false
	}
	comparability.SameTrackContext = quickCompareYes
	if len(trackWarnings) > 0 {
		comparability.Confidence = quickConfidenceMedium
	} else {
		comparability.Confidence = quickConfidenceHigh
	}
	currentDuration := quickLapDurationMS(currentSamples)
	previousComparable := quickComparableLapSamples(previousSamples, currentDuration)
	if len(previousComparable) >= quickLapMinSamples {
		previousSamples = previousComparable
	}
	currentEvents := quickEventsInWindow(events, currentSamples)
	previousEvents := quickEventsInWindow(events, previousSamples)
	currentGroups := BuildSessionIssueGroups(currentEvents, nil)
	previousGroups := BuildSessionIssueGroups(previousEvents, nil)
	return quickLapAnalysis{
		current:         quickLapSummary(currentLapNumber, currentSamples, currentEvents, currentGroups),
		previous:        quickLapSummary(previousLapNumber, previousSamples, previousEvents, previousGroups),
		currentSamples:  currentSamples,
		previousSamples: previousSamples,
		currentEvents:   currentEvents,
		previousEvents:  previousEvents,
		comparability:   comparability,
	}, true
}

func quickComparableLapSamples(samples []telemetry.NormalizedTelemetry, durationMS int64) []telemetry.NormalizedTelemetry {
	if durationMS <= 0 || len(samples) == 0 {
		return samples
	}
	out := make([]telemetry.NormalizedTelemetry, 0, len(samples))
	durationSeconds := float64(durationMS) / 1000
	startMS := samples[0].TimeMS
	for _, sample := range samples {
		if sample.CurrentLap > 0 {
			if sample.CurrentLap <= durationSeconds+0.25 {
				out = append(out, sample)
			}
			continue
		}
		if sample.TimeMS-startMS <= durationMS+250 {
			out = append(out, sample)
		}
	}
	if len(out) == 0 {
		return samples
	}
	return out
}

func quickDefaultComparability() QuickComparability {
	return QuickComparability{
		SameVehicleClass: quickCompareUnknown,
		SameTrackContext: quickCompareUnknown,
		Confidence:       quickConfidenceLow,
		Warnings:         []string{},
	}
}

func quickSameVehicleClass(baseline SessionVehicleSnapshot, current SessionVehicleSnapshot) string {
	if baseline.CarOrdinal == nil || current.CarOrdinal == nil || strings.TrimSpace(baseline.CarClass) == "" || strings.TrimSpace(current.CarClass) == "" {
		return quickCompareUnknown
	}
	if *baseline.CarOrdinal != *current.CarOrdinal || !strings.EqualFold(strings.TrimSpace(baseline.CarClass), strings.TrimSpace(current.CarClass)) {
		return quickCompareNo
	}
	return quickCompareYes
}

func quickTrackContextComparable(previous []telemetry.NormalizedTelemetry, current []telemetry.NormalizedTelemetry) (bool, []string) {
	warnings := []string{}
	if len(previous) < quickLapMinSamples || len(current) < quickLapMinSamples {
		return false, []string{"quick_lap_data_insufficient"}
	}
	if quickHasLapClockReset(previous) || quickHasLapClockReset(current) {
		return false, []string{"quick_lap_clock_reset"}
	}
	prevStart, prevEnd, prevOK := quickRaceTimeRange(previous)
	curStart, curEnd, curOK := quickRaceTimeRange(current)
	if prevOK && curOK {
		if curStart+0.25 < prevEnd || curEnd+0.25 < curStart || prevEnd+0.25 < prevStart {
			return false, []string{"quick_race_time_reset"}
		}
		return true, warnings
	}
	warnings = append(warnings, "quick_race_time_missing")
	return true, warnings
}

func quickHasLapClockReset(samples []telemetry.NormalizedTelemetry) bool {
	last := -1.0
	for _, sample := range samples {
		if sample.CurrentLap <= 0 {
			continue
		}
		if last >= 0 && sample.CurrentLap+0.5 < last {
			return true
		}
		last = sample.CurrentLap
	}
	return false
}

func quickRaceTimeRange(samples []telemetry.NormalizedTelemetry) (float64, float64, bool) {
	start := 0.0
	end := 0.0
	found := false
	for _, sample := range samples {
		if sample.CurrentRaceTime <= 0 {
			continue
		}
		if !found {
			start = sample.CurrentRaceTime
			end = sample.CurrentRaceTime
			found = true
			continue
		}
		if sample.CurrentRaceTime+0.5 < end {
			return start, end, false
		}
		end = sample.CurrentRaceTime
	}
	return start, end, found
}

func quickLapSummary(lapNumber int, samples []telemetry.NormalizedTelemetry, events []telemetry.DetectedEvent, groups []SessionIssueGroup) QuickLapSummary {
	summary := QuickLapSummary{
		LapNumber:   lapNumber,
		SampleCount: len(samples),
		DurationMS:  quickLapDurationMS(samples),
		EventCount:  len(events),
		IssueScore:  issueGroupScoreSum(groups),
	}
	if len(samples) == 0 {
		return summary
	}
	speedSum := 0.0
	for _, sample := range samples {
		speedSum += sample.SpeedKmh
		if sample.SpeedKmh > summary.MaxSpeedKmh {
			summary.MaxSpeedKmh = sample.SpeedKmh
		}
	}
	summary.AvgSpeedKmh = speedSum / float64(len(samples))
	return summary
}

func quickLapDurationMS(samples []telemetry.NormalizedTelemetry) int64 {
	if len(samples) == 0 {
		return 0
	}
	maxCurrentLap := 0.0
	for _, sample := range samples {
		if sample.CurrentLap > maxCurrentLap {
			maxCurrentLap = sample.CurrentLap
		}
	}
	if maxCurrentLap > 0 {
		return int64(maxCurrentLap * 1000)
	}
	duration := samples[len(samples)-1].TimeMS - samples[0].TimeMS
	if duration < 0 {
		return 0
	}
	return duration
}

func quickEventsInWindow(events []telemetry.DetectedEvent, samples []telemetry.NormalizedTelemetry) []telemetry.DetectedEvent {
	if len(events) == 0 || len(samples) == 0 {
		return nil
	}
	start := samples[0].TimeMS
	end := samples[len(samples)-1].TimeMS
	out := make([]telemetry.DetectedEvent, 0, len(events))
	for _, event := range events {
		if event.EndMS < start || event.StartMS > end {
			continue
		}
		out = append(out, event)
	}
	return out
}

func quickFieldKeysForAction(action WholeCarAdjustment, profile *TuneProfile) []string {
	if profile != nil {
		fields := tunePlanFields(action, profile)
		out := make([]string, 0, len(fields))
		for _, field := range fields {
			out = append(out, field.key)
		}
		return out
	}
	return quickFieldKeysForItem(action.Item, int(action.Evidence["gear"]+0.5), "")
}

func quickFieldKeysForItem(item string, gear int, drivetrain string) []string {
	switch item {
	case "front_tire_pressure":
		return []string{"frontTirePressure"}
	case "rear_tire_pressure":
		return []string{"rearTirePressure"}
	case "gear_1", "gear_2", "gear_3", "gear_4", "gear_5", "gear_6", "gear_7", "gear_8", "gear_9", "gear_10":
		return []string{"gear" + item[len("gear_"):]}
	case "current_gear":
		if gear >= 1 && gear <= 10 {
			return []string{fmt.Sprintf("gear%d", gear)}
		}
	case "final_drive":
		return []string{"finalDrive"}
	case "brake_balance":
		return []string{"brakeBalance"}
	case "brake_pressure":
		return []string{"brakePressure"}
	case "front_diff_accel":
		return []string{"frontDiffAccel"}
	case "front_diff_decel":
		return []string{"frontDiffDecel"}
	case "rear_diff_accel":
		return []string{"rearDiffAccel"}
	case "rear_diff_decel":
		return []string{"rearDiffDecel"}
	case "drive_diff_accel":
		return quickDrivenKeys(drivetrain, "frontDiffAccel", "rearDiffAccel")
	case "drive_diff_decel":
		return quickDrivenKeys(drivetrain, "frontDiffDecel", "rearDiffDecel")
	case "drive_tire_pressure":
		return quickDrivenKeys(drivetrain, "frontTirePressure", "rearTirePressure")
	case "tire_pressure":
		return []string{"frontTirePressure", "rearTirePressure"}
	case "front_arb":
		return []string{"frontArb"}
	case "rear_arb":
		return []string{"rearArb"}
	case "front_rebound":
		return []string{"frontRebound"}
	case "rear_rebound":
		return []string{"rearRebound"}
	case "front_camber":
		return []string{"frontCamber"}
	case "front_and_rear_aero":
		return []string{"frontAero", "rearAero"}
	case "ride_height":
		return []string{"frontRideHeight", "rearRideHeight"}
	case "spring_rate":
		return []string{"frontSpring", "rearSpring"}
	case "bump":
		return []string{"frontBump", "rearBump"}
	}
	return nil
}

func quickDrivenKeys(drivetrain string, front string, rear string) []string {
	switch strings.ToUpper(strings.TrimSpace(drivetrain)) {
	case "FWD":
		return []string{front}
	case "RWD":
		return []string{rear}
	case "AWD":
		return []string{front, rear}
	default:
		return nil
	}
}

func quickVehicleSnapshot(samples []telemetry.NormalizedTelemetry, current *telemetry.NormalizedTelemetry) SessionVehicleSnapshot {
	if current != nil && current.CarOrdinal > 0 {
		return quickSnapshot(current)
	}
	for i := len(samples) - 1; i >= 0; i-- {
		if samples[i].CarOrdinal > 0 {
			frame := samples[i]
			return quickSnapshot(&frame)
		}
	}
	return SessionVehicleSnapshot{}
}

func quickSnapshot(frame *telemetry.NormalizedTelemetry) SessionVehicleSnapshot {
	if frame == nil {
		return SessionVehicleSnapshot{}
	}
	var ordinal *int64
	if frame.CarOrdinal > 0 {
		value := int64(frame.CarOrdinal)
		ordinal = &value
	}
	var pi *int64
	if frame.CarPI > 0 {
		value := int64(frame.CarPI)
		pi = &value
	}
	var cylinders *int64
	if frame.NumCylinders > 0 {
		value := int64(frame.NumCylinders)
		cylinders = &value
	}
	return SessionVehicleSnapshot{
		CarOrdinal:   ordinal,
		CarClass:     frame.CarClass,
		CarPI:        pi,
		Drivetrain:   frame.Drivetrain,
		NumCylinders: cylinders,
	}
}

func quickSummarizeGameMode(samples []telemetry.NormalizedTelemetry) string {
	counts := map[string]int{}
	for _, sample := range samples {
		mode := telemetry.NormalizeGameMode(sample.GameMode)
		if mode == telemetry.GameModeUnknown {
			continue
		}
		counts[mode]++
	}
	if counts[telemetry.GameModeRace] > 0 && counts[telemetry.GameModeFreeRoam] > 0 {
		return "mixed"
	}
	bestMode := telemetry.GameModeUnknown
	bestCount := 0
	for mode, count := range counts {
		if count > bestCount {
			bestMode = mode
			bestCount = count
		}
	}
	return bestMode
}
