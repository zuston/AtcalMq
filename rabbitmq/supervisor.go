package rabbitmq

import (
	"net/rpc"
	"net"
	"net/http"
	"github.com/zuston/AtcalMq/util"
	"io/ioutil"
	"time"
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/json-iterator/go"
	"strconv"
)

/**
get the rabbit mq info from the different original datasource
exposed to rpc
 */

type SupervisorObj struct {
	QueueName string
	Overstock string
	UnitHandlerAblitity int
}

var (
	apiUsername string
	apiPassword string
	apiAddress string
)

var (
// 当前队列名下的挤压数据量
 queueNameCountMapper map[string]string
// 当前队列名下的单位处理数据量
 queueUnitHandlerMapper map[string]int
 queueNameContainer []string
)

var zlloger *util.Logger

func init(){
	// 初始化队列数组，便于curl获取
	//queueNameContainer = make([]string,1)

	configMapper,_ := util.ConfigReader("/opt/mq.cfg")
	apiUsername = configMapper["username"]
	apiPassword = configMapper["password"]
	apiAddress = configMapper["apiaddress"]

	zlloger, _ = util.NewLogger(util.DEBUG_LEVEL,"/tmp/supervisor.log")

	zlloger.SetDebug()

	// 初始化mapper
	queueNameCountMapper = make(map[string]string,1)
	queueUnitHandlerMapper = make(map[string]int,1)
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

func AddSupervisorQueue(queueName string){
	queueNameContainer = append(queueNameContainer,queueName)
}

func Supervisor(){

	// 获取当前队列积压数据量
	// 循环获取
	for {
		for _,v := range queueNameContainer{
			queuename := v
			jsonResp,err := rabbitCliApi(queuename)
			if err!=nil {
				errorFormat := fmt.Sprintf("supervisor timely got [%s] error : [%s]",queuename,err)
				zlloger.Error(errorFormat)
				util.WechatNotify(errorFormat)
				continue
			}
			countV := gjson.Get(jsonResp,"messages")
			originalCount := queueNameCountMapper[queuename]
			queueNameCountMapper[queuename] = countV.String()

			// 设置unitHandlerAbility
			originalCountInt, _ := strconv.Atoi(originalCount)
			countVInt, _ := strconv.Atoi(countV.String())
			queueUnitHandlerMapper[queuename] = originalCountInt-countVInt

		}
		time.Sleep(60*time.Second)
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

func (w *Watcher) GetUnitHandlerAbility(queueName string, result *string) error{
	*result = strconv.Itoa(queueUnitHandlerMapper[queueName])
	return nil
}

// 获取所有队列情况，返回json
func (w *Watcher) GetAll(tag string, result *string) error{
	var supervisorEntities []SupervisorObj
	for _,v := range queueNameContainer{
		tempSupervisorObj := SupervisorObj{
			Overstock:queueNameCountMapper[v],
			QueueName:v,
			UnitHandlerAblitity:queueUnitHandlerMapper[v],
		}
		supervisorEntities = append(supervisorEntities, tempSupervisorObj)
	}
	b, err := jsoniter.Marshal(supervisorEntities)
	if err!=nil {
		zlloger.Error("json encode error : %s",err)
		return err
	}
	*result = string(b)
	return nil
}

func NewWatcher(){
	// wechat handler queue
	go util.NotifyHandlerQueue()
	// 循环刷新数据
	go Supervisor()
	// rpc 暴露接口
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

