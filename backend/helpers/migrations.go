package helpers

import (
	"go-gin-auth/config"
	"go-gin-auth/internal/brand"
	"go-gin-auth/internal/category"
	"go-gin-auth/internal/doctor"
	"go-gin-auth/internal/drug_category"
	"go-gin-auth/internal/incomingProducts"
	"go-gin-auth/internal/nonpbf"
	"go-gin-auth/internal/opname"
	"go-gin-auth/internal/outgoingProducts"
	"go-gin-auth/internal/patient"
	"go-gin-auth/internal/pbf"
	"go-gin-auth/internal/product"
	"go-gin-auth/internal/shift"
	"go-gin-auth/internal/stock"
	"go-gin-auth/internal/stock_correction"
	storagelocation "go-gin-auth/internal/storage_location"
	"go-gin-auth/internal/supplier"
	"go-gin-auth/internal/unit"
	"go-gin-auth/model"
)

func MigrateDB() error {
	db := config.DB

	err := db.AutoMigrate(
		&model.User{},
		&model.ActivityLog{},
		&model.SystemConfig{},
		&product.Product{},
		&unit.Unit{},
		&category.Category{},
		&model.AuditLog{},
		&model.Transaksi{},
		&opname.StockOpname{},
		&opname.StockOpnameDetail{},
		&incomingProducts.IncomingProduct{},
		&incomingProducts.IncomingProductDetail{},
		&stock.Stock{},
		&outgoingProducts.OutgoingProduct{},
		&outgoingProducts.OutgoingProductDetail{},
		&brand.Brand{},
		&storagelocation.StorageLocation{},
		&supplier.Supplier{},
		&doctor.Doctor{},
		&patient.Patient{},
		&drug_category.DrugCategory{},
		&shift.Shift{},
		&stock_correction.StockCorrection{},
		&pbf.IncomingPBF{}, &pbf.IncomingPBFDetail{},
		&nonpbf.IncomingNonPBF{}, &nonpbf.IncomingNonPBFDetail{},
	)

	if err != nil {
		return err
	}

	if db.Dialector.Name() != "sqlite" {
		if db.Migrator().HasColumn(&product.Product{}, "storage_location") {
			if !db.Migrator().HasColumn(&product.Product{}, "StorageLocationID") {
				err = db.Migrator().AddColumn(&product.Product{}, "StorageLocationID")
				if err != nil {
					return err
				}
			}
		}
		if db.Migrator().HasColumn(&product.Product{}, "brand") {
			if !db.Migrator().HasColumn(&product.Product{}, "BrandID") {
				err = db.Migrator().AddColumn(&product.Product{}, "BrandID")
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
