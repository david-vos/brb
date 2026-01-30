package services

import (
	"log"
	"os/exec"
)

// BrowserOpener defines the interface for opening browsers
type BrowserOpener interface {
	OpenBrowser(browserPath string, url string)
}

// RealBrowserOpener is the production implementation that actually opens browsers
type RealBrowserOpener struct{}

// NewRealBrowserOpener creates a new RealBrowserOpener
func NewRealBrowserOpener() *RealBrowserOpener {
	return &RealBrowserOpener{}
}

// OpenBrowser opens a URL in the specified browser
func (r *RealBrowserOpener) OpenBrowser(browserPath string, url string) {
	cmd := exec.Command("open", "-a", browserPath, url)
	err := cmd.Run()
	if err == nil {
		return
	}
	if fallbackErr := exec.Command("open", url).Run(); fallbackErr != nil {
		log.Printf("Failed to open browser %s: %v; fallback failed: %v", browserPath, err, fallbackErr)
	}
}
