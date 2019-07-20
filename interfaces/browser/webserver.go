package browser

import (
	"fmt"
	"net/http"
	"os/exec"
	"runtime"

	"github.com/ShrewdSpirit/credman/data"
	"github.com/ShrewdSpirit/credman/gui"
	"github.com/labstack/echo"
)

func Run() {
	e := echo.New()
	e.HideBanner = true

	e.GET("/", func(ctx echo.Context) error {
		return ctx.HTMLBlob(http.StatusOK, gui.Root.Get("app-html").Bytes())
	})

	e.GET("/app.js", func(ctx echo.Context) error {
		return ctx.Blob(http.StatusOK, "text/javascript", gui.Root.Get("app-js").Bytes())
	})

	e.GET("/app.css", func(ctx echo.Context) error {
		return ctx.Blob(http.StatusOK, "text/css", gui.Root.Get("app-css").Bytes())
	})

	e.GET("/favicon.ico", func(ctx echo.Context) error {
		return ctx.Blob(http.StatusOK, "image/x-icon", gui.Root.Get("favicon").Bytes())
	})

	done := make(chan bool)
	go func() {
		e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", data.Config.WebInterfacePort)))
		done <- true
	}()
	openbrowser()
	<-done
}

func openbrowser() error {
	var err error
	url := fmt.Sprintf("http://localhost:%d", data.Config.WebInterfacePort)

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	return err
}
