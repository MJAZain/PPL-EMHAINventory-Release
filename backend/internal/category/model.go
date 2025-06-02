package category

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Name          string         `gorm:"type:varchar(100);not null;uniqueIndex;comment:Nama Kategori" json:"name" form:"name"`
	Description   string         `gorm:"type:varchar(255);comment:Deskripsi" json:"description" form:"description"`
	CreatedBy     uint           `gorm:"not null;comment:ID Pengguna Pembuat" json:"created_by"`
	CreatedByName string         `gorm:"-" json:"created_by_name,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedBy     uint           `gorm:"comment:ID Pengguna Pengubah" json:"updated_by"`
	UpdatedByName string         `gorm:"-" json:"updated_by_name,omitempty"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedBy     uint           `gorm:"comment:ID Pengguna Penghapus" json:"deleted_by"`
	DeletedByName string         `gorm:"-" json:"deleted_by_name,omitempty"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
