---
name: language-practice-temporal
description: Teach Temporal durable execution in Go through deterministic Workflows, Activities, Workers, Clients, testing, and local operations.
compatibility: Designed for the language-practice Go track and local Temporal development; live services, downloads, credentials, and production changes require explicit user action.
metadata:
  version: "1"
  scope: "Temporal Go learning track"
  importance: "high"
---

# Temporal Go learning track

Use this skill when working on `/home/michael/src/github/language-practice/go/temporal`, explaining Temporal concepts, debugging the local Temporal application, designing tests, or preparing the local CLI prerequisite.

The goal is durable understanding of Temporal's execution model, not merely a
working Hello World. Keep the learner writing the core Workflow and Activity
code in `handwrite` or `collaborate` mode. Use `assist` for repetitive local
tooling and configuration, while explaining the protocol, failure model, and
resource lifecycle that the tooling hides.

## Current repository facts

Observed on 2026-10-26:

| Item | Current value | Boundary |
| --- | --- | --- |
| Go module | `github.com/mchandschuh/language-practice/go` | repository fact |
| Go directive | `1.26.8` | repository fact, not a cross-machine support promise |
| Temporal Go SDK | `go.temporal.io/sdk v1.50.0` | current module dependency; verify with `go.mod` before relying on newer APIs |
| Local service | `localhost:7233` | Temporal CLI development-server default |
| Local UI | `http://localhost:8233` | Temporal CLI development-server default |
| CLI installer | `/home/michael/src/shell/init/install_temporal` | separate shell repository; user-run only |
| CLI installer release | Temporal CLI `v1.9.1` | pinned installer observation; refresh deliberately when upgrading |

The SDK and CLI versions are independent. Do not infer that the CLI release,
Temporal Server version, and Go SDK version are interchangeable. Check the
repository manifest and official release notes when changing any of them.

## Core execution model

Temporal persists Workflow Event History in a Temporal Service. A Worker polls
a Task Queue, replays Workflow history to reconstruct Workflow state, and runs
Activity code when the history schedules Activity work. A Client starts and
interacts with Workflow Executions; it is not the Worker and it is not the
Workflow implementation.

```text
Client / starter
        │ start, signal, query, get result
        ▼
Temporal Service ──durable Event History──┐
        │                                  │
        │ Workflow and Activity tasks     │ replay/recovery
        ▼                                  ▼
Worker polling a Task Queue ───────> Workflow + Activity code
```

Keep these boundaries explicit:

- **Workflow Definition:** deterministic orchestration and durable state
  transitions.
- **Activity Definition:** ordinary Go code for external effects and
  failure-prone work.
- **Worker:** long-lived process that registers types and polls a Task Queue.
- **Client:** process-level connection used to start or interact with
  executions.
- **Temporal Service:** durable coordination and Event History, local or
  remote.

## Non-negotiable Workflow rules

Workflow code may be replayed many times. The result of a Workflow Task must be
consistent with its Event History. Treat determinism as a correctness
requirement, not a performance detail.

Use Temporal's Workflow APIs inside Workflow code:

| Do not use directly in a Workflow | Use instead |
| --- | --- |
| `context.Context` | `workflow.Context` |
| `time.Now()` | `workflow.Now(ctx)` |
| `time.Sleep()` | `workflow.Sleep(ctx, duration)` |
| standard `log` or `fmt.Println` | `workflow.GetLogger(ctx)` |
| native `go` | `workflow.Go()` |
| native channels | `workflow.Channel` |
| native `select` | `workflow.Selector` |
| random/UUID generation | a recorded Side Effect or an Activity |
| network, database, filesystem, or other external I/O | an Activity |

Native goroutines and channels are not allowed in Workflow code. Do not use
replay detection to change business behavior. Keep Workflow parameters,
returns, Activity parameters, and Activity results serializable.

Use `workflow.ExecuteActivity` and a Future to schedule and await an Activity.
Set the required Activity timeout deliberately; the Hello World example uses a
`StartToCloseTimeout` of ten seconds. Return errors explicitly and propagate
them unless the exercise is specifically teaching recovery or compensation.

## Activity rules

Activities are the boundary for non-deterministic or external work: HTTP calls,
database operations, email, filesystem access, random values, and calls to
other services. Activity code receives standard `context.Context`, and should
use `activity.GetLogger(ctx)` for Activity logs.

