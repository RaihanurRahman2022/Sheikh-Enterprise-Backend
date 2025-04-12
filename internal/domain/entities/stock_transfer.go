package entities

import (
	"time"

	"github.com/google/uuid"
)

type StockTransfer struct {
	Base
	ProductID        uuid.UUID `gorm:"type:uuid;not null" json:"product_id"`
	ToShopID         uuid.UUID `gorm:"type:uuid;not null" json:"to_shop_id"`
	Quantity         int       `gorm:"not null" json:"quantity"`
	TransferDateTime time.Time `gorm:"not null" json:"transfer_datetime"`
	TransferredByID  uuid.UUID `gorm:"type:uuid;not null" json:"transferred_by_id"`

	// Relations
	Product       *Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	ToShop        *Shop    `gorm:"foreignKey:ToShopID" json:"to_shop,omitempty"`
	TransferredBy *User    `gorm:"foreignKey:TransferredByID" json:"transferred_by,omitempty"`
}
