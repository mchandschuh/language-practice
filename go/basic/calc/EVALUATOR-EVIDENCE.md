# Evaluator evidence report

## Question and scope

This report measures the existing sequential evaluator against the reusable,
type-dispatched channel evaluator for a fixed local workload. It is a bounded
learning experiment, not a capacity target or production performance claim.

The comparison changes one implementation variable at a time. The channel
benchmarks compare admission capacities `1` and `16`; both use the current
runtime's single worker. The sequential benchmark is the correctness oracle.

## Workload and method

The parsed corpus contains these expressions and inputs:

| Expression | Input |
| --- | --- |
| `x * x + 2` | `x=3` |
| `(9 - 4) / 5 + 2 * 3` | empty |
| `-x + 3` | `x=4` |
| `x++ + ++x` | `x=5` |

Expressions are parsed once before timing. The expected value for every corpus
entry is computed with `Evaluator{}.EvaluateWithEnvironment`. Each benchmark
iteration selects the next corpus entry in round-robin order. The sequential
path clones the input through the oracle API; the channel path submits the same
immutable AST and receives one terminal result, including its request-owned
environment. Results are checked against the precomputed oracle value.

`BenchmarkChannelEvaluator` measures one request from submission through result
receipt. `BenchmarkChannelBatch` submits a fixed batch of 16 requests and then
receives all results. The batch benchmark exposes bounded admission effects;
capacity `1` blocks later submissions while the single worker progresses,
whereas capacity `16` admits the whole batch. Neither measurement isolates
internal worker service time. No warmup phase is added beyond the Go benchmark
harness's normal calibration, and setup/teardown are outside the timed region.

All benchmarks call `ReportAllocs`. The raw output is intentionally not checked
in; this report records the bounded five-sample run below.

## Reproduction

Environment:

```text
Go: go1.26.8
OS/architecture: linux/arm64
Module: github.com/mchandschuh/language-practice/go
Worker count: 1
Queue capacities: 1 and 16
Batch size: 16
```

Command:

```bash
cd go
go test -bench 'Benchmark(SequentialOracle|ChannelEvaluator|ChannelBatch)' \
  -benchmem -benchtime=100ms -count=5 ./basic/calc/math
```

The bounded command completed successfully in 5.595 seconds. Representative
raw results from the five samples were:

| Benchmark | ns/op range | B/op | allocs/op |
| --- | ---: | ---: | ---: |
| SequentialOracle | 341.6–505.8 | 204 | 1 |
| ChannelEvaluator, capacity 1 | 6,222–16,740 | 2,860 | 35 |
| ChannelEvaluator, capacity 16 | 6,355–9,261 | 2,859–2,860 | 35 |
| ChannelBatch, capacity 1 | 110,203–265,609 | 45,761–45,763 | 572 |
| ChannelBatch, capacity 16 | 252,801–270,944 | 45,768–45,771 | 572 |

The ranges show substantial local timing noise, especially in the short
sequential benchmark and the capacity-1 batch run. Allocation counts were
stable. The channel path is materially more expensive for this one-worker,
small-expression workload; the channel abstraction does not imply a speedup.
The batch result also shows that a larger queue can increase total batch
completion time in this serial worker configuration because it allows all work
to queue before collection rather than applying admission backpressure.

## Correctness and validation evidence

The benchmark corpus's expected values are computed by the sequential oracle,
and every timed channel result is checked against that value. The existing
bounded reuse and lifecycle tests additionally cover immutable AST reuse,
request-local mutation isolation, cancellation, shutdown, and repeated close.

Commands run after the final benchmark/report edits:

```text
cd go && gofmt -w basic/calc/math/evaluator_benchmark_test.go: passed
cd go && go test ./basic/calc/math -count=1: passed
cd go && go test -bench 'Benchmark(SequentialOracle|ChannelEvaluator|ChannelBatch)' -benchmem -benchtime=100ms -count=5 ./basic/calc/math: passed
cd go && go test ./... -count=1: passed
cd go && go vet ./...: passed
cd go && go test -race ./basic/calc/math -count=1: passed
git diff --check: passed
```

## Limitations

- These are local measurements on one Linux/ARM64 host and one Go toolchain.
- The workload contains four small expressions and a fixed input distribution;
  it is not representative application traffic.
- The runtime currently owns one worker, so this report provides no evidence of
  scaling with worker count or parallel execution.
- Submission-to-result time includes channel, admission, scheduling, cloning,
  and result delivery overhead. The batch experiment shows queueing effects but
  cannot separate internal service time from all other costs.
- Five samples with 100 ms per benchmark provide bounded observations, not
  statistically robust latency distributions or tail guarantees.
- Allocation evidence is benchmark-harness evidence for these code paths, not a
  memory budget or deployment cost estimate.
- Passing tests and race checks cover the exercised inputs and schedules; they
  do not prove universal race freedom, leak freedom, portability, production
  capacity, or production readiness.

## Conclusion

The sequential evaluator remains the appropriate correctness baseline. For the
measured fixed workload, the channel evaluator preserves the oracle values but
adds substantial allocation and timing overhead, as expected from the current
single-worker channel architecture. Queue capacity changes admission behavior
and observed batch completion time, not the evaluator's semantics. No
optimization or production-capacity conclusion is justified by this report.