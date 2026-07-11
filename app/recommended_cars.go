package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"fh6worker/internal/storage"
)

type RecommendedCarsFileResult struct {
	Path  string `json:"path"`
	Count int    `json:"count"`
}

type RecommendedCarsFileSelection struct {
	Path      string   `json:"path"`
	Exists    bool     `json:"exists"`
	Version   string   `json:"version"`
	IDs       []string `json:"ids"`
	TuneCodes []string `json:"tuneCodes"`
	Count     int      `json:"count"`
}

type recommendedCarsExportedCar struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	UseCase      string   `json:"useCase"`
	CarClass     string   `json:"carClass"`
	PI           int      `json:"pi"`
	Drivetrain   string   `json:"drivetrain"`
	TireCompound string   `json:"tireCompound"`
	TuneCode     string   `json:"tuneCode,omitempty"`
	TuneCodes    []string `json:"tuneCodes"`
	ImageSrc     string   `json:"imageSrc,omitempty"`
	Tags         []string `json:"tags"`
	Reason       string   `json:"reason"`
}

func (a *App) ListRecommendedCars() ([]storage.RecommendedCar, error) {
	if err := a.ensureStore(); err != nil {
		return nil, err
	}
	return a.store.ListRecommendedCars()
}

func (a *App) SaveRecommendedCar(input storage.RecommendedCarInput) (*storage.RecommendedCar, error) {
	if err := a.ensureStore(); err != nil {
		return nil, err
	}
	return a.store.SaveRecommendedCar(input)
}

func (a *App) SaveRecommendedCarRecord(input storage.RecommendedCarInput, previousID string) (*storage.RecommendedCar, error) {
	if err := a.ensureStore(); err != nil {
		return nil, err
	}
	return a.store.SaveRecommendedCarRecord(input, previousID)
}

func (a *App) DeleteRecommendedCar(id string) error {
	if err := a.ensureStore(); err != nil {
		return err
	}
	return a.store.DeleteRecommendedCar(id)
}

func (a *App) DeleteAllRecommendedCars() (int64, error) {
	if err := a.ensureStore(); err != nil {
		return 0, err
	}
	return a.store.DeleteAllRecommendedCars()
}

func (a *App) ExportRecommendedCarsFile(version string) (*RecommendedCarsFileResult, error) {
	if err := a.ensureStore(); err != nil {
		return nil, err
	}
	cars, err := a.store.ListRecommendedCars()
	if err != nil {
		return nil, err
	}
	inputs := make([]storage.RecommendedCarInput, 0, len(cars))
	for _, car := range cars {
		inputs = append(inputs, car.RecommendedCarInput)
	}
	return writeRecommendedCarsFile(inputs, version)
}

func (a *App) SaveRecommendedCarsFile(cars []storage.RecommendedCarInput, version string) (*RecommendedCarsFileResult, error) {
	return writeRecommendedCarsFile(cars, version)
}

func (a *App) LoadRecommendedCarsFileSelection() (*RecommendedCarsFileSelection, error) {
	target, err := recommendedCarsFilePath()
	if err != nil {
		return nil, err
	}
	raw, err := os.ReadFile(target)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &RecommendedCarsFileSelection{
				Path:   target,
				Exists: false,
				IDs:    []string{},
			}, nil
		}
		return nil, err
	}
	return parseRecommendedCarsFileSelection(raw, target)
}

func writeRecommendedCarsFile(cars []storage.RecommendedCarInput, version string) (*RecommendedCarsFileResult, error) {
	if len(cars) == 0 {
		return nil, errors.New("recommended cars list is empty")
	}
	version = strings.TrimSpace(version)
	if version == "" {
		version = time.Now().Format("2006-01-02") + "-001"
	}
	normalized := make([]storage.RecommendedCarInput, 0, len(cars))
	seenIDs := map[string]bool{}
	for index, car := range cars {
		next, err := storage.NormalizeRecommendedCarInput(car)
		if err != nil {
			return nil, fmt.Errorf("car %d: %w", index+1, err)
		}
		if seenIDs[next.ID] {
			return nil, fmt.Errorf("car %d: duplicate id %q", index+1, next.ID)
		}
		seenIDs[next.ID] = true
		normalized = append(normalized, next)
	}
	target, err := recommendedCarsFilePath()
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return nil, err
	}
	if err := os.WriteFile(target, []byte(formatRecommendedCarsJSON(version, normalized)), 0o644); err != nil {
		return nil, err
	}
	return &RecommendedCarsFileResult{Path: target, Count: len(mergeRecommendedCarsForExport(normalized))}, nil
}

func recommendedCarsFilePath() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for dir := wd; ; dir = filepath.Dir(dir) {
		candidateDir := filepath.Join(dir, "weChatApp", "miniprogram", "data")
		if info, err := os.Stat(candidateDir); err == nil && info.IsDir() {
			return filepath.Join(candidateDir, "recommendedCars.json"), nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
	}
	return "", errors.New("weChatApp/miniprogram/data directory was not found")
}

