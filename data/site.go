package data

type Site struct {
	Name   string            `json:"n"`
	Fields map[string]string `json:"f"`
}

func NewSite(name string) *Site {
	return &Site{
		Name:   name,
		Fields: make(map[string]string),
	}
}
