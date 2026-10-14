---
name: language-practice-git
description: Maintain safe, reviewable Git history for language-practice work through WI/GH work-item branches, validated checkpoint commits, exact declared subjects, and verified integration. Use whenever repository files are inspected, changed, tested, committed, branched, merged, rebased, or reviewed.
compatibility: Designed for the language-practice repository with Git, a current work-item record, and non-interactive validation commands available.
metadata:
  version: "1"
  scope: "repository-wide Git hygiene"
---

# Language-practice Git hygiene

This skill owns the repository's branch, commit, staging, integration, and Git
evidence policy. The global `git` skill owns universal Git safety and pager
rules. `list-work-item` owns work-item identity and ordered-step authority.
Explicit user instructions override these defaults.

## Repository state before work

Before editing or relying on Git state:

```bash
GIT_PAGER=cat PAGER=cat git --no-pager rev-parse --show-toplevel
GIT_PAGER=cat PAGER=cat git --no-pager status --short --branch
GIT_PAGER=cat PAGER=cat git --no-pager branch --show-current
GIT_PAGER=cat PAGER=cat git --no-pager log --oneline --decorate -n 10
```

Resolve the repository root from Git, not from the shell's current directory.
Preserve unrelated user changes. Stop if the repository is in a merge, rebase,
or cherry-pick state until that state is understood.

Before implementation, identify:

- the current base branch;
- the authoritative `WI-NNNN` or `GH-NNNN` work-item record;
- the absolute work-item record directory and root README;
- the ordered step identifier `SS`, when the work is step-based;
- the intended branch name;
- the committed root-documentation commit;
- the selected step's exact declared commit subject;
- the validation commands and evidence tier;
- explicit user overrides or non-goals.

Do not invent a work-item identifier. If a required work item or root record is
absent, report the blocker rather than implementing or creating a record through
this skill.

An active, committed work item and a valid work-item-linked topic branch are
mandatory before implementation. A branch alone is insufficient. The complete
root work-item documentation must be an ancestor of `HEAD`, and the selected
step must own the target files and declare one exact commit subject.

## Topic branches

Work on a topic branch created from the recorded base branch. Do not make
ordinary task changes directly on `main`, `master`, or another base branch.

Use the canonical lower-case form `wi-NNNN-<slug>` for local work items and
`gh-NNNN-<slug>` for GitHub-originated work items. The slug is short,
lowercase, and hyphenated. Do not add a category prefix such as `feature/` or
`docs/` in front of this form.

```text
wi-0002-bounded-go-workers
wi-0007-rust-error-propagation
gh-0042-auth-timeout
```

Validate branch names with:

```bash
git check-ref-format --branch '<branch-name>'
```

Keep topic branches short-lived and integrate frequently. Do not create a
second branch for the same work item unless the work-item record requires it.
Record the original base branch before switching branches.

The first state-changing action for a new work item is creating this branch.
Then create and validate the work-item documentation, commit that root record,
and only then implement a step.
Local records use the required `.docs/work-items/WI-NNNN/` directory. GitHub-
originated records use `.docs/work-items/GH-NNNN-<slug>/`. Do not silently
rename an existing authoritative record.

## Commit early and often

Make a commit after each smallest coherent, reviewable checkpoint. A normal
checkpoint commit must leave the repository in a valid state:

- the intended files are the only files staged;
- no conflict markers or accidental generated files are present;
- applicable formatting, type-check, build, and test commands pass;
- related tests are included with behavior changes;
- the commit is understandable and reversible;
- the message records enough context for a future maintainer.

Do not commit arbitrary partial edits merely because time has passed, and do not
mix unrelated cleanup, refactoring, formatting, or renaming into the checkpoint.
If the user explicitly requests a WIP, intentionally failing, or exploratory
commit, record that override and its limitations in the commit body and report.

## Work-item record commit

Creating or updating the root documentation for a work item under
`.docs/work-items/WI-NNNN/` or `.docs/work-items/GH-NNNN-<slug>/` is a
documentation commit, not an implementation
step. A root work-item record has no `SS` identifier, so its subject deliberately
omits the step segment:

```text
KEY: doc(repo): short imperative description
```

For example:

```text
WI-0001: doc(repo): prepare ordered skill delivery plan
GH-0042: doc(repo): record GitHub issue delivery plan
```

Use this form only for the work-item README, ordered step files, prompt
catalogues, corrective-step records, advisor review sections, and other
explicitly scoped work-item documentation. Do not hide implementation, source,
test, or configuration changes in a root-record commit. The body must explain
the work-item scope, ordering decision, dependencies, and documentation
validation.

## Required implementation-step commit subject

For an implementation-bearing work-item step, use exactly:

```text
KEY: SS: feat(lang): short imperative description
```

Examples:

```text
WI-0002: 03: feat(go): add bounded worker exercise
WI-0007: 02: bug(rust): reject invalid borrowed input
GH-0042: 01: doc(repo): explain runtime validation
WI-0011: 04: test(python): cover generator exhaustion
```

Rules:

- `KEY` must exactly match the authoritative `WI-NNNN` or `GH-NNNN` record.
- `SS` must match the ordered step identifier, including required zero-padding.
- Use a valid type such as `feat`, `bug`, `doc`, `test`, `refactor`, `perf`,
  `build`, `ci`, or `chore`.
