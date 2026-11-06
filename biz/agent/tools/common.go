package tools

import (
	"context"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cxd-scu/bbot/biz/infra/card"
	"github.com/cxd-scu/bbot/biz/model"
)

func InitCommonTools() error {
	var err error
	FallbackTool, err = utils.InferTool(
		"fallback",
		"兜底工具：当用户请求与资源位配置无关、无法归类到任何可处理任务，或需要引导用户转人工时调用。"+
			"输入为无法处理的原因（reason），返回一段引导用户回到可用能力的回复话术。",
		fallbackTool,
	)
	if err != nil {
		return err
	}
	return nil
}

var (
	FallbackTool tool.InvokableTool
)

func fallbackTool(ctx context.Context, req model.FallbackReq) (model.FallbackResp, error) {
	return model.FallbackResp{
		Reply: "抱歉，该请求暂不在我的处理范围内，请从澄清卡中选择你的意图。",
		Card:  card.ClarifyCard(ctx),
	}, nil
}
