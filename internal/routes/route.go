package routes

import (
	"net/http"

	"github.com/yemiwebby/go-cassandra-ledger/internal/handlers/health"
	"github.com/yemiwebby/go-cassandra-ledger/internal/handlers/ledger"
)

func LedgerRoutes(mux *http.ServeMux, lh *ledger.LedgerHandler) {
	mux.HandleFunc("/healthz", health.HealthCheckHandler)
	mux.HandleFunc("/transaction", lh.TransactionHandler)
	mux.HandleFunc("/balance", lh.BalanceHandler)
}
