package tools

import (
	"context"
	"strings"
	"testing"

	"github.com/cxd-scu/bbot/biz/model"
)

type fakeRetriever struct {
	docs []*model.Document
	err  error
}

func (f *fakeRetriever) Retrieve(ctx context.Context, query string, topK int) ([]*model.Document, error) {
	return f.docs, f.err
}

func TestRetrieveTool(t *testing.T) {
	fake := &fakeRetriever{docs: []*model.Document{{Content: "值班时间是 9:00-18:00"}}}
	tl, err := NewRetrieveTool(fake)
	if err != nil {
		t.Fatal(err)
	}

	out, err := tl.InvokableRun(context.Background(), `{"query":"值班时间","top_k":3}`)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "值班时间是 9:00-18:00") {
		t.Fatalf("输出应包含检索内容，got=%q", out)
	}
}

func TestRetrieveTool_DefaultTopK(t *testing.T) {
	fake := &fakeRetriever{docs: []*model.Document{{Content: "x"}}}
	tl, err := NewRetrieveTool(fake)
	if err != nil {
		t.Fatal(err)
	}

	out, err := tl.InvokableRun(context.Background(), `{"query":"x"}`)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "x") {
		t.Fatalf("缺省 top_k 时应返回文档，got=%q", out)
	}
}
