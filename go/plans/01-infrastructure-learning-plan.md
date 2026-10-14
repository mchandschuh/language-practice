# Go Infrastructure Learning Plan

This plan adapts the accepted Python behavioral baseline to Go. It preserves
the infrastructure objective, experiment contract, evidence boundaries, and
20-stage progression while teaching Go's packages, values, interfaces,
explicit errors, goroutines, channels, contexts, and standard tooling. It does
not translate Python classes, exceptions, or generators into Go terminology.

The track develops Go writing fluency in service of dependable infrastructure:

- platform and backend foundations;
- distributed request lifecycles;
- concurrency, queues, backpressure, and cancellation;
- caching and consistency;
- timeouts, retries, rate limiting, isolation, and fallbacks;
- API and contract design;
- observability and debugging;
- throughput, latency distributions, and tail behavior;
- progressive delivery, verification, and rollback; and
- reusable libraries and operational maturity.

The goal is not to become a syntax typist or to reproduce another language's
design. The goal is to write enough Go fluently to investigate, model, measure,
debug, and improve dependable systems.

## Learner profile and teaching strategy

The learner can read most programming languages but currently writes Go slowly
and is not yet fluent. Use short, repeated writing loops before introducing
large infrastructure abstractions:

1. Read a small specification and the relevant official documentation.
2. Predict behavior and identify a measurable question.
3. Handwrite a small implementation.
4. Run a focused test or experiment.
5. Inspect output, diagnostics, failures, logs, or visualization.
6. Change one important variable and rerun the same experiment.
7. Explain what changed, why, and what remains unproven.
8. Write a concise GitHub-compatible experiment report.

Every exercise records prerequisites, learner ownership, permitted assistance,
observable validation, edge cases, failure modes, and a reflection checkpoint.

## Assistance modes

### Handwrite

The learner writes the core code. The agent may provide requirements, official
API references, small syntax hints, test targets, diagnostic interpretation,
and review feedback, but not a drop-in solution. Use this for functions,
slices/maps, interfaces, errors, `context`, channels, and basic concurrency.

### Collaborate

The learner remains the primary implementer. The agent may propose designs,
inspect a partial attempt, suggest hypotheses and measurements, and show small
targeted fragments. Use this for API design, failure analysis, concurrency
diagnosis, profiling interpretation, and tradeoffs.

### Assist

The agent may generate repetitive scaffolding after the contract is understood,
but must explain what it measures or does and preserve a learning checkpoint.
Use this for fixtures, adapters, benchmark harnesses, report formatting,
serialization, and repetitive CLI or load-generator plumbing.

## Version and dependency policy

The repository currently has no Go module, `go.mod`, `go.sum`, declared Go
version, or local Go toolchain. The Go specification currently reports Go 1.27;
that is a documentation observation, not this repository's execution contract.
Record the actual Go version, operating system, module state, dependencies, and
commands in every future experiment report. Do not create a module or add
dependencies merely to implement this plan.

Prefer the standard library while learning. Add an external module only for a
deliberate objective, after recording its version and reproducibility boundary.
Use the `go` command for module changes and review `go.mod` and `go.sum`.

`Effective Go` remains a useful idiom source, but its official notice says it
was written for the 2009 release and is not actively updated; it does not fully
cover generics, modules, or newer libraries. Use the specification and current
official package documentation for those subjects.

Go is strongly typed and garbage-collected, with explicit concurrency support.
Values such as arrays and structs contain their data, while slices, maps,
channels, functions, pointers, and interfaces can refer to shared underlying
state. Make nil behavior, copying, method sets, blocking, cancellation,
package initialization, and goroutine ownership observable rather than
assuming they match another language.

## Infrastructure concept map

