package sales

// Interface untuk Sales Regular Repository
type SalesRegularRepository interface {
	GetAllSalesRegular(limit, offset int) ([]SalesRegular, int64, error)
	GetSalesRegularByID(id uint) (*SalesRegular, error)
	CreateSalesRegular(data *SalesRegular) error
	UpdateSalesRegular(data *SalesRegular) error
	DeleteSalesRegular(id uint) error
}
