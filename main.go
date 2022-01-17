package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	var err error

	//create db connection
	db, err := getDB()
	if err != nil {
		logFatal(err)
	}

	//migrate db
	err = migrate(db)
	if err != nil {
		logFatal(err)
	}

	//create server
	srv := server{
		db:  db,
		gin: gin.Default(),
	}

	//setup server routes
	srv.setupRoutes()

	//run server and log error if something goes wrong
	err = srv.gin.Run()
	if err != nil {
		logFatal(err)
	}
	return
}
