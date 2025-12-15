package mcp

import (
	"context"
	"fmt"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func middleware(next mcp.MethodHandler) mcp.MethodHandler {
	return func(
		ctx context.Context,
		method string,
		req mcp.Request,
	) (mcp.Result, error) {
		fmt.Println("MCP method started",
			"method", method,
			"session_id", req.GetSession().ID(),
			"has_params", req.GetParams() != nil,
		)
		// Log more for tool calls.
		if ctr, ok := req.(*mcp.CallToolRequest); ok {
			fmt.Println("Calling tool",
				"name", ctr.Params.Name,
				"args", ctr.Params.Arguments)
		}

		start := time.Now()
		result, err := next(ctx, method, req)
		duration := time.Since(start)
		if err != nil {
			fmt.Println("MCP method failed",
				"method", method,
				"session_id", req.GetSession().ID(),
				"duration_ms", duration.Milliseconds(),
				"err", err,
			)
		} else {
			fmt.Println("MCP method completed",
				"method", method,
				"session_id", req.GetSession().ID(),
				"duration_ms", duration.Milliseconds(),
				"has_result", result != nil,
			)
			// Log more for tool results.
			if ctr, ok := result.(*mcp.CallToolResult); ok {
				fmt.Println("tool result",
					"isError", ctr.IsError,
					"structuredContent", ctr.StructuredContent)
			}
		}
		return result, err
	}
}
