package storage

import (
	"math"
	"strings"
	"testing"
)

func TestGenerateRoadStaticTuneBaselineMinimalInput(t *testing.T) {
	result, err := GenerateRoadStaticTuneBaseline(RoadStaticTuneBaselineInput{
		PI:             700,
		Drivetrain:     "AWD",
		WeightKG:       1400,
		FrontWeightPct: 54,
	})
	if err != nil {
		t.Fatalf("generate baseline: %v", err)
	}
	if result.ProfileDraft.UseCase != "Road" || result.ProfileDraft.CarClass != "A" {
		t.Fatalf("draft identity = %#v", result.ProfileDraft)
	}
	if result.ProfileDraft.CarName != "" {
		t.Fatalf("empty car name should be allowed in preview, got %q", result.ProfileDraft.CarName)
	}
	if result.ProfileDraft.FrontTirePressure == nil || *result.ProfileDraft.FrontTirePressure != 2.17 {
		t.Fatalf("front tire pressure = %#v", result.ProfileDraft.FrontTirePressure)
	}
	if result.ProfileDraft.Gear1 != nil || result.ProfileDraft.FinalDrive != nil {
		t.Fatalf("gearing should not be generated: final=%v gear1=%v", result.ProfileDraft.FinalDrive, result.ProfileDraft.Gear1)
	}
	if result.ProfileDraft.FrontRideHeight != nil || result.ProfileDraft.RearRideHeight != nil || result.ProfileDraft.FrontAero != nil || result.ProfileDraft.RearAero != nil {
		t.Fatalf("ride height and aero should not be precise generated fields: %#v", result.ProfileDraft)
	}
	if result.ProfileDraft.FrontDiffAccel == nil || result.ProfileDraft.RearDiffAccel == nil || result.ProfileDraft.CenterDiffBalance == nil {
		t.Fatalf("AWD diff baseline incomplete: %#v", result.ProfileDraft)
	}
	if len(result.GeneratedFields) == 0 || !hasGeneratedField(result.GeneratedFields, "frontSpring") || hasGeneratedField(result.GeneratedFields, "frontRideHeight") || hasGeneratedField(result.GeneratedFields, "frontAero") || !hasSkippedField(result.SkippedFields, "gear10") {
		t.Fatalf("generated=%#v skipped=%#v", result.GeneratedFields, result.SkippedFields)
	}
}

func TestGenerateRoadStaticTuneBaselineRejectsPIOutOfRange(t *testing.T) {
	_, err := GenerateRoadStaticTuneBaseline(RoadStaticTuneBaselineInput{
		PI:             99,
		Drivetrain:     "AWD",
		WeightKG:       1400,
		FrontWeightPct: 54,
	})
	if err == nil || !strings.Contains(err.Error(), "PI must be between 100 and 999") {
		t.Fatalf("expected PI range error, got %v", err)
	}
}

func TestGenerateRoadStaticTuneBaselineRejectsInvalidTireCompound(t *testing.T) {
	_, err := GenerateRoadStaticTuneBaseline(RoadStaticTuneBaselineInput{
		PI:             700,
		Drivetrain:     "AWD",
		TireCompound:   "magic",
		WeightKG:       1400,
		FrontWeightPct: 54,
	})
	if err == nil || !strings.Contains(err.Error(), "tire compound must be one of") {
		t.Fatalf("expected tire compound error, got %v", err)
	}
}

func TestGenerateRoadStaticTuneBaselineRejectsInvalidVehicleInputs(t *testing.T) {
	cases := []struct {
		name        string
		input       RoadStaticTuneBaselineInput
		wantMessage string
	}{
		{
			name: "weight below range",
			input: RoadStaticTuneBaselineInput{
				PI:             700,
				Drivetrain:     "AWD",
				WeightKG:       299,
				FrontWeightPct: 54,
			},
			wantMessage: "weight kg must be an integer between 300 and 3000",
		},
		{
			name: "weight decimal",
			input: RoadStaticTuneBaselineInput{
				PI:             700,
				Drivetrain:     "AWD",
				WeightKG:       1400.5,
				FrontWeightPct: 54,
			},
			wantMessage: "weight kg must be an integer between 300 and 3000",
		},
		{
			name: "front weight decimal",
			input: RoadStaticTuneBaselineInput{
				PI:             700,
				Drivetrain:     "AWD",
				WeightKG:       1400,
				FrontWeightPct: 54.5,
			},
			wantMessage: "front weight percentage must be an integer between 1 and 99",
		},
		{
			name: "front weight out of range",
			input: RoadStaticTuneBaselineInput{
				PI:             700,
				Drivetrain:     "AWD",
				WeightKG:       1400,
				FrontWeightPct: 100,
			},
			wantMessage: "front weight percentage must be an integer between 1 and 99",
		},
		{
			name: "balance bias low",
			input: RoadStaticTuneBaselineInput{
				PI:             700,
				Drivetrain:     "AWD",
				WeightKG:       1400,
				FrontWeightPct: 54,
				BalanceBias:    49,
			},
			wantMessage: "balance bias must be an integer between 50 and 150",
		},
		{
			name: "stiffness bias decimal",
			input: RoadStaticTuneBaselineInput{
				PI:             700,
				Drivetrain:     "AWD",
				WeightKG:       1400,
				FrontWeightPct: 54,
				StiffnessBias:  100.5,
			},
			wantMessage: "stiffness bias must be an integer between 50 and 150",
		},
		{
			name: "speed bias high",
			input: RoadStaticTuneBaselineInput{
				PI:             700,
				Drivetrain:     "AWD",
				WeightKG:       1400,
				FrontWeightPct: 54,
				SpeedBias:      151,
			},
			wantMessage: "speed bias must be an integer between 50 and 150",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := GenerateRoadStaticTuneBaseline(tc.input)
			if err == nil || !strings.Contains(err.Error(), tc.wantMessage) {
				t.Fatalf("expected %q, got %v", tc.wantMessage, err)
			}
		})
	}
}

func TestGenerateRoadStaticTuneBaselineRejectsInvalidGearingInputs(t *testing.T) {
	redline := 7200.0
	gearCount := int64(6)
	tireDiameter := 67.0
	targetTopSpeed := 280.0
	cases := []struct {
		name        string
		mutate      func()
		wantMessage string
	}{
		{
			name: "redline decimal",
			mutate: func() {
				redline = 7200.5
			},
			wantMessage: "redline RPM must be an integer between 1000 and 20000",
		},
		{
			name: "gear count high",
			mutate: func() {
				gearCount = 11
			},
			wantMessage: "gear count must be between 2 and 10",
		},
		{
			name: "tire diameter low",
			mutate: func() {
				tireDiameter = 39.9
			},
			wantMessage: "tire diameter must be between 40 cm and 120 cm",
		},
		{
			name: "target speed decimal",
			mutate: func() {
				targetTopSpeed = 280.5
			},
			wantMessage: "target top speed must be an integer between 1 and 600",
		},
		{
			name: "target speed high",
			mutate: func() {
				targetTopSpeed = 601
			},
			wantMessage: "target top speed must be an integer between 1 and 600",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			redline = 7200
			gearCount = 6
			tireDiameter = 67
			targetTopSpeed = 280
			tc.mutate()
			_, err := GenerateRoadStaticTuneBaseline(RoadStaticTuneBaselineInput{
				PI:                800,
				Drivetrain:        "RWD",
				WeightKG:          1350,
				FrontWeightPct:    53,
				RedlineRPM:        &redline,
				GearCount:         &gearCount,
				TireDiameterCm:    &tireDiameter,
				TargetTopSpeedKmh: &targetTopSpeed,
			})
			if err == nil || !strings.Contains(err.Error(), tc.wantMessage) {
				t.Fatalf("expected %q, got %v", tc.wantMessage, err)
			}
		})
	}
}

func TestGenerateRoadStaticTuneBaselineRejectsInvalidDriftTargetSpeed(t *testing.T) {
	redline := 7200.0
	gearCount := int64(6)
	tireDiameter := 67.0
	targetDriftSpeed := 220.0
	_, err := GenerateRoadStaticTuneBaseline(RoadStaticTuneBaselineInput{
		UseCase:           "Drift",
		PI:                700,
		Drivetrain:        "RWD",
		WeightKG:          1400,
		FrontWeightPct:    54,
		RedlineRPM:        &redline,
		GearCount:         &gearCount,
		TireDiameterCm:    &tireDiameter,
		TargetTopSpeedKmh: &targetDriftSpeed,
	})
	if err == nil || !strings.Contains(err.Error(), "target drift speed must be an integer between 40 and 180") {
		t.Fatalf("expected drift target speed error, got %v", err)
	}
}

