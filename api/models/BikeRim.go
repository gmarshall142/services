package models

import (
	"errors"
	"gorm.io/gorm"
)

type Tabler interface {
	TableName() string
}

func (BikeRim) TableName() string {
	return "bikerims"
}

type BikeRim struct {
	ID          uint   `gorm:"primary_key;auto_increment" json:"id"`
	Name        string `gorm:"size:60;column:name;not null" json:"name"`
	Description string `gorm:"size:128;column:description" json:"description"`
	Diameter    uint   `gorm:"type:integer;column:diameter" json:"diameter"`
}

func (u *BikeRim) FindAllBikeRims(db *gorm.DB) (*[]BikeRim, error) {
	var err error
	var bikeRims []BikeRim
	err = db.Debug().Model(&BikeRim{}).Limit(100).Find(&bikeRims).Error
	if err != nil {
		return &[]BikeRim{}, err
	}
	return &bikeRims, err
}

func (u *BikeRim) FindBikeRimByID(db *gorm.DB, uid uint32) (*BikeRim, error) {
	var err error
	err = db.Debug().Model(BikeRim{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &BikeRim{}, err
	}
	if err == gorm.ErrRecordNotFound {
		return &BikeRim{}, errors.New("Bike Rim Not Found")
	}
	return u, err
}
