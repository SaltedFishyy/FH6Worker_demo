package telemetry

import (
	"testing"
	"time"
)

func TestRingBufferPrunesOldFrames(t *testing.T) {
	buffer := NewRingBuffer(30 * time.Second)
	base := time.Now()

	buffer.Add(base.Add(-31*time.Second), NormalizedTelemetry{SpeedKmh: 10})
	buffer.Add(base, NormalizedTelemetry{SpeedKmh: 20})

	if got := buffer.Len(); got != 1 {
		t.Fatalf("buffer length = %d, want 1", got)
	}
	latest := buffer.Latest()
	if latest == nil || latest.SpeedKmh != 20 {
		t.Fatalf("latest frame = %#v, want speed 20", latest)
	}
}

func TestRingBufferAggregatesAtTenHz(t *testing.T) {
	buffer := NewRingBuffer(30 * time.Second)
	base := time.Now().Add(-time.Second)
	buffer.Add(base, NormalizedTelemetry{ReceivedAt: base.Format(time.RFC3339Nano), SpeedKmh: 100, Rpm: 4000, Gear: 2, WheelFL: NormalizedWheelTelemetry{CombinedSlip: 0.4}})
	buffer.Add(base.Add(50*time.Millisecond), NormalizedTelemetry{ReceivedAt: base.Add(50 * time.Millisecond).Format(time.RFC3339Nano), SpeedKmh: 120, Rpm: 5000, Gear: 2, WheelFL: NormalizedWheelTelemetry{CombinedSlip: 0.8}})
	buffer.Add(base.Add(150*time.Millisecond), NormalizedTelemetry{ReceivedAt: base.Add(150 * time.Millisecond).Format(time.RFC3339Nano), SpeedKmh: 80, Rpm: 3000, Gear: 3, WheelFL: NormalizedWheelTelemetry{CombinedSlip: 0.2}})

	got := buffer.Aggregated(2, 100*time.Millisecond)
	if len(got) != 2 {
		t.Fatalf("aggregated frame count = %d, want 2", len(got))
	}
	if got[0].SpeedKmh != 110 {
		t.Fatalf("first bucket speed = %f, want 110", got[0].SpeedKmh)
	}
	if got[0].Rpm != 4500 {
		t.Fatalf("first bucket rpm = %f, want 4500", got[0].Rpm)
	}
	if got[1].Gear != 3 {
		t.Fatalf("second bucket gear = %d, want 3", got[1].Gear)
	}
}
