package app

import (
	"context"

	"github.com/cxd-scu/bbot/biz/agent/graph"
	"github.com/cxd-scu/bbot/biz/model"
)

// GatewayNlu 路由意图，只返回目标 agent 名称。
func GatewayNlu(ctx context.Context, query string) (string, error) {
	out, err := GatewayNluFull(ctx, query)
	if err != nil {
		return "", err
	}
	return out.TargetAgent, nil
}

// GatewayNluFull 返回完整的路由结果（target_agent / reason / confidence），供评测与调试使用。
func GatewayNluFull(ctx context.Context, query string) (*model.GatewayAgentOutput, error) {
	ret, err := graph.GatewayRunner.Invoke(ctx, map[string]interface{}{
		"query": query,
	})
	if err != nil {
		return nil, err
	}

	out := &model.GatewayAgentOutput{}
	if v, ok := ret.Params["target_agent"].(string); ok {
		out.TargetAgent = v
	}
	if v, ok := ret.Params["reason"].(string); ok {
		out.Reason = v
	}
	if v, ok := ret.Params["confidence"].(float64); ok {
		out.Confidence = v
	}
	return out, nil
}
