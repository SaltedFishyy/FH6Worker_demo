package storage

import "fh6worker/internal/telemetry"

type SessionIssueSummary struct {
	SessionID          int64               `json:"sessionId"`
	BaselineSession    *TelemetrySession   `json:"baselineSession,omitempty"`
	BaselineStatus     string              `json:"baselineStatus"`
	RecentChangeFields []string            `json:"recentChangeFields"`
	Groups             []SessionIssueGroup `json:"groups"`
	GearPower          GearPowerDiagnostic `json:"gearPower"`
	WholeCarPlan       WholeCarTuningPlan  `json:"wholeCarPlan"`
}

type QuickDiagnostic struct {
	Status               string                 `json:"status"`
	ComparisonStatus     string                 `json:"comparisonStatus"`
	UpdatedAt            string                 `json:"updatedAt"`
	SampleCount          int                    `json:"sampleCount"`
	EventCount           int                    `json:"eventCount"`
	GameMode             string                 `json:"gameMode"`
	DriverMode           string                 `json:"driverMode"`
	DriverModeConfidence float64                `json:"driverModeConfidence"`
	Vehicle              SessionVehicleSnapshot `json:"vehicle"`
	Comparability        QuickComparability     `json:"comparability"`
	CurrentLap           *QuickLapSummary       `json:"currentLap,omitempty"`
	PreviousLap          *QuickLapSummary       `json:"previousLap,omitempty"`
	Groups               []SessionIssueGroup    `json:"groups"`
	GearPower            GearPowerDiagnostic    `json:"gearPower"`
	Suggestions          []QuickSuggestion      `json:"suggestions"`
	MissingProfileFields []string               `json:"missingProfileFields"`
}

type TireModelDiagnostic struct {
	Status        string                 `json:"status"`
	UpdatedAt     string                 `json:"updatedAt"`
	SampleCount   int                    `json:"sampleCount"`
	WindowMS      int64                  `json:"windowMs"`
	GameMode      string                 `json:"gameMode"`
	Phase         string                 `json:"phase"`
	PhaseDetail   TirePhaseDiagnostic    `json:"phaseDetail"`
	DataQuality   TireDataQuality        `json:"dataQuality"`
	GripLimit     TireGripLimit          `json:"gripLimit"`
	LimitType     string                 `json:"limitType"`
	Confidence    string                 `json:"confidence"`
	Summary       string                 `json:"summary"`
	Explanation   string                 `json:"explanation"`
	Warnings      []string               `json:"warnings"`
	Wheels        []TireWheelDiagnostic  `json:"wheels"`
	FrontAxle     TireAxleDiagnostic     `json:"frontAxle"`
	RearAxle      TireAxleDiagnostic     `json:"rearAxle"`
	LeftRight     TireSideBalance        `json:"leftRight"`
	GForce        GForceDiagnostic       `json:"gForce"`
	Camber        CamberInference        `json:"camber"`
	PowerToTire   PowerToTireDiagnostic  `json:"powerToTire"`
	BrakeToTire   BrakeToTireDiagnostic  `json:"brakeToTire"`
	IssueAnalysis TireIssueAnalysis      `json:"issueAnalysis"`
	IssueAdvice   TireIssueAdvice        `json:"issueAdvice"`
	Hints         []TireModelHint        `json:"hints"`
	Evidence      map[string]float64     `json:"evidence"`
	Vehicle       SessionVehicleSnapshot `json:"vehicle"`
}

type TirePhaseDiagnostic struct {
	CurrentPhase   string             `json:"currentPhase"`
	SecondaryPhase string             `json:"secondaryPhase"`
	StablePhase    string             `json:"stablePhase"`
	PhaseStability string             `json:"phaseStability"`
	ScoreMargin    float64            `json:"scoreMargin"`
	Confidence     string             `json:"confidence"`
	Scores         map[string]float64 `json:"scores"`
	Evidence       map[string]float64 `json:"evidence"`
	WindowMS       int64              `json:"windowMs"`
	SampleCount    int                `json:"sampleCount"`
}

type TireDataQuality struct {
	Status             string             `json:"status"`
	Confidence         string             `json:"confidence"`
	SampleCount        int                `json:"sampleCount"`
	DynamicSampleCount int                `json:"dynamicSampleCount"`
	SpeedSignal        string             `json:"speedSignal"`
	GForceSignal       string             `json:"gForceSignal"`
	SlipSignal         string             `json:"slipSignal"`
	InputSignal        string             `json:"inputSignal"`
	Reasons            []string           `json:"reasons"`
	Evidence           map[string]float64 `json:"evidence"`
}