func TestGenerateRoadStaticTuneBaselineRallyAWD(t *testing.T) {
	result, err := GenerateRoadStaticTuneBaseline(RoadStaticTuneBaselineInput{
		UseCase:        "拉力",
		PI:             700,
		Drivetrain:     "AWD",
		WeightKG:       1400,
		FrontWeightPct: 54,
	})
	if err != nil {
		t.Fatalf("generate rally baseline: %v", err)
	}
	if result.ProfileDraft.UseCase != "Rally" || result.ProfileDraft.VersionName != "Rally Baseline" {
		t.Fatalf("rally identity = %#v", result.ProfileDraft)
	}
	requireFloatPtrNear(t, "front tire pressure", result.ProfileDraft.FrontTirePressure, 2.03, 0.001)
	requireFloatPtrNear(t, "rear tire pressure", result.ProfileDraft.RearTirePressure, 2.03, 0.001)
	requireFloatPtrNear(t, "front camber", result.ProfileDraft.FrontCamber, -1.2, 0.001)
	requireFloatPtrNear(t, "rear camber", result.ProfileDraft.RearCamber, -0.8, 0.001)
	requireFloatPtrNear(t, "caster", result.ProfileDraft.Caster, 5.8, 0.001)
	requireFloatPtrNear(t, "front arb", result.ProfileDraft.FrontARB, 10, 0.001)
	requireFloatPtrNear(t, "rear arb", result.ProfileDraft.RearARB, 12, 0.001)
	requireFloatPtrNear(t, "brake pressure", result.ProfileDraft.BrakePressure, 95, 0.001)
	requireFloatPtrNear(t, "front diff accel", result.ProfileDraft.FrontDiffAccel, 30, 0.001)
	requireFloatPtrNear(t, "rear diff accel", result.ProfileDraft.RearDiffAccel, 60, 0.001)
	requireFloatPtrNear(t, "center diff", result.ProfileDraft.CenterDiffBalance, 65, 0.001)
	if !hasTierRecommendation(result.TierRecommendations, "frontRideHeight", "medium", false) {
		t.Fatalf("rally ride height tier missing: %#v", result.TierRecommendations)
	}
}

func TestGenerateRoadStaticTuneBaselineOffroadAWD(t *testing.T) {
	result, err := GenerateRoadStaticTuneBaseline(RoadStaticTuneBaselineInput{
		UseCase:        "Offroad",
		PI:             700,
		Drivetrain:     "AWD",
		WeightKG:       1400,
		FrontWeightPct: 54,
	})
	if err != nil {
		t.Fatalf("generate offroad baseline: %v", err)
	}
	requireFloatPtrNear(t, "front tire pressure", result.ProfileDraft.FrontTirePressure, 2.00, 0.001)
	requireFloatPtrNear(t, "rear tire pressure", result.ProfileDraft.RearTirePressure, 2.00, 0.001)
	requireFloatPtrNear(t, "front camber", result.ProfileDraft.FrontCamber, -0.8, 0.001)
	requireFloatPtrNear(t, "rear camber", result.ProfileDraft.RearCamber, -0.5, 0.001)
	requireFloatPtrNear(t, "front arb", result.ProfileDraft.FrontARB, 6, 0.001)
	requireFloatPtrNear(t, "rear arb", result.ProfileDraft.RearARB, 8, 0.001)
	requireFloatPtrNear(t, "brake pressure", result.ProfileDraft.BrakePressure, 95, 0.001)
	requireFloatPtrNear(t, "front diff accel", result.ProfileDraft.FrontDiffAccel, 35, 0.001)
	requireFloatPtrNear(t, "rear diff accel", result.ProfileDraft.RearDiffAccel, 70, 0.001)
	requireFloatPtrNear(t, "center diff", result.ProfileDraft.CenterDiffBalance, 60, 0.001)
	if !hasTierRecommendation(result.TierRecommendations, "frontRideHeight", "high", false) {
		t.Fatalf("offroad ride height tier missing: %#v", result.TierRecommendations)
	}
}

func TestGenerateRoadStaticTuneBaselineDragByDrivetrain(t *testing.T) {
	cases := []struct {
		name          string
		drivetrain    string
		frontPressure float64
		rearPressure  float64
		frontDiff     *float64
		rearDiff      *float64
		centerDiff    *float64
		brakeBalance  float64
	}{
		{name: "rwd", drivetrain: "RWD", frontPressure: 2.20, rearPressure: 1.72, rearDiff: testFloatPtr(85), brakeBalance: 50},
		{name: "fwd", drivetrain: "FWD", frontPressure: 1.72, rearPressure: 2.20, frontDiff: testFloatPtr(85), brakeBalance: 58},
		{name: "awd", drivetrain: "AWD", frontPressure: 1.72, rearPressure: 1.72, frontDiff: testFloatPtr(35), rearDiff: testFloatPtr(85), centerDiff: testFloatPtr(75), brakeBalance: 54},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := GenerateRoadStaticTuneBaseline(RoadStaticTuneBaselineInput{
				UseCase:        "Drag",
				PI:             700,
				Drivetrain:     tc.drivetrain,
				WeightKG:       1400,
				FrontWeightPct: 54,
			})
			if err != nil {
				t.Fatalf("generate drag baseline: %v", err)
			}
			requireFloatPtrNear(t, "front tire pressure", result.ProfileDraft.FrontTirePressure, tc.frontPressure, 0.001)
			requireFloatPtrNear(t, "rear tire pressure", result.ProfileDraft.RearTirePressure, tc.rearPressure, 0.001)
			requireFloatPtrNear(t, "front camber", result.ProfileDraft.FrontCamber, 0, 0.001)
			requireFloatPtrNear(t, "rear camber", result.ProfileDraft.RearCamber, 0, 0.001)
			requireFloatPtrNear(t, "caster", result.ProfileDraft.Caster, 6, 0.001)
			requireFloatPtrNear(t, "brake balance", result.ProfileDraft.BrakeBalance, tc.brakeBalance, 0.001)
			if tc.frontDiff != nil {
				requireFloatPtrNear(t, "front diff accel", result.ProfileDraft.FrontDiffAccel, *tc.frontDiff, 0.001)
			}
			if tc.rearDiff != nil {
				requireFloatPtrNear(t, "rear diff accel", result.ProfileDraft.RearDiffAccel, *tc.rearDiff, 0.001)
			}
			if tc.centerDiff != nil {
				requireFloatPtrNear(t, "center diff", result.ProfileDraft.CenterDiffBalance, *tc.centerDiff, 0.001)
			}
		})
	}
}

func TestGenerateRoadStaticTuneBaselineDragGearing(t *testing.T) {
	redline := 7200.0
	gearCount := int64(6)
	tireDiameter := 67.0
	targetTrapSpeed := 260.0
	result, err := GenerateRoadStaticTuneBaseline(RoadStaticTuneBaselineInput{
		UseCase:           "Drag",
		PI:                800,
		Drivetrain:        "RWD",
		WeightKG:          1350,
		FrontWeightPct:    53,
		RedlineRPM:        &redline,
		GearCount:         &gearCount,
		TireDiameterCm:    &tireDiameter,
		TargetTopSpeedKmh: &targetTrapSpeed,
	})
	if err != nil {
		t.Fatalf("generate drag gearing: %v", err)
	}
	if result.ProfileDraft.FinalDrive == nil || result.ProfileDraft.Gear1 == nil || result.ProfileDraft.Gear6 == nil {
		t.Fatalf("drag gearing missing: %#v", result.ProfileDraft)
	}
	topRPMRatio := gearRPMAtSpeedKmh(targetTrapSpeed, tireDiameter, *result.ProfileDraft.FinalDrive, *result.ProfileDraft.Gear6) / redline
	if topRPMRatio < 0.93 || topRPMRatio > 0.99 {
		t.Fatalf("drag top gear rpm ratio = %.3f, want near 0.96", topRPMRatio)
	}
	if result.ProfileDraft.Gear7 != nil || !hasSkippedField(result.SkippedFields, "gear7") {
		t.Fatalf("locked gear should be skipped: gear7=%#v skipped=%#v", result.ProfileDraft.Gear7, result.SkippedFields)
	}
}

func TestGenerateRoadStaticTuneBaselineSpringFormulaV11(t *testing.T) {
	result, err := GenerateRoadStaticTuneBaseline(RoadStaticTuneBaselineInput{
		CarName:        "Spring Formula",
		PI:             700,
		Drivetrain:     "AWD",
		WeightKG:       1717,
		FrontWeightPct: 52,
	})
	if err != nil {
		t.Fatalf("generate baseline: %v", err)
	}
	requireFloatPtrNear(t, "front spring", result.ProfileDraft.FrontSpring, 145.8, 0.05)
	requireFloatPtrNear(t, "rear spring", result.ProfileDraft.RearSpring, 138.5, 0.05)
	requireFloatPtrNear(t, "front rebound", result.ProfileDraft.FrontRebound, 11.9, 0.05)
	requireFloatPtrNear(t, "rear rebound", result.ProfileDraft.RearRebound, 11.8, 0.05)
	requireFloatPtrNear(t, "front bump", result.ProfileDraft.FrontBump, 7.1, 0.05)
	requireFloatPtrNear(t, "rear bump", result.ProfileDraft.RearBump, 7.1, 0.05)
}

