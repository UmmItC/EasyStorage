package utilities

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// DownloadCount represents the download count for each file category
type DownloadCount struct {
	sync.Mutex
	Counts map[string]map[string]int `json:"counts"`
}

// CategoryMutex represents a mutex associated with each category
type CategoryMutex struct {
	sync.Mutex
}

// Global variables to store download counts and category mutexes
var downloadCount DownloadCount
var categoryMutexes map[string]*CategoryMutex

// Global variable to store allowed extensions
var allowedExtensions []string

// Initialize download counts, category mutexes and allowed extensions.
func init() {
	downloadCount = DownloadCount{
		Counts: make(map[string]map[string]int),
	}
	categoryMutexes = make(map[string]*CategoryMutex)
	readDownloadCountsFromJSON("download_counts.json")
  readAllowedExtensionsFromJSON("allowed_extensions.json")
}

// Read download counts from JSON file
func readDownloadCountsFromJSON(filename string) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading JSON file: %v\n", err)
		return
	}

  // Check if the JSON data is empty
  if len(data) == 0 {
      fmt.Println("JSON file is empty. Waiting someone download the file.")
      return
  }

	err = json.Unmarshal(data, &downloadCount)
	if err != nil {
		fmt.Printf("Error decoding JSON data: %v\n", err)
		return
	}
}

// Write download counts to JSON file with indentation
func writeDownloadCountsToJSON(filename string) {
	jsonData, err := json.MarshalIndent(downloadCount, "", "    ")
	if err != nil {
		fmt.Printf("Error encoding JSON data: %v\n", err)
		return
	}

	err = ioutil.WriteFile(filename, jsonData, 0644)
	if err != nil {
		fmt.Printf("Error writing JSON file: %v\n", err)
		return
	}
}

// DownloadFileHandler handles the /download endpoint
func DownloadFileHandler(w http.ResponseWriter, r *http.Request) {
	// Get the filename from the URL query parameter
	file := r.URL.Query().Get("file")

	// Root directory containing PDF and EXE directories
	rootDir := "files_download"

	// Construct the full file path
	filePath := filepath.Join(rootDir, file)

	// Generate a random filename with the current date and extension
	randomFilename := generateRandomFilename(file)

	// Set the Content-Disposition header to specify the filename for the browser
	w.Header().Set("Content-Disposition", "attachment; filename="+randomFilename)

	// Serve the file for download with the random filename
	http.ServeFile(w, r, filePath)

	// Get file extension and category
	ext := filepath.Ext(file)
	category := getCategory(ext)

	// Lock the category mutex
	categoryMutex, ok := getCategoryMutex(category)
	if !ok {
		fmt.Printf("Error: Category mutex not found for %s\n", category)
		return
	}
	categoryMutex.Lock()
	defer categoryMutex.Unlock()

	// Initialize category if it doesn't exist
	if _, ok := downloadCount.Counts[category]; !ok {
		downloadCount.Counts[category] = make(map[string]int)
	}

	// Increment the download count for the file
	downloadCount.Counts[category][file]++

	// Log the download with the random filename and download count
	fmt.Printf("File '%s' downloaded as '%s'. Total downloads: %d\n", file, randomFilename, downloadCount.Counts[category][file])

	// Write download counts to JSON file
	writeDownloadCountsToJSON("download_counts.json")
}

// generateRandomFilename generates a random filename with the current date and the same extension as the original file
func generateRandomFilename(originalFilename string) string {
	// Generate a random string
	randomBytes := make([]byte, 8)
	_, err := rand.Read(randomBytes)
	if err != nil {
		// If an error occurs, fallback to a simple timestamp-based filename
		// You may want to handle this error differently in a real application
		return "download"
	}
	randomString := hex.EncodeToString(randomBytes)

	// Get the extension of the original file
	extension := filepath.Ext(originalFilename)

	// Get the current date and format it as YYYYMMDD
	currentDate := time.Now().Format("20060102")

	// Concatenate the random string, current date, and the extension
	randomFilename := randomString + "-" + currentDate + extension

	return randomFilename
}

// Read allowed extensions from JSON file
func readAllowedExtensionsFromJSON(filename string) {
    data, err := ioutil.ReadFile(filename)
    if err != nil {
        fmt.Printf("Error reading JSON file: %v\n", err)
        return
    }

    // Check if the JSON data is empty
    if len(data) == 0 {
        fmt.Println("JSON file is empty")
        return
    }

    var jsonData map[string][]string
    if err := json.Unmarshal(data, &jsonData); err != nil {
        fmt.Printf("Error decoding JSON data: %v\n", err)
        return
    }

    allowedExtensions = jsonData["allowed_extensions"]
}


// getCategory returns the category of a file based on its extension
func getCategory(ext string) string {
    for _, allowedExt := range allowedExtensions {
        if ext == allowedExt {
            switch strings.ToLower(ext) {
            case ".exe":
                return "EXE"
            case ".jpg":
                return "JPG"
            }
        }
    }
    return "" // Return empty string if extension is not allowed
}

// getCategoryMutex returns the mutex associated with the given category
func getCategoryMutex(category string) (*CategoryMutex, bool) {
	mutex, ok := categoryMutexes[category]
	if !ok {
		mutex = &CategoryMutex{}
		categoryMutexes[category] = mutex
	}
	return mutex, true
}
