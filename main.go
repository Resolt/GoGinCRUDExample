package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	ginlogrus "github.com/toorop/gin-logrus"
)

func main() {
	var err error

	log := logrus.StandardLogger()

	//create db connection
	db, err := getDB()
	if err != nil {
		log.Fatal(err)
	}

	//migrate db
	err = migrate(db)
	if err != nil {
		log.Fatal(err)
	}

	//connect amqp
	th, err := getTaskhandler()
	if err != nil {
		log.Fatal(err)
	}

	r := gin.New()
	r.Use(ginlogrus.Logger(log), gin.Recovery())

	//create server and setup routes
	srv := &server{
		db:  db,
		ge:  r,
		th:  th,
		log: log,
	}
	srv.setupRoutes()

	//run server and log error if something goes wrong
	err = srv.ge.Run()
	if err != nil {
		log.Fatal(err)
	}
	return
}
