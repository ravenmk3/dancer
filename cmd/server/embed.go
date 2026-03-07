package main

import (
	"embed"
	"io"
	"io/fs"
	"net/http"

	"github.com/labstack/echo/v4"
)

//go:embed all:static
var staticFiles embed.FS

// getStaticFS returns the static file system
func getStaticFS() (fs.FS, error) {
	return fs.Sub(staticFiles, "static")
}

// registerStaticRoutes registers static file routes
func registerStaticRoutes(e *echo.Echo) {
	staticFS, err := getStaticFS()
	if err != nil {
		// Static files not available
		return
	}

	// Serve static files at root path
	e.GET("/*", func(c echo.Context) error {
		path := c.Request().URL.Path
		if path == "/" {
			path = "/index.html"
		}

		// Remove leading slash
		cleanPath := path[1:]
		
		// Try to open the file
		file, err := staticFS.Open(cleanPath)
		if err != nil {
			// File not found, serve index.html for SPA routing
			return serveIndexHTML(c, staticFS)
		}
		defer file.Close()

		// Check if it's a directory
		stat, err := file.Stat()
		if err != nil {
			return serveIndexHTML(c, staticFS)
		}
		
		if stat.IsDir() {
			return serveIndexHTML(c, staticFS)
		}

		// Serve the file
		content, err := io.ReadAll(file)
		if err != nil {
			return serveIndexHTML(c, staticFS)
		}

		// Set content type
		contentType := http.DetectContentType(content)
		return c.Blob(http.StatusOK, contentType, content)
	})
}

// serveIndexHTML serves index.html for SPA routing
func serveIndexHTML(c echo.Context, staticFS fs.FS) error {
	file, err := staticFS.Open("index.html")
	if err != nil {
		return c.String(http.StatusNotFound, "index.html not found")
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to read index.html")
	}

	return c.HTMLBlob(http.StatusOK, content)
}
