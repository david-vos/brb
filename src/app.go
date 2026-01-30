package src

import (
	"browserRedirectBar/src/services"
	_ "embed"

	"github.com/getlantern/systray"
)

//go:embed icon.png
var iconData []byte

// App represents the main application
type App struct {
	configService  *services.ConfigService
	patternService *services.PatternService
	browserService *services.BrowserService
	menuService    *services.MenuService
	urlChan        chan string
}

// NewApp creates a new App instance
func NewApp() (*App, error) {
	configService, err := services.NewConfigService()
	if err != nil {
		return nil, err
	}

	config := configService.GetConfig()
	patternService := services.NewPatternService(config)
	browserService := services.NewBrowserService()
	defaultBrowserService := services.NewDefaultBrowserService()
	urlChan := make(chan string, 10)
	configPath := configService.GetConfigPath()

	var menuService *services.MenuService
	menuService = services.NewMenuService(configPath, urlChan, func(url string) {
		browserPath := patternService.FindBrowserForURL(url)
		if browserPath == "" {
			browserPath = configService.GetConfig().DefaultBrowserURL
		}
		if browserPath == "" {
			browserPath = "/Applications/Safari.app"
		}
		browserService.OpenBrowser(browserPath, url)
	}, configService, func() {
		if err := configService.Load(); err != nil {
			menuService.ShowConfigError(err.Error())
		} else {
			menuService.ClearConfigError()
			patternService.UpdateConfig(configService.GetConfig())
		}
	}, defaultBrowserService)

	return &App{
		configService:  configService,
		patternService: patternService,
		browserService: browserService,
		menuService:    menuService,
		urlChan:        urlChan,
	}, nil
}

// Run starts the menu bar application
func (a *App) Run() {
	systray.Run(a.onReady, a.onExit)
}

// URLChan returns the channel used to receive URLs (e.g. from command line when launched as default browser)
func (a *App) URLChan() chan string {
	return a.urlChan
}

// HandleURL finds the appropriate browser for the URL and opens it (used by tests)
func (a *App) HandleURL(url string) {
	config := a.configService.GetConfig()
	browserPath := a.patternService.FindBrowserForURL(url)
	if browserPath == "" {
		browserPath = config.DefaultBrowserURL
	}
	if browserPath == "" {
		browserPath = "/Applications/Safari.app"
	}
	a.browserService.OpenBrowser(browserPath, url)
}

// onReady is called when the systray is ready (run loop is active)
func (a *App) onReady() {
	services.SetupAppleEventHandler(a.urlChan)
	a.menuService.OnReady(iconData)
}

// onExit is called when the systray exits
func (a *App) onExit() {
	a.menuService.OnExit()
}
