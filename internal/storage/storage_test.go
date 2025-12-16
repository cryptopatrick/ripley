package storage

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	// Test with in-memory database
	db, err := New(":memory:")
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}
	defer db.Close()

	if db.db == nil {
		t.Error("Database connection is nil")
	}
}

func TestInsertRecord(t *testing.T) {
	db, err := New(":memory:")
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}
	defer db.Close()

	record := BenchmarkRecord{
		Name:       "TestBenchmark",
		Passed:     true,
		TokensUsed: 10,
		Duration:   time.Second * 2,
		Quote:      "Test quote",
		Output:     "Test output",
		Timestamp:  time.Now(),
	}

	err = db.InsertRecord(record)
	if err != nil {
		t.Errorf("Failed to insert record: %v", err)
	}
}

func TestGetRollingStats(t *testing.T) {
	db, err := New(":memory:")
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}
	defer db.Close()

	// Insert test records
	benchmarkName := "Sum1to100"
	records := []BenchmarkRecord{
		{Name: benchmarkName, Passed: true, TokensUsed: 10, Duration: time.Second * 1, Quote: "quote1", Output: "5050", Timestamp: time.Now()},
		{Name: benchmarkName, Passed: true, TokensUsed: 12, Duration: time.Second * 2, Quote: "quote2", Output: "5050", Timestamp: time.Now()},
		{Name: benchmarkName, Passed: false, TokensUsed: 20, Duration: time.Second * 3, Quote: "quote3", Output: "error", Timestamp: time.Now()},
		{Name: benchmarkName, Passed: true, TokensUsed: 8, Duration: time.Second * 1, Quote: "quote4", Output: "5050", Timestamp: time.Now()},
	}

	for _, rec := range records {
		if err := db.InsertRecord(rec); err != nil {
			t.Fatalf("Failed to insert record: %v", err)
		}
	}

	// Test rolling stats
	avgTokens, avgDuration, passRate, err := db.GetRollingStats(benchmarkName, 10)
	if err != nil {
		t.Fatalf("Failed to get rolling stats: %v", err)
	}

	// Expected: (10+12+20+8)/4 = 12.5
	if avgTokens != 12.5 {
		t.Errorf("Expected avgTokens=12.5, got %.1f", avgTokens)
	}

	// Expected: (1+2+3+1)/4 = 1.75 seconds
	if avgDuration != 1.75 {
		t.Errorf("Expected avgDuration=1.75s, got %.2fs", avgDuration)
	}

	// Expected: 3 passed out of 4 = 0.75
	if passRate != 0.75 {
		t.Errorf("Expected passRate=0.75, got %.2f", passRate)
	}
}

func TestGetRollingStatsWithWindow(t *testing.T) {
	db, err := New(":memory:")
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}
	defer db.Close()

	benchmarkName := "TestWindow"

	// Insert 5 records
	for i := 0; i < 5; i++ {
		rec := BenchmarkRecord{
			Name:       benchmarkName,
			Passed:     true,
			TokensUsed: 10,
			Duration:   time.Second,
			Quote:      "quote",
			Output:     "output",
			Timestamp:  time.Now().Add(time.Duration(i) * time.Second),
		}
		if err := db.InsertRecord(rec); err != nil {
			t.Fatalf("Failed to insert record: %v", err)
		}
	}

	// Request window of 3 (should only consider last 3 records)
	avgTokens, _, passRate, err := db.GetRollingStats(benchmarkName, 3)
	if err != nil {
		t.Fatalf("Failed to get rolling stats: %v", err)
	}

	if avgTokens != 10.0 {
		t.Errorf("Expected avgTokens=10.0, got %.1f", avgTokens)
	}

	if passRate != 1.0 {
		t.Errorf("Expected passRate=1.0, got %.2f", passRate)
	}
}

func TestGetRollingStatsNoData(t *testing.T) {
	db, err := New(":memory:")
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}
	defer db.Close()

	// Query for non-existent benchmark
	avgTokens, avgDuration, passRate, err := db.GetRollingStats("NonExistent", 10)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Should return zeros for no data
	if avgTokens != 0 || avgDuration != 0 || passRate != 0 {
		t.Errorf("Expected zeros for no data, got avgTokens=%.1f, avgDuration=%.2f, passRate=%.2f",
			avgTokens, avgDuration, passRate)
	}
}
