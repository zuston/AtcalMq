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
	QUEUE_DATA_CENTERLOAD = "ane_its_ai_data_centerLoad_queue"

	QUEUE_BIZ_ORDER = "ane_its_ai_biz_order_queue"

	QUEUE_BIZ_EWB = "ane_its_ai_biz_ewb_queue"

	QUEUE_BASIC_ROUTE = "ane_its_ai_basic_route_queue"

	QUEUE_BIZ_EWBSLIST = "ane_its_ai_biz_ewbsList_queue"

	QUEUE_DATA_SITELOAD = "ane_its_ai_data_siteLoad_queue"

	QUEUE_DATA_CENTERUNLOAD = "ane_its_ai_data_centerUnload_queue"

	QUEUE_DATA_CENTERPALLET = "ane_its_ai_data_center Pallet_queue"

	QUEUE_DATA_CENTERSORT = "ane_its_ai_data_centerSort_queue"

	QUEUE_BASIC_AREA = "ane_its_ai_basic_area_queue"

	QUEUE_BASIC_ATTEND = "ane_its_ai_basic_attend_queue"

	QUEUE_TRIGGER_SITESEND = "ane_its_ai_trigger_siteSend_queue"

	QUEUE_TRIGGER_SITEUPLOAD = "ane_its_ai_trigger_siteUpload_queue"

	QUEUE_TRIGGER_INOROUT = "ane_its_ai_trigger_inOrOut_queue"

	QUEUE_TRIGGER_STAYORLEAVE = "ane_its_ai_trigger_stayOrLeave_queue"

	QUEUE_TRIGGER_CENTERTRANSPORT = "ane_its_ai_trigger_centerTransport_queue"

	QUEUE_BASIC_SITE = "ane_its_ai_basic_site_queue"

	QUEUE_BASIC_VEHICLELINE = "ane_its_ai_basic_vehicleLine_queue"

	QUEUE_BASIC_PLATFORM = "ane_its_ai_basic_platform_queue"
)

var lloger *util.Logger

func init(){
	// init the log instance
	lloger, _ = util.NewLogger(util.DEBUG_LEVEL, "/tmp/ane.log")
	// the log status is debug state
	lloger.SetDebug()

	// set the mq config by reading the mq config file
	lloger.Info("ready to read the mq config")
	configMapper,err := util.ConfigReader("./mq.cfg")
	if err!=nil {
		configMapper, _ = util.ConfigReader("/opt/mq.cfg")
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

	// 2.5.5 centerload queue
	cf.Register(QUEUE_DATA_CENTERLOAD, core.CenterLoadHandler)


	// BIZ ORDER 2.1
	cf.Register(QUEUE_BIZ_ORDER, core.BizOrderHandler)

	// BIZ EWB  2.2
	cf.Register(QUEUE_BIZ_EWB, core.BizEwbHandler)

	// BASIC ROUTE 2.3
	cf.Register(QUEUE_BASIC_ROUTE, core.BasicRouteHandler)

	// BIZ EWBSLIST 2.4
	cf.Register(QUEUE_BIZ_EWBSLIST, core.BizEwbslistHandler)

	// DATA SITELOAD 2.5.1
	cf.Register(QUEUE_DATA_SITELOAD, core.DataSiteLoadHandler)

	// 2.5.2
	cf.Register(QUEUE_DATA_CENTERUNLOAD, core.DataCenterUnloadHandler)

	// 2.5.3
	cf.Register(QUEUE_DATA_CENTERPALLET, core.DataCenterPalletHandler)

	// 2.5.4
	cf.Register(QUEUE_DATA_CENTERSORT, core.DataCenterSortHandler)


	// 2.6
	cf.Register(QUEUE_BASIC_AREA, core.BasicAreaHandler)

	//2.7
	cf.Register(QUEUE_BASIC_ATTEND, core.BasicAttendHandler)

	//2.8.1
	cf.Register(QUEUE_TRIGGER_SITESEND, core.TriggerSiteSendHandler)

	//2.8.2
	cf.Register(QUEUE_TRIGGER_SITEUPLOAD, core.TriggerSiteUploadHandler)

	//2.8.3
	cf.Register(QUEUE_TRIGGER_INOROUT, core.TriggerInOrOutHandler)

	//2.8.4
	cf.Register(QUEUE_TRIGGER_STAYORLEAVE, core.TriggerStayOrLeaveHandler)

	// 2.8.5
	cf.Register(QUEUE_TRIGGER_CENTERTRANSPORT, core.TriggerCenterTransportHandler)

	//2.9
	cf.Register(QUEUE_BASIC_SITE, core.BasicSiteHandler)

	//2.10
	cf.Register(QUEUE_BASIC_PLATFORM, core.BasicPlatFormHandler)

	//2.11
	cf.Register(QUEUE_BASIC_VEHICLELINE, core.BasicVehicleLineHandler)


	go cf.Handle()
	// provide the rpc service, expose the port 9898
	core.NewWatcher()

	// main process blocking
	select {}
}

// producer demo
func initProducer(){
	// send the msg to producer rabbitmq
	pf, _ := rabbitmq.NewProducerFactory(mq_uri,exchange,exchange_type,false)
	// put the msg to channel
	pf.Publish("","")
	go pf.Handle()
}