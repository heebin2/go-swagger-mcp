package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/heebin2/go-swagger-mcp/internal/mcp"
)

func main() {

	SERVER := os.Getenv("SWAGGER_MCP_SERVER_URL")
	conn, err := mcp.NewServer("go-swagger-mcp", "v0.0.1", []string{
		SERVER,
	})
	if err != nil {
		panic(err)
	}

	if err := conn.Open(); err != nil {
		panic(fmt.Sprintf("failed to open MCP server connection: %v", err))
	}
	defer conn.Close()

	GracefullyShutdown()
}

func GracefullyShutdown() os.Signal {
	// Listen for termination signals
	signals := make(chan os.Signal, 1)
	defer close(signals)

	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	// Wait for the termination signal
	return <-signals
}
