package card

import (
	"context"

	"github.com/cxd-scu/bbot/biz/constant"
	"github.com/cxd-scu/bbot/biz/model"
)

// ClarifyCard 生成澄清卡（mock）：列出可选意图，让用户选择。
// 用于 gateway 返回 need_classify、或 L1 agent 调用 fallback 时返回给用户。
//
// TODO: 接入真实卡片渲染与交互，选项列表可由配置/服务发现动态生成。
func ClarifyCard(ctx context.Context) *model.ClarifyCard {
	return &model.ClarifyCard{
		Type:    "clarify_card",
		Title:   "请选择你的意图",
		Message: "我没有完全理解你的需求，请选择你想做的事情：",
		Options: []model.CardOption{
			{AgentName: constant.SlotAgent, Label: "资源位查询与配置"},
			{AgentName: constant.OncallAgent, Label: "oncall 信息查询"},
		},
	}
}
