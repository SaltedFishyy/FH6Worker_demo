package main

import (
	"encoding/json"
	"testing"

	"fh6worker/internal/storage"
)

func TestFormatRecommendedCarsJSON(t *testing.T) {
	cars := []storage.RecommendedCarInput{{
		ID:                "porsche-911-gt3-2021-road-a800-413829605",
		Name:              "2021 Porsche 911 GT3",
		UseCase:           "Road",
		UseCaseLabel:      "road-label-db-only",
		PI:                800,
		CarClass:          "A",
		Drivetrain:        "AWD",
		TireCompound:      "sport",
		TireCompoundLabel: "sport-label-db-only",
		WeightKG:          1435,
		FrontWeightPct:    39,
		TuneCode:          "413829605",
		ImageSrc:          "https://example.com/car.png",
		Tags:              []string{"grip", "speed"},
		Reason:            "stable at speed",
	}}
	output := formatRecommendedCarsJSON("2026-05-28-001", cars)
	var payload struct {
		Version string `json:"version"`
		Cars    []struct {
			ID                string   `json:"id"`
			Name              string   `json:"name"`
			UseCase           string   `json:"useCase"`
			UseCaseLabel      string   `json:"useCaseLabel"`
			CarClass          string   `json:"carClass"`
			PI                int      `json:"pi"`
			Drivetrain        string   `json:"drivetrain"`
			TireCompound      string   `json:"tireCompound"`
			TireCompoundLabel string   `json:"tireCompoundLabel"`
			WeightKG          float64  `json:"weightKG"`
			TuneCode          string   `json:"tuneCode"`
			TuneCodes         []string `json:"tuneCodes"`
			ImageSrc          string   `json:"imageSrc"`
			Tags              []string `json:"tags"`
			Reason            string   `json:"reason"`
		} `json:"cars"`
	}
	if err := json.Unmarshal([]byte(output), &payload); err != nil {
		t.Fatalf("output is not valid JSON: %v\n%s", err, output)
	}
	if payload.Version != "2026-05-28-001" || len(payload.Cars) != 1 {
		t.Fatalf("payload header = %#v", payload)
	}
	car := payload.Cars[0]
	if car.ID != "porsche-911-gt3-2021-road-a800" || car.Name != "2021 Porsche 911 GT3" || car.PI != 800 {
		t.Fatalf("exported car = %#v", car)
	}
	if car.TuneCode != "413829605" || len(car.TuneCodes) != 1 || car.TuneCodes[0] != "413829605" {
		t.Fatalf("exported tune codes = %#v", car)
	}
	if car.UseCaseLabel != "" || car.TireCompoundLabel != "" || car.WeightKG != 0 {
		t.Fatalf("exported JSON leaked database-only fields: %#v", car)
	}
	if len(car.Tags) != 2 || car.Tags[0] != "grip" || car.Reason != "stable at speed" {
		t.Fatalf("exported tags/reason = %#v", car)
	}
}

