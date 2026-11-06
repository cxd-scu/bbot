package llm

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cxd-scu/bbot/biz/constant"
)

// VLModel 返回多模态（视觉语言）大模型，用于图片识别与提问优化。
func VLModel() *openai.ChatModel {
	return vlModel
}

var vlModel *openai.ChatModel

func initVLModel(ctx context.Context) error {
	var err error
	vlModel, err = openai.NewChatModel(ctx, &openai.ChatModelConfig{
		BaseURL:     constant.VLBaseURL,
		APIKey:      constant.VLAPIKey,
		Model:       constant.VLModelName,
		Temperature: new(float32), // 0，保证识别结果稳定可复现
	})
	if err != nil {
		return fmt.Errorf("failed to init vision-language model: %w", err)
	}
	return nil
}
