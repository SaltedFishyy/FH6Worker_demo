package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"fh6worker/internal/advisor"
	"fh6worker/internal/storage"
	"fh6worker/internal/telemetry"
)

type App struct {
	ctx                   context.Context
	telemetry             *telemetry.Service
	store                 *storage.Store
	storeErr              error
	mu                    sync.Mutex
	activeSessionID       int64
	activeSessionStart    time.Time
	activeBaselineTrackID int64
	analysisMode          string
	tireSampleDir         string
	tuneWebServer         *TuneWebServer
	tuneHarvestCancel     *tuneHarvestCancellation
}

type tuneHarvestCancellation struct {
	cancel context.CancelFunc
}

const (
	analysisModeNone         = "none"
	analysisModeQuick        = "quick"
	analysisModeExpert       = "expert"
	analysisModeProfessional = "professional"
	analysisModeTireLab      = "tire_lab"
	analysisModeTrack        = "track_capture"
	analysisModeBaseline     = "track_baseline"
)

func NewApp() *App {
	store, err := storage.OpenDefault()
	if err == nil {
		err = store.CleanupLegacySessions()
	}
	return &App{telemetry: telemetry.NewService(), store: store, storeErr: err, tuneWebServer: NewTuneWebServer()}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) shutdown(ctx context.Context) {
	_ = a.StopTuneHarvest()
	_ = a.StopTelemetryReplay()
	_ = a.StopTelemetry()
	_ = a.StopTuneWebServer()
	_ = a.store.Close()
}

func (a *App) StartTelemetry(address string, port int) error {
	return a.StartTelemetryWithConditions(address, port, storage.DefaultTestConditions())
}

func (a *App) StartTelemetryWithConditions(address string, port int, conditions storage.TestConditions) error {
	if err := a.ensureStore(); err != nil {
		return err
	}
	a.mu.Lock()
	if a.activeSessionID != 0 {
		a.mu.Unlock()
		return fmt.Errorf("telemetry session is already active")
	}
	a.mu.Unlock()

	startedAt := time.Now().UTC()
	recordingPath, err := newRecordingPath(startedAt)
	if err != nil {
		return err
	}
	active, err := a.store.GetActiveTuneProfile()
	if err != nil {
		return err
	}
	if err := a.applyRuleConfig(active); err != nil {
		return err
	}
	session, err := a.createTelemetrySession(startedAt, recordingPath, active, sessionSnapshot(a.telemetry.Current(), active), conditions)
	if err != nil {
		return err
	}
	if err := a.telemetry.StartWithOptions(telemetry.StartOptions{
		Address:             address,
		Port:                port,
		RecordingPath:       recordingPath,
		RecordingLimitBytes: telemetry.DefaultRecordingLimit,
	}); err != nil {
		_ = a.store.DeleteTelemetrySession(session.ID)
		_ = os.Remove(recordingPath)
		return err
	}
	a.mu.Lock()
	a.activeSessionID = session.ID
	a.activeSessionStart = parseTimeOrNow(session.StartedAt)
	a.activeBaselineTrackID = 0
	a.analysisMode = analysisModeExpert
	a.mu.Unlock()
	return nil
}

func (a *App) StartQuickTelemetry(address string, port int) error {
	if err := a.ensureStore(); err != nil {
		return err
	}
	a.mu.Lock()
	if a.activeSessionID != 0 {
		a.mu.Unlock()
		return fmt.Errorf("expert telemetry session is already active")
	}
	a.mu.Unlock()
	active, err := a.store.GetActiveTuneProfile()
	if err != nil {
		return err
	}
	if err := a.applyRuleConfig(active); err != nil {
		return err
	}
	if err := a.telemetry.StartWithOptions(telemetry.StartOptions{
		Address: address,
		Port:    port,
	}); err != nil {
		a.mu.Lock()
		if a.analysisMode != analysisModeQuick {
			a.analysisMode = analysisModeNone
		}
		a.mu.Unlock()
		return err
	}
	a.mu.Lock()
	a.activeBaselineTrackID = 0
	a.analysisMode = analysisModeQuick
	a.mu.Unlock()
	return nil
}

func (a *App) StartProfessionalTelemetry(address string, port int) error {
	if err := a.ensureStore(); err != nil {
		return err
	}
	a.mu.Lock()
	if a.activeSessionID != 0 {
		a.mu.Unlock()
		return fmt.Errorf("legacy telemetry session is already active")
	}
	a.mu.Unlock()
	active, err := a.store.GetActiveTuneProfile()
	if err != nil {
		return err
	}
	if err := a.applyRuleConfig(active); err != nil {
		return err
	}
	if err := a.telemetry.StartWithOptions(telemetry.StartOptions{
		Address: address,
		Port:    port,
	}); err != nil {
		a.mu.Lock()
		if a.analysisMode != analysisModeProfessional {
			a.analysisMode = analysisModeNone
		}
		a.mu.Unlock()
		return err
	}
	a.mu.Lock()
	a.activeBaselineTrackID = 0
	a.analysisMode = analysisModeProfessional
	a.mu.Unlock()
	return nil
}

func (a *App) StartTireModelTelemetry(address string, port int) error {
	a.mu.Lock()
	if a.activeSessionID != 0 {
		a.mu.Unlock()
		return fmt.Errorf("expert telemetry session is already active")
	}
	a.mu.Unlock()
	if err := a.telemetry.StartWithOptions(telemetry.StartOptions{
		Address: address,
		Port:    port,
	}); err != nil {
		a.mu.Lock()
		if a.analysisMode != analysisModeTireLab {
			a.analysisMode = analysisModeNone
		}
		a.mu.Unlock()
		return err
	}
	a.mu.Lock()
	a.activeBaselineTrackID = 0
	a.analysisMode = analysisModeTireLab
	a.mu.Unlock()
	return nil
}

func (a *App) StartTrackCaptureTelemetry(address string, port int) error {
	a.mu.Lock()
	if a.activeSessionID != 0 {
		a.mu.Unlock()
		return fmt.Errorf("expert telemetry session is already active")
	}
	a.mu.Unlock()
	if err := a.telemetry.StartWithOptions(telemetry.StartOptions{
		Address: address,
		Port:    port,
	}); err != nil {
		a.mu.Lock()
		if a.analysisMode != analysisModeTrack {
			a.analysisMode = analysisModeNone
		}
		a.mu.Unlock()
		return err
	}
	a.mu.Lock()
	a.activeBaselineTrackID = 0
	a.analysisMode = analysisModeTrack
	a.mu.Unlock()
	return nil
}

