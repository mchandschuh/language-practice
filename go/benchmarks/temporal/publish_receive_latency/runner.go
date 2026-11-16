package publishlatency

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	enums "go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/interceptor"
	"go.temporal.io/sdk/log"
	"go.temporal.io/sdk/worker"
)

// RunResult contains all raw records and the machine-readable summary for one
// unique run ID. Warm-ups remain in Records and are marked Warmup=true.
type RunResult struct {
	RunID       string
	Records     []Measurement
	Summaries   []ScenarioSummary
	Environment Environment
}

type harness struct {
	config                Config
	runID                 string
	client                client.Client
	registry              *TimestampRegistry
	workflowWorker        worker.Worker
	activityWorker        worker.Worker
	workflowQueue         string
	activityQueue         string
	activityCoordinatorID string
	signalWorkflowID      string
	signalWorkflowRunID   string
}

// Run connects to Temporal, starts and verifies workers, runs the selected
// scenario/mode matrix without overlap within each publication stream, and
// returns raw records plus summaries. It does not write files; callers can use
// Export to write after the run has ended.
func Run(ctx context.Context, config Config) (RunResult, error) {
	if err := config.validate(); err != nil {
		return RunResult{}, err
	}
	runID, err := newRunID()
	if err != nil {
		return RunResult{}, err
	}

	clientOptions := client.Options{
		HostPort:  config.HostPort,
		Namespace: config.Namespace,
		Identity:  "temporal-publish-receive-benchmark/" + runID,
		Logger:    discardLogger{},
	}
	temporalClient, err := client.DialContext(ctx, clientOptions)
	if err != nil {
		return RunResult{}, fmt.Errorf("dial Temporal: %w", err)
	}
	defer temporalClient.Close()

	registry := NewTimestampRegistry()
	workerInterceptor := &timingWorkerInterceptor{registry: registry}
	workflowOptions := worker.Options{
		MaxConcurrentWorkflowTaskPollers:       config.WorkflowPollers,
		MaxConcurrentWorkflowTaskExecutionSize: config.WorkflowExecutionConcurrency,
		Interceptors:                           []interceptor.WorkerInterceptor{workerInterceptor},
		EnableLoggingInReplay:                  false,
	}
	activityOptions := worker.Options{
		MaxConcurrentActivityTaskPollers:   config.ActivityPollers,
		MaxConcurrentActivityExecutionSize: config.ActivityExecutionConcurrency,
		Interceptors:                       []interceptor.WorkerInterceptor{workerInterceptor},
		EnableLoggingInReplay:              false,
		DisableWorkflowWorker:              true,
	}

	workflowQueue := config.WorkflowTaskQueue
	activityQueue := config.ActivityTaskQueue
	if workflowQueue == "" {
		workflowQueue = "publish-latency/" + runID + "/workflow"
	}
	if activityQueue == "" {
		activityQueue = "publish-latency/" + runID + "/activity"
	}

	workflowWorker := worker.New(temporalClient, workflowQueue, workflowOptions)
	workflowWorker.RegisterWorkflow(StartWorkflow)
	workflowWorker.RegisterWorkflow(ActivityCoordinatorWorkflow)
	workflowWorker.RegisterWorkflow(SignalWorkflow)
	activityWorker := worker.New(temporalClient, activityQueue, activityOptions)
	activityWorker.RegisterActivityWithOptions((&activities{registry: registry}).MinimalActivity, activity.RegisterOptions{Name: "MinimalActivity"})

	if err := workflowWorker.Start(); err != nil {
		return RunResult{}, fmt.Errorf("start workflow worker: %w", err)
	}
	defer workflowWorker.Stop()
	if err := activityWorker.Start(); err != nil {
		return RunResult{}, fmt.Errorf("start activity worker: %w", err)
	}
	defer activityWorker.Stop()

	h := &harness{
		config:         config,
		runID:          runID,
		client:         temporalClient,
		registry:       registry,
		workflowWorker: workflowWorker,
		activityWorker: activityWorker,
		workflowQueue:  workflowQueue,
		activityQueue:  activityQueue,
	}
	if err := h.waitReady(ctx); err != nil {
		return RunResult{}, err
	}
	if err := h.startLongLivedWorkflows(ctx); err != nil {
		return RunResult{}, err
	}
	defer h.cleanup(context.Background())

	scenarios := selectedScenarios(config.Scenario)
	modes := selectedModes(config.Mode)
	for _, scenario := range scenarios {
		for _, mode := range modes {
			if err := h.runWorkload(ctx, scenario, mode, true, config.WarmupIterations); err != nil {
				return RunResult{}, err
			}
			if err := h.runWorkload(ctx, scenario, mode, false, config.Iterations); err != nil {
				return RunResult{}, err
			}
		}
	}

	environment, err := captureEnvironment(ctx, temporalClient, config, runID, workflowQueue, activityQueue)
	if err != nil {
		return RunResult{}, err
	}
	records := registry.Records()
	return RunResult{
		RunID:       runID,
		Records:     records,
		Summaries:   BuildSummaries(records),
		Environment: environment,
	}, nil
}

