package pullcore

import (
	"github.com/streadway/amqp"
	"github.com/tsuna/gohbase"
	"github.com/zuston/AtcalMq/util"
	"github.com/ivpusic/grpool"
	"fmt"
	"github.com/zuston/AtcalMq/rabbitmq"
	"runtime"
)

const HANDLER_LOG_PATH  = "/tmp/AneHandler.log"
const HBASE_INI_PATH = "/opt/hbase.ini"


// singleton
type HbaseConn struct{
	Client gohbase.Client
	Name string
}

// 主存储和备份存储 hbase connection list
var HconnContainersMapper map[string]*HbaseConn

var hlogger *util.Logger

func InitHconn(){

	HconnContainersMapper := make(map[string]*HbaseConn)

	productionZkMapper,err := util.NewConfigReader(HBASE_INI_PATH,"production")
	util.CheckPanic(err)
	developmentZkMapper, err := util.NewConfigReader(HBASE_INI_PATH,"development")
	util.CheckPanic(err)

	proZk := productionZkMapper["zkquorum"]
	proHconn := &HbaseConn{}
	proHconn.Client = gohbase.NewClient(proZk)
	proHconn.Name = "production"
	HconnContainersMapper["production"] = proHconn

	if len(developmentZkMapper)!=0 {
		devZk := developmentZkMapper["zkquorum"]
		devConn := &HbaseConn{}
		devConn.Client = gohbase.NewClient(devZk)
		devConn.Name = "development"
		HconnContainersMapper["development"] = devConn
	}

	fmt.Println(HconnContainersMapper)
}

func init(){
	InitHconn()

	hlogger, _ = util.NewLogger(util.INFO_LEVEL,HANDLER_LOG_PATH)
	hlogger.SetDebug()

	// 协程池，设置运行池
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

		backuper.Info(string(msg.Body))
		list := ModelGen(msg.Body)
		for _,v := range list{
			// 统计埋点
			rabbitmq.StatisticsChan <- queue
			for _, hconn := range HconnContainersMapper{
				pool.JobQueue <- SaveModelGen(v,queue,hconn)
			}
		}
	}
}


// 调试 handler 示例代码
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
			pool.JobQueue <- SaveModelGen(v,queue,HconnContainersMapper["production"])

			for key,value := range v{
				hlogger.Info("%s=%s",key,string(value))
			}
			return
		}
		return
	}
}








