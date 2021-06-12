package lever

import (
	"io/fs"
	"log"
	"net/http"
	"os"
)


func EmbededDir(dir string, fsys fs.FS) http.FileSystem {

	fsys, err := fs.Sub(fsys, dir)
	if err != nil {
		log.Fatalf("unable to get sub-filesystem : %s", err.Error())
	}
	return http.FS(fsys)
}

func Dir(dir string) http.FileSystem {
	return http.FS(os.DirFS(dir))
}