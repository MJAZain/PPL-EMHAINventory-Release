package category

import (
	"go-gin-auth/config"
	"go-gin-auth/model"
	"time"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	CreateCategory(category Category) (Category, error)
	GetCategories() ([]Category, error)
	GetCategoryByID(ID uint) (Category, error)
	UpdateCategory(ID uint, category Category) (Category, error)
	DeleteCategory(ID uint, category Category) error
}

type repository struct {
	db *gorm.DB
}

func NewCategoryRepository() *repository {
	return &repository{db: config.DB}
}

func (r *repository) CreateCategory(category Category) (Category, error) {
	err := r.db.Create(&category).Error
	if err != nil {
		return category, err
	}

	return category, nil
}

func (r *repository) GetCategories() ([]Category, error) {
	var categories []Category

	if err := r.db.Find(&categories).Error; err != nil {
		return nil, err
	}

	for i := range categories {
		var createdByUser model.User
		if categories[i].CreatedBy > 0 {
			if err := r.db.Select("full_name").First(&createdByUser, categories[i].CreatedBy).Error; err == nil {
				categories[i].CreatedByName = createdByUser.FullName
			}
		}

		if categories[i].UpdatedBy > 0 {
			var updatedByUser model.User
			if err := r.db.Select("full_name").First(&updatedByUser, categories[i].UpdatedBy).Error; err == nil {
				categories[i].UpdatedByName = updatedByUser.FullName
			}
		}
	}

	return categories, nil
}

func (r *repository) GetCategoryByID(ID uint) (Category, error) {
	var category Category
	err := r.db.First(&category, ID).Error
	if err != nil {
		return category, err
	}

	var createdByUser model.User
	if category.CreatedBy > 0 {
		if err := r.db.Select("full_name").First(&createdByUser, category.CreatedBy).Error; err == nil {
			category.CreatedByName = createdByUser.FullName
		}
	}

	if category.UpdatedBy > 0 {
		var updatedByUser model.User
		if err := r.db.Select("full_name").First(&updatedByUser, category.UpdatedBy).Error; err == nil {
			category.UpdatedByName = updatedByUser.FullName
		}
	}

	return category, nil
}

func (r *repository) UpdateCategory(ID uint, category Category) (Category, error) {
	var existingCategory Category
	err := r.db.First(&existingCategory, ID).Error
	if err != nil {
		return category, err
	}

	existingCategory.Name = category.Name
	existingCategory.Description = category.Description

	err = r.db.Save(&existingCategory).Error
	if err != nil {
		return category, err
	}

	return existingCategory, nil
}

func (r *repository) DeleteCategory(ID uint, category Category) error {
	var existingCategory Category
	err := r.db.First(&existingCategory, ID).Error
	if err != nil {
		return err
	}

	existingCategory.Name = category.Name
	existingCategory.DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}

	err = r.db.Save(&existingCategory).Error
	if err != nil {
		return err
	}

	return nil
}
