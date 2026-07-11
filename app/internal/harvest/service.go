package harvest

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io"
	"math"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"

	"fh6worker/internal/storage"
)

const (
	SourceCODMunity = "codmunity"
	SourceForzaFire = "forzafire"
	SourceJSRSheet  = "jsr_chronic_sheet"

	forzaFireBuildsURL = "https://www.forzafire.com/builds"
	jsrSheetCSVURL     = "https://docs.google.com/spreadsheets/d/15VOzHQ1VRJMPPPhWeXOKGzLt-0d4UTyN7KyBWVZ5Djk/export?format=csv&gid=2065485225"

	tuneHarvestCanceledMessage = "harvest cancelled"
)

type codmunityEndpoint struct {
	Key string
	URL string
}

var codmunityEndpoints = []codmunityEndpoint{
	{Key: "road", URL: "https://codmunity.gg/forza/road"},
	{Key: "road-purist", URL: "https://codmunity.gg/forza/road-purist"},
	{Key: "dirt", URL: "https://codmunity.gg/forza/dirt"},
	{Key: "cross-country", URL: "https://codmunity.gg/forza/cross-country"},
	{Key: "drag", URL: "https://codmunity.gg/forza/drag"},
	{Key: "drift-rwd", URL: "https://codmunity.gg/forza/drift-rwd"},
	{Key: "drift-awd", URL: "https://codmunity.gg/forza/pr-stunts-drift-awd"},
	{Key: "touge-purist", URL: "https://codmunity.gg/forza/touge-purist"},
	{Key: "pr-stunts", URL: "https://codmunity.gg/forza/pr-stunts"},
}

type collectedSource struct {
	Cars       []storage.FH6CarInput
	Candidates []storage.TuneHarvestCandidateInput
	Warnings   []string
}

func Run(ctx context.Context, store *storage.Store, options storage.TuneHarvestOptions) (*storage.TuneHarvestRunResult, error) {
	if store == nil {
		return nil, errors.New("store is required")
	}
	options.Sources = normalizeSources(options.Sources)
	if options.LimitPerSource <= 0 {
		options.LimitPerSource = 80
	}
	if options.LimitPerSource > 500 {
		options.LimitPerSource = 500
	}

	var run *storage.TuneHarvestRun
	var err error
	if !options.DryRun {
		run, err = store.CreateTuneHarvestRun(storage.TuneHarvestRunInput{Sources: options.Sources, DryRun: options.DryRun})
		if err != nil {
			return nil, err
		}
	}

	client := &http.Client{Timeout: 30 * time.Second}
	var collected collectedSource
	for _, source := range options.Sources {
		if ctx.Err() != nil {
			return finishRunAfterCancel(store, run, len(collected.Candidates), 0, 0, 0, 0, collected.Warnings, nil)
		}
		next, err := collectSource(ctx, client, source, options.LimitPerSource)
		collected.Cars = append(collected.Cars, next.Cars...)
		collected.Candidates = append(collected.Candidates, next.Candidates...)
		collected.Warnings = append(collected.Warnings, next.Warnings...)
		if ctx.Err() != nil {
			return finishRunAfterCancel(store, run, len(collected.Candidates), 0, 0, 0, 0, collected.Warnings, nil)
		}
		if err != nil {
			collected.Warnings = append(collected.Warnings, fmt.Sprintf("%s: %v", source, err))
			continue
		}
	}

	existingCars, err := store.ListFH6Cars()
	if err != nil {
		return nil, err
	}
	carsByID := map[string]storage.FH6Car{}
	for _, car := range existingCars {
		carsByID[car.CarID] = car
	}
	for _, input := range collected.Cars {
		if ctx.Err() != nil {
			return finishRunAfterCancel(store, run, len(collected.Candidates), 0, 0, 0, 0, collected.Warnings, nil)
		}
		normalized, err := storage.NormalizeFH6CarInput(input)
		if err != nil {
			collected.Warnings = append(collected.Warnings, fmt.Sprintf("car library: %v", err))
			continue
		}
		carsByID[normalized.CarID] = storage.FH6Car{FH6CarInput: normalized}
	}
	cars := make([]storage.FH6Car, 0, len(carsByID))
	for _, car := range carsByID {
		cars = append(cars, car)
	}
	sort.Slice(cars, func(i, j int) bool {
		if cars[i].Make == cars[j].Make {
			return cars[i].Model < cars[j].Model
		}
		return cars[i].Make < cars[j].Make
	})

	if !options.DryRun {
		if ctx.Err() != nil {
			return finishRunAfterCancel(store, run, len(collected.Candidates), 0, 0, 0, 0, collected.Warnings, nil)
		}
		if _, err := store.SaveFH6Cars(collected.Cars); err != nil {
			return finishRunAfterError(store, run, err, len(collected.Candidates), collected.Warnings)
		}
	}

	enrichedCandidates := make([]storage.TuneHarvestCandidateInput, 0, len(collected.Candidates))
	for _, candidate := range collected.Candidates {
		if ctx.Err() != nil {
			return finishRunAfterCancel(store, run, len(enrichedCandidates), 0, 0, 0, 0, collected.Warnings, nil)
		}
		candidate = enrichCandidate(candidate, cars)
		if run != nil {
			candidate.RunID = run.ID
		}
		enrichedCandidates = append(enrichedCandidates, candidate)
	}
	enrichedCandidates = dedupeCandidatesByShareCode(enrichedCandidates)

	result := &storage.TuneHarvestRunResult{
		Run:        run,
		Candidates: make([]storage.TuneHarvestCandidate, 0, len(enrichedCandidates)),
		Found:      len(enrichedCandidates),
		Warnings:   collected.Warnings,
	}
	for _, candidate := range enrichedCandidates {
		if ctx.Err() != nil {
			return finishRunAfterCancel(store, run, result.Found, result.Saved, result.Rejected, result.Pending, result.Imported, result.Warnings, result.Candidates)
		}
		if candidate.Status == storage.TuneHarvestCandidateRejected {
			result.Rejected++
		} else {
			result.Pending++
		}
		if options.DryRun {
			result.Candidates = append(result.Candidates, storage.TuneHarvestCandidate{TuneHarvestCandidateInput: candidate})
			continue
		}
		saved, err := store.UpsertTuneHarvestCandidate(candidate)
		if err != nil {
			result.Warnings = append(result.Warnings, fmt.Sprintf("%s %s: %v", candidate.Source, candidate.ShareCode, err))
			continue
		}
		result.Saved++
		result.Candidates = append(result.Candidates, *saved)
	}

	if run != nil {
		status := storage.TuneHarvestRunComplete
		message := strings.Join(result.Warnings, "\n")
		run, err = store.FinishTuneHarvestRun(run.ID, status, message, result.Found, result.Saved, result.Rejected, result.Pending, result.Imported)
		if err != nil {
			return nil, err
		}
		result.Run = run
	}
	return result, nil
}

