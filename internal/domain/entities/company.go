package entities

type Company struct {
	Base
	Name    string `gorm:"type:varchar(255);not null" json:"name"`
	Address string `gorm:"type:text" json:"address"`
	Phone   string `gorm:"type:varchar(20)" json:"phone"`
	Email   string `gorm:"type:varchar(255);unique" json:"email"`
	Slogan  string `gorm:"type:text" json:"slogan"`
	Remarks string `gorm:"type:text" json:"remarks"`
	Shops   []Shop `gorm:"foreignKey:CompanyID" json:"shops,omitempty"`
}
