package main

import (
	"github.com/streadway/amqp"
	"fmt"
	"github.com/zuston/AtcalMq/util"
)

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
	QUEUE_1 = "ane_its_ai_data_centerLoad_queue"
)

type Consumer map[string]interface{}

type ConsumerChan <-chan amqp.Delivery

// consumers factory
type ConsumerFactory struct {
	conn *amqp.Connection
	// 注册消费队列
	registerChan chan Consumer
	// 处理报错信息
	done chan error
	exchange string
	exchangeType string
}

func NewConsumerFactory(url string, exchange string, exchangeType string) (*ConsumerFactory,error){
	cf := &ConsumerFactory{
		conn:nil,
		registerChan:make(chan Consumer,1000),
		done:make(chan error,1000),
		exchange:exchange,
		exchangeType:exchangeType,
	}
	// start the dial
	lloger.Debug("ready to dial, dail url : %s",url)
	var err error
	cf.conn, err = amqp.Dial(url)
	if err!=nil {
		return nil,fmt.Errorf("DIAL ERROR : ",err)
	}
	// dial success
	lloger.Debug("dial success")

	return cf, nil
}

func (cf *ConsumerFactory) register(queueName string, f func(msgChan <-chan amqp.Delivery)) error {
	lloger.Info("%s register to factory",queueName)
	tempMap := make(Consumer,1)
	tempMap[queueName] = f
	cf.registerChan <- tempMap
	return nil
}

// close all the queue
func (cf *ConsumerFactory) closeAll(){

}

// close the declared Consumer
func (cf *ConsumerFactory) close(queueName string){

}


func (cf *ConsumerFactory) handle() {
	// 1. handle register queue consumer
	// 2. do the consume msg of consumer's action
	for consumer := range cf.registerChan {
		for k,v := range consumer {
			//fmt.Println(k)
			//v.(func(msg string))("sd")
			queueName := k
			handleFunc := v

			// start the acheive the channel
			lloger.Debug("ready to acheive the channel")
			channel, err := cf.conn.Channel()
			if err!=nil {
				lloger.Error("!acheive the channel error, queueName : %s, error : %s",queueName, err)
				continue
			}
			// channel success
			lloger.Debug("acheive the channel success")

			// declare exchange
			lloger.Debug("ready to declare the exchange, exchange name : %s", exchange)
			if err = channel.ExchangeDeclare(
				cf.exchange,
				cf.exchangeType,
				true,
				false,
				false,
				false,
				nil,
			); err!=nil {
				lloger.Error("!declare the exchange error, queueName : %s, error : %s",queueName, err)
				continue
			}
			// declare exchange success
			lloger.Debug("declare the exchange success")

			// declare queue
			lloger.Debug("ready to declare the queue, queue name : %s", queueName)
			queue, err := channel.QueueDeclare(
				queueName,
				true,
				false,
				false,
				false,
				nil,
			)
			if err!=nil {
				lloger.Error("!declare the queue error, queueName : %s, error : %s",queueName, err)
				continue
			}
			// queue declare success
			lloger.Debug("declare the queue success")


			// set the channel bind to the declared queue
			lloger.Debug("ready to bind the declared queue")
			if err := channel.QueueBind(
				queue.Name,
				queueName,
				cf.exchange,
				false,
				nil,
			); err!=nil {
				lloger.Error("!bind the queue error, queueName : %s, error : %s",queueName, err)
				continue
			}
			// channel bind success
			lloger.Debug("bind the declared queue success")

			// start the consume
			lloger.Debug("%s ready to consume the queue",queueName)
			deliveries, err := channel.Consume(
				queue.Name,
				queue.Name,
				false,
				false,
				false,
				false,
				nil,
			)

			if err!=nil {
				lloger.Error("%s fail to consume the queue", queueName)
				continue
			}

			//log.Println(queue.Name,queue.Consumers,queue.Messages)
			go handleFunc.(func(msgChan <-chan amqp.Delivery))(deliveries)
		}
	}
}


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
	cf, err := NewConsumerFactory(mq_uri,exchange,exchange_type)

	if err!=nil {
		panic("fail to connect to the message queue of rabbitmq")
	}

	cf.register(QUEUE_1, func(msgChan <-chan amqp.Delivery) {
		i := 1
		for msg := range msgChan{
			if i==2 {
				return
			}
			fmt.Println("accept the info : ",i)
			i++
			msg.Ack(false)
		}
	})

	cf.handle()
}