func finishRunAfterCancel(store *storage.Store, run *storage.TuneHarvestRun, found int, saved int, rejected int, pending int, imported int, warnings []string, candidates []storage.TuneHarvestCandidate) (*storage.TuneHarvestRunResult, error) {
	message := strings.TrimSpace(strings.Join(appendCancelWarning(warnings), "\n"))
	result := &storage.TuneHarvestRunResult{
		Run:        run,
		Candidates: candidates,
		Found:      found,
		Saved:      saved,
		Rejected:   rejected,
		Pending:    pending,
		Imported:   imported,
		Warnings:   appendCancelWarning(warnings),
	}
	if run == nil {
		return result, nil
	}
	finished, err := store.FinishTuneHarvestRun(run.ID, storage.TuneHarvestRunCanceled, message, found, saved, rejected, pending, imported)
	if err != nil {
		return nil, err
	}
	result.Run = finished
	return result, nil
}

func appendCancelWarning(warnings []string) []string {
	out := append([]string{}, warnings...)
	for _, warning := range out {
		if warning == tuneHarvestCanceledMessage {
			return out
		}
	}
	return append(out, tuneHarvestCanceledMessage)
}

func finishRunAfterError(store *storage.Store, run *storage.TuneHarvestRun, cause error, found int, warnings []string) (*storage.TuneHarvestRunResult, error) {
	if run == nil {
		return nil, cause
	}
	message := strings.TrimSpace(strings.Join(append(warnings, cause.Error()), "\n"))
	finished, err := store.FinishTuneHarvestRun(run.ID, storage.TuneHarvestRunFailed, message, found, 0, 0, 0, 0)
	if err != nil {
		return nil, err
	}
	return &storage.TuneHarvestRunResult{Run: finished, Found: found, Warnings: append(warnings, cause.Error())}, cause
}

func dedupeCandidatesByShareCode(candidates []storage.TuneHarvestCandidateInput) []storage.TuneHarvestCandidateInput {
	indexByCode := map[string]int{}
	out := make([]storage.TuneHarvestCandidateInput, 0, len(candidates))
	for _, candidate := range candidates {
		code := storage.NormalizeTuneShareCode(candidate.ShareCode)
		if code == "" {
			continue
		}
		candidate.ShareCode = code
		if index, ok := indexByCode[code]; ok {
			if betterDuplicateCandidate(candidate, out[index]) {
				out[index] = candidate
			}
			continue
		}
		indexByCode[code] = len(out)
		out = append(out, candidate)
	}
	return out
}

func betterDuplicateCandidate(candidate storage.TuneHarvestCandidateInput, current storage.TuneHarvestCandidateInput) bool {
	candidateScore := duplicateCandidateScore(candidate)
	currentScore := duplicateCandidateScore(current)
	if candidateScore != currentScore {
		return candidateScore > currentScore
	}
	if candidate.MatchScore != current.MatchScore {
		return candidate.MatchScore > current.MatchScore
	}
	return len(candidate.RawJSON) > len(current.RawJSON)
}

