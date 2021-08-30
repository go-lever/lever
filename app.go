package lever

import (
	"fmt"

	"github.com/go-lever/lever/config"
	"github.com/kataras/iris/v12"
)

type App struct {
	*iris.Application
	tlsConfig *tlsConfig
	devMode   bool
	localMode bool
}

func NewApp() *App {
	return &App{
		Application: iris.Default(),
		tlsConfig:   newTLSConfig(),
		devMode:     config.DevMode(),
		localMode:   config.LocalMode(),
	}
}

func (app *App) Run() {

	if app.localMode {
		app.runDev()
		return
	}

	if app.devMode {
		app.runDev()
		return
	}

	app.runProd()
}

func (app *App) runDev() {
	app.tlsConfig.generateDevCertificate()
	app.Application.Run(iris.TLS(fmt.Sprintf(":%s", app.tlsConfig.port), app.tlsConfig.certFile, app.tlsConfig.KeyFile))
}

func (app *App) runProd() {
	app.Application.Run(iris.AutoTLS(fmt.Sprintf(":%s", app.tlsConfig.port), app.tlsConfig.domain, app.tlsConfig.email))
}
