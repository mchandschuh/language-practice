# Rust Infrastructure Learning Plan

This plan adapts the shared infrastructure-learning objective to Rust. It
preserves the behavioral progression of the accepted Python plan while treating
ownership, borrowing, lifetimes, explicit error values, compiler diagnostics,
and Rust's concurrency model as Rust concepts rather than Python equivalents.

The learner can read most programming languages but writes Rust slowly. The
goal is writing fluency in service of dependable infrastructure: request
lifecycle design, queues and backpressure, reliability guardrails, observability,
performance, progressive delivery, and operational review.

## Learner profile and teaching strategy

Each exercise follows this loop: read a small specification and official
documentation; predict behavior and compiler diagnostics; handwrite a small
piece; run a focused check or experiment; inspect failures and measurements;
change one variable; and explain what changed and what remains unproven.
Compiler diagnostics are learning artifacts. Ownership, borrowing, lifetimes,
`Drop`, and resource ownership come before `clone`, reference counting, or
interior mutability.

## Assistance modes

### Handwrite

The learner writes core functions, state models, ownership examples, enums,
traits, error paths, iterators, tests, threads, channels, and synchronization.
The agent may provide requirements, official references, tiny syntax hints,
test targets, compiler-error interpretation, and review feedback, but not a
drop-in solution.

### Collaborate

The learner remains the primary implementer while the agent helps design APIs,
inspect partial attempts, isolate failures, interpret measurements, and compare
tradeoffs. This mode covers dependency boundaries, retry policies, metrics,
profiling, and concurrency diagnosis.

### Assist

The agent may provide repetitive fixtures, adapters, report formatting,
benchmark harnesses, and operational plumbing after the behavior is understood.
The learner must still explain what the generated material measures, its
failure modes, and its limits.

Every exercise records prerequisites, learner-owned work, permitted assistance,
observable validation, edge cases, failure modes, and a reflection checkpoint.

## Version and dependency policy

Use the repository's declared Rust toolchain and edition when they exist. This
repository currently declares neither a Rust toolchain nor a Cargo edition, and
no local Rust toolchain is an execution contract for this plan. The current
Rust Book states that its text assumes Rust 1.97.0 or later and Edition 2024;
that is a documentation observation, not a repository guarantee. Use Cargo's
standard commands when a future exercise supplies a manifest. Do not create a
manifest, dependencies, experiments, or generated artifacts as part of this
plan. Any async runtime, serialization crate, tracing crate, profiler, or other
external facility must be selected and versioned explicitly by that future
exercise.

## Infrastructure concept map

| Infrastructure concern | Rust practice |
| --- | --- |
| State and ownership | structs, enums, moves, borrows, lifetimes, `Drop`, and explicit ownership |
| API contracts | public types, newtypes, traits, generics, exhaustive patterns, `Result`/`Option` |
| Resource lifecycle | scope-based cleanup, ownership transfer, `Drop`, explicit shutdown |
| Concurrency | threads, scoped threads, channels, locks, `Send`, `Sync`, and join handles |
| Reliability | deadlines, typed errors, bounded retries, idempotency, backpressure, fallbacks |
| Tooling | Cargo, rustfmt, compiler diagnostics, tests, clippy, and documented external tools |
| Operations | bounded logs/metrics, queue depth, cancellation, latency distributions, rollout gates |

## The experiment contract

Every experiment answers a question with a stated hypothesis, fixed workload,
controlled inputs, observable outputs, and a reproduction command once code is
authorized. Start with a sequential baseline, then vary one relevant variable.
Record Rust/toolchain and dependency versions, platform, seed, duration,
workload, commands, failures, and limitations. Keep raw measurements separate
from summaries. A successful compile or test proves only that stated check; it
does not prove understanding, deadlock freedom, production capacity, or
performance.

## Load-testing ladder

1. **Micro experiment:** one deterministic input, one question, and compiler or
   test diagnostics.
2. **Sequential baseline:** fixed fake dependencies and bounded input.
3. **Controlled concurrency:** vary thread/task count, queue capacity, and
   cancellation timing while preserving the contract.
4. **Failure and overload:** inject slow, malformed, rejected, and failed work;
   measure queue depth, rejection, retries, shutdown, and recovery.
5. **Bounded local load:** sweep arrival rate and service capacity with an
   explicit stop condition. Never present this as production-capacity evidence.

## Visualization ladder

Begin with tables and ASCII timelines for ownership, request lifecycle,
resource cleanup, queue depth, and state transitions. Then preserve raw JSON or
CSV measurements before making a Markdown, SVG, or other figure. Show latency
quantiles, throughput, error classes, retries, cancellation, and worker count
without smoothing away failures or unexplained outliers.

## GitHub experiment documentation

Future experiments should use the smallest required Rust-specific directory,
never empty scaffolding for the whole curriculum. A report should contain:

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

Reports must preserve what was expected, what happened, what was confusing,
what changed the mental model, and what remains to test. Do not commit Cargo
caches, credentials, private data, unbounded logs, or generated artifacts that
are not part of an explicitly authorized experiment.

