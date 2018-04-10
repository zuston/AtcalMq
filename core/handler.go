package core

import (
	"github.com/streadway/amqp"
	"github.com/tsuna/gohbase"
	"github.com/zuston/AtcalMq/util"
	"github.com/zuston/AtcalMq/core/object"
	"encoding/json"
	"runtime"
	"github.com/ivpusic/grpool"
)

// singleton
type HbaseConn struct{
	Client gohbase.Client
}

// 人，货，车，站点
// 人 ： columnFamily----> uid:[uid]
// ewoNo -----> 插入 queueName(columnFamily) ----> uid:[uid]

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

 	TABLE_NAME_BASIC_ROUTE = "route"
 )

/**
	 consume the mq msg
 */
func CenterLoadHandler(msgChan <-chan amqp.Delivery){

	// 备份hdfs的文件
	backuper,_ := util.NewLogger(util.INFO_LEVEL,"/tmp/backup/centerload.backup")

	var clList []object.CenterLoadObj

	pool := grpool.NewPool(100, 100)
	defer pool.Release()

	for msg := range msgChan{
		msg.Ack(false)

		backuper.Info(string(msg.Body))

		json.Unmarshal([]byte(string(msg.Body)), &clList)

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


			pool.JobQueue <- SaveModelGen(v,QUEUE_DATA_CENTERLOAD)

		}

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

	backuper,_ := util.NewLogger(util.INFO_LEVEL,"/tmp/backup/centerUnload.backup")

	var culist []object.CenterUnLoadObj

	pool := grpool.NewPool(100, 100)
	defer pool.Release()

	for msg := range msgChan{
		msg.Ack(false)
		backuper.Info(string(msg.Body))
		json.Unmarshal([]byte(string(msg.Body)), &culist)

		for _,v := range culist{
			pool.JobQueue <- SaveModelGen(v,QUEUE_DATA_CENTERUNLOAD)
		}
	}
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






