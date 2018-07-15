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
	"os/exec"
	"bytes"
	"sync"
	"strings"
)

/**
get the rabbit mq info from the different original datasource
exposed to rpc
 */

type PullSupervisorObj struct {
	// 队列名
	QueueName string
	// 积压数
	Overstock string
	// 单位时间处理能力
	UnitHandlerAblitity int
	// 24小时内处理数据量
	TotalHandlerCount	int64
}

type PushSupervisorObj struct {
	QueueName string
	Overstock string
	UnitHandlerAblitity int
}

// rpc回传 struct
type RpcObj struct {
	// consumer 队列状态表
	PullSupervisorObjs []PullSupervisorObj
	// producer 队列状态表
	PushSupervisorObjs []PushSupervisorObj
	// 守护程序启动时间
	StartTime string
	// 守护程序运行时间
	Runtime time.Duration
	// 消费队列数目
	PullQueueNumber int
	// push队列数目
	PushQueueNumber int
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

var (
	// 放置队列名
	StatisticsChan chan string
	// 统计容器
	statisticsContainer map[string]int64

	lock *sync.Mutex
)

var (
	// 启动时间
	startTime time.Time
)

var zlloger *util.Logger

func init(){
	// 初始化队列数组，便于curl获取
	//queueNameContainer = make([]string,1)

	startTime = time.Now()

	// 从上一级读取
	// todo 最好从上一级传入，解耦
	configMapper,_ := util.NewConfigReader(ANE_CONFIG_PATH,"console")
	apiUsername = configMapper["username"]
	apiPassword = configMapper["password"]
	apiAddress = configMapper["apiaddress"]

	zlloger, _ = util.NewLogger(RABBITMQ_LOG_LEVEL,RABBITMQ_SUPERVISOR_LOGGER_PATH)

	zlloger.SetDebug()

	// 初始化mapper
	queueNameCountMapper = make(map[string]string,1)
	queueUnitHandlerMapper = make(map[string]int,1)

	statisticsContainer = make(map[string]int64,20)
	lock = new(sync.Mutex)
	StatisticsChan = make(chan string,1000)
	return
}

// 提供给 consumer 进行注册
func AddSupervisorQueue(queueName string){
	queueNameContainer = append(queueNameContainer,queueName)
}

func Statistics(){

	// 定时器
	minTicker := time.NewTicker(time.Second*30*60)
	dayTicker := time.NewTicker(time.Hour*24)

	go func() {
		for _ = range minTicker.C{
			lock.Lock()
			defer lock.Unlock()
			statisticNotify()
		}
	}()

	go func() {
		for _ = range dayTicker.C{
			statisticNotify()
			statisticsContainer = make(map[string]int64,20)
		}
	}()

	// 计数
	for queueName := range StatisticsChan{
		if v, ok := statisticsContainer[queueName];ok{
			statisticsContainer[queueName] = v+1
		}else {
			statisticsContainer[queueName] = 1
		}
	}

}
func statisticNotify() {
	pushList := make([]string,len(statisticsContainer))
	for key,value := range statisticsContainer{
		pushList = append(pushList, key+" : "+fmt.Sprint(value))
	}
	util.BarkNotify(strings.Join(pushList,"\n"))
	util.WechatNotify(strings.Join(pushList,"\n"))
}

func getPullQueueNumber() int{
	return len(queueNameContainer)
}

// post 调用
func rabbitStatus(url string)(string,error){
	client := &http.Client{}
	apiUrl := url
	req, err := http.NewRequest("GET", apiUrl, nil)
	req.SetBasicAuth(apiUsername, apiPassword)
	resp, err := client.Do(req)
	if err != nil{
		return "",err
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	s := string(bodyText)
	return s,nil
}

// shell 调用
func execShell(commandLine string) (string,error){
	cmd := exec.Command("/bin/bash","-c",commandLine)
	var out bytes.Buffer

	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return out.String(),err
	}
	return out.String(),nil
}

func Supervisor(){
	// 分钟获取
	//timer := time.NewTicker(time.Minute)
	// 获取当前队列积压数据量
	// 循环获取
	for {
		execCommandLine := fmt.Sprintf(`curl  -u %s:%s %s`,apiUsername,apiPassword,apiAddress)
		jsonResp, err := execShell(execCommandLine)

		if err!=nil {
			errorFormat := fmt.Sprintf("supervisor timely got error : [%s]",err)
			zlloger.Error(errorFormat)
			util.BarkNotify(errorFormat)
			util.WechatNotify(errorFormat)
			continue
		}
		countList := gjson.Get(jsonResp,"#.messages")
		nameList := gjson.Get(jsonResp,"#.name")
		for i,v := range countList.Array(){
			currentV := v
			name := nameList.Array()[i].String()
			originalV := queueNameCountMapper[name]
			queueNameCountMapper[name] = currentV.String()

			originalCountInt, _ := strconv.Atoi(originalV)
			countVInt, _ := strconv.Atoi(currentV.String())
			queueUnitHandlerMapper[name] = originalCountInt-countVInt
		}
		time.Sleep(time.Minute)
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
	b, err := jsonGen()
	if err!=nil {
		return err
	}
	*result = string(b)
	return nil
}

func jsonGen() (string, error){
	var rpcObj = RpcObj{}

	// pullcore queue status
	var supervisorEntities []PullSupervisorObj
	for _,v := range queueNameContainer{
		tempSupervisorObj := PullSupervisorObj{
			Overstock:queueNameCountMapper[v],
			QueueName:v,
			UnitHandlerAblitity:queueUnitHandlerMapper[v],
			TotalHandlerCount:statisticsContainer[v],
		}
		supervisorEntities = append(supervisorEntities, tempSupervisorObj)
	}

	// 求注册队列与获取的队列的差集
	pushNames := differentSet(queueNameContainer,queueNameCountMapper)
	var pushEntities []PushSupervisorObj
	for _,pushName := range pushNames{
		tempPush := PushSupervisorObj{
			QueueName:pushName,
			Overstock:queueNameCountMapper[pushName],
			UnitHandlerAblitity:queueUnitHandlerMapper[pushName],
		}
		pushEntities = append(pushEntities,tempPush)
	}

	rpcObj.PullSupervisorObjs = supervisorEntities
	rpcObj.PullQueueNumber = getPullQueueNumber()
	rpcObj.PushQueueNumber = len(pushNames)
	rpcObj.PushSupervisorObjs = pushEntities
	rpcObj.Runtime = time.Since(startTime)
	rpcObj.StartTime = startTime.Format("2006-01-02 15:04:05")

	b, err := jsoniter.Marshal(rpcObj)
	if err!=nil {
		zlloger.Error("json encode error : %s",err)
		return "",err
	}
	return string(b),nil
}

// 差集
func differentSet(i []string, i2 map[string]string) []string {
	var list []string
	sameMap := make(map[string]bool,10)
	for _,v := range i{
		if _,ok := i2[v];ok {
			sameMap[v] = false
		}
	}

	for key,_ := range i2{
		if _,ok := sameMap[key];!ok{
			list = append(list,key)
		}
	}
	return list
}

func RpcService(){
	watcher := new(Watcher)
	rpc.Register(watcher)
	rpc.HandleHTTP()
	lhandler, err := net.Listen("tcp",RABBITMQ_SUPERVISOR_RPC_URL)
	if err != nil{
		zlloger.Error("listen to the URL error : %s",err)
		return
	}
	zlloger.Info("Supervisor RPC listen to "+RABBITMQ_SUPERVISOR_RPC_URL)
	http.Serve(lhandler,nil)
}


func NewWatcher(){
	// 定时获取数据
	go Supervisor()

	// 统计结果发送 wechat
	go Statistics()

	// rpc service expose
	go RpcService()
}

