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
	"io"
	"sync"
	"strings"
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

var (
	// 放置队列名
	StatisticsChan chan string
	// 统计容器
	statisticsContainer map[string]int64

	lock *sync.Mutex
)

var zlloger *util.Logger

func init(){
	// 初始化队列数组，便于curl获取
	//queueNameContainer = make([]string,1)

	configMapper,_ := util.ConfigReader(ANE_CONFIG_PATH)
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

func statistics(){

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
	util.WechatNotify(strings.Join(pushList,"\n"))
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

func AddSupervisorQueue(queueName string){
	queueNameContainer = append(queueNameContainer,queueName)
}

func Supervisor(){
	// 获取当前队列积压数据量
	// 循环获取
	for {
		execCommandLine := fmt.Sprintf(`curl  -u %s:%s %s`,apiUsername,apiPassword,apiAddress)
		jsonResp, err := execShell(execCommandLine)

		//jsonResp, err := rabbitStatus(apiAddress)

		if err!=nil {
			errorFormat := fmt.Sprintf("supervisor timely got error : [%s]",err)
			zlloger.Error(errorFormat)
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
	b, err := jsonGen()
	if err!=nil {
		return err
	}
	*result = string(b)
	return nil
}

func jsonGen() (string, error){
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
		return "",err
	}
	return string(b),nil
}

/**
http service
 */
// hello world, the web server
func GetAll(w http.ResponseWriter, req *http.Request) {
	json, _ := jsonGen()
	io.WriteString(w, json)
}

func initHttp(){
	http.HandleFunc("/getall", GetAll)
	err := http.ListenAndServe(":8909", nil)
	if err != nil {
		zlloger.Error("ListenAndServe: ", err)
	}
}


func NewWatcher(){
	// 循环刷新数据
	go Supervisor()

	// http service
	go initHttp()

	// 统计
	go statistics()

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

