package main

import "github.com/gin-gonic/gin"

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

	//connect amqp
	th, err := getTaskhandler()
	if err != nil {
		logFatal(err)
	}

	//create server and setup routes
	srv := &server{
		db: db,
		ge: gin.Default(),
		th: th,
	}
	srv.setupRoutes()

	//run server and log error if something goes wrong
	err = srv.ge.Run()
	if err != nil {
		logFatal(err)
	}
	return
}
