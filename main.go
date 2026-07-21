package main

import (
	"log"
	"os"
	"runtime"
	"sync/atomic"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

const (
	appName     = "Tingly DS"
	deepSeekURL = "https://chat.deepseek.com/"
)

func main() {
	var quitting atomic.Bool
	var app *application.App

	app = application.New(application.Options{
		Name:        appName,
		Description: "A lightweight desktop shell for DeepSeek Chat",
		Icon:        deepSeekColorIcon,
		Mac: application.MacOptions{
			// Keep the app in the Dock while retaining the menu-bar tray entry.
			ActivationPolicy: application.ActivationPolicyRegular,
		},
		Windows: application.WindowsOptions{
			DisableQuitOnLastWindowClosed: true,
		},
		ShouldQuit: func() bool {
			quitting.Store(true)
			return true
		},
		RawMessageHandler: func(window application.Window, message string, origin *application.OriginInfo) {
			target, err := externalURLFromMessage(window, message, origin)
			if err != nil {
				return
			}
			if err := app.Browser.OpenURL(target); err != nil {
				log.Printf("open external URL: %v", err)
			}
		},
	})

	window := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Name:                       "deepseek",
		Title:                      "DeepSeek",
		Width:                      1100,
		Height:                     760,
		MinWidth:                   640,
		MinHeight:                  520,
		URL:                        deepSeekURL,
		BackgroundColour:           application.NewRGB(248, 250, 252),
		JS:                         externalLinkBridgeJS,
		HideOnEscape:               true,
		DevToolsEnabled:            false,
		DefaultContextMenuDisabled: false,
		EnableFileDrop:             false,
		Permissions: map[application.PermissionType]application.Permission{
			application.PermissionMicrophone:    application.PermissionDefault,
			application.PermissionCamera:        application.PermissionDefault,
			application.PermissionGeolocation:   application.PermissionDefault,
			application.PermissionNotifications: application.PermissionDefault,
			application.PermissionClipboardRead: application.PermissionDefault,
		},
	})

	// Preserve the WebView and its session when the user closes the window.
	// A real Quit action sets quitting first, allowing Wails to destroy it.
	window.RegisterHook(events.Common.WindowClosing, func(event *application.WindowEvent) {
		if quitting.Load() {
			return
		}
		event.Cancel()
		window.Hide()
	})

	tray := app.SystemTray.New()
	tray.SetTooltip(appName)
	if runtime.GOOS == "darwin" {
		// Template icons let macOS recolour the official mark for light and dark
		// menu bars without maintaining two platform-specific variants.
		tray.SetTemplateIcon(deepSeekTemplateIcon)
	} else {
		tray.SetIcon(deepSeekColorIcon)
	}

	menu := app.NewMenu()
	menu.Add("显示 DeepSeek").OnClick(func(_ *application.Context) {
		showWindow(window)
	})
	menu.Add("隐藏窗口").OnClick(func(_ *application.Context) {
		window.Hide()
	})
	menu.Add("刷新").OnClick(func(_ *application.Context) {
		window.Reload()
	})
	menu.Add("在浏览器中打开").OnClick(func(_ *application.Context) {
		if err := app.Browser.OpenURL(deepSeekURL); err != nil {
			log.Printf("open DeepSeek in browser: %v", err)
		}
	})
	menu.AddSeparator()
	menu.Add("退出").OnClick(func(_ *application.Context) {
		quitting.Store(true)
		app.Quit()
	})

	tray.SetMenu(menu)
	tray.OnClick(func() {
		if window.IsVisible() {
			window.Hide()
			return
		}
		showWindow(window)
	})

	if err := app.Run(); err != nil {
		log.Printf("run %s: %v", appName, err)
		os.Exit(1)
	}
}

func showWindow(window application.Window) {
	window.Show()
	window.Focus()
}
