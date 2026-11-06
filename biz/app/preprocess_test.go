package app

import (
	"context"
	"testing"

	"github.com/cxd-scu/bbot/biz/model"
)

// 图片路径需要真实多模态模型，这里只覆盖无需模型的「无图片原样透传」分支；
// 图片的消息构造与 JSON 解析在 biz/agent/chain 包内单测。
func TestNormalizeQuery_NoImage(t *testing.T) {
	ctx := context.Background()

	req := &model.UserReq{Query: "帮我配个首页banner"}
	got, err := NormalizeQuery(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	if got != "帮我配个首页banner" {
		t.Fatalf("期望原样返回，got=%q", got)
	}
}
