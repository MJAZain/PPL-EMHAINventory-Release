package utils

import (
	"fmt"
	"strconv"
)

func ConvertProductID(productIDStr string) (uint, error) {
	idUint64, err := strconv.ParseUint(productIDStr, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid product ID: %w", err)
	}
	return uint(idUint64), nil
}
