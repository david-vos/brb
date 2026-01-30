// Wrapper so cgo compiles the native code in nativeHelpers/
// (cgo only compiles C/ObjC files in the same directory as the .go file)
#import "nativeHelpers/appleEventBridge.c"
#import "nativeHelpers/appleEventHandler.m"
