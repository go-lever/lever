package main

import (
	"embed"
	"log"
	"net/http"

	"github.com/go-lever/lever"
	"github.com/kataras/iris/v12"
)

//go:embed public/*.json public/*.js public/*.css
var Assets embed.FS

//go:embed templates/*.gohtml
var HTML embed.FS

func main() {
	webapp := lever.NewWebApp(
		lever.EmbededDir("public", Assets),
		"/assets",
		lever.EmbededDir("templates", HTML))

	webapp.Get("/", index)

	webapp.Run()
}

func index(ctx iris.Context) {

	pushAssetsHandler(ctx)

	ctx.ViewLayout("layout")
	ctx.View("index")
}

func pushAssetsHandler(ctx iris.Context) {

	log.Println("PushAssets is called")

	if pusher, ok := ctx.ResponseWriter().Naive().(http.Pusher); ok {
		for _, asset := range []string{"/assets/app.js", "/assets/app.css" } {
			err := pusher.Push(asset, nil)
			if err != nil {
				if err == iris.ErrPushNotSupported {
					ctx.StopWithText(iris.StatusHTTPVersionNotSupported, "HTTP/2 push not supported.")
				} else {
					ctx.StopWithError(iris.StatusInternalServerError, err)
				}
				return
			}
		}

	}
}