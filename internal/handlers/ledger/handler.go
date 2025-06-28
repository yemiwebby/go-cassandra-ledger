package ledger

import "github.com/yemiwebby/go-cassandra-ledger/internal/service"

type LedgerHandler struct {
	Engine *service.Engine
}

func NewLedgerHandler(engine *service.Engine) *LedgerHandler {
	return &LedgerHandler{Engine: engine}
}
