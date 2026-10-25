package calc

import "context"

type binaryWorker func(context.Context, BinaryOperatorJob) EvaluationResult

var binaryWorkers = map[BinaryOperator]struct {
	name string
	work binaryWorker
}{
	Add:      {name: "add", work: addWorker},
	Subtract: {name: "subtract", work: subtractWorker},
	Multiply: {name: "multiply", work: multiplyWorker},
	Divide:   {name: "divide", work: divideWorker},
}

func binaryWorkerName(operator BinaryOperator) (string, bool) {
	worker, ok := binaryWorkers[operator]
	if !ok {
		return "", false
	}
	return worker.name, true
}

func dispatchBinaryWorker(ctx context.Context, job BinaryOperatorJob) (EvaluationResult, <-chan struct{}) {
	result := NewResultChannel()
	done := make(chan struct{})
	jobs := make(chan BinaryOperatorJob, 1)
	if ctx == nil {
		return EvaluationResult{Environment: job.Environment, Err: ErrNilContext}, closedSignal()
	}
	worker, ok := binaryWorkers[job.Operator]
	if !ok {
		return EvaluationResult{Environment: job.Environment, Err: errInvalidBinaryOperator(job.Operator)}, closedSignal()
	}

	go func() {
		defer close(done)
		select {
		case received := <-jobs:
			_ = DeliverResult(received.Context, result, worker.work(received.Context, received))
		case <-ctx.Done():
		}
	}()

	select {
	case jobs <- job:
	case <-ctx.Done():
	}
	select {
	case terminal := <-result:
		return terminal, done
	case <-ctx.Done():
		return EvaluationResult{Environment: job.Environment, Err: ctx.Err()}, done
	}
}

func addWorker(_ context.Context, job BinaryOperatorJob) EvaluationResult {
	return EvaluationResult{Value: job.Left + job.Right, Environment: job.Environment}
}

func subtractWorker(_ context.Context, job BinaryOperatorJob) EvaluationResult {
	return EvaluationResult{Value: job.Left - job.Right, Environment: job.Environment}
}

func multiplyWorker(_ context.Context, job BinaryOperatorJob) EvaluationResult {
	return EvaluationResult{Value: job.Left * job.Right, Environment: job.Environment}
}

// divideWorker preserves the floating-point division semantics of Evaluator.
// In particular, zero divisors produce the corresponding IEEE-754 value.
func divideWorker(_ context.Context, job BinaryOperatorJob) EvaluationResult {
	return EvaluationResult{Value: job.Left / job.Right, Environment: job.Environment}
}
