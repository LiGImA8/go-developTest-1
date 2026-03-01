package main

import (
	"context"

	"minigate/pkg/config"
	"minigate/pkg/db"
	"minigate/pkg/logging"
	"minigate/services/log/internal"

	"go.uber.org/zap"
)

func main() {
	log := logging.New("log-service")
	defer log.Sync()
	zap.ReplaceGlobals(log)

	mysql, err := db.Open(config.MustDSN())
	if err != nil {
		log.Fatal("open mysql failed", zap.Error(err))
	}
	defer mysql.Close()

	bindEndpoint := config.GetEnv("ZMQ_BIND_ENDPOINT", "tcp://*:5557")
	consumer, err := internal.NewConsumer(mysql, bindEndpoint)
	if err != nil {
		log.Fatal("start zmq consumer failed", zap.Error(err), zap.String("endpoint", bindEndpoint))
	}
	defer consumer.Close()

	log.Info("log service listening for zmq messages", zap.String("endpoint", bindEndpoint))
	if err := consumer.Run(context.Background()); err != nil {
		log.Fatal("consumer stopped", zap.Error(err))
	}
}
