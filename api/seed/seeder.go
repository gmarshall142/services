package seed

import (
	"github.com/gmarshall142/services/api/models"
	"gorm.io/gorm"
	"log"
)

var users = []models.User{
	models.User{
		FirstName: "Steven",
		LastName:  "Victor",
		Email:     "steven@gmail.com",
		Password:  "password",
	},
	models.User{
		FirstName: "Martin Luther",
		LastName:  "Martin Luther",
		Email:     "luther@gmail.com",
		Password:  "password",
	},
}

func Load(db *gorm.DB) {

	//err := db.Debug().AutoMigrate(&models.User{}).Error
	//if err != nil {
	//	log.Fatalf("cannot migrate table: %v", err)
	//}

	for i, _ := range users {
		err := users[i].BeforeSave()
		if err != nil {
			log.Fatal(err)
		}
		err = db.Debug().Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
		//posts[i].AuthorID = users[i].ID
		//
		//err = db.Debug().Model(&models.Post{}).Create(&posts[i]).Error
		//if err != nil {
		//	log.Fatalf("cannot seed posts table: %v", err)
		//}
	}
}
