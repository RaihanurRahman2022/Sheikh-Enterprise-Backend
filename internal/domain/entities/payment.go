package entities

import (
	"time"

	"github.com/google/uuid"
)

type PaymentEntityType string

const (
	PaymentEntityTypeSupplier PaymentEntityType = "SUPPLIER"
	PaymentEntityTypeCustomer PaymentEntityType = "CUSTOMER"
)

type Payment struct {
	Base
	Type            PaymentEntityType `gorm:"type:varchar(20);not null" json:"type"`
	SupplierID      *uuid.UUID        `gorm:"type:uuid" json:"supplier_id,omitempty"`
	CustomerID      *uuid.UUID        `gorm:"type:uuid" json:"customer_id,omitempty"`
	Amount          float64           `gorm:"type:decimal(10,2);not null" json:"amount"`
	PaymentDateTime time.Time         `gorm:"not null" json:"payment_datetime"`
	Remarks         string            `gorm:"type:text" json:"remarks"`

	// Relations
	Supplier *Supplier `gorm:"foreignKey:SupplierID" json:"supplier,omitempty"`
	Customer *Customer `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
}
