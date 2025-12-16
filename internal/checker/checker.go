package checker

import (
    "bytes"
    "fmt"
    "os/exec"
    "strings"
    "time"

    "github.com/cryptopatrick/ripley/internal/ripley"
    "github.com/cryptopatrick/ripley/internal/storage"
)

type Result struct {
    Name       string
    Passed     bool
    TokensUsed int
    Duration   time.Duration
    Quote      string
    Effort     string // "good", "medium", "poor"
    Output     string
}

// Determine effort category based on passed status, tokens, and duration
func categorizeEffort(r Result, b Benchmark) string {
    // Failed benchmarks are always poor effort
    if !r.Passed {
        return "poor"
    }

    // Passed within limits is good
    if r.TokensUsed <= b.MaxTokens && r.Duration.Seconds() <= float64(b.MaxDuration) {
        return "good"
    }

    // Medium if slightly exceeded tokens or duration (but still passed)
    if r.TokensUsed <= b.MaxTokens*2 && r.Duration.Seconds() <= float64(b.MaxDuration)*2 {
        return "medium"
    }

    return "poor"
}

// Run a single benchmark using Claude CLI
func RunClaudeBenchmark(b Benchmark, db *storage.Storage) Result {
    start := time.Now()

    cmd := exec.Command(
        "claude",
        "--model", "Sonnet",
        "--fresh",
        "--max-tokens", fmt.Sprintf("%d", b.MaxTokens),
    )
    cmd.Stdin = strings.NewReader(b.Prompt)

    var out bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &out

    err := cmd.Start()
    if err != nil {
        r := Result{Name: b.Name, Passed: false, Effort: "poor", Quote: ripley.RandomQuoteByEffort("poor"), Output: err.Error()}
        saveResult(r, db)
        return r
    }

    done := make(chan error)
    go func() { done <- cmd.Wait() }()

    var r Result
    select {
    case <-time.After(time.Duration(b.MaxDuration) * time.Second):
        _ = cmd.Process.Kill()
        duration := time.Since(start)
        r = Result{Name: b.Name, Passed: false, Duration: duration, Output: "Timed out"}
    case err := <-done:
        duration := time.Since(start)
        output := out.String()
        tokensUsed := len(strings.Fields(output))
        passed := tokensUsed <= b.MaxTokens && duration.Seconds() <= float64(b.MaxDuration)
        if err != nil {
            passed = false
        }

        r = Result{
            Name:       b.Name,
            Passed:     passed,
            TokensUsed: tokensUsed,
            Duration:   duration,
            Output:     strings.TrimSpace(output),
        }
    }

    // Determine effort and assign Ripley quote
    r.Effort = categorizeEffort(r, b)
    r.Quote = ripley.RandomQuoteByEffort(r.Effort)

    saveResult(r, db)
    return r
}

// Save benchmark result to DB if storage is provided
func saveResult(r Result, db *storage.Storage) {
    if db != nil {
        _ = db.InsertRecord(storage.BenchmarkRecord{
            Name:       r.Name,
            Passed:     r.Passed,
            TokensUsed: r.TokensUsed,
            Duration:   r.Duration,
            Quote:      r.Quote,
            Output:     r.Output,
            Timestamp:  time.Now(),
        })
    }
}

// Run all benchmarks
func RunBenchmarks(db *storage.Storage) []Result {
    var results []Result
    for _, b := range Benchmarks {
        results = append(results, RunClaudeBenchmark(b, db))
    }
    return results
}

// Print results with Ripley-style quotes
func PrintResults(results []Result) {
    for _, r := range results {
        status := "PASS"
        if !r.Passed {
            status = "FAIL"
        }
        fmt.Printf("[%s] %s | Effort: %s | Tokens: %d | Duration: %s\nQuote: %s\nOutput: %s\n\n",
            status, r.Name, r.Effort, r.TokensUsed, r.Duration, r.Quote, r.Output)
    }
}
