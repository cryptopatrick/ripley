package checker

import (
	"testing"
	"time"
)

func TestCategorizeEffort(t *testing.T) {
	benchmark := Benchmark{
		Name:        "TestBench",
		Prompt:      "test",
		MaxTokens:   10,
		MaxDuration: 5,
	}

	tests := []struct {
		name     string
		result   Result
		expected string
	}{
		{
			name: "Good effort - passed within limits",
			result: Result{
				Name:       "TestBench",
				Passed:     true,
				TokensUsed: 8,
				Duration:   3 * time.Second,
			},
			expected: "good",
		},
		{
			name: "Good effort - exact limits",
			result: Result{
				Name:       "TestBench",
				Passed:     true,
				TokensUsed: 10,
				Duration:   5 * time.Second,
			},
			expected: "good",
		},
		{
			name: "Medium effort - slightly over tokens",
			result: Result{
				Name:       "TestBench",
				Passed:     true,
				TokensUsed: 15,
				Duration:   6 * time.Second,
			},
			expected: "medium",
		},
		{
			name: "Medium effort - at 2x limit",
			result: Result{
				Name:       "TestBench",
				Passed:     true,
				TokensUsed: 20,
				Duration:   10 * time.Second,
			},
			expected: "medium",
		},
		{
			name: "Poor effort - failed",
			result: Result{
				Name:       "TestBench",
				Passed:     false,
				TokensUsed: 5,
				Duration:   2 * time.Second,
			},
			expected: "poor",
		},
		{
			name: "Poor effort - way over limits",
			result: Result{
				Name:       "TestBench",
				Passed:     true,
				TokensUsed: 50,
				Duration:   20 * time.Second,
			},
			expected: "poor",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			effort := categorizeEffort(tt.result, benchmark)
			if effort != tt.expected {
				t.Errorf("categorizeEffort() = %v, want %v", effort, tt.expected)
			}
		})
	}
}

func TestBenchmarks(t *testing.T) {
	// Verify benchmarks are properly defined
	if len(Benchmarks) == 0 {
		t.Error("No benchmarks defined")
	}

	for _, b := range Benchmarks {
		t.Run(b.Name, func(t *testing.T) {
			if b.Name == "" {
				t.Error("Benchmark name is empty")
			}
			if b.Prompt == "" {
				t.Error("Benchmark prompt is empty")
			}
			if b.MaxTokens <= 0 {
				t.Errorf("MaxTokens should be positive, got %d", b.MaxTokens)
			}
			if b.MaxDuration <= 0 {
				t.Errorf("MaxDuration should be positive, got %d", b.MaxDuration)
			}
		})
	}
}

func TestResultStruct(t *testing.T) {
	// Test that Result struct can be properly initialized
	result := Result{
		Name:       "Test",
		Passed:     true,
		TokensUsed: 10,
		Duration:   time.Second,
		Quote:      "test quote",
		Effort:     "good",
		Output:     "test output",
	}

	if result.Name != "Test" {
		t.Errorf("Expected Name='Test', got '%s'", result.Name)
	}
	if !result.Passed {
		t.Error("Expected Passed=true")
	}
	if result.Effort != "good" {
		t.Errorf("Expected Effort='good', got '%s'", result.Effort)
	}
}