func TestGenerateRoadStaticTuneBaselineRoadMVPV12(t *testing.T) {
	result, err := GenerateRoadStaticTuneBaseline(RoadStaticTuneBaselineInput{
		CarName:        "Road MVP",
		PI:             700,
		Drivetrain:     "AWD",
		WeightKG:       1717,
		FrontWeightPct: 52,
	})
	if err != nil {
		t.Fatalf("generate baseline: %v", err)
	}
	if result.ProfileDraft.CarClass != "A" {
		t.Fatalf("car class = %q, want A", result.ProfileDraft.CarClass)
	}
	requireFloatPtrNear(t, "front tire pressure", result.ProfileDraft.FrontTirePressure, 2.24, 0.001)
	requireFloatPtrNear(t, "rear tire pressure", result.ProfileDraft.RearTirePressure, 2.21, 0.001)
	requireFloatPtrNear(t, "front camber", result.ProfileDraft.FrontCamber, -1.5, 0.001)
	requireFloatPtrNear(t, "rear camber", result.ProfileDraft.RearCamber, -1.0, 0.001)
	requireFloatPtrNear(t, "front toe", result.ProfileDraft.FrontToe, 0, 0.001)
	requireFloatPtrNear(t, "rear toe", result.ProfileDraft.RearToe, 0, 0.001)
	requireFloatPtrNear(t, "caster", result.ProfileDraft.Caster, 5.6, 0.001)
	requireFloatPtrNear(t, "front spring", result.ProfileDraft.FrontSpring, 145.8, 0.05)
	requireFloatPtrNear(t, "rear spring", result.ProfileDraft.RearSpring, 138.5, 0.05)
	requireFloatPtrNear(t, "front rebound", result.ProfileDraft.FrontRebound, 11.9, 0.05)
	requireFloatPtrNear(t, "rear rebound", result.ProfileDraft.RearRebound, 11.8, 0.05)
	requireFloatPtrNear(t, "front bump", result.ProfileDraft.FrontBump, 7.1, 0.05)
	requireFloatPtrNear(t, "rear bump", result.ProfileDraft.RearBump, 7.1, 0.05)
	requireFloatPtrNear(t, "front arb", result.ProfileDraft.FrontARB, 25, 0.001)
	requireFloatPtrNear(t, "rear arb", result.ProfileDraft.RearARB, 33, 0.001)
	requireFloatPtrNear(t, "brake balance", result.ProfileDraft.BrakeBalance, 54, 0.001)
	requireFloatPtrNear(t, "brake pressure", result.ProfileDraft.BrakePressure, 100, 0.001)
	requireFloatPtrNear(t, "front diff accel", result.ProfileDraft.FrontDiffAccel, 21, 0.001)
	requireFloatPtrNear(t, "front diff decel", result.ProfileDraft.FrontDiffDecel, 6, 0.001)
	requireFloatPtrNear(t, "rear diff accel", result.ProfileDraft.RearDiffAccel, 54, 0.001)
	requireFloatPtrNear(t, "rear diff decel", result.ProfileDraft.RearDiffDecel, 25, 0.001)
	requireFloatPtrNear(t, "center diff balance", result.ProfileDraft.CenterDiffBalance, 65, 0.001)
}

func TestGenerateRoadStaticTuneBaselineTirePressureFormula(t *testing.T) {
	cases := []struct {
		name       string
		pi         int64
		drivetrain string
		compound   string
		weightKG   float64
		frontPct   float64
		wantFront  float64
		wantRear   float64
	}{
		{name: "default sport AWD", pi: 600, drivetrain: "AWD", weightKG: 1460, frontPct: 60, wantFront: 2.17, wantRear: 2.14},
		{name: "heavy AWD sport", pi: 700, drivetrain: "AWD", compound: "sport", weightKG: 1717, frontPct: 52, wantFront: 2.24, wantRear: 2.21},
		{name: "light RWD sport", pi: 900, drivetrain: "RWD", compound: "sport", weightKG: 1050, frontPct: 50, wantFront: 2.17, wantRear: 2.11},
		{name: "FWD sport", pi: 700, drivetrain: "FWD", compound: "sport", weightKG: 1300, frontPct: 60, wantFront: 2.24, wantRear: 2.14},
		{name: "AWD semi slick heavy", pi: 800, drivetrain: "AWD", compound: "semi", weightKG: 1717, frontPct: 52, wantFront: 2.28, wantRear: 2.25},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := GenerateRoadStaticTuneBaseline(RoadStaticTuneBaselineInput{
				PI:             tc.pi,
				Drivetrain:     tc.drivetrain,
				TireCompound:   tc.compound,
				WeightKG:       tc.weightKG,
				FrontWeightPct: tc.frontPct,
			})
			if err != nil {
				t.Fatalf("generate baseline: %v", err)
			}
			requireFloatPtrNear(t, "front tire pressure", result.ProfileDraft.FrontTirePressure, tc.wantFront, 0.001)
			requireFloatPtrNear(t, "rear tire pressure", result.ProfileDraft.RearTirePressure, tc.wantRear, 0.001)
		})
	}
}

func TestGenerateRoadStaticTuneBaselineTirePressureIgnoresBalanceBias(t *testing.T) {
	base := RoadStaticTuneBaselineInput{
		PI:             700,
		Drivetrain:     "AWD",
		WeightKG:       1400,
		FrontWeightPct: 54,
	}
	neutral, err := GenerateRoadStaticTuneBaseline(base)
	if err != nil {
		t.Fatalf("generate neutral baseline: %v", err)
	}
	base.BalanceBias = 150
	agile, err := GenerateRoadStaticTuneBaseline(base)
	if err != nil {
		t.Fatalf("generate agile baseline: %v", err)
	}
	if *agile.ProfileDraft.FrontTirePressure != *neutral.ProfileDraft.FrontTirePressure || *agile.ProfileDraft.RearTirePressure != *neutral.ProfileDraft.RearTirePressure {
		t.Fatalf("balance bias should not change tire pressure: neutral %.2f/%.2f agile %.2f/%.2f",
			*neutral.ProfileDraft.FrontTirePressure, *neutral.ProfileDraft.RearTirePressure, *agile.ProfileDraft.FrontTirePressure, *agile.ProfileDraft.RearTirePressure)
	}
	base.BalanceBias = 50
	stable, err := GenerateRoadStaticTuneBaseline(base)
	if err != nil {
		t.Fatalf("generate stable baseline: %v", err)
	}
	if *stable.ProfileDraft.FrontTirePressure != *neutral.ProfileDraft.FrontTirePressure || *stable.ProfileDraft.RearTirePressure != *neutral.ProfileDraft.RearTirePressure {
		t.Fatalf("balance bias should not change tire pressure: neutral %.2f/%.2f stable %.2f/%.2f",
			*neutral.ProfileDraft.FrontTirePressure, *neutral.ProfileDraft.RearTirePressure, *stable.ProfileDraft.FrontTirePressure, *stable.ProfileDraft.RearTirePressure)
	}
}

func TestGenerateRoadStaticTuneBaselineDriftTirePressureRules(t *testing.T) {
	base := RoadStaticTuneBaselineInput{
		UseCase:        "Drift",
		PI:             700,
		Drivetrain:     "RWD",
		WeightKG:       1400,
		FrontWeightPct: 54,
	}
	neutral, err := GenerateRoadStaticTuneBaseline(base)
	if err != nil {
		t.Fatalf("generate neutral drift baseline: %v", err)
	}
	requireFloatPtrNear(t, "front drift tire pressure", neutral.ProfileDraft.FrontTirePressure, 2.10, 0.001)
	requireFloatPtrNear(t, "rear drift tire pressure", neutral.ProfileDraft.RearTirePressure, 2.50, 0.001)
	if diff := *neutral.ProfileDraft.RearTirePressure - *neutral.ProfileDraft.FrontTirePressure; math.Abs(diff-0.40) > 0.03 {
		t.Fatalf("drift tire pressure split = %.2f, want about 0.40", diff)
	}

	heavy := base
	heavy.WeightKG = 1900
	heavyResult, err := GenerateRoadStaticTuneBaseline(heavy)
	if err != nil {
		t.Fatalf("generate heavy drift baseline: %v", err)
	}
	requireFloatPtrNear(t, "heavy front drift tire pressure", heavyResult.ProfileDraft.FrontTirePressure, 2.20, 0.001)
	requireFloatPtrNear(t, "heavy rear drift tire pressure", heavyResult.ProfileDraft.RearTirePressure, 2.60, 0.001)
	if diff := *heavyResult.ProfileDraft.RearTirePressure - *heavyResult.ProfileDraft.FrontTirePressure; math.Abs(diff-0.40) > 0.03 {
		t.Fatalf("heavy drift tire pressure split = %.2f, want about 0.40", diff)
	}

	biased := base
	biased.BalanceBias = 150
	biased.StiffnessBias = 150
	biased.SpeedBias = 150
	biasedResult, err := GenerateRoadStaticTuneBaseline(biased)
	if err != nil {
		t.Fatalf("generate biased drift baseline: %v", err)
	}
	requireFloatPtrNear(t, "biased front drift tire pressure", biasedResult.ProfileDraft.FrontTirePressure, *neutral.ProfileDraft.FrontTirePressure, 0.001)
	requireFloatPtrNear(t, "biased rear drift tire pressure", biasedResult.ProfileDraft.RearTirePressure, *neutral.ProfileDraft.RearTirePressure, 0.001)
}

