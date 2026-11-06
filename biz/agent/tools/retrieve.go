package tools

import (
	"context"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cxd-scu/bbot/biz/infra/retriever"
	"github.com/cxd-scu/bbot/biz/model"
)

// NewRetrieveTool 创建知识库检索工具，把 Retriever 封装为 Agent 可调用的 tool。
func NewRetrieveTool(r retriever.Retriever) (tool.InvokableTool, error) {
	return utils.InferTool(
		"retrieve_knowledge",
		"从 oncall 知识库检索相关内容。当用户提问值班、流程、操作手册等 oncall 相关信息时，调用此工具获取知识库中的相关文档，再据此回答。"+
			"输入 query 为检索查询词（应提炼用户问题中的关键信息），top_k 为返回文档数量（缺省 5）。"+
			"返回检索到的文档片段列表（documents）。"+
			"只依据返回的文档回答，不得编造；检索结果为空时如实告知用户。",
		func(ctx context.Context, req model.RetrieveReq) (model.RetrieveResp, error) {
			topK := req.TopK
			if topK <= 0 {
				topK = 5
			}
			docs, err := r.Retrieve(ctx, req.Query, topK)
			if err != nil {
				return model.RetrieveResp{}, err
			}
			return model.RetrieveResp{Documents: docs}, nil
		},
	)
}
