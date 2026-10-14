# Linux bootstrap contract

## Purpose and learner profile

This document is the shared contract for preparing the five language-practice
tracks on a Linux host. It is written for a learner who can use a terminal,
read shell commands, and inspect a diff, but who should understand an
installation choice before running it. The intended mode is **assist** for
repetitive, visible setup work and **collaborate** for choosing versions,
installation ownership, and dependency boundaries.

The learner runs a selected script explicitly from the repository root. The
iterator and reviewer may inspect and statically validate scripts, but never run
them or install software. A successful installation is not evidence of project
correctness, learning, production readiness, or support for an unrecorded
version.

## Host and Linux assumptions

The contract assumes a supported Linux distribution, a POSIX shell with Bash
available for the executable scripts, a working network only when the learner
chooses a download-based method, and enough disk space for the selected
toolchain. The exact distribution, architecture, shell, proxy, certificates,
and package repositories must be recorded by the learner when troubleshooting.
No repository support version is inferred from the host.

Host prerequisites are separate from a language toolchain. A distribution
package may require an administrator and `sudo`; a user-local installer or
version manager may instead write below the user's home directory. The learner
must choose one method, read its commands, confirm the destination and
permissions, and run it in their own terminal.

## State changes, permissions, and visible commands

Every command in a language script must be visible and classifiable before it is
run:

| Class | Typical state | Permission boundary |
| --- | --- | --- |
| Observe | `command --version`, `command -v`, environment inspection | read-only |
| User-local | a home-directory toolchain, virtual environment, project cache, or dependency directory | current user; no `sudo` |
| Project-local | manifests, lockfiles, generated build directories, or a language environment under the track | repository and explicitly named disposable paths |
| System-wide | distribution packages, `/usr`, `/opt`, shared PATH configuration | administrator only, with an explicit `sudo` command and user approval |
| Network | downloads, repository access, registry access, or package resolution | user chooses when to permit it; URL and source must be visible |

Scripts must not hide commands behind encoded payloads, remote pipes, implicit
profile edits, or an unannounced global install. `sudo` is never automatic or
authorized for the iterator. Shell-profile mutation is an optional, explicit
user action outside the iterator; prefer a command or documented session-local
PATH change that the learner can see. A script must fail clearly rather than
silently widening permissions.

## Version selection and source authority

Supported versions are repository decisions, not measurements of this host.
Before running a script, record the chosen language/toolchain version, target
architecture, installation source, and (where available) a release checksum or
lockfile. If the repository has not made a version decision, the script and this
contract must say so instead of inventing one. Distribution packages, official
vendor archives, user-local version managers, and project-local package
managers are alternatives, not interchangeable evidence.

Use primary documentation for the selected method and preserve the exact source
URL in the language script or its handoff. Useful official starting points are:

- Python `venv`: <https://docs.python.org/3/library/venv.html> and the Python
  Packaging User Guide;
- Go installation: <https://go.dev/doc/install>;
- Rust installation and toolchain management:
  <https://www.rust-lang.org/tools/install>;
- Node.js installation guidance:
  <https://nodejs.org/en/download/package-manager>;
- CMake tutorial and build concepts:
  <https://cmake.org/cmake/help/latest/guide/tutorial/index.html>.

These links are source authorities, not a repository version pin. Their current
web content can change; the repository records only deliberate decisions and
observations available at the time of the change.

## Per-language bootstrap matrix

Each row is a separate user-run executable target. Step 01 owns this shared
contract only; steps 02–06 own the corresponding scripts and must not be
substituted for one another.

