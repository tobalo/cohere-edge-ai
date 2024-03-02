package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"tobalo/zeticuli-demo/pkg/synopsis"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	apiKey := os.Getenv("CO_API_KEY") // Fetch API key from environment
	if apiKey == "" {
		log.Fatal("Missing CO_API_KEY environment variable")
	}

	client := synopsis.NewClient(apiKey)

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
