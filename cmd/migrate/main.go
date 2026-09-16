package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	"ajirascan/internal/database"
)

func main() {
	directory := flag.String("dir", "migrations", "migration directory")
	flag.Parse()
	url := os.Getenv("DATABASE_URL")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	db, err := database.Open(ctx, url)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if err := database.ApplyMigrations(ctx, db.SQL, *directory); err != nil {
		log.Fatal(err)
	}
	log.Print("migrations applied")
}
