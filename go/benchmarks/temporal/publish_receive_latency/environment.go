package publishlatency

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	enums "go.temporal.io/api/enums/v1"
	taskqueue "go.temporal.io/api/taskqueue/v1"
	workflowservice "go.temporal.io/api/workflowservice/v1"
	"go.temporal.io/sdk/client"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

// Environment records the conditions that can materially affect a local
// Temporal latency run. Persistence and database durability are intentionally
// descriptive inputs: Temporal does not provide a portable API for discovering
// the database's fsync policy, and this benchmark never changes it.
type Environment struct {
	CapturedAt            string `json:"captured_at"`
	TemporalServerVersion string `json:"temporal_server_version"`
	TemporalSDKVersion    string `json:"temporal_sdk_version"`
	PersistenceBackend    string `json:"persistence_backend"`
	DatabaseDurability    string `json:"database_durability"`

	OperatingSystem       string `json:"operating_system"`
	Architecture          string `json:"architecture"`
	CPUModel              string `json:"cpu_model"`
	LogicalProcessorCount int    `json:"logical_processor_count"`
	AvailableMemoryBytes  uint64 `json:"available_memory_bytes"`
	RuntimeVersion        string `json:"runtime_version"`

	Worker     WorkerConfiguration    `json:"worker_configuration"`
	TaskQueues TaskQueueConfiguration `json:"task_queues"`
}

// WorkerConfiguration captures the configured SDK worker knobs, including the
// SDK constraint that poller/slot values of one are invalid.
type WorkerConfiguration struct {
	WorkflowPollers              int `json:"workflow_pollers"`
	ActivityPollers              int `json:"activity_pollers"`
	WorkflowExecutionConcurrency int `json:"workflow_execution_concurrency"`
	ActivityExecutionConcurrency int `json:"activity_execution_concurrency"`
	PublicationConcurrency       int `json:"publication_concurrency"`
}

// TaskQueueConfiguration includes configured queues and observed poller
// identities. Fairness is not overridden by the benchmark.
type TaskQueueConfiguration struct {
	Workflow         string   `json:"workflow"`
	Activity         string   `json:"activity"`
	Fairness         string   `json:"fairness"`
	WorkflowFairness string   `json:"workflow_fairness"`
	ActivityFairness string   `json:"activity_fairness"`
	WorkflowPollers  []Poller `json:"workflow_pollers"`
	ActivityPollers  []Poller `json:"activity_pollers"`
}

// Poller is a stable export form for the server's DescribeTaskQueue response.
type Poller struct {
	Identity       string `json:"identity"`
	LastAccessTime string `json:"last_access_time,omitempty"`
}

func captureEnvironment(ctx context.Context, temporalClient client.Client, config Config, runID, workflowQueue, activityQueue string) (Environment, error) {
	serverVersion := "unknown/not reported"
	response, err := temporalClient.WorkflowService().GetSystemInfo(ctx, &workflowservice.GetSystemInfoRequest{})
	if err == nil && response.GetServerVersion() != "" {
		serverVersion = response.GetServerVersion()
	}
	if err != nil {
		return Environment{}, fmt.Errorf("get Temporal server version: %w", err)
	}

	workflowPollers, err := describePollers(ctx, temporalClient, workflowQueue, enums.TASK_QUEUE_TYPE_WORKFLOW)
	if err != nil {
		return Environment{}, fmt.Errorf("describe workflow task queue: %w", err)
	}
	activityPollers, err := describePollers(ctx, temporalClient, activityQueue, enums.TASK_QUEUE_TYPE_ACTIVITY)
	if err != nil {
		return Environment{}, fmt.Errorf("describe activity task queue: %w", err)
	}

	fairness := config.TaskQueueFairness
	workflowFairness := queueConfig(ctx, temporalClient, config.Namespace, workflowQueue, enums.TASK_QUEUE_TYPE_WORKFLOW)
	activityFairness := queueConfig(ctx, temporalClient, config.Namespace, activityQueue, enums.TASK_QUEUE_TYPE_ACTIVITY)
	if workflowFairness != "" {
		fairness = workflowFairness
	}
	return Environment{
		CapturedAt:            time.Now().UTC().Format(time.RFC3339Nano),
		TemporalServerVersion: serverVersion,
		TemporalSDKVersion:    sdkVersion(),
		PersistenceBackend:    config.PersistenceBackend,
		DatabaseDurability:    config.DatabaseDurability,
		OperatingSystem:       operatingSystem(),
		Architecture:          runtime.GOARCH,
		CPUModel:              cpuModel(),
		LogicalProcessorCount: runtime.NumCPU(),
		AvailableMemoryBytes:  availableMemoryBytes(),
		RuntimeVersion:        runtime.Version(),
		Worker: WorkerConfiguration{
			WorkflowPollers:              config.WorkflowPollers,
			ActivityPollers:              config.ActivityPollers,
			WorkflowExecutionConcurrency: config.WorkflowExecutionConcurrency,
			ActivityExecutionConcurrency: config.ActivityExecutionConcurrency,
			PublicationConcurrency:       config.Concurrency,
		},
		TaskQueues: TaskQueueConfiguration{
			Workflow:         workflowQueue,
			Activity:         activityQueue,
			Fairness:         fairness,
			WorkflowFairness: workflowFairness,
			ActivityFairness: activityFairness,
			WorkflowPollers:  workflowPollers,
			ActivityPollers:  activityPollers,
		},
	}, nil
}

