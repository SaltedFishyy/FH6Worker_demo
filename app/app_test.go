package main

import (
	"net"
	"path/filepath"
	"strings"
	"testing"

	"fh6worker/internal/storage"
	"fh6worker/internal/telemetry"
)

func TestStartQuickTelemetryDoesNotCreateSession(t *testing.T) {
	store, err := storage.Open(filepath.Join(t.TempDir(), "quick.db"))
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	defer store.Close()
	app := &App{telemetry: telemetry.NewService(), store: store}
	port := freeUDPPort(t)
	if err := app.StartQuickTelemetry("127.0.0.1", port); err != nil {
		t.Fatalf("start quick: %v", err)
	}
	status := app.GetTelemetryStatus()
	if !status.Running || status.AnalysisMode != analysisModeQuick || status.RecordingActive || status.RecordingPackets != 0 {
		t.Fatalf("status = %#v, want quick mode without recording", status)
	}
	if err := app.StopTelemetry(); err != nil {
		t.Fatalf("stop quick: %v", err)
	}
	sessions, err := store.ListTelemetrySessions(10)
	if err != nil {
		t.Fatalf("list sessions: %v", err)
	}
	if len(sessions) != 0 {
		t.Fatalf("sessions = %#v, want no persisted session in quick mode", sessions)
	}
	status = app.GetTelemetryStatus()
	if status.AnalysisMode != analysisModeQuick {
		t.Fatalf("analysis mode after stop = %q, want retained quick diagnostic mode", status.AnalysisMode)
	}
}

func TestStartTireModelTelemetryDoesNotCreateSession(t *testing.T) {
	store, err := storage.Open(filepath.Join(t.TempDir(), "tirelab.db"))
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	defer store.Close()
	app := &App{telemetry: telemetry.NewService(), store: store}
	if err := startTireModelOnFreePort(t, app); err != nil {
		t.Fatalf("start tire model: %v", err)
	}
	status := app.GetTelemetryStatus()
	if !status.Running || status.AnalysisMode != analysisModeTireLab || status.RecordingActive || status.RecordingPackets != 0 {
		t.Fatalf("status = %#v, want tire lab mode without recording", status)
	}
	if err := app.StopTelemetry(); err != nil {
		t.Fatalf("stop tire model: %v", err)
	}
	sessions, err := store.ListTelemetrySessions(10)
	if err != nil {
		t.Fatalf("list sessions: %v", err)
	}
	if len(sessions) != 0 {
		t.Fatalf("sessions = %#v, want no persisted session in tire model mode", sessions)
	}
	status = app.GetTelemetryStatus()
	if status.AnalysisMode != analysisModeTireLab {
		t.Fatalf("analysis mode after stop = %q, want retained tire lab diagnostic mode", status.AnalysisMode)
	}
}

func TestStartTrackCaptureTelemetryDoesNotCreateSession(t *testing.T) {
	store, err := storage.Open(filepath.Join(t.TempDir(), "trackcapture.db"))
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	defer store.Close()
	app := &App{telemetry: telemetry.NewService(), store: store}
	if err := startTrackCaptureOnFreePort(t, app); err != nil {
		t.Fatalf("start track capture: %v", err)
	}
	status := app.GetTelemetryStatus()
	if !status.Running || status.AnalysisMode != analysisModeTrack || status.RecordingActive || status.RecordingPackets != 0 {
		t.Fatalf("status = %#v, want track capture mode without recording", status)
	}
	if err := app.StopTelemetry(); err != nil {
		t.Fatalf("stop track capture: %v", err)
	}
	sessions, err := store.ListTelemetrySessions(10)
	if err != nil {
		t.Fatalf("list sessions: %v", err)
	}
	if len(sessions) != 0 {
		t.Fatalf("sessions = %#v, want no persisted session in track capture mode", sessions)
	}
	status = app.GetTelemetryStatus()
	if status.AnalysisMode != analysisModeTrack {
		t.Fatalf("analysis mode after stop = %q, want retained track capture mode", status.AnalysisMode)
	}
}

func TestStartTrackBaselineTelemetryDoesNotCreateSession(t *testing.T) {
	store, err := storage.Open(filepath.Join(t.TempDir(), "trackbaseline.db"))
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	defer store.Close()
	track, err := store.CreateBenchmarkTrack(storage.BenchmarkTrackInput{
		Name:      "Baseline Track",
		TrackType: "sprint",
		Polyline:  []storage.BenchmarkPoint{{X: 0, Z: 0}, {X: 100, Z: 0}},
	})
	if err != nil {
		t.Fatalf("create track: %v", err)
	}
	app := &App{telemetry: telemetry.NewService(), store: store}
	if err := startTrackBaselineOnFreePort(t, app, track.ID); err != nil {
		t.Fatalf("start track baseline: %v", err)
	}
	status := app.GetTelemetryStatus()
	if !status.Running || status.AnalysisMode != analysisModeBaseline || status.RecordingActive || status.RecordingPackets != 0 {
		t.Fatalf("status = %#v, want track baseline mode without recording", status)
	}
	if err := app.StopTrackBaselineTelemetry(); err != nil {
		t.Fatalf("stop track baseline: %v", err)
	}
	sessions, err := store.ListTelemetrySessions(10)
	if err != nil {
		t.Fatalf("list sessions: %v", err)
	}
	if len(sessions) != 0 {
		t.Fatalf("sessions = %#v, want no persisted session in track baseline mode", sessions)
	}
	status = app.GetTelemetryStatus()
	if status.AnalysisMode != analysisModeBaseline {
		t.Fatalf("analysis mode after stop = %q, want retained track baseline mode", status.AnalysisMode)
	}
}

