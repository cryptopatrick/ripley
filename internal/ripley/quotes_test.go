package ripley

import (
	"testing"
)

func TestRandomQuoteByEffort(t *testing.T) {
	tests := []struct {
		effort        string
		expectedArray []string
	}{
		{"good", GoodEffortQuotes},
		{"medium", MediumEffortQuotes},
		{"poor", PoorEffortQuotes},
	}

	for _, tt := range tests {
		t.Run(tt.effort, func(t *testing.T) {
			quote := RandomQuoteByEffort(tt.effort)

			// Check that quote is from the correct array
			found := false
			for _, q := range tt.expectedArray {
				if q == quote {
					found = true
					break
				}
			}

			if !found {
				t.Errorf("RandomQuoteByEffort(%s) returned unexpected quote: %s", tt.effort, quote)
			}

			if quote == "" {
				t.Errorf("RandomQuoteByEffort(%s) returned empty string", tt.effort)
			}
		})
	}
}

func TestRandomQuoteByEffortDefault(t *testing.T) {
	// Test with invalid effort value - should return from all categories
	quote := RandomQuoteByEffort("invalid")

	if quote == "" {
		t.Error("RandomQuoteByEffort with invalid effort returned empty string")
	}

	// Verify it's from one of the categories
	allQuotes := append(append(GoodEffortQuotes, MediumEffortQuotes...), PoorEffortQuotes...)
	found := false
	for _, q := range allQuotes {
		if q == quote {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("RandomQuoteByEffort(invalid) returned quote not in any category: %s", quote)
	}
}

func TestQuoteArraysNotEmpty(t *testing.T) {
	if len(GoodEffortQuotes) == 0 {
		t.Error("GoodEffortQuotes is empty")
	}

	if len(MediumEffortQuotes) == 0 {
		t.Error("MediumEffortQuotes is empty")
	}

	if len(PoorEffortQuotes) == 0 {
		t.Error("PoorEffortQuotes is empty")
	}
}

func TestQuotesAreUnique(t *testing.T) {
	// Test that good quotes are unique
	seen := make(map[string]bool)
	for _, quote := range GoodEffortQuotes {
		if seen[quote] {
			t.Errorf("Duplicate quote in GoodEffortQuotes: %s", quote)
		}
		seen[quote] = true
	}

	// Test that medium quotes are unique
	seen = make(map[string]bool)
	for _, quote := range MediumEffortQuotes {
		if seen[quote] {
			t.Errorf("Duplicate quote in MediumEffortQuotes: %s", quote)
		}
		seen[quote] = true
	}

	// Test that poor quotes are unique
	seen = make(map[string]bool)
	for _, quote := range PoorEffortQuotes {
		if seen[quote] {
			t.Errorf("Duplicate quote in PoorEffortQuotes: %s", quote)
		}
		seen[quote] = true
	}
}

func TestQuotesNotEmpty(t *testing.T) {
	// Test that no quote is an empty string
	allQuotes := append(append(GoodEffortQuotes, MediumEffortQuotes...), PoorEffortQuotes...)

	for _, quote := range allQuotes {
		if quote == "" {
			t.Error("Found empty quote in quote arrays")
		}
	}
}
