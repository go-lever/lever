package lever

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

type dirFS struct {
	rootDir string
	FS      fs.FS
}

func NewDirFS(rootDir string, fsys fs.FS) *dirFS {
	return &dirFS{
		rootDir: rootDir,
		FS:      fsys,
	}
}

func (dirFS *dirFS) getRootFS() fs.FS {
	fsys, err := fs.Sub(dirFS.FS, dirFS.rootDir)
	if err != nil {
		log.Fatalf("unable to get sub-filesystem : %s", err.Error())
	}
	return fsys
}

func (dirFS *dirFS) getDirFS(dir string) fs.FS {
	fsys, err := fs.Sub(dirFS.FS, dir)
	if err != nil {
		log.Fatalf("unable to get sub-filesystem : %s", err.Error())
	}
	return fsys
}

func (dirFS *dirFS) getRootDir() fs.FS {
	return os.DirFS(dirFS.rootDir)
}

func (dirFS *dirFS) getDir(dir string) fs.FS {
	return os.DirFS(filepath.Join(dirFS.rootDir, dir))
}
