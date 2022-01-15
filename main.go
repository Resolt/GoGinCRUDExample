package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	s := server{db: getDBConn()}
	s.setupRoutes()
	if err := s.r.Run(); err != nil {
		log.Fatal(err)
	}
}

func getDBConn() (db *gorm.DB) {
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	fmt.Println(dbPort)
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Copenhagen", dbHost, dbUser, dbPass, dbName, dbPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	return
}

type server struct {
	db *gorm.DB
	r  *gin.Engine
}

func (s *server) setupRoutes() {
	s.r = gin.Default()
	s.r.GET("/basicget", s.getBasicGet())
}

func (s *server) getBasicGet() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, map[string]string{"test": "test"})
	}
}
