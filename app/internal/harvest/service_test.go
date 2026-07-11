package harvest

import (
	"context"
	"encoding/json"
	"path/filepath"
	"testing"

	"fh6worker/internal/storage"
)

func TestMatchCarHandlesAliasesAndAccents(t *testing.T) {
	carInput, err := storage.NormalizeFH6CarInput(storage.FH6CarInput{
		Year:      2020,
		Make:      "BMW",
		Model:     "M2 Competition Coupé",
		Source:    SourceCODMunity,
		SourceRef: "321",
	})
	if err != nil {
		t.Fatalf("normalize car: %v", err)
	}
	match := MatchCar(storage.TuneHarvestCandidateInput{
		Source:  SourceJSRSheet,
		Year:    2020,
		Make:    "BMW",
		Model:   "M2 Comp",
		CarName: "BMW M2 Comp",
	}, []storage.FH6Car{{FH6CarInput: carInput}})
	if match.CarID == "" || match.Score < 0.88 {
		t.Fatalf("match = %+v, want confident BMW M2 Competition match", match)
	}
}

func TestCODMunityStateExtractionFindsCarTunings(t *testing.T) {
	state := map[string]any{
		"nested": map[string]any{
			"cars": []any{
				map[string]any{
					"ID":           32.0,
					"CarName":      "Vulcan",
					"Manufacturer": "Aston Martin",
					"Year":         2016.0,
					"PI Value":     884.0,
					"stats":        map[string]any{"Drive": "RWD"},
					"Tunings": []any{
						map[string]any{
							"Sharecode":  "953 733 070",
							"Playstyle":  "Road",
							"Class":      "S2",
							"PI Value":   900.0,
							"Tuner":      "K1Z Jumpy",
							"Drivetrain": "AWD",
							"Comment":    "Slick Tires",
						},
					},
				},
			},
		},
	}
	encoded, _ := json.Marshal(state)
	body := `<script id="serverApp-state" type="application/json">` + string(encoded) + `</script>`
	script, err := extractScriptJSON(body, "serverApp-state")
	if err != nil {
		t.Fatalf("extract script: %v", err)
	}
	var decoded any
	if err := json.Unmarshal([]byte(script), &decoded); err != nil {
		t.Fatalf("decode state: %v", err)
	}
	cars := findCODMunityCars(decoded)
	if len(cars) != 1 {
		t.Fatalf("cars = %d, want 1", len(cars))
	}
	tuning := mapValue(sliceValue(cars[0]["Tunings"])[0])
	if got := storage.NormalizeTuneShareCode(stringValue(tuning["Sharecode"])); got != "953733070" {
		t.Fatalf("share code = %q", got)
	}
}

func TestForzaFireDetailRequiresTunerCodeField(t *testing.T) {
	link := forzaFireBuildLink{
		URL:        "https://www.forzafire.com/build/2020-bmw-m2-competition-coup-14142",
		Path:       "/build/2020-bmw-m2-competition-coup-14142",
		Search:     "2020 BMW M2 Competition Coupé Grip Race Tune A700 GhostXD",
		CarClass:   "A",
		RaceType:   "Road Racing",
		Difficulty: "Easy",
	}
	body := `<div class="builder__code"><span>Tuner Code (In Game)</span><input type="text" name="tuner_code" value="828106826" readonly /></div>
<script type="text/javascript"> var loadBuild = {"car":{"id":"748","name":"2020 BMW M2 Competition Coup\u00e9"},"url":"2020-bmw-m2-competition-coup-14142","name":"Grip Race Tune A700","tags":{"race_type":["Road Racing"],"difficulty":["Easy"]},"build":{"tires":{"Tire Compound":"Stock"}}}; </script>`
	candidate, car, ok := parseForzaFireDetail(body, link)
	if !ok {
		t.Fatal("parseForzaFireDetail returned !ok")
	}
	if candidate.ShareCode != "828106826" || candidate.UseCase != "Road" || candidate.PI != 700 || candidate.TireCompound != "stock" {
		t.Fatalf("candidate = %+v", candidate)
	}
	if car.Year != 2020 || car.Make != "BMW" || car.Model != "M2 Competition Coupé" {
		t.Fatalf("car = %+v", car)
	}

	bodyWithoutTunerCode := `<input type="text" name="build_link" value="www.forzafire.com/build/2020-bmw-m2-competition-coup-14142" />`
	if _, _, ok := parseForzaFireDetail(bodyWithoutTunerCode, link); ok {
		t.Fatal("detail without tuner_code should not produce a candidate")
	}
}

func TestLooksLikeNonTuneCode(t *testing.T) {
	if !looksLikeNonTuneCode("blueprint share code 123 456 789") {
		t.Fatal("blueprint context should be rejected")
	}
	if looksLikeNonTuneCode("road tune share code 123 456 789 slick tires") {
		t.Fatal("normal tune context should not be rejected")
	}
}

func TestDedupeCandidatesByShareCodeKeepsBestContext(t *testing.T) {
	candidates := dedupeCandidatesByShareCode([]storage.TuneHarvestCandidateInput{
		{
			Source:    SourceJSRSheet,
			RawKey:    "jsr:1",
			ShareCode: "123 456 789",
			Status:    storage.TuneHarvestCandidateRejected,
		},
		{
			Source:       SourceCODMunity,
			RawKey:       "codmunity:1",
			ShareCode:    "123-456-789",
			CarName:      "2020 BMW M2 Competition",
			MatchedCarID: "2020-bmw-m2-competition",
			MatchScore:   0.94,
			Status:       storage.TuneHarvestCandidatePending,
			Tuner:        "Ghost",
		},
	})
	if len(candidates) != 1 {
		t.Fatalf("candidates = %d, want 1", len(candidates))
	}
	if candidates[0].Source != SourceCODMunity || candidates[0].ShareCode != "123456789" || candidates[0].Tuner != "Ghost" {
		t.Fatalf("candidate = %#v", candidates[0])
	}
}

func TestRunCancelledBeforeCollectionMarksRunCancelled(t *testing.T) {
	store, err := storage.Open(filepath.Join(t.TempDir(), "harvest.db"))
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() { _ = store.Close() })

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	result, err := Run(ctx, store, storage.TuneHarvestOptions{
		Sources: []string{SourceCODMunity},
		DryRun:  false,
	})
	if err != nil {
		t.Fatalf("run cancelled harvest: %v", err)
	}
	if result.Run == nil || result.Run.Status != storage.TuneHarvestRunCanceled {
		t.Fatalf("run status = %#v, want cancelled", result.Run)
	}
	if result.Found != 0 || result.Saved != 0 {
		t.Fatalf("counts = found %d saved %d, want zero", result.Found, result.Saved)
	}
	if len(result.Warnings) != 1 || result.Warnings[0] != tuneHarvestCanceledMessage {
		t.Fatalf("warnings = %#v, want cancellation message", result.Warnings)
	}
}
