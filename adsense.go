package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	adsenseScript = `<script async src="https://pagead2.googlesyndication.com/pagead/js/adsbygoogle.js?client=ca-pub-3534718780470570" crossorigin="anonymous"></script>`
	checkedHTML   = "checked_html.txt"
)

func main() {
	repoDir := "." // Start from the current directory

	// Load checked HTML files from the record
	checked := loadCheckedHTML()

	// Walk through the repository directory and its subdirectories
	err := filepath.Walk(repoDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error accessing path %s: %v\n", path, err)
			return nil // Continue walking
		}

		// Skip directories
		if info.IsDir() {
			return nil // Continue walking
		}

		// Check if the file is an HTML file and hasn't been checked before
		if strings.HasSuffix(strings.ToLower(path), ".html") && !contains(checked, path) {
			// Read the HTML file
			htmlContent, err := ioutil.ReadFile(path)
			if err != nil {
				fmt.Printf("Error reading file %s: %v\n", path, err)
				return nil // Continue walking
			}

			// Check if Google Adsense script is present
			if !strings.Contains(string(htmlContent), adsenseScript) {
				// Append Google Adsense script to HTML content
				updatedContent := appendAdsenseScript(htmlContent)

				// Write updated HTML content back to file
				err := ioutil.WriteFile(path, updatedContent, 0644)
				if err != nil {
					fmt.Printf("Error writing file %s: %v\n", path, err)
					return nil // Continue walking
				}

				fmt.Println("Google Adsense script added to:", path)
			}

			// Record the checked HTML file
			recordCheckedHTML(path)
		}

		return nil // Continue walking
	})

	if err != nil {
		fmt.Println("Error walking through directory:", err)
	}
}

// Function to append Google Adsense script to HTML content
func appendAdsenseScript(content []byte) []byte {
	// Find position to insert script before </head> tag
	headIndex := strings.LastIndex(string(content), "</head>")
	if headIndex == -1 {
		// If </head> tag not found, append script to end of file
		return append(content, []byte(adsenseScript)...)
	}

	// Insert script before </head> tag
	return []byte(string(content[:headIndex]) + adsenseScript + string(content[headIndex:]))
}

// Function to load checked HTML files from the record
func loadCheckedHTML() []string {
	var checked []string

	// Check if the record file exists
	if _, err := os.Stat(checkedHTML); err == nil {
		// Read checked HTML files from the record
		data, err := ioutil.ReadFile(checkedHTML)
		if err != nil {
			fmt.Println("Error reading record file:", err)
			return checked
		}

		checked = strings.Split(string(data), "\n")
	}

	return checked
}

// Function to record checked HTML files
func recordCheckedHTML(filePath string) {
	// Open the record file in append mode
	file, err := os.OpenFile(checkedHTML, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening record file:", err)
		return
	}
	defer file.Close()

	// Write the file path to the record file
	if _, err := file.WriteString(filePath + "\n"); err != nil {
		fmt.Println("Error writing to record file:", err)
	}
}

// Function to check if a slice contains a string
func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}
