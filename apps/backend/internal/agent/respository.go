package agent

import (
	"context"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

type Agent struct {
	client *anthropic.Client
}

func NewAgent() *Agent {
	client := anthropic.NewClient(option.WithAPIKey(""))
	return &Agent{ client: &client }
}

// converts the Agent Message to the anthropic.MessageParam
func convertMessagesToConversation(userMessage string, messages []Message) []anthropic.MessageParam {
	conversation := []anthropic.MessageParam{}

	for _, msg := range messages {
		if msg.Role == RoleUser {
			message := anthropic.NewUserMessage(anthropic.NewTextBlock(msg.Content))
			conversation = append(conversation, message)
		} else if msg.Role == RoleAI {
			message := anthropic.NewAssistantMessage(anthropic.NewTextBlock(msg.Content))
			conversation = append(conversation, message)
		}
	}

	// Add the user's message to the conversation
	conversation = append(conversation, anthropic.NewUserMessage(anthropic.NewTextBlock(userMessage)))

	return conversation
}

func (a *Agent) GenerateResponse(ctx context.Context, userMessage string, messages []Message) (Message, error) {
	conversation := convertMessagesToConversation(userMessage, messages)

	// run the anthropic client to generate a response
	response, err := a.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model: anthropic.ModelClaudeOpus4_0,
		Messages: conversation,
	})

	if err != nil {
		return Message{}, err
	}

	return Message{
		ID:        response.ID,
		Content:   response.Content[0].Text,
		Role:    RoleAI,
		JsonInput: map[string]interface{}{},
	}, nil
}	