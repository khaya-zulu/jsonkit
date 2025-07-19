package chat

import (
	"backend/internal/agent"
	"context"
	"fmt"
)

type StorageRepository interface {
	// CreateChatMessage creates a new chat message in the storage.
	CreateChatMessage(message NewChatMessage) (Message, error)
	// CreateChat creates a new chat in the storage.
	CreateChat(m NewChatMessage) (Message, error)
	// IsChatExists checks if a chat with the given ID exists.
	IsChatExists(chatID string) (bool, error)
}

type AgentRepository interface {
	GenerateResponse(ctx context.Context, userMessage string, messages []agent.Message) (agent.Message, error)
}

type Service interface {
	NewChatMessage(userMessage NewChatMessage, messages []agent.Message) (Message, error)
}

type service struct {
	storage StorageRepository
	agent AgentRepository
}

func NewService(storage StorageRepository, agent AgentRepository) Service {
	return &service{storage, agent}
}

// CreateChatMessage creates a new message in the chat with the given ID and user message.
func (s *service) NewChatMessage(userMessage NewChatMessage, messages []agent.Message) (Message, error) {
	chatId := userMessage.ChatId

	newChatMessage := NewChatMessage{
		ChatId: chatId,
		Content: userMessage.Content,
		Role: agent.RoleUser,
		JsonInput: userMessage.JsonInput,
	}

	// Check if chat exists and create message accordingly
	isChatExists, err := s.storage.IsChatExists(chatId)
	if err != nil {
		fmt.Printf("Error checking if chat exists: %v\n", err)
		return Message{}, err
	}

	if !isChatExists {
		_, err = s.storage.CreateChat(newChatMessage)
	} else {
		_, err = s.storage.CreateChatMessage(newChatMessage)
	}

	if err != nil {
		fmt.Printf("Error creating chat message: %v\n", err)
		return Message{}, err
	}

	// Append the new message to the messages slice
	agentMessage, err := s.agent.GenerateResponse(context.Background(), userMessage.Content, messages)
	if err != nil {
		fmt.Printf("Error generating agent response: %v\n", err)
		return Message{}, err
	}

	newMessage, err := s.storage.CreateChatMessage(NewChatMessage{
		ChatId:    chatId,
		Content:   agentMessage.Content,
		Role:      agentMessage.Role,
		JsonInput: agentMessage.JsonInput,
	})

	if err != nil {
		fmt.Printf("Error creating chat: %v\n", err)
		return Message{}, err
	}

	return newMessage, nil
}