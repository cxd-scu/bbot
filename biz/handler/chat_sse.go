package handler

import (
	"context"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/sse"
	bizapp "github.com/cxd-scu/bbot/biz/app"
	"github.com/cxd-scu/bbot/biz/constant"
	"github.com/cxd-scu/bbot/biz/infra/card"
	"github.com/cxd-scu/bbot/biz/infra/guard"
	"github.com/cxd-scu/bbot/biz/model"
)

// HandleChatSSE 流式对话入口：预处理 + 路由后，以 SSE 逐 token 返回 agent 输出。
func HandleChatSSE(ctx context.Context, c *app.RequestContext) {
	var req model.UserReq
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(400, map[string]any{"error": err.Error()})
		return
	}

	w := sse.NewWriter(c)
	defer w.Close()

	emit := func(delta string) error {
		return w.WriteEvent("", "message", []byte(delta))
	}

	// 0. 输入护栏：检测 prompt 注入，命中即拒绝
	if r := guard.Check(req.Query); r.Suspicious {
		_ = emit("检测到可能的提示词注入，已拒绝处理。请重新描述你的需求。")
		return
	}

	// 1. 预处理（含图片识别）
	input, err := bizapp.NormalizeQuery(ctx, &req)
	if err != nil {
		_ = emit("处理失败: " + err.Error())
		return
	}

	// 2. 意图路由
	agentName := req.AgentName
	if agentName == "" {
		targetAgent, err := bizapp.GatewayNlu(ctx, input)
		if err != nil {
			_ = emit("路由失败: " + err.Error())
			return
		}
		if targetAgent == constant.NeedClassify {
			cardJSON, _ := sonic.MarshalString(card.ClarifyCard(ctx))
			_ = w.WriteEvent("", "card", []byte(cardJSON))
			return
		}
		agentName = targetAgent
	}

	// 3. 流式执行对应 L1 agent
	switch agentName {
	case constant.SlotAgent:
		_ = bizapp.StreamSlotAgent(ctx, input, emit)
	case constant.OncallAgent:
		_ = bizapp.StreamOncallAgent(ctx, input, emit)
	default:
		_ = emit("不支持的 agent: " + agentName)
	}
}
