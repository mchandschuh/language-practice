package calc

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"
)

type unsupportedExpression struct{}

func (*unsupportedExpression) expressionNode() {}

func TestEvaluatorRuntimeDispatchesAllExpressionTypes(t *testing.T) {
	evaluator, err := NewEvaluator(context.Background(), 2)
	if err != nil {
		t.Fatal(err)
	}
	defer evaluator.Close()

	tests := []struct {
		name       string
		expression MathExpression
		input      Environment
		want       float64
	}{
		{"number", &NumberExpression{Value: 7}, nil, 7},
		{"identifier", &IdentifierExpression{Name: "x"}, Environment{"x": 8}, 8},
		{"unary", &UnaryExpression{Operator: Negate, Operand: &NumberExpression{Value: 3}}, nil, -3},
		{"binary", &BinaryExpression{Operator: Add, Left: &NumberExpression{Value: 2}, Right: &NumberExpression{Value: 5}}, nil, 7},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := evaluator.Submit(context.Background(), test.expression, test.input)
			if err != nil {
				t.Fatal(err)
			}
			got := receiveResult(t, result)
			if got.Err != nil || got.Value != test.want {
				t.Fatalf("result = %+v, want %v", got, test.want)
			}
		})
	}
}

func TestEvaluatorRuntimeDispatchesNestedExpressionsWithoutSingleWorkerDeadlock(t *testing.T) {
	evaluator, err := NewEvaluator(context.Background(), 1)
	if err != nil {
		t.Fatal(err)
	}
	defer evaluator.Close()

	expression := &UnaryExpression{
		Operator: Negate,
		Operand: &BinaryExpression{
			Operator: Add,
			Left:     &NumberExpression{Value: 2},
			Right:    &UnaryExpression{Operator: Negate, Operand: &NumberExpression{Value: 5}},
		},
	}
	result, err := evaluator.Submit(context.Background(), expression, nil)
	if err != nil {
		t.Fatal(err)
	}
	got := receiveResult(t, result)
	if got.Err != nil || got.Value != 3 {
		t.Fatalf("nested result = %+v, want 3", got)
	}
}

func TestEvaluatorRuntimeDispatchesDeepBoundedExpression(t *testing.T) {
	var expression MathExpression = &NumberExpression{Value: 1}
	for i := 0; i < maxExpressionDepth; i++ {
		expression = &UnaryExpression{Operator: Positive, Operand: expression}
	}
	evaluator, err := NewEvaluator(context.Background(), 1)
	if err != nil {
		t.Fatal(err)
	}
	defer evaluator.Close()
	result, err := evaluator.Submit(context.Background(), expression, nil)
	if err != nil {
		t.Fatal(err)
	}
	if got := receiveResult(t, result); got.Err != nil || got.Value != 1 {
		t.Fatalf("bounded result = %+v", got)
	}
}

func TestEvaluatorRuntimeRejectsExpressionBeyondDepthLimit(t *testing.T) {
	var expression MathExpression = &NumberExpression{Value: 1}
	for i := 0; i <= maxExpressionDepth; i++ {
		expression = &UnaryExpression{Operator: Positive, Operand: expression}
	}
	evaluator, err := NewEvaluator(context.Background(), 1)
	if err != nil {
		t.Fatal(err)
	}
	defer evaluator.Close()
	result, err := evaluator.Submit(context.Background(), expression, nil)
	if err != nil {
		t.Fatal(err)
	}
	got := receiveResult(t, result)
	if !errors.Is(got.Err, errExpressionDepth) {
		t.Fatalf("depth error = %v, want %v", got.Err, errExpressionDepth)
	}
}

func TestEvaluatorRuntimePropagatesChildErrorsAndUnsupportedNodes(t *testing.T) {
	evaluator, err := NewEvaluator(context.Background(), 2)
	if err != nil {
		t.Fatal(err)
	}
	defer evaluator.Close()

	missing := &BinaryExpression{Operator: Add, Left: &IdentifierExpression{Name: "missing"}, Right: &NumberExpression{Value: 1}}
	result, err := evaluator.Submit(context.Background(), missing, nil)
	if err != nil {
		t.Fatal(err)
	}
	got := receiveResult(t, result)
	if got.Err == nil || !strings.Contains(got.Err.Error(), `identifier "missing" is not defined`) {
		t.Fatalf("child error = %v", got.Err)
	}

	result, err = evaluator.Submit(context.Background(), &unsupportedExpression{}, nil)
	if err != nil {
		t.Fatal(err)
	}
	got = receiveResult(t, result)
	if got.Err == nil || !strings.Contains(got.Err.Error(), "unsupported expression node") {
		t.Fatalf("unsupported error = %v", got.Err)
	}
}

