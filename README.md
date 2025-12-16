<h1 align="center">
  <br>
    <img 
      src="https://github.com/cryptopatrick/factory/blob/master/img/100days/ripley.png" 
      alt="Title" 
      width="200"
    />
  <br>
  RIPLEY<br>
  <h3 align="center">"Are you actually trying, or just giving me surface-level bullshit?"</h3>
</h1>


<p align="center">
    <a href="LICENSE" target="_blank">
    <img src="https://img.shields.io/github/license/sulu/sulu.svg" alt="GitHub license"/>
  </a>
</p>

<b>Authors's bio:</b> 👋😀 Hi, I'm CryptoPatrick! I'm currently enrolled as an 
Undergraduate student in Mathematics, at Chalmers & the University of Gothenburg, Sweden. <br>
If you have any questions or need more info, then please <a href="https://discord.gg/T8EWmJZpCB">join my Discord Channel: AiMath</a>

---


## What is Ripley?
**Ripley** is a lightweight Go daemon for **AI liveness and effort testing**, inspired by Ellen Ripley from *Alien*. Just as Ripley follows procedure and exposes systems that claim they're fine but aren't, this daemon monitors Claude Code AI (Sonnet 4.5) with regular benchmarks and Ripley-style feedback.

> Current version only supports Sonnet 4.5, but Opus is on the way.


## Features

- **Automated Benchmarking**:  
  Runs simple, deterministic tests against Claude CLI at regular intervals
- **Benchmarks are Token Limited**:  
  Ripley cares deeply about your tokens, and does her utmost to avoid wasting even a single token. Bench mark will not burn more than a token limit (default is 200 tokens, but you can set that in config.yaml - see below)
- **Effort Categorization**:  
  Classifies test results as "good", "medium", or "poor" based on token usage, duration, and correctness
- **Ripley-Style Quotes**:  
  Provides feedback on test results, with Ripley's characteristic calm, procedural, and no-nonsense tone
- **SQLite Logging**:  
  Persists all benchmark results for historical analysis
- **Rolling Statistics**:  
  Tracks performance trends over the last N runs
- **Configurable**:  
  YAML-based configuration for intervals, thresholds, and more
- **Extensible**:  
  Easy to add new benchmarks or integrate with other AI models


## Prerequisites

- **Go 1.22+**
- **Claude CLI** installed and accessible in your PATH
  - The daemon uses the Claude Code CLI for benchmarking
  - Install from: [github.com/anthropics/claude-code](https://github.com/anthropics/claude-code)

## Installation

```bash
git clone https://github.com/cryptopatrick/ripley.git
cd ripley
make build
```

This builds two binaries:
- `ripleyd` - The main daemon
- `ripleyctl` - CLI tool with enhanced monitoring

## Configuration

Copy the example configuration file and customize it:

```bash
cp config.yaml.example config.yaml
```

Edit `config.yaml`:

```yaml
daemon:
  interval: "30m"              # How often to run benchmarks
  db_path: "./ripley.db"       # SQLite database location

claude:
  model: "Sonnet"              # Claude model to test
  default_max_tokens: 200      # Default token limit

monitoring:
  rolling_window: 10           # Number of runs for statistics
  warning_threshold: 0.7       # Alert if pass rate < 70%
```

> If no `config.yaml` is found, the daemon uses sensible defaults - which are??? TODO add details.

## Usage

### Running the Daemon

```bash
# Using make
make run-daemon

# Using the script
./scripts/run-daemon.sh

# Direct execution
./ripleyd
```

### Running the CLI Tool

The CLI tool provides the same functionality as the daemon but with additional warnings when performance drops below thresholds:

```bash
# Using make
make run-cli

# Using the script
./scripts/run-cli.sh

# Direct execution
./ripleyctl
```

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage
```
---

## Example Output

```
Ripley daemon started with Sonnet...
Database: ./ripley.db | Interval: 30m0s

=== Running Claude Code Liveness & Effort Check ===
[PASS] Sum1to100 | Effort: good | Tokens: 7 | Duration: 1.2s
Quote: Now we're getting somewhere — that's the baseline competence I expect.
Output: 5050

[PASS] PalindromeCheck | Effort: good | Tokens: 4 | Duration: 0.9s
Quote: Precision, speed, and efficiency — looks like someone's awake.
Output: true

[PASS] SimpleArithmetic | Effort: good | Tokens: 3 | Duration: 0.8s
Quote: You followed procedure. I approve.
Output: 105

[PASS] ListReverse | Effort: medium | Tokens: 12 | Duration: 1.4s
Quote: Decent. I'll allow it… this time.
Output: [5, 4, 3, 2, 1]

=== Rolling Statistics (Last 10 Runs) ===
✓ Sum1to100 | Avg Tokens: 7.2 | Avg Duration: 1.18s | Pass Rate: 100%
✓ PalindromeCheck | Avg Tokens: 4.1 | Avg Duration: 0.92s | Pass Rate: 100%
✓ SimpleArithmetic | Avg Tokens: 3.0 | Avg Duration: 0.85s | Pass Rate: 100%
⚠ ListReverse | Avg Tokens: 13.4 | Avg Duration: 1.52s | Pass Rate: 65%
```

## Database Schema

Results are stored in SQLite with the following schema:

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

## Adding New Benchmarks

Edit `internal/checker/benchmarks.go`:

```go
var Benchmarks = []Benchmark{
    {
        Name:        "YourBenchmark",
        Prompt:      "Your prompt here",
        MaxTokens:   10,
        MaxDuration: 5, // seconds
    },
    // ... existing benchmarks
}
```

## Project Structure

```
ripley/
├── main.go                    # Daemon entry point
├── cmd/
│   └── ripleyctl/             # CLI tool with warnings
├── internal/
│   ├── checker/               # Benchmark execution logic
│   ├── config/                # Configuration management
│   ├── ripley/                # Ripley quotes
│   └── storage/               # SQLite persistence
├── scripts/                   # Helper scripts
├── config.yaml.example        # Configuration template
├── Makefile                   # Build automation
└── README.md                  # This file
```

## Development

See [DEVELOPER.md](./DEVELOPER.md) for detailed development documentation including:
- Architecture overview
- Adding new benchmarks
- Extending effort categorization
- Integrating with other AI models
- Contributing guidelines

## Make Targets

```bash
make build          # Build daemon and CLI binaries
make test           # Run all tests
make test-coverage  # Run tests with coverage report
make clean          # Remove build artifacts
make install        # Install binaries to $GOPATH/bin
make run-daemon     # Build and run the daemon
make run-cli        # Build and run the CLI tool
make fmt            # Format code
make lint           # Run linter
make help           # Show all targets
```

## Author
**CryptoPatrick**

## License

MIT

## Roadmap

- [x] Core benchmarking and storage
- [x] Configuration system
- [x] Rolling statistics
- [x] Automated testing
- [ ] Docker support
- [ ] GitHub Actions CI
- [ ] Slack/Discord notifications
- [ ] Web dashboard
- [ ] Multi-model support (Opus, Haiku)
- [ ] Custom benchmark DSL

## Acknowledgments

Inspired by Ellen Ripley from *Alien* - the only one actually doing any work.🔫
