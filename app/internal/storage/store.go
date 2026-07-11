package storage

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"fh6worker/internal/telemetry"

	_ "modernc.org/sqlite"
)

type Store struct {
	db              *sql.DB
	path            string
	knowledgeMu     sync.RWMutex
	roadKnowledge   *RoadTuningKnowledge
	knowledgeStatus RoadTuningKnowledgeStatus
}

const (
	benchmarkTrackTypeAuto    = "auto"
	benchmarkTrackTypeCircuit = "circuit"
	benchmarkTrackTypeSprint  = "sprint"

	benchmarkExtractionAutoBestLap = "auto_best_lap"
	benchmarkExtractionFirstLap    = "first_lap"
	benchmarkExtractionFullSegment = "full_segment"

	defaultGateWidthMeters = 30.0
	defaultGateDepthMeters = 20.0
	defaultGateRadius      = 20.0
	circuitAutoDistance    = 40.0
	minCircuitRouteMeters  = 200.0
	minCircuitDurationMS   = int64(30000)
)

func OpenDefault() (*Store, error) {
	primary, primaryErr := defaultDatabasePath()
	if primaryErr == nil {
		store, err := Open(primary)
		if err == nil {
			return store, nil
		}
	}
	fallback := filepath.Join("data", "fh6worker.db")
	return Open(fallback)
}

func Open(path string) (*Store, error) {
	if strings.TrimSpace(path) == "" {
		return nil, errors.New("database path is required")
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	store := &Store{db: db, path: path}
	if err := store.Migrate(); err != nil {
		_ = db.Close()
		return nil, err
	}
	_ = store.ReloadTuningKnowledge()
	return store, nil
}

func (s *Store) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	return s.db.Close()
}

func (s *Store) Path() string {
	if s == nil {
		return ""
	}
	return s.path
}

