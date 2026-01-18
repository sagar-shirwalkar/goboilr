package testdata

import (
	j "encoding/json"
)

// gen:new
type ComplexStruct struct {
	// Embedded field (No tag, but affects constructor)
	Base

	Name string `gen:"get,set"`
	Age  int    `gen:"get,set,val"`

	// Complex types
	Tags     []string     `gen:"get"`
	Metadata j.RawMessage `gen:"get"`
	Ptr      *int         `gen:"set"`
}

type Base struct {
	ID string `gen:"get"`
}
