package services

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation -framework AppKit
#import <Foundation/Foundation.h>
#import <AppKit/AppKit.h>

int openSystemSettingsDefaultBrowser() {
	NSURL *url = [NSURL URLWithString:@"x-apple.systempreferences:com.apple.Desktop-Settings"];
	if (url == nil) {
		return -1;
	}

	NSWorkspace *workspace = [NSWorkspace sharedWorkspace];
	BOOL success = [workspace openURL:url];
	return success ? 0 : -1;
}
*/
import "C"

import "log"

// DefaultBrowserService handles requesting default browser status
type DefaultBrowserService struct{}

// NewDefaultBrowserService creates a new DefaultBrowserService instance
func NewDefaultBrowserService() *DefaultBrowserService {
	return &DefaultBrowserService{}
}

// RequestDefaultBrowser opens System Settings to the default web browser pane
func (dbs *DefaultBrowserService) RequestDefaultBrowser() error {
	result := C.openSystemSettingsDefaultBrowser()
	if result != 0 {
		log.Printf("Could not open System Settings: %d", result)
	}
	return nil
}