Design Activities to be idempotent or to tolerate duplicate attempts. Temporal
may retry an Activity after a timeout, Worker failure, or transient error. A
successful local call is not evidence that the operation is safe to repeat.

Teach these controls in this order:

1. `StartToCloseTimeout` for a bounded Activity attempt.
2. Retry policy and the difference between retryable and non-retryable errors.
3. `ScheduleToStartTimeout`, `ScheduleToCloseTimeout`, and heartbeat timeout
   when queue delay, total duration, or long-running progress matters.
4. Cancellation and heartbeat details for long-running work.

Keep payloads small. Activity inputs and outputs are recorded in Event History;
large histories increase replay and recovery cost. Do not put credentials,
tokens, or sensitive data in ordinary payloads or logs. Introduce a Data
Converter, encryption/codec, or external storage only when the exercise has a
clear reason and the learner understands the data boundary.

## Worker and Task Queue rules

Create one heavyweight Client and Worker per process unless the exercise is
explicitly testing a different lifecycle. Register every Workflow and Activity
type the Worker may receive before calling `Run`.

```go
c, err := client.Dial(client.Options{})
if err != nil {
	log.Fatalln("Unable to create client", err)
}
defer c.Close()

w := worker.New(c, hello.TaskQueue, worker.Options{})
w.RegisterWorkflow(hello.HelloWorkflow)
w.RegisterActivity(hello.HelloActivity)

if err := w.Run(worker.InterruptCh()); err != nil {
	log.Fatalln("Unable to start worker", err)
}
```

`Run(worker.InterruptCh())` blocks and shuts down on an interrupt, allowing the
Worker to stop polling and finish in-flight work according to its options. For
more controlled lifecycle tests, use `Start` and `Stop` explicitly.

A Task Queue is routing, not a type registry. All Workers polling the same
queue must register compatible Workflow and Activity types. A starter and
Worker must use exactly the same queue string. A task can be delivered to any
Worker polling that queue; a Worker that did not register the received type
fails the task.

## Client and Workflow Execution rules

The Client connects to the Temporal Service and should be created once per
process and closed on exit. It must not be initialized or used inside a
Workflow. It is valid for an Activity to use a Client when that is the intended
external interaction.

For a starter, make the routing and identity visible:

```go
options := client.StartWorkflowOptions{
	ID:        "hello-workflow",
	TaskQueue: hello.TaskQueue,
}

workflowRun, err := c.ExecuteWorkflow(
	context.Background(),
	options,
	hello.HelloWorkflow,
	"World",
)
```

Then use the returned `WorkflowRun`:

```go
log.Println("Started workflow", "WorkflowID", workflowRun.GetID(), "RunID", workflowRun.GetRunID())

var result string
if err := workflowRun.Get(context.Background(), &result); err != nil {
	return err
}
```

A Workflow ID is an application-level identity, such as an order number. A Run
ID identifies one run in a Workflow Execution chain. Only one execution with a
given Workflow ID may be open in a Namespace at a time. Do not use a Run ID for
business identity. Teach duplicate-start behavior explicitly before relying on
the default conflict/reuse policy; set the appropriate `StartWorkflowOptions`
policy when a duplicate is an error or when an existing execution should be
used.

## Testing ladder

Every test should state what it proves and what it does not prove.

### 1. In-memory Workflow and Activity tests

Use `testsuite.WorkflowTestSuite` and a `TestWorkflowEnvironment` or
`TestActivityEnvironment` for deterministic, fast tests of Workflow logic and
Activity behavior. Register real or mocked Activities intentionally. These
tests do not prove Client connectivity, Worker registration in a real process,
Task Queue routing, serialization against a live Service, or shutdown behavior.

Current repository examples:

```bash
cd /home/michael/src/github/language-practice/go
go test ./temporal/hello/... -count=1
go vet ./temporal/hello/...
```

### 2. Local integration test

Run a disposable local Temporal Service, a real Worker, and a Client starter.
This proves process wiring, Client connectivity, registration, Task Queue
routing, serialization, Activity execution, and result retrieval. It is slower
and requires an external local dependency, so do not make ordinary unit tests
depend on it.

### 3. Replay testing

Workflow replay checks that current Workflow code remains compatible with a
representative Event History. Use `worker.WorkflowReplayer` and histories
downloaded from a controlled local or test Namespace. Replay is especially
important after changing Workflow orchestration, Activity ordering, timers,
signals, or child Workflows.

