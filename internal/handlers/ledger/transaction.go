package ledger

import (
	"net/http"

	"github.com/yemiwebby/go-cassandra-ledger/internal/models/entry"
	"github.com/yemiwebby/go-cassandra-ledger/internal/models/input"
	"github.com/yemiwebby/go-cassandra-ledger/internal/pkg/utils/handlerutils"
)

func (lh *LedgerHandler) TransactionHandler(w http.ResponseWriter, r *http.Request) {
	// Validate HTTP method and JSON Content Type Header
	if !handlerutils.ValidatePostJSON(w, r) {
		return
	}

	defer r.Body.Close()

	var entrySetInput input.EntrySetInput
	if !handlerutils.DecodeJSONBody(w, r, &entrySetInput) {
		return
	}

	entrySet := entry.NewLedgerEntry(entrySetInput)

	if ok, err := entrySet.IsBalanced(); !ok {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := lh.Engine.ProcessEntrySet(entrySet); err != nil {
		http.Error(w, "failed to process entry set", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("EntrySet recorded\n"))
}
