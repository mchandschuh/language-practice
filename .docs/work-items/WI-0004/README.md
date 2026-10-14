---
schema: language-practice-work-item/v1
id: WI-0004
status: IN_PROGRESS
type: feature
title: Document reproducible Linux toolchain bootstrap
repository: /home/michael/src/github/language-practice
branch: wi-0004-toolchain-bootstrap
authority: local work-item records under this directory
---

# WI-0004 — Document reproducible Linux toolchain bootstrap

## Authority and status

- **Raw status:** `IN_PROGRESS`
- **Derived status:** `planned`
- **Authority:** this README and the ordered `steps/` records
- **Repository:** `/home/michael/src/github/language-practice`
- **Branch:** `wi-0004-toolchain-bootstrap`
- **Baseline:** `63095e9` (`Merge branch wi-0003-learning-plans`)
- **Predecessor:** WI-0003, merged by `63095e9` with `--no-ff`
- **Root record commit:** `WI-0004: doc(repo): prepare toolchain bootstrap plan`
- **Implementation commits:** one commit per ordered setup-documentation step
  after this directory is committed

## Scope

Document a safe, reproducible Linux bootstrap path for the five learning tracks
created by WI-0003:

- Python;
- Go;
- Rust;
- TypeScript/Node; and
- C++ under `cpp/`.

The bootstrap work must distinguish four layers:

1. host prerequisites and system packages;
2. language toolchain installation and version selection;
3. project-local environments, manifests, dependencies, and lockfiles; and
4. verification, troubleshooting, cleanup, and rollback.

The setup must be understandable by a learner and reproducible by another
developer without silently changing the global system. Each language script
must state its supported-version decision, installation alternatives, exact
verification commands, project bootstrap commands, dependency/reproducibility
policy, permissions boundary, failure modes, cleanup path, and evidence limits.

Each language-specific step owns one executable script at the repository-relative
path `./{lang}/bootstrap.sh`:

- `./python/bootstrap.sh`;
- `./go/bootstrap.sh`;
- `./rust/bootstrap.sh`;
- `./typescript/bootstrap.sh`; and
- `./cpp/bootstrap.sh`.

The scripts are intended for the user to inspect and run explicitly from a
terminal, with `sudo` when the selected installation method requires it. The
iterator worker and advisor must never execute these scripts, invoke `sudo`, or
install a toolchain while processing this work item.

The shared guide must explain how to choose between distribution packages,
user-local installers, version managers, and project-local package managers. It
must not pretend that a successful installer command proves compiler support,
project correctness, or learner understanding.

## Non-goals

- Do not execute the bootstrap scripts, installers, `sudo`, privileged package
  operations, network downloads, `curl | sh`, global package installs, or
  destructive cleanup while preparing, iterating, or reviewing this work item.
- The scripts may contain the explicit commands needed for a user-run bootstrap,
  but those commands must be visible, bounded, non-interactive where practical,
  and never hidden behind an implicit worker or advisor action.
- Do not mutate the host toolchains, shell profiles, PATH, package registries,
  credentials, or global configuration as part of the work-item documentation.
- Do not invent supported language versions from local observations. Select and
  document versions deliberately, or mark the decision open.
- Do not add application code, exercises, solutions, service dependencies,
  production configuration, or deployment automation.
- Do not commit virtual environments, compiled binaries, package caches,
  `node_modules`, Cargo `target`, Go build output, private keys, tokens, or
  downloaded installers.
- Do not require a single installation method when Linux distribution and
  permission constraints materially differ; document safe alternatives and the
  verification common to them.
- Do not claim that installing a toolchain proves a learning plan, runtime
  behavior, production readiness, security readiness, or performance.
- Do not push or merge the WI-0004 topic branch as part of this preparation.

## Safety and learning contract

The documentation mode is `assist` for repetitive setup commands and
configuration scaffolding, with a `collaborate` checkpoint for choosing versions,
installation ownership, and dependency boundaries. The learner must inspect and
understand each command before running it. Core verification and troubleshooting
remain observable learning checkpoints.

Every setup section must identify:

- prerequisite knowledge and host assumptions;
- the learner's decisions and commands;
- whether a command changes user-local or system state;
- required permissions and the smallest safe scope;
- expected observable output and version checks;
- failure modes and recovery without deleting unrelated user state;
- how to reproduce the environment later; and
- what the check does not prove.

Prefer dry-run, version-reporting, temporary directories, user-local prefixes,
project-local environments, and explicit confirmation before mutation. Never
hide a privileged or destructive action behind a convenience script.

## Bootstrap matrix

The final documents must preserve these language distinctions:

| Track | Toolchain boundary | Project/dependency boundary | Required verification |
| --- | --- | --- | --- |
| Python | Existing interpreter plus an isolated `venv`; interpreter selection is separate from package installation | `venv`, `pip`, project metadata, requirements or lockfile policy | `python --version`, `python -m venv`, `python -m pip --version`, package/import checks |
| Go | Go distribution and PATH/GOPATH tool locations | `go.mod`, `go.sum`, `go mod tidy`, module-local commands | `go version`, `go env`, `go test`, `go vet`, and module inspection |
| Rust | `rustup` toolchain/channel and components | Cargo manifest, `Cargo.lock`, profiles, workspaces | `rustc --version`, `cargo --version`, `rustup show`, `cargo check/test/fmt/clippy` |
| TypeScript/Node | Node/npm installation and version-management choice | `package.json`, local TypeScript dev dependency, lockfile, `tsconfig.json` | `node --version`, `npm --version`, local `npx tsc --version`, type-check/test/runtime checks |
| C++ | Distribution/compiler choice, standard version, CMake/build generator | CMake project, presets, out-of-source build directory, test runner | `g++ --version` or `clang++ --version`, `cmake --version`, configure/build/test checks |