func TestEvaluatorRuntimeDispatchPreservesRequestStateAndIndependentResults(t *testing.T) {
	evaluator, err := NewEvaluator(context.Background(), 4)
	if err != nil {
		t.Fatal(err)
	}
	defer evaluator.Close()

	expression, err := Parse("x++ + ++x")
	if err != nil {
		t.Fatal(err)
	}
	results := make([]<-chan EvaluationResult, 2)
	for i, value := range []float64{5, 10} {
		results[i], err = evaluator.Submit(context.Background(), expression, Environment{"x": value})
		if err != nil {
			t.Fatal(err)
		}
	}
	for i, result := range results {
		got := receiveResult(t, result)
		want := []float64{12, 22}[i]
		if got.Err != nil || got.Value != want || got.Environment["x"] != []float64{7, 12}[i] {
			t.Fatalf("request %d result = %+v", i, got)
		}
	}
}

func TestDispatchCancellationDoesNotBlock(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	result := NewResultChannel()
	dispatchExpression(ctx, &NumberExpression{Value: 1}, Environment{}, result, 0)
	select {
	case got := <-result:
		t.Fatalf("cancelled dispatch delivered result: %+v", got)
	case <-time.After(20 * time.Millisecond):
	}
}

func TestUnaryWorkersSelectEveryOperator(t *testing.T) {
	tests := []struct {
		operator UnaryOperator
		wantName string
		want     float64
	}{
		{Positive, "positive", 5},
		{Negate, "negate", -5},
		{PrefixDecrement, "prefix decrement", 4},
		{PostfixDecrement, "postfix decrement", 5},
		{PrefixIncrement, "prefix increment", 6},
		{PostfixIncrement, "postfix increment", 5},
	}
	for _, test := range tests {
		t.Run(test.wantName, func(t *testing.T) {
			name, ok := unaryWorkerName(test.operator)
			if !ok || name != test.wantName {
				t.Fatalf("worker = %q, %v; want %q", name, ok, test.wantName)
			}
			environment := Environment{"x": 5}
			job := UnaryOperatorJob{Context: context.Background(), Operator: test.operator, Value: 5, Environment: environment, Target: "x", Assignable: true}
			got, done := dispatchUnaryWorker(context.Background(), job)
			select {
			case <-done:
			case <-time.After(time.Second):
				t.Fatal("unary worker did not stop")
			}
			if got.Err != nil || got.Value != test.want {
				t.Fatalf("result = %+v, want %v", got, test.want)
			}
			if test.operator >= PrefixDecrement && environment["x"] != 4 && environment["x"] != 6 {
				t.Fatalf("final environment = %v", environment)
			}
		})
	}
}

func TestUnaryWorkersPropagateCancellationAndInvalidOperators(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	got, done := dispatchUnaryWorker(ctx, UnaryOperatorJob{Context: ctx, Operator: Positive, Value: 1, Environment: Environment{}})
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("cancelled unary worker did not stop")
	}
	if !errors.Is(got.Err, context.Canceled) {
		t.Fatalf("cancellation error = %v", got.Err)
	}

	got, done = dispatchUnaryWorker(context.Background(), UnaryOperatorJob{Operator: UnaryOperator(99), Environment: Environment{}})
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("invalid unary dispatch did not complete")
	}
	if got.Err == nil || !strings.Contains(got.Err.Error(), "unsupported unary operator") {
		t.Fatalf("invalid operator result = %+v", got)
	}
}

func TestUnaryWorkersPreserveLiteralAndRejectNonAssignableMutation(t *testing.T) {
	evaluator, err := NewEvaluator(context.Background(), 1)
	if err != nil {
		t.Fatal(err)
	}
	defer evaluator.Close()

	for _, operator := range []UnaryOperator{PrefixIncrement, PostfixIncrement, PrefixDecrement, PostfixDecrement} {
		result, err := evaluator.Submit(context.Background(), &UnaryExpression{Operator: operator, Operand: &NumberExpression{Value: 5}}, nil)
		if err != nil {
			t.Fatal(err)
		}
		got := receiveResult(t, result)
		if got.Err != nil {
			t.Fatalf("literal %d error = %v", operator, got.Err)
		}
	}

	result, err := evaluator.Submit(context.Background(), &UnaryExpression{Operator: PrefixIncrement, Operand: &UnaryExpression{Operator: Positive, Operand: &NumberExpression{Value: 5}}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	got := receiveResult(t, result)
	if got.Err == nil || !strings.Contains(got.Err.Error(), "requires an identifier or literal operand") {
		t.Fatalf("non-assignable result = %+v", got)
	}
}
