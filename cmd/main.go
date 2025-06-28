package main

import (
	"log"

	"github.com/yemiwebby/go-cassandra-ledger/internal/kafka"
)

func main() {
	log.Println("Starting go-cassandra-ledger...")

	go kafka.StartConsumer() // Kafka consumer goroutine
	
}
