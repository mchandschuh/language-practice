package calc

import "fmt"

// Evaluator evaluates an immutable expression tree sequentially. Evaluator is
// intentionally stateless, so one instance can be reused for multiple trees.
type Evaluator struct{}

// Evaluate returns the value of expression or an error for a nil, unsupported,
// or malformed expression node. Binary operands are evaluated left first and
// right second.
func (e Evaluator) Evaluate(expression MathExpression) (float64, error) {
	switch expression := expression.(type) {
	case *NumberExpression:
		if expression == nil {
			return 0, fmt.Errorf("cannot evaluate a nil number expression")
		}
		return expression.Value, nil
	case *UnaryExpression:
		if expression == nil {
			return 0, fmt.Errorf("cannot evaluate a nil unary expression")
		}
		value, err := e.Evaluate(expression.Operand)
		if err != nil {
			return 0, err
		}
		switch expression.Operator {
		case Positive:
			return value, nil
		case Negate:
			return -value, nil
		case PrefixDecrement, PostfixDecrement:
			return value - 1, nil
		case PrefixIncrement, PostfixIncrement:
			return value + 1, nil
		default:
			return 0, fmt.Errorf("unsupported unary operator: %d", expression.Operator)
		}
	case *BinaryExpression:
		if expression == nil {
			return 0, fmt.Errorf("cannot evaluate a nil binary expression")
		}
		left, err := e.Evaluate(expression.Left)
		if err != nil {
			return 0, err
		}
		right, err := e.Evaluate(expression.Right)
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
		return 0, fmt.Errorf("unsupported expression node %T", expression)
	}
}

// Evaluate is the convenience form of Evaluator{}.Evaluate.
func Evaluate(expression MathExpression) (float64, error) {
	return (Evaluator{}).Evaluate(expression)
}
