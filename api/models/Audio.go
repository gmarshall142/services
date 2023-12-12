package models

import (
	"errors"
	"fmt"
	"github.com/lib/pq"
	"gorm.io/gorm"
	"html"
	"strings"
	"time"
)

type Audio struct {
	ID            uint           `gorm:"primary_key;auto_increment" json:"id"`
	DiscogsID     uint           `gorm:"column:discogsid;" json:"discogsid"`
	Title         string         `gorm:"size:128;column:title;not null;" json:"title"`
	SortName      string         `gorm:"size:128;column:sortname;not null;" json:"sortname"`
	ImageUrl      string         `gorm:"size:256;column:imageurl;not null;" json:"imageurl"`
	Genres        pq.StringArray `gorm:"type:string[];column:genres" json:"genres"`
	Artists       pq.StringArray `gorm:"type:string[];column:artists" json:"artists"`
	CreatedAt     time.Time      `gorm:"default:CURRENT_TIMESTAMP;column:createdat" json:"createdat"`
	UpdatedAt     time.Time      `gorm:"default:CURRENT_TIMESTAMP;column:updatedat" json:"updatedat"`
	AudioFormatId uint           `gorm:"type:integer;column:audioformatid" json:"audioformatid"`
	AudioFormat   AudioFormat    `gorm:"foreignKey:AudioFormatId"`
	Notes         string         `gorm:"size:60;column:notes;" json:"notes"`
	Catno         string         `gorm:"size:40;column:catno;" json:"catno"`
	Barcode       string         `gorm:"size:60;column:barcode;" json:"barcode"`
	Year          string         `gorm:"size:4;column:year" json:"year"`
	AudioTracks   []AudioTrack   `gorm:"foreignKey:AudioId"`
}

func (obj *Audio) Prepare() {
	obj.ID = 0
	obj.Title = html.EscapeString(strings.TrimSpace(obj.Title))
	obj.CreatedAt = time.Now()
	obj.UpdatedAt = time.Now()
}

func (obj *Audio) Validate(action string) error {
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

func (obj *Audio) SaveAudio(db *gorm.DB) (*Audio, error) {
	var err error
	tracks := obj.AudioTracks
	obj.AudioTracks = []AudioTrack{}
	err = db.Debug().Create(&obj).Error
	if err != nil {
		return &Audio{}, err
	}

	if obj.ID != 0 {
		// Save Tracks
		if len(tracks) > 0 {
			fmt.Println(tracks)
			for _, track := range tracks {
				track.AudioId = obj.ID
				_, err2 := track.SaveAudioTrack(db)
				if err2 != nil {
					return &Audio{}, err2
				}
			}
		}

		err = db.Debug().Model(&AudioFormat{}).Where("id = ?", obj.AudioFormatId).Take(&obj.AudioFormat).Error
		if err != nil {
			return &Audio{}, err
		}
	}

	return obj, nil
}

func (obj *Audio) FindAllAudios(db *gorm.DB) (*[]Audio, error) {
	var err error
	var audios []Audio
	err = db.Debug().Model(&Audio{}).Preload("AudioFormat").Find(&audios).Error
	if err != nil {
		return &[]Audio{}, err
	}
	return &audios, err
}

func (obj *Audio) FindAudioByID(db *gorm.DB, uid uint32) (*Audio, error) {
	var err error
	err = db.Debug().Model(Audio{}).Preload("AudioFormat").Where("id = ?", uid).Take(&obj).Error
	if err != nil {
		return &Audio{}, err
	}
	if err == gorm.ErrRecordNotFound {
		return &Audio{}, errors.New("Audio Not Found")
	}
	return obj, err
}

func (obj *Audio) FindAllAudiosByTitle(db *gorm.DB, title string) (*[]Audio, error) {
	var err error
	var audios []Audio
	titleStr := "%" + title + "%"
	err = db.Debug().Model(&Audio{}).Limit(100).Preload("AudioFormat").Where("title LIKE ?", titleStr).Find(&audios).Error
	if err != nil {
		return &[]Audio{}, err
	}
	return &audios, err
}

func (obj *Audio) UpdateAudio(db *gorm.DB, uid uint32) (*Audio, error) {
	db = db.Debug().Model(&Audio{}).Where("id = ?", uid).Take(&Audio{}).UpdateColumns(
		map[string]interface{}{
			"title":         obj.Title,
			"discogsid":     obj.DiscogsID,
			"imageurl":      obj.ImageUrl,
			"genres":        obj.Genres,
			"artists":       obj.Artists,
			"audioformatid": obj.AudioFormatId,
			"catno":         obj.Catno,
			"barcode":       obj.Barcode,
			"year":          obj.Year,
			"sortname":      obj.SortName,
			"notes":         obj.Notes,
		},
	)
	if db.Error != nil {
		return &Audio{}, db.Error
	}
	// This is the display the updated audio
	obj, err := obj.FindAudioByID(db, uid)
	return obj, err
}

func (obj *Audio) DeleteAudio(db *gorm.DB, uid uint32) (int64, error) {
	db = db.Debug().Model(&Audio{}).Where("id = ?", uid).Take(&Audio{}).Delete(&Audio{})
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

func (obj *Audio) FindAudioByDiscogsSearch(params string) (*Audio, error) {
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

	obj.DiscogsID = rawObj.DiscogsID
	obj.Title = rawObj.Title
	obj.ImageUrl = rawObj.ImageUrl
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
