package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	benchmark "github.com/mchandschuh/language-practice/go/benchmarks/temporal/publish_receive_latency"
)

func main() {
	defaults := benchmark.DefaultConfig()
	config := defaults
	flag.StringVar(&config.HostPort, "host-port", defaults.HostPort, "Temporal frontend gRPC host:port")
	flag.StringVar(&config.Namespace, "namespace", defaults.Namespace, "Temporal namespace")
	flag.IntVar(&config.Iterations, "iterations", defaults.Iterations, "measured iterations per scenario and mode")
	flag.IntVar(&config.WarmupIterations, "warmup-iterations", defaults.WarmupIterations, "warm-up iterations per scenario and mode")
	flag.IntVar(&config.Concurrency, "concurrency", defaults.Concurrency, "maximum in-flight publications for concurrent and sustained modes")
	flag.DurationVar(&config.OperationTimeout, "operation-timeout", defaults.OperationTimeout, "per-operation receive/completion timeout")
	flag.StringVar(&config.WorkflowTaskQueue, "workflow-task-queue", "", "workflow task queue; empty creates a run-specific queue")
	flag.StringVar(&config.ActivityTaskQueue, "activity-task-queue", "", "activity task queue; empty creates a run-specific queue")
	flag.IntVar(&config.WorkflowPollers, "workflow-pollers", defaults.WorkflowPollers, "workflow task pollers; Temporal SDK requires 0/default or at least 2")
	flag.IntVar(&config.ActivityPollers, "activity-pollers", defaults.ActivityPollers, "activity task pollers; 0 uses SDK default")
	flag.IntVar(&config.WorkflowExecutionConcurrency, "workflow-execution-concurrency", defaults.WorkflowExecutionConcurrency, "workflow task execution slots; 0 uses SDK default")
	flag.IntVar(&config.ActivityExecutionConcurrency, "activity-execution-concurrency", defaults.ActivityExecutionConcurrency, "activity execution slots; 0 uses SDK default")
	flag.StringVar(&config.Scenario, "scenario", defaults.Scenario, "all, start, activity, signal, or full scenario name")
	flag.StringVar(&config.Mode, "mode", defaults.Mode, "all, sequential, concurrent, or sustained")
	flag.StringVar(&config.OutputDir, "output-dir", defaults.OutputDir, "directory for run-id.csv and run-id.json")
	flag.StringVar(&config.PersistenceBackend, "persistence-backend", defaults.PersistenceBackend, "descriptive backend label; never changes server settings")
	flag.StringVar(&config.DatabaseDurability, "database-durability", defaults.DatabaseDurability, "descriptive durability label; never changes server settings")
	flag.StringVar(&config.TaskQueueFairness, "task-queue-fairness", defaults.TaskQueueFairness, "descriptive task-queue fairness label; benchmark does not override fairness")
	flag.Parse()

	ctx := context.Background()
	result, err := benchmark.Run(ctx, config)
	if err != nil {
		log.Fatal(err)
	}
	csvPath, jsonPath, err := benchmark.Export(result, config.OutputDir)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("run_id=%s\nraw_csv=%s\nsummary_json=%s\n", result.RunID, csvPath, jsonPath)
	printSummary(result)
}

func printSummary(result benchmark.RunResult) {
	for _, summary := range result.Summaries {
		if summary.Warmup {
			continue
		}
		stats := summary.PublishToReceive
		fmt.Printf("scenario=%s mode=%s success=%d failed=%d timed_out=%d p50_us=%.3f p99_us=%.3f\n",
			summary.Scenario,
			summary.Mode,
			summary.SuccessfulIterations,
			summary.FailedIterations,
			summary.TimedOutIterations,
			stats.MedianUS,
			stats.P99US,
		)
	}
}
