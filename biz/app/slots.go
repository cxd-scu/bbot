package app

import (
	"context"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/schema"
	"github.com/cxd-scu/bbot/biz/agent/graph"
	"github.com/cxd-scu/bbot/biz/constant"
	"github.com/cxd-scu/bbot/biz/infra/cache"
	"github.com/cxd-scu/bbot/biz/infra/repo"
	"github.com/cxd-scu/bbot/biz/model"
	"github.com/cxd-scu/bbot/biz/utils"
)

func CallSlotAgent(ctx context.Context, agentID int64) (*model.IOParam, error) {
	agentInstance, err := repo.GetAgentByID(ctx, agentID)
	if err != nil {
		return nil, err
	}

	sessionID := agentInstance.SessionID
	session := cache.GetOrCreateSession(sessionID)

	var ret *adk.AsyncIterator[*adk.TypedAgentEvent[*schema.Message]]

	if !session.Pending {
		ret = graph.SlotAgentRunner.Run(ctx, []*schema.Message{
			{
				Role:    schema.User,
				Content: agentInstance.Input,
			},
		},
			adk.WithCheckPointID(session.CheckPointID),
			adk.WithSessionValues(map[string]any{constant.SessionIDKey: sessionID}),
		)
	} else {
		ret, err = graph.SlotAgentRunner.ResumeWithParams(ctx, session.CheckPointID, &adk.ResumeParams{
			Targets: map[string]any{
				"input": agentInstance.Input,
			},
		}, adk.WithSessionValues(map[string]any{constant.SessionIDKey: sessionID}))
		if err != nil {
			return nil, err
		}
	}

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

	// fallback 触发的澄清卡优先返回
	if c := utils.ExtractClarifyCard(final); c != nil {
		return &model.IOParam{Params: map[string]interface{}{"card": c}}, nil
	}
	return utils.ParseJSONToIOParam(ctx, final)
}
