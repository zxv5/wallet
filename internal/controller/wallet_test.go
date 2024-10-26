package controller

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"wallet/internal/gin/middleware"
	"wallet/internal/gin/middleware/jwt"
	"wallet/internal/service"
	"wallet/internal/types"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupWalletRouter(svcCtx *service.ServiceContext) *gin.Engine {
	walletController := NewWalletController(svcCtx)

	gin.SetMode(gin.TestMode)
	router := gin.Default()

	conf := getTestConfig()
	jwt := jwt.Jwt(&conf.Jwt)
	router.Use(middleware.Base()...)

	router.POST("/wallet/deposit", jwt, walletController.Deposit)
	router.POST("/wallet/withdraw", jwt, walletController.Withdraw)
	router.POST("/wallet/transfer", jwt, walletController.Transfer)
	router.GET("/wallet/balance", jwt, walletController.Balance)
	router.GET("/wallet/record", jwt, walletController.Record)

	return router
}

func getAuthorization(id int64) string {
	conf := getTestConfig()
	token, _ := jwt.Sign(&conf.Jwt, &types.UserInfo{
		ID:        id,
		FirstName: "1",
		LastName:  "1",
	})

	return token
}

func TestWalletControllerDepositSuccess(t *testing.T) {
	var (
		id    int64 = 2000001
		email       = "TestWalletControllerDepositSuccess@email.com"
	)

	svcCtx := setupSvcCtx(t)
	addDBData(t, svcCtx, id, email)

	router := setupWalletRouter(svcCtx)

	depositReq := `{"amount":100}`
	req, _ := http.NewRequest(http.MethodPost, "/wallet/deposit", bytes.NewBuffer([]byte(depositReq)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+getAuthorization(id))

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestWalletControllerDepositFailure(t *testing.T) {
	var (
		id    int64 = 2000002
		email       = "TestWalletControllerDepositFailure@email.com"
	)

	svcCtx := setupSvcCtx(t)
	addDBData(t, svcCtx, id, email)

	router := setupWalletRouter(svcCtx)

	depositReq := `{"amount":-100}`
	req, _ := http.NewRequest(http.MethodPost, "/wallet/deposit", bytes.NewBuffer([]byte(depositReq)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+getAuthorization(id))

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

}

func TestWalletControllerWithdrawSuccess(t *testing.T) {
	var (
		id    int64 = 2000003
		email       = "TestWalletControllerWithdrawSuccess@email.com"
	)

	svcCtx := setupSvcCtx(t)
	addDBData(t, svcCtx, id, email)

	router := setupWalletRouter(svcCtx)

	depositReq := `{"amount":100}`
	req, _ := http.NewRequest(http.MethodPost, "/wallet/deposit", bytes.NewBuffer([]byte(depositReq)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+getAuthorization(id))

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	withdrawReq := `{"amount":50}`
	req, _ = http.NewRequest(http.MethodPost, "/wallet/withdraw", bytes.NewBuffer([]byte(withdrawReq)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+getAuthorization(id))

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestWalletControllerTransferSuccess(t *testing.T) {
	var (
		id     int64 = 2000004
		id2    int64 = 2000005
		email        = "TestWalletControllerTransferSuccess@email.com"
		email2       = "TestWalletControllerTransferSuccess2@email.com"
	)

	svcCtx := setupSvcCtx(t)
	addDBData(t, svcCtx, id, email)
	addDBData(t, svcCtx, id2, email2)

	router := setupWalletRouter(svcCtx)

	depositReq := `{"amount":100}`
	req, _ := http.NewRequest(http.MethodPost, "/wallet/deposit", bytes.NewBuffer([]byte(depositReq)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+getAuthorization(id))

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	transferReq := fmt.Sprintf(`{"amount":30, "target_user_id":%d}`, id2)
	req, _ = http.NewRequest(http.MethodPost, "/wallet/transfer", bytes.NewBuffer([]byte(transferReq)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+getAuthorization(id))

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestWalletControllerBalanceSuccess(t *testing.T) {
	var (
		id    int64 = 2000006
		email       = "TestWalletControllerBalanceSuccess@email.com"
	)

	svcCtx := setupSvcCtx(t)
	addDBData(t, svcCtx, id, email)

	router := setupWalletRouter(svcCtx)

	req, _ := http.NewRequest(http.MethodGet, "/wallet/balance", nil)
	req.Header.Set("Authorization", "Bearer "+getAuthorization(id))

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestWalletControllerRecordSuccess(t *testing.T) {
	var (
		id    int64 = 2000007
		email       = "TestWalletControllerRecordSuccess@email.com"
	)

	svcCtx := setupSvcCtx(t)
	addDBData(t, svcCtx, id, email)

	router := setupWalletRouter(svcCtx)

	req, _ := http.NewRequest(http.MethodGet, "/wallet/record?page=1&size=10", nil)
	req.Header.Set("Authorization", "Bearer "+getAuthorization(id))

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
