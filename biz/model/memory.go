package model

import (
	"fmt"

	"github.com/cloudwego/eino/schema"
)

type Turn struct {
	UserMsg      *schema.Message // role: user
	AssistantMsg *schema.Message // role: assistant,只存最终回复,不含 ToolCalls
}

type Session struct {
	SessionID    string
	Summary      string
	Window       []Turn
	CheckPointID string
	Pending      bool
}

type SummarizeInput struct {
	OldSummary string
	ToFold     []Turn
}

// BuildContext 组装一次 conductor 调用所需要的完整 messages
func (s *Session) BuildContext() []*schema.Message {
	// 预估容量：system + summary(可选) + window*2 + current
	capacity := 2 + len(s.Window)*2
	messages := make([]*schema.Message, 0, capacity)

	if s.Summary != "" {
		messages = append(messages, schema.SystemMessage(
			fmt.Sprintf("此前对话摘要：\n%s", s.Summary),
		))
	}

	for _, turn := range s.Window {
		if turn.UserMsg == nil || turn.AssistantMsg == nil {
			continue // 防御性跳过脏数据，不让一条坏 Turn 拖垮整次组装
		}
		messages = append(messages, turn.UserMsg, turn.AssistantMsg)
	}

	return messages
}
