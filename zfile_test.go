package main

import (
	"testing"
	"fmt"
	"github.com/zuston/AtcalMq/util"
)

func TestWalkDir(t *testing.T){
	fmt.Println(util.WalkDir("/Users/zuston/goDev/src/github.com/zuston/AtcalMq",".model"))
}
