package main

import "github.com/zuston/AtcalMq/util"

func main(){
	util.NewConfigReader("/Users/zuston/goDev/src/github.com/zuston/AtcalMq/mq.ini","rabbitmq")
}
