package spa

import (
	"io"
	"io/fs"
	"net/http"
	"path"
	"strings"
)

type SPAHandler struct {
	staticFS  fs.FS
	indexPath string
}

// NewHandler creates a new SPAHandler.
// staticFS is the filesystem containing the static assets (e.g., embedded public folder).
// indexPath is the path to the index.html file within that filesystem (e.g., "index.html").
func NewHandler(staticFS fs.FS, indexPath string) *SPAHandler {
	return &SPAHandler{
		staticFS:  staticFS,
		indexPath: indexPath,
	}
}

func (h *SPAHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	upath := strings.TrimPrefix(path.Clean(r.URL.Path), "/")

	// 1. Try to get file info from the static filesystem
	f, err := h.staticFS.Open(upath)
	if err == nil {
		stat, err := f.Stat()
		if err == nil && !stat.IsDir() {
			err = f.Close()
			if err != nil {
				http.Error(w, "Failed to close file", http.StatusInternalServerError)
				return
			}
			// File exists and is not a dir - serve it directly
			http.FileServer(http.FS(h.staticFS)).ServeHTTP(w, r)
			return
		}
		err = f.Close()
		if err != nil {
			http.Error(w, "Failed to close file", http.StatusInternalServerError)
			return
		}
	}

	// 2. FALLBACK: Serve index.html
	indexFile, err := h.staticFS.Open(h.indexPath)
	if err != nil {
		http.Error(w, "index.html not found", http.StatusInternalServerError)
		return
	}
	defer func() {
		err = indexFile.Close()
		if err != nil {
			http.Error(w, "Failed to close file", http.StatusInternalServerError)
		}
	}()

	stat, err := indexFile.Stat()
	if err != nil {
		http.Error(w, "Failed to stat index.html", http.StatusInternalServerError)
		return
	}

	// Cache headers for SPA entry point
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Use ServeContent which handles Range requests and ModTime correctly
	// The type assertion to io.ReadSeeker works for embedded files
	http.ServeContent(w, r, "index.html", stat.ModTime(), indexFile.(io.ReadSeeker))
}