func duplicateCandidateScore(candidate storage.TuneHarvestCandidateInput) int {
	score := 0
	if candidate.Status != storage.TuneHarvestCandidateRejected {
		score += 100
	}
	if candidate.MatchedCarID != "" {
		score += 60
	}
	if strings.TrimSpace(candidate.CarName+candidate.Make+candidate.Model+candidate.SourceCarID) != "" {
		score += 30
	}
	for _, value := range []string{candidate.UseCase, candidate.CarClass, candidate.Drivetrain, candidate.TireCompound, candidate.Tuner, candidate.TuneName, candidate.BestFor, candidate.Notes} {
		if strings.TrimSpace(value) != "" {
			score++
		}
	}
	if candidate.PI > 0 {
		score++
	}
	return score
}

func collectSource(ctx context.Context, client *http.Client, source string, limit int) (collectedSource, error) {
	switch source {
	case SourceCODMunity:
		return collectCODMunity(ctx, client)
	case SourceJSRSheet:
		return collectJSRSheet(ctx, client)
	case SourceForzaFire:
		return collectForzaFire(ctx, client, limit)
	default:
		return collectedSource{}, fmt.Errorf("unsupported source %q", source)
	}
}

func normalizeSources(sources []string) []string {
	seen := map[string]bool{}
	out := []string{}
	for _, source := range sources {
		source = strings.ToLower(strings.TrimSpace(source))
		switch source {
		case "cod", "codmunity":
			source = SourceCODMunity
		case "jsr", "jsr_chronic", "sheet", "google_sheet", "jsr_chronic_sheet":
			source = SourceJSRSheet
		case "forzafire", "forza_fire":
			source = SourceForzaFire
		}
		if source == "" || seen[source] {
			continue
		}
		seen[source] = true
		out = append(out, source)
	}
	if len(out) == 0 {
		return []string{SourceJSRSheet, SourceCODMunity}
	}
	return out
}

func collectJSRSheet(ctx context.Context, client *http.Client) (collectedSource, error) {
	body, err := fetchText(ctx, client, jsrSheetCSVURL)
	if err != nil {
		return collectedSource{}, err
	}
	reader := csv.NewReader(strings.NewReader(body))
	reader.FieldsPerRecord = -1
	rows, err := reader.ReadAll()
	if err != nil {
		return collectedSource{}, err
	}
	if len(rows) == 0 {
		return collectedSource{}, errors.New("empty CSV")
	}
	headers := map[string]int{}
	for index, header := range rows[0] {
		headers[normalizeHeader(header)] = index
	}
	var out collectedSource
	for rowIndex, row := range rows[1:] {
		if ctx.Err() != nil {
			return out, ctx.Err()
		}
		year := parseIntField(csvField(row, headers, "year"))
		makeName := csvField(row, headers, "make")
		model := csvField(row, headers, "model")
		if year > 0 && makeName != "" && model != "" {
			out.Cars = append(out.Cars, storage.FH6CarInput{
				Year: year, Make: makeName, Model: model,
				Source: SourceJSRSheet, SourceRef: fmt.Sprintf("row:%d", rowIndex+2),
			})
		}
		code := storage.NormalizeTuneShareCode(csvField(row, headers, "share code"))
		if code == "" {
			continue
		}
		raw := rowToMap(rows[0], row)
		rawJSON, _ := json.Marshal(raw)
		buildType := csvField(row, headers, "build type")
		classText := csvField(row, headers, "class")
		tuneName := csvField(row, headers, "tune name")
		bestFor := csvField(row, headers, "best for")
		notes := csvField(row, headers, "notes")
		carClass, pi := extractClassPI(strings.Join([]string{classText, buildType, tuneName}, " "))
		out.Candidates = append(out.Candidates, storage.TuneHarvestCandidateInput{
			Source:       SourceJSRSheet,
			SourceRef:    fmt.Sprintf("row:%d", rowIndex+2),
			SourceURL:    jsrSheetCSVURL,
			RawKey:       fmt.Sprintf("jsr:%d:%s:%d:%s:%s:%s", rowIndex+2, code, year, slug(makeName), slug(model), slug(tuneName)),
			ShareCode:    code,
			Year:         year,
			Make:         makeName,
			Model:        model,
			CarName:      fullCarName(year, makeName, model),
			UseCase:      inferUseCase(strings.Join([]string{buildType, tuneName, bestFor, notes}, " ")),
			CarClass:     carClass,
			PI:           pi,
			Drivetrain:   normalizeDrivetrain(csvField(row, headers, "drivetrain")),
			TireCompound: inferTireCompound(strings.Join([]string{buildType, tuneName, bestFor, notes}, " ")),
			TuneName:     tuneName,
			BestFor:      bestFor,
			Difficulty:   csvField(row, headers, "difficulty"),
			Notes:        notes,
			RawJSON:      string(rawJSON),
		})
	}
	return out, nil
}

func collectCODMunity(ctx context.Context, client *http.Client) (collectedSource, error) {
	var out collectedSource
	for _, endpoint := range codmunityEndpoints {
		if ctx.Err() != nil {
			return out, ctx.Err()
		}
		next, err := collectCODMunityEndpoint(ctx, client, endpoint)
		if ctx.Err() != nil {
			out.Cars = append(out.Cars, next.Cars...)
			out.Candidates = append(out.Candidates, next.Candidates...)
			out.Warnings = append(out.Warnings, next.Warnings...)
			return out, ctx.Err()
		}
		if err != nil {
			out.Warnings = append(out.Warnings, fmt.Sprintf("codmunity %s: %v", endpoint.Key, err))
			continue
		}
		out.Cars = append(out.Cars, next.Cars...)
		out.Candidates = append(out.Candidates, next.Candidates...)
		out.Warnings = append(out.Warnings, next.Warnings...)
	}
	if len(out.Candidates) == 0 && len(out.Warnings) > 0 {
		return out, errors.New(strings.Join(out.Warnings, "; "))
	}
	return out, nil
}

