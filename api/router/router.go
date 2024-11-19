package router

import (
	"github.com/gin-gonic/gin"

	"docstor-be/internal/document"
	"docstor-be/internal/user"
)

func NewRouter(userHandler *user.UserHandler, documentHandler *document.DocumentHandler) *gin.Engine {
	// create new gin router
	r := gin.Default()

	// Setup user routes
	authRoutes := r.Group("/auth")
	{
		authRoutes.POST("/login", userHandler.Login)
		authRoutes.POST("/register", userHandler.Register)
	}

	// Setup document routes
	documentRoutes := r.Group("/document")
	{
		documentRoutes.GET("/ws", documentHandler.HandleWebSocketConnection)
	}
	return r
}
