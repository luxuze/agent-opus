package ai

// Message represents a chat message
type Message struct {
	Role    string // system, user, assistant
	Content string
}

// ChatRequest represents a request to generate a chat response
type ChatRequest struct {
	Model       string
	Messages    []Message
	Temperature float32
	MaxTokens   int
	Stream      bool
}

// ChatResponse represents a response from the AI model
type ChatResponse struct {
	Content      string
	FinishReason string
	TokensUsed   int
	Model        string
}

// AIService defines the interface for AI model interactions
type AIService interface {
	Chat(request ChatRequest) (*ChatResponse, error)
	ChatStream(request ChatRequest) (<-chan string, <-chan error)
}
