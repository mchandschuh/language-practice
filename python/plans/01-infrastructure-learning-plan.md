# Python Infrastructure Learning Plan

This is the Python adaptation of the historical draft at commit
`670bafa6a5621f4f14ef91b8d71b6d844478fbb3`. It preserves the shared
infrastructure objective and evidence contract while stating Python behavior
only within documented, version-appropriate boundaries.

This track develops Python writing fluency while building intuition for
dependable infrastructure, platform engineering, backend systems, distributed
systems, reliability engineering, and operational tooling:

- platform foundations and shared frameworks;
- high concurrency and unpredictable traffic;
- caching, consistency, queueing, and backpressure;
- reliability guardrails, rate limiting, load shedding, isolation, and fallbacks;
- throughput, latency distributions, tail behavior, and predictable cost;
- observability and debugging systems;
- safe rollout, automated verification, and rollback;
- stable APIs and enforceable contracts; and
- end-to-end ownership from design through operational maturity.

The goal is not to become a Python syntax typist by skipping understanding. The
goal is to be able to write enough Python fluently to investigate, model,
measure, debug, and improve dependable systems.

## Learner profile and teaching strategy

The learner can read and understand most programming languages but currently
writes Python slowly and is not yet fluent. The plan therefore uses short,
repeated writing loops before introducing large infrastructure abstractions.

For most exercises:

1. Read a small specification and relevant official documentation.
2. Predict the behavior and identify a measurable question.
3. Handwrite a small implementation.
4. Run a focused test or experiment.
5. Inspect the output, failure, logs, or visualization.
6. Make one change and rerun the same experiment.
7. Explain what changed, why, and what remains unproven.
8. Write a concise GitHub-compatible experiment report.

Do not optimize for typing speed alone. Optimize for repeated retrieval of
common Python shapes, clear boundaries, and the ability to reason from evidence.

## Assistance modes

Every exercise declares one mode:

### Handwrite

The learner writes the core code. The agent may provide requirements, questions,
API references, tiny syntax hints, test targets, traceback interpretation, and
review feedback, but not a drop-in solution.

Use this for functions, state models, exception paths, context managers, queue
semantics, retry policies, cache behavior, async coordination, contract
validation, and rollout state machines.

### Collaborate

The learner remains the primary implementer. The agent may propose designs,
inspect a partial attempt, suggest experiments, identify hypotheses, and show
small targeted fragments.

Use this for API design, failure analysis, metrics design, profiling
interpretation, concurrency diagnosis, and system tradeoffs.

### Assist

The agent may generate repetitive scaffolding after the contract is understood.
The agent must explain what the generated code measures or does and preserve a
learning checkpoint.

Use this for fixtures, load generators, report formatting, plotting boilerplate,
serialization, adapters, and repetitive CLI plumbing.

Every exercise also records the learner's contract: prerequisite concepts, what
the learner writes or inspects, permitted assistance, an observable validation,
edge cases and failure modes, and a short explanation or reflection checkpoint.

## Version and dependency policy

The repository currently has no Python manifest, lockfile, or declared Python
version. The observed local interpreter is Python 3.14.4. Treat that as an
environment observation, not yet as the repository contract.

Before adding exercises that depend on version-specific behavior:

1. Decide and document the supported Python version.
2. Create the environment with `python -m venv` rather than modifying the
   global interpreter.
3. Prefer the standard library while learning core behavior.
4. Add third-party dependencies only for a deliberate learning objective such
   as visualization or a realistic load generator.
5. Record the Python version, operating system, dependency versions, and command
   in each experiment report.

Python 3.14.4 is only a local observation. The repository has no declared Python
version, so Python 3.15 documentation or newer features must not be treated as
available until a repository version is selected and the feature is checked
against it. In particular, executor and profiling behavior must be described
for the version actually selected, not inferred from a newer interpreter.

Annotations and type hints describe intended contracts to readers and static
tools; Python does not use them as general runtime validation. Boundary code
must validate untrusted values explicitly, and experiments must distinguish a
static diagnosis from a runtime observation.

