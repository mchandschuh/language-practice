---
schema: language-practice-step/v1
id: WI-0005
step: "02"
status: planned
skill: language-practice-go
depends_on: ["01"]
target: go/basic/calc/math/
---

# WI-0005 step 02 — Add reusable input environments

## Objective

Add identifiers and isolated per-evaluation inputs so one parsed AST can be
evaluated repeatedly with different values. Define variable mutation semantics
before any goroutine or channel implementation relies on them.

## Commit subject

```text
WI-0005: 02: feat(calc): add reusable input environments
```

## Scope and ownership

- Own only `go/basic/calc/math/` and this step's handoff update.
- Add an identifier AST node and the parser syntax required for valid names.
- Define the public input/environment type and ownership rules.
- Ensure each evaluation receives isolated state; caller-owned input must not be
  mutated unless the API explicitly returns a new state.
- Define and implement prefix/postfix increment/decrement for identifiers.
- Do not add channels, worker goroutines, evaluator lifecycle, CLI flags, or
  benchmarks.

## Design checkpoint

Record:

- identifier grammar and invalid-name errors;
- whether inputs are copied at admission or represented by an explicit owned
  environment;
- whether evaluation returns only a value or value plus final state;
- exact prefix/postfix semantics and assignment requirements; and
- how missing identifiers and non-assignable increment operands fail.

The design must distinguish immutable AST structure from mutable per-request
state and must explain why two concurrent requests cannot share a mutable map.

## Implementation checkpoint

Implement the identifier node, environment/state type, parser support, and
sequential evaluator behavior. Preserve step 01's literal semantics. Keep all
state associated with a request rather than storing it on the AST or a package
global.

## Test and evidence checkpoint

Add deterministic tests proving:

- `x * x + 2` parses and evaluates with at least two input environments;
- the same AST can be evaluated repeatedly without reparsing;
- missing identifiers produce stable errors;
- `++x` and `x++` return the documented values and final state;
- incrementing a literal or otherwise non-assignable expression follows the
  documented rejection/compatibility rule;
- caller input is unchanged when isolation is promised; and
- sequential package tests, `go vet`, and formatting pass.

The evidence gate for step 03 is an executable repeated-evaluation test and a
state-isolation test. It does not prove concurrent safety until step 10.

## Review checkpoint

Review identifier grammar, environment copying/ownership, mutation semantics,
AST reuse, missing-input errors, and tests that compare both returned values and
final state. Reject hidden global state or tests that only evaluate one input.

## Commit checkpoint

Stage only the selected math files and this handoff, validate the cached diff,
and commit with the exact declared subject. Record the evidence that unlocks
typed jobs.

## Prerequisite and next gate

Step 01 must be proven and reviewed. Step 03 requires this step's committed
environment contract and passing sequential evidence.

## Agent handoff notes

Reserved for design, implementation, test output, review findings, evidence
tier, limitations, and next action.