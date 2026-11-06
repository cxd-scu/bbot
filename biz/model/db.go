package model

// AiAgent 表示一个 Agent 实例：单个用户与单个 L1 agent（如 slot/oncall）的完整对话。
// 主键为 (UID, AgentName)：同一用户对不同的 L1 agent 各有一个实例。
type AiAgent struct {
	AgentID   int64  `json:"agent_id"`   // Agent 实例唯一标识
	UID       string `json:"uid"`        // 用户唯一标识（字符串）
	Channel   string `json:"channel"`    // 渠道标识
	AgentName string `json:"agent_name"` // L1 agent 名称（slot / oncall）
	SessionID string `json:"session_id"` // 关联的会话 ID
	Input     string `json:"input"`      // 当前轮输入内容
	InputType int    `json:"input_type"` // 输入类型
}
