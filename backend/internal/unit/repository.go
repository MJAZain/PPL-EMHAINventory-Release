package unit

import (
	"go-gin-auth/config"
	"go-gin-auth/model"
	"time"

	"gorm.io/gorm"
)

type UnitRepository interface {
	CreateUnit(unit Unit) (Unit, error)
	GetUnits() ([]Unit, error)
	GetUnitByID(ID uint) (Unit, error)
	UpdateUnit(ID uint, unit Unit) (Unit, error)
	DeleteUnit(ID uint, unit Unit) error
}

type repository struct {
	db *gorm.DB
}

func NewUnitRepository() *repository {
	return &repository{db: config.DB}
}

func (r *repository) CreateUnit(unit Unit) (Unit, error) {
	err := r.db.Create(&unit).Error
	if err != nil {
		return unit, err
	}

	return unit, nil
}

func (r *repository) GetUnits() ([]Unit, error) {
	var units []Unit

	if err := r.db.Find(&units).Error; err != nil {
		return nil, err
	}

	for i := range units {
		var createdByUser model.User
		if units[i].CreatedBy > 0 {
			if err := r.db.Select("full_name").First(&createdByUser, units[i].CreatedBy).Error; err == nil {
				units[i].CreatedByName = createdByUser.FullName
			}
		}

		if units[i].UpdatedBy > 0 {
			var updatedByUser model.User
			if err := r.db.Select("full_name").First(&updatedByUser, units[i].UpdatedBy).Error; err == nil {
				units[i].UpdatedByName = updatedByUser.FullName
			}
		}
	}

	return units, nil
}

func (r *repository) GetUnitByID(ID uint) (Unit, error) {
	var unit Unit
	err := r.db.First(&unit, ID).Error
	if err != nil {
		return unit, err
	}

	var createdByUser model.User
	if unit.CreatedBy > 0 {
		if err := r.db.Select("full_name").First(&createdByUser, unit.CreatedBy).Error; err == nil {
			unit.CreatedByName = createdByUser.FullName
		}
	}

	if unit.UpdatedBy > 0 {
		var updatedByUser model.User
		if err := r.db.Select("full_name").First(&updatedByUser, unit.UpdatedBy).Error; err == nil {
			unit.UpdatedByName = updatedByUser.FullName
		}
	}

	return unit, nil
}

func (r *repository) UpdateUnit(ID uint, unit Unit) (Unit, error) {
	var existingUnit Unit
	err := r.db.First(&existingUnit, ID).Error
	if err != nil {
		return unit, err
	}

	existingUnit.Name = unit.Name
	existingUnit.Description = unit.Description

	err = r.db.Save(&existingUnit).Error
	if err != nil {
		return unit, err
	}

	return existingUnit, nil
}

func (r *repository) DeleteUnit(ID uint, unit Unit) error {
	var existingUnit Unit
	err := r.db.First(&existingUnit, ID).Error
	if err != nil {
		return err
	}

	existingUnit.Name = unit.Name
	existingUnit.DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}

	err = r.db.Save(&existingUnit).Error
	if err != nil {
		return err
	}

	return nil
}
