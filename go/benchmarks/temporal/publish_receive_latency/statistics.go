package publishlatency

import (
	"math"
	"sort"
	"time"
)

// DurationStats reports successful latency samples in microseconds and
// milliseconds. Percentiles use linear interpolation between ordered samples.
type DurationStats struct {
	Count               int     `json:"count"`
	MinimumUS           float64 `json:"minimum_us"`
	MeanUS              float64 `json:"mean_us"`
	MedianUS            float64 `json:"median_p50_us"`
	P90US               float64 `json:"p90_us"`
	P95US               float64 `json:"p95_us"`
	P99US               float64 `json:"p99_us"`
	P999US              float64 `json:"p99_9_us"`
	MaximumUS           float64 `json:"maximum_us"`
	StandardDeviationUS float64 `json:"standard_deviation_us"`

	MinimumMS           float64 `json:"minimum_ms"`
	MeanMS              float64 `json:"mean_ms"`
	MedianMS            float64 `json:"median_p50_ms"`
	P90MS               float64 `json:"p90_ms"`
	P95MS               float64 `json:"p95_ms"`
	P99MS               float64 `json:"p99_ms"`
	P999MS              float64 `json:"p99_9_ms"`
	MaximumMS           float64 `json:"maximum_ms"`
	StandardDeviationMS float64 `json:"standard_deviation_ms"`
}

// ScenarioSummary separates warm-up data from measured data and keeps failures
// visible instead of silently excluding them.
type ScenarioSummary struct {
	Scenario                   string        `json:"scenario"`
	Mode                       string        `json:"mode"`
	Warmup                     bool          `json:"warmup"`
	SuccessfulIterations       int           `json:"successful_iterations"`
	FailedIterations           int           `json:"failed_iterations"`
	TimedOutIterations         int           `json:"timed_out_iterations"`
	FailedOrTimedOutIterations int           `json:"failed_or_timed_out_iterations"`
	PublishToReceive           DurationStats `json:"publish_to_receive"`
	PublishAPILatency          DurationStats `json:"publish_api_latency"`
	PostAcknowledgmentDelivery DurationStats `json:"post_acknowledgment_delivery"`
	TotalExecution             DurationStats `json:"total_execution"`
	ScheduleToReceive          DurationStats `json:"schedule_to_receive"`
}

// BuildSummaries groups records by scenario, publication mode, and warm-up
// marker. Only successful records contribute latency samples; counts include
// every record.
func BuildSummaries(records []Measurement) []ScenarioSummary {
	type key struct {
		scenario string
		mode     string
		warmup   bool
	}
	groups := make(map[key][]Measurement)
	for _, record := range records {
		groupKey := key{scenario: record.Scenario, mode: record.Mode, warmup: record.Warmup}
		groups[groupKey] = append(groups[groupKey], record)
	}
	keys := make([]key, 0, len(groups))
	for groupKey := range groups {
		keys = append(keys, groupKey)
	}
	sort.Slice(keys, func(i, j int) bool {
		if keys[i].scenario != keys[j].scenario {
			return keys[i].scenario < keys[j].scenario
		}
		if keys[i].mode != keys[j].mode {
			return keys[i].mode < keys[j].mode
		}
		return !keys[i].warmup && keys[j].warmup
	})

	result := make([]ScenarioSummary, 0, len(keys))
	for _, groupKey := range keys {
		group := groups[groupKey]
		summary := ScenarioSummary{Scenario: groupKey.scenario, Mode: groupKey.mode, Warmup: groupKey.warmup}
		publishToReceive := make([]time.Duration, 0, len(group))
		publishAPI := make([]time.Duration, 0, len(group))
		postAck := make([]time.Duration, 0, len(group))
		total := make([]time.Duration, 0, len(group))
		scheduleToReceive := make([]time.Duration, 0, len(group))
		for _, record := range group {
			switch record.Status {
			case statusSuccess:
				summary.SuccessfulIterations++
				if !record.T0.IsZero() && !record.T2.IsZero() {
					publishToReceive = append(publishToReceive, record.PublishToReceive)
				}
				if !record.T0.IsZero() && !record.T1.IsZero() {
					publishAPI = append(publishAPI, record.PublishAPILatency)
				}
				if !record.T1.IsZero() && !record.T2.IsZero() {
					postAck = append(postAck, record.PostAcknowledgmentDelivery)
				}
				if !record.T0.IsZero() && !record.T3.IsZero() {
					total = append(total, record.TotalExecution)
				}
				if !record.ScheduleInitiated.IsZero() && !record.T2.IsZero() {
					scheduleToReceive = append(scheduleToReceive, record.ScheduleToReceive)
				}
			case statusTimeout:
				summary.TimedOutIterations++
			default:
				summary.FailedIterations++
			}
		}
		summary.FailedOrTimedOutIterations = summary.FailedIterations + summary.TimedOutIterations
		summary.PublishToReceive = calculateDurationStats(publishToReceive)
		summary.PublishAPILatency = calculateDurationStats(publishAPI)
		summary.PostAcknowledgmentDelivery = calculateDurationStats(postAck)
		summary.TotalExecution = calculateDurationStats(total)
		summary.ScheduleToReceive = calculateDurationStats(scheduleToReceive)
		result = append(result, summary)
	}
	return result
}

func calculateDurationStats(durations []time.Duration) DurationStats {
	if len(durations) == 0 {
		return DurationStats{}
	}
	values := make([]float64, len(durations))
	for i, duration := range durations {
		values[i] = float64(duration) / float64(time.Microsecond)
	}
	sort.Float64s(values)
	var sum float64
	for _, value := range values {
		sum += value
	}
	mean := sum / float64(len(values))
	var squaredError float64
	for _, value := range values {
		delta := value - mean
		squaredError += delta * delta
	}
	standardDeviation := math.Sqrt(squaredError / float64(len(values)))
	stats := DurationStats{
		Count:               len(values),
		MinimumUS:           values[0],
		MeanUS:              mean,
		MedianUS:            percentile(values, 0.50),
		P90US:               percentile(values, 0.90),
		P95US:               percentile(values, 0.95),
		P99US:               percentile(values, 0.99),
		P999US:              percentile(values, 0.999),
		MaximumUS:           values[len(values)-1],
		StandardDeviationUS: standardDeviation,
	}
	stats.MinimumMS = stats.MinimumUS / 1000
	stats.MeanMS = stats.MeanUS / 1000
	stats.MedianMS = stats.MedianUS / 1000
	stats.P90MS = stats.P90US / 1000
	stats.P95MS = stats.P95US / 1000
	stats.P99MS = stats.P99US / 1000
	stats.P999MS = stats.P999US / 1000
	stats.MaximumMS = stats.MaximumUS / 1000
	stats.StandardDeviationMS = stats.StandardDeviationUS / 1000
	return stats
}

func percentile(values []float64, probability float64) float64 {
	if len(values) == 0 {
		return 0
	}
	if probability <= 0 {
		return values[0]
	}
	if probability >= 1 {
		return values[len(values)-1]
	}
	position := probability * float64(len(values)-1)
	lower := int(position)
	upper := lower + 1
	if upper >= len(values) {
		return values[lower]
	}
	fraction := position - float64(lower)
	return values[lower] + fraction*(values[upper]-values[lower])
}
