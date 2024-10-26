package core

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Mocking the Context
func setupResponseContext() (*Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	return New(c), recorder
}

func TestSendOk(t *testing.T) {
	c, recorder := setupResponseContext()
	c.SendOk("test data")

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.JSONEq(t, `{"code":200,"data":"test data"}`, recorder.Body.String())
}

func TestSendErr(t *testing.T) {
	c, recorder := setupResponseContext()
	c.SendErr("test error")

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.JSONEq(t, `{"code":500,"message":"test error"}`, recorder.Body.String())
}

func TestSendNotFound(t *testing.T) {
	c, recorder := setupResponseContext()
	c.SendNotFound()

	assert.Equal(t, 404, recorder.Code)
	assert.JSONEq(t, `{"code":404,"message":"Not Found"}`, recorder.Body.String())
}

func TestAppend(t *testing.T) {
	c, _ := setupResponseContext()
	c.Append("key", "value")

	assert.Equal(t, "value", c.C.GetStringMap(ExtraProperty)["key"])
}
