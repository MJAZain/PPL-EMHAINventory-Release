// Helper functions
package prescription

import (
	"fmt"
	"go-gin-auth/internal/stock"
	"math/rand"
	"time"
)

func generateTransactionCode() string {
	return fmt.Sprintf("TRX-%d-%s", time.Now().Unix(), generateRandomString(4))
}

func generateRandomString(length int) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

// ValidateUpdateRequest validates the update request
func (s *PrescriptionSaleService) ValidateUpdateRequest(req *CreatePrescriptionSaleRequest) error {
	if req.PrescriptionNo == "" {
		return fmt.Errorf("prescription number is required")
	}

	if len(req.Items) == 0 {
		return fmt.Errorf("at least one item is required")
	}

	for i, item := range req.Items {
		if item.ProductID == 0 {
			return fmt.Errorf("product ID is required for item %d", i+1)
		}
		if item.Quantity <= 0 {
			return fmt.Errorf("quantity must be positive for item %d", i+1)
		}
		if item.Price < 0 {
			return fmt.Errorf("price cannot be negative for item %d", i+1)
		}
	}

	return nil
}

// GetStockSummary returns stock summary for debugging
func (s *PrescriptionSaleService) GetStockSummary(productIDs []uint) (map[uint]int, error) {
	var stocks []stock.Stock
	err := s.db.Where("product_id IN ?", productIDs).Find(&stocks).Error
	if err != nil {
		return nil, err
	}

	summary := make(map[uint]int)
	for _, stock := range stocks {
		summary[stock.ProductID] = stock.Quantity
	}

	return summary, nil
}
