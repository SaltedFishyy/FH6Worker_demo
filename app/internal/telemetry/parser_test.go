package telemetry

import (
	"encoding/binary"
	"math"
	"testing"
	"time"
)

func TestParseRejectsInvalidLength(t *testing.T) {
	parser := DefaultParser()
	_, err := parser.Parse(make([]byte, PacketLength-1), time.Now())
	if err == nil {
		t.Fatal("expected invalid length error")
	}
}

func TestParseAndNormalizeCoreFields(t *testing.T) {
	spec := DefaultPacketSpec()
	data := make([]byte, spec.Length)
	o := spec.Offsets

	putI32(data, o.IsRaceOn, 1)
	putU32(data, o.TimestampMS, 12345)
	putF32(data, o.EngineIdleRpm, 1000)
	putF32(data, o.EngineMaxRpm, 8000)
	putF32(data, o.CurrentEngineRpm, 4500)
	putF32(data, o.Speed, 50)
	putI32(data, o.CarOrdinal, 123456)
	putI32(data, o.CarClass, 4)
	putI32(data, o.CarPI, 860)
	putI32(data, o.DrivetrainType, 2)
	putI32(data, o.NumCylinders, 6)
	putU32(data, o.CarCategory, 12)
	data[o.Accel] = 128
	data[o.Brake] = 64
	data[o.Gear] = 3
	data[o.Steer] = 0xC0
	putF32(data, o.SuspensionTravelFL, 0.5)
	putF32(data, o.SuspensionTravelFR, 0.6)
	putF32(data, o.SuspensionTravelRL, 0.7)
	putF32(data, o.SuspensionTravelRR, 0.8)
	putF32(data, o.TireSlipRatioFL, 0.2)
	putF32(data, o.TireSlipRatioFR, 0.4)
	putF32(data, o.TireSlipRatioRL, 1.2)
	putF32(data, o.TireSlipRatioRR, 1.4)
	putI32(data, o.WheelOnRumbleStripFL, 1)
	putF32(data, o.WheelInPuddleDepthFL, 0.25)
	putF32(data, o.SurfaceRumbleFL, 0.33)
	putF32(data, o.TireSlipAngleFL, 0.1)
	putF32(data, o.TireSlipAngleFR, 0.3)
	putF32(data, o.TireSlipAngleRL, 0.5)
	putF32(data, o.TireSlipAngleRR, 0.7)
	putF32(data, o.TireCombinedSlipFL, 0.8)
	putF32(data, o.TireCombinedSlipFR, 0.9)
	putF32(data, o.TireCombinedSlipRL, 1.4)
	putF32(data, o.TireCombinedSlipRR, 1.6)
	putF32(data, o.SuspensionTravelMetersFL, 0.045)
	putF32(data, o.Power, 345000)
	putF32(data, o.Torque, 520)
	putF32(data, o.Boost, 1.2)
	putF32(data, o.Fuel, 0.75)
	putF32(data, o.DistanceTraveled, 1234)
	putF32(data, o.CurrentRaceTime, 88.5)
	putF32(data, o.PositionX, 10)
	putF32(data, o.PositionY, 20)
	putF32(data, o.PositionZ, 30)
	putF32(data, o.TireTempFL, 82)
	putF32(data, o.TireTempFR, 84)
	putF32(data, o.TireTempRL, 91)
	putF32(data, o.TireTempRR, 93)
	putU16(data, o.LapNumber, 2)
	data[o.RacePosition] = 4
	data[o.Clutch] = 255
	data[o.HandBrake] = 32

	parser := DefaultParser()
	frame, err := parser.Parse(data, time.Unix(10, 0))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	normalized := NormalizeFrame(frame)
	assertClose(t, normalized.SpeedKmh, 180, "speed kmh")
	assertClose(t, normalized.SpeedFieldKmh, 180, "speed field kmh")
	assertClose(t, normalized.RpmRatio, 0.5, "rpm ratio")
	assertClose(t, normalized.Throttle01, 128.0/255.0, "throttle")
	assertClose(t, normalized.Brake01, 64.0/255.0, "brake")
	assertClose(t, normalized.Steer01, -0.5, "steer")
	assertClose(t, normalized.RearSlipRatioAvg, 1.3, "rear slip ratio")
	assertClose(t, normalized.FrontCombinedSlipAvg, 0.85, "front combined slip")
	assertClose(t, normalized.TireTempRearAvg, 92, "rear tire temp")
	assertClose(t, normalized.WheelFL.RumbleStrip, 1, "rumble strip")
	assertClose(t, normalized.WheelFL.PuddleDepth, 0.25, "puddle depth")
	assertClose(t, normalized.WheelFL.SurfaceRumble, 0.33, "surface rumble")
	assertClose(t, normalized.WheelFL.SuspensionTravelMeters, 0.045, "suspension meters")
	assertClose(t, normalized.Power, 345000, "power")
	assertClose(t, normalized.Torque, 520, "torque")
	assertClose(t, normalized.Boost, 1.2, "boost")
	assertClose(t, normalized.Fuel, 0.75, "fuel")
	assertClose(t, normalized.DistanceTraveled, 1234, "distance")
	assertClose(t, normalized.CurrentRaceTime, 88.5, "race time")
	assertClose(t, normalized.PositionX, 10, "position x")
	assertClose(t, normalized.Clutch01, 1, "clutch")
	assertClose(t, normalized.HandBrake01, 32.0/255.0, "handbrake")
	if normalized.Gear != 3 {
		t.Fatalf("gear = %d, want 3", normalized.Gear)
	}
	if normalized.LapNumber != 2 || normalized.RacePosition != 4 {
		t.Fatalf("race metadata = lap %d position %d", normalized.LapNumber, normalized.RacePosition)
	}
	if normalized.SpeedSource != "packet" {
		t.Fatalf("speed source = %q, want packet", normalized.SpeedSource)
	}
	if normalized.CarOrdinal != 123456 || normalized.CarClass != "S1" || normalized.CarPI != 860 || normalized.Drivetrain != "AWD" || normalized.NumCylinders != 6 || normalized.CarCategory != 12 || normalized.CarCategoryName != "复古超级跑车" {
		t.Fatalf("vehicle metadata = %#v", normalized)
	}
	if !normalized.IsRaceOn {
		t.Fatal("expected race-on flag")
	}
}

