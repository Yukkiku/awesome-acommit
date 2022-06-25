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
		fmt.Printf("Error running git 