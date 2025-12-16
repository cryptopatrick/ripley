Product Requirements Document (PRD)

Project Name: Ripley Daemon Refactor & Documentation AI Agent

Author: [Your Name]
Date: 2025-12-16

1. Purpose

The purpose of this PRD is to define requirements for an AI agent capable of:

Cleaning up the existing Go codebase for the Ripley daemon

Refactoring it for modularity, maintainability, and readability

Writing automated tests to verify functionality

Documenting the code and generating developer-friendly guidance

Preparing the project for collaboration and sharing with other developers

This ensures the Ripley daemon is production-ready, testable, and understandable by a broader development team.

2. Background

The Ripley daemon currently:

Runs liveness and effort benchmarks against Claude Code AI (Sonnet 4.5)

Logs results to SQLite

Prints Ripley-style quotes based on performance

Tracks rolling statistics and alerts when benchmarks underperform

Challenges with the current codebase:

Mix of simulation and real API logic

Limited modularity and inconsistent structure

Sparse inline documentation

No automated test coverage

Hard-coded configuration values (model, token limits, etc.)

Some functions perform multiple responsibilities, violating Separation of Concerns

3. Objectives

The AI agent should:

Code Cleanup & Refactoring

Organize files/modules consistently

Apply Go idiomatic patterns and best practices

Remove duplicate or unused code

Separate concerns: CLI, benchmarking, storage, logging, and quotes

Testing

Generate unit tests for each module (checker, ripley quotes, storage, main)

Include mock Claude Code API for offline testing

Ensure coverage for token counting, effort categorization, and SQLite logging

Documentation

Add clear inline comments for all functions

Generate a README.md with:

Project description

Setup instructions

Usage instructions (daemon & CLI)

Example output

Create a developer guide for extending benchmarks or integrating with other models

Configuration & Portability

Replace hard-coded parameters with a config file or environment variables

Ensure it runs on macOS and Linux

Provide defaults for testing vs. production runs

Code Sharing & Packaging

Structure the repo for easy cloning, building, and contribution

Add .gitignore and go.mod dependencies

Optional: add Makefile or shell script for local development

4. Deliverables

Refactored Go codebase with modular design

Unit tests and test harness for automated verification

Inline code documentation (comments for functions and structs)

README.md and developer guide

Configuration file template (config.yaml or .env)

Example scripts for running daemon and CLI

Optional: instructions for Docker or containerized environment

5. Success Criteria

The project will be considered successful if:

All benchmarks and SQLite logging still function correctly after refactoring

Effort categorization (good, medium, poor) and Ripley quotes remain accurate

The code passes all automated unit tests

Documentation is clear, concise, and sufficient for a new developer to run and extend the project

Configuration is externalized and easy to modify

Repository is ready to share publicly or internally

6. Requirements
Functional Requirements

AI agent must parse and refactor Go code automatically

Generate unit tests with mocked Claude sessions

Insert inline comments and docstrings

Create README and developer guide automatically

Non-Functional Requirements

Refactored code must be idiomatic Go, readable, and maintainable

Tests should run on macOS and Linux

Logging to SQLite must remain fully functional

Effort quotes system must remain intact

7. Constraints

The AI agent must work with existing Go 1.22 code

Must not remove Ripley-specific logic (quotes, effort categorization, rolling statistics)

Should maintain backward compatibility with previous benchmark definitions

Mocking must avoid making real API calls during test execution

8. Timeline & Milestones
Milestone	Description	Estimated Effort
Code Cleanup	Remove duplicates, reorganize files, separate concerns	1 day
Refactor	Apply idiomatic Go patterns, modularize	2 days
Unit Testing	Generate automated tests, include mocks	2 days
Documentation	Inline comments, README, developer guide	1 day
Verification	Run all benchmarks, check SQLite logging, print outputs	0.5 day
Packaging	Final repo structure, .gitignore, config templates	0.5 day
9. Notes

AI agent should treat Ripley quotes as core feature and never remove or simplify them

Effort categorization thresholds should remain configurable

SQLite logging schema should remain compatible with existing data