package services

import "log"

// BrowserService handles browser operations
type BrowserService struct {
	opener BrowserOpener
}

// NewBrowserService creates a new BrowserService instance with a real browser opener
func NewBrowserService() *BrowserService {
	return &BrowserService{
		opener: NewRealBrowserOpener(),
	}
}

// NewBrowserServiceWithOpener creates a new BrowserService with a custom opener (for testing)
func NewBrowserServiceWithOpener(opener BrowserOpener) *BrowserService {
	return &BrowserService{
		opener: opener,
	}
}

// OpenBrowser opens a URL in the specified browser
func (bs *BrowserService) OpenBrowser(browserPath string, url string) {
	if bs.opener == nil {
		log.Printf("BrowserOpener is nil")
		return
	}
	bs.opener.OpenBrowser(browserPath, url)
}