type TireGripLimit struct {
	Type            string             `json:"type"`
	LimitedAxle     string             `json:"limitedAxle"`
	LimitedWheels   []string           `json:"limitedWheels"`
	PrimaryEvidence string             `json:"primaryEvidence"`
	Confidence      string             `json:"confidence"`
	Reason          string             `json:"reason"`
	FrontRearDelta  float64            `json:"frontRearDelta"`
	DrivenDelta     float64            `json:"drivenDelta"`
	LeftRightDelta  float64            `json:"leftRightDelta"`
	Evidence        map[string]float64 `json:"evidence"`
}

type TireIssueAnalysis struct {
	Status       string             `json:"status"`
	UpdatedAt    string             `json:"updatedAt"`
	WindowMS     int64              `json:"windowMs"`
	SampleCount  int                `json:"sampleCount"`
	SegmentCount int                `json:"segmentCount"`
	GroupCount   int                `json:"groupCount"`
	Segments     []TireIssueSegment `json:"segments"`
	Groups       []TireIssueGroup   `json:"groups"`
	Warnings     []string           `json:"warnings"`
}

type TireIssueAdvice struct {
	Status                string                  `json:"status"`
	UpdatedAt             string                  `json:"updatedAt"`
	Confidence            string                  `json:"confidence"`
	BasedOnIssueUpdatedAt string                  `json:"basedOnIssueUpdatedAt"`
	IssueGroupCount       int                     `json:"issueGroupCount"`
	PriorityActions       []TireIssueAdviceAction `json:"priorityActions"`
	Groups                []TireIssueAdviceGroup  `json:"groups"`
	Warnings              []string                `json:"warnings"`
}

type TireIssueAdviceGroup struct {
	IssueGroupID  string                  `json:"issueGroupId"`
	IssueType     string                  `json:"issueType"`
	Phase         string                  `json:"phase"`
	OperationTags []string                `json:"operationTags"`
	LimitedAxle   string                  `json:"limitedAxle"`
	DriftSource   string                  `json:"driftSource"`
	PrimaryCause  string                  `json:"primaryCause"`
	ShouldTune    bool                    `json:"shouldTune"`
	Priority      int                     `json:"priority"`
	Confidence    string                  `json:"confidence"`
	Evidence      map[string]float64      `json:"evidence"`
	Actions       []TireIssueAdviceAction `json:"actions"`
}

type TireIssueAdviceAction struct {
	ID              string   `json:"id"`
	IssueGroupID    string   `json:"issueGroupId"`
	Layer           string   `json:"layer"`
	Category        string   `json:"category"`
	Scope           string   `json:"scope"`
	Direction       string   `json:"direction"`
	RelatedFields   []string `json:"relatedFields"`
	Rationale       string   `json:"rationale"`
	VerifyEvidence  []string `json:"verifyEvidence"`
	Confidence      string   `json:"confidence"`
	MissingInputs   []string `json:"missingInputs"`
	ConflictReason  string   `json:"conflictReason"`
	TuneRecommended bool     `json:"tuneRecommended"`
}

type TireIssueSegment struct {
	ID            string             `json:"id"`
	Type          string             `json:"type"`
	Phase         string             `json:"phase"`
	OperationTags []string           `json:"operationTags"`
	DriftSource   string             `json:"driftSource"`
	LimitType     string             `json:"limitType"`
	LimitedAxle   string             `json:"limitedAxle"`
	LimitedWheels []string           `json:"limitedWheels"`
	StartMS       int64              `json:"startMs"`
	EndMS         int64              `json:"endMs"`
	DurationMS    int64              `json:"durationMs"`
	SampleCount   int                `json:"sampleCount"`
	SpeedMinKmh   float64            `json:"speedMinKmh"`
	SpeedMaxKmh   float64            `json:"speedMaxKmh"`
	SpeedAvgKmh   float64            `json:"speedAvgKmh"`
	Confidence    string             `json:"confidence"`
	DataQuality   string             `json:"dataQuality"`
	RiskLevel     string             `json:"riskLevel"`
	Evidence      map[string]float64 `json:"evidence"`
	Reason        string             `json:"reason"`
}

type TireIssueGroup struct {
	ID                     string             `json:"id"`
	Type                   string             `json:"type"`
	Phase                  string             `json:"phase"`
	OperationTags          []string           `json:"operationTags"`
	DriftSource            string             `json:"driftSource"`
	LimitType              string             `json:"limitType"`
	LimitedAxle            string             `json:"limitedAxle"`
	LimitedWheels          []string           `json:"limitedWheels"`
	Count                  int                `json:"count"`
	TotalDurationMS        int64              `json:"totalDurationMs"`
	SpeedMinKmh            float64            `json:"speedMinKmh"`
	SpeedMaxKmh            float64            `json:"speedMaxKmh"`
	SpeedAvgKmh            float64            `json:"speedAvgKmh"`
	Confidence             string             `json:"confidence"`
	DataQuality            string             `json:"dataQuality"`
	RiskLevel              string             `json:"riskLevel"`
	RepresentativeEvidence map[string]float64 `json:"representativeEvidence"`
	SegmentIDs             []string           `json:"segmentIds"`
	Reason                 string             `json:"reason"`
}

