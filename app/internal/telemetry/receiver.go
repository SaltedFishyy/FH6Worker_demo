package telemetry

import (
	"fmt"
	"net"
	"sync"
	"time"
)

type Receiver struct {
	address string
	port    int
	handler func(RawPacket)

	mu     sync.RWMutex
	conn   *net.UDPConn
	done   chan struct{}
	status receiverStatus
}

type receiverStatus struct {
	running        bool
	rawPackets     uint64
	validPackets   uint64
	invalidPackets uint64
	lastDatagramAt time.Time
	lastBytes      int
	lastRemote     string
	lastPacketAt   time.Time
	lastError      string
}

func NewReceiver(address string, port int, handler func(RawPacket)) *Receiver {
	if address == "" {
		address = DefaultAddress
	}
	if port == 0 {
		port = DefaultPort
	}
	return &Receiver{
		address: address,
		port:    port,
		handler: handler,
		done:    make(chan struct{}),
	}
}

func (r *Receiver) Start() error {
	r.mu.Lock()
	if r.status.running {
		r.mu.Unlock()
		return fmt.Errorf("telemetry listener is already running")
	}
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", r.address, r.port))
	if err != nil {
		r.status.lastError = err.Error()
		r.mu.Unlock()
		return err
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		r.status.lastError = err.Error()
		r.mu.Unlock()
		return err
	}
	r.conn = conn
	r.done = make(chan struct{})
	r.status.running = true
	r.status.lastError = ""
	r.mu.Unlock()

	go r.readLoop(conn)
	return nil
}

func (r *Receiver) Stop() error {
	r.mu.Lock()
	if !r.status.running {
		r.mu.Unlock()
		return nil
	}
	close(r.done)
	err := r.conn.Close()
	r.status.running = false
	r.conn = nil
	r.mu.Unlock()
	return err
}

func (r *Receiver) Snapshot() TelemetryStatus {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var lastPacketAt string
	if !r.status.lastPacketAt.IsZero() {
		lastPacketAt = r.status.lastPacketAt.UTC().Format(time.RFC3339Nano)
	}
	var lastDatagramAt string
	if !r.status.lastDatagramAt.IsZero() {
		lastDatagramAt = r.status.lastDatagramAt.UTC().Format(time.RFC3339Nano)
	}
	return TelemetryStatus{
		Running:            r.status.running,
		Address:            r.address,
		Port:               r.port,
		PacketLength:       PacketLength,
		RawPackets:         r.status.rawPackets,
		ValidPackets:       r.status.validPackets,
		InvalidPackets:     r.status.invalidPackets,
		LastDatagramAt:     lastDatagramAt,
		LastDatagramBytes:  r.status.lastBytes,
		LastDatagramRemote: r.status.lastRemote,
		LastPacketAt:       lastPacketAt,
		LastError:          r.status.lastError,
	}
}

func (r *Receiver) readLoop(conn *net.UDPConn) {
	buffer := make([]byte, 1500)
	for {
		select {
		case <-r.done:
			return
		default:
		}

		_ = conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
		n, remoteAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				continue
			}
			select {
			case <-r.done:
				return
			default:
			}
			r.setError(err)
			continue
		}

		now := time.Now()
		r.recordDatagram(now, n, remoteAddr)

		if n != PacketLength {
			r.incrementInvalid(n, remoteAddr)
			continue
		}

		data := make([]byte, n)
		copy(data, buffer[:n])
		r.incrementValid(now)

		if r.handler != nil {
			r.handler(RawPacket{
				Timestamp: now,
				Data:      data,
				Addr:      remoteAddr.String(),
			})
		}
	}
}

func (r *Receiver) recordDatagram(at time.Time, bytes int, remote *net.UDPAddr) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.status.rawPackets++
	r.status.lastDatagramAt = at
	r.status.lastBytes = bytes
	if remote != nil {
		r.status.lastRemote = remote.String()
	}
}

func (r *Receiver) incrementValid(at time.Time) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.status.validPackets++
	r.status.lastPacketAt = at
	r.status.lastError = ""
}

func (r *Receiver) incrementInvalid(bytes int, remote *net.UDPAddr) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.status.invalidPackets++
	remoteText := ""
	if remote != nil {
		remoteText = " from " + remote.String()
	}
	r.status.lastError = fmt.Sprintf("unexpected UDP packet size %d%s, expected %d", bytes, remoteText, PacketLength)
}

func (r *Receiver) setError(err error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.status.lastError = err.Error()
}
