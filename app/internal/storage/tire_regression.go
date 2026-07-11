package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"fh6worker/internal/telemetry"
)

const (
	defaultTireRegressionWindowSeconds = 15
	minTireRegressionWindowSeconds     = 5
	maxTireRegressionWindowSeconds     = 30
)

type TireDiagnosticSnapshot struct {
	GeneratedAt   string                 `json:"generatedAt"`
	Status        string                 `json:"status"`
	SampleCount   int                    `json:"sampleCount"`
	WindowMS      int64                  `json:"windowMs"`
	Vehicle       SessionVehicleSnapshot `json:"vehicle"`
	DataQuality   TireSnapshotQuality    `json:"dataQuality"`
	Phase         TireSnapshotPhase      `json:"phase"`
	GripLimit     TireSnapshotGripLimit  `json:"gripLimit"`
	IssueAnalysis TireIssueAnalysis      `json:"issueAnalysis"`
	IssueAdvice   TireIssueAdvice        `json:"issueAdvice"`
	Risks         []string               `json:"risks"`
	Power         TireSnapshotSubsystem  `json:"power"`
	Brake         TireSnapshotSubsystem  `json:"brake"`
	Evidence      map[string]float64     `json:"evidence"`
}

type TireSnapshotQuality struct {
	Status             string   `json:"status"`
	Confidence         string   `json:"confidence"`
	SampleCount        int      `json:"sampleCount"`
	DynamicSampleCount int      `json:"dynamicSampleCount"`
	Reasons            []string `json:"reasons"`
}

type TireSnapshotPhase struct {
	Current     string  `json:"current"`
	Stable      string  `json:"stable"`
	Secondary   string  `json:"secondary"`
	Stability   string  `json:"stability"`
	Confidence  string  `json:"confidence"`
	ScoreMargin float64 `json:"scoreMargin"`
}

type TireSnapshotGripLimit struct {
	Type            string   `json:"type"`
	LimitedAxle     string   `json:"limitedAxle"`
	LimitedWheels   []string `json:"limitedWheels"`
	Confidence      string   `json:"confidence"`
	PrimaryEvidence string   `json:"primaryEvidence"`
	Reason          string   `json:"reason"`
}

type TireSnapshotSubsystem struct {
	Status      string             `json:"status"`
	Summary     string             `json:"summary"`
	Confidence  string             `json:"confidence"`
	Explanation string             `json:"explanation"`
	Evidence    map[string]float64 `json:"evidence"`
}

type TireRegressionExpectation struct {
	AllowedPhases      []string `json:"allowedPhases"`
	RequiredGripTypes  []string `json:"requiredGripTypes"`
	AllowedAxles       []string `json:"allowedAxles"`
	ForbiddenGripTypes []string `json:"forbiddenGripTypes"`
	MinDataQuality     string   `json:"minDataQuality"`
	Notes              string   `json:"notes"`
}

type TireRegressionSampleInput struct {
	Name          string                    `json:"name"`
	Scenario      string                    `json:"scenario"`
	WindowSeconds int                       `json:"windowSeconds"`
	Expected      TireRegressionExpectation `json:"expected"`
}

type TireRegressionSample struct {
	ID            string                          `json:"id"`
	Name          string                          `json:"name"`
	Scenario      string                          `json:"scenario"`
	CreatedAt     string                          `json:"createdAt"`
	WindowSeconds int                             `json:"windowSeconds"`
	Vehicle       SessionVehicleSnapshot          `json:"vehicle"`
	SampleCount   int                             `json:"sampleCount"`
	Samples       []telemetry.NormalizedTelemetry `json:"samples"`
	Snapshot      TireDiagnosticSnapshot          `json:"snapshot"`
	Expected      TireRegressionExpectation       `json:"expected"`
}

type TireRegressionSampleSummary struct {
	ID            string                    `json:"id"`
	Name          string                    `json:"name"`
	Scenario      string                    `json:"scenario"`
	CreatedAt     string                    `json:"createdAt"`
	WindowSeconds int                       `json:"windowSeconds"`
	Vehicle       SessionVehicleSnapshot    `json:"vehicle"`
	SampleCount   int                       `json:"sampleCount"`
	Expected      TireRegressionExpectation `json:"expected"`
}

