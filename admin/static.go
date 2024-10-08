package admin

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"path/filepath"

	"github.com/rangidev/rangi/config"
)

var (
	//go:embed static
	staticFS embed.FS
)

func NewStaticServer(config *config.Config) (http.Handler, error) {
	if config.EnableTemplateDevelopment {
		return http.FileServer(http.Dir(filepath.Join(config.ExecutableDir, "admin/static"))), nil
	}
	staticSubFS, err := fs.Sub(staticFS, "static")
	if err != nil {
		return nil, fmt.Errorf("could not get filesystem sub directory: %v", err)
	}
	return http.FileServer(http.FS(staticSubFS)), nil
}
