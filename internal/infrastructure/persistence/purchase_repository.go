
package persistence

import (
	"Sheikh-Enterprise-Backend/internal/domain/entities"
	"gorm.io/gorm"
)

type PurchaseRepository interface {
	BaseRepository[entities.PurchaseInvoice]
	GetPurchasesWithFilters(filters map[string]interface{}, sorts []string, page, pageSize int) ([]entities.PurchaseInvoice, int64, error)
}

type purchaseRepository struct {
	BaseRepositoryImpl[entities.PurchaseInvoice]
}

func NewPurchaseRepository(db *gorm.DB) PurchaseRepository {
	return &purchaseRepository{
		BaseRepositoryImpl: BaseRepositoryImpl[entities.PurchaseInvoice]{DB: db},
	}
}

func (r *purchaseRepository) GetPurchasesWithFilters(filters map[string]interface{}, sorts []string, page, pageSize int) ([]entities.PurchaseInvoice, int64, error) {
	var purchases []entities.PurchaseInvoice
	var total int64

	query := r.DB.Model(&entities.PurchaseInvoice{}).
		Preload("Supplier").
		Preload("EntryBy").
		Preload("PurchaseDetails").
		Preload("PurchaseDetails.Product").
		Where("is_marked_to_delete = ?", false)

	// Apply filters
	for field, value := range filters {
		switch field {
		case "supplier_id", "entry_by_id":
			query = query.Where(field+" = ?", value)
		case "payment_type":
			query = query.Where("payment_type = ?", value)
		case "min_total":
			query = query.Where("total >= ?", value)
		case "max_total":
			query = query.Where("total <= ?", value)
		case "date_from":
			query = query.Where("purchase_datetime >= ?", value)
		case "date_to":
			query = query.Where("purchase_datetime <= ?", value)
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
	err = query.Offset(offset).Limit(pageSize).Find(&purchases).Error
	if err != nil {
		return nil, 0, err
	}

	return purchases, total, nil
}
