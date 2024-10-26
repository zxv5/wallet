package dao

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"testing"
	"wallet/internal/config"
	"wallet/pkg/e"

	"github.com/stretchr/testify/assert"
)

func getTestDBConfig() config.DBCfg {
	cwd, _ := os.Getwd()
	migrationsPath := filepath.Join(cwd, "..", "..", "etc", "migrations")

	dbConfig := config.DBCfg{
		Host:           "43.198.184.4",
		Port:           35333,
		User:           "postgres",
		Password:       "123456",
		DBName:         "wallet_test",
		MaxIdleConns:   10,
		MaxOpenConns:   100,
		MigrationsPath: migrationsPath,
	}

	return dbConfig
}

func setupTestDB(t *testing.T) *Dao {
	dbConfig := getTestDBConfig()
	config := &config.Config{
		DBCfg: dbConfig,
	}
	dao := NewDao(config)

	t.Cleanup(func() {
		if dao.DB != nil {
			dao.Close()
		}
	})

	return dao
}

func addDBData(t *testing.T, dao *Dao, id int64, email string) {
	tx, err := dao.Begin()
	if err != nil {
		t.Fatalf("failed to begin transaction: %v", err)
	}

	userQuery := `INSERT INTO "user" ("id", "email", "password", "first_name", "last_name") VALUES ($1, $2, '123456', 'first_name_a', 'last_name_a');`
	if _, err := tx.Exec(userQuery, id, email); err != nil {
		tx.Rollback()
		t.Fatalf("failed to insert user: %v", err)
	}

	walletQuery := `INSERT INTO "wallet" ("id", "user_id", "balance") VALUES ($1, $2, 0);`
	if _, err := tx.Exec(walletQuery, id, id); err != nil {
		tx.Rollback()
		t.Fatalf("failed to insert wallet: %v", err)
	}

	walletRecordQuery := `INSERT INTO "wallet_record" ("wallet_id", "amount", "transaction_type", "describe") VALUES ($1, $2, $3, $4);`
	if _, err := tx.Exec(walletRecordQuery, id, 1, 1, "describe"); err != nil {
		tx.Rollback()
		t.Fatalf("failed to insert wallet record: %v", err)
	}

	if err := tx.Commit(); err != nil {
		t.Fatalf("failed to commit transaction: %v", err)
	}

	t.Cleanup(func() {
		if _, err := dao.Exec(`DELETE FROM "user" WHERE "id" = $1;`, id); err != nil {
			log.Fatalf("failed to reset test database: %v", err)
		}

		if _, err := dao.Exec(`DELETE FROM "wallet" WHERE "id" = $1;`, id); err != nil {
			log.Fatalf("failed to reset test database: %v", err)
		}

		if _, err := dao.Exec(`DELETE FROM "wallet_record" WHERE "id" = $1;`, id); err != nil {
			log.Fatalf("failed to reset test database: %v", err)
		}
	})
}

func mockTransactionFunction(ctx context.Context) e.Codes {
	return nil
}

func TestTransactionSuccess(t *testing.T) {
	dao := setupTestDB(t)

	ctx := context.Background()
	err := dao.Transaction(ctx, func(ctx context.Context) e.Codes {
		if _, err := dao.WithContext(ctx).Exec(`INSERT INTO "user" ("id", "email", "password", "first_name", "last_name") VALUES ($1, $2, '123456', 'first_name_b', 'last_name_b');`, 10086, "TestTransactionRollback@email.com"); err != nil {
			return e.SQLErr
		}
		return nil
	})
	assert.NoError(t, err)

	var count int
	err2 := dao.WithContext(ctx).QueryRow(`SELECT COUNT(*) FROM "user" WHERE "id" = $1`, 10086).Scan(&count)
	assert.NoError(t, err2)

	assert.Equal(t, 1, count, "Expected one record to be inserted")

	// Clean up: delete the inserted record
	_, err2 = dao.Exec(`DELETE FROM "user" WHERE "id" = $1`, 10086)
	assert.NoError(t, err2)
}

func TestTransactionRollback(t *testing.T) {
	dao := setupTestDB(t)

	ctx := context.Background()
	err := dao.Transaction(ctx, func(ctx context.Context) e.Codes {
		dao.WithContext(ctx).Exec(`INSERT INTO "user" ("id", "email", "password", "first_name", "last_name") VALUES ($1, $2, '123456', 'first_name_b', 'last_name_b');`, 10088, "TestTransactionRollback@email.com")
		return e.SQLErr
	})
	assert.Error(t, err)

	var count int
	err2 := dao.WithContext(ctx).QueryRow(`SELECT COUNT(*) FROM "user" WHERE "id" = $1`, 10088).Scan(&count)
	assert.NoError(t, err2)

	// Assert that the count is 0, meaning the record was not inserted
	assert.Equal(t, 0, count, "Expected no record to be inserted due to rollback")

	if count > 0 {
		_, err2 = dao.Exec(`DELETE FROM "user" WHERE "id" = $1`, 10088)
		if err2 != nil {
			t.Fatalf("failed to delete user: %v", err2)
		}
	}
}