| Infrastructure concept | Go practice |
| --- | --- |
| Platform foundations | packages, modules, structs, interfaces, safe zero values |
| Distributed systems | local listeners, request lifecycles, explicit errors |
| High concurrency | goroutines, channels, `sync`, `context` |
| Caching and consistency | synchronized maps, TTL, invalidation, freshness |
| Queueing/backpressure | bounded channels, worker ownership, rejection, shutdown |
| Reliability guardrails | deadlines, retries, rate limits, isolation, fallbacks |
| Performance primitives | benchmarks, allocations, profiles, traces, distributions |
| Reusable workflows | small packages, consumer-owned interfaces, safe defaults |
| Observability/debugging | structured events, errors, pprof, traces, runtime signals |
| Safe change management | flags, rollout gates, verification, rollback simulation |
| API/contract design | exported types, interfaces, JSON boundaries, validation |

## The experiment contract

Every experiment answers a question; runnable code or an attractive graph is
not sufficient.

```text
Question:
Hypothesis:
System under test:
Workload:
Variables:
Controls:
Measurements:
Expected behavior:
Failure/stop conditions:
Validation:
Limitations:
```

Start with a sequential or single-worker baseline when studying concurrency.
Change one important variable at a time. Use fixed seeds, input sizes,
controlled clocks, bounded durations, and enough repetitions to expose noise.
Preserve raw observations separately from interpretation. Stop when error rate,
queue growth, memory, or duration crosses the declared safety limit. Never load
test production or a third-party service without explicit authorization.

Each observation must be labeled as a language guarantee, standard-library
guarantee, implementation detail, or local workload result. A passing command
proves only the stated check; `go test`, `go vet`, and the race detector do not
prove understanding, universal race freedom, shutdown correctness, or capacity.

## Load-testing ladder

### Level 0 — Micro experiment

Use a function or small component. Measure correctness, elapsed time,
allocations where relevant, or state transitions with a focused test or small
Go executable.

### Level 1 — Component workload

Exercise a parser, cache, retry policy, queue, or worker pool with controlled
inputs. Record throughput, latency, failures, queue depth, and cleanup state.
Use `go test -bench` only after the behavior is understood.

### Level 2 — Local service simulation

Use local processes, `net` listeners, or a fake dependency. Vary arrival rate,
service time, dependency errors, deadlines, and worker count. Keep it bounded
and disposable.

### Level 3 — Optional tool-assisted load test

Only after component behavior is understood may an external load or
visualization tool be considered. Pin it, explain the reason, and preserve a
standard-library fallback where practical. A local load generator is not
production-capacity evidence.

## Visualization ladder

Start with sorted tables, ASCII timelines, queue-depth snapshots, request
lifecycle diagrams, Markdown summaries, and bounded CSV. When a graph adds
understanding, prefer reproducible SVG or self-contained HTML plus source data
and a regeneration command.

Visualize latency distributions and tails, throughput versus arrival rate,
queue depth, accepted/rejected/retried/timed-out/completed work, worker
utilization, cache outcomes, allocations, rollout exposure, and dependency
timelines. State units, workload, sample count, measured/simulated/inferred
status, aggregation, noise, and limitations. Do not smooth away failures or
discard outliers without explanation.

## GitHub experiment documentation

Future experiments may use the smallest required structure under:

```text
go/experiments/<section>/<experiment-name>/
├── README.md
├── experiment.go
├── experiment_test.go
├── results/
└── figures/
```

Do not create this tree during plan creation. Each substantial report should
contain:

```markdown
# Experiment title
## Question
## Hypothesis
## System and workload
## Method
## Commands
## Results
## Visualization
## Interpretation
## Surprising behavior
## Limitations
## Reproduction
## Follow-up experiment
```

Record the learner's expectation, observation, confusion, evidence that changed
the mental model, rejected approach, and next test. Keep raw data bounded and
free of credentials, private data, and unbounded logs. Record Go/tool versions,
platform, module/dependency versions, seed, duration, workload, and command.

## Curriculum

## GO-00 — Go writing loop

