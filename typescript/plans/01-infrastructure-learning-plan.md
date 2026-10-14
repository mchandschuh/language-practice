# TypeScript Infrastructure Learning Plan

This plan adapts the shared infrastructure objective and behavioral sequence from
the accepted Python plan without translating Python APIs. It develops TypeScript
writing fluency for dependable platform, backend, distributed, and operational
systems while keeping TypeScript's static model separate from JavaScript and
Node's runtime behavior.

The learner should become able to model, implement, measure, debug, and improve
systems involving request lifecycles, queues, backpressure, cancellation,
caching, consistency, timeouts, retries, rate limits, isolation, fallbacks,
observability, tail latency, rollout, rollback, and reusable libraries.

## Learner profile and teaching strategy

The learner reads most programming languages but writes TypeScript slowly. Each
exercise therefore uses a small repeated loop: read the specification and
official documentation, predict behavior, handwrite a focused change, run a
bounded check or experiment, inspect output and diagnostics, change one variable,
rerun, and explain both the result and what remains unproven. The goal is
retrieval of common TypeScript shapes and infrastructure reasoning, not typing
speed or framework familiarity.

## Assistance modes

### Handwrite

The learner writes core types, functions, state transitions, narrowing,
generics, promise coordination, and failure paths. Assistance may provide
requirements, questions, official references, tiny syntax hints, test targets,
compiler-error interpretation, and review, but not a drop-in solution.

### Collaborate

The learner remains the primary implementer while the agent helps design API
boundaries, isolate a runtime/module failure, choose measurements, interpret
profiles, and compare tradeoffs. Small targeted fragments are allowed; a full
replacement is not the default.

### Assist

After the behavior is understood, the agent may generate repetitive fixtures,
load-driver plumbing, serialization adapters, report formatting, and plotting
boilerplate. The report must state what the generated code measures, what it
perturbs, and what it cannot prove.

Every exercise records prerequisites, learner-owned work, permitted assistance,
observable behavior, validation, edge cases, failure modes, and a reflection
checkpoint. A successful type check or test proves only its stated check.

## Version and dependency policy

No repository TypeScript version, package manifest, lockfile, `tsconfig.json`,
package script, test runner, bundler, or dependency contract is available. No
local `tsc` executable is available. Node 24.21.0 is a local observation only,
not a repository contract. Future exercises must name their selected compiler,
Node version, package manager, scripts, module mode, and external dependencies;
none are selected or installed by this plan. TypeScript documentation is an
authority for language and compiler concepts, while the selected Node runtime is
the authority for host behavior.

## Infrastructure concept map

| TypeScript concern | Infrastructure question | Evidence boundary |
| --- | --- | --- |
| Types, narrowing, unions, generics | Which states and contracts are valid? | Static model; external data still needs runtime validation. |
| Modules and declarations | What is the public package boundary? | Type declarations describe an interface; loading is a runtime concern. |
| Promises and event loop | What is scheduled, ordered, delayed, or rejected? | Fresh-process runtime traces are required; types do not prove scheduling. |
| Abort signals and workers | Who owns cancellation and work lifecycle? | Application protocol and Node host behavior must be tested. |
| `tsconfig` and scripts | Which source, emit, module, and test policies apply? | Configuration inspection is not proof of emitted or deployed behavior. |
| Metrics and load tests | What happens under queueing and overload? | Local measurements do not prove production capacity or tail behavior. |

## The experiment contract

For each experiment, state the hypothesis, prerequisites, fixed input, workload,
resource, observable behavior, metric and unit, expected invariant, failure
path, cleanup, and evidence tier. Predict before running. Establish a sequential
or single-request baseline, then vary one relevant factor such as concurrency,
queue capacity, payload size, cancellation timing, retry policy, or module mode.
Keep tests deterministic, bounded, offline, and free of credentials or private
data. Separate static type checking, emitted JavaScript inspection, fresh-process
module loading, runtime validation, and measured behavior. Label inferred versus
observed transitions and record warm-up, sampling, machine variance, and
workload limits. No experiment, manifest, dependency, or generated artifact is
created by this plan.

## Load-testing ladder