## Curriculum

## RUST-00 — Rust writing loop and diagnostics

**Infrastructure connection:** establishes the edit–check–test–measure loop for
all later systems work.

**Concepts:** expressions, bindings, mutability, functions, types, ownership
diagnostics, `rustfmt`, and `cargo check`.

**Learner work:** handwrite small transformations, predict move and mutability
errors, run focused checks, and record the minimal correction and its limit.

## RUST-01 — Data structures, ownership, and invariants

**Infrastructure connection:** models queues, caches, counters, configuration,
and lifecycle state.

**Concepts:** `Vec`, `String`, slices, `HashMap`, structs, moves, shared and
mutable borrows, aliasing, and the explicit cost of `clone`.

**Learner work:** handwrite a stateful collection and snapshot/update operations;
experiment with borrow-checker failures and document enforced invariants.

## RUST-02 — Functions, closures, iterators, and side effects

**Infrastructure connection:** request handlers, retry policies, cache actions,
and worker callbacks.

**Concepts:** ownership in signatures, borrowed versus owned values, closures,
iterator adapters, lazy work, generic bounds, and side effects.

**Learner work:** handwrite pure and iterator-based transformations, predict
allocation and ownership behavior, then validate with focused tests or inspection.

## RUST-03 — Modules, crates, Cargo, and reproducibility

**Infrastructure connection:** reusable libraries, package boundaries, dependency
control, and reproducible builds.

**Concepts:** modules, visibility, crates, packages, manifests, lockfiles,
features, profiles, configuration, and Cargo commands.

**Learner work:** in a future authorized Cargo exercise, split a library and
binary, inspect the dependency graph, and report toolchain/edition limits.

## RUST-04 — Errors, cleanup, and resource lifetime

**Infrastructure connection:** cleanup of files, sockets, processes, and partial
initialization during failure.

**Concepts:** `Result`, `Option`, `?`, typed errors, error context, `Drop`, scope,
ownership transfer, and partial initialization.

**Learner work:** handwrite a resource wrapper and failure-injection paths;
measure cleanup on success and failure and distinguish `Drop` guarantees from
external-resource behavior.

## RUST-05 — Types, enums, traits, and runtime contracts

**Infrastructure connection:** stable APIs, invalid-state prevention, and
dependency substitution.

**Concepts:** enums, exhaustive pattern matching, newtypes, traits, generics,
trait objects, `Send`/`Sync` bounds, and runtime validation of external data.

**Learner work:** model request states and dependency interfaces, test malformed
external input, and record compile-time versus runtime guarantees.

## RUST-06 — Testing and regression control

**Infrastructure connection:** automated verification and regression containment.

**Concepts:** unit/integration tests, deterministic fixtures, negative paths,
cleanup, `cargo test`, `cargo fmt --check`, and `cargo clippy`.

**Learner work:** handwrite tests for parsing, state preservation, retries,
idempotency, queue shutdown, timeout, and cancellation; break each behavior and
confirm the test fails.

## RUST-07 — Files, structured data, subprocesses, and safe tooling

**Infrastructure connection:** safe operational automation and service tooling.

**Concepts:** `std::fs`, buffered I/O, `Command`, exit status, child ownership,
timeouts, cleanup, malformed data, and external serialization boundaries.

**Learner work:** handwrite a safe argument-based command wrapper and compare
successful, slow, missing, malformed, and failing children in a lifecycle report.

## RUST-08 — Sockets and request lifecycles

**Infrastructure connection:** service interfaces, framing, partial I/O, and
dependency failure.

**Concepts:** `TcpListener`, `TcpStream`, blocking, partial reads/writes,
framing, deadlines, connection ownership, and malformed messages.

**Learner work:** implement a tiny length-prefixed local protocol; load-test
bounded connections and measure completions, timeouts, errors, and latency.

## RUST-09 — Asynchronous execution fundamentals

**Infrastructure connection:** high concurrency and unpredictable I/O patterns.

**Concepts:** `async fn`, futures, `await`, polling, executor/runtime boundary,
drop-based cancellation, pinning where relevant, `Send`/`Sync`, and blocking work.

**Learner work:** compare a synchronous baseline with an explicitly selected
external async runtime; measure cancellation and shutdown and label runtime
guarantees separately from language guarantees.

## RUST-10 — Threads and parallel work

**Infrastructure connection:** throughput, worker systems, and parallel failures.

**Concepts:** `std::thread`, scoped threads, `move` closures, `Send`, `Sync`,
`JoinHandle`, synchronization, and panic propagation.

**Learner work:** compare sequential and threaded deterministic workloads,
predict a non-`Send` compiler error, and sweep worker count including join cost.

## RUST-11 — Channels, worker pools, backpressure, and shutdown

**Infrastructure connection:** queueing, overload, load shedding, and graceful
shutdown.

**Concepts:** channel ownership, bounded versus unbounded queues, disconnect
semantics, worker lifecycle, `JoinHandle`, shutdown ordering, and backpressure.

