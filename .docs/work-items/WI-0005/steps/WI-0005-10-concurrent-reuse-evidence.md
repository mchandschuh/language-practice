---
schema: language-practice-step/v1
id: WI-0005
step: "10"
status: implemented-unproven
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

### Design checkpoint

The test workload is bounded: 16 concurrent submissions for immutable AST reuse,
12 sequential submissions for request-local mutation isolation, a four-expression
sequential-oracle corpus, four concurrently admitted requests during shutdown,
and 64 evaluator construction/close cycles for the leak observation. The runtime
admission capacity is set to the workload size for reuse tests and four for the
shutdown test. Each request receives a distinct `Environment` input, and the
runtime's existing admission path must clone it before evaluation. Results are
matched by request index rather than completion order.

The reuse tests parse one AST before starting requests. The isolation expression
`x++ + ++x` verifies both returned values and request-local final state while
checking that caller-owned input maps remain unchanged. The fixed oracle corpus
covers precedence, unary arithmetic, prefix mutation, and postfix mutation; the
sequential `Evaluator` computes the expected value and final environment before
the channel evaluator is submitted. Shutdown uses the existing package barrier
to hold active work, waits for all three queued submissions to be admitted, then
closes the evaluator and verifies terminal `ErrEvaluatorClosed` results.

Every evaluator is explicitly closed. Result waits and admission waits use a
one-second timeout. Repeated-close observation compares `runtime.NumGoroutine`
against the pre-test count with a tolerance of two goroutines after 64
construction/close cycles. This is a bounded leak observation, not a universal
goroutine-leak proof. The tests do not claim performance improvement or parallel
worker execution; the current runtime owns one evaluator worker.

### Implementation summary

Added `concurrent_reuse_test.go` under the selected math package. Also updated
the existing `lifecycle_test.go` admitted-close regression test to use its
barrier deterministically; the test had a scheduling race where `Close` could
complete before work began. The tests cover concurrent reuse of one immutable
AST, isolation of mutable request environments,
agreement with an independent sequential evaluator corpus, concurrent admission
and cancel-first shutdown, and bounded repeated construction/close. No production
code was changed because no reproducible correctness or lifecycle defect was
found.

### Validation and evidence

- `cd go && gofmt -w basic/calc/math/concurrent_reuse_test.go`: passed.
- `cd go && go test ./basic/calc/math -count=1`: passed.
- `cd go && go test ./basic/calc/math -run 'Concurrent|Oracle|Repeated' -count=20`: passed.
- `cd go && go test -race ./basic/calc/math -count=1`: passed.
- `cd go && go vet ./basic/calc/math`: passed.
- `git diff --check`: passed.
- `cd go && go test ./... -count=1` remains blocked by the pre-existing,
  out-of-scope `basic/calc/main.go:24:24` call to `parsed.Evaluate`.
- `cd go && go vet ./...` reports the same pre-existing `parsed.Evaluate`
  compile failure.
- `cd go && go test -race ./... -count=1` reaches the math package successfully
  but fails the overall command on the same unrelated compile error.

Focused tests passed on the repository's Go 1.26.8 module environment. Evidence
proves the exercised bounded workloads, request isolation, oracle comparisons,
shutdown paths, package race schedules, and worker-count tolerance. It does not
prove arbitrary scheduling, universal race freedom, universal leak freedom,
performance, or repository-wide race/test/vet success while the unrelated CLI
compile error remains.

### Review checkpoint

- Fixture ASTs are parsed once and never mutated by the tests.
- Input maps are distinct per request and are checked after evaluation.
- Results are indexed rather than assumed to complete in submission order.
- Barriers and wait groups provide deterministic synchronization; waits have
  explicit bounds.
- The sequential oracle uses the direct evaluator rather than channel internals.
- Shutdown checks active and queued terminal results and waits for admission and
  close completion.
- Leak measurement uses an explicit tolerance and is reported only as an
  observation.
- Review found no production defect or test finding worth addressing within the
  selected scope. An unrelated untracked `.docs/work-items/WI-0004/` path was
  present before/alongside this work and is intentionally untouched.

### Next action

Advisor review, final scope inspection, and commit with the exact declared
subject. Step 11 remains gated on this step's committed reuse, race, and cleanup
evidence.