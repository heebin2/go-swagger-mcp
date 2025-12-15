package swagger

// Spec represents a simplified OpenAPI/Swagger specification
type Spec struct {
	OpenAPI string                 `json:"openapi,omitempty"`
	Swagger string                 `json:"swagger,omitempty"`
	Info    Info                   `json:"info"`
	Paths   map[string]PathItem    `json:"paths"`
	Servers []Server               `json:"servers,omitempty"`
	Raw     map[string]interface{} `json:"-"`
}

func (s Spec) Summary() string {
	summary := ""
	if s.Info.Title != "" {
		summary = s.Info.Title
		if s.Info.Description != "" {
			summary += " - " + s.Info.Description
		}
		if s.Info.Version != "" {
			summary += " (v" + s.Info.Version + ")"
		}
		summary += " with " + string(rune(len(s.Paths))) + " endpoints"
	}

	return summary
}

type Info struct {
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Version     string `json:"version"`
}

type Server struct {
	URL         string `json:"url"`
	Description string `json:"description,omitempty"`
}

type PathItem struct {
	Get     *Operation `json:"get,omitempty"`
	Post    *Operation `json:"post,omitempty"`
	Put     *Operation `json:"put,omitempty"`
	Delete  *Operation `json:"delete,omitempty"`
	Patch   *Operation `json:"patch,omitempty"`
	Options *Operation `json:"options,omitempty"`
	Head    *Operation `json:"head,omitempty"`
}

type Operation struct {
	Summary     string              `json:"summary,omitempty"`
	Description string              `json:"description,omitempty"`
	OperationID string              `json:"operationId,omitempty"`
	Tags        []string            `json:"tags,omitempty"`
	Parameters  []Parameter         `json:"parameters,omitempty"`
	RequestBody *RequestBody        `json:"requestBody,omitempty"`
	Responses   map[string]Response `json:"responses,omitempty"`
}

type Parameter struct {
	Name        string      `json:"name"`
	In          string      `json:"in"`
	Description string      `json:"description,omitempty"`
	Required    bool        `json:"required,omitempty"`
	Schema      interface{} `json:"schema,omitempty"`
}

type RequestBody struct {
	Description string               `json:"description,omitempty"`
	Required    bool                 `json:"required,omitempty"`
	Content     map[string]MediaType `json:"content,omitempty"`
}

type MediaType struct {
	Schema interface{} `json:"schema,omitempty"`
}

type Response struct {
	Description string               `json:"description"`
	Content     map[string]MediaType `json:"content,omitempty"`
}

type APIRecommendation struct {
	Path        string   `json:"path"`
	Method      string   `json:"method"`
	Summary     string   `json:"summary"`
	Description string   `json:"description"`
	Score       float64  `json:"score"`
	Reasons     []string `json:"reasons"`
}
