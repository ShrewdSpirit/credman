package browser

import (
	"fmt"
	"net/http"
	"os/exec"
	"runtime"

	"github.com/ShrewdSpirit/credman/data"
	"github.com/ShrewdSpirit/credman/gui"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var server *echo.Echo

func Run(silent bool) {
	server = echo.New()
	server.HideBanner = true

	server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:9000", fmt.Sprintf("http://localhost:%d", data.Config.WebInterfacePort)},
		AllowMethods: []string{http.MethodPost},
	}))

	server.GET("/", func(ctx echo.Context) error {
		return ctx.HTMLBlob(http.StatusOK, gui.Root.Get("app-html").Bytes())
	})

	server.GET("/app.js", func(ctx echo.Context) error {
		return ctx.Blob(http.StatusOK, "text/javascript", gui.Root.Get("app-js").Bytes())
	})

	server.GET("/app.css", func(ctx echo.Context) error {
		return ctx.Blob(http.StatusOK, "text/css", gui.Root.Get("app-css").Bytes())
	})

	server.GET("/favicon.ico", func(ctx echo.Context) error {
		return ctx.Blob(http.StatusOK, "image/x-icon", gui.Root.Get("favicon").Bytes())
	})

	server.POST("/invoke", invokeHandler)

	initInvokeHandler()

	done := make(chan bool)
	go func() {
		server.Logger.Fatal(server.Start(fmt.Sprintf(":%d", data.Config.WebInterfacePort)))
		done <- true
	}()

	if !silent {
		openbrowser()
	}

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
		err = fmt.Errorf("Unsupported platform")
	}

	return err
}