- Use `go`, `python`, `rust`, `typescript`, or an appropriate scope such as
  `repo` for `lang`.
- Use imperative wording, a specific description, and no trailing period.
- The complete subject line must not exceed 80 characters; target 72 or fewer
  when practical.
- Do not use vague subjects such as `update`, `changes`, or `fix` alone.

For non-step work, use `KEY: 00: ...` only if the repository's work-item
record explicitly defines `00` as the applicable step. Otherwise stop and use
the work-item's declared step identity rather than inventing one.

The actual implementation commit subject must equal the selected step's declared
subject byte-for-byte. A related or semantically equivalent subject is still a
delivery finding. A step with a missing, malformed, or mismatched subject is
blocked before implementation.

Advisor review notes are documentation, not implementation. Commit them with:

```text
KEY: doc(repo): record advisor review for step SS
```

Stage only the relevant step-file review section and explicitly scoped work-item
documentation. Do not put source or tests in that commit.

## Commit body

Non-trivial commits must contain a blank line after the subject and explain why
the change exists, not merely restate the diff. Use this structure:

```text
WI-0002: 03: feat(go): add bounded worker exercise

Why:
Describe the requirement, defect, learning objective, or rationale.

What:
Describe the repository-visible behavior or documentation change.

How:
Describe the implementation approach and important mechanics.

Decisions:
Record meaningful alternatives, tradeoffs, and rejected approaches.

Validation:
- `go test ./...` — passed
- `go vet ./...` — passed
- `git diff --check` — passed

Refs:
- WI-0002
```

Wrap body lines at 80 characters, preferably 72. Include exact validation
commands and results. Record unproven checks, blockers, design references, and
follow-up work. Use `Refs:` or `Fixes:` only when the relationship is accurate;
do not claim that a commit alone proves resolution.

## Staging and pre-commit review

Stage explicit paths, then inspect the complete staged state:

```bash
GIT_PAGER=cat PAGER=cat git --no-pager status --short
GIT_PAGER=cat PAGER=cat git --no-pager diff
GIT_PAGER=cat PAGER=cat git --no-pager diff --staged
GIT_PAGER=cat PAGER=cat git --no-pager diff --cached --check
```

Do not stage credentials, private keys, tokens, `.env` files, caches, binaries,
downloaded data, logs, iterator artifacts, temporary files, or unrelated user
changes. Check for conflict markers before committing:

```bash
grep -RInE '^(<<<<<<<|=======|>>>>>>>)' -- . ':!/.git'
```

If a staged secret or unrelated change is found, stop, inspect, and report it.
Do not use broad staging to avoid deciding which files belong to the commit.

## Validation matrix

Run the smallest complete validation set after the final relevant edit and
before each ordinary checkpoint commit:

| Change | Minimum evidence |
| --- | --- |
| Source | formatter plus compiler, type checker, or build |
| Behavior | focused tests plus the repository test suite |
| Bugfix | reproduction plus regression test when reproducible |
| Concurrency | race detector, stress test, sanitizer, or equivalent |
| Performance | representative benchmark or profile before and after |
| Skill/docs | `skills-ref validate`, reference checks, and `git diff --check` |
| Integration | complete applicable validation after the merge |

Use `proven`, `implemented-unproven`, `blocked`, and `deferred` accurately.
A commit hash establishes provenance; it does not prove correctness, production
behavior, scalability, or learner understanding.

## Rebase and `--no-ff` integration

`--no-ff` applies to `git merge`, not to `git rebase`. When repository policy
requires integration:

1. Confirm the topic branch is unpublished or that rewriting it is explicitly
   authorized.
2. Ensure the working tree is clean and record the current base commit.
3. Rebase the unpublished topic branch onto the current base when appropriate.
4. Resolve conflicts deliberately and rerun all required validation.
5. Switch to the original base branch.
6. Merge with:

   ```bash
   git merge --no-ff <topic-branch>
   ```

7. Run post-merge validation.
8. Inspect the resulting graph, merge commit, status, and diff.
9. Delete the topic branch only after successful verification and when no
   follow-up work depends on it.

Before integration, additionally verify that the root documentation commit is
present, all implementation subjects match their declared steps, every relevant
step has committed advisor review notes, all major findings are resolved or
explicitly accepted by the operator, and no unrelated worktree paths remain.

Never rebase a published/shared branch by default. Never force-push as ordinary
completion. If a rebase or merge conflicts, preserve the conflict state until
it is understood, remove all markers, validate the result, and report any
unresolved issue as blocked.

## Learning-oriented Git operations

When teaching an unfamiliar Git operation, briefly explain:

- what it changes;
- which reference or object it affects;
- whether it rewrites history;
- whether it is safe on a shared branch;
- how to inspect the result;
- how to recover safely;
- what evidence confirms success.

Safe routine operations may be automated once the learner understands them.
History-rewriting and destructive operations remain explicit and visible.

## Completion report

Report:

- work-item and step identifier;
- original base branch and topic branch;
- commits created and their subjects;
- whether a rebase occurred;
- whether integration used `--no-ff`;
- exact validation commands and results;
- final commit or merge hash;
- whether the topic branch was deleted;
- unproven, blocked, deferred, or user-overridden work.