**Learner work:** handwrite a worker pool with explicit termination and per-task
errors; vary arrival rate and capacity and measure rejection, timeout, and
recovery.

## RUST-12 — Timeouts, retries, rate limits, isolation, and fallbacks

**Infrastructure connection:** reliability guardrails and controlled degradation.

**Concepts:** `Duration`, `Instant`, deadlines, typed transient/permanent errors,
bounded backoff, jitter, idempotency, isolation, and fallback results.

**Learner work:** implement a deterministic fake dependency and `Result`-based
retry policy; measure attempts, delay, deadline, duplicate work, and fallback.

## RUST-13 — Caching, TTL, invalidation, and consistency

**Infrastructure connection:** latency reduction, stale data, and source/cache
consistency.

**Concepts:** ownership of cached values, `HashMap`, `Instant`, TTL, invalidation,
stale reads, synchronization boundaries, and `Option` misses.

**Learner work:** handwrite a bounded cache with expiration and invalidation;
experiment with hit, miss, stale, refresh, and source/cache disagreement.

## RUST-14 — API and contract design

**Infrastructure connection:** stable interfaces, compatibility, validation,
pagination, and idempotency.

**Concepts:** public types, newtypes, traits, `Result`/`Option`, ownership in
signatures, compatibility, and runtime protocol validation.

**Learner work:** design and test a request/response API with invalid input,
pagination boundaries, duplicate requests, and classified errors; publish a
contract matrix.

## RUST-15 — Reusable library design and misuse resistance

**Infrastructure connection:** shared frameworks, safe defaults, and independent
evolution.

**Concepts:** visibility, public API review, generic versus dynamic dispatch,
builders/state types, documentation tests, and explicit `unsafe` boundaries.

**Learner work:** refactor a small exercise into a narrow library API and review
ownership, error, shutdown, and abstraction contracts for misuse resistance.

## RUST-16 — Observability and debugging systems

**Infrastructure connection:** diagnosis of failures, saturation, latency, and
shutdown defects.

**Concepts:** `Debug`, error context, bounded structured logs, correlation IDs,
metrics, queue depth, compiler diagnostics, and `tracing` as an external choice.

**Learner work:** instrument the smallest useful boundary; measure lifecycle,
retries, cancellation, queue depth, and cleanup without sensitive data.

## RUST-17 — Performance, profiling, throughput, and tail latency

**Infrastructure connection:** scalability, cost, bottlenecks, and tail behavior.

**Concepts:** allocation and cloning, contention, warm-up, sampling, quantiles,
`cargo bench` limits, and benchmark/profiler tools as explicit boundaries.

**Learner work:** establish a sequential baseline, vary one variable, and report
throughput and latency distributions with machine, workload, and variance limits.

## RUST-18 — Progressive delivery and rollback simulation

**Infrastructure connection:** rollout gates, verification, and safe rollback.

**Concepts:** typed rollout state machines, exhaustive transitions, thresholds,
configuration boundaries, rollback, and recovery.

**Learner work:** handwrite a deterministic canary simulator, inject failed
verification, measure gate decisions, and visualize transitions and rollback.

## RUST-19 — Safe request-platform capstone

**Infrastructure connection:** integrates request handling, ownership, errors,
bounded concurrency, backpressure, observability, performance, rollout, and
rollback.

**Concepts:** typed state, traits for fakes, `Result`/`Option`, threads/channels,
shutdown, cancellation, and an explicitly selected async runtime if needed.

**Learner work:** build in slices from synchronous contract through fake
dependency, timeout/error model, retries, bounded queue, workers, rate limit,
fallback, instrumentation, measurement, rollout, and rollback; publish the
complete GitHub-readable report and operational review.

## Official documentation anchors

- Rust Book: `https://doc.rust-lang.org/book/`
- Rust Reference: `https://doc.rust-lang.org/reference/`
- Cargo Book: `https://doc.rust-lang.org/cargo/`
- Rust API Guidelines: `https://rust-lang.github.io/api-guidelines/`
- Rust standard library: `https://doc.rust-lang.org/std/`
- Rust threads: `https://doc.rust-lang.org/std/thread/`
- Rust channels: `https://doc.rust-lang.org/std/sync/mpsc/`
- Rust synchronization: `https://doc.rust-lang.org/std/sync/`
- Rust `Result`: `https://doc.rust-lang.org/std/result/`
- Rust `Option`: `https://doc.rust-lang.org/std/option/`

The Book currently documents Rust 1.97 and Edition 2024. The Reference follows
the latest stable release, and standard-library/Cargo pages are versioned with
their releases; with no repository toolchain or edition, these are anchors for
research, not an execution contract.

## Evidence boundaries

This document establishes scope, ordering, source references, and intended
learning evidence. Future local experiments can establish behavior only for
their stated workload, implementation, toolchain, and environment. They do not
establish production capacity, cloud behavior, provider reliability, security
readiness, universal race/deadlock freedom, cancellation correctness, or learner
understanding. Memory safety and successful compilation are not proofs of
logical correctness, performance, or operational readiness.