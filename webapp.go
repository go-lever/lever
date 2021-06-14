package lever

import (
	"log"
	"net/http"

	"github.com/go-lever/encore"
	"github.com/kataras/iris/v12"
)


func NewWebApp(assetFS http.FileSystem, assetPath string, templateFS http.FileSystem) *App {
	app := NewApp()

	app.assetRenderer = encore.NewRenderer(assetFS, assetPath)

	app.RegisterView(
		iris.HTML(templateFS, ".gohtml").
			Reload(true).
			Funcs(app.assetRenderer.FuncMap()))

	app.HandleDir(assetPath, assetFS)

	return app
}

func (app *App) PushAssetsHandler(ctx iris.Context) {

	log.Println("PushAssets is called")

	if pusher, ok := ctx.ResponseWriter().Naive().(http.Pusher); ok {
		log.Println("yolo")
		for _, asset := range []string{"/assets/app.js", "/assets/app.css" } {
			err := pusher.Push(asset, nil)
			if err != nil {
				if err == iris.ErrPushNotSupported {
					log.Println("Push is Not Supported")
					ctx.StopWithText(iris.StatusHTTPVersionNotSupported, "HTTP/2 push not supported.")
				} else {
					log.Println("ERROR", err.Error())
					ctx.StopWithError(iris.StatusInternalServerError, err)
				}
				return
			}
		}

	}
}

