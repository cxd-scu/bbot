package eval

import (
	"context"
	"fmt"
	"strings"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/eino/schema"
	"github.com/cxd-scu/bbot/biz/infra/llm"
	"github.com/cxd-scu/bbot/biz/model"
	"github.com/cxd-scu/bbot/biz/utils"
)

// FaithfulnessCase 一条 RAG 忠实度评测用例。
type FaithfulnessCase struct {
	Name         string
	Query        string
	Docs         []*model.Document // 检索到的文档
	Answer       string            // 待判断的回答
	WantFaithful bool              // 期望：忠实 / 不忠实
}

// FaithfulnessCaseResult 单条用例的判定结果。
type FaithfulnessCaseResult struct {
	Case        FaithfulnessCase
	GotFaithful bool
	Reason      string
	Correct     bool
}

// FaithfulnessCaseError 运行失败的用例。
type FaithfulnessCaseError struct {
	Case FaithfulnessCase
	Err  error
}

// FaithfulnessReport 忠实度评测报告。
type FaithfulnessReport struct {
	Total   int
	Correct int
	Results []FaithfulnessCaseResult
	Errors  []FaithfulnessCaseError
}

// Accuracy 返回裁判准确率。
func (r *FaithfulnessReport) Accuracy() float64 {
	if r.Total == 0 {
		return 0
	}
	return float64(r.Correct) / float64(r.Total)
}

// EvalFaithfulness 用 LLM-as-judge 逐条判断回答是否忠实于检索文档，并与期望标签对比。
func EvalFaithfulness(ctx context.Context, cases []FaithfulnessCase) (*FaithfulnessReport, error) {
	report := &FaithfulnessReport{Total: len(cases)}

	for _, c := range cases {
		faithful, reason, err := judgeFaithfulness(ctx, c.Query, c.Docs, c.Answer)
		if err != nil {
			report.Errors = append(report.Errors, FaithfulnessCaseError{Case: c, Err: err})
			continue
		}

		res := FaithfulnessCaseResult{
			Case:        c,
			GotFaithful: faithful,
			Reason:      reason,
			Correct:     faithful == c.WantFaithful,
		}
		report.Results = append(report.Results, res)
		if res.Correct {
			report.Correct++
		}
	}

	return report, nil
}

// judgeFaithfulness 用 LLM 裁判判断回答是否忠实于检索文档。
func judgeFaithfulness(ctx context.Context, query string, docs []*model.Document, answer string) (bool, string, error) {
	userContent := fmt.Sprintf("用户问题：%s\n\n检索到的文档：\n%s\n\n模型回答：%s", query, formatDocs(docs), answer)
	messages := []*schema.Message{
		schema.SystemMessage(faithfulnessJudgePrompt),
		schema.UserMessage(userContent),
	}

	resp, err := llm.DsPro().Generate(ctx, messages)
	if err != nil {
		return false, "", fmt.Errorf("faithfulness judge generate: %w", err)
	}

	clean := utils.ExtractJSON(resp.Content)
	var j faithfulnessJudgement
	if err := sonic.UnmarshalString(clean, &j); err != nil {
		return false, "", fmt.Errorf("faithfulness judge parse: %w, raw=%s", err, resp.Content)
	}
	return j.Faithful, j.Reason, nil
}

// faithfulnessJudgement 裁判输出。
type faithfulnessJudgement struct {
	Faithful bool   `json:"faithful"`
	Reason   string `json:"reason"`
}

func formatDocs(docs []*model.Document) string {
	var sb strings.Builder
	idx := 0
	for _, d := range docs {
		if d == nil || strings.TrimSpace(d.Content) == "" {
			continue
		}
		if idx > 0 {
			sb.WriteString("\n")
		}
		idx++
		sb.WriteString(fmt.Sprintf("[%d] %s\n%s", idx, d.Title, d.Content))
	}
	return sb.String()
}

const faithfulnessJudgePrompt = `你是一个 RAG 回答忠实度裁判。给定用户问题、检索到的文档、以及模型生成的回答，判断回答是否完全忠实于检索文档。

判定规则：
- 忠实（faithful=true）：回答中的所有事实性陈述都能在检索文档中找到依据，没有编造、没有添油加醋。
- 不忠实（faithful=false）：回答中出现了检索文档里没有的信息、数字、结论，或与文档矛盾。

只输出 JSON，不要输出其他内容：
{"faithful": true, "reason": "一句话说明判断依据"}`

// FormatFaithfulnessReport 格式化忠实度评测报告。
func FormatFaithfulnessReport(report *FaithfulnessReport) string {
	var sb strings.Builder
	sb.WriteString("RAG 忠实度评测结果（LLM-as-judge）\n")
	sb.WriteString(fmt.Sprintf("总数: %d  正确: %d  裁判准确率: %.2f%%\n", report.Total, report.Correct, report.Accuracy()*100))

	sb.WriteString("\n用例明细:\n")
	for _, r := range report.Results {
		status := "✓"
		if !r.Correct {
			status = "✗"
		}
		sb.WriteString(fmt.Sprintf("  %s [%s] 期望忠实=%v 裁判=%v reason=%q\n",
			status, r.Case.Name, r.Case.WantFaithful, r.GotFaithful, r.Reason))
	}

	if len(report.Errors) > 0 {
		sb.WriteString("\n运行失败:\n")
		for _, e := range report.Errors {
			sb.WriteString(fmt.Sprintf("  [%s] err=%v\n", e.Case.Name, e.Err))
		}
	}

	return sb.String()
}
