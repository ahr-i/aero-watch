package handler

import (
	"net/http"
	"os"
	"path"
	"strings"
)

func spaHandler(distDir string) http.Handler {
	fileServer := http.FileServer(http.Dir(distDir))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestPath := strings.TrimPrefix(path.Clean(r.URL.Path), "/")
		if requestPath == "." {
			requestPath = ""
		}

		if requestPath != "" {
			filePath := path.Join(distDir, requestPath)
			info, err := os.Stat(filePath)
			if err == nil && !info.IsDir() {
				fileServer.ServeHTTP(w, r)
				return
			}
		}

		http.ServeFile(w, r, path.Join(distDir, "index.html"))
	})
}
