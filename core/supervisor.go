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
)

/**
get the rabbit mq info from the different original datasource
exposed to rpc
 */

var consumerContainer map[string]string

func Supervisor(queue *amqp.Queue){
	consumerContainer = make(map[string]string,100)
	for {
		consumerContainer[queue.Name] = strconv.Itoa(queue.Messages)
		fmt.Println(consumerContainer["ane_its_ai_data_centerLoad_queue"])
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
	zlloger.Info("RPC listen to port : 9898")
	http.Serve(lhandler,nil)
}