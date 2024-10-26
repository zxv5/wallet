package dao

import (
	"database/sql"
	"fmt"
	"wallet/internal/config"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func createDatabaseIfNotExists(dbConfig *config.DBCfg) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=template1 port=%d sslmode=disable TimeZone=UTC",
		dbConfig.Host, dbConfig.User, dbConfig.Password, dbConfig.Port)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM pg_database WHERE datname = $1);`

	err = db.QueryRow(query, dbConfig.DBName).Scan(&exists)
	if err != nil {
		panic("Unable to create database if not exists")
	}

	if !exists {
		_, err := db.Exec("CREATE DATABASE " + dbConfig.DBName + " ENCODING 'UTF8' LC_COLLATE = 'en_US.utf8' LC_CTYPE = 'en_US.utf8' TEMPLATE = template0;")
		if err != nil {
			panic("Unable to create database if not exists")
		}
	}
}

func migrations(db *sql.DB, filePath string) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		panic("Migrations WithInstance error" + err.Error())
	}

	m, err := migrate.NewWithDatabaseInstance("file://"+filePath, "postgres", driver)
	if err != nil {
		panic("Migrations error" + err.Error())
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		panic("Migrations Up error" + err.Error())
	}
}
