package main

import (
	"log"
	"server2/app/server"
)

func main() {
	server := server.NewMemoryServer()
	log.Fatal(server.ListenAndServe())
}
