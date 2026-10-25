---
schema: language-practice-step/v1
id: WI-0005
step: "05"
status: implemented-unproven
skill: language-practice-go
depends_on: ["04"]
target: go/basic/calc/math/
---

# WI-0005 step 05 — Dispatch expression types

## Objective

Replace the lifecycle baseline's direct computation with type-centric dispatch
for number, identifier, unary, and binary AST nodes. Establish recursive child
job orchestration without yet splitting every operator into its own worker.

## Commit subject

```text
WI-0005: 05: feat(calc): dispatch expression types
```

## Scope and ownership

- Own only `go/basic/calc/math/` and this step's handoff update.
- Add the root type dispatcher and typed handlers for all supported AST nodes.
- Recursively submit child expressions through the evaluator boundary.
- Preserve per-request environment state and error propagation.
- Keep operator computation in the type handlers for this step; operator worker
  channels are step 06/07 work.
- Do not solve nested dispatch by silently increasing worker counts or by adding
  unbounded goroutines.

## Design checkpoint

Record the type-dispatch table, recursion strategy, child-result ownership, and
the mechanism that prevents a dispatcher from blocking while waiting for a child
job that it alone must receive. Include a bounded nesting/depth decision and a
documented failure for unsupported AST types.

## Implementation checkpoint

Implement type dispatch and recursive orchestration against the lifecycle and
typed-job contracts. Preserve the sequential evaluator's values, state, and
errors. Keep each child result one-shot and cancellation-aware.

## Test and evidence checkpoint

Add tests proving:

- each supported AST type reaches its handler;
- nested unary and binary expressions complete;
- deeply nested but bounded expressions do not deadlock;
- child errors reach the original request;
- unsupported node types return an error; and
- multiple independent requests can be in flight without cross-result mixing.

Use bounded test deadlines and an explicit regression test for the single-worker
nested-dispatch deadlock hazard. Run focused tests, full suite, vet, and format.
The next step is blocked if nested dispatch relies on a timing accident.

## Review checkpoint

Review type coverage, recursive scheduling, request-state propagation, child
channel ownership, deadlock avoidance, and error context. Confirm that the tests
would fail if a child job were routed to a blocked worker.

## Commit checkpoint

Stage only type-dispatch implementation/tests and this handoff, validate the
cached diff, and commit with the exact declared subject.

## Prerequisite and next gate

Step 04 must prove lifecycle behavior. Step 06 requires all type-dispatch tests,
including nested-expression evidence, to pass.

## Agent handoff notes

Reserved for dispatch design, implementation, tests, deadlock evidence, review,
limitations, and next action.

### Design checkpoint

The dispatcher uses this table:

| AST type | Handler | Behavior |
| --- | --- | --- |
| `*NumberExpression` | number handler | returns the literal value |
| `*IdentifierExpression` | identifier handler | reads request-local state and reports missing names |
| `*UnaryExpression` | unary handler | recursively dispatches its operand, then applies the unary operator |
| `*BinaryExpression` | binary handler | recursively dispatches left then right, then applies the binary operator |
| other/nil | unsupported handler | returns an explicit error result |

Recursive dispatch is synchronous and reentrant. A child is dispatched directly
through the typed boundary instead of being submitted to the root queue. This is
the deliberate single-worker deadlock prevention mechanism: a type handler never
waits for a child that only the blocked worker can receive. Child result channels
are capacity-one and owned by the parent handler; children send one terminal
result, never close the channel, and all delivery remains context-aware.

The dispatcher enforces a bounded nesting depth of 1024 AST edges. Expressions at
the bound are accepted and deeper expressions return a stable depth error. The
request environment is shared by recursive handlers within one request so
increment/decrement ordering and final state match the sequential evaluator;
the lifecycle boundary still clones caller input for every root request.

### Implementation summary

Added `dispatch.go` with typed number, identifier, unary, and binary handlers,
recursive child orchestration, cancellation-aware result waits, unsupported-node
errors, and the explicit depth bound. The runtime worker now enters through the
dispatcher rather than calling the sequential evaluator directly. Operator
calculation remains in the type handlers; operator worker channels remain
deferred to steps 06 and 07.

Added `dispatch_test.go` covering all supported node types, nested expressions
with a single worker, bounded and over-limit nesting, child error propagation,
unsupported nodes, request-local mutation with independent concurrent requests,
and cancellation behavior.

### Validation and review

- `cd go && gofmt -w basic/calc/math/dispatch.go basic/calc/math/dispatch_test.go basic/calc/math/lifecycle.go`: passed.
- `cd go && go test ./basic/calc/math -count=1`: passed.
- `cd go && go test ./basic/calc/math -race -count=1`: passed.
- `cd go && go vet ./basic/calc/math`: passed.
- `git diff --check`: passed.
- `cd go && go test ./... -count=1` and `go vet ./...` remain blocked by the
  pre-existing out-of-scope `basic/calc/main.go:24:24` call to `parsed.Evaluate`.
- Review confirmed that nested children are not routed to the root queue, result
  ownership remains one-shot and cancellation-aware, environments remain
  request-local, unsupported nodes and depth overflow return errors, and no
  additional goroutines or worker-count changes were introduced.

Evidence tier: focused dispatch tests, race test, package vet, and explicit
single-worker deadlock regression. Repository-wide validation remains limited by
the known caller failure. Step 06 may now introduce unary operator workers while
preserving this recursive dispatch and ownership contract.