package rabbitmq

import (
	"github.com/streadway/amqp"
	"fmt"
	"github.com/zuston/AtcalMq/util"
	"github.com/zuston/AtcalMq/core"
)

const (
	CONSUMER_LOGGER_PATH = "/tmp/aneConsumer.log"
)

type Consumer map[string]interface{}

type ConsumerChan <-chan amqp.Delivery

// consumers factory
type ConsumerFactory struct {
	// zlog
	zloger *util.Logger
	conn *amqp.Connection
	// 注册消费队列
	registerChan chan Consumer
	// 处理报错信息
	done chan error
	exchange string
	exchangeType string
}

func NewConsumerFactory(url string, exchange string, exchangeType string) (*ConsumerFactory,error){
	zlogger, _ := util.NewLogger(util.DEBUG_LEVEL,CONSUMER_LOGGER_PATH)
	zlogger.SetDebug()

	cf := &ConsumerFactory{
		zloger:zlogger,
		conn:nil,
		registerChan:make(chan Consumer,1000),
		done:make(chan error,1000),
		exchange:exchange,
		exchangeType:exchangeType,
	}
	// start the dial
	zlogger.Debug("ready to dial, dail url : %s",url)
	var err error
	cf.conn, err = amqp.Dial(url)
	if err!=nil {
		return nil,fmt.Errorf("DIAL ERROR : ",err)
	}
	// dial success
	zlogger.Debug("dial success")

	return cf, nil
}

func (cf *ConsumerFactory) Register(queueName string, f func(msgChan <-chan amqp.Delivery)) error {
	// add the register queueName to Supervisor component
	core.AddSupervisorQueue(queueName)
	// set the special queue color
	cf.zloger.Info("%c[1;40;32m%s%c[0m register to factory",0x1B,queueName,0x1B)
	tempMap := make(Consumer,1)
	tempMap[queueName] = f
	cf.registerChan <- tempMap
	return nil
}

// close all the queue
func (cf *ConsumerFactory) CloseAll(){
	// todo
}

// close the declared Consumer
func (cf *ConsumerFactory) Close(queueName string){
	// todo
}


func (cf *ConsumerFactory) Handle() {
	// 1. handle register queue consumer
	// 2. do the consume msg of consumer's action
	for consumer := range cf.registerChan {
		for k,v := range consumer {
			//fmt.Println(k)
			//v.(func(msg string))("sd")
			queueName := k
			handleFunc := v

			// start the acheive the channel
			cf.zloger.Debug("ready to acheive the channel")
			channel, err := cf.conn.Channel()
			if err!=nil {
				cf.zloger.Error("!acheive the channel error, queueName : %s, error : %s",queueName, err)
				continue
			}
			// channel success
			cf.zloger.Debug("acheive the channel success")

			// declare exchange
			cf.zloger.Debug("ready to declare the exchange, exchange name : %s", cf.exchange)
			if err = channel.ExchangeDeclare(
				cf.exchange,
				cf.exchangeType,
				true,
				false,
				false,
				false,
				nil,
			); err!=nil {
				cf.zloger.Error("!declare the exchange error, queueName : %s, error : %s",queueName, err)
				continue
			}
			// declare exchange success
			cf.zloger.Debug("declare the exchange success")

			// declare queue
			cf.zloger.Debug("ready to declare the queue, queue name : %s", queueName)
			queue, err := channel.QueueDeclare(
				queueName,
				true,
				false,
				false,
				false,
				nil,
			)
			if err!=nil {
				cf.zloger.Error("!declare the queue error, queueName : %s, error : %s",queueName, err)
				continue
			}
			// queue declare success
			cf.zloger.Debug("declare the queue success")


			// set the channel bind to the declared queue
			cf.zloger.Debug("ready to bind the declared queue")
			if err := channel.QueueBind(
				queue.Name,
				queueName,
				cf.exchange,
				false,
				nil,
			); err!=nil {
				cf.zloger.Error("!bind the queue error, queueName : %s, error : %s",queueName, err)
				continue
			}
			// channel bind success
			cf.zloger.Debug("bind the declared queue success")

			// 设置预读机制，防止处理过慢，导致connection reset
			err = channel.Qos(
				10,
				0,
				false,
			)

			// start the consume
			cf.zloger.Debug("%s ready to consume the queue",queueName)
			deliveries, err := channel.Consume(
				queue.Name,
				fmt.Sprintf("counsumer-%s",queueName),
				false,
				false,
				false,
				false,
				nil,
			)

			if err!=nil {
				cf.zloger.Error("%s fail to consume the queue", queueName)
				continue
			}

			//log.Println(queue.Name,queue.Consumers,queue.Messages)
			// provider info to supervisor
			// remove it because of getting the data from api, like
			// curl -i -u sitrab:sitrab123456 http://58.215.167.31:15672/api/queues/its-test/ane_its_ai_data_centerLoad_queue
			//go core.Supervisor(channel,queueName)

			go handleFunc.(func(msgChan <-chan amqp.Delivery))(deliveries)
		}
	}
}
