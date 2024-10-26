package core

import (
	"net/http"
	"reflect"
	"wallet/pkg/e"
	"wallet/pkg/logger"

	"github.com/gin-gonic/gin"
)

const _BusinessResponse = "_business_response_"
const ExtraProperty = "_extraProperty_"

// response .
type response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// Context .
type Context struct {
	C        *gin.Context
	httpCode int
	response *response
}

// New .
func New(c *gin.Context) *Context {
	res := response{}
	return &Context{C: c, response: &res}
}

// HttpCode .
func (c *Context) HttpCode(httpCode int) *Context {
	c.httpCode = httpCode
	return c
}

// Code .
func (c *Context) Code(code e.Code) *Context {
	c.response.Code = code.Code()
	c.response.Message = code.Message()
	return c
}

// Message .
func (c *Context) Message(message string) *Context {
	c.response.Message = message
	return c
}

// Data .
func (c *Context) Data(data interface{}) *Context {
	c.response.Data = data
	return c
}

// success Normal return
func (c *Context) success(data ...interface{}) {
	if c.httpCode == 0 {
		c.httpCode = 200
	}

	if c.response.Code == 0 {
		c.response.Code = 200
	}

	if len(data) > 0 {
		c.response.Data = data[0]
	}

	if len(data) > 1 {
		logger.Warn("Too many arguments in call to success...", data)
	}

	if c.response.Data == nil {
		c.response.Data = gin.H{}
	}

	res := gin.H{"code": c.response.Code, "data": c.response.Data}

	value := c.C.GetStringMap(ExtraProperty)
	for k, v := range value {
		res[k] = v
	}

	c.C.JSON(c.httpCode, res)
}

// fail Exception Return
func (c *Context) fail(params ...interface{}) {
	c.errMsg(params...)

	if c.httpCode == 0 {
		c.httpCode = 200
	}

	if c.response.Code == 0 {
		c.Code(e.ServerErr)
	}

	// logger.Warnf("http://%s%s -> Response Error: %+v => Params: %+v", c.C.Request.Host, c.C.Request.URL.Path, *c.response, params)

	c.C.Set(_BusinessResponse, c.response)

	res := gin.H{"code": c.response.Code, "message": c.response.Message}
	c.C.AbortWithStatusJSON(c.httpCode, res)
}

/**
 * response success
 */

// SendOk .
func (c *Context) SendOk(data ...interface{}) {
	c.success(data...)
}

// CreateOk .
func (c *Context) CreateOk(data ...interface{}) {
	c.Code(201).success(data...)
}

// DeleteOk .
func (c *Context) DeleteOk(data ...interface{}) {
	c.Code(204).success(data...)
}

/**
 * response error
 */

// SendErr .
func (c *Context) SendErr(params ...interface{}) {
	c.HttpCode(http.StatusBadRequest).fail(params...)
}

// SendNotLogin .
func (c *Context) SendNotLogin(params ...interface{}) {
	args := append(params, e.AuthErr)
	c.HttpCode(http.StatusUnauthorized).fail(args...)
}

// SendNotRole .
func (c *Context) SendNotRole(params ...interface{}) {
	args := append(params, e.PermissionDenied)
	c.HttpCode(http.StatusForbidden).fail(args...)
}

// SendNotFound .
func (c *Context) SendNotFound(params ...interface{}) {
	args := append(params, e.NotFound)
	c.HttpCode(http.StatusNotFound).fail(args...)
}

// SendPanic .
func (c *Context) SendPanic(err error, params ...interface{}) {
	c.C.Set("ignorePainc", true)
	c.HttpCode(http.StatusInternalServerError).fail(params...)
	panic(err)
}

// Append Return to the extended field
func (c *Context) Append(key string, val interface{}) *Context {
	value := c.C.GetStringMap(ExtraProperty)
	if value == nil {
		value = make(map[string]interface{})
	}
	value[key] = val

	c.C.Set(ExtraProperty, value)
	return c
}

// errMsg Error message handling
func (c *Context) errMsg(codeWithMsg ...interface{}) {
	code, message, l := e.ServerErr, "", len(codeWithMsg)
	if c.response.Code != 0 {
		code = e.Int(c.response.Code)
	}

	for _, value := range codeWithMsg {
		switch value := value.(type) {
		case e.Codes:
			code = e.Code(value.Code())
			message = value.Message()
			continue
		case error:
			message = value.Error()
			continue
		}

		v := reflect.ValueOf(value)
		switch v.Kind() {
		case reflect.Int:
			code = e.Code(v.Int())
		case reflect.String:
			message = v.String()
		default:
			logger.Warn("errMsg arguments must be Error, Codes, Int or String...", codeWithMsg)
		}
	}

	if l > 3 {
		logger.Warn("Too many arguments in call to errMsg...", codeWithMsg)
	}

	if code != 0 {
		c.Code(code)
	}

	if message != "" {
		c.Message(message)
	}
}
