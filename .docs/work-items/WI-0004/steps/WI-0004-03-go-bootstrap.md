---
schema: language-practice-step/v1
id: WI-0004
step: "03"
status: implemented-unproven
skill: language-practice-go
depends_on: ["01"]
target: go/bootstrap.sh
---

# WI-0004 step 03 — Document Go bootstrap setup

## Objective

Write the executable `./go/bootstrap.sh` for Linux Go installation choices,
PATH/GOPATH tool locations, module initialization, dependency recording,
verification, and cleanup boundaries.

## Commit subject

```text
WI-0004: 03: doc(go): document bootstrap setup
```

## Required boundaries

- Write `./go/bootstrap.sh` as the sole implementation target for this step; the
  user will inspect and run it explicitly from a terminal.
- Distinguish distribution/archive installation from user-local tool locations
  and document permissions without requiring a global mutation.
- Explain `go version`, `go env`, `go mod init`, `go mod tidy`, `go.sum`, and
  bounded `go test`/`go vet` verification.
- Explain module path selection as a project decision, not an installer detail.
- Record that no local Go toolchain or repository Go version is currently
  available; Go 1.27 is only a documentation observation.
- Do not execute the script, install Go, change PATH, initialize a module in the
  repository, or download modules while processing this step.

## Validation

```bash
test -x go/bootstrap.sh
bash -n go/bootstrap.sh
grep -q 'go version' go/bootstrap.sh
grep -q 'go env' go/bootstrap.sh
grep -q 'go mod init' go/bootstrap.sh
grep -q 'go mod tidy' go/bootstrap.sh
grep -q 'go.sum\|go\.mod' go/bootstrap.sh
grep -q 'sudo\|apt\|package' go/bootstrap.sh
git diff --check
```

The script must link to or cite the official Go installation and module guidance
in its comments or the shared `BOOTSTRAP.md`. Static checks do not execute the
script or prove Go installation, module downloads, or build correctness.

## Prerequisite

Step 01 must be proven and reviewed. No installation, PATH mutation, or module
download is authorized by this step.

## Agent handoff notes

### Agent execution

- Agent identity: API assistant
- Agent/session ID: not applicable; implemented directly after the prior iterator
  stopped on the rewritten-history gate
- Provider/model: not applicable
- Capability: repository implementation and static validation
- Use case: execute WI-0004 step 03
- Repository root: `/home/michael/src/github/language-practice`
- Work-item ID: `WI-0004`
- Step file: `.docs/work-items/WI-0004/steps/WI-0004-03-go-bootstrap.md`
- Execution mode: direct implementation with static-only validation
- Effective task constraints: write only `go/bootstrap.sh` plus this handoff;
  never execute the script, install Go, change PATH, initialize a module, or
  download modules; use the declared commit subject.

### Work performed

- Added executable `./go/bootstrap.sh` as the sole implementation target.
- Documented visible distribution-package and checksum-verified official archive
  installation alternatives, user-local/system permissions, PATH/GOPATH behavior,
  module-path selection, `go.mod`/`go.sum`, and cleanup/rollback boundaries.
- Made `go mod init`, `go mod tidy`, `go test ./...`, and `go vet ./...` explicit
  opt-in actions rather than implicit bootstrap mutations.
- Recorded Go 1.27 as a documentation observation only; no repository Go support
  version was invented.
- No Go executable, installer, package manager, `sudo`, archive download, module
  command, PATH/profile mutation, test, vet, or cleanup command was executed.

### Files changed

- `go/bootstrap.sh`
- `.docs/work-items/WI-0004/steps/WI-0004-03-go-bootstrap.md`

### Validation and evidence

- `test -x go/bootstrap.sh` — PASS.
- `bash -n go/bootstrap.sh` — PASS.
- Required `go version`, `go env`, `go mod init`, `go mod tidy`, module-file,
  and package/install-method content checks — PASS.
- Additional static audit for official references, checksum enforcement, visible
  package/archive commands, no profile mutation, opt-in project actions, and
  conflict markers — PASS.
- `git diff --check` — PASS.
- Evidence tier: static/documentation evidence only. Go installation, archive
  authenticity beyond a user-run checksum command, module resolution, project
  tests, vet results, host compatibility, and learner understanding remain
  unproven.

### Exact next actions

1. Review the staged diff and commit only this script and handoff.
2. Have the learner choose an installation method, version, architecture, source,
   checksum, and module path before running the script.
3. Run `./go/bootstrap.sh` explicitly from the repository root and record host,
   project, dependency, test, and vet evidence separately.

## Advisor review

### Review round 1 — implementation static review

- Advisor: API assistant
- Work-item: `WI-0004`
- Step: `03`
- Reviewed snapshot: topic branch `wi-0004-toolchain-bootstrap`, based on `2d9c898`
- Reviewed at (UTC): `2026-10-14`

### Findings and disposition

- Severity: none for static scope and requirements.
- Finding: The missing Go bootstrap target is now present and executable. The
  implementation keeps installation, module, dependency, test, and vet mutations
  visible and bounded, and does not execute them during repository processing.
- Decision: implemented-unproven; user-run validation required.