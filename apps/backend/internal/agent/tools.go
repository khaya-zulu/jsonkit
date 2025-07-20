package agent

import (
	"encoding/json"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/invopop/jsonschema"
	"github.com/itchyny/gojq"
)

type ToolDefinition struct {
	Name string `json:"name"`
	Description string `json:"description"`
	InputSchema  anthropic.ToolInputSchemaParam `json:"inputSchema"`
	Function func(input json.RawMessage) (interface{}, error)
}

func GenerateSchema[T any]() anthropic.ToolInputSchemaParam {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference: true,
	}

	var v T

	schema := reflector.Reflect(v)

	return anthropic.ToolInputSchemaParam{
		Properties: schema.Properties,
	}
}

// jqLang tool
var PerformJqLangDefinition = ToolDefinition{
	Name: "process_json",
	Description: "Process JSON data using jqLang",
	InputSchema: GenerateSchema[JqLangInput](),
	Function: ParseJsonInput,
}

type JqLangInput struct {
	Input any `json:"input" jsonschema_description:"The JSON data to process"`
	Query string          `json:"query" jsonschema_description:"The jq query to apply to the JSON data"`
}

func parseJsonInput(jsonInput any, query string) (any, error) {
	parsedQuery, err := gojq.Parse(query)

	if err != nil {
		fmt.Printf("Error parsing jq query: %v\n", err)
		return nil, err
	}

	var result any

	iterator := parsedQuery.Run(jsonInput)
	for {
		value, ok := iterator.Next()

		if !ok {
			break
		}
		
		if err, ok := value.(error); ok {
			if err, ok := err.(*gojq.HaltError); ok && err.Value() == nil {
				break
			}
		}

		result = value
	}

	return result, nil
}

func ParseJsonInput (input json.RawMessage) (interface{}, error) {
	jsonInput := JqLangInput{}

	err := json.Unmarshal(input, &jsonInput)
	if err != nil {
		fmt.Printf("Error unmarshalling jq input: %v\n", err)
		return nil, err
	}

	result, err := parseJsonInput(jsonInput.Input, jsonInput.Query)

	return result, err
}

