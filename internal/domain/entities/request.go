package entities

// LoginRequest represents the login request body
type LoginRequest struct {
	Username string `json:"username" binding:"required,username"`
	Password string `json:"password" binding:"required,password"`
}

// RegisterRequest represents the registration request body
type RegisterRequest struct {
	Username  string `json:"username" binding:"required,username"`
	Password  string `json:"password" binding:"required,password"`
	Email     string `json:"email" binding:"required,email"`
	Phone     string `json:"phone" binding:"required,phone"`
	FirstName string `json:"first_name" binding:"required,min=2,max=50"`
	LastName  string `json:"last_name" binding:"required,min=2,max=50"`
	Role      string `json:"role" binding:"required,oneof=admin manager staff"`
	ShopID    string `json:"shop_id" binding:"omitempty,uuid"`
}

// UpdateUserRequest represents the update user request body
type UpdateUserRequest struct {
	Email     string `json:"email" binding:"omitempty,email"`
	Phone     string `json:"phone" binding:"omitempty,phone"`
	FirstName string `json:"first_name" binding:"omitempty,min=2,max=50"`
	LastName  string `json:"last_name" binding:"omitempty,min=2,max=50"`
	Active    *bool  `json:"active" binding:"omitempty"`
}

// ChangePasswordRequest represents the change password request body
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required,password"`
	NewPassword string `json:"new_password" binding:"required,password,nefield=OldPassword"`
}

// CreateProductRequest represents the create product request body
type CreateProductRequest struct {
	Code           string  `json:"code" binding:"required,min=3,max=50"`
	Name           string  `json:"name" binding:"required,min=3,max=100"`
	Style          string  `json:"style" binding:"required,min=2,max=50"`
	MasterCategory string  `json:"master_category" binding:"required,min=2,max=50"`
	SubCategory    string  `json:"sub_category" binding:"required,min=2,max=50"`
	Color          string  `json:"color" binding:"required,min=2,max=30"`
	Size           string  `json:"size" binding:"required,min=1,max=20"`
	PurchasePrice  float64 `json:"purchase_price" binding:"required,min=0"`
	SalesPrice     float64 `json:"sales_price" binding:"required,min=0,gtefield=PurchasePrice"`
	SalesType      string  `json:"sales_type" binding:"required,oneof=retail wholesale"`
	ShopID         string  `json:"shop_id" binding:"required,uuid"`
}

// CreateSaleRequest represents the create sale request body
type CreateSaleRequest struct {
	CustomerID  string               `json:"customer_id" binding:"required,uuid"`
	ShopID      string               `json:"shop_id" binding:"required,uuid"`
	SaleType    string               `json:"sale_type" binding:"required,oneof=retail wholesale"`
	PaymentType string               `json:"payment_type" binding:"required,oneof=cash card mobile"`
	Items       []SaleItemRequest    `json:"items" binding:"required,min=1,dive"`
	Payments    []SalePaymentRequest `json:"payments" binding:"required,min=1,dive"`
	Discount    float64              `json:"discount" binding:"min=0"`
	Note        string               `json:"note" binding:"max=500"`
}

// SaleItemRequest represents a sale item in the create sale request
type SaleItemRequest struct {
	ProductID string  `json:"product_id" binding:"required,uuid"`
	Quantity  int     `json:"quantity" binding:"required,min=1"`
	UnitPrice float64 `json:"unit_price" binding:"required,min=0"`
	Discount  float64 `json:"discount" binding:"min=0"`
}

// SalePaymentRequest represents a payment in the create sale request
type SalePaymentRequest struct {
	Amount      float64 `json:"amount" binding:"required,min=0"`
	PaymentType string  `json:"payment_type" binding:"required,oneof=cash card mobile"`
	Reference   string  `json:"reference" binding:"max=100"`
}
