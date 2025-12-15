package mcp

type ListSwaggerInput struct {
}

type ListSwaggersOutput struct {
	Swaggers []SwaggerRegistryItem `json:"swaggers" jsonschema:"list of available swagger specifications"`
}

type GetSwaggerOutput struct {
	ID      string         `json:"id"`
	Name    string         `json:"name"`
	Spec    map[string]any `json:"spec" jsonschema:"full swagger/openapi specification"`
	Summary string         `json:"summary"`
}

type GetSwaggerInput struct {
	ID string `json:"id" jsonschema:"swagger specification ID to retrieve"`
}