func TestGenerateRoadStaticTuneBaselineDriftRWD(t *testing.T) {
	result, err := GenerateRoadStaticTuneBaseline(RoadStaticTuneBaselineInput{
		UseCase:        "Drift",
		PI:             700,
		Drivetrain:     "RWD",
		WeightKG:       1400,
		FrontWeightPct: 54,
	})
	if err != nil {
		t.Fatalf("generate drift baseline: %v", err)
	}
	if result.ProfileDraft.UseCase != "Drift" || result.ProfileDraft.VersionName != "Drift Baseline" {
		t.Fatalf("drift identity = %#v", result.ProfileDraft)
	}
	requireFloatPtrNear(t, "front tire pressure", result.ProfileDraft.FrontTirePressure, 2.10, 0.001)
	requireFloatPtrNear(t, "rear tire pressure", result.ProfileDraft.RearTirePressure, 2.50, 0.001)
	requireFloatPtrNear(t, "front camber", result.ProfileDraft.FrontCamber, -4.5, 0.001)
	requireFloatPtrNear(t, "rear camber", result.ProfileDraft.RearCamber, -0.5, 0.001)
	requireFloatPtrNear(t, "front toe", result.ProfileDraft.FrontToe, 1.0, 0.001)
	requireFloatPtrNear(t, "rear toe", result.ProfileDraft.RearToe, 0.0, 0.001)
	requireFloatPtrNear(t, "caster", result.ProfileDraft.Caster, 7.0, 0.001)
	requireFloatPtrNear(t, "front spring", result.ProfileDraft.FrontSpring, 94.4, 0.05)
	requireFloatPtrNear(t, "rear spring", result.ProfileDraft.RearSpring, 83.0, 0.05)
	requireFloatPtrNear(t, "front rebound", result.ProfileDraft.FrontRebound, 9.8, 0.05)
	requireFloatPtrNear(t, "rear rebound", result.ProfileDraft.RearRebound, 9.7, 0.05)
	requireFloatPtrNear(t, "front bump", result.ProfileDraft.FrontBump, 5.9, 0.05)
	requireFloatPtrNear(t, "rear bump", result.ProfileDraft.RearBump, 5.8, 0.05)
	requireFloatPtrNear(t, "front arb", result.ProfileDraft.FrontARB, 19, 0.001)
	requireFloatPtrNear(t, "rear arb", result.ProfileDraft.RearARB, 23, 0.001)
	requireFloatPtrNear(t, "brake balance", result.ProfileDraft.BrakeBalance, 48, 0.001)
	requireFloatPtrNear(t, "rear diff accel", result.ProfileDraft.RearDiffAccel, 100, 0.001)
	requireFloatPtrNear(t, "rear diff decel", result.ProfileDraft.RearDiffDecel, 100, 0.001)
	if !hasTierRecommendation(result.TierRecommendations, "frontAero", "low", false) {
		t.Fatalf("drift aero should be a low tier recommendation: %#v", result.TierRecommendations)
	}
}

func TestGenerateRoadStaticTuneBaselineDriftNonRWDDoesNotInventDifferential(t *testing.T) {
	result, err := GenerateRoadStaticTuneBaseline(RoadStaticTuneBaselineInput{
		UseCase:        "Drift",
		PI:             700,
		Drivetrain:     "AWD",
		WeightKG:       1400,
		FrontWeightPct: 54,
	})
	if err != nil {
		t.Fatalf("generate drift AWD baseline: %v", err)
	}
	if result.ProfileDraft.FrontDiffAccel != nil || result.ProfileDraft.RearDiffAccel != nil || result.ProfileDraft.CenterDiffBalance != nil {
		t.Fatalf("drift AWD should not generate differential values: %#v", result.ProfileDraft)
	}
	if !hasSkippedField(result.SkippedFields, "centerDiffBalance") || !hasSkippedField(result.SkippedFields, "rearDiffAccel") {
		t.Fatalf("drift AWD should explain skipped diff fields: %#v", result.SkippedFields)
	}
	if !hasGeneratedField(result.GeneratedFields, "frontSpring") || !hasGeneratedField(result.GeneratedFields, "frontArb") {
		t.Fatalf("drift AWD should still generate core fields: %#v", result.GeneratedFields)
	}
}

func TestGenerateRoadStaticTuneBaselineRejectsUnsupportedUseCase(t *testing.T) {
	_, err := GenerateRoadStaticTuneBaseline(RoadStaticTuneBaselineInput{
		UseCase:        "TimeAttack",
		PI:             700,
		Drivetrain:     "AWD",
		WeightKG:       1400,
		FrontWeightPct: 54,
	})
	if err == nil || !strings.Contains(err.Error(), "use case must be Road, Drift, Rally, Offroad, or Drag") {
		t.Fatalf("expected unsupported use case error, got %v", err)
	}
}

func TestGenerateRoadStaticTuneBaselineDifferentialFormula(t *testing.T) {
	t.Run("class and drivetrain affect diff baseline", func(t *testing.T) {
		bClass, err := GenerateRoadStaticTuneBaseline(RoadStaticTuneBaselineInput{
			PI:             600,
			Drivetrain:     "AWD",
			WeightKG:       1460,
			FrontWeightPct: 60,
		})
		if err != nil {
			t.Fatalf("generate B AWD baseline: %v", err)
		}
		s2Class, err := GenerateRoadStaticTuneBaseline(RoadStaticTuneBaselineInput{
			PI:             900,
			Drivetrain:     "AWD",
			WeightKG:       1460,
			FrontWeightPct: 60,
		})
		if err != nil {
			t.Fatalf("generate S2 AWD baseline: %v", err)
		}
		if *s2Class.ProfileDraft.RearDiffAccel <= *bClass.ProfileDraft.RearDiffAccel || *s2Class.ProfileDraft.CenterDiffBalance <= *bClass.ProfileDraft.CenterDiffBalance {
			t.Fatalf("higher PI AWD should use more rear lock and rear-biased center split: B=%#v S2=%#v", bClass.ProfileDraft, s2Class.ProfileDraft)
		}
	})

	t.Run("balance bias moves diff toward stability or rotation", func(t *testing.T) {
		base := RoadStaticTuneBaselineInput{
			PI:             700,
			Drivetrain:     "AWD",
			WeightKG:       1400,
			FrontWeightPct: 54,
		}
		neutral, err := GenerateRoadStaticTuneBaseline(base)
		if err != nil {
			t.Fatalf("generate neutral baseline: %v", err)
		}
		base.BalanceBias = 150
		agile, err := GenerateRoadStaticTuneBaseline(base)
		if err != nil {
			t.Fatalf("generate agile baseline: %v", err)
		}
		if *agile.ProfileDraft.FrontDiffAccel >= *neutral.ProfileDraft.FrontDiffAccel || *agile.ProfileDraft.RearDiffAccel <= *neutral.ProfileDraft.RearDiffAccel || *agile.ProfileDraft.CenterDiffBalance <= *neutral.ProfileDraft.CenterDiffBalance {
			t.Fatalf("agile bias should lower front lock and raise rear/center lock: neutral=%#v agile=%#v", neutral.ProfileDraft, agile.ProfileDraft)
		}
		base.BalanceBias = 50
		stable, err := GenerateRoadStaticTuneBaseline(base)
		if err != nil {
			t.Fatalf("generate stable baseline: %v", err)
		}
		if *stable.ProfileDraft.FrontDiffAccel <= *neutral.ProfileDraft.FrontDiffAccel || *stable.ProfileDraft.RearDiffAccel >= *neutral.ProfileDraft.RearDiffAccel || *stable.ProfileDraft.CenterDiffBalance >= *neutral.ProfileDraft.CenterDiffBalance {
			t.Fatalf("stable bias should raise front lock and lower rear/center lock: neutral=%#v stable=%#v", neutral.ProfileDraft, stable.ProfileDraft)
		}
	})
}

func TestGenerateRoadStaticTuneBaselineDefaultBiasMatchesNeutral(t *testing.T) {
	neutral, err := GenerateRoadStaticTuneBaseline(RoadStaticTuneBaselineInput{
		PI:             700,
		Drivetrain:     "AWD",
		WeightKG:       1717,
		FrontWeightPct: 52,
	})
	if err != nil {
		t.Fatalf("generate neutral baseline: %v", err)
	}
	explicit, err := GenerateRoadStaticTuneBaseline(RoadStaticTuneBaselineInput{
		PI:             700,
		Drivetrain:     "AWD",
		WeightKG:       1717,
		FrontWeightPct: 52,
		BalanceBias:    100,
		StiffnessBias:  100,
		SpeedBias:      100,
	})
	if err != nil {
		t.Fatalf("generate explicit neutral baseline: %v", err)
	}
	requireFloatPtrNear(t, "front arb", explicit.ProfileDraft.FrontARB, *neutral.ProfileDraft.FrontARB, 0.001)
	requireFloatPtrNear(t, "rear arb", explicit.ProfileDraft.RearARB, *neutral.ProfileDraft.RearARB, 0.001)
	requireFloatPtrNear(t, "front spring", explicit.ProfileDraft.FrontSpring, *neutral.ProfileDraft.FrontSpring, 0.001)
	requireFloatPtrNear(t, "front rebound", explicit.ProfileDraft.FrontRebound, *neutral.ProfileDraft.FrontRebound, 0.001)
}

