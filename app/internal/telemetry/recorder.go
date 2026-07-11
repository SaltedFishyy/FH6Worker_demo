package telemetry

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	RecordingMagic        = "FH6UDP1"
	DefaultRecordingLimit = int64(128 * 1024 * 1024)
	recordHeaderSize      = 12
)

type RecordingSnapshot struct {
	Path       string `json:"path"`
	Active     bool   `json:"active"`
	Packets    int64  `json:"packets"`
	Bytes      int64  `json:"bytes"`
	LimitBytes int64  `json:"limitBytes"`
	Truncated  bool   `json:"truncated"`
}

type RecordingPacket struct {
	Timestamp time.Time
	Data      []byte
}

type Recorder struct {
	mu     sync.Mutex
	file   *os.File
	path   string
	limit  int64
	bytes  int64
	count  int64
	closed bool
	cut    bool
}

func NewRecorder(path string, limitBytes int64) (*Recorder, error) {
	if limitBytes <= 0 {
		limitBytes = DefaultRecordingLimit
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, err
	}
	file, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	n, err := file.Write([]byte(RecordingMagic))
	if err != nil {
		_ = file.Close()
		return nil, err
	}
	return &Recorder{
		file:  file,
		path:  path,
		limit: limitBytes,
		bytes: int64(n),
	}, nil
}

func (r *Recorder) Write(packet RawPacket) error {
	if r == nil {
		return nil
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.closed || r.cut {
		return nil
	}
	if len(packet.Data) != PacketLength {
		return nil
	}
	nextSize := int64(recordHeaderSize + len(packet.Data))
	if r.bytes+nextSize > r.limit {
		r.cut = true
		return nil
	}

	var header [recordHeaderSize]byte
	binary.LittleEndian.PutUint64(header[0:8], uint64(packet.Timestamp.UTC().UnixNano()))
	binary.LittleEndian.PutUint32(header[8:12], uint32(len(packet.Data)))
	if _, err := r.file.Write(header[:]); err != nil {
		return err
	}
	if _, err := r.file.Write(packet.Data); err != nil {
		return err
	}
	r.bytes += nextSize
	r.count++
	return nil
}

func (r *Recorder) Close() error {
	if r == nil {
		return nil
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.closed {
		return nil
	}
	r.closed = true
	if r.file == nil {
		return nil
	}
	err := r.file.Close()
	r.file = nil
	return err
}

func (r *Recorder) Snapshot() RecordingSnapshot {
	if r == nil {
		return RecordingSnapshot{LimitBytes: DefaultRecordingLimit}
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	return RecordingSnapshot{
		Path:       r.path,
		Active:     !r.closed && !r.cut,
		Packets:    r.count,
		Bytes:      r.bytes,
		LimitBytes: r.limit,
		Truncated:  r.cut,
	}
}

func ReadRecording(path string, emit func(RecordingPacket) error) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	magic := make([]byte, len(RecordingMagic))
	if _, err := io.ReadFull(file, magic); err != nil {
		return err
	}
	if string(magic) != RecordingMagic {
		return fmt.Errorf("invalid recording magic: %q", string(magic))
	}

	var header [recordHeaderSize]byte
	for {
		if _, err := io.ReadFull(file, header[:]); err != nil {
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				return nil
			}
			return err
		}
		ts := int64(binary.LittleEndian.Uint64(header[0:8]))
		size := int(binary.LittleEndian.Uint32(header[8:12]))
		if size <= 0 || size > 1500 {
			return fmt.Errorf("invalid recording packet size: %d", size)
		}
		data := make([]byte, size)
		if _, err := io.ReadFull(file, data); err != nil {
			return err
		}
		if err := emit(RecordingPacket{Timestamp: time.Unix(0, ts).UTC(), Data: data}); err != nil {
			return err
		}
	}
}

func LoadRecording(path string) ([]RecordingPacket, error) {
	var packets []RecordingPacket
	err := ReadRecording(path, func(packet RecordingPacket) error {
		packets = append(packets, packet)
		return nil
	})
	return packets, err
}
