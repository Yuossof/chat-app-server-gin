package database

import (
	"log"

	"github.com/Yuossof/messaging-app-server/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := "postgresql://msg-app_owner:npg_yGvQIpKiM0d7@ep-holy-poetry-a511maev-pooler.us-east-2.aws.neon.tech/msg-app?sslmode=require"

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("db connection error", err)
	}

	err = DB.AutoMigrate(&models.User{}, &models.Message{}, &models.DirectConversation{})
	if err != nil {
		log.Fatal("migration error", err)
	}

	log.Println("connected")
}
