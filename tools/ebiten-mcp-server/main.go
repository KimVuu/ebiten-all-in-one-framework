package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	ebitenmcp "github.com/kimyechan/ebiten-aio-framework/libs/go/ebiten-mcp"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const defaultBridgeAddr = "127.0.0.1:47831"
const defaultHTTPListenAddr = "127.0.0.1:47830"
const defaultHTTPPath = "/mcp"

func main() {
	addrFlag := flag.String("addr", "", "ebiten debug bridge address")
	transportFlag := flag.String("transport", "", "mcp transport: stdio or streamable-http")
	listenFlag := flag.String("listen", "", "streamable HTTP listen address")
	pathFlag := flag.String("path", "", "streamable HTTP path")
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	server := ebitenmcp.New(ebitenmcp.NewBridgeClient(resolveBridgeAddr(*addrFlag)))
	if err := run(ctx, server, resolveTransport(*transportFlag), resolveListenAddr(*listenFlag), normalizeHTTPPath(resolveHTTPPath(*pathFlag))); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context, server *ebitenmcp.Server, transport string, listenAddr string, httpPath string) error {
	switch transport {
	case "stdio":
		return server.ServeStdio(ctx, os.Stdin, os.Stdout)
	case "streamable-http":
		return serveStreamableHTTP(ctx, server, listenAddr, httpPath)
	default:
		return &invalidTransportError{transport: transport}
	}
}

func serveStreamableHTTP(ctx context.Context, server *ebitenmcp.Server, listenAddr string, httpPath string) error {
	httpServer := &http.Server{
		Addr:    listenAddr,
		Handler: newHTTPHandler(server, httpPath),
	}

	errCh := make(chan error, 1)
	go func() {
		err := httpServer.ListenAndServe()
		if err == nil || err == http.ErrServerClosed {
			errCh <- nil
			return
		}
		errCh <- err
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			return err
		}
		return <-errCh
	case err := <-errCh:
		return err
	}
}

func newHTTPHandler(server *ebitenmcp.Server, httpPath string) http.Handler {
	mux := http.NewServeMux()
	mux.Handle(httpPath, server.StreamableHTTPHandler(&mcp.StreamableHTTPOptions{
		JSONResponse: true,
	}))
	mux.HandleFunc("/healthz", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(map[string]any{
			"ok":        true,
			"transport": "streamable-http",
			"path":      httpPath,
		})
	})
	return mux
}

func resolveBridgeAddr(flagValue string) string {
	if flagValue != "" {
		return flagValue
	}
	if value := os.Getenv("EBITEN_DEBUG_ADDR"); value != "" {
		return value
	}
	return defaultBridgeAddr
}

func resolveTransport(flagValue string) string {
	if flagValue != "" {
		return flagValue
	}
	if value := os.Getenv("EBITEN_MCP_TRANSPORT"); value != "" {
		return value
	}
	return "stdio"
}

func resolveListenAddr(flagValue string) string {
	if flagValue != "" {
		return flagValue
	}
	if value := os.Getenv("EBITEN_MCP_LISTEN_ADDR"); value != "" {
		return value
	}
	return defaultHTTPListenAddr
}

func resolveHTTPPath(flagValue string) string {
	if flagValue != "" {
		return flagValue
	}
	if value := os.Getenv("EBITEN_MCP_HTTP_PATH"); value != "" {
		return value
	}
	return defaultHTTPPath
}

func normalizeHTTPPath(path string) string {
	path = strings.TrimSpace(path)
	if path == "" {
		return defaultHTTPPath
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	if len(path) > 1 {
		path = strings.TrimRight(path, "/")
	}
	return path
}

type invalidTransportError struct {
	transport string
}

func (err *invalidTransportError) Error() string {
	return "unsupported transport: " + err.transport
}
