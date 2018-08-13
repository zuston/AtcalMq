package main

import (
	_ "github.com/ziutek/mymysql/native" // Native engine

	"fmt"
	"github.com/zuston/AtcalMq/core"
)

func main() {
	fmt.Println(core.MixedBasicInfoTableNames)
}