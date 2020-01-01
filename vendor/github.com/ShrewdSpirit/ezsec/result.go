package ezsec

import (
	"github.com/gopherjs/gopherjs/js"
)

type result struct {
	*js.Object
}

func newResult() result {
	r := result{
		Object: js.Global.Get("Object").New(),
	}

	r.Object.Set("error", nil)
	r.Object.Set("value", nil)

	return r
}

func (r result) SetError(err error) result {
	r.Object.Set("error", err.Error())
	return r
}

func (r result) SetValue(value interface{}) result {
	r.Object.Set("value", value)
	return r
}
