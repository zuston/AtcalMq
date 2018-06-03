package main

import (
	"github.com/zuston/AtcalMq/util"
	"github.com/zuston/AtcalMq/core/pullcore"
)

func main(){
	util.NewConfigReader("/Users/zuston/goDev/src/github.com/zuston/AtcalMq/core/pullcore/hbase.ini","production")

	pullcore.InitHconn()
}
