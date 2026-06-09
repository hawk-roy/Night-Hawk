package main

import (
	"log"

	"github.com/hawk-roy/Night-Hawk/internal/config"
	"github.com/hawk-roy/Night-Hawk/internal/db"
	"github.com/hawk-roy/Night-Hawk/internal/router"
)

func main() {
	cfg := config.Load()

	mysqlDB, err := db.NewMySQL(cfg.MySQL)
	if err != nil {
		log.Fatal("failed to connect mysql: ", err)
	}
	defer mysqlDB.Close()

	log.Println("mysql connected successfully")

	r := router.NewRouter(mysqlDB)

	log.Println("go-order-service is running on :8500")
	if err := r.Run(":8500"); err != nil {
		log.Fatal(err)
	}
	log.Println("Done")
}