1. Fixed single request with a known valid and invalid payload.
2. Repeated sequential requests with success, rejection, and cleanup counts.
3. Small bounded concurrency with latency and error distributions.
4. Increasing queue depth and producer/consumer imbalance.
5. Capacity, overload, backpressure, timeout, retry, and rate-limit thresholds.
6. Cancellation at different lifecycle points and orderly worker shutdown.
7. A bounded Node-worker comparison, recording serialization and startup costs.

Never infer production capacity, scalability, cost, or tail guarantees from a
local run. Preserve meaningful work and report machine and workload limits.

## Visualization ladder

Begin with Markdown tables for inputs, states, outcomes, durations, and errors.
Then use ASCII timelines for event-loop turns, queue depth, retries, and
cancellation. Add Mermaid state or sequence diagrams for request lifecycles and
module boundaries. Only after the question is stable may a future exercise use a
small, explicitly selected plotting or tracing tool. Keep raw data bounded and
reviewable; do not commit generated figures or logs here.

## GitHub experiment documentation

Each future report states the question, prediction, environment and selected
versions, command, input and workload, baseline, one changed variable, observed
results, error and cleanup behavior, visualization, interpretation, surprise,
limitations, and next experiment. Link official documentation and distinguish
documented facts, local observations, design judgments, and learner claims.
Do not upload secrets, tokens, private data, unbounded logs, or claims of
production readiness. A report records evidence; it does not replace review.

## Curriculum

## TS-00 — TypeScript writing loop and diagnostics

**Infrastructure connection:** Small reliable transformations are the prerequisite for every later request, queue, and tooling exercise.

**Concepts:** Expressions, statements, `let`/`const`, everyday primitive types, functions, control flow, inference, compiler diagnostics, and JavaScript output.

**Learner-owned work:** Handwrite short transformations from specifications, predict outputs, fix an edge case, and explain one diagnostic without using an escape hatch.

**Experiment/measurement:** Compare fixed-input implementations only after checking equal behavior; report a table of outputs and bounded timings.

**Limitations/reflection:** A successful type check does not prove emitted JavaScript or runtime behavior.

## TS-01 — Data structures, aliasing, and state

**Infrastructure connection:** Arrays, objects, `Map`, `Set`, counters, caches, queues, configuration snapshots, and lifecycle state hold platform data.

**Concepts:** Object shapes, optional properties, mutation, references, copying, structural typing, readonly views, invariants, and discriminated state.

**Learner-owned work:** Implement a bounded FIFO, counters, an immutable-looking snapshot boundary, and a lifecycle model; identify which references may mutate.

**Experiment/measurement:** Vary capacity and producer/consumer rates; record depth, rejection, ordering, and invariant preservation as an ASCII timeline.

**Limitations/reflection:** `readonly` is a static constraint and does not freeze nested runtime objects.

## TS-02 — Functions as contracts

**Infrastructure connection:** Handlers, validators, retry decisions, cache operations, callbacks, and pipeline stages are function contracts.

**Concepts:** Parameter and return types, overloads, generic functions, callbacks, iterables, generators, closures, side effects, and `unknown` results.

**Learner-owned work:** Write parsers, retry-decision functions, callbacks, and small pipeline stages with explicit preconditions and postconditions.

**Experiment/measurement:** Inject malformed input and dependency failure; compare thrown errors, rejected promises, tagged results, and fallback values.

**Limitations/reflection:** A function signature cannot enforce behavior from an untyped caller or external service.

## TS-03 — Modules, packages, environments, and CLIs

**Infrastructure connection:** Module boundaries determine reusable libraries, CLIs, configuration, deployment packaging, and reproducibility.

**Concepts:** `import`/`export`, default versus named exports, package entry points, `package.json` scripts as an external boundary, environment variables, ESM/CommonJS, and import-time effects.

**Learner-owned work:** Design a small module graph and CLI contract, documenting the intended `tsconfig`, package, runtime, and test commands without creating them here.

**Experiment/measurement:** In a future selected environment, type-check, emit, and load the same boundary in a fresh process under its declared module mode.

**Limitations/reflection:** TypeScript module syntax and Node module resolution are related but not identical authorities.

## TS-04 — Errors, cleanup, and resource lifetime

**Infrastructure connection:** Failed requests must release timers, streams, files, workers, and queue ownership without hiding failure.

**Concepts:** `unknown` in catches, custom error shapes, rejected promises, `finally`, cleanup ownership, idempotent teardown, and failure propagation.

