package handler

import (
	"context"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/cxd-scu/bbot/biz/app"
	"github.com/cxd-scu/bbot/biz/constant"
	"github.com/cxd-scu/bbot/biz/infra/card"
	"github.com/cxd-scu/bbot/biz/infra/guard"
	"github.com/cxd-scu/bbot/biz/infra/repo"
	"github.com/cxd-scu/bbot/biz/model"
	"go.uber.org/zap"
)

func HandleGatewayReq(ctx context.Context, req *model.UserReq) (*model.IOParam, error) {
	// 0. 输入护栏：检测 prompt 注入，命中即拒绝
	if r := guard.Check(req.Query); r.Suspicious {
		return &model.IOParam{
			Params: map[string]interface{}{
				"answer": "检测到可能的提示词注入，已拒绝处理。请重新描述你的需求。",
			},
		}, nil
	}

	// 1. 预处理输入（含图片识别 → 优化 query）
	input, err := app.NormalizeQuery(ctx, req)
	if err != nil {
		return nil, err
	}

	// 2. 意图路由（未指定 agent 时）
	agentName := req.AgentName
	if agentName == "" {
		targetAgent, err := app.GatewayNlu(ctx, input)
		if err != nil {
			return nil, err
		}
		if targetAgent == constant.NeedClassify {
			// 无法归类：返回澄清卡让用户选择意图
			return &model.IOParam{
				Params: map[string]interface{}{"card": card.ClarifyCard(ctx)},
			}, nil
		}
		agentName = targetAgent
		req.AgentName = agentName
	}

	// 3. 获取/创建实例（主键 uid + agentName）
	agent, err := repo.GetAgentByUIDAndName(ctx, req.UID, agentName)
	if err != nil {
		logger.Error("🚨 获取 Agent 失败", zap.Error(err))
		return nil, err
	}
	if agent == nil {
		agent, err = repo.CreateAgent(ctx, req.UID, agentName, req.Channel)
		if err != nil {
			logger.Error("🚨 创建 Agent 失败", zap.Error(err))
			return nil, err
		}
	}

	// 4. 落地本轮输入
	agent.Input = input
	if len(req.Images) > 0 {
		agent.InputType = int(constant.InputTypeCard)
	} else {
		agent.InputType = int(constant.InputTypeText)
	}

	// 5. 分发到对应 L1 agent
	return ConductAgent(ctx, agent.AgentID)
}
