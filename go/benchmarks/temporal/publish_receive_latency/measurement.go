// Package publishlatency implements a live Temporal publish-to-worker-receive
// latency benchmark.
//
// The package deliberately keeps the timing registry in the process that owns
// both the Temporal client and the workers. A time.Time captured by time.Now
// contains Go's monotonic clock reading, so durations are calculated before
// export while that clock origin is still shared by every participant.
package publishlatency

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"
)

const (
	// ScenarioStart measures Workflow Start -> first Workflow Task dispatch.
	ScenarioStart = "workflow-start-to-workflow-task-receive"
	// ScenarioActivity measures a running workflow scheduling an Activity -> Activity dispatch.
	ScenarioActivity = "activity-scheduling-to-activity-receive"
	// ScenarioSignal measures Signal API -> signal processing in an idle Workflow.
	ScenarioSignal = "signal-to-workflow-receive"

	// ModeSequential has one in-flight operation and is the default integrity baseline.
	ModeSequential = "sequential"
	// ModeConcurrent has the configured number of in-flight operations.
	ModeConcurrent = "concurrent"
	// ModeSustained publishes continuously at configured concurrency.
	ModeSustained = "sustained"

	statusPending = "pending"
	statusSuccess = "success"
	statusFailure = "failure"
	statusTimeout = "timeout"
)

// Config controls a benchmark run. Zero-valued queue names are replaced with
// unique names derived from the run ID, which prevents stale local workers or
// tasks from contaminating a run.
type Config struct {
	HostPort  string
	Namespace string

	Iterations       int
	WarmupIterations int
	Concurrency      int
	OperationTimeout time.Duration

	WorkflowTaskQueue string
	ActivityTaskQueue string

	WorkflowPollers              int
	ActivityPollers              int
	WorkflowExecutionConcurrency int
	ActivityExecutionConcurrency int

	Scenario  string
	Mode      string
	OutputDir string

	// These fields are descriptive only. The benchmark never changes server or
	// database durability settings.
	PersistenceBackend string
	DatabaseDurability string
	TaskQueueFairness  string
}

// DefaultConfig returns the documented command-line defaults.
func DefaultConfig() Config {
	return Config{
		HostPort:           "localhost:7233",
		Namespace:          "default",
		Iterations:         10000,
		WarmupIterations:   500,
		Concurrency:        1,
		OperationTimeout:   10 * time.Second,
		WorkflowPollers:    2,
		ActivityPollers:    2,
		Scenario:           "all",
		Mode:               "all",
		OutputDir:          "./results",
		PersistenceBackend: "unknown/not reported",
		DatabaseDurability: "unknown/not reported",
		TaskQueueFairness:  "default/not overridden by benchmark",
	}
}

func (c Config) validate() error {
	if c.HostPort == "" {
		return errors.New("host port must not be empty")
	}
	if c.Namespace == "" {
		return errors.New("namespace must not be empty")
	}
	if c.Iterations <= 0 {
		return fmt.Errorf("iterations must be positive, got %d", c.Iterations)
	}
	if c.WarmupIterations < 0 {
		return fmt.Errorf("warm-up iterations must not be negative, got %d", c.WarmupIterations)
	}
	if c.Concurrency <= 0 {
		return fmt.Errorf("concurrency must be positive, got %d", c.Concurrency)
	}
	if c.OperationTimeout <= 0 {
		return fmt.Errorf("operation timeout must be positive, got %s", c.OperationTimeout)
	}
	if c.WorkflowPollers == 1 {
		return errors.New("Temporal Go SDK does not support exactly one workflow-task poller; use 0/default or at least 2")
	}
	if c.WorkflowPollers < 0 || c.ActivityPollers < 0 {
		return errors.New("poller counts must not be negative")
	}
	if c.WorkflowExecutionConcurrency == 1 {
		return errors.New("Temporal Go SDK does not support exactly one workflow-task execution slot; use 0/default or at least 2")
	}
	if c.WorkflowExecutionConcurrency < 0 || c.ActivityExecutionConcurrency < 0 {
		return errors.New("worker execution concurrency must not be negative")
	}
	if !validScenario(c.Scenario) {
		return fmt.Errorf("unknown scenario %q", c.Scenario)
	}
	if !validMode(c.Mode) {
		return fmt.Errorf("unknown publication mode %q", c.Mode)
	}
	return nil
}

func validScenario(value string) bool {
	return value == "all" || value == ScenarioStart || value == ScenarioActivity || value == ScenarioSignal || value == "start" || value == "activity" || value == "signal"
}

func validMode(value string) bool {
	return value == "all" || value == ModeSequential || value == ModeConcurrent || value == ModeSustained
}

// Operation describes the non-timestamp identity of one publication.
type Operation struct {
	RunID       string `json:"run_id"`
	OperationID string `json:"operation_id"`
	Scenario    string `json:"scenario"`
	Mode        string `json:"mode"`
	Warmup      bool   `json:"warmup"`
}

