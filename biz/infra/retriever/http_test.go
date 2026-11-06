package retriever

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTPRetriever_Retrieve(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/retrieve" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("unexpected method: %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"documents":[{"id":"1","title":"值班安排","content":"值班时间是 9:00-18:00","score":0.9}]}`))
	}))
	defer srv.Close()

	r := NewHTTPRetriever(srv.URL)
	docs, err := r.Retrieve(context.Background(), "值班时间", 3)
	if err != nil {
		t.Fatal(err)
	}
	if len(docs) != 1 {
		t.Fatalf("期望 1 条文档，got=%d", len(docs))
	}
	if docs[0].Content != "值班时间是 9:00-18:00" {
		t.Fatalf("content 不符，got=%q", docs[0].Content)
	}
	if docs[0].Title != "值班安排" {
		t.Fatalf("title 不符，got=%q", docs[0].Title)
	}
}

func TestHTTPRetriever_Error(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	r := NewHTTPRetriever(srv.URL)
	if _, err := r.Retrieve(context.Background(), "x", 3); err == nil {
		t.Fatal("期望非 200 返回错误")
	}
}
