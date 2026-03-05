package models

// SiteGenerationRequest represents a request to generate a Site class
type SiteGenerationRequest struct {
	Name        string   `json:"name" validate:"required,min=1"`
	BaseURL     string   `json:"baseUrl" validate:"required,url"`
	Pages       []string `json:"pages" validate:"required,min=1"`
	Flows       []string `json:"flows,omitempty"`
	Description string   `json:"description,omitempty"`
}

// PageGenerationRequest represents a request to generate a Page class
type PageGenerationRequest struct {
	Name        string            `json:"name" validate:"required,min=1"`
	Path        string            `json:"path" validate:"required"`
	Elements    []ElementDefinition `json:"elements" validate:"required,min=1,dive"`
	Description string            `json:"description,omitempty"`
}

// ElementDefinition represents an element on a page
type ElementDefinition struct {
	Name     string `json:"name" validate:"required,min=1"`
	Selector string `json:"selector" validate:"required,min=1"`
	Type     string `json:"type,omitempty" validate:"omitempty,oneof=button input link select checkbox radio"`
}


// FlowGenerationRequest represents a request to generate a Flow class
type FlowGenerationRequest struct {
	Name        string   `json:"name" validate:"required,min=1"`
	Steps       []string `json:"steps" validate:"required,min=1"`
	Description string   `json:"description,omitempty"`
}

// StepDefinitionRequest represents a request to generate Cucumber step definitions
type StepDefinitionRequest struct {
	Steps       []StepDefinition `json:"steps" validate:"required,min=1,dive"`
	Language    string           `json:"language,omitempty" validate:"omitempty,oneof=typescript javascript"`
	Description string           `json:"description,omitempty"`
}

// StepDefinition represents a single Cucumber step
type StepDefinition struct {
	Type       string `json:"type" validate:"required,oneof=Given When Then"`
	Pattern    string `json:"pattern" validate:"required,min=1"`
	Parameters string `json:"parameters,omitempty"`
}

// ElementGenerationRequest represents a request to generate a custom Element class
type ElementGenerationRequest struct {
	Name         string   `json:"name" validate:"required,min=1"`
	RootSelector string   `json:"rootSelector" validate:"required"`
	Methods      []string `json:"methods,omitempty"`
	Description  string   `json:"description,omitempty"`
}

// ServiceGenerationRequest represents a request to generate an API service class
type ServiceGenerationRequest struct {
	Name        string              `json:"name" validate:"required,min=1"`
	BaseURL     string              `json:"baseUrl" validate:"required,url"`
	Endpoints   []EndpointDefinition `json:"endpoints" validate:"required,min=1,dive"`
	Description string              `json:"description,omitempty"`
}

// EndpointDefinition represents an API endpoint
type EndpointDefinition struct {
	Name   string `json:"name" validate:"required,min=1"`
	Method string `json:"method" validate:"required,oneof=GET POST PUT PATCH DELETE"`
	Path   string `json:"path" validate:"required"`
}

// AnalysisRequest represents a request to analyze framework structure
type AnalysisRequest struct {
	ProjectPath string   `json:"projectPath" validate:"required"`
	Patterns    []string `json:"patterns,omitempty"`
}

// ValidationRequest represents a request to validate code patterns
type ValidationRequest struct {
	FilePath string   `json:"filePath" validate:"required"`
	Rules    []string `json:"rules,omitempty"`
}

// MigrationRequest represents a request to migrate code to the framework
type MigrationRequest struct {
	SourceCode string `json:"sourceCode" validate:"required"`
	Type       string `json:"type" validate:"required,oneof=page test"`
}

// DocumentationRequest represents a request to generate project documentation
type DocumentationRequest struct {
	ProjectPath string `json:"projectPath" validate:"required"`
	OutputPath  string `json:"outputPath" validate:"required"`
}

// TestExecutionRequest represents a request to execute tests
type TestExecutionRequest struct {
	ProjectPath string `json:"projectPath" validate:"required"`
	Command     string `json:"command,omitempty"`     // e.g. "npx playwright test"
	Filter      string `json:"filter,omitempty"`      // e.g. "login" (grep)
	Env         map[string]string `json:"env,omitempty"`
}
