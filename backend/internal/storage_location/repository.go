package storagelocation

import (
	"go-gin-auth/config"
	"go-gin-auth/model"
	"go-gin-auth/pkg/pagination"
	"strings"

	"gorm.io/gorm"
)

type StorageLocationRepository interface {
	CreateStorageLocation(storageLocation StorageLocation) (StorageLocation, error)
	GetStorageLocations(page, limit int, search string) ([]StorageLocation, int64, error)
	GetStorageLocationByID(ID uint) (StorageLocation, error)
	UpdateStorageLocation(ID uint, storageLocation StorageLocation) (StorageLocation, error)
	DeleteStorageLocation(ID uint, storageLocation StorageLocation) error
}

type repository struct {
	db *gorm.DB
}

func NewStorageLocationRepository() StorageLocationRepository {
	return &repository{db: config.DB}
}

func (r *repository) CreateStorageLocation(storageLocation StorageLocation) (StorageLocation, error) {
	err := r.db.Create(&storageLocation).Error
	if err != nil {
		return storageLocation, err
	}
	return r.loadAssociatedUserNames(storageLocation)
}

func (r *repository) GetStorageLocations(page, limit int, search string) ([]StorageLocation, int64, error) {
	var storageLocations []StorageLocation
	var totalData int64

	query := r.db.Model(&StorageLocation{})

	if search != "" {
		searchTerm := "%" + strings.ToLower(search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(description) LIKE ?", searchTerm, searchTerm)
	}

	err := query.Count(&totalData).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Scopes(pagination.PaginateScope(page, limit)).Order("created_at desc").Find(&storageLocations).Error
	if err != nil {
		return nil, totalData, err
	}

	for i := range storageLocations {
		storageLocations[i], _ = r.loadAssociatedUserNames(storageLocations[i])
	}

	return storageLocations, totalData, nil
}

func (r *repository) GetStorageLocationByID(ID uint) (StorageLocation, error) {
	var storageLocation StorageLocation
	err := r.db.First(&storageLocation, ID).Error
	if err != nil {
		return storageLocation, err
	}
	return r.loadAssociatedUserNames(storageLocation)
}

func (r *repository) UpdateStorageLocation(ID uint, slInput StorageLocation) (StorageLocation, error) {
	var existingSL StorageLocation
	err := r.db.First(&existingSL, ID).Error
	if err != nil {
		return slInput, err
	}

	existingSL.Name = slInput.Name
	existingSL.Description = slInput.Description
	existingSL.UpdatedBy = slInput.UpdatedBy

	err = r.db.Save(&existingSL).Error
	if err != nil {
		return existingSL, err
	}
	return r.loadAssociatedUserNames(existingSL)
}

func (r *repository) DeleteStorageLocation(ID uint, slInput StorageLocation) error {
	var existingSL StorageLocation
	err := r.db.First(&existingSL, ID).Error
	if err != nil {
		return err
	}

	updateData := map[string]interface{}{
		"deleted_by": slInput.DeletedBy,
		"updated_by": slInput.DeletedBy,
	}

	if err := r.db.Model(&existingSL).Updates(updateData).Error; err != nil {
		return err
	}

	err = r.db.Delete(&existingSL).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) loadAssociatedUserNames(sl StorageLocation) (StorageLocation, error) {
	if sl.CreatedBy > 0 {
		var createdByUser model.User
		if err := r.db.Select("full_name").First(&createdByUser, sl.CreatedBy).Error; err == nil {
			sl.CreatedByName = createdByUser.FullName
		}
	}

	if sl.UpdatedBy > 0 {
		var updatedByUser model.User
		if err := r.db.Select("full_name").First(&updatedByUser, sl.UpdatedBy).Error; err == nil {
			sl.UpdatedByName = updatedByUser.FullName
		}
	}

	if sl.DeletedBy > 0 {
		var deletedByUser model.User
		if err := r.db.Select("full_name").First(&deletedByUser, sl.DeletedBy).Error; err == nil {
			sl.DeletedByName = deletedByUser.FullName
		}
	}
	return sl, nil
}
