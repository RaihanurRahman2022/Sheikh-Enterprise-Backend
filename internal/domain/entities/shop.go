package entities

import "github.com/google/uuid"

type Shop struct {
	Base
	ShopID       uuid.UUID `gorm:"type:uuid;primary_key" json:"shop_id"`
	CompanyID    uuid.UUID `gorm:"type:uuid" json:"company_id"`
	Name         string    `json:"name" binding:"required"`
	Address      string    `json:"address" binding:"required"`
	Phone        string    `json:"phone" binding:"required"`
	Email        string    `json:"email" binding:"required,email"`
	ManagerName  string    `json:"manager_name" binding:"required"`
	ManagerPhone string    `json:"manager_phone" binding:"required"`
	Remarks      string    `json:"remarks"`
	Company      Company   `gorm:"foreignKey:CompanyID;references:ID" json:"company,omitempty"`
}
