---
schema: language-practice-step/v1
id: WI-0005
step: "04"
status: implemented-unproven
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

### Design checkpoint

`EvaluatorRuntime` owns a derived context, cancellation function, one root worker,
and a bounded root queue. `NewEvaluator(parent, capacity)` requires a positive
capacity; `DefaultQueueCapacity` is 16 for callers that want the documented
default. `Submit(ctx, expression, input)` clones the input environment, returns a
capacity-one result channel, and blocks only while admission is full. Admission
selects on the request context and evaluator shutdown, and the lifecycle mutex
prevents a send from racing with close. `ErrEvaluatorClosed` rejects new work.

The worker delegates computation to the existing sequential evaluator. The
evaluator context is used for lifecycle cancellation, while each request context
controls result delivery. `Close` is idempotent: it marks the runtime closed,
cancels the owned context, drains queued jobs with `ErrEvaluatorClosed`, and
waits on the worker `WaitGroup`. A request already running finishes the bounded
sequential calculation and is reported as closed if shutdown began; queued work
is rejected during draining; new work is rejected immediately. No queue is sent
to after shutdown state is set, and the queue is never closed because the worker
owns its lifecycle through the shutdown signal.

Lifecycle tests use one-second timeout bounds only to diagnose hangs. Normal
ordering uses result channels and explicit `Close`; goroutine evidence compares
before/after observations with a bounded tolerance and is not a universal leak
claim.

### Implementation summary

Added `lifecycle.go` with `EvaluatorRuntime`, constructor, bounded root
admission, cancellation-aware submission, idempotent close, queue draining, and
worker joining. Added `lifecycle_test.go` covering reusable AST evaluation,
multiple results, bounded admission/cancellation, submit-after-close,
repeated-close safety, close cancellation, and worker termination observation.
The runtime intentionally does not add type-specific or operator-specific worker
channels; those remain later-step work.

### Validation and review

- `cd go && gofmt -w basic/calc/math/lifecycle.go basic/calc/math/lifecycle_test.go`: passed.
- `cd go && go test ./basic/calc/math -count=1`: passed.
- `cd go && go vet ./basic/calc/math`: passed.
- `git diff --check`: passed.
- `cd go && go test ./... -count=1` remains blocked by the pre-existing,
  out-of-scope `basic/calc/main.go:24:24` call to `parsed.Evaluate`.
- `cd go && go vet ./...` reports the same pre-existing out-of-scope caller.
- Review checked lifecycle ownership, bounded queue admission, cancellation
  selection, close ordering, repeated close, queue draining, result ownership,
  worker joining, and selected-path scope. No type/operator worker decomposition
  was introduced.

Evidence tier: focused lifecycle implementation, tests, and package vet; the
repository-wide evidence remains limited by the known caller failure. The next
action is step 05 only after this lifecycle baseline is reviewed and accepted.