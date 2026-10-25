package calc

import (
	"context"
	"fmt"
)

const maxExpressionDepth = 1024

var errExpressionDepth = fmt.Errorf("expression nesting exceeds %d levels", maxExpressionDepth)

// dispatchEvaluate evaluates an expression through the typed job boundary.
// Recursive dispatch is synchronous and reentrant: child jobs are handled by
// the current worker instead of being sent back to a queue that it may be
// blocked waiting on.
func dispatchEvaluate(ctx context.Context, expression MathExpression, environment Environment) (float64, Environment, error) {
	result := NewResultChannel()
	dispatchExpression(ctx, expression, environment, result, 0)
	terminal, err := receiveDispatchResult(ctx, result)
	if err != nil {
		return 0, environment, err
	}
	return terminal.Value, terminal.Environment, terminal.Err
}

func receiveDispatchResult(ctx context.Context, result <-chan EvaluationResult) (EvaluationResult, error) {
	select {
	case terminal := <-result:
		return terminal, nil
	case <-ctx.Done():
		return EvaluationResult{}, ctx.Err()
	}
}

func dispatchExpression(ctx context.Context, expression MathExpression, environment Environment, result chan<- EvaluationResult, depth int) {
	if ctx == nil {
		_ = DeliverResult(context.Background(), result, EvaluationResult{Environment: environment, Err: ErrNilContext})
		return
	}
	if err := ctx.Err(); err != nil {
		_ = DeliverResult(ctx, result, EvaluationResult{Environment: environment, Err: err})
		return
	}
	if depth > maxExpressionDepth {
		_ = DeliverResult(ctx, result, EvaluationResult{Environment: environment, Err: errExpressionDepth})
		return
	}

	switch expression := expression.(type) {
	case *NumberExpression:
		handleNumber(ctx, expression, environment, result)
	case *IdentifierExpression:
		handleIdentifier(ctx, expression, environment, result)
	case *UnaryExpression:
		handleUnary(ctx, expression, environment, result, depth)
	case *BinaryExpression:
		handleBinary(ctx, expression, environment, result, depth)
	default:
		_ = DeliverResult(ctx, result, EvaluationResult{Environment: environment, Err: fmt.Errorf("unsupported expression node %T", expression)})
	}
}

func handleNumber(ctx context.Context, expression *NumberExpression, environment Environment, result chan<- EvaluationResult) {
	if expression == nil {
		_ = DeliverResult(ctx, result, EvaluationResult{Environment: environment, Err: fmt.Errorf("cannot evaluate a nil number expression")})
		return
	}
	_ = DeliverResult(ctx, result, EvaluationResult{Value: expression.Value, Environment: environment})
}

func handleIdentifier(ctx context.Context, expression *IdentifierExpression, environment Environment, result chan<- EvaluationResult) {
	if expression == nil {
		_ = DeliverResult(ctx, result, EvaluationResult{Environment: environment, Err: fmt.Errorf("cannot evaluate a nil identifier expression")})
		return
	}
	value, ok := environment[expression.Name]
	if !ok {
		_ = DeliverResult(ctx, result, EvaluationResult{Environment: environment, Err: fmt.Errorf("identifier %q is not defined", expression.Name)})
		return
	}
	_ = DeliverResult(ctx, result, EvaluationResult{Value: value, Environment: environment})
}

func handleUnary(ctx context.Context, expression *UnaryExpression, environment Environment, result chan<- EvaluationResult, depth int) {
	if expression == nil {
		_ = DeliverResult(ctx, result, EvaluationResult{Environment: environment, Err: fmt.Errorf("cannot evaluate a nil unary expression")})
		return
	}

	mutation := expression.Operator == PrefixIncrement || expression.Operator == PostfixIncrement || expression.Operator == PrefixDecrement || expression.Operator == PostfixDecrement
	job := UnaryOperatorJob{Context: ctx, Operator: expression.Operator, Environment: environment, Result: result}
	if mutation {
		if identifier, ok := expression.Operand.(*IdentifierExpression); ok {
			if identifier == nil {
				_ = DeliverResult(ctx, result, EvaluationResult{Environment: environment, Err: fmt.Errorf("cannot increment or decrement a nil identifier expression")})
				return
			}
			value, ok := environment[identifier.Name]
			if !ok {
				_ = DeliverResult(ctx, result, EvaluationResult{Environment: environment, Err: fmt.Errorf("identifier %q is not defined", identifier.Name)})
				return
			}
			job.Value, job.Target, job.Assignable = value, identifier.Name, true
			terminal, done := dispatchUnaryWorker(ctx, job)
			<-done
			_ = DeliverResult(ctx, result, terminal)
			return
		}
		if number, ok := expression.Operand.(*NumberExpression); ok && number != nil {
			job.Value = number.Value
			terminal, done := dispatchUnaryWorker(ctx, job)
			<-done
			_ = DeliverResult(ctx, result, terminal)
			return
		}
		_ = DeliverResult(ctx, result, EvaluationResult{Environment: environment, Err: fmt.Errorf("increment or decrement requires an identifier or literal operand")})
		return
	}

	child := NewResultChannel()
	dispatchExpression(ctx, expression.Operand, environment, child, depth+1)
	operand, err := receiveDispatchResult(ctx, child)
	if err != nil {
		_ = DeliverResult(ctx, result, EvaluationResult{Environment: environment, Err: err})
		return
	}
	if operand.Err != nil {
		_ = DeliverResult(ctx, result, operand)
		return
	}
	job.Context, job.Value = ctx, operand.Value
	terminal, done := dispatchUnaryWorker(ctx, job)
	<-done
	_ = DeliverResult(ctx, result, terminal)
}

func errInvalidUnaryOperator(operator UnaryOperator) error {
	return fmt.Errorf("unsupported unary operator: %d", operator)
}

func handleBinary(ctx context.Context, expression *BinaryExpression, environment Environment, result chan<- EvaluationResult, depth int) {
	if expression == nil {
		_ = DeliverResult(ctx, result, EvaluationResult{Environment: environment, Err: fmt.Errorf("cannot evaluate a nil binary expression")})
		return
	}

	leftResult := NewResultChannel()
	dispatchExpression(ctx, expression.Left, environment, leftResult, depth+1)
	left, err := receiveDispatchResult(ctx, leftResult)
	if err != nil {
		_ = DeliverResult(ctx, result, EvaluationResult{Environment: environment, Err: err})
		return
	}
	if left.Err != nil {
		_ = DeliverResult(ctx, result, left)
		return
	}
	rightResult := NewResultChannel()
	dispatchExpression(ctx, expression.Right, environment, rightResult, depth+1)
	right, err := receiveDispatchResult(ctx, rightResult)
	if err != nil {
		_ = DeliverResult(ctx, result, EvaluationResult{Environment: environment, Err: err})
		return
	}
	if right.Err != nil {
		_ = DeliverResult(ctx, result, right)
		return
	}

	var value float64
	switch expression.Operator {
	case Add:
		value = left.Value + right.Value
	case Subtract:
		value = left.Value - right.Value
	case Multiply:
		value = left.Value * right.Value
	case Divide:
		value = left.Value / right.Value
	default:
		_ = DeliverResult(ctx, result, EvaluationResult{Environment: environment, Err: fmt.Errorf("unsupported binary operator: %d", expression.Operator)})
		return
	}
	_ = DeliverResult(ctx, result, EvaluationResult{Value: value, Environment: environment})
}
