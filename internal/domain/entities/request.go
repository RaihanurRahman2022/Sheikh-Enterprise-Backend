package entities

import "time"

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

// UpdatePasswordRequest represents the request body for updating a user's password
type UpdatePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,password"`
}

// CreateSupplierRequest represents the request body for creating a new supplier
type CreateSupplierRequest struct {
	Name    string `json:"name" binding:"required,min=3,max=100"`
	Address string `json:"address" binding:"required"`
	Phone   string `json:"phone" binding:"required,phone"`
	Email   string `json:"email" binding:"required,email"`
	Remarks string `json:"remarks" binding:"max=500"`
}

// CreateCompanyRequest represents the request body for creating a new company
type CreateCompanyRequest struct {
	Name    string `json:"name" binding:"required"`
	Address string `json:"address" binding:"required"`
	Phone   string `json:"phone" binding:"required"`
	Email   string `json:"email" binding:"required,email"`
	Slogan  string `json:"slogan"`
	Remarks string `json:"remarks"`
}

// CreateShopRequest represents the request body for creating a new shop
type CreateShopRequest struct {
	CompanyID    string `json:"company_id" binding:"required,uuid"`
	Name         string `json:"name" binding:"required"`
	Address      string `json:"address" binding:"required"`
	Phone        string `json:"phone" binding:"required"`
	Email        string `json:"email" binding:"required,email"`
	ManagerName  string `json:"manager_name" binding:"required"`
	ManagerPhone string `json:"manager_phone" binding:"required"`
	Remarks      string `json:"remarks"`
}

// Stock Transfer Requests
type CreateStockTransferRequest struct {
	FromShopID       string    `json:"from_shop_id" binding:"required,uuid"`
	ToShopID         string    `json:"to_shop_id" binding:"required,uuid"`
	ProductID        string    `json:"product_id" binding:"required,uuid"`
	Quantity         int       `json:"quantity" binding:"required,min=1"`
	TransferDateTime time.Time `json:"transfer_datetime" binding:"required"`
	Remarks          string    `json:"remarks"`
}

type UpdateStockTransferRequest struct {
	Status  string `json:"status" binding:"required,oneof=pending completed cancelled"`
	Remarks string `json:"remarks"`
}

type StockTransferFilter struct {
	FromShopID string    `form:"from_shop_id"`
	ToShopID   string    `form:"to_shop_id"`
	ProductID  string    `form:"product_id"`
	Status     string    `form:"status"`
	StartDate  time.Time `form:"start_date"`
	EndDate    time.Time `form:"end_date"`
	Page       int       `form:"page,default=1" binding:"min=1"`
	PageSize   int       `form:"page_size,default=10" binding:"min=1,max=100"`
}
