export namespace main {
	
	export class RecommendedCarsFileResult {
	    path: string;
	    count: number;
	
	    static createFrom(source: any = {}) {
	        return new RecommendedCarsFileResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.count = source["count"];
	    }
	}
	export class RecommendedCarsFileSelection {
	    path: string;
	    exists: boolean;
	    version: string;
	    ids: string[];
	    tuneCodes: string[];
	    count: number;
	
	    static createFrom(source: any = {}) {
	        return new RecommendedCarsFileSelection(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.exists = source["exists"];
	        this.version = source["version"];
	        this.ids = source["ids"];
	        this.tuneCodes = source["tuneCodes"];
	        this.count = source["count"];
	    }
	}
	export class TuneWebServerStatus {
	    running: boolean;
	    port: number;
	    url: string;
	    lanAddress: string;
	    lastError: string;
	
	    static createFrom(source: any = {}) {
	        return new TuneWebServerStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.running = source["running"];
	        this.port = source["port"];
	        this.url = source["url"];
	        this.lanAddress = source["lanAddress"];
	        this.lastError = source["lastError"];
	    }
	}

}

export namespace storage {
	
	export class BaselineGeneratedField {
	    fieldKey: string;
	    group: string;
	    value?: number;
	    unit: string;
	    reason: string;
	    defaultSelected: boolean;
	
	    static createFrom(source: any = {}) {
	        return new BaselineGeneratedField(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.fieldKey = source["fieldKey"];
	        this.group = source["group"];
	        this.value = source["value"];
	        this.unit = source["unit"];
	        this.reason = source["reason"];
	        this.defaultSelected = source["defaultSelected"];
	    }
	}
	export class BaselineSkippedField {
	    fieldKey: string;
	    group: string;
	    reason: string;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new BaselineSkippedField(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.fieldKey = source["fieldKey"];
	        this.group = source["group"];
	        this.reason = source["reason"];
	        this.message = source["message"];
	    }
	}
	export class BaselineTierRecommendation {
	    fieldKey: string;
	    group: string;
	    tier: string;
	    reason: string;
	    applicable: boolean;
	
	    static createFrom(source: any = {}) {
	        return new BaselineTierRecommendation(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.fieldKey = source["fieldKey"];
	        this.group = source["group"];
	        this.tier = source["tier"];
	        this.reason = source["reason"];
	        this.applicable = source["applicable"];
	    }
	}
	export class BenchmarkPoint {
	    x: number;
	    y: number;
	    z: number;
	
	    static createFrom(source: any = {}) {
	        return new BenchmarkPoint(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.x = source["x"];
	        this.y = source["y"];
	        this.z = source["z"];
	    }
	}
	export class BenchmarkGate {
	    center: BenchmarkPoint;
	    directionX: number;
	    directionZ: number;
	    widthMeters: number;
	    depthMeters: number;
	
