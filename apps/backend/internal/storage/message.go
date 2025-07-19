package storage

import "gorm.io/gorm"

type Message struct {
	gorm.Model
	ChatId    string
	Content   string
	Role    string
	JsonInput map[string]interface{} `gorm:"serializer:json"`
}