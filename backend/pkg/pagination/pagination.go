package pagination

import (
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	DefaultPage  = 1
	DefaultLimit = 10
)

type PaginationResult struct {
	TotalData    int64       `json:"total_data"`
	TotalPages   int         `json:"total_pages"`
	CurrentPage  int         `json:"current_page"`
	NextPage     *int        `json:"next_page"`
	PrevPage     *int        `json:"prev_page"`
	LimitPerPage int         `json:"limit_per_page"`
	Data         interface{} `json:"data"`
}

func PaginateScope(page, limit int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := (page - 1) * limit
		return db.Offset(offset).Limit(limit)
	}
}

func GetPaginationParams(c *gin.Context) (page, limit, offset int) {
	pageStr := c.DefaultQuery("page", strconv.Itoa(DefaultPage))
	limitStr := c.DefaultQuery("limit", strconv.Itoa(DefaultLimit))

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = DefaultPage
	}

	limit, err = strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = DefaultLimit
	}

	offset = (page - 1) * limit
	return page, limit, offset
}

func CreatePaginationResult(data interface{}, totalData int64, page, limit int) PaginationResult {
	totalPages := int(math.Ceil(float64(totalData) / float64(limit)))

	var nextPage *int
	if page < totalPages {
		val := page + 1
		nextPage = &val
	}

	var prevPage *int
	if page > 1 && page <= totalPages {
		val := page - 1
		prevPage = &val
	}

	if totalData == 0 {
		page = 1
		totalPages = 0
		nextPage = nil
		prevPage = nil
	} else if page > totalPages && totalPages > 0 {
		page = totalPages
	}

	return PaginationResult{
		TotalData:    totalData,
		TotalPages:   totalPages,
		CurrentPage:  page,
		NextPage:     nextPage,
		PrevPage:     prevPage,
		LimitPerPage: limit,
		Data:         data,
	}
}
