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
	// Get the absolute path to prevent directory traversal
	upath := r.URL.Path
	if !strings.HasPrefix(upath, "/") {
		upath = "/" + upath
	}
	upath = path.Clean(upath)

	// Check if the file exists in the static filesystem
	f, err := h.staticFS.Open(strings.TrimPrefix(upath, "/"))
	if err == nil {
		// File exists, serve it
		defer func() { _ = f.Close() }()
		stat, err := f.Stat()
		if err == nil {
			// If it's a directory, we want to fall through to index.html (unless it has an index.html inside, but SvelteKit usually handles this)
			// For a typical SPA build, we just want to serve files if they are files.
			if !stat.IsDir() {
				http.FileServer(http.FS(h.staticFS)).ServeHTTP(w, r)
				return
			}
		}
	}

	// File does not exist or is a directory, serve index.html
	indexFile, err := h.staticFS.Open(h.indexPath)
	if err != nil {
		http.Error(w, "index.html not found", http.StatusInternalServerError)
		return
	}
	defer func() { _ = indexFile.Close() }()

	// Serve index.html
	// We read it into memory or use http.ServeContent.
	// Since it's embedded, we can use ServeContent with the file.
	stat, err := indexFile.Stat()
	if err != nil {
		http.Error(w, "Failed to stat index.html", http.StatusInternalServerError)
		return
	}

	// Prevent caching of index.html so updates are seen immediately
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	http.ServeContent(w, r, "index.html", stat.ModTime(), indexFile.(io.ReadSeeker))
}
