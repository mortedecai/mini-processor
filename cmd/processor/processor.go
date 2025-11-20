package main

import (
	"go.uber.org/zap"

	"github.com/censys/scan-takehome/internal/database"
	"github.com/censys/scan-takehome/internal/processor"

	// import the psql database for the registration side effect
	_ "github.com/censys/scan-takehome/internal/database/psql"
)

func main() {
	zap.ReplaceGlobals(zap.L().Named("processor"))
	db, err := database.New()
	// If database initialization fails, we panic since we can't proceed
	if err != nil {
		panic(err)
	}
	// Ensure the database is closed on exit
	defer db.Close()
	proc, err := processor.New(processor.ConfigFromEnv(), db)
	if err != nil {
		panic(err)
	}
	proc.Start()
}
