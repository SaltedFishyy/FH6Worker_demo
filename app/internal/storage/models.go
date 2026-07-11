package storage

type TuneProfile struct {
	ID        int64  `json:"id"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`

	CarName      string `json:"carName"`
	CarOrdinal   *int64 `json:"carOrdinal,omitempty"`
	CarCategory  *int64 `json:"carCategory,omitempty"`
	CarClass     string `json:"carClass"`
	PI           *int64 `json:"pi,omitempty"`
	Drivetrain   string `json:"drivetrain"`
	NumCylinders *int64 `json:"numCylinders,omitempty"`
	UseCase      string `json:"useCase"`
	VersionName  string `json:"versionName"`

	PowerKW              *float64 `json:"powerKW,omitempty"`
	TorqueNM             *float64 `json:"torqueNM,omitempty"`
	WeightKG             *float64 `json:"weightKG,omitempty"`
	FrontWeightPct       *float64 `json:"frontWeightPct,omitempty"`
	PowerToWeightKWPerKG *float64 `json:"powerToWeightKWPerKG,omitempty"`
	PeakTorqueRPM        *float64 `json:"peakTorqueRPM,omitempty"`
	PeakPowerRPM         *float64 `json:"peakPowerRPM,omitempty"`
	RedlineRPM           *float64 `json:"redlineRPM,omitempty"`

	FrontTirePressure *float64 `json:"frontTirePressure,omitempty"`
	RearTirePressure  *float64 `json:"rearTirePressure,omitempty"`

	FinalDrive *float64 `json:"finalDrive,omitempty"`
	Gear1      *float64 `json:"gear1,omitempty"`
	Gear2      *float64 `json:"gear2,omitempty"`
	Gear3      *float64 `json:"gear3,omitempty"`
	Gear4      *float64 `json:"gear4,omitempty"`
	Gear5      *float64 `json:"gear5,omitempty"`
	Gear6      *float64 `json:"gear6,omitempty"`
	Gear7      *float64 `json:"gear7,omitempty"`
	Gear8      *float64 `json:"gear8,omitempty"`
	Gear9      *float64 `json:"gear9,omitempty"`
	Gear10     *float64 `json:"gear10,omitempty"`

	FrontCamber *float64 `json:"frontCamber,omitempty"`
	RearCamber  *float64 `json:"rearCamber,omitempty"`
	FrontToe    *float64 `json:"frontToe,omitempty"`
	RearToe     *float64 `json:"rearToe,omitempty"`
	Caster      *float64 `json:"caster,omitempty"`

	FrontARB *float64 `json:"frontArb,omitempty"`
	RearARB  *float64 `json:"rearArb,omitempty"`

	FrontSpring     *float64 `json:"frontSpring,omitempty"`
	RearSpring      *float64 `json:"rearSpring,omitempty"`
	FrontRideHeight *float64 `json:"frontRideHeight,omitempty"`
	RearRideHeight  *float64 `json:"rearRideHeight,omitempty"`

	FrontRebound *float64 `json:"frontRebound,omitempty"`
	RearRebound  *float64 `json:"rearRebound,omitempty"`
	FrontBump    *float64 `json:"frontBump,omitempty"`
	RearBump     *float64 `json:"rearBump,omitempty"`

	FrontAero   *float64 `json:"frontAero,omitempty"`
	RearAero    *float64 `json:"rearAero,omitempty"`
	AeroBalance *float64 `json:"aeroBalance,omitempty"`

	BrakeBalance  *float64 `json:"brakeBalance,omitempty"`
	BrakePressure *float64 `json:"brakePressure,omitempty"`

	FrontDiffAccel    *float64 `json:"frontDiffAccel,omitempty"`
	FrontDiffDecel    *float64 `json:"frontDiffDecel,omitempty"`
	RearDiffAccel     *float64 `json:"rearDiffAccel,omitempty"`
	RearDiffDecel     *float64 `json:"rearDiffDecel,omitempty"`
	CenterDiffBalance *float64 `json:"centerDiffBalance,omitempty"`

	Notes string `json:"notes"`
}