The concurrency model is selected per experiment. Begin with a synchronous
baseline. Use `asyncio` tasks for cooperative I/O concurrency, threads or
thread executors for suitable blocking I/O, and processes for selected CPU-bound
or isolation workloads. The interpreter/runtime implementation and selected
Python version limit what can run in parallel; an async event loop does not make
CPU-bound work parallel, and a passing concurrent run does not prove absence of
races, deadlocks, or shutdown defects.

## Infrastructure concept map

| Infrastructure concept | Python practice |
| --- | --- |
| Platform foundations | reusable modules, package boundaries, protocols, safe defaults |
| Distributed systems | local dependency simulators, process boundaries, request lifecycles |
| High concurrency | `asyncio`, threads, processes, executors, queues |
| Caching and consistency | TTL caches, invalidation, stale reads, snapshots, duplicate refresh |
| Queueing/backpressure | bounded queues, worker pools, overload rejection, graceful shutdown |
| Reliability guardrails | timeouts, retries, rate limiting, load shedding, isolation, fallbacks |
| Performance primitives | repeatable workloads, throughput, latency distributions, tail analysis |
| Reusable workflows | reusable libraries and CLIs that make safe behavior easy |
| Observability/debugging | structured logs, metrics, correlation IDs, traces, failure reports |
| Safe change management | feature flags, rollout gates, verification, rollback simulation |
| API/contract design | annotations, typed boundaries, explicit runtime validation, compatibility, contract tests |
| Shared primitive impact | one well-designed primitive serving multiple fake consumers |

## The experiment contract

Every experiment must answer a question. It is not enough to produce code that
runs or a graph that looks interesting.

Each experiment should define:

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

### Experiment rules

- Start with a sequential or single-worker baseline when concurrency is being
  studied.
- Change one important variable at a time before exploring combinations.
- Use fixed seeds, fixed input sizes, controlled clocks, and bounded durations.
- Warm up before measuring steady-state behavior when appropriate.
- Record environment and version information.
- Preserve raw observations separately from interpretation.
- Repeat enough times to expose noise; do not report one lucky run as a result.
- Stop load tests when error rate, queue growth, memory, or duration crosses the
  declared safety limit.
- Never run a load test against a production or third-party service without
  explicit authorization and a reviewed scope.
- Keep generated results free of credentials, personal data, private endpoints,
  and unbounded logs.
- State whether each observation is a language guarantee, a standard-library
  guarantee, an implementation detail, or only a result of the local workload.

### What to compare

Depending on the section, compare:

- sequential versus concurrent execution;
- bounded versus unbounded queues;
- blocking versus nonblocking submission;
- retry versus no retry;
- cached versus uncached access;
- stale-while-revalidate versus strict freshness;
- one dependency versus an isolated dependency;
- baseline versus optimized code;
- before-rollout versus after-rollout behavior;
- graceful degradation versus fail-fast behavior.

## Load-testing ladder

Load testing begins with a local deterministic simulator, not a web-scale
framework or an external service.

### Level 0 — Micro experiment

Use a function or small component. Measure elapsed time, result correctness,
allocations where relevant, or state transitions.

Typical tools:

- `time.perf_counter()`;
- `timeit` for small snippets;
- `unittest` or a small executable harness;
- `statistics` for simple summaries.

### Level 1 — Component workload

Exercise a queue, cache, parser, retry policy, or executor with controlled
inputs. Record throughput, latency, failures, queue depth, and cleanup state.

### Level 2 — Local service simulation

Use local processes, sockets, or a fake dependency. Vary arrival rate, service
time, dependency errors, timeouts, and worker count. Keep the test bounded and
disposable.

### Level 3 — Optional tool-assisted load test

Only after the component behavior is understood may the plan introduce a
third-party load or visualization tool. Pin the dependency, document why it is
needed, and preserve a standard-library fallback where practical.

Do not confuse a local load generator with production capacity evidence.

## Visualization ladder

Visualization is part of understanding, not decoration.

### Text and Markdown first

Begin with:

