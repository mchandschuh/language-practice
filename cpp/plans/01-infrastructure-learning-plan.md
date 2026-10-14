# C++ Infrastructure Learning Plan

This plan develops C++ writing fluency in service of dependable infrastructure.
It preserves the shared infrastructure objective and behavioral sequence used by
the accepted Python baseline, but teaches C++'s value semantics, object lifetime,
RAII, ownership, generic programming, standard-library facilities, concurrency,
and undefined-behavior boundaries rather than translating Python APIs.

The track connects language practice to platform and backend foundations,
distributed request lifecycles, queues and backpressure, caching and
consistency, timeouts and retries, API contracts, observability, tail latency,
progressive delivery, rollback, and operational maturity. The plan is a learning
contract, not an implementation repository or a claim of production readiness.

## Learner profile and teaching strategy

The learner can read most programming languages but currently writes C++ slowly
and is building fluency with modern C++. Each exercise uses a short loop:

1. Read a small specification and the relevant authoritative documentation.
2. Predict behavior, ownership, lifetime, failure, and a measurable question.
3. Handwrite the smallest core implementation.
4. Run a focused test or bounded experiment.
5. Inspect diagnostics, output, traces, measurements, or visualization.
6. Change one variable and repeat the same experiment.
7. Explain what the language guarantees, what the library guarantees, and what
   remains an implementation detail or unproven assumption.
8. Write a concise GitHub-compatible report.

Compilation is not understanding, and a passing local experiment is not
production evidence. Warnings, sanitizers, tests, and benchmarks are separate
evidence layers.

## Assistance modes

Every exercise declares one mode and keeps the learner responsible for the
contract and explanation.

### Handwrite

The learner writes the core C++ code. The agent may provide requirements,
official references, tiny syntax hints, compiler-error interpretation, test
targets, and review feedback, but not a drop-in solution. Use this for values,
initialization, ownership, classes, templates, error handling, queues, locks,
and small algorithms.

### Collaborate

The learner remains the primary implementer. The agent may propose designs,
inspect a partial attempt, identify lifetime or synchronization hazards, suggest
experiments, and show small targeted fragments. Use this for API design,
failure analysis, generic interfaces, contention diagnosis, and report
interpretation.

### Assist

The agent may generate repetitive scaffolding after the contract is understood,
but must explain what it measures or does and preserve a learner checkpoint. Use
this for fixtures, bounded load generators, report formatting, plotting
boilerplate, adapters, and repetitive CLI plumbing. Socket, JSON, logging,
metrics, plotting, package-manager, build-system, and async-runtime code is
always an explicitly selected external boundary, not an assumed C++ facility.

Every exercise records prerequisites, learner-owned work, permitted assistance,
observable validation, edge cases, failure modes, and a short reflection.

## Version and tooling policy

The repository declares no C++ standard, compiler version, build system,
package manager, dependency, or asynchronous runtime. Do not turn a local
compiler command into a repository contract. GCC 15.2 and Clang 21.1.8 are local
observations only. When an exercise uses a standard-version feature, name the
minimum intended standard for that exercise, verify compiler support separately,
and record portability limits.

ISO C++'s getting-started material identifies multiple compiler choices and
points learners toward modern tours and reference material. This plan therefore
uses standard-library-first examples and treats build configuration, dependency
resolution, socket APIs, JSON implementations, structured logging, metrics,
plotting, and asynchronous runtimes as external tooling choices.

## Infrastructure concept map

| Infrastructure concern | C++ practice | Boundary |
| --- | --- | --- |
| Values and state | value types, references, pointers, invariants, containers | language and standard library |
| Resource ownership | RAII, destructors, smart pointers, explicit ownership | language plus library types |
| API contracts | classes, abstract interfaces, concepts, error policy | learner-designed interface |
| Request paths | algorithms, ranges, exceptions or error returns | standard library; sockets/JSON external |
| Concurrency | threads, `jthread`, futures, locks, condition variables, atomics | standard library; async runtime external |
| Reliability | bounded queues, cancellation, timeout/retry policy, shutdown | application contract |
| Observability | event schema, metrics, traces, measurements, plots | logging/metrics/plotting tools external |
| Delivery | verification gates, rollout and rollback simulation | local experiment, not deployment evidence |

## The experiment contract

Every experiment answers one question with a stated hypothesis, prerequisites,
fixed inputs, workload, environment, command, result, interpretation, and
limitation. First make behavior correct and observable; only then compare
performance. Preserve failures, meaningful outliers, ownership transitions,
queue depth, cancellation, timeout, retry, and shutdown outcomes.

For C++, distinguish compiler acceptance, warnings, sanitizer findings, test
results, and runtime measurements. A sanitizer run can expose selected classes
of defects but cannot prove absence of undefined behavior or races. A benchmark
on one machine and workload cannot prove production capacity, cost, scalability,
or tail behavior. Never intentionally execute an uncontrolled invalid-memory or
data-race experiment; use a minimal reviewed reproduction and document that it
is undefined or erroneous.

## Load-testing ladder

1. **Micro experiment:** one deterministic input and one language or library
   question; use a table of expected and observed behavior.
2. **Sequential baseline:** one request path, one worker, bounded data, explicit
   cleanup, and a repeatable command.
3. **Controlled variation:** vary one input such as payload size, queue capacity,
   allocation strategy, or failure rate.
4. **Concurrency ladder:** add workers, contention, cancellation, and shutdown
   ordering while recording queue depth, completed work, rejected work, errors,
   and latency percentiles.
5. **Failure ladder:** inject timeouts, dependency failures, duplicate work,
   backpressure, and restart/rollback events with a clear stop condition.

All levels are local evidence for the stated implementation, machine, and
workload. They are not production-capacity or universal race-freedom evidence.

## Visualization ladder

Begin with Markdown tables, compiler diagnostics, ownership/lifetime timelines,
state-transition tables, and ASCII queue-depth timelines. Then use bounded CSV
or JSON summaries, histograms, percentile tables, concurrency timelines, and
only then a plot selected through external tooling. Preserve raw observations,
units, clock source, seed, duration, workload, compiler/toolchain, and command.
Do not smooth failures or discard outliers without explaining why. Structured
logging, metrics, tracing, and plotting libraries are external choices.

## GitHub experiment documentation

Future experiments may use the smallest required structure under
`cpp/experiments/<section>/<experiment-name>/`; this plan creates none:

```text
README.md
source files and tests appropriate to the experiment
results/raw.json or measurements.csv
figures/README.md and selected figures
```

Each report contains:

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

Reports state what the learner expected, what happened, what changed the mental
model, what was tried or rejected, and what remains to test. Do not commit
credentials, private data, caches, unbounded logs, dependencies, manifests, or
generated artifacts as part of this plan.

## Curriculum

## CPP-00 — C++ writing loop and diagnostics

**Infrastructure connection:** establishes the writing and diagnostic loop used by every later request, worker, and tool.

**Concepts:** expressions, statements, functions, control flow, names, values, `const`, compiler diagnostics, warnings, and prediction before execution.

**Learner-owned work:** handwrite small transformations and a command-line-independent test harness; predict outputs, fix an edge case, and explain one warning.

**Experiment/measurement:** compare two correct implementations on fixed inputs only after checking equivalent results; record timing method and why the result is limited.

## CPP-01 — Values, references, pointers, and state

**Infrastructure connection:** makes request state, counters, configuration snapshots, and queue items understandable and ownership-visible.

**Concepts:** value semantics, lvalues and rvalues, references, pointers, `nullptr`, `const`, aliasing, mutation, copying, and non-owning versus owning handles.

**Learner-owned work:** handwrite a bounded FIFO and a state snapshot using values, a `const` reference, and a pointer; annotate who may mutate and who owns each object.

**Experiment/measurement:** vary aliasing and copying; inspect state invariants and explain why a pointer or reference does not by itself transfer ownership.

## CPP-02 — Initialization, scope, lifetime, and undefined behavior

**Infrastructure connection:** prevents dangling request context, use-after-scope, invalid configuration, and lifetime bugs in service components.

**Concepts:** default/value/list initialization, scope, storage duration, object lifetime, temporary lifetime, references, dangling pointers, bounds, evaluation, and undefined behavior.

**Learner-owned work:** handwrite small examples with explicit initialization and scope; classify each access as valid, implementation-defined, unspecified, or undefined without relying on folklore.

**Experiment/measurement:** use warnings and a sanitizer on a reviewed minimal reproduction; record that one detected defect is evidence of that run, not proof of complete safety.

## CPP-03 — Translation units, interfaces, and tooling boundaries

**Infrastructure connection:** separates reusable platform libraries from build orchestration and prevents hidden dependency or initialization coupling.

**Concepts:** declarations and definitions, headers, namespaces, linkage, translation units, the one-definition rule, command-line interfaces, and configuration.

**Learner-owned work:** handwrite a tiny library boundary with a pure interface and separate implementation; document the selected compiler invocation without making it repository policy.

**Experiment/measurement:** compare startup and dependency behavior for explicit versus import-like initialization; record build-system, package-manager, JSON, and socket choices as external tooling decisions.

## CPP-04 — RAII, destructors, and resource lifetime

**Infrastructure connection:** gives files, locks, sockets, transactions, and request scopes deterministic cleanup during success, failure, and cancellation.

**Concepts:** constructors, destructors, scope exit, RAII, exception safety, acquisition failure, cleanup ordering, rule of zero, and resource wrappers.

**Learner-owned work:** handwrite a small scope-bound resource wrapper or use an appropriate standard RAII type; demonstrate normal return and failure cleanup without a global owner.

**Experiment/measurement:** record acquisition, use, and release events under early returns and exceptions; explain which cleanup is guaranteed and which external resource contract remains unproven.

## CPP-05 — Classes, invariants, and abstract interfaces

**Infrastructure connection:** supports stable request handlers, dependency boundaries, state machines, and reusable platform components.

**Concepts:** class invariants, access control, constructors, member initialization, const member functions, virtual functions, abstract classes, `override`, and `final`.

**Learner-owned work:** handwrite a class with an invariant and a narrow abstract dependency interface; test valid, invalid, and repeated operations.

**Experiment/measurement:** substitute a fake dependency and compare observable state transitions; report which polymorphism and ownership decisions are design judgments rather than language guarantees.

## CPP-06 — Copy/move behavior and ownership

**Infrastructure connection:** controls data transfer, queue handoff, configuration snapshots, and accidental expensive or invalid request-state copies.

**Concepts:** copy constructor and assignment, move constructor and assignment, copy elision, moved-from objects, rule of zero/three/five, and value transfer.

**Learner-owned work:** handwrite a small value type and trace its copy/move operations around containers and function calls; state the post-move contract.

**Experiment/measurement:** count or log operations for lvalue and rvalue paths, then measure only equivalent work; explain why optimization and elision are not a portable semantic promise in every context.

## CPP-07 — Smart pointers and explicit ownership

**Infrastructure connection:** makes component lifetime, dependency ownership, cache entries, and shutdown responsibility explicit.

**Concepts:** `unique_ptr`, `shared_ptr`, `weak_ptr`, custom deleters, ownership transfer, reference cycles, raw non-owning pointers/references, and lifetime extension.

**Learner-owned work:** handwrite ownership diagrams and a small component graph using the narrowest suitable smart pointer; include destruction order and cycle analysis.

**Experiment/measurement:** observe destruction and retained objects across success and shutdown paths; explain why `shared_ptr` is not a general solution for ownership design or thread safety.

## CPP-08 — Containers, iterators, and external data boundaries

**Infrastructure connection:** supports bounded queues, caches, request payloads, indexes, and safe handling of data received from dependencies.

**Concepts:** `vector`, `deque`, associative and unordered containers, strings, views, iterators, capacity, invalidation, complexity, and container thread-safety boundaries.

**Learner-owned work:** handwrite a bounded queue or cache using a suitable standard container; document invalidation and ownership of every returned view or iterator.

**Experiment/measurement:** vary capacity and workload; record depth, rejection, allocation observations, and invalidation cases. Sockets and JSON parsing remain external library/tooling choices.

## CPP-09 — Algorithms, ranges, and request transformations

**Infrastructure connection:** expresses filtering, routing, batching, indexing, and validation without hiding work or ownership in ad hoc loops.

**Concepts:** standard algorithms, predicates, projections, iterators, ranges, views, lazy evaluation, composition, concepts, and invalidation.

**Learner-owned work:** handwrite an algorithm pipeline for valid and invalid request records, first with conventional algorithms and then with ranges where supported.

**Experiment/measurement:** compare correctness, allocations, and latency for equivalent workloads; preserve the sequential baseline and distinguish measured implementation behavior from standard guarantees.

## CPP-10 — Exceptions, `noexcept`, and error-return alternatives

**Infrastructure connection:** defines how request failure, dependency failure, cancellation, and recovery cross API boundaries.

**Concepts:** throwing and catching, stack unwinding, exception safety, `noexcept`, termination consequences, error codes, `optional`, and result/expected-style error returns where supported.

**Learner-owned work:** handwrite the same small operation with an exception contract and an explicit error-return contract; test invalid input, dependency failure, and cleanup.

**Experiment/measurement:** compare propagation, logging, and caller handling; choose one policy for one boundary and explain why mixing policies invisibly is unsafe.

## CPP-11 — Threads, `jthread`, futures, and task ownership

**Infrastructure connection:** introduces worker ownership, parallel request work, cancellation propagation, and graceful shutdown.

**Concepts:** `std::thread`, `std::jthread`, stop tokens, futures, promises, `async`, joinability, task lifetime, exception propagation, and scheduling limits.

**Learner-owned work:** handwrite a sequential baseline, then a bounded worker example with explicit join or stop behavior; identify the owner of every task and result.

**Experiment/measurement:** vary worker count and cancellation timing; record completed, abandoned, and failed work, plus shutdown latency. An async runtime is an external choice, not implied by coroutines or futures.

## CPP-12 — Mutexes, condition variables, queues, and backpressure

**Infrastructure connection:** models bounded work queues, producer/consumer flow, overload, blocking, rejection, and graceful drain.

**Concepts:** mutexes, lock guards, condition variables, predicates, spurious wakeups, bounded capacity, notification, lock scope, and deadlock avoidance.

**Learner-owned work:** handwrite a bounded producer/consumer queue with success, full, stop, and drain states; test normal completion, cancellation, and shutdown ordering.

**Experiment/measurement:** vary producer/consumer rates and capacity; visualize queue depth, blocked time, rejection, and completion while preserving a single-worker baseline.

## CPP-13 — Atomics, memory ordering, and data-race boundaries

**Infrastructure connection:** supports counters, flags, metrics, and low-level coordination while making unsafe shared state visible.

**Concepts:** atomic objects, data races, happens-before, memory orders, compare/exchange, false sharing, lock-free claims, and the limits of tests and sanitizers.

**Learner-owned work:** handwrite a stop flag and an atomic counter, then state why unrelated compound state still needs synchronization.

**Experiment/measurement:** compare a mutex-protected and atomic counter under controlled contention; use a relevant sanitizer or stress run when available and report that passing runs do not prove race freedom.

## CPP-14 — Reliability guardrails

**Infrastructure connection:** turns dependency failures and overload into bounded timeout, retry, rate-limit, load-shed, fallback, and recovery behavior.

**Concepts:** deadlines, cancellation, bounded retries with jitter, idempotency, duplicate work, rate limiting, isolation, fallback, and shutdown under partial failure.

**Learner-owned work:** handwrite a fake dependency and a policy layer with explicit timeout, retry, rejection, and fallback states; never retry a non-idempotent effect without a deliberate contract.

**Experiment/measurement:** inject transient and permanent failures, vary queue capacity and retry budget, and report errors, retries, duplicate effects, latency, and remaining uncertainty.

## CPP-15 — API and reusable library design

**Infrastructure connection:** makes shared platform primitives safe to reuse across handlers, workers, tools, and operational workflows.

**Concepts:** concepts, templates, generic interfaces, regular/value types, dependency injection, contracts, namespace boundaries, ABI/API considerations, and error policy.

**Learner-owned work:** handwrite a constrained generic component and an abstract or value-based dependency boundary; document its ownership, lifetime, complexity, and invalid inputs.

**Experiment/measurement:** instantiate the component with multiple types and a fake dependency; inspect diagnostics and behavior while avoiding claims about ABI portability or universal concept satisfaction.

## CPP-16 — Observability and debugging systems

**Infrastructure connection:** makes request identity, queue depth, retries, failures, contention, and shutdown diagnosable.

**Concepts:** event schemas, bounded cardinality, correlation identifiers, clock choice, assertions, sanitizers, debugger/profiler selection, metrics, traces, and structured logging.

**Learner-owned work:** handwrite a minimal diagnostic boundary with redacted, stable fields and a report that separates compiler, sanitizer, test, and measurement evidence.

**Experiment/measurement:** inject a known failure and locate it from bounded observations; structured logging, metrics, tracing, and plotting implementations are external tooling choices.

## CPP-17 — Performance and tail latency

**Infrastructure connection:** links allocation, locality, contention, queueing, throughput, and tail behavior to capacity decisions without premature optimization.

**Concepts:** workload design, warm-up, variance, percentiles, allocation, cache locality, contention, batching, throughput, and latency distributions.

**Learner-owned work:** handwrite a sequential baseline and one measured change; verify equivalent results and preserve meaningful work before timing.

**Experiment/measurement:** run a bounded load ladder and report p50/p95/p99 or an appropriate alternative, machine/compiler details, noise, and why the result does not prove production capacity.

## CPP-18 — Progressive delivery and rollback simulation

**Infrastructure connection:** practices verification gates, staged rollout, compatibility, rollback, and recovery from a bad change.

**Concepts:** versioned contracts, feature flags, health signals, canary decisions, invariant checks, rollback state, durable versus in-memory state, and failure injection.

**Learner-owned work:** handwrite a local rollout state machine with a fake dependency and explicit promote, pause, reject, and rollback paths.

**Experiment/measurement:** inject verification failures and restarts; record state transitions and recovery time while clearly labeling the simulation as not deployment-system evidence.

## CPP-19 — Safe request-platform capstone

**Infrastructure connection:** integrates a bounded request path, ownership, reliability, concurrency, observability, measurement, rollout, and operational review.

**Concepts:** classes and abstract interfaces, templates and concepts, containers, algorithms and ranges, exceptions or error returns, RAII, smart pointers, threads, `jthread`, futures, mutexes, condition variables, atomics, cancellation, shutdown, and explicit coroutine/runtime boundaries.

**Learner-owned work:** handwrite the core synchronous contract, fake dependency, error model, bounded queue, worker lifecycle, rate limit, fallback, observations, rollout state, and rollback. Use assist mode only for repetitive harness or report scaffolding after the contract is understood.

**Experiment/measurement:** build slices from sequential behavior through controlled concurrency and failure injection; report correctness evidence, sanitizer/test evidence, latency and queue measurements, visualization, surprises, and limitations. Coroutines may be a language-only boundary exercise or require a separately selected external runtime; they do not supply sockets, scheduling, or I/O by themselves.

## Official documentation anchors

- ISO C++ getting started: <https://isocpp.org/get-started>
- C++ Core Guidelines: <https://isocpp.github.io/CppCoreGuidelines/CppCoreGuidelines>
- cppreference language reference: <https://en.cppreference.com/w/cpp/language>
- cppreference object lifetime: <https://en.cppreference.com/w/cpp/language/lifetime>
- cppreference undefined behavior: <https://en.cppreference.com/w/cpp/language/ub>
- cppreference RAII: <https://en.cppreference.com/w/cpp/language/raii>
- cppreference smart pointers: <https://en.cppreference.com/w/cpp/memory>
- cppreference containers: <https://en.cppreference.com/w/cpp/container>
- cppreference algorithms: <https://en.cppreference.com/w/cpp/algorithm>
- cppreference ranges: <https://en.cppreference.com/w/cpp/ranges>
- cppreference concepts: <https://en.cppreference.com/w/cpp/language/constraints>
- cppreference exceptions and `noexcept`: <https://en.cppreference.com/w/cpp/language/exceptions>
- cppreference threads and `jthread`: <https://en.cppreference.com/w/cpp/thread>
- cppreference futures: <https://en.cppreference.com/w/cpp/thread/future>
- cppreference mutexes and condition variables: <https://en.cppreference.com/w/cpp/thread/mutex>
- cppreference atomics: <https://en.cppreference.com/w/cpp/atomic>
- cppreference coroutines: <https://en.cppreference.com/w/cpp/language/coroutines>

These sources are anchors, not proof that a learner ran an experiment or that a
particular compiler supports every feature. Feature availability and portability
must be checked for the exercise's explicitly selected standard and toolchain.

## Evidence boundaries

This document establishes the intended objective, ordering, source references,
experiments, and evidence contract. It does not prove C++ compilation,
portability, defined behavior, sanitizer coverage, race freedom, benchmark
validity, learner understanding, production capacity, security readiness,
provider reliability, or deployment safety. Local observations apply only to the
stated compiler, machine, implementation, workload, and run. No experiments,
manifests, dependencies, build configuration, external libraries, or generated
artifacts are created by this plan.