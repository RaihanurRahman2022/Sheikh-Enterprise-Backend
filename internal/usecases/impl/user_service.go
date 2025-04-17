package usecases

import (
	"errors"

	"Sheikh-Enterprise-Backend/internal/domain/entities"
	repository "Sheikh-Enterprise-Backend/internal/infrastructure/persistence"

	"github.com/google/uuid"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type UserService interface {
	GetUserByID(id uuid.UUID) (*entities.User, error)
	UpdateUser(user *entities.User) error
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) GetUserByID(id uuid.UUID) (*entities.User, error) {
	return s.userRepo.GetByID(id)
}

func (s *userService) UpdateUser(user *entities.User) error {
	return s.userRepo.UpdateDetails(user)
}
