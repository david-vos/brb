package services

// BrowserConfig represents a browser configuration with URL patterns
type BrowserConfig struct {
	Patterns      []string `json:"patterns"`      // Simple string matching (case-insensitive)
	RegexPatterns []string `json:"regexPatterns"` // Regex pattern matching
	BrowserURL    string   `json:"browserURL"`    // Path to browser application (e.g., "/Applications/Google Chrome.app")
}

// Config represents the application configuration
type Config struct {
	Browsers          []BrowserConfig `json:"browsers"`
	DefaultBrowserURL string          `json:"defaultBrowserURL"` // Path to default browser application
}