	    static createFrom(source: any = {}) {
	        return new BenchmarkGate(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.center = this.convertValues(source["center"], BenchmarkPoint);
	        this.directionX = source["directionX"];
	        this.directionZ = source["directionZ"];
	        this.widthMeters = source["widthMeters"];
	        this.depthMeters = source["depthMeters"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class BenchmarkRun {
	    id: number;
	    sessionId: number;
	    trackId: number;
	    trackName: string;
	    startMs: number;
	    endMs: number;
	    durationMs: number;
	    confidence: number;
	    avgSpeedKmh?: number;
	    maxSpeedKmh?: number;
	    routeProgress01?: number;
	    geometryLengthMeters?: number;
	    trackLengthErrorPct?: number;
	    distanceTraveledDeltaMeters?: number;
	    currentRaceTimeDeltaSeconds?: number;
	    avgLateralErrorMeters?: number;
	    maxLateralErrorMeters?: number;
	    warningFlags: string;
	    eventCount: number;
	    driverMode: string;
	    driverModeConfidence: number;
	    driverModeEvidenceJson: string;
	    valid: boolean;
	    createdAt: string;
	
	    static createFrom(source: any = {}) {
	        return new BenchmarkRun(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.sessionId = source["sessionId"];
	        this.trackId = source["trackId"];
	        this.trackName = source["trackName"];
	        this.startMs = source["startMs"];
	        this.endMs = source["endMs"];
	        this.durationMs = source["durationMs"];
	        this.confidence = source["confidence"];
	        this.avgSpeedKmh = source["avgSpeedKmh"];
	        this.maxSpeedKmh = source["maxSpeedKmh"];
	        this.routeProgress01 = source["routeProgress01"];
	        this.geometryLengthMeters = source["geometryLengthMeters"];
	        this.trackLengthErrorPct = source["trackLengthErrorPct"];
	        this.distanceTraveledDeltaMeters = source["distanceTraveledDeltaMeters"];
	        this.currentRaceTimeDeltaSeconds = source["currentRaceTimeDeltaSeconds"];
	        this.avgLateralErrorMeters = source["avgLateralErrorMeters"];
	        this.maxLateralErrorMeters = source["maxLateralErrorMeters"];
	        this.warningFlags = source["warningFlags"];
	        this.eventCount = source["eventCount"];
	        this.driverMode = source["driverMode"];
	        this.driverModeConfidence = source["driverModeConfidence"];
	        this.driverModeEvidenceJson = source["driverModeEvidenceJson"];
	        this.valid = source["valid"];
	        this.createdAt = source["createdAt"];
	    }
	}
	export class BenchmarkTrack {
	    id: number;
	    name: string;
	    sourceMode: string;
	    trackType: string;
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
	    sourceSessionId?: number;
	    lapCountObserved: number;
	    notes: string;
	    createdAt: string;
	    updatedAt: string;
	
	    static createFrom(source: any = {}) {
	        return new BenchmarkTrack(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.sourceMode = source["sourceMode"];
	        this.trackType = source["trackType"];
	        this.start = this.convertValues(source["start"], BenchmarkPoint);
	        this.end = this.convertValues(source["end"], BenchmarkPoint);
	        this.startRadius = source["startRadius"];
	        this.endRadius = source["endRadius"];
	        this.directionX = source["directionX"];
	        this.directionZ = source["directionZ"];
	        this.startGate = this.convertValues(source["startGate"], BenchmarkGate);
	        this.finishGate = this.convertValues(source["finishGate"], BenchmarkGate);
	        this.checkpoints = this.convertValues(source["checkpoints"], BenchmarkPoint);
	        this.routeLengthMeters = source["routeLengthMeters"];
	        this.hasDrivingLine = source["hasDrivingLine"];
	        this.polyline = this.convertValues(source["polyline"], BenchmarkPoint);
	        this.sourceSessionId = source["sourceSessionId"];
	        this.lapCountObserved = source["lapCountObserved"];
	        this.notes = source["notes"];
	        this.createdAt = source["createdAt"];
	        this.updatedAt = source["updatedAt"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class BenchmarkTrackExtractionInput {
	    sessionId: number;
	    name: string;
	    trackType: string;
	    extractionMode: string;
	    startGate?: BenchmarkGate;
	    finishGate?: BenchmarkGate;
	
	    static createFrom(source: any = {}) {
	        return new BenchmarkTrackExtractionInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.sessionId = source["sessionId"];
	        this.name = source["name"];
	        this.trackType = source["trackType"];
	        this.extractionMode = source["extractionMode"];
	        this.startGate = this.convertValues(source["startGate"], BenchmarkGate);
	        this.finishGate = this.convertValues(source["finishGate"], BenchmarkGate);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class BenchmarkTrackInput {
	    name: string;
	    sourceMode: string;
	    trackType: string;
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
	    sourceSessionId?: number;
	    lapCountObserved: number;
	    notes: string;
	
	    static createFrom(source: any = {}) {
	        return new BenchmarkTrackInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.sourceMode = source["sourceMode"];
	        this.trackType = source["trackType"];
	        this.start = this.convertValues(source["start"], BenchmarkPoint);
	        this.end = this.convertValues(source["end"], BenchmarkPoint);
	        this.startRadius = source["startRadius"];
	        this.endRadius = source["endRadius"];
	        this.directionX = source["directionX"];
	        this.directionZ = source["directionZ"];
	        this.startGate = this.convertValues(source["startGate"], BenchmarkGate);
	        this.finishGate = this.convertValues(source["finishGate"], BenchmarkGate);
	        this.checkpoints = this.convertValues(source["checkpoints"], BenchmarkPoint);
	        this.routeLengthMeters = source["routeLengthMeters"];
	        this.hasDrivingLine = source["hasDrivingLine"];
	        this.polyline = this.convertValues(source["polyline"], BenchmarkPoint);
	        this.sourceSessionId = source["sourceSessionId"];
	        this.lapCountObserved = source["lapCountObserved"];
	        this.notes = source["notes"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class BrakeToTireDiagnostic {
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
	
	    static createFrom(source: any = {}) {
	        return new BrakeToTireDiagnostic(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.status = source["status"];
	        this.summary = source["summary"];
	        this.explanation = source["explanation"];
	        this.confidence = source["confidence"];
	        this.sampleCount = source["sampleCount"];
	        this.brakeSampleCount = source["brakeSampleCount"];
	        this.averageBrake = source["averageBrake"];
	        this.peakBrake = source["peakBrake"];
	        this.averageHandBrake = source["averageHandBrake"];
	        this.peakHandBrake = source["peakHandBrake"];
	        this.averageSpeedKmh = source["averageSpeedKmh"];
	        this.speedDeltaKmh = source["speedDeltaKmh"];
	        this.averageSteer = source["averageSteer"];
	        this.averageDecelMps2 = source["averageDecelMps2"];
	        this.averageDecelG = source["averageDecelG"];
	        this.peakDecelG = source["peakDecelG"];
	        this.averagePlaneG = source["averagePlaneG"];
	        this.peakPlaneG = source["peakPlaneG"];
	        this.frontSlipRatioP90 = source["frontSlipRatioP90"];
	        this.rearSlipRatioP90 = source["rearSlipRatioP90"];
	        this.frontCombinedSlipP90 = source["frontCombinedSlipP90"];
	        this.rearCombinedSlipP90 = source["rearCombinedSlipP90"];
	        this.frontRearSlipDelta = source["frontRearSlipDelta"];
	        this.trailBraking = source["trailBraking"];
	        this.handbrakeActive = source["handbrakeActive"];
	        this.evidence = source["evidence"];
	    }
	}
	export class TireModelHint {
	    code: string;
	    severity: string;
	    direction: string;
	    reason: string;
	
	    static createFrom(source: any = {}) {
	        return new TireModelHint(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.code = source["code"];
	        this.severity = source["severity"];
	        this.direction = source["direction"];
	        this.reason = source["reason"];
	    }
	}
	export class CamberInference {
	    status: string;
	    confidence: string;
	    frontState: string;
	    rearState: string;
	    summary: string;
	    explanation: string;
	    warnings: string[];
	    hints: TireModelHint[];
	    evidence: Record<string, number>;
	
	    static createFrom(source: any = {}) {
	        return new CamberInference(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.status = source["status"];
	        this.confidence = source["confidence"];
	        this.frontState = source["frontState"];
	        this.rearState = source["rearState"];
	        this.summary = source["summary"];
	        this.explanation = source["explanation"];
	        this.warnings = source["warnings"];
	        this.hints = this.convertValues(source["hints"], TireModelHint);
	        this.evidence = source["evidence"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class FH6Car {
	    carId: string;
	    year: number;
	    make: string;
	    model: string;
	    alias: string[];
	    basePi: number;
	    drivetrainDefault: string;
	    source: string;
	    sourceRef: string;
	    createdAt: string;
	    updatedAt: string;
	
	    static createFrom(source: any = {}) {
	        return new FH6Car(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.carId = source["carId"];
	        this.year = source["year"];
	        this.make = source["make"];
	        this.model = source["model"];
	        this.alias = source["alias"];
	        this.basePi = source["basePi"];
	        this.drivetrainDefault = source["drivetrainDefault"];
	        this.source = source["source"];
	        this.sourceRef = source["sourceRef"];
	        this.createdAt = source["createdAt"];
	        this.updatedAt = source["updatedAt"];
	    }
	}
	export class GForcePoint {
	    timeMs: number;
	    xG: number;
	    yG: number;
	    zG: number;
	    totalG: number;
	
	    static createFrom(source: any = {}) {
	        return new GForcePoint(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.timeMs = source["timeMs"];
	        this.xG = source["xG"];
	        this.yG = source["yG"];
	        this.zG = source["zG"];
	        this.totalG = source["totalG"];
	    }
	}
	export class GForceDiagnostic {
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
	
	    static createFrom(source: any = {}) {
	        return new GForceDiagnostic(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.source = source["source"];
	        this.axisMapping = source["axisMapping"];
	        this.currentXG = source["currentXG"];
	        this.currentYG = source["currentYG"];
	        this.currentZG = source["currentZG"];
	        this.currentTotalG = source["currentTotalG"];
	        this.avgAbsXG = source["avgAbsXG"];
	        this.avgAbsYG = source["avgAbsYG"];
	        this.avgAbsZG = source["avgAbsZG"];
	        this.avgTotalG = source["avgTotalG"];
	        this.peakAbsXG = source["peakAbsXG"];
	        this.peakAbsYG = source["peakAbsYG"];
	        this.peakAbsZG = source["peakAbsZG"];
	        this.peakTotalG = source["peakTotalG"];
	        this.dominantAxis = source["dominantAxis"];
	        this.series = this.convertValues(source["series"], GForcePoint);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class GearPowerBand {
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
	
	    static createFrom(source: any = {}) {
	        return new GearPowerBand(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.gear = source["gear"];
	        this.sampleCount = source["sampleCount"];
	        this.highLoadSampleCount = source["highLoadSampleCount"];
	        this.speedMinKmh = source["speedMinKmh"];
	        this.speedMaxKmh = source["speedMaxKmh"];
	        this.speedAvgKmh = source["speedAvgKmh"];
	        this.rpmMin = source["rpmMin"];
	        this.rpmMax = source["rpmMax"];
	        this.rpmAvg = source["rpmAvg"];
	        this.rpmRatioMin = source["rpmRatioMin"];
	        this.rpmRatioMax = source["rpmRatioMax"];
	        this.rpmRatioAvg = source["rpmRatioAvg"];
	        this.inPowerBandRpmMin = source["inPowerBandRpmMin"];
	        this.inPowerBandRpmMax = source["inPowerBandRpmMax"];
	        this.inPowerBandRatioMin = source["inPowerBandRatioMin"];
	        this.inPowerBandRatioMax = source["inPowerBandRatioMax"];
	        this.throttleAvg = source["throttleAvg"];
	        this.accelAvgMps2 = source["accelAvgMps2"];
	        this.accelMaxMps2 = source["accelMaxMps2"];
	        this.speedPer1000RpmKmh = source["speedPer1000RpmKmh"];
	        this.shiftAfterRPM = source["shiftAfterRPM"];
	        this.shiftDropRPM = source["shiftDropRPM"];
	        this.frontSlipAvg = source["frontSlipAvg"];
	        this.rearSlipAvg = source["rearSlipAvg"];
	        this.frontTractionLimitedPct = source["frontTractionLimitedPct"];
	        this.rearTractionLimitedPct = source["rearTractionLimitedPct"];
	        this.belowPowerBandPercent = source["belowPowerBandPercent"];
	        this.inPowerBandPercent = source["inPowerBandPercent"];
	        this.abovePowerBandPercent = source["abovePowerBandPercent"];
	        this.lowRpmHighLoadPercent = source["lowRpmHighLoadPercent"];
	        this.highRpmHighLoadPercent = source["highRpmHighLoadPercent"];
	        this.tractionLimitedPercent = source["tractionLimitedPercent"];
	        this.finding = source["finding"];
	    }
	}
	export class GearPowerComparisonRow {
	    item: string;
	    gear?: number;
	    beforeValue?: number;
	    afterValue?: number;
	    deltaValue?: number;
	    beforeSpeedMaxKmh?: number;
	    afterSpeedMaxKmh?: number;
	    speedMaxDeltaKmh?: number;
	    beforeInPowerBandPct?: number;
	    afterInPowerBandPct?: number;
	    inPowerBandDeltaPct?: number;
	    beforeTractionLimitPct?: number;
	    afterTractionLimitPct?: number;
	    tractionLimitDeltaPct?: number;
	    beforeFinding?: string;
	    afterFinding?: string;
	
	    static createFrom(source: any = {}) {
	        return new GearPowerComparisonRow(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.item = source["item"];
	        this.gear = source["gear"];
	        this.beforeValue = source["beforeValue"];
	        this.afterValue = source["afterValue"];
	        this.deltaValue = source["deltaValue"];
	        this.beforeSpeedMaxKmh = source["beforeSpeedMaxKmh"];
	        this.afterSpeedMaxKmh = source["afterSpeedMaxKmh"];
	        this.speedMaxDeltaKmh = source["speedMaxDeltaKmh"];
	        this.beforeInPowerBandPct = source["beforeInPowerBandPct"];
	        this.afterInPowerBandPct = source["afterInPowerBandPct"];
	        this.inPowerBandDeltaPct = source["inPowerBandDeltaPct"];
	        this.beforeTractionLimitPct = source["beforeTractionLimitPct"];
	        this.afterTractionLimitPct = source["afterTractionLimitPct"];
	        this.tractionLimitDeltaPct = source["tractionLimitDeltaPct"];
	        this.beforeFinding = source["beforeFinding"];
	        this.afterFinding = source["afterFinding"];
	    }
	}
	export class GearPowerComparison {
	    type: string;
	    status: string;
	    baselineSessionId?: number;
	    rows: GearPowerComparisonRow[];
	
	    static createFrom(source: any = {}) {
	        return new GearPowerComparison(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.type = source["type"];
	        this.status = source["status"];
	        this.baselineSessionId = source["baselineSessionId"];
	        this.rows = this.convertValues(source["rows"], GearPowerComparisonRow);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class GearPowerDiagnostic {
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
	    recommendedActions: telemetry.SuggestedAction[];
	    evidence: Record<string, number>;
	
	    static createFrom(source: any = {}) {
	        return new GearPowerDiagnostic(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.status = source["status"];
	        this.summary = source["summary"];
	        this.launchFinding = source["launchFinding"];
	        this.topSpeedFinding = source["topSpeedFinding"];
	        this.powerKW = source["powerKW"];
	        this.torqueNM = source["torqueNM"];
	        this.weightKG = source["weightKG"];
	        this.frontWeightPct = source["frontWeightPct"];
	        this.powerToWeightKWPerKG = source["powerToWeightKWPerKG"];
	        this.powerToWeightBand = source["powerToWeightBand"];
	        this.peakTorqueRPM = source["peakTorqueRPM"];
	        this.peakPowerRPM = source["peakPowerRPM"];
	        this.redlineRPM = source["redlineRPM"];
	        this.powerBandStartRPM = source["powerBandStartRPM"];
	        this.powerBandEndRPM = source["powerBandEndRPM"];
	        this.powerBandSource = source["powerBandSource"];
	        this.confidence = source["confidence"];
	        this.strategyMode = source["strategyMode"];
	        this.globalGearIssueCount = source["globalGearIssueCount"];
	        this.usableGearCount = source["usableGearCount"];
	        this.globalGearIssueRatio = source["globalGearIssueRatio"];
	        this.tractionLimitedPercent = source["tractionLimitedPercent"];
	        this.lowRpmHighLoadPercent = source["lowRpmHighLoadPercent"];
	        this.highRpmHighLoadPercent = source["highRpmHighLoadPercent"];
	        this.gears = this.convertValues(source["gears"], GearPowerBand);
	        this.comparisons = this.convertValues(source["comparisons"], GearPowerComparison);
	        this.recommendedActions = this.convertValues(source["recommendedActions"], telemetry.SuggestedAction);
	        this.evidence = source["evidence"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class IssueEvidence {
	    min: number;
	    max: number;
	    avg: number;
	    count: number;
	
	    static createFrom(source: any = {}) {
	        return new IssueEvidence(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.min = source["min"];
	        this.max = source["max"];
	        this.avg = source["avg"];
	        this.count = source["count"];
	    }
	}
	export class PowerToTireDiagnostic {
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
	
	    static createFrom(source: any = {}) {
	        return new PowerToTireDiagnostic(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.status = source["status"];
	        this.summary = source["summary"];
	        this.explanation = source["explanation"];
	        this.confidence = source["confidence"];
	        this.sampleCount = source["sampleCount"];
	        this.highThrottleSampleCount = source["highThrottleSampleCount"];
	        this.drivetrain = source["drivetrain"];
	        this.drivenAxle = source["drivenAxle"];
	        this.currentPowerKW = source["currentPowerKW"];
	        this.averagePowerKW = source["averagePowerKW"];
	        this.maxPowerKW = source["maxPowerKW"];
	        this.currentTorqueNM = source["currentTorqueNM"];
	        this.averageTorqueNM = source["averageTorqueNM"];
	        this.maxTorqueNM = source["maxTorqueNM"];
	        this.currentRPM = source["currentRPM"];
	        this.averageRPM = source["averageRPM"];
	        this.currentRPMRatio = source["currentRPMRatio"];
	        this.averageRPMRatio = source["averageRPMRatio"];
	        this.currentGear = source["currentGear"];
	        this.averageThrottle = source["averageThrottle"];
	        this.averageSpeedKmh = source["averageSpeedKmh"];
	        this.speedDeltaKmh = source["speedDeltaKmh"];
	        this.averageAccelMps2 = source["averageAccelMps2"];
	        this.averageAccelG = source["averageAccelG"];
	        this.peakAccelG = source["peakAccelG"];
	        this.frontSlipRatioP90 = source["frontSlipRatioP90"];
	        this.rearSlipRatioP90 = source["rearSlipRatioP90"];
	        this.drivenSlipRatioP90 = source["drivenSlipRatioP90"];
	        this.drivenSlipRatioHighPct = source["drivenSlipRatioHighPct"];
	        this.rpmLowHighThrottlePct = source["rpmLowHighThrottlePct"];
	        this.rpmHighHighThrottlePct = source["rpmHighHighThrottlePct"];
	        this.powerSignalAvailable = source["powerSignalAvailable"];
	        this.tractionLimited = source["tractionLimited"];
	        this.evidence = source["evidence"];
	    }
	}
	export class ProfessionalPipelineConfig {
	    detectorId: string;
	    decisionerId: string;
	    interpreterId: string;
	
	    static createFrom(source: any = {}) {
	        return new ProfessionalPipelineConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.detectorId = source["detectorId"];
	        this.decisionerId = source["decisionerId"];
	        this.interpreterId = source["interpreterId"];
	    }
	}
	export class TuningAdvice {
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
	
	    static createFrom(source: any = {}) {
	        return new TuningAdvice(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.decisionId = source["decisionId"];
	        this.problemId = source["problemId"];
	        this.layer = source["layer"];
	        this.category = source["category"];
	        this.scope = source["scope"];
	        this.direction = source["direction"];
	        this.relatedFields = source["relatedFields"];
	        this.rationale = source["rationale"];
	        this.verifyEvidence = source["verifyEvidence"];
	        this.trustLevel = source["trustLevel"];
	        this.missingInputs = source["missingInputs"];
	        this.conflictReason = source["conflictReason"];
	        this.canApply = source["canApply"];
	        this.documentSources = source["documentSources"];
	        this.evidence = source["evidence"];
	    }
	}
	export class TuningAdviceSet {
	    interpreterId: string;
	    status: string;
	    advice: TuningAdvice[];
	    documentSources: string[];
	    warnings: string[];
	
	    static createFrom(source: any = {}) {
	        return new TuningAdviceSet(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.interpreterId = source["interpreterId"];
	        this.status = source["status"];
	        this.advice = this.convertValues(source["advice"], TuningAdvice);
	        this.documentSources = source["documentSources"];
	        this.warnings = source["warnings"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class TuningDecision {
	    id: string;
	    problemId: string;
	    phase: string;
	    primaryCause: string;
	    shouldTune: boolean;
	    confidence: string;
	    rationale: string;
	    documentContext: string;
	    evidence: Record<string, number>;
	
	    static createFrom(source: any = {}) {
	        return new TuningDecision(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.problemId = source["problemId"];
	        this.phase = source["phase"];
	        this.primaryCause = source["primaryCause"];
	        this.shouldTune = source["shouldTune"];
	        this.confidence = source["confidence"];
	        this.rationale = source["rationale"];
	        this.documentContext = source["documentContext"];
	        this.evidence = source["evidence"];
	    }
	}
	export class TuningDecisionSet {
	    decisionerId: string;
	    status: string;
	    decisions: TuningDecision[];
	    warnings: string[];
	
	    static createFrom(source: any = {}) {
	        return new TuningDecisionSet(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.decisionerId = source["decisionerId"];
	        this.status = source["status"];
	        this.decisions = this.convertValues(source["decisions"], TuningDecision);
	        this.warnings = source["warnings"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class TuningProblem {
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
	
	    static createFrom(source: any = {}) {
	        return new TuningProblem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.sourceId = source["sourceId"];
	        this.family = source["family"];
	        this.type = source["type"];
	        this.phase = source["phase"];
	        this.operationTags = source["operationTags"];
	        this.limitedAxle = source["limitedAxle"];
	        this.limitedWheels = source["limitedWheels"];
	        this.severity = source["severity"];
	        this.confidence = source["confidence"];
	        this.riskLevel = source["riskLevel"];
	        this.count = source["count"];
	        this.durationMs = source["durationMs"];
	        this.summary = source["summary"];
	        this.reason = source["reason"];
	        this.evidence = source["evidence"];
	    }
	}
	export class TuningProblemSet {
	    detectorId: string;
	    status: string;
	    problems: TuningProblem[];
	    warnings: string[];
	
	    static createFrom(source: any = {}) {
	        return new TuningProblemSet(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.detectorId = source["detectorId"];
	        this.status = source["status"];
	        this.problems = this.convertValues(source["problems"], TuningProblem);
	        this.warnings = source["warnings"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class SessionVehicleSnapshot {
	    carOrdinal?: number;
	    carClass: string;
	    carPi?: number;
	    drivetrain: string;
	    numCylinders?: number;
	
	    static createFrom(source: any = {}) {
	        return new SessionVehicleSnapshot(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.carOrdinal = source["carOrdinal"];
	        this.carClass = source["carClass"];
	        this.carPi = source["carPi"];
	        this.drivetrain = source["drivetrain"];
	        this.numCylinders = source["numCylinders"];
	    }
	}
	export class TuningPipelineSourceSummary {
	    sourceType: string;
	    sessionId?: number;
	    sampleCount: number;
	    eventCount: number;
	    vehicle: SessionVehicleSnapshot;
	    gameMode: string;
	    driverMode: string;
	    label: string;
	
	    static createFrom(source: any = {}) {
	        return new TuningPipelineSourceSummary(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.sourceType = source["sourceType"];
	        this.sessionId = source["sessionId"];
	        this.sampleCount = source["sampleCount"];
	        this.eventCount = source["eventCount"];
	        this.vehicle = this.convertValues(source["vehicle"], SessionVehicleSnapshot);
	        this.gameMode = source["gameMode"];
	        this.driverMode = source["driverMode"];
	        this.label = source["label"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class TuningPipelineRunResult {
	    status: string;
	    updatedAt: string;
	    sourceSummary: TuningPipelineSourceSummary;
	    problemSet: TuningProblemSet;
	    decisionSet: TuningDecisionSet;
	    adviceSet: TuningAdviceSet;
	    warnings: string[];
	
	    static createFrom(source: any = {}) {
	        return new TuningPipelineRunResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.status = source["status"];
	        this.updatedAt = source["updatedAt"];
	        this.sourceSummary = this.convertValues(source["sourceSummary"], TuningPipelineSourceSummary);
	        this.problemSet = this.convertValues(source["problemSet"], TuningProblemSet);
	        this.decisionSet = this.convertValues(source["decisionSet"], TuningDecisionSet);
	        this.adviceSet = this.convertValues(source["adviceSet"], TuningAdviceSet);
	        this.warnings = source["warnings"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ProfessionalTuningDiagnostic {
	    status: string;
	    updatedAt: string;
	    config: ProfessionalPipelineConfig;
	    pipeline?: TuningPipelineRunResult;
	    warnings: string[];
	
	    static createFrom(source: any = {}) {
	        return new ProfessionalTuningDiagnostic(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.status = source["status"];
	        this.updatedAt = source["updatedAt"];
	        this.config = this.convertValues(source["config"], ProfessionalPipelineConfig);
	        this.pipeline = this.convertValues(source["pipeline"], TuningPipelineRunResult);
	        this.warnings = source["warnings"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class QuickComparability {
	    sameVehicleClass: string;
	    sameTrackContext: string;
	    confidence: string;
	    warnings: string[];
	    baselineVehicle: SessionVehicleSnapshot;
	    currentVehicle: SessionVehicleSnapshot;
	
	    static createFrom(source: any = {}) {
	        return new QuickComparability(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.sameVehicleClass = source["sameVehicleClass"];
	        this.sameTrackContext = source["sameTrackContext"];
	        this.confidence = source["confidence"];
	        this.warnings = source["warnings"];
	        this.baselineVehicle = this.convertValues(source["baselineVehicle"], SessionVehicleSnapshot);
	        this.currentVehicle = this.convertValues(source["currentVehicle"], SessionVehicleSnapshot);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class QuickSuggestion {
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
	
	    static createFrom(source: any = {}) {
	        return new QuickSuggestion(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.family = source["family"];
	        this.source = source["source"];
	        this.confidence = source["confidence"];
	        this.trustLevel = source["trustLevel"];
	        this.adviceLayer = source["adviceLayer"];
	        this.category = source["category"];
	        this.item = source["item"];
	        this.direction = source["direction"];
	        this.amount = source["amount"];
	        this.reason = source["reason"];
	        this.rationale = source["rationale"];
	        this.nextStep = source["nextStep"];
	        this.fieldKeys = source["fieldKeys"];
	        this.missingInputs = source["missingInputs"];
	        this.canApply = source["canApply"];
	        this.blockedReason = source["blockedReason"];
	    }
	}
	export class SessionIssueGroup {
	    id: string;
	    family: string;
	    severity: string;
	    segment: string;
	    eventTypes: string[];
	    eventIds: string[];
	    events: telemetry.DetectedEvent[];
	    eventCount: number;
	    totalDurationMs: number;
	    firstStartMs: number;
	    lastEndMs: number;
	    evidence: Record<string, IssueEvidence>;
	    primaryActions: telemetry.SuggestedAction[];
	    comparison: string;
	    baselineEventCount: number;
	    baselineTotalDurationMs: number;
	    relatedRecentChanges: string[];
	    prioritizeTuning: boolean;
	    adjustmentStrategy: string;
	    feedbackDirective: string;
	
	    static createFrom(source: any = {}) {
	        return new SessionIssueGroup(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.family = source["family"];
	        this.severity = source["severity"];
	        this.segment = source["segment"];
	        this.eventTypes = source["eventTypes"];
	        this.eventIds = source["eventIds"];
	        this.events = this.convertValues(source["events"], telemetry.DetectedEvent);
	        this.eventCount = source["eventCount"];
	        this.totalDurationMs = source["totalDurationMs"];
	        this.firstStartMs = source["firstStartMs"];
	        this.lastEndMs = source["lastEndMs"];
	        this.evidence = this.convertValues(source["evidence"], IssueEvidence, true);
	        this.primaryActions = this.convertValues(source["primaryActions"], telemetry.SuggestedAction);
	        this.comparison = source["comparison"];
	        this.baselineEventCount = source["baselineEventCount"];
	        this.baselineTotalDurationMs = source["baselineTotalDurationMs"];
	        this.relatedRecentChanges = source["relatedRecentChanges"];
	        this.prioritizeTuning = source["prioritizeTuning"];
	        this.adjustmentStrategy = source["adjustmentStrategy"];
	        this.feedbackDirective = source["feedbackDirective"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class QuickLapSummary {
	    lapNumber: number;
	    sampleCount: number;
	    durationMs: number;
	    avgSpeedKmh: number;
	    maxSpeedKmh: number;
	    eventCount: number;
	    issueScore: number;
	
	    static createFrom(source: any = {}) {
	        return new QuickLapSummary(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.lapNumber = source["lapNumber"];
	        this.sampleCount = source["sampleCount"];
	        this.durationMs = source["durationMs"];
	        this.avgSpeedKmh = source["avgSpeedKmh"];
	        this.maxSpeedKmh = source["maxSpeedKmh"];
	        this.eventCount = source["eventCount"];
	        this.issueScore = source["issueScore"];
	    }
	}
	export class QuickDiagnostic {
	    status: string;
	    comparisonStatus: string;
	    updatedAt: string;
	    sampleCount: number;
	    eventCount: number;
	    gameMode: string;
	    driverMode: string;
	    driverModeConfidence: number;
	    vehicle: SessionVehicleSnapshot;
	    comparability: QuickComparability;
	    currentLap?: QuickLapSummary;
	    previousLap?: QuickLapSummary;
	    groups: SessionIssueGroup[];
	    gearPower: GearPowerDiagnostic;
	    suggestions: QuickSuggestion[];
	    missingProfileFields: string[];
	
	    static createFrom(source: any = {}) {
	        return new QuickDiagnostic(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.status = source["status"];
	        this.comparisonStatus = source["comparisonStatus"];
	        this.updatedAt = source["updatedAt"];
	        this.sampleCount = source["sampleCount"];
	        this.eventCount = source["eventCount"];
	        this.gameMode = source["gameMode"];
	        this.driverMode = source["driverMode"];
	        this.driverModeConfidence = source["driverModeConfidence"];
	        this.vehicle = this.convertValues(source["vehicle"], SessionVehicleSnapshot);
	        this.comparability = this.convertValues(source["comparability"], QuickComparability);
	        this.currentLap = this.convertValues(source["currentLap"], QuickLapSummary);
	        this.previousLap = this.convertValues(source["previousLap"], QuickLapSummary);
	        this.groups = this.convertValues(source["groups"], SessionIssueGroup);
	        this.gearPower = this.convertValues(source["gearPower"], GearPowerDiagnostic);
	        this.suggestions = this.convertValues(source["suggestions"], QuickSuggestion);
	        this.missingProfileFields = source["missingProfileFields"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	
	export class RecommendedCar {
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
	    createdAt: string;
	    updatedAt: string;
	
	    static createFrom(source: any = {}) {
	        return new RecommendedCar(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.useCase = source["useCase"];
	        this.useCaseLabel = source["useCaseLabel"];
	        this.pi = source["pi"];
	        this.carClass = source["carClass"];
	        this.drivetrain = source["drivetrain"];
	        this.tireCompound = source["tireCompound"];
	        this.tireCompoundLabel = source["tireCompoundLabel"];
	        this.weightKG = source["weightKG"];
	        this.frontWeightPct = source["frontWeightPct"];
	        this.tuneCode = source["tuneCode"];
	        this.imageSrc = source["imageSrc"];
	        this.tags = source["tags"];
	        this.reason = source["reason"];
	        this.createdAt = source["createdAt"];
	        this.updatedAt = source["updatedAt"];
	    }
	}
	export class RecommendedCarInput {
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
	
	    static createFrom(source: any = {}) {
	        return new RecommendedCarInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.useCase = source["useCase"];
	        this.useCaseLabel = source["useCaseLabel"];
	        this.pi = source["pi"];
	        this.carClass = source["carClass"];
	        this.drivetrain = source["drivetrain"];
	        this.tireCompound = source["tireCompound"];
	        this.tireCompoundLabel = source["tireCompoundLabel"];
	        this.weightKG = source["weightKG"];
	        this.frontWeightPct = source["frontWeightPct"];
	        this.tuneCode = source["tuneCode"];
	        this.imageSrc = source["imageSrc"];
	        this.tags = source["tags"];
	        this.reason = source["reason"];
	    }
	}
	export class RetestMetric {
	    key: string;
	    current: number;
	    baseline: number;
	    delta: number;
	    direction: string;
	    status: string;
	
	    static createFrom(source: any = {}) {
	        return new RetestMetric(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.key = source["key"];
	        this.current = source["current"];
	        this.baseline = source["baseline"];
	        this.delta = source["delta"];
	        this.direction = source["direction"];
	        this.status = source["status"];
	    }
	}
	export class TunePlanDraftAction {
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
	
	    static createFrom(source: any = {}) {
	        return new TunePlanDraftAction(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.family = source["family"];
	        this.source = source["source"];
	        this.confidence = source["confidence"];
	        this.adviceLayer = source["adviceLayer"];
	        this.trustLevel = source["trustLevel"];
	        this.trustReasons = source["trustReasons"];
	        this.missingInputs = source["missingInputs"];
	        this.retestGuard = source["retestGuard"];
	        this.rationale = source["rationale"];
	        this.conflictReason = source["conflictReason"];
	        this.category = source["category"];
	        this.item = source["item"];
	        this.fieldKey = source["fieldKey"];
	        this.direction = source["direction"];
	        this.reason = source["reason"];
	        this.currentValue = source["currentValue"];
	        this.targetValue = source["targetValue"];
	        this.delta = source["delta"];
	        this.unit = source["unit"];
	        this.step = source["step"];
	        this.canApply = source["canApply"];
	        this.blockedReason = source["blockedReason"];
	    }
	}
	export class TelemetrySession {
	    id: number;
	    tuneProfileId?: number;
	    tuneSnapshotJson: string;
	    tuneName: string;
	    sessionName: string;
	    trackName: string;
	    mode: string;
	    gameMode: string;
	    startedAt: string;
	    endedAt: string;
	    durationMs: number;
	    bestLapMs?: number;
	    avgSpeedKmh?: number;
	    maxSpeedKmh?: number;
	    eventCount: number;
	    sampleCount: number;
	    recordingPath: string;
	    recordingPackets: number;
	    recordingBytes: number;
	    recordingTruncated: boolean;
	    carOrdinal?: number;
	    carClass: string;
	    carPi?: number;
	    drivetrain: string;
	    numCylinders?: number;
	    driverMode: string;
	    driverModeConfidence: number;
	    driverModeEvidenceJson: string;
	    brakeAssist: string;
	    steeringAssist: string;
	    tractionControl: string;
	    stabilityControl: string;
	    shifting: string;
	    launchControl: string;
	    driverFeedbackJson: string;
	    notes: string;
	
	    static createFrom(source: any = {}) {
	        return new TelemetrySession(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.tuneProfileId = source["tuneProfileId"];
	        this.tuneSnapshotJson = source["tuneSnapshotJson"];
	        this.tuneName = source["tuneName"];
	        this.sessionName = source["sessionName"];
	        this.trackName = source["trackName"];
	        this.mode = source["mode"];
	        this.gameMode = source["gameMode"];
	        this.startedAt = source["startedAt"];
	        this.endedAt = source["endedAt"];
	        this.durationMs = source["durationMs"];
	        this.bestLapMs = source["bestLapMs"];
	        this.avgSpeedKmh = source["avgSpeedKmh"];
	        this.maxSpeedKmh = source["maxSpeedKmh"];
	        this.eventCount = source["eventCount"];
	        this.sampleCount = source["sampleCount"];
	        this.recordingPath = source["recordingPath"];
	        this.recordingPackets = source["recordingPackets"];
	        this.recordingBytes = source["recordingBytes"];
	        this.recordingTruncated = source["recordingTruncated"];
	        this.carOrdinal = source["carOrdinal"];
	        this.carClass = source["carClass"];
	        this.carPi = source["carPi"];
	        this.drivetrain = source["drivetrain"];
	        this.numCylinders = source["numCylinders"];
	        this.driverMode = source["driverMode"];
	        this.driverModeConfidence = source["driverModeConfidence"];
	        this.driverModeEvidenceJson = source["driverModeEvidenceJson"];
	        this.brakeAssist = source["brakeAssist"];
	        this.steeringAssist = source["steeringAssist"];
	        this.tractionControl = source["tractionControl"];
	        this.stabilityControl = source["stabilityControl"];
	        this.shifting = source["shifting"];
	        this.launchControl = source["launchControl"];
	        this.driverFeedbackJson = source["driverFeedbackJson"];
	        this.notes = source["notes"];
	    }
	}
	export class RetestEvaluation {
	    sessionId: number;
	    baselineSession?: TelemetrySession;
	    status: string;
	    summary: string;
	    confidence: string;
	    baselineReason: string;
	    changedFields: string[];
	    changeSourceSessionId?: number;
	    rollbackActions: TunePlanDraftAction[];
	    metricSummary: string[];
	    metrics: RetestMetric[];
	
	    static createFrom(source: any = {}) {
	        return new RetestEvaluation(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.sessionId = source["sessionId"];
	        this.baselineSession = this.convertValues(source["baselineSession"], TelemetrySession);
	        this.status = source["status"];
	        this.summary = source["summary"];
	        this.confidence = source["confidence"];
	        this.baselineReason = source["baselineReason"];
	        this.changedFields = source["changedFields"];
	        this.changeSourceSessionId = source["changeSourceSessionId"];
	        this.rollbackActions = this.convertValues(source["rollbackActions"], TunePlanDraftAction);
	        this.metricSummary = source["metricSummary"];
	        this.metrics = this.convertValues(source["metrics"], RetestMetric);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class RoadEvaluationAttribution {
	    type: string;
	    eventType?: string;
	    count: number;
	    severity?: string;
	    priority: number;
	    message: string;
	    prioritizeTuning: boolean;
	
	    static createFrom(source: any = {}) {
	        return new RoadEvaluationAttribution(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.type = source["type"];
	        this.eventType = source["eventType"];
	        this.count = source["count"];
	        this.severity = source["severity"];
	        this.priority = source["priority"];
	        this.message = source["message"];
	        this.prioritizeTuning = source["prioritizeTuning"];
	    }
	}
	export class SessionComparisonMetric {
	    key: string;
	    label: string;
	    unit: string;
	    left: number;
	    right: number;
	    delta: number;
	    higherIsBetter: boolean;
	
	    static createFrom(source: any = {}) {
	        return new SessionComparisonMetric(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.key = source["key"];
	        this.label = source["label"];
	        this.unit = source["unit"];
	        this.left = source["left"];
	        this.right = source["right"];
	        this.delta = source["delta"];
	        this.higherIsBetter = source["higherIsBetter"];
	    }
	}
	export class RoadSessionEvaluation {
	    session: TelemetrySession;
	    track?: BenchmarkTrack;
	    bestRun?: BenchmarkRun;
	    baselineRun?: BenchmarkRun;
	    baselineSession?: TelemetrySession;
	    baselineStatus: string;
	    paperPerformanceScore: number;
	    playerFitScore: number;
	    riskScore: number;
	    overallVerdict: string;
	    attributions: RoadEvaluationAttribution[];
	    notes: string[];
	
	    static createFrom(source: any = {}) {
	        return new RoadSessionEvaluation(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.session = this.convertValues(source["session"], TelemetrySession);
	        this.track = this.convertValues(source["track"], BenchmarkTrack);
	        this.bestRun = this.convertValues(source["bestRun"], BenchmarkRun);
	        this.baselineRun = this.convertValues(source["baselineRun"], BenchmarkRun);
	        this.baselineSession = this.convertValues(source["baselineSession"], TelemetrySession);
	        this.baselineStatus = source["baselineStatus"];
	        this.paperPerformanceScore = source["paperPerformanceScore"];
	        this.playerFitScore = source["playerFitScore"];
	        this.riskScore = source["riskScore"];
	        this.overallVerdict = source["overallVerdict"];
	        this.attributions = this.convertValues(source["attributions"], RoadEvaluationAttribution);
	        this.notes = source["notes"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class RoadEvaluationComparison {
	    left: RoadSessionEvaluation;
	    right: RoadSessionEvaluation;
	    metrics: SessionComparisonMetric[];
	    verdict: string;
	    notes: string[];
	
	    static createFrom(source: any = {}) {
	        return new RoadEvaluationComparison(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.left = this.convertValues(source["left"], RoadSessionEvaluation);
	        this.right = this.convertValues(source["right"], RoadSessionEvaluation);
	        this.metrics = this.convertValues(source["metrics"], SessionComparisonMetric);
	        this.verdict = source["verdict"];
	        this.notes = source["notes"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class RoadStaticTuneBaselineInput {
	    carName: string;
	    versionName: string;
	    useCase?: string;
	    carOrdinal?: number;
	    carCategory?: number;
	    pi: number;
	    drivetrain: string;
	    tireCompound?: string;
	    weightKG: number;
	    frontWeightPct: number;
	    powerKW?: number;
	    torqueNM?: number;
	    redlineRPM?: number;
	    gearCount?: number;
	    tireDiameterCm?: number;
	    targetTopSpeedKmh?: number;
	    frontRideHeightMinCm?: number;
	    frontRideHeightMaxCm?: number;
	    rearRideHeightMinCm?: number;
	    rearRideHeightMaxCm?: number;
	    frontAeroMinKgf?: number;
	    frontAeroMaxKgf?: number;
	    rearAeroMinKgf?: number;
	    rearAeroMaxKgf?: number;
	    frontRideHeightAdjustable?: boolean;
	    rearRideHeightAdjustable?: boolean;
	    frontAeroAdjustable?: boolean;
	    rearAeroAdjustable?: boolean;
	    balanceBias?: number;
	    stiffnessBias?: number;
	    speedBias?: number;
	
	    static createFrom(source: any = {}) {
	        return new RoadStaticTuneBaselineInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.carName = source["carName"];
	        this.versionName = source["versionName"];
	        this.useCase = source["useCase"];
	        this.carOrdinal = source["carOrdinal"];
	        this.carCategory = source["carCategory"];
	        this.pi = source["pi"];
	        this.drivetrain = source["drivetrain"];
	        this.tireCompound = source["tireCompound"];
	        this.weightKG = source["weightKG"];
	        this.frontWeightPct = source["frontWeightPct"];
	        this.powerKW = source["powerKW"];
	        this.torqueNM = source["torqueNM"];
	        this.redlineRPM = source["redlineRPM"];
	        this.gearCount = source["gearCount"];
	        this.tireDiameterCm = source["tireDiameterCm"];
	        this.targetTopSpeedKmh = source["targetTopSpeedKmh"];
	        this.frontRideHeightMinCm = source["frontRideHeightMinCm"];
	        this.frontRideHeightMaxCm = source["frontRideHeightMaxCm"];
	        this.rearRideHeightMinCm = source["rearRideHeightMinCm"];
	        this.rearRideHeightMaxCm = source["rearRideHeightMaxCm"];
	        this.frontAeroMinKgf = source["frontAeroMinKgf"];
	        this.frontAeroMaxKgf = source["frontAeroMaxKgf"];
	        this.rearAeroMinKgf = source["rearAeroMinKgf"];
	        this.rearAeroMaxKgf = source["rearAeroMaxKgf"];
	        this.frontRideHeightAdjustable = source["frontRideHeightAdjustable"];
	        this.rearRideHeightAdjustable = source["rearRideHeightAdjustable"];
	        this.frontAeroAdjustable = source["frontAeroAdjustable"];
	        this.rearAeroAdjustable = source["rearAeroAdjustable"];
	        this.balanceBias = source["balanceBias"];
	        this.stiffnessBias = source["stiffnessBias"];
	        this.speedBias = source["speedBias"];
	    }
	}
	export class RoadStaticTuneBaselineApplyInput {
	    createNew: boolean;
	    targetProfileId: number;
	    baselineInput: RoadStaticTuneBaselineInput;
	    selectedFieldKeys: string[];
	
	    static createFrom(source: any = {}) {
	        return new RoadStaticTuneBaselineApplyInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.createNew = source["createNew"];
	        this.targetProfileId = source["targetProfileId"];
	        this.baselineInput = this.convertValues(source["baselineInput"], RoadStaticTuneBaselineInput);
	        this.selectedFieldKeys = source["selectedFieldKeys"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class TuneProfile {
	    id: number;
	    createdAt: string;
	    updatedAt: string;
	    carName: string;
	    carOrdinal?: number;
	    carCategory?: number;
	    carClass: string;
	    pi?: number;
	    drivetrain: string;
	    numCylinders?: number;
	    useCase: string;
	    versionName: string;
	    powerKW?: number;
	    torqueNM?: number;
	    weightKG?: number;
	    frontWeightPct?: number;
	    powerToWeightKWPerKG?: number;
	    peakTorqueRPM?: number;
	    peakPowerRPM?: number;
	    redlineRPM?: number;
	    frontTirePressure?: number;
	    rearTirePressure?: number;
	    finalDrive?: number;
	    gear1?: number;
	    gear2?: number;
	    gear3?: number;
	    gear4?: number;
	    gear5?: number;
	    gear6?: number;
	    gear7?: number;
	    gear8?: number;
	    gear9?: number;
	    gear10?: number;
	    frontCamber?: number;
	    rearCamber?: number;
	    frontToe?: number;
	    rearToe?: number;
	    caster?: number;
	    frontArb?: number;
	    rearArb?: number;
	    frontSpring?: number;
	    rearSpring?: number;
	    frontRideHeight?: number;
	    rearRideHeight?: number;
	    frontRebound?: number;
	    rearRebound?: number;
	    frontBump?: number;
	    rearBump?: number;
	    frontAero?: number;
	    rearAero?: number;
	    aeroBalance?: number;
	    brakeBalance?: number;
	    brakePressure?: number;
	    frontDiffAccel?: number;
	    frontDiffDecel?: number;
	    rearDiffAccel?: number;
	    rearDiffDecel?: number;
	    centerDiffBalance?: number;
	    notes: string;
	
	    static createFrom(source: any = {}) {
	        return new TuneProfile(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.createdAt = source["createdAt"];
	        this.updatedAt = source["updatedAt"];
	        this.carName = source["carName"];
	        this.carOrdinal = source["carOrdinal"];
	        this.carCategory = source["carCategory"];
	        this.carClass = source["carClass"];
	        this.pi = source["pi"];
	        this.drivetrain = source["drivetrain"];
	        this.numCylinders = source["numCylinders"];
	        this.useCase = source["useCase"];
	        this.versionName = source["versionName"];
	        this.powerKW = source["powerKW"];
	        this.torqueNM = source["torqueNM"];
	        this.weightKG = source["weightKG"];
	        this.frontWeightPct = source["frontWeightPct"];
	        this.powerToWeightKWPerKG = source["powerToWeightKWPerKG"];
	        this.peakTorqueRPM = source["peakTorqueRPM"];
	        this.peakPowerRPM = source["peakPowerRPM"];
	        this.redlineRPM = source["redlineRPM"];
	        this.frontTirePressure = source["frontTirePressure"];
	        this.rearTirePressure = source["rearTirePressure"];
	        this.finalDrive = source["finalDrive"];
	        this.gear1 = source["gear1"];
	        this.gear2 = source["gear2"];
	        this.gear3 = source["gear3"];
	        this.gear4 = source["gear4"];
	        this.gear5 = source["gear5"];
	        this.gear6 = source["gear6"];
	        this.gear7 = source["gear7"];
	        this.gear8 = source["gear8"];
	        this.gear9 = source["gear9"];
	        this.gear10 = source["gear10"];
	        this.frontCamber = source["frontCamber"];
	        this.rearCamber = source["rearCamber"];
	        this.frontToe = source["frontToe"];
	        this.rearToe = source["rearToe"];
	        this.caster = source["caster"];
	        this.frontArb = source["frontArb"];
	        this.rearArb = source["rearArb"];
	        this.frontSpring = source["frontSpring"];
	        this.rearSpring = source["rearSpring"];
	        this.frontRideHeight = source["frontRideHeight"];
	        this.rearRideHeight = source["rearRideHeight"];
	        this.frontRebound = source["frontRebound"];
	        this.rearRebound = source["rearRebound"];
	        this.frontBump = source["frontBump"];
	        this.rearBump = source["rearBump"];
	        this.frontAero = source["frontAero"];
	        this.rearAero = source["rearAero"];
	        this.aeroBalance = source["aeroBalance"];
	        this.brakeBalance = source["brakeBalance"];
	        this.brakePressure = source["brakePressure"];
	        this.frontDiffAccel = source["frontDiffAccel"];
	        this.frontDiffDecel = source["frontDiffDecel"];
	        this.rearDiffAccel = source["rearDiffAccel"];
	        this.rearDiffDecel = source["rearDiffDecel"];
	        this.centerDiffBalance = source["centerDiffBalance"];
	        this.notes = source["notes"];
	    }
	}
	export class RoadStaticTuneBaselineApplyResult {
	    profile: TuneProfile;
	    appliedFields: string[];
	    skippedFields: BaselineSkippedField[];
	
	    static createFrom(source: any = {}) {
	        return new RoadStaticTuneBaselineApplyResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.profile = this.convertValues(source["profile"], TuneProfile);
	        this.appliedFields = source["appliedFields"];
	        this.skippedFields = this.convertValues(source["skippedFields"], BaselineSkippedField);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class TuneProfileInput {
	    carName: string;
	    carOrdinal?: number;
	    carCategory?: number;
	    carClass: string;
	    pi?: number;
	    drivetrain: string;
	    numCylinders?: number;
	    useCase: string;
	    versionName: string;
	    powerKW?: number;
	    torqueNM?: number;
	    weightKG?: number;
	    frontWeightPct?: number;
	    powerToWeightKWPerKG?: number;
	    peakTorqueRPM?: number;
	    peakPowerRPM?: number;
	    redlineRPM?: number;
	    frontTirePressure?: number;
	    rearTirePressure?: number;
	    finalDrive?: number;
	    gear1?: number;
	    gear2?: number;
	    gear3?: number;
	    gear4?: number;
	    gear5?: number;
	    gear6?: number;
	    gear7?: number;
	    gear8?: number;
	    gear9?: number;
	    gear10?: number;
	    frontCamber?: number;
	    rearCamber?: number;
	    frontToe?: number;
	    rearToe?: number;
	    caster?: number;
	    frontArb?: number;
	    rearArb?: number;
	    frontSpring?: number;
	    rearSpring?: number;
	    frontRideHeight?: number;
	    rearRideHeight?: number;
	    frontRebound?: number;
	    rearRebound?: number;
	    frontBump?: number;
	    rearBump?: number;
	    frontAero?: number;
	    rearAero?: number;
	    aeroBalance?: number;
	    brakeBalance?: number;
	    brakePressure?: number;
	    frontDiffAccel?: number;
	    frontDiffDecel?: number;
	    rearDiffAccel?: number;
	    rearDiffDecel?: number;
	    centerDiffBalance?: number;
	    notes: string;
	
	    static createFrom(source: any = {}) {
	        return new TuneProfileInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.carName = source["carName"];
	        this.carOrdinal = source["carOrdinal"];
	        this.carCategory = source["carCategory"];
	        this.carClass = source["carClass"];
	        this.pi = source["pi"];
	        this.drivetrain = source["drivetrain"];
	        this.numCylinders = source["numCylinders"];
	        this.useCase = source["useCase"];
	        this.versionName = source["versionName"];
	        this.powerKW = source["powerKW"];
	        this.torqueNM = source["torqueNM"];
	        this.weightKG = source["weightKG"];
	        this.frontWeightPct = source["frontWeightPct"];
	        this.powerToWeightKWPerKG = source["powerToWeightKWPerKG"];
	        this.peakTorqueRPM = source["peakTorqueRPM"];
	        this.peakPowerRPM = source["peakPowerRPM"];
	        this.redlineRPM = source["redlineRPM"];
	        this.frontTirePressure = source["frontTirePressure"];
	        this.rearTirePressure = source["rearTirePressure"];
	        this.finalDrive = source["finalDrive"];
	        this.gear1 = source["gear1"];
	        this.gear2 = source["gear2"];
	        this.gear3 = source["gear3"];
	        this.gear4 = source["gear4"];
	        this.gear5 = source["gear5"];
	        this.gear6 = source["gear6"];
	        this.gear7 = source["gear7"];
	        this.gear8 = source["gear8"];
	        this.gear9 = source["gear9"];
	        this.gear10 = source["gear10"];
	        this.frontCamber = source["frontCamber"];
	        this.rearCamber = source["rearCamber"];
	        this.frontToe = source["frontToe"];
	        this.rearToe = source["rearToe"];
	        this.caster = source["caster"];
	        this.frontArb = source["frontArb"];
	        this.rearArb = source["rearArb"];
	        this.frontSpring = source["frontSpring"];
	        this.rearSpring = source["rearSpring"];
	        this.frontRideHeight = source["frontRideHeight"];
	        this.rearRideHeight = source["rearRideHeight"];
	        this.frontRebound = source["frontRebound"];
	        this.rearRebound = source["rearRebound"];
	        this.frontBump = source["frontBump"];
	        this.rearBump = source["rearBump"];
	        this.frontAero = source["frontAero"];
	        this.rearAero = source["rearAero"];
	        this.aeroBalance = source["aeroBalance"];
	        this.brakeBalance = source["brakeBalance"];
	        this.brakePressure = source["brakePressure"];
	        this.frontDiffAccel = source["frontDiffAccel"];
	        this.frontDiffDecel = source["frontDiffDecel"];
	        this.rearDiffAccel = source["rearDiffAccel"];
	        this.rearDiffDecel = source["rearDiffDecel"];
	        this.centerDiffBalance = source["centerDiffBalance"];
	        this.notes = source["notes"];
	    }
	}
	export class RoadStaticTuneBaselineResult {
	    profileDraft: TuneProfileInput;
	    confidence: string;
	    generatedFields: BaselineGeneratedField[];
	    tierRecommendations: BaselineTierRecommendation[];
	    skippedFields: BaselineSkippedField[];
	    warnings: string[];
	    nextTestPlan: string[];
	
	    static createFrom(source: any = {}) {
	        return new RoadStaticTuneBaselineResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.profileDraft = this.convertValues(source["profileDraft"], TuneProfileInput);
	        this.confidence = source["confidence"];
	        this.generatedFields = this.convertValues(source["generatedFields"], BaselineGeneratedField);
	        this.tierRecommendations = this.convertValues(source["tierRecommendations"], BaselineTierRecommendation);
	        this.skippedFields = this.convertValues(source["skippedFields"], BaselineSkippedField);
	        this.warnings = source["warnings"];
	        this.nextTestPlan = source["nextTestPlan"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class StrategyAnalysisHint {
	    level: string;
	    message: string;
	    eventType?: string;
	    family?: string;
	
	    static createFrom(source: any = {}) {
	        return new StrategyAnalysisHint(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.level = source["level"];
	        this.message = source["message"];
	        this.eventType = source["eventType"];
	        this.family = source["family"];
	    }
	}
	export class StrategyIssueAggregate {
	    family: string;
	    eventCount: number;
	    sessionCount: number;
	    severity: string;
	    recommendation: string;
	
	    static createFrom(source: any = {}) {
	        return new StrategyIssueAggregate(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.family = source["family"];
	        this.eventCount = source["eventCount"];
	        this.sessionCount = source["sessionCount"];
	        this.severity = source["severity"];
	        this.recommendation = source["recommendation"];
	    }
	}
	export class StrategyEventDistribution {
	    type: string;
	    count: number;
	    severity: string;
	
	    static createFrom(source: any = {}) {
	        return new StrategyEventDistribution(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.type = source["type"];
	        this.count = source["count"];
	        this.severity = source["severity"];
	    }
	}
	export class StrategyTemplate {
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
	
	    static createFrom(source: any = {}) {
	        return new StrategyTemplate(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.carClass = source["carClass"];
	        this.drivetrain = source["drivetrain"];
	        this.useCase = source["useCase"];
	        this.gameMode = source["gameMode"];
	        this.isDefault = source["isDefault"];
	        this.enabledEventCount = source["enabledEventCount"];
	        this.totalEventCount = source["totalEventCount"];
	        this.updatedAt = source["updatedAt"];
	    }
	}
	export class RoadStrategyAnalysis {
	    template: StrategyTemplate;
	    sessionIds: number[];
	    sessionCount: number;
	    totalEvents: number;
	    eventDistribution: StrategyEventDistribution[];
	    issueGroups: StrategyIssueAggregate[];
	    hints: StrategyAnalysisHint[];
	
	    static createFrom(source: any = {}) {
	        return new RoadStrategyAnalysis(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.template = this.convertValues(source["template"], StrategyTemplate);
	        this.sessionIds = source["sessionIds"];
	        this.sessionCount = source["sessionCount"];
	        this.totalEvents = source["totalEvents"];
	        this.eventDistribution = this.convertValues(source["eventDistribution"], StrategyEventDistribution);
	        this.issueGroups = this.convertValues(source["issueGroups"], StrategyIssueAggregate);
	        this.hints = this.convertValues(source["hints"], StrategyAnalysisHint);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class RoadTuningKnowledgeStatus {
	    loadedAt: string;
	    sourcePath: string;
	    lastError: string;
	    symptomCount: number;
	    actionCount: number;
	    usingFallback: boolean;
	
	    static createFrom(source: any = {}) {
	        return new RoadTuningKnowledgeStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.loadedAt = source["loadedAt"];
	        this.sourcePath = source["sourcePath"];
	        this.lastError = source["lastError"];
	        this.symptomCount = source["symptomCount"];
	        this.actionCount = source["actionCount"];
	        this.usingFallback = source["usingFallback"];
	    }
	}
	export class RoadTuningDecisionAction {
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
	
	    static createFrom(source: any = {}) {
	        return new RoadTuningDecisionAction(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.role = source["role"];
	        this.family = source["family"];
	        this.source = source["source"];
	        this.confidence = source["confidence"];
	        this.trustLevel = source["trustLevel"];
	        this.adviceLayer = source["adviceLayer"];
	        this.category = source["category"];
	        this.item = source["item"];
	        this.fieldKey = source["fieldKey"];
	        this.direction = source["direction"];
	        this.amount = source["amount"];
	        this.unit = source["unit"];
	        this.reason = source["reason"];
	        this.rationale = source["rationale"];
	        this.conflictReason = source["conflictReason"];
	        this.canAutoApply = source["canAutoApply"];
	        this.blockedReason = source["blockedReason"];
	        this.evidence = source["evidence"];
	    }
	}
	export class RoadTuningDecision {
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
	    relatedIssueGroup?: SessionIssueGroup;
	    evidence: Record<string, number>;
	    actions: RoadTuningDecisionAction[];
	    retestFocus: string[];
	    knowledgeStatus: RoadTuningKnowledgeStatus;
	
	    static createFrom(source: any = {}) {
	        return new RoadTuningDecision(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.sessionId = source["sessionId"];
	        this.status = source["status"];
	        this.symptomId = source["symptomId"];
	        this.phase = source["phase"];
	        this.symptom = source["symptom"];
	        this.primaryCause = source["primaryCause"];
	        this.confidence = source["confidence"];
	        this.fitVerdict = source["fitVerdict"];
	        this.reason = source["reason"];
	        this.rollbackRecommended = source["rollbackRecommended"];
	        this.relatedIssueGroup = this.convertValues(source["relatedIssueGroup"], SessionIssueGroup);
	        this.evidence = source["evidence"];
	        this.actions = this.convertValues(source["actions"], RoadTuningDecisionAction);
	        this.retestFocus = source["retestFocus"];
	        this.knowledgeStatus = this.convertValues(source["knowledgeStatus"], RoadTuningKnowledgeStatus);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	
	export class RuleThresholdProfile {
	    id: number;
	    name: string;
	    carClass: string;
	    drivetrain: string;
	    useCase: string;
	    gameMode: string;
	    configJson: string;
	    isDefault: boolean;
	    createdAt: string;
	    updatedAt: string;
	
	    static createFrom(source: any = {}) {
	        return new RuleThresholdProfile(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.carClass = source["carClass"];
	        this.drivetrain = source["drivetrain"];
	        this.useCase = source["useCase"];
	        this.gameMode = source["gameMode"];
	        this.configJson = source["configJson"];
	        this.isDefault = source["isDefault"];
	        this.createdAt = source["createdAt"];
	        this.updatedAt = source["updatedAt"];
	    }
	}
	export class RuleThresholdProfileInput {
	    name: string;
	    carClass: string;
	    drivetrain: string;
	    useCase: string;
	    gameMode: string;
	    configJson: string;
	
	    static createFrom(source: any = {}) {
	        return new RuleThresholdProfileInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.carClass = source["carClass"];
	        this.drivetrain = source["drivetrain"];
	        this.useCase = source["useCase"];
	        this.gameMode = source["gameMode"];
	        this.configJson = source["configJson"];
	    }
	}
	export class SessionEventComparison {
	    type: string;
	    left: number;
	    right: number;
	    delta: number;
	
	    static createFrom(source: any = {}) {
	        return new SessionEventComparison(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.type = source["type"];
	        this.left = source["left"];
	        this.right = source["right"];
	        this.delta = source["delta"];
	    }
	}
	export class SessionComparison {
	    leftSession: TelemetrySession;
	    rightSession: TelemetrySession;
	    metrics: SessionComparisonMetric[];
	    eventTypes: SessionEventComparison[];
	    comparabilityWarnings: string[];
	
	    static createFrom(source: any = {}) {
	        return new SessionComparison(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.leftSession = this.convertValues(source["leftSession"], TelemetrySession);
	        this.rightSession = this.convertValues(source["rightSession"], TelemetrySession);
	        this.metrics = this.convertValues(source["metrics"], SessionComparisonMetric);
	        this.eventTypes = this.convertValues(source["eventTypes"], SessionEventComparison);
	        this.comparabilityWarnings = source["comparabilityWarnings"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	
	
	export class TuningConflict {
	    key: string;
	    keptItem: string;
	    droppedItem: string;
	    reason: string;
	
	    static createFrom(source: any = {}) {
	        return new TuningConflict(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.key = source["key"];
	        this.keptItem = source["keptItem"];
	        this.droppedItem = source["droppedItem"];
	        this.reason = source["reason"];
	    }
	}
	export class WholeCarAdjustment {
	    priority: number;
	    family: string;
	    source: string;
	    confidence: string;
	    category: string;
	    item: string;
	    direction: string;
	    amount: string;
	    reason: string;
	    evidence: Record<string, number>;
	
	    static createFrom(source: any = {}) {
	        return new WholeCarAdjustment(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.priority = source["priority"];
	        this.family = source["family"];
	        this.source = source["source"];
	        this.confidence = source["confidence"];
	        this.category = source["category"];
	        this.item = source["item"];
	        this.direction = source["direction"];
	        this.amount = source["amount"];
	        this.reason = source["reason"];
	        this.evidence = source["evidence"];
	    }
	}
	export class WholeCarTuningPlan {
	    strategy: string;
	    confidence: string;
	    summary: string;
	    actions: WholeCarAdjustment[];
	    conflicts: TuningConflict[];
	    notes: string[];
	
	    static createFrom(source: any = {}) {
	        return new WholeCarTuningPlan(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.strategy = source["strategy"];
	        this.confidence = source["confidence"];
	        this.summary = source["summary"];
	        this.actions = this.convertValues(source["actions"], WholeCarAdjustment);
	        this.conflicts = this.convertValues(source["conflicts"], TuningConflict);
	        this.notes = source["notes"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class SessionIssueSummary {
	    sessionId: number;
	    baselineSession?: TelemetrySession;
	    baselineStatus: string;
	    recentChangeFields: string[];
	    groups: SessionIssueGroup[];
	    gearPower: GearPowerDiagnostic;
	    wholeCarPlan: WholeCarTuningPlan;
	
	    static createFrom(source: any = {}) {
	        return new SessionIssueSummary(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.sessionId = source["sessionId"];
	        this.baselineSession = this.convertValues(source["baselineSession"], TelemetrySession);
	        this.baselineStatus = source["baselineStatus"];
	        this.recentChangeFields = source["recentChangeFields"];
	        this.groups = this.convertValues(source["groups"], SessionIssueGroup);
	        this.gearPower = this.convertValues(source["gearPower"], GearPowerDiagnostic);
	        this.wholeCarPlan = this.convertValues(source["wholeCarPlan"], WholeCarTuningPlan);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	
	
	
	
	
	export class TestConditions {
	    driverMode: string;
	    brakeAssist: string;
	    steeringAssist: string;
	    tractionControl: string;
	    stabilityControl: string;
	    shifting: string;
	    launchControl: string;
	
	    static createFrom(source: any = {}) {
	        return new TestConditions(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.driverMode = source["driverMode"];
	        this.brakeAssist = source["brakeAssist"];
	        this.steeringAssist = source["steeringAssist"];
	        this.tractionControl = source["tractionControl"];
	        this.stabilityControl = source["stabilityControl"];
	        this.shifting = source["shifting"];
	        this.launchControl = source["launchControl"];
	    }
	}
	export class TireAxleDiagnostic {
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
	
	    static createFrom(source: any = {}) {
	        return new TireAxleDiagnostic(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.combinedSlipAvg = source["combinedSlipAvg"];
	        this.combinedSlipMax = source["combinedSlipMax"];
	        this.combinedSlipP90 = source["combinedSlipP90"];
	        this.combinedSlipHighPct = source["combinedSlipHighPct"];
	        this.slipRatioAvg = source["slipRatioAvg"];
	        this.slipRatioMax = source["slipRatioMax"];
	        this.slipRatioP90 = source["slipRatioP90"];
	        this.slipRatioHighPct = source["slipRatioHighPct"];
	        this.slipAngleAvg = source["slipAngleAvg"];
	        this.slipAngleMax = source["slipAngleMax"];
	        this.slipAngleP90 = source["slipAngleP90"];
	        this.tireTempAvg = source["tireTempAvg"];
	        this.tireTempMax = source["tireTempMax"];
	        this.suspensionTravelAvg = source["suspensionTravelAvg"];
	        this.suspensionTravelMax = source["suspensionTravelMax"];
	        this.suspensionOffsetPctAvg = source["suspensionOffsetPctAvg"];
	        this.suspensionOffsetPctMax = source["suspensionOffsetPctMax"];
	        this.limitScore = source["limitScore"];
	        this.gripState = source["gripState"];
	    }
	}
	export class TireDataQuality {
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
	
	    static createFrom(source: any = {}) {
	        return new TireDataQuality(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.status = source["status"];
	        this.confidence = source["confidence"];
	        this.sampleCount = source["sampleCount"];
	        this.dynamicSampleCount = source["dynamicSampleCount"];
	        this.speedSignal = source["speedSignal"];
	        this.gForceSignal = source["gForceSignal"];
	        this.slipSignal = source["slipSignal"];
	        this.inputSignal = source["inputSignal"];
	        this.reasons = source["reasons"];
	        this.evidence = source["evidence"];
	    }
	}
	export class TireSnapshotSubsystem {
	    status: string;
	    summary: string;
	    confidence: string;
	    explanation: string;
	    evidence: Record<string, number>;
	
	    static createFrom(source: any = {}) {
	        return new TireSnapshotSubsystem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.status = source["status"];
	        this.summary = source["summary"];
	        this.confidence = source["confidence"];
	        this.explanation = source["explanation"];
	        this.evidence = source["evidence"];
	    }
	}
	export class TireIssueAdviceGroup {
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
	
	    static createFrom(source: any = {}) {
	        return new TireIssueAdviceGroup(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.issueGroupId = source["issueGroupId"];
	        this.issueType = source["issueType"];
	        this.phase = source["phase"];
	        this.operationTags = source["operationTags"];
	        this.limitedAxle = source["limitedAxle"];
	        this.driftSource = source["driftSource"];
	        this.primaryCause = source["primaryCause"];
	        this.shouldTune = source["shouldTune"];
	        this.priority = source["priority"];
	        this.confidence = source["confidence"];
	        this.evidence = source["evidence"];
	        this.actions = this.convertValues(source["actions"], TireIssueAdviceAction);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class TireIssueAdviceAction {
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
	
	    static createFrom(source: any = {}) {
	        return new TireIssueAdviceAction(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.issueGroupId = source["issueGroupId"];
	        this.layer = source["layer"];
	        this.category = source["category"];
	        this.scope = source["scope"];
	        this.direction = source["direction"];
	        this.relatedFields = source["relatedFields"];
	        this.rationale = source["rationale"];
	        this.verifyEvidence = source["verifyEvidence"];
	        this.confidence = source["confidence"];
	        this.missingInputs = source["missingInputs"];
	        this.conflictReason = source["conflictReason"];
	        this.tuneRecommended = source["tuneRecommended"];
	    }
	}
	export class TireIssueAdvice {
	    status: string;
	    updatedAt: string;
	    confidence: string;
	    basedOnIssueUpdatedAt: string;
	    issueGroupCount: number;
	    priorityActions: TireIssueAdviceAction[];
	    groups: TireIssueAdviceGroup[];
	    warnings: string[];
	
	    static createFrom(source: any = {}) {
	        return new TireIssueAdvice(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.status = source["status"];
	        this.updatedAt = source["updatedAt"];
	        this.confidence = source["confidence"];
	        this.basedOnIssueUpdatedAt = source["basedOnIssueUpdatedAt"];
	        this.issueGroupCount = source["issueGroupCount"];
	        this.priorityActions = this.convertValues(source["priorityActions"], TireIssueAdviceAction);
	        this.groups = this.convertValues(source["groups"], TireIssueAdviceGroup);
	        this.warnings = source["warnings"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class TireIssueGroup {
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
	
	    static createFrom(source: any = {}) {
	        return new TireIssueGroup(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.type = source["type"];
	        this.phase = source["phase"];
	        this.operationTags = source["operationTags"];
	        this.driftSource = source["driftSource"];
	        this.limitType = source["limitType"];
	        this.limitedAxle = source["limitedAxle"];
	        this.limitedWheels = source["limitedWheels"];
	        this.count = source["count"];
	        this.totalDurationMs = source["totalDurationMs"];
	        this.speedMinKmh = source["speedMinKmh"];
	        this.speedMaxKmh = source["speedMaxKmh"];
	        this.speedAvgKmh = source["speedAvgKmh"];
	        this.confidence = source["confidence"];
	        this.dataQuality = source["dataQuality"];
	        this.riskLevel = source["riskLevel"];
	        this.representativeEvidence = source["representativeEvidence"];
	        this.segmentIds = source["segmentIds"];
	        this.reason = source["reason"];
	    }
	}
	export class TireIssueSegment {
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
	
	    static createFrom(source: any = {}) {
	        return new TireIssueSegment(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.type = source["type"];
	        this.phase = source["phase"];
	        this.operationTags = source["operationTags"];
	        this.driftSource = source["driftSource"];
	        this.limitType = source["limitType"];
	        this.limitedAxle = source["limitedAxle"];
	        this.limitedWheels = source["limitedWheels"];
	        this.startMs = source["startMs"];
	        this.endMs = source["endMs"];
	        this.durationMs = source["durationMs"];
	        this.sampleCount = source["sampleCount"];
	        this.speedMinKmh = source["speedMinKmh"];
	        this.speedMaxKmh = source["speedMaxKmh"];
	        this.speedAvgKmh = source["speedAvgKmh"];
	        this.confidence = source["confidence"];
	        this.dataQuality = source["dataQuality"];
	        this.riskLevel = source["riskLevel"];
	        this.evidence = source["evidence"];
	        this.reason = source["reason"];
	    }
	}
	export class TireIssueAnalysis {
	    status: string;
	    updatedAt: string;
	    windowMs: number;
	    sampleCount: number;
	    segmentCount: number;
	    groupCount: number;
	    segments: TireIssueSegment[];
	    groups: TireIssueGroup[];
	    warnings: string[];
	
	    static createFrom(source: any = {}) {
	        return new TireIssueAnalysis(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.status = source["status"];
	        this.updatedAt = source["updatedAt"];
	        this.windowMs = source["windowMs"];
	        this.sampleCount = source["sampleCount"];
	        this.segmentCount = source["segmentCount"];
	        this.groupCount = source["groupCount"];
	        this.segments = this.convertValues(source["segments"], TireIssueSegment);
	        this.groups = this.convertValues(source["groups"], TireIssueGroup);
	        this.warnings = source["warnings"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class TireSnapshotGripLimit {
	    type: string;
	    limitedAxle: string;
	    limitedWheels: string[];
	    confidence: string;
	    primaryEvidence: string;
	    reason: string;
	
	    static createFrom(source: any = {}) {
	        return new TireSnapshotGripLimit(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.type = source["type"];
	        this.limitedAxle = source["limitedAxle"];
	        this.limitedWheels = source["limitedWheels"];
	        this.confidence = source["confidence"];
	        this.primaryEvidence = source["primaryEvidence"];
	        this.reason = source["reason"];
	    }
	}
	export class TireSnapshotPhase {
	    current: string;
	    stable: string;
	    secondary: string;
	    stability: string;
	    confidence: string;
	    scoreMargin: number;
	
	    static createFrom(source: any = {}) {
	        return new TireSnapshotPhase(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.current = source["current"];
	        this.stable = source["stable"];
	        this.secondary = source["secondary"];
	        this.stability = source["stability"];
	        this.confidence = source["confidence"];
	        this.scoreMargin = source["scoreMargin"];
	    }
	}
	export class TireSnapshotQuality {
	    status: string;
	    confidence: string;
	    sampleCount: number;
	    dynamicSampleCount: number;
	    reasons: string[];
	
	    static createFrom(source: any = {}) {
	        return new TireSnapshotQuality(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.status = source["status"];
	        this.confidence = source["confidence"];
	        this.sampleCount = source["sampleCount"];
	        this.dynamicSampleCount = source["dynamicSampleCount"];
	        this.reasons = source["reasons"];
	    }
	}
	export class TireDiagnosticSnapshot {
	    generatedAt: string;
	    status: string;
	    sampleCount: number;
	    windowMs: number;
	    vehicle: SessionVehicleSnapshot;
	    dataQuality: TireSnapshotQuality;
	    phase: TireSnapshotPhase;
	    gripLimit: TireSnapshotGripLimit;
	    issueAnalysis: TireIssueAnalysis;
	    issueAdvice: TireIssueAdvice;
	    risks: string[];
	    power: TireSnapshotSubsystem;
	    brake: TireSnapshotSubsystem;
	    evidence: Record<string, number>;
	
	    static createFrom(source: any = {}) {
	        return new TireDiagnosticSnapshot(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.generatedAt = source["generatedAt"];
	        this.status = source["status"];
	        this.sampleCount = source["sampleCount"];
	        this.windowMs = source["windowMs"];
	        this.vehicle = this.convertValues(source["vehicle"], SessionVehicleSnapshot);
	        this.dataQuality = this.convertValues(source["dataQuality"], TireSnapshotQuality);
	        this.phase = this.convertValues(source["phase"], TireSnapshotPhase);
	        this.gripLimit = this.convertValues(source["gripLimit"], TireSnapshotGripLimit);
	        this.issueAnalysis = this.convertValues(source["issueAnalysis"], TireIssueAnalysis);
	        this.issueAdvice = this.convertValues(source["issueAdvice"], TireIssueAdvice);
	        this.risks = source["risks"];
	        this.power = this.convertValues(source["power"], TireSnapshotSubsystem);
	        this.brake = this.convertValues(source["brake"], TireSnapshotSubsystem);
	        this.evidence = source["evidence"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class TireGripLimit {
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
	
	    static createFrom(source: any = {}) {
	        return new TireGripLimit(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.type = source["type"];
	        this.limitedAxle = source["limitedAxle"];
	        this.limitedWheels = source["limitedWheels"];
	        this.primaryEvidence = source["primaryEvidence"];
	        this.confidence = source["confidence"];
	        this.reason = source["reason"];
	        this.frontRearDelta = source["frontRearDelta"];
	        this.drivenDelta = source["drivenDelta"];
	        this.leftRightDelta = source["leftRightDelta"];
	        this.evidence = source["evidence"];
	    }
	}
	
	
	
	
	
	
	export class TireSideBalance {
	    leftCombinedSlipAvg: number;
	    rightCombinedSlipAvg: number;
	    delta: number;
	    state: string;
	
	    static createFrom(source: any = {}) {
	        return new TireSideBalance(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.leftCombinedSlipAvg = source["leftCombinedSlipAvg"];
	        this.rightCombinedSlipAvg = source["rightCombinedSlipAvg"];
	        this.delta = source["delta"];
	        this.state = source["state"];
	    }
	}
	export class TireWheelDiagnostic {
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
	
	    static createFrom(source: any = {}) {
	        return new TireWheelDiagnostic(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.position = source["position"];
	        this.combinedSlipAvg = source["combinedSlipAvg"];
	        this.combinedSlipMax = source["combinedSlipMax"];
	        this.combinedSlipP90 = source["combinedSlipP90"];
	        this.combinedSlipHighPct = source["combinedSlipHighPct"];
	        this.slipRatioAvg = source["slipRatioAvg"];
	        this.slipRatioMax = source["slipRatioMax"];
	        this.slipRatioP90 = source["slipRatioP90"];
	        this.slipRatioHighPct = source["slipRatioHighPct"];
	        this.slipAngleAvg = source["slipAngleAvg"];
	        this.slipAngleMax = source["slipAngleMax"];
	        this.slipAngleP90 = source["slipAngleP90"];
	        this.tireTempAvg = source["tireTempAvg"];
	        this.tireTempMax = source["tireTempMax"];
	        this.suspensionTravelAvg = source["suspensionTravelAvg"];
	        this.suspensionTravelMax = source["suspensionTravelMax"];
	        this.suspensionOffsetPctAvg = source["suspensionOffsetPctAvg"];
	        this.suspensionOffsetPctMax = source["suspensionOffsetPctMax"];
	        this.suspensionTravelMetersAvg = source["suspensionTravelMetersAvg"];
	        this.suspensionTravelMetersMax = source["suspensionTravelMetersMax"];
	        this.gripState = source["gripState"];
	    }
	}
	export class TirePhaseDiagnostic {
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
	
	    static createFrom(source: any = {}) {
	        return new TirePhaseDiagnostic(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.currentPhase = source["currentPhase"];
	        this.secondaryPhase = source["secondaryPhase"];
	        this.stablePhase = source["stablePhase"];
	        this.phaseStability = source["phaseStability"];
	        this.scoreMargin = source["scoreMargin"];
	        this.confidence = source["confidence"];
	        this.scores = source["scores"];
	        this.evidence = source["evidence"];
	        this.windowMs = source["windowMs"];
	        this.sampleCount = source["sampleCount"];
	    }
	}
	export class TireModelDiagnostic {
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
	    issueAnalysis: TireIssueAnalysis;
	    issueAdvice: TireIssueAdvice;
	    hints: TireModelHint[];
	    evidence: Record<string, number>;
	    vehicle: SessionVehicleSnapshot;
	
	    static createFrom(source: any = {}) {
	        return new TireModelDiagnostic(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.status = source["status"];
	        this.updatedAt = source["updatedAt"];
	        this.sampleCount = source["sampleCount"];
	        this.windowMs = source["windowMs"];
	        this.gameMode = source["gameMode"];
	        this.phase = source["phase"];
	        this.phaseDetail = this.convertValues(source["phaseDetail"], TirePhaseDiagnostic);
	        this.dataQuality = this.convertValues(source["dataQuality"], TireDataQuality);
	        this.gripLimit = this.convertValues(source["gripLimit"], TireGripLimit);
	        this.limitType = source["limitType"];
	        this.confidence = source["confidence"];
	        this.summary = source["summary"];
	        this.explanation = source["explanation"];
	        this.warnings = source["warnings"];
	        this.wheels = this.convertValues(source["wheels"], TireWheelDiagnostic);
	        this.frontAxle = this.convertValues(source["frontAxle"], TireAxleDiagnostic);
	        this.rearAxle = this.convertValues(source["rearAxle"], TireAxleDiagnostic);
	        this.leftRight = this.convertValues(source["leftRight"], TireSideBalance);
	        this.gForce = this.convertValues(source["gForce"], GForceDiagnostic);
	        this.camber = this.convertValues(source["camber"], CamberInference);
	        this.powerToTire = this.convertValues(source["powerToTire"], PowerToTireDiagnostic);
	        this.brakeToTire = this.convertValues(source["brakeToTire"], BrakeToTireDiagnostic);
	        this.issueAnalysis = this.convertValues(source["issueAnalysis"], TireIssueAnalysis);
	        this.issueAdvice = this.convertValues(source["issueAdvice"], TireIssueAdvice);
	        this.hints = this.convertValues(source["hints"], TireModelHint);
	        this.evidence = source["evidence"];
	        this.vehicle = this.convertValues(source["vehicle"], SessionVehicleSnapshot);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	
	export class TireRegressionExpectation {
	    allowedPhases: string[];
	    requiredGripTypes: string[];
	    allowedAxles: string[];
	    forbiddenGripTypes: string[];
	    minDataQuality: string;
	    notes: string;
	
	    static createFrom(source: any = {}) {
	        return new TireRegressionExpectation(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.allowedPhases = source["allowedPhases"];
	        this.requiredGripTypes = source["requiredGripTypes"];
	        this.allowedAxles = source["allowedAxles"];
	        this.forbiddenGripTypes = source["forbiddenGripTypes"];
	        this.minDataQuality = source["minDataQuality"];
	        this.notes = source["notes"];
	    }
	}
	export class TireRegressionResult {
	    sampleId: string;
	    name: string;
	    scenario: string;
	    passed: boolean;
	    status: string;
	    failures: string[];
	    expected: TireRegressionExpectation;
	    actual: TireDiagnosticSnapshot;
	
	    static createFrom(source: any = {}) {
	        return new TireRegressionResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.sampleId = source["sampleId"];
	        this.name = source["name"];
	        this.scenario = source["scenario"];
	        this.passed = source["passed"];
	        this.status = source["status"];
	        this.failures = source["failures"];
	        this.expected = this.convertValues(source["expected"], TireRegressionExpectation);
	        this.actual = this.convertValues(source["actual"], TireDiagnosticSnapshot);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class TireRegressionSample {
	    id: string;
	    name: string;
	    scenario: string;
	    createdAt: string;
	    windowSeconds: number;
	    vehicle: SessionVehicleSnapshot;
	    sampleCount: number;
	    samples: telemetry.NormalizedTelemetry[];
	    snapshot: TireDiagnosticSnapshot;
	    expected: TireRegressionExpectation;
	
	    static createFrom(source: any = {}) {
	        return new TireRegressionSample(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.scenario = source["scenario"];
	        this.createdAt = source["createdAt"];
	        this.windowSeconds = source["windowSeconds"];
	        this.vehicle = this.convertValues(source["vehicle"], SessionVehicleSnapshot);
	        this.sampleCount = source["sampleCount"];
	        this.samples = this.convertValues(source["samples"], telemetry.NormalizedTelemetry);
	        this.snapshot = this.convertValues(source["snapshot"], TireDiagnosticSnapshot);
	        this.expected = this.convertValues(source["expected"], TireRegressionExpectation);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class TireRegressionSampleInput {
	    name: string;
	    scenario: string;
	    windowSeconds: number;
	    expected: TireRegressionExpectation;
	
	    static createFrom(source: any = {}) {
	        return new TireRegressionSampleInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.scenario = source["scenario"];
	        this.windowSeconds = source["windowSeconds"];
	        this.expected = this.convertValues(source["expected"], TireRegressionExpectation);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class TireRegressionSampleSummary {
	    id: string;
	    name: string;
	    scenario: string;
	    createdAt: string;
	    windowSeconds: number;
	    vehicle: SessionVehicleSnapshot;
	    sampleCount: number;
	    expected: TireRegressionExpectation;
	
	    static createFrom(source: any = {}) {
	        return new TireRegressionSampleSummary(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.scenario = source["scenario"];
	        this.createdAt = source["createdAt"];
	        this.windowSeconds = source["windowSeconds"];
	        this.vehicle = this.convertValues(source["vehicle"], SessionVehicleSnapshot);
	        this.sampleCount = source["sampleCount"];
	        this.expected = this.convertValues(source["expected"], TireRegressionExpectation);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	
	
	
	
	
	export class TrackRunContext {
	    run: BenchmarkRun;
	    session: TelemetrySession;
	    vehicle: TrackVehicleKey;
	
	    static createFrom(source: any = {}) {
	        return new TrackRunContext(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.run = this.convertValues(source["run"], BenchmarkRun);
	        this.session = this.convertValues(source["session"], TelemetrySession);
	        this.vehicle = this.convertValues(source["vehicle"], TrackVehicleKey);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class TrackVehicleKey {
	    carOrdinal?: number;
	    carClass: string;
	    carPi?: number;
	    drivetrain: string;
	    label: string;
	
	    static createFrom(source: any = {}) {
	        return new TrackVehicleKey(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.carOrdinal = source["carOrdinal"];
	        this.carClass = source["carClass"];
	        this.carPi = source["carPi"];
	        this.drivetrain = source["drivetrain"];
	        this.label = source["label"];
	    }
	}
	export class TrackAutoBaseline {
	    vehicle: TrackVehicleKey;
	    bestRun: TrackRunContext;
	    recentRuns: TrackRunContext[];
	    runCount: number;
	
	    static createFrom(source: any = {}) {
	        return new TrackAutoBaseline(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.vehicle = this.convertValues(source["vehicle"], TrackVehicleKey);
	        this.bestRun = this.convertValues(source["bestRun"], TrackRunContext);
	        this.recentRuns = this.convertValues(source["recentRuns"], TrackRunContext);
	        this.runCount = source["runCount"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class TrackBaselineRun {
	    id: number;
	    trackId: number;
	    vehicle: TrackVehicleKey;
	    startMs: number;
	    endMs: number;
	    durationMs: number;
	    confidence: number;
	    avgSpeedKmh?: number;
	    maxSpeedKmh?: number;
	    routeProgress01?: number;
	    geometryLengthMeters?: number;
	    trackLengthErrorPct?: number;
	    distanceTraveledDeltaMeters?: number;
	    currentRaceTimeDeltaSeconds?: number;
	    avgLateralErrorMeters?: number;
	    maxLateralErrorMeters?: number;
	    warningFlags: string;
	    eventCount: number;
	    driverMode: string;
	    driverModeConfidence: number;
	    driverModeEvidenceJson: string;
	    valid: boolean;
	    gameMode: string;
	    createdAt: string;
	
	    static createFrom(source: any = {}) {
	        return new TrackBaselineRun(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.trackId = source["trackId"];
	        this.vehicle = this.convertValues(source["vehicle"], TrackVehicleKey);
	        this.startMs = source["startMs"];
	        this.endMs = source["endMs"];
	        this.durationMs = source["durationMs"];
	        this.confidence = source["confidence"];
	        this.avgSpeedKmh = source["avgSpeedKmh"];
	        this.maxSpeedKmh = source["maxSpeedKmh"];
	        this.routeProgress01 = source["routeProgress01"];
	        this.geometryLengthMeters = source["geometryLengthMeters"];
	        this.trackLengthErrorPct = source["trackLengthErrorPct"];
	        this.distanceTraveledDeltaMeters = source["distanceTraveledDeltaMeters"];
	        this.currentRaceTimeDeltaSeconds = source["currentRaceTimeDeltaSeconds"];
	        this.avgLateralErrorMeters = source["avgLateralErrorMeters"];
	        this.maxLateralErrorMeters = source["maxLateralErrorMeters"];
	        this.warningFlags = source["warningFlags"];
	        this.eventCount = source["eventCount"];
	        this.driverMode = source["driverMode"];
	        this.driverModeConfidence = source["driverModeConfidence"];
	        this.driverModeEvidenceJson = source["driverModeEvidenceJson"];
	        this.valid = source["valid"];
	        this.gameMode = source["gameMode"];
	        this.createdAt = source["createdAt"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class TrackMergeCandidate {
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
	
	    static createFrom(source: any = {}) {
	        return new TrackMergeCandidate(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.track = this.convertValues(source["track"], BenchmarkTrack);
	        this.matchLevel = source["matchLevel"];
	        this.lengthErrorPct = source["lengthErrorPct"];
	        this.startDistanceMeters = source["startDistanceMeters"];
	        this.endDistanceMeters = source["endDistanceMeters"];
	        this.shapeSimilarity = source["shapeSimilarity"];
	        this.routeFitAvgErrorMeters = source["routeFitAvgErrorMeters"];
	        this.routeFitP90ErrorMeters = source["routeFitP90ErrorMeters"];
	        this.routeFitScore = source["routeFitScore"];
	        this.directionMatched = source["directionMatched"];
	        this.reverseMatched = source["reverseMatched"];
	        this.reason = source["reason"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class TrackBaselineSaveResult {
	    track: BenchmarkTrack;
	    baseline: TrackBaselineRun;
	    action: string;
	    matchCandidate?: TrackMergeCandidate;
	
	    static createFrom(source: any = {}) {
	        return new TrackBaselineSaveResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.track = this.convertValues(source["track"], BenchmarkTrack);
	        this.baseline = this.convertValues(source["baseline"], TrackBaselineRun);
	        this.action = source["action"];
	        this.matchCandidate = this.convertValues(source["matchCandidate"], TrackMergeCandidate);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class TrackVehicleReference {
	    vehicle: TrackVehicleKey;
	    bestAutoBaseline?: TrackRunContext;
	    bestTrackBaseline?: TrackBaselineRun;
	    recentRuns: TrackRunContext[];
	    recentBaselineRuns: TrackBaselineRun[];
	    validRunCount: number;
	    autoRunCount: number;
	    baselineRunCount: number;
	    avgSpeedKmh?: number;
	    maxSpeedKmh?: number;
	    eventCount: number;
	
	    static createFrom(source: any = {}) {
	        return new TrackVehicleReference(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.vehicle = this.convertValues(source["vehicle"], TrackVehicleKey);
	        this.bestAutoBaseline = this.convertValues(source["bestAutoBaseline"], TrackRunContext);
	        this.bestTrackBaseline = this.convertValues(source["bestTrackBaseline"], TrackBaselineRun);
	        this.recentRuns = this.convertValues(source["recentRuns"], TrackRunContext);
	        this.recentBaselineRuns = this.convertValues(source["recentBaselineRuns"], TrackBaselineRun);
	        this.validRunCount = source["validRunCount"];
	        this.autoRunCount = source["autoRunCount"];
	        this.baselineRunCount = source["baselineRunCount"];
	        this.avgSpeedKmh = source["avgSpeedKmh"];
	        this.maxSpeedKmh = source["maxSpeedKmh"];
	        this.eventCount = source["eventCount"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class TrackProfile {
	    track: BenchmarkTrack;
	    autoBaselines: TrackAutoBaseline[];
	    vehicleReferences: TrackVehicleReference[];
	    recentRuns: TrackRunContext[];
	    warnings: string[];
	
	    static createFrom(source: any = {}) {
	        return new TrackProfile(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.track = this.convertValues(source["track"], BenchmarkTrack);
	        this.autoBaselines = this.convertValues(source["autoBaselines"], TrackAutoBaseline);
	        this.vehicleReferences = this.convertValues(source["vehicleReferences"], TrackVehicleReference);
	        this.recentRuns = this.convertValues(source["recentRuns"], TrackRunContext);
	        this.warnings = source["warnings"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	
	
	export class TuneAdjustmentExplanation {
	    category: string;
	    item: string;
	    detail: string;
	    description: string;
	
	    static createFrom(source: any = {}) {
	        return new TuneAdjustmentExplanation(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.category = source["category"];
	        this.item = source["item"];
	        this.detail = source["detail"];
	        this.description = source["description"];
	    }
	}
	export class TuneFieldInfluence {
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
	
	    static createFrom(source: any = {}) {
	        return new TuneFieldInfluence(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.fieldKey = source["fieldKey"];
	        this.category = source["category"];
	        this.labelZh = source["labelZh"];
	        this.labelEn = source["labelEn"];
	        this.influenceType = source["influenceType"];
	        this.scope = source["scope"];
	        this.phases = source["phases"];
	        this.tireMetrics = source["tireMetrics"];
	        this.evidenceKeys = source["evidenceKeys"];
	        this.sideEffects = source["sideEffects"];
	        this.conditions = source["conditions"];
	        this.summaryZh = source["summaryZh"];
	        this.summaryEn = source["summaryEn"];
	    }
	}
	export class TuneHarvestCandidate {
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
	
	    static createFrom(source: any = {}) {
	        return new TuneHarvestCandidate(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.runId = source["runId"];
	        this.source = source["source"];
	        this.sourceRef = source["sourceRef"];
	        this.sourceUrl = source["sourceUrl"];
	        this.sourceCarId = source["sourceCarId"];
	        this.rawKey = source["rawKey"];
	        this.shareCode = source["shareCode"];
	        this.year = source["year"];
	        this.make = source["make"];
	        this.model = source["model"];
	        this.carName = source["carName"];
	        this.matchedCarId = source["matchedCarId"];
	        this.matchScore = source["matchScore"];
	        this.matchReason = source["matchReason"];
	        this.useCase = source["useCase"];
	        this.carClass = source["carClass"];
	        this.pi = source["pi"];
	        this.drivetrain = source["drivetrain"];
	        this.tireCompound = source["tireCompound"];
	        this.tuner = source["tuner"];
	        this.tuneName = source["tuneName"];
	        this.bestFor = source["bestFor"];
	        this.difficulty = source["difficulty"];
	        this.notes = source["notes"];
	        this.rawJson = source["rawJson"];
	        this.status = source["status"];
	        this.rejectionReason = source["rejectionReason"];
	        this.createdAt = source["createdAt"];
	        this.updatedAt = source["updatedAt"];
	    }
	}
	export class TuneHarvestOptions {
	    sources: string[];
	    dryRun: boolean;
	    limitPerSource: number;
	
	    static createFrom(source: any = {}) {
	        return new TuneHarvestOptions(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.sources = source["sources"];
	        this.dryRun = source["dryRun"];
	        this.limitPerSource = source["limitPerSource"];
	    }
	}
	export class TuneHarvestRun {
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
	
	    static createFrom(source: any = {}) {
	        return new TuneHarvestRun(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.startedAt = source["startedAt"];
	        this.finishedAt = source["finishedAt"];
	        this.sources = source["sources"];
	        this.dryRun = source["dryRun"];
	        this.status = source["status"];
	        this.message = source["message"];
	        this.foundCount = source["foundCount"];
	        this.savedCount = source["savedCount"];
	        this.rejectedCount = source["rejectedCount"];
	        this.pendingCount = source["pendingCount"];
	        this.importedCount = source["importedCount"];
	    }
	}
	export class TuneHarvestRunResult {
	    run?: TuneHarvestRun;
	    candidates: TuneHarvestCandidate[];
	    found: number;
	    saved: number;
	    rejected: number;
	    pending: number;
	    imported: number;
	    warnings: string[];
	
	    static createFrom(source: any = {}) {
	        return new TuneHarvestRunResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.run = this.convertValues(source["run"], TuneHarvestRun);
	        this.candidates = this.convertValues(source["candidates"], TuneHarvestCandidate);
	        this.found = source["found"];
	        this.saved = source["saved"];
	        this.rejected = source["rejected"];
	        this.pending = source["pending"];
	        this.imported = source["imported"];
	        this.warnings = source["warnings"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class TunePlanApplyInput {
	    sessionId: number;
	    selectedActionIds: string[];
	
	    static createFrom(source: any = {}) {
	        return new TunePlanApplyInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.sessionId = source["sessionId"];
	        this.selectedActionIds = source["selectedActionIds"];
	    }
	}
	export class TunePlanApplyResult {
	    profile: TuneProfile;
	    appliedActions: TunePlanDraftAction[];
	    changedFields: string[];
	
	    static createFrom(source: any = {}) {
	        return new TunePlanApplyResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.profile = this.convertValues(source["profile"], TuneProfile);
	        this.appliedActions = this.convertValues(source["appliedActions"], TunePlanDraftAction);
	        this.changedFields = source["changedFields"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class TunePlanDraft {
	    sessionId: number;
	    tuneProfileId?: number;
	    status: string;
	    summary: string;
	    actions: TunePlanDraftAction[];
	    conflicts: TuningConflict[];
	
	    static createFrom(source: any = {}) {
	        return new TunePlanDraft(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.sessionId = source["sessionId"];
	        this.tuneProfileId = source["tuneProfileId"];
	        this.status = source["status"];
	        this.summary = source["summary"];
	        this.actions = this.convertValues(source["actions"], TunePlanDraftAction);
	        this.conflicts = this.convertValues(source["conflicts"], TuningConflict);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	
	
	export class TuneProfileSessionStat {
	    tuneProfileId: number;
	    sessionCount: number;
	    lastStartedAt: string;
	
	    static createFrom(source: any = {}) {
	        return new TuneProfileSessionStat(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.tuneProfileId = source["tuneProfileId"];
	        this.sessionCount = source["sessionCount"];
	        this.lastStartedAt = source["lastStartedAt"];
	    }
	}
	export class TuneProfileSnapshot {
	    id: number;
	    tuneProfileId: number;
	    sessionId?: number;
	    changedAt: string;
	    changeReason: string;
	    before: TuneProfile;
	    after: TuneProfile;
	    changedFields: string[];
	    changeJson: string;
	
	    static createFrom(source: any = {}) {
	        return new TuneProfileSnapshot(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.tuneProfileId = source["tuneProfileId"];
	        this.sessionId = source["sessionId"];
	        this.changedAt = source["changedAt"];
	        this.changeReason = source["changeReason"];
	        this.before = this.convertValues(source["before"], TuneProfile);
	        this.after = this.convertValues(source["after"], TuneProfile);
	        this.changedFields = source["changedFields"];
	        this.changeJson = source["changeJson"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class TuneToTireInfluenceMap {
	    version: string;
	    items: TuneFieldInfluence[];
	
	    static createFrom(source: any = {}) {
	        return new TuneToTireInfluenceMap(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.version = source["version"];
	        this.items = this.convertValues(source["items"], TuneFieldInfluence);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	
	
	
	
	export class TuningPipelineCombination {
	    sourceType: string;
	    detectorId: string;
	    decisionerId: string;
	    interpreterId: string;
	    description: string;
	
	    static createFrom(source: any = {}) {
	        return new TuningPipelineCombination(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.sourceType = source["sourceType"];
	        this.detectorId = source["detectorId"];
	        this.decisionerId = source["decisionerId"];
	        this.interpreterId = source["interpreterId"];
	        this.description = source["description"];
	    }
	}
	export class TuningPipelineComponent {
	    id: string;
	    name: string;
	    description: string;
	    sourceTypes: string[];
	    compatibleWith: string[];
	    tags: string[];
	
	    static createFrom(source: any = {}) {
	        return new TuningPipelineComponent(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.sourceTypes = source["sourceTypes"];
	        this.compatibleWith = source["compatibleWith"];
	        this.tags = source["tags"];
	    }
	}
	export class TuningPipelineCatalog {
	    sourceTypes: TuningPipelineComponent[];
	    detectors: TuningPipelineComponent[];
	    decisioners: TuningPipelineComponent[];
	    interpreters: TuningPipelineComponent[];
	    defaultCombinations: TuningPipelineCombination[];
	
	    static createFrom(source: any = {}) {
	        return new TuningPipelineCatalog(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.sourceTypes = this.convertValues(source["sourceTypes"], TuningPipelineComponent);
	        this.detectors = this.convertValues(source["detectors"], TuningPipelineComponent);
	        this.decisioners = this.convertValues(source["decisioners"], TuningPipelineComponent);
	        this.interpreters = this.convertValues(source["interpreters"], TuningPipelineComponent);
	        this.defaultCombinations = this.convertValues(source["defaultCombinations"], TuningPipelineCombination);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	
	export class TuningPipelineRunInput {
	    sourceType: string;
	    sessionId?: number;
	    detectorId: string;
	    decisionerId: string;
	    interpreterId: string;
	
	    static createFrom(source: any = {}) {
	        return new TuningPipelineRunInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.sourceType = source["sourceType"];
	        this.sessionId = source["sessionId"];
	        this.detectorId = source["detectorId"];
	        this.decisionerId = source["decisionerId"];
	        this.interpreterId = source["interpreterId"];
	    }
	}
	
	
	
	
	export class UpgradeUnlockRule {
	    category: string;
	    upgradeName: string;
	    unlocks: string;
	
	    static createFrom(source: any = {}) {
	        return new UpgradeUnlockRule(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.category = source["category"];
	        this.upgradeName = source["upgradeName"];
	        this.unlocks = source["unlocks"];
	    }
	}
	

}

export namespace telemetry {
	
	export class SuggestedAction {
	    priority: number;
	    category: string;
	    item: string;
	    direction: string;
	    amount: string;
	    reason: string;
	
	    static createFrom(source: any = {}) {
	        return new SuggestedAction(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.priority = source["priority"];
	        this.category = source["category"];
	        this.item = source["item"];
	        this.direction = source["direction"];
	        this.amount = source["amount"];
	        this.reason = source["reason"];
	    }
	}
	export class DetectedEvent {
	    id: string;
	    type: string;
	    severity: string;
	    startMs: number;
	    endMs: number;
	    durationMs: number;
	    segment: string;
	    evidence: Record<string, number>;
	    suggestedActions: SuggestedAction[];
	
	    static createFrom(source: any = {}) {
	        return new DetectedEvent(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.type = source["type"];
	        this.severity = source["severity"];
	        this.startMs = source["startMs"];
	        this.endMs = source["endMs"];
	        this.durationMs = source["durationMs"];
	        this.segment = source["segment"];
	        this.evidence = source["evidence"];
	        this.suggestedActions = this.convertValues(source["suggestedActions"], SuggestedAction);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class NetworkInterface {
	    name: string;
	    displayName: string;
	    address: string;
	    isLoopback: boolean;
	    isPrivate: boolean;
	    isUp: boolean;
	
	    static createFrom(source: any = {}) {
	        return new NetworkInterface(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.displayName = source["displayName"];
	        this.address = source["address"];
	        this.isLoopback = source["isLoopback"];
	        this.isPrivate = source["isPrivate"];
	        this.isUp = source["isUp"];
	    }
	}
	export class NormalizedWheelTelemetry {
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
	
	    static createFrom(source: any = {}) {
	        return new NormalizedWheelTelemetry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.slipRatio = source["slipRatio"];
	        this.slipAngle = source["slipAngle"];
	        this.combinedSlip = source["combinedSlip"];
	        this.tireTemp = source["tireTemp"];
	        this.suspensionTravel = source["suspensionTravel"];
	        this.suspensionTravelMeters = source["suspensionTravelMeters"];
	        this.wheelRotationSpeed = source["wheelRotationSpeed"];
	        this.rumbleStrip = source["rumbleStrip"];
	        this.puddleDepth = source["puddleDepth"];
	        this.surfaceRumble = source["surfaceRumble"];
	    }
	}
	export class NormalizedTelemetry {
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
	    wheelFL: NormalizedWheelTelemetry;
	    wheelFR: NormalizedWheelTelemetry;
	    wheelRL: NormalizedWheelTelemetry;
	    wheelRR: NormalizedWheelTelemetry;
	
	    static createFrom(source: any = {}) {
	        return new NormalizedTelemetry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.receivedAt = source["receivedAt"];
	        this.timeMs = source["timeMs"];
	        this.isRaceOn = source["isRaceOn"];
	        this.gameMode = source["gameMode"];
	        this.speedKmh = source["speedKmh"];
	        this.speedFieldKmh = source["speedFieldKmh"];
	        this.velocitySpeedKmh = source["velocitySpeedKmh"];
	        this.speedSource = source["speedSource"];
	        this.rpm = source["rpm"];
	        this.rpmRatio = source["rpmRatio"];
	        this.engineMaxRpm = source["engineMaxRpm"];
	        this.engineIdleRpm = source["engineIdleRpm"];
	        this.gear = source["gear"];
	        this.accelerationX = source["accelerationX"];
	        this.accelerationY = source["accelerationY"];
	        this.accelerationZ = source["accelerationZ"];
	        this.velocityX = source["velocityX"];
	        this.velocityY = source["velocityY"];
	        this.velocityZ = source["velocityZ"];
	        this.yaw = source["yaw"];
	        this.pitch = source["pitch"];
	        this.roll = source["roll"];
	        this.power = source["power"];
	        this.torque = source["torque"];
	        this.positionX = source["positionX"];
	        this.positionY = source["positionY"];
	        this.positionZ = source["positionZ"];
	        this.boost = source["boost"];
	        this.fuel = source["fuel"];
	        this.distanceTraveled = source["distanceTraveled"];
	        this.bestLap = source["bestLap"];
	        this.lastLap = source["lastLap"];
	        this.currentLap = source["currentLap"];
	        this.currentRaceTime = source["currentRaceTime"];
	        this.lapNumber = source["lapNumber"];
	        this.racePosition = source["racePosition"];
	        this.smashableVelDiff = source["smashableVelDiff"];
	        this.smashableMass = source["smashableMass"];
	        this.carOrdinal = source["carOrdinal"];
	        this.carClassId = source["carClassId"];
	        this.carClass = source["carClass"];
	        this.carPi = source["carPi"];
	        this.drivetrainType = source["drivetrainType"];
	        this.drivetrain = source["drivetrain"];
	        this.numCylinders = source["numCylinders"];
	        this.carCategory = source["carCategory"];
	        this.carCategoryName = source["carCategoryName"];
	        this.throttle01 = source["throttle01"];
	        this.brake01 = source["brake01"];
	        this.clutch01 = source["clutch01"];
	        this.handBrake01 = source["handBrake01"];
	        this.steer01 = source["steer01"];
	        this.drivingLine01 = source["drivingLine01"];
	        this.aiBrakeDifference01 = source["aiBrakeDifference01"];
	        this.frontSlipRatioAvg = source["frontSlipRatioAvg"];
	        this.rearSlipRatioAvg = source["rearSlipRatioAvg"];
	        this.frontSlipAngleAvg = source["frontSlipAngleAvg"];
	        this.rearSlipAngleAvg = source["rearSlipAngleAvg"];
	        this.frontCombinedSlipAvg = source["frontCombinedSlipAvg"];
	        this.rearCombinedSlipAvg = source["rearCombinedSlipAvg"];
	        this.tireTempFrontAvg = source["tireTempFrontAvg"];
	        this.tireTempRearAvg = source["tireTempRearAvg"];
	        this.suspensionFrontAvg = source["suspensionFrontAvg"];
	        this.suspensionRearAvg = source["suspensionRearAvg"];
	        this.yawRate = source["yawRate"];
	        this.pitchRate = source["pitchRate"];
	        this.rollRate = source["rollRate"];
	        this.wheelFL = this.convertValues(source["wheelFL"], NormalizedWheelTelemetry);
	        this.wheelFR = this.convertValues(source["wheelFR"], NormalizedWheelTelemetry);
	        this.wheelRL = this.convertValues(source["wheelRL"], NormalizedWheelTelemetry);
	        this.wheelRR = this.convertValues(source["wheelRR"], NormalizedWheelTelemetry);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	
	export class TelemetryReplayStatus {
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
	
	    static createFrom(source: any = {}) {
	        return new TelemetryReplayStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.running = source["running"];
	        this.paused = source["paused"];
	        this.sessionId = source["sessionId"];
	        this.speed = source["speed"];
	        this.positionMs = source["positionMs"];
	        this.durationMs = source["durationMs"];
	        this.progress01 = source["progress01"];
	        this.packetIndex = source["packetIndex"];
	        this.packetCount = source["packetCount"];
	        this.lastError = source["lastError"];
	    }
	}
	export class TelemetryStatus {
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
	
	    static createFrom(source: any = {}) {
	        return new TelemetryStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.running = source["running"];
	        this.mode = source["mode"];
	        this.analysisMode = source["analysisMode"];
	        this.address = source["address"];
	        this.port = source["port"];
	        this.packetLength = source["packetLength"];
	        this.rawPackets = source["rawPackets"];
	        this.validPackets = source["validPackets"];
	        this.invalidPackets = source["invalidPackets"];
	        this.parseErrors = source["parseErrors"];
	        this.lastDatagramAt = source["lastDatagramAt"];
	        this.lastDatagramBytes = source["lastDatagramBytes"];
	        this.lastDatagramRemote = source["lastDatagramRemote"];
	        this.lastPacketAt = source["lastPacketAt"];
	        this.lastError = source["lastError"];
	        this.hasCurrentFrame = source["hasCurrentFrame"];
	        this.recordingActive = source["recordingActive"];
	        this.recordingBytes = source["recordingBytes"];
	        this.recordingLimitBytes = source["recordingLimitBytes"];
	        this.recordingPackets = source["recordingPackets"];
	        this.recordingTruncated = source["recordingTruncated"];
	    }
	}

}

