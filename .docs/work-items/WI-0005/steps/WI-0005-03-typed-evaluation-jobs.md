---
schema: language-practice-step/v1
id: WI-0005
step: "03"
status: planned
skill: language-practice-go
depends_on: ["02"]
target: go/basic/calc/math/
---

# WI-0005 step 03 — Define typed evaluation jobs

## Objective

Define the message protocol used by future evaluator workers: typed jobs,
explicit result/error values, request state, context propagation, and channel
direction. This step may add protocol helpers but must not start the full worker
runtime.

## Commit subject

```text
WI-0005: 03: feat(calc): define typed evaluation jobs
```

## Scope and ownership

- Own only `go/basic/calc/math/` and this step's handoff update.
- Define root expression jobs and the distinct unary, binary, number, identifier,
  and operator job shapes needed by later steps.
- Define a result type that can carry a value, final request state where needed,
  and an error without logging-only failure propagation.
- Use `context.Context` in jobs that may block or be cancelled.
- Use directional channels in worker signatures and document who creates, sends,
  receives, and closes each channel.
- Do not add long-lived goroutines, worker pools, or public lifecycle methods.

## Design checkpoint

Record the protocol table:

| Message | Producer | Consumer | Result owner | Close owner |
| --- | --- | --- | --- | --- |
| expression job | public evaluator | type dispatcher | submitting request | creator/no close for one-shot |
| unary job | unary dispatcher | unary operator worker | parent evaluator | creator |
| binary job | binary dispatcher | binary operator worker | parent evaluator | creator |

The design must explain why a result channel carries errors, why one-shot result
channels are buffered, and why workers never close caller-owned channels.

## Implementation checkpoint

Add the protocol types and small validation/send helpers without changing the
sequential semantics. Keep result delivery cancellation-aware and prevent a
worker from blocking forever when the request context is done.

## Test and evidence checkpoint

Add tests proving:

- valid jobs preserve expression, inputs, context, and result destination;
- successful results carry values and failed results carry errors;
- a cancelled context prevents or aborts delivery according to the contract;
- a one-shot buffered result can receive exactly one result without a receiver
  race; and
- protocol validation rejects nil/invalid inputs in the documented way.

Run focused package tests, the full Go suite, `go vet`, and formatting. The
evidence gate for step 04 is a documented ownership table plus tests for success,
error, and cancellation delivery; no worker-lifecycle claim is allowed yet.

## Review checkpoint

Review channel directions, close ownership, result cardinality, context use,
error preservation, and whether protocol helpers accidentally introduce hidden
goroutines or unbounded queues.

## Commit checkpoint

Stage only protocol implementation/tests and this handoff, inspect the cached
diff, and commit with the exact declared subject.

## Prerequisite and next gate

Step 02 must be proven and reviewed. Step 04 requires the typed protocol tests
and ownership table to pass.

## Agent handoff notes

Reserved for protocol design, implementation, exact test commands/results,
review, evidence tier, limitations, and next action.