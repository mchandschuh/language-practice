package publishlatency

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

func TestTimestampRegistryCorrelatesConcurrentOperations(t *testing.T) {
	base := time.Now()
	var clockMu sync.Mutex
	clockValue := base
	registry := newTimestampRegistryWithClock(func() time.Time {
		clockMu.Lock()
		defer clockMu.Unlock()
		value := clockValue
		clockValue = clockValue.Add(time.Microsecond)
		return value
	})

	const operationCount = 64
	var waitGroup sync.WaitGroup
	for i := 0; i < operationCount; i++ {
		waitGroup.Add(1)
		go func(i int) {
			defer waitGroup.Done()
			operationID := "concurrent/" + time.Duration(i).String()
			if err := registry.Register(Operation{RunID: "run", OperationID: operationID, Scenario: ScenarioStart, Mode: ModeConcurrent}); err != nil {
				t.Errorf("register %s: %v", operationID, err)
				return
			}
			if err := registry.BeginPublish(operationID); err != nil {
				t.Errorf("begin %s: %v", operationID, err)
				return
			}
			if !registry.MarkWorkflowReceive(operationID, false) {
				t.Errorf("mark receive %s returned false", operationID)
			}
			if err := registry.FinishPublish(operationID, nil); err != nil {
				t.Errorf("finish %s: %v", operationID, err)
				return
			}
			if err := registry.MarkSuccess(operationID); err != nil {
				t.Errorf("mark success %s: %v", operationID, err)
			}
		}(i)
	}
	waitGroup.Wait()

	records := registry.Records()
	if got, want := len(records), operationCount; got != want {
		t.Fatalf("record count = %d, want %d", got, want)
	}
	for _, record := range records {
		if record.Status != statusSuccess {
			t.Errorf("operation %s status = %q, want success", record.OperationID, record.Status)
		}
		if record.T0.IsZero() || record.T1.IsZero() || record.T2.IsZero() {
			t.Errorf("operation %s missing timestamp: %+v", record.OperationID, record)
		}
	}
}

func TestTimestampRegistryPreservesDeliveryBeforePublishAcknowledgment(t *testing.T) {
	base := time.Now()
	steps := []time.Time{base, base.Add(10 * time.Microsecond), base.Add(20 * time.Microsecond)}
	index := 0
	registry := newTimestampRegistryWithClock(func() time.Time {
		value := steps[index]
		index++
		return value
	})
	operationID := "delivery-before-ack"
	if err := registry.Register(Operation{RunID: "run", OperationID: operationID, Scenario: ScenarioSignal, Mode: ModeSequential}); err != nil {
		t.Fatal(err)
	}
	if err := registry.BeginPublish(operationID); err != nil {
		t.Fatal(err)
	}
	if !registry.MarkWorkflowReceive(operationID, false) {
		t.Fatal("receive was not recorded")
	}
	if err := registry.FinishPublish(operationID, nil); err != nil {
		t.Fatal(err)
	}
	if err := registry.MarkSuccess(operationID); err != nil {
		t.Fatal(err)
	}
	record, err := registry.Snapshot(operationID)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := record.PostAcknowledgmentDelivery, -10*time.Microsecond; got != want {
		t.Fatalf("post-ack delivery = %s, want %s", got, want)
	}
}

func TestTimestampRegistryIgnoresReplayAndDuplicateReceive(t *testing.T) {
	registry := NewTimestampRegistry()
	operationID := "replay/operation"
	if err := registry.Register(Operation{RunID: "run", OperationID: operationID, Scenario: ScenarioSignal, Mode: ModeSequential}); err != nil {
		t.Fatal(err)
	}
	if err := registry.BeginPublish(operationID); err != nil {
		t.Fatal(err)
	}
	if registry.MarkWorkflowReceive(operationID, true) {
		t.Fatal("replay created a receive event")
	}
	if !registry.MarkWorkflowReceive(operationID, false) {
		t.Fatal("first non-replay receive was not recorded")
	}
	if registry.MarkWorkflowReceive(operationID, false) {
		t.Fatal("duplicate receive was counted")
	}
}

func TestTimestampRegistryTimeoutAndErrorStayVisible(t *testing.T) {
	registry := NewTimestampRegistry()
	for _, test := range []struct {
		name       string
		publishErr error
		wantStatus string
	}{
		{name: "publish failure", publishErr: errors.New("service unavailable"), wantStatus: statusFailure},
		{name: "publish timeout", publishErr: context.DeadlineExceeded, wantStatus: statusTimeout},
		{name: "timeout", wantStatus: statusTimeout},
	} {
		t.Run(test.name, func(t *testing.T) {
			operationID := "failure/" + test.name
			if err := registry.Register(Operation{RunID: "run", OperationID: operationID, Scenario: ScenarioStart, Mode: ModeSequential}); err != nil {
				t.Fatal(err)
			}
			if err := registry.BeginPublish(operationID); err != nil {
				t.Fatal(err)
			}
			if err := registry.FinishPublish(operationID, test.publishErr); err != nil {
				t.Fatal(err)
			}
			if test.publishErr == nil {
				if err := registry.MarkTimeout(operationID); err != nil {
					t.Fatal(err)
				}
			}
			record, err := registry.Snapshot(operationID)
			if err != nil {
				t.Fatal(err)
			}
			if record.Status != test.wantStatus {
				t.Fatalf("status = %q, want %q", record.Status, test.wantStatus)
			}
			if record.Error == "" {
				t.Fatal("failure/timeout error was not recorded")
			}
		})
	}
}

func TestWaitReceiveHonorsContext(t *testing.T) {
	registry := NewTimestampRegistry()
	operationID := "wait/timeout"
	if err := registry.Register(Operation{RunID: "run", OperationID: operationID, Scenario: ScenarioStart, Mode: ModeSequential}); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := registry.WaitReceive(ctx, operationID); !errors.Is(err, context.Canceled) {
		t.Fatalf("wait error = %v, want context.Canceled", err)
	}
}
