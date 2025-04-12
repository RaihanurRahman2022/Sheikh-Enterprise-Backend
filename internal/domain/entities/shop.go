package entities

import "github.com/google/uuid"

type Shop struct {
	Base
	CompanyID    uuid.UUID   `gorm:"type:uuid;not null" json:"company_id"`
	Name         string      `gorm:"type:varchar(255);not null" json:"name"`
	Address      string      `gorm:"type:text" json:"address"`
	Phone        string      `gorm:"type:varchar(20)" json:"phone"`
	Email        string      `gorm:"type:varchar(255)" json:"email"`
	ManagerName  string      `gorm:"type:varchar(255)" json:"manager_name"`
	ManagerPhone string      `gorm:"type:varchar(20)" json:"manager_phone"`
	Remarks      string      `gorm:"type:text" json:"remarks"`
	Company      *Company    `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Inventories  []Inventory `gorm:"foreignKey:ShopID" json:"inventories,omitempty"`
}
