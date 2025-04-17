package entities

type Company struct {
	Base
	Name    string `json:"name" binding:"required"`
	Address string `json:"address" binding:"required"`
	Phone   string `json:"phone" binding:"required"`
	Email   string `json:"email" binding:"required,email"`
	Slogan  string `json:"slogan"`
	Remarks string `json:"remarks"`
	Shops   []Shop `gorm:"foreignKey:CompanyID" json:"shops,omitempty"`
}
