package vllm

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/joaopandolfi/core"
	"github.com/joaopandolfi/vllm/client"
)

// Provider implements the LLMProvider interface for VLLM
type Provider struct {
	host string
	port int

	model *core.Model

	// client is the internal VLLM HTTP client
	client *client.VLLMClient

	logger *logr.Logger
}

type ProviderOpts struct {
	BaseURL string
	Port    int
	Logger  *logr.Logger
}

// NewProvider creates a new vllm provider
//
// TODO:
// - need to handle base URL better (with trailing slashes, etc.)
// - need to construct actual URL using baseURL, port, etc.
func NewProvider(opts *ProviderOpts) *Provider {
	opts.Logger.Info("Creating new Provider")

	client := client.NewVLLMClient(
		client.WithBaseURL(opts.BaseURL),
	)

	return &Provider{
		host:   opts.BaseURL,
		port:   opts.Port,
		client: client,
		logger: opts.Logger,
	}
}

func (p *Provider) GetCapabilities(ctx context.Context) (*core.Capabilities, error) {
	p.logger.Info("Fetching capabilities")

	// Placeholder for future implementation
	p.logger.Info("GetCapabilities method is not implemented yet")
	return nil, nil
}

func (p *Provider) UseModel(ctx context.Context, model *core.Model) error {
	p.logger.Info("Setting model", "modelID", model.ID)
	p.model = model
	return nil
}

// convertMessageToVLLMMessage converts a core.Message to a client.Message
func convertMessageToVLLMMessage(msg *core.Message) *client.Message {
	return &client.Message{
		Role:    string(msg.Role),
		Content: msg.Content,
	}
}

// convertManyMessagesToVLLMMessages converts a slice of core.Message to a slice of client.Message
func convertManyMessagesToVLLMMessages(msgs []*core.Message) []*client.Message {
	vllmMsgs := make([]*client.Message, len(msgs))
	for i, msg := range msgs {
		vllmMsgs[i] = convertMessageToVLLMMessage(msg)
	}
	return vllmMsgs
}

// Generate implements the LLMProvider interface for basic responses
func (p *Provider) Generate(ctx context.Context, opts *core.GenerateOptions) (*core.Message, error) {
	p.logger.Info("Generate request received", "modelID", p.model.ID)
	vllmMessages := convertManyMessagesToVLLMMessages(opts.Messages)

	// Convert tools into Ollama's format
	vllmTools := make([]*client.Tool, len(opts.Tools))
	p.logger.Info("Converting tools to VLLM format", "toolCount", len(opts.Tools))

	for i, t := range opts.Tools {
		vllmTools[i] = &client.Tool{
			Type: "function",
			Function: client.ToolFunction{
				Name:        t.Name,
				Description: t.Description,
				Parameters:  t.JSONSchema,
			},
		}
	}

	resp, err := p.client.Chat(ctx, &client.ChatRequest{
		Model:    p.model.ID,
		Messages: vllmMessages,
		Tools:    vllmTools,
	})
	if err != nil {
		p.logger.V(-1).Info(err.Error(), "client error", err)
		return nil, fmt.Errorf("error calling client chat method: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no choices returned from model")
	}

	toolCalls := []*core.ToolCall{}
	for _, toolCall := range resp.Choices[0].Message.ToolCalls {
		toolCalls = append(toolCalls, &core.ToolCall{
			Name:      toolCall.Function.Name,
			Arguments: json.RawMessage([]byte(toolCall.Function.Arguments)),
		})
	}

	return &core.Message{
		Role:      core.AssistantMessageRole,
		Content:   resp.Choices[0].Message.Content,
		ToolCalls: toolCalls,
	}, nil
}

// GenerateStream streams the response token by token
func (p *Provider) GenerateStream(ctx context.Context, opts *core.GenerateOptions) (<-chan *core.Message, <-chan string, <-chan error) {
	p.logger.Info("Stream generation not implemented yet")
	return nil, nil, nil
}
