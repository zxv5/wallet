package recovery

import (
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

func setupRouter(logger *zap.SugaredLogger) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(Recovery(logger, true)) // Enable stack trace
	return r
}

func TestRecoveryPanicHandled(t *testing.T) {
	logger := zaptest.NewLogger(t).Sugar()
	router := setupRouter(logger)

	router.GET("/panic", func(c *gin.Context) {
		panic("test panic")
	})

	req, _ := http.NewRequest("GET", "/panic", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
}

func TestRecoveryBrokenPipe(t *testing.T) {
	logger := zaptest.NewLogger(t).Sugar()
	router := setupRouter(logger)

	router.GET("/broken-pipe", func(c *gin.Context) {
		// Simulate broken pipe error
		panic(&net.OpError{Err: &os.SyscallError{Syscall: "write", Err: os.ErrClosed}})
	})

	req, _ := http.NewRequest("GET", "/broken-pipe", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
}

func TestRecoveryNoPanic(t *testing.T) {
	logger := zaptest.NewLogger(t).Sugar()
	router := setupRouter(logger)

	router.GET("/no-panic", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req, _ := http.NewRequest("GET", "/no-panic", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.JSONEq(t, `{"message": "success"}`, resp.Body.String())
}
