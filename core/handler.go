package core

import (
	"github.com/streadway/amqp"
	"fmt"
)

/**
	 handle the mq msg
 */


func CenterLoadHandler(msgChan <-chan amqp.Delivery){
	i := 1
	for msg := range msgChan{
		if i==2 {
			return
		}
		fmt.Println("accept the info : ",i)
		i++
		msg.Ack(false)
	}
}



