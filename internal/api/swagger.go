package api

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"
)

// Arquivos do Swagger UI 5.29.0, licença Apache-2.0, em swaggerui/LICENSE.
//
//go:embed swaggerui/*
var swaggerFiles embed.FS

func swaggerHandler() http.Handler {
	sub, err := fs.Sub(swaggerFiles, "swaggerui")
	if err != nil {
		panic(err)
	}
	files := http.FileServer(http.FS(sub))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/swagger" {
			http.Redirect(w, r, "/swagger/", http.StatusMovedPermanently)
			return
		}
		path := strings.TrimPrefix(r.URL.Path, "/swagger")
		if path == "" {
			path = "/"
		}
		if strings.HasSuffix(path, ".yaml") {
			w.Header().Set("Content-Type", "application/yaml; charset=utf-8")
		}
		clone := r.Clone(r.Context())
		clone.URL.Path = path
		files.ServeHTTP(w, clone)
	})
}
