---
schema: language-practice-step/v1
id: WI-0005
step: "08"
status: planned
skill: language-practice-reliability-concurrency
depends_on: ["07"]
target: go/basic/calc/math/
---

# WI-0005 step 08 — Bound evaluator admission

## Objective

Make queue capacity, worker count, admission, backpressure, and overload
behavior explicit. The evaluator must not accept unbounded work merely because
the public API is convenient.

## Commit subject

```text
WI-0005: 08: feat(calc): bound evaluator admission
```

## Scope and ownership

- Own only `go/basic/calc/math/` and this step's handoff update.
- Add bounded queue configuration or documented defaults.
- Define whether full queues block with context, reject immediately, or apply a
  bounded policy; return a typed or sentinel error for rejection.
- Preserve all worker and operator behavior from step 07.
- Do not add retries or silently drop accepted jobs.

## Design checkpoint

Record queue capacities, worker counts, admission order, rejection behavior,
backpressure metrics/observations available to local tests, and the failure
budget for overload. Explain how internal child jobs are accounted for so a
parent cannot bypass the bound through recursive submission.

## Implementation checkpoint

Implement bounded admission and configuration validation. Keep sends and receives
context-aware and ensure rejection does not leave a result channel waiting
forever.

## Test and evidence checkpoint

Add deterministic tests proving:

- a queue accepts work up to its documented bound;
- full admission blocks or rejects exactly as documented;
- cancellation releases a blocked admission attempt;
- rejected work produces an explicit result/error and no leaked goroutine;
- accepted work is neither duplicated nor silently lost; and
- nested evaluation cannot bypass the admission bound.

Use barriers/channels rather than sleeps. Run focused queue tests, full suite,
vet, formatting, and a bounded test that records queue/rejection outcomes. The
next step is blocked if queue capacity is accidental or rejection is silent.

## Review checkpoint

Review backpressure semantics, recursive admission, error cardinality, fairness
claims, worker sizing, and cleanup on rejected/cancelled work.

## Commit checkpoint

Stage only admission/queue implementation/tests and this handoff, validate the
cached diff, and commit with the exact declared subject.

## Prerequisite and next gate

Step 07 must prove all operator workers. Step 09 requires bounded admission and
rejection evidence before adding final shutdown handling.

## Agent handoff notes

Reserved for queue design, implementation, test output, overload evidence,
review, limitations, and next action.