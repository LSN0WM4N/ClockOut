package main

import (
	"ClockOut/internal/database"
	"ClockOut/internal/logger"
)

func main() {
	logger.Print("MAIN", "Starting...")
	db, err := database.Open("./data/data.db")
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	if err := database.Init(db); err != nil {
		logger.Fatal(err)
	}

	if err != nil {
		logger.Fatal(err)
	}

	logger.Print("MAIN", "Done :)")
}
