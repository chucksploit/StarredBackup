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

	// Fetch and clone starred repositories
	page := 1
	perPage := 100
	for {
		repositories, err := fetchStarredRepositories(*owner, page, perPage)
		if err != nil {
			log.Fatalf("Failed to fetch starred repositories: %v", err)
		}

		if len(repositories) == 0 {
			break
		}

		for _, repo := range repositories {
			cloneCmd := exec.Command("git", "clone", repo.CloneURL, filepath.Join(*directory, repo.Name))
			if err := cloneCmd.Run(); err != nil {
				log.Printf("Failed to clone repository %s: %v", repo.Name, err)
			} else {
				log.Printf("Cloned repository: %s", repo.Name)
			}
		}

		page++
	}
}

func fetchStarredRepositories(owner string, page, perPage int) ([]Repository, error) {
	url := fmt.Sprintf("https://api.github.com/users/%s/starred?page=%d&per_page=%d", owner, page, perPage)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get starred repositories: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var repositories []Repository
	if err := json.Unmarshal(body, &repositories); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %v", err)
	}

	return repositories, nil
}
