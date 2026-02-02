package main

import (
	"log"

	"ap1-final-mini-moodle/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
