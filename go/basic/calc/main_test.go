package main

import (
	"bytes"
	"context"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestRunEvaluatesExpression(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		wantOutput string
	}{
		{name: "addition", args: []string{"1", "+", "2"}, wantOutput: "3\n"},
		{name: "precedence", args: []string{"1", "+", "2", "*", "3"}, wantOutput: "7\n"},
		{name: "parentheses", args: []string{"(", "1", "+", "2", ")", "*", "3"}, wantOutput: "9\n"},
		{name: "prefix increment literal", args: []string{"++", "2"}, wantOutput: "3\n"},
		{name: "postfix increment literal", args: []string{"2++"}, wantOutput: "2\n"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			if status := run(context.Background(), test.args, &stdout, &stderr); status != 0 {
				t.Fatalf("run() status = %d, stderr = %q", status, stderr.String())
			}
			if got := stdout.String(); got != test.wantOutput {
				t.Errorf("stdout = %q, want %q", got, test.wantOutput)
			}
			if stderr.Len() != 0 {
				t.Errorf("stderr = %q, want empty", stderr.String())
			}
		})
	}
}

func TestRunPreservesArgumentBoundaries(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if status := run(context.Background(), []string{"1 +", "2"}, &stdout, &stderr); status != 0 {
		t.Fatalf("run() status = %d, stderr = %q", status, stderr.String())
	}
	if got, want := stdout.String(), "3\n"; got != want {
		t.Errorf("stdout = %q, want %q", got, want)
	}
}

func TestRunReportsUsageForMissingExpression(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if status := run(context.Background(), nil, &stdout, &stderr); status != 2 {
		t.Fatalf("run() status = %d, want 2", status)
	}
	if got, want := stderr.String(), "usage: calc <expression>\n"; got != want {
		t.Errorf("stderr = %q, want %q", got, want)
	}
	if stdout.Len() != 0 {
		t.Errorf("stdout = %q, want empty", stdout.String())
	}
}

func TestRunReportsParseAndEvaluationErrors(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{name: "parse error", args: []string{"1", "+"}, want: "expected a number"},
		{name: "undefined identifier", args: []string{"missing"}, want: `identifier "missing" is not defined`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			if status := run(context.Background(), test.args, &stdout, &stderr); status == 0 {
				t.Fatal("run() returned success for an invalid expression")
			}
			if !strings.Contains(stderr.String(), test.want) {
				t.Errorf("stderr = %q, want substring %q", stderr.String(), test.want)
			}
			if stdout.Len() != 0 {
				t.Errorf("stdout = %q, want empty", stdout.String())
			}
		})
	}
}

func TestRunReportsCanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	var stdout, stderr bytes.Buffer
	if status := run(ctx, []string{"1", "+", "2"}, &stdout, &stderr); status == 0 {
		t.Fatal("run() returned success for a canceled context")
	}
	if !strings.Contains(stderr.String(), context.Canceled.Error()) {
		t.Errorf("stderr = %q, want cancellation error", stderr.String())
	}
}

func TestInitHelperRunsCalcAndCalcTests(t *testing.T) {
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	repoRoot := filepath.Clean(filepath.Join(filepath.Dir(currentFile), "../../.."))

	tests := []struct {
		name string
		args []string
	}{
		{name: "run calc", args: []string{"run", "calc", "1", "+", "2"}},
		{name: "quoted run calc", args: []string{"run", "calc", "1 + 2"}},
		{name: "test calc", args: []string{"test", "calc", "-run", "TestRunEvaluatesExpression", "-count=1"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			commandArgs := append([]string{"-c", "source ./init; \"$@\"", "init-test"}, test.args...)
			command := exec.Command("bash", commandArgs...)
			command.Dir = repoRoot
			output, err := command.CombinedOutput()
			if err != nil {
				t.Fatalf("helper failed: %v\noutput:\n%s", err, output)
			}
			if test.name != "test calc" && !strings.Contains(string(output), "3") {
				t.Errorf("output = %q, want result 3", output)
			}
		})
	}
}
