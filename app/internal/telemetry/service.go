package telemetry

import (
	"fmt"
	"net"
	"sync"
	"time"
)

type StartOptions struct {
	Address             string
	Port                int
	RecordingPath       string
	RecordingLimitBytes int64
}

type ReplayOptions struct {
	SessionID int64
	Path      string
	Speed     float64
	Config    RuleConfig
}

type Service struct {
	mu          sync.RWMutex
	parser      Parser
	buffer      *RingBuffer
	receiver    *Receiver
	recorder    *Recorder
	engine      *RuleEngine
	parseErrors uint64
	current     *NormalizedTelemetry
	summary     TelemetrySummary
	speedSum    float64
	samples     []NormalizedTelemetry
	lastSample  time.Time
	gameMode    GameModeTracker

	mode              string
	lastRecording     RecordingSnapshot
	replayStop        chan struct{}
	replaySeek        chan int64
	replaySessionID   int64
	replaySpeed       float64
	replayPaused      bool
	replayPositionMS  int64
	replayDurationMS  int64
	replayPacketIndex int
	replayPacketCount int
	replayPackets     uint64
	replayLastPacket  time.Time
	replayLastError   string
}

func NewService() *Service {
	return &Service{
		parser:   DefaultParser(),
		buffer:   NewRingBuffer(30 * time.Second),
		engine:   NewRuleEngine(),
		gameMode: NewGameModeTracker(),
	}
}

func (s *Service) Start(address string, port int) error {
	return s.StartWithOptions(StartOptions{Address: address, Port: port})
}

func (s *Service) StartWithOptions(options StartOptions) error {
	if options.Address == "" {
		options.Address = DefaultAddress
	}
	if options.Port == 0 {
		options.Port = DefaultPort
	}
	if options.Port < 1 || options.Port > 65535 {
		return fmt.Errorf("invalid UDP port: %d", options.Port)
	}
	if net.ParseIP(options.Address) == nil {
		return fmt.Errorf("invalid listen address: %s", options.Address)
	}

	s.mu.Lock()
	if s.receiver != nil && s.receiver.Snapshot().Running {
		s.mu.Unlock()
		return fmt.Errorf("telemetry listener is already running")
	}
	if s.mode == "replay" {
		s.mu.Unlock()
		return fmt.Errorf("telemetry replay is already running")
	}
	var recorder *Recorder
	if options.RecordingPath != "" {
		created, err := NewRecorder(options.RecordingPath, options.RecordingLimitBytes)
		if err != nil {
			s.mu.Unlock()
			return err
		}
		recorder = created
	}
	s.resetRuntimeLocked()
	s.recorder = recorder
	if recorder != nil {
		s.lastRecording = recorder.Snapshot()
	} else {
		s.lastRecording = RecordingSnapshot{LimitBytes: DefaultRecordingLimit}
	}
	receiver := NewReceiver(options.Address, options.Port, s.handlePacket)
	s.receiver = receiver
	s.mode = "udp"
	s.mu.Unlock()

	if err := receiver.Start(); err != nil {
		s.mu.Lock()
		if s.recorder != nil {
			_ = s.recorder.Close()
		}
		s.recorder = nil
		s.receiver = nil
		s.mode = "idle"
		s.mu.Unlock()
		return err
	}
	return nil
}

func (s *Service) Stop() error {
	s.mu.RLock()
	receiver := s.receiver
	s.mu.RUnlock()
	var err error
	if receiver != nil {
		err = receiver.Stop()
	}
	s.mu.Lock()
	if s.recorder != nil {
		closeErr := s.recorder.Close()
		s.lastRecording = s.recorder.Snapshot()
		s.recorder = nil
		if err == nil {
			err = closeErr
		}
	}
	s.mode = "idle"
	s.mu.Unlock()
	return err
}

func (s *Service) Status() TelemetryStatus {
	s.mu.RLock()
	receiver := s.receiver
	parseErrors := s.parseErrors
	hasCurrent := s.current != nil
	mode := s.mode
	recording := s.lastRecording
	if s.recorder != nil {
		recording = s.recorder.Snapshot()
	}
	replayPackets := s.replayPackets
	replayLastPacket := s.replayLastPacket
	replayLastError := s.replayLastError
	s.mu.RUnlock()

	status := TelemetryStatus{
		Address:      DefaultAddress,
		Port:         DefaultPort,
		PacketLength: PacketLength,
		Mode:         "idle",
	}
	if receiver != nil {
		status = receiver.Snapshot()
		status.Mode = "udp"
	}
	if mode == "replay" {
		status.Running = true
		status.Mode = "replay"
		status.ValidPackets = replayPackets
		status.RawPackets = replayPackets
		status.InvalidPackets = 0
		status.Address = "replay"
		status.Port = 0
		status.LastError = replayLastError
		if !replayLastPacket.IsZero() {
			status.LastPacketAt = replayLastPacket.UTC().Format(time.RFC3339Nano)
		}
	}
	status.ParseErrors = parseErrors
	status.HasCurrentFrame = hasCurrent
	status.RecordingActive = recording.Active
	status.RecordingBytes = recording.Bytes
	status.RecordingLimitBytes = recording.LimitBytes
	status.RecordingPackets = recording.Packets
	status.RecordingTruncated = recording.Truncated
	if status.Mode == "" {
		status.Mode = "idle"
	}
	return status
}

