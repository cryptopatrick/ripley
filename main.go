package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/cryptopatrick/ripley/internal/checker"
	"github.com/cryptopatrick/ripley/internal/config"
	"github.com/cryptopatrick/ripley/internal/storage"
)

func main() {
	// Load configuration
	var cfg *config.Config
	configPath := "config.yaml"

	if _, err := os.Stat(configPath); err == nil {
		cfg, err = config.Load(configPath)
		if err != nil {
			log.Fatalf("Failed to load config: %v", err)
		}
		fmt.Printf("Loaded configuration from %s\n", configPath)
	} else {
		cfg = config.LoadWithDefaults()
		fmt.Println("Using default configuration (config.yaml not found)")
	}

	interval, err := cfg.GetInterval()
	if err != nil {
		log.Fatalf("Invalid interval configuration: %v", err)
	}

	// Initialize storage
	db, err := storage.New(cfg.Daemon.DBPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	fmt.Printf("Ripley daemon started with %s...\n", cfg.Claude.Model)
	fmt.Printf("Database: %s | Interval: %v\n\n", cfg.Daemon.DBPath, interval)

	for {
		fmt.Println("=== Running Claude Code Liveness & Effort Check ===")
		results := checker.RunBenchmarks(db)
		checker.PrintResults(results)

		// Show rolling statistics
		fmt.Printf("\n=== Rolling Statistics (Last %d Runs) ===\n", cfg.Monitoring.RollingWindow)
		for _, b := range checker.Benchmarks {
			avgTokens, avgDuration, passRate, err := db.GetRollingStats(b.Name, cfg.Monitoring.RollingWindow)
			if err != nil {
				log.Printf("Error getting stats for %s: %v", b.Name, err)
				continue
			}

			status := "✓"
			if passRate < cfg.Monitoring.WarningThreshold {
				status = "⚠"
			}

			fmt.Printf("%s %s | Avg Tokens: %.1f | Avg Duration: %.2fs | Pass Rate: %.0f%%\n",
				status, b.Name, avgTokens, avgDuration, passRate*100)
		}
		fmt.Println()

		time.Sleep(interval)
	}
}
