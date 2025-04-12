
package persistence

import (
	"Sheikh-Enterprise-Backend/internal/domain/entities"
	"gorm.io/gorm"
)

type CustomerRepository interface {
	BaseRepository[entities.Customer]
}

type customerRepository struct {
	BaseRepositoryImpl[entities.Customer]
}

func NewCustomerRepository(db *gorm.DB) CustomerRepository {
	return &customerRepository{
		BaseRepositoryImpl: BaseRepositoryImpl[entities.Customer]{DB: db},
	}
}
