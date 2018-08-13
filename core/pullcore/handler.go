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
// 数据备份路径点
//const BACKUPER_PATH = "/opt/aneBackup/"

// 更改数据备份点
const BACKUPER_PATH = "/extdata/aneBackup/"

const (
	POOL_WORK_NUMBER = 100
	POOL_QUEUE_LEN = 1000
)


// singleton
type HbaseConn struct{
	Client gohbase.Client
	Name string
}

// 主存储和备份存储 hbase connection list
var HconnContainersMapper map[string]*HbaseConn

var hlogger *util.Logger

// 单个传递,目前采用的方式
var singleConn *HbaseConn

func InitHconn(){

	HconnContainersMapper = make(map[string]*HbaseConn)

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

	singleConn = proHconn
}


// 实例化备份的 backuper
func instanceBackerUper(queueName string) *util.LLogger{
	backuper := util.NewLog4Go()
	backuper.SetType(1)
	backuper.SetRollingDaily(BACKUPER_PATH,queueName)
	return backuper
}

// 一些配置参数启动
func init(){
	InitHconn()

	hlogger, _ = util.NewLogger(util.INFO_LEVEL,HANDLER_LOG_PATH)
	hlogger.SetDebug()

	// 协程池，设置运行池
	numCpus := runtime.NumCPU()
	runtime.GOMAXPROCS(numCpus)
}

/**
最基础的数据存储函数
 */
func BasicHandler(queue string, msgChan <-chan amqp.Delivery){
	//logPath := fmt.Sprintf("/tmp/backup/%s.backup",queue)
	//oldBackuper,_ := util.NewLogger(util.INFO_LEVEL,logPath)
	newBackuper := instanceBackerUper(queue)

	pool := grpool.NewPool(100, 1000)
	defer pool.Release()

	hlogger.Info("[%s] enter the handler",queue)

	hlogger.Debug("HconnContainersMapper len is [%d]",len(HconnContainersMapper))

	for msg := range msgChan{
		msg.Ack(true)

		//oldBackuper.Info(string(msg.Body))
		newBackuper.Println(string(msg.Body))
		list := ModelGen(msg.Body)
		for _,v := range list{
			// 统计埋点
			rabbitmq.StatisticsChan <- queue
			for name, _ := range HconnContainersMapper{
				pool.JobQueue <- SaveModelGen(v,queue,name)
			}
		}
	}
}


/**
	调试 handler 示例代码
  */
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
			pool.JobQueue <- SaveModelGen(v,queue,"production")

		}
	}
}


/**
	新数据存储结构(二元组关系)存储模型
 */
 func MultiSavingHandler(qn string, msgChan <- chan amqp.Delivery){
	 pool := grpool.NewPool(POOL_WORK_NUMBER, POOL_QUEUE_LEN)
	 defer pool.Release()
	 hlogger.Info("processing the dual saving model handler...")

	 for msg := range msgChan{
	 	msg.Ack(true)
	 	infos := ModelGen(msg.Body)
	 	for _, info := range infos{
	 		pool.JobQueue <- MultiSavingModel(qn,info)
		}
	 }
 }





