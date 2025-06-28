package main

import (
	"log"

	"github.com/yemiwebby/go-cassandra-ledger/internal/api"
	"github.com/yemiwebby/go-cassandra-ledger/internal/handlers/ledger"
	"github.com/yemiwebby/go-cassandra-ledger/internal/service"
	"github.com/yemiwebby/go-cassandra-ledger/internal/store/database/cassandra"
)

func main() {
	log.Println("Starting go-cassandra-ledger...")

	dbStore, err := cassandra.NewCassandraStore()
	if err != nil {
		log.Fatalf("Failed to connect to Cassandra: %v", err)
	}

	engine := service.NewEngine(dbStore)
	ledgerHandler := ledger.NewLedgerHandler(engine)

	// go kafka.StartConsumer() // Kafka consumer goroutine

	api.StartServer(ledgerHandler)
}
