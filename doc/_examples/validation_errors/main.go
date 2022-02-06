package main

import (
	"github.com/go-lever/lever"
	"github.com/kataras/iris/v12"
)

var schema = `
{
	"type": "object",
	"properties": {
		"firstName": {
			"type": "string",
		},
		"lastName": {
			"type": "string",
		},
		"age": {
			"type": "integer",
			"minimum": 0
		}
	}
}
`

func main() {
	app := lever.NewDefaultApp()
	app.Post("/entity", lever.UseSchema(schema), postEntityHandler)

	app.Run()
}

func postEntityHandler(ctx iris.Context) {
	ctx.Application().Logger().Info("payload is valid")
	//put your business logic here
}