func (a *App) StartTrackBaselineTelemetry(trackID int64, address string, port int) error {
	if err := a.ensureStore(); err != nil {
		return err
	}
	if trackID > 0 {
		if _, err := a.store.GetBenchmarkTrack(trackID); err != nil {
			return err
		}
	}
	a.mu.Lock()
	if a.activeSessionID != 0 {
		a.mu.Unlock()
		return fmt.Errorf("expert telemetry session is already active")
	}
	a.mu.Unlock()
	if err := a.telemetry.StartWithOptions(telemetry.StartOptions{
		Address: address,
		Port:    port,
	}); err != nil {
		a.mu.Lock()
		if a.analysisMode != analysisModeBaseline {
			a.analysisMode = analysisModeNone
			a.activeBaselineTrackID = 0
		}
		a.mu.Unlock()
		return err
	}
	a.mu.Lock()
	a.analysisMode = analysisModeBaseline
	a.activeBaselineTrackID = trackID
	a.mu.Unlock()
	return nil
}

func (a *App) StopTrackBaselineTelemetry() error {
	return a.StopTelemetry()
}

func (a *App) SaveTrackBaselineCapture(trackID int64) (*storage.TrackBaselineRun, error) {
	if err := a.ensureStore(); err != nil {
		return nil, err
	}
	a.mu.Lock()
	mode := a.analysisMode
	activeTrackID := a.activeBaselineTrackID
	a.mu.Unlock()
	if mode != analysisModeBaseline {
		return nil, fmt.Errorf("track baseline capture is not active")
	}
	if activeTrackID != 0 && activeTrackID != trackID {
		return nil, fmt.Errorf("track baseline capture target mismatch")
	}
	return a.store.SaveTrackBaselineCapture(trackID, a.telemetry.Samples(), a.telemetry.Events())
}

func (a *App) SaveTrackBaselineCaptureAuto(preferredTrackID int64, name string, trackType string) (*storage.TrackBaselineSaveResult, error) {
	if err := a.ensureStore(); err != nil {
		return nil, err
	}
	a.mu.Lock()
	mode := a.analysisMode
	a.mu.Unlock()
	if mode != analysisModeBaseline {
		return nil, fmt.Errorf("track baseline capture is not active")
	}
	return a.store.SaveTrackBaselineCaptureAuto(preferredTrackID, name, trackType, a.telemetry.Samples(), a.telemetry.Events())
}

func (a *App) DeleteTrackBaselineRun(id int64) error {
	if err := a.ensureStore(); err != nil {
		return err
	}
	return a.store.DeleteTrackBaselineRun(id)
}

func (a *App) GetTestConditionDefaults() (storage.TestConditions, error) {
	if err := a.ensureStore(); err != nil {
		return storage.TestConditions{}, err
	}
	return a.store.GetTestConditionDefaults()
}

func (a *App) SaveTestConditionDefaults(conditions storage.TestConditions) (storage.TestConditions, error) {
	if err := a.ensureStore(); err != nil {
		return storage.TestConditions{}, err
	}
	return a.store.SaveTestConditionDefaults(conditions)
}

func (a *App) StopTelemetry() error {
	a.mu.Lock()
	mode := a.analysisMode
	a.mu.Unlock()
	stopErr := a.telemetry.Stop()
	var finishErr error
	if mode == analysisModeExpert {
		finishErr = a.finishTelemetrySession()
	} else {
		a.mu.Lock()
		if a.analysisMode == analysisModeQuick || a.analysisMode == analysisModeProfessional || a.analysisMode == analysisModeTireLab || a.analysisMode == analysisModeTrack || a.analysisMode == analysisModeBaseline {
			// Keep the in-memory diagnostics available after stopping.
		} else {
			a.analysisMode = analysisModeNone
		}
		a.mu.Unlock()
	}
	if stopErr != nil {
		return stopErr
	}
	return finishErr
}

func (a *App) GetTelemetryStatus() telemetry.TelemetryStatus {
	status := a.telemetry.Status()
	a.mu.Lock()
	mode := a.analysisMode
	a.mu.Unlock()
	if mode == "" {
		mode = analysisModeNone
	}
	status.AnalysisMode = mode
	return status
}

func (a *App) GetCurrentTelemetry() *telemetry.NormalizedTelemetry {
	return a.telemetry.Current()
}

func (a *App) GetRecentTelemetry(seconds int) []telemetry.NormalizedTelemetry {
	return a.telemetry.Recent(seconds)
}

func (a *App) GetDetectedEvents() []telemetry.DetectedEvent {
	return a.telemetry.Events()
}

func (a *App) GetQuickDiagnostic() (*storage.QuickDiagnostic, error) {
	if err := a.ensureStore(); err != nil {
		return nil, err
	}
	active, _ := a.store.GetActiveTuneProfile()
	current := a.telemetry.Current()
	diagnostic := storage.BuildQuickDiagnostic(a.telemetry.Samples(), a.telemetry.Events(), current, active)
	return &diagnostic, nil
}

func (a *App) GetTireModelDiagnostic() storage.TireModelDiagnostic {
	return storage.BuildTireModelDiagnostic(a.telemetry.Samples(), a.telemetry.Current())
}

func (a *App) GetTireIssueAnalysis() storage.TireIssueAnalysis {
	return storage.BuildTireIssueAnalysis(a.telemetry.Samples())
}

func (a *App) GetTireIssueAdvice() storage.TireIssueAdvice {
	return storage.BuildTireIssueAdvice(a.telemetry.Samples())
}

func (a *App) GetTireDiagnosticSnapshot() storage.TireDiagnosticSnapshot {
	return storage.BuildTireDiagnosticSnapshotFromSamples(a.telemetry.Samples(), a.telemetry.Current())
}

func (a *App) ListTuningModelPipelines() storage.TuningPipelineCatalog {
	return storage.ListTuningModelPipelines()
}

func (a *App) RunTuningModelPipeline(input storage.TuningPipelineRunInput) (*storage.TuningPipelineRunResult, error) {
	if err := a.ensureStore(); err != nil {
		return nil, err
	}
	return a.store.RunTuningModelPipeline(input, a.telemetry.Samples())
}

