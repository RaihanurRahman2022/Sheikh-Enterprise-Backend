package entities

import (
	"time"

	"github.com/google/uuid"
)

type SalesType string

const (
	SalesTypeRetail    SalesType = "retail"
	SalesTypeWholesale SalesType = "wholesale"
)

type Product struct {
	ID             uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Code           string     `json:"code" gorm:"uniqueIndex;not null"`
	Name           string     `json:"name" gorm:"not null"`
	Style          string     `json:"style" gorm:"not null"`
	MasterCategory string     `json:"master_category" gorm:"not null"`
	SubCategory    string     `json:"sub_category" gorm:"not null"`
	Color          string     `json:"color" gorm:"not null"`
	Size           string     `json:"size" gorm:"not null"`
	PurchasePrice  float64    `json:"purchase_price" gorm:"not null"`
	SalesPrice     float64    `json:"sales_price" gorm:"not null"`
	SalesType      SalesType  `json:"sales_type" gorm:"type:varchar(20);not null"`
	ShopID         uuid.UUID  `json:"shop_id" gorm:"type:uuid;not null"`
	Remarks        string     `json:"remarks" gorm:"type:text"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}
