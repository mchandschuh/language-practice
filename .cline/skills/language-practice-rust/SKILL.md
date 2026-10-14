---
name: language-practice-rust
description: Teach idiomatic Rust conventions and semantics for exercises and services in the Rust track. Use when working with ownership, borrowing, lifetimes, enums, pattern matching, traits, errors, crates, tests, concurrency, or async code.
compatibility: Requires an installed Rust toolchain and Cargo; use the repository's declared edition and toolchain when one is added.
metadata:
  version: "1"
  scope: "rust track"
---

# Rust track conventions

Use the Rust edition and toolchain configured by the repository. Treat compiler
diagnostics as a central learning resource, not as an obstacle to silence.

## Coding conventions

- Use `rustfmt`, Cargo conventions, clear names, and small modules.
- Model domain states with enums and types when that makes invalid states harder
  to represent; use pattern matching deliberately and check exhaustiveness.
- Understand ownership, moves, borrows, mutability, lifetimes, and drop before
  adding clones, reference-counting, or interior mutability.
- Prefer `Result`/`Option` and meaningful error context. Do not use `unwrap` or
  `expect` as a production error strategy; when used pedagogically, explain
  the precondition and failure behavior.
- Use traits for behavior and generics when they improve the contract; avoid
  abstraction that hides the ownership or dispatch lesson.
- Keep unsafe code out of ordinary exercises unless the exercise explicitly
  studies its contract and validation.

## Learning focus

Handwrite examples for ownership, borrowing, lifetimes, enums, traits, error
handling, iterators, threads, channels, and synchronization. Ask the learner to
predict compiler errors and then explain the minimal correction.

For async code, identify the executor/runtime, cancellation behavior, `Send`
and `Sync` requirements, pinning constraints where relevant, and shutdown path.
The language guarantees memory safety properties; it does not guarantee that a
program is deadlock-free or logically correct.

## Validation

Use Cargo's documented commands, commonly `cargo fmt --check`, `cargo check`,
`cargo test`, and `cargo clippy` where the repository adopts them. Run focused
tests and benchmarks against the final source snapshot. Do not treat compiler
success as proof of runtime behavior or performance.

## Official anchors

- `https://doc.rust-lang.org/book/`
- `https://rust-lang.github.io/api-guidelines/`
- the Rust Reference, Cargo documentation, and standard-library documentation
