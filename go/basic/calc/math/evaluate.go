package calc

import "fmt"

// Evaluate evaluates an arithmetic expression tree.
func Evaluate(expression MathExpression) (float64, error) {
	switch expression := expression.(type) {
	case *NumberExpression:
		return expression.Value, nil
	case UnaryExpression:
		operand, err := Evaluate(expression.Operand)
		if err != nil {
			return 0, err
		}
		switch expression.Operator {
		case Positive, PrefixIncrement, PostfixIncrement:
			if expression.Operator == Positive {
				return operand, nil
			}
			return operand + 1, nil
		case Negate, PrefixDecrement, PostfixDecrement:
			if expression.Operator == Negate {
				return -operand, nil
			}
			return operand - 1, nil
		default:
			return 0, fmt.Errorf("unsupported unary operator: %d", expression.Operator)
		}
	case BinaryExpression:
		left, err := Evaluate(expression.Left)
		if err != nil {
			return 0, err
		}
		right, err := Evaluate(expression.Right)
		if err != nil {
			return 0, err
		}
		switch expression.Operator {
		case Add:
			return left + right, nil
		case Subtract:
			return left - right, nil
		case Multiply:
			return left * right, nil
		case Divide:
			return left / right, nil
		default:
			return 0, fmt.Errorf("unsupported binary operator: %d", expression.Operator)
		}
	default:
		return 0, fmt.Errorf("unsupported expression type %T", expression)
	}
}
