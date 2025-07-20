package storage

import "gorm.io/gorm"

type Message struct {
	gorm.Model
	ID        string                 `gorm:"primaryKey"`
	ChatId    string                 `json:"chatId"`
	Content   string                 `json:"content"`
	Role      string                 `json:"role"`
	JsonInput map[string]interface{} `gorm:"serializer:json"`
}