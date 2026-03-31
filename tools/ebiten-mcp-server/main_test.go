package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	ebitenmcp "github.com/kimyechan/ebiten-aio-framework/libs/go/ebiten-mcp"
)

func TestNewHTTPHandlerServesHealthAndMCPPath(t *testing.T) {
	server := ebitenmcp.New(ebitenmcp.NewBridgeClient("http://127.0.0.1:47831"))
	handler := newHTTPHandler(server, "/mcp")

	healthRequest := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	healthRecorder := httptest.NewRecorder()
	handler.ServeHTTP(healthRecorder, healthRequest)
	if healthRecorder.Code != http.StatusOK {
		t.Fatalf("health status = %d, want %d", healthRecorder.Code, http.StatusOK)
	}

	mcpRequest := httptest.NewRequest(http.MethodPost, "/mcp", nil)
	mcpRequest.Header.Set("Accept", "application/json,text/event-stream")
	mcpRequest.Header.Set("Content-Type", "application/json")
	mcpRequest.Header.Set("Mcp-Protocol-Version", "2025-03-26")
	mcpRecorder := httptest.NewRecorder()
	handler.ServeHTTP(mcpRecorder, mcpRequest)
	if mcpRecorder.Code == http.StatusNotFound {
		t.Fatalf("expected mcp path to be routed")
	}
}

func TestNormalizeHTTPPath(t *testing.T) {
	cases := map[string]string{
		"":        "/mcp",
		"/mcp":    "/mcp",
		"mcp":     "/mcp",
		"/x/y/":   "/x/y",
		"foo/bar": "/foo/bar",
	}

	for input, want := range cases {
		if got := normalizeHTTPPath(input); got != want {
			t.Fatalf("normalizeHTTPPath(%q) = %q, want %q", input, got, want)
		}
	}
}