func (h *harness) waitReady(ctx context.Context) error {
	readyCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	for {
		if _, err := h.client.CheckHealth(readyCtx, &client.CheckHealthRequest{}); err != nil {
			if readyCtx.Err() != nil {
				return fmt.Errorf("Temporal health readiness: %w", readyCtx.Err())
			}
			time.Sleep(25 * time.Millisecond)
			continue
		}
		workflowReady, workflowErr := h.taskQueueHasPoller(readyCtx, h.workflowQueue, enums.TASK_QUEUE_TYPE_WORKFLOW)
		activityReady, activityErr := h.taskQueueHasPoller(readyCtx, h.activityQueue, enums.TASK_QUEUE_TYPE_ACTIVITY)
		if workflowErr == nil && activityErr == nil && workflowReady && activityReady {
			return nil
		}
		if readyCtx.Err() != nil {
			return fmt.Errorf("worker poller readiness: %w (workflow=%v activity=%v)", readyCtx.Err(), workflowErr, activityErr)
		}
		time.Sleep(25 * time.Millisecond)
	}
}

func (h *harness) taskQueueHasPoller(ctx context.Context, taskQueue string, taskQueueType enums.TaskQueueType) (bool, error) {
	response, err := h.client.DescribeTaskQueue(ctx, taskQueue, taskQueueType)
	if err != nil {
		return false, err
	}
	return len(response.GetPollers()) > 0, nil
}

func (h *harness) startLongLivedWorkflows(ctx context.Context) error {
	activityWorkflowID := "publish-latency-activity-coordinator/" + h.runID
	activityRun, err := h.client.ExecuteWorkflow(ctx, workflowStartOptions(activityWorkflowID, h.workflowQueue), ActivityCoordinatorWorkflow)
	if err != nil {
		return fmt.Errorf("start Activity coordinator workflow: %w", err)
	}
	h.activityCoordinatorID = activityRun.GetID()
	if err := h.waitForWorkflowReady(ctx, activityRun.GetID(), activityRun.GetRunID()); err != nil {
		return fmt.Errorf("Activity coordinator readiness: %w", err)
	}

	signalWorkflowID := "publish-latency-signal-workflow/" + h.runID
	signalRun, err := h.client.ExecuteWorkflow(ctx, workflowStartOptions(signalWorkflowID, h.workflowQueue), SignalWorkflow)
	if err != nil {
		return fmt.Errorf("start signal workflow: %w", err)
	}
	h.signalWorkflowID = signalRun.GetID()
	h.signalWorkflowRunID = signalRun.GetRunID()
	if err := h.waitForWorkflowReady(ctx, signalRun.GetID(), signalRun.GetRunID()); err != nil {
		return fmt.Errorf("signal workflow readiness: %w", err)
	}
	return nil
}

func (h *harness) waitForWorkflowReady(ctx context.Context, workflowID, runID string) error {
	readyCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	for {
		_, err := h.client.QueryWorkflow(readyCtx, workflowID, runID, readyQuery)
		if err == nil {
			return nil
		}
		if readyCtx.Err() != nil {
			return readyCtx.Err()
		}
		time.Sleep(25 * time.Millisecond)
	}
}

