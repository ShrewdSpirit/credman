package methods

import (
	"fmt"
)

type Test struct {
	Arg int `json:"p0"`
}

func (s Test) Do() (MethodResult, error) {
	fmt.Println("Received test number", s.Arg)
	return MethodResult{
		"result": s.Arg * 2,
	}, nil
}
