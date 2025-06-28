package ledger

import (
	"fmt"
	"net/http"
)

func (lh *LedgerHandler) BalanceHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "This is for the balance handler")
}
