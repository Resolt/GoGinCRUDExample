package main

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

	//create server and setup routes
	srv := createServer(db)
	srv.setupRoutes()

	//run server and log error if something goes wrong
	err = srv.gin.Run()
	if err != nil {
		logFatal(err)
	}
	return
}
