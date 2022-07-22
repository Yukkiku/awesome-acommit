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
		fmt.Pr