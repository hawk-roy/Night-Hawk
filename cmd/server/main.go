package main

import (
	"log"

	"github.com/teng-lei-tfs/Night-Hawk/internal/router"
)

func main() {
	r := router.NewRouter()

	log.Println("go-order-service is running on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