func TestNormalizeUsesVelocitySpeedWhenPacketSpeedDisagrees(t *testing.T) {
	frame := TelemetryFrame{
		ReceivedAt:  time.Unix(20, 0),
		TimestampMS: 2000,
		Speed:       3,
		VelocityX:   20,
	}

	normalized := NormalizeFrame(frame)
	assertClose(t, normalized.SpeedKmh, 72, "velocity speed kmh")
	assertClose(t, normalized.SpeedFieldKmh, 10.8, "packet speed kmh")
	assertClose(t, normalized.VelocitySpeedKmh, 72, "velocity debug kmh")
	if normalized.SpeedSource != "velocity" {
		t.Fatalf("speed source = %q, want velocity", normalized.SpeedSource)
	}
}

func putF32(data []byte, off int, value float32) {
	binary.LittleEndian.PutUint32(data[off:off+4], math.Float32bits(value))
}

func putI32(data []byte, off int, value int32) {
	binary.LittleEndian.PutUint32(data[off:off+4], uint32(value))
}

func putU32(data []byte, off int, value uint32) {
	binary.LittleEndian.PutUint32(data[off:off+4], value)
}

func putU16(data []byte, off int, value uint16) {
	binary.LittleEndian.PutUint16(data[off:off+2], value)
}

func assertClose(t *testing.T, got, want float64, name string) {
	t.Helper()
	if math.Abs(got-want) > 0.0001 {
		t.Fatalf("%s = %f, want %f", name, got, want)
	}
}
