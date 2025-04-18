package dto

import (
	"Sheikh-Enterprise-Backend/internal/domain/entities"
	"time"

	"github.com/google/uuid"
)

type CompanyDTO struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email"`
	Slogan    string    `json:"slogan"`
	Remarks   string    `json:"remarks"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ToCompanyDTO(company *entities.Company) CompanyDTO {
	return CompanyDTO{
		ID:        company.ID,
		Name:      company.Name,
		Address:   company.Address,
		Phone:     company.Phone,
		Email:     company.Email,
		Slogan:    company.Slogan,
		Remarks:   company.Remarks,
		CreatedAt: company.CreatedAt,
		UpdatedAt: company.UpdatedAt,
	}
}