type TireRegressionResult struct {
	SampleID string                    `json:"sampleId"`
	Name     string                    `json:"name"`
	Scenario string                    `json:"scenario"`
	Passed   bool                      `json:"passed"`
	Status   string                    `json:"status"`
	Failures []string                  `json:"failures"`
	Expected TireRegressionExpectation `json:"expected"`
	Actual   TireDiagnosticSnapshot    `json:"actual"`
}

func DefaultTireRegressionSampleDir() (string, error) {
	if base, err := os.UserConfigDir(); err == nil {
		dir := filepath.Join(base, "FH6Worker", "tire_model_samples")
		if err := os.MkdirAll(dir, 0755); err == nil {
			return dir, nil
		}
	}
	dir := filepath.Join("data", "tire_model_samples")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	return dir, nil
}

func BuildTireDiagnosticSnapshotFromSamples(samples []telemetry.NormalizedTelemetry, current *telemetry.NormalizedTelemetry) TireDiagnosticSnapshot {
	diag := BuildTireModelDiagnostic(samples, current)
	return BuildTireDiagnosticSnapshot(diag)
}

func BuildTireDiagnosticSnapshot(diag TireModelDiagnostic) TireDiagnosticSnapshot {
	risks := make([]string, 0, len(diag.Warnings))
	for _, warning := range diag.Warnings {
		if strings.TrimSpace(warning) != "" {
			risks = append(risks, warning)
		}
	}
	return TireDiagnosticSnapshot{
		GeneratedAt: time.Now().UTC().Format(time.RFC3339Nano),
		Status:      diag.Status,
		SampleCount: diag.SampleCount,
		WindowMS:    diag.WindowMS,
		Vehicle:     diag.Vehicle,
		DataQuality: TireSnapshotQuality{
			Status:             diag.DataQuality.Status,
			Confidence:         diag.DataQuality.Confidence,
			SampleCount:        diag.DataQuality.SampleCount,
			DynamicSampleCount: diag.DataQuality.DynamicSampleCount,
			Reasons:            append([]string(nil), diag.DataQuality.Reasons...),
		},
		Phase: TireSnapshotPhase{
			Current:     diag.PhaseDetail.CurrentPhase,
			Stable:      diag.PhaseDetail.StablePhase,
			Secondary:   diag.PhaseDetail.SecondaryPhase,
			Stability:   diag.PhaseDetail.PhaseStability,
			Confidence:  diag.PhaseDetail.Confidence,
			ScoreMargin: diag.PhaseDetail.ScoreMargin,
		},
		GripLimit: TireSnapshotGripLimit{
			Type:            diag.GripLimit.Type,
			LimitedAxle:     diag.GripLimit.LimitedAxle,
			LimitedWheels:   append([]string(nil), diag.GripLimit.LimitedWheels...),
			Confidence:      diag.GripLimit.Confidence,
			PrimaryEvidence: diag.GripLimit.PrimaryEvidence,
			Reason:          diag.GripLimit.Reason,
		},
		IssueAnalysis: diag.IssueAnalysis,
		IssueAdvice:   diag.IssueAdvice,
		Risks:         risks,
		Power: TireSnapshotSubsystem{
			Status:      diag.PowerToTire.Status,
			Summary:     diag.PowerToTire.Summary,
			Confidence:  diag.PowerToTire.Confidence,
			Explanation: diag.PowerToTire.Explanation,
			Evidence: map[string]float64{
				"driven_slip_ratio_p90": diag.PowerToTire.DrivenSlipRatioP90,
				"average_accel_g":       diag.PowerToTire.AverageAccelG,
				"peak_accel_g":          diag.PowerToTire.PeakAccelG,
				"high_throttle_samples": float64(diag.PowerToTire.HighThrottleSampleCount),
			},
		},
		Brake: TireSnapshotSubsystem{
			Status:      diag.BrakeToTire.Status,
			Summary:     diag.BrakeToTire.Summary,
			Confidence:  diag.BrakeToTire.Confidence,
			Explanation: diag.BrakeToTire.Explanation,
			Evidence: map[string]float64{
				"front_slip_ratio_p90": diag.BrakeToTire.FrontSlipRatioP90,
				"rear_slip_ratio_p90":  diag.BrakeToTire.RearSlipRatioP90,
				"average_decel_g":      diag.BrakeToTire.AverageDecelG,
				"peak_decel_g":         diag.BrakeToTire.PeakDecelG,
				"brake_samples":        float64(diag.BrakeToTire.BrakeSampleCount),
			},
		},
		Evidence: cloneFloatMap(diag.Evidence),
	}
}