**Infrastructure connection:** A reliable writing loop is the prerequisite for
every later service, diagnostic, and operational tool.

**Concepts:** declarations, expressions, control flow, functions, zero values,
formatting, package compilation, and compiler diagnostics.

**Learner-owned work:** Handwrite 5–20 line transformations, predict output,
fix an edge case, and reimplement a variation from memory. Use `handwrite` mode.

**Experiment/measurement:** Compare fixed-input implementations for correctness
before timing them; record a Markdown table of inputs, outputs, and times.

**Reflection:** Explain what the compiler established and what the comparison did
not establish about production performance.

## GO-01 — Data structures and state

**Infrastructure connection:** Slices, maps, and structs model queues, caches,
counters, configuration, and lifecycle state.

**Concepts:** arrays versus slices, length/capacity, `append`, maps, structs,
zero values, nil slices/maps, copying, and shared backing arrays.

**Learner-owned work:** Handwrite a bounded FIFO, counters, a configuration
snapshot, and a small lifecycle model; explain who may mutate each value.

**Experiment/measurement:** Vary capacity and producer/consumer rates; observe
blocking, rejection, growth, aliasing, and invariants after each run.

**Reflection:** Distinguish documented value/reference behavior from local
observations; do not equate slices or maps with another language's collections.

## GO-02 — Functions as contracts

**Infrastructure connection:** Functions define request handlers, parsers,
retry decisions, validation boundaries, and worker operations.

**Concepts:** signatures, multiple returns, closures, callbacks, method values,
side effects, and explicit iterator-like functions or channels rather than
Python generator translations.

**Learner-owned work:** Handwrite parsers, retry-decision functions, callbacks,
and pipeline stages with clear preconditions and postconditions.

**Experiment/measurement:** Inject malformed input and dependency failures;
compare explicit `error` returns and fallback values without hiding the cause.

**Reflection:** Explain why a function, interface, or channel is the appropriate
boundary and what error paths remain untested.

## GO-03 — Packages, modules, CLIs, and reproducibility

**Infrastructure connection:** Package boundaries support reusable libraries,
shared tools, controlled initialization, and dependency reproducibility.

**Concepts:** packages, imports, exported identifiers, `internal`, `main`,
`go.mod`, module paths, the `go` command, configuration, and CLI parsing.

**Learner-owned work:** Split a small program into packages, separate CLI
parsing from pure logic, and document a module boundary in `collaborate` mode.

**Experiment/measurement:** Compare startup behavior with and without import-time
work and map package coupling; report module and environment details.

**Reflection:** Record that no module or dependency manifest exists in this
repository and that no reproducibility claim can be made yet.

## GO-04 — Errors, cleanup, and resource lifetime

**Infrastructure connection:** Explicit failure paths and cleanup contain file,
socket, dependency, and partial-initialization failures.

**Concepts:** `error`, wrapping and comparison, sentinel/type errors, `defer`,
`io.Closer`, partial initialization, and cancellation cleanup.

**Learner-owned work:** Handwrite a resource wrapper and cleanup after each
possible acquisition failure; return contextual errors without destroying causes.

**Experiment/measurement:** Inject failures at acquisition and release
boundaries; record open resources and whether cleanup is idempotent.

**Reflection:** Distinguish explicit error returns from exceptions and record
which cleanup guarantees are documented versus merely observed.

## GO-05 — Types, interfaces, and runtime contracts

**Infrastructure connection:** Static types and small interfaces stabilize APIs,
permit dependency substitution, and expose boundary mistakes early.

**Concepts:** defined types, structs, interfaces, method sets, pointer/value
receivers, type assertions and switches, generics with version caveats, and
runtime validation of decoded or untrusted data.

**Learner-owned work:** Handwrite request/response types, consumer-owned small
interfaces, fakes, and explicit boundary validators.

**Experiment/measurement:** Pass malformed runtime data through a statically
valid boundary and substitute two interface implementations; record compile-
time guarantees separately from runtime validation.

