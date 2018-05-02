package main

import (
	"flag"
	"fmt"
)

func main() {
	test := flag.String("test","hahha","")
	flag.Parse()
	fmt.Println(*test)
}