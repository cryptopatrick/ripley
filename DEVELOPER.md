# Ripley Developer Guide

This guide provides detailed information for developers who want to extend, modify, or contribute to the Ripley daemon project.

## Table of Contents

- [Architecture Overview](#architecture-overview)
- [Project Structure](#project-structure)
- [Adding New Benchmarks](#adding-new-benchmarks)
- [Extending Effort Categorization](#extending-effort-categorization)
- [Integrating with Other AI Models](#integrating-with-other-ai-models)
- [Storage Schema Details](#storage-schema-details)
- [Configuration System](#configuration-system)
- [Testing](#testing)
- [Contributing Guidelines](#contributing-guidelines)

## Architecture Overview

Ripley follows a modular, layered architecture:

```
┌─────────────────────────────────────────┐
│  main.go / cmd/ripleyctl/main.go        │  Entry Points
├─────────────────────────────────────────┤
│  internal/config/                       │  Configuration Layer
├─────────────────────────────────────────┤
│  internal/checker/                      │  Business Logic
│  - benchmarks.go (definitions)          │
│  - checker.go (execution)               │
├─────────────────────────────────────────┤
│  internal/ripley/                       │  Domain Logic
│  - quotes.go (effort feedback)          │
├─────────────────────────────────────────┤
│  internal/storage/                      │  Persistence Layer
│  - storage.go (SQLite)                  │
└─────────────────────────────────────────┘
```

### Key Components

1. **Checker**: Executes benchmarks against Claude CLI and categorizes effort
2. **Storage**: Persists results to SQLite and provides rolling statistics
3. **Config**: Manages YAML configuration with sensible defaults
4. **Ripley**: Provides themed quotes based on performance
5. **Main/CLI**: Entry points that orchestrate the components

## Project Structure

```
ripley/
├── main.go                           # Daemon entry point
├── cmd/
│   └── ripleyctl/
│       └── main.go                   # CLI tool with warnings
├── internal/
│   ├── checker/
│   │   ├── benchmarks.go             # Benchmark definitions
│   │   ├── checker.go                # Execution logic
│   │   └── checker_test.go           # Tests
│   ├── config/
│   │   ├── config.go                 # Config management
│   │   └── config_test.go            # Tests
│   ├── ripley/
│   │   ├── quotes.go                 # Ripley quotes
│   │   └── quotes_test.go            # Tests
│   └── storage/
│       ├── storage.go                # SQLite persistence
│       └── storage_test.go           # Tests
├── scripts/
│   ├── run-daemon.sh                 # Daemon launcher
│   └── run-cli.sh                    # CLI launcher
├── config.yaml.example               # Config template
├── Makefile                          # Build automation
├── README.md                         # User documentation
└── DEVELOPER.md                      # This file
```

## Adding New Benchmarks

Benchmarks are defined in `internal/checker/benchmarks.go`. Each benchmark tests a specific capability.

### Step 1: Define the Benchmark

```go
var Benchmarks = []Benchmark{
    // ... existing benchmarks
    {
        Name:        "YourBenchmarkName",
        Prompt:      "The prompt to send to Claude",
        MaxTokens:   20,  // Maximum allowed tokens
        MaxDuration: 10,  // Maximum allowed seconds
    },
}
```

### Step 2: Consider Effort Thresholds

The categorization logic in `checker.go` automatically handles effort scoring:

- **Good**: Passed + within MaxTokens and MaxDuration
- **Medium**: Passed + within 2x MaxTokens and 2x MaxDuration
- **Poor**: Failed or exceeded 2x limits

### Step 3: Test Your Benchmark

Run the daemon and verify the benchmark executes correctly:

```bash
make build
./ripleyd
```

### Example: Adding a JSON Parsing Benchmark

```go
{
    Name:        "JSONParse",
    Prompt:      "Parse this JSON and return the 'name' field: {\"name\":\"Ripley\",\"role\":\"Officer\"}",
    MaxTokens:   10,
    MaxDuration: 5,
}
```

## Extending Effort Categorization

The effort categorization logic is in `internal/checker/checker.go`:

```go
func categorizeEffort(r Result, b Benchmark) string {
    // Failed benchmarks are always poor effort
    if !r.Passed {
        return "poor"
    }

    // Passed within limits is good
    if r.TokensUsed <= b.MaxTokens && r.Duration.Seconds() <= float64(b.MaxDuration) {
        return "good"
    }

    // Medium if slightly exceeded (but still passed)
    if r.TokensUsed <= b.MaxTokens*2 && r.Duration.Seconds() <= float64(b.MaxDuration)*2 {
        return "medium"
    }

    return "poor"
}
```

### Customizing Categorization

You can modify this logic to:

1. **Add more categories**: Introduce "excellent" or "critical" levels
2. **Weight factors differently**: Prioritize duration over tokens
3. **Add custom rules**: Check output format, length, or content

Example with weighted scoring:

```go
func categorizeEffort(r Result, b Benchmark) string {
    if !r.Passed {
        return "poor"
    }

    // Calculate scores (0-100)
    tokenScore := float64(r.TokensUsed) / float64(b.MaxTokens) * 100
    durationScore := r.Duration.Seconds() / float64(b.MaxDuration) * 100

    // Weighted average (60% tokens, 40% duration)
    overallScore := (tokenScore * 0.6) + (durationScore * 0.4)

    switch {
    case overallScore <= 80:
        return "excellent"
    case overallScore <= 100:
        return "good"
    case overallScore <= 150:
        return "medium"
    default:
        return "poor"
    }
}
```

## Integrating with Other AI Models

Currently, Ripley uses the Claude CLI. To support other models:

### Option 1: Parameterize the CLI Command

Modify `internal/checker/checker.go`:

```go
func RunModelBenchmark(b Benchmark, modelCLI string, db *storage.Storage) Result {
    cmd := exec.Command(modelCLI, "--prompt", b.Prompt)
    // ... rest of execution logic
}
```

### Option 2: Create Model-Specific Adapters

```go
// internal/checker/adapters.go
type ModelAdapter interface {
    Execute(prompt string, maxTokens int) (output string, err error)
}

type ClaudeAdapter struct{}
func (a *ClaudeAdapter) Execute(prompt string, maxTokens int) (string, error) {
    // Claude CLI execution
}

type OpenAIAdapter struct{}
func (a *OpenAIAdapter) Execute(prompt string, maxTokens int) (string, error) {
    // OpenAI API call
}
```

### Option 3: Use HTTP APIs Directly

Replace CLI execution with HTTP calls:

```go
func RunAPIBenchmark(b Benchmark, apiURL string, apiKey string, db *storage.Storage) Result {
    req := &APIRequest{
        Prompt:    b.Prompt,
        MaxTokens: b.MaxTokens,
    }

    resp, err := http.Post(apiURL, "application/json", marshalReq(req))
    // ... handle response
}
```

## Storage Schema Details

The SQLite schema is defined in `internal/storage/storage.go`:

```sql
CREATE TABLE benchmarks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    passed BOOLEAN NOT NULL,
    tokens_used INTEGER NOT NULL,
    duration_ms INTEGER NOT NULL,
    quote TEXT NOT NULL,
    output TEXT,
    timestamp DATETIME NOT NULL
);
```

### Indexes

```sql
CREATE INDEX idx_benchmarks_name ON benchmarks(name);
CREATE INDEX idx_benchmarks_timestamp ON benchmarks(timestamp);
```

### Adding Custom Fields

To add new fields (e.g., `model` or `temperature`):

1. Update the `BenchmarkRecord` struct:

```go
type BenchmarkRecord struct {
    Name       string
    Model      string  // New field
    Passed     bool
    TokensUsed int
    // ... other fields
}
```

2. Update the schema:

```sql
CREATE TABLE benchmarks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    model TEXT,  -- New column
    passed BOOLEAN NOT NULL,
    -- ... other columns
);
```

3. Update `InsertRecord` method to include the new field.

### Querying Historical Data

Access the database directly for custom queries:

```bash
sqlite3 ripley.db

# View all results
SELECT name, passed, tokens_used, timestamp FROM benchmarks ORDER BY timestamp DESC LIMIT 20;

# Calculate average performance per benchmark
SELECT name, AVG(tokens_used), AVG(duration_ms), COUNT(*) FROM benchmarks GROUP BY name;

# Find failing benchmarks
SELECT name, output, timestamp FROM benchmarks WHERE passed = 0 ORDER BY timestamp DESC;
```

## Configuration System

Configuration is managed in `internal/config/config.go` using YAML.

### Adding New Config Options

1. Update the `Config` struct:

```go
type Config struct {
    Daemon struct {
        Interval string
        DBPath   string
        Timeout  int  // New option
    }
    // ... other sections
}
```

2. Update `config.yaml.example`:

```yaml
daemon:
  interval: "30m"
  db_path: "./ripley.db"
  timeout: 300  # New option (seconds)
```

3. Update validation in `validate()` method:

```go
func (c *Config) validate() error {
    // ... existing validation

    if c.Daemon.Timeout <= 0 {
        return fmt.Errorf("daemon.timeout must be positive")
    }

    return nil
}
```

4. Update `LoadWithDefaults()`:

```go
func LoadWithDefaults() *Config {
    cfg := &Config{}
    cfg.Daemon.Interval = "30m"
    cfg.Daemon.DBPath = "./ripley.db"
    cfg.Daemon.Timeout = 300  // New default
    // ...
}
```

## Testing

### Running Tests

```bash
# All tests
make test

# Specific package
go test ./internal/checker -v

# With coverage
make test-coverage
```

### Writing Tests

Follow the existing patterns in `*_test.go` files:

```go
func TestYourFeature(t *testing.T) {
    tests := []struct {
        name     string
        input    YourInput
        expected YourOutput
    }{
        {
            name:     "test case 1",
            input:    /* ... */,
            expected: /* ... */,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := YourFunction(tt.input)
            if result != tt.expected {
                t.Errorf("Expected %v, got %v", tt.expected, result)
            }
        })
    }
}
```

### Test Best Practices

1. **Use in-memory SQLite for storage tests**: `:memory:` for speed
2. **Table-driven tests**: Multiple test cases in one function
3. **Mock external dependencies**: Don't call real Claude CLI in tests
4. **Test error paths**: Not just happy paths

## Contributing Guidelines

### Before Submitting

1. **Run tests**: `make test`
2. **Format code**: `make fmt`
3. **Run linter**: `make lint`
4. **Update documentation**: README.md and this file if needed
5. **Add tests**: For new features

### Code Style

- Follow standard Go conventions
- Use `gofmt` for formatting
- Write clear, descriptive variable names
- Add comments for exported functions
- Keep functions focused and small

### Commit Messages

Use conventional commit format:

```
feat: add JSON parsing benchmark
fix: correct effort categorization for edge case
docs: update installation instructions
test: add tests for rolling statistics
refactor: simplify config validation logic
```

### Pull Request Process

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Commit with descriptive messages
5. Push to your fork
6. Open a Pull Request

### Getting Help

- **Issues**: [github.com/cryptopatrick/ripley/issues](https://github.com/cryptopatrick/ripley/issues)
- **Discussions**: Use GitHub Discussions for questions
- **Email**: Contact the maintainer for sensitive topics

## Future Enhancements

Ideas for contribution:

- Docker containerization
- Prometheus metrics export
- Web dashboard with charts
- Slack/Discord notifications
- Multi-model comparison mode
- Benchmark result visualization
- Custom DSL for benchmark definition
- Performance regression detection

---

**Remember**: Like Ripley, we follow procedure, expose problems, and maintain baseline competence. No shortcuts, no excuses.
