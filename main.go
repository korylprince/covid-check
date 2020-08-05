package main

import (
	"log"
	"time"

	"github.com/kelseyhightower/envconfig"
)

func main() {
	config := new(Config)
	if err := envconfig.Process("", config); err != nil {
		log.Fatalln("ERROR: Unable to process configuration:", err)
	}

	for {
		users, err := GetADUsers(config)
		if err != nil {
			log.Println("ERROR: Unable to get AD users:", err)
			continue
		}

		n, err := UpsertUsers(config, users)
		if err != nil {
			log.Println("ERROR: Unable to get upsert users:", err)
			continue
		}

		log.Println("INFO: Synced", n, "users")

		time.Sleep(config.SyncInterval)
	}
}
