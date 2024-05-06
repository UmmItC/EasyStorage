package main

import (
    "fmt"
    "net/http"
    "EasyStorage/utilities"
)

func main() {
    // Register the fileList handler for the "/list" URL
    http.HandleFunc("/list", utilities.FileListHandler)
    // Register the downloadFile handler for the "/download" URL
    http.HandleFunc("/download", utilities.DownloadFileHandler)

    // Start the server on port 8080
    fmt.Println("Server is listening on port 8080...")
    http.ListenAndServe(":8080", nil)
}

