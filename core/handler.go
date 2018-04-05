package core

import (
	"github.com/streadway/amqp"
	"log"
	"encoding/json"
	"github.com/tsuna/gohbase"
	"github.com/tsuna/gohbase/hrpc"
	"context"
)

// singleton
type HbaseConn struct{
	Client gohbase.Client
}

var Hconn *HbaseConn

func init(){
	Hconn = &HbaseConn{}
	Hconn.Client = gohbase.NewClient("slave4,slave2,slave3")
}

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
	var clList []CenterLoadObj
	for msg := range msgChan{
		if i==10 {return}
		log.Printf("accept the info : %d",i)
		log.Printf("got %dB deliveries, [%v]", len(msg.Body),msg.DeliveryTag)

		json.Unmarshal([]byte(string(msg.Body)), &clList)

		/**
		todo
		need to maintain the threadPool
		 */
		for _,v := range clList{
			values := map[string]map[string][]byte{"base": map[string][]byte{"ScanMan": []byte(v.ScanMan)}}
			putRequest, err := hrpc.NewPutStr(context.Background(), "centerload", v.EwbsListNo, values)
			if err!=nil {
				log.Printf("hbase put error : %s",err)
				continue
			}
			_, err = Hconn.Client.Put(putRequest)
			if err!=nil {
				log.Printf("hbase finally put error : %s",err)
				continue
			}
		}
		i++
		msg.Ack(false)
	}
}



