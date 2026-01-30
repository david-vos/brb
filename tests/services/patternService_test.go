package services

import (
	"browserRedirectBar/src/services"
	"testing"
)

func TestPatternService_FindBrowserForURL_SimplePatterns(t *testing.T) {
	// Setup test config
	testConfig := services.Config{
		Browsers: []services.BrowserConfig{
			{
				Patterns:   []string{"github.com", "gitlab.com"},
				BrowserURL: "/Applications/Chrome.app",
			},
			{
				Patterns:   []string{"localhost"},
				BrowserURL: "/Applications/Firefox.app",
			},
		},
		DefaultBrowserURL: "/Applications/Safari.app",
	}

	service := services.NewPatternService(testConfig)

	tests := []struct {
		name     string
		url      string
		expected string
	}{
		{
			name:     "GitHub URL matches Chrome",
			url:      "https://github.com/user/repo",
			expected: "/Applications/Chrome.app",
		},
		{
			name:     "GitLab URL matches Chrome",
			url:      "https://gitlab.com/project",
			expected: "/Applications/Chrome.app",
		},
		{
			name:     "Localhost URL matches Firefox",
			url:      "http://localhost:3000",
			expected: "/Applications/Firefox.app",
		},
		{
			name:     "No match returns empty",
			url:      "https://example.com",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.FindBrowserForURL(tt.url)
			if result != tt.expected {
				t.Errorf("FindBrowserForURL(%q) = %q, want %q", tt.url, result, tt.expected)
			}
		})
	}
}

func TestPatternService_FindBrowserForURL_RegexPatterns(t *testing.T) {
	// Setup test config with regex patterns
	testConfig := services.Config{
		Browsers: []services.BrowserConfig{
			{
				RegexPatterns: []string{
					"^https://.*\\.internal\\..*",
					"^https://.*\\.dev\\.local",
				},
				BrowserURL: "/Applications/Firefox.app",
			},
			{
				RegexPatterns: []string{
					".*\\.staging\\..*",
				},
				BrowserURL: "/Applications/Chrome.app",
			},
		},
		DefaultBrowserURL: "/Applications/Safari.app",
	}

	service := services.NewPatternService(testConfig)

	tests := []struct {
		name     string
		url      string
		expected string
	}{
		{
			name:     "Internal domain matches Firefox",
			url:      "https://app.internal.company.com",
			expected: "/Applications/Firefox.app",
		},
		{
			name:     "Dev local matches Firefox",
			url:      "https://api.dev.local",
			expected: "/Applications/Firefox.app",
		},
		{
			name:     "Staging domain matches Chrome",
			url:      "https://app.staging.example.com",
			expected: "/Applications/Chrome.app",
		},
		{
			name:     "No regex match returns empty",
			url:      "https://example.com",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.FindBrowserForURL(tt.url)
			if result != tt.expected {
				t.Errorf("FindBrowserForURL(%q) = %q, want %q", tt.url, result, tt.expected)
			}
		})
	}
}

func TestPatternService_FindBrowserForURL_MixedPatterns(t *testing.T) {
	// Setup test config with both simple and regex patterns
	testConfig := services.Config{
		Browsers: []services.BrowserConfig{
			{
				Patterns:      []string{"github.com"},
				RegexPatterns: []string{"^https://.*\\.internal\\..*"},
				BrowserURL:    "/Applications/Chrome.app",
			},
		},
		DefaultBrowserURL: "/Applications/Safari.app",
	}

	service := services.NewPatternService(testConfig)

	tests := []struct {
		name     string
		url      string
		expected string
	}{
		{
			name:     "Simple pattern matches first",
			url:      "https://github.com/user/repo",
			expected: "/Applications/Chrome.app",
		},
		{
			name:     "Regex pattern also matches",
			url:      "https://app.internal.company.com",
			expected: "/Applications/Chrome.app",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.FindBrowserForURL(tt.url)
			if result != tt.expected {
				t.Errorf("FindBrowserForURL(%q) = %q, want %q", tt.url, result, tt.expected)
			}
		})
	}
}

func TestPatternService_InvalidRegexPattern(t *testing.T) {
	// Setup test config with invalid regex
	testConfig := services.Config{
		Browsers: []services.BrowserConfig{
			{
				RegexPatterns: []string{
					"[invalid regex[", // Invalid regex pattern
				},
				BrowserURL: "/Applications/Firefox.app",
			},
		},
		DefaultBrowserURL: "/Applications/Safari.app",
	}

	service := services.NewPatternService(testConfig)

	// This should not panic and should return empty (no match)
	result := service.FindBrowserForURL("https://example.com")
	if result != "" {
		t.Errorf("Invalid regex should not match, got %q", result)
	}
}
