package usecases

import (
	"Sheikh-Enterprise-Backend/internal/domain/entities"
	repository "Sheikh-Enterprise-Backend/internal/infrastructure/persistence"

	"github.com/google/uuid"
)

type ShopService interface {
	GetShops(page, pageSize int, filters map[string]interface{}, sorts []string) ([]entities.Shop, int64, error)
	GetShopByID(id uuid.UUID) (*entities.Shop, error)
	GetShopsByCompanyID(companyID uuid.UUID) ([]entities.Shop, error)
	CreateShop(shop *entities.Shop) error
	UpdateShop(shop *entities.Shop) error
	DeleteShop(id uuid.UUID) error
}

type shopService struct {
	shopRepo repository.ShopRepository
}

func NewShopService(shopRepo repository.ShopRepository) ShopService {
	return &shopService{
		shopRepo: shopRepo,
	}
}

func (s *shopService) GetShops(page, pageSize int, filters map[string]interface{}, sorts []string) ([]entities.Shop, int64, error) {
	return s.shopRepo.GetShopsWithFilters(filters, sorts, page, pageSize)
}

func (s *shopService) GetShopByID(id uuid.UUID) (*entities.Shop, error) {
	return s.shopRepo.GetByID(id)
}

func (s *shopService) GetShopsByCompanyID(companyID uuid.UUID) ([]entities.Shop, error) {
	return s.shopRepo.GetShopsByCompanyID(companyID.String())
}

func (s *shopService) CreateShop(shop *entities.Shop) error {
	return s.shopRepo.Create(shop)
}

func (s *shopService) UpdateShop(shop *entities.Shop) error {
	return s.shopRepo.Update(shop)
}

func (s *shopService) DeleteShop(id uuid.UUID) error {
	return s.shopRepo.Delete(id)
}
