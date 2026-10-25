---
schema: language-practice-step/v1
id: WI-0005
step: "01"
status: planned
skill: language-practice-go
depends_on: []
target: go/basic/calc/math/
---

# WI-0005 step 01 — Establish the sequential evaluator contract

## Objective

Separate the calculator AST from its evaluation runtime and establish a small,
deterministic sequential evaluator baseline. This step defines the expression
semantics that every later channel implementation must preserve.

## Commit subject

```text
WI-0005: 01: refactor(calc): establish sequential evaluator contract
```

## Scope and ownership

- Own only `go/basic/calc/math/` and this step's handoff update.
- Normalize the existing expression nodes to the selected pointer/value policy;
  do not mix pointer and value implementations in the interface.
- Make `MathExpression` describe AST structure rather than binding it to one
  evaluation strategy.
- Add a sequential `Evaluator` or equivalent package-owned evaluator with an
  explicit result/error contract.
- Preserve parsing of literals, arithmetic precedence, parentheses, unary plus
  and minus, and the currently supported increment/decrement syntax.
- Do not add identifiers, mutable environments, goroutines, channels, or CLI
  changes in this step.

## Design checkpoint

Before implementation, record in the handoff:

- whether AST nodes are pointers or values and why;
- whether `MathExpression` includes `Evaluate` or is data-only;
- the evaluator's error behavior for unsupported nodes/operators;
- the documented semantics of increment/decrement on literal expressions; and
- the sequential evaluation order for binary operands.

The design must state that later evaluators are compared against this baseline,
not against an implementation detail of a worker.

## Implementation checkpoint

Implement only the selected AST/evaluator changes. Keep exported identifiers
documented. Return errors instead of panicking for unsupported operators or node
types. Keep the parser's existing error positions and precedence behavior unless
the handoff records a deliberate contract correction.

## Test and evidence checkpoint

Add or update deterministic package tests proving:

- number, unary, binary, precedence, parentheses, and increment/decrement
  evaluation;
- malformed expressions still return errors;
- unsupported operator/node paths return errors rather than hanging or panicking;
- the parser returns the selected AST representation consistently; and
- `go test ./basic/calc/math -count=1`, `go test ./... -count=1`, `go vet ./...`,
  and `gofmt` pass.

The evidence gate for step 02 is a passing sequential baseline and a written
semantic table in the handoff. A passing test does not prove future channel
behavior or race freedom.

## Review checkpoint

The advisor must inspect AST ownership, receiver/method sets, error paths,
parser compatibility, test assertions, and the exact diff scope. Reject the step
if it introduces concurrency or leaves the sequential contract implicit.

## Commit checkpoint

Stage only the selected math files and this handoff, run `git diff --cached --check`,
verify the exact declared subject, and commit before step 02 becomes eligible.

## Prerequisite and next gate

The WI-0005 root documentation commit must be an ancestor of `HEAD`. Step 02
requires this step's implementation commit, advisor review, passing tests, and a
clean worktree.

## Agent handoff notes

Reserved for design notes, implementation summary, exact commands/results,
evidence tier, limitations, review rounds, and the exact next action. Do not
claim channel, cancellation, race, or performance evidence in this step.

### Design checkpoint

- AST nodes are pointers (`*NumberExpression`, `*UnaryExpression`, and
  `*BinaryExpression`). This gives the parser one consistent representation,
  avoids value/pointer method-set ambiguity, and makes the immutable tree
  reusable by later evaluators.
- `MathExpression` is a data-only marker interface. It describes AST structure
  and deliberately does not include `Evaluate`; the sequential `Evaluator` owns
  runtime behavior and returns `(float64, error)`.
- Unsupported node types, nil nodes, malformed child links, and unsupported
  operators return errors rather than panicking or blocking.
- Increment/decrement on literal expressions are numeric exercise semantics:
  prefix and postfix increment return the literal plus one, and prefix and
  postfix decrement return the literal minus one. They do not mutate the AST or
  imply mutable-variable semantics; identifier state is deferred to step 02.
- Binary operands are evaluated sequentially, left operand before right
  operand. Later evaluators are compared with this sequential semantic
  baseline, not with an implementation detail of a worker.

### Implementation summary

Added the stateless package-owned `Evaluator` and convenience `Evaluate`
function, moved all evaluation dispatch out of AST nodes, normalized parser
output to pointer nodes, and updated parser/evaluator tests for the contract.
No identifiers, environments, goroutines, channels, or CLI changes were made.

### Evidence and limitations

Focused tests and repository-wide Go validation are required by this step. The
current checkout has no `go.mod`, so module-mode `go test ./...` and `go vet
./...` cannot run until the repository's existing Go module setup is restored.
No channel, cancellation, race, leak, or performance behavior is claimed.

### Validation and review

- `gofmt -w go/basic/calc/math/*.go`: passed.
- `GO111MODULE=off go test ./go/basic/calc/math -count=1`: passed.
- `go test ./basic/calc/math -count=1`, `go test ./... -count=1`, and
  `go vet ./...`: blocked by the checkout having no `go.mod`.
- `GO111MODULE=off go test ./go/... -count=1` and `GO111MODULE=off go vet
  ./go/...`: blocked by the pre-existing CLI import path not resolving from
  this checkout outside its expected GOPATH/module layout.
- `git diff --check`: passed.

Review found and corrected value-node expectations left in parser tests, the
CLI's stale node `Evaluate` call was verified to remain outside this step's
scope, and evaluator recursion was kept on one stateless receiver. The exact
changed paths are this handoff plus files under `go/basic/calc/math/`.

Evidence tier: focused sequential package correctness and AST representation
evidence. The next action is to use this committed baseline as the contract for
step 02's reusable input environment; restore/define the repository Go module
layout before claiming repository-wide validation.