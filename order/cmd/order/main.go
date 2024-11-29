package main

import (
	"log"
	"time"

	"github.com/AhmedShaabanElhdad/goMicroService-grpc-graphQl/account"
	"github.com/AhmedShaabanElhdad/goMicroService-grpc-graphQl/catalog"
	"github.com/AhmedShaabanElhdad/goMicroService-grpc-graphQl/order"
	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
)

type config struct {
	AccountURL  string `envconfig:"ACCOUNT_SERVICE_URL"`
	CatalogURL  string `envconfig:"CATALOG_SERVICE_URL"`
	DatabaseURL string `envconfig:"DATABASE_URL"`
	PORT        int    `envconfig:"port"`
}

func main() {
	var cfg config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	catalogClient, err := catalog.NewClient(cfg.CatalogURL)
	if err != nil {
		log.Fatal(err)
	}

	accountClient, err := account.NewClient(cfg.AccountURL)
	if err != nil {
		log.Fatal(err)
	}

	var repo order.Repository
	retry.ForeverSleep(2*time.Second, func(i int) error {
		repo, err = order.NewPostgressRepository(cfg.DatabaseURL)
		if err != nil {
			log.Fatal(err)
		}
		return nil
	})
	defer repo.Close()
	log.Println("Listening on port 8080...")

	service := order.NewOrderService(repo)

	log.Fatal(order.ListenAndServe(service, cfg.PORT, *catalogClient, *accountClient))
}
