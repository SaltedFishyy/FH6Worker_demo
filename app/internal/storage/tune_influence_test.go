package storage

import "testing"

func TestTuneToTireInfluenceMapCoversCoreFields(t *testing.T) {
	influenceMap := GetTuneToTireInfluenceMap()
	if influenceMap.Version == "" {
		t.Fatal("expected version")
	}
	items := map[string]TuneFieldInfluence{}
	for _, item := range influenceMap.Items {
		items[item.FieldKey] = item
		if len(item.Scope) == 0 {
			t.Fatalf("%s missing scope", item.FieldKey)
		}
		if len(item.Phases) == 0 {
			t.Fatalf("%s missing phases", item.FieldKey)
		}
		if len(item.TireMetrics) == 0 {
			t.Fatalf("%s missing tire metrics", item.FieldKey)
		}
		if len(item.EvidenceKeys) == 0 {
			t.Fatalf("%s missing evidence keys", item.FieldKey)
		}
	}

	for _, key := range []string{
		"frontTirePressure", "rearTirePressure", "finalDrive", "gear1", "gear10",
		"frontCamber", "rearCamber", "frontToe", "rearToe", "caster",
		"frontArb", "rearArb", "frontSpring", "rearSpring", "frontRideHeight", "rearRideHeight",
		"frontRebound", "rearRebound", "frontBump", "rearBump", "frontAero", "rearAero", "aeroBalance",
		"brakeBalance", "brakePressure", "frontDiffAccel", "frontDiffDecel", "rearDiffAccel", "rearDiffDecel", "centerDiffBalance",
	} {
		if _, ok := items[key]; !ok {
			t.Fatalf("missing influence for %s", key)
		}
	}
}

func TestExplainTuneFieldInfluenceFrontARB(t *testing.T) {
	item, err := ExplainTuneFieldInfluence("frontArb")
	if err != nil {
		t.Fatal(err)
	}
	if item.InfluenceType != "indirect" {
		t.Fatalf("frontArb influence type = %s, want indirect", item.InfluenceType)
	}
	if !containsInfluenceString(item.Scope, "front_axle") {
		t.Fatalf("frontArb scope = %#v, want front_axle", item.Scope)
	}
	if !containsInfluenceString(item.TireMetrics, "combined_slip") || !containsInfluenceString(item.TireMetrics, "slip_angle") {
		t.Fatalf("frontArb metrics = %#v, want lateral grip evidence", item.TireMetrics)
	}
	if !containsInfluenceString(item.EvidenceKeys, "front_rear_slip_delta") {
		t.Fatalf("frontArb evidence = %#v, want load transfer/balance evidence", item.EvidenceKeys)
	}
}

func TestExplainTuneFieldInfluenceRearDiffAccel(t *testing.T) {
	item, err := ExplainTuneFieldInfluence("rearDiffAccel")
	if err != nil {
		t.Fatal(err)
	}
	if !containsInfluenceString(item.Scope, "driven_wheels") || !containsInfluenceString(item.Scope, "rear_axle") {
		t.Fatalf("rearDiffAccel scope = %#v, want driven rear wheels", item.Scope)
	}
	if !containsInfluenceString(item.Phases, "corner_exit") {
		t.Fatalf("rearDiffAccel phases = %#v, want corner_exit", item.Phases)
	}
	if !containsInfluenceString(item.TireMetrics, "slip_ratio") {
		t.Fatalf("rearDiffAccel metrics = %#v, want slip_ratio", item.TireMetrics)
	}
	if !containsInfluenceString(item.EvidenceKeys, "rear_slip_ratio_p90") {
		t.Fatalf("rearDiffAccel evidence = %#v, want rear_slip_ratio_p90", item.EvidenceKeys)
	}
}

func TestExplainTuneFieldInfluenceUnknown(t *testing.T) {
	if _, err := ExplainTuneFieldInfluence("notAField"); err == nil {
		t.Fatal("expected error for unknown field")
	}
}

func containsInfluenceString(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}