- sorted tables;
- ASCII timelines;
- queue-depth snapshots;
- request lifecycle diagrams;
- Markdown result summaries;
- CSV files committed when small and useful.

This keeps the first experiment dependency-free and makes the data easy to
review in GitHub.

### Static files next

When a graph materially improves understanding, produce a small deterministic
artifact such as:

- SVG;
- a self-contained HTML report;
- a carefully scoped PNG;
- a CSV plus a documented plotting command.

Generated images should be reproducible from committed code and small input
data. Prefer source data and a regeneration command over opaque binary output.

### What to visualize

- request latency distributions and tail percentiles;
- throughput as arrival rate changes;
- queue depth over time;
- accepted, rejected, retried, timed-out, and completed work;
- worker utilization and idle time;
- cache hit/miss/stale/eviction rates;
- memory allocation snapshots;
- rollout exposure and verification failures;
- dependency call timelines.

### Visualization honesty

Every visualization must state:

- what each axis and unit means;
- the workload and sample count;
- whether values are measured, simulated, or inferred;
- aggregation choices and omitted data;
- known noise and limitations.

Never smooth away failures, discard outliers without explanation, or use a graph
to imply production behavior that was not measured.

## GitHub experiment documentation

Each substantial experiment should eventually be documented in a directory such
as:

```text
python/experiments/<section>/<experiment-name>/
├── README.md
├── experiment.py
├── test_experiment.py
├── results/
│   ├── raw.json
│   ├── summary.json
│   └── measurements.csv
└── figures/
    ├── README.md
    └── latency.svg
```

Do not create every directory during the initial writing exercises. Add the
smallest structure that helps the current experiment.

### Experiment README template

Each report should include:

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

The report should make the learner's experience visible: what they expected,
what actually happened, what confused them, what the evidence changed, and what
they would try next.

### Artifact rules

- Keep raw data bounded and reviewable.
- Prefer JSON/CSV/Markdown/SVG over large opaque files.
- Record the command, Python version, platform, seed, duration, and workload.
- Do not commit `.venv`, caches, secrets, credentials, private data, or huge
  generated logs.
- Add a regeneration command for every committed figure or summary.
- Keep generated artifacts separate from source and tests.
- Document whether an artifact is canonical evidence or an exploratory output.

## Curriculum

## PY-00 — Python writing loop

**Infrastructure connection:** prerequisite for every later section.

**Concepts:** expressions, statements, indentation, names, values, function
calls, conditionals, loops, small transformations, and traceback reading.

**Handwrite exercises:** produce 5–20 line functions from specifications; predict
outputs; fix one failing edge case; reimplement a small variation from memory.

**Experiment:** compare two implementations on fixed inputs. First verify that
they produce the same result; only then measure them. Record whether the
performance difference is large enough to matter.

**Visualization:** begin with a Markdown table of inputs, outputs, and measured
times. No plotting dependency is needed.

**Report:** one short README describing prediction, result, surprise, and limit.

## PY-01 — Data structures and state

**Infrastructure connection:** queues, caches, counters, configuration, and lifecycle
state.

**Concepts:** built-in collections, `deque`, `defaultdict`, `Counter`, mutation,
aliasing, object identity, copying, dataclasses, and explicit states. Identity
(`is`) answers whether two names refer to the same object; equality (`==`) answers
whether values compare equal. Aliasing means shared references, so mutation may
be visible through every alias; copying or an immutable snapshot changes that
boundary rather than merely changing a type annotation.

**Handwrite exercises:** implement a bounded FIFO, request counters, a mutable
configuration snapshot, and a small lifecycle model.

**Experiment:** vary queue capacity and producer/consumer rates. Observe when
the queue blocks, rejects, or grows. Record state invariants after every run.

**Visualization:** render queue contents and depth as an ASCII timeline before
creating a graph.

**Report:** explain which objects are shared, who may mutate them, and what a
snapshot protects.

## PY-02 — Functions as contracts

**Infrastructure connection:** request handlers, retry policies, cache operations, worker
functions, and validation boundaries.