func queueConfig(ctx context.Context, temporalClient client.Client, namespace, taskQueue string, taskQueueType enums.TaskQueueType) string {
	response, err := temporalClient.WorkflowService().DescribeTaskQueue(ctx, &workflowservice.DescribeTaskQueueRequest{
		Namespace:     namespace,
		TaskQueue:     &taskqueue.TaskQueue{Name: taskQueue, Kind: enums.TASK_QUEUE_KIND_NORMAL},
		TaskQueueType: taskQueueType,
		ReportConfig:  true,
	})
	if err != nil || response.GetConfig() == nil {
		return ""
	}
	config := response.GetConfig()
	if len(config.GetFairnessWeightOverrides()) > 0 || config.GetFairnessKeysRateLimitDefault() != nil {
		return "server-configured fairness settings present"
	}
	return "server-reported default fairness configuration"
}

func describePollers(ctx context.Context, temporalClient client.Client, taskQueue string, taskQueueType enums.TaskQueueType) ([]Poller, error) {
	response, err := temporalClient.DescribeTaskQueue(ctx, taskQueue, taskQueueType)
	if err != nil {
		return nil, err
	}
	pollers := make([]Poller, 0, len(response.GetPollers()))
	for _, poller := range response.GetPollers() {
		pollers = append(pollers, Poller{
			Identity:       poller.GetIdentity(),
			LastAccessTime: timestampString(poller.GetLastAccessTime()),
		})
	}
	return pollers, nil
}

func timestampString(timestamp *timestamppb.Timestamp) string {
	if timestamp == nil {
		return ""
	}
	return timestamp.AsTime().UTC().Format(time.RFC3339Nano)
}

func sdkVersion() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "unknown/not reported"
	}
	for _, dependency := range info.Deps {
		if dependency.Path == "go.temporal.io/sdk" {
			return dependency.Version
		}
	}
	return "unknown/not reported"
}

func operatingSystem() string {
	if data, err := os.ReadFile("/etc/os-release"); err == nil {
		for _, line := range strings.Split(string(data), "\n") {
			if strings.HasPrefix(line, "PRETTY_NAME=") {
				return strings.Trim(strings.TrimPrefix(line, "PRETTY_NAME="), "\"")
			}
		}
	}
	return runtime.GOOS
}

func cpuModel() string {
	data, err := os.ReadFile("/proc/cpuinfo")
	if err != nil {
		return "unknown/not reported"
	}
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "model name") || strings.HasPrefix(line, "Model") {
			if _, value, ok := strings.Cut(line, ":"); ok {
				return strings.TrimSpace(value)
			}
		}
	}
	return "unknown/not reported"
}

func availableMemoryBytes() uint64 {
	data, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return 0
	}
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "MemAvailable:") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				value, err := strconv.ParseUint(fields[1], 10, 64)
				if err == nil {
					return value * 1024
				}
			}
		}
	}
	return 0
}