type TuneProfileInput struct {
	CarName      string `json:"carName"`
	CarOrdinal   *int64 `json:"carOrdinal,omitempty"`
	CarCategory  *int64 `json:"carCategory,omitempty"`
	CarClass     string `json:"carClass"`
	PI           *int64 `json:"pi,omitempty"`
	Drivetrain   string `json:"drivetrain"`
	NumCylinders *int64 `json:"numCylinders,omitempty"`
	UseCase      string `json:"useCase"`
	VersionName  string `json:"versionName"`

	PowerKW              *float64 `json:"powerKW,omitempty"`
	TorqueNM             *float64 `json:"torqueNM,omitempty"`
	WeightKG             *float64 `json:"weightKG,omitempty"`
	FrontWeightPct       *float64 `json:"frontWeightPct,omitempty"`
	PowerToWeightKWPerKG *float64 `json:"powerToWeightKWPerKG,omitempty"`
	PeakTorqueRPM        *float64 `json:"peakTorqueRPM,omitempty"`
	PeakPowerRPM         *float64 `json:"peakPowerRPM,omitempty"`
	RedlineRPM           *float64 `json:"redlineRPM,omitempty"`

	FrontTirePressure *float64 `json:"frontTirePressure,omitempty"`
	RearTirePressure  *float64 `json:"rearTirePressure,omitempty"`

	FinalDrive *float64 `json:"finalDrive,omitempty"`
	Gear1      *float64 `json:"gear1,omitempty"`
	Gear2      *float64 `json:"gear2,omitempty"`
	Gear3      *float64 `json:"gear3,omitempty"`
	Gear4      *float64 `json:"gear4,omitempty"`
	Gear5      *float64 `json:"gear5,omitempty"`
	Gear6      *float64 `json:"gear6,omitempty"`
	Gear7      *float64 `json:"gear7,omitempty"`
	Gear8      *float64 `json:"gear8,omitempty"`
	Gear9      *float64 `json:"gear9,omitempty"`
	Gear10     *float64 `json:"gear10,omitempty"`

	FrontCamber *float64 `json:"frontCamber,omitempty"`
	RearCamber  *float64 `json:"rearCamber,omitempty"`
	FrontToe    *float64 `json:"frontToe,omitempty"`
	RearToe     *float64 `json:"rearToe,omitempty"`
	Caster      *float64 `json:"caster,omitempty"`

	FrontARB *float64 `json:"frontArb,omitempty"`
	RearARB  *float64 `json:"rearArb,omitempty"`

	FrontSpring     *float64 `json:"frontSpring,omitempty"`
	RearSpring      *float64 `json:"rearSpring,omitempty"`
	FrontRideHeight *float64 `json:"frontRideHeight,omitempty"`
	RearRideHeight  *float64 `json:"rearRideHeight,omitempty"`

	FrontRebound *float64 `json:"frontRebound,omitempty"`
	RearRebound  *float64 `json:"rearRebound,omitempty"`
	FrontBump    *float64 `json:"frontBump,omitempty"`
	RearBump     *float64 `json:"rearBump,omitempty"`

	FrontAero   *float64 `json:"frontAero,omitempty"`
	RearAero    *float64 `json:"rearAero,omitempty"`
	AeroBalance *float64 `json:"aeroBalance,omitempty"`

	BrakeBalance  *float64 `json:"brakeBalance,omitempty"`
	BrakePressure *float64 `json:"brakePressure,omitempty"`

	FrontDiffAccel    *float64 `json:"frontDiffAccel,omitempty"`
	FrontDiffDecel    *float64 `json:"frontDiffDecel,omitempty"`
	RearDiffAccel     *float64 `json:"rearDiffAccel,omitempty"`
	RearDiffDecel     *float64 `json:"rearDiffDecel,omitempty"`
	CenterDiffBalance *float64 `json:"centerDiffBalance,omitempty"`

	Notes string `json:"notes"`
}

