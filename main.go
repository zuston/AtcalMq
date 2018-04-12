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

var lloger *util.Logger

var (
	testTag = flag.Bool("test",false,"read the mq config.default is producing environment")
	showTag = flag.Bool("show",false,"just open the rpc, not consume the data")
)

var (
	configPath = "./mq.cfg"
	configPathBackUp = "/opt/mq.cfg"
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
	lloger, _ = util.NewLogger(util.DEBUG_LEVEL, "/tmp/ane.log")
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
}


func main(){
	// show the queue msg info, through the ConsoleQueueChan
	// todo: ...

	cf, err := rabbitmq.NewConsumerFactory(mq_uri,exchange,exchange_type)

	if err!=nil {
		panic("fail to connect to the message queue of rabbitmq")
	}

	cf.RegisterAll(core.BasicInfoTableNames,core.BasicHandler)

	/*

	// 2.5.5 centerload queue
	cf.Register(core.QUEUE_DATA_CENTERLOAD, core.BasicHandler)

	// 2.5.2
	cf.Register(core.QUEUE_DATA_CENTERUNLOAD, core.BasicHandler)

	// 2.5.4
	cf.Register(core.QUEUE_DATA_CENTERSORT, core.BasicHandler)
	// 2.5.3
	cf.Register(core.QUEUE_DATA_CENTERPALLET, core.BasicHandler)

	// 2.8.5
	cf.Register(core.QUEUE_TRIGGER_CENTERTRANSPORT, core.BasicHandler)

	// BIZ EWBSLIST 2.4
	cf.Register(core.QUEUE_BIZ_EWBSLIST, core.BasicHandler)

	/**
	// BIZ ORDER 2.1
	cf.Register(core.QUEUE_BIZ_ORDER, core.BizOrderHandler)

	// BIZ EWB  2.2
	cf.Register(core.QUEUE_BIZ_EWB, core.BizEwbHandler)

	// BASIC ROUTE 2.3
	cf.Register(core.QUEUE_BASIC_ROUTE, core.BasicRouteHandler)


	// DATA SITELOAD 2.5.1
	cf.Register(core.QUEUE_DATA_SITELOAD, core.DataSiteLoadHandler)

	// 2.6
	cf.Register(core.QUEUE_BASIC_AREA, core.BasicAreaHandler)

	//2.7
	cf.Register(core.QUEUE_BASIC_ATTEND, core.BasicAttendHandler)

	//2.8.1
	cf.Register(core.QUEUE_TRIGGER_SITESEND, core.TriggerSiteSendHandler)

	//2.8.2
	cf.Register(core.QUEUE_TRIGGER_SITEUPLOAD, core.TriggerSiteUploadHandler)

	//2.8.3
	cf.Register(core.QUEUE_TRIGGER_INOROUT, core.TriggerInOrOutHandler)

	//2.8.4
	cf.Register(core.QUEUE_TRIGGER_STAYORLEAVE, core.TriggerStayOrLeaveHandler)

	//2.9
	cf.Register(core.QUEUE_BASIC_SITE, core.BasicSiteHandler)

	//2.10
	cf.Register(core.QUEUE_BASIC_PLATFORM, core.BasicPlatFormHandler)

	//2.11
	cf.Register(core.QUEUE_BASIC_VEHICLELINE, core.BasicVehicleLineHandler)

	*/

	go cf.Handle()
	// provide the rpc service, expose the port 9898
	rabbitmq.NewWatcher()

	select {
	}
}

// producer demo
func initProducer(){
	// send the msg to producer rabbitmq
	pf, _ := rabbitmq.NewProducerFactory(mq_uri,exchange,exchange_type,false)
	// put the msg to channel
	pf.Publish("","")
	go pf.Handle()
}