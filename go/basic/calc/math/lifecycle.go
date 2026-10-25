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
	caller      context.Context
	cancel      context.CancelFunc
	stopCaller  func() bool
	expression  MathExpression
	environment Environment
	result      chan EvaluationResult
}

// EvaluatorRuntime owns bounded root admission and one worker. The worker
// delegates expression computation to the channel-based dispatcher.
type EvaluatorRuntime struct {
	context context.Context
	cancel  context.CancelFunc
	queue   chan evaluatorRequest
	slots   chan struct{}
	done    chan struct{}
	mu      sync.Mutex
	closed  bool
	workers sync.WaitGroup
	// beforeRun is only used by package tests to hold a worker at a
	// deterministic barrier. It is nil in normal operation.
	beforeRun func()
}

// NewEvaluator starts an evaluator with one worker and a bounded admission
// capacity. Capacity counts all admitted requests, including the request
// currently being evaluated; it is not merely the buffered queue length. The
// caller owns the returned result channels; Close is idempotent and waits for
// the worker.
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
		slots: make(chan struct{}, capacity), done: make(chan struct{}),
	}
	e.workers.Add(1)
	go e.worker()
	return e, nil
}

// Submit admits one request, blocking only while the configured admission
// capacity is full. The caller context and evaluator shutdown both interrupt
// admission. A successfully returned result channel represents exactly one
// admitted request.
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
	requestContext, cancel := context.WithCancel(e.context)
	stopCaller := context.AfterFunc(ctx, cancel)
	request := evaluatorRequest{
		context: requestContext, caller: ctx, cancel: cancel, stopCaller: stopCaller,
		expression: expression, environment: input.Clone(), result: result,
	}
	select {
	case e.slots <- struct{}{}:
	case <-ctx.Done():
		stopCaller()
		cancel()
		return nil, ctx.Err()
	case <-e.done:
		stopCaller()
		cancel()
		return nil, ErrEvaluatorClosed
	case <-e.context.Done():
		stopCaller()
		cancel()
		return nil, ErrEvaluatorClosed
	}

	e.mu.Lock()
	defer e.mu.Unlock()
	if e.closed || e.context.Err() != nil {
		request.stopCaller()
		request.cancel()
		e.releaseSlot()
		return nil, ErrEvaluatorClosed
	}
	e.queue <- request
	return result, nil
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
			request.stopCaller()
			request.cancel()
			deliverTerminal(request.result, EvaluationResult{Environment: request.environment, Err: ErrEvaluatorClosed})
			e.releaseSlot()
		default:
			return
		}
	}
}

func (e *EvaluatorRuntime) run(request evaluatorRequest) {
	if e.beforeRun != nil {
		e.beforeRun()
	}
	value, environment, err := dispatchEvaluate(request.context, request.expression, request.environment)
	request.stopCaller()
	request.cancel()
	if e.isClosed() && isCancellation(err) && request.caller.Err() == nil {
		err = ErrEvaluatorClosed
	}
	deliverTerminal(request.result, EvaluationResult{Value: value, Environment: environment, Err: err})
	e.releaseSlot()
}

func isCancellation(err error) bool {
	return errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)
}

// deliverTerminal is used after a request has been admitted. Result channels
// are capacity one and owned by the request, so terminal delivery cannot block
// and must not be suppressed by cancellation of the work context.
func deliverTerminal(destination chan<- EvaluationResult, result EvaluationResult) {
	destination <- result
}

func (e *EvaluatorRuntime) isClosed() bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.closed
}

func (e *EvaluatorRuntime) releaseSlot() {
	select {
	case <-e.slots:
	default:
		panic("evaluator admission slot released without admission")
	}
}
