---
name: language-practice-bootstrap
description: Initialize or extend the language-practice repository without inventing conventions, versions, or learning objectives. Use when setting up the Go, Python, Rust, or TypeScript tracks, adding toolchain configuration, or establishing repository commands.
compatibility: Designed for the language-practice repository and a coding agent with read-only inspection, file editing, and non-interactive command execution.
metadata:
  version: "1"
  scope: "repository initialization"
---

# Bootstrap the learning repository

Use this skill only for repository setup. It does not authorize implementing
the exercises themselves.

## Before editing

1. Inspect the working tree and preserve any existing user changes.
2. Inspect all four track directories and repository-level files.
3. Identify installed toolchains, but do not treat installed versions as the
   repository's supported versions until they are deliberately documented.
4. Research current official setup guidance for any toolchain being added.
5. State the intended learning outcome, supported command, and explicit
   non-goals before adding configuration.

## Setup rules

- Keep each language track independently runnable where practical.
- Do not add a framework, dependency, formatter, linter, or test runner merely
  because it is popular. Prefer the language toolchain and an existing local
  convention.
- Do not generate exercise solutions as part of scaffolding.
- Keep generated output, caches, virtual environments, binaries, secrets, and
  downloaded data out of version control.
- Make commands deterministic, bounded, and non-interactive.
- Prefer configuration that makes learning-relevant behavior visible instead of
  hiding it behind a large framework.

## Completion gate

Before reporting bootstrap complete:

- inspect every created file;
- run the narrowest available format, type-check, build, or test command;
- verify that the command is documented and reproducible;
- run `git diff --check`;
- report what the setup proves and what remains unproven.

Do not claim that a toolchain is supported because only installation succeeded.