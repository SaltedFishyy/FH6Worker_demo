package storage

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"regexp"
	"strings"
	"unicode"
)

type RecommendedCarInput struct {
	ID                string   `json:"id"`
	Name              string   `json:"name"`
	UseCase           string   `json:"useCase"`
	UseCaseLabel      string   `json:"useCaseLabel"`
	PI                int      `json:"pi"`
	CarClass          string   `json:"carClass"`
	Drivetrain        string   `json:"drivetrain"`
	TireCompound      string   `json:"tireCompound"`
	TireCompoundLabel string   `json:"tireCompoundLabel"`
	WeightKG          float64  `json:"weightKG"`
	FrontWeightPct    float64  `json:"frontWeightPct"`
	TuneCode          string   `json:"tuneCode"`
	ImageSrc          string   `json:"imageSrc,omitempty"`
	Tags              []string `json:"tags"`
	Reason            string   `json:"reason"`
}

type RecommendedCar struct {
	RecommendedCarInput
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

func NormalizeRecommendedCarInput(input RecommendedCarInput) (RecommendedCarInput, error) {
	input.ID = strings.TrimSpace(input.ID)
	input.Name = strings.TrimSpace(input.Name)
	input.UseCase = strings.TrimSpace(input.UseCase)
	input.CarClass = strings.ToUpper(strings.TrimSpace(input.CarClass))
	input.Drivetrain = strings.ToUpper(strings.TrimSpace(input.Drivetrain))
	input.TireCompound = strings.TrimSpace(input.TireCompound)
	input.TuneCode = NormalizeTuneShareCode(input.TuneCode)
	input.ImageSrc = strings.TrimSpace(input.ImageSrc)
	input.Reason = strings.TrimSpace(input.Reason)
	if input.Name == "" {
		return input, errors.New("name is required")
	}
	useCase, useCaseLabel, ok := normalizeRecommendedUseCase(input.UseCase)
	if !ok {
		return input, errors.New("unsupported use case")
	}
	input.UseCase = useCase
	input.UseCaseLabel = useCaseLabel
	if input.PI < 100 || input.PI > 999 {
		return input, errors.New("pi must be between 100 and 999")
	}
	if input.CarClass == "" {
		return input, errors.New("car class is required")
	}
	if input.Drivetrain == "" {
		return input, errors.New("drivetrain is required")
	}
	tireCompound, tireCompoundLabel, ok := normalizeRecommendedTireCompound(input.TireCompound)
	if !ok {
		return input, errors.New("unsupported tire compound")
	}
	input.TireCompound = tireCompound
	input.TireCompoundLabel = tireCompoundLabel
	if !finiteNumber(input.WeightKG) || input.WeightKG < 0 {
		return input, errors.New("weightKG must be a non-negative number")
	}
	if !finiteNumber(input.FrontWeightPct) || input.FrontWeightPct < 0 || input.FrontWeightPct >= 100 {
		return input, errors.New("frontWeightPct must be empty or between 1 and 99")
	}
	if input.TuneCode == "" {
		return input, errors.New("tuneCode is required")
	}
	tags := make([]string, 0, len(input.Tags))
	for _, tag := range input.Tags {
		tag = strings.TrimSpace(tag)
		if tag != "" {
			tags = append(tags, tag)
		}
	}
	input.Tags = tags
	input.ID = GenerateRecommendedCarID(input)
	return input, nil
}

func finiteNumber(value float64) bool {
	return !math.IsNaN(value) && !math.IsInf(value, 0)
}

func normalizeRecommendedUseCase(value string) (string, string, bool) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "road", "公路":
		return "Road", "公路", true
	case "rally", "拉力":
		return "Rally", "拉力", true
	case "offroad", "越野":
		return "Offroad", "越野", true
	case "drift", "漂移":
		return "Drift", "漂移", true
	case "drag", "直线":
		return "Drag", "直线", true
	default:
		return "", "", false
	}
}

