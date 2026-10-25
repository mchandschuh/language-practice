---
schema: language-practice-step/v1
id: WI-0005
step: "07"
status: implemented-unproven
skill: language-practice-reliability-concurrency
depends_on: ["06"]
target: go/basic/calc/math/
---

# WI-0005 step 07 — Dispatch binary operators

## Objective

Decompose binary evaluation into dedicated add, subtract, multiply, and divide
operator workers while preserving operand ordering, short-circuit error
delivery, request state, and nested expression correctness.

## Commit subject

```text
WI-0005: 07: feat(calc): dispatch binary operators
```

## Scope and ownership

- Own only `go/basic/calc/math/` and this step's handoff update.
- Add binary-expression orchestration and typed workers for all four operators.
- Evaluate both child expressions according to the documented ordering and
  cancellation contract before dispatching the operator job.
- Preserve left/right ordering for subtraction and division.
- Define divide-by-zero behavior rather than relying on an accidental float
  result if the language contract requires an error.
- Do not tune queue capacity or add benchmark claims here.

## Design checkpoint

Record operand scheduling, result joining, operator job ownership, division/error
semantics, and the behavior when one child fails while the other is still in
flight. Identify how the implementation avoids leaking the other child.

## Implementation checkpoint

Implement binary orchestration and operator workers with typed jobs and bounded
result channels. Preserve step 06 unary behavior and the sequential baseline.

## Test and evidence checkpoint

Add deterministic tests proving:

- all binary operators produce baseline-equivalent values;
- subtraction and division preserve left/right order;
- precedence and nested binary trees complete;
- child and operator errors reach the caller without deadlock;
- cancellation during either child does not strand the sibling; and
- operator workers stop under the documented lifecycle path.

Run focused binary tests, full suite, vet, formatting, and the existing CLI
smoke test if it is still within the package contract. Step 08 requires both
normal and failed binary jobs to have evidence.

## Review checkpoint

Review operand ownership, join behavior, error/cancellation fan-in, division
semantics, worker shutdown, and comparison with the sequential oracle.

## Commit checkpoint

Stage only binary-worker implementation/tests and this handoff, verify the
cached scope and whitespace, and commit with the exact declared subject.

## Prerequisite and next gate

Step 06 must be proven and reviewed. Step 08 requires all binary operator and
nested-expression evidence to pass.

## Agent handoff notes

Reserved for binary design, implementation, tests, error/cancellation evidence,
review, limitations, and next action.

### Design checkpoint

Binary children remain synchronously scheduled in left-to-right order. The
binary handler owns both child result channels and does not schedule the right
child until the left child has delivered a successful result. Therefore a
failed or cancelled child has no concurrently running sibling to leak; the
handler returns the terminal child error through its parent result channel.
After both values are joined, the handler creates a typed `BinaryOperatorJob`
containing the request context, operator, ordered left/right values, and
request-local environment. A one-shot worker receives that job through a
capacity-one channel, sends one terminal result through a capacity-one channel,
and signals completion through a bounded `done` channel. Context cancellation
interrupts job delivery and result delivery, and the caller waits for worker
completion before returning. Workers never close caller-owned channels or
mutate the AST.

The dispatch table maps `Add`, `Subtract`, `Multiply`, and `Divide` to separate
small worker functions. Division preserves the existing sequential evaluator's
IEEE-754 behavior: a zero divisor produces the corresponding floating-point
value rather than an accidental channel-specific error. This keeps the channel
runtime equivalent to the sequential oracle.

### Implementation summary

Added `binary_workers.go` with the four typed worker paths, dispatch table,
context-aware one-shot lifecycle, worker-name observability, and ordered
arithmetic. Replaced the inline binary operator switch in `dispatch.go` with
`BinaryOperatorJob` dispatch while retaining synchronous recursive child
evaluation and bounded result delivery.

Added deterministic tests for all four workers and their values, subtraction and
division operand order, nested/precedence oracle equivalence, child and invalid
operator errors, cancellation, and worker completion. Existing single-worker
nested dispatch and evaluator shutdown tests continue to cover the lifecycle
integration path.

### Validation and review

- `cd go && gofmt -w basic/calc/math/binary_workers.go basic/calc/math/dispatch.go basic/calc/math/dispatch_test.go`: passed.
- `cd go && go test ./basic/calc/math -count=1`: passed.
- `cd go && go test ./basic/calc/math -race -count=1`: passed.
- `cd go && go vet ./basic/calc/math`: passed.
- `git diff --check`: passed.
- Review confirmed all four operators use typed workers, left/right ordering is
  preserved, child errors short-circuit before operator dispatch, cancellation
  has bounded worker completion, and no queue capacity or benchmark changes were
  introduced.
- `cd go && go test ./... -count=1` and `go vet ./...` remain expected to be
  blocked by the pre-existing out-of-scope `basic/calc/main.go:24:24` call to
  `parsed.Evaluate`; this step does not modify that file.

Evidence tier: focused binary tests, package race test, package vet, formatting,
and diff checks. This does not prove repository-wide race freedom, universal
goroutine leak freedom, or performance. The next action is advisor review and
step 08 bounded-admission work.