type TelemetrySession struct {
	ID                     int64    `json:"id"`
	TuneProfileID          *int64   `json:"tuneProfileId,omitempty"`
	TuneSnapshotJSON       string   `json:"tuneSnapshotJson"`
	TuneName               string   `json:"tuneName"`
	SessionName            string   `json:"sessionName"`
	TrackName              string   `json:"trackName"`
	Mode                   string   `json:"mode"`
	GameMode               string   `json:"gameMode"`
	StartedAt              string   `json:"startedAt"`
	EndedAt                string   `json:"endedAt"`
	DurationMS             int64    `json:"durationMs"`
	BestLapMS              *int64   `json:"bestLapMs,omitempty"`
	AvgSpeedKmh            *float64 `json:"avgSpeedKmh,omitempty"`
	MaxSpeedKmh            *float64 `json:"maxSpeedKmh,omitempty"`
	EventCount             int64    `json:"eventCount"`
	SampleCount            int64    `json:"sampleCount"`
	RecordingPath          string   `json:"recordingPath"`
	RecordingPackets       int64    `json:"recordingPackets"`
	RecordingBytes         int64    `json:"recordingBytes"`
	RecordingTruncated     bool     `json:"recordingTruncated"`
	CarOrdinal             *int64   `json:"carOrdinal,omitempty"`
	CarClass               string   `json:"carClass"`
	CarPI                  *int64   `json:"carPi,omitempty"`
	Drivetrain             string   `json:"drivetrain"`
	NumCylinders           *int64   `json:"numCylinders,omitempty"`
	DriverMode             string   `json:"driverMode"`
	DriverModeConfidence   float64  `json:"driverModeConfidence"`
	DriverModeEvidenceJSON string   `json:"driverModeEvidenceJson"`
	BrakeAssist            string   `json:"brakeAssist"`
	SteeringAssist         string   `json:"steeringAssist"`
	TractionControl        string   `json:"tractionControl"`
	StabilityControl       string   `json:"stabilityControl"`
	Shifting               string   `json:"shifting"`
	LaunchControl          string   `json:"launchControl"`
	DriverFeedbackJSON     string   `json:"driverFeedbackJson"`
	Notes                  string   `json:"notes"`
}

type SessionVehicleSnapshot struct {
	CarOrdinal   *int64 `json:"carOrdinal,omitempty"`
	CarClass     string `json:"carClass"`
	CarPI        *int64 `json:"carPi,omitempty"`
	Drivetrain   string `json:"drivetrain"`
	NumCylinders *int64 `json:"numCylinders,omitempty"`
}

type TestConditions struct {
	DriverMode       string `json:"driverMode"`
	BrakeAssist      string `json:"brakeAssist"`
	SteeringAssist   string `json:"steeringAssist"`
	TractionControl  string `json:"tractionControl"`
	StabilityControl string `json:"stabilityControl"`
	Shifting         string `json:"shifting"`
	LaunchControl    string `json:"launchControl"`
}

type DriverModeDetection struct {
	Mode       string             `json:"mode"`
	Confidence float64            `json:"confidence"`
	Summary    string             `json:"summary"`
	Evidence   map[string]float64 `json:"evidence"`
}

type SessionStartInput struct {
	TuneProfileID    *int64         `json:"tuneProfileId,omitempty"`
	TuneSnapshotJSON string         `json:"tuneSnapshotJson"`
	SessionName      string         `json:"sessionName"`
	TrackName        string         `json:"trackName"`
	Mode             string         `json:"mode"`
	GameMode         string         `json:"gameMode"`
	StartedAt        string         `json:"startedAt"`
	RecordingPath    string         `json:"recordingPath"`
	TestConditions   TestConditions `json:"testConditions"`
	SessionVehicleSnapshot
}

type TuneProfileSnapshot struct {
	ID            int64       `json:"id"`
	TuneProfileID int64       `json:"tuneProfileId"`
	SessionID     *int64      `json:"sessionId,omitempty"`
	ChangedAt     string      `json:"changedAt"`
	ChangeReason  string      `json:"changeReason"`
	Before        TuneProfile `json:"before"`
	After         TuneProfile `json:"after"`
	ChangedFields []string    `json:"changedFields"`
	ChangeJSON    string      `json:"changeJson"`
}

type SessionFinalizeInput struct {
	SessionID           int64               `json:"sessionId"`
	EndedAt             string              `json:"endedAt"`
	DurationMS          int64               `json:"durationMs"`
	AvgSpeedKmh         *float64            `json:"avgSpeedKmh,omitempty"`
	MaxSpeedKmh         *float64            `json:"maxSpeedKmh,omitempty"`
	RecordingPackets    int64               `json:"recordingPackets"`
	RecordingBytes      int64               `json:"recordingBytes"`
	RecordingTruncated  bool                `json:"recordingTruncated"`
	GameMode            string              `json:"gameMode"`
	DriverModeDetection DriverModeDetection `json:"driverModeDetection"`
	Notes               string              `json:"notes"`
	SessionVehicleSnapshot
}

