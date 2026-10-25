package calc

import (
	"strings"
	"testing"
)

func TestEvaluateExpressions(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		want       float64
	}{
		{name: "number", expression: "42", want: 42},
		{name: "precedence", expression: "1 + 2 * 3", want: 7},
		{name: "parentheses", expression: "(1 + 2) * 3", want: 9},
		{name: "unary operators", expression: "-2 + +3", want: 1},
		{name: "prefix increment", expression: "++5", want: 6},
		{name: "postfix increment", expression: "5++", want: 6},
		{name: "prefix decrement", expression: "--5", want: 4},
		{name: "postfix decrement", expression: "5--", want: 4},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expression, err := Parse(test.expression)
			if err != nil {
				t.Fatalf("Parse(%q) returned an unexpected error: %v", test.expression, err)
			}

			got, err := Evaluate(expression)
			if err != nil {
				t.Fatalf("Evaluate(Parse(%q)) returned an unexpected error: %v", test.expression, err)
			}
			if got != test.want {
				t.Errorf("Evaluate(Parse(%q)) = %v, want %v", test.expression, got, test.want)
			}
		})
	}
}

func TestEvaluateRejectsUnsupportedNodesAndOperators(t *testing.T) {
	tests := []MathExpression{
		(*NumberExpression)(nil),
		&UnaryExpression{Operator: UnaryOperator(99), Operand: &NumberExpression{Value: 1}},
		&BinaryExpression{Operator: BinaryOperator(99), Left: &NumberExpression{Value: 1}, Right: &NumberExpression{Value: 2}},
		&UnaryExpression{Operator: Positive},
		&BinaryExpression{Operator: Add, Left: &NumberExpression{Value: 1}},
	}

	for _, expression := range tests {
		t.Run("unsupported", func(t *testing.T) {
			if _, err := Evaluate(expression); err == nil {
				t.Fatalf("Evaluate(%T) returned nil error", expression)
			}
		})
	}
}

func TestEvaluateWithEnvironmentReusesASTAndIsolatesInput(t *testing.T) {
	expression, err := Parse("x * x + 2")
	if err != nil {
		t.Fatalf("Parse returned an unexpected error: %v", err)
	}

	input := Environment{"x": 3}
	got, final, err := EvaluateWithEnvironment(expression, input)
	if err != nil {
		t.Fatalf("first evaluation returned an unexpected error: %v", err)
	}
	if got != 11 {
		t.Errorf("first evaluation = %v, want 11", got)
	}
	if input["x"] != 3 || final["x"] != 3 {
		t.Errorf("first evaluation changed state unexpectedly: input=%v final=%v", input, final)
	}

	input["x"] = 10
	got, final, err = EvaluateWithEnvironment(expression, input)
	if err != nil {
		t.Fatalf("second evaluation returned an unexpected error: %v", err)
	}
	if got != 102 {
		t.Errorf("second evaluation = %v, want 102", got)
	}
	if final["x"] != 10 {
		t.Errorf("second final x = %v, want 10", final["x"])
	}

	final["x"] = 99
	if input["x"] != 10 {
		t.Errorf("returned final environment aliases input: input=%v", input)
	}
}

func TestEvaluateWithEnvironmentIdentifierMutation(t *testing.T) {
	tests := []struct {
		expression string
		wantValue  float64
		wantFinal  float64
	}{
		{expression: "++x", wantValue: 6, wantFinal: 6},
		{expression: "x++", wantValue: 5, wantFinal: 6},
		{expression: "--x", wantValue: 4, wantFinal: 4},
		{expression: "x--", wantValue: 5, wantFinal: 4},
	}

	for _, test := range tests {
		t.Run(test.expression, func(t *testing.T) {
			expression, err := Parse(test.expression)
			if err != nil {
				t.Fatalf("Parse returned an unexpected error: %v", err)
			}
			input := Environment{"x": 5}
			got, final, err := EvaluateWithEnvironment(expression, input)
			if err != nil {
				t.Fatalf("evaluation returned an unexpected error: %v", err)
			}
			if got != test.wantValue || final["x"] != test.wantFinal {
				t.Errorf("evaluation = (%v, %v), want (%v, %v)", got, final["x"], test.wantValue, test.wantFinal)
			}
			if input["x"] != 5 {
				t.Errorf("evaluation mutated caller input: %v", input)
			}
		})
	}
}

func TestEvaluateWithEnvironmentRejectsMissingAndNonAssignableIdentifiers(t *testing.T) {
	expression, err := Parse("x + 1")
	if err != nil {
		t.Fatalf("Parse returned an unexpected error: %v", err)
	}
	if _, _, err := EvaluateWithEnvironment(expression, Environment{}); err == nil || !strings.Contains(err.Error(), `identifier "x" is not defined`) {
		t.Fatalf("missing identifier error = %v, want stable identifier error", err)
	}

	for _, expressionText := range []string{"++(x + 1)"} {
		expression, err := Parse(expressionText)
		if err != nil {
			t.Fatalf("Parse(%q) returned an unexpected error: %v", expressionText, err)
		}
		if _, _, err := EvaluateWithEnvironment(expression, Environment{}); err == nil {
			t.Errorf("EvaluateWithEnvironment(%q) returned nil error", expressionText)
		}
	}

	expression, err = Parse("++(x + 1)")
	if err != nil {
		t.Fatalf("Parse returned an unexpected error: %v", err)
	}
	if _, _, err := EvaluateWithEnvironment(expression, Environment{"x": 1}); err == nil {
		t.Fatal("non-assignable increment returned nil error")
	}
}
