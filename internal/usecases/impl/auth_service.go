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
)

type AuthService interface {
	Login(username, password string) (string, error)
	ChangePassword(userID string, oldPassword, newPassword string) error
	Register(req *entities.RegisterRequest) error
}

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{
		userRepo: userRepo,
	}
}

func (s *authService) Login(username, password string) (string, error) {
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	if !user.Active {
		return "", ErrUserInactive
	}

	if !utils.CheckPassword(password, user.Password) {
		return "", ErrInvalidCredentials
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.ID, user.Username, string(user.Role), user.ShopID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *authService) ChangePassword(userID string, oldPassword, newPassword string) error {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return err
	}
	user, err := s.userRepo.GetByID(uid)
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

	return s.userRepo.UpdatePassword(user.ID, hashedPassword)
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

	return s.userRepo.Create(user)
}
