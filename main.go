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

	//create server
	srv := createServer(db)

	//setup server routes
	srv.setupRoutes()

	//run server and log error if something goes wrong
	err = srv.gin.Run()
	if err != nil {
		logFatal(err)
	}
	return
}
