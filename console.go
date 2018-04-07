package main


import (
	"net/rpc"
	"github.com/zuston/AtcalMq/util"
	"time"
)


const (
	QUEUE_1 = "ane_its_ai_data_centerLoad_queue"
)

var zllogger *util.Logger
func init(){
	zllogger, _ = util.NewLogger(util.DEBUG_LEVEL,"/tmp/console.log")
	zllogger.SetDebug()
}


func main() {

	client, err := rpc.DialHTTP("tcp","127.0.0.1:9898")

	if err!=nil {
		zllogger.Error("connect to the rpc server error : %s",err)
		return
	}
	zllogger.Info("connect to the rpc server success")
	var reply string

	for{
		err = client.Call("Watcher.GetOverStock",QUEUE_1, &reply)
		if err!=nil {
			zllogger.Error("call function error : ",err)
			continue
		}
		zllogger.Info("queueName : %s, overstock number : %s",QUEUE_1,reply)
		time.Sleep(10*time.Second)
	}
	return
}