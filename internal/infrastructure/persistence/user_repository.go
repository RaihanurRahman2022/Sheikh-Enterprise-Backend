package persistence

import (
	"Sheikh-Enterprise-Backend/internal/domain/entities"

	"gorm.io/gorm"
)

type UserRepository interface {
	BaseRepository[entities.User]
	UpdateDetails(user *entities.User) error
}

type userRepository struct {
	BaseRepositoryImpl[entities.User]
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		BaseRepositoryImpl: BaseRepositoryImpl[entities.User]{DB: db},
	}
}

func (r *userRepository) UpdateDetails(user *entities.User) error {
	return r.DB.Model(&entities.User{}).
		Where("id = ?", user.ID).
		Updates(map[string]interface{}{
			"first_name": user.FirstName,
			"last_name":  user.LastName,
			"email":      user.Email,
			"phone":      user.Phone,
			"shop_id":    user.ShopID,
		}).Error
}
