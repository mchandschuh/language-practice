package publishlatency

import "testing"

func TestConfigRejectsUnsupportedTemporalWorkflowValues(t *testing.T) {
	config := DefaultConfig()
	config.WorkflowPollers = 1
	if err := config.validate(); err == nil {
		t.Fatal("workflow poller count 1 was accepted")
	}
	config = DefaultConfig()
	config.WorkflowExecutionConcurrency = 1
	if err := config.validate(); err == nil {
		t.Fatal("workflow execution concurrency 1 was accepted")
	}
}

func TestSelectedScenariosAndModes(t *testing.T) {
	if got := selectedScenarios("activity"); len(got) != 1 || got[0] != ScenarioActivity {
		t.Fatalf("activity selection = %v", got)
	}
	if got := selectedModes(ModeSustained); len(got) != 1 || got[0] != ModeSustained {
		t.Fatalf("sustained selection = %v", got)
	}
	if got := selectedScenarios("all"); len(got) != 3 {
		t.Fatalf("all scenario selection = %v", got)
	}
}
