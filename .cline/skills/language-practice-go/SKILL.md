---
name: language-practice-go
description: Teach idiomatic Go conventions and behavior for exercises and services in the Go track. Use when working with Go source, packages, modules, tests, errors, interfaces, goroutines, channels, contexts, or Go diagnostics.
compatibility: Requires an installed Go toolchain for execution; use the repository's declared Go version when one is added.
metadata:
  version: "1"
  scope: "go track"
---

# Go track conventions

Use official Go documentation and standard-library examples as the primary
authority. Begin with the Go version and module configuration actually present.

## Coding conventions

- Run `gofmt`/the repository formatter rather than hand-formatting.
- Use short, clear names and idiomatic package boundaries; avoid package names
  that repeat the package at every call site.
- Prefer simple composition and small interfaces defined near their consumers.
- Return errors explicitly, add useful context without destroying the cause,
  and decide deliberately between wrapping, comparison, and sentinel errors.
- Keep exported identifiers documented when they form a public package API.
- Use modules and the `go` command for dependency changes; inspect the resulting
  `go.mod`/`go.sum` and do not add dependencies without a reason.

## Learning focus

Handwrite small examples for slices/maps, interfaces, method sets, error
handling, goroutines, channels, and `context`. Make nil behavior, copying,
escape, blocking, cancellation, and package initialization observable.

Do not translate class hierarchies or exception-heavy designs directly into Go.
Explain when a concrete type, interface, function, channel, or `context` is the
right abstraction.

## Validation and diagnostics

Use the repository's documented commands, commonly including `go test ./...`
and `go vet ./...` where applicable. Use the race detector for race-relevant
tests. Use benchmarks, `pprof`, execution tracing, and runtime statistics only
for a stated diagnostic question; report their perturbation and scope.

Do not claim that a successful build proves goroutine shutdown, race freedom,
or production performance.

## Official anchors

- `https://go.dev/doc/effective_go`
- `https://go.dev/doc/modules/managing-dependencies`
- `https://go.dev/doc/diagnostics`
- the Go specification and standard-library documentation
