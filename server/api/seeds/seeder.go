package seed

import (
	"github.com/Hsmnasiri/http_monitoring/server/api/models"

	"github.com/jinzhu/gorm"
)

var users = []models.User{
	{
		Nickname: "Steven victor",
		Email:    "steven@gmail.com",
		Password: "password",
	},
	{
		Nickname: "Martin Luther",
		Email:    "luther@gmail.com",
		Password: "password",
	},
}

func Load(db *gorm.DB) {

	// err := db.Debug().DropTableIfExists(&models.Urls{}, &models.User{}).Error
	// if err != nil {
	// 	log.Fatalf("cannot drop table: %v", err)
	// }
	// err = db.Debug().AutoMigrate(&models.User{}, &models.Urls{}).Error
	// if err != nil {
	// 	log.Fatalf("cannot migrate table: %v", err)
	// }

	// err = db.Debug().Model(&models.Urls{}).AddForeignKey("author_id", "users(id)", "cascade", "cascade").Error
	// if err != nil {
	// 	log.Fatalf("attaching foreign key error: %v", err)
	// }

	// for i, _ := range users {
	// 	err = db.Debug().Model(&models.User{}).Create(&users[i]).Error
	// 	if err != nil {
	// 		log.Fatalf("cannot seed users table: %v", err)
	// 	}
	// 	Urls[i].AuthorID = users[i].ID

	// 	err = db.Debug().Model(&models.Urls{}).Create(&Urls[i]).Error
	// 	if err != nil {
	// 		log.Fatalf("cannot seed Urls table: %v", err)
	// 	}
	// }
}
