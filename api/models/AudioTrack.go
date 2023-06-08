package models

import (
	"errors"
	"gorm.io/gorm"
)

func (AudioTrack) TableName() string {
	return "audiotracks"
}

type AudioTrack struct {
	ID       uint   `gorm:"primary_key;auto_increment" json:"id"`
	Title    string `gorm:"size:60;column:name;not null" json:"title"`
	Duration uint   `gorm:"column:duration" json:"duration"`
	Position string `gorm:"size:10;column:position" json:"position"`
	AudioId  uint   `gorm:"type:integer;column:audioid" json:"audioid"`
}

//func (u *AudioFormat) FindAllAudioFormats(db *gorm.DB) (*[]AudioFormat, error) {
//	var err error
//	var audioFormats []AudioFormat
//	err = db.Debug().Model(&AudioFormat{}).Limit(100).Find(&audioFormats).Error
//	if err != nil {
//		return &[]AudioFormat{}, err
//	}
//	return &audioFormats, err
//}

func (u *AudioTrack) FindAudioTracksByAudioID(db *gorm.DB, uid uint32) (*[]AudioTrack, error) {
	var err error
	var audioTracks []AudioTrack
	err = db.Debug().Model(&AudioTrack{}).Where("audioid = ?", uid).Find(&audioTracks).Error
	if err != nil {
		return nil, err
	}
	if err == gorm.ErrRecordNotFound {
		return nil, errors.New("audio tracks Not Found")
	}
	return &audioTracks, err
}
