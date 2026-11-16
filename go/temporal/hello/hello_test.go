package hello

import (
	"testing"

	"go.temporal.io/sdk/testsuite"
)

func TestHelloActivity(t *testing.T) {
	var testSuite testsuite.WorkflowTestSuite
	env := testSuite.NewTestActivityEnvironment()
	env.RegisterActivity(HelloActivity)

	value, err := env.ExecuteActivity(HelloActivity, "World")
	if err != nil {
		t.Fatalf("executing activity: %v", err)
	}

	var got string
	if err := value.Get(&got); err != nil {
		t.Fatalf("getting activity result: %v", err)
	}

	if want := "Hello World!"; got != want {
		t.Fatalf("activity result = %q, want %q", got, want)
	}
}

func TestHelloWorkflow(t *testing.T) {
	var testSuite testsuite.WorkflowTestSuite
	env := testSuite.NewTestWorkflowEnvironment()

	env.RegisterActivity(HelloActivity)
	env.ExecuteWorkflow(HelloWorkflow, "World")

	if !env.IsWorkflowCompleted() {
		t.Fatal("workflow did not complete")
	}

	if err := env.GetWorkflowError(); err != nil {
		t.Fatalf("workflow failed: %v", err)
	}

	var got string
	if err := env.GetWorkflowResult(&got); err != nil {
		t.Fatalf("getting workflow result: %v", err)
	}

	if want := "Hello World!"; got != want {
		t.Fatalf("workflow result = %q, want %q", got, want)
	}
}
