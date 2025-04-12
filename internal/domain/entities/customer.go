package entities

type Customer struct {
	Base
	Name    string `gorm:"type:varchar(255);not null" json:"name"`
	Address string `gorm:"type:text" json:"address"`
	Phone   string `gorm:"type:varchar(20)" json:"phone"`
	Email   string `gorm:"type:varchar(255)" json:"email"`
	Remarks string `gorm:"type:text" json:"remarks"`
}
