---
schema: language-practice-step/v1
id: WI-0005
step: "06"
status: implemented-unproven
skill: language-practice-reliability-concurrency
depends_on: ["05"]
target: go/basic/calc/math/
---

# WI-0005 step 06 — Dispatch unary operators

## Objective

Decompose unary evaluation into dedicated operator workers for positive, negate,
prefix decrement, postfix decrement, prefix increment, and postfix increment.
Preserve the type dispatcher and request-local state while making operator
selection observable and testable.

## Commit subject

```text
WI-0005: 06: feat(calc): dispatch unary operators
```

## Scope and ownership

- Own only `go/basic/calc/math/` and this step's handoff update.
- Add unary-expression orchestration and one dedicated worker path per supported
  unary operator, or a clearly documented equivalent dispatch table.
- Keep operator workers small: receive a typed value/state job, compute, and
  return one result/error.
- Preserve true prefix/postfix semantics for identifiers and reject unsupported
  non-assignable mutations according to step 02.
- Do not add binary operator workers or queue tuning here.

## Design checkpoint

Record the operator dispatch table, state mutation ownership, result ordering,
and why prefix/postfix workers cannot mutate shared AST or shared input maps.
Explain whether workers are long-lived or one-shot and how they stop.

## Implementation checkpoint

Implement the unary dispatcher and dedicated workers using directional channels,
context-aware sends/receives, and typed results. Preserve all step 05 nested
dispatch behavior.

## Test and evidence checkpoint

Add deterministic tests proving:

- every unary operator selects the intended worker and returns the expected value;
- prefix and postfix identifier updates differ in returned value but share the
  documented final state;
- unary child errors and cancellation propagate;
- literal/non-assignable mutation behavior remains explicit; and
- repeated and nested unary expressions complete without deadlock.

Run focused unary tests, the full suite, vet, and formatting. Include a bounded
worker shutdown assertion for the new worker paths. This evidence unlocks binary
operator decomposition but does not prove system-wide race freedom.

## Review checkpoint

Review operator coverage, state ownership, prefix/postfix semantics, channel
closure, error propagation, and worker shutdown. Reject duplicated hidden
evaluation logic that could diverge from the sequential contract.

## Commit checkpoint

Stage only unary-worker implementation/tests and this handoff, inspect the
cached diff, and commit with the exact declared subject.

## Prerequisite and next gate

Step 05 must prove type dispatch and nested orchestration. Step 07 requires all
unary worker tests and review evidence.

## Agent handoff notes

Reserved for operator design, implementation, test output, lifecycle evidence,
review, limitations, and next action.

### Design checkpoint

The unary dispatch table maps each supported operator to one small worker path:

| Operator | Worker | Result |
| --- | --- | --- |
| `Positive` | `positiveWorker` | operand unchanged |
| `Negate` | `negateWorker` | negated operand |
| `PrefixDecrement` | `prefixDecrementWorker` | update request state, return updated value |
| `PostfixDecrement` | `postfixDecrementWorker` | update request state, return previous value |
| `PrefixIncrement` | `prefixIncrementWorker` | update request state, return updated value |
| `PostfixIncrement` | `postfixIncrementWorker` | update request state, return previous value |

The unary handler owns child evaluation and ordering. It creates a typed
`UnaryOperatorJob` after the operand is available and sends it through a
capacity-one directional channel to the selected worker. Mutation ownership is
explicit: the job carries the request-local environment, target name, and an
assignability flag. Workers never mutate the AST or caller-owned input maps.
Literal mutation retains the step 01 numeric compatibility behavior; compound
or otherwise non-assignable mutation is rejected before worker dispatch.

Workers are one-shot goroutines. Each invocation owns its input/result channels,
selects sends and receives against the request context, emits at most one
terminal result, and reports completion through a bounded shutdown signal. The
unary handler waits for that signal before returning, and cancellation causes
the worker to stop without waiting for a receiver. No worker closes a caller
owned result channel.

### Implementation summary

Added `unary_workers.go` with the six explicit unary worker functions, a typed
dispatch table, context-aware one-shot worker lifecycle, and request-state-aware
mutation. Extended `UnaryOperatorJob` with target and assignability metadata.
Updated the type dispatcher so all unary arithmetic and mutation computation
flows through these workers while preserving synchronous recursive child
dispatch and the existing single-worker deadlock avoidance.

Added deterministic tests for every dispatch-table entry and result, worker
completion, cancellation, invalid operators, literal compatibility, and
non-assignable mutation. Existing nested, child-error, request-state, and
single-worker tests remain in the focused package suite.

### Validation and review

- `cd go && gofmt -w basic/calc/math/*.go`: passed.
- `cd go && go test ./basic/calc/math -count=1`: passed.
- `cd go && go test ./basic/calc/math -race -count=1`: passed.
- `cd go && go vet ./basic/calc/math`: passed.
- `git diff --check`: passed.
- `cd go && go test ./... -count=1` and `go vet ./...` remain blocked by the
  pre-existing out-of-scope `basic/calc/main.go:24:24` call to `parsed.Evaluate`.
- Review confirmed coverage of all six operators, request-local mutation,
  prefix/postfix ordering, context-aware delivery, bounded one-shot shutdown,
  preserved recursive dispatch, and no binary-worker or queue changes.

Evidence tier: focused unary tests, race test, package vet, formatting, and
bounded worker completion assertions. This does not prove system-wide race
freedom or universal goroutine leak freedom. The next action is advisor review
and step 07 binary-worker decomposition.