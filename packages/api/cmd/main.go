package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"api/internal/database/db"
	"api/internal/handler"
	"api/internal/router"
	"api/internal/service"

	_ "github.com/mattn/go-sqlite3"
	"github.com/robfig/cron/v3"
)

func runCron(ctx context.Context, queries *db.Queries, csService *service.CsService) {
	ipAddress := os.Getenv("IP_ADDRESS")
	c := cron.New()

	_, err := c.AddFunc("* * * * *", func() {
		log.Println("Running cron job...")
		servers, _ := queries.GetServers(ctx)

		for _, s := range servers {
			serverStatus, err := csService.GetServerStatus(ctx, service.CsServerStatusPayload{IpAddress: fmt.Sprintf("%s:%d", ipAddress, s.Port)})
			if err != nil {
				err := csService.DestroyServer(ctx, s.ID, s.Port)
				log.Fatalf("Error getting server status %s", err.Error())
			}

			if !s.ExpiresAt.Time.Before(time.Now().UTC()) || s.IsDemo.Bool {
				continue
			}

			if s.ExpiresAt.Time.Before(time.Now().UTC()) && serverStatus.PlayerInfo.Count == 0 {
				log.Printf("Destroying server...")
				err := csService.DestroyServer(ctx, s.ID, s.Port)
				if err != nil {
					fmt.Printf("Error destroying server %s", err.Error())
				}
				continue
			}

			// Extend the server for 30 minutes, since there are players playing on it
			log.Printf("Extending server...")
			queries.UpdateServerExpiration(ctx, db.UpdateServerExpirationParams{
				ExpiresAt: sql.NullTime{Time: time.Now().UTC().Add(30 * time.Minute), Valid: true},
				ID:        s.ID,
			})
		}
	})
	if err != nil {
		fmt.Println("Error scheduling job:", err)
		return
	}

	c.Start()
}

func main() {
	rootPath := os.Getenv("ROOT_PATH")
	dbFile := fmt.Sprintf("%s/packages/api/cs.db", rootPath)
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

	// Make FKs work
	database.Exec("PRAGMA foreign_keys = ON")

	queries := db.New(database)

	db.SeedPorts(ctx, queries)

	csService := service.NewCsService(queries)

	runCron(ctx, queries, csService)

	csHandler := handler.NewCsHandler(csService)

	r := router.NewRouter(csHandler)

	log.Printf("API Server listening on port %d\n", PORT)

	err = http.ListenAndServe(fmt.Sprintf(":%d", PORT), r)
	if err != nil {
		panic("Error starting HTTP server")
	}
}
