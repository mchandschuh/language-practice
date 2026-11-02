package main

import (
	"context"
	"log"

	enums "go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/client"

	"github.com/mchandschuh/language-practice/go/temporal/hello"
)

func startWorkflowOptions() client.StartWorkflowOptions {
	return client.StartWorkflowOptions{
		ID:                                       "hello-workflow",
		TaskQueue:                                hello.TaskQueue,
		WorkflowIDConflictPolicy:                 enums.WORKFLOW_ID_CONFLICT_POLICY_FAIL,
		WorkflowIDReusePolicy:                    enums.WORKFLOW_ID_REUSE_POLICY_REJECT_DUPLICATE,
		WorkflowExecutionErrorWhenAlreadyStarted: true,
	}
}

func main() {
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}

	defer c.Close()

	options := startWorkflowOptions()

	workflowRun, err := c.ExecuteWorkflow(
		context.Background(),
		options,
		hello.HelloWorkflow,
		"World",
	)

	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}

	log.Println("Started workflow",
		"WorkflowID", workflowRun.GetID(),
		"RunID", workflowRun.GetRunID(),
	)

	var result string
	if err := workflowRun.Get(context.Background(), &result); err != nil {
		log.Fatalln("Unable to get workflow result", err)
	}

	log.Println("Workflow result:", result)

}
