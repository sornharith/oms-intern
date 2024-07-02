package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"os"
)

type dataSources struct {
	DB *sqlx.DB
}

// InitDS establishes connections to fields in dataSources
func initDS() (*dataSources, error) {
	//log := logger.Setup()
	//log.Printf("Initializing data sources\n")

	// load env variables
	pgHost := os.Getenv("PG_HOST")
	pgPort := os.Getenv("PG_PORT")
	pgUser := os.Getenv("PG_USER")
	pgPassword := os.Getenv("PG_PASSWORD")
	pgDB := os.Getenv("PG_DB")
	pgSSL := os.Getenv("PG_SSL")

	pgConnString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", pgHost, pgPort, pgUser, pgPassword, pgDB, pgSSL)

	log.Printf("Connecting to PostgreSQL at %s:%s", pgHost, pgPort)
	db, err := sqlx.Open("postgres", pgConnString)

	if err != nil {
		log.Fatal("error opening db:", err)
		return nil, fmt.Errorf("error opening db: %w", err)
	}

	// Verify database connection is working
	if err := db.Ping(); err != nil {
		log.Fatal("error pinging db:", err)
		return nil, fmt.Errorf("error connecting to db: %w", err)
	}

	return &dataSources{
		DB: db,
	}, nil
}

// close to be used in graceful server shutdown
func (d *dataSources) close() error {
	if err := d.DB.Close(); err != nil {
		return fmt.Errorf("error closing PostgreSQL: %w", err)
	}

	return nil
}
