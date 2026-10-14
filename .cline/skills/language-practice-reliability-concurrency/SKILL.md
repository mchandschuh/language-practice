---
name: language-practice-reliability-concurrency
description: Teach safe concurrent and failure-prone design across Go, Python, Rust, and TypeScript. Use for threads, goroutines, async tasks, queues, retries, cancellation, timeouts, backpressure, resource ownership, and recovery.
compatibility: Designed for local exercises and bounded service simulations; it does not authorize live production changes or external mutations.
metadata:
  version: "1"
  scope: "concurrency and reliability"
---

# Concurrency and reliability

Treat concurrency and distributed behavior as explicit contracts. The agent must
not hide scheduling, ownership, failure, or cleanup behind a framework when the
exercise is intended to teach those concepts.

## Design checklist

For each exercise, identify:

- work ownership and lifecycle;
- shared state and synchronization boundary;
- cancellation and timeout propagation;
- queue capacity and backpressure behavior;
- error propagation and retry policy;
- idempotency and duplicate work;
- ordering guarantees;
- resource cleanup and shutdown behavior;
- observability needed to diagnose hangs, leaks, and overload.

Use bounded retries with jitter only when appropriate, distinguish transient
from permanent failures, and avoid retrying non-idempotent effects without a
deliberate design.

## Experiments and validation

Start with a sequential or single-worker baseline when it clarifies the concept.
Then vary scheduling, contention, cancellation timing, queue size, failures,
and shutdown order. Test both normal completion and rejection/abort paths.

Use the language's race detector, sanitizer, model-checking aid, async test
tool, or tracing facility when available and relevant. A test that passes once
does not prove the absence of a race or deadlock.

## Language distinctions

Explain the language's actual model rather than translating terminology from
another track:

- Go: goroutines, channels, `context`, synchronization, and race detection;
- Python: threads, processes, async tasks, the interpreter/runtime model, and
  the limits of concurrency for the selected workload;
- Rust: ownership, borrowing, `Send`/`Sync`, channels, locks, scoped tasks,
  and the selected async runtime;
- TypeScript: JavaScript event-loop behavior, promises, cancellation patterns,
  worker boundaries, and the host runtime's I/O model.

Do not imply that memory safety, type checking, or an event loop automatically
proves application-level correctness.