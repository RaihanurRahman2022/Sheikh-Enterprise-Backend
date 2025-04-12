package entities

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	RoleAdmin   UserRole = "admin"
	RoleManager UserRole = "manager"
	RoleStaff   UserRole = "staff"
)

type User struct {
	ID        uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Username  string     `json:"username" gorm:"uniqueIndex;not null"`
	Password  string     `json:"-" gorm:"not null"`
	Email     string     `json:"email" gorm:"uniqueIndex;not null"`
	Phone     string     `json:"phone" gorm:"not null"`
	FirstName string     `json:"first_name" gorm:"not null"`
	LastName  string     `json:"last_name" gorm:"not null"`
	Role      UserRole   `json:"role" gorm:"type:varchar(20);not null"`
	ShopID    *uuid.UUID `json:"shop_id,omitempty" gorm:"type:uuid"`
	Active    bool       `json:"active" gorm:"not null;default:true"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}
