package usecases

import (
	"Sheikh-Enterprise-Backend/internal/domain/entities"
	repository "Sheikh-Enterprise-Backend/internal/infrastructure/persistence"

	"github.com/google/uuid"
)

type CompanyService interface {
	GetCompanies(page, pageSize int, filters map[string]interface{}, sorts []string) ([]entities.Company, int64, error)
	GetCompanyByID(id uuid.UUID) (*entities.Company, error)
	CreateCompany(company *entities.Company) error
	UpdateCompany(company *entities.Company) error
	DeleteCompany(id uuid.UUID) error
}

type companyService struct {
	companyRepo repository.CompanyRepository
}

func NewCompanyService(companyRepo repository.CompanyRepository) CompanyService {
	return &companyService{
		companyRepo: companyRepo,
	}
}

func (s *companyService) GetCompanies(page, pageSize int, filters map[string]interface{}, sorts []string) ([]entities.Company, int64, error) {
	return s.companyRepo.GetCompaniesWithFilters(filters, sorts, page, pageSize)
}

func (s *companyService) GetCompanyByID(id uuid.UUID) (*entities.Company, error) {
	return s.companyRepo.GetByID(id)
}

func (s *companyService) CreateCompany(company *entities.Company) error {
	return s.companyRepo.Create(company)
}

func (s *companyService) UpdateCompany(company *entities.Company) error {
	return s.companyRepo.Update(company)
}

func (s *companyService) DeleteCompany(id uuid.UUID) error {
	return s.companyRepo.Delete(id)
}
