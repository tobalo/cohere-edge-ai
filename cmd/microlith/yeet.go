package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// CohereAPIClient with Generate method and API key
type CohereAPIClient struct {
	apiKey string
}

// NewClient returns a new CohereAPIClient
func NewClient(apiKey string) *CohereAPIClient {
	return &CohereAPIClient{apiKey: apiKey}
}

// Generate sends a POST request to the Cohere generate endpoint
func (c *CohereAPIClient) Generate(prompt string) (string, error) {
	url := "https://api.cohere.ai/v1/generate"
	payload := fmt.Sprintf(`{
        "model": "command-nightly", 
        "prompt": "%s",
        "max_tokens": 50,
        "temperature": 0.9, 
        "stop_sequences": ["--"]
    }`, prompt)

	req, err := http.NewRequest("POST", url, strings.NewReader(payload))
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %w", err)
	}

	// Extract the synopsis
	var responseMap map[string]interface{}
	if err := json.Unmarshal(body, &responseMap); err != nil {
		return "", fmt.Errorf("error parsing response: %w", err)
	}
	generations, ok := responseMap["generations"].([]interface{})
	if !ok {
		return "", fmt.Errorf("error: 'generations' field not found or not an array")
	}

	firstGeneration, ok := generations[0].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("error: first generation not a map")
	}

	textValue, ok := firstGeneration["text"].(string)
	if !ok {
		return "", fmt.Errorf("error: 'text' field not found or not a string")
	}

	synopsis := textValue
	fmt.Println("###\nSYNOPSIS\n###\n" + synopsis)
	return synopsis, nil
}

// SynopsisFunction fetches a synopsis and returns it
func (c *CohereAPIClient) SynopsisFunction(prompt string) string {
	synopsis, err := c.Generate(prompt)
	if err != nil {
		log.Println("Error getting synopsis:", err)
		return "Error: Could not fetch synopsis" // Provide some fallback text
	}
	return synopsis
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	apiKey := os.Getenv("CO_API_KEY") // Fetch API key from environment
	if apiKey == "" {
		log.Fatal("Missing CO_API_KEY environment variable")
	}

	client := NewClient(apiKey)

	// Define the flag for getting prompt input
	inputPrompt := flag.String("i", "", "Input prompt for the AI")
	inputPromptLong := flag.String("input", "", "Input prompt for the AI")

	// Parse the command-line flags
	flag.Parse()

	// Determine if an input flag was used, else assume console input
	var prompt string
	if *inputPrompt != "" || *inputPromptLong != "" {
		if *inputPrompt != "" {
			prompt = *inputPrompt
		} else {
			prompt = *inputPromptLong
		}
	} else {
		fmt.Print("Enter your prompt: ")
		fmt.Scanln(&prompt) // Get input from the console
	}
	edgePrompt := prompt
	synopsis := client.SynopsisFunction(edgePrompt)
	fmt.Println(synopsis)
}
