package main

import (
	"fmt"
	"github.com/zuston/AtcalMq/core"
)

func main() {
	for i:=0;i<20;i++{
		go func() {
			fmt.Println(core.UidGen())
		}()
	}
	select {

	}
}