package storage

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

const (
	TuneHarvestCandidatePending  = "pending"
	TuneHarvestCandidateRejected = "rejected"
	TuneHarvestCandidateImported = "imported"

	TuneHarvestRunRunning  = "running"
	TuneHarvestRunComplete = "complete"
	TuneHarvestRunFailed   = "failed"
	TuneHarvestRunCanceled = "cancelled"
)

type FH6CarInput struct {
	CarID             string   `json:"carId"`
	Year              int      `json:"year"`
	Make              string   `json:"make"`
	Model             string   `json:"model"`
	Aliases           []string `json:"alias"`
	BasePI            int      `json:"basePi"`
	DrivetrainDefault string   `json:"drivetrainDefault"`
	Source            string   `json:"source"`
	SourceRef         string   `json:"sourceRef"`
}

type FH6Car struct {
	FH6CarInput
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type TuneHarvestOptions struct {
	Sources        []string `json:"sources"`
	DryRun         bool     `json:"dryRun"`
	LimitPerSource int      `json:"limitPerSource"`
}

type TuneHarvestRunInput struct {
	Sources []string `json:"sources"`
	DryRun  bool     `json:"dryRun"`
}

type TuneHarvestRun struct {
	ID            int64    `json:"id"`
	StartedAt     string   `json:"startedAt"`
	FinishedAt    string   `json:"finishedAt"`
	Sources       []string `json:"sources"`
	DryRun        bool     `json:"dryRun"`
	Status        string   `json:"status"`
	Message       string   `json:"message"`
	FoundCount    int      `json:"foundCount"`
	SavedCount    int      `json:"savedCount"`
	RejectedCount int      `json:"rejectedCount"`
	PendingCount  int      `json:"pendingCount"`
	ImportedCount int      `json:"importedCount"`
}

type TuneHarvestCandidateInput struct {
	RunID           int64   `json:"runId"`
	Source          string  `json:"source"`
	SourceRef       string  `json:"sourceRef"`
	SourceURL       string  `json:"sourceUrl"`
	SourceCarID     string  `json:"sourceCarId"`
	RawKey          string  `json:"rawKey"`
	ShareCode       string  `json:"shareCode"`
	Year            int     `json:"year"`
	Make            string  `json:"make"`
	Model           string  `json:"model"`
	CarName         string  `json:"carName"`
	MatchedCarID    string  `json:"matchedCarId"`
	MatchScore      float64 `json:"matchScore"`
	MatchReason     string  `json:"matchReason"`
	UseCase         string  `json:"useCase"`
	CarClass        string  `json:"carClass"`
	PI              int     `json:"pi"`
	Drivetrain      string  `json:"drivetrain"`
	TireCompound    string  `json:"tireCompound"`
	Tuner           string  `json:"tuner"`
	TuneName        string  `json:"tuneName"`
	BestFor         string  `json:"bestFor"`
	Difficulty      string  `json:"difficulty"`
	Notes           string  `json:"notes"`
	RawJSON         string  `json:"rawJson"`
	Status          string  `json:"status"`
	RejectionReason string  `json:"rejectionReason"`
}

type TuneHarvestCandidate struct {
	ID int64 `json:"id"`
	TuneHarvestCandidateInput
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type TuneHarvestRunResult struct {
	Run        *TuneHarvestRun        `json:"run,omitempty"`
	Candidates []TuneHarvestCandidate `json:"candidates"`
	Found      int                    `json:"found"`
	Saved      int                    `json:"saved"`
	Rejected   int                    `json:"rejected"`
	Pending    int                    `json:"pending"`
	Imported   int                    `json:"imported"`
	Warnings   []string               `json:"warnings"`
}

func NormalizeFH6CarInput(input FH6CarInput) (FH6CarInput, error) {
	input.CarID = strings.TrimSpace(input.CarID)
	input.Make = strings.TrimSpace(input.Make)
	input.Model = strings.TrimSpace(input.Model)
	input.DrivetrainDefault = strings.ToUpper(strings.TrimSpace(input.DrivetrainDefault))
	input.Source = strings.TrimSpace(input.Source)
	input.SourceRef = strings.TrimSpace(input.SourceRef)
	if input.Year <= 0 {
		return input, errors.New("year is required")
	}
	if input.Make == "" {
		return input, errors.New("make is required")
	}
	if input.Model == "" {
		return input, errors.New("model is required")
	}
	if input.CarID == "" {
		input.CarID = GenerateFH6CarID(input.Year, input.Make, input.Model)
	}
	aliases := map[string]bool{}
	for _, alias := range DefaultFH6CarAliases(input.Year, input.Make, input.Model) {
		alias = strings.TrimSpace(alias)
		if alias != "" {
			aliases[alias] = true
		}
	}
	for _, alias := range input.Aliases {
		alias = strings.TrimSpace(alias)
		if alias != "" {
			aliases[alias] = true
		}
	}
	input.Aliases = make([]string, 0, len(aliases))
	for alias := range aliases {
		input.Aliases = append(input.Aliases, alias)
	}
	sort.Strings(input.Aliases)
	return input, nil
}

func GenerateFH6CarID(year int, makeName string, model string) string {
	parts := []string{"fh6", strconv.Itoa(year), slugifyASCII(makeName), slugifyASCII(model)}
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.Trim(part, "-")
		if part != "" {
			out = append(out, part)
		}
	}
	return strings.Join(out, ":")
}

func DefaultFH6CarAliases(year int, makeName string, model string) []string {
	model = strings.TrimSpace(model)
	makeName = strings.TrimSpace(makeName)
	if model == "" {
		return nil
	}
	aliases := []string{
		model,
		strings.Join([]string{makeName, model}, " "),
		strings.Join([]string{strconv.Itoa(year), makeName, model}, " "),
	}
	replacements := []struct {
		old string
		new string
	}{
		{old: " Competition ", new: " Comp "},
		{old: " Coupe", new: ""},
		{old: " Coupé", new: ""},
		{old: " Forza Edition", new: " FE"},
		{old: " Type-R", new: " Type R"},
		{old: " Type S", new: " Type-S"},
	}
	for _, replacement := range replacements {
		if strings.Contains(model, strings.TrimSpace(replacement.old)) || strings.HasSuffix(model, replacement.old) {
			aliases = append(aliases, strings.TrimSpace(strings.ReplaceAll(model, replacement.old, replacement.new)))
			aliases = append(aliases, strings.TrimSpace(strings.Join([]string{makeName, strings.ReplaceAll(model, replacement.old, replacement.new)}, " ")))
		}
	}
	return aliases
}

func (s *Store) SaveFH6Cars(inputs []FH6CarInput) (int, error) {
	if len(inputs) == 0 {
		return 0, nil
	}
	tx, err := s.db.Begin()
	if err != nil {
		return 0, err
	}
	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback()
		}
	}()
	now := nowText()
	count := 0
	for index, input := range inputs {
		normalized, err := NormalizeFH6CarInput(input)
		if err != nil {
			return 0, fmt.Errorf("car %d: %w", index+1, err)
		}
		aliasJSON, err := json.Marshal(normalized.Aliases)
		if err != nil {
			return 0, err
		}
		_, err = tx.Exec(`INSERT INTO fh6_car (
			car_id, year, make, model, alias_json, base_pi, drivetrain_default, source, source_ref, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(car_id) DO UPDATE SET
			year = excluded.year,
			make = excluded.make,
			model = excluded.model,
			alias_json = excluded.alias_json,
			base_pi = CASE WHEN excluded.base_pi > 0 THEN excluded.base_pi ELSE fh6_car.base_pi END,
			drivetrain_default = CASE WHEN excluded.drivetrain_default <> '' THEN excluded.drivetrain_default ELSE fh6_car.drivetrain_default END,
			source = excluded.source,
			source_ref = excluded.source_ref,
			updated_at = excluded.updated_at`,
			normalized.CarID, normalized.Year, normalized.Make, normalized.Model, string(aliasJSON),
			normalized.BasePI, normalized.DrivetrainDefault, normalized.Source, normalized.SourceRef, now, now)
		if err != nil {
			return 0, err
		}
		count++
	}
	if err := tx.Commit(); err != nil {
		return 0, err
	}
	committed = true
	return count, nil
}

