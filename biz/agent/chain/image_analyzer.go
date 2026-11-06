package chain

import (
	"context"
	"fmt"
	"strings"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/cxd-scu/bbot/biz/infra/llm"
	"github.com/cxd-scu/bbot/biz/model"
	"github.com/cxd-scu/bbot/biz/utils"
)

// AnalyzeImages 用多模态大模型分析「图片 + 用户文字」，直接输出结构化 JSON（model.ImageContent）。
// 其中 OptimizedQuery 为结合图片内容优化后的提问，供 gateway 意图路由使用。
func AnalyzeImages(ctx context.Context, input *model.ImageAnalyzeInput) (*model.ImageContent, error) {
	return imageAnalyzer.Invoke(ctx, input)
}

var imageAnalyzer compose.Runnable[*model.ImageAnalyzeInput, *model.ImageContent]

func initImageAnalyzer(ctx context.Context) error {
	var err error
	chain := compose.NewChain[*model.ImageAnalyzeInput, *model.ImageContent]()
	chain.
		AppendLambda(compose.InvokableLambda(imageAnalyzeTmpl), compose.WithNodeName("image_analyze_template")).
		AppendChatModel(llm.VLModel(), compose.WithNodeName("image_analyze_model")).
		AppendLambda(compose.InvokableLambda(parseImageContent), compose.WithNodeName("image_analyze_parser"))

	imageAnalyzer, err = chain.Compile(ctx)
	return err
}

// imageAnalyzeTmpl 构造多模态消息：系统提示 + 用户文字/图片。
func imageAnalyzeTmpl(ctx context.Context, input *model.ImageAnalyzeInput) ([]*schema.Message, error) {
	parts := make([]schema.MessageInputPart, 0, len(input.ImageURLs)+1)

	if strings.TrimSpace(input.UserText) != "" {
		parts = append(parts, schema.MessageInputPart{
			Type: schema.ChatMessagePartTypeText,
			Text: input.UserText,
		})
	}

	for _, u := range input.ImageURLs {
		url := u
		parts = append(parts, schema.MessageInputPart{
			Type: schema.ChatMessagePartTypeImageURL,
			Image: &schema.MessageInputImage{
				MessagePartCommon: schema.MessagePartCommon{
					URL: &url,
				},
			},
		})
	}

	return []*schema.Message{
		schema.SystemMessage(imageAnalyzerPrompt),
		{
			Role:                  schema.User,
			UserInputMultiContent: parts,
		},
	}, nil
}

// parseImageContent 解析模型输出的 JSON 为 model.ImageContent。
func parseImageContent(ctx context.Context, msg *schema.Message) (*model.ImageContent, error) {
	clean := utils.ExtractJSON(msg.Content)
	res := &model.ImageContent{}
	if err := sonic.UnmarshalString(clean, res); err != nil {
		return nil, fmt.Errorf("parse image content json: %w, raw: %s", err, clean)
	}
	return res, nil
}

const imageAnalyzerPrompt = `你是运营助手的图片理解模块。用户可能发来图片和文字，请结合两者，只输出一个 JSON（不要输出其他内容，不要用 markdown 代码块）。

请完成：
1. 判断图片类型 type，取其一：config_screenshot(配置截图) / error_screenshot(报错截图) / schedule(值班表) / doc(文档) / design(设计图) / other(其他)。
2. 提取图片中的关键信息到 fields（按类型而定，例如 city_ids、image_url、open_url、error_message 等）。URL、ID、错误码等关键字段必须逐字提取，不得改写、不得推测；拿不准就写"未识别"。
3. 结合用户文字与图片信息，输出 optimized_query：一条完整、可路由、可回答的提问，补全模糊指代。

输出格式：
{
  "type": "config_screenshot",
  "fields": { ... },
  "optimized_query": "..."
}`
