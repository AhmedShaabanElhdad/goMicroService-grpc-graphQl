package main

import (
	"log"
	"time"

	"github.com/AhmedShaabanElhdad/goMicroService-grpc-graphQl/catalog"
	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
)

type Config struct {
	ElasticSearchUrl string `envconfig:"ELASTICSEARCH_URL"`
	DatabaseUrl      string `envconfig:"DATABASE_URL"`
	PORT             int    `envconfig:"PORT"`
}

func main() {

	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	var repository catalog.Repository
	retry.ForeverSleep(2*time.Second, func(i int) (err error) {
		repository, err = catalog.NewElasticRepository(cfg.ElasticSearchUrl)
		if err != nil {
			log.Println(err)
		}
		return
	})
	defer repository.Close()

	service := catalog.NewService(repository)

	log.Fatal(catalog.ListenAndServeGrpc(service, cfg.PORT))

}
