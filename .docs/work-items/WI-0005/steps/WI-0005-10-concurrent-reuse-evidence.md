---
schema: language-practice-step/v1
id: WI-0005
step: "10"
status: planned
skill: language-practice-testing
depends_on: ["09"]
target: go/basic/calc/math/
---

# WI-0005 step 10 — Verify concurrent AST reuse

## Objective

Prove that one immutable AST can serve multiple concurrent evaluations with
isolated inputs and results. Add bounded stress, race, leak, and sequential-oracle
comparisons without turning a test run into an unbounded load test.

## Commit subject

```text
WI-0005: 10: test(calc): verify concurrent evaluator reuse
```

## Scope and ownership

- Own only `go/basic/calc/math/` and this step's handoff update.
- Add reusable test fixtures, deterministic barriers, bounded concurrency tests,
  race coverage, and leak observations.
- Keep production behavior unchanged unless a reproducible test exposes a
  lifecycle defect within the existing contract.
- Do not claim performance improvement from concurrency in this step.

## Design checkpoint

Record the test workload, concurrency count, input isolation rule, oracle
comparison strategy, timeout/leak tolerance, race command, and what the test
cannot prove. Define how each test run cleans up its evaluator.

## Implementation checkpoint

Implement only the test harness and narrowly scoped fixes required by a
reproduced correctness/lifecycle defect. Keep tests deterministic and bounded.

## Test and evidence checkpoint

Add tests proving:

- one AST evaluated concurrently with distinct input environments returns the
  correct value for every request;
- no request receives another request's result or mutated state;
- sequential and channel evaluators agree across a fixed expression corpus;
- cancellation and shutdown remain correct under concurrent admission;
- repeated evaluator construction/close does not show unbounded goroutine growth
  within the documented tolerance; and
- `go test -race ./...` passes for the relevant module scope.

Run focused tests, full tests, vet, race tests, and the bounded stress command.
Record workload size, timeout, platform, and limitations. A successful race run
is evidence for tested schedules and inputs, not proof of all race freedom.

## Review checkpoint

Review fixture isolation, synchronization, flake resistance, race output, leak
measurement limits, oracle validity, and whether any production fix is justified
by a reproducible failure.

## Commit checkpoint

Stage only the selected tests/fixes and this handoff, validate cached scope, and
commit with the exact declared subject.

## Prerequisite and next gate

Step 09 must prove reliable cancellation and shutdown. Step 11 requires this
step's reuse, race, and cleanup evidence to be reviewed.

## Agent handoff notes

Reserved for workload design, test implementation, exact commands/results,
race/leak evidence, review, limitations, and next action.