func collectCODMunityEndpoint(ctx context.Context, client *http.Client, endpoint codmunityEndpoint) (collectedSource, error) {
	body, err := fetchText(ctx, client, endpoint.URL)
	if err != nil {
		return collectedSource{}, err
	}
	state, err := extractScriptJSON(body, "serverApp-state")
	if err != nil {
		return collectedSource{}, err
	}
	var root any
	if err := json.Unmarshal([]byte(state), &root); err != nil {
		return collectedSource{}, err
	}
	carMaps := findCODMunityCars(root)
	var out collectedSource
	for _, carMap := range carMaps {
		if ctx.Err() != nil {
			return out, ctx.Err()
		}
		year := intValue(carMap["Year"])
		makeName := stringValue(carMap["Manufacturer"])
		model := stringValue(carMap["CarName"])
		sourceCarID := stringValue(carMap["ID"])
		if sourceCarID == "" {
			sourceCarID = strconv.Itoa(intValue(carMap["ID"]))
		}
		stats := mapValue(carMap["stats"])
		basePI := intValue(carMap["PI Value"])
		drive := normalizeDrivetrain(stringValue(stats["Drive"]))
		if year > 0 && makeName != "" && model != "" {
			out.Cars = append(out.Cars, storage.FH6CarInput{
				Year: year, Make: makeName, Model: model, BasePI: basePI, DrivetrainDefault: drive,
				Source: SourceCODMunity, SourceRef: sourceCarID,
			})
		}
		tunings := sliceValue(carMap["Tunings"])
		for _, item := range tunings {
			if ctx.Err() != nil {
				return out, ctx.Err()
			}
			tuning := mapValue(item)
			code := storage.NormalizeTuneShareCode(firstString(tuning, "Sharecode", "ShareCode", "shareCode"))
			if code == "" {
				continue
			}
			comment := stringValue(tuning["Comment"])
			carClass := strings.ToUpper(strings.TrimSpace(stringValue(tuning["Class"])))
			pi := intValue(tuning["PI Value"])
			rawJSON, _ := json.Marshal(map[string]any{"car": carMap, "tuning": tuning})
			playstyle := stringValue(tuning["Playstyle"])
			useCase := inferUseCase(strings.Join([]string{endpoint.Key, playstyle}, " "))
			out.Candidates = append(out.Candidates, storage.TuneHarvestCandidateInput{
				Source:       SourceCODMunity,
				SourceRef:    stringValue(tuning["_id"]),
				SourceURL:    endpoint.URL,
				SourceCarID:  sourceCarID,
				RawKey:       fmt.Sprintf("codmunity:%s:%s:%s:%s", sourceCarID, code, slug(playstyle), strings.ToLower(carClass)),
				ShareCode:    code,
				Year:         year,
				Make:         makeName,
				Model:        model,
				CarName:      fullCarName(year, makeName, model),
				UseCase:      useCase,
				CarClass:     carClass,
				PI:           pi,
				Drivetrain:   normalizeDrivetrain(stringValue(tuning["Drivetrain"])),
				TireCompound: inferTireCompound(comment),
				Tuner:        stringValue(tuning["Tuner"]),
				TuneName:     strings.TrimSpace(strings.Join([]string{endpoint.Key, stringValue(tuning["Meta"])}, " / ")),
				BestFor:      playstyle,
				Notes:        comment,
				RawJSON:      string(rawJSON),
			})
		}
	}
	return out, nil
}

func collectForzaFire(ctx context.Context, client *http.Client, limit int) (collectedSource, error) {
	body, err := fetchText(ctx, client, forzaFireBuildsURL)
	if err != nil {
		return collectedSource{}, err
	}
	links := parseForzaFireBuildLinks(body)
	if len(links) == 0 {
		return collectedSource{}, errors.New("no build links found")
	}
	if limit > 0 && len(links) > limit {
		links = links[:limit]
	}
	var out collectedSource
	for _, link := range links {
		if ctx.Err() != nil {
			return out, ctx.Err()
		}
		detail, err := fetchText(ctx, client, link.URL)
		if err != nil {
			if ctx.Err() != nil {
				return out, ctx.Err()
			}
			out.Warnings = append(out.Warnings, fmt.Sprintf("forzafire detail %s: %v", link.URL, err))
			continue
		}
		candidate, car, ok := parseForzaFireDetail(detail, link)
		if !ok {
			continue
		}
		if car.Year > 0 && car.Make != "" && car.Model != "" {
			out.Cars = append(out.Cars, car)
		}
		out.Candidates = append(out.Candidates, candidate)
	}
	return out, nil
}

