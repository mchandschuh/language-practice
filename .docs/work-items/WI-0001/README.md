---
schema: language-practice-work-item/v1
id: WI-0001
status: DESIGNED
type: feature
title: Initialize the language-practice Agent Skills set
repository: /home/michael/src/github/language-practice
branch: wi-0001-initialize-skills
authority: local work-item records under this directory
---

# WI-0001 — Initialize the language-practice Agent Skills set

## Authority and status

- **Raw status:** `DESIGNED`
- **Derived status:** `designed`
- **Authority:** this README and the ordered `steps/` records
- **Repository:** `/home/michael/src/github/language-practice`
- **Branch:** `wi-0001-initialize-skills`
- **Baseline:** `0d67ff8` (`Initial commit`)
- **Root record commit:** `WI-0001: doc(repo): prepare ordered skill delivery plan`
- **Implementation commits:** one commit per ordered step after this directory is committed

## Scope

Create the initial repository-specific Agent Skills set for the four learning
tracks: Go, Python, Rust, and TypeScript. The skills establish:

- repository bootstrap and development process;
- learning-first handwrite/collaborate/assist boundaries;
- Git hygiene and work-item delivery;
- bounded `cline-iterate` orchestration;
- authoritative documentation research;
- testing and debugging practice;
- observability, performance, reliability, and concurrency practice;
- infrastructure-oriented design; and
- language-specific and cross-language guidance.

The existing skill files under `.cline/skills/` are the implementation targets.
The work-item records are documentation and delivery coordination; they do not
replace the skill files or claim that the skills have been committed until the
ordered steps are processed.

## Non-goals

- Do not implement language exercises or application code.
- Do not add language toolchain manifests yet.
- Do not install dependencies or change global tool configuration as part of
  this work item.
- Do not commit unrelated changes, including the global skills repository's
  `text-extraction` work.
- Do not claim runtime, learner, or production evidence from static skill
  validation.

## Delivery contract

The root work-item documentation is committed first with the record-level
subject shown above. Each implementation-bearing step is then processed in
ascending order and committed individually using:

```text
WI-0001: SS: feat(lang): short imperative description
```

Each step commit must stage only its declared skill file and the selected step
record's execution update, if that update is made during processing. The `.docs`
directory commit is the only commit in this initial preparation phase that may
contain the complete work-item directory.

After each iterator round, the advisor appends an `## Advisor review` section
to the relevant step record and commits that documentation separately before
integration. If a major finding requires another round, the corrective step and
updated dependency contract must be documented and committed before the next
iterator invocation.

## Dependency-ordered steps

| Step | Skill | Depends on | Commit subject |
| ---: | --- | --- | --- |
| 01 | `language-practice-git` | none | `WI-0001: 01: feat(repo): add Git hygiene skill` |
| 02 | `language-practice-process` | 01 | `WI-0001: 02: feat(repo): add development process skill` |
| 03 | `language-practice-learning-first` | 02 | `WI-0001: 03: feat(repo): add learning-first skill` |
| 04 | `language-practice-bootstrap` | 02, 03 | `WI-0001: 04: feat(repo): add bootstrap skill` |
| 05 | `language-practice-research` | 02, 03 | `WI-0001: 05: feat(repo): add research skill` |
| 06 | `language-practice-cline-iterate` | 01, 02, 03 | `WI-0001: 06: feat(repo): add cline iterate skill` |
| 07 | `language-practice-testing` | 02, 03, 05 | `WI-0001: 07: feat(repo): add testing skill` |
| 08 | `language-practice-debugging` | 02, 03, 05, 07 | `WI-0001: 08: feat(repo): add debugging skill` |
| 09 | `language-practice-observability-performance` | 05, 07, 08 | `WI-0001: 09: feat(repo): add observability skill` |
| 10 | `language-practice-reliability-concurrency` | 05, 07, 08, 09 | `WI-0001: 10: feat(repo): add reliability skill` |
| 11 | `language-practice-infrastructure-design` | 02, 03, 05, 07, 09, 10 | `WI-0001: 11: feat(repo): add infrastructure design skill` |
| 12 | `language-practice-go` | 05, 07, 08, 09, 10 | `WI-0001: 12: feat(go): add Go conventions skill` |
| 13 | `language-practice-python` | 05, 07, 08, 09, 10 | `WI-0001: 13: feat(python): add Python conventions skill` |
| 14 | `language-practice-rust` | 05, 07, 08, 09, 10 | `WI-0001: 14: feat(rust): add Rust conventions skill` |
| 15 | `language-practice-typescript` | 05, 07, 08, 09, 10 | `WI-0001: 15: feat(typescript): add TypeScript conventions skill` |
| 16 | `language-practice-cross-language-comparison` | 05, 12, 13, 14, 15 | `WI-0001: 16: feat(repo): add cross language skill` |

The order is intentionally serial even where independent language steps could
be parallelized: each step is a separately reviewable Git checkpoint, and the
requested delivery mode is ordered processing.

## Validation and evidence

Root-record validation before the documentation commit:

```bash
npx --yes skills-ref@0.1.5 validate .cline/skills/language-practice-git
python3 <work-item-structure-check>
git diff --check
```

The root record proves only that the delivery plan and step records are
documented. It does not prove that any target skill is committed, loaded by a
runtime, or behaviorally evaluated.

Each later step must validate its target skill with `skills-ref`, check relative
references and whitespace, inspect its staged diff, and record exact evidence
and limitations before the corresponding commit.

## Not proven or blocked

- No implementation step has been processed at the time this record is written.
- No target skill has been committed by WI-0001 yet.
- No runtime activation, prompt-evaluation, learner-outcome, or language
  toolchain evidence is claimed.
