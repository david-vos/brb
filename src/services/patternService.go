package services

import (
	"regexp"
	"strings"
	"sync"
)

// PatternService handles URL pattern matching
type PatternService struct {
	config         Config
	regexCache     map[string]*regexp.Regexp
	regexCacheLock sync.RWMutex
}

// NewPatternService creates a new PatternService instance
func NewPatternService(config Config) *PatternService {
	return &PatternService{
		config:     config,
		regexCache: make(map[string]*regexp.Regexp),
	}
}

// UpdateConfig updates the configuration used for pattern matching
func (ps *PatternService) UpdateConfig(config Config) {
	ps.config = config
}

// FindBrowserForURL finds the appropriate browser for a given URL based on patterns
func (ps *PatternService) FindBrowserForURL(url string) string {
	urlLower := strings.ToLower(url)

	for _, browserConfig := range ps.config.Browsers {
		for _, pattern := range browserConfig.Patterns {
			if strings.Contains(urlLower, strings.ToLower(pattern)) {
				return browserConfig.BrowserURL
			}
		}
		for _, regexPattern := range browserConfig.RegexPatterns {
			compiled, err := ps.getCompiledRegex(regexPattern)
			if err != nil {
				continue
			}
			if compiled.MatchString(url) {
				return browserConfig.BrowserURL
			}
		}
	}
	return ""
}

// getCompiledRegex returns a compiled regex, using cache for performance
func (ps *PatternService) getCompiledRegex(pattern string) (*regexp.Regexp, error) {
	ps.regexCacheLock.RLock()
	if cached, exists := ps.regexCache[pattern]; exists {
		ps.regexCacheLock.RUnlock()
		return cached, nil
	}
	ps.regexCacheLock.RUnlock()

	ps.regexCacheLock.Lock()
	defer ps.regexCacheLock.Unlock()

	if cached, exists := ps.regexCache[pattern]; exists {
		return cached, nil
	}
	compiled, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	ps.regexCache[pattern] = compiled
	return compiled, nil
}
