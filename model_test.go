package main

import (
	"testing"
	"github.com/zuston/AtcalMq/core"
)

type Person struct {
	Name string
	Age int
	Year int
}

func TestModelMapGen(t *testing.T){
	v := &Person{
		Name:"zuston",
		Age:24,
		Year:1993,
	}
	mapper := core.ModelMapperGen(v)
	t.Log(string(mapper["Year"]),string(mapper["Age"]),string(mapper["Name"]))
}

func TestUidGen(t *testing.T){
	core.UidGen()
}

func ModelGen(t *testing.T){
	jsonline := `[{"name":"zuston","age":123,"period":1.23}]`
	core.ModelGen([]byte(jsonline))
}