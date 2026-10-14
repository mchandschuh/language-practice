---
schema: language-practice-step/v1
id: WI-0004
step: "06"
status: proven
skill: language-practice-infrastructure-design
depends_on: ["01"]
target: cpp/bootstrap.sh
---

# WI-0004 step 06 — Document C++ bootstrap setup

## Objective

Write the executable `./cpp/bootstrap.sh` for Linux C++ compiler selection,
standard-version decisions, CMake installation, out-of-source configuration,
presets, build/test verification, and cleanup boundaries.

## Commit subject

```text
WI-0004: 06: doc(cpp): document bootstrap setup
```

## Required boundaries

- Write `./cpp/bootstrap.sh` as the sole implementation target for this step;
  the user will inspect and run it explicitly from a terminal.
- Distinguish using distribution GCC/Clang packages from building a compiler
  from source; do not prescribe source builds for ordinary learner setup.
- Require an explicit C++ standard and compiler decision before exercises depend
  on features such as concepts, ranges, `jthread`, or coroutines.
- Explain CMake configure/build/test, out-of-source build directories,
  `CMakePresets.json` versus user-local presets, and toolchain-file boundaries.
- Verify compiler versions, CMake version, selected standard, configure/build,
  and tests without claiming portability from one host.
- Record local GCC 15.2 and Clang 21.1.8, with CMake unavailable, as
- Do not execute the script, install packages, run CMake, or mutate build state
  while processing this step.

## Validation

```bash
test -x cpp/bootstrap.sh
bash -n cpp/bootstrap.sh
grep -q 'g++ --version\|clang++ --version' cpp/bootstrap.sh
grep -q 'cmake --version' cpp/bootstrap.sh
grep -q 'CMakePresets.json\|cmake' cpp/bootstrap.sh
grep -q 'out-of-source\|out of source\|build' cpp/bootstrap.sh
grep -q 'apt\|package\|GCC 15.2' cpp/bootstrap.sh
git diff --check
```

The script must link to or cite official CMake, GCC, Clang, and ISO C++ setup
sources in its comments or the shared `BOOTSTRAP.md`. Static checks do not
execute the script or prove compiler, CMake, standard, generator, or portability
behavior.

## Prerequisite

Step 01 must be proven and reviewed. No compiler/package installation, CMake
execution, or build-directory mutation is authorized by this step.

## Agent handoff notes

### Agent execution

- Agent identity: API assistant
- Agent/session ID: not applicable; implemented directly after prior iterator
  artifacts were found but not present in the current checkout
- Provider/model: not applicable
- Capability: repository implementation and static validation
- Use case: execute WI-0004 step 06
- Repository root: `/home/michael/src/github/language-practice`
- Work-item ID: `WI-0004`
- Step file: `.docs/work-items/WI-0004/steps/WI-0004-06-cpp-bootstrap.md`
- Execution mode: direct implementation with static-only validation
- Effective task constraints: write only `cpp/bootstrap.sh` plus this handoff;
  never execute the script, install packages, run CMake/compiler/build/test,
  access the network, mutate build state, or edit shell profiles; use the
  declared commit subject.

### Work performed

- Added executable `./cpp/bootstrap.sh` as the sole implementation target.
- Documented explicit `g++`/`clang++` selection, C++17/20/23 selection, and the
  fact that local GCC 15.2 and Clang 21.1.8 are observations rather than support
  commitments.
- Added visible apt, dnf, and pacman prerequisite alternatives for compiler and
  CMake installation; compiler source builds and user-local archives remain
  explicit user-managed alternatives outside automation.
- Added compiler/CMake verification and an opt-in out-of-source CMake configure,
  selected-standard check, bounded build, and CTest path.
- Documented `CMakePresets.json` versus `CMakeUserPresets.json`, toolchain-file
  boundaries, cleanup, rollback, and evidence limitations.
- No compiler, CMake, package manager, `sudo`, network, build directory, project
  configure/build/test, or cleanup command was executed.

### Files changed

- `cpp/bootstrap.sh`
- `.docs/work-items/WI-0004/steps/WI-0004-06-cpp-bootstrap.md`

### Validation and evidence

- `test -x cpp/bootstrap.sh` — PASS.
- `bash -n cpp/bootstrap.sh` — PASS.
- Required compiler, CMake, presets, build, package, and observation content
  checks — PASS.
- Additional static audit for visible package commands, standard/compiler choices,
  no shell-profile mutation, opt-in project checks, out-of-source paths, and
  conflict markers — PASS.
- `git diff --check` — PASS.
- Evidence tier: static/documentation evidence only. Compiler/CMake installation,
  CMake configuration, build/test behavior, portability, package authenticity,
  host compatibility, and learner understanding remain unproven.

### Exact next actions

1. Review the staged diff and commit only this script and handoff.
2. Have the learner choose a compiler, minimum C++ standard, generator/preset,
   CMake source/build paths, and any toolchain file before running the script.
3. Run `./cpp/bootstrap.sh` explicitly from the repository root; use
   `--check-project` only when a real CMake project exists and record host/project
   evidence separately.

## Advisor review

### Review round 1 — implementation static review

- Advisor: API assistant
- Work-item: `WI-0004`
- Step: `06`
- Reviewed snapshot: topic branch `wi-0004-toolchain-bootstrap`, based on `ad93be0`
- Reviewed at (UTC): `2026-10-14`

### Findings and disposition

- Severity: none for static scope and requirements.
- Finding: The missing C++ bootstrap target is now present and executable. The
  implementation keeps installation and project-build mutations visible, bounded,
  and opt-in, and does not execute them during repository processing.
- Decision: proven for static script scope; user-reported execution outcome
  reconciled below. Portability, project correctness, and learner understanding
  remain unproven.

### Review round 2 — user-run outcome reconciliation

- Advisor: API assistant
- Work-item: `WI-0004`
- Step: `06`
- Reviewed at (UTC): `2026-10-14T22:38:05Z`

### Findings and disposition

- Severity: none for the reported outcome; evidence limitation remains.
- Finding: The operator reports that `cpp/bootstrap.sh` completed successfully
  when run by the user. No command transcript or independent host/project
  artifact was supplied, so this is user-reported host evidence only.
- Decision: proven for static script scope and user-reported execution outcome;
  portability, project correctness, and learner understanding remain unproven.