**Learner-owned work:** Model success, permanent failure, transient failure, cancellation, and cleanup paths; handwrite a lifecycle state transition.

**Experiment/measurement:** Inject failure before, during, and after an effect; measure settled operations, cleanup count, and retained work.

**Limitations/reflection:** Type annotations do not make cleanup automatic and cannot classify arbitrary runtime errors by themselves.

## TS-05 — Types, narrowing, unions, generics, and runtime contracts

**Infrastructure connection:** Explicit valid states reduce ambiguous request, response, queue, and rollout behavior.

**Concepts:** Everyday types, `unknown`, control-flow narrowing, type predicates, discriminated unions, exhaustive checks, generic constraints, `any`, assertions, non-null assertions, and erased types.

**Learner-owned work:** Handwrite a tagged request/result model, generic bounded collection, and exhaustive dispatcher; keep assertions at named boundaries.

**Experiment/measurement:** Parse external JSON as `unknown`, compare unchecked assertions with explicit runtime validation, and record accepted/rejected cases.

**Limitations/reflection:** TypeScript types are erased; `.d.ts` and compile-time narrowing do not validate network data or JavaScript behavior at runtime.

## TS-06 — Testing and regression control

**Infrastructure connection:** Tests protect request contracts, invalid-input handling, queue invariants, and operational regressions.

**Concepts:** Unit versus integration boundaries, deterministic fixtures, type checking, emitted-code checks, fresh-process checks, invalid inputs, cleanup, and test-runner boundaries.

**Learner-owned work:** Specify success, malformed input, empty/boundary, repeated execution, ordering, cancellation, and cleanup assertions before implementing focused tests.

**Experiment/measurement:** Run the repository's selected package script, compiler configuration, tests, and a fresh-process runtime check as separate evidence.

**Limitations/reflection:** A passing test suite does not prove race freedom, production capacity, or learner understanding; no test runner is selected in this repository.

## TS-07 — Files, structured data, and safe tooling

**Infrastructure connection:** Configuration, manifests, fixtures, logs, and subprocesses cross trust and serialization boundaries.

**Concepts:** Node filesystem/process APIs, JSON parsing, `unknown`, runtime schemas, encoding, path handling, exit status, and shell-injection boundaries.

**Learner-owned work:** Design a safe configuration loader and subprocess contract with explicit validation and cleanup; do not trust a cast from parsed JSON.

**Experiment/measurement:** Compare valid, missing, malformed, oversized, and unexpected fields; measure failure classification and retained temporary state.

**Limitations/reflection:** Type declarations describe expected shapes but cannot sanitize files, command arguments, or subprocess output.

## TS-08 — Request lifecycles, framing, and partial I/O

**Infrastructure connection:** Network requests encounter partial reads, malformed messages, timeouts, and peer failure.

**Concepts:** Node socket/stream boundaries, buffers, framing, backpressure, event ordering, request state unions, malformed messages, and cleanup.

**Learner-owned work:** Model a framed request lifecycle and its valid transitions; handwrite validation and an explicit close/error path.

**Experiment/measurement:** Feed fragmented, combined, truncated, and malformed frames; record state, bytes accepted, rejection, and cleanup.

**Limitations/reflection:** A well-typed parser does not prove the peer obeys the protocol or that the host stream implementation meets a production workload.

## TS-09 — Promises and event-loop scheduling

**Infrastructure connection:** Async request work, timers, I/O callbacks, retries, and shutdown depend on scheduling and rejection behavior.

**Concepts:** Promises, `async`/`await`, microtasks, timers, event-loop turns, ordering, concurrency without parallel JavaScript execution, and unhandled rejection boundaries.

**Learner-owned work:** Predict and then handwrite a small schedule with sequential, parallel, and dependent operations; handle every rejection.

**Experiment/measurement:** Trace labeled synchronous, microtask, timer, and I/O events in a fresh process; compare predicted and observed order.

**Limitations/reflection:** TypeScript does not guarantee scheduling, fairness, latency, or absence of event-loop blockage.

## TS-10 — Node worker boundaries

**Infrastructure connection:** CPU-heavy or isolated work may cross a Node worker boundary with explicit ownership and shutdown.

