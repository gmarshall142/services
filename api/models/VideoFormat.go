package models

import (
	"errors"
	"gorm.io/gorm"
)

func (VideoFormat) TableName() string {
	return "videoformats"
}

type VideoFormat struct {
	ID          uint   `gorm:"primary_key;auto_increment" json:"id"`
	Name        string `gorm:"size:60;column:name;not null" json:"name"`
	Description string `gorm:"size:128;column:description" json:"description"`
}

func (u *VideoFormat) FindAllVideoFormats(db *gorm.DB) (*[]VideoFormat, error) {
	var err error
	var videoFormats []VideoFormat
	err = db.Debug().Model(&VideoFormat{}).Limit(100).Find(&videoFormats).Error
	if err != nil {
		return &[]VideoFormat{}, err
	}
	return &videoFormats, err
}

func (u *VideoFormat) FindVideoFormatByID(db *gorm.DB, uid uint32) (*VideoFormat, error) {
	var err error
	err = db.Debug().Model(VideoFormat{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &VideoFormat{}, err
	}
	if err == gorm.ErrRecordNotFound {
		return &VideoFormat{}, errors.New("Video Format Not Found")
	}
	return u, err
}