func (s *Service) NetworkInterfaces() []NetworkInterface {
	return ListNetworkInterfaces()
}

func (s *Service) SetRuleConfig(config RuleConfig) {
	s.engine.SetConfig(config)
}

func (s *Service) Current() *NormalizedTelemetry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.current == nil {
		return nil
	}
	frame := *s.current
	return &frame
}

func (s *Service) Recent(seconds int) []NormalizedTelemetry {
	if seconds <= 0 {
		seconds = 5
	}
	return s.buffer.Since(seconds)
}

func (s *Service) Events() []DetectedEvent {
	return s.engine.Events()
}

func (s *Service) Summary() TelemetrySummary {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.summary
}

func (s *Service) Samples() []NormalizedTelemetry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]NormalizedTelemetry, len(s.samples))
	copy(out, s.samples)
	return out
}

func (s *Service) Recording() RecordingSnapshot {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.recorder != nil {
		return s.recorder.Snapshot()
	}
	return s.lastRecording
}

func (s *Service) Replay(path string, speed float64) error {
	return s.ReplayWithOptions(ReplayOptions{Path: path, Speed: speed})
}

func (s *Service) ReplayWithOptions(options ReplayOptions) error {
	speed := options.Speed
	if speed <= 0 {
		speed = 1
	}
	packets, err := LoadRecording(options.Path)
	if err != nil {
		return err
	}
	if len(packets) == 0 {
		return fmt.Errorf("recording has no packets")
	}
	duration := packets[len(packets)-1].Timestamp.Sub(packets[0].Timestamp).Milliseconds()
	if duration < 0 {
		duration = 0
	}
	s.mu.Lock()
	if s.receiver != nil && s.receiver.Snapshot().Running {
		s.mu.Unlock()
		return fmt.Errorf("telemetry listener is already running")
	}
	if s.mode == "replay" {
		s.mu.Unlock()
		return fmt.Errorf("telemetry replay is already running")
	}
	s.engine.SetConfig(options.Config)
	s.resetRuntimeLocked()
	stop := make(chan struct{})
	seek := make(chan int64, 1)
	s.replayStop = stop
	s.replaySeek = seek
	s.replaySessionID = options.SessionID
	s.replaySpeed = speed
	s.replayPaused = false
	s.replayPositionMS = 0
	s.replayDurationMS = duration
	s.replayPacketIndex = 0
	s.replayPacketCount = len(packets)
	s.replayPackets = 0
	s.replayLastPacket = time.Time{}
	s.replayLastError = ""
	s.mode = "replay"
	s.mu.Unlock()

	go s.replayLoop(packets, speed, stop, seek)
	return nil
}

func (s *Service) StopReplay() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.mode != "replay" || s.replayStop == nil {
		return nil
	}
	close(s.replayStop)
	s.replayStop = nil
	s.replaySeek = nil
	s.mode = "idle"
	return nil
}

func (s *Service) PauseReplay() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.mode != "replay" {
		return fmt.Errorf("telemetry replay is not running")
	}
	s.replayPaused = true
	return nil
}

func (s *Service) ResumeReplay() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.mode != "replay" {
		return fmt.Errorf("telemetry replay is not running")
	}
	s.replayPaused = false
	return nil
}

func (s *Service) SeekReplay(positionMS int64) error {
	s.mu.RLock()
	seek := s.replaySeek
	running := s.mode == "replay"
	duration := s.replayDurationMS
	s.mu.RUnlock()
	if !running || seek == nil {
		return fmt.Errorf("telemetry replay is not running")
	}
	if positionMS < 0 {
		positionMS = 0
	}
	if duration > 0 && positionMS > duration {
		positionMS = duration
	}
	select {
	case seek <- positionMS:
	default:
		select {
		case <-seek:
		default:
		}
		seek <- positionMS
	}
	return nil
}

func (s *Service) ReplayStatus() TelemetryReplayStatus {
	s.mu.RLock()
	defer s.mu.RUnlock()
	status := TelemetryReplayStatus{
		Running:     s.mode == "replay",
		Paused:      s.replayPaused,
		SessionID:   s.replaySessionID,
		Speed:       s.replaySpeed,
		PositionMS:  s.replayPositionMS,
		DurationMS:  s.replayDurationMS,
		PacketIndex: s.replayPacketIndex,
		PacketCount: s.replayPacketCount,
		LastError:   s.replayLastError,
	}
	if status.DurationMS > 0 {
		status.Progress01 = clamp01(float64(status.PositionMS) / float64(status.DurationMS))
	}
	return status
}

func (s *Service) handlePacket(packet RawPacket) {
	s.processPacket(packet, true)
}

func (s *Service) processPacket(packet RawPacket, record bool) {
	if record {
		s.mu.RLock()
		recorder := s.recorder
		s.mu.RUnlock()
		if recorder != nil {
			_ = recorder.Write(packet)
		}
	}

	frame, err := s.parser.Parse(packet.Data, packet.Timestamp)
	if err != nil {
		s.mu.Lock()
		s.parseErrors++
		s.mu.Unlock()
		return
	}

	normalized := NormalizeFrame(frame)
	s.mu.Lock()
	normalized.GameMode = s.gameMode.Observe(normalized, packet.Timestamp)
	s.mu.Unlock()
	s.buffer.Add(packet.Timestamp, normalized)

	s.mu.Lock()
	s.current = &normalized
	s.summary.SampleCount++
	s.speedSum += normalized.SpeedKmh
	s.summary.AvgSpeedKmh = s.speedSum / float64(s.summary.SampleCount)
	if normalized.SpeedKmh > s.summary.MaxSpeedKmh {
		s.summary.MaxSpeedKmh = normalized.SpeedKmh
	}
	if s.lastSample.IsZero() || packet.Timestamp.Sub(s.lastSample) >= defaultAggregateBucket {
		s.samples = append(s.samples, normalized)
		s.lastSample = packet.Timestamp
	}
	if s.mode == "replay" {
		s.replayPackets++
		s.replayLastPacket = packet.Timestamp
	}
	s.mu.Unlock()

	s.engine.Observe(normalized)
}

func (s *Service) replayLoop(packets []RecordingPacket, speed float64, stop <-chan struct{}, seek <-chan int64) {
	base := packets[0].Timestamp
	index := 0
	for index < len(packets) {
		if nextIndex, ok := s.waitReplayTurn(packets, index, speed, stop, seek, base); !ok {
			return
		} else if nextIndex >= 0 {
			index = nextIndex
		}

		packet := packets[index]
		position := packet.Timestamp.Sub(base).Milliseconds()
		if position < 0 {
			position = 0
		}
		s.mu.Lock()
		s.replayPacketIndex = index
		s.replayPositionMS = position
		s.mu.Unlock()
		s.processPacket(RawPacket{Timestamp: time.Now(), Data: packet.Data, Addr: "replay"}, false)
		index++
	}

	s.mu.Lock()
	if s.mode == "replay" {
		s.mode = "idle"
		s.replayStop = nil
		s.replaySeek = nil
		s.replayPaused = false
	}
	s.mu.Unlock()
}

func (s *Service) waitReplayTurn(packets []RecordingPacket, index int, speed float64, stop <-chan struct{}, seek <-chan int64, base time.Time) (int, bool) {
	if index > 0 {
		delay := packets[index].Timestamp.Sub(packets[index-1].Timestamp)
		if delay > 0 {
			timer := time.NewTimer(time.Duration(float64(delay) / speed))
			defer timer.Stop()
			for {
				select {
				case <-stop:
					return -1, false
				case position := <-seek:
					return s.applyReplaySeek(packets, position, base), true
				case <-timer.C:
					goto afterDelay
				}
			}
		}
	}

afterDelay:
	for {
		s.mu.RLock()
		paused := s.replayPaused
		s.mu.RUnlock()
		if !paused {
			return -1, true
		}
		select {
		case <-stop:
			return -1, false
		case position := <-seek:
			return s.applyReplaySeek(packets, position, base), true
		case <-time.After(50 * time.Millisecond):
		}
	}
}

func (s *Service) applyReplaySeek(packets []RecordingPacket, positionMS int64, base time.Time) int {
	target := base.Add(time.Duration(positionMS) * time.Millisecond)
	index := 0
	for i, packet := range packets {
		if !packet.Timestamp.Before(target) {
			index = i
			break
		}
		index = i
	}
	s.mu.Lock()
	s.resetRuntimeLocked()
	s.replayPackets = uint64(index)
	s.replayPacketIndex = index
	s.replayPositionMS = packets[index].Timestamp.Sub(base).Milliseconds()
	if s.replayPositionMS < 0 {
		s.replayPositionMS = 0
	}
	s.mu.Unlock()
	return index
}

func (s *Service) resetRuntimeLocked() {
	s.buffer.Reset()
	s.engine.Reset()
	s.parseErrors = 0
	s.current = nil
	s.summary = TelemetrySummary{}
	s.speedSum = 0
	s.samples = nil
	s.lastSample = time.Time{}
	s.gameMode.Reset()
}
