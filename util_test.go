package main

import (
	"testing"
	"fmt"
	"github.com/zuston/AtcalMq/util"
)

func TestLog(t *testing.T) {
	lloger, err := util.NewLogger(util.DEBUG_LEVEL, "/tmp/test.log")

	if err == nil {
		lloger.SetDebug()
		lloger.Error("hello world")
	}
}

func TestWechatNotify(t *testing.T) {
	channel := make(chan int, 10)
	go util.HandlerQueue()
	util.WechatNotify("waht the world is!")
	<-channel
}


func TestHello(t *testing.T){
	fmt.Println("hello world")
}