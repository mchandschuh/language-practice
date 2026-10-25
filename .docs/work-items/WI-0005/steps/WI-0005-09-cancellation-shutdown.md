---
schema: language-practice-step/v1
id: WI-0005
step: "09"
status: planned
skill: language-practice-reliability-concurrency
depends_on: ["08"]
target: go/basic/calc/math/
---

# WI-0005 step 09 — Make evaluator shutdown reliable

## Objective

Complete context cancellation, error propagation, shutdown ordering, worker
termination, and cleanup for the channel evaluator. This is the reliability
gate before concurrency reuse and race testing.

## Commit subject

```text
WI-0005: 09: bug(calc): make evaluator shutdown reliable
```

## Scope and ownership

- Own only `go/basic/calc/math/` and this step's handoff update.
- Define cancellation for queued, running, child, and operator jobs.
- Make shutdown cancel first, reject new work, wait for owned workers, and follow
  the documented repeated-close behavior.
- Preserve root causes when returning errors; do not replace them with generic
  worker-exited messages.
- Add no new operator semantics or benchmark claims.

## Design checkpoint

Record the shutdown state machine, ownership of every goroutine and channel,
close/cancel ordering, behavior for in-flight jobs, error precedence when both
work and context fail, and bounded shutdown timeout expectations.

## Implementation checkpoint

Implement cancellation-aware worker loops, shutdown coordination, error fan-in,
and cleanup. Ensure every goroutine has a clear exit path for success, error,
queue rejection, cancellation, and close.

## Test and evidence checkpoint

Add bounded tests proving:

- cancellation of queued work;
- cancellation of running unary and binary work;
- cancellation while child jobs are in flight;
- operator error propagation with preserved context;
- close while idle, queued, and active;
- all owned workers terminate before close returns;
- repeated close and submit-after-close behavior; and
- no result channel is left waiting after abort.

Run `go test ./...`, `go vet ./...`, focused lifecycle tests, and a scoped
`go test -race` run. Use explicit timeouts as failure bounds and record that
passing tests do not prove universal leak freedom. Step 10 cannot begin without
both success and abort evidence.

## Review checkpoint

Review cancellation propagation, shutdown state transitions, error precedence,
worker join guarantees, channel closure, timeout handling, and race output.

## Commit checkpoint

Stage only reliability implementation/tests and this handoff, inspect cached
scope and whitespace, and commit with the exact declared subject.

## Prerequisite and next gate

Step 08 must prove bounded admission. Step 10 requires this step's shutdown and
race evidence to be reviewed and committed.

## Agent handoff notes

Reserved for lifecycle design, implementation, cancellation/shutdown tests,
race evidence, review, limitations, and next action.