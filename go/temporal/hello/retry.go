package hello

import (
	"context"
	"fmt"
	"time"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

const retryMaximumAttempts int32 = 3

// RetryRequest controls the deliberately failing attempts in RetryActivity.
//
// The request is serializable because it crosses the Workflow/Activity
// boundary. PermanentFailure demonstrates an error that should not be retried.
type RetryRequest struct {
	FailUntilAttempt int32
	PermanentFailure bool
}

// RetryWorkflow runs an Activity with a bounded retry policy.
func RetryWorkflow(ctx workflow.Context, request RetryRequest) (string, error) {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Millisecond,
			BackoffCoefficient: 1,
			MaximumInterval:    time.Millisecond,
			MaximumAttempts:    retryMaximumAttempts,
		},
	})

	var result string
	err := workflow.ExecuteActivity(ctx, RetryActivity, request).Get(ctx, &result)
	return result, err
}

// RetryActivity simulates transient and permanent failures.
//
// This Activity intentionally has no external side effect. A real Activity
// must make an external effect idempotent because every attempt can run the
// effect again after a timeout or Worker failure.
func RetryActivity(ctx context.Context, request RetryRequest) (string, error) {
	attempt := activity.GetInfo(ctx).Attempt
	if attempt <= request.FailUntilAttempt {
		message := fmt.Sprintf("failure on attempt %d", attempt)
		if request.PermanentFailure {
			return "", temporal.NewNonRetryableApplicationError(message, "PermanentFailure", nil)
		}
		return "", temporal.NewApplicationError(message, "TransientFailure")
	}

	return fmt.Sprintf("succeeded on attempt %d", attempt), nil
}
