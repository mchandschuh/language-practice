package publishlatency

import (
	"context"
	"errors"
	"time"

	enums "go.temporal.io/api/enums/v1"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/converter"
	"go.temporal.io/sdk/interceptor"
	"go.temporal.io/sdk/workflow"
)

const (
	activityTriggerSignal = "publish-latency-activity-trigger"
	signalReceiveSignal   = "publish-latency-signal"
	readyQuery            = "publish-latency-ready"
)

// StartPayload is the minimal correlated input for a started Workflow.
type StartPayload struct {
	OperationID string
}

// ActivityPayload identifies both the publication and the Activity Task Queue.
type ActivityPayload struct {
	OperationID string
	TaskQueue   string
}

// SignalPayload identifies a signal publication.
type SignalPayload struct {
	OperationID string
}

// StartWorkflow is intentionally empty business logic. Its inbound worker
// interceptor captures the first Workflow Task receive before this function is
// invoked, and its return is the optional T3 boundary.
func StartWorkflow(ctx workflow.Context, payload StartPayload) error {
	if payload.OperationID == "" {
		return errors.New("empty operation ID")
	}
	return nil
}

// ActivityCoordinatorWorkflow is an existing, idle Workflow. Each trigger
// signal causes it to schedule exactly one minimal Activity. Workflow.Go is
// used so controlled concurrent publication can create concurrent Activity
// Tasks; this is an explicit concurrency experiment, not the sequential mode.
func ActivityCoordinatorWorkflow(ctx workflow.Context) error {
	if err := workflow.SetQueryHandler(ctx, readyQuery, func() (bool, error) {
		return true, nil
	}); err != nil {
		return err
	}

	triggerChannel := workflow.GetSignalChannel(ctx, activityTriggerSignal)
	for {
		var payload ActivityPayload
		if !triggerChannel.Receive(ctx, &payload) {
			return nil
		}
		if payload.OperationID == "" || payload.TaskQueue == "" {
			continue
		}
		activityPayload := payload
		workflow.Go(ctx, func(activityCtx workflow.Context) {
			activityCtx = workflow.WithActivityOptions(activityCtx, workflow.ActivityOptions{
				TaskQueue:           activityPayload.TaskQueue,
				StartToCloseTimeout: 10 * time.Second,
			})
			var ignored struct{}
			_ = workflow.ExecuteActivity(activityCtx, "MinimalActivity", activityPayload).Get(activityCtx, &ignored)
		})
	}
}

// SignalWorkflow is an idle Workflow that consumes every benchmark signal.
// Signal receive timing is captured in the inbound interceptor before the
// signal is placed on this SDK channel.
func SignalWorkflow(ctx workflow.Context) error {
	if err := workflow.SetQueryHandler(ctx, readyQuery, func() (bool, error) {
		return true, nil
	}); err != nil {
		return err
	}

	signalChannel := workflow.GetSignalChannel(ctx, signalReceiveSignal)
	for {
		var payload SignalPayload
		if !signalChannel.Receive(ctx, &payload) {
			return nil
		}
	}
}

type activities struct {
	registry *TimestampRegistry
}

// MinimalActivity has no artificial delay or business logic. The first
// operation is the correlated receive timestamp; the final operation records
// T3. The SDK has already decoded the Activity input before invoking this
// method, so T2 is an approximation of task receipt rather than raw poll
// response arrival.
func (a *activities) MinimalActivity(_ context.Context, payload ActivityPayload) (struct{}, error) {
	a.registry.MarkReceive(payload.OperationID)
	a.registry.MarkComplete(payload.OperationID)
	return struct{}{}, nil
}

type timingWorkerInterceptor struct {
	interceptor.WorkerInterceptorBase
	registry *TimestampRegistry
}

func (i *timingWorkerInterceptor) InterceptWorkflow(
	ctx workflow.Context,
	next interceptor.WorkflowInboundInterceptor,
) interceptor.WorkflowInboundInterceptor {
	return &timingWorkflowInboundInterceptor{
		WorkflowInboundInterceptorBase: interceptor.WorkflowInboundInterceptorBase{Next: next},
		registry:                       i.registry,
	}
}