type UpgradeUnlockRule struct {
	Category    string `json:"category"`
	UpgradeName string `json:"upgradeName"`
	Unlocks     string `json:"unlocks"`
}

type TuneAdjustmentExplanation struct {
	Category    string `json:"category"`
	Item        string `json:"item"`
	Detail      string `json:"detail"`
	Description string `json:"description"`
}

type RuleThresholdProfile struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	CarClass   string `json:"carClass"`
	Drivetrain string `json:"drivetrain"`
	UseCase    string `json:"useCase"`
	GameMode   string `json:"gameMode"`
	ConfigJSON string `json:"configJson"`
	IsDefault  bool   `json:"isDefault"`
	CreatedAt  string `json:"createdAt"`
	UpdatedAt  string `json:"updatedAt"`
}

type RuleThresholdProfileInput struct {
	Name       string `json:"name"`
	CarClass   string `json:"carClass"`
	Drivetrain string `json:"drivetrain"`
	UseCase    string `json:"useCase"`
	GameMode   string `json:"gameMode"`
	ConfigJSON string `json:"configJson"`
}

type BenchmarkPoint struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

type BenchmarkGate struct {
	Center      BenchmarkPoint `json:"center"`
	DirectionX  float64        `json:"directionX"`
	DirectionZ  float64        `json:"directionZ"`
	WidthMeters float64        `json:"widthMeters"`
	DepthMeters float64        `json:"depthMeters"`
}

type BenchmarkTrackInput struct {
	Name              string           `json:"name"`
	SourceMode        string           `json:"sourceMode"`
	TrackType         string           `json:"trackType"`
	Start             BenchmarkPoint   `json:"start"`
	End               BenchmarkPoint   `json:"end"`
	StartRadius       float64          `json:"startRadius"`
	EndRadius         float64          `json:"endRadius"`
	DirectionX        float64          `json:"directionX"`
	DirectionZ        float64          `json:"directionZ"`
	StartGate         BenchmarkGate    `json:"startGate"`
	FinishGate        BenchmarkGate    `json:"finishGate"`
	Checkpoints       []BenchmarkPoint `json:"checkpoints"`
	RouteLengthMeters float64          `json:"routeLengthMeters"`
	HasDrivingLine    bool             `json:"hasDrivingLine"`
	Polyline          []BenchmarkPoint `json:"polyline"`
	SourceSessionID   *int64           `json:"sourceSessionId,omitempty"`
	LapCountObserved  int              `json:"lapCountObserved"`
	Notes             string           `json:"notes"`
}

type BenchmarkTrackExtractionInput struct {
	SessionID      int64          `json:"sessionId"`
	Name           string         `json:"name"`
	TrackType      string         `json:"trackType"`
	ExtractionMode string         `json:"extractionMode"`
	StartGate      *BenchmarkGate `json:"startGate,omitempty"`
	FinishGate     *BenchmarkGate `json:"finishGate,omitempty"`
}

