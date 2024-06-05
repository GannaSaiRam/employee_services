package main

import (
	"log"

	"github.com/GannaSaiRam/employee_services/api"
	"github.com/GannaSaiRam/employee_services/employee"
)

// This main function initially intiates database
// Create table if doesn't exist
func main() {
	store, err := employee.NewPGStore()
	if err != nil {
		log.Fatal(err)
	}
	if err := store.Init(); err != nil {
		log.Fatal("Creation of table got failed:", err)
	}
	server := api.StartServer(":8000", store)
	server.Run()
}
