package calc

import "fmt"

// Environment contains the input values for one evaluation request.
// Evaluator copies it before evaluating, so callers retain ownership of their
// map and can safely reuse it for another request.
type Environment map[string]float64

// Clone returns an independent copy of the environment.
func (environment Environment) Clone() Environment {
	clone := make(Environment, len(environment))
	for name, value := range environment {
		clone[name] = value
	}
	return clone
}

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

// EvaluateWithEnvironment evaluates expression using an isolated copy of
// input and returns both the value and the final request state. Identifier
// increments and decrements mutate only that private state. Literal
// increments and decrements retain the numeric exercise semantics established
// by Evaluate; other non-assignable increment operands return an error.
func (e Evaluator) EvaluateWithEnvironment(expression MathExpression, input Environment) (float64, Environment, error) {
	environment := input.Clone()
	value, err := e.evaluate(expression, environment)
	return value, environment, err
}

func (e Evaluator) evaluate(expression MathExpression, environment Environment) (float64, error) {
	switch expression := expression.(type) {
	case *NumberExpression:
		if expression == nil {
			return 0, fmt.Errorf("cannot evaluate a nil number expression")
		}
		return expression.Value, nil
	case *IdentifierExpression:
		if expression == nil {
			return 0, fmt.Errorf("cannot evaluate a nil identifier expression")
		}
		value, ok := environment[expression.Name]
		if !ok {
			return 0, fmt.Errorf("identifier %q is not defined", expression.Name)
		}
		return value, nil
	case *UnaryExpression:
		if expression == nil {
			return 0, fmt.Errorf("cannot evaluate a nil unary expression")
		}
		if expression.Operator == PrefixIncrement || expression.Operator == PostfixIncrement || expression.Operator == PrefixDecrement || expression.Operator == PostfixDecrement {
			if identifier, ok := expression.Operand.(*IdentifierExpression); ok {
				if identifier == nil {
					return 0, fmt.Errorf("cannot increment or decrement a nil identifier expression")
				}
				value, ok := environment[identifier.Name]
				if !ok {
					return 0, fmt.Errorf("identifier %q is not defined", identifier.Name)
				}
				updated := value
				if expression.Operator == PrefixIncrement || expression.Operator == PostfixIncrement {
					updated++
				} else {
					updated--
				}
				environment[identifier.Name] = updated
				if expression.Operator == PostfixIncrement || expression.Operator == PostfixDecrement {
					return value, nil
				}
				return updated, nil
			}
			if number, ok := expression.Operand.(*NumberExpression); ok && number != nil {
				return e.Evaluate(expression)
			}
			return 0, fmt.Errorf("increment or decrement requires an identifier or literal operand")
		}
		value, err := e.evaluate(expression.Operand, environment)
		if err != nil {
			return 0, err
		}
		switch expression.Operator {
		case Positive:
			return value, nil
		case Negate:
			return -value, nil
		default:
			return 0, fmt.Errorf("unsupported unary operator: %d", expression.Operator)
		}
	case *BinaryExpression:
		if expression == nil {
			return 0, fmt.Errorf("cannot evaluate a nil binary expression")
		}
		left, err := e.evaluate(expression.Left, environment)
		if err != nil {
			return 0, err
		}
		right, err := e.evaluate(expression.Right, environment)
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

// EvaluateWithEnvironment evaluates expression with input values using a
// stateless Evaluator.
func EvaluateWithEnvironment(expression MathExpression, input Environment) (float64, Environment, error) {
	return (Evaluator{}).EvaluateWithEnvironment(expression, input)
}

// Evaluate is the convenience form of Evaluator{}.Evaluate.
func Evaluate(expression MathExpression) (float64, error) {
	return (Evaluator{}).Evaluate(expression)
}
