package ignorerr

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(Ignorerr())
	return r
}

func TestIgnorerrPanicHandled(t *testing.T) {
	router := setupRouter()

	router.GET("/panic", func(c *gin.Context) {
		c.Set("ignorePainc", true)
		panic("test panic")
	})

	req, _ := http.NewRequest("GET", "/panic", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestIgnorerrNoPanic(t *testing.T) {
	router := setupRouter()
	router.GET("/no-panic", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req, _ := http.NewRequest("GET", "/no-panic", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.JSONEq(t, `{"message": "success"}`, resp.Body.String())
}
