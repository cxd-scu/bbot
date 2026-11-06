package app

import (
	"context"
	"errors"
	"io"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/schema"
	"github.com/cxd-scu/bbot/biz/agent/graph"
)

// StreamDelta 流式回调：收到一段增量文本。
type StreamDelta func(delta string) error

// StreamSlotAgent 流式运行 slot agent，逐段回调 token 增量。
func StreamSlotAgent(ctx context.Context, query string, emit StreamDelta) error {
	return streamAgent(graph.SlotAgentStreamRunner.Run(ctx, []*schema.Message{
		{Role: schema.User, Content: query},
	}), emit)
}

// StreamOncallAgent 流式运行 oncall agent，逐段回调 token 增量。
func StreamOncallAgent(ctx context.Context, query string, emit StreamDelta) error {
	return streamAgent(graph.OncallStreamRunner.Run(ctx, []*schema.Message{
		{Role: schema.User, Content: query},
	}), emit)
}

// streamAgent 迭代 agent 事件流，把模型输出的 token 增量逐个交给 emit。
func streamAgent(iter *adk.AsyncIterator[*adk.TypedAgentEvent[*schema.Message]], emit StreamDelta) error {
	for {
		ev, ok := iter.Next()
		if !ok {
			return nil
		}
		if ev.Err != nil {
			return ev.Err
		}
		if ev.Output == nil || ev.Output.MessageOutput == nil {
			continue
		}

		mv := ev.Output.MessageOutput
		if !mv.IsStreaming || mv.Role != schema.Assistant || mv.MessageStream == nil {
			continue
		}
		if err := emitMessageStream(mv.MessageStream, emit); err != nil {
			return err
		}
	}
}

// emitMessageStream 消费一段模型输出流，把每块内容交给 emit。
func emitMessageStream(stream *schema.StreamReader[*schema.Message], emit StreamDelta) error {
	defer stream.Close()
	for {
		chunk, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return err
		}
		if chunk != nil && chunk.Content != "" {
			if err := emit(chunk.Content); err != nil {
				return err
			}
		}
	}
}
