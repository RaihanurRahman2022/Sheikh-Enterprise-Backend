package usecases

import (
	"errors"

	"Sheikh-Enterprise-Backend/internal/domain/entities"
	repository "Sheikh-Enterprise-Backend/internal/infrastructure/persistence"
	"Sheikh-Enterprise-Backend/pkg/utils"

	"github.com/google/uuid"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type UserService interface {
	GetUserDetails(id uuid.UUID) (*entities.User, error)
	UpdateUserDetails(user *entities.User) error
	UpdateUserPassword(id uuid.UUID, oldPassword, newPassword string) error
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) GetUserDetails(id uuid.UUID) (*entities.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (s *userService) UpdateUserDetails(user *entities.User) error {
	return s.userRepo.UpdateDetails(user)
}

func (s *userService) UpdateUserPassword(id uuid.UUID, oldPassword, newPassword string) error {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return ErrUserNotFound
	}

	if !utils.CheckPassword(oldPassword, user.Password) {
		return ErrInvalidCredentials
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	return s.userRepo.UpdatePassword(id, hashedPassword)
}