type GForceDiagnostic struct {
	Source        string        `json:"source"`
	AxisMapping   string        `json:"axisMapping"`
	CurrentXG     float64       `json:"currentXG"`
	CurrentYG     float64       `json:"currentYG"`
	CurrentZG     float64       `json:"currentZG"`
	CurrentTotalG float64       `json:"currentTotalG"`
	AvgAbsXG      float64       `json:"avgAbsXG"`
	AvgAbsYG      float64       `json:"avgAbsYG"`
	AvgAbsZG      float64       `json:"avgAbsZG"`
	AvgTotalG     float64       `json:"avgTotalG"`
	PeakAbsXG     float64       `json:"peakAbsXG"`
	PeakAbsYG     float64       `json:"peakAbsYG"`
	PeakAbsZG     float64       `json:"peakAbsZG"`
	PeakTotalG    float64       `json:"peakTotalG"`
	DominantAxis  string        `json:"dominantAxis"`
	Series        []GForcePoint `json:"series"`
}

type GForcePoint struct {
	TimeMS int64   `json:"timeMs"`
	XG     float64 `json:"xG"`
	YG     float64 `json:"yG"`
	ZG     float64 `json:"zG"`
	TotalG float64 `json:"totalG"`
}

type CamberInference struct {
	Status      string             `json:"status"`
	Confidence  string             `json:"confidence"`
	FrontState  string             `json:"frontState"`
	RearState   string             `json:"rearState"`
	Summary     string             `json:"summary"`
	Explanation string             `json:"explanation"`
	Warnings    []string           `json:"warnings"`
	Hints       []TireModelHint    `json:"hints"`
	Evidence    map[string]float64 `json:"evidence"`
}

type PowerToTireDiagnostic struct {
	Status                  string             `json:"status"`
	Summary                 string             `json:"summary"`
	Explanation             string             `json:"explanation"`
	Confidence              string             `json:"confidence"`
	SampleCount             int                `json:"sampleCount"`
	HighThrottleSampleCount int                `json:"highThrottleSampleCount"`
	Drivetrain              string             `json:"drivetrain"`
	DrivenAxle              string             `json:"drivenAxle"`
	CurrentPowerKW          float64            `json:"currentPowerKW"`
	AveragePowerKW          float64            `json:"averagePowerKW"`
	MaxPowerKW              float64            `json:"maxPowerKW"`
	CurrentTorqueNM         float64            `json:"currentTorqueNM"`
	AverageTorqueNM         float64            `json:"averageTorqueNM"`
	MaxTorqueNM             float64            `json:"maxTorqueNM"`
	CurrentRPM              float64            `json:"currentRPM"`
	AverageRPM              float64            `json:"averageRPM"`
	CurrentRPMRatio         float64            `json:"currentRPMRatio"`
	AverageRPMRatio         float64            `json:"averageRPMRatio"`
	CurrentGear             int                `json:"currentGear"`
	AverageThrottle         float64            `json:"averageThrottle"`
	AverageSpeedKmh         float64            `json:"averageSpeedKmh"`
	SpeedDeltaKmh           float64            `json:"speedDeltaKmh"`
	AverageAccelMps2        float64            `json:"averageAccelMps2"`
	AverageAccelG           float64            `json:"averageAccelG"`
	PeakAccelG              float64            `json:"peakAccelG"`
	FrontSlipRatioP90       float64            `json:"frontSlipRatioP90"`
	RearSlipRatioP90        float64            `json:"rearSlipRatioP90"`
	DrivenSlipRatioP90      float64            `json:"drivenSlipRatioP90"`
	DrivenSlipRatioHighPct  float64            `json:"drivenSlipRatioHighPct"`
	RPMLowHighThrottlePct   float64            `json:"rpmLowHighThrottlePct"`
	RPMHighHighThrottlePct  float64            `json:"rpmHighHighThrottlePct"`
	PowerSignalAvailable    bool               `json:"powerSignalAvailable"`
	TractionLimited         bool               `json:"tractionLimited"`
	Evidence                map[string]float64 `json:"evidence"`
}