func fetchText(ctx context.Context, client *http.Client, targetURL string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, targetURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "FH6Worker tune harvest/1.0")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	data, err := io.ReadAll(io.LimitReader(resp.Body, 16*1024*1024))
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func enrichCandidate(candidate storage.TuneHarvestCandidateInput, cars []storage.FH6Car) storage.TuneHarvestCandidateInput {
	candidate.ShareCode = storage.NormalizeTuneShareCode(candidate.ShareCode)
	candidate.UseCase = inferUseCase(candidate.UseCase)
	if candidate.UseCase == "" {
		candidate.UseCase = inferUseCase(strings.Join([]string{candidate.TuneName, candidate.BestFor, candidate.Notes}, " "))
	}
	if candidate.CarClass == "" || candidate.PI == 0 {
		carClass, pi := extractClassPI(strings.Join([]string{candidate.CarClass, candidate.TuneName, candidate.BestFor, candidate.Notes}, " "))
		if candidate.CarClass == "" {
			candidate.CarClass = carClass
		}
		if candidate.PI == 0 {
			candidate.PI = pi
		}
	}
	if candidate.PI == 0 && candidate.CarClass != "" {
		candidate.PI = classDefaultPI(candidate.CarClass)
	}
	if candidate.Drivetrain == "" {
		candidate.Drivetrain = normalizeDrivetrain(strings.Join([]string{candidate.TuneName, candidate.BestFor, candidate.Notes}, " "))
	}
	if candidate.TireCompound == "" {
		candidate.TireCompound = inferTireCompound(strings.Join([]string{candidate.TuneName, candidate.BestFor, candidate.Notes}, " "))
	}
	if candidate.CarName == "" {
		candidate.CarName = fullCarName(candidate.Year, candidate.Make, candidate.Model)
	}
	contextText := strings.Join([]string{candidate.TuneName, candidate.BestFor, candidate.Notes, candidate.RawJSON}, " ")
	if looksLikeNonTuneCode(contextText) {
		candidate.Status = storage.TuneHarvestCandidateRejected
		candidate.RejectionReason = "non_tune_code_context"
		return candidate
	}
	if strings.TrimSpace(candidate.CarName+candidate.Make+candidate.Model+candidate.SourceCarID) == "" {
		candidate.Status = storage.TuneHarvestCandidatePending
		candidate.MatchReason = "missing_vehicle_context"
		return candidate
	}
	match := MatchCar(candidate, cars)
	candidate.MatchedCarID = match.CarID
	candidate.MatchScore = match.Score
	candidate.MatchReason = match.Reason
	if candidate.Status == "" {
		candidate.Status = storage.TuneHarvestCandidatePending
	}
	return candidate
}

type CarMatch struct {
	CarID  string
	Score  float64
	Reason string
}

func MatchCar(candidate storage.TuneHarvestCandidateInput, cars []storage.FH6Car) CarMatch {
	best := CarMatch{}
	second := 0.0
	for _, car := range cars {
		score := scoreCarMatch(candidate, car)
		if score > best.Score {
			second = best.Score
			best = CarMatch{CarID: car.CarID, Score: score}
		} else if score > second {
			second = score
		}
	}
	if best.Score < 0.88 {
		best.CarID = ""
		best.Reason = "below_threshold"
		return best
	}
	if best.Score-second < 0.025 {
		best.CarID = ""
		best.Reason = "ambiguous_match"
		return best
	}
	best.Score = math.Round(best.Score*1000) / 1000
	best.Reason = "matched"
	return best
}

func scoreCarMatch(candidate storage.TuneHarvestCandidateInput, car storage.FH6Car) float64 {
	yearScore := 0.35
	if candidate.Year > 0 {
		if candidate.Year == car.Year {
			yearScore = 1
		} else {
			yearScore = 0
		}
	}
	candidateCarText := normalizeText(strings.Join([]string{candidate.CarName, candidate.Make, candidate.Model}, " "))
	carMake := normalizeText(car.Make)
	makeScore := 0.0
	if normalizeText(candidate.Make) == carMake && carMake != "" {
		makeScore = 1
	} else if candidateCarText != "" && containsTokenSequence(candidateCarText, carMake) {
		makeScore = 1
	}

	names := []string{candidate.Model, candidate.CarName}
	modelScore := 0.0
	aliases := append([]string{car.Model, fullCarName(car.Year, car.Make, car.Model)}, car.Aliases...)
	for _, name := range names {
		for _, alias := range aliases {
			modelScore = math.Max(modelScore, textSimilarity(name, alias))
		}
	}
	sourceScore := 0.0
	if candidate.Source != "" && candidate.Source == car.Source && candidate.SourceCarID != "" && candidate.SourceCarID == car.SourceRef {
		sourceScore = 1
	}
	return yearScore*0.25 + makeScore*0.25 + modelScore*0.40 + sourceScore*0.10
}

