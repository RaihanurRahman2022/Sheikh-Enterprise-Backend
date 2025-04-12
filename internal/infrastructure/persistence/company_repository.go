
package persistence

import (
	"Sheikh-Enterprise-Backend/internal/domain/entities"
	"gorm.io/gorm"
)

type CompanyRepository interface {
	BaseRepository[entities.Company]
}

type companyRepository struct {
	BaseRepositoryImpl[entities.Company]
}

func NewCompanyRepository(db *gorm.DB) CompanyRepository {
	return &companyRepository{
		BaseRepositoryImpl: BaseRepositoryImpl[entities.Company]{DB: db},
	}
}
