package rabbitmq

import (
	"github.com/zuston/AtcalMq/util"
	"github.com/streadway/amqp"
)

const (
	LOGGER_PATH = "/tmp/aneProducer.log"
)

type MsgEntity struct {
	queueName string
	msg string
}

type ProducerFactory struct {
	// log instance
	zlogger *util.Logger

	// setting the rabbitmq variables
	amqpUri string
	exchange string
	exchangeType string
	reliable bool

	// producer connection
	conn *amqp.Connection
	// setting channel
	channel *amqp.Channel

	// msgEntity channel
	msgChan chan MsgEntity
}

func NewProducerFactory(uri string, exchange string, exchangeType string, reliable bool) (*ProducerFactory, error) {

	zlloger, _ := util.NewLogger(util.DEBUG_LEVEL,LOGGER_PATH)

	cp := &ProducerFactory{
		zlogger:zlloger,

		amqpUri:uri,
		exchange:exchange,
		exchangeType:exchangeType,
		reliable:reliable,

		conn:nil,
		channel:nil,
		msgChan:make(chan MsgEntity, 1000),
	}

	zlloger.Info("ready to dial, amqp uri : [%s]",uri)
	var err error
	cp.conn, err = amqp.Dial(uri)
	if err != nil {
		zlloger.Error("dial error : %s",err)
		return nil, err
	}

	zlloger.Info("dial success")

	zlloger.Info("ready to get the channel")
	cp.channel, err = cp.conn.Channel()
	if err!=nil {
		zlloger.Error("acheive the channel error, err : %s",err)
		return nil, err
	}
	zlloger.Info("got the channel, decarling [%s] Exchange [%s]",exchangeType,exchange)

	if err := cp.channel.ExchangeDeclare(
		exchange,
		exchangeType,
		true,
		false,
		false,
		false,
		nil,
	);err!=nil {
		zlloger.Error("declaring the channel error, err : %s",err)
		return nil,err
	}

	zlloger.Info("declared the channel")
	if reliable {
		//todo
	}

	return cp, nil
}

func (cp *ProducerFactory) register(queueName string) bool{

	//cp.zlogger.Info("ready to init channel, channelName : [%s]", queueName)
	//channel, err := cp.conn.Channel()
	//if err!=nil {
	//	cp.zlogger.Error("acheive the channel error, err : %s",err)
	//	return false
	//}
	//cp.zlogger.Info("got [%s] Channel, decarling the [%s] Exchange [%s]", queueName,cp.exchangeType, cp.exchange)
	//if err := channel.ExchangeDeclare(
	//	cp.exchange,
	//	cp.exchangeType,
	//	true,
	//	false,
	//	false,
	//	false,
	//	nil,
	//);err!=nil {
	//	cp.zlogger.Error("declaring the channel error, err : %s",err)
	//	return false
	//}
	//
	//if cp.reliable {
	//	// todo
	//}
	//
	//
	return false
}

func (cp *ProducerFactory) Publish(queueName string, msg string){
	entity := MsgEntity{
		queueName:queueName,
		msg:msg,
	}
	cp.msgChan <- entity
}


func (cp *ProducerFactory) Handle(){
	cp.zlogger.Info("ready to publish the rabbitmq...")
	for msgEntity := range cp.msgChan{
		queueName := msgEntity.queueName
		msgBody := msgEntity.msg
		if err := cp.channel.Publish(
			cp.exchange,   // publish to an exchange
			queueName, // routing to 0 or more queues
			false,      // mandatory
			false,      // immediate
			amqp.Publishing{
				Headers:         amqp.Table{},
				ContentType:     "text/plain",
				ContentEncoding: "",
				Body:            []byte(msgBody),
				DeliveryMode:    amqp.Transient, // 1=non-persistent, 2=persistent
				Priority:        0,              // 0-9
				// a bunch of application/implementation-specific fields
			},
		); err != nil {
			cp.zlogger.Error("Exchange Publish Error, routingKey : %s, error : %s", queueName, err)
		}
	}
}