**Concepts:** signatures, keyword-only parameters, defaults, pure functions,
callbacks, generators, preconditions, postconditions, and side effects.

**Handwrite exercises:** write parsers, retry-decision functions, callbacks, and
small pipeline stages.

**Experiment:** inject malformed input and dependency failures. Compare explicit
error returns, exceptions, and fallback values without hiding the distinction.

**Visualization:** produce a simple request-state transition table or Mermaid
diagram in the README; label inferred transitions separately from observed ones.

## PY-03 — Modules, packages, environments, and CLIs

**Infrastructure connection:** shared libraries, frameworks, reusable tools, and
safe reusable workflows.

**Concepts:** imports, packages, import-time side effects, `__main__`, virtual
environments, installation, CLI parsing, and configuration.

**Handwrite exercises:** turn a script into a package; separate CLI parsing from
pure logic; create a disposable `.venv`; build a small JSON CLI.

**Experiment:** measure startup behavior with and without import-time work. Use
an import graph or dependency table to identify coupling.

**Visualization:** document the package/module dependency graph. Use
`graphlib.TopologicalSorter` to validate dependency order and demonstrate cycle
detection.

**Report:** record Python version, environment creation command, package
boundary, and reproducibility limitations.

## PY-04 — Exceptions, cleanup, and resource lifetime

**Infrastructure connection:** dependency isolation, process/socket/file cleanup, and
failure containment.

**Concepts:** exception hierarchy, chaining, narrow catches, `try/finally`,
context managers, async context managers, `ExitStack`, partial initialization,
and cancellation cleanup.

**Handwrite exercises:** implement a temporary resource context manager, cleanup
after partial initialization, and an async resource wrapper.

**Experiment:** inject failures at each acquisition and release boundary. Record
which resources remain open and whether cleanup is idempotent.

**Visualization:** create a resource-lifetime timeline showing acquire, use,
failure, cleanup, and release.

**Report:** distinguish documented guarantees, observed cleanup behavior, and
unproven external-resource assumptions.

## PY-05 — Types and runtime contracts

**Infrastructure connection:** stable APIs, enforceable contracts, dependency substitution,
and independent evolution.

**Concepts:** annotations, generics, `TypedDict`, protocols, callable/coroutine
types, type guards, static checking, and runtime validation.

**Handwrite exercises:** define request/response contracts, cache protocols, fake
dependencies, and boundary validators.

**Experiment:** pass malformed runtime data through a valid annotation and show
why annotations alone do not enforce the contract.

**Visualization:** document the contract as a request/response table with valid,
invalid, and ambiguous cases.

## PY-06 — Testing and regression control

**Infrastructure connection:** regression controls and automated verification.

**Concepts:** `unittest`, fixtures, cleanup, isolation, deterministic clocks,
fake dependencies, subtests, negative paths, and contract tests.

**Handwrite exercises:** write tests for parsers, state preservation, retry
limits, idempotency, queue shutdown, timeout behavior, and cancellation.

**Experiment:** deliberately break the implementation and confirm each test
actually fails. Capture the failing and passing commands in the report.

**Visualization:** show a compact test matrix mapping behavior to evidence;
avoid coverage graphs until the learner understands assertion quality.

**Report:** state exactly what the tests prove and what they do not prove.

## PY-07 — Files, JSON, subprocesses, and safe tooling

**Infrastructure connection:** developer productivity tooling, service integration, and
safe operational automation.

**Concepts:** file boundaries, JSON validation, subprocess arguments, exit codes,
timeouts, environment variables, shell safety, and child cleanup.

**Handwrite exercises:** build a command wrapper with safe argument handling,
nonzero-exit reporting, timeout cleanup, and malformed-output handling.

**Experiment:** compare successful, slow, missing, malformed, and failing child
processes. Record process state and captured output separately.

**Visualization:** use a process timeline showing parent/child lifecycle,
timeout, termination, and cleanup.

## PY-08 — Sockets and request lifecycles

**Infrastructure connection:** interfaces, data paths, request lifecycle, and dependency
failure.

**Concepts:** blocking/nonblocking sockets, timeouts, partial reads/writes,
framing, connection lifecycle, and malformed messages.

