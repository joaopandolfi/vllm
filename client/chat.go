package client

// ChatRequest represents a request to the chat completions endpoint
type ChatRequest struct {
	Model    string     `json:"model"`
	Messages []*Message `json:"messages"`
	Tools    []*Tool    `json:"tools,omitempty"`
	Stream   *bool      `json:"stream,omitempty"`
}

// ChatResponse represents a response from the chat completions endpoint
type ChatResponse struct {
	ID      string    `json:"id"`
	Object  string    `json:"object"`
	Created int64     `json:"created"`
	Model   string    `json:"model"`
	Choices []*Choice `json:"choices"`
	Done    bool      `json:"done,omitempty"` // ollama specific
}

// Choice represents a choice in the chat response
type Choice struct {
	Index        int      `json:"index"`
	Message      *Message `json:"message"`
	Delta        *Message `json:"delta"`
	FinishReason string   `json:"finish_reason"`
}