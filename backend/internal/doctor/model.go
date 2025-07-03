package doctor

type Doctor struct {
	ID              uint   `gorm:"primaryKey" json:"id"`
	FullName        string `gorm:"type:varchar(512);not null" json:"full_name" form:"full_name"`
	STRNumber       string `gorm:"type:varchar(512);index" json:"str_number,omitempty" form:"str_number"`
	Specialization  string `gorm:"type:varchar(255);not null" json:"specialization" form:"specialization"`
	PhoneNumber     string `gorm:"type:varchar(512);not null" json:"phone_number" form:"phone_number"`
	PracticeAddress string `gorm:"type:text" json:"practice_address,omitempty" form:"practice_address"`
	Email           string `gorm:"type:varchar(512);uniqueIndex" json:"email,omitempty" form:"email"`
	Status          string `gorm:"type:varchar(20);not null;default:'Aktif'" json:"status" form:"status"`
}
