---
schema: language-practice-step/v1
id: WI-0004
step: "06"
status: planned
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

Reserved for worker evidence, exact validation, and limitations.

## Advisor review

The advisor appends review rounds here after the worker completes.