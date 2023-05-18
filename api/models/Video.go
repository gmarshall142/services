package models

import (
	"encoding/json"
	"fmt"
	"golang.org/x/text/currency"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
	"os"
)

type Video struct {
	ID     uint   `gorm:"primary_key;auto_increment" json:"id"`
	ImdbID string `gorm:"size:9;column:imdb;not null;" json:"imdb"`
	Name   string `gorm:"size:60;column:name;not null;" json:"name"`
}

//ChainRings pq.Int64Array `gorm:"type:integer[];column:chainrings" json:"chainrings"`
//Cogs       pq.Int64Array `gorm:"type:integer[];column:cogs" json:"cogs"`
//TireWidth  uint          `gorm:"type:integer;column:tirewidth" json:"tirewidth"`
//CreatedAt  time.Time     `gorm:"default:CURRENT_TIMESTAMP;column:createdat" json:"createdat"`
//UpdatedAt  time.Time     `gorm:"default:CURRENT_TIMESTAMP;column:updatedat" json:"updatedat"`
//BikeRimId  uint          `gorm:"type:integer;column:bikerimid" json:"bikerimid"`
//BikeRim    BikeRim       `gorm:"foreignKey:BikeRimId"`

type TitleText struct {
	Text     string `json:"text"`
	TypeName string `json:"__typename"`
}

type BaseInfo struct {
	ID        string    `json:"id"`
	TitleText TitleText `json:"titleText"`
}

type BaseInfoResults struct {
	Results BaseInfo `json:"results"`
}

type PrincipalCast struct {
	ID currency.Unit
}

type MoviesDB struct {
}

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
	obj, err = getMoviesDbRecord(id)
	if err != nil {
		return &Video{}, err
	}
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

func getMoviesDbRecord(id string) (*Video, error) {
	url := "https://moviesdatabase.p.rapidapi.com/titles/" + id + "?info=base_info"
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Print(err.Error())
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-RapidAPI-Key", os.Getenv("RAPID_API_KEY"))
	req.Header.Add("X-RapidAPI-Host", "moviesdatabase.p.rapidapi.com")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Print(err.Error())
	}

	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Print(err.Error())
	}

	var responseObject BaseInfoResults
	json.Unmarshal(bodyBytes, &responseObject)
	fmt.Printf("API Response as struct %+v\n", responseObject)
	video := Video{}
	video.ImdbID = responseObject.Results.ID
	video.Name = responseObject.Results.TitleText.Text

	//body, err := ioutil.ReadAll(r.Body)
	//if err != nil {
	//	return models.Bike{}, err
	//}
	//bike := models.Bike{}
	//err = json.Unmarshal(body, &bike)
	//if err != nil {
	//	return models.Bike{}, err
	//}
	//bike.Prepare()
	//err = bike.Validate("update")
	//if err != nil {
	//	return models.Bike{}, err
	//}
	//return bike, nil
	return &video, nil
}
