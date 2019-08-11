package methods

import (
	"github.com/ShrewdSpirit/credman/data"
)

type GetInfo struct{}

func (s GetInfo) Do() (MethodResult, error) {
	return MethodResult{
		"version":    data.Version,
		"commithash": data.GitCommit,
		"motto":      "Safeguard your credentials",
	}, nil
}