func (a *App) GetProfessionalPipelineConfig() (storage.ProfessionalPipelineConfig, error) {
	if err := a.ensureStore(); err != nil {
		return storage.ProfessionalPipelineConfig{}, err
	}
	return a.store.GetProfessionalPipelineConfig()
}

func (a *App) SaveProfessionalPipelineConfig(input storage.ProfessionalPipelineConfig) (storage.ProfessionalPipelineConfig, error) {
	if err := a.ensureStore(); err != nil {
		return storage.ProfessionalPipelineConfig{}, err
	}
	return a.store.SaveProfessionalPipelineConfig(input)
}

func (a *App) GetProfessionalTuningDiagnostic() (*storage.ProfessionalTuningDiagnostic, error) {
	if err := a.ensureStore(); err != nil {
		return nil, err
	}
	config, err := a.store.GetProfessionalPipelineConfig()
	if err != nil {
		return nil, err
	}
	diagnostic := &storage.ProfessionalTuningDiagnostic{
		Status:    "ready",
		UpdatedAt: time.Now().UTC().Format(time.RFC3339),
		Config:    config,
	}
	warnings := storage.ValidateTuningPipelineCombination(storage.TuningSourceTireLabCurrent, config.DetectorID, config.DecisionerID, config.InterpreterID)
	if len(warnings) > 0 {
		diagnostic.Warnings = append(diagnostic.Warnings, warnings...)
	}
	result, err := a.store.RunTuningModelPipeline(storage.TuningPipelineRunInput{
		SourceType:    storage.TuningSourceTireLabCurrent,
		DetectorID:    config.DetectorID,
		DecisionerID:  config.DecisionerID,
		InterpreterID: config.InterpreterID,
	}, a.telemetry.Samples())
	if err != nil {
		return nil, err
	}
	diagnostic.Pipeline = result
	diagnostic.Warnings = append(diagnostic.Warnings, result.Warnings...)
	if result.Status != "" {
		diagnostic.Status = result.Status
	}
	return diagnostic, nil
}

func (a *App) CleanupLegacySessions() error {
	if err := a.ensureStore(); err != nil {
		return err
	}
	return a.store.CleanupLegacySessions()
}

func (a *App) SaveTireRegressionSample(input storage.TireRegressionSampleInput) (*storage.TireRegressionSample, error) {
	a.mu.Lock()
	mode := a.analysisMode
	a.mu.Unlock()
	if mode != analysisModeTireLab {
		return nil, fmt.Errorf("Tire Lab mode is required before saving a regression sample")
	}
	dir, err := a.tireRegressionSampleDir()
	if err != nil {
		return nil, err
	}
	return storage.SaveTireRegressionSample(dir, input, a.telemetry.Samples(), a.telemetry.Current())
}

func (a *App) ListTireRegressionSamples() ([]storage.TireRegressionSampleSummary, error) {
	dir, err := a.tireRegressionSampleDir()
	if err != nil {
		return nil, err
	}
	return storage.ListTireRegressionSamples(dir)
}

func (a *App) GetTireRegressionSample(id string) (*storage.TireRegressionSample, error) {
	dir, err := a.tireRegressionSampleDir()
	if err != nil {
		return nil, err
	}
	return storage.GetTireRegressionSample(dir, id)
}

func (a *App) UpdateTireRegressionSampleExpectation(id string, expected storage.TireRegressionExpectation) error {
	dir, err := a.tireRegressionSampleDir()
	if err != nil {
		return err
	}
	return storage.UpdateTireRegressionSampleExpectation(dir, id, expected)
}

func (a *App) DeleteTireRegressionSample(id string) error {
	dir, err := a.tireRegressionSampleDir()
	if err != nil {
		return err
	}
	return storage.DeleteTireRegressionSample(dir, id)
}

func (a *App) RunTireRegressionSample(id string) (*storage.TireRegressionResult, error) {
	dir, err := a.tireRegressionSampleDir()
	if err != nil {
		return nil, err
	}
	return storage.RunTireRegressionSample(dir, id)
}

func (a *App) RunAllTireRegressionSamples() ([]storage.TireRegressionResult, error) {
	dir, err := a.tireRegressionSampleDir()
	if err != nil {
		return nil, err
	}
	return storage.RunAllTireRegressionSamples(dir)
}

func (a *App) GetNetworkInterfaces() []telemetry.NetworkInterface {
	return a.telemetry.NetworkInterfaces()
}

func (a *App) ListTuneProfiles() ([]storage.TuneProfile, error) {
	if err := a.ensureStore(); err != nil {
		return nil, err
	}
	return a.store.ListTuneProfiles()
}

func (a *App) GetTuneProfile(id int64) (*storage.TuneProfile, error) {
	if err := a.ensureStore(); err != nil {
		return nil, err
	}
	return a.store.GetTuneProfile(id)
}

func (a *App) ListTuneProfilesForVehicle(carOrdinal int64, carClass string) ([]storage.TuneProfile, error) {
	if err := a.ensureStore(); err != nil {
		return nil, err
	}
	return a.store.ListTuneProfilesForVehicle(carOrdinal, carClass)
}

func (a *App) ListUpgradeUnlockRules() []storage.UpgradeUnlockRule {
	return storage.ListUpgradeUnlockRules()
}

func (a *App) ListTuneAdjustmentExplanations() []storage.TuneAdjustmentExplanation {
	return storage.ListTuneAdjustmentExplanations()
}

func (a *App) GetTuneToTireInfluenceMap() storage.TuneToTireInfluenceMap {
	return storage.GetTuneToTireInfluenceMap()
}

func (a *App) GenerateRoadStaticTuneBaseline(input storage.RoadStaticTuneBaselineInput) (*storage.RoadStaticTuneBaselineResult, error) {
	return storage.GenerateRoadStaticTuneBaseline(input)
}

func (a *App) StartTuneWebServer(port int) error {
	if a.tuneWebServer == nil {
		a.tuneWebServer = NewTuneWebServer()
	}
	return a.tuneWebServer.Start(port)
}

func (a *App) StopTuneWebServer() error {
	if a.tuneWebServer == nil {
		return nil
	}
	return a.tuneWebServer.Stop()
}

func (a *App) GetTuneWebServerStatus() TuneWebServerStatus {
	if a.tuneWebServer == nil {
		return TuneWebServerStatus{}
	}
	return a.tuneWebServer.Status()
}

