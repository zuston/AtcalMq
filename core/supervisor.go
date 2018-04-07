package core

import (
	"net/rpc"
	"net"
	"net/http"
	"github.com/zuston/AtcalMq/util"
	"io/ioutil"
	"time"
	"fmt"
	"github.com/tidwall/gjson"
)

/**
get the rabbit mq info from the different original datasource
exposed to rpc
 */

const (
	QUEUE_CENTERLOAD = "ane_its_ai_data_centerLoad_queue"
)

var apiUsername string
var apiPassword string
var apiAddress string

var queueNameCountMapper map[string]string

var queueNameArr []string

var zlloger *util.Logger

func init(){
	// 初始化队列数组，便于curl获取
	queueNameArr = make([]string,2)
	queueNameArr = append(queueNameArr,QUEUE_CENTERLOAD)

	configMapper,_ := util.ConfigReader("/opt/mq.cfg")
	apiUsername = configMapper["username"]
	apiPassword = configMapper["password"]
	apiAddress = configMapper["apiaddress"]

	zlloger, _ = util.NewLogger(util.DEBUG_LEVEL,"/tmp/supervisor.log")

	zlloger.SetDebug()

	// 初始化mapper
	queueNameCountMapper = make(map[string]string,1)
}


func rabbitCliApi(queueName string) (string,error) {
	client := &http.Client{}
	apiUrl := fmt.Sprintf("%s/%s",apiAddress,queueName)
	req, err := http.NewRequest("GET", apiUrl, nil)
	req.SetBasicAuth(apiUsername, apiPassword)
	resp, err := client.Do(req)
	if err != nil{
		zlloger.Error("rabbitCli api request error : %s",err)
		return "",err
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	s := string(bodyText)
	return s,nil
}

func Supervisor(){
	// 获取当前队列积压数据量
	// 循环获取
	for {
		for _,v := range queueNameArr{
			queuename := v
			jsonResp,err := rabbitCliApi(queuename)
			if err!=nil {
				zlloger.Error("supervisor timely got [%s] error : [%s]",queuename,err)
				continue
			}
			countV := gjson.Get(jsonResp,"backing_queue_status.len")
			queueNameCountMapper[queuename] = countV.String()
		}
		time.Sleep(10*time.Second)
	}
}

/**
rpc function
 */
type Watcher int

func (w *Watcher) GetOverStock(queueName string, result *string) error {
	*result = queueNameCountMapper[queueName]
	return nil
}

func NewWatcher(){
	go Supervisor()
	watcher := new(Watcher)
	rpc.Register(watcher)
	rpc.HandleHTTP()
	lhandler, err := net.Listen("tcp",":9898")
	if err != nil{
		zlloger.Error("listen to the port error : %s",err)
		return
	}
	zlloger.Info("Supervisor RPC listen to port : 9898")
	http.Serve(lhandler,nil)
}

