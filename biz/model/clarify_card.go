package model

// ClarifyCard 澄清卡：当无法确定用户意图（gateway 返回 need_classify）或 L1 agent 无法处理
// （调用 fallback）时，返回给用户选择意图。
type ClarifyCard struct {
	Type    string       `json:"type"`    // 固定 "clarify_card"
	Title   string       `json:"title"`   // 卡片标题
	Message string       `json:"message"` // 提示文案
	Options []CardOption `json:"options"` // 意图选项
}

// CardOption 澄清卡上的一个意图选项。
type CardOption struct {
	AgentName string `json:"agent_name"` // 对应的 agent 名（slot / oncall）
	Label     string `json:"label"`      // 展示文案
}
