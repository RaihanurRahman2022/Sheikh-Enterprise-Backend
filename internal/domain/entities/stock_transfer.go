package entities

import (
	"time"

	"github.com/google/uuid"
)

type StockTransferStatus string

const (
	StatusPending   StockTransferStatus = "PENDING"
	StatusApproved  StockTransferStatus = "APPROVED"
	StatusInTransit StockTransferStatus = "IN_TRANSIT"
	StatusCompleted StockTransferStatus = "COMPLETED"
	StatusCancelled StockTransferStatus = "CANCELLED"
	StatusRejected  StockTransferStatus = "REJECTED"
)

type StockTransfer struct {
	ID                uuid.UUID           `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	FromShopID        *uuid.UUID          `gorm:"type:uuid" json:"from_shop_id"` // Nullable for central stock
	ToShopID          uuid.UUID           `gorm:"type:uuid;not null" json:"to_shop_id"`
	ProductID         uuid.UUID           `gorm:"type:uuid;not null" json:"product_id"`
	Quantity          int                 `json:"quantity" gorm:"check:quantity > 0"`
	Status            StockTransferStatus `json:"status" gorm:"type:varchar(20);default:'PENDING'"`
	TransferDateTime  time.Time           `json:"transfer_datetime"`
	TransferredBy     uuid.UUID           `gorm:"type:uuid;not null" json:"transferred_by"`
	Remarks           string              `json:"remarks"`
	CreatedAt         time.Time           `json:"created_at"`
	UpdatedAt         time.Time           `json:"updated_at"`
	ApprovedBy        *uuid.UUID          `gorm:"type:uuid" json:"approved_by"`
	ApprovedAt        *time.Time          `json:"approved_at"`
	CompletedBy       *uuid.UUID          `gorm:"type:uuid" json:"completed_by"`
	CompletedAt       *time.Time          `json:"completed_at"`
	RejectedBy        *uuid.UUID          `gorm:"type:uuid" json:"rejected_by"`
	RejectedAt        *time.Time          `json:"rejected_at"`
	RejectionReason   string              `json:"rejection_reason"`
	FromShop          *Shop               `gorm:"foreignKey:FromShopID;references:ShopID" json:"from_shop,omitempty"`
	ToShop            *Shop               `gorm:"foreignKey:ToShopID;references:ShopID" json:"to_shop,omitempty"`
	Product           *Product            `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	TransferredByUser *User               `gorm:"foreignKey:TransferredBy" json:"transferred_by_user,omitempty"`
}

type StockTransferHistory struct {
	Base
	StockTransferID uuid.UUID           `json:"stock_transfer_id" gorm:"type:uuid"`
	Status          StockTransferStatus `json:"status" gorm:"type:varchar(20)"`
	ChangedBy       uuid.UUID           `json:"changed_by" gorm:"type:uuid"`
	Remarks         string              `json:"remarks"`
	StockTransfer   StockTransfer       `gorm:"foreignKey:StockTransferID" json:"stock_transfer"`
}
