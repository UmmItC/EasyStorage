package utilities

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
  "time"
)

// DownloadFileHandler handles the /download endpoint
func DownloadFileHandler(w http.ResponseWriter, r *http.Request) {
	// Get the filename from the URL query parameter
	file := r.URL.Query().Get("file")

	// Detect the Linux distribution
	distro, err := getDistro()
	if err != nil {
		http.Error(w, "Error detecting distro", http.StatusInternalServerError)
		return
	}

	// Determine the root directory based on the detected distribution
	var rootDir string
	switch distro {
	case "ubuntu":
		rootDir = "/var/www/html"
	case "arch":
		rootDir = "/usr/share/nginx/html"
	default:
		http.Error(w, "Unsupported distro", http.StatusInternalServerError)
		return
	}

	// Construct the full file path
	filePath := filepath.Join(rootDir, file)

	// Log the file path
	fmt.Printf("File path: %s\n", filePath)

	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Printf("File not found: %s\n", filePath)
		http.NotFound(w, r)
		return
	}

	// Open the file
	f, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "Error opening file", http.StatusInternalServerError)
		return
	}
	defer f.Close()

	// Set the Content-Disposition header to specify the filename for the browser
	filename := filepath.Base(file)
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)

	// Serve the file for download
	http.ServeContent(w, r, filename, time.Now(), f)
}

