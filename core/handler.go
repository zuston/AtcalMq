package core

import (
	"github.com/streadway/amqp"
	"encoding/json"
	"github.com/tsuna/gohbase"
	"github.com/tsuna/gohbase/hrpc"
	"context"
	"github.com/zuston/AtcalMq/util"
	"fmt"
)

// singleton
type HbaseConn struct{
	Client gohbase.Client
}

var Hconn *HbaseConn

var hlogger *util.Logger

func init(){
	Hconn = &HbaseConn{}
	Hconn.Client = gohbase.NewClient("slave4,slave2,slave3")

	hlogger, _ = util.NewLogger(util.DEBUG_LEVEL,"/tmp/handler.log")
	hlogger.SetDebug()
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
		hlogger.Debug("accept the info : %d",i)

		json.Unmarshal([]byte(string(msg.Body)), &clList)
		zlloger.Debug("handler queueName : [%s], unmarshal the json len : %d","centerLoad",len(clList))
		/**
		todo
		need to maintain the threadPool
		 */
		for _,v := range clList{

			// generated code by generate.go file
			columnsMapper := map[string][]byte{
				"PlatformCode":[]byte(v.PlatformCode),
				"ScanType":[]byte(fmt.Sprintf("%d",v.ScanType)),
				"EwbNo":[]byte(v.EwbNo),
				"HewbNo":[]byte(v.HewbNo),
				"NextSiteCode":[]byte(v.NextSiteCode),
				"Weight":[]byte(fmt.Sprintf("%d",v.Weight)),
				"DataType":[]byte(fmt.Sprintf("%d",v.DataType)),
				"OperatorCode":[]byte(v.OperatorCode),
				"SiteCode":[]byte(v.SiteCode),
				"ScanTime":[]byte(v.ScanTime),
				"Volume":[]byte(fmt.Sprintf("%.4f",v.Volume)),
				"SiteId":[]byte(fmt.Sprintf("%d",v.SiteId)),
				"EwbsListNo":[]byte(v.EwbsListNo),
				"NextSiteId":[]byte(fmt.Sprintf("%d",v.NextSiteId)),
				"ScanMan":[]byte(v.ScanMan),
			}


			values := map[string]map[string][]byte{"base": columnsMapper}
			putRequest, err := hrpc.NewPutStr(context.Background(), "centerload", v.EwbNo, values)
			if err!=nil {
				hlogger.Error("hbase put error : %s",err)
				continue
			}
			_, err = Hconn.Client.Put(putRequest)
			if err!=nil {
				hlogger.Error("hbase finally put error : %s",err)
				continue
			}
		}
		i++
		msg.Ack(false)
	}
}



