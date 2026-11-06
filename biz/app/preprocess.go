package app

import (
	"context"
	"strings"

	"github.com/cxd-scu/bbot/biz/agent/chain"
	"github.com/cxd-scu/bbot/biz/model"
)

// NormalizeQuery 对用户输入做预处理，产出可路由、可回答的纯文本提问。
//   - 无图片：原样返回用户文字。
//   - 有图片：调用多模态大模型直接输出结构化 JSON（含优化提问），返回优化后的提问。
//
// 下游意图路由（GatewayNlu）只消费纯文本，slot/oncall 无需感知图片。
func NormalizeQuery(ctx context.Context, req *model.UserReq) (string, error) {
	if len(req.Images) == 0 {
		return req.Query, nil
	}

	urls := make([]string, 0, len(req.Images))
	for _, img := range req.Images {
		urls = append(urls, img.URL)
	}

	content, err := chain.AnalyzeImages(ctx, &model.ImageAnalyzeInput{
		UserText:  req.Query,
		ImageURLs: urls,
	})
	if err != nil {
		return "", err
	}

	// 兜底：多模态模型未能产出优化提问时，退回用户原文字。
	if content == nil || strings.TrimSpace(content.OptimizedQuery) == "" {
		return req.Query, nil
	}
	return content.OptimizedQuery, nil
}
