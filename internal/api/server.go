package api

import (
	"log"
	"net/http"

	"github.com/yemiwebby/go-cassandra-ledger/internal/handlers/ledger"
	"github.com/yemiwebby/go-cassandra-ledger/internal/routes"
)

func StartServer(lh *ledger.LedgerHandler) {
	mux := http.NewServeMux()

	routes.LedgerRoutes(mux, lh)

	log.Println("API server running on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
