package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	loop := true

	postgresDSN := os.Getenv("POSTGRES_DSN")

	if postgresDSN == "" {
		postgresDSN = "postgres://postgres:postgres@postgres:5432/postgres?sslmode=disable"
		fmt.Println("missing env variable POSTGRES_DSN, using default value")
	}

	for loop {
		log.Println("Init wait for pg")
		time.Sleep(3 * time.Second)
		db, err := sqlx.Connect("postgres", postgresDSN)
		if err != nil {
			log.Println("PG DSN ", postgresDSN)
			log.Fatalf("Failed connect to postgres: %v", err)
		}

		if _, err := db.Query("SELECT NOW()"); err == nil {
			loop = false
		}
	}
	log.Println("DB ready")
}
