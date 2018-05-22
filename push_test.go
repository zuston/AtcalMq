package main

import (
	"testing"
	"github.com/zuston/AtcalMq/rabbitmq"
	"github.com/zuston/AtcalMq/util"
	"github.com/streadway/amqp"
	"log"
)


func TestModelMapper(t *testing.T){
	//ParseTranslationTable("./optional.model")
}

func TestPush(t *testing.T){
	option()
	select {

	}
}

var info = "hello world"

// producer demo
func option(){
	testQueueName := "zuston_p"
	configMapper, _ := util.ConfigReader("/opt/mq.cfg")
	mq_uri := configMapper["mq_uri"]
	exchange := configMapper["exchange"]
	exchange_type := configMapper["exchange_type"]
	// send the msg to producer rabbitmq
	pf, _ := rabbitmq.NewProducerFactory(mq_uri,exchange,exchange_type,false)
	// put the msg to channel
	pf.Publish(testQueueName,info)
	go pf.Handle()

	cf, err := rabbitmq.NewConsumerFactory(mq_uri,exchange,exchange_type,false)

	if err!=nil {
		panic("fail to connect to the message queue of rabbitmq")
	}

	cf.Register(testQueueName,pushHandlerTest)

	go cf.Handle()
}


func pushHandlerTest(queue string, msgChan <-chan amqp.Delivery){
	for msg := range msgChan{
		msg.Ack(true)
		if string(msg.Body)==info {
			log.Println("+测试成功")
		}
	}
}