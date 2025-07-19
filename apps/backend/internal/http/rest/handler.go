package rest

import (
	"backend/internal/agent"
	"backend/internal/chat"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Handler(c chat.Service) *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(
		cors.Config{
			// todo: restrict to just the client URL
			AllowOrigins: []string{"*"},
			AllowMethods: []string{"GET", "POST"},
			AllowHeaders: []string{"*"},
		},
	))

	r.POST("/chat", createChat(c))

	return r
}

func createChat(c chat.Service) func(c *gin.Context) {
	return func(ctx *gin.Context) {
		// todo: extend this to accept JsonInput
		var req chat.NewChatMessageRequest

		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(400, gin.H{"error": "Invalid request"})
			return
		}

		newMessage := chat.NewChatMessage{
			Content:   req.Content,
			ChatId:    req.ChatId,
			Role:      agent.RoleUser,
			JsonInput: req.JsonInput,
		}

		newChat, err := c.NewChatMessage(newMessage, req.Messages)
		if err != nil {
			ctx.JSON(500, gin.H{"error": "Failed to create chat"})
			return
		}

		ctx.JSON(201, newChat)
	}
}