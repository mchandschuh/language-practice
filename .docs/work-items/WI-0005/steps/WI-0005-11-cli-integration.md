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