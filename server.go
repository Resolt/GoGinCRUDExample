package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type server struct {
	db *gorm.DB
	ge *gin.Engine
	th *taskhandler
}

//Setup the routes of the API
func (s *server) setupRoutes() {
	s.ge.GET("/users", s.handleUsersGet())
	s.ge.POST("/users/:name", s.handleUserCreate())
	s.ge.DELETE("/users/:name", s.handleUserDelete())
	s.ge.GET("/users/:name/posts", s.handleUserPostsGet())
	s.ge.POST("/users/:name/posts/:title", s.handleUserPostCreate())
	s.ge.DELETE("/users/:name/posts/:title", s.handleUserPostDelete())
	s.ge.POST("/tasks/print", s.handleTaskPrint())
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

//Create user endpoint
func (s *server) handleUserCreate() gin.HandlerFunc {
	type response struct {
		UserID uint `json:"user_id"`
	}

	return func(c *gin.Context) {
		name := c.Param("name")
		user := User{Name: name}
		result := s.db.Create(&user)
		if err := result.Error; err != nil {
			if errIsDBUniqueViolation(err) {
				c.JSON(
					http.StatusConflict,
					gin.H{"detail": fmt.Sprintf("user already exists: %s", name)},
				)
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
		if err := result.Error; err != nil {
			logError(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if result.RowsAffected == 0 {
			c.AbortWithStatusJSON(
				http.StatusNotFound,
				gin.H{"detail": fmt.Sprintf("user does not exists: %s", name)},
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
		user := User{}
		result := s.db.Where("name = ?", name).Find(&user)
		if err := result.Error; err != nil {
			logError(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if result.RowsAffected == 0 {
			c.AbortWithStatusJSON(
				http.StatusNotFound,
				gin.H{"detail": fmt.Sprintf("user does not exists: %s", name)},
			)
			return
		}
		posts := []Post{}
		result = s.db.Where("user_id = ?", user.ID).Find(&posts)
		if err := result.Error; err != nil {
			logError(err)
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
		user := User{}
		result := s.db.Where("name = ?", name).Find(&user)
		if err = result.Error; err != nil {
			logError(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if result.RowsAffected == 0 {
			c.AbortWithStatusJSON(
				http.StatusNotFound,
				gin.H{"detail": fmt.Sprintf("user does not exists: %s", name)},
			)
			return
		}
		title := c.Param("title")
		post := Post{
			UserID: user.ID,
			User:   user,
			Title:  title,
			Text:   req.Text,
		}
		result = s.db.Create(&post)
		if err := result.Error; err != nil {
			if errIsDBUniqueViolation(err) {
				c.AbortWithStatusJSON(http.StatusConflict, gin.H{
					"detail": "user already has post with this title",
				})
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

//Delete post endpoint
func (s *server) handleUserPostDelete() gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("name")
		usr := User{}
		result := s.db.Where("name = ?", name).Find(&usr)
		if err := result.Error; err != nil {
			logError(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
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
		if err := result.Error; err != nil {
			logError(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if result.RowsAffected == 0 {
			c.AbortWithStatusJSON(
				http.StatusNotFound,
				gin.H{"detail": "post not found"},
			)
			return
		}
		result = s.db.Delete(&post)
		if err := result.Error; err != nil {
			logError(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.JSON(http.StatusOK, gin.H{"detail": "OK"})
		return
	}
}

func (s *server) handleTaskPrint() gin.HandlerFunc {
	type request struct {
		Task string `json:"task" binding:"required"`
	}

	return func(c *gin.Context) {
		r := request{}
		err := c.ShouldBindJSON(&r)
		if err != nil {
			logError(err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		err = s.th.sendTask(s.th.queueName, r.Task)
		if err != nil {
			logError(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.Status(http.StatusOK)
		return
	}
}
