package main

import (
	"fmt"
	"os"
	"time"

	"github.com/lib/pq"
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
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return
	}
	sqlDB, err := db.DB()
	if err != nil {
		return
	}
	sqlDB.SetConnMaxLifetime(time.Hour)
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

//****
//Postgres errors
//****

const (
	dbErrUniqueViolation = pq.ErrorCode("23505")
)

func isErrorCode(err error, errcode pq.ErrorCode) bool {
	if pgerr, ok := err.(*pq.Error); ok {
		return pgerr.Code == errcode
	}
	return false
}

func errIsDbUniqueViolation(err error) bool {
	return isErrorCode(err, dbErrUniqueViolation)
}
