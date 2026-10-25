---
schema: language-practice-step/v1
id: WI-0005
step: "09"
status: implemented-unproven
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

### Design checkpoint

`EvaluatorRuntime` remains a single-worker, bounded-admission runtime. The
runtime context owns evaluator shutdown; each admitted request receives a child
context derived from that runtime context, and `context.AfterFunc` links caller
cancellation to the request child. The request stores both contexts so work
cancellation and terminal-result reporting have separate ownership: dispatch
and operator workers observe the request context, while an admitted request's
capacity-one result channel receives exactly one terminal value even if the
work context has just been canceled.

Shutdown transitions the runtime to closed under the mutex, closes the
submission wake-up signal, cancels the runtime context, drains queued requests
with `ErrEvaluatorClosed`, releases each admission slot once, and waits for the
owned worker. Active work observes the derived request cancellation and returns
before `Close` returns. New submissions are rejected after the state transition;
repeated `Close` calls only wait on the same worker join and are idempotent.

Error precedence preserves concrete evaluation/operator/child errors. A
shutdown-generated cancellation is translated to `ErrEvaluatorClosed` only
when the caller itself did not cancel and no more specific error exists. Caller
cancellation remains `context.Canceled` or `context.DeadlineExceeded`. Result
delivery for admitted requests is non-blocking because result channels are
capacity one and are owned by the request; channels are not closed by workers.

### Implementation summary

Added per-request cancellation ownership and caller-cancellation linkage in
`lifecycle.go`. Shutdown now cancels active work through the request context,
drains queued results independently of the canceled work context, preserves
root causes instead of overwriting every active error, and cleans up admission
callbacks on all rejected paths. Added deterministic lifecycle coverage for
close with active and queued requests, worker join ordering, and terminal result
delivery. Existing unary/binary worker cancellation and invalid-operator tests
continue to cover child/operator abort behavior.

### Validation and review

- `cd go && gofmt -w basic/calc/math/lifecycle.go basic/calc/math/lifecycle_test.go`: passed.
- `cd go && go test ./basic/calc/math -count=1`: passed.
- `cd go && go test ./basic/calc/math -race -count=1`: passed.
- `cd go && go vet ./basic/calc/math`: passed.
- `git diff --check`: passed.
- Review confirmed runtime shutdown is cancel-first, queued requests receive an
  explicit terminal error, active requests are joined before `Close` returns,
  caller cancellation remains distinguishable, and admission callbacks/tokens
  are released on rejection, draining, and completion.
- `cd go && go test ./... -count=1` and `cd go && go vet ./...` remain blocked by
  the pre-existing out-of-scope `basic/calc/main.go:24:24` call to
  `parsed.Evaluate`; this step does not modify that file.

Evidence tier: focused deterministic lifecycle tests, package race tests,
package vet, formatting, and diff checks. This proves the exercised shutdown,
queued-drain, cancellation, and worker-join scenarios, not universal goroutine
leak freedom, arbitrary scheduling, or repository-wide race freedom. The next
action is advisor review and step 10 concurrent reuse evidence.

### Review checkpoint

- Reviewed request-context ownership, shutdown ordering, queued draining,
  terminal-result cardinality, admission-token release, error precedence, and
  repeated-close behavior.
- No findings worth addressing remain within the selected scope.
- Repository-wide test and vet failures are pre-existing and isolated to
  `go/basic/calc/main.go`; they are not concealed by this step and remain an
  explicit evidence limitation.