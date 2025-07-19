package storage

import "gorm.io/gorm"

type Chat struct {
	gorm.Model
	ID string `json:"id"`
	Title     string `json:"title"`
	JsonInput map[string]interface{} `gorm:"serializer:json"`
}
