package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Shop struct {
	ShopID       uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"shop_id"`
	CompanyID    uuid.UUID      `gorm:"type:uuid;not null" json:"company_id"`
	Name         string         `json:"name" binding:"required"`
	Address      string         `json:"address" binding:"required"`
	Phone        string         `json:"phone" binding:"required"`
	Email        string         `json:"email" binding:"required,email"`
	ManagerName  string         `json:"manager_name" binding:"required"`
	ManagerPhone string         `json:"manager_phone" binding:"required"`
	Remarks      string         `json:"remarks"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Company      Company        `gorm:"foreignKey:CompanyID;references:ID" json:"company,omitempty"`
}
