package http

import (
	"reflect"

	"github.com/gin-gonic/gin"
)

type (
	HandlerType int
	BindFunc    func(c *gin.Context, _ reflect.Type, req *reflect.Value) error
)

const (
	JSON HandlerType = iota
	QUERY
)

func WrapperHandlerJson(controller interface{}) gin.HandlerFunc {
	return wrapperHandler(controller, JSON)
}

func WrapperHandlerQuery(controller interface{}) gin.HandlerFunc {
	return wrapperHandler(controller, QUERY)
}

func wrapperHandler(controller interface{}, typeName HandlerType) gin.HandlerFunc {
	return func(c *gin.Context) {
		var bindFunc BindFunc
		switch typeName {
		case JSON:
			bindFunc = shouldBindJSON
		case QUERY:
			bindFunc = shouldBindQuery
		default:
			bindFunc = shouldBindJSON
		}

		typ := reflect.TypeOf(controller)
		var params = []reflect.Value{}
		if typ.NumIn() > 0 {
			param, err := buildParameter(c, typ.In(0), bindFunc)
			if err != nil {
				ValidateFail(map[string]interface{}{}, c)
				return
			}
			params = append(params, param)
		}

		rets := reflect.ValueOf(controller).Call(params)
		if len(rets) == 0 {
			OkWithData(nil, c)
			return
		} else if len(rets) == 1 {
			if !rets[0].IsNil() {
				FailWithError(rets[0].Interface().(error), c)
				return
			}
		} else {
			if !rets[1].IsNil() {
				FailWithError(rets[1].Interface().(error), c)
				return
			}
		}
		OkWithData(rets[0].Interface(), c)
	}
}

func buildParameter(c *gin.Context, tp reflect.Type, bindFunc BindFunc) (reflect.Value, error) {
	reqStructType := tp.Elem()
	req := reflect.New(reqStructType)
	err := bindFunc(c, reqStructType, &req)
	return reflect.ValueOf(req.Interface()), err
}

func shouldBindJSON(c *gin.Context, _ reflect.Type, req *reflect.Value) error {
	return c.ShouldBindJSON(req.Interface())
}

func shouldBindQuery(c *gin.Context, _ reflect.Type, req *reflect.Value) error {
	return c.ShouldBindQuery(req.Interface())
}