func TestFormatRecommendedCarsJSONMergesSameCarUseCaseClass(t *testing.T) {
	cars := []storage.RecommendedCarInput{
		{
			ID:           "porsche-911-gt3-2021-road-a800-413829605",
			Name:         "2021 Porsche 911 GT3",
			UseCase:      "Road",
			PI:           800,
			CarClass:     "A",
			Drivetrain:   "AWD",
			TireCompound: "sport",
			TuneCode:     "413829605",
			Tags:         []string{"grip"},
			Reason:       "one",
		},
		{
			ID:           "porsche-911-gt3-2021-road-a800-123456789",
			Name:         "2021 Porsche 911 GT3",
			UseCase:      "Road",
			PI:           799,
			CarClass:     "A",
			Drivetrain:   "AWD",
			TireCompound: "sport",
			TuneCode:     "123456789",
			Tags:         []string{"speed"},
			Reason:       "two",
		},
		{
			ID:           "porsche-911-gt3-2021-road-a800-987654321",
			Name:         "2021 Porsche 911 GT3",
			UseCase:      "Road",
			PI:           800,
			CarClass:     "A",
			Drivetrain:   "AWD",
			TireCompound: "sport",
			TuneCode:     "987654321",
		},
	}
	output := formatRecommendedCarsJSON("2026-05-28-001", cars)
	var payload struct {
		Cars []struct {
			ID        string   `json:"id"`
			PI        int      `json:"pi"`
			TuneCode  string   `json:"tuneCode"`
			TuneCodes []string `json:"tuneCodes"`
			Tags      []string `json:"tags"`
		} `json:"cars"`
	}
	if err := json.Unmarshal([]byte(output), &payload); err != nil {
		t.Fatalf("output is not valid JSON: %v\n%s", err, output)
	}
	if len(payload.Cars) != 1 {
		t.Fatalf("exported cars = %#v", payload.Cars)
	}
	car := payload.Cars[0]
	if car.ID != "porsche-911-gt3-2021-road-a800" || car.PI != 800 {
		t.Fatalf("merged identity = %#v", car)
	}
	wantCodes := []string{"413829605", "123456789", "987654321"}
	if car.TuneCode != wantCodes[0] || len(car.TuneCodes) != len(wantCodes) {
		t.Fatalf("merged tune codes = %#v", car)
	}
	for index, want := range wantCodes {
		if car.TuneCodes[index] != want {
			t.Fatalf("merged tune codes = %#v", car.TuneCodes)
		}
	}
	if len(car.Tags) != 2 {
		t.Fatalf("merged tags = %#v", car.Tags)
	}
}

func TestParseRecommendedCarsFileSelection(t *testing.T) {
	raw := []byte(`{
  "version": "2026-05-28-001",
  "cars": [
    {"id": "porsche-911-gt3-2021-road-a800", "tuneCodes": ["413 829 605", "123456789"]},
    {"id": " "},
    {"id": "porsche-911-gt3-2021-road-a800"},
    {"id": "honda-civic-type-r-2018-road-b700", "tuneCode": "987 654 321"}
  ]
}`)
	selection, err := parseRecommendedCarsFileSelection(raw, `D:\FH6Worker\weChatApp\miniprogram\data\recommendedCars.json`)
	if err != nil {
		t.Fatalf("parse selection: %v", err)
	}
	if !selection.Exists || selection.Version != "2026-05-28-001" || selection.Count != 2 {
		t.Fatalf("selection header = %#v", selection)
	}
	if len(selection.IDs) != 2 || selection.IDs[0] != "porsche-911-gt3-2021-road-a800" || selection.IDs[1] != "honda-civic-type-r-2018-road-b700" {
		t.Fatalf("selection IDs = %#v", selection.IDs)
	}
	if len(selection.TuneCodes) != 3 || selection.TuneCodes[0] != "413829605" || selection.TuneCodes[1] != "123456789" || selection.TuneCodes[2] != "987654321" {
		t.Fatalf("selection tune codes = %#v", selection.TuneCodes)
	}
}

func TestNormalizeRecommendedCarInput(t *testing.T) {
	car, err := storage.NormalizeRecommendedCarInput(storage.RecommendedCarInput{
		Name:         "2018 Honda Civic Type R",
		UseCase:      "Road",
		PI:           700,
		CarClass:     " b ",
		Drivetrain:   " fwd ",
		TireCompound: "sport",
		TuneCode:     "927164038",
		Tags:         []string{"front", "", " stable "},
	})
	if err != nil {
		t.Fatalf("normalize input: %v", err)
	}
	if car.ID != "honda-civic-type-r-2018-road-b700-927164038" || car.CarClass != "B" || car.Drivetrain != "FWD" {
		t.Fatalf("normalized car = %#v", car)
	}
	if car.UseCaseLabel == "" || car.TireCompoundLabel == "" {
		t.Fatalf("linked labels = %q / %q", car.UseCaseLabel, car.TireCompoundLabel)
	}
	if len(car.Tags) != 2 || car.Tags[1] != "stable" {
		t.Fatalf("normalized tags = %#v", car.Tags)
	}
}
