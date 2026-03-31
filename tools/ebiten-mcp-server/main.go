package main

import (
	"context"
	"flag"
	"log"
	"os"

	ebitenmcp "github.com/kimyechan/ebiten-aio-framework/libs/go/ebiten-mcp"
)

const defaultBridgeAddr = "127.0.0.1:47831"

func main() {
	addrFlag := flag.String("addr", "", "ebiten debug bridge address")
	flag.Parse()

	server := ebitenmcp.New(ebitenmcp.NewBridgeClient(resolveBridgeAddr(*addrFlag)))
	if err := server.ServeStdio(context.Background(), os.Stdin, os.Stdout); err != nil {
		log.Fatal(err)
	}
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