// Measurement is an immutable snapshot of one operation. Timestamp fields are
// wall-clock representations for export; Duration fields were calculated from
// the shared monotonic clock before the snapshot was made.
type Measurement struct {
	Operation
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`

	T0                time.Time `json:"t0"`
	T1                time.Time `json:"t1"`
	T2                time.Time `json:"t2"`
	T3                time.Time `json:"t3"`
	ScheduleInitiated time.Time `json:"schedule_initiated"`

	PublishToReceive           time.Duration `json:"-"`
	PublishAPILatency          time.Duration `json:"-"`
	PostAcknowledgmentDelivery time.Duration `json:"-"`
	TotalExecution             time.Duration `json:"-"`
	ScheduleToReceive          time.Duration `json:"-"`
}

type measurementEntry struct {
	Measurement
	doneReceive    chan struct{}
	doneComplete   chan struct{}
	receiveClosed  bool
	completeClosed bool
}

// TimestampRegistry is a thread-safe, pre-registration-first timestamp store.
// Register before BeginPublish; this makes a receive event impossible to
// associate with an unrelated operation even when concurrent publications race.
type TimestampRegistry struct {
	mu      sync.Mutex
	clock   func() time.Time
	entries map[string]*measurementEntry
}

// NewTimestampRegistry constructs an empty registry using time.Now.
func NewTimestampRegistry() *TimestampRegistry {
	return &TimestampRegistry{
		clock:   time.Now,
		entries: make(map[string]*measurementEntry),
	}
}

// newTimestampRegistryWithClock is intentionally package-private so tests can
// make ordering and signed post-acknowledgment durations deterministic.
func newTimestampRegistryWithClock(clock func() time.Time) *TimestampRegistry {
	return &TimestampRegistry{clock: clock, entries: make(map[string]*measurementEntry)}
}

// Register pre-registers an operation before its publish API is invoked.
func (r *TimestampRegistry) Register(operation Operation) error {
	if operation.OperationID == "" {
		return errors.New("operation ID must not be empty")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.entries[operation.OperationID]; exists {
		return fmt.Errorf("operation %q already registered", operation.OperationID)
	}
	r.entries[operation.OperationID] = &measurementEntry{
		Measurement:  Measurement{Operation: operation, Status: statusPending},
		doneReceive:  make(chan struct{}),
		doneComplete: make(chan struct{}),
	}
	return nil
}

// BeginPublish records T0. Call it immediately before invoking the Temporal
// publish API, after Register and all non-critical setup are complete.
func (r *TimestampRegistry) BeginPublish(operationID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	entry, err := r.entryLocked(operationID)
	if err != nil {
		return err
	}
	if !entry.T0.IsZero() {
		return fmt.Errorf("operation %q already began publishing", operationID)
	}
	entry.T0 = r.clock()
	return nil
}

// FinishPublish records T1 immediately after the publish API returns. An API
// error remains visible in the raw result and is not converted into a timeout.
func (r *TimestampRegistry) FinishPublish(operationID string, publishErr error) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	entry, err := r.entryLocked(operationID)
	if err != nil {
		return err
	}
	if entry.T0.IsZero() {
		return fmt.Errorf("operation %q has no T0", operationID)
	}
	if !entry.T1.IsZero() {
		return fmt.Errorf("operation %q already has T1", operationID)
	}
	entry.T1 = r.clock()
	if publishErr != nil {
		if errors.Is(publishErr, context.DeadlineExceeded) {
			entry.Status = statusTimeout
		} else {
			entry.Status = statusFailure
		}
		entry.Error = publishErr.Error()
	}
	return nil
}

// MarkReceive records the first correlated receive event only. Later calls are
// ignored, which is the replay-safety invariant for workflow receive events.
func (r *TimestampRegistry) MarkReceive(operationID string) bool {
	return r.markReceive(operationID, false)
}

// MarkWorkflowReceive records a Workflow receive unless the SDK says the
// current execution is replaying. It is used by the worker interceptor and is
// directly testable without starting a Temporal server.
func (r *TimestampRegistry) MarkWorkflowReceive(operationID string, replaying bool) bool {
	if !r.isWorkflowReceive(operationID) {
		return false
	}
	return r.markReceive(operationID, replaying)
}

func (r *TimestampRegistry) markReceive(operationID string, replaying bool) bool {
	if replaying {
		return false
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	entry, ok := r.entries[operationID]
	if !ok || !entry.T2.IsZero() {
		return false
	}
	entry.T2 = r.clock()
	if !entry.receiveClosed {
		close(entry.doneReceive)
		entry.receiveClosed = true
	}
	return true
}

func (r *TimestampRegistry) isWorkflowReceive(operationID string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	entry, ok := r.entries[operationID]
	return ok && (entry.Scenario == ScenarioStart || entry.Scenario == ScenarioSignal)
}

// MarkSchedule records the SDK-level workflow ExecuteActivity initiation. It
// is not the server's durable ActivityTaskScheduled event.
func (r *TimestampRegistry) MarkSchedule(operationID string, replaying bool) bool {
	if replaying {
		return false
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	entry, ok := r.entries[operationID]
	if !ok || !entry.ScheduleInitiated.IsZero() {
		return false
	}
	entry.ScheduleInitiated = r.clock()
	return true
}

// MarkComplete records T3 at the end of the minimal user handler.
func (r *TimestampRegistry) MarkComplete(operationID string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	entry, ok := r.entries[operationID]
	if !ok || !entry.T3.IsZero() {
		return false
	}
	entry.T3 = r.clock()
	if !entry.completeClosed {
		close(entry.doneComplete)
		entry.completeClosed = true
	}
	return true
}

// MarkSuccess marks an operation successful only if it has a receive event.
func (r *TimestampRegistry) MarkSuccess(operationID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	entry, err := r.entryLocked(operationID)
	if err != nil {
		return err
	}
	if entry.Status == statusFailure || entry.Status == statusTimeout {
		return nil
	}
	if entry.T2.IsZero() {
		return fmt.Errorf("operation %q has no receive timestamp", operationID)
	}
	entry.Status = statusSuccess
	return nil
}

// MarkTimeout records a timeout separately from an API or worker failure.
func (r *TimestampRegistry) MarkTimeout(operationID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	entry, err := r.entryLocked(operationID)
	if err != nil {
		return err
	}
	if entry.Status == statusPending {
		entry.Status = statusTimeout
		entry.Error = "receive deadline exceeded"
	}
	return nil
}

// MarkFailure records a publication or worker error when no latency sample can
// be produced. It is used, for example, when the trigger Signal for the
// Activity scenario fails before the Workflow reaches ExecuteActivity.
func (r *TimestampRegistry) MarkFailure(operationID string, failure error) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	entry, err := r.entryLocked(operationID)
	if err != nil {
		return err
	}
	if entry.Status == statusPending {
		entry.Status = statusFailure
		if failure != nil {
			entry.Error = failure.Error()
		} else {
			entry.Error = "operation failed"
		}
	}
	return nil
}

// WaitReceive waits for T2 without writing logs or allocating timing data in
// the critical path.
func (r *TimestampRegistry) WaitReceive(ctx context.Context, operationID string) error {
	r.mu.Lock()
	entry, err := r.entryLocked(operationID)
	if err != nil {
		r.mu.Unlock()
		return err
	}
	done := entry.doneReceive
	alreadyReceived := !entry.T2.IsZero()
	r.mu.Unlock()
	if alreadyReceived {
		return nil
	}
	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// WaitComplete waits for T3 when a handler exposes one.
func (r *TimestampRegistry) WaitComplete(ctx context.Context, operationID string) error {
	r.mu.Lock()
	entry, err := r.entryLocked(operationID)
	if err != nil {
		r.mu.Unlock()
		return err
	}
	done := entry.doneComplete
	alreadyComplete := !entry.T3.IsZero()
	r.mu.Unlock()
	if alreadyComplete {
		return nil
	}
	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Snapshot returns one operation with durations calculated from monotonic
// timestamps. Missing timestamps remain zero and produce zero durations.
func (r *TimestampRegistry) Snapshot(operationID string) (Measurement, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	entry, err := r.entryLocked(operationID)
	if err != nil {
		return Measurement{}, err
	}
	return deriveMeasurement(entry.Measurement), nil
}

// Records returns all registered operations sorted by operation ID. Sorting is
// performed after releasing the registry lock and therefore cannot affect the
// timing critical section.
func (r *TimestampRegistry) Records() []Measurement {
	r.mu.Lock()
	result := make([]Measurement, 0, len(r.entries))
	for _, entry := range r.entries {
		result = append(result, deriveMeasurement(entry.Measurement))
	}
	r.mu.Unlock()
	sort.Slice(result, func(i, j int) bool { return result[i].OperationID < result[j].OperationID })
	return result
}

func (r *TimestampRegistry) entryLocked(operationID string) (*measurementEntry, error) {
	entry, ok := r.entries[operationID]
	if !ok {
		return nil, fmt.Errorf("operation %q is not registered", operationID)
	}
	return entry, nil
}

func deriveMeasurement(measurement Measurement) Measurement {
	if !measurement.T0.IsZero() && !measurement.T2.IsZero() {
		measurement.PublishToReceive = measurement.T2.Sub(measurement.T0)
	}
	if !measurement.T0.IsZero() && !measurement.T1.IsZero() {
		measurement.PublishAPILatency = measurement.T1.Sub(measurement.T0)
	}
	if !measurement.T1.IsZero() && !measurement.T2.IsZero() {
		// This is intentionally signed: receive can race ahead of API return.
		measurement.PostAcknowledgmentDelivery = measurement.T2.Sub(measurement.T1)
	}
	if !measurement.T0.IsZero() && !measurement.T3.IsZero() {
		measurement.TotalExecution = measurement.T3.Sub(measurement.T0)
	}
	if !measurement.ScheduleInitiated.IsZero() && !measurement.T2.IsZero() {
		measurement.ScheduleToReceive = measurement.T2.Sub(measurement.ScheduleInitiated)
	}
	return measurement
}
