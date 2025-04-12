
package persistence

import (
	"Sheikh-Enterprise-Backend/internal/domain/entities"
	"gorm.io/gorm"
)

type StockTransferRepository interface {
	BaseRepository[entities.StockTransfer]
	GetStockTransfersWithFilters(filters map[string]interface{}, sorts []string, page, pageSize int) ([]entities.StockTransfer, int64, error)
}

type stockTransferRepository struct {
	BaseRepositoryImpl[entities.StockTransfer]
}

func NewStockTransferRepository(db *gorm.DB) StockTransferRepository {
	return &stockTransferRepository{
		BaseRepositoryImpl: BaseRepositoryImpl[entities.StockTransfer]{DB: db},
	}
}

func (r *stockTransferRepository) GetStockTransfersWithFilters(filters map[string]interface{}, sorts []string, page, pageSize int) ([]entities.StockTransfer, int64, error) {
	var transfers []entities.StockTransfer
	var total int64

	query := r.DB.Model(&entities.StockTransfer{}).
		Preload("Product").
		Preload("ToShop").
		Preload("TransferredBy").
		Where("is_marked_to_delete = ?", false)

	// Apply filters
	for field, value := range filters {
		switch field {
		case "product_id", "to_shop_id", "transferred_by_id":
			query = query.Where(field+" = ?", value)
		case "date_from":
			query = query.Where("transfer_datetime >= ?", value)
		case "date_to":
			query = query.Where("transfer_datetime <= ?", value)
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
	err = query.Offset(offset).Limit(pageSize).Find(&transfers).Error
	if err != nil {
		return nil, 0, err
	}

	return transfers, total, nil
}