func textSimilarity(a string, b string) float64 {
	na := normalizeText(a)
	nb := normalizeText(b)
	if na == "" || nb == "" {
		return 0
	}
	if na == nb {
		return 1
	}
	if containsTokenSequence(na, nb) || containsTokenSequence(nb, na) {
		return 0.92
	}
	at := tokenSet(na)
	bt := tokenSet(nb)
	if len(at) == 0 || len(bt) == 0 {
		return 0
	}
	intersection := 0
	for token := range at {
		if bt[token] {
			intersection++
		}
	}
	union := len(at) + len(bt) - intersection
	if union == 0 {
		return 0
	}
	return float64(intersection) / float64(union)
}

func containsTokenSequence(haystack string, needle string) bool {
	haystack = " " + normalizeText(haystack) + " "
	needle = " " + normalizeText(needle) + " "
	return strings.Contains(haystack, needle)
}

func tokenSet(value string) map[string]bool {
	out := map[string]bool{}
	for _, token := range strings.Fields(normalizeText(value)) {
		if token != "" && token != "the" {
			out[token] = true
		}
	}
	return out
}

func normalizeText(value string) string {
	replacer := strings.NewReplacer(
		"é", "e", "è", "e", "ê", "e", "ë", "e", "É", "e",
		"á", "a", "à", "a", "ä", "a", "â", "a",
		"ó", "o", "ö", "o", "ô", "o",
		"ú", "u", "ü", "u",
		"í", "i", "ï", "i",
		"ñ", "n",
		"’", "'", "‘", "'", "“", "\"", "”", "\"",
	)
	value = strings.ToLower(replacer.Replace(value))
	var b strings.Builder
	lastSpace := false
	for _, r := range value {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			b.WriteRune(r)
			lastSpace = false
			continue
		}
		if !lastSpace {
			b.WriteByte(' ')
			lastSpace = true
		}
	}
	out := strings.Join(strings.Fields(b.String()), " ")
	out = strings.ReplaceAll(out, "competition", "comp")
	out = strings.ReplaceAll(out, "coupé", "coupe")
	return strings.TrimSpace(out)
}

func extractScriptJSON(body string, id string) (string, error) {
	pattern := regexp.MustCompile(`(?is)<script[^>]+id=["']` + regexp.QuoteMeta(id) + `["'][^>]*>(.*?)</script>`)
	match := pattern.FindStringSubmatch(body)
	if len(match) != 2 {
		return "", fmt.Errorf("script %q not found", id)
	}
	return html.UnescapeString(strings.TrimSpace(match[1])), nil
}

func findCODMunityCars(root any) []map[string]any {
	var cars []map[string]any
	var walk func(any)
	walk = func(value any) {
		switch v := value.(type) {
		case map[string]any:
			if _, ok := v["Tunings"]; ok && stringValue(v["CarName"]) != "" && stringValue(v["Manufacturer"]) != "" {
				cars = append(cars, v)
				return
			}
			for _, child := range v {
				walk(child)
			}
		case []any:
			for _, child := range v {
				walk(child)
			}
		}
	}
	walk(root)
	return cars
}

type forzaFireBuildLink struct {
	URL        string
	Path       string
	Search     string
	CarClass   string
	RaceType   string
	Difficulty string
}

var forzaFireBuildLinkPattern = regexp.MustCompile(`(?is)<a\b[^>]*href=["'](/build/[^"']+)["'][^>]*browse-builds__item[^>]*>`)

func parseForzaFireBuildLinks(body string) []forzaFireBuildLink {
	matches := forzaFireBuildLinkPattern.FindAllStringSubmatch(body, -1)
	seen := map[string]bool{}
	links := []forzaFireBuildLink{}
	for _, match := range matches {
		if len(match) < 2 {
			continue
		}
		path := html.UnescapeString(match[1])
		if seen[path] {
			continue
		}
		seen[path] = true
		tag := match[0]
		links = append(links, forzaFireBuildLink{
			URL:        "https://www.forzafire.com" + path,
			Path:       path,
			Search:     attrValue(tag, "data-search"),
			CarClass:   strings.ToUpper(strings.TrimSpace(attrValue(tag, "data-class"))),
			RaceType:   attrValue(tag, "data-race-type"),
			Difficulty: attrValue(tag, "data-difficulty"),
		})
	}
	return links
}

