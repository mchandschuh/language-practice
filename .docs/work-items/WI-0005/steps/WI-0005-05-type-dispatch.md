---
schema: language-practice-step/v1
id: WI-0005
step: "05"
status: planned
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