func TestGenerateRoadStaticTuneBaselineBalanceBiasAdjustsAxleRatio(t *testing.T) {
	base := RoadStaticTuneBaselineInput{
		PI:             700,
		Drivetrain:     "AWD",
		WeightKG:       1717,
		FrontWeightPct: 52,
	}
	neutral, err := GenerateRoadStaticTuneBaseline(base)
	if err != nil {
		t.Fatalf("generate neutral baseline: %v", err)
	}
	base.BalanceBias = 150
	agile, err := GenerateRoadStaticTuneBaseline(base)
	if err != nil {
		t.Fatalf("generate agile baseline: %v", err)
	}
	if *agile.ProfileDraft.FrontARB >= *neutral.ProfileDraft.FrontARB || *agile.ProfileDraft.RearARB <= *neutral.ProfileDraft.RearARB {
		t.Fatalf("agile balance should reduce front ARB and raise rear ARB: neutral front/rear %.2f/%.2f agile %.2f/%.2f",
			*neutral.ProfileDraft.FrontARB, *neutral.ProfileDraft.RearARB, *agile.ProfileDraft.FrontARB, *agile.ProfileDraft.RearARB)
	}
	base.BalanceBias = 50
	stable, err := GenerateRoadStaticTuneBaseline(base)
	if err != nil {
		t.Fatalf("generate stable baseline: %v", err)
	}
	if *stable.ProfileDraft.FrontARB <= *neutral.ProfileDraft.FrontARB || *stable.ProfileDraft.RearARB >= *neutral.ProfileDraft.RearARB {
		t.Fatalf("stable balance should raise front ARB and reduce rear ARB: neutral front/rear %.2f/%.2f stable %.2f/%.2f",
			*neutral.ProfileDraft.FrontARB, *neutral.ProfileDraft.RearARB, *stable.ProfileDraft.FrontARB, *stable.ProfileDraft.RearARB)
	}
}

func TestGenerateRoadStaticTuneBaselineStiffnessBiasScalesSupport(t *testing.T) {
	base := RoadStaticTuneBaselineInput{
		PI:             800,
		Drivetrain:     "RWD",
		WeightKG:       1350,
		FrontWeightPct: 53,
	}
	neutral, err := GenerateRoadStaticTuneBaseline(base)
	if err != nil {
		t.Fatalf("generate neutral baseline: %v", err)
	}
	base.StiffnessBias = 150
	hard, err := GenerateRoadStaticTuneBaseline(base)
	if err != nil {
		t.Fatalf("generate hard baseline: %v", err)
	}
	if *hard.ProfileDraft.FrontSpring <= *neutral.ProfileDraft.FrontSpring || *hard.ProfileDraft.FrontRebound <= *neutral.ProfileDraft.FrontRebound || *hard.ProfileDraft.FrontARB <= *neutral.ProfileDraft.FrontARB {
		t.Fatalf("hard stiffness should raise support: neutral=%#v hard=%#v", neutral.ProfileDraft, hard.ProfileDraft)
	}
	base.StiffnessBias = 50
	soft, err := GenerateRoadStaticTuneBaseline(base)
	if err != nil {
		t.Fatalf("generate soft baseline: %v", err)
	}
	if *soft.ProfileDraft.FrontSpring >= *neutral.ProfileDraft.FrontSpring || *soft.ProfileDraft.FrontRebound >= *neutral.ProfileDraft.FrontRebound || *soft.ProfileDraft.FrontARB >= *neutral.ProfileDraft.FrontARB {
		t.Fatalf("soft stiffness should lower support: neutral=%#v soft=%#v", neutral.ProfileDraft, soft.ProfileDraft)
	}
}

func TestGenerateRoadStaticTuneBaselineSpeedBiasAdjustsGeneratedGearingOnly(t *testing.T) {
	redline := 7200.0
	gearCount := int64(6)
	tireDiameter := 67.0
	targetTopSpeed := 280.0
	base := RoadStaticTuneBaselineInput{
		PI:                800,
		Drivetrain:        "RWD",
		WeightKG:          1350,
		FrontWeightPct:    53,
		RedlineRPM:        &redline,
		GearCount:         &gearCount,
		TireDiameterCm:    &tireDiameter,
		TargetTopSpeedKmh: &targetTopSpeed,
	}
	neutral, err := GenerateRoadStaticTuneBaseline(base)
	if err != nil {
		t.Fatalf("generate neutral baseline: %v", err)
	}
	base.SpeedBias = 50
	topSpeed, err := GenerateRoadStaticTuneBaseline(base)
	if err != nil {
		t.Fatalf("generate top speed baseline: %v", err)
	}
	if *topSpeed.ProfileDraft.FinalDrive >= *neutral.ProfileDraft.FinalDrive {
		t.Fatalf("speed bias 50 should lower final drive: neutral final %.2f top %.2f",
			*neutral.ProfileDraft.FinalDrive, *topSpeed.ProfileDraft.FinalDrive)
	}
	if *topSpeed.ProfileDraft.Gear1 != *neutral.ProfileDraft.Gear1 || *topSpeed.ProfileDraft.Gear6 != *neutral.ProfileDraft.Gear6 {
		t.Fatalf("speed bias 50 should preserve gear spacing: neutral gear1/gear6 %.2f/%.2f top %.2f/%.2f",
			*neutral.ProfileDraft.Gear1, *neutral.ProfileDraft.Gear6, *topSpeed.ProfileDraft.Gear1, *topSpeed.ProfileDraft.Gear6)
	}
	base.SpeedBias = 150
	accel, err := GenerateRoadStaticTuneBaseline(base)
	if err != nil {
		t.Fatalf("generate acceleration baseline: %v", err)
	}
	if *accel.ProfileDraft.FinalDrive <= *neutral.ProfileDraft.FinalDrive {
		t.Fatalf("speed bias 150 should raise final drive: neutral final %.2f accel %.2f",
			*neutral.ProfileDraft.FinalDrive, *accel.ProfileDraft.FinalDrive)
	}
	if *accel.ProfileDraft.Gear1 != *neutral.ProfileDraft.Gear1 || *accel.ProfileDraft.Gear6 != *neutral.ProfileDraft.Gear6 {
		t.Fatalf("speed bias 150 should preserve gear spacing: neutral gear1/gear6 %.2f/%.2f accel %.2f/%.2f",
			*neutral.ProfileDraft.Gear1, *neutral.ProfileDraft.Gear6, *accel.ProfileDraft.Gear1, *accel.ProfileDraft.Gear6)
	}

	noGearing, err := GenerateRoadStaticTuneBaseline(RoadStaticTuneBaselineInput{
		PI:             800,
		Drivetrain:     "RWD",
		WeightKG:       1350,
		FrontWeightPct: 53,
		SpeedBias:      150,
	})
	if err != nil {
		t.Fatalf("generate no-gearing baseline: %v", err)
	}
	if noGearing.ProfileDraft.FinalDrive != nil || noGearing.ProfileDraft.Gear1 != nil {
		t.Fatalf("speed bias should not generate gearing when gearing inputs are missing: %#v", noGearing.ProfileDraft)
	}
}

func TestGenerateRoadStaticTuneBaselineDampingUsesSpringFactor(t *testing.T) {
	result, err := GenerateRoadStaticTuneBaseline(RoadStaticTuneBaselineInput{
		CarName:        "High PI Damping",
		PI:             900,
		Drivetrain:     "AWD",
		WeightKG:       1400,
		FrontWeightPct: 50,
	})
	if err != nil {
		t.Fatalf("generate baseline: %v", err)
	}
	requireFloatPtrNear(t, "front spring", result.ProfileDraft.FrontSpring, 148.8, 0.15)
	requireFloatPtrNear(t, "rear spring", result.ProfileDraft.RearSpring, 148.8, 0.15)
	requireFloatPtrNear(t, "front rebound", result.ProfileDraft.FrontRebound, 14.0, 0.05)
	requireFloatPtrNear(t, "rear rebound", result.ProfileDraft.RearRebound, 13.8, 0.05)
	requireFloatPtrNear(t, "front bump", result.ProfileDraft.FrontBump, 8.4, 0.05)
	requireFloatPtrNear(t, "rear bump", result.ProfileDraft.RearBump, 8.3, 0.05)
}

func TestGenerateRoadStaticTuneBaselineSpringFormulaByClass(t *testing.T) {
	cases := []struct {
		name      string
		weightKG  float64
		frontPct  float64
		pi        int64
		wantFront float64
		wantRear  float64
	}{
		{name: "S1 800", weightKG: 1500, frontPct: 55, pi: 800, wantFront: 151.2, wantRear: 132.7},
		{name: "S2 900", weightKG: 1400, frontPct: 50, pi: 900, wantFront: 148.8, wantRear: 148.8},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := GenerateRoadStaticTuneBaseline(RoadStaticTuneBaselineInput{
				CarName:        tc.name,
				PI:             tc.pi,
				Drivetrain:     "AWD",
				WeightKG:       tc.weightKG,
				FrontWeightPct: tc.frontPct,
			})
			if err != nil {
				t.Fatalf("generate baseline: %v", err)
			}
			requireFloatPtrNear(t, "front spring", result.ProfileDraft.FrontSpring, tc.wantFront, 0.15)
			requireFloatPtrNear(t, "rear spring", result.ProfileDraft.RearSpring, tc.wantRear, 0.15)
			if !isStepValue(*result.ProfileDraft.FrontSpring, 0.1) || !isStepValue(*result.ProfileDraft.RearSpring, 0.1) {
				t.Fatalf("spring values should use 0.1 step: front=%v rear=%v", *result.ProfileDraft.FrontSpring, *result.ProfileDraft.RearSpring)
			}
		})
	}
}

