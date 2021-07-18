package main

import (
	"log"

	"server/internal/server"
)

func main() {
	err := server.Run()
	if err != nil {
		log.Println(err)
	}
}
