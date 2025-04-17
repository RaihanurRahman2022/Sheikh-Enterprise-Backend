package persistence

import (
	"Sheikh-Enterprise-Backend/internal/domain/entities"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CustomerRepository interface {
	BaseRepository[entities.Customer]
	GetByPhone(phone string) (*entities.Customer, error)
	GetCustomerSales(customerID uuid.UUID) ([]entities.SalesInvoice, error)
}

type customerRepository struct {
	BaseRepositoryImpl[entities.Customer]
}

func NewCustomerRepository(db *gorm.DB) CustomerRepository {
	return &customerRepository{
		BaseRepositoryImpl: BaseRepositoryImpl[entities.Customer]{DB: db},
	}
}

func (r *customerRepository) GetByPhone(phone string) (*entities.Customer, error) {
	var customer entities.Customer
	err := r.DB.Where("phone = ? AND is_marked_to_delete = ?", phone, false).First(&customer).Error
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

func (r *customerRepository) GetCustomerSales(customerID uuid.UUID) ([]entities.SalesInvoice, error) {
	var sales []entities.SalesInvoice
	err := r.DB.Where("customer_id = ?", customerID).Find(&sales).Error
	if err != nil {
		return nil, err
	}
	return sales, nil
}
