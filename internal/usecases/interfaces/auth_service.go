package interfaces

import (
	"errors"

	models "Sheikh-Enterprise-Backend/internal/domain/entities"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserInactive       = errors.New("user is inactive")
)

// AuthService defines the interface for authentication operations
type AuthService interface {
	Login(username, password string) (string, error)
	Register(req *models.RegisterRequest) error
	ChangePassword(userID string, oldPassword, newPassword string) error
}
