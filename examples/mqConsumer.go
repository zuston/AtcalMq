package main

import (
	"flag"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"github.com/zuston/AtcalMq/util"
)

var (
	exchange     = flag.String("exchange", "ane_its_exchange", "Durable, non-auto-deleted AMQP exchange name")
	exchangeType = flag.String("exchange-type", "direct", "Exchange type - direct|fanout|topic|x-custom")
	consumerTag  = flag.String("consumer-tag", "simple-consumer", "AMQP consumer tag (should not be blank)")
)

var uri string
var queue string
var bindingKey string

func init() {
	flag.Parse()
	configMapper,_ := util.ConfigReader("./mq.cfg")

	uri = configMapper["mq_uri"]
	queue = "ane_its_ai_biz_order_queue"
	bindingKey = queue
}

func main() {
	_, err := NewConsumer(uri, *exchange, *exchangeType, queue, bindingKey, *consumerTag)
	if err != nil {
		log.Fatalf("%s", err)
	}


	//if err := c.Shutdown(); err != nil {
	//	log.Fatalf("error during shutdown: %s", err)
	//}
	select {

	}
}

type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	tag     string
	done    chan error
}

func NewConsumer(amqpURI, exchange, exchangeType, queueName, key, ctag string) (*Consumer, error) {
	c := &Consumer{
		conn:    nil,
		channel: nil,
		tag:     ctag,
		done:    make(chan error),
	}

	var err error

	log.Printf("dialing %q", amqpURI)
	c.conn, err = amqp.Dial(amqpURI)
	if err != nil {
		return nil, fmt.Errorf("Dial: %s", err)
	}



	log.Printf("got Connection, getting Channel")
	c.channel, err = c.conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("Channel: %s", err)
	}

	log.Printf("got Channel, declaring Exchange (%q)", exchange)
	if err = c.channel.ExchangeDeclare(
		exchange,     // name of the exchange
		exchangeType, // type
		true,         // durable
		false,        // delete when complete
		false,        // internal
		false,        // noWait
		nil,          // arguments
	); err != nil {
		return nil, fmt.Errorf("Exchange Declare: %s", err)
	}

	log.Printf("declared Exchange, declaring Queue %q", queueName)
	queue, err := c.channel.QueueDeclare(
		queueName, // name of the queue
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // noWait
		nil,       // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("Queue Declare: %s", err)
	}

	log.Printf("declared Queue (%q %d messages, %d consumers), binding to Exchange (key %q)",
		queue.Name, queue.Messages, queue.Consumers, key)

	if err = c.channel.QueueBind(
		queue.Name, // name of the queue
		key,        // bindingKey
		exchange,   // sourceExchange
		false,      // noWait
		nil,        // arguments
	); err != nil {
		return nil, fmt.Errorf("Queue Bind: %s", err)
	}

	//err = c.channel.Qos(
	//	10,
	//	0,
	//	false,
	//)

	log.Printf("Queue bound to Exchange, starting Consume (consumer tag %q)", c.tag)
	deliveries, err := c.channel.Consume(
		queue.Name, // name
		c.tag,      // consumerTag,
		false,      // noAck
		false,      // exclusive
		false,      // noLocal
		false,      // noWait
		nil,        // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("Queue Consume: %s", err)
	}

	go handle(deliveries, c.done)

	go func() {
		fmt.Printf("conn closing: %s", <-c.conn.NotifyClose(make(chan *amqp.Error)))
		fmt.Printf("channel closing: %s", <-c.channel.NotifyClose(make(chan *amqp.Error)))
	}()


	return c, nil
}

func (c *Consumer) Shutdown() error {
	// will close() the deliveries channel
	if err := c.channel.Cancel(c.tag, true); err != nil {
		return fmt.Errorf("Consumer cancel failed: %s", err)
	}

	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("AMQP connection close error: %s", err)
	}

	defer log.Printf("AMQP shutdown OK")

	// wait for handle() to exit
	return <-c.done
}

func (c *Consumer) Supervisor() {
	select {
	case <-c.done:
		c.channel.Close()
		c.conn.Close()

	}
}

func handle(deliveries <-chan amqp.Delivery, done chan error) {
	i := 1
	for d := range deliveries {
		if i==2 {
			break
		}
		log.Printf("index : %d",i)
		fmt.Println(string(d.Body))
		d.Ack(false)
		//time.Sleep(10*time.Second)
		i++
	}
	log.Printf("handle: deliveries auto closed")
	done <- nil
}
