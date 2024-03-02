package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

// CohereAPIClient is a struct that defines the Cohere API client
type CohereAPIClient struct {
	apiKey string
}

// NewClient returns a new CohereAPIClient with the given API key
func NewClient(apiKey string) *CohereAPIClient {
	return &CohereAPIClient{apiKey: apiKey}
}

// Request performs a GET request to the Cohere API
func (c *CohereAPIClient) Request(endpoint string, params map[string]string) (*http.Response, error) {
	// Create a URL with the endpoint and parameters
	url := fmt.Sprintf("https://api.cohere.com/v2/%s?api_key=%s", endpoint, c.apiKey)
	for key, value := range params {
		url += fmt.Sprintf("&%s=%s", key, value)
	}

	// Make the request
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("could not make request: %w", err)
	}

	return resp, nil
}

// Optimized ExampleFunction with http.Client reuse, timeout, and batching
func (c *CohereAPIClient) ExampleFunction() {
	// Define the parameters for the request
	params := map[string]string{
		"model":       "chat",
		"prompt":      "Hello!",
		"api_key":     c.apiKey,
		"n_turns":     "1",
		"temperature": "neutral",
	}

	client := &http.Client{
		Timeout: 10 * time.Second, // Set a reasonable timeout, adjust as needed
	}

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ { // Example of making multiple requests
		wg.Add(1)
		go func() {
			defer wg.Done()
			response, err := client.Get(url) // Reuse http.Client
			if err != nil {
				log.Println("Error in request:", err)
				return
			}
			defer response.Body.Close()

			responseBody, err := ioutil.ReadAll(response.Body)
			if err != nil {
				log.Println("Error reading response:", err)
				return
			}

			fmt.Println(string(responseBody))
		}()
	}

	wg.Wait()
}

func main() {
	// Create a new CohereAPIClient with a sample API key
	client := NewClient("YOUR_API_KEY_HERE")

	// Call the ExampleFunction to interact with the chatbot
	client.ExampleFunction()
}
