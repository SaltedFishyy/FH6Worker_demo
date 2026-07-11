package telemetry

import (
	"encoding/binary"
	"fmt"
	"math"
	"time"
)

type Parser struct {
	spec PacketSpec
}

func NewParser(spec PacketSpec) Parser {
	if spec.Length == 0 {
		spec = DefaultPacketSpec()
	}
	return Parser{spec: spec}
}

func DefaultParser() Parser {
	return NewParser(DefaultPacketSpec())
}

func (p Parser) Parse(data []byte, receivedAt time.Time) (TelemetryFrame, error) {
	if len(data) != p.spec.Length {
		return TelemetryFrame{}, fmt.Errorf("invalid telemetry packet length: got %d bytes, want %d", len(data), p.spec.Length)
	}

	r := packetReader{data: data}
	o := p.spec.Offsets

	frame := TelemetryFrame{
		ReceivedAt:                  receivedAt,
		IsRaceOn:                    r.i32(o.IsRaceOn),
		TimestampMS:                 r.u32(o.TimestampMS),
		EngineMaxRpm:                r.f32(o.EngineMaxRpm),
		EngineIdleRpm:               r.f32(o.EngineIdleRpm),
		CurrentEngineRpm:            r.f32(o.CurrentEngineRpm),
		AccelerationX:               r.f32(o.AccelerationX),
		AccelerationY:               r.f32(o.AccelerationY),
		AccelerationZ:               r.f32(o.AccelerationZ),
		VelocityX:                   r.f32(o.VelocityX),
		VelocityY:                   r.f32(o.VelocityY),
		VelocityZ:                   r.f32(o.VelocityZ),
		AngularVelocityX:            r.f32(o.AngularVelocityX),
		AngularVelocityY:            r.f32(o.AngularVelocityY),
		AngularVelocityZ:            r.f32(o.AngularVelocityZ),
		Yaw:                         r.f32(o.Yaw),
		Pitch:                       r.f32(o.Pitch),
		Roll:                        r.f32(o.Roll),
		Speed:                       r.f32(o.Speed),
		Power:                       r.f32(o.Power),
		Torque:                      r.f32(o.Torque),
		CarOrdinal:                  r.i32(o.CarOrdinal),
		CarClass:                    r.i32(o.CarClass),
		CarPI:                       r.i32(o.CarPI),
		DrivetrainType:              r.i32(o.DrivetrainType),
		NumCylinders:                r.i32(o.NumCylinders),
		CarCategory:                 r.u32(o.CarCategory),
		SmashableVelDiff:            r.f32(o.SmashableVelDiff),
		SmashableMass:               r.f32(o.SmashableMass),
		PositionX:                   r.f32(o.PositionX),
		PositionY:                   r.f32(o.PositionY),
		PositionZ:                   r.f32(o.PositionZ),
		Boost:                       r.f32(o.Boost),
		Fuel:                        r.f32(o.Fuel),
		DistanceTraveled:            r.f32(o.DistanceTraveled),
		BestLap:                     r.f32(o.BestLap),
		LastLap:                     r.f32(o.LastLap),
		CurrentLap:                  r.f32(o.CurrentLap),
		CurrentRaceTime:             r.f32(o.CurrentRaceTime),
		LapNumber:                   r.u16(o.LapNumber),
		RacePosition:                r.u8(o.RacePosition),
		Accel:                       r.u8(o.Accel),
		Brake:                       r.u8(o.Brake),
		Clutch:                      r.u8(o.Clutch),
		HandBrake:                   r.u8(o.HandBrake),
		Gear:                        r.u8(o.Gear),
		Steer:                       r.i8(o.Steer),
		NormalizedDrivingLine:       r.i8(o.NormalizedDrivingLine),
		NormalizedAIBrakeDifference: r.i8(o.NormalizedAIBrakeDifference),
		WheelFL:                     r.wheel(o.SuspensionTravelFL, o.TireSlipRatioFL, o.WheelRotationSpeedFL, o.WheelOnRumbleStripFL, o.WheelInPuddleDepthFL, o.SurfaceRumbleFL, o.TireSlipAngleFL, o.TireCombinedSlipFL, o.SuspensionTravelMetersFL, o.TireTempFL),
		WheelFR:                     r.wheel(o.SuspensionTravelFR, o.TireSlipRatioFR, o.WheelRotationSpeedFR, o.WheelOnRumbleStripFR, o.WheelInPuddleDepthFR, o.SurfaceRumbleFR, o.TireSlipAngleFR, o.TireCombinedSlipFR, o.SuspensionTravelMetersFR, o.TireTempFR),
		WheelRL:                     r.wheel(o.SuspensionTravelRL, o.TireSlipRatioRL, o.WheelRotationSpeedRL, o.WheelOnRumbleStripRL, o.WheelInPuddleDepthRL, o.SurfaceRumbleRL, o.TireSlipAngleRL, o.TireCombinedSlipRL, o.SuspensionTravelMetersRL, o.TireTempRL),
		WheelRR:                     r.wheel(o.SuspensionTravelRR, o.TireSlipRatioRR, o.WheelRotationSpeedRR, o.WheelOnRumbleStripRR, o.WheelInPuddleDepthRR, o.SurfaceRumbleRR, o.TireSlipAngleRR, o.TireCombinedSlipRR, o.SuspensionTravelMetersRR, o.TireTempRR),
	}

	if r.err != nil {
		return TelemetryFrame{}, r.err
	}

	return frame, nil
}

