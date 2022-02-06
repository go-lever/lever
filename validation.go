package lever

import (
	"fmt"

	"github.com/kataras/iris/v12"
	"github.com/santhosh-tekuri/jsonschema/v5"
)

// UseSchema is a per-route middleware that validates the input payload following
// the provided Json Schema
// The aim of this middleware is to avoid data basic validation in your own middleware.
//
// example :
//
//  app.Post("/entity", lever.UseSchema(`{"type": "object"}`), postEntityHandler)
//  ...
//  func postEntityHandler(ctx iris.Context) {
// 	 ctx.Application().Logger().Info("payload is valid")
// 	 //put your business logic here
//  }
func UseSchema(schema string) iris.Handler {
	return func(ctx iris.Context) {
		sch, err := jsonschema.CompileString(ctx.RouteName(), schema)
		if err != nil {
			ctx.StopWithJSON(iris.StatusBadRequest, MultiErrorsJSON(err))
		}

		var b interface{}
		if err := ctx.ReadJSON(&b); err != nil {
			ctx.StopWithJSON(iris.StatusBadRequest, MultiErrorsJSON(fmt.Errorf("cannot parse JSON : %s", err.Error())))
			return
		}

		if err := sch.Validate(b); err != nil {
			ctx.StopWithJSON(iris.StatusBadRequest, MultiErrorsJSON(fmt.Errorf("invalid request body : %s", err.Error())))
		}
	}
}
