package repository

import (
	"github.com/KrittayotToin/quote-vote-backend/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (repo *UserRepository) Create(user model.User) (model.User, error) {
	// oldPassword := user.Passwordhash
	hashedPassword, err := hashPassword(user.PasswordHash)
	user.PasswordHash = hashedPassword
	if err != nil {
		return model.User{}, err
	}
	if err := repo.DB.Create(&user).Error; err != nil {
		return model.User{}, err
	}
	return user, nil
}