func NormalizeFrame(frame TelemetryFrame) NormalizedTelemetry {
	fl := normalizeWheel(frame.WheelFL)
	fr := normalizeWheel(frame.WheelFR)
	rl := normalizeWheel(frame.WheelRL)
	rr := normalizeWheel(frame.WheelRR)
	speedKmh, speedFieldKmh, velocitySpeedKmh, speedSource := normalizeSpeed(frame)

	return NormalizedTelemetry{
		ReceivedAt:       frame.ReceivedAt.UTC().Format(time.RFC3339Nano),
		TimeMS:           int64(frame.TimestampMS),
		IsRaceOn:         frame.IsRaceOn != 0,
		SpeedKmh:         speedKmh,
		SpeedFieldKmh:    speedFieldKmh,
		VelocitySpeedKmh: velocitySpeedKmh,
		SpeedSource:      speedSource,
		Rpm:              float64(frame.CurrentEngineRpm),
		RpmRatio:         rpmRatio(frame.EngineIdleRpm, frame.EngineMaxRpm, frame.CurrentEngineRpm),
		EngineMaxRpm:     float64(frame.EngineMaxRpm),
		EngineIdleRpm:    float64(frame.EngineIdleRpm),
		Gear:             int(frame.Gear),

		AccelerationX: float64(frame.AccelerationX),
		AccelerationY: float64(frame.AccelerationY),
		AccelerationZ: float64(frame.AccelerationZ),
		VelocityX:     float64(frame.VelocityX),
		VelocityY:     float64(frame.VelocityY),
		VelocityZ:     float64(frame.VelocityZ),
		Yaw:           float64(frame.Yaw),
		Pitch:         float64(frame.Pitch),
		Roll:          float64(frame.Roll),
		Power:         float64(frame.Power),
		Torque:        float64(frame.Torque),
		PositionX:     float64(frame.PositionX),
		PositionY:     float64(frame.PositionY),
		PositionZ:     float64(frame.PositionZ),

		Boost:            float64(frame.Boost),
		Fuel:             float64(frame.Fuel),
		DistanceTraveled: float64(frame.DistanceTraveled),
		BestLap:          float64(frame.BestLap),
		LastLap:          float64(frame.LastLap),
		CurrentLap:       float64(frame.CurrentLap),
		CurrentRaceTime:  float64(frame.CurrentRaceTime),
		LapNumber:        int(frame.LapNumber),
		RacePosition:     int(frame.RacePosition),
		SmashableVelDiff: float64(frame.SmashableVelDiff),
		SmashableMass:    float64(frame.SmashableMass),

		CarOrdinal:      int(frame.CarOrdinal),
		CarClassID:      int(frame.CarClass),
		CarClass:        carClassName(frame.CarClass),
		CarPI:           int(frame.CarPI),
		DrivetrainType:  int(frame.DrivetrainType),
		Drivetrain:      drivetrainName(frame.DrivetrainType),
		NumCylinders:    int(frame.NumCylinders),
		CarCategory:     int(frame.CarCategory),
		CarCategoryName: CarCategoryName(int(frame.CarCategory)),

		Throttle01:          input01(frame.Accel),
		Brake01:             input01(frame.Brake),
		Clutch01:            input01(frame.Clutch),
		HandBrake01:         input01(frame.HandBrake),
		Steer01:             steer01(frame.Steer),
		DrivingLine01:       steer01(frame.NormalizedDrivingLine),
		AIBrakeDifference01: steer01(frame.NormalizedAIBrakeDifference),

		FrontSlipRatioAvg:    avg(fl.SlipRatio, fr.SlipRatio),
		RearSlipRatioAvg:     avg(rl.SlipRatio, rr.SlipRatio),
		FrontSlipAngleAvg:    avg(fl.SlipAngle, fr.SlipAngle),
		RearSlipAngleAvg:     avg(rl.SlipAngle, rr.SlipAngle),
		FrontCombinedSlipAvg: avg(fl.CombinedSlip, fr.CombinedSlip),
		RearCombinedSlipAvg:  avg(rl.CombinedSlip, rr.CombinedSlip),

		TireTempFrontAvg:   avg(fl.TireTemp, fr.TireTemp),
		TireTempRearAvg:    avg(rl.TireTemp, rr.TireTemp),
		SuspensionFrontAvg: avg(fl.SuspensionTravel, fr.SuspensionTravel),
		SuspensionRearAvg:  avg(rl.SuspensionTravel, rr.SuspensionTravel),

		PitchRate: float64(frame.AngularVelocityX),
		YawRate:   float64(frame.AngularVelocityY),
		RollRate:  float64(frame.AngularVelocityZ),

		WheelFL: fl,
		WheelFR: fr,
		WheelRL: rl,
		WheelRR: rr,
	}
}

