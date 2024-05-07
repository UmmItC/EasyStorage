package utilities

import (
    "fmt"
    "io/ioutil"
    "strings"
)

// getDistro detects the Linux distribution
func getDistro() (string, error) {
    // Read the contents of the /etc/os-release file
    data, err := ioutil.ReadFile("/etc/os-release")
    if err != nil {
        return "", err
    }

    // Parse the contents to find the ID field
    lines := strings.Split(string(data), "\n")
    for _, line := range lines {
        if strings.HasPrefix(line, "ID=") {
            parts := strings.Split(line, "=")
            if len(parts) == 2 {
                return strings.ToLower(strings.Trim(parts[1], `"`)), nil
            }
        }
    }

    return "", fmt.Errorf("distribution not found")
}

