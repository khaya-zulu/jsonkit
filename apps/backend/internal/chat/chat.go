package chat

import "backend/internal/agent"

type Chat struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
	Messages  []Message `json:"messages"`
}

type Message struct {
	ID        uint     `json:"id"`
	Content   string     `json:"content"`
	Role      agent.Role `json:"role"`
	JsonInput map[string]interface{} `json:"json_input,omitempty"`
	ChatId    string                 `json:"chat_id,omitempty"`
}

type NewChatMessage struct {
	Content   string     `json:"content"`
	Role      agent.Role `json:"role"`
	JsonInput map[string]interface{} `json:"json_input,omitempty"`
	ChatId    string                 `json:"chat_id,omitempty"`
}