func (s *Store) CleanupLegacySessions() error {
	if s == nil || s.db == nil {
		return nil
	}
	rows, err := s.db.Query(`SELECT COALESCE(recording_path, '') FROM telemetry_session WHERE COALESCE(recording_path, '') <> ''`)
	if err != nil {
		return err
	}
	recordings := []string{}
	for rows.Next() {
		var path string
		if err := rows.Scan(&path); err != nil {
			_ = rows.Close()
			return err
		}
		if strings.TrimSpace(path) != "" {
			recordings = append(recordings, path)
		}
	}
	if err := rows.Close(); err != nil {
		return err
	}
	for _, path := range recordings {
		if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
			return err
		}
	}
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	if _, err := tx.Exec(`UPDATE tune_change_log SET session_id = NULL WHERE session_id IS NOT NULL`); err != nil {
		_ = tx.Rollback()
		return err
	}
	for _, statement := range []string{
		`DELETE FROM detected_event`,
		`DELETE FROM telemetry_sample_agg`,
		`DELETE FROM benchmark_run`,
		`DELETE FROM telemetry_session`,
	} {
		if _, err := tx.Exec(statement); err != nil {
			_ = tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func defaultDatabasePath() (string, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, "FH6Worker", "fh6worker.db"), nil
}

func (s *Store) Migrate() error {
	statements := []string{
		`PRAGMA foreign_keys = ON`,
		`PRAGMA busy_timeout = 5000`,
		`CREATE TABLE IF NOT EXISTS app_setting (
			key TEXT PRIMARY KEY,
			value TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS tune_profile (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			car_name TEXT NOT NULL,
			car_ordinal INTEGER,
			car_category INTEGER,
			car_class TEXT,
			pi INTEGER,
			drivetrain TEXT,
			num_cylinders INTEGER,
			use_case TEXT,
			version_name TEXT,
			power_kw REAL,
			torque_nm REAL,
			weight_kg REAL,
			front_weight_pct REAL,
			power_to_weight_kw_per_kg REAL,
			peak_torque_rpm REAL,
			peak_power_rpm REAL,
			redline_rpm REAL,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL,
			front_tire_pressure REAL,
			rear_tire_pressure REAL,
			final_drive REAL,
			gear_1 REAL,
			gear_2 REAL,
			gear_3 REAL,
			gear_4 REAL,
			gear_5 REAL,
			gear_6 REAL,
			gear_7 REAL,
			gear_8 REAL,
			gear_9 REAL,
			gear_10 REAL,
			front_camber REAL,
			rear_camber REAL,
			front_toe REAL,
			rear_toe REAL,
			caster REAL,
			front_arb REAL,
			rear_arb REAL,
			front_spring REAL,
			rear_spring REAL,
			front_ride_height REAL,
			rear_ride_height REAL,
			front_rebound REAL,
			rear_rebound REAL,
			front_bump REAL,
			rear_bump REAL,
			front_aero REAL,
			rear_aero REAL,
			aero_balance REAL,
			brake_balance REAL,
			brake_pressure REAL,
			front_diff_accel REAL,
			front_diff_decel REAL,
			rear_diff_accel REAL,
			rear_diff_decel REAL,
			center_diff_balance REAL,
			notes TEXT
		)`,
		`CREATE TABLE IF NOT EXISTS telemetry_session (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			tune_profile_id INTEGER,
			tune_snapshot_json TEXT,
			session_name TEXT,
			track_name TEXT,
			mode TEXT,
			game_mode TEXT,
			started_at TEXT,
			ended_at TEXT,
			duration_ms INTEGER,
			best_lap_ms INTEGER,
			avg_speed_kmh REAL,
			max_speed_kmh REAL,
			recording_path TEXT,
			recording_packets INTEGER DEFAULT 0,
			recording_bytes INTEGER DEFAULT 0,
			recording_truncated INTEGER DEFAULT 0,
			car_ordinal INTEGER,
			car_class TEXT,
			car_pi INTEGER,
			drivetrain TEXT,
			num_cylinders INTEGER,
			driver_mode TEXT,
			driver_mode_confidence REAL DEFAULT 0,
			driver_mode_evidence_json TEXT,
			brake_assist TEXT,
			steering_assist TEXT,
			traction_control TEXT,
			stability_control TEXT,
			shifting TEXT,
			launch_control TEXT,
			driver_feedback_json TEXT,
			notes TEXT,
			FOREIGN KEY (tune_profile_id) REFERENCES tune_profile(id) ON DELETE SET NULL
		)`,
		`CREATE TABLE IF NOT EXISTS detected_event (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			session_id INTEGER NOT NULL,
			event_type TEXT NOT NULL,
			severity TEXT NOT NULL,
			segment TEXT,
			start_ms INTEGER,
			end_ms INTEGER,
			duration_ms INTEGER,
			evidence_json TEXT NOT NULL,
			suggestion_json TEXT,
			created_at TEXT NOT NULL,
			FOREIGN KEY (session_id) REFERENCES telemetry_session(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS telemetry_sample_agg (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			session_id INTEGER NOT NULL,
			timestamp_ms INTEGER,
			game_mode TEXT,
			is_race_on INTEGER,
			speed_kmh REAL,
			rpm REAL,
			rpm_ratio REAL,
			gear INTEGER,
			throttle REAL,
			brake REAL,
			steer REAL,
			front_slip_ratio REAL,
			rear_slip_ratio REAL,
			front_combined_slip REAL,
			rear_combined_slip REAL,
			front_tire_temp REAL,
			rear_tire_temp REAL,
			front_suspension REAL,
			rear_suspension REAL,
			yaw_rate REAL,
			pitch_rate REAL,
			roll_rate REAL,
			speed_field_kmh REAL,
			velocity_speed_kmh REAL,
			speed_source TEXT,
			position_x REAL,
			position_y REAL,
			position_z REAL,
			distance_traveled REAL,
			best_lap REAL,
			last_lap REAL,
			current_lap REAL,
			current_race_time REAL,
			lap_number INTEGER,
			race_position INTEGER,
			driving_line REAL,
			car_ordinal INTEGER,
			car_category INTEGER,
			car_class TEXT,
			car_pi INTEGER,
			drivetrain TEXT,
			num_cylinders INTEGER,
			FOREIGN KEY (session_id) REFERENCES telemetry_session(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS car_identity (
			car_ordinal INTEGER PRIMARY KEY,
			car_name TEXT NOT NULL,
			updated_at TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS recommended_car (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			use_case TEXT NOT NULL,
			use_case_label TEXT NOT NULL,
			pi INTEGER NOT NULL,
			car_class TEXT NOT NULL,
			drivetrain TEXT NOT NULL,
			tire_compound TEXT NOT NULL,
			tire_compound_label TEXT NOT NULL,
			weight_kg REAL NOT NULL,
			front_weight_pct REAL NOT NULL,
			tune_code TEXT NOT NULL,
			image_src TEXT,
			tags_json TEXT NOT NULL,
			reason TEXT,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS fh6_car (
			car_id TEXT PRIMARY KEY,
			year INTEGER NOT NULL,
			make TEXT NOT NULL,
			model TEXT NOT NULL,
			alias_json TEXT NOT NULL,
			base_pi INTEGER,
			drivetrain_default TEXT,
			source TEXT,
			source_ref TEXT,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS tune_harvest_run (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			started_at TEXT NOT NULL,
			finished_at TEXT,
			sources_json TEXT NOT NULL,
			dry_run INTEGER DEFAULT 0,
			status TEXT NOT NULL,
			message TEXT,
			found_count INTEGER DEFAULT 0,
			saved_count INTEGER DEFAULT 0,
			rejected_count INTEGER DEFAULT 0,
			pending_count INTEGER DEFAULT 0,
			imported_count INTEGER DEFAULT 0
		)`,
		`CREATE TABLE IF NOT EXISTS tune_harvest_candidate (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			run_id INTEGER,
			source TEXT NOT NULL,
			source_ref TEXT,
			source_url TEXT,
			source_car_id TEXT,
			raw_key TEXT NOT NULL,
			share_code TEXT NOT NULL,
			year INTEGER,
			make TEXT,
			model TEXT,
			car_name TEXT,
			matched_car_id TEXT,
			match_score REAL DEFAULT 0,
			match_reason TEXT,
			use_case TEXT,
			car_class TEXT,
			pi INTEGER,
			drivetrain TEXT,
			tire_compound TEXT,
			tuner TEXT,
			tune_name TEXT,
			best_for TEXT,
			difficulty TEXT,
			notes TEXT,
			raw_json TEXT NOT NULL,
			status TEXT NOT NULL,
			rejection_reason TEXT,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL,
			FOREIGN KEY (run_id) REFERENCES tune_harvest_run(id) ON DELETE SET NULL,
			FOREIGN KEY (matched_car_id) REFERENCES fh6_car(car_id) ON DELETE SET NULL,
			UNIQUE(source, raw_key)
		)`,
		`CREATE TABLE IF NOT EXISTS tune_change_log (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			tune_profile_id INTEGER NOT NULL,
			session_id INTEGER,
			changed_at TEXT NOT NULL,
			change_reason TEXT,
			change_json TEXT NOT NULL,
			FOREIGN KEY (tune_profile_id) REFERENCES tune_profile(id) ON DELETE CASCADE,
			FOREIGN KEY (session_id) REFERENCES telemetry_session(id) ON DELETE SET NULL
		)`,
		`CREATE TABLE IF NOT EXISTS rule_threshold_profile (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			car_class TEXT,
			drivetrain TEXT,
			use_case TEXT,
			game_mode TEXT,
			config_json TEXT NOT NULL,
			is_default INTEGER DEFAULT 0,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS benchmark_track (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			source_mode TEXT,
			track_type TEXT,
			start_x REAL,
			start_y REAL,
			start_z REAL,
			end_x REAL,
			end_y REAL,
			end_z REAL,
			start_radius REAL,
			end_radius REAL,
			direction_x REAL,
			direction_z REAL,
			route_length_meters REAL,
			has_driving_line INTEGER DEFAULT 0,
			start_gate_json TEXT,
			finish_gate_json TEXT,
			checkpoints_json TEXT,
			source_session_id INTEGER,
			lap_count_observed INTEGER DEFAULT 0,
			polyline_json TEXT NOT NULL,
			notes TEXT,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS benchmark_run (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			session_id INTEGER NOT NULL,
			track_id INTEGER NOT NULL,
			start_ms INTEGER,
			end_ms INTEGER,
			duration_ms INTEGER,
			confidence REAL,
			avg_speed_kmh REAL,
			max_speed_kmh REAL,
			route_progress_01 REAL,
			geometry_length_meters REAL,
			track_length_error_pct REAL,
			distance_traveled_delta_meters REAL,
			current_race_time_delta_seconds REAL,
			avg_lateral_error_meters REAL,
			max_lateral_error_meters REAL,
			warning_flags TEXT,
			event_count INTEGER DEFAULT 0,
			driver_mode TEXT,
			driver_mode_confidence REAL DEFAULT 0,
			driver_mode_evidence_json TEXT,
			valid INTEGER DEFAULT 0,
			created_at TEXT NOT NULL,
			FOREIGN KEY (session_id) REFERENCES telemetry_session(id) ON DELETE CASCADE,
			FOREIGN KEY (track_id) REFERENCES benchmark_track(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS track_baseline_run (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			track_id INTEGER NOT NULL,
			start_ms INTEGER,
			end_ms INTEGER,
			duration_ms INTEGER,
			confidence REAL,
			avg_speed_kmh REAL,
			max_speed_kmh REAL,
			route_progress_01 REAL,
			geometry_length_meters REAL,
			track_length_error_pct REAL,
			distance_traveled_delta_meters REAL,
			current_race_time_delta_seconds REAL,
			avg_lateral_error_meters REAL,
			max_lateral_error_meters REAL,
			warning_flags TEXT,
			event_count INTEGER DEFAULT 0,
			driver_mode TEXT,
			driver_mode_confidence REAL DEFAULT 0,
			driver_mode_evidence_json TEXT,
			valid INTEGER DEFAULT 0,
			car_ordinal INTEGER,
			car_class TEXT,
			car_pi INTEGER,
			drivetrain TEXT,
			game_mode TEXT,
			created_at TEXT NOT NULL,
			FOREIGN KEY (track_id) REFERENCES benchmark_track(id) ON DELETE CASCADE
		)`,
	}
	for _, statement := range statements {
		if _, err := s.db.Exec(statement); err != nil {
			return err
		}
	}
	if err := s.dedupeTuneHarvestCandidatesByShareCode(); err != nil {
		return err
	}
	if _, err := s.db.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS idx_tune_harvest_candidate_share_code_unique ON tune_harvest_candidate(share_code)`); err != nil {
		return err
	}
	for _, column := range []struct {
		name       string
		definition string
	}{
		{name: "car_category", definition: "INTEGER"},
		{name: "num_cylinders", definition: "INTEGER"},
		{name: "power_kw", definition: "REAL"},
		{name: "torque_nm", definition: "REAL"},
		{name: "weight_kg", definition: "REAL"},
		{name: "front_weight_pct", definition: "REAL"},
		{name: "power_to_weight_kw_per_kg", definition: "REAL"},
		{name: "peak_torque_rpm", definition: "REAL"},
		{name: "peak_power_rpm", definition: "REAL"},
		{name: "redline_rpm", definition: "REAL"},
		{name: "gear_9", definition: "REAL"},
		{name: "gear_10", definition: "REAL"},
	} {
		if err := s.ensureColumn("tune_profile", column.name, column.definition); err != nil {
			return err
		}
	}
	for _, column := range []struct {
		name       string
		definition string
	}{
		{name: "recording_path", definition: "TEXT"},
		{name: "tune_snapshot_json", definition: "TEXT"},
		{name: "recording_packets", definition: "INTEGER DEFAULT 0"},
		{name: "recording_bytes", definition: "INTEGER DEFAULT 0"},
		{name: "recording_truncated", definition: "INTEGER DEFAULT 0"},
		{name: "game_mode", definition: "TEXT"},
		{name: "car_ordinal", definition: "INTEGER"},
		{name: "car_class", definition: "TEXT"},
		{name: "car_pi", definition: "INTEGER"},
		{name: "drivetrain", definition: "TEXT"},
		{name: "num_cylinders", definition: "INTEGER"},
		{name: "driver_mode", definition: "TEXT"},
		{name: "driver_mode_confidence", definition: "REAL DEFAULT 0"},
		{name: "driver_mode_evidence_json", definition: "TEXT"},
		{name: "brake_assist", definition: "TEXT"},
		{name: "steering_assist", definition: "TEXT"},
		{name: "traction_control", definition: "TEXT"},
		{name: "stability_control", definition: "TEXT"},
		{name: "shifting", definition: "TEXT"},
		{name: "launch_control", definition: "TEXT"},
		{name: "driver_feedback_json", definition: "TEXT"},
	} {
		if err := s.ensureColumn("telemetry_session", column.name, column.definition); err != nil {
			return err
		}
	}
	for _, column := range []struct {
		name       string
		definition string
	}{
		{name: "speed_field_kmh", definition: "REAL"},
		{name: "velocity_speed_kmh", definition: "REAL"},
		{name: "speed_source", definition: "TEXT"},
		{name: "game_mode", definition: "TEXT"},
		{name: "is_race_on", definition: "INTEGER"},
		{name: "position_x", definition: "REAL"},
		{name: "position_y", definition: "REAL"},
		{name: "position_z", definition: "REAL"},
		{name: "distance_traveled", definition: "REAL"},
		{name: "best_lap", definition: "REAL"},
		{name: "last_lap", definition: "REAL"},
		{name: "current_lap", definition: "REAL"},
		{name: "current_race_time", definition: "REAL"},
		{name: "lap_number", definition: "INTEGER"},
		{name: "race_position", definition: "INTEGER"},
		{name: "driving_line", definition: "REAL"},
		{name: "car_ordinal", definition: "INTEGER"},
		{name: "car_category", definition: "INTEGER"},
		{name: "car_class", definition: "TEXT"},
		{name: "car_pi", definition: "INTEGER"},
		{name: "drivetrain", definition: "TEXT"},
		{name: "num_cylinders", definition: "INTEGER"},
	} {
		if err := s.ensureColumn("telemetry_sample_agg", column.name, column.definition); err != nil {
			return err
		}
	}
	if err := s.ensureColumn("rule_threshold_profile", "game_mode", "TEXT"); err != nil {
		return err
	}
	for _, column := range []struct {
		name       string
		definition string
	}{
		{name: "track_type", definition: "TEXT"},
		{name: "start_gate_json", definition: "TEXT"},
		{name: "finish_gate_json", definition: "TEXT"},
		{name: "checkpoints_json", definition: "TEXT"},
		{name: "source_session_id", definition: "INTEGER"},
		{name: "lap_count_observed", definition: "INTEGER DEFAULT 0"},
	} {
		if err := s.ensureColumn("benchmark_track", column.name, column.definition); err != nil {
			return err
		}
	}
	if err := s.backfillBenchmarkTrackFields(); err != nil {
		return err
	}
	for _, column := range []struct {
		name       string
		definition string
	}{
		{name: "route_progress_01", definition: "REAL"},
		{name: "geometry_length_meters", definition: "REAL"},
		{name: "track_length_error_pct", definition: "REAL"},
		{name: "distance_traveled_delta_meters", definition: "REAL"},
		{name: "current_race_time_delta_seconds", definition: "REAL"},
		{name: "avg_lateral_error_meters", definition: "REAL"},
		{name: "max_lateral_error_meters", definition: "REAL"},
		{name: "warning_flags", definition: "TEXT"},
		{name: "driver_mode_confidence", definition: "REAL DEFAULT 0"},
		{name: "driver_mode_evidence_json", definition: "TEXT"},
	} {
		if err := s.ensureColumn("benchmark_run", column.name, column.definition); err != nil {
			return err
		}
	}
	if err := s.ensureDefaultRuleThresholdProfile(); err != nil {
		return err
	}
	if err := s.ensureRoadRacingRuleThresholdProfile(); err != nil {
		return err
	}
	return nil
}

func (s *Store) ensureColumn(table, column, definition string) error {
	rows, err := s.db.Query(`PRAGMA table_info(` + table + `)`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var cid int
		var name, columnType string
		var notNull int
		var defaultValue any
		var pk int
		if err := rows.Scan(&cid, &name, &columnType, &notNull, &defaultValue, &pk); err != nil {
			return err
		}
		if name == column {
			return rows.Err()
		}
	}
	if err := rows.Err(); err != nil {
		return err
	}
	_, err = s.db.Exec(fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s", table, column, definition))
	return err
}

func (s *Store) backfillBenchmarkTrackFields() error {
	tracks, err := s.ListBenchmarkTracks()
	if err != nil {
		return err
	}
	for _, track := range tracks {
		startGateJSON, finishGateJSON, checkpointsJSON, err := benchmarkTrackAuxJSON(track.BenchmarkTrackInput)
		if err != nil {
			return err
		}
		polylineJSON, err := json.Marshal(track.Polyline)
		if err != nil {
			return err
		}
		if _, err := s.db.Exec(`UPDATE benchmark_track SET
			track_type = ?, start_gate_json = ?, finish_gate_json = ?, checkpoints_json = ?,
			lap_count_observed = ?, route_length_meters = ?, polyline_json = ?, updated_at = updated_at
			WHERE id = ?`,
			track.TrackType, startGateJSON, finishGateJSON, checkpointsJSON, track.LapCountObserved, track.RouteLengthMeters, string(polylineJSON), track.ID,
		); err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) ListTuneProfiles() ([]TuneProfile, error) {
	rows, err := s.db.Query(`SELECT ` + profileSelectColumns + ` FROM tune_profile ORDER BY updated_at DESC, id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var profiles []TuneProfile
	for rows.Next() {
		profile, err := scanProfile(rows)
		if err != nil {
			return nil, err
		}
		profiles = append(profiles, *profile)
	}
	return profiles, rows.Err()
}

func (s *Store) GetTuneProfile(id int64) (*TuneProfile, error) {
	row := s.db.QueryRow(`SELECT `+profileSelectColumns+` FROM tune_profile WHERE id = ?`, id)
	return scanProfile(row)
}

func (s *Store) ListTuneProfilesForVehicle(carOrdinal int64, carClass string) ([]TuneProfile, error) {
	if carOrdinal <= 0 {
		return []TuneProfile{}, nil
	}
	carClass = strings.TrimSpace(carClass)
	query := `SELECT ` + profileSelectColumns + ` FROM tune_profile WHERE car_ordinal = ?`
	args := []any{carOrdinal}
	if carClass != "" {
		query += ` AND LOWER(COALESCE(car_class, '')) = LOWER(?)`
		args = append(args, carClass)
	}
	query += ` ORDER BY updated_at DESC, id DESC`
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var profiles []TuneProfile
	for rows.Next() {
		profile, err := scanProfile(rows)
		if err != nil {
			return nil, err
		}
		profiles = append(profiles, *profile)
	}
	return profiles, rows.Err()
}

func (s *Store) CreateTuneProfile(input TuneProfileInput) (*TuneProfile, error) {
	if strings.TrimSpace(input.CarName) == "" {
		return nil, errors.New("car name is required")
	}
	input = normalizeTuneProfilePower(input)
	now := nowText()
	result, err := s.db.Exec(`INSERT INTO tune_profile (
		car_name, car_ordinal, car_category, car_class, pi, drivetrain, num_cylinders, use_case, version_name,
		power_kw, torque_nm, weight_kg, front_weight_pct, power_to_weight_kw_per_kg, peak_torque_rpm, peak_power_rpm, redline_rpm, created_at, updated_at,
		front_tire_pressure, rear_tire_pressure, final_drive, gear_1, gear_2, gear_3, gear_4, gear_5, gear_6, gear_7, gear_8, gear_9, gear_10,
		front_camber, rear_camber, front_toe, rear_toe, caster, front_arb, rear_arb,
		front_spring, rear_spring, front_ride_height, rear_ride_height,
		front_rebound, rear_rebound, front_bump, rear_bump,
		front_aero, rear_aero, aero_balance, brake_balance, brake_pressure,
		front_diff_accel, front_diff_decel, rear_diff_accel, rear_diff_decel, center_diff_balance, notes
	) VALUES (`+placeholders(58)+`)`, profileInsertArgs(input, now, now)...)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	if err := s.upsertCarIdentity(input); err != nil {
		return nil, err
	}
	return s.GetTuneProfile(id)
}

func (s *Store) UpdateTuneProfile(id int64, input TuneProfileInput) (*TuneProfile, error) {
	return s.updateTuneProfile(id, input, "update")
}

func (s *Store) updateTuneProfile(id int64, input TuneProfileInput, reason string) (*TuneProfile, error) {
	if strings.TrimSpace(input.CarName) == "" {
		return nil, errors.New("car name is required")
	}
	input = normalizeTuneProfilePower(input)
	before, err := s.GetTuneProfile(id)
	if err != nil {
		return nil, err
	}
	now := nowText()
	args := append(profileUpdateArgs(input, now), id)
	result, err := s.db.Exec(`UPDATE tune_profile SET
		car_name = ?, car_ordinal = ?, car_category = ?, car_class = ?, pi = ?, drivetrain = ?, num_cylinders = ?, use_case = ?, version_name = ?,
		power_kw = ?, torque_nm = ?, weight_kg = ?, front_weight_pct = ?, power_to_weight_kw_per_kg = ?, peak_torque_rpm = ?, peak_power_rpm = ?, redline_rpm = ?, updated_at = ?,
		front_tire_pressure = ?, rear_tire_pressure = ?, final_drive = ?, gear_1 = ?, gear_2 = ?, gear_3 = ?, gear_4 = ?, gear_5 = ?, gear_6 = ?, gear_7 = ?, gear_8 = ?, gear_9 = ?, gear_10 = ?,
		front_camber = ?, rear_camber = ?, front_toe = ?, rear_toe = ?, caster = ?, front_arb = ?, rear_arb = ?,
		front_spring = ?, rear_spring = ?, front_ride_height = ?, rear_ride_height = ?,
		front_rebound = ?, rear_rebound = ?, front_bump = ?, rear_bump = ?,
		front_aero = ?, rear_aero = ?, aero_balance = ?, brake_balance = ?, brake_pressure = ?,
		front_diff_accel = ?, front_diff_decel = ?, rear_diff_accel = ?, rear_diff_decel = ?, center_diff_balance = ?, notes = ?
		WHERE id = ?`, args...)
	if err != nil {
		return nil, err
	}
	if changed, _ := result.RowsAffected(); changed == 0 {
		return nil, sql.ErrNoRows
	}
	if err := s.upsertCarIdentity(input); err != nil {
		return nil, err
	}
	after, err := s.GetTuneProfile(id)
	if err != nil {
		return nil, err
	}
	changedFields, err := changedTuneProfileFields(before, after)
	if err != nil {
		return nil, err
	}
	if len(changedFields) > 0 {
		if err := s.insertTuneProfileSnapshot(before, after, reason, nil, changedFields); err != nil {
			return nil, err
		}
		if err := s.pruneTuneProfileSnapshots(id, 5); err != nil {
			return nil, err
		}
	}
	return after, nil
}

func (s *Store) DuplicateTuneProfile(id int64, versionName string) (*TuneProfile, error) {
	profile, err := s.GetTuneProfile(id)
	if err != nil {
		return nil, err
	}
	input := profile.ToInput()
	if strings.TrimSpace(versionName) != "" {
		input.VersionName = strings.TrimSpace(versionName)
	} else if strings.TrimSpace(input.VersionName) != "" {
		input.VersionName += " Copy"
	} else {
		input.VersionName = "Copy"
	}
	return s.CreateTuneProfile(input)
}

func (s *Store) ListTuneProfileSnapshots(profileID int64) ([]TuneProfileSnapshot, error) {
	if profileID <= 0 {
		return []TuneProfileSnapshot{}, nil
	}
	rows, err := s.db.Query(`SELECT id, tune_profile_id, session_id, changed_at, COALESCE(change_reason, ''), change_json
		FROM tune_change_log WHERE tune_profile_id = ? ORDER BY changed_at DESC, id DESC LIMIT 5`, profileID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	snapshots := []TuneProfileSnapshot{}
	for rows.Next() {
		snapshot, err := scanTuneProfileSnapshot(rows)
		if err != nil {
			return nil, err
		}
		snapshots = append(snapshots, snapshot)
	}
	return snapshots, rows.Err()
}

func (s *Store) RestoreTuneProfileSnapshot(snapshotID int64) (*TuneProfile, error) {
	if snapshotID <= 0 {
		return nil, errors.New("snapshot id is required")
	}
	row := s.db.QueryRow(`SELECT id, tune_profile_id, session_id, changed_at, COALESCE(change_reason, ''), change_json
		FROM tune_change_log WHERE id = ?`, snapshotID)
	snapshot, err := scanTuneProfileSnapshot(row)
	if err != nil {
		return nil, err
	}
	if snapshot.Before.ID == 0 {
		return nil, errors.New("snapshot does not contain a restorable profile")
	}
	return s.updateTuneProfile(snapshot.TuneProfileID, snapshot.Before.ToInput(), "restore")
}

func (s *Store) insertTuneProfileSnapshot(before *TuneProfile, after *TuneProfile, reason string, sessionID *int64, changedFields []string) error {
	if before == nil || after == nil {
		return nil
	}
	payload := tuneProfileSnapshotPayload{
		Before:        *before,
		After:         *after,
		ChangedFields: changedFields,
	}
	raw, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	_, err = s.db.Exec(`INSERT INTO tune_change_log (tune_profile_id, session_id, changed_at, change_reason, change_json)
		VALUES (?, ?, ?, ?, ?)`, after.ID, nullableInt(sessionID), nowText(), reason, string(raw))
	return err
}

func (s *Store) pruneTuneProfileSnapshots(profileID int64, keep int) error {
	if keep <= 0 {
		keep = 5
	}
	_, err := s.db.Exec(`DELETE FROM tune_change_log
		WHERE tune_profile_id = ?
		AND id NOT IN (
			SELECT id FROM tune_change_log WHERE tune_profile_id = ? ORDER BY changed_at DESC, id DESC LIMIT ?
		)`, profileID, profileID, keep)
	return err
}

func (s *Store) DeleteTuneProfile(id int64) error {
	result, err := s.db.Exec(`DELETE FROM tune_profile WHERE id = ?`, id)
	if err != nil {
		return err
	}
	if changed, _ := result.RowsAffected(); changed == 0 {
		return sql.ErrNoRows
	}
	active, err := s.GetActiveTuneProfile()
	if err == nil && active != nil && active.ID == id {
		return s.clearActiveTuneProfile()
	}
	return nil
}

func (s *Store) SetActiveTuneProfile(id int64) error {
	if id <= 0 {
		return s.clearActiveTuneProfile()
	}
	if _, err := s.GetTuneProfile(id); err != nil {
		return err
	}
	_, err := s.db.Exec(`INSERT INTO app_setting(key, value) VALUES('active_tune_profile_id', ?)
		ON CONFLICT(key) DO UPDATE SET value = excluded.value`, fmt.Sprintf("%d", id))
	return err
}

func (s *Store) GetActiveTuneProfile() (*TuneProfile, error) {
	var value string
	err := s.db.QueryRow(`SELECT value FROM app_setting WHERE key = 'active_tune_profile_id'`).Scan(&value)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var id int64
	if _, err := fmt.Sscanf(value, "%d", &id); err != nil || id <= 0 {
		return nil, nil
	}
	profile, err := s.GetTuneProfile(id)
	if errors.Is(err, sql.ErrNoRows) {
		_ = s.clearActiveTuneProfile()
		return nil, nil
	}
	return profile, err
}

func (s *Store) clearActiveTuneProfile() error {
	_, err := s.db.Exec(`DELETE FROM app_setting WHERE key = 'active_tune_profile_id'`)
	return err
}

const testConditionDefaultsKey = "test_condition_defaults"

func DefaultTestConditions() TestConditions {
	return TestConditions{
		DriverMode:       "unknown",
		BrakeAssist:      "unknown",
		SteeringAssist:   "unknown",
		TractionControl:  "unknown",
		StabilityControl: "unknown",
		Shifting:         "unknown",
		LaunchControl:    "unknown",
	}
}

func NormalizeTestConditions(input TestConditions) TestConditions {
	return TestConditions{
		DriverMode:       NormalizeDriverMode(input.DriverMode),
		BrakeAssist:      normalizeAllowed(input.BrakeAssist, "unknown", "assisted", "abs_on", "abs_off"),
		SteeringAssist:   normalizeAllowed(input.SteeringAssist, "unknown", "auto", "assisted", "standard", "simulation"),
		TractionControl:  normalizeAllowed(input.TractionControl, "unknown", "on", "off"),
		StabilityControl: normalizeAllowed(input.StabilityControl, "unknown", "on", "off"),
		Shifting:         normalizeAllowed(input.Shifting, "unknown", "automatic", "manual"),
		LaunchControl:    normalizeAllowed(input.LaunchControl, "unknown", "on", "off"),
	}
}

func SessionTestConditions(session TelemetrySession) TestConditions {
	return NormalizeTestConditions(TestConditions{
		DriverMode:       session.DriverMode,
		BrakeAssist:      session.BrakeAssist,
		SteeringAssist:   session.SteeringAssist,
		TractionControl:  session.TractionControl,
		StabilityControl: session.StabilityControl,
		Shifting:         session.Shifting,
		LaunchControl:    session.LaunchControl,
	})
}

func TestConditionsContainUnknown(conditions TestConditions) bool {
	normalized := NormalizeTestConditions(conditions)
	return normalized.BrakeAssist == "unknown" ||
		normalized.SteeringAssist == "unknown" ||
		normalized.TractionControl == "unknown" ||
		normalized.StabilityControl == "unknown" ||
		normalized.Shifting == "unknown" ||
		normalized.LaunchControl == "unknown"
}

func TestConditionsEqual(left TestConditions, right TestConditions) bool {
	l := NormalizeTestConditions(left)
	r := NormalizeTestConditions(right)
	return l.DriverMode == r.DriverMode &&
		l.BrakeAssist == r.BrakeAssist &&
		l.SteeringAssist == r.SteeringAssist &&
		l.TractionControl == r.TractionControl &&
		l.StabilityControl == r.StabilityControl &&
		l.Shifting == r.Shifting &&
		l.LaunchControl == r.LaunchControl
}

func (s *Store) GetTestConditionDefaults() (TestConditions, error) {
	var value string
	err := s.db.QueryRow(`SELECT value FROM app_setting WHERE key = ?`, testConditionDefaultsKey).Scan(&value)
	if errors.Is(err, sql.ErrNoRows) {
		return DefaultTestConditions(), nil
	}
	if err != nil {
		return TestConditions{}, err
	}
	var conditions TestConditions
	if err := json.Unmarshal([]byte(value), &conditions); err != nil {
		return DefaultTestConditions(), nil
	}
	return NormalizeTestConditions(conditions), nil
}

func (s *Store) SaveTestConditionDefaults(conditions TestConditions) (TestConditions, error) {
	normalized := NormalizeTestConditions(conditions)
	payload, err := json.Marshal(normalized)
	if err != nil {
		return TestConditions{}, err
	}
	_, err = s.db.Exec(`INSERT INTO app_setting(key, value) VALUES(?, ?)
		ON CONFLICT(key) DO UPDATE SET value = excluded.value`, testConditionDefaultsKey, string(payload))
	return normalized, err
}

func normalizeAllowed(value string, allowed ...string) string {
	normalized := strings.ToLower(strings.TrimSpace(value))
	for _, candidate := range allowed {
		if normalized == candidate {
			return candidate
		}
	}
	return "unknown"
}

func changedTuneProfileFields(before *TuneProfile, after *TuneProfile) ([]string, error) {
	if before == nil || after == nil {
		return nil, nil
	}
	beforeMap, err := tuneProfileInputJSONMap(before.ToInput())
	if err != nil {
		return nil, err
	}
	afterMap, err := tuneProfileInputJSONMap(after.ToInput())
	if err != nil {
		return nil, err
	}
	keySet := map[string]struct{}{}
	for key := range beforeMap {
		keySet[key] = struct{}{}
	}
	for key := range afterMap {
		keySet[key] = struct{}{}
	}
	keys := make([]string, 0, len(keySet))
	for key := range keySet {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	changed := make([]string, 0)
	for _, key := range keys {
		if string(beforeMap[key]) != string(afterMap[key]) {
			changed = append(changed, key)
		}
	}
	return changed, nil
}

func tuneProfileInputJSONMap(input TuneProfileInput) (map[string]json.RawMessage, error) {
	raw, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}
	values := map[string]json.RawMessage{}
	if err := json.Unmarshal(raw, &values); err != nil {
		return nil, err
	}
	return values, nil
}

func (s *Store) CreateTelemetrySession(input SessionStartInput) (*TelemetrySession, error) {
	startedAt := input.StartedAt
	if strings.TrimSpace(startedAt) == "" {
		startedAt = nowText()
	}
	conditions := NormalizeTestConditions(input.TestConditions)
	result, err := s.db.Exec(`INSERT INTO telemetry_session (
		tune_profile_id, tune_snapshot_json, session_name, track_name, mode, game_mode, started_at, duration_ms, recording_path,
		car_ordinal, car_class, car_pi, drivetrain, num_cylinders,
		driver_mode, brake_assist, steering_assist, traction_control, stability_control, shifting, launch_control
	) VALUES (?, ?, ?, ?, ?, ?, ?, 0, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		nullableInt(input.TuneProfileID),
		strings.TrimSpace(input.TuneSnapshotJSON),
		input.SessionName,
		input.TrackName,
		input.Mode,
		telemetry.NormalizeGameMode(input.GameMode),
		startedAt,
		input.RecordingPath,
		nullableInt(input.CarOrdinal),
		strings.TrimSpace(input.CarClass),
		nullableInt(input.CarPI),
		strings.TrimSpace(input.Drivetrain),
		nullableInt(input.NumCylinders),
		driverModeUnknown,
		conditions.BrakeAssist,
		conditions.SteeringAssist,
		conditions.TractionControl,
		conditions.StabilityControl,
		conditions.Shifting,
		conditions.LaunchControl,
	)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	return s.GetTelemetrySession(id)
}

func (s *Store) DeleteTelemetrySession(id int64) error {
	_, err := s.db.Exec(`DELETE FROM telemetry_session WHERE id = ?`, id)
	return err
}

func (s *Store) FinalizeTelemetrySession(input SessionFinalizeInput, events []telemetry.DetectedEvent, samples []telemetry.NormalizedTelemetry) (*TelemetrySession, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	carClass := strings.TrimSpace(input.CarClass)
	drivetrain := strings.TrimSpace(input.Drivetrain)
	driverDetection := input.DriverModeDetection
	driverDetection.Mode = NormalizeDriverMode(driverDetection.Mode)
	if driverDetection.Mode == "" {
		driverDetection.Mode = driverModeUnknown
	}
	if driverDetection.Evidence == nil {
		driverDetection.Evidence = map[string]float64{}
	}
	driverEvidenceJSON, err := json.Marshal(driverDetection)
	if err != nil {
		return nil, err
	}
	_, err = tx.Exec(`UPDATE telemetry_session SET
		ended_at = ?, duration_ms = ?, avg_speed_kmh = ?, max_speed_kmh = ?,
		recording_packets = ?, recording_bytes = ?, recording_truncated = ?, game_mode = ?, notes = ?,
		car_ordinal = COALESCE(?, car_ordinal),
		car_class = CASE WHEN ? <> '' THEN ? ELSE car_class END,
		car_pi = COALESCE(?, car_pi),
		drivetrain = CASE WHEN ? <> '' THEN ? ELSE drivetrain END,
		num_cylinders = COALESCE(?, num_cylinders),
		driver_mode = ?, driver_mode_confidence = ?, driver_mode_evidence_json = ?
		WHERE id = ?`,
		input.EndedAt,
		input.DurationMS,
		nullableFloat(input.AvgSpeedKmh),
		nullableFloat(input.MaxSpeedKmh),
		input.RecordingPackets,
		input.RecordingBytes,
		boolInt(input.RecordingTruncated),
		telemetry.NormalizeGameMode(input.GameMode),
		input.Notes,
		nullableInt(input.CarOrdinal),
		carClass,
		carClass,
		nullableInt(input.CarPI),
		drivetrain,
		drivetrain,
		nullableInt(input.NumCylinders),
		driverDetection.Mode,
		driverDetection.Confidence,
		string(driverEvidenceJSON),
		input.SessionID,
	)
	if err != nil {
		return nil, err
	}
	if err := replaceEvents(tx, input.SessionID, events); err != nil {
		return nil, err
	}
	if err := replaceSamples(tx, input.SessionID, samples); err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return s.GetTelemetrySession(input.SessionID)
}

func (s *Store) ListTelemetrySessions(limit int) ([]TelemetrySession, error) {
	if limit <= 0 || limit > 500 {
		limit = 100
	}
	rows, err := s.db.Query(sessionSelectSQL+` GROUP BY ts.id ORDER BY ts.started_at DESC, ts.id DESC LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []TelemetrySession
	for rows.Next() {
		session, err := scanSession(rows)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, *session)
	}
	return sessions, rows.Err()
}

func (s *Store) GetTelemetrySession(id int64) (*TelemetrySession, error) {
	row := s.db.QueryRow(sessionSelectSQL+` WHERE ts.id = ? GROUP BY ts.id`, id)
	return scanSession(row)
}

func (s *Store) BindTelemetrySessionTuneProfile(sessionID int64, tuneProfileID int64) (*TelemetrySession, error) {
	if sessionID <= 0 {
		return nil, errors.New("session id is required")
	}
	if tuneProfileID <= 0 {
		return nil, errors.New("tune profile id is required")
	}

	session, err := s.GetTelemetrySession(sessionID)
	if err != nil {
		return nil, err
	}
	profile, err := s.GetTuneProfile(tuneProfileID)
	if err != nil {
		return nil, err
	}

	sessionClass := strings.TrimSpace(session.CarClass)
	if session.CarOrdinal == nil || *session.CarOrdinal <= 0 || sessionClass == "" {
		return nil, errors.New("cannot verify vehicle match: telemetry session has no vehicle snapshot")
	}

	profileClass := strings.TrimSpace(profile.CarClass)
	if profile.CarOrdinal == nil || *profile.CarOrdinal <= 0 || profileClass == "" {
		return nil, errors.New("tune profile vehicle identity is incomplete")
	}

	if *session.CarOrdinal != *profile.CarOrdinal || !strings.EqualFold(sessionClass, profileClass) {
		return nil, fmt.Errorf("tune profile does not match telemetry vehicle: session %d/%s, profile %d/%s", *session.CarOrdinal, sessionClass, *profile.CarOrdinal, profileClass)
	}

	result, err := s.db.Exec(`UPDATE telemetry_session SET tune_profile_id = ? WHERE id = ?`, tuneProfileID, sessionID)
	if err != nil {
		return nil, err
	}
	if changed, _ := result.RowsAffected(); changed == 0 {
		return nil, sql.ErrNoRows
	}
	return s.GetTelemetrySession(sessionID)
}

func (s *Store) GetSessionEvents(sessionID int64) ([]telemetry.DetectedEvent, error) {
	rows, err := s.db.Query(`SELECT id, event_type, severity, segment, start_ms, end_ms, duration_ms, evidence_json, suggestion_json
		FROM detected_event WHERE session_id = ? ORDER BY start_ms ASC, id ASC`, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := []telemetry.DetectedEvent{}
	for rows.Next() {
		event, err := scanEvent(rows)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, rows.Err()
}

func (s *Store) GetSessionTelemetrySamples(sessionID int64, limit int) ([]telemetry.NormalizedTelemetry, error) {
	if limit <= 0 || limit > 10000 {
		limit = 2000
	}
	rows, err := s.db.Query(`SELECT timestamp_ms, speed_kmh, rpm, rpm_ratio, gear, throttle, brake, steer,
		front_slip_ratio, rear_slip_ratio, front_combined_slip, rear_combined_slip,
		front_tire_temp, rear_tire_temp, front_suspension, rear_suspension,
		yaw_rate, pitch_rate, roll_rate, speed_field_kmh, velocity_speed_kmh, speed_source,
		game_mode, is_race_on, position_x, position_y, position_z, distance_traveled, best_lap, last_lap,
		current_lap, current_race_time, lap_number, race_position, driving_line,
		car_ordinal, car_category, car_class, car_pi, drivetrain, num_cylinders
		FROM telemetry_sample_agg WHERE session_id = ? ORDER BY timestamp_ms ASC, id ASC LIMIT ?`, sessionID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	samples := []telemetry.NormalizedTelemetry{}
	for rows.Next() {
		sample, err := scanSample(rows)
		if err != nil {
			return nil, err
		}
		samples = append(samples, sample)
	}
	return samples, rows.Err()
}

func (s *Store) ResolveCarNameByOrdinal(carOrdinal int64) (string, error) {
	if carOrdinal <= 0 {
		return "", nil
	}
	var name string
	err := s.db.QueryRow(`SELECT car_name FROM car_identity WHERE car_ordinal = ?`, carOrdinal).Scan(&name)
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}
	return name, err
}

func (s *Store) ListBenchmarkTracks() ([]BenchmarkTrack, error) {
	rows, err := s.db.Query(benchmarkTrackSelectSQL + ` ORDER BY updated_at DESC, id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	tracks := []BenchmarkTrack{}
	for rows.Next() {
		track, err := scanBenchmarkTrack(rows)
		if err != nil {
			return nil, err
		}
		tracks = append(tracks, *track)
	}
	return tracks, rows.Err()
}

func (s *Store) GetBenchmarkTrack(id int64) (*BenchmarkTrack, error) {
	row := s.db.QueryRow(benchmarkTrackSelectSQL+` WHERE id = ?`, id)
	return scanBenchmarkTrack(row)
}

func (s *Store) CreateBenchmarkTrack(input BenchmarkTrackInput) (*BenchmarkTrack, error) {
	normalized, err := normalizeBenchmarkTrackInput(input)
	if err != nil {
		return nil, err
	}
	polylineJSON, err := json.Marshal(normalized.Polyline)
	if err != nil {
		return nil, err
	}
	startGateJSON, finishGateJSON, checkpointsJSON, err := benchmarkTrackAuxJSON(normalized)
	if err != nil {
		return nil, err
	}
	now := nowText()
	result, err := s.db.Exec(`INSERT INTO benchmark_track (
		name, source_mode, track_type, start_x, start_y, start_z, end_x, end_y, end_z,
		start_radius, end_radius, direction_x, direction_z, route_length_meters,
		has_driving_line, start_gate_json, finish_gate_json, checkpoints_json, source_session_id,
		lap_count_observed, polyline_json, notes, created_at, updated_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		normalized.Name, normalized.SourceMode, normalized.TrackType, normalized.Start.X, normalized.Start.Y, normalized.Start.Z, normalized.End.X, normalized.End.Y, normalized.End.Z,
		normalized.StartRadius, normalized.EndRadius, normalized.DirectionX, normalized.DirectionZ, normalized.RouteLengthMeters,
		boolInt(normalized.HasDrivingLine), startGateJSON, finishGateJSON, checkpointsJSON, nullableInt(normalized.SourceSessionID),
		normalized.LapCountObserved, string(polylineJSON), normalized.Notes, now, now,
	)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	return s.GetBenchmarkTrack(id)
}

func (s *Store) UpdateBenchmarkTrack(id int64, input BenchmarkTrackInput) (*BenchmarkTrack, error) {
	normalized, err := normalizeBenchmarkTrackInput(input)
	if err != nil {
		return nil, err
	}
	polylineJSON, err := json.Marshal(normalized.Polyline)
	if err != nil {
		return nil, err
	}
	startGateJSON, finishGateJSON, checkpointsJSON, err := benchmarkTrackAuxJSON(normalized)
	if err != nil {
		return nil, err
	}
	result, err := s.db.Exec(`UPDATE benchmark_track SET
		name = ?, source_mode = ?, track_type = ?, start_x = ?, start_y = ?, start_z = ?, end_x = ?, end_y = ?, end_z = ?,
		start_radius = ?, end_radius = ?, direction_x = ?, direction_z = ?, route_length_meters = ?,
		has_driving_line = ?, start_gate_json = ?, finish_gate_json = ?, checkpoints_json = ?, source_session_id = ?,
		lap_count_observed = ?, polyline_json = ?, notes = ?, updated_at = ?
		WHERE id = ?`,
		normalized.Name, normalized.SourceMode, normalized.TrackType, normalized.Start.X, normalized.Start.Y, normalized.Start.Z, normalized.End.X, normalized.End.Y, normalized.End.Z,
		normalized.StartRadius, normalized.EndRadius, normalized.DirectionX, normalized.DirectionZ, normalized.RouteLengthMeters,
		boolInt(normalized.HasDrivingLine), startGateJSON, finishGateJSON, checkpointsJSON, nullableInt(normalized.SourceSessionID),
		normalized.LapCountObserved, string(polylineJSON), normalized.Notes, nowText(), id,
	)
	if err != nil {
		return nil, err
	}
	if changed, _ := result.RowsAffected(); changed == 0 {
		return nil, sql.ErrNoRows
	}
	return s.GetBenchmarkTrack(id)
}

func (s *Store) DeleteBenchmarkTrack(id int64) error {
	result, err := s.db.Exec(`DELETE FROM benchmark_track WHERE id = ?`, id)
	if err != nil {
		return err
	}
	if changed, _ := result.RowsAffected(); changed == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (s *Store) CreateBenchmarkTrackFromSession(sessionID int64, name string) (*BenchmarkTrack, error) {
	return s.ExtractBenchmarkTrackFromSession(BenchmarkTrackExtractionInput{
		SessionID:      sessionID,
		Name:           name,
		TrackType:      benchmarkTrackTypeAuto,
		ExtractionMode: benchmarkExtractionAutoBestLap,
	})
}

func (s *Store) ExtractBenchmarkTrackFromSession(input BenchmarkTrackExtractionInput) (*BenchmarkTrack, error) {
	trackInput, err := s.buildBenchmarkTrackInputFromExtraction(input)
	if err != nil {
		return nil, err
	}
	return s.CreateBenchmarkTrack(trackInput)
}

func (s *Store) ReextractBenchmarkTrack(trackID int64, input BenchmarkTrackExtractionInput) (*BenchmarkTrack, error) {
	current, err := s.GetBenchmarkTrack(trackID)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(input.Name) == "" {
		input.Name = current.Name
	}
	if strings.TrimSpace(input.TrackType) == "" {
		input.TrackType = current.TrackType
	}
	if strings.TrimSpace(input.ExtractionMode) == "" {
		input.ExtractionMode = benchmarkExtractionAutoBestLap
	}
	if input.StartGate == nil && !emptyGate(current.StartGate) {
		gate := current.StartGate
		input.StartGate = &gate
	}
	if input.FinishGate == nil && !emptyGate(current.FinishGate) {
		gate := current.FinishGate
		input.FinishGate = &gate
	}
	trackInput, err := s.buildBenchmarkTrackInputFromExtraction(input)
	if err != nil {
		return nil, err
	}
	return s.UpdateBenchmarkTrack(trackID, trackInput)
}

func (s *Store) buildBenchmarkTrackInputFromExtraction(input BenchmarkTrackExtractionInput) (BenchmarkTrackInput, error) {
	session, err := s.GetTelemetrySession(input.SessionID)
	if err != nil {
		return BenchmarkTrackInput{}, err
	}
	samples, err := s.GetSessionTelemetrySamples(input.SessionID, 10000)
	if err != nil {
		return BenchmarkTrackInput{}, err
	}
	if strings.TrimSpace(input.Name) == "" {
		input.Name = session.SessionName
	}
	trackInput, err := buildBenchmarkTrackInputFromSamples(input.Name, session.GameMode, samples, input)
	if err != nil {
		return BenchmarkTrackInput{}, err
	}
	trackInput.SourceSessionID = &input.SessionID
	return trackInput, nil
}

func (s *Store) AnalyzeSessionBenchmarkRuns(sessionID int64) ([]BenchmarkRun, error) {
	session, err := s.GetTelemetrySession(sessionID)
	if err != nil {
		return nil, err
	}
	samples, err := s.GetSessionTelemetrySamples(sessionID, 10000)
	if err != nil {
		return nil, err
	}
	events, err := s.GetSessionEvents(sessionID)
	if err != nil {
		return nil, err
	}
	tracks, err := s.ListBenchmarkTracks()
	if err != nil {
		return nil, err
	}
	if _, err := s.db.Exec(`DELETE FROM benchmark_run WHERE session_id = ?`, sessionID); err != nil {
		return nil, err
	}
	for _, track := range tracks {
		runs := analyzeBenchmarkRuns(sessionID, track, samples)
		if len(runs) == 0 {
			continue
		}
		for _, run := range runs {
			eventCount, err := s.countSessionEventsBetween(sessionID, run.StartMS, run.EndMS)
			if err != nil {
				return nil, err
			}
			run.EventCount = eventCount
			detection := DetectDriverMode(samplesBetween(samples, run.StartMS, run.EndMS), eventsBetween(events, run.StartMS, run.EndMS), session.GameMode)
			if detection.Mode == driverModeUnknown {
				detection = DriverModeDetection{
					Mode:       NormalizeDriverMode(session.DriverMode),
					Confidence: session.DriverModeConfidence,
					Summary:    "session_detection_fallback",
					Evidence:   map[string]float64{},
				}
			}
			run.DriverMode = NormalizeDriverMode(detection.Mode)
			run.DriverModeConfidence = detection.Confidence
			if raw, err := json.Marshal(detection); err == nil {
				run.DriverModeEvidenceJSON = string(raw)
			}
			if _, err := s.insertBenchmarkRun(run); err != nil {
				return nil, err
			}
		}
	}
	return s.listBenchmarkRunsForSession(sessionID)
}

func (s *Store) ListBenchmarkRuns(trackID int64, limit int) ([]BenchmarkRun, error) {
	if limit <= 0 || limit > 500 {
		limit = 100
	}
	rows, err := s.db.Query(benchmarkRunSelectSQL+` WHERE br.track_id = ? ORDER BY br.duration_ms ASC, br.created_at DESC, br.id DESC LIMIT ?`, trackID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanBenchmarkRuns(rows)
}

func (s *Store) listBenchmarkRunsForSession(sessionID int64) ([]BenchmarkRun, error) {
	rows, err := s.db.Query(benchmarkRunSelectSQL+` WHERE br.session_id = ? ORDER BY br.valid DESC, br.duration_ms ASC, br.confidence DESC, br.id DESC`, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanBenchmarkRuns(rows)
}

func (s *Store) insertBenchmarkRun(run BenchmarkRun) (*BenchmarkRun, error) {
	now := nowText()
	result, err := s.db.Exec(`INSERT INTO benchmark_run (
		session_id, track_id, start_ms, end_ms, duration_ms, confidence, avg_speed_kmh, max_speed_kmh,
		route_progress_01, geometry_length_meters, track_length_error_pct, distance_traveled_delta_meters,
		current_race_time_delta_seconds, avg_lateral_error_meters, max_lateral_error_meters, warning_flags,
		event_count, driver_mode, driver_mode_confidence, driver_mode_evidence_json, valid, created_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		run.SessionID, run.TrackID, run.StartMS, run.EndMS, run.DurationMS, run.Confidence,
		nullableFloat(run.AvgSpeedKmh), nullableFloat(run.MaxSpeedKmh),
		nullableFloat(run.RouteProgress01), nullableFloat(run.GeometryLengthMeters), nullableFloat(run.TrackLengthErrorPct), nullableFloat(run.DistanceTraveledDeltaMeters),
		nullableFloat(run.CurrentRaceTimeDeltaSeconds), nullableFloat(run.AvgLateralErrorMeters), nullableFloat(run.MaxLateralErrorMeters), run.WarningFlags,
		run.EventCount, run.DriverMode, run.DriverModeConfidence, run.DriverModeEvidenceJSON, boolInt(run.Valid), now,
	)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	row := s.db.QueryRow(benchmarkRunSelectSQL+` WHERE br.id = ?`, id)
	saved, err := scanBenchmarkRun(row)
	if err != nil {
		return nil, err
	}
	return &saved, nil
}

func (s *Store) countSessionEventsBetween(sessionID int64, startMS int64, endMS int64) (int64, error) {
	var count int64
	err := s.db.QueryRow(`SELECT COUNT(*) FROM detected_event WHERE session_id = ? AND COALESCE(start_ms, 0) <= ? AND COALESCE(end_ms, start_ms, 0) >= ?`, sessionID, endMS, startMS).Scan(&count)
	return count, err
}

func samplesBetween(samples []telemetry.NormalizedTelemetry, startMS int64, endMS int64) []telemetry.NormalizedTelemetry {
	if endMS <= startMS {
		return nil
	}
	out := make([]telemetry.NormalizedTelemetry, 0, len(samples))
	for _, sample := range samples {
		if sample.TimeMS >= startMS && sample.TimeMS <= endMS {
			out = append(out, sample)
		}
	}
	return out
}

func eventsBetween(events []telemetry.DetectedEvent, startMS int64, endMS int64) []telemetry.DetectedEvent {
	if endMS <= startMS {
		return nil
	}
	out := make([]telemetry.DetectedEvent, 0, len(events))
	for _, event := range events {
		if event.StartMS <= endMS && event.EndMS >= startMS {
			out = append(out, event)
		}
	}
	return out
}

func (s *Store) ListRuleThresholdProfiles() ([]RuleThresholdProfile, error) {
	rows, err := s.db.Query(`SELECT id, name, COALESCE(car_class, ''), COALESCE(drivetrain, ''), COALESCE(use_case, ''), COALESCE(game_mode, ''), config_json, COALESCE(is_default, 0), created_at, updated_at
		FROM rule_threshold_profile ORDER BY is_default DESC, updated_at DESC, id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var profiles []RuleThresholdProfile
	for rows.Next() {
		profile, err := scanRuleThresholdProfile(rows)
		if err != nil {
			return nil, err
		}
		profiles = append(profiles, profile)
	}
	return profiles, rows.Err()
}

func (s *Store) CreateRuleThresholdProfile(input RuleThresholdProfileInput) (*RuleThresholdProfile, error) {
	configJSON, err := normalizeRuleConfigJSON(input.ConfigJSON)
	if err != nil {
		return nil, err
	}
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return nil, errors.New("rule threshold profile name is required")
	}
	now := nowText()
	result, err := s.db.Exec(`INSERT INTO rule_threshold_profile(name, car_class, drivetrain, use_case, game_mode, config_json, is_default, created_at, updated_at)
		VALUES(?, ?, ?, ?, ?, ?, 0, ?, ?)`, name, strings.TrimSpace(input.CarClass), strings.TrimSpace(input.Drivetrain), strings.TrimSpace(input.UseCase), normalizeOptionalGameMode(input.GameMode), configJSON, now, now)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	return s.GetRuleThresholdProfile(id)
}

func (s *Store) UpdateRuleThresholdProfile(id int64, input RuleThresholdProfileInput) (*RuleThresholdProfile, error) {
	configJSON, err := normalizeRuleConfigJSON(input.ConfigJSON)
	if err != nil {
		return nil, err
	}
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return nil, errors.New("rule threshold profile name is required")
	}
	result, err := s.db.Exec(`UPDATE rule_threshold_profile SET name = ?, car_class = ?, drivetrain = ?, use_case = ?, game_mode = ?, config_json = ?, updated_at = ? WHERE id = ?`,
		name, strings.TrimSpace(input.CarClass), strings.TrimSpace(input.Drivetrain), strings.TrimSpace(input.UseCase), normalizeOptionalGameMode(input.GameMode), configJSON, nowText(), id)
	if err != nil {
		return nil, err
	}
	if changed, _ := result.RowsAffected(); changed == 0 {
		return nil, sql.ErrNoRows
	}
	return s.GetRuleThresholdProfile(id)
}

func (s *Store) DeleteRuleThresholdProfile(id int64) error {
	profile, err := s.GetRuleThresholdProfile(id)
	if err != nil {
		return err
	}
	if profile.IsDefault {
		return errors.New("default rule threshold profile cannot be deleted")
	}
	result, err := s.db.Exec(`DELETE FROM rule_threshold_profile WHERE id = ?`, id)
	if err != nil {
		return err
	}
	if changed, _ := result.RowsAffected(); changed == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (s *Store) ResetRuleThresholdProfile(id int64) (*RuleThresholdProfile, error) {
	profile, err := s.GetRuleThresholdProfile(id)
	if err != nil {
		return nil, err
	}
	configJSON, err := ruleConfigJSONForProfile(profile)
	if err != nil {
		return nil, err
	}
	result, err := s.db.Exec(`UPDATE rule_threshold_profile SET config_json = ?, updated_at = ? WHERE id = ?`, configJSON, nowText(), id)
	if err != nil {
		return nil, err
	}
	if changed, _ := result.RowsAffected(); changed == 0 {
		return nil, sql.ErrNoRows
	}
	return s.GetRuleThresholdProfile(id)
}

func (s *Store) GetRuleThresholdProfile(id int64) (*RuleThresholdProfile, error) {
	row := s.db.QueryRow(`SELECT id, name, COALESCE(car_class, ''), COALESCE(drivetrain, ''), COALESCE(use_case, ''), COALESCE(game_mode, ''), config_json, COALESCE(is_default, 0), created_at, updated_at
		FROM rule_threshold_profile WHERE id = ?`, id)
	profile, err := scanRuleThresholdProfile(row)
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (s *Store) MatchRuleThresholdProfile(profile *TuneProfile) (*RuleThresholdProfile, telemetry.RuleConfig, error) {
	profiles, err := s.ListRuleThresholdProfiles()
	if err != nil {
		return nil, telemetry.DefaultRuleConfig(), err
	}
	bestScore := -1
	var best *RuleThresholdProfile
	for i := range profiles {
		candidate := &profiles[i]
		score, ok := ruleProfileMatchScore(candidate, profile)
		if !ok {
			continue
		}
		if score > bestScore {
			bestScore = score
			best = candidate
		}
	}
	if best == nil {
		return nil, telemetry.DefaultRuleConfig(), nil
	}
	config, err := parseRuleConfigJSON(best.ConfigJSON)
	if err != nil {
		return nil, telemetry.DefaultRuleConfig(), err
	}
	return best, config, nil
}

func (s *Store) ListTuneProfileSessionStats() ([]TuneProfileSessionStat, error) {
	rows, err := s.db.Query(`SELECT tune_profile_id, COUNT(*), COALESCE(MAX(started_at), '') FROM telemetry_session WHERE tune_profile_id IS NOT NULL GROUP BY tune_profile_id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var stats []TuneProfileSessionStat
	for rows.Next() {
		var stat TuneProfileSessionStat
		if err := rows.Scan(&stat.TuneProfileID, &stat.SessionCount, &stat.LastStartedAt); err != nil {
			return nil, err
		}
		stats = append(stats, stat)
	}
	return stats, rows.Err()
}

func (p TuneProfile) ToInput() TuneProfileInput {
	return TuneProfileInput{
		CarName: p.CarName, CarOrdinal: p.CarOrdinal, CarCategory: p.CarCategory, CarClass: p.CarClass, PI: p.PI, Drivetrain: p.Drivetrain, NumCylinders: p.NumCylinders, UseCase: p.UseCase, VersionName: p.VersionName,
		PowerKW: p.PowerKW, TorqueNM: p.TorqueNM, WeightKG: p.WeightKG, FrontWeightPct: p.FrontWeightPct, PowerToWeightKWPerKG: p.PowerToWeightKWPerKG,
		PeakTorqueRPM: p.PeakTorqueRPM, PeakPowerRPM: p.PeakPowerRPM, RedlineRPM: p.RedlineRPM,
		FrontTirePressure: p.FrontTirePressure, RearTirePressure: p.RearTirePressure, FinalDrive: p.FinalDrive,
		Gear1: p.Gear1, Gear2: p.Gear2, Gear3: p.Gear3, Gear4: p.Gear4, Gear5: p.Gear5, Gear6: p.Gear6, Gear7: p.Gear7, Gear8: p.Gear8, Gear9: p.Gear9, Gear10: p.Gear10,
		FrontCamber: p.FrontCamber, RearCamber: p.RearCamber, FrontToe: p.FrontToe, RearToe: p.RearToe, Caster: p.Caster,
		FrontARB: p.FrontARB, RearARB: p.RearARB,
		FrontSpring: p.FrontSpring, RearSpring: p.RearSpring, FrontRideHeight: p.FrontRideHeight, RearRideHeight: p.RearRideHeight,
		FrontRebound: p.FrontRebound, RearRebound: p.RearRebound, FrontBump: p.FrontBump, RearBump: p.RearBump,
		FrontAero: p.FrontAero, RearAero: p.RearAero, AeroBalance: p.AeroBalance,
		BrakeBalance: p.BrakeBalance, BrakePressure: p.BrakePressure,
		FrontDiffAccel: p.FrontDiffAccel, FrontDiffDecel: p.FrontDiffDecel, RearDiffAccel: p.RearDiffAccel, RearDiffDecel: p.RearDiffDecel, CenterDiffBalance: p.CenterDiffBalance,
		Notes: p.Notes,
	}
}

const profileSelectColumns = `id, car_name, car_ordinal, car_category, car_class, pi, drivetrain, num_cylinders, use_case, version_name, created_at, updated_at,
	power_kw, torque_nm, weight_kg, front_weight_pct, power_to_weight_kw_per_kg, peak_torque_rpm, peak_power_rpm, redline_rpm,
	front_tire_pressure, rear_tire_pressure, final_drive, gear_1, gear_2, gear_3, gear_4, gear_5, gear_6, gear_7, gear_8, gear_9, gear_10,
	front_camber, rear_camber, front_toe, rear_toe, caster, front_arb, rear_arb,
	front_spring, rear_spring, front_ride_height, rear_ride_height,
	front_rebound, rear_rebound, front_bump, rear_bump,
	front_aero, rear_aero, aero_balance, brake_balance, brake_pressure,
	front_diff_accel, front_diff_decel, rear_diff_accel, rear_diff_decel, center_diff_balance, notes`

const sessionSelectSQL = `SELECT ts.id, ts.tune_profile_id, COALESCE(ts.tune_snapshot_json, ''), COALESCE(tp.car_name, ''), COALESCE(ts.session_name, ''), COALESCE(ts.track_name, ''), COALESCE(ts.mode, ''), COALESCE(ts.game_mode, 'unknown'),
	COALESCE(ts.started_at, ''), COALESCE(ts.ended_at, ''), COALESCE(ts.duration_ms, 0), ts.best_lap_ms, ts.avg_speed_kmh, ts.max_speed_kmh,
	COALESCE(ts.recording_path, ''), COALESCE(ts.recording_packets, 0), COALESCE(ts.recording_bytes, 0), COALESCE(ts.recording_truncated, 0),
	ts.car_ordinal, COALESCE(ts.car_class, ''), ts.car_pi, COALESCE(ts.drivetrain, ''), ts.num_cylinders,
	COALESCE(ts.driver_mode, 'unknown'), COALESCE(ts.driver_mode_confidence, 0), COALESCE(ts.driver_mode_evidence_json, ''), COALESCE(ts.brake_assist, 'unknown'), COALESCE(ts.steering_assist, 'unknown'),
	COALESCE(ts.traction_control, 'unknown'), COALESCE(ts.stability_control, 'unknown'), COALESCE(ts.shifting, 'unknown'), COALESCE(ts.launch_control, 'unknown'),
	COALESCE(ts.driver_feedback_json, ''), COALESCE(ts.notes, ''), COUNT(DISTINCT de.id), COUNT(DISTINCT tsa.id)
	FROM telemetry_session ts
	LEFT JOIN tune_profile tp ON tp.id = ts.tune_profile_id
	LEFT JOIN detected_event de ON de.session_id = ts.id
	LEFT JOIN telemetry_sample_agg tsa ON tsa.session_id = ts.id`

const benchmarkTrackSelectSQL = `SELECT id, name, COALESCE(source_mode, ''), COALESCE(track_type, ''), start_x, start_y, start_z, end_x, end_y, end_z,
	COALESCE(start_radius, 20), COALESCE(end_radius, 20), COALESCE(direction_x, 0), COALESCE(direction_z, 0),
	COALESCE(route_length_meters, 0), COALESCE(has_driving_line, 0), COALESCE(start_gate_json, ''), COALESCE(finish_gate_json, ''),
	COALESCE(checkpoints_json, ''), source_session_id, COALESCE(lap_count_observed, 0), COALESCE(polyline_json, '[]'),
	COALESCE(notes, ''), created_at, updated_at
	FROM benchmark_track`

const benchmarkRunSelectSQL = `SELECT br.id, br.session_id, br.track_id, COALESCE(bt.name, ''), COALESCE(br.start_ms, 0), COALESCE(br.end_ms, 0),
	COALESCE(br.duration_ms, 0), COALESCE(br.confidence, 0), br.avg_speed_kmh, br.max_speed_kmh,
	br.route_progress_01, br.geometry_length_meters, br.track_length_error_pct, br.distance_traveled_delta_meters,
	br.current_race_time_delta_seconds, br.avg_lateral_error_meters, br.max_lateral_error_meters, COALESCE(br.warning_flags, ''),
	COALESCE(br.event_count, 0), COALESCE(br.driver_mode, 'unknown'), COALESCE(br.driver_mode_confidence, 0), COALESCE(br.driver_mode_evidence_json, ''), COALESCE(br.valid, 0), COALESCE(br.created_at, '')
	FROM benchmark_run br
	LEFT JOIN benchmark_track bt ON bt.id = br.track_id`

type scanner interface {
	Scan(dest ...any) error
}

type tuneProfileSnapshotPayload struct {
	Before        TuneProfile `json:"before"`
	After         TuneProfile `json:"after"`
	ChangedFields []string    `json:"changedFields"`
}

func TuneProfileSnapshotJSON(profile *TuneProfile) (string, error) {
	if profile == nil {
		return "", nil
	}
	raw, err := json.Marshal(profile)
	if err != nil {
		return "", err
	}
	return string(raw), nil
}

func ParseTuneProfileSnapshotJSON(value string) (*TuneProfile, error) {
	if strings.TrimSpace(value) == "" {
		return nil, nil
	}
	var profile TuneProfile
	if err := json.Unmarshal([]byte(value), &profile); err != nil {
		return nil, err
	}
	if profile.ID == 0 && strings.TrimSpace(profile.CarName) == "" {
		return nil, nil
	}
	return &profile, nil
}

func scanRuleThresholdProfile(row scanner) (RuleThresholdProfile, error) {
	var profile RuleThresholdProfile
	var isDefault int
	err := row.Scan(&profile.ID, &profile.Name, &profile.CarClass, &profile.Drivetrain, &profile.UseCase, &profile.GameMode, &profile.ConfigJSON, &isDefault, &profile.CreatedAt, &profile.UpdatedAt)
	profile.IsDefault = isDefault != 0
	return profile, err
}

func scanTuneProfileSnapshot(row scanner) (TuneProfileSnapshot, error) {
	var snapshot TuneProfileSnapshot
	var sessionID sql.NullInt64
	var raw string
	if err := row.Scan(&snapshot.ID, &snapshot.TuneProfileID, &sessionID, &snapshot.ChangedAt, &snapshot.ChangeReason, &raw); err != nil {
		return snapshot, err
	}
	var payload tuneProfileSnapshotPayload
	if strings.TrimSpace(raw) != "" {
		if err := json.Unmarshal([]byte(raw), &payload); err != nil {
			return snapshot, err
		}
	}
	snapshot.SessionID = intPtr(sessionID)
	snapshot.Before = payload.Before
	snapshot.After = payload.After
	snapshot.ChangedFields = payload.ChangedFields
	snapshot.ChangeJSON = raw
	return snapshot, nil
}

func scanProfile(row scanner) (*TuneProfile, error) {
	var p TuneProfile
	var carOrdinal, carCategory, pi, numCylinders sql.NullInt64
	var carClass, drivetrain, useCase, versionName, notes sql.NullString
	var floats [46]sql.NullFloat64
	err := row.Scan(
		&p.ID, &p.CarName, &carOrdinal, &carCategory, &carClass, &pi, &drivetrain, &numCylinders, &useCase, &versionName, &p.CreatedAt, &p.UpdatedAt,
		&floats[0], &floats[1], &floats[2], &floats[3], &floats[4], &floats[5], &floats[6], &floats[7],
		&floats[8], &floats[9], &floats[10], &floats[11], &floats[12], &floats[13], &floats[14], &floats[15], &floats[16], &floats[17], &floats[18], &floats[19], &floats[20],
		&floats[21], &floats[22], &floats[23], &floats[24], &floats[25], &floats[26], &floats[27],
		&floats[28], &floats[29], &floats[30], &floats[31],
		&floats[32], &floats[33], &floats[34], &floats[35],
		&floats[36], &floats[37], &floats[38], &floats[39], &floats[40],
		&floats[41], &floats[42], &floats[43], &floats[44], &floats[45], &notes,
	)
	if err != nil {
		return nil, err
	}
	p.CarOrdinal = intPtr(carOrdinal)
	p.CarCategory = intPtr(carCategory)
	p.CarClass = carClass.String
	p.PI = intPtr(pi)
	p.Drivetrain = drivetrain.String
	p.NumCylinders = intPtr(numCylinders)
	p.UseCase = useCase.String
	p.VersionName = versionName.String
	p.PowerKW = floatPtr(floats[0])
	p.TorqueNM = floatPtr(floats[1])
	p.WeightKG = floatPtr(floats[2])
	p.FrontWeightPct = floatPtr(floats[3])
	p.PowerToWeightKWPerKG = floatPtr(floats[4])
	p.PeakTorqueRPM = floatPtr(floats[5])
	p.PeakPowerRPM = floatPtr(floats[6])
	p.RedlineRPM = floatPtr(floats[7])
	p.FrontTirePressure = floatPtr(floats[8])
	p.RearTirePressure = floatPtr(floats[9])
	p.FinalDrive = floatPtr(floats[10])
	p.Gear1 = floatPtr(floats[11])
	p.Gear2 = floatPtr(floats[12])
	p.Gear3 = floatPtr(floats[13])
	p.Gear4 = floatPtr(floats[14])
	p.Gear5 = floatPtr(floats[15])
	p.Gear6 = floatPtr(floats[16])
	p.Gear7 = floatPtr(floats[17])
	p.Gear8 = floatPtr(floats[18])
	p.Gear9 = floatPtr(floats[19])
	p.Gear10 = floatPtr(floats[20])
	p.FrontCamber = floatPtr(floats[21])
	p.RearCamber = floatPtr(floats[22])
	p.FrontToe = floatPtr(floats[23])
	p.RearToe = floatPtr(floats[24])
	p.Caster = floatPtr(floats[25])
	p.FrontARB = floatPtr(floats[26])
	p.RearARB = floatPtr(floats[27])
	p.FrontSpring = floatPtr(floats[28])
	p.RearSpring = floatPtr(floats[29])
	p.FrontRideHeight = floatPtr(floats[30])
	p.RearRideHeight = floatPtr(floats[31])
	p.FrontRebound = floatPtr(floats[32])
	p.RearRebound = floatPtr(floats[33])
	p.FrontBump = floatPtr(floats[34])
	p.RearBump = floatPtr(floats[35])
	p.FrontAero = floatPtr(floats[36])
	p.RearAero = floatPtr(floats[37])
	p.AeroBalance = floatPtr(floats[38])
	p.BrakeBalance = floatPtr(floats[39])
	p.BrakePressure = floatPtr(floats[40])
	p.FrontDiffAccel = floatPtr(floats[41])
	p.FrontDiffDecel = floatPtr(floats[42])
	p.RearDiffAccel = floatPtr(floats[43])
	p.RearDiffDecel = floatPtr(floats[44])
	p.CenterDiffBalance = floatPtr(floats[45])
	p.Notes = notes.String
	return &p, nil
}

func scanSession(row scanner) (*TelemetrySession, error) {
	var s TelemetrySession
	var tuneProfileID, bestLap, carOrdinal, carPI, numCylinders sql.NullInt64
	var carClass, drivetrain, driverMode, brakeAssist, steeringAssist, tractionControl, stabilityControl, shifting, launchControl sql.NullString
	var avgSpeed, maxSpeed, driverModeConfidence sql.NullFloat64
	var recordingTruncated int
	err := row.Scan(
		&s.ID, &tuneProfileID, &s.TuneSnapshotJSON, &s.TuneName, &s.SessionName, &s.TrackName, &s.Mode, &s.GameMode, &s.StartedAt, &s.EndedAt, &s.DurationMS, &bestLap, &avgSpeed, &maxSpeed,
		&s.RecordingPath, &s.RecordingPackets, &s.RecordingBytes, &recordingTruncated,
		&carOrdinal, &carClass, &carPI, &drivetrain, &numCylinders,
		&driverMode, &driverModeConfidence, &s.DriverModeEvidenceJSON, &brakeAssist, &steeringAssist, &tractionControl, &stabilityControl, &shifting, &launchControl,
		&s.DriverFeedbackJSON, &s.Notes, &s.EventCount, &s.SampleCount,
	)
	if err != nil {
		return nil, err
	}
	s.TuneProfileID = intPtr(tuneProfileID)
	s.BestLapMS = intPtr(bestLap)
	s.AvgSpeedKmh = floatPtr(avgSpeed)
	s.MaxSpeedKmh = floatPtr(maxSpeed)
	s.RecordingTruncated = recordingTruncated != 0
	s.CarOrdinal = intPtr(carOrdinal)
	s.CarClass = carClass.String
	s.CarPI = intPtr(carPI)
	s.Drivetrain = drivetrain.String
	s.NumCylinders = intPtr(numCylinders)
	s.DriverModeConfidence = floatFromNull(driverModeConfidence)
	s.GameMode = telemetry.NormalizeGameMode(s.GameMode)
	conditions := NormalizeTestConditions(TestConditions{
		DriverMode:       driverMode.String,
		BrakeAssist:      brakeAssist.String,
		SteeringAssist:   steeringAssist.String,
		TractionControl:  tractionControl.String,
		StabilityControl: stabilityControl.String,
		Shifting:         shifting.String,
		LaunchControl:    launchControl.String,
	})
	s.DriverMode = conditions.DriverMode
	s.BrakeAssist = conditions.BrakeAssist
	s.SteeringAssist = conditions.SteeringAssist
	s.TractionControl = conditions.TractionControl
	s.StabilityControl = conditions.StabilityControl
	s.Shifting = conditions.Shifting
	s.LaunchControl = conditions.LaunchControl
	return &s, nil
}

func scanBenchmarkTrack(row scanner) (*BenchmarkTrack, error) {
	var track BenchmarkTrack
	var hasDrivingLine int
	var startGateJSON, finishGateJSON, checkpointsJSON, polylineJSON string
	var sourceSessionID sql.NullInt64
	err := row.Scan(
		&track.ID, &track.Name, &track.SourceMode, &track.TrackType,
		&track.Start.X, &track.Start.Y, &track.Start.Z,
		&track.End.X, &track.End.Y, &track.End.Z,
		&track.StartRadius, &track.EndRadius, &track.DirectionX, &track.DirectionZ,
		&track.RouteLengthMeters, &hasDrivingLine, &startGateJSON, &finishGateJSON,
		&checkpointsJSON, &sourceSessionID, &track.LapCountObserved, &polylineJSON,
		&track.Notes, &track.CreatedAt, &track.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	track.SourceMode = telemetry.NormalizeGameMode(track.SourceMode)
	track.HasDrivingLine = hasDrivingLine != 0
	track.SourceSessionID = intPtr(sourceSessionID)
	if err := json.Unmarshal([]byte(polylineJSON), &track.Polyline); err != nil {
		return nil, err
	}
	if strings.TrimSpace(startGateJSON) != "" {
		if err := json.Unmarshal([]byte(startGateJSON), &track.StartGate); err != nil {
			return nil, err
		}
	}
	if strings.TrimSpace(finishGateJSON) != "" {
		if err := json.Unmarshal([]byte(finishGateJSON), &track.FinishGate); err != nil {
			return nil, err
		}
	}
	if strings.TrimSpace(checkpointsJSON) != "" {
		if err := json.Unmarshal([]byte(checkpointsJSON), &track.Checkpoints); err != nil {
			return nil, err
		}
	}
	normalized, err := normalizeBenchmarkTrackInput(track.BenchmarkTrackInput)
	if err != nil {
		return nil, err
	}
	track.BenchmarkTrackInput = normalized
	return &track, nil
}

func scanBenchmarkRun(row scanner) (BenchmarkRun, error) {
	var run BenchmarkRun
	var avgSpeed, maxSpeed sql.NullFloat64
	var routeProgress, geometryLength, trackLengthError, distanceDelta, raceTimeDelta, avgLateral, maxLateral, driverModeConfidence sql.NullFloat64
	var valid int
	err := row.Scan(
		&run.ID, &run.SessionID, &run.TrackID, &run.TrackName, &run.StartMS, &run.EndMS,
		&run.DurationMS, &run.Confidence, &avgSpeed, &maxSpeed,
		&routeProgress, &geometryLength, &trackLengthError, &distanceDelta,
		&raceTimeDelta, &avgLateral, &maxLateral, &run.WarningFlags,
		&run.EventCount, &run.DriverMode, &driverModeConfidence, &run.DriverModeEvidenceJSON, &valid, &run.CreatedAt,
	)
	run.AvgSpeedKmh = floatPtr(avgSpeed)
	run.MaxSpeedKmh = floatPtr(maxSpeed)
	run.RouteProgress01 = floatPtr(routeProgress)
	run.GeometryLengthMeters = floatPtr(geometryLength)
	run.TrackLengthErrorPct = floatPtr(trackLengthError)
	run.DistanceTraveledDeltaMeters = floatPtr(distanceDelta)
	run.CurrentRaceTimeDeltaSeconds = floatPtr(raceTimeDelta)
	run.AvgLateralErrorMeters = floatPtr(avgLateral)
	run.MaxLateralErrorMeters = floatPtr(maxLateral)
	run.Valid = valid != 0
	run.DriverModeConfidence = floatFromNull(driverModeConfidence)
	if strings.TrimSpace(run.DriverMode) == "" {
		run.DriverMode = "unknown"
	}
	return run, err
}

func scanBenchmarkRuns(rows *sql.Rows) ([]BenchmarkRun, error) {
	runs := []BenchmarkRun{}
	for rows.Next() {
		run, err := scanBenchmarkRun(rows)
		if err != nil {
			return nil, err
		}
		runs = append(runs, run)
	}
	return runs, rows.Err()
}

func scanEvent(row scanner) (telemetry.DetectedEvent, error) {
	var id int64
	var event telemetry.DetectedEvent
	var evidenceJSON, suggestionJSON string
	if err := row.Scan(&id, &event.Type, &event.Severity, &event.Segment, &event.StartMS, &event.EndMS, &event.DurationMS, &evidenceJSON, &suggestionJSON); err != nil {
		return event, err
	}
	event.ID = fmt.Sprintf("db-%d", id)
	if err := json.Unmarshal([]byte(evidenceJSON), &event.Evidence); err != nil {
		return event, err
	}
	if strings.TrimSpace(suggestionJSON) != "" {
		if err := json.Unmarshal([]byte(suggestionJSON), &event.SuggestedActions); err != nil {
			return event, err
		}
		event.SuggestedActions = telemetry.NormalizeSuggestedActions(event.Type, event.SuggestedActions)
	}
	return event, nil
}

func replaceEvents(tx *sql.Tx, sessionID int64, events []telemetry.DetectedEvent) error {
	if _, err := tx.Exec(`DELETE FROM detected_event WHERE session_id = ?`, sessionID); err != nil {
		return err
	}
	for _, event := range events {
		evidenceJSON, err := json.Marshal(event.Evidence)
		if err != nil {
			return err
		}
		suggestionJSON, err := json.Marshal(event.SuggestedActions)
		if err != nil {
			return err
		}
		if _, err := tx.Exec(`INSERT INTO detected_event (
			session_id, event_type, severity, segment, start_ms, end_ms, duration_ms, evidence_json, suggestion_json, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			sessionID, event.Type, event.Severity, event.Segment, event.StartMS, event.EndMS, event.DurationMS, string(evidenceJSON), string(suggestionJSON), nowText(),
		); err != nil {
			return err
		}
	}
	return nil
}

func replaceSamples(tx *sql.Tx, sessionID int64, samples []telemetry.NormalizedTelemetry) error {
	if _, err := tx.Exec(`DELETE FROM telemetry_sample_agg WHERE session_id = ?`, sessionID); err != nil {
		return err
	}
	stmt, err := tx.Prepare(`INSERT INTO telemetry_sample_agg (
		session_id, timestamp_ms, speed_kmh, rpm, rpm_ratio, gear, throttle, brake, steer,
		front_slip_ratio, rear_slip_ratio, front_combined_slip, rear_combined_slip,
		front_tire_temp, rear_tire_temp, front_suspension, rear_suspension,
		yaw_rate, pitch_rate, roll_rate, speed_field_kmh, velocity_speed_kmh, speed_source,
		game_mode, is_race_on, position_x, position_y, position_z, distance_traveled, best_lap, last_lap,
		current_lap, current_race_time, lap_number, race_position, driving_line,
		car_ordinal, car_category, car_class, car_pi, drivetrain, num_cylinders
	) VALUES (` + placeholders(42) + `)`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, sample := range samples {
		if _, err := stmt.Exec(
			sessionID, sample.TimeMS, sample.SpeedKmh, sample.Rpm, sample.RpmRatio, sample.Gear, sample.Throttle01, sample.Brake01, sample.Steer01,
			sample.FrontSlipRatioAvg, sample.RearSlipRatioAvg, sample.FrontCombinedSlipAvg, sample.RearCombinedSlipAvg,
			sample.TireTempFrontAvg, sample.TireTempRearAvg, sample.SuspensionFrontAvg, sample.SuspensionRearAvg,
			sample.YawRate, sample.PitchRate, sample.RollRate, sample.SpeedFieldKmh, sample.VelocitySpeedKmh, sample.SpeedSource,
			telemetry.NormalizeGameMode(sample.GameMode), boolInt(sample.IsRaceOn),
			sample.PositionX, sample.PositionY, sample.PositionZ, sample.DistanceTraveled, sample.BestLap, sample.LastLap,
			sample.CurrentLap, sample.CurrentRaceTime, zeroAsNil(sample.LapNumber), zeroAsNil(sample.RacePosition), sample.DrivingLine01,
			zeroAsNil(sample.CarOrdinal), zeroAsNil(sample.CarCategory), sample.CarClass, zeroAsNil(sample.CarPI), sample.Drivetrain, zeroAsNil(sample.NumCylinders),
		); err != nil {
			return err
		}
	}
	return nil
}

func scanSample(row scanner) (telemetry.NormalizedTelemetry, error) {
	var sample telemetry.NormalizedTelemetry
	var speedSource, gameMode, carClass, drivetrain sql.NullString
	var isRaceOn, lapNumber, racePosition, carOrdinal, carCategory, carPI, numCylinders sql.NullInt64
	var speedField, velocitySpeed, positionX, positionY, positionZ, distanceTraveled, bestLap, lastLap, currentLap, currentRaceTime, drivingLine sql.NullFloat64
	err := row.Scan(
		&sample.TimeMS, &sample.SpeedKmh, &sample.Rpm, &sample.RpmRatio, &sample.Gear, &sample.Throttle01, &sample.Brake01, &sample.Steer01,
		&sample.FrontSlipRatioAvg, &sample.RearSlipRatioAvg, &sample.FrontCombinedSlipAvg, &sample.RearCombinedSlipAvg,
		&sample.TireTempFrontAvg, &sample.TireTempRearAvg, &sample.SuspensionFrontAvg, &sample.SuspensionRearAvg,
		&sample.YawRate, &sample.PitchRate, &sample.RollRate, &speedField, &velocitySpeed, &speedSource,
		&gameMode, &isRaceOn, &positionX, &positionY, &positionZ, &distanceTraveled, &bestLap, &lastLap,
		&currentLap, &currentRaceTime, &lapNumber, &racePosition, &drivingLine,
		&carOrdinal, &carCategory, &carClass, &carPI, &drivetrain, &numCylinders,
	)
	if err != nil {
		return sample, err
	}
	if speedField.Valid {
		sample.SpeedFieldKmh = speedField.Float64
	}
	if velocitySpeed.Valid {
		sample.VelocitySpeedKmh = velocitySpeed.Float64
	}
	sample.SpeedSource = speedSource.String
	sample.GameMode = telemetry.NormalizeGameMode(gameMode.String)
	sample.IsRaceOn = isRaceOn.Valid && isRaceOn.Int64 != 0
	sample.PositionX = floatFromNull(positionX)
	sample.PositionY = floatFromNull(positionY)
	sample.PositionZ = floatFromNull(positionZ)
	sample.DistanceTraveled = floatFromNull(distanceTraveled)
	sample.BestLap = floatFromNull(bestLap)
	sample.LastLap = floatFromNull(lastLap)
	sample.CurrentLap = floatFromNull(currentLap)
	sample.CurrentRaceTime = floatFromNull(currentRaceTime)
	sample.LapNumber = intFromNull(lapNumber)
	sample.RacePosition = intFromNull(racePosition)
	sample.DrivingLine01 = floatFromNull(drivingLine)
	sample.CarOrdinal = intFromNull(carOrdinal)
	sample.CarCategory = intFromNull(carCategory)
	sample.CarCategoryName = telemetry.CarCategoryName(sample.CarCategory)
	sample.CarClass = carClass.String
	sample.CarPI = intFromNull(carPI)
	sample.Drivetrain = drivetrain.String
	sample.NumCylinders = intFromNull(numCylinders)
	return sample, nil
}

func normalizeBenchmarkTrackInput(input BenchmarkTrackInput) (BenchmarkTrackInput, error) {
	input.Name = strings.TrimSpace(input.Name)
	if input.Name == "" {
		return input, errors.New("benchmark track name is required")
	}
	if len(input.Polyline) < 2 {
		return input, errors.New("benchmark track requires at least two route points")
	}
	if input.StartRadius <= 0 {
		input.StartRadius = defaultGateRadius
	}
	if input.EndRadius <= 0 {
		input.EndRadius = defaultGateRadius
	}
	input.SourceMode = telemetry.NormalizeGameMode(input.SourceMode)
	input.TrackType = normalizeBenchmarkTrackType(input.TrackType)
	if input.TrackType == benchmarkTrackTypeAuto {
		input.TrackType = inferBenchmarkTrackType(input.Polyline)
	}
	if input.TrackType == benchmarkTrackTypeCircuit {
		polyline, observed := extractFirstCircuitPolyline(input.Polyline, input.StartGate)
		input.Polyline = polyline
		if input.LapCountObserved == 0 {
			input.LapCountObserved = observed
		}
	}
	if len(input.Polyline) < 2 {
		return input, errors.New("benchmark track route is too short")
	}
	input.Start = input.Polyline[0]
	if input.TrackType == benchmarkTrackTypeCircuit {
		input.End = input.Start
	} else {
		input.End = input.Polyline[len(input.Polyline)-1]
	}
	dirX, dirZ := input.DirectionX, input.DirectionZ
	if math.Hypot(dirX, dirZ) < 0.001 {
		dirX, dirZ = input.StartGate.DirectionX, input.StartGate.DirectionZ
	}
	if math.Hypot(dirX, dirZ) < 0.001 {
		dirX, dirZ = directionFromPolyline(input.Polyline)
	}
	dirLen := math.Hypot(dirX, dirZ)
	if dirLen < 0.001 {
		return input, errors.New("benchmark track direction cannot be determined")
	}
	input.DirectionX = dirX / dirLen
	input.DirectionZ = dirZ / dirLen
	input.StartGate = normalizeBenchmarkGate(input.StartGate, input.Start, input.DirectionX, input.DirectionZ)
	if input.TrackType == benchmarkTrackTypeCircuit {
		input.FinishGate = input.StartGate
	} else {
		endDirX, endDirZ := directionAtEnd(input.Polyline)
		if math.Hypot(input.FinishGate.DirectionX, input.FinishGate.DirectionZ) < 0.001 {
			input.FinishGate.DirectionX = endDirX
			input.FinishGate.DirectionZ = endDirZ
		}
		input.FinishGate = normalizeBenchmarkGate(input.FinishGate, input.End, input.FinishGate.DirectionX, input.FinishGate.DirectionZ)
	}
	if len(input.Checkpoints) == 0 {
		input.Checkpoints = benchmarkCheckpoints(input.Polyline)
	}
	input.RouteLengthMeters = routeLength(input.Polyline)
	input.Notes = strings.TrimSpace(input.Notes)
	return input, nil
}

func buildBenchmarkTrackInputFromSamples(name string, sourceMode string, samples []telemetry.NormalizedTelemetry, extraction BenchmarkTrackExtractionInput) (BenchmarkTrackInput, error) {
	points := sampledRoutePointsFromSamples(samples)
	if len(points) < 2 {
		return BenchmarkTrackInput{}, errors.New("session does not contain enough position samples to build a benchmark track")
	}
	hasDrivingLine := false
	modeCounts := map[string]int{}
	for _, sample := range samples {
		if math.Abs(sample.DrivingLine01) > 0.05 {
			hasDrivingLine = true
		}
		mode := telemetry.NormalizeGameMode(sample.GameMode)
		if mode != telemetry.GameModeUnknown {
			modeCounts[mode]++
		}
	}
	if telemetry.NormalizeGameMode(sourceMode) == telemetry.GameModeUnknown {
		sourceMode = dominantGameMode(modeCounts)
	}
	trackType := normalizeBenchmarkTrackType(extraction.TrackType)
	rawPoints := sampledPointsOnly(points)
	if trackType == benchmarkTrackTypeAuto {
		trackType = inferBenchmarkTrackType(rawPoints)
	}
	mode := normalizeBenchmarkExtractionMode(extraction.ExtractionMode)
	segment := rawPoints
	lapCount := 0
	if mode != benchmarkExtractionFullSegment {
		switch trackType {
		case benchmarkTrackTypeCircuit:
			var err error
			segment, lapCount, err = extractCircuitSegment(points, extraction, mode)
			if err != nil {
				return BenchmarkTrackInput{}, err
			}
		case benchmarkTrackTypeSprint:
			segment = extractSprintSegment(points, extraction)
		}
	}
	polyline := simplifyPolyline(segment, 8, 800)
	if len(polyline) < 2 {
		return BenchmarkTrackInput{}, errors.New("session route is too short to build a benchmark track")
	}
	startGate := BenchmarkGate{}
	if extraction.StartGate != nil {
		startGate = *extraction.StartGate
	}
	finishGate := BenchmarkGate{}
	if extraction.FinishGate != nil {
		finishGate = *extraction.FinishGate
	}
	return BenchmarkTrackInput{
		Name:             name,
		SourceMode:       sourceMode,
		TrackType:        trackType,
		StartRadius:      defaultGateRadius,
		EndRadius:        defaultGateRadius,
		StartGate:        startGate,
		FinishGate:       finishGate,
		HasDrivingLine:   hasDrivingLine,
		Polyline:         polyline,
		LapCountObserved: lapCount,
	}, nil
}

func routePointsFromSamples(samples []telemetry.NormalizedTelemetry) []BenchmarkPoint {
	points := make([]BenchmarkPoint, 0, len(samples))
	var last BenchmarkPoint
	for _, sample := range samples {
		if !sampleHasVehiclePosition(sample) {
			continue
		}
		point := BenchmarkPoint{X: sample.PositionX, Y: sample.PositionY, Z: sample.PositionZ}
		if len(points) == 0 || distanceXZ(point, last) >= 1 || sample.SpeedKmh > 1 {
			points = append(points, point)
			last = point
		}
	}
	return points
}

type sampledBenchmarkPoint struct {
	Point  BenchmarkPoint
	Sample telemetry.NormalizedTelemetry
	Index  int
}

func sampledRoutePointsFromSamples(samples []telemetry.NormalizedTelemetry) []sampledBenchmarkPoint {
	points := make([]sampledBenchmarkPoint, 0, len(samples))
	var last BenchmarkPoint
	for i, sample := range samples {
		if !sampleHasVehiclePosition(sample) {
			continue
		}
		point := BenchmarkPoint{X: sample.PositionX, Y: sample.PositionY, Z: sample.PositionZ}
		if len(points) == 0 || distanceXZ(point, last) >= 1 || sample.SpeedKmh > 1 {
			points = append(points, sampledBenchmarkPoint{Point: point, Sample: sample, Index: i})
			last = point
		}
	}
	return points
}

func sampledPointsOnly(points []sampledBenchmarkPoint) []BenchmarkPoint {
	out := make([]BenchmarkPoint, 0, len(points))
	for _, point := range points {
		out = append(out, point.Point)
	}
	return out
}

func analyzeBenchmarkRuns(sessionID int64, track BenchmarkTrack, samples []telemetry.NormalizedTelemetry) []BenchmarkRun {
	if len(track.Polyline) < 2 {
		return nil
	}
	points := sampledRoutePointsFromSamples(samples)
	if len(points) < 2 {
		return nil
	}
	checkpoints := track.Checkpoints
	if len(checkpoints) == 0 {
		checkpoints = benchmarkCheckpoints(track.Polyline)
	}
	checkpointRadius := math.Max(50, math.Max(track.StartGate.WidthMeters, track.FinishGate.WidthMeters)*1.6)
	startGate := normalizeBenchmarkGate(track.StartGate, track.Start, track.DirectionX, track.DirectionZ)
	finishGate := normalizeBenchmarkGate(track.FinishGate, track.End, track.DirectionX, track.DirectionZ)
	if track.TrackType == benchmarkTrackTypeCircuit {
		finishGate = startGate
	}

	var runs []BenchmarkRun
	var startPoint sampledBenchmarkPoint
	startPointIndex := 0
	running := false
	runHasDrivingLine := false
	maxSpeed := 0.0
	speedSum := 0.0
	speedCount := 0
	runDistance := 0.0
	checkpointIndex := 0
	var runPoints []sampledBenchmarkPoint

	for i := 1; i < len(points); i++ {
		prev := points[i-1]
		current := points[i]
		if !running {
			if gateCrossed(prev.Point, current.Point, startGate) {
				startPoint = prev
				startPointIndex = i - 1
				running = true
				runHasDrivingLine = math.Abs(prev.Sample.DrivingLine01) > 0.05 || math.Abs(current.Sample.DrivingLine01) > 0.05
				maxSpeed = math.Max(prev.Sample.SpeedKmh, current.Sample.SpeedKmh)
				speedSum = prev.Sample.SpeedKmh + current.Sample.SpeedKmh
				speedCount = 2
				runDistance = distanceXZ(prev.Point, current.Point)
				checkpointIndex = advanceCheckpointIndex(checkpointIndex, checkpoints, current.Point, checkpointRadius)
				runPoints = []sampledBenchmarkPoint{prev, current}
				continue
			}
			continue
		}

		runDistance += distanceXZ(prev.Point, current.Point)
		runPoints = append(runPoints, current)
		if math.Abs(current.Sample.DrivingLine01) > 0.05 {
			runHasDrivingLine = true
		}
		speedSum += current.Sample.SpeedKmh
		speedCount++
		if current.Sample.SpeedKmh > maxSpeed {
			maxSpeed = current.Sample.SpeedKmh
		}
		checkpointIndex = advanceCheckpointIndex(checkpointIndex, checkpoints, current.Point, checkpointRadius)

		finished := gateCrossed(prev.Point, current.Point, finishGate)
		if track.TrackType == benchmarkTrackTypeCircuit && !finished {
			finished = closeCircuitReturn(current.Point, finishGate, runDistance)
		}
		if !finished {
			continue
		}
		duration := sampleTimeMS(current.Sample, current.Index) - sampleTimeMS(startPoint.Sample, startPoint.Index)
		if duration <= 0 {
			duration = int64(i-startPointIndex) * 100
		}
		if duration <= 0 {
			resetBenchmarkRunState(&running, &runHasDrivingLine, &maxSpeed, &speedSum, &speedCount, &runDistance, &checkpointIndex)
			runPoints = nil
			continue
		}
		if track.TrackType == benchmarkTrackTypeCircuit && !validCircuitRun(track, duration, runDistance) {
			continue
		}
		run, ok := buildBenchmarkRun(sessionID, track, startPoint, current, runPoints, duration, speedSum, speedCount, maxSpeed, runHasDrivingLine)
		if ok {
			avg := 0.0
			if speedCount > 0 {
				avg = speedSum / float64(speedCount)
			}
			run.AvgSpeedKmh = &avg
			run.MaxSpeedKmh = &maxSpeed
			runs = append(runs, run)
		}
		if track.TrackType != benchmarkTrackTypeCircuit {
			resetBenchmarkRunState(&running, &runHasDrivingLine, &maxSpeed, &speedSum, &speedCount, &runDistance, &checkpointIndex)
			runPoints = nil
			continue
		}
		startPoint = current
		startPointIndex = i
		runHasDrivingLine = math.Abs(current.Sample.DrivingLine01) > 0.05
		maxSpeed = current.Sample.SpeedKmh
		speedSum = current.Sample.SpeedKmh
		speedCount = 1
		runDistance = 0
		checkpointIndex = 0
		runPoints = []sampledBenchmarkPoint{current}
	}
	return runs
}

func buildBenchmarkRun(sessionID int64, track BenchmarkTrack, startPoint sampledBenchmarkPoint, endPoint sampledBenchmarkPoint, runPoints []sampledBenchmarkPoint, duration int64, speedSum float64, speedCount int, maxSpeed float64, runHasDrivingLine bool) (BenchmarkRun, bool) {
	if len(runPoints) < 2 || len(track.Polyline) < 2 {
		return BenchmarkRun{}, false
	}
	diagnostics := benchmarkRunDiagnostics(track, runPoints)
	if diagnostics.RouteProgress01 < 0.82 {
		diagnostics.addWarning("route_progress_low")
		return BenchmarkRun{}, false
	}
	if track.RouteLengthMeters > 0 && (diagnostics.GeometryLengthMeters < track.RouteLengthMeters*0.75 || diagnostics.GeometryLengthMeters > track.RouteLengthMeters*1.35) {
		diagnostics.addWarning("geometry_length_mismatch")
		return BenchmarkRun{}, false
	}
	if diagnostics.DistanceTraveledDeltaMeters != nil && diagnostics.GeometryLengthMeters > 1 {
		deltaError := math.Abs(*diagnostics.DistanceTraveledDeltaMeters-diagnostics.GeometryLengthMeters) / diagnostics.GeometryLengthMeters
		if deltaError > 0.20 {
			diagnostics.addWarning("distance_traveled_mismatch")
		}
	}
	if diagnostics.AvgLateralErrorMeters > 35 || diagnostics.MaxLateralErrorMeters > 120 {
		diagnostics.addWarning("route_deviation")
	}

	routeScore := clamp01(diagnostics.RouteProgress01)
	lateralScore := clamp01(1 - diagnostics.AvgLateralErrorMeters/60)
	lengthScore := clamp01(1 - math.Abs(diagnostics.TrackLengthErrorPct)/35)
	lineScore := 0.0
	if track.HasDrivingLine && runHasDrivingLine {
		lineScore = 1
	}
	confidence := 0.30 + 0.35*routeScore + 0.20*lateralScore + 0.10*lineScore + 0.05*lengthScore
	confidence = clamp01(confidence)
	if confidence < 0.70 {
		return BenchmarkRun{}, false
	}
	avg := 0.0
	if speedCount > 0 {
		avg = speedSum / float64(speedCount)
	}
	return BenchmarkRun{
		SessionID:                   sessionID,
		TrackID:                     track.ID,
		StartMS:                     sampleTimeMS(startPoint.Sample, startPoint.Index),
		EndMS:                       sampleTimeMS(endPoint.Sample, endPoint.Index),
		DurationMS:                  duration,
		Confidence:                  confidence,
		AvgSpeedKmh:                 &avg,
		MaxSpeedKmh:                 &maxSpeed,
		RouteProgress01:             float64Ptr(diagnostics.RouteProgress01),
		GeometryLengthMeters:        float64Ptr(diagnostics.GeometryLengthMeters),
		TrackLengthErrorPct:         float64Ptr(diagnostics.TrackLengthErrorPct),
		DistanceTraveledDeltaMeters: diagnostics.DistanceTraveledDeltaMeters,
		CurrentRaceTimeDeltaSeconds: diagnostics.CurrentRaceTimeDeltaSeconds,
		AvgLateralErrorMeters:       float64Ptr(diagnostics.AvgLateralErrorMeters),
		MaxLateralErrorMeters:       float64Ptr(diagnostics.MaxLateralErrorMeters),
		WarningFlags:                strings.Join(diagnostics.WarningFlags, ","),
		DriverMode:                  "unknown",
		Valid:                       true,
	}, true
}

type benchmarkRunDiagnosticValues struct {
	RouteProgress01             float64
	GeometryLengthMeters        float64
	TrackLengthErrorPct         float64
	DistanceTraveledDeltaMeters *float64
	CurrentRaceTimeDeltaSeconds *float64
	AvgLateralErrorMeters       float64
	MaxLateralErrorMeters       float64
	WarningFlags                []string
}

func (values *benchmarkRunDiagnosticValues) addWarning(flag string) {
	for _, existing := range values.WarningFlags {
		if existing == flag {
			return
		}
	}
	values.WarningFlags = append(values.WarningFlags, flag)
}

func benchmarkRunDiagnostics(track BenchmarkTrack, runPoints []sampledBenchmarkPoint) benchmarkRunDiagnosticValues {
	values := benchmarkRunDiagnosticValues{}
	values.GeometryLengthMeters = sampledRouteLength(runPoints)
	if track.RouteLengthMeters > 0 {
		values.TrackLengthErrorPct = (values.GeometryLengthMeters/track.RouteLengthMeters - 1) * 100
	}
	if len(runPoints) >= 2 {
		first := runPoints[0].Sample
		last := runPoints[len(runPoints)-1].Sample
		if finiteFloat(first.DistanceTraveled) && finiteFloat(last.DistanceTraveled) {
			delta := last.DistanceTraveled - first.DistanceTraveled
			if delta >= 0 {
				values.DistanceTraveledDeltaMeters = &delta
			}
		}
		if finiteFloat(first.CurrentRaceTime) && finiteFloat(last.CurrentRaceTime) {
			delta := last.CurrentRaceTime - first.CurrentRaceTime
			if delta >= 0 {
				values.CurrentRaceTimeDeltaSeconds = &delta
			}
		}
	}
	projection := newTrackProjection(track.Polyline)
	maxProgress := 0.0
	lateralSum := 0.0
	lateralCount := 0
	for _, point := range runPoints {
		progress, lateral := projection.project(point.Point)
		if progress > maxProgress {
			maxProgress = progress
		}
		lateralSum += lateral
		lateralCount++
		if lateral > values.MaxLateralErrorMeters {
			values.MaxLateralErrorMeters = lateral
		}
	}
	values.RouteProgress01 = maxProgress
	if lateralCount > 0 {
		values.AvgLateralErrorMeters = lateralSum / float64(lateralCount)
	}
	return values
}

type trackProjection struct {
	points     []BenchmarkPoint
	cumulative []float64
	length     float64
}

func newTrackProjection(points []BenchmarkPoint) trackProjection {
	projection := trackProjection{points: points, cumulative: make([]float64, len(points))}
	for i := 1; i < len(points); i++ {
		projection.length += distanceXZ(points[i-1], points[i])
		projection.cumulative[i] = projection.length
	}
	return projection
}

func (projection trackProjection) project(point BenchmarkPoint) (float64, float64) {
	if len(projection.points) < 2 || projection.length <= 0 {
		return 0, 0
	}
	bestDistance := math.MaxFloat64
	bestProgressMeters := 0.0
	for i := 1; i < len(projection.points); i++ {
		a := projection.points[i-1]
		b := projection.points[i]
		segLen := distanceXZ(a, b)
		if segLen <= 0 {
			continue
		}
		t := ((point.X-a.X)*(b.X-a.X) + (point.Z-a.Z)*(b.Z-a.Z)) / (segLen * segLen)
		t = clamp01(t)
		closest := BenchmarkPoint{X: a.X + (b.X-a.X)*t, Z: a.Z + (b.Z-a.Z)*t}
		distance := distanceXZ(point, closest)
		if distance < bestDistance {
			bestDistance = distance
			bestProgressMeters = projection.cumulative[i-1] + segLen*t
		}
	}
	if bestDistance == math.MaxFloat64 {
		bestDistance = 0
	}
	return clamp01(bestProgressMeters / projection.length), bestDistance
}

func sampledRouteLength(points []sampledBenchmarkPoint) float64 {
	total := 0.0
	for i := 1; i < len(points); i++ {
		total += distanceXZ(points[i-1].Point, points[i].Point)
	}
	return total
}

func clamp01(value float64) float64 {
	if value < 0 {
		return 0
	}
	if value > 1 {
		return 1
	}
	return value
}

func float64Ptr(value float64) *float64 {
	return &value
}

func normalizeBenchmarkTrackType(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case benchmarkTrackTypeCircuit:
		return benchmarkTrackTypeCircuit
	case benchmarkTrackTypeSprint:
		return benchmarkTrackTypeSprint
	case benchmarkTrackTypeAuto, "":
		return benchmarkTrackTypeAuto
	default:
		return benchmarkTrackTypeAuto
	}
}

func normalizeBenchmarkExtractionMode(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case benchmarkExtractionFirstLap:
		return benchmarkExtractionFirstLap
	case benchmarkExtractionFullSegment:
		return benchmarkExtractionFullSegment
	case benchmarkExtractionAutoBestLap, "":
		return benchmarkExtractionAutoBestLap
	default:
		return benchmarkExtractionAutoBestLap
	}
}

func inferBenchmarkTrackType(points []BenchmarkPoint) string {
	if len(points) < 2 {
		return benchmarkTrackTypeSprint
	}
	if distanceXZ(points[0], points[len(points)-1]) <= circuitAutoDistance && routeLength(points) >= minCircuitRouteMeters {
		return benchmarkTrackTypeCircuit
	}
	return benchmarkTrackTypeSprint
}

func extractFirstCircuitPolyline(points []BenchmarkPoint, gate BenchmarkGate) ([]BenchmarkPoint, int) {
	if len(points) < 3 {
		return points, 0
	}
	dirX, dirZ := directionFromPolyline(points)
	startGate := normalizeBenchmarkGate(gate, points[0], dirX, dirZ)
	distanceSinceStart := 0.0
	observed := 0
	for i := 1; i < len(points); i++ {
		distanceSinceStart += distanceXZ(points[i-1], points[i])
		if distanceSinceStart < minCircuitRouteMeters {
			continue
		}
		if gateCrossed(points[i-1], points[i], startGate) || closeCircuitReturn(points[i], startGate, distanceSinceStart) {
			observed++
			return points[:i+1], observed
		}
	}
	if distanceXZ(points[0], points[len(points)-1]) <= circuitAutoDistance && distanceSinceStart >= minCircuitRouteMeters {
		return points, 1
	}
	return points, observed
}

func extractCircuitSegment(points []sampledBenchmarkPoint, extraction BenchmarkTrackExtractionInput, mode string) ([]BenchmarkPoint, int, error) {
	if len(points) < 3 {
		return nil, 0, errors.New("session does not contain enough route points to extract a circuit")
	}
	if segments := lapNumberCircuitSegments(points); len(segments) > 0 {
		selected := selectCircuitSegment(segments, mode)
		segment := make([]BenchmarkPoint, 0, selected.end-selected.start+1)
		for _, point := range points[selected.start : selected.end+1] {
			segment = append(segment, point.Point)
		}
		return segment, len(segments), nil
	}
	startGate := inferredStartGate(points, extraction.StartGate)
	var segments []circuitSegment
	startIndex := 0
	distanceSinceStart := 0.0
	totalDistance := 0.0
	for i := 1; i < len(points); i++ {
		stepDistance := distanceXZ(points[i-1].Point, points[i].Point)
		totalDistance += stepDistance
		distanceSinceStart += stepDistance
		if distanceSinceStart < minCircuitRouteMeters {
			continue
		}
		if !(gateCrossed(points[i-1].Point, points[i].Point, startGate) || closeCircuitReturn(points[i].Point, startGate, distanceSinceStart)) {
			continue
		}
		duration := sampleTimeMS(points[i].Sample, points[i].Index) - sampleTimeMS(points[startIndex].Sample, points[startIndex].Index)
		minDuration := minCircuitDurationMS
		if totalDistance < 1000 {
			minDuration = 1000
		}
		if duration < minDuration {
			continue
		}
		segments = append(segments, circuitSegment{start: startIndex, end: i, duration: duration})
		startIndex = i
		distanceSinceStart = 0
	}
	if len(segments) == 0 {
		return nil, 0, errors.New("no complete circuit lap was detected in this session")
	}
	selected := selectCircuitSegment(segments, mode)
	segment := make([]BenchmarkPoint, 0, selected.end-selected.start+1)
	for _, point := range points[selected.start : selected.end+1] {
		segment = append(segment, point.Point)
	}
	return segment, len(segments), nil
}

type circuitSegment struct {
	start    int
	end      int
	duration int64
}

func selectCircuitSegment(segments []circuitSegment, mode string) circuitSegment {
	selected := segments[0]
	if mode == benchmarkExtractionAutoBestLap {
		for _, candidate := range segments[1:] {
			if candidate.duration > 0 && candidate.duration < selected.duration {
				selected = candidate
			}
		}
	}
	return selected
}

func lapNumberCircuitSegments(points []sampledBenchmarkPoint) []circuitSegment {
	if len(points) < 3 {
		return nil
	}
	boundaries := make([]int, 0)
	lastLap, ok := validSampleLapNumber(points[0].Sample)
	for i := 1; i < len(points); i++ {
		currentLap, currentOK := validSampleLapNumber(points[i].Sample)
		if currentOK && ok && currentLap > lastLap {
			boundaries = append(boundaries, i)
		}
		if currentOK {
			lastLap = currentLap
			ok = true
		}
	}
	if len(boundaries) < 2 {
		return nil
	}
	segments := make([]circuitSegment, 0, len(boundaries)-1)
	for i := 1; i < len(boundaries); i++ {
		start := boundaries[i-1]
		end := boundaries[i]
		if end-start < 2 {
			continue
		}
		distance := sampledRouteLength(points[start : end+1])
		if distance < minCircuitRouteMeters {
			continue
		}
		duration := sampleTimeMS(points[end].Sample, points[end].Index) - sampleTimeMS(points[start].Sample, points[start].Index)
		minDuration := minCircuitDurationMS
		if distance < 1000 {
			minDuration = 1000
		}
		if duration < minDuration {
			continue
		}
		segments = append(segments, circuitSegment{start: start, end: end, duration: duration})
	}
	return segments
}

func validSampleLapNumber(sample telemetry.NormalizedTelemetry) (int, bool) {
	if sample.LapNumber <= 0 {
		return 0, false
	}
	return sample.LapNumber, true
}

func extractSprintSegment(points []sampledBenchmarkPoint, extraction BenchmarkTrackExtractionInput) []BenchmarkPoint {
	if len(points) < 2 {
		return nil
	}
	startGate := inferredStartGate(points, extraction.StartGate)
	finishGate := inferredFinishGate(points, extraction.FinishGate)
	startIndex := 0
	started := false
	for i := 1; i < len(points); i++ {
		if !started {
			if gateCrossed(points[i-1].Point, points[i].Point, startGate) {
				startIndex = i - 1
				started = true
			}
			continue
		}
		if gateCrossed(points[i-1].Point, points[i].Point, finishGate) {
			segment := make([]BenchmarkPoint, 0, i-startIndex+1)
			for _, point := range points[startIndex : i+1] {
				segment = append(segment, point.Point)
			}
			if len(segment) >= 2 {
				return segment
			}
		}
	}
	return sampledPointsOnly(points)
}

func inferredStartGate(points []sampledBenchmarkPoint, override *BenchmarkGate) BenchmarkGate {
	center := points[0].Point
	dirX, dirZ := directionFromSampledPoints(points)
	if override != nil {
		return normalizeBenchmarkGate(*override, center, dirX, dirZ)
	}
	return normalizeBenchmarkGate(BenchmarkGate{}, center, dirX, dirZ)
}

func inferredFinishGate(points []sampledBenchmarkPoint, override *BenchmarkGate) BenchmarkGate {
	center := points[len(points)-1].Point
	dirX, dirZ := directionAtEnd(sampledPointsOnly(points))
	if override != nil {
		return normalizeBenchmarkGate(*override, center, dirX, dirZ)
	}
	return normalizeBenchmarkGate(BenchmarkGate{}, center, dirX, dirZ)
}

func normalizeBenchmarkGate(gate BenchmarkGate, fallbackCenter BenchmarkPoint, fallbackDirX float64, fallbackDirZ float64) BenchmarkGate {
	if emptyPoint(gate.Center) {
		gate.Center = fallbackCenter
	}
	dirX, dirZ := gate.DirectionX, gate.DirectionZ
	if math.Hypot(dirX, dirZ) < 0.001 {
		dirX, dirZ = fallbackDirX, fallbackDirZ
	}
	dirLen := math.Hypot(dirX, dirZ)
	if dirLen < 0.001 {
		dirX, dirZ, dirLen = 1, 0, 1
	}
	gate.DirectionX = dirX / dirLen
	gate.DirectionZ = dirZ / dirLen
	if gate.WidthMeters <= 0 {
		gate.WidthMeters = defaultGateWidthMeters
	}
	if gate.DepthMeters <= 0 {
		gate.DepthMeters = defaultGateDepthMeters
	}
	return gate
}

func gateCrossed(prev BenchmarkPoint, current BenchmarkPoint, gate BenchmarkGate) bool {
	gate = normalizeBenchmarkGate(gate, gate.Center, gate.DirectionX, gate.DirectionZ)
	movementDot := directionDot(current.X-prev.X, current.Z-prev.Z, gate.DirectionX, gate.DirectionZ)
	if movementDot <= 0.15 {
		return false
	}
	halfWidth := math.Max(gate.WidthMeters/2, 1)
	if math.Min(gateLateralAbs(prev, gate), gateLateralAbs(current, gate)) > halfWidth {
		return false
	}
	prevProgress := gateProgress(prev, gate)
	currentProgress := gateProgress(current, gate)
	if prevProgress < 0 && currentProgress >= 0 {
		return true
	}
	halfDepth := math.Max(gate.DepthMeters/2, 3)
	return math.Abs(prevProgress) <= halfDepth && currentProgress >= halfDepth
}

func gateProgress(point BenchmarkPoint, gate BenchmarkGate) float64 {
	return (point.X-gate.Center.X)*gate.DirectionX + (point.Z-gate.Center.Z)*gate.DirectionZ
}

func gateLateralAbs(point BenchmarkPoint, gate BenchmarkGate) float64 {
	return math.Abs((point.X-gate.Center.X)*(-gate.DirectionZ) + (point.Z-gate.Center.Z)*gate.DirectionX)
}

func closeCircuitReturn(point BenchmarkPoint, gate BenchmarkGate, distanceSinceStart float64) bool {
	if distanceSinceStart < minCircuitRouteMeters {
		return false
	}
	return distanceXZ(point, gate.Center) <= math.Max(gate.WidthMeters/2, circuitAutoDistance)
}

func validCircuitRun(track BenchmarkTrack, duration int64, distance float64) bool {
	minDuration := minCircuitDurationMS
	if track.RouteLengthMeters < 1000 {
		minDuration = 1000
	}
	if duration < minDuration {
		return false
	}
	if track.RouteLengthMeters <= 0 {
		return distance >= minCircuitRouteMeters
	}
	return distance >= track.RouteLengthMeters*0.92 && distance <= track.RouteLengthMeters*1.12
}

func advanceCheckpointIndex(current int, checkpoints []BenchmarkPoint, point BenchmarkPoint, radius float64) int {
	if current >= len(checkpoints) {
		return current
	}
	if distanceXZ(point, checkpoints[current]) <= radius {
		return current + 1
	}
	return current
}

func resetBenchmarkRunState(running *bool, runHasDrivingLine *bool, maxSpeed *float64, speedSum *float64, speedCount *int, runDistance *float64, checkpointIndex *int) {
	*running = false
	*runHasDrivingLine = false
	*maxSpeed = 0
	*speedSum = 0
	*speedCount = 0
	*runDistance = 0
	*checkpointIndex = 0
}

func sampleHasVehiclePosition(sample telemetry.NormalizedTelemetry) bool {
	if sample.CarOrdinal <= 0 {
		return false
	}
	return finiteFloat(sample.PositionX) && finiteFloat(sample.PositionY) && finiteFloat(sample.PositionZ)
}

func finiteFloat(value float64) bool {
	return !math.IsNaN(value) && !math.IsInf(value, 0)
}

func simplifyPolyline(points []BenchmarkPoint, minDistance float64, maxPoints int) []BenchmarkPoint {
	if len(points) <= 2 {
		return points
	}
	out := []BenchmarkPoint{points[0]}
	last := points[0]
	for _, point := range points[1 : len(points)-1] {
		if distanceXZ(point, last) >= minDistance {
			out = append(out, point)
			last = point
		}
	}
	out = append(out, points[len(points)-1])
	if maxPoints <= 0 || len(out) <= maxPoints {
		return out
	}
	reduced := make([]BenchmarkPoint, 0, maxPoints)
	for i := 0; i < maxPoints; i++ {
		idx := int(math.Round(float64(i) * float64(len(out)-1) / float64(maxPoints-1)))
		reduced = append(reduced, out[idx])
	}
	return reduced
}

func routeLength(points []BenchmarkPoint) float64 {
	total := 0.0
	for i := 1; i < len(points); i++ {
		total += distanceXZ(points[i-1], points[i])
	}
	return total
}

func directionFromPolyline(points []BenchmarkPoint) (float64, float64) {
	if len(points) < 2 {
		return 0, 0
	}
	start := points[0]
	for _, point := range points[1:] {
		dx := point.X - start.X
		dz := point.Z - start.Z
		if math.Hypot(dx, dz) > 5 {
			return dx, dz
		}
	}
	end := points[len(points)-1]
	return end.X - start.X, end.Z - start.Z
}

func directionAtEnd(points []BenchmarkPoint) (float64, float64) {
	if len(points) < 2 {
		return 0, 0
	}
	end := points[len(points)-1]
	for i := len(points) - 2; i >= 0; i-- {
		dx := end.X - points[i].X
		dz := end.Z - points[i].Z
		if math.Hypot(dx, dz) > 5 {
			return dx, dz
		}
	}
	start := points[0]
	return end.X - start.X, end.Z - start.Z
}

func directionFromSampledPoints(points []sampledBenchmarkPoint) (float64, float64) {
	if len(points) < 2 {
		return 0, 0
	}
	start := points[0].Point
	for _, point := range points[1:] {
		dx := point.Point.X - start.X
		dz := point.Point.Z - start.Z
		if math.Hypot(dx, dz) > 5 {
			return dx, dz
		}
	}
	end := points[len(points)-1].Point
	return end.X - start.X, end.Z - start.Z
}

func benchmarkCheckpoints(points []BenchmarkPoint) []BenchmarkPoint {
	if len(points) < 5 {
		return nil
	}
	return []BenchmarkPoint{
		points[len(points)/4],
		points[len(points)/2],
		points[(len(points)*3)/4],
	}
}

func checkpointHitRatio(hitCount int, total int) float64 {
	if total == 0 {
		return 1
	}
	if hitCount < 0 {
		hitCount = 0
	}
	if hitCount > total {
		hitCount = total
	}
	return float64(hitCount) / float64(total)
}

func distanceXZ(a BenchmarkPoint, b BenchmarkPoint) float64 {
	return math.Hypot(a.X-b.X, a.Z-b.Z)
}

func directionDot(dx, dz, dirX, dirZ float64) float64 {
	length := math.Hypot(dx, dz)
	if length < 0.001 {
		return 0
	}
	dirLength := math.Hypot(dirX, dirZ)
	if dirLength < 0.001 {
		return 0
	}
	return (dx/length)*(dirX/dirLength) + (dz/length)*(dirZ/dirLength)
}

func sampleTimeMS(sample telemetry.NormalizedTelemetry, index int) int64 {
	if sample.TimeMS > 0 {
		return sample.TimeMS
	}
	return int64(index) * 100
}

func dominantGameMode(counts map[string]int) string {
	best := telemetry.GameModeUnknown
	bestCount := 0
	for _, mode := range []string{telemetry.GameModeRace, telemetry.GameModeFreeRoam, telemetry.GameModeMenu} {
		if counts[mode] > bestCount {
			best = mode
			bestCount = counts[mode]
		}
	}
	return best
}

func (s *Store) upsertCarIdentity(input TuneProfileInput) error {
	if input.CarOrdinal == nil || *input.CarOrdinal <= 0 || strings.TrimSpace(input.CarName) == "" {
		return nil
	}
	_, err := s.db.Exec(`INSERT INTO car_identity(car_ordinal, car_name, updated_at) VALUES(?, ?, ?)
		ON CONFLICT(car_ordinal) DO UPDATE SET car_name = excluded.car_name, updated_at = excluded.updated_at`,
		*input.CarOrdinal, strings.TrimSpace(input.CarName), nowText())
	return err
}

func profileInsertArgs(input TuneProfileInput, createdAt, updatedAt string) []any {
	input = normalizeTuneProfilePower(input)
	return []any{
		strings.TrimSpace(input.CarName), nullableInt(input.CarOrdinal), nullableInt(input.CarCategory), input.CarClass, nullableInt(input.PI), input.Drivetrain, nullableInt(input.NumCylinders), input.UseCase, input.VersionName,
		nullableFloat(input.PowerKW), nullableFloat(input.TorqueNM), nullableFloat(input.WeightKG), nullableFloat(input.FrontWeightPct), nullableFloat(input.PowerToWeightKWPerKG), nullableFloat(input.PeakTorqueRPM), nullableFloat(input.PeakPowerRPM), nullableFloat(input.RedlineRPM), createdAt, updatedAt,
		nullableFloat(input.FrontTirePressure), nullableFloat(input.RearTirePressure), nullableFloat(input.FinalDrive), nullableFloat(input.Gear1), nullableFloat(input.Gear2), nullableFloat(input.Gear3), nullableFloat(input.Gear4), nullableFloat(input.Gear5), nullableFloat(input.Gear6), nullableFloat(input.Gear7), nullableFloat(input.Gear8), nullableFloat(input.Gear9), nullableFloat(input.Gear10),
		nullableFloat(input.FrontCamber), nullableFloat(input.RearCamber), nullableFloat(input.FrontToe), nullableFloat(input.RearToe), nullableFloat(input.Caster), nullableFloat(input.FrontARB), nullableFloat(input.RearARB),
		nullableFloat(input.FrontSpring), nullableFloat(input.RearSpring), nullableFloat(input.FrontRideHeight), nullableFloat(input.RearRideHeight),
		nullableFloat(input.FrontRebound), nullableFloat(input.RearRebound), nullableFloat(input.FrontBump), nullableFloat(input.RearBump),
		nullableFloat(input.FrontAero), nullableFloat(input.RearAero), nullableFloat(input.AeroBalance), nullableFloat(input.BrakeBalance), nullableFloat(input.BrakePressure),
		nullableFloat(input.FrontDiffAccel), nullableFloat(input.FrontDiffDecel), nullableFloat(input.RearDiffAccel), nullableFloat(input.RearDiffDecel), nullableFloat(input.CenterDiffBalance), input.Notes,
	}
}

func profileUpdateArgs(input TuneProfileInput, updatedAt string) []any {
	args := profileInsertArgs(input, "", updatedAt)
	return append(args[:17], args[18:]...)
}

func normalizeTuneProfilePower(input TuneProfileInput) TuneProfileInput {
	input.PowerKW = positiveFloatPtr(input.PowerKW)
	input.TorqueNM = positiveFloatPtr(input.TorqueNM)
	input.WeightKG = positiveFloatPtr(input.WeightKG)
	input.FrontWeightPct = boundedFloatPtr(input.FrontWeightPct, 0, 100)
	input.PeakTorqueRPM = positiveFloatPtr(input.PeakTorqueRPM)
	input.PeakPowerRPM = positiveFloatPtr(input.PeakPowerRPM)
	input.RedlineRPM = positiveFloatPtr(input.RedlineRPM)
	if input.PowerKW != nil && input.WeightKG != nil && *input.PowerKW > 0 && *input.WeightKG > 0 {
		value := math.Round((*input.PowerKW / *input.WeightKG)*10000) / 10000
		input.PowerToWeightKWPerKG = &value
	} else {
		input.PowerToWeightKWPerKG = nil
	}
	return input
}

func positiveFloatPtr(value *float64) *float64 {
	if value == nil || *value <= 0 || math.IsNaN(*value) || math.IsInf(*value, 0) {
		return nil
	}
	return value
}

func boundedFloatPtr(value *float64, min float64, max float64) *float64 {
	value = positiveFloatPtr(value)
	if value == nil {
		return nil
	}
	if *value < min {
		clamped := min
		return &clamped
	}
	if *value > max {
		clamped := max
		return &clamped
	}
	return value
}

func placeholders(count int) string {
	parts := make([]string, count)
	for i := range parts {
		parts[i] = "?"
	}
	return strings.Join(parts, ", ")
}

func benchmarkTrackAuxJSON(input BenchmarkTrackInput) (string, string, string, error) {
	startGateJSON, err := json.Marshal(input.StartGate)
	if err != nil {
		return "", "", "", err
	}
	finishGateJSON, err := json.Marshal(input.FinishGate)
	if err != nil {
		return "", "", "", err
	}
	checkpointsJSON, err := json.Marshal(input.Checkpoints)
	if err != nil {
		return "", "", "", err
	}
	return string(startGateJSON), string(finishGateJSON), string(checkpointsJSON), nil
}

func nullableFloat(value *float64) any {
	if value == nil {
		return nil
	}
	return *value
}

func nullableInt(value *int64) any {
	if value == nil {
		return nil
	}
	return *value
}

func zeroAsNil(value int) any {
	if value == 0 {
		return nil
	}
	return value
}

func emptyPoint(point BenchmarkPoint) bool {
	return math.Abs(point.X) < 0.000001 && math.Abs(point.Y) < 0.000001 && math.Abs(point.Z) < 0.000001
}

func emptyGate(gate BenchmarkGate) bool {
	return emptyPoint(gate.Center) && math.Hypot(gate.DirectionX, gate.DirectionZ) < 0.001 && gate.WidthMeters <= 0 && gate.DepthMeters <= 0
}

func boolInt(value bool) int {
	if value {
		return 1
	}
	return 0
}

func floatPtr(value sql.NullFloat64) *float64 {
	if !value.Valid {
		return nil
	}
	return &value.Float64
}

func floatFromNull(value sql.NullFloat64) float64 {
	if !value.Valid {
		return 0
	}
	return value.Float64
}

func intPtr(value sql.NullInt64) *int64 {
	if !value.Valid {
		return nil
	}
	return &value.Int64
}

func intFromNull(value sql.NullInt64) int {
	if !value.Valid {
		return 0
	}
	return int(value.Int64)
}

func nowText() string {
	return time.Now().UTC().Format(time.RFC3339Nano)
}

func (s *Store) ensureDefaultRuleThresholdProfile() error {
	var count int
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM rule_threshold_profile WHERE is_default = 1`).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	configJSON, err := defaultRuleConfigJSON()
	if err != nil {
		return err
	}
	now := nowText()
	_, err = s.db.Exec(`INSERT INTO rule_threshold_profile(name, car_class, drivetrain, use_case, config_json, is_default, created_at, updated_at)
		VALUES('Default', '', '', '', ?, 1, ?, ?)`, configJSON, now, now)
	return err
}

func (s *Store) ensureRoadRacingRuleThresholdProfile() error {
	const markerKey = "road_racing_rule_profile_initialized"
	var marker string
	err := s.db.QueryRow(`SELECT value FROM app_setting WHERE key = ?`, markerKey).Scan(&marker)
	if err == nil && marker == "1" {
		return nil
	}
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	var count int
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM rule_threshold_profile WHERE name = 'Road Racing' AND LOWER(COALESCE(use_case, '')) = 'road'`).Scan(&count); err != nil {
		return err
	}
	if count == 0 {
		configJSON, err := roadRacingRuleConfigJSON()
		if err != nil {
			return err
		}
		now := nowText()
		if _, err := s.db.Exec(`INSERT INTO rule_threshold_profile(name, car_class, drivetrain, use_case, game_mode, config_json, is_default, created_at, updated_at)
			VALUES('Road Racing', '', '', 'Road', '', ?, 0, ?, ?)`, configJSON, now, now); err != nil {
			return err
		}
	}
	_, err = s.db.Exec(`INSERT INTO app_setting(key, value) VALUES(?, '1')
		ON CONFLICT(key) DO UPDATE SET value = excluded.value`, markerKey)
	return err
}

func defaultRuleConfigJSON() (string, error) {
	data, err := json.MarshalIndent(telemetry.DefaultRuleConfig(), "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func roadRacingRuleConfigJSON() (string, error) {
	data, err := json.MarshalIndent(telemetry.RoadRacingRuleConfig(), "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func ruleConfigJSONForProfile(profile *RuleThresholdProfile) (string, error) {
	if profile != nil && !profile.IsDefault && strings.EqualFold(strings.TrimSpace(profile.UseCase), "Road") {
		return roadRacingRuleConfigJSON()
	}
	return defaultRuleConfigJSON()
}

func normalizeRuleConfigJSON(value string) (string, error) {
	if strings.TrimSpace(value) == "" {
		return defaultRuleConfigJSON()
	}
	config, err := parseRuleConfigJSON(value)
	if err != nil {
		return "", err
	}
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func parseRuleConfigJSON(value string) (telemetry.RuleConfig, error) {
	var config telemetry.RuleConfig
	if err := json.Unmarshal([]byte(value), &config); err != nil {
		return telemetry.RuleConfig{}, err
	}
	return telemetry.NormalizeRuleConfig(config), nil
}

func normalizeOptionalGameMode(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	return telemetry.NormalizeGameMode(value)
}

func ruleProfileMatchScore(candidate *RuleThresholdProfile, profile *TuneProfile) (int, bool) {
	if candidate == nil {
		return 0, false
	}
	if candidate.IsDefault {
		return 0, true
	}
	if profile == nil {
		return 0, false
	}
	score := 0
	if !matchOptional(candidate.Drivetrain, profile.Drivetrain) {
		return 0, false
	}
	if strings.TrimSpace(candidate.Drivetrain) != "" {
		score += 1
	}
	if !matchOptional(candidate.CarClass, profile.CarClass) {
		return 0, false
	}
	if strings.TrimSpace(candidate.CarClass) != "" {
		score += 2
	}
	if !matchOptional(candidate.UseCase, profile.UseCase) {
		return 0, false
	}
	if strings.TrimSpace(candidate.UseCase) != "" {
		score += 4
	}
	return score, true
}

func matchOptional(candidate string, actual string) bool {
	candidate = strings.TrimSpace(strings.ToLower(candidate))
	if candidate == "" {
		return true
	}
	return candidate == strings.TrimSpace(strings.ToLower(actual))
}
