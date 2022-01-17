package main

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func getDBConn() (db *gorm.DB) {
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Copenhagen", dbHost, dbUser, dbPass, dbName, dbPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logFatal(err)
	}
	return
}

func migrate(db *gorm.DB) (err error) {
	if err = db.AutoMigrate(&User{}); err != nil {
		return
	}
	if err = db.AutoMigrate(&Post{}); err != nil {
		return
	}
	return
}

//User model
type User struct {
	gorm.Model
	Name string `gorm:"unique"`
}

func (u *User) BeforeDelete(tx *gorm.DB) (err error) {
	posts := []Post{}
	result := tx.Where("user_id = ?", u.ID).Find(&posts)
	if result.Error != nil && result.RowsAffected != 0 {
		err = result.Error
		return
	}
	if result.RowsAffected > 0 {
		tx.Delete(&posts)
	}
	return
}

//Post model
type Post struct {
	gorm.Model
	User   User
	UserID uint   `gorm:"uniqueIndex:post_idx"`
	Title  string `gorm:"uniqueIndex:post_idx"`
	Text   string
}
