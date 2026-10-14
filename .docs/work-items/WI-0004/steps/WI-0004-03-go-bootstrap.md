---
schema: language-practice-step/v1
id: WI-0004
step: "03"
status: planned
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

Reserved for worker evidence, exact validation, and limitations.

## Advisor review

The advisor appends review rounds here after the worker completes.