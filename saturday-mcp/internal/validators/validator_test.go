package validators

import (
	"testing"

	"github.com/orieken/saturday-mcp/internal/models"
)

func TestValidator(t *testing.T) {
	validator := NewValidator()

	t.Run("Valid PageGenerationRequest", func(t *testing.T) {
		req := models.PageGenerationRequest{
			Name: "login",
			Path: "/login",
			Elements: []models.ElementDefinition{
				{Name: "username", Selector: "#username"},
				{Name: "password", Selector: "#password"},
			},
		}

		err := validator.Validate(req)
		if err != nil {
			t.Errorf("Expected valid request, got error: %v", err)
		}
	})

	t.Run("Invalid PageGenerationRequest - missing name", func(t *testing.T) {
		req := models.PageGenerationRequest{
			Path: "/login",
			Elements: []models.ElementDefinition{
				{Name: "username", Selector: "#username"},
			},
		}

		err := validator.Validate(req)
		if err == nil {
			t.Error("Expected validation error for missing name")
		}
	})

	t.Run("Invalid PageGenerationRequest - empty elements", func(t *testing.T) {
		req := models.PageGenerationRequest{
			Name:     "login",
			Path:     "/login",
			Elements: []models.ElementDefinition{},
		}

		err := validator.Validate(req)
		if err == nil {
			t.Error("Expected validation error for empty elements")
		}
	})

	t.Run("Valid SiteGenerationRequest", func(t *testing.T) {
		req := models.SiteGenerationRequest{
			Name:    "ecommerce",
			BaseURL: "https://example.com",
			Pages:   []string{"home", "product"},
		}

		err := validator.Validate(req)
		if err != nil {
			t.Errorf("Expected valid request, got error: %v", err)
		}
	})

	t.Run("Invalid SiteGenerationRequest - invalid URL", func(t *testing.T) {
		req := models.SiteGenerationRequest{
			Name:    "ecommerce",
			BaseURL: "not-a-url",
			Pages:   []string{"home"},
		}

		err := validator.Validate(req)
		if err == nil {
			t.Error("Expected validation error for invalid URL")
		}
	})

	t.Run("Invalid SiteGenerationRequest - empty pages", func(t *testing.T) {
		req := models.SiteGenerationRequest{
			Name:    "ecommerce",
			BaseURL: "https://example.com",
			Pages:   []string{},
		}

		err := validator.Validate(req)
		if err == nil {
			t.Error("Expected validation error for empty pages")
		}
	})

	t.Run("Valid ElementDefinition with type", func(t *testing.T) {
		elem := models.ElementDefinition{
			Name:     "submitBtn",
			Selector: "button[type='submit']",
			Type:     "button",
		}

		err := validator.Validate(elem)
		if err != nil {
			t.Errorf("Expected valid element, got error: %v", err)
		}
	})

	t.Run("Invalid ElementDefinition - invalid type", func(t *testing.T) {
		elem := models.ElementDefinition{
			Name:     "elem",
			Selector: "#elem",
			Type:     "invalid",
		}

		err := validator.Validate(elem)
		if err == nil {
			t.Error("Expected validation error for invalid type")
		}
	})

	t.Run("Valid FlowGenerationRequest", func(t *testing.T) {
		req := models.FlowGenerationRequest{
			Name:  "checkout",
			Steps: []string{"addToCart", "enterShipping", "payment"},
		}

		err := validator.Validate(req)
		if err != nil {
			t.Errorf("Expected valid request, got error: %v", err)
		}
	})

	t.Run("Valid ServiceGenerationRequest", func(t *testing.T) {
		req := models.ServiceGenerationRequest{
			Name:    "userService",
			BaseURL: "https://api.example.com",
			Endpoints: []models.EndpointDefinition{
				{Name: "getUser", Method: "GET", Path: "/users/:id"},
				{Name: "createUser", Method: "POST", Path: "/users"},
			},
		}

		err := validator.Validate(req)
		if err != nil {
			t.Errorf("Expected valid request, got error: %v", err)
		}
	})

	t.Run("Invalid EndpointDefinition - invalid method", func(t *testing.T) {
		endpoint := models.EndpointDefinition{
			Name:   "test",
			Method: "INVALID",
			Path:   "/test",
		}

		err := validator.Validate(endpoint)
		if err == nil {
			t.Error("Expected validation error for invalid HTTP method")
		}
	})
}

func TestValidationError(t *testing.T) {
	validator := NewValidator()

	req := models.PageGenerationRequest{
		// Missing required fields
		Elements: []models.ElementDefinition{},
	}

	err := validator.Validate(req)
	if err == nil {
		t.Fatal("Expected validation error")
	}

	valErr, ok := err.(*ValidationError)
	if !ok {
		t.Fatal("Expected ValidationError type")
	}

	errors := valErr.GetErrors()
	if len(errors) == 0 {
		t.Error("Expected validation errors map to have entries")
	}

	errorMsg := valErr.Error()
	if len(errorMsg) == 0 {
		t.Error("Expected error message to be non-empty")
	}
}

func TestValidateField(t *testing.T) {
	validator := NewValidator()

	t.Run("Valid URL", func(t *testing.T) {
		err := validator.ValidateField("https://example.com", "url")
		if err != nil {
			t.Errorf("Expected valid URL, got error: %v", err)
		}
	})

	t.Run("Invalid URL", func(t *testing.T) {
		err := validator.ValidateField("not-a-url", "url")
		if err == nil {
			t.Error("Expected validation error for invalid URL")
		}
	})

	t.Run("Required field", func(t *testing.T) {
		err := validator.ValidateField("", "required")
		if err == nil {
			t.Error("Expected validation error for empty required field")
		}
	})
}
