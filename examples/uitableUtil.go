package main

import (
	"github.com/gosuri/uitable"
	"fmt"
	"github.com/gosuri/uilive"
)

func main(){

	writer := uilive.New()
	// start listening for updates and render
	writer.Start()

	table := uitable.New()
	table.MaxColWidth = 80
	table.Wrap = true // wrap columns

	for i:=0;i<10;i++{
		table.AddRow("Name:", i)
		table.AddRow("Birthday:", i+10)
		table.AddRow("Bio:", "distruted system")
		fmt.Fprint(writer,table)
	}
	writer.Stop()
}