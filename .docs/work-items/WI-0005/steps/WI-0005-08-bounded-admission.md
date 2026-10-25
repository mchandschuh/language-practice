---
schema: language-practice-step/v1
id: WI-0005
step: "08"
status: implemented-unproven
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

### Design checkpoint

`EvaluatorRuntime` uses one worker and a caller-configured positive admission
capacity. The capacity is the total number of admitted but unfinished root
requests, including the request currently executing; it is not merely the
buffered queue length. A capacity-sized token channel accounts for admission,
while the request queue is buffered to the same size. The worker consumes a
request from the queue without releasing its token, and releases that token
only after terminal result delivery or explicit shutdown draining.

Admission validates the caller context and expression, waits for an available
token, and then serializes the closed-state check and queue send under the
evaluator mutex. A full bound therefore blocks with context-aware backpressure;
caller cancellation returns the caller context error and does not create an
accepted result channel. Shutdown rejects new admission with
`ErrEvaluatorClosed`, drains already queued requests with that explicit error,
and waits for the worker. Accepted requests are delivered at most once and are
never silently dropped. FIFO fairness is not claimed: the contract guarantees
bounded admission and cancellation/shutdown interruption only.

Nested expression jobs remain synchronous within their already-admitted parent
request. They do not acquire additional root admission tokens, so recursive
child evaluation cannot bypass the public request bound. The existing AST depth
limit remains the separate protection against excessively deep trees. Local
backpressure observations are blocked admission, cancellation of a blocked
submit, and subsequent admission after a token is released; no production
metrics or fairness guarantees are implied. The overload failure budget is one
explicit error per rejected request and zero leaked admission goroutines.

### Implementation summary

Replaced the previous queue-space notification loop with token-based admission
accounting in `lifecycle.go`. Added deterministic package tests that hold the
worker behind channel barriers, fill the exact total-admission bound, verify a
third submission blocks until a token is released, verify cancellation releases
a blocked submit, and verify nested evaluation completes with capacity one.
Existing close/drain and result-cardinality tests remain in place.

### Validation and review

- `cd go && gofmt -w basic/calc/math/lifecycle.go basic/calc/math/lifecycle_test.go`: passed.
- `cd go && go test ./basic/calc/math -count=1`: passed.
- `cd go && go test ./basic/calc/math -race -count=1`: passed.
- `cd go && go vet ./basic/calc/math`: passed.
- Focused admission/nested/close tests with `-count=20`: passed.
- Package race tests with `-count=5`: passed.
- `git diff --check`: passed.
- Review confirmed token accounting prevents accepted work from exceeding the
  configured total bound, blocked admission is context-aware, nested evaluation
  reuses the parent admission, and queued shutdown work receives an explicit
  error. FIFO fairness is intentionally not claimed.
- `cd go && go test ./... -count=1` and `cd go && go vet ./...` remain blocked by
  the pre-existing out-of-scope `basic/calc/main.go:24:24` call to
  `parsed.Evaluate`; this step does not modify that file.

Evidence tier: focused deterministic tests, repeated package tests, package race
tests, package vet, formatting, and diff checks. This proves the tested local
admission and cancellation behavior but not repository-wide race freedom,
universal goroutine leak freedom, fairness, or performance. Next action is
advisor review and step 09 shutdown work.