package rabbitmq

import (
	"github.com/zuston/AtcalMq/util"
	"fmt"
)

// 将一些 rabbitmq 下的log路径配置信息定义在此

var (
	// rabbitmq 日志路径
	RABBIT_MQ_LOG_PATH_PREFIX = "/extdata/log"

	// consumer 日志路径名
	RABBITMQ_CONSUMER_LOGGER_PATH = fmt.Sprintf("%s/RbtConsumer.log",RABBIT_MQ_LOG_PATH_PREFIX)

	// producer
	RABBITMQ_PRODUCER_LOGGER_PATH = fmt.Sprintf("%s/RbtProducer.log",RABBIT_MQ_LOG_PATH_PREFIX)

	// supervisor
	RABBITMQ_SUPERVISOR_LOGGER_PATH = fmt.Sprintf("%s/RbtSuperviosr.log",RABBIT_MQ_LOG_PATH_PREFIX)

	// log level
	RABBITMQ_LOG_LEVEL = util.INFO_LEVEL

	// supervisor 中的配置连接文件路径
	ANE_CONFIG_PATH = "/opt/mq.ini"

	// supervisor rpc port
	RABBITMQ_SUPERVISOR_RPC_URL = "127.0.0.1:9898"
)
