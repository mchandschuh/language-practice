---
schema: language-practice-step/v1
id: WI-0004
step: "05"
status: planned
skill: language-practice-typescript
depends_on: ["01"]
target: typescript/bootstrap.sh
---

# WI-0004 step 05 — Document TypeScript/Node bootstrap setup

## Objective

Write the executable `./typescript/bootstrap.sh` for Linux Node/npm version
selection, local TypeScript installation, package metadata and lockfiles,
`tsconfig`, verification, and runtime/module boundaries.

## Commit subject

```text
WI-0004: 05: doc(typescript): document bootstrap setup
```

## Required boundaries

- Write `./typescript/bootstrap.sh` as the sole implementation target for this
  step; the user will inspect and run it explicitly from a terminal.
- Distinguish Node installation/version management from project-local npm
  dependencies and TypeScript installation.
- Prefer a project-local TypeScript dependency and `npx tsc` over a global
  compiler for reproducibility; explain why a global install is only an optional
  one-off tool.
- Verify `node --version`, `npm --version`, local `npx tsc --version`, package
  lockfile state, `tsconfig`, type checking, and a fresh-process runtime check.
- Record local Node 24.21.0/npm 11.19.0 and the absence of local `tsc` as
  observations, not repository support decisions.
- Do not execute the script, install Node/npm/TypeScript, modify shell profiles,
  or download npm packages while processing this step.

## Validation

```bash
test -x typescript/bootstrap.sh
bash -n typescript/bootstrap.sh
grep -q 'node --version' typescript/bootstrap.sh
grep -q 'npm --version' typescript/bootstrap.sh
grep -q 'npx tsc\|typescript' typescript/bootstrap.sh
grep -q 'package-lock.json\|lockfile' typescript/bootstrap.sh
grep -q 'tsconfig\|tsc' typescript/bootstrap.sh
grep -q 'nvm\|Node 24.21.0\|apt\|package' typescript/bootstrap.sh
git diff --check
```

The script must link to or cite official Node/npm and TypeScript setup guidance
in its comments or the shared `BOOTSTRAP.md`. Static checks do not execute the
script or prove Node, npm, registry, or TypeScript behavior.

## Prerequisite

Step 01 must be proven and reviewed. No Node installation, npm registry access,
or dependency download is authorized by this step.

## Agent handoff notes

Reserved for worker evidence, exact validation, and limitations.

## Advisor review

The advisor appends review rounds here after the worker completes.