func startTireModelOnFreePort(t *testing.T, app *App) error {
	t.Helper()
	var lastErr error
	for i := 0; i < 20; i++ {
		err := app.StartTireModelTelemetry("127.0.0.1", freeUDPPort(t))
		if err == nil {
			return nil
		}
		lastErr = err
		if !strings.Contains(err.Error(), "Only one usage") && !strings.Contains(err.Error(), "address already in use") {
			return err
		}
	}
	return lastErr
}

func startTrackCaptureOnFreePort(t *testing.T, app *App) error {
	t.Helper()
	var lastErr error
	for i := 0; i < 20; i++ {
		err := app.StartTrackCaptureTelemetry("127.0.0.1", freeUDPPort(t))
		if err == nil {
			return nil
		}
		lastErr = err
		if !strings.Contains(err.Error(), "Only one usage") && !strings.Contains(err.Error(), "address already in use") {
			return err
		}
	}
	return lastErr
}

func startTrackBaselineOnFreePort(t *testing.T, app *App, trackID int64) error {
	t.Helper()
	var lastErr error
	for i := 0; i < 20; i++ {
		err := app.StartTrackBaselineTelemetry(trackID, "127.0.0.1", freeUDPPort(t))
		if err == nil {
			return nil
		}
		lastErr = err
		if !strings.Contains(err.Error(), "Only one usage") && !strings.Contains(err.Error(), "address already in use") {
			return err
		}
	}
	return lastErr
}

func TestTireModelTelemetryBlockedByQuickListener(t *testing.T) {
	store, err := storage.Open(filepath.Join(t.TempDir(), "tirelab_block.db"))
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	defer store.Close()
	app := &App{telemetry: telemetry.NewService(), store: store}
	if err := app.StartQuickTelemetry("127.0.0.1", freeUDPPort(t)); err != nil {
		t.Fatalf("start quick: %v", err)
	}
	defer app.StopTelemetry()
	if err := app.StartTireModelTelemetry("127.0.0.1", freeUDPPort(t)); err == nil {
		t.Fatalf("start tire model while quick listener is active: got nil error")
	}
}

func TestTrackCaptureTelemetryBlockedByQuickListener(t *testing.T) {
	store, err := storage.Open(filepath.Join(t.TempDir(), "trackcapture_block.db"))
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	defer store.Close()
	app := &App{telemetry: telemetry.NewService(), store: store}
	if err := app.StartQuickTelemetry("127.0.0.1", freeUDPPort(t)); err != nil {
		t.Fatalf("start quick: %v", err)
	}
	defer app.StopTelemetry()
	if err := app.StartTrackCaptureTelemetry("127.0.0.1", freeUDPPort(t)); err == nil {
		t.Fatalf("start track capture while quick listener is active: got nil error")
	}
}

func TestComparabilityWarnings(t *testing.T) {
	known := storage.TelemetrySession{
		GameMode:         telemetry.GameModeFreeRoam,
		DriverMode:       "player",
		BrakeAssist:      "abs_on",
		SteeringAssist:   "simulation",
		TractionControl:  "off",
		StabilityControl: "off",
		Shifting:         "manual",
		LaunchControl:    "off",
	}
	if warnings := comparabilityWarnings(known, known); len(warnings) != 0 {
		t.Fatalf("warnings for equal known sessions = %#v", warnings)
	}

	unknown := known
	unknown.DriverMode = "unknown"
	if warnings := comparabilityWarnings(known, unknown); contains(warnings, "test_conditions_unknown") || !contains(warnings, "driver_mode_mismatch") {
		t.Fatalf("warnings for unknown/mismatch = %#v", warnings)
	}

	other := known
	other.GameMode = telemetry.GameModeRace
	other.BrakeAssist = "abs_off"
	other.TractionControl = "on"
	want := []string{"game_mode_mismatch", "brake_assist_mismatch", "traction_control_mismatch"}
	warnings := comparabilityWarnings(known, other)
	for _, item := range want {
		if !contains(warnings, item) {
			t.Fatalf("warnings missing %q: %#v", item, warnings)
		}
	}
}

func contains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

func freeUDPPort(t *testing.T) int {
	t.Helper()
	conn, err := net.ListenPacket("udp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen udp: %v", err)
	}
	defer conn.Close()
	return conn.LocalAddr().(*net.UDPAddr).Port
}
