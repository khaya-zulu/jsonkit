package agent

type Role string

const (
	RoleAI   Role = "AI"
	RoleUser Role	 = "User"
)

// agents repository idea of a message
type Message struct {
	ID        string     `json:"id"`
	Content   string     `json:"content"`
	Role      Role       `json:"role"`
	JsonInput map[string]interface{} `json:"json_input,omitempty"`
}