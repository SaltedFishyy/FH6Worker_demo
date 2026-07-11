package main

import (
	"encoding/binary"
	"flag"
	"log"
	"math"
	"net"
	"time"

	"fh6worker/internal/telemetry"
)

func main() {
	addr := flag.String("addr", "127.0.0.1:5301", "UDP target address")
	interval := flag.Duration("interval", 100*time.Millisecond, "packet interval")
	flag.Parse()

	conn, err := net.Dial("udp", *addr)
	if err != nil {
		log.Fatalf("dial UDP: %v", err)
	}
	defer conn.Close()

	start := time.Now()
	ticker := time.NewTicker(*interval)
	defer ticker.Stop()

	for now := range ticker.C {
		elapsed := now.Sub(start)
		packet := syntheticPacket(elapsed)
		if _, err := conn.Write(packet); err != nil {
			log.Printf("write UDP packet: %v", err)
		}
	}
}

func syntheticPacket(elapsed time.Duration) []byte {
	spec := telemetry.DefaultPacketSpec()
	o := spec.Offsets
	data := make([]byte, spec.Length)
	seconds := elapsed.Seconds()
	speedMS := 25 + 10*math.Sin(seconds)
	throttle := 150 + int(80*math.Sin(seconds*0.5))
	if throttle < 0 {
		throttle = 0
	}
	if throttle > 255 {
		throttle = 255
	}

	putI32(data, o.IsRaceOn, 1)
	putU32(data, o.TimestampMS, uint32(elapsed.Milliseconds()))
	putF32(data, o.EngineIdleRpm, 900)
	putF32(data, o.EngineMaxRpm, 8200)
	putF32(data, o.CurrentEngineRpm, float32(3500+2500*math.Abs(math.Sin(seconds))))
	putF32(data, o.VelocityZ, float32(speedMS))
	putF32(data, o.Speed, float32(speedMS))
	putI32(data, o.CarOrdinal, 214748)
	putI32(data, o.CarClass, 4)
	putI32(data, o.CarPI, 862)
	putI32(data, o.DrivetrainType, 2)
	putI32(data, o.NumCylinders, 6)
	putU32(data, o.CarCategory, 12)
	putF32(data, o.AngularVelocityY, float32(0.25*math.Sin(seconds)))
	putF32(data, o.SuspensionTravelFL, 0.45)
	putF32(data, o.SuspensionTravelFR, 0.47)
	putF32(data, o.SuspensionTravelRL, 0.53)
	putF32(data, o.SuspensionTravelRR, 0.55)
	putF32(data, o.TireSlipRatioFL, 0.28)
	putF32(data, o.TireSlipRatioFR, 0.30)
	putF32(data, o.TireSlipRatioRL, float32(0.45+0.2*math.Abs(math.Sin(seconds))))
	putF32(data, o.TireSlipRatioRR, float32(0.48+0.22*math.Abs(math.Cos(seconds))))
	putF32(data, o.TireSlipAngleFL, 0.08)
	putF32(data, o.TireSlipAngleFR, 0.09)
	putF32(data, o.TireSlipAngleRL, 0.14)
	putF32(data, o.TireSlipAngleRR, 0.13)
	putF32(data, o.TireCombinedSlipFL, 0.38)
	putF32(data, o.TireCombinedSlipFR, 0.40)
	putF32(data, o.TireCombinedSlipRL, 0.72)
	putF32(data, o.TireCombinedSlipRR, 0.78)
	putF32(data, o.TireTempFL, 82)
	putF32(data, o.TireTempFR, 83)
	putF32(data, o.TireTempRL, 89)
	putF32(data, o.TireTempRR, 90)
	data[o.Accel] = byte(throttle)
	data[o.Brake] = 0
	data[o.Gear] = 3
	data[o.Steer] = byte(int8(24 * math.Sin(seconds)))
	return data
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
