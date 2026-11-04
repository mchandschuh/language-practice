# Temporal publish-to-worker-receive latency benchmark

This benchmark measures actual elapsed time in a single Go process containing
the Temporal Client, Workflow Worker, and Activity Worker. It targets a local
Temporal installation and is intended to expose Temporal scheduling and task
delivery overhead, not to estimate performance from Temporal metrics.

The implementation is under:

```text
/home/michael/src/github/language-practice/go/benchmarks/temporal/publish_receive_latency
```

It uses the repository's Go module and Temporal Go SDK `v1.50.0`; no new
dependency or external infrastructure is introduced.

## Run against a local Temporal installation

From the repository root, start the local development service in another
terminal. The repository's helper is also available:

```bash
cd /home/michael/src/github/language-practice
source ./init
temporal_dev_start
```

The equivalent Temporal CLI command is:

```bash
temporal server start-dev
```

The default frontend is `localhost:7233`, Namespace `default`, and the default
development server uses an in-memory persistence store. Do not infer that this
is durable or representative of a production database. If you use a persistent
local database, configure it deliberately (for example with
`temporal_dev_start --persistent` or `--db-filename`) and pass descriptive
metadata to the benchmark. The benchmark never changes database durability.

Run a small smoke test first:

```bash
cd /home/michael/src/github/language-practice/go
go run ./benchmarks/temporal/publish_receive_latency/cmd/temporal-publish-receive \
  -iterations 10 \
  -warmup-iterations 2 \
  -scenario all \
  -mode sequential \
  -output-dir /tmp/temporal-publish-receive
```

Run the documented default workload:

```bash
cd /home/michael/src/github/language-practice/go
go run ./benchmarks/temporal/publish_receive_latency/cmd/temporal-publish-receive \
  -output-dir ./benchmarks/temporal/results
```

The default is 500 warm-up and 10,000 measured iterations per selected
scenario and publication mode. The default `-mode all` runs sequential,
controlled concurrent, and sustained workloads. This can create substantial
Workflow histories; start with a small smoke test and choose a deliberate local
Namespace/retention policy.

## Configuration examples

### Workflow Start only, sequential baseline

```bash
go run ./benchmarks/temporal/publish_receive_latency/cmd/temporal-publish-receive \
  -scenario start \
  -mode sequential \
  -iterations 10000 \
  -warmup-iterations 500 \
  -workflow-pollers 2 \
  -output-dir ./benchmarks/temporal/results
```

### Controlled concurrent and sustained publication

```bash
go run ./benchmarks/temporal/publish_receive_latency/cmd/temporal-publish-receive \
  -scenario activity \
  -mode all \
  -concurrency 8 \
  -workflow-pollers 4 \
  -activity-pollers 4 \
  -workflow-execution-concurrency 8 \
  -activity-execution-concurrency 8 \
  -output-dir ./benchmarks/temporal/results
```

### Persistent local service metadata

```bash
go run ./benchmarks/temporal/publish_receive_latency/cmd/temporal-publish-receive \
  -iterations 1000 \
  -warmup-iterations 100 \
  -persistence-backend "Temporal CLI dev sqlite" \
  -database-durability "configured externally; not changed by benchmark" \
  -task-queue-fairness "server default" \
  -output-dir ./benchmarks/temporal/results
```

`-workflow-pollers 1` and
`-workflow-execution-concurrency 1` are rejected because the selected SDK
requires zero/default or at least two for those settings. Activity poller and
execution values of zero use SDK defaults. The CLI creates run-specific queues
when queue flags are empty, preventing stale tasks or another local worker from
being accidentally included.

## Measurement boundaries

All timestamps are captured with `time.Now()` in one process. Go's
`time.Time` carries a monotonic reading, and durations are calculated before
timestamps are serialized to wall-clock strings. There is no subtraction of
independently captured monotonic clocks.

For every operation, the registry is populated before publication:

| Timestamp | Boundary |
| --- | --- |
| `T0` | Immediately before the Client calls the Temporal publish API. |
| `T1` | Immediately after that API returns. API response time is not task delivery time. |
| `T2` | First correlated receive boundary in the Worker. Workflow Start uses the Workflow inbound interceptor before user Workflow code. Activity uses the first instruction in `MinimalActivity`. Signal uses the signal inbound interceptor before the SDK places the signal on the Workflow channel. |
| `T3` | End of the minimal Workflow or Activity handler where applicable. Idle Signal Workflow has no useful completion boundary. |
| `ScheduleInitiated` | SDK Workflow outbound `ExecuteActivity` initiation for the Activity scenario. This is not the server's durable `ActivityTaskScheduled` event. |

Reported durations are:

- **Publish-to-receive:** `T2 - T0`, the primary metric.
- **Publish API latency:** `T1 - T0`, only the Client publish call.
- **Post-acknowledgment delivery:** signed `T2 - T1`; a negative result means the Worker boundary raced ahead of the API response.
- **Total execution:** `T3 - T0` when `T3` exists.
- **Schedule-to-receive:** `T2 - ScheduleInitiated`, an additional Activity diagnostic.

