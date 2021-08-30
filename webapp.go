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
	buildDir = "build"
	imgDir   = "img"
)

type WebApp struct {
	*App
	assetsFS    *dirFS
	templatesFS *dirFS
	assetsPath  string
}

func NewWebApp(assetsFS *dirFS, assetPath string, templatesFS *dirFS) *WebApp {
	webapp := &WebApp{
		App:         NewApp(),
		assetsFS:    assetsFS,
		templatesFS: templatesFS,
		assetsPath:  assetPath,
	}

	encoreRenderer := encore.NewRenderer(webapp.assetsBuild(), filepath.Join(assetPath, buildDir))
	webapp.RegisterView(
		iris.HTML(http.FS(webapp.templates()), ".gohtml").
			Reload(true).
			Funcs(encoreRenderer.FuncMap()).
			Funcs(webapp.FuncMap()).
			Funcs(sprig.FuncMap()))

	webapp.HandleDir(assetPath, http.FS(webapp.assets()))

	return webapp
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
		return wa.assetsFS.getDir(buildDir)
	}

	fsys, err := fs.Sub(wa.assets(), buildDir)
	if err != nil {
		panic(err)
	}

	return fsys
}