type BrakeToTireDiagnostic struct {
	Status               string             `json:"status"`
	Summary              string             `json:"summary"`
	Explanation          string             `json:"explanation"`
	Confidence           string             `json:"confidence"`
	SampleCount          int                `json:"sampleCount"`
	BrakeSampleCount     int                `json:"brakeSampleCount"`
	AverageBrake         float64            `json:"averageBrake"`
	PeakBrake            float64            `json:"peakBrake"`
	AverageHandBrake     float64            `json:"averageHandBrake"`
	PeakHandBrake        float64            `json:"peakHandBrake"`
	AverageSpeedKmh      float64            `json:"averageSpeedKmh"`
	SpeedDeltaKmh        float64            `json:"speedDeltaKmh"`
	AverageSteer         float64            `json:"averageSteer"`
	AverageDecelMps2     float64            `json:"averageDecelMps2"`
	AverageDecelG        float64            `json:"averageDecelG"`
	PeakDecelG           float64            `json:"peakDecelG"`
	AveragePlaneG        float64            `json:"averagePlaneG"`
	PeakPlaneG           float64            `json:"peakPlaneG"`
	FrontSlipRatioP90    float64            `json:"frontSlipRatioP90"`
	RearSlipRatioP90     float64            `json:"rearSlipRatioP90"`
	FrontCombinedSlipP90 float64            `json:"frontCombinedSlipP90"`
	RearCombinedSlipP90  float64            `json:"rearCombinedSlipP90"`
	FrontRearSlipDelta   float64            `json:"frontRearSlipDelta"`
	TrailBraking         bool               `json:"trailBraking"`
	HandbrakeActive      bool               `json:"handbrakeActive"`
	Evidence             map[string]float64 `json:"evidence"`
}

type TireWheelDiagnostic struct {
	Position               string  `json:"position"`
	CombinedSlipAvg        float64 `json:"combinedSlipAvg"`
	CombinedSlipMax        float64 `json:"combinedSlipMax"`
	CombinedSlipP90        float64 `json:"combinedSlipP90"`
	CombinedSlipHighPct    float64 `json:"combinedSlipHighPct"`
	SlipRatioAvg           float64 `json:"slipRatioAvg"`
	SlipRatioMax           float64 `json:"slipRatioMax"`
	SlipRatioP90           float64 `json:"slipRatioP90"`
	SlipRatioHighPct       float64 `json:"slipRatioHighPct"`
	SlipAngleAvg           float64 `json:"slipAngleAvg"`
	SlipAngleMax           float64 `json:"slipAngleMax"`
	SlipAngleP90           float64 `json:"slipAngleP90"`
	TireTempAvg            float64 `json:"tireTempAvg"`
	TireTempMax            float64 `json:"tireTempMax"`
	SuspensionTravelAvg    float64 `json:"suspensionTravelAvg"`
	SuspensionTravelMax    float64 `json:"suspensionTravelMax"`
	SuspensionOffsetPctAvg float64 `json:"suspensionOffsetPctAvg"`
	SuspensionOffsetPctMax float64 `json:"suspensionOffsetPctMax"`
	SuspensionTravelMAvg   float64 `json:"suspensionTravelMetersAvg"`
	SuspensionTravelMMax   float64 `json:"suspensionTravelMetersMax"`
	GripState              string  `json:"gripState"`
}

type TireAxleDiagnostic struct {
	Name                   string  `json:"name"`
	CombinedSlipAvg        float64 `json:"combinedSlipAvg"`
	CombinedSlipMax        float64 `json:"combinedSlipMax"`
	CombinedSlipP90        float64 `json:"combinedSlipP90"`
	CombinedSlipHighPct    float64 `json:"combinedSlipHighPct"`
	SlipRatioAvg           float64 `json:"slipRatioAvg"`
	SlipRatioMax           float64 `json:"slipRatioMax"`
	SlipRatioP90           float64 `json:"slipRatioP90"`
	SlipRatioHighPct       float64 `json:"slipRatioHighPct"`
	SlipAngleAvg           float64 `json:"slipAngleAvg"`
	SlipAngleMax           float64 `json:"slipAngleMax"`
	SlipAngleP90           float64 `json:"slipAngleP90"`
	TireTempAvg            float64 `json:"tireTempAvg"`
	TireTempMax            float64 `json:"tireTempMax"`
	SuspensionTravelAvg    float64 `json:"suspensionTravelAvg"`
	SuspensionTravelMax    float64 `json:"suspensionTravelMax"`
	SuspensionOffsetPctAvg float64 `json:"suspensionOffsetPctAvg"`
	SuspensionOffsetPctMax float64 `json:"suspensionOffsetPctMax"`
	LimitScore             float64 `json:"limitScore"`
	GripState              string  `json:"gripState"`
}

