package services

/*
#cgo CFLAGS: -x objective-c -I${SRCDIR}
#cgo LDFLAGS: -framework Foundation -framework AppKit -framework Carbon
// Native implementation is in nativeHelpers/; compiled via appleEventNative.m
extern void setupAppleEventHandler();
extern void bridgeSendURLToGo(const char* url);
*/
import "C"

import (
	"log"
)

var urlChannel chan string

//export sendURLToGo
func sendURLToGo(urlCStr *C.char) {
	if urlCStr == nil {
		return
	}
	url := C.GoString(urlCStr)
	if urlChannel == nil {
		return
	}
	select {
	case urlChannel <- url:
	default:
		log.Printf("Apple Event: URL channel full, dropping URL")
	}
}

// SetupAppleEventHandler sets up the Apple Event handler to receive URLs when the app is already running
func SetupAppleEventHandler(urlChan chan string) {
	urlChannel = urlChan
	C.setupAppleEventHandler()
}
