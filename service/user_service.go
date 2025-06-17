package service

import (
	iface "github.com/KrittayotToin/quote-vote-backend/interfaces"
	"github.com/KrittayotToin/quote-vote-backend/model"
)

type UserService struct {
	repo iface.UserInterface
}

func (s *UserService) Create(user model.User) (model.User, error) {
	return s.repo.Create(user)
}

func NewUserService(repo iface.UserInterface) iface.UserInterface {
	return &UserService{repo: repo}
}