type TireSideBalance struct {
	LeftCombinedSlipAvg  float64 `json:"leftCombinedSlipAvg"`
	RightCombinedSlipAvg float64 `json:"rightCombinedSlipAvg"`
	Delta                float64 `json:"delta"`
	State                string  `json:"state"`
}

type TireModelHint struct {
	Code      string `json:"code"`
	Severity  string `json:"severity"`
	Direction string `json:"direction"`
	Reason    string `json:"reason"`
}

type QuickComparability struct {
	SameVehicleClass string                 `json:"sameVehicleClass"`
	SameTrackContext string                 `json:"sameTrackContext"`
	Confidence       string                 `json:"confidence"`
	Warnings         []string               `json:"warnings"`
	BaselineVehicle  SessionVehicleSnapshot `json:"baselineVehicle"`
	CurrentVehicle   SessionVehicleSnapshot `json:"currentVehicle"`
}

type QuickLapSummary struct {
	LapNumber   int     `json:"lapNumber"`
	SampleCount int     `json:"sampleCount"`
	DurationMS  int64   `json:"durationMs"`
	AvgSpeedKmh float64 `json:"avgSpeedKmh"`
	MaxSpeedKmh float64 `json:"maxSpeedKmh"`
	EventCount  int     `json:"eventCount"`
	IssueScore  float64 `json:"issueScore"`
}

type QuickSuggestion struct {
	Family        string   `json:"family"`
	Source        string   `json:"source"`
	Confidence    string   `json:"confidence"`
	TrustLevel    string   `json:"trustLevel"`
	AdviceLayer   string   `json:"adviceLayer"`
	Category      string   `json:"category"`
	Item          string   `json:"item"`
	Direction     string   `json:"direction"`
	Amount        string   `json:"amount"`
	Reason        string   `json:"reason"`
	Rationale     string   `json:"rationale"`
	NextStep      string   `json:"nextStep"`
	FieldKeys     []string `json:"fieldKeys"`
	MissingInputs []string `json:"missingInputs"`
	CanApply      bool     `json:"canApply"`
	BlockedReason string   `json:"blockedReason"`
}

type SessionIssueGroup struct {
	ID                      string                      `json:"id"`
	Family                  string                      `json:"family"`
	Severity                string                      `json:"severity"`
	Segment                 string                      `json:"segment"`
	EventTypes              []string                    `json:"eventTypes"`
	EventIDs                []string                    `json:"eventIds"`
	Events                  []telemetry.DetectedEvent   `json:"events"`
	EventCount              int                         `json:"eventCount"`
	TotalDurationMS         int64                       `json:"totalDurationMs"`
	FirstStartMS            int64                       `json:"firstStartMs"`
	LastEndMS               int64                       `json:"lastEndMs"`
	Evidence                map[string]IssueEvidence    `json:"evidence"`
	PrimaryActions          []telemetry.SuggestedAction `json:"primaryActions"`
	Comparison              string                      `json:"comparison"`
	BaselineEventCount      int                         `json:"baselineEventCount"`
	BaselineTotalDurationMS int64                       `json:"baselineTotalDurationMs"`
	RelatedRecentChanges    []string                    `json:"relatedRecentChanges"`
	PrioritizeTuning        bool                        `json:"prioritizeTuning"`
	AdjustmentStrategy      string                      `json:"adjustmentStrategy"`
	FeedbackDirective       string                      `json:"feedbackDirective"`
}

type IssueEvidence struct {
	Min   float64 `json:"min"`
	Max   float64 `json:"max"`
	Avg   float64 `json:"avg"`
	Count int     `json:"count"`
}

type GearPowerDiagnostic struct {
	Status                 string                      `json:"status"`
	Summary                string                      `json:"summary"`
	LaunchFinding          string                      `json:"launchFinding"`
	TopSpeedFinding        string                      `json:"topSpeedFinding"`
	PowerKW                float64                     `json:"powerKW"`
	TorqueNM               float64                     `json:"torqueNM"`
	WeightKG               float64                     `json:"weightKG"`
	FrontWeightPct         float64                     `json:"frontWeightPct"`
	PowerToWeightKWPerKG   float64                     `json:"powerToWeightKWPerKG"`
	PowerToWeightBand      string                      `json:"powerToWeightBand"`
	PeakTorqueRPM          float64                     `json:"peakTorqueRPM"`
	PeakPowerRPM           float64                     `json:"peakPowerRPM"`
	RedlineRPM             float64                     `json:"redlineRPM"`
	PowerBandStartRPM      float64                     `json:"powerBandStartRPM"`
	PowerBandEndRPM        float64                     `json:"powerBandEndRPM"`
	PowerBandSource        string                      `json:"powerBandSource"`
	Confidence             string                      `json:"confidence"`
	StrategyMode           string                      `json:"strategyMode"`
	GlobalGearIssueCount   int                         `json:"globalGearIssueCount"`
	UsableGearCount        int                         `json:"usableGearCount"`
	GlobalGearIssueRatio   float64                     `json:"globalGearIssueRatio"`
	TractionLimitedPercent float64                     `json:"tractionLimitedPercent"`
	LowRpmHighLoadPercent  float64                     `json:"lowRpmHighLoadPercent"`
	HighRpmHighLoadPercent float64                     `json:"highRpmHighLoadPercent"`
	Gears                  []GearPowerBand             `json:"gears"`
	Comparisons            []GearPowerComparison       `json:"comparisons"`
	RecommendedActions     []telemetry.SuggestedAction `json:"recommendedActions"`
	Evidence               map[string]float64          `json:"evidence"`
}

