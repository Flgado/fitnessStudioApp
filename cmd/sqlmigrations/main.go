package main

import (
	"database/sql"
	"log"

	"github.com/Flgado/fitnessStudioApp/config"
	"github.com/Flgado/fitnessStudioApp/utils"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	configPath := utils.GetMigrationConfigPath()
	cfgFile, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("LoadConfig: %v", err)
	}

	cfg, err := config.ParseMigrationConfig(cfgFile)
	if err != nil {
		log.Fatalf("ParseConfig: %v", err)
	}

	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		db.Close()
		log.Panic(err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://cmd//sqlmigrations/migrations",
		"postgres", driver)

	if err != nil {
		log.Fatal(err)
	}

	// Attempt to rollback
	if err := m.Down(); err != nil && err != migrate.ErrNoChange {
		log.Printf("Rollback failed: %v", err)
	}

	// Reapply migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}
}
