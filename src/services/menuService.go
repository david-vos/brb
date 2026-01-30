package services

import (
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/getlantern/systray"
)

// MenuService handles the menu bar setup and interactions
type MenuService struct {
	configPath            string
	urlChan               chan string
	handleURL             func(string)
	configService         *ConfigService
	onConfigUpdated       func() // Callback to reload config when default browser is changed
	defaultBrowserService *DefaultBrowserService
	browserMenuItems      map[*systray.MenuItem]string // Map of menu items to browser paths
	browsers              []BrowserInfo                // List of detected browsers
	mSetDefault           *systray.MenuItem            // Parent menu item for "Set Default Browser"
	mConfigError          *systray.MenuItem            // Menu item shown when config has errors
	configError           string                       // Current config error message
}

// NewMenuService creates a new MenuService instance
func NewMenuService(configPath string, urlChan chan string, handleURL func(string), configService *ConfigService, onConfigUpdated func(), defaultBrowserService *DefaultBrowserService) *MenuService {
	return &MenuService{
		configPath:            configPath,
		urlChan:               urlChan,
		handleURL:             handleURL,
		configService:         configService,
		onConfigUpdated:       onConfigUpdated,
		defaultBrowserService: defaultBrowserService,
	}
}

// OnReady sets up the menu bar when systray is ready
func (ms *MenuService) OnReady(icon []byte) {
	systray.SetIcon(icon)
	systray.SetTooltip("Browser Redirect Bar")

	// Detect browsers
	detector := NewBrowserDetector()
	browsers := detector.DetectBrowsers()

	// Check for config errors on startup
	ms.checkConfigErrors()

	// Add menu items
	mQuit := systray.AddMenuItem("Quit", "Quit Browser Redirect Bar")
	systray.AddSeparator()

	// Add config error menu item if there's an error (initially hidden)
	ms.mConfigError = systray.AddMenuItem("Config Error - Click for details", "Configuration file has errors")
	ms.mConfigError.Hide()
	systray.AddSeparator()

	mSetAsDefault := systray.AddMenuItem("Set as Default Browser", "Request this app to be the default browser")
	systray.AddSeparator()
	mSetDefault := systray.AddMenuItem("Set Default Browser", "Choose default browser for all requests")
	systray.AddSeparator()
	mReloadConfig := systray.AddMenuItem("Reload Config", "Reload configuration from disk")
	mConfig := systray.AddMenuItem("Go to Config File", "Open config file in Finder")

	// Store references for later updates
	ms.browsers = browsers
	ms.mSetDefault = mSetDefault
	ms.browserMenuItems = make(map[*systray.MenuItem]string)

	// Create initial submenu items
	ms.updateBrowserMenuItems()

	// Handle menu item clicks
	go func() {
		for {
			select {
			case url := <-ms.urlChan:
				ms.handleURL(url)
			case <-mQuit.ClickedCh:
				systray.Quit()
				return
			case <-mSetAsDefault.ClickedCh:
				if ms.defaultBrowserService != nil {
					_ = ms.defaultBrowserService.RequestDefaultBrowser()
				}
			case <-mReloadConfig.ClickedCh:
				ms.reloadConfig()
			case <-mConfig.ClickedCh:
				ms.openConfigFile()
			case <-ms.mConfigError.ClickedCh:
				// Show error notification when user clicks the error menu item
				ms.showConfigErrorNotification()
			}
		}
	}()

}

// OnExit is called when the systray exits
func (ms *MenuService) OnExit() {
	// Cleanup if needed
}

// updateBrowserMenuItems updates the browser menu items with current default browser checkmarks
func (ms *MenuService) updateBrowserMenuItems() {
	for menuItem := range ms.browserMenuItems {
		menuItem.Hide()
	}
	ms.browserMenuItems = make(map[*systray.MenuItem]string)

	currentDefault := ms.configService.GetConfig().DefaultBrowserURL

	if len(ms.browsers) == 0 {
		mNoBrowsers := ms.mSetDefault.AddSubMenuItem("No browsers detected", "")
		mNoBrowsers.Disable()
		return
	}

	for _, browser := range ms.browsers {
		menuText := browser.Name
		if browser.Path == currentDefault {
			menuText += " [default]"
		}
		menuItem := ms.mSetDefault.AddSubMenuItem(menuText, "Set as default browser")
		ms.browserMenuItems[menuItem] = browser.Path

		go func(item *systray.MenuItem, path string) {
			for range item.ClickedCh {
				if err := ms.configService.SetDefaultBrowser(path); err != nil {
					ms.ShowConfigError(fmt.Sprintf("Cannot set default browser: %v", err))
					continue
				}
				ms.ClearConfigError()
				if ms.onConfigUpdated != nil {
					ms.onConfigUpdated()
				}
				ms.updateBrowserMenuItems()
			}
		}(menuItem, browser.Path)
	}
}

// checkConfigErrors checks if there are config errors and updates the menu
func (ms *MenuService) checkConfigErrors() {
	// Try to load the config to see if there are errors
	err := ms.configService.Load()
	if err != nil {
		ms.configError = err.Error()
		if ms.mConfigError != nil {
			ms.mConfigError.Show()
		}
		// Show notification
		ms.showConfigErrorNotification()
	} else {
		ms.configError = ""
		if ms.mConfigError != nil {
			ms.mConfigError.Hide()
		}
	}
}

// showConfigErrorNotification shows a macOS notification about the config error
func (ms *MenuService) showConfigErrorNotification() {
	errorMsg := "Invalid JSON in config file"
	if ms.configError != "" {
		errorMsg = ms.configError
	}

	// Use osascript to show a macOS notification
	script := fmt.Sprintf(`display notification "%s\n\nPlease fix the JSON syntax in:\n%s" with title "Browser Redirect Bar - Config Error"`,
		errorMsg, ms.configPath)
	_ = exec.Command("osascript", "-e", script).Run()
}

// ShowConfigError displays a config error in the menu and shows a notification
func (ms *MenuService) ShowConfigError(errorMsg string) {
	ms.configError = errorMsg
	if ms.mConfigError != nil {
		ms.mConfigError.Show()
	}
	ms.showConfigErrorNotification()
}

// ClearConfigError clears the config error from the menu
func (ms *MenuService) ClearConfigError() {
	ms.configError = ""
	if ms.mConfigError != nil {
		ms.mConfigError.Hide()
	}
}

// reloadConfig reloads the configuration from disk
func (ms *MenuService) reloadConfig() {
	if err := ms.configService.Load(); err != nil {
		ms.ShowConfigError(fmt.Sprintf("Failed to reload config: %v", err))
		return
	}
	ms.ClearConfigError()
	if ms.onConfigUpdated != nil {
		ms.onConfigUpdated()
	}
	ms.updateBrowserMenuItems()
	_ = exec.Command("osascript", "-e", `display notification "Configuration reloaded" with title "Browser Redirect Bar"`).Run()
}

// openConfigFile opens the config file location in Finder
func (ms *MenuService) openConfigFile() {
	_ = exec.Command("open", filepath.Dir(ms.configPath)).Run()
}
