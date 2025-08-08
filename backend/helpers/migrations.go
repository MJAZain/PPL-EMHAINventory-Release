package helpers

import (
	"go-gin-auth/config"
	"go-gin-auth/internal/adjustment"
	"go-gin-auth/internal/brand"
	"go-gin-auth/internal/category"
	"go-gin-auth/internal/doctor"
	"go-gin-auth/internal/drug_category"
	"go-gin-auth/internal/expense"
	"go-gin-auth/internal/expense_type"
	"go-gin-auth/internal/incomingProducts"
	"go-gin-auth/internal/nonpbf"
	"go-gin-auth/internal/opname"
	"go-gin-auth/internal/outgoingProducts"
	"go-gin-auth/internal/patient"
	"go-gin-auth/internal/pbf"
	"go-gin-auth/internal/prescription"
	"go-gin-auth/internal/product"
	"go-gin-auth/internal/sales"
	"go-gin-auth/internal/shift"
	"go-gin-auth/internal/stock"
	"go-gin-auth/internal/stock_correction"
	storagelocation "go-gin-auth/internal/storage_location"
	"go-gin-auth/internal/supplier"
	"go-gin-auth/internal/unit"
	"go-gin-auth/model"
	"go-gin-auth/service"
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
		&prescription.PrescriptionSale{},
		&prescription.PrescriptionItem{},
		&sales.SalesRegular{},
		&sales.SalesRegularItem{},
		&adjustment.StockAdjustment{},
		&expense_type.ExpenseType{},
		&expense.Expense{},
	)

	var count int64
	err = db.Model(&model.User{}).Where("email = ?", "admin@admin.com").Count(&count).Error
	if err != nil {
		return err
	}

	if count == 0 {
		password, _ := service.HashPassword("admin")
		user := model.User{
			FullName: "Admin User",
			Email:    "admin@admin.com",
			Password: password,
			Role:     "admin",
			Active:   true,
		}

		err = db.Create(&user).Error
		if err != nil {
			return err
		}
	}

	var sysConfigCount int64
	err = db.Model(&model.SystemConfig{}).Count(&sysConfigCount).Error
	if err != nil {
		return err
	}

	if sysConfigCount == 0 {
		systemConfigs := []model.SystemConfig{
			{
				MaxFailedLogin:  5,  // Batas maksimal login gagal
				LockoutDuration: 30, // Durasi lockout dalam menit
			},
		}

		err = db.Create(&systemConfigs).Error
		if err != nil {
			return err
		}
	}

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
