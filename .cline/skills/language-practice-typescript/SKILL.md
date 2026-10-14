---
name: language-practice-typescript
description: Teach TypeScript type modeling, strictness, modules, JavaScript runtime behavior, and Node tooling for exercises in the TypeScript track. Use when working with TypeScript source, compiler configuration, declarations, async code, tests, or package modules.
compatibility: Requires a TypeScript-capable JavaScript runtime and the repository's declared package/tool versions when they are added.
metadata:
  version: "1"
  scope: "typescript track"
---

# TypeScript track conventions

Inspect `package.json`, lockfiles, `tsconfig.json`, module settings, and runtime
before choosing commands. TypeScript's static model and JavaScript's runtime are
separate authorities.

## Coding conventions

- Prefer strict compiler settings when the repository can adopt them, and
  understand which checks each option enables rather than copying a preset
  blindly.
- Model valid states with unions, discriminated unions, interfaces, and
  generics. Use narrowing and exhaustive checks instead of unchecked casts.
- Treat `any`, non-null assertions, assertions, and `@ts-ignore` as explicit
  escape hatches that require a reason and a boundary test.
- Keep modules and imports consistent with the configured ESM/CommonJS/runtime
  model. A successful type check does not prove that the runtime resolves or
  loads the module.
- Understand that `.d.ts` files describe an interface; they do not validate
  JavaScript behavior at runtime.
- Validate untrusted data at runtime even when it has a TypeScript type.

## Learning focus

Handwrite examples for everyday types, functions, object types, narrowing,
generics, discriminated unions, modules, and declaration boundaries. Make the
difference between compile-time erasure, runtime JavaScript semantics, and
external validation visible.

For asynchronous code, explain promises, event-loop scheduling, rejection
handling, cancellation conventions, and the host runtime's I/O behavior. Do not
describe TypeScript's types as a concurrency or memory-safety guarantee.

## Validation

Use the repository's package scripts and compiler configuration. Run the
formatter, type checker, tests, and a fresh-process/runtime check when the
change affects imports, emitted code, or package behavior. Keep dependencies
deliberate and lockfile changes reviewable.

## Official anchors

- `https://www.typescriptlang.org/docs/handbook/`
- `https://www.typescriptlang.org/tsconfig`
- `https://www.typescriptlang.org/docs/handbook/2/narrowing.html`
- `https://www.typescriptlang.org/docs/handbook/declaration-files/introduction.html`
- the JavaScript runtime's official module and API documentation