func parseRecommendedCarsFileSelection(raw []byte, path string) (*RecommendedCarsFileSelection, error) {
	var payload struct {
		Version string `json:"version"`
		Cars    []struct {
			ID        string   `json:"id"`
			TuneCode  string   `json:"tuneCode"`
			TuneCodes []string `json:"tuneCodes"`
		} `json:"cars"`
	}
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil, fmt.Errorf("parse recommendedCars.json: %w", err)
	}
	ids := make([]string, 0, len(payload.Cars))
	tuneCodes := []string{}
	seenIDs := map[string]bool{}
	seenTuneCodes := map[string]bool{}
	for _, car := range payload.Cars {
		id := strings.TrimSpace(car.ID)
		if id != "" && !seenIDs[id] {
			seenIDs[id] = true
			ids = append(ids, id)
		}
		for _, code := range append([]string{car.TuneCode}, car.TuneCodes...) {
			code = storage.NormalizeTuneShareCode(code)
			if code == "" || seenTuneCodes[code] {
				continue
			}
			seenTuneCodes[code] = true
			tuneCodes = append(tuneCodes, code)
		}
	}
	return &RecommendedCarsFileSelection{
		Path:      path,
		Exists:    true,
		Version:   strings.TrimSpace(payload.Version),
		IDs:       ids,
		TuneCodes: tuneCodes,
		Count:     len(ids),
	}, nil
}

func formatRecommendedCarsJSON(version string, cars []storage.RecommendedCarInput) string {
	payload := struct {
		Version string                       `json:"version"`
		Cars    []recommendedCarsExportedCar `json:"cars"`
	}{
		Version: strings.TrimSpace(version),
		Cars:    mergeRecommendedCarsForExport(cars),
	}
	encoded, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return "{}\n"
	}
	return string(encoded) + "\n"
}

func mergeRecommendedCarsForExport(cars []storage.RecommendedCarInput) []recommendedCarsExportedCar {
	type aggregate struct {
		car       recommendedCarsExportedCar
		codeSet   map[string]bool
		tagSet    map[string]bool
		reasonSet map[string]bool
	}
	merged := []recommendedCarsExportedCar{}
	byKey := map[string]*aggregate{}
	for _, car := range cars {
		key := recommendedCarsMergeKey(car)
		entry := byKey[key]
		if entry == nil {
			entry = &aggregate{
				car: recommendedCarsExportedCar{
					ID:           recommendedCarsExportID(car),
					Name:         car.Name,
					UseCase:      car.UseCase,
					CarClass:     car.CarClass,
					PI:           car.PI,
					Drivetrain:   car.Drivetrain,
					TireCompound: car.TireCompound,
					TuneCodes:    []string{},
					ImageSrc:     car.ImageSrc,
					Tags:         []string{},
					Reason:       strings.TrimSpace(car.Reason),
				},
				codeSet:   map[string]bool{},
				tagSet:    map[string]bool{},
				reasonSet: map[string]bool{},
			}
			if entry.car.Reason != "" {
				entry.reasonSet[entry.car.Reason] = true
			}
			byKey[key] = entry
			merged = append(merged, entry.car)
		}
		if car.PI > entry.car.PI {
			entry.car.PI = car.PI
			entry.car.ID = recommendedCarsExportID(car)
		}
		if entry.car.ImageSrc == "" && car.ImageSrc != "" {
			entry.car.ImageSrc = car.ImageSrc
		}
		if entry.car.Drivetrain == "" && car.Drivetrain != "" {
			entry.car.Drivetrain = car.Drivetrain
		}
		if entry.car.TireCompound == "" && car.TireCompound != "" {
			entry.car.TireCompound = car.TireCompound
		}
		code := storage.NormalizeTuneShareCode(car.TuneCode)
		if code != "" && !entry.codeSet[code] {
			entry.codeSet[code] = true
			entry.car.TuneCodes = append(entry.car.TuneCodes, code)
		}
		for _, tag := range car.Tags {
			tag = strings.TrimSpace(tag)
			if tag == "" || entry.tagSet[tag] {
				continue
			}
			entry.tagSet[tag] = true
			entry.car.Tags = append(entry.car.Tags, tag)
		}
		reason := strings.TrimSpace(car.Reason)
		if reason != "" && !entry.reasonSet[reason] {
			entry.reasonSet[reason] = true
			if entry.car.Reason == "" {
				entry.car.Reason = reason
			}
		}
		entry.car.TuneCode = ""
		if len(entry.car.TuneCodes) > 0 {
			entry.car.TuneCode = entry.car.TuneCodes[0]
		}
		for index := range merged {
			if recommendedCarsMergeKeyFromExport(merged[index]) == key {
				merged[index] = entry.car
				break
			}
		}
	}
	return merged
}

func recommendedCarsMergeKey(car storage.RecommendedCarInput) string {
	return strings.Join([]string{
		recommendedCarsNormalizeMergeText(car.Name),
		car.UseCase,
		car.CarClass,
	}, "|")
}

func recommendedCarsMergeKeyFromExport(car recommendedCarsExportedCar) string {
	return strings.Join([]string{
		recommendedCarsNormalizeMergeText(car.Name),
		car.UseCase,
		car.CarClass,
	}, "|")
}

func recommendedCarsExportID(car storage.RecommendedCarInput) string {
	withoutCode := car
	withoutCode.ID = ""
	withoutCode.TuneCode = ""
	return storage.GenerateRecommendedCarID(withoutCode)
}

func recommendedCarsNormalizeMergeText(value string) string {
	return strings.ToLower(strings.Join(strings.Fields(strings.TrimSpace(value)), " "))
}
