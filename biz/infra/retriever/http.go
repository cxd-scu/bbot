package retriever

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/bytedance/sonic"
	"github.com/cxd-scu/bbot/biz/model"
)

// HTTPRetriever 通过 HTTP 调用知识库检索服务。
//
// 假设的知识库 HTTP 契约（TODO: 与真实知识库接口对齐后调整）：
//   - 请求：POST {baseURL}/retrieve，body: {"query": "...", "top_k": N}
//   - 响应：{"documents": [{"id": "...", "title": "...", "content": "...", "score": 0.9}]}
type HTTPRetriever struct {
	baseURL      string
	retrievePath string
	client       *http.Client
}

// NewHTTPRetriever 创建 HTTP 知识库检索器。
func NewHTTPRetriever(baseURL string) *HTTPRetriever {
	return &HTTPRetriever{
		baseURL:      baseURL,
		retrievePath: "/retrieve",
		client:       &http.Client{Timeout: 5 * time.Second},
	}
}

type retrieveRequest struct {
	Query string `json:"query"`
	TopK  int    `json:"top_k"`
}

type retrieveResponse struct {
	Documents []*model.Document `json:"documents"`
}

// Retrieve 调用知识库 HTTP 接口检索相关文档。
func (r *HTTPRetriever) Retrieve(ctx context.Context, query string, topK int) ([]*model.Document, error) {
	body, err := sonic.MarshalString(retrieveRequest{Query: query, TopK: topK})
	if err != nil {
		return nil, fmt.Errorf("retriever marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, r.baseURL+r.retrievePath, bytes.NewBufferString(body))
	if err != nil {
		return nil, fmt.Errorf("retriever build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("retriever http request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("retriever http status: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("retriever read response: %w", err)
	}

	var out retrieveResponse
	if err := sonic.Unmarshal(data, &out); err != nil {
		return nil, fmt.Errorf("retriever decode response: %w", err)
	}
	return out.Documents, nil
}
