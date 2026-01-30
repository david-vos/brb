#include "_cgo_export.h"

// C bridge function that Objective-C can call
// This calls the Go-exported function
void bridgeSendURLToGo(const char* url) {
    sendURLToGo((char*)url);
}
