package synopsis

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	shared "tobalo/v1/synopsis/pkg/shared"

	"github.com/nats-io/nats.go"
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
        "prompt": "Provide an executive synopsis and intelligence report of the following raw data: %s",
        "max_tokens": 500,
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
	return synopsis, nil
}

// SynopsisFunction fetches a synopsis and returns it
func (c *CohereAPIClient) SynopsisFunction(prompt string) string {
	synopsis, err := c.Generate(prompt)
	if err != nil {
		log.Println("Error getting synopsis:", err)
		return "Error: Could not fetch synopsis" // Provide some fallback text
	}
	Publish(synopsis)
	return synopsis
}

func Publish(msg string) {
	url := os.Getenv("NATS_URL")
	if url == "" {
		url = nats.DefaultURL
	}

	nc, _ := nats.Connect(url)

	defer nc.Drain()

	nc.Publish(shared.SYNOPSIS_SUB, []byte(msg))
}
