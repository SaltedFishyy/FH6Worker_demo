package telemetry

import "time"

const (
	DefaultAddress = "0.0.0.0"
	DefaultPort    = 5301
	PacketLength   = 324
)

type RawPacket struct {
	Timestamp time.Time
	Data      []byte
	Addr      string
}

type WheelFrame struct {
	SuspensionTravel       float32
	TireSlipRatio          float32
	WheelRotationRate      float32
	OnRumbleStrip          int32
	InPuddleDepth          float32
	SurfaceRumble          float32
	TireSlipAngle          float32
	TireCombinedSlip       float32
	SuspensionTravelMeters float32
	TireTemp               float32
}

type TelemetryFrame struct {
	ReceivedAt  time.Time
	TimestampMS uint32
	IsRaceOn    int32

	EngineMaxRpm     float32
	EngineIdleRpm    float32
	CurrentEngineRpm float32

	AccelerationX float32
	AccelerationY float32
	AccelerationZ float32

	VelocityX float32
	VelocityY float32
	VelocityZ float32

	AngularVelocityX float32
	AngularVelocityY float32
	AngularVelocityZ float32

	Yaw   float32
	Pitch float32
	Roll  float32

	WheelFL WheelFrame
	WheelFR WheelFrame
	WheelRL WheelFrame
	WheelRR WheelFrame

	Speed  float32
	Power  float32
	Torque float32

	CarOrdinal       int32
	CarClass         int32
	CarPI            int32
	DrivetrainType   int32
	NumCylinders     int32
	CarCategory      uint32
	SmashableVelDiff float32
	SmashableMass    float32
	PositionX        float32
	PositionY        float32
	PositionZ        float32
	Boost            float32
	Fuel             float32
	DistanceTraveled float32
	BestLap          float32
	LastLap          float32
	CurrentLap       float32
	CurrentRaceTime  float32
	LapNumber        uint16
	RacePosition     uint8

	Accel                       uint8
	Brake                       uint8
	Clutch                      uint8
	HandBrake                   uint8
	Gear                        uint8
	Steer                       int8
	NormalizedDrivingLine       int8
	NormalizedAIBrakeDifference int8
}

type NormalizedWheelTelemetry struct {
	SlipRatio              float64 `json:"slipRatio"`
	SlipAngle              float64 `json:"slipAngle"`
	CombinedSlip           float64 `json:"combinedSlip"`
	TireTemp               float64 `json:"tireTemp"`
	SuspensionTravel       float64 `json:"suspensionTravel"`
	SuspensionTravelMeters float64 `json:"suspensionTravelMeters"`
	WheelRotationSpeed     float64 `json:"wheelRotationSpeed"`
	RumbleStrip            float64 `json:"rumbleStrip"`
	PuddleDepth            float64 `json:"puddleDepth"`
	SurfaceRumble          float64 `json:"surfaceRumble"`
}

type NormalizedTelemetry struct {
	ReceivedAt string `json:"receivedAt"`
	TimeMS     int64  `json:"timeMs"`
	IsRaceOn   bool   `json:"isRaceOn"`
	GameMode   string `json:"gameMode"`

	SpeedKmh         float64 `json:"speedKmh"`
	SpeedFieldKmh    float64 `json:"speedFieldKmh"`
	VelocitySpeedKmh float64 `json:"velocitySpeedKmh"`
	SpeedSource      string  `json:"speedSource"`
	Rpm              float64 `json:"rpm"`
	RpmRatio         float64 `json:"rpmRatio"`
	EngineMaxRpm     float64 `json:"engineMaxRpm"`
	EngineIdleRpm    float64 `json:"engineIdleRpm"`
	Gear             int     `json:"gear"`

	AccelerationX float64 `json:"accelerationX"`
	AccelerationY float64 `json:"accelerationY"`
	AccelerationZ float64 `json:"accelerationZ"`
	VelocityX     float64 `json:"velocityX"`
	VelocityY     float64 `json:"velocityY"`
	VelocityZ     float64 `json:"velocityZ"`
	Yaw           float64 `json:"yaw"`
	Pitch         float64 `json:"pitch"`
	Roll          float64 `json:"roll"`
	Power         float64 `json:"power"`
	Torque        float64 `json:"torque"`
	PositionX     float64 `json:"positionX"`
	PositionY     float64 `json:"positionY"`
	PositionZ     float64 `json:"positionZ"`

	Boost            float64 `json:"boost"`
	Fuel             float64 `json:"fuel"`
	DistanceTraveled float64 `json:"distanceTraveled"`
	BestLap          float64 `json:"bestLap"`
	LastLap          float64 `json:"lastLap"`
	CurrentLap       float64 `json:"currentLap"`
	CurrentRaceTime  float64 `json:"currentRaceTime"`
	LapNumber        int     `json:"lapNumber"`
	RacePosition     int     `json:"racePosition"`
	SmashableVelDiff float64 `json:"smashableVelDiff"`
	SmashableMass    float64 `json:"smashableMass"`

	CarOrdinal      int    `json:"carOrdinal"`
	CarClassID      int    `json:"carClassId"`
	CarClass        string `json:"carClass"`
	CarPI           int    `json:"carPi"`
	DrivetrainType  int    `json:"drivetrainType"`
	Drivetrain      string `json:"drivetrain"`
	NumCylinders    int    `json:"numCylinders"`
	CarCategory     int    `json:"carCategory"`
	CarCategoryName string `json:"carCategoryName"`

	Throttle01          float64 `json:"throttle01"`
	Brake01             float64 `json:"brake01"`
	Clutch01            float64 `json:"clutch01"`
	HandBrake01         float64 `json:"handBrake01"`
	Steer01             float64 `json:"steer01"`
	DrivingLine01       float64 `json:"drivingLine01"`
	AIBrakeDifference01 float64 `json:"aiBrakeDifference01"`

	FrontSlipRatioAvg    float64 `json:"frontSlipRatioAvg"`
	RearSlipRatioAvg     float64 `json:"rearSlipRatioAvg"`
	FrontSlipAngleAvg    float64 `json:"frontSlipAngleAvg"`
	RearSlipAngleAvg     float64 `json:"rearSlipAngleAvg"`
	FrontCombinedSlipAvg float64 `json:"frontCombinedSlipAvg"`
	RearCombinedSlipAvg  float64 `json:"rearCombinedSlipAvg"`

	TireTempFrontAvg   float64 `json:"tireTempFrontAvg"`
	TireTempRearAvg    float64 `json:"tireTempRearAvg"`
	SuspensionFrontAvg float64 `json:"suspensionFrontAvg"`
	SuspensionRearAvg  float64 `json:"suspensionRearAvg"`

	YawRate   float64 `json:"yawRate"`
	PitchRate float64 `json:"pitchRate"`
	RollRate  float64 `json:"rollRate"`

	WheelFL NormalizedWheelTelemetry `json:"wheelFL"`
	WheelFR NormalizedWheelTelemetry `json:"wheelFR"`
	WheelRL NormalizedWheelTelemetry `json:"wheelRL"`
	WheelRR NormalizedWheelTelemetry `json:"wheelRR"`
}

