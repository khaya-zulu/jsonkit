package agent

import (
	"backend/internal/utils"
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
	tools := []ToolDefinition{ParseJsonDefinition}

	return &Agent{ client, tools }
}

func (a *Agent) Run(ctx context.Context, userMessage string, jsonInput interface{}, messages []Message) (Message, error) {
	conversation := convertMessagesToConversation(userMessage, messages)

	// run the anthropic client to generate a response
	response, err := a.runInference(ctx, jsonInput, conversation)

	if err != nil {
		return Message{}, err
	}	
	 
	toolResultList := []anthropic.ContentBlockParamUnion{}
	messageTools := []MessageTool{}

	// Call all the tools in the response
	for _, block := range response.Content {
		switch block := block.AsAny().(type) {
			case anthropic.ToolUseBlock:
				var toolResult anthropic.ContentBlockParamUnion
				if block.Name == ParseJsonDefinition.Name {
				   modifiedInput, _ := utils.AddFieldToRawMessage(block.Input, "Input", jsonInput)
				   _, toolResult = a.executeTool(block.ID, block.Name, modifiedInput)
				} else {
				   _, toolResult = a.executeTool(block.ID, block.Name, block.Input)
				}

				messageTools = append(messageTools, MessageTool{
					ID:     block.ID,
					Input:  block.Input,
					Result: toolResult,
					Name:   block.Name,
				})

				toolResultList = append(toolResultList, toolResult)
		}
	}

	var textResponse string

	if len(toolResultList) > 0 {
		conversation = append(conversation, response.ToParam())

		conversation = append(conversation, anthropic.NewUserMessage(toolResultList...))
		response, err = a.runInference(ctx, jsonInput, conversation)

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
		ToolCalls: messageTools,
	}, nil
}

func generateSystemPrompt(jsonInput interface{}) string {
    return fmt.Sprintf(
        `# JSON Data Processing Agent

You are a specialized agent designed to help users work with JSON data efficiently. Your primary capabilities include generating CloudWatch Insights queries and performing jq operations on JSON data.

## Your Role
- Analyze and process JSON data provided by the user
- Generate CloudWatch Insights queries based on JSON log data
- Execute jq operations to filter, transform, and extract data from JSON
- Provide clear explanations of your operations and results

## JSON Data Context
The following JSON data has been provided for analysis:
` + "```json\n%s\n```" + `

## Guidelines
- Always validate JSON structure before processing
- Provide clear, executable jq commands with explanations
- Generate CloudWatch Insights queries that are optimized for the specific data structure
- Explain the purpose and expected output of each operation
- Handle errors gracefully and suggest alternative approaches
- When generating CloudWatch queries, consider common log analysis patterns like:
  - Error rate analysis
  - Performance metrics
  - Request volume trends
  - Status code distributions
  - Field-specific filtering and aggregation

## Response Format
When providing solutions:
1. State what operation you're performing
2. Execute the appropriate tool (results will be displayed in the UI)
3. Explain key components of the operation
4. Interpret the results and highlight important findings
5. Suggest variations or related operations when relevant

## Example Interactions
- "Generate a CloudWatch query to find all errors in the last hour"
- "Use jq to extract all user IDs from this JSON"
- "Create a query to analyze response times by endpoint"
- "Filter the JSON to show only records with status code 500"

You should be proactive in suggesting useful operations based on the JSON structure and common analysis patterns.`,
        jsonInput,
    )
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

// runInference runs the inference using the anthropic client and returns the response.
func (a *Agent) runInference(ctx context.Context, jsonInput interface{}, conversation []anthropic.MessageParam) (*anthropic.Message, error) {
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
		System: []anthropic.TextBlockParam{
			{
				Text: generateSystemPrompt(jsonInput),
			},
		},
		MaxTokens: int64(1024),
		Tools: anthropicTools,
	})

	return response, err
}

// executeTool executes the tool with the given name and input, returning the result and a tool result block.
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