func TestGenerateRoadStaticTuneBaselineFrontHeavyRoadBalanceCompensation(t *testing.T) {
	result, err := GenerateRoadStaticTuneBaseline(RoadStaticTuneBaselineInput{
		CarName:        "Front Heavy AWD",
		PI:             600,
		Drivetrain:     "AWD",
		WeightKG:       1460,
		FrontWeightPct: 60,
	})
	if err != nil {
		t.Fatalf("generate baseline: %v", err)
	}
	requireFloatPtrNear(t, "front spring", result.ProfileDraft.FrontSpring, 113.5, 0.05)
	requireFloatPtrNear(t, "rear spring", result.ProfileDraft.RearSpring, 86.9, 0.05)
	requireFloatPtrNear(t, "front rebound", result.ProfileDraft.FrontRebound, 10.1, 0.05)
	requireFloatPtrNear(t, "rear rebound", result.ProfileDraft.RearRebound, 10.9, 0.05)
	requireFloatPtrNear(t, "front bump", result.ProfileDraft.FrontBump, 6.1, 0.05)
	requireFloatPtrNear(t, "rear bump", result.ProfileDraft.RearBump, 6.5, 0.05)
	requireFloatPtrNear(t, "front arb", result.ProfileDraft.FrontARB, 23, 0.001)
	requireFloatPtrNear(t, "rear arb", result.ProfileDraft.RearARB, 33, 0.001)
	requireFloatPtrNear(t, "front camber", result.ProfileDraft.FrontCamber, -1.4, 0.001)
	requireFloatPtrNear(t, "rear camber", result.ProfileDraft.RearCamber, -0.7, 0.001)
	requireFloatPtrNear(t, "front toe", result.ProfileDraft.FrontToe, 0, 0.001)
	requireFloatPtrNear(t, "rear toe", result.ProfileDraft.RearToe, 0, 0.001)
	requireFloatPtrNear(t, "caster", result.ProfileDraft.Caster, 5.3, 0.001)
	if *result.ProfileDraft.RearRebound <= *result.ProfileDraft.FrontRebound || *result.ProfileDraft.RearBump <= *result.ProfileDraft.FrontBump {
		t.Fatalf("front-heavy compensation should keep rear damping responsive: %#v", result.ProfileDraft)
	}
}

func TestGenerateRoadStaticTuneBaselineAlignmentFormula(t *testing.T) {
	t.Run("agile balance adds front toe out and rotation camber", func(t *testing.T) {
		result, err := GenerateRoadStaticTuneBaseline(RoadStaticTuneBaselineInput{
			PI:             700,
			Drivetrain:     "AWD",
			WeightKG:       1717,
			FrontWeightPct: 52,
			BalanceBias:    150,
		})
		if err != nil {
			t.Fatalf("generate baseline: %v", err)
		}
		requireFloatPtrNear(t, "front camber", result.ProfileDraft.FrontCamber, -1.6, 0.001)
		requireFloatPtrNear(t, "rear camber", result.ProfileDraft.RearCamber, -0.9, 0.001)
		requireFloatPtrNear(t, "front toe", result.ProfileDraft.FrontToe, -0.1, 0.001)
		requireFloatPtrNear(t, "rear toe", result.ProfileDraft.RearToe, 0, 0.001)
	})
	t.Run("stable balance adds rear toe in", func(t *testing.T) {
		result, err := GenerateRoadStaticTuneBaseline(RoadStaticTuneBaselineInput{
			PI:             700,
			Drivetrain:     "AWD",
			WeightKG:       1717,
			FrontWeightPct: 52,
			BalanceBias:    50,
		})
		if err != nil {
			t.Fatalf("generate baseline: %v", err)
		}
		requireFloatPtrNear(t, "front camber", result.ProfileDraft.FrontCamber, -1.4, 0.001)
		requireFloatPtrNear(t, "rear camber", result.ProfileDraft.RearCamber, -1.1, 0.001)
		requireFloatPtrNear(t, "front toe", result.ProfileDraft.FrontToe, 0, 0.001)
		requireFloatPtrNear(t, "rear toe", result.ProfileDraft.RearToe, 0.1, 0.001)
	})
	t.Run("drivetrain modifiers", func(t *testing.T) {
		fwd, err := GenerateRoadStaticTuneBaseline(RoadStaticTuneBaselineInput{
			PI:             600,
			Drivetrain:     "FWD",
			WeightKG:       1300,
			FrontWeightPct: 50,
		})
		if err != nil {
			t.Fatalf("generate FWD baseline: %v", err)
		}
		rwd, err := GenerateRoadStaticTuneBaseline(RoadStaticTuneBaselineInput{
			PI:             600,
			Drivetrain:     "RWD",
			WeightKG:       1300,
			FrontWeightPct: 50,
		})
		if err != nil {
			t.Fatalf("generate RWD baseline: %v", err)
		}
		awd, err := GenerateRoadStaticTuneBaseline(RoadStaticTuneBaselineInput{
			PI:             600,
			Drivetrain:     "AWD",
			WeightKG:       1300,
			FrontWeightPct: 50,
		})
		if err != nil {
			t.Fatalf("generate AWD baseline: %v", err)
		}
		requireFloatPtrNear(t, "FWD front camber", fwd.ProfileDraft.FrontCamber, -1.3, 0.001)
		requireFloatPtrNear(t, "FWD rear camber", fwd.ProfileDraft.RearCamber, -0.7, 0.001)
		requireFloatPtrNear(t, "FWD caster", fwd.ProfileDraft.Caster, 5.2, 0.001)
		requireFloatPtrNear(t, "RWD rear camber", rwd.ProfileDraft.RearCamber, -0.9, 0.001)
		requireFloatPtrNear(t, "AWD front camber", awd.ProfileDraft.FrontCamber, -1.2, 0.001)
		requireFloatPtrNear(t, "AWD rear camber", awd.ProfileDraft.RearCamber, -0.8, 0.001)
	})
}

func TestGenerateRoadStaticTuneBaselineGearing(t *testing.T) {
	redline := 7200.0
	gearCount := int64(6)
	tireDiameter := 67.0
	targetTopSpeed := 280.0
	result, err := GenerateRoadStaticTuneBaseline(RoadStaticTuneBaselineInput{
		CarName:           "Gear Starter",
		PI:                800,
		Drivetrain:        "RWD",
		WeightKG:          1350,
		FrontWeightPct:    53,
		RedlineRPM:        &redline,
		GearCount:         &gearCount,
		TireDiameterCm:    &tireDiameter,
		TargetTopSpeedKmh: &targetTopSpeed,
	})
	if err != nil {
		t.Fatalf("generate baseline: %v", err)
	}
	if result.ProfileDraft.RedlineRPM == nil || *result.ProfileDraft.RedlineRPM != redline {
		t.Fatalf("redline = %#v", result.ProfileDraft.RedlineRPM)
	}
	if result.ProfileDraft.FinalDrive == nil || result.ProfileDraft.Gear1 == nil || result.ProfileDraft.Gear6 == nil {
		t.Fatalf("generated gearing incomplete: %#v", result.ProfileDraft)
	}
	if result.ProfileDraft.Gear7 != nil || hasGeneratedField(result.GeneratedFields, "gear7") {
		t.Fatalf("gear7 should be treated as locked: %#v", result.ProfileDraft.Gear7)
	}
	gotTopSpeed := gearTopSpeedKmh(redline, tireDiameter, *result.ProfileDraft.FinalDrive, *result.ProfileDraft.Gear6)
	if diff := absFloat(gotTopSpeed - targetTopSpeed); diff > 3 {
		t.Fatalf("top speed %.2f km/h differs from target %.2f by %.2f", gotTopSpeed, targetTopSpeed, diff)
	}
	retentions := []float64{
		gearShiftRetention(result.ProfileDraft.Gear1, result.ProfileDraft.Gear2),
		gearShiftRetention(result.ProfileDraft.Gear2, result.ProfileDraft.Gear3),
		gearShiftRetention(result.ProfileDraft.Gear3, result.ProfileDraft.Gear4),
		gearShiftRetention(result.ProfileDraft.Gear4, result.ProfileDraft.Gear5),
		gearShiftRetention(result.ProfileDraft.Gear5, result.ProfileDraft.Gear6),
	}
	for i := 1; i < len(retentions); i++ {
		if retentions[i] <= retentions[i-1] {
			t.Fatalf("6-speed retention should increase from low to high gears: %#v", retentions)
		}
	}
	if retentions[0] < 0.70 || retentions[0] > 0.73 {
		t.Fatalf("1->2 retention = %.3f, want about 71-73%%", retentions[0])
	}
	if retentions[len(retentions)-1] < 0.85 || retentions[len(retentions)-1] > 0.87 {
		t.Fatalf("5->6 retention = %.3f, want about 85-86%%", retentions[len(retentions)-1])
	}
}