### 4. End-to-end testing

Use end-to-end tests sparingly for the complete Service/Worker/Client contract.
Keep them bounded, isolated, and disposable. The official Go SDK testing guide
generally favors integration tests for most Worker-backed behavior.

## Versioning and change safety

Workflow code can outlive the process and deployment that started it. Changing
Workflow orchestration can cause a non-determinism failure when old Event
History is replayed. Do not casually reorder or remove calls such as Activities,
Timers, Child Workflows, Signals, or Side Effects.

For beginner work, first learn to recognize the determinism boundary and run
in-memory tests. Defer `workflow.GetVersion` patching and Worker Deployment
Versioning until there is a concrete compatibility scenario. When the topic is
introduced, teach replay and an explicit migration plan together; do not hide a
Workflow change behind a green compile.

## Local CLI and service operations

The separate shell repository owns the user-run installer:

```text
/home/michael/src/shell/init/install_temporal
```

The installer is intentionally not part of the Go module and is not
automatically added to `install_all`. It installs the pinned Temporal CLI
release into a user-local directory, verifies official checksums, does not edit
shell profiles, and does not invoke `sudo`. Read the script before running it.

After installation, start the disposable local Service explicitly:

```bash
temporal server start-dev
```

The development Service normally listens on `localhost:7233`, exposes the Web
UI at `http://localhost:8233`, creates the default Namespace, and uses an
in-memory database. To preserve local data across restarts, use a deliberate
database file such as:

```bash
temporal server start-dev --db-filename temporal.db
```

The Service process is blocking; stop it with `Ctrl-C`. Do not treat a local
development server as production infrastructure. Do not connect exercises to
Temporal Cloud or another remote Namespace without explicit configuration,
credentials handling, and scope.

Useful local inspection commands include:

```bash
temporal workflow list
temporal workflow describe --workflow-id hello-workflow
temporal workflow show --workflow-id hello-workflow
temporal operator namespace list
```

Use the repository's `init` helper from the language-practice root after the
Service is running:

```bash
source ./init
run hello/worker
run hello/starter
```

The Worker is long-running. The starter should submit a Workflow and exit after
retrieving its result. Run them in separate terminals.

## Failure-oriented debugging checklist

When a local execution fails, classify the boundary before changing code:

| Symptom | First check |
| --- | --- |
| Client connection refused | Is `temporal server start-dev` running on `localhost:7233`? |
| Workflow remains pending | Do starter and Worker use the same exact Task Queue? |
| Unregistered type failure | Did the Worker register the Workflow and every Activity type? |
| Activity timeout | Is the timeout appropriate, and is the Activity idempotent under retry? |
| Duplicate Workflow start | What Workflow ID conflict/reuse policy was intended? |
| Non-determinism error | Did Workflow orchestration change, and can the old Event History replay? |
| Missing or invalid result | Is the result serializable and decoded into the matching Go type? |
| Unexpected duplicate side effect | Was external work placed in an Activity and made idempotent? |
| Repeated logs | Is Workflow logging using `workflow.GetLogger` rather than standard logging? |

Capture the exact command, Namespace, Task Queue, Workflow ID, Run ID, error
type, and relevant CLI history. Do not paste credentials, payload secrets, or
unbounded histories into logs or exercise files.

## Suggested learning progression

1. **Hello World:** one Activity, one Workflow, one Worker, one starter, one
   Task Queue, and in-memory tests.
2. **Execution visibility:** Workflow ID, Run ID, CLI list/describe/show, UI,
   and the Event History sequence.
3. **Failures:** Activity timeout, retry policy, application versus timeout and
   cancellation errors, and idempotency.
4. **Multiple Activities:** sequential and deterministic parallel execution;
   use `workflow.Go`, `workflow.Channel`, and `workflow.Selector` only after the
   sequential version is understood.
5. **Messages:** Signals, Queries, Updates, and cancellation with explicit
   lifecycle tests.
6. **Durable time:** timers, sleeps, schedules, and Continue-As-New while
   observing Event History growth.
7. **Composition:** Child Workflows and external interactions.
8. **Change safety:** replay tests, patching, and Worker Deployment Versioning.
9. **Operations:** structured logging, metrics/tracing, graceful shutdown,
   Worker tuning, namespaces, and only then remote/cloud configuration.

Do not skip directly from Hello World to production deployment, Cloud
credentials, encryption, custom payload codecs, or performance tuning.

