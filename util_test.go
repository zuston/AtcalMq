package main

import (
	"testing"
	"github.com/zuston/AtcalMq/util"
)

func TestLog(t *testing.T){
	lloger, err := util.NewLogger(1,"/tmp/test.log")

	if err== nil {
		lloger.SetDebug()
		lloger.Error("hello world")
	}
}


func TestWechatNotify(t *testing.T){
	util.WechatNotify("fuck the dog")
}