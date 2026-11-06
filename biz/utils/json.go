package utils

import (
	"context"
	"fmt"
	"strings"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/eino/schema"
	"github.com/cxd-scu/bbot/biz/constant"
	"github.com/cxd-scu/bbot/biz/model"
)

// ExtractClarifyCard 从 agent 最终消息中提取澄清卡（fallback 触发时）。无卡返回 nil。
func ExtractClarifyCard(msg *schema.Message) *model.ClarifyCard {
	if msg == nil {
		return nil
	}
	clean := ExtractJSON(msg.Content)
	var wrapper struct {
		Card *model.ClarifyCard `json:"card"`
	}
	if err := sonic.UnmarshalString(clean, &wrapper); err != nil || wrapper.Card == nil {
		return nil
	}
	return wrapper.Card
}

// ExtractJSON 从模型输出文本中裁剪出首个 JSON 对象（首个 { 到最后一个 }）。
func ExtractJSON(content string) string {
	start := strings.Index(content, "{")
	end := strings.LastIndex(content, "}")
	if start == -1 || end == -1 || start >= end {
		return content
	}
	return content[start : end+1]
}

func ParseJSONToIOParam(ctx context.Context, msg *schema.Message) (*model.IOParam, error) {
	// 裁剪出干净的 JSON 字符串
	cleanJSON := ExtractJSON(msg.Content)

	res := &model.IOParam{
		Params:    make(map[string]interface{}),
		NodeState: constant.NodeStateInit,
	}

	// 1. 将整个 JSON 的所有字段解析到 Params map 中
	if err := sonic.UnmarshalString(cleanJSON, &res.Params); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON into Params: %w, raw_clipped: %s", err, cleanJSON)
	}

	return res, nil
}

func ToString(v any) string {
	s, _ := sonic.MarshalString(v)
	return s
}
