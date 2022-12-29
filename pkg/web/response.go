package web

import (
	"net/http"

	e "go-seven/pkg/error"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

func result(err *e.CustomError, data interface{}, c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		err.Code,
		data,
		err.Msg,
	})
}

func Ok(c *gin.Context) {
	result(e.ErrCodeOK, map[string]interface{}{}, c)
}

func OkWithData(data interface{}, c *gin.Context) {
	result(e.ErrCodeOK, data, c)
}

func Fail(c *gin.Context) {
	result(e.ErrCodeInternalError, map[string]interface{}{}, c)
}

func FailWithData(err error, data interface{}, c *gin.Context) {
	errCode, ok := err.(*e.CustomError)
	if !ok {
		errCode = e.ErrCodeInternalError
	}
	result(errCode, data, c)
}

func FailWithError(err error, c *gin.Context) {
	FailWithData(err, map[string]interface{}{}, c)
}

func ValidateFail(data interface{}, c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		e.ErrCodeInvalidParameter.Code,
		data,
		e.ErrCodeInvalidParameter.Msg,
	})
}
