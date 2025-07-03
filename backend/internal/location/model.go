package location

type Province struct {
	ID   int    `json:"id"`
	Name string `json:"province"`
}

type Regency struct {
	ID         int    `json:"id"`
	ProvinceID int    `json:"province_id"`
	Name       string `json:"regency"`
}
