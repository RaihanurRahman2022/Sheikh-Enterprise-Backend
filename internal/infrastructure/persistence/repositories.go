package persistence

import "gorm.io/gorm"

// Repositories groups all repository instances
type Repositories struct {
	User      UserRepository
	Product   ProductRepository
	Sales     SalesRepository
	Purchase  PurchaseRepository
	Supplier  SupplierRepository
	Company   CompanyRepository
	Shop      ShopRepository
	Auth      AuthRepository
	Analytics AnalyticsRepository
	Customer  CustomerRepository
	Inventory InventoryRepository
}

func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		User:      NewUserRepository(db),
		Product:   NewProductRepository(db),
		Sales:     NewSalesRepository(db),
		Purchase:  NewPurchaseRepository(db),
		Supplier:  NewSupplierRepository(db),
		Company:   NewCompanyRepository(db),
		Shop:      NewShopRepository(db),
		Auth:      NewAuthRepository(db),
		Analytics: NewAnalyticsRepository(db),
		Customer:  NewCustomerRepository(db),
		Inventory: NewInventoryRepository(db),
	}
}
