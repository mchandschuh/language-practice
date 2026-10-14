---
name: language-practice-infrastructure-design
description: Frame language exercises around reliable infrastructure concerns without skipping fundamentals. Use for services, APIs, queues, storage, caching, scheduling, deployment boundaries, configuration, and operational design.
compatibility: Designed for local designs, simulations, and bounded integration exercises; live deployment and infrastructure mutation require explicit authorization.
metadata:
  version: "1"
  scope: "infrastructure-oriented learning"
---

# Infrastructure-oriented design

Use this skill to connect language practice to the work of building dependable
systems. Prefer a small executable experiment over an architecture diagram that
has no observable behavior.

## Design sequence

1. State the user-visible contract and failure budget.
2. Identify data ownership, consistency needs, and durability requirements.
3. Define API, queue, storage, cache, worker, and scheduler boundaries only as
   needed by the exercise.
4. Identify timeouts, retries, idempotency, backpressure, rate limits, and
   graceful degradation.
5. Define logs, metrics, traces, health signals, and operator actions.
6. Model dependency failure, partial failure, restart, duplicate delivery, and
   recovery.
7. Validate the design with a bounded local simulation or test.

## Scope and safety

- Prefer fakes, local services, temporary directories, and disposable data.
- Do not deploy, apply infrastructure, publish a service, access production,
  or incur cost without explicit user scope.
- Do not put credentials, private endpoints, customer data, or copied production
  configuration into exercises or logs.
- Distinguish a design claim, a local simulation, a provider-backed plan, and a
  real operational result.

## Learning focus

In `handwrite` mode, keep the service boundary small enough that the learner can
implement the important behavior. In `assist` mode, generated boilerplate must
not obscure the protocol, failure model, or resource lifecycle being learned.

Require the learner to explain tradeoffs such as consistency versus latency,
retries versus duplicate work, caching versus staleness, and simplicity versus
operational visibility.