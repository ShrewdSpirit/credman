package browser

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"

	m "github.com/ShrewdSpirit/credman/interfaces/browser/methods"
	"github.com/labstack/echo/v4"
)

var methods = make(map[string]interface{})

type invokeData struct {
	Method string        `json:"method"`
	Args   []interface{} `json:"args"`
}

func init() {
	methods["getinfo"] = m.GetInfo
}

func invokeErrorResult(ctx echo.Context, message string, err error) error {
	return ctx.String(http.StatusBadRequest, fmt.Sprintf("%s: %s", message, err.Error()))
}

func invokeSuccessJsonResult(ctx echo.Context, r interface{}) error {
	return ctx.JSON(http.StatusOK, r)
}

func invokeHandler(ctx echo.Context) error {
	body, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		return invokeErrorResult(ctx, "Failed to read request body", err)
	}

	// encrypted := ctx.Request().Header.Get("Encrypted") == "true"

	var data invokeData
	if err := json.Unmarshal(body, &data); err != nil {
		return invokeErrorResult(ctx, "Invalid request JSON", err)
	}

	methodFunc, ok := methods[data.Method]
	if !ok {
		return invokeErrorResult(ctx, "Invalid method", errors.New("Method "+data.Method+" not found"))
	}

	methodRef := reflect.ValueOf(methodFunc)
	methodType := methodRef.Type()
	numParams := methodType.NumIn()
	if len(data.Args) != numParams {
		return invokeErrorResult(ctx, "Call error", errors.New("Invalid number of arguments to method "+data.Method))
	}

	params := make([]reflect.Value, numParams)
	for i := 0; i < numParams; i++ {
		p := methodType.In(i)
		argType := reflect.TypeOf(data.Args[i])
		if p != argType {
			return invokeErrorResult(ctx, "Call error", errors.New(
				fmt.Sprintf("Invalid argument type for method %s. Expected %s received %s",
					data.Method, p.Name(), argType.Name())))
		}
		params[i] = reflect.ValueOf(data.Args[i])
	}

	results := methodRef.Call(params)
	if len(results) > 0 {
		if err, ok := results[len(results)-1].Interface().(error); ok && err != nil {
			return invokeErrorResult(ctx, "Call error", err)
		}

		if _, ok := results[0].Interface().(error); !ok {
			return invokeSuccessJsonResult(ctx, results[0].Interface())
		}
	}

	return invokeSuccessJsonResult(ctx, map[string]string{})
}
