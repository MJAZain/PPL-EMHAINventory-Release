package storagelocation

import (
	"errors"
	"strings"
)

type StorageLocationService struct {
	repository StorageLocationRepository
}

func NewStorageLocationService(repo StorageLocationRepository) *StorageLocationService {
	return &StorageLocationService{repository: repo}
}

func (s *StorageLocationService) CreateStorageLocation(sl StorageLocation) (StorageLocation, error) {
	if strings.TrimSpace(sl.Name) == "" {
		return StorageLocation{}, errors.New("nama wajib diisi")
	}

	return s.repository.CreateStorageLocation(sl)
}

func (s *StorageLocationService) GetStorageLocations(page, limit int, search string) ([]StorageLocation, int64, error) {
	return s.repository.GetStorageLocations(page, limit, search)
}

func (s *StorageLocationService) GetStorageLocationByID(ID uint) (StorageLocation, error) {
	return s.repository.GetStorageLocationByID(ID)
}

func (s *StorageLocationService) UpdateStorageLocation(ID uint, sl StorageLocation) (StorageLocation, error) {
	if strings.TrimSpace(sl.Name) == "" {
		return StorageLocation{}, errors.New("nama wajib diisi untuk pembaruan")
	}

	_, err := s.repository.GetStorageLocationByID(ID)
	if err != nil {
		return StorageLocation{}, errors.New("lokasi penyimpanan tidak ditemukan")
	}

	return s.repository.UpdateStorageLocation(ID, sl)
}

func (s *StorageLocationService) DeleteStorageLocation(ID uint, sl StorageLocation) error {
	return s.repository.DeleteStorageLocation(ID, sl)
}
