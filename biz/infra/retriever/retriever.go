package retriever

import (
	"context"

	"github.com/cxd-scu/bbot/biz/model"
)

// Retriever 知识库检索接口：给定查询，返回相关的文档片段。
//
// 抽象出该接口是为了解耦「知识库实现」与「oncall agent」：知识库现在通过 HTTP 提供，
// 后续若更换为本地向量库、ES 等实现，只需新增一个实现，agent 侧无需改动。
type Retriever interface {
	Retrieve(ctx context.Context, query string, topK int) ([]*model.Document, error)
}
