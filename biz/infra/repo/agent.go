package repo

import (
	"context"
	"sync"

	"github.com/cxd-scu/bbot/biz/model"
	"github.com/cxd-scu/bbot/biz/utils"
)

// 内存实现（最简版）：仅用于打通流程，真实存储（DB/缓存）后续替换。
// 进程内 map 存储，加读写锁保证并发安全。

var (
	mu     sync.RWMutex
	byID         = make(map[int64]*model.AiAgent)  // key: AgentID
	byKey        = make(map[string]*model.AiAgent) // key: uid + ":" + agentName
	nextID int64 = 1
)

func agentKey(uid, agentName string) string {
	return uid + ":" + agentName
}

// CreateAgent 创建 Agent 实例；同一 (uid, agentName) 已存在时直接复用，保证幂等。
func CreateAgent(ctx context.Context, uid string, agentName, channel string) (*model.AiAgent, error) {
	mu.Lock()
	defer mu.Unlock()

	key := agentKey(uid, agentName)
	if a, ok := byKey[key]; ok {
		return a, nil
	}

	a := &model.AiAgent{
		AgentID:   nextID,
		UID:       uid,
		Channel:   channel,
		AgentName: agentName,
		SessionID: utils.GenerateSessionID(nextID),
	}
	nextID++
	byID[a.AgentID] = a
	byKey[key] = a
	return a, nil
}

// GetAgentByID 按 AgentID 查询，不存在返回 (nil, nil)。
func GetAgentByID(ctx context.Context, agentID int64) (*model.AiAgent, error) {
	mu.RLock()
	defer mu.RUnlock()
	return byID[agentID], nil
}

// GetAgentByUIDAndName 按 (uid, agentName) 查询，不存在返回 (nil, nil)。
func GetAgentByUIDAndName(ctx context.Context, uid, agentName string) (*model.AiAgent, error) {
	mu.RLock()
	defer mu.RUnlock()
	return byKey[agentKey(uid, agentName)], nil
}
