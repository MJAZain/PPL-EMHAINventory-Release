package supplier

type Supplier struct {
	ID            uint   `gorm:"primaryKey" json:"id"`
	Name          string `gorm:"type:varchar(255);not null" json:"name" form:"name"`
	Type          string `gorm:"type:varchar(100);not null" json:"type" form:"type"`
	Address       string `gorm:"type:text;not null" json:"address" form:"address"`
	Phone         string `gorm:"type:varchar(20);not null" json:"phone" form:"phone"`
	Email         string `gorm:"type:varchar(255)" json:"email,omitempty" form:"email"`
	ContactPerson string `gorm:"type:varchar(255);not null" json:"contact_person" form:"contact_person"`
	ContactNumber string `gorm:"type:varchar(20);not null" json:"contact_number" form:"contact_number"`
	Status        string `gorm:"type:varchar(20);not null;default:'Aktif'" json:"status" form:"status"`

	ProvinceID string `gorm:"type:varchar(10);not null" json:"province_id" form:"province_id"`
	CityID     string `gorm:"type:varchar(10);not null" json:"city_id" form:"city_id"`

	Province string `gorm:"-" json:"province"`
	City     string `gorm:"-" json:"city"`
}
