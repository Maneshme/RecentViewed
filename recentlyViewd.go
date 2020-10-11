package main

import (
	"log"

	"recentViewed/controller"
)

func main() {
	log.Println("Running recently viewed service")
	controller.HandleRequests()
}