**Concepts:** `worker_threads`, messages, structured clone/serialization, transfer boundaries, worker lifecycle, errors, termination, and event-loop versus worker execution.

**Learner-owned work:** Design a typed message protocol and worker ownership state machine, including startup failure, cancellation request, and termination.

**Experiment/measurement:** Compare a sequential baseline with a bounded worker case; measure startup, serialization, throughput, and termination cleanup.

**Limitations/reflection:** Typed messages do not validate hostile or stale messages at runtime, and workers do not automatically make shared resources safe.

## TS-11 — Queues, worker pools, overload, and shutdown

**Infrastructure connection:** Bounded queues and explicit worker ownership protect services from overload and shutdown leaks.

**Concepts:** Queue capacity, backpressure, worker-pool state, promises, fairness, rejection, cancellation, drain, and graceful versus forced shutdown.

**Learner-owned work:** Handwrite a bounded work protocol with accepted, rejected, running, completed, failed, and cancelled states.

**Experiment/measurement:** Vary producer rate, consumer rate, capacity, and shutdown timing; record queue depth, drops, latency, and unfinished work.

**Limitations/reflection:** The event loop is not a queue policy; application-level capacity and ownership must be explicit.

## TS-12 — Timeouts, retries, rate limits, and fallbacks

**Infrastructure connection:** Reliability guardrails constrain slow, failing, duplicate, and overloaded dependencies.

**Concepts:** `AbortController`/`AbortSignal`, timeout races, retry classification, jitter, idempotency, token/window limits, isolation, and fallback contracts.

**Learner-owned work:** Design a cancellation-propagating request policy that distinguishes transient, permanent, and caller-aborted outcomes.

**Experiment/measurement:** Vary delay, failure class, retry count, cancellation timing, and rate; measure attempts, deadlines, dropped work, and cleanup.

**Limitations/reflection:** Abort signals are conventions that cooperating code must observe; they do not forcibly stop arbitrary synchronous work.

## TS-13 — Caching, TTL, invalidation, and consistency

**Infrastructure connection:** Caches trade latency and load for freshness, memory, and consistency risk.

**Concepts:** Generic cache entries, keys, TTL clocks, stale data, invalidation, stampede control, serialization, mutation, and concurrent refresh.

**Learner-owned work:** Model fresh, stale, missing, refreshing, failed, and invalidated states with explicit ownership and value validation.

**Experiment/measurement:** Vary TTL, hit rate, invalidation timing, and dependency delay; record freshness, latency, origin load, and duplicate refreshes.

**Limitations/reflection:** A static cache value type cannot prove serialized data or freshness; local time and process scope limit conclusions.

## TS-14 — API contracts, compatibility, and idempotency

**Infrastructure connection:** Stable APIs make request platforms evolvable, paginatable, and safe under retries.

**Concepts:** Discriminated responses, generic pagination, versioning, compatibility, declaration surfaces, idempotency keys, validation, and error contracts.

**Learner-owned work:** Handwrite a versioned API model with exhaustive response handling and an idempotent operation contract.

**Experiment/measurement:** Replay requests, vary page boundaries and unknown fields, and compare old/new consumers under valid and invalid responses.

**Limitations/reflection:** Compatibility of declarations does not prove compatibility of emitted JavaScript, wire behavior, or deployment configuration.

## TS-15 — Reusable libraries and declaration boundaries

**Infrastructure connection:** Shared libraries reduce duplicated reliability policy while making misuse and package boundaries visible.

**Concepts:** Generic APIs, module exports, `package.json` entry points, `.d.ts`, declaration emit, `unknown` at boundaries, and API misuse resistance.

**Learner-owned work:** Design a small library surface, document invariants and escape hatches, and predict its generated declaration boundary.

**Experiment/measurement:** In a selected toolchain, compare source type checking, declaration output, emitted JavaScript, and fresh-process consumer loading.

**Limitations/reflection:** A `.d.ts` file describes an API for tools; it does not implement, validate, or secure the underlying JavaScript.

## TS-16 — Observability and debugging systems

**Infrastructure connection:** Operators need bounded, structured evidence for request failures, queues, cancellation, and event-loop delay.

**Concepts:** Structured logs, metrics, traces, correlation identifiers, explicit units, bounded cardinality, redaction, `Error` causes, and diagnostic channels.

