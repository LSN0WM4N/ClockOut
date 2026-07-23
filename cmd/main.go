package main

import (
	"ClockOut/internal/core"
	"ClockOut/internal/database"
	"ClockOut/internal/dispatcher"
	"ClockOut/internal/listener"
	"ClockOut/internal/logger"
)

func main() {
	logger.Print("main", "Starting...")
	db, err := database.Open("./data/data.db")
	if err != nil {
		logger.Fatal("main", err)
	}
	defer db.Close()

	if err := database.Init(db); err != nil {
		logger.Error("main", "Something went wrong during the database initialization.")
		logger.Fatal("main", err)
	}

	eventsCh := make(chan core.Event)
	dispatcher := dispatcher.NewDispatcher(eventsCh)

	dispatcher.RegisterHandler("event", func(event core.Event) {
		employeeId := event.GetPayload()
		logger.Print("main", "Catch event [", event.GetType(), "] with payload", employeeId)
	})

	dispatcher.Start()
	defer dispatcher.Stop()

	listener.Init(eventsCh)

	logger.Print("main", "Done :)")
}
