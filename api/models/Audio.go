package models

import (
	"github.com/lib/pq"
	"net/url"
	"strings"
	"time"
)

type Audio struct {
	ID            uint           `gorm:"primary_key;auto_increment" json:"id"`
	MasterID      uint           `gorm:"column:masterid;" json:"masterid"`
	Title         string         `gorm:"size:128;column:title;not null;" json:"title"`
	ImageUrl      string         `gorm:"size:256;column:imageurl;not null;" json:"imageurl"`
	ImageWidth    uint           `gorm:"column:imagewidth" json:"imagewidth"`
	ImageHeight   uint           `gorm:"column:imageheight" json:"imageheight"`
	Genres        pq.StringArray `gorm:"type:string[];column:genres" json:"genres"`
	Artists       pq.StringArray `gorm:"type:string[];column:artists" json:"artists"`
	CreatedAt     time.Time      `gorm:"default:CURRENT_TIMESTAMP;column:createdat" json:"createdat"`
	UpdatedAt     time.Time      `gorm:"default:CURRENT_TIMESTAMP;column:updatedat" json:"updatedat"`
	AudioFormatId uint           `gorm:"type:integer;column:audioformatid" json:"audioformatid"`
	AudioFormat   AudioFormat    `gorm:"foreignKey:AudioFormatId"`
	Notes         string         `gorm:"size:60;column:notes;" json:"notes"`
	Catno         string         `gorm:"size:10;column:catno;" json:"catno"`
	Barcode       string         `gorm:"size:20;column:barcode;" json:"barcode"`
	Year          string         `gorm:"size:4;column:year" json:"year"`
	AudioTracks   []AudioTrack   `gorm:"foreignKey:AudioId"`
}

//	func (obj *Video) Prepare() {
//		obj.ID = 0
//		obj.Title = html.EscapeString(strings.TrimSpace(obj.Title))
//		obj.CreatedAt = time.Now()
//		obj.UpdatedAt = time.Now()
//	}
//
//	func (obj *Video) Validate(action string) error {
//		switch strings.ToLower(action) {
//		case "update":
//			if obj.Title == "" {
//				return errors.New("Required Title")
//			}
//			return nil
//		default:
//			if obj.Title == "" {
//				return errors.New("Required Title")
//			}
//			return nil
//		}
//	}
//
//	func (obj *Video) SaveVideo(db *gorm.DB) (*Video, error) {
//		var err error
//		err = db.Debug().Create(&obj).Error
//		if err != nil {
//			return &Video{}, err
//		}
//
//		if obj.ID != 0 {
//			err = db.Debug().Model(&VideoFormat{}).Where("id = ?", obj.VideoFormatId).Take(&obj.VideoFormat).Error
//			if err != nil {
//				return &Video{}, err
//			}
//		}
//
//		return obj, nil
//	}
//
//	func (obj *Video) FindAllVideos(db *gorm.DB) (*[]Video, error) {
//		var err error
//		var videos []Video
//		err = db.Debug().Model(&Video{}).Preload("VideoFormat").Find(&videos).Error
//		if err != nil {
//			return &[]Video{}, err
//		}
//		return &videos, err
//	}
//
//	func (obj *Video) FindVideoByID(db *gorm.DB, uid uint32) (*Video, error) {
//		var err error
//		err = db.Debug().Model(Video{}).Preload("VideoFormat").Where("id = ?", uid).Take(&obj).Error
//		if err != nil {
//			return &Video{}, err
//		}
//		if err == gorm.ErrRecordNotFound {
//			return &Video{}, errors.New("Video Not Found")
//		}
//		return obj, err
//	}
//
//	func (obj *Video) FindAllVideosByTitle(db *gorm.DB, title string) (*[]Video, error) {
//		var err error
//		var videos []Video
//		titleStr := "%" + title + "%"
//		err = db.Debug().Model(&Video{}).Limit(100).Preload("VideoFormat").Where("title LIKE ?", titleStr).Find(&videos).Error
//		if err != nil {
//			return &[]Video{}, err
//		}
//		return &videos, err
//	}
//
//	func (obj *Video) UpdateVideo(db *gorm.DB, uid uint32) (*Video, error) {
//		db = db.Debug().Model(&Video{}).Where("id = ?", uid).Take(&Video{}).UpdateColumns(
//			map[string]interface{}{
//				"title":         obj.Title,
//				"imdbid":        obj.ImdbID,
//				"imageurl":      obj.ImageUrl,
//				"imagewidth":    obj.ImageWidth,
//				"imageheight":   obj.ImageHeight,
//				"runtime":       obj.Runtime,
//				"genres":        obj.Genres,
//				"plot":          obj.Plot,
//				"actors":        obj.Actors,
//				"videoformatid": obj.VideoFormatId,
//				"directors":     obj.Directors,
//			},
//		)
//		if db.Error != nil {
//			return &Video{}, db.Error
//		}
//		// This is the display the updated bike
//		obj, err := obj.FindVideoByID(db, uid)
//		return obj, err
//	}
//
// func (obj *Video) DeleteVideo(db *gorm.DB, uid uint32) (int64, error) {
//
//		db = db.Debug().Model(&Video{}).Where("id = ?", uid).Take(&Video{}).Delete(&Video{})
//
//		if db.Error != nil {
//			return 0, db.Error
//		}
//		return db.RowsAffected, nil
//	}
func (obj *Audio) FindAudioByDiscogsSearch(params url.Values) (*Audio, error) {
	var err error
	//err = db.Debug().Model(Bike{}).Preload("BikeRim").Where("id = ?", uid).Take(&obj).Error
	//if err != nil {
	//	return &Video{}, err
	//}
	//if err == gorm.ErrRecordNotFound {
	//	return &Video{}, errors.New("Video Not Found")
	//}
	rawObj, err := getDiscogsRecord(params)
	if err != nil {
		return &Audio{}, err
	}

	obj.MasterID = rawObj.MasterID
	obj.Title = rawObj.Title
	obj.ImageUrl = rawObj.ImageUrl
	obj.ImageWidth = rawObj.ImageWidth
	obj.ImageHeight = rawObj.ImageHeight
	for _, genre := range rawObj.Genres {
		obj.Genres = append(obj.Genres, strings.ToLower(genre))
	}
	obj.Artists = rawObj.Artists
	obj.Catno = rawObj.Catno
	obj.Barcode = rawObj.Barcode
	obj.Year = rawObj.Year
	obj.AudioTracks = rawObj.AudioTracks

	return obj, err
}