**Learner-owned work:** Define event names and metrics for one request lifecycle, including timeout, retry, cancellation, queue depth, and shutdown.

**Experiment/measurement:** Instrument the smallest boundary and compare normal, rejected, cancelled, and overloaded runs without logging secrets or raw untrusted data.

**Limitations/reflection:** Diagnostic output is evidence about a workload and instrumentation path, not proof of absence of hidden failures.

## TS-17 — Performance, profiling, and tail latency

**Infrastructure connection:** Throughput, resource use, and tail latency determine safe capacity and user-visible reliability.

**Concepts:** Baselines, warm-up, throughput, latency distributions, p95/p99, event-loop delay, allocations, serialization, worker overhead, and sampling noise.

**Learner-owned work:** Predict the bottleneck and define workload, metric, unit, and acceptance boundary before selecting a benchmark or profiler.

**Experiment/measurement:** Change one factor after a baseline; record distributions, CPU/memory observations, queueing, and confidence limits rather than only averages.

**Limitations/reflection:** A local benchmark cannot establish production capacity, scalability, cost, or cross-machine tail behavior.

## TS-18 — Progressive delivery, verification, and rollback

**Infrastructure connection:** Changes to types, emitted modules, runtime validation, and operational policy require gates and reversible rollout.

**Concepts:** Verification stages, compatibility checks, feature flags, canaries, package/module mode, fresh-process smoke tests, rollback state, and evidence promotion.

**Learner-owned work:** Design a rollout state machine with static, emitted-code, runtime, load, and rollback gates; identify who owns each decision.

**Experiment/measurement:** Simulate success, failed verification, cancellation, and rollback; record state transitions and whether old behavior remains available.

**Limitations/reflection:** A simulated rollback does not prove deployment orchestration, data compatibility, or operator readiness.

## TS-19 — Safe request-platform capstone and operational review

**Infrastructure connection:** The capstone combines typed contracts, runtime validation, async scheduling, bounded work, cancellation, observability, load, and rollback.

**Concepts:** End-to-end request lifecycle, discriminated outcomes, module/package boundaries, promise and worker ownership, queue invariants, idempotency, evidence review, and operational limits.

**Learner-owned work:** Design and handwrite the core request-platform contract and review plan; use collaboration for architecture and assist only for repetitive harness/report plumbing.

**Experiment/measurement:** Run the approved ladder from single request through bounded overload, cancellation, worker shutdown, fresh-process module loading, and rollback simulation; report results and evidence boundaries.

**Limitations/reflection:** The capstone can demonstrate a bounded local contract only. It cannot prove TypeScript or Node version compatibility, production readiness, capacity, security, learner understanding, or absence of concurrency defects without separately selected evidence.

## Official documentation anchors

- TypeScript Handbook: <https://www.typescriptlang.org/docs/handbook/>; see Everyday Types, More on Functions, Object Types, Generics, and Modules.
- Narrowing: <https://www.typescriptlang.org/docs/handbook/2/narrowing.html>.
- TSConfig Reference: <https://www.typescriptlang.org/tsconfig>; see `strict`, `module`, `moduleResolution`, `target`, `declaration`, and `noEmit`.
- Declaration Files: <https://www.typescriptlang.org/docs/handbook/declaration-files/introduction.html>.
- Node CommonJS modules: <https://nodejs.org/api/modules.html>.
- Node ECMAScript modules: <https://nodejs.org/api/esm.html>.
- Node timers and scheduling: <https://nodejs.org/api/timers.html>.
- Node globals and `AbortSignal`: <https://nodejs.org/api/globals.html>.
- Node worker threads: <https://nodejs.org/api/worker_threads.html>.

These sources document language, compiler, and host behavior; their current
online versions do not establish a repository toolchain version. Any selected
exercise must record its actual compiler, Node runtime, package configuration,
and test environment.

## Evidence boundaries

This document proves only that a scoped learning-plan structure and intended
evidence contract were written. It does not prove TypeScript compilation,
emitted JavaScript, ESM/CommonJS resolution, runtime validation, promise or
event-loop behavior, worker correctness, cancellation correctness, benchmark
results, race/deadlock freedom, learner understanding, production capacity,
security, or production readiness. No experiments, manifests, dependencies,
toolchain setup, or generated artifacts are created here.