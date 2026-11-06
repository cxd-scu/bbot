package eval

import (
	"strings"
	"testing"
)

func TestComputeConfusion(t *testing.T) {
	results := []GatewayCaseResult{
		{Case: GatewayCase{Name: "a", WantAgent: "slot"}, GotAgent: "slot", Correct: true},
		{Case: GatewayCase{Name: "b", WantAgent: "slot"}, GotAgent: "oncall", Correct: false},
		{Case: GatewayCase{Name: "c", WantAgent: "oncall"}, GotAgent: "oncall", Correct: true},
	}

	confusion := ComputeConfusion(results)

	if confusion["slot"]["slot"] != 1 {
		t.Fatalf("slot->slot 应为 1，got=%v", confusion["slot"])
	}
	if confusion["slot"]["oncall"] != 1 {
		t.Fatalf("slot->oncall 应为 1，got=%v", confusion["slot"])
	}
	if confusion["oncall"]["oncall"] != 1 {
		t.Fatalf("oncall->oncall 应为 1，got=%v", confusion["oncall"])
	}
}

func TestFormatGatewayReport(t *testing.T) {
	report := &GatewayReport{
		Total:    2,
		Correct:  1,
		Accuracy: 0.5,
		Details: []GatewayCaseResult{
			{Case: GatewayCase{Name: "a", Query: "配banner", WantAgent: "slot"}, GotAgent: "slot", Correct: true},
			{Case: GatewayCase{Name: "b", Query: "今天天气", WantAgent: "need_classify"}, GotAgent: "oncall", Correct: false},
		},
	}

	s := FormatGatewayReport(report)
	for _, want := range []string{"准确率: 50.00%", "slot → slot:1", "need_classify → oncall:1", "[b]"} {
		if !strings.Contains(s, want) {
			t.Fatalf("报告应包含 %q，got:\n%s", want, s)
		}
	}
}