func normalizeRecommendedTireCompound(value string) (string, string, bool) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "stock", "原厂", "原厂/街胎":
		return "stock", "原厂", true
	case "street", "街胎":
		return "street", "街胎", true
	case "sport", "运动":
		return "sport", "运动", true
	case "semi", "半热熔":
		return "semi", "半热熔", true
	case "slick", "热熔胎", "光头胎":
		return "slick", "热熔胎", true
	case "rally", "拉力":
		return "rally", "拉力", true
	case "offroad", "越野":
		return "offroad", "越野", true
	case "drift", "漂移":
		return "drift", "漂移", true
	case "drag", "直线":
		return "drag", "直线", true
	case "snow", "雪地":
		return "snow", "雪地", true
	default:
		return "", "", false
	}
}

func GenerateRecommendedCarID(input RecommendedCarInput) string {
	namePart := slugifyRecommendedCarName(input.Name)
	if namePart == "" {
		namePart = "car"
	}
	useCase := strings.ToLower(input.UseCase)
	classPI := strings.ToLower(input.CarClass) + fmt.Sprintf("%d", input.PI)
	parts := []string{namePart, useCase, classPI}
	if code := NormalizeTuneShareCode(input.TuneCode); code != "" {
		parts = append(parts, code)
	}
	return strings.Join(parts, "-")
}

var yearPattern = regexp.MustCompile(`^(19|20)\d{2}$`)

func slugifyRecommendedCarName(name string) string {
	words := strings.Fields(strings.TrimSpace(name))
	if len(words) > 1 && yearPattern.MatchString(words[0]) {
		words = append(words[1:], words[0])
	}
	return slugifyASCII(strings.Join(words, " "))
}

func slugifyASCII(value string) string {
	var b strings.Builder
	lastHyphen := false
	for _, r := range strings.ToLower(value) {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
			lastHyphen = false
			continue
		}
		if unicode.IsSpace(r) || r == '-' || r == '_' || r == '/' {
			if !lastHyphen && b.Len() > 0 {
				b.WriteRune('-')
				lastHyphen = true
			}
		}
	}
	return strings.Trim(b.String(), "-")
}

