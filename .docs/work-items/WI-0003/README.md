---
schema: language-practice-work-item/v1
id: WI-0003
status: COMPLETE
type: feature
title: Create comparable infrastructure learning plans
repository: /home/michael/src/github/language-practice
branch: wi-0003-learning-plans
authority: local work-item records under this directory
---

# WI-0003 — Create comparable infrastructure learning plans

## Authority and status

- **Raw status:** `COMPLETE`
- **Derived status:** `complete`
- **Authority:** this README and the ordered `steps/` records
- **Repository:** `/home/michael/src/github/language-practice`
- **Branch:** `wi-0003-learning-plans`
- **Baseline:** `96bef8f` (`Merge branch wi-0002-language-gitignore`)
- **Source draft:** `670bafa6a5621f4f14ef91b8d71b6d844478fbb3`
  (`python/README.md` from a two-parent WIP commit)
- **Root record commit:** `WI-0003: doc(repo): prepare learning plans`
- **Implementation commits:** one commit per ordered plan step after this
  directory is committed

## Scope

Create one learning-plan document for each selected language:

- Python;
- Go;
- Rust;
- TypeScript; and
- C++ under the repository directory name `cpp/`.

The Python document adapts the draft at `670bafa` without merging or copying
that WIP commit. The five documents share the draft's infrastructure objective,
experiment contract, evidence boundaries, and ordered 20-stage curriculum while
preserving each language's actual semantics, standard-library boundaries,
tooling, and concurrency model.

The plans target writing fluency in service of dependable infrastructure:

- platform and backend foundations;
- distributed request lifecycles;
- concurrency, queues, backpressure, and cancellation;
- caching and consistency;
- timeouts, retries, rate limiting, isolation, and fallbacks;
- API and contract design;
- observability and debugging;
- throughput, latency distributions, and tail behavior;
- progressive delivery, verification, and rollback; and
- reusable libraries and operational maturity.

Each plan must contain the following shared sections:

1. track purpose and infrastructure objective;
2. learner profile and teaching strategy;
3. `handwrite`, `collaborate`, and `assist` modes;
4. version and dependency policy;
5. infrastructure concept map;
6. experiment contract;
7. load-testing ladder;
8. visualization ladder;
9. GitHub experiment/report protocol;
10. the 20-stage curriculum;
11. language-specific official documentation anchors; and
12. evidence boundaries.

Every curriculum stage must state its infrastructure connection, concepts,
learner-owned work, experiment or measurement, and limitations or reflection
checkpoint where applicable.

## Non-goals

- Do not merge, cherry-pick, or otherwise import the WIP commit `670bafa`.
- Do not implement exercises, application code, tests, manifests, dependencies,
  or toolchain setup.
- Do not install language tools or select undeclared repository versions as if
  they were contracts.
- Do not force Python APIs or terminology into the other language plans.
- Do not claim that plan documentation proves runtime behavior, benchmarks,
  learner understanding, production capacity, or production readiness.
- Do not create future `experiments/` trees or generated result artifacts.
- Do not select C++ libraries, a build system, a JSON implementation, a socket
  framework, or an async runtime without labeling it as an external boundary.
- Do not push or merge the topic branch as part of this work item.

## Learning contract

Core language concepts and reasoning remain `handwrite` tasks. Design,
diagnosis, measurement interpretation, and tradeoff analysis use
`collaborate`. Repetitive fixtures, adapters, report formatting, and harness
plumbing may use `assist` only after the behavior being measured is understood.

Each future exercise described by a plan must identify prerequisites, learner
ownership, permitted assistance, observable behavior, validation, edge cases,
failure modes, and a short reflection checkpoint. A successful compile or test
is evidence about the stated check, not proof of understanding.

## Shared curriculum ordering

The language-specific plans adapt these same behavioral areas in order:

| Stage | Shared area |
| ---: | --- |
| 00 | Writing loop, prediction, small transformations, diagnostics |
| 01 | Data structures, mutation/aliasing/ownership, state, invariants |
| 02 | Functions as contracts, callbacks, iterators/generators, side effects |
| 03 | Modules/packages, environments, CLIs, configuration, reproducibility |
| 04 | Errors/exceptions, cleanup, resource lifetime |
| 05 | Types and runtime contracts |
| 06 | Testing and regression control |
| 07 | Files, structured data, subprocesses, safe tooling |
| 08 | Sockets/request lifecycles, framing, partial I/O, malformed messages |
| 09 | Asynchronous execution fundamentals |
| 10 | Threads, processes, executors, futures, or parallel work |
| 11 | Queues, worker pools, overload, backpressure, shutdown |
| 12 | Timeouts, retries, rate limiting, isolation, fallbacks |
| 13 | Caching, TTL, invalidation, stale data, consistency |
| 14 | API and contract design, compatibility, pagination, idempotency |
| 15 | Reusable library design and misuse resistance |
| 16 | Observability and debugging systems |
| 17 | Performance, profiling, throughput, distributions, tail latency |
| 18 | Progressive delivery, verification gates, rollback simulation |
| 19 | Safe request-platform capstone and operational review |

