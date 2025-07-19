package storage

import (
	"backend/internal/agent"
	"backend/internal/chat"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Storage struct {
	db *gorm.DB
}

func NewStorage() *Storage {
	db, err := gorm.Open(sqlite.Open("jsonkit.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&Chat{}, &Message{})

	return &Storage{db}
}

func toChatMessage(m Message) chat.Message {
	return chat.Message{
		ID:        m.ID,
		Content:   m.Content,
		Role:      agent.Role(m.Role),
		JsonInput: m.JsonInput,
		ChatId:    m.ChatID,
	}
}

func (s *Storage) CreateChat(m chat.NewChatMessage) (chat.Message, error) {
	newChat := Chat{ID: m.ChatId,  Title: "New Message", JsonInput: m.JsonInput }
	
	// Create the chat in the database
	if err := s.db.Create(&newChat).Error; err != nil {
		return chat.Message{}, err
	}

	newMessage := Message{
		ChatID:  newChat.ID,
		Content:  m.Content,
		Role:   string(m.Role),
		JsonInput: m.JsonInput,
	}

	// Create the message associated with the chat
	s.db.Create(&newMessage)

	return toChatMessage(newMessage), nil
}

func (s *Storage) CreateChatMessage(m chat.NewChatMessage) (chat.Message, error) {
	newMessage := Message{
		ChatID:  m.ChatId,
		Content:  m.Content,
		Role:     string(m.Role),
		JsonInput: m.JsonInput,
	}

	s.db.Create(&newMessage)

	return toChatMessage(newMessage), nil
}

func (s *Storage) IsChatExists(chatID string) (bool, error) {
	var count int64
	if err := s.db.Model(&Chat{}).Where("id = ?", chatID).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}