# Browser Redirect Bar (brb)
Don't you hate it when your something or someone requires you to use a specific browser when accessing a specific site? 
I fixed it by creating a "default browser" that redirects you to the preferred browser for that site.

### Quick Start
Simply run to get started:
```bash
go run build.go
```

The app bundle can then be set as your default browser through System Settings.

This will automatically:
- Build the app bundle
- Install it to `~/Applications/`
- Register it with Launch Services

### Development Build
For development/testing without creating an app bundle:

```bash
go build -o brb main.go
./brb https://github.com  # Test with a URL
```

## Configuration

The configuration file is located at `~/.brb/config.json`. It will be created automatically with default settings on first run.

### Example config.json

Here's a complete example configuration file:

```json
{
  "browsers": [
    {
      "patterns": ["github.com", "gitlab.com"],
      "browserURL": "/Applications/Google Chrome.app"
    },
    {
      "patterns": ["localhost", "127.0.0.1"],
      "browserURL": "/Applications/Google Chrome.app"
    },
    {
      "regexPatterns": [
        "^https://.*\\.internal\\.company\\.com",
        "^https://.*\\.dev\\.local"
      ],
      "browserURL": "/Applications/Firefox.app"
    },
    {
      "patterns": ["stackoverflow.com", "reddit.com"],
      "browserURL": "/Applications/Arc.app"
    }
  ],
  "defaultBrowserURL": "/Applications/Safari.app"
}
```

### Configuration Structure

- **`browsers`**: An array of browser configurations, each containing:
  - **`patterns`** (optional): Array of simple string patterns to match in URLs (case-insensitive)
  - **`regexPatterns`** (optional): Array of regular expression patterns for more complex matching
  - **`browserURL`**: Full path to the browser application (e.g., `/Applications/Google Chrome.app`)
- **`defaultBrowserURL`**: Path to the browser that will be used when no patterns match. This can also be set via the menu bar: **Set Default Browser** → select a browser.

### Pattern Types

- **`patterns`**: Simple string matching (case-insensitive). Checks if the URL contains the pattern string.
- **`regexPatterns`**: Regular expression matching. Uses Go's regexp package. More powerful but requires valid regex syntax.

You can use both `patterns` and `regexPatterns` in the same browser configuration. The first match wins.
You can use any browser installed on your system by providing its full application path. The app can also auto-detect browsers installed in `/Applications` or `~/Applications` - use the **Set Default Browser** menu item to see detected browsers and set one as default.

## Development

### Manual Development

```bash
# Run the app directly
go run main.go

# Test URL handling
go run main.go https://github.com
```

### Manual Build
If you prefer to build manually:

```bash
mkdir -p BrowserRedirectBar.app/Contents/MacOS
mkdir -p BrowserRedirectBar.app/Contents/Resources
go build -o BrowserRedirectBar.app/Contents/MacOS/BrowserRedirectBar main.go
cp Info.plist BrowserRedirectBar.app/Contents/
chmod +x BrowserRedirectBar.app/Contents/MacOS/BrowserRedirectBar
```


## Limitations
- When in a browser clicking on a url ( new tab or not ) even though it matches a given browser config, Will not be detected
