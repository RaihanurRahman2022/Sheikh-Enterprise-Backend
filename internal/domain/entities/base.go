package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Base struct {
	ID               uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	CreatedAt        time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt        time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	IsMarkedToDelete bool           `gorm:"not null;default:false" json:"is_marked_to_delete"`
}
