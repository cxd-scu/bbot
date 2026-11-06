package eval

import (
	"fmt"
	"sort"
	"strings"
)

// GatewayCaseResult 单条用例的运行结果。
type GatewayCaseResult struct {
	Case       GatewayCase
	GotAgent   string
	Confidence float64
	Reason     string
	Correct    bool
}

// GatewayCaseError 运行失败的用例（LLM 调用报错）。
type GatewayCaseError struct {
	Case GatewayCase
	Err  error
}

// GatewayReport gateway 评测报告。
type GatewayReport struct {
	Total    int
	Correct  int
	Accuracy float64
	Details  []GatewayCaseResult
	Errors   []GatewayCaseError
}

// ComputeConfusion 由运行结果构建混淆矩阵：want -> got -> 次数。
func ComputeConfusion(results []GatewayCaseResult) map[string]map[string]int {
	confusion := map[string]map[string]int{}
	for _, r := range results {
		if confusion[r.Case.WantAgent] == nil {
			confusion[r.Case.WantAgent] = map[string]int{}
		}
		confusion[r.Case.WantAgent][r.GotAgent]++
	}
	return confusion
}

// FormatGatewayReport 格式化评测报告为可读文本。
func FormatGatewayReport(report *GatewayReport) string {
	var sb strings.Builder
	sb.WriteString("Gateway 路由评测结果\n")
	sb.WriteString(fmt.Sprintf("总数: %d  正确: %d  准确率: %.2f%%\n", report.Total, report.Correct, report.Accuracy*100))

	// 混淆矩阵
	sb.WriteString("\n混淆矩阵（期望 → 实际）:\n")
	confusion := ComputeConfusion(report.Details)
	wants := make([]string, 0, len(confusion))
	for k := range confusion {
		wants = append(wants, k)
	}
	sort.Strings(wants)
	for _, want := range wants {
		gots := make([]string, 0, len(confusion[want]))
		for k := range confusion[want] {
			gots = append(gots, k)
		}
		sort.Strings(gots)
		parts := make([]string, 0, len(gots))
		for _, got := range gots {
			parts = append(parts, fmt.Sprintf("%s:%d", got, confusion[want][got]))
		}
		sb.WriteString(fmt.Sprintf("  %s → %s\n", want, strings.Join(parts, ", ")))
	}

	// 错误用例
	var wrong []GatewayCaseResult
	for _, r := range report.Details {
		if !r.Correct {
			wrong = append(wrong, r)
		}
	}
	if len(wrong) > 0 {
		sb.WriteString("\n错误用例:\n")
		for _, r := range wrong {
			sb.WriteString(fmt.Sprintf("  [%s] %q 期望=%s 实际=%s (confidence=%.2f, reason=%q)\n",
				r.Case.Name, r.Case.Query, r.Case.WantAgent, r.GotAgent, r.Confidence, r.Reason))
		}
	}

	// 运行失败
	if len(report.Errors) > 0 {
		sb.WriteString("\n运行失败:\n")
		for _, e := range report.Errors {
			sb.WriteString(fmt.Sprintf("  [%s] %q err=%v\n", e.Case.Name, e.Case.Query, e.Err))
		}
	}

	return sb.String()
}
