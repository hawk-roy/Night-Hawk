package main

import (
	"log"

	"github.com/hawk-roy/Night-Hawk/internal/cache"
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

	redisClient, err := cache.NewRedisClient(cfg.Redis)
	if err != nil {
		log.Fatal("failed to connect redis: ", err)
	}
	defer redisClient.Close()

	log.Println("redis connected successfully")

	r := router.NewRouter(mysqlDB, redisClient)

	log.Println("go-order-service is running on :8500")
	if err := r.Run(":8500"); err != nil {
		log.Fatal(err)
	}
	log.Println("Done")
}