type BenchmarkTrack struct {
	ID int64 `json:"id"`
	BenchmarkTrackInput
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type BenchmarkRun struct {
	ID                          int64    `json:"id"`
	SessionID                   int64    `json:"sessionId"`
	TrackID                     int64    `json:"trackId"`
	TrackName                   string   `json:"trackName"`
	StartMS                     int64    `json:"startMs"`
	EndMS                       int64    `json:"endMs"`
	DurationMS                  int64    `json:"durationMs"`
	Confidence                  float64  `json:"confidence"`
	AvgSpeedKmh                 *float64 `json:"avgSpeedKmh,omitempty"`
	MaxSpeedKmh                 *float64 `json:"maxSpeedKmh,omitempty"`
	RouteProgress01             *float64 `json:"routeProgress01,omitempty"`
	GeometryLengthMeters        *float64 `json:"geometryLengthMeters,omitempty"`
	TrackLengthErrorPct         *float64 `json:"trackLengthErrorPct,omitempty"`
	DistanceTraveledDeltaMeters *float64 `json:"distanceTraveledDeltaMeters,omitempty"`
	CurrentRaceTimeDeltaSeconds *float64 `json:"currentRaceTimeDeltaSeconds,omitempty"`
	AvgLateralErrorMeters       *float64 `json:"avgLateralErrorMeters,omitempty"`
	MaxLateralErrorMeters       *float64 `json:"maxLateralErrorMeters,omitempty"`
	WarningFlags                string   `json:"warningFlags"`
	EventCount                  int64    `json:"eventCount"`
	DriverMode                  string   `json:"driverMode"`
	DriverModeConfidence        float64  `json:"driverModeConfidence"`
	DriverModeEvidenceJSON      string   `json:"driverModeEvidenceJson"`
	Valid                       bool     `json:"valid"`
	CreatedAt                   string   `json:"createdAt"`
}

type TrackBaselineRun struct {
	ID                          int64           `json:"id"`
	TrackID                     int64           `json:"trackId"`
	Vehicle                     TrackVehicleKey `json:"vehicle"`
	StartMS                     int64           `json:"startMs"`
	EndMS                       int64           `json:"endMs"`
	DurationMS                  int64           `json:"durationMs"`
	Confidence                  float64         `json:"confidence"`
	AvgSpeedKmh                 *float64        `json:"avgSpeedKmh,omitempty"`
	MaxSpeedKmh                 *float64        `json:"maxSpeedKmh,omitempty"`
	RouteProgress01             *float64        `json:"routeProgress01,omitempty"`
	GeometryLengthMeters        *float64        `json:"geometryLengthMeters,omitempty"`
	TrackLengthErrorPct         *float64        `json:"trackLengthErrorPct,omitempty"`
	DistanceTraveledDeltaMeters *float64        `json:"distanceTraveledDeltaMeters,omitempty"`
	CurrentRaceTimeDeltaSeconds *float64        `json:"currentRaceTimeDeltaSeconds,omitempty"`
	AvgLateralErrorMeters       *float64        `json:"avgLateralErrorMeters,omitempty"`
	MaxLateralErrorMeters       *float64        `json:"maxLateralErrorMeters,omitempty"`
	WarningFlags                string          `json:"warningFlags"`
	EventCount                  int64           `json:"eventCount"`
	DriverMode                  string          `json:"driverMode"`
	DriverModeConfidence        float64         `json:"driverModeConfidence"`
	DriverModeEvidenceJSON      string          `json:"driverModeEvidenceJson"`
	Valid                       bool            `json:"valid"`
	GameMode                    string          `json:"gameMode"`
	CreatedAt                   string          `json:"createdAt"`
}

type TrackBaselineSaveResult struct {
	Track          BenchmarkTrack       `json:"track"`
	Baseline       TrackBaselineRun     `json:"baseline"`
	Action         string               `json:"action"`
	MatchCandidate *TrackMergeCandidate `json:"matchCandidate,omitempty"`
}

type TrackVehicleKey struct {
	CarOrdinal *int64 `json:"carOrdinal,omitempty"`
	CarClass   string `json:"carClass"`
	CarPI      *int64 `json:"carPi,omitempty"`
	Drivetrain string `json:"drivetrain"`
	Label      string `json:"label"`
}

type TrackRunContext struct {
	Run     BenchmarkRun     `json:"run"`
	Session TelemetrySession `json:"session"`
	Vehicle TrackVehicleKey  `json:"vehicle"`
}

type TrackAutoBaseline struct {
	Vehicle    TrackVehicleKey   `json:"vehicle"`
	BestRun    TrackRunContext   `json:"bestRun"`
	RecentRuns []TrackRunContext `json:"recentRuns"`
	RunCount   int               `json:"runCount"`
}

type TrackVehicleReference struct {
	Vehicle            TrackVehicleKey    `json:"vehicle"`
	BestAutoBaseline   *TrackRunContext   `json:"bestAutoBaseline,omitempty"`
	BestTrackBaseline  *TrackBaselineRun  `json:"bestTrackBaseline,omitempty"`
	RecentRuns         []TrackRunContext  `json:"recentRuns"`
	RecentBaselineRuns []TrackBaselineRun `json:"recentBaselineRuns"`
	ValidRunCount      int                `json:"validRunCount"`
	AutoRunCount       int                `json:"autoRunCount"`
	BaselineRunCount   int                `json:"baselineRunCount"`
	AvgSpeedKmh        *float64           `json:"avgSpeedKmh,omitempty"`
	MaxSpeedKmh        *float64           `json:"maxSpeedKmh,omitempty"`
	EventCount         int64              `json:"eventCount"`
}

type TrackProfile struct {
	Track             BenchmarkTrack          `json:"track"`
	AutoBaselines     []TrackAutoBaseline     `json:"autoBaselines"`
	VehicleReferences []TrackVehicleReference `json:"vehicleReferences"`
	RecentRuns        []TrackRunContext       `json:"recentRuns"`
	Warnings          []string                `json:"warnings"`
}

type TrackMergeCandidate struct {
	Track                  BenchmarkTrack `json:"track"`
	MatchLevel             string         `json:"matchLevel"`
	LengthErrorPct         float64        `json:"lengthErrorPct"`
	StartDistanceMeters    float64        `json:"startDistanceMeters"`
	EndDistanceMeters      float64        `json:"endDistanceMeters"`
	ShapeSimilarity        float64        `json:"shapeSimilarity"`
	RouteFitAvgErrorMeters float64        `json:"routeFitAvgErrorMeters"`
	RouteFitP90ErrorMeters float64        `json:"routeFitP90ErrorMeters"`
	RouteFitScore          float64        `json:"routeFitScore"`
	DirectionMatched       bool           `json:"directionMatched"`
	ReverseMatched         bool           `json:"reverseMatched"`
	Reason                 string         `json:"reason"`
}

type TuneProfileSessionStat struct {
	TuneProfileID int64  `json:"tuneProfileId"`
	SessionCount  int64  `json:"sessionCount"`
	LastStartedAt string `json:"lastStartedAt"`
}

type SessionComparison struct {
	LeftSession           TelemetrySession          `json:"leftSession"`
	RightSession          TelemetrySession          `json:"rightSession"`
	Metrics               []SessionComparisonMetric `json:"metrics"`
	EventTypes            []SessionEventComparison  `json:"eventTypes"`
	ComparabilityWarnings []string                  `json:"comparabilityWarnings"`
}

type SessionComparisonMetric struct {
	Key            string  `json:"key"`
	Label          string  `json:"label"`
	Unit           string  `json:"unit"`
	Left           float64 `json:"left"`
	Right          float64 `json:"right"`
	Delta          float64 `json:"delta"`
	HigherIsBetter bool    `json:"higherIsBetter"`
}

type SessionEventComparison struct {
	Type  string `json:"type"`
	Left  int    `json:"left"`
	Right int    `json:"right"`
	Delta int    `json:"delta"`
}

type RoadSessionEvaluation struct {
	Session               TelemetrySession            `json:"session"`
	Track                 *BenchmarkTrack             `json:"track,omitempty"`
	BestRun               *BenchmarkRun               `json:"bestRun,omitempty"`
	BaselineRun           *BenchmarkRun               `json:"baselineRun,omitempty"`
	BaselineSession       *TelemetrySession           `json:"baselineSession,omitempty"`
	BaselineStatus        string                      `json:"baselineStatus"`
	PaperPerformanceScore float64                     `json:"paperPerformanceScore"`
	PlayerFitScore        float64                     `json:"playerFitScore"`
	RiskScore             float64                     `json:"riskScore"`
	OverallVerdict        string                      `json:"overallVerdict"`
	Attributions          []RoadEvaluationAttribution `json:"attributions"`
	Notes                 []string                    `json:"notes"`
}

type RoadEvaluationAttribution struct {
	Type             string `json:"type"`
	EventType        string `json:"eventType,omitempty"`
	Count            int    `json:"count"`
	Severity         string `json:"severity,omitempty"`
	Priority         int    `json:"priority"`
	Message          string `json:"message"`
	PrioritizeTuning bool   `json:"prioritizeTuning"`
}

type RoadEvaluationComparison struct {
	Left    RoadSessionEvaluation     `json:"left"`
	Right   RoadSessionEvaluation     `json:"right"`
	Metrics []SessionComparisonMetric `json:"metrics"`
	Verdict string                    `json:"verdict"`
	Notes   []string                  `json:"notes"`
}
