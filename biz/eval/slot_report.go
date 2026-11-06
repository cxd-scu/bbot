package eval

import (
	"fmt"
	"strings"
)

// FormatSlotReport 格式化 slot 评测报告为可读文本。
func FormatSlotReport(report *SlotReport) string {
	var sb strings.Builder
	sb.WriteString("Slot Agent 端到端评测结果\n")
	sb.WriteString(fmt.Sprintf("总数: %d  正确: %d  准确率: %.2f%%\n", report.Total, report.Correct, report.Accuracy()*100))

	// 明细
	sb.WriteString("\n用例明细:\n")
	for _, r := range report.Results {
		status := "✓"
		if !r.Correct {
			status = "✗"
		}
		sb.WriteString(fmt.Sprintf("  %s [%s] %q", status, r.Case.Name, r.Case.Query))
		if len(r.ToolsCalled) > 0 {
			sb.WriteString(fmt.Sprintf(" 调用=%v", r.ToolsCalled))
		}
		if r.Case.WantClarify {
			sb.WriteString(fmt.Sprintf(" 澄清=%v", r.Clarified))
		}
		if r.Case.WantSetting != nil {
			if r.Generated != nil {
				sb.WriteString(fmt.Sprintf(" 生成=%v", r.Generated))
			} else {
				sb.WriteString(" 生成=<nil>")
			}
		}
		if len(r.Diffs) > 0 {
			sb.WriteString(fmt.Sprintf(" 差异=%v", r.Diffs))
		}
		sb.WriteString("\n")
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
