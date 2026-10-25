package calc

import (
	"context"
	"strconv"
	"testing"
)

type benchmarkCase struct {
	expression MathExpression
	input      Environment
	wantValue  float64
}

var benchmarkCorpus = mustBenchmarkCorpus()

func mustBenchmarkCorpus() []benchmarkCase {
	definitions := []struct {
		expression string
		input      Environment
	}{
		{expression: "x * x + 2", input: Environment{"x": 3}},
		{expression: "(9 - 4) / 5 + 2 * 3", input: nil},
		{expression: "-x + 3", input: Environment{"x": 4}},
		{expression: "x++ + ++x", input: Environment{"x": 5}},
	}

	corpus := make([]benchmarkCase, 0, len(definitions))
	for _, definition := range definitions {
		expression, err := Parse(definition.expression)
		if err != nil {
			panic(err)
		}
		value, _, err := (Evaluator{}).EvaluateWithEnvironment(expression, definition.input)
		if err != nil {
			panic(err)
		}
		corpus = append(corpus, benchmarkCase{
			expression: expression,
			input:      definition.input,
			wantValue:  value,
		})
	}
	return corpus
}

func BenchmarkSequentialOracle(b *testing.B) {
	evaluator := Evaluator{}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		work := benchmarkCorpus[i%len(benchmarkCorpus)]
		value, _, err := evaluator.EvaluateWithEnvironment(work.expression, work.input)
		if err != nil {
			b.Fatal(err)
		}
		if value != work.wantValue {
			b.Fatalf("value = %v, want %v", value, work.wantValue)
		}
	}
}

func BenchmarkChannelEvaluator(b *testing.B) {
	for _, capacity := range []int{1, DefaultQueueCapacity} {
		b.Run("capacity-"+benchmarkIntString(capacity), func(b *testing.B) {
			evaluator, err := NewEvaluator(context.Background(), capacity)
			if err != nil {
				b.Fatal(err)
			}
			b.Cleanup(func() { _ = evaluator.Close() })
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				work := benchmarkCorpus[i%len(benchmarkCorpus)]
				results, err := evaluator.Submit(context.Background(), work.expression, work.input)
				if err != nil {
					b.Fatal(err)
				}
				result := <-results
				if result.Err != nil {
					b.Fatal(result.Err)
				}
				if result.Value != work.wantValue {
					b.Fatalf("value = %v, want %v", result.Value, work.wantValue)
				}
			}
		})
	}
}

func BenchmarkChannelBatch(b *testing.B) {
	const batchSize = 16
	for _, capacity := range []int{1, DefaultQueueCapacity} {
		b.Run("capacity-"+benchmarkIntString(capacity), func(b *testing.B) {
			evaluator, err := NewEvaluator(context.Background(), capacity)
			if err != nil {
				b.Fatal(err)
			}
			b.Cleanup(func() { _ = evaluator.Close() })
			b.ReportAllocs()
			b.ResetTimer()
			for iteration := 0; iteration < b.N; iteration++ {
				results := make([]<-chan EvaluationResult, batchSize)
				for i := range results {
					work := benchmarkCorpus[(iteration*batchSize+i)%len(benchmarkCorpus)]
					results[i], err = evaluator.Submit(context.Background(), work.expression, work.input)
					if err != nil {
						b.Fatal(err)
					}
				}
				for i, resultChannel := range results {
					work := benchmarkCorpus[(iteration*batchSize+i)%len(benchmarkCorpus)]
					result := <-resultChannel
					if result.Err != nil {
						b.Fatal(result.Err)
					}
					if result.Value != work.wantValue {
						b.Fatalf("value = %v, want %v", result.Value, work.wantValue)
					}
				}
			}
		})
	}
}

func benchmarkIntString(value int) string {
	return strconv.Itoa(value)
}
