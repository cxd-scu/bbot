package model

// Document 知识库检索返回的文档片段。
type Document struct {
	ID      string         `json:"id,omitempty"`
	Title   string         `json:"title,omitempty"`
	Content string         `json:"content"`
	Score   float64        `json:"score,omitempty"`
	Extra   map[string]any `json:"extra,omitempty"`
}

// RetrieveReq 知识库检索工具（retrieve_knowledge）的输入。
type RetrieveReq struct {
	Query string `json:"query" jsonschema_description:"检索查询词，应提炼用户问题中的关键信息"`
	TopK  int    `json:"top_k,omitempty" jsonschema_description:"返回文档数量，缺省 5"`
}

// RetrieveResp 知识库检索工具的输出。
type RetrieveResp struct {
	Documents []*Document `json:"documents" jsonschema_description:"检索到的文档片段列表"`
}
