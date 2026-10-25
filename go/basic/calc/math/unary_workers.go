package calc

import "context"

type unaryWorker func(context.Context, UnaryOperatorJob) EvaluationResult

var unaryWorkers = map[UnaryOperator]struct {
	name string
	work unaryWorker
}{
	Positive:         {name: "positive", work: positiveWorker},
	Negate:           {name: "negate", work: negateWorker},
	PrefixDecrement:  {name: "prefix decrement", work: prefixDecrementWorker},
	PostfixDecrement: {name: "postfix decrement", work: postfixDecrementWorker},
	PrefixIncrement:  {name: "prefix increment", work: prefixIncrementWorker},
	PostfixIncrement: {name: "postfix increment", work: postfixIncrementWorker},
}

func unaryWorkerName(operator UnaryOperator) (string, bool) {
	worker, ok := unaryWorkers[operator]
	if !ok {
		return "", false
	}
	return worker.name, true
}

func dispatchUnaryWorker(ctx context.Context, job UnaryOperatorJob) (EvaluationResult, <-chan struct{}) {
	result := NewResultChannel()
	done := make(chan struct{})
	jobs := make(chan UnaryOperatorJob, 1)
	if ctx == nil {
		return EvaluationResult{Environment: job.Environment, Err: ErrNilContext}, closedSignal()
	}
	worker, ok := unaryWorkers[job.Operator]
	if !ok {
		return EvaluationResult{Environment: job.Environment, Err: errInvalidUnaryOperator(job.Operator)}, closedSignal()
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

func closedSignal() <-chan struct{} {
	done := make(chan struct{})
	close(done)
	return done
}

func positiveWorker(_ context.Context, job UnaryOperatorJob) EvaluationResult {
	return EvaluationResult{Value: job.Value, Environment: job.Environment}
}

func negateWorker(_ context.Context, job UnaryOperatorJob) EvaluationResult {
	return EvaluationResult{Value: -job.Value, Environment: job.Environment}
}

func prefixDecrementWorker(_ context.Context, job UnaryOperatorJob) EvaluationResult {
	return mutateWorker(job, -1, true)
}

func postfixDecrementWorker(_ context.Context, job UnaryOperatorJob) EvaluationResult {
	return mutateWorker(job, -1, false)
}

func prefixIncrementWorker(_ context.Context, job UnaryOperatorJob) EvaluationResult {
	return mutateWorker(job, 1, true)
}

func postfixIncrementWorker(_ context.Context, job UnaryOperatorJob) EvaluationResult {
	return mutateWorker(job, 1, false)
}

func mutateWorker(job UnaryOperatorJob, delta float64, prefix bool) EvaluationResult {
	updated := job.Value + delta
	if job.Assignable {
		job.Environment[job.Target] = updated
	}
	value := job.Value
	if prefix {
		value = updated
	}
	return EvaluationResult{Value: value, Environment: job.Environment}
}
