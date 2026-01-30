package tests

import (
	"browserRedirectBar/src/services"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockBrowserOpener is a testify mock for BrowserOpener
type MockBrowserOpener struct {
	mock.Mock
}

// OpenBrowser implements BrowserOpener interface
func (m *MockBrowserOpener) OpenBrowser(browserPath string, url string) {
	m.Called(browserPath, url)
}

func TestApp_HandleURL_DefaultBrowserFallback(t *testing.T) {
	// Setup test config with no matching patterns
	testConfig := services.Config{
		Browsers: []services.BrowserConfig{
			{
				Patterns:   []string{"github.com"},
				BrowserURL: "/Applications/Chrome.app",
			},
		},
		DefaultBrowserURL: "/Applications/Safari.app",
	}

	// Create services
	configService := &services.ConfigService{}
	configService.SetConfig(testConfig)
	patternService := services.NewPatternService(testConfig)

	// Create mock browser opener using testify
	mockOpener := new(MockBrowserOpener)
	browserService := services.NewBrowserServiceWithOpener(mockOpener)

	// Test URL that doesn't match any pattern
	testURL := "https://example.com"

	// Set up mock expectation
	mockOpener.On("OpenBrowser", testConfig.DefaultBrowserURL, testURL).Return()

	// Test through services directly (integration test)
	// Find browser path
	browserPath := patternService.FindBrowserForURL(testURL)
	if browserPath == "" {
		browserPath = testConfig.DefaultBrowserURL
	}

	// Open browser
	browserService.OpenBrowser(browserPath, testURL)

	// Verify mock was called correctly
	mockOpener.AssertExpectations(t)

	// Verify that patternService returns empty for non-matching URLs
	browser := patternService.FindBrowserForURL(testURL)
	assert.Empty(t, browser, "Expected no match for non-matching URL")
}