func (s *Store) ListRecommendedCars() ([]RecommendedCar, error) {
	rows, err := s.db.Query(`SELECT id, name, use_case, use_case_label, pi, car_class, drivetrain, tire_compound, tire_compound_label,
		weight_kg, front_weight_pct, tune_code, COALESCE(image_src, ''), COALESCE(tags_json, '[]'), COALESCE(reason, ''), created_at, updated_at
		FROM recommended_car ORDER BY use_case, car_class, pi DESC, name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	cars := []RecommendedCar{}
	for rows.Next() {
		car, err := scanRecommendedCar(rows)
		if err != nil {
			return nil, err
		}
		cars = append(cars, car)
	}
	return cars, rows.Err()
}

func (s *Store) GetRecommendedCar(id string) (*RecommendedCar, error) {
	row := s.db.QueryRow(`SELECT id, name, use_case, use_case_label, pi, car_class, drivetrain, tire_compound, tire_compound_label,
		weight_kg, front_weight_pct, tune_code, COALESCE(image_src, ''), COALESCE(tags_json, '[]'), COALESCE(reason, ''), created_at, updated_at
		FROM recommended_car WHERE id = ?`, strings.TrimSpace(id))
	car, err := scanRecommendedCar(row)
	if err != nil {
		return nil, err
	}
	return &car, nil
}

func (s *Store) SaveRecommendedCar(input RecommendedCarInput) (*RecommendedCar, error) {
	previousID := strings.TrimSpace(input.ID)
	if previousID != "" {
		var count int
		if err := s.db.QueryRow(`SELECT COUNT(*) FROM recommended_car WHERE id = ?`, previousID).Scan(&count); err != nil {
			return nil, err
		}
		if count == 0 {
			previousID = ""
		}
	}
	return s.SaveRecommendedCarRecord(input, previousID)
}

func (s *Store) SaveRecommendedCarRecord(input RecommendedCarInput, previousID string) (*RecommendedCar, error) {
	normalized, err := NormalizeRecommendedCarInput(input)
	if err != nil {
		return nil, err
	}
	tagsJSON, err := json.Marshal(normalized.Tags)
	if err != nil {
		return nil, err
	}
	previousID = strings.TrimSpace(previousID)
	now := nowText()

	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback()
		}
	}()

	createdAt := now
	if previousID != "" {
		if err := tx.QueryRow(`SELECT created_at FROM recommended_car WHERE id = ?`, previousID).Scan(&createdAt); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, fmt.Errorf("recommended car %q not found", previousID)
			}
			return nil, err
		}
	}
	if previousID == "" || normalized.ID != previousID {
		var existingID string
		err := tx.QueryRow(`SELECT id FROM recommended_car WHERE id = ?`, normalized.ID).Scan(&existingID)
		if err == nil {
			return nil, fmt.Errorf("recommended car id %q already exists", normalized.ID)
		}
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}
	var tuneCodeOwner string
	err = tx.QueryRow(`SELECT id FROM recommended_car WHERE tune_code = ? LIMIT 1`, normalized.TuneCode).Scan(&tuneCodeOwner)
	if err == nil && tuneCodeOwner != previousID {
		return nil, fmt.Errorf("tuneCode %q already exists", normalized.TuneCode)
	}
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	if previousID != "" && normalized.ID != previousID {
		if _, err := tx.Exec(`DELETE FROM recommended_car WHERE id = ?`, previousID); err != nil {
			return nil, err
		}
	}

	_, err = tx.Exec(`INSERT INTO recommended_car (
		id, name, use_case, use_case_label, pi, car_class, drivetrain, tire_compound, tire_compound_label,
		weight_kg, front_weight_pct, tune_code, image_src, tags_json, reason, created_at, updated_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(id) DO UPDATE SET
		name = excluded.name,
		use_case = excluded.use_case,
		use_case_label = excluded.use_case_label,
		pi = excluded.pi,
		car_class = excluded.car_class,
		drivetrain = excluded.drivetrain,
		tire_compound = excluded.tire_compound,
		tire_compound_label = excluded.tire_compound_label,
		weight_kg = excluded.weight_kg,
		front_weight_pct = excluded.front_weight_pct,
		tune_code = excluded.tune_code,
		image_src = excluded.image_src,
		tags_json = excluded.tags_json,
		reason = excluded.reason,
		updated_at = excluded.updated_at`,
		normalized.ID, normalized.Name, normalized.UseCase, normalized.UseCaseLabel, normalized.PI, normalized.CarClass, normalized.Drivetrain,
		normalized.TireCompound, normalized.TireCompoundLabel, normalized.WeightKG, normalized.FrontWeightPct, normalized.TuneCode,
		normalized.ImageSrc, string(tagsJSON), normalized.Reason, createdAt, now)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	committed = true
	return s.GetRecommendedCar(normalized.ID)
}

func (s *Store) DeleteRecommendedCar(id string) error {
	result, err := s.db.Exec(`DELETE FROM recommended_car WHERE id = ?`, strings.TrimSpace(id))
	if err != nil {
		return err
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("recommended car %q not found", strings.TrimSpace(id))
	}
	return nil
}

func (s *Store) DeleteAllRecommendedCars() (int64, error) {
	result, err := s.db.Exec(`DELETE FROM recommended_car`)
	if err != nil {
		return 0, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return affected, nil
}

type recommendedCarScanner interface {
	Scan(dest ...any) error
}

func scanRecommendedCar(scanner recommendedCarScanner) (RecommendedCar, error) {
	var car RecommendedCar
	var tagsJSON string
	err := scanner.Scan(
		&car.ID, &car.Name, &car.UseCase, &car.UseCaseLabel, &car.PI, &car.CarClass, &car.Drivetrain,
		&car.TireCompound, &car.TireCompoundLabel, &car.WeightKG, &car.FrontWeightPct, &car.TuneCode,
		&car.ImageSrc, &tagsJSON, &car.Reason, &car.CreatedAt, &car.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return car, err
		}
		return car, err
	}
	if strings.TrimSpace(tagsJSON) != "" {
		_ = json.Unmarshal([]byte(tagsJSON), &car.Tags)
	}
	if car.Tags == nil {
		car.Tags = []string{}
	}
	return car, nil
}
