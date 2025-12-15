package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage:")
		fmt.Println("  go run main.go list              - List all swagger specifications")
		fmt.Println("  go run main.go get <swagger-id>  - Get specific swagger specification")
		os.Exit(1)
	}

	command := os.Args[1]

	ctx := context.Background()

	// Create MCP client
	mcpClient := mcp.NewClient(&mcp.Implementation{
		Name:    "swagger-cli-client",
		Version: "1.0.0",
	}, nil)

	// Connect to MCP server
	mcpURL := "http://localhost:8080/mcp/sse"
	session, err := mcpClient.Connect(ctx, &mcp.SSEClientTransport{Endpoint: mcpURL}, nil)
	if err != nil {
		fmt.Printf("Failed to connect to MCP server: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		if err := session.Close(); err != nil {
			fmt.Printf("Warning: failed to close session: %v\n", err)
		}
	}()

	fmt.Println("✓ Connected to MCP server")

	switch command {
	case "list":
		listSwaggers(ctx, session)
	case "get":
		if len(os.Args) < 3 {
			fmt.Println("Error: swagger-id is required")
			fmt.Println("Usage: go run main.go get <swagger-id>")
			os.Exit(1)
		}
		swaggerID := os.Args[2]
		getSwagger(ctx, session, swaggerID)
	default:
		fmt.Printf("Unknown command: %s\n", command)
		fmt.Println("Available commands: list, get")
		os.Exit(1)
	}
}

func listSwaggers(ctx context.Context, session *mcp.ClientSession) {
	fmt.Println("\n📋 Listing swagger specifications...")

	result, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      "list_swagger",
		Arguments: map[string]interface{}{},
	})
	if err != nil {
		fmt.Printf("Failed to list swaggers: %v\n", err)
		os.Exit(1)
	}

	// Parse and display results
	if len(result.Content) > 0 {
		if textContent, ok := result.Content[0].(*mcp.TextContent); ok {
			var response map[string]interface{}
			if err := json.Unmarshal([]byte(textContent.Text), &response); err == nil {
				if swaggers, ok := response["swaggers"].([]interface{}); ok {
					fmt.Printf("\n✓ Found %d swagger specification(s):\n\n", len(swaggers))
					for i, swagger := range swaggers {
						if s, ok := swagger.(map[string]interface{}); ok {
							fmt.Printf("%d. ID: %s\n", i+1, s["id"])
							fmt.Printf("   Name: %s\n", s["name"])
							fmt.Printf("   Summary: %s\n", s["summary"])
							fmt.Println()
						}
					}
					return
				}
			}
			// Fallback: just print the text
			fmt.Println(textContent.Text)
		}
	}
}

func getSwagger(ctx context.Context, session *mcp.ClientSession, swaggerID string) {
	fmt.Printf("\n📄 Getting swagger specification: %s...\n", swaggerID)

	result, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name: "get_swagger",
		Arguments: map[string]interface{}{
			"id": swaggerID,
		},
	})
	if err != nil {
		fmt.Printf("Failed to get swagger: %v\n", err)
		os.Exit(1)
	}

	// Parse and display results
	if len(result.Content) > 0 {
		if textContent, ok := result.Content[0].(*mcp.TextContent); ok {
			var response map[string]interface{}
			if err := json.Unmarshal([]byte(textContent.Text), &response); err == nil {
				// Pretty print the JSON
				prettyJSON, _ := json.MarshalIndent(response, "", "  ")
				fmt.Println("\n✓ Swagger specification:")
				fmt.Println(string(prettyJSON))
				return
			}
			// Fallback: just print the text
			fmt.Println(textContent.Text)
		}
	}
}