func SaveTireRegressionSample(dir string, input TireRegressionSampleInput, samples []telemetry.NormalizedTelemetry, current *telemetry.NormalizedTelemetry) (*TireRegressionSample, error) {
	if strings.TrimSpace(dir) == "" {
		return nil, errors.New("sample directory is required")
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	windowSeconds := clampTireRegressionWindow(input.WindowSeconds)
	window := tireRegressionWindowSamples(samples, current, windowSeconds)
	if len(window) == 0 {
		return nil, errors.New("no Tire Lab telemetry samples available")
	}
	snapshot := BuildTireDiagnosticSnapshotFromSamples(window, &window[len(window)-1])
	expected := normalizeTireRegressionExpectation(input.Expected, snapshot)
	now := time.Now().UTC()
	name := strings.TrimSpace(input.Name)
	if name == "" {
		name = fmt.Sprintf("Tire sample %s", now.Format("20060102-150405"))
	}
	scenario := strings.TrimSpace(input.Scenario)
	if scenario == "" {
		scenario = "unclassified"
	}
	sample := &TireRegressionSample{
		ID:            tireRegressionSampleID(now, name),
		Name:          name,
		Scenario:      scenario,
		CreatedAt:     now.Format(time.RFC3339Nano),
		WindowSeconds: windowSeconds,
		Vehicle:       snapshot.Vehicle,
		SampleCount:   len(window),
		Samples:       window,
		Snapshot:      snapshot,
		Expected:      expected,
	}
	if err := writeTireRegressionSample(dir, sample); err != nil {
		return nil, err
	}
	return sample, nil
}

func ListTireRegressionSamples(dir string) ([]TireRegressionSampleSummary, error) {
	samples, err := loadAllTireRegressionSamples(dir)
	if err != nil {
		return nil, err
	}
	summaries := make([]TireRegressionSampleSummary, 0, len(samples))
	for _, sample := range samples {
		summaries = append(summaries, TireRegressionSampleSummary{
			ID:            sample.ID,
			Name:          sample.Name,
			Scenario:      sample.Scenario,
			CreatedAt:     sample.CreatedAt,
			WindowSeconds: sample.WindowSeconds,
			Vehicle:       sample.Vehicle,
			SampleCount:   sample.SampleCount,
			Expected:      sample.Expected,
		})
	}
	sort.SliceStable(summaries, func(i, j int) bool {
		return summaries[i].CreatedAt > summaries[j].CreatedAt
	})
	return summaries, nil
}

func GetTireRegressionSample(dir string, id string) (*TireRegressionSample, error) {
	return readTireRegressionSample(dir, id)
}

func UpdateTireRegressionSampleExpectation(dir string, id string, expected TireRegressionExpectation) error {
	sample, err := readTireRegressionSample(dir, id)
	if err != nil {
		return err
	}
	sample.Expected = normalizeTireRegressionExpectation(expected, sample.Snapshot)
	return writeTireRegressionSample(dir, sample)
}

func DeleteTireRegressionSample(dir string, id string) error {
	path, err := tireRegressionSamplePath(dir, id)
	if err != nil {
		return err
	}
	if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return nil
}

func RunTireRegressionSample(dir string, id string) (*TireRegressionResult, error) {
	sample, err := readTireRegressionSample(dir, id)
	if err != nil {
		return nil, err
	}
	return EvaluateTireRegressionSample(sample), nil
}

func RunAllTireRegressionSamples(dir string) ([]TireRegressionResult, error) {
	samples, err := loadAllTireRegressionSamples(dir)
	if err != nil {
		return nil, err
	}
	results := make([]TireRegressionResult, 0, len(samples))
	for _, sample := range samples {
		results = append(results, *EvaluateTireRegressionSample(sample))
	}
	return results, nil
}

func EvaluateTireRegressionSample(sample *TireRegressionSample) *TireRegressionResult {
	result := &TireRegressionResult{
		SampleID: sample.ID,
		Name:     sample.Name,
		Scenario: sample.Scenario,
		Expected: sample.Expected,
		Actual:   BuildTireDiagnosticSnapshotFromSamples(sample.Samples, nil),
	}
	failures := evaluateTireRegressionExpectation(sample.Expected, result.Actual)
	result.Failures = failures
	result.Passed = len(failures) == 0
	if result.Passed {
		result.Status = "passed"
	} else {
		result.Status = "failed"
	}
	return result
}

func evaluateTireRegressionExpectation(expected TireRegressionExpectation, actual TireDiagnosticSnapshot) []string {
	failures := []string{}
	if len(expected.AllowedPhases) > 0 && !containsAnyString(expected.AllowedPhases, actual.Phase.Current, actual.Phase.Stable) {
		failures = append(failures, "phase_mismatch")
	}
	if len(expected.RequiredGripTypes) > 0 && !containsStringValue(expected.RequiredGripTypes, actual.GripLimit.Type) {
		failures = append(failures, "required_grip_missing")
	}
	if len(expected.AllowedAxles) > 0 && !containsStringValue(expected.AllowedAxles, actual.GripLimit.LimitedAxle) {
		failures = append(failures, "limited_axle_mismatch")
	}
	if len(expected.ForbiddenGripTypes) > 0 && containsStringValue(expected.ForbiddenGripTypes, actual.GripLimit.Type) {
		failures = append(failures, "forbidden_grip_detected")
	}
	if expected.MinDataQuality != "" && tireDataQualityRank(actual.DataQuality.Status) < tireDataQualityRank(expected.MinDataQuality) {
		failures = append(failures, "data_quality_below_minimum")
	}
	return failures
}

func normalizeTireRegressionExpectation(expected TireRegressionExpectation, snapshot TireDiagnosticSnapshot) TireRegressionExpectation {
	if len(expected.AllowedPhases) == 0 {
		if snapshot.Phase.Current != "" && snapshot.Phase.Current != "unknown" {
			expected.AllowedPhases = []string{snapshot.Phase.Current}
		}
	}
	if len(expected.RequiredGripTypes) == 0 {
		if snapshot.GripLimit.Type != "" && snapshot.GripLimit.Type != "no_limit_detected" {
			expected.RequiredGripTypes = []string{snapshot.GripLimit.Type}
		}
	}
	if len(expected.AllowedAxles) == 0 {
		if snapshot.GripLimit.LimitedAxle != "" && snapshot.GripLimit.LimitedAxle != "none" {
			expected.AllowedAxles = []string{snapshot.GripLimit.LimitedAxle}
		}
	}
	if strings.TrimSpace(expected.MinDataQuality) == "" {
		if snapshot.DataQuality.Status == "invalid" {
			expected.MinDataQuality = "invalid"
		} else {
			expected.MinDataQuality = "low_confidence"
		}
	}
	expected.AllowedPhases = normalizeStringList(expected.AllowedPhases)
	expected.RequiredGripTypes = normalizeStringList(expected.RequiredGripTypes)
	expected.AllowedAxles = normalizeStringList(expected.AllowedAxles)
	expected.ForbiddenGripTypes = normalizeStringList(expected.ForbiddenGripTypes)
	expected.MinDataQuality = strings.TrimSpace(expected.MinDataQuality)
	return expected
}

func tireRegressionWindowSamples(samples []telemetry.NormalizedTelemetry, current *telemetry.NormalizedTelemetry, seconds int) []telemetry.NormalizedTelemetry {
	if len(samples) == 0 {
		if current == nil {
			return nil
		}
		return []telemetry.NormalizedTelemetry{*current}
	}
	out := make([]telemetry.NormalizedTelemetry, 0, len(samples))
	last := samples[len(samples)-1]
	if last.TimeMS > 0 {
		cutoff := last.TimeMS - int64(seconds*1000)
		for _, sample := range samples {
			if sample.TimeMS >= cutoff {
				out = append(out, sample)
			}
		}
	} else {
		count := seconds * 10
		if count <= 0 || count > len(samples) {
			count = len(samples)
		}
		out = append(out, samples[len(samples)-count:]...)
	}
	if len(out) == 0 {
		out = append(out, last)
	}
	return out
}

func clampTireRegressionWindow(seconds int) int {
	if seconds <= 0 {
		seconds = defaultTireRegressionWindowSeconds
	}
	if seconds < minTireRegressionWindowSeconds {
		return minTireRegressionWindowSeconds
	}
	if seconds > maxTireRegressionWindowSeconds {
		return maxTireRegressionWindowSeconds
	}
	return seconds
}

func loadAllTireRegressionSamples(dir string) ([]*TireRegressionSample, error) {
	if strings.TrimSpace(dir) == "" {
		return nil, errors.New("sample directory is required")
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	files, err := filepath.Glob(filepath.Join(dir, "*.json"))
	if err != nil {
		return nil, err
	}
	samples := make([]*TireRegressionSample, 0, len(files))
	for _, file := range files {
		raw, err := os.ReadFile(file)
		if err != nil {
			return nil, err
		}
		var sample TireRegressionSample
		if err := json.Unmarshal(raw, &sample); err != nil {
			return nil, err
		}
		if sample.ID == "" {
			sample.ID = strings.TrimSuffix(filepath.Base(file), filepath.Ext(file))
		}
		samples = append(samples, &sample)
	}
	sort.SliceStable(samples, func(i, j int) bool {
		return samples[i].CreatedAt > samples[j].CreatedAt
	})
	return samples, nil
}

func readTireRegressionSample(dir string, id string) (*TireRegressionSample, error) {
	path, err := tireRegressionSamplePath(dir, id)
	if err != nil {
		return nil, err
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var sample TireRegressionSample
	if err := json.Unmarshal(raw, &sample); err != nil {
		return nil, err
	}
	return &sample, nil
}

func writeTireRegressionSample(dir string, sample *TireRegressionSample) error {
	path, err := tireRegressionSamplePath(dir, sample.ID)
	if err != nil {
		return err
	}
	raw, err := json.MarshalIndent(sample, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, raw, 0644)
}

func tireRegressionSamplePath(dir string, id string) (string, error) {
	if strings.TrimSpace(dir) == "" {
		return "", errors.New("sample directory is required")
	}
	clean := sanitizeTireRegressionID(id)
	if clean == "" {
		return "", errors.New("sample id is required")
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	return filepath.Join(dir, clean+".json"), nil
}

func tireRegressionSampleID(now time.Time, name string) string {
	base := strings.ToLower(strings.TrimSpace(name))
	base = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(base, "-")
	base = strings.Trim(base, "-")
	if base == "" {
		base = "sample"
	}
	if len(base) > 40 {
		base = base[:40]
		base = strings.Trim(base, "-")
	}
	return fmt.Sprintf("%s-%d", base, now.UnixNano())
}

func sanitizeTireRegressionID(id string) string {
	clean := regexp.MustCompile(`[^A-Za-z0-9._-]+`).ReplaceAllString(strings.TrimSpace(id), "")
	return strings.Trim(clean, ".-_")
}

func normalizeStringList(values []string) []string {
	seen := map[string]bool{}
	out := []string{}
	for _, value := range values {
		normalized := strings.TrimSpace(value)
		if normalized == "" || seen[normalized] {
			continue
		}
		seen[normalized] = true
		out = append(out, normalized)
	}
	return out
}

func containsAnyString(values []string, candidates ...string) bool {
	for _, candidate := range candidates {
		if containsStringValue(values, candidate) {
			return true
		}
	}
	return false
}

func containsStringValue(values []string, target string) bool {
	target = strings.TrimSpace(target)
	for _, value := range values {
		if strings.TrimSpace(value) == target {
			return true
		}
	}
	return false
}

func tireDataQualityRank(status string) int {
	switch strings.TrimSpace(status) {
	case "valid":
		return 2
	case "low_confidence":
		return 1
	case "invalid":
		return 0
	default:
		return -1
	}
}

func cloneFloatMap(in map[string]float64) map[string]float64 {
	if len(in) == 0 {
		return map[string]float64{}
	}
	out := make(map[string]float64, len(in))
	for key, value := range in {
		out[key] = value
	}
	return out
}
