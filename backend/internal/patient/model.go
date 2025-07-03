package patient

type Patient struct {
	ID             uint   `gorm:"primaryKey" json:"id"`
	FullName       string `gorm:"type:varchar(512);not null" json:"full_name" form:"full_name"`
	Gender         string `gorm:"type:char(1);not null" json:"gender" form:"gender"`
	PlaceOfBirth   string `gorm:"type:varchar(100);not null" json:"place_of_birth" form:"place_of_birth"`
	DateOfBirth    string `gorm:"type:varchar(512);not null" json:"-" form:"date_of_birth"`
	Address        string `gorm:"type:text" json:"address" form:"address"`
	PhoneNumber    string `gorm:"type:varchar(512)" json:"phone_number" form:"phone_number"`
	PatientType    string `gorm:"type:varchar(50);not null" json:"patient_type" form:"patient_type"`
	IdentityNumber string `gorm:"type:varchar(512);index" json:"identity_number" form:"identity_number"`
	GuarantorName  string `gorm:"type:varchar(512)" json:"guarantor_name,omitempty" form:"guarantor_name"`
	Status         string `gorm:"type:varchar(20);not null;default:'Aktif'" json:"status" form:"status"`
	Age            int    `gorm:"-" json:"age,omitempty"`

	DecryptedDateOfBirth string `gorm:"-" json:"date_of_birth"`
}
