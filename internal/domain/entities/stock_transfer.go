package entities

import (
	"time"

	"github.com/google/uuid"
)

// StockTransferStatus represents the status of a stock transfer
type StockTransferStatus string

const (
	StatusPending   StockTransferStatus = "PENDING"
	StatusApproved  StockTransferStatus = "APPROVED"
	StatusInTransit StockTransferStatus = "IN_TRANSIT"
	StatusCompleted StockTransferStatus = "COMPLETED"
	StatusCancelled StockTransferStatus = "CANCELLED"
	StatusRejected  StockTransferStatus = "REJECTED"
)

// StockTransfer represents a transfer of stock between shops
type StockTransfer struct {
	ID               uuid.UUID           `gorm:"type:varchar(36);primary_key" json:"id"`
	FromShopID       uuid.UUID           `gorm:"type:varchar(36)" json:"from_shop_id"`
	ToShopID         uuid.UUID           `gorm:"type:varchar(36)" json:"to_shop_id"`
	ProductID        uuid.UUID           `gorm:"type:varchar(36)" json:"product_id"`
	Quantity         int                 `json:"quantity"`
	Status           StockTransferStatus `json:"status" gorm:"type:varchar(20);default:'PENDING'"`
	TransferDateTime time.Time           `json:"transfer_datetime"`
	TransferredBy    uuid.UUID           `gorm:"type:varchar(36)" json:"transferred_by"`
	Remarks          string              `json:"remarks"`
	CreatedAt        time.Time           `json:"created_at"`
	UpdatedAt        time.Time           `json:"updated_at"`
	ApprovedBy       *uuid.UUID          `json:"approved_by" gorm:"type:char(36)"`
	ApprovedAt       *time.Time          `json:"approved_at"`
	CompletedBy      *uuid.UUID          `json:"completed_by" gorm:"type:char(36)"`
	CompletedAt      *time.Time          `json:"completed_at"`
	RejectedBy       *uuid.UUID          `json:"rejected_by" gorm:"type:char(36)"`
	RejectedAt       *time.Time          `json:"rejected_at"`
	RejectionReason  string              `json:"rejection_reason"`

	// Relationships
	FromShop          *Shop    `gorm:"foreignKey:FromShopID" json:"from_shop,omitempty"`
	ToShop            *Shop    `gorm:"foreignKey:ToShopID" json:"to_shop,omitempty"`
	Product           *Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	TransferredByUser *User    `gorm:"foreignKey:TransferredBy" json:"transferred_by_user,omitempty"`
}

// StockTransferHistory represents the history of status changes for a stock transfer
type StockTransferHistory struct {
	Base
	StockTransferID uuid.UUID           `json:"stock_transfer_id" gorm:"type:char(36)"`
	Status          StockTransferStatus `json:"status" gorm:"type:varchar(20)"`
	ChangedBy       uuid.UUID           `json:"changed_by" gorm:"type:char(36)"`
	Remarks         string              `json:"remarks"`
	StockTransfer   StockTransfer       `json:"stock_transfer" gorm:"foreignKey:StockTransferID"`
}
