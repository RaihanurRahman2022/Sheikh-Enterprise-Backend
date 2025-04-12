
package usecases

import (
	"Sheikh-Enterprise-Backend/internal/domain/entities"
	repository "Sheikh-Enterprise-Backend/internal/infrastructure/persistence"
	"github.com/google/uuid"
)

type SupplierService interface {
	GetSuppliers(page, pageSize int, filters map[string]interface{}, sorts []string) ([]entities.Supplier, int64, error)
	GetSupplierByID(id uuid.UUID) (*entities.Supplier, error)
	CreateSupplier(supplier *entities.Supplier) error
	UpdateSupplier(supplier *entities.Supplier) error
	DeleteSupplier(id uuid.UUID) error
}

type supplierService struct {
	supplierRepo repository.SupplierRepository
}

func NewSupplierService(supplierRepo repository.SupplierRepository) SupplierService {
	return &supplierService{
		supplierRepo: supplierRepo,
	}
}

func (s *supplierService) GetSuppliers(page, pageSize int, filters map[string]interface{}, sorts []string) ([]entities.Supplier, int64, error) {
	return s.supplierRepo.GetSuppliersWithFilters(filters, sorts, page, pageSize)
}

func (s *supplierService) GetSupplierByID(id uuid.UUID) (*entities.Supplier, error) {
	return s.supplierRepo.GetByID(id)
}

func (s *supplierService) CreateSupplier(supplier *entities.Supplier) error {
	return s.supplierRepo.Create(supplier)
}

func (s *supplierService) UpdateSupplier(supplier *entities.Supplier) error {
	return s.supplierRepo.Update(supplier)
}

func (s *supplierService) DeleteSupplier(id uuid.UUID) error {
	return s.supplierRepo.Delete(id)
}