**Handwrite exercises:** implement a tiny length-prefixed local protocol and
handle partial frames, disconnects, slow peers, and invalid messages.

**Load experiment:** use a bounded local client generator to vary concurrent
connections and request arrival rate. Stop on error-rate, duration, or resource
limits.

**Visualization:** plot or tabulate active connections, completed requests,
timeouts, error classes, and latency percentiles over time.

## PY-09 — `asyncio` fundamentals

**Infrastructure connection:** high concurrency and unpredictable request patterns.

**Concepts:** event loops, coroutines, tasks, streams, async queues,
synchronization, timeouts, cancellation, and blocking work.

**Handwrite exercises:** convert a sequential I/O simulator to tasks; add
timeouts; cancel tasks; use a bounded async queue; implement fan-out/fan-in.

**Load experiment:** compare sequential, unconstrained concurrent, and bounded
concurrent versions using the same deterministic fake dependencies.

**Measurements:** total throughput, per-request latency, p50/p95/p99-like
quantiles, cancellation count, queue depth, and dependency failures.

**Visualization:** produce a request timeline and latency-distribution report.
Use SVG/HTML/PNG only after preserving raw JSON/CSV measurements.

## PY-10 — Executors and parallel work

**Infrastructure connection:** throughput, worker systems, and concurrency/failure
handling.

**Concepts:** thread, process, and interpreter executors where supported;
futures, lifecycle, cancellation limits, deadlocks, and serialization.

**Handwrite exercises:** run the same workload sequentially, with threads, and
with processes; diagnose a deliberate deadlock; compare I/O-bound and CPU-bound
work.

**Load experiment:** sweep worker count over a small fixed range. Include startup,
serialization, service, and shutdown time rather than reporting only task time.

**Visualization:** show throughput and latency against worker count, marking the
point where additional workers stop helping or increase tail latency.

**Version note:** interpreter-executor and free-threading behavior must be gated
by the repository's selected Python version. Do not assume Python 3.15 behavior
while the environment is currently Python 3.14.4.

## PY-11 — Queues, worker pools, and backpressure

**Infrastructure connection:** queueing, backpressure, load shedding, and graceful
degradation.

**Concepts:** producer/consumer design, bounded queues, rejection, `task_done`,
`join`, worker lifecycle, graceful shutdown, overload, and queue termination.

**Handwrite exercises:** implement fixed workers, bounded capacity, timeout-aware
submission, per-task failure capture, and explicit shutdown.

**Load experiment:** sweep arrival rate across service capacity. Observe stable
operation, queue growth, rejection, timeout, and recovery after load falls.

**Visualization:** create a time-series chart or Markdown sparkline for arrival
rate, service rate, queue depth, accepted work, rejected work, and completion.

**Report:** explain why unbounded queues can hide overload and why `qsize()` or a
single snapshot is not sufficient evidence of system health.

## PY-12 — Reliability guardrails

**Infrastructure connection:** timeouts, retries, rate limiting, load shedding,
dependency isolation, circuit breakers, and safe fallbacks.

**Concepts:** timeout budgets, retry classification and amplification, token-bucket
rate limiting, bulkheads, circuit states, fallback contracts, and idempotency.

**Handwrite exercises:** implement retry policy functions, a fake clock,
token-bucket rate limiter, dependency bulkhead, fallback, and circuit state
machine.

**Failure-injection experiment:** vary transient failure, permanent failure,
slow response, malformed response, and cancellation. Measure retry amplification,
recovery time, rejected work, and user-visible errors.

**Visualization:** show state transitions and a fault timeline. Annotate when a
guardrail activates and deactivates.

**Report:** explicitly answer whether each operation is safe to retry and whether
duplicate work can occur.

## PY-13 — Caching and consistency

**Infrastructure connection:** caching and consistency fundamentals.

**Concepts:** TTL, invalidation, stale reads, snapshots, duplicate refresh,
bounded capacity, eviction, cache-aside, and write-through.

