package entities

import (
	"time"

	"github.com/google/uuid"
)

type SalesInvoice struct {
	Base
	ShopID       uuid.UUID  `gorm:"type:uuid;not null" json:"shop_id"`
	CustomerID   *uuid.UUID `gorm:"type:uuid" json:"customer_id,omitempty"`
	SalesByID    uuid.UUID  `gorm:"type:uuid;not null" json:"sales_by_id"`
	SaleDateTime time.Time  `gorm:"not null" json:"sale_datetime"`
	Total        float64    `gorm:"type:decimal(10,2);not null" json:"total"`
	Discount     float64    `gorm:"type:decimal(10,2);not null;default:0" json:"discount"`
	DiscountByID *uuid.UUID `gorm:"type:uuid" json:"discount_by_id,omitempty"`
	Remarks      string     `gorm:"type:text" json:"remarks"`

	// Relations
	Shop         *Shop         `gorm:"foreignKey:ShopID" json:"shop,omitempty"`
	Customer     *Customer     `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	SalesBy      *User         `gorm:"foreignKey:SalesByID" json:"sales_by,omitempty"`
	DiscountBy   *User         `gorm:"foreignKey:DiscountByID" json:"discount_by,omitempty"`
	SalesDetails []SalesDetail `gorm:"foreignKey:InvoiceID" json:"sales_details,omitempty"`
}

type SalesDetail struct {
	Base
	InvoiceID  uuid.UUID `gorm:"type:uuid;not null" json:"invoice_id"`
	ProductID  uuid.UUID `gorm:"type:uuid;not null" json:"product_id"`
	Quantity   int       `gorm:"not null" json:"quantity"`
	SalesPrice float64   `gorm:"type:decimal(10,2);not null" json:"sales_price"`
	Subtotal   float64   `gorm:"type:decimal(10,2);not null" json:"subtotal"`

	// Relations
	SalesInvoice *SalesInvoice `gorm:"foreignKey:InvoiceID" json:"sales_invoice,omitempty"`
	Product      *Product      `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}
