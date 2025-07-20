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
	ID        string     `json:"id"`
	Content   string     `json:"content"`
	Role      agent.Role `json:"role"`
	JsonInput map[string]interface{} `json:"jsonInput,omitempty"`
	ChatId    string                 `json:"chatId,omitempty"`
}

type NewChatMessage struct {
	Content   string     `json:"content"`
	Role      agent.Role `json:"role"`
	JsonInput map[string]interface{} `json:"json_input,omitempty"`
	ChatId    string                 `json:"chatId"`
}

type NewChatMessageRequest struct {
	Content string `json:"content"`
	ChatId string `json:"chatId"`
	JsonInput map[string]interface{} `json:"jsonInput,omitempty"`
	Messages []agent.Message `json:"messages"`
}