package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

type Repository struct {
	Name     string `json:"name"`
	CloneURL string `json:"clone_url"`
}

func main() {
	// Define command-line flags
	owner := flag.String("owner", "", "GitHub account owner")
	directory := flag.String("directory", "", "Directory to clone repositories")
	flag.Parse()

	if *owner == "" || *directory == "" {
		fmt.Println("Please provide the GitHub account owner and directory")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Create the directory if it doesn't exist
	if err := os.MkdirAll(*directory, 0755); err != nil {
		log.Fatalf("Failed to create directory: %v", err)
	}

	// Make a GET request to the GitHub API to get the starred repositories
	url := fmt.Sprintf("https://api.github.com/users/%s/starred", *owner)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Failed to get starred repositories: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	// Parse the JSON response
	var repositories []Repository
	if err := json.Unmarshal(body, &repositories); err != nil {
		log.Fatalf("Failed to parse JSON response: %v", err)
	}

	// Clone each repository
	for _, repo := range repositories {
		cloneCmd := exec.Command("git", "clone", repo.CloneURL, filepath.Join(*directory, repo.Name))
		if err := cloneCmd.Run(); err != nil {
			log.Printf("Failed to clone repository %s: %v", repo.Name, err)
		} else {
			log.Printf("Cloned repository: %s", repo.Name)
		}
	}
}
