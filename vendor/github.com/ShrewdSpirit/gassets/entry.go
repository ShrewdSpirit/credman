package gassets

import (
	"bytes"
	"strings"
)

type Entry struct {
	isDir           bool
	name            string
	data            []byte
	children        []*Entry
	parent          *Entry
	overrideAliases []string
}

func NewEntry(isDir bool, name string, data []byte, overrideAliases []string) *Entry {
	return &Entry{
		isDir:           isDir,
		name:            name,
		data:            data,
		overrideAliases: overrideAliases,
	}
}

func (e *Entry) AddChild(entry *Entry) *Entry {
	if e.children == nil {
		e.children = make([]*Entry, 0)
	}

	e.children = append(e.children, entry)
	entry.parent = e

	return entry
}

func (e *Entry) Name() string { return e.name }

func (e *Entry) IsDir() bool { return e.isDir }

func (e *Entry) Ls() []*Entry { return e.children }

func (e *Entry) Bytes() []byte { return e.data }

func (e *Entry) String() string { return string(e.Bytes()) }

func (e *Entry) Reader() *bytes.Reader { return bytes.NewReader(e.Bytes()) }

func (e *Entry) Path() string {
	parts := make([]string, 1)
	parts[0] = e.name

	ce := e
	for ce != nil {
		parts = append(parts, ce.name)
		ce = ce.parent
	}

	for i := len(parts)/2 - 1; i >= 0; i-- {
		opp := len(parts) - 1 - i
		parts[i], parts[opp] = parts[opp], parts[i]
	}

	return strings.Join(parts, "/")
}

func (e *Entry) Get(path string) *Entry {
	path = strings.Trim(path, "/")
	parts := strings.Split(path, "/")

	if e.children == nil || len(e.children) == 0 {
		return nil
	}

	current := e
	for _, part := range parts {
		found := false
		for _, c := range current.children {
			if c.name == part {
				current = c
				found = true
				break
			}
		}

		if !found {
			return nil
		}
	}

	if current == e {
		return nil
	}

	return current
}
