package ledger

import (
	"log"
	"os"
	"path/filepath"

	"github.com/yemiwebby/go-cassandra-ledger/internal/models/config"
	"github.com/yemiwebby/go-cassandra-ledger/internal/service"
)

type LedgerHandler struct {
	Engine             *service.Engine
	BalanceDefinitions map[string]config.ResolvedBalanceDefinition
}

func NewLedgerHandler(engine *service.Engine) *LedgerHandler {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get working dir: %v", err)
	}

	addressPath := filepath.Join(cwd, "configs", "address_config.yml")
	balancePath := filepath.Join(cwd, "configs", "balance_definitions.yml")

	resolvedDefs, err := config.LoadResolvedBalanceDefinitions(addressPath, balancePath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	return &LedgerHandler{
		Engine:             engine,
		BalanceDefinitions: resolvedDefs,
	}
}
