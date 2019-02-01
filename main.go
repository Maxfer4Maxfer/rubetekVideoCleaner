package main

import (
	"log"
)

func main() {

	// create a new app object
	app, err := newApp()
	if err != nil {
		log.Fatalf("Something goes wrong: %v", err)
	}

	// get cli arguments
	app.initializeCLIArrgs()

	// execute application
	app.start()

}
