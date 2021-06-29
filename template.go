package lever

import (
	"html/template"
	"path/filepath"
)

func (wa *WebApp) FuncMap() template.FuncMap {
	return template.FuncMap{
		"img": wa.img,
	}
}

func (wa *WebApp) img(file string) string {
	return filepath.Join(wa.assetsPath, imgDir, file)
}
