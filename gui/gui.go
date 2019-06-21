// +build gui

package gui

import (
	"fmt"
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
		window.Eval(fmt.Sprintf("window.AppVersion='%s';window.CommitHash='%s'", data.Version, data.GitCommit))
		window.Eval(assets.Root.Get("app-js").String())
	})

	window.Run()
}

func invokeCallback(window webview.WebView, data string) {

}
