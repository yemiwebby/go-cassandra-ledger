package routes

import (
	"net/http"

	"github.com/yemiwebby/go-cassandra-ledger/internal/handlers/health"
	"github.com/yemiwebby/go-cassandra-ledger/internal/handlers/ledger"
)

func LedgerRoutes(mux *http.ServeMux, lh *ledger.LedgerHandler) {
	mux.HandleFunc("/health", health.HealthCheckHandler)
	mux.HandleFunc("/transactions", lh.TransactionHandler)
	mux.HandleFunc("/balance", lh.BalanceHandler)
}