The matrix is a design contract, not evidence that any tool is installed.

## Research and version record

The setup documents must cite and distinguish official facts from repository
decisions and local observations:

| Topic | Primary source | Current fit | Limitation |
| --- | --- | --- | --- |
| Python environments/packages | Python `venv`; Python Packaging User Guide | Current docs; local Python is 3.14.4 | Repository Python version and dependency model are undeclared |
| Go installation/modules | Go install guide; Go module tutorial | Current official docs; no local Go toolchain | Installation paths and permissions differ between archive and distribution methods |
| Rust installation/Cargo | rustup; Cargo Book | Current Rust tooling docs; no local Rust toolchain | Channel, target, component, and edition decisions remain open |
| Node/npm/TypeScript | Node download page; npm install guide; TypeScript download/setup docs | Node docs report v24.21.0 LTS; local Node is 24.21.0/npm 11.19.0 | TypeScript version and package manager are undeclared; local `tsc` is unavailable |
| C++/CMake | ISO C++ getting started; GCC/Clang docs; CMake tutorial and presets | Local GCC 15.2, Clang 21.1.8; no CMake | C++ standard, generator, and distribution policy remain open |

Research notes must not turn the current online latest release into an implicit
repository contract. Each implementation step must either select a version and
record why, or leave the choice explicitly open for operator decision.

## Dependency-ordered steps

| Step | Target | Depends on | Commit subject |
| ---: | --- | --- | --- |
| 01 | `BOOTSTRAP.md` | none | `WI-0004: 01: doc(repo): add bootstrap contract` |
| 02 | `python/bootstrap.sh` | 01 | `WI-0004: 02: doc(python): document bootstrap setup` |
| 03 | `go/bootstrap.sh` | 01 | `WI-0004: 03: doc(go): document bootstrap setup` |
| 04 | `rust/bootstrap.sh` | 01 | `WI-0004: 04: doc(rust): document bootstrap setup` |
| 05 | `typescript/bootstrap.sh` | 01 | `WI-0004: 05: doc(typescript): document bootstrap setup` |
| 06 | `cpp/bootstrap.sh` | 01 | `WI-0004: 06: doc(cpp): document bootstrap setup` |

The steps are intentionally processed serially even though steps 02–06 each
depend only on step 01. Step 01 defines shared safety, version, permissions,
verification, and reproducibility language. Steps 02–06 adapt that contract to
one track each and own separate files. No step authorizes execution of the
installation commands it documents.

## Delivery contract

The root work-item documentation is committed first with the record-level
subject shown above. Each setup-documentation step is then processed in ascending
order through the bounded iterator with `--parallelism 1`, clean-workspace
enforcement, and no forced recovery.

Each worker owns only its selected script target and selected step handoff. The
advisor reviews commands for unsafe implicit mutation, version ambiguity,
unsupported claims, missing rollback/verification, non-idempotent behavior, and
scope expansion. Review notes are appended and committed before the next
dependent step or integration.

This work item includes writing executable bootstrap scripts, but not executing
them. Static validation may inspect syntax, permissions, command ordering, and
declared paths. System installation, privileged mutation, and external downloads
remain user-run actions requiring explicit terminal execution outside the
iterator.

## Validation and evidence

Before the root record commit:

```bash
git check-ref-format --branch wi-0004-toolchain-bootstrap
git diff --check
```

After the root commit, each selected step must pass the WI preflight and verify:

- the target and selected step handoff are the only changed paths;
- the target script is executable, has valid Bash syntax, and contains the
  language-specific commands required by its step;
- version choices, host assumptions, permissions, and mutation boundaries are
  explicit;
- verification commands are bounded and distinguish version reporting,
  project checks, and runtime evidence;
- scripts are safe to inspect without running and do not hide `sudo`, network,
  or destructive actions;
- cleanup and rollback do not remove unrelated user state;
- official sources are linked and version limitations are recorded; and
- `git diff --cached --check` passes.

The final review must inspect all six setup documents, all declared subjects,
all advisor reviews, the final validation output, the commit graph, and the
clean worktree. It must report any commands not executed because they would
install software, require privileges, access the network, or mutate global state.

## Not proven or blocked

- No bootstrap script, toolchain installation, or dependency download is
  executed by the iterator or advisor, so host-installation evidence is
  intentionally unavailable.
- Current local observations are Python 3.14.4, Node 24.21.0/npm 11.19.0,
  GCC 15.2, Clang 21.1.8, and unavailable Go/rustc/Cargo/tsc/CMake.
- Supported versions, package-manager policy, C++ standard, CMake generator,
  and exact dependency lockfile policy must be decided in the implementation
  steps or explicitly left open.
- Static script checks can establish syntax, scope, and command visibility; they
  cannot prove successful installation on another developer's Linux host.