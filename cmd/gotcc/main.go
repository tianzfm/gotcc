package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"gotcc/internal/config"
	eng "gotcc/internal/engine"
)

func main() {
	// Load configuration (optional)
	_, _ = config.LoadConfig("configs/example.yaml")

	dsn := os.Getenv("FLOW_DB_DSN")
	if dsn == "" {
		log.Fatalf("FLOW_DB_DSN env required (e.g. user:pass@tcp(host:3306)/db)")
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer db.Close()

	loader := eng.SQLFlowLoader(db)
	engine, err := eng.NewEngine(loader)
	if err != nil {
		log.Fatalf("create engine: %v", err)
	}

	// Example: trigger a transaction (replace FLOW_ID and params as needed)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	flowID := os.Getenv("FLOW_ID")
	if flowID != "" {
		if txID, err := engine.ExecuteTransaction(ctx, flowID, map[string]interface{}{"example": true}); err != nil {
			log.Printf("ExecuteTransaction error: %v", err)
		} else {
			log.Printf("ExecuteTransaction succeeded tx=%s", txID)
		}
	} else {
		log.Printf("Engine ready. Set FLOW_ID to run a demo transaction.")
		select {}
	}
}
