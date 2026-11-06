package eval

import (
	"context"
	"fmt"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/eino/schema"
	"github.com/cxd-scu/bbot/biz/agent/graph"
	"github.com/cxd-scu/bbot/biz/model"
)

// SlotCaseResult 单条 slot 用例的运行结果。
type SlotCaseResult struct {
	Case SlotEvalCase

	Clarified   bool               // 是否触发澄清（need_clarify / interrupt）
	ToolsCalled []string           // 模型调用的工具名序列
	Generated   *model.SlotSetting // generate_setting 收到的配置（模型生成的结果）

	FinalMessage string   // 最终回复文本
	Diffs        []string // 与期望配置的字段差异
	Correct      bool
}

// SlotCaseError 运行失败的用例。
type SlotCaseError struct {
	Case SlotEvalCase
	Err  error
}

// SlotReport slot agent 评测报告。
type SlotReport struct {
	Total   int
	Correct int
	Results []SlotCaseResult
	Errors  []SlotCaseError
}

// Accuracy 返回准确率。
func (r *SlotReport) Accuracy() float64 {
	if r.Total == 0 {
		return 0
	}
	return float64(r.Correct) / float64(r.Total)
}

// EvalSlotAgent 运行 slot agent 端到端评测集。
// 直接跑 SlotAgentRunner（不接会话/仓库），逐条观察工具调用与最终产出。
func EvalSlotAgent(ctx context.Context, cases []SlotEvalCase) (*SlotReport, error) {
	report := &SlotReport{Total: len(cases)}

	for _, c := range cases {
		res, err := runSlotCase(ctx, c)
		if err != nil {
			report.Errors = append(report.Errors, SlotCaseError{Case: c, Err: err})
			continue
		}

		judge(&res)
		report.Results = append(report.Results, res)
		if res.Correct {
			report.Correct++
		}
	}

	return report, nil
}

// runSlotCase 跑单条用例，观察 agent 的工具调用与产出。
func runSlotCase(ctx context.Context, c SlotEvalCase) (SlotCaseResult, error) {
	res := SlotCaseResult{Case: c}

	iter := graph.SlotAgentRunner.Run(ctx, []*schema.Message{
		{Role: schema.User, Content: c.Query},
	})

	for {
		ev, ok := iter.Next()
		if !ok {
			break
		}
		if ev.Err != nil {
			return res, ev.Err
		}
		if ev.Action != nil && ev.Action.Interrupted != nil {
			res.Clarified = true
		}
		if ev.Output == nil || ev.Output.MessageOutput == nil || ev.Output.MessageOutput.Message == nil {
			continue
		}

		m := ev.Output.MessageOutput.Message
		res.FinalMessage = m.Content
		for _, tc := range m.ToolCalls {
			name := tc.Function.Name
			res.ToolsCalled = append(res.ToolsCalled, name)
			switch name {
			case "need_clarify":
				res.Clarified = true
			case "generate_setting":
				if res.Generated == nil {
					var s model.SlotSetting
					if err := sonic.UnmarshalString(tc.Function.Arguments, &s); err == nil {
						res.Generated = &s
					}
				}
			}
		}
	}

	return res, nil
}

// judge 按用例期望判定对错，并填充差异信息。
func judge(res *SlotCaseResult) {
	switch {
	case res.Case.WantClarify:
		res.Correct = res.Clarified
	case res.Case.WantSetting != nil:
		if res.Generated == nil {
			res.Correct = false
			res.Diffs = []string{"未调用 generate_setting 或未生成配置"}
		} else {
			res.Diffs = compareSetting(res.Generated, res.Case.WantSetting)
			res.Correct = len(res.Diffs) == 0
		}
	default:
		res.Correct = true
	}
}

// compareSetting 字段级对比两个配置，仅比对 want 中明确设置的字段，返回差异列表。
func compareSetting(got, want *model.SlotSetting) []string {
	var diffs []string
	if want.Page != "" && got.Page != want.Page {
		diffs = append(diffs, fmt.Sprintf("page: got=%q want=%q", got.Page, want.Page))
	}
	if want.ResourceType != "" && got.ResourceType != want.ResourceType {
		diffs = append(diffs, fmt.Sprintf("resource_type: got=%q want=%q", got.ResourceType, want.ResourceType))
	}
	if len(want.CityIDs) > 0 && !equalInt64s(got.CityIDs, want.CityIDs) {
		diffs = append(diffs, fmt.Sprintf("city_ids: got=%v want=%v", got.CityIDs, want.CityIDs))
	}
	if len(want.EffectSource) > 0 && !equalStrings(got.EffectSource, want.EffectSource) {
		diffs = append(diffs, fmt.Sprintf("effect_source: got=%v want=%v", got.EffectSource, want.EffectSource))
	}
	if want.ImageURL != "" && got.ImageURL != want.ImageURL {
		diffs = append(diffs, fmt.Sprintf("image_url: got=%q want=%q", got.ImageURL, want.ImageURL))
	}
	if want.OpenURL != "" && got.OpenURL != want.OpenURL {
		diffs = append(diffs, fmt.Sprintf("open_url: got=%q want=%q", got.OpenURL, want.OpenURL))
	}
	if want.StartTime != "" && got.StartTime != want.StartTime {
		diffs = append(diffs, fmt.Sprintf("start_time: got=%q want=%q", got.StartTime, want.StartTime))
	}
	if want.EndTime != "" && got.EndTime != want.EndTime {
		diffs = append(diffs, fmt.Sprintf("end_time: got=%q want=%q", got.EndTime, want.EndTime))
	}
	if want.IsLimit != nil && !boolPtrEqual(got.IsLimit, want.IsLimit) {
		diffs = append(diffs, fmt.Sprintf("is_limit: got=%v want=%v", got.IsLimit, want.IsLimit))
	}
	return diffs
}

func equalInt64s(a, b []int64) bool {
	if len(a) != len(b) {
		return false
	}
	ma := map[int64]int{}
	for _, v := range a {
		ma[v]++
	}
	for _, v := range b {
		if ma[v] == 0 {
			return false
		}
		ma[v]--
	}
	return true
}

func equalStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	ma := map[string]int{}
	for _, v := range a {
		ma[v]++
	}
	for _, v := range b {
		if ma[v] == 0 {
			return false
		}
		ma[v]--
	}
	return true
}

func boolPtrEqual(a, b *bool) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}
