package entities

import (
	"github.com/google/uuid"
)

// Inventory represents the stock level of a product in a shop
type Inventory struct {
	Base
	ProductID uuid.UUID `json:"product_id" gorm:"type:uuid;not null"`
	ShopID    uuid.UUID `json:"shop_id" gorm:"type:uuid;not null"`
	Quantity  int       `json:"quantity" gorm:"not null;default:0"`
	Product   *Product  `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	Shop      *Shop     `json:"shop,omitempty" gorm:"foreignKey:ShopID"`
}
