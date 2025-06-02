package brand

import (
	"go-gin-auth/config"
	"go-gin-auth/model"
	"go-gin-auth/pkg/pagination"
	"strings"

	"gorm.io/gorm"
)

type BrandRepository interface {
	CreateBrand(brand Brand) (Brand, error)
	GetBrands(page, limit int, search string) ([]Brand, int64, error)
	GetBrandByID(ID uint) (Brand, error)
	UpdateBrand(ID uint, brand Brand) (Brand, error)
	DeleteBrand(ID uint, brand Brand) error
}

type repository struct {
	db *gorm.DB
}

func NewBrandRepository() BrandRepository {
	return &repository{db: config.DB}
}

func (r *repository) CreateBrand(brand Brand) (Brand, error) {
	err := r.db.Create(&brand).Error
	if err != nil {
		return brand, err
	}
	return r.loadAssociatedUserNames(brand)
}

func (r *repository) GetBrands(page, limit int, search string) ([]Brand, int64, error) {
	var brands []Brand
	var totalData int64

	query := r.db.Model(&Brand{})

	if search != "" {
		searchTerm := "%" + strings.ToLower(search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(description) LIKE ?", searchTerm, searchTerm)
	}

	err := query.Count(&totalData).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Scopes(pagination.PaginateScope(page, limit)).Order("created_at desc").Find(&brands).Error
	if err != nil {
		return nil, totalData, err
	}

	for i := range brands {
		brands[i], _ = r.loadAssociatedUserNames(brands[i])
	}

	return brands, totalData, nil
}

func (r *repository) GetBrandByID(ID uint) (Brand, error) {
	var brand Brand
	err := r.db.First(&brand, ID).Error
	if err != nil {
		return brand, err
	}
	return r.loadAssociatedUserNames(brand)
}

func (r *repository) UpdateBrand(ID uint, brandInput Brand) (Brand, error) {
	var existingBrand Brand
	err := r.db.First(&existingBrand, ID).Error
	if err != nil {
		return brandInput, err
	}

	existingBrand.Name = brandInput.Name
	existingBrand.Description = brandInput.Description
	existingBrand.UpdatedBy = brandInput.UpdatedBy

	err = r.db.Save(&existingBrand).Error
	if err != nil {
		return existingBrand, err
	}
	return r.loadAssociatedUserNames(existingBrand)
}

func (r *repository) DeleteBrand(ID uint, brandInput Brand) error {
	var existingBrand Brand
	err := r.db.First(&existingBrand, ID).Error
	if err != nil {
		return err
	}

	updateData := map[string]interface{}{
		"deleted_by": brandInput.DeletedBy,
		"updated_by": brandInput.DeletedBy,
	}

	if err := r.db.Model(&existingBrand).Updates(updateData).Error; err != nil {
		return err
	}

	err = r.db.Delete(&existingBrand).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) loadAssociatedUserNames(b Brand) (Brand, error) {
	if b.CreatedBy > 0 {
		var createdByUser model.User
		if err := r.db.Select("full_name").First(&createdByUser, b.CreatedBy).Error; err == nil {
			b.CreatedByName = createdByUser.FullName
		}
	}

	if b.UpdatedBy > 0 {
		var updatedByUser model.User
		if err := r.db.Select("full_name").First(&updatedByUser, b.UpdatedBy).Error; err == nil {
			b.UpdatedByName = updatedByUser.FullName
		}
	}

	if b.DeletedBy > 0 {
		var deletedByUser model.User
		if err := r.db.Select("full_name").First(&deletedByUser, b.DeletedBy).Error; err == nil {
			b.DeletedByName = deletedByUser.FullName
		}
	}
	return b, nil
}
