package main

import (
	"testing"

	enums "go.temporal.io/api/enums/v1"

	"github.com/mchandschuh/language-practice/go/temporal/hello"
)

func TestStartWorkflowOptionsRejectDuplicateWorkflowID(t *testing.T) {
	options := startWorkflowOptions()

	if got, want := options.ID, "hello-workflow"; got != want {
		t.Fatalf("workflow ID = %q, want %q", got, want)
	}
	if got, want := options.TaskQueue, hello.TaskQueue; got != want {
		t.Fatalf("task queue = %q, want %q", got, want)
	}
	if got, want := options.WorkflowIDConflictPolicy, enums.WORKFLOW_ID_CONFLICT_POLICY_FAIL; got != want {
		t.Fatalf("workflow ID conflict policy = %v, want %v", got, want)
	}
	if got, want := options.WorkflowIDReusePolicy, enums.WORKFLOW_ID_REUSE_POLICY_REJECT_DUPLICATE; got != want {
		t.Fatalf("workflow ID reuse policy = %v, want %v", got, want)
	}
	if !options.WorkflowExecutionErrorWhenAlreadyStarted {
		t.Fatal("workflow execution should report duplicate-start errors")
	}
}
