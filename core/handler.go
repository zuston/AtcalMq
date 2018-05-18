package core

import (
	"github.com/streadway/amqp"
	"github.com/tsuna/gohbase"
	"github.com/zuston/AtcalMq/util"
	"runtime"
	"github.com/ivpusic/grpool"
	"fmt"
	"github.com/zuston/AtcalMq/rabbitmq"
)

const HANDLER_LOG_PATH  = "/tmp/AneHandler.log"

// singleton
type HbaseConn struct{
	Client gohbase.Client
}

var Hconn *HbaseConn

var hlogger *util.Logger

func init(){
	Hconn = &HbaseConn{}
	Hconn.Client = gohbase.NewClient("slave4,slave2,slave3")

	hlogger, _ = util.NewLogger(util.INFO_LEVEL,HANDLER_LOG_PATH)
	hlogger.SetDebug()

	// 协程池
	numCpus := runtime.NumCPU()
	runtime.GOMAXPROCS(numCpus)
}

func BasicHandler(queue string, msgChan <-chan amqp.Delivery){
	logPath := fmt.Sprintf("/tmp/backup/%s.backup",queue)
	backuper,_ := util.NewLogger(util.INFO_LEVEL,logPath)

	pool := grpool.NewPool(100, 1000)
	defer pool.Release()

	hlogger.Info("[%s] enter the handler",queue)

	for msg := range msgChan{
		msg.Ack(true)
		hlogger.Debug("+++[%s] first enter",queue)

		backuper.Info(string(msg.Body))
		list := ModelGen(msg.Body)
		for _,v := range list{
			// 统计埋点
			rabbitmq.StatisticsChan <- queue
			pool.JobQueue <- SaveModelGen(v,queue)
		}
	}
}

func TestHandler(queue string, msgChan <-chan amqp.Delivery){
	logPath := fmt.Sprintf("/tmp/backup/%s.backup",queue)
	backuper,_ := util.NewLogger(util.DEBUG_LEVEL,logPath)

	pool := grpool.NewPool(10, 100)
	defer pool.Release()

	hlogger.Info("[%s] enter the handler",queue)

	for msg := range msgChan{
		msg.Ack(false)
		hlogger.Debug("+++[%s] first enter",queue)

		backuper.Info(string(msg.Body))
		list := ModelGen(msg.Body)
		for _,v := range list{
			hlogger.Debug("+++[%s] handle one",queue)
			// 统计埋点
			rabbitmq.StatisticsChan <- queue
			pool.JobQueue <- SaveModelGen(v,queue)

			for key,value := range v{
				hlogger.Info("%s=%s",key,string(value))
			}
			return
		}
		return
	}
}