**Handwrite exercises:** build a monotonic-clock TTL cache, bounded eviction,
stale-while-revalidate behavior, and single-flight refresh.

**Experiment:** vary hit rate, TTL, refresh latency, invalidation timing, and
concurrent requests. Compare latency, freshness, origin load, and stale results.

**Visualization:** report hit/miss/stale/refresh/eviction counters and freshness
over time.

**Report:** state the consistency contract in plain language and identify the
failure mode when the cache and source disagree.

## PY-14 — API and contract design

**Infrastructure connection:** simple APIs, enforceable contracts, safe defaults, and
independent evolution.

**Concepts:** request/response models, validation, error contracts, versioning,
compatibility, pagination, idempotency, injection, and adapters.

**Handwrite exercises:** design a local service contract, validate requests,
separate transport/domain logic, add an idempotency key, and write contract
tests against two implementations.

**Experiment:** replay old clients against new server behavior. Record accepted,
rejected, ambiguous, and degraded cases.

**Visualization:** publish a contract matrix and dependency-boundary diagram.

## PY-15 — Golden-path library design

**Infrastructure connection:** shared frameworks, safe reusable workflows, shared
primitive impact,
and hard-to-misuse interfaces.

**Concepts:** library boundaries, safe defaults, extension points, compatibility,
documentation, examples, and misuse resistance.

**Handwrite exercises:** build a reusable request/queue/cache primitive with two
fake consumers. Keep the consumer code simpler than the underlying primitive.

**Experiment:** compare two consumer implementations before and after adopting
the reusable workflow. Measure code duplication, failure behavior, and operational
signals—not only line count.

**Visualization:** show the before/after dependency and request-flow diagrams.

**Report:** identify which operational lesson became a default and which escape
hatches remain available.

## PY-16 — Observability and debugging systems

**Infrastructure connection:** instrumentation, metrics models, and investigative tooling.

**Concepts:** logger hierarchy, levels, handlers, contextual fields, correlation
IDs, structured events, metrics, queue signals, and symptom-to-cause diagnosis.

**Handwrite exercises:** replace `print()` with module loggers; add request/task
IDs; emit events for retries, timeouts, fallback, queue transitions, and
shutdown.

**Experiment:** inject a latency spike or dependency failure and diagnose it from
logs and metrics without reading the implementation first.

**Visualization:** produce a request timeline, dependency breakdown, latency
histogram/quantile table, and error-class summary.

**Report:** distinguish symptom, evidence, hypothesis, confirmed cause, and
remaining uncertainty.

## PY-17 — Performance and tail latency

**Infrastructure connection:** scalability primitives, tail latency, throughput, bottleneck
analysis, and cost predictability.

**Concepts:** workload design, warm-up, variance, throughput, p50/p95/p99-like
quantiles, queueing delay, service time, profiling, memory allocation, and
benchmark limitations.

**Handwrite exercises:** design a fixed-workload benchmark, record raw samples,
separate service time from queueing delay, and explain one profiling result before
changing the implementation.

**Experiments:**

- use `timeit` for small isolated snippets;
- use a repeatable workload for component measurements;
- use `statistics` for mean, median, quantiles, and spread;
- use `tracemalloc` snapshots to investigate Python allocation behavior;
- use the version-appropriate profiler only after establishing a baseline;
- compare one changed variable at a time.

**Visualization:** preserve raw samples and produce a latency distribution,
throughput curve, queueing/service-time comparison, and memory snapshot table.

**Report:** explain noise, outliers, warm-up, machine variance, workload shape,
and why the result does not prove production capacity.

## PY-18 — Progressive delivery and rollback simulation

**Infrastructure connection:** progressive delivery, automated verification, rollout gates,
operational validation, and fast rollback.

**Concepts:** feature-flag evaluation, exposure cohorts, verification thresholds,
halt conditions, compatible configuration, rollback state, and recovery time.

**Handwrite exercises:** implement feature flags, staged exposure, verification
thresholds, automatic halt, compatible configuration, and state-safe rollback.

**Experiment:** expose a change to 1%, 10%, 50%, and 100% of simulated traffic.
Inject elevated errors or tail latency and verify that the gate halts or rolls
back exposure.

**Visualization:** show exposure percentage, error rate, tail latency, gate
state, and rollback time on one annotated timeline.

**Report:** document the rollback trigger, evidence available to the gate, and
what the simulation cannot prove about deployment systems.

## PY-19 — Safe request platform capstone

**Infrastructure connection:** integrating safe request handling, reliability
guardrails, observability, performance evidence, progressive delivery, and
operational ownership.

**Concepts:** stable request contracts, dependency seams, bounded work, failure
containment, measurable service objectives, rollout gates, and rollback.

**Handwrite exercises:** assemble the capstone in slices, keeping the learner
responsible for the request contract, failure model, queue behavior, guardrails,
and final operational explanation.

**Experiment:** run a bounded local workload through the capstone, inject
dependency failures and overload, and compare accepted, rejected, retried,
timed-out, completed, and rolled-back work.

Build a small local Python platform primitive providing:

- typed request configuration;
- a stable dependency interface;
- timeout handling;
- bounded retries;
- rate limiting;
- bounded queueing and backpressure;
- fallback behavior;
- structured logs and metrics;
- deterministic tests;
- progressive rollout simulation;
- automated verification gates; and
- rollback.

Build it in slices: synchronous contract, fake dependency, timeout/error model,
retry policy, bounded queue, workers, rate limiting, fallback, observability,
performance measurement, rollout, rollback, and final operational review.

The capstone should produce a GitHub-readable report with architecture,
experiments, raw-data references, visualizations, failures encountered,
decisions, validation, and explicit limitations.

## Artifact and GitHub documentation protocol

For each substantial experiment, use:

```text
python/experiments/<section>/<experiment-name>/
├── README.md
├── experiment.py
├── test_experiment.py
├── results/
│   ├── raw.json
│   ├── summary.json
│   └── measurements.csv
└── figures/
    ├── README.md
    └── latency.svg
```

Start with only the files needed for the current experiment. Do not create
empty scaffolding for the entire curriculum.

### Required report sections

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

The report must preserve the learner's experience:

- what was expected;
- what happened;
- what was confusing;
- what evidence changed the mental model;
- what was tried and rejected;
- what should be tested next.

### Artifact rules

- Keep raw data bounded and reviewable.
- Prefer JSON, CSV, Markdown, and SVG over large opaque binaries.
- Record Python version, platform, dependencies, seed, duration, workload, and
  command.
- Do not commit `.venv`, caches, credentials, private data, or unbounded logs.
- Include a regeneration command for every committed figure or summary.
- State whether each artifact is canonical evidence or exploratory output.
- Never smooth away failures or discard outliers without documenting why.

## Standard-library anchors

- Tutorial: `https://docs.python.org/3/tutorial/`
- Language Reference: `https://docs.python.org/3/reference/`
- Style: `https://peps.python.org/pep-0008/`
- Typing: `https://docs.python.org/3/library/typing.html`
- Async I/O: `https://docs.python.org/3/library/asyncio.html`
- Executors: `https://docs.python.org/3/library/concurrent.futures.html`
- Queues: `https://docs.python.org/3/library/queue.html`
- Context managers: `https://docs.python.org/3/library/contextlib.html`
- Subprocesses: `https://docs.python.org/3/library/subprocess.html`
- Testing: `https://docs.python.org/3/library/unittest.html`
- Logging: `https://docs.python.org/3/howto/logging.html`
- Timing: `https://docs.python.org/3/library/timeit.html`
- Statistics: `https://docs.python.org/3/library/statistics.html`
- Memory tracing: `https://docs.python.org/3/library/tracemalloc.html`
- Dependency order: `https://docs.python.org/3/library/graphlib.html`
- Environments: `https://docs.python.org/3/library/venv.html`

## Evidence boundaries

Local experiments can establish behavior for the stated workload, environment,
and implementation. They do not establish production capacity, cloud behavior,
provider reliability, real traffic distributions, security readiness, or user
understanding. Every report must preserve those limitations.
