package hello

import (
	"context"
	"strings"
	"testing"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/converter"
	"go.temporal.io/sdk/testsuite"
)

func TestRetryWorkflowRetriesTransientActivityFailure(t *testing.T) {
	var testSuite testsuite.WorkflowTestSuite
	env := testSuite.NewTestWorkflowEnvironment()
	var attempts []int32
	env.SetOnActivityStartedListener(func(info *activity.Info, _ context.Context, _ converter.EncodedValues) {
		attempts = append(attempts, info.Attempt)
	})
	env.RegisterActivity(RetryActivity)

	env.ExecuteWorkflow(RetryWorkflow, RetryRequest{FailUntilAttempt: 2})

	if err := env.GetWorkflowError(); err != nil {
		t.Fatalf("workflow failed: %v", err)
	}

	var got string
	if err := env.GetWorkflowResult(&got); err != nil {
		t.Fatalf("getting workflow result: %v", err)
	}
	if want := "succeeded on attempt 3"; got != want {
		t.Fatalf("workflow result = %q, want %q", got, want)
	}
	if want := []int32{1, 2, 3}; !equalInt32Slices(attempts, want) {
		t.Fatalf("activity attempts = %v, want %v", attempts, want)
	}
}

func TestRetryWorkflowStopsAfterMaximumAttempts(t *testing.T) {
	var testSuite testsuite.WorkflowTestSuite
	env := testSuite.NewTestWorkflowEnvironment()
	var attempts []int32
	env.SetOnActivityStartedListener(func(info *activity.Info, _ context.Context, _ converter.EncodedValues) {
		attempts = append(attempts, info.Attempt)
	})
	env.RegisterActivity(RetryActivity)

	env.ExecuteWorkflow(RetryWorkflow, RetryRequest{FailUntilAttempt: 3})

	if err := env.GetWorkflowError(); err == nil {
		t.Fatal("workflow succeeded, want retry exhaustion error")
	} else if !strings.Contains(err.Error(), "failure on attempt 3") {
		t.Fatalf("workflow error = %q, want final attempt in error", err)
	}
	if want := []int32{1, 2, 3}; !equalInt32Slices(attempts, want) {
		t.Fatalf("activity attempts = %v, want %v", attempts, want)
	}
}

func TestRetryWorkflowDoesNotRetryPermanentActivityFailure(t *testing.T) {
	var testSuite testsuite.WorkflowTestSuite
	env := testSuite.NewTestWorkflowEnvironment()
	var attempts []int32
	env.SetOnActivityStartedListener(func(info *activity.Info, _ context.Context, _ converter.EncodedValues) {
		attempts = append(attempts, info.Attempt)
	})
	env.RegisterActivity(RetryActivity)

	env.ExecuteWorkflow(RetryWorkflow, RetryRequest{
		FailUntilAttempt: 1,
		PermanentFailure: true,
	})

	if err := env.GetWorkflowError(); err == nil {
		t.Fatal("workflow succeeded, want permanent failure")
	} else if !strings.Contains(err.Error(), "failure on attempt 1") {
		t.Fatalf("workflow error = %q, want permanent failure in error", err)
	}
	if want := []int32{1}; !equalInt32Slices(attempts, want) {
		t.Fatalf("activity attempts = %v, want %v", attempts, want)
	}
}

func equalInt32Slices(got, want []int32) bool {
	if len(got) != len(want) {
		return false
	}
	for i := range got {
		if got[i] != want[i] {
			return false
		}
	}
	return true
}
