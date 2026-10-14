---
schema: language-practice-step/v1
id: WI-0004
step: "04"
status: planned
skill: language-practice-rust
depends_on: ["01"]
target: rust/bootstrap.sh
---

# WI-0004 step 04 — Document Rust bootstrap setup

## Objective

Write the executable `./rust/bootstrap.sh` for Linux Rust installation and
toolchain selection through rustup, Cargo project creation, components, targets,
lockfiles, verification, and cleanup.

## Commit subject

```text
WI-0004: 04: doc(rust): document bootstrap setup
```

## Required boundaries

- Write `./rust/bootstrap.sh` as the sole implementation target for this step;
  the user will inspect and run it explicitly from a terminal.
- Explain stable/channel selection, `rustup show`, toolchain/component/target
  inspection, and the consequences of changing the default toolchain.
- Explain `cargo new`, `Cargo.toml`, `Cargo.lock`, `cargo check`, `cargo test`,
  `cargo fmt`, and `cargo clippy` as separate checks.
- State when a project should pin a toolchain or edition and identify unresolved
  repository decisions explicitly.
- Record that rustc/Cargo are unavailable locally and the Rust Book's 1.97/
  Edition 2024 context is not an execution contract.
- Do not execute the script, run rustup, download components, create Cargo
  projects, or change the user's default toolchain while processing this step.

## Validation

```bash
test -x rust/bootstrap.sh
bash -n rust/bootstrap.sh
grep -q 'rustup show' rust/bootstrap.sh
grep -q 'cargo new' rust/bootstrap.sh
grep -q 'Cargo.toml\|cargo' rust/bootstrap.sh
grep -q 'Cargo.lock\|--locked' rust/bootstrap.sh
grep -q 'rustup\|curl\|apt\|package' rust/bootstrap.sh
git diff --check
```

The script must link to or cite rustup and the official Cargo Book in its
comments or the shared `BOOTSTRAP.md`. Static checks do not execute rustup or
prove toolchain, component, target, or Cargo behavior.

## Prerequisite

Step 01 must be proven and reviewed. No rustup installation, component download,
or default-toolchain mutation is authorized by this step.

## Agent handoff notes

Reserved for worker evidence, exact validation, and limitations.

## Advisor review

The advisor appends review rounds here after the worker completes.