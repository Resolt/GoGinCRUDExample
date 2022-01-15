package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type server struct {
	db  *gorm.DB
	gin *gin.Engine
}

func (s *server) setupRoutes() {
	s.gin.GET("/users", s.handleUsersGet())
	s.gin.POST("/users/:name", s.handleUserCreate())
	s.gin.DELETE("/users/:name", s.handleUserDelete())
}

func (s *server) handleUsersGet() gin.HandlerFunc {
	type response struct {
		Users []string `json:"users"`
	}

	return func(c *gin.Context) {
		users := []User{}
		err := s.db.Find(&users).Error
		if err != nil {
			logError(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		userNames := []string{}
		for _, user := range users {
			userNames = append(userNames, user.Name)
		}
		c.JSON(http.StatusOK, response{Users: userNames})
		return
	}
}

func (s *server) handleUserCreate() gin.HandlerFunc {
	type response struct {
		UserID uint `json:"user_id"`
	}

	return func(c *gin.Context) {
		name := c.Param("name")
		user := User{}
		result := s.db.Where("name = ?", name).Find(&user)
		if result.RowsAffected > 0 {
			c.JSON(
				http.StatusConflict,
				gin.H{"detail": "user with name already exists"},
			)
			return
		}
		user.Name = name
		result = s.db.Create(&user)
		if err := result.Error; err != nil {
			logError(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.JSON(http.StatusOK, response{UserID: user.ID})
		return
	}
}

func (s *server) handleUserDelete() gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("name")
		user := User{}
		result := s.db.Where("name = ?", name).Find(&user)
		if result.RowsAffected == 0 {
			c.JSON(
				http.StatusNotFound,
				gin.H{"detail": "user with name does not exists"},
			)
			return
		}
		result = s.db.Delete(&user)
		if err := result.Error; err != nil {
			logError(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.JSON(http.StatusOK, gin.H{"detail": "OK"})
		return
	}
}
