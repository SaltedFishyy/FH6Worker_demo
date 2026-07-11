package telemetry

import (
	"sync"
	"time"
)

const defaultAggregateBucket = 100 * time.Millisecond

type RingBuffer struct {
	mu     sync.RWMutex
	window time.Duration
	items  []bufferItem
}

type bufferItem struct {
	at    time.Time
	frame NormalizedTelemetry
}

func NewRingBuffer(window time.Duration) *RingBuffer {
	if window <= 0 {
		window = 30 * time.Second
	}
	return &RingBuffer{window: window}
}

func (b *RingBuffer) Add(at time.Time, frame NormalizedTelemetry) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.items = append(b.items, bufferItem{at: at, frame: frame})
	b.pruneLocked(at.Add(-b.window))
}

func (b *RingBuffer) Reset() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.items = nil
}

func (b *RingBuffer) Latest() *NormalizedTelemetry {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if len(b.items) == 0 {
		return nil
	}
	frame := b.items[len(b.items)-1].frame
	return &frame
}

func (b *RingBuffer) Since(seconds int) []NormalizedTelemetry {
	if seconds <= 0 {
		return nil
	}
	return b.Aggregated(seconds, defaultAggregateBucket)
}

func (b *RingBuffer) Aggregated(seconds int, bucketSize time.Duration) []NormalizedTelemetry {
	if seconds <= 0 {
		return nil
	}
	if bucketSize <= 0 {
		bucketSize = defaultAggregateBucket
	}

	b.mu.RLock()
	defer b.mu.RUnlock()

	if len(b.items) == 0 {
		return nil
	}

	cutoff := time.Now().Add(-time.Duration(seconds) * time.Second)
	var buckets []aggregateBucket
	for _, item := range b.items {
		if item.at.Before(cutoff) {
			continue
		}
		if len(buckets) == 0 || item.at.Sub(buckets[len(buckets)-1].start) >= bucketSize {
			buckets = append(buckets, aggregateBucket{start: item.at})
		}
		buckets[len(buckets)-1].add(item.frame)
	}

	result := make([]NormalizedTelemetry, 0, len(buckets))
	for _, bucket := range buckets {
		if bucket.count == 0 {
			continue
		}
		result = append(result, bucket.frame())
	}
	return result
}

func (b *RingBuffer) Len() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.items)
}

func (b *RingBuffer) pruneLocked(cutoff time.Time) {
	idx := 0
	for idx < len(b.items) && b.items[idx].at.Before(cutoff) {
		idx++
	}
	if idx == 0 {
		return
	}
	copy(b.items, b.items[idx:])
	b.items = b.items[:len(b.items)-idx]
}

type aggregateBucket struct {
	start time.Time
	count float64
	acc   NormalizedTelemetry
}

func (b *aggregateBucket) add(frame NormalizedTelemetry) {
	b.count++
	b.acc.ReceivedAt = frame.ReceivedAt
	b.acc.TimeMS = frame.TimeMS
	b.acc.IsRaceOn = frame.IsRaceOn
	b.acc.GameMode = frame.GameMode
	b.acc.Gear = frame.Gear
	b.acc.SpeedSource = frame.SpeedSource
	b.acc.CarOrdinal = frame.CarOrdinal
	b.acc.CarClassID = frame.CarClassID
	b.acc.CarClass = frame.CarClass
	b.acc.CarPI = frame.CarPI
	b.acc.DrivetrainType = frame.DrivetrainType
	b.acc.Drivetrain = frame.Drivetrain
	b.acc.NumCylinders = frame.NumCylinders
	b.acc.CarCategory = frame.CarCategory
	b.acc.CarCategoryName = frame.CarCategoryName
	b.acc.PositionX = frame.PositionX
	b.acc.PositionY = frame.PositionY
	b.acc.PositionZ = frame.PositionZ
	b.acc.DistanceTraveled = frame.DistanceTraveled
	b.acc.BestLap = frame.BestLap
	b.acc.LastLap = frame.LastLap
	b.acc.CurrentLap = frame.CurrentLap
	b.acc.CurrentRaceTime = frame.CurrentRaceTime
	b.acc.LapNumber = frame.LapNumber
	b.acc.RacePosition = frame.RacePosition
	b.acc.SmashableVelDiff = frame.SmashableVelDiff
	b.acc.SmashableMass = frame.SmashableMass

	b.acc.SpeedKmh += frame.SpeedKmh
	b.acc.SpeedFieldKmh += frame.SpeedFieldKmh
	b.acc.VelocitySpeedKmh += frame.VelocitySpeedKmh
	b.acc.Rpm += frame.Rpm
	b.acc.RpmRatio += frame.RpmRatio
	b.acc.EngineMaxRpm += frame.EngineMaxRpm
	b.acc.EngineIdleRpm += frame.EngineIdleRpm
	b.acc.AccelerationX += frame.AccelerationX
	b.acc.AccelerationY += frame.AccelerationY
	b.acc.AccelerationZ += frame.AccelerationZ
	b.acc.VelocityX += frame.VelocityX
	b.acc.VelocityY += frame.VelocityY
	b.acc.VelocityZ += frame.VelocityZ
	b.acc.Yaw += frame.Yaw
	b.acc.Pitch += frame.Pitch
	b.acc.Roll += frame.Roll
	b.acc.Power += frame.Power
	b.acc.Torque += frame.Torque
	b.acc.Boost += frame.Boost
	b.acc.Fuel += frame.Fuel
	b.acc.Throttle01 += frame.Throttle01
	b.acc.Brake01 += frame.Brake01
	b.acc.Clutch01 += frame.Clutch01
	b.acc.HandBrake01 += frame.HandBrake01
	b.acc.Steer01 += frame.Steer01
	b.acc.DrivingLine01 += frame.DrivingLine01
	b.acc.AIBrakeDifference01 += frame.AIBrakeDifference01
	b.acc.FrontSlipRatioAvg += frame.FrontSlipRatioAvg
	b.acc.RearSlipRatioAvg += frame.RearSlipRatioAvg
	b.acc.FrontSlipAngleAvg += frame.FrontSlipAngleAvg
	b.acc.RearSlipAngleAvg += frame.RearSlipAngleAvg
	b.acc.FrontCombinedSlipAvg += frame.FrontCombinedSlipAvg
	b.acc.RearCombinedSlipAvg += frame.RearCombinedSlipAvg
	b.acc.TireTempFrontAvg += frame.TireTempFrontAvg
	b.acc.TireTempRearAvg += frame.TireTempRearAvg
	b.acc.SuspensionFrontAvg += frame.SuspensionFrontAvg
	b.acc.SuspensionRearAvg += frame.SuspensionRearAvg
	b.acc.YawRate += frame.YawRate
	b.acc.PitchRate += frame.PitchRate
	b.acc.RollRate += frame.RollRate
	b.acc.WheelFL = addWheel(b.acc.WheelFL, frame.WheelFL)
	b.acc.WheelFR = addWheel(b.acc.WheelFR, frame.WheelFR)
	b.acc.WheelRL = addWheel(b.acc.WheelRL, frame.WheelRL)
	b.acc.WheelRR = addWheel(b.acc.WheelRR, frame.WheelRR)
}

