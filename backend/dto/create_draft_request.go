package dto

import (
	"encoding/json"
	"fmt"
	"time"
)

type CreateDraftRequest struct {
	Tanggal DateOnly `json:"opname_date" binding:"required"`
	Catatan string   `json:"notes"`
}

type UpdateDraftRequest struct {
	Tanggal DateOnly `json:"opname_date" binding:"required"`
	Catatan string   `json:"notes"`
}

type DateOnly struct {
	time.Time
}

const dateFormat = "2006-01-02"

func (d *DateOnly) UnmarshalJSON(b []byte) error {
	s := string(b)
	s = s[1 : len(s)-1] // remove quotes
	t, err := time.Parse(dateFormat, s)
	if err != nil {
		return fmt.Errorf("format tanggal harus YYYY-MM-DD")
	}
	d.Time = t
	return nil
}

func (d DateOnly) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Format(dateFormat))
}
