package telemetry

import (
	"net"
	"path/filepath"
	"testing"
	"time"
)

func TestServiceStartClearsSessionState(t *testing.T) {
	service := NewService()
	service.parseErrors = 3
	service.current = &NormalizedTelemetry{SpeedKmh: 120}
	service.buffer.Add(time.Now(), NormalizedTelemetry{SpeedKmh: 120})
	observeRepeated(service.engine, baseRuleFrame(NormalizedTelemetry{
		Gear:              1,
		SpeedKmh:          35,
		Throttle01:        0.92,
		RpmRatio:          0.75,
		RearSlipRatioAvg:  1.35,
		FrontSlipRatioAvg: 0.25,
	}), 4)

	oldReceiver := NewReceiver("127.0.0.1", 1, nil)
	oldReceiver.status.validPackets = 4
	oldReceiver.status.invalidPackets = 2
	oldReceiver.status.lastError = "previous error"
	service.receiver = oldReceiver

	startServiceOnFreePort(t, service)
	defer func() {
		if err := service.Stop(); err != nil {
			t.Fatalf("stop failed: %v", err)
		}
	}()

	if service.Current() != nil {
		t.Fatal("expected current telemetry to be cleared")
	}
	if got := len(service.Recent(5)); got != 0 {
		t.Fatalf("recent count = %d, want 0", got)
	}
	if got := len(service.Events()); got != 0 {
		t.Fatalf("event count = %d, want 0", got)
	}

	status := service.Status()
	if status.ValidPackets != 0 || status.InvalidPackets != 0 || status.ParseErrors != 0 {
		t.Fatalf("counters = valid %d invalid %d parse %d, want all zero", status.ValidPackets, status.InvalidPackets, status.ParseErrors)
	}
	if status.LastPacketAt != "" {
		t.Fatalf("last packet = %q, want empty", status.LastPacketAt)
	}
	if status.LastError != "" {
		t.Fatalf("last error = %q, want empty", status.LastError)
	}
}

func TestServiceRecordsSamplesAndReplay(t *testing.T) {
	service := NewService()
	recordingPath := filepath.Join(t.TempDir(), "session.fh6udp")
	port := freeUDPPort(t)
	if err := service.StartWithOptions(StartOptions{Address: "127.0.0.1", Port: port, RecordingPath: recordingPath, RecordingLimitBytes: 1024 * 1024}); err != nil {
		t.Fatalf("start service: %v", err)
	}
	packet := RawPacket{Timestamp: time.Now(), Data: serviceTestPacket(10)}
	service.handlePacket(packet)
	if err := service.Stop(); err != nil {
		t.Fatalf("stop service: %v", err)
	}
	recording := service.Recording()
	if recording.Packets != 1 || recording.Truncated {
		t.Fatalf("recording = %#v", recording)
	}
	if samples := service.Samples(); len(samples) != 1 || samples[0].SpeedKmh != 36 || samples[0].GameMode != GameModeFreeRoam {
		t.Fatalf("samples = %#v", samples)
	}

	replay := NewService()
	if err := replay.Replay(recordingPath, 4); err != nil {
		t.Fatalf("replay: %v", err)
	}
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) && replay.Status().Mode == "replay" {
		time.Sleep(10 * time.Millisecond)
	}
	current := replay.Current()
	if current == nil || current.SpeedKmh != 36 || current.GameMode != GameModeFreeRoam {
		t.Fatalf("replay current = %#v", current)
	}
	if replay.Status().Mode != "idle" {
		t.Fatalf("replay mode = %q, want idle", replay.Status().Mode)
	}
}

func TestServiceReplayPauseResumeAndSeek(t *testing.T) {
	recordingPath := filepath.Join(t.TempDir(), "seek.fh6udp")
	recorder, err := NewRecorder(recordingPath, 1024*1024)
	if err != nil {
		t.Fatalf("new recorder: %v", err)
	}
	start := time.Now().UTC()
	for i, speed := range []float32{10, 20, 30} {
		if err := recorder.Write(RawPacket{Timestamp: start.Add(time.Duration(i) * time.Second), Data: serviceTestPacket(speed)}); err != nil {
			t.Fatalf("write packet: %v", err)
		}
	}
	if err := recorder.Close(); err != nil {
		t.Fatalf("close recorder: %v", err)
	}

	service := NewService()
	if err := service.ReplayWithOptions(ReplayOptions{SessionID: 42, Path: recordingPath, Speed: 1}); err != nil {
		t.Fatalf("replay: %v", err)
	}
	waitForSpeed(t, service, 36)
	if err := service.PauseReplay(); err != nil {
		t.Fatalf("pause: %v", err)
	}
	if status := service.ReplayStatus(); !status.Running || !status.Paused || status.SessionID != 42 || status.PacketCount != 3 {
		t.Fatalf("paused status = %#v", status)
	}
	if err := service.SeekReplay(1000); err != nil {
		t.Fatalf("seek: %v", err)
	}
	waitForSpeed(t, service, 72)
	if status := service.ReplayStatus(); status.PositionMS < 900 {
		t.Fatalf("seek status = %#v", status)
	}
	if err := service.ResumeReplay(); err != nil {
		t.Fatalf("resume: %v", err)
	}
	if err := service.StopReplay(); err != nil {
		t.Fatalf("stop replay: %v", err)
	}
}

func waitForSpeed(t *testing.T, service *Service, speed float64) {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		current := service.Current()
		if current != nil && current.SpeedKmh == speed {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("did not observe speed %.1f, current=%#v status=%#v", speed, service.Current(), service.ReplayStatus())
}

func startServiceOnFreePort(t *testing.T, service *Service) {
	t.Helper()

	var lastErr error
	for i := 0; i < 10; i++ {
		port := freeUDPPort(t)
		if err := service.Start("127.0.0.1", port); err == nil {
			return
		} else {
			lastErr = err
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("start failed after retries: %v", lastErr)
}

func freeUDPPort(t *testing.T) int {
	t.Helper()

	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("resolve free UDP port: %v", err)
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		t.Fatalf("listen free UDP port: %v", err)
	}
	port := conn.LocalAddr().(*net.UDPAddr).Port
	if err := conn.Close(); err != nil {
		t.Fatalf("close free UDP port probe: %v", err)
	}
	time.Sleep(10 * time.Millisecond)

	return port
}

func serviceTestPacket(speedMS float32) []byte {
	spec := DefaultPacketSpec()
	data := make([]byte, spec.Length)
	o := spec.Offsets
	putI32(data, o.IsRaceOn, 1)
	putU32(data, o.TimestampMS, 1000)
	putF32(data, o.EngineIdleRpm, 900)
	putF32(data, o.EngineMaxRpm, 8000)
	putF32(data, o.CurrentEngineRpm, 3000)
	putF32(data, o.Speed, speedMS)
	putI32(data, o.CarOrdinal, 1001)
	putI32(data, o.CarClass, 4)
	putI32(data, o.CarPI, 850)
	data[o.Accel] = 128
	data[o.Gear] = 2
	return data
}
