#import <Foundation/Foundation.h>
#import <AppKit/AppKit.h>
#import <Carbon/Carbon.h>

// Forward declaration of the bridge function
extern void bridgeSendURLToGo(const char* url);

// Dedicated handler for GetURL Apple Events (not the app delegate)
@interface GetURLHandler : NSObject
@end

@implementation GetURLHandler

- (void)handleGetURL:(NSAppleEventDescriptor *)event withReplyEvent:(NSAppleEventDescriptor *)replyEvent {
    NSAppleEventDescriptor *urlDescriptor = [event paramDescriptorForKeyword:keyDirectObject];
    if (urlDescriptor) {
        NSString *urlString = [urlDescriptor stringValue];
        if (urlString) {
            const char *urlCStr = [urlString UTF8String];
            if (urlCStr) {
                bridgeSendURLToGo(urlCStr);
            }
        }
    }
}

@end

static GetURLHandler *getURLHandler = nil;

void setupAppleEventHandler() {
    @autoreleasepool {
        if (getURLHandler == nil) {
            getURLHandler = [[GetURLHandler alloc] init];
            // Register directly with NSAppleEventManager - do NOT use app delegate.
            // This way we don't conflict with systray and the handler is registered immediately.
            [[NSAppleEventManager sharedAppleEventManager]
                setEventHandler:getURLHandler
                andSelector:@selector(handleGetURL:withReplyEvent:)
                forEventClass:kInternetEventClass
                andEventID:kAEGetURL];
        }
    }
}
