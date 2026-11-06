package handler

import (
	"context"
	"errors"

	"github.com/cxd-scu/bbot/biz/app"
	"github.com/cxd-scu/bbot/biz/constant"
	"github.com/cxd-scu/bbot/biz/infra/repo"
	"github.com/cxd-scu/bbot/biz/model"
)

// ConductAgent Mock rpc call
func ConductAgent(ctx context.Context, agentID int64) (*model.IOParam, error) {
	agentInstance, err := repo.GetAgentByID(ctx, agentID)
	if err != nil {
		return nil, err
	}
	if agentInstance == nil {
		return nil, errors.New("agent not found")
	}

	var result *model.IOParam
	switch agentInstance.AgentName {
	case constant.SlotAgent:
		result, _ = app.CallSlotAgent(ctx, agentID)
	case constant.OncallAgent:
		result, _ = app.CallOncallAgent(ctx, agentID)
	default:
		return nil, errors.New("agent type not supported")
	}
	return result, nil
}