func parseForzaFireDetail(body string, link forzaFireBuildLink) (storage.TuneHarvestCandidateInput, storage.FH6CarInput, bool) {
	code := storage.NormalizeTuneShareCode(firstRegex(body, `(?is)<input[^>]+name=["']tuner_code["'][^>]+value=["']([^"']+)["']`))
	if code == "" {
		return storage.TuneHarvestCandidateInput{}, storage.FH6CarInput{}, false
	}
	loadBuildJSON, ok := extractJSObject(body, "loadBuild")
	if !ok {
		return storage.TuneHarvestCandidateInput{}, storage.FH6CarInput{}, false
	}
	var payload map[string]any
	if err := json.Unmarshal([]byte(loadBuildJSON), &payload); err != nil {
		return storage.TuneHarvestCandidateInput{}, storage.FH6CarInput{}, false
	}
	carMap := mapValue(payload["car"])
	buildMap := mapValue(payload["build"])
	tireMap := mapValue(buildMap["tires"])
	carName := stringValue(carMap["name"])
	year, makeName, model := splitCarName(carName)
	if carName == "" {
		carName = link.Search
		year, makeName, model = splitCarName(link.Search)
	}
	tuneName := stringValue(payload["name"])
	classText, pi := extractClassPI(strings.Join([]string{link.CarClass, link.Search, tuneName}, " "))
	if classText == "" {
		classText = link.CarClass
	}
	sourceRef := stringValue(carMap["id"])
	rawJSON, _ := json.Marshal(payload)
	car := storage.FH6CarInput{
		Year: year, Make: makeName, Model: model,
		Source: SourceForzaFire, SourceRef: sourceRef,
	}
	candidate := storage.TuneHarvestCandidateInput{
		Source:       SourceForzaFire,
		SourceRef:    strings.TrimPrefix(link.Path, "/build/"),
		SourceURL:    link.URL,
		SourceCarID:  sourceRef,
		RawKey:       fmt.Sprintf("forzafire:%s:%s", strings.TrimPrefix(link.Path, "/build/"), code),
		ShareCode:    code,
		Year:         year,
		Make:         makeName,
		Model:        model,
		CarName:      carName,
		UseCase:      inferUseCase(strings.Join([]string{link.RaceType, link.Search, tuneName}, " ")),
		CarClass:     classText,
		PI:           pi,
		Drivetrain:   normalizeDrivetrain(strings.Join([]string{link.Search, tuneName}, " ")),
		TireCompound: inferTireCompound(stringValue(tireMap["Tire Compound"])),
		TuneName:     tuneName,
		BestFor:      link.RaceType,
		Difficulty:   link.Difficulty,
		RawJSON:      string(rawJSON),
	}
	return candidate, car, true
}

func extractJSObject(body string, varName string) (string, bool) {
	startPattern := regexp.MustCompile(`var\s+` + regexp.QuoteMeta(varName) + `\s*=`)
	loc := startPattern.FindStringIndex(body)
	if loc == nil {
		return "", false
	}
	start := strings.Index(body[loc[1]:], "{")
	if start < 0 {
		return "", false
	}
	start += loc[1]
	depth := 0
	inString := false
	escapeNext := false
	for index := start; index < len(body); index++ {
		ch := body[index]
		if inString {
			if escapeNext {
				escapeNext = false
				continue
			}
			if ch == '\\' {
				escapeNext = true
				continue
			}
			if ch == '"' {
				inString = false
			}
			continue
		}
		switch ch {
		case '"':
			inString = true
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return body[start : index+1], true
			}
		}
	}
	return "", false
}

func attrValue(tag string, name string) string {
	pattern := regexp.MustCompile(`(?is)\b` + regexp.QuoteMeta(name) + `=["']([^"']*)["']`)
	match := pattern.FindStringSubmatch(tag)
	if len(match) != 2 {
		return ""
	}
	return html.UnescapeString(strings.TrimSpace(match[1]))
}

func firstRegex(value string, pattern string) string {
	re := regexp.MustCompile(pattern)
	match := re.FindStringSubmatch(value)
	if len(match) < 2 {
		return ""
	}
	return html.UnescapeString(strings.TrimSpace(match[1]))
}

func normalizeHeader(value string) string {
	return strings.Join(strings.Fields(strings.ToLower(strings.TrimSpace(value))), " ")
}

func csvField(row []string, headers map[string]int, name string) string {
	index, ok := headers[normalizeHeader(name)]
	if !ok || index < 0 || index >= len(row) {
		return ""
	}
	return strings.TrimSpace(row[index])
}

func rowToMap(headers []string, row []string) map[string]string {
	out := map[string]string{}
	for index, header := range headers {
		if index < len(row) {
			out[header] = row[index]
		}
	}
	return out
}

func parseIntField(value string) int {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0
	}
	match := regexp.MustCompile(`\d+`).FindString(value)
	if match == "" {
		return 0
	}
	parsed, _ := strconv.Atoi(match)
	return parsed
}

func inferUseCase(value string) string {
	lower := strings.ToLower(value)
	switch {
	case strings.Contains(lower, "drift"):
		return "Drift"
	case strings.Contains(lower, "drag"):
		return "Drag"
	case strings.Contains(lower, "cross-country"), strings.Contains(lower, "cross country"), strings.Contains(lower, "offroad"), strings.Contains(lower, "off-road"):
		return "Offroad"
	case strings.Contains(lower, "rally"), strings.Contains(lower, "dirt"):
		return "Rally"
	case strings.Contains(lower, "road"), strings.Contains(lower, "street"), strings.Contains(lower, "touge"), strings.Contains(lower, "purist"), strings.Contains(lower, "pr stunts"), strings.Contains(lower, "pr-stunts"), strings.Contains(lower, "rival"), strings.Contains(lower, "grip"), strings.Contains(lower, "all around"), strings.Contains(lower, "meta"):
		return "Road"
	default:
		return ""
	}
}