The plans compare concepts rather than spellings. An interface, protocol, trait,
duck-typed object, promise, channel, task, exception, or result type may have
different guarantees and must not be presented as a one-to-one equivalent.

## Research and version record

The plans must link to concise, authoritative sources instead of copying large
documentation passages:

| Track | Primary anchors | Version or evidence limitation |
| --- | --- | --- |
| Python | Tutorial, Language Reference, PEP 8, typing, standard library | No repository version; local observation is Python 3.14.4 |
| Go | Tour/tutorials, Specification, modules, diagnostics, standard library | No local toolchain; current specification reports Go 1.27; Effective Go is historical and does not cover newer modules/generics fully |
| Rust | Rust Book, Reference, Cargo Book, API Guidelines, standard library | No local toolchain; current Book reports Rust 1.97 and Edition 2024 as its assumption |
| TypeScript | Handbook, TSConfig, declaration files, Node runtime documentation | No repository TypeScript version or local `tsc`; local Node observation is 24.21.0 |
| C++ | ISO C++ getting-started material, C++ Core Guidelines, cppreference | GCC 15.2 and Clang 21.1.8 are local observations; standard version and external tooling remain to be selected |

Documentation facts, local observations, design judgments, and open questions
must remain distinguishable in the plan text.

## Dependency-ordered steps

| Step | Target | Depends on | Commit subject |
| ---: | --- | --- | --- |
| 01 | `python/plans/01-infrastructure-learning-plan.md` | none | `WI-0003: 01: doc(python): add infrastructure learning plan` |
| 02 | `go/plans/01-infrastructure-learning-plan.md` | 01 | `WI-0003: 02: doc(go): add infrastructure learning plan` |
| 03 | `rust/plans/01-infrastructure-learning-plan.md` | 01 | `WI-0003: 03: doc(rust): add infrastructure learning plan` |
| 04 | `typescript/plans/01-infrastructure-learning-plan.md` | 01 | `WI-0003: 04: doc(typescript): add infrastructure learning plan` |
| 05 | `cpp/plans/01-infrastructure-learning-plan.md` | 01 | `WI-0003: 05: doc(cpp): add infrastructure learning plan` |

The steps are intentionally serial. Step 01 establishes the direct adaptation
of the historical Python draft; steps 02–05 use that shared behavioral baseline
but own separate language-specific files.

## Delivery contract

The root work-item documentation is committed first with the record-level
subject shown above. Each plan step is then processed in ascending order through
the bounded iterator with `--parallelism 1`, clean-workspace enforcement, and no
forced recovery. The worker owns only its selected plan and its selected step
handoff update.

After each iterator run, an advisor reviews the selected plan, implementation
commit, staged scope, validation output, and run receipt. The advisor appends an
`## Advisor review` section to that step record and commits the review
separately before the next dependent step or integration.

The exact implementation subject must match the selected step's declaration
byte-for-byte. Review commits use the repository's documentation-review
convention and do not change the implementation subject.

## Validation and evidence

Before the root record commit:

```bash
git check-ref-format --branch wi-0003-learning-plans
node /home/michael/src/github/language-practice/.cline/skills/language-practice-cline-iterate/scripts/work_item_preflight.mjs \
  --root /home/michael/src/github/language-practice \
  --work-item WI-0003 \
  --step 01 \
  --target python/plans/01-infrastructure-learning-plan.md
git diff --check
```

The preflight must be rerun for each selected step after the root record is
committed. Each plan step must additionally verify that:

- its target exists and is the only implementation target;
- all required shared sections and stages `00` through `19` are present and
  ordered;
- language-specific documentation anchors and version limitations are present;
- language guarantees are not overstated and external facilities are marked;
- the staged diff contains only the selected plan and selected step handoff; and
- `git diff --cached --check` passes.

The final review must inspect all five plans, all declared implementation
subjects, all advisor reviews, the final validation output, the commit graph, and
the clean worktree. No merge or push is authorized by this record.

## Not proven or blocked

- Plan content can establish scope, ordering, source references, and intended
  learning evidence only.
- No language runtime, compiler, test suite, benchmark, learner-outcome, or
  production evidence is claimed by this work item.
- Go, Rust, and TypeScript compiler/toolchain checks may remain unavailable in
  the current environment and must be reported rather than inferred.
- The first TypeScript iterator attempt failed with an OpenRouter `ECONNRESET`
  before writing repository files. A changed-prompt retry completed successfully;
  both receipts remain preserved outside the repository.

## Completion record

- **Completed status:** all five plan steps are proven for their scoped static
  documentation and Git provenance contracts.
- **Implementation commits:** `6700348`, `73c9078`, `bae1528`, `5fe6583`, and
  `96c3646`, each matching its declared step subject exactly.
- **Advisor review commits:** `063dc5a`, `7b815a8`, `8732cb1`, and the final
  review commit for step 05.
- **Integration:** no merge, rebase, or push performed; the topic branch remains
  `wi-0003-learning-plans`.