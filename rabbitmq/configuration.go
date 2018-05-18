package rabbitmq

import "github.com/zuston/AtcalMq/util"

// 将一些 rabbitmq 下的log路径配置信息定义在此

const (
	// consumer 日志路径名
	RABBITMQ_CONSUMER_LOGGER_PATH = "/tmp/RbtConsumer.log"

	// producer
	RABBITMQ_PRODUCER_LOGGER_PATH = "/tmp/RbtProducer.log"

	// supervisor
	RABBITMQ_SUPERVISOR_LOGGER_PATH = "/tmp/RbtSuperviosr.log"

	// log level
	RABBITMQ_LOG_LEVEL = util.DEBUG_LEVEL

	// supervisor 中的配置连接文件路径
	ANE_CONFIG_PATH = "/opt/mq.cfg"
)
