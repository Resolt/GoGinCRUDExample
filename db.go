package main

import (
	"fmt"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//Get Gorm DB connection
func getDB() (db *gorm.DB, err error) {
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Copenhagen", dbHost, dbUser, dbPass, dbName, dbPort)
	conn := postgres.Open(dsn)
	db, err = gorm.Open(conn, &gorm.Config{})
	return
}

//migrate database models
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
	ID        uint   `gorm:"primarykey"`
	Name      string `gorm:"unique"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

//Gorm hook for deleting posts of user before deleting user
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
	ID        uint   `gorm:"primarykey"`
	User      User   ``
	UserID    uint   `gorm:"uniqueIndex:post_idx"`
	Title     string `gorm:"uniqueIndex:post_idx"`
	Text      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