type packetReader struct {
	data []byte
	err  error
}

func (r *packetReader) check(off int, size int) bool {
	if r.err != nil {
		return false
	}
	if off < 0 || off+size > len(r.data) {
		r.err = fmt.Errorf("packet spec reads outside packet: offset=%d size=%d length=%d", off, size, len(r.data))
		return false
	}
	return true
}

func (r *packetReader) f32(off int) float32 {
	if !r.check(off, 4) {
		return 0
	}
	return math.Float32frombits(binary.LittleEndian.Uint32(r.data[off : off+4]))
}

func (r *packetReader) i32(off int) int32 {
	if !r.check(off, 4) {
		return 0
	}
	return int32(binary.LittleEndian.Uint32(r.data[off : off+4]))
}

func (r *packetReader) u32(off int) uint32 {
	if !r.check(off, 4) {
		return 0
	}
	return binary.LittleEndian.Uint32(r.data[off : off+4])
}

func (r *packetReader) u16(off int) uint16 {
	if !r.check(off, 2) {
		return 0
	}
	return binary.LittleEndian.Uint16(r.data[off : off+2])
}

func (r *packetReader) u8(off int) uint8 {
	if !r.check(off, 1) {
		return 0
	}
	return r.data[off]
}

func (r *packetReader) i8(off int) int8 {
	if !r.check(off, 1) {
		return 0
	}
	return int8(r.data[off])
}

func (r *packetReader) wheel(suspensionOff, slipRatioOff, rotationOff, rumbleOff, puddleOff, surfaceRumbleOff, slipAngleOff, combinedSlipOff, suspensionMetersOff, tempOff int) WheelFrame {
	return WheelFrame{
		SuspensionTravel:       r.f32(suspensionOff),
		TireSlipRatio:          r.f32(slipRatioOff),
		WheelRotationRate:      r.f32(rotationOff),
		OnRumbleStrip:          r.i32(rumbleOff),
		InPuddleDepth:          r.f32(puddleOff),
		SurfaceRumble:          r.f32(surfaceRumbleOff),
		TireSlipAngle:          r.f32(slipAngleOff),
		TireCombinedSlip:       r.f32(combinedSlipOff),
		SuspensionTravelMeters: r.f32(suspensionMetersOff),
		TireTemp:               r.f32(tempOff),
	}
}

func normalizeWheel(w WheelFrame) NormalizedWheelTelemetry {
	return NormalizedWheelTelemetry{
		SlipRatio:              float64(w.TireSlipRatio),
		SlipAngle:              float64(w.TireSlipAngle),
		CombinedSlip:           float64(w.TireCombinedSlip),
		TireTemp:               float64(w.TireTemp),
		SuspensionTravel:       float64(w.SuspensionTravel),
		SuspensionTravelMeters: float64(w.SuspensionTravelMeters),
		WheelRotationSpeed:     float64(w.WheelRotationRate),
		RumbleStrip:            float64(w.OnRumbleStrip),
		PuddleDepth:            float64(w.InPuddleDepth),
		SurfaceRumble:          float64(w.SurfaceRumble),
	}
}

