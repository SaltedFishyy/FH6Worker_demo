package telemetry

type PacketSpec struct {
	Name    string
	Length  int
	Offsets PacketOffsets
}

type PacketOffsets struct {
	IsRaceOn         int
	TimestampMS      int
	EngineMaxRpm     int
	EngineIdleRpm    int
	CurrentEngineRpm int

	AccelerationX int
	AccelerationY int
	AccelerationZ int
	VelocityX     int
	VelocityY     int
	VelocityZ     int

	AngularVelocityX int
	AngularVelocityY int
	AngularVelocityZ int
	Yaw              int
	Pitch            int
	Roll             int

	SuspensionTravelFL int
	SuspensionTravelFR int
	SuspensionTravelRL int
	SuspensionTravelRR int

	TireSlipRatioFL int
	TireSlipRatioFR int
	TireSlipRatioRL int
	TireSlipRatioRR int

	WheelRotationSpeedFL int
	WheelRotationSpeedFR int
	WheelRotationSpeedRL int
	WheelRotationSpeedRR int

	WheelOnRumbleStripFL int
	WheelOnRumbleStripFR int
	WheelOnRumbleStripRL int
	WheelOnRumbleStripRR int

	WheelInPuddleDepthFL int
	WheelInPuddleDepthFR int
	WheelInPuddleDepthRL int
	WheelInPuddleDepthRR int

	SurfaceRumbleFL int
	SurfaceRumbleFR int
	SurfaceRumbleRL int
	SurfaceRumbleRR int

	TireSlipAngleFL int
	TireSlipAngleFR int
	TireSlipAngleRL int
	TireSlipAngleRR int

	TireCombinedSlipFL int
	TireCombinedSlipFR int
	TireCombinedSlipRL int
	TireCombinedSlipRR int

	SuspensionTravelMetersFL int
	SuspensionTravelMetersFR int
	SuspensionTravelMetersRL int
	SuspensionTravelMetersRR int

	CarOrdinal     int
	CarClass       int
	CarPI          int
	DrivetrainType int
	NumCylinders   int
	CarCategory    int

	SmashableVelDiff int
	SmashableMass    int

	PositionX int
	PositionY int
	PositionZ int
	Speed     int
	Power     int
	Torque    int

	TireTempFL int
	TireTempFR int
	TireTempRL int
	TireTempRR int

	Boost            int
	Fuel             int
	DistanceTraveled int

	BestLap         int
	LastLap         int
	CurrentLap      int
	CurrentRaceTime int
	LapNumber       int
	RacePosition    int

	Accel                       int
	Brake                       int
	Clutch                      int
	HandBrake                   int
	Gear                        int
	Steer                       int
	NormalizedDrivingLine       int
	NormalizedAIBrakeDifference int
}

func DefaultPacketSpec() PacketSpec {
	return PacketSpec{
		Name:   "forza-horizon-data-out-324",
		Length: PacketLength,
		Offsets: PacketOffsets{
			IsRaceOn:         0,
			TimestampMS:      4,
			EngineMaxRpm:     8,
			EngineIdleRpm:    12,
			CurrentEngineRpm: 16,

			AccelerationX: 20,
			AccelerationY: 24,
			AccelerationZ: 28,
			VelocityX:     32,
			VelocityY:     36,
			VelocityZ:     40,

			AngularVelocityX: 44,
			AngularVelocityY: 48,
			AngularVelocityZ: 52,
			Yaw:              56,
			Pitch:            60,
			Roll:             64,

			SuspensionTravelFL: 68,
			SuspensionTravelFR: 72,
			SuspensionTravelRL: 76,
			SuspensionTravelRR: 80,

			TireSlipRatioFL: 84,
			TireSlipRatioFR: 88,
			TireSlipRatioRL: 92,
			TireSlipRatioRR: 96,

			WheelRotationSpeedFL: 100,
			WheelRotationSpeedFR: 104,
			WheelRotationSpeedRL: 108,
			WheelRotationSpeedRR: 112,

			WheelOnRumbleStripFL: 116,
			WheelOnRumbleStripFR: 120,
			WheelOnRumbleStripRL: 124,
			WheelOnRumbleStripRR: 128,

			WheelInPuddleDepthFL: 132,
			WheelInPuddleDepthFR: 136,
			WheelInPuddleDepthRL: 140,
			WheelInPuddleDepthRR: 144,

			SurfaceRumbleFL: 148,
			SurfaceRumbleFR: 152,
			SurfaceRumbleRL: 156,
			SurfaceRumbleRR: 160,

			TireSlipAngleFL: 164,
			TireSlipAngleFR: 168,
			TireSlipAngleRL: 172,
			TireSlipAngleRR: 176,

			TireCombinedSlipFL: 180,
			TireCombinedSlipFR: 184,
			TireCombinedSlipRL: 188,
			TireCombinedSlipRR: 192,

			SuspensionTravelMetersFL: 196,
			SuspensionTravelMetersFR: 200,
			SuspensionTravelMetersRL: 204,
			SuspensionTravelMetersRR: 208,

			CarOrdinal:     212,
			CarClass:       216,
			CarPI:          220,
			DrivetrainType: 224,
			NumCylinders:   228,
			CarCategory:    232,

			SmashableVelDiff: 236,
			SmashableMass:    240,

			PositionX: 244,
			PositionY: 248,
			PositionZ: 252,
			Speed:     256,
			Power:     260,
			Torque:    264,

			TireTempFL: 268,
			TireTempFR: 272,
			TireTempRL: 276,
			TireTempRR: 280,

			Boost:            284,
			Fuel:             288,
			DistanceTraveled: 292,

			BestLap:         296,
			LastLap:         300,
			CurrentLap:      304,
			CurrentRaceTime: 308,
			LapNumber:       312,
			RacePosition:    314,

			Accel:                       315,
			Brake:                       316,
			Clutch:                      317,
			HandBrake:                   318,
			Gear:                        319,
			Steer:                       320,
			NormalizedDrivingLine:       321,
			NormalizedAIBrakeDifference: 322,
		},
	}
}
