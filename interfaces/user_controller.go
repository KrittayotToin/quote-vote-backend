package iface

import (
	"github.com/KrittayotToin/quote-vote-backend/model"
)

// UserInterface interface defines methods for interacting with user data
type UserInterface interface {
	Create(user model.User) (model.User, error)
}