func (i *timingWorkerInterceptor) InterceptActivity(
	ctx context.Context,
	next interceptor.ActivityInboundInterceptor,
) interceptor.ActivityInboundInterceptor {
	return &timingActivityInboundInterceptor{
		ActivityInboundInterceptorBase: interceptor.ActivityInboundInterceptorBase{Next: next},
		registry:                       i.registry,
	}
}

type timingWorkflowInboundInterceptor struct {
	interceptor.WorkflowInboundInterceptorBase
	registry *TimestampRegistry
}

func (i *timingWorkflowInboundInterceptor) Init(next interceptor.WorkflowOutboundInterceptor) error {
	return i.Next.Init(&timingWorkflowOutboundInterceptor{
		WorkflowOutboundInterceptorBase: interceptor.WorkflowOutboundInterceptorBase{Next: next},
		registry:                        i.registry,
	})
}

func (i *timingWorkflowInboundInterceptor) ExecuteWorkflow(
	ctx workflow.Context,
	in *interceptor.ExecuteWorkflowInput,
) (any, error) {
	replaying := workflow.IsReplaying(ctx)
	operationID := operationIDFromArgs(in.Args)
	if operationID != "" {
		i.registry.MarkWorkflowReceive(operationID, replaying)
	}
	result, err := i.Next.ExecuteWorkflow(ctx, in)
	if operationID != "" && !replaying {
		i.registry.MarkComplete(operationID)
	}
	return result, err
}

func (i *timingWorkflowInboundInterceptor) HandleSignal(
	ctx workflow.Context,
	in *interceptor.HandleSignalInput,
) error {
	var payload SignalPayload
	if err := converter.GetDefaultDataConverter().FromPayloads(in.Arg, &payload); err == nil && payload.OperationID != "" {
		i.registry.MarkWorkflowReceive(payload.OperationID, workflow.IsReplaying(ctx))
	}
	return i.Next.HandleSignal(ctx, in)
}

type timingWorkflowOutboundInterceptor struct {
	interceptor.WorkflowOutboundInterceptorBase
	registry *TimestampRegistry
}

func (i *timingWorkflowOutboundInterceptor) ExecuteActivity(
	ctx workflow.Context,
	activityType string,
	args ...any,
) workflow.Future {
	if !workflow.IsReplaying(ctx) {
		if payload, ok := activityPayloadFromArgs(args); ok {
			_ = i.registry.MarkSchedule(payload.OperationID, false)
		}
	}
	return i.Next.ExecuteActivity(ctx, activityType, args...)
}

type timingActivityInboundInterceptor struct {
	interceptor.ActivityInboundInterceptorBase
	registry *TimestampRegistry
}

func (i *timingActivityInboundInterceptor) ExecuteActivity(
	ctx context.Context,
	in *interceptor.ExecuteActivityInput,
) (any, error) {
	// The actual user Activity method also records T2. This interceptor is kept
	// as a no-op wrapper so the worker can use one interceptor chain for all
	// task types without moving the Activity boundary earlier than documented.
	return i.Next.ExecuteActivity(ctx, in)
}

func operationIDFromArgs(args []any) string {
	for _, arg := range args {
		switch payload := arg.(type) {
		case StartPayload:
			return payload.OperationID
		case *StartPayload:
			if payload != nil {
				return payload.OperationID
			}
		}
	}
	return ""
}

func activityPayloadFromArgs(args []any) (ActivityPayload, bool) {
	for _, arg := range args {
		switch payload := arg.(type) {
		case ActivityPayload:
			return payload, true
		case *ActivityPayload:
			if payload != nil {
				return *payload, true
			}
		}
	}
	return ActivityPayload{}, false
}

func workflowStartOptions(id, taskQueue string) client.StartWorkflowOptions {
	return client.StartWorkflowOptions{
		ID:                       id,
		TaskQueue:                taskQueue,
		WorkflowIDReusePolicy:    enums.WORKFLOW_ID_REUSE_POLICY_REJECT_DUPLICATE,
		WorkflowIDConflictPolicy: enums.WORKFLOW_ID_CONFLICT_POLICY_FAIL,
	}
}
