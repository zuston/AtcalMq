package core

import (
	"github.com/streadway/amqp"
	"github.com/tsuna/gohbase"
	"github.com/zuston/AtcalMq/util"
	"github.com/zuston/AtcalMq/core/object"
	"encoding/json"
	"runtime"
	"github.com/ivpusic/grpool"
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

	hlogger, _ = util.NewLogger(util.INFO_LEVEL,"/tmp/handler.log")
	hlogger.SetDebug()

	// 协程池
	numCpus := runtime.NumCPU()
	runtime.GOMAXPROCS(numCpus)
}

func BasicHandler(queue string, msgChan <-chan amqp.Delivery){
	logPath := fmt.Sprintf("/tmp/backup/%d.backup",queue)
	backuper,_ := util.NewLogger(util.INFO_LEVEL,logPath)

	pool := grpool.NewPool(100, 100)
	defer pool.Release()

	for msg := range msgChan{
		msg.Ack(false)

		backuper.Info(string(msg.Body))
		list := ModelGen(msg.Body)
		for _,v := range list{
			pool.JobQueue <- SaveModelGen(v,queue)
		}
	}
}


/**
	 consume the mq msg
 */
func CenterLoadHandler(msgChan <-chan amqp.Delivery){

	backuper,_ := util.NewLogger(util.INFO_LEVEL,"/tmp/backup/centerload.backup")

	var clList []object.CenterLoadObj

	pool := grpool.NewPool(100, 100)
	defer pool.Release()

	for msg := range msgChan{
		msg.Ack(false)

		backuper.Info(string(msg.Body))

		json.Unmarshal([]byte(string(msg.Body)), &clList)

		for _,v := range clList{
			pool.JobQueue <- SaveModelGen(v,QUEUE_DATA_CENTERLOAD)
		}

	}
}




func BizOrderHandler(msgChan <-chan amqp.Delivery){
	select {
	case <-msgChan:
	}
}


func BizEwbHandler(msgChan <-chan amqp.Delivery){
	select {
	case <-msgChan:
	}
}

func BasicRouteHandler(msgChan <-chan amqp.Delivery){
	select {
	case <-msgChan:
	}
}

func BizEwbslistHandler(msgChan <-chan amqp.Delivery){
	backuper,_ := util.NewLogger(util.INFO_LEVEL,"/tmp/backup/ewbslist.backup")

	var culist []object.EwbsListObj

	pool := grpool.NewPool(100, 100)
	defer pool.Release()

	for msg := range msgChan{
		msg.Ack(false)
		backuper.Info(string(msg.Body))
		json.Unmarshal([]byte(string(msg.Body)), &culist)

		for _,v := range culist{
			pool.JobQueue <- SaveModelGen(v,QUEUE_BIZ_EWBSLIST)
		}
	}
}

func DataSiteLoadHandler(msgChan <-chan amqp.Delivery){
	select {
	case <-msgChan:
	}
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
	backuper,_ := util.NewLogger(util.INFO_LEVEL,"/tmp/backup/centerPallet.backup")

	var culist []object.CenterPalletObj

	pool := grpool.NewPool(100, 100)
	defer pool.Release()

	for msg := range msgChan{
		msg.Ack(false)
		backuper.Info(string(msg.Body))
		json.Unmarshal([]byte(string(msg.Body)), &culist)

		for _,v := range culist{
			pool.JobQueue <- SaveModelGen(v,QUEUE_DATA_CENTERPALLET)
		}
	}
}


func DataCenterSortHandler(msgChan <-chan amqp.Delivery){
	backuper,_ := util.NewLogger(util.INFO_LEVEL,"/tmp/backup/centerSort.backup")

	var culist []object.CenterSortObj

	pool := grpool.NewPool(100, 100)
	defer pool.Release()

	for msg := range msgChan{
		msg.Ack(false)
		backuper.Info(string(msg.Body))
		json.Unmarshal([]byte(string(msg.Body)), &culist)

		for _,v := range culist{
			pool.JobQueue <- SaveModelGen(v,QUEUE_DATA_CENTERSORT)
		}
	}
}


func BasicAreaHandler(msgChan <-chan amqp.Delivery){
	select {
	case <-msgChan:
	}
}

func BasicAttendHandler(msgChan <-chan amqp.Delivery){
	select {
	case <-msgChan:
	}
}


func TriggerSiteSendHandler(msgChan <-chan amqp.Delivery){
	select {
	case <-msgChan:
	}
}

func TriggerSiteUploadHandler(msgChan <-chan amqp.Delivery){
	select {
	case <-msgChan:
	}
}

func TriggerInOrOutHandler(msgChan <-chan amqp.Delivery){
	select {
	case <-msgChan:
	}
}


func TriggerStayOrLeaveHandler(msgChan <-chan amqp.Delivery){
	select {

	}
}

func TriggerCenterTransportHandler(msgChan <-chan amqp.Delivery){
	backuper,_ := util.NewLogger(util.INFO_LEVEL,"/tmp/backup/centerTransport.backup")

	var culist []object.CenterTransportObj

	pool := grpool.NewPool(100, 100)
	defer pool.Release()

	for msg := range msgChan{
		msg.Ack(false)
		backuper.Info(string(msg.Body))
		json.Unmarshal([]byte(string(msg.Body)), &culist)

		for _,v := range culist{
			pool.JobQueue <- SaveModelGen(v,QUEUE_TRIGGER_CENTERTRANSPORT)
		}
	}
}

func BasicSiteHandler(msgChan <-chan amqp.Delivery){
	select {

	}
}

func BasicVehicleLineHandler(msgChan <-chan amqp.Delivery){
	select {

	}
}


func BasicPlatFormHandler(msgChan <-chan amqp.Delivery){
	select {

	}
}






