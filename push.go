package main

import (
	"github.com/zuston/AtcalMq/util"
	"github.com/zuston/AtcalMq/rabbitmq"
	"time"
	"fmt"
	"encoding/json"
	"io/ioutil"
	"github.com/zuston/AtcalMq/core/object"
	"strings"
	"flag"
	"github.com/zuston/AtcalMq/core/pushcore"
	"github.com/zuston/AtcalMq/core"
	"github.com/streadway/amqp"
)


var logger *util.Logger
var mqUri string
var mqExchange string
var mqExchangeType string

const(
	MODEL_SUFFIX = "model"
)

const PUSH_LOG_PATH  = "/tmp/AnePush.log"

var(
	err error
	duration = flag.Duration("duration",time.Hour,"please enter the pushcore duration")
	optionPath = flag.String("option","/opt/optional.model","the pushcore queue column name")
	verifyTime = flag.String("verifyTime","","the pushcore time")
)

func init(){
	flag.Parse()

	logger,err = util.NewLogger(util.DEBUG_LEVEL,PUSH_LOG_PATH)
	util.CheckPanic(err)
	logger.SetDebug()
}

// 定时向对方推送计算好的数据
func main(){

	// 1. 初始化producer配置
	// 2. 读取抓取推送的配置文件
	// 3. 设定定时器，定时推送
	// 4. 推送信息

	configMapper, _ := util.NewConfigReader("/opt/mq.ini","rabbitmq")
	mqUri = configMapper["mq_uri"]
	mqExchange = configMapper["exchange"]
	mqExchangeType = configMapper["exchange_type"]
	pf, _ := rabbitmq.NewProducerFactory(mqUri,mqExchange,mqExchangeType,true)

	//*duration = 60*time.Second
	//*optionPath = "./optional.model"
	//*verifyTime = "2017-05-20 23:50:31"
	timelyPublish(pf,*duration,*optionPath,*verifyTime)
}


func timelyPublish(pf *rabbitmq.ProducerFactory, duration time.Duration, optionPath string, verifyTime string) {
	// 解析映射文件
	translationT := ParseTranslationTable(optionPath)
	// 为了防止新建队列无法推送消息，先进行消费
	declaredQueues := getDeclaredQueues(translationT)
	preventAction(declaredQueues)
	ticker := time.NewTicker(duration)

	for _ = range ticker.C{
		pusher(pf,translationT,verifyTime)
	}
}

// 防御性编程
func preventAction(queueList []string) {
	cf, err := rabbitmq.NewConsumerFactory(mqUri,mqExchange,mqExchangeType,false)
	util.CheckPanic(err)
	cf.RegisterAll(queueList, func(queue string, msgChan <-chan amqp.Delivery) {
	})
	logger.Info("%s has declared",queueList)
	cf.CloseAll()
}

func getDeclaredQueues(i map[string]map[string]string) []string {
	var list []string
	for key, _ := range i{
		_, queueName := parseComponentKey(key)
		list = append(list,queueName)
	}
	return list
}

// 获取数据，推送
func pusher(producerFactory *rabbitmq.ProducerFactory, translationT map[string]map[string]string, verifytime string) {
	for componentKey, convertMapper := range translationT{
		dbName, queueName := parseComponentKey(componentKey)
		pushJsonList := pushcore.Get(core.MONGO_TAG,dbName,convertMapper,verifytime)
		notifyInfo := fmt.Sprintf("[%s] - [%s] 当次推送消息数量：%d",queueName,dbName,len(pushJsonList))
		logger.Info(notifyInfo)
		//pushLine := "["+strings.Join(pushJsonList,",")+"]"
		for _, info := range pushJsonList{
			producerFactory.Publish(queueName,info)
		}
		util.BarkNotify(notifyInfo)
	}
}

func parseComponentKey(componentKey string) (string, string) {
	arr := strings.Split(componentKey,"*")
	return arr[0],arr[1]
}

// 形式为 ：{"tablename":{"name":"zname","age":"zage"}}
func ParseTranslationTable(modelPath string) map[string]map[string]string {
	// 兼容传递目录或者直接是文件
	files := []string{modelPath}
	if util.IsDir(modelPath) {
		files = util.WalkDir(modelPath,MODEL_SUFFIX)
	}
	modelMappers := make(map[string]map[string]string)

	for _,file := range files{
		data, err := ioutil.ReadFile(file)
		checkPanic(err)
		datajson := []byte(data)

		var models []object.PushMapperObj
		err = json.Unmarshal(datajson, &models)
		checkPanic(err)


		for _,model := range models{
			tablename := model.TABLENAME
			queuename := model.QUEUENAME

			componentKey := fmt.Sprintf("%s*%s",tablename,queuename)
			modelMapper := make(map[string]string)
			for _,relation := range model.REALTIONS{
				modelMapper[relation.CN] = relation.QN
			}

			modelMappers[componentKey] = modelMapper
		}
	}
	return modelMappers
}

func checkPanic(err error) {
	if err!=nil {
		logger.Error("读取model映射文件报错")
		panic(err)
	}	
}