**Reflection:** Explain why an interface is not a class hierarchy and what the
compiler cannot validate about external data.

## GO-06 — Testing and regression control

**Infrastructure connection:** Tests provide repeatable verification for APIs,
invariants, failure paths, and operational guardrails.

**Concepts:** package tests, table-driven tests, subtests, temporary
directories, fakes, deterministic clocks, negative paths, examples, and
benchmarks.

**Learner-owned work:** Handwrite tests for parsers, error contracts, retry
limits, idempotency, queue shutdown, deadlines, and cancellation.

**Experiment/measurement:** Deliberately break each behavior and confirm the
test fails before restoring it; use `go test` and capture exact commands.

**Reflection:** State precisely what the tests prove and why passing tests do
not prove race freedom, learner understanding, or production readiness.

## GO-07 — Files, JSON, subprocesses, and safe tooling

**Infrastructure connection:** Safe operational tooling integrates services and
automation without hiding process, output, or cleanup failures.

**Concepts:** `io`, `os`, `encoding/json`, argument boundaries, exit status,
environment variables, `os/exec`, context-bound subprocesses, and output limits.

**Learner-owned work:** Handwrite a command wrapper with safe arguments,
non-zero-exit reporting, timeout cleanup, and malformed-output handling.

**Experiment/measurement:** Compare successful, slow, missing, malformed, and
failing child processes; record process state and captured output separately.

**Reflection:** Identify shell and external-tool assumptions that remain outside
Go's standard-library guarantees.

## GO-08 — Sockets and request lifecycles

**Infrastructure connection:** Transport boundaries expose framing, partial I/O,
timeouts, malformed messages, and dependency failure.

**Concepts:** `net`, listeners, `net.Conn`, deadlines, partial reads/writes,
framing, connection closure, and malformed messages.

**Learner-owned work:** Handwrite a local client/server protocol using the
standard library and keep transport, parsing, and domain logic separate.

**Experiment/measurement:** Inject partial frames, slow peers, malformed input,
disconnects, and deadline expiry; measure accepted, rejected, and abandoned
requests.

**Reflection:** Separate protocol guarantees from behavior observed on one local
network stack and workload.

## GO-09 — Asynchronous execution fundamentals

**Infrastructure connection:** Goroutines and cancellation support concurrent
request lifecycles without losing ownership or shutdown responsibility.

**Concepts:** goroutines, channel directions, blocking sends/receives,
`select`, `context`, cancellation, ownership, leaks, and happens-before
reasoning.

**Learner-owned work:** Handwrite a cancellable worker and request lifecycle;
the learner owns when work starts, stops, and reports errors.

**Experiment/measurement:** Compare a sequential baseline with a goroutine and
channel version; measure cancellation latency and check for leaked goroutines.

**Reflection:** Explain why goroutine creation is not unlimited parallelism or
proof of correct shutdown.

## GO-10 — Threads, goroutines, and parallel work

**Infrastructure connection:** Controlled parallel work affects throughput,
contention, isolation, and resource usage.

**Concepts:** goroutines versus OS threads, `sync.WaitGroup`, mutexes, atomics,
scheduling, shared memory, worker sizing, and the race detector.

**Learner-owned work:** Handwrite synchronized shared-state and message-passing
variants; choose the simpler ownership model and justify it.

**Experiment/measurement:** Compare sequential, goroutine, and bounded-worker
variants for correctness, throughput, allocations, and shutdown; use
`go test -race` for race-relevant tests.

**Reflection:** A race-detector run covers observed schedules and inputs only;
it does not prove universal race freedom.

## GO-11 — Queues, worker pools, overload, and shutdown

**Infrastructure connection:** Bounded queues and explicit worker ownership make
backpressure, overload rejection, and graceful shutdown observable.

**Concepts:** buffered/unbuffered channels, bounded admission, rejection,
worker ownership, close semantics, draining, cancellation, and shutdown order.

