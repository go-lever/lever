package lever

import (
	"net/http"

	"github.com/go-lever/encore"
	"github.com/kataras/iris/v12"
)


func NewWebApp(assetFS http.FileSystem, assetPath string, templateFS http.FileSystem) *App {
	app := NewApp()

	assetRenderer := encore.NewRenderer(assetFS, assetPath)
	app.RegisterView(
		iris.HTML(templateFS, ".gohtml").
			Reload(true).
			Funcs(assetRenderer.FuncMap()))

	app.HandleDir(assetPath, assetFS)

	return app
}

