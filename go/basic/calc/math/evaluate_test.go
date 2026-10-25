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

			got, err := expression.Evaluate()
			if err != nil {
				t.Fatalf("Evaluate(Parse(%q)) returned an unexpected error: %v", test.expression, err)
			}
			if got != test.want {
				t.Errorf("Evaluate(Parse(%q)) = %v, want %v", test.expression, got, test.want)
			}
		})
	}
}
