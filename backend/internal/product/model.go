package product

import (
	"go-gin-auth/internal/brand"
	"go-gin-auth/internal/category"
	"go-gin-auth/internal/drug_category"
	storagelocation "go-gin-auth/internal/storage_location"
	"go-gin-auth/internal/unit"
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID                     uint                            `gorm:"primaryKey" json:"id"`
	Name                   string                          `gorm:"type:varchar(255);not null;comment:Nama Obat" json:"name" form:"name"`
	Code                   string                          `gorm:"type:varchar(50);not null;comment:Kode/SKU" json:"code" form:"code"`
	Barcode                string                          `gorm:"type:varchar(100);not null;comment:Barcode" json:"barcode" form:"barcode"`
	CategoryID             uint                            `gorm:"not null;comment:ID Kategori" json:"category_id" form:"category_id"`
	Category               category.Category               `gorm:"-" json:"category"`
	UnitID                 uint                            `gorm:"not null;comment:ID Satuan" json:"unit_id" form:"unit_id"`
	Unit                   unit.Unit                       `gorm:"-" json:"unit"`
	SellingPrice           float64                         `gorm:"type:decimal(15,2);not null;comment:Harga Jual" json:"selling_price" form:"selling_price"`
	StorageLocationID      uint                            `gorm:"not null;comment:ID Lokasi Penyimpanan;default:1" json:"storage_location_id" form:"storage_location_id"`
	StorageLocation        storagelocation.StorageLocation `gorm:"-" json:"storage_location"`
	BrandID                uint                            `gorm:"not null;comment:ID Brand;default:1" json:"brand_id" form:"brand_id"`
	Brand                  brand.Brand                     `gorm:"-" json:"brand"`
	CreatedBy              uint                            `gorm:"not null;comment:ID Pengguna Pembuat" json:"created_by"`
	CreatedAt              time.Time                       `json:"created_at"`
	UpdatedBy              uint                            `gorm:"not null;comment:ID Pengguna Pengubah" json:"updated_by"`
	UpdatedAt              time.Time                       `json:"updated_at"`
	DeletedBy              uint                            `gorm:"comment:ID Pengguna Penghapus" json:"deleted_by"`
	DeletedAt              gorm.DeletedAt                  `gorm:"index" json:"deleted_at"`
	DosageDescription      string                          `gorm:"type:text;comment:Deskripsi Dosis" json:"dosage_description" form:"dosage_description"`
	CompositionDescription string                          `gorm:"type:text;comment:Deskripsi Komposisi Obat" json:"composition_description" form:"composition_description"`
	DrugCategoryID         uint                            `gorm:"not null;comment:ID Kategori Obat;default:1" json:"drug_category_id" form:"drug_category_id"`
	DrugCategory           drug_category.DrugCategory      `gorm:"-" json:"drug_category"`
	MinStock               int                             `gorm:"not null;default:0;comment:Stok Minimum" json:"min_stock" form:"min_stock"`
}
