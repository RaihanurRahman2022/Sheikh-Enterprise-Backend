package usecases

import (
	"errors"

	"Sheikh-Enterprise-Backend/internal/domain/entities"
	repository "Sheikh-Enterprise-Backend/internal/infrastructure/persistence"
	"Sheikh-Enterprise-Backend/pkg/utils"

	"github.com/google/uuid"
)

var (
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrUserInactive       = errors.New("user is inactive")
	ErrInvalidToken       = errors.New("invalid token")
)

type AuthService interface {
	Login(username, password string) (string, string, error)
	ChangePassword(userID string, oldPassword, newPassword string) error
	Register(req *entities.RegisterRequest) error
	RefreshToken(refreshToken string) (string, string, error)
}

type authService struct {
	authRepo repository.AuthRepository
}

func NewAuthService(authRepo repository.AuthRepository) AuthService {
	return &authService{
		authRepo: authRepo,
	}
}

func (s *authService) Login(username, password string) (string, string, error) {
	user, err := s.authRepo.GetUserByUsername(username)
	if err != nil {
		return "", "", ErrInvalidCredentials
	}

	if !user.Active {
		return "", "", ErrUserInactive
	}

	if !utils.CheckPassword(password, user.Password) {
		return "", "", ErrInvalidCredentials
	}

	// Generate access token (24 hours)
	accessToken, err := utils.GenerateJWT(user.ID, user.Username, string(user.Role), user.ShopID)
	if err != nil {
		return "", "", err
	}

	// Generate refresh token (7 days)
	refreshToken, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *authService) ChangePassword(userID string, oldPassword, newPassword string) error {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return err
	}

	user, err := s.authRepo.GetUserByID(uid)
	if err != nil {
		return err
	}

	if !utils.CheckPassword(oldPassword, user.Password) {
		return ErrInvalidCredentials
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	return s.authRepo.UpdatePassword(uid, hashedPassword)
}

func (s *authService) Register(req *entities.RegisterRequest) error {
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return err
	}

	user := &entities.User{
		Username:  req.Username,
		Password:  hashedPassword,
		Email:     req.Email,
		Phone:     req.Phone,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      entities.UserRole(req.Role),
		Active:    true,
	}

	if req.ShopID != "" {
		shopID, err := uuid.Parse(req.ShopID)
		if err != nil {
			return err
		}
		user.ShopID = &shopID
	}

	return s.authRepo.Create(user)
}

func (s *authService) RefreshToken(refreshToken string) (string, string, error) {
	claims, err := utils.ValidateRefreshToken(refreshToken)
	if err != nil {
		return "", "", ErrInvalidToken
	}

	user, err := s.authRepo.GetUserByID(claims.UserID)
	if err != nil {
		return "", "", ErrInvalidToken
	}

	if !user.Active {
		return "", "", ErrUserInactive
	}

	// Generate new access token (24 hours)
	accessToken, err := utils.GenerateJWT(user.ID, user.Username, string(user.Role), user.ShopID)
	if err != nil {
		return "", "", err
	}

	// Generate new refresh token (7 days)
	newRefreshToken, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		return "", "", err
	}

	return accessToken, newRefreshToken, nil
}