type GearPowerBand struct {
	Gear                    int     `json:"gear"`
	SampleCount             int     `json:"sampleCount"`
	HighLoadSampleCount     int     `json:"highLoadSampleCount"`
	SpeedMinKmh             float64 `json:"speedMinKmh"`
	SpeedMaxKmh             float64 `json:"speedMaxKmh"`
	SpeedAvgKmh             float64 `json:"speedAvgKmh"`
	RpmMin                  float64 `json:"rpmMin"`
	RpmMax                  float64 `json:"rpmMax"`
	RpmAvg                  float64 `json:"rpmAvg"`
	RpmRatioMin             float64 `json:"rpmRatioMin"`
	RpmRatioMax             float64 `json:"rpmRatioMax"`
	RpmRatioAvg             float64 `json:"rpmRatioAvg"`
	InPowerBandRpmMin       float64 `json:"inPowerBandRpmMin"`
	InPowerBandRpmMax       float64 `json:"inPowerBandRpmMax"`
	InPowerBandRatioMin     float64 `json:"inPowerBandRatioMin"`
	InPowerBandRatioMax     float64 `json:"inPowerBandRatioMax"`
	ThrottleAvg             float64 `json:"throttleAvg"`
	AccelAvgMps2            float64 `json:"accelAvgMps2"`
	AccelMaxMps2            float64 `json:"accelMaxMps2"`
	SpeedPer1000RpmKmh      float64 `json:"speedPer1000RpmKmh"`
	ShiftAfterRPM           float64 `json:"shiftAfterRPM"`
	ShiftDropRPM            float64 `json:"shiftDropRPM"`
	FrontSlipAvg            float64 `json:"frontSlipAvg"`
	RearSlipAvg             float64 `json:"rearSlipAvg"`
	FrontTractionLimitedPct float64 `json:"frontTractionLimitedPct"`
	RearTractionLimitedPct  float64 `json:"rearTractionLimitedPct"`
	BelowPowerBandPercent   float64 `json:"belowPowerBandPercent"`
	InPowerBandPercent      float64 `json:"inPowerBandPercent"`
	AbovePowerBandPercent   float64 `json:"abovePowerBandPercent"`
	LowRpmHighLoadPercent   float64 `json:"lowRpmHighLoadPercent"`
	HighRpmHighLoadPercent  float64 `json:"highRpmHighLoadPercent"`
	TractionLimitedPercent  float64 `json:"tractionLimitedPercent"`
	Finding                 string  `json:"finding"`
}

type GearPowerComparison struct {
	Type              string                   `json:"type"`
	Status            string                   `json:"status"`
	BaselineSessionID int64                    `json:"baselineSessionId,omitempty"`
	Rows              []GearPowerComparisonRow `json:"rows"`
}

type GearPowerComparisonRow struct {
	Item                   string  `json:"item"`
	Gear                   int     `json:"gear,omitempty"`
	BeforeValue            float64 `json:"beforeValue,omitempty"`
	AfterValue             float64 `json:"afterValue,omitempty"`
	DeltaValue             float64 `json:"deltaValue,omitempty"`
	BeforeSpeedMaxKmh      float64 `json:"beforeSpeedMaxKmh,omitempty"`
	AfterSpeedMaxKmh       float64 `json:"afterSpeedMaxKmh,omitempty"`
	SpeedMaxDeltaKmh       float64 `json:"speedMaxDeltaKmh,omitempty"`
	BeforeInPowerBandPct   float64 `json:"beforeInPowerBandPct,omitempty"`
	AfterInPowerBandPct    float64 `json:"afterInPowerBandPct,omitempty"`
	InPowerBandDeltaPct    float64 `json:"inPowerBandDeltaPct,omitempty"`
	BeforeTractionLimitPct float64 `json:"beforeTractionLimitPct,omitempty"`
	AfterTractionLimitPct  float64 `json:"afterTractionLimitPct,omitempty"`
	TractionLimitDeltaPct  float64 `json:"tractionLimitDeltaPct,omitempty"`
	BeforeFinding          string  `json:"beforeFinding,omitempty"`
	AfterFinding           string  `json:"afterFinding,omitempty"`
}

