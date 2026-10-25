---
schema: language-practice-step/v1
id: WI-0005
step: "04"
status: planned
skill: language-practice-reliability-concurrency
depends_on: ["03"]
target: go/basic/calc/math/
---

# WI-0005 step 04 — Add evaluator lifecycle

## Objective

Introduce the `Evaluator` runtime boundary with a root submission queue, explicit
construction, request admission, cancellation propagation, and shutdown. The
runtime may initially delegate computation to the sequential baseline; lifecycle
correctness comes before type/operator worker decomposition.

## Commit subject

```text
WI-0005: 04: feat(calc): add evaluator lifecycle
```

## Scope and ownership

- Own only `go/basic/calc/math/` and this step's handoff update.
- Add `NewEvaluator`, public evaluation submission, and `Close` or equivalent
  lifecycle methods.
- Own worker goroutines and wait for them during shutdown.
- Define queue capacity, admission behavior, repeated close behavior, and whether
  evaluation after close returns a sentinel/error.
- Do not add type-specific or operator-specific worker channels yet.

## Design checkpoint

Record:

- evaluator-owned context and cancellation source;
- root queue capacity and admission/backpressure behavior;
- startup and shutdown order;
- what happens to queued, running, and newly submitted jobs during close; and
- the timeout policy used by lifecycle tests.

The design must prohibit sending to a queue after shutdown and must not depend on
sleep-based timing for correctness.

## Implementation checkpoint

Implement the evaluator constructor, root admission path, lifecycle state, and
worker join behavior. Ensure all blocking operations select on context or
shutdown. Make `Close` safe for the documented repeated-call behavior.

## Test and evidence checkpoint

Add bounded tests proving:

- one request completes through the evaluator boundary;
- multiple requests complete without losing results;
- admission obeys the documented queue capacity;
- cancellation unblocks a waiting request;
- close returns and all owned workers stop;
- submit-after-close returns the documented error; and
- repeated close does not panic or deadlock.

Use explicit timeouts only as failure bounds. Capture a before/after goroutine
observation with a documented tolerance, but do not claim universal leak freedom.
Run focused tests, `go test ./...`, `go vet`, and formatting. This evidence
unlocks type dispatch only if normal completion and shutdown both pass.

## Review checkpoint

Review lifecycle ownership, queue behavior, cancellation races, close ordering,
worker joins, and tests for both success and abort. A worker that can outlive
`Close` is a blocker.

## Commit checkpoint

Stage only lifecycle implementation/tests and this handoff, verify cached
whitespace and scope, and commit with the exact declared subject.

## Prerequisite and next gate

Step 03 must be proven and reviewed. Step 05 cannot begin with a runtime that has
an untested shutdown or cancellation path.

## Agent handoff notes

Reserved for lifecycle design, implementation, test output, worker evidence,
review, limitations, and next action.