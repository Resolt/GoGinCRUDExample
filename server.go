package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type server struct {
	db  *gorm.DB
	gin *gin.Engine
}

func createServer(db *gorm.DB) (s *server) {
	s = &server{
		db:  db,
		gin: gin.Default(),
	}
	return
}

//Setup the routes of the API
func (s *server) setupRoutes() {
	s.gin.GET("/users", s.handleUsersGet())
	s.gin.POST("/users/:name", s.handleUserCreate())
	s.gin.DELETE("/users/:name", s.handleUserDelete())
	s.gin.GET("/users/:name/posts", s.handleUserPostsGet())
	s.gin.POST("/users/:name/posts/:title", s.handleUserPostCreate())
	s.gin.DELETE("/users/:name/posts/:title", s.handleUserPostDelete())
}

//Get users endpoint
func (s *server) handleUsersGet() gin.HandlerFunc {
	type response struct {
		Users []string `json:"users"`
	}

	return func(c *gin.Context) {
		users := []User{}
		err := s.db.Find(&users).Error
		if err != nil {
			if strings.Contains(err.Error(), "broken pipe") {
				db, err := getDB()
				if err != nil {
					logError(err)
					c.AbortWithStatus(http.StatusInternalServerError)
					return
				}
				s.db = db
				c.Handler()(c)
				return
			} else {
				logError(err)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
		}
		userNames := []string{}
		for _, user := range users {
			userNames = append(userNames, user.Name)
		}
		c.JSON(http.StatusOK, response{Users: userNames})
		return
	}
}

//Create user endpoint
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
				gin.H{"detail": fmt.Sprintf("user already exists: %s", name)},
			)
			return
		}
		user.Name = name
		err := s.db.Create(&user).Error
		if err != nil {
			if strings.Contains(err.Error(), "broken pipe") {
				db, err := getDB()
				if err != nil {
					logError(err)
					c.AbortWithStatus(http.StatusInternalServerError)
					return
				}
				s.db = db
				c.Handler()(c)
				return
			} else {
				logError(err)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
		}
		c.JSON(http.StatusOK, response{UserID: user.ID})
		return
	}
}

//Delete user endpoint
func (s *server) handleUserDelete() gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("name")
		user := User{}
		result := s.db.Where("name = ?", name).Find(&user)
		if result.RowsAffected == 0 {
			c.AbortWithStatusJSON(
				http.StatusNotFound,
				gin.H{"detail": fmt.Sprintf("user does not exists: %s", name)},
			)
			return
		}
		err := s.db.Delete(&user).Error
		if err != nil {
			if strings.Contains(err.Error(), "broken pipe") {
				db, err := getDB()
				if err != nil {
					logError(err)
					c.AbortWithStatus(http.StatusInternalServerError)
					return
				}
				s.db = db
				c.Handler()(c)
				return
			} else {
				logError(err)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
		}
		c.JSON(http.StatusOK, gin.H{"detail": "OK"})
		return
	}
}

//Get posts of user endpoint
func (s *server) handleUserPostsGet() gin.HandlerFunc {
	type postResponse struct {
		Title string `json:"title" binding:"required"`
		Text  string `json:"text" binding:"required"`
	}
	type response struct {
		Posts []postResponse `json:"posts" binding:"required"`
	}
	return func(c *gin.Context) {
		name := c.Param("name")
		usr := User{}
		result := s.db.Where("name = ?", name).Find(&usr)
		if result.RowsAffected == 0 {
			c.AbortWithStatusJSON(
				http.StatusNotFound,
				gin.H{"detail": fmt.Sprintf("user does not exists: %s", name)},
			)
			return
		}
		posts := []Post{}
		result = s.db.Where("user_id = ?", usr.ID).Find(&posts)
		err := result.Error
		if err != nil && result.RowsAffected != 0 {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		resp := response{}
		for _, post := range posts {
			resp.Posts = append(resp.Posts, postResponse{
				Title: post.Title, Text: post.Text,
			})
		}
		c.JSON(http.StatusOK, resp)
		return
	}
}

//Create post for user endpoint
func (s *server) handleUserPostCreate() gin.HandlerFunc {
	type request struct {
		Text string `json:"text" binding:"required"`
	}
	return func(c *gin.Context) {
		req := request{}
		err := c.ShouldBindJSON(&req)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
				"detail": "post must have fields: title, text",
			})
			return
		}
		name := c.Param("name")
		usr := User{}
		result := s.db.Where("name = ?", name).Find(&usr)
		if result.RowsAffected == 0 {
			c.AbortWithStatusJSON(
				http.StatusNotFound,
				gin.H{"detail": fmt.Sprintf("user does not exists: %s", name)},
			)
			return
		}
		title := c.Param("title")
		post := Post{}
		result = s.db.Where("user_id = ? and title = ?", usr.ID, title).Find(&post)
		if result.RowsAffected > 0 {
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{
				"detail": "user already has post with this title",
			})
			return
		}
		post.UserID = usr.ID
		post.User = usr
		post.Title = title
		post.Text = req.Text
		err = s.db.Create(&post).Error
		if err != nil {
			logError(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.JSON(http.StatusOK, gin.H{"detail": "OK"})
		return
	}
}

//Delete post endpoint
func (s *server) handleUserPostDelete() gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("name")
		usr := User{}
		result := s.db.Where("name = ?", name).Find(&usr)
		if result.RowsAffected == 0 {
			c.AbortWithStatusJSON(
				http.StatusNotFound,
				gin.H{"detail": fmt.Sprintf("user does not exists: %s", name)},
			)
			return
		}
		title := c.Param("title")
		post := Post{}
		result = s.db.Where("user_id = ? and title = ?", usr.ID, title).Find(&post)
		if result.RowsAffected == 0 {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"detail": "post not found"})
			return
		} else if result.Error != nil {
			logError(result.Error)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		result = s.db.Delete(&post)
		err := result.Error
		if err != nil {
			logError(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.JSON(http.StatusOK, gin.H{"detail": "OK"})
		return
	}
}
