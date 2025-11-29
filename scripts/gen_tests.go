package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func main() {
	filePath := flag.String("file", "", "Path to the Go file to generate tests for")
	flag.Parse()

	if *filePath == "" {
		log.Fatal("Please provide a file path using -file flag")
	}

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatal("GEMINI_API_KEY environment variable is not set")
	}

	content, err := ioutil.ReadFile(*filePath)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")

	// Configure generation config if needed, e.g. temperature
	// model.SetTemperature(0.2)

	prompt := fmt.Sprintf(`You are an expert Go developer. Generate comprehensive unit tests for the following Go code using the standard 'testing' package. 
Output ONLY the code for the test file, including package declaration and imports. 
Do not include markdown code blocks or any other text.

Code:
%s`, string(content))

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		log.Fatal(err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		log.Fatal("No content generated")
	}

	// Extract text from response
	var testContent string
	for _, part := range resp.Candidates[0].Content.Parts {
		if txt, ok := part.(genai.Text); ok {
			testContent += string(txt)
		}
	}

	// Clean up potential markdown code blocks
	testContent = strings.TrimPrefix(testContent, "```go")
	testContent = strings.TrimPrefix(testContent, "```")
	testContent = strings.TrimSuffix(testContent, "```")
	// Sometimes the model might put the language identifier on the first line
	lines := strings.Split(testContent, "\n")
	if len(lines) > 0 && strings.HasPrefix(lines[0], "```") {
		lines = lines[1:]
		testContent = strings.Join(lines, "\n")
	}
	// Remove trailing backticks if any
	if len(lines) > 0 && strings.HasPrefix(lines[len(lines)-1], "```") {
		lines = lines[:len(lines)-1]
		testContent = strings.Join(lines, "\n")
	}

	dir := filepath.Dir(*filePath)
	baseName := filepath.Base(*filePath)
	ext := filepath.Ext(baseName)
	testFileName := strings.TrimSuffix(baseName, ext) + "_test.go"
	testFilePath := filepath.Join(dir, testFileName)

	err = ioutil.WriteFile(testFilePath, []byte(testContent), 0644)
	if err != nil {
		log.Fatalf("Failed to write test file: %v", err)
	}

	fmt.Printf("Generated tests in %s\n", testFilePath)
}
