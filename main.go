package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Choice struct {
	Message Message `json:"message"`
}

type JsonResponse struct {
	Choices []Choice `json:"choices"`
}

var (
	apiKey = os.Getenv("OPENAI_API_KEY")
)

func main() {

	diff, err := getStagedDiff()
	if err != nil {
		fmt.Printf("Error running git diff --staged: %v\n", err)
		return
	}
	result, err := generateText(diff)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	text, err := parseResponse(result)
	if err != nil {
		fmt.Printf("Error unmarshalling JSON: %v\n", err)
		return
	}
	err = commitWithEditor(text)
	if err != nil {
		fmt.Printf("Error running git commit: %v\n", err)
		return
	}
}

func getStagedDiff() (string, error) {
	cmd := exec.Command("git", "diff", "--staged")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	result := strings.TrimSpace(string(output))
	if result == "" {
		return "", fmt.Errorf("No staged changes.")
	}

	return result, nil
}

func fetchPrompt() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	filePath := filepath.Join(home, ".config", "acommit", "prompt.txt")

	content, err := os.ReadFile(filePath)
	if os.IsNotExist(err) {
		fmt.Println("No prompt file found. Creating one at ~/.config/acommit/prompt.txt")
		const prompt = "You are to act as the author of a commit message in git. Your mission is to create clean and comprehensive commit messages in the gitmoji convention with emoji and explain why a change was done. I'll send you an output of 'git diff --staged' command, and you convert it into a commit message. Add a short description of WHY the changes are done after the commit message. Don't start it with 'This commit', just describe the changes. Use the present tense. Commit title must not be longer than 74 characters."
		// mkdir -p
		err = os.MkdirAll(filepath.Dir(filePath), 0755)
		if err != nil {
			return "", err
		}
		file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
		if err != nil {
			return "", err
		}
		defer file.C