package core

import (
	"github.com/streadway/amqp"
	"github.com/tsuna/gohbase"
	"github.com/zuston/AtcalMq/util"
	"github.com/zuston/AtcalMq/core/object"
	"encoding/json"
	"fmt"
	"github.com/tsuna/gohbase/hrpc"
	"context"
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

	hlogger, _ = util.NewLogger(util.INFO_LEVEL,"/tmp/handler.log")
	hlogger.SetDebug()
}


/**
	 consume the mq msg
 */
func CenterLoadHandler(msgChan <-chan amqp.Delivery){

	//i := 1
	//for d := range msgChan {
	//	log.Printf("index : %d",i)
	//	log.Printf(
	//		"got %dB delivery: [%v]",
	//		len(d.Body),
	//		d.DeliveryTag,
	//	)
	//	d.Ack(false)
	//	i++
	//}
	//log.Printf("handle: deliveries channel closed")
	//return


	// 备份hdfs的文件
	backuper,_ := util.NewLogger(util.INFO_LEVEL,"/tmp/backup/centerload.backup")

	i := 1
	var clList []object.CenterLoadObj

	for msg := range msgChan{
		//if i==10 {return}

		msg.Ack(false)
		hlogger.Debug("ack success")

		// 进入备份，以log的方式
		backuper.Info(string(msg.Body))

		json.Unmarshal([]byte(string(msg.Body)), &clList)
		//hlogger.Debug("handler queueName : [%s], unmarshal the json len : %d","centerLoad",len(clList))
		hlogger.Debug("queue [%s] accept the info : %d, len : %d","centerLoad",i,len(clList))


		for k,v := range clList{

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
			putRequest, err := hrpc.NewPutStr(context.Background(), "centerload", fmt.Sprintf("%s#%s",v.ScanTime,v.EwbNo), values)
			if err!=nil {
				hlogger.Error("hbase put error : %s",err)
				continue
			}
			_, err = Hconn.Client.Put(putRequest)
			if err!=nil {
				hlogger.Error("hbase finally put error : %s",err)
				continue
			}
			// 埋点监控
			HandlerAbility(QUEUE_CENTERLOAD)
			hlogger.Debug("%d put",k)
		}

		hlogger.Debug("all put success")
		i++
	}


}



func BizOrderHandler(msgChan <-chan amqp.Delivery){

}


func BizEwbHandler(msgChan <-chan amqp.Delivery){

}

func BasicRouteHandler(msgChan <-chan amqp.Delivery){

}

func BizEwbslistHandler(msgChan <-chan amqp.Delivery){

}

func DataSiteLoadHandler(msgChan <-chan amqp.Delivery){

}

func DataCenterUnloadHandler(msgChan <-chan amqp.Delivery){

}


func DataCenterPalletHandler(msgChan <-chan amqp.Delivery){

}


func DataCenterSortHandler(msgChan <-chan amqp.Delivery){

}


func BasicAreaHandler(msgChan <-chan amqp.Delivery){

}

func BasicAttendHandler(msgChan <-chan amqp.Delivery){

}


func TriggerSiteSendHandler(msgChan <-chan amqp.Delivery){

}

func TriggerSiteUploadHandler(msgChan <-chan amqp.Delivery){

}

func TriggerInOrOutHandler(msgChan <-chan amqp.Delivery){

}


func TriggerStayOrLeaveHandler(msgChan <-chan amqp.Delivery){

}

func TriggerCenterTransportHandler(msgChan <-chan amqp.Delivery){

}

func BasicSiteHandler(msgChan <-chan amqp.Delivery){

}

func BasicVehicleLineHandler(msgChan <-chan amqp.Delivery){

}


func BasicPlatFormHandler(msgChan <-chan amqp.Delivery){

}






