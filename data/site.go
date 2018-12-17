package data

import (
	"github.com/ShrewdSpirit/credman/utility"
	"github.com/nu7hatch/gouuid"
)

type Site struct {
	Name   string            `json:"n"`
	Id     string            `json:"u"`
	Fields map[string]string `json:"f"`
}

func NewSite(name string) *Site {
	uuid, _ := uuid.NewV4()
	return &Site{
		Name:   name,
		Id:     utility.Hash(uuid.String()),
		Fields: make(map[string]string),
	}
}
