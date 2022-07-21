package lever

import (
	"fmt"

	"github.com/go-lever/lever/config"
	"github.com/kataras/iris/v12"
	"github.com/rs/cors"
)

type App struct {
	*iris.Application
	tlsConfig *tlsConfig
	devMode   bool
	localMode bool
}

// AppOptions handles the App configuration.
// AppOptions should be used when creating a new App with custom parameters like
// * Cors enabled/disabled
type AppOptions struct {
	Cors bool
}

// NewApp creates a new App with the given options
func NewApp(options *AppOptions) *App {
	app := &App{
		Application: iris.Default(),
		tlsConfig:   newTLSConfig(),
		devMode:     config.DevMode(),
		localMode:   config.LocalMode(),
	}

	if options.Cors {
		c := cors.New(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowCredentials: true,
			AllowedMethods:   []string{"HEAD", "GET", "POST", "PUT", "DELETE"},
			// Enable Debugging only in dev mode.
			Debug: app.devMode,
		})
		app.WrapRouter(c.ServeHTTP)
	}

	return app
	
}

// NewDefaultApp creates a new App with the following default parameters
// * CORS are disabled
func NewDefaultApp() *App {
	return NewApp(&AppOptions{})
}

// Run runs the app with TLS (https). Serving http request through https is the default mode. 
// When the app is ran locally or on a remote server that doesn't have a public domain name, a certificate is generated automatically
func (app *App) Run() {

	if app.localMode {
		app.runDev()
		return
	}

	if app.devMode {
		app.runDev()
		return
	}

	if !app.devMode && app.tlsConfig.domain == "" {
		app.runDev()
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
