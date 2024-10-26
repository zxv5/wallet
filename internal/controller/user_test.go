package controller

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"wallet/internal/config"
	"wallet/internal/service"

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

func setupSvcCtx(t *testing.T) *service.ServiceContext {
	conf := getTestConfig()

	svcCtx := service.NewServiceContext(&conf)

	t.Cleanup(func() {
		if svcCtx.Dao.DB != nil {
			svcCtx.Dao.Close()
		}
	})

	return svcCtx
}

func addDBData(t *testing.T, svcCtx *service.ServiceContext, id int64, email string) {
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

func TestUserControllerLoginSuccess(t *testing.T) {
	var (
		id    int64 = 1000001
		email       = "TestUserControllerLoginSuccess@email.com"
	)

	svcCtx := setupSvcCtx(t)
	addDBData(t, svcCtx, id, email)
	userController := NewUserController(svcCtx)

	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.POST("/user/login", userController.Login)

	loginReq := fmt.Sprintf(`{"email":"%s", "password":"123456"}`, email)
	req, _ := http.NewRequest(http.MethodPost, "/user/login", bytes.NewBuffer([]byte(loginReq)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUserControllerLoginFailure(t *testing.T) {
	var (
		id    int64 = 1000002
		email       = "TestUserControllerLoginFailure@email.com"
	)

	svcCtx := setupSvcCtx(t)
	addDBData(t, svcCtx, id, email)
	userController := NewUserController(svcCtx)

	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.POST("/user/login", userController.Login)

	loginReq := `{"email":"a@email.com", "password":"wrongpassword"}`
	req, _ := http.NewRequest(http.MethodPost, "/user/login", bytes.NewBuffer([]byte(loginReq)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
