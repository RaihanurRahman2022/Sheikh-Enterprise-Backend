package database

import (
	"log"

	"Sheikh-Enterprise-Backend/internal/domain/entities"

	"gorm.io/gorm"
)

// RunMigrations handles all database migrations
func RunMigrations(db *gorm.DB) error {
	log.Println("Running database migrations...")

	// List of entities to migrate
	entities := []interface{}{
		&entities.User{},
		&entities.Company{},
		&entities.Shop{},
		&entities.Product{},
		&entities.Inventory{},
		&entities.Supplier{},
		&entities.Customer{},
		&entities.PurchaseInvoice{},
		&entities.PurchaseDetail{},
		&entities.SalesInvoice{},
		&entities.SalesDetail{},
		&entities.StockTransfer{},
		&entities.Payment{},
	}

	// Run migrations
	for _, model := range entities {
		if err := db.AutoMigrate(model); err != nil {
			return err
		}
	}

	log.Println("Database migrations completed successfully")
	return nil
}
