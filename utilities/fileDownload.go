package utilities

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type DownloadCounts map[string]map[string]int

func readDownloadCounts(filename string) (DownloadCounts, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			// If the file doesn't exist, return an empty DownloadCounts
			return make(DownloadCounts), nil
		}
		return nil, err
	}

	// Check if the file is empty
	if len(data) == 0 {
		// If the file is empty, return an empty DownloadCounts
		return make(DownloadCounts), nil
	}

	var counts DownloadCounts
	err = json.Unmarshal(data, &counts)
	if err != nil {
		return nil, err
	}
	return counts, nil
}


func updateDownloadCounts(filename, fileType, file string) error {
	counts, err := readDownloadCounts(filename)
	if err != nil {
		return err
	}
	if counts == nil {
		counts = make(DownloadCounts)
	}
	if counts[fileType] == nil {
		counts[fileType] = make(map[string]int)
	}
	counts[fileType][file]++
	data, err := json.MarshalIndent(counts, "", "    ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

// GenerateRandomFilenameWithDate generates a random filename with the current date and the same extension as the original file
func GenerateRandomFilenameWithDate(originalFilename string) (string, error) {
	// Generate random bytes
	randomBytes := make([]byte, 8)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	// Convert random bytes to hex
	randomHex := hex.EncodeToString(randomBytes)

	// Get the extension of the original filename
	ext := filepath.Ext(originalFilename)

	// Get current date in YYYYMMDD format
	currentDate := time.Now().Format("20060102")

	// Construct filename with random hex, current date, and original extension
	filename := fmt.Sprintf("%s-%s%s", randomHex, currentDate, ext)
	return filename, nil
}

// DownloadFileHandler handles the /download endpoint
func DownloadFileHandler(w http.ResponseWriter, r *http.Request) {
	// Get the filename from the URL query parameter
	file := r.URL.Query().Get("file")

	// Generate random output filename with the current date and the same extension as the original file
	randomFilename, err := GenerateRandomFilenameWithDate(file)
	if err != nil {
		http.Error(w, "Error generating random filename", http.StatusInternalServerError)
		return
	}

	// Increment download count for the file
	ext := filepath.Ext(file)
	fileType := ext[1:]
	err = updateDownloadCounts("download_counts.json", fileType, file)
	if err != nil {
		fmt.Println("Error updating download count:", err)
	}

	// Set the Content-Disposition header to specify the filename for the browser
	w.Header().Set("Content-Disposition", "attachment; filename="+randomFilename)

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
		rootDir = "/var/www/ummit_storage"
	case "arch":
		rootDir = "/usr/share/nginx/html"
	default:
		http.Error(w, "Unsupported distro", http.StatusInternalServerError)
		return
	}

	// Construct the full file path
	filePath := filepath.Join(rootDir, file)

	// Log the file path
	fmt.Printf("Server file path: %s\n", filePath) // Printing server file path

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

	// Serve the file for download
	http.ServeContent(w, r, randomFilename, time.Now(), f)
}

