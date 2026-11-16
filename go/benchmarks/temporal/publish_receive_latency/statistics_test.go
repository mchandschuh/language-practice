package publishlatency

import (
	"math"
	"testing"
	"time"
)

func TestBuildSummariesCalculatesPercentilesAndCountsFailures(t *testing.T) {
	base := time.Now()
	records := make([]Measurement, 0, 6)
	for i := 1; i <= 4; i++ {
		t0 := base.Add(time.Duration(i) * time.Second)
		t2 := t0.Add(time.Duration(i*10) * time.Microsecond)
		records = append(records, Measurement{
			Operation:                  Operation{RunID: "run", OperationID: "op" + time.Duration(i).String(), Scenario: ScenarioStart, Mode: ModeSequential},
			Status:                     statusSuccess,
			T0:                         t0,
			T1:                         t0.Add(time.Microsecond),
			T2:                         t2,
			T3:                         t2.Add(time.Microsecond),
			PublishToReceive:           t2.Sub(t0),
			PublishAPILatency:          time.Microsecond,
			PostAcknowledgmentDelivery: 9 * time.Microsecond,
			TotalExecution:             11 * time.Microsecond,
		})
	}
	records = append(records,
		Measurement{Operation: Operation{Scenario: ScenarioStart, Mode: ModeSequential}, Status: statusFailure},
		Measurement{Operation: Operation{Scenario: ScenarioStart, Mode: ModeSequential}, Status: statusTimeout},
		Measurement{Operation: Operation{Scenario: ScenarioStart, Mode: ModeSequential, Warmup: true}, Status: statusSuccess, PublishToReceive: time.Microsecond},
	)

	summaries := BuildSummaries(records)
	if len(summaries) != 2 {
		t.Fatalf("summary count = %d, want 2", len(summaries))
	}
	var measured ScenarioSummary
	for _, summary := range summaries {
		if !summary.Warmup {
			measured = summary
		}
	}
	if measured.SuccessfulIterations != 4 || measured.FailedIterations != 1 || measured.TimedOutIterations != 1 {
		t.Fatalf("counts = success %d failure %d timeout %d", measured.SuccessfulIterations, measured.FailedIterations, measured.TimedOutIterations)
	}
	if got, want := measured.PublishToReceive.MinimumUS, 10.0; got != want {
		t.Fatalf("minimum = %v, want %v", got, want)
	}
	if got, want := measured.PublishToReceive.MaximumUS, 40.0; got != want {
		t.Fatalf("maximum = %v, want %v", got, want)
	}
	if got, want := measured.PublishToReceive.MedianUS, 25.0; got != want {
		t.Fatalf("median = %v, want %v", got, want)
	}
	if got, want := measured.PublishToReceive.MeanUS, 25.0; got != want {
		t.Fatalf("mean = %v, want %v", got, want)
	}
	if got, want := measured.PublishToReceive.P90US, 37.0; got != want {
		t.Fatalf("p90 = %v, want %v", got, want)
	}
	if got, want := measured.PublishToReceive.P999US, 39.97; got != want {
		t.Fatalf("p99.9 = %v, want %v", got, want)
	}
	if got, want := measured.PublishToReceive.StandardDeviationUS, math.Sqrt(125); math.Abs(got-want) > 1e-9 {
		t.Fatalf("standard deviation = %v, want %v", got, want)
	}
}