func input01(value uint8) float64 {
	return clamp01(float64(value) / 255.0)
}

func steer01(value int8) float64 {
	if value < 0 {
		return math.Max(-1, float64(value)/128.0)
	}
	return math.Min(1, float64(value)/127.0)
}

func rpmRatio(idle, max, current float32) float64 {
	if max <= idle {
		return 0
	}
	return clamp01(float64(current-idle) / float64(max-idle))
}

func normalizeSpeed(frame TelemetryFrame) (valueKmh, speedFieldKmh, velocitySpeedKmh float64, source string) {
	speedFieldKmh = float64(frame.Speed) * 3.6
	velocitySpeedKmh = math.Sqrt(
		float64(frame.VelocityX)*float64(frame.VelocityX)+
			float64(frame.VelocityY)*float64(frame.VelocityY)+
			float64(frame.VelocityZ)*float64(frame.VelocityZ),
	) * 3.6

	fieldValid := plausibleSpeed(speedFieldKmh)
	velocityValid := plausibleSpeed(velocitySpeedKmh)
	if velocityValid && velocitySpeedKmh > 0.5 && (!fieldValid || speedDisagrees(speedFieldKmh, velocitySpeedKmh)) {
		return velocitySpeedKmh, speedFieldKmh, velocitySpeedKmh, "velocity"
	}
	if fieldValid {
		return speedFieldKmh, speedFieldKmh, velocitySpeedKmh, "packet"
	}
	if velocityValid {
		return velocitySpeedKmh, speedFieldKmh, velocitySpeedKmh, "velocity"
	}
	return 0, speedFieldKmh, velocitySpeedKmh, "none"
}

func plausibleSpeed(value float64) bool {
	return !math.IsNaN(value) && !math.IsInf(value, 0) && value >= 0 && value <= 700
}

func speedDisagrees(a, b float64) bool {
	diff := math.Abs(a - b)
	baseline := math.Max(math.Max(math.Abs(a), math.Abs(b)), 1)
	return diff > 25 && diff/baseline > 0.2
}

func carClassName(value int32) string {
	switch value {
	case 0:
		return "D"
	case 1:
		return "C"
	case 2:
		return "B"
	case 3:
		return "A"
	case 4:
		return "S1"
	case 5:
		return "S2"
	case 6:
		return "R"
	case 7:
		return "X"
	default:
		return ""
	}
}

func drivetrainName(value int32) string {
	switch value {
	case 0:
		return "FWD"
	case 1:
		return "RWD"
	case 2:
		return "AWD"
	default:
		return ""
	}
}

func CarCategoryName(value int) string {
	switch value {
	case 11:
		return "现代超级跑车"
	case 12:
		return "复古超级跑车"
	case 13:
		return "顶级超跑"
	case 14:
		return "复古超级轿车"
	case 16:
		return "多功能英雄车"
	case 17:
		return "复古跑车"
	case 18:
		return "现代跑车"
	case 19:
		return "现代超级轿车"
	case 20:
		return "经典赛车"
	case 21:
		return "狂热跑车"
	case 22:
		return "稀有经典车"
	case 23:
		return "高性能掀背车"
	case 24:
		return "复古高性能掀背车"
	case 25:
		return "超高性能掀背车"
	case 26:
		return "极限赛道玩具"
	case 28:
		return "经典肌肉车"
	case 29:
		return "改装车和特制车"
	case 30:
		return "复古肌肉车"
	case 31:
		return "现代肌肉车"
	case 32:
		return "复古拉力赛车"
	case 33:
		return "经典拉力赛车"
	case 34:
		return "拉力赛怪物"
	case 35:
		return "现代拉力赛车"
	case 36:
		return "GT赛车"
	case 37:
		return "超级豪华旅行车"
	case 38:
		return "无限制越野车"
	case 39:
		return "运动型多功能英雄车"
	case 40:
		return "越野车"
	case 41:
		return "无限制沙滩车"
	case 42:
		return "经典跑车"
	case 43:
		return "赛道玩具"
	case 46:
		return "沙滩车"
	case 47:
		return "漂移赛车"
	case 48:
		return "皮卡和四轮驱动车"
	case 49:
		return "多功能车"
	case 50:
		return "家庭多用"
	case 51:
		return "复古赛车"
	default:
		return ""
	}
}

func avg(a, b float64) float64 {
	return (a + b) / 2
}

func clamp01(value float64) float64 {
	if value < 0 {
		return 0
	}
	if value > 1 {
		return 1
	}
	return value
}
