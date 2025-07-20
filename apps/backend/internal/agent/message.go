package agent

type Role string

const (
	RoleAssistant Role = "Assistant"  // Changed from "AI" to "Assistant"
	RoleUser      Role = "User"
)

// agents repository idea of a message
type Message struct {
	ID        string     `json:"id"`
	Content   string     `json:"content"`
	Role      Role       `json:"role"`
	JsonInput map[string]interface{} `json:"jsonInput,omitempty"`
}