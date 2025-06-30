package ledger

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/yemiwebby/go-cassandra-ledger/internal/models/config"
)

func (lh *LedgerHandler) BalanceHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	start := r.URL.Query().Get("start")
	end := r.URL.Query().Get("end")

	def, ok := lh.BalanceDefinitions[name]
	if !ok {
		http.Error(w, "unknown balance name", http.StatusBadRequest)
		return
	}

	startTs, _ := time.Parse(time.RFC3339, start)
	endTs, _ := time.Parse(time.RFC3339, end)

	var mu sync.Mutex
	var total float64
	var wg sync.WaitGroup
	errCh := make(chan error, len(def.Addresses))

	for _, addr := range def.Addresses {
		wg.Add(1)
		go func(a config.LedgerAddress) {
			defer wg.Done()

			entries, err := lh.Engine.Store.GetEntries(a, def.TimeAxis, startTs, endTs)
			if err != nil {
				errCh <- err
				return
			}

			var sum float64
			for _, e := range entries {
				switch e.Type {
				case "credit":
					sum += e.Amount
				case "debit":
					sum -= e.Amount
				}
			}

			mu.Lock()
			total += sum
			mu.Unlock()
		}(addr)
	}

	wg.Wait()
	select {
	case err := <-errCh:
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	default:
	}

	resp := map[string]interface{}{
		"balance_name": name,
		"start":        start,
		"end":          end,
		"amount":       total,
		"currency":     def.Addresses[0].Currency,
	}
	json.NewEncoder(w).Encode(resp)
}