**Learner-owned work:** Handwrite a bounded worker pool with explicit stop and
drain behavior; test blocked sends, closed-channel mistakes, and worker leaks.

**Experiment/measurement:** Vary arrival rate, service time, and worker count;
measure depth, rejection, completion, cancellation, and shutdown time.

**Reflection:** State which queue properties are intentional contracts and which
are implementation observations.

## GO-12 — Timeouts, retries, rate limits, isolation, and fallbacks

**Infrastructure connection:** Reliability guardrails bound dependency failure,
duplicate work, overload, and degraded responses.

**Concepts:** context deadlines, timers/tickers, retry budgets, exponential
backoff, jitter, idempotency, token/leaky-bucket models, bulkheads, and errors.

**Learner-owned work:** Handwrite a bounded retry policy and rate limiter around
a fake dependency, retaining cancellation and error causes.

**Experiment/measurement:** Inject timeouts, transient/permanent errors, slow
responses, and duplicate execution; record retries, delays, load, and fallback.

**Reflection:** Retry success does not prove downstream idempotency or service
reliability; explain the declared failure budget.

## GO-13 — Caching, TTL, invalidation, and consistency

**Infrastructure connection:** Caches trade freshness and origin load for
latency, so consistency and invalidation must be explicit.

**Concepts:** synchronized maps, monotonic time, TTL, invalidation, stale reads,
refresh ownership, hit/miss/stale/eviction states, and snapshots.

**Learner-owned work:** Handwrite a cache with a stated freshness contract and
safe concurrent access; do not hide synchronization behind an unexplained tool.

**Experiment/measurement:** Vary TTL, invalidation timing, refresh latency, and
concurrent requests; measure freshness, origin calls, latency, and stale data.

**Reflection:** Describe what happens when cache and source disagree and which
consistency result is only a local observation.

## GO-14 — API and contract design

**Infrastructure connection:** Stable Go APIs and wire contracts support safe
evolution, compatibility, pagination, and idempotent requests.

**Concepts:** exported API design, interfaces, structs, JSON tags, validation,
error contracts, versioning, adapters, pagination, and idempotency keys.

**Learner-owned work:** Handwrite a local service contract and two implementations
behind a consumer-owned interface, with contract tests.

**Experiment/measurement:** Replay old clients against changed behavior and
classify accepted, rejected, ambiguous, and degraded requests.

**Reflection:** Distinguish compile-time Go API compatibility from wire-format
compatibility and runtime validation.

## GO-15 — Reusable library design and misuse resistance

**Infrastructure connection:** Small reusable packages turn operational lessons
into safe defaults without forcing consumers into a framework.

**Concepts:** package boundaries, small interfaces, constructors, zero-value
usability, safe defaults, options, documentation, examples, and compatibility.

**Learner-owned work:** Handwrite a reusable request, queue, or cache primitive
with two fake consumers; keep consumer code simpler than the primitive.

**Experiment/measurement:** Compare duplicated consumers with the reusable
workflow for failure behavior, operational signals, and change cost, not just
line count.

**Reflection:** Identify defaults, escape hatches, and misuse cases; reject a
class hierarchy or oversized interface when composition is clearer.

## GO-16 — Observability and debugging systems

**Infrastructure connection:** Instrumentation makes request failures,
contention, queue growth, and shutdown defects diagnosable.

**Concepts:** structured events, request IDs, wrapped errors, metrics, goroutine
dumps, `runtime/pprof`, `net/http/pprof`, execution tracing, `go vet`, and
compiler/test diagnostics.

**Learner-owned work:** Handwrite instrumentation for retries, timeouts,
fallbacks, queue transitions, and shutdown; preserve the learner's hypothesis.

**Experiment/measurement:** Inject latency or dependency failure and diagnose
from logs, metrics, profiles, or traces before reading implementation details.

