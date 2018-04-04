package main

import (
	"github.com/zuston/AtcalMq/util"
	"github.com/zuston/AtcalMq/core"
	"github.com/zuston/AtcalMq/rabbitmq"
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

const (
	QUEUE_CENTERLOAD = "ane_its_ai_data_centerLoad_queue"
)




var lloger *util.Logger

func init(){
	// init the log instance
	lloger, _ = util.NewLogger(util.DEBUG_LEVEL, "/tmp/ane.log")
	// the log status is debug state
	lloger.SetDebug()

	// set the mq config by reading the mq config file
	lloger.Info("ready to read the mq config")
	configMapper,_ := util.ConfigReader("./mq.cfg")

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

	// centerload queue
	cf.Register(QUEUE_CENTERLOAD, core.CenterLoadHandler)

	go cf.Handle()

	// provide the rpc service, expose the port 9898
	core.NewWatcher()

	// send the msg to producer rabbitmq
	pf, err := rabbitmq.NewProducerFactory(mq_uri,exchange,exchange_type,false)
	// put the msg to channel
	pf.Publish("","")
	go pf.Handle()

	// main process blocking
	select {}
}