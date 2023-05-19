package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
	"strings"
)

type Video struct {
	ID          uint           `gorm:"primary_key;auto_increment" json:"id"`
	ImdbID      string         `gorm:"size:9;column:imdb;not null;" json:"imdb"`
	Name        string         `gorm:"size:60;column:name;not null;" json:"name"`
	ImageUrl    string         `gorm:"size:256;column:imageurl;not null;" json:"imageurl"`
	ImageWidth  uint           `gorm:"column:imagewidth" json:"imagewidth"`
	ImageHeight uint           `gorm:"column:imageheight" json:"imageheight"`
	Runtime     uint           `gorm:"column:runtime" json:"runtime"`
	Genres      pq.StringArray `gorm:"type:string[];column:genres" json:"genres"`
	Plot        string         `gorm:"size:512;column:plot;" json:"plot"`
	Actors      pq.StringArray `gorm:"type:string[];column:actors" json:"actors"`
}

//ChainRings pq.Int64Array `gorm:"type:integer[];column:chainrings" json:"chainrings"`
//Cogs       pq.Int64Array `gorm:"type:integer[];column:cogs" json:"cogs"`
//TireWidth  uint          `gorm:"type:integer;column:tirewidth" json:"tirewidth"`
//CreatedAt  time.Time     `gorm:"default:CURRENT_TIMESTAMP;column:createdat" json:"createdat"`
//UpdatedAt  time.Time     `gorm:"default:CURRENT_TIMESTAMP;column:updatedat" json:"updatedat"`
//BikeRimId  uint          `gorm:"type:integer;column:bikerimid" json:"bikerimid"`
//BikeRim    BikeRim       `gorm:"foreignKey:BikeRimId"`

//func (obj *Bike) Prepare() {
//	obj.ID = 0
//	obj.Name = html.EscapeString(strings.TrimSpace(obj.Name))
//	obj.CreatedAt = time.Now()
//	obj.UpdatedAt = time.Now()
//}

//func (obj *Bike) Validate(action string) error {
//	switch strings.ToLower(action) {
//	case "update":
//		if obj.Name == "" {
//			return errors.New("Required Name")
//		}
//		return nil
//	default:
//		if obj.Name == "" {
//			return errors.New("Required Name")
//		}
//		return nil
//	}
//}

//func (obj *Bike) SaveBike(db *gorm.DB) (*Bike, error) {
//	var err error
//	err = db.Debug().Create(&obj).Error
//	if err != nil {
//		return &Bike{}, err
//	}
//
//	if obj.ID != 0 {
//		err = db.Debug().Model(&BikeRim{}).Where("id = ?", obj.BikeRimId).Take(&obj.BikeRim).Error
//		if err != nil {
//			return &Bike{}, err
//		}
//	}
//
//	return obj, nil
//}

//func (obj *Bike) FindAllBikes(db *gorm.DB) (*[]Bike, error) {
//	var err error
//	var bikes []Bike
//	err = db.Debug().Model(&Bike{}).Limit(100).Preload("BikeRim").Find(&bikes).Error
//	if err != nil {
//		return &[]Bike{}, err
//	}
//	return &bikes, err
//}

func (obj *Video) FindVideoByImdbID(db *gorm.DB, id string) (*Video, error) {
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
	obj.Name = rawObj.Name
	obj.ImageUrl = rawObj.ImageUrl
	obj.ImageWidth = rawObj.ImageWidth
	obj.ImageHeight = rawObj.ImageHeight
	obj.Runtime = rawObj.Runtime
	for _, genre := range rawObj.Genres {
		obj.Genres = append(obj.Genres, strings.ToLower(genre))
	}
	obj.Plot = rawObj.Plot
	obj.Actors = rawObj.Actors

	return obj, err
}

//func (obj *Bike) UpdateBike(db *gorm.DB, uid uint32) (*Bike, error) {
//	db = db.Debug().Model(&Bike{}).Where("id = ?", uid).Take(&Bike{}).UpdateColumns(
//		map[string]interface{}{
//			"name":       obj.Name,
//			"chainrings": obj.ChainRings,
//			"cogs":       obj.Cogs,
//			"tirewidth":  obj.TireWidth,
//			"bikerimid":  obj.BikeRimId,
//		},
//	)
//	if db.Error != nil {
//		return &Bike{}, db.Error
//	}
//	// This is the display the updated bike
//	obj, err := obj.FindBikeByID(db, uid)
//	return obj, err
//}

//func (obj *Bike) DeleteBike(db *gorm.DB, uid uint32) (int64, error) {
//
//	db = db.Debug().Model(&Bike{}).Where("id = ?", uid).Take(&Bike{}).Delete(&Bike{})
//
//	if db.Error != nil {
//		return 0, db.Error
//	}
//	return db.RowsAffected, nil
//}