## Research record

| Topic | Source | Claim used | Version fit | Limitation |
| --- | --- | --- | --- | --- |
| Go SDK map | [Go SDK developer guide](https://docs.temporal.io/develop/go) | Official links for setup, Workflows, Activities, Workers, Client, testing, and examples | Current docs consulted 2026-10-26 | Pages can change paths; follow the linked current page |
| Local setup | [Go local setup](https://docs.temporal.io/develop/go/set-up-your-local-go) | SDK setup, `temporal server start-dev`, `localhost:7233`, UI `8233`, three-process Hello World | Current docs; tutorial text may mention older Go versions | Repository uses Go 1.26.8; verify local toolchain separately |
| Workflow rules | [Workflow basics](https://docs.temporal.io/develop/go/workflows/basics) | `workflow.Context`, serializable data, deterministic APIs, replay-safe logging and time | Current Go SDK docs | Advanced determinism cases require replay experiments |
| Activity execution | [Activity execution](https://docs.temporal.io/develop/go/activities/execution) | Required timeout, Future/Get, retries, idempotency, Event History cost | Current Go SDK docs | Timeout/retry values are application-specific |
| Activity design | [Activity basics](https://docs.temporal.io/develop/go/activities/basics) | External work belongs in Activities; parameters/results are serialized | Current Go SDK docs | Payload limits and codecs need a focused exercise |
| Worker lifecycle | [Run a Worker](https://docs.temporal.io/develop/go/workers/run-worker-process) | Client/Worker creation, registration, Task Queue routing, interrupt shutdown | Current Go SDK docs | Production Worker tuning is deferred |
| Client | [Temporal Client](https://docs.temporal.io/develop/go/client/temporal-client) | `Dial`, `ExecuteWorkflow`, `WorkflowRun`, IDs, result retrieval, config boundaries | Current Go SDK docs; environment config support noted from SDK v1.28.0 | Cloud authentication is out of current scope |
| IDs | [Workflow ID and Run ID](https://docs.temporal.io/workflow-execution/workflowid-runid) | Application identity, run identity, duplicate/conflict/reuse semantics | Current platform docs | Exact policy depends on Namespace retention and application contract |
| Testing | [Go testing](https://docs.temporal.io/develop/go/best-practices/testing-suite) | In-memory tests, integration tests, end-to-end tests, replay | Current Go SDK docs | Live server tests require an explicit disposable service |
| Version safety | [Go versioning](https://docs.temporal.io/develop/go/workflows/versioning) | Long-lived histories require replay-compatible changes and versioning | Current Go SDK docs | Defer patching/deployment versioning until needed |
| CLI | [Temporal CLI](https://docs.temporal.io/cli) and [CLI basics](https://docs.temporal.io/cli/common-operations) | Local development server and inspection commands | Current CLI docs | Installer pins a separately observed CLI release |
| Official examples | [samples-go Hello World](https://github.com/temporalio/samples-go/tree/main/helloworld) | Current Worker/Starter/Workflow/Activity/test/replay layout | Current repository | Sample may evolve independently of this learning repo |
| Official code examples | [Go sandbox](https://github.com/temporalio/documentation-sdk-code-examples/tree/main/go-sandbox) | Documentation-linked focused examples | Current docs examples | Examples are references, not a repository dependency |
| Historical template | [Hello World template](https://github.com/temporalio/hello-world-project-template-go) | Earlier simple project layout | Archived; old SDK (`v1.17.0`) | Do not use its versions or Docker setup as current authority |
| SDK source | [Temporal Go SDK](https://github.com/temporalio/sdk-go) and [API docs](https://pkg.go.dev/go.temporal.io/sdk) | Release/API source and `workflowcheck` direction | Current release observed `v1.50.0` | Confirm APIs against the repository's selected module version |

## What validation proves

`go test ./...` and `go vet ./...` prove the current Go packages compile and the
repository tests/vet checks pass. In-memory Temporal tests prove selected
Workflow/Activity behavior under the SDK test environment. A successful Worker
build does not prove a live Service is reachable, Task Queue routing works, a
Workflow is deterministic under all histories, external effects are idempotent,
or shutdown is race-free.

Live local execution proves only the selected local Service/Worker/Client path.
It does not prove production capacity, availability, security, Cloud behavior,
or compatibility with an untested Server/SDK version pair.