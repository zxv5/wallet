package core

import (
	"bytes"
	"io"
	"net/http"
	"reflect"
	"testing"

	"wallet/internal/types"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Mocking the Context
func setupBindContext() (*Context, *gin.Context) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(nil)
	return New(c), c
}

func TestBindWidth(t *testing.T) {
	c, _ := setupBindContext()
	form := new(struct {
		Name string `form:"name"`
	})

	bindFunc := func(obj interface{}) error {
		reflect.ValueOf(obj).Elem().FieldByName("Name").SetString("John")
		return nil
	}

	result := c.BindWidth(form, bindFunc)
	assert.NotNil(t, result)
	assert.Equal(t, "John", form.Name)
}

func TestBindUserInfo(t *testing.T) {
	c, _ := setupBindContext()
	form := &types.UserInfo{}

	c.C.Set("__userinfo", types.UserInfo{ID: 1, FirstName: "John", LastName: "Doe"})

	result := c.BindUserInfo(form)
	assert.NotNil(t, result)
	assert.Equal(t, "John", form.FirstName)
}

func TestBindQuery(t *testing.T) {
	c, _ := setupBindContext()
	form := new(struct {
		QueryParam string `form:"query_param"`
	})

	c.C.Request, _ = http.NewRequest("GET", "/?query_param=test", nil)

	result := c.BindQuery(form)
	assert.NotNil(t, result)
	assert.Equal(t, "test", form.QueryParam)
}

func TestBindUri(t *testing.T) {
	c, _ := setupBindContext()
	form := new(struct {
		UriParam string `uri:"uri_param"`
	})

	c.C.Params = gin.Params{gin.Param{Key: "uri_param", Value: "test-uri"}}

	result := c.BindUri(form)
	assert.NotNil(t, result)
	assert.Equal(t, "test-uri", form.UriParam)
}

func TestBind(t *testing.T) {
	c, _ := setupBindContext()
	form := new(struct {
		Name string `form:"name"`
	})

	body := bytes.NewBufferString(`{"name": "John"}`)
	c.C.Request, _ = http.NewRequest("POST", "/", body)
	c.C.Request.Header.Set("Content-Type", "application/json")
	c.C.Request.Body = io.NopCloser(body)

	result := c.Bind(form)
	assert.NotNil(t, result)
	assert.Equal(t, "John", form.Name)
}
