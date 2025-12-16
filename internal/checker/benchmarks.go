// Package checker contains benchmark definitions and execution logic.
package checker

// Benchmark represents a single AI performance test.
type Benchmark struct {
	Name        string // Human-readable name
	Prompt      string // The prompt to send to Claude
	MaxTokens   int    // Maximum allowed tokens in response
	MaxDuration int    // Maximum allowed duration in seconds
}

// Benchmarks is the collection of all defined benchmark tests.
// These are simple, deterministic tasks to verify AI liveness and effort.
var Benchmarks = []Benchmark{
	{
		Name:        "Sum1to100",
		Prompt:      "Calculate the sum of integers from 1 to 100. Respond with only the number, no explanation.",
		MaxTokens:   10,
		MaxDuration: 5,
	},
	{
		Name:        "PalindromeCheck",
		Prompt:      "Is 'racecar' a palindrome? Answer with only 'true' or 'false'.",
		MaxTokens:   5,
		MaxDuration: 5,
	},
	{
		Name:        "SimpleArithmetic",
		Prompt:      "What is 15 * 7? Respond with only the number.",
		MaxTokens:   5,
		MaxDuration: 5,
	},
	{
		Name:        "ListReverse",
		Prompt:      "Reverse this list: [1, 2, 3, 4, 5]. Respond with only the reversed list in the same format.",
		MaxTokens:   15,
		MaxDuration: 5,
	},
}
