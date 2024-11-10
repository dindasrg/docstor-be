package router

import (
	"github.com/gin-gonic/gin"

	"docstor-be/internal/user"
)

type Router struct {
	UserHandler *user.UserHandler
}

func NewRouter(userHandler *user.UserHandler) *gin.Engine {
	// create new gin router
	r := gin.Default()

	// auth routes
	r.POST("/login", userHandler.Login)
	r.POST("/register", userHandler.Register)

	return r
}