### Scenario A: Workflow Start -> Workflow Task Receive

The publisher calls `ExecuteWorkflow` for a new Workflow ID. `T2` is recorded
when the first Workflow Task reaches the SDK Workflow inbound interceptor. This
includes service persistence, History transfer, Matching, poll response
handling, SDK decoding, and interceptor dispatch up to that boundary. It does
not include the user's Workflow body before the interceptor; `T3` is recorded
after the empty Workflow returns.

### Scenario B: Activity Scheduling -> Activity Receive

An existing coordinator Workflow receives a correlated trigger Signal. The
publisher's `T0/T1` surround that Signal API, which is the publication used to
ask the existing Workflow to schedule work. `ScheduleInitiated` is recorded
later, inside the Workflow SDK outbound `ExecuteActivity` call. The server's
durable scheduling event and transfer to the Activity Task Queue occur between
that point and Activity `T2`. Therefore the primary `T2-T0` includes the trigger
Signal's scheduling path; `ScheduleToReceive` is the closer SDK-scheduling to
Activity-handler metric. The Activity handler's `T2` is an approximation of task
receipt: SDK poll response processing and argument decoding occur before the
handler is called.

### Scenario C: Signal -> Workflow Receive

The publisher signals an already-running idle Workflow. `T2` is recorded for
the correlated signal at the SDK Workflow signal inbound boundary. This
scenario necessarily includes a new Workflow Task being scheduled, polled,
replayed as needed, and dispatched. It is reported separately from Workflow
Start and Activity delivery.

Workflow replay is excluded from unique receive counts with
`workflow.IsReplaying(ctx)`, and the registry also accepts only the first
non-replay timestamp for an operation. Replay exclusion is especially relevant
to signal handling and the Activity coordinator's evolving history.

## Outputs

After the workload completes, the command writes:

- `<run-id>.csv`: raw per-operation records, including run ID, operation ID,
  scenario, warm-up marker, status, error, all timestamps, and durations in
  nanoseconds/microseconds/milliseconds.
- `<run-id>.json`: environment metadata and summaries.

Warm-up records are retained in the CSV and summarized separately with
`warmup: true`; they are never silently discarded. Successful, failed, and
timed-out iterations are counted separately. Percentiles are linear
interpolation over sorted successful samples and include minimum, mean, p50,
p90, p95, p99, p99.9, maximum, and population standard deviation in both
microseconds and milliseconds.

The environment document records Temporal Server version, SDK version, OS, CPU
model, logical processor count, available memory, Go runtime, worker settings,
queue names, observed pollers, and server-reported task-queue fairness config
when available. Persistence backend and database durability are labels supplied
by the operator because they cannot be discovered portably from the SDK. No
durability or fairness settings are changed automatically.

## Methodology and interpretation

- `sequential` has one operation in flight and is the baseline for intrinsic
  delivery overhead.
- `concurrent` uses the configured number of in-flight operations.
- `sustained` continuously feeds the same bounded concurrent publisher for the
  full iteration count. Its latency includes any backlog created by the chosen
  worker capacity; compare it with sequential and inspect timeout/failure
  counts before attributing tails to Temporal internals.
- No sleep is inserted between measured iterations. Readiness polling and
  worker warm-up happen before measured operations.
- Temporal SDK logging is discarded by the live harness so console I/O does not
  occur during measurement. Export and the short final CLI summary happen only
  after all operations finish.
- Worker readiness requires a healthy Temporal frontend, observed pollers from
  `DescribeTaskQueue`, and successful ready queries on the two long-lived
  Workflows. Worker startup alone is not treated as proof of an outstanding
  long poll.
- The harness discards Temporal SDK logs during the workload. This reduces
  console-I/O perturbation; it also means fatal errors are surfaced through the
  returned run error and the per-operation CSV rather than live log lines.
- A single process removes cross-process clock-origin ambiguity. Running the
  Client or Worker separately would require an explicitly synchronized clock
  strategy and is not supported by this implementation.
- These are local same-host measurements. They do not establish production
  capacity, cross-host network latency, availability, or database durability.

## Validation

From `/home/michael/src/github/language-practice/go`:

```bash
gofmt -w benchmarks/temporal/publish_receive_latency
go test ./benchmarks/temporal/publish_receive_latency -count=1
go test -race ./benchmarks/temporal/publish_receive_latency -count=1
go vet ./benchmarks/temporal/publish_receive_latency/...
go test ./... -count=1
```

The package tests cover concurrent correlation, signed delivery-before-ack
handling, replay/duplicate suppression, timeout/error classification,
percentiles, warm-up separation, configuration validation, and CSV/JSON export.
They do not claim a Temporal performance result and do not replace a live
server run.