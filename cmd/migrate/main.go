package main

import (
	"context"
	"flag"
	"log"
	"time"

	"ajirascan/internal/config"
	"ajirascan/internal/database"
)

func main() {
	directory := flag.String("dir", "migrations", "migration directory")
	flag.Parse()
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	db, err := database.Open(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if err := database.ApplyMigrations(ctx, db.SQL, *directory); err != nil {
		log.Fatal(err)
	}
	log.Print("migrations applied")
}
