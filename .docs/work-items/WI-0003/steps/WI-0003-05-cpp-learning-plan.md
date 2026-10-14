---
schema: language-practice-step/v1
id: WI-0003
step: "05"
status: proven
skill: language-practice-cross-language-comparison
depends_on: ["01"]
target: cpp/plans/01-infrastructure-learning-plan.md
---

# WI-0003 step 05 — Add the C++ infrastructure learning plan

## Objective

Create a modern C++ learning plan aligned with the Python baseline's 20
behavioral areas while teaching value semantics, lifetime, RAII, ownership,
templates, standard-library facilities, concurrency, and undefined-behavior
boundaries instead of assuming garbage collection or copying Python APIs.

## Commit subject

```text
WI-0003: 05: doc(cpp): add C++ infrastructure learning plan
```

## Scope

- Add only `cpp/plans/01-infrastructure-learning-plan.md` and this selected step
  record's execution handoff update.
- Keep the shared objective, experiment contract, report protocol, evidence
  boundaries, and stage ordering aligned with step 01.
- Adapt each stage to modern C++ language and standard-library concepts.
- Use `cpp/` as the repository directory name while naming the language C++ in
  learner-facing text.

## Required language boundaries

- Teach values, references, pointers, initialization, scope, object lifetime,
  RAII, destructors, smart pointers, copy/move behavior, and ownership.
- Include classes, abstract interfaces, templates, concepts, containers,
  algorithms, ranges, exceptions, `noexcept`, and error-return alternatives
  where the selected standard supports them.
- Cover threads, `jthread`, futures, mutexes, condition variables, atomics,
  shutdown, and explicit coroutine/runtime boundaries without claiming a
  standard facility is universally available across toolchains.
- Mark sockets, JSON, structured logging, metrics, plotting, package managers,
  build systems, and async runtimes as external library/tooling choices.
- Record GCC 15.2 and Clang 21.1.8 as local observations; no C++ standard or
  build system is currently the repository contract.

## Validation

```bash
test -f cpp/plans/01-infrastructure-learning-plan.md
grep -nE '^## (Infrastructure concept map|The experiment contract|Load-testing ladder|Visualization ladder|GitHub experiment documentation|Curriculum|Official documentation anchors|Evidence boundaries)$' cpp/plans/01-infrastructure-learning-plan.md
grep -nE '^## CPP-(0[0-9]|1[0-9]) — ' cpp/plans/01-infrastructure-learning-plan.md
grep -q 'RAII' cpp/plans/01-infrastructure-learning-plan.md
grep -q 'GCC 15.2' cpp/plans/01-infrastructure-learning-plan.md
git diff --check
```

The stage check must find exactly one heading for each stage `CPP-00` through
`CPP-19`. Static documentation checks do not prove C++ compilation, defined
behavior, sanitizer coverage, portability, or production capacity.

## Prerequisite

Step 01 must have its implementation and advisor review committed and accepted.
The WI-0003 root documentation commit must be an ancestor of `HEAD`, and the
worktree must be clean before the iterator run.

## Agent handoff notes

This section is reserved for the iterator worker. Record loaded skills, run
reference, files changed, validation, evidence tier, and limitations here.

### Worker handoff

- Skills loaded: `language-practice-process`, `language-practice-git`, `git`,
  `language-practice-learning-first`, `language-practice-research`,
  `language-practice-infrastructure-design`,
  `language-practice-observability-performance`,
  `language-practice-reliability-concurrency`,
  `language-practice-cross-language-comparison`, and
  `language-practice-testing`.
- Run reference: direct bounded Act execution for WI-0003 step 05; no provider
  iterator run identifier was available in this execution context.
- Files changed: `cpp/plans/01-infrastructure-learning-plan.md` and this
  selected step handoff only.
- Work performed: created the C++ adaptation with the shared infrastructure
  objective, learner profile, `handwrite`/`collaborate`/`assist` modes,
  experiment/load/visualization/report contracts, evidence boundaries, and
  official ISO C++ getting-started, C++ Core Guidelines, and cppreference
  anchors. The plan contains exactly 20 ordered `CPP-00` through `CPP-19`
  stages covering values, references, pointers, initialization, scope, object
  lifetime, undefined behavior, RAII, destructors, smart pointers, copy/move
  behavior, ownership, classes, abstract interfaces, templates, concepts,
  containers, algorithms, ranges, exceptions, `noexcept`, error-return
  alternatives, threads, `jthread`, futures, mutexes, condition variables,
  atomics, shutdown, and explicit coroutine/runtime boundaries. Sockets, JSON,
  structured logging, metrics, plotting, package managers, build systems, and
  async runtimes are labeled external choices. GCC 15.2 and Clang 21.1.8 are
  recorded as local observations; no standard or build system is made a
  repository contract. No experiments, manifests, dependencies, or generated
  artifacts were created.
- Validation: the exact commands declared by this step passed. The additional
  structural semantic check passed with exactly 20 ordered stage headings and
  complete `Infrastructure connection`, `Concepts`, and learner-work or
  experiment/measurement fields for every stage. The staged diff was restricted
  to the selected plan and this selected step handoff; cached whitespace
  validation passed.
- Evidence tier: proven for scoped static documentation structure, source-anchor
  inclusion, and Git scope/whitespace checks; implemented-unproven for C++
  compilation, standard-version portability, defined behavior, sanitizer
  coverage, race freedom, benchmarks, learner understanding, and production
  behavior.
- Limitations: the repository has no declared C++ standard, compiler contract,
  build system, package manager, dependency set, socket/JSON/logging/metrics/
  plotting library, or async runtime. GCC `15.2` and Clang `21.1.8` remain local
  observations only. No experiments, manifests, dependencies, toolchain setup,
  or generated artifacts were created.

## Advisor review

### Review round 1

- Advisor: API assistant (identity not available)
- Work-item: `WI-0003`
- Step: `05`
- Iterator run: `/tmp/language-practice-cline-runs/WI-0003-step05/WI-0003-05-cpp-learning-plan.md/20261014_043137/`
- Implementation commit: `96c3646f2f79b48352de3e69cb69109c46733815`
- Reviewed snapshot: topic branch `wi-0003-learning-plans`, commit `96c3646`
- Reviewed at (UTC): `2026-10-14T04:36:00Z`

### Findings

- Severity: none
- Finding: The C++ plan preserves the accepted shared behavioral progression,
  provides all 20 ordered `CPP-00` through `CPP-19` stages, and makes values,
  lifetime, RAII, ownership, generic programming, standard-library facilities,
  concurrency, undefined behavior, and external tooling boundaries explicit.
- Evidence: semantic per-stage contract check, exact step validation commands,
  implementation commit path inspection, `git show --check`, selected-step
  preflight, and iterator receipt `20261014_043137` with completed status.
- Affected paths: `cpp/plans/01-infrastructure-learning-plan.md` and this
  selected step record.

### Disposition

- Decision: accepted
- Follow-up: Step 05 is proven for scoped static documentation and Git
  provenance. C++ compilation, standard-version portability, defined behavior,
  sanitizer coverage, race freedom, benchmarks, learner understanding, and
  production behavior remain unproven.