package model

import "github.com/cxd-scu/bbot/biz/constant"

type UserInput struct {
	Type constant.InputType `json:"type"`
	Body string             `json:"body"`
}

// AgentContext 包含 Agent 运行时的会话状态与上下文信息。
type AgentContext struct {

	// -- 基础标识 (8字节/指针对齐字段优先) --

	AgentID int64 `json:"agent_id"` // Agent 唯一标识
	UID     int64 `json:"uid"`      // 用户 唯一标识

	// -- 引用指针 (8字节) --

	UserInput   *UserInput `json:"user_input,omitempty"` // 用户当前输入的上下文
	ShortMemory *Session   `json:"memory,omitempty"`     // 历史记忆/上下文状态

	// -- 字符串字段 (16字节) --

	SessionID string `json:"session_id"` // 会话 ID
	AgentName string `json:"agent_name"` // Agent 名称
	Channel   string `json:"channel"`    // 渠道标识
}

type IOParam struct {
	NodeState constant.NodeState     `json:"node_state"`
	Params    map[string]interface{} `json:"params"`
}
