---
name: language-practice-observability-performance
description: Teach measurement, profiling, benchmarking, tracing, logging, and metrics without confusing diagnostic evidence with production readiness. Use for performance work, latency analysis, resource investigations, or operational visibility.
compatibility: Designed for local experiments in Go, Python, Rust, and TypeScript with the language's documented diagnostic tools.
metadata:
  version: "1"
  scope: "observability and performance"
---

# Measurement and diagnostics

Measure before optimizing. First define the question, workload, resource, and
acceptable evidence; then choose the least invasive diagnostic tool.

## Investigation order

1. Define the metric: latency, throughput, allocations, memory, CPU, I/O,
   contention, startup time, error rate, or queue depth.
2. Establish a representative, repeatable baseline.
3. Instrument the smallest useful boundary with structured, non-sensitive data.
4. Choose a tool appropriate to the question: benchmark, profiler, trace,
   runtime statistics, logs, metrics, or debugger.
5. Change one relevant variable and repeat the measurement.
6. Explain noise, warm-up, sampling, machine variance, and workload limits.

## Operational signals

Prefer stable event names, correlation identifiers that contain no secrets,
explicit units, bounded cardinality, and timestamps with a stated clock source.
Redact credentials, tokens, personal data, and untrusted raw content. Do not
log every value merely because it is available.

For concurrent or distributed exercises, make cancellation, queue depth,
retries, timeouts, dropped work, and shutdown visible where they matter.

## Learning boundary

Require the learner to predict the bottleneck and interpret the evidence before
showing an optimization in `handwrite` or `collaborate` mode. In `assist` mode,
generated instrumentation or harness code must still explain what it measures,
what it perturbs, and what it cannot establish.

Do not claim a local benchmark proves production capacity, tail latency,
scalability, or cost. Do not optimize a benchmark by changing the contract or
discarding meaningful work without documenting the change.

## Tool anchors

Use current official documentation for the language and runtime. Go's official
diagnostics guidance distinguishes profiling, tracing, debugging, and runtime
statistics. Python's standard library includes performance measurement,
logging, and profiling tools. Rust and TypeScript require explicit selection of
the compiler, runtime, benchmark, profiler, and async/runtime tooling in use.