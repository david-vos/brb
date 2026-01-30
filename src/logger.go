package src

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

// InitLogger initializes logging to both stderr and a log file at ~/.brb/brb.log
func InitLogger() (func(), error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	logDir := filepath.Join(homeDir, ".brb")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}
	logPath := filepath.Join(logDir, "brb.log")
	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	multi := io.MultiWriter(os.Stderr, f)
	log.SetOutput(multi)
	log.SetFlags(log.Ldate | log.Ltime)
	cleanup := func() {
		err := f.Close()
		if err != nil {
			return
		}
	}
	return cleanup, nil
}