type WholeCarTuningPlan struct {
	Strategy   string               `json:"strategy"`
	Confidence string               `json:"confidence"`
	Summary    string               `json:"summary"`
	Actions    []WholeCarAdjustment `json:"actions"`
	Conflicts  []TuningConflict     `json:"conflicts"`
	Notes      []string             `json:"notes"`
}

type WholeCarAdjustment struct {
	Priority   int                `json:"priority"`
	Family     string             `json:"family"`
	Source     string             `json:"source"`
	Confidence string             `json:"confidence"`
	Category   string             `json:"category"`
	Item       string             `json:"item"`
	Direction  string             `json:"direction"`
	Amount     string             `json:"amount"`
	Reason     string             `json:"reason"`
	Evidence   map[string]float64 `json:"evidence"`
}

type TuningConflict struct {
	Key         string `json:"key"`
	KeptItem    string `json:"keptItem"`
	DroppedItem string `json:"droppedItem"`
	Reason      string `json:"reason"`
}

type TunePlanDraft struct {
	SessionID     int64                 `json:"sessionId"`
	TuneProfileID *int64                `json:"tuneProfileId,omitempty"`
	Status        string                `json:"status"`
	Summary       string                `json:"summary"`
	Actions       []TunePlanDraftAction `json:"actions"`
	Conflicts     []TuningConflict      `json:"conflicts"`
}

type TunePlanDraftAction struct {
	ID             string   `json:"id"`
	Family         string   `json:"family"`
	Source         string   `json:"source"`
	Confidence     string   `json:"confidence"`
	AdviceLayer    string   `json:"adviceLayer"`
	TrustLevel     string   `json:"trustLevel"`
	TrustReasons   []string `json:"trustReasons"`
	MissingInputs  []string `json:"missingInputs"`
	RetestGuard    string   `json:"retestGuard"`
	Rationale      string   `json:"rationale"`
	ConflictReason string   `json:"conflictReason"`
	Category       string   `json:"category"`
	Item           string   `json:"item"`
	FieldKey       string   `json:"fieldKey"`
	Direction      string   `json:"direction"`
	Reason         string   `json:"reason"`
	CurrentValue   *float64 `json:"currentValue,omitempty"`
	TargetValue    *float64 `json:"targetValue,omitempty"`
	Delta          *float64 `json:"delta,omitempty"`
	Unit           string   `json:"unit"`
	Step           float64  `json:"step"`
	CanApply       bool     `json:"canApply"`
	BlockedReason  string   `json:"blockedReason"`
}

type TunePlanApplyInput struct {
	SessionID         int64    `json:"sessionId"`
	SelectedActionIDs []string `json:"selectedActionIds"`
}

type TunePlanApplyResult struct {
	Profile        TuneProfile           `json:"profile"`
	AppliedActions []TunePlanDraftAction `json:"appliedActions"`
	ChangedFields  []string              `json:"changedFields"`
}

type RetestEvaluation struct {
	SessionID             int64                 `json:"sessionId"`
	BaselineSession       *TelemetrySession     `json:"baselineSession,omitempty"`
	Status                string                `json:"status"`
	Summary               string                `json:"summary"`
	Confidence            string                `json:"confidence"`
	BaselineReason        string                `json:"baselineReason"`
	ChangedFields         []string              `json:"changedFields"`
	ChangeSourceSessionID *int64                `json:"changeSourceSessionId,omitempty"`
	RollbackActions       []TunePlanDraftAction `json:"rollbackActions"`
	MetricSummary         []string              `json:"metricSummary"`
	Metrics               []RetestMetric        `json:"metrics"`
}

type RetestMetric struct {
	Key       string  `json:"key"`
	Current   float64 `json:"current"`
	Baseline  float64 `json:"baseline"`
	Delta     float64 `json:"delta"`
	Direction string  `json:"direction"`
	Status    string  `json:"status"`
}

type StrategyTemplate struct {
	ID                int64  `json:"id"`
	Name              string `json:"name"`
	CarClass          string `json:"carClass"`
	Drivetrain        string `json:"drivetrain"`
	UseCase           string `json:"useCase"`
	GameMode          string `json:"gameMode"`
	IsDefault         bool   `json:"isDefault"`
	EnabledEventCount int    `json:"enabledEventCount"`
	TotalEventCount   int    `json:"totalEventCount"`
	UpdatedAt         string `json:"updatedAt"`
}

