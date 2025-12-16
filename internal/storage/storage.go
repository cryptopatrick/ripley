// Package storage provides SQLite-based persistence for benchmark results.
package storage

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// BenchmarkRecord represents a single benchmark execution result.
type BenchmarkRecord struct {
	Name       string
	Passed     bool
	TokensUsed int
	Duration   time.Duration
	Quote      string
	Output     string
	Timestamp  time.Time
}

// Storage wraps a SQLite database connection for benchmark data.
type Storage struct {
	db *sql.DB
}

const schema = `
CREATE TABLE IF NOT EXISTS benchmarks (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	passed BOOLEAN NOT NULL,
	tokens_used INTEGER NOT NULL,
	duration_ms INTEGER NOT NULL,
	quote TEXT NOT NULL,
	output TEXT,
	timestamp DATETIME NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_benchmarks_name ON benchmarks(name);
CREATE INDEX IF NOT EXISTS idx_benchmarks_timestamp ON benchmarks(timestamp);
`

// New creates or opens a SQLite database at the given path and initializes the schema.
// Returns a Storage instance ready for use.
func New(dbPath string) (*Storage, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Initialize schema
	if _, err := db.Exec(schema); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return &Storage{db: db}, nil
}

// Close closes the underlying database connection.
func (s *Storage) Close() error {
	return s.db.Close()
}

// InsertRecord saves a benchmark result to the database.
func (s *Storage) InsertRecord(record BenchmarkRecord) error {
	query := `
		INSERT INTO benchmarks (name, passed, tokens_used, duration_ms, quote, output, timestamp)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err := s.db.Exec(
		query,
		record.Name,
		record.Passed,
		record.TokensUsed,
		record.Duration.Milliseconds(),
		record.Quote,
		record.Output,
		record.Timestamp,
	)

	if err != nil {
		return fmt.Errorf("failed to insert record: %w", err)
	}

	return nil
}

// GetRollingStats computes aggregate statistics for a benchmark over the last N runs.
// Returns average tokens used, average duration in seconds, pass rate (0.0-1.0), and any error.
func (s *Storage) GetRollingStats(benchmarkName string, window int) (avgTokens, avgDuration, passRate float64, err error) {
	query := `
		SELECT
			COALESCE(AVG(tokens_used), 0) as avg_tokens,
			COALESCE(AVG(duration_ms), 0) as avg_duration_ms,
			COALESCE(AVG(CASE WHEN passed = 1 THEN 1.0 ELSE 0.0 END), 0) as pass_rate
		FROM (
			SELECT tokens_used, duration_ms, passed
			FROM benchmarks
			WHERE name = ?
			ORDER BY timestamp DESC
			LIMIT ?
		)
	`

	row := s.db.QueryRow(query, benchmarkName, window)

	var avgDurationMs float64
	err = row.Scan(&avgTokens, &avgDurationMs, &passRate)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("failed to query rolling stats: %w", err)
	}

	// Convert milliseconds to seconds for duration
	avgDuration = avgDurationMs / 1000.0

	return avgTokens, avgDuration, passRate, nil
}
