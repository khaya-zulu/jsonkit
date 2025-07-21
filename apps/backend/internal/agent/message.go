package agent

type Role string

const (
	RoleAssistant Role = "Assistant"  // Changed from "AI" to "Assistant"
	RoleUser      Role = "User"
)

type MessageTool struct {
	ID          string `json:"id"`
	Input       any    `json:"input"`
	Result      any    `json:"result"`
	Name		string `json:"name"`
}

// agents repository idea of a message
type Message struct {
	ID        string     `json:"id"`
	Content   string     `json:"content"`
	Role      Role       `json:"role"`
	ToolCalls []MessageTool `json:"toolCalls,omitempty"`
	JsonInput map[string]interface{} `json:"jsonInput,omitempty"`
}