package main

import (
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"browserRedirectBar/src"
)

func main() {
	cleanup, err := src.InitLogger()
	if err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}
	defer cleanup()

	app, err := src.NewApp()
	if err != nil {
		log.Fatal("Failed to initialize app:", err)
	}

	// When launched as default browser, macOS passes the URL as the first argument
	for _, arg := range os.Args[1:] {
		if arg == "" {
			continue
		}
		var urlStr string
		if strings.HasPrefix(arg, "http://") || strings.HasPrefix(arg, "https://") || strings.HasPrefix(arg, "file://") {
			urlStr = arg
		} else if strings.HasSuffix(strings.ToLower(arg), ".html") || strings.HasSuffix(strings.ToLower(arg), ".htm") || strings.HasSuffix(strings.ToLower(arg), ".xhtml") {
			absPath, err := filepath.Abs(arg)
			if err != nil {
				continue
			}
			urlStr = (&url.URL{Scheme: "file", Path: absPath}).String()
		}
		if urlStr != "" {
			select {
			case app.URLChan() <- urlStr:
			default:
				// Channel full, drop URL
			}
		}
	}

	app.Run()
}
