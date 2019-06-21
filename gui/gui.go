package gui

import (
	"net/url"

	"github.com/ShrewdSpirit/credman/data"
	"github.com/ShrewdSpirit/credman/gui/assets"
	"github.com/zserge/webview"
)

func Open() {
	data.SaveConfiguration()

	window := webview.New(webview.Settings{
		Title:                  "credman",
		Width:                  800,
		Height:                 600,
		Resizable:              false,
		URL:                    `data:text/html,` + url.PathEscape(assets.Root.Get("app-html").String()),
		Debug:                  true,
		ExternalInvokeCallback: invokeCallback,
	})
	defer window.Exit()

	window.Dispatch(func() {
		window.Eval(assets.Root.Get("app-js").String())
	})

	window.Run()
}

func invokeCallback(window webview.WebView, data string) {

}
