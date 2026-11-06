package llm

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino-ext/components/model/openai"
	toolutils "github.com/cloudwego/eino/components/tool/utils"
	"github.com/cxd-scu/bbot/biz/constant"
	"github.com/cxd-scu/bbot/biz/model"
)

func DsFlash() *openai.ChatModel {
	return flashModel
}

func DsPro() *openai.ChatModel {
	return proModel
}

// GatewayModel 返回 gateway 专用的模型：通过 response_format 强制结构化输出。
func GatewayModel() *openai.ChatModel {
	return gatewayModel
}

var (
	flashModel   *openai.ChatModel
	proModel     *openai.ChatModel
	gatewayModel *openai.ChatModel
)

func dsModelConfig(model string, temperature float32) *openai.ChatModelConfig {
	return &openai.ChatModelConfig{
		BaseURL:     "https://api.deepseek.com",
		APIKey:      constant.DSAPIKey,
		Model:       model,
		Temperature: new(temperature),
	}
}

// initGatewayModel 初始化 gateway 模型：由 Go 结构体生成 JSON Schema，
// 通过 response_format=json_schema 让模型在解码层面就被约束输出结构。
func initGatewayModel(ctx context.Context) error {
	paramsOneOf, err := toolutils.GoStruct2ParamsOneOf[model.GatewayAgentOutput]()
	if err != nil {
		return fmt.Errorf("gateway output schema: %w", err)
	}
	js, err := paramsOneOf.ToJSONSchema()
	if err != nil {
		return fmt.Errorf("gateway output jsonschema: %w", err)
	}

	gatewayModel, err = openai.NewChatModel(ctx, &openai.ChatModelConfig{
		BaseURL:     "https://api.deepseek.com",
		APIKey:      constant.DSAPIKey,
		Model:       "deepseek-flash",
		Temperature: new(float32), // 0，稳定可复现
		ResponseFormat: &openai.ChatCompletionResponseFormat{
			Type: openai.ChatCompletionResponseFormatTypeJSONSchema,
			JSONSchema: &openai.ChatCompletionResponseFormatJSONSchema{
				Name:        "gateway_output",
				Description: "意图路由结果",
				Strict:      true,
				JSONSchema:  js,
			},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to init gateway model: %w", err)
	}
	return nil
}

func InitModel(ctx context.Context) {
	var err error
	flashModel, err = openai.NewChatModel(ctx, dsModelConfig("deepseek-flash", 0.1))
	if err != nil {
		panic(fmt.Errorf("failed to init flash model: %v", err))
	}

	proModel, err = openai.NewChatModel(ctx, dsModelConfig("deepseek-v4-pro", 0.1))
	if err != nil {
		panic(fmt.Errorf("failed to init pro model: %v", err))
	}

	if err = initVLModel(ctx); err != nil {
		panic(err)
	}

	if err = initGatewayModel(ctx); err != nil {
		panic(err)
	}
}