type RoadStrategyAnalysis struct {
	Template          StrategyTemplate            `json:"template"`
	SessionIDs        []int64                     `json:"sessionIds"`
	SessionCount      int                         `json:"sessionCount"`
	TotalEvents       int                         `json:"totalEvents"`
	EventDistribution []StrategyEventDistribution `json:"eventDistribution"`
	IssueGroups       []StrategyIssueAggregate    `json:"issueGroups"`
	Hints             []StrategyAnalysisHint      `json:"hints"`
}

type StrategyEventDistribution struct {
	Type     string `json:"type"`
	Count    int    `json:"count"`
	Severity string `json:"severity"`
}

type StrategyIssueAggregate struct {
	Family         string `json:"family"`
	EventCount     int    `json:"eventCount"`
	SessionCount   int    `json:"sessionCount"`
	Severity       string `json:"severity"`
	Recommendation string `json:"recommendation"`
}

type StrategyAnalysisHint struct {
	Level     string `json:"level"`
	Message   string `json:"message"`
	EventType string `json:"eventType,omitempty"`
	Family    string `json:"family,omitempty"`
}

type RoadTuningKnowledgeStatus struct {
	LoadedAt      string `json:"loadedAt"`
	SourcePath    string `json:"sourcePath"`
	LastError     string `json:"lastError"`
	SymptomCount  int    `json:"symptomCount"`
	ActionCount   int    `json:"actionCount"`
	UsingFallback bool   `json:"usingFallback"`
}

type RoadSymptomCard struct {
	ID             string   `json:"id"`
	Phase          string   `json:"phase"`
	Symptom        string   `json:"symptom"`
	PrimaryCause   string   `json:"primaryCause"`
	EventTypes     []string `json:"eventTypes"`
	EvidenceRule   string   `json:"evidenceRule"`
	Priority       int      `json:"priority"`
	Source         string   `json:"source"`
	Enabled        bool     `json:"enabled"`
	SourcePriority int      `json:"sourcePriority"`
}

type RoadActionCard struct {
	SymptomID             string  `json:"symptomId"`
	Rank                  int     `json:"rank"`
	Category              string  `json:"category"`
	Item                  string  `json:"item"`
	FieldKey              string  `json:"fieldKey"`
	Direction             string  `json:"direction"`
	Amount                float64 `json:"amount"`
	Unit                  string  `json:"unit"`
	DrivetrainRestriction string  `json:"drivetrainRestriction"`
	CanAutoApply          bool    `json:"canAutoApply"`
	Description           string  `json:"description"`
	Source                string  `json:"source"`
	SourcePriority        int     `json:"sourcePriority"`
}

type RoadTuningDecision struct {
	SessionID           int64                      `json:"sessionId"`
	Status              string                     `json:"status"`
	SymptomID           string                     `json:"symptomId"`
	Phase               string                     `json:"phase"`
	Symptom             string                     `json:"symptom"`
	PrimaryCause        string                     `json:"primaryCause"`
	Confidence          string                     `json:"confidence"`
	FitVerdict          string                     `json:"fitVerdict"`
	Reason              string                     `json:"reason"`
	RollbackRecommended bool                       `json:"rollbackRecommended"`
	RelatedIssueGroup   *SessionIssueGroup         `json:"relatedIssueGroup,omitempty"`
	Evidence            map[string]float64         `json:"evidence"`
	Actions             []RoadTuningDecisionAction `json:"actions"`
	RetestFocus         []string                   `json:"retestFocus"`
	KnowledgeStatus     RoadTuningKnowledgeStatus  `json:"knowledgeStatus"`
}

type RoadTuningDecisionAction struct {
	ID             string             `json:"id"`
	Role           string             `json:"role"`
	Family         string             `json:"family"`
	Source         string             `json:"source"`
	Confidence     string             `json:"confidence"`
	TrustLevel     string             `json:"trustLevel"`
	AdviceLayer    string             `json:"adviceLayer"`
	Category       string             `json:"category"`
	Item           string             `json:"item"`
	FieldKey       string             `json:"fieldKey"`
	Direction      string             `json:"direction"`
	Amount         string             `json:"amount"`
	Unit           string             `json:"unit"`
	Reason         string             `json:"reason"`
	Rationale      string             `json:"rationale"`
	ConflictReason string             `json:"conflictReason"`
	CanAutoApply   bool               `json:"canAutoApply"`
	BlockedReason  string             `json:"blockedReason"`
	Evidence       map[string]float64 `json:"evidence"`
}
