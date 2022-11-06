package assets

import (
	"embed"
	"io/fs"
	"net/http"
	"os"
	"path"
)

//go:embed frontend/index.html
//go:embed frontend/dist/*
var Assets embed.FS

type fsFunc func(name string) (fs.File, error)

func (f fsFunc) Open(name string) (fs.File, error) {
	return f(name)
}

func AssetHandler(prefix, root string) http.Handler {
	handler := fsFunc(func(name string) (fs.File, error) {
		assetPath := path.Join(root, name)

		f, err := Assets.Open(assetPath)
		if os.IsNotExist(err) {
			return Assets.Open("frontend/index.html")
		}

		return f, err
	})

	return http.StripPrefix(prefix, http.FileServer(http.FS(handler)))
}
