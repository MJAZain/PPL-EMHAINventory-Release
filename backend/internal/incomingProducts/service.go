package incomingProducts

import (
	"errors"
	"go-gin-auth/internal/stock"
)

type Service interface {
	CreateIncomingProduct(incomingProduct *IncomingProduct, details []IncomingProductDetail) error
	GetAllIncomingProducts() ([]IncomingProduct, error)
	GetIncomingProductByID(id uint) (*IncomingProduct, error)
	GetIncomingProductDetails(incomingProductID uint) ([]IncomingProductDetail, error)
	UpdateIncomingProduct(id uint, incomingProduct *IncomingProduct) error
	UpdateIncomingProductDetails(details []IncomingProductDetail) error
	DeleteIncomingProduct(id uint) error
}

type service struct {
	repository      Repository
	repositoryStock stock.Repository
}

func NewService() *service {
	return &service{repository: NewRepository(), repositoryStock: stock.NewRepository()}
}

func (s *service) CreateIncomingProduct(incomingProduct *IncomingProduct, details []IncomingProductDetail) error {
	// Validasi
	if incomingProduct.Date == "" {
		return errors.New("tanggal tidak boleh kosong")
	}
	if incomingProduct.Supplier == "" {
		return errors.New("supplier tidak boleh kosong")
	}
	if incomingProduct.NoFaktur == "" {
		return errors.New("nomor faktur tidak boleh kosong")
	}
	if incomingProduct.PaymentStatus == "" {
		return errors.New("status pembayaran tidak boleh kosong")
	}
	if len(details) == 0 {
		return errors.New("detail produk masuk tidak boleh kosong")
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

		err := s.repositoryStock.UpdateProductStock(details[i].ProductID, details[i].Quantity, true)
		if err != nil {
			return errors.New("gagal memperbarui stok produk")
		}
	}

	return s.repository.Create(incomingProduct, details)
}

func (s *service) GetAllIncomingProducts() ([]IncomingProduct, error) {
	return s.repository.GetAll()
}

func (s *service) GetIncomingProductByID(id uint) (*IncomingProduct, error) {
	return s.repository.GetByID(id)
}

func (s *service) GetIncomingProductDetails(incomingProductID uint) ([]IncomingProductDetail, error) {
	return s.repository.GetDetailsByIncomingProductID(incomingProductID)
}

func (s *service) UpdateIncomingProduct(id uint, incomingProduct *IncomingProduct) error {
	// Validasi
	if incomingProduct.Date == "" {
		return errors.New("tanggal tidak boleh kosong")
	}
	if incomingProduct.Supplier == "" {
		return errors.New("supplier tidak boleh kosong")
	}
	if incomingProduct.NoFaktur == "" {
		return errors.New("nomor faktur tidak boleh kosong")
	}
	if incomingProduct.PaymentStatus == "" {
		return errors.New("status pembayaran tidak boleh kosong")
	}

	return s.repository.Update(id, incomingProduct)
}

func (s *service) UpdateIncomingProductDetails(details []IncomingProductDetail) error {
	// Validasi
	if len(details) == 0 {
		return errors.New("detail produk masuk tidak boleh kosong")
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

		existingDetail, err := s.repository.GetDetailByIncomingProductID(details[i].IncomingProductID)
		if err != nil {
			return errors.New("gagal mendapatkan detail produk masuk")
		}

		err = s.repositoryStock.UpdateProductStock(details[i].ProductID, existingDetail.Quantity, false)
		if err != nil {
			return errors.New("gagal memperbarui stok produk")
		}

		err = s.repositoryStock.UpdateProductStock(details[i].ProductID, details[i].Quantity, true)
		if err != nil {
			return errors.New("gagal memperbarui stok produk")
		}
	}

	return s.repository.UpdateDetails(details)
}

func (s *service) DeleteIncomingProduct(id uint) error {
	details, err := s.repository.GetDetailsByIncomingProductID(id)
	if err != nil {
		return err
	}

	for _, detail := range details {
		if err := s.repositoryStock.UpdateProductStock(detail.ProductID, detail.Quantity, false); err != nil {
			return err
		}
	}

	return s.repository.Delete(id)
}
