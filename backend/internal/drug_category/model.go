package drug_category

type DrugCategory struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"type:varchar(255);not null" json:"name" form:"name"`
	Description string `gorm:"type:text" json:"description,omitempty" form:"description"`
	Status      string `gorm:"type:varchar(20);not null;default:'Aktif'" json:"status" form:"status"`
}
