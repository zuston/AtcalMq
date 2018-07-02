package rabbitmq

import (
	"github.com/streadway/amqp"
	"fmt"
	"github.com/zuston/AtcalMq/util"
	"time"
)

const (
	CONSUMER_TAG_PREFIX = "CONSUMER-"
)

type Consumer struct{
	queueName string
	handlerFunc interface{}
}

type ConsumerChan <-chan amqp.Delivery

// consumers factory
type ConsumerFactory struct {
	// zlog
	zloger *util.Logger
	// 总的一个链接 connection,（已废弃）
	conn *amqp.Connection
	// 注册消费队列
	registerChan chan *Consumer
	// 处理关闭信息
	done bool
	exchange string
	exchangeType string
	amqpUrl string
	// 监视器，防止处理队列挂掉，重新再连接
	restartChan chan string
	// registerMapper
	registerMapper map[string]interface{}
	// register rabbitmq connection and channel struct mapper
	registerConnMapper map[string]*RabbitConn
}

type RabbitConn struct {
	conn *amqp.Connection
	channel *amqp.Channel
}

func NewConsumerFactory(url string, exchange string, exchangeType string, isSupervisor bool) (*ConsumerFactory,error){
	zlogger, _ := util.NewLogger(RABBITMQ_LOG_LEVEL,RABBITMQ_CONSUMER_LOGGER_PATH)
	zlogger.SetDebug()

	cf := &ConsumerFactory{
		zloger:zlogger,
		conn:nil,
		registerChan:make(chan *Consumer,20),
		done:false,
		exchange:exchange,
		exchangeType:exchangeType,
		amqpUrl:url,
		restartChan:make(chan string,20),
		registerMapper:make(map[string]interface{},20),
		registerConnMapper:make(map[string]*RabbitConn,10),
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

	// start the handle function
	go cf.Handle()

	// supervisor goroutinue
	if isSupervisor {
		go NewWatcher()
	}

	return cf, nil
}

func (cf *ConsumerFactory) Register(queueName string, f interface{}) error {
	// add the register queueName to Supervisor component
	AddSupervisorQueue(queueName)
	cf.zloger.Info("%c[1;40;32m%s%c[0m register to factory",0x1B,queueName,0x1B)

	consumerReg := &Consumer{
		queueName:queueName,
		handlerFunc:f,
	}

	cf.registerChan <- consumerReg
	cf.registerMapper[queueName] = f
	return nil
}

func (cf *ConsumerFactory) RegisterAll(queueList []string,f interface{}) error{
	for _,v := range queueList{
		cf.Register(v,f)
	}
	return nil
}

// close all the queue
func (cf *ConsumerFactory) CloseAll(){
	// todo
	cf.done = true
	close(cf.restartChan)
	close(cf.registerChan)
	cf.conn.Close()
}

// close the declared Consumer
func (cf *ConsumerFactory) Close(queueName string){
	// todo
}


func (cf *ConsumerFactory) Handle() {

	// 监听异步执行的任务,挂掉则重启
	go func(){
		for stopQueueName := range cf.restartChan{
			//cf.registerConnMapper[stopQueueName].channel.Cancel(setConsumerTag(stopQueueName),false)
			//cf.registerConnMapper[stopQueueName].conn.Close()

			//cf.zloger.Error("[%s] stop",stopQueueName)

			time.Sleep(10*time.Second)
			f := cf.registerMapper[stopQueueName]

			consumerReg := &Consumer{
				queueName:stopQueueName,
				handlerFunc:f,
			}
			// 重启任务
			cf.zloger.Info("[%s] restart the channel",stopQueueName)
			util.BarkNotify(fmt.Sprintf("[%s] restart the channel",stopQueueName))
			cf.registerChan <- consumerReg

		}
	}()

	// 1. 处理注册queue
	// 2. 异步consume数据
	for consumer := range cf.registerChan {
		queueName := consumer.queueName
		handleFunc := consumer.handlerFunc

		cf.zloger.Debug("ready to dial...")
		amqpConn, err := amqp.Dial(cf.amqpUrl)
		if err!=nil {
			cf.zloger.Error("[%s] dial amqp Connection error : %s",queueName,err)
			util.BarkNotify(fmt.Sprintf("[%s] dial amqp Connection error : %s",queueName,err))
			continue
		}

		// start the acheive the channel
		cf.zloger.Debug("ready to acheive the channel")
		channel, err := amqpConn.Channel()
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
			5,
			0,
			false,
		)

		// start the consume
		cf.zloger.Info("%c[1;40;32m%s%c[0m ready to consume the queue",0x1B,queueName,0x1B)
		deliveries, err := channel.Consume(
			queue.Name,
			setConsumerTag(queueName),
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

		cf.registerConnMapper[queueName] = &RabbitConn{
			conn:amqpConn,
			channel:channel,
		}

		go func(){
			// 处理数据业务流程
			handleFunc.(func(name string,msgc <- chan amqp.Delivery))(queueName,deliveries)
			// 结束通知,便于重启
			// 判断是否关闭通道
			if !cf.done {
				cf.restartChan <- queueName
			}
		}()
	}
}

func setConsumerTag(queueName string) string{
	return fmt.Sprintf("%s%s",CONSUMER_TAG_PREFIX,queueName)
}
