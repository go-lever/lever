package main

import (
	"embed"

	"github.com/go-lever/lever"
	"github.com/kataras/iris/v12"
)

//go:embed public/*.json public/*.js public/*.css
var Assets embed.FS

//go:embed templates/*.gohtml
var HTML embed.FS

func main() {

	webapp := lever.NewDefaultWebApp(
		lever.NewDirFS("public", Assets),
		lever.NewDirFS("templates", HTML))

	webapp.Get("/", index)

	webapp.Run()
}

func index(ctx iris.Context) {
	ctx.ViewLayout("layout")
	ctx.View("index")
}