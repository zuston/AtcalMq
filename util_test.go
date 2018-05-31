package main

import (
	"testing"
	"fmt"
	"github.com/zuston/AtcalMq/util"
)


/**
go test -v
go test -v -test.run TestLog
 */


func TestLog(t *testing.T) {
	lloger, err := util.NewLogger(util.DEBUG_LEVEL, "/tmp/test.log")

	if err == nil {
		lloger.SetDebug()
		lloger.Error("hello world")
	}
}

func TestWechatNotify(t *testing.T) {
	//channel := make(chan int, 10)
	go util.NotifyHandlerQueue()
	util.WechatNotify("waht the world is!")
	//<-channel
}


func TestCofigReader(t *testing.T){
	mapper, _ := util.NewConfigReader("./mq.ini","rabbitmq")
	fmt.Println(mapper)
}
