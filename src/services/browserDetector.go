package services

import (
	"os"
	"path/filepath"
)

// BrowserInfo represents information about a detected browser
type BrowserInfo struct {
	Name string // Display name (e.g., "Google Chrome")
	Path string // Full path (e.g., "/Applications/Google Chrome.app")
}

// knownBrowser defines an app bundle name and its display name in the menu
type knownBrowser struct {
	AppName     string // e.g. "Google Chrome.app"
	DisplayName string // e.g. "Google Chrome"
}

// knownBrowsers is the single source for detection order and display-name mapping
var knownBrowsers = []knownBrowser{
	{"Safari.app", "Safari"},
	{"Google Chrome.app", "Google Chrome"},
	{"Firefox.app", "Firefox"},
	{"Microsoft Edge.app", "Microsoft Edge"},
	{"Edge.app", "Edge"},
	{"Brave Browser.app", "Brave"},
	{"Arc.app", "Arc"},
	{"Zen Browser.app", "Zen"},
	{"Island.app", "Island"},
	{"Opera.app", "Opera"},
	{"Vivaldi.app", "Vivaldi"},
	{"Orion.app", "Orion"},
	{"Chrome.app", "Chrome"},
	{"Chromium.app", "Chromium"},
	{"Firefox Developer Edition.app", "Firefox Developer"},
	{"Firefox Nightly.app", "Firefox Nightly"},
	{"Opera Developer.app", "Opera Developer"},
	{"Opera Beta.app", "Opera Beta"},
	{"Chrome Beta.app", "Chrome Beta"},
	{"Chrome Dev.app", "Chrome Dev"},
	{"Chrome Canary.app", "Chrome Canary"},
	{"Edge Beta.app", "Edge Beta"},
	{"Edge Dev.app", "Edge Dev"},
	{"Edge Canary.app", "Edge Canary"},
}

// BrowserDetector handles browser detection on macOS
type BrowserDetector struct{}

// NewBrowserDetector creates a new BrowserDetector instance
func NewBrowserDetector() *BrowserDetector {
	return &BrowserDetector{}
}

// DetectBrowsers scans common browser locations and returns a list of found browsers
func (bd *BrowserDetector) DetectBrowsers() []BrowserInfo {
	var browsers []BrowserInfo

	// Common browser locations on macOS
	searchPaths := []string{
		"/Applications",
		filepath.Join(os.Getenv("HOME"), "Applications"),
	}

	// Track found browsers to avoid duplicates
	foundPaths := make(map[string]bool)

	for _, searchPath := range searchPaths {
		for _, b := range knownBrowsers {
			browserPath := filepath.Join(searchPath, b.AppName)

			info, err := os.Stat(browserPath)
			if err != nil || !info.IsDir() {
				continue
			}

			normalizedPath := filepath.Clean(browserPath)
			if foundPaths[normalizedPath] {
				continue
			}
			foundPaths[normalizedPath] = true

			browsers = append(browsers, BrowserInfo{
				Name: bd.displayName(b.AppName, b.DisplayName),
				Path: normalizedPath,
			})
		}
	}

	return browsers
}

// displayName returns the display name for the app
func (bd *BrowserDetector) displayName(_, displayName string) string {
	return displayName
}
