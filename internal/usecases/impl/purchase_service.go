
package usecases

import (
	"Sheikh-Enterprise-Backend/internal/domain/entities"
	repository "Sheikh-Enterprise-Backend/internal/infrastructure/persistence"
	"github.com/google/uuid"
)

type PurchaseService interface {
	GetPurchases(page, pageSize int, filters map[string]interface{}, sorts []string) ([]entities.PurchaseInvoice, int64, error)
	GetPurchaseByID(id uuid.UUID) (*entities.PurchaseInvoice, error)
	CreatePurchase(purchase *entities.PurchaseInvoice) error
	DeletePurchase(id uuid.UUID) error
}

type purchaseService struct {
	purchaseRepo repository.PurchaseRepository
}

func NewPurchaseService(purchaseRepo repository.PurchaseRepository) PurchaseService {
	return &purchaseService{
		purchaseRepo: purchaseRepo,
	}
}

func (s *purchaseService) GetPurchases(page, pageSize int, filters map[string]interface{}, sorts []string) ([]entities.PurchaseInvoice, int64, error) {
	return s.purchaseRepo.GetPurchasesWithFilters(filters, sorts, page, pageSize)
}

func (s *purchaseService) GetPurchaseByID(id uuid.UUID) (*entities.PurchaseInvoice, error) {
	return s.purchaseRepo.GetByID(id)
}

func (s *purchaseService) CreatePurchase(purchase *entities.PurchaseInvoice) error {
	// Calculate total
	var total float64
	for i := range purchase.PurchaseDetails {
		detail := &purchase.PurchaseDetails[i]
		total += detail.PurchasePrice * float64(detail.Quantity)
	}
	purchase.Total = total
	return s.purchaseRepo.Create(purchase)
}

func (s *purchaseService) DeletePurchase(id uuid.UUID) error {
	return s.purchaseRepo.Delete(id)
}
