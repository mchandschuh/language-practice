---
schema: language-practice-work-item/v1
id: WI-0002
status: COMPLETE
type: maintenance
title: Ignore language-toolchain artifacts
repository: /home/michael/src/github/language-practice
branch: wi-0002-language-gitignore
authority: local work-item records under this directory
---

# WI-0002 — Ignore language-toolchain artifacts

## Authority and status

- **Raw status:** `COMPLETE`
- **Derived status:** `complete`
- **Authority:** this README and the ordered `steps/` records
- **Repository:** `/home/michael/src/github/language-practice`
- **Branch:** `wi-0002-language-gitignore`
- **Baseline:** `28474a5` (`Merge branch wi-0001-gitignore`)
- **Root record commit:** `WI-0002: doc(repo): prepare language ignore plan`
- **Implementation commits:** one commit per ordered step after this directory is committed

## Scope

Add conservative root `.gitignore` patterns for generated artifacts from the
selected learning tracks: Go, Python, Rust, and TypeScript/Node. Preserve the
existing WI-0001 rule for `.cline-iterations/`.

## Non-goals

- Do not ignore source files, tests, documentation, or committed lockfiles.
- Do not ignore `.env.example`; only local secret-bearing environment files are
  excluded when a pattern is justified by the selected tooling.
- Do not add language manifests, dependencies, exercises, or toolchain setup.
- Do not modify the existing `.cline-iterations/` rule except to verify it still
  works.

## Delivery contract

The root work-item documentation is committed first with the record-level
subject shown above. The implementation step is then processed serially and
committed with its exact declared subject. The advisor appends review notes to
the step and commits them separately before integration.

## Dependency-ordered steps

| Step | Target | Depends on | Commit subject |
| ---: | --- | --- | --- |
| 01 | repository `.gitignore` | none | `WI-0002: 01: chore(repo): ignore language artifacts` |

## Validation and evidence

Root-record validation before the documentation commit:

```bash
python3 <work-item-structure-check>
git diff --check
```

The implementation step must validate each intended ignore pattern with
`git check-ignore -v`, verify that representative source/configuration files
remain visible, inspect the staged diff, and run `git diff --check`.

The work item proves Git ignore matching and scoped documentation only. It does
not prove language-toolchain behavior, builds, learner outcomes, or production
behavior.

## Not proven or blocked

- Step 01 is proven by the committed ignore matrix and advisor review.
- No language toolchain manifests or runtime evidence are in scope.