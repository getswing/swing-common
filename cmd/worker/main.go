package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"getswing.app/player-service/internal/config"
	"getswing.app/player-service/internal/mq"
)

const queueName = "todos.created"

func main() {
	cfg := config.Load()
	log.Printf("[worker] starting log_level=%s", cfg.LogLevel)

	// oopen this if needed for handle consummer
	// gormDB, sqlDB, err := db.ConnectPostgres(cfg)
	// if err != nil {
	// 	log.Fatalf("[worker] db connect error: %v", err)
	// }
	// defer sqlDB.Close()

	// MQ
	conn, ch, err := mq.Connect(cfg.RabbitURL)
	if err != nil {
		log.Fatalf("[worker] rabbitmq connect error: %v", err)
	}
	defer conn.Close()
	defer ch.Close()

	msgs, err := mq.Consume(ch, queueName)
	if err != nil {
		log.Fatalf("[worker] consume error: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Signal handling
	go func() {
		sigC := make(chan os.Signal, 1)
		signal.Notify(sigC, syscall.SIGINT, syscall.SIGTERM)
		<-sigC
		log.Println("[worker] signal received, stopping...")
		cancel()
	}()

	// Process messages
	for {
		select {
		case <-ctx.Done():
			return
		case d, ok := <-msgs:
			if !ok {
				log.Println("[worker] channel closed")
				return
			}

			_ = d.Ack(false)
		}
	}
}
