---
schema: language-practice-step/v1
id: WI-0005
step: "11"
status: planned
skill: language-practice-go
depends_on: ["10"]
target: go/basic/calc/
---

# WI-0005 step 11 — Integrate the channel evaluator with the CLI

## Objective

Integrate the reviewed evaluator with the calculator command and repository
helpers without changing library semantics. The command must expose successful
results, parse/evaluation errors, cancellation/exit behavior, and the existing
`run calc`/`test calc` workflow.

## Commit subject

```text
WI-0005: 11: feat(calc): integrate channel evaluator
```

## Scope and ownership

- Own only `go/basic/calc/` and this step's handoff update.
- Update the CLI to construct/use the reviewed evaluator and clean it up.
- Preserve argument joining/quoting behavior for expressions passed through
  `/home/michael/src/github/language-practice/init`.
- Add or update command tests without changing unrelated app dispatch.
- Do not change evaluator internals except for an integration defect reproduced
  by a command-level test.

## Design checkpoint

Record CLI input grammar, evaluator construction/close ownership, output format,
exit statuses, error stream behavior, and how a command timeout/cancellation is
represented. Include the exact `run calc` and `test calc` commands.

## Implementation checkpoint

Implement the CLI integration, command tests, and only the minimal init-helper
change required to invoke the evaluator. Keep generated binaries out of Git.

## Test and evidence checkpoint

Add command-level tests proving:

- `run calc 1 + 2` returns `3`;
- precedence, parentheses, identifiers/inputs if exposed, and increment
  semantics produce documented output;
- malformed expressions return a nonzero status and useful stderr;
- missing expressions return the usage status;
- quoted expression arguments preserve intended input boundaries; and
- `test calc` runs the calculator package suite.

Run `bash -n init`, the helper commands, `go test ./...`, `go vet ./...`, and
focused command tests. Step 12 requires successful and failed CLI evidence.

## Review checkpoint

Review CLI/evaluator ownership, process exit status, stdout/stderr contracts,
argument forwarding, cleanup, generated-file handling, and package scope.

## Commit checkpoint

Stage only the selected CLI/helper files and tests plus this handoff, inspect the
cached diff, and commit with the exact declared subject.

## Prerequisite and next gate

Step 10 must prove concurrent reuse and race evidence. Step 12 requires the
command contract and test helper to be stable and reviewed.

## Agent handoff notes

Reserved for CLI design, implementation, command output, tests, review,
limitations, and next action.

### Design checkpoint

The CLI grammar is `run calc [expression-argument...]`. The Go command joins all
arguments with one space, preserving the existing `init` helper's `"$@"`
forwarding behavior. Thus both `run calc 1 + 2` and `run calc "1 + 2"`
evaluate the same expression. The command does not expose identifier bindings;
an identifier therefore reports the evaluator's undefined-identifier error.

The command constructs one `EvaluatorRuntime` with `context.Background()` and
`DefaultQueueCapacity` after parsing, owns it for the invocation, and defers
`Close`. It submits the parsed AST with an empty request environment, receives
the single terminal result, and never closes evaluator-owned channels. A
successful result is printed as a decimal value plus newline on stdout. Missing
arguments print `usage: calc <expression>` to stderr and return status 2. Parse,
submission, evaluation, and cancellation errors print their concrete error to
stderr and return status 1; stdout remains empty on failures. The testable runner
accepts a context so cancellation is represented by the evaluator error and
status 1, while the production command currently uses a background context and
does not add a signal or timeout flag.

The exact helper commands are:

```bash
source ./init
run calc 1 + 2
run calc "1 + 2"
test calc -run TestRunEvaluatesExpression -count=1
```

### Implementation summary

Refactored `main.go` into a thin `main` wrapper and testable `run` function.
The command now parses once, constructs and cleans up the channel evaluator,
submits the AST, maps terminal errors to stderr/status 1, and writes successful
values to stdout. Added command tests for arithmetic, precedence, parentheses,
prefix/postfix literal increment semantics, argument joining, missing input,
parse errors, undefined identifiers, canceled contexts, and fresh-process
`init` helper coverage for both `run calc` and `test calc`. No `init` change was
needed because its existing argument forwarding already preserves boundaries.

### Validation and review

- `gofmt -w go/basic/calc/main.go go/basic/calc/main_test.go`: passed.
- `bash -n init`: passed.
- `cd go && go test ./basic/calc/... -count=1`: passed.
- `cd go && go vet ./basic/calc/...`: passed.
- `bash -c 'source ./init; run calc 1 + 2'`: passed, output `3`.
- `bash -c 'source ./init; run calc "1 + 2"'`: passed, output `3`.
- `bash -c 'source ./init; test calc -run TestRunEvaluatesExpression -count=1'`: passed.
- `cd go && go test ./... -count=1`: passed.
- `cd go && go vet ./...`: passed.

Review found and corrected two test-only issues during validation: the parser's
documented postfix literal semantics return the original value (`2++` returns
`2`), and the helper subprocess root path was corrected to the repository root.
Final review confirmed evaluator ownership is deferred and joined, stdout and
stderr are separated, exit statuses are explicit, arguments are forwarded
without changing `init`, and no evaluator internals were modified.

Evidence tier: focused deterministic command tests, fresh-process helper/CLI
audits, package tests, repository-wide tests, vet, shell syntax, formatting, and
static review. This proves the tested local CLI contracts and helper workflow;
it does not prove OS signal handling because no signal/timeout interface is
exposed, nor universal process-level cancellation behavior.

### Next action

Advisor review and step 12 evaluator benchmarks/evidence report. The command
contract and `test calc` helper are now stable for that report.