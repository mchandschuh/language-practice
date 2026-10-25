---
schema: language-practice-work-item/v1
id: WI-0005
status: PLANNED
type: feature
title: Build a reusable channel-driven AST evaluator
repository: /home/michael/src/github/language-practice
branch: wi-0005-channel-evaluator
authority: local work-item records under this directory
---

# WI-0005 — Build a reusable channel-driven AST evaluator

## Authority and status

- **Raw status:** `PLANNED`
- **Derived status:** `planned`
- **Authority:** this README and the ordered `steps/` records
- **Repository:** `/home/michael/src/github/language-practice`
- **Branch:** `wi-0005-channel-evaluator`
- **Baseline:** `e6abe32` (`go: calc: Move Evaluate to MathExpression`)
- **Predecessor:** the WI-0005 vision commit `47ae32a` (`WI-0005: doc(repo): define telemetry platform vision`)
- **Root record commit:** `WI-0005: doc(repo): define channel-driven AST evaluator plan`
- **Implementation commits:** one exact subject per ordered implementation step,
  after this root record is committed

## Objective

Turn the existing calculator parser and AST into a reusable evaluator pipeline.
The same parsed expression tree must be evaluable repeatedly with different input
values. The final implementation will expose a bounded local Go evaluator whose
internal work is decomposed into typed jobs and small operator workers connected
by channels.

The work is intentionally staged. A correct sequential evaluator and explicit
expression semantics must exist before channels and goroutines are introduced.
Each later step is blocked until the previous step has a committed implementation,
passing tests, an advisor review, and evidence sufficient to support the next
design boundary.

## User-visible contract

The calculator must eventually support a flow equivalent to:

```text
parse "x * x + 2" once
evaluate the same AST with x=3  -> 11
evaluate the same AST with x=10 -> 102
```

The channel evaluator must preserve the sequential evaluator's observable
semantics for the supported expression set. It must return values and errors
through explicit result types, accept `context.Context` cancellation, apply
bounded admission/backpressure, and shut down owned goroutines without leaks in
the tested scenarios.

## Design decisions

### AST is reusable data

`MathExpression` is an AST marker interface, not the evaluation engine. AST nodes
are normalized to pointer nodes and remain immutable during evaluation. The
evaluator owns runtime state, input bindings, channels, cancellation, and worker
lifecycle. This prevents one evaluation request from contaminating another and
allows multiple requests to reuse the same tree.

The planned node set is:

- `NumberExpression`;
- `IdentifierExpression`;
- `UnaryExpression`; and
- `BinaryExpression`.

### Per-request state

Public input bindings are treated as read-only. Each evaluation request receives
an isolated state copy. Prefix and postfix increment/decrement may update that
request-local state only. They must not mutate a caller-owned map or shared AST.

The step that introduces identifiers must define and test the observable
distinction:

```text
x = 5; ++x -> value 6, final x 6
x = 5; x++ -> value 5, final x 6
```

Literal operands remain valid for the learning exercise only if their behavior
is documented and tested. The implementation must not imply that incrementing a
literal is a real mutable-variable operation.

### Typed jobs and results

Workers do not pass bare `float64` values when an error or cancellation can
occur. The design uses an explicit result shape similar to:

```go
type EvaluationResult struct {
    Value float64
    State EvaluationState
    Err   error
}
```

The exact public names may be refined by step 01, but the contract must remain
explicit. A one-shot result channel is created by the submitting owner, is
buffered sufficiently for one result, and is not closed by a worker that did not
create it.

### Type and operator dispatch

The final internal pipeline is intentionally educational:

```text
root expression job
  -> type dispatcher
     -> number / identifier evaluator
     -> unary-expression evaluator
        -> positive / negate / prefix / postfix operator evaluator
     -> binary-expression evaluator
        -> add / subtract / multiply / divide evaluator
```

The type dispatcher owns the mapping from AST node type to typed job. Unary and
binary evaluators recursively submit child expression jobs and then dispatch
operator jobs. The design must avoid a single worker blocking forever while
waiting for a nested job that only that same worker can receive.

### Ownership and lifecycle

The `Evaluator` owns its internal queues, workers, cancellation context, and
`sync.WaitGroup`. Callers submit through a public method rather than writing to
internal channels. Workers receive from channels they do not close. Shutdown
cancels first, prevents new work, drains or rejects bounded pending work
according to the documented contract, waits for workers, and is idempotent or
returns a documented repeated-close result.

Every blocking send and receive in the channel implementation must have a
cancellation or shutdown story. Queue capacities and worker counts are explicit
configuration or documented defaults, not accidental constants.

### Evidence policy

Each step has five mandatory checkpoints:

1. **Design:** record the contract, ownership, failure behavior, and API before
   implementation.
2. **Implementation:** change only the selected step's owned code and tests.
3. **Test/evidence:** run deterministic focused tests, the applicable package
   suite, and the declared diagnostics. A passing test proves only its stated
   behavior.
4. **Review:** inspect the diff, test output, scope, error paths, lifecycle, and
   evidence limits; append an advisor review to the step record.
5. **Commit:** stage only the selected target and step handoff, verify the exact
   declared subject, and commit before the next dependent step.

No step is considered proven merely because its code compiles. Concurrency steps
must include cancellation, shutdown, leak, ordering, backpressure, and race
evidence appropriate to their scope.

## Scope

- Refactor the Go calculator under `go/basic/calc/` into reusable AST data,
  input/state types, a sequential evaluator, and a channel-driven evaluator.
- Preserve parser correctness for arithmetic, precedence, parentheses, unary
  operators, and the existing increment/decrement syntax while making variable
  semantics explicit.
- Add deterministic unit, integration, race, leak, cancellation, and bounded
  benchmark evidence as the relevant step introduces each behavior.
- Keep the existing module and `init` app helpers usable; integrate the final
  evaluator with `run calc` and `test calc` only after the library contracts are
  stable.
- Keep every step independently reviewable with a small owned target and an
  exact commit subject.

## Non-goals and safety boundaries

- Do not introduce external dependencies; use the standard library, including
  `context`, `sync`, `testing`, and `runtime` diagnostics where justified.
- Do not build a general-purpose programming language, optimizer, bytecode VM,
  distributed evaluator, network service, or persistent variable store.
- Do not claim that more goroutines prove better performance. The channel design
  is a bounded learning experiment and must retain a sequential baseline.
- Do not share mutable input maps between concurrent evaluation requests.
- Do not close channels from the wrong owner, send on closed channels, leak
  workers on cancellation, or hide errors in goroutine logs.
- Do not use sleeps as the primary synchronization mechanism in tests. Use
  barriers, channels, contexts, bounded timeouts, and observable state.
- Do not run unbounded stress tests, production services, network calls, or
  destructive cleanup as part of this work item.
- Do not change unrelated language tracks, generated binaries, iterator logs,
  credentials, or host configuration.

## Ordered steps and evidence gates

Steps are serial. The next step may begin only after the prior step's design,
implementation, tests, review, exact commit, and clean-worktree handoff are
recorded. Each step record repeats its own required checkpoints so a step can be
processed and reviewed independently.

| Step | Target scope | Depends on | Commit subject |
| ---: | --- | --- | --- |
| 01 | `go/basic/calc/math/` sequential AST/evaluator contract | none | `WI-0005: 01: refactor(calc): establish sequential evaluator contract` |
| 02 | `go/basic/calc/math/` identifiers and per-request state | 01 | `WI-0005: 02: feat(calc): add reusable input environments` |
| 03 | `go/basic/calc/math/` typed jobs and result/error protocol | 02 | `WI-0005: 03: feat(calc): define typed evaluation jobs` |
| 04 | `go/basic/calc/math/` evaluator lifecycle and root submission | 03 | `WI-0005: 04: feat(calc): add evaluator lifecycle` |
| 05 | `go/basic/calc/math/` type dispatch and recursive orchestration | 04 | `WI-0005: 05: feat(calc): dispatch expression types` |
| 06 | `go/basic/calc/math/` unary and operator workers | 05 | `WI-0005: 06: feat(calc): dispatch unary operators` |
| 07 | `go/basic/calc/math/` binary and operator workers | 06 | `WI-0005: 07: feat(calc): dispatch binary operators` |
| 08 | `go/basic/calc/math/` bounded queues and backpressure | 07 | `WI-0005: 08: feat(calc): bound evaluator admission` |
| 09 | `go/basic/calc/math/` cancellation, errors, and shutdown | 08 | `WI-0005: 09: bug(calc): make evaluator shutdown reliable` |
| 10 | `go/basic/calc/math/` reuse, races, leaks, and stress evidence | 09 | `WI-0005: 10: test(calc): verify concurrent evaluator reuse` |
| 11 | `go/basic/calc/` CLI and `init` integration | 10 | `WI-0005: 11: feat(calc): integrate channel evaluator` |
| 12 | `go/basic/calc/` benchmarks and evaluator report | 11 | `WI-0005: 12: doc(calc): record evaluator evidence` |

## Delivery contract

The root documentation must be committed before any implementation step. Each
step is processed serially with clean-workspace enforcement and no forced dirty
workspace recovery. The worker owns only the selected target scope and its own
step handoff. The advisor reviews the implementation, staged scope, exact
validation output, and evidence limits before the next step is eligible.

Every step must leave a runnable baseline. If a step cannot meet its evidence
gate, mark it `blocked` or `implemented-unproven`; do not weaken the test or
silently advance the dependency chain.

## Validation and evidence requirements

Before the root record commit:

```bash
git check-ref-format --branch wi-0005-channel-evaluator
git diff --check
```

For every step, the worker and advisor must verify:

- the selected target scope and step handoff are the only changed paths;
- design decisions and non-goals are recorded before implementation claims;
- `gofmt`, focused tests, the applicable package test suite, `go vet`, and
  `git diff --check` pass where applicable;
- the test output demonstrates the behavior needed by the next step;
- cancellation, shutdown, queue, error, and race claims are limited to the
  commands actually run;
- the exact declared commit subject is used; and
- the step handoff records implementation, tests, review, evidence tier,
  limitations, and the exact next action.

Concurrency evidence must use bounded timeouts and deterministic synchronization.
Steps 09–10 additionally require `go test -race ./...` or a narrowly scoped race
equivalent, explicit goroutine-count/leak checks with tolerances, and tests for
both successful completion and cancellation/abort paths. Benchmarks in step 12
are measurements, not correctness proof.

## Not proven or blocked

- This README and the step records define intended contracts only; they do not
  prove implementation correctness, race freedom, absence of all goroutine
  leaks, performance, portability, or learner understanding.
- The current calculator AST evaluates literal expressions but has no identifier
  environment or channel evaluator. Those are future step deliverables.
- Prefix/postfix increment currently produce numeric results for literal ASTs;
  true mutable-binding semantics remain unimplemented until step 02.
- The channel architecture may be more complex and slower than direct recursion.
  Step 12 must measure rather than assume its value.
- No external services, network dependencies, production load, or cross-host
  behavior is part of this work item.

## Completion record

Reserved for the final review. It must list all 12 implementation commits,
advisor review evidence, final validation commands, remaining limitations, and
whether the topic branch was integrated.