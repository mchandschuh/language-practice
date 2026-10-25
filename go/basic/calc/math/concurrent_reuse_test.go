package calc

import (
	"context"
	"errors"
	"runtime"
	"sync"
	"testing"
	"time"
)

const concurrentReuseTimeout = time.Second

func TestEvaluatorRuntimeConcurrentlyReusesImmutableAST(t *testing.T) {
	expression, err := Parse("x * x + 2")
	if err != nil {
		t.Fatal(err)
	}

	const workload = 16
	evaluator, err := NewEvaluator(context.Background(), workload)
	if err != nil {
		t.Fatal(err)
	}
	defer evaluator.Close()

	start := make(chan struct{})
	results := make([]<-chan EvaluationResult, workload)
	inputs := make([]Environment, workload)
	errs := make(chan error, workload)
	var workers sync.WaitGroup
	for i := 0; i < workload; i++ {
		i := i
		inputs[i] = Environment{"x": float64(i)}
		workers.Add(1)
		go func() {
			defer workers.Done()
			<-start
			result, submitErr := evaluator.Submit(context.Background(), expression, inputs[i])
			if submitErr != nil {
				errs <- submitErr
				return
			}
			results[i] = result
		}()
	}
	close(start)
	workers.Wait()
	close(errs)
	for err := range errs {
		t.Fatal(err)
	}

	for i, result := range results {
		got := receiveResultWithin(t, result, concurrentReuseTimeout)
		want := float64(i*i + 2)
		if got.Err != nil || got.Value != want {
			t.Fatalf("request %d result = %+v, want value %v", i, got, want)
		}
		if got.Environment["x"] != float64(i) {
			t.Fatalf("request %d environment = %v, want x=%d", i, got.Environment, i)
		}
		if inputs[i]["x"] != float64(i) {
			t.Fatalf("request %d mutated caller input: %v", i, inputs[i])
		}
	}
}

func TestEvaluatorRuntimeConcurrentMutationKeepsRequestStateIsolated(t *testing.T) {
	expression, err := Parse("x++ + ++x")
	if err != nil {
		t.Fatal(err)
	}

	const workload = 12
	evaluator, err := NewEvaluator(context.Background(), workload)
	if err != nil {
		t.Fatal(err)
	}
	defer evaluator.Close()

	results := make([]<-chan EvaluationResult, workload)
	inputs := make([]Environment, workload)
	for i := range results {
		value := float64(i + 1)
		inputs[i] = Environment{"x": value}
		results[i], err = evaluator.Submit(context.Background(), expression, inputs[i])
		if err != nil {
			t.Fatal(err)
		}
	}

	for i, result := range results {
		got := receiveResultWithin(t, result, concurrentReuseTimeout)
		initial := float64(i + 1)
		if got.Err != nil || got.Value != 2*initial+2 || got.Environment["x"] != initial+2 {
			t.Fatalf("request %d result = %+v, want value=%v x=%v", i, got, 2*initial+2, initial+2)
		}
		if inputs[i]["x"] != initial {
			t.Fatalf("request %d mutated caller input: %v", i, inputs[i])
		}
	}

}

func TestEvaluatorRuntimeMatchesSequentialOracleCorpus(t *testing.T) {
	tests := []struct {
		name  string
		expr  string
		input Environment
	}{
		{name: "precedence", expr: "(9 - 4) / 5 + 2 * 3"},
		{name: "unary", expr: "-x + 3", input: Environment{"x": 4}},
		{name: "prefix mutation", expr: "++x", input: Environment{"x": 5}},
		{name: "postfix mutation", expr: "x++", input: Environment{"x": 5}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expression, err := Parse(test.expr)
			if err != nil {
				t.Fatal(err)
			}
			wantValue, wantEnvironment, err := (Evaluator{}).EvaluateWithEnvironment(expression, test.input)
			if err != nil {
				t.Fatal(err)
			}

			evaluator, err := NewEvaluator(context.Background(), 1)
			if err != nil {
				t.Fatal(err)
			}
			defer evaluator.Close()
			result, err := evaluator.Submit(context.Background(), expression, test.input)
			if err != nil {
				t.Fatal(err)
			}
			got := receiveResultWithin(t, result, concurrentReuseTimeout)
			if got.Err != nil || got.Value != wantValue {
				t.Fatalf("channel result = %+v, sequential value = %v", got, wantValue)
			}
			if got.Environment["x"] != wantEnvironment["x"] {
				t.Fatalf("channel environment = %v, sequential environment = %v", got.Environment, wantEnvironment)
			}
		})
	}
}

