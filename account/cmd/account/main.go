package main

import (
	"log"
	"time"

	"github.com/AhmedShaabanElhdad/goMicroService-grpc-graphQl/account"
	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
)

type Config struct {
	DatabaseUrl string `envconfig:"DATABASE_URL"`
}

func main() {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	// make dependency injection for repository
	var repository account.Repository
	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		repository, err = account.NewPostgressRepository(cfg.DatabaseUrl)
		if err != nil {
			log.Println(err)
		}
		return
	})
	defer repository.Close()
	log.Println("Listening on port 8080...")

	// make dependency injection for Service
	service := account.NewService(repository)

	// start server
	log.Fatal(account.ListenGRPC(service, 8080))
}
