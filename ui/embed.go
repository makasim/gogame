package ui

import (
	"embed"
	"io/fs"
)

//go:embed public
var publicFS embed.FS

func PublicFS() fs.FS {
	fsys, _ := fs.Sub(publicFS, "public")
	return fsys
}
