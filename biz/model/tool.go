package model

// 通用工具（非 slot 专属）的输入/输出契约类型。
// 所有工具用到的结构体统一收敛到 model 包，便于作为契约单一事实来源。

// FallbackReq fallback 工具的输入：说明为何无法处理当前请求。
type FallbackReq struct {
	Reason string `json:"reason" jsonschema_description:"说明为何无法处理当前请求，或引导用户转人工/联系谁"`
}

// FallbackResp fallback 工具的输出：兜底回复话术，可携带澄清卡。
type FallbackResp struct {
	Reply string       `json:"reply" jsonschema_description:"回复给用户的兜底话术"`
	Card  *ClarifyCard `json:"card,omitempty" jsonschema_description:"需要用户选择意图时返回的澄清卡"`
}
