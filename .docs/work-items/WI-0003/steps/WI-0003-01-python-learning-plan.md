---
schema: language-practice-step/v1
id: WI-0003
step: "01"
status: proven
skill: language-practice-python
depends_on: []
target: python/plans/01-infrastructure-learning-plan.md
---

# WI-0003 step 01 — Add the Python infrastructure learning plan

## Objective

Create the canonical Python adaptation of the learning-plan draft at commit
`670bafa6a5621f4f14ef91b8d71b6d844478fbb3`, using its infrastructure objective,
experiment contract, evidence boundaries, report protocol, and ordered `PY-00`
through `PY-19` curriculum as the shared baseline for the remaining plans.

## Commit subject

```text
WI-0003: 01: doc(python): add infrastructure learning plan
```

## Scope

- Add only `python/plans/01-infrastructure-learning-plan.md` and this selected
  step record's execution handoff update.
- Preserve the learner profile and `handwrite`/`collaborate`/`assist` modes.
- Preserve the experiment, load, visualization, reporting, and evidence
  contracts from the historical draft.
- Include all 20 shared curriculum stages with Python-specific APIs and runtime
  limits.
- State that the repository has no declared Python version; local Python 3.14.4
  is an observation, not a repository contract.

## Required language boundaries

- Distinguish Python annotations from runtime validation.
- Explain object identity, aliasing, mutation, cleanup, exceptions, and the
  selected concurrency model rather than hiding them behind frameworks.
- Gate version-specific features against the selected repository version rather
  than assuming Python 3.15 documentation applies to local execution.
- Prefer standard-library-first experiments and explicitly bound any future
  third-party tooling.

## Validation

```bash
test -f python/plans/01-infrastructure-learning-plan.md
grep -nE '^## (Infrastructure concept map|The experiment contract|Load-testing ladder|Visualization ladder|GitHub experiment documentation|Curriculum|Standard-library anchors|Evidence boundaries)$' python/plans/01-infrastructure-learning-plan.md
grep -nE '^## PY-(0[0-9]|1[0-9]) — ' python/plans/01-infrastructure-learning-plan.md
grep -q '670bafa6a5621f4f14ef91b8d71b6d844478fbb3' python/plans/01-infrastructure-learning-plan.md
grep -q 'Python 3.14.4' python/plans/01-infrastructure-learning-plan.md
git diff --check
```

The stage check must find exactly one heading for each stage `PY-00` through
`PY-19`. Static documentation checks do not prove Python execution, learner
understanding, or production behavior.

## Prerequisite

The WI-0003 root documentation commit must be an ancestor of `HEAD`. The current
branch must be `wi-0003-learning-plans`, and the worktree must be clean before
the iterator run.

## Agent handoff notes

This section is reserved for the iterator worker. It must record the selected
skills, run reference, files changed, exact validation, evidence tier, and
limitations without copying raw transcripts or secrets.

### Worker handoff

- Skills: `language-practice-process`, `language-practice-git`,
  `language-practice-learning-first`, `language-practice-research`,
  `language-practice-python`, `language-practice-observability-performance`,
  `language-practice-reliability-concurrency`, and `language-practice-testing`.
- Run reference: direct bounded implementation for WI-0003 step 01.
- Files changed: `python/plans/01-infrastructure-learning-plan.md` and this
  selected step handoff.
- Evidence tier: implemented-unproven documentation evidence. The document
  records intended experiments and boundaries; it does not prove execution,
  learner understanding, or production behavior.
- Validation: the exact commands in this step were run before commit; the
  additional one-heading-per-stage check, staged diff inspection, and cached
  whitespace check also passed.
- Limitations: the repository still has no declared Python version; Python
  3.14.4 remains only a local observation. No experiments, dependencies,
  manifests, or generated artifacts were created.

## Advisor review

### Review round 1

- Advisor: API assistant (identity not available)
- Work-item: `WI-0003`
- Step: `01`
- Iterator run: `/tmp/language-practice-cline-runs/WI-0003-step01-rerun/WI-0003-01-python-learning-plan.md/20261014_041453/`
- Implementation commit: `6700348bcb5ad7e5d318761d698c84eb1ff5ed0a`
- Reviewed snapshot: topic branch `wi-0003-learning-plans`, commit `6700348`
- Reviewed at (UTC): `2026-10-14T04:20:00Z`

### Findings

- Severity: minor
- Finding: The worker's initial Python plan had all 20 stage headings but
  omitted explicit `Concepts` content for `PY-12` and `PY-18`, explicit
  handwrite work for `PY-17`, and the complete stage contract for `PY-19`.
- Evidence: the advisor's per-stage contract check identified the omissions;
  the amended implementation adds the missing fields, and the corrected check
  passes for every stage `PY-00` through `PY-19`.
- Affected paths: `python/plans/01-infrastructure-learning-plan.md` and this
  selected step record.

### Disposition

- Decision: accepted
- Follow-up: The amended implementation is accepted for the static plan
  contract. The plan documents intended experiments only; Python execution,
  learner understanding, and production behavior remain unproven.

### Validation evidence

- The exact selected-step validation commands passed.
- A custom structural check passed with exactly 20 ordered stage headings and
  complete infrastructure connection, concepts, and learner-work/experiment
  fields for every stage.
- Required shared sections, draft provenance, Python 3.14.4 limitation, and
  official documentation anchors were confirmed.
- `git show --check 6700348` passed, and the implementation commit contains
  only the target plan and this selected step handoff.
- The post-review worktree was clean before this review edit.
- Evidence tier: proven for scoped static documentation and Git provenance;
  implemented-unproven for all future language behavior, experiments, learner
  outcomes, and production claims.