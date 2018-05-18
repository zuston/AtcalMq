package main

import (
	"github.com/zuston/AtcalMq/util"
	"github.com/zuston/AtcalMq/core"
	"github.com/zuston/AtcalMq/rabbitmq"
	"flag"
)

// declare the variable
var (
	mq_uri = ""
	exchange = ""
	exchange_type = ""
)

const (
	URL = "mq_uri"
	EXCHANGE = "exchange"
	EXCHANGE_TYPE = "exchange_type"
)

const MAIN_LOG_PATH  = "/tmp/AneMain.log"

var lloger *util.Logger

var (
	testTag = flag.Bool("test",false,"read the mq config.default is producing environment")
	showTag = flag.Bool("show",false,"just open the rpc, not consume the data")
	queueName = flag.String("queue","all","just set the test queue consume")
	debugTag = flag.Bool("debug",false,"set the log output level")
)

var (
	configPath = "./mq.cfg"
	configPathBackUp = "/opt/mq.cfg"
)

var (
	settingConsumeQueueTag = false
)

func init(){
	flag.Parse()
	if *showTag {
		for _,v := range core.BasicInfoTableNames{
			rabbitmq.AddSupervisorQueue(v)
		}
		rabbitmq.NewWatcher()
		return
	}

	// init the log instance
	lloger, _ = util.NewLogger(util.INFO_LEVEL,MAIN_LOG_PATH)
	// the log status is debug state
	lloger.SetDebug()

	// set the mq config by reading the mq config file
	lloger.Info("ready to read the mq config")
	// 测试mq账号
	if *testTag {
		configPath = "./tmq.cfg"
		configPathBackUp = "/opt/tmq.cfg"
	}
	configMapper,err := util.ConfigReader(configPath)
	if err!=nil {
		configMapper, _ = util.ConfigReader(configPathBackUp)
	}

	mq_uri = configMapper[URL]
	exchange = configMapper[EXCHANGE]
	exchange_type = configMapper[EXCHANGE_TYPE]

	lloger.Info("config uri=[%s], exchange=[%s], exchangeType=[%s]", mq_uri, exchange, exchange_type)

	if *queueName!="all" {
		settingConsumeQueueTag = true
	}
}


func main(){

	cf, err := rabbitmq.NewConsumerFactory(mq_uri,exchange,exchange_type)

	if err!=nil {
		panic("fail to connect to the message queue of rabbitmq")
	}

	if !settingConsumeQueueTag {
		cf.RegisterAll(core.BasicInfoTableNames,core.BasicHandler)
	}else {
		// 调试使用
		cf.Register(*queueName,core.TestHandler)
	}

	go cf.Handle()
	// provide the rpc service, expose the port 9898
	rabbitmq.NewWatcher()

	select {
	}
}

