package calc

import "testing"

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
