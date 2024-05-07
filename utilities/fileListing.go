package utilities

import (
    "net/http"
    "os"
    "fmt"
    "strings"
    "path/filepath"
    "sort"
)

// FileListHandler handles the /list endpoint
func FileListHandler(w http.ResponseWriter, r *http.Request) {
    // Root directory containing files
    var rootDir string
    distro, err := getDistro()
    if err != nil {
        http.Error(w, "Error detecting distro", http.StatusInternalServerError)
        return
    }
    if distro == "ubuntu" {
        rootDir = "/var/www/html"
    } else if distro == "arch" {
        rootDir = "/usr/share/nginx/html"
    } else {
        http.Error(w, "Unsupported distro", http.StatusInternalServerError)
        return
    }

    // Allowed extensions
    allowedExtensions := map[string]bool{
        ".jpg": true,
        ".png": true,
        ".exe": true,
    }

    // Write the HTML header
    w.Header().Set("Content-Type", "text/html")
    fmt.Fprintf(w, "<h1>File List:</h1>")

    // Map to store files by extension
    fileMap := make(map[string][]string)

    // Iterate over the files in the root directory
    err = filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
        // Skip if there's an error or it's a directory
        if err != nil || info.IsDir() {
            return nil
        }

        // Get file extension
        ext := strings.ToLower(filepath.Ext(info.Name()))

        // Check if the file extension is allowed
        if allowedExtensions[ext] {
            // Add file path to the map
            fileMap[ext] = append(fileMap[ext], strings.TrimPrefix(path, rootDir))
        }

        return nil
    })

    if err != nil {
        http.Error(w, "Error reading directory", http.StatusInternalServerError)
        return
    }

    // Sort file extensions alphabetically
    sortedExts := make([]string, 0, len(fileMap))
    for ext := range fileMap {
        sortedExts = append(sortedExts, ext)
    }
    sort.Strings(sortedExts)

    // Write the list of files by extension
    for _, ext := range sortedExts {
        fmt.Fprintf(w, "<h2>%s Files:</h2>", strings.ToUpper(strings.TrimLeft(ext, ".")))
        fmt.Fprintf(w, "<ul>")
        for _, file := range fileMap[ext] {
            // Escape special characters in filenames to prevent XSS
            filename := strings.ReplaceAll(filepath.Base(file), "<", "&lt;")
            filename = strings.ReplaceAll(filename, ">", "&gt;")
            // Write the HTML list item with download link
            fmt.Fprintf(w, `<li><a href="/download?file=%s" download>%s</a></li>`, filepath.Join("/", strings.TrimPrefix(file, rootDir)), filename)
        }
        fmt.Fprintf(w, "</ul>")
    }
}

