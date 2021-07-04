package lever

import (
	"html/template"
	"path/filepath"
)

func (wa *WebApp) FuncMap() template.FuncMap {
	return template.FuncMap{
		"asset": wa.asset,
	}
}

func (wa *WebApp) asset(path string) string {
	return filepath.Join(wa.assetsPath, path)
}
