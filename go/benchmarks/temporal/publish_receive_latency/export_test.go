package publishlatency

import (
	"encoding/csv"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestExportWritesRawCSVAndJSONSummary(t *testing.T) {
	directory := t.TempDir()
	result := RunResult{
		RunID: "run-test",
		Records: []Measurement{{
			Operation:        Operation{RunID: "run-test", OperationID: "operation", Scenario: ScenarioStart, Mode: ModeSequential},
			Status:           statusSuccess,
			T0:               time.Now(),
			T1:               time.Now().Add(time.Microsecond),
			T2:               time.Now().Add(2 * time.Microsecond),
			PublishToReceive: 2 * time.Microsecond,
		}},
	}
	result.Summaries = BuildSummaries(result.Records)
	csvPath, jsonPath, err := Export(result, filepath.Join(directory, "nested"))
	if err != nil {
		t.Fatal(err)
	}
	csvFile, err := os.Open(csvPath)
	if err != nil {
		t.Fatal(err)
	}
	rows, err := csv.NewReader(csvFile).ReadAll()
	_ = csvFile.Close()
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 2 || len(rows[0]) == 0 {
		t.Fatalf("CSV rows = %d, want header plus one record", len(rows))
	}
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		t.Fatal(err)
	}
	var document struct {
		RunID       string `json:"run_id"`
		RecordCount int    `json:"record_count"`
	}
	if err := json.Unmarshal(data, &document); err != nil {
		t.Fatal(err)
	}
	if document.RunID != result.RunID || document.RecordCount != 1 {
		t.Fatalf("JSON document = %+v", document)
	}
}
