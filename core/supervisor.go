package core

import (
	"github.com/streadway/amqp"
	"time"
	"strconv"
	"net/rpc"
	"net"
	"fmt"
	"net/http"
	"github.com/zuston/AtcalMq/util"
	"log"
)

/**
get the rabbit mq info from the different original datasource
exposed to rpc
 */

var consumerContainer map[string]string

func Supervisor(channel *amqp.Channel, queueName string){
	consumerContainer = make(map[string]string,100)
	for {
		queue, err := channel.QueueDeclare(
			queueName,
			true,
			false,
			false,
			false,
			nil,
		)
		if err!=nil {
			log.Printf("%s",err)
		}
		fmt.Println(queue.Messages)
		consumerContainer[queue.Name] = strconv.Itoa(queue.Messages)
		time.Sleep(10*time.Second)
	}
}

type Watcher int

func (w *Watcher) GetInfo(key string, result *string) error {
	*result = consumerContainer[key]
	return nil
}

func NewWatcher(){

	zlloger, _ := util.NewLogger(util.DEBUG_LEVEL,"/tmp/supervisor.log")
	zlloger.SetDebug()

	watcher := new(Watcher)
	rpc.Register(watcher)
	rpc.HandleHTTP()
	lhandler, err := net.Listen("tcp",":9898")
	if err != nil{
		zlloger.Error("listen to the port error : %s",err)
		return
	}
	zlloger.Info("Console RPC listen to port : 9898")
	http.Serve(lhandler,nil)
}