func TestEvaluatorRuntimeConcurrentAdmissionAndShutdown(t *testing.T) {
	evaluator, err := NewEvaluator(context.Background(), 4)
	if err != nil {
		t.Fatal(err)
	}

	started := make(chan struct{})
	release := make(chan struct{})
	var barrier sync.Once
	evaluator.beforeRun = func() {
		barrier.Do(func() {
			close(started)
			<-release
		})
	}

	first, err := evaluator.Submit(context.Background(), &NumberExpression{Value: 1}, nil)
	if err != nil {
		t.Fatal(err)
	}
	<-started

	const queued = 3
	results := make([]<-chan EvaluationResult, queued)
	submitErrors := make(chan error, queued)
	admitted := make(chan struct{}, queued)
	var submit sync.WaitGroup
	for i := range results {
		i := i
		submit.Add(1)
		go func() {
			defer submit.Done()
			result, submitErr := evaluator.Submit(context.Background(), &NumberExpression{Value: float64(i + 2)}, nil)
			results[i] = result
			if submitErr != nil {
				submitErrors <- submitErr
				return
			}
			admitted <- struct{}{}
		}()
	}
	for i := 0; i < queued; i++ {
		select {
		case <-admitted:
		case <-time.After(concurrentReuseTimeout):
			t.Fatal("timed out admitting queued request")
		}
	}
	closeDone := make(chan error, 1)
	go func() { closeDone <- evaluator.Close() }()
	select {
	case err := <-closeDone:
		t.Fatalf("Close returned before active work was released: %v", err)
	case <-time.After(20 * time.Millisecond):
	}
	close(release)
	submit.Wait()
	close(submitErrors)
	for err := range submitErrors {
		t.Fatal(err)
	}
	if err := <-closeDone; err != nil {
		t.Fatal(err)
	}
	if got := receiveResultWithin(t, first, concurrentReuseTimeout); !errors.Is(got.Err, ErrEvaluatorClosed) {
		t.Fatalf("active result = %+v", got)
	}
	for i, result := range results {
		if result == nil {
			t.Fatalf("queued request %d was not admitted", i)
		}
		if got := receiveResultWithin(t, result, concurrentReuseTimeout); !errors.Is(got.Err, ErrEvaluatorClosed) {
			t.Fatalf("queued result %d = %+v", i, got)
		}
	}
}

func TestEvaluatorRuntimeRepeatedCloseDoesNotAccumulateWorkers(t *testing.T) {
	baseline := runtime.NumGoroutine()
	const workload = 64
	for i := 0; i < workload; i++ {
		evaluator, err := NewEvaluator(context.Background(), 1)
		if err != nil {
			t.Fatal(err)
		}
		if err := evaluator.Close(); err != nil {
			t.Fatal(err)
		}
	}

	deadline := time.Now().Add(concurrentReuseTimeout)
	for runtime.NumGoroutine() > baseline+2 && time.Now().Before(deadline) {
		runtime.Gosched()
	}
	if got := runtime.NumGoroutine(); got > baseline+2 {
		t.Fatalf("goroutine count grew beyond tolerance: baseline=%d after=%d tolerance=2", baseline, got)
	}
}

func receiveResultWithin(t *testing.T, result <-chan EvaluationResult, timeout time.Duration) EvaluationResult {
	t.Helper()
	select {
	case value := <-result:
		return value
	case <-time.After(timeout):
		t.Fatalf("timed out waiting for evaluation result after %s", timeout)
		return EvaluationResult{}
	}
}
