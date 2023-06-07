package models

import (
	"gorm.io/gorm"
)

func (AudioFormat) TableName() string {
	return "audioformats"
}

type AudioFormat struct {
	ID          uint   `gorm:"primary_key;auto_increment" json:"id"`
	Name        string `gorm:"size:60;column:name;not null" json:"name"`
	Description string `gorm:"size:128;column:description" json:"description"`
}

func (u *AudioFormat) FindAllAudioFormats(db *gorm.DB) (*[]AudioFormat, error) {
	var err error
	var audioFormats []AudioFormat
	err = db.Debug().Model(&AudioFormat{}).Limit(100).Find(&audioFormats).Error
	if err != nil {
		return &[]AudioFormat{}, err
	}
	return &audioFormats, err
}

//func (u *VideoFormat) FindVideoFormatByID(db *gorm.DB, uid uint32) (*VideoFormat, error) {
//	var err error
//	err = db.Debug().Model(VideoFormat{}).Where("id = ?", uid).Take(&u).Error
//	if err != nil {
//		return &VideoFormat{}, err
//	}
//	if err == gorm.ErrRecordNotFound {
//		return &VideoFormat{}, errors.New("Video Format Not Found")
//	}
//	return u, err
//}
