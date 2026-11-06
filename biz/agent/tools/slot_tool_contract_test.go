package tools

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/cloudwego/eino/schema"
)

// 工具契约测试：锁定 slot agent 各工具的输入 JSON Schema 关键约束，
// 防止后续改动破坏「语义字段契约」或「后端标识 Go 侧推导」的设计约定。

func toolJSONSchema(t *testing.T, ti schema.ToolInfo) map[string]any {
	t.Helper()
	js, err := ti.ParamsOneOf.ToJSONSchema()
	if err != nil {
		t.Fatalf("ToJSONSchema: %v", err)
	}
	b, err := js.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal schema: %v", err)
	}
	return m
}

func TestSlotToolContract(t *testing.T) {
	if err := initSlotTools(); err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()

	infoOf := func(tool interface {
		Info(ctx context.Context) (*schema.ToolInfo, error)
	}) schema.ToolInfo {
		info, err := tool.Info(ctx)
		if err != nil {
			t.Fatalf("Info: %v", err)
		}
		return *info
	}

	t.Run("generate_setting", func(t *testing.T) {
		m := toolJSONSchema(t, infoOf(GenerateSettingTool))
		props := m["properties"].(map[string]any)

		// 必填：page + resource_type 唯一确定目标资源位
		assertRequired(t, m, "page", "resource_type")

		// 枚举约束
		assertEnum(t, props["page"], "home", "new_house", "second_hand")
		assertEnum(t, props["resource_type"], "banner", "widget", "popup")

		// 语义契约：不暴露后端标识字段（应由 Go 侧推导）
		assertAbsent(t, props, "resource_id", "page_code", "resource_name", "banner_name", "resource_type_code")
	})

	t.Run("query_setting", func(t *testing.T) {
		m := toolJSONSchema(t, infoOf(QuerySettingTool))
		// 所有筛选条件均可选
		if req, ok := m["required"]; ok {
			t.Fatalf("query_setting 不应有必填字段，got required=%v", req)
		}
	})

	t.Run("need_clarify", func(t *testing.T) {
		m := toolJSONSchema(t, infoOf(NeedClarifyTool))
		assertRequired(t, m, "questions")
	})

	t.Run("validate_setting", func(t *testing.T) {
		m := toolJSONSchema(t, infoOf(ValidateSettingTool))
		assertRequired(t, m, "page", "resource_type")
	})
}

func assertRequired(t *testing.T, m map[string]any, keys ...string) {
	t.Helper()
	req, ok := m["required"].([]any)
	if !ok {
		t.Fatalf("缺少 required 数组")
	}
	have := map[string]bool{}
	for _, r := range req {
		have[r.(string)] = true
	}
	for _, k := range keys {
		if !have[k] {
			t.Fatalf("期望字段 %q 为必填，实际 required=%v", k, req)
		}
	}
}

func assertEnum(t *testing.T, prop any, want ...string) {
	t.Helper()
	p := prop.(map[string]any)
	enum, ok := p["enum"].([]any)
	if !ok {
		t.Fatalf("字段缺少 enum，prop=%v", p)
	}
	if len(enum) != len(want) {
		t.Fatalf("enum 长度不一致，got=%v want=%v", enum, want)
	}
	have := map[string]bool{}
	for _, e := range enum {
		have[e.(string)] = true
	}
	for _, w := range want {
		if !have[w] {
			t.Fatalf("enum 缺少 %q，got=%v", w, enum)
		}
	}
}

func assertAbsent(t *testing.T, props map[string]any, keys ...string) {
	t.Helper()
	for _, k := range keys {
		if _, ok := props[k]; ok {
			t.Fatalf("语义契约被破坏：输入 schema 不应暴露后端标识字段 %q", k)
		}
	}
}
