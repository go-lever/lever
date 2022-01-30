package main

import (
	"fmt"
	"net/http"

	"github.com/go-lever/lever"
	"github.com/kataras/iris/v12"
)

func main() {
	app := lever.NewApp()
	app.Get("/error-response", handleError)
	app.Get("/problem-response", handleProblem)

	app.Run()
}

func handleError(ctx iris.Context) {
	err := validate()
	ctx.StopWithJSON(http.StatusBadRequest, lever.MultiErrorsJSON(err))
}

func handleProblem(ctx iris.Context) {
	err := validate()
	ctx.StopWithProblem(http.StatusBadRequest,
		iris.NewProblem().Title("validation error").Key("validation-constrains", lever.MultiErrorsSlice(err)))
}

func validate() error {
	err := lever.NewMultiErrors().Append(fmt.Errorf("field #1 should not be empty"))
	err.Append(fmt.Errorf("field #2 should contains a number"))

	return err
}