func (b aggregateBucket) frame() NormalizedTelemetry {
	frame := b.acc
	div := b.count
	frame.SpeedKmh /= div
	frame.SpeedFieldKmh /= div
	frame.VelocitySpeedKmh /= div
	frame.Rpm /= div
	frame.RpmRatio /= div
	frame.EngineMaxRpm /= div
	frame.EngineIdleRpm /= div
	frame.AccelerationX /= div
	frame.AccelerationY /= div
	frame.AccelerationZ /= div
	frame.VelocityX /= div
	frame.VelocityY /= div
	frame.VelocityZ /= div
	frame.Yaw /= div
	frame.Pitch /= div
	frame.Roll /= div
	frame.Power /= div
	frame.Torque /= div
	frame.Boost /= div
	frame.Fuel /= div
	frame.Throttle01 /= div
	frame.Brake01 /= div
	frame.Clutch01 /= div
	frame.HandBrake01 /= div
	frame.Steer01 /= div
	frame.DrivingLine01 /= div
	frame.AIBrakeDifference01 /= div
	frame.FrontSlipRatioAvg /= div
	frame.RearSlipRatioAvg /= div
	frame.FrontSlipAngleAvg /= div
	frame.RearSlipAngleAvg /= div
	frame.FrontCombinedSlipAvg /= div
	frame.RearCombinedSlipAvg /= div
	frame.TireTempFrontAvg /= div
	frame.TireTempRearAvg /= div
	frame.SuspensionFrontAvg /= div
	frame.SuspensionRearAvg /= div
	frame.YawRate /= div
	frame.PitchRate /= div
	frame.RollRate /= div
	frame.WheelFL = divideWheel(frame.WheelFL, div)
	frame.WheelFR = divideWheel(frame.WheelFR, div)
	frame.WheelRL = divideWheel(frame.WheelRL, div)
	frame.WheelRR = divideWheel(frame.WheelRR, div)
	return frame
}

func addWheel(a, b NormalizedWheelTelemetry) NormalizedWheelTelemetry {
	return NormalizedWheelTelemetry{
		SlipRatio:              a.SlipRatio + b.SlipRatio,
		SlipAngle:              a.SlipAngle + b.SlipAngle,
		CombinedSlip:           a.CombinedSlip + b.CombinedSlip,
		TireTemp:               a.TireTemp + b.TireTemp,
		SuspensionTravel:       a.SuspensionTravel + b.SuspensionTravel,
		SuspensionTravelMeters: a.SuspensionTravelMeters + b.SuspensionTravelMeters,
		WheelRotationSpeed:     a.WheelRotationSpeed + b.WheelRotationSpeed,
		RumbleStrip:            a.RumbleStrip + b.RumbleStrip,
		PuddleDepth:            a.PuddleDepth + b.PuddleDepth,
		SurfaceRumble:          a.SurfaceRumble + b.SurfaceRumble,
	}
}

func divideWheel(w NormalizedWheelTelemetry, div float64) NormalizedWheelTelemetry {
	return NormalizedWheelTelemetry{
		SlipRatio:              w.SlipRatio / div,
		SlipAngle:              w.SlipAngle / div,
		CombinedSlip:           w.CombinedSlip / div,
		TireTemp:               w.TireTemp / div,
		SuspensionTravel:       w.SuspensionTravel / div,
		SuspensionTravelMeters: w.SuspensionTravelMeters / div,
		WheelRotationSpeed:     w.WheelRotationSpeed / div,
		RumbleStrip:            w.RumbleStrip / div,
		PuddleDepth:            w.PuddleDepth / div,
		SurfaceRumble:          w.SurfaceRumble / div,
	}
}
