package db

import (
	"context"
	"log"
)

var availablePorts = []int64{27015, 27016, 27017, 27018, 27019, 27020, 27021}

func SeedPorts(ctx context.Context, db *Queries) {
	log.Println("Seeding ports...")
	for _, port := range availablePorts {
		db.InsertPort(ctx, port)
	}
}
