package main

// Import the necessary packages
import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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

// ExampleFunction is a function that uses the Cohere API to interact with a chatbot
func (c *CohereAPIClient) ExampleFunction() {
	// Define the parameters for the request
	params := map[string]string{
		"model":       "chat",
		"prompt":      "Hello!",
		"api_key":     c.apiKey,
		"n_turns":     "1",
		"temperature": "neutral",
	}

	// Make the request to the Cohere API
	response, err := c.Request("generate", params)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	// Read the response body as a string
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Print the response
	fmt.Println(string(responseBody))
}

// ExampleMain is the main function for the chatbot example
func Main() {
	// Create a new CohereAPIClient with a sample API key (you would replace this with your actual API key)
	client := NewClient("YOUR_API_KEY_HERE")

	// Call the ExampleFunction to interact with the chatbot
	client.ExampleFunction()
}
