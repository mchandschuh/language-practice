package calc

import (
	"context"
	"errors"
	"runtime"
	"testing"
	"time"
)

func receiveResult(t *testing.T, result <-chan EvaluationResult) EvaluationResult {
	t.Helper()
	select {
	case value := <-result:
		return value
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for evaluation result")
		return EvaluationResult{}
	}
}

func TestEvaluatorRuntimeSubmitCompletesAndReusesAST(t *testing.T) {
	evaluator, err := NewEvaluator(context.Background(), 2)
	if err != nil {
		t.Fatal(err)
	}
	defer evaluator.Close()
	expression, err := Parse("x * x + 2")
	if err != nil {
		t.Fatal(err)
	}
	for _, test := range []struct {
		input float64
		want  float64
	}{{3, 11}, {10, 102}} {
		result, err := evaluator.Submit(context.Background(), expression, Environment{"x": test.input})
		if err != nil {
			t.Fatal(err)
		}
		got := receiveResult(t, result)
		if got.Err != nil || got.Value != test.want {
			t.Fatalf("result = %+v, want value %v", got, test.want)
		}
	}
}

func TestEvaluatorRuntimeMultipleRequests(t *testing.T) {
	evaluator, err := NewEvaluator(context.Background(), 3)
	if err != nil {
		t.Fatal(err)
	}
	defer evaluator.Close()
	results := make([]<-chan EvaluationResult, 8)
	for i := range results {
		results[i], err = evaluator.Submit(context.Background(), &NumberExpression{Value: float64(i)}, Environment{})
		if err != nil {
			t.Fatal(err)
		}
	}
	for i, result := range results {
		got := receiveResult(t, result)
		if got.Err != nil || got.Value != float64(i) {
			t.Fatalf("request %d result = %+v", i, got)
		}
	}
}

func TestEvaluatorRuntimeAdmissionCancellationAndCapacity(t *testing.T) {
	evaluator, err := NewEvaluator(context.Background(), 1)
	if err != nil {
		t.Fatal(err)
	}
	defer evaluator.Close()
	first, err := evaluator.Submit(context.Background(), &NumberExpression{Value: 1}, Environment{})
	if err != nil {
		t.Fatal(err)
	}
	second, err := evaluator.Submit(context.Background(), &NumberExpression{Value: 2}, Environment{})
	if err != nil {
		t.Fatal(err)
	}
	if receiveResult(t, first).Err != nil {
		t.Fatal("first request failed")
	}
	if receiveResult(t, second).Err != nil {
		t.Fatal("second request failed")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := evaluator.Submit(ctx, &NumberExpression{Value: 3}, Environment{}); !errors.Is(err, context.Canceled) {
		t.Fatalf("cancelled admission error = %v", err)
	}
}

func TestEvaluatorRuntimeCloseRejectsNewRequestsAndIsIdempotent(t *testing.T) {
	before := runtime.NumGoroutine()
	evaluator, err := NewEvaluator(context.Background(), 1)
	if err != nil {
		t.Fatal(err)
	}
	if err := evaluator.Close(); err != nil {
		t.Fatal(err)
	}
	if err := evaluator.Close(); err != nil {
		t.Fatal(err)
	}
	if _, err := evaluator.Submit(context.Background(), &NumberExpression{Value: 1}, Environment{}); !errors.Is(err, ErrEvaluatorClosed) {
		t.Fatalf("submit after close error = %v", err)
	}
	deadline := time.Now().Add(time.Second)
	for runtime.NumGoroutine() > before && time.Now().Before(deadline) {
		runtime.Gosched()
	}
	if runtime.NumGoroutine() > before {
		t.Fatalf("worker goroutine did not stop: before=%d after=%d", before, runtime.NumGoroutine())
	}
}

func TestEvaluatorRuntimeCloseCancelsAdmittedRequest(t *testing.T) {
	evaluator, err := NewEvaluator(context.Background(), 1)
	if err != nil {
		t.Fatal(err)
	}
	result, err := evaluator.Submit(context.Background(), &NumberExpression{Value: 7}, Environment{})
	if err != nil {
		t.Fatal(err)
	}
	if err := evaluator.Close(); err != nil {
		t.Fatal(err)
	}
	if got := receiveResult(t, result); !errors.Is(got.Err, ErrEvaluatorClosed) {
		t.Fatalf("close result = %+v", got)
	}
}
