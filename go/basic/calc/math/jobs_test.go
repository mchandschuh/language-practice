package calc

import (
	"context"
	"errors"
	"testing"
)

func TestTypedJobsPreserveRequest(t *testing.T) {
	ctx := context.WithValue(context.Background(), "request", "one")
	expression := &BinaryExpression{Operator: Add}
	environment := Environment{"x": 2}
	result := NewResultChannel()
	job := ExpressionJob{Context: ctx, Expression: expression, Environment: environment, Result: result}

	if err := job.Validate(); err != nil {
		t.Fatalf("Validate returned error: %v", err)
	}
	if job.Context != ctx || job.Expression != expression || job.Result != result || job.Environment["x"] != 2 {
		t.Fatalf("job did not preserve request fields: %#v", job)
	}

	if err := (UnaryJob{Context: ctx, Expression: &UnaryExpression{Operand: &NumberExpression{Value: 1}}, Environment: environment, Result: result}).Validate(); err != nil {
		t.Fatalf("UnaryJob validation failed: %v", err)
	}
	if err := (BinaryJob{Context: ctx, Expression: expression, Environment: environment, Result: result}).Validate(); err != nil {
		t.Fatalf("BinaryJob validation failed: %v", err)
	}
	if err := (NumberJob{Context: ctx, Expression: &NumberExpression{Value: 1}, Environment: environment, Result: result}).Validate(); err != nil {
		t.Fatalf("NumberJob validation failed: %v", err)
	}
	if err := (IdentifierJob{Context: ctx, Expression: &IdentifierExpression{Name: "x"}, Environment: environment, Result: result}).Validate(); err != nil {
		t.Fatalf("IdentifierJob validation failed: %v", err)
	}
	if err := (UnaryOperatorJob{Context: ctx, Environment: environment, Result: result}).Validate(); err != nil {
		t.Fatalf("UnaryOperatorJob validation failed: %v", err)
	}
	if err := (BinaryOperatorJob{Context: ctx, Environment: environment, Result: result}).Validate(); err != nil {
		t.Fatalf("BinaryOperatorJob validation failed: %v", err)
	}
}

func TestEvaluationResultCarriesValueStateAndError(t *testing.T) {
	state := Environment{"x": 3}
	result := EvaluationResult{Value: 7, Environment: state}
	if result.Value != 7 || result.Environment["x"] != 3 || result.Err != nil {
		t.Fatalf("success result = %#v", result)
	}

	want := errors.New("division by zero")
	failed := EvaluationResult{Err: want}
	if !errors.Is(failed.Err, want) {
		t.Fatalf("failure result lost error: %#v", failed)
	}
}

func TestDeliverResultIsBufferedOneShotAndCancellationAware(t *testing.T) {
	destination := NewResultChannel()
	want := EvaluationResult{Value: 42}
	if err := DeliverResult(context.Background(), destination, want); err != nil {
		t.Fatalf("DeliverResult returned error: %v", err)
	}
	select {
	case got := <-destination:
		if got.Value != want.Value {
			t.Fatalf("result = %#v, want %#v", got, want)
		}
	default:
		t.Fatal("buffered result was not immediately available")
	}

	canceled, cancel := context.WithCancel(context.Background())
	cancel()
	if err := DeliverResult(canceled, NewResultChannel(), EvaluationResult{}); !errors.Is(err, context.Canceled) {
		t.Fatalf("canceled delivery error = %v, want context.Canceled", err)
	}
}

func TestTypedJobValidationRejectsInvalidInputs(t *testing.T) {
	ctx := context.Background()
	result := NewResultChannel()
	environment := Environment{}
	tests := []struct {
		name string
		job  ExpressionJob
		want error
	}{
		{name: "nil expression", job: ExpressionJob{Context: ctx, Environment: environment, Result: result}, want: ErrNilExpression},
		{name: "nil context", job: ExpressionJob{Context: nil, Expression: &NumberExpression{}, Environment: environment, Result: result}, want: ErrNilContext},
		{name: "nil result", job: ExpressionJob{Context: ctx, Expression: &NumberExpression{}, Environment: environment}, want: ErrNilResult},
		{name: "nil environment", job: ExpressionJob{Context: ctx, Expression: &NumberExpression{}, Result: result}, want: ErrNilEnvironment},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := test.job.Validate(); err == nil || err.Error() != test.want.Error() {
				t.Fatalf("Validate error = %v, want %v", err, test.want)
			}
		})
	}
}
