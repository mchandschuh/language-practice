---
schema: language-practice-step/v1
id: WI-0005
step: "07"
status: planned
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