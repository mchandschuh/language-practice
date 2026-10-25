package calc

import (
	"context"
	"errors"
)

var (
	ErrNilContext      = errors.New("evaluation job has nil context")
	ErrNilExpression   = errors.New("evaluation job has nil expression")
	ErrNilResult       = errors.New("evaluation job has nil result channel")
	ErrNilEnvironment  = errors.New("evaluation job has nil environment")
	ErrInvalidOperator = errors.New("evaluation job has invalid operator")
)

// EvaluationResult is the terminal value of one request. Err is the source of
// truth for failed evaluation; workers must not replace it with logging.
// Environment is the request-owned final state when evaluation mutates it.
type EvaluationResult struct {
	Value       float64
	Environment Environment
	Err         error
}

// NewResultChannel creates the one-shot destination owned by the request
// submitter. Its capacity lets a worker deliver one result without racing a
// receiver. The channel is intentionally never closed by a worker.
func NewResultChannel() chan EvaluationResult {
	return make(chan EvaluationResult, 1)
}

// ExpressionJob is the root message sent by a public evaluator to a type
// dispatcher. The dispatcher receives it; the submitting request owns Result.
type ExpressionJob struct {
	Context     context.Context
	Expression  MathExpression
	Environment Environment
	Result      chan<- EvaluationResult
}

// NumberJob evaluates a numeric expression.
type NumberJob struct {
	Context     context.Context
	Expression  *NumberExpression
	Environment Environment
	Result      chan<- EvaluationResult
}

// IdentifierJob evaluates an identifier expression.
type IdentifierJob struct {
	Context     context.Context
	Expression  *IdentifierExpression
	Environment Environment
	Result      chan<- EvaluationResult
}

// UnaryJob evaluates a unary expression and is consumed by unary dispatch.
type UnaryJob struct {
	Context     context.Context
	Expression  *UnaryExpression
	Environment Environment
	Result      chan<- EvaluationResult
}

// BinaryJob evaluates a binary expression and is consumed by binary dispatch.
type BinaryJob struct {
	Context     context.Context
	Expression  *BinaryExpression
	Environment Environment
	Result      chan<- EvaluationResult
}

// UnaryOperatorJob is the typed message for a unary operator worker.
type UnaryOperatorJob struct {
	Context     context.Context
	Operator    UnaryOperator
	Value       float64
	Environment Environment
	Result      chan<- EvaluationResult
}

// BinaryOperatorJob is the typed message for a binary operator worker.
type BinaryOperatorJob struct {
	Context     context.Context
	Operator    BinaryOperator
	Left        float64
	Right       float64
	Environment Environment
	Result      chan<- EvaluationResult
}

func validateCommon(ctx context.Context, environment Environment, result chan<- EvaluationResult) error {
	if ctx == nil {
		return ErrNilContext
	}
	if result == nil {
		return ErrNilResult
	}
	if environment == nil {
		return ErrNilEnvironment
	}
	return nil
}

func (job ExpressionJob) Validate() error {
	if job.Expression == nil {
		return ErrNilExpression
	}
	if err := validateCommon(job.Context, job.Environment, job.Result); err != nil {
		return err
	}
	return nil
}

func (job NumberJob) Validate() error {
	if job.Expression == nil {
		return ErrNilExpression
	}
	return validateCommon(job.Context, job.Environment, job.Result)
}

func (job IdentifierJob) Validate() error {
	if job.Expression == nil {
		return ErrNilExpression
	}
	return validateCommon(job.Context, job.Environment, job.Result)
}

func (job UnaryJob) Validate() error {
	if job.Expression == nil {
		return ErrNilExpression
	}
	return validateCommon(job.Context, job.Environment, job.Result)
}

func (job BinaryJob) Validate() error {
	if job.Expression == nil {
		return ErrNilExpression
	}
	return validateCommon(job.Context, job.Environment, job.Result)
}

func (job UnaryOperatorJob) Validate() error {
	if job.Operator > PostfixIncrement {
		return ErrInvalidOperator
	}
	return validateCommon(job.Context, job.Environment, job.Result)
}

func (job BinaryOperatorJob) Validate() error {
	if job.Operator > Divide {
		return ErrInvalidOperator
	}
	return validateCommon(job.Context, job.Environment, job.Result)
}

// DeliverResult sends exactly one terminal result unless ctx is already done
// or becomes done while the destination is unavailable. It creates no
// goroutines and never closes destination, because the caller owns that
// channel.
func DeliverResult(ctx context.Context, destination chan<- EvaluationResult, result EvaluationResult) error {
	if ctx == nil {
		return ErrNilContext
	}
	if destination == nil {
		return ErrNilResult
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	select {
	case destination <- result:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
