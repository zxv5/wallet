package core

import (
	"errors"
	"reflect"
	"wallet/internal/types"
	"wallet/pkg/e"
	"wallet/pkg/logger"
)

// BindWidth .
func (c *Context) BindWidth(form interface{}, f func(obj interface{}) error) *Context {
	if reflect.TypeOf(form).Kind() != reflect.Ptr {
		c.SendPanic(errors.New("Bind Error, Must Ptr"), e.CheckErr)
	}

	if err := f(form); err != nil {
		logger.Warnf("BindWidth validator Error: %s", err.Error())
		msg := GetValidationError(err, form)
		c.SendPanic(errors.New("Bind Error"), e.ParamsErr, msg)
	}

	return c
}

// Bind .
func (c *Context) Bind(form interface{}) *Context {
	return c.BindWidth(form, c.C.ShouldBind)
}

// BindQuery .
func (c *Context) BindQuery(form interface{}) *Context {
	return c.BindWidth(form, c.C.ShouldBindQuery)
}

// BindUri .
func (c *Context) BindUri(form interface{}) *Context {
	return c.BindWidth(form, c.C.ShouldBindUri)
}

// BindUserInfo
func (c *Context) BindUserInfo(form *types.UserInfo) *Context {
	if userInfo, ok := c.C.Get("__userinfo"); ok {
		form.ID = userInfo.(types.UserInfo).ID
		form.FirstName = userInfo.(types.UserInfo).FirstName
		form.LastName = userInfo.(types.UserInfo).LastName
		form.Gender = userInfo.(types.UserInfo).Gender
		form.Status = userInfo.(types.UserInfo).Status
		form.CreatedAt = userInfo.(types.UserInfo).CreatedAt
	} else {
		c.SendPanic(e.AuthErr)
	}

	return c
}