func inferTireCompound(value string) string {
	lower := strings.ToLower(value)
	switch {
	case strings.Contains(lower, "semi"):
		return "semi"
	case strings.Contains(lower, "slick"), strings.Contains(lower, "sllck"):
		return "slick"
	case strings.Contains(lower, "rally"):
		return "rally"
	case strings.Contains(lower, "offroad"), strings.Contains(lower, "off-road"):
		return "offroad"
	case strings.Contains(lower, "drift"):
		return "drift"
	case strings.Contains(lower, "drag"):
		return "drag"
	case strings.Contains(lower, "snow"):
		return "snow"
	case strings.Contains(lower, "sport"):
		return "sport"
	case strings.Contains(lower, "street"):
		return "street"
	case strings.Contains(lower, "stock"):
		return "stock"
	default:
		return ""
	}
}

func normalizeDrivetrain(value string) string {
	upper := strings.ToUpper(value)
	switch {
	case strings.Contains(upper, "AWD"):
		return "AWD"
	case strings.Contains(upper, "FWD"):
		return "FWD"
	case strings.Contains(upper, "RWD"):
		return "RWD"
	default:
		return ""
	}
}

var classPIPattern = regexp.MustCompile(`(?i)\b(S1|S2|[ABCDRX])\s*[- ]?(\d{3})?\b`)

func extractClassPI(value string) (string, int) {
	match := classPIPattern.FindStringSubmatch(value)
	if len(match) == 0 {
		return "", 0
	}
	carClass := strings.ToUpper(match[1])
	pi := 0
	if len(match) > 2 && match[2] != "" {
		pi, _ = strconv.Atoi(match[2])
	}
	if pi == 0 {
		pi = classDefaultPI(carClass)
	}
	return carClass, pi
}

func classDefaultPI(carClass string) int {
	switch strings.ToUpper(strings.TrimSpace(carClass)) {
	case "D":
		return 400
	case "C":
		return 500
	case "B":
		return 600
	case "A":
		return 700
	case "S1":
		return 800
	case "S2":
		return 900
	case "R":
		return 998
	case "X":
		return 999
	default:
		return 0
	}
}

func looksLikeNonTuneCode(value string) bool {
	lower := strings.ToLower(value)
	negative := []string{"livery", "paint", "vinyl", "design", "blueprint", "eventlab", "route creator", "challenge card", "涂装", "蓝图", "赛事蓝图"}
	for _, word := range negative {
		if strings.Contains(lower, word) {
			return true
		}
	}
	return false
}

func fullCarName(year int, makeName string, model string) string {
	parts := []string{}
	if year > 0 {
		parts = append(parts, strconv.Itoa(year))
	}
	if strings.TrimSpace(makeName) != "" {
		parts = append(parts, strings.TrimSpace(makeName))
	}
	if strings.TrimSpace(model) != "" {
		parts = append(parts, strings.TrimSpace(model))
	}
	return strings.Join(parts, " ")
}

var knownMakes = []string{
	"AMG Transport Dynamics",
	"Alfa Romeo",
	"Aston Martin",
	"Sierra Sierra Enterprises",
	"Mercedes-AMG",
	"Mercedes-Benz",
	"Formula Drift",
	"Hot Wheels",
	"Land Rover",
	"MINI",
	"BMW",
	"Acura",
	"Abarth",
	"Audi",
	"Ford",
	"Honda",
	"Mazda",
	"Nissan",
	"Porsche",
	"Toyota",
}

func splitCarName(value string) (int, string, string) {
	words := strings.Fields(strings.TrimSpace(value))
	if len(words) < 3 {
		return 0, "", strings.TrimSpace(value)
	}
	year, err := strconv.Atoi(words[0])
	if err != nil || year < 1900 {
		return 0, "", strings.TrimSpace(value)
	}
	rest := strings.TrimSpace(strings.TrimPrefix(value, words[0]))
	for _, makeName := range knownMakes {
		if strings.HasPrefix(strings.ToLower(rest), strings.ToLower(makeName)+" ") {
			return year, makeName, strings.TrimSpace(rest[len(makeName):])
		}
	}
	return year, words[1], strings.TrimSpace(strings.Join(words[2:], " "))
}

func slug(value string) string {
	value = normalizeText(value)
	value = strings.ReplaceAll(value, " ", "-")
	return strings.Trim(value, "-")
}

func stringValue(value any) string {
	switch v := value.(type) {
	case string:
		return strings.TrimSpace(v)
	case json.Number:
		return v.String()
	case float64:
		if math.Trunc(v) == v {
			return strconv.Itoa(int(v))
		}
		return fmt.Sprintf("%v", v)
	case int:
		return strconv.Itoa(v)
	default:
		return ""
	}
}

func intValue(value any) int {
	switch v := value.(type) {
	case int:
		return v
	case float64:
		return int(v)
	case json.Number:
		parsed, _ := v.Int64()
		return int(parsed)
	case string:
		return parseIntField(v)
	default:
		return 0
	}
}

func firstString(values map[string]any, keys ...string) string {
	for _, key := range keys {
		if value := stringValue(values[key]); value != "" {
			return value
		}
	}
	return ""
}

func mapValue(value any) map[string]any {
	if v, ok := value.(map[string]any); ok {
		return v
	}
	return map[string]any{}
}

func sliceValue(value any) []any {
	if v, ok := value.([]any); ok {
		return v
	}
	return nil
}
