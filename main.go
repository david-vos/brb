package main

import (
	"log"
	"os"
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
		if strings.HasPrefix(arg, "http://") || strings.HasPrefix(arg, "https://") {
			select {
			case app.URLChan() <- arg:
			default:
				// Channel full, drop URL
			}
		}
	}

	app.Run()
}
