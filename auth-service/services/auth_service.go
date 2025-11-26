package services

import (
	"auth-service/database"
	"auth-service/models"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(req models.AuthRequest) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := models.User{
		Email:    req.Email,
		Password: string(hashedPassword),
		Name:     req.Name,
		Role:     "customer",
	}

	result := database.DB.Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func LoginUser(req models.AuthRequest) (*models.User, error) {
	var user models.User
	result := database.DB.Where("email = ?", req.Email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	result := database.DB.First(&user, userID)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	result := database.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}