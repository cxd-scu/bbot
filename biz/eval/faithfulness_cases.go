package eval

import "github.com/cxd-scu/bbot/biz/model"

// DefaultFaithfulnessCases 默认的 RAG 忠实度评测集。
// 用于校准 LLM 裁判（judge）能否正确区分「忠实」与「编造」的回答；
// 后续接入真实 agent 回答后，同样用它评测 agent 的忠实度。
var DefaultFaithfulnessCases = []FaithfulnessCase{
	{
		Name:         "忠实回答",
		Query:        "值班时间是什么",
		Docs:         []*model.Document{{Title: "值班安排", Content: "值班时间是每天 9:00 到 18:00"}},
		Answer:       "值班时间是每天 9:00 到 18:00。",
		WantFaithful: true,
	},
	{
		Name:         "编造时间与额外信息",
		Query:        "值班时间是什么",
		Docs:         []*model.Document{{Title: "值班安排", Content: "值班时间是每天 9:00 到 18:00"}},
		Answer:       "值班时间是每天 8:00 到 20:00，而且值班期间可以远程办公。",
		WantFaithful: false,
	},
	{
		Name:         "无依据回答",
		Query:        "故障上报流程",
		Docs:         []*model.Document{{Title: "值班安排", Content: "值班时间是每天 9:00 到 18:00"}},
		Answer:       "故障上报需要先联系运维总监，再填写纸质表格。",
		WantFaithful: false,
	},
	{
		Name:         "部分编造",
		Query:        "告警怎么处理",
		Docs:         []*model.Document{{Title: "告警处理", Content: "收到告警后，先确认告警级别，再通知值班人员"}},
		Answer:       "收到告警后，先确认告警级别，再通知值班人员，并且必须在 5 分钟内响应。",
		WantFaithful: false,
	},
	{
		Name:         "忠实但简要转述",
		Query:        "告警怎么处理",
		Docs:         []*model.Document{{Title: "告警处理", Content: "收到告警后，先确认告警级别，再通知值班人员"}},
		Answer:       "收到告警要先确认级别，然后通知值班人员。",
		WantFaithful: true,
	},
}
