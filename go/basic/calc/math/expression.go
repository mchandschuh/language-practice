package calc

import "fmt"

type MathExpression interface {
	expressionNode()
	Evaluate() (float64, error)
}

type NumberExpression struct {
	Value float64
}

func (e NumberExpression) Evaluate() (float64, error) {
	return e.Value, nil
}

func (*NumberExpression) expressionNode() {}

type UnaryExpression struct {
	Operator UnaryOperator
	Operand  MathExpression
}

func (UnaryExpression) expressionNode() {}

func (e UnaryExpression) Evaluate() (float64, error) {
	value, err := e.Operand.Evaluate()
	if err != nil {
		return 0, err
	}

	switch e.Operator {
	case Positive:
		return value, nil
	case Negate:
		return -value, nil
	case PrefixDecrement, PostfixDecrement:
		return value - 1, nil
	case PrefixIncrement, PostfixIncrement:
		return value + 1, nil
	}

	return 0, fmt.Errorf("Invalid UnaryOperator: %d", e.Operator)
}

type BinaryExpression struct {
	Operator BinaryOperator
	Left     MathExpression
	Right    MathExpression
}

func (BinaryExpression) expressionNode() {}

func (e BinaryExpression) Evaluate() (float64, error) {
	left, err := e.Left.Evaluate()
	if err != nil {
		return 0, err
	}
	right, err := e.Right.Evaluate()
	if err != nil {
		return 0, err
	}
	switch e.Operator {
	case Add:
		return left + right, nil
	case Subtract:
		return left - right, nil
	case Multiply:
		return left * right, nil
	case Divide:
		return left / right, nil
	}

	return 0, fmt.Errorf("Invalid operator: %d", e.Operator)
}

type MathExpressionJob struct {
	Result     chan float64
	Expression MathExpression
}