**Reflection:** State whether each result is documented, locally observed, or
inferred; diagnostics can perturb one another and must be scoped.

## GO-17 — Performance, profiling, throughput, and tail latency

**Infrastructure connection:** Workload design and profiles reveal bottlenecks,
queueing delay, tail behavior, and cost-relevant tradeoffs.

**Concepts:** `testing.B`, `go test -bench`, allocations, `b.ReportAllocs`,
warmup, p50/p95/p99-like summaries, service versus queueing time, `pprof`, and
`go tool trace`.

**Learner-owned work:** Handwrite a fixed-workload benchmark, preserve raw
samples, separate service from queueing time, and explain one profile first.

**Experiment/measurement:** Compare one changed variable at a time across
sequential and bounded-worker variants; report noise, outliers, and machine
variance.

**Reflection:** Profilers and traces perturb execution; local measurements do
not establish production capacity or a universal optimization.

## GO-18 — Progressive delivery and rollback simulation

**Infrastructure connection:** Rollout gates and rollback contain regressions
while preserving compatible state and observable verification.

**Concepts:** feature flags, cohorts, compatibility, thresholds, halt
conditions, configuration snapshots, rollback state, and recovery time.

**Learner-owned work:** Handwrite a local rollout state machine with staged
exposure, verification, automatic halt, and state-safe rollback.

**Experiment/measurement:** Expose simulated traffic at 1%, 10%, 50%, and 100%;
inject error and tail-latency regressions and record gate and rollback timing.

**Reflection:** Document which signals the gate could see and what the local
simulation cannot prove about deployment systems.

## GO-19 — Safe request-platform capstone

**Infrastructure connection:** The capstone integrates request contracts,
dependency seams, reliability guardrails, observability, performance evidence,
progressive delivery, rollback, and operational ownership.

**Concepts:** packages, interfaces, explicit errors, `context`, bounded queues,
worker ownership, rate limits, retries, fallbacks, metrics, logs, tests,
profiling, rollout, and rollback.

**Learner-owned work:** Assemble the platform in slices. The learner owns the
request contract, failure model, queue behavior, guardrails, measurements, and
final operational explanation; assistance may supply only understood plumbing.

**Experiment/measurement:** Run a bounded local workload, inject dependency
failures and overload, and compare accepted, rejected, retried, timed-out,
completed, and rolled-back work.

**Reflection:** Produce an architecture and GitHub report with raw-data
references, visualizations, failures, decisions, validation, and explicit
limitations. The plan itself creates no implementation or experiment artifact.

## Official documentation anchors

- Go specification: `https://go.dev/ref/spec`
- Effective Go and its historical coverage notice: `https://go.dev/doc/effective_go`
- Go modules and dependency management: `https://go.dev/doc/modules/managing-dependencies`
- Go diagnostics: `https://go.dev/doc/diagnostics`
- Go command reference: `https://go.dev/cmd/go/`
- Standard-library documentation: `https://pkg.go.dev/std`
- `context`: `https://pkg.go.dev/context`
- `sync`: `https://pkg.go.dev/sync`
- `testing`: `https://pkg.go.dev/testing`
- `net`: `https://pkg.go.dev/net`
- `encoding/json`: `https://pkg.go.dev/encoding/json`
- `runtime/pprof`: `https://pkg.go.dev/runtime/pprof`
- `net/http/pprof`: `https://pkg.go.dev/net/http/pprof`

The specification and package documentation are the authority for current
semantics. Effective Go supplies historical idioms, not complete current
module, generic, or library guidance.

## Evidence boundaries

This document establishes intended scope, ordering, source references, and
learning evidence. It does not prove Go compilation, runtime behavior, race
freedom, benchmark validity, learner understanding, production capacity,
security readiness, or production readiness. Future local experiments establish
only behavior for their stated implementation, workload, environment, and
toolchain. Keep documented guarantees, local observations, design judgments,
and open questions visibly distinct.
