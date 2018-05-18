package main

import (
	"time"
	"fmt"
)

func main(){
	minTicker := time.NewTicker(time.Second*10)
	for _ = range minTicker.C{
		fmt.Println("hello world")
	}
}
