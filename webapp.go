package lever

import (
	"io/fs"
	"net/http"
	"path/filepath"

	"github.com/Masterminds/sprig"
	"github.com/go-lever/encore"
	"github.com/kataras/iris/v12"
)

const (
	defaultBuildDir = "build"
	defaultAssetPath = "/assets"
)

type WebApp struct {
	*App
	assetsFS    *dirFS
	templatesFS *dirFS
	assetsPath  string
}

// WebAppOptions handles the WebApp configuration.
// WebAppOptions should be used when creating a new WebApp with custom parameters.
type WebAppOptions struct {
	AppOptions
	assetPath string
}

// NewWebApp creates a new WebApp with the given options
func NewWebApp(assetsFS *dirFS, templatesFS *dirFS, options *WebAppOptions) *WebApp {
	webapp := &WebApp{
		App:         NewApp(&options.AppOptions),
		assetsFS:    assetsFS,
		templatesFS: templatesFS,
		assetsPath:  options.assetPath,
	}

	if options.assetPath == "" {
		webapp.assetsPath = defaultAssetPath
	}

	encoreRenderer := encore.NewRenderer(webapp.assetsBuild(), filepath.Join(webapp.assetsPath, defaultBuildDir))
	webapp.RegisterView(
		iris.HTML(http.FS(webapp.templates()), ".gohtml").
			Reload(true).
			Funcs(encoreRenderer.FuncMap()).
			Funcs(webapp.FuncMap()).
			Funcs(sprig.FuncMap()))

	webapp.HandleDir(webapp.assetsPath, http.FS(webapp.assets()))

	return webapp

}

// NewDefaultWebApp creates a new WebApp with the following default parameters
// * CORS are disabled
// * Assets are served on the default /assets path
func NewDefaultWebApp(assetsFS *dirFS, templatesFS *dirFS) *WebApp {
	return NewWebApp(assetsFS, templatesFS, &WebAppOptions{})
}

func (wa *WebApp) templates() fs.FS {
	if wa.devMode {
		return wa.templatesFS.getRootDir()
	}

	return wa.templatesFS.getRootFS()
}

func (wa *WebApp) assets() fs.FS {
	if wa.devMode {
		return wa.assetsFS.getRootDir()
	}

	return wa.assetsFS.getRootFS()
}

func (wa *WebApp) assetsBuild() fs.FS {
	if wa.devMode {
		return wa.assetsFS.getDir(defaultBuildDir)
	}

	fsys, err := fs.Sub(wa.assets(), defaultBuildDir)
	if err != nil {
		panic(err)
	}

	return fsys
}
