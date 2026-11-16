package hello

import (
	"context"
	"fmt"
	"time"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/workflow"
)

// TaskQueue is the queue polled by the hello/worker
const TaskQueue = "hello"

// HelloWorkflow schedules HelloActivity and returns its greeting
func HelloWorkflow(ctx workflow.Context, name string) (string, error) {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	})

	var result string
	future := workflow.ExecuteActivity(ctx, HelloActivity, name)
	err := future.Get(ctx, &result)
	return result, err
}

// HelloActivity returns a greeting for a provided name
func HelloActivity(ctx context.Context, name string) (string, error) {
	log := activity.GetLogger(ctx)
	result := fmt.Sprintf("Hello %s!", name)
	log.Info(result)
	return result, nil
}
