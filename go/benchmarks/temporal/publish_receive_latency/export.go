package publishlatency

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"
)

// Export writes one raw CSV and one JSON summary document. Files are written
// only after the workload has ended, so formatting and filesystem latency are
// outside all measured intervals.
func Export(result RunResult, outputDir string) (csvPath, jsonPath string, err error) {
	if outputDir == "" {
		return "", "", fmt.Errorf("output directory must not be empty")
	}
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return "", "", fmt.Errorf("create output directory: %w", err)
	}
	records := append([]Measurement(nil), result.Records...)
	sort.Slice(records, func(i, j int) bool { return records[i].OperationID < records[j].OperationID })
	csvPath = filepath.Join(outputDir, result.RunID+".csv")
	file, err := os.Create(csvPath)
	if err != nil {
		return "", "", fmt.Errorf("create CSV: %w", err)
	}
	writer := csv.NewWriter(file)
	header := []string{
		"run_id", "operation_id", "scenario", "mode", "warmup", "status", "error",
		"t0", "t1", "t2", "t3", "schedule_initiated",
		"publish_to_receive_ns", "publish_to_receive_us", "publish_to_receive_ms",
		"publish_api_latency_ns", "publish_api_latency_us", "publish_api_latency_ms",
		"post_acknowledgment_delivery_ns", "post_acknowledgment_delivery_us", "post_acknowledgment_delivery_ms",
		"total_execution_ns", "total_execution_us", "total_execution_ms",
		"schedule_to_receive_ns", "schedule_to_receive_us", "schedule_to_receive_ms",
	}
	if err := writer.Write(header); err != nil {
		_ = file.Close()
		return "", "", fmt.Errorf("write CSV header: %w", err)
	}
	for _, record := range records {
		row := []string{
			record.RunID,
			record.OperationID,
			record.Scenario,
			record.Mode,
			strconv.FormatBool(record.Warmup),
			record.Status,
			record.Error,
			timestampExport(record.T0), timestampExport(record.T1), timestampExport(record.T2), timestampExport(record.T3), timestampExport(record.ScheduleInitiated),
			durationNS(record.PublishToReceive), durationUS(record.PublishToReceive), durationMS(record.PublishToReceive),
			durationNS(record.PublishAPILatency), durationUS(record.PublishAPILatency), durationMS(record.PublishAPILatency),
			durationNS(record.PostAcknowledgmentDelivery), durationUS(record.PostAcknowledgmentDelivery), durationMS(record.PostAcknowledgmentDelivery),
			durationNS(record.TotalExecution), durationUS(record.TotalExecution), durationMS(record.TotalExecution),
			durationNS(record.ScheduleToReceive), durationUS(record.ScheduleToReceive), durationMS(record.ScheduleToReceive),
		}
		if err := writer.Write(row); err != nil {
			_ = file.Close()
			return "", "", fmt.Errorf("write CSV row: %w", err)
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		_ = file.Close()
		return "", "", fmt.Errorf("flush CSV: %w", err)
	}
	if err := file.Close(); err != nil {
		return "", "", fmt.Errorf("close CSV: %w", err)
	}

	document := struct {
		RunID       string            `json:"run_id"`
		GeneratedAt string            `json:"generated_at"`
		Environment Environment       `json:"environment"`
		Summaries   []ScenarioSummary `json:"summaries"`
		RawCSV      string            `json:"raw_csv"`
		RecordCount int               `json:"record_count"`
	}{
		RunID:       result.RunID,
		GeneratedAt: time.Now().UTC().Format(time.RFC3339Nano),
		Environment: result.Environment,
		Summaries:   result.Summaries,
		RawCSV:      csvPath,
		RecordCount: len(result.Records),
	}
	jsonData, err := json.MarshalIndent(document, "", "  ")
	if err != nil {
		return "", "", fmt.Errorf("marshal JSON summary: %w", err)
	}
	jsonPath = filepath.Join(outputDir, result.RunID+".json")
	if err := os.WriteFile(jsonPath, append(jsonData, '\n'), 0o644); err != nil {
		return "", "", fmt.Errorf("write JSON summary: %w", err)
	}
	return csvPath, jsonPath, nil
}

func timestampExport(value time.Time) string {
	if value.IsZero() {
		return ""
	}
	return value.UTC().Format(time.RFC3339Nano)
}

func durationNS(value time.Duration) string { return strconv.FormatInt(value.Nanoseconds(), 10) }
func durationUS(value time.Duration) string {
	return strconv.FormatFloat(float64(value)/float64(time.Microsecond), 'f', 3, 64)
}
func durationMS(value time.Duration) string {
	return strconv.FormatFloat(float64(value)/float64(time.Millisecond), 'f', 6, 64)
}