func (a *App) ApplyRoadStaticTuneBaseline(input storage.RoadStaticTuneBaselineApplyInput) (*storage.RoadStaticTuneBaselineApplyResult, error) {
	if err := a.ensureStore(); err != nil {
		return nil, err
	}
	return a.store.ApplyRoadStaticTuneBaseline(input)
}

func (a *App) ExplainTuneFieldInfluence(fieldKey string) (*storage.TuneFieldInfluence, error) {
	item, err := storage.ExplainTuneFieldInfluence(fieldKey)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (a *App) CreateTuneProfile(input storage.TuneProfileInput) (*storage.TuneProfile, error) {
	if err := a.ensureStore(); err != nil {
		return nil, err
	}
	return a.store.CreateTuneProfile(input)
}

func (a *App) UpdateTuneProfile(id int64, input storage.TuneProfileInput) (*storage.TuneProfile, error) {
	if err := a.ensureStore(); err != nil {
		return nil, err
	}
	return a.store.UpdateTuneProfile(id, input)
}

func (a *App) ListTuneProfileSnapshots(profileID int64) ([]storage.TuneProfileSnapshot, error) {
	if err := a.ensureStore(); err != nil {
		return nil, err
	}
	return a.store.ListTuneProfileSnapshots(profileID)
}

func (a *App) RestoreTuneProfileSnapshot(snapshotID int64) (*storage.TuneProfile, error) {
	if err := a.ensureStore(); err != nil {
		return nil, err
	}
	return a.store.RestoreTuneProfileSnapshot(snapshotID)
}

func (a *App) DuplicateTuneProfile(id int64, versionName string) (*storage.TuneProfile, error) {
	if err := a.ensureStore(); err != nil {
		return nil, err
	}
	return a.store.DuplicateTuneProfile(id, versionName)
}

func (a *App) DeleteTuneProfile(id int64) error {
	if err := a.ensureStore(); err != nil {
		return err
	}
	return a.store.DeleteTuneProfile(id)
}

func (a *App) SetActiveTuneProfile(id int64) error {
	if err := a.ensureStore(); err != nil {
		return err
	}
	return a.store.SetActiveTuneProfile(id)
}

func (a *App) GetActiveTuneProfile() (*storage.TuneProfile, error) {
	if err := a.ensureStore(); err != nil {
		return nil, err
	}
	return a.store.GetActiveTuneProfile()
}

func (a *App) ListTuneProfileSessionStats() ([]storage.TuneProfileSessionStat, error) {
	if err := a.ensureStore(); err != nil {
		return nil, err
	}
	return a.store.ListTuneProfileSessionStats()
}

func (a *App) ListRuleThresholdProfiles() ([]storage.RuleThresholdProfile, error) {
	if err := a.ensureStore(); err != nil {
		return nil, err
	}
	return a.store.ListRuleThresholdProfiles()
}

func (a *App) ListStrategyTemplates() ([]storage.StrategyTemplate, error) {
	if err := a.ensureStore(); err != nil {
		return nil, err
	}
	return a.store.ListStrategyTemplates()
}

func (a *App) AnalyzeRoadStrategySessions(sessionIDs []int64, strategyTemplateID int64) (*storage.RoadStrategyAnalysis, error) {
	if err := a.ensureStore(); err != nil {
		return nil, err
	}
	return a.store.AnalyzeRoadStrategySessions(sessionIDs, strategyTemplateID)
}

func (a *App) CreateRuleThresholdProfile(input storage.RuleThresholdProfileInput) (*storage.RuleThresholdProfile, error) {
	if err := a.ensureStore(); err != nil {
		return nil, err
	}
	profile, err := a.store.CreateRuleThresholdProfile(input)
	if err == nil {
		_ = a.applyActiveRuleConfig()
	}
	return profile, err
}

func (a *App) UpdateRuleThresholdProfile(id int64, input storage.RuleThresholdProfileInput) (*storage.RuleThresholdProfile, error) {
	if err := a.ensureStore(); err != nil {
		return nil, err
	}
	profile, err := a.store.UpdateRuleThresholdProfile(id, input)
	if err == nil {
		_ = a.applyActiveRuleConfig()
	}
	return profile, err
}

func (a *App) DeleteRuleThresholdProfile(id int64) error {
	if err := a.ensureStore(); err != nil {
		return err
	}
	if err := a.store.DeleteRuleThresholdProfile(id); err != nil {
		return err
	}
	return a.applyActiveRuleConfig()
}

func (a *App) ResetRuleThresholdProfile(id int64) (*storage.RuleThresholdProfile, error) {
	if err := a.ensureStore(); err != nil {
		return nil, err
	}
	profile, err := a.store.ResetRuleThresholdProfile(id)
	if err == nil {
		_ = a.applyActiveRuleConfig()
	}
	return profile, err
}

func (a *App) ListTelemetrySessions(limit int) ([]storage.TelemetrySession, error) {
	if err := a.ensureStore(); err != nil {
		return nil, err
	}
	return a.store.ListTelemetrySessions(limit)
}

func (a *App) GetTelemetrySession(id int64) (*storage.TelemetrySession, error) {
	if err := a.ensureStore(); err != nil {
		return nil, err
	}
	return a.store.GetTelemetrySession(id)
}

func (a *App) BindTelemetrySessionTuneProfile(sessionID int64, tuneProfileID int64) (*storage.TelemetrySession, error) {
	if err := a.ensureStore(); err != nil {
		return nil, err
	}
	return a.store.BindTelemetrySessionTuneProfile(sessionID, tuneProfileID)
}

func (a *App) DeleteTelemetrySession(id int64) error {
	if err := a.ensureStore(); err != nil {
		return err
	}
	if a.telemetry.Status().Mode == "replay" {
		return fmt.Errorf("stop replay before deleting a telemetry session")
	}
	a.mu.Lock()
	activeSessionID := a.activeSessionID
	a.mu.Unlock()
	if activeSessionID == id {
		return fmt.Errorf("stop telemetry before deleting the active session")
	}
	session, err := a.store.GetTelemetrySession(id)
	if err != nil {
		return err
	}
	if err := a.store.DeleteTelemetrySession(id); err != nil {
		return err
	}
	return removeRecordingFile(session.RecordingPath)
}

func (a *App) GetSessionEvents(sessionID int64) ([]telemetry.DetectedEvent, error) {
	if err := a.ensureStore(); err != nil {
		return nil, err
	}
	return a.store.GetSessionEvents(sessionID)
}

func (a *App) GetSessionIssueSummary(sessionID int64) (*storage.SessionIssueSummary, error) {
	if err := a.ensureStore(); err != nil {
		return nil, err
	}
	return a.store.GetSessionIssueSummary(sessionID)
}

func (a *App) GetRoadTuningDecision(sessionID int64) (*storage.RoadTuningDecision, error) {
	if err := a.ensureStore(); err != nil {
		return nil, err
	}
	return a.store.GetRoadTuningDecision(sessionID)
}

func (a *App) ReloadTuningKnowledge() error {
	if err := a.ensureStore(); err != nil {
		return err
	}
	return a.store.ReloadTuningKnowledge()
}

func (a *App) GetRoadTuningKnowledgeStatus() storage.RoadTuningKnowledgeStatus {
	if err := a.ensureStore(); err != nil {
		return storage.RoadTuningKnowledgeStatus{LastError: err.Error()}
	}
	return a.store.RoadTuningKnowledgeStatus()
}

func (a *App) GetTunePlanDraft(sessionID int64) (*storage.TunePlanDraft, error) {
	if err := a.ensureStore(); err != nil {
		return nil, err
	}
	return a.store.GetTunePlanDraft(sessionID)
}

func (a *App) ApplyTunePlanDraft(input storage.TunePlanApplyInput) (*storage.TunePlanApplyResult, error) {
	if err := a.ensureStore(); err != nil {
		return nil, err
	}
	if a.telemetry.Status().Running || a.telemetry.ReplayStatus().Running {
		return nil, fmt.Errorf("stop telemetry or replay before applying a tune plan")
	}
	return a.store.ApplyTunePlanDraft(input)
}

func (a *App) GetRetestEvaluation(sessionID int64) (*storage.RetestEvaluation, error) {
	if err := a.ensureStore(); err != nil {
		return nil, err
	}
	return a.store.GetRetestEvaluation(sessionID)
}

func (a *App) GetSessionTelemetrySamples(sessionID int64, limit int) ([]telemetry.NormalizedTelemetry, error) {
	if err := a.ensureStore(); err != nil {
		return nil, err
	}
	return a.store.GetSessionTelemetrySamples(sessionID, limit)
}

func (a *App) ListBenchmarkTracks() ([]storage.BenchmarkTrack, error) {
	if err := a.ensureStore(); err != nil {
		return nil, err
	}
	return a.store.ListBenchmarkTracks()
}

func (a *App) ListTrackProfiles() ([]storage.TrackProfile, error) {
	if err := a.ensureStore(); err != nil {
		return nil, err
	}
	return a.store.ListTrackProfiles()
}

func (a *App) GetTrackProfile(trackID int64) (*storage.TrackProfile, error) {
	if err := a.ensureStore(); err != nil {
		return nil, err
	}
	return a.store.GetTrackProfile(trackID)
}

func (a *App) FindSimilarBenchmarkTracks(input storage.BenchmarkTrackInput) ([]storage.TrackMergeCandidate, error) {
	if err := a.ensureStore(); err != nil {
		return nil, err
	}
	return a.store.FindSimilarBenchmarkTracks(input)
}

func (a *App) MergeBenchmarkTrackInput(trackID int64, input storage.BenchmarkTrackInput) (*storage.BenchmarkTrack, error) {
	if err := a.ensureStore(); err != nil {
		return nil, err
	}
	return a.store.MergeBenchmarkTrackInput(trackID, input)
}

func (a *App) RenameBenchmarkTrack(trackID int64, name string) (*storage.BenchmarkTrack, error) {
	if err := a.ensureStore(); err != nil {
		return nil, err
	}
	return a.store.RenameBenchmarkTrack(trackID, name)
}

func (a *App) GetBenchmarkTrack(id int64) (*storage.BenchmarkTrack, error) {
	if err := a.ensureStore(); err != nil {
		return nil, err
	}
	return a.store.GetBenchmarkTrack(id)
}

func (a *App) CreateBenchmarkTrack(input storage.BenchmarkTrackInput) (*storage.BenchmarkTrack, error) {
	if err := a.ensureStore(); err != nil {
		return nil, err
	}
	return a.store.CreateBenchmarkTrack(input)
}

func (a *App) UpdateBenchmarkTrack(id int64, input storage.BenchmarkTrackInput) (*storage.BenchmarkTrack, error) {
	if err := a.ensureStore(); err != nil {
		return nil, err
	}
	return a.store.UpdateBenchmarkTrack(id, input)
}

func (a *App) DeleteBenchmarkTrack(id int64) error {
	if err := a.ensureStore(); err != nil {
		return err
	}
	return a.store.DeleteBenchmarkTrack(id)
}

func (a *App) CreateBenchmarkTrackFromSession(sessionID int64, name string) (*storage.BenchmarkTrack, error) {
	if err := a.ensureStore(); err != nil {
		return nil, err
	}
	return a.store.CreateBenchmarkTrackFromSession(sessionID, name)
}

func (a *App) ExtractBenchmarkTrackFromSession(input storage.BenchmarkTrackExtractionInput) (*storage.BenchmarkTrack, error) {
	if err := a.ensureStore(); err != nil {
		return nil, err
	}
	return a.store.ExtractBenchmarkTrackFromSession(input)
}

func (a *App) ReextractBenchmarkTrack(trackID int64, input storage.BenchmarkTrackExtractionInput) (*storage.BenchmarkTrack, error) {
	if err := a.ensureStore(); err != nil {
		return nil, err
	}
	return a.store.ReextractBenchmarkTrack(trackID, input)
}

func (a *App) AnalyzeSessionBenchmarkRuns(sessionID int64) ([]storage.BenchmarkRun, error) {
	if err := a.ensureStore(); err != nil {
		return nil, err
	}
	return a.store.AnalyzeSessionBenchmarkRuns(sessionID)
}

func (a *App) ListBenchmarkRuns(trackID int64, limit int) ([]storage.BenchmarkRun, error) {
	if err := a.ensureStore(); err != nil {
		return nil, err
	}
	return a.store.ListBenchmarkRuns(trackID, limit)
}

func (a *App) EvaluateRoadSession(sessionID int64) (*storage.RoadSessionEvaluation, error) {
	if err := a.ensureStore(); err != nil {
		return nil, err
	}
	return a.store.EvaluateRoadSession(sessionID)
}

func (a *App) CompareRoadEvaluations(leftID int64, rightID int64) (*storage.RoadEvaluationComparison, error) {
	if err := a.ensureStore(); err != nil {
		return nil, err
	}
	return a.store.CompareRoadEvaluations(leftID, rightID)
}

func (a *App) CompareTelemetrySessions(leftID int64, rightID int64) (storage.SessionComparison, error) {
	if err := a.ensureStore(); err != nil {
		return storage.SessionComparison{}, err
	}
	leftSession, err := a.store.GetTelemetrySession(leftID)
	if err != nil {
		return storage.SessionComparison{}, err
	}
	rightSession, err := a.store.GetTelemetrySession(rightID)
	if err != nil {
		return storage.SessionComparison{}, err
	}
	leftSamples, err := a.store.GetSessionTelemetrySamples(leftID, 10000)
	if err != nil {
		return storage.SessionComparison{}, err
	}
	rightSamples, err := a.store.GetSessionTelemetrySamples(rightID, 10000)
	if err != nil {
		return storage.SessionComparison{}, err
	}
	leftEvents, err := a.store.GetSessionEvents(leftID)
	if err != nil {
		return storage.SessionComparison{}, err
	}
	rightEvents, err := a.store.GetSessionEvents(rightID)
	if err != nil {
		return storage.SessionComparison{}, err
	}
	leftStats := summarizeComparison(leftSamples, leftEvents)
	rightStats := summarizeComparison(rightSamples, rightEvents)
	return storage.SessionComparison{
		LeftSession:           *leftSession,
		RightSession:          *rightSession,
		ComparabilityWarnings: comparabilityWarnings(*leftSession, *rightSession),
		Metrics: []storage.SessionComparisonMetric{
			metric("sample_count", "Samples", "", leftStats.sampleCount, rightStats.sampleCount, true),
			metric("event_count", "Events", "", leftStats.eventCount, rightStats.eventCount, false),
			metric("avg_speed", "Average speed", "km/h", leftStats.avgSpeed, rightStats.avgSpeed, true),
			metric("max_speed", "Max speed", "km/h", leftStats.maxSpeed, rightStats.maxSpeed, true),
			metric("max_rpm", "Max RPM", "rpm", leftStats.maxRpm, rightStats.maxRpm, false),
			metric("avg_throttle", "Average throttle", "%", leftStats.avgThrottle*100, rightStats.avgThrottle*100, false),
			metric("avg_brake", "Average brake", "%", leftStats.avgBrake*100, rightStats.avgBrake*100, false),
			metric("max_front_slip", "Max front combined slip", "", leftStats.maxFrontSlip, rightStats.maxFrontSlip, false),
			metric("max_rear_slip", "Max rear combined slip", "", leftStats.maxRearSlip, rightStats.maxRearSlip, false),
			metric("avg_tire_temp", "Average tire temp", "deg", leftStats.avgTireTemp, rightStats.avgTireTemp, false),
			metric("bottom_out_events", "Bottom-out events", "", leftStats.bottomOutEvents, rightStats.bottomOutEvents, false),
		},
		EventTypes: compareEventTypes(leftEvents, rightEvents),
	}, nil
}

func comparabilityWarnings(left storage.TelemetrySession, right storage.TelemetrySession) []string {
	warnings := make([]string, 0, 8)
	if telemetry.NormalizeGameMode(left.GameMode) != telemetry.NormalizeGameMode(right.GameMode) {
		warnings = append(warnings, "game_mode_mismatch")
	}
	leftConditions := storage.SessionTestConditions(left)
	rightConditions := storage.SessionTestConditions(right)
	if storage.TestConditionsContainUnknown(leftConditions) || storage.TestConditionsContainUnknown(rightConditions) {
		warnings = append(warnings, "test_conditions_unknown")
	}
	if leftConditions.DriverMode != rightConditions.DriverMode {
		warnings = append(warnings, "driver_mode_mismatch")
	}
	if leftConditions.BrakeAssist != rightConditions.BrakeAssist {
		warnings = append(warnings, "brake_assist_mismatch")
	}
	if leftConditions.SteeringAssist != rightConditions.SteeringAssist {
		warnings = append(warnings, "steering_assist_mismatch")
	}
	if leftConditions.TractionControl != rightConditions.TractionControl {
		warnings = append(warnings, "traction_control_mismatch")
	}
	if leftConditions.StabilityControl != rightConditions.StabilityControl {
		warnings = append(warnings, "stability_control_mismatch")
	}
	if leftConditions.Shifting != rightConditions.Shifting {
		warnings = append(warnings, "shifting_mismatch")
	}
	if leftConditions.LaunchControl != rightConditions.LaunchControl {
		warnings = append(warnings, "launch_control_mismatch")
	}
	return warnings
}

func (a *App) ReplayTelemetrySession(sessionID int64, speed float64) error {
	if err := a.ensureStore(); err != nil {
		return err
	}
	a.mu.Lock()
	if a.activeSessionID != 0 {
		a.mu.Unlock()
		return fmt.Errorf("telemetry session is already active")
	}
	a.mu.Unlock()
	session, err := a.store.GetTelemetrySession(sessionID)
	if err != nil {
		return err
	}
	if strings.TrimSpace(session.RecordingPath) == "" || session.RecordingPackets == 0 {
		return fmt.Errorf("session has no recording")
	}
	var profile *storage.TuneProfile
	if session.TuneProfileID != nil {
		profile, _ = a.store.GetTuneProfile(*session.TuneProfileID)
	}
	_, config, err := a.store.MatchRuleThresholdProfile(profile)
	if err != nil {
		return err
	}
	return a.telemetry.ReplayWithOptions(telemetry.ReplayOptions{SessionID: session.ID, Path: session.RecordingPath, Speed: speed, Config: config})
}

func (a *App) StopTelemetryReplay() error {
	return a.telemetry.StopReplay()
}

func (a *App) PauseTelemetryReplay() error {
	return a.telemetry.PauseReplay()
}

func (a *App) ResumeTelemetryReplay() error {
	return a.telemetry.ResumeReplay()
}

func (a *App) SeekTelemetryReplay(positionMS int64) error {
	return a.telemetry.SeekReplay(positionMS)
}

func (a *App) GetTelemetryReplayStatus() telemetry.TelemetryReplayStatus {
	return a.telemetry.ReplayStatus()
}

func (a *App) ResolveCarNameByOrdinal(carOrdinal int64) (string, error) {
	if err := a.ensureStore(); err != nil {
		return "", err
	}
	return a.store.ResolveCarNameByOrdinal(carOrdinal)
}

func (a *App) GenerateTuningReport(sessionID int64, language string) (string, error) {
	if err := a.ensureStore(); err != nil {
		return "", err
	}
	session, err := a.store.GetTelemetrySession(sessionID)
	if err != nil {
		return "", err
	}
	events, err := a.store.GetSessionEvents(sessionID)
	if err != nil {
		return "", err
	}
	var profile *storage.TuneProfile
	if snapshotProfile, err := storage.ParseTuneProfileSnapshotJSON(session.TuneSnapshotJSON); err == nil && snapshotProfile != nil {
		profile = snapshotProfile
	} else if session.TuneProfileID != nil {
		profile, _ = a.store.GetTuneProfile(*session.TuneProfileID)
	}
	evaluation, _ := a.store.EvaluateRoadSession(sessionID)
	issueSummary, _ := a.store.GetSessionIssueSummary(sessionID)
	decision, _ := a.store.GetRoadTuningDecision(sessionID)
	return advisor.GenerateTuningReportWithRoadDecision(*session, profile, events, language, evaluation, issueSummary, decision), nil
}

func (a *App) ensureStore() error {
	if a.storeErr != nil {
		return a.storeErr
	}
	if a.store == nil {
		return fmt.Errorf("storage is not available")
	}
	return nil
}

func (a *App) tireRegressionSampleDir() (string, error) {
	if strings.TrimSpace(a.tireSampleDir) != "" {
		if err := os.MkdirAll(a.tireSampleDir, 0755); err != nil {
			return "", err
		}
		return a.tireSampleDir, nil
	}
	return storage.DefaultTireRegressionSampleDir()
}

func (a *App) createTelemetrySession(startedAt time.Time, recordingPath string, active *storage.TuneProfile, snapshot storage.SessionVehicleSnapshot, conditions storage.TestConditions) (*storage.TelemetrySession, error) {
	var tuneProfileID *int64
	if active != nil {
		id := active.ID
		tuneProfileID = &id
	}
	tuneSnapshotJSON, err := storage.TuneProfileSnapshotJSON(active)
	if err != nil {
		return nil, err
	}
	input := storage.SessionStartInput{
		TuneProfileID:    tuneProfileID,
		TuneSnapshotJSON: tuneSnapshotJSON,
		SessionName:      "Session " + startedAt.Local().Format("2006-01-02 15:04:05"),
		Mode:             "Data Out",
		GameMode:         telemetry.GameModeUnknown,
		StartedAt:        startedAt.Format(time.RFC3339Nano),
		RecordingPath:    recordingPath,
		TestConditions:   storage.NormalizeTestConditions(conditions),
	}
	input.SessionVehicleSnapshot = snapshot
	return a.store.CreateTelemetrySession(input)
}

func (a *App) finishTelemetrySession() error {
	if err := a.ensureStore(); err != nil {
		return err
	}
	a.mu.Lock()
	sessionID := a.activeSessionID
	startedAt := a.activeSessionStart
	if sessionID == 0 {
		a.mu.Unlock()
		return nil
	}
	a.activeSessionID = 0
	a.activeSessionStart = time.Time{}
	a.mu.Unlock()

	endedAt := time.Now().UTC()
	if startedAt.IsZero() {
		startedAt = endedAt
	}
	duration := endedAt.Sub(startedAt).Milliseconds()
	if duration < 0 {
		duration = 0
	}
	summary := a.telemetry.Summary()
	var avgSpeed, maxSpeed *float64
	if summary.SampleCount > 0 {
		avg := summary.AvgSpeedKmh
		max := summary.MaxSpeedKmh
		avgSpeed = &avg
		maxSpeed = &max
	}
	recording := a.telemetry.Recording()
	samples := a.telemetry.Samples()
	gameMode := summarizeSessionGameMode(samples)
	if gameMode == telemetry.GameModeUnknown {
		if current := a.telemetry.Current(); current != nil {
			currentMode := telemetry.NormalizeGameMode(current.GameMode)
			if currentMode == telemetry.GameModeFreeRoam || currentMode == telemetry.GameModeRace {
				gameMode = currentMode
			}
		}
	}
	snapshot := sessionSnapshotFromSamples(samples)
	if !hasSessionSnapshot(snapshot) {
		snapshot = sessionSnapshot(a.telemetry.Current(), nil)
	}
	events := a.telemetry.Events()
	driverDetection := storage.DetectDriverMode(samples, events, gameMode)
	_, err := a.store.FinalizeTelemetrySession(storage.SessionFinalizeInput{
		SessionID:              sessionID,
		EndedAt:                endedAt.Format(time.RFC3339Nano),
		DurationMS:             duration,
		AvgSpeedKmh:            avgSpeed,
		MaxSpeedKmh:            maxSpeed,
		RecordingPackets:       recording.Packets,
		RecordingBytes:         recording.Bytes,
		RecordingTruncated:     recording.Truncated,
		GameMode:               gameMode,
		DriverModeDetection:    driverDetection,
		SessionVehicleSnapshot: snapshot,
	}, events, samples)
	a.mu.Lock()
	if a.analysisMode == analysisModeExpert {
		a.analysisMode = analysisModeNone
	}
	a.mu.Unlock()
	return err
}

func newRecordingPath(startedAt time.Time) (string, error) {
	dir, err := recordingDir()
	if err != nil {
		return "", err
	}
	name := "fh6_" + startedAt.Format("20060102_150405.000000000") + ".fh6udp"
	return filepath.Join(dir, name), nil
}

func recordingDir() (string, error) {
	if base, err := os.UserConfigDir(); err == nil {
		dir := filepath.Join(base, "FH6Worker", "recordings")
		if err := os.MkdirAll(dir, 0755); err == nil {
			return dir, nil
		}
	}
	dir := filepath.Join("data", "recordings")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	return dir, nil
}

func removeRecordingFile(path string) error {
	path = strings.TrimSpace(path)
	if path == "" {
		return nil
	}
	clean := filepath.Clean(path)
	base := filepath.Base(clean)
	if filepath.Ext(base) != ".fh6udp" || !strings.HasPrefix(base, "fh6_") {
		return nil
	}
	if err := os.Remove(clean); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func parseTimeOrNow(value string) time.Time {
	if parsed, err := time.Parse(time.RFC3339Nano, value); err == nil {
		return parsed
	}
	return time.Now().UTC()
}

func sessionSnapshot(frame *telemetry.NormalizedTelemetry, fallback *storage.TuneProfile) storage.SessionVehicleSnapshot {
	var snapshot storage.SessionVehicleSnapshot
	if frame != nil {
		snapshot = storage.SessionVehicleSnapshot{
			CarOrdinal:   positiveIntPtr(frame.CarOrdinal),
			CarClass:     strings.TrimSpace(frame.CarClass),
			CarPI:        positiveIntPtr(frame.CarPI),
			Drivetrain:   strings.TrimSpace(frame.Drivetrain),
			NumCylinders: positiveIntPtr(frame.NumCylinders),
		}
	}
	if fallback != nil {
		if snapshot.CarOrdinal == nil {
			snapshot.CarOrdinal = fallback.CarOrdinal
		}
		if strings.TrimSpace(snapshot.CarClass) == "" {
			snapshot.CarClass = fallback.CarClass
		}
		if snapshot.CarPI == nil {
			snapshot.CarPI = fallback.PI
		}
		if strings.TrimSpace(snapshot.Drivetrain) == "" {
			snapshot.Drivetrain = fallback.Drivetrain
		}
		if snapshot.NumCylinders == nil {
			snapshot.NumCylinders = fallback.NumCylinders
		}
	}
	return snapshot
}

func sessionSnapshotFromSamples(samples []telemetry.NormalizedTelemetry) storage.SessionVehicleSnapshot {
	for i := len(samples) - 1; i >= 0; i-- {
		snapshot := sessionSnapshot(&samples[i], nil)
		if hasSessionSnapshot(snapshot) {
			return snapshot
		}
	}
	return storage.SessionVehicleSnapshot{}
}

func hasSessionSnapshot(snapshot storage.SessionVehicleSnapshot) bool {
	return snapshot.CarOrdinal != nil || strings.TrimSpace(snapshot.CarClass) != "" || snapshot.CarPI != nil || strings.TrimSpace(snapshot.Drivetrain) != "" || snapshot.NumCylinders != nil
}

func summarizeSessionGameMode(samples []telemetry.NormalizedTelemetry) string {
	hasRace := false
	hasFreeRoam := false
	for _, sample := range samples {
		switch telemetry.NormalizeGameMode(sample.GameMode) {
		case telemetry.GameModeRace:
			hasRace = true
		case telemetry.GameModeFreeRoam:
			hasFreeRoam = true
		}
	}
	if hasRace && hasFreeRoam {
		return telemetry.GameModeMixed
	}
	if hasRace {
		return telemetry.GameModeRace
	}
	if hasFreeRoam {
		return telemetry.GameModeFreeRoam
	}
	return telemetry.GameModeUnknown
}

func positiveIntPtr(value int) *int64 {
	if value <= 0 {
		return nil
	}
	out := int64(value)
	return &out
}

func (a *App) applyActiveRuleConfig() error {
	active, err := a.store.GetActiveTuneProfile()
	if err != nil {
		return err
	}
	return a.applyRuleConfig(active)
}

func (a *App) applyRuleConfig(profile *storage.TuneProfile) error {
	_, config, err := a.store.MatchRuleThresholdProfile(profile)
	if err != nil {
		return err
	}
	a.telemetry.SetRuleConfig(config)
	return nil
}

type comparisonStats struct {
	sampleCount     float64
	eventCount      float64
	avgSpeed        float64
	maxSpeed        float64
	maxRpm          float64
	avgThrottle     float64
	avgBrake        float64
	maxFrontSlip    float64
	maxRearSlip     float64
	avgTireTemp     float64
	bottomOutEvents float64
}

func summarizeComparison(samples []telemetry.NormalizedTelemetry, events []telemetry.DetectedEvent) comparisonStats {
	stats := comparisonStats{sampleCount: float64(len(samples)), eventCount: float64(len(events))}
	var speedSum, throttleSum, brakeSum, tireTempSum float64
	for _, sample := range samples {
		speedSum += sample.SpeedKmh
		throttleSum += sample.Throttle01
		brakeSum += sample.Brake01
		tireTempSum += (sample.TireTempFrontAvg + sample.TireTempRearAvg) / 2
		if sample.SpeedKmh > stats.maxSpeed {
			stats.maxSpeed = sample.SpeedKmh
		}
		if sample.Rpm > stats.maxRpm {
			stats.maxRpm = sample.Rpm
		}
		if sample.FrontCombinedSlipAvg > stats.maxFrontSlip {
			stats.maxFrontSlip = sample.FrontCombinedSlipAvg
		}
		if sample.RearCombinedSlipAvg > stats.maxRearSlip {
			stats.maxRearSlip = sample.RearCombinedSlipAvg
		}
	}
	if len(samples) > 0 {
		n := float64(len(samples))
		stats.avgSpeed = speedSum / n
		stats.avgThrottle = throttleSum / n
		stats.avgBrake = brakeSum / n
		stats.avgTireTemp = tireTempSum / n
	}
	for _, event := range events {
		if event.Type == "suspension_bottom_out" {
			stats.bottomOutEvents++
		}
	}
	return stats
}

func metric(key, label, unit string, left, right float64, higherIsBetter bool) storage.SessionComparisonMetric {
	return storage.SessionComparisonMetric{Key: key, Label: label, Unit: unit, Left: left, Right: right, Delta: right - left, HigherIsBetter: higherIsBetter}
}

func compareEventTypes(leftEvents, rightEvents []telemetry.DetectedEvent) []storage.SessionEventComparison {
	left := map[string]int{}
	right := map[string]int{}
	keys := map[string]bool{}
	for _, event := range leftEvents {
		left[event.Type]++
		keys[event.Type] = true
	}
	for _, event := range rightEvents {
		right[event.Type]++
		keys[event.Type] = true
	}
	ordered := make([]string, 0, len(keys))
	for key := range keys {
		ordered = append(ordered, key)
	}
	sort.Strings(ordered)
	out := make([]storage.SessionEventComparison, 0, len(ordered))
	for _, key := range ordered {
		out = append(out, storage.SessionEventComparison{Type: key, Left: left[key], Right: right[key], Delta: right[key] - left[key]})
	}
	return out
}
