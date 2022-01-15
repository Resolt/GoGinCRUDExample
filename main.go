package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	srv := server{
		db:  getDBConn(),
		gin: gin.Default(),
	}
	if err := migrate(srv.db); err != nil {
		logFatal(err)
	}
	srv.setupRoutes()
	if err := srv.gin.Run(); err != nil {
		logFatal(err)
	}
}
