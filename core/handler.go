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
	"runtime"
	"github.com/ivpusic/grpool"
	"github.com/json-iterator/go"
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

	// 协程池
	numCpus := runtime.NumCPU()
	runtime.GOMAXPROCS(numCpus)
}

/**
装车数据
hbase name : scandata
 */
 const(
 	TABLE_NAME_SCAN_DATA = "scandata"
 	COLUMN_FAMILY_NAME = "base"
 )

/**
	 consume the mq msg
 */
func CenterLoadHandler(msgChan <-chan amqp.Delivery, chanTag chan string){

	qualifierSuffix := "centerLoad"

	// 备份hdfs的文件
	backuper,_ := util.NewLogger(util.INFO_LEVEL,"/tmp/backup/centerload.backup")

	index := 1
	var clList []object.CenterLoadObj

	pool := grpool.NewPool(100, 100)
	defer pool.Release()

	var queueName string
	for msg := range msgChan{
		if index==1 {
			queueName = GetQueueNameFromConsumerTag(msg.ConsumerTag)
		}

		msg.Ack(false)
		hlogger.Debug("ack success")

		// 进入备份，以log的方式
		backuper.Info(string(msg.Body))

		json.Unmarshal([]byte(string(msg.Body)), &clList)
		//hlogger.Debug("handler queueName : [%s], unmarshal the json len : %d","centerLoad",len(clList))
		hlogger.Debug("queue [%s] accept the info : %d, len : %d","centerLoad",index,len(clList))

		for _,v := range clList{

			//insertCellDataQualifier := fmt.Sprintf("%s#%s",v.ScanTime,qualifierSuffix)

			// generated code by generate.go file
			//columnsMapper := map[string][]byte{
			//	"PlatformCode":[]byte(v.PlatformCode),
			//	"ScanType":[]byte(fmt.Sprintf("%d",v.ScanType)),
			//	"EwbNo":[]byte(v.EwbNo),
			//	"HewbNo":[]byte(v.HewbNo),
			//	"NextSiteCode":[]byte(v.NextSiteCode),
			//	"Weight":[]byte(fmt.Sprintf("%d",v.Weight)),
			//	"DataType":[]byte(fmt.Sprintf("%d",v.DataType)),
			//	"OperatorCode":[]byte(v.OperatorCode),
			//	"SiteCode":[]byte(v.SiteCode),
			//	"ScanTime":[]byte(v.ScanTime),
			//	"Volume":[]byte(fmt.Sprintf("%.4f",v.Volume)),
			//	"SiteId":[]byte(fmt.Sprintf("%d",v.SiteId)),
			//	"EwbsListNo":[]byte(v.EwbsListNo),
			//	"NextSiteId":[]byte(fmt.Sprintf("%d",v.NextSiteId)),
			//	"ScanMan":[]byte(v.ScanMan),
			//}

			convertJson,_ := jsoniter.Marshal(v)
			columnsMapper := map[string][]byte{
				qualifierSuffix : []byte(convertJson),
			}

			pool.JobQueue <- func() {
				values := map[string]map[string][]byte{COLUMN_FAMILY_NAME: columnsMapper}
				putRequest, err := hrpc.NewPutStr(context.Background(), TABLE_NAME_SCAN_DATA, v.EwbsListNo, values)
				if err!=nil {
					hlogger.Error("hbase put error : %s",err)
					return
				}
				_, err = Hconn.Client.Put(putRequest)
				if err!=nil {
					hlogger.Error("hbase finally put error : %s",err)
					return
				}
			}

		}

		hlogger.Debug("all put success")
		index++
	}
	hlogger.Info("[%s] stop",queueName)
	// 设置守护进程来监视，否则直接连接
	util.WechatNotify(fmt.Sprintf("[%s] stop",queueName))

	chanTag <- queueName
}




func BizOrderHandler(msgChan <-chan amqp.Delivery, stopChan chan string){

}


func BizEwbHandler(msgChan <-chan amqp.Delivery, stopChan chan string){

}

func BasicRouteHandler(msgChan <-chan amqp.Delivery, stopChan chan string){

}

func BizEwbslistHandler(msgChan <-chan amqp.Delivery, stopChan chan string){

}

func DataSiteLoadHandler(msgChan <-chan amqp.Delivery, stopChan chan string){

}

func DataCenterUnloadHandler(msgChan <-chan amqp.Delivery, stopChan chan string){

	qualifierSuffix := "centerUnLoad"

	backuper,_ := util.NewLogger(util.INFO_LEVEL,"/tmp/backup/centerUnload.backup")

	var culist []object.CenterUnLoadObj

	pool := grpool.NewPool(100, 100)
	defer pool.Release()

	var queueName string

	index := 1
	for msg := range msgChan{
		if index==1 {
			queueName = GetQueueNameFromConsumerTag(msg.ConsumerTag)
		}
		index ++
		msg.Ack(false)

		backuper.Info(string(msg.Body))

		json.Unmarshal([]byte(string(msg.Body)), &culist)

		for _,v := range culist{
			convertJson,_ := jsoniter.Marshal(v)
			columnsMapper := map[string][]byte{
				qualifierSuffix : []byte(convertJson),
			}

			pool.JobQueue <- func() {
				values := map[string]map[string][]byte{COLUMN_FAMILY_NAME: columnsMapper}
				putRequest, err := hrpc.NewPutStr(context.Background(), TABLE_NAME_SCAN_DATA, v.EwbsListNo, values)
				if err!=nil {
					hlogger.Error("[%s] hbase put error : %s",queueName,err)
					return
				}
				_, err = Hconn.Client.Put(putRequest)
				if err!=nil {
					hlogger.Error("[%s] hbase finally put error : %s",queueName,err)
					return
				}
			}
		}
	}
	notifyLine := fmt.Sprintf("[%s] stop",queueName)
	hlogger.Info(notifyLine)
	util.WechatNotify(notifyLine)
	stopChan <- queueName
}


func DataCenterPalletHandler(msgChan <-chan amqp.Delivery, stopChan chan string){

}


func DataCenterSortHandler(msgChan <-chan amqp.Delivery, stopChan chan string){

}


func BasicAreaHandler(msgChan <-chan amqp.Delivery, stopChan chan string){

}

func BasicAttendHandler(msgChan <-chan amqp.Delivery, stopChan chan string){

}


func TriggerSiteSendHandler(msgChan <-chan amqp.Delivery, stopChan chan string){

}

func TriggerSiteUploadHandler(msgChan <-chan amqp.Delivery, stopChan chan string){

}

func TriggerInOrOutHandler(msgChan <-chan amqp.Delivery, stopChan chan string){

}


func TriggerStayOrLeaveHandler(msgChan <-chan amqp.Delivery, stopChan chan string){

}

func TriggerCenterTransportHandler(msgChan <-chan amqp.Delivery, stopChan chan string){

}

func BasicSiteHandler(msgChan <-chan amqp.Delivery, stopChan chan string){

}

func BasicVehicleLineHandler(msgChan <-chan amqp.Delivery, stopChan chan string){

}


func BasicPlatFormHandler(msgChan <-chan amqp.Delivery, stopChan chan string){

}