| Track | Sole script target | Local project boundary | Minimum visible verification |
| --- | --- | --- | --- |
| Python | `./python/bootstrap.sh` | disposable `.venv` and declared Python metadata/dependencies | interpreter, `venv`, `pip`, imports, and one bounded project command |
| Go | `./go/bootstrap.sh` | module and explicitly named build/test output | `go version`, module inspection, and a bounded build or test |
| Rust | `./rust/bootstrap.sh` | Cargo manifest/lockfile and disposable `target/` | `rustc`, `cargo`, manifest/lockfile inspection, and a bounded check/test |
| TypeScript/Node | `./typescript/bootstrap.sh` | package manifest, lockfile, and disposable `node_modules` | Node/npm, package metadata, compiler invocation, and a bounded project command |
| C++ | `./cpp/bootstrap.sh` | source/build directory and declared CMake or compiler configuration | compiler, CMake/generator if used, configure/build, and a bounded test if present |

The matrix does not claim that any target currently exists or that a toolchain
is installed. Language steps must state their chosen version policy, package
and dependency policy, official sources, and exact commands without claiming
more than their evidence supports.

## Idempotence, failure, cleanup, and rollback

Scripts should be safe to inspect and should be rerunnable. They must detect
whether a prerequisite or project-local environment already exists, avoid
blindly overwriting user files, and make recreation an explicit choice. A
second run may verify an existing compatible state; it must not silently mix
versions or append duplicate PATH/profile entries. Generated environments,
build directories, caches, and downloaded archives must be named and excluded
from source control.

Failure boundaries must be stated before mutation. Commands should stop on
unexpected errors, report the failed operation and destination, and leave
already-installed user or system state for the learner to inspect. Rollback is
limited to artifacts the script created in its named project-local or
user-local scope. Cleanup must not remove unrelated home-directory, system,
repository, or credential state. System package removal and broad cleanup are
manual user decisions, not automatic recovery. If a partial installer cannot
be safely reversed, document that boundary rather than promising rollback.

## Verification and evidence tiers

Verification commands must be bounded, visible, and separated by what they
prove:

1. **Static evidence:** file ownership, executable permission, shell syntax,
   command visibility, declared paths, and `git diff --check`.
2. **Host observation:** the learner's command path and reported versions after
   installation. This is evidence about that host, not a repository support
   decision.
3. **Project evidence:** a local environment, manifest/lockfile, bounded
   configure/build/check/test, and expected output for the selected track.
4. **Learning evidence:** the learner can explain the choice, permissions,
   reproducibility inputs, failure boundary, and cleanup path. Automation cannot
   establish this tier.

The iterator may use tier 1 only. No tier proves successful installation on a
different Linux host, network availability, package authenticity beyond the
checks actually performed, future reproducibility, security, performance, or
learner understanding. Host-installation evidence is intentionally unproven
until a user runs a selected script outside this iterator.

## Troubleshooting and reporting contract

When a command fails, report the track and script, Linux distribution and
architecture, shell, exact visible command (redacting secrets), exit status,
relevant output, selected source/version, permission context, and whether any
files or packages were changed. Also report the last successful verification
tier. Do not paste credentials, tokens, private URLs, or complete environment
dumps. Resolve one boundary at a time: prerequisite, toolchain, project-local
environment, dependency resolution, then project verification.

## Ownership and static-validation rules

The five executable targets are exactly:

```text
./python/bootstrap.sh
./go/bootstrap.sh
./rust/bootstrap.sh
./typescript/bootstrap.sh
./cpp/bootstrap.sh
```

Each later step owns only its listed script plus its selected handoff notes.
Scripts may contain user-run installation, `sudo`, package-manager, and network
commands when necessary, but those commands must remain explicit, bounded, and
reviewable. Static validation may inspect permissions, Bash syntax, command
ordering, declared paths, and dangerous command visibility. It must not execute
a script, installer, package manager, network download, `sudo`, shell-profile
mutation, global install, or destructive cleanup.

## Non-goals and unproven claims

This contract does not select a Linux distribution, language version, compiler
standard, CMake generator, package registry, dependency set, lockfile format,
global installation policy, or shell profile. It does not add exercises,
solutions, application code, service dependencies, or production deployment.
It does not claim that the current host's observed versions are supported, that
all five toolchains are installed, or that any installer will work on another
host. Those decisions and limitations must remain explicit in the
language-specific handoffs.