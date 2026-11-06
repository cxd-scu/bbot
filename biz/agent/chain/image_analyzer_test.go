package chain

import (
	"context"
	"testing"

	"github.com/cloudwego/eino/schema"
	"github.com/cxd-scu/bbot/biz/model"
)

func TestImageAnalyzeTmpl(t *testing.T) {
	ctx := context.Background()

	msgs, err := imageAnalyzeTmpl(ctx, &model.ImageAnalyzeInput{
		UserText:  "帮我看看这个配置",
		ImageURLs: []string{"https://x.com/a.png", "https://x.com/b.png"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(msgs) != 2 {
		t.Fatalf("期望 2 条消息（system + user），got=%d", len(msgs))
	}

	if msgs[0].Role != schema.System {
		t.Fatalf("第一条应为 system，got=%v", msgs[0].Role)
	}

	user := msgs[1]
	if user.Role != schema.User {
		t.Fatalf("第二条应为 user，got=%v", user.Role)
	}
	if len(user.UserInputMultiContent) != 3 {
		t.Fatalf("期望 3 个 part（1 text + 2 image），got=%d", len(user.UserInputMultiContent))
	}

	// 第一个 part 是 text
	if user.UserInputMultiContent[0].Type != schema.ChatMessagePartTypeText {
		t.Fatalf("第 1 个 part 应为 text，got=%v", user.UserInputMultiContent[0].Type)
	}
	if user.UserInputMultiContent[0].Text != "帮我看看这个配置" {
		t.Fatalf("text 内容不符，got=%q", user.UserInputMultiContent[0].Text)
	}

	// 后两个 part 是 image_url，URL 正确
	for i, want := range []string{"https://x.com/a.png", "https://x.com/b.png"} {
		p := user.UserInputMultiContent[i+1]
		if p.Type != schema.ChatMessagePartTypeImageURL {
			t.Fatalf("第 %d 个 part 应为 image_url，got=%v", i+2, p.Type)
		}
		if p.Image == nil || p.Image.URL == nil || *p.Image.URL != want {
			t.Fatalf("第 %d 个 image URL 不符，got=%v", i+2, p.Image)
		}
	}
}

func TestImageAnalyzeTmpl_NoText(t *testing.T) {
	ctx := context.Background()

	msgs, err := imageAnalyzeTmpl(ctx, &model.ImageAnalyzeInput{
		ImageURLs: []string{"https://x.com/a.png"},
	})
	if err != nil {
		t.Fatal(err)
	}
	user := msgs[1]
	if len(user.UserInputMultiContent) != 1 {
		t.Fatalf("无文字时只应有 1 个 image part，got=%d", len(user.UserInputMultiContent))
	}
	if user.UserInputMultiContent[0].Type != schema.ChatMessagePartTypeImageURL {
		t.Fatalf("应为 image_url，got=%v", user.UserInputMultiContent[0].Type)
	}
}

func TestParseImageContent(t *testing.T) {
	ctx := context.Background()

	raw := "```json\n{\"type\":\"config_screenshot\",\"fields\":{\"city_ids\":[100000]},\"optimized_query\":\"首页banner配置\"}\n```"
	got, err := parseImageContent(ctx, &schema.Message{Content: raw})
	if err != nil {
		t.Fatal(err)
	}

	if got.Type != "config_screenshot" {
		t.Fatalf("type 解析不符，got=%q", got.Type)
	}
	if got.OptimizedQuery != "首页banner配置" {
		t.Fatalf("optimized_query 解析不符，got=%q", got.OptimizedQuery)
	}
	if got.Fields == nil {
		t.Fatal("fields 应为非 nil")
	}
}
