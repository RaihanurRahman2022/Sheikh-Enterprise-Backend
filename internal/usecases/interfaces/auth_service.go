package interfaces

import (
	"errors"

	models "Sheikh-Enterprise-Backend/internal/domain/entities"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserInactive       = errors.New("user is inactive")
	ErrInvalidToken       = errors.New("invalid token")
)

// AuthService defines the interface for authentication operations
type AuthService interface {
	Login(username, password string) (string, string, error) // Returns access token and refresh token
	Register(req *models.RegisterRequest) error
	ChangePassword(userID string, oldPassword, newPassword string) error
	RefreshToken(refreshToken string) (string, string, error) // Returns new access token and refresh token
}
