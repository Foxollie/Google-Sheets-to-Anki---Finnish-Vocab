package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

// DownloadPrivateFile authenticates with a token to download a file from a private repo.
func DownloadPrivateFile(user, repo, branch, filePath, localDest, token string) error {
	// 1. Build the same raw URL structure
	url := fmt.Sprintf(
		"https://api.github.com/repos/%s/%s/contents/%s?ref=%s",
		user, repo, filePath, branch,
	)

	// 2. Create the network request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// 3. Inject your GitHub Personal Access Token for authentication
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github.raw")

	// 4. Execute the authenticated request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("network request failed: %w", err)
	}
	defer resp.Body.Close()

	// Capture authentication or typo errors (e.g., 401 Unauthorized or 404 Not Found)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status: %s", resp.Status)
	}

	// 5. Stream the private file directly to disk
	out, err := os.Create(localDest)
	if err != nil {
		return fmt.Errorf("failed to create local file: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save file content: %w", err)
	}

	return nil
}

func GitPull() {
	token := "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" //this is VERY secure
	// For security, always pull tokens from environment variables instead of hardcoding them
	/*token := "os.Getenv("GITHUB_TOKEN")"
	if token == "" {
		fmt.Println("Error: GITHUB_TOKEN environment variable is not set")
		return
	}*/

	username := "Foxollie"
	repository := "Google-Sheets-to-Anki---Finnish-Vocab"
	branch := "main"
	remotePath := "data/spreadsheet.csv"
	localPath := "vocab.csv"

	fmt.Println("Downloading file from private repository...")
	err := DownloadPrivateFile(username, repository, branch, remotePath, localPath, token)
	if err != nil {
		fmt.Printf("Execution Error: %v\n", err)
		return
	}

	fmt.Println("Private file successfully downloaded and saved!")
}
