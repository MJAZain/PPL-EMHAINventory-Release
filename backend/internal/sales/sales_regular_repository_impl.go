package sales

import (
	"gorm.io/gorm"
)

type salesRegularRepository struct {
	db *gorm.DB
}

// Konstruktor
func NewSalesRegularRepository(db *gorm.DB) SalesRegularRepository {
	return &salesRegularRepository{db: db}
}

// Ambil semua transaksi penjualan reguler dengan paginasi
func (r *salesRegularRepository) GetAllSalesRegular(limit, offset int) ([]SalesRegular, int64, error) {
	var sales []SalesRegular
	var total int64

	query := r.db.Model(&SalesRegular{})

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("Items").
		Order("transaction_date DESC").
		Limit(limit).
		Offset(offset).
		Find(&sales).Error; err != nil {
		return nil, 0, err
	}

	return sales, total, nil
}

// Ambil detail satu transaksi
func (r *salesRegularRepository) GetSalesRegularByID(id uint) (*SalesRegular, error) {
	var sale SalesRegular
	if err := r.db.Preload("Items").First(&sale, id).Error; err != nil {
		return nil, err
	}
	return &sale, nil
}

// Tambah transaksi baru
func (r *salesRegularRepository) CreateSalesRegular(data *SalesRegular) error {
	return r.db.Create(data).Error
}

// Update transaksi
func (r *salesRegularRepository) UpdateSalesRegular(data *SalesRegular) error {
	return r.db.Save(data).Error
}

// Hapus transaksi (soft delete)
func (r *salesRegularRepository) DeleteSalesRegular(id uint) error {
	return r.db.Delete(&SalesRegular{}, id).Error
}
