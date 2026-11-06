package app

import (
	"context"

	"github.com/cloudwego/eino/schema"
	"github.com/cxd-scu/bbot/biz/agent/graph"
	"github.com/cxd-scu/bbot/biz/infra/repo"
	"github.com/cxd-scu/bbot/biz/model"
	"github.com/cxd-scu/bbot/biz/utils"
)

// CallOncallAgent 处理 oncall 信息查询：检索知识库 → 生成回答。
func CallOncallAgent(ctx context.Context, agentID int64) (*model.IOParam, error) {
	agentInstance, err := repo.GetAgentByID(ctx, agentID)
	if err != nil {
		return nil, err
	}

	ret := graph.OncallRunner.Run(ctx, []*schema.Message{
		{Role: schema.User, Content: agentInstance.Input},
	})

	var final *schema.Message
	for {
		ev, ok := ret.Next()
		if !ok {
			break
		}
		if ev.Err != nil {
			return nil, ev.Err
		}
		if ev.Action != nil && ev.Action.Interrupted != nil {
			return &model.IOParam{}, nil
		}
		if ev.Output != nil && ev.Output.MessageOutput.Message != nil {
			final = ev.Output.MessageOutput.Message
		}
	}

	if final == nil {
		return &model.IOParam{}, nil
	}
	// fallback 触发的澄清卡优先返回
	if c := utils.ExtractClarifyCard(final); c != nil {
		return &model.IOParam{Params: map[string]interface{}{"card": c}}, nil
	}
	return &model.IOParam{
		Params: map[string]interface{}{
			"answer": final.Content,
		},
	}, nil
}
