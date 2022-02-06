package lever_test

import (
	"net/http"
	"testing"

	"github.com/go-lever/lever"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/httptest"
)

func newApp(path, schema string) *lever.App {
	app := lever.NewDefaultApp()

	app.Post(path, lever.UseSchema(schema), func(ctx iris.Context) {
		ctx.Text("ok")
	})

	return app
}

func newAppWithParty(path, schema string) *lever.App {
	app := lever.NewDefaultApp()

	party := app.Party("/party")
	{
		party.Post(path, lever.UseSchema(schema), func(ctx iris.Context) {
			ctx.Text("ok")
		})
	}

	return app
}

func TestPostInvalidBody(t *testing.T) {
	// instance := `{"foo": "bar"}`>
	app := newApp("/api", `{"type": "object"}`)

	e := httptest.New(t, app.Application)

	e.POST("/api").WithJSON([]string{"yolo", "yoli"}).
		Expect().
		Status(http.StatusBadRequest).
		JSON().
		Object().
		ContainsKey("errors")
}

func TestPostInvalidJSON(t *testing.T) {
	// instance := `{"foo": "bar"}`>
	app := newApp("/api", `{"type": "object"}`)

	e := httptest.New(t, app.Application)

	e.POST("/api").WithText("yolo").
		Expect().
		Status(http.StatusBadRequest).
		JSON().
		Object().
		ContainsKey("errors")
}

func TestAppWithParty(t *testing.T) {

	app := newAppWithParty("/test", `{"type": "object"}`)

	e := httptest.New(t, app.Application)

	e.POST("/party/test").WithJSON(struct {
		foo string
	}{
		foo: "bar",
	}).
		Expect().
		Status(http.StatusOK)
}

func TestAppWithPartyAndInvalidPayload(t *testing.T) {

	app := newAppWithParty("/test", `{"type": "object"}`)

	e := httptest.New(t, app.Application)

	e.POST("/party/test").WithText("yolo").
		Expect().
		Status(http.StatusBadRequest)
}



