package usecases

// Services groups all service instances
type Services struct {
	Auth     AuthService
	User     UserService
	Product  ProductService
	Sales    SalesService
	Purchase PurchaseService
	Supplier SupplierService
	Company  CompanyService
	Shop     ShopService
}
