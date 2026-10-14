---
name: language-practice-cross-language-comparison
description: Compare one programming concept across Go, Python, Rust, and TypeScript while preserving each language's real semantics and conventions. Use for parallel exercises, porting analysis, tradeoff discussions, or language selection.
compatibility: Designed for the four language tracks in this repository and authoritative documentation for each selected version.
metadata:
  version: "1"
  scope: "cross-language pedagogy"
---

# Cross-language comparison

Compare concepts, not spellings. Start with one behavioral contract and show how
each language models it, what it guarantees, and what remains a runtime or
library concern.

## Comparison method

1. State the shared problem and acceptance behavior.
2. Identify the language feature or standard-library abstraction used in each
   track.
3. Implement or inspect the smallest equivalent example.
4. Compare types, memory/resource ownership, errors, concurrency, tooling, and
   runtime behavior.
5. Run equivalent tests and document meaningful differences in setup and
   evidence.
6. Record where the comparison is intentionally not equivalent.

Do not force a one-to-one mapping. An interface, protocol, trait, duck-typed
object, promise, channel, task, exception, or result type may have different
semantics even when the names sound similar.

## Learning rules

Use `handwrite` for the core examples when the purpose is intuition. Use
`collaborate` to identify translations and mismatches. Use `assist` for
repetitive harnesses only after the shared contract is understood.

Require a short explanation of:

- static versus dynamic checking;
- memory and resource ownership;
- error representation and propagation;
- concurrency and cancellation model;
- package/module boundaries;
- tooling and diagnostic workflow;
- performance and portability limits.

Research each language's official documentation rather than importing a style
rule from another track.