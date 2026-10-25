---
schema: language-practice-step/v1
id: WI-0005
step: "12"
status: implemented-unproven
skill: language-practice-observability-performance
depends_on: ["11"]
target: go/basic/calc/
---

# WI-0005 step 12 — Record evaluator evidence

## Objective

Measure and document the sequential, type-dispatched, and channel-driven
evaluators under a fixed bounded workload. Explain correctness evidence,
throughput/latency observations, allocation cost, queueing effects, and the
limitations of local measurements.

## Commit subject

```text
WI-0005: 12: doc(calc): record evaluator evidence
```

## Scope and ownership

- Own only the calculator benchmark/report files under `go/basic/calc/` and this
  step's handoff update.
- Add fixed-input benchmarks or a bounded experiment plus a Markdown report.
- Compare one variable at a time: evaluator implementation, worker count, or
  queue capacity.
- Preserve a sequential correctness oracle and do not optimize by weakening
  error, cancellation, ordering, or shutdown guarantees.
- Do not claim production capacity, portability, or universal performance.

## Design checkpoint

Record workload expressions, input distribution, iteration count/duration,
warmup/benchmark command, platform/toolchain, queue/worker configuration,
metrics, aggregation, and known perturbations. Separate service computation
time from queueing/admission time where the instrumentation permits.

## Implementation checkpoint

Add only bounded benchmark/report artifacts and instrumentation needed to answer
the stated measurement question. Keep raw data small, reproducible, and free of
credentials or unbounded logs.

## Test and evidence checkpoint

Run and record:

- correctness tests against the sequential oracle;
- `go test -bench` with fixed benchmark inputs;
- allocation evidence where relevant;
- a bounded sequential-versus-channel comparison;
- the final `go test ./...`, `go vet ./...`, and relevant race command; and
- exact commands, versions, platform, workload, and limitations.

The evidence gate is a reproducible local measurement report, not a claim that
the channel evaluator is faster or production-ready.

## Review checkpoint

Review benchmark fairness, workload disclosure, units, sample/iteration count,
outliers, queueing interpretation, raw-data bounds, and whether conclusions are
limited to the measured host and inputs.

## Commit checkpoint

Stage only the selected benchmark/report files and this handoff, validate cached
scope and whitespace, and commit with the exact declared subject. Append the
final advisor review only after the post-commit validation is complete.

## Prerequisite and next gate

Step 11 must prove the CLI and final test helper. This is the final implementation
step; WI-0005 completion still requires review of all step commits, handoffs,
advisor findings, and the final evidence limits.

## Agent handoff notes

Reserved for measurement design, implementation, commands/results, report scope,
review, evidence tier, limitations, and finalization action.

### Design checkpoint

The fixed corpus is four parsed expressions: `x * x + 2` with `x=3`,
`(9 - 4) / 5 + 2 * 3` with an empty environment, `-x + 3` with `x=4`, and
`x++ + ++x` with `x=5`. Parsing and expected-value calculation occur before the
timed region. The sequential `Evaluator{}.EvaluateWithEnvironment` path is the
oracle. Channel requests use the same immutable ASTs, cloned inputs, and one
terminal result per request. Queue capacities 1 and 16 are compared; the
runtime has one worker. A single-request benchmark measures submission through
result receipt, while a fixed batch of 16 exposes bounded admission effects.
`ReportAllocs` records allocation cost. Internal worker service time cannot be
isolated by the current public API, so queueing observations are explicitly
caller-boundary observations.

### Implementation summary

Added `go/basic/calc/math/evaluator_benchmark_test.go` with fixed-input
sequential, single-request channel, and batch channel benchmarks. Added
`go/basic/calc/EVALUATOR-EVIDENCE.md` documenting workload, method, environment,
bounded benchmark output, correctness relationship, queue interpretation, and
limitations. No production evaluator or CLI code changed.

### Validation and evidence

- `cd go && gofmt -w basic/calc/math/evaluator_benchmark_test.go`: passed.
- `cd go && go test ./basic/calc/math -count=1`: passed.
- `cd go && go test -bench 'Benchmark(SequentialOracle|ChannelEvaluator|ChannelBatch)' -benchmem -benchtime=100ms -count=5 ./basic/calc/math`: passed; five bounded samples recorded in `go/basic/calc/EVALUATOR-EVIDENCE.md`.
- `cd go && go test ./... -count=1`: passed.
- `cd go && go vet ./...`: passed.
- `cd go && go test -race ./basic/calc/math -count=1`: passed.
- `git diff --check`: passed.

Evidence tier is static review, deterministic unit/oracle tests, bounded local
benchmark measurement, package race testing, repository-wide tests, and vet.
It proves only the tested expressions, inputs, schedules, host, toolchain, and
single-worker runtime configuration. It does not prove production capacity,
tail latency, portability, scaling, universal race/leak freedom, or that the
channel evaluator is faster.

### Review checkpoint

Reviewed benchmark fairness, parse/setup boundaries, oracle validation, queue
capacity interpretation, stable allocation counts, noisy timing ranges, bounded
raw-data disclosure, and evidence limitations. The batch comparison is not
presented as service-time measurement. No findings worth addressing remain in
the selected benchmark/report scope.

### Next action

Inspect the final scoped diff and commit the benchmark, report, and this handoff
with the exact subject `WI-0005: 12: doc(calc): record evaluator evidence`.