package main

import (
	"ClockOut/internal/database"
	"ClockOut/internal/logger"
)

func main() {
	logger.Print("main", "Starting...")
	db, err := database.Open("./data/data.db")
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	if err := database.Init(db); err != nil {
		logger.Error("main", "Something went wrong during the database initialization.")
		logger.Fatal(err)
	}

	logger.Print("main", "Done :)")
}
