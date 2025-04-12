package entities

import "github.com/google/uuid"

type Inventory struct {
	Base
	ProductID uuid.UUID `gorm:"type:uuid;not null" json:"product_id"`
	ShopID    uuid.UUID `gorm:"type:uuid" json:"shop_id"`
	Quantity  int       `gorm:"not null;default:0" json:"quantity"`
	Product   *Product  `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Shop      *Shop     `gorm:"foreignKey:ShopID" json:"shop,omitempty"`
}
