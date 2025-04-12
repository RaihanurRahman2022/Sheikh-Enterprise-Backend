
package impl

import (
	"Sheikh-Enterprise-Backend/internal/domain/entities"
	"Sheikh-Enterprise-Backend/internal/infrastructure/persistence"
	"github.com/google/uuid"
)

type CustomerService struct {
	customerRepo persistence.CustomerRepository
}

func NewCustomerService(customerRepo persistence.CustomerRepository) *CustomerService {
	return &CustomerService{customerRepo: customerRepo}
}

func (s *CustomerService) GetAllCustomers() ([]entities.Customer, error) {
	return s.customerRepo.GetAll()
}

func (s *CustomerService) GetCustomerByID(id uuid.UUID) (*entities.Customer, error) {
	return s.customerRepo.GetByID(id)
}

func (s *CustomerService) CreateCustomer(customer *entities.Customer) error {
	return s.customerRepo.Create(customer)
}

func (s *CustomerService) UpdateCustomer(customer *entities.Customer) error {
	return s.customerRepo.Update(customer)
}

func (s *CustomerService) DeleteCustomer(id uuid.UUID) error {
	return s.customerRepo.Delete(id)
}
