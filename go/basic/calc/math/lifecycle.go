package calc

import (
	"context"
	"errors"
	"sync"
)

var (
	ErrEvaluatorClosed      = errors.New("evaluator is closed")
	ErrInvalidQueueCapacity = errors.New("evaluator queue capacity must be positive")
)

const DefaultQueueCapacity = 16

type evaluatorRequest struct {
	context     context.Context
	expression  MathExpression
	environment Environment
	result      chan EvaluationResult
}

// EvaluatorRuntime owns a bounded root queue and one worker. It deliberately
// delegates computation to the sequential evaluator until later steps add
// type-specific workers.
type EvaluatorRuntime struct {
	context context.Context
	cancel  context.CancelFunc
	queue   chan evaluatorRequest
	space   chan struct{}
	done    chan struct{}
	mu      sync.Mutex
	closed  bool
	workers sync.WaitGroup
}

// NewEvaluator starts an evaluator with a bounded root queue. The caller owns
// the returned result channels; Close is idempotent and waits for the worker.
func NewEvaluator(parent context.Context, capacity int) (*EvaluatorRuntime, error) {
	if parent == nil {
		parent = context.Background()
	}
	if capacity <= 0 {
		return nil, ErrInvalidQueueCapacity
	}
	ctx, cancel := context.WithCancel(parent)
	e := &EvaluatorRuntime{
		context: ctx, cancel: cancel, queue: make(chan evaluatorRequest, capacity),
		space: make(chan struct{}, 1), done: make(chan struct{}),
	}
	e.workers.Add(1)
	go e.worker()
	return e, nil
}

// Submit admits one request, blocking only while the root queue is full. The
// caller context and evaluator shutdown both interrupt admission.
func (e *EvaluatorRuntime) Submit(ctx context.Context, expression MathExpression, input Environment) (<-chan EvaluationResult, error) {
	if ctx == nil {
		return nil, ErrNilContext
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if expression == nil {
		return nil, ErrNilExpression
	}
	result := NewResultChannel()
	request := evaluatorRequest{context: ctx, expression: expression, environment: input.Clone(), result: result}
	for {
		e.mu.Lock()
		if e.closed || e.context.Err() != nil {
			e.mu.Unlock()
			return nil, ErrEvaluatorClosed
		}
		select {
		case e.queue <- request:
			e.mu.Unlock()
			return result, nil
		default:
			e.mu.Unlock()
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-e.done:
			return nil, ErrEvaluatorClosed
		case <-e.context.Done():
			return nil, ErrEvaluatorClosed
		case <-e.space:
		}
	}
}

// Close prevents new admission, cancels owned context, drains admitted jobs
// with ErrEvaluatorClosed, and waits until the worker has stopped. Requests
// already being evaluated finish their sequential computation before Close
// returns; their result delivery remains request-context aware.
func (e *EvaluatorRuntime) Close() error {
	e.mu.Lock()
	if !e.closed {
		e.closed = true
		close(e.done)
		e.cancel()
	}
	e.mu.Unlock()
	e.workers.Wait()
	return nil
}

func (e *EvaluatorRuntime) worker() {
	defer e.workers.Done()
	for {
		select {
		case request := <-e.queue:
			e.signalSpace()
			e.run(request)
		case <-e.done:
			e.drain()
			return
		case <-e.context.Done():
			e.mu.Lock()
			e.closed = true
			e.mu.Unlock()
			e.drain()
			return
		}
	}
}

func (e *EvaluatorRuntime) drain() {
	for {
		select {
		case request := <-e.queue:
			e.signalSpace()
			_ = DeliverResult(request.context, request.result, EvaluationResult{Environment: request.environment, Err: ErrEvaluatorClosed})
		default:
			return
		}
	}
}

func (e *EvaluatorRuntime) run(request evaluatorRequest) {
	value, environment, err := (Evaluator{}).EvaluateWithEnvironment(request.expression, request.environment)
	if e.isClosed() {
		err = ErrEvaluatorClosed
	}
	_ = DeliverResult(request.context, request.result, EvaluationResult{Value: value, Environment: environment, Err: err})
}

func (e *EvaluatorRuntime) isClosed() bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.closed
}

func (e *EvaluatorRuntime) signalSpace() {
	select {
	case e.space <- struct{}{}:
	default:
	}
}
