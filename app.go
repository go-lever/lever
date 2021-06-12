package lever

import (
	"github.com/go-lever/lever/config"
	"github.com/kataras/iris/v12"
)

type App struct {
	*iris.Application
	tlsConfig *tlsConfig
	devMode bool
}

func NewApp() *App{
	return &App{
		Application: iris.Default(),
		tlsConfig: newTLSConfig(),
		devMode: config.DevMode(),
	}
}

func (app *App) Run(){
	if app.devMode {
		app.Application.Run(iris.TLS(app.tlsConfig.port, app.tlsConfig.certFile, app.tlsConfig.KeyFile))
	} else {
		app.Application.Run(iris.AutoTLS(app.tlsConfig.port, app.tlsConfig.domain, app.tlsConfig.email))
	}
}