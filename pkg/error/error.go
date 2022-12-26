package error

import "fmt"

type CustomError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

var (
	ErrCodeOK = NewCustomError(200, "success")

	ErrCodeInvalidParameter = NewCustomError(400, "parameters error")

	ErrCodeInvalidToken = NewCustomError(403001, "token invalid")
	ErrCodeNoToken      = NewCustomError(403002, "illegal access")
	ErrCodeTokenExpired = NewCustomError(403003, "token expired")

	ErrCodeInternalError = NewCustomError(500, "fail")
)

func (e *CustomError) Error() string {
	return fmt.Sprintf("code: %d, msg: %s", e.Code, e.Msg)
}

func NewCustomError(code int, msg string) *CustomError {
	return &CustomError{
		Code: code,
		Msg:  msg,
	}
}
