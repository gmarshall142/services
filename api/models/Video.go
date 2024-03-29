package models

import (
	"errors"
	"github.com/lib/pq"
	"gorm.io/gorm"
	"html"
	"strings"
	"time"
)

type Video struct {
	ID            uint           `gorm:"primary_key;auto_increment" json:"id"`
	ImdbID        string         `gorm:"size:9;column:imdbid;not null;" json:"imdbid"`
	Title         string         `gorm:"size:128;column:title;not null;" json:"title"`
	ImageUrl      string         `gorm:"size:256;column:imageurl;not null;" json:"imageurl"`
	ImageWidth    uint           `gorm:"column:imagewidth" json:"imagewidth"`
	ImageHeight   uint           `gorm:"column:imageheight" json:"imageheight"`
	Runtime       uint           `gorm:"column:runtime" json:"runtime"`
	Genres        pq.StringArray `gorm:"type:string[];column:genres" json:"genres"`
	Plot          string         `gorm:"size:1024;column:plot;" json:"plot"`
	Actors        pq.StringArray `gorm:"type:string[];column:actors" json:"actors"`
	Directors     pq.StringArray `gorm:"type:string[];column:directors" json:"directors"`
	VideoFormatId uint           `gorm:"type:integer;column:videoformatid" json:"videoformatid"`
	VideoFormat   VideoFormat    `gorm:"foreignKey:VideoFormatId"`
	CreatedAt     time.Time      `gorm:"default:CURRENT_TIMESTAMP;column:createdat" json:"createdat"`
	UpdatedAt     time.Time      `gorm:"default:CURRENT_TIMESTAMP;column:updatedat" json:"updatedat"`
}

func (obj *Video) Prepare() {
	obj.ID = 0
	obj.Title = html.EscapeString(strings.TrimSpace(obj.Title))
	obj.CreatedAt = time.Now()
	obj.UpdatedAt = time.Now()
}

func (obj *Video) Validate(action string) error {
	switch strings.ToLower(action) {
	case "update":
		if obj.Title == "" {
			return errors.New("Required Title")
		}
		return nil
	default:
		if obj.Title == "" {
			return errors.New("Required Title")
		}
		return nil
	}
}

func (obj *Video) SaveVideo(db *gorm.DB) (*Video, error) {
	var err error
	err = db.Debug().Create(&obj).Error
	if err != nil {
		return &Video{}, err
	}

	if obj.ID != 0 {
		err = db.Debug().Model(&VideoFormat{}).Where("id = ?", obj.VideoFormatId).Take(&obj.VideoFormat).Error
		if err != nil {
			return &Video{}, err
		}
	}

	return obj, nil
}

func (obj *Video) FindAllVideos(db *gorm.DB) (*[]Video, error) {
	var err error
	var videos []Video
	err = db.Debug().Model(&Video{}).Preload("VideoFormat").Find(&videos).Error
	if err != nil {
		return &[]Video{}, err
	}
	return &videos, err
}

func (obj *Video) FindVideoByID(db *gorm.DB, uid uint32) (*Video, error) {
	var err error
	err = db.Debug().Model(Video{}).Preload("VideoFormat").Where("id = ?", uid).Take(&obj).Error
	if err != nil {
		return &Video{}, err
	}
	if err == gorm.ErrRecordNotFound {
		return &Video{}, errors.New("Video Not Found")
	}
	return obj, err
}

func (obj *Video) FindAllVideosByTitle(db *gorm.DB, title string) (*[]Video, error) {
	var err error
	var videos []Video
	titleStr := "%" + title + "%"
	err = db.Debug().Model(&Video{}).Limit(100).Preload("VideoFormat").Where("title LIKE ?", titleStr).Find(&videos).Error
	if err != nil {
		return &[]Video{}, err
	}
	return &videos, err
}

func (obj *Video) UpdateVideo(db *gorm.DB, uid uint32) (*Video, error) {
	db = db.Debug().Model(&Video{}).Where("id = ?", uid).Take(&Video{}).UpdateColumns(
		map[string]interface{}{
			"title":         obj.Title,
			"imdbid":        obj.ImdbID,
			"imageurl":      obj.ImageUrl,
			"imagewidth":    obj.ImageWidth,
			"imageheight":   obj.ImageHeight,
			"runtime":       obj.Runtime,
			"genres":        obj.Genres,
			"plot":          obj.Plot,
			"actors":        obj.Actors,
			"videoformatid": obj.VideoFormatId,
			"directors":     obj.Directors,
		},
	)
	if db.Error != nil {
		return &Video{}, db.Error
	}
	// This is the display the updated bike
	obj, err := obj.FindVideoByID(db, uid)
	return obj, err
}

func (obj *Video) DeleteVideo(db *gorm.DB, uid uint32) (int64, error) {

	db = db.Debug().Model(&Video{}).Where("id = ?", uid).Take(&Video{}).Delete(&Video{})

	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

func (obj *Video) FindVideoByImdbID(id string) (*Video, error) {
	var err error
	//err = db.Debug().Model(Bike{}).Preload("BikeRim").Where("id = ?", uid).Take(&obj).Error
	//if err != nil {
	//	return &Video{}, err
	//}
	//if err == gorm.ErrRecordNotFound {
	//	return &Video{}, errors.New("Video Not Found")
	//}
	rawObj, err := getMoviesDbRecord(id)
	if err != nil {
		return &Video{}, err
	}

	obj.ImdbID = rawObj.ImdbID
	obj.Title = rawObj.Name
	obj.ImageUrl = rawObj.ImageUrl
	obj.ImageWidth = rawObj.ImageWidth
	obj.ImageHeight = rawObj.ImageHeight
	obj.Runtime = rawObj.Runtime
	for _, genre := range rawObj.Genres {
		obj.Genres = append(obj.Genres, strings.ToLower(genre))
	}
	obj.Plot = rawObj.Plot
	obj.Actors = rawObj.Actors
	obj.Directors = rawObj.Directors

	return obj, err
}
