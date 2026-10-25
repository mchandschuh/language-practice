---
schema: language-practice-step/v1
id: WI-0005
step: "06"
status: planned
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