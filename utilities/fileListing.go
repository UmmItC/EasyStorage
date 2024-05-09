package utilities

import (
    "net/http"
    "os"
    "fmt"
    "strings"
    "path/filepath"
    "sort"
    "encoding/json"
    "io/ioutil"
)

// FileListHandler handles the /list endpoint
func FileListHandler(w http.ResponseWriter, r *http.Request) {
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
        rootDir = "/var/www/ummit_storage" // Example path for Ubuntu
    case "arch":
        rootDir = "/usr/share/nginx/html" // Example path for Arch Linux
    default:
        http.Error(w, "Unsupported distro", http.StatusInternalServerError)
        return
    }

    // Read allowed extensions from JSON file
    allowedExtensions, err := readAllowedExtensions("allowed_extensions.json")
    if err != nil {
        http.Error(w, "Error reading allowed extensions", http.StatusInternalServerError)
        return
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
            fileMap[ext] = append(fileMap[ext], filepath.Join("/", strings.TrimPrefix(path, rootDir)))
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
            fmt.Fprintf(w, `<li><a href="/download?file=%s" download>%s</a></li>`, file, filename)
        }
        fmt.Fprintf(w, "</ul>")
    }
}

// Function to read allowed extensions from JSON file
func readAllowedExtensions(filename string) (map[string]bool, error) {
    var data map[string][]string

    // Read JSON file
    file, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, err
    }

    // Unmarshal JSON data
    err = json.Unmarshal(file, &data)
    if err != nil {
        return nil, err
    }

    // Convert list of extensions to map for easy lookup
    allowedExtensions := make(map[string]bool)
    for _, ext := range data["allowed_extensions"] {
        allowedExtensions[ext] = true
    }

    return allowedExtensions, nil
}
