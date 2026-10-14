---
name: language-practice-python
description: Teach readable, idiomatic Python conventions and runtime behavior for exercises and services in the Python track. Use when working with Python modules, data structures, exceptions, typing, iterators, generators, packages, tests, or async code.
compatibility: Requires an installed Python 3 toolchain; use the repository's declared version and environment manager when one is added.
metadata:
  version: "1"
  scope: "python track"
---

# Python track conventions

Use the Python version, project configuration, and environment actually present.
PEP 8 and the official tutorial guide defaults, but project-local consistency
takes precedence when deliberately documented.

## Coding conventions

- Prefer readable names, four-space indentation, explicit imports, small
  functions, and docstrings for public behavior.
- Use Python data structures and comprehensions when they clarify the operation;
  do not compress code merely to be clever.
- Understand mutability, identity, aliasing, scope, default arguments, and
  iteration before introducing abstractions.
- Catch specific exceptions, preserve useful context, and avoid broad catches
  that hide defects. Use cleanup constructs for resources.
- Treat type hints as information for tools, not runtime enforcement. Use them
  to make contracts clearer without claiming they replace validation.
- Use virtual environments and the repository's package manager. Do not modify
  the global interpreter as an implicit setup step.

## Learning focus

Handwrite examples for data structures, functions, modules, exceptions, classes,
iterators, generators, and context managers. Make laziness, mutation, object
identity, exception flow, and floating-point limitations observable.

For async or threaded code, explain the selected runtime model and what work it
can actually parallelize. Do not import assumptions from Go or Rust.

## Validation

Use the repository's documented formatter, type checker, test runner, and
package commands. Keep tests deterministic and isolated. Use profiling and
performance measurement only after defining a workload and baseline.

## Official anchors

- `https://docs.python.org/3/tutorial/`
- `https://peps.python.org/pep-0008/`
- `https://docs.python.org/3/library/typing.html`
- the Python language reference and standard-library documentation
