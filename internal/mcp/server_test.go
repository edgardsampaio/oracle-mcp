package mcp

import (
	"context"
	"io"
	"log/slog"
	"net/http/httptest"
	"testing"
	"time"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestHandshakeListsDimensionTools(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	ts := httptest.NewServer(Handler(logger, nil, nil, "", nil, 1000))
	t.Cleanup(ts.Close)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	t.Cleanup(cancel)

	client := sdk.NewClient(&sdk.Implementation{Name: "oracle-mcp-test", Version: "0.0.1"}, nil)
	session, err := client.Connect(ctx, &sdk.StreamableClientTransport{Endpoint: ts.URL}, nil)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	t.Cleanup(func() { _ = session.Close() })

	info := session.InitializeResult().ServerInfo
	if info == nil || info.Name != serverName {
		t.Fatalf("server info = %#v", info)
	}
	tools, err := session.ListTools(ctx, nil)
	if err != nil {
		t.Fatalf("list tools: %v", err)
	}
	want := map[string]bool{"boas_vindas": true}
	for _, name := range businessTools {
		want[name] = true
	}
	if len(tools.Tools) != len(want) {
		t.Fatalf("tools = %d", len(tools.Tools))
	}
	for _, tool := range tools.Tools {
		if !want[tool.Name] {
			t.Fatalf("tool inesperada %s", tool.Name)
		}
	}
}
