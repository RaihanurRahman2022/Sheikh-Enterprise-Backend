package persistence

import (
	"Sheikh-Enterprise-Backend/internal/domain/entities"

	"gorm.io/gorm"
)

type CompanyRepository interface {
	BaseRepository[entities.Company]
	GetCompaniesWithFilters(filters map[string]interface{}, sorts []string, page, pageSize int) ([]entities.Company, int64, error)
}

type companyRepository struct {
	BaseRepositoryImpl[entities.Company]
}

func NewCompanyRepository(db *gorm.DB) CompanyRepository {
	return &companyRepository{
		BaseRepositoryImpl: BaseRepositoryImpl[entities.Company]{DB: db},
	}
}

func (r *companyRepository) GetCompaniesWithFilters(filters map[string]interface{}, sorts []string, page, pageSize int) ([]entities.Company, int64, error) {
	var companies []entities.Company
	var total int64

	query := r.DB.Model(&entities.Company{}).Where("is_marked_to_delete = ?", false)

	// Apply filters
	for field, value := range filters {
		switch field {
		case "name", "phone", "email":
			query = query.Where(field+" LIKE ?", "%"+value.(string)+"%")
		}
	}

	// Count total before pagination
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Apply sorting
	for _, sort := range sorts {
		query = query.Order(sort)
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	err = query.Offset(offset).Limit(pageSize).Find(&companies).Error
	if err != nil {
		return nil, 0, err
	}

	return companies, total, nil
}