func TestGenerateRoadStaticTuneBaselineDriftGearingTargetsCoreGearRPM(t *testing.T) {
	redline := 7200.0
	gearCount := int64(6)
	tireDiameter := 67.0
	targetDriftSpeed := 100.0
	result, err := GenerateRoadStaticTuneBaseline(RoadStaticTuneBaselineInput{
		UseCase:           "Drift",
		PI:                700,
		Drivetrain:        "RWD",
		WeightKG:          1400,
		FrontWeightPct:    54,
		RedlineRPM:        &redline,
		GearCount:         &gearCount,
		TireDiameterCm:    &tireDiameter,
		TargetTopSpeedKmh: &targetDriftSpeed,
	})
	if err != nil {
		t.Fatalf("generate drift baseline: %v", err)
	}
	if result.ProfileDraft.FinalDrive == nil || result.ProfileDraft.Gear3 == nil {
		t.Fatalf("drift gearing incomplete: %#v", result.ProfileDraft)
	}
	coreRPMRatio := gearRPMAtSpeedKmh(targetDriftSpeed, tireDiameter, *result.ProfileDraft.FinalDrive, *result.ProfileDraft.Gear3) / redline
	if coreRPMRatio < 0.85 || coreRPMRatio > 0.92 {
		t.Fatalf("drift core gear rpm ratio = %.3f, want 0.85-0.92", coreRPMRatio)
	}
	if retention := gearShiftRetention(result.ProfileDraft.Gear1, result.ProfileDraft.Gear2); retention < 0.67 || retention > 0.69 {
		t.Fatalf("drift 1->2 retention = %.3f, want about 68%%", retention)
	}
	if retention := gearShiftRetention(result.ProfileDraft.Gear3, result.ProfileDraft.Gear4); retention < 0.81 || retention > 0.83 {
		t.Fatalf("drift 3->4 retention = %.3f, want about 82%%", retention)
	}
}

func TestGenerateRoadStaticTuneBaselineDriftSpeedBiasOnlyAdjustsFinalDrive(t *testing.T) {
	redline := 7200.0
	gearCount := int64(6)
	tireDiameter := 67.0
	targetDriftSpeed := 100.0
	base := RoadStaticTuneBaselineInput{
		UseCase:           "Drift",
		PI:                700,
		Drivetrain:        "RWD",
		WeightKG:          1400,
		FrontWeightPct:    54,
		RedlineRPM:        &redline,
		GearCount:         &gearCount,
		TireDiameterCm:    &tireDiameter,
		TargetTopSpeedKmh: &targetDriftSpeed,
	}
	neutral, err := GenerateRoadStaticTuneBaseline(base)
	if err != nil {
		t.Fatalf("generate neutral drift baseline: %v", err)
	}
	base.SpeedBias = 150
	accel, err := GenerateRoadStaticTuneBaseline(base)
	if err != nil {
		t.Fatalf("generate accel drift baseline: %v", err)
	}
	base.SpeedBias = 50
	faster, err := GenerateRoadStaticTuneBaseline(base)
	if err != nil {
		t.Fatalf("generate faster drift baseline: %v", err)
	}
	if !(faster.ProfileDraft.FinalDrive != nil && neutral.ProfileDraft.FinalDrive != nil && accel.ProfileDraft.FinalDrive != nil) {
		t.Fatalf("expected generated final drive values")
	}
	if !(*faster.ProfileDraft.FinalDrive < *neutral.ProfileDraft.FinalDrive && *accel.ProfileDraft.FinalDrive > *neutral.ProfileDraft.FinalDrive) {
		t.Fatalf("drift speed bias should lower/raise final drive: fast %.2f neutral %.2f accel %.2f",
			*faster.ProfileDraft.FinalDrive, *neutral.ProfileDraft.FinalDrive, *accel.ProfileDraft.FinalDrive)
	}
	requireFloatPtrNear(t, "drift gear1 unchanged", accel.ProfileDraft.Gear1, *neutral.ProfileDraft.Gear1, 0.001)
	requireFloatPtrNear(t, "drift gear3 unchanged", accel.ProfileDraft.Gear3, *neutral.ProfileDraft.Gear3, 0.001)
	requireFloatPtrNear(t, "drift gear6 unchanged", faster.ProfileDraft.Gear6, *neutral.ProfileDraft.Gear6, 0.001)
}

func TestGenerateRoadStaticTuneBaselineTenSpeedGearingUsesWiderLowGears(t *testing.T) {
	redline := 8000.0
	gearCount := int64(10)
	tireDiameter := 67.0
	targetTopSpeed := 360.0
	result, err := GenerateRoadStaticTuneBaseline(RoadStaticTuneBaselineInput{
		CarName:           "Ten Speed Gear Starter",
		PI:                900,
		Drivetrain:        "AWD",
		WeightKG:          1450,
		FrontWeightPct:    51,
		RedlineRPM:        &redline,
		GearCount:         &gearCount,
		TireDiameterCm:    &tireDiameter,
		TargetTopSpeedKmh: &targetTopSpeed,
	})
	if err != nil {
		t.Fatalf("generate baseline: %v", err)
	}
	if result.ProfileDraft.Gear10 == nil {
		t.Fatalf("gear10 should be generated for 10-speed input: %#v", result.ProfileDraft)
	}
	firstRetention := gearShiftRetention(result.ProfileDraft.Gear1, result.ProfileDraft.Gear2)
	lastRetention := gearShiftRetention(result.ProfileDraft.Gear9, result.ProfileDraft.Gear10)
	if firstRetention < 0.74 || firstRetention > 0.76 {
		t.Fatalf("10-speed 1->2 retention = %.3f, want about 75%%", firstRetention)
	}
	if lastRetention < 0.91 || lastRetention > 0.93 {
		t.Fatalf("10-speed 9->10 retention = %.3f, want about 92%%", lastRetention)
	}
	if lastRetention-firstRetention < 0.12 {
		t.Fatalf("10-speed low gears should be meaningfully wider than high gears: first %.3f last %.3f", firstRetention, lastRetention)
	}
}

func TestGenerateRoadStaticTuneBaselineSkipsIncompleteGearing(t *testing.T) {
	redline := 7200.0
	result, err := GenerateRoadStaticTuneBaseline(RoadStaticTuneBaselineInput{
		CarName:        "Gear Starter",
		PI:             800,
		Drivetrain:     "RWD",
		WeightKG:       1350,
		FrontWeightPct: 53,
		RedlineRPM:     &redline,
	})
	if err != nil {
		t.Fatalf("generate baseline: %v", err)
	}
	if result.ProfileDraft.FinalDrive != nil || result.ProfileDraft.Gear1 != nil {
		t.Fatalf("incomplete gearing input should not generate ratios: %#v", result.ProfileDraft)
	}
	if !hasGeneratedField(result.GeneratedFields, "frontSpring") || !hasSkippedField(result.SkippedFields, "finalDrive") {
		t.Fatalf("baseline fields should continue while gearing is skipped: generated=%#v skipped=%#v", result.GeneratedFields, result.SkippedFields)
	}
}

func TestGenerateRoadStaticTuneBaselineTierRecommendations(t *testing.T) {
	result, err := GenerateRoadStaticTuneBaseline(RoadStaticTuneBaselineInput{
		CarName:                   "Road Starter",
		PI:                        700,
		Drivetrain:                "RWD",
		WeightKG:                  1300,
		FrontWeightPct:            52,
		FrontRideHeightAdjustable: true,
		RearRideHeightAdjustable:  true,
		FrontAeroAdjustable:       true,
	})
	if err != nil {
		t.Fatalf("generate baseline: %v", err)
	}
	if !hasTierRecommendation(result.TierRecommendations, "frontRideHeight", "low", true) {
		t.Fatalf("missing front ride tier: %#v", result.TierRecommendations)
	}
	if !hasTierRecommendation(result.TierRecommendations, "frontAero", "medium", true) {
		t.Fatalf("missing front aero tier: %#v", result.TierRecommendations)
	}
	if !hasTierRecommendation(result.TierRecommendations, "rearAero", "medium", false) || !hasSkippedField(result.SkippedFields, "rearAero") {
		t.Fatalf("missing non-adjustable rear aero state: tiers=%#v skipped=%#v", result.TierRecommendations, result.SkippedFields)
	}
	if hasGeneratedField(result.GeneratedFields, "frontRideHeight") || hasGeneratedField(result.GeneratedFields, "frontAero") {
		t.Fatalf("tier-only fields should not be generated: %#v", result.GeneratedFields)
	}
}

