---
schema: language-practice-step/v1
id: WI-0004
step: "02"
status: planned
skill: language-practice-python
depends_on: ["01"]
target: python/bootstrap.sh
---

# WI-0004 step 02 — Document Python bootstrap setup

## Objective

Write the executable `./python/bootstrap.sh` for a Linux Python setup that
separates interpreter selection, disposable `venv` creation, package
installation, project metadata, dependency recording, and verification.

## Commit subject

```text
WI-0004: 02: doc(python): document bootstrap setup
```

## Required boundaries

- Write `./python/bootstrap.sh` as the sole implementation target for this step;
  the user will inspect and run it explicitly from a terminal.
- Explain `python -m venv`, activation/deactivation, recreation, and why the
  environment is not committed or treated as movable project state.
- Distinguish `python -m pip` from global pip and explain requirements or lockfile
  decisions without inventing a package-management contract.
- Verify interpreter, venv, pip, imports, and a bounded project command.
- Identify system package, user-local, and version-manager alternatives without
  silently selecting one.
- The iterator must not execute the script, install Python, create a venv,
  download packages, or modify shell state while processing this step; those
  are user-run actions of the default bootstrap path.

## Validation

```bash
test -x python/bootstrap.sh
bash -n python/bootstrap.sh
grep -q 'python -m venv' python/bootstrap.sh
grep -q 'python -m pip' python/bootstrap.sh
grep -q 'recreate\|already exists' python/bootstrap.sh
grep -q 'sudo\|apt\|package' python/bootstrap.sh
grep -q -- '--install-method\|installing Python toolchain' python/bootstrap.sh
git diff --check
```

The script must link to or cite the official Python `venv` and Python Packaging
User Guide sources in its comments or the shared `BOOTSTRAP.md`, and state that
local Python 3.14.4 is only an observation. Static checks do not run the script
or prove Python installation, package downloads, or environment correctness.

## Prerequisite

Step 01 must be proven and reviewed. No installation or dependency download is
authorized for the iterator; the user-run script installs missing Python by
default.

## Agent handoff notes

Implemented only `./python/bootstrap.sh`. The script installs missing Python
host prerequisites by default through a visible detected distribution package
command, then keeps the interpreter, disposable `venv`, package installation,
optional dependency recording, and bounded verification explicit. It does not
select a repository support version, run installation during this iterator,
mutate shell profiles, or infer a project command that is not present.

Static validation performed by the worker:

```text
test -x python/bootstrap.sh                 PASS
bash -n python/bootstrap.sh                 PASS
grep -q 'python -m venv' python/bootstrap.sh PASS
grep -q 'python -m pip' python/bootstrap.sh  PASS
grep -q 'recreate\|already exists' python/bootstrap.sh PASS
grep -q 'sudo\|apt\|package' python/bootstrap.sh PASS
grep -q -- '--install-method\|installing Python toolchain' python/bootstrap.sh PASS
git diff --check                            PASS
```

These are tier-1 checks only. The script was not executed; Python installation,
venv creation, package downloads, host compatibility, imports on the target
host, and project correctness remain unproven. No package manager, network,
sudo, installer, build, test, or cleanup command was run by the worker.

## Advisor review

### Review round 1

- Advisor: API assistant (identity not available)
- Work-item: `WI-0004`
- Step: `02`
- Iterator run: `/home/michael/src/github/language-practice/.cline-iterations/WI-0004-02-python-bootstrap.md/20261014_050313/`
- Implementation commit: `7d0b320`
- Reviewed snapshot: topic branch `wi-0004-toolchain-bootstrap`, commit `7d0b320`
- Reviewed at (UTC): `2026-10-14T05:12:11Z`

### Findings

- Severity: blocker
- Finding: `--venv` and `--record` are checked with lexical absolute/prefix
  comparisons but are not canonicalized before mutation. A user-provided path
  containing `..` can therefore pass the prefix check while resolving outside
  the intended repository or disposable `.venv` scope. The `--recreate` guard
  has the same lexical-path weakness.
- Evidence: static inspection of `python/bootstrap.sh` lines 80–109 and 137–140;
  `bash -n`, executable-mode, content, and whitespace checks passed, but no
  runtime path-escape test was run because scripts are user-run only.
- Affected paths: `python/bootstrap.sh` and this step handoff.

### Disposition

- Decision: blocked
- Follow-up: Add canonical path resolution or an equivalent component-safe
  confinement check for every user-controlled path before accepting this step.
  Do not merge WI-0004 or ask the user to run this script until repaired and
  re-reviewed.

### Review round 2

- Advisor: API assistant (identity not available)
- Work-item: `WI-0004`
- Step: `02`
- Iterator run: `/home/michael/src/github/language-practice/.cline-iterations/WI-0004-02-python-bootstrap.md/20261014_050313/`
- Implementation commit: `7d0b320`
- Reviewed snapshot: topic branch `wi-0004-toolchain-bootstrap`, commit `c3321d1`
- Reviewed at (UTC): `2026-10-14T05:15:00Z`

### Findings

- Severity: none
- Finding: The operator explicitly accepted user-selected paths outside the
  repository as intentional for this learning repository. The prior path-
  confinement finding is therefore a rejected design preference, not a defect.
  The script's visible `--venv`, `--record`, and `--recreate` choices remain the
  user's responsibility; no hidden cleanup or automatic external mutation is
  introduced by the script.
- Evidence: operator decision in the current task; static script review,
  executable mode, `bash -n`, declared validation, exact implementation subject,
  and preserved iterator run `20261014_050313`. The script was not executed.
- Affected paths: `python/bootstrap.sh` and this step handoff.

### Disposition

- Decision: accepted
- Follow-up: Step 02 is proven for static script scope and provenance. User-run
  Python installation, environment creation, package behavior, external-path
  behavior, and project correctness remain unproven.

### Review round 3 — corrective iterator review

- Advisor: API assistant
- Work-item: `WI-0004`
- Step: `02`
- Iterator run: `/home/michael/src/github/language-practice/.cline-iterations/WI-0004-02-python-bootstrap.md/20261014_053455/`
- Reviewed snapshot: topic branch `wi-0004-toolchain-bootstrap`, fixup snapshot before round 3 repair
- Reviewed at (UTC): `2026-10-14T05:35:33Z`

### Findings

- Severity: major
- Finding: The default path could select an interpreter with `venv` but without
  usable pip support, then begin project-local venv setup before detecting that
  host prerequisite gap.
- Evidence: iterator result `initialResult` in `session.json`; static review of
  the pre-repair `python/bootstrap.sh` path; `bash -n`, executable-mode, and
  `git diff --check` passed. No installer, package manager, sudo, network, or
  script execution occurred.
- Affected paths: `python/bootstrap.sh` and this step handoff.

### Disposition

- Decision: corrective step completed; pending final static re-review
- Follow-up: The repair adds a host `python -m pip --version` check and repeats
  the visible package installation path before creating the project venv. The
  previously accepted external-path policy remains unchanged. User-run Python
  installation, pip behavior, venv creation, and project correctness remain
  unproven.