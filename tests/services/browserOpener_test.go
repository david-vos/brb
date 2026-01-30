package services

import (
	"browserRedirectBar/src/services"
	"testing"

	"github.com/stretchr/testify/mock"
)

// MockBrowserOpener is mock for BrowserOpener
type MockBrowserOpener struct {
	mock.Mock
}

// OpenBrowser implements BrowserOpener interface
func (m *MockBrowserOpener) OpenBrowser(browserPath string, url string) {
	m.Called(browserPath, url)
}

func TestBrowserService_OpenBrowser(t *testing.T) {
	mockOpener := new(MockBrowserOpener)
	service := services.NewBrowserServiceWithOpener(mockOpener)

	testBrowserPath := "/Applications/Chrome.app"
	testURL := "https://example.com"

	mockOpener.On("OpenBrowser", testBrowserPath, testURL).Return()

	service.OpenBrowser(testBrowserPath, testURL)

	mockOpener.AssertExpectations(t)
}

func TestBrowserService_OpenBrowser_NilOpener(t *testing.T) {
	service := &services.BrowserService{}
	// Access the private field through reflection or test the public API
	// For now, test that it doesn't panic
	service.OpenBrowser("/Applications/Chrome.app", "https://example.com")
}