func TestGenerateRoadStaticTuneBaselineDrivetrainDiffs(t *testing.T) {
	cases := []struct {
		name       string
		drive      string
		wantFront  bool
		wantRear   bool
		wantCenter bool
		frontARB   float64
		rearARB    float64
		frontAccel *float64
		frontDecel *float64
		rearAccel  *float64
		rearDecel  *float64
	}{
		{name: "FWD", drive: "FWD", wantFront: true, frontARB: 12, rearARB: 35, frontAccel: testFloatPtr(27), frontDecel: testFloatPtr(7)},
		{name: "RWD", drive: "RWD", wantRear: true, frontARB: 22, rearARB: 31, rearAccel: testFloatPtr(50), rearDecel: testFloatPtr(31)},
		{name: "AWD", drive: "AWD", wantFront: true, wantRear: true, wantCenter: true, frontARB: 26, rearARB: 35, frontAccel: testFloatPtr(22), frontDecel: testFloatPtr(7), rearAccel: testFloatPtr(56), rearDecel: testFloatPtr(27)},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := GenerateRoadStaticTuneBaseline(RoadStaticTuneBaselineInput{
				CarName:        tc.name,
				PI:             800,
				Drivetrain:     tc.drive,
				WeightKG:       1300,
				FrontWeightPct: 52,
			})
			if err != nil {
				t.Fatalf("generate baseline: %v", err)
			}
			if (result.ProfileDraft.FrontDiffAccel != nil) != tc.wantFront {
				t.Fatalf("front diff presence = %v", result.ProfileDraft.FrontDiffAccel != nil)
			}
			if (result.ProfileDraft.RearDiffAccel != nil) != tc.wantRear {
				t.Fatalf("rear diff presence = %v", result.ProfileDraft.RearDiffAccel != nil)
			}
			if (result.ProfileDraft.CenterDiffBalance != nil) != tc.wantCenter {
				t.Fatalf("center diff presence = %v", result.ProfileDraft.CenterDiffBalance != nil)
			}
			requireFloatPtrNear(t, "front arb", result.ProfileDraft.FrontARB, tc.frontARB, 0.001)
			requireFloatPtrNear(t, "rear arb", result.ProfileDraft.RearARB, tc.rearARB, 0.001)
			if tc.frontAccel != nil {
				requireFloatPtrNear(t, "front diff accel", result.ProfileDraft.FrontDiffAccel, *tc.frontAccel, 0.001)
			}
			if tc.frontDecel != nil {
				requireFloatPtrNear(t, "front diff decel", result.ProfileDraft.FrontDiffDecel, *tc.frontDecel, 0.001)
			}
			if tc.rearAccel != nil {
				requireFloatPtrNear(t, "rear diff accel", result.ProfileDraft.RearDiffAccel, *tc.rearAccel, 0.001)
			}
			if tc.rearDecel != nil {
				requireFloatPtrNear(t, "rear diff decel", result.ProfileDraft.RearDiffDecel, *tc.rearDecel, 0.001)
			}
		})
	}
}

func TestApplyRoadStaticTuneBaselineCreateAndUpdateSnapshot(t *testing.T) {
	store := openTestStore(t)
	input := RoadStaticTuneBaselineInput{
		PI:             700,
		Drivetrain:     "RWD",
		WeightKG:       1350,
		FrontWeightPct: 53,
	}
	created, err := store.ApplyRoadStaticTuneBaseline(RoadStaticTuneBaselineApplyInput{
		CreateNew:     true,
		BaselineInput: input,
		SelectedFieldKeys: []string{
			"frontTirePressure",
			"rearTirePressure",
			"frontArb",
		},
	})
	if err != nil {
		t.Fatalf("create from baseline: %v", err)
	}
	if created.Profile.ID == 0 || created.Profile.FrontTirePressure == nil || created.Profile.RearTirePressure == nil || created.Profile.FrontARB == nil {
		t.Fatalf("created profile = %#v", created.Profile)
	}
	if created.Profile.CarName != "Quick Tune A700 RWD" {
		t.Fatalf("auto profile name = %q", created.Profile.CarName)
	}
	if created.Profile.RearARB != nil {
		t.Fatalf("unselected rear ARB should remain empty: %#v", created.Profile.RearARB)
	}

	updated, err := store.ApplyRoadStaticTuneBaseline(RoadStaticTuneBaselineApplyInput{
		TargetProfileID: created.Profile.ID,
		BaselineInput:   input,
		SelectedFieldKeys: []string{
			"rearArb",
		},
	})
	if err != nil {
		t.Fatalf("update from baseline: %v", err)
	}
	if updated.Profile.RearARB == nil {
		t.Fatalf("updated profile missing rear ARB: %#v", updated.Profile)
	}
	snapshots, err := store.ListTuneProfileSnapshots(created.Profile.ID)
	if err != nil {
		t.Fatalf("list snapshots: %v", err)
	}
	if len(snapshots) != 1 || snapshots[0].ChangeReason != "road_static_baseline" {
		t.Fatalf("snapshots = %#v", snapshots)
	}
}

func TestApplyRoadStaticTuneBaselineRejectsForgedField(t *testing.T) {
	store := openTestStore(t)
	_, err := store.ApplyRoadStaticTuneBaseline(RoadStaticTuneBaselineApplyInput{
		CreateNew: true,
		BaselineInput: RoadStaticTuneBaselineInput{
			CarName:        "Generated Road",
			PI:             700,
			Drivetrain:     "RWD",
			WeightKG:       1350,
			FrontWeightPct: 53,
		},
		SelectedFieldKeys: []string{"gear1"},
	})
	if err == nil {
		t.Fatal("expected forged/non-generated field to fail")
	}
}

func TestApplyRoadStaticTuneBaselineWritesSelectedGearing(t *testing.T) {
	store := openTestStore(t)
	redline := 7200.0
	gearCount := int64(6)
	tireDiameter := 67.0
	targetTopSpeed := 280.0
	created, err := store.ApplyRoadStaticTuneBaseline(RoadStaticTuneBaselineApplyInput{
		CreateNew: true,
		BaselineInput: RoadStaticTuneBaselineInput{
			CarName:           "Generated Gear Road",
			PI:                800,
			Drivetrain:        "RWD",
			WeightKG:          1350,
			FrontWeightPct:    53,
			RedlineRPM:        &redline,
			GearCount:         &gearCount,
			TireDiameterCm:    &tireDiameter,
			TargetTopSpeedKmh: &targetTopSpeed,
		},
		SelectedFieldKeys: []string{"finalDrive", "gear1", "gear6", "redlineRPM"},
	})
	if err != nil {
		t.Fatalf("create with gearing: %v", err)
	}
	if created.Profile.FinalDrive == nil || created.Profile.Gear1 == nil || created.Profile.Gear6 == nil || created.Profile.RedlineRPM == nil {
		t.Fatalf("created profile missing selected gearing: %#v", created.Profile)
	}
	if created.Profile.Gear7 != nil {
		t.Fatalf("unselected/locked gear7 should remain empty: %#v", created.Profile.Gear7)
	}
}

func TestApplyRoadStaticTuneBaselineRejectsLockedGear(t *testing.T) {
	store := openTestStore(t)
	redline := 7200.0
	gearCount := int64(6)
	tireDiameter := 67.0
	targetTopSpeed := 280.0
	_, err := store.ApplyRoadStaticTuneBaseline(RoadStaticTuneBaselineApplyInput{
		CreateNew: true,
		BaselineInput: RoadStaticTuneBaselineInput{
			CarName:           "Generated Gear Road",
			PI:                800,
			Drivetrain:        "RWD",
			WeightKG:          1350,
			FrontWeightPct:    53,
			RedlineRPM:        &redline,
			GearCount:         &gearCount,
			TireDiameterCm:    &tireDiameter,
			TargetTopSpeedKmh: &targetTopSpeed,
		},
		SelectedFieldKeys: []string{"gear7"},
	})
	if err == nil {
		t.Fatal("expected locked gear to fail")
	}
}

func TestApplyRoadStaticTuneBaselineRejectsTierOnlyField(t *testing.T) {
	store := openTestStore(t)
	_, err := store.ApplyRoadStaticTuneBaseline(RoadStaticTuneBaselineApplyInput{
		CreateNew: true,
		BaselineInput: RoadStaticTuneBaselineInput{
			CarName:                   "Generated Road",
			PI:                        700,
			Drivetrain:                "RWD",
			WeightKG:                  1350,
			FrontWeightPct:            53,
			FrontRideHeightAdjustable: true,
		},
		SelectedFieldKeys: []string{"frontRideHeight"},
	})
	if err == nil {
		t.Fatal("expected tier-only field to fail")
	}
}

func hasGeneratedField(fields []BaselineGeneratedField, key string) bool {
	for _, field := range fields {
		if field.FieldKey == key {
			return true
		}
	}
	return false
}

func hasSkippedField(fields []BaselineSkippedField, key string) bool {
	for _, field := range fields {
		if field.FieldKey == key {
			return true
		}
	}
	return false
}

func hasTierRecommendation(fields []BaselineTierRecommendation, key string, tier string, applicable bool) bool {
	for _, field := range fields {
		if field.FieldKey == key && field.Tier == tier && field.Applicable == applicable {
			return true
		}
	}
	return false
}

func requireFloatPtrNear(t *testing.T, label string, got *float64, want float64, tolerance float64) {
	t.Helper()
	if got == nil {
		t.Fatalf("%s is nil, want %.3f", label, want)
	}
	if diff := absFloat(*got - want); diff > tolerance {
		t.Fatalf("%s = %.3f, want %.3f ± %.3f", label, *got, want, tolerance)
	}
}

func testFloatPtr(value float64) *float64 {
	return &value
}

func isStepValue(value float64, step float64) bool {
	steps := value / step
	return absFloat(steps-math.Round(steps)) < 0.000001
}

func gearTopSpeedKmh(redlineRPM float64, tireDiameterCm float64, finalDrive float64, gearRatio float64) float64 {
	circumferenceMeters := 3.141592653589793 * tireDiameterCm / 100
	wheelRPM := redlineRPM / (finalDrive * gearRatio)
	return wheelRPM * circumferenceMeters * 60 / 1000
}

func gearRPMAtSpeedKmh(speedKmh float64, tireDiameterCm float64, finalDrive float64, gearRatio float64) float64 {
	circumferenceMeters := 3.141592653589793 * tireDiameterCm / 100
	wheelRPM := speedKmh * 1000 / 60 / circumferenceMeters
	return wheelRPM * finalDrive * gearRatio
}

func gearShiftRetention(lowerGear *float64, higherGear *float64) float64 {
	if lowerGear == nil || higherGear == nil || *lowerGear == 0 {
		return 0
	}
	return *higherGear / *lowerGear
}
