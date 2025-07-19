package rest

import (
	"backend/internal/agent"
	"backend/internal/chat"

	"github.com/gin-gonic/gin"
)

func Handler(c chat.Service) *gin.Engine {
	r := gin.Default()

	r.POST("/chats", createChat(c))

	return r
}

func createChat(c chat.Service) func(c *gin.Context) {
	return func(ctx *gin.Context) {
		// todo: extend this to accept JsonInput
		var req chat.NewChatMessage

		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(400, gin.H{"error": "Invalid request"})
			return
		}

		newChat, err := c.NewChatMessage(req, []agent.Message{})
		if err != nil {
			ctx.JSON(500, gin.H{"error": "Failed to create chat"})
			return
		}

		ctx.JSON(201, newChat)
	}
}