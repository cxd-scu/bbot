package chain

import (
	"context"
	"fmt"
	"strings"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/cxd-scu/bbot/biz/infra/llm"
	"github.com/cxd-scu/bbot/biz/model"
)

func ShortMemorySummary(ctx context.Context, req *model.SummarizeInput) (string, error) {
	return shortMemorySummarizer.Invoke(ctx, req)
}

var shortMemorySummarizer compose.Runnable[*model.SummarizeInput, string]

func initShortMemorySummarizer(ctx context.Context) error {
	var err error
	chain := compose.NewChain[*model.SummarizeInput, string]()
	chain.
		AppendLambda(compose.InvokableLambda(summarizeTmpl), compose.WithNodeName("summarize_template")).
		AppendChatModel(llm.DsFlash(), compose.WithNodeName("summarize_model")).
		AppendLambda(compose.InvokableLambda(extractContentAsString), compose.WithNodeName("extract_summary_text"))

	shortMemorySummarizer, err = chain.Compile(ctx)
	return err
}

func summarizeTmpl(ctx context.Context, input *model.SummarizeInput) ([]*schema.Message, error) {
	oldSummaryPart := "（无,这是首次总结）"
	if input.OldSummary != "" {
		oldSummaryPart = input.OldSummary
	}
	userContent := fmt.Sprintf(
		"之前的总结:\n%s\n\n以下为新的需要总结的对话:\n%s",
		oldSummaryPart,
		formatTurns(input.ToFold),
	)

	return []*schema.Message{
		schema.SystemMessage(summarizerPrompt),
		schema.UserMessage(userContent),
	}, nil
}

func extractContentAsString(ctx context.Context, msg *schema.Message) (string, error) {
	return msg.Content, nil
}

func formatTurns(turns []model.Turn) string {
	var sb strings.Builder

	for _, t := range turns {
		if t.UserMsg == nil || t.AssistantMsg == nil {
			continue
		}
		sb.WriteString("用户: ")
		sb.WriteString(t.UserMsg.Content)
		sb.WriteByte('\n')

		sb.WriteString("助手: ")
		sb.WriteString(t.AssistantMsg.Content)
		sb.WriteString("\n\n")
	}

	return sb.String()
}

const summarizerPrompt = `你是对话记忆总结助手。把下面的对话压缩成简洁的摘要，保留关键事实、用户意图、已确认的参数和未完成的事项，供后续对话继续使用。摘要用中文，控制在几句话内。`
