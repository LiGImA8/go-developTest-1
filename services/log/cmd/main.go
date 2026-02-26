package main

import (
	"context"
	"log"

	"minigate/pkg/config"
	"minigate/pkg/db"
	"minigate/services/log/internal"
)

func main() {
	mysql, err := db.Open(config.MustDSN())
	if err != nil {
		log.Fatal(err)
	}
	defer mysql.Close()

	consumer, err := internal.NewConsumer(mysql, config.GetEnv("ZMQ_BIND_ENDPOINT", "tcp://*:5557"))
	if err != nil {
		log.Fatal(err)
	}
	defer consumer.Close()

	log.Println("log service listening for zmq messages")
	log.Fatal(consumer.Run(context.Background()))
}
