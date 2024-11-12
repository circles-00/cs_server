package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"api/internal/database/db"
	"api/internal/handler"
	"api/internal/router"
	"api/internal/service"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	const dbFile string = "cs.db"
	var ddl string

	PORT := 5000
	ctx := context.Background()

	database, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		panic(err)
	}

	if _, err := database.ExecContext(ctx, ddl); err != nil {
		panic(err)
	}

	queries := db.New(database)
	db.SeedPorts(ctx, queries)

	csService := service.NewCsService(queries)
	csHandler := handler.NewCsHandler(csService)

	r := router.NewRouter(csHandler)

	log.Printf("Server listening on port %d\n", PORT)

	err = http.ListenAndServe(fmt.Sprintf(":%d", PORT), r)
	if err != nil {
		panic("Error starting HTTP server")
	}
}
