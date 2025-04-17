package handlers

// Handlers groups all HTTP handlers
type Handlers struct {
	Auth     *AuthHandler
	User     *UserHandler
	Product  *ProductHandler
	Sales    *SalesHandler
	Purchase *PurchaseHandler
	Supplier *SupplierHandler
	Company  *CompanyHandler
	Shop     *ShopHandler
}
