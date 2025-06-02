package outgoingProducts

import (
	"errors"
	"go-gin-auth/internal/stock"
)

type Service interface {
	CreateOutgoingProduct(outgoingProduct *OutgoingProduct, details []OutgoingProductDetail) error
	GetAllOutgoingProducts() ([]OutgoingProduct, error)
	GetOutgoingProductByID(id uint) (*OutgoingProduct, error)
	GetOutgoingProductDetails(outgoingProductID uint) ([]OutgoingProductDetail, error)
	UpdateOutgoingProduct(id uint, outgoingProduct *OutgoingProduct) error
	UpdateOutgoingProductDetails(details []OutgoingProductDetail) error
	DeleteOutgoingProduct(id uint) error
}

type service struct {
	repository      Repository
	repositoryStock stock.Repository
}

func NewService() *service {
	return &service{repository: NewRepository(), repositoryStock: stock.NewRepository()}
}

func (s *service) CreateOutgoingProduct(outgoingProduct *OutgoingProduct, details []OutgoingProductDetail) error {
	// Validasi
	if outgoingProduct.Date == "" {
		return errors.New("tanggal tidak boleh kosong")
	}
	if outgoingProduct.Customer == "" {
		return errors.New("customer tidak boleh kosong")
	}
	if outgoingProduct.NoFaktur == "" {
		return errors.New("nomor faktur tidak boleh kosong")
	}
	if outgoingProduct.PaymentStatus == "" {
		return errors.New("status pembayaran tidak boleh kosong")
	}
	if len(details) == 0 {
		return errors.New("detail produk keluar tidak boleh kosong")
	}

	// Hitung total setiap detail
	for i := range details {
		if details[i].ProductID == 0 {
			return errors.New("id produk tidak boleh kosong")
		}
		if details[i].Quantity <= 0 {
			return errors.New("kuantitas harus lebih dari 0")
		}
		if details[i].Price <= 0 {
			return errors.New("harga harus lebih dari 0")
		}

		// Hitung total
		details[i].Total = float64(details[i].Quantity) * details[i].Price

		// Kurangi stok produk (false untuk mengurangi)
		err := s.repositoryStock.UpdateProductStock(details[i].ProductID, details[i].Quantity, false)
		if err != nil {
			return errors.New("gagal memperbarui stok produk")
		}
	}

	return s.repository.Create(outgoingProduct, details)
}

func (s *service) GetAllOutgoingProducts() ([]OutgoingProduct, error) {
	return s.repository.GetAll()
}

func (s *service) GetOutgoingProductByID(id uint) (*OutgoingProduct, error) {
	return s.repository.GetByID(id)
}

func (s *service) GetOutgoingProductDetails(outgoingProductID uint) ([]OutgoingProductDetail, error) {
	return s.repository.GetDetailsByOutgoingProductID(outgoingProductID)
}

func (s *service) UpdateOutgoingProduct(id uint, outgoingProduct *OutgoingProduct) error {
	// Validasi
	if outgoingProduct.Date == "" {
		return errors.New("tanggal tidak boleh kosong")
	}
	if outgoingProduct.Customer == "" {
		return errors.New("customer tidak boleh kosong")
	}
	if outgoingProduct.NoFaktur == "" {
		return errors.New("nomor faktur tidak boleh kosong")
	}
	if outgoingProduct.PaymentStatus == "" {
		return errors.New("status pembayaran tidak boleh kosong")
	}

	return s.repository.Update(id, outgoingProduct)
}

func (s *service) UpdateOutgoingProductDetails(details []OutgoingProductDetail) error {
	// Validasi
	if len(details) == 0 {
		return errors.New("detail produk keluar tidak boleh kosong")
	}

	// Hitung total setiap detail
	for i := range details {
		if details[i].ProductID == 0 {
			return errors.New("id produk tidak boleh kosong")
		}
		if details[i].Quantity <= 0 {
			return errors.New("kuantitas harus lebih dari 0")
		}
		if details[i].Price <= 0 {
			return errors.New("harga harus lebih dari 0")
		}

		// Hitung total
		details[i].Total = float64(details[i].Quantity) * details[i].Price

		existingDetail, err := s.repository.GetDetailByOutgoingProductID(details[i].OutgoingProductID)
		if err != nil {
			return errors.New("gagal mendapatkan detail produk keluar")
		}

		// Kembalikan stok produk lama sebelum update (true untuk menambah)
		err = s.repositoryStock.UpdateProductStock(details[i].ProductID, existingDetail.Quantity, true)
		if err != nil {
			return errors.New("gagal memperbarui stok produk")
		}

		// Kurangi stok dengan jumlah baru (false untuk mengurangi)
		err = s.repositoryStock.UpdateProductStock(details[i].ProductID, details[i].Quantity, false)
		if err != nil {
			return errors.New("gagal memperbarui stok produk")
		}
	}

	return s.repository.UpdateDetails(details)
}

func (s *service) DeleteOutgoingProduct(id uint) error {
	details, err := s.repository.GetDetailsByOutgoingProductID(id)
	if err != nil {
		return err
	}

	for _, detail := range details {
		// Kembalikan stok saat menghapus produk keluar (true untuk menambah)
		if err := s.repositoryStock.UpdateProductStock(detail.ProductID, detail.Quantity, true); err != nil {
			return err
		}
	}

	return s.repository.Delete(id)
}
