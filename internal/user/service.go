package user

import (
	"context"
	"errors"
	"log"

	"firebase.google.com/go/v4/auth"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserService struct {
	DB       *gorm.DB
	FireAuth *auth.Client
}

func (s *UserService) Login(email, password string) (token string, err error) {
	var user User

	// get user from db
	err = s.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", errors.New("user with email does not exist")
		}
		log.Printf("failed to get user by email from database: %v", err)
		return "", errors.New("internal server error")
	}

	// check password
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("incorrect password")
	}

	// generate firebase custom token
	token, err = s.FireAuth.CustomToken(context.Background(), user.ID)
	if err != nil {
		log.Printf("failed to generate custom token: %v", err)
		return "", errors.New("internal server error")
	}

	return
}

func (s *UserService) Register(email, password string) (token string, err error) {
	var user User

	// check if user with the email already exists
	err = s.DB.Raw("SELECT id, email, password FROM users WHERE email = ?", email).Scan(&user).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			log.Printf("failed to get user by email from database: %v", err)
			return "", errors.New("internal server error")
		}
	}

	if user.ID != "" {
		return "", errors.New("user with email already exists")
	}

	// generate a uuid for the new user
	uid := uuid.New().String()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("failed to hash password: %v", err)
		return "", errors.New("internal server error")
	}

	// create a new user in DB
	user.ID = uid
	user.Email = email
	user.Password = string(hashedPassword)

	err = s.DB.Create(&user).Error
	if err != nil {
		log.Printf("failed to insert user into database: %v", err)
		return "", errors.New("internal server error")
	}

	// create custom token
	token, err = s.FireAuth.CustomToken(context.Background(), uid)
	if err != nil {
		log.Printf("failed to generate custom token: %v", err)
		return "", errors.New("internal server error")
	}

	return
}
