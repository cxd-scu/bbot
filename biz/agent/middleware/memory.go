package middleware

import (
	"context"
	"sync"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/schema"
	"github.com/cxd-scu/bbot/biz/agent/chain"
	"github.com/cxd-scu/bbot/biz/constant"
	"github.com/cxd-scu/bbot/biz/infra/cache"
	"github.com/cxd-scu/bbot/biz/model"
)

var ShortMemoryManager = &ShortMemoryMiddleware{
	BaseChatModelAgentMiddleware: &adk.BaseChatModelAgentMiddleware{},
}

type ShortMemoryMiddleware struct {
	*adk.BaseChatModelAgentMiddleware
}

func (m *ShortMemoryMiddleware) BeforeModelRewriteState(ctx context.Context, state *adk.ChatModelAgentState, mc *adk.ModelContext) (context.Context, *adk.ChatModelAgentState, error) {
	// 已注入过历史，跳过
	if _, ok := adk.GetSessionValue(ctx, constant.InjectIndexKey); ok {
		return ctx, state, nil
	}

	// 会话 ID 未注入（未开启记忆），直接跳过，避免类型断言 panic
	sessionIDRaw, ok := adk.GetSessionValue(ctx, constant.SessionIDKey)
	if !ok {
		return ctx, state, nil
	}
	sessionID, _ := sessionIDRaw.(string)
	if sessionID == "" {
		return ctx, state, nil
	}

	session := cache.GetOrCreateSession(sessionID)

	insertAt := insertIndex(state.Messages)
	historyMsgs := session.BuildContext()
	if len(historyMsgs) == 0 {
		return ctx, state, nil
	}

	newMsgs := make([]*schema.Message, 0, len(state.Messages)+len(historyMsgs))
	newMsgs = append(newMsgs, state.Messages[:insertAt]...)
	newMsgs = append(newMsgs, historyMsgs...)
	newMsgs = append(newMsgs, state.Messages[insertAt:]...)

	state.Messages = newMsgs
	adk.AddSessionValue(ctx, constant.InjectIndexKey, len(historyMsgs)+insertAt)
	logger.Info("[MemoryInject] session=%s turn_injected=true history_msg_count=%d final_msg_count=%d",
		sessionID, len(historyMsgs), len(newMsgs))

	return ctx, state, nil
}

func (m *ShortMemoryMiddleware) AfterAgent(
	ctx context.Context,
	state *adk.ChatModelAgentState,
) (context.Context, error) {
	sessionIDRaw, ok := adk.GetSessionValue(ctx, constant.SessionIDKey)
	if !ok {
		return ctx, nil
	}
	sessionID, _ := sessionIDRaw.(string)
	if sessionID == "" {
		return ctx, nil
	}

	injectIndexRaw, ok := adk.GetSessionValue(ctx, constant.InjectIndexKey)
	if !ok {
		return ctx, nil
	}
	injectIndex, _ := injectIndexRaw.(int)
	if injectIndex < 0 || injectIndex >= len(state.Messages) {
		return ctx, nil
	}

	turn := model.Turn{}
	for _, m := range state.Messages[injectIndex:] {
		switch m.Role {
		case schema.Assistant:
			// 跳过工具调用消息，只存最终回复（不含 ToolCalls）
			if len(m.ToolCalls) > 0 {
				continue
			}
			turn.AssistantMsg = m
		case schema.User:
			turn.UserMsg = m
		}
	}

	// 只存完整一轮（user + assistant）
	if turn.UserMsg != nil && turn.AssistantMsg != nil {
		AppendTurnToSession(ctx, sessionID, turn)
	}

	return ctx, nil
}

func insertIndex(msgs []*schema.Message) int {
	for i, m := range msgs {
		if m.Role == schema.System {
			return i + 1
		}
	}
	return len(msgs)
}

const (
	threshold  = 7
	keepRecent = 3
)

func AppendTurnToSession(ctx context.Context, sessionID string, turn model.Turn) {
	mu := lockFor(sessionID)
	mu.Lock()
	s := cache.GetOrCreateSession(sessionID)
	s.Window = append(s.Window, turn)
	needCompress := len(s.Window) >= threshold
	cache.Set(sessionID, s) // 每次append都显式落盘
	mu.Unlock()

	if needCompress {
		go compressSession(ctx, sessionID)
	}
}

func compressSession(ctx context.Context, sessionID string) {
	mu := lockFor(sessionID)

	mu.Lock()
	s := cache.GetOrCreateSession(sessionID)
	if len(s.Window) < threshold {
		mu.Unlock()
		return // 已被别的goroutine处理过，双重检查
	}
	foldCount := len(s.Window) - keepRecent
	toFold := make([]model.Turn, foldCount)
	copy(toFold, s.Window[:foldCount]) // 拷贝出来，避免持锁期间等LLM
	oldSummary := s.Summary
	mu.Unlock()

	// 不持锁调用LLM，避免长时间阻塞其他并发的append
	newSummary, err := chain.ShortMemorySummary(ctx, &model.SummarizeInput{
		ToFold:     toFold,
		OldSummary: oldSummary,
	})
	if err != nil {
		return // 压缩失败，window继续超预算，下次append再触发
	}

	mu.Lock()
	defer mu.Unlock()
	s = cache.GetOrCreateSession(sessionID)
	// 压缩期间可能有新turn被append进来，只移除已经被折叠的那部分，不能整体覆盖
	if len(s.Window) >= foldCount {
		s.Window = s.Window[foldCount:]
	}
	s.Summary = newSummary
	cache.Set(sessionID, s)
}

var sessionLocks sync.Map // sessionID -> *sync.Mutex
func lockFor(sessionID string) *sync.Mutex {
	l, _ := sessionLocks.LoadOrStore(sessionID, &sync.Mutex{})
	return l.(*sync.Mutex)
}