type TelemetryStatus struct {
	Running             bool   `json:"running"`
	Mode                string `json:"mode"`
	AnalysisMode        string `json:"analysisMode"`
	Address             string `json:"address"`
	Port                int    `json:"port"`
	PacketLength        int    `json:"packetLength"`
	RawPackets          uint64 `json:"rawPackets"`
	ValidPackets        uint64 `json:"validPackets"`
	InvalidPackets      uint64 `json:"invalidPackets"`
	ParseErrors         uint64 `json:"parseErrors"`
	LastDatagramAt      string `json:"lastDatagramAt"`
	LastDatagramBytes   int    `json:"lastDatagramBytes"`
	LastDatagramRemote  string `json:"lastDatagramRemote"`
	LastPacketAt        string `json:"lastPacketAt"`
	LastError           string `json:"lastError"`
	HasCurrentFrame     bool   `json:"hasCurrentFrame"`
	RecordingActive     bool   `json:"recordingActive"`
	RecordingBytes      int64  `json:"recordingBytes"`
	RecordingLimitBytes int64  `json:"recordingLimitBytes"`
	RecordingPackets    int64  `json:"recordingPackets"`
	RecordingTruncated  bool   `json:"recordingTruncated"`
}

type TelemetryReplayStatus struct {
	Running     bool    `json:"running"`
	Paused      bool    `json:"paused"`
	SessionID   int64   `json:"sessionId"`
	Speed       float64 `json:"speed"`
	PositionMS  int64   `json:"positionMs"`
	DurationMS  int64   `json:"durationMs"`
	Progress01  float64 `json:"progress01"`
	PacketIndex int     `json:"packetIndex"`
	PacketCount int     `json:"packetCount"`
	LastError   string  `json:"lastError"`
}

type TelemetrySummary struct {
	SampleCount int64   `json:"sampleCount"`
	AvgSpeedKmh float64 `json:"avgSpeedKmh"`
	MaxSpeedKmh float64 `json:"maxSpeedKmh"`
}

type SuggestedAction struct {
	Priority  int    `json:"priority"`
	Category  string `json:"category"`
	Item      string `json:"item"`
	Direction string `json:"direction"`
	Amount    string `json:"amount"`
	Reason    string `json:"reason"`
}

type DetectedEvent struct {
	ID               string             `json:"id"`
	Type             string             `json:"type"`
	Severity         string             `json:"severity"`
	StartMS          int64              `json:"startMs"`
	EndMS            int64              `json:"endMs"`
	DurationMS       int64              `json:"durationMs"`
	Segment          string             `json:"segment"`
	Evidence         map[string]float64 `json:"evidence"`
	SuggestedActions []SuggestedAction  `json:"suggestedActions"`
}

type NetworkInterface struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Address     string `json:"address"`
	IsLoopback  bool   `json:"isLoopback"`
	IsPrivate   bool   `json:"isPrivate"`
	IsUp        bool   `json:"isUp"`
}
