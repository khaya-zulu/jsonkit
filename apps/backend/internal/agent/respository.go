package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

type Agent struct {
	client anthropic.Client
	tools  []ToolDefinition
}

func NewAgent() *Agent {
	client := anthropic.NewClient(option.WithAPIKey(os.Getenv("ANTHROPIC_API_KEY")))
	tools := []ToolDefinition{PerformJqLangDefinition}

	return &Agent{ client, tools }
}

// converts the Agent Message to the anthropic.MessageParam
func convertMessagesToConversation(userMessage string, messages []Message) []anthropic.MessageParam {
	conversation := []anthropic.MessageParam{}

	for _, msg := range messages {
		var message anthropic.MessageParam
		switch msg.Role {
			case RoleUser:
				message = anthropic.NewUserMessage(anthropic.NewTextBlock(msg.Content))
			case RoleAssistant:
				message = anthropic.NewAssistantMessage(anthropic.NewTextBlock(msg.Content))
		}

		conversation = append(conversation, message)
	}

	// Add the user's message to the conversation
	conversation = append(conversation, anthropic.NewUserMessage(anthropic.NewTextBlock(userMessage)))

	return conversation
}

func (a *Agent) GenerateResponse(ctx context.Context, userMessage string, messages []Message) (Message, error) {
	conversation := convertMessagesToConversation(userMessage, messages)

	// run the anthropic client to generate a response
	response, err := a.runInference(ctx, conversation)

	if err != nil {
		return Message{}, err
	}	
	 
	toolResultList := []anthropic.ContentBlockParamUnion{}

	// Call all the tools in the response
	for _, block := range response.Content {
		switch block := block.AsAny().(type) {
			case anthropic.ToolUseBlock:
				_, toolResult := a.executeTool(block.ID, block.Name, block.Input)
				toolResultList = append(toolResultList, toolResult)
		}
	}

	var textResponse string

	if len(toolResultList) > 0 {
		conversation = append(conversation, response.ToParam())

		conversation = append(conversation, anthropic.NewUserMessage(toolResultList...))
		response, err = a.runInference(ctx, conversation)

		if err != nil {
			return Message{}, err
		}
		
		// find the first text response in the content blocks
		for _, block := range response.Content {
			if block.Type == "text" {
				textResponse = block.AsText().Text
				break
			}
		}
	} else {
		textResponse = response.Content[0].AsText().Text
	}

	return Message{
		ID:        response.ID,
		Content:   textResponse,
		Role:      RoleAssistant,
		JsonInput: map[string]interface{}{},
	}, nil
}

func (a *Agent) runInference(ctx context.Context, conversation []anthropic.MessageParam) (*anthropic.Message, error) {
	anthropicTools := []anthropic.ToolUnionParam{}

	for _, tool := range a.tools {
		anthropicTools = append(anthropicTools, anthropic.ToolUnionParam{
			OfTool: &anthropic.ToolParam{
				Name:        tool.Name,
				Description: anthropic.String(tool.Description),
				InputSchema: tool.InputSchema,
			},
		})
	}

	// run the anthropic client to generate a response
	response, err := a.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model: anthropic.ModelClaudeOpus4_0,
		Messages: conversation,
		MaxTokens: int64(1024),
		Tools: anthropicTools,
	})

	return response, err
}

func (a *Agent) executeTool(id string, toolName string, input json.RawMessage) (interface{}, anthropic.ContentBlockParamUnion) {
	var toolDef ToolDefinition
	var isFound bool

	for _, tool := range a.tools {
		if tool.Name == toolName {
			toolDef = tool
			isFound = true
			break
		}
	}

	if !isFound {
		println("Tool not found: " + toolName)
		return nil, anthropic.NewToolResultBlock(id, "Tool not found", true)
	}

	fmt.Printf("🔥 %s", id)

	toolResult, err := toolDef.Function(input)
	if err != nil {
		println("Error executing tool: " + err.Error())
		return nil, anthropic.NewToolResultBlock(id, err.Error(), true)
	}

	result, err := json.Marshal(toolResult)
	if err != nil {
		panic(err)
	}

	return toolResult, anthropic.NewToolResultBlock(id, string(result), false)
}