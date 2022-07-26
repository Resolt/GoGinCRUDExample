package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	ginlogrus "github.com/toorop/gin-logrus"
)

func main() {
	var err error

	lr := logrus.StandardLogger()
	lr.SetFormatter(&logrus.JSONFormatter{})

	//create db connection
	db, err := getDB()
	if err != nil {
		lr.Fatal(err)
	}

	//migrate db
	err = migrate(db)
	if err != nil {
		lr.Fatal(err)
	}

	//connect amqp
	th, err := getTaskhandler(lr)
	if err != nil {
		lr.Fatal(err)
	}

	r := gin.New()
	r.Use(ginlogrus.Logger(lr), gin.Recovery())

	//create server and setup routes
	srv := &server{
		db:  db,
		r:   r,
		th:  th,
		log: lr,
	}
	srv.setupRoutes()

	//run server and log error
	err = srv.r.Run()
	if err != nil {
		lr.Fatal(err)
	}
	return
}
