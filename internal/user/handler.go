package user

import (
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UserService *UserService
}

func NewUserHandler(UserService *UserService) *UserHandler {
	return &UserHandler{UserService}
}

func (u *UserHandler) Login(ctx *gin.Context) {
	// get the email and password from the request body
	var loginData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := ctx.ShouldBindJSON(&loginData); err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{"error": "Invalid request body"})
	}

	if loginData.Email == "" || loginData.Password == "" {
		ctx.AbortWithStatusJSON(400, gin.H{"error": "Email and password are required"})
		return
	}

	token, err := u.UserService.Login(loginData.Email, loginData.Password)
	if err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"token": token})
}

func (u *UserHandler) Register(ctx *gin.Context) {
	// Get the email and password from the request body
	var registrationData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := ctx.ShouldBindJSON(&registrationData); err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{"error": "Invalid request body"})
	}

	if registrationData.Email == "" || registrationData.Password == "" {
		ctx.AbortWithStatusJSON(400, gin.H{"error": "Email and password are required"})
		return
	}

	token, err := u.UserService.Register(registrationData.Email, registrationData.Password)
	if err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"token": token})
}
