package service

import (
	"log"
	"os"
	"path/filepath"
	"testing"
	"wallet/internal/config"
	"wallet/pkg/e"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func getTestConfig() config.Config {
	dbConfig := config.Config{
		Server: config.Server{
			RunMode:      "test",
			Host:         "127.0.0.1",
			Port:         3000,
			ReadTimeout:  60,
			WriteTimeout: 60,
		},
		DBCfg: getTestDBConfig(),
		Jwt: config.Jwt{
			Secret: "secret",
			Exp:    3600,
		},
	}

	return dbConfig
}

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

func setupSvcCtx(t *testing.T) *ServiceContext {
	conf := getTestConfig()

	svcCtx := NewServiceContext(&conf)

	t.Cleanup(func() {
		if svcCtx.Dao.DB != nil {
			svcCtx.Dao.Close()
		}
	})

	return svcCtx
}

func addDBData(t *testing.T, svcCtx *ServiceContext, id int64, email string) {
	tx, err := svcCtx.Dao.Begin()
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
		if _, err := svcCtx.Dao.Exec(`DELETE FROM "user" WHERE "id" = $1;`, id); err != nil {
			log.Fatalf("failed to reset test database: %v", err)
		}

		if _, err := svcCtx.Dao.Exec(`DELETE FROM "wallet" WHERE "id" = $1;`, id); err != nil {
			log.Fatalf("failed to reset test database: %v", err)
		}

		if _, err := svcCtx.Dao.Exec(`DELETE FROM "wallet_record" WHERE "id" = $1;`, id); err != nil {
			log.Fatalf("failed to reset test database: %v", err)
		}
	})
}

func TestLoginSuccess(t *testing.T) {
	var (
		id    int64 = 100001
		email       = "TestLoginSuccess@email.com"
	)

	svcCtx := setupSvcCtx(t)
	addDBData(t, svcCtx, id, email)
	ctx, _ := gin.CreateTestContext(nil)

	userService := NewUserService(ctx, svcCtx)
	output, err := userService.Login(email, "123456")

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, email, output.Info.Email)
	assert.NotEmpty(t, output.Token)
}

func TestLoginUserNotFound(t *testing.T) {
	svcCtx := setupSvcCtx(t)
	ctx, _ := gin.CreateTestContext(nil)

	userService := NewUserService(ctx, svcCtx)
	output, err := userService.Login("notfound@example.com", "password")

	assert.Nil(t, output)
	assert.Equal(t, e.NotFound, err)
}

func TestLoginInvalidPassword(t *testing.T) {
	var (
		id    int64 = 100002
		email       = "TestLoginInvalidPassword@email.com"
	)

	svcCtx := setupSvcCtx(t)
	addDBData(t, svcCtx, id, email)
	ctx, _ := gin.CreateTestContext(nil)

	userService := NewUserService(ctx, svcCtx)

	output, err := userService.Login(email, "wrong_password")

	// 7. 断言期望
	assert.Nil(t, output)
	assert.Equal(t, e.UsernameOrPasswordErr, err)
}
