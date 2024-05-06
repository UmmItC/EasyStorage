package utilities

import (
    "net/http"
    "os"
    "fmt"
    "strings"
    "path/filepath"
)

// FileListHandler handles the /list endpoint
func FileListHandler(w http.ResponseWriter, r *http.Request) {
    // Root directory containing PDF and EXE directories
    rootDir := "files_download"

    // Directories to search for files
    searchDirs := []string{"exe", "JPG"}

    // Allowed extensions
    allowedExtensions := map[string]bool{
        ".exe": true,
        ".jpg": true,
    }

    // Write the HTML header
    w.Header().Set("Content-Type", "text/html")
    fmt.Fprintf(w, "<h1>PDF and EXE files:</h1>")

    // Iterate over the search directories
    for _, dir := range searchDirs {
        // Get the directory path
        dirPath := filepath.Join(rootDir, dir)

        // Get the list of files in the directory
        files, err := os.ReadDir(dirPath)
        if err != nil {
            // Skip if directory not found
            continue
        }

        // Write the directory name as a header
        fmt.Fprintf(w, "<h2>%s Files:</h2>", strings.ToUpper(dir))

        // Write the list of files in the directory
        fmt.Fprintf(w, "<ul>")
        for _, file := range files {
            // Ignore directories
            if file.IsDir() {
                continue
            }

            // Get file extension
            ext := strings.ToLower(filepath.Ext(file.Name()))

            // Check if the file extension is allowed
            if allowedExtensions[ext] {
                // Escape special characters in filenames to prevent XSS
                filename := strings.ReplaceAll(file.Name(), "<", "&lt;")
                filename = strings.ReplaceAll(filename, ">", "&gt;")
                // Write the HTML list item with download link
                fmt.Fprintf(w, `<li><a href="/download?file=%s" download>%s</a></li>`, filepath.Join(dir, file.Name()), filename)
            }
        }
        fmt.Fprintf(w, "</ul>")
    }
}
