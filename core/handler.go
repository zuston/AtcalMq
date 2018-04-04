package core

import (
	"github.com/streadway/amqp"
	"log"
	"fmt"
)


type CenterLoadObj struct {
	PlatformCode string
	ScanType int
	SiteCode string
	ScanTime string
	SiteId int
	Volume float32
	Weight int
	DataType int
	EwbNo string
	EwbsListNo string
	NextSiteCode string
	NextSiteId int
	OperatorCode string
	ScanMan string
	HewbNo string
}

/**
	 consume the mq msg
 */
func CenterLoadHandler(msgChan <-chan amqp.Delivery){
	i := 1
	for msg := range msgChan{
		if i==2 {return}
		log.Printf("accept the info : %d",i)
		log.Printf("got %dB deliveries, [%v]", len(msg.Body),msg.DeliveryTag)
		fmt.Println(string(msg.Body))
		i++
		msg.Ack(false)
	}
}



