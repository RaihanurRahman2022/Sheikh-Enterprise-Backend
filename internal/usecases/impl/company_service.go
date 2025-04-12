
package impl

import (
	"Sheikh-Enterprise-Backend/internal/domain/entities"
	"Sheikh-Enterprise-Backend/internal/infrastructure/persistence"
	"github.com/google/uuid"
)

type CompanyService struct {
	companyRepo persistence.CompanyRepository
}

func NewCompanyService(companyRepo persistence.CompanyRepository) *CompanyService {
	return &CompanyService{companyRepo: companyRepo}
}

func (s *CompanyService) GetAllCompanies() ([]entities.Company, error) {
	return s.companyRepo.GetAll()
}

func (s *CompanyService) GetCompanyByID(id uuid.UUID) (*entities.Company, error) {
	return s.companyRepo.GetByID(id)
}

func (s *CompanyService) CreateCompany(company *entities.Company) error {
	return s.companyRepo.Create(company)
}

func (s *CompanyService) UpdateCompany(company *entities.Company) error {
	return s.companyRepo.Update(company)
}

func (s *CompanyService) DeleteCompany(id uuid.UUID) error {
	return s.companyRepo.Delete(id)
}
