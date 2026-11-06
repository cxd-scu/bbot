package model

type UserReq struct {
	Query     string  `json:"query"`
	Images    []Image `json:"images,omitempty"` // 用户可能直接发图，按 URL 处理
	UID       string  `json:"uid"`
	Channel   string  `json:"channel"`
	AgentName string  `json:"agent_name"`
}

// Image 用户输入中的图片，仅按 URL 处理。
type Image struct {
	URL string `json:"url"`
}

// ImageContent 多模态大模型对「图片 + 用户文字」的结构化识别结果，直接以 JSON 输出。
// 多模态模型一次完成「看图 + 信息提取 + 提问优化」，无需单独 OCR 环节。
type ImageContent struct {
	Type           string         `json:"type"`            // 图片类型：config_screenshot / error_screenshot / schedule / doc / design / other
	Fields         map[string]any `json:"fields"`          // 提取到的结构化字段，按图片类型而定（如 city_ids、image_url、error_message 等）
	OptimizedQuery string         `json:"optimized_query"` // 结合用户文字优化后的完整提问，供 gateway 意图路由使用
}

// ImageAnalyzeInput 图片分析链的输入：用户文字 + 图片 URL 列表。
type ImageAnalyzeInput struct {
	UserText  string   `json:"user_text"`
	ImageURLs []string `json:"image_urls"`
}

type GatewayAgentOutput struct {
	TargetAgent string  `json:"target_agent"`
	Reason      string  `json:"reason"`
	Confidence  float64 `json:"confidence"`
}
