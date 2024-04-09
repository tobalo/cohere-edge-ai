## Synopsis: Edge AI Assistant

This tool demonstrates how to use the AI at the edge for raw data to natural language executive synopsis.

**Design**
```mermaid
graph TD
    subgraph Edge AI Service
    EdgeInput[Edge Input] --> PromptInputModule[Process Edge Input]
    PromptInputModule -->|Raw Data Input| MultiModalAPI
    MultiModalAPI[Cohere/GPT/LLaMa2/Mistral/Claude API] --> ResponseProcessingModule[Process API Response]
    ResponseProcessingModule --> OutputModule[Generate Synopsis]
    OutputModule --> EdgeOutput[Edge Output]
    OutputModule -->|Synopsis Output| PublishToNats[Publish to nats.io]
    
    %% Back and forth data flow for feedback or iterative refinement
    MultiModalAPI -->|Request/Response| PromptInputModule
    ResponseProcessingModule -->|Refined Data| PromptInputModule
    end

   subgraph Mothership
    PublishToNats --> HybridCloud[NATS Server]
    HybridCloud --> EmailService[Email Out]
    HybridCloud --> ArchiveService[Archive]
    HybridCloud --> IntegrationService[Integrate with Systems]
    end
    %% Additional description for Multi-Modal API interaction
    class MultiModalAPI multiModalAPIStyle;
    classDef multiModalAPIStyle fill:#f9f,stroke:#333,stroke-width:2px;

  
```

**How It Works**

* The core of the chatbot is a Go program named "synopsis.go".
* It utilizes the Cohere API to generate text responses that are relevant to raw data context at the edge.
* A frontend (e.g., built with a template like the UFO Alien Template) provides a user interface for interacting with the chatbot.

**File Structure**
```
edge-ai-service/
├── pkg/
│   └── synopsis/service.go
|   └── shared/subjects.go
├── cmd/microlith/
│       └── synopsis.go
├── .env
├── Dockerfile
└── README.md
```

**Requirements**

* A Cohere API key ([https://cohere.ai](https://cohere.ai))
* Go programming language ([https://go.dev/](https://go.dev/))
* A frontend web framework or HTML template (optional, for the user interface)
* Nats.io Edge Messaging Fabric Technology ([https://nats.io](https://nats.io))

**Installation**

1. **Obtain a Cohere API key:** Sign up for a Cohere account and get your API key.
2. **Install Go:** Follow the instructions at [https://go.dev/doc/install](https://go.dev/doc/install)
3. **Clone or download this project:** This will give you the `synopsis.go` file.
4. **Set Environment Variable:** Set the `CO_API_KEY` environment variable with your Cohere API key.
   * **Linux/macOS:** `export CO_API_KEY=your_api_key`
   * **Windows:** Use the System Properties settings.

**Development and Testing**
```bash
# Initialize default nats server
nats-server
# Subscribe to synopsis messages
nats sub "edge.synopsis"
```
```bash
go mod download
go run ./cmd/microlith/synopsis.go
```

**Building and Running**

**1. Build the Go binary:**
```bash
go build ./cmd/microlith/synopsis.go
```

This will create an executable file named synopsis (or synopsis.exe on Windows)

2. Run the chatbot:

```bash
./synopsis -i "<co: 0>341 2003-10-11T22:14:15.003Z edge.machine.oci-rover.ai su - ID47 - BOM'su root' failed for lonvick on /dev/pts/8"
```

Use code with caution.
The chatbot will start listening for input from the command line.

Frontend Integration (Optional)

Choose a frontend web framework (e.g., React, Vue.js, Svelte) or use a pre-built template like the UFO Alien Template.

Implement logic in your frontend to:

Send user queries to the synopsis executable (potentially running on a server).
Display the chatbot's responses from Cohere within the user interface.
