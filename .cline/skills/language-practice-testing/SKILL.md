---
name: language-practice-testing
description: Design and run deterministic tests that teach the public behavior, invariants, failure paths, and limits of language exercises. Use when adding, changing, reviewing, or validating tests or testable behavior.
compatibility: Designed for Go, Python, Rust, and TypeScript projects using the repository's declared or language-standard test tooling.
metadata:
  version: "1"
  scope: "testing and validation"
---

# Testing as a learning tool

Tests should prove a stated contract and help the learner understand the
behavior. Do not maximize line coverage at the expense of useful experiments.

## Before writing tests

1. State the behavior, invariant, or failure that the test should establish.
2. Identify the public boundary under test and avoid coupling assertions to
   private implementation details.
3. Decide whether the test is a unit, integration, command-level, performance,
   or manual learning check.
4. Choose fixed inputs, controlled time, isolated temporary state, and explicit
   cleanup where applicable.
5. Identify what the test cannot prove.

## Required cases

Cover the smallest complete set relevant to the exercise:

- representative success;
- invalid input and expected error behavior;
- boundary values and empty cases;
- malformed or partial external data;
- state preservation after rejection;
- repeated execution, idempotency, or ordering when relevant;
- cancellation, timeout, race, or resource cleanup for concurrent code;
- a regression case for every reproducible bug.

Assert observable behavior, error type/message contract where stable, exit
status, output shape, side effects, and cleanup. Keep test names explanatory.

## Test execution

- Run the repository's documented command after the final relevant edit.
- Do not treat a stale build or previously passing command as current evidence.
- Keep ordinary tests offline, bounded, deterministic, and independent of real
  credentials or production services.
- Use fakes or dependency seams before live integrations.
- Run race detectors, sanitizers, benchmarks, or fresh-process checks only when
  the behavior requires that evidence and report their separate scope.

## Learning checkpoint

After a non-trivial test, explain:

- which behavior it proves;
- why the fixture exposes that behavior;
- which tempting but weaker assertion was avoided;
- what remains unproven.

Do not generate a complete test suite in `handwrite` mode. Provide test cases,
assertion targets, and feedback while the learner writes the core test.