package telemetry

import (
	"path/filepath"
	"testing"
	"time"
)

func TestRecorderRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "session.fh6udp")
	recorder, err := NewRecorder(path, 1024*1024)
	if err != nil {
		t.Fatalf("new recorder: %v", err)
	}
	packet := RawPacket{Timestamp: time.Unix(100, 123), Data: make([]byte, PacketLength)}
	packet.Data[0] = 7
	if err := recorder.Write(packet); err != nil {
		t.Fatalf("write packet: %v", err)
	}
	if err := recorder.Close(); err != nil {
		t.Fatalf("close recorder: %v", err)
	}

	var packets []RecordingPacket
	if err := ReadRecording(path, func(packet RecordingPacket) error {
		packets = append(packets, packet)
		return nil
	}); err != nil {
		t.Fatalf("read recording: %v", err)
	}
	if len(packets) != 1 {
		t.Fatalf("packet count = %d, want 1", len(packets))
	}
	if !packets[0].Timestamp.Equal(packet.Timestamp.UTC()) || len(packets[0].Data) != PacketLength || packets[0].Data[0] != 7 {
		t.Fatalf("packet = %#v", packets[0])
	}
}

func TestRecorderLimitTruncatesWithoutWritingOverflow(t *testing.T) {
	path := filepath.Join(t.TempDir(), "limited.fh6udp")
	recorder, err := NewRecorder(path, int64(len(RecordingMagic)+recordHeaderSize+PacketLength))
	if err != nil {
		t.Fatalf("new recorder: %v", err)
	}
	packet := RawPacket{Timestamp: time.Now(), Data: make([]byte, PacketLength)}
	if err := recorder.Write(packet); err != nil {
		t.Fatalf("write first packet: %v", err)
	}
	if err := recorder.Write(packet); err != nil {
		t.Fatalf("write second packet: %v", err)
	}
	snapshot := recorder.Snapshot()
	if snapshot.Packets != 1 || !snapshot.Truncated || snapshot.Active {
		t.Fatalf("snapshot = %#v", snapshot)
	}
	if err := recorder.Close(); err != nil {
		t.Fatalf("close recorder: %v", err)
	}
}
