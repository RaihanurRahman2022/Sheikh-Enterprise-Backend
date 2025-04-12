package entities

import (
	"time"

	"github.com/google/uuid"
)

type PaymentType string

const (
	PaymentTypeCash   PaymentType = "CASH"
	PaymentTypeCredit PaymentType = "CREDIT"
)

type PurchaseInvoice struct {
	Base
	SupplierID       uuid.UUID   `gorm:"type:uuid;not null" json:"supplier_id"`
	PurchaseDateTime time.Time   `gorm:"not null" json:"purchase_datetime"`
	Total            float64     `gorm:"type:decimal(10,2);not null" json:"total"`
	PaymentType      PaymentType `gorm:"type:varchar(20);not null" json:"payment_type"`
	EntryByID        uuid.UUID   `gorm:"type:uuid;not null" json:"entry_by_id"`
	Remarks          string      `gorm:"type:text" json:"remarks"`

	// Relations
	Supplier        *Supplier        `gorm:"foreignKey:SupplierID" json:"supplier,omitempty"`
	EntryBy         *User            `gorm:"foreignKey:EntryByID" json:"entry_by,omitempty"`
	PurchaseDetails []PurchaseDetail `gorm:"foreignKey:PurchaseInvoiceID" json:"purchase_details,omitempty"`
}

type PurchaseDetail struct {
	Base
	PurchaseInvoiceID uuid.UUID `gorm:"type:uuid;not null" json:"purchase_invoice_id"`
	ProductID         uuid.UUID `gorm:"type:uuid;not null" json:"product_id"`
	Quantity          int       `gorm:"not null" json:"quantity"`
	PurchasePrice     float64   `gorm:"type:decimal(10,2);not null" json:"purchase_price"`

	// Relations
	PurchaseInvoice *PurchaseInvoice `gorm:"foreignKey:PurchaseInvoiceID" json:"purchase_invoice,omitempty"`
	Product         *Product         `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}
