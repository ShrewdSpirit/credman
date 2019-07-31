package methods

import (
	"github.com/ShrewdSpirit/credman/data"
)

type ResultGetInfo struct {
	Version    string `json:"version"`
	CommitHash string `json:"commithash"`
	Motto      string `json:"motto"`
}

func GetInfo() ResultGetInfo {
	return ResultGetInfo{
		Version:    data.Version,
		CommitHash: data.GitCommit,
		Motto:      "Safeguard your credentials",
	}
}
