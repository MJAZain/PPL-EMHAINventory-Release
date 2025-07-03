package location

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"sync"
)

var (
	provincesByID     map[string]Province
	regenciesByID     map[string]Regency
	provinceRegencies map[string]map[string]bool
	allProvinces      []Province
	once              sync.Once
)

type Service interface {
	GetAllProvinces() ([]Province, error)
	GetRegenciesByProvinceID(provinceID string) ([]Regency, error)
	ValidateLocation(provinceID, cityID string) (bool, error)
	GetLocationNames(provinceID, cityID string) (provinceName, cityName string)
}

type service struct{}

func NewService() Service {
	once.Do(loadLocationData)
	return &service{}
}

func loadLocationData() {
	var provincesData []Province
	var regenciesData []Regency

	provinceBytes, err := os.ReadFile("data/provinces.json")
	if err != nil {
		log.Fatalf("Gagal membaca file provinces.json: %v", err)
	}
	json.Unmarshal(provinceBytes, &provincesData)

	regencyBytes, err := os.ReadFile("data/regencies.json")
	if err != nil {
		log.Fatalf("Gagal membaca file regencies.json: %v", err)
	}
	json.Unmarshal(regencyBytes, &regenciesData)

	provincesByID = make(map[string]Province)
	regenciesByID = make(map[string]Regency)
	provinceRegencies = make(map[string]map[string]bool)
	allProvinces = provincesData

	for _, p := range provincesData {
		provincesByID[strconv.Itoa(p.ID)] = p
		provinceRegencies[strconv.Itoa(p.ID)] = make(map[string]bool)
	}

	for _, r := range regenciesData {
		regenciesByID[strconv.Itoa(r.ID)] = r
		if _, ok := provinceRegencies[strconv.Itoa(r.ProvinceID)]; ok {
			provinceRegencies[strconv.Itoa(r.ProvinceID)][strconv.Itoa(r.ID)] = true
		}
	}
	log.Println("Data lokasi berhasil dimuat dan diindeks.")
}

func (s *service) ValidateLocation(provinceID, cityID string) (bool, error) {
	if cities, ok := provinceRegencies[provinceID]; ok {
		if _, cityExists := cities[cityID]; cityExists {
			return true, nil
		}
	}
	return false, nil
}

func (s *service) GetLocationNames(provinceID, cityID string) (string, string) {
	provinceName := ""
	cityName := ""
	if province, ok := provincesByID[provinceID]; ok {
		provinceName = province.Name
	}
	if city, ok := regenciesByID[cityID]; ok {
		cityName = city.Name
	}
	return provinceName, cityName
}

func (s *service) GetAllProvinces() ([]Province, error) {
	return allProvinces, nil
}

func (s *service) GetRegenciesByProvinceID(provinceID string) ([]Regency, error) {
	var filteredRegencies []Regency
	provinceIDs, _ := strconv.Atoi(provinceID)
	for _, r := range regenciesByID {
		if r.ProvinceID == provinceIDs {
			filteredRegencies = append(filteredRegencies, r)
		}
	}
	return filteredRegencies, nil
}
