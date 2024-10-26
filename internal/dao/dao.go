package dao

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"wallet/internal/config"

	"wallet/pkg/e"
)

type ctxTransactionKey struct{}

type Dao struct {
	config *config.Config
	*sql.DB
}

func NewDao(c *config.Config) *Dao {
	dbConfig := c.DBCfg
	createDatabaseIfNotExists(&dbConfig)

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=UTC",
		dbConfig.Host, dbConfig.User, dbConfig.Password, dbConfig.DBName, dbConfig.Port)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("open postgresql fail: %v", err)
	}

	db.SetMaxIdleConns(dbConfig.MaxIdleConns)
	db.SetMaxOpenConns(dbConfig.MaxOpenConns)

	migrations(db, dbConfig.MigrationsPath)

	return &Dao{
		config: c,
		DB:     db,
	}
}

// Define a unified database interface
type Database interface {
	QueryRow(query string, args ...interface{}) *sql.Row
	Query(query string, args ...interface{}) (*sql.Rows, error)
	Exec(query string, args ...any) (sql.Result, error)
}

// DBWrapper is used to encapsulate *sql.DB
type DBWrapper struct {
	*sql.DB
}

// TxWrapper is used to encapsulate *sql.Tx
type TxWrapper struct {
	*sql.Tx
}

func (db *DBWrapper) QueryRow(query string, args ...interface{}) *sql.Row {
	return db.DB.QueryRow(query, args...)
}

func (tx *TxWrapper) QueryRow(query string, args ...interface{}) *sql.Row {
	return tx.Tx.QueryRow(query, args...)
}

func (db *DBWrapper) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return db.DB.Query(query, args...)
}

func (tx *TxWrapper) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return tx.Tx.Query(query, args...)
}

func (db *DBWrapper) Exec(query string, args ...any) (sql.Result, error) {
	return db.DB.Exec(query, args...)
}

func (tx *TxWrapper) Exec(query string, args ...any) (sql.Result, error) {
	return tx.Tx.Exec(query, args...)
}

func (dao *Dao) WithContext(ctx context.Context) Database {
	iface := ctx.Value(ctxTransactionKey{})

	if iface != nil {
		if tx, ok := iface.(*sql.Tx); ok {
			return &TxWrapper{tx} // Returns the transaction wrapper
		}
	}

	return &DBWrapper{dao.DB} // Returns the database wrapper

}

func (dao *Dao) Transaction(ctx context.Context, runInTransaction func(ctx context.Context) e.Codes) e.Codes {
	iface := ctx.Value(ctxTransactionKey{})

	if iface != nil {
		if tx, ok := iface.(*sql.Tx); ok {
			if err := runInTransaction(ctx); err != nil {
				tx.Rollback()
				return err
			}
			return nil
		}
	}

	tx, err := dao.Begin()
	if err != nil {
		return e.SQLErr
	}

	ctx = context.WithValue(ctx, ctxTransactionKey{}, tx)

	if err := runInTransaction(ctx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}
