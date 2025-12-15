package mcp

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/heebin2/go-swagger-mcp/internal/swagger"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type SwaggerRegistryItem struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Version     string `json:"version"`
}

type Server struct {
	specs map[string]*swagger.Spec
	conn  *mcp.Server
	route *http.Server
	wg    sync.WaitGroup
}

func NewServer(name, version string, urls []string) (*Server, error) {

	// add swagger specs
	m := make(map[string]*swagger.Spec, len(urls))
	for _, url := range urls {
		spec, err := swagger.NewParser().Load(url)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to load %s", url)
		}

		m[url] = spec
	}

	return &Server{
		specs: m,
		conn: mcp.NewServer(&mcp.Implementation{
			Name:    name,
			Title:   name,
			Version: version,
		}, &mcp.ServerOptions{HasTools: true}),
	}, nil
}

func (s *Server) listSwagger(ctx context.Context, _ *mcp.CallToolRequest, input ListSwaggerInput) (
	*mcp.CallToolResult, ListSwaggersOutput, error) {

	select {
	case <-ctx.Done():
		return nil, ListSwaggersOutput{}, ctx.Err()
	default:
	}

	var swaggers []SwaggerRegistryItem
	for id, spec := range s.specs {
		swaggers = append(swaggers, SwaggerRegistryItem{
			ID:          id,
			Name:        spec.Info.Title,
			Description: spec.Info.Description,
			Version:     spec.Info.Version,
		})
	}

	return nil, ListSwaggersOutput{
		Swaggers: swaggers,
	}, nil
}

func (s *Server) getSwagger(ctx context.Context, _ *mcp.CallToolRequest, input GetSwaggerInput) (
	*mcp.CallToolResult, GetSwaggerOutput, error) {

	select {
	case <-ctx.Done():
		return nil, GetSwaggerOutput{}, ctx.Err()
	default:
	}

	spec, exists := s.specs[input.ID]
	if !exists {
		return nil, GetSwaggerOutput{}, errors.Errorf("swagger spec %s not found", input.ID)
	}

	// Convert spec to map for output
	specMap := map[string]any{
		"info":  spec.Info,
		"paths": spec.Paths,
	}

	if spec.Servers != nil {
		specMap["servers"] = spec.Servers
	}

	return nil, GetSwaggerOutput{
		ID:      input.ID,
		Name:    spec.Info.Title,
		Spec:    specMap,
		Summary: spec.Summary(),
	}, nil
}

// loggingMiddleware logs HTTP requests with method, path, duration, and status
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a custom ResponseWriter to capture status code
		lrw := &loggingResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		// Log request
		log.Printf("[MCP] --> %s %s %s", r.Method, r.URL.Path, r.RemoteAddr)

		// Call the next handler
		next.ServeHTTP(lrw, r)

		// Log response with duration
		duration := time.Since(start)
		log.Printf("[MCP] <-- %s %s %d %v", r.Method, r.URL.Path, lrw.statusCode, duration)
	})
}

// loggingResponseWriter wraps http.ResponseWriter to capture status code
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (s *Server) Open() error {
	mcp.AddTool(s.conn, &mcp.Tool{
		Name:        "list_swagger",
		Description: "List all available Swagger/OpenAPI specifications in the registry",
	}, s.listSwagger)

	mcp.AddTool(s.conn, &mcp.Tool{
		Name:        "get_swagger",
		Description: "Get a Swagger/OpenAPI specification by its ID from the registry",
	}, s.getSwagger)

	handler := mcp.NewSSEHandler(func(request *http.Request) *mcp.Server {
		log.Printf("[MCP] SSE connection request from: %s", request.RemoteAddr)
		if request.URL.String() == "/mcp/sse" {
			return s.conn
		}

		return s.conn
	}, nil)

	s.conn.AddReceivingMiddleware(middleware)

	s.route = &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		if err := s.route.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Println("HTTP server error:", err)
		}
	}()

	return nil
}

func (s *Server) Close() {
	if err := s.route.Close(); err != nil {
		fmt.Println("failed to close HTTP server:", err)
	}

	s.wg.Wait()
}
