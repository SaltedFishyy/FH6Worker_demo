import {useEffect, useMemo, useRef, useState, type ReactNode, type RefObject} from 'react';
import {Activity, AlertTriangle, Copy as CopyIcon, FileText, Gauge, Lock, Pencil, Plus, Power, Radio, Save, Square, Thermometer, Trash2, Waves, Zap} from 'lucide-react';
import './App.css';
import appIcon from './assets/images/icon.png';
import {
    AnalyzeRoadStrategySessions,
    AnalyzeSessionBenchmarkRuns,
    ApplyRoadStaticTuneBaseline,
    ApplyTunePlanDraft,
    BindTelemetrySessionTuneProfile,
    ClearTuneHarvestCandidates,
    CreateBenchmarkTrack,
    CreateTuneProfile,
    CompareTelemetrySessions,
    EvaluateRoadSession,
    DeleteAllRecommendedCars,
    DeleteBenchmarkTrack,
    DeleteRecommendedCar,
    DeleteTrackBaselineRun,
    DeleteTireRegressionSample,
    CreateRuleThresholdProfile,
    DeleteTelemetrySession,
    DeleteRuleThresholdProfile,
    DeleteTuneProfile,
    DuplicateTuneProfile,
    ExplainTuneFieldInfluence,
    ExtractBenchmarkTrackFromSession,
    FindSimilarBenchmarkTracks,
    GenerateRoadStaticTuneBaseline,
    GenerateTuningReport,
    GetActiveTuneProfile,
    GetCurrentTelemetry,
    GetNetworkInterfaces,
    GetProfessionalPipelineConfig,
    GetProfessionalTuningDiagnostic,
    GetQuickDiagnostic,
    GetRecentTelemetry,
    GetRoadTuningDecision,
    GetRoadTuningKnowledgeStatus,
    GetSessionEvents,
    GetSessionIssueSummary,
    GetSessionTelemetrySamples,
    GetTireModelDiagnostic,
    GetTuneToTireInfluenceMap,
    ListTuningModelPipelines,
    GetRetestEvaluation,
    GetTestConditionDefaults,
    GetTelemetryReplayStatus,
    GetTelemetryStatus,
    GetTuneWebServerStatus,
    GetTunePlanDraft,
    GetTireRegressionSample,
    GetTrackProfile,
    ListBenchmarkTracks,
    ListBenchmarkRuns,
    ListRuleThresholdProfiles,
    ListRecommendedCars,
    ListStrategyTemplates,
    ListTelemetrySessions,
    ListTuneAdjustmentExplanations,
    ListTuneProfilesForVehicle,
    ListTuneProfileSnapshots,
    ListTuneProfileSessionStats,
    ListTuneProfiles,
    ListTrackProfiles,
    ListTireRegressionSamples,
    ListUpgradeUnlockRules,
    LoadRecommendedCarsFileSelection,
    MergeBenchmarkTrackInput,
    PauseTelemetryReplay,
    ReplayTelemetrySession,
    ReextractBenchmarkTrack,
    ReloadTuningKnowledge,
    RenameBenchmarkTrack,
    ResetRuleThresholdProfile,
    ResolveCarNameByOrdinal,
    ResumeTelemetryReplay,
    RestoreTuneProfileSnapshot,
    RunAllTireRegressionSamples,
    RunTuneHarvest,
    RunTuningModelPipeline,
    RunTireRegressionSample,
    SearchTuneHarvestCandidates,
    SaveTestConditionDefaults,
    SaveProfessionalPipelineConfig,
    SaveRecommendedCar,
    SaveRecommendedCarRecord,
    SaveRecommendedCarsFile,
    SaveTireRegressionSample,
    SeekTelemetryReplay,
    SetActiveTuneProfile,
    StartQuickTelemetry,
    StartProfessionalTelemetry,
    StartTireModelTelemetry,
    StartTuneWebServer,
    StartTrackBaselineTelemetry,
    StartTrackCaptureTelemetry,
    StartTelemetryWithConditions,
    StopTelemetry,
    StopTrackBaselineTelemetry,
    StopTelemetryReplay,
    StopTuneHarvest,
    StopTuneWebServer,
    SaveTrackBaselineCaptureAuto,
    UpdateTireRegressionSampleExpectation,
    UpdateRuleThresholdProfile,
    UpdateTuneHarvestCandidateStatus,
    UpdateTuneProfile,
} from "../wailsjs/go/main/App";

type WheelTelemetry = {
    slipRatio: number;
    slipAngle: number;
    combinedSlip: number;
    tireTemp: number;
    suspensionTravel: number;
    suspensionTravelMeters: number;
    wheelRotationSpeed: number;
    rumbleStrip: number;
    puddleDepth: number;
    surfaceRumble: number;
};

type TelemetryFrame = {
    receivedAt: string;
    timeMs: number;
    isRaceOn: boolean;
    gameMode: string;
    speedKmh: number;
    speedFieldKmh: number;
    velocitySpeedKmh: number;
    speedSource: string;
    rpm: number;
    rpmRatio: number;
    engineMaxRpm: number;
    engineIdleRpm: number;
    gear: number;
    accelerationX: number;
    accelerationY: number;
    accelerationZ: number;
    velocityX: number;
    velocityY: number;
    velocityZ: number;
    yaw: number;
    pitch: number;
    roll: number;
    power: number;
    torque: number;
    positionX: number;
    positionY: number;
    positionZ: number;
    boost: number;
    fuel: number;
    distanceTraveled: number;
    bestLap: number;
    lastLap: number;
    currentLap: number;
    currentRaceTime: number;
    lapNumber: number;
    racePosition: number;
    smashableVelDiff: number;
    smashableMass: number;
    carOrdinal: number;
    carClassId: number;
    carClass: string;
    carPi: number;
    drivetrainType: number;
    drivetrain: string;
    numCylinders: number;
    carCategory: number;
    carCategoryName: string;
    throttle01: number;
    brake01: number;
    clutch01: number;
    handBrake01: number;
    steer01: number;
    drivingLine01: number;
    aiBrakeDifference01: number;
    frontSlipRatioAvg: number;
    rearSlipRatioAvg: number;
    frontSlipAngleAvg: number;
    rearSlipAngleAvg: number;
    frontCombinedSlipAvg: number;
    rearCombinedSlipAvg: number;
    tireTempFrontAvg: number;
    tireTempRearAvg: number;
    suspensionFrontAvg: number;
    suspensionRearAvg: number;
    yawRate: number;
    pitchRate: number;
    rollRate: number;
    wheelFL: WheelTelemetry;
    wheelFR: WheelTelemetry;
    wheelRL: WheelTelemetry;
    wheelRR: WheelTelemetry;
};

type TelemetryStatus = {
    running: boolean;
    mode: string;
    analysisMode: string;
    address: string;
    port: number;
    packetLength: number;
    rawPackets: number;
    validPackets: number;
    invalidPackets: number;
    parseErrors: number;
    lastDatagramAt: string;
    lastDatagramBytes: number;
    lastDatagramRemote: string;
    lastPacketAt: string;
    lastError: string;
    hasCurrentFrame: boolean;
    recordingActive: boolean;
    recordingBytes: number;
    recordingLimitBytes: number;
    recordingPackets: number;
    recordingTruncated: boolean;
};

type TelemetryReplayStatus = {
    running: boolean;
    paused: boolean;
    sessionId: number;
    speed: number;
    positionMs: number;
    durationMs: number;
    progress01: number;
    packetIndex: number;
    packetCount: number;
    lastError: string;
};

type NetworkInterface = {
    name: string;
    displayName: string;
    address: string;
    isLoopback: boolean;
    isPrivate: boolean;
    isUp: boolean;
};

type SuggestedAction = {
    priority: number;
    category: string;
    item: string;
    direction: string;
    amount: string;
    reason: string;
};

type DetectedEvent = {
    id: string;
    type: string;
    severity: 'low' | 'medium' | 'high' | string;
    startMs: number;
    endMs: number;
    durationMs: number;
    segment: string;
    evidence: Record<string, number>;
    suggestedActions: SuggestedAction[];
};

type TuneProfileInput = {
    carName: string;
    carOrdinal?: number | null;
    carCategory?: number | null;
    carClass: string;
    pi?: number | null;
    drivetrain: string;
    numCylinders?: number | null;
    useCase: string;
    versionName: string;
    powerKW?: number | null;
    torqueNM?: number | null;
    weightKG?: number | null;
    frontWeightPct?: number | null;
    powerToWeightKWPerKG?: number | null;
    peakTorqueRPM?: number | null;
    peakPowerRPM?: number | null;
    redlineRPM?: number | null;
    frontTirePressure?: number | null;
    rearTirePressure?: number | null;
    finalDrive?: number | null;
    gear1?: number | null;
    gear2?: number | null;
    gear3?: number | null;
    gear4?: number | null;
    gear5?: number | null;
    gear6?: number | null;
    gear7?: number | null;
    gear8?: number | null;
    gear9?: number | null;
    gear10?: number | null;
    frontCamber?: number | null;
    rearCamber?: number | null;
    frontToe?: number | null;
    rearToe?: number | null;
    caster?: number | null;
    frontArb?: number | null;
    rearArb?: number | null;
    frontSpring?: number | null;
    rearSpring?: number | null;
    frontRideHeight?: number | null;
    rearRideHeight?: number | null;
    frontRebound?: number | null;
    rearRebound?: number | null;
    frontBump?: number | null;
    rearBump?: number | null;
    frontAero?: number | null;
    rearAero?: number | null;
    aeroBalance?: number | null;
    brakeBalance?: number | null;
    brakePressure?: number | null;
    frontDiffAccel?: number | null;
    frontDiffDecel?: number | null;
    rearDiffAccel?: number | null;
    rearDiffDecel?: number | null;
    centerDiffBalance?: number | null;
    notes: string;
};

type TuneProfile = TuneProfileInput & {
    id: number;
    createdAt: string;
    updatedAt: string;
};

type RoadStaticTuneBaselineInput = {
    carName: string;
    versionName?: string;
    useCase?: string;
    carOrdinal?: number | null;
    carCategory?: number | null;
    pi: number;
    drivetrain: string;
    tireCompound?: string;
    weightKG: number;
    frontWeightPct: number;
    powerKW?: number | null;
    torqueNM?: number | null;
    redlineRPM?: number | null;
    gearCount?: number | null;
    tireDiameterCm?: number | null;
    targetTopSpeedKmh?: number | null;
    frontRideHeightMinCm?: number | null;
    frontRideHeightMaxCm?: number | null;
    rearRideHeightMinCm?: number | null;
    rearRideHeightMaxCm?: number | null;
    frontAeroMinKgf?: number | null;
    frontAeroMaxKgf?: number | null;
    rearAeroMinKgf?: number | null;
    rearAeroMaxKgf?: number | null;
    frontRideHeightAdjustable?: boolean;
    rearRideHeightAdjustable?: boolean;
    frontAeroAdjustable?: boolean;
    rearAeroAdjustable?: boolean;
    balanceBias?: number;
    stiffnessBias?: number;
    speedBias?: number;
};

type RoadStaticTuneBaselineForm = Record<keyof RoadStaticTuneBaselineInput, string> & {
    useCase: string;
    gearingEnabled: string;
    tireWidthMm: string;
    tireAspectRatio: string;
    tireRimInches: string;
};
type RoadStaticTuneBaselineFormKey = keyof RoadStaticTuneBaselineForm;
type QuickTuneFieldErrors = Partial<Record<RoadStaticTuneBaselineFormKey, string>>;

type BaselineGeneratedField = {
    fieldKey: keyof TuneProfileInput | string;
    group: string;
    value?: number | null;
    unit: string;
    reason: string;
    defaultSelected: boolean;
};

type BaselineSkippedField = {
    fieldKey: keyof TuneProfileInput | string;
    group: string;
    reason: string;
    message: string;
};

type BaselineTierRecommendation = {
    fieldKey: keyof TuneProfileInput | string;
    group: string;
    tier: 'low' | 'medium' | 'high' | string;
    reason: string;
    applicable: boolean;
};

type RoadStaticTuneBaselineResult = {
    profileDraft: TuneProfileInput;
    confidence: string;
    generatedFields: BaselineGeneratedField[];
    tierRecommendations: BaselineTierRecommendation[];
    skippedFields: BaselineSkippedField[];
    warnings: string[];
    nextTestPlan: string[];
};

type RoadStaticTuneBaselineApplyResult = {
    profile: TuneProfile;
    appliedFields: string[];
    skippedFields: BaselineSkippedField[];
};

type RecommendedCarInput = {
    id: string;
    name: string;
    useCase: string;
    useCaseLabel: string;
    pi: number;
    carClass: string;
    drivetrain: string;
    tireCompound: string;
    tireCompoundLabel: string;
    weightKG: number;
    frontWeightPct: number;
    tuneCode: string;
    imageSrc?: string;
    tags: string[];
    reason: string;
};

type RecommendedCar = RecommendedCarInput & {
    createdAt: string;
    updatedAt: string;
};

type RecommendedCarForm = {
    id: string;
    name: string;
    useCase: string;
    useCaseLabel: string;
    pi: string;
    carClass: string;
    drivetrain: string;
    tireCompound: string;
    tireCompoundLabel: string;
    weightKG: string;
    frontWeightPct: string;
    tuneCode: string;
    imageSrc: string;
    tags: string;
    reason: string;
};

type RecommendedCarsFileResult = {
    path: string;
    count: number;
};

type RecommendedCarsFileSelection = {
    path: string;
    exists: boolean;
    version: string;
    ids: string[];
    tuneCodes: string[];
    count: number;
};

type TuneHarvestOptions = {
    sources: string[];
    dryRun: boolean;
    limitPerSource: number;
};

type TuneHarvestRun = {
    id: number;
    startedAt: string;
    finishedAt: string;
    sources: string[];
    dryRun: boolean;
    status: string;
    message: string;
    foundCount: number;
    savedCount: number;
    rejectedCount: number;
    pendingCount: number;
    importedCount: number;
};

type TuneHarvestCandidate = {
    id: number;
    runId: number;
    source: string;
    sourceRef: string;
    sourceUrl: string;
    sourceCarId: string;
    rawKey: string;
    shareCode: string;
    year: number;
    make: string;
    model: string;
    carName: string;
    matchedCarId: string;
    matchScore: number;
    matchReason: string;
    useCase: string;
    carClass: string;
    pi: number;
    drivetrain: string;
    tireCompound: string;
    tuner: string;
    tuneName: string;
    bestFor: string;
    difficulty: string;
    notes: string;
    rawJson: string;
    status: string;
    rejectionReason: string;
    createdAt: string;
    updatedAt: string;
};

type TuneHarvestRunResult = {
    run?: TuneHarvestRun | null;
    candidates: TuneHarvestCandidate[];
    found: number;
    saved: number;
    rejected: number;
    pending: number;
    imported: number;
    warnings: string[];
};

type TuneHarvestSourceState = Record<'jsr_chronic_sheet' | 'codmunity' | 'forzafire', boolean>;

const tuneHarvestCandidateListLimit = 5000;

type TuneProfileSnapshot = {
    id: number;
    tuneProfileId: number;
    sessionId?: number | null;
    changedAt: string;
    changeReason: string;
    before: TuneProfile;
    after: TuneProfile;
    changedFields: string[];
    changeJson: string;
};

type TelemetrySession = {
    id: number;
    tuneProfileId?: number | null;
    tuneSnapshotJson: string;
    tuneName: string;
    sessionName: string;
    trackName: string;
    mode: string;
    gameMode: string;
    startedAt: string;
    endedAt: string;
    durationMs: number;
    bestLapMs?: number | null;
    avgSpeedKmh?: number | null;
    maxSpeedKmh?: number | null;
    eventCount: number;
    sampleCount: number;
    recordingPath: string;
    recordingPackets: number;
    recordingBytes: number;
    recordingTruncated: boolean;
    carOrdinal?: number | null;
    carClass: string;
    carPi?: number | null;
    drivetrain: string;
    numCylinders?: number | null;
    driverMode: string;
    driverModeConfidence: number;
    driverModeEvidenceJson: string;
    brakeAssist: string;
    steeringAssist: string;
    tractionControl: string;
    stabilityControl: string;
    shifting: string;
    launchControl: string;
    driverFeedbackJson?: string;
    notes: string;
};

type TestConditions = {
    driverMode: string;
    brakeAssist: string;
    steeringAssist: string;
    tractionControl: string;
    stabilityControl: string;
    shifting: string;
    launchControl: string;
};

type TuneProfileSessionStat = {
    tuneProfileId: number;
    sessionCount: number;
    lastStartedAt: string;
};

type RuleThresholdProfileInput = {
    name: string;
    carClass: string;
    drivetrain: string;
    useCase: string;
    gameMode?: string;
    configJson: string;
};

type RuleThresholdProfile = RuleThresholdProfileInput & {
    id: number;
    isDefault: boolean;
    createdAt: string;
    updatedAt: string;
};

type IssueEvidence = {
    min: number;
    max: number;
    avg: number;
    count: number;
};

type GearPowerBand = {
    gear: number;
    sampleCount: number;
    highLoadSampleCount: number;
    speedMinKmh: number;
    speedMaxKmh: number;
    speedAvgKmh: number;
    rpmMin: number;
    rpmMax: number;
    rpmAvg: number;
    rpmRatioMin: number;
    rpmRatioMax: number;
    rpmRatioAvg: number;
    inPowerBandRpmMin: number;
    inPowerBandRpmMax: number;
    inPowerBandRatioMin: number;
    inPowerBandRatioMax: number;
    throttleAvg: number;
    accelAvgMps2: number;
    accelMaxMps2: number;
    speedPer1000RpmKmh: number;
    shiftAfterRPM: number;
    shiftDropRPM: number;
    frontSlipAvg: number;
    rearSlipAvg: number;
    frontTractionLimitedPct: number;
    rearTractionLimitedPct: number;
    belowPowerBandPercent: number;
    inPowerBandPercent: number;
    abovePowerBandPercent: number;
    lowRpmHighLoadPercent: number;
    highRpmHighLoadPercent: number;
    tractionLimitedPercent: number;
    finding: string;
};

type GearPowerComparisonRow = {
    item: string;
    gear: number;
    beforeValue: number;
    afterValue: number;
    deltaValue: number;
    beforeSpeedMaxKmh: number;
    afterSpeedMaxKmh: number;
    speedMaxDeltaKmh: number;
    beforeInPowerBandPct: number;
    afterInPowerBandPct: number;
    inPowerBandDeltaPct: number;
    beforeTractionLimitPct: number;
    afterTractionLimitPct: number;
    tractionLimitDeltaPct: number;
    beforeFinding: string;
    afterFinding: string;
};

type GearPowerComparison = {
    type: string;
    status: string;
    baselineSessionId: number;
    rows: GearPowerComparisonRow[];
};

type GearPowerDiagnostic = {
    status: string;
    summary: string;
    launchFinding: string;
    topSpeedFinding: string;
    powerKW: number;
    torqueNM: number;
    weightKG: number;
    frontWeightPct: number;
    powerToWeightKWPerKG: number;
    powerToWeightBand: string;
    peakTorqueRPM: number;
    peakPowerRPM: number;
    redlineRPM: number;
    powerBandStartRPM: number;
    powerBandEndRPM: number;
    powerBandSource: string;
    confidence: string;
    strategyMode: string;
    globalGearIssueCount: number;
    usableGearCount: number;
    globalGearIssueRatio: number;
    tractionLimitedPercent: number;
    lowRpmHighLoadPercent: number;
    highRpmHighLoadPercent: number;
    gears: GearPowerBand[];
    comparisons: GearPowerComparison[];
    recommendedActions: SuggestedAction[];
    evidence: Record<string, number>;
};

type WholeCarAdjustment = SuggestedAction & {
    family: string;
    source: string;
    confidence: string;
    evidence: Record<string, number>;
};

type TuningConflict = {
    key: string;
    keptItem: string;
    droppedItem: string;
    reason: string;
};

type WholeCarTuningPlan = {
    strategy: string;
    confidence: string;
    summary: string;
    actions: WholeCarAdjustment[];
    conflicts: TuningConflict[];
    notes: string[];
};

type TunePlanDraftAction = {
    id: string;
    family: string;
    source: string;
    confidence: string;
    adviceLayer: string;
    trustLevel: string;
    trustReasons: string[];
    missingInputs: string[];
    retestGuard: string;
    rationale: string;
    conflictReason: string;
    category: string;
    item: string;
    fieldKey: string;
    direction: string;
    reason: string;
    currentValue?: number;
    targetValue?: number;
    delta?: number;
    unit: string;
    step: number;
    canApply: boolean;
    blockedReason: string;
};

type TunePlanDraft = {
    sessionId: number;
    tuneProfileId?: number | null;
    status: string;
    summary: string;
    actions: TunePlanDraftAction[];
    conflicts: TuningConflict[];
};

type RoadTuningKnowledgeStatus = {
    loadedAt: string;
    sourcePath: string;
    lastError: string;
    symptomCount: number;
    actionCount: number;
    usingFallback: boolean;
};

type RoadTuningDecisionAction = {
    id: string;
    role: string;
    family: string;
    source: string;
    confidence: string;
    trustLevel: string;
    adviceLayer: string;
    category: string;
    item: string;
    fieldKey: string;
    direction: string;
    amount: string;
    unit: string;
    reason: string;
    rationale: string;
    conflictReason: string;
    canAutoApply: boolean;
    blockedReason: string;
    evidence: Record<string, number>;
};

type RoadTuningDecision = {
    sessionId: number;
    status: string;
    symptomId: string;
    phase: string;
    symptom: string;
    primaryCause: string;
    confidence: string;
    fitVerdict: string;
    reason: string;
    rollbackRecommended: boolean;
    relatedIssueGroup?: SessionIssueGroup | null;
    evidence: Record<string, number>;
    actions: RoadTuningDecisionAction[];
    retestFocus: string[];
    knowledgeStatus: RoadTuningKnowledgeStatus;
};

type RetestMetric = {
    key: string;
    current: number;
    baseline: number;
    delta: number;
    direction: string;
    status: string;
};

type RetestEvaluation = {
    sessionId: number;
    baselineSession?: TelemetrySession | null;
    status: string;
    summary: string;
    confidence: string;
    baselineReason: string;
    changedFields: string[];
    changeSourceSessionId?: number | null;
    rollbackActions: TunePlanDraftAction[];
    metricSummary: string[];
    metrics: RetestMetric[];
};

type SessionIssueGroup = {
    id: string;
    family: string;
    severity: string;
    segment: string;
    eventTypes: string[];
    eventIds: string[];
    events: DetectedEvent[];
    eventCount: number;
    totalDurationMs: number;
    firstStartMs: number;
    lastEndMs: number;
    evidence: Record<string, IssueEvidence>;
    primaryActions: SuggestedAction[];
    comparison: string;
    baselineEventCount: number;
    baselineTotalDurationMs: number;
    relatedRecentChanges: string[];
    prioritizeTuning: boolean;
    adjustmentStrategy?: string;
    feedbackDirective?: string;
};

type SessionIssueSummary = {
    sessionId: number;
    baselineSession?: TelemetrySession | null;
    baselineStatus: string;
    recentChangeFields: string[];
    groups: SessionIssueGroup[];
    gearPower: GearPowerDiagnostic;
    wholeCarPlan: WholeCarTuningPlan;
};

type QuickLapSummary = {
    lapNumber: number;
    sampleCount: number;
    durationMs: number;
    avgSpeedKmh: number;
    maxSpeedKmh: number;
    eventCount: number;
    issueScore: number;
};

type QuickSuggestion = {
    family: string;
    source: string;
    confidence: string;
    trustLevel: string;
    adviceLayer: string;
    category: string;
    item: string;
    direction: string;
    amount: string;
    reason: string;
    rationale: string;
    nextStep: string;
    fieldKeys: string[];
    missingInputs: string[];
    canApply: boolean;
    blockedReason: string;
};

type QuickVehicleSnapshot = {
    carOrdinal?: number | null;
    carClass: string;
    carPi?: number | null;
    drivetrain: string;
    numCylinders?: number | null;
};

type QuickComparability = {
    sameVehicleClass: string;
    sameTrackContext: string;
    confidence: string;
    warnings: string[];
    baselineVehicle: QuickVehicleSnapshot;
    currentVehicle: QuickVehicleSnapshot;
};

type QuickDiagnostic = {
    status: string;
    comparisonStatus: string;
    updatedAt: string;
    sampleCount: number;
    eventCount: number;
    gameMode: string;
    driverMode: string;
    driverModeConfidence: number;
    vehicle: QuickVehicleSnapshot;
    comparability: QuickComparability;
    currentLap?: QuickLapSummary | null;
    previousLap?: QuickLapSummary | null;
    groups: SessionIssueGroup[];
    gearPower: GearPowerDiagnostic;
    suggestions: QuickSuggestion[];
    missingProfileFields: string[];
};

type TireWheelDiagnostic = {
    position: string;
    combinedSlipAvg: number;
    combinedSlipMax: number;
    combinedSlipP90: number;
    combinedSlipHighPct: number;
    slipRatioAvg: number;
    slipRatioMax: number;
    slipRatioP90: number;
    slipRatioHighPct: number;
    slipAngleAvg: number;
    slipAngleMax: number;
    slipAngleP90: number;
    tireTempAvg: number;
    tireTempMax: number;
    suspensionTravelAvg: number;
    suspensionTravelMax: number;
    suspensionOffsetPctAvg: number;
    suspensionOffsetPctMax: number;
    suspensionTravelMetersAvg: number;
    suspensionTravelMetersMax: number;
    gripState: string;
};

type TireAxleDiagnostic = {
    name: string;
    combinedSlipAvg: number;
    combinedSlipMax: number;
    combinedSlipP90: number;
    combinedSlipHighPct: number;
    slipRatioAvg: number;
    slipRatioMax: number;
    slipRatioP90: number;
    slipRatioHighPct: number;
    slipAngleAvg: number;
    slipAngleMax: number;
    slipAngleP90: number;
    tireTempAvg: number;
    tireTempMax: number;
    suspensionTravelAvg: number;
    suspensionTravelMax: number;
    suspensionOffsetPctAvg: number;
    suspensionOffsetPctMax: number;
    limitScore: number;
    gripState: string;
};

type TireSideBalance = {
    leftCombinedSlipAvg: number;
    rightCombinedSlipAvg: number;
    delta: number;
    state: string;
};

type GForceDiagnostic = {
    source: string;
    axisMapping: string;
    currentXG: number;
    currentYG: number;
    currentZG: number;
    currentTotalG: number;
    avgAbsXG: number;
    avgAbsYG: number;
    avgAbsZG: number;
    avgTotalG: number;
    peakAbsXG: number;
    peakAbsYG: number;
    peakAbsZG: number;
    peakTotalG: number;
    dominantAxis: string;
    series: GForcePoint[];
};

type GForcePoint = {
    timeMs: number;
    xG: number;
    yG: number;
    zG: number;
    totalG: number;
};

type CamberInference = {
    status: string;
    confidence: string;
    frontState: string;
    rearState: string;
    summary: string;
    explanation: string;
    warnings: string[];
    hints: TireModelHint[];
    evidence: Record<string, number>;
};

type PowerToTireDiagnostic = {
    status: string;
    summary: string;
    explanation: string;
    confidence: string;
    sampleCount: number;
    highThrottleSampleCount: number;
    drivetrain: string;
    drivenAxle: string;
    currentPowerKW: number;
    averagePowerKW: number;
    maxPowerKW: number;
    currentTorqueNM: number;
    averageTorqueNM: number;
    maxTorqueNM: number;
    currentRPM: number;
    averageRPM: number;
    currentRPMRatio: number;
    averageRPMRatio: number;
    currentGear: number;
    averageThrottle: number;
    averageSpeedKmh: number;
    speedDeltaKmh: number;
    averageAccelMps2: number;
    averageAccelG: number;
    peakAccelG: number;
    frontSlipRatioP90: number;
    rearSlipRatioP90: number;
    drivenSlipRatioP90: number;
    drivenSlipRatioHighPct: number;
    rpmLowHighThrottlePct: number;
    rpmHighHighThrottlePct: number;
    powerSignalAvailable: boolean;
    tractionLimited: boolean;
    evidence: Record<string, number>;
};

type BrakeToTireDiagnostic = {
    status: string;
    summary: string;
    explanation: string;
    confidence: string;
    sampleCount: number;
    brakeSampleCount: number;
    averageBrake: number;
    peakBrake: number;
    averageHandBrake: number;
    peakHandBrake: number;
    averageSpeedKmh: number;
    speedDeltaKmh: number;
    averageSteer: number;
    averageDecelMps2: number;
    averageDecelG: number;
    peakDecelG: number;
    averagePlaneG: number;
    peakPlaneG: number;
    frontSlipRatioP90: number;
    rearSlipRatioP90: number;
    frontCombinedSlipP90: number;
    rearCombinedSlipP90: number;
    frontRearSlipDelta: number;
    trailBraking: boolean;
    handbrakeActive: boolean;
    evidence: Record<string, number>;
};

type TireModelHint = {
    code: string;
    severity: string;
    direction: string;
    reason: string;
};

type TireDataQuality = {
    status: string;
    confidence: string;
    sampleCount: number;
    dynamicSampleCount: number;
    speedSignal: string;
    gForceSignal: string;
    slipSignal: string;
    inputSignal: string;
    reasons: string[];
    evidence: Record<string, number>;
};

type TireGripLimit = {
    type: string;
    limitedAxle: string;
    limitedWheels: string[];
    primaryEvidence: string;
    confidence: string;
    reason: string;
    frontRearDelta: number;
    drivenDelta: number;
    leftRightDelta: number;
    evidence: Record<string, number>;
};

type TirePhaseDiagnostic = {
    currentPhase: string;
    secondaryPhase: string;
    stablePhase: string;
    phaseStability: string;
    scoreMargin: number;
    confidence: string;
    scores: Record<string, number>;
    evidence: Record<string, number>;
    windowMs: number;
    sampleCount: number;
};

type TireIssueSegment = {
    id: string;
    type: string;
    phase: string;
    operationTags: string[];
    driftSource: string;
    limitType: string;
    limitedAxle: string;
    limitedWheels: string[];
    startMs: number;
    endMs: number;
    durationMs: number;
    sampleCount: number;
    speedMinKmh: number;
    speedMaxKmh: number;
    speedAvgKmh: number;
    confidence: string;
    dataQuality: string;
    riskLevel: string;
    evidence: Record<string, number>;
    reason: string;
};

type TireIssueGroup = {
    id: string;
    type: string;
    phase: string;
    operationTags: string[];
    driftSource: string;
    limitType: string;
    limitedAxle: string;
    limitedWheels: string[];
    count: number;
    totalDurationMs: number;
    speedMinKmh: number;
    speedMaxKmh: number;
    speedAvgKmh: number;
    confidence: string;
    dataQuality: string;
    riskLevel: string;
    representativeEvidence: Record<string, number>;
    segmentIds: string[];
    reason: string;
};

type TireIssueAnalysis = {
    status: string;
    updatedAt: string;
    windowMs: number;
    sampleCount: number;
    segmentCount: number;
    groupCount: number;
    segments: TireIssueSegment[];
    groups: TireIssueGroup[];
    warnings: string[];
};

type TireIssueAdviceAction = {
    id: string;
    issueGroupId: string;
    layer: string;
    category: string;
    scope: string;
    direction: string;
    relatedFields: string[];
    rationale: string;
    verifyEvidence: string[];
    confidence: string;
    missingInputs: string[];
    conflictReason: string;
    tuneRecommended: boolean;
};

type TireIssueAdviceGroup = {
    issueGroupId: string;
    issueType: string;
    phase: string;
    operationTags: string[];
    limitedAxle: string;
    driftSource: string;
    primaryCause: string;
    shouldTune: boolean;
    priority: number;
    confidence: string;
    evidence: Record<string, number>;
    actions: TireIssueAdviceAction[];
};

type TireIssueAdvice = {
    status: string;
    updatedAt: string;
    confidence: string;
    basedOnIssueUpdatedAt: string;
    issueGroupCount: number;
    priorityActions: TireIssueAdviceAction[];
    groups: TireIssueAdviceGroup[];
    warnings: string[];
};

type TireModelDiagnostic = {
    status: string;
    updatedAt: string;
    sampleCount: number;
    windowMs: number;
    gameMode: string;
    phase: string;
    phaseDetail: TirePhaseDiagnostic;
    dataQuality: TireDataQuality;
    gripLimit: TireGripLimit;
    limitType: string;
    confidence: string;
    summary: string;
    explanation: string;
    warnings: string[];
    wheels: TireWheelDiagnostic[];
    frontAxle: TireAxleDiagnostic;
    rearAxle: TireAxleDiagnostic;
    leftRight: TireSideBalance;
    gForce: GForceDiagnostic;
    camber: CamberInference;
    powerToTire: PowerToTireDiagnostic;
    brakeToTire: BrakeToTireDiagnostic;
    issueAnalysis?: TireIssueAnalysis;
    issueAdvice?: TireIssueAdvice;
    hints: TireModelHint[];
    evidence: Record<string, number>;
    vehicle: QuickVehicleSnapshot;
};

type TireDiagnosticSnapshot = {
    generatedAt: string;
    status: string;
    sampleCount: number;
    windowMs: number;
    vehicle: QuickVehicleSnapshot;
    dataQuality: {
        status: string;
        confidence: string;
        sampleCount: number;
        dynamicSampleCount: number;
        reasons: string[];
    };
    phase: {
        current: string;
        stable: string;
        secondary: string;
        stability: string;
        confidence: string;
        scoreMargin: number;
    };
    gripLimit: {
        type: string;
        limitedAxle: string;
        limitedWheels: string[];
        confidence: string;
        primaryEvidence: string;
        reason: string;
    };
    issueAnalysis?: TireIssueAnalysis;
    issueAdvice?: TireIssueAdvice;
    risks: string[];
    power: {
        status: string;
        summary: string;
        confidence: string;
        explanation: string;
        evidence: Record<string, number>;
    };
    brake: {
        status: string;
        summary: string;
        confidence: string;
        explanation: string;
        evidence: Record<string, number>;
    };
    evidence: Record<string, number>;
};

type TireRegressionExpectation = {
    allowedPhases: string[];
    requiredGripTypes: string[];
    allowedAxles: string[];
    forbiddenGripTypes: string[];
    minDataQuality: string;
    notes: string;
};

type TireRegressionSampleInput = {
    name: string;
    scenario: string;
    windowSeconds: number;
    expected: TireRegressionExpectation;
};

type TireRegressionSampleSummary = {
    id: string;
    name: string;
    scenario: string;
    createdAt: string;
    windowSeconds: number;
    vehicle: QuickVehicleSnapshot;
    sampleCount: number;
    expected: TireRegressionExpectation;
};

type TireRegressionSample = TireRegressionSampleSummary & {
    samples: TelemetryFrame[];
    snapshot: TireDiagnosticSnapshot;
};

type TireRegressionResult = {
    sampleId: string;
    name: string;
    scenario: string;
    passed: boolean;
    status: string;
    failures: string[];
    expected: TireRegressionExpectation;
    actual: TireDiagnosticSnapshot;
};

type TireRegressionSaveFormState = {
    name: string;
    scenario: string;
    windowSeconds: string;
};

type TireRegressionExpectedFormState = {
    allowedPhases: string;
    requiredGripTypes: string;
    allowedAxles: string;
    forbiddenGripTypes: string;
    minDataQuality: string;
    notes: string;
};

type TuningPipelineComponent = {
    id: string;
    name: string;
    description: string;
    sourceTypes?: string[];
    compatibleWith?: string[];
    tags?: string[];
};

type TuningPipelineCombination = {
    sourceType: string;
    detectorId: string;
    decisionerId: string;
    interpreterId: string;
    description: string;
};

type TuningPipelineCatalog = {
    sourceTypes: TuningPipelineComponent[];
    detectors: TuningPipelineComponent[];
    decisioners: TuningPipelineComponent[];
    interpreters: TuningPipelineComponent[];
    defaultCombinations: TuningPipelineCombination[];
};

type TuningPipelineRunInput = {
    sourceType: string;
    sessionId?: number;
    detectorId: string;
    decisionerId: string;
    interpreterId: string;
};

type TuningPipelineSourceSummary = {
    sourceType: string;
    sessionId?: number;
    sampleCount: number;
    eventCount: number;
    vehicle: QuickVehicleSnapshot;
    gameMode: string;
    driverMode: string;
    label: string;
};

type TuningProblem = {
    id: string;
    sourceId: string;
    family: string;
    type: string;
    phase: string;
    operationTags: string[];
    limitedAxle: string;
    limitedWheels: string[];
    severity: string;
    confidence: string;
    riskLevel: string;
    count: number;
    durationMs: number;
    summary: string;
    reason: string;
    evidence: Record<string, number>;
};

type TuningProblemSet = {
    detectorId: string;
    status: string;
    problems: TuningProblem[];
    warnings: string[];
};

type TuningDecision = {
    id: string;
    problemId: string;
    phase: string;
    primaryCause: string;
    shouldTune: boolean;
    confidence: string;
    rationale: string;
    documentContext: string;
    evidence: Record<string, number>;
};

type TuningDecisionSet = {
    decisionerId: string;
    status: string;
    decisions: TuningDecision[];
    warnings: string[];
};

type TuningAdvice = {
    id: string;
    decisionId: string;
    problemId: string;
    layer: string;
    category: string;
    scope: string;
    direction: string;
    relatedFields: string[];
    rationale: string;
    verifyEvidence: string[];
    trustLevel: string;
    missingInputs: string[];
    conflictReason: string;
    canApply: boolean;
    documentSources: string[];
    evidence: Record<string, number>;
};

type TuningAdviceSet = {
    interpreterId: string;
    status: string;
    advice: TuningAdvice[];
    documentSources: string[];
    warnings: string[];
};

type TuningPipelineRunResult = {
    status: string;
    updatedAt: string;
    sourceSummary: TuningPipelineSourceSummary;
    problemSet: TuningProblemSet;
    decisionSet: TuningDecisionSet;
    adviceSet: TuningAdviceSet;
    warnings: string[];
};

type ProfessionalPipelineConfig = {
    detectorId: string;
    decisionerId: string;
    interpreterId: string;
};

type ProfessionalTuningDiagnostic = {
    status: string;
    updatedAt: string;
    config: ProfessionalPipelineConfig;
    pipeline?: TuningPipelineRunResult | null;
    warnings: string[];
};

type TuneFieldInfluence = {
    fieldKey: string;
    category: string;
    labelZh: string;
    labelEn: string;
    influenceType: string;
    scope: string[];
    phases: string[];
    tireMetrics: string[];
    evidenceKeys: string[];
    sideEffects: string[];
    conditions: string[];
    summaryZh: string;
    summaryEn: string;
};

type TuneToTireInfluenceMap = {
    version: string;
    items: TuneFieldInfluence[];
};

type StrategyTemplate = {
    id: number;
    name: string;
    carClass: string;
    drivetrain: string;
    useCase: string;
    gameMode: string;
    isDefault: boolean;
    enabledEventCount: number;
    totalEventCount: number;
    updatedAt: string;
};

type StrategyEventDistribution = {
    type: string;
    count: number;
    severity: string;
};

type StrategyIssueAggregate = {
    family: string;
    eventCount: number;
    sessionCount: number;
    severity: string;
    recommendation: string;
};

type StrategyAnalysisHint = {
    level: string;
    message: string;
    eventType?: string;
    family?: string;
};

type RoadStrategyAnalysis = {
    template: StrategyTemplate;
    sessionIds: number[];
    sessionCount: number;
    totalEvents: number;
    eventDistribution: StrategyEventDistribution[];
    issueGroups: StrategyIssueAggregate[];
    hints: StrategyAnalysisHint[];
};

type UpgradeUnlockRule = {
    category: string;
    upgradeName: string;
    unlocks: string;
};

type TuneAdjustmentExplanation = {
    category: string;
    item: string;
    detail: string;
    description: string;
};

type SessionComparisonMetric = {
    key: string;
    label: string;
    unit: string;
    left: number;
    right: number;
    delta: number;
    higherIsBetter: boolean;
};

type SessionEventComparison = {
    type: string;
    left: number;
    right: number;
    delta: number;
};

type SessionComparison = {
    leftSession: TelemetrySession;
    rightSession: TelemetrySession;
    metrics: SessionComparisonMetric[];
    eventTypes: SessionEventComparison[];
    comparabilityWarnings: string[];
};

type BenchmarkPoint = {
    x: number;
    y: number;
    z: number;
};

type TrackCapturePoint = BenchmarkPoint & {
    lapNumber: number;
    currentLap: number;
    currentRaceTime: number;
};

type BenchmarkTrackType = 'auto' | 'circuit' | 'sprint';
type BenchmarkExtractionMode = 'auto_best_lap' | 'first_lap' | 'full_segment';

type BenchmarkGate = {
    center: BenchmarkPoint;
    directionX: number;
    directionZ: number;
    widthMeters: number;
    depthMeters: number;
};

type BenchmarkTrackInput = {
    name: string;
    sourceMode: string;
    trackType: BenchmarkTrackType;
    start: BenchmarkPoint;
    end: BenchmarkPoint;
    startRadius: number;
    endRadius: number;
    directionX: number;
    directionZ: number;
    startGate: BenchmarkGate;
    finishGate: BenchmarkGate;
    checkpoints: BenchmarkPoint[];
    routeLengthMeters: number;
    hasDrivingLine: boolean;
    polyline: BenchmarkPoint[];
    sourceSessionId?: number | null;
    lapCountObserved: number;
    notes: string;
};

type BenchmarkTrackExtractionInput = {
    sessionId: number;
    name: string;
    trackType: BenchmarkTrackType;
    extractionMode: BenchmarkExtractionMode;
    startGate?: BenchmarkGate;
    finishGate?: BenchmarkGate;
};

type BenchmarkTrack = BenchmarkTrackInput & {
    id: number;
    createdAt: string;
    updatedAt: string;
};

type BenchmarkRun = {
    id: number;
    sessionId: number;
    trackId: number;
    trackName: string;
    startMs: number;
    endMs: number;
    durationMs: number;
    confidence: number;
    avgSpeedKmh?: number | null;
    maxSpeedKmh?: number | null;
    routeProgress01?: number | null;
    geometryLengthMeters?: number | null;
    trackLengthErrorPct?: number | null;
    distanceTraveledDeltaMeters?: number | null;
    currentRaceTimeDeltaSeconds?: number | null;
    avgLateralErrorMeters?: number | null;
    maxLateralErrorMeters?: number | null;
    warningFlags: string;
    eventCount: number;
    driverMode: string;
    driverModeConfidence: number;
    driverModeEvidenceJson: string;
    valid: boolean;
    createdAt: string;
};

type TrackBaselineRun = Omit<BenchmarkRun, 'sessionId' | 'trackName'> & {
    vehicle: TrackVehicleKey;
    gameMode: string;
};

type TrackBaselineSaveResult = {
    track: BenchmarkTrack;
    baseline: TrackBaselineRun;
    action: string;
    matchCandidate?: TrackMergeCandidate | null;
};

type TrackVehicleKey = {
    carOrdinal?: number | null;
    carClass: string;
    carPi?: number | null;
    drivetrain: string;
    label: string;
};

type TrackRunContext = {
    run: BenchmarkRun;
    session: TelemetrySession;
    vehicle: TrackVehicleKey;
};

type TrackAutoBaseline = {
    vehicle: TrackVehicleKey;
    bestRun: TrackRunContext;
    recentRuns: TrackRunContext[];
    runCount: number;
};

type TrackVehicleReference = {
    vehicle: TrackVehicleKey;
    bestAutoBaseline?: TrackRunContext | null;
    bestTrackBaseline?: TrackBaselineRun | null;
    recentRuns: TrackRunContext[];
    recentBaselineRuns: TrackBaselineRun[];
    validRunCount: number;
    autoRunCount: number;
    baselineRunCount: number;
    avgSpeedKmh?: number | null;
    maxSpeedKmh?: number | null;
    eventCount: number;
};

type TrackProfile = {
    track: BenchmarkTrack;
    autoBaselines: TrackAutoBaseline[];
    vehicleReferences: TrackVehicleReference[];
    recentRuns: TrackRunContext[];
    warnings: string[];
};

type TrackMergeCandidate = {
    track: BenchmarkTrack;
    matchLevel: string;
    lengthErrorPct: number;
    startDistanceMeters: number;
    endDistanceMeters: number;
    shapeSimilarity: number;
    routeFitAvgErrorMeters: number;
    routeFitP90ErrorMeters: number;
    routeFitScore: number;
    directionMatched: boolean;
    reverseMatched: boolean;
    reason: string;
};

type RoadEvaluationAttribution = {
    type: string;
    eventType?: string;
    count: number;
    severity?: string;
    priority: number;
    message: string;
    prioritizeTuning: boolean;
};

type TuneWebServerStatus = {
    running: boolean;
    port: number;
    url: string;
    lanAddress: string;
    lastError: string;
};

type RoadSessionEvaluation = {
    session: TelemetrySession;
    track?: BenchmarkTrack | null;
    bestRun?: BenchmarkRun | null;
    baselineRun?: BenchmarkRun | null;
    baselineSession?: TelemetrySession | null;
    baselineStatus: string;
    paperPerformanceScore: number;
    playerFitScore: number;
    riskScore: number;
    overallVerdict: string;
    attributions: RoadEvaluationAttribution[];
    notes: string[];
};

type TrackCaptureState = {
    recording: boolean;
    name: string;
    points: TrackCapturePoint[];
    hasDrivingLine: boolean;
};

type ViewName = 'quick' | 'tire_lab' | 'tire_regression' | 'model_pipeline_lab' | 'track_profiles' | 'tune_generator' | 'remote_tune' | 'expert' | 'reports' | 'developer';
type DeveloperToolName = 'do_fields' | 'tire_lab' | 'tire_regression' | 'model_pipeline_lab' | 'track_profiles' | 'recommended_cars' | 'tune_harvest' | 'strategy';
type GameMode = 'unknown' | 'menu' | 'free_roam' | 'race' | 'mixed';

type PendingStartChoice = {
    profiles: TuneProfile[];
    address: string;
    port: number;
    carOrdinal: number;
    carClass: string;
    carPi: number;
};

type PendingProfileMismatch = {
    candidates: TuneProfile[];
    address: string;
    port: number;
    telemetry: TelemetryFrame;
    profile: TuneProfile;
};

type PendingSessionBind = {
    session: TelemetrySession;
    profiles: TuneProfile[];
};

type PendingTrackMerge = {
    input: BenchmarkTrackInput;
    candidates: TrackMergeCandidate[];
};

type StartValidationResult =
    | { action: 'start'; notice?: string }
    | { action: 'deferred' };

type TelemetryStartMode = 'quick' | 'professional' | 'tire_lab' | 'track_capture' | 'expert';

const COPY = {
    en: {
        title: 'FH6 Vehicle Tuning Tool',
        receiving: 'Receiving',
        idle: 'Idle',
        networkAdapter: 'Network adapter',
        udpPort: 'UDP Port',
        start: 'Start',
        stop: 'Stop',
        ready: 'Ready for FH6 UDP Data Out packets.',
        listening: (address: string, port: number) => `Listening on ${address}:${port}`,
        stopped: 'Listener stopped',
        gameTargetTitle: 'Game Data Out target',
        gameTargetPrefix: 'Set FH6 Data Out IP to this PC address, for example',
        gameTargetMiddle: 'and port',
        gameTargetSuffix: 'The game PC address',
        gameTargetSuffix2: 'is the sender, not the target.',
        allInterfaces: 'All interfaces',
        loopback: 'loopback',
        lan: 'LAN',
        ip: 'IP',
        speed: 'Speed',
        rpm: 'RPM',
        gear: 'Gear',
        yawRate: 'Yaw Rate',
        driverInputs: 'Driver Inputs',
        raceDataActive: 'Race mode',
        waitingRaceState: 'Menu / transition',
        freeRoamMode: 'Free roam',
        menuMode: 'Menu / transition',
        mixedMode: 'Mixed',
        unknownMode: 'Unknown',
        notApplicable: 'N/A in this mode',
        throttle: 'Throttle',
        brake: 'Brake',
        handbrake: 'Handbrake',
        steering: 'Steering',
        validPackets: 'Valid packets',
        invalidPackets: 'Invalid packets',
        parseErrors: 'Parse errors',
        wheelState: 'Wheel State',
        wheelSubtitle: 'Slip, temperature, suspension',
        frontLeft: 'Front Left',
        frontRight: 'Front Right',
        rearLeft: 'Rear Left',
        rearRight: 'Rear Right',
        ratio: 'Ratio',
        angle: 'Angle',
        combined: 'Combined',
        temp: 'Temp',
        susp: 'Susp',
        suspensionOffsetPct: 'Suspension offset %',
        realtimeTrend: 'Realtime Trend',
        trendSubtitle: '10 Hz aggregate, last 8 seconds',
        rpmLoad: 'RPM Load',
        endpoint: 'Endpoint',
        mode: 'Mode',
        udpMode: 'UDP',
        rawPackets: 'UDP packets',
        lastUdpPacket: 'Last UDP packet',
        lastUdpRemote: 'Last UDP sender',
        replayMode: 'Replay',
        idleMode: 'Idle',
        packetSize: 'Packet size',
        lastPacket: 'Last packet',
        speedSource: 'Speed source',
        speedSourcePacket: 'Speed field',
        speedSourceVelocity: 'Velocity vector',
        speedSourceNone: 'No valid speed',
        packetSpeed: 'Packet speed',
        velocitySpeed: 'Velocity speed',
        vehicleId: 'Vehicle ID',
        vehicleCategory: 'Category',
        classPi: 'Class / PI',
        drivetrainCylinders: 'Drivetrain / cylinders',
        engineOutput: 'Power / torque',
        boostFuel: 'Boost / fuel',
        worldPosition: 'World position',
        lapRace: 'Lap / race',
        clutchHandbrake: 'Clutch / handbrake',
        surfaceSignals: 'Surface signals',
        recording: 'Recording',
        recordingSize: 'Recording size',
        recordingPackets: 'Recorded packets',
        recordingLimit: 'Recording limit',
        recordingTruncated: 'Recording truncated',
        recordingReady: 'Recording ready',
        samplesSaved: 'Saved samples',
        replay: 'Replay',
        stopReplay: 'Stop replay',
        replaySpeed: 'Replay speed',
        noRecording: 'No recording',
        historicalTrend: 'Historical Trend',
        inputsTrend: 'Inputs',
        frontRearSlip: 'Front / Rear slip',
        testLaunchpad: 'Test Launchpad',
        testLaunchpadSubtitle: 'Select the correct tune profile and test conditions before recording.',
        analysisMode: 'Analysis mode',
        quickDiagnosis: 'Quick Diagnosis',
        expertTuning: 'Expert Tuning',
        quickModeNote: 'Real-time diagnosis only. No session, recording, or replay will be saved.',
        expertModeNote: 'Full tuning workflow with saved sessions, recordings, reports, and tune plan drafts.',
        quickNoHistory: 'Quick mode: no history or replay',
        quickDiagnosticTitle: 'Quick Real-time Diagnosis',
        quickDiagnosticEmpty: 'Start quick diagnosis and receive telemetry to see live issues.',
        quickSuggestions: 'Suggestion directions',
        noQuickSuggestions: 'No directional suggestions yet.',
        quickSuggestionReason: 'Why',
        quickSuggestionNextStep: 'Next step',
        adviceLayerLabels: {
            rollback: 'Rollback',
            primary: 'Primary issue',
            powertrain: 'Powertrain',
            support: 'Support',
            alternative: 'Alternative',
        },
        quickNextStepLabels: {
            bind_tune_profile_for_values: 'Bind or create a tune profile for exact values.',
            fill_or_unlock_tune_fields: 'Fill or unlock the related tune fields first.',
            collect_power_samples: 'Collect more full-throttle acceleration samples.',
            use_expert_for_concrete_values: 'Use Expert Tuning for exact values and one-click apply.',
        },
        missingProfileFields: 'Fields needed for concrete values',
        sameVehicleClass: 'Same car / class',
        sameTrackContext: 'Same track context',
        comparisonConfidence: 'Comparison confidence',
        comparabilityLabels: {
            yes: 'Yes',
            no: 'No',
            unknown: 'Unknown',
        },
        quickConfidenceLabels: {
            high: 'High',
            medium: 'Medium',
            low: 'Low',
            invalid: 'Invalid',
        },
        quickWarningLabels: {
            quick_lap_data_insufficient: 'Lap data is incomplete. Rolling-window diagnosis is shown.',
            quick_non_race_track_unknown: 'Free roam or non-race telemetry cannot confirm the same track.',
            quick_vehicle_or_class_changed: 'Vehicle ID or class changed during this quick diagnosis.',
            quick_vehicle_class_unknown: 'Vehicle ID or class is incomplete.',
            quick_lap_clock_reset: 'Lap timer reset was detected.',
            quick_race_time_reset: 'Race timer reset was detected.',
            quick_race_time_missing: 'Race timer is missing; comparison confidence is reduced.',
            quick_track_context_unknown: 'Track or race context cannot be confirmed.',
        },
        currentLap: 'Current lap',
        previousLap: 'Previous lap',
        issueScore: 'Issue score',
        quickRollingWindow: 'Lap data is not complete. Showing rolling-window diagnosis.',
        quickComparisonStatuses: {
            lap_comparison: 'Current lap vs previous lap',
            rolling_window_only: 'Rolling window',
            no_data: 'No data',
        },
        testConditionWarning: 'Driver mode is unknown. Auto baseline and player-fit conclusions will be limited.',
        activeProfileMismatch: 'Current telemetry vehicle does not match the selected tune profile.',
        activeProfileReady: 'Tune profile and telemetry vehicle are ready for this test.',
        coreTelemetry: 'Core Telemetry',
        startupFocus: 'Advanced diagnostics, track extraction, and rule thresholds are in Developer Mode.',
        never: 'Never',
        languageLabel: 'Language',
        currentProfile: 'Current tune',
        noProfile: 'No tune profile',
        quickTab: 'Quick Diagnosis',
        tireLabTab: 'Tire Model Lab',
        tireRegressionTab: 'Tire Regression Lab',
        modelPipelineTab: 'Model Pipeline Lab',
        trackProfilesTab: 'Track Profiles',
        recommendedCarsTab: 'Recommended Cars JSON',
        tuneGeneratorTab: 'Quick Tune',
        remoteTuneTab: 'Remote Tune',
        expertTab: 'Professional Tuning',
        dashboardTab: 'Quick Diagnosis',
        profilesTab: 'Expert Tuning',
        reportsTab: 'Reports',
        strategyTab: 'Strategy Templates',
        developerTab: 'Developer Mode',
        developerDoDiagnostics: 'DO diagnostics',
        developerStrategyConfig: 'Strategies / pipeline config',
        recommendedCarsTitle: 'Recommended cars JSON generator',
        recommendedCarsSubtitle: 'Maintain mini-program recommendation data in the local database, then export weChatApp/miniprogram/data/recommendedCars.json.',
        recommendedCarsHint: 'Export uses the database list below and replaces the current recommendedCars.json file.',
        recommendedCarsAdd: 'Save to database',
        recommendedCarsClear: 'Clear form',
        recommendedCarsGenerate: 'Export file',
        recommendedCarsPending: 'Database cars',
        recommendedCarsNoItems: 'No recommended cars in the database yet.',
        recommendedCarsSearch: 'Search recommended cars',
        recommendedCarsSelectVisible: 'Select visible',
        recommendedCarsClearSelection: 'Clear selection',
        recommendedCarsRefresh: 'Refresh',
        recommendedCarsSelected: (selected: number, total: number) => `${selected} selected / ${total} total`,
        recommendedCarsNoSearchResults: 'No matching recommended cars.',
        recommendedCarsCreatedAt: 'Last modified (UTC+8)',
        recommendedCarsImageSrc: 'Image',
        recommendedCarsNew: 'Add car',
        recommendedCarsFormTitleNew: 'Add recommended car',
        recommendedCarsFormTitleEdit: 'Edit recommended car',
        recommendedCarsDetailTitle: 'Recommended car details',
        recommendedCarsCancel: 'Cancel',
        recommendedCarsExportEmpty: 'Select at least one recommended car before exporting.',
        recommendedCarsSaved: (count: number, path: string) => `Generated ${count} cars: ${path}`,
        recommendedCarsDbSaved: 'Recommended car saved.',
        recommendedCarsDbDeleted: 'Recommended car deleted.',
        recommendedCarsDbDeletedAll: (count: number) => `Deleted ${count} database cars.`,
        recommendedCarsTarget: 'Output file',
        recommendedCarsVersion: 'JSON version',
        recommendedCarsFileCurrent: 'Current JSON recommendations',
        recommendedCarsFileFound: (count: number, version: string) => `${count} cars in recommendedCars.json${version ? ` / ${version}` : ''}`,
        recommendedCarsFileNotFound: 'recommendedCars.json has not been generated yet.',
        recommendedCarsFileMissing: (count: number) => `${count} file IDs are not in the database.`,
        recommendedCarsInFile: 'Recommended',
        recommendedCarsAutoId: 'ID is generated automatically from name, use case, class, and PI.',
        recommendedCarsOptionalMeta: 'Weight and front weight are optional database notes and are not exported to JSON.',
        recommendedCarsEdit: 'Edit',
        recommendedCarsDuplicate: 'Duplicate',
        recommendedCarsRemove: 'Remove',
        recommendedCarsDeleteAll: 'Delete all database cars',
        recommendedCarsDeleteAllConfirm: 'Delete all recommended cars from the local database? This will not delete recommendedCars.json.',
        recommendedCarsCopied: 'Duplicated as a new car. Fill a new tuneCode before saving.',
        recommendedCarsDuplicateTuneCode: 'tuneCode already exists.',
        recommendedCarsDuplicateIdentity: 'This vehicle identity already exists. Change name, use case, class, or PI, or edit the existing record.',
        recommendedCarsTagsHint: 'Separate tags with commas, for example: grip, high speed, easy.',
        tuneHarvestTab: 'Tune Code Harvest',
        tuneHarvestTitle: 'Tune code harvest',
        tuneHarvestSubtitle: 'Collect FH6 tune share codes into a review pool before importing them into recommended cars.',
        tuneHarvestSources: 'Sources',
        tuneHarvestSourceLabels: {
            jsr_chronic_sheet: 'JSR Chronic Sheet',
            codmunity: 'CODMunity',
            forzafire: 'ForzaFire details',
        },
        tuneHarvestDryRun: 'Dry run',
        tuneHarvestDryRunHint: 'Dry run shows extracted candidates without writing them to the local database.',
        tuneHarvestLimit: 'ForzaFire detail limit',
        tuneHarvestRun: 'Run harvest',
        tuneHarvestStop: 'Stop harvest',
        tuneHarvestStopping: 'Stopping harvest...',
        tuneHarvestRefresh: 'Refresh pool',
        tuneHarvestClear: 'Clear pool',
        tuneHarvestClearConfirm: 'Delete every candidate in the review pool?',
        tuneHarvestCleared: (count: number) => `Cleared ${count} candidates from the review pool.`,
        tuneHarvestCandidates: 'Review pool',
        tuneHarvestNoCandidates: 'No harvest candidates found.',
        tuneHarvestSelectedSourcesRequired: 'Select at least one source.',
        tuneHarvestResult: (found: number, saved: number, pending: number, rejected: number) => `Found ${found}, saved ${saved}, pending ${pending}, rejected ${rejected}.`,
        tuneHarvestStopped: 'Harvest stopped.',
        tuneHarvestStatusFilter: 'Status filter',
        tuneHarvestSearch: 'Search',
        tuneHarvestSearchPlaceholder: 'Code, vehicle, tuner, source, use case...',
        tuneHarvestStatusLabels: {
            all: 'All',
            pending: 'Pending',
            rejected: 'Rejected',
            imported: 'Imported',
        },
        tuneHarvestUseCandidate: 'Use',
        tuneHarvestReject: 'Reject',
        tuneHarvestRestore: 'Restore',
        tuneHarvestSource: 'Source',
        tuneHarvestVehicle: 'Vehicle',
        tuneHarvestCode: 'Code',
        tuneHarvestContext: 'Context',
        tuneHarvestMatch: 'Match',
        tuneHarvestStatus: 'Status',
        tuneHarvestWarnings: 'Warnings',
        tuneHarvestImported: 'Candidate marked as imported.',
        tuneHarvestRejected: 'Candidate rejected.',
        tuneHarvestRestored: 'Candidate restored to pending.',
        tuneHarvestCopiedToRecommended: 'Candidate loaded into the recommended car form.',
        developerMode: 'Developer Mode',
        tireLab: 'Tire Model Lab',
        tireLabTitle: 'Four-wheel Tire Model Test',
        tireLabSubtitle: 'Experimental tire-centered model. It does not affect existing quick diagnosis, reports, or tune drafts.',
        tireLabNoPersistence: 'Tire model lab: no session, recording, database writes, or replay.',
        tireLabEmpty: 'Start Tire Model Lab and receive telemetry to inspect four-wheel grip balance.',
        tireRegressionTitle: 'Tire Regression Lab',
        tireRegressionSubtitle: 'Independent sample validation for the tire model. It does not affect quick diagnosis, reports, or tune drafts.',
        tireRegressionSaveCurrent: 'Save current Tire Lab window',
        tireRegressionSampleName: 'Sample name',
        tireRegressionScenario: 'Scenario',
        tireRegressionWindowSeconds: 'Window seconds',
        tireRegressionSave: 'Save sample',
        tireRegressionRunOne: 'Run sample',
        tireRegressionRunAll: 'Run all samples',
        tireRegressionDelete: 'Delete sample',
        tireRegressionUpdateExpected: 'Save expectation',
        tireRegressionSamples: 'Regression samples',
        tireRegressionExpected: 'Expected rule',
        tireRegressionResults: 'Run results',
        tireRegressionNoSamples: 'No regression samples yet. Save a Tire Lab window first.',
        tireRegressionNoSelection: 'Select a sample to inspect or edit expectations.',
        tireRegressionAllowedPhases: 'Allowed phases',
        tireRegressionRequiredGrip: 'Required grip limits',
        tireRegressionAllowedAxles: 'Allowed axles',
        tireRegressionForbiddenGrip: 'Forbidden grip limits',
        tireRegressionMinQuality: 'Minimum data quality',
        tireRegressionNotes: 'Notes',
        tireRegressionActual: 'Actual output',
        tireRegressionPassed: 'Passed',
        tireRegressionFailed: 'Failed',
        tireRegressionFailures: 'Failures',
        tireRegressionCsvHint: 'Comma-separated IDs, for example corner_exit, traction_limit, rear.',
        tireRegressionRequiresTireLab: 'Saving samples requires Tire Model Lab data in the current memory window.',
        tireRegressionSaved: 'Tire regression sample saved.',
        tireRegressionExpectationSaved: 'Expectation saved.',
        tireRegressionDeleted: 'Sample deleted.',
        tireRegressionRan: 'Regression run complete.',
        tireRegressionFailureLabels: {
            phase_mismatch: 'Phase did not match expectation',
            required_grip_missing: 'Required grip limit was not detected',
            limited_axle_mismatch: 'Limited axle did not match expectation',
            forbidden_grip_detected: 'Forbidden grip limit was detected',
            data_quality_below_minimum: 'Data quality is below the minimum',
        },
        modelPipelineTitle: 'Model Pipeline Lab',
        modelPipelineSubtitle: 'Experimental read-only runner for detector, decisioner, and interpreter combinations. It does not modify sessions or tune profiles.',
        modelPipelineSource: 'Data source',
        modelPipelineDetector: 'Detector',
        modelPipelineDecisioner: 'Decisioner',
        modelPipelineInterpreter: 'Interpreter',
        modelPipelineSession: 'Telemetry session',
        modelPipelineRun: 'Run pipeline',
        modelPipelineRunComplete: 'Pipeline run complete.',
        modelPipelineNoCatalog: 'Pipeline catalog is not loaded yet.',
        modelPipelineNoResult: 'Choose a source and model combination, then run the pipeline.',
        modelPipelineProblems: '问题集',
        modelPipelineDecisions: '决策结果',
        modelPipelineAdvice: '解释器建议',
        modelPipelineWarnings: 'Compatibility warnings',
        modelPipelineDocs: 'Document sources',
        modelPipelineSourceSummary: 'Source summary',
        modelPipelineExplainOnly: 'Explain-only. No numeric write values or tune-profile changes are produced.',
        modelPipelineStatus: 'Status',
        modelPipelineConfidence: 'Confidence',
        modelPipelineShouldTune: 'Should tune',
        modelPipelineEvidence: 'Evidence',
        modelPipelineStatusLabels: {
            ready: 'Ready',
            no_data: 'No data',
            incompatible: 'Incompatible',
        },
        modelPipelineAdviceCategoryLabels: {
            power_to_tire: 'Power to tire',
            platform_aero: 'Platform / aero',
            mechanical_grip: 'Mechanical grip',
            observe: 'Observe',
        },
        modelPipelineAdviceDirectionLabels: {
            reduce_wheel_torque_or_drive_lock: 'Reduce wheel torque or drive lock',
            verify_brake_balance_pressure_and_decel_lock: 'Verify brake balance, pressure, and decel lock',
            verify_platform_then_use_tiered_ride_height_or_aero: 'Verify platform, then use tiered ride height or aero',
            verify_pressure_window_and_slip_heat: 'Verify pressure window and slip heat',
            rebalance_tire_contact_and_load_transfer: 'Rebalance tire contact and load transfer',
            collect_more_evidence_before_tuning: 'Collect more evidence before tuning',
        },
        modelPipelineScopeLabels: {
            driven_wheels: 'Driven wheels',
            front_rear_balance: 'Front / rear balance',
            front_rear_platform: 'Front / rear platform',
            front_rear_tires: 'Front / rear tires',
            vehicle: 'Vehicle',
        },
        modelPipelineRationaleLabels: {
            tire_problem_group_decision: 'Decision generated from tire issue groups',
            legacy_problem_fallback: 'Legacy problem fallback decision',
            legacy_road_decision_v1: 'Legacy road decision model',
            tire_lab_problem_groups_v1: 'Tire Lab issue groups',
            docs_v12_interpreter_outputs_no_write_values: 'Docs v1.2 interpreter outputs directions only, not writable numeric values',
            'Docs v1.2 treats gearing and differential as Forza slider/display levers; verify whether torque delivery exceeds tire traction before changing chassis balance.': 'Docs v1.2 treats gearing and differential as Forza slider/display levers; verify whether torque delivery exceeds tire traction before changing chassis balance.',
            'Docs v1.2 uses conservative road brake baselines by drivetrain; explain braking issues through balance, pressure, and decel lock before changing unrelated systems.': 'Docs v1.2 uses conservative road brake baselines by drivetrain; explain braking issues through balance, pressure, and decel lock before changing unrelated systems.',
            'Docs v1.2 keeps ride height and aero as low/medium/high tier explanations, not precise write values; use them only after verifying platform evidence.': 'Docs v1.2 keeps ride height and aero as low/medium/high tier explanations, not precise write values; use them only after verifying platform evidence.',
            'Docs v1.4 uses BAR as the primary tire-pressure unit; thermal issues need pressure, slip, and temperature evidence before any 0.02-0.03 BAR correction is trusted.': 'Docs v1.4 uses BAR as the primary tire-pressure unit; thermal issues need pressure, slip, and temperature evidence before any 0.02-0.03 BAR correction is trusted.',
            'Docs v1.2 makes tire pressure, alignment, anti-roll bars, springs, and damping explicit Forza display/sliders; use them as grouped levers around tire contact and load transfer.': 'Docs v1.2 makes tire pressure, alignment, anti-roll bars, springs, and damping explicit Forza display/sliders; use them as grouped levers around tire contact and load transfer.',
            'Docs v1.2/v1.4 make tire pressure, alignment, anti-roll bars, springs, and damping explicit Forza display/sliders; use BAR-first tire pressure only when slip and heat evidence support it.': 'Docs v1.2/v1.4 make tire pressure, alignment, anti-roll bars, springs, and damping explicit Forza display/sliders; use BAR-first tire pressure only when slip and heat evidence support it.',
            'Docs v1.2 interpreter did not find a specific road baseline lever for this decision; keep it as an observation.': 'Docs v1.2 interpreter did not find a specific road baseline lever for this decision; keep it as an observation.',
        },
        modelPipelineEvidenceLabels: {
            driven_wheel_slip_ratio: 'Driven wheel slip ratio',
            throttle: 'Throttle',
            rpm: 'RPM',
            speed_gain: 'Speed gain',
            gear: 'Gear',
            brake: 'Brake',
            deceleration_g: 'Deceleration G',
            front_slip_ratio: 'Front slip ratio',
            rear_slip_ratio: 'Rear slip ratio',
            speed_band: 'Speed band',
            suspension_offset: 'Suspension offset',
            combined_slip: 'Combined slip',
            g_force: 'G force',
            front_tire_temp: 'Front tire temp',
            rear_tire_temp: 'Rear tire temp',
            front_combined_slip: 'Front combined slip',
            rear_combined_slip: 'Rear combined slip',
            slip_angle: 'Slip angle',
            steer: 'Steer',
            sample_quality: 'Sample quality',
            phase: 'Phase',
            speed: 'Speed',
            inputs: 'Inputs',
        },
        tuneGeneratorTitle: 'Quick Tune',
        tuneGeneratorSubtitle: 'Generate a neutral Road baseline from minimum vehicle inputs. No AI, telemetry, or saved session is required.',
        tuneGeneratorValueHint: 'Enter the vehicle information shown in Forza Horizon 6.',
        tuneGeneratorNoHistory: 'Generated baselines are only saved when you create or apply a tune profile.',
        remoteTuneTitle: 'Remote Tune',
        remoteTuneSubtitle: 'Start a LAN page for iPhone/iPad quick tuning preview. It does not copy results or save tune profiles.',
        remoteTuneStatus: 'Status',
        remoteTunePort: 'Web port',
        remoteTuneLanAddress: 'LAN address',
        remoteTuneUrl: 'Access URL',
        remoteTuneStart: 'Start remote page',
        remoteTuneStop: 'Stop remote page',
        remoteTuneRunning: 'Remote tune page is running.',
        remoteTuneStopped: 'Remote tune page is stopped.',
        remoteTuneReadOnly: 'MVP is preview-only: no copy result, no profile saving, no database writes.',
        remoteTuneDeviceHint: 'Open this URL from Safari on an iPhone/iPad in the same LAN.',
        remoteTunePortInvalid: 'Port must be an integer from 1 to 65535.',
        quickTuneInputButton: 'Input vehicle info',
        quickTuneInputTitle: 'Vehicle information',
        quickTuneLastSummary: 'Last quick tune',
        quickTuneNoResult: 'No quick tune generated yet',
        quickTuneUseCase: 'Use case',
        quickTuneTireCompound: 'Tire compound',
        quickTuneUnsupportedUseCase: 'This MVP supports Road, Drift, Rally, Offroad, and Drag static baseline generation.',
        quickTuneDriftRwdPreferred: 'Drift quick tune is RWD-first. FWD/AWD will generate core setup values, but differential values are skipped.',
        quickTuneValidationSummary: 'Fix the highlighted vehicle information before generating.',
        quickTuneIntegerRange: (label: string, min: number, max: number) => `${label} must be an integer from ${min} to ${max}.`,
        quickTunePositiveInteger: (label: string) => `${label} must be a positive integer.`,
        quickTuneTireSizeInvalid: 'Tire size must use integer width, aspect ratio, and rim values.',
        quickTuneGearingToggle: 'Enable static gearing calculation',
        quickTuneDriftGearingHint: 'Drift gearing targets high RPM in the core drift gear at the selected speed.',
        quickTuneDragGearingHint: 'Drag gearing targets strong acceleration through the selected terminal speed.',
        quickTuneTargetDriftSpeed: 'Target drift speed',
        quickTuneTargetDragSpeed: 'Target terminal speed',
        quickTuneCarryToProfessional: 'Use in Professional Tuning',
        quickTuneCarriedToProfessional: 'Quick tune parameters loaded into Professional Tuning. Save there only when you are ready.',
        quickTuneBiasTitle: 'Baseline bias sliders',
        quickTuneBiasHint: 'Start from the neutral Road baseline, then apply a conservative whole-car bias before saving or carrying to Professional Tuning.',
        quickTuneBalance: 'Balance',
        quickTuneBalanceLeft: 'Stable',
        quickTuneBalanceRight: 'Agile',
        quickTuneStiffness: 'Stiffness',
        quickTuneStiffnessLeft: 'Soft',
        quickTuneStiffnessRight: 'Hard',
        quickTuneSpeed: 'Speed',
        quickTuneSpeedLeft: 'Top speed',
        quickTuneSpeedRight: 'Acceleration',
        quickTuneNeutral: 'Neutral',
        quickTuneSpeedDisabled: 'Enable static gearing and generate gear ratios to adjust speed bias.',
        quickTuneDrivetrainLabels: {
            FWD: 'FWD',
            AWD: 'AWD',
            RWD: 'RWD',
        },
        quickTuneTireCompoundLabels: {
            stock: 'Stock',
            street: 'Street',
            sport: 'Sport',
            semi: 'Semi-slick',
            slick: 'Slick',
            rally: 'Rally',
            offroad: 'Offroad',
            drift: 'Drift',
            drag: 'Drag',
            snow: 'Snow',
        },
        tuneGeneratorMinimum: 'Minimum inputs',
        tuneGeneratorAdvanced: 'Advanced optional inputs',
        tuneGeneratorPreview: 'Baseline preview',
        tuneGeneratorFields: 'Generated fields',
        tuneGeneratorSkipped: 'Not generated in MVP',
        tuneGeneratorNextTest: 'Next test plan',
        tuneGeneratorGenerate: 'Generate preview',
        tuneGeneratorCreate: 'Save as professional profile',
        tuneGeneratorApply: 'Apply to existing professional profile',
        tuneGeneratorReset: 'Reset inputs',
        tuneGeneratorOpenExpert: 'Open in Professional Tuning',
        tuneGeneratorTarget: 'Target tune profile',
        tuneGeneratorNoPreview: 'Enter the minimum inputs and generate a Road baseline preview.',
        tuneGeneratorNoTarget: 'Select a target tune profile before applying.',
        tuneGeneratorNoSelection: 'Select at least one generated field.',
        tuneGeneratorCreated: 'Baseline tune profile created.',
        tuneGeneratorApplied: 'Baseline applied to the selected tune profile.',
        tuneGeneratorConfidence: 'Confidence',
        tuneGeneratorSelectedCount: 'Selected fields',
        tuneGeneratorRangeHint: 'Leave range fields empty when the car-specific slider range is unknown.',
        tuneGeneratorCarName: 'Car name',
        tuneGeneratorVersionName: 'Version',
        tuneGeneratorCarOrdinal: 'Car ID',
        tuneGeneratorCarCategory: 'Car category',
        tuneGeneratorPI: 'PI',
        tuneGeneratorDrivetrain: 'Drivetrain',
        tuneGeneratorWeight: 'Weight',
        tuneGeneratorFrontWeight: 'Front weight',
        tuneGeneratorPower: 'Power',
        tuneGeneratorTorque: 'Torque',
        tuneGeneratorGearingInputs: 'Static gearing',
        tuneGeneratorGearingHint: 'Optional. Generates final drive and 1-N gears from redline, gear count, tire size, and target top speed.',
        tuneGeneratorRedlineRPM: 'Redline RPM',
        tuneGeneratorGearCount: 'Gear count',
        tuneGeneratorTireDiameter: 'Tire size',
        tuneGeneratorTargetTopSpeed: 'Target top speed',
        tuneGeneratorRideRange: 'Ride height range',
        tuneGeneratorAeroRange: 'Aero range',
        tuneGeneratorAdjustableHint: 'Tier recommendations do not write exact values. Set the matching slider manually in game.',
        tuneGeneratorTierRecommendations: 'Tier recommendations',
        tuneGeneratorTierManual: 'Manual tier only; not applied to the tune profile.',
        tuneGeneratorFrontRideAdjustable: 'Front ride height adjustable',
        tuneGeneratorRearRideAdjustable: 'Rear ride height adjustable',
        tuneGeneratorFrontAeroAdjustable: 'Front aero adjustable',
        tuneGeneratorRearAeroAdjustable: 'Rear aero adjustable',
        tuneGeneratorTierLabels: {
            low: 'Low',
            medium: 'Medium',
            high: 'High',
        },
        tuneGeneratorMin: 'Min',
        tuneGeneratorMax: 'Max',
        tuneGeneratorReason: 'Reason',
        tireLabMatrix: 'Four-wheel grip matrix',
        tireLabAxleBalance: 'Front vs rear grip balance',
        tireLabFrontAxle: 'Front axle',
        tireLabRearAxle: 'Rear axle',
        tireLabLimitType: 'Tire limit state',
        tireLabPhase: 'Phase',
        tireLabConfidence: 'Model confidence',
        tireLabExplanation: 'Model explanation',
        tireLabHints: 'Interpretation directions',
        tireLabWarnings: 'Risks / warnings',
        tireLabLeftRight: 'Left / right balance',
        tireLabWindow: 'Window',
        tireDataQuality: 'Data quality',
        tireGripLimit: 'Grip limit type',
        tireGripRelationships: 'Four-tire relationships',
        tirePhaseStability: 'Phase stability',
        tireStablePhase: 'Stable phase',
        tireScoreMargin: 'Score margin',
        tireDynamicSamples: 'Dynamic samples',
        tireSignalQuality: 'Signal quality',
        tireLimitedAxle: 'Limited axle',
        tireLimitedWheels: 'Limited wheels',
        tirePrimaryEvidence: 'Primary evidence',
        tireReason: 'Reason',
        tirePhaseCurrent: 'Current phase',
        tirePhaseSecondary: 'Secondary phase',
        tirePhaseEvidence: 'Phase evidence',
        tirePhaseScores: 'Phase scores',
        tirePhaseSpeedDelta: 'Speed delta',
        tirePhaseSpeedReference: 'Speed reference',
        tirePhaseSpeedBand: 'Speed band',
        tirePhaseThrottleDelta: 'Throttle delta',
        tirePhaseSteerDelta: 'Steer delta',
        tirePhaseBrake: 'Brake avg / peak',
        tirePhaseHandbrake: 'Handbrake avg / peak',
        tirePhasePlaneG: 'Plane G avg / peak',
        tirePhaseDecelG: 'Decel G avg / peak',
        tirePhaseAccelG: 'Accel G avg / peak',
        tireIssueGroups: 'Tire issue groups',
        tireIssueSegments: 'Issue segments',
        tireIssueNoGroups: 'No aggregated tire issues in this window.',
        tireIssueNoSegments: 'No confirmed issue segments.',
        tireIssueCount: 'Occurrences',
        tireIssueDuration: 'Total duration',
        tireIssueSpeedRange: 'Speed range',
        tireIssueOperations: 'Operation tags',
        tireIssueDriftSource: 'Drift source',
        tireIssueRisk: 'Risk',
        tireIssueEvidence: 'Representative evidence',
        tireIssueAdvice: 'Repair directions',
        tireIssueExperimentHint: 'Experimental explanation only; it will not write tune profiles or affect official reports.',
        tireIssueNoAdvice: 'No repair direction can be produced from the current tire issue groups.',
        tireIssuePriorityAdvice: 'Priority repair directions',
        tireIssueGroupAdvice: 'Directions by issue group',
        tireIssuePrimaryCause: 'Primary cause hypothesis',
        tireIssueShouldTune: 'Tune direction',
        tireIssueNoTune: 'No tuning action',
        tireIssueRelatedFields: 'Related tune levers',
        tireIssueVerifyEvidence: 'Evidence to verify',
        tireIssueConflict: 'Conflict note',
        tireIssueMissingInputs: 'Missing conditions',
        tireAdviceLayerLabels: {
            primary: 'Primary',
            alternative: 'Alternative',
            check: 'Check',
            observe: 'Observe',
        },
        tireAdviceCategoryLabels: {
            tire_pressure: 'Tire pressure',
            alignment: 'Alignment',
            antiroll: 'Anti-roll bars',
            spring_damping: 'Springs / damping',
            aero_platform: 'Aero / platform',
            brake: 'Brake',
            differential: 'Differential',
            gearing: 'Gearing',
            driver_input: 'Driver input',
            data_quality: 'Data quality',
        },
        tireAdviceDirectionLabels: {
            increase_front_high_speed_support: 'Increase front high-speed support',
            increase_front_mechanical_grip: 'Increase front mechanical grip',
            check_front_contact_patch: 'Check front contact patch',
            increase_rear_stability: 'Increase rear stability',
            check_rear_contact_patch: 'Check rear contact patch',
            reduce_four_wheel_lateral_load: 'Reduce four-wheel lateral load',
            reduce_overlap_input: 'Reduce overlapping inputs',
            reduce_drive_lock: 'Reduce drive lock',
            reduce_wheel_torque: 'Reduce wheel torque',
            move_brake_balance_rearward: 'Move brake balance rearward',
            move_brake_balance_forward: 'Move brake balance forward',
            check_decel_lock: 'Check decel differential',
            check_front_brake_platform: 'Check front brake platform',
            check_brake_balance: 'Check brake balance',
            check_platform: 'Check platform',
            check_temperature_window: 'Check temperature window',
            check_left_right: 'Check left/right evidence',
            continue_sampling: 'Continue sampling',
            avoid_tuning: 'Do not tune for this behavior',
        },
        tireAdviceCauseLabels: {
            data_not_reliable: 'Data is not reliable enough',
            driver_handbrake_drift: 'Driver-initiated handbrake drift',
            driver_weight_transfer_drift: 'Driver-induced weight transfer drift',
            front_high_speed_lateral_limit: 'Front tires reach high-speed lateral limit',
            front_mechanical_lateral_limit: 'Front mechanical grip limit',
            rear_lateral_stability_limit: 'Rear lateral stability limit',
            four_wheel_lateral_limit: 'Four tires are laterally loaded',
            drive_torque_exceeds_tire_grip: 'Drive torque exceeds tire grip',
            driven_wheel_longitudinal_slip: 'Driven wheels show longitudinal slip',
            rear_brake_or_decel_instability: 'Rear brake/deceleration instability',
            front_brake_overload: 'Front brake overload',
            combined_longitudinal_lateral_overload: 'Combined longitudinal and lateral overload',
            platform_travel_or_load_risk: 'Platform travel/load risk',
            tire_temperature_risk: 'Tire temperature risk',
            left_right_signal_or_load_imbalance: 'Left/right signal or load imbalance',
            unknown_tire_issue: 'Unknown tire issue',
        },
        tireAdviceRationaleLabels: {
            front_high_speed_lateral_limit_prioritize_platform: 'At high speed, verify front aero/platform support before changing low-speed mechanical balance.',
            front_high_speed_lateral_limit_secondary_mechanical_grip: 'If platform evidence is clean, verify front mechanical grip and contact patch.',
            front_lateral_limit_prioritize_mechanical_grip: 'Front lateral slip dominates at this speed, so start from front mechanical grip levers.',
            front_lateral_limit_verify_alignment: 'Confirm front contact patch before treating it as only an anti-roll issue.',
            rear_lateral_limit_prioritize_stability: 'Rear lateral slip dominates; prioritize rear stability and platform response.',
            rear_lateral_limit_verify_alignment: 'Confirm rear contact patch and toe/camber behavior.',
            four_wheel_lateral_limit_reduce_platform_load: 'Both axles are near lateral load; reduce platform load before single-axle changes.',
            four_wheel_lateral_limit_verify_driver_overlap: 'Verify steering/throttle/brake overlap before tuning around the behavior.',
            front_traction_limit_reduce_accel_diff: 'Front driven tires are slipping under power; check front acceleration differential first.',
            rear_traction_limit_reduce_accel_diff: 'Rear driven tires are slipping under power; check rear acceleration differential first.',
            driven_traction_limit_balance_diff: 'Driven wheels are traction limited; balance differential locking across the driven axle(s).',
            front_traction_limit_check_gearing: 'If differential changes do not solve it, verify wheel torque from gearing.',
            rear_traction_limit_check_gearing: 'If differential changes do not solve it, verify wheel torque from gearing.',
            driven_traction_limit_check_gearing: 'Verify final drive and low gears are not overloading the driven tires.',
            rear_braking_limit_move_balance_forward: 'Rear tires lock or slide under braking; move brake load forward or reduce pressure.',
            rear_braking_limit_check_decel_diff: 'Rear decel differential can add instability during braking.',
            front_braking_limit_move_balance_rearward: 'Front tires overload under braking; reduce front brake load or pressure.',
            front_braking_limit_check_platform: 'Front platform support can make braking load transfer abrupt.',
            combined_limit_trail_brake_reduce_overlap: 'Brake plus steering is overloading the tire set.',
            combined_limit_verify_brake_balance: 'If overlap is intentional, verify brake balance and pressure.',
            combined_limit_corner_exit_reduce_overlap: 'Throttle plus steering is overloading the tire set.',
            combined_limit_verify_drive_lock: 'If exit overlap is intentional, verify drive lock and torque split.',
            combined_limit_reduce_combined_load: 'The issue is combined load rather than one simple axle limit.',
            combined_limit_verify_platform: 'Confirm platform stability before changing a single tire lever.',
            platform_risk_check_travel_and_damping: 'Suspension travel/platform risk should be verified separately from tire slip.',
            thermal_risk_check_pressure_and_slip: 'Tire temperature is a risk signal; verify pressure and sustained slip.',
            left_right_imbalance_verify_route_and_sensor: 'Left/right differences are evidence only; verify route direction and sensor consistency.',
            data_quality_continue_sampling: 'Collect more dynamic samples before making a tuning decision.',
            handbrake_drift_driver_behavior: 'Handbrake-initiated drift is a driver behavior, not a repair target.',
            scandinavian_flick_driver_behavior: 'A flick/weight-transfer drift is driver behavior unless it repeats unintentionally.',
            low_confidence_verify_before_tuning: 'The signal is low confidence; verify the pattern before tuning.',
            unknown_issue_continue_sampling: 'The issue is not classified strongly enough; continue sampling.',
        },
        tireDataQualityLabels: {
            valid: 'Valid',
            low_confidence: 'Low confidence',
            invalid: 'Invalid',
            ok: 'OK',
            low: 'Low',
            flat: 'Flat',
        },
        tireDataQualityReasonLabels: {
            tire_data_no_samples: 'No tire model samples yet.',
            tire_data_menu_or_no_vehicle: 'Menu/no vehicle telemetry; tire limit cannot be judged.',
            tire_data_sample_insufficient: 'Sample count is too low.',
            tire_data_dynamic_sample_insufficient: 'Dynamic tire-load samples are insufficient.',
            tire_data_stationary: 'Vehicle is stationary.',
            tire_data_speed_low: 'Speed signal is too low for limit judgment.',
            tire_data_g_force_flat: 'G-force signal is nearly flat.',
            tire_data_slip_signal_flat: 'Wheel slip signal is flat or missing.',
            tire_data_input_low: 'Driver input signal is too low.',
        },
        tireGripLimitLabels: {
            lateral_limit: 'Lateral grip limit',
            traction_limit: 'Traction limit',
            braking_limit: 'Braking grip limit',
            combined_limit: 'Combined four-tire limit',
            balanced_near_limit: 'Balanced near limit',
            no_limit_detected: 'No tire limit detected',
            risk: 'Risk',
            data_quality: 'Data quality',
        },
        tireIssueTypeLabels: {
            lateral_limit: 'Lateral limit',
            traction_limit: 'Traction limit',
            braking_limit: 'Braking limit',
            combined_limit: 'Combined limit',
            platform_risk: 'Platform risk',
            thermal_risk: 'Thermal risk',
            left_right_imbalance: 'Left/right imbalance',
            data_insufficient: 'Data insufficient',
        },
        tireOperationTagLabels: {
            throttle_on: 'Throttle on',
            throttle_steady: 'Throttle steady',
            throttle_lift: 'Throttle lift',
            light_brake: 'Light brake',
            heavy_brake: 'Heavy brake',
            handbrake_active: 'Handbrake active',
            steer_increasing: 'Steer increasing',
            steer_holding: 'Steer holding',
            steer_unwinding: 'Steer unwinding',
            speed_rising: 'Speed rising',
            speed_falling: 'Speed falling',
        },
        tireDriftSourceLabels: {
            handbrake_initiated: 'Handbrake initiated',
            power_oversteer: 'Power oversteer',
            scandinavian_flick: 'Scandinavian flick',
            lift_off_oversteer: 'Lift-off oversteer',
            unknown_oversteer: 'Unknown oversteer',
        },
        tireGripReasonLabels: {
            tire_grip_no_dynamic_limit: 'Dynamic slip evidence does not show a dominant tire limit.',
            tire_grip_stationary: 'Stationary state; tire limit is not evaluated.',
            tire_grip_no_dynamic_load: 'No dynamic tire-load samples.',
            tire_grip_data_invalid: 'Data quality is invalid.',
            tire_grip_lateral_slip: 'Slip-angle evidence points to a lateral grip limit.',
            tire_grip_power_slip: 'Driven tire slip points to a traction limit.',
            tire_grip_braking_slip: 'Brake slip evidence points to a braking grip limit.',
            tire_grip_handbrake_slip: 'Handbrake input is causing rear slip.',
            tire_grip_four_wheel_combined: 'Front and rear tires are both near the combined slip limit.',
            tire_grip_near_limit: 'One axle is near the combined slip limit.',
        },
        tireAxleLabels: {
            none: 'None',
            front: 'Front axle',
            rear: 'Rear axle',
            both: 'Both axles',
            driven: 'Driven wheels',
        },
        tirePhaseStabilityLabels: {
            stable: 'Stable',
            transition: 'Transition',
            low_confidence: 'Low confidence',
        },
        powerToTireTitle: 'Power output → tire traction',
        powerToTireSubtitle: 'Experimental DO-only power delivery check; tune profile power inputs are not required.',
        powerToTireStatus: 'Power landing state',
        powerToTireDrivenAxle: 'Driven axle',
        powerToTirePower: 'Power',
        powerToTireTorque: 'Torque',
        powerToTireRPM: 'RPM / ratio',
        powerToTireGear: 'Gear',
        powerToTireThrottle: 'Throttle',
        powerToTireDrivenSlip: 'Driven slip P90',
        powerToTireAccel: 'Acceleration',
        powerToTireSamples: 'High-throttle samples',
        powerToTireSignal: 'Power signal',
        powerToTireAvailable: 'Available',
        powerToTireUnavailable: 'Unavailable',
        powerToTireSummaryLabels: {
            power_to_tire_no_data: 'No power-to-tire data',
            power_to_tire_insufficient: 'Insufficient high-throttle samples',
            power_to_tire_low_throttle: 'Low throttle',
            power_to_tire_power_signal_unavailable: 'Power / torque signal unavailable',
            traction_over_power: 'Power exceeds driven-tire traction',
            rpm_below_useful_range: 'RPM below useful range',
            rpm_too_high_or_gear_short: 'RPM too high / gear too short',
            power_not_reaching_ground: 'Power is not reaching the ground',
            power_landing_ok: 'Power delivery looks usable',
        },
        powerToTireExplanationLabels: {
            power_to_tire_waiting_for_samples: 'Waiting for enough DO samples with speed, RPM, gear, power and tire slip.',
            power_to_tire_need_high_throttle: 'Collect more high-throttle samples before judging power delivery.',
            power_to_tire_needs_more_high_throttle: 'Collect more high-throttle samples before judging power delivery.',
            power_to_tire_low_throttle_explanation: 'Throttle is too low for a power-to-traction conclusion.',
            power_to_tire_power_signal_unavailable_explanation: 'The DO power or torque signal is zero/abnormal, so the result is shown as low confidence.',
            power_to_tire_traction_over_power_explanation: 'Driven tire slip is high while acceleration is weak or unstable, so torque is exceeding usable tire traction.',
            traction_over_power_explanation: 'Driven tire slip is high while acceleration is weak or unstable, so torque is exceeding usable tire traction.',
            power_to_tire_rpm_below_explanation: 'High throttle, low driven slip and weak acceleration point to RPM below the useful range or a gear that is too long.',
            rpm_below_useful_range_explanation: 'High throttle, low driven slip and weak acceleration point to RPM below the useful range or a gear that is too long.',
            power_to_tire_rpm_high_explanation: 'RPM stays near the upper range while acceleration fades, so the current gear may be too short or shift timing is too late.',
            rpm_too_high_or_gear_short_explanation: 'RPM stays near the upper range while acceleration fades, so the current gear may be too short or shift timing is too late.',
            power_to_tire_not_reaching_ground_explanation: 'Power is being requested, but speed and G response are weak; check traction, surface and power delivery before gearing.',
            power_not_reaching_ground_explanation: 'Power is being requested, but speed and G response are weak; check traction, surface and power delivery before gearing.',
            power_to_tire_ok_explanation: 'Power, RPM and driven tire slip are not showing a dominant power-delivery bottleneck in this window.',
            power_landing_ok_explanation: 'Power, RPM and driven tire slip are not showing a dominant power-delivery bottleneck in this window.',
        },
        powerToTireDrivenAxleLabels: {
            front: 'Front driven wheels',
            rear: 'Rear driven wheels',
            all: 'All driven wheels',
            unknown: 'Unknown driven wheels',
        },
        brakeToTireTitle: 'Brake input → tire grip',
        brakeToTireSubtitle: 'Experimental DO-only braking grip check; speed delta is primary, raw G is supporting evidence.',
        brakeToTireStatus: 'Braking grip state',
        brakeToTireBrake: 'Brake',
        brakeToTireHandbrake: 'Handbrake',
        brakeToTireSpeed: 'Speed / delta',
        brakeToTireSteer: 'Steer',
        brakeToTireDecel: 'Decel G',
        brakeToTirePlaneG: 'Raw X/Z G',
        brakeToTireFrontSlip: 'Front brake slip P90',
        brakeToTireRearSlip: 'Rear brake slip P90',
        brakeToTireFrontCombined: 'Front combined P90',
        brakeToTireRearCombined: 'Rear combined P90',
        brakeToTireSamples: 'Brake samples',
        brakeToTireTrail: 'Trail braking',
        brakeToTireHandbrakeActive: 'Handbrake active',
        brakeToTireSummaryLabels: {
            brake_to_tire_no_data: 'No brake-to-tire data',
            brake_to_tire_insufficient: 'Insufficient brake samples',
            brake_to_tire_low_brake: 'Low brake input',
            front_brake_lock_tendency: 'Front lock tendency',
            rear_brake_lock_tendency: 'Rear lock tendency',
            trail_brake_front_overload: 'Trail braking front overload',
            handbrake_rear_slide: 'Handbrake rear slide',
            brake_not_slowing_effectively: 'Brake input is not slowing effectively',
            brake_landing_ok: 'Braking grip looks usable',
        },
        brakeToTireExplanationLabels: {
            brake_to_tire_waiting_for_samples: 'Waiting for DO samples with brake input, speed, wheel slip and G evidence.',
            brake_to_tire_need_brake_samples: 'Collect more braking samples before judging brake grip.',
            brake_to_tire_low_brake_explanation: 'Brake and handbrake inputs are too low for a brake-to-tire conclusion.',
            brake_to_tire_front_lock_explanation: 'Front wheel slip ratio is dominant under braking, so the front tires are consuming too much longitudinal grip.',
            brake_to_tire_rear_lock_explanation: 'Rear wheel slip ratio is dominant under braking, so the rear tires are losing braking stability.',
            brake_to_tire_trail_brake_front_overload_explanation: 'Brake and steering overlap while front combined slip is high, pointing to front overload on entry.',
            brake_to_tire_handbrake_rear_slide_explanation: 'Handbrake input is active while rear slip rises; treat this as driver-induced rear rotation evidence.',
            brake_to_tire_not_slowing_effectively_explanation: 'Brake input is high, but speed-based deceleration is weak and tires are not clearly locked.',
            brake_to_tire_ok_explanation: 'Brake input, deceleration and front/rear slip do not show a dominant braking bottleneck in this window.',
        },
        tuneInfluenceMap: 'Tune-to-tire influence map',
        tuneInfluenceMapHint: 'Read-only explanation of how tune fields affect the four tires.',
        tuneInfluenceNoData: 'Influence map is not loaded yet.',
        tuneInfluenceButton: 'Influence',
        tuneInfluenceModalTitle: 'Tune field influence',
        tuneInfluenceType: 'Influence type',
        tuneInfluenceScope: 'Scope',
        tuneInfluencePhase: 'Phase',
        tuneInfluenceMetrics: 'Tire metrics',
        tuneInfluenceEvidence: 'Telemetry evidence',
        tuneInfluenceSideEffects: 'Common side effects',
        tuneInfluenceConditions: 'Conditions',
        tuneInfluenceTypeLabels: {
            direct: 'Direct',
            indirect: 'Indirect',
        },
        tuneInfluenceCategoryLabels: {
            tire: 'Tires',
            gearing: 'Gearing',
            alignment: 'Alignment',
            antiroll: 'Anti-roll bars',
            springs: 'Springs / ride height',
            damping: 'Damping',
            aero: 'Aero',
            brake: 'Brakes',
            differential: 'Differential',
        },
        tuneInfluenceScopeLabels: {
            front_axle: 'Front axle',
            rear_axle: 'Rear axle',
            all_wheels: 'Four wheels',
            driven_wheels: 'Driven wheels',
            left_right_balance: 'Left/right evidence',
        },
        tuneInfluencePhaseLabels: {
            all_phases: 'All phases',
            launch: 'Launch',
            braking: 'Braking',
            corner_entry: 'Corner entry',
            sustained_cornering: 'Sustained cornering',
            corner_exit: 'Corner exit',
            straight_power: 'Straight power',
            high_speed_corner: 'High-speed corner',
            transition: 'Transition',
            kerb_impact: 'Kerb / bump',
            coast: 'Coast',
        },
        tuneInfluenceMetricLabels: {
            tire_temp: 'Tire temperature',
            combined_slip: 'Combined slip',
            slip_angle: 'Slip angle',
            slip_ratio: 'Slip ratio',
            suspension_offset: 'Suspension offset',
            g_force: 'G force',
            yaw_response: 'Yaw response',
            speed_rpm: 'Speed / RPM',
            wheel_torque: 'Wheel torque',
            brake_slip: 'Brake slip',
        },
        tuneInfluenceSideEffectLabels: {
            can_increase_tire_temp: 'May raise tire temperature',
            can_reduce_stability: 'May reduce stability',
            can_affect_acceleration: 'Affects acceleration',
            can_affect_top_speed: 'Affects top speed',
            can_mask_camber_issue: 'Can mask camber evidence',
            can_reduce_opposite_axle_grip: 'Can reduce opposite-axle grip',
            can_create_understeer: 'Can create understeer',
            can_create_oversteer: 'Can create oversteer',
            can_increase_lockup_risk: 'Can increase lockup risk',
        },
        tuneInfluenceConditionLabels: {
            high_load_only: 'Use loaded samples',
            throttle_sensitive: 'Throttle-sensitive',
            drivetrain_sensitive: 'Drivetrain-sensitive',
            unlocked_only: 'Requires unlocked field',
            speed_sensitive: 'Speed-sensitive',
            aero_speed_sensitive: 'Aero speed-sensitive',
            brake_sensitive: 'Brake-sensitive',
        },
        gForceDiagnostics: 'G-force diagnostics',
        gForceCurrent: 'Current G',
        gForceAverage: 'Average G',
        gForcePeak: 'Peak G',
        gForceTotal: 'Total G',
        gForcePlane: 'X/Z plane G',
        gForceCircleScale: 'Outer circle',
        gForceDominantAxis: 'Dominant axis',
        gForceAxisMapping: 'Axis mapping',
        gForceAxisMappingLabels: {
            raw_packet_axes_unverified: 'Raw AccelerationX/Y/Z axes; circle uses X/Z with X reversed',
        },
        gForceChart: 'Live G-force graph',
        camberInference: 'Camber inference',
        camberFront: 'Front camber',
        camberRear: 'Rear camber',
        camberCorneringSamples: 'Cornering samples',
        camberStateLabels: {
            unknown: 'Unknown',
            stable: 'Stable',
            monitor: 'Monitor',
            likely_needs_more_negative: 'Likely needs more negative camber',
            platform_limited: 'Platform first',
            thermal_limited: 'Temperature first',
        },
        camberSummaryLabels: {
            camber_inference_insufficient: 'Not enough sustained-cornering samples.',
            camber_inference_front_needs_more_negative: 'Front axle may need more negative camber.',
            camber_inference_rear_needs_more_negative: 'Rear axle may need more negative camber.',
            camber_inference_both_axles_need_more_negative: 'Both axles may need more negative camber.',
            camber_inference_platform_first: 'Platform/suspension issue should be solved before camber.',
            camber_inference_temperature_first: 'Tire temperature issue should be solved before camber.',
            camber_inference_monitor: 'Camber is not the dominant evidence in this window.',
        },
        camberExplanationLabels: {
            camber_inference_needs_cornering: 'Collect sustained-cornering samples with meaningful steering before judging camber.',
            camber_inference_slip_angle_explanation: 'Because FH Data Out has one tire temperature per tire, this uses slip angle, combined slip and G-load as low-confidence camber evidence.',
            camber_inference_platform_explanation: 'Suspension/platform limits can imitate camber problems, so platform stability should be checked first.',
            camber_inference_temperature_explanation: 'Tire temperature is already limiting grip, so pressure/thermal balance should be checked before camber.',
            camber_inference_monitor_explanation: 'Slip angle and combined slip do not strongly point to camber in the current window.',
        },
        tireLabLimitLabels: {
            unknown: 'Unknown',
            no_dynamic_load: 'No dynamic load',
            stationary: 'Stationary',
            balanced: 'Balanced',
            balanced_near_limit: 'Balanced near limit',
            front_limited: 'Front tires limited',
            rear_limited: 'Rear tires limited',
            four_wheel_limited: 'Four tires limited',
            traction_limited: 'Driven tires traction-limited',
            thermal_limited: 'Tire temperature limited',
            platform_limited: 'Platform / suspension limited',
        },
        tireLabPhaseLabels: {
            unknown: 'Unknown',
            stationary: 'Stationary',
            handbrake: 'Handbrake',
            launch: 'Launch',
            straight_decel: 'Straight decel',
            braking: 'Braking',
            light_braking: 'Light braking',
            corner_entry: 'Corner entry',
            low_speed_corner: 'Low-speed corner',
            mid_speed_corner: 'Mid-speed corner',
            sustained_cornering: 'Sustained cornering',
            corner_exit: 'Corner exit',
            straight_power: 'Straight power',
            high_speed_corner: 'High-speed corner',
            drift: 'Drift',
        },
        tireLabGripStateLabels: {
            unknown: 'Unknown',
            stable: 'Stable',
            warning: 'Warning',
            limit: 'At limit',
        },
        tireLabSummaryLabels: {
            tire_model_no_data: 'No telemetry samples yet.',
            tire_model_stationary: 'Vehicle is stationary; tire limit analysis is paused.',
            tire_model_no_dynamic_load: 'No dynamic tire-load samples in the current window.',
            tire_model_balanced: 'The four tires are not showing a dominant grip limit.',
            tire_model_balanced_near_limit: 'The car is near the grip limit, but the load is reasonably balanced.',
            tire_model_front_limited: 'The front axle reaches the lateral grip limit first.',
            tire_model_front_brake_limited: 'The front axle is overloaded while braking into the corner.',
            tire_model_rear_limited: 'The rear axle reaches the lateral grip limit first.',
            tire_model_rear_power_limited: 'The rear axle is overloaded during power application.',
            tire_model_rear_handbrake_limited: 'The rear axle is sliding under handbrake input.',
            tire_model_rear_traction_limited: 'Rear driven tires show excessive longitudinal slip.',
            tire_model_front_traction_limited: 'Front driven tires show excessive longitudinal slip.',
            tire_model_drive_traction_limited: 'Driven tires show excessive longitudinal slip.',
            tire_model_four_wheel_limited: 'All four tires are near or beyond the grip limit.',
            tire_model_thermal_limited: 'Tire temperature or front/rear tire temperature spread is limiting grip.',
            tire_model_platform_limited: 'Suspension travel or platform control is likely limiting tire contact.',
        },
        tireLabExplanationLabels: {
            tire_model_waiting_for_samples: 'Waiting for valid 324-byte telemetry samples.',
            tire_model_stationary_explanation: 'Speed and inputs are near zero, so slip, temperature, and suspension values are shown as telemetry only and do not trigger a tire warning.',
            tire_model_no_dynamic_load_explanation: 'The current window does not include enough speed, steering, braking, throttle, or G-load to judge tire limit.',
            tire_model_balanced_explanation: 'No single tire group dominates the slip evidence in the current window.',
            tire_model_balanced_near_limit_explanation: 'Both axles are working hard; reduce speed or collect a cleaner sample before tuning.',
            tire_model_front_explanation: 'Front tire combined slip is higher than rear tire slip, indicating front grip is the current bottleneck.',
            tire_model_front_brake_explanation: 'Braking and steering together are loading the front tires beyond available grip.',
            tire_model_rear_explanation: 'Rear tire combined slip is higher than front tire slip, indicating rear stability is the current bottleneck.',
            tire_model_rear_power_explanation: 'Throttle is active while rear slip rises, so rear grip or power delivery is the bottleneck.',
            tire_model_rear_handbrake_explanation: 'Handbrake input is active while rear slip rises, so this sample should be treated as a driver-induced rotation event.',
            tire_model_rear_traction_explanation: 'Rear slip ratio is dominant under throttle, so wheel torque is exceeding rear tire traction.',
            tire_model_front_traction_explanation: 'Front slip ratio is dominant under throttle, so wheel torque is exceeding front tire traction.',
            tire_model_drive_traction_explanation: 'Driven tire slip ratio is high under throttle, so traction should be solved before gearing.',
            tire_model_four_wheel_explanation: 'Front and rear tires are both near the limit; this points to total grip, speed, aero, or platform rather than one axle only.',
            tire_model_thermal_explanation: 'High tire temperature or a large front/rear temperature spread can reduce available grip.',
            tire_model_platform_explanation: 'Suspension travel is nearly exhausted, which can interrupt tire contact and make grip inconsistent.',
        },
        tireLabWarningLabels: {
            tire_model_no_data: 'No tire model samples are available yet.',
            tire_model_sample_insufficient: 'Sample count is low; treat the model result as provisional.',
            sample_insufficient: 'Sample count is low; treat the model result as provisional.',
            thermal_risk: 'Tire temperature risk is present, but it is not treated as tire limit without slip evidence.',
            platform_risk: 'Suspension/platform risk is present, but it is not treated as tire limit without slip evidence.',
            left_right_imbalance: 'Left/right tire slip differs strongly. This is shown as evidence only, not as a single-side tuning action.',
            tire_model_left_right_imbalance: 'Left/right tire slip differs strongly. This is shown as evidence only, not as a single-side tuning action.',
            g_force_axis_mapping_unverified: 'G-force uses raw packet X/Y/Z axes. Longitudinal/lateral/vertical mapping still needs calibration.',
            camber_inference_no_three_point_temps: 'Camber inference is indirect because Data Out exposes only one tire temperature per tire.',
        },
        tireLabHintDirections: {
            improve_front_grip_or_reduce_entry_load: 'Improve front axle grip or reduce entry load.',
            improve_rear_grip_or_smooth_rotation: 'Improve rear axle grip or smooth rotation.',
            reduce_wheel_torque_or_improve_driven_tire_grip: 'Reduce wheel torque at the driven tires or improve driven tire grip.',
            reduce_speed_or_increase_total_grip: 'Reduce entry speed or increase total grip/platform support.',
            bring_tires_back_to_temperature_window: 'Bring tire temperatures back into a usable window.',
            restore_suspension_travel_and_platform_control: 'Restore suspension travel and platform control.',
            collect_moving_tire_load_samples: 'Start moving to collect loaded tire samples.',
            collect_more_representative_corner_and_power_samples: 'Collect more representative corner and power samples.',
            consider_more_negative_front_camber: 'Consider slightly more negative front camber only after validating with repeatable cornering samples.',
            consider_more_negative_rear_camber: 'Consider slightly more negative rear camber only after validating with repeatable cornering samples.',
            use_camber_inference_as_low_confidence_evidence: 'Use camber inference as low-confidence evidence, not as a direct adjustment command.',
        },
        tireLabHintLabels: {
            front_axle_grip: 'Front axle grip',
            rear_axle_grip: 'Rear axle grip',
            driven_tire_traction: 'Driven tire traction',
            whole_car_grip: 'Whole-car grip',
            tire_temperature: 'Tire temperature',
            platform_stability: 'Platform stability',
            observe: 'Observe',
            front_camber_check: 'Front camber check',
            rear_camber_check: 'Rear camber check',
            camber_observe: 'Camber observation',
        },
        fieldDiagnostics: 'Telemetry Field Diagnostics',
        speedCalibration: 'Speed Calibration',
        vehicleMetadata: 'Vehicle Metadata',
        enginePower: 'Engine / Power',
        motionPose: 'Motion / Pose',
        raceLapData: 'Race / Lap',
        auxiliaryFields: 'Auxiliary Fields',
        fieldName: 'Field',
        fieldValue: 'Value',
        fieldUnit: 'Unit',
        fieldSource: 'Source',
        fieldRange: 'Expected',
        fieldState: 'State',
        ok: 'OK',
        checkValue: 'Check',
        noCurrentFrame: 'No current telemetry frame.',
        replayTimeline: 'Replay Timeline',
        pauseReplay: 'Pause',
        resumeReplay: 'Resume',
        replayPosition: 'Position',
        sessionCompare: 'Session Compare',
        sessionMode: 'Session mode',
        modeCompareWarning: 'These sessions use different telemetry modes. Compare tuning events carefully.',
        advancedSettings: 'Advanced settings',
        testConditions: 'Test conditions',
        restoreUnknown: 'Restore unknown',
        driverMode: 'Driver detection',
        brakeAssist: 'Brake',
        steeringAssist: 'Steering',
        tractionControl: 'TCS',
        stabilityControl: 'STM',
        shifting: 'Shifting',
        launchControl: 'Launch control',
        assists: 'Assists',
        comparabilityWarnings: {
            game_mode_mismatch: 'Telemetry modes differ; compare these sessions as reference only.',
            test_conditions_unknown: 'One or both sessions have unknown test conditions, lowering comparison confidence.',
            driver_mode_mismatch: 'Driver modes differ.',
            brake_assist_mismatch: 'Brake assist settings differ.',
            steering_assist_mismatch: 'Steering settings differ.',
            traction_control_mismatch: 'Traction control settings differ.',
            stability_control_mismatch: 'Stability control settings differ.',
            shifting_mismatch: 'Shifting settings differ.',
            launch_control_mismatch: 'Launch control settings differ.',
        },
        testConditionValues: {
            unknown: 'Unknown',
            player: 'Player',
            auto: 'Auto driver',
            assisted: 'Assisted',
            abs_on: 'ABS on',
            abs_off: 'ABS off',
            standard: 'Standard',
            simulation: 'Simulation',
            on: 'On',
            off: 'Off',
            automatic: 'Automatic',
            manual: 'Manual',
        },
        leftSession: 'Left session',
        rightSession: 'Right session',
        compare: 'Compare',
        metric: 'Metric',
        left: 'Left',
        right: 'Right',
        delta: 'Delta',
        eventDistribution: 'Event distribution',
        profileCompare: 'Profile Version Compare',
        openProfileCompare: 'Compare versions',
        profileA: 'Profile A',
        profileB: 'Profile B',
        changedFields: 'Changed fields',
        recentChanges: 'Recent changes',
        noRecentChanges: 'No saved changes yet.',
        expertWorkspace: 'Professional Tuning Workspace',
        expertStartHint: 'Edit a tune profile, start live telemetry, then review in-memory problems, decisions, and interpreter advice. No session or recording is saved.',
        professionalDiagnosticTitle: 'Live tuning analysis',
        professionalDiagnosticHint: 'Uses the detector, decisioner, and interpreter selected in Developer Mode.',
        professionalDiagnosticEmpty: 'Start professional telemetry to see live vehicle problems and advice.',
        professionalMergedTitle: 'Merged tuning diagnosis',
        professionalMergedHint: 'Problems, decisions, interpreter advice, and concrete adjustments are grouped by issue.',
        professionalMergedEmpty: 'No merged diagnosis yet.',
        professionalMergedProblem: 'Problem',
        professionalMergedDecision: 'Decision',
        professionalMergedAdvice: 'Advice',
        professionalMergedAdjustments: 'Adjustment values',
        professionalMergedNoAdjustments: 'No concrete value yet. Fill the related tune fields or use this as a directional suggestion.',
        professionalTuneInstruction: 'How to tune',
        professionalChangeAmount: 'Change amount',
        selectProfileToEdit: 'Select a tune profile on the left, or create a new one.',
        newProfileModalTitle: 'Create tune profile',
        newProfileModalHint: 'Create the base profile first, then complete detailed tuning values in Expert Tuning.',
        createAndEdit: 'Create and edit',
        fillTelemetryIntoDraft: 'Fill vehicle identity from telemetry',
        snapshotChangedCount: (count: number) => `${count} changed field${count === 1 ? '' : 's'}`,
        compareWithCurrent: 'Compare with current',
        restoreSnapshot: 'Restore',
        snapshotRestored: 'Snapshot restored.',
        restoreSnapshotConfirm: 'Restore this snapshot to the current tune profile?',
        profileSnapshotCompare: 'Current vs Recent Change',
        snapshotBefore: 'Snapshot',
        currentSettings: 'Current settings',
        moreActions: 'More actions',
        noChanges: 'No changed fields.',
        close: 'Close',
        chooseProfile: 'Choose tune profile',
        profileChoiceTitle: 'Select tune for this telemetry session',
        profileChoiceHint: 'Multiple tune profiles match the current vehicle. Choose the one used for this run.',
        noMatchingProfile: 'No tune profile matches the current vehicle yet. Telemetry will start without a bound profile.',
        profileMatchUnavailable: 'No current telemetry frame is available, so vehicle matching cannot be checked.',
        profileMismatchTitle: 'Telemetry vehicle mismatch',
        profileMismatchHint: 'The selected tune profile does not match the current telemetry vehicle. Choose a matching profile, use no tune profile, or cancel.',
        telemetryVehicle: 'Telemetry vehicle',
        currentTuneVehicle: 'Selected tune profile',
        chooseMatchingProfile: 'Choose matching profile',
        clearProfileAndStart: 'No tune profile',
        profileSessions: 'Sessions',
        recentSession: 'Recent',
        ruleThresholds: 'Rule Thresholds',
        strategyTemplates: 'Strategy Templates',
        strategyAnalysis: 'Five-session strategy analysis',
        strategyTemplate: 'Strategy template',
        selectedSessions: 'Selected sessions',
        runStrategyAnalysis: 'Run analysis',
        strategyAnalysisEmpty: 'Select up to 5 sessions and a strategy template to analyze rule matching.',
        strategyRecommendation: 'Recommendation',
        strategyHints: 'Analysis hints',
        enabledEvents: 'Enabled events',
        totalEvents: 'Total events',
        ruleProfiles: 'Rule profiles',
        ruleName: 'Name',
        ruleCarClass: 'Class match',
        ruleDrivetrain: 'Drivetrain match',
        ruleUseCase: 'Use case match',
        ruleConfigJson: 'Config JSON',
        resetDefaults: 'Reset defaults',
        createProfile: 'Create profile',
        updateProfile: 'Update profile',
        newProfile: 'New profile',
        duplicate: 'Duplicate',
        delete: 'Delete',
        deleteSession: 'Delete session',
        deleteSessionConfirm: (name: string) => `Delete telemetry session "${name}" and its replay recording?`,
        deleteSessionBlocked: 'Stop replay or telemetry before deleting a session.',
        setActive: 'Set active',
        active: 'Active',
        profileList: 'Profiles',
        profileForm: 'Tune Profile',
        fillFromTelemetry: 'Fill from telemetry',
        telemetryFillUnavailable: 'No current telemetry frame to fill from.',
        telemetryFilled: 'Telemetry vehicle fields filled.',
        requiredCarName: 'Car name is required.',
        saved: 'Saved',
        saveAction: 'Save',
        deleted: 'Deleted',
        reportSessions: 'Telemetry Sessions',
        generateReport: 'Generate report',
        noSessions: 'No saved sessions yet.',
        noProfiles: 'No tune profiles yet.',
        profileIdentity: 'Vehicle identity',
        profileIdentityHint: 'Used to match telemetry sessions, reports, and tune versions.',
        profileTelemetryMatch: 'Matches current telemetry vehicle.',
        profileTelemetryMismatch: 'Does not match current telemetry vehicle.',
        profileTelemetryUnavailable: 'No current telemetry vehicle to compare.',
        profileEditorActions: 'Profile actions',
        markdownReport: 'Markdown Report',
        reportPlaceholder: 'Select a saved session to generate a report.',
        reportDecisionTitle: 'Tuning Decision Report',
        reportStatus: 'Report status',
        issueAdvice: 'Problems and advice',
        wholeCarPlan: 'Whole-car tuning plan',
        wholeCarPlanEmpty: 'No whole-car adjustment plan yet. Bind a tune profile and collect comparable samples.',
        roadTuningDecision: 'Road tuning decision',
        roadTuningDecisionEmpty: 'No clear road tuning decision yet. Use comparable Road sessions and telemetry evidence.',
        primaryIssue: 'Primary issue',
        primaryCause: 'Primary cause',
        cornerPhase: 'Phase',
        driverFitVerdict: 'Retest verdict',
        rollbackRecommended: 'Rollback recommended',
        rollbackRecommendedHint: 'The previous comparable test got worse after related changes. Roll back or reverse fine-tune before adding more changes.',
        retestFocus: 'Retest focus',
        knowledgeSource: 'Knowledge source',
        knowledgeFallback: 'Fallback rules',
        autoApplicable: 'Can apply',
        manualCheck: 'Manual check',
        optional: 'Optional',
        knowledgeStatus: 'Road model rules',
        reloadKnowledge: 'Reload rules',
        knowledgeReloaded: 'Tuning knowledge reloaded.',
        knowledgeSymptoms: 'Symptoms',
        knowledgeActions: 'Actions',
        roadDecisionStatusLabels: {
            ready: 'Ready',
            rollback_recommended: 'Rollback first',
            no_matching_symptom: 'No matching symptom',
            insufficient_data: 'Insufficient data',
            profile_unbound: 'No tune profile',
            knowledge_error: 'Rule load error',
        },
        roadActionRoles: {
            primary: 'Primary',
            support: 'Support',
            alternative: 'Alternative',
        },
        roadPhaseLabels: {
            launch: 'Launch',
            braking: 'Braking',
            corner_entry: 'Corner entry',
            mid_corner: 'Sustained cornering',
            corner_exit: 'Corner exit',
            cornering: 'Cornering',
            high_speed: 'High speed',
            platform: 'Platform',
            tires: 'Tires',
            power: 'Power',
        },
        driverFitVerdictLabels: {
            unknown: 'Unknown',
            improved: 'Telemetry improved',
            worsened: 'Telemetry worsened',
            unchanged: 'Unchanged',
            insufficient_data: 'Insufficient data',
        },
        speedBandLabels: {
            1: 'Low speed',
            2: 'Mid speed',
            3: 'High speed',
        },
        retestFocusLabels: {
            same_car: 'Same car',
            same_track: 'Same track',
            same_driver_mode: 'Same driver mode',
            verify_rollback_before_new_changes: 'Verify rollback first',
            gear_power_window: 'Gear power window',
            road_launch_wheelspin: 'Launch wheelspin',
            road_launch_bog_down: 'Launch bog down',
            road_entry_understeer: 'Corner-entry understeer',
            road_entry_understeer_low_speed: 'Low-speed entry understeer',
            road_entry_understeer_mid_speed: 'Mid-speed entry understeer',
            road_entry_understeer_high_speed: 'High-speed entry understeer',
            road_mid_understeer: 'Sustained-corner understeer',
            road_power_understeer: 'Power understeer',
            road_exit_oversteer: 'Exit oversteer',
            road_lift_snap_oversteer: 'Lift-off snap oversteer',
            road_front_brake_lockup: 'Front brake lockup',
            road_rear_brake_lockup: 'Rear brake lockup',
            road_high_speed_slide: 'High-speed slide',
            road_bottom_out: 'Bottom-out',
            road_tire_overheat: 'Tire temperature',
        },
        roadSymptomLabels: {
            road_launch_wheelspin: 'Launch wheelspin',
            road_launch_bog_down: 'Launch bog down',
            road_entry_understeer: 'Corner-entry understeer',
            road_entry_understeer_low_speed: 'Low-speed entry understeer',
            road_entry_understeer_mid_speed: 'Mid-speed entry understeer',
            road_entry_understeer_high_speed: 'High-speed entry understeer',
            road_mid_understeer: 'Sustained-corner understeer',
            road_power_understeer: 'Power understeer',
            road_exit_oversteer: 'Power oversteer',
            road_lift_snap_oversteer: 'Lift-off snap oversteer',
            road_front_brake_lockup: 'Front brake lockup',
            road_rear_brake_lockup: 'Rear brake lockup',
            road_high_speed_slide: 'High-speed four-wheel slide',
            road_bottom_out: 'Suspension bottom-out',
            road_tire_overheat: 'Tire overheat / imbalance',
            road_gearing_power: 'Gear power window',
        },
        tunePlanDraft: 'Tune plan draft',
        tunePlanDraftEmpty: 'No applicable tune plan actions for this session.',
        applyTunePlan: 'Apply selected to tune profile',
        tunePlanApplied: 'Tune plan applied to the bound profile.',
        tunePlanStatusLabels: {
            ready: 'Ready',
            no_actions: 'No applicable actions',
            profile_unbound: 'No bound tune profile',
            vehicle_mismatch: 'Vehicle mismatch',
            cannot_verify_vehicle: 'Cannot verify vehicle match',
        },
        tunePlanBlockedReasons: {
            field_locked_or_blank: 'Blank or locked field',
            vehicle_mismatch: 'Vehicle mismatch',
            cannot_verify_vehicle: 'Cannot verify vehicle match',
            profile_unbound: 'No bound tune profile',
            no_numeric_adjustment: 'No numeric adjustment',
            manual_review_required: 'Manual review required',
            no_change: 'No change',
            rollback_first: 'Retest got worse. Apply rollback first.',
            duplicate_action_removed: 'Duplicate suggestion removed',
            same_field_direction_conflict: 'Conflicting direction was resolved',
        },
        tunePlanTrust: 'Trust',
        tunePlanMissingInputs: 'Missing inputs',
        tunePlanRetestGuard: 'Retest guard',
        tunePlanTrustLevels: {
            high: 'High trust',
            medium: 'Medium trust',
            low: 'Low trust',
            blocked: 'Blocked',
        },
        tunePlanTrustReasons: {
            low_retest_confidence: 'Low retest confidence',
            retest_worsened: 'Retest got worse',
            recent_tune_plan_apply: 'From last applied tune plan',
            missing_current_value: 'Current value is missing',
            model_confidence_high: 'High model confidence',
            model_confidence_medium: 'Medium model confidence',
            model_confidence_low: 'Low model confidence',
            source_gear_power_diagnostic: 'Gear power diagnostic',
            source_road_tuning_model: 'Road tuning model',
            source_local_rule_report: 'Local rule report',
        },
        tunePlanMissingInputLabels: {
            tune_profile: 'Tune profile',
            current_tune_value: 'Current tune value',
            profile_power_band: 'Power-band RPM inputs',
            high_load_samples: 'High-throttle samples',
        },
        retestConfidence: 'Retest confidence',
        retestBaselineReason: 'Baseline',
        retestChangedFields: 'Recent tune changes',
        retestRollbackActions: 'Rollback suggestions',
        retestResult: 'Retest result',
        retestEmpty: 'No comparable previous test yet.',
        retestMetricLabels: {
            issue_score: 'Problem score',
            event_count: 'Event count',
            event_duration_ms: 'Event duration',
            avg_speed_kmh: 'Average speed',
            max_speed_kmh: 'Max speed',
            best_run_duration_ms: 'Best segment time',
            risk_score: 'Risk score',
            front_tire_temp: 'Front tire temp',
            rear_tire_temp: 'Rear tire temp',
            gear_problem_count: 'Gear power issues',
        },
        retestStatusLabels: {
            improved: 'Improved',
            worsened: 'Worsened',
            unchanged: 'Unchanged',
            insufficient_data: 'Insufficient data',
        },
        retestBaselineReasons: {
            matched_profile_track_driver: 'Same profile, track, and driver mode',
            matched_vehicle_class_usecase_driver: 'Same vehicle/class/use case and driver mode',
            missing_comparison_baseline: 'No comparable previous session',
            unavailable: 'Unavailable',
        },
        planStrategy: 'Strategy',
        planConfidence: 'Confidence',
        gearPowerDiagnostic: 'Gear power diagnostic',
        gearPowerDiagnosticHint: 'Works in race or free roam when there are enough high-throttle RPM, gear, speed, and slip samples.',
        gearPowerWhyNoAdvice: 'Why no gear advice',
        gearPowerNeedSamples: 'Collect more acceleration samples before judging gear ratios.',
        gearPowerNeedHighLoad: 'Need more clean high-throttle samples in unlocked gears.',
        gearPowerNoUnlockedGears: 'No filled or unlocked gear ratios were found in the active tune profile.',
        gearPowerFallbackLowConfidence: 'Power band is using RPM-ratio fallback; fill peak torque RPM, peak power RPM, and redline RPM for stronger advice.',
        gearPowerTractionFirst: 'Traction is limiting power delivery first, so differential or grip changes should be checked before changing gears.',
        gearPowerNoAdvice: 'No gear-ratio change is recommended from the current samples.',
        gearStrategyMode: 'Gear strategy',
        gearStrategyIssueCount: 'Issue gears',
        quickGearAdviceReadOnly: 'Quick diagnosis only gives direction. Bind a tune profile in Expert Tuning for concrete values and one-click apply.',
        gearPowerComparisons: 'Gear comparisons',
        gearTelemetryComparison: 'Measured gear performance',
        gearTuneComparison: 'Gear ratio setting changes',
        gearComparisonBefore: 'Before',
        gearComparisonAfter: 'After',
        gearComparisonDelta: 'Delta',
        gearComparisonUnavailable: 'No comparable gear data yet.',
        gearComparisonStatuses: {
            ready: 'Ready',
            missing_baseline: 'No previous comparable session',
            no_matching_gears: 'No matching gear samples',
            no_changed_gears: 'No recent gear-ratio changes',
            profile_unbound: 'No bound tune profile',
        },
        gearPowerSummary: 'Power window',
        powerBandTarget: 'Target RPM band',
        powerBandSource: 'Band source',
        diagnosticConfidence: 'Confidence',
        speedRange: 'Speed range',
        highestObserved: 'Max observed',
        inPowerBand: 'In band',
        inPowerBandCoverage: 'In-band coverage',
        acceleration: 'Acceleration',
        shiftAfter: 'Shift after',
        powerGearTest: 'Power / gear test',
        powerGearTestHint: 'For reliable gearing checks, run one clean full-throttle pull in each unlocked gear, then review sample quality and shift RPM here.',
        powerToWeight: 'Power to weight',
        tractionLimited: 'Traction limited',
        gearFinding: 'Finding',
        planConflicts: 'Resolved conflicts',
        planStrategies: {
            rollback_first: 'Rollback first',
            coarse_whole_car: 'Coarse whole-car pass',
            targeted_whole_car: 'Targeted whole-car pass',
        },
        planSummaries: {
            whole_car_template: 'Use the highest-impact combined changes first, then fine tune from the next run.',
            rollback_before_more_changes: 'Recent related changes made the result worse. Roll back before adding more changes.',
            no_clear_whole_car_action: 'No clear whole-car action was found from this session.',
        },
        planConfidenceLabels: {
            high: 'High',
            medium: 'Medium',
            low: 'Low',
            needs_profile: 'Needs tune profile',
        },
        gearFindings: {
            not_enough_samples: 'Not enough samples',
            not_enough_high_load: 'Not enough high-throttle samples',
            no_unlocked_gear_samples: 'No unlocked gear samples',
            gearing_window_ok: 'Gearing window OK',
            gearing_adjustment_needed: 'Gearing adjustment needed',
            traction_limited_power: 'Traction-limited power',
            global_too_long: 'Overall gearing too long, final drive first',
            global_too_short: 'Overall gearing too short, final drive first',
            single_gear_too_long: 'Individual gears too long',
            single_gear_too_short: 'Individual gears too short',
            single_gear_mixed: 'Mixed individual gear issues',
            traction_limited_low_gears: 'Traction limited, solve grip first',
            top_speed_limited: 'Top-speed gearing first',
            ok: 'OK',
            too_long: 'Too long under load',
            too_short: 'Too short under load',
            traction_limited: 'Traction limited',
            top_speed_limited_by_gearing: 'Top speed limited by gearing',
            top_speed_bog_down: 'Top gear too long',
            top_speed_ok: 'Top speed gearing OK',
            launch_wheelspin: 'Launch wheelspin',
            launch_bog_down: 'Launch bog down',
        },
        powerBandSources: {
            profile_power_band: 'Profile RPM inputs',
            telemetry_engine_max_rpm: 'Telemetry max RPM fallback',
            rpm_ratio_fallback: 'RPM ratio fallback',
        },
        cornerOperationStateLabels: {
            '1': 'Trail braking / braking',
            '2': 'Coasting',
            '3': 'Maintenance throttle',
            '4': 'Power on',
        },
        issueGroups: 'Merged problem groups',
        noIssueGroups: 'No merged problem groups for this session.',
        issueGroupAdviceTitle: 'Problem group advice',
        issueGroupEvents: 'Events in this group',
        issueGroupEvidence: 'Evidence range',
        issueGroupPrimaryAdvice: 'Primary tuning advice',
        issueStrategy: 'Adjustment strategy',
        issueStrategyLabels: {
            rollback_first: 'Rollback related changes first',
            coarse_combination: 'Coarse combined adjustment',
            medium_combination: 'Medium combined adjustment',
            fine_tune: 'Fine tune',
        },
        feedbackDirectiveLabels: {
            rollback_related_changes: 'The related last change made this worse. Roll back part of that change before adding more.',
            keep_direction_then_fine_tune: 'The direction improved the result. Continue with smaller steps.',
            avoid_more_same_direction: 'The issue did not clearly improve. Avoid more of the same direction until retested.',
        },
        issueGroupComparison: 'Compared with previous test',
        issueBaseline: 'Comparison baseline',
        issueRecentChanges: 'Recent tune changes',
        concreteProfileRequired: 'Bind a tune profile to generate concrete adjustment values.',
        noReportIssues: 'No saved problem events for this session.',
        profileBoundStatus: 'Tune profile bound',
        profileUnboundStatus: 'Tune profile not bound',
        sessionProfileUnboundHint: 'Bind a matching tune profile so the report can use the correct setup context.',
        driverModeUnknownReportHint: 'Driver detection confidence is low. Auto baseline and player-fit conclusions will be limited.',
        baselineMissingReportHint: 'Missing matched auto baseline. Run the same car and class on the same standard segment with auto driver.',
        standardSegmentMissingHint: 'No valid standard segment was detected. Create or match a benchmark track first.',
        advancedReportDetails: 'Advanced report details',
        expand: 'Expand',
        collapse: 'Collapse',
        playbackAndTimeline: 'Replay timeline',
        rawMarkdown: 'Raw Markdown',
        bindSessionProfile: 'Bind tune profile',
        changeSessionProfile: 'Change tune profile',
        sessionProfileBindTitle: 'Bind tune profile to session',
        sessionProfileBindHint: 'Only tune profiles matching this session vehicle can be selected.',
        sessionVehicle: 'Session vehicle',
        matchingTuneProfiles: 'Matching tune profiles',
        noSessionProfileMatches: 'No tune profile matches this session vehicle. Create one from telemetry first.',
        sessionProfileBound: 'Session tune profile updated. Generate the report again to use the new profile.',
        benchmarkTracks: 'Benchmark Tracks',
        trackProfilesTitle: 'Track Profiles',
        trackProfilesSubtitle: 'Capture route data and review auto-driving baselines by vehicle.',
        trackCaptureMode: 'Track capture',
        trackCaptureNoHistory: 'Track capture does not save test sessions, recordings, samples, or replay data.',
        trackData: 'Track data',
        vehicleReferences: 'Vehicle references',
        autoBaselines: 'Auto-driving baselines',
        noAutoBaselines: 'No valid auto-driving baseline is available for this track.',
        noVehicleReferences: 'No vehicle reference runs are available for this track.',
        bestAutoBaseline: 'Best auto baseline',
        baselineVehicle: 'Baseline vehicle',
        baselineRunCount: 'Baseline runs',
        validRuns: 'Valid runs',
        autoRuns: 'Auto runs',
        recentBenchmarkRuns: 'Recent benchmark runs',
        routeCompletion: 'Route completion',
        baselineWarnings: 'Baseline warnings',
        similarTrackFound: 'Similar track found',
        similarTrackHint: 'This route looks close to an existing track. Merge to preserve existing vehicle baselines, or save as a new track.',
        mergeIntoExistingTrack: 'Merge into existing track',
        saveAsNewTrack: 'Save as new track',
        startDistance: 'Start distance',
        endDistance: 'End distance',
        shapeSimilarity: 'Shape similarity',
        routeFitAvgError: 'Route fit avg error',
        routeFitP90Error: 'Route fit P90 error',
        matchLevel: 'Match level',
        strongMatch: 'Strong match',
        mediumMatch: 'Review match',
        autoMergedTrack: 'Track auto-merged',
        renameTrack: 'Rename track',
        trackRenamed: 'Track renamed.',
        trackBaselines: 'Vehicle baselines',
        trackBaselineCapture: 'Track baseline capture',
        startTrackBaseline: 'Start baseline capture',
        saveTrackBaseline: 'Save baseline',
        trackBaselineSaved: 'Track baseline saved.',
        trackBaselineAutoMatched: 'Baseline saved to matched track',
        trackBaselineAutoCreated: 'Baseline saved and new track created',
        trackBaselineAutoArchiveHint: 'Baseline capture auto-matches the route to an existing track; if no strong match is found, it creates a new track.',
        trackBaselineStopped: 'Track baseline capture stopped.',
        trackBaselineDeleted: 'Track baseline deleted.',
        trackBaselineNoSession: 'Vehicle baseline capture does not save test sessions or replay data.',
        confirmDeleteBaseline: 'Delete this vehicle baseline?',
        trackProfileWarnings: {
            no_auto_baseline: 'No valid auto-driving baseline has been detected for this track.',
            baseline_vehicle_identity_missing: 'Some auto runs are missing vehicle ID, class, or PI and were not grouped.',
        },
        trackBuilder: 'Track Builder',
        trackName: 'Track name',
        startCapture: 'Start capture',
        stopCapture: 'Stop capture',
        saveTrack: 'Save track',
        fromSession: 'Create from session',
        analyzeTrackRuns: 'Analyze track runs',
        capturePoints: 'Points',
        routeLength: 'Route length',
        drivingLineSignal: 'Driving line signal',
        detected: 'Detected',
        notDetected: 'Not detected',
        noTracks: 'No benchmark tracks yet.',
        noTrackPoints: 'No route points captured.',
        trackType: 'Track type',
        autoTrackType: 'Auto',
        circuitTrack: 'Circuit',
        sprintTrack: 'Point-to-point',
        extractionMode: 'Extraction mode',
        autoBestLap: 'Auto best lap',
        firstLap: 'First lap',
        fullSegment: 'Full segment',
        observedLaps: 'Observed laps',
        startGate: 'Start gate',
        finishGate: 'Finish gate',
        checkpoints: 'Checkpoints',
        setStartGate: 'Set start gate',
        setFinishGate: 'Set finish gate',
        clearGates: 'Clear gates',
        reextractTrack: 'Re-extract track',
        trackSaved: 'Benchmark track saved.',
        trackDeleted: 'Benchmark track deleted.',
        benchmarkRuns: 'Benchmark runs',
        noBenchmarkRuns: 'No benchmark run matched this session.',
        roadEvaluation: 'Road Racing Evaluation',
        roadEvaluationEmpty: 'No standard segment evaluation is available yet. Create or match a benchmark track first.',
        insufficientData: 'Insufficient data',
        paperPerformanceScore: 'Paper performance',
        playerFitScore: 'Player fit',
        riskScore: 'Risk',
        autoBaseline: 'Auto baseline',
        bestPlayerRun: 'Best player run',
        missingAutoBaselineHint: 'Only player performance can be reviewed right now. Add an auto-driver run with the same car, class, and standard segment to judge paper baseline.',
        prioritizeTuning: 'Prioritize tuning',
        evaluationContext: 'Evaluation context',
        evaluationContextEmpty: 'No evaluation attribution is linked to this issue yet.',
        prioritizeTuningYes: 'Prioritize tuning review',
        prioritizeTuningNo: 'Review driving style and repeatability before changing tune',
        roadVerdicts: {
            good_fit: 'Good fit',
            fast_but_risky: 'Fast but risky',
            paper_fast_not_fit: 'Paper-fast but not a fit',
            needs_tuning: 'Needs tuning',
            insufficient_data: 'Insufficient data',
        },
        roadBaselineStatuses: {
            matched_auto_baseline: 'Matched auto baseline',
            self_auto_baseline: 'This session is the auto baseline',
            missing_auto_baseline: 'Missing auto baseline',
            missing_vehicle_identity: 'Missing vehicle identity',
            no_valid_standard_run: 'No valid standard segment',
            no_standard_track: 'No standard track',
        },
        roadAttributions: {
            tune_issue: 'Tune issue',
            style_fit_issue: 'Driving style fit',
            driver_execution_issue: 'Driver execution',
            data_gap: 'Data gap',
        },
        roadAttributionMessages: {
            event_pattern: 'This issue is repeated enough to affect the road evaluation.',
            route_deviation: 'The standard segment detected route deviation.',
            route_progress_low: 'The run did not complete enough of the standard route.',
            geometry_length_mismatch: 'The measured geometry length differs from the saved standard route.',
            distance_traveled_mismatch: 'The Data Out distance field does not match the route geometry.',
            missing_auto_baseline: 'There is no matched auto-driver baseline for this vehicle and route.',
            missing_vehicle_identity: 'The session has no vehicle snapshot, so baseline matching is limited.',
            no_valid_standard_run: 'No valid benchmark run was detected.',
            no_standard_track: 'No standard track is available for this session.',
        },
        confidence: 'Confidence',
        bestRun: 'Best run',
        eventsSaved: 'Saved events',
        avgSpeed: 'Avg speed',
        maxSpeed: 'Max speed',
        duration: 'Duration',
        routeProgress: 'Route progress',
        sourceSession: 'Source session',
        updatedAt: 'Updated',
        geometryLength: 'Geometry length',
        lengthError: 'Length error',
        distanceDelta: 'DO distance delta',
        raceTimeDelta: 'DO time delta',
        lateralError: 'Lateral error',
        gateWidth: 'Gate width',
        gateDepth: 'Gate depth',
        warnings: 'Warnings',
        noWarnings: 'No warnings',
        trackSavedDetails: (id: number, type: string, length: number, session: string, laps: number) => `Track #${id} saved: ${type}, ${length.toFixed(0)} m, source ${session || '--'}, observed laps ${laps}.`,
        warningLabels: {
            distance_traveled_mismatch: 'Distance field mismatch',
            route_deviation: 'Route deviation',
            route_progress_low: 'Low route progress',
            geometry_length_mismatch: 'Geometry length mismatch',
        },
        useCases: {
            Road: 'Road',
            Rally: 'Rally',
            Drift: 'Drift',
            Offroad: 'Offroad',
            Drag: 'Drag',
            Wet: 'Wet',
            Test: 'Test',
        },
        fieldGroups: {
            vehicle: 'Vehicle',
            power: 'Power and weight',
            tire: 'Tires',
            gearing: 'Gearing',
            alignment: 'Alignment',
            antiroll: 'Anti-roll bars',
            springs: 'Springs and ride height',
            damping: 'Damping',
            aero: 'Aero',
            brake: 'Brakes',
            differential: 'Differential',
            notes: 'Notes',
        },
        eventTimeline: 'Event Timeline',
        eventSubtitle: 'Local rule detections from the current session',
        noEvents: 'No detected events in this session.',
        eventEvidence: 'Evidence',
        eventSuggestions: 'Initial suggestions',
        eventAdviceTitle: 'Tuning advice',
        advicePlaceholder: 'No exact tuning advice yet. Later versions will combine the tune profile, driving style, and same-condition comparisons to generate a focused recommendation.',
        tuningNote: 'Tuning note',
        eventDuration: 'Duration',
        eventSegment: 'Segment',
        eventStarted: 'Start',
        severityLabel: 'Severity',
        severityLow: 'Low',
        severityMedium: 'Medium',
        severityHigh: 'High',
        durationMsUnit: 'ms',
        durationSecondUnit: 's',
        events: {
            launch_wheelspin: 'Launch wheelspin',
            launch_bog_down: 'Launch bog down',
            short_gear: 'Short gear',
            long_gear_bog_down: 'Long gear bog down',
            top_speed_limited_by_gearing: 'Top speed limited by gearing',
            front_brake_lockup: 'Front brake lockup',
            rear_brake_lockup: 'Rear brake lockup',
            corner_entry_understeer: 'Corner entry understeer',
            mid_corner_understeer: 'Sustained-corner understeer',
            corner_exit_oversteer: 'Corner exit oversteer',
            power_understeer: 'Power understeer',
            snap_oversteer: 'Snap oversteer',
            high_speed_four_wheel_slide: 'High-speed four-wheel slide',
            tire_overheat: 'Tire overheat',
            tire_temp_imbalance: 'Tire temperature imbalance',
            suspension_bottom_out: 'Suspension bottom out',
        },
        segments: {
            launch: 'Launch',
            acceleration: 'Acceleration',
            braking: 'Braking',
            corner_entry: 'Corner entry',
            mid_corner: 'Sustained cornering',
            corner_exit: 'Corner exit',
            cornering: 'Cornering',
            high_speed_corner: 'High-speed corner',
            tire: 'Tire',
            suspension: 'Suspension',
        },
        evidenceLabels: {
            speed_kmh: 'Speed',
            speed_min_kmh: 'Min speed',
            speed_avg_kmh: 'Average speed',
            speed_max_kmh: 'Max speed',
            speed_band: 'Speed band',
            gear: 'Gear',
            throttle: 'Throttle',
            brake: 'Brake',
            steer_abs: 'Steering',
            front_slip_ratio: 'Front slip ratio',
            rear_slip_ratio: 'Rear slip ratio',
            max_slip_ratio: 'Max slip ratio',
            rpm_ratio: 'RPM ratio',
            front_combined_slip: 'Front combined slip',
            rear_combined_slip: 'Rear combined slip',
            slip_delta: 'Front-rear slip delta',
            corner_operation_state: 'Corner operation state',
            yaw_rate_abs: 'Yaw rate',
            max_suspension_travel: 'Max suspension travel',
            front_suspension: 'Front suspension',
            rear_suspension: 'Rear suspension',
            pitch_rate_abs: 'Pitch rate',
            roll_rate_abs: 'Roll rate',
            front_tire_temp: 'Front tire temp',
            rear_tire_temp: 'Rear tire temp',
            tire_temp_delta: 'Tire temp delta',
        },
        actionCategories: {
            aero: 'Aero',
            alignment: 'Alignment',
            brake: 'Brake',
            damping: 'Damping',
            differential: 'Differential',
            gearing: 'Gearing',
            rollback: 'Rollback',
            suspension: 'Suspension',
            tire: 'Tire',
        },
        actionItems: {
            brake_balance: 'Brake balance',
            brake_pressure: 'Brake pressure',
            bump: 'Bump damping',
            current_gear: 'Current gear',
            drive_diff_accel: 'Drive diff accel',
            drive_tire_pressure: 'Drive tire pressure',
            final_drive: 'Final drive',
            front_diff_accel: 'Front diff accel',
            front_diff_decel: 'Front diff decel',
            front_tire_pressure: 'Front tire pressure',
            front_and_rear_aero: 'Front and rear aero',
            front_arb: 'Front anti-roll bar',
            front_camber: 'Front camber',
            front_rebound: 'Front rebound',
            gear_1: '1st gear',
            gear_2: '2nd gear',
            gear_3: '3rd gear',
            gear_4: '4th gear',
            gear_5: '5th gear',
            gear_6: '6th gear',
            gear_7: '7th gear',
            gear_8: '8th gear',
            gear_9: '9th gear',
            gear_10: '10th gear',
            rear_arb: 'Rear anti-roll bar',
            rear_diff_accel: 'Rear diff accel',
            rear_diff_decel: 'Rear diff decel',
            rear_rebound: 'Rear rebound',
            rear_tire_pressure: 'Rear tire pressure',
            ride_height: 'Ride height',
            spring_rate: 'Spring rate',
            tire_pressure: 'Tire pressure',
        },
        actionDirections: {
            check: 'Check',
            decrease: 'Decrease',
            increase: 'Increase',
        },
        actionAmounts: {
            direction_only: 'direction only',
            'one small step': 'one small step',
            'slightly more negative': 'slightly more negative',
            'avoid bottoming': 'avoid bottoming',
        },
        actionReasons: {
            'avoid hitting the top of the gear too early': 'Avoid hitting the top of the gear too early',
            'help the engine stay in the power band': 'Help the engine stay in the power band',
            'increase front grip on entry': 'Increase front grip on entry',
            'increase front grip on steady cornering': 'Increase front grip during sustained cornering',
            'increase front grip under power': 'Increase front grip under power',
            'increase high-speed grip': 'Increase high-speed grip',
            'increase launch traction': 'Increase launch traction',
            'increase rear grip': 'Increase rear grip',
            'increase tire contact patch': 'Increase tire contact patch',
            'improve front tire contact in cornering': 'Improve front tire contact in cornering',
            'improve rear compliance under braking': 'Improve rear compliance under braking',
            'lengthen all gears if multiple gears are short': 'Lengthen all gears if multiple gears are short',
            'let the front tires load more smoothly': 'Let the front tires load more smoothly',
            'make threshold braking easier': 'Make threshold braking easier',
            'prevent aero and suspension instability': 'Prevent aero and suspension instability',
            'reduce bottoming frequency': 'Reduce bottoming frequency',
            'reduce driven-wheel slip': 'Reduce driven-wheel slip',
            'reduce front lockup tendency': 'Reduce front lockup tendency',
            'reduce power oversteer': 'Reduce power oversteer',
            'reduce rear lockup tendency': 'Reduce rear lockup tendency',
            'reduce wheel torque during launch': 'Reduce wheel torque during launch',
            'reduce wheel torque on exit': 'Reduce wheel torque on exit',
            'reduce power-on understeer': 'Reduce power-on understeer',
            'reduce sustained tire scrub': 'Reduce sustained tire scrub',
            'restore suspension travel': 'Restore suspension travel',
            'rotate the car more in steady cornering': 'Rotate the car more during sustained cornering',
            'make rear response less abrupt': 'Make rear response less abrupt',
            'stabilize the rear axle while off throttle': 'Stabilize the rear axle while off throttle',
            'reduce tire overheating tendency': 'Reduce tire overheating tendency',
            'balance front and rear tire temperatures': 'Balance front and rear tire temperatures',
            'rebalance axle load transfer': 'Rebalance axle load transfer',
            'rollback half of the last related change': 'Rollback half of the last related change',
            'stabilize tire temperature and contact patch': 'Stabilize tire temperature and contact patch',
            'shorten launch gearing': 'Shorten launch gearing',
            'shorten road acceleration gearing': 'Shorten road acceleration gearing',
            rollback_retest_worsened: 'Restore the value before the last applied tune plan',
            half_reverse_retest_worsened: 'Move half way back from the last applied tune plan',
            'stabilize the rear axle while braking': 'Stabilize the rear axle while braking',
            'support compression on impacts': 'Support compression on impacts',
            'increase top speed headroom': 'Increase top speed headroom',
            'verify aero drag is not limiting top speed': 'Verify aero drag is not limiting top speed',
            'verify traction is not limiting exit drive': 'Verify traction is not limiting exit drive',
        },
    },
    zh: {
        tireIssueGroups: '轮胎问题聚合',
        tireIssueSegments: '问题片段',
        tireIssueNoGroups: '当前窗口没有聚合后的轮胎问题。',
        tireIssueNoSegments: '暂无确认的问题片段。',
        tireIssueCount: '出现次数',
        tireIssueDuration: '累计持续时间',
        tireIssueSpeedRange: '速度范围',
        tireIssueOperations: '操作标签',
        tireIssueDriftSource: '漂移来源',
        tireIssueRisk: '风险',
        tireIssueEvidence: '代表证据',
        tireIssueAdvice: '修复建议',
        tireIssueExperimentHint: '实验解释层，不写入调校档案，也不影响正式报告。',
        tireIssueNoAdvice: '当前轮胎问题组无法生成修复方向。',
        tireIssuePriorityAdvice: '优先修复方向',
        tireIssueGroupAdvice: '按问题组展开的建议',
        tireIssuePrimaryCause: '主因假设',
        tireIssueShouldTune: '建议调校',
        tireIssueNoTune: '不建议调校',
        tireIssueRelatedFields: '相关调校杠杆',
        tireIssueVerifyEvidence: '需要验证的证据',
        tireIssueConflict: '冲突说明',
        tireIssueMissingInputs: '缺失条件',
        tireAdviceLayerLabels: {
            primary: '主方向',
            alternative: '备选方向',
            check: '检查',
            observe: '观察',
        },
        tireAdviceCategoryLabels: {
            tire_pressure: '胎压',
            alignment: '轮胎定位',
            antiroll: '防倾杆',
            spring_damping: '弹簧 / 阻尼',
            aero_platform: '空力 / 平台',
            brake: '刹车',
            differential: '差速器',
            gearing: '齿比',
            driver_input: '驾驶输入',
            data_quality: '数据可信度',
        },
        tireAdviceDirectionLabels: {
            increase_front_high_speed_support: '增加前轴高速支撑',
            increase_front_mechanical_grip: '增加前轴机械抓地',
            check_front_contact_patch: '检查前轮接地状态',
            increase_rear_stability: '增加后轴稳定性',
            check_rear_contact_patch: '检查后轮接地状态',
            reduce_four_wheel_lateral_load: '降低四轮横向负载',
            reduce_overlap_input: '减少复合输入',
            reduce_drive_lock: '降低驱动锁止',
            reduce_wheel_torque: '降低轮上扭矩',
            move_brake_balance_rearward: '制动力平衡向后移',
            move_brake_balance_forward: '制动力平衡向前移',
            check_decel_lock: '检查减速差速锁止',
            check_front_brake_platform: '检查前轴制动平台',
            check_brake_balance: '检查制动力平衡',
            check_platform: '检查平台支撑',
            check_temperature_window: '检查胎温窗口',
            check_left_right: '检查左右侧证据',
            continue_sampling: '继续采样',
            avoid_tuning: '该行为不作为修车目标',
        },
        tireAdviceCauseLabels: {
            data_not_reliable: '数据可信度不足',
            driver_handbrake_drift: '玩家主动手刹漂移',
            driver_weight_transfer_drift: '玩家主动重心转移漂移',
            front_high_speed_lateral_limit: '前轮先到高速横向极限',
            front_mechanical_lateral_limit: '前轴机械抓地不足',
            rear_lateral_stability_limit: '后轴横向稳定不足',
            four_wheel_lateral_limit: '四轮同时承受横向负载',
            drive_torque_exceeds_tire_grip: '动力输出超过轮胎牵引',
            driven_wheel_longitudinal_slip: '驱动轮纵向滑移',
            rear_brake_or_decel_instability: '后轴制动/减速不稳定',
            front_brake_overload: '前轴制动负载过高',
            combined_longitudinal_lateral_overload: '纵向和横向复合负载过高',
            platform_travel_or_load_risk: '平台行程或负载风险',
            tire_temperature_risk: '胎温风险',
            left_right_signal_or_load_imbalance: '左右侧信号或负载不均',
            unknown_tire_issue: '未知轮胎问题',
        },
        tireAdviceRationaleLabels: {
            front_high_speed_lateral_limit_prioritize_platform: '高速下先验证前轴空力和平台支撑，再调整低速机械平衡。',
            front_high_speed_lateral_limit_secondary_mechanical_grip: '如果平台证据正常，再检查前轴机械抓地和接地状态。',
            front_lateral_limit_prioritize_mechanical_grip: '该速度下前轮横向滑移占主导，优先看前轴机械抓地杠杆。',
            front_lateral_limit_verify_alignment: '先确认前轮接地状态，避免把问题简单归为防倾杆。',
            rear_lateral_limit_prioritize_stability: '后轮横向滑移占主导，优先看后轴稳定和平台响应。',
            rear_lateral_limit_verify_alignment: '确认后轮束角、外倾角和接地状态。',
            four_wheel_lateral_limit_reduce_platform_load: '前后轴都接近横向负载极限，应先降低整车平台负载。',
            four_wheel_lateral_limit_verify_driver_overlap: '先确认方向、油门、刹车是否存在复合输入。',
            front_traction_limit_reduce_accel_diff: '前驱动轮给油打滑，优先检查前侧加速差速。',
            rear_traction_limit_reduce_accel_diff: '后驱动轮给油打滑，优先检查后侧加速差速。',
            driven_traction_limit_balance_diff: '驱动轮牵引受限，检查驱动轴差速锁止和扭矩分配。',
            front_traction_limit_check_gearing: '如果差速不能解决，再检查齿比导致的轮上扭矩。',
            rear_traction_limit_check_gearing: '如果差速不能解决，再检查齿比导致的轮上扭矩。',
            driven_traction_limit_check_gearing: '检查终传比和低挡是否让驱动轮负载过高。',
            rear_braking_limit_move_balance_forward: '后轮制动滑移，制动力应更偏前或降低压力。',
            rear_braking_limit_check_decel_diff: '后侧减速差速可能增加制动阶段不稳定。',
            front_braking_limit_move_balance_rearward: '前轮制动负载过高，降低前轴制动负担或压力。',
            front_braking_limit_check_platform: '前轴平台支撑可能让制动载荷转移过于突然。',
            combined_limit_trail_brake_reduce_overlap: '刹车和转向叠加导致轮胎复合负载过高。',
            combined_limit_verify_brake_balance: '如果带刹入弯是主动操作，需要验证制动力平衡和压力。',
            combined_limit_corner_exit_reduce_overlap: '给油和转向叠加导致轮胎复合负载过高。',
            combined_limit_verify_drive_lock: '如果出弯给油是主动操作，需要验证驱动锁止和扭矩分配。',
            combined_limit_reduce_combined_load: '这是复合负载问题，不是单一轴向极限。',
            combined_limit_verify_platform: '先确认平台稳定，再调整单个轮胎杠杆。',
            platform_risk_check_travel_and_damping: '悬挂平台风险需要和轮胎滑移分开验证。',
            thermal_risk_check_pressure_and_slip: '胎温是风险信号，需要结合胎压和持续滑移验证。',
            left_right_imbalance_verify_route_and_sensor: '左右差只作为证据，先确认路线方向和传感器一致性。',
            data_quality_continue_sampling: '需要更多动态样本后再判断是否调校。',
            handbrake_drift_driver_behavior: '手刹触发的漂移属于驾驶行为，不作为修车目标。',
            scandinavian_flick_driver_behavior: '钟摆/重心转移漂移属于驾驶行为，除非非主动重复出现。',
            low_confidence_verify_before_tuning: '当前信号低可信，先验证模式再调校。',
            unknown_issue_continue_sampling: '问题分类不够明确，继续采样。',
        },
        tireIssueTypeLabels: {
            lateral_limit: '横向极限',
            traction_limit: '牵引极限',
            braking_limit: '制动极限',
            combined_limit: '综合极限',
            platform_risk: '平台风险',
            thermal_risk: '胎温风险',
            left_right_imbalance: '左右差风险',
            data_insufficient: '数据不足',
        },
        tireOperationTagLabels: {
            throttle_on: '给油',
            throttle_steady: '稳油',
            throttle_lift: '松油',
            light_brake: '轻刹',
            heavy_brake: '重刹',
            handbrake_active: '手刹',
            steer_increasing: '方向增加',
            steer_holding: '持续转向',
            steer_unwinding: '回正',
            speed_rising: '速度上升',
            speed_falling: '速度下降',
        },
        tireDriftSourceLabels: {
            handbrake_initiated: '手刹触发',
            power_oversteer: '动力甩尾',
            scandinavian_flick: '钟摆漂移',
            lift_off_oversteer: '松油甩尾',
            unknown_oversteer: '未知过度转向',
        },
        rawPackets: 'UDP 包',
        lastUdpPacket: '最近 UDP 包',
        lastUdpRemote: '最近 UDP 来源',
        title: 'FH6车辆调校工具',
        receiving: '接收中',
        idle: '未监听',
        networkAdapter: '网卡',
        udpPort: 'UDP 端口',
        start: '开始',
        stop: '停止',
        ready: '等待 FH6 Data Out UDP 遥测包。',
        listening: (address: string, port: number) => `正在监听 ${address}:${port}`,
        stopped: '监听已停止',
        gameTargetTitle: '游戏 Data Out 目标地址',
        gameTargetPrefix: 'FH6 Data Out IP 应设置为本机遥测端地址，例如',
        gameTargetMiddle: '端口',
        gameTargetSuffix: '游戏端地址',
        gameTargetSuffix2: '是发送方，不是目标地址。',
        allInterfaces: '所有网卡',
        loopback: '本机回环',
        lan: '局域网',
        ip: 'IP',
        speed: '速度',
        rpm: '转速',
        gear: '挡位',
        yawRate: '横摆角速度',
        driverInputs: '驾驶输入',
        raceDataActive: '比赛模式',
        waitingRaceState: '菜单 / 过场',
        freeRoamMode: '漫游模式',
        menuMode: '菜单 / 过场',
        mixedMode: '混合模式',
        unknownMode: '未知',
        notApplicable: '当前模式不适用',
        throttle: '油门',
        brake: '刹车',
        handbrake: '手刹',
        steering: '方向',
        validPackets: '有效包',
        invalidPackets: '异常包',
        parseErrors: '解析错误',
        wheelState: '车轮状态',
        wheelSubtitle: '滑移、胎温、悬挂',
        frontLeft: '左前轮',
        frontRight: '右前轮',
        rearLeft: '左后轮',
        rearRight: '右后轮',
        ratio: '滑移率',
        angle: '滑移角',
        combined: '综合滑移',
        temp: '胎温',
        susp: '悬挂',
        suspensionOffsetPct: '悬挂偏距 %',
        realtimeTrend: '实时趋势',
        trendSubtitle: '10Hz 聚合，最近 8 秒',
        rpmLoad: '转速负载',
        endpoint: '监听端点',
        mode: '模式',
        udpMode: 'UDP',
        replayMode: '回放',
        idleMode: '空闲',
        packetSize: '包大小',
        lastPacket: '最后包',
        speedSource: '速度来源',
        speedSourcePacket: 'Speed 字段',
        speedSourceVelocity: '速度向量',
        speedSourceNone: '无有效速度',
        packetSpeed: '包内速度',
        velocitySpeed: '向量速度',
        vehicleId: '车辆 ID',
        vehicleCategory: '车辆分类',
        classPi: '等级 / PI',
        drivetrainCylinders: '传动 / 气缸',
        engineOutput: '功率 / 扭矩',
        boostFuel: '增压 / 燃油',
        worldPosition: '世界坐标',
        lapRace: '圈数 / 名次',
        clutchHandbrake: '离合 / 手刹',
        surfaceSignals: '路面信号',
        recording: '录制',
        recordingSize: '录制大小',
        recordingPackets: '录制包数',
        recordingLimit: '录制上限',
        recordingTruncated: '录制已截断',
        recordingReady: '录制可用',
        samplesSaved: '保存样本',
        replay: '回放',
        stopReplay: '停止回放',
        replaySpeed: '回放速度',
        noRecording: '无录制',
        historicalTrend: '历史趋势',
        inputsTrend: '输入',
        frontRearSlip: '前后滑移',
        testLaunchpad: '测试启动台',
        testLaunchpadSubtitle: '开始录制前，先确认调校档案和测试条件。',
        analysisMode: '分析模式',
        quickDiagnosis: '快速诊断',
        expertTuning: '专家调校',
        quickModeNote: '只做实时诊断，不保存会话、录制或回放。',
        expertModeNote: '完整调校流程：保存会话、录制、报告、草稿应用和复测对比。',
        quickNoHistory: '快速模式：不保存历史或回放',
        quickDiagnosticTitle: '快速实时诊断',
        quickDiagnosticEmpty: '启动快速诊断并接收遥测后显示实时问题。',
        quickSuggestions: '建议方向',
        noQuickSuggestions: '暂无方向性建议。',
        quickSuggestionReason: '原因',
        quickSuggestionNextStep: '下一步',
        adviceLayerLabels: {
            rollback: '回退',
            primary: '主问题',
            powertrain: '动力总成',
            support: '辅助项',
            alternative: '备选项',
        },
        quickNextStepLabels: {
            bind_tune_profile_for_values: '绑定或新建调校档案后生成具体数值。',
            fill_or_unlock_tune_fields: '先填写或解锁相关调校字段。',
            collect_power_samples: '采集更多全油门加速样本。',
            use_expert_for_concrete_values: '进入专家调校获取具体数值和一键应用。',
        },
        missingProfileFields: '生成具体数值需要补充的档案字段',
        sameVehicleClass: '同车同级',
        sameTrackContext: '同一赛道上下文',
        comparisonConfidence: '对比可信度',
        comparabilityLabels: {
            yes: '是',
            no: '否',
            unknown: '未知',
        },
        quickConfidenceLabels: {
            high: '高',
            medium: '中',
            low: '低',
            invalid: '不可比',
        },
        quickWarningLabels: {
            quick_lap_data_insufficient: '圈数据不完整，已显示滚动窗口诊断。',
            quick_non_race_track_unknown: '漫游或非比赛遥测无法确认同一赛道。',
            quick_vehicle_or_class_changed: '本次快速诊断中车辆 ID 或等级发生变化。',
            quick_vehicle_class_unknown: '车辆 ID 或等级信息不完整。',
            quick_lap_clock_reset: '检测到圈计时重置。',
            quick_race_time_reset: '检测到比赛时间重置。',
            quick_race_time_missing: '缺少比赛时间，对比可信度降低。',
            quick_track_context_unknown: '无法确认赛道或比赛上下文。',
        },
        currentLap: '当前圈',
        previousLap: '上一圈',
        issueScore: '问题分',
        quickRollingWindow: '圈速数据不完整，当前显示最近窗口诊断。',
        quickComparisonStatuses: {
            lap_comparison: '当前圈 vs 上一圈',
            rolling_window_only: '滚动窗口',
            no_data: '暂无数据',
        },
        testConditionWarning: '驾驶方式为未知，自动驾驶基线和玩家适配结论会受限。',
        activeProfileMismatch: '当前遥测车辆与所选调校档案不匹配。',
        activeProfileReady: '调校档案与遥测车辆已匹配，可用于本次测试。',
        coreTelemetry: '核心遥测',
        startupFocus: '高级诊断、赛道提取和规则阈值已移至开发者模式。',
        never: '无',
        languageLabel: '语言',
        currentProfile: '当前调校',
        noProfile: '未选择调校档案',
        quickTab: '快速诊断',
        tireLabTab: '轮胎模型测试',
        tireRegressionTab: '样本验证',
        modelPipelineTab: '模型管线实验',
        trackProfilesTab: '赛道档案',
        recommendedCarsTab: '推荐车辆 JSON',
        tuneGeneratorTab: '快速调校',
        remoteTuneTab: '远程调校',
        expertTab: '专业调校',
        dashboardTab: '快速诊断',
        profilesTab: '专家调校',
        reportsTab: '测试报告',
        strategyTab: '策略模板',
        developerTab: '开发者模式',
        developerDoDiagnostics: 'DO 诊断',
        developerStrategyConfig: '策略 / 管线配置',
        recommendedCarsTitle: '推荐车辆 JSON 生成器',
        recommendedCarsSubtitle: '在本地数据库维护小程序推荐车辆信息，并导出 weChatApp/miniprogram/data/recommendedCars.json。',
        recommendedCarsHint: '导出会使用下方数据库列表，并替换当前 recommendedCars.json 文件。',
        recommendedCarsAdd: '保存到数据库',
        recommendedCarsClear: '清空表单',
        recommendedCarsGenerate: '导出文件',
        recommendedCarsPending: '数据库车辆',
        recommendedCarsNoItems: '数据库中暂无推荐车辆。',
        recommendedCarsSearch: '搜索推荐车辆',
        recommendedCarsSelectVisible: '选择当前结果',
        recommendedCarsClearSelection: '清除选择',
        recommendedCarsRefresh: '刷新',
        recommendedCarsSelected: (selected: number, total: number) => `已选择 ${selected} / 共 ${total} 条`,
        recommendedCarsNoSearchResults: '没有匹配的推荐车辆。',
        recommendedCarsCreatedAt: '最后修改时间（UTC+8）',
        recommendedCarsImageSrc: '图片',
        recommendedCarsNew: '新增车辆',
        recommendedCarsFormTitleNew: '新增推荐车辆',
        recommendedCarsFormTitleEdit: '编辑推荐车辆',
        recommendedCarsDetailTitle: '推荐车辆完整信息',
        recommendedCarsCancel: '取消',
        recommendedCarsExportEmpty: '请先选择至少一条推荐车辆再导出。',
        recommendedCarsSaved: (count: number, path: string) => `已生成 ${count} 台车辆：${path}`,
        recommendedCarsDbSaved: '推荐车辆已保存。',
        recommendedCarsDbDeleted: '推荐车辆已删除。',
        recommendedCarsDbDeletedAll: (count: number) => `已删除 ${count} 条数据库车辆。`,
        recommendedCarsTarget: '输出文件',
        recommendedCarsVersion: 'JSON 版本号',
        recommendedCarsFileCurrent: '当前 JSON 推荐状态',
        recommendedCarsFileFound: (count: number, version: string) => `recommendedCars.json 中已有 ${count} 台推荐车辆${version ? ` / ${version}` : ''}`,
        recommendedCarsFileNotFound: '尚未生成 recommendedCars.json。',
        recommendedCarsFileMissing: (count: number) => `文件中有 ${count} 个 ID 未在数据库中找到。`,
        recommendedCarsInFile: '推荐中',
        recommendedCarsAutoId: 'ID 会根据车辆名、用途、等级和 PI 自动生成。',
        recommendedCarsOptionalMeta: '车重和前轮重量分配为可选数据库辅助信息，不会导出到 JSON。',
        recommendedCarsEdit: '编辑',
        recommendedCarsDuplicate: '复制',
        recommendedCarsRemove: '移除',
        recommendedCarsDeleteAll: '全部删除数据库车辆',
        recommendedCarsDeleteAllConfirm: '确定要删除本地数据库中的全部推荐车辆吗？该操作不会删除 recommendedCars.json。',
        recommendedCarsCopied: '已复制为新车辆，请填写新的 tuneCode 后保存。',
        recommendedCarsDuplicateTuneCode: 'tuneCode 已存在。',
        recommendedCarsDuplicateIdentity: '该车辆身份已存在，请修改车辆名、用途、等级或 PI，或直接编辑原记录。',
        recommendedCarsTagsHint: '标签用逗号分隔，例如：抓地, 高速, 易上手。',
        tuneHarvestTab: '调校码采集',
        tuneHarvestTitle: '调校分享码采集',
        tuneHarvestSubtitle: '从公开来源收集 FH6 调校分享码，先进入审核池，再导入推荐车辆库。',
        tuneHarvestSources: '信息源',
        tuneHarvestSourceLabels: {
            jsr_chronic_sheet: 'JSR Chronic 表格',
            codmunity: 'CODMunity',
            forzafire: 'ForzaFire 详情页',
        },
        tuneHarvestDryRun: '只预览',
        tuneHarvestDryRunHint: '只预览会显示抽取候选，但不写入本地数据库。',
        tuneHarvestLimit: 'ForzaFire 详情页上限',
        tuneHarvestRun: '开始采集',
        tuneHarvestStop: '停止采集',
        tuneHarvestStopping: '正在停止采集...',
        tuneHarvestRefresh: '刷新审核池',
        tuneHarvestClear: '清空审核池',
        tuneHarvestClearConfirm: '确定删除审核池里的全部候选吗？',
        tuneHarvestCleared: (count: number) => `已清空审核池，删除 ${count} 条候选。`,
        tuneHarvestCandidates: '审核池',
        tuneHarvestNoCandidates: '暂无采集候选。',
        tuneHarvestSelectedSourcesRequired: '请至少选择一个信息源。',
        tuneHarvestResult: (found: number, saved: number, pending: number, rejected: number) => `发现 ${found} 条，保存 ${saved} 条，待确认 ${pending} 条，拒绝 ${rejected} 条。`,
        tuneHarvestStopped: '采集已停止。',
        tuneHarvestStatusFilter: '状态筛选',
        tuneHarvestSearch: '搜索',
        tuneHarvestSearchPlaceholder: '分享码、车辆、调校作者、来源、用途...',
        tuneHarvestStatusLabels: {
            all: '全部',
            pending: '待确认',
            rejected: '已拒绝',
            imported: '已导入',
        },
        tuneHarvestUseCandidate: '使用',
        tuneHarvestReject: '拒绝',
        tuneHarvestRestore: '恢复',
        tuneHarvestSource: '来源',
        tuneHarvestVehicle: '车辆',
        tuneHarvestCode: '分享码',
        tuneHarvestContext: '上下文',
        tuneHarvestMatch: '匹配',
        tuneHarvestStatus: '状态',
        tuneHarvestWarnings: '警告',
        tuneHarvestImported: '候选已标记为导入。',
        tuneHarvestRejected: '候选已拒绝。',
        tuneHarvestRestored: '候选已恢复为待确认。',
        tuneHarvestCopiedToRecommended: '候选已填入推荐车辆表单。',
        developerMode: '开发者模式',
        tireLab: '轮胎模型测试',
        tireLabTitle: '四轮轮胎模型测试',
        tireLabSubtitle: '实验性的轮胎中心模型，不影响现有快速诊断、测试报告或调校草稿。',
        tireLabNoPersistence: '轮胎模型测试：不保存会话、不录制、不写数据库、不支持回放。',
        tireLabEmpty: '启动轮胎模型测试并接收遥测后显示四轮抓地平衡。',
        tireRegressionTitle: '轮胎回归样本验证',
        tireRegressionSubtitle: '独立的轮胎模型样本验证工具，不影响快速诊断、测试报告或调校草稿。',
        tireRegressionSaveCurrent: '保存当前 Tire Lab 窗口',
        tireRegressionSampleName: '样本名称',
        tireRegressionScenario: '场景',
        tireRegressionWindowSeconds: '窗口秒数',
        tireRegressionSave: '保存样本',
        tireRegressionRunOne: '运行样本',
        tireRegressionRunAll: '运行全部样本',
        tireRegressionDelete: '删除样本',
        tireRegressionUpdateExpected: '保存期望',
        tireRegressionSamples: '回归样本',
        tireRegressionExpected: '期望规则',
        tireRegressionResults: '运行结果',
        tireRegressionNoSamples: '暂无回归样本。请先保存一个 Tire Lab 窗口。',
        tireRegressionNoSelection: '选择一个样本后查看或编辑期望规则。',
        tireRegressionAllowedPhases: '允许阶段',
        tireRegressionRequiredGrip: '必须检测的抓地极限',
        tireRegressionAllowedAxles: '允许受限轴',
        tireRegressionForbiddenGrip: '禁止误报的抓地极限',
        tireRegressionMinQuality: '最低数据可信度',
        tireRegressionNotes: '备注',
        tireRegressionActual: '实际输出',
        tireRegressionPassed: '通过',
        tireRegressionFailed: '失败',
        tireRegressionFailures: '失败原因',
        tireRegressionCsvHint: '用英文逗号分隔，例如 corner_exit, traction_limit, rear。',
        tireRegressionRequiresTireLab: '保存样本需要当前内存中已有轮胎模型测试数据。',
        tireRegressionSaved: '已保存轮胎回归样本。',
        tireRegressionExpectationSaved: '已保存期望规则。',
        tireRegressionDeleted: '已删除样本。',
        tireRegressionRan: '回归验证完成。',
        tireRegressionFailureLabels: {
            phase_mismatch: '阶段与期望不匹配',
            required_grip_missing: '未检测到必须出现的抓地极限',
            limited_axle_mismatch: '受限轴与期望不匹配',
            forbidden_grip_detected: '出现了禁止误报的抓地极限',
            data_quality_below_minimum: '数据可信度低于最低要求',
        },
        modelPipelineTitle: '模型管线实验',
        modelPipelineSubtitle: '独立只读运行检测模型、决策器和解释器组合，不修改会话或调校档案。',
        modelPipelineSource: '数据源',
        modelPipelineDetector: '检测模型',
        modelPipelineDecisioner: '决策器',
        modelPipelineInterpreter: '解释器',
        modelPipelineSession: '遥测会话',
        modelPipelineRun: '运行管线',
        modelPipelineRunComplete: '模型管线运行完成。',
        modelPipelineNoCatalog: '模型管线目录尚未加载。',
        modelPipelineNoResult: '选择数据源和模型组合后运行管线。',
        modelPipelineProblems: '问题集',
        modelPipelineDecisions: '决策结果',
        modelPipelineAdvice: '解释器建议',
        modelPipelineWarnings: '兼容性提示',
        modelPipelineDocs: '文档来源',
        modelPipelineSourceSummary: '数据源摘要',
        modelPipelineExplainOnly: '仅解释，不生成可写入数值，也不会修改调校档案。',
        modelPipelineStatus: '状态',
        modelPipelineConfidence: '可信度',
        modelPipelineShouldTune: '是否建议调校',
        modelPipelineEvidence: '证据',
        modelPipelineStatusLabels: {
            ready: '就绪',
            no_data: '无数据',
            incompatible: '组合不兼容',
        },
        modelPipelineAdviceCategoryLabels: {
            power_to_tire: '动力落地',
            platform_aero: '平台 / 空力',
            mechanical_grip: '机械抓地',
            observe: '观察',
        },
        modelPipelineAdviceDirectionLabels: {
            reduce_wheel_torque_or_drive_lock: '降低轮上扭矩或驱动锁止',
            verify_brake_balance_pressure_and_decel_lock: '检查制动力平衡、压力和减速差速',
            verify_platform_then_use_tiered_ride_height_or_aero: '先验证平台，再按档位调整车高或空力',
            verify_pressure_window_and_slip_heat: '检查胎压窗口和滑移升温',
            rebalance_tire_contact_and_load_transfer: '重新平衡轮胎接地和载荷转移',
            collect_more_evidence_before_tuning: '继续采样后再调校',
        },
        modelPipelineScopeLabels: {
            driven_wheels: '驱动轮',
            front_rear_balance: '前后平衡',
            front_rear_platform: '前后平台',
            front_rear_tires: '前后轮胎',
            vehicle: '整车',
        },
        modelPipelineRationaleLabels: {
            tire_problem_group_decision: '由轮胎问题组生成的调校判断',
            legacy_problem_fallback: '旧规则问题的兜底判断',
            legacy_road_decision_v1: '旧公路决策模型',
            tire_lab_problem_groups_v1: '轮胎模型问题组',
            docs_v12_interpreter_outputs_no_write_values: 'Docs v1.2 解释器只输出方向，不生成可写入数值',
            'Docs v1.2 treats gearing and differential as Forza slider/display levers; verify whether torque delivery exceeds tire traction before changing chassis balance.': 'Docs v1.2 将齿比和差速视为 Forza 滑块 / 显示值杠杆；先确认动力是否超过轮胎牵引，再调整底盘平衡。',
            'Docs v1.2 uses conservative road brake baselines by drivetrain; explain braking issues through balance, pressure, and decel lock before changing unrelated systems.': 'Docs v1.2 按驱动形式使用保守公路刹车基线；先从制动力平衡、压力和减速差速解释制动问题，再考虑无关系统。',
            'Docs v1.2 keeps ride height and aero as low/medium/high tier explanations, not precise write values; use them only after verifying platform evidence.': 'Docs v1.2 将车高和空力保留为低 / 中 / 高档解释，不输出精确写入值；需要先验证平台证据。',
            'Docs v1.4 uses BAR as the primary tire-pressure unit; thermal issues need pressure, slip, and temperature evidence before any 0.02-0.03 BAR correction is trusted.': 'Docs v1.4 使用 BAR 作为胎压主单位；胎温问题需要同时有胎压、滑移和温度证据，才建议 0.02-0.03 BAR 微调。',
            'Docs v1.2 makes tire pressure, alignment, anti-roll bars, springs, and damping explicit Forza display/sliders; use them as grouped levers around tire contact and load transfer.': 'Docs v1.2 将胎压、定位、防倾杆、弹簧和阻尼明确为 Forza 显示值 / 滑块；应作为围绕轮胎接地和载荷转移的一组杠杆使用。',
            'Docs v1.2/v1.4 make tire pressure, alignment, anti-roll bars, springs, and damping explicit Forza display/sliders; use BAR-first tire pressure only when slip and heat evidence support it.': 'Docs v1.2/v1.4 将胎压、定位、防倾杆、弹簧和阻尼明确为 Forza 显示值 / 滑块；只有滑移和温度证据支持时，才优先使用 BAR 胎压建议。',
            'Docs v1.2 interpreter did not find a specific road baseline lever for this decision; keep it as an observation.': 'Docs v1.2 解释器未找到对应的公路基线调校杠杆，暂作为观察项。',
        },
        modelPipelineEvidenceLabels: {
            driven_wheel_slip_ratio: '驱动轮滑移率',
            throttle: '油门',
            rpm: '转速',
            speed_gain: '速度增长',
            gear: '挡位',
            brake: '刹车',
            deceleration_g: '减速 G',
            front_slip_ratio: '前轮滑移率',
            rear_slip_ratio: '后轮滑移率',
            speed_band: '速度区间',
            suspension_offset: '悬挂偏距',
            combined_slip: '综合滑移',
            g_force: 'G 力',
            front_tire_temp: '前轮胎温',
            rear_tire_temp: '后轮胎温',
            front_combined_slip: '前轮综合滑移',
            rear_combined_slip: '后轮综合滑移',
            slip_angle: '滑移角',
            steer: '转向',
            sample_quality: '样本质量',
            phase: '阶段',
            speed: '速度',
            inputs: '输入',
        },
        tuneGeneratorTitle: '快速调校',
        tuneGeneratorSubtitle: '用最小车辆参数生成中性公路基线。不接 AI、不依赖遥测、不保存测试会话。',
        tuneGeneratorValueHint: '请输入Forza Horizon 6中车辆信息',
        tuneGeneratorNoHistory: '生成结果只有在创建或应用调校档案时才会保存。',
        remoteTuneTitle: '远程调校',
        remoteTuneSubtitle: '开启局域网 iPhone/iPad 快速调校预览页面。不支持复制结果，不支持保存调校档案。',
        remoteTuneStatus: '状态',
        remoteTunePort: 'Web 端口',
        remoteTuneLanAddress: '局域网地址',
        remoteTuneUrl: '访问地址',
        remoteTuneStart: '开启远程页面',
        remoteTuneStop: '关闭远程页面',
        remoteTuneRunning: '远程调校页面运行中。',
        remoteTuneStopped: '远程调校页面未开启。',
        remoteTuneReadOnly: 'MVP 仅生成预览：不复制结果、不保存档案、不写数据库。',
        remoteTuneDeviceHint: '请在同一局域网内用 iPhone/iPad Safari 打开该地址。',
        remoteTunePortInvalid: '端口必须为 1-65535 范围内的整数。',
        quickTuneInputButton: '输入车辆信息',
        quickTuneInputTitle: '车辆信息',
        quickTuneLastSummary: '上次快速调校',
        quickTuneNoResult: '暂无快速调校结果',
        quickTuneUseCase: '用途',
        quickTuneTireCompound: '轮胎类型',
        quickTuneUnsupportedUseCase: '当前 MVP 支持公路、漂移、拉力、越野和直线静态基线生成。',
        quickTuneDriftRwdPreferred: '漂移快速调校以后驱为优先。前驱/四驱会生成基础参数，但会跳过差速器数值。',
        quickTuneValidationSummary: '请先修正红色高亮的车辆信息。',
        quickTuneIntegerRange: (label: string, min: number, max: number) => `${label}必须为 ${min}-${max} 范围内的整数。`,
        quickTunePositiveInteger: (label: string) => `${label}必须为正整数。`,
        quickTuneTireSizeInvalid: '轮胎尺寸必须填写整数宽度、扁平比和轮毂尺寸。',
        quickTuneGearingToggle: '开启齿轮设置',
        quickTuneDriftGearingHint: '漂移齿比会围绕核心漂移挡生成，目标是在该速度附近保持高转。',
        quickTuneDragGearingHint: '直线齿比围绕目标尾速生成，优先保证加速和换挡后的转速衔接。',
        quickTuneTargetDriftSpeed: '目标漂移速度',
        quickTuneTargetDragSpeed: '目标终点速度',
        quickTuneCarryToProfessional: '使用专业调校',
        quickTuneCarriedToProfessional: '已将快速调校参数带入专业调校。确认后可在专业调校中保存。',
        quickTuneBiasTitle: '基线偏置滑块',
        quickTuneBiasHint: '先生成中性公路基线，再用整车偏置进行保守修正；保存或带入专业调校时使用修正后的参数。',
        quickTuneBalance: '平衡',
        quickTuneBalanceLeft: '稳定',
        quickTuneBalanceRight: '灵活',
        quickTuneStiffness: '硬度',
        quickTuneStiffnessLeft: '软',
        quickTuneStiffnessRight: '硬',
        quickTuneSpeed: '速度',
        quickTuneSpeedLeft: '极速',
        quickTuneSpeedRight: '加速',
        quickTuneNeutral: '中性',
        quickTuneSpeedDisabled: '开启齿轮设置并生成齿比后可调整速度偏置。',
        quickTuneDrivetrainLabels: {
            FWD: '前驱',
            AWD: '四驱',
            RWD: '后驱',
        },
        quickTuneTireCompoundLabels: {
            stock: '原厂胎',
            street: '街胎',
            sport: '运动胎',
            semi: '半热熔',
            slick: '光头胎',
            rally: '拉力胎',
            offroad: '越野胎',
            drift: '漂移胎',
            drag: '直线胎',
            snow: '雪地胎',
        },
        tuneGeneratorMinimum: '最小输入',
        tuneGeneratorAdvanced: '高级可选输入',
        tuneGeneratorPreview: '基线预览',
        tuneGeneratorFields: '生成字段',
        tuneGeneratorSkipped: 'MVP 暂不生成',
        tuneGeneratorNextTest: '下一轮测试',
        tuneGeneratorGenerate: '生成预览',
        tuneGeneratorCreate: '保存为专业档案',
        tuneGeneratorApply: '应用到已有专业档案',
        tuneGeneratorReset: '重置输入',
        tuneGeneratorOpenExpert: '进入专业调校',
        tuneGeneratorTarget: '目标调校档案',
        tuneGeneratorNoPreview: '输入最小参数后生成公路基线预览。',
        tuneGeneratorNoTarget: '应用前请选择目标调校档案。',
        tuneGeneratorNoSelection: '请至少选择一个生成字段。',
        tuneGeneratorCreated: '已创建基线调校档案。',
        tuneGeneratorApplied: '已应用到目标调校档案。',
        tuneGeneratorConfidence: '可信度',
        tuneGeneratorSelectedCount: '已选字段',
        tuneGeneratorRangeHint: '不知道车辆专属滑条范围时，范围字段可留空。',
        tuneGeneratorCarName: '车辆名称',
        tuneGeneratorVersionName: '版本',
        tuneGeneratorCarOrdinal: '车辆 ID',
        tuneGeneratorCarCategory: '车辆分类',
        tuneGeneratorPI: 'PI',
        tuneGeneratorDrivetrain: '驱动',
        tuneGeneratorWeight: '车重',
        tuneGeneratorFrontWeight: '前轮重量分配',
        tuneGeneratorPower: '功率',
        tuneGeneratorTorque: '扭矩',
        tuneGeneratorGearingInputs: '静态齿比',
        tuneGeneratorGearingHint: '可选。根据红线转速、挡位数、轮胎尺寸和目标极速生成终传比与 1-N 挡。',
        tuneGeneratorRedlineRPM: '红线转速',
        tuneGeneratorGearCount: '挡位数',
        tuneGeneratorTireDiameter: '轮胎尺寸',
        tuneGeneratorTargetTopSpeed: '目标极速',
        tuneGeneratorRideRange: '车高范围',
        tuneGeneratorAeroRange: '空力范围',
        tuneGeneratorAdjustableHint: '档位建议不会写入精确数值，请在游戏中按档位手动设置。',
        tuneGeneratorTierRecommendations: '档位建议',
        tuneGeneratorTierManual: '仅手动档位建议，不会应用到调校档案。',
        tuneGeneratorFrontRideAdjustable: '前侧车身高度可调',
        tuneGeneratorRearRideAdjustable: '后侧车身高度可调',
        tuneGeneratorFrontAeroAdjustable: '前侧下压力可调',
        tuneGeneratorRearAeroAdjustable: '后侧下压力可调',
        tuneGeneratorTierLabels: {
            low: '低档',
            medium: '中档',
            high: '高档',
        },
        tuneGeneratorMin: '最小',
        tuneGeneratorMax: '最大',
        tuneGeneratorReason: '依据',
        tireLabMatrix: '四轮抓地矩阵',
        tireLabAxleBalance: '前后轴抓地平衡',
        tireLabFrontAxle: '前轴',
        tireLabRearAxle: '后轴',
        tireLabLimitType: '轮胎极限状态',
        tireLabPhase: '阶段',
        tireLabConfidence: '模型可信度',
        tireLabExplanation: '模型解释',
        tireLabHints: '解释方向',
        tireLabWarnings: '风险/提示',
        tireLabLeftRight: '左右平衡',
        tireLabWindow: '窗口',
        tireDataQuality: '数据可信度',
        tireGripLimit: '抓地极限类型',
        tireGripRelationships: '四轮关系',
        tirePhaseStability: '阶段稳定性',
        tireStablePhase: '稳定阶段',
        tireScoreMargin: '评分差距',
        tireDynamicSamples: '动态样本',
        tireSignalQuality: '信号质量',
        tireLimitedAxle: '受限车轴',
        tireLimitedWheels: '受限轮胎',
        tirePrimaryEvidence: '主要证据',
        tireReason: '原因',
        tirePhaseCurrent: '当前阶段',
        tirePhaseSecondary: '次要阶段',
        tirePhaseEvidence: '阶段证据',
        tirePhaseScores: '阶段评分',
        tirePhaseSpeedDelta: '速度变化',
        tirePhaseSpeedReference: '速度基准',
        tirePhaseSpeedBand: '速度区间',
        tirePhaseThrottleDelta: '油门变化',
        tirePhaseSteerDelta: '方向变化',
        tirePhaseBrake: '刹车 平均/峰值',
        tirePhaseHandbrake: '手刹 平均/峰值',
        tirePhasePlaneG: '平面 G 平均/峰值',
        tirePhaseDecelG: '减速 G 平均/峰值',
        tirePhaseAccelG: '加速 G 平均/峰值',
        tireDataQualityLabels: {
            valid: '有效',
            low_confidence: '低可信',
            invalid: '无效',
            ok: '正常',
            low: '偏低',
            flat: '无变化',
        },
        tireDataQualityReasonLabels: {
            tire_data_no_samples: '暂无轮胎模型样本。',
            tire_data_menu_or_no_vehicle: '菜单或无车辆遥测，无法判断轮胎极限。',
            tire_data_sample_insufficient: '样本数量不足。',
            tire_data_dynamic_sample_insufficient: '动态轮胎负载样本不足。',
            tire_data_stationary: '车辆处于静止状态。',
            tire_data_speed_low: '速度信号过低，无法判断极限。',
            tire_data_g_force_flat: 'G 力信号接近无变化。',
            tire_data_slip_signal_flat: '车轮滑移信号为 0 或缺失。',
            tire_data_input_low: '驾驶输入信号过低。',
        },
        tireGripLimitLabels: {
            lateral_limit: '横向抓地极限',
            traction_limit: '牵引抓地极限',
            braking_limit: '制动抓地极限',
            combined_limit: '四轮综合极限',
            balanced_near_limit: '均衡接近极限',
            no_limit_detected: '未检测到轮胎极限',
        },
        tireGripReasonLabels: {
            tire_grip_no_dynamic_limit: '动态滑移证据未显示明确轮胎极限。',
            tire_grip_stationary: '静止状态不评估轮胎极限。',
            tire_grip_no_dynamic_load: '没有动态轮胎负载样本。',
            tire_grip_data_invalid: '数据质量无效。',
            tire_grip_lateral_slip: '滑移角证据指向横向抓地极限。',
            tire_grip_power_slip: '驱动轮滑移率证据指向牵引极限。',
            tire_grip_braking_slip: '制动滑移证据指向制动抓地极限。',
            tire_grip_handbrake_slip: '手刹输入导致后轮滑移。',
            tire_grip_four_wheel_combined: '前后轴都接近综合滑移极限。',
            tire_grip_near_limit: '单轴接近综合滑移极限。',
        },
        tireAxleLabels: {
            none: '无',
            front: '前轴',
            rear: '后轴',
            both: '前后轴',
            driven: '驱动轮',
        },
        tirePhaseStabilityLabels: {
            stable: '稳定',
            transition: '过渡',
            low_confidence: '低可信',
        },
        powerToTireTitle: '动力输出 → 轮胎牵引',
        powerToTireSubtitle: '实验性 DO 包动力落地检查，不依赖调校档案功率输入。',
        powerToTireStatus: '动力落地状态',
        powerToTireDrivenAxle: '驱动轮',
        powerToTirePower: '功率',
        powerToTireTorque: '扭矩',
        powerToTireRPM: '转速 / 比例',
        powerToTireGear: '挡位',
        powerToTireThrottle: '油门',
        powerToTireDrivenSlip: '驱动轮滑移 P90',
        powerToTireAccel: '加速度',
        powerToTireSamples: '高油门样本',
        powerToTireSignal: '功率信号',
        powerToTireAvailable: '可用',
        powerToTireUnavailable: '不可用',
        powerToTireSummaryLabels: {
            power_to_tire_no_data: '暂无动力落地数据',
            power_to_tire_insufficient: '高油门样本不足',
            power_to_tire_low_throttle: '油门不足',
            power_to_tire_power_signal_unavailable: '功率/扭矩信号不可用',
            traction_over_power: '动力超过驱动轮牵引',
            rpm_below_useful_range: '转速低于有效区间',
            rpm_too_high_or_gear_short: '转速过高 / 齿比偏短',
            power_not_reaching_ground: '动力没有有效落地',
            power_landing_ok: '动力落地状态可用',
        },
        powerToTireExplanationLabels: {
            power_to_tire_waiting_for_samples: '等待包含速度、转速、挡位、功率和轮胎滑移的 DO 样本。',
            power_to_tire_need_high_throttle: '需要更多高油门样本后再判断动力落地。',
            power_to_tire_needs_more_high_throttle: '需要更多高油门样本后再判断动力落地。',
            power_to_tire_low_throttle_explanation: '当前油门不足，不能形成动力到牵引的判断。',
            power_to_tire_power_signal_unavailable_explanation: 'DO 功率或扭矩信号为 0 或异常，因此仅显示低可信状态。',
            power_to_tire_traction_over_power_explanation: '驱动轮滑移较高且加速度弱或不稳定，说明轮上扭矩超过可用牵引。',
            traction_over_power_explanation: '驱动轮滑移较高且加速度弱或不稳定，说明轮上扭矩超过可用牵引。',
            power_to_tire_rpm_below_explanation: '高油门、低滑移但加速弱，指向转速低于有效区间或齿比偏长。',
            rpm_below_useful_range_explanation: '高油门、低滑移但加速弱，指向转速低于有效区间或齿比偏长。',
            power_to_tire_rpm_high_explanation: '转速长期接近上限且加速衰减，当前挡位可能偏短或换挡过晚。',
            rpm_too_high_or_gear_short_explanation: '转速长期接近上限且加速衰减，当前挡位可能偏短或换挡过晚。',
            power_to_tire_not_reaching_ground_explanation: '动力已有请求，但速度和 G 力响应偏弱，应先检查牵引、路面和动力释放，再判断齿比。',
            power_not_reaching_ground_explanation: '动力已有请求，但速度和 G 力响应偏弱，应先检查牵引、路面和动力释放，再判断齿比。',
            power_to_tire_ok_explanation: '当前窗口中功率、转速和驱动轮滑移没有显示明显动力落地瓶颈。',
            power_landing_ok_explanation: '当前窗口中功率、转速和驱动轮滑移没有显示明显动力落地瓶颈。',
        },
        powerToTireDrivenAxleLabels: {
            front: '前驱动轮',
            rear: '后驱动轮',
            all: '四驱驱动轮',
            unknown: '未知驱动轮',
        },
        brakeToTireTitle: '刹车输入 → 轮胎抓地',
        brakeToTireSubtitle: '实验性 DO 包制动抓地检查；速度变化为主证据，原始 G 力为辅助证据。',
        brakeToTireStatus: '制动抓地状态',
        brakeToTireBrake: '刹车',
        brakeToTireHandbrake: '手刹',
        brakeToTireSpeed: '速度 / 变化',
        brakeToTireSteer: '方向',
        brakeToTireDecel: '减速 G',
        brakeToTirePlaneG: '原始 X/Z G',
        brakeToTireFrontSlip: '前轮制动滑移 P90',
        brakeToTireRearSlip: '后轮制动滑移 P90',
        brakeToTireFrontCombined: '前轮综合滑移 P90',
        brakeToTireRearCombined: '后轮综合滑移 P90',
        brakeToTireSamples: '制动样本',
        brakeToTireTrail: '带刹转向',
        brakeToTireHandbrakeActive: '手刹激活',
        brakeToTireSummaryLabels: {
            brake_to_tire_no_data: '暂无制动抓地数据',
            brake_to_tire_insufficient: '制动样本不足',
            brake_to_tire_low_brake: '刹车输入不足',
            front_brake_lock_tendency: '前轮抱死倾向',
            rear_brake_lock_tendency: '后轮抱死倾向',
            trail_brake_front_overload: '带刹入弯前轮超载',
            handbrake_rear_slide: '手刹导致后轮滑移',
            brake_not_slowing_effectively: '刹车没有有效减速',
            brake_landing_ok: '制动抓地状态可用',
        },
        brakeToTireExplanationLabels: {
            brake_to_tire_waiting_for_samples: '等待包含刹车、速度、轮胎滑移和 G 力证据的 DO 样本。',
            brake_to_tire_need_brake_samples: '需要更多制动样本后再判断刹车抓地。',
            brake_to_tire_low_brake_explanation: '刹车和手刹输入都不足，不能形成刹车到轮胎的判断。',
            brake_to_tire_front_lock_explanation: '制动时前轮滑移率占主导，说明前轮纵向抓地被过度消耗。',
            brake_to_tire_rear_lock_explanation: '制动时后轮滑移率占主导，说明后轴制动稳定性不足。',
            brake_to_tire_trail_brake_front_overload_explanation: '刹车和转向叠加时前轮综合滑移较高，指向入弯前轮超载。',
            brake_to_tire_handbrake_rear_slide_explanation: '手刹输入激活且后轮滑移升高，应视为驾驶者主动制造后轴旋转的证据。',
            brake_to_tire_not_slowing_effectively_explanation: '刹车输入较高，但基于速度变化的减速度偏弱，轮胎也没有明显抱死。',
            brake_to_tire_ok_explanation: '当前窗口中刹车输入、减速度和前后轮滑移没有显示明显制动瓶颈。',
        },
        tuneInfluenceMap: '调校参数影响图',
        tuneInfluenceMapHint: '只读解释层：显示调校字段如何影响四条轮胎。',
        tuneInfluenceNoData: '影响图尚未加载。',
        tuneInfluenceButton: '影响',
        tuneInfluenceModalTitle: '调校字段影响',
        tuneInfluenceType: '影响类型',
        tuneInfluenceScope: '影响范围',
        tuneInfluencePhase: '主要阶段',
        tuneInfluenceMetrics: '轮胎指标',
        tuneInfluenceEvidence: '遥测证据',
        tuneInfluenceSideEffects: '常见副作用',
        tuneInfluenceConditions: '适用条件',
        tuneInfluenceTypeLabels: {
            direct: '直接',
            indirect: '间接',
        },
        tuneInfluenceCategoryLabels: {
            tire: '轮胎',
            gearing: '齿比',
            alignment: '轮胎定位',
            antiroll: '防倾杆',
            springs: '弹簧 / 车高',
            damping: '阻尼',
            aero: '空气动力',
            brake: '刹车',
            differential: '差速器',
        },
        tuneInfluenceScopeLabels: {
            front_axle: '前轴',
            rear_axle: '后轴',
            all_wheels: '四轮',
            driven_wheels: '驱动轮',
            left_right_balance: '左右差证据',
        },
        tuneInfluencePhaseLabels: {
            all_phases: '全部阶段',
            launch: '起步',
            braking: '制动',
            corner_entry: '入弯',
            sustained_cornering: '持续转向',
            corner_exit: '出弯',
            straight_power: '直线加速',
            high_speed_corner: '高速弯',
            transition: '重心转移',
            kerb_impact: '路肩 / 颠簸',
            coast: '松油滑行',
        },
        tuneInfluenceMetricLabels: {
            tire_temp: '胎温',
            combined_slip: '综合滑移',
            slip_angle: '滑移角',
            slip_ratio: '滑移率',
            suspension_offset: '悬挂偏距',
            g_force: 'G 力',
            yaw_response: '车身旋转响应',
            speed_rpm: '速度 / 转速',
            wheel_torque: '轮上扭矩',
            brake_slip: '制动滑移',
        },
        tuneInfluenceSideEffectLabels: {
            can_increase_tire_temp: '可能提高胎温',
            can_reduce_stability: '可能降低稳定性',
            can_affect_acceleration: '影响加速',
            can_affect_top_speed: '影响极速',
            can_mask_camber_issue: '可能掩盖外倾角证据',
            can_reduce_opposite_axle_grip: '可能降低另一轴抓地',
            can_create_understeer: '可能制造推头',
            can_create_oversteer: '可能制造甩尾',
            can_increase_lockup_risk: '可能增加抱死风险',
        },
        tuneInfluenceConditionLabels: {
            high_load_only: '需要有负载样本',
            throttle_sensitive: '受油门影响',
            drivetrain_sensitive: '受传动形式影响',
            unlocked_only: '需要已解锁调校项',
            speed_sensitive: '受速度影响',
            aero_speed_sensitive: '高速空力相关',
            brake_sensitive: '受刹车输入影响',
        },
        gForceDiagnostics: 'G 力诊断',
        gForceCurrent: '当前 G',
        gForceAverage: '平均 G',
        gForcePeak: '峰值 G',
        gForceTotal: '合成 G',
        gForcePlane: 'X/Z 平面 G',
        gForceCircleScale: '外圈',
        gForceDominantAxis: '主导轴',
        gForceAxisMapping: '轴向映射',
        gForceAxisMappingLabels: {
            raw_packet_axes_unverified: '原始 AccelerationX/Y/Z 轴；圆图使用 X/Z，且 X 方向已反转。',
        },
        gForceChart: '实时 G 力图',
        camberInference: '外倾角推断',
        camberFront: '前轮外倾角',
        camberRear: '后轮外倾角',
        camberCorneringSamples: '过弯样本',
        camberStateLabels: {
            unknown: '未知',
            stable: '稳定',
            monitor: '观察',
            likely_needs_more_negative: '可能需要更多负外倾',
            platform_limited: '先看平台',
            thermal_limited: '先看胎温',
        },
        camberSummaryLabels: {
            camber_inference_insufficient: '持续过弯样本不足。',
            camber_inference_front_needs_more_negative: '前轴可能需要更多负外倾。',
            camber_inference_rear_needs_more_negative: '后轴可能需要更多负外倾。',
            camber_inference_both_axles_need_more_negative: '前后轴都可能需要更多负外倾。',
            camber_inference_platform_first: '应先解决平台/悬挂问题，再判断外倾角。',
            camber_inference_temperature_first: '应先解决胎温问题，再判断外倾角。',
            camber_inference_monitor: '当前窗口中外倾角不是主导证据。',
        },
        camberExplanationLabels: {
            camber_inference_needs_cornering: '需要采集有明显转向的持续过弯样本后再判断外倾角。',
            camber_inference_slip_angle_explanation: '因为 Data Out 每条轮胎只有一个胎温，本判断使用滑移角、综合滑移和 G 负载作为低置信度外倾角证据。',
            camber_inference_platform_explanation: '悬挂/平台问题会伪装成外倾角问题，因此应先检查平台稳定性。',
            camber_inference_temperature_explanation: '胎温已在限制抓地，应先检查胎压/热平衡再看外倾角。',
            camber_inference_monitor_explanation: '当前窗口的滑移角和综合滑移没有强烈指向外倾角。',
        },
        tireLabLimitLabels: {
            unknown: '未知',
            no_dynamic_load: '无动态负载',
            stationary: '静止',
            balanced: '平衡',
            balanced_near_limit: '接近极限但较平衡',
            front_limited: '前轮受限',
            rear_limited: '后轮受限',
            four_wheel_limited: '四轮受限',
            traction_limited: '驱动轮牵引受限',
            thermal_limited: '胎温受限',
            platform_limited: '平台/悬挂受限',
        },
        tireLabPhaseLabels: {
            unknown: '未知',
            stationary: '静止',
            handbrake: '手刹',
            launch: '起步',
            braking: '制动',
            corner_entry: '入弯',
            low_speed_corner: '低速弯',
            mid_speed_corner: '中速弯',
            sustained_cornering: '持续转向',
            corner_exit: '出弯',
            straight_power: '直线加速',
            high_speed_corner: '高速弯',
        },
        tireLabGripStateLabels: {
            unknown: '未知',
            stable: '稳定',
            warning: '警告',
            limit: '到达极限',
        },
        tireLabSummaryLabels: {
            tire_model_no_data: '暂无遥测样本。',
            tire_model_stationary: '车辆处于静止状态，暂停轮胎极限判断。',
            tire_model_no_dynamic_load: '当前窗口没有足够的动态轮胎负载样本。',
            tire_model_balanced: '当前四条轮胎没有明显单一抓地瓶颈。',
            tire_model_balanced_near_limit: '车辆已接近抓地极限，但前后负载相对平衡。',
            tire_model_front_limited: '前轴先到达横向抓地极限。',
            tire_model_front_brake_limited: '带刹入弯时前轴负载过高。',
            tire_model_rear_limited: '后轴先到达横向抓地极限。',
            tire_model_rear_power_limited: '给油时后轴负载过高。',
            tire_model_rear_handbrake_limited: '手刹输入下后轴发生滑移。',
            tire_model_rear_traction_limited: '后驱动轮纵向滑移过高。',
            tire_model_front_traction_limited: '前驱动轮纵向滑移过高。',
            tire_model_drive_traction_limited: '驱动轮纵向滑移过高。',
            tire_model_four_wheel_limited: '四条轮胎都接近或超过抓地极限。',
            tire_model_thermal_limited: '胎温或前后胎温差正在限制抓地。',
            tire_model_platform_limited: '悬挂行程或车身平台控制可能限制轮胎贴地。',
        },
        tireLabExplanationLabels: {
            tire_model_waiting_for_samples: '等待有效 324-byte 遥测样本。',
            tire_model_stationary_explanation: '速度和驾驶输入接近 0，此时滑移、胎温和悬挂只作为遥测显示，不触发轮胎警告。',
            tire_model_no_dynamic_load_explanation: '当前窗口缺少足够的速度、转向、制动、油门或 G 力负载，无法判断轮胎极限。',
            tire_model_balanced_explanation: '当前窗口内没有某一组轮胎明显主导滑移证据。',
            tire_model_balanced_near_limit_explanation: '前后轴都在高负载工作，建议先降低速度或采集更干净样本再调校。',
            tire_model_front_explanation: '前轮综合滑移高于后轮，说明当前瓶颈是前轮抓地。',
            tire_model_front_brake_explanation: '制动和转向叠加使前轮负载超过可用抓地。',
            tire_model_rear_explanation: '后轮综合滑移高于前轮，说明当前瓶颈是后轴稳定性。',
            tire_model_rear_power_explanation: '给油时后轮滑移升高，说明后轴抓地或动力输出是瓶颈。',
            tire_model_rear_handbrake_explanation: '手刹输入激活时后轮滑移升高，这类样本应视为驾驶者主动制造旋转。',
            tire_model_rear_traction_explanation: '给油时后轮滑移率主导，轮上扭矩超过后轮牵引能力。',
            tire_model_front_traction_explanation: '给油时前轮滑移率主导，轮上扭矩超过前轮牵引能力。',
            tire_model_drive_traction_explanation: '给油时驱动轮滑移率过高，应先解决牵引再判断齿比。',
            tire_model_four_wheel_explanation: '前后轮都接近极限，更像总抓地、速度、空力或平台问题，而不是单独某一轴。',
            tire_model_thermal_explanation: '胎温过高或前后胎温差过大会降低可用抓地。',
            tire_model_platform_explanation: '悬挂行程接近耗尽，可能破坏轮胎贴地并造成抓地不稳定。',
        },
        tireLabWarningLabels: {
            tire_model_no_data: '暂无轮胎模型样本。',
            tire_model_sample_insufficient: '样本数量偏少，当前结论仅作临时参考。',
            sample_insufficient: '样本数量偏少，当前结论仅作临时参考。',
            thermal_risk: '存在胎温风险，但没有滑移证据时不视为轮胎极限。',
            platform_risk: '存在悬挂/平台风险，但没有滑移证据时不视为轮胎极限。',
            left_right_imbalance: '左右轮滑移差异较大。这里只作为证据显示，不生成单侧调校动作。',
            tire_model_left_right_imbalance: '左右轮滑移差异较大。这里只作为证据显示，不生成单侧调校动作。',
            g_force_axis_mapping_unverified: 'G 力使用原始数据包 X/Y/Z 轴，纵向/横向/垂向映射仍需实测校准。',
            camber_inference_no_three_point_temps: '外倾角判断为间接推断，因为 Data Out 每条轮胎只有一个胎温。',
        },
        tireLabHintDirections: {
            improve_front_grip_or_reduce_entry_load: '提升前轴抓地，或降低入弯负载。',
            improve_rear_grip_or_smooth_rotation: '提升后轴抓地，或让车尾响应更平顺。',
            reduce_wheel_torque_or_improve_driven_tire_grip: '降低驱动轮轮上扭矩，或提升驱动轮抓地。',
            reduce_speed_or_increase_total_grip: '降低入弯速度，或提升整体抓地/平台支撑。',
            bring_tires_back_to_temperature_window: '让胎温回到可用窗口。',
            restore_suspension_travel_and_platform_control: '恢复悬挂行程和车身平台控制。',
            collect_moving_tire_load_samples: '开始行驶后采集有负载的轮胎样本。',
            collect_more_representative_corner_and_power_samples: '继续采集更有代表性的弯道和动力样本。',
            consider_more_negative_front_camber: '仅在可重复过弯样本验证后，考虑略微增加前轮负外倾。',
            consider_more_negative_rear_camber: '仅在可重复过弯样本验证后，考虑略微增加后轮负外倾。',
            use_camber_inference_as_low_confidence_evidence: '把外倾角推断作为低置信度证据，不作为直接调校命令。',
        },
        tireLabHintLabels: {
            front_axle_grip: '前轴抓地',
            rear_axle_grip: '后轴抓地',
            driven_tire_traction: '驱动轮牵引',
            whole_car_grip: '整车抓地',
            tire_temperature: '胎温',
            platform_stability: '平台稳定性',
            observe: '继续观察',
            front_camber_check: '前轮外倾角检查',
            rear_camber_check: '后轮外倾角检查',
            camber_observe: '外倾角观察',
        },
        fieldDiagnostics: '遥测字段诊断',
        speedCalibration: '速度校准',
        vehicleMetadata: '车辆元数据',
        enginePower: '发动机 / 动力',
        motionPose: '运动 / 姿态',
        raceLapData: '比赛 / 圈速',
        auxiliaryFields: '辅助字段',
        fieldName: '字段',
        fieldValue: '数值',
        fieldUnit: '单位',
        fieldSource: '来源',
        fieldRange: '期望范围',
        fieldState: '状态',
        ok: '正常',
        checkValue: '检查',
        noCurrentFrame: '暂无当前遥测帧。',
        replayTimeline: '回放时间轴',
        pauseReplay: '暂停',
        resumeReplay: '继续',
        replayPosition: '位置',
        sessionCompare: '会话对比',
        sessionMode: '会话模式',
        modeCompareWarning: '这两个会话的遥测模式不同，请谨慎比较调校事件。',
        advancedSettings: '高级设置',
        testConditions: '测试条件',
        restoreUnknown: '恢复未知',
        driverMode: '驾驶识别',
        brakeAssist: '刹车',
        steeringAssist: '转向',
        tractionControl: 'TCS',
        stabilityControl: 'STM',
        shifting: '换挡',
        launchControl: '起跑控制',
        assists: '辅助设置',
        comparabilityWarnings: {
            game_mode_mismatch: '遥测模式不同，该对比只能作为参考。',
            test_conditions_unknown: '一个或两个会话存在未知测试条件，对比可信度降低。',
            driver_mode_mismatch: '驾驶方式不同。',
            brake_assist_mismatch: '刹车辅助设置不同。',
            steering_assist_mismatch: '转向设置不同。',
            traction_control_mismatch: '牵引力控制设置不同。',
            stability_control_mismatch: '稳定控制设置不同。',
            shifting_mismatch: '换挡设置不同。',
            launch_control_mismatch: '起跑控制设置不同。',
        },
        testConditionValues: {
            unknown: '未知',
            player: '玩家',
            auto: '自动驾驶',
            assisted: '辅助',
            abs_on: 'ABS 开',
            abs_off: 'ABS 关',
            standard: '标准',
            simulation: '拟真',
            on: '开启',
            off: '关闭',
            automatic: '自动',
            manual: '手动',
        },
        leftSession: '左侧会话',
        rightSession: '右侧会话',
        compare: '对比',
        metric: '指标',
        left: '左侧',
        right: '右侧',
        delta: '差值',
        eventDistribution: '事件分布',
        profileCompare: '档案版本对比',
        openProfileCompare: '版本对比',
        profileA: '档案 A',
        profileB: '档案 B',
        changedFields: '变化项',
        recentChanges: '最近修改',
        noRecentChanges: '暂无保存的修改记录。',
        expertWorkspace: '专业调校工作区',
        expertStartHint: '编辑调校档案，开启实时遥测后查看本轮内存中的车辆问题、决策和解释器建议；不保存会话或录制。',
        professionalDiagnosticTitle: '实时调校分析',
        professionalDiagnosticHint: '使用开发者模式中选择的检测模型、决策器和解释器。',
        professionalDiagnosticEmpty: '启动专业调校遥测后显示实时车辆问题与修改建议。',
        professionalMergedTitle: '合并调校诊断',
        professionalMergedHint: '按问题合并显示问题、决策、解释器建议和对应修改值。',
        professionalMergedEmpty: '暂无合并诊断结果。',
        professionalMergedProblem: '问题',
        professionalMergedDecision: '决策',
        professionalMergedAdvice: '建议',
        professionalMergedAdjustments: '对应修改值',
        professionalMergedNoAdjustments: '暂无具体修改值。请先填写相关调校字段，或仅按方向建议处理。',
        professionalTuneInstruction: '怎么调',
        professionalChangeAmount: '改多少',
        selectProfileToEdit: '请从左侧选择调校档案，或新建一个档案。',
        newProfileModalTitle: '新建调校档案',
        newProfileModalHint: '先创建基础档案，再在专家调校中补全详细调校数值。',
        createAndEdit: '创建并编辑',
        fillTelemetryIntoDraft: '从遥测填充车辆身份',
        snapshotChangedCount: (count: number) => `${count} 个变化项`,
        compareWithCurrent: '与当前设置对比',
        restoreSnapshot: '恢复',
        snapshotRestored: '已恢复修改快照。',
        restoreSnapshotConfirm: '确认把该快照恢复为当前调校档案？',
        profileSnapshotCompare: '当前设置 vs 最近修改',
        snapshotBefore: '快照',
        currentSettings: '当前设置',
        moreActions: '更多操作',
        noChanges: '没有变化项。',
        close: '关闭',
        chooseProfile: '选择调校档案',
        profileChoiceTitle: '选择本次测试档案',
        profileChoiceHint: '当前车辆匹配到多个调校档案，请选择本次测试使用的档案。',
        noMatchingProfile: '当前车辆暂未匹配到调校档案，将以未绑定档案开始监听。',
        profileMatchUnavailable: '暂无当前遥测帧，无法校验车辆匹配。',
        profileMismatchTitle: '遥测车辆不匹配',
        profileMismatchHint: '当前选择的调校档案与遥测车辆不一致。请选择匹配档案、未选择调校档案，或取消。',
        telemetryVehicle: '遥测车辆',
        currentTuneVehicle: '当前调校档案',
        chooseMatchingProfile: '选择匹配档案',
        clearProfileAndStart: '未选择调校档案',
        profileSessions: '测试会话',
        recentSession: '最近测试',
        ruleThresholds: '规则阈值',
        strategyTemplates: '策略模板',
        strategyAnalysis: '五会话策略分析',
        strategyTemplate: '策略模板',
        selectedSessions: '已选会话',
        runStrategyAnalysis: '运行分析',
        strategyAnalysisEmpty: '选择最多 5 个会话和一个策略模板，用于分析规则匹配。',
        strategyRecommendation: '策略建议',
        strategyHints: '分析提示',
        enabledEvents: '启用事件',
        totalEvents: '事件总数',
        ruleProfiles: '阈值配置',
        ruleName: '名称',
        ruleCarClass: '等级匹配',
        ruleDrivetrain: '传动匹配',
        ruleUseCase: '用途匹配',
        ruleConfigJson: '配置 JSON',
        resetDefaults: '恢复默认',
        createProfile: '新增档案',
        updateProfile: '保存修改',
        newProfile: '新建',
        duplicate: '复制',
        delete: '删除',
        deleteSession: '删除会话',
        deleteSessionConfirm: (name: string) => `删除测试会话“${name}”及其回放录制？`,
        deleteSessionBlocked: '请先停止回放或监听，再删除会话。',
        setActive: '设为当前',
        active: '当前',
        profileList: '档案列表',
        profileForm: '调校档案',
        fillFromTelemetry: '从遥测填充',
        telemetryFillUnavailable: '当前没有可用于填充的遥测帧。',
        telemetryFilled: '已填充遥测车辆字段。',
        requiredCarName: '车辆名称必填。',
        saved: '已保存',
        saveAction: '保存',
        deleted: '已删除',
        reportSessions: '测试会话',
        generateReport: '生成报告',
        noSessions: '暂无已保存会话。',
        noProfiles: '暂无调校档案。',
        profileIdentity: '车辆身份',
        profileIdentityHint: '用于匹配遥测会话、测试报告和调校版本。',
        profileTelemetryMatch: '与当前遥测车辆匹配。',
        profileTelemetryMismatch: '与当前遥测车辆不匹配。',
        profileTelemetryUnavailable: '暂无当前遥测车辆可用于对比。',
        profileEditorActions: '档案操作',
        markdownReport: 'Markdown 报告',
        reportPlaceholder: '选择一个已保存会话生成报告。',
        reportDecisionTitle: '调校决策报告',
        reportStatus: '报告状态',
        issueAdvice: '问题项与建议',
        wholeCarPlan: '整车调校方案',
        wholeCarPlanEmpty: '暂无整车调校方案。请绑定调校档案并采集可对比样本。',
        roadTuningDecision: '公路调校决策',
        roadTuningDecisionEmpty: '暂无明确公路调校决策。请使用可对比的公路测试和遥测证据。',
        primaryIssue: '主问题',
        primaryCause: '主因',
        cornerPhase: '阶段',
        driverFitVerdict: '驾驶适配',
        rollbackRecommended: '建议先回退',
        rollbackRecommendedHint: '上一次相关修改后对比结果变差，先回退或反向微调，再继续加新改动。',
        retestFocus: '复测观察点',
        knowledgeSource: '规则来源',
        knowledgeFallback: '备用规则',
        autoApplicable: '可应用',
        manualCheck: '人工确认',
        optional: '可选',
        knowledgeStatus: '公路模型规则',
        reloadKnowledge: '重载规则',
        knowledgeReloaded: '调校知识已重新加载。',
        knowledgeSymptoms: '症状',
        knowledgeActions: '动作',
        roadDecisionStatusLabels: {
            ready: '可用',
            rollback_recommended: '优先回退',
            no_matching_symptom: '无匹配症状',
            insufficient_data: '数据不足',
            profile_unbound: '未绑定档案',
            knowledge_error: '规则加载错误',
        },
        roadActionRoles: {
            primary: '主调整',
            support: '辅助调整',
            alternative: '备选调整',
        },
        roadPhaseLabels: {
            launch: '起步',
            braking: '制动',
            corner_entry: '入弯',
            mid_corner: '持续过弯',
            corner_exit: '出弯',
            cornering: '过弯',
            high_speed: '高速',
            platform: '车身平台',
            tires: '轮胎',
            power: '动力',
        },
        driverFitVerdictLabels: {
            unknown: '未知',
            improved: '遥测改善',
            worsened: '遥测恶化',
            unchanged: '基本持平',
            insufficient_data: '数据不足',
        },
        retestFocusLabels: {
            same_car: '同车',
            same_track: '同赛道',
            same_driver_mode: '同驾驶方式',
            verify_rollback_before_new_changes: '先验证回退',
            gear_power_window: '齿比动力窗口',
            road_launch_wheelspin: '起步打滑',
            road_launch_bog_down: '起步憋转',
            road_entry_understeer: '入弯推头',
            road_entry_understeer_low_speed: '低速入弯推头',
            road_entry_understeer_mid_speed: '中速入弯推头',
            road_entry_understeer_high_speed: '高速入弯推头',
            road_mid_understeer: '持续过弯推头',
            road_power_understeer: '给油推头',
            road_exit_oversteer: '出弯甩尾',
            road_lift_snap_oversteer: '松油甩尾',
            road_front_brake_lockup: '前轮抱死',
            road_rear_brake_lockup: '后轮抱死',
            road_high_speed_slide: '高速四轮侧滑',
            road_bottom_out: '悬挂触底',
            road_tire_overheat: '胎温异常',
        },
        roadSymptomLabels: {
            road_launch_wheelspin: '起步打滑',
            road_launch_bog_down: '起步憋转',
            road_entry_understeer: '入弯推头',
            road_entry_understeer_low_speed: '低速入弯推头',
            road_entry_understeer_mid_speed: '中速入弯推头',
            road_entry_understeer_high_speed: '高速入弯推头',
            road_mid_understeer: '持续过弯推头',
            road_power_understeer: '给油推头',
            road_exit_oversteer: '给油甩尾',
            road_lift_snap_oversteer: '松油甩尾',
            road_front_brake_lockup: '前轮抱死',
            road_rear_brake_lockup: '后轮抱死',
            road_high_speed_slide: '高速四轮侧滑',
            road_bottom_out: '悬挂触底',
            road_tire_overheat: '胎温过热 / 温差异常',
            road_gearing_power: '齿比动力窗口',
        },
        tunePlanDraft: '调校方案草稿',
        tunePlanDraftEmpty: '该会话暂无可应用的调校方案动作。',
        applyTunePlan: '应用所选项到调校档案',
        tunePlanApplied: '调校方案已应用到绑定档案。',
        tunePlanStatusLabels: {
            ready: '可应用',
            no_actions: '无可应用动作',
            profile_unbound: '未绑定调校档案',
            vehicle_mismatch: '车辆不匹配',
            cannot_verify_vehicle: '无法校验车辆匹配',
        },
        tunePlanBlockedReasons: {
            field_locked_or_blank: '字段为空或锁定',
            vehicle_mismatch: '车辆不匹配',
            cannot_verify_vehicle: '无法校验车辆匹配',
            profile_unbound: '未绑定调校档案',
            no_numeric_adjustment: '无数值调整',
            manual_review_required: '需要人工确认',
            no_change: '无变化',
            rollback_first: '复测变差，先应用回退建议。',
            duplicate_action_removed: '已移除重复建议',
            same_field_direction_conflict: '已消解同字段方向冲突',
        },
        tunePlanTrust: '可信度',
        tunePlanMissingInputs: '缺失条件',
        tunePlanRetestGuard: '复测保护',
        tunePlanTrustLevels: {
            high: '高可信',
            medium: '中可信',
            low: '低可信',
            blocked: '已阻止',
        },
        tunePlanTrustReasons: {
            low_retest_confidence: '复测可信度较低',
            retest_worsened: '复测结果变差',
            recent_tune_plan_apply: '来自上次应用的调校方案',
            missing_current_value: '当前值缺失',
            model_confidence_high: '模型高可信',
            model_confidence_medium: '模型中可信',
            model_confidence_low: '模型低可信',
            source_gear_power_diagnostic: '齿比动力诊断',
            source_road_tuning_model: '公路调校模型',
            source_local_rule_report: '本地规则报告',
        },
        tunePlanMissingInputLabels: {
            tune_profile: '调校档案',
            current_tune_value: '当前调校值',
            profile_power_band: '动力区间 RPM 输入',
            high_load_samples: '高油门样本',
        },
        retestConfidence: '复测可信度',
        retestBaselineReason: '对比基线',
        retestChangedFields: '最近调校修改',
        retestRollbackActions: '回退建议',
        retestResult: '复测结果',
        retestEmpty: '暂无可对比的上一轮测试。',
        retestMetricLabels: {
            issue_score: '问题总分',
            event_count: '事件数量',
            event_duration_ms: '事件持续',
            avg_speed_kmh: '平均速度',
            max_speed_kmh: '最高速度',
            best_run_duration_ms: '最佳赛段时间',
            risk_score: '风险分',
            front_tire_temp: '前胎温',
            rear_tire_temp: '后胎温',
            gear_problem_count: '齿比动力问题',
        },
        retestStatusLabels: {
            improved: '改善',
            worsened: '恶化',
            unchanged: '持平',
            insufficient_data: '数据不足',
        },
        retestBaselineReasons: {
            matched_profile_track_driver: '同档案、同赛道、同驾驶方式',
            matched_vehicle_class_usecase_driver: '同车同级同用途和驾驶方式',
            missing_comparison_baseline: '没有可比上一轮',
            unavailable: '不可用',
        },
        planStrategy: '策略',
        planConfidence: '可信度',
        gearPowerDiagnostic: '齿比动力诊断',
        gearPowerDiagnosticHint: '比赛或漫游都可诊断，只要有足够的高油门、转速、挡位、速度和滑移样本。',
        gearPowerWhyNoAdvice: '为什么没有齿比建议',
        gearPowerNeedSamples: '需要先采集更多加速样本，再判断齿比。',
        gearPowerNeedHighLoad: '需要更多干净的高油门样本，并且样本来自已解锁/已填写的挡位。',
        gearPowerNoUnlockedGears: '当前专家档案中没有已填写或已解锁的齿比。',
        gearPowerFallbackLowConfidence: '当前动力区间使用转速比例回退。填写峰值扭矩转速、峰值功率转速和红线转速后，建议会更可靠。',
        gearPowerTractionFirst: '当前首先受牵引限制，应先检查差速器或抓地力，再改齿比。',
        gearPowerNoAdvice: '当前样本不建议修改齿比。',
        gearStrategyMode: '齿比策略',
        gearStrategyIssueCount: '问题挡位',
        quickGearAdviceReadOnly: '快速诊断只给方向。需要具体数值和一键应用时，请在专家调校中绑定档案。',
        gearPowerComparisons: '齿比对比',
        gearTelemetryComparison: '实测齿比表现',
        gearTuneComparison: '齿比设置变化',
        gearComparisonBefore: '对比前',
        gearComparisonAfter: '当前',
        gearComparisonDelta: '变化',
        gearComparisonUnavailable: '暂无可对比的齿比数据。',
        gearComparisonStatuses: {
            ready: '可用',
            missing_baseline: '缺少上一轮可比会话',
            no_matching_gears: '没有匹配挡位样本',
            no_changed_gears: '最近没有齿比修改',
            profile_unbound: '未绑定调校档案',
        },
        gearPowerSummary: '动力区间',
        powerBandTarget: '目标转速区间',
        powerBandSource: '区间来源',
        diagnosticConfidence: '可信度',
        speedRange: '速度范围',
        highestObserved: '最高实测',
        inPowerBand: '动力区间内',
        inPowerBandCoverage: '动力区间覆盖率',
        acceleration: '加速度',
        shiftAfter: '升挡后转速',
        powerGearTest: '动力 / 齿比测试',
        powerGearTestHint: '为了稳定判断齿比，请每个已解锁挡位做一次干净的全油门拉速，再在这里查看样本质量与升挡转速。',
        powerToWeight: '功率比重',
        tractionLimited: '牵引受限',
        gearFinding: '判断',
        planConflicts: '已解决冲突',
        planStrategies: {
            rollback_first: '优先回退',
            coarse_whole_car: '整车大步调整',
            targeted_whole_car: '整车定向调整',
        },
        planSummaries: {
            whole_car_template: '先执行影响最大的组合修改，再根据下一轮结果微调。',
            rollback_before_more_changes: '最近相关修改后结果变差，先回退，再继续调整。',
            no_clear_whole_car_action: '该会话暂未形成明确整车修改方案。',
        },
        planConfidenceLabels: {
            high: '高',
            medium: '中',
            low: '低',
            needs_profile: '需要绑定调校档案',
        },
        gearFindings: {
            not_enough_samples: '样本不足',
            not_enough_high_load: '高油门样本不足',
            no_unlocked_gear_samples: '没有已解锁挡位样本',
            gearing_window_ok: '齿比动力区间正常',
            gearing_adjustment_needed: '需要调整齿比动力区间',
            traction_limited_power: '动力输出受牵引限制',
            global_too_long: '整套齿比偏长，优先调整终传比',
            global_too_short: '整套齿比偏短，优先调整终传比',
            single_gear_too_long: '个别挡位偏长',
            single_gear_too_short: '个别挡位偏短',
            single_gear_mixed: '个别挡位问题混合',
            traction_limited_low_gears: '牵引受限，先解决抓地',
            top_speed_limited: '优先处理极速齿比',
            ok: '正常',
            too_long: '负载下齿比偏长',
            too_short: '负载下齿比偏短',
            traction_limited: '牵引受限',
            top_speed_limited_by_gearing: '极速受齿比限制',
            top_speed_bog_down: '高挡齿比偏长',
            top_speed_ok: '高速齿比正常',
            launch_wheelspin: '起步打滑',
            launch_bog_down: '起步憋转',
        },
        powerBandSources: {
            profile_power_band: '档案 RPM 输入',
            telemetry_engine_max_rpm: '遥测最高转速回退',
            rpm_ratio_fallback: '转速比例回退',
        },
        cornerOperationStateLabels: {
            '1': '带刹 / 制动',
            '2': '松油滑行',
            '3': '轻油维持',
            '4': '给油',
        },
        issueGroups: '合并问题组',
        noIssueGroups: '该会话暂无合并问题组。',
        issueGroupAdviceTitle: '问题组调校建议',
        issueGroupEvents: '组内事件',
        issueGroupEvidence: '证据范围',
        issueGroupPrimaryAdvice: '主要调校建议',
        issueStrategy: '调整策略',
        issueStrategyLabels: {
            rollback_first: '先回退相关修改',
            coarse_combination: '大步组合调整',
            medium_combination: '中等组合调整',
            fine_tune: '小步微调',
        },
        feedbackDirectiveLabels: {
            rollback_related_changes: '相关的上次修改后问题变差，先回退部分修改，不继续同方向加码。',
            keep_direction_then_fine_tune: '当前方向已有改善，后续改用更小步进微调。',
            avoid_more_same_direction: '问题没有明确改善，先避免继续同方向调整，建议复测确认。',
        },
        issueGroupComparison: '相对上次测试',
        issueBaseline: '对比基线',
        issueRecentChanges: '最近调校修改',
        concreteProfileRequired: '需绑定调校档案后生成具体调整量。',
        noReportIssues: '该会话没有保存的问题事件。',
        profileBoundStatus: '已绑定调校档案',
        profileUnboundStatus: '未绑定调校档案',
        sessionProfileUnboundHint: '请绑定匹配的调校档案，让报告使用正确的设置上下文。',
        driverModeUnknownReportHint: '驾驶识别置信度不足，自动驾驶基线和玩家适配结论会受限。',
        baselineMissingReportHint: '缺少匹配的自动驾驶基线。请用同车、同级、同标准赛段录制一次自动驾驶。',
        standardSegmentMissingHint: '未检测到有效标准赛段。请先创建或匹配标准赛道。',
        advancedReportDetails: '高级报告详情',
        expand: '展开',
        collapse: '收起',
        playbackAndTimeline: '回放时间轴',
        rawMarkdown: 'Markdown 原文',
        bindSessionProfile: '绑定调校档案',
        changeSessionProfile: '更改调校档案',
        sessionProfileBindTitle: '绑定测试会话调校档案',
        sessionProfileBindHint: '只能选择与该会话车辆匹配的调校档案。',
        sessionVehicle: '会话车辆',
        matchingTuneProfiles: '匹配调校档案',
        noSessionProfileMatches: '没有与该会话车辆匹配的调校档案，请先从遥测新建或填充档案。',
        sessionProfileBound: '测试会话调校档案已更新，请重新生成报告以使用新档案。',
        benchmarkTracks: '标准赛道',
        trackProfilesTitle: '赛道档案',
        trackProfilesSubtitle: '采集赛道数据，并按车辆查看自动驾驶基线。',
        trackCaptureMode: '赛道采集',
        trackCaptureNoHistory: '赛道采集不保存测试会话、录制、样本或回放数据。',
        trackData: '赛道数据',
        vehicleReferences: '车辆参考信息',
        autoBaselines: '自动驾驶基线',
        noAutoBaselines: '该赛道暂无有效自动驾驶基线。',
        noVehicleReferences: '该赛道暂无车辆参考记录。',
        bestAutoBaseline: '最佳自动驾驶基线',
        baselineVehicle: '基线车辆',
        baselineRunCount: '基线次数',
        validRuns: '有效通过',
        autoRuns: '自动驾驶',
        recentBenchmarkRuns: '最近通过记录',
        routeCompletion: '路线完成度',
        baselineWarnings: '基线提示',
        similarTrackFound: '发现相似赛道',
        similarTrackHint: '本次路线与已有赛道接近。合并可保留已有车辆基线，也可以另存为新赛道。',
        mergeIntoExistingTrack: '合并到已有赛道',
        saveAsNewTrack: '另存为新赛道',
        startDistance: '起点距离',
        endDistance: '终点距离',
        shapeSimilarity: '形状相似度',
        routeFitAvgError: '路线拟合平均误差',
        routeFitP90Error: '路线拟合 P90 误差',
        matchLevel: '匹配等级',
        strongMatch: '强匹配',
        mediumMatch: '需确认匹配',
        autoMergedTrack: '赛道已自动合并',
        renameTrack: '重命名赛道',
        trackRenamed: '赛道已重命名。',
        trackBaselines: '车辆基线',
        trackBaselineCapture: '车辆基线采集',
        startTrackBaseline: '开始基线采集',
        saveTrackBaseline: '保存基线',
        trackBaselineSaved: '车辆基线已保存。',
        trackBaselineAutoMatched: '基线已保存到匹配赛道',
        trackBaselineAutoCreated: '基线已保存，并已创建新赛道',
        trackBaselineAutoArchiveHint: '基线采集会自动匹配本次路线到已有赛道；没有强匹配时会创建新赛道。',
        trackBaselineStopped: '车辆基线采集已停止。',
        trackBaselineDeleted: '车辆基线已删除。',
        trackBaselineNoSession: '车辆基线采集不保存测试会话或回放数据。',
        confirmDeleteBaseline: '删除这条车辆基线？',
        trackProfileWarnings: {
            no_auto_baseline: '该赛道尚未检测到有效自动驾驶基线。',
            baseline_vehicle_identity_missing: '部分自动驾驶记录缺少车辆 ID、等级或 PI，未纳入分组。',
        },
        trackBuilder: '赛道绘制器',
        trackName: '赛道名称',
        startCapture: '开始采集',
        stopCapture: '停止采集',
        saveTrack: '保存赛道',
        fromSession: '从会话生成',
        analyzeTrackRuns: '分析赛道通过',
        capturePoints: '采集点',
        routeLength: '路线长度',
        drivingLineSignal: '辅助线信号',
        detected: '已检测',
        notDetected: '未检测',
        noTracks: '暂无标准赛道。',
        noTrackPoints: '暂无路线采集点。',
        trackType: '赛道类型',
        autoTrackType: '自动',
        circuitTrack: '环道',
        sprintTrack: '点到点',
        extractionMode: '提取模式',
        autoBestLap: '自动最佳圈',
        firstLap: '第一圈',
        fullSegment: '全段',
        observedLaps: '观察圈数',
        startGate: '起点 Gate',
        finishGate: '终点 Gate',
        checkpoints: '检查点',
        setStartGate: '设置起点 Gate',
        setFinishGate: '设置终点 Gate',
        clearGates: '清除 Gate',
        reextractTrack: '重新提取赛道',
        trackSaved: '标准赛道已保存。',
        trackDeleted: '标准赛道已删除。',
        benchmarkRuns: '赛道通过记录',
        noBenchmarkRuns: '该会话暂未匹配到标准赛道。',
        roadEvaluation: '公路赛车评估',
        roadEvaluationEmpty: '暂无标准赛段评估。请先创建或匹配标准赛道。',
        insufficientData: '数据不足',
        paperPerformanceScore: '纸面性能',
        playerFitScore: '玩家适配',
        riskScore: '失控风险',
        autoBaseline: '自动驾驶基线',
        bestPlayerRun: '最佳玩家通过',
        missingAutoBaselineHint: '当前只可查看玩家表现。请使用同车、同级、同标准赛段录制一次自动驾驶，才能判断纸面基线。',
        prioritizeTuning: '优先调校',
        evaluationContext: '评估上下文',
        evaluationContextEmpty: '该问题项暂未关联到评估归因。',
        prioritizeTuningYes: '建议优先检查调校',
        prioritizeTuningNo: '先复查驾驶风格和重复稳定性，再决定是否改调校',
        roadVerdicts: {
            good_fit: '好车：快且可控',
            fast_but_risky: '速度快但风险高',
            paper_fast_not_fit: '纸面快但不适配',
            needs_tuning: '需要继续调校',
            insufficient_data: '数据不足',
        },
        roadBaselineStatuses: {
            matched_auto_baseline: '已匹配自动驾驶基线',
            self_auto_baseline: '本会话是自动驾驶基线',
            missing_auto_baseline: '缺少自动驾驶基线',
            missing_vehicle_identity: '缺少车辆身份',
            no_valid_standard_run: '没有有效标准赛段',
            no_standard_track: '没有标准赛道',
        },
        roadAttributions: {
            tune_issue: '调校问题',
            style_fit_issue: '驾驶风格适配',
            driver_execution_issue: '驾驶执行',
            data_gap: '数据缺口',
        },
        roadAttributionMessages: {
            event_pattern: '该问题重复出现，已影响公路评估。',
            route_deviation: '标准赛段检测到路线偏离。',
            route_progress_low: '本次通过没有完成足够的标准路线。',
            geometry_length_mismatch: '实测几何长度与保存的标准路线存在差异。',
            distance_traveled_mismatch: 'Data Out 距离字段与路线几何不一致。',
            missing_auto_baseline: '当前车辆和路线缺少匹配的自动驾驶基线。',
            missing_vehicle_identity: '会话缺少车辆快照，基线匹配能力受限。',
            no_valid_standard_run: '未检测到有效的标准赛段通过。',
            no_standard_track: '该会话暂无可用标准赛道。',
        },
        confidence: '置信度',
        bestRun: '最佳通过',
        eventsSaved: '保存事件',
        avgSpeed: '平均速度',
        maxSpeed: '最高速度',
        duration: '持续时间',
        routeProgress: '路线进度',
        sourceSession: '来源会话',
        updatedAt: '更新时间',
        geometryLength: '几何长度',
        lengthError: '长度偏差',
        distanceDelta: 'DO 距离增量',
        raceTimeDelta: 'DO 计时增量',
        lateralError: '横向偏离',
        gateWidth: 'Gate 宽度',
        gateDepth: 'Gate 深度',
        warnings: '警告',
        noWarnings: '无警告',
        trackSavedDetails: (id: number, type: string, length: number, session: string, laps: number) => `赛道 #${id} 已保存：${type}，${length.toFixed(0)} m，来源 ${session || '--'}，观察圈数 ${laps}。`,
        warningLabels: {
            distance_traveled_mismatch: '距离字段偏差',
            route_deviation: '路线偏离',
            route_progress_low: '进度不足',
            geometry_length_mismatch: '几何长度偏差',
        },
        useCases: {
            Road: '公路',
            Rally: '拉力',
            Drift: '漂移',
            Offroad: '越野',
            Drag: '直线',
            Wet: '雨天',
            Test: '测试',
        },
        fieldGroups: {
            vehicle: '车辆信息',
            power: '动力与重量',
            tire: '轮胎',
            gearing: '齿比',
            alignment: '轮胎定位',
            antiroll: '防倾杆',
            springs: '弹簧与车高',
            damping: '阻尼',
            aero: '空气动力学设置',
            brake: '刹车',
            differential: '差速器',
            notes: '备注',
        },
        eventTimeline: '事件时间线',
        eventSubtitle: '当前会话的本地规则识别结果',
        noEvents: '当前会话暂无事件。',
        eventEvidence: '证据',
        eventSuggestions: '初始建议',
        eventAdviceTitle: '调校建议',
        advicePlaceholder: '暂无明确调校建议。后续将结合调校档案、驾驶风格和同条件对比生成聚焦建议。',
        tuningNote: '调校说明',
        eventDuration: '持续',
        eventSegment: '阶段',
        eventStarted: '开始',
        severityLabel: '严重度',
        severityLow: '低',
        severityMedium: '中',
        severityHigh: '高',
        durationMsUnit: '毫秒',
        durationSecondUnit: '秒',
        events: {
            launch_wheelspin: '起步打滑',
            launch_bog_down: '起步憋转',
            short_gear: '挡位过短',
            long_gear_bog_down: '长齿比憋转',
            top_speed_limited_by_gearing: '极速受齿比限制',
            front_brake_lockup: '前轮抱死',
            rear_brake_lockup: '后轮抱死',
            corner_entry_understeer: '入弯推头',
            mid_corner_understeer: '持续过弯推头',
            corner_exit_oversteer: '出弯甩尾',
            power_understeer: '动力推头',
            snap_oversteer: '突然甩尾',
            high_speed_four_wheel_slide: '高速四轮侧滑',
            tire_overheat: '轮胎过热',
            tire_temp_imbalance: '胎温不均',
            suspension_bottom_out: '悬挂触底',
        },
        segments: {
            launch: '起步',
            acceleration: '加速',
            braking: '制动',
            corner_entry: '入弯',
            mid_corner: '持续过弯',
            corner_exit: '出弯',
            cornering: '过弯',
            high_speed_corner: '高速弯',
            tire: '轮胎',
            suspension: '悬挂',
        },
        evidenceLabels: {
            speed_kmh: '速度',
            speed_min_kmh: '最低速度',
            speed_avg_kmh: '平均速度',
            speed_max_kmh: '最高速度',
            speed_band: '速度区间',
            gear: '挡位',
            throttle: '油门',
            brake: '刹车',
            steer_abs: '方向输入',
            front_slip_ratio: '前轮滑移率',
            rear_slip_ratio: '后轮滑移率',
            max_slip_ratio: '最大滑移率',
            rpm_ratio: '转速比例',
            front_combined_slip: '前轮综合滑移',
            rear_combined_slip: '后轮综合滑移',
            slip_delta: '前后滑移差',
            corner_operation_state: '过弯操作状态',
            yaw_rate_abs: '横摆角速度',
            max_suspension_travel: '最大悬挂行程',
            front_suspension: '前悬挂行程',
            rear_suspension: '后悬挂行程',
            pitch_rate_abs: '俯仰角速度',
            roll_rate_abs: '侧倾角速度',
            front_tire_temp: '前胎温',
            rear_tire_temp: '后胎温',
            tire_temp_delta: '胎温差',
        },
        actionCategories: {
            aero: '空气动力',
            alignment: '定位',
            brake: '刹车',
            damping: '阻尼',
            differential: '差速器',
            gearing: '齿比',
            rollback: '回退',
            suspension: '悬挂',
            tire: '轮胎',
        },
        actionItems: {
            brake_balance: '刹车平衡',
            brake_pressure: '刹车压力',
            bump: '压缩阻尼',
            current_gear: '当前挡位齿比',
            drive_diff_accel: '驱动轮加速差速',
            drive_tire_pressure: '驱动轮胎压',
            final_drive: '终传比',
            front_diff_accel: '前差速加速',
            front_diff_decel: '前差速减速',
            front_tire_pressure: '前胎压',
            front_and_rear_aero: '前后下压力',
            front_arb: '前防倾杆',
            front_camber: '前轮外倾角',
            front_rebound: '前回弹阻尼',
            gear_1: '1 挡齿比',
            gear_2: '2 挡齿比',
            gear_3: '3 挡齿比',
            gear_4: '4 挡齿比',
            gear_5: '5 挡齿比',
            gear_6: '6 挡齿比',
            gear_7: '7 挡齿比',
            gear_8: '8 挡齿比',
            gear_9: '9 挡齿比',
            gear_10: '10 挡齿比',
            rear_arb: '后防倾杆',
            rear_diff_accel: '后差速加速',
            rear_diff_decel: '后差速减速',
            rear_rebound: '后回弹阻尼',
            rear_tire_pressure: '后胎压',
            ride_height: '车身高度',
            spring_rate: '弹簧硬度',
            tire_pressure: '胎压',
        },
        actionDirections: {
            check: '检查',
            decrease: '降低',
            increase: '提高',
        },
        actionAmounts: {
            direction_only: '仅方向',
            'one small step': '一小格',
            'slightly more negative': '略微增加负外倾',
            'avoid bottoming': '避免触底',
            '1%-2% rearward': '向后 1%-2%',
            '1%-2% forward': '向前 1%-2%',
        },
        actionReasons: {
            'avoid hitting the top of the gear too early': '避免过早顶到当前挡位转速上限',
            'help the engine stay in the power band': '帮助发动机保持在动力区间',
            'increase front grip on entry': '提高入弯阶段前轮抓地',
            'increase front grip on steady cornering': '提高持续过弯时的前轮抓地',
            'increase front grip under power': '提高给油时前轮抓地',
            'increase high-speed grip': '提高高速抓地',
            'increase launch traction': '提高起步牵引力',
            'increase rear grip': '提高后轮抓地',
            'increase tire contact patch': '增加轮胎接地面积',
            'improve front tire contact in cornering': '改善持续过弯时的前轮接地状态',
            'improve rear compliance under braking': '改善制动时后轴贴服性',
            'lengthen all gears if multiple gears are short': '如果多个挡位都偏短，整体拉长齿比',
            'let the front tires load more smoothly': '让前轮载荷转移更平顺',
            'make threshold braking easier': '降低临界刹车控制难度',
            'prevent aero and suspension instability': '避免下压力和悬挂状态不稳定',
            'reduce bottoming frequency': '减少触底频率',
            'reduce driven-wheel slip': '降低驱动轮滑移',
            'reduce front lockup tendency': '降低前轮抱死倾向',
            'reduce power oversteer': '降低动力甩尾倾向',
            'reduce rear lockup tendency': '降低后轮抱死倾向',
            'reduce wheel torque during launch': '降低起步时轮上扭矩',
            'reduce wheel torque on exit': '降低出弯时轮上扭矩',
            'reduce power-on understeer': '降低给油推头倾向',
            'reduce sustained tire scrub': '降低持续轮胎擦滑',
            'restore suspension travel': '恢复悬挂可用行程',
            'rotate the car more in steady cornering': '增强持续过弯时的车身转向响应',
            'make rear response less abrupt': '让车尾响应更平顺',
            'stabilize the rear axle while off throttle': '稳定收油时后轴',
            'reduce tire overheating tendency': '降低轮胎过热倾向',
            'balance front and rear tire temperatures': '平衡前后胎温',
            'rebalance axle load transfer': '重新平衡前后轴载荷转移',
            'rollback half of the last related change': '回退上次相关修改的一半',
            'stabilize tire temperature and contact patch': '稳定胎温和接地状态',
            'shorten launch gearing': '缩短起步齿比',
            'shorten road acceleration gearing': '缩短公路加速齿比',
            rollback_retest_worsened: '恢复上次应用调校方案前的数值',
            half_reverse_retest_worsened: '按上次应用调校方案的一半幅度反向微调',
            'stabilize the rear axle while braking': '稳定制动时的后轴',
            'support compression on impacts': '提高冲击压缩阶段支撑',
            'increase top speed headroom': '增加极速转速余量',
            'verify aero drag is not limiting top speed': '确认空阻没有限制极速',
            'verify traction is not limiting exit drive': '确认牵引力没有限制出弯加速',
        },
    },
} as const;

type Lang = keyof typeof COPY;
type Copy = (typeof COPY)[Lang];

const emptyWheel: WheelTelemetry = {
    slipRatio: 0,
    slipAngle: 0,
    combinedSlip: 0,
    tireTemp: 0,
    suspensionTravel: 0,
    suspensionTravelMeters: 0,
    wheelRotationSpeed: 0,
    rumbleStrip: 0,
    puddleDepth: 0,
    surfaceRumble: 0,
};

const emptyStatus: TelemetryStatus = {
    running: false,
    mode: 'idle',
    analysisMode: 'none',
    address: '0.0.0.0',
    port: 5301,
    packetLength: 324,
    rawPackets: 0,
    validPackets: 0,
    invalidPackets: 0,
    parseErrors: 0,
    lastDatagramAt: '',
    lastDatagramBytes: 0,
    lastDatagramRemote: '',
    lastPacketAt: '',
    lastError: '',
    hasCurrentFrame: false,
    recordingActive: false,
    recordingBytes: 0,
    recordingLimitBytes: 128 * 1024 * 1024,
    recordingPackets: 0,
    recordingTruncated: false,
};

const emptyReplayStatus: TelemetryReplayStatus = {
    running: false,
    paused: false,
    sessionId: 0,
    speed: 1,
    positionMs: 0,
    durationMs: 0,
    progress01: 0,
    packetIndex: 0,
    packetCount: 0,
    lastError: '',
};

const unknownTestConditions: TestConditions = {
    driverMode: 'unknown',
    brakeAssist: 'unknown',
    steeringAssist: 'unknown',
    tractionControl: 'unknown',
    stabilityControl: 'unknown',
    shifting: 'unknown',
    launchControl: 'unknown',
};

const testConditionOptions: Record<keyof TestConditions, string[]> = {
    driverMode: ['unknown', 'player', 'auto'],
    brakeAssist: ['unknown', 'assisted', 'abs_on', 'abs_off'],
    steeringAssist: ['unknown', 'auto', 'assisted', 'standard', 'simulation'],
    tractionControl: ['unknown', 'on', 'off'],
    stabilityControl: ['unknown', 'on', 'off'],
    shifting: ['unknown', 'automatic', 'manual'],
    launchControl: ['unknown', 'on', 'off'],
};

type ProfileField = {
    key: keyof TuneProfileInput;
    group: keyof Copy['fieldGroups'];
    label: { en: string; zh: string };
    unit?: { en: string; zh: string };
    kind: 'text' | 'number' | 'textarea' | 'select';
    step?: string;
    readOnly?: boolean;
};

const tuneUseCaseValues = ['Road', 'Rally', 'Drift', 'Offroad', 'Drag', 'Wet', 'Test'] as const;

const tuneUseCaseAliases: Record<string, string> = {
    road: 'Road',
    公路: 'Road',
    rally: 'Rally',
    拉力: 'Rally',
    drift: 'Drift',
    漂移: 'Drift',
    offroad: 'Offroad',
    'off-road': 'Offroad',
    越野: 'Offroad',
    drag: 'Drag',
    直线: 'Drag',
    wet: 'Wet',
    雨天: 'Wet',
    test: 'Test',
    测试: 'Test',
};

const profileFields: ProfileField[] = [
    {key: 'carName', group: 'vehicle', label: {en: 'Car name', zh: '车辆名称'}, kind: 'text'},
    {key: 'versionName', group: 'vehicle', label: {en: 'Version', zh: '版本'}, kind: 'text'},
    {key: 'carOrdinal', group: 'vehicle', label: {en: 'Car ordinal', zh: '车辆 ID'}, kind: 'number'},
    {key: 'carCategory', group: 'vehicle', label: {en: 'Car category', zh: '车辆分类'}, kind: 'number'},
    {key: 'carClass', group: 'vehicle', label: {en: 'Class', zh: '等级'}, kind: 'text'},
    {key: 'pi', group: 'vehicle', label: {en: 'PI', zh: 'PI'}, kind: 'number'},
    {key: 'drivetrain', group: 'vehicle', label: {en: 'Drivetrain', zh: '驱动'}, kind: 'text'},
    {key: 'numCylinders', group: 'vehicle', label: {en: 'Cylinders', zh: '气缸数'}, kind: 'number'},
    {key: 'useCase', group: 'vehicle', label: {en: 'Use case', zh: '用途'}, kind: 'select'},
    {key: 'powerKW', group: 'power', label: {en: 'Power', zh: '功率'}, unit: {en: 'kW', zh: '千瓦'}, kind: 'number', step: '1'},
    {key: 'torqueNM', group: 'power', label: {en: 'Torque', zh: '扭矩'}, unit: {en: 'Nm', zh: '牛米'}, kind: 'number', step: '1'},
    {key: 'weightKG', group: 'power', label: {en: 'Weight', zh: '车重'}, unit: {en: 'kg', zh: '千克'}, kind: 'number', step: '1'},
    {key: 'frontWeightPct', group: 'power', label: {en: 'Front weight', zh: '前轮重量分配'}, unit: {en: '%', zh: '%'}, kind: 'number', step: '0.1'},
    {key: 'powerToWeightKWPerKG', group: 'power', label: {en: 'Power to weight', zh: '功率比重'}, unit: {en: 'kW/kg', zh: '千瓦/千克'}, kind: 'number', step: '0.0001', readOnly: true},
    {key: 'peakTorqueRPM', group: 'power', label: {en: 'Peak torque RPM', zh: '峰值扭矩转速'}, unit: {en: 'rpm', zh: '转/分'}, kind: 'number', step: '1'},
    {key: 'peakPowerRPM', group: 'power', label: {en: 'Peak power RPM', zh: '峰值功率转速'}, unit: {en: 'rpm', zh: '转/分'}, kind: 'number', step: '1'},
    {key: 'redlineRPM', group: 'power', label: {en: 'Redline RPM', zh: '红线转速'}, unit: {en: 'rpm', zh: '转/分'}, kind: 'number', step: '1'},
    {key: 'frontTirePressure', group: 'tire', label: {en: 'Front tire pressure', zh: '前侧胎压'}, unit: {en: 'BAR', zh: '巴'}, kind: 'number', step: '0.01'},
    {key: 'rearTirePressure', group: 'tire', label: {en: 'Rear tire pressure', zh: '后侧胎压'}, unit: {en: 'BAR', zh: '巴'}, kind: 'number', step: '0.01'},
    {key: 'finalDrive', group: 'gearing', label: {en: 'Final drive', zh: '终传比'}, kind: 'number', step: '0.01'},
    {key: 'gear1', group: 'gearing', label: {en: '1st gear', zh: '1 挡'}, kind: 'number', step: '0.01'},
    {key: 'gear2', group: 'gearing', label: {en: '2nd gear', zh: '2 挡'}, kind: 'number', step: '0.01'},
    {key: 'gear3', group: 'gearing', label: {en: '3rd gear', zh: '3 挡'}, kind: 'number', step: '0.01'},
    {key: 'gear4', group: 'gearing', label: {en: '4th gear', zh: '4 挡'}, kind: 'number', step: '0.01'},
    {key: 'gear5', group: 'gearing', label: {en: '5th gear', zh: '5 挡'}, kind: 'number', step: '0.01'},
    {key: 'gear6', group: 'gearing', label: {en: '6th gear', zh: '6 挡'}, kind: 'number', step: '0.01'},
    {key: 'gear7', group: 'gearing', label: {en: '7th gear', zh: '7 挡'}, kind: 'number', step: '0.01'},
    {key: 'gear8', group: 'gearing', label: {en: '8th gear', zh: '8 挡'}, kind: 'number', step: '0.01'},
    {key: 'gear9', group: 'gearing', label: {en: '9th gear', zh: '9 挡'}, kind: 'number', step: '0.01'},
    {key: 'gear10', group: 'gearing', label: {en: '10th gear', zh: '10 挡'}, kind: 'number', step: '0.01'},
    {key: 'frontCamber', group: 'alignment', label: {en: 'Front camber', zh: '前侧外倾角'}, unit: {en: '°', zh: '度'}, kind: 'number', step: '0.1'},
    {key: 'rearCamber', group: 'alignment', label: {en: 'Rear camber', zh: '后侧外倾角'}, unit: {en: '°', zh: '度'}, kind: 'number', step: '0.1'},
    {key: 'frontToe', group: 'alignment', label: {en: 'Front toe', zh: '前侧束角'}, unit: {en: '°', zh: '度'}, kind: 'number', step: '0.1'},
    {key: 'rearToe', group: 'alignment', label: {en: 'Rear toe', zh: '后侧束角'}, unit: {en: '°', zh: '度'}, kind: 'number', step: '0.1'},
    {key: 'caster', group: 'alignment', label: {en: 'Front caster', zh: '前轮后倾角'}, unit: {en: '°', zh: '度'}, kind: 'number', step: '0.1'},
    {key: 'frontArb', group: 'antiroll', label: {en: 'Front ARB', zh: '前侧防倾杆'}, kind: 'number', step: '0.1'},
    {key: 'rearArb', group: 'antiroll', label: {en: 'Rear ARB', zh: '后侧防倾杆'}, kind: 'number', step: '0.1'},
    {key: 'frontSpring', group: 'springs', label: {en: 'Front spring', zh: '前侧弹簧'}, unit: {en: 'kgf/mm', zh: 'kgf/mm'}, kind: 'number', step: '0.1'},
    {key: 'rearSpring', group: 'springs', label: {en: 'Rear spring', zh: '后侧弹簧'}, unit: {en: 'kgf/mm', zh: 'kgf/mm'}, kind: 'number', step: '0.1'},
    {key: 'frontRideHeight', group: 'springs', label: {en: 'Front ride height', zh: '前侧车身高度'}, unit: {en: 'cm', zh: '厘米'}, kind: 'number', step: '0.1'},
    {key: 'rearRideHeight', group: 'springs', label: {en: 'Rear ride height', zh: '后侧车身高度'}, unit: {en: 'cm', zh: '厘米'}, kind: 'number', step: '0.1'},
    {key: 'frontRebound', group: 'damping', label: {en: 'Front rebound', zh: '前侧回弹硬度'}, kind: 'number', step: '0.1'},
    {key: 'rearRebound', group: 'damping', label: {en: 'Rear rebound', zh: '后侧回弹硬度'}, kind: 'number', step: '0.1'},
    {key: 'frontBump', group: 'damping', label: {en: 'Front bump', zh: '前侧压缩硬度'}, kind: 'number', step: '0.1'},
    {key: 'rearBump', group: 'damping', label: {en: 'Rear bump', zh: '后侧压缩硬度'}, kind: 'number', step: '0.1'},
    {key: 'frontAero', group: 'aero', label: {en: 'Front aero', zh: '前侧下压力'}, unit: {en: 'kgf', zh: '千克力'}, kind: 'number', step: '1'},
    {key: 'rearAero', group: 'aero', label: {en: 'Rear aero', zh: '后侧下压力'}, unit: {en: 'kgf', zh: '千克力'}, kind: 'number', step: '1'},
    {key: 'aeroBalance', group: 'aero', label: {en: 'Aero balance', zh: '空力平衡'}, kind: 'number'},
    {key: 'brakeBalance', group: 'brake', label: {en: 'Brake balance', zh: '制动力平衡'}, unit: {en: '%', zh: '%'}, kind: 'number', step: '1'},
    {key: 'brakePressure', group: 'brake', label: {en: 'Brake pressure', zh: '制动力压力'}, unit: {en: '%', zh: '%'}, kind: 'number', step: '1'},
    {key: 'frontDiffAccel', group: 'differential', label: {en: 'Front diff accel', zh: '前侧加速'}, unit: {en: '%', zh: '%'}, kind: 'number', step: '1'},
    {key: 'frontDiffDecel', group: 'differential', label: {en: 'Front diff decel', zh: '前侧减速'}, unit: {en: '%', zh: '%'}, kind: 'number', step: '1'},
    {key: 'rearDiffAccel', group: 'differential', label: {en: 'Rear diff accel', zh: '后侧加速'}, unit: {en: '%', zh: '%'}, kind: 'number', step: '1'},
    {key: 'rearDiffDecel', group: 'differential', label: {en: 'Rear diff decel', zh: '后侧减速'}, unit: {en: '%', zh: '%'}, kind: 'number', step: '1'},
    {key: 'centerDiffBalance', group: 'differential', label: {en: 'Center diff balance', zh: '中央平衡'}, unit: {en: '%', zh: '%'}, kind: 'number', step: '1'},
    {key: 'notes', group: 'notes', label: {en: 'Notes', zh: '备注'}, kind: 'textarea'},
];

const coreTuneGroupOrder = ['tire', 'alignment', 'antiroll', 'springs', 'damping', 'aero', 'brake', 'differential'] as Array<keyof Copy['fieldGroups']>;

const coreTuneFieldOrder = [
    'frontTirePressure',
    'rearTirePressure',
    'frontCamber',
    'rearCamber',
    'frontToe',
    'rearToe',
    'caster',
    'frontArb',
    'rearArb',
    'frontSpring',
    'rearSpring',
    'frontRideHeight',
    'rearRideHeight',
    'frontRebound',
    'rearRebound',
    'frontBump',
    'rearBump',
    'frontAero',
    'rearAero',
    'brakeBalance',
    'brakePressure',
    'frontDiffAccel',
    'frontDiffDecel',
    'rearDiffAccel',
    'rearDiffDecel',
    'centerDiffBalance',
] as Array<keyof TuneProfileInput>;

const coreTuneFieldOrderSet = new Set<string>(coreTuneFieldOrder.map(String));

const emptyProfileInput: TuneProfileInput = {
    carName: '',
    carClass: '',
    drivetrain: '',
    useCase: '',
    versionName: '',
    notes: '',
};

const emptyRoadBaselineForm: RoadStaticTuneBaselineForm = {
    carName: '',
    versionName: 'Road Baseline',
    carOrdinal: '',
    carCategory: '',
    pi: '',
    drivetrain: 'RWD',
    tireCompound: 'sport',
    weightKG: '',
    frontWeightPct: '',
    powerKW: '',
    torqueNM: '',
    redlineRPM: '',
    gearCount: '',
    tireDiameterCm: '',
    targetTopSpeedKmh: '',
    tireWidthMm: '',
    tireAspectRatio: '',
    tireRimInches: '',
    frontRideHeightMinCm: '',
    frontRideHeightMaxCm: '',
    rearRideHeightMinCm: '',
    rearRideHeightMaxCm: '',
    frontAeroMinKgf: '',
    frontAeroMaxKgf: '',
    rearAeroMinKgf: '',
    rearAeroMaxKgf: '',
    frontRideHeightAdjustable: 'true',
    rearRideHeightAdjustable: 'true',
    frontAeroAdjustable: 'true',
    rearAeroAdjustable: 'true',
    useCase: 'Road',
    gearingEnabled: 'false',
    balanceBias: '100',
    stiffnessBias: '100',
    speedBias: '100',
};

const quickTuneStorageKey = 'fh6-quick-tune-state-v1';
const quickTuneUseCases = ['Road', 'Drift', 'Rally', 'Offroad', 'Drag'] as const;
const quickTuneTireCompounds = ['stock', 'street', 'sport', 'semi', 'slick', 'rally', 'offroad', 'drift', 'drag', 'snow'] as const;
const recommendedUseCaseLabels: Record<string, string> = {
    Road: '公路',
    Drift: '漂移',
    Rally: '拉力',
    Offroad: '越野',
    Drag: '直线',
};
const recommendedTireCompoundLabels: Record<string, string> = {
    stock: '原厂',
    street: '街胎',
    sport: '运动',
    semi: '半热熔',
    slick: '热熔胎',
    rally: '拉力',
    offroad: '越野',
    drift: '漂移',
    drag: '直线',
    snow: '雪地',
};
const recommendedClassDefaultPI: Record<string, string> = {
    D: '400',
    C: '500',
    B: '600',
    A: '700',
    S1: '800',
    S2: '900',
    R: '998',
    X: '999',
};
const recommendedPIDefaultClass = Object.entries(recommendedClassDefaultPI).reduce<Record<string, string>>((acc, [carClass, pi]) => {
    acc[pi] = carClass;
    return acc;
}, {});
const emptyRecommendedCarForm: RecommendedCarForm = {
    id: '',
    name: '',
    useCase: 'Road',
    useCaseLabel: '公路',
    pi: '700',
    carClass: 'A',
    drivetrain: 'AWD',
    tireCompound: 'sport',
    tireCompoundLabel: '运动',
    weightKG: '',
    frontWeightPct: '',
    tuneCode: '',
    imageSrc: '',
    tags: '',
    reason: '',
};

function defaultRecommendedCarsVersion() {
    const date = new Date();
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const day = String(date.getDate()).padStart(2, '0');
    return `${year}-${month}-${day}-001`;
}

type QuickTuneStoredState = {
    form?: RoadStaticTuneBaselineForm;
    result?: RoadStaticTuneBaselineResult | null;
    selectedFields?: string[];
    targetProfileId?: number;
};

function recommendedCarInputFromForm(form: RecommendedCarForm): RecommendedCarInput {
    const pi = parseRequiredNumber(form.pi, 'PI');
    const weightKG = parseOptionalRecommendedNumber(form.weightKG, 'weightKG');
    const frontWeightPct = parseOptionalRecommendedNumber(form.frontWeightPct, 'frontWeightPct');
    const input: RecommendedCarInput = {
        id: form.id.trim(),
        name: form.name.trim(),
        useCase: form.useCase.trim(),
        useCaseLabel: recommendedUseCaseLabels[form.useCase] || form.useCaseLabel.trim(),
        pi,
        carClass: form.carClass.trim().toUpperCase(),
        drivetrain: form.drivetrain.trim().toUpperCase(),
        tireCompound: form.tireCompound.trim(),
        tireCompoundLabel: recommendedTireCompoundLabels[form.tireCompound] || form.tireCompoundLabel.trim(),
        weightKG,
        frontWeightPct,
        tuneCode: normalizeTuneCodeInput(form.tuneCode),
        imageSrc: form.imageSrc.trim() || undefined,
        tags: form.tags.split(',').map(tag => tag.trim()).filter(Boolean),
        reason: form.reason.trim(),
    };
    const required: Array<[keyof RecommendedCarInput, string]> = [
        ['name', 'name'],
        ['useCase', 'useCase'],
        ['carClass', 'carClass'],
        ['drivetrain', 'drivetrain'],
        ['tireCompound', 'tireCompound'],
        ['tuneCode', 'tuneCode'],
    ];
    const missing = required.find(([key]) => !String(input[key] || '').trim());
    if (missing) {
        throw new Error(`${missing[1]} is required`);
    }
    if (pi < 100 || pi > 999) {
        throw new Error('PI must be between 100 and 999');
    }
    if (weightKG < 0) {
        throw new Error('weightKG must be empty or greater than 0');
    }
    if (frontWeightPct < 0 || frontWeightPct >= 100) {
        throw new Error('frontWeightPct must be empty or between 1 and 99');
    }
    return input;
}

function recommendedCarFormFromCar(car: RecommendedCarInput): RecommendedCarForm {
    return {
        id: car.id || '',
        name: car.name || '',
        useCase: car.useCase || 'Road',
        useCaseLabel: car.useCaseLabel || '',
        pi: String(car.pi || ''),
        carClass: car.carClass || '',
        drivetrain: car.drivetrain || '',
        tireCompound: car.tireCompound || '',
        tireCompoundLabel: car.tireCompoundLabel || '',
        weightKG: String(car.weightKG || ''),
        frontWeightPct: String(car.frontWeightPct || ''),
        tuneCode: car.tuneCode || '',
        imageSrc: car.imageSrc || '',
        tags: (car.tags || []).join(', '),
        reason: car.reason || '',
    };
}

function recommendedCarFormFromHarvestCandidate(candidate: TuneHarvestCandidate): RecommendedCarForm {
    const carClass = candidate.carClass || classFromPI(candidate.pi) || 'A';
    const pi = candidate.pi || Number(recommendedClassDefaultPI[carClass] || 700);
    const useCase = (quickTuneUseCases as readonly string[]).includes(candidate.useCase) ? candidate.useCase : 'Road';
    const tireCompound = (quickTuneTireCompounds as readonly string[]).includes(candidate.tireCompound) ? candidate.tireCompound : 'sport';
    const tags = [
        candidate.source,
        candidate.tuner,
        candidate.difficulty,
        candidate.bestFor,
    ].filter(Boolean);
    return {
        ...emptyRecommendedCarForm,
        id: '',
        name: candidate.carName || [candidate.year || '', candidate.make, candidate.model].filter(Boolean).join(' '),
        useCase,
        useCaseLabel: recommendedUseCaseLabels[useCase] || '',
        pi: String(pi),
        carClass,
        drivetrain: candidate.drivetrain || 'AWD',
        tireCompound,
        tireCompoundLabel: recommendedTireCompoundLabels[tireCompound] || '',
        tuneCode: formatTuneHarvestShareCode(candidate.shareCode),
        tags: tags.join(', '),
        reason: [
            candidate.tuneName,
            candidate.bestFor,
            candidate.notes,
            candidate.sourceUrl,
        ].filter(Boolean).join(' / '),
    };
}

function classFromPI(pi: number) {
    if (!pi) {
        return '';
    }
    if (pi <= 400) return 'D';
    if (pi <= 500) return 'C';
    if (pi <= 600) return 'B';
    if (pi <= 700) return 'A';
    if (pi <= 800) return 'S1';
    if (pi <= 900) return 'S2';
    if (pi <= 998) return 'R';
    return 'X';
}

function normalizeTuneCodeInput(value: string) {
    const digits = String(value || '').replace(/\D/g, '');
    return digits.length === 9 ? digits : String(value || '').trim();
}

function formatTuneHarvestShareCode(value: string) {
    const digits = String(value || '').replace(/\D/g, '');
    if (digits.length !== 9) {
        return value || '';
    }
    return `${digits.slice(0, 3)} ${digits.slice(3, 6)} ${digits.slice(6)}`;
}

function tuneHarvestStatusLabel(terms: Copy, status: string) {
    if (status === 'all' || status === 'pending' || status === 'rejected' || status === 'imported') {
        return terms.tuneHarvestStatusLabels[status];
    }
    return status || '--';
}

function tuneHarvestSearchDigits(value: string) {
    return String(value || '').replace(/\D/g, '');
}

function candidateMatchesTuneHarvestSearch(candidate: TuneHarvestCandidate, search: string) {
    const terms = String(search || '').trim().toLowerCase().split(/\s+/).filter(Boolean);
    if (terms.length === 0) {
        return true;
    }
    const searchable = [
        candidate.source,
        candidate.sourceRef,
        candidate.sourceUrl,
        candidate.sourceCarId,
        candidate.rawKey,
        candidate.shareCode,
        formatTuneHarvestShareCode(candidate.shareCode),
        candidate.year ? String(candidate.year) : '',
        candidate.make,
        candidate.model,
        candidate.carName,
        candidate.matchedCarId,
        candidate.matchReason,
        candidate.useCase,
        candidate.carClass,
        candidate.pi ? String(candidate.pi) : '',
        candidate.drivetrain,
        candidate.tireCompound,
        candidate.tuner,
        candidate.tuneName,
        candidate.bestFor,
        candidate.difficulty,
        candidate.notes,
        candidate.status,
        candidate.rejectionReason,
    ].filter(Boolean).join(' ').toLowerCase();
    const shareCodeDigits = tuneHarvestSearchDigits(candidate.shareCode);
    return terms.every(term => {
        if (searchable.includes(term)) {
            return true;
        }
        const digits = tuneHarvestSearchDigits(term);
        return Boolean(digits && shareCodeDigits.includes(digits));
    });
}

function recommendedCarExportInput(car: RecommendedCar): RecommendedCarInput {
    return {
        id: car.id,
        name: car.name,
        useCase: car.useCase,
        useCaseLabel: car.useCaseLabel,
        pi: car.pi,
        carClass: car.carClass,
        drivetrain: car.drivetrain,
        tireCompound: car.tireCompound,
        tireCompoundLabel: car.tireCompoundLabel,
        weightKG: car.weightKG || 0,
        frontWeightPct: car.frontWeightPct || 0,
        tuneCode: car.tuneCode,
        imageSrc: car.imageSrc,
        tags: car.tags || [],
        reason: car.reason || '',
    };
}

function generateRecommendedCarPreviewID(input: RecommendedCarInput) {
    const words = input.name
        .trim()
        .split(/\s+/)
        .filter(Boolean);
    if (words.length > 1 && /^(19|20)\d{2}$/.test(words[0])) {
        words.push(words.shift() || '');
    }
    const namePart = words
        .join(' ')
        .toLowerCase()
        .replace(/[^a-z0-9]+/g, '-')
        .replace(/^-+|-+$/g, '') || 'car';
    const tuneCode = normalizeTuneCodeInput(input.tuneCode);
    return [namePart, input.useCase.toLowerCase(), `${input.carClass.toLowerCase()}${input.pi}`, tuneCode].filter(Boolean).join('-');
}

function formatUTC8DateTime(value: string, language: Lang) {
    if (!value) {
        return '--';
    }
    const date = new Date(value);
    if (Number.isNaN(date.getTime())) {
        return value;
    }
    return new Intl.DateTimeFormat(language === 'zh' ? 'zh-CN' : 'en-US', {
        timeZone: 'Asia/Shanghai',
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit',
        hour12: false,
    }).format(date);
}

function parseRequiredNumber(value: string, label: string) {
    const parsed = Number(value);
    if (!Number.isFinite(parsed)) {
        throw new Error(`${label} must be a number`);
    }
    return parsed;
}

function parseOptionalRecommendedNumber(value: string, label: string) {
    if (!value.trim()) {
        return 0;
    }
    return parseRequiredNumber(value, label);
}

const emptyRuleThresholdInput: RuleThresholdProfileInput = {
    name: '',
    carClass: '',
    drivetrain: '',
    useCase: '',
    gameMode: '',
    configJson: '',
};

const defaultPipelineInput: TuningPipelineRunInput = {
    sourceType: 'tire_lab_current',
    detectorId: 'tire_lab_problem_groups_v1',
    decisionerId: 'tire_problem_decision_v1',
    interpreterId: 'road_baseline_docs_v12_interpreter_v1',
};

function App() {
    const [language, setLanguage] = useState<Lang>(() => {
        const saved = window.localStorage.getItem('fh6-ui-language');
        return saved === 'zh' ? 'zh' : 'en';
    });
    const [interfaces, setInterfaces] = useState<NetworkInterface[]>([]);
    const [selectedAddress, setSelectedAddress] = useState('0.0.0.0');
    const [port, setPort] = useState(5301);
    const [status, setStatus] = useState<TelemetryStatus>(emptyStatus);
    const [replayStatus, setReplayStatus] = useState<TelemetryReplayStatus>(emptyReplayStatus);
    const [current, setCurrent] = useState<TelemetryFrame | null>(null);
    const [recent, setRecent] = useState<TelemetryFrame[]>([]);
    const [quickDiagnostic, setQuickDiagnostic] = useState<QuickDiagnostic | null>(null);
    const [tireModelDiagnostic, setTireModelDiagnostic] = useState<TireModelDiagnostic | null>(null);
    const [view, setView] = useState<ViewName>('tune_generator');
    const [developerTool, setDeveloperTool] = useState<DeveloperToolName>('do_fields');
    const [recommendedCarForm, setRecommendedCarForm] = useState<RecommendedCarForm>(emptyRecommendedCarForm);
    const [recommendedCars, setRecommendedCars] = useState<RecommendedCar[]>([]);
    const [recommendedCarsResult, setRecommendedCarsResult] = useState<RecommendedCarsFileResult | null>(null);
    const [recommendedCarsFileSelection, setRecommendedCarsFileSelection] = useState<RecommendedCarsFileSelection | null>(null);
    const [recommendedCarsVersion, setRecommendedCarsVersion] = useState(defaultRecommendedCarsVersion);
    const [recommendedCarsFormError, setRecommendedCarsFormError] = useState('');
    const [selectedRecommendedCarIds, setSelectedRecommendedCarIds] = useState<string[]>([]);
    const [recommendedCarFormOpen, setRecommendedCarFormOpen] = useState(false);
    const [editingRecommendedCarId, setEditingRecommendedCarId] = useState('');
    const [recommendedCarDetail, setRecommendedCarDetail] = useState<RecommendedCar | null>(null);
    const [tuneHarvestCandidates, setTuneHarvestCandidates] = useState<TuneHarvestCandidate[]>([]);
    const [tuneHarvestResult, setTuneHarvestResult] = useState<TuneHarvestRunResult | null>(null);
    const [tuneHarvestSources, setTuneHarvestSources] = useState<TuneHarvestSourceState>({
        jsr_chronic_sheet: true,
        codmunity: true,
        forzafire: false,
    });
    const [tuneHarvestDryRun, setTuneHarvestDryRun] = useState(false);
    const [tuneHarvestLimit, setTuneHarvestLimit] = useState('80');
    const [tuneHarvestStatusFilter, setTuneHarvestStatusFilter] = useState('pending');
    const [tuneHarvestSearch, setTuneHarvestSearch] = useState('');
    const [tuneHarvestRunning, setTuneHarvestRunning] = useState(false);
    const [tuneHarvestStopping, setTuneHarvestStopping] = useState(false);
    const [importingHarvestCandidateId, setImportingHarvestCandidateId] = useState<number | null>(null);
    const [profiles, setProfiles] = useState<TuneProfile[]>([]);
    const [profileSessionStats, setProfileSessionStats] = useState<TuneProfileSessionStat[]>([]);
    const [activeProfile, setActiveProfile] = useState<TuneProfile | null>(null);
    const [profileForm, setProfileForm] = useState<TuneProfileInput>(emptyProfileInput);
    const [editingProfileId, setEditingProfileId] = useState<number | null>(null);
    const [profileSnapshots, setProfileSnapshots] = useState<TuneProfileSnapshot[]>([]);
    const [roadBaselineForm, setRoadBaselineForm] = useState<RoadStaticTuneBaselineForm>(() => loadQuickTuneStoredState().form || emptyRoadBaselineForm);
    const [roadBaselineResult, setRoadBaselineResult] = useState<RoadStaticTuneBaselineResult | null>(() => loadQuickTuneStoredState().result || null);
    const [selectedBaselineFields, setSelectedBaselineFields] = useState<string[]>(() => loadQuickTuneStoredState().selectedFields || []);
    const [roadBaselineTargetProfileId, setRoadBaselineTargetProfileId] = useState(() => loadQuickTuneStoredState().targetProfileId || 0);
    const [roadBaselineAdvancedOpen, setRoadBaselineAdvancedOpen] = useState(false);
    const [quickTuneInputOpen, setQuickTuneInputOpen] = useState(false);
    const [quickTuneFieldErrors, setQuickTuneFieldErrors] = useState<QuickTuneFieldErrors>({});
    const roadBaselineBiasDebounceRef = useRef<number | null>(null);
    const [tuneWebStatus, setTuneWebStatus] = useState<TuneWebServerStatus>({running: false, port: 0, url: '', lanAddress: '', lastError: ''});
    const [tuneWebPort, setTuneWebPort] = useState('8787');
    const [sessions, setSessions] = useState<TelemetrySession[]>([]);
    const [selectedSessionId, setSelectedSessionId] = useState<number | null>(null);
    const [compareLeftId, setCompareLeftId] = useState<number | null>(null);
    const [compareRightId, setCompareRightId] = useState<number | null>(null);
    const [sessionComparison, setSessionComparison] = useState<SessionComparison | null>(null);
    const [sessionEvents, setSessionEvents] = useState<DetectedEvent[]>([]);
    const [sessionSamples, setSessionSamples] = useState<TelemetryFrame[]>([]);
    const [benchmarkTracks, setBenchmarkTracks] = useState<BenchmarkTrack[]>([]);
    const [sessionBenchmarkRuns, setSessionBenchmarkRuns] = useState<BenchmarkRun[]>([]);
    const [roadEvaluation, setRoadEvaluation] = useState<RoadSessionEvaluation | null>(null);
    const [sessionIssueSummary, setSessionIssueSummary] = useState<SessionIssueSummary | null>(null);
    const [roadTuningDecision, setRoadTuningDecision] = useState<RoadTuningDecision | null>(null);
    const [tunePlanDraft, setTunePlanDraft] = useState<TunePlanDraft | null>(null);
    const [selectedTunePlanActionIds, setSelectedTunePlanActionIds] = useState<string[]>([]);
    const [retestEvaluation, setRetestEvaluation] = useState<RetestEvaluation | null>(null);
    const [knowledgeStatus, setKnowledgeStatus] = useState<RoadTuningKnowledgeStatus | null>(null);
    const [selectedTrackId, setSelectedTrackId] = useState<number | null>(null);
    const [selectedTrackRuns, setSelectedTrackRuns] = useState<BenchmarkRun[]>([]);
    const [trackProfiles, setTrackProfiles] = useState<TrackProfile[]>([]);
    const [selectedTrackProfile, setSelectedTrackProfile] = useState<TrackProfile | null>(null);
    const [tireRegressionSamples, setTireRegressionSamples] = useState<TireRegressionSampleSummary[]>([]);
    const [selectedTireRegressionId, setSelectedTireRegressionId] = useState('');
    const [selectedTireRegressionSample, setSelectedTireRegressionSample] = useState<TireRegressionSample | null>(null);
    const [tireRegressionResults, setTireRegressionResults] = useState<TireRegressionResult[]>([]);
    const [tireRegressionSaveForm, setTireRegressionSaveForm] = useState<TireRegressionSaveFormState>({name: '', scenario: '', windowSeconds: '15'});
    const [tireRegressionExpectedForm, setTireRegressionExpectedForm] = useState<TireRegressionExpectedFormState>({
        allowedPhases: '',
        requiredGripTypes: '',
        allowedAxles: '',
        forbiddenGripTypes: '',
        minDataQuality: 'low_confidence',
        notes: '',
    });
    const [pipelineCatalog, setPipelineCatalog] = useState<TuningPipelineCatalog | null>(null);
    const [pipelineInput, setPipelineInput] = useState<TuningPipelineRunInput>(defaultPipelineInput);
    const [pipelineResult, setPipelineResult] = useState<TuningPipelineRunResult | null>(null);
    const [professionalPipelineConfig, setProfessionalPipelineConfig] = useState<ProfessionalPipelineConfig>({
        detectorId: defaultPipelineInput.detectorId,
        decisionerId: defaultPipelineInput.decisionerId,
        interpreterId: defaultPipelineInput.interpreterId,
    });
    const [professionalDiagnostic, setProfessionalDiagnostic] = useState<ProfessionalTuningDiagnostic | null>(null);
    const [trackCapture, setTrackCapture] = useState<TrackCaptureState>({recording: false, name: '', points: [], hasDrivingLine: false});
    const [replaySpeed, setReplaySpeed] = useState(1);
    const [ruleProfiles, setRuleProfiles] = useState<RuleThresholdProfile[]>([]);
    const [strategyTemplates, setStrategyTemplates] = useState<StrategyTemplate[]>([]);
    const [strategyAnalysis, setStrategyAnalysis] = useState<RoadStrategyAnalysis | null>(null);
    const [strategySessionIds, setStrategySessionIds] = useState<number[]>([]);
    const [selectedStrategyTemplateId, setSelectedStrategyTemplateId] = useState<number | null>(null);
    const [ruleForm, setRuleForm] = useState<RuleThresholdProfileInput>(emptyRuleThresholdInput);
    const [editingRuleId, setEditingRuleId] = useState<number | null>(null);
    const [upgradeUnlockRules, setUpgradeUnlockRules] = useState<UpgradeUnlockRule[]>([]);
    const [tuneExplanations, setTuneExplanations] = useState<TuneAdjustmentExplanation[]>([]);
    const [tuneInfluenceMap, setTuneInfluenceMap] = useState<TuneToTireInfluenceMap | null>(null);
    const [testConditions, setTestConditions] = useState<TestConditions>(unknownTestConditions);
    const [showAdvancedSettings, setShowAdvancedSettings] = useState(false);
    const [showNewProfileModal, setShowNewProfileModal] = useState(false);
    const [newProfileDraft, setNewProfileDraft] = useState<TuneProfileInput>({...emptyProfileInput, useCase: 'Road'});
    const [pendingStartChoice, setPendingStartChoice] = useState<PendingStartChoice | null>(null);
    const [pendingProfileMismatch, setPendingProfileMismatch] = useState<PendingProfileMismatch | null>(null);
    const [pendingSessionBind, setPendingSessionBind] = useState<PendingSessionBind | null>(null);
    const [pendingTrackMerge, setPendingTrackMerge] = useState<PendingTrackMerge | null>(null);
    const [reportMarkdown, setReportMarkdown] = useState('');
    const [message, setMessage] = useState('');
    const [busy, setBusy] = useState(false);
    const profileFormRef = useRef<HTMLDivElement>(null);
    const t = COPY[language];

    useEffect(() => {
        let mounted = true;

        async function refresh() {
            try {
                const [nextStatus, nextReplayStatus, nextCurrent, nextRecent, nextInterfaces, nextTuneWebStatus] = await Promise.all([
                    GetTelemetryStatus(),
                    GetTelemetryReplayStatus(),
                    GetCurrentTelemetry(),
                    GetRecentTelemetry(8),
                    GetNetworkInterfaces(),
                    GetTuneWebServerStatus(),
                ]);
                if (!mounted) {
                    return;
                }
                const normalizedStatus = nextStatus as TelemetryStatus;
                const nextQuickDiagnostic = normalizedStatus.analysisMode === 'quick'
                    ? await GetQuickDiagnostic().catch(() => null)
                    : null;
                const nextProfessionalDiagnostic = normalizedStatus.analysisMode === 'professional'
                    ? await GetProfessionalTuningDiagnostic().catch(() => null)
                    : null;
                const nextTireModelDiagnostic = normalizedStatus.analysisMode === 'tire_lab'
                    ? await GetTireModelDiagnostic().catch(() => null)
                    : null;
                if (!mounted) {
                    return;
                }
                const networkInterfaces = (nextInterfaces || []) as NetworkInterface[];
                const normalizedCurrent = (nextCurrent || null) as unknown as TelemetryFrame | null;
                setStatus(normalizedStatus);
                setReplayStatus((nextReplayStatus || emptyReplayStatus) as TelemetryReplayStatus);
                setCurrent(normalizedCurrent);
                setRecent((nextRecent || []) as unknown as TelemetryFrame[]);
                setInterfaces(networkInterfaces);
                setTuneWebStatus((nextTuneWebStatus || {running: false, port: 0, url: '', lanAddress: '', lastError: ''}) as TuneWebServerStatus);
                setQuickDiagnostic((nextQuickDiagnostic || null) as QuickDiagnostic | null);
                setProfessionalDiagnostic((nextProfessionalDiagnostic || null) as ProfessionalTuningDiagnostic | null);
                setTireModelDiagnostic((nextTireModelDiagnostic || null) as TireModelDiagnostic | null);
                setSelectedAddress(currentAddress => {
                    if (currentAddress !== '0.0.0.0') {
                        return currentAddress;
                    }
                    return preferredListenAddress(networkInterfaces);
                });
            } catch (error) {
                if (mounted) {
                    setMessage(error instanceof Error ? error.message : String(error));
                }
            }
        }

        refresh();
        const id = window.setInterval(refresh, 250);
        return () => {
            mounted = false;
            window.clearInterval(id);
        };
    }, []);

    useEffect(() => {
        window.localStorage.setItem('fh6-ui-language', language);
    }, [language]);

    useEffect(() => {
        setSelectedRecommendedCarIds(ids => ids.filter(id => recommendedCars.some(car => car.id === id)));
    }, [recommendedCars]);

    function applyRecommendedCarsFileSelection(cars: RecommendedCar[], selection: RecommendedCarsFileSelection | null) {
        setRecommendedCarsFileSelection(selection);
        if (selection?.exists && selection.version) {
            setRecommendedCarsVersion(selection.version);
        }
        const knownIDs = new Set(cars.map(car => car.id));
        const selectedTuneCodes = new Set((selection?.tuneCodes || []).map(normalizeTuneCodeInput).filter(Boolean));
        const selectedIDs = new Set((selection?.ids || []).filter(id => knownIDs.has(id)));
        cars.forEach(car => {
            if (selectedTuneCodes.has(normalizeTuneCodeInput(car.tuneCode))) {
                selectedIDs.add(car.id);
            }
        });
        setSelectedRecommendedCarIds(cars.filter(car => selectedIDs.has(car.id)).map(car => car.id));
    }

    useEffect(() => {
        if (view === 'developer' && developerTool === 'tune_harvest') {
            void refreshTuneHarvestCandidates(false);
        }
    }, [view, developerTool, tuneHarvestStatusFilter, tuneHarvestSearch]);

    useEffect(() => {
        saveQuickTuneStoredState({
            form: roadBaselineForm,
            result: roadBaselineResult,
            selectedFields: selectedBaselineFields,
            targetProfileId: roadBaselineTargetProfileId,
        });
    }, [roadBaselineForm, roadBaselineResult, selectedBaselineFields, roadBaselineTargetProfileId]);

    useEffect(() => () => {
        if (roadBaselineBiasDebounceRef.current !== null) {
            window.clearTimeout(roadBaselineBiasDebounceRef.current);
        }
    }, []);

    useEffect(() => {
        if (!selectedTireRegressionId) {
            setSelectedTireRegressionSample(null);
            setTireRegressionExpectedForm({
                allowedPhases: '',
                requiredGripTypes: '',
                allowedAxles: '',
                forbiddenGripTypes: '',
                minDataQuality: 'low_confidence',
                notes: '',
            });
            return;
        }
        let mounted = true;
        GetTireRegressionSample(selectedTireRegressionId)
            .then(sample => {
                if (!mounted) {
                    return;
                }
                const normalized = (sample || null) as TireRegressionSample | null;
                setSelectedTireRegressionSample(normalized);
                setTireRegressionExpectedForm(expectationToForm(normalized?.expected));
            })
            .catch(error => {
                if (mounted) {
                    setMessage(error instanceof Error ? error.message : String(error));
                }
            });
        return () => {
            mounted = false;
        };
    }, [selectedTireRegressionId]);

    useEffect(() => {
        setSelectedTunePlanActionIds((tunePlanDraft?.actions || []).filter(action => action.canApply).map(action => action.id));
    }, [tunePlanDraft]);

    useEffect(() => {
        loadMetadata();
    }, []);

    useEffect(() => {
        loadTestConditionDefaults();
    }, []);

    useEffect(() => {
        loadProfileSnapshots(editingProfileId);
    }, [editingProfileId]);

    useEffect(() => {
        if (!trackCapture.recording || !current || !hasVehicleTelemetry(current)) {
            return;
        }
        const point = telemetryPoint(current);
        setTrackCapture(capture => {
            if (!capture.recording) {
                return capture;
            }
            const last = capture.points[capture.points.length - 1];
            if (last && pointDistanceXZ(last, point) < 2) {
                return {
                    ...capture,
                    hasDrivingLine: capture.hasDrivingLine || Math.abs(current.drivingLine01) > 0.05,
                };
            }
            return {
                ...capture,
                points: [...capture.points, point],
                hasDrivingLine: capture.hasDrivingLine || Math.abs(current.drivingLine01) > 0.05,
            };
        });
    }, [current, trackCapture.recording]);

    useEffect(() => {
        let mounted = true;
        async function loadRuns() {
            if (!selectedTrackId) {
                setSelectedTrackRuns([]);
                setSelectedTrackProfile(null);
                return;
            }
            try {
                const [runs, profile] = await Promise.all([
                    ListBenchmarkRuns(selectedTrackId, 50),
                    GetTrackProfile(selectedTrackId),
                ]);
                if (mounted) {
                    setSelectedTrackRuns((runs || []) as BenchmarkRun[]);
                    setSelectedTrackProfile((profile || null) as TrackProfile | null);
                }
            } catch (error) {
                if (mounted) {
                    setMessage(error instanceof Error ? error.message : String(error));
                }
            }
        }
        loadRuns();
        return () => {
            mounted = false;
        };
    }, [selectedTrackId]);

    const speedPoints = useMemo(() => sparklinePoints(recent.map(item => item.speedKmh)), [recent]);
    const rpmPoints = useMemo(() => sparklinePoints(recent.map(item => item.rpmRatio * 100)), [recent]);
    const dataOutTargetAddress = selectedAddress === '0.0.0.0'
        ? preferredListenAddress(interfaces)
        : selectedAddress;
    const currentGameMode = current?.gameMode || 'unknown';
    const normalizedTestConditions = normalizeTestConditions(testConditions);
    const activeProfileMatchState = activeProfile ? profileFormTelemetryState(profileToInput(activeProfile), current) : 'unknown';
    async function startExpertTelemetryNow(address: string, udpPort: number, notice = '') {
        setCurrent(null);
        setRecent([]);
        setQuickDiagnostic(null);
        setProfessionalDiagnostic(null);
        await StartProfessionalTelemetry('0.0.0.0', udpPort);
        setMessage(notice ? `${notice} ${t.listening('0.0.0.0', udpPort)}` : t.listening('0.0.0.0', udpPort));
    }

    async function startQuickTelemetryNow(address: string, udpPort: number) {
        setCurrent(null);
        setRecent([]);
        setQuickDiagnostic(null);
        await StartQuickTelemetry('0.0.0.0', udpPort);
        setMessage(`${t.quickDiagnosis}: ${t.listening('0.0.0.0', udpPort)}`);
    }

    async function startTireModelTelemetryNow(address: string, udpPort: number) {
        setCurrent(null);
        setRecent([]);
        setQuickDiagnostic(null);
        setTireModelDiagnostic(null);
        await StartTireModelTelemetry(address, udpPort);
        setMessage(`${t.tireLab}: ${t.listening(address, udpPort)}`);
    }

    async function startTrackCaptureTelemetryNow(address: string, udpPort: number) {
        setCurrent(null);
        setRecent([]);
        setQuickDiagnostic(null);
        setTireModelDiagnostic(null);
        await StartTrackCaptureTelemetry(address, udpPort);
        setMessage(`${t.trackCaptureMode}: ${t.listening(address, udpPort)}`);
    }

    async function startTrackBaselineTelemetryNow(trackId: number, address: string, udpPort: number) {
        setCurrent(null);
        setRecent([]);
        setQuickDiagnostic(null);
        setTireModelDiagnostic(null);
        await StartTrackBaselineTelemetry(trackId, address, udpPort);
        setMessage(`${t.trackBaselineCapture}: ${t.listening(address, udpPort)}`);
    }

    async function validateTelemetryProfileBeforeStart(address: string, udpPort: number): Promise<StartValidationResult> {
        const telemetry = current;
        if (!hasTelemetryVehicleIdentity(telemetry)) {
            return {action: 'start', notice: t.profileMatchUnavailable};
        }

        if (activeProfile && profileMatchesTelemetry(activeProfile, telemetry)) {
            return {action: 'start'};
        }

        const candidates = (await ListTuneProfilesForVehicle(telemetry.carOrdinal, telemetry.carClass) || []) as TuneProfile[];
        if (activeProfile) {
            setPendingProfileMismatch({
                candidates,
                address,
                port: udpPort,
                telemetry,
                profile: activeProfile,
            });
            return {action: 'deferred'};
        }

        if (candidates.length === 1) {
            await SetActiveTuneProfile(candidates[0].id);
            setActiveProfile(candidates[0]);
            return {action: 'start'};
        }

        if (candidates.length > 1) {
            setPendingStartChoice({
                profiles: candidates,
                address,
                port: udpPort,
                carOrdinal: telemetry.carOrdinal,
                carClass: telemetry.carClass,
                carPi: telemetry.carPi,
            });
            return {action: 'deferred'};
        }

        await SetActiveTuneProfile(0);
        setActiveProfile(null);
        return {action: 'start', notice: t.noMatchingProfile};
    }

    async function start(mode: TelemetryStartMode) {
        setBusy(true);
        setMessage('');
        try {
            const udpPort = Number(port);
            if (mode === 'quick') {
                await startQuickTelemetryNow(selectedAddress, udpPort);
                return;
            }
            if (mode === 'tire_lab') {
                await startTireModelTelemetryNow(selectedAddress, udpPort);
                return;
            }
            if (mode === 'track_capture') {
                await startTrackCaptureTelemetryNow(selectedAddress, udpPort);
                return;
            }
            if (mode === 'professional') {
                await startExpertTelemetryNow(selectedAddress, udpPort);
                return;
            }
            const validation = await validateTelemetryProfileBeforeStart(selectedAddress, udpPort);
            if (validation.action === 'deferred') {
                return;
            }
            await startExpertTelemetryNow(selectedAddress, udpPort, validation.notice || '');
        } catch (error) {
            setMessage(error instanceof Error ? error.message : String(error));
        } finally {
            setBusy(false);
        }
    }

    async function startWithProfile(profile: TuneProfile) {
        if (!pendingStartChoice) {
            return;
        }
        setBusy(true);
        setMessage('');
        try {
            await SetActiveTuneProfile(profile.id);
            setActiveProfile(profile);
            const pending = pendingStartChoice;
            setPendingStartChoice(null);
            await startExpertTelemetryNow(pending.address, pending.port);
        } catch (error) {
            setMessage(error instanceof Error ? error.message : String(error));
        } finally {
            setBusy(false);
        }
    }

    function chooseMatchingProfileFromMismatch() {
        if (!pendingProfileMismatch || pendingProfileMismatch.candidates.length === 0) {
            return;
        }
        const pending = pendingProfileMismatch;
        setPendingProfileMismatch(null);
        setPendingStartChoice({
            profiles: pending.candidates,
            address: pending.address,
            port: pending.port,
            carOrdinal: pending.telemetry.carOrdinal,
            carClass: pending.telemetry.carClass,
            carPi: pending.telemetry.carPi,
        });
    }

    async function clearProfileAndStartFromMismatch() {
        if (!pendingProfileMismatch) {
            return;
        }
        setBusy(true);
        setMessage('');
        try {
            const pending = pendingProfileMismatch;
            await SetActiveTuneProfile(0);
            setActiveProfile(null);
            setPendingProfileMismatch(null);
            await startExpertTelemetryNow(pending.address, pending.port);
        } catch (error) {
            setMessage(error instanceof Error ? error.message : String(error));
        } finally {
            setBusy(false);
        }
    }

    async function stop() {
        setBusy(true);
        setMessage('');
        try {
            if (status.mode === 'replay') {
                await StopTelemetryReplay();
            } else {
                await StopTelemetry();
            }
            if (status.analysisMode === 'expert' || status.mode === 'replay') {
                await loadMetadata();
            }
            if (status.analysisMode === 'track_capture') {
                setTrackCapture(capture => ({...capture, recording: false}));
            }
            setMessage(t.stopped);
        } catch (error) {
            setMessage(error instanceof Error ? error.message : String(error));
        } finally {
            setBusy(false);
        }
    }

    async function startTuneWebServer() {
        const parsedPort = parseStrictInteger(tuneWebPort);
        if (parsedPort === null || parsedPort < 1 || parsedPort > 65535) {
            setMessage(t.remoteTunePortInvalid);
            return;
        }
        setBusy(true);
        setMessage('');
        try {
            await StartTuneWebServer(parsedPort);
            const nextStatus = await GetTuneWebServerStatus();
            setTuneWebStatus((nextStatus || {running: false, port: 0, url: '', lanAddress: '', lastError: ''}) as TuneWebServerStatus);
            setMessage(t.remoteTuneRunning);
        } catch (error) {
            const nextStatus = await GetTuneWebServerStatus().catch(() => null);
            if (nextStatus) {
                setTuneWebStatus(nextStatus as TuneWebServerStatus);
            }
            setMessage(error instanceof Error ? error.message : String(error));
        } finally {
            setBusy(false);
        }
    }

    async function stopTuneWebServer() {
        setBusy(true);
        setMessage('');
        try {
            await StopTuneWebServer();
            const nextStatus = await GetTuneWebServerStatus();
            setTuneWebStatus((nextStatus || {running: false, port: 0, url: '', lanAddress: '', lastError: ''}) as TuneWebServerStatus);
            setMessage(t.remoteTuneStopped);
        } catch (error) {
            setMessage(error instanceof Error ? error.message : String(error));
        } finally {
            setBusy(false);
        }
    }

    function newProfile() {
        setNewProfileDraft({...emptyProfileInput, useCase: 'Road'});
        setShowNewProfileModal(true);
    }

    function fillNewProfileDraftFromTelemetry() {
        if (!current) {
            setMessage(t.telemetryFillUnavailable);
            return;
        }
        setNewProfileDraft(draft => withCalculatedPowerToWeight({
            ...draft,
            carOrdinal: current.carOrdinal || draft.carOrdinal,
            carClass: current.carClass || draft.carClass,
            pi: current.carPi || draft.pi,
            drivetrain: current.drivetrain || draft.drivetrain,
            numCylinders: current.numCylinders || draft.numCylinders,
        }));
    }

    async function createProfileFromDraft() {
        if (!newProfileDraft.carName.trim()) {
            setMessage(t.requiredCarName);
            return;
        }
        setBusy(true);
        try {
            const created = await CreateTuneProfile(cleanProfileInput(newProfileDraft) as never) as TuneProfile;
            await SetActiveTuneProfile(created.id);
            setShowNewProfileModal(false);
            setActiveProfile(created);
            setEditingProfileId(created.id);
            setProfileForm(profileToInput(created));
            await loadMetadata();
            await loadProfileSnapshots(created.id);
            setView('expert');
            window.setTimeout(() => profileFormRef.current?.scrollIntoView({behavior: 'smooth', block: 'start'}), 0);
            setMessage(t.saved);
        } catch (error) {
            setMessage(error instanceof Error ? error.message : String(error));
        } finally {
            setBusy(false);
        }
    }

    async function loadTestConditionDefaults() {
        try {
            const defaults = await GetTestConditionDefaults();
            setTestConditions(normalizeTestConditions((defaults || unknownTestConditions) as TestConditions));
        } catch (error) {
            setMessage(error instanceof Error ? error.message : String(error));
        }
    }

    async function loadTireRegressionSamples(preferredId = selectedTireRegressionId) {
        const loaded = (await ListTireRegressionSamples() || []) as TireRegressionSampleSummary[];
        setTireRegressionSamples(loaded);
        const targetId = preferredId && loaded.some(sample => sample.id === preferredId)
            ? preferredId
            : loaded[0]?.id || '';
        setSelectedTireRegressionId(targetId);
        if (!targetId) {
            setSelectedTireRegressionSample(null);
        }
        return loaded;
    }

    async function saveCurrentTireRegressionSample() {
        setBusy(true);
        try {
            const input: TireRegressionSampleInput = {
                name: tireRegressionSaveForm.name,
                scenario: tireRegressionSaveForm.scenario,
                windowSeconds: Math.trunc(Number(tireRegressionSaveForm.windowSeconds) || 15),
                expected: emptyTireRegressionExpectation(),
            };
            const saved = await SaveTireRegressionSample(input as never) as TireRegressionSample;
            await loadTireRegressionSamples(saved.id);
            setSelectedTireRegressionSample(saved);
            setTireRegressionExpectedForm(expectationToForm(saved.expected));
            setMessage(t.tireRegressionSaved);
        } catch (error) {
            setMessage(error instanceof Error ? error.message : String(error));
        } finally {
            setBusy(false);
        }
    }

    async function updateSelectedTireRegressionExpectation() {
        if (!selectedTireRegressionId) {
            return;
        }
        setBusy(true);
        try {
            await UpdateTireRegressionSampleExpectation(selectedTireRegressionId, expectationFromForm(tireRegressionExpectedForm) as never);
            const updated = await GetTireRegressionSample(selectedTireRegressionId) as TireRegressionSample;
            setSelectedTireRegressionSample(updated);
            setTireRegressionExpectedForm(expectationToForm(updated.expected));
            await loadTireRegressionSamples(selectedTireRegressionId);
            setMessage(t.tireRegressionExpectationSaved);
        } catch (error) {
            setMessage(error instanceof Error ? error.message : String(error));
        } finally {
            setBusy(false);
        }
    }

    async function runSelectedTireRegressionSample() {
        if (!selectedTireRegressionId) {
            return;
        }
        setBusy(true);
        try {
            const result = await RunTireRegressionSample(selectedTireRegressionId) as TireRegressionResult;
            setTireRegressionResults([result]);
            setMessage(t.tireRegressionRan);
        } catch (error) {
            setMessage(error instanceof Error ? error.message : String(error));
        } finally {
            setBusy(false);
        }
    }

    async function runAllTireRegressionSamples() {
        setBusy(true);
        try {
            const results = (await RunAllTireRegressionSamples() || []) as TireRegressionResult[];
            setTireRegressionResults(results);
            setMessage(t.tireRegressionRan);
        } catch (error) {
            setMessage(error instanceof Error ? error.message : String(error));
        } finally {
            setBusy(false);
        }
    }

    async function deleteSelectedTireRegressionSample() {
        if (!selectedTireRegressionId) {
            return;
        }
        setBusy(true);
        try {
            await DeleteTireRegressionSample(selectedTireRegressionId);
            setSelectedTireRegressionSample(null);
            setTireRegressionResults(results => results.filter(result => result.sampleId !== selectedTireRegressionId));
            await loadTireRegressionSamples('');
            setMessage(t.tireRegressionDeleted);
        } catch (error) {
            setMessage(error instanceof Error ? error.message : String(error));
        } finally {
            setBusy(false);
        }
    }

    async function runModelPipeline() {
        setBusy(true);
        try {
            const input: TuningPipelineRunInput = {
                ...pipelineInput,
                sessionId: pipelineInput.sourceType === 'telemetry_session' ? (pipelineInput.sessionId || selectedSessionId || 0) : 0,
            };
            const result = await RunTuningModelPipeline(input as never) as TuningPipelineRunResult;
            setPipelineResult(result || null);
            setMessage(result?.warnings?.[0] || t.modelPipelineRunComplete);
        } catch (error) {
            setMessage(error instanceof Error ? error.message : String(error));
        } finally {
            setBusy(false);
        }
    }

    async function saveProfessionalPipelineConfig() {
        setBusy(true);
        try {
            const saved = await SaveProfessionalPipelineConfig(professionalPipelineConfig as never) as ProfessionalPipelineConfig;
            setProfessionalPipelineConfig(saved);
            setMessage(t.saved);
        } catch (error) {
            setMessage(error instanceof Error ? error.message : String(error));
        } finally {
            setBusy(false);
        }
    }

    async function refreshTuneHarvestCandidates(showBusy = true) {
        if (showBusy) {
            setBusy(true);
        }
        try {
            const next = await SearchTuneHarvestCandidates(tuneHarvestStatusFilter, tuneHarvestSearch, tuneHarvestCandidateListLimit);
            setTuneHarvestCandidates((next || []) as TuneHarvestCandidate[]);
        } catch (error) {
            setMessage(error instanceof Error ? error.message : String(error));
        } finally {
            if (showBusy) {
                setBusy(false);
            }
        }
    }

    async function runTuneHarvest() {
        if (tuneHarvestRunning) {
            return;
        }
        const sources = Object.entries(tuneHarvestSources)
            .filter(([, enabled]) => enabled)
            .map(([source]) => source);
        if (sources.length === 0) {
            setMessage(t.tuneHarvestSelectedSourcesRequired);
            return;
        }
        const parsedLimit = Math.trunc(Number(tuneHarvestLimit) || 80);
        setTuneHarvestRunning(true);
        setTuneHarvestStopping(false);
        setBusy(true);
        setMessage('');
        try {
            const result = await RunTuneHarvest({
                sources,
                dryRun: tuneHarvestDryRun,
                limitPerSource: Math.max(1, Math.min(500, parsedLimit)),
            } as TuneHarvestOptions as never) as TuneHarvestRunResult;
            setTuneHarvestResult(result || null);
            if (tuneHarvestDryRun) {
                setTuneHarvestCandidates((result?.candidates || []) as TuneHarvestCandidate[]);
            } else {
                await refreshTuneHarvestCandidates(false);
            }
            if (result?.run?.status === 'cancelled' || (result?.warnings || []).includes('harvest cancelled')) {
                setMessage(t.tuneHarvestStopped);
            } else {
                setMessage(t.tuneHarvestResult(result?.found || 0, result?.saved || 0, result?.pending || 0, result?.rejected || 0));
            }
        } catch (error) {
            setMessage(error instanceof Error ? error.message : String(error));
        } finally {
            setTuneHarvestStopping(false);
            setTuneHarvestRunning(false);
            setBusy(false);
        }
    }

    async function stopTuneHarvest() {
        if (!tuneHarvestRunning || tuneHarvestStopping) {
            return;
        }
        setTuneHarvestStopping(true);
        try {
            await StopTuneHarvest();
            setMessage(t.tuneHarvestStopping);
        } catch (error) {
            setTuneHarvestStopping(false);
            setMessage(error instanceof Error ? error.message : String(error));
        }
    }

    async function clearTuneHarvestCandidates() {
        if (!window.confirm(t.tuneHarvestClearConfirm)) {
            return;
        }
        setBusy(true);
        try {
            const count = await ClearTuneHarvestCandidates();
            setTuneHarvestCandidates([]);
            setTuneHarvestResult(null);
            setImportingHarvestCandidateId(null);
            setMessage(t.tuneHarvestCleared(Number(count || 0)));
        } catch (error) {
            setMessage(error instanceof Error ? error.message : String(error));
        } finally {
            setBusy(false);
        }
    }

    async function updateTuneHarvestCandidate(candidate: TuneHarvestCandidate, status: string, reason: string) {
        setBusy(true);
        try {
            await UpdateTuneHarvestCandidateStatus(candidate.id, status, reason);
            await refreshTuneHarvestCandidates(false);
            setMessage(status === 'rejected' ? t.tuneHarvestRejected : t.tuneHarvestRestored);
        } catch (error) {
            setMessage(error instanceof Error ? error.message : String(error));
        } finally {
            setBusy(false);
        }
    }

    function useTuneHarvestCandidate(candidate: TuneHarvestCandidate) {
        setImportingHarvestCandidateId(candidate.id);
        setEditingRecommendedCarId('');
        setRecommendedCarsFormError('');
        setRecommendedCarForm(recommendedCarFormFromHarvestCandidate(candidate));
        setRecommendedCarFormOpen(true);
        setDeveloperTool('recommended_cars');
        setMessage(t.tuneHarvestCopiedToRecommended);
    }

    function changeRecommendedCarForm(field: keyof RecommendedCarForm, value: string) {
        setRecommendedCarsFormError('');
        setRecommendedCarForm(form => {
            const next = {...form, [field]: value};
    if (field === 'name' || field === 'useCase' || field === 'carClass' || field === 'pi' || field === 'tuneCode') {
        next.id = '';
    }
            if (field === 'useCase') {
                next.useCaseLabel = recommendedUseCaseLabels[value] || '';
            }
            if (field === 'tireCompound') {
                next.tireCompoundLabel = recommendedTireCompoundLabels[value] || '';
            }
            if (field === 'carClass') {
                const previousDefault = recommendedClassDefaultPI[form.carClass];
                if (!form.pi || form.pi === previousDefault) {
                    next.pi = recommendedClassDefaultPI[value] || form.pi;
                }
            }
            if (field === 'pi') {
                const defaultClass = recommendedPIDefaultClass[value.trim()];
                if (defaultClass) {
                    next.carClass = defaultClass;
                }
            }
            return next;
        });
    }

    function openRecommendedCarCreate() {
        setEditingRecommendedCarId('');
        setRecommendedCarsFormError('');
        setRecommendedCarForm(emptyRecommendedCarForm);
        setRecommendedCarFormOpen(true);
    }

    function clearRecommendedCarForm() {
        setEditingRecommendedCarId('');
        setRecommendedCarsFormError('');
        setRecommendedCarForm(emptyRecommendedCarForm);
    }

    async function refreshRecommendedCars() {
        setBusy(true);
        try {
            const [nextCars, nextSelection] = await Promise.all([
                ListRecommendedCars(),
                LoadRecommendedCarsFileSelection().catch(() => null),
            ]);
            const loadedCars = (nextCars || []) as RecommendedCar[];
            setRecommendedCars(loadedCars);
            applyRecommendedCarsFileSelection(loadedCars, (nextSelection || null) as RecommendedCarsFileSelection | null);
            setRecommendedCarsResult(null);
        } catch (error) {
            setMessage(error instanceof Error ? error.message : String(error));
        } finally {
            setBusy(false);
        }
    }

    async function addRecommendedCar() {
        try {
            setRecommendedCarsFormError('');
            const previousId = editingRecommendedCarId;
            const input = recommendedCarInputFromForm(recommendedCarForm);
            const duplicateTuneCode = recommendedCars.find(car => normalizeTuneCodeInput(car.tuneCode) === input.tuneCode && car.id !== previousId);
            if (duplicateTuneCode) {
                throw new Error(t.recommendedCarsDuplicateTuneCode);
            }
            const previewID = generateRecommendedCarPreviewID(input);
            const duplicateIdentity = recommendedCars.find(car => car.id === previewID && car.id !== previousId);
            if (duplicateIdentity) {
                throw new Error(t.recommendedCarsDuplicateIdentity);
            }
            setBusy(true);
            const saved = await SaveRecommendedCarRecord(input as never, previousId) as RecommendedCar;
            if (importingHarvestCandidateId) {
                await UpdateTuneHarvestCandidateStatus(importingHarvestCandidateId, 'imported', saved.id);
                setImportingHarvestCandidateId(null);
            }
            const nextCars = await ListRecommendedCars();
            setRecommendedCars((nextCars || []) as RecommendedCar[]);
            setRecommendedCarsResult(null);
            setSelectedRecommendedCarIds(ids => Array.from(new Set([...ids.filter(id => id !== previousId), saved.id])));
            setRecommendedCarFormOpen(false);
            setEditingRecommendedCarId('');
            setRecommendedCarForm(emptyRecommendedCarForm);
            setMessage(t.recommendedCarsDbSaved);
        } catch (error) {
            const text = error instanceof Error ? error.message : String(error);
            setRecommendedCarsFormError(text);
            setMessage(text);
        } finally {
            setBusy(false);
        }
    }

    async function saveRecommendedCars() {
        setBusy(true);
        try {
            const selectedCars = recommendedCars.filter(car => selectedRecommendedCarIds.includes(car.id));
            if (selectedCars.length === 0) {
                throw new Error(t.recommendedCarsExportEmpty);
            }
            const result = await SaveRecommendedCarsFile(selectedCars.map(recommendedCarExportInput) as never, recommendedCarsVersion) as RecommendedCarsFileResult;
            const nextSelection = await LoadRecommendedCarsFileSelection();
            setRecommendedCarsResult(result);
            applyRecommendedCarsFileSelection(recommendedCars, (nextSelection || null) as RecommendedCarsFileSelection | null);
            setMessage(t.recommendedCarsSaved(result.count, result.path));
        } catch (error) {
            setMessage(error instanceof Error ? error.message : String(error));
        } finally {
            setBusy(false);
        }
    }

    async function deleteRecommendedCar(id: string) {
        setBusy(true);
        try {
            await DeleteRecommendedCar(id);
            const nextCars = await ListRecommendedCars();
            setRecommendedCars((nextCars || []) as RecommendedCar[]);
            setSelectedRecommendedCarIds(ids => ids.filter(item => item !== id));
            setRecommendedCarsResult(null);
            setMessage(t.recommendedCarsDbDeleted);
        } catch (error) {
            setMessage(error instanceof Error ? error.message : String(error));
        } finally {
            setBusy(false);
        }
    }

    async function deleteAllRecommendedCars() {
        if (!window.confirm(t.recommendedCarsDeleteAllConfirm)) {
            return;
        }
        setBusy(true);
        try {
            const deleted = await DeleteAllRecommendedCars();
            const nextSelection = await LoadRecommendedCarsFileSelection().catch(() => null);
            setRecommendedCars([]);
            applyRecommendedCarsFileSelection([], (nextSelection || null) as RecommendedCarsFileSelection | null);
            setRecommendedCarsResult(null);
            setRecommendedCarDetail(null);
            setRecommendedCarFormOpen(false);
            setEditingRecommendedCarId('');
            setRecommendedCarsFormError('');
            setMessage(t.recommendedCarsDbDeletedAll(Number(deleted || 0)));
        } catch (error) {
            setMessage(error instanceof Error ? error.message : String(error));
        } finally {
            setBusy(false);
        }
    }

    function editRecommendedCar(car: RecommendedCar) {
        setEditingRecommendedCarId(car.id);
        setRecommendedCarsFormError('');
        setRecommendedCarForm(recommendedCarFormFromCar(car));
        setRecommendedCarFormOpen(true);
        setRecommendedCarsResult(null);
    }

    function duplicateRecommendedCar(car: RecommendedCar) {
        setEditingRecommendedCarId('');
        setRecommendedCarsFormError(t.recommendedCarsCopied);
        setRecommendedCarForm({
            ...recommendedCarFormFromCar(car),
            id: '',
            tuneCode: '',
        });
        setRecommendedCarFormOpen(true);
        setRecommendedCarsResult(null);
    }

    async function loadMetadata(preferredSessionId = selectedSessionId) {
        void preferredSessionId;
        try {
            const [nextProfiles, nextActive, nextProfileStats, nextRuleProfiles, nextStrategyTemplates, nextUpgradeRules, nextTuneExplanations, nextInfluenceMap, nextTracks, nextTrackProfiles, nextTireRegressionSamples, nextRecommendedCars, nextRecommendedCarsFileSelection, nextKnowledgeStatus, nextPipelineCatalog, nextProfessionalConfig] = await Promise.all([
                ListTuneProfiles(),
                GetActiveTuneProfile(),
                ListTuneProfileSessionStats(),
                ListRuleThresholdProfiles(),
                ListStrategyTemplates(),
                ListUpgradeUnlockRules(),
                ListTuneAdjustmentExplanations(),
                GetTuneToTireInfluenceMap(),
                ListBenchmarkTracks(),
                ListTrackProfiles(),
                ListTireRegressionSamples(),
                ListRecommendedCars(),
                LoadRecommendedCarsFileSelection().catch(() => null),
                GetRoadTuningKnowledgeStatus(),
                ListTuningModelPipelines(),
                GetProfessionalPipelineConfig(),
            ]);
            setProfiles((nextProfiles || []) as TuneProfile[]);
            setActiveProfile((nextActive || null) as TuneProfile | null);
            setProfileSessionStats((nextProfileStats || []) as TuneProfileSessionStat[]);
            setUpgradeUnlockRules((nextUpgradeRules || []) as UpgradeUnlockRule[]);
            setTuneExplanations((nextTuneExplanations || []) as TuneAdjustmentExplanation[]);
            setTuneInfluenceMap((nextInfluenceMap || null) as TuneToTireInfluenceMap | null);
            const loadedTracks = (nextTracks || []) as BenchmarkTrack[];
            setBenchmarkTracks(loadedTracks);
            const loadedTrackProfiles = (nextTrackProfiles || []) as TrackProfile[];
            setTrackProfiles(loadedTrackProfiles);
            const loadedTireSamples = (nextTireRegressionSamples || []) as TireRegressionSampleSummary[];
            setTireRegressionSamples(loadedTireSamples);
            const loadedRecommendedCars = (nextRecommendedCars || []) as RecommendedCar[];
            setRecommendedCars(loadedRecommendedCars);
            applyRecommendedCarsFileSelection(loadedRecommendedCars, (nextRecommendedCarsFileSelection || null) as RecommendedCarsFileSelection | null);
            setPipelineCatalog((nextPipelineCatalog || null) as TuningPipelineCatalog | null);
            setProfessionalPipelineConfig((nextProfessionalConfig || {
                detectorId: defaultPipelineInput.detectorId,
                decisionerId: defaultPipelineInput.decisionerId,
                interpreterId: defaultPipelineInput.interpreterId,
            }) as ProfessionalPipelineConfig);
            setSelectedTireRegressionId(currentId => currentId || loadedTireSamples[0]?.id || '');
            setSelectedTrackId(currentTrackId => currentTrackId || loadedTracks[0]?.id || null);
            setSelectedTrackProfile(currentProfile => {
                const targetId = currentProfile?.track.id || selectedTrackId || loadedTracks[0]?.id || null;
                return loadedTrackProfiles.find(profile => profile.track.id === targetId) || loadedTrackProfiles[0] || null;
            });
            const loadedRuleProfiles = (nextRuleProfiles || []) as RuleThresholdProfile[];
            setRuleProfiles(loadedRuleProfiles);
            setKnowledgeStatus((nextKnowledgeStatus || null) as RoadTuningKnowledgeStatus | null);
            const loadedStrategyTemplates = (nextStrategyTemplates || []) as StrategyTemplate[];
            setStrategyTemplates(loadedStrategyTemplates);
            setSelectedStrategyTemplateId(currentId => currentId || loadedStrategyTemplates[0]?.id || null);
            if (!editingRuleId && ruleForm.configJson === '') {
                const defaultRule = loadedRuleProfiles.find(profile => profile.isDefault) || loadedRuleProfiles[0];
                if (defaultRule) {
                    setRuleForm(ruleProfileToInput(defaultRule));
                    setEditingRuleId(defaultRule.id);
                }
            }
            setSessions([]);
            setSelectedSessionId(null);
            setSessionEvents([]);
            setSessionSamples([]);
            setSessionBenchmarkRuns([]);
            setRoadEvaluation(null);
            setSessionIssueSummary(null);
            setRoadTuningDecision(null);
            setTunePlanDraft(null);
            setRetestEvaluation(null);
            setReportMarkdown('');
        } catch (error) {
            setMessage(error instanceof Error ? error.message : String(error));
        }
    }

    async function loadProfileSnapshots(profileId: number | null) {
        if (!profileId) {
            setProfileSnapshots([]);
            return;
        }
        try {
            const snapshots = await ListTuneProfileSnapshots(profileId);
            setProfileSnapshots((snapshots || []) as TuneProfileSnapshot[]);
        } catch (error) {
            setProfileSnapshots([]);
            setMessage(error instanceof Error ? error.message : String(error));
        }
    }

    async function saveProfile() {
        if (!profileForm.carName.trim()) {
            setMessage(t.requiredCarName);
            return;
        }
        setBusy(true);
        try {
            let savedProfile: TuneProfile;
            if (editingProfileId) {
                savedProfile = await UpdateTuneProfile(editingProfileId, cleanProfileInput(profileForm) as never) as TuneProfile;
            } else {
                const created = await CreateTuneProfile(cleanProfileInput(profileForm) as never);
                savedProfile = created as TuneProfile;
                await SetActiveTuneProfile(savedProfile.id);
            }
            setProfileForm(profileToInput(savedProfile));
            setEditingProfileId(savedProfile.id);
            await loadMetadata();
            await loadProfileSnapshots(savedProfile.id);
            setMessage(t.saved);
        } catch (error) {
            setMessage(error instanceof Error ? error.message : String(error));
        } finally {
            setBusy(false);
        }
    }

    async function generateRoadBaselinePreview(formOverride = roadBaselineForm, options: {quiet?: boolean; closeInput?: boolean; preserveSelection?: boolean} = {}) {
        const activeForm = formOverride;
        if (!quickTuneUseCaseSupported(activeForm.useCase)) {
            setMessage(t.quickTuneUnsupportedUseCase);
            return;
        }
        const validation = validateQuickTuneForm(activeForm, t);
        setQuickTuneFieldErrors(validation.errors);
        if (validation.messages.length > 0) {
            setMessage(validation.messages[0]);
            return;
        }
        setBusy(true);
        try {
            const input = roadBaselineInputFromForm(activeForm);
            const result = await GenerateRoadStaticTuneBaseline(input as never) as RoadStaticTuneBaselineResult;
            setRoadBaselineResult(result);
            const defaultFields = result.generatedFields.filter(field => field.defaultSelected).map(field => String(field.fieldKey));
            if (options.preserveSelection) {
                const availableFields = new Set(result.generatedFields.map(field => String(field.fieldKey)));
                setSelectedBaselineFields(currentFields => {
                    if (currentFields.length === 0) {
                        return [];
                    }
                    const preserved = currentFields.filter(field => availableFields.has(field));
                    return preserved.length > 0 ? preserved : defaultFields;
                });
            } else {
                setSelectedBaselineFields(defaultFields);
            }
            if (options.closeInput !== false) {
                setQuickTuneInputOpen(false);
            }
            if (!options.quiet) {
                setMessage(t.tuneGeneratorPreview);
            }
        } catch (error) {
            setMessage(error instanceof Error ? error.message : String(error));
        } finally {
            setBusy(false);
        }
    }

    async function applyRoadBaseline(createNew: boolean) {
        if (!roadBaselineResult) {
            setMessage(t.tuneGeneratorNoPreview);
            return;
        }
        if (!createNew && !roadBaselineTargetProfileId) {
            setMessage(t.tuneGeneratorNoTarget);
            return;
        }
        if (selectedBaselineFields.length === 0) {
            setMessage(t.tuneGeneratorNoSelection);
            return;
        }
        setBusy(true);
        try {
            const input = roadBaselineInputFromForm(roadBaselineForm);
            const result = await ApplyRoadStaticTuneBaseline({
                createNew,
                targetProfileId: createNew ? 0 : roadBaselineTargetProfileId,
                baselineInput: input,
                selectedFieldKeys: selectedBaselineFields,
            } as never) as RoadStaticTuneBaselineApplyResult;
            await SetActiveTuneProfile(result.profile.id);
            setRoadBaselineTargetProfileId(result.profile.id);
            setEditingProfileId(result.profile.id);
            setProfileForm(profileToInput(result.profile));
            await loadMetadata();
            await loadProfileSnapshots(result.profile.id);
            setMessage(createNew ? t.tuneGeneratorCreated : t.tuneGeneratorApplied);
        } catch (error) {
            setMessage(error instanceof Error ? error.message : String(error));
        } finally {
            setBusy(false);
        }
    }

    function resetRoadBaseline() {
        setRoadBaselineForm(emptyRoadBaselineForm);
        setRoadBaselineResult(null);
        setSelectedBaselineFields([]);
        setRoadBaselineTargetProfileId(0);
        setRoadBaselineAdvancedOpen(false);
        setQuickTuneFieldErrors({});
        setMessage('');
    }

    function updateRoadBaselineField(key: RoadStaticTuneBaselineFormKey, value: string) {
        setRoadBaselineForm(currentForm => {
            const nextForm = {...currentForm, [key]: value};
            if (key === 'useCase') {
                const defaultTireByUseCase: Record<string, string> = {
                    Road: 'sport',
                    Drift: 'drift',
                    Rally: 'rally',
                    Offroad: 'offroad',
                    Drag: 'drag',
                };
                const defaultVersionByUseCase: Record<string, string> = {
                    Road: 'Road Baseline',
                    Drift: 'Drift Baseline',
                    Rally: 'Rally Baseline',
                    Offroad: 'Offroad Baseline',
                    Drag: 'Drag Baseline',
                };
                const knownDefaultTires = new Set(Object.values(defaultTireByUseCase));
                if (!currentForm.tireCompound || knownDefaultTires.has(currentForm.tireCompound)) {
                    nextForm.tireCompound = defaultTireByUseCase[value] || 'sport';
                }
                if (!currentForm.versionName || Object.values(defaultVersionByUseCase).includes(currentForm.versionName)) {
                    nextForm.versionName = defaultVersionByUseCase[value] || 'Road Baseline';
                }
            }
            return nextForm;
        });
        setQuickTuneFieldErrors(currentErrors => {
            if (!currentErrors[key]) {
                return currentErrors;
            }
            const nextErrors = {...currentErrors};
            delete nextErrors[key];
            return nextErrors;
        });
    }

    function updateRoadBaselineBiasField(key: 'balanceBias' | 'stiffnessBias' | 'speedBias', value: string) {
        const nextForm = {...roadBaselineForm, [key]: value};
        setRoadBaselineForm(nextForm);
        setQuickTuneFieldErrors(currentErrors => {
            if (!currentErrors[key]) {
                return currentErrors;
            }
            const nextErrors = {...currentErrors};
            delete nextErrors[key];
            return nextErrors;
        });
        if (!roadBaselineResult) {
            return;
        }
        if (roadBaselineBiasDebounceRef.current !== null) {
            window.clearTimeout(roadBaselineBiasDebounceRef.current);
        }
        roadBaselineBiasDebounceRef.current = window.setTimeout(() => {
            void generateRoadBaselinePreview(nextForm, {quiet: true, closeInput: false, preserveSelection: true});
        }, 180);
    }

    function carryRoadBaselineToProfessional() {
        if (!roadBaselineResult) {
            setMessage(t.tuneGeneratorNoPreview);
            return;
        }
        const draft = {
            ...roadBaselineResult.profileDraft,
            carName: roadBaselineResult.profileDraft.carName || quickTuneAutoProfileName(roadBaselineResult.profileDraft),
        };
        setEditingProfileId(null);
        setProfileForm(draft);
        setView('expert');
        window.setTimeout(() => profileFormRef.current?.scrollIntoView({behavior: 'smooth', block: 'start'}), 0);
        setMessage(t.quickTuneCarriedToProfessional);
    }

    function toggleBaselineField(fieldKey: string) {
        setSelectedBaselineFields(currentFields => currentFields.includes(fieldKey)
            ? currentFields.filter(key => key !== fieldKey)
            : [...currentFields, fieldKey]);
    }

    async function editProfile(profile: TuneProfile) {
        setEditingProfileId(profile.id);
        setProfileForm(profileToInput(profile));
        await loadProfileSnapshots(profile.id);
        setView('expert');
    }

    async function duplicateProfile(profile: TuneProfile) {
        setBusy(true);
        try {
            await DuplicateTuneProfile(profile.id, `${profile.versionName || profile.carName} Copy`);
            await loadMetadata();
            setMessage(t.saved);
        } catch (error) {
            setMessage(error instanceof Error ? error.message : String(error));
        } finally {
            setBusy(false);
        }
    }

    async function restoreProfileSnapshot(snapshot: TuneProfileSnapshot) {
        if (!window.confirm(t.restoreSnapshotConfirm)) {
            return;
        }
        setBusy(true);
        try {
            const restored = await RestoreTuneProfileSnapshot(snapshot.id) as TuneProfile;
            setEditingProfileId(restored.id);
            setProfileForm(profileToInput(restored));
            await loadMetadata();
            await loadProfileSnapshots(restored.id);
            setMessage(t.snapshotRestored);
        } catch (error) {
            setMessage(error instanceof Error ? error.message : String(error));
        } finally {
            setBusy(false);
        }
    }

    async function deleteProfile(profile: TuneProfile) {
        if (!window.confirm(`${t.delete} ${profile.carName}?`)) {
            return;
        }
        setBusy(true);
        try {
            await DeleteTuneProfile(profile.id);
            if (editingProfileId === profile.id) {
                setEditingProfileId(null);
                setProfileForm(emptyProfileInput);
                setProfileSnapshots([]);
            }
            await loadMetadata();
            setMessage(t.deleted);
        } catch (error) {
            setMessage(error instanceof Error ? error.message : String(error));
        } finally {
            setBusy(false);
        }
    }

    async function setProfileActive(profileId: number) {
        setBusy(true);
        try {
            await SetActiveTuneProfile(profileId);
            await loadMetadata();
        } catch (error) {
            setMessage(error instanceof Error ? error.message : String(error));
        } finally {
            setBusy(false);
        }
    }

    async function selectExpertProfile(profileId: number) {
        await setProfileActive(profileId);
        const selected = profiles.find(profile => profile.id === profileId) || null;
        if (selected) {
            await editProfile(selected);
            return;
        }
        setEditingProfileId(null);
        setProfileForm(emptyProfileInput);
        setProfileSnapshots([]);
    }

    async function fillProfileFromTelemetry() {
        if (!current) {
            setMessage(t.telemetryFillUnavailable);
            return;
        }
        let boundName = '';
        if (current.carOrdinal > 0) {
            try {
                boundName = await ResolveCarNameByOrdinal(current.carOrdinal);
            } catch {
                boundName = '';
            }
        }
        setProfileForm(currentForm => ({
            ...currentForm,
            carName: currentForm.carName || boundName,
            carOrdinal: optionalPositiveInt(current.carOrdinal),
            carCategory: optionalNonNegativeInt(current.carCategory),
            carClass: current.carClass || currentForm.carClass,
            pi: optionalPositiveInt(current.carPi),
            drivetrain: current.drivetrain || currentForm.drivetrain,
            numCylinders: optionalPositiveInt(current.numCylinders),
        }));
        setMessage(t.telemetryFilled);
    }

    async function loadReportSessionData(sessionId: number) {
        const [
            eventsForSession,
            samplesForSession,
            runsForSession,
            issueSummaryForSession,
            decisionForSession,
            draftForSession,
            retestForSession,
        ] = await Promise.all([
            GetSessionEvents(sessionId),
            GetSessionTelemetrySamples(sessionId, 2000),
            AnalyzeSessionBenchmarkRuns(sessionId),
            GetSessionIssueSummary(sessionId).catch(() => null),
            GetRoadTuningDecision(sessionId).catch(() => null),
            GetTunePlanDraft(sessionId).catch(() => null),
            GetRetestEvaluation(sessionId).catch(() => null),
        ]);
        const evaluationForSession = await EvaluateRoadSession(sessionId).catch(() => null);
        setSessionEvents((eventsForSession || []) as DetectedEvent[]);
        setSessionSamples((samplesForSession || []) as unknown as TelemetryFrame[]);
        setSessionBenchmarkRuns((runsForSession || []) as BenchmarkRun[]);
        setRoadEvaluation((evaluationForSession || null) as RoadSessionEvaluation | null);
        setSessionIssueSummary((issueSummaryForSession || null) as SessionIssueSummary | null);
        setRoadTuningDecision((decisionForSession || null) as RoadTuningDecision | null);
        setTunePlanDraft((draftForSession || null) as TunePlanDraft | null);
        setRetestEvaluation((retestForSession || null) as RetestEvaluation | null);
    }

    async function selectSession(sessionId: number) {
        setSelectedSessionId(sessionId);
        setReportMarkdown('');
        try {
            await loadReportSessionData(sessionId);
        } catch (error) {
            setMessage(error instanceof Error ? error.message : String(error));
        }
    }

    async function generateReport(sessionId: number) {
        setBusy(true);
        try {
            const report = await GenerateTuningReport(sessionId, language);
            await loadReportSessionData(sessionId);
            setReportMarkdown(report || '');
        } catch (error) {
            setMessage(error instanceof Error ? error.message : String(error));
        } finally {
            setBusy(false);
        }
    }

    function toggleTunePlanAction(actionId: string) {
        setSelectedTunePlanActionIds(currentIds => (
            currentIds.includes(actionId)
                ? currentIds.filter(id => id !== actionId)
                : [...currentIds, actionId]
        ));
    }

    async function applyTunePlan(sessionId: number) {
        if (selectedTunePlanActionIds.length === 0) {
            setMessage(t.tunePlanDraftEmpty);
            return;
        }
        setBusy(true);
        setMessage('');
        try {
            const result = await ApplyTunePlanDraft({sessionId, selectedActionIds: selectedTunePlanActionIds}) as { profile: TuneProfile };
            if (result?.profile) {
                setEditingProfileId(result.profile.id);
                setProfileForm(profileToInput(result.profile));
                await loadProfileSnapshots(result.profile.id);
            }
            await loadMetadata(sessionId);
            setReportMarkdown('');
            setMessage(t.tunePlanApplied);
        } catch (error) {
            setMessage(error instanceof Error ? error.message : String(error));
        } finally {
            setBusy(false);
        }
    }

    async function reloadTuningKnowledge() {
        setBusy(true);
        setMessage('');
        try {
            await ReloadTuningKnowledge();
            const nextStatus = await GetRoadTuningKnowledgeStatus();
            setKnowledgeStatus((nextStatus || null) as RoadTuningKnowledgeStatus | null);
            if (selectedSessionId) {
                await loadReportSessionData(selectedSessionId);
            }
            setMessage(t.knowledgeReloaded);
        } catch (error) {
            const nextStatus = await GetRoadTuningKnowledgeStatus().catch(() => null);
            setKnowledgeStatus((nextStatus || null) as RoadTuningKnowledgeStatus | null);
            setMessage(error instanceof Error ? error.message : String(error));
        } finally {
            setBusy(false);
        }
    }

    async function replaySession(sessionId: number) {
        setBusy(true);
        setMessage('');
        try {
            await ReplayTelemetrySession(sessionId, replaySpeed);
            setReplayStatus(await GetTelemetryReplayStatus() as TelemetryReplayStatus);
            setMessage(`${t.replay} ${replaySpeed}x`);
        } catch (error) {
            setMessage(error instanceof Error ? error.message : String(error));
        } finally {
            setBusy(false);
        }
    }

    async function deleteSession(session: TelemetrySession) {
        if (status.running) {
            setMessage(t.deleteSessionBlocked);
            return;
        }
        const name = session.sessionName || `#${session.id}`;
        if (!window.confirm(t.deleteSessionConfirm(name))) {
            return;
        }
        setBusy(true);
        setMessage('');
        try {
            await DeleteTelemetrySession(session.id);
            if (selectedSessionId === session.id) {
                setSelectedSessionId(null);
                setSessionEvents([]);
                setSessionSamples([]);
                setSessionBenchmarkRuns([]);
                setRoadEvaluation(null);
                setSessionIssueSummary(null);
                setRoadTuningDecision(null);
                setTunePlanDraft(null);
                setRetestEvaluation(null);
                setReportMarkdown('');
            }
            await loadMetadata(null);
            setMessage(t.deleted);
        } catch (error) {
            setMessage(error instanceof Error ? error.message : String(error));
        } finally {
            setBusy(false);
        }
    }

    async function openSessionBind(session: TelemetrySession) {
        setBusy(true);
        setMessage('');
        try {
            const candidates = session.carOrdinal && session.carClass
                ? await ListTuneProfilesForVehicle(session.carOrdinal, session.carClass)
                : [];
            setPendingSessionBind({session, profiles: (candidates || []) as TuneProfile[]});
        } catch (error) {
            setMessage(error instanceof Error ? error.message : String(error));
        } finally {
            setBusy(false);
        }
    }

    async function bindSessionToProfile(profile: TuneProfile) {
        if (!pendingSessionBind) {
            return;
        }
        setBusy(true);
        setMessage('');
        try {
            const sessionId = pendingSessionBind.session.id;
            await BindTelemetrySessionTuneProfile(sessionId, profile.id);
            setPendingSessionBind(null);
            setReportMarkdown('');
            await loadMetadata(sessionId);
            setMessage(t.sessionProfileBound);
        } catch (error) {
            setMessage(error instanceof Error ? error.message : String(error));
        } finally {
            setBusy(false);
        }
    }

    function startTrackCapture() {
        setTrackCapture(capture => ({...capture, recording: true, points: [], hasDrivingLine: false}));
        setMessage('');
    }

    function stopTrackCapture() {
        setTrackCapture(capture => ({...capture, recording: false}));
    }

    async function saveCapturedTrack(trackType: BenchmarkTrackType, extractionMode: BenchmarkExtractionMode, gateWidth: number, gateDepth: number, startGate?: BenchmarkGate, finishGate?: BenchmarkGate) {
        if (trackCapture.points.length < 2) {
            setMessage(t.noTrackPoints);
            return;
        }
        setBusy(true);
        setMessage('');
        try {
            const input = capturedTrackInput(trackCapture, current?.gameMode || 'unknown', trackType, extractionMode, gateWidth, gateDepth, startGate, finishGate);
            const candidates = (await FindSimilarBenchmarkTracks(input as never) || []) as TrackMergeCandidate[];
            const strong = candidates.find(candidate => candidate.matchLevel === 'strong');
            if (strong) {
                const updated = await MergeBenchmarkTrackInput(strong.track.id, input as never) as BenchmarkTrack;
                setTrackCapture({recording: false, name: '', points: [], hasDrivingLine: false});
                setSelectedTrackId(updated.id);
                setPendingTrackMerge(null);
                await loadMetadata(selectedSessionId);
                setMessage(`${t.autoMergedTrack}: ${updated.name}`);
                return;
            }
            const manualCandidates = candidates.filter(candidate => candidate.matchLevel === 'medium');
            if (manualCandidates.length > 0) {
                setPendingTrackMerge({input, candidates: manualCandidates});
                return;
            }
            await saveTrackInputAsNew(input, '');
        } catch (error) {
            setMessage(error instanceof Error ? error.message : String(error));
        } finally {
            setBusy(false);
        }
    }

    async function renameTrack(track: BenchmarkTrack) {
        const nextName = window.prompt(t.renameTrack, track.name);
        if (nextName === null || nextName.trim() === '' || nextName.trim() === track.name) {
            return;
        }
        setBusy(true);
        setMessage('');
        try {
            const updated = await RenameBenchmarkTrack(track.id, nextName.trim()) as BenchmarkTrack;
            setSelectedTrackId(updated.id);
            await loadMetadata(selectedSessionId);
            setMessage(t.trackRenamed);
        } catch (error) {
            setMessage(error instanceof Error ? error.message : String(error));
        } finally {
            setBusy(false);
        }
    }

    async function startTrackBaselineCapture() {
        setBusy(true);
        setMessage('');
        try {
            await startTrackBaselineTelemetryNow(0, selectedAddress, Number(port));
        } catch (error) {
            setMessage(error instanceof Error ? error.message : String(error));
        } finally {
            setBusy(false);
        }
    }

    async function stopTrackBaselineCapture() {
        setBusy(true);
        setMessage('');
        try {
            await StopTrackBaselineTelemetry();
            setMessage(t.trackBaselineStopped);
        } catch (error) {
            setMessage(error instanceof Error ? error.message : String(error));
        } finally {
            setBusy(false);
        }
    }

    async function saveTrackBaseline() {
        setBusy(true);
        setMessage('');
        try {
            const preferredTrack = selectedTrackProfile?.track;
            const result = await SaveTrackBaselineCaptureAuto(
                preferredTrack?.id || 0,
                trackCapture.name || preferredTrack?.name || '',
                preferredTrack?.trackType || 'auto',
            ) as TrackBaselineSaveResult;
            setSelectedTrackId(result.track.id);
            await loadMetadata(selectedSessionId);
            setMessage(result.action === 'matched_existing'
                ? `${t.trackBaselineAutoMatched}: ${result.track.name}`
                : `${t.trackBaselineAutoCreated}: ${result.track.name}`);
        } catch (error) {
            setMessage(error instanceof Error ? error.message : String(error));
        } finally {
            setBusy(false);
        }
    }

    async function deleteTrackBaseline(run: TrackBaselineRun) {
        if (!window.confirm(t.confirmDeleteBaseline)) {
            return;
        }
        setBusy(true);
        setMessage('');
        try {
            await DeleteTrackBaselineRun(run.id);
            await loadMetadata(selectedSessionId);
            setMessage(t.trackBaselineDeleted);
        } catch (error) {
            setMessage(error instanceof Error ? error.message : String(error));
        } finally {
            setBusy(false);
        }
    }

    async function saveTrackInputAsNew(input: BenchmarkTrackInput, sourceSession: string) {
        const saved = await CreateBenchmarkTrack(input as never) as BenchmarkTrack;
        setTrackCapture({recording: false, name: '', points: [], hasDrivingLine: false});
        setSelectedTrackId(saved.id);
        setPendingTrackMerge(null);
        await loadMetadata(selectedSessionId);
        setMessage(formatTrackSavedMessage(saved, sourceSession, t));
    }

    async function savePendingTrackAsNew() {
        if (!pendingTrackMerge) {
            return;
        }
        setBusy(true);
        setMessage('');
        try {
            await saveTrackInputAsNew(pendingTrackMerge.input, '');
        } catch (error) {
            setMessage(error instanceof Error ? error.message : String(error));
        } finally {
            setBusy(false);
        }
    }

    async function mergePendingTrack(candidate: TrackMergeCandidate) {
        if (!pendingTrackMerge) {
            return;
        }
        setBusy(true);
        setMessage('');
        try {
            const updated = await MergeBenchmarkTrackInput(candidate.track.id, pendingTrackMerge.input as never) as BenchmarkTrack;
            setTrackCapture({recording: false, name: '', points: [], hasDrivingLine: false});
            setSelectedTrackId(updated.id);
            setPendingTrackMerge(null);
            await loadMetadata(selectedSessionId);
            setMessage(formatTrackSavedMessage(updated, candidate.track.name, t));
        } catch (error) {
            setMessage(error instanceof Error ? error.message : String(error));
        } finally {
            setBusy(false);
        }
    }

    async function createTrackFromSession(sessionId: number, name: string, trackType: BenchmarkTrackType, extractionMode: BenchmarkExtractionMode, gateWidth: number, gateDepth: number, startGate?: BenchmarkGate, finishGate?: BenchmarkGate) {
        if (!sessionId) {
            return;
        }
        setBusy(true);
        setMessage('');
        try {
            const input: BenchmarkTrackExtractionInput = {sessionId, name, trackType, extractionMode};
            input.startGate = startGate ? applyGateSize(startGate, gateWidth, gateDepth) : emptyGateWithSize(gateWidth, gateDepth);
            input.finishGate = finishGate ? applyGateSize(finishGate, gateWidth, gateDepth) : emptyGateWithSize(gateWidth, gateDepth);
            const saved = await ExtractBenchmarkTrackFromSession(input as never) as BenchmarkTrack;
            setSelectedTrackId(saved.id);
            await loadMetadata(sessionId);
            setMessage(formatTrackSavedMessage(saved, sourceSessionName(sessions, sessionId), t));
        } catch (error) {
            setMessage(error instanceof Error ? error.message : String(error));
        } finally {
            setBusy(false);
        }
    }

    async function reextractTrack(trackId: number, sessionId: number, name: string, trackType: BenchmarkTrackType, extractionMode: BenchmarkExtractionMode, gateWidth: number, gateDepth: number, startGate?: BenchmarkGate, finishGate?: BenchmarkGate) {
        if (!trackId || !sessionId) {
            return;
        }
        setBusy(true);
        setMessage('');
        try {
            const input: BenchmarkTrackExtractionInput = {sessionId, name, trackType, extractionMode};
            input.startGate = startGate ? applyGateSize(startGate, gateWidth, gateDepth) : emptyGateWithSize(gateWidth, gateDepth);
            input.finishGate = finishGate ? applyGateSize(finishGate, gateWidth, gateDepth) : emptyGateWithSize(gateWidth, gateDepth);
            const updated = await ReextractBenchmarkTrack(trackId, input as never) as BenchmarkTrack;
            await loadMetadata(sessionId);
            setSelectedTrackId(updated.id);
            setMessage(formatTrackSavedMessage(updated, sourceSessionName(sessions, sessionId), t));
        } catch (error) {
            setMessage(error instanceof Error ? error.message : String(error));
        } finally {
            setBusy(false);
        }
    }

    async function deleteTrack(track: BenchmarkTrack) {
        if (!window.confirm(`${t.delete} ${track.name}?`)) {
            return;
        }
        setBusy(true);
        setMessage('');
        try {
            await DeleteBenchmarkTrack(track.id);
            setSelectedTrackId(currentTrackId => currentTrackId === track.id ? null : currentTrackId);
            setSelectedTrackRuns([]);
            setSelectedTrackProfile(null);
            await loadMetadata(selectedSessionId);
            setMessage(t.trackDeleted);
        } catch (error) {
            setMessage(error instanceof Error ? error.message : String(error));
        } finally {
            setBusy(false);
        }
    }

    async function selectTrack(trackId: number | null) {
        setSelectedTrackId(trackId);
        if (!trackId) {
            setSelectedTrackRuns([]);
            setSelectedTrackProfile(null);
            return;
        }
        try {
            const [runs, profile] = await Promise.all([
                ListBenchmarkRuns(trackId, 50),
                GetTrackProfile(trackId),
            ]);
            setSelectedTrackRuns((runs || []) as BenchmarkRun[]);
            setSelectedTrackProfile((profile || null) as TrackProfile | null);
        } catch (error) {
            setMessage(error instanceof Error ? error.message : String(error));
        }
    }

    async function analyzeSelectedSessionRuns() {
        if (!selectedSessionId) {
            return;
        }
        setBusy(true);
        try {
            const runs = await AnalyzeSessionBenchmarkRuns(selectedSessionId);
            setSessionBenchmarkRuns((runs || []) as BenchmarkRun[]);
            setRoadEvaluation((await EvaluateRoadSession(selectedSessionId).catch(() => null) || null) as RoadSessionEvaluation | null);
            if (selectedTrackId) {
                const [trackRuns, profile] = await Promise.all([
                    ListBenchmarkRuns(selectedTrackId, 50),
                    GetTrackProfile(selectedTrackId),
                ]);
                setSelectedTrackRuns((trackRuns || []) as BenchmarkRun[]);
                setSelectedTrackProfile((profile || null) as TrackProfile | null);
            }
            setTrackProfiles((await ListTrackProfiles() || []) as TrackProfile[]);
        } catch (error) {
            setMessage(error instanceof Error ? error.message : String(error));
        } finally {
            setBusy(false);
        }
    }

    async function pauseReplay() {
        setBusy(true);
        try {
            await PauseTelemetryReplay();
            setReplayStatus(await GetTelemetryReplayStatus() as TelemetryReplayStatus);
        } catch (error) {
            setMessage(error instanceof Error ? error.message : String(error));
        } finally {
            setBusy(false);
        }
    }

    async function resumeReplay() {
        setBusy(true);
        try {
            await ResumeTelemetryReplay();
            setReplayStatus(await GetTelemetryReplayStatus() as TelemetryReplayStatus);
        } catch (error) {
            setMessage(error instanceof Error ? error.message : String(error));
        } finally {
            setBusy(false);
        }
    }

    async function seekReplay(positionMs: number) {
        try {
            await SeekTelemetryReplay(Math.round(positionMs));
            setReplayStatus(await GetTelemetryReplayStatus() as TelemetryReplayStatus);
        } catch (error) {
            setMessage(error instanceof Error ? error.message : String(error));
        }
    }

    async function compareSessions() {
        if (!compareLeftId || !compareRightId) {
            return;
        }
        setBusy(true);
        try {
            const result = await CompareTelemetrySessions(compareLeftId, compareRightId);
            setSessionComparison(result as SessionComparison);
        } catch (error) {
            setMessage(error instanceof Error ? error.message : String(error));
        } finally {
            setBusy(false);
        }
    }

    async function saveRuleProfile() {
        setBusy(true);
        try {
            if (editingRuleId) {
                await UpdateRuleThresholdProfile(editingRuleId, ruleForm as never);
            } else {
                const created = await CreateRuleThresholdProfile(ruleForm as never);
                setEditingRuleId((created as RuleThresholdProfile).id);
            }
            await loadMetadata(selectedSessionId);
            setMessage(t.saved);
        } catch (error) {
            setMessage(error instanceof Error ? error.message : String(error));
        } finally {
            setBusy(false);
        }
    }

    async function resetRuleProfile(id: number) {
        setBusy(true);
        try {
            const reset = await ResetRuleThresholdProfile(id);
            setRuleForm(ruleProfileToInput(reset as RuleThresholdProfile));
            setEditingRuleId((reset as RuleThresholdProfile).id);
            await loadMetadata(selectedSessionId);
            setMessage(t.saved);
        } catch (error) {
            setMessage(error instanceof Error ? error.message : String(error));
        } finally {
            setBusy(false);
        }
    }

    async function deleteRuleProfile(id: number) {
        const profile = ruleProfiles.find(item => item.id === id);
        if (!profile || profile.isDefault || !window.confirm(`${t.delete} ${profile.name}?`)) {
            return;
        }
        setBusy(true);
        try {
            await DeleteRuleThresholdProfile(id);
            setEditingRuleId(null);
            setRuleForm(emptyRuleThresholdInput);
            await loadMetadata(selectedSessionId);
            setMessage(t.deleted);
        } catch (error) {
            setMessage(error instanceof Error ? error.message : String(error));
        } finally {
            setBusy(false);
        }
    }

    function toggleStrategySession(sessionId: number) {
        setStrategySessionIds(currentIds => {
            if (currentIds.includes(sessionId)) {
                return currentIds.filter(id => id !== sessionId);
            }
            if (currentIds.length >= 5) {
                setMessage(t.strategyAnalysisEmpty);
                return currentIds;
            }
            return [...currentIds, sessionId];
        });
    }

    async function runStrategyAnalysis() {
        if (!selectedStrategyTemplateId || strategySessionIds.length === 0) {
            setMessage(t.strategyAnalysisEmpty);
            return;
        }
        setBusy(true);
        try {
            const result = await AnalyzeRoadStrategySessions(strategySessionIds, selectedStrategyTemplateId);
            setStrategyAnalysis((result || null) as RoadStrategyAnalysis | null);
        } catch (error) {
            setMessage(error instanceof Error ? error.message : String(error));
        } finally {
            setBusy(false);
        }
    }

    return (
        <main className="shell">
            <header className="topbar">
                <div className="app-title">
                    <img className="app-icon" src={appIcon} alt="" aria-hidden="true"/>
                    <h1 className="app-title-text">{t.title}</h1>
                </div>
                <div className="topbar-actions">
                    <div className="language-switch" aria-label={t.languageLabel}>
                        <button className={language === 'en' ? 'active' : ''} onClick={() => setLanguage('en')}>EN</button>
                        <button className={language === 'zh' ? 'active' : ''} onClick={() => setLanguage('zh')}>中文</button>
                    </div>
                    <div className={`status-pill ${status.running ? 'online' : 'offline'}`}>
                        <Radio size={18}/>
                        <span>{status.running ? t.receiving : t.idle}</span>
                    </div>
                </div>
            </header>

            <nav className="view-tabs">
                <button className={view === 'tune_generator' ? 'active' : ''} onClick={() => setView('tune_generator')}>{t.tuneGeneratorTab}</button>
                <button className={view === 'remote_tune' ? 'active' : ''} onClick={() => setView('remote_tune')}>{t.remoteTuneTab}</button>
                <button className={view === 'expert' ? 'active' : ''} onClick={() => setView('expert')}>{t.expertTab}</button>
                <button className={view === 'developer' ? 'active' : ''} onClick={() => setView('developer')}>{t.developerTab}</button>
            </nav>

            {view === 'quick' && (
                <>
                    <TelemetryControlStrip
                        interfaces={interfaces}
                        selectedAddress={selectedAddress}
                        port={port}
                        status={status}
                        busy={busy}
                        message={message}
                        language={language}
                        terms={t}
                        targetAddress={dataOutTargetAddress}
                        onAddressChange={setSelectedAddress}
                        onPortChange={setPort}
                        onStart={() => start('track_capture')}
                        onStop={stop}
                    />
                    <section className="panel launchpad-panel">
                        <div className="panel-heading">
                            <div>
                                <h2>{t.quickDiagnosticTitle}</h2>
                                <span>{t.quickModeNote}</span>
                            </div>
                            <div className="launchpad-heading-actions">
                                <span>{t.quickNoHistory}</span>
                            </div>
                        </div>
                        <div className="launchpad-grid">
                            <TextStat label={t.analysisMode} value={t.quickDiagnosis}/>
                            <TextStat label={t.sessionMode} value={gameModeLabel(currentGameMode, t)}/>
                            <TextStat label={t.vehicleId} value={formatOptionalInt(current?.carOrdinal)}/>
                            <TextStat label={t.classPi} value={current ? `${current.carClass || '--'} / ${formatOptionalInt(current.carPi)}` : '--'}/>
                            <TextStat label={t.recording} value={t.quickNoHistory}/>
                        </div>
                        <div className="status-alerts">
                            <div className="status-alert ok">{t.quickNoHistory}</div>
                        </div>
                    </section>

                    <QuickDiagnosticPanel diagnostic={quickDiagnostic} terms={t} language={language}/>

                    <section className="metrics-grid">
                        <MetricCard icon={<Gauge size={22}/>} label={t.speed} value={formatNumber(current?.speedKmh, 0)} unit="km/h"/>
                        <MetricCard icon={<Zap size={22}/>} label={t.rpm} value={formatNumber(current?.rpm, 0)} unit="rpm"/>
                        <MetricCard icon={<Activity size={22}/>} label={t.gear} value={current ? String(current.gear) : '--'} unit=""/>
                        <MetricCard icon={<Waves size={22}/>} label={t.yawRate} value={formatNumber(current?.yawRate, 2)} unit="rad/s"/>
                    </section>

                    <section className="dashboard-grid">
                        <div className="panel gauges-panel">
                            <div className="panel-heading">
                                <h2>{t.driverInputs}</h2>
                                <span>{gameModeLabel(currentGameMode, t)}</span>
                            </div>
                            <ControlBar label={t.throttle} value={current?.throttle01 ?? 0} tone="green"/>
                            <ControlBar label={t.brake} value={current?.brake01 ?? 0} tone="red"/>
                            <ControlBar label={t.handbrake} value={current?.handBrake01 ?? 0} tone="red"/>
                            <SteerMeter label={t.steering} value={current?.steer01 ?? 0}/>
                            <div className="stats-row">
                                <Stat label={t.rawPackets} value={status.rawPackets}/>
                                <Stat label={t.validPackets} value={status.validPackets}/>
                                <Stat label={t.invalidPackets} value={status.invalidPackets}/>
                                <Stat label={t.parseErrors} value={status.parseErrors}/>
                            </div>
                        </div>

                        <div className="panel trend-panel">
                            <div className="panel-heading">
                                <h2>{t.coreTelemetry}</h2>
                                <span>{t.trendSubtitle}</span>
                            </div>
                            <Trend title={t.speed} unit="km/h" points={speedPoints}/>
                            <Trend title={t.rpmLoad} unit="%" points={rpmPoints}/>
                            <div className="connection-details">
                                <div>
                                    <span>{t.endpoint}</span>
                                    <strong>{status.address}:{status.port}</strong>
                                </div>
                                <div>
                                    <span>{t.mode}</span>
                                    <strong>{modeLabel(status.mode, t)}</strong>
                                </div>
                                <div>
                                    <span>{t.packetSize}</span>
                                    <strong>{status.packetLength} bytes</strong>
                                </div>
                                <div>
                                    <span>{t.lastPacket}</span>
                                    <strong>{formatTime(status.lastPacketAt, t.never)}</strong>
                                </div>
                                <div>
                                    <span>{t.lastUdpPacket}</span>
                                    <strong>{status.lastDatagramAt ? `${formatTime(status.lastDatagramAt, t.never)} / ${status.lastDatagramBytes} bytes` : t.never}</strong>
                                </div>
                                <div>
                                    <span>{t.lastUdpRemote}</span>
                                    <strong>{status.lastDatagramRemote || '--'}</strong>
                                </div>
                                <div>
                                    <span>{t.vehicleId}</span>
                                    <strong>{formatOptionalInt(current?.carOrdinal)}</strong>
                                </div>
                                <div>
                                    <span>{t.classPi}</span>
                                    <strong>{current ? `${current.carClass || '--'} / ${formatOptionalInt(current.carPi)}` : '--'}</strong>
                                </div>
                                <div>
                                    <span>{t.drivetrainCylinders}</span>
                                    <strong>{current ? `${current.drivetrain || '--'} / ${formatOptionalInt(current.numCylinders)}` : '--'}</strong>
                                </div>
                                <div>
                                    <span>{t.recordingSize}</span>
                                    <strong>{formatBytes(status.recordingBytes)} / {formatBytes(status.recordingLimitBytes)}</strong>
                                </div>
                            </div>
                        </div>
                    </section>
                </>
            )}

            {view === 'tire_lab' && (
                <>
                    <TelemetryControlStrip
                        interfaces={interfaces}
                        selectedAddress={selectedAddress}
                        port={port}
                        status={status}
                        busy={busy}
                        message={message}
                        language={language}
                        terms={t}
                        targetAddress={dataOutTargetAddress}
                        onAddressChange={setSelectedAddress}
                        onPortChange={setPort}
                        onStart={() => start('tire_lab')}
                        onStop={stop}
                    />
                    <TireModelLabView diagnostic={tireModelDiagnostic} status={status} current={current} terms={t} language={language} influenceMap={tuneInfluenceMap}/>
                </>
            )}

            {view === 'tire_regression' && (
                <TireRegressionLabView
                    samples={tireRegressionSamples}
                    selectedId={selectedTireRegressionId}
                    selectedSample={selectedTireRegressionSample}
                    results={tireRegressionResults}
                    saveForm={tireRegressionSaveForm}
                    expectedForm={tireRegressionExpectedForm}
                    busy={busy}
                    terms={t}
                    language={language}
                    onSaveFormChange={setTireRegressionSaveForm}
                    onExpectedFormChange={setTireRegressionExpectedForm}
                    onSelect={setSelectedTireRegressionId}
                    onSaveCurrent={saveCurrentTireRegressionSample}
                    onUpdateExpected={updateSelectedTireRegressionExpectation}
                    onRunSelected={runSelectedTireRegressionSample}
                    onRunAll={runAllTireRegressionSamples}
                    onDelete={deleteSelectedTireRegressionSample}
                />
            )}

            {view === 'model_pipeline_lab' && (
                <ModelPipelineLabView
                    catalog={pipelineCatalog}
                    input={pipelineInput}
                    result={pipelineResult}
                    sessions={sessions}
                    selectedSessionId={selectedSessionId}
                    busy={busy}
                    terms={t}
                    language={language}
                    onInputChange={setPipelineInput}
                    onRun={runModelPipeline}
                />
            )}

            {view === 'track_profiles' && (
                <>
                    <TelemetryControlStrip
                        interfaces={interfaces}
                        selectedAddress={selectedAddress}
                        port={port}
                        status={status}
                        busy={busy}
                        message={message}
                        language={language}
                        terms={t}
                        targetAddress={dataOutTargetAddress}
                        onAddressChange={setSelectedAddress}
                        onPortChange={setPort}
                        onStart={() => start('track_capture')}
                        onStop={stop}
                    />
                    <TrackProfilesView
                        current={current}
                        status={status}
                        tracks={benchmarkTracks}
                        profiles={trackProfiles}
                        selectedTrackId={selectedTrackId}
                        selectedTrackProfile={selectedTrackProfile}
                        capture={trackCapture}
                        terms={t}
                        busy={busy}
                        onCaptureNameChange={(name) => setTrackCapture(capture => ({...capture, name}))}
                        onStartCapture={startTrackCapture}
                        onStopCapture={stopTrackCapture}
                        onSaveCapture={saveCapturedTrack}
                        onSelectTrack={selectTrack}
                        onRenameTrack={renameTrack}
                        onDeleteTrack={deleteTrack}
                        onStartBaseline={startTrackBaselineCapture}
                        onStopBaseline={stopTrackBaselineCapture}
                        onSaveBaseline={saveTrackBaseline}
                        onDeleteBaseline={deleteTrackBaseline}
                    />
                </>
            )}

            {view === 'tune_generator' && (
                <TuneGeneratorView
                    form={roadBaselineForm}
                    result={roadBaselineResult}
                    profiles={profiles}
                    selectedFields={selectedBaselineFields}
                    targetProfileId={roadBaselineTargetProfileId}
                    advancedOpen={roadBaselineAdvancedOpen}
                    fieldErrors={quickTuneFieldErrors}
                    busy={busy}
                    terms={t}
                    language={language}
                    onFieldChange={updateRoadBaselineField}
                    onBiasChange={updateRoadBaselineBiasField}
                    onFieldBlur={(key, step) => setRoadBaselineForm(currentForm => ({...currentForm, [key]: formatBaselineFormNumber(currentForm[key], step)}))}
                    onAdvancedToggle={() => setRoadBaselineAdvancedOpen(value => !value)}
                    inputOpen={quickTuneInputOpen}
                    onInputOpen={() => setQuickTuneInputOpen(true)}
                    onInputClose={() => setQuickTuneInputOpen(false)}
                    onGenerate={() => generateRoadBaselinePreview()}
                    onCreate={() => applyRoadBaseline(true)}
                    onApply={() => applyRoadBaseline(false)}
                    onReset={resetRoadBaseline}
                    onToggleField={toggleBaselineField}
                    onTargetChange={setRoadBaselineTargetProfileId}
                    onOpenExpert={carryRoadBaselineToProfessional}
                />
            )}

            {view === 'remote_tune' && (
                <RemoteTuneView
                    status={tuneWebStatus}
                    port={tuneWebPort}
                    busy={busy}
                    terms={t}
                    onPortChange={setTuneWebPort}
                    onStart={startTuneWebServer}
                    onStop={stopTuneWebServer}
                />
            )}

            {view === 'expert' && (
                <>
                    <TelemetryControlStrip
                        interfaces={interfaces}
                        selectedAddress={selectedAddress}
                        port={port}
                        status={status}
                        busy={busy}
                        message={message}
                        language={language}
                        terms={t}
                        targetAddress={dataOutTargetAddress}
                        onAddressChange={setSelectedAddress}
                        onPortChange={setPort}
                        onStart={() => start('professional')}
                        onStop={stop}
                    />
                    <section className="panel expert-workflow-panel">
                        <div className="panel-heading">
                            <div>
                                <h2>{t.expertWorkspace}</h2>
                                <span>{t.expertStartHint}</span>
                            </div>
                            <label className="profile-picker embedded">
                                <span>{t.currentProfile}</span>
                                <select
                                    value={activeProfile?.id || 0}
                                    onChange={(event) => selectExpertProfile(Number(event.target.value))}
                                    disabled={busy}
                                >
                                    <option value={0}>{t.noProfile}</option>
                                    {profiles.map(profile => (
                                        <option key={profile.id} value={profile.id}>
                                            {[profile.carName, profile.versionName, profile.carClass, localizedUseCase(profile.useCase, t)].filter(Boolean).join(' / ')}
                                        </option>
                                    ))}
                                </select>
                            </label>
                        </div>
                        <div className="launchpad-grid">
                            <TextStat label={t.analysisMode} value={t.expertTab}/>
                            <TextStat label={t.currentProfile} value={activeProfile ? formatTuneProfileVehicle(activeProfile, t) : t.noProfile}/>
                            <TextStat label={t.modelPipelineDetector} value={professionalPipelineConfig.detectorId}/>
                            <TextStat label={t.modelPipelineDecisioner} value={professionalPipelineConfig.decisionerId}/>
                            <TextStat label={t.modelPipelineInterpreter} value={professionalPipelineConfig.interpreterId}/>
                        </div>
                        <div className="status-alerts">
                            {activeProfile && activeProfileMatchState === 'mismatch' && (
                                <div className="status-alert warn">{t.activeProfileMismatch}</div>
                            )}
                            {activeProfile && activeProfileMatchState === 'match' && (
                                <div className="status-alert ok">{t.activeProfileReady}</div>
                            )}
                        </div>
                    </section>
                    <ProfessionalTuningDiagnosticPanel
                        diagnostic={professionalDiagnostic}
                        profile={profileInputToProfile(profileForm, activeProfile)}
                        terms={t}
                        language={language}
                    />
                    <ProfileManager
                        profiles={profiles}
                        activeProfile={activeProfile}
                        sessionStats={profileSessionStats}
                        snapshots={profileSnapshots}
                        profileForm={profileForm}
                        current={current}
                        formRef={profileFormRef}
                        editingProfileId={editingProfileId}
                        busy={busy}
                        language={language}
                        terms={t}
                        influenceMap={tuneInfluenceMap}
                        onFieldChange={(field, value) => setProfileForm(currentForm => withCalculatedPowerToWeight({...currentForm, [field]: value}))}
                        onNew={newProfile}
                        onEdit={editProfile}
                        onDuplicate={duplicateProfile}
                        onRestoreSnapshot={restoreProfileSnapshot}
                        onDelete={deleteProfile}
                        onSetActive={setProfileActive}
                        onFillFromTelemetry={fillProfileFromTelemetry}
                        canFillFromTelemetry={!!current}
                        onSave={saveProfile}
                    />
                </>
            )}

            {view === 'reports' && (
                <ReportsView
                    sessions={sessions}
                    profiles={profiles}
                    tuneExplanations={tuneExplanations}
                    selectedSessionId={selectedSessionId}
                    events={sessionEvents}
                    issueSummary={sessionIssueSummary}
                    roadTuningDecision={roadTuningDecision}
                    tunePlanDraft={tunePlanDraft}
                    selectedTunePlanActionIds={selectedTunePlanActionIds}
                    retestEvaluation={retestEvaluation}
                    samples={sessionSamples}
                    benchmarkRuns={sessionBenchmarkRuns}
                    roadEvaluation={roadEvaluation}
                    reportMarkdown={reportMarkdown}
                    status={status}
                    replaySpeed={replaySpeed}
                    replayStatus={replayStatus}
                    compareLeftId={compareLeftId}
                    compareRightId={compareRightId}
                    comparison={sessionComparison}
                    language={language}
                    terms={t}
                    onSelectSession={selectSession}
                    onGenerateReport={generateReport}
                    onReplay={replaySession}
                    onDeleteSession={deleteSession}
                    onBindSession={openSessionBind}
                    onToggleTunePlanAction={toggleTunePlanAction}
                    onApplyTunePlan={applyTunePlan}
                    onReplaySpeedChange={setReplaySpeed}
                    onStopReplay={stop}
                    onPauseReplay={pauseReplay}
                    onResumeReplay={resumeReplay}
                    onSeekReplay={seekReplay}
                    onCompareLeftChange={setCompareLeftId}
                    onCompareRightChange={setCompareRightId}
                    onCompare={compareSessions}
                    busy={busy}
                />
            )}

            {view === 'developer' && (
                <>
                    <section className="panel">
                        <div className="panel-heading">
                            <div>
                                <h2>{t.developerTab}</h2>
                                <span>{t.modelPipelineExplainOnly}</span>
                            </div>
                        </div>
                        <div className="view-tabs nested-tabs">
                            <button className={developerTool === 'do_fields' ? 'active' : ''} onClick={() => setDeveloperTool('do_fields')}>{t.developerDoDiagnostics}</button>
                            <button className={developerTool === 'tire_lab' ? 'active' : ''} onClick={() => setDeveloperTool('tire_lab')}>{t.tireLabTab}</button>
                            <button className={developerTool === 'tire_regression' ? 'active' : ''} onClick={() => setDeveloperTool('tire_regression')}>{t.tireRegressionTab}</button>
                            <button className={developerTool === 'model_pipeline_lab' ? 'active' : ''} onClick={() => setDeveloperTool('model_pipeline_lab')}>{t.modelPipelineTab}</button>
                            <button className={developerTool === 'track_profiles' ? 'active' : ''} onClick={() => setDeveloperTool('track_profiles')}>{t.trackProfilesTab}</button>
                            <button className={developerTool === 'recommended_cars' ? 'active' : ''} onClick={() => setDeveloperTool('recommended_cars')}>{t.recommendedCarsTab}</button>
                            <button className={developerTool === 'tune_harvest' ? 'active' : ''} onClick={() => setDeveloperTool('tune_harvest')}>{t.tuneHarvestTab}</button>
                            <button className={developerTool === 'strategy' ? 'active' : ''} onClick={() => setDeveloperTool('strategy')}>{t.developerStrategyConfig}</button>
                        </div>
                    </section>

                    {developerTool === 'tire_lab' && (
                        <>
                            <TelemetryControlStrip
                                interfaces={interfaces}
                                selectedAddress={selectedAddress}
                                port={port}
                                status={status}
                                busy={busy}
                                message={message}
                                language={language}
                                terms={t}
                                targetAddress={dataOutTargetAddress}
                                onAddressChange={setSelectedAddress}
                                onPortChange={setPort}
                                onStart={() => start('tire_lab')}
                                onStop={stop}
                            />
                            <TireModelLabView diagnostic={tireModelDiagnostic} status={status} current={current} terms={t} language={language} influenceMap={tuneInfluenceMap}/>
                        </>
                    )}

                    {developerTool === 'tire_regression' && (
                        <TireRegressionLabView
                            samples={tireRegressionSamples}
                            selectedId={selectedTireRegressionId}
                            selectedSample={selectedTireRegressionSample}
                            results={tireRegressionResults}
                            saveForm={tireRegressionSaveForm}
                            expectedForm={tireRegressionExpectedForm}
                            busy={busy}
                            terms={t}
                            language={language}
                            onSaveFormChange={setTireRegressionSaveForm}
                            onExpectedFormChange={setTireRegressionExpectedForm}
                            onSelect={setSelectedTireRegressionId}
                            onSaveCurrent={saveCurrentTireRegressionSample}
                            onUpdateExpected={updateSelectedTireRegressionExpectation}
                            onRunSelected={runSelectedTireRegressionSample}
                            onRunAll={runAllTireRegressionSamples}
                            onDelete={deleteSelectedTireRegressionSample}
                        />
                    )}

                    {developerTool === 'model_pipeline_lab' && (
                        <ModelPipelineLabView
                            catalog={pipelineCatalog}
                            input={pipelineInput}
                            result={pipelineResult}
                            sessions={sessions}
                            selectedSessionId={selectedSessionId}
                            busy={busy}
                            terms={t}
                            language={language}
                            onInputChange={setPipelineInput}
                            onRun={runModelPipeline}
                        />
                    )}

                    {developerTool === 'track_profiles' && (
                        <>
                            <TelemetryControlStrip
                                interfaces={interfaces}
                                selectedAddress={selectedAddress}
                                port={port}
                                status={status}
                                busy={busy}
                                message={message}
                                language={language}
                                terms={t}
                                targetAddress={dataOutTargetAddress}
                                onAddressChange={setSelectedAddress}
                                onPortChange={setPort}
                                onStart={() => start('track_capture')}
                                onStop={stop}
                            />
                            <TrackProfilesView
                                current={current}
                                status={status}
                                tracks={benchmarkTracks}
                                profiles={trackProfiles}
                                selectedTrackId={selectedTrackId}
                                selectedTrackProfile={selectedTrackProfile}
                                capture={trackCapture}
                                terms={t}
                                busy={busy}
                                onCaptureNameChange={(name) => setTrackCapture(capture => ({...capture, name}))}
                                onStartCapture={startTrackCapture}
                                onStopCapture={stopTrackCapture}
                                onSaveCapture={saveCapturedTrack}
                                onSelectTrack={selectTrack}
                                onRenameTrack={renameTrack}
                                onDeleteTrack={deleteTrack}
                                onStartBaseline={startTrackBaselineCapture}
                                onStopBaseline={stopTrackBaselineCapture}
                                onSaveBaseline={saveTrackBaseline}
                                onDeleteBaseline={deleteTrackBaseline}
                            />
                        </>
                    )}

                    {developerTool === 'recommended_cars' && (
                        <RecommendedCarsGeneratorView
                            form={recommendedCarForm}
                            cars={recommendedCars}
                            result={recommendedCarsResult}
                            fileSelection={recommendedCarsFileSelection}
                            version={recommendedCarsVersion}
                            terms={t}
                            busy={busy}
                            formError={recommendedCarsFormError}
                            language={language}
                            selectedIds={selectedRecommendedCarIds}
                            formOpen={recommendedCarFormOpen}
                            isEditing={Boolean(editingRecommendedCarId)}
                            detailCar={recommendedCarDetail}
                            onFormChange={changeRecommendedCarForm}
                            onVersionChange={(value) => setRecommendedCarsVersion(value)}
                            onSelectionChange={setSelectedRecommendedCarIds}
                            onOpenCreate={openRecommendedCarCreate}
                            onRefresh={refreshRecommendedCars}
                            onCloseForm={() => {
                                setRecommendedCarFormOpen(false);
                                setRecommendedCarsFormError('');
                                setImportingHarvestCandidateId(null);
                            }}
                            onAdd={addRecommendedCar}
                            onClearForm={clearRecommendedCarForm}
                            onEdit={editRecommendedCar}
                            onDuplicate={duplicateRecommendedCar}
                            onShowDetail={setRecommendedCarDetail}
                            onCloseDetail={() => setRecommendedCarDetail(null)}
                            onRemove={deleteRecommendedCar}
                            onDeleteAll={deleteAllRecommendedCars}
                            onSave={saveRecommendedCars}
                        />
                    )}

                    {developerTool === 'tune_harvest' && (
                        <TuneHarvestView
                            candidates={tuneHarvestCandidates}
                            result={tuneHarvestResult}
                            sources={tuneHarvestSources}
                            dryRun={tuneHarvestDryRun}
                            limit={tuneHarvestLimit}
                            statusFilter={tuneHarvestStatusFilter}
                            search={tuneHarvestSearch}
                            terms={t}
                            language={language}
                            busy={busy}
                            running={tuneHarvestRunning}
                            stopping={tuneHarvestStopping}
                            onSourcesChange={setTuneHarvestSources}
                            onDryRunChange={setTuneHarvestDryRun}
                            onLimitChange={setTuneHarvestLimit}
                            onStatusFilterChange={setTuneHarvestStatusFilter}
                            onSearchChange={setTuneHarvestSearch}
                            onRun={runTuneHarvest}
                            onStop={stopTuneHarvest}
                            onRefresh={() => refreshTuneHarvestCandidates()}
                            onClear={clearTuneHarvestCandidates}
                            onUse={useTuneHarvestCandidate}
                            onReject={(candidate) => updateTuneHarvestCandidate(candidate, 'rejected', 'manual_reject')}
                            onRestore={(candidate) => updateTuneHarvestCandidate(candidate, 'pending', '')}
                        />
                    )}

                    {(developerTool === 'do_fields' || developerTool === 'strategy') && (
                        <>
                            {developerTool === 'strategy' && (
                                <ProfessionalPipelineConfigPanel
                                    catalog={pipelineCatalog}
                                    config={professionalPipelineConfig}
                                    terms={t}
                                    busy={busy}
                                    onChange={setProfessionalPipelineConfig}
                                    onSave={saveProfessionalPipelineConfig}
                                />
                            )}
                            <DeveloperModeView
                                current={current}
                                status={status}
                                replayStatus={replayStatus}
                                issueSummary={sessionIssueSummary}
                                sessions={sessions}
                                ruleProfiles={ruleProfiles}
                                strategyTemplates={strategyTemplates}
                                strategyAnalysis={strategyAnalysis}
                                strategySessionIds={strategySessionIds}
                                selectedStrategyTemplateId={selectedStrategyTemplateId}
                                knowledgeStatus={knowledgeStatus}
                                ruleForm={ruleForm}
                                editingRuleId={editingRuleId}
                                language={language}
                                terms={t}
                                busy={busy}
                                onRuleSelect={(profile) => {
                                    setEditingRuleId(profile.id);
                                    setRuleForm(ruleProfileToInput(profile));
                                }}
                                onRuleNew={() => {
                                    setEditingRuleId(null);
                                    setRuleForm({...emptyRuleThresholdInput, configJson: ruleProfiles.find(profile => profile.isDefault)?.configJson || ''});
                                }}
                                onRuleFormChange={(field, value) => setRuleForm(currentForm => ({...currentForm, [field]: value}))}
                                onRuleSave={saveRuleProfile}
                                onRuleReset={() => editingRuleId && resetRuleProfile(editingRuleId)}
                                onRuleDelete={() => editingRuleId && deleteRuleProfile(editingRuleId)}
                                onStrategyTemplateChange={setSelectedStrategyTemplateId}
                                onToggleStrategySession={toggleStrategySession}
                                onRunStrategyAnalysis={runStrategyAnalysis}
                                onReloadKnowledge={reloadTuningKnowledge}
                            />
                        </>
                    )}
                </>
            )}

            {showNewProfileModal && (
                <NewProfileModal
                    draft={newProfileDraft}
                    current={current}
                    terms={t}
                    language={language}
                    busy={busy}
                    onChange={(field, value) => setNewProfileDraft(draft => withCalculatedPowerToWeight({...draft, [field]: value}))}
                    onFillFromTelemetry={fillNewProfileDraftFromTelemetry}
                    onCreate={createProfileFromDraft}
                    onCancel={() => setShowNewProfileModal(false)}
                />
            )}

            {pendingStartChoice && (
                <ProfileChoiceModal
                    pending={pendingStartChoice}
                    terms={t}
                    busy={busy}
                    onChoose={startWithProfile}
                    onCancel={() => setPendingStartChoice(null)}
                />
            )}

            {pendingProfileMismatch && (
                <ProfileMismatchModal
                    pending={pendingProfileMismatch}
                    terms={t}
                    busy={busy}
                    onChooseMatching={chooseMatchingProfileFromMismatch}
                    onClearAndStart={clearProfileAndStartFromMismatch}
                    onCancel={() => setPendingProfileMismatch(null)}
                />
            )}

            {pendingSessionBind && (
                <SessionTuneProfileModal
                    pending={pendingSessionBind}
                    terms={t}
                    busy={busy}
                    onChoose={bindSessionToProfile}
                    onCancel={() => setPendingSessionBind(null)}
                />
            )}

            {pendingTrackMerge && (
                <TrackMergeModal
                    pending={pendingTrackMerge}
                    terms={t}
                    busy={busy}
                    onMerge={mergePendingTrack}
                    onSaveNew={savePendingTrackAsNew}
                    onCancel={() => setPendingTrackMerge(null)}
                />
            )}
        </main>
    )
}

function MetricCard({icon, label, value, unit}: { icon: JSX.Element; label: string; value: string; unit: string }) {
    return (
        <div className="metric-card">
            <div className="metric-icon">{icon}</div>
            <span>{label}</span>
            <strong>{value}</strong>
            {unit && <small>{unit}</small>}
        </div>
    );
}

function ProfessionalTuningDiagnosticPanel({
    diagnostic,
    profile,
    terms,
    language,
}: {
    diagnostic: ProfessionalTuningDiagnostic | null;
    profile: TuneProfile | null;
    terms: Copy;
    language: Lang;
}) {
    const pipeline = diagnostic?.pipeline || null;
    const rows = pipeline ? buildProfessionalDiagnosticRows(pipeline) : [];
    return (
        <section className="panel">
            <div className="panel-heading">
                <div>
                    <h2>{terms.professionalDiagnosticTitle}</h2>
                    <span>{terms.professionalDiagnosticHint}</span>
                </div>
                <span>{diagnostic?.status || terms.modelPipelineNoResult}</span>
            </div>
            {!pipeline ? (
                <div className="empty-events advice-placeholder">{terms.professionalDiagnosticEmpty}</div>
            ) : (
                <div className="stacked-panels">
                    <section className="quick-section">
                        <div className="panel-heading compact">
                            <h2>{terms.professionalMergedTitle}</h2>
                            <span>{terms.professionalMergedHint}</span>
                        </div>
                        {rows.length === 0 ? (
                            <div className="empty-events">{terms.professionalMergedEmpty}</div>
                        ) : (
                            <div className="quick-suggestion-list professional-diagnostic-list">
                                {rows.map(row => (
                                    <ProfessionalDiagnosticRowView
                                        key={row.key}
                                        row={row}
                                        profile={profile}
                                        terms={terms}
                                        language={language}
                                    />
                                ))}
                            </div>
                        )}
                    </section>
                    {(diagnostic?.warnings?.length || 0) > 0 && (
                        <div className="status-alert warn">{diagnostic?.warnings?.slice(0, 3).join(' / ')}</div>
                    )}
                </div>
            )}
        </section>
    );
}

type ProfessionalDiagnosticRow = {
    key: string;
    problems: TuningProblem[];
    decisions: TuningDecision[];
    advice: TuningAdvice[];
    count: number;
    durationMs: number;
    confidence: string;
    riskLevel: string;
};

function ProfessionalDiagnosticRowView({
    row,
    profile,
    terms,
    language,
}: {
    row: ProfessionalDiagnosticRow;
    profile: TuneProfile | null;
    terms: Copy;
    language: Lang;
}) {
    const problem = row.problems[0];
    const decision = bestProfessionalDecision(row.decisions);
    const advice = sortProfessionalAdvice(row.advice, problem, decision, profile).slice(0, 3);
    const adjustments = uniqueProfessionalAdjustments(
        advice.flatMap(item => tuningAdviceConcreteAdjustments(item, problem, decision, profile, language))
    ).slice(0, 4);
    const title = problem
        ? pipelineProblemLabel(problem, terms)
        : advice[0]
            ? `${pipelineAdviceCategoryLabel(advice[0].category, terms)} / ${pipelineAdviceDirectionLabel(advice[0].direction, terms)}`
            : decision
                ? pipelineCauseLabel(decision.primaryCause || decision.id, terms)
                : terms.professionalMergedProblem;
    const instruction = formatProfessionalAdviceInstruction(advice[0], problem, decision, terms, language);
    const operationTags = Array.from(new Set(row.problems.flatMap(item => item.operationTags || []))).slice(0, 3);
    const phaseSummary = professionalRowPhaseSummary(row, terms);
    const axleSummary = professionalRowAxleSummary(row, terms);
    return (
        <div className="quick-suggestion-row professional-diagnostic-row">
            <strong>{terms.professionalMergedProblem}: {title}</strong>
            <span>
                {phaseSummary}
                {' / '}{axleSummary}
                {' / '}{terms.modelPipelineConfidence}: {pipelineConfidenceLabel(row.confidence, terms)}
                {' / '}{terms.eventsSaved}: {row.count}
                {' / '}{formatDuration(row.durationMs || 0, terms)}
            </span>
            {operationTags.length > 0 && (
                <small>{terms.tireIssueOperations}: {operationTags.map(tag => localizedLabel(tag, terms.tireOperationTagLabels)).join(' / ')}</small>
            )}
            {decision && (
                <small>
                    {terms.professionalMergedDecision}: {pipelineCauseLabel(decision.primaryCause || decision.id, terms)}
                    {' / '}{terms.modelPipelineShouldTune}: {decision.shouldTune ? terms.comparabilityLabels.yes : terms.comparabilityLabels.no}
                </small>
            )}
            <small className="professional-tune-instruction">
                {terms.professionalTuneInstruction}: {instruction}
            </small>
            {adjustments.length > 0 ? (
                <div className="professional-adjustments">
                    <span>{terms.professionalChangeAmount}</span>
                    {adjustments.map(adjustment => (
                        <small key={adjustment.fieldKey}>
                            {formatProfessionalAdjustment(adjustment, language)}
                        </small>
                    ))}
                </div>
            ) : (
                <small className="muted-text">
                    {terms.professionalMergedNoAdjustments}
                    {missingProfessionalAdjustmentFields(advice, profile).length > 0
                        ? `: ${missingProfessionalAdjustmentFields(advice, profile).map(key => profileFieldDisplayName(key, language)).join(' / ')}`
                        : ''}
                </small>
            )}
            <EvidenceInline evidence={problem?.evidence || decision?.evidence || advice[0]?.evidence} terms={terms} language={language}/>
        </div>
    );
}

function buildProfessionalDiagnosticRows(pipeline: TuningPipelineRunResult): ProfessionalDiagnosticRow[] {
    const problems = pipeline.problemSet?.problems || [];
    const decisions = pipeline.decisionSet?.decisions || [];
    const advice = pipeline.adviceSet?.advice || [];
    const rows = new Map<string, ProfessionalDiagnosticRow>();

    const ensure = (key: string) => {
        let row = rows.get(key);
        if (!row) {
            row = {key, problems: [], decisions: [], advice: [], count: 0, durationMs: 0, confidence: 'low', riskLevel: 'low'};
            rows.set(key, row);
        }
        return row;
    };

    const problemKeyById = new Map<string, string>();
    for (const problem of problems) {
        const key = professionalProblemMergeKey(problem);
        problemKeyById.set(problem.id, key);
        if (problem.sourceId) {
            problemKeyById.set(problem.sourceId, key);
        }
        const row = ensure(key);
        row.problems.push(problem);
        row.count += Math.max(1, Number(problem.count || 0));
        row.durationMs += Math.max(0, Number(problem.durationMs || 0));
        row.confidence = higherProfessionalConfidence(row.confidence, problem.confidence);
        row.riskLevel = higherProfessionalRisk(row.riskLevel, problem.riskLevel || problem.severity);
    }

    for (const decision of decisions) {
        const key = problemKeyById.get(decision.problemId) || `decision:${decision.problemId || decision.id}`;
        const row = ensure(key);
        row.decisions.push(decision);
        row.confidence = higherProfessionalConfidence(row.confidence, decision.confidence);
    }

    for (const item of advice) {
        const key = problemKeyById.get(item.problemId) || `decision:${item.problemId || item.decisionId || item.id}`;
        const row = ensure(key);
        if (!row.advice.some(existing => existing.id === item.id)) {
            row.advice.push(item);
        }
        row.confidence = higherProfessionalConfidence(row.confidence, item.trustLevel);
    }

    return Array.from(rows.values())
        .map(row => ({
            ...row,
            count: row.count || row.problems.length || row.advice.length || row.decisions.length,
            durationMs: row.durationMs || 0,
        }))
        .sort((left, right) => (
            professionalRiskRank(right.riskLevel) - professionalRiskRank(left.riskLevel)
            || professionalConfidenceRank(right.confidence) - professionalConfidenceRank(left.confidence)
            || right.count - left.count
            || right.durationMs - left.durationMs
        ));
}

function professionalProblemMergeKey(problem: TuningProblem) {
    return [
        normalizeProfessionalProblemFamily(problem.family || problem.type || problem.id),
        normalizeProfessionalAxleGroup(problem.limitedAxle, problem.limitedWheels),
    ].join('|');
}

function normalizeProfessionalProblemFamily(value: string) {
    const key = String(value || '').toLowerCase();
    if (key.includes('lateral') || key.includes('understeer') || key.includes('oversteer')) {
        return 'lateral_balance';
    }
    if (key.includes('traction') || key.includes('power')) {
        return 'traction_limit';
    }
    if (key.includes('brake') || key.includes('decel')) {
        return 'braking_limit';
    }
    if (key.includes('platform') || key.includes('suspension') || key.includes('bottom')) {
        return 'platform_risk';
    }
    if (key.includes('thermal') || key.includes('temperature')) {
        return 'thermal_risk';
    }
    return key || 'unknown';
}

function normalizeProfessionalAxleGroup(axle: string | undefined, wheels: string[] | undefined) {
    const key = String(axle || '').toLowerCase();
    if (key === 'front' || key === 'rear' || key === 'driven' || key === 'all') {
        return key;
    }
    const wheelText = (wheels || []).join(' ').toLowerCase();
    if (wheelText.includes('fl') || wheelText.includes('fr') || wheelText.includes('front')) {
        return 'front';
    }
    if (wheelText.includes('rl') || wheelText.includes('rr') || wheelText.includes('rear')) {
        return 'rear';
    }
    return key || 'none';
}

function professionalRowPhaseSummary(row: ProfessionalDiagnosticRow, terms: Copy) {
    const phases = Array.from(new Set([
        ...row.problems.map(item => item.phase),
        ...row.decisions.map(item => item.phase),
    ].filter(Boolean)));
    if (phases.length === 0) {
        return pipelinePhaseLabel(undefined, terms);
    }
    const visible = phases.slice(0, 3).map(phase => pipelinePhaseLabel(phase, terms));
    const remaining = phases.length - visible.length;
    return remaining > 0 ? `${visible.join(' / ')} +${remaining}` : visible.join(' / ');
}

function professionalRowAxleSummary(row: ProfessionalDiagnosticRow, terms: Copy) {
    const axles = Array.from(new Set(row.problems.map(item => item.limitedAxle).filter(Boolean)));
    if (axles.length === 0) {
        return pipelineAxleLabel(undefined, terms);
    }
    const visible = axles.slice(0, 2).map(axle => pipelineAxleLabel(axle, terms));
    const remaining = axles.length - visible.length;
    return remaining > 0 ? `${visible.join(' / ')} +${remaining}` : visible.join(' / ');
}

function sortProfessionalAdvice(advice: TuningAdvice[], problem: TuningProblem | undefined, decision: TuningDecision | null, profile: TuneProfile | null) {
    return [...advice].sort((left, right) => {
        const leftConcrete = tuningAdviceConcreteAdjustments(left, problem, decision, profile, 'en').length;
        const rightConcrete = tuningAdviceConcreteAdjustments(right, problem, decision, profile, 'en').length;
        return rightConcrete - leftConcrete
            || professionalConfidenceRank(right.trustLevel) - professionalConfidenceRank(left.trustLevel)
            || Number(right.canApply) - Number(left.canApply);
    });
}

function bestProfessionalDecision(decisions: TuningDecision[]) {
    return [...decisions].sort((left, right) => (
        Number(right.shouldTune) - Number(left.shouldTune)
        || professionalConfidenceRank(right.confidence) - professionalConfidenceRank(left.confidence)
    ))[0] || null;
}

function professionalConfidenceRank(value: string | undefined) {
    const key = String(value || '').toLowerCase();
    if (key === 'high' || key === 'valid') return 3;
    if (key === 'medium') return 2;
    if (key === 'low' || key === 'low_confidence') return 1;
    return 0;
}

function professionalRiskRank(value: string | undefined) {
    const key = String(value || '').toLowerCase();
    if (key === 'critical' || key === 'severe' || key === 'high') return 3;
    if (key === 'medium') return 2;
    if (key === 'low') return 1;
    return 0;
}

function higherProfessionalConfidence(left: string, right: string | undefined) {
    return professionalConfidenceRank(right) > professionalConfidenceRank(left) ? String(right) : left;
}

function higherProfessionalRisk(left: string, right: string | undefined) {
    return professionalRiskRank(right) > professionalRiskRank(left) ? String(right) : left;
}

type ProfessionalConcreteAdjustment = {
    fieldKey: string;
    label: string;
    current: number;
    target: number;
    delta: number;
    step: number;
};

function tuningAdviceConcreteAdjustments(
    advice: TuningAdvice,
    problem: TuningProblem | undefined,
    decision: TuningDecision | null,
    profile: TuneProfile | null,
    language: Lang
): ProfessionalConcreteAdjustment[] {
    if (!profile) {
        return [];
    }
    const keys = professionalAdviceCandidateFields(advice, problem, profile).slice(0, 8);
    return keys.flatMap(key => {
        const field = profileFields.find(item => String(item.key) === key);
        if (!field || field.kind !== 'number') {
            return [];
        }
        const current = profile[field.key];
        if (typeof current !== 'number' || !Number.isFinite(current)) {
            return [];
        }
        const step = Number(field.step || '1') || 1;
        const delta = professionalAdviceDelta(key, advice, problem, decision, step);
        if (!delta) {
            return [];
        }
        const roundedDelta = roundToStep(delta, step);
        return [{
            fieldKey: key,
            label: profileFieldLabel(field, language),
            current,
            target: roundToStep(current + roundedDelta, step),
            delta: roundedDelta,
            step,
        }];
    });
}

function professionalAdviceCandidateFields(advice: TuningAdvice, problem: TuningProblem | undefined, profile: TuneProfile) {
    const declared = (advice.relatedFields || []).filter(key => profileFields.some(field => String(field.key) === key));
    const text = professionalAdviceText(advice, problem, null);
    const drivetrain = String(profile.drivetrain || '').trim().toUpperCase();
    const drivenAccel = drivetrain === 'FWD'
        ? ['frontDiffAccel']
        : drivetrain === 'RWD'
            ? ['rearDiffAccel']
            : ['rearDiffAccel', 'frontDiffAccel'];
    const drivenDecel = drivetrain === 'FWD'
        ? ['frontDiffDecel']
        : drivetrain === 'RWD'
            ? ['rearDiffDecel']
            : ['rearDiffDecel', 'frontDiffDecel'];
    const extras: string[] = [];
    if (text.includes('torque') || text.includes('drive_lock') || text.includes('power_to_tire') || text.includes('traction')) {
        extras.push(...drivenAccel, 'finalDrive', 'gear1', 'gear2');
    }
    if (text.includes('brake')) {
        extras.push('brakeBalance', 'brakePressure', ...drivenDecel);
    }
    if (problem?.limitedAxle === 'front' || text.includes('front') || text.includes('mechanical_grip') || text.includes('load_transfer')) {
        extras.push('frontArb', 'rearArb', 'frontCamber', 'frontTirePressure', 'frontRebound', 'frontBump');
    }
    if (problem?.limitedAxle === 'rear' || text.includes('rear') || text.includes('stability')) {
        extras.push('rearArb', 'frontArb', 'rearCamber', 'rearTirePressure', 'rearRebound', 'rearBump');
    }
    if (text.includes('platform') || text.includes('aero')) {
        extras.push('frontSpring', 'rearSpring', 'frontRebound', 'rearRebound', 'frontBump', 'rearBump');
    }
    return Array.from(new Set([...declared, ...extras]));
}

function professionalAdviceText(advice: TuningAdvice, problem: TuningProblem | undefined, decision: TuningDecision | null) {
    return [
        advice.category,
        advice.direction,
        advice.scope,
        advice.layer,
        advice.rationale,
        decision?.primaryCause,
        problem?.family,
        problem?.type,
        problem?.phase,
        problem?.limitedAxle,
    ].filter(Boolean).join(' ').toLowerCase();
}

function professionalAdviceDelta(
    fieldKey: string,
    advice: TuningAdvice,
    problem: TuningProblem | undefined,
    decision: TuningDecision | null,
    step: number
) {
    const text = professionalAdviceText(advice, problem, decision);
    const frontLimited = problem?.limitedAxle === 'front' || text.includes('front');
    const rearLimited = problem?.limitedAxle === 'rear' || text.includes('rear');

    if (text.includes('collect_more') || text.includes('observe') || text.includes('check_temperature')) {
        return 0;
    }
    if (text.includes('reduce_wheel_torque') || text.includes('drive_lock') || text.includes('traction')) {
        if (fieldKey.endsWith('DiffAccel')) return -2;
        if (fieldKey === 'finalDrive' || /^gear\d+$/.test(fieldKey)) return -Math.max(step, 0.03);
    }
    if (text.includes('brake')) {
        if (fieldKey === 'brakePressure') return -1;
        if (fieldKey === 'brakeBalance') {
            return frontLimited ? -1 : rearLimited ? 1 : 0;
        }
        if (fieldKey.endsWith('DiffDecel')) return -1;
    }
    if (text.includes('platform') || text.includes('aero')) {
        if (fieldKey.includes('Spring')) return 0.5;
        if (fieldKey.includes('Rebound') || fieldKey.includes('Bump')) return 0.1;
    }
    if (text.includes('mechanical_grip') || text.includes('load_transfer') || text.includes('rebalance')) {
        if (frontLimited) {
            if (fieldKey === 'frontArb') return -0.5;
            if (fieldKey === 'rearArb') return 0.5;
            if (fieldKey === 'frontCamber') return -0.1;
            if (fieldKey === 'frontTirePressure') return -0.01;
            if (fieldKey === 'frontRebound' || fieldKey === 'frontBump') return -0.1;
        }
        if (rearLimited) {
            if (fieldKey === 'rearArb') return -0.5;
            if (fieldKey === 'frontArb') return 0.5;
            if (fieldKey === 'rearCamber') return -0.1;
            if (fieldKey === 'rearTirePressure') return -0.01;
            if (fieldKey === 'rearRebound' || fieldKey === 'rearBump') return -0.1;
        }
    }
    return 0;
}

function uniqueProfessionalAdjustments(adjustments: ProfessionalConcreteAdjustment[]) {
    const seen = new Set<string>();
    return adjustments.filter(item => {
        if (seen.has(item.fieldKey)) {
            return false;
        }
        seen.add(item.fieldKey);
        return true;
    });
}

function missingProfessionalAdjustmentFields(advice: TuningAdvice[], profile: TuneProfile | null) {
    if (!profile) {
        return [];
    }
    return Array.from(new Set(advice.flatMap(item => item.relatedFields || []))).filter(key => {
        const field = profileFields.find(item => String(item.key) === key);
        return field?.kind === 'number' && typeof profile[field.key] !== 'number';
    }).slice(0, 4);
}

function formatProfessionalAdjustment(adjustment: ProfessionalConcreteAdjustment, language: Lang) {
    const verb = adjustment.delta < 0
        ? (language === 'zh' ? '降低' : 'decrease by')
        : (language === 'zh' ? '增加' : 'increase by');
    if (language === 'zh') {
        return `将 ${adjustment.label} 从 ${formatTuneNumber(adjustment.current, adjustment.step)} 调到 ${formatTuneNumber(adjustment.target, adjustment.step)}（${verb} ${formatTuneNumber(Math.abs(adjustment.delta), adjustment.step)}）`;
    }
    return `Set ${adjustment.label} from ${formatTuneNumber(adjustment.current, adjustment.step)} to ${formatTuneNumber(adjustment.target, adjustment.step)} (${verb} ${formatTuneNumber(Math.abs(adjustment.delta), adjustment.step)})`;
}

function TuneHarvestView({
    candidates,
    result,
    sources,
    dryRun,
    limit,
    statusFilter,
    search,
    terms,
    language,
    busy,
    running,
    stopping,
    onSourcesChange,
    onDryRunChange,
    onLimitChange,
    onStatusFilterChange,
    onSearchChange,
    onRun,
    onStop,
    onRefresh,
    onClear,
    onUse,
    onReject,
    onRestore,
}: {
    candidates: TuneHarvestCandidate[];
    result: TuneHarvestRunResult | null;
    sources: TuneHarvestSourceState;
    dryRun: boolean;
    limit: string;
    statusFilter: string;
    search: string;
    terms: Copy;
    language: Lang;
    busy: boolean;
    running: boolean;
    stopping: boolean;
    onSourcesChange: (sources: TuneHarvestSourceState) => void;
    onDryRunChange: (dryRun: boolean) => void;
    onLimitChange: (limit: string) => void;
    onStatusFilterChange: (status: string) => void;
    onSearchChange: (search: string) => void;
    onRun: () => void;
    onStop: () => void;
    onRefresh: () => void;
    onClear: () => void;
    onUse: (candidate: TuneHarvestCandidate) => void;
    onReject: (candidate: TuneHarvestCandidate) => void;
    onRestore: (candidate: TuneHarvestCandidate) => void;
}) {
    const statusCandidates = statusFilter === 'all'
        ? candidates
        : candidates.filter(candidate => candidate.status === statusFilter);
    const visibleCandidates = statusCandidates.filter(candidate => candidateMatchesTuneHarvestSearch(candidate, search));
    const warnings = (result?.warnings || []).filter(warning => warning !== 'harvest cancelled');
    const sourceKeys = Object.keys(sources) as Array<keyof TuneHarvestSourceState>;
    return (
        <section className="panel tune-harvest-panel">
            <div className="panel-heading">
                <div>
                    <h2>{terms.tuneHarvestTitle}</h2>
                    <span>{terms.tuneHarvestSubtitle}</span>
                </div>
                <div className="tune-harvest-heading-actions">
                    <button className="action primary" type="button" disabled={busy || running} onClick={onRun}>
                        <Radio size={16}/> {terms.tuneHarvestRun}
                    </button>
                    <button className="action secondary danger" type="button" disabled={!running || stopping} onClick={onStop}>
                        <Square size={16}/> {stopping ? terms.tuneHarvestStopping : terms.tuneHarvestStop}
                    </button>
                </div>
            </div>
            <div className="notice-card tune-harvest-controls">
                <strong>{terms.tuneHarvestSources}</strong>
                <div className="source-toggle-row">
                    {sourceKeys.map(source => (
                        <label key={source} className="source-toggle">
                            <input
                                type="checkbox"
                                checked={sources[source]}
                                onChange={event => onSourcesChange({...sources, [source]: event.target.checked})}
                            />
                            <span>{terms.tuneHarvestSourceLabels[source]}</span>
                        </label>
                    ))}
                </div>
                <div className="profile-field-grid compact">
                    <label className="profile-field">
                        <span className="profile-field-label">{terms.tuneHarvestLimit}</span>
                        <input type="number" min="1" max="500" value={limit} onChange={event => onLimitChange(event.target.value)}/>
                    </label>
                    <label className="profile-field">
                        <span className="profile-field-label">{terms.tuneHarvestStatusFilter}</span>
                        <select value={statusFilter} onChange={event => onStatusFilterChange(event.target.value)}>
                            {(['pending', 'rejected', 'imported', 'all'] as const).map(status => (
                                <option key={status} value={status}>{tuneHarvestStatusLabel(terms, status)}</option>
                            ))}
                        </select>
                    </label>
                    <label className="profile-field">
                        <span className="profile-field-label">{terms.tuneHarvestSearch}</span>
                        <input
                            type="search"
                            value={search}
                            placeholder={terms.tuneHarvestSearchPlaceholder}
                            onChange={event => onSearchChange(event.target.value)}
                        />
                    </label>
                    <label className="profile-field checkbox-field">
                        <span className="profile-field-label">{terms.tuneHarvestDryRun}</span>
                        <input type="checkbox" checked={dryRun} onChange={event => onDryRunChange(event.target.checked)}/>
                    </label>
                </div>
                <small>{terms.tuneHarvestDryRunHint}</small>
                {result && (
                    <small>{terms.tuneHarvestResult(result.found, result.saved, result.pending, result.rejected)}</small>
                )}
                {warnings.length > 0 && (
                    <details className="harvest-warnings">
                        <summary>{terms.tuneHarvestWarnings} ({warnings.length})</summary>
                        {warnings.slice(0, 8).map((warning, index) => <span key={`${warning}-${index}`}>{warning}</span>)}
                    </details>
                )}
            </div>
            <div className="panel-subsection">
                <div className="panel-heading compact">
                    <div>
                        <h2>{terms.tuneHarvestCandidates}</h2>
                        <span>{visibleCandidates.length} / {statusCandidates.length}</span>
                    </div>
                    <div className="tune-harvest-pool-actions">
                        <button className="small-action" type="button" disabled={busy || dryRun} onClick={onRefresh}>
                            {terms.tuneHarvestRefresh}
                        </button>
                        <button className="small-action danger" type="button" disabled={busy || candidates.length === 0} onClick={onClear}>
                            {terms.tuneHarvestClear}
                        </button>
                    </div>
                </div>
                {visibleCandidates.length === 0 ? (
                    <p className="empty-state">{terms.tuneHarvestNoCandidates}</p>
                ) : (
                    <div className="vehicle-reference-table tune-harvest-table">
                        <div className="vehicle-reference-row tune-harvest-row header">
                            <span>{terms.tuneHarvestSource}</span>
                            <span>{terms.tuneHarvestVehicle}</span>
                            <span>{terms.tuneHarvestCode}</span>
                            <span>PI</span>
                            <span>{terms.tuneHarvestContext}</span>
                            <span>{terms.tuneHarvestMatch}</span>
                            <span>{terms.tuneHarvestStatus}</span>
                            <span></span>
                        </div>
                        {visibleCandidates.map(candidate => (
                            <div className={`vehicle-reference-row tune-harvest-row ${candidate.status}`} key={`${candidate.id || candidate.rawKey}-${candidate.shareCode}`}>
                                <span>
                                    {candidate.sourceUrl ? (
                                        <a href={candidate.sourceUrl} target="_blank" rel="noreferrer">{terms.tuneHarvestSourceLabels[candidate.source as keyof TuneHarvestSourceState] || candidate.source}</a>
                                    ) : (
                                        terms.tuneHarvestSourceLabels[candidate.source as keyof TuneHarvestSourceState] || candidate.source
                                    )}
                                </span>
                                <strong title={candidate.matchedCarId || candidate.matchReason}>
                                    {candidate.carName || [candidate.year || '', candidate.make, candidate.model].filter(Boolean).join(' ') || '--'}
                                </strong>
                                <span>{formatTuneHarvestShareCode(candidate.shareCode)}</span>
                                <span>{candidate.carClass || '--'} {candidate.pi || ''} / {candidate.drivetrain || '--'} / {candidate.tireCompound || '--'}</span>
                                <span title={[candidate.tuneName, candidate.bestFor, candidate.notes].filter(Boolean).join(' / ')}>
                                    {[candidate.tuneName, candidate.bestFor, candidate.tuner].filter(Boolean).join(' / ') || '--'}
                                </span>
                                <span>{candidate.matchedCarId ? `${Math.round((candidate.matchScore || 0) * 100)}%` : candidate.matchReason || '--'}</span>
                                <span>{tuneHarvestStatusLabel(terms, candidate.status)}</span>
                                <span className="recommended-car-actions">
                                    <button className="small-action" type="button" disabled={busy || candidate.status === 'rejected'} onClick={() => onUse(candidate)}>
                                        {terms.tuneHarvestUseCandidate}
                                    </button>
                                    {candidate.status === 'rejected' ? (
                                        <button className="small-action" type="button" disabled={busy} onClick={() => onRestore(candidate)}>
                                            {terms.tuneHarvestRestore}
                                        </button>
                                    ) : candidate.status !== 'imported' && (
                                        <button className="small-action" type="button" disabled={busy} onClick={() => onReject(candidate)}>
                                            {terms.tuneHarvestReject}
                                        </button>
                                    )}
                                </span>
                                <small className="tune-harvest-updated">{formatUTC8DateTime(candidate.updatedAt, language)}</small>
                            </div>
                        ))}
                    </div>
                )}
            </div>
        </section>
    );
}

function RecommendedCarsGeneratorView({
    form,
    cars,
    result,
    fileSelection,
    version,
    terms,
    busy,
    formError,
    language,
    selectedIds,
    formOpen,
    isEditing,
    detailCar,
    onFormChange,
    onVersionChange,
    onSelectionChange,
    onOpenCreate,
    onRefresh,
    onCloseForm,
    onAdd,
    onClearForm,
    onEdit,
    onDuplicate,
    onShowDetail,
    onCloseDetail,
    onRemove,
    onDeleteAll,
    onSave,
}: {
    form: RecommendedCarForm;
    cars: RecommendedCar[];
    result: RecommendedCarsFileResult | null;
    fileSelection: RecommendedCarsFileSelection | null;
    version: string;
    terms: Copy;
    busy: boolean;
    formError: string;
    language: Lang;
    selectedIds: string[];
    formOpen: boolean;
    isEditing: boolean;
    detailCar: RecommendedCar | null;
    onFormChange: (field: keyof RecommendedCarForm, value: string) => void;
    onVersionChange: (value: string) => void;
    onSelectionChange: (ids: string[]) => void;
    onOpenCreate: () => void;
    onRefresh: () => void;
    onCloseForm: () => void;
    onAdd: () => void;
    onClearForm: () => void;
    onEdit: (car: RecommendedCar) => void;
    onDuplicate: (car: RecommendedCar) => void;
    onShowDetail: (car: RecommendedCar) => void;
    onCloseDetail: () => void;
    onRemove: (id: string) => void;
    onDeleteAll: () => void;
    onSave: () => void;
}) {
    const [search, setSearch] = useState('');
    const normalizedSearch = search.trim().toLowerCase();
    const filteredCars = useMemo(() => {
        if (!normalizedSearch) {
            return cars;
        }
        return cars.filter(car => [
            car.id,
            car.name,
            car.useCase,
            car.useCaseLabel,
            car.carClass,
            String(car.pi),
            car.drivetrain,
            car.tireCompound,
            car.tireCompoundLabel,
            car.tuneCode,
            car.tags?.join(' '),
            car.reason,
        ].some(value => String(value || '').toLowerCase().includes(normalizedSearch)));
    }, [cars, normalizedSearch]);
    const selectedSet = useMemo(() => new Set(selectedIds), [selectedIds]);
    const fileSelectedSet = useMemo(() => new Set(fileSelection?.ids || []), [fileSelection]);
    const fileSelectedTuneCodeSet = useMemo(() => new Set((fileSelection?.tuneCodes || []).map(normalizeTuneCodeInput).filter(Boolean)), [fileSelection]);
    const missingFileIDs = useMemo(() => {
        const knownIDs = new Set(cars.map(car => car.id));
        const knownTuneCodes = new Set(cars.map(car => normalizeTuneCodeInput(car.tuneCode)).filter(Boolean));
        const fileTuneCodes = (fileSelection?.tuneCodes || []).map(normalizeTuneCodeInput).filter(Boolean);
        if (fileTuneCodes.some(code => knownTuneCodes.has(code))) {
            return [];
        }
        return (fileSelection?.ids || []).filter(id => !knownIDs.has(id));
    }, [cars, fileSelection]);
    const isCarInFile = (car: RecommendedCar) => fileSelectedSet.has(car.id) || fileSelectedTuneCodeSet.has(normalizeTuneCodeInput(car.tuneCode));
    const visibleIds = filteredCars.map(car => car.id);
    const allVisibleSelected = visibleIds.length > 0 && visibleIds.every(id => selectedSet.has(id));
    const toggleSelected = (id: string) => {
        onSelectionChange(selectedSet.has(id) ? selectedIds.filter(item => item !== id) : [...selectedIds, id]);
    };
    const toggleVisible = () => {
        if (allVisibleSelected) {
            onSelectionChange(selectedIds.filter(id => !visibleIds.includes(id)));
            return;
        }
        onSelectionChange(Array.from(new Set([...selectedIds, ...visibleIds])));
    };
    return (
        <section className="panel">
            <div className="panel-heading">
                <div>
                    <h2>{terms.recommendedCarsTitle}</h2>
                    <span>{terms.recommendedCarsSubtitle}</span>
                </div>
                <button className="action primary" type="button" disabled={busy || selectedIds.length === 0} onClick={onSave}>
                    <Save size={16}/> {terms.recommendedCarsGenerate}
                </button>
            </div>
            <div className="notice-card">
                <strong>{terms.recommendedCarsTarget}</strong>
                <span>weChatApp/miniprogram/data/recommendedCars.json</span>
                <label className="profile-field">
                    <span className="profile-field-label">{terms.recommendedCarsVersion}</span>
                    <input value={version} onChange={event => onVersionChange(event.target.value)} placeholder="2026-05-28-001"/>
                </label>
                <small>{terms.recommendedCarsHint}</small>
                <div className="recommended-car-file-state">
                    <strong>{terms.recommendedCarsFileCurrent}</strong>
                    <span>
                        {fileSelection?.exists
                            ? terms.recommendedCarsFileFound(fileSelection.count, fileSelection.version)
                            : terms.recommendedCarsFileNotFound}
                    </span>
                    {fileSelection?.path && <small>{fileSelection.path}</small>}
                    {missingFileIDs.length > 0 && <small className="warning-text">{terms.recommendedCarsFileMissing(missingFileIDs.length)}</small>}
                </div>
                {result && <small>{terms.recommendedCarsSaved(result.count, result.path)}</small>}
            </div>
            {false && (
                <>
            {formError && (
                <div className="quick-tune-error-summary">
                    <strong>{formError}</strong>
                </div>
            )}
            <div className="profile-field-grid">
                <label className="profile-field">
                    <span className="profile-field-label">name</span>
                    <input value={form.name} onChange={event => onFormChange('name', event.target.value)} placeholder="2021 Porsche 911 GT3"/>
                </label>
                <label className="profile-field">
                    <span className="profile-field-label">useCase</span>
                    <select value={form.useCase} onChange={event => onFormChange('useCase', event.target.value)}>
                        {quickTuneUseCases.map(useCase => <option key={useCase} value={useCase}>{useCase}</option>)}
                    </select>
                </label>
                <label className="profile-field">
                    <span className="profile-field-label">useCaseLabel</span>
                    <input value={recommendedUseCaseLabels[form.useCase] || form.useCaseLabel} readOnly/>
                </label>
                <label className="profile-field">
                    <span className="profile-field-label">PI</span>
                    <input type="number" min="100" max="999" value={form.pi} onChange={event => onFormChange('pi', event.target.value)}/>
                </label>
                <label className="profile-field">
                    <span className="profile-field-label">carClass</span>
                    <select value={form.carClass} onChange={event => onFormChange('carClass', event.target.value)}>
                        {['D', 'C', 'B', 'A', 'S1', 'S2', 'R', 'X'].map(value => <option key={value} value={value}>{value}</option>)}
                    </select>
                </label>
                <label className="profile-field">
                    <span className="profile-field-label">drivetrain</span>
                    <select value={form.drivetrain} onChange={event => onFormChange('drivetrain', event.target.value)}>
                        {['FWD', 'AWD', 'RWD'].map(value => <option key={value} value={value}>{value}</option>)}
                    </select>
                </label>
                <label className="profile-field">
                    <span className="profile-field-label">tireCompound</span>
                    <select value={form.tireCompound} onChange={event => onFormChange('tireCompound', event.target.value)}>
                        {quickTuneTireCompounds.map(value => <option key={value} value={value}>{value}</option>)}
                    </select>
                </label>
                <label className="profile-field">
                    <span className="profile-field-label">tireCompoundLabel</span>
                    <input value={recommendedTireCompoundLabels[form.tireCompound] || form.tireCompoundLabel} readOnly/>
                </label>
                <label className="profile-field">
                    <span className="profile-field-label">weightKG</span>
                    <input type="number" min="0" value={form.weightKG} onChange={event => onFormChange('weightKG', event.target.value)}/>
                </label>
                <label className="profile-field">
                    <span className="profile-field-label">frontWeightPct</span>
                    <input type="number" min="0" max="99" value={form.frontWeightPct} onChange={event => onFormChange('frontWeightPct', event.target.value)}/>
                </label>
                <label className="profile-field">
                    <span className="profile-field-label">tuneCode</span>
                    <input value={form.tuneCode} onChange={event => onFormChange('tuneCode', event.target.value)} placeholder="413829605"/>
                </label>
                <label className="profile-field wide">
                    <span className="profile-field-label">imageSrc</span>
                    <input value={form.imageSrc} onChange={event => onFormChange('imageSrc', event.target.value)} placeholder="https://..."/>
                </label>
                <label className="profile-field">
                    <span className="profile-field-label">tags</span>
                    <input value={form.tags} onChange={event => onFormChange('tags', event.target.value)} placeholder={terms.recommendedCarsTagsHint}/>
                </label>
                <label className="profile-field wide">
                    <span className="profile-field-label">reason</span>
                    <textarea value={form.reason} onChange={event => onFormChange('reason', event.target.value)} placeholder="前驱容错高，适合熟悉轮胎定位和差速器反馈。"/>
                </label>
                <p className="baseline-tier-hint">{terms.recommendedCarsAutoId}</p>
                <p className="baseline-tier-hint">{terms.recommendedCarsOptionalMeta}</p>
            </div>
            <div className="form-actions">
                <button className="action primary" type="button" onClick={onAdd}>
                    <Plus size={16}/> {terms.recommendedCarsAdd}
                </button>
                <button className="action" type="button" onClick={onClearForm}>{terms.recommendedCarsClear}</button>
            </div>
                </>
            )}
            <div className="form-actions">
                <button className="action primary" type="button" onClick={onOpenCreate}>
                    <Plus size={16}/> {terms.recommendedCarsNew}
                </button>
            </div>
            {formOpen && (
                <RecommendedCarFormModal
                    form={form}
                    terms={terms}
                    busy={busy}
                    formError={formError}
                    isEditing={isEditing}
                    onFormChange={onFormChange}
                    onAdd={onAdd}
                    onClearForm={onClearForm}
                    onClose={onCloseForm}
                />
            )}
            {detailCar && (
                <RecommendedCarDetailModal
                    car={detailCar}
                    terms={terms}
                    language={language}
                    onClose={onCloseDetail}
                />
            )}
            <div className="panel-subsection">
                <div className="panel-heading compact">
                    <div>
                        <h2>{terms.recommendedCarsPending}</h2>
                        <span>{terms.recommendedCarsSelected(selectedIds.length, cars.length)}</span>
                    </div>
                </div>
                <div className="recommended-car-toolbar">
                    <label className="profile-field">
                        <span className="profile-field-label">{terms.recommendedCarsSearch}</span>
                        <input value={search} onChange={event => setSearch(event.target.value)} placeholder="Porsche / Road / 413829605"/>
                    </label>
                    <div className="recommended-car-toolbar-actions">
                        <button className="small-action" type="button" disabled={busy} onClick={onRefresh}>
                            {terms.recommendedCarsRefresh}
                        </button>
                        <button className="small-action" type="button" disabled={filteredCars.length === 0} onClick={toggleVisible}>
                            {terms.recommendedCarsSelectVisible}
                        </button>
                        <button className="small-action" type="button" disabled={selectedIds.length === 0} onClick={() => onSelectionChange([])}>
                            {terms.recommendedCarsClearSelection}
                        </button>
                        <button className="small-action danger" type="button" disabled={busy || cars.length === 0} onClick={onDeleteAll}>
                            {terms.recommendedCarsDeleteAll}
                        </button>
                    </div>
                </div>
                {cars.length === 0 ? (
                    <p className="empty-state">{terms.recommendedCarsNoItems}</p>
                ) : filteredCars.length === 0 ? (
                    <p className="empty-state">{terms.recommendedCarsNoSearchResults}</p>
                ) : (
                    <div className="vehicle-reference-table">
                        <div className="vehicle-reference-row recommended-car-row header">
                            <span>
                                <input type="checkbox" checked={allVisibleSelected} onChange={toggleVisible}/>
                            </span>
                            <span>ID</span>
                            <span>Name</span>
                            <span>Use Case</span>
                            <span>PI</span>
                            <span>Code</span>
                            <span>{terms.recommendedCarsImageSrc}</span>
                            <span>{terms.recommendedCarsCreatedAt}</span>
                            <span></span>
                        </div>
                        {filteredCars.map(car => (
                            <div className="vehicle-reference-row recommended-car-row" key={car.id}>
                                <span>
                                    <input type="checkbox" checked={selectedSet.has(car.id)} onChange={() => toggleSelected(car.id)}/>
                                </span>
                                <strong className="recommended-car-id-cell">
                                    <button className="recommended-car-id-button" type="button" onClick={() => onShowDetail(car)}>{car.id}</button>
                                    {isCarInFile(car) && <span className="recommended-car-file-chip">{terms.recommendedCarsInFile}</span>}
                                </strong>
                                <span>{car.name}</span>
                                <span>{car.useCaseLabel} / {car.useCase}</span>
                                <span>{car.carClass} {car.pi} / {car.drivetrain}</span>
                                <span>{car.tuneCode}</span>
                                <span title={car.imageSrc}>{car.imageSrc || '--'}</span>
                                <span>{formatUTC8DateTime(car.updatedAt, language)}</span>
                                <span className="recommended-car-actions">
                                    <button className="icon-button" type="button" onClick={() => onEdit(car)} aria-label={terms.recommendedCarsEdit}>
                                        <Pencil size={16}/>
                                    </button>
                                    <button className="icon-button" type="button" onClick={() => onDuplicate(car)} aria-label={terms.recommendedCarsDuplicate}>
                                        <CopyIcon size={16}/>
                                    </button>
                                    <button className="icon-button" type="button" onClick={() => onRemove(car.id)} aria-label={terms.recommendedCarsRemove}>
                                        <Trash2 size={16}/>
                                    </button>
                                </span>
                            </div>
                        ))}
                    </div>
                )}
            </div>
        </section>
    );
}

function RecommendedCarFormModal({
    form,
    terms,
    busy,
    formError,
    isEditing,
    onFormChange,
    onAdd,
    onClearForm,
    onClose,
}: {
    form: RecommendedCarForm;
    terms: Copy;
    busy: boolean;
    formError: string;
    isEditing: boolean;
    onFormChange: (field: keyof RecommendedCarForm, value: string) => void;
    onAdd: () => void;
    onClearForm: () => void;
    onClose: () => void;
}) {
    return (
        <div className="modal-backdrop" role="presentation">
            <section className="modal-card recommended-car-modal" role="dialog" aria-modal="true" aria-label={isEditing ? terms.recommendedCarsFormTitleEdit : terms.recommendedCarsFormTitleNew}>
                <div className="panel-heading">
                    <div>
                        <h2>{isEditing ? terms.recommendedCarsFormTitleEdit : terms.recommendedCarsFormTitleNew}</h2>
                        <span>{terms.recommendedCarsAutoId}</span>
                    </div>
                </div>
                {formError && (
                    <div className="quick-tune-error-summary">
                        <strong>{formError}</strong>
                    </div>
                )}
                <div className="profile-field-grid">
                    <label className="profile-field">
                        <span className="profile-field-label">name</span>
                        <input value={form.name} onChange={event => onFormChange('name', event.target.value)} placeholder="2021 Porsche 911 GT3"/>
                    </label>
                    <label className="profile-field">
                        <span className="profile-field-label">useCase</span>
                        <select value={form.useCase} onChange={event => onFormChange('useCase', event.target.value)}>
                            {quickTuneUseCases.map(useCase => <option key={useCase} value={useCase}>{useCase}</option>)}
                        </select>
                    </label>
                    <label className="profile-field">
                        <span className="profile-field-label">useCaseLabel</span>
                        <input value={recommendedUseCaseLabels[form.useCase] || form.useCaseLabel} readOnly/>
                    </label>
                    <label className="profile-field">
                        <span className="profile-field-label">PI</span>
                        <input type="number" min="100" max="999" value={form.pi} onChange={event => onFormChange('pi', event.target.value)}/>
                    </label>
                    <label className="profile-field">
                        <span className="profile-field-label">carClass</span>
                        <select value={form.carClass} onChange={event => onFormChange('carClass', event.target.value)}>
                            {['D', 'C', 'B', 'A', 'S1', 'S2', 'R', 'X'].map(value => <option key={value} value={value}>{value}</option>)}
                        </select>
                    </label>
                    <label className="profile-field">
                        <span className="profile-field-label">drivetrain</span>
                        <select value={form.drivetrain} onChange={event => onFormChange('drivetrain', event.target.value)}>
                            {['FWD', 'AWD', 'RWD'].map(value => <option key={value} value={value}>{value}</option>)}
                        </select>
                    </label>
                    <label className="profile-field">
                        <span className="profile-field-label">tireCompound</span>
                        <select value={form.tireCompound} onChange={event => onFormChange('tireCompound', event.target.value)}>
                            {quickTuneTireCompounds.map(value => <option key={value} value={value}>{value}</option>)}
                        </select>
                    </label>
                    <label className="profile-field">
                        <span className="profile-field-label">tireCompoundLabel</span>
                        <input value={recommendedTireCompoundLabels[form.tireCompound] || form.tireCompoundLabel} readOnly/>
                    </label>
                    <label className="profile-field">
                        <span className="profile-field-label">weightKG</span>
                        <input type="number" min="0" value={form.weightKG} onChange={event => onFormChange('weightKG', event.target.value)}/>
                    </label>
                    <label className="profile-field">
                        <span className="profile-field-label">frontWeightPct</span>
                        <input type="number" min="0" max="99" value={form.frontWeightPct} onChange={event => onFormChange('frontWeightPct', event.target.value)}/>
                    </label>
                    <label className="profile-field">
                        <span className="profile-field-label">tuneCode</span>
                        <input value={form.tuneCode} onChange={event => onFormChange('tuneCode', event.target.value)} placeholder="413829605"/>
                    </label>
                    <label className="profile-field wide">
                        <span className="profile-field-label">imageSrc</span>
                        <input value={form.imageSrc} onChange={event => onFormChange('imageSrc', event.target.value)} placeholder="https://..."/>
                    </label>
                    <label className="profile-field">
                        <span className="profile-field-label">tags</span>
                        <input value={form.tags} onChange={event => onFormChange('tags', event.target.value)} placeholder={terms.recommendedCarsTagsHint}/>
                    </label>
                    <label className="profile-field wide">
                        <span className="profile-field-label">reason</span>
                        <textarea value={form.reason} onChange={event => onFormChange('reason', event.target.value)}/>
                    </label>
                    <p className="baseline-tier-hint">{terms.recommendedCarsOptionalMeta}</p>
                </div>
                <div className="modal-actions">
                    <button className="action primary" type="button" disabled={busy} onClick={onAdd}>
                        <Save size={16}/> {terms.recommendedCarsAdd}
                    </button>
                    <button className="action" type="button" disabled={busy} onClick={onClearForm}>{terms.recommendedCarsClear}</button>
                    <button className="action" type="button" disabled={busy} onClick={onClose}>{terms.recommendedCarsCancel}</button>
                </div>
            </section>
        </div>
    );
}

function RecommendedCarDetailModal({car, terms, language, onClose}: { car: RecommendedCar; terms: Copy; language: Lang; onClose: () => void }) {
    const rows: Array<[string, string]> = [
        ['id', car.id],
        ['name', car.name],
        ['useCase', `${car.useCaseLabel} / ${car.useCase}`],
        ['carClass / pi', `${car.carClass} / ${car.pi}`],
        ['drivetrain', car.drivetrain],
        ['tireCompound', `${car.tireCompoundLabel} / ${car.tireCompound}`],
        ['weightKG', car.weightKG ? String(car.weightKG) : '--'],
        ['frontWeightPct', car.frontWeightPct ? `${car.frontWeightPct}%` : '--'],
        ['tuneCode', car.tuneCode],
        ['imageSrc', car.imageSrc || '--'],
        ['tags', (car.tags || []).join(', ') || '--'],
        ['reason', car.reason || '--'],
        ['createdAt', formatUTC8DateTime(car.createdAt, language)],
        ['updatedAt', formatUTC8DateTime(car.updatedAt, language)],
    ];
    return (
        <div className="modal-backdrop" role="presentation">
            <section className="modal-card recommended-car-modal" role="dialog" aria-modal="true" aria-label={terms.recommendedCarsDetailTitle}>
                <div className="panel-heading">
                    <div>
                        <h2>{terms.recommendedCarsDetailTitle}</h2>
                        <span>{car.name}</span>
                    </div>
                </div>
                <div className="recommended-car-detail-grid">
                    {rows.map(([label, value]) => (
                        <div key={label}>
                            <span>{label}</span>
                            <strong title={value}>{value}</strong>
                        </div>
                    ))}
                </div>
                <div className="modal-actions">
                    <button className="action primary" type="button" onClick={onClose}>{terms.close}</button>
                </div>
            </section>
        </div>
    );
}

function ProfessionalPipelineConfigPanel({
    catalog,
    config,
    terms,
    busy,
    onChange,
    onSave,
}: {
    catalog: TuningPipelineCatalog | null;
    config: ProfessionalPipelineConfig;
    terms: Copy;
    busy: boolean;
    onChange: (config: ProfessionalPipelineConfig) => void;
    onSave: () => void;
}) {
    return (
        <section className="panel">
            <div className="panel-heading">
                <div>
                    <h2>{terms.developerStrategyConfig}</h2>
                    <span>{terms.professionalDiagnosticHint}</span>
                </div>
                <button className="small-action" type="button" onClick={onSave} disabled={busy || !catalog}>{terms.saveAction}</button>
            </div>
            {!catalog ? (
                <div className="empty-events advice-placeholder">{terms.modelPipelineNoCatalog}</div>
            ) : (
                <div className="rule-form">
                    <label>
                        <span>{terms.modelPipelineDetector}</span>
                        <select value={config.detectorId} onChange={(event) => onChange({...config, detectorId: event.target.value})}>
                            {catalog.detectors.map(item => <option value={item.id} key={item.id}>{item.name}</option>)}
                        </select>
                    </label>
                    <label>
                        <span>{terms.modelPipelineDecisioner}</span>
                        <select value={config.decisionerId} onChange={(event) => onChange({...config, decisionerId: event.target.value})}>
                            {catalog.decisioners.map(item => <option value={item.id} key={item.id}>{item.name}</option>)}
                        </select>
                    </label>
                    <label>
                        <span>{terms.modelPipelineInterpreter}</span>
                        <select value={config.interpreterId} onChange={(event) => onChange({...config, interpreterId: event.target.value})}>
                            {catalog.interpreters.map(item => <option value={item.id} key={item.id}>{item.name}</option>)}
                        </select>
                    </label>
                </div>
            )}
        </section>
    );
}

function TelemetryControlStrip({
    interfaces,
    selectedAddress,
    port,
    status,
    busy,
    message,
    language,
    terms,
    targetAddress,
    onAddressChange,
    onPortChange,
    onStart,
    onStop,
}: {
    interfaces: NetworkInterface[];
    selectedAddress: string;
    port: number;
    status: TelemetryStatus;
    busy: boolean;
    message: string;
    language: Lang;
    terms: Copy;
    targetAddress: string;
    onAddressChange: (address: string) => void;
    onPortChange: (port: number) => void;
    onStart: () => void;
    onStop: () => void;
}) {
    return (
        <>
            <section className="control-strip">
                <label className="address-field">
                    <span>{terms.networkAdapter}</span>
                    <select
                        value={selectedAddress}
                        onChange={(event) => onAddressChange(event.target.value)}
                        disabled={status.running || busy}
                    >
                        {interfaces.map((item) => (
                            <option key={`${item.name}-${item.address}`} value={item.address}>
                                {interfaceLabel(item, language)}
                            </option>
                        ))}
                    </select>
                </label>
                <label className="port-field">
                    <span>{terms.udpPort}</span>
                    <input
                        value={port}
                        min={1}
                        max={65535}
                        type="number"
                        onChange={(event) => onPortChange(Number(event.target.value))}
                        disabled={status.running || busy}
                    />
                </label>
                <button className="action primary" onClick={onStart} disabled={status.running || busy}>
                    <Power size={18}/>
                    {terms.start}
                </button>
                <button className="action secondary" onClick={onStop} disabled={!status.running || busy}>
                    <Square size={17}/>
                    {terms.stop}
                </button>
                <div className="message-line">{message || status.lastError || terms.ready}</div>
            </section>

            <section className="target-card">
                <strong>{terms.gameTargetTitle}</strong>
                <span>{terms.gameTargetPrefix} <b>{targetAddress}</b>, {terms.gameTargetMiddle} <b>{port}</b>.</span>
            </section>
        </>
    );
}

function NewProfileModal({
    draft,
    current,
    terms,
    language,
    busy,
    onChange,
    onFillFromTelemetry,
    onCreate,
    onCancel,
}: {
    draft: TuneProfileInput;
    current: TelemetryFrame | null;
    terms: Copy;
    language: Lang;
    busy: boolean;
    onChange: (field: keyof TuneProfileInput, value: string | number | null) => void;
    onFillFromTelemetry: () => void;
    onCreate: () => void;
    onCancel: () => void;
}) {
    return (
        <div className="modal-backdrop" role="presentation">
            <section className="modal-card new-profile-modal" role="dialog" aria-modal="true" aria-label={terms.newProfileModalTitle}>
                <div className="panel-heading">
                    <div>
                        <h2>{terms.newProfileModalTitle}</h2>
                        <span>{terms.newProfileModalHint}</span>
                    </div>
                    <button className="small-action" type="button" onClick={onCancel} disabled={busy}>{terms.close}</button>
                </div>
                <div className="rule-form">
                    <label>
                        <span>{profileFieldDisplayName('carName', language)}</span>
                        <input value={draft.carName || ''} onChange={(event) => onChange('carName', event.target.value)}/>
                    </label>
                    <label>
                        <span>{profileFieldDisplayName('versionName', language)}</span>
                        <input value={draft.versionName || ''} onChange={(event) => onChange('versionName', event.target.value)}/>
                    </label>
                    <label>
                        <span>{profileFieldDisplayName('useCase', language)}</span>
                        <select value={draft.useCase || 'Road'} onChange={(event) => onChange('useCase', event.target.value)}>
                            {tuneUseCaseValues.map(useCase => (
                                <option value={useCase} key={useCase}>{terms.useCases[useCase]}</option>
                            ))}
                        </select>
                    </label>
                    <label>
                        <span>{terms.vehicleId}</span>
                        <input type="number" value={draft.carOrdinal ?? ''} onChange={(event) => onChange('carOrdinal', parseOptionalNumber(event.target.value))}/>
                    </label>
                    <label>
                        <span>{terms.classPi}</span>
                        <input value={draft.carClass || ''} onChange={(event) => onChange('carClass', event.target.value)}/>
                    </label>
                    <label>
                        <span>PI</span>
                        <input type="number" value={draft.pi ?? ''} onChange={(event) => onChange('pi', parseOptionalNumber(event.target.value))}/>
                    </label>
                    <label>
                        <span>{profileFieldDisplayName('drivetrain', language)}</span>
                        <input value={draft.drivetrain || ''} onChange={(event) => onChange('drivetrain', event.target.value)}/>
                    </label>
                    <label>
                        <span>{profileFieldDisplayName('numCylinders', language)}</span>
                        <input type="number" value={draft.numCylinders ?? ''} onChange={(event) => onChange('numCylinders', parseOptionalNumber(event.target.value))}/>
                    </label>
                </div>
                <div className="form-actions">
                    <button className="small-action" type="button" onClick={onFillFromTelemetry} disabled={busy || !current}>
                        <Gauge size={15}/>
                        {terms.fillTelemetryIntoDraft}
                    </button>
                    <button className="action primary" type="button" onClick={onCreate} disabled={busy}>
                        <Plus size={16}/>
                        {terms.createAndEdit}
                    </button>
                </div>
            </section>
        </div>
    );
}

function ProfileChoiceModal({
    pending,
    terms,
    busy,
    onChoose,
    onCancel,
}: {
    pending: PendingStartChoice;
    terms: Copy;
    busy: boolean;
    onChoose: (profile: TuneProfile) => void;
    onCancel: () => void;
}) {
    return (
        <div className="modal-backdrop" role="presentation">
            <section className="modal-card profile-choice-modal" role="dialog" aria-modal="true" aria-label={terms.profileChoiceTitle}>
                <div className="panel-heading">
                    <div>
                        <h2>{terms.profileChoiceTitle}</h2>
                        <span>{terms.vehicleId}: {pending.carOrdinal} / {terms.classPi}: {pending.carClass} / {formatOptionalInt(pending.carPi)}</span>
                    </div>
                    <button className="small-action" type="button" onClick={onCancel} disabled={busy}>{terms.close}</button>
                </div>
                <p className="modal-hint">{terms.profileChoiceHint}</p>
                <div className="profile-choice-list">
                    {pending.profiles.map(profile => (
                        <button className="profile-row" key={profile.id} type="button" onClick={() => onChoose(profile)} disabled={busy}>
                            <strong>{profile.carName}</strong>
                            <span>{[profile.versionName, profile.carClass, profile.drivetrain, localizedUseCase(profile.useCase, terms)].filter(Boolean).join(' / ') || '--'}</span>
                        </button>
                    ))}
                </div>
            </section>
        </div>
    );
}

function ProfileMismatchModal({
    pending,
    terms,
    busy,
    onChooseMatching,
    onClearAndStart,
    onCancel,
}: {
    pending: PendingProfileMismatch;
    terms: Copy;
    busy: boolean;
    onChooseMatching: () => void;
    onClearAndStart: () => void;
    onCancel: () => void;
}) {
    return (
        <div className="modal-backdrop" role="presentation">
            <section className="modal-card profile-choice-modal" role="dialog" aria-modal="true" aria-label={terms.profileMismatchTitle}>
                <div className="panel-heading">
                    <div>
                        <h2>{terms.profileMismatchTitle}</h2>
                        <span>{terms.profileMismatchHint}</span>
                    </div>
                    <button className="small-action" type="button" onClick={onCancel} disabled={busy}>{terms.close}</button>
                </div>
                <div className="profile-mismatch-body">
                    <div className="mismatch-card">
                        <strong>{terms.telemetryVehicle}</strong>
                        <span>{formatTelemetryVehicle(pending.telemetry, terms)}</span>
                    </div>
                    <div className="mismatch-card">
                        <strong>{terms.currentTuneVehicle}</strong>
                        <span>{formatTuneProfileVehicle(pending.profile, terms)}</span>
                    </div>
                </div>
                <div className="modal-actions">
                    <button
                        className="small-action"
                        type="button"
                        onClick={onChooseMatching}
                        disabled={busy || pending.candidates.length === 0}
                    >
                        {terms.chooseMatchingProfile}
                    </button>
                    <button className="small-action" type="button" onClick={onClearAndStart} disabled={busy}>
                        {terms.clearProfileAndStart}
                    </button>
                    <button className="small-action" type="button" onClick={onCancel} disabled={busy}>
                        {terms.close}
                    </button>
                </div>
            </section>
        </div>
    );
}

function SessionTuneProfileModal({
    pending,
    terms,
    busy,
    onChoose,
    onCancel,
}: {
    pending: PendingSessionBind;
    terms: Copy;
    busy: boolean;
    onChoose: (profile: TuneProfile) => void;
    onCancel: () => void;
}) {
    return (
        <div className="modal-backdrop" role="presentation">
            <section className="modal-card profile-choice-modal" role="dialog" aria-modal="true" aria-label={terms.sessionProfileBindTitle}>
                <div className="panel-heading">
                    <div>
                        <h2>{terms.sessionProfileBindTitle}</h2>
                        <span>{terms.sessionProfileBindHint}</span>
                    </div>
                    <button className="small-action" type="button" onClick={onCancel} disabled={busy}>{terms.close}</button>
                </div>
                <div className="profile-mismatch-body">
                    <div className="mismatch-card">
                        <strong>{terms.sessionVehicle}</strong>
                        <span>{formatSessionVehicle(pending.session, terms)}</span>
                    </div>
                    <div className="mismatch-card">
                        <strong>{terms.currentTuneVehicle}</strong>
                        <span>{pending.session.tuneName || terms.noProfile}</span>
                    </div>
                </div>
                <div className="panel-heading compact">
                    <h2>{terms.matchingTuneProfiles}</h2>
                    <span>{pending.profiles.length}</span>
                </div>
                {pending.profiles.length === 0 ? (
                    <div className="empty-events">{terms.noSessionProfileMatches}</div>
                ) : (
                    <div className="profile-choice-list">
                        {pending.profiles.map(profile => (
                            <button className="profile-row" key={profile.id} type="button" onClick={() => onChoose(profile)} disabled={busy}>
                                <strong>{profile.carName}</strong>
                                <span>{formatTuneProfileVehicle(profile, terms)}</span>
                                {pending.session.tuneProfileId === profile.id && <small>{terms.active}</small>}
                            </button>
                        ))}
                    </div>
                )}
            </section>
        </div>
    );
}

function TrackMergeModal({
    pending,
    terms,
    busy,
    onMerge,
    onSaveNew,
    onCancel,
}: {
    pending: PendingTrackMerge;
    terms: Copy;
    busy: boolean;
    onMerge: (candidate: TrackMergeCandidate) => void;
    onSaveNew: () => void;
    onCancel: () => void;
}) {
    return (
        <div className="modal-backdrop" role="presentation">
            <section className="modal-card track-merge-modal" role="dialog" aria-modal="true" aria-label={terms.similarTrackFound}>
                <div className="panel-heading">
                    <div>
                        <h2>{terms.similarTrackFound}</h2>
                        <span>{terms.similarTrackHint}</span>
                    </div>
                    <button className="small-action" type="button" onClick={onCancel} disabled={busy}>{terms.close}</button>
                </div>
                <div className="profile-choice-list">
                    {pending.candidates.map(candidate => (
                        <button className="profile-row" key={candidate.track.id} type="button" onClick={() => onMerge(candidate)} disabled={busy}>
                            <strong>{candidate.track.name}</strong>
                            <span>
                                {terms.matchLevel}: {trackMatchLevelLabel(candidate.matchLevel, terms)} / {terms.lengthError}: {formatNumber(candidate.lengthErrorPct, 1)}%
                            </span>
                            <small>
                                {terms.routeFitAvgError}: {formatMeters(candidate.routeFitAvgErrorMeters)} / {terms.routeFitP90Error}: {formatMeters(candidate.routeFitP90ErrorMeters)} / {terms.shapeSimilarity}: {formatPercentValue(candidate.routeFitScore)}
                            </small>
                        </button>
                    ))}
                </div>
                <div className="modal-actions">
                    <button className="action secondary" type="button" onClick={onSaveNew} disabled={busy}>{terms.saveAsNewTrack}</button>
                    <button className="small-action" type="button" onClick={onCancel} disabled={busy}>{terms.close}</button>
                </div>
            </section>
        </div>
    );
}

function TestConditionsPanel({
    terms,
    conditions,
    expanded,
    compact = false,
    disabled,
    onToggle,
    onChange,
    onReset,
}: {
    terms: Copy;
    conditions: TestConditions;
    expanded: boolean;
    compact?: boolean;
    disabled: boolean;
    onToggle: () => void;
    onChange: (field: keyof TestConditions, value: string) => void;
    onReset: () => void;
}) {
    const fields: Array<{ key: keyof TestConditions; label: string }> = [
        {key: 'brakeAssist', label: terms.brakeAssist},
        {key: 'steeringAssist', label: terms.steeringAssist},
        {key: 'tractionControl', label: terms.tractionControl},
        {key: 'stabilityControl', label: terms.stabilityControl},
        {key: 'shifting', label: terms.shifting},
        {key: 'launchControl', label: terms.launchControl},
    ];
    if (compact && !expanded) {
        return null;
    }
    return (
        <section className={`advanced-card${compact ? ' launchpad-settings' : ''}`}>
            <div className="advanced-card-heading">
                <button className="small-action" type="button" onClick={onToggle}>
                    {expanded ? terms.close : terms.testConditions}
                </button>
                <span>{terms.testConditions}: {formatTestConditionsCompact(conditions, terms)}</span>
            </div>
            {expanded && (
                <>
                    <div className="test-condition-grid">
                        {fields.map(field => (
                            <label key={field.key}>
                                <span>{field.label}</span>
                                <select
                                    value={conditions[field.key] || 'unknown'}
                                    onChange={(event) => onChange(field.key, event.target.value)}
                                    disabled={disabled}
                                >
                                    {testConditionOptions[field.key].map(option => (
                                        <option key={option} value={option}>{testConditionLabel(option, terms)}</option>
                                    ))}
                                </select>
                            </label>
                        ))}
                    </div>
                    <div className="advanced-actions">
                        <button className="small-action" type="button" onClick={onReset} disabled={disabled}>{terms.restoreUnknown}</button>
                    </div>
                </>
            )}
        </section>
    );
}

function ModelPipelineLabView({
    catalog,
    input,
    result,
    sessions,
    selectedSessionId,
    busy,
    terms,
    language,
    onInputChange,
    onRun,
}: {
    catalog: TuningPipelineCatalog | null;
    input: TuningPipelineRunInput;
    result: TuningPipelineRunResult | null;
    sessions: TelemetrySession[];
    selectedSessionId: number | null;
    busy: boolean;
    terms: Copy;
    language: Lang;
    onInputChange: (input: TuningPipelineRunInput) => void;
    onRun: () => void;
}) {
    const sourceTypes = catalog?.sourceTypes || [];
    const detectors = catalog?.detectors || [];
    const decisioners = catalog?.decisioners || [];
    const interpreters = catalog?.interpreters || [];
    const selectedSourceType = input.sourceType || 'tire_lab_current';
    const selectedSession = input.sessionId || selectedSessionId || sessions[0]?.id || 0;
    const update = (patch: Partial<TuningPipelineRunInput>) => onInputChange({...input, ...patch});

    return (
        <section className="panel">
            <div className="panel-heading">
                <div>
                    <h2>{terms.modelPipelineTitle}</h2>
                    <span>{terms.modelPipelineSubtitle}</span>
                </div>
                <button className="action" type="button" onClick={onRun} disabled={busy || !catalog}>
                    {terms.modelPipelineRun}
                </button>
            </div>
            {!catalog ? (
                <div className="empty-events">{terms.modelPipelineNoCatalog}</div>
            ) : (
                <>
                    <div className="generator-form compact">
                        <label>
                            <span>{terms.modelPipelineSource}</span>
                            <select
                                value={selectedSourceType}
                                onChange={(event) => {
                                    const sourceType = event.target.value;
                                    const defaults = catalog.defaultCombinations.find(item => item.sourceType === sourceType);
                                    update({
                                        sourceType,
                                        sessionId: sourceType === 'telemetry_session' ? selectedSession : 0,
                                        detectorId: defaults?.detectorId || input.detectorId,
                                        decisionerId: defaults?.decisionerId || input.decisionerId,
                                        interpreterId: defaults?.interpreterId || input.interpreterId,
                                    });
                                }}
                            >
                                {sourceTypes.map(component => (
                                    <option key={component.id} value={component.id}>{pipelineComponentName(component)}</option>
                                ))}
                            </select>
                        </label>
                        {selectedSourceType === 'telemetry_session' && (
                            <label>
                                <span>{terms.modelPipelineSession}</span>
                                <select value={selectedSession || 0} onChange={(event) => update({sessionId: Number(event.target.value)})}>
                                    {sessions.map(session => (
                                        <option key={session.id} value={session.id}>
                                            #{session.id} · {session.sessionName || session.tuneName || terms.noProfile}
                                        </option>
                                    ))}
                                </select>
                            </label>
                        )}
                        <label>
                            <span>{terms.modelPipelineDetector}</span>
                            <select value={input.detectorId} onChange={(event) => update({detectorId: event.target.value})}>
                                {detectors.map(component => (
                                    <option key={component.id} value={component.id}>{pipelineComponentName(component)}</option>
                                ))}
                            </select>
                        </label>
                        <label>
                            <span>{terms.modelPipelineDecisioner}</span>
                            <select value={input.decisionerId} onChange={(event) => update({decisionerId: event.target.value})}>
                                {decisioners.map(component => (
                                    <option key={component.id} value={component.id}>{pipelineComponentName(component)}</option>
                                ))}
                            </select>
                        </label>
                        <label>
                            <span>{terms.modelPipelineInterpreter}</span>
                            <select value={input.interpreterId} onChange={(event) => update({interpreterId: event.target.value})}>
                                {interpreters.map(component => (
                                    <option key={component.id} value={component.id}>{pipelineComponentName(component)}</option>
                                ))}
                            </select>
                        </label>
                    </div>
                    <div className="status-alert ok">{terms.modelPipelineExplainOnly}</div>
                    {!result ? (
                        <div className="empty-events">{terms.modelPipelineNoResult}</div>
                    ) : (
                        <div className="quick-diagnostic-stack">
                            <section className="quick-section">
                                <div className="panel-heading compact">
                                    <h2>{terms.modelPipelineSourceSummary}</h2>
                                    <span>{pipelineStatusLabel(result.status, terms)}</span>
                                </div>
                                <div className="launchpad-grid">
                                    <TextStat label={terms.modelPipelineSource} value={result.sourceSummary?.sourceType || '--'}/>
                                    <TextStat label={terms.modelPipelineSession} value={result.sourceSummary?.sessionId ? `#${result.sourceSummary.sessionId}` : '--'}/>
                                    <TextStat label={terms.samplesSaved} value={formatOptionalInt(result.sourceSummary?.sampleCount, true)}/>
                                    <TextStat label={terms.eventsSaved} value={formatOptionalInt(result.sourceSummary?.eventCount, true)}/>
                                    <TextStat label={terms.vehicleId} value={formatOptionalInt(result.sourceSummary?.vehicle?.carOrdinal ?? undefined)}/>
                                    <TextStat label={terms.classPi} value={`${result.sourceSummary?.vehicle?.carClass || '--'} / ${formatOptionalInt(result.sourceSummary?.vehicle?.carPi ?? undefined)}`}/>
                                </div>
                            </section>
                            <PipelineWarnings warnings={result.warnings} terms={terms}/>
                            <PipelineProblemsCard problemSet={result.problemSet} terms={terms} language={language}/>
                            <PipelineDecisionsCard decisionSet={result.decisionSet} terms={terms} language={language}/>
                            <PipelineAdviceCard adviceSet={result.adviceSet} terms={terms} language={language}/>
                        </div>
                    )}
                </>
            )}
        </section>
    );
}

function PipelineWarnings({warnings, terms}: { warnings: string[]; terms: Copy }) {
    if (!warnings || warnings.length === 0) {
        return null;
    }
    return (
        <section className="quick-section">
            <div className="panel-heading compact">
                <h2>{terms.modelPipelineWarnings}</h2>
                <span>{warnings.length}</span>
            </div>
            <div className="status-alerts">
                {warnings.map((warning, index) => (
                    <div className="status-alert warn" key={`${warning}-${index}`}>{warning}</div>
                ))}
            </div>
        </section>
    );
}

function pipelineStatusLabel(value: string | undefined, terms: Copy) {
    return localizedLabel(value || 'no_data', terms.modelPipelineStatusLabels);
}

function pipelineProblemLabel(problem: TuningProblem, terms: Copy) {
    return localizedLabel(problem.family || problem.type || problem.id, terms.tireIssueTypeLabels);
}

function pipelinePhaseLabel(value: string | undefined, terms: Copy) {
    return localizedLabel(value || 'unknown', terms.tireLabPhaseLabels);
}

function pipelineAxleLabel(value: string | undefined, terms: Copy) {
    return localizedLabel(value || 'none', terms.tireAxleLabels);
}

function pipelineConfidenceLabel(value: string | undefined, terms: Copy) {
    return localizedLabel(value || 'low', terms.quickConfidenceLabels);
}

function pipelineCauseLabel(value: string | undefined, terms: Copy) {
    return localizedLabel(value || 'unknown_tire_issue', terms.tireAdviceCauseLabels);
}

function pipelineAdviceCategoryLabel(value: string | undefined, terms: Copy) {
    const labels = terms.modelPipelineAdviceCategoryLabels as Record<string, string>;
    return labels[value || ''] || localizedLabel(value || 'data_quality', terms.tireAdviceCategoryLabels);
}

function pipelineAdviceDirectionLabel(value: string | undefined, terms: Copy) {
    const labels = terms.modelPipelineAdviceDirectionLabels as Record<string, string>;
    return labels[value || ''] || localizedLabel(value || 'continue_sampling', terms.tireAdviceDirectionLabels);
}

function pipelineAdviceScopeLabel(value: string | undefined, terms: Copy) {
    const labels = terms.modelPipelineScopeLabels as Record<string, string>;
    return labels[value || ''] || localizedLabel(value || 'none', terms.tireAxleLabels);
}

function pipelineRationaleLabel(value: string | undefined, terms: Copy) {
    if (!value) {
        return '--';
    }
    const pipelineLabels = terms.modelPipelineRationaleLabels as Record<string, string>;
    const tireLabels = terms.tireAdviceRationaleLabels as Record<string, string>;
    return pipelineLabels[value] || tireLabels[value] || formatEvidenceKey(value);
}

function pipelineEvidenceLabel(value: string, terms: Copy) {
    const labels = terms.modelPipelineEvidenceLabels as Record<string, string>;
    return labels[value] || formatEvidenceKey(value);
}

function pipelineRelatedFields(fields: string[], language: Lang) {
    return fields.map(field => profileFieldDisplayName(field, language)).join(' / ') || '--';
}

function formatProfessionalAdviceInstruction(
    advice: TuningAdvice | undefined,
    problem: TuningProblem | undefined,
    decision: TuningDecision | null,
    terms: Copy,
    language: Lang
) {
    if (!advice) {
        return decision ? pipelineCauseLabel(decision.primaryCause || decision.id, terms) : '--';
    }
    const narrative = professionalAdviceNarrative(advice, problem, decision, language);
    const fields = (advice.relatedFields || []).slice(0, 5);
    const fieldText = fields.length > 0 ? pipelineRelatedFields(fields, language) : '';
    return fieldText ? `${narrative} (${fieldText})` : narrative;
}

function professionalAdviceNarrative(advice: TuningAdvice, problem: TuningProblem | undefined, decision: TuningDecision | null, language: Lang) {
    const zh = language === 'zh';
    const text = professionalAdviceText(advice, problem, decision);
    const frontLimited = problem?.limitedAxle === 'front' || text.includes('front');
    const rearLimited = problem?.limitedAxle === 'rear' || text.includes('rear');

    if (text.includes('collect_more')) {
        return zh ? '继续采集数据，当前证据不足，暂不建议改车' : 'Collect more data first; evidence is not strong enough to tune yet';
    }
    if (text.includes('reduce_wheel_torque') || text.includes('drive_lock') || text.includes('traction')) {
        return zh ? '先降低轮上扭矩或驱动轮锁止，减少给油时轮胎打滑' : 'Reduce wheel torque or drive lock first to limit throttle-on tire slip';
    }
    if (text.includes('brake')) {
        if (frontLimited) {
            return zh ? '刹车负载偏向前轮，优先减小前轮抱死倾向' : 'Brake load is front-biased; reduce front lock tendency first';
        }
        if (rearLimited) {
            return zh ? '刹车时后轮稳定性不足，优先减小后轮抱死/甩尾倾向' : 'Rear stability is weak under braking; reduce rear lock or slide tendency first';
        }
        return zh ? '先检查制动力平衡和制动力压力，再看减速差速' : 'Check brake balance and pressure first, then decel differential';
    }
    if (text.includes('platform') || text.includes('aero')) {
        return zh ? '先提高高速/平台支撑，车高和空力按档位手动验证' : 'Improve high-speed platform support first; validate ride height and aero by tier';
    }
    if (text.includes('temperature') || text.includes('thermal') || text.includes('pressure_window')) {
        return zh ? '先确认胎压窗口和滑移升温，不直接按单帧胎温改车' : 'Verify pressure window and slip heat before tuning from tire temperature alone';
    }
    if (text.includes('mechanical_grip') || text.includes('load_transfer') || text.includes('rebalance')) {
        if (frontLimited) {
            return zh ? '增加前轴机械抓地，减少入弯/过弯推头' : 'Increase front mechanical grip to reduce entry or mid-corner understeer';
        }
        if (rearLimited) {
            return zh ? '增加后轴稳定性，减少过弯/出弯甩尾' : 'Increase rear stability to reduce cornering or exit oversteer';
        }
        return zh ? '重新平衡前后机械抓地和载荷转移' : 'Rebalance front/rear mechanical grip and load transfer';
    }
    return `${pipelineAdviceCategoryLabel(advice.category, COPY[language])} / ${pipelineAdviceDirectionLabel(advice.direction, COPY[language])}`;
}

function PipelineProblemsCard({problemSet, terms, language}: { problemSet: TuningProblemSet; terms: Copy; language: Lang }) {
    const problems = problemSet?.problems || [];
    return (
        <section className="quick-section">
            <div className="panel-heading compact">
                <h2>{terms.modelPipelineProblems}</h2>
                <span>{problemSet?.detectorId || '--'} / {terms.modelPipelineStatus}: {pipelineStatusLabel(problemSet?.status, terms)}</span>
            </div>
            {problems.length === 0 ? (
                <div className="empty-events">{terms.noEvents}</div>
            ) : (
                <div className="quick-suggestion-list">
                    {problems.map(problem => (
                        <div className="quick-suggestion-row" key={problem.id}>
                            <strong>{pipelineProblemLabel(problem, terms)}</strong>
                            <span>{pipelinePhaseLabel(problem.phase, terms)} / {pipelineAxleLabel(problem.limitedAxle, terms)} / {terms.modelPipelineConfidence}: {pipelineConfidenceLabel(problem.confidence, terms)}</span>
                            {(problem.operationTags || []).length > 0 && (
                                <small>{terms.tireIssueOperations}: {problem.operationTags.map(tag => localizedLabel(tag, terms.tireOperationTagLabels)).join(' / ')}</small>
                            )}
                            <small>{terms.eventsSaved}: {problem.count} / {formatDuration(problem.durationMs || 0, terms)}</small>
                            <EvidenceInline evidence={problem.evidence} terms={terms} language={language}/>
                        </div>
                    ))}
                </div>
            )}
        </section>
    );
}

function PipelineDecisionsCard({decisionSet, terms, language}: { decisionSet: TuningDecisionSet; terms: Copy; language: Lang }) {
    const decisions = decisionSet?.decisions || [];
    return (
        <section className="quick-section">
            <div className="panel-heading compact">
                <h2>{terms.modelPipelineDecisions}</h2>
                <span>{decisionSet?.decisionerId || '--'} / {terms.modelPipelineStatus}: {pipelineStatusLabel(decisionSet?.status, terms)}</span>
            </div>
            {decisions.length === 0 ? (
                <div className="empty-events">{terms.noEvents}</div>
            ) : (
                <div className="quick-suggestion-list">
                    {decisions.map(decision => (
                        <div className="quick-suggestion-row" key={decision.id}>
                            <strong>{pipelineCauseLabel(decision.primaryCause || decision.id, terms)}</strong>
                            <span>{pipelinePhaseLabel(decision.phase, terms)} / {terms.modelPipelineConfidence}: {pipelineConfidenceLabel(decision.confidence, terms)}</span>
                            <small>{terms.modelPipelineShouldTune}: {decision.shouldTune ? terms.comparabilityLabels.yes : terms.comparabilityLabels.no}</small>
                            <small>{pipelineRationaleLabel(decision.rationale || decision.documentContext, terms)}</small>
                            <EvidenceInline evidence={decision.evidence} terms={terms} language={language}/>
                        </div>
                    ))}
                </div>
            )}
        </section>
    );
}

function PipelineAdviceCard({adviceSet, terms, language}: { adviceSet: TuningAdviceSet; terms: Copy; language: Lang }) {
    const advice = adviceSet?.advice || [];
    return (
        <section className="quick-section">
            <div className="panel-heading compact">
                <h2>{terms.modelPipelineAdvice}</h2>
                <span>{adviceSet?.interpreterId || '--'} / {terms.modelPipelineStatus}: {pipelineStatusLabel(adviceSet?.status, terms)}</span>
            </div>
            {(adviceSet?.documentSources || []).length > 0 && (
                <div className="status-alert ok">{terms.modelPipelineDocs}: {adviceSet.documentSources.join(' / ')}</div>
            )}
            {advice.length === 0 ? (
                <div className="empty-events">{terms.noQuickSuggestions}</div>
            ) : (
                <div className="quick-suggestion-list">
                    {advice.map(item => (
                        <div className="quick-suggestion-row" key={item.id}>
                            <strong>{pipelineAdviceCategoryLabel(item.category, terms)} / {pipelineAdviceDirectionLabel(item.direction, terms)}</strong>
                            <span>{pipelineAdviceScopeLabel(item.scope, terms)} / {terms.modelPipelineConfidence}: {pipelineConfidenceLabel(item.trustLevel, terms)}</span>
                            <small>{pipelineRelatedFields(item.relatedFields || [], language)}</small>
                            <small>{pipelineRationaleLabel(item.rationale, terms)}</small>
                            {(item.verifyEvidence || []).length > 0 && (
                                <small>{terms.tireIssueVerifyEvidence}: {item.verifyEvidence.map(key => pipelineEvidenceLabel(key, terms)).join(' / ')}</small>
                            )}
                            {(item.documentSources || []).length > 0 && <small>{terms.modelPipelineDocs}: {item.documentSources.join(' / ')}</small>}
                            {item.conflictReason && <small className="warn-text">{pipelineRationaleLabel(item.conflictReason, terms)}</small>}
                            <EvidenceInline evidence={item.evidence} terms={terms} language={language}/>
                        </div>
                    ))}
                </div>
            )}
        </section>
    );
}

function EvidenceInline({evidence, terms}: { evidence?: Record<string, number>; terms: Copy; language: Lang }) {
    const entries = Object.entries(evidence || {}).slice(0, 4);
    if (entries.length === 0) {
        return null;
    }
    return (
        <small>
            {terms.modelPipelineEvidence}: {entries.map(([key, value]) => `${pipelineEvidenceLabel(key, terms)}=${formatEvidenceValue(value)}`).join(' / ')}
        </small>
    );
}

function pipelineComponentName(component: TuningPipelineComponent) {
    return component.name || component.id;
}

function TireRegressionLabView({
    samples,
    selectedId,
    selectedSample,
    results,
    saveForm,
    expectedForm,
    busy,
    terms,
    language,
    onSaveFormChange,
    onExpectedFormChange,
    onSelect,
    onSaveCurrent,
    onUpdateExpected,
    onRunSelected,
    onRunAll,
    onDelete,
}: {
    samples: TireRegressionSampleSummary[];
    selectedId: string;
    selectedSample: TireRegressionSample | null;
    results: TireRegressionResult[];
    saveForm: TireRegressionSaveFormState;
    expectedForm: TireRegressionExpectedFormState;
    busy: boolean;
    terms: Copy;
    language: Lang;
    onSaveFormChange: (form: TireRegressionSaveFormState) => void;
    onExpectedFormChange: (form: TireRegressionExpectedFormState) => void;
    onSelect: (id: string) => void;
    onSaveCurrent: () => void;
    onUpdateExpected: () => void;
    onRunSelected: () => void;
    onRunAll: () => void;
    onDelete: () => void;
}) {
    return (
        <section className="panel tire-regression-panel">
            <div className="panel-heading">
                <div>
                    <h2>{terms.tireRegressionTitle}</h2>
                    <span>{terms.tireRegressionSubtitle}</span>
                </div>
                <button className="primary-action" disabled={busy || samples.length === 0} onClick={onRunAll}>
                    <Activity size={16}/>{terms.tireRegressionRunAll}
                </button>
            </div>

            <div className="status-alerts">
                <div className="status-alert warn">{terms.tireRegressionRequiresTireLab}</div>
            </div>

            <div className="quick-section">
                <div className="section-title-row">
                    <div>
                        <h3>{terms.tireRegressionSaveCurrent}</h3>
                        <span>{terms.tireRegressionCsvHint}</span>
                    </div>
                    <button className="primary-action" disabled={busy} onClick={onSaveCurrent}>
                        <Save size={16}/>{terms.tireRegressionSave}
                    </button>
                </div>
                <div className="quick-comparability-grid">
                    <label>
                        <span>{terms.tireRegressionSampleName}</span>
                        <input value={saveForm.name} onChange={event => onSaveFormChange({...saveForm, name: event.target.value})}/>
                    </label>
                    <label>
                        <span>{terms.tireRegressionScenario}</span>
                        <input value={saveForm.scenario} onChange={event => onSaveFormChange({...saveForm, scenario: event.target.value})}/>
                    </label>
                    <label>
                        <span>{terms.tireRegressionWindowSeconds}</span>
                        <input type="number" min={5} max={30} value={saveForm.windowSeconds} onChange={event => onSaveFormChange({...saveForm, windowSeconds: event.target.value})}/>
                    </label>
                </div>
            </div>

            <div className="dashboard-grid">
                <div className="quick-section">
                    <div className="section-title-row">
                        <div>
                            <h3>{terms.tireRegressionSamples}</h3>
                            <span>{samples.length}</span>
                        </div>
                    </div>
                    {samples.length === 0 ? (
                        <div className="empty-events advice-placeholder">{terms.tireRegressionNoSamples}</div>
                    ) : (
                        <div className="rule-profile-list">
                            {samples.map(sample => (
                                <button
                                    className={`profile-row ${selectedId === sample.id ? 'active' : ''}`}
                                    key={sample.id}
                                    onClick={() => onSelect(sample.id)}
                                >
                                    <strong>{sample.name}</strong>
                                    <span>{sample.scenario || '--'}</span>
                                    <span>{vehicleSnapshotLabel(sample.vehicle, terms)}</span>
                                    <span>{sample.sampleCount} / {sample.windowSeconds}s</span>
                                </button>
                            ))}
                        </div>
                    )}
                </div>

                <div className="quick-section">
                    <div className="section-title-row">
                        <div>
                            <h3>{terms.tireRegressionExpected}</h3>
                            <span>{selectedSample ? selectedSample.name : terms.tireRegressionNoSelection}</span>
                        </div>
                        <div className="timeline-actions">
                            <button disabled={busy || !selectedSample} onClick={onRunSelected}><Activity size={15}/>{terms.tireRegressionRunOne}</button>
                            <button disabled={busy || !selectedSample} onClick={onUpdateExpected}><Save size={15}/>{terms.tireRegressionUpdateExpected}</button>
                            <button disabled={busy || !selectedSample} onClick={onDelete}><Trash2 size={15}/>{terms.tireRegressionDelete}</button>
                        </div>
                    </div>
                    {!selectedSample ? (
                        <div className="empty-events advice-placeholder">{terms.tireRegressionNoSelection}</div>
                    ) : (
                        <>
                            <div className="quick-comparability-grid">
                                <label>
                                    <span>{terms.tireRegressionAllowedPhases}</span>
                                    <input value={expectedForm.allowedPhases} onChange={event => onExpectedFormChange({...expectedForm, allowedPhases: event.target.value})}/>
                                </label>
                                <label>
                                    <span>{terms.tireRegressionRequiredGrip}</span>
                                    <input value={expectedForm.requiredGripTypes} onChange={event => onExpectedFormChange({...expectedForm, requiredGripTypes: event.target.value})}/>
                                </label>
                                <label>
                                    <span>{terms.tireRegressionAllowedAxles}</span>
                                    <input value={expectedForm.allowedAxles} onChange={event => onExpectedFormChange({...expectedForm, allowedAxles: event.target.value})}/>
                                </label>
                                <label>
                                    <span>{terms.tireRegressionForbiddenGrip}</span>
                                    <input value={expectedForm.forbiddenGripTypes} onChange={event => onExpectedFormChange({...expectedForm, forbiddenGripTypes: event.target.value})}/>
                                </label>
                                <label>
                                    <span>{terms.tireRegressionMinQuality}</span>
                                    <select value={expectedForm.minDataQuality} onChange={event => onExpectedFormChange({...expectedForm, minDataQuality: event.target.value})}>
                                        <option value="invalid">{localizedLabel('invalid', terms.tireDataQualityLabels)}</option>
                                        <option value="low_confidence">{localizedLabel('low_confidence', terms.tireDataQualityLabels)}</option>
                                        <option value="valid">{localizedLabel('valid', terms.tireDataQualityLabels)}</option>
                                    </select>
                                </label>
                            </div>
                            <label className="full-width-field">
                                <span>{terms.tireRegressionNotes}</span>
                                <textarea value={expectedForm.notes} onChange={event => onExpectedFormChange({...expectedForm, notes: event.target.value})}/>
                            </label>
                            <TireSnapshotSummary snapshot={selectedSample.snapshot} terms={terms} language={language}/>
                        </>
                    )}
                </div>
            </div>

            <div className="quick-section">
                <div className="section-title-row">
                    <div>
                        <h3>{terms.tireRegressionResults}</h3>
                        <span>{results.length}</span>
                    </div>
                </div>
                {results.length === 0 ? (
                    <div className="empty-events advice-placeholder">{terms.tireRegressionNoSelection}</div>
                ) : (
                    <div className="rule-profile-list">
                        {results.map(result => (
                            <TireRegressionResultRow result={result} terms={terms} language={language} key={result.sampleId}/>
                        ))}
                    </div>
                )}
            </div>
        </section>
    );
}

function TireRegressionResultRow({result, terms, language}: { result: TireRegressionResult; terms: Copy; language: Lang }) {
    const statusClass = result.passed ? 'ok' : 'warn';
    return (
        <div className={`status-alert ${statusClass}`}>
            <strong>{result.name} - {result.passed ? terms.tireRegressionPassed : terms.tireRegressionFailed}</strong>
            <span>{result.failures.length === 0 ? terms.tireRegressionPassed : result.failures.map(failure => localizedLabel(failure, terms.tireRegressionFailureLabels) || failure).join(' / ')}</span>
            <span>{terms.tireRegressionActual}: {snapshotShortLabel(result.actual, terms)}</span>
            <TireIssueAnalysisCard analysis={result.actual?.issueAnalysis} terms={terms} compact/>
            <TireIssueAdviceCard advice={result.actual?.issueAdvice} terms={terms} language={language} compact/>
        </div>
    );
}

function TireSnapshotSummary({snapshot, terms, language}: { snapshot: TireDiagnosticSnapshot; terms: Copy; language: Lang }) {
    if (!snapshot) {
        return null;
    }
    return (
        <div className="quick-section-subtle">
            <strong>{terms.tireRegressionActual}</strong>
            <div className="quick-comparability-grid">
                <TextStat label={terms.vehicleId} value={vehicleSnapshotLabel(snapshot.vehicle, terms)}/>
                <TextStat label={terms.tireDataQuality} value={localizedLabel(snapshot.dataQuality?.status || 'invalid', terms.tireDataQualityLabels)}/>
                <TextStat label={terms.tirePhaseCurrent} value={localizedLabel(snapshot.phase?.current || 'unknown', terms.tireLabPhaseLabels)}/>
                <TextStat label={terms.tireStablePhase} value={localizedLabel(snapshot.phase?.stable || 'unknown', terms.tireLabPhaseLabels)}/>
                <TextStat label={terms.tireGripLimit} value={localizedLabel(snapshot.gripLimit?.type || 'no_limit_detected', terms.tireGripLimitLabels)}/>
                <TextStat label={terms.tireLimitedAxle} value={localizedLabel(snapshot.gripLimit?.limitedAxle || 'none', terms.tireAxleLabels)}/>
                <TextStat label={terms.powerToTireStatus} value={localizedLabel(snapshot.power?.summary || '', terms.powerToTireSummaryLabels) || snapshot.power?.summary || '--'}/>
                <TextStat label={terms.brakeToTireStatus} value={localizedLabel(snapshot.brake?.summary || '', terms.brakeToTireSummaryLabels) || snapshot.brake?.summary || '--'}/>
            </div>
            <TireIssueAnalysisCard analysis={snapshot.issueAnalysis} terms={terms} compact/>
            <TireIssueAdviceCard advice={snapshot.issueAdvice} terms={terms} language={language} compact/>
        </div>
    );
}

function TireIssueAnalysisCard({analysis, terms, compact = false}: {
    analysis?: TireIssueAnalysis | null;
    terms: Copy;
    compact?: boolean;
}) {
    const groups = analysis?.groups || [];
    const segments = analysis?.segments || [];
    return (
        <div className="quick-section tire-issue-analysis-card">
            <div className="section-title-row">
                <div>
                    <h3>{terms.tireIssueGroups}</h3>
                    <span>{analysis ? `${analysis.groupCount || groups.length} / ${analysis.segmentCount || segments.length}` : terms.tireIssueNoGroups}</span>
                </div>
                <strong>{analysis?.updatedAt ? formatTime(analysis.updatedAt, terms.never) : '--'}</strong>
            </div>
            {groups.length === 0 ? (
                <div className="empty-events">{terms.tireIssueNoGroups}</div>
            ) : (
                <div className="quick-suggestion-list">
                    {groups.slice(0, compact ? 3 : 8).map(group => (
                        <div className="quick-suggestion-row tire-issue-row" key={group.id}>
                            <div>
                                <strong>{localizedLabel(group.type, terms.tireIssueTypeLabels)} · {localizedLabel(group.phase || 'unknown', terms.tireLabPhaseLabels)}</strong>
                                <span>
                                    {terms.tireIssueOperations}: {formatTireIssueTags(group.operationTags, terms)}
                                    {group.driftSource ? ` / ${localizedLabel(group.driftSource, terms.tireDriftSourceLabels)}` : ''}
                                </span>
                                <small>{formatTireIssueEvidence(group.representativeEvidence, terms)}</small>
                            </div>
                            <small>
                                {terms.tireIssueCount}: {group.count} · {terms.tireIssueDuration}: {formatDuration(group.totalDurationMs || 0, terms)}
                                <br/>
                                {terms.tireIssueSpeedRange}: {formatNumber(group.speedMinKmh, 0)}-{formatNumber(group.speedMaxKmh, 0)} km/h · {localizedLabel(group.limitedAxle || 'none', terms.tireAxleLabels)}
                            </small>
                        </div>
                    ))}
                </div>
            )}
            {!compact && (
                <details className="quick-section-subtle">
                    <summary>{terms.tireIssueSegments} ({segments.length})</summary>
                    {segments.length === 0 ? (
                        <div className="empty-events">{terms.tireIssueNoSegments}</div>
                    ) : (
                        <div className="quick-suggestion-list">
                            {segments.slice(0, 20).map(segment => (
                                <div className="quick-suggestion-row tire-issue-row" key={segment.id}>
                                    <div>
                                        <strong>{localizedLabel(segment.type, terms.tireIssueTypeLabels)} · {localizedLabel(segment.phase || 'unknown', terms.tireLabPhaseLabels)}</strong>
                                        <span>{formatTireIssueTags(segment.operationTags, terms)}</span>
                                    </div>
                                    <small>
                                        {formatDuration(segment.durationMs || 0, terms)} · {formatNumber(segment.speedMinKmh, 0)}-{formatNumber(segment.speedMaxKmh, 0)} km/h · {localizedLabel(segment.confidence || 'low', terms.quickConfidenceLabels)}
                                    </small>
                                </div>
                            ))}
                        </div>
                    )}
                </details>
            )}
        </div>
    );
}

function TireIssueAdviceCard({advice, terms, language, compact = false}: {
    advice?: TireIssueAdvice | null;
    terms: Copy;
    language: Lang;
    compact?: boolean;
}) {
    const priorityActions = advice?.priorityActions || [];
    const groups = advice?.groups || [];
    const visibleGroups = compact ? groups.slice(0, 2) : groups.slice(0, 8);
    return (
        <div className="quick-section tire-issue-advice-card">
            <div className="section-title-row">
                <div>
                    <h3>{terms.tireIssueAdvice}</h3>
                    <span>{terms.tireIssueExperimentHint}</span>
                </div>
                <strong>{advice?.confidence ? localizedLabel(advice.confidence, terms.quickConfidenceLabels) : '--'}</strong>
            </div>
            {!advice || (priorityActions.length === 0 && groups.length === 0) ? (
                <div className="empty-events">{terms.tireIssueNoAdvice}</div>
            ) : (
                <>
                    <div className="quick-suggestion-list">
                        {priorityActions.slice(0, 3).map(action => (
                            <TireIssueAdviceActionRow action={action} terms={terms} language={language} key={action.id}/>
                        ))}
                    </div>
                    {!compact && (
                        <details className="quick-section-subtle">
                            <summary>{terms.tireIssueGroupAdvice} ({groups.length})</summary>
                            {visibleGroups.length === 0 ? (
                                <div className="empty-events">{terms.tireIssueNoAdvice}</div>
                            ) : (
                                <div className="quick-suggestion-list">
                                    {visibleGroups.map(group => (
                                        <div className="quick-suggestion-row tire-issue-row" key={group.issueGroupId}>
                                            <div>
                                                <strong>
                                                    {localizedLabel(group.issueType, terms.tireIssueTypeLabels)} · {localizedLabel(group.phase || 'unknown', terms.tireLabPhaseLabels)}
                                                </strong>
                                                <span>{terms.tireIssuePrimaryCause}: {localizedLabel(group.primaryCause, terms.tireAdviceCauseLabels) || group.primaryCause}</span>
                                                <small>{terms.tireIssueOperations}: {formatTireIssueTags(group.operationTags, terms)}</small>
                                            </div>
                                            <small>
                                                {terms.tireIssueShouldTune}: {group.shouldTune ? terms.comparabilityLabels.yes : terms.comparabilityLabels.no}
                                                <br/>
                                                {terms.tireLimitedAxle}: {localizedLabel(group.limitedAxle || 'none', terms.tireAxleLabels)} · {localizedLabel(group.confidence || 'low', terms.quickConfidenceLabels)}
                                            </small>
                                            <div className="tire-advice-actions-inline">
                                                {(group.actions || []).map(action => (
                                                    <TireIssueAdviceActionRow action={action} terms={terms} language={language} key={action.id}/>
                                                ))}
                                            </div>
                                        </div>
                                    ))}
                                </div>
                            )}
                        </details>
                    )}
                </>
            )}
        </div>
    );
}

function TireIssueAdviceActionRow({action, terms, language}: { action: TireIssueAdviceAction; terms: Copy; language: Lang }) {
    const fields = (action.relatedFields || []).map(key => profileFieldDisplayName(key, language)).filter(Boolean);
    const verify = (action.verifyEvidence || []).map(formatEvidenceKey).join(' / ');
    return (
        <div className="quick-suggestion-row tire-advice-row">
            <div>
                <strong>
                    {localizedLabel(action.layer || 'primary', terms.tireAdviceLayerLabels)}
                    {' · '}
                    {localizedLabel(action.category || 'data_quality', terms.tireAdviceCategoryLabels)}
                    {' · '}
                    {localizedLabel(action.direction || 'observe', terms.tireAdviceDirectionLabels) || action.direction}
                </strong>
                <span>{localizedLabel(action.rationale || '', terms.tireAdviceRationaleLabels) || action.rationale || '--'}</span>
                {fields.length > 0 && <small>{terms.tireIssueRelatedFields}: {fields.join(' / ')}</small>}
                {verify && <small>{terms.tireIssueVerifyEvidence}: {verify}</small>}
                {(action.missingInputs || []).length > 0 && <small>{terms.tireIssueMissingInputs}: {action.missingInputs.join(' / ')}</small>}
                {action.conflictReason && <small>{terms.tireIssueConflict}: {localizedLabel(action.conflictReason, terms.tireAdviceRationaleLabels) || action.conflictReason}</small>}
            </div>
            <small>
                {localizedLabel(action.scope || 'none', terms.tireAxleLabels)}
                <br/>
                {localizedLabel(action.confidence || 'low', terms.quickConfidenceLabels)}
                <br/>
                {action.tuneRecommended ? terms.tireIssueShouldTune : terms.tireIssueNoTune}
            </small>
        </div>
    );
}

function formatTireIssueTags(tags: string[] = [], terms: Copy) {
    const labels = tags.map(tag => localizedLabel(tag, terms.tireOperationTagLabels)).filter(Boolean);
    return labels.length > 0 ? labels.join(' / ') : '--';
}

function formatTireIssueEvidence(evidence: Record<string, number> = {}, terms: Copy) {
    const preferred = [
        'front_combined_slip_p90',
        'rear_combined_slip_p90',
        'front_slip_ratio_p90',
        'rear_slip_ratio_p90',
        'front_slip_angle_p90',
        'rear_slip_angle_p90',
        'avg_speed_kmh',
        'avg_throttle',
        'avg_brake',
    ];
    const entries = preferred
        .filter(key => Number.isFinite(evidence[key]))
        .slice(0, 4)
        .map(key => `${formatEvidenceKey(key)} ${formatEvidenceValueForKey(key, evidence[key], terms)}`);
    return entries.join(' / ') || '--';
}

function TireModelLabView({diagnostic, status, current, terms, language, influenceMap}: {
    diagnostic: TireModelDiagnostic | null;
    status: TelemetryStatus;
    current: TelemetryFrame | null;
    terms: Copy;
    language: Lang;
    influenceMap: TuneToTireInfluenceMap | null;
}) {
    const hasData = Boolean(diagnostic && diagnostic.sampleCount > 0);
    const wheels = diagnostic?.wheels || [];
    return (
        <section className="panel tire-lab-panel">
            <div className="panel-heading">
                <div>
                    <h2>{terms.tireLabTitle}</h2>
                    <span>{terms.tireLabSubtitle}</span>
                </div>
                <span>{diagnostic?.updatedAt ? formatTime(diagnostic.updatedAt, terms.never) : terms.never}</span>
            </div>
            <div className="status-alerts">
                <div className="status-alert ok">{terms.tireLabNoPersistence}</div>
            </div>
            {!hasData ? (
                <div className="empty-events advice-placeholder">{terms.tireLabEmpty}</div>
            ) : (
                <>
                    <div className="launchpad-grid tire-lab-summary-grid">
                        <TextStat label={terms.tireLabLimitType} value={localizedLabel(diagnostic?.limitType || 'unknown', terms.tireLabLimitLabels)}/>
                        <TextStat label={terms.tirePhaseCurrent} value={localizedLabel(diagnostic?.phaseDetail?.currentPhase || diagnostic?.phase || 'unknown', terms.tireLabPhaseLabels)}/>
                        <TextStat label={terms.tirePhaseSecondary} value={localizedLabel(diagnostic?.phaseDetail?.secondaryPhase || 'unknown', terms.tireLabPhaseLabels)}/>
                        <TextStat label={terms.tireLabConfidence} value={localizedLabel(diagnostic?.phaseDetail?.confidence || diagnostic?.confidence || 'low', terms.quickConfidenceLabels)}/>
                        <TextStat label={terms.samplesSaved} value={`${diagnostic?.sampleCount || 0}`}/>
                        <TextStat label={terms.tireLabWindow} value={formatDuration(diagnostic?.windowMs || 0, terms)}/>
                        <TextStat label={terms.sessionMode} value={gameModeLabel(diagnostic?.gameMode || current?.gameMode || 'unknown', terms)}/>
                        <TextStat label={terms.vehicleId} value={formatOptionalInt(diagnostic?.vehicle?.carOrdinal || current?.carOrdinal)}/>
                        <TextStat label={terms.classPi} value={current ? `${current.carClass || '--'} / ${formatOptionalInt(current.carPi)}` : `${diagnostic?.vehicle?.carClass || '--'} / ${formatOptionalInt(diagnostic?.vehicle?.carPi || undefined)}`}/>
                    </div>

                    <div className="quick-section tire-lab-explanation">
                        <h3>{terms.tireLabExplanation}</h3>
                        <strong>{localizedLabel(diagnostic?.summary || '', terms.tireLabSummaryLabels)}</strong>
                        <span>{localizedLabel(diagnostic?.explanation || '', terms.tireLabExplanationLabels)}</span>
                    </div>

                    <TireDataQualityCard quality={diagnostic!.dataQuality} terms={terms}/>

                    <TireGripLimitCard gripLimit={diagnostic!.gripLimit} terms={terms}/>

                    <TireGripRelationshipsCard diagnostic={diagnostic!} terms={terms}/>

                    <TireModelGForceCard gForce={diagnostic!.gForce} terms={terms}/>

                    <TirePhaseEvidenceCard phase={diagnostic!.phaseDetail} terms={terms}/>

                    <TireIssueAnalysisCard analysis={diagnostic!.issueAnalysis} terms={terms}/>

                    <TireIssueAdviceCard advice={diagnostic!.issueAdvice} terms={terms} language={language}/>

                    <PowerToTireCard diagnostic={diagnostic!.powerToTire} terms={terms}/>

                    <BrakeToTireCard diagnostic={diagnostic!.brakeToTire} terms={terms}/>

                    <TireModelCamberCard camber={diagnostic!.camber} terms={terms}/>

                    <div className="quick-section">
                        <h3>{terms.tireLabMatrix}</h3>
                        <div className="tire-lab-wheel-grid">
                            {wheels.map(wheel => (
                                <TireModelWheelCard wheel={wheel} terms={terms} key={wheel.position}/>
                            ))}
                        </div>
                    </div>

                    <div className="quick-section">
                        <h3>{terms.tireLabAxleBalance}</h3>
                        <div className="quick-comparability-grid">
                            <TireAxleCard axle={diagnostic!.frontAxle} terms={terms}/>
                            <TireAxleCard axle={diagnostic!.rearAxle} terms={terms}/>
                            <TextStat
                                label={terms.tireLabLeftRight}
                                value={`${localizedLabel(diagnostic!.leftRight.state, terms.tireLabGripStateLabels)} / ${formatNumber(diagnostic!.leftRight.delta, 2)}`}
                            />
                        </div>
                    </div>

                    {(diagnostic?.warnings || []).length > 0 && (
                        <div className="quick-section">
                            <h3>{terms.tireLabWarnings}</h3>
                            <div className="status-alerts">
                                {(diagnostic?.warnings || []).map(warning => (
                                    <div className="status-alert warn" key={warning}>{localizedLabel(warning, terms.tireLabWarningLabels)}</div>
                                ))}
                            </div>
                        </div>
                    )}

                    <div className="quick-section">
                        <h3>{terms.tireLabHints}</h3>
                        {(diagnostic?.hints || []).length === 0 ? (
                            <div className="empty-events">{terms.tireLabEmpty}</div>
                        ) : (
                            <div className="quick-suggestion-list">
                                {(diagnostic?.hints || []).map((hint, index) => (
                                    <div className="quick-suggestion-row" key={`${hint.code}-${index}`}>
                                        <div>
                                            <strong>{localizedLabel(hint.code, terms.tireLabHintLabels)}</strong>
                                            <span>{localizedLabel(hint.severity, terms.tireLabGripStateLabels)}</span>
                                        </div>
                                        <small>{localizedLabel(hint.direction, terms.tireLabHintDirections) || hint.direction}</small>
                                    </div>
                                ))}
                            </div>
                        )}
                    </div>

                    <div className="connection-details tire-lab-connection">
                        <div>
                            <span>{terms.endpoint}</span>
                            <strong>{status.address}:{status.port}</strong>
                        </div>
                        <div>
                            <span>{terms.lastPacket}</span>
                            <strong>{formatTime(status.lastPacketAt, terms.never)}</strong>
                        </div>
                        <div>
                            <span>{terms.rawPackets}</span>
                            <strong>{status.rawPackets}</strong>
                        </div>
                        <div>
                            <span>{terms.validPackets}</span>
                            <strong>{status.validPackets}</strong>
                        </div>
                    </div>
                </>
            )}
            <TuneInfluenceMapSection influenceMap={influenceMap} terms={terms} language={language}/>
        </section>
    );
}

function TuneInfluenceMapSection({influenceMap, terms, language}: {
    influenceMap: TuneToTireInfluenceMap | null;
    terms: Copy;
    language: Lang;
}) {
    const grouped = useMemo(() => {
        const groups = new Map<string, TuneFieldInfluence[]>();
        for (const item of influenceMap?.items || []) {
            groups.set(item.category, [...(groups.get(item.category) || []), item]);
        }
        return Array.from(groups.entries());
    }, [influenceMap]);
    return (
        <details className="quick-section tune-influence-section">
            <summary>
                <span>{terms.tuneInfluenceMap}</span>
                <small>{terms.tuneInfluenceMapHint}</small>
            </summary>
            {grouped.length === 0 ? (
                <div className="empty-events">{terms.tuneInfluenceNoData}</div>
            ) : (
                <div className="tune-influence-map-grid">
                    {grouped.map(([category, items]) => (
                        <div className="tune-influence-category-card" key={category}>
                            <h4>{localizedLabel(category, terms.tuneInfluenceCategoryLabels) || category}</h4>
                            <div className="tune-influence-field-list">
                                {items.map(item => (
                                    <div className="tune-influence-map-row" key={item.fieldKey}>
                                        <div>
                                            <strong>{language === 'zh' ? item.labelZh || item.labelEn : item.labelEn || item.labelZh}</strong>
                                            <span>{language === 'zh' ? item.labelEn || item.labelZh : item.labelZh || item.labelEn}</span>
                                        </div>
                                        <TuneInfluenceChipGroup title={terms.tuneInfluenceType} values={[item.influenceType]} labels={terms.tuneInfluenceTypeLabels}/>
                                        <TuneInfluenceChipGroup title={terms.tuneInfluenceScope} values={item.scope} labels={terms.tuneInfluenceScopeLabels}/>
                                        <TuneInfluenceChipGroup title={terms.tuneInfluencePhase} values={item.phases} labels={terms.tuneInfluencePhaseLabels}/>
                                        <TuneInfluenceChipGroup title={terms.tuneInfluenceMetrics} values={item.tireMetrics} labels={terms.tuneInfluenceMetricLabels}/>
                                        <small>{language === 'zh' ? item.summaryZh || item.summaryEn : item.summaryEn || item.summaryZh}</small>
                                    </div>
                                ))}
                            </div>
                        </div>
                    ))}
                </div>
            )}
        </details>
    );
}

function TireDataQualityCard({quality, terms}: { quality?: TireDataQuality; terms: Copy }) {
    if (!quality) {
        return null;
    }
    const statusClass = quality.status === 'valid' ? 'ok' : quality.status === 'invalid' ? 'danger' : 'warn';
    return (
        <div className="quick-section tire-lab-quality-card">
            <div className="section-title-row">
                <div>
                    <h3>{terms.tireDataQuality}</h3>
                    <span>{(quality.reasons || []).slice(0, 2).map(reason => localizedLabel(reason, terms.tireDataQualityReasonLabels)).join(' / ') || terms.tireSignalQuality}</span>
                </div>
                <strong>{localizedLabel(quality.status || 'invalid', terms.tireDataQualityLabels)}</strong>
            </div>
            <div className="quick-comparability-grid">
                <TextStat label={terms.tireLabConfidence} value={localizedLabel(quality.confidence || 'low', terms.quickConfidenceLabels)}/>
                <TextStat label={terms.tireDynamicSamples} value={`${quality.dynamicSampleCount || 0} / ${quality.sampleCount || 0}`}/>
                <TextStat label={terms.speed} value={localizedLabel(quality.speedSignal || 'low', terms.tireDataQualityLabels)}/>
                <TextStat label={terms.gForceDiagnostics} value={localizedLabel(quality.gForceSignal || 'flat', terms.tireDataQualityLabels)}/>
                <TextStat label={terms.ratio} value={localizedLabel(quality.slipSignal || 'flat', terms.tireDataQualityLabels)}/>
                <TextStat label={terms.driverInputs} value={localizedLabel(quality.inputSignal || 'low', terms.tireDataQualityLabels)}/>
            </div>
            {(quality.reasons || []).length > 0 && (
                <div className="status-alerts">
                    {quality.reasons.map(reason => (
                        <div className={`status-alert ${statusClass}`} key={reason}>{localizedLabel(reason, terms.tireDataQualityReasonLabels) || reason}</div>
                    ))}
                </div>
            )}
        </div>
    );
}

function TireGripLimitCard({gripLimit, terms}: { gripLimit?: TireGripLimit; terms: Copy }) {
    if (!gripLimit) {
        return null;
    }
    const wheels = (gripLimit.limitedWheels || []).map(position => tireWheelPositionLabel(position, terms)).join(' / ') || '--';
    return (
        <div className="quick-section tire-lab-grip-limit-card">
            <div className="section-title-row">
                <div>
                    <h3>{terms.tireGripLimit}</h3>
                    <span>{localizedLabel(gripLimit.reason || '', terms.tireGripReasonLabels) || gripLimit.reason || '--'}</span>
                </div>
                <strong>{localizedLabel(gripLimit.type || 'no_limit_detected', terms.tireGripLimitLabels)}</strong>
            </div>
            <div className="quick-comparability-grid">
                <TextStat label={terms.tireLabConfidence} value={localizedLabel(gripLimit.confidence || 'low', terms.quickConfidenceLabels)}/>
                <TextStat label={terms.tireLimitedAxle} value={localizedLabel(gripLimit.limitedAxle || 'none', terms.tireAxleLabels)}/>
                <TextStat label={terms.tireLimitedWheels} value={wheels}/>
                <TextStat label={terms.tirePrimaryEvidence} value={formatEvidenceKey(gripLimit.primaryEvidence || '--')}/>
                <TextStat label={terms.tireLabAxleBalance} value={formatSignedNumber(gripLimit.frontRearDelta || 0, 2)}/>
                <TextStat label={terms.tireLabLeftRight} value={formatNumber(gripLimit.leftRightDelta, 2)}/>
            </div>
        </div>
    );
}

function TireGripRelationshipsCard({diagnostic, terms}: { diagnostic: TireModelDiagnostic; terms: Copy }) {
    return (
        <div className="quick-section tire-lab-relationships-card">
            <h3>{terms.tireGripRelationships}</h3>
            <div className="quick-comparability-grid">
                <TextStat label={terms.tireLabFrontAxle} value={`${formatNumber(diagnostic.frontAxle.combinedSlipP90, 2)} / ${formatNumber(diagnostic.frontAxle.slipAngleP90, 2)} / ${formatNumber(diagnostic.frontAxle.slipRatioP90, 2)}`}/>
                <TextStat label={terms.tireLabRearAxle} value={`${formatNumber(diagnostic.rearAxle.combinedSlipP90, 2)} / ${formatNumber(diagnostic.rearAxle.slipAngleP90, 2)} / ${formatNumber(diagnostic.rearAxle.slipRatioP90, 2)}`}/>
                <TextStat label={terms.tireLabLeftRight} value={`${localizedLabel(diagnostic.leftRight.state, terms.tireLabGripStateLabels)} / ${formatNumber(diagnostic.leftRight.delta, 2)}`}/>
                <TextStat label={terms.suspensionOffsetPct} value={`${formatRawPercent(Math.max(diagnostic.frontAxle.suspensionOffsetPctMax, diagnostic.rearAxle.suspensionOffsetPctMax))}`}/>
            </div>
        </div>
    );
}

function TireModelWheelCard({wheel, terms}: { wheel: TireWheelDiagnostic; terms: Copy }) {
    return (
        <div className={`wheel-card ${wheelGripCss(wheel.gripState)}`}>
            <div className="wheel-title">
                <span>{tireWheelPositionLabel(wheel.position, terms)}</span>
                <span>{localizedLabel(wheel.gripState, terms.tireLabGripStateLabels)}</span>
            </div>
            <div className="wheel-values tire-lab-wheel-values">
                <span>{terms.combined} <strong>{formatNumber(wheel.combinedSlipAvg, 2)} / {formatNumber(wheel.combinedSlipMax, 2)}</strong></span>
                <span>{terms.ratio} <strong>{formatNumber(wheel.slipRatioAvg, 2)} / {formatNumber(wheel.slipRatioMax, 2)}</strong></span>
                <span>{terms.angle} <strong>{formatNumber(wheel.slipAngleAvg, 2)} / {formatNumber(wheel.slipAngleMax, 2)}</strong></span>
                <span>{terms.temp} <strong>{formatNumber(wheel.tireTempAvg, 0)} / {formatNumber(wheel.tireTempMax, 0)}</strong></span>
                <span>{terms.suspensionOffsetPct} <strong>{formatRawPercent(wheel.suspensionOffsetPctAvg)} / {formatRawPercent(wheel.suspensionOffsetPctMax)}</strong></span>
            </div>
        </div>
    );
}

function TireModelGForceCard({gForce, terms}: { gForce: GForceDiagnostic; terms: Copy }) {
    return (
        <div className="quick-section tire-lab-gforce-card">
            <div className="section-title-row">
                <div>
                    <h3>{terms.gForceChart}</h3>
                    <span>{localizedLabel(gForce.axisMapping, terms.gForceAxisMappingLabels)}</span>
                </div>
                <strong>{terms.gForceDominantAxis}: {String(gForce.dominantAxis || '--').toUpperCase()}</strong>
            </div>
            <GForceCircle gForce={gForce} terms={terms}/>
            <div className="quick-comparability-grid">
                <TextStat label={`${terms.gForceCurrent} X/Y/Z`} value={`${formatGValue(gForce.currentXG)} / ${formatGValue(gForce.currentYG)} / ${formatGValue(gForce.currentZG)}`}/>
                <TextStat label={terms.gForceTotal} value={formatGValue(gForce.currentTotalG)}/>
                <TextStat label={`${terms.gForceAverage} X/Y/Z`} value={`${formatGValue(gForce.avgAbsXG)} / ${formatGValue(gForce.avgAbsYG)} / ${formatGValue(gForce.avgAbsZG)}`}/>
                <TextStat label={terms.gForcePeak} value={formatGValue(gForce.peakTotalG)}/>
                <TextStat label={`${terms.gForcePeak} X/Y/Z`} value={`${formatGValue(gForce.peakAbsXG)} / ${formatGValue(gForce.peakAbsYG)} / ${formatGValue(gForce.peakAbsZG)}`}/>
                <TextStat label={terms.gForceAxisMapping} value={gForce.source || '--'}/>
            </div>
        </div>
    );
}

function GForceCircle({gForce, terms}: { gForce: GForceDiagnostic; terms: Copy }) {
    const scale = 1.5;
    const x = clamp(-gForce.currentXG / scale, -1, 1);
    const z = clamp(-gForce.currentZG / scale, -1, 1);
    const pointX = 100 + x * 72;
    const pointY = 100 - z * 72;
    const planeG = Math.sqrt(gForce.currentXG * gForce.currentXG + gForce.currentZG * gForce.currentZG);
    return (
        <div className="gforce-circle-layout">
            <svg className="gforce-circle" viewBox="0 0 200 200" role="img" aria-label={terms.gForceChart}>
                <circle cx="100" cy="100" r="72" className="gforce-ring outer"/>
                <circle cx="100" cy="100" r="48" className="gforce-ring"/>
                <circle cx="100" cy="100" r="24" className="gforce-ring"/>
                <line x1="24" y1="100" x2="176" y2="100" className="gforce-axis"/>
                <line x1="100" y1="24" x2="100" y2="176" className="gforce-axis"/>
                <text x="178" y="96" className="gforce-axis-label">-X</text>
                <text x="12" y="96" className="gforce-axis-label">+X</text>
                <text x="104" y="30" className="gforce-axis-label">-Z</text>
                <text x="104" y="182" className="gforce-axis-label">+Z</text>
                <line x1="100" y1="100" x2={pointX} y2={pointY} className="gforce-vector"/>
                <circle cx={pointX} cy={pointY} r="5" className="gforce-dot"/>
                <circle cx="100" cy="100" r="3" className="gforce-center"/>
            </svg>
            <div className="gforce-circle-stats">
                <TextStat label={terms.gForcePlane} value={formatGValue(planeG)}/>
                <TextStat label={terms.gForceCircleScale} value={`${scale.toFixed(1)} g`}/>
                <TextStat label="X / Z" value={`${formatGValue(gForce.currentXG)} / ${formatGValue(gForce.currentZG)}`}/>
                <TextStat label="Y" value={formatGValue(gForce.currentYG)}/>
            </div>
        </div>
    );
}

function PowerToTireCard({diagnostic, terms}: { diagnostic?: PowerToTireDiagnostic; terms: Copy }) {
    if (!diagnostic) {
        return null;
    }
    const statusClass = diagnostic.status === 'ready' ? 'ok' : 'warn';
    const summary = localizedLabel(diagnostic.summary || 'power_to_tire_no_data', terms.powerToTireSummaryLabels);
    const explanation = localizedLabel(diagnostic.explanation || 'power_to_tire_waiting_for_samples', terms.powerToTireExplanationLabels);
    return (
        <div className="quick-section tire-lab-power-card">
            <div className="section-title-row">
                <div>
                    <h3>{terms.powerToTireTitle}</h3>
                    <span>{terms.powerToTireSubtitle}</span>
                </div>
                <strong>{localizedLabel(diagnostic.confidence || 'low', terms.quickConfidenceLabels)}</strong>
            </div>
            <div className="quick-comparability-grid">
                <TextStat label={terms.powerToTireStatus} value={summary}/>
                <TextStat label={terms.powerToTireDrivenAxle} value={localizedLabel(diagnostic.drivenAxle || 'unknown', terms.powerToTireDrivenAxleLabels)}/>
                <TextStat label={terms.powerToTirePower} value={`${formatNumber(diagnostic.currentPowerKW, 0)} / ${formatNumber(diagnostic.averagePowerKW, 0)} kW`}/>
                <TextStat label={terms.powerToTireTorque} value={`${formatNumber(diagnostic.currentTorqueNM, 0)} / ${formatNumber(diagnostic.averageTorqueNM, 0)} Nm`}/>
                <TextStat label={terms.powerToTireRPM} value={`${formatNumber(diagnostic.currentRPM, 0)} / ${formatPercentValue(diagnostic.averageRPMRatio)}`}/>
                <TextStat label={terms.powerToTireGear} value={diagnostic.currentGear > 0 ? String(diagnostic.currentGear) : '--'}/>
                <TextStat label={terms.powerToTireThrottle} value={formatPercentValue(diagnostic.averageThrottle)}/>
                <TextStat label={terms.powerToTireDrivenSlip} value={`${formatNumber(diagnostic.drivenSlipRatioP90, 2)} / ${formatPercentValue(diagnostic.drivenSlipRatioHighPct)}`}/>
                <TextStat label={terms.powerToTireAccel} value={`${formatGValue(diagnostic.averageAccelG)} / ${formatGValue(diagnostic.peakAccelG)}`}/>
                <TextStat label={terms.powerToTireSamples} value={`${diagnostic.highThrottleSampleCount} / ${diagnostic.sampleCount}`}/>
                <TextStat label={terms.powerToTireSignal} value={diagnostic.powerSignalAvailable ? terms.powerToTireAvailable : terms.powerToTireUnavailable}/>
            </div>
            <div className={`status-alert ${statusClass}`}>{explanation}</div>
        </div>
    );
}

function TirePhaseEvidenceCard({phase, terms}: { phase?: TirePhaseDiagnostic; terms: Copy }) {
    if (!phase) {
        return null;
    }
    const evidence = phase.evidence || {};
    const scores = Object.entries(phase.scores || {})
        .filter(([, score]) => Number.isFinite(score) && score > 0)
        .sort((left, right) => right[1] - left[1])
        .slice(0, 4);
    return (
        <div className="quick-section tire-lab-phase-card">
            <div className="section-title-row">
                <div>
                    <h3>{terms.tirePhaseEvidence}</h3>
                    <span>
                        {localizedLabel(phase.currentPhase || 'unknown', terms.tireLabPhaseLabels)}
                        {phase.secondaryPhase && phase.secondaryPhase !== 'unknown' ? ` / ${localizedLabel(phase.secondaryPhase, terms.tireLabPhaseLabels)}` : ''}
                    </span>
                </div>
                <strong>{localizedLabel(phase.confidence || 'low', terms.quickConfidenceLabels)}</strong>
            </div>
            <div className="quick-comparability-grid">
                <TextStat label={terms.tirePhaseCurrent} value={localizedLabel(phase.currentPhase || 'unknown', terms.tireLabPhaseLabels)}/>
                <TextStat label={terms.tirePhaseSecondary} value={localizedLabel(phase.secondaryPhase || 'unknown', terms.tireLabPhaseLabels)}/>
                <TextStat label={terms.tireStablePhase} value={localizedLabel(phase.stablePhase || 'unknown', terms.tireLabPhaseLabels)}/>
                <TextStat label={terms.tirePhaseStability} value={localizedLabel(phase.phaseStability || 'low_confidence', terms.tirePhaseStabilityLabels)}/>
                <TextStat label={terms.tireScoreMargin} value={formatPercentValue(phase.scoreMargin || 0)}/>
                <TextStat label={terms.tireLabWindow} value={`${formatDuration(phase.windowMs || 0, terms)} / ${phase.sampleCount || 0}`}/>
                <TextStat label={terms.tirePhaseSpeedDelta} value={`${formatNumber(evidence.avg_speed_kmh, 0)} km/h / ${formatSignedNumber(evidence.speed_delta_kmh || 0, 1)} km/h`}/>
                <TextStat label={terms.tirePhaseSpeedReference} value={formatTirePhaseSpeedReference(evidence.speed_reference_kmh, evidence.speed_band_confidence, terms)}/>
                <TextStat label={terms.tirePhaseSpeedBand} value={formatTirePhaseSpeedBand(evidence.speed_band, evidence.speed_band_confidence, terms)}/>
                <TextStat label={terms.powerToTireThrottle} value={formatPercentValue(evidence.avg_throttle)}/>
                <TextStat label={terms.tirePhaseThrottleDelta} value={`${formatSignedNumber((evidence.throttle_delta || 0) * 100, 0)}%`}/>
                <TextStat label={terms.tirePhaseBrake} value={`${formatPercentValue(evidence.avg_brake)} / ${formatPercentValue(evidence.peak_brake)}`}/>
                <TextStat label={terms.tirePhaseHandbrake} value={`${formatPercentValue(evidence.avg_handbrake)} / ${formatPercentValue(evidence.peak_handbrake)}`}/>
                <TextStat label={terms.brakeToTireSteer} value={formatPercentValue(evidence.avg_steer)}/>
                <TextStat label={terms.tirePhaseSteerDelta} value={`${formatSignedNumber((evidence.steer_delta || 0) * 100, 0)}%`}/>
                <TextStat label={terms.tirePhasePlaneG} value={`${formatGValue(evidence.avg_plane_g)} / ${formatGValue(evidence.peak_plane_g)}`}/>
                <TextStat label={terms.tirePhaseDecelG} value={`${formatGValue(evidence.avg_decel_g)} / ${formatGValue(evidence.peak_decel_g)}`}/>
                <TextStat label={terms.tirePhaseAccelG} value={`${formatGValue(evidence.avg_accel_g)} / ${formatGValue(evidence.peak_accel_g)}`}/>
            </div>
            {scores.length > 0 && (
                <div className="quick-suggestion-list">
                    <strong>{terms.tirePhaseScores}</strong>
                    {scores.map(([phaseName, score]) => (
                        <div className="quick-suggestion-row" key={phaseName}>
                            <div>
                                <strong>{localizedLabel(phaseName, terms.tireLabPhaseLabels)}</strong>
                            </div>
                            <small>{formatPercentValue(score)}</small>
                        </div>
                    ))}
                </div>
            )}
        </div>
    );
}

function BrakeToTireCard({diagnostic, terms}: { diagnostic?: BrakeToTireDiagnostic; terms: Copy }) {
    if (!diagnostic) {
        return null;
    }
    const statusClass = diagnostic.status === 'ready' ? 'ok' : 'warn';
    const summary = localizedLabel(diagnostic.summary || 'brake_to_tire_no_data', terms.brakeToTireSummaryLabels);
    const explanation = localizedLabel(diagnostic.explanation || 'brake_to_tire_waiting_for_samples', terms.brakeToTireExplanationLabels);
    return (
        <div className="quick-section tire-lab-brake-card">
            <div className="section-title-row">
                <div>
                    <h3>{terms.brakeToTireTitle}</h3>
                    <span>{terms.brakeToTireSubtitle}</span>
                </div>
                <strong>{localizedLabel(diagnostic.confidence || 'low', terms.quickConfidenceLabels)}</strong>
            </div>
            <div className="quick-comparability-grid">
                <TextStat label={terms.brakeToTireStatus} value={summary}/>
                <TextStat label={terms.brakeToTireBrake} value={`${formatPercentValue(diagnostic.averageBrake)} / ${formatPercentValue(diagnostic.peakBrake)}`}/>
                <TextStat label={terms.brakeToTireHandbrake} value={`${formatPercentValue(diagnostic.averageHandBrake)} / ${formatPercentValue(diagnostic.peakHandBrake)}`}/>
                <TextStat label={terms.brakeToTireSpeed} value={`${formatNumber(diagnostic.averageSpeedKmh, 0)} km/h / ${formatSignedNumber(diagnostic.speedDeltaKmh, 1)} km/h`}/>
                <TextStat label={terms.brakeToTireSteer} value={formatPercentValue(Math.abs(diagnostic.averageSteer))}/>
                <TextStat label={terms.brakeToTireDecel} value={`${formatGValue(diagnostic.averageDecelG)} / ${formatGValue(diagnostic.peakDecelG)}`}/>
                <TextStat label={terms.brakeToTirePlaneG} value={`${formatGValue(diagnostic.averagePlaneG)} / ${formatGValue(diagnostic.peakPlaneG)}`}/>
                <TextStat label={terms.brakeToTireFrontSlip} value={formatNumber(diagnostic.frontSlipRatioP90, 2)}/>
                <TextStat label={terms.brakeToTireRearSlip} value={formatNumber(diagnostic.rearSlipRatioP90, 2)}/>
                <TextStat label={terms.brakeToTireFrontCombined} value={formatNumber(diagnostic.frontCombinedSlipP90, 2)}/>
                <TextStat label={terms.brakeToTireRearCombined} value={formatNumber(diagnostic.rearCombinedSlipP90, 2)}/>
                <TextStat label={terms.brakeToTireSamples} value={`${diagnostic.brakeSampleCount} / ${diagnostic.sampleCount}`}/>
                <TextStat label={terms.brakeToTireTrail} value={diagnostic.trailBraking ? terms.detected : terms.notDetected}/>
                <TextStat label={terms.brakeToTireHandbrakeActive} value={diagnostic.handbrakeActive ? terms.detected : terms.notDetected}/>
            </div>
            <div className={`status-alert ${statusClass}`}>{explanation}</div>
        </div>
    );
}

function TireModelCamberCard({camber, terms}: { camber: CamberInference; terms: Copy }) {
    return (
        <div className="quick-section tire-lab-camber-card">
            <div className="section-title-row">
                <div>
                    <h3>{terms.camberInference}</h3>
                    <span>{localizedLabel(camber.explanation || '', terms.camberExplanationLabels)}</span>
                </div>
                <strong>{localizedLabel(camber.confidence || 'low', terms.quickConfidenceLabels)}</strong>
            </div>
            <div className="quick-comparability-grid">
                <TextStat label={terms.camberFront} value={localizedLabel(camber.frontState || 'unknown', terms.camberStateLabels)}/>
                <TextStat label={terms.camberRear} value={localizedLabel(camber.rearState || 'unknown', terms.camberStateLabels)}/>
                <TextStat label={terms.camberCorneringSamples} value={formatOptionalInt(camber.evidence?.cornering_samples, true)}/>
                <TextStat label={terms.tireLabExplanation} value={localizedLabel(camber.summary || '', terms.camberSummaryLabels)}/>
            </div>
            {(camber.warnings || []).length > 0 && (
                <div className="status-alerts">
                    {(camber.warnings || []).map(warning => (
                        <div className="status-alert warn" key={warning}>{localizedLabel(warning, terms.tireLabWarningLabels)}</div>
                    ))}
                </div>
            )}
            {(camber.hints || []).length > 0 && (
                <div className="quick-suggestion-list">
                    {(camber.hints || []).map((hint, index) => (
                        <div className="quick-suggestion-row" key={`${hint.code}-${index}`}>
                            <div>
                                <strong>{localizedLabel(hint.code, terms.tireLabHintLabels)}</strong>
                                <span>{localizedLabel(hint.severity, terms.tireLabGripStateLabels)}</span>
                            </div>
                            <small>{localizedLabel(hint.direction, terms.tireLabHintDirections) || hint.direction}</small>
                        </div>
                    ))}
                </div>
            )}
        </div>
    );
}

function TireAxleCard({axle, terms}: { axle: TireAxleDiagnostic; terms: Copy }) {
    return (
        <div className="quick-lap-card tire-axle-card">
            <strong>{axle.name === 'front' ? terms.tireLabFrontAxle : terms.tireLabRearAxle}</strong>
            <div className="mini-stat-grid">
                <TextStat label={terms.combined} value={`${formatNumber(axle.combinedSlipAvg, 2)} / ${formatNumber(axle.combinedSlipMax, 2)}`}/>
                <TextStat label={terms.ratio} value={formatNumber(axle.slipRatioAvg, 2)}/>
                <TextStat label={terms.angle} value={formatNumber(axle.slipAngleAvg, 2)}/>
                <TextStat label={terms.suspensionOffsetPct} value={`${formatRawPercent(axle.suspensionOffsetPctAvg)} / ${formatRawPercent(axle.suspensionOffsetPctMax)}`}/>
                <TextStat label={terms.tireLabLimitType} value={localizedLabel(axle.gripState, terms.tireLabGripStateLabels)}/>
            </div>
        </div>
    );
}

function QuickDiagnosticPanel({diagnostic, terms, language}: { diagnostic: QuickDiagnostic | null; terms: Copy; language: Lang }) {
    const groups = (diagnostic?.groups || []).slice(0, 5);
    const suggestions = (diagnostic?.suggestions || []).slice(0, 5);
    const missingFields = diagnostic?.missingProfileFields || [];
    const comparability = diagnostic?.comparability;
    const canShowLapComparison = diagnostic?.comparisonStatus === 'lap_comparison' && ['high', 'medium'].includes(comparability?.confidence || '');
    return (
        <section className="panel quick-diagnostic-panel">
            <div className="panel-heading">
                <div>
                    <h2>{terms.quickDiagnosticTitle}</h2>
                    <span>{diagnostic ? localizedLabel(diagnostic.comparisonStatus, terms.quickComparisonStatuses) : terms.quickDiagnosticEmpty}</span>
                </div>
                <span>{diagnostic?.updatedAt ? formatTime(diagnostic.updatedAt, terms.never) : terms.never}</span>
            </div>
            {!diagnostic || diagnostic.sampleCount === 0 ? (
                <div className="empty-events advice-placeholder">{terms.quickDiagnosticEmpty}</div>
            ) : (
                <>
                    <div className="launchpad-grid">
                        <Stat label={terms.samplesSaved} value={diagnostic.sampleCount}/>
                        <Stat label={terms.eventsSaved} value={diagnostic.eventCount}/>
                        <TextStat label={terms.sessionMode} value={gameModeLabel(diagnostic.gameMode, terms)}/>
                        <TextStat label={terms.driverMode} value={`${testConditionLabel(diagnostic.driverMode, terms)} / ${(diagnostic.driverModeConfidence * 100).toFixed(0)}%`}/>
                    </div>
                    {comparability && (
                        <>
                            <div className="quick-comparability-grid">
                                <TextStat label={terms.sameVehicleClass} value={localizedLabel(comparability.sameVehicleClass, terms.comparabilityLabels)}/>
                                <TextStat label={terms.sameTrackContext} value={localizedLabel(comparability.sameTrackContext, terms.comparabilityLabels)}/>
                                <TextStat label={terms.comparisonConfidence} value={localizedLabel(comparability.confidence, terms.quickConfidenceLabels)}/>
                            </div>
                            {(comparability.warnings || []).length > 0 && (
                                <div className="status-alerts quick-warning-list">
                                    {comparability.warnings.map(warning => (
                                        <div className="status-alert warn" key={warning}>{localizedLabel(warning, terms.quickWarningLabels)}</div>
                                    ))}
                                </div>
                            )}
                        </>
                    )}
                    {canShowLapComparison && diagnostic.currentLap && diagnostic.previousLap ? (
                        <div className="quick-lap-grid">
                            <QuickLapCard title={terms.currentLap} lap={diagnostic.currentLap} terms={terms}/>
                            <QuickLapCard title={terms.previousLap} lap={diagnostic.previousLap} terms={terms}/>
                        </div>
                    ) : (
                        <div className="empty-events advice-placeholder">{terms.quickRollingWindow}</div>
                    )}
                    <QuickGearPowerDiagnosticCard gearPower={diagnostic.gearPower} terms={terms}/>
                    <div className="quick-section">
                        <h3>{terms.issueGroups}</h3>
                        {groups.length === 0 ? (
                            <div className="empty-events">{terms.noIssueGroups}</div>
                        ) : (
                            <div className="quick-issue-list">
                                {groups.map(group => (
                                    <div className="quick-issue-row" key={group.id}>
                                        <div>
                                            <strong>{issueFamilyLabel(group.family, terms)}</strong>
                                            <span>{localizedLabel(group.severity, terms.planConfidenceLabels)} / {terms.eventsSaved}: {group.eventCount} / {formatDuration(group.totalDurationMs, terms)}</span>
                                        </div>
                                        {canShowLapComparison && <small>{terms.issueGroupComparison}: {issueComparisonLabel(group.comparison, terms)}</small>}
                                    </div>
                                ))}
                            </div>
                        )}
                    </div>
                    <div className="quick-section">
                        <h3>{terms.quickSuggestions}</h3>
                        {suggestions.length === 0 ? (
                            <div className="empty-events">{terms.noQuickSuggestions}</div>
                        ) : (
                            <div className="quick-suggestion-list">
                                {suggestions.map((suggestion, index) => (
                                    <div className="quick-suggestion-row" key={`${suggestion.item}-${suggestion.direction}-${index}`}>
                                        <div>
                                            <strong>{localizedLabel(suggestion.category, terms.actionCategories)} / {localizedLabel(suggestion.item, terms.actionItems)}</strong>
                                            <span>
                                                {localizedLabel(suggestion.direction, terms.actionDirections)}
                                                {' / '}
                                                {localizedLabel(suggestion.adviceLayer || 'primary', terms.adviceLayerLabels)}
                                                {' / '}
                                                {localizedLabel(suggestion.trustLevel || suggestion.confidence || 'medium', terms.tunePlanTrustLevels)}
                                            </span>
                                        </div>
                                        <small>{terms.quickSuggestionReason}: {localizedLabel(suggestion.rationale || suggestion.reason, terms.actionReasons) || suggestion.rationale || suggestion.reason}</small>
                                        {(suggestion.missingInputs || []).length > 0 && (
                                            <small className="warn-text">{terms.tunePlanMissingInputs}: {(suggestion.missingInputs || []).map(item => localizedLabel(item, terms.tunePlanMissingInputLabels) || item).join(' / ')}</small>
                                        )}
                                        {suggestion.nextStep && <small>{terms.quickSuggestionNextStep}: {localizedLabel(suggestion.nextStep, terms.quickNextStepLabels) || suggestion.nextStep}</small>}
                                    </div>
                                ))}
                            </div>
                        )}
                    </div>
                    {missingFields.length > 0 && (
                        <div className="quick-section">
                            <h3>{terms.missingProfileFields}</h3>
                            <div className="chip-list">
                                {missingFields.map(field => (
                                    <span key={field}>{profileFieldDisplayName(field, language)}</span>
                                ))}
                            </div>
                        </div>
                    )}
                </>
            )}
        </section>
    );
}

function RemoteTuneView({
    status,
    port,
    busy,
    terms,
    onPortChange,
    onStart,
    onStop,
}: {
    status: TuneWebServerStatus;
    port: string;
    busy: boolean;
    terms: Copy;
    onPortChange: (value: string) => void;
    onStart: () => void;
    onStop: () => void;
}) {
    return (
        <section className="tune-generator-page">
            <section className="panel launchpad-panel">
                <div className="panel-heading">
                    <div>
                        <h2>{terms.remoteTuneTitle}</h2>
                        <span>{terms.remoteTuneSubtitle}</span>
                    </div>
                    <span>{status.running ? terms.remoteTuneRunning : terms.remoteTuneStopped}</span>
                </div>

                <div className="status-alerts">
                    <div className="status-alert warn">{terms.remoteTuneReadOnly}</div>
                    <div className="status-alert ok">{terms.remoteTuneDeviceHint}</div>
                    {status.lastError && <div className="status-alert danger">{status.lastError}</div>}
                </div>

                <div className="baseline-form-grid">
                    <label className="profile-field">
                        <span>{terms.remoteTunePort}</span>
                        <input
                            type="number"
                            min={1}
                            max={65535}
                            step={1}
                            value={port}
                            disabled={status.running || busy}
                            onChange={event => onPortChange(event.target.value)}
                        />
                    </label>
                    <label className="profile-field">
                        <span>{terms.remoteTuneUrl}</span>
                        <input value={status.url || '--'} readOnly/>
                    </label>
                </div>

                <div className="launchpad-grid">
                    <TextStat label={terms.remoteTuneStatus} value={status.running ? terms.remoteTuneRunning : terms.remoteTuneStopped}/>
                    <TextStat label={terms.remoteTunePort} value={status.port ? String(status.port) : port || '--'}/>
                    <TextStat label={terms.remoteTuneLanAddress} value={status.lanAddress || '--'}/>
                    <TextStat label={terms.remoteTuneUrl} value={status.url || '--'}/>
                </div>

                <div className="form-actions">
                    {status.running ? (
                        <button className="action secondary" type="button" disabled={busy} onClick={onStop}>
                            <Square size={17}/>
                            {terms.remoteTuneStop}
                        </button>
                    ) : (
                        <button className="action primary" type="button" disabled={busy} onClick={onStart}>
                            <Power size={17}/>
                            {terms.remoteTuneStart}
                        </button>
                    )}
                </div>
            </section>
        </section>
    );
}

function TuneGeneratorView({
    form,
    result,
    profiles,
    selectedFields,
    targetProfileId,
    advancedOpen,
    fieldErrors,
    busy,
    terms,
    language,
    onFieldChange,
    onBiasChange,
    onFieldBlur,
    onAdvancedToggle,
    inputOpen,
    onInputOpen,
    onInputClose,
    onGenerate,
    onCreate,
    onApply,
    onReset,
    onToggleField,
    onTargetChange,
    onOpenExpert,
}: {
    form: RoadStaticTuneBaselineForm;
    result: RoadStaticTuneBaselineResult | null;
    profiles: TuneProfile[];
    selectedFields: string[];
    targetProfileId: number;
    advancedOpen: boolean;
    fieldErrors: QuickTuneFieldErrors;
    busy: boolean;
    terms: Copy;
    language: Lang;
    onFieldChange: (key: RoadStaticTuneBaselineFormKey, value: string) => void;
    onBiasChange: (key: 'balanceBias' | 'stiffnessBias' | 'speedBias', value: string) => void;
    onFieldBlur: (key: RoadStaticTuneBaselineFormKey, step: number) => void;
    onAdvancedToggle: () => void;
    inputOpen: boolean;
    onInputOpen: () => void;
    onInputClose: () => void;
    onGenerate: () => void;
    onCreate: () => void;
    onApply: () => void;
    onReset: () => void;
    onToggleField: (key: string) => void;
    onTargetChange: (id: number) => void;
    onOpenExpert: () => void;
}) {
    type QuickTunePreviewEntry =
        | { type: 'generated'; fieldKey: string; group: string; field: BaselineGeneratedField }
        | { type: 'tier'; fieldKey: string; group: string; field: BaselineTierRecommendation };
    const previewByGroup = ([
        ...(result?.generatedFields || []).map(field => ({
            type: 'generated' as const,
            fieldKey: String(field.fieldKey),
            group: String(field.group),
            field,
        })),
        ...(result?.tierRecommendations || []).map(field => ({
            type: 'tier' as const,
            fieldKey: String(field.fieldKey),
            group: String(field.group),
            field,
        })),
    ] as QuickTunePreviewEntry[]).reduce((acc, field) => {
        acc[field.group] = [...(acc[field.group] || []), field];
        return acc;
    }, {} as Record<string, QuickTunePreviewEntry[]>);
    Object.keys(previewByGroup).forEach(group => {
        previewByGroup[group] = [...previewByGroup[group]].sort((left, right) => profileFieldOrder(left.fieldKey) - profileFieldOrder(right.fieldKey));
    });
    const generatedByGroup = (result?.generatedFields || []).reduce((acc, field) => {
        acc[field.group] = [...(acc[field.group] || []), field];
        return acc;
    }, {} as Record<string, BaselineGeneratedField[]>);
    Object.keys(generatedByGroup).forEach(group => {
        generatedByGroup[group] = [...generatedByGroup[group]].sort((left, right) => profileFieldOrder(String(left.fieldKey)) - profileFieldOrder(String(right.fieldKey)));
    });
    const skippedByGroup = (result?.skippedFields || []).reduce((acc, field) => {
        acc[field.group] = [...(acc[field.group] || []), field];
        return acc;
    }, {} as Record<string, BaselineSkippedField[]>);
    Object.keys(skippedByGroup).forEach(group => {
        skippedByGroup[group] = [...skippedByGroup[group]].sort((left, right) => profileFieldOrder(String(left.fieldKey)) - profileFieldOrder(String(right.fieldKey)));
    });
    const generatorCoreGroupOrder = ['tire', 'gearing', 'alignment', 'antiroll', 'springs', 'damping', 'brake', 'aero', 'differential', 'power'] as Array<keyof Copy['fieldGroups']>;
    const groupOrder = [
        ...generatorCoreGroupOrder.filter(group => previewByGroup[group]?.length),
        ...Object.keys(previewByGroup)
            .filter(group => !generatorCoreGroupOrder.includes(group as keyof Copy['fieldGroups']) && previewByGroup[group]?.length)
            .sort(),
    ];
    const skippedGroupOrder = generatorCoreGroupOrder;
    const tiersByGroup = (result?.tierRecommendations || []).reduce((acc, field) => {
        acc[field.group] = [...(acc[field.group] || []), field];
        return acc;
    }, {} as Record<string, BaselineTierRecommendation[]>);
    Object.keys(tiersByGroup).forEach(group => {
        tiersByGroup[group] = [...tiersByGroup[group]].sort((left, right) => profileFieldOrder(String(left.fieldKey)) - profileFieldOrder(String(right.fieldKey)));
    });
    const tierGroupOrder = ['springs', 'aero'] as Array<keyof Copy['fieldGroups']>;
    const orderedSkippedGroups = [
        ...skippedGroupOrder.filter(group => skippedByGroup[group]?.length),
        ...Object.keys(skippedByGroup)
            .filter(group => !skippedGroupOrder.includes(group as keyof Copy['fieldGroups']) && skippedByGroup[group]?.length)
            .sort(),
    ];
    const selectedSet = new Set(selectedFields);
    const targetProfile = profiles.find(profile => profile.id === targetProfileId) || null;

    return (
        <section className="tune-generator-page">
            <div className="quick-tune-toolbar">
                <div>
                    <button className="action primary" type="button" disabled={busy} onClick={onInputOpen}>
                        <Gauge size={17}/>
                        {terms.quickTuneInputButton}
                    </button>
                    <button className="action secondary" type="button" disabled={busy} onClick={onReset}>
                        {terms.tuneGeneratorReset}
                    </button>
                </div>
                <div className="quick-tune-summary">
                    <strong>{result ? terms.quickTuneLastSummary : terms.quickTuneNoResult}</strong>
                    <span>
                        {result
                            ? `${result.profileDraft.carClass || '--'} / PI ${result.profileDraft.pi ?? '--'} / ${result.profileDraft.drivetrain || '--'} / ${selectedFields.length}`
                            : terms.tuneGeneratorNoHistory}
                    </span>
                </div>
            </div>

            {inputOpen && (
                <QuickTuneInputModal
                    form={form}
                    busy={busy}
                    terms={terms}
                    errors={fieldErrors}
                    onFieldChange={onFieldChange}
                    onGenerate={onGenerate}
                    onCancel={onInputClose}
                />
            )}

            <section className="panel tune-generator-preview">
                <div className="panel-heading">
                    <div>
                        <h2>{terms.tuneGeneratorPreview}</h2>
                        <span>{result ? `${terms.tuneGeneratorFields}: ${result.generatedFields.length}` : terms.tuneGeneratorNoPreview}</span>
                    </div>
                </div>
                {!result ? (
                    <div className="empty-events advice-placeholder">{terms.tuneGeneratorNoPreview}</div>
                ) : (
                    <>
                        <QuickTuneBiasPanel
                            form={form}
                            terms={terms}
                            hasGeneratedGearing={Boolean(generatedByGroup.gearing?.length)}
                            busy={busy}
                            onBiasChange={onBiasChange}
                        />

                        {groupOrder.map(group => previewByGroup[group]?.length ? (
                            <div className="baseline-preview-group" key={group}>
                                <h3>{terms.fieldGroups[group as keyof Copy['fieldGroups']] || group}</h3>
                                <div className="baseline-field-list">
                                    {previewByGroup[group].map(entry => {
                                        const profileField = profileFields.find(item => String(item.key) === entry.fieldKey);
                                        const value = entry.type === 'generated'
                                            ? formatBaselineGeneratedValue(entry.field.value, entry.fieldKey, profileField, language)
                                            : (localizedLabel(entry.field.tier, terms.tuneGeneratorTierLabels) || entry.field.tier);
                                        return (
                                            <label className={`baseline-generated-field ${entry.type === 'tier' ? 'tier-only' : ''}`} key={`${entry.type}-${entry.fieldKey}`}>
                                                {entry.type === 'generated' ? (
                                                    <input
                                                        type="checkbox"
                                                        checked={selectedSet.has(entry.fieldKey)}
                                                        onChange={() => onToggleField(entry.fieldKey)}
                                                    />
                                                ) : (
                                                    <span className="baseline-tier-marker" aria-hidden="true">~</span>
                                                )}
                                                <div>
                                                    <strong>{profileFieldDisplayName(entry.fieldKey, language)}</strong>
                                                    <span>{value}</span>
                                                </div>
                                            </label>
                                        );
                                    })}
                                </div>
                            </div>
                        ) : null)}

                        <div className="baseline-preview-group">
                            <h3>{terms.tuneGeneratorSkipped}</h3>
                            <div className="baseline-skipped-list">
                                {Object.keys(skippedByGroup).length === 0 ? (
                                    <div className="empty-events">--</div>
                                ) : (
                                    orderedSkippedGroups.map(group => (
                                        <div key={group}>
                                            <strong>{terms.fieldGroups[group as keyof Copy['fieldGroups']] || group}</strong>
                                            {(skippedByGroup[group] || []).map(field => (
                                                <span key={String(field.fieldKey)}>{profileFieldDisplayName(String(field.fieldKey), language)}: {field.message}</span>
                                            ))}
                                        </div>
                                    ))
                                )}
                            </div>
                        </div>

                        <div className="baseline-apply-card">
                            <label className="profile-picker">
                                <span>{terms.tuneGeneratorTarget}</span>
                                <select value={targetProfileId} onChange={event => onTargetChange(Number(event.target.value))}>
                                    <option value={0}>{terms.noProfile}</option>
                                    {profiles.map(profile => (
                                        <option key={profile.id} value={profile.id}>{[profile.carName, profile.versionName, profile.carClass].filter(Boolean).join(' / ')}</option>
                                    ))}
                                </select>
                            </label>
                            {targetProfile && <small>{formatTuneProfileVehicle(targetProfile, terms)}</small>}
                            <div className="form-actions">
                                <button className="action secondary" type="button" disabled={busy || !result} onClick={onOpenExpert}>{terms.quickTuneCarryToProfessional}</button>
                                <button className="action secondary" type="button" disabled={busy || !result} onClick={onCreate}><Plus size={17}/>{terms.tuneGeneratorCreate}</button>
                                <button className="action primary" type="button" disabled={busy || !result || !targetProfileId} onClick={onApply}><Save size={17}/>{terms.tuneGeneratorApply}</button>
                            </div>
                        </div>
                    </>
                )}
            </section>
        </section>
    );

    return (
        <section className="tune-generator-page">
            <section className="panel tune-generator-hero">
                <div className="panel-heading">
                    <div>
                        <h2>{terms.tuneGeneratorTitle}</h2>
                        <span>{terms.tuneGeneratorSubtitle}</span>
                        <span>{terms.tuneGeneratorValueHint}</span>
                    </div>
                    <span>{terms.tuneGeneratorNoHistory}</span>
                </div>
                <div className="launchpad-grid">
                    <TextStat label={terms.analysisMode} value={terms.tuneGeneratorTab}/>
                    <TextStat label={terms.sessionMode} value="Road"/>
                    <TextStat label={terms.tuneGeneratorConfidence} value={localizedLabel(result?.confidence || 'medium', terms.quickConfidenceLabels) || '--'}/>
                    <TextStat label={terms.tuneGeneratorSelectedCount} value={`${selectedFields.length}`}/>
                </div>
            </section>

            <div className="tune-generator-layout">
                <section className="panel tune-generator-form">
                    <div className="panel-heading">
                        <div>
                            <h2>{terms.tuneGeneratorMinimum}</h2>
                            <span>{terms.tuneGeneratorRangeHint}</span>
                        </div>
                    </div>
                    <div className="baseline-form-grid">
                        <BaselineInputField label={terms.tuneGeneratorCarName} value={form.carName} required onChange={value => onFieldChange('carName', value)}/>
                        <BaselineInputField label={terms.tuneGeneratorVersionName} value={form.versionName} onChange={value => onFieldChange('versionName', value)}/>
                        <BaselineInputField label={terms.tuneGeneratorPI} value={form.pi} required type="number" step={1} onChange={value => onFieldChange('pi', value)} onBlur={() => onFieldBlur('pi', 1)}/>
                        <label className="profile-field">
                            <span>{terms.tuneGeneratorDrivetrain}</span>
                            <select value={form.drivetrain} onChange={event => onFieldChange('drivetrain', event.target.value)}>
                                <option value="FWD">FWD</option>
                                <option value="RWD">RWD</option>
                                <option value="AWD">AWD</option>
                            </select>
                        </label>
                        <BaselineInputField label={`${terms.tuneGeneratorWeight} (kg)`} value={form.weightKG} required type="number" step={1} onChange={value => onFieldChange('weightKG', value)} onBlur={() => onFieldBlur('weightKG', 1)}/>
                        <BaselineInputField label={`${terms.tuneGeneratorFrontWeight} (%)`} value={form.frontWeightPct} required type="number" step={0.1} onChange={value => onFieldChange('frontWeightPct', value)} onBlur={() => onFieldBlur('frontWeightPct', 0.1)}/>
                    </div>

                    <details className="baseline-advanced" open={advancedOpen}>
                        <summary onClick={(event) => {
                            event.preventDefault();
                            onAdvancedToggle();
                        }}>{terms.tuneGeneratorAdvanced}</summary>
                        <div className="baseline-form-grid">
                            <BaselineInputField label={terms.tuneGeneratorCarOrdinal} value={form.carOrdinal} type="number" step={1} onChange={value => onFieldChange('carOrdinal', value)} onBlur={() => onFieldBlur('carOrdinal', 1)}/>
                            <BaselineInputField label={terms.tuneGeneratorCarCategory} value={form.carCategory} type="number" step={1} onChange={value => onFieldChange('carCategory', value)} onBlur={() => onFieldBlur('carCategory', 1)}/>
                            <BaselineInputField label={`${terms.tuneGeneratorPower} (kW)`} value={form.powerKW} type="number" step={1} onChange={value => onFieldChange('powerKW', value)} onBlur={() => onFieldBlur('powerKW', 1)}/>
                            <BaselineInputField label={`${terms.tuneGeneratorTorque} (Nm)`} value={form.torqueNM} type="number" step={1} onChange={value => onFieldChange('torqueNM', value)} onBlur={() => onFieldBlur('torqueNM', 1)}/>
                        </div>
                        <div className="baseline-range-grid">
                            <strong>{terms.tuneGeneratorGearingInputs}</strong>
                            <span className="baseline-tier-hint">{terms.tuneGeneratorGearingHint}</span>
                            <BaselineInputField label={`${terms.tuneGeneratorRedlineRPM} (rpm)`} value={form.redlineRPM} type="number" step={1} onChange={value => onFieldChange('redlineRPM', value)} onBlur={() => onFieldBlur('redlineRPM', 1)}/>
                            <BaselineInputField label={terms.tuneGeneratorGearCount} value={form.gearCount} type="number" step={1} onChange={value => onFieldChange('gearCount', value)} onBlur={() => onFieldBlur('gearCount', 1)}/>
                            <BaselineInputField label={`${terms.tuneGeneratorTireDiameter} (cm)`} value={form.tireDiameterCm} type="number" step={0.1} onChange={value => onFieldChange('tireDiameterCm', value)} onBlur={() => onFieldBlur('tireDiameterCm', 0.1)}/>
                            <BaselineInputField label={`${terms.tuneGeneratorTargetTopSpeed} (km/h)`} value={form.targetTopSpeedKmh} type="number" step={1} onChange={value => onFieldChange('targetTopSpeedKmh', value)} onBlur={() => onFieldBlur('targetTopSpeedKmh', 1)}/>
                        </div>
                        <div className="baseline-range-grid">
                            <strong>{terms.tuneGeneratorTierRecommendations}</strong>
                            <span className="baseline-tier-hint">{terms.tuneGeneratorAdjustableHint}</span>
                            <BaselineToggleField label={terms.tuneGeneratorFrontRideAdjustable} checked={form.frontRideHeightAdjustable === 'true'} onChange={checked => onFieldChange('frontRideHeightAdjustable', checked ? 'true' : 'false')}/>
                            <BaselineToggleField label={terms.tuneGeneratorRearRideAdjustable} checked={form.rearRideHeightAdjustable === 'true'} onChange={checked => onFieldChange('rearRideHeightAdjustable', checked ? 'true' : 'false')}/>
                            <BaselineToggleField label={terms.tuneGeneratorFrontAeroAdjustable} checked={form.frontAeroAdjustable === 'true'} onChange={checked => onFieldChange('frontAeroAdjustable', checked ? 'true' : 'false')}/>
                            <BaselineToggleField label={terms.tuneGeneratorRearAeroAdjustable} checked={form.rearAeroAdjustable === 'true'} onChange={checked => onFieldChange('rearAeroAdjustable', checked ? 'true' : 'false')}/>
                        </div>
                    </details>

                    <div className="form-actions">
                        <button className="action secondary" type="button" disabled={busy} onClick={onReset}>{terms.tuneGeneratorReset}</button>
                    <button className="action primary" type="button" disabled={busy} onClick={() => onGenerate()}><Gauge size={17}/>{terms.tuneGeneratorGenerate}</button>
                    </div>
                </section>

                <section className="panel tune-generator-preview">
                    <div className="panel-heading">
                        <div>
                            <h2>{terms.tuneGeneratorPreview}</h2>
                            <span>{result ? `${terms.tuneGeneratorFields}: ${result?.generatedFields?.length || 0}` : terms.tuneGeneratorNoPreview}</span>
                        </div>
                    </div>
                    {!result ? (
                        <div className="empty-events advice-placeholder">{terms.tuneGeneratorNoPreview}</div>
                    ) : (
                        <>
                            {groupOrder.map(group => generatedByGroup[group]?.length ? (
                                <div className="baseline-preview-group" key={group}>
                                    <h3>{terms.fieldGroups[group as keyof Copy['fieldGroups']] || group}</h3>
                                    <div className="baseline-field-list">
                                        {generatedByGroup[group].map(field => {
                                            const profileField = profileFields.find(item => String(item.key) === field.fieldKey);
                                            return (
                                                <label className="baseline-generated-field" key={String(field.fieldKey)}>
                                                    <input
                                                        type="checkbox"
                                                        checked={selectedSet.has(String(field.fieldKey))}
                                                        onChange={() => onToggleField(String(field.fieldKey))}
                                                    />
                                                    <div>
                                                        <strong>{profileFieldDisplayName(String(field.fieldKey), language)}</strong>
                                                        <span>{formatBaselineGeneratedValue(field.value, String(field.fieldKey), profileField, language)}</span>
                                                        <small>{terms.tuneGeneratorReason}: {field.reason}</small>
                                                    </div>
                                                </label>
                                            );
                                        })}
                                    </div>
                                </div>
                            ) : null)}

                            {((result?.tierRecommendations || []).length > 0) && (
                                <div className="baseline-preview-group">
                                    <h3>{terms.tuneGeneratorTierRecommendations}</h3>
                                    <div className="baseline-tier-list">
                                        {tierGroupOrder.map(group => tiersByGroup[group]?.length ? (
                                            <div key={group}>
                                                <strong>{terms.fieldGroups[group as keyof Copy['fieldGroups']] || group}</strong>
                                                {tiersByGroup[group].map(field => (
                                                    <div className={`baseline-tier-row ${field.applicable ? '' : 'disabled'}`} key={String(field.fieldKey)}>
                                                        <span>{profileFieldDisplayName(String(field.fieldKey), language)}</span>
                                                        <em>{localizedLabel(field.tier, terms.tuneGeneratorTierLabels) || field.tier}</em>
                                                        <small>{field.applicable ? terms.tuneGeneratorTierManual : field.reason}</small>
                                                        <small>{terms.tuneGeneratorReason}: {field.reason}</small>
                                                    </div>
                                                ))}
                                            </div>
                                        ) : null)}
                                    </div>
                                </div>
                            )}

                            <div className="baseline-preview-group">
                                <h3>{terms.tuneGeneratorSkipped}</h3>
                                <div className="baseline-skipped-list">
                                    {Object.keys(skippedByGroup).length === 0 ? (
                                        <div className="empty-events">--</div>
                                    ) : (
                                        orderedSkippedGroups.map(group => (
                                            <div key={group}>
                                                <strong>{terms.fieldGroups[group as keyof Copy['fieldGroups']] || group}</strong>
                                                {(skippedByGroup[group] || []).map(field => (
                                                    <span key={String(field.fieldKey)}>{profileFieldDisplayName(String(field.fieldKey), language)}: {field.message}</span>
                                                ))}
                                            </div>
                                        ))
                                    )}
                                </div>
                            </div>

                            <div className="baseline-preview-group">
                                <h3>{terms.tuneGeneratorNextTest}</h3>
                                <ul className="baseline-next-test">
                                    {(result?.nextTestPlan || []).map(item => <li key={item}>{item}</li>)}
                                </ul>
                            </div>

                            <div className="baseline-apply-card">
                                <label className="profile-picker">
                                    <span>{terms.tuneGeneratorTarget}</span>
                                    <select value={targetProfileId} onChange={event => onTargetChange(Number(event.target.value))}>
                                        <option value={0}>{terms.noProfile}</option>
                                        {profiles.map(profile => (
                                            <option key={profile.id} value={profile.id}>{[profile.carName, profile.versionName, profile.carClass].filter(Boolean).join(' / ')}</option>
                                        ))}
                                    </select>
                                </label>
                                {targetProfile && <small>{formatTuneProfileVehicle(targetProfile as TuneProfile, terms)}</small>}
                                <div className="form-actions">
                                    <button className="action secondary" type="button" disabled={busy || !result} onClick={onCreate}><Plus size={17}/>{terms.tuneGeneratorCreate}</button>
                                    <button className="action primary" type="button" disabled={busy || !result || !targetProfileId} onClick={onApply}><Save size={17}/>{terms.tuneGeneratorApply}</button>
                                    <button className="action secondary" type="button" onClick={onOpenExpert}>{terms.tuneGeneratorOpenExpert}</button>
                                </div>
                            </div>
                        </>
                    )}
                </section>
            </div>
        </section>
    );
}

function BaselineInputField({label, value, required, type = 'text', step, inputMode, error, onChange, onBlur}: {
    label: string;
    value: string;
    required?: boolean;
    type?: string;
    step?: number;
    inputMode?: 'none' | 'text' | 'tel' | 'url' | 'email' | 'numeric' | 'decimal' | 'search';
    error?: string;
    onChange: (value: string) => void;
    onBlur?: () => void;
}) {
    return (
        <label className={`profile-field ${error ? 'has-error' : ''}`}>
            <span>{label}{required ? ' *' : ''}</span>
            <input type={type} step={step} inputMode={inputMode} value={value} onChange={event => onChange(event.target.value)} onBlur={onBlur}/>
            {error && <small className="field-error">{error}</small>}
        </label>
    );
}

function TireSizeInputField({
    label,
    width,
    aspect,
    rim,
    required,
    errors,
    onFieldChange,
}: {
    label: string;
    width: string;
    aspect: string;
    rim: string;
    required?: boolean;
    errors: Pick<QuickTuneFieldErrors, 'tireWidthMm' | 'tireAspectRatio' | 'tireRimInches'>;
    onFieldChange: (key: RoadStaticTuneBaselineFormKey, value: string) => void;
}) {
    const error = errors.tireWidthMm || errors.tireAspectRatio || errors.tireRimInches || '';
    return (
        <label className={`profile-field tire-size-field ${error ? 'has-error' : ''}`}>
            <span>{label}{required ? ' *' : ''}</span>
            <div className="tire-size-input">
                <input
                    type="number"
                    step={1}
                    inputMode="numeric"
                    value={width}
                    placeholder="245"
                    aria-label={`${label} width`}
                    onChange={event => onFieldChange('tireWidthMm', event.target.value)}
                />
                <span>/</span>
                <input
                    type="number"
                    step={1}
                    inputMode="numeric"
                    value={aspect}
                    placeholder="35"
                    aria-label={`${label} aspect ratio`}
                    onChange={event => onFieldChange('tireAspectRatio', event.target.value)}
                />
                <span>R</span>
                <input
                    type="number"
                    step={1}
                    inputMode="numeric"
                    value={rim}
                    placeholder="19"
                    aria-label={`${label} rim`}
                    onChange={event => onFieldChange('tireRimInches', event.target.value)}
                />
            </div>
            {error && <small className="field-error">{error}</small>}
        </label>
    );
}

function QuickTuneBiasPanel({
    form,
    terms,
    hasGeneratedGearing,
    busy,
    onBiasChange,
}: {
    form: RoadStaticTuneBaselineForm;
    terms: Copy;
    hasGeneratedGearing: boolean;
    busy: boolean;
    onBiasChange: (key: 'balanceBias' | 'stiffnessBias' | 'speedBias', value: string) => void;
}) {
    const gearingAvailable = form.gearingEnabled === 'true' && hasGeneratedGearing;
    return (
        <div className="quick-tune-bias-card">
            <div className="quick-tune-bias-heading">
                <div>
                    <h3>{terms.quickTuneBiasTitle}</h3>
                    <span>{terms.quickTuneBiasHint}</span>
                </div>
            </div>
            <div className="quick-tune-bias-grid">
                <QuickTuneBiasSlider
                    label={terms.quickTuneBalance}
                    value={form.balanceBias}
                    left={terms.quickTuneBalanceLeft}
                    center={terms.quickTuneNeutral}
                    right={terms.quickTuneBalanceRight}
                    disabled={busy}
                    onChange={value => onBiasChange('balanceBias', value)}
                />
                <QuickTuneBiasSlider
                    label={terms.quickTuneStiffness}
                    value={form.stiffnessBias}
                    left={terms.quickTuneStiffnessLeft}
                    center={terms.quickTuneNeutral}
                    right={terms.quickTuneStiffnessRight}
                    disabled={busy}
                    onChange={value => onBiasChange('stiffnessBias', value)}
                />
                <QuickTuneBiasSlider
                    label={terms.quickTuneSpeed}
                    value={form.speedBias}
                    left={terms.quickTuneSpeedLeft}
                    center={terms.quickTuneNeutral}
                    right={terms.quickTuneSpeedRight}
                    disabled={busy || !gearingAvailable}
                    disabledHint={gearingAvailable ? '' : terms.quickTuneSpeedDisabled}
                    onChange={value => onBiasChange('speedBias', value)}
                />
            </div>
        </div>
    );
}

function QuickTuneBiasSlider({
    label,
    value,
    left,
    center,
    right,
    disabled,
    disabledHint,
    onChange,
}: {
    label: string;
    value: string;
    left: string;
    center: string;
    right: string;
    disabled: boolean;
    disabledHint?: string;
    onChange: (value: string) => void;
}) {
    const currentValue = quickTuneBiasSliderValue(value);
    return (
        <label className={`quick-tune-bias-slider ${disabled ? 'disabled' : ''}`}>
            <span>
                <strong>{label}</strong>
                <em>{currentValue}</em>
            </span>
            <input
                type="range"
                min={50}
                max={150}
                step={1}
                value={currentValue}
                disabled={disabled}
                onChange={event => onChange(String(quickTuneBiasDetentValue(Number(event.target.value))))}
            />
            <small>
                <span>{left}</span>
                <span>{center}</span>
                <span>{right}</span>
            </small>
            {disabledHint && <small className="quick-tune-bias-disabled">{disabledHint}</small>}
        </label>
    );
}

function quickTuneBiasSliderValue(value: string) {
    const parsed = Number(value);
    if (!Number.isFinite(parsed)) {
        return 100;
    }
    return Math.min(150, Math.max(50, Math.round(parsed)));
}

function quickTuneBiasDetentValue(value: number) {
    if (!Number.isFinite(value)) {
        return 100;
    }
    const rounded = Math.min(150, Math.max(50, Math.round(value)));
    return Math.abs(rounded - 100) <= 3 ? 100 : rounded;
}

function quickTuneUseCaseSupported(useCase: string) {
    return ['Road', 'Drift', 'Rally', 'Offroad', 'Drag'].includes(useCase);
}

function quickTuneGearingHint(useCase: string, terms: Copy) {
    if (useCase === 'Drift') {
        return terms.quickTuneDriftGearingHint;
    }
    if (useCase === 'Drag') {
        return terms.quickTuneDragGearingHint;
    }
    return terms.tuneGeneratorGearingHint;
}

function quickTuneTargetSpeedLabel(useCase: string, terms: Copy) {
    if (useCase === 'Drift') {
        return terms.quickTuneTargetDriftSpeed;
    }
    if (useCase === 'Drag') {
        return terms.quickTuneTargetDragSpeed;
    }
    return terms.tuneGeneratorTargetTopSpeed;
}

function QuickTuneInputModal({
    form,
    busy,
    terms,
    errors,
    onFieldChange,
    onGenerate,
    onCancel,
}: {
    form: RoadStaticTuneBaselineForm;
    busy: boolean;
    terms: Copy;
    errors: QuickTuneFieldErrors;
    onFieldChange: (key: RoadStaticTuneBaselineFormKey, value: string) => void;
    onGenerate: () => void;
    onCancel: () => void;
}) {
    const gearingEnabled = form.gearingEnabled === 'true';
    const unsupportedUseCase = !quickTuneUseCaseSupported(form.useCase);
    const driftRwdWarning = form.useCase === 'Drift' && form.drivetrain !== 'RWD';
    const gearingHint = quickTuneGearingHint(form.useCase, terms);
    const targetSpeedLabel = quickTuneTargetSpeedLabel(form.useCase, terms);
    const errorMessages = Array.from(new Set(Object.values(errors).filter(Boolean) as string[]));
    return (
        <div className="modal-backdrop" role="presentation">
            <section className="modal-card quick-tune-modal" role="dialog" aria-modal="true" aria-label={terms.quickTuneInputTitle}>
                <div className="panel-heading">
                    <div>
                        <h2>{terms.quickTuneInputTitle}</h2>
                        <span>{terms.tuneGeneratorValueHint}</span>
                    </div>
                    <button className="small-action" type="button" onClick={onCancel} disabled={busy}>{terms.close}</button>
                </div>
                {errorMessages.length > 0 && (
                    <div className="quick-tune-error-summary">
                        <strong>{terms.quickTuneValidationSummary}</strong>
                        <span>{errorMessages[0]}</span>
                    </div>
                )}
                <div className="baseline-form-grid">
                    <BaselineInputField label={terms.tuneGeneratorCarName} value={form.carName} onChange={value => onFieldChange('carName', value)}/>
                    <BaselineInputField label={`${terms.tuneGeneratorWeight} (kg)`} value={form.weightKG} required type="number" step={1} inputMode="numeric" error={errors.weightKG} onChange={value => onFieldChange('weightKG', value)}/>
                    <BaselineInputField label={`${terms.tuneGeneratorFrontWeight} (%)`} value={form.frontWeightPct} required type="number" step={1} inputMode="numeric" error={errors.frontWeightPct} onChange={value => onFieldChange('frontWeightPct', value)}/>
                    <BaselineInputField label={terms.tuneGeneratorPI} value={form.pi} required type="number" step={1} inputMode="numeric" error={errors.pi} onChange={value => onFieldChange('pi', value)}/>
                    <label className="profile-field">
                        <span>{terms.quickTuneUseCase}</span>
                        <select value={form.useCase} onChange={event => onFieldChange('useCase', event.target.value)}>
                            {quickTuneUseCases.map(useCase => (
                                <option key={useCase} value={useCase}>{terms.useCases[useCase] || useCase}</option>
                            ))}
                        </select>
                    </label>
                    <label className="profile-field">
                        <span>{terms.quickTuneTireCompound}</span>
                        <select value={form.tireCompound || 'sport'} onChange={event => onFieldChange('tireCompound', event.target.value)}>
                            {quickTuneTireCompounds.map(compound => (
                                <option key={compound} value={compound}>{terms.quickTuneTireCompoundLabels[compound] || compound}</option>
                            ))}
                        </select>
                    </label>
                    <div className="profile-field">
                        <span>{terms.tuneGeneratorDrivetrain}</span>
                        <div className="segmented-control">
                            {(['FWD', 'AWD', 'RWD'] as const).map(value => (
                                <button
                                    key={value}
                                    type="button"
                                    className={form.drivetrain === value ? 'active' : ''}
                                    onClick={() => onFieldChange('drivetrain', value)}
                                >
                                    {terms.quickTuneDrivetrainLabels[value]}
                                </button>
                            ))}
                        </div>
                    </div>
                </div>

                {unsupportedUseCase && <div className="status-alert warn">{terms.quickTuneUnsupportedUseCase}</div>}
                {driftRwdWarning && <div className="status-alert warn">{terms.quickTuneDriftRwdPreferred}</div>}

                <div className="baseline-range-grid quick-tune-gearing">
                    <BaselineToggleField label={terms.quickTuneGearingToggle} checked={gearingEnabled} onChange={checked => onFieldChange('gearingEnabled', checked ? 'true' : 'false')}/>
                    {gearingEnabled && <span className="baseline-tier-hint">{gearingHint}</span>}
                    {gearingEnabled && (
                        <>
                            <BaselineInputField label={`${terms.tuneGeneratorRedlineRPM} (rpm)`} value={form.redlineRPM} required type="number" step={1} inputMode="numeric" error={errors.redlineRPM} onChange={value => onFieldChange('redlineRPM', value)}/>
                            <BaselineInputField label={terms.tuneGeneratorGearCount} value={form.gearCount} required type="number" step={1} inputMode="numeric" error={errors.gearCount} onChange={value => onFieldChange('gearCount', value)}/>
                            <TireSizeInputField
                                label={terms.tuneGeneratorTireDiameter}
                                width={form.tireWidthMm}
                                aspect={form.tireAspectRatio}
                                rim={form.tireRimInches}
                                required
                                errors={{
                                    tireWidthMm: errors.tireWidthMm,
                                    tireAspectRatio: errors.tireAspectRatio,
                                    tireRimInches: errors.tireRimInches,
                                }}
                                onFieldChange={onFieldChange}
                            />
                            <BaselineInputField label={`${targetSpeedLabel} (km/h)`} value={form.targetTopSpeedKmh} required type="number" step={1} inputMode="numeric" error={errors.targetTopSpeedKmh} onChange={value => onFieldChange('targetTopSpeedKmh', value)}/>
                        </>
                    )}
                </div>

                <div className="form-actions">
                    <button className="action secondary" type="button" onClick={onCancel} disabled={busy}>{terms.close}</button>
                    <button className="action primary" type="button" onClick={() => onGenerate()} disabled={busy || unsupportedUseCase}>
                        <Gauge size={17}/>
                        {terms.tuneGeneratorGenerate}
                    </button>
                </div>
            </section>
        </div>
    );
}

function BaselineToggleField({label, checked, onChange}: {
    label: string;
    checked: boolean;
    onChange: (checked: boolean) => void;
}) {
    return (
        <label className="baseline-toggle-field">
            <input type="checkbox" checked={checked} onChange={event => onChange(event.target.checked)}/>
            <span>{label}</span>
        </label>
    );
}

function QuickLapCard({title, lap, terms}: { title: string; lap: QuickLapSummary; terms: Copy }) {
    return (
        <div className="quick-lap-card">
            <strong>{title} #{lap.lapNumber}</strong>
            <div className="mini-stat-grid">
                <TextStat label={terms.samplesSaved} value={`${lap.sampleCount}`}/>
                <TextStat label={terms.duration} value={formatDuration(lap.durationMs, terms)}/>
                <TextStat label={terms.avgSpeed} value={`${formatNumber(lap.avgSpeedKmh, 0)} km/h`}/>
                <TextStat label={terms.issueScore} value={formatNumber(lap.issueScore, 0)}/>
            </div>
        </div>
    );
}

function GearSpeedValue({gear, terms}: { gear: GearPowerBand; terms: Copy }) {
    return (
        <strong className="gear-speed-value">
            <span>{formatGearSpeedRange(gear)}</span>
            <small>{terms.highestObserved}: {formatGearMaxObserved(gear)}</small>
        </strong>
    );
}

function QuickGearPowerDiagnosticCard({gearPower, terms}: { gearPower: GearPowerDiagnostic | undefined; terms: Copy }) {
    const gearRows = (gearPower?.gears || []).filter(gear => gear.sampleCount > 0).slice(0, 10);
    const actions = gearPower?.recommendedActions || [];
    const targetMin = gearPower?.powerBandStartRPM || 0;
    const targetMax = gearPower?.powerBandEndRPM || 0;
    const targetLabel = targetMin > 0 && targetMax > 0
        ? `${formatNumber(targetMin, 0)} - ${formatNumber(targetMax, 0)} rpm`
        : `${formatPercentValue(gearPower?.evidence?.power_band_target_min ?? 0.55)} - ${formatPercentValue(gearPower?.evidence?.power_band_target_max ?? 0.9)}`;
    const reasonList = gearPowerNoAdviceReasons(gearPower, terms);
    return (
        <div className="quick-section quick-gear-power-card">
            <div className="section-title-row">
                <div>
                    <h3>{terms.gearPowerDiagnostic}</h3>
                    <span>{terms.gearPowerDiagnosticHint}</span>
                </div>
                <strong>{localizedLabel(gearPower?.summary || 'not_enough_samples', terms.gearFindings)}</strong>
            </div>
            <div className="quick-comparability-grid">
                <TextStat label={terms.powerBandTarget} value={targetLabel}/>
                <TextStat label={terms.powerBandSource} value={localizedLabel(gearPower?.powerBandSource || '', terms.powerBandSources) || '--'}/>
                <TextStat label={terms.diagnosticConfidence} value={localizedLabel(gearPower?.confidence || '', terms.planConfidenceLabels) || '--'}/>
                <TextStat label={terms.gearStrategyMode} value={localizedLabel(gearPower?.strategyMode || '', terms.gearFindings) || '--'}/>
                <TextStat label={terms.gearStrategyIssueCount} value={formatGearStrategyIssueCount(gearPower)}/>
                <TextStat label={terms.tractionLimited} value={formatPercentValue(gearPower?.tractionLimitedPercent || 0)}/>
            </div>
            {gearRows.length > 0 ? (
                <div className="diagnostic-table compact-diagnostic quick-gear-table">
                    <div className="diagnostic-row diagnostic-head">
                        <span>{terms.gear}</span>
                        <span>{terms.speedRange}</span>
                        <span>RPM</span>
                        <span>{terms.inPowerBand}</span>
                        <span>{terms.tractionLimited}</span>
                        <span>{terms.gearFinding}</span>
                        <span>{terms.samplesSaved}</span>
                    </div>
                    {gearRows.map(gear => (
                        <div className={`diagnostic-row ${gear.finding !== 'ok' ? 'warn' : ''}`} key={gear.gear}>
                            <span>{gear.gear}</span>
                            <GearSpeedValue gear={gear} terms={terms}/>
                            <span>{gear.rpmAvg ? formatNumber(gear.rpmAvg, 0) : formatPercentValue(gear.rpmRatioAvg)}</span>
                            <span>{formatGearInPowerBandRange(gear)}</span>
                            <span>{formatPercentValue(gear.tractionLimitedPercent || 0)}</span>
                            <span>{localizedLabel(gear.finding, terms.gearFindings)}</span>
                            <span>{gear.highLoadSampleCount}/{gear.sampleCount}</span>
                        </div>
                    ))}
                </div>
            ) : (
                <div className="empty-events advice-placeholder">{terms.gearPowerNeedSamples}</div>
            )}
            {actions.length > 0 ? (
                <div className="quick-suggestion-list">
                    {actions.map((action, index) => (
                        <div className="quick-suggestion-row" key={`${action.category}-${action.item}-${action.direction}-${index}`}>
                            <div>
                                <strong>{localizedLabel(action.category, terms.actionCategories)} / {localizedLabel(action.item, terms.actionItems)}</strong>
                                <span>{localizedLabel(action.direction, terms.actionDirections)} / {localizedLabel('direction_only', terms.actionAmounts)}</span>
                            </div>
                            <small>{localizedLabel(action.reason, terms.actionReasons)}</small>
                        </div>
                    ))}
                    <div className="empty-events advice-placeholder">{terms.quickGearAdviceReadOnly}</div>
                </div>
            ) : (
                <div className="quick-section-subtle">
                    <strong>{terms.gearPowerWhyNoAdvice}</strong>
                    {reasonList.map(reason => <span key={reason}>{reason}</span>)}
                </div>
            )}
        </div>
    );
}

function ControlBar({label, value, tone}: { label: string; value: number; tone: 'green' | 'red' }) {
    const percent = clamp(value, 0, 1) * 100;
    return (
        <div className="control-bar">
            <div className="control-label">
                <span>{label}</span>
                <strong>{percent.toFixed(0)}%</strong>
            </div>
            <div className="bar-track">
                <div className={`bar-fill ${tone}`} style={{width: `${percent}%`}}/>
            </div>
        </div>
    );
}

function SteerMeter({label, value}: { label: string; value: number }) {
    const clamped = clamp(value, -1, 1);
    return (
        <div className="steer-meter">
            <div className="control-label">
                <span>{label}</span>
                <strong>{clamped.toFixed(2)}</strong>
            </div>
            <div className="steer-track">
                <div className="steer-center"/>
                <div className="steer-marker" style={{left: `${(clamped + 1) * 50}%`}}/>
            </div>
        </div>
    );
}

function WheelCard({label, wheel, terms}: { label: string; wheel: WheelTelemetry; terms: Copy }) {
    const health = wheelHealth(wheel.combinedSlip);
    return (
        <div className={`wheel-card ${health}`}>
            <div className="wheel-title">
                <span>{label}</span>
                <Thermometer size={16}/>
            </div>
            <div className="wheel-values">
                <span>{terms.ratio} <strong>{wheel.slipRatio.toFixed(2)}</strong></span>
                <span>{terms.angle} <strong>{wheel.slipAngle.toFixed(2)}</strong></span>
                <span>{terms.combined} <strong>{wheel.combinedSlip.toFixed(2)}</strong></span>
                <span>{terms.temp} <strong>{wheel.tireTemp.toFixed(0)}</strong></span>
                <span>{terms.susp} <strong>{(wheel.suspensionTravel * 100).toFixed(0)}%</strong></span>
                <span>{terms.surfaceSignals} <strong>{wheel.rumbleStrip.toFixed(0)} / {wheel.puddleDepth.toFixed(2)}</strong></span>
            </div>
        </div>
    );
}

function Stat({label, value}: { label: string; value: number }) {
    return (
        <div className="stat">
            <span>{label}</span>
            <strong>{value}</strong>
        </div>
    );
}

function TextStat({label, value}: { label: string; value: string }) {
    return (
        <div className="stat">
            <span>{label}</span>
            <strong>{value}</strong>
        </div>
    );
}

function Trend({title, unit, points}: { title: string; unit: string; points: string }) {
    return (
        <div className="trend">
            <div className="trend-title">
                <span>{title}</span>
                <small>{unit}</small>
            </div>
            <svg viewBox="0 0 300 80" role="img" aria-label={`${title} trend`}>
                <polyline points={points} fill="none" stroke="currentColor" strokeWidth="3" strokeLinecap="round" strokeLinejoin="round"/>
            </svg>
        </div>
    );
}

function ProfileManager({
    profiles,
    activeProfile,
    sessionStats,
    snapshots,
    profileForm,
    current,
    formRef,
    editingProfileId,
    busy,
    language,
    terms,
    influenceMap,
    onFieldChange,
    onNew,
    onEdit,
    onDuplicate,
    onRestoreSnapshot,
    onDelete,
    onSetActive,
    onFillFromTelemetry,
    canFillFromTelemetry,
    onSave,
}: {
    profiles: TuneProfile[];
    activeProfile: TuneProfile | null;
    sessionStats: TuneProfileSessionStat[];
    snapshots: TuneProfileSnapshot[];
    profileForm: TuneProfileInput;
    current: TelemetryFrame | null;
    formRef: RefObject<HTMLDivElement>;
    editingProfileId: number | null;
    busy: boolean;
    language: Lang;
    terms: Copy;
    influenceMap: TuneToTireInfluenceMap | null;
    onFieldChange: (field: keyof TuneProfileInput, value: string | number | null) => void;
    onNew: () => void;
    onEdit: (profile: TuneProfile) => void;
    onDuplicate: (profile: TuneProfile) => void;
    onRestoreSnapshot: (snapshot: TuneProfileSnapshot) => void;
    onDelete: (profile: TuneProfile) => void;
    onSetActive: (id: number) => void;
    onFillFromTelemetry: () => void;
    canFillFromTelemetry: boolean;
    onSave: () => void;
}) {
    const [snapshotForCompare, setSnapshotForCompare] = useState<TuneProfileSnapshot | null>(null);
    const [selectedInfluence, setSelectedInfluence] = useState<TuneFieldInfluence | null>(null);
    const [editingNumberField, setEditingNumberField] = useState<{ key: keyof TuneProfileInput; value: string } | null>(null);
    const fieldsByGroup = profileFields.reduce((acc, field) => {
        acc[field.group] = [...(acc[field.group] || []), field];
        return acc;
    }, {} as Record<keyof Copy['fieldGroups'], ProfileField[]>);
    Object.keys(fieldsByGroup).forEach(group => {
        const key = group as keyof Copy['fieldGroups'];
        fieldsByGroup[key] = [...(fieldsByGroup[key] || [])].sort((left, right) => profileFieldOrder(String(left.key)) - profileFieldOrder(String(right.key)));
    });
    const profileTopGroupOrder = ['vehicle', 'power'] as Array<keyof Copy['fieldGroups']>;
    const aeroAdvancedFields = (fieldsByGroup.aero || []).filter(field => !coreTuneFieldOrderSet.has(String(field.key)));
    const profileGroupSections = [
        ...profileTopGroupOrder.map(group => ({group, fields: fieldsByGroup[group] || []})),
        ...coreTuneGroupOrder.map(group => ({
            group,
            fields: (fieldsByGroup[group] || []).filter(field => coreTuneFieldOrderSet.has(String(field.key))),
        })),
        {group: 'gearing' as keyof Copy['fieldGroups'], fields: fieldsByGroup.gearing || []},
        ...(aeroAdvancedFields.length ? [{group: 'aero' as keyof Copy['fieldGroups'], fields: aeroAdvancedFields}] : []),
        {group: 'notes' as keyof Copy['fieldGroups'], fields: fieldsByGroup.notes || []},
    ].filter(section => section.fields.length > 0);
    const statByProfile = new Map(sessionStats.map(stat => [stat.tuneProfileId, stat]));
    const editingProfile = profiles.find(profile => profile.id === editingProfileId) || null;
    const currentComparableProfile = profileInputToProfile(profileForm, editingProfile);
    const profileDiffs = compareProfiles(snapshotForCompare?.before || null, snapshotForCompare ? currentComparableProfile : null, language);
    const profileTelemetryState = profileFormTelemetryState(profileForm, current);
    const numberInputValue = (field: ProfileField, value: unknown) => (
        field.kind === 'number' && editingNumberField?.key === field.key ? editingNumberField.value : profileInputValue(field, value)
    );
    const onNumberFocus = (field: ProfileField, value: unknown) => {
        if (field.kind === 'number' && !field.readOnly) {
            setEditingNumberField({key: field.key, value: value === undefined || value === null ? '' : String(value)});
        }
    };
    const onNumberChange = (field: ProfileField, raw: string) => {
        setEditingNumberField({key: field.key, value: raw});
        onFieldChange(field.key, parseOptionalNumber(raw));
    };
    const onNumberBlur = (field: ProfileField) => {
        setEditingNumberField(current => current?.key === field.key ? null : current);
    };
    const influenceByField = useMemo(() => new Map((influenceMap?.items || []).map(item => [item.fieldKey, item])), [influenceMap]);
    const openFieldInfluence = async (fieldKey: string) => {
        const cached = influenceByField.get(fieldKey);
        if (cached) {
            setSelectedInfluence(cached);
            return;
        }
        try {
            const item = await ExplainTuneFieldInfluence(fieldKey);
            setSelectedInfluence(item as TuneFieldInfluence);
        } catch {
            setSelectedInfluence(null);
        }
    };

    useEffect(() => {
        setEditingNumberField(null);
    }, [editingProfileId]);

    return (
        <section className="profiles-layout">
            <div className="panel profile-list-panel">
                <div className="panel-heading">
                    <h2>{terms.profileList}</h2>
                    <button className="small-action" type="button" onClick={onNew}>
                        <Plus size={15}/>
                        {terms.newProfile}
                    </button>
                </div>
                <div className="profile-list">
                    {profiles.length === 0 ? (
                        <div className="empty-events profile-empty">
                            <span>{terms.noProfiles}</span>
                        </div>
                    ) : profiles.map(profile => (
                        <div className={`profile-row ${activeProfile?.id === profile.id ? 'active' : ''}`} key={profile.id}>
                            <div>
                                <strong>{profile.carName}</strong>
                                <span>{[profile.versionName, profile.carClass, profile.drivetrain, localizedUseCase(profile.useCase, terms)].filter(Boolean).join(' / ') || '--'}</span>
                                <small>{terms.profileSessions}: {statByProfile.get(profile.id)?.sessionCount || 0} / {terms.recentSession}: {formatTime(statByProfile.get(profile.id)?.lastStartedAt || '', '--')}</small>
                            </div>
                            <div className="profile-row-actions">
                                {activeProfile?.id === profile.id && <span className="active-chip">{terms.active}</span>}
                                <button type="button" onClick={() => onEdit(profile)} title={terms.updateProfile}><Pencil size={15}/></button>
                            </div>
                        </div>
                    ))}
                </div>
            </div>

            <div className="panel profile-form-panel" ref={formRef}>
                {!editingProfileId ? (
                    <div className="profile-editor-empty">
                        <h2>{terms.expertWorkspace}</h2>
                        <p>{terms.selectProfileToEdit}</p>
                        <button className="action primary" type="button" onClick={onNew}>
                            <Plus size={16}/>
                            {terms.newProfile}
                        </button>
                    </div>
                ) : (
                    <>
                <div className="panel-heading profile-form-heading">
                    <div>
                        <h2>{terms.profileForm}</h2>
                        <span>{editingProfileId ? terms.updateProfile : terms.createProfile}</span>
                    </div>
                    <div className="profile-form-heading-actions">
                        <button className="small-action" type="button" onClick={onFillFromTelemetry} disabled={busy || !canFillFromTelemetry}>
                            <Gauge size={15}/>
                            {terms.fillFromTelemetry}
                        </button>
                        <button className="small-action" type="button" onClick={onSave} disabled={busy}>
                            <Save size={15}/>
                            {editingProfileId ? terms.updateProfile : terms.createProfile}
                        </button>
                        <button className="small-action" type="button" onClick={() => editingProfile && onSetActive(editingProfile.id)} disabled={busy || !editingProfile || activeProfile?.id === editingProfile.id}>
                            {terms.setActive}
                        </button>
                        <button className="small-action danger" type="button" onClick={() => editingProfile && onDelete(editingProfile)} disabled={busy || !editingProfile}>
                            <Trash2 size={15}/>
                            {terms.delete}
                        </button>
                    </div>
                </div>
                <div className={`profile-identity-card ${profileTelemetryState}`}>
                    <div>
                        <strong>{terms.profileIdentity}</strong>
                        <span>{terms.profileIdentityHint}</span>
                    </div>
                    <div className="identity-grid">
                        <TextStat label={terms.vehicleId} value={formatOptionalInt(profileForm.carOrdinal ?? undefined)}/>
                        <TextStat label={terms.classPi} value={`${profileForm.carClass || '--'} / ${formatOptionalInt(profileForm.pi ?? undefined)}`}/>
                        <TextStat label={terms.drivetrainCylinders} value={`${profileForm.drivetrain || '--'} / ${formatOptionalInt(profileForm.numCylinders ?? undefined)}`}/>
                        <TextStat label={terms.profileForm} value={[profileForm.carName, profileForm.versionName, localizedUseCase(profileForm.useCase, terms)].filter(Boolean).join(' / ') || '--'}/>
                    </div>
                    <div className="identity-status">
                        {profileTelemetryState === 'match' ? terms.profileTelemetryMatch : profileTelemetryState === 'mismatch' ? terms.profileTelemetryMismatch : terms.profileTelemetryUnavailable}
                    </div>
                </div>
                {editingProfile && (
                    <section className="profile-snapshot-section">
                        <div className="section-title-row">
                            <div>
                                <h3>{terms.recentChanges}</h3>
                                <span>{terms.openProfileCompare}</span>
                            </div>
                            <details className="profile-more-actions">
                                <summary>{terms.moreActions}</summary>
                                <button className="small-action" type="button" onClick={() => onDuplicate(editingProfile)} disabled={busy}>
                                    <CopyIcon size={15}/>
                                    {terms.duplicate}
                                </button>
                            </details>
                        </div>
                        {snapshots.length === 0 ? (
                            <div className="empty-events">{terms.noRecentChanges}</div>
                        ) : (
                            <div className="snapshot-list">
                                {snapshots.map(snapshot => (
                                    <div className="snapshot-row" key={snapshot.id}>
                                        <div>
                                            <strong>{formatTime(snapshot.changedAt, '--')}</strong>
                                            <span>{terms.snapshotChangedCount(snapshot.changedFields.length)}</span>
                                            <small>{formatSnapshotChangedFields(snapshot.changedFields, language)}</small>
                                        </div>
                                        <div className="snapshot-actions">
                                            <button className="small-action" type="button" onClick={() => setSnapshotForCompare(snapshot)}>
                                                {terms.compareWithCurrent}
                                            </button>
                                            <button className="small-action" type="button" onClick={() => onRestoreSnapshot(snapshot)} disabled={busy}>
                                                {terms.restoreSnapshot}
                                            </button>
                                        </div>
                                    </div>
                                ))}
                            </div>
                        )}
                    </section>
                )}
                {profileGroupSections.map(section => (
                    <div className="profile-field-group" key={`${section.group}-${section.fields.map(field => String(field.key)).join('-')}`}>
                        <h3>{terms.fieldGroups[section.group]}</h3>
                        <div className="profile-field-grid">
                            {section.fields.map(field => {
                                const value = profileForm[field.key];
                                const locked = isLockedProfileField(field, value);
                                const className = `${field.kind === 'textarea' ? 'profile-field wide' : 'profile-field'}${locked ? ' locked' : ''}`;
                                return (
                                    <label className={className} key={field.key}>
                                        <span className="profile-field-label">
                                            <span>{profileFieldLabel(field, language)}</span>
                                            {locked && (
                                                <em className="lock-chip" title={profileLockedHint(language)}>
                                                    <Lock size={12}/>
                                                    {profileLockedLabel(language)}
                                                </em>
                                            )}
                                            {influenceByField.has(String(field.key)) && (
                                                <button
                                                    className="field-inline-action"
                                                    type="button"
                                                    onClick={(event) => {
                                                        event.preventDefault();
                                                        void openFieldInfluence(String(field.key));
                                                    }}
                                                >
                                                    {terms.tuneInfluenceButton}
                                                </button>
                                            )}
                                        </span>
                                        {field.kind === 'textarea' ? (
                                            <textarea
                                                value={String(value ?? '')}
                                                onChange={(event) => onFieldChange(field.key, event.target.value)}
                                            />
                                        ) : field.kind === 'select' ? (
                                            <select
                                                value={String(value ?? '')}
                                                onChange={(event) => onFieldChange(field.key, event.target.value)}
                                            >
                                                <option value="">--</option>
                                                {tuneUseCaseValues.map(useCase => (
                                                    <option value={useCase} key={useCase}>{terms.useCases[useCase]}</option>
                                                ))}
                                            </select>
                                        ) : (
                                            <input
                                                value={numberInputValue(field, value)}
                                                type={field.kind}
                                                step={field.kind === 'number' ? profileFieldStep(field) : undefined}
                                                readOnly={field.readOnly}
                                                disabled={field.readOnly}
                                                placeholder={locked ? profileLockedHint(language) : undefined}
                                                onFocus={() => onNumberFocus(field, value)}
                                                onBlur={() => onNumberBlur(field)}
                                                onChange={(event) => field.kind === 'number' ? onNumberChange(field, event.target.value) : onFieldChange(field.key, event.target.value)}
                                            />
                                        )}
                                    </label>
                                );
                            })}
                        </div>
                    </div>
                ))}
                    </>
                )}
            </div>
            {snapshotForCompare && (
                <div className="modal-backdrop" role="presentation">
                    <section className="modal-card profile-compare-modal" role="dialog" aria-modal="true" aria-label={terms.profileSnapshotCompare}>
                        <div className="panel-heading">
                            <div>
                                <h2>{terms.profileSnapshotCompare}</h2>
                                <span>{formatTime(snapshotForCompare.changedAt, '--')}</span>
                            </div>
                            <button className="small-action" type="button" onClick={() => setSnapshotForCompare(null)}>{terms.close}</button>
                        </div>
                        <div className="diff-list">
                            {profileDiffs.length === 0 ? <div className="empty-events">{terms.noChanges}</div> : profileDiffs.map(diff => (
                                <div className="diff-row" key={diff.key}>
                                    <span>{diff.label}</span>
                                    <strong title={terms.snapshotBefore}>{diff.left}</strong>
                                    <strong title={terms.currentSettings}>{diff.right}</strong>
                                </div>
                            ))}
                        </div>
                    </section>
                </div>
            )}
            {selectedInfluence && (
                <TuneFieldInfluenceModal
                    influence={selectedInfluence}
                    language={language}
                    terms={terms}
                    onClose={() => setSelectedInfluence(null)}
                />
            )}
        </section>
    );
}

function TuneFieldInfluenceModal({influence, language, terms, onClose}: {
    influence: TuneFieldInfluence;
    language: Lang;
    terms: Copy;
    onClose: () => void;
}) {
    const label = language === 'zh' ? influence.labelZh || influence.labelEn : influence.labelEn || influence.labelZh;
    const summary = language === 'zh' ? influence.summaryZh || influence.summaryEn : influence.summaryEn || influence.summaryZh;
    return (
        <div className="modal-backdrop" role="presentation">
            <section className="modal-card tune-influence-modal" role="dialog" aria-modal="true" aria-label={terms.tuneInfluenceModalTitle}>
                <div className="panel-heading">
                    <div>
                        <h2>{terms.tuneInfluenceModalTitle}</h2>
                        <span>{label}</span>
                    </div>
                    <button className="small-action" type="button" onClick={onClose}>{terms.close}</button>
                </div>
                <p className="tune-influence-summary">{summary}</p>
                <TuneInfluenceDetailGrid influence={influence} terms={terms}/>
            </section>
        </div>
    );
}

function TuneInfluenceDetailGrid({influence, terms}: { influence: TuneFieldInfluence; terms: Copy }) {
    return (
        <div className="tune-influence-detail-grid">
            <TuneInfluenceChipGroup title={terms.tuneInfluenceType} values={[influence.influenceType]} labels={terms.tuneInfluenceTypeLabels}/>
            <TuneInfluenceChipGroup title={terms.tuneInfluenceScope} values={influence.scope} labels={terms.tuneInfluenceScopeLabels}/>
            <TuneInfluenceChipGroup title={terms.tuneInfluencePhase} values={influence.phases} labels={terms.tuneInfluencePhaseLabels}/>
            <TuneInfluenceChipGroup title={terms.tuneInfluenceMetrics} values={influence.tireMetrics} labels={terms.tuneInfluenceMetricLabels}/>
            <TuneInfluenceChipGroup title={terms.tuneInfluenceEvidence} values={influence.evidenceKeys}/>
            <TuneInfluenceChipGroup title={terms.tuneInfluenceSideEffects} values={influence.sideEffects} labels={terms.tuneInfluenceSideEffectLabels}/>
            <TuneInfluenceChipGroup title={terms.tuneInfluenceConditions} values={influence.conditions} labels={terms.tuneInfluenceConditionLabels}/>
        </div>
    );
}

function TuneInfluenceChipGroup({title, values, labels}: {
    title: string;
    values: string[];
    labels?: Record<string, string>;
}) {
    return (
        <div className="tune-influence-chip-group">
            <span>{title}</span>
            <div className="tune-influence-chip-row">
                {(values || []).length === 0 ? (
                    <em>--</em>
                ) : values.map(value => (
                    <strong className="tune-influence-chip" key={value}>{localizedLabel(value, labels) || value}</strong>
                ))}
            </div>
        </div>
    );
}

function ReportsView({
    sessions,
    profiles,
    tuneExplanations,
    selectedSessionId,
    events,
    issueSummary,
    roadTuningDecision,
    tunePlanDraft,
    selectedTunePlanActionIds,
    retestEvaluation,
    samples,
    benchmarkRuns,
    roadEvaluation,
    reportMarkdown,
    status,
    replaySpeed,
    replayStatus,
    compareLeftId,
    compareRightId,
    comparison,
    language,
    terms,
    onSelectSession,
    onGenerateReport,
    onReplay,
    onDeleteSession,
    onBindSession,
    onToggleTunePlanAction,
    onApplyTunePlan,
    onReplaySpeedChange,
    onStopReplay,
    onPauseReplay,
    onResumeReplay,
    onSeekReplay,
    onCompareLeftChange,
    onCompareRightChange,
    onCompare,
    busy,
}: {
    sessions: TelemetrySession[];
    profiles: TuneProfile[];
    tuneExplanations: TuneAdjustmentExplanation[];
    selectedSessionId: number | null;
    events: DetectedEvent[];
    issueSummary: SessionIssueSummary | null;
    roadTuningDecision: RoadTuningDecision | null;
    tunePlanDraft: TunePlanDraft | null;
    selectedTunePlanActionIds: string[];
    retestEvaluation: RetestEvaluation | null;
    samples: TelemetryFrame[];
    benchmarkRuns: BenchmarkRun[];
    roadEvaluation: RoadSessionEvaluation | null;
    reportMarkdown: string;
    status: TelemetryStatus;
    replaySpeed: number;
    replayStatus: TelemetryReplayStatus;
    compareLeftId: number | null;
    compareRightId: number | null;
    comparison: SessionComparison | null;
    language: Lang;
    terms: Copy;
    onSelectSession: (id: number) => void;
    onGenerateReport: (id: number) => void;
    onReplay: (id: number) => void;
    onDeleteSession: (session: TelemetrySession) => void;
    onBindSession: (session: TelemetrySession) => void;
    onToggleTunePlanAction: (actionId: string) => void;
    onApplyTunePlan: (sessionId: number) => void;
    onReplaySpeedChange: (speed: number) => void;
    onStopReplay: () => void;
    onPauseReplay: () => void;
    onResumeReplay: () => void;
    onSeekReplay: (positionMs: number) => void;
    onCompareLeftChange: (id: number | null) => void;
    onCompareRightChange: (id: number | null) => void;
    onCompare: () => void;
    busy: boolean;
}) {
    const selectedSession = sessions.find(session => session.id === selectedSessionId) || null;
    const isReplay = status.mode === 'replay';
    const canReplay = !!selectedSession && selectedSession.recordingPackets > 0 && !!selectedSession.recordingPath && !status.running && !busy;
    const [adviceEventId, setAdviceEventId] = useState('');
    const [adviceGroupId, setAdviceGroupId] = useState('');
    const [expandedSections, setExpandedSections] = useState<Record<string, boolean>>({
        advanced: false,
        replay: false,
        benchmark: false,
        compare: false,
        trends: false,
        markdown: false,
    });
    const adviceEvent = events.find(event => event.id === adviceEventId) || null;
    const adviceGroup = issueSummary?.groups?.find(group => group.id === adviceGroupId) || null;
    const selectedSessionProfile = parseSessionTuneSnapshot(selectedSession) || profiles.find(profile => profile.id === selectedSession?.tuneProfileId) || null;
    const reportAlerts = selectedSession ? reportStatusAlerts(selectedSession, roadEvaluation, terms) : [];
    const toggleSection = (key: string) => setExpandedSections(current => ({...current, [key]: !current[key]}));

    useEffect(() => {
        if (adviceEventId && !events.some(event => event.id === adviceEventId)) {
            setAdviceEventId('');
        }
    }, [events, adviceEventId]);

    useEffect(() => {
        if (adviceGroupId && !(issueSummary?.groups || []).some(group => group.id === adviceGroupId)) {
            setAdviceGroupId('');
        }
    }, [issueSummary, adviceGroupId]);

    return (
        <section className="reports-layout">
            <div className="panel sessions-panel">
                <div className="panel-heading">
                    <h2>{terms.reportSessions}</h2>
                    <span>{sessions.length}</span>
                </div>
                {sessions.length === 0 ? (
                    <div className="empty-events">{terms.noSessions}</div>
                ) : (
                    <div className="session-list">
                        {sessions.map(session => (
                            <button
                                key={session.id}
                                type="button"
                                className={`session-row ${selectedSessionId === session.id ? 'selected' : ''}`}
                                onClick={() => onSelectSession(session.id)}
                            >
                                <strong>{session.sessionName || `#${session.id}`}</strong>
                                <span>{session.tuneName || terms.noProfile}</span>
                                <small>
                                    {formatTime(session.startedAt, '--')} / {terms.eventsSaved}: {session.eventCount} / {terms.samplesSaved}: {session.sampleCount || 0}
                                </small>
                                <small>{terms.sessionMode}: {gameModeLabel(session.gameMode, terms)}</small>
                                <small>{terms.driverMode}: {formatDriverModeDetection(session, terms)}</small>
                                <small>{formatSessionVehicle(session, terms)}</small>
                                <small>{terms.assists}: {formatAssistSummary(session, terms)}</small>
                                <small>{session.recordingPackets > 0 ? `${terms.recordingPackets}: ${session.recordingPackets}` : terms.noRecording}</small>
                            </button>
                        ))}
                    </div>
                )}
            </div>

            <div className="panel report-panel">
                <div className="panel-heading">
                    <h2>{terms.reportDecisionTitle}</h2>
                    <div className="report-actions">
                        <button
                            className="small-action"
                            type="button"
                            disabled={!selectedSession || busy}
                            onClick={() => selectedSession && onBindSession(selectedSession)}
                        >
                            <Pencil size={15}/>
                            {selectedSession?.tuneProfileId ? terms.changeSessionProfile : terms.bindSessionProfile}
                        </button>
                        <button
                            className="small-action"
                            type="button"
                            disabled={!selectedSession || busy}
                            onClick={() => selectedSession && onGenerateReport(selectedSession.id)}
                        >
                            <FileText size={15}/>
                            {terms.generateReport}
                        </button>
                        <button
                            className="small-action danger"
                            type="button"
                            disabled={!selectedSession || busy || status.running}
                            onClick={() => selectedSession && onDeleteSession(selectedSession)}
                            title={terms.deleteSession}
                        >
                            <Trash2 size={15}/>
                            {terms.deleteSession}
                        </button>
                    </div>
                </div>

                {selectedSession && (
                    <div className="report-summary">
                        <Stat label={terms.duration} value={Math.round((selectedSession.durationMs || 0) / 1000)}/>
                        <Stat label={terms.eventsSaved} value={selectedSession.eventCount}/>
                        <Stat label={terms.avgSpeed} value={Math.round(selectedSession.avgSpeedKmh || 0)}/>
                        <Stat label={terms.maxSpeed} value={Math.round(selectedSession.maxSpeedKmh || 0)}/>
                        <Stat label={terms.samplesSaved} value={selectedSession.sampleCount || 0}/>
                        <Stat label={terms.recordingPackets} value={selectedSession.recordingPackets || 0}/>
                        <TextStat label={terms.sessionMode} value={gameModeLabel(selectedSession.gameMode, terms)}/>
                        <TextStat label={terms.driverMode} value={formatDriverModeDetection(selectedSession, terms)}/>
                        <TextStat label={terms.assists} value={formatAssistSummary(selectedSession, terms)}/>
                        <TextStat label={terms.recordingSize} value={formatBytes(selectedSession.recordingBytes || 0)}/>
                        <TextStat label={terms.recording} value={selectedSession.recordingTruncated ? terms.recordingTruncated : selectedSession.recordingPackets > 0 ? terms.recordingReady : terms.noRecording}/>
                    </div>
                )}
                {reportAlerts.length > 0 && (
                    <div className="report-status-list">
                        {reportAlerts.map(alert => (
                            <div className={`status-alert ${alert.tone}`} key={alert.key}>
                                <strong>{alert.title}</strong>
                                <span>{alert.message}</span>
                            </div>
                        ))}
                    </div>
                )}
                {selectedSession && (
                    <RoadEvaluationCard evaluation={roadEvaluation} terms={terms}/>
                )}
                {selectedSession && (
                    <RoadTuningDecisionCard decision={roadTuningDecision} profile={selectedSessionProfile} terms={terms} language={language}/>
                )}
                {selectedSession && (
                    <GearPowerDiagnosticCard summary={issueSummary} terms={terms}/>
                )}
                {selectedSession && (
                    <TunePlanDraftCard
                        session={selectedSession}
                        draft={tunePlanDraft}
                        selectedActionIds={selectedTunePlanActionIds}
                        language={language}
                        terms={terms}
                        busy={busy}
                        status={status}
                        onToggleAction={onToggleTunePlanAction}
                        onApply={onApplyTunePlan}
                    />
                )}
                {selectedSession && (
                    <RetestEvaluationCard evaluation={retestEvaluation} terms={terms} language={language}/>
                )}
                <div className="panel-inner issue-section">
                    <div className="panel-heading compact">
                        <h2>{terms.issueGroups}</h2>
                        <span>{issueSummary?.groups?.length || events.length}</span>
                    </div>
                    {(issueSummary?.groups || []).length > 0 ? (
                        <div className="report-event-list">
                            {(issueSummary?.groups || []).slice(0, 3).map(group => (
                                <button className="report-event-button" key={group.id} type="button" onClick={() => setAdviceGroupId(group.id)}>
                                    <span className="report-event-main">
                                        <strong>{issueFamilyLabel(group.family, terms)}</strong>
                                        <small>{terms.eventsSaved}: {group.eventCount} / {formatDuration(group.totalDurationMs, terms)} / {terms.issueGroupComparison}: {issueComparisonLabel(group.comparison, terms)}</small>
                                        {group.adjustmentStrategy && <small>{terms.issueStrategy}: {localizedLabel(group.adjustmentStrategy, terms.issueStrategyLabels)}</small>}
                                    </span>
                                    <span className={`severity-badge ${group.severity}`}>{severityLabel(group.severity, terms)}</span>
                                </button>
                            ))}
                            {(issueSummary?.groups || []).length > 3 && (
                                <span className="empty-events">{terms.advancedReportDetails}: {(issueSummary?.groups || []).length - 3}</span>
                            )}
                        </div>
                    ) : events.length > 0 ? (
                        <div className="report-event-list">
                        {events.map(event => (
                            <button className="report-event-button" key={event.id} type="button" onClick={() => setAdviceEventId(event.id)}>
                                <span className="report-event-main">
                                    <strong>{eventLabel(event.type, terms)}</strong>
                                    <small>{terms.eventStarted} {formatEventOffset(event.startMs, terms)} / {formatDuration(event.durationMs, terms)}</small>
                                </span>
                                <span className={`severity-badge ${event.severity}`}>{severityLabel(event.severity, terms)}</span>
                            </button>
                        ))}
                        </div>
                    ) : (
                        <div className="empty-events">{terms.noIssueGroups}</div>
                    )}
                </div>

                <CollapsibleSection title={terms.advancedReportDetails} badge={terms.replay} open={expandedSections.advanced} terms={terms} onToggle={() => toggleSection('advanced')}>
                <div className="advanced-report-stack">
                <CollapsibleSection title={terms.playbackAndTimeline} badge={`${formatDuration(replayStatus.positionMs, terms)} / ${formatDuration(replayStatus.durationMs, terms)}`} open={expandedSections.replay} terms={terms} onToggle={() => toggleSection('replay')}>
                    <div className="report-actions advanced-report-actions">
                        <label className="replay-speed">
                            <span>{terms.replaySpeed}</span>
                            <select value={replaySpeed} onChange={(event) => onReplaySpeedChange(Number(event.target.value))} disabled={busy || status.running}>
                                <option value={1}>1x</option>
                                <option value={2}>2x</option>
                                <option value={4}>4x</option>
                            </select>
                        </label>
                        {isReplay ? (
                            <button className="small-action" type="button" onClick={onStopReplay} disabled={busy}>
                                <Square size={15}/>
                                {terms.stopReplay}
                            </button>
                        ) : (
                            <button className="small-action" type="button" disabled={!canReplay} onClick={() => selectedSession && onReplay(selectedSession.id)}>
                                <Power size={15}/>
                                {terms.replay}
                            </button>
                        )}
                    </div>
                    <div className="replay-timeline compact-timeline">
                        <input
                            type="range"
                            min={0}
                            max={Math.max(replayStatus.durationMs, 1)}
                            value={Math.min(replayStatus.positionMs, Math.max(replayStatus.durationMs, 1))}
                            disabled={!replayStatus.running || busy}
                            onChange={(event) => onSeekReplay(Number(event.target.value))}
                        />
                        <div className="timeline-actions">
                            {replayStatus.paused ? (
                                <button className="small-action" type="button" onClick={onResumeReplay} disabled={!replayStatus.running || busy}>{terms.resumeReplay}</button>
                            ) : (
                                <button className="small-action" type="button" onClick={onPauseReplay} disabled={!replayStatus.running || busy}>{terms.pauseReplay}</button>
                            )}
                            <span>{terms.recordingPackets}: {replayStatus.packetIndex + 1}/{replayStatus.packetCount || 0}</span>
                        </div>
                        {events.length > 0 && (
                            <div className="timeline-events">
                                {events.map(event => (
                                    <button key={event.id} type="button" onClick={() => onSeekReplay(event.startMs)} disabled={!replayStatus.running || busy}>
                                        {eventLabel(event.type, terms)} +{formatDuration(event.startMs, terms)}
                                    </button>
                                ))}
                            </div>
                        )}
                    </div>
                </CollapsibleSection>

                <CollapsibleSection title={terms.benchmarkRuns} badge={String(benchmarkRuns.length)} open={expandedSections.benchmark} terms={terms} onToggle={() => toggleSection('benchmark')}>
                    {benchmarkRuns.length === 0 ? (
                        <div className="empty-events">{terms.noBenchmarkRuns}</div>
                    ) : (
                        <div className="report-event-list">
                            {benchmarkRuns.map(run => (
                                <span key={run.id}>
                                    {run.trackName || `#${run.trackId}`} / {formatDuration(run.durationMs, terms)} / {terms.confidence}: {formatPercentValue(run.confidence)} / {terms.routeProgress}: {formatPercentValue(run.routeProgress01 ?? undefined)} / {terms.lengthError}: {formatSignedPercent(run.trackLengthErrorPct)}
                                    {run.warningFlags ? ` / ${terms.warnings}: ${localizedWarningFlags(run.warningFlags, terms).join(' / ')}` : ''}
                                </span>
                            ))}
                        </div>
                    )}
                </CollapsibleSection>

                <CollapsibleSection title={terms.sessionCompare} badge={terms.delta} open={expandedSections.compare} terms={terms} onToggle={() => toggleSection('compare')}>
                    <div className="compare-selects">
                        <label>
                            <span>{terms.leftSession}</span>
                            <select value={compareLeftId || 0} onChange={(event) => onCompareLeftChange(Number(event.target.value) || null)}>
                                <option value={0}>--</option>
                                {sessions.map(session => <option key={session.id} value={session.id}>{session.sessionName || `#${session.id}`}</option>)}
                            </select>
                        </label>
                        <label>
                            <span>{terms.rightSession}</span>
                            <select value={compareRightId || 0} onChange={(event) => onCompareRightChange(Number(event.target.value) || null)}>
                                <option value={0}>--</option>
                                {sessions.map(session => <option key={session.id} value={session.id}>{session.sessionName || `#${session.id}`}</option>)}
                            </select>
                        </label>
                        <button className="small-action" type="button" onClick={onCompare} disabled={!compareLeftId || !compareRightId || busy}>{terms.compare}</button>
                    </div>
                    {comparison && (
                        <>
                            <div className="report-event-list">
                                <span>{terms.left}: {formatTestConditionsCompact(sessionTestConditions(comparison.leftSession), terms)}</span>
                                <span>{terms.right}: {formatTestConditionsCompact(sessionTestConditions(comparison.rightSession), terms)}</span>
                            </div>
                            {(comparison.comparabilityWarnings || []).length > 0 && (
                                <div className="empty-events">
                                    {(comparison.comparabilityWarnings || []).map(warning => (
                                        <span key={warning}>{comparabilityWarningLabel(warning, terms)}</span>
                                    ))}
                                </div>
                            )}
                            <div className="comparison-table">
                                <div className="comparison-head"><span>{terms.metric}</span><span>{terms.left}</span><span>{terms.right}</span><span>{terms.delta}</span></div>
                                {comparison.metrics.map(metric => (
                                    <div className="comparison-row" key={metric.key}>
                                        <span>{metric.label}</span>
                                        <strong>{formatComparisonValue(metric.left, metric.unit)}</strong>
                                        <strong>{formatComparisonValue(metric.right, metric.unit)}</strong>
                                        <strong className={metric.delta === 0 ? '' : metric.delta > 0 === metric.higherIsBetter ? 'good' : 'warn'}>{formatSigned(metric.delta, metric.unit)}</strong>
                                    </div>
                                ))}
                            </div>
                            {comparison.eventTypes.length > 0 && (
                                <div className="report-event-list">
                                    {comparison.eventTypes.map(item => <span key={item.type}>{eventLabel(item.type, terms)}: {item.left} / {item.right} ({formatSigned(item.delta, '')})</span>)}
                                </div>
                            )}
                        </>
                    )}
                </CollapsibleSection>

                <CollapsibleSection title={terms.historicalTrend} badge={samples.length ? String(samples.length) : '--'} open={expandedSections.trends} terms={terms} onToggle={() => toggleSection('trends')}>
                    {samples.length > 0 ? (
                        <div className="history-trends">
                        <h3>{terms.historicalTrend}</h3>
                        <Trend title={terms.speed} unit="km/h" points={sparklinePoints(samples.map(item => item.speedKmh))}/>
                        <Trend title={terms.rpmLoad} unit="%" points={sparklinePoints(samples.map(item => item.rpmRatio * 100))}/>
                        <Trend title={terms.throttle} unit="%" points={sparklinePoints(samples.map(item => item.throttle01 * 100))}/>
                        <Trend title={terms.brake} unit="%" points={sparklinePoints(samples.map(item => item.brake01 * 100))}/>
                        <Trend title={terms.steering} unit="%" points={sparklinePoints(samples.map(item => Math.abs(item.steer01) * 100))}/>
                        <Trend title={terms.frontRearSlip} unit="" points={sparklinePoints(samples.map(item => Math.max(item.frontCombinedSlipAvg, item.rearCombinedSlipAvg)))}/>
                        </div>
                    ) : (
                        <div className="empty-events">{terms.reportPlaceholder}</div>
                    )}
                </CollapsibleSection>

                <CollapsibleSection title={terms.rawMarkdown} badge={reportMarkdown ? terms.saved : '--'} open={expandedSections.markdown} terms={terms} onToggle={() => toggleSection('markdown')}>
                    <pre className="markdown-preview">{reportMarkdown || terms.reportPlaceholder}</pre>
                </CollapsibleSection>
                </div>
                </CollapsibleSection>
            </div>
            {adviceEvent && (
                <ReportEventAdviceModal event={adviceEvent} profile={selectedSessionProfile} tuneExplanations={tuneExplanations} roadEvaluation={roadEvaluation} language={language} terms={terms} onClose={() => setAdviceEventId('')}/>
            )}
            {adviceGroup && (
                <ReportIssueGroupAdviceModal group={adviceGroup} profile={selectedSessionProfile} tuneExplanations={tuneExplanations} language={language} terms={terms} onClose={() => setAdviceGroupId('')}/>
            )}
        </section>
    );
}

function CollapsibleSection({
    title,
    badge,
    open,
    terms,
    onToggle,
    children,
}: {
    title: string;
    badge?: string;
    open: boolean;
    terms: Copy;
    onToggle: () => void;
    children: ReactNode;
}) {
    return (
        <section className={`report-collapsible ${open ? 'open' : ''}`}>
            <button className="report-collapsible-toggle" type="button" onClick={onToggle}>
                <span>{title}</span>
                <strong>{badge || (open ? terms.collapse : terms.expand)}</strong>
                <small>{open ? terms.collapse : terms.expand}</small>
            </button>
            {open && <div className="report-collapsible-body">{children}</div>}
        </section>
    );
}

function RoadEvaluationCard({evaluation, terms}: { evaluation: RoadSessionEvaluation | null; terms: Copy }) {
    if (!evaluation) {
        return (
            <div className="panel-inner road-evaluation-card">
                <div className="panel-heading compact">
                    <h2>{terms.roadEvaluation}</h2>
                    <span>{terms.insufficientData}</span>
                </div>
                <div className="empty-events">{terms.roadEvaluationEmpty}</div>
            </div>
        );
    }
    const bestRun = evaluation.bestRun || null;
    const baselineRun = evaluation.baselineRun || null;
    return (
        <div className="panel-inner road-evaluation-card">
            <div className="panel-heading compact">
                <h2>{terms.roadEvaluation}</h2>
                <span>{roadVerdictLabel(evaluation.overallVerdict, terms)}</span>
            </div>
            <div className="road-score-grid">
                <TextStat label={terms.paperPerformanceScore} value={formatScore(evaluation.paperPerformanceScore)}/>
                <TextStat label={terms.playerFitScore} value={formatScore(evaluation.playerFitScore)}/>
                <TextStat label={terms.riskScore} value={formatScore(evaluation.riskScore)}/>
                <TextStat label={terms.autoBaseline} value={roadBaselineStatusLabel(evaluation.baselineStatus, terms)}/>
            </div>
            <div className="road-run-summary">
                <div>
                    <span>{terms.bestPlayerRun}</span>
                    <strong>{bestRun ? `${bestRun.trackName || `#${bestRun.trackId}`} / ${formatDuration(bestRun.durationMs, terms)} / ${terms.confidence}: ${formatPercentValue(bestRun.confidence)}` : '--'}</strong>
                </div>
                <div>
                    <span>{terms.autoBaseline}</span>
                    <strong>{baselineRun ? `${baselineRun.trackName || `#${baselineRun.trackId}`} / ${formatDuration(baselineRun.durationMs, terms)}` : roadBaselineStatusLabel(evaluation.baselineStatus, terms)}</strong>
                </div>
            </div>
            {evaluation.baselineStatus === 'missing_auto_baseline' && (
                <div className="empty-events advice-placeholder">{terms.missingAutoBaselineHint}</div>
            )}
            {(evaluation.attributions || []).length > 0 && (
                <div className="attribution-list">
                    {(evaluation.attributions || []).map((item, index) => (
                        <span key={`${item.type}-${item.eventType || item.message}-${index}`}>
                            {roadAttributionLabel(item.type, terms)}
                            {item.eventType ? ` / ${eventLabel(item.eventType, terms)}` : ''}
                            {item.count ? ` x${item.count}` : ''}
                            {item.prioritizeTuning ? ` / ${terms.prioritizeTuning}` : ''}
                        </span>
                    ))}
                </div>
            )}
        </div>
    );
}

function GearPowerDiagnosticCard({summary, terms}: { summary: SessionIssueSummary | null; terms: Copy }) {
    const gearPower = summary?.gearPower;
    const gearRows = (gearPower?.gears || []).filter(gear => gear.sampleCount > 0).slice(0, 10);
    const actions = gearPower?.recommendedActions || [];
    const targetMin = gearPower?.powerBandStartRPM || 0;
    const targetMax = gearPower?.powerBandEndRPM || 0;
    const targetLabel = targetMin > 0 && targetMax > 0
        ? `${formatNumber(targetMin, 0)} - ${formatNumber(targetMax, 0)} rpm`
        : `${formatPercentValue(gearPower?.evidence?.power_band_target_min ?? 0.55)} - ${formatPercentValue(gearPower?.evidence?.power_band_target_max ?? 0.9)}`;
    return (
        <div className="panel-inner issue-section gear-power-card">
            <div className="panel-heading compact">
                <h2>{terms.gearPowerDiagnostic}</h2>
                <span>{localizedLabel(gearPower?.summary || '', terms.gearFindings) || localizedLabel(gearPower?.summary || '', terms.planSummaries) || terms.insufficientData}</span>
            </div>
            {!gearPower || gearPower.status === 'insufficient_data' ? (
                <div className="empty-events">{terms.reportPlaceholder}</div>
            ) : (
                <>
                    <div className="report-summary">
                        <TextStat label={terms.gearPowerSummary} value={localizedLabel(gearPower.summary || '', terms.gearFindings) || localizedLabel(gearPower.summary || '', terms.planSummaries)}/>
                        <TextStat label={terms.gearFinding} value={[gearPower.launchFinding, gearPower.topSpeedFinding].filter(Boolean).map(item => localizedLabel(item || '', terms.gearFindings)).join(' / ') || '--'}/>
                        <TextStat label={terms.powerBandTarget} value={targetLabel}/>
                        <TextStat label={terms.powerBandSource} value={localizedLabel(gearPower.powerBandSource || '', terms.powerBandSources) || '--'}/>
                        <TextStat label={terms.diagnosticConfidence} value={localizedLabel(gearPower.confidence || '', terms.planConfidenceLabels) || '--'}/>
                        <TextStat label={terms.gearStrategyMode} value={localizedLabel(gearPower.strategyMode || '', terms.gearFindings) || '--'}/>
                        <TextStat label={terms.gearStrategyIssueCount} value={formatGearStrategyIssueCount(gearPower)}/>
                        <TextStat label={terms.powerToWeight} value={gearPower.powerToWeightKWPerKG ? `${gearPower.powerToWeightKWPerKG.toFixed(4)} kW/kg` : '--'}/>
                        <TextStat label={terms.tractionLimited} value={formatPercentValue(gearPower.tractionLimitedPercent || 0)}/>
                    </div>
                    {gearRows.length > 0 ? (
                        <div className="diagnostic-table compact-diagnostic">
                            <div className="diagnostic-row diagnostic-head">
                                <span>{terms.gear}</span>
                                <span>{terms.speedRange}</span>
                                <span>RPM</span>
                                <span>{terms.inPowerBand}</span>
                                <span>{terms.acceleration}</span>
                                <span>{terms.shiftAfter}</span>
                                <span>{terms.tractionLimited}</span>
                                <span>{terms.gearFinding}</span>
                                <span>{terms.samplesSaved}</span>
                            </div>
                            {gearRows.map(gear => (
                                <div className={`diagnostic-row ${gear.finding !== 'ok' ? 'warn' : ''}`} key={gear.gear}>
                                    <span>{gear.gear}</span>
                                    <GearSpeedValue gear={gear} terms={terms}/>
                                    <span>{gear.rpmAvg ? formatNumber(gear.rpmAvg, 0) : formatPercentValue(gear.rpmRatioAvg)}</span>
                                    <span>{formatGearInPowerBandRange(gear)}</span>
                                    <span>{formatNumber(gear.accelAvgMps2 || 0, 2)} m/s²</span>
                                    <span>{gear.shiftAfterRPM ? formatNumber(gear.shiftAfterRPM, 0) : '--'}</span>
                                    <span>{formatPercentValue(gear.tractionLimitedPercent || 0)}</span>
                                    <span>{localizedLabel(gear.finding, terms.gearFindings)}</span>
                                    <span>{gear.sampleCount}</span>
                                </div>
                            ))}
                        </div>
                    ) : (
                        <div className="empty-events">{terms.reportPlaceholder}</div>
                    )}
                    {actions.length > 0 && (
                        <div className="report-event-list">
                            {actions.map((action, index) => (
                                <span key={`${action.category}-${action.item}-${action.direction}-${index}`}>
                                    {localizedLabel(action.category, terms.actionCategories)} / {localizedLabel(action.item, terms.actionItems)}:
                                    {' '}{localizedLabel(action.direction, terms.actionDirections)} {localizedLabel(action.amount, terms.actionAmounts) || action.amount}
                                </span>
                            ))}
                        </div>
                    )}
                    <GearPowerComparisonList comparisons={gearPower.comparisons || []} terms={terms}/>
                </>
            )}
        </div>
    );
}

function GearPowerComparisonList({comparisons, terms}: { comparisons: GearPowerComparison[]; terms: Copy }) {
    if (!comparisons || comparisons.length === 0) {
        return null;
    }
    const telemetryComparison = comparisons.find(item => item.type === 'session_telemetry');
    const tuneComparison = comparisons.find(item => item.type === 'tune_settings');
    return (
        <div className="event-section gear-comparison-section">
            <h3>{terms.gearPowerComparisons}</h3>
            {telemetryComparison && (
                <GearTelemetryComparison comparison={telemetryComparison} terms={terms}/>
            )}
            {tuneComparison && (
                <GearTuneComparison comparison={tuneComparison} terms={terms}/>
            )}
        </div>
    );
}

function GearTelemetryComparison({comparison, terms}: { comparison: GearPowerComparison; terms: Copy }) {
    const rows = comparison.rows || [];
    return (
        <div className="gear-comparison-card">
            <div className="section-title-row">
                <div>
                    <h3>{terms.gearTelemetryComparison}</h3>
                    <span>{localizedLabel(comparison.status || '', terms.gearComparisonStatuses)}</span>
                </div>
            </div>
            {comparison.status !== 'ready' || rows.length === 0 ? (
                <div className="empty-events advice-placeholder">{localizedLabel(comparison.status || '', terms.gearComparisonStatuses) || terms.gearComparisonUnavailable}</div>
            ) : (
                <div className="diagnostic-table compact-diagnostic gear-comparison-table telemetry-comparison-table">
                    <div className="diagnostic-row diagnostic-head">
                        <span>{terms.gear}</span>
                        <span>{terms.highestObserved}</span>
                        <span>{terms.inPowerBandCoverage}</span>
                        <span>{terms.tractionLimited}</span>
                        <span>{terms.gearFinding}</span>
                    </div>
                    {rows.map(row => (
                        <div className="diagnostic-row" key={`gear-telemetry-${row.item}-${row.gear}`}>
                            <span>{localizedLabel(row.item, terms.actionItems)}</span>
                            <strong>{formatGearTelemetryDelta(row.beforeSpeedMaxKmh, row.afterSpeedMaxKmh, row.speedMaxDeltaKmh, 'km/h')}</strong>
                            <span>{formatGearTelemetryDelta(row.beforeInPowerBandPct, row.afterInPowerBandPct, row.inPowerBandDeltaPct, '%')}</span>
                            <span>{formatGearTelemetryDelta(row.beforeTractionLimitPct, row.afterTractionLimitPct, row.tractionLimitDeltaPct, '%')}</span>
                            <span>{localizedLabel(row.beforeFinding || 'ok', terms.gearFindings)} {'->'} {localizedLabel(row.afterFinding || 'ok', terms.gearFindings)}</span>
                        </div>
                    ))}
                </div>
            )}
        </div>
    );
}

function GearTuneComparison({comparison, terms}: { comparison: GearPowerComparison; terms: Copy }) {
    const rows = comparison.rows || [];
    return (
        <div className="gear-comparison-card">
            <div className="section-title-row">
                <div>
                    <h3>{terms.gearTuneComparison}</h3>
                    <span>{localizedLabel(comparison.status || '', terms.gearComparisonStatuses)}</span>
                </div>
            </div>
            {comparison.status !== 'ready' || rows.length === 0 ? (
                <div className="empty-events advice-placeholder">{localizedLabel(comparison.status || '', terms.gearComparisonStatuses) || terms.gearComparisonUnavailable}</div>
            ) : (
                <div className="diagnostic-table compact-diagnostic gear-comparison-table tune-comparison-table">
                    <div className="diagnostic-row diagnostic-head">
                        <span>{terms.gear}</span>
                        <span>{terms.gearComparisonBefore}</span>
                        <span>{terms.gearComparisonAfter}</span>
                        <span>{terms.gearComparisonDelta}</span>
                    </div>
                    {rows.map(row => (
                        <div className="diagnostic-row" key={`gear-tune-${row.item}-${row.gear}`}>
                            <span>{localizedLabel(row.item, terms.actionItems)}</span>
                            <strong>{formatNumber(row.beforeValue, 2)}</strong>
                            <span>{formatNumber(row.afterValue, 2)}</span>
                            <span>{formatSignedNumber(row.deltaValue, 2)}</span>
                        </div>
                    ))}
                </div>
            )}
        </div>
    );
}

function WholeCarPlanCard({summary, profile, terms, language}: { summary: SessionIssueSummary | null; profile: TuneProfile | null; terms: Copy; language: Lang }) {
    const plan = summary?.wholeCarPlan;
    const gearPower = summary?.gearPower;
    const actions = plan?.actions || [];
    const gearRows = (gearPower?.gears || []).filter(gear => gear.sampleCount > 0).slice(0, 10);
    return (
        <div className="panel-inner issue-section">
            <div className="panel-heading compact">
                <h2>{terms.wholeCarPlan}</h2>
                <span>{localizedLabel(plan?.strategy || '', terms.planStrategies)}</span>
            </div>
            {!plan || actions.length === 0 ? (
                <div className="empty-events">{terms.wholeCarPlanEmpty}</div>
            ) : (
                <>
                    <div className="report-summary">
                        <TextStat label={terms.planStrategy} value={localizedLabel(plan.strategy, terms.planStrategies)}/>
                        <TextStat label={terms.planConfidence} value={localizedLabel(plan.confidence, terms.planConfidenceLabels)}/>
                        <TextStat label={terms.gearPowerSummary} value={localizedLabel(gearPower?.summary || '', terms.gearFindings) || localizedLabel(gearPower?.summary || '', terms.planSummaries)}/>
                        <TextStat label={terms.gearFinding} value={[gearPower?.launchFinding, gearPower?.topSpeedFinding].filter(Boolean).map(item => localizedLabel(item || '', terms.gearFindings)).join(' / ') || '--'}/>
                        <TextStat label={terms.powerToWeight} value={gearPower?.powerToWeightKWPerKG ? `${gearPower.powerToWeightKWPerKG.toFixed(4)} kW/kg` : '--'}/>
                        <TextStat label={terms.tractionLimited} value={formatPercentValue(gearPower?.tractionLimitedPercent || 0)}/>
                    </div>
                    <div className="suggestion-list">
                        {actions.map(action => (
                            <div key={`${action.priority}-${action.category}-${action.item}-${action.direction}`} className="suggestion">
                                <strong>{localizedLabel(action.category, terms.actionCategories)} / {localizedLabel(action.item, terms.actionItems)}</strong>
                                <span>{formatConcreteSuggestion(action, wholeCarActionEvent(action), profile, language, terms)}</span>
                                <small>{issueFamilyLabel(action.family, terms)} / {localizedLabel(action.reason, terms.actionReasons)}</small>
                            </div>
                        ))}
                    </div>
                    {(plan.conflicts || []).length > 0 && (
                        <div className="event-section">
                            <h3>{terms.planConflicts}</h3>
                            <div className="report-event-list">
                                {plan.conflicts.map(conflict => (
                                    <span key={`${conflict.key}-${conflict.keptItem}`}>{conflict.key}: {conflict.keptItem} / {conflict.droppedItem}</span>
                                ))}
                            </div>
                        </div>
                    )}
                    {gearRows.length > 0 && (
                        <div className="event-section">
                            <h3>{terms.gearPowerDiagnostic}</h3>
                            <div className="diagnostic-table compact-diagnostic">
                                <div className="diagnostic-row diagnostic-head">
                                    <span>{terms.gear}</span>
                                    <span>{terms.speedRange}</span>
                                    <span>{terms.rpmLoad}</span>
                                    <span>{terms.throttle}</span>
                                    <span>{terms.tractionLimited}</span>
                                    <span>{terms.gearFinding}</span>
                                    <span>{terms.samplesSaved}</span>
                                </div>
                                {gearRows.map(gear => (
                                    <div className={`diagnostic-row ${gear.finding !== 'ok' ? 'warn' : ''}`} key={gear.gear}>
                                        <span>{gear.gear}</span>
                                        <GearSpeedValue gear={gear} terms={terms}/>
                                        <span>{formatPercentValue(gear.rpmRatioAvg)}</span>
                                        <span>{formatPercentValue(gear.throttleAvg)}</span>
                                        <span>{formatPercentValue(gear.tractionLimitedPercent || 0)}</span>
                                        <span>{localizedLabel(gear.finding, terms.gearFindings)}</span>
                                        <span>{gear.sampleCount}</span>
                                    </div>
                                ))}
                            </div>
                        </div>
                    )}
                </>
            )}
        </div>
    );
}

function RoadTuningDecisionCard({decision, profile, terms, language}: { decision: RoadTuningDecision | null; profile: TuneProfile | null; terms: Copy; language: Lang }) {
    if (!decision || !decision.status || ['no_matching_symptom', 'insufficient_data', 'knowledge_error'].includes(decision.status)) {
        return (
            <div className="panel-inner whole-car-plan">
                <div className="panel-heading compact">
                    <h2>{terms.roadTuningDecision}</h2>
                    <span>{decision ? localizedLabel(decision.status, terms.roadDecisionStatusLabels) : '--'}</span>
                </div>
                <div className="empty-events">{decision?.knowledgeStatus?.lastError || terms.roadTuningDecisionEmpty}</div>
            </div>
        );
    }
    return (
        <div className="panel-inner whole-car-plan road-decision-card">
            <div className="panel-heading compact">
                <h2>{terms.roadTuningDecision}</h2>
                <span>{localizedLabel(decision.confidence || 'medium', terms.planConfidenceLabels)}</span>
            </div>
            <div className="decision-summary-grid">
                <TextStat label={terms.primaryIssue} value={decision.symptom || localizedLabel(decision.symptomId, terms.roadSymptomLabels) || decision.symptomId}/>
                <TextStat label={terms.cornerPhase} value={localizedLabel(decision.phase, terms.roadPhaseLabels) || decision.phase || '--'}/>
                <TextStat label={terms.primaryCause} value={decision.primaryCause || decision.reason || '--'}/>
                <TextStat label={terms.driverFitVerdict} value={localizedLabel(decision.fitVerdict, terms.driverFitVerdictLabels)}/>
            </div>
            <div className="evidence-grid">
                {decisionEvidenceEntries(decision).map(([key, value]) => (
                    <div key={key}>
                        <span>{localizedLabel(key, terms.evidenceLabels)}</span>
                        <strong>{formatDecisionEvidenceValue(key, value, language)}</strong>
                    </div>
                ))}
            </div>
            {decision.rollbackRecommended && (
                <div className="status-alert warn">
                    <strong>{terms.rollbackRecommended}</strong>
                    <span>{terms.rollbackRecommendedHint}</span>
                </div>
            )}
            {profile == null && <div className="empty-events">{terms.concreteProfileRequired}</div>}
            {decision.actions.length === 0 ? (
                <div className="empty-events">{terms.roadTuningDecisionEmpty}</div>
            ) : (
                <div className="plan-action-list">
                    {decision.actions.map(action => {
                        const draftPreview = profile && action.canAutoApply ? previewDecisionAction(action, profile, language) : null;
                        return (
                            <div className="plan-action-card" key={action.id}>
                                <div>
                                    <strong>{localizedLabel(action.role, terms.roadActionRoles)} · {localizedLabel(action.item, terms.actionItems)}</strong>
                                    <span>{draftPreview || `${localizedLabel(action.direction, terms.actionDirections)} ${action.amount}${action.unit ? ` ${action.unit}` : ''}`}</span>
                                    <small>
                                        {localizedLabel(action.adviceLayer || 'primary', terms.adviceLayerLabels)}
                                        {' / '}
                                        {localizedLabel(action.trustLevel || action.confidence || 'medium', terms.tunePlanTrustLevels)}
                                    </small>
                                    <small>{localizedLabel(action.rationale || action.reason, terms.actionReasons) || action.rationale || action.reason || '--'}</small>
                                    {action.conflictReason && <small className="warn-text">{localizedLabel(action.conflictReason, terms.tunePlanBlockedReasons) || action.conflictReason}</small>}
                                    {!action.canAutoApply && <small>{localizedLabel(action.blockedReason || 'manual_review_required', terms.tunePlanBlockedReasons)}</small>}
                                </div>
                                <span className={`severity-badge ${action.canAutoApply ? 'medium' : 'low'}`}>{action.canAutoApply ? terms.autoApplicable : terms.manualCheck}</span>
                            </div>
                        );
                    })}
                </div>
            )}
            {decision.retestFocus.length > 0 && (
                <div className="report-event-list">
                    <span>{terms.retestFocus}: {decision.retestFocus.map(item => localizedLabel(item, terms.retestFocusLabels)).join(' / ')}</span>
                </div>
            )}
            {decision.knowledgeStatus && (
                <div className="report-event-list">
                    <span>{terms.knowledgeSource}: {decision.knowledgeStatus.symptomCount} / {decision.knowledgeStatus.actionCount}{decision.knowledgeStatus.usingFallback ? ` · ${terms.knowledgeFallback}` : ''}</span>
                </div>
            )}
        </div>
    );
}

function TunePlanDraftCard({
    session,
    draft,
    selectedActionIds,
    language,
    terms,
    busy,
    status,
    onToggleAction,
    onApply,
}: {
    session: TelemetrySession;
    draft: TunePlanDraft | null;
    selectedActionIds: string[];
    language: Lang;
    terms: Copy;
    busy: boolean;
    status: TelemetryStatus;
    onToggleAction: (actionId: string) => void;
    onApply: (sessionId: number) => void;
}) {
    const actions = draft?.actions || [];
    const canApply = !!draft && draft.status === 'ready' && selectedActionIds.length > 0 && !busy && !status.running;
    return (
        <div className="panel-inner issue-section tune-plan-draft-card">
            <div className="panel-heading compact">
                <h2>{terms.tunePlanDraft}</h2>
                <span>{localizedLabel(draft?.status || 'no_actions', terms.tunePlanStatusLabels)}</span>
            </div>
            {!draft || actions.length === 0 ? (
                <div className="empty-events">{terms.tunePlanDraftEmpty}</div>
            ) : (
                <>
                    {draft.summary && <div className="empty-events advice-placeholder">{localizedLabel(draft.summary, terms.planSummaries)}</div>}
                    <div className="tune-plan-actions">
                        {actions.map(action => {
                            const checked = selectedActionIds.includes(action.id);
                            const blocked = !action.canApply;
                            return (
                                <label className={`tune-plan-action ${blocked ? 'blocked' : ''}`} key={action.id}>
                                    <input
                                        type="checkbox"
                                        checked={checked}
                                        disabled={blocked || busy || status.running}
                                        onChange={() => onToggleAction(action.id)}
                                    />
                                    <span className="tune-plan-action-main">
                                        <strong>{profileFieldDisplayName(action.fieldKey, language)}</strong>
                                        <span className="tune-plan-values">
                                            {formatDraftValue(action.currentValue, action.step, action.unit)}
                                            <span aria-hidden="true">-&gt;</span>
                                            {formatDraftValue(action.targetValue, action.step, action.unit)}
                                            {action.delta !== undefined && <em>{formatDraftDelta(action.delta, action.step, action.unit)}</em>}
                                        </span>
                                        <small>{issueFamilyLabel(action.family, terms)} / {localizedLabel(action.reason, terms.actionReasons)}</small>
                                        <small>
                                            {localizedLabel(action.adviceLayer || 'primary', terms.adviceLayerLabels)}
                                            {' / '}
                                            {terms.tunePlanTrust}: {localizedLabel(action.trustLevel || action.confidence || 'medium', terms.tunePlanTrustLevels)}
                                            {action.retestGuard && ` / ${terms.tunePlanRetestGuard}: ${localizedLabel(action.retestGuard, terms.tunePlanBlockedReasons) || action.retestGuard}`}
                                        </small>
                                        {action.rationale && action.rationale !== action.reason && (
                                            <small>{localizedLabel(action.rationale, terms.actionReasons) || action.rationale}</small>
                                        )}
                                        {(action.trustReasons || []).length > 0 && (
                                            <small>{(action.trustReasons || []).map(reason => localizedLabel(reason, terms.tunePlanTrustReasons) || reason).join(' / ')}</small>
                                        )}
                                        {(action.missingInputs || []).length > 0 && (
                                            <small className="warn-text">{terms.tunePlanMissingInputs}: {(action.missingInputs || []).map(item => localizedLabel(item, terms.tunePlanMissingInputLabels) || item).join(' / ')}</small>
                                        )}
                                        {blocked && <small className="warn-text">{localizedLabel(action.blockedReason, terms.tunePlanBlockedReasons)}</small>}
                                        {action.conflictReason && <small className="warn-text">{localizedLabel(action.conflictReason, terms.tunePlanBlockedReasons) || action.conflictReason}</small>}
                                    </span>
                                </label>
                            );
                        })}
                    </div>
                    {(draft.conflicts || []).length > 0 && (
                        <div className="event-section">
                            <h3>{terms.planConflicts}</h3>
                            <div className="report-event-list">
                                {draft.conflicts.map(conflict => (
                                    <span key={`${conflict.key}-${conflict.keptItem}-${conflict.droppedItem}`}>{conflict.key}: {conflict.keptItem} / {conflict.droppedItem}</span>
                                ))}
                            </div>
                        </div>
                    )}
                    <div className="draft-actions-footer">
                        <button className="small-action" type="button" disabled={!canApply} onClick={() => onApply(session.id)}>
                            <Save size={15}/>
                            {terms.applyTunePlan}
                        </button>
                    </div>
                </>
            )}
        </div>
    );
}

function RetestEvaluationCard({evaluation, terms, language}: { evaluation: RetestEvaluation | null; terms: Copy; language: Lang }) {
    const metrics = evaluation?.metrics || [];
    const rollbackActions = evaluation?.rollbackActions || [];
    return (
        <div className="panel-inner issue-section retest-evaluation-card">
            <div className="panel-heading compact">
                <h2>{terms.retestResult}</h2>
                <span>{localizedLabel(evaluation?.status || 'insufficient_data', terms.retestStatusLabels)}</span>
            </div>
            {!evaluation || metrics.length === 0 ? (
                <div className="empty-events">{terms.retestEmpty}</div>
            ) : (
                <>
                    {evaluation.baselineSession && (
                        <div className="empty-events advice-placeholder">
                            {terms.left}: {evaluation.baselineSession.sessionName || `#${evaluation.baselineSession.id}`}
                        </div>
                    )}
                    <div className="report-summary">
                        <TextStat label={terms.retestConfidence} value={localizedLabel(evaluation.confidence || 'low', terms.planConfidenceLabels) || '--'}/>
                        <TextStat label={terms.retestBaselineReason} value={localizedLabel(evaluation.baselineReason || '', terms.retestBaselineReasons) || '--'}/>
                        <TextStat label={terms.retestChangedFields} value={(evaluation.changedFields || []).map(field => profileFieldDisplayName(field, language)).join(' / ') || '--'}/>
                    </div>
                    {(evaluation.metricSummary || []).length > 0 && (
                        <div className="report-event-list">
                            {(evaluation.metricSummary || []).map(item => <span key={item}>{formatRetestMetricSummary(item, terms)}</span>)}
                        </div>
                    )}
                    {rollbackActions.length > 0 && (
                        <div className="event-section">
                            <h3>{terms.retestRollbackActions}</h3>
                            <div className="report-event-list">
                                {rollbackActions.map(action => (
                                    <span key={action.id}>
                                        {profileFieldDisplayName(action.fieldKey, language)}: {formatDraftValue(action.currentValue, action.step, action.unit)} -&gt; {formatDraftValue(action.targetValue, action.step, action.unit)}
                                    </span>
                                ))}
                            </div>
                        </div>
                    )}
                    <div className="retest-metric-list">
                        {metrics.map(metric => (
                            <div className="retest-metric-row" key={metric.key}>
                                <span>{localizedLabel(metric.key, terms.retestMetricLabels)}</span>
                                <strong>{formatRetestMetricValue(metric.key, metric.baseline, terms)} -&gt; {formatRetestMetricValue(metric.key, metric.current, terms)}</strong>
                                <em className={metric.status === 'improved' ? 'good' : metric.status === 'worsened' ? 'warn' : ''}>
                                    {formatRetestMetricDelta(metric.key, metric.delta, terms)} / {localizedLabel(metric.status, terms.retestStatusLabels)}
                                </em>
                            </div>
                        ))}
                    </div>
                </>
            )}
        </div>
    );
}

function wholeCarActionEvent(action: WholeCarAdjustment): DetectedEvent {
    return {
        id: `${action.source}-${action.priority}`,
        type: action.family,
        severity: action.confidence === 'high' ? 'high' : 'medium',
        startMs: 0,
        endMs: 0,
        durationMs: 0,
        segment: '',
        evidence: action.evidence || {},
        suggestedActions: [action],
    };
}

function ReportEventAdviceModal({event, profile, tuneExplanations, roadEvaluation, language, terms, onClose}: { event: DetectedEvent; profile: TuneProfile | null; tuneExplanations: TuneAdjustmentExplanation[]; roadEvaluation: RoadSessionEvaluation | null; language: Lang; terms: Copy; onClose: () => void }) {
    const suggestions = event.suggestedActions || [];
    const attribution = (roadEvaluation?.attributions || []).find(item => item.eventType === event.type) || null;
    return (
        <div className="modal-backdrop" role="presentation">
            <section className="modal-card report-event-advice-modal" role="dialog" aria-modal="true" aria-label={terms.eventAdviceTitle}>
                <div className="panel-heading">
                    <div>
                        <h2>{terms.eventAdviceTitle}</h2>
                        <span>{eventLabel(event.type, terms)}</span>
                    </div>
                    <button className="small-action" type="button" onClick={onClose}>{terms.close}</button>
                </div>
                <div className="event-details">
                    <div className="event-details-head">
                        <div>
                            <span>{terms.eventSegment}</span>
                            <strong>{localizedLabel(event.segment, terms.segments)}</strong>
                        </div>
                        <div>
                            <span>{terms.eventDuration}</span>
                            <strong>{formatDuration(event.durationMs, terms)}</strong>
                        </div>
                        <div>
                            <span>{terms.severityLabel}</span>
                            <strong>{severityLabel(event.severity, terms)}</strong>
                        </div>
                        <div>
                            <span>{terms.eventStarted}</span>
                            <strong>{formatEventOffset(event.startMs, terms)}</strong>
                        </div>
                    </div>

                    <div className="event-section">
                        <h3>{terms.eventEvidence}</h3>
                        <div className="evidence-grid">
                            {Object.entries(event.evidence || {}).map(([key, value]) => (
                                <div key={key}>
                                    <span>{localizedLabel(key, terms.evidenceLabels)}</span>
                                    <strong>{formatEvidenceValueForKey(key, value, terms)}</strong>
                                </div>
                            ))}
                        </div>
                    </div>

                    <div className="event-section">
                        <h3>{terms.evaluationContext}</h3>
                        {attribution ? (
                            <div className="evaluation-context">
                                <span>{roadAttributionLabel(attribution.type, terms)}</span>
                                <strong>{attribution.prioritizeTuning ? terms.prioritizeTuningYes : terms.prioritizeTuningNo}</strong>
                                <small>{roadAttributionMessage(attribution, terms)}</small>
                            </div>
                        ) : (
                            <div className="empty-events advice-placeholder">{terms.evaluationContextEmpty}</div>
                        )}
                    </div>

                    <div className="event-section">
                        <h3>{terms.eventSuggestions}</h3>
                        {suggestions.length === 0 ? (
                            <div className="empty-events advice-placeholder">{terms.advicePlaceholder}</div>
                        ) : (
                            <div className="suggestion-list">
                                {suggestions.map(action => (
                                    <div key={`${action.priority}-${action.category}-${action.item}`} className="suggestion">
                                        <strong>{localizedLabel(action.category, terms.actionCategories)} / {localizedLabel(action.item, terms.actionItems)}</strong>
                                        <span>{formatConcreteSuggestion(action, event, profile, language, terms)}</span>
                                        <small>{localizedLabel(action.reason, terms.actionReasons)}</small>
                                        {formatTuneExplanationNote(action, event, tuneExplanations) && (
                                            <small><b>{terms.tuningNote}:</b> {formatTuneExplanationNote(action, event, tuneExplanations)}</small>
                                        )}
                                    </div>
                                ))}
                            </div>
                        )}
                    </div>
                </div>
            </section>
        </div>
    );
}

function ReportIssueGroupAdviceModal({group, profile, tuneExplanations, language, terms, onClose}: { group: SessionIssueGroup; profile: TuneProfile | null; tuneExplanations: TuneAdjustmentExplanation[]; language: Lang; terms: Copy; onClose: () => void }) {
    const representativeEvent = issueGroupRepresentativeEvent(group);
    return (
        <div className="modal-backdrop" role="presentation">
            <section className="modal-card report-event-advice-modal" role="dialog" aria-modal="true" aria-label={terms.issueGroupAdviceTitle}>
                <div className="panel-heading">
                    <div>
                        <h2>{terms.issueGroupAdviceTitle}</h2>
                        <span>{issueFamilyLabel(group.family, terms)}</span>
                    </div>
                    <button className="small-action" type="button" onClick={onClose}>{terms.close}</button>
                </div>
                <div className="event-details">
                    <div className="event-details-head">
                        <div>
                            <span>{terms.eventsSaved}</span>
                            <strong>{group.eventCount}</strong>
                        </div>
                        <div>
                            <span>{terms.eventDuration}</span>
                            <strong>{formatDuration(group.totalDurationMs, terms)}</strong>
                        </div>
                        <div>
                            <span>{terms.severityLabel}</span>
                            <strong>{severityLabel(group.severity, terms)}</strong>
                        </div>
                        <div>
                            <span>{terms.issueGroupComparison}</span>
                            <strong>{issueComparisonLabel(group.comparison, terms)}</strong>
                        </div>
                    </div>

                    {(group.adjustmentStrategy || group.feedbackDirective) && (
                        <div className="event-section">
                            <h3>{terms.issueStrategy}</h3>
                            <div className="empty-events advice-placeholder">
                                {group.adjustmentStrategy && <span>{localizedLabel(group.adjustmentStrategy, terms.issueStrategyLabels)}</span>}
                                {group.feedbackDirective && <span>{localizedLabel(group.feedbackDirective, terms.feedbackDirectiveLabels)}</span>}
                            </div>
                        </div>
                    )}

                    {(group.relatedRecentChanges || []).length > 0 && (
                        <div className="event-section">
                            <h3>{terms.issueRecentChanges}</h3>
                            <div className="attribution-list">
                                {group.relatedRecentChanges.map(field => <span key={field}>{profileFieldDisplayName(field, language)}</span>)}
                            </div>
                        </div>
                    )}

                    <div className="event-section">
                        <h3>{terms.issueGroupPrimaryAdvice}</h3>
                        {!profile ? (
                            <div className="empty-events advice-placeholder">{terms.concreteProfileRequired}</div>
                        ) : (group.primaryActions || []).length === 0 ? (
                            <div className="empty-events advice-placeholder">{terms.advicePlaceholder}</div>
                        ) : (
                            <div className="suggestion-list">
                                {(group.primaryActions || []).map(action => (
                                    <div key={`${action.priority}-${action.category}-${action.item}`} className="suggestion">
                                        <strong>{localizedLabel(action.category, terms.actionCategories)} / {localizedLabel(action.item, terms.actionItems)}</strong>
                                        <span>{formatConcreteSuggestion(action, representativeEvent, profile, language, terms)}</span>
                                        <small>{localizedLabel(action.reason, terms.actionReasons)}</small>
                                        {formatTuneExplanationNote(action, representativeEvent, tuneExplanations) && (
                                            <small><b>{terms.tuningNote}:</b> {formatTuneExplanationNote(action, representativeEvent, tuneExplanations)}</small>
                                        )}
                                    </div>
                                ))}
                            </div>
                        )}
                    </div>

                    <div className="event-section">
                        <h3>{terms.issueGroupEvidence}</h3>
                        <div className="evidence-grid">
                            {Object.entries(group.evidence || {}).slice(0, 8).map(([key, stat]) => (
                                <div key={key}>
                                    <span>{localizedLabel(key, terms.evidenceLabels)}</span>
                                    <strong>{formatEvidenceValueForKey(key, stat.avg, terms)} ({formatEvidenceValueForKey(key, stat.min, terms)}-{formatEvidenceValueForKey(key, stat.max, terms)})</strong>
                                </div>
                            ))}
                        </div>
                    </div>

                    <div className="event-section">
                        <h3>{terms.issueGroupEvents}</h3>
                        <div className="report-event-list">
                            {(group.events || []).map(event => (
                                <span key={event.id}>
                                    {eventLabel(event.type, terms)} / {severityLabel(event.severity, terms)} / {formatDuration(event.durationMs, terms)}
                                </span>
                            ))}
                        </div>
                    </div>
                </div>
            </section>
        </div>
    );
}

function DeveloperModeView({
    current,
    status,
    replayStatus,
    issueSummary,
    sessions,
    ruleProfiles,
    strategyTemplates,
    strategyAnalysis,
    strategySessionIds,
    selectedStrategyTemplateId,
    knowledgeStatus,
    ruleForm,
    editingRuleId,
    language,
    terms,
    busy,
    onRuleSelect,
    onRuleNew,
    onRuleFormChange,
    onRuleSave,
    onRuleReset,
    onRuleDelete,
    onStrategyTemplateChange,
    onToggleStrategySession,
    onRunStrategyAnalysis,
    onReloadKnowledge,
}: {
    current: TelemetryFrame | null;
    status: TelemetryStatus;
    replayStatus: TelemetryReplayStatus;
    issueSummary: SessionIssueSummary | null;
    sessions: TelemetrySession[];
    ruleProfiles: RuleThresholdProfile[];
    strategyTemplates: StrategyTemplate[];
    strategyAnalysis: RoadStrategyAnalysis | null;
    strategySessionIds: number[];
    selectedStrategyTemplateId: number | null;
    knowledgeStatus: RoadTuningKnowledgeStatus | null;
    ruleForm: RuleThresholdProfileInput;
    editingRuleId: number | null;
    language: Lang;
    terms: Copy;
    busy: boolean;
    onRuleSelect: (profile: RuleThresholdProfile) => void;
    onRuleNew: () => void;
    onRuleFormChange: (field: keyof RuleThresholdProfileInput, value: string) => void;
    onRuleSave: () => void;
    onRuleReset: () => void;
    onRuleDelete: () => void;
    onStrategyTemplateChange: (id: number | null) => void;
    onToggleStrategySession: (id: number) => void;
    onRunStrategyAnalysis: () => void;
    onReloadKnowledge: () => void;
}) {
    const groups = developerFieldGroups(current, status, replayStatus, terms, language);
    return (
        <section className="developer-layout">
            <div className="panel developer-panel">
                <div className="panel-heading">
                    <h2>{terms.fieldDiagnostics}</h2>
                    <span>{current ? formatTime(current.receivedAt, terms.never) : terms.noCurrentFrame}</span>
                </div>
                {groups.map(group => (
                    <div className="diagnostic-group" key={group.title}>
                        <h3>{group.title}</h3>
                        <div className="diagnostic-table">
                            <div className="diagnostic-row diagnostic-head">
                                <span>{terms.fieldName}</span>
                                <span>{terms.fieldValue}</span>
                                <span>{terms.fieldUnit}</span>
                                <span>{terms.fieldSource}</span>
                                <span>{terms.fieldRange}</span>
                                <span>{terms.fieldState}</span>
                            </div>
                            {group.rows.map(row => (
                                <div className={`diagnostic-row ${row.warn ? 'warn' : ''}`} key={`${group.title}-${row.name}`}>
                                    <span>{row.name}</span>
                                    <strong>{row.value}</strong>
                                    <span>{row.unit || '--'}</span>
                                    <span>{row.source}</span>
                                    <span>{row.range}</span>
                                    <span>{row.warn ? terms.checkValue : terms.ok}</span>
                                </div>
                            ))}
                        </div>
                    </div>
                ))}
            </div>

            <div className="panel rule-panel">
                <PowerGearTestPanel summary={issueSummary} terms={terms}/>
                <StrategyTemplatePanel
                    sessions={sessions}
                    templates={strategyTemplates}
                    analysis={strategyAnalysis}
                    selectedSessionIds={strategySessionIds}
                    selectedTemplateId={selectedStrategyTemplateId}
                    terms={terms}
                    busy={busy}
                    onTemplateChange={onStrategyTemplateChange}
                    onToggleSession={onToggleStrategySession}
                    onRunAnalysis={onRunStrategyAnalysis}
                />
                <div className="panel-heading">
                    <h2>{terms.knowledgeStatus}</h2>
                    <button className="small-action" type="button" onClick={onReloadKnowledge} disabled={busy}>
                        <Activity size={15}/>
                        {terms.reloadKnowledge}
                    </button>
                </div>
                <div className="report-summary">
                    <TextStat label={terms.knowledgeSymptoms} value={String(knowledgeStatus?.symptomCount ?? 0)}/>
                    <TextStat label={terms.knowledgeActions} value={String(knowledgeStatus?.actionCount ?? 0)}/>
                    <TextStat label={terms.knowledgeSource} value={knowledgeStatus?.sourcePath || '--'}/>
                    <TextStat label={terms.reportStatus} value={knowledgeStatus?.lastError ? terms.knowledgeFallback : terms.ok}/>
                </div>
                {knowledgeStatus?.lastError && <div className="empty-events">{knowledgeStatus.lastError}</div>}
                <div className="panel-heading">
                    <h2>{terms.ruleThresholds}</h2>
                    <button className="small-action" type="button" onClick={onRuleNew} disabled={busy}>
                        <Plus size={15}/>
                        {terms.newProfile}
                    </button>
                </div>
                <div className="rule-profile-list">
                    {ruleProfiles.map(profile => (
                        <button key={profile.id} className={`rule-profile-row ${editingRuleId === profile.id ? 'selected' : ''}`} type="button" onClick={() => onRuleSelect(profile)}>
                            <strong>{profile.name}</strong>
                            <span>{[profile.drivetrain, profile.carClass, profile.useCase].filter(Boolean).join(' / ') || (profile.isDefault ? 'Default' : '--')}</span>
                        </button>
                    ))}
                </div>
                <div className="rule-form">
                    <label><span>{terms.ruleName}</span><input value={ruleForm.name} onChange={(event) => onRuleFormChange('name', event.target.value)}/></label>
                    <label><span>{terms.ruleCarClass}</span><input value={ruleForm.carClass} onChange={(event) => onRuleFormChange('carClass', event.target.value)}/></label>
                    <label><span>{terms.ruleDrivetrain}</span><input value={ruleForm.drivetrain} onChange={(event) => onRuleFormChange('drivetrain', event.target.value)}/></label>
                    <label><span>{terms.ruleUseCase}</span><input value={ruleForm.useCase} onChange={(event) => onRuleFormChange('useCase', event.target.value)}/></label>
                    <label className="wide"><span>{terms.ruleConfigJson}</span><textarea value={ruleForm.configJson} onChange={(event) => onRuleFormChange('configJson', event.target.value)}/></label>
                    <div className="form-actions">
                        <button className="action secondary" type="button" onClick={onRuleReset} disabled={!editingRuleId || busy}>{terms.resetDefaults}</button>
                        <button className="action secondary" type="button" onClick={onRuleDelete} disabled={!editingRuleId || busy || ruleProfiles.find(profile => profile.id === editingRuleId)?.isDefault}>{terms.delete}</button>
                        <button className="action primary" type="button" onClick={onRuleSave} disabled={busy}>{editingRuleId ? terms.updateProfile : terms.createProfile}</button>
                    </div>
                </div>
            </div>
        </section>
    );
}

function PowerGearTestPanel({summary, terms}: { summary: SessionIssueSummary | null; terms: Copy }) {
    const gearPower = summary?.gearPower;
    const rows = (gearPower?.gears || []).filter(gear => gear.sampleCount > 0).slice(0, 10);
    return (
        <div className="panel-inner issue-section">
            <div className="panel-heading compact">
                <h2>{terms.powerGearTest}</h2>
                <span>{localizedLabel(gearPower?.confidence || '', terms.planConfidenceLabels) || terms.insufficientData}</span>
            </div>
            <div className="empty-events">{terms.powerGearTestHint}</div>
            {gearPower && gearPower.status !== 'insufficient_data' && (
                <div className="report-summary">
                    <TextStat label={terms.powerBandTarget} value={gearPower.powerBandStartRPM && gearPower.powerBandEndRPM ? `${formatNumber(gearPower.powerBandStartRPM, 0)} - ${formatNumber(gearPower.powerBandEndRPM, 0)} rpm` : '--'}/>
                    <TextStat label={terms.powerBandSource} value={localizedLabel(gearPower.powerBandSource || '', terms.powerBandSources) || '--'}/>
                    <TextStat label={terms.samplesSaved} value={String(Math.round(gearPower.evidence?.power_band_high_load_samples || 0))}/>
                    <TextStat label={terms.gearFinding} value={localizedLabel(gearPower.summary || '', terms.gearFindings) || '--'}/>
                </div>
            )}
            {rows.length > 0 && (
                <div className="diagnostic-table compact-diagnostic">
                    <div className="diagnostic-row diagnostic-head">
                        <span>{terms.gear}</span>
                        <span>{terms.speedRange}</span>
                        <span>{terms.samplesSaved}</span>
                        <span>{terms.inPowerBand}</span>
                        <span>{terms.acceleration}</span>
                        <span>{terms.shiftAfter}</span>
                        <span>{terms.gearFinding}</span>
                    </div>
                    {rows.map(gear => (
                        <div className={`diagnostic-row ${gear.finding !== 'ok' ? 'warn' : ''}`} key={`power-test-${gear.gear}`}>
                            <span>{gear.gear}</span>
                            <GearSpeedValue gear={gear} terms={terms}/>
                            <strong>{gear.highLoadSampleCount}</strong>
                            <span>{formatGearInPowerBandRange(gear)}</span>
                            <span>{formatNumber(gear.accelAvgMps2 || 0, 2)} m/s²</span>
                            <span>{gear.shiftAfterRPM ? formatNumber(gear.shiftAfterRPM, 0) : '--'}</span>
                            <span>{localizedLabel(gear.finding, terms.gearFindings)}</span>
                        </div>
                    ))}
                </div>
            )}
        </div>
    );
}

function StrategyTemplatesView({
    sessions,
    templates,
    analysis,
    selectedSessionIds,
    selectedTemplateId,
    terms,
    busy,
    onTemplateChange,
    onToggleSession,
    onRunAnalysis,
}: {
    sessions: TelemetrySession[];
    templates: StrategyTemplate[];
    analysis: RoadStrategyAnalysis | null;
    selectedSessionIds: number[];
    selectedTemplateId: number | null;
    terms: Copy;
    busy: boolean;
    onTemplateChange: (id: number | null) => void;
    onToggleSession: (id: number) => void;
    onRunAnalysis: () => void;
}) {
    return (
        <section className="strategy-view">
            <div className="panel">
                <div className="panel-heading">
                    <div>
                        <h2>{terms.strategyTemplates}</h2>
                        <span>{terms.strategyAnalysisEmpty}</span>
                    </div>
                </div>
                <StrategyTemplatePanel
                    sessions={sessions}
                    templates={templates}
                    analysis={analysis}
                    selectedSessionIds={selectedSessionIds}
                    selectedTemplateId={selectedTemplateId}
                    terms={terms}
                    busy={busy}
                    onTemplateChange={onTemplateChange}
                    onToggleSession={onToggleSession}
                    onRunAnalysis={onRunAnalysis}
                />
            </div>
        </section>
    );
}

function StrategyTemplatePanel({
    sessions,
    templates,
    analysis,
    selectedSessionIds,
    selectedTemplateId,
    terms,
    busy,
    onTemplateChange,
    onToggleSession,
    onRunAnalysis,
}: {
    sessions: TelemetrySession[];
    templates: StrategyTemplate[];
    analysis: RoadStrategyAnalysis | null;
    selectedSessionIds: number[];
    selectedTemplateId: number | null;
    terms: Copy;
    busy: boolean;
    onTemplateChange: (id: number | null) => void;
    onToggleSession: (id: number) => void;
    onRunAnalysis: () => void;
}) {
    return (
        <section className="strategy-panel">
            <div className="panel-heading compact">
                <h2>{terms.strategyTemplates}</h2>
                <span>{templates.length}</span>
            </div>
            <div className="rule-profile-list">
                {templates.map(template => (
                    <button key={template.id} className={`rule-profile-row ${selectedTemplateId === template.id ? 'selected' : ''}`} type="button" onClick={() => onTemplateChange(template.id)}>
                        <strong>{template.name}</strong>
                        <span>{[template.drivetrain, template.carClass, template.useCase, gameModeLabel(template.gameMode as GameMode, terms)].filter(Boolean).join(' / ')}</span>
                        <small>{terms.enabledEvents}: {template.enabledEventCount}/{template.totalEventCount}</small>
                    </button>
                ))}
            </div>

            <div className="panel-heading compact">
                <h2>{terms.strategyAnalysis}</h2>
                <span>{selectedSessionIds.length}/5</span>
            </div>
            <div className="strategy-session-list">
                {sessions.slice(0, 20).map(session => (
                    <label key={session.id} className="strategy-session-row">
                        <input
                            type="checkbox"
                            checked={selectedSessionIds.includes(session.id)}
                            disabled={busy || (!selectedSessionIds.includes(session.id) && selectedSessionIds.length >= 5)}
                            onChange={() => onToggleSession(session.id)}
                        />
                        <span>{session.sessionName || `#${session.id}`}</span>
                        <small>{formatTime(session.startedAt, '--')} / {session.tuneName || terms.noProfile}</small>
                    </label>
                ))}
            </div>
            <div className="form-actions">
                <button className="action primary" type="button" onClick={onRunAnalysis} disabled={busy || !selectedTemplateId || selectedSessionIds.length === 0}>{terms.runStrategyAnalysis}</button>
            </div>
            {!analysis ? (
                <div className="empty-events">{terms.strategyAnalysisEmpty}</div>
            ) : (
                <div className="strategy-analysis-result">
                    <div className="report-summary">
                        <TextStat label={terms.strategyTemplate} value={analysis.template.name}/>
                        <Stat label={terms.reportSessions} value={analysis.sessionCount}/>
                        <Stat label={terms.totalEvents} value={analysis.totalEvents}/>
                        <Stat label={terms.eventDistribution} value={analysis.eventDistribution.length}/>
                    </div>
                    {analysis.issueGroups.length > 0 && (
                        <div className="report-event-list">
                            {analysis.issueGroups.map(group => (
                                <span key={group.family}>
                                    {issueFamilyLabel(group.family, terms)} / {terms.eventsSaved}: {group.eventCount} / {terms.reportSessions}: {group.sessionCount} / {terms.strategyRecommendation}: {strategyRecommendationLabel(group.recommendation, terms)}
                                </span>
                            ))}
                        </div>
                    )}
                    {analysis.eventDistribution.length > 0 && (
                        <div className="report-event-list">
                            {analysis.eventDistribution.map(item => (
                                <span key={item.type}>{eventLabel(item.type, terms)}: {item.count} / {severityLabel(item.severity, terms)}</span>
                            ))}
                        </div>
                    )}
                    {analysis.hints.length > 0 && (
                        <div className="empty-events">
                            {analysis.hints.map((hint, index) => (
                                <span key={`${hint.message}-${hint.eventType || hint.family || index}`}>{strategyHintLabel(hint, terms)}</span>
                            ))}
                        </div>
                    )}
                </div>
            )}
        </section>
    );
}

function TrackProfilesView({
    current,
    status,
    tracks,
    profiles,
    selectedTrackId,
    selectedTrackProfile,
    capture,
    terms,
    busy,
    onCaptureNameChange,
    onStartCapture,
    onStopCapture,
    onSaveCapture,
    onSelectTrack,
    onRenameTrack,
    onDeleteTrack,
    onStartBaseline,
    onStopBaseline,
    onSaveBaseline,
    onDeleteBaseline,
}: {
    current: TelemetryFrame | null;
    status: TelemetryStatus;
    tracks: BenchmarkTrack[];
    profiles: TrackProfile[];
    selectedTrackId: number | null;
    selectedTrackProfile: TrackProfile | null;
    capture: TrackCaptureState;
    terms: Copy;
    busy: boolean;
    onCaptureNameChange: (name: string) => void;
    onStartCapture: () => void;
    onStopCapture: () => void;
    onSaveCapture: (trackType: BenchmarkTrackType, extractionMode: BenchmarkExtractionMode, gateWidth: number, gateDepth: number, startGate?: BenchmarkGate, finishGate?: BenchmarkGate) => void;
    onSelectTrack: (trackId: number | null) => void;
    onRenameTrack: (track: BenchmarkTrack) => void;
    onDeleteTrack: (track: BenchmarkTrack) => void;
    onStartBaseline: () => void;
    onStopBaseline: () => void;
    onSaveBaseline: () => void;
    onDeleteBaseline: (run: TrackBaselineRun) => void;
}) {
    const selectedTrack = tracks.find(track => track.id === selectedTrackId) || selectedTrackProfile?.track || null;
    const selectedProfile = selectedTrackProfile || profiles.find(profile => profile.track.id === selectedTrack?.id) || null;
    const [trackType, setTrackType] = useState<BenchmarkTrackType>('auto');
    const isTrackCaptureActive = status.running && status.analysisMode === 'track_capture';
    const isBaselineActive = status.running && status.analysisMode === 'track_baseline';
    const isBaselineMode = status.analysisMode === 'track_baseline';
    const routeMeters = capture.points.length > 1 ? pointsRouteLength(capture.points) : selectedTrack?.routeLengthMeters || 0;
    const previewPoints = capture.points.length > 0 ? capture.points : selectedTrack?.polyline || [];
    return (
        <section className="track-profiles-layout">
            <div className="panel track-profile-sidebar">
                <div className="panel-heading">
                    <div>
                        <h2>{terms.trackProfilesTitle}</h2>
                        <span>{terms.trackProfilesSubtitle}</span>
                    </div>
                </div>
                <div className="status-alerts">
                    <div className="status-alert ok">{terms.trackCaptureNoHistory}</div>
                </div>
                <div className="track-list">
                    {tracks.length === 0 ? (
                        <div className="empty-events">{terms.noTracks}</div>
                    ) : tracks.map(track => (
                        <button key={track.id} className={`rule-profile-row ${selectedTrack?.id === track.id ? 'selected' : ''}`} type="button" onClick={() => onSelectTrack(track.id)}>
                            <strong>{track.name}</strong>
                            <span>{trackTypeLabel(track.trackType as BenchmarkTrackType, terms)} / {track.routeLengthMeters.toFixed(0)} m</span>
                            <small>{terms.trackBaselines}: {trackBaselineCount(profiles.find(profile => profile.track.id === track.id))}</small>
                        </button>
                    ))}
                </div>
            </div>

            <div className="panel track-profile-main">
                <div className="panel-heading">
                    <div>
                        <h2>{selectedTrack?.name || terms.trackProfilesTitle}</h2>
                        <span>{terms.trackData}</span>
                    </div>
                    {selectedTrack && (
                        <div className="profile-actions">
                            <button className="small-action" type="button" onClick={() => onRenameTrack(selectedTrack)} disabled={busy}>{terms.renameTrack}</button>
                            <button className="small-action danger" type="button" onClick={() => onDeleteTrack(selectedTrack)} disabled={busy}>{terms.delete}</button>
                        </div>
                    )}
                </div>
                <div className="track-capture-actions">
                    <div className="track-action-card">
                        <div className="panel-heading compact">
                            <h2>{terms.trackCaptureMode}</h2>
                            <span>{terms.trackCaptureNoHistory}</span>
                        </div>
                        <div className="track-action-fields">
                            <label>
                                <span>{terms.trackName}</span>
                                <input value={capture.name} onChange={(event) => onCaptureNameChange(event.target.value)}/>
                            </label>
                            <label>
                                <span>{terms.trackType}</span>
                                <select value={trackType} onChange={(event) => setTrackType(event.target.value as BenchmarkTrackType)}>
                                    <option value="auto">{terms.autoTrackType}</option>
                                    <option value="circuit">{terms.circuitTrack}</option>
                                    <option value="sprint">{terms.sprintTrack}</option>
                                </select>
                            </label>
                        </div>
                        <div className="form-actions">
                            {capture.recording ? (
                                <button className="action secondary" type="button" onClick={onStopCapture} disabled={busy}>{terms.stopCapture}</button>
                            ) : (
                                <button className="action secondary" type="button" onClick={onStartCapture} disabled={busy || !isTrackCaptureActive || !current}>{terms.startCapture}</button>
                            )}
                            <button className="action primary" type="button" onClick={() => onSaveCapture(trackType, 'auto_best_lap', 30, 20)} disabled={busy || capture.points.length < 2}>{terms.saveTrack}</button>
                        </div>
                    </div>
                    <div className="track-action-card">
                        <div className="panel-heading compact">
                            <h2>{terms.trackBaselineCapture}</h2>
                            <span>{terms.trackBaselineNoSession}</span>
                        </div>
                        <div className="status-alerts">
                            <div className="status-alert ok">{terms.trackBaselineAutoArchiveHint}</div>
                        </div>
                        <div className="form-actions">
                            {isBaselineActive ? (
                                <>
                                    <button className="action secondary" type="button" onClick={onStopBaseline} disabled={busy}>{terms.stopCapture}</button>
                                    <button className="action primary" type="button" onClick={onSaveBaseline} disabled={busy}>{terms.saveTrackBaseline}</button>
                                </>
                            ) : isBaselineMode ? (
                                <button className="action primary" type="button" onClick={onSaveBaseline} disabled={busy}>{terms.saveTrackBaseline}</button>
                            ) : (
                                <button className="action secondary" type="button" onClick={onStartBaseline} disabled={busy}>{terms.startTrackBaseline}</button>
                            )}
                        </div>
                    </div>
                </div>
                <TrackProfileSummary profile={selectedProfile} selectedTrack={selectedTrack} terms={terms}/>
                <div className="track-profile-summary">
                    <div className="panel-heading compact">
                        <h2>{terms.trackData}</h2>
                        <span>{terms.capturePoints}: {capture.points.length}</span>
                    </div>
                    <div className="track-stats">
                        <Stat label={terms.capturePoints} value={capture.points.length}/>
                        <TextStat label={terms.routeLength} value={`${routeMeters.toFixed(0)} m`}/>
                        <TextStat label={terms.trackType} value={trackTypeLabel(trackType, terms)}/>
                        <Stat label={terms.observedLaps} value={countCapturedLapIncrements(capture.points)}/>
                    </div>
                    <TrackPreview points={previewPoints} startGate={selectedTrack?.startGate} finishGate={selectedTrack?.finishGate} checkpoints={selectedTrack?.checkpoints || []} terms={terms}/>
                </div>
                <TrackBaselinePanel
                    profile={selectedProfile}
                    selectedTrack={selectedTrack}
                    terms={terms}
                    busy={busy}
                    baselineActive={isBaselineActive}
                    baselineMode={isBaselineMode}
                    onStartBaseline={onStartBaseline}
                    onStopBaseline={onStopBaseline}
                    onSaveBaseline={onSaveBaseline}
                    onDeleteBaseline={onDeleteBaseline}
                />
            </div>
        </section>
    );
}

function TrackProfileSummary({profile, selectedTrack, terms}: { profile: TrackProfile | null; selectedTrack: BenchmarkTrack | null; terms: Copy }) {
    const track = profile?.track || selectedTrack;
    if (!track) {
        return <div className="empty-events">{terms.noTracks}</div>;
    }
    return (
        <div className="track-profile-summary">
            <div className="report-summary">
                <TextStat label={terms.trackType} value={trackTypeLabel(track.trackType as BenchmarkTrackType, terms)}/>
                <TextStat label={terms.routeLength} value={`${track.routeLengthMeters.toFixed(0)} m`}/>
                <Stat label={terms.observedLaps} value={track.lapCountObserved || 0}/>
                <TextStat label={terms.startDistance} value={`${formatPointShort(track.start)}`}/>
                <TextStat label={terms.endDistance} value={`${formatPointShort(track.end)}`}/>
                <TextStat label={terms.updatedAt} value={formatTime(track.updatedAt, '--')}/>
            </div>
            {profile && profile.warnings.length > 0 && (
                <div className="status-alerts">
                    {profile.warnings.map(warning => (
                        <div className="status-alert warn" key={warning}>
                            {localizedLabel(warning, terms.trackProfileWarnings) || warning}
                        </div>
                    ))}
                </div>
            )}
            <div className="panel-heading compact">
                <h2>{terms.vehicleReferences}</h2>
                <span>{trackBaselineCount(profile)}</span>
            </div>
            {!profile || profile.vehicleReferences.length === 0 ? (
                <div className="empty-events">{terms.noVehicleReferences}</div>
            ) : (
                <div className="vehicle-reference-table">
                    <div className="vehicle-reference-row header">
                        <span>{terms.baselineVehicle}</span>
                        <span>{terms.avgSpeed}</span>
                        <span>{terms.maxSpeed}</span>
                    </div>
                    {profile.vehicleReferences.map(reference => (
                        <div className="vehicle-reference-row" key={trackVehicleKeyLabel(reference.vehicle)}>
                            <strong>{reference.vehicle.label || trackVehicleKeyLabel(reference.vehicle)}</strong>
                            <span>{formatNumber(reference.avgSpeedKmh || undefined, 0)} km/h</span>
                            <span>{formatNumber(reference.maxSpeedKmh || undefined, 0)} km/h</span>
                        </div>
                    ))}
                </div>
            )}
        </div>
    );
}

function TrackBaselinePanel({
    profile,
    selectedTrack,
    terms,
    busy,
    baselineActive,
    baselineMode,
    onStartBaseline,
    onStopBaseline,
    onSaveBaseline,
    onDeleteBaseline,
}: {
    profile: TrackProfile | null;
    selectedTrack: BenchmarkTrack | null;
    terms: Copy;
    busy: boolean;
    baselineActive: boolean;
    baselineMode: boolean;
    onStartBaseline: () => void;
    onStopBaseline: () => void;
    onSaveBaseline: () => void;
    onDeleteBaseline: (run: TrackBaselineRun) => void;
}) {
    const baselines = (profile?.vehicleReferences || []).flatMap(reference => reference.recentBaselineRuns || []);
    return (
        <div className="track-profile-summary">
            <div className="panel-heading compact">
                <h2>{terms.trackBaselines}</h2>
                <span>{terms.trackBaselineNoSession}</span>
            </div>
            <div className="form-actions">
                            {baselineActive ? (
                                <>
                                    <button className="action secondary" type="button" onClick={onStopBaseline} disabled={busy}>{terms.stopCapture}</button>
                                    <button className="action primary" type="button" onClick={onSaveBaseline} disabled={busy}>{terms.saveTrackBaseline}</button>
                                </>
                            ) : baselineMode ? (
                                <button className="action primary" type="button" onClick={onSaveBaseline} disabled={busy}>{terms.saveTrackBaseline}</button>
                            ) : (
                                <button className="action secondary" type="button" onClick={onStartBaseline} disabled={busy}>{terms.startTrackBaseline}</button>
                            )}
            </div>
            {baselines.length === 0 ? (
                <div className="empty-events">{terms.noVehicleReferences}</div>
            ) : (
                <div className="track-run-list">
                    {baselines.map(run => (
                        <div className="run-row" key={run.id}>
                            <strong>{run.vehicle.label || trackVehicleKeyLabel(run.vehicle)}</strong>
                            <div className="run-diagnostics">
                                <span>{formatDuration(run.durationMs, terms)}</span>
                                <span>{terms.driverMode}: {testConditionLabel(run.driverMode, terms)} / {formatPercentValue(run.driverModeConfidence)}</span>
                                <span>{terms.avgSpeed}: {formatNumber(run.avgSpeedKmh || undefined, 0)} km/h</span>
                                <span>{terms.maxSpeed}: {formatNumber(run.maxSpeedKmh || undefined, 0)} km/h</span>
                                <span>{terms.confidence}: {formatPercentValue(run.confidence)}</span>
                            </div>
                            <button className="small-action danger" type="button" onClick={() => onDeleteBaseline(run)} disabled={busy}>{terms.delete}</button>
                        </div>
                    ))}
                </div>
            )}
        </div>
    );
}

function TrackBuilderPanel({
    current,
    sessions,
    tracks,
    selectedTrack,
    selectedTrackRuns,
    capture,
    terms,
    busy,
    onCaptureNameChange,
    onStartCapture,
    onStopCapture,
    onSaveCapture,
    onCreateFromSession,
    onReextractTrack,
    onSelectTrack,
    onDeleteTrack,
    onAnalyzeSessionRuns,
    showTrackList = true,
}: {
    current: TelemetryFrame | null;
    sessions: TelemetrySession[];
    tracks: BenchmarkTrack[];
    selectedTrack: BenchmarkTrack | null;
    selectedTrackRuns: BenchmarkRun[];
    capture: TrackCaptureState;
    terms: Copy;
    busy: boolean;
    onCaptureNameChange: (name: string) => void;
    onStartCapture: () => void;
    onStopCapture: () => void;
    onSaveCapture: (trackType: BenchmarkTrackType, extractionMode: BenchmarkExtractionMode, gateWidth: number, gateDepth: number, startGate?: BenchmarkGate, finishGate?: BenchmarkGate) => void;
    onCreateFromSession: (sessionId: number, name: string, trackType: BenchmarkTrackType, extractionMode: BenchmarkExtractionMode, gateWidth: number, gateDepth: number, startGate?: BenchmarkGate, finishGate?: BenchmarkGate) => void;
    onReextractTrack: (trackId: number, sessionId: number, name: string, trackType: BenchmarkTrackType, extractionMode: BenchmarkExtractionMode, gateWidth: number, gateDepth: number, startGate?: BenchmarkGate, finishGate?: BenchmarkGate) => void;
    onSelectTrack: (trackId: number | null) => void;
    onDeleteTrack: (track: BenchmarkTrack) => void;
    onAnalyzeSessionRuns: () => void;
    showTrackList?: boolean;
}) {
    const [sourceSessionId, setSourceSessionId] = useState<number>(sessions[0]?.id || 0);
    const [sourceTrackName, setSourceTrackName] = useState('');
    const [trackType, setTrackType] = useState<BenchmarkTrackType>('auto');
    const [extractionMode, setExtractionMode] = useState<BenchmarkExtractionMode>('auto_best_lap');
    const [startGate, setStartGate] = useState<BenchmarkGate | undefined>(undefined);
    const [finishGate, setFinishGate] = useState<BenchmarkGate | undefined>(undefined);
    const [gateWidth, setGateWidth] = useState(30);
    const [gateDepth, setGateDepth] = useState(20);
    const sourceSession = sessions.find(session => session.id === sourceSessionId) || null;
    const previewPoints = capture.points.length > 0 ? capture.points : selectedTrack?.polyline || [];
    const routeMeters = capture.points.length > 1 ? pointsRouteLength(capture.points) : selectedTrack?.routeLengthMeters || 0;
    const hasDrivingLine = capture.points.length > 0 ? capture.hasDrivingLine : !!selectedTrack?.hasDrivingLine;
    const previewStartGate = startGate || selectedTrack?.startGate;
    const previewFinishGate = finishGate || selectedTrack?.finishGate;
    const previewCheckpoints = selectedTrack?.checkpoints || [];

    useEffect(() => {
        if (!sourceSessionId && sessions[0]?.id) {
            setSourceSessionId(sessions[0].id);
        }
    }, [sessions, sourceSessionId]);

    useEffect(() => {
        if (selectedTrack) {
            setTrackType((selectedTrack.trackType || 'auto') as BenchmarkTrackType);
            setGateWidth(selectedTrack.startGate?.widthMeters || 30);
            setGateDepth(selectedTrack.startGate?.depthMeters || 20);
        }
    }, [selectedTrack?.id]);

    return (
        <section className="track-builder">
            <div className="panel-heading compact">
                <h2>{terms.trackBuilder}</h2>
                <span>{terms.benchmarkTracks}</span>
            </div>
            <div className="track-stats">
                <TextStat label={terms.worldPosition} value={current ? `${formatNumber(current.positionX, 0)}, ${formatNumber(current.positionY, 0)}, ${formatNumber(current.positionZ, 0)}` : '--'}/>
                <Stat label={terms.capturePoints} value={capture.points.length}/>
                <TextStat label={terms.routeLength} value={`${routeMeters.toFixed(0)} m`}/>
                <TextStat label={terms.drivingLineSignal} value={hasDrivingLine ? terms.detected : terms.notDetected}/>
                <TextStat label={terms.trackType} value={selectedTrack ? trackTypeLabel(selectedTrack.trackType as BenchmarkTrackType, terms) : trackTypeLabel(trackType, terms)}/>
                <TextStat label={terms.extractionMode} value={extractionModeLabel(extractionMode, terms)}/>
                <Stat label={terms.observedLaps} value={selectedTrack?.lapCountObserved || 0}/>
            </div>
            <TrackPreview points={previewPoints} startGate={previewStartGate} finishGate={previewFinishGate} checkpoints={previewCheckpoints} terms={terms}/>
            <div className="rule-form">
                <label>
                    <span>{terms.trackName}</span>
                    <input value={capture.name} onChange={(event) => onCaptureNameChange(event.target.value)}/>
                </label>
                <label>
                    <span>{terms.trackType}</span>
                    <select value={trackType} onChange={(event) => setTrackType(event.target.value as BenchmarkTrackType)}>
                        <option value="auto">{terms.autoTrackType}</option>
                        <option value="circuit">{terms.circuitTrack}</option>
                        <option value="sprint">{terms.sprintTrack}</option>
                    </select>
                </label>
                <label>
                    <span>{terms.extractionMode}</span>
                    <select value={extractionMode} onChange={(event) => setExtractionMode(event.target.value as BenchmarkExtractionMode)}>
                        <option value="auto_best_lap">{terms.autoBestLap}</option>
                        <option value="first_lap">{terms.firstLap}</option>
                        <option value="full_segment">{terms.fullSegment}</option>
                    </select>
                </label>
                <label>
                    <span>{terms.gateWidth}</span>
                    <input type="number" min={5} step={1} value={gateWidth} onChange={(event) => setGateWidth(Number(event.target.value) || 30)}/>
                </label>
                <label>
                    <span>{terms.gateDepth}</span>
                    <input type="number" min={5} step={1} value={gateDepth} onChange={(event) => setGateDepth(Number(event.target.value) || 20)}/>
                </label>
                <div className="form-actions">
                    <button className="small-action" type="button" onClick={() => current && setStartGate(gateFromCurrent(current, gateWidth, gateDepth))} disabled={!current}>{terms.setStartGate}</button>
                    <button className="small-action" type="button" onClick={() => current && setFinishGate(gateFromCurrent(current, gateWidth, gateDepth))} disabled={!current}>{terms.setFinishGate}</button>
                    <button className="small-action" type="button" onClick={() => { setStartGate(undefined); setFinishGate(undefined); }}>{terms.clearGates}</button>
                </div>
                <div className="form-actions">
                    {capture.recording ? (
                        <button className="action secondary" type="button" onClick={onStopCapture} disabled={busy}>{terms.stopCapture}</button>
                    ) : (
                        <button className="action secondary" type="button" onClick={onStartCapture} disabled={busy || !current}>{terms.startCapture}</button>
                    )}
                    <button className="action primary" type="button" onClick={() => onSaveCapture(trackType, extractionMode, gateWidth, gateDepth, startGate, finishGate)} disabled={busy || capture.points.length < 2}>{terms.saveTrack}</button>
                </div>
            </div>
            <div className="rule-form">
                <label>
                    <span>{terms.reportSessions}</span>
                    <select value={sourceSessionId || 0} onChange={(event) => setSourceSessionId(Number(event.target.value))}>
                        <option value={0}>--</option>
                        {sessions.map(session => <option key={session.id} value={session.id}>{session.sessionName || `#${session.id}`}</option>)}
                    </select>
                </label>
                <label>
                    <span>{terms.trackName}</span>
                    <input value={sourceTrackName} onChange={(event) => setSourceTrackName(event.target.value)} placeholder={sourceSession?.sessionName || ''}/>
                </label>
                <div className="form-actions">
                    <button className="action secondary" type="button" onClick={() => onCreateFromSession(sourceSessionId, sourceTrackName || sourceSession?.sessionName || '', trackType, extractionMode, gateWidth, gateDepth, startGate, finishGate)} disabled={busy || !sourceSessionId}>{terms.fromSession}</button>
                    <button className="action secondary" type="button" onClick={() => selectedTrack && onReextractTrack(selectedTrack.id, sourceSessionId, sourceTrackName || selectedTrack.name, trackType, extractionMode, gateWidth, gateDepth, startGate, finishGate)} disabled={busy || !selectedTrack || !sourceSessionId}>{terms.reextractTrack}</button>
                    <button className="action secondary" type="button" onClick={onAnalyzeSessionRuns} disabled={busy || sessions.length === 0}>{terms.analyzeTrackRuns}</button>
                </div>
            </div>
            {showTrackList && (
                <div className="track-list">
                    {tracks.length === 0 ? (
                        <div className="empty-events">{terms.noTracks}</div>
                    ) : tracks.map(track => (
                        <button key={track.id} className={`rule-profile-row ${selectedTrack?.id === track.id ? 'selected' : ''}`} type="button" onClick={() => onSelectTrack(track.id)}>
                            <strong>{track.name}</strong>
                            <span>{trackTypeLabel(track.trackType as BenchmarkTrackType, terms)} / {gameModeLabel(track.sourceMode as GameMode, terms)} / {track.routeLengthMeters.toFixed(0)} m / {terms.observedLaps}: {track.lapCountObserved || 0}</span>
                        </button>
                    ))}
                </div>
            )}
            {selectedTrack && (
                <div className="track-run-list">
                    <div className="panel-heading compact">
                        <h2>{terms.benchmarkRuns}</h2>
                        <button className="small-action danger" type="button" onClick={() => onDeleteTrack(selectedTrack)} disabled={busy}>{terms.delete}</button>
                    </div>
                    {selectedTrackRuns.length === 0 ? <div className="empty-events">{terms.noBenchmarkRuns}</div> : selectedTrackRuns.map(run => (
                        <div className="run-row" key={run.id}>
                            <strong>{formatDuration(run.durationMs, terms)}</strong>
                            <div className="run-diagnostics">
                                <span>{terms.confidence}: {formatPercentValue(run.confidence)}</span>
                                <span>{terms.avgSpeed}: {formatNumber(run.avgSpeedKmh || undefined, 0)} km/h</span>
                                <BenchmarkRunDiagnostics run={run} terms={terms}/>
                            </div>
                        </div>
                    ))}
                </div>
            )}
        </section>
    );
}

function TrackPreview({
    points,
    startGate,
    finishGate,
    checkpoints,
    terms,
}: {
    points: BenchmarkPoint[];
    startGate?: BenchmarkGate;
    finishGate?: BenchmarkGate;
    checkpoints?: BenchmarkPoint[];
    terms: Copy;
}) {
    const path = routeSvgPath(points);
    return (
        <div className="track-preview">
            {points.length < 2 ? (
                <span>{terms.noTrackPoints}</span>
            ) : (
                <svg viewBox="0 0 320 180" role="img" aria-label={terms.benchmarkTracks}>
                    <path d={path} fill="none" stroke="#7dd6c5" strokeWidth="3" strokeLinecap="round" strokeLinejoin="round"/>
                    {(checkpoints || []).map((point, index) => {
                        const p = routeSvgPoint(point, points);
                        return <circle key={`checkpoint-${index}`} cx={p.x} cy={p.y} r="3.5" fill="#f5c05b"/>;
                    })}
                    {startGate && <GateMarker gate={startGate} points={points} color="#62d26f" label={terms.startGate}/>}
                    {finishGate && <GateMarker gate={finishGate} points={points} color="#ff8a82" label={terms.finishGate}/>}
                    <circle cx={routeSvgPoint(points[0], points).x} cy={routeSvgPoint(points[0], points).y} r="5" fill="#62d26f"/>
                    <circle cx={routeSvgPoint(points[points.length - 1], points).x} cy={routeSvgPoint(points[points.length - 1], points).y} r="5" fill="#ff8a82"/>
                </svg>
            )}
        </div>
    );
}

function GateMarker({gate, points, color, label}: { gate: BenchmarkGate; points: BenchmarkPoint[]; color: string; label: string }) {
    const center = routeSvgPoint(gate.center, points);
    const dx = gate.directionX || 1;
    const dz = gate.directionZ || 0;
    const length = Math.hypot(dx, dz) || 1;
    const nx = -dz / length;
    const nz = dx / length;
    const scale = Math.max(8, Math.min(28, (gate.widthMeters || 30) * 0.25));
    const a = routeSvgPoint({x: gate.center.x + nx * scale, y: gate.center.y, z: gate.center.z + nz * scale}, points);
    const b = routeSvgPoint({x: gate.center.x - nx * scale, y: gate.center.y, z: gate.center.z - nz * scale}, points);
    return (
        <g>
            <title>{label}</title>
            <line x1={a.x} y1={a.y} x2={b.x} y2={b.y} stroke={color} strokeWidth="2.5" strokeLinecap="round"/>
            <circle cx={center.x} cy={center.y} r="3" fill={color}/>
        </g>
    );
}

function BenchmarkRunDiagnostics({run, terms}: { run: BenchmarkRun; terms: Copy }) {
    const warnings = localizedWarningFlags(run.warningFlags, terms);
    return (
        <div className="run-diagnostic-grid">
            <span>{terms.routeProgress}: {formatPercentValue(run.routeProgress01 ?? undefined)}</span>
            <span>{terms.geometryLength}: {formatMeters(run.geometryLengthMeters)}</span>
            <span>{terms.lengthError}: {formatSignedPercent(run.trackLengthErrorPct)}</span>
            <span>{terms.lateralError}: {formatMeters(run.avgLateralErrorMeters)} / {formatMeters(run.maxLateralErrorMeters)}</span>
            <span>{terms.distanceDelta}: {formatMeters(run.distanceTraveledDeltaMeters)}</span>
            <span>{terms.raceTimeDelta}: {formatSecondsValue(run.currentRaceTimeDeltaSeconds)}</span>
            <span className={warnings.length > 0 ? 'warn' : ''}>{terms.warnings}: {warnings.length > 0 ? warnings.join(' / ') : terms.noWarnings}</span>
        </div>
    );
}

function EventTimeline({
    events,
    selectedEvent,
    selectedEventId,
    onSelect,
    terms,
}: {
    events: DetectedEvent[];
    selectedEvent: DetectedEvent | null;
    selectedEventId: string;
    onSelect: (id: string) => void;
    terms: Copy;
}) {
    const orderedEvents = [...events].reverse();

    return (
        <section className="panel events-panel">
            <div className="panel-heading">
                <div className="event-title">
                    <AlertTriangle size={19}/>
                    <h2>{terms.eventTimeline}</h2>
                </div>
                <span>{terms.eventSubtitle}</span>
            </div>
            {orderedEvents.length === 0 ? (
                <div className="empty-events">{terms.noEvents}</div>
            ) : (
                <div className="events-layout">
                    <div className="event-list" role="list">
                        {orderedEvents.map(event => (
                            <button
                                key={event.id}
                                className={`event-row ${selectedEventId === event.id ? 'selected' : ''}`}
                                onClick={() => onSelect(event.id)}
                                type="button"
                            >
                                <span className="event-row-main">
                                    <strong>{eventLabel(event.type, terms)}</strong>
                                    <small>{terms.eventStarted} {formatEventOffset(event.startMs, terms)}</small>
                                </span>
                                <span className={`severity-badge ${event.severity}`}>
                                    {severityLabel(event.severity, terms)}
                                </span>
                                <span className="event-duration">{formatDuration(event.durationMs, terms)}</span>
                            </button>
                        ))}
                    </div>

                    {selectedEvent && (
                        <div className="event-details">
                            <div className="event-details-head">
                                <div>
                                    <span>{terms.eventSegment}</span>
                                    <strong>{localizedLabel(selectedEvent.segment, terms.segments)}</strong>
                                </div>
                                <div>
                                    <span>{terms.eventDuration}</span>
                                    <strong>{formatDuration(selectedEvent.durationMs, terms)}</strong>
                                </div>
                                <div>
                                    <span>{terms.severityLabel}</span>
                                    <strong>{severityLabel(selectedEvent.severity, terms)}</strong>
                                </div>
                            </div>

                            <div className="event-section">
                                <h3>{terms.eventEvidence}</h3>
                                <div className="evidence-grid">
                                    {Object.entries(selectedEvent.evidence || {}).map(([key, value]) => (
                                        <div key={key}>
                                            <span>{localizedLabel(key, terms.evidenceLabels)}</span>
                                            <strong>{formatEvidenceValueForKey(key, value, terms)}</strong>
                                        </div>
                                    ))}
                                </div>
                            </div>

                            <div className="event-section">
                                <h3>{terms.eventSuggestions}</h3>
                                <div className="suggestion-list">
                                    {(selectedEvent.suggestedActions || []).map(action => (
                                        <div key={`${action.priority}-${action.category}-${action.item}`} className="suggestion">
                                            <strong>{localizedLabel(action.category, terms.actionCategories)} / {localizedLabel(action.item, terms.actionItems)}</strong>
                                            <span>{localizedLabel(action.direction, terms.actionDirections)} {localizedLabel(action.amount, terms.actionAmounts)}</span>
                                            <small>{localizedLabel(action.reason, terms.actionReasons)}</small>
                                        </div>
                                    ))}
                                </div>
                            </div>
                        </div>
                    )}
                </div>
            )}
        </section>
    );
}

type DiagnosticRow = {
    name: string;
    value: string;
    unit: string;
    source: string;
    range: string;
    warn?: boolean;
};

const DIAGNOSTIC_FIELD_LABELS: Record<string, Record<Lang, string>> = {
    mode: {en: 'Mode', zh: '模式'},
    gameMode: {en: 'Game mode', zh: '游戏模式'},
    validPackets: {en: 'Valid packets', zh: '有效包'},
    parseErrors: {en: 'Parse errors', zh: '解析错误'},
    replayProgress: {en: 'Replay progress', zh: '回放进度'},
    speedKmh: {en: 'Speed', zh: '速度'},
    speedFieldKmh: {en: 'Packet speed', zh: '包内速度'},
    velocitySpeedKmh: {en: 'Velocity speed', zh: '向量速度'},
    speedDifference: {en: 'Speed difference', zh: '速度差异'},
    carOrdinal: {en: 'Vehicle ID', zh: '车辆 ID'},
    carCategory: {en: 'Vehicle category', zh: '车辆分类'},
    carClassPi: {en: 'Class / PI', zh: '等级 / PI'},
    drivetrainCylinders: {en: 'Drivetrain / cylinders', zh: '传动 / 气缸'},
    rpm: {en: 'RPM', zh: '转速'},
    rpmRatio: {en: 'RPM ratio', zh: '转速比例'},
    engineIdleRpm: {en: 'Idle RPM', zh: '怠速转速'},
    engineMaxRpm: {en: 'Max RPM', zh: '最高转速'},
    power: {en: 'Power', zh: '功率'},
    torque: {en: 'Torque', zh: '扭矩'},
    boostFuel: {en: 'Boost / fuel', zh: '增压 / 燃油'},
    accelerationXYZ: {en: 'Acceleration XYZ', zh: '加速度 XYZ'},
    velocityXYZ: {en: 'Velocity XYZ', zh: '速度向量 XYZ'},
    yawPitchRoll: {en: 'Yaw / pitch / roll', zh: '偏航 / 俯仰 / 侧倾'},
    yawRatePitchRateRoll: {en: 'Yaw / pitch / roll rate', zh: '偏航 / 俯仰 / 侧倾角速度'},
    positionXYZ: {en: 'Position XYZ', zh: '位置 XYZ'},
    throttle: {en: 'Throttle', zh: '油门'},
    brake: {en: 'Brake', zh: '刹车'},
    clutch: {en: 'Clutch', zh: '离合'},
    handBrake: {en: 'Handbrake', zh: '手刹'},
    steer: {en: 'Steering', zh: '方向'},
    isRaceOn: {en: 'Race mode flag', zh: '比赛模式标记'},
    lapRacePosition: {en: 'Lap / race position', zh: '圈数 / 名次'},
    lapTimes: {en: 'Best / last / current lap', zh: '最佳 / 上圈 / 当前圈'},
    currentRaceTime: {en: 'Race time', zh: '比赛时间'},
    distanceTraveled: {en: 'Distance traveled', zh: '行驶距离'},
    drivingLine: {en: 'Driving line', zh: '辅助线'},
    aiBrakeDifference: {en: 'AI brake difference', zh: 'AI 刹车差值'},
    smashableVelDiff: {en: 'Smashable velocity diff', zh: '可破坏物速度差'},
    smashableMass: {en: 'Smashable mass', zh: '可破坏物质量'},
};

const WHEEL_DIAGNOSTIC_CORNERS: Record<string, Record<Lang, string>> = {
    FL: {en: 'Front left', zh: '左前轮'},
    FR: {en: 'Front right', zh: '右前轮'},
    RL: {en: 'Rear left', zh: '左后轮'},
    RR: {en: 'Rear right', zh: '右后轮'},
};

const WHEEL_DIAGNOSTIC_FIELDS: Record<string, Record<Lang, string>> = {
    slip: {en: 'slip ratio / angle / combined', zh: '滑移率 / 滑移角 / 综合滑移'},
    temp: {en: 'tire temperature', zh: '胎温'},
    suspension: {en: 'suspension travel', zh: '悬挂行程'},
    surface: {en: 'surface signals', zh: '路面信号'},
};

function developerFieldGroups(current: TelemetryFrame | null, status: TelemetryStatus, replayStatus: TelemetryReplayStatus, terms: Copy, language: Lang) {
    const speedDiff = current ? Math.abs(current.speedFieldKmh - current.velocitySpeedKmh) : 0;
    const speedBase = current ? Math.max(Math.abs(current.speedKmh), 1) : 1;
    const speedDiffPct = (speedDiff / speedBase) * 100;
    const gameMode = current?.gameMode || 'unknown';
    const raceFieldsApplicable = gameMode === 'race';
    const raceFieldRange = raceFieldsApplicable ? '>= 0' : terms.notApplicable;
    const f = (value: number | undefined, digits = 2) => formatNumber(value, digits);
    const p = (value: number | undefined) => formatPercentValue(value);
    const d = (name: string, value: string, unit: string, source: string, range: string, warn = false) => diag(language, name, value, unit, source, range, warn);
    return [
        {
            title: terms.mode,
            rows: [
                d('mode', modeLabel(status.mode, terms), '', 'TelemetryStatus.mode', 'idle/udp/replay'),
                d('gameMode', gameModeLabel(gameMode, terms), '', 'backend GameModeTracker', 'menu/free roam/race'),
                d('validPackets', String(status.validPackets), '', 'TelemetryStatus.validPackets', '>= 0'),
                d('parseErrors', String(status.parseErrors), '', 'TelemetryStatus.parseErrors', '0 preferred', status.parseErrors > 0),
                d('replayProgress', `${Math.round(replayStatus.progress01 * 100)}%`, '%', 'TelemetryReplayStatus.progress01', '0-100'),
            ],
        },
        {
            title: terms.speedCalibration,
            rows: [
                d('speedKmh', f(current?.speedKmh, 1), 'km/h', speedSourceRawField(current?.speedSource), '0-700'),
                d('speedFieldKmh', f(current?.speedFieldKmh, 1), 'km/h', 'Speed', '0-700'),
                d('velocitySpeedKmh', f(current?.velocitySpeedKmh, 1), 'km/h', 'VelocityX/VelocityY/VelocityZ', '0-700'),
                d('speedDifference', current ? `${speedDiff.toFixed(1)} / ${speedDiffPct.toFixed(0)}%` : '--', 'km/h', 'speedFieldKmh vs velocitySpeedKmh', '< 25 km/h or < 20%', current ? speedDiff > 25 && speedDiffPct > 20 : false),
            ],
        },
        {
            title: terms.vehicleMetadata,
            rows: [
                d('carOrdinal', formatOptionalInt(current?.carOrdinal), '', 'CarOrdinal', '> 0'),
                d('carCategory', current ? formatCategory(current) : '--', '', 'CarCategory', 'mapped'),
                d('carClassPi', current ? `${current.carClass || '--'} / ${formatOptionalInt(current.carPi)}` : '--', '', 'CarClass/CarPI', 'D-X / 100-999'),
                d('drivetrainCylinders', current ? `${current.drivetrain || '--'} / ${formatOptionalInt(current.numCylinders)}` : '--', '', 'DrivetrainType/NumCylinders', 'FWD/RWD/AWD'),
            ],
        },
        {
            title: terms.enginePower,
            rows: [
                d('rpm', f(current?.rpm, 0), 'rpm', 'CurrentEngineRpm', 'idle-max'),
                d('rpmRatio', p(current?.rpmRatio), '%', 'CurrentEngineRpm/EngineIdleRpm/EngineMaxRpm', '0-100%'),
                d('engineIdleRpm', f(current?.engineIdleRpm, 0), 'rpm', 'EngineIdleRpm', '> 0'),
                d('engineMaxRpm', f(current?.engineMaxRpm, 0), 'rpm', 'EngineMaxRpm', '> idle'),
                d('power', f(current?.power, 0), 'W', 'Power', 'varies'),
                d('torque', f(current?.torque, 0), 'Nm', 'Torque', 'varies'),
                d('boostFuel', current ? `${f(current.boost, 2)} / ${p(current.fuel)}` : '--', '', 'Boost/Fuel', 'fuel 0-100%'),
            ],
        },
        {
            title: terms.motionPose,
            rows: [
                d('accelerationXYZ', current ? `${f(current.accelerationX)} / ${f(current.accelerationY)} / ${f(current.accelerationZ)}` : '--', 'm/s2', 'AccelerationX/AccelerationY/AccelerationZ', 'varies'),
                d('velocityXYZ', current ? `${f(current.velocityX)} / ${f(current.velocityY)} / ${f(current.velocityZ)}` : '--', 'm/s', 'VelocityX/VelocityY/VelocityZ', 'varies'),
                d('yawPitchRoll', current ? `${f(current.yaw)} / ${f(current.pitch)} / ${f(current.roll)}` : '--', 'rad', 'Yaw/Pitch/Roll', 'varies'),
                d('yawRatePitchRateRoll', current ? `${f(current.yawRate)} / ${f(current.pitchRate)} / ${f(current.rollRate)}` : '--', 'rad/s', 'AngularVelocityX/Y/Z', 'varies'),
                d('positionXYZ', current ? `${f(current.positionX, 0)} / ${f(current.positionY, 0)} / ${f(current.positionZ, 0)}` : '--', 'm', 'PositionX/PositionY/PositionZ', 'world'),
            ],
        },
        {
            title: terms.driverInputs,
            rows: [
                d('throttle', p(current?.throttle01), '%', 'Accel', '0-100%'),
                d('brake', p(current?.brake01), '%', 'Brake', '0-100%'),
                d('clutch', p(current?.clutch01), '%', 'Clutch', '0-100%'),
                d('handBrake', p(current?.handBrake01), '%', 'HandBrake', '0-100%'),
                d('steer', f(current?.steer01, 2), '', 'Steer', '-1 to 1'),
            ],
        },
        {
            title: terms.wheelState,
            rows: wheelDiagnostics(current, terms, language),
        },
        {
            title: terms.raceLapData,
            rows: [
                d('isRaceOn', current ? String(current.isRaceOn) : '--', '', 'IsRaceOn', 'mode hint only'),
                d('lapRacePosition', current ? `${formatOptionalInt(current.lapNumber, true)} / ${formatOptionalInt(current.racePosition)}` : '--', '', 'LapNumber/RacePosition', raceFieldRange),
                d('lapTimes', current ? `${f(current.bestLap, 2)} / ${f(current.lastLap, 2)} / ${f(current.currentLap, 2)}` : '--', 's', 'BestLap/LastLap/CurrentLap', raceFieldRange),
                d('currentRaceTime', f(current?.currentRaceTime, 2), 's', 'CurrentRaceTime', raceFieldRange),
                d('distanceTraveled', f(current?.distanceTraveled, 1), 'm', 'DistanceTraveled', raceFieldRange),
            ],
        },
        {
            title: terms.auxiliaryFields,
            rows: [
                d('drivingLine', f(current?.drivingLine01, 2), '', 'NormalizedDrivingLine', '-1 to 1'),
                d('aiBrakeDifference', f(current?.aiBrakeDifference01, 2), '', 'NormalizedAIBrakeDifference', '-1 to 1'),
                d('smashableVelDiff', f(current?.smashableVelDiff, 2), '', 'SmashableVelDiff', 'varies'),
                d('smashableMass', f(current?.smashableMass, 2), '', 'SmashableMass', 'varies'),
            ],
        },
    ];
}

function wheelDiagnostics(current: TelemetryFrame | null, terms: Copy, language: Lang) {
    const wheels = [
        ['FL', current?.wheelFL],
        ['FR', current?.wheelFR],
        ['RL', current?.wheelRL],
        ['RR', current?.wheelRR],
    ] as const;
    const d = (name: string, value: string, unit: string, source: string, range: string, warn = false) => diag(language, name, value, unit, source, range, warn);
    return wheels.flatMap(([prefix, wheel]) => [
        d(`${prefix}_slip`, wheel ? `${wheel.slipRatio.toFixed(2)} / ${wheel.slipAngle.toFixed(2)} / ${wheel.combinedSlip.toFixed(2)}` : '--', '', `TireSlipRatio${prefix}/TireSlipAngle${prefix}/TireCombinedSlip${prefix}`, 'combined < 1 preferred', !!wheel && wheel.combinedSlip > 1),
        d(`${prefix}_temp`, wheel ? wheel.tireTemp.toFixed(0) : '--', 'deg', `TireTemp${prefix}`, 'surface dependent'),
        d(`${prefix}_suspension`, wheel ? `${(wheel.suspensionTravel * 100).toFixed(0)}% / ${wheel.suspensionTravelMeters.toFixed(3)}m` : '--', '', `SuspensionTravel${prefix}/SuspensionTravelMeters${prefix}`, '< 95% preferred', !!wheel && wheel.suspensionTravel > 0.95),
        d(`${prefix}_surface`, wheel ? `${wheel.rumbleStrip.toFixed(0)} / ${wheel.puddleDepth.toFixed(2)} / ${wheel.surfaceRumble.toFixed(2)}` : '--', '', `WheelOnRumbleStrip${prefix}/WheelInPuddleDepth${prefix}/SurfaceRumble${prefix}`, 'varies'),
    ]);
}

function diag(language: Lang, name: string, value: string, unit: string, source: string, range: string, warn = false): DiagnosticRow {
    return {name: diagnosticFieldLabel(name, language), value, unit, source, range, warn};
}

function diagnosticFieldLabel(name: string, language: Lang) {
    const wheelMatch = name.match(/^(FL|FR|RL|RR)_(slip|temp|suspension|surface)$/);
    if (wheelMatch) {
        const corner = WHEEL_DIAGNOSTIC_CORNERS[wheelMatch[1]]?.[language] || wheelMatch[1];
        const field = WHEEL_DIAGNOSTIC_FIELDS[wheelMatch[2]]?.[language] || wheelMatch[2];
        return `${corner} ${field}`;
    }
    return DIAGNOSTIC_FIELD_LABELS[name]?.[language] || name;
}

function speedSourceRawField(value: string | undefined) {
    if (value === 'velocity') {
        return 'VelocityX/VelocityY/VelocityZ';
    }
    if (value === 'packet') {
        return 'Speed';
    }
    return '--';
}

function sparklinePoints(values: number[]) {
    if (values.length === 0) {
        return '0,70 300,70';
    }
    const width = 300;
    const height = 80;
    const min = Math.min(...values);
    const max = Math.max(...values);
    const range = Math.max(max - min, 1);
    return values.map((value, index) => {
        const x = values.length === 1 ? width : (index / (values.length - 1)) * width;
        const y = height - 10 - ((value - min) / range) * (height - 20);
        return `${x.toFixed(1)},${y.toFixed(1)}`;
    }).join(' ');
}

function wheelHealth(combinedSlip: number) {
    if (combinedSlip >= 1) {
        return 'danger';
    }
    if (combinedSlip >= 0.75) {
        return 'warning';
    }
    return 'stable';
}

function wheelGripCss(state: string | undefined) {
    if (state === 'limit') {
        return 'danger';
    }
    if (state === 'warning') {
        return 'warning';
    }
    return 'stable';
}

function tireWheelPositionLabel(position: string, terms: Copy) {
    if (position === 'front_left') {
        return terms.frontLeft;
    }
    if (position === 'front_right') {
        return terms.frontRight;
    }
    if (position === 'rear_left') {
        return terms.rearLeft;
    }
    if (position === 'rear_right') {
        return terms.rearRight;
    }
    return formatEvidenceKey(position);
}

function formatNumber(value: number | undefined, digits: number) {
    if (value === undefined || Number.isNaN(value)) {
        return '--';
    }
    return value.toFixed(digits);
}

function formatTirePhaseSpeedReference(value: number | undefined, confidence: number | undefined, terms: Copy) {
    const reference = value !== undefined && Number.isFinite(value) && value > 0 ? `${formatNumber(value, 0)} km/h` : '--';
    const confidenceLabel = (confidence || 0) >= 0.5 ? terms.quickConfidenceLabels.high : terms.quickConfidenceLabels.low;
    return `${reference} / ${confidenceLabel}`;
}

function formatTirePhaseSpeedBand(value: number | undefined, confidence: number | undefined, terms: Copy) {
    const rounded = Math.round(value || 0);
    const labels = terms === COPY.zh
        ? {1: '低速', 2: '中速', 3: '高速'}
        : {1: 'Low speed', 2: 'Mid speed', 3: 'High speed'};
    const band = labels[rounded as 1 | 2 | 3] || '--';
    const confidenceText = (confidence || 0) >= 0.5
        ? (terms === COPY.zh ? '动态基准' : 'dynamic reference')
        : (terms === COPY.zh ? '固定回退' : 'fixed fallback');
    return `${band} / ${confidenceText}`;
}

function emptyTireRegressionExpectation(): TireRegressionExpectation {
    return {
        allowedPhases: [],
        requiredGripTypes: [],
        allowedAxles: [],
        forbiddenGripTypes: [],
        minDataQuality: '',
        notes: '',
    };
}

function expectationToForm(expected?: TireRegressionExpectation | null): TireRegressionExpectedFormState {
    return {
        allowedPhases: listToCsv(expected?.allowedPhases || []),
        requiredGripTypes: listToCsv(expected?.requiredGripTypes || []),
        allowedAxles: listToCsv(expected?.allowedAxles || []),
        forbiddenGripTypes: listToCsv(expected?.forbiddenGripTypes || []),
        minDataQuality: expected?.minDataQuality || 'low_confidence',
        notes: expected?.notes || '',
    };
}

function expectationFromForm(form: TireRegressionExpectedFormState): TireRegressionExpectation {
    return {
        allowedPhases: csvToList(form.allowedPhases),
        requiredGripTypes: csvToList(form.requiredGripTypes),
        allowedAxles: csvToList(form.allowedAxles),
        forbiddenGripTypes: csvToList(form.forbiddenGripTypes),
        minDataQuality: form.minDataQuality || 'low_confidence',
        notes: form.notes || '',
    };
}

function csvToList(value: string) {
    return value
        .split(',')
        .map(item => item.trim())
        .filter(Boolean);
}

function listToCsv(values: string[]) {
    return (values || []).join(', ');
}

function vehicleSnapshotLabel(vehicle: QuickVehicleSnapshot | undefined, terms: Copy) {
    if (!vehicle) {
        return '--';
    }
    const parts = [
        vehicle.carOrdinal ? `ID ${vehicle.carOrdinal}` : '',
        vehicle.carClass || '',
        vehicle.carPi ? `PI ${vehicle.carPi}` : '',
        vehicle.drivetrain || '',
    ].filter(Boolean);
    return parts.join(' / ') || terms.noProfile;
}

function snapshotShortLabel(snapshot: TireDiagnosticSnapshot, terms: Copy) {
    if (!snapshot) {
        return '--';
    }
    const phase = localizedLabel(snapshot.phase?.current || 'unknown', terms.tireLabPhaseLabels);
    const grip = localizedLabel(snapshot.gripLimit?.type || 'no_limit_detected', terms.tireGripLimitLabels);
    const axle = localizedLabel(snapshot.gripLimit?.limitedAxle || 'none', terms.tireAxleLabels);
    return `${phase} / ${grip} / ${axle}`;
}

function formatOptionalInt(value: number | undefined, allowZero = false) {
    if (value === undefined || !Number.isFinite(value) || (!allowZero && value <= 0) || (allowZero && value < 0)) {
        return '--';
    }
    return Math.round(value).toString();
}

function formatPercentValue(value: number | undefined) {
    if (value === undefined || !Number.isFinite(value)) {
        return '--';
    }
    return `${Math.round(clamp(value, 0, 1) * 100)}%`;
}

function formatRawPercent(value: number | undefined, digits = 0) {
    if (value === undefined || !Number.isFinite(value)) {
        return '--';
    }
    return `${value.toFixed(digits)}%`;
}

function formatSignedPercent(value: number | null | undefined) {
    if (value === undefined || value === null || !Number.isFinite(value)) {
        return '--';
    }
    return `${value > 0 ? '+' : ''}${value.toFixed(1)}%`;
}

function formatGValue(value: number | null | undefined) {
    if (value === undefined || value === null || !Number.isFinite(value)) {
        return '--';
    }
    return `${value >= 0 ? '+' : ''}${value.toFixed(2)} g`;
}

function formatMeters(value: number | null | undefined) {
    if (value === undefined || value === null || !Number.isFinite(value)) {
        return '--';
    }
    return `${value.toFixed(value >= 100 ? 0 : 1)} m`;
}

function formatSecondsValue(value: number | null | undefined) {
    if (value === undefined || value === null || !Number.isFinite(value)) {
        return '--';
    }
    return `${value.toFixed(2)} s`;
}

function localizedWarningFlags(value: string | undefined, terms: Copy) {
    return String(value || '')
        .split(',')
        .map(flag => flag.trim())
        .filter(Boolean)
        .map(flag => terms.warningLabels[flag as keyof Copy['warningLabels']] || flag);
}

function formatCategory(frame: TelemetryFrame | null) {
    if (!frame) {
        return '--';
    }
    const id = formatOptionalInt(frame.carCategory, true);
    return frame.carCategoryName ? `${frame.carCategoryName} (${id})` : id;
}

function speedSourceLabel(value: string | undefined, terms: Copy) {
    if (value === 'velocity') {
        return terms.speedSourceVelocity;
    }
    if (value === 'packet') {
        return terms.speedSourcePacket;
    }
    return terms.speedSourceNone;
}

function modeLabel(value: string | undefined, terms: Copy) {
    if (value === 'udp') {
        return terms.udpMode;
    }
    if (value === 'replay') {
        return terms.replayMode;
    }
    return terms.idleMode;
}

function gameModeLabel(mode: string | undefined, terms: Copy) {
    const normalized = normalizeGameMode(mode);
    if (normalized === 'race') {
        return terms.raceDataActive;
    }
    if (normalized === 'free_roam') {
        return terms.freeRoamMode;
    }
    if (normalized === 'menu') {
        return terms.menuMode;
    }
    if (normalized === 'mixed') {
        return terms.mixedMode;
    }
    return terms.unknownMode;
}

function normalizeGameMode(mode: string | undefined): GameMode {
    if (mode === 'race' || mode === 'free_roam' || mode === 'menu' || mode === 'mixed') {
        return mode;
    }
    return 'unknown';
}

function hasVehicleTelemetry(frame: TelemetryFrame) {
    return frame.carOrdinal > 0 && Number.isFinite(frame.positionX) && Number.isFinite(frame.positionZ);
}

function recordingLabel(status: TelemetryStatus, terms: Copy) {
    if (status.recordingTruncated) {
        return terms.recordingTruncated;
    }
    if (status.recordingActive) {
        return terms.recording;
    }
    if (status.recordingPackets > 0) {
        return terms.recordingReady;
    }
    return terms.noRecording;
}

function formatBytes(value: number | undefined) {
    const bytes = value || 0;
    if (bytes >= 1024 * 1024) {
        return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
    }
    if (bytes >= 1024) {
        return `${(bytes / 1024).toFixed(1)} KB`;
    }
    return `${bytes} B`;
}

function formatTime(value: string, fallback: string) {
    if (!value) {
        return fallback;
    }
    return new Date(value).toLocaleTimeString();
}

function eventLabel(type: string, terms: Copy) {
    return localizedLabel(type, terms.events);
}

function severityLabel(value: string, terms: Copy) {
    if (value === 'high') {
        return terms.severityHigh;
    }
    if (value === 'medium') {
        return terms.severityMedium;
    }
    return terms.severityLow;
}

function formatDuration(ms: number, terms: Copy) {
    if (!Number.isFinite(ms) || ms <= 0) {
        return `0 ${terms.durationMsUnit}`;
    }
    if (ms < 1000) {
        return `${Math.round(ms)} ${terms.durationMsUnit}`;
    }
    return `${(ms / 1000).toFixed(1)} ${terms.durationSecondUnit}`;
}

function formatEventOffset(ms: number, terms: Copy) {
    return `+${formatDuration(ms, terms)}`;
}

function localizedLabel(value: string, labels: Record<string, string> = {}) {
    return labels[value] || formatEvidenceKey(value);
}

function formatScore(value: number | undefined) {
    if (value === undefined || !Number.isFinite(value)) {
        return '--';
    }
    return `${Math.round(clamp(value, 0, 100))}/100`;
}

function roadVerdictLabel(value: string | undefined, terms: Copy) {
    return localizedLabel(value || 'insufficient_data', terms.roadVerdicts);
}

function roadBaselineStatusLabel(value: string | undefined, terms: Copy) {
    return localizedLabel(value || 'no_valid_standard_run', terms.roadBaselineStatuses);
}

function roadAttributionLabel(value: string | undefined, terms: Copy) {
    return localizedLabel(value || 'data_gap', terms.roadAttributions);
}

function roadAttributionMessage(attribution: RoadEvaluationAttribution, terms: Copy) {
    return localizedLabel(attribution.message || attribution.eventType || attribution.type, terms.roadAttributionMessages);
}

function issueFamilyLabel(value: string | undefined, terms: Copy) {
    const labels: Record<string, string> = {
        launch_traction: terms === COPY.zh ? '起步牵引' : 'Launch traction',
        gearing_acceleration: terms === COPY.zh ? '齿比 / 加速' : 'Gearing / acceleration',
        brake_balance: terms === COPY.zh ? '刹车平衡' : 'Brake balance',
        corner_entry_balance: terms === COPY.zh ? '入弯平衡' : 'Corner entry balance',
        mid_corner_balance: terms === COPY.zh ? '持续过弯平衡' : 'Sustained cornering balance',
        corner_exit_power: terms === COPY.zh ? '出弯动力' : 'Corner exit power',
        suspension_platform: terms === COPY.zh ? '悬挂平台' : 'Suspension platform',
        tire_temperature_stability: terms === COPY.zh ? '轮胎温度 / 稳定性' : 'Tire temperature / stability',
        driver_execution: terms === COPY.zh ? '驾驶执行' : 'Driver execution',
    };
    return labels[value || ''] || formatEvidenceKey(value || '');
}

function issueComparisonLabel(value: string | undefined, terms: Copy) {
    const zh = terms === COPY.zh;
    switch (value) {
        case 'improved':
            return zh ? '改善' : 'Improved';
        case 'worsened':
            return zh ? '恶化' : 'Worsened';
        case 'unchanged':
            return zh ? '持平' : 'Unchanged';
        default:
            return zh ? '无法比较' : 'Unavailable';
    }
}

function strategyRecommendationLabel(value: string | undefined, terms: Copy) {
    const zh = terms === COPY.zh;
    switch (value) {
        case 'increase_adjustment_step_not_detection_threshold':
            return zh ? '重复问题确认，建议增加调校步幅，不提高检测阈值' : 'Repeated issue confirmed; increase tuning step, not detection threshold';
        case 'consider_raising_threshold_or_duration':
            return zh ? '疑似过度匹配，建议提高阈值或持续时间' : 'Possible overmatch; raise threshold or duration';
        case 'consider_lowering_threshold_if_driver_feedback_confirms':
            return zh ? '若遥测样本确认，可考虑降低阈值' : 'Consider lowering threshold if telemetry samples confirm';
        case 'repeated_issue_confirmed':
            return zh ? '重复问题已确认' : 'Repeated issue confirmed';
        default:
            return zh ? '保持当前阈值' : 'Keep current thresholds';
    }
}

function strategyHintLabel(hint: StrategyAnalysisHint, terms: Copy) {
    const subject = hint.family ? issueFamilyLabel(hint.family, terms) : hint.eventType ? eventLabel(hint.eventType, terms) : '';
    return `${subject ? `${subject}: ` : ''}${strategyRecommendationLabel(hint.message, terms)}`;
}

function issueGroupRepresentativeEvent(group: SessionIssueGroup): DetectedEvent {
    const first = (group.events || [])[0];
    if (first) {
        return first;
    }
    const evidence = Object.fromEntries(Object.entries(group.evidence || {}).map(([key, stat]) => [key, stat.avg]));
    return {
        id: group.id,
        type: group.eventTypes?.[0] || group.family,
        severity: group.severity,
        startMs: group.firstStartMs,
        endMs: group.lastEndMs,
        durationMs: group.totalDurationMs,
        segment: group.segment,
        evidence,
        suggestedActions: group.primaryActions || [],
    };
}

function profileFieldDisplayName(key: string, language: Lang) {
    const field = profileFields.find(item => String(item.key) === key);
    return field ? profileFieldLabel(field, language) : formatEvidenceKey(key);
}

function profileFieldOrder(key: string) {
    const coreIndex = coreTuneFieldOrder.findIndex(item => String(item) === key);
    if (coreIndex >= 0) {
        return coreIndex;
    }
    const index = profileFields.findIndex(item => String(item.key) === key);
    return index >= 0 ? coreTuneFieldOrder.length + index : Number.MAX_SAFE_INTEGER;
}

function reportStatusAlerts(session: TelemetrySession, evaluation: RoadSessionEvaluation | null, terms: Copy) {
    const alerts: Array<{ key: string; title: string; message: string; tone: 'ok' | 'warn' }> = [];
    if (session.tuneProfileId) {
        alerts.push({key: 'profile-bound', title: terms.profileBoundStatus, message: session.tuneName || terms.noProfile, tone: 'ok'});
    } else {
        alerts.push({key: 'profile-unbound', title: terms.profileUnboundStatus, message: terms.sessionProfileUnboundHint, tone: 'warn'});
    }
    if (sessionTestConditions(session).driverMode === 'unknown' || (session.driverModeConfidence || 0) < 0.5) {
        alerts.push({key: 'driver-unknown', title: terms.driverMode, message: terms.driverModeUnknownReportHint, tone: 'warn'});
    }
    if (!evaluation?.bestRun) {
        alerts.push({key: 'no-standard-run', title: terms.benchmarkRuns, message: terms.standardSegmentMissingHint, tone: 'warn'});
    } else if (evaluation.baselineStatus === 'missing_auto_baseline') {
        alerts.push({key: 'missing-baseline', title: terms.autoBaseline, message: terms.baselineMissingReportHint, tone: 'warn'});
    }
    return alerts;
}

function localizedUseCase(value: string | undefined, terms: Copy) {
    if (!value) {
        return '';
    }
    return terms.useCases[value as keyof Copy['useCases']] || value;
}

function normalizeTestConditions(input: Partial<TestConditions> | null | undefined): TestConditions {
    const safe = input || {};
    return {
        driverMode: normalizeOption(safe.driverMode, testConditionOptions.driverMode),
        brakeAssist: normalizeOption(safe.brakeAssist, testConditionOptions.brakeAssist),
        steeringAssist: normalizeOption(safe.steeringAssist, testConditionOptions.steeringAssist),
        tractionControl: normalizeOption(safe.tractionControl, testConditionOptions.tractionControl),
        stabilityControl: normalizeOption(safe.stabilityControl, testConditionOptions.stabilityControl),
        shifting: normalizeOption(safe.shifting, testConditionOptions.shifting),
        launchControl: normalizeOption(safe.launchControl, testConditionOptions.launchControl),
    };
}

function normalizeOption(value: string | undefined, options: string[]) {
    const normalized = String(value || '').trim().toLowerCase();
    return options.includes(normalized) ? normalized : 'unknown';
}

function sessionTestConditions(session: TelemetrySession): TestConditions {
    return normalizeTestConditions({
        driverMode: session.driverMode,
        brakeAssist: session.brakeAssist,
        steeringAssist: session.steeringAssist,
        tractionControl: session.tractionControl,
        stabilityControl: session.stabilityControl,
        shifting: session.shifting,
        launchControl: session.launchControl,
    });
}

function testConditionLabel(value: string | undefined, terms: Copy) {
    const key = String(value || 'unknown') as keyof Copy['testConditionValues'];
    return terms.testConditionValues[key] || terms.testConditionValues.unknown;
}

function formatDriverModeDetection(session: TelemetrySession, terms: Copy) {
    const label = testConditionLabel(session.driverMode, terms);
    const confidence = Number(session.driverModeConfidence || 0);
    return confidence > 0 ? `${label} / ${(confidence * 100).toFixed(0)}%` : label;
}

function formatTestConditionsCompact(conditions: Partial<TestConditions> | null | undefined, terms: Copy) {
    const normalized = normalizeTestConditions(conditions);
    return [
        `${terms.brakeAssist}: ${testConditionLabel(normalized.brakeAssist, terms)}`,
        `TCS: ${testConditionLabel(normalized.tractionControl, terms)}`,
        `STM: ${testConditionLabel(normalized.stabilityControl, terms)}`,
        `${terms.shifting}: ${testConditionLabel(normalized.shifting, terms)}`,
    ].join(' / ');
}

function formatAssistSummary(session: TelemetrySession, terms: Copy) {
    const conditions = sessionTestConditions(session);
    return [
        `${terms.brakeAssist}: ${testConditionLabel(conditions.brakeAssist, terms)}`,
        `${terms.steeringAssist}: ${testConditionLabel(conditions.steeringAssist, terms)}`,
        `TCS: ${testConditionLabel(conditions.tractionControl, terms)}`,
        `STM: ${testConditionLabel(conditions.stabilityControl, terms)}`,
        `${terms.shifting}: ${testConditionLabel(conditions.shifting, terms)}`,
        `${terms.launchControl}: ${testConditionLabel(conditions.launchControl, terms)}`,
    ].join(' / ');
}

function comparabilityWarningLabel(value: string, terms: Copy) {
    return terms.comparabilityWarnings[value as keyof Copy['comparabilityWarnings']] || value;
}

function telemetryPoint(frame: TelemetryFrame): TrackCapturePoint {
    return {
        x: frame.positionX,
        y: frame.positionY,
        z: frame.positionZ,
        lapNumber: frame.lapNumber || 0,
        currentLap: frame.currentLap || 0,
        currentRaceTime: frame.currentRaceTime || 0,
    };
}

function capturedTrackInput(
    capture: TrackCaptureState,
    sourceMode: string,
    trackType: BenchmarkTrackType,
    extractionMode: BenchmarkExtractionMode,
    gateWidth: number,
    gateDepth: number,
    startGate?: BenchmarkGate,
    finishGate?: BenchmarkGate,
): BenchmarkTrackInput {
    const lapSegment = trackType !== 'sprint' && extractionMode !== 'full_segment'
        ? capturedCircuitSegmentFromLapIncrements(capture.points, extractionMode)
        : null;
    const inferredTrackType = trackType === 'auto'
        ? lapSegment ? 'circuit' : inferClientTrackType(capture.points)
        : trackType;
    const points = inferredTrackType === 'circuit' && lapSegment ? lapSegment.points : capture.points;
    const first = points[0];
    const last = points[points.length - 1];
    const direction = routeDirection(points);
    const simplified = simplifyClientPoints(points, 8, 800);
    const normalizedStartGate = startGate ? applyGateSize(startGate, gateWidth, gateDepth) : clientGate(first, direction.x, direction.z, gateWidth, gateDepth);
    const normalizedFinishGate = inferredTrackType === 'circuit'
        ? normalizedStartGate
        : finishGate ? applyGateSize(finishGate, gateWidth, gateDepth) : clientGate(last, directionAtClientEnd(points).x, directionAtClientEnd(points).z, gateWidth, gateDepth);
    return {
        name: capture.name.trim() || `Track ${new Date().toLocaleString()}`,
        sourceMode,
        trackType: inferredTrackType,
        start: first,
        end: inferredTrackType === 'circuit' ? first : last,
        startRadius: 20,
        endRadius: 20,
        directionX: direction.x,
        directionZ: direction.z,
        startGate: normalizedStartGate,
        finishGate: normalizedFinishGate,
        checkpoints: clientCheckpoints(simplified),
        routeLengthMeters: pointsRouteLength(points),
        hasDrivingLine: capture.hasDrivingLine,
        polyline: simplified,
        lapCountObserved: lapSegment?.lapCountObserved || countCapturedLapIncrements(capture.points),
        notes: '',
    };
}

function capturedCircuitSegmentFromLapIncrements(points: TrackCapturePoint[], extractionMode: BenchmarkExtractionMode) {
    const boundaries = capturedLapBoundaries(points);
    if (boundaries.length < 2) {
        return null;
    }
    const segments: {points: TrackCapturePoint[]; durationMs: number}[] = [];
    for (let index = 1; index < boundaries.length; index++) {
        const start = boundaries[index - 1];
        const end = boundaries[index];
        if (end - start < 2) {
            continue;
        }
        const segment = points.slice(start, end + 1);
        if (pointsRouteLength(segment) < 200) {
            continue;
        }
        const durationMs = capturedRaceDurationMs(segment[0], segment[segment.length - 1]);
        segments.push({points: segment, durationMs});
    }
    if (!segments.length) {
        return null;
    }
    const selected = extractionMode === 'auto_best_lap'
        ? segments.reduce((best, segment) => {
            if (segment.durationMs > 0 && (!best.durationMs || segment.durationMs < best.durationMs)) {
                return segment;
            }
            return best;
        }, segments[0])
        : segments[0];
    return {points: selected.points, lapCountObserved: segments.length};
}

function capturedLapBoundaries(points: TrackCapturePoint[]) {
    const boundaries: number[] = [];
    let lastLap = validCapturedLap(points[0]);
    for (let index = 1; index < points.length; index++) {
        const lap = validCapturedLap(points[index]);
        if (lap !== null && lastLap !== null && lap > lastLap) {
            boundaries.push(index);
        }
        if (lap !== null) {
            lastLap = lap;
        }
    }
    return boundaries;
}

function countCapturedLapIncrements(points: TrackCapturePoint[]) {
    return capturedLapBoundaries(points).length;
}

function validCapturedLap(point?: TrackCapturePoint) {
    if (!point || !Number.isFinite(point.lapNumber) || point.lapNumber <= 0) {
        return null;
    }
    return point.lapNumber;
}

function capturedRaceDurationMs(start: TrackCapturePoint, end: TrackCapturePoint) {
    if (Number.isFinite(start.currentRaceTime) && Number.isFinite(end.currentRaceTime) && end.currentRaceTime > start.currentRaceTime) {
        return (end.currentRaceTime - start.currentRaceTime) * 1000;
    }
    if (Number.isFinite(start.currentLap) && Number.isFinite(end.currentLap) && end.currentLap > start.currentLap) {
        return (end.currentLap - start.currentLap) * 1000;
    }
    return 0;
}

function simplifyClientPoints(points: BenchmarkPoint[], minDistance: number, maxPoints: number) {
    if (points.length <= 2) {
        return points;
    }
    const output = [points[0]];
    let last = points[0];
    for (const point of points.slice(1, -1)) {
        if (pointDistanceXZ(point, last) >= minDistance) {
            output.push(point);
            last = point;
        }
    }
    output.push(points[points.length - 1]);
    if (output.length <= maxPoints) {
        return output;
    }
    return Array.from({length: maxPoints}, (_, index) => {
        const sourceIndex = Math.round(index * (output.length - 1) / (maxPoints - 1));
        return output[sourceIndex];
    });
}

function pointsRouteLength(points: BenchmarkPoint[]) {
    return points.reduce((total, point, index) => index === 0 ? total : total + pointDistanceXZ(points[index - 1], point), 0);
}

function pointDistanceXZ(a: BenchmarkPoint, b: BenchmarkPoint) {
    return Math.hypot(a.x - b.x, a.z - b.z);
}

function routeDirection(points: BenchmarkPoint[]) {
    if (points.length < 2) {
        return {x: 0, z: 0};
    }
    const start = points[0];
    const candidate = points.find(point => pointDistanceXZ(point, start) > 5) || points[points.length - 1];
    const dx = candidate.x - start.x;
    const dz = candidate.z - start.z;
    const length = Math.hypot(dx, dz) || 1;
    return {x: dx / length, z: dz / length};
}

function directionAtClientEnd(points: BenchmarkPoint[]) {
    if (points.length < 2) {
        return {x: 1, z: 0};
    }
    const end = points[points.length - 1];
    const candidate = [...points].reverse().find(point => pointDistanceXZ(point, end) > 5) || points[0];
    const dx = end.x - candidate.x;
    const dz = end.z - candidate.z;
    const length = Math.hypot(dx, dz) || 1;
    return {x: dx / length, z: dz / length};
}

function inferClientTrackType(points: BenchmarkPoint[]): BenchmarkTrackType {
    if (points.length < 2) {
        return 'sprint';
    }
    return pointDistanceXZ(points[0], points[points.length - 1]) <= 40 && pointsRouteLength(points) >= 200 ? 'circuit' : 'sprint';
}

function clientGate(center: BenchmarkPoint, directionX: number, directionZ: number, widthMeters = 30, depthMeters = 20): BenchmarkGate {
    const length = Math.hypot(directionX, directionZ) || 1;
    return {
        center,
        directionX: directionX / length,
        directionZ: directionZ / length,
        widthMeters,
        depthMeters,
    };
}

function emptyGateWithSize(widthMeters: number, depthMeters: number): BenchmarkGate {
    return clientGate({x: 0, y: 0, z: 0}, 0, 0, widthMeters, depthMeters);
}

function applyGateSize(gate: BenchmarkGate, widthMeters: number, depthMeters: number): BenchmarkGate {
    return {...gate, widthMeters, depthMeters};
}

function gateFromCurrent(frame: TelemetryFrame, widthMeters: number, depthMeters: number): BenchmarkGate {
    const speedLength = Math.hypot(frame.velocityX, frame.velocityZ);
    const direction = speedLength > 0.5
        ? {x: frame.velocityX / speedLength, z: frame.velocityZ / speedLength}
        : {x: 1, z: 0};
    return clientGate(telemetryPoint(frame), direction.x, direction.z, widthMeters, depthMeters);
}

function clientCheckpoints(points: BenchmarkPoint[]) {
    if (points.length < 5) {
        return [];
    }
    return [points[Math.floor(points.length / 4)], points[Math.floor(points.length / 2)], points[Math.floor(points.length * 3 / 4)]];
}

function trackTypeLabel(type: BenchmarkTrackType | string | undefined, terms: Copy) {
    if (type === 'circuit') {
        return terms.circuitTrack;
    }
    if (type === 'sprint') {
        return terms.sprintTrack;
    }
    return terms.autoTrackType;
}

function trackMatchLevelLabel(value: string, terms: Copy) {
    if (value === 'strong') {
        return terms.strongMatch;
    }
    if (value === 'medium') {
        return terms.mediumMatch;
    }
    return value || '--';
}

function sourceSessionName(sessions: TelemetrySession[], sessionId: number) {
    const session = sessions.find(item => item.id === sessionId);
    return session?.sessionName || `#${sessionId}`;
}

function trackVehicleKeyLabel(vehicle: TrackVehicleKey) {
    return vehicle.label || [
        vehicle.carOrdinal ? `ID ${vehicle.carOrdinal}` : '',
        vehicle.carClass,
        vehicle.carPi ? `PI ${vehicle.carPi}` : '',
        vehicle.drivetrain,
    ].filter(Boolean).join(' / ') || '--';
}

function trackBaselineCount(profile?: TrackProfile | null) {
    return (profile?.vehicleReferences || []).reduce((total, reference) => total + (reference.baselineRunCount || 0), 0);
}

function formatReferenceBestDuration(reference: TrackVehicleReference, terms: Copy) {
    if (reference.bestTrackBaseline) {
        return formatDuration(reference.bestTrackBaseline.durationMs, terms);
    }
    if (reference.bestAutoBaseline) {
        return formatDuration(reference.bestAutoBaseline.run.durationMs, terms);
    }
    return terms.noAutoBaselines;
}

function formatReferenceBestConfidence(reference: TrackVehicleReference) {
    if (reference.bestTrackBaseline) {
        return formatPercentValue(reference.bestTrackBaseline.confidence);
    }
    if (reference.bestAutoBaseline) {
        return formatPercentValue(reference.bestAutoBaseline.run.confidence);
    }
    return '--';
}

function formatPointShort(point?: BenchmarkPoint) {
    if (!point) {
        return '--';
    }
    return `${formatNumber(point.x, 0)}, ${formatNumber(point.z, 0)}`;
}

function formatTrackSavedMessage(track: BenchmarkTrack, sourceSession: string, terms: Copy) {
    return terms.trackSavedDetails(track.id, trackTypeLabel(track.trackType as BenchmarkTrackType, terms), track.routeLengthMeters || 0, sourceSession, track.lapCountObserved || 0);
}

function extractionModeLabel(mode: BenchmarkExtractionMode | string | undefined, terms: Copy) {
    if (mode === 'first_lap') {
        return terms.firstLap;
    }
    if (mode === 'full_segment') {
        return terms.fullSegment;
    }
    return terms.autoBestLap;
}

function routeSvgPath(points: BenchmarkPoint[]) {
    return points.map((point, index) => {
        const svgPoint = routeSvgPoint(point, points);
        return `${index === 0 ? 'M' : 'L'} ${svgPoint.x.toFixed(1)} ${svgPoint.y.toFixed(1)}`;
    }).join(' ');
}

function routeSvgPoint(point: BenchmarkPoint, allPoints: BenchmarkPoint[]) {
    const xs = allPoints.map(item => item.x);
    const zs = allPoints.map(item => item.z);
    const minX = Math.min(...xs);
    const maxX = Math.max(...xs);
    const minZ = Math.min(...zs);
    const maxZ = Math.max(...zs);
    const width = Math.max(maxX - minX, 1);
    const height = Math.max(maxZ - minZ, 1);
    const padding = 14;
    return {
        x: padding + ((point.x - minX) / width) * (320 - padding * 2),
        y: padding + ((point.z - minZ) / height) * (180 - padding * 2),
    };
}

function hasTelemetryVehicleIdentity(frame: TelemetryFrame | null): frame is TelemetryFrame {
    return Boolean(frame && frame.carOrdinal > 0 && frame.carClass.trim());
}

function profileMatchesTelemetry(profile: TuneProfile, frame: TelemetryFrame) {
    return Number(profile.carOrdinal || 0) === frame.carOrdinal
        && normalizeVehicleClass(profile.carClass) === normalizeVehicleClass(frame.carClass);
}

function profileFormTelemetryState(profile: TuneProfileInput, frame: TelemetryFrame | null): 'match' | 'mismatch' | 'unknown' {
    if (!hasTelemetryVehicleIdentity(frame)) {
        return 'unknown';
    }
    if (!profile.carOrdinal || !profile.carClass) {
        return 'unknown';
    }
    return Number(profile.carOrdinal || 0) === frame.carOrdinal
        && normalizeVehicleClass(profile.carClass) === normalizeVehicleClass(frame.carClass)
        ? 'match'
        : 'mismatch';
}

function normalizeVehicleClass(value: string | undefined) {
    return String(value || '').trim().toUpperCase();
}

function formatTelemetryVehicle(frame: TelemetryFrame, terms: Copy) {
    const parts = [
        `${terms.vehicleId}: ${frame.carOrdinal}`,
        `${terms.classPi}: ${frame.carClass || '--'} / ${formatOptionalInt(frame.carPi)}`,
    ];
    if (frame.drivetrain || frame.numCylinders) {
        parts.push(`${terms.drivetrainCylinders}: ${frame.drivetrain || '--'} / ${formatOptionalInt(frame.numCylinders)}`);
    }
    return parts.join(' / ');
}

function formatTuneProfileVehicle(profile: TuneProfile, terms: Copy) {
    const parts = [
        profile.carName || terms.noProfile,
        profile.versionName,
    ];
    if (profile.carOrdinal) {
        parts.push(`${terms.vehicleId}: ${profile.carOrdinal}`);
    }
    if (profile.carClass || profile.pi) {
        parts.push(`${terms.classPi}: ${profile.carClass || '--'} / ${formatOptionalInt(profile.pi || undefined)}`);
    }
    if (profile.drivetrain || profile.numCylinders) {
        parts.push(`${terms.drivetrainCylinders}: ${profile.drivetrain || '--'} / ${formatOptionalInt(profile.numCylinders || undefined)}`);
    }
    if (profile.useCase) {
        parts.push(localizedUseCase(profile.useCase, terms));
    }
    return parts.filter(Boolean).join(' / ');
}

function normalizeUseCaseValue(value: unknown) {
    const raw = String(value ?? '').trim();
    if (!raw) {
        return '';
    }
    const exact = tuneUseCaseValues.find(item => item === raw);
    if (exact) {
        return exact;
    }
    return tuneUseCaseAliases[raw.toLowerCase()] || tuneUseCaseAliases[raw] || '';
}

function formatSessionVehicle(session: TelemetrySession, terms: Copy) {
    const parts = [];
    if (session.carOrdinal) {
        parts.push(`${terms.vehicleId}: ${session.carOrdinal}`);
    }
    if (session.carClass || session.carPi) {
        parts.push(`${terms.classPi}: ${session.carClass || '--'} / ${formatOptionalInt(session.carPi || undefined)}`);
    }
    if (session.drivetrain || session.numCylinders) {
        parts.push(`${terms.drivetrainCylinders}: ${session.drivetrain || '--'} / ${formatOptionalInt(session.numCylinders || undefined)}`);
    }
    return parts.join(' / ') || terms.noProfile;
}

function parseSessionTuneSnapshot(session: TelemetrySession | null): TuneProfile | null {
    if (!session?.tuneSnapshotJson) {
        return null;
    }
    try {
        const parsed = JSON.parse(session.tuneSnapshotJson) as TuneProfile;
        return parsed && (parsed.id || parsed.carName) ? parsed : null;
    } catch {
        return null;
    }
}

function formatEvidenceKey(value: string) {
    return value
        .replace(/_/g, ' ')
        .replace(/\b\w/g, character => character.toUpperCase());
}

function formatEvidenceValue(value: number) {
    if (!Number.isFinite(value)) {
        return '--';
    }
    if (Math.abs(value) >= 100) {
        return value.toFixed(0);
    }
    return value.toFixed(2);
}

function formatEvidenceValueForKey(key: string, value: number, terms: Copy) {
    if (key === 'corner_operation_state') {
        return localizedLabel(String(Math.round(value)), terms.cornerOperationStateLabels);
    }
    return formatEvidenceValue(value);
}

function formatGearSpeedRange(gear: GearPowerBand) {
    const min = Number.isFinite(gear.speedMinKmh) ? Math.max(0, gear.speedMinKmh) : 0;
    const max = Number.isFinite(gear.speedMaxKmh) ? Math.max(0, gear.speedMaxKmh) : 0;
    if (max <= 0) {
        return '--';
    }
    if (Math.abs(max - min) < 1) {
        return `${formatNumber(max, 0)} km/h`;
    }
    return `${formatNumber(min, 0)}-${formatNumber(max, 0)} km/h`;
}

function formatGearInPowerBandRange(gear: GearPowerBand) {
    const rpmMin = Number.isFinite(gear.inPowerBandRpmMin) ? Math.max(0, gear.inPowerBandRpmMin) : 0;
    const rpmMax = Number.isFinite(gear.inPowerBandRpmMax) ? Math.max(0, gear.inPowerBandRpmMax) : 0;
    if (rpmMax > 0) {
        if (Math.abs(rpmMax - rpmMin) < 1) {
            return `${formatNumber(rpmMax, 0)} rpm`;
        }
        return `${formatNumber(rpmMin, 0)}-${formatNumber(rpmMax, 0)} rpm`;
    }

    const ratioMin = Number.isFinite(gear.inPowerBandRatioMin) ? Math.max(0, gear.inPowerBandRatioMin) : 0;
    const ratioMax = Number.isFinite(gear.inPowerBandRatioMax) ? Math.max(0, gear.inPowerBandRatioMax) : 0;
    if (ratioMax > 0) {
        if (Math.abs(ratioMax - ratioMin) < 0.005) {
            return formatPercentValue(ratioMax);
        }
        return `${formatPercentValue(ratioMin)}-${formatPercentValue(ratioMax)}`;
    }

    return '--';
}

function formatGearMaxObserved(gear: GearPowerBand) {
    if (!Number.isFinite(gear.speedMaxKmh) || gear.speedMaxKmh <= 0) {
        return '--';
    }
    return `${formatNumber(gear.speedMaxKmh, 0)} km/h`;
}

function formatGearStrategyIssueCount(gearPower: GearPowerDiagnostic | undefined) {
    if (!gearPower || !gearPower.usableGearCount) {
        return '--';
    }
    const issueCount = gearPower.globalGearIssueCount || 0;
    const ratio = Number.isFinite(gearPower.globalGearIssueRatio)
        ? gearPower.globalGearIssueRatio
        : issueCount / gearPower.usableGearCount;
    return `${issueCount}/${gearPower.usableGearCount} (${formatPercentValue(ratio)})`;
}

function formatGearTelemetryDelta(before: number, after: number, delta: number, unit: 'km/h' | '%') {
    if (unit === '%') {
        return `${formatPercentValue(before || 0)} -> ${formatPercentValue(after || 0)} (${formatSignedPercentDelta(delta || 0)})`;
    }
    return `${formatNumber(before || 0, 0)} -> ${formatNumber(after || 0, 0)} ${unit} (${formatSignedNumber(delta || 0, 0)} ${unit})`;
}

function formatSignedNumber(value: number, digits = 1) {
    if (!Number.isFinite(value)) {
        return '--';
    }
    const prefix = value > 0 ? '+' : '';
    return `${prefix}${value.toFixed(digits)}`;
}

function formatSignedPercentDelta(value: number) {
    if (!Number.isFinite(value)) {
        return '--';
    }
    const percent = value * 100;
    const prefix = percent > 0 ? '+' : '';
    return `${prefix}${percent.toFixed(0)}pp`;
}

function gearPowerNoAdviceReasons(gearPower: GearPowerDiagnostic | undefined, terms: Copy) {
    if (!gearPower || gearPower.status === 'insufficient_data') {
        if (gearPower?.summary === 'no_unlocked_gear_samples') {
            return [terms.gearPowerNoUnlockedGears];
        }
        return [terms.gearPowerNeedSamples];
    }
    const evidence = gearPower.evidence || {};
    const totalSamples = evidence.power_band_total_samples || 0;
    const highLoadSamples = evidence.power_band_high_load_samples || 0;
    const usableGears = evidence.power_band_usable_gears || 0;
    const reasons: string[] = [];
    const addReason = (reason: string) => {
        if (!reasons.includes(reason)) {
            reasons.push(reason);
        }
    };
    if (gearPower.summary === 'no_unlocked_gear_samples') {
        addReason(terms.gearPowerNoUnlockedGears);
    }
    if (gearPower.summary === 'not_enough_high_load') {
        addReason(terms.gearPowerNeedHighLoad);
    }
    if (totalSamples <= 0) {
        addReason(terms.gearPowerNeedSamples);
    }
    if (usableGears <= 0) {
        addReason(terms.gearPowerNoUnlockedGears);
    }
    if (highLoadSamples < 4) {
        addReason(terms.gearPowerNeedHighLoad);
    }
    if (gearPower.powerBandSource === 'rpm_ratio_fallback') {
        addReason(terms.gearPowerFallbackLowConfidence);
    }
    if (gearPower.summary === 'traction_limited_power') {
        addReason(terms.gearPowerTractionFirst);
    }
    if (reasons.length === 0) {
        addReason(terms.gearPowerNoAdvice);
    }
    return reasons;
}

function decisionEvidenceEntries(decision: RoadTuningDecision) {
    return ['speed_band', 'speed_avg_kmh', 'speed_max_kmh', 'speed_min_kmh']
        .filter(key => typeof decision.evidence?.[key] === 'number')
        .map(key => [key, decision.evidence[key]] as [string, number]);
}

function formatDecisionEvidenceValue(key: string, value: number, language: Lang) {
    if (!Number.isFinite(value)) {
        return '--';
    }
    if (key === 'speed_band') {
        const band = Math.round(value);
        if (band === 1) {
            return language === 'zh' ? '低速' : 'Low speed';
        }
        if (band === 2) {
            return language === 'zh' ? '中速' : 'Mid speed';
        }
        if (band === 3) {
            return language === 'zh' ? '高速' : 'High speed';
        }
        return '--';
    }
    if (key.endsWith('_kmh')) {
        return `${value.toFixed(0)} km/h`;
    }
    return formatEvidenceValue(value);
}

function parseOptionalNumber(value: string) {
    if (value.trim() === '') {
        return null;
    }
    const parsed = Number(value);
    return Number.isFinite(parsed) ? parsed : null;
}

function parseStrictInteger(value: string) {
    const trimmed = value.trim();
    if (!/^\d+$/.test(trimmed)) {
        return null;
    }
    const parsed = Number(trimmed);
    return Number.isSafeInteger(parsed) ? parsed : null;
}

function parseRequiredInteger(value: string, label: string) {
    const parsed = parseStrictInteger(value);
    if (parsed === null) {
        throw new Error(`${label} must be an integer`);
    }
    return parsed;
}

function validateIntegerRange(
    form: RoadStaticTuneBaselineForm,
    errors: QuickTuneFieldErrors,
    messages: string[],
    key: RoadStaticTuneBaselineFormKey,
    label: string,
    min: number,
    max: number,
    terms: Copy,
) {
    const parsed = parseStrictInteger(form[key]);
    if (parsed === null || parsed < min || parsed > max) {
        const message = terms.quickTuneIntegerRange(label, min, max);
        errors[key] = message;
        messages.push(message);
    }
}

function tireDiameterCmFromSize(widthMm: number, aspectRatio: number, rimInches: number) {
    return (rimInches * 25.4 + 2 * widthMm * (aspectRatio / 100)) / 10;
}

function validateQuickTuneForm(form: RoadStaticTuneBaselineForm, terms: Copy) {
    const errors: QuickTuneFieldErrors = {};
    const messages: string[] = [];
    validateIntegerRange(form, errors, messages, 'weightKG', `${terms.tuneGeneratorWeight} (kg)`, 300, 3000, terms);
    validateIntegerRange(form, errors, messages, 'frontWeightPct', `${terms.tuneGeneratorFrontWeight} (%)`, 1, 99, terms);
    validateIntegerRange(form, errors, messages, 'pi', terms.tuneGeneratorPI, 100, 999, terms);
    validateIntegerRange(form, errors, messages, 'balanceBias', terms.quickTuneBalance, 50, 150, terms);
    validateIntegerRange(form, errors, messages, 'stiffnessBias', terms.quickTuneStiffness, 50, 150, terms);
    validateIntegerRange(form, errors, messages, 'speedBias', terms.quickTuneSpeed, 50, 150, terms);

    if (form.gearingEnabled === 'true') {
        validateIntegerRange(form, errors, messages, 'redlineRPM', `${terms.tuneGeneratorRedlineRPM} (rpm)`, 1000, 20000, terms);
        validateIntegerRange(form, errors, messages, 'gearCount', terms.tuneGeneratorGearCount, 2, 10, terms);
        const targetSpeedLabel = quickTuneTargetSpeedLabel(form.useCase, terms);
        const targetMin = form.useCase === 'Drift' || form.useCase === 'Drag' ? 40 : 1;
        const targetMax = form.useCase === 'Drift' ? 180 : form.useCase === 'Drag' ? 450 : 600;
        validateIntegerRange(form, errors, messages, 'targetTopSpeedKmh', `${targetSpeedLabel} (km/h)`, targetMin, targetMax, terms);

        const width = parseStrictInteger(form.tireWidthMm);
        const aspect = parseStrictInteger(form.tireAspectRatio);
        const rim = parseStrictInteger(form.tireRimInches);
        if (width === null || aspect === null || rim === null || tireDiameterCmFromSize(width, aspect, rim) < 40 || tireDiameterCmFromSize(width, aspect, rim) > 120) {
            const message = terms.quickTuneTireSizeInvalid;
            errors.tireWidthMm = errors.tireWidthMm || message;
            errors.tireAspectRatio = errors.tireAspectRatio || message;
            errors.tireRimInches = errors.tireRimInches || message;
            messages.push(message);
        }
    }

    return {errors, messages};
}

function roadBaselineInputFromForm(form: RoadStaticTuneBaselineForm): RoadStaticTuneBaselineInput {
    const gearingEnabled = form.gearingEnabled === 'true';
    const pi = parseRequiredInteger(form.pi, 'PI');
    if (pi < 100 || pi > 999) {
        throw new Error('PI must be between 100 and 999');
    }
    const redlineRPM = gearingEnabled ? parseRequiredInteger(form.redlineRPM, 'redlineRPM') : null;
    const gearCount = gearingEnabled ? parseRequiredInteger(form.gearCount, 'gearCount') : null;
    const tireWidthMm = gearingEnabled ? parseRequiredInteger(form.tireWidthMm, 'tireWidthMm') : null;
    const tireAspectRatio = gearingEnabled ? parseRequiredInteger(form.tireAspectRatio, 'tireAspectRatio') : null;
    const tireRimInches = gearingEnabled ? parseRequiredInteger(form.tireRimInches, 'tireRimInches') : null;
    const tireDiameterCm = gearingEnabled && tireWidthMm !== null && tireAspectRatio !== null && tireRimInches !== null
        ? tireDiameterCmFromSize(tireWidthMm, tireAspectRatio, tireRimInches)
        : null;
    const targetTopSpeedKmh = gearingEnabled ? parseRequiredInteger(form.targetTopSpeedKmh, 'targetTopSpeedKmh') : null;
    return {
        carName: form.carName.trim(),
        versionName: form.versionName.trim(),
        useCase: form.useCase || 'Road',
        carOrdinal: parseOptionalInteger(form.carOrdinal),
        carCategory: parseOptionalInteger(form.carCategory),
        pi,
        drivetrain: form.drivetrain,
        tireCompound: form.tireCompound || 'sport',
        weightKG: parseRequiredInteger(form.weightKG, 'weightKG'),
        frontWeightPct: parseRequiredInteger(form.frontWeightPct, 'frontWeightPct'),
        powerKW: parseOptionalNumber(form.powerKW),
        torqueNM: parseOptionalNumber(form.torqueNM),
        redlineRPM,
        gearCount,
        tireDiameterCm,
        targetTopSpeedKmh,
        frontRideHeightMinCm: parseOptionalNumber(form.frontRideHeightMinCm),
        frontRideHeightMaxCm: parseOptionalNumber(form.frontRideHeightMaxCm),
        rearRideHeightMinCm: parseOptionalNumber(form.rearRideHeightMinCm),
        rearRideHeightMaxCm: parseOptionalNumber(form.rearRideHeightMaxCm),
        frontAeroMinKgf: parseOptionalNumber(form.frontAeroMinKgf),
        frontAeroMaxKgf: parseOptionalNumber(form.frontAeroMaxKgf),
        rearAeroMinKgf: parseOptionalNumber(form.rearAeroMinKgf),
        rearAeroMaxKgf: parseOptionalNumber(form.rearAeroMaxKgf),
        frontRideHeightAdjustable: form.frontRideHeightAdjustable === 'true',
        rearRideHeightAdjustable: form.rearRideHeightAdjustable === 'true',
        frontAeroAdjustable: form.frontAeroAdjustable === 'true',
        rearAeroAdjustable: form.rearAeroAdjustable === 'true',
        balanceBias: parseRequiredInteger(form.balanceBias, 'balanceBias'),
        stiffnessBias: parseRequiredInteger(form.stiffnessBias, 'stiffnessBias'),
        speedBias: parseRequiredInteger(form.speedBias, 'speedBias'),
    };
}

function quickTuneAutoProfileName(draft: TuneProfileInput) {
    const carClass = draft.carClass || '';
    const pi = draft.pi ? String(draft.pi) : '';
    const drivetrain = draft.drivetrain || '';
    const prefix = draft.useCase === 'Drift'
        ? 'Quick Drift'
        : draft.useCase === 'Rally'
            ? 'Quick Rally'
            : draft.useCase === 'Offroad'
                ? 'Quick Offroad'
                : draft.useCase === 'Drag'
                    ? 'Quick Drag'
                    : 'Quick Tune';
    return [prefix, `${carClass}${pi}`.trim(), drivetrain].filter(Boolean).join(' ');
}

function loadQuickTuneStoredState(): QuickTuneStoredState {
    try {
        const raw = window.localStorage.getItem(quickTuneStorageKey);
        if (!raw) {
            return {};
        }
        const parsed = JSON.parse(raw) as QuickTuneStoredState;
        return {
            form: parsed.form ? {...emptyRoadBaselineForm, ...parsed.form} : undefined,
            result: parsed.result || null,
            selectedFields: Array.isArray(parsed.selectedFields) ? parsed.selectedFields : undefined,
            targetProfileId: typeof parsed.targetProfileId === 'number' ? parsed.targetProfileId : undefined,
        };
    } catch {
        return {};
    }
}

function saveQuickTuneStoredState(state: QuickTuneStoredState) {
    try {
        window.localStorage.setItem(quickTuneStorageKey, JSON.stringify(state));
    } catch {
        // Ignore storage quota/private mode failures; quick tune still works in memory.
    }
}

function parseOptionalInteger(value: string) {
    const parsed = parseOptionalNumber(value);
    return parsed === null ? null : Math.trunc(parsed);
}

function formatBaselineFormNumber(value: string, step: number) {
    if (value.trim() === '') {
        return '';
    }
    const parsed = Number(value);
    if (!Number.isFinite(parsed)) {
        return value;
    }
    return parsed.toFixed(decimalPlacesForStep(step));
}

function profileFieldStep(field: ProfileField) {
    return field.step || 'any';
}

function decimalPlacesForStep(step: number) {
    if (!Number.isFinite(step) || step <= 0) {
        return 0;
    }
    const text = String(step);
    const exponent = text.match(/e-(\d+)$/i);
    if (exponent) {
        return Number(exponent[1]) || 0;
    }
    const dot = text.indexOf('.');
    return dot >= 0 ? text.length - dot - 1 : 0;
}

function profileInputValue(field: ProfileField, value: unknown) {
    if (value === undefined || value === null || value === '') {
        return '';
    }
    if (field.kind !== 'number' || typeof value !== 'number' || !Number.isFinite(value)) {
        return String(value);
    }
    const step = Number(field.step || '');
    if (!Number.isFinite(step) || step <= 0) {
        return String(value);
    }
    return value.toFixed(decimalPlacesForStep(step));
}

function profileFieldLabel(field: ProfileField, language: Lang) {
    const unit = field.unit?.[language];
    return unit ? `${field.label[language]} (${unit})` : field.label[language];
}

function isLockedProfileField(field: ProfileField, value: unknown) {
    return field.kind === 'number' && field.group !== 'vehicle' && field.group !== 'power' && field.group !== 'tire' && (value === undefined || value === null || value === '');
}

function profileLockedLabel(language: Lang) {
    return language === 'zh' ? '锁定' : 'Locked';
}

function profileLockedHint(language: Lang) {
    return language === 'zh' ? '未填写，按锁定处理' : 'Blank, treated as locked';
}

function optionalPositiveInt(value: number | undefined) {
    if (value === undefined || !Number.isFinite(value) || value <= 0) {
        return null;
    }
    return Math.round(value);
}

function optionalNonNegativeInt(value: number | undefined) {
    if (value === undefined || !Number.isFinite(value) || value < 0) {
        return null;
    }
    return Math.round(value);
}

function cleanProfileInput(input: TuneProfileInput): TuneProfileInput {
    const clean: TuneProfileInput = {...input, carName: input.carName.trim()};
    const mutable = clean as Record<string, unknown>;
    for (const field of profileFields) {
        const value = clean[field.key];
        if (field.kind === 'text' || field.kind === 'textarea' || field.kind === 'select') {
            mutable[field.key] = field.key === 'useCase' ? normalizeUseCaseValue(value) : String(value ?? '').trim();
        }
        if (field.kind === 'number' && (value === undefined || value === null || (typeof value === 'number' && Number.isNaN(value)))) {
            mutable[field.key] = null;
        }
    }
    return withCalculatedPowerToWeight(clean);
}

function withCalculatedPowerToWeight(input: TuneProfileInput): TuneProfileInput {
    const next = {...input};
    const power = typeof next.powerKW === 'number' && Number.isFinite(next.powerKW) && next.powerKW > 0 ? next.powerKW : null;
    const weight = typeof next.weightKG === 'number' && Number.isFinite(next.weightKG) && next.weightKG > 0 ? next.weightKG : null;
    next.powerToWeightKWPerKG = power && weight ? Number((power / weight).toFixed(4)) : null;
    return next;
}

function profileToInput(profile: TuneProfile): TuneProfileInput {
    const input = cleanProfileInput(profile);
    return {...input};
}

function profileInputToProfile(input: TuneProfileInput, base?: TuneProfile | null): TuneProfile {
    return {
        ...(base || {id: 0, createdAt: '', updatedAt: ''}),
        ...cleanProfileInput(input),
    };
}

function formatSnapshotChangedFields(fields: string[], language: Lang) {
    if (!fields.length) {
        return COPY[language].noChanges;
    }
    const labels = fields.slice(0, 4).map(field => {
        const profileField = profileFields.find(item => String(item.key) === field);
        return profileField ? profileFieldLabel(profileField, language) : formatEvidenceKey(field);
    });
    const remaining = fields.length - labels.length;
    return remaining > 0 ? `${labels.join(' / ')} +${remaining}` : labels.join(' / ');
}

function formatConcreteSuggestion(action: SuggestedAction, event: DetectedEvent, profile: TuneProfile | null, language: Lang, terms: Copy) {
    const adjustments = concreteSuggestionAdjustments(action, event, profile, language);
    if (adjustments.length === 0) {
        return `${localizedLabel(action.direction, terms.actionDirections)} ${localizedLabel(action.amount, terms.actionAmounts)}`;
    }
    return adjustments.map(adjustment => {
        const verb = concreteDirectionLabel(adjustment.delta, action, language);
        return `${adjustment.label}: ${formatTuneNumber(adjustment.current, adjustment.step)} -> ${formatTuneNumber(adjustment.target, adjustment.step)}, ${verb} ${formatTuneNumber(Math.abs(adjustment.delta), adjustment.step)}`;
    }).join(language === 'zh' ? '；' : '; ');
}

function previewDecisionAction(action: RoadTuningDecisionAction, profile: TuneProfile, language: Lang) {
    const event = {
        id: action.id,
        type: action.family,
        severity: 'medium',
        startMs: 0,
        endMs: 0,
        durationMs: 0,
        segment: '',
        evidence: action.evidence || {},
        suggestedActions: [],
    } as DetectedEvent;
    const suggested = {
        priority: 0,
        category: action.category,
        item: action.item,
        direction: action.direction,
        amount: action.amount,
        reason: action.reason,
    };
    const adjustments = concreteSuggestionAdjustments(suggested, event, profile, language);
    if (adjustments.length === 0) {
        return null;
    }
    return adjustments.map(adjustment => {
        const verb = concreteDirectionLabel(adjustment.delta, suggested, language);
        return `${adjustment.label}: ${formatTuneNumber(adjustment.current, adjustment.step)} -> ${formatTuneNumber(adjustment.target, adjustment.step)}, ${verb} ${formatTuneNumber(Math.abs(adjustment.delta), adjustment.step)}`;
    }).join(language === 'zh' ? '，' : '; ');
}

function formatTuneExplanationNote(action: SuggestedAction, event: DetectedEvent, explanations: TuneAdjustmentExplanation[]) {
    const keys = tuneExplanationKeysForAction(action.item, Math.round(event.evidence?.gear || 0));
    const seen = new Set<string>();
    const notes: string[] = [];
    keys.forEach(key => {
        const explanation = explanations.find(item => tuneExplanationKey(item) === key);
        if (explanation?.description && !seen.has(explanation.description)) {
            seen.add(explanation.description);
            notes.push(explanation.description);
        }
    });
    return notes.join('；');
}

function tuneExplanationKeysForAction(item: string, gear: number) {
    switch (item) {
        case 'gear_1':
            return [makeTuneExplanationKey('齿轮', '前进档', '1档')];
        case 'current_gear':
            return gear >= 1 && gear <= 10 ? [makeTuneExplanationKey('齿轮', '前进档', `${gear}档`)] : [];
        case 'final_drive':
            return [makeTuneExplanationKey('齿轮', '前进档', '最终传动')];
        case 'brake_balance':
            return [makeTuneExplanationKey('刹车', '制动力', '平衡')];
        case 'brake_pressure':
            return [makeTuneExplanationKey('刹车', '制动力', '压力')];
        case 'rear_diff_accel':
        case 'drive_diff_accel':
            return [makeTuneExplanationKey('差速器', '后侧', '加速')];
        case 'rear_diff_decel':
            return [makeTuneExplanationKey('差速器', '后侧', '减速')];
        case 'drive_tire_pressure':
        case 'tire_pressure':
            return [makeTuneExplanationKey('轮胎', '胎压', '前侧'), makeTuneExplanationKey('轮胎', '胎压', '后侧')];
        case 'front_arb':
            return [makeTuneExplanationKey('防倾杆', '防倾杆', '前侧')];
        case 'rear_arb':
            return [makeTuneExplanationKey('防倾杆', '防倾杆', '后侧')];
        case 'front_rebound':
            return [makeTuneExplanationKey('阻尼', '回弹硬度', '前侧')];
        case 'rear_rebound':
            return [makeTuneExplanationKey('阻尼', '回弹硬度', '后侧')];
        case 'front_camber':
            return [makeTuneExplanationKey('轮胎定位', '外倾角', '前侧')];
        case 'front_and_rear_aero':
            return [makeTuneExplanationKey('空气动力学设置', '下压力', '前侧'), makeTuneExplanationKey('空气动力学设置', '下压力', '后侧')];
        case 'ride_height':
            return [makeTuneExplanationKey('弹簧', '车身高度', '前侧'), makeTuneExplanationKey('弹簧', '车身高度', '后侧')];
        case 'spring_rate':
            return [makeTuneExplanationKey('弹簧', '弹簧', '前侧'), makeTuneExplanationKey('弹簧', '弹簧', '后侧')];
        case 'bump':
            return [makeTuneExplanationKey('阻尼', '压缩硬度', '前侧'), makeTuneExplanationKey('阻尼', '压缩硬度', '后侧')];
        default:
            return [];
    }
}

function tuneExplanationKey(item: TuneAdjustmentExplanation) {
    return makeTuneExplanationKey(item.category, item.item, item.detail);
}

function makeTuneExplanationKey(category: string, item: string, detail: string) {
    return `${category}/${item}/${detail}`;
}

type ConcreteSuggestionAdjustment = {
    label: string;
    current: number;
    target: number;
    delta: number;
    step: number;
};

function concreteSuggestionAdjustments(action: SuggestedAction, event: DetectedEvent, profile: TuneProfile | null, language: Lang): ConcreteSuggestionAdjustment[] {
    if (!profile) {
        return [];
    }
    const fields = actionTuneFields(action, event, profile);
    return fields.flatMap(field => {
        const raw = profile[field.key];
        if (typeof raw !== 'number' || !Number.isFinite(raw)) {
            return [];
        }
        const step = Number(field.step || '1') || 1;
        const delta = actionDelta(action, raw, step, field.unit?.en === '%' ? '%' : '');
        if (!delta) {
            return [];
        }
        return [{
            label: profileFieldLabel(field, language),
            current: raw,
            target: roundToStep(raw + delta, step),
            delta: roundToStep(delta, step),
            step,
        }];
    });
}

function actionTuneFields(action: SuggestedAction, event: DetectedEvent, profile: TuneProfile): ProfileField[] {
    const field = (key: keyof TuneProfileInput) => profileFields.find(item => item.key === key);
    const pick = (...keys: Array<keyof TuneProfileInput>) => keys.map(field).filter(Boolean) as ProfileField[];
    const driven = (front: keyof TuneProfileInput, rear: keyof TuneProfileInput) => {
        const drivetrain = String(profile.drivetrain || '').trim().toUpperCase();
        if (drivetrain === 'FWD') {
            return pick(front);
        }
        if (drivetrain === 'RWD') {
            return pick(rear);
        }
        return pick(front, rear);
    };
    switch (action.item) {
        case 'front_tire_pressure':
            return pick('frontTirePressure');
        case 'rear_tire_pressure':
            return pick('rearTirePressure');
        case 'gear_1':
            return pick('gear1');
        case 'gear_2':
            return pick('gear2');
        case 'gear_3':
            return pick('gear3');
        case 'gear_4':
            return pick('gear4');
        case 'gear_5':
            return pick('gear5');
        case 'gear_6':
            return pick('gear6');
        case 'gear_7':
            return pick('gear7');
        case 'gear_8':
            return pick('gear8');
        case 'gear_9':
            return pick('gear9');
        case 'gear_10':
            return pick('gear10');
        case 'current_gear': {
            const gear = Math.round(event.evidence?.gear || 0);
            const key = (`gear${gear}`) as keyof TuneProfileInput;
            return gear >= 1 && gear <= 10 ? pick(key) : [];
        }
        case 'final_drive':
            return pick('finalDrive');
        case 'brake_balance':
            return pick('brakeBalance');
        case 'brake_pressure':
            return pick('brakePressure');
        case 'rear_diff_accel':
            return pick('rearDiffAccel');
        case 'rear_diff_decel':
            return pick('rearDiffDecel');
        case 'front_diff_accel':
            return pick('frontDiffAccel');
        case 'front_diff_decel':
            return pick('frontDiffDecel');
        case 'drive_diff_accel':
            return driven('frontDiffAccel', 'rearDiffAccel');
        case 'drive_tire_pressure':
            return driven('frontTirePressure', 'rearTirePressure');
        case 'tire_pressure':
            return pick('frontTirePressure', 'rearTirePressure');
        case 'front_arb':
            return pick('frontArb');
        case 'rear_arb':
            return pick('rearArb');
        case 'front_rebound':
            return pick('frontRebound');
        case 'rear_rebound':
            return pick('rearRebound');
        case 'front_camber':
            return pick('frontCamber');
        case 'front_and_rear_aero':
            return pick('frontAero', 'rearAero');
        case 'ride_height':
            return pick('frontRideHeight', 'rearRideHeight');
        case 'spring_rate':
            return pick('frontSpring', 'rearSpring');
        case 'bump':
            return pick('frontBump', 'rearBump');
        default:
            return [];
    }
}

function actionDelta(action: SuggestedAction, current: number, step: number, unit: string) {
    const amount = String(action.amount || '').trim().toLowerCase();
    const direction = String(action.direction || '').trim().toLowerCase();
    let magnitude = 0;
    if (amount === 'slightly more negative') {
        return -step;
    }
    if (amount === 'one small step' || amount === 'avoid bottoming' || amount === '0.5 psi') {
        magnitude = step;
    } else if (amount.includes('%')) {
        const base = firstAmountNumber(amount);
        if (!base) {
            return 0;
        }
        magnitude = unit === '%' ? base : Math.max(step, Math.abs(current) * base / 100);
    } else {
        magnitude = firstAmountNumber(amount);
    }
    magnitude = Math.max(step, roundToStep(magnitude, step));
    if (direction === 'decrease') {
        return -magnitude;
    }
    if (direction === 'increase') {
        return magnitude;
    }
    if (direction === 'check' && (amount.includes('bottom') || amount.includes('small step'))) {
        return magnitude;
    }
    return 0;
}

function firstAmountNumber(value: string) {
    const match = value.match(/\d+(?:\.\d+)?/);
    return match ? Number(match[0]) : 0;
}

function roundToStep(value: number, step: number) {
    if (!step) {
        return value;
    }
    const decimals = step >= 1 ? 0 : step >= 0.1 ? 1 : 2;
    return Number((Math.round(value / step) * step).toFixed(decimals + 2));
}

function formatTuneNumber(value: number, step: number) {
    const decimals = step >= 1 ? 0 : step >= 0.1 ? 1 : 2;
    return value.toFixed(decimals);
}

function formatDraftValue(value: number | null | undefined, step: number, unit: string) {
    if (value === undefined || value === null || !Number.isFinite(value)) {
        return '--';
    }
    const body = formatTuneNumber(value, step || 0.1);
    if (!unit) {
        return body;
    }
    return unit === '%' ? `${body}%` : `${body} ${unit}`;
}

function formatDraftDelta(value: number | null | undefined, step: number, unit: string) {
    if (value === undefined || value === null || !Number.isFinite(value)) {
        return '';
    }
    const sign = value > 0 ? '+' : '';
    return `${sign}${formatDraftValue(value, step, unit)}`;
}

function formatRetestMetricValue(key: string, value: number, terms: Copy) {
    if (!Number.isFinite(value)) {
        return '--';
    }
    if (key.endsWith('_duration_ms')) {
        return formatDuration(value, terms);
    }
    if (key.includes('temp')) {
        return `${value.toFixed(1)} °C`;
    }
    if (key === 'avg_speed_kmh' || key === 'max_speed_kmh') {
        return `${value.toFixed(1)} km/h`;
    }
    if (key === 'event_count' || key === 'gear_problem_count') {
        return String(Math.round(value));
    }
    return Math.abs(value) >= 100 ? value.toFixed(0) : value.toFixed(1);
}

function formatRetestMetricDelta(key: string, value: number, terms: Copy) {
    if (!Number.isFinite(value)) {
        return '--';
    }
    const sign = value > 0 ? '+' : '';
    if (key.endsWith('_duration_ms')) {
        return `${sign}${formatDuration(Math.abs(value), terms)}`;
    }
    if (key.includes('temp')) {
        return `${sign}${value.toFixed(1)} °C`;
    }
    if (key === 'avg_speed_kmh' || key === 'max_speed_kmh') {
        return `${sign}${value.toFixed(1)} km/h`;
    }
    if (key === 'event_count' || key === 'gear_problem_count') {
        return `${sign}${Math.round(value)}`;
    }
    return `${sign}${Math.abs(value) >= 100 ? value.toFixed(0) : value.toFixed(1)}`;
}

function formatRetestMetricSummary(value: string, terms: Copy) {
    const [key, status] = value.split(':');
    const label = localizedLabel(key || '', terms.retestMetricLabels) || key;
    const statusLabel = localizedLabel(status || '', terms.retestStatusLabels) || status;
    return `${label}: ${statusLabel}`;
}

function concreteDirectionLabel(delta: number, action: SuggestedAction, language: Lang) {
    if (action.amount === 'slightly more negative') {
        return language === 'zh' ? '增加负外倾' : 'more negative by';
    }
    if (delta < 0) {
        return language === 'zh' ? '降低' : 'decrease by';
    }
    return language === 'zh' ? '增加' : 'increase by';
}

function ruleProfileToInput(profile: RuleThresholdProfile): RuleThresholdProfileInput {
    return {
        name: profile.name,
        carClass: profile.carClass || '',
        drivetrain: profile.drivetrain || '',
        useCase: profile.useCase || '',
        gameMode: profile.gameMode || '',
        configJson: profile.configJson || '',
    };
}

function compareProfiles(left: TuneProfile | null, right: TuneProfile | null, language: Lang) {
    if (!left || !right) {
        return [];
    }
    return profileFields
        .map(field => {
            const leftValue = field.key === 'useCase' ? localizedUseCase(String(left[field.key] || ''), COPY[language]) || '--' : formatProfileValue(left[field.key], field, language);
            const rightValue = field.key === 'useCase' ? localizedUseCase(String(right[field.key] || ''), COPY[language]) || '--' : formatProfileValue(right[field.key], field, language);
            return {key: String(field.key), label: profileFieldLabel(field, language), left: leftValue, right: rightValue};
        })
        .filter(item => item.left !== item.right);
}

function formatProfileValue(value: unknown, field?: ProfileField, language: Lang = 'en') {
    if (value === undefined || value === null || value === '') {
        if (field && isLockedProfileField(field, value)) {
            return profileLockedLabel(language);
        }
        return '--';
    }
    if (typeof value === 'number') {
        return field ? profileInputValue(field, value) : (Number.isInteger(value) ? String(value) : value.toFixed(2));
    }
    return String(value);
}

function formatBaselineGeneratedValue(value: unknown, fieldKey: string, field?: ProfileField, language: Lang = 'en') {
    if (value === undefined || value === null || value === '') {
        return '--';
    }
    const base = formatProfileValue(value, field, language);
    if ((fieldKey === 'frontTirePressure' || fieldKey === 'rearTirePressure') && typeof value === 'number' && Number.isFinite(value)) {
        return `${base} (≈ ${formatNumber(value * 14.5038, 1)} PSI)`;
    }
    return base;
}

function formatComparisonValue(value: number, unit: string) {
    const body = Math.abs(value) >= 100 ? value.toFixed(0) : value.toFixed(2);
    return unit ? `${body} ${unit}` : body;
}

function formatSigned(value: number, unit: string) {
    const sign = value > 0 ? '+' : '';
    return `${sign}${formatComparisonValue(value, unit)}`;
}

function clamp(value: number, min: number, max: number) {
    return Math.max(min, Math.min(max, value));
}

function preferredListenAddress(items: NetworkInterface[]) {
    const candidates = items.filter(item => item.isUp && item.isPrivate && !item.isLoopback && item.address !== '0.0.0.0');
    const wlan = candidates.find(item => {
        const name = `${item.name} ${item.displayName}`.toLowerCase();
        return name.includes('wlan') || name.includes('wi-fi') || name.includes('wifi') || name.includes('wireless') || name.includes('无线') || name.includes('無線');
    });
    return wlan?.address || candidates[0]?.address || '0.0.0.0';
}

function interfaceLabel(item: NetworkInterface, language: Lang) {
    const t = COPY[language];
    const suffix = item.isLoopback ? t.loopback : item.isPrivate ? t.lan : t.ip;
    if (item.address === '0.0.0.0') {
        return `${t.allInterfaces} (0.0.0.0)`;
    }
    return `${item.displayName} - ${item.address} (${suffix})`;
}

export default App