func (s *Store) ListFH6Cars() ([]FH6Car, error) {
	rows, err := s.db.Query(`SELECT car_id, year, make, model, COALESCE(alias_json, '[]'), COALESCE(base_pi, 0),
		COALESCE(drivetrain_default, ''), COALESCE(source, ''), COALESCE(source_ref, ''), created_at, updated_at
		FROM fh6_car ORDER BY make, model, year`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	cars := []FH6Car{}
	for rows.Next() {
		car, err := scanFH6Car(rows)
		if err != nil {
			return nil, err
		}
		cars = append(cars, car)
	}
	return cars, rows.Err()
}

func (s *Store) CreateTuneHarvestRun(input TuneHarvestRunInput) (*TuneHarvestRun, error) {
	sources := cleanTuneHarvestSources(input.Sources)
	sourcesJSON, err := json.Marshal(sources)
	if err != nil {
		return nil, err
	}
	now := nowText()
	result, err := s.db.Exec(`INSERT INTO tune_harvest_run(started_at, sources_json, dry_run, status, message)
		VALUES (?, ?, ?, ?, '')`, now, string(sourcesJSON), boolInt(input.DryRun), TuneHarvestRunRunning)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	return s.GetTuneHarvestRun(id)
}

func (s *Store) GetTuneHarvestRun(id int64) (*TuneHarvestRun, error) {
	row := s.db.QueryRow(`SELECT id, started_at, COALESCE(finished_at, ''), COALESCE(sources_json, '[]'), dry_run,
		status, COALESCE(message, ''), found_count, saved_count, rejected_count, pending_count, imported_count
		FROM tune_harvest_run WHERE id = ?`, id)
	run, err := scanTuneHarvestRun(row)
	if err != nil {
		return nil, err
	}
	return &run, nil
}

func (s *Store) FinishTuneHarvestRun(id int64, status string, message string, found int, saved int, rejected int, pending int, imported int) (*TuneHarvestRun, error) {
	status = strings.TrimSpace(status)
	if status == "" {
		status = TuneHarvestRunComplete
	}
	_, err := s.db.Exec(`UPDATE tune_harvest_run SET finished_at = ?, status = ?, message = ?,
		found_count = ?, saved_count = ?, rejected_count = ?, pending_count = ?, imported_count = ?
		WHERE id = ?`, nowText(), status, strings.TrimSpace(message), found, saved, rejected, pending, imported, id)
	if err != nil {
		return nil, err
	}
	return s.GetTuneHarvestRun(id)
}

func (s *Store) dedupeTuneHarvestCandidatesByShareCode() error {
	rows, err := s.db.Query(`SELECT share_code FROM tune_harvest_candidate GROUP BY share_code HAVING COUNT(*) > 1`)
	if err != nil {
		return err
	}
	defer rows.Close()
	duplicateCodes := []string{}
	for rows.Next() {
		var shareCode string
		if err := rows.Scan(&shareCode); err != nil {
			return err
		}
		duplicateCodes = append(duplicateCodes, shareCode)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	for _, shareCode := range duplicateCodes {
		var keepID int64
		err := s.db.QueryRow(`SELECT id FROM tune_harvest_candidate
			WHERE share_code = ?
			ORDER BY
				CASE status WHEN 'imported' THEN 0 WHEN 'pending' THEN 1 WHEN 'rejected' THEN 2 ELSE 3 END,
				COALESCE(match_score, 0) DESC,
				updated_at DESC,
				id DESC
			LIMIT 1`, shareCode).Scan(&keepID)
		if err != nil {
			return err
		}
		if _, err := s.db.Exec(`DELETE FROM tune_harvest_candidate WHERE share_code = ? AND id <> ?`, shareCode, keepID); err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) UpsertTuneHarvestCandidate(input TuneHarvestCandidateInput) (*TuneHarvestCandidate, error) {
	normalized, err := NormalizeTuneHarvestCandidateInput(input)
	if err != nil {
		return nil, err
	}
	now := nowText()
	if _, err := s.getTuneHarvestCandidateBySourceRawKey(normalized.Source, normalized.RawKey); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	} else if err == nil {
		return s.upsertTuneHarvestCandidateBySourceRawKey(normalized, now)
	}
	if existing, err := s.getTuneHarvestCandidateByShareCode(normalized.ShareCode); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	} else if err == nil {
		return s.mergeTuneHarvestCandidateDuplicate(existing.ID, normalized, now)
	}
	return s.upsertTuneHarvestCandidateBySourceRawKey(normalized, now)
}

func (s *Store) upsertTuneHarvestCandidateBySourceRawKey(normalized TuneHarvestCandidateInput, now string) (*TuneHarvestCandidate, error) {
	_, err := s.db.Exec(`INSERT INTO tune_harvest_candidate (
		run_id, source, source_ref, source_url, source_car_id, raw_key, share_code, year, make, model, car_name,
		matched_car_id, match_score, match_reason, use_case, car_class, pi, drivetrain, tire_compound, tuner,
		tune_name, best_for, difficulty, notes, raw_json, status, rejection_reason, created_at, updated_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(source, raw_key) DO UPDATE SET
		run_id = excluded.run_id,
		source_ref = excluded.source_ref,
		source_url = excluded.source_url,
		source_car_id = excluded.source_car_id,
		share_code = excluded.share_code,
		year = excluded.year,
		make = excluded.make,
		model = excluded.model,
		car_name = excluded.car_name,
		matched_car_id = excluded.matched_car_id,
		match_score = excluded.match_score,
		match_reason = excluded.match_reason,
		use_case = excluded.use_case,
		car_class = excluded.car_class,
		pi = excluded.pi,
		drivetrain = excluded.drivetrain,
		tire_compound = excluded.tire_compound,
		tuner = excluded.tuner,
		tune_name = excluded.tune_name,
		best_for = excluded.best_for,
		difficulty = excluded.difficulty,
		notes = excluded.notes,
		raw_json = excluded.raw_json,
		status = CASE WHEN tune_harvest_candidate.status IN ('imported', 'rejected') THEN tune_harvest_candidate.status ELSE excluded.status END,
		rejection_reason = CASE WHEN tune_harvest_candidate.status IN ('imported', 'rejected') THEN tune_harvest_candidate.rejection_reason ELSE excluded.rejection_reason END,
		updated_at = excluded.updated_at`,
		nullInt64(normalized.RunID), normalized.Source, normalized.SourceRef, normalized.SourceURL, normalized.SourceCarID, normalized.RawKey,
		normalized.ShareCode, nullInt(normalized.Year), normalized.Make, normalized.Model, normalized.CarName,
		emptyStringAsNil(normalized.MatchedCarID), normalized.MatchScore, normalized.MatchReason, normalized.UseCase,
		normalized.CarClass, nullInt(normalized.PI), normalized.Drivetrain, normalized.TireCompound, normalized.Tuner,
		normalized.TuneName, normalized.BestFor, normalized.Difficulty, normalized.Notes, normalized.RawJSON,
		normalized.Status, normalized.RejectionReason, now, now)
	if err != nil {
		return nil, err
	}
	return s.getTuneHarvestCandidateBySourceRawKey(normalized.Source, normalized.RawKey)
}

func (s *Store) mergeTuneHarvestCandidateDuplicate(id int64, normalized TuneHarvestCandidateInput, now string) (*TuneHarvestCandidate, error) {
	_, err := s.db.Exec(`UPDATE tune_harvest_candidate SET
		run_id = COALESCE(?, run_id),
		source_ref = COALESCE(NULLIF(source_ref, ''), NULLIF(?, ''), ''),
		source_url = COALESCE(NULLIF(source_url, ''), NULLIF(?, ''), ''),
		source_car_id = COALESCE(NULLIF(source_car_id, ''), NULLIF(?, ''), ''),
		year = CASE WHEN COALESCE(year, 0) = 0 THEN ? ELSE year END,
		make = COALESCE(NULLIF(make, ''), NULLIF(?, ''), ''),
		model = COALESCE(NULLIF(model, ''), NULLIF(?, ''), ''),
		car_name = COALESCE(NULLIF(car_name, ''), NULLIF(?, ''), ''),
		matched_car_id = CASE WHEN ? > COALESCE(match_score, 0) THEN NULLIF(?, '') ELSE matched_car_id END,
		match_score = CASE WHEN ? > COALESCE(match_score, 0) THEN ? ELSE match_score END,
		match_reason = CASE WHEN ? > COALESCE(match_score, 0) THEN ? ELSE COALESCE(NULLIF(match_reason, ''), NULLIF(?, ''), '') END,
		use_case = COALESCE(NULLIF(use_case, ''), NULLIF(?, ''), ''),
		car_class = COALESCE(NULLIF(car_class, ''), NULLIF(?, ''), ''),
		pi = CASE WHEN COALESCE(pi, 0) = 0 THEN ? ELSE pi END,
		drivetrain = COALESCE(NULLIF(drivetrain, ''), NULLIF(?, ''), ''),
		tire_compound = COALESCE(NULLIF(tire_compound, ''), NULLIF(?, ''), ''),
		tuner = COALESCE(NULLIF(tuner, ''), NULLIF(?, ''), ''),
		tune_name = COALESCE(NULLIF(tune_name, ''), NULLIF(?, ''), ''),
		best_for = COALESCE(NULLIF(best_for, ''), NULLIF(?, ''), ''),
		difficulty = COALESCE(NULLIF(difficulty, ''), NULLIF(?, ''), ''),
		notes = COALESCE(NULLIF(notes, ''), NULLIF(?, ''), ''),
		raw_json = COALESCE(NULLIF(raw_json, '{}'), NULLIF(?, ''), raw_json),
		status = CASE WHEN status IN ('imported', 'rejected') THEN status WHEN ? = 'rejected' THEN status ELSE ? END,
		rejection_reason = CASE WHEN status IN ('imported', 'rejected') THEN rejection_reason WHEN ? = 'rejected' THEN rejection_reason ELSE ? END,
		updated_at = ?
		WHERE id = ?`,
		nullInt64(normalized.RunID), normalized.SourceRef, normalized.SourceURL, normalized.SourceCarID,
		nullInt(normalized.Year), normalized.Make, normalized.Model, normalized.CarName,
		normalized.MatchScore, normalized.MatchedCarID, normalized.MatchScore, normalized.MatchScore,
		normalized.MatchScore, normalized.MatchReason, normalized.MatchReason,
		normalized.UseCase, normalized.CarClass, nullInt(normalized.PI), normalized.Drivetrain, normalized.TireCompound,
		normalized.Tuner, normalized.TuneName, normalized.BestFor, normalized.Difficulty, normalized.Notes, normalized.RawJSON,
		normalized.Status, normalized.Status, normalized.Status, normalized.RejectionReason, now, id)
	if err != nil {
		return nil, err
	}
	return s.getTuneHarvestCandidateByID(id)
}

func (s *Store) getTuneHarvestCandidateBySourceRawKey(source string, rawKey string) (*TuneHarvestCandidate, error) {
	row := s.db.QueryRow(`SELECT `+tuneHarvestCandidateColumns+` FROM tune_harvest_candidate WHERE source = ? AND raw_key = ?`, source, rawKey)
	candidate, err := scanTuneHarvestCandidate(row)
	if err != nil {
		return nil, err
	}
	return &candidate, nil
}

func (s *Store) getTuneHarvestCandidateByShareCode(shareCode string) (*TuneHarvestCandidate, error) {
	row := s.db.QueryRow(`SELECT `+tuneHarvestCandidateColumns+` FROM tune_harvest_candidate WHERE share_code = ? ORDER BY id ASC LIMIT 1`, shareCode)
	candidate, err := scanTuneHarvestCandidate(row)
	if err != nil {
		return nil, err
	}
	return &candidate, nil
}

func (s *Store) getTuneHarvestCandidateByID(id int64) (*TuneHarvestCandidate, error) {
	row := s.db.QueryRow(`SELECT `+tuneHarvestCandidateColumns+` FROM tune_harvest_candidate WHERE id = ?`, id)
	candidate, err := scanTuneHarvestCandidate(row)
	if err != nil {
		return nil, err
	}
	return &candidate, nil
}

func NormalizeTuneHarvestCandidateInput(input TuneHarvestCandidateInput) (TuneHarvestCandidateInput, error) {
	input.Source = strings.ToLower(strings.TrimSpace(input.Source))
	input.SourceRef = strings.TrimSpace(input.SourceRef)
	input.SourceURL = strings.TrimSpace(input.SourceURL)
	input.SourceCarID = strings.TrimSpace(input.SourceCarID)
	input.RawKey = strings.TrimSpace(input.RawKey)
	input.ShareCode = NormalizeTuneShareCode(input.ShareCode)
	input.Make = strings.TrimSpace(input.Make)
	input.Model = strings.TrimSpace(input.Model)
	input.CarName = strings.TrimSpace(input.CarName)
	input.MatchedCarID = strings.TrimSpace(input.MatchedCarID)
	input.MatchReason = strings.TrimSpace(input.MatchReason)
	input.UseCase = strings.TrimSpace(input.UseCase)
	input.CarClass = strings.ToUpper(strings.TrimSpace(input.CarClass))
	input.Drivetrain = strings.ToUpper(strings.TrimSpace(input.Drivetrain))
	input.TireCompound = strings.TrimSpace(input.TireCompound)
	input.Tuner = strings.TrimSpace(input.Tuner)
	input.TuneName = strings.TrimSpace(input.TuneName)
	input.BestFor = strings.TrimSpace(input.BestFor)
	input.Difficulty = strings.TrimSpace(input.Difficulty)
	input.Notes = strings.TrimSpace(input.Notes)
	input.Status = strings.TrimSpace(input.Status)
	input.RejectionReason = strings.TrimSpace(input.RejectionReason)
	if input.Source == "" {
		return input, errors.New("source is required")
	}
	if input.ShareCode == "" {
		return input, errors.New("share code is required")
	}
	if input.RawKey == "" {
		input.RawKey = input.ShareCode
	}
	if input.RawJSON == "" {
		input.RawJSON = "{}"
	}
	if input.Status == "" {
		input.Status = TuneHarvestCandidatePending
	}
	return input, nil
}

var tuneShareCodePattern = regexp.MustCompile(`(?m)(?:^|[^\d])(\d{3})[ -]?(\d{3})[ -]?(\d{3})(?:[^\d]|$)`)

func NormalizeTuneShareCode(value string) string {
	match := tuneShareCodePattern.FindStringSubmatch(strings.TrimSpace(value))
	if len(match) != 4 {
		return ""
	}
	return match[1] + match[2] + match[3]
}

func FormatTuneShareCode(value string) string {
	code := NormalizeTuneShareCode(value)
	if code == "" {
		return strings.TrimSpace(value)
	}
	return code[:3] + " " + code[3:6] + " " + code[6:]
}

func (s *Store) ListTuneHarvestCandidates(status string, limit int) ([]TuneHarvestCandidate, error) {
	return s.SearchTuneHarvestCandidates(status, "", limit)
}

func (s *Store) SearchTuneHarvestCandidates(status string, search string, limit int) ([]TuneHarvestCandidate, error) {
	status = strings.TrimSpace(status)
	search = strings.TrimSpace(search)
	if limit <= 0 || limit > 5000 {
		limit = 5000
	}
	query := `SELECT ` + tuneHarvestCandidateColumns + ` FROM tune_harvest_candidate`
	args := []any{}
	filters := []string{}
	if status != "" && status != "all" {
		filters = append(filters, `status = ?`)
		args = append(args, status)
	}
	for _, term := range tuneHarvestSearchTerms(search) {
		condition, conditionArgs := tuneHarvestSearchCondition(term)
		filters = append(filters, condition)
		args = append(args, conditionArgs...)
	}
	if len(filters) > 0 {
		query += ` WHERE ` + strings.Join(filters, ` AND `)
	}
	query += ` ORDER BY updated_at DESC, id DESC LIMIT ?`
	args = append(args, limit)
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	candidates := []TuneHarvestCandidate{}
	for rows.Next() {
		candidate, err := scanTuneHarvestCandidate(rows)
		if err != nil {
			return nil, err
		}
		candidates = append(candidates, candidate)
	}
	return candidates, rows.Err()
}

func tuneHarvestSearchTerms(search string) []string {
	search = strings.ToLower(strings.TrimSpace(search))
	if search == "" {
		return nil
	}
	seen := map[string]bool{}
	terms := []string{}
	for _, term := range strings.Fields(search) {
		term = strings.TrimSpace(term)
		if term == "" || seen[term] {
			continue
		}
		seen[term] = true
		terms = append(terms, term)
	}
	if len(terms) == 0 && !seen[search] {
		terms = append(terms, search)
	}
	return terms
}

func tuneHarvestSearchCondition(term string) (string, []any) {
	columns := []string{
		"source", "source_ref", "source_url", "source_car_id", "raw_key", "share_code",
		"CAST(year AS TEXT)", "make", "model", "car_name", "matched_car_id", "match_reason",
		"use_case", "car_class", "CAST(pi AS TEXT)", "drivetrain", "tire_compound", "tuner",
		"tune_name", "best_for", "difficulty", "notes", "status", "rejection_reason",
	}
	like := "%" + term + "%"
	parts := make([]string, 0, len(columns)+1)
	args := make([]any, 0, len(columns)+1)
	for _, column := range columns {
		parts = append(parts, "LOWER(COALESCE("+column+", '')) LIKE ?")
		args = append(args, like)
	}
	if digits := tuneHarvestSearchDigits(term); digits != "" {
		parts = append(parts, "share_code LIKE ?")
		args = append(args, "%"+digits+"%")
	}
	return "(" + strings.Join(parts, " OR ") + ")", args
}

func tuneHarvestSearchDigits(value string) string {
	var builder strings.Builder
	for _, r := range value {
		if r >= '0' && r <= '9' {
			builder.WriteRune(r)
		}
	}
	return builder.String()
}

func (s *Store) UpdateTuneHarvestCandidateStatus(id int64, status string, reason string) (*TuneHarvestCandidate, error) {
	status = strings.TrimSpace(status)
	if status == "" {
		return nil, errors.New("status is required")
	}
	if status != TuneHarvestCandidatePending && status != TuneHarvestCandidateRejected && status != TuneHarvestCandidateImported {
		return nil, fmt.Errorf("unsupported candidate status %q", status)
	}
	result, err := s.db.Exec(`UPDATE tune_harvest_candidate SET status = ?, rejection_reason = ?, updated_at = ? WHERE id = ?`,
		status, strings.TrimSpace(reason), nowText(), id)
	if err != nil {
		return nil, err
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return nil, fmt.Errorf("tune harvest candidate %d not found", id)
	}
	row := s.db.QueryRow(`SELECT `+tuneHarvestCandidateColumns+` FROM tune_harvest_candidate WHERE id = ?`, id)
	candidate, err := scanTuneHarvestCandidate(row)
	if err != nil {
		return nil, err
	}
	return &candidate, nil
}

func (s *Store) ClearTuneHarvestCandidates() (int64, error) {
	result, err := s.db.Exec(`DELETE FROM tune_harvest_candidate`)
	if err != nil {
		return 0, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return affected, nil
}

func cleanTuneHarvestSources(sources []string) []string {
	seen := map[string]bool{}
	out := []string{}
	for _, source := range sources {
		source = strings.ToLower(strings.TrimSpace(source))
		if source == "" || seen[source] {
			continue
		}
		seen[source] = true
		out = append(out, source)
	}
	return out
}

type fh6CarScanner interface {
	Scan(dest ...any) error
}

func scanFH6Car(scanner fh6CarScanner) (FH6Car, error) {
	var car FH6Car
	var aliasJSON string
	err := scanner.Scan(&car.CarID, &car.Year, &car.Make, &car.Model, &aliasJSON, &car.BasePI,
		&car.DrivetrainDefault, &car.Source, &car.SourceRef, &car.CreatedAt, &car.UpdatedAt)
	if err != nil {
		return car, err
	}
	_ = json.Unmarshal([]byte(aliasJSON), &car.Aliases)
	if car.Aliases == nil {
		car.Aliases = []string{}
	}
	return car, nil
}

type tuneHarvestRunScanner interface {
	Scan(dest ...any) error
}

func scanTuneHarvestRun(scanner tuneHarvestRunScanner) (TuneHarvestRun, error) {
	var run TuneHarvestRun
	var sourcesJSON string
	var dryRun int
	err := scanner.Scan(&run.ID, &run.StartedAt, &run.FinishedAt, &sourcesJSON, &dryRun, &run.Status, &run.Message,
		&run.FoundCount, &run.SavedCount, &run.RejectedCount, &run.PendingCount, &run.ImportedCount)
	if err != nil {
		return run, err
	}
	_ = json.Unmarshal([]byte(sourcesJSON), &run.Sources)
	run.DryRun = dryRun != 0
	if run.Sources == nil {
		run.Sources = []string{}
	}
	return run, nil
}

const tuneHarvestCandidateColumns = `id, COALESCE(run_id, 0), source, COALESCE(source_ref, ''), COALESCE(source_url, ''), COALESCE(source_car_id, ''),
	raw_key, share_code, COALESCE(year, 0), COALESCE(make, ''), COALESCE(model, ''), COALESCE(car_name, ''),
	COALESCE(matched_car_id, ''), COALESCE(match_score, 0), COALESCE(match_reason, ''), COALESCE(use_case, ''),
	COALESCE(car_class, ''), COALESCE(pi, 0), COALESCE(drivetrain, ''), COALESCE(tire_compound, ''), COALESCE(tuner, ''),
	COALESCE(tune_name, ''), COALESCE(best_for, ''), COALESCE(difficulty, ''), COALESCE(notes, ''), COALESCE(raw_json, '{}'),
	status, COALESCE(rejection_reason, ''), created_at, updated_at`

type tuneHarvestCandidateScanner interface {
	Scan(dest ...any) error
}

func scanTuneHarvestCandidate(scanner tuneHarvestCandidateScanner) (TuneHarvestCandidate, error) {
	var candidate TuneHarvestCandidate
	err := scanner.Scan(
		&candidate.ID, &candidate.RunID, &candidate.Source, &candidate.SourceRef, &candidate.SourceURL, &candidate.SourceCarID,
		&candidate.RawKey, &candidate.ShareCode, &candidate.Year, &candidate.Make, &candidate.Model, &candidate.CarName,
		&candidate.MatchedCarID, &candidate.MatchScore, &candidate.MatchReason, &candidate.UseCase,
		&candidate.CarClass, &candidate.PI, &candidate.Drivetrain, &candidate.TireCompound, &candidate.Tuner,
		&candidate.TuneName, &candidate.BestFor, &candidate.Difficulty, &candidate.Notes, &candidate.RawJSON,
		&candidate.Status, &candidate.RejectionReason, &candidate.CreatedAt, &candidate.UpdatedAt,
	)
	return candidate, err
}

func nullInt(value int) any {
	if value == 0 {
		return nil
	}
	return value
}

func nullInt64(value int64) any {
	if value == 0 {
		return nil
	}
	return value
}

func emptyStringAsNil(value string) any {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	return value
}
