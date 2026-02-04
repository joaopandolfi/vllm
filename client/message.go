package client

// Message represents a chat message
type Message struct {
	Role      string      `json:"role"`
	Content   string      `json:"content"`
	ToolCalls []*ToolCall `json:"tool_calls,omitempty"`
}
