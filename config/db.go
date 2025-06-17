package config

import (
	"fmt"

	"github.com/KrittayotToin/quote-vote-backend/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() *gorm.DB {
	dsn := "host=localhost user=quoteuser password=password123 dbname=quote_vote port=5432 sslmode=disable TimeZone=Asia/Bangkok"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Println("Failed to connect database", err)
		panic("failed to connect database")
	}

	fmt.Println("Database connected")
	db.AutoMigrate(&model.User{}, &model.Quote{}, &model.Vote{})
	return db
}
