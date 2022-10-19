package models

import (
	"errors"
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Bike struct {
	ID         uint          `gorm:"primary_key;auto_increment" json:"id"`
	Name       string        `gorm:"size:60;column:name;not null;" json:"name"`
	ChainRings pq.Int64Array `gorm:"type:integer[];column:chainrings" json:"chainrings"`
	Cogs       pq.Int64Array `gorm:"type:integer[];column:cogs" json:"cogs"`
	TireWidth  uint          `gorm:"type:integer;column:tirewidth" json:"tirewidth"`
	CreatedAt  time.Time     `gorm:"default:CURRENT_TIMESTAMP;column:createdat" json:"createdat"`
	UpdatedAt  time.Time     `gorm:"default:CURRENT_TIMESTAMP;column:updatedat" json:"updatedat"`
	BikeRimId  uint          `gorm:"type:integer;column:bikerimid" json:"bikerimid"`
	BikeRim    BikeRim       `gorm:"foreignKey:BikeRimId"`
}

//func (u *User) BeforeSave() error {
//	hashedPassword, err := Hash(u.Password)
//	if err != nil {
//		return err
//	}
//	u.Password = string(hashedPassword)
//	return nil
//}
//
//func (u *User) Prepare() {
//	u.ID = 0
//	u.FirstName = html.EscapeString(strings.TrimSpace(u.FirstName))
//	u.MI = html.EscapeString(strings.TrimSpace(u.MI))
//	u.LastName = html.EscapeString(strings.TrimSpace(u.LastName))
//	u.Phone = html.EscapeString(strings.TrimSpace(u.Phone))
//	u.Email = html.EscapeString(strings.TrimSpace(u.Email))
//	u.CreatedAt = time.Now()
//	u.UpdatedAt = time.Now()
//}

func (u *Bike) FindAllBikes(db *gorm.DB) (*[]Bike, error) {
	var err error
	var bikes []Bike
	err = db.Debug().Model(&Bike{}).Limit(100).Preload("BikeRim").Find(&bikes).Error
	if err != nil {
		return &[]Bike{}, err
	}
	return &bikes, err
}

func (u *Bike) FindBikeByID(db *gorm.DB, uid uint32) (*Bike, error) {
	var err error
	err = db.Debug().Model(Bike{}).Preload("BikeRim").Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &Bike{}, err
	}
	if err == gorm.ErrRecordNotFound {
		return &Bike{}, errors.New("Bike Not Found")
	}
	return u, err
}

//func (u *User) Validate(action string) error {
//	switch strings.ToLower(action) {
//	case "update":
//		if u.FirstName == "" {
//			return errors.New("Required First Name")
//		}
//		if u.LastName == "" {
//			return errors.New("Required Last Name")
//		}
//		if u.Password == "" {
//			return errors.New("Required Password")
//		}
//		if u.Email == "" {
//			return errors.New("Required Email")
//		}
//		if err := checkmail.ValidateFormat(u.Email); err != nil {
//			return errors.New("Invalid Email")
//		}
//
//		return nil
//	case "login":
//		if u.Password == "" {
//			return errors.New("Required Password")
//		}
//		if u.Email == "" {
//			return errors.New("Required Email")
//		}
//		if err := checkmail.ValidateFormat(u.Email); err != nil {
//			return errors.New("Invalid Email")
//		}
//		return nil
//
//	default:
//		if u.FirstName == "" {
//			return errors.New("Required First Name")
//		}
//		if u.LastName == "" {
//			return errors.New("Required Last Name")
//		}
//		if u.Password == "" {
//			return errors.New("Required Password")
//		}
//		if u.Email == "" {
//			return errors.New("Required Email")
//		}
//		if err := checkmail.ValidateFormat(u.Email); err != nil {
//			return errors.New("Invalid Email")
//		}
//		return nil
//	}
//}

//func (u *User) SaveUser(db *gorm.DB) (*User, error) {
//	var err error
//	// To hash the password
//	err = u.BeforeSave()
//	if err != nil {
//		log.Fatal(err)
//	}
//	err = db.Debug().Create(&u).Error
//	if err != nil {
//		return &User{}, err
//	}
//	return u, nil
//}

//func (u *User) UpdateAUser(db *gorm.DB, uid uint32) (*User, error) {
//
//	// To hash the password
//	err := u.BeforeSave()
//	if err != nil {
//		log.Fatal(err)
//	}
//	db = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&User{}).UpdateColumns(
//		map[string]interface{}{
//			"password":  u.Password,
//			"firstname": u.FirstName,
//			"mi":        u.MI,
//			"lastname":  u.LastName,
//			"phone":     u.Phone,
//			"email":     u.Email,
//		},
//	)
//	if db.Error != nil {
//		return &User{}, db.Error
//	}
//	// This is the display the updated user
//	err = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&u).Error
//	if err != nil {
//		return &User{}, err
//	}
//	return u, nil
//}

//func (u *User) DeleteAUser(db *gorm.DB, uid uint32) (int64, error) {
//
//	db = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&User{}).Delete(&User{})
//
//	if db.Error != nil {
//		return 0, db.Error
//	}
//	return db.RowsAffected, nil
//}
