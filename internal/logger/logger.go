package logger

import (
	"ClockOut/internal/utils"
	"os"

	"fmt"
	"log"
	"strings"
)

// Simple Print function, add Debug, silent and so in the future
func Print(context string, args ...string) {
	fmt.Fprintf(os.Stdout, "[INFO][%s][%s] %s.\n", utils.Now(), strings.ToUpper(context), args)
}

// Prints the error to the STDERR
func Error(context string, args ...string) {
	fmt.Fprintf(os.Stderr, "[ERROR][%s][%s] %s.\n", utils.Now(), strings.ToUpper(context), args)
}

func Debug(context string, args ...string) {
	// Make it lowercase for a simpler compare
	debugVar := strings.ToLower(os.Getenv("DEBUG"))
	if debugVar == "true" || debugVar == "1" {
		fmt.Fprintf(os.Stdout, "[DEBUG][%s][%s] %s.\n", utils.Now(), strings.ToUpper(context), args)
	}
}

func Fatal(err error) {
	log.Fatal(err)
}

// !TODO: Make the logger :/
