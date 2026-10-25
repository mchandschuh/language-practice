package calc

import (
	"reflect"
	"testing"
)

func TestParseExpressions(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		want       MathExpression
	}{
		{
			name:       "number",
			expression: "42",
			want:       &NumberExpression{Value: 42},
		},
		{
			name:       "identifier",
			expression: "value2",
			want:       &IdentifierExpression{Name: "value2"},
		},
		{
			name:       "underscored identifier",
			expression: "_value",
			want:       &IdentifierExpression{Name: "_value"},
		},
		{
			name:       "decimal with exponent",
			expression: "1.25e2",
			want:       &NumberExpression{Value: 125},
		},
		{
			name:       "addition",
			expression: "1 + 2",
			want: &BinaryExpression{
				Operator: Add,
				Left:     &NumberExpression{Value: 1},
				Right:    &NumberExpression{Value: 2},
			},
		},
		{
			name:       "precedence",
			expression: "1 + 2 * 3",
			want: &BinaryExpression{
				Operator: Add,
				Left:     &NumberExpression{Value: 1},
				Right: &BinaryExpression{
					Operator: Multiply,
					Left:     &NumberExpression{Value: 2},
					Right:    &NumberExpression{Value: 3},
				},
			},
		},
		{
			name:       "parentheses override precedence",
			expression: "(1 + 2) * 3",
			want: &BinaryExpression{
				Operator: Multiply,
				Left: &BinaryExpression{
					Operator: Add,
					Left:     &NumberExpression{Value: 1},
					Right:    &NumberExpression{Value: 2},
				},
				Right: &NumberExpression{Value: 3},
			},
		},
		{
			name:       "left associative subtraction",
			expression: "10 - 3 - 2",
			want: &BinaryExpression{
				Operator: Subtract,
				Left: &BinaryExpression{
					Operator: Subtract,
					Left:     &NumberExpression{Value: 10},
					Right:    &NumberExpression{Value: 3},
				},
				Right: &NumberExpression{Value: 2},
			},
		},
		{
			name:       "unary minus",
			expression: "-2 * 3",
			want: &BinaryExpression{
				Operator: Multiply,
				Left:     &UnaryExpression{Operator: Negate, Operand: &NumberExpression{Value: 2}},
				Right:    &NumberExpression{Value: 3},
			},
		},
		{
			name:       "unary plus",
			expression: "+5",
			want: &UnaryExpression{
				Operator: Positive,
				Operand:  &NumberExpression{Value: 5},
			},
		},
		{
			name:       "nested and spaced expression",
			expression: " \t(2.5 + 1.5) / (3 - 1)\n",
			want: &BinaryExpression{
				Operator: Divide,
				Left: &BinaryExpression{
					Operator: Add,
					Left:     &NumberExpression{Value: 2.5},
					Right:    &NumberExpression{Value: 1.5},
				},
				Right: &BinaryExpression{
					Operator: Subtract,
					Left:     &NumberExpression{Value: 3},
					Right:    &NumberExpression{Value: 1},
				},
			},
		},
		{
			name:       "prefix increment",
			expression: "++5",
			want: &UnaryExpression{
				Operator: PrefixIncrement,
				Operand:  &NumberExpression{Value: 5},
			},
		},
		{
			name:       "postfix increment",
			expression: "5++",
			want: &UnaryExpression{
				Operator: PostfixIncrement,
				Operand:  &NumberExpression{Value: 5},
			},
		},
		{
			name:       "prefix decrement",
			expression: "--5",
			want: &UnaryExpression{
				Operator: PrefixDecrement,
				Operand:  &NumberExpression{Value: 5},
			},
		},
		{
			name:       "postfix decrement",
			expression: "5--",
			want: &UnaryExpression{
				Operator: PostfixDecrement,
				Operand:  &NumberExpression{Value: 5},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := parse(test.expression)
			if err != nil {
				t.Fatalf("parse(%q) returned an unexpected error: %v", test.expression, err)
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("parse(%q) = %#v, want %#v", test.expression, got, test.want)
			}
		})
	}
}

func TestParseRejectsInvalidExpressions(t *testing.T) {
	tests := []string{
		"",
		" ",
		"1 +",
		"+",
		"1 2",
		"1 ** 2",
		"1 / / 2",
		"(1 + 2",
		"1 + 2)",
		".",
		"1e",
		"not a number",
		"2x",
		"$x",
	}

	for _, expression := range tests {
		t.Run(expression, func(t *testing.T) {
			if _, err := parse(expression); err == nil {
				t.Fatalf("parse(%q) returned nil error", expression)
			}
		})
	}
}
