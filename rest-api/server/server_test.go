package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"app/datasources"
)

func TestGetStatus(t *testing.T) {
	app := NewServer(context.Background(), &datasources.DataSources{})

	req := httptest.NewRequest("GET", "/api/status", nil)
	resp := httptest.NewRecorder()
	app.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, resp.Code)
	}
	if resp.Body.String() != "ok" {
		t.Fatalf("expected body %q, got %q", "ok", resp.Body.String())
	}
}