func (h *harness) runWorkload(parent context.Context, scenario, mode string, warmup bool, iterations int) error {
	if iterations == 0 {
		return nil
	}
	if mode == ModeSequential || h.config.Concurrency == 1 {
		for i := 0; i < iterations; i++ {
			if err := h.runOne(parent, scenario, mode, warmup, i); err != nil {
				return err
			}
		}
		return nil
	}
	if mode == ModeConcurrent {
		for first := 0; first < iterations; first += h.config.Concurrency {
			last := first + h.config.Concurrency
			if last > iterations {
				last = iterations
			}
			var waitGroup sync.WaitGroup
			var firstErr error
			var errorMu sync.Mutex
			for iteration := first; iteration < last; iteration++ {
				waitGroup.Add(1)
				go func(iteration int) {
					defer waitGroup.Done()
					if err := h.runOne(parent, scenario, mode, warmup, iteration); err != nil {
						errorMu.Lock()
						if firstErr == nil {
							firstErr = err
						}
						errorMu.Unlock()
					}
				}(iteration)
			}
			waitGroup.Wait()
			if firstErr != nil {
				return firstErr
			}
		}
		return nil
	}

	// Sustained mode continuously feeds a bounded worker pool. Unlike the
	// concurrent mode's barrier between batches, the next operation is made
	// available as soon as a publication worker is free, which exposes backlog
	// and saturation effects.
	jobs := make(chan int)
	var waitGroup sync.WaitGroup
	workerCount := h.config.Concurrency
	if workerCount > iterations {
		workerCount = iterations
	}
	var firstErr error
	var errorMu sync.Mutex
	for workerIndex := 0; workerIndex < workerCount; workerIndex++ {
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			for iteration := range jobs {
				if err := h.runOne(parent, scenario, mode, warmup, iteration); err != nil {
					errorMu.Lock()
					if firstErr == nil {
						firstErr = err
					}
					errorMu.Unlock()
				}
			}
		}()
	}
	for iteration := 0; iteration < iterations; iteration++ {
		jobs <- iteration
	}
	close(jobs)
	waitGroup.Wait()
	return firstErr
}

func (h *harness) runOne(parent context.Context, scenario, mode string, warmup bool, iteration int) error {
	operationID := fmt.Sprintf("%s/%s/%s/%t/%06d", h.runID, scenario, mode, warmup, iteration)
	operation := Operation{RunID: h.runID, OperationID: operationID, Scenario: scenario, Mode: mode, Warmup: warmup}
	if err := h.registry.Register(operation); err != nil {
		return err
	}
	operationCtx, cancel := context.WithTimeout(parent, h.config.OperationTimeout)
	defer cancel()
	if err := h.registry.BeginPublish(operationID); err != nil {
		return err
	}
	var publishErr error
	switch scenario {
	case ScenarioStart:
		_, publishErr = h.client.ExecuteWorkflow(operationCtx, workflowStartOptions("publish-latency-start/"+operationID, h.workflowQueue), StartWorkflow, StartPayload{OperationID: operationID})
	case ScenarioActivity:
		publishErr = h.client.SignalWorkflow(operationCtx, h.activityCoordinatorID, "", activityTriggerSignal, ActivityPayload{OperationID: operationID, TaskQueue: h.activityQueue})
	case ScenarioSignal:
		publishErr = h.client.SignalWorkflow(operationCtx, h.signalWorkflowID, h.signalWorkflowRunID, signalReceiveSignal, SignalPayload{OperationID: operationID})
	default:
		return fmt.Errorf("unsupported scenario %q", scenario)
	}
	if err := h.registry.FinishPublish(operationID, publishErr); err != nil {
		return err
	}
	if publishErr != nil {
		_ = h.registry.MarkFailure(operationID, publishErr)
		return nil
	}
	if err := h.registry.WaitReceive(operationCtx, operationID); err != nil {
		_ = h.registry.MarkTimeout(operationID)
		return nil
	}
	if scenario != ScenarioSignal {
		if err := h.registry.WaitComplete(operationCtx, operationID); err != nil {
			_ = h.registry.MarkTimeout(operationID)
			return nil
		}
	}
	return h.registry.MarkSuccess(operationID)
}

func (h *harness) cleanup(ctx context.Context) {
	if h.activityCoordinatorID != "" {
		_ = h.client.TerminateWorkflow(ctx, h.activityCoordinatorID, "", "benchmark cleanup")
	}
	if h.signalWorkflowID != "" {
		_ = h.client.TerminateWorkflow(ctx, h.signalWorkflowID, h.signalWorkflowRunID, "benchmark cleanup")
	}
}

func selectedScenarios(value string) []string {
	switch value {
	case "start":
		return []string{ScenarioStart}
	case "activity":
		return []string{ScenarioActivity}
	case "signal":
		return []string{ScenarioSignal}
	case ScenarioStart, ScenarioActivity, ScenarioSignal:
		return []string{value}
	default:
		return []string{ScenarioStart, ScenarioActivity, ScenarioSignal}
	}
}

func selectedModes(value string) []string {
	switch value {
	case ModeSequential, ModeConcurrent, ModeSustained:
		return []string{value}
	default:
		return []string{ModeSequential, ModeConcurrent, ModeSustained}
	}
}

func newRunID() (string, error) {
	bytes := make([]byte, 12)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("generate run ID: %w", err)
	}
	return "run-" + hex.EncodeToString(bytes), nil
}

var _ log.Logger = discardLogger{}

type discardLogger struct{}

func (discardLogger) Debug(string, ...any) {}
func (discardLogger) Info(string, ...any)  {}
func (discardLogger) Warn(string, ...any)  {}
func (discardLogger) Error(string, ...any) {}
