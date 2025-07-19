package storage

import "gorm.io/gorm"

type Message struct {
	gorm.Model
	ChatID    string
	Content   string
	Role    string
	JsonInput map[string]interface{}
}