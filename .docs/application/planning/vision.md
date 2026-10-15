---
meta:
  schema: language-practice-application-vision/v1
  title: Distributed learning-systems telemetry platform
  status: draft
  repository: /home/michael/src/github/language-practice
  updated: 2026-10-14
  audience:
    - learner
    - future contributors
    - technical reviewers
    - potential employers
  purpose: >-
    Preserve the five language learning plans while giving their later stages a
    shared, cumulative destination: a distributed application that runs
    meaningful telemetry-analysis workloads through Python, Go, Rust,
    TypeScript, and C++ services.
  authority:
    learning_plans:
      - python/plans/01-infrastructure-learning-plan.md
      - go/plans/01-infrastructure-learning-plan.md
      - rust/plans/01-infrastructure-learning-plan.md
      - typescript/plans/01-infrastructure-learning-plan.md
      - cpp/plans/01-infrastructure-learning-plan.md
    shared_bootstrap_contract: BOOTSTRAP.md
    related_work_items:
      - .docs/work-items/WI-0003/README.md
      - .docs/work-items/WI-0004/README.md
  north_star:
    name: Learning Systems Telemetry Platform
    short_name: LSTP
    primary_contract: >-
      Given a bounded telemetry dataset and an analysis scenario, produce an
      explainable report about errors, latency, queueing, retries, saturation,
      anomalies, and service-level objective outcomes.
    architectural_shape: message-driven distributed application
    first_domain: synthetic service telemetry and incident/performance analysis
  evidence_policy:
    default_mode: local, bounded, deterministic, offline-first
    never_infer: >-
      A passing compile, local benchmark, successful package installation,
      generated report, or model-produced explanation does not by itself prove
      portability, production readiness, security, scalability, or learner
      understanding.
    required_labels:
      - documented fact
      - local observation
      - design decision
      - experiment result
      - unproven limitation
---

# Learning Systems Telemetry Platform

## 1. Executive vision

The repository should develop toward one coherent distributed application rather
than five unrelated collections of language exercises.

The application accepts a bounded telemetry dataset, runs an analysis workflow,
and produces an explainable report. The workflow itself becomes the shared
learning surface for the five language tracks:

```text
user/client
  -> TypeScript gateway
  -> Go workflow coordinator
  -> Python normalization and policy service
  -> C++ compute worker
  -> Rust execution ledger and state projection
  -> TypeScript report/status view
```

The message broker is an implementation boundary, not the product. Messages
must carry meaningful commands, results, state transitions, and telemetry. The
system should not become a chain of services that only deserialize, forward, and
log opaque payloads.

The product should answer a useful question:

> What happened in this workload, how did the system behave under pressure or
> failure, and what evidence supports that conclusion?

The learning plans remain the authority for language progression. This document
adds continuity between their later infrastructure stages and the final
application. It does not replace the plans with a framework-first curriculum.

## 2. Why this application is the right integration target

The five plans already converge on the same infrastructure concerns:

- request lifecycles and partial I/O;
- data and API contracts;
- errors and resource cleanup;
- asynchronous and parallel work;
- queues, workers, overload, and backpressure;
- timeouts, retries, rate limits, isolation, and fallbacks;
- caching, TTL, invalidation, and consistency;
- observability, tracing, metrics, and debugging;
- throughput, latency distributions, and tail behavior;
- progressive delivery and rollback; and
- a safe request-platform capstone.

The unified application gives each concern a durable home. A queue exercise can
later become a service queue. A retry experiment can later protect a real
message handler. A cache experiment can later serve policy data. A state-machine
exercise can later become the execution ledger. A C++ performance experiment can
later become the compute stage.

This prevents two opposite failure modes:

1. **Tutorial fragmentation:** many isolated examples that never accumulate into
   a system a reviewer can run or understand.
2. **Premature distributed complexity:** five thin services introduced before the
   learner understands the data structures, errors, ownership, and concurrency
   underneath them.

The intended progression is local intuition first, shared contracts second,
process boundaries third, real messaging fourth, and operational experiments
last.

## 3. User-visible product contract

The first product contract is a telemetry-analysis job.

### Input

A user submits a scenario and bounded telemetry records. A record may look like:

```json
{
  "timestamp": "2026-10-14T22:00:01Z",
  "service": "checkout",
  "route": "/payments",
  "status": 503,
  "duration_ms": 842,
  "request_id": "req-123"
}
```

The job also carries:

- a stable `work_id`;
- an idempotency key;
- a scenario name;
- a deadline;
- workload configuration;
- optional failure-injection settings;
- the language implementation or version under test.

### Output

The final report should contain, at minimum:

```json
{
  "work_id": "work-123",
  "scenario": "bounded-queue-overload",
  "summary": {
    "request_count": 100000,
    "error_rate": 0.037,
    "p50_latency_ms": 18,
    "p95_latency_ms": 220,
    "p99_latency_ms": 910,
    "slo_status": "violated"
  },
  "slow_routes": [],
  "error_windows": [],
  "queue_observations": [],
  "retry_observations": [],
  "anomalies": [],
  "evidence": {
    "dataset": "...",
    "run_id": "...",
    "services": [],
    "limitations": []
  }
}
```

The report is not merely a dashboard. It must preserve enough raw and summarized
evidence for the learner to explain what happened and what did not happen.

## 4. Service ownership

Each language owns meaningful work. The roles are defaults for architecture and
comparison, not permanent restrictions on what a language can express.

| Language | Service | Owns | Plan overlap |
| --- | --- | --- | --- |
| TypeScript | Gateway and report client | CLI/HTTP boundary, runtime validation, job submission, cancellation, status/report presentation | modules, erased types, runtime validation, promises, event loop, workers, API compatibility |
| Go | Coordinator | partitioning, scheduling, retries, deadlines, idempotency, fan-out/fan-in, aggregation, admission control | contexts, goroutines, channels, worker pools, backpressure, diagnostics, profiling |
| Python | Normalizer/policy service | NDJSON parsing, record validation, normalization, enrichment, policy lookup, cache behavior, dependency simulation | dynamic boundaries, exceptions, asyncio, executors, queues, caches, rate limits |
| Rust | Execution ledger | legal state transitions, append-only history, replay, stale/duplicate event rejection, status projection | ownership, enums, traits, `Result`, `Send`/`Sync`, channels, shutdown, typed APIs |
| C++ | Compute worker | histograms, percentiles, aggregation, anomaly features, bounded CPU workers, performance measurements | values, RAII, move/copy, templates, threads, `jthread`, queues, atomics, CMake, sanitizers |

### TypeScript gateway

The gateway accepts user input, validates it at runtime, assigns identifiers,
publishes `SubmitWork`, consumes status events, and renders the report. It must
make the difference between TypeScript’s static model and Node’s runtime
behavior visible.

### Go coordinator

The coordinator owns the workflow and retry budget. It partitions telemetry,
dispatches work, propagates deadlines, applies bounded concurrency, aggregates
partial results, and emits state transitions.

### Python normalizer/policy service

The Python service parses and validates telemetry, classifies records, removes or
reports duplicates, enriches records, and optionally consults a slow/failing
policy dependency through a cache and fallback boundary.

### Rust ledger

The ledger validates transitions such as `Submitted -> Accepted -> Running ->
Completed`, rejects stale or illegal transitions, records an append-only history,
and rebuilds a status projection by replaying events.

### C++ compute worker

The C++ worker consumes normalized batches and calculates counts, error rates,
latency distributions, top slow routes, rolling-window aggregates, and anomaly
features. It should own a bounded worker pool and expose measurable CPU and
queue behavior.

## 5. Message contract

The first shared contract should be language-neutral JSON with explicit schemas.
The protocol is more important than any particular broker client library.

### Envelope

```json
{
  "message_id": "msg-456",
  "message_type": "ComputeRequested",
  "schema_version": 1,
  "producer": "go-coordinator",
  "occurred_at": "2026-10-14T22:00:00Z",
  "correlation_id": "work-123",
  "causation_id": "msg-455",
  "idempotency_key": "client-operation-123",
  "attempt": 1,
  "deadline": "2026-10-14T22:00:10Z",
  "traceparent": "...",
  "payload": {}
}
```

Required meanings:

- `message_id` identifies one message instance;
- `message_type` identifies a command or event;
- `schema_version` makes compatibility explicit;
- `correlation_id` connects one end-to-end job;
- `causation_id` explains why a message exists;
- `idempotency_key` protects effects from duplicate delivery;
- `attempt` makes retries observable;
- `deadline` bounds work;
- `traceparent` connects the distributed diagnostic view.

### Initial message catalog

Commands:

- `SubmitWork`;
- `PolicyCheckRequested`;
- `NormalizeBatchRequested`;
- `ComputeRequested`;
- `CancelWork`;
- `RecordTransition`.

Events:

- `WorkAccepted`;
- `WorkRejected`;
- `WorkStarted`;
- `StepStarted`;
- `StepCompleted`;
- `StepFailed`;
- `RetryScheduled`;
- `WorkTimedOut`;
- `WorkCancelled`;
- `WorkCompleted`;
- `WorkDeadLettered`.

The contract should eventually live under:

```text
contracts/
├── envelope.schema.json
├── messages/
├── examples/
├── invalid/
└── compatibility/
```

Do not generate opaque language bindings before the learner has hand-written at
least one boundary mapping in each language.

## 6. Delivery and state semantics

### Initial delivery guarantee

Use at-least-once delivery. Do not promise exactly-once processing.

The system must tolerate:

- duplicate messages;
- redelivery after a consumer crash;
- out-of-order events outside an explicitly scoped ordering key;
- transient and permanent failures;
- poison messages;
- cancellation racing with completion;
- shutdown while work is queued or active.

Handlers must be idempotent wherever an effect may repeat. Acknowledgement
timing must be visible in experiments.

### State machine

The shared logical workflow states are:

```text
Submitted
  -> Accepted
  -> Normalizing
  -> Normalized
  -> Computing
  -> Aggregating
  -> Completed

Any active state may also reach:
  Failed, TimedOut, Cancelled, or DeadLettered
```

Rust owns the authoritative transition rules, but every service should be able
to represent the state relevant to its boundary.

### Data ownership

No service reaches into another service’s private storage.

| Data | Owner |
| --- | --- |
| request/read view | TypeScript gateway |
| workflow state/retry budget | Go coordinator |
| policy/enrichment cache | Python service |
| transition history/projection | Rust ledger |
| compute-local configuration/results | C++ worker |
| delivery/redelivery/dead-letter state | broker |

## 7. Scenario catalog

Every scenario should have a product outcome, a learning objective, a bounded
workload, observable signals, and a reflection prompt.

### `validate-only`

Submit malformed, incomplete, duplicated, and unknown-version telemetry. Produce
an accepted/rejected report with classified errors.

Teaches parsing, contracts, runtime validation, state transitions, and tests.

### `enrich-and-score`

Normalize records in Python, apply policy, compute aggregates in C++, and record
the workflow in Rust.

Teaches the first complete cross-language vertical slice.

### `bounded-queue-overload`

Increase producer rate above service capacity. Observe queue growth, admission
rejection, latency, recovery, and shutdown.

Teaches bounded queues, backpressure, overload, and capacity limits.

### `timeout-and-retry`

Inject slow, transient, and permanent policy failures. Observe timeout budgets,
retry classification, backoff, fallback, and retry amplification.

Teaches reliability guardrails and failure budgets.

### `duplicate-delivery`

Deliver the same compute command more than once. Confirm idempotency and inspect
the ledger history.

Teaches at-least-once semantics, duplicate effects, acknowledgements, and state
integrity.

### `fanout-and-aggregate`

Partition one dataset into many batches, fail one partition, and compare fail-fast,
retry-only-failed-partition, and degraded-result policies.

Teaches ownership, aggregation, partial failure, cancellation, and consistency.

### `cache-freshness`

Vary policy TTL, cache hits, invalidation, refresh delay, stale reads, and refresh
failure.

Teaches freshness contracts, cache stampede behavior, consistency, and source/cache
disagreement.

### `cancel-and-shutdown`

Cancel work before admission, while queued, during Python processing, during C++
computation, and during retry backoff.

Teaches cooperative cancellation, resource cleanup, shutdown ordering, and races
between cancellation and completion.

### `rollout-and-rollback`

Route a controlled percentage of work to a second compute algorithm. Halt on an
error or latency threshold and restore the previous version.

Teaches compatibility, rollout gates, verification, rollback, and operational
decision-making.

## 8. Layered implementation path

The application should grow through five transport/complexity tiers.

### Tier 0 — local pure behavior

No networking or broker. Handwrite validators, state transitions, queue models,
retry functions, cache policies, and aggregation functions.

Assistance mode: **handwrite**.

### Tier 1 — deterministic fake dependency

Use a fake dependency that can return success, delay, transient error, permanent
error, malformed output, or duplicate effects.

Assistance mode: **handwrite**, then **collaborate**.

### Tier 2 — in-memory fake transport

Implement a small deterministic transport with publish, subscribe, acknowledge,
reject, duplicate, reorder, delay, capacity, and shutdown controls.

Assistance mode: **collaborate**. The learner writes handlers and failure policy;
repetitive harness plumbing may use **assist** after the semantics are understood.

### Tier 3 — local process boundaries

Run selected services as separate local processes. Study framing, partial I/O,
startup, exit codes, disconnects, timeouts, and cleanup.

Assistance mode: **collaborate**.

### Tier 4 — real broker

Only after fake transport semantics are understood, introduce RabbitMQ, NATS, or
another explicitly selected broker. Study acknowledgements, prefetch, redelivery,
dead lettering, consumer concurrency, and broker restart.

Assistance mode: **assist** for infrastructure setup, **handwrite/collaborate**
for handlers and failure behavior.

### Tier 5 — integrated operational experiment

Run all services with controlled workloads, failure injection, observability,
performance measurement, and rollout/rollback simulation.

Assistance mode: **collaborate** for architecture and interpretation.

## 9. Learning-plan mapping

The existing plans are not discarded. Their stages map to platform work as follows.

| Plan band | Shared platform contribution | Learning mode |
| --- | --- | --- |
| 00–03 foundations | record values, request models, result types, state, validators | handwrite |
| 04–07 boundaries | message parsing, errors, serialization, fixtures, local handler APIs | handwrite/collaborate |
| 08 request lifecycle | framing, partial input, timeout, malformed message, cleanup | handwrite/collaborate |
| 09–10 concurrency | sequential vs unbounded vs bounded processing | handwrite |
| 11 queues | shared logical queue states and shutdown behavior | handwrite/collaborate |
| 12 reliability | timeout, retry, idempotency, rate limits, fallback, dead letter | collaborate |
| 13 caching | policy/enrichment/status caches and freshness contracts | collaborate |
| 14 contracts | schemas, versioning, compatibility, replay, unknown fields | collaborate |
| 15 libraries | language-local reusable reliability primitives | handwrite/collaborate |
| 16 observability | shared event vocabulary and distributed lifecycle reconstruction | assist/collaborate |
| 17 performance | fixed workload, queue wait, service time, tail latency | collaborate |
| 18 rollout | versioned compute/policy implementation and rollback state machine | collaborate |
| 19 capstone | complete telemetry-analysis platform and operational report | collaborate/assist |

### Language-specific capstone contributions

#### Python

The Python capstone owns record normalization, runtime validation, enrichment,
fake dependencies, caching, and async I/O. It should make exceptions, cleanup,
event-loop behavior, queues, and dynamic-boundary validation observable.

#### Go

The Go capstone owns partition scheduling, workflow orchestration, deadlines,
retries, idempotency, bounded worker pools, fan-out/fan-in, and aggregation. It
should make goroutine ownership, context cancellation, queue capacity, and race
diagnostics observable.

#### Rust

The Rust capstone owns legal state transitions, event replay, stale/duplicate
event rejection, and an execution projection. It should make ownership, enums,
typed errors, channels, shutdown, and `Send`/`Sync` boundaries observable.

#### TypeScript

The TypeScript capstone owns user interaction, runtime validation, request/status
models, cancellation, module boundaries, and report rendering. It should make
the difference between static types, emitted JavaScript, Node module loading, and
event-loop scheduling observable.

#### C++

The C++ capstone owns CPU-heavy aggregation, percentile/histogram computation,
bounded compute workers, and performance comparison. It should make RAII,
value/move behavior, worker ownership, queues, condition variables, atomics,
allocation, locality, and compiler/toolchain effects observable.

## 10. Canonical telemetry

Every service emits bounded structured events with common semantic fields:

```json
{
  "event_name": "work.completed",
  "event_version": 1,
  "occurred_at": "2026-10-14T22:00:01.123Z",
  "service": "go-coordinator",
  "language": "go",
  "scenario_id": "bounded-queue-overload",
  "concept_ids": ["GO-11"],
  "work_id": "work-123",
  "message_id": "msg-456",
  "correlation_id": "corr-123",
  "causation_id": "msg-455",
  "attempt": 1,
  "outcome": "success",
  "duration_ms": 42,
  "attributes": {
    "queue_depth": 12,
    "queue_wait_ms": 8,
    "worker_id": "worker-2"
  }
}
```

The analyzer must preserve raw observations and distinguish:

- event facts;
- metrics;
- diagnostic logs;
- trace relationships;
- derived interpretations;
- unproven conclusions.

Shared fields should include `work_id`, `message_id`, `correlation_id`,
`causation_id`, `attempt`, `service`, `queue`, `message_type`, `outcome`,
`duration_ms`, `deadline`, and `error_class` where relevant.

## 11. Test and experiment architecture

### Domain tests

Test pure policies and invariants:

- legal state transitions;
- malformed input;
- retry classification;
- idempotency;
- queue admission;
- cache freshness;
- aggregate correctness;
- rollout decisions.

### Contract fixture tests

All language implementations consume common fixtures:

- valid request;
- malformed request;
- unknown schema version;
- expired deadline;
- duplicate command;
- partial result;
- retryable failure;
- permanent failure.

### Fake-transport tests

Exercise duplicate delivery, reordering, delay, acknowledgement failure,
requeue, dead lettering, capacity, and shutdown deterministically.

### Process integration tests

Run selected services as local processes and validate startup, framing,
disconnect, timeout, malformed messages, and cleanup.

### Broker tests

Only after the fake transport is understood, validate acknowledgement,
redelivery, prefetch, dead lettering, broker restart, and consumer restart.

### Performance experiments

Use fixed workloads and vary one important variable at a time. Record throughput,
queue wait, service time, p50/p95/p99 latency, retries, failures, worker count,
payload size, cache hit rate, machine, toolchain, and workload limitations.

## 12. Evidence and learning contract

Every exercise must record:

- concept and prerequisites;
- what the learner writes or inspects;
- assistance mode;
- observable behavior and validation command;
- edge cases and failure modes;
- reflection question;
- evidence tier;
- what remains unproven.

The project uses four evidence tiers:

1. **Static:** source shape, schemas, syntax, configuration, and ownership.
2. **Local behavioral:** deterministic tests or simulations for one implementation.
3. **Host/project:** a user-run service, toolchain, project, or broker on one host.
4. **Learning/operational:** learner explanation, tradeoff reasoning, and bounded
   operational interpretation.

No lower tier substitutes for a higher tier.

## 13. Repository direction

The current repository contains language plans and bootstrap/toolchain material.
The application should be added incrementally rather than as one large rewrite:

```text
contracts/
├── envelope.schema.json
├── messages/
├── examples/
├── invalid/
└── compatibility/

scenarios/
├── validate-only/
├── enrich-and-score/
├── bounded-queue-overload/
├── timeout-and-retry/
├── duplicate-delivery/
├── fanout-and-aggregate/
├── cache-freshness/
├── cancel-and-shutdown/
└── rollout-and-rollback/

python/app/
go/app/
rust/app/
typescript/app/
cpp/app/

docs/application/
├── planning/
├── contracts/
├── experiments/
└── operations/
```

Do not create empty trees for future work. Add only the files needed for the
current vertical slice or experiment.

## 14. Staged delivery plan

### Phase A — architecture and contract

Define the product contract, service ownership, message envelope, lifecycle
states, telemetry fields, failure semantics, and first scenario.

### Phase B — first useful vertical slice

Implement:

```text
TypeScript gateway
  -> Go coordinator
  -> Python normalization
  -> Go aggregation
  -> TypeScript report
```

Use fixtures and an in-memory fake transport first. The first slice should
produce a real telemetry-analysis report, not only a health check.

### Phase C — ledger and compute

Add Rust state validation/replay and C++ aggregation/percentile computation.
Compare the C++ result with a reference implementation and preserve the
comparison report.

### Phase D — reliability and real transport

Add bounded queues, retries, idempotency, cancellation, dead-letter behavior,
and the selected broker after fake-transport semantics are understood.

### Phase E — observability and performance

Reconstruct distributed lifecycles from events, run fixed workloads, measure
queue wait/service time/tail latency, and diagnose injected failures.

### Phase F — rollout and capstone

Run two algorithm versions, simulate staged rollout and rollback, then publish
the complete architecture, experiment, evidence, and limitation report.

## 15. Open decisions

These decisions should be made in a dedicated architecture work item rather
than smuggled into an individual language exercise:

- final product naming;
- RabbitMQ versus NATS or another broker;
- JSON-only versus a later binary protocol;
- schema registry/versioning approach;
- persistence strategy for the Rust ledger;
- local process orchestration method;
- telemetry storage and report format;
- whether the first user interface is CLI, HTTP, or both;
- exact SLO/report vocabulary;
- C++ percentile/anomaly algorithm policy;
- whether a second domain scenario is needed after telemetry analysis.

## 16. Explicit non-goals

- Do not begin with five framework-heavy production services.
- Do not require cloud deployment, Kubernetes, managed brokers, or production
  credentials.
- Do not make the broker the learning objective before queue semantics are
  understood in memory.
- Do not force one logging, async, serialization, or testing library across all
  languages.
- Do not treat generated code as learner understanding.
- Do not claim that the final local platform is production-ready.
- Do not replace the language plans with a distributed-systems-only curriculum.

## 17. Immediate next actions

1. Create a dedicated architecture work item for the Learning Systems Telemetry
   Platform.
2. Define the first `validate-only` and `enrich-and-score` scenarios.
3. Write the envelope and first message schemas plus valid/invalid fixtures.
4. Define the deterministic fake transport and its duplicate/reorder/delay/
   capacity/shutdown controls.
5. Implement the first vertical slice: TypeScript → Go → Python → Go →
   TypeScript.
6. Add one learning contract and one experiment report for every new platform
   capability.
7. Add Rust and C++ only after the first slice has meaningful input, output,
   tests, and observability.

## 18. Decision memory

The durable decisions captured by this document are:

- the project has one distributed telemetry-analysis destination;
- the five languages own different meaningful service capabilities;
- messages and telemetry are shared contracts, implementations remain idiomatic;
- at-least-once delivery and idempotency are explicit learning concepts;
- local fake transports precede real broker infrastructure;
- every scenario must perform useful work and produce an explainable result;
- the original learning plans remain authoritative and are woven into the
  platform by stage, not discarded;
- experiments must preserve raw evidence and limitations;
- early exercises stay small and hand-written;
- generated plumbing is allowed only after the learner understands the contract;
- production claims, portability claims, and learner outcomes require separate
  evidence.

## 19. Final memory model

### Current repository state

- Five language-specific infrastructure learning plans already exist.
- Each plan has a 20-stage progression and a request-platform capstone.
- The plans intentionally leave broker, build, serialization, async-runtime, and
  dependency choices open.
- Bootstrap scripts and toolchain documentation are complete, but application code
  has not yet been added.

### Intended architecture

- TypeScript: user boundary and report client.
- Go: orchestration, scheduling, retries, and aggregation.
- Python: normalization, enrichment, policy, and cache experiments.
- Rust: typed execution ledger, replay, and state projection.
- C++: CPU-bound telemetry aggregation and performance experiments.

### Intended first useful outcome

Given a bounded telemetry dataset, the system produces a report with normalized
records, error rates, latency distributions, queue observations, retry outcomes,
and explicit limitations.

### Safety boundaries

- local and offline-first;
- bounded workloads;
- fake dependencies before live infrastructure;
- no production data or credentials;
- no deployment or cloud spending without separate authorization;
- no hidden package, broker, or shell-profile mutations.

### Unknowns

- broker selection;
- schema/tooling format;
- persistence choice;
- exact first vertical-slice repository layout;
- supported runtime versions;
- operational report format;
- whether and when to add a second domain scenario.