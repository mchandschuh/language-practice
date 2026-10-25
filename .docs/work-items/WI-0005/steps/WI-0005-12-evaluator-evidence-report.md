---
schema: language-practice-step/v1
id: WI-0005
step: "12"
status: planned
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