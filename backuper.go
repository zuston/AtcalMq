package main

import (
	"log"
	"flag"
	"strings"
	"github.com/zuston/AtcalMq/util"
	"sync"
	"os"
	"bufio"
	"fmt"
	"time"
	"github.com/zuston/AtcalMq/rabbitmq"
	"github.com/streadway/amqp"
	"github.com/zuston/AtcalMq/core"
)

/**
将日志文件打入mq中程序
 */

// curl -u ane:ane1106 http://202.120.117.179:15672/api/queues/its-pds

const BACKUPER_DEFAULT_PATH = "/opt/aneBackup/"
//const BACKUPER_DEFAULT_PATH = "/temp/"

// 备份通道
var infoChannel chan backuperStruct

const TIME_FORMAT = "2006-01-02"

const (
	MQ_SECTION = "amqp"
	MQ_URL = "mq_uri"
	MQ_EXCHANGE = "exchange"
	MQ_TYPE = "exchange_type"
)


var (
	backuperPath = flag.String("path",BACKUPER_DEFAULT_PATH,"choose the backuper path")
	// 日期之间以 "," 来进行分割
	filterSuffix = flag.String("filter","","choose the filter backuper path suffix")
	// mq 配置文件位置
	mqConfig = flag.String("config","/opt/mq.ini","enter the rabbitmq backuper config file file position")
	// 推送至二元关系存储模型 还是 备份恢复队列
	savingTag = flag.Bool("dual",true,"set the dual relation saving model")
)

var barrier sync.WaitGroup

var backupPusher *rabbitmq.ProducerFactory

var done chan bool

var haveDone chan bool


func init(){
	flag.Parse()
	infoChannel = make(chan backuperStruct,1000)
	done = make(chan bool,1)
	haveDone = make(chan bool,1)


	mqConfigs, err := util.NewConfigReader(*mqConfig,MQ_SECTION)
	backupPusher, err = rabbitmq.NewProducerFactory(mqConfigs[MQ_URL],mqConfigs[MQ_EXCHANGE],mqConfigs[MQ_TYPE],false)
	util.CheckPanic(err)

	if *savingTag  && len(*filterSuffix)==0 {
		*filterSuffix = generateFilterQueue()
	}

}

func generateFilterQueue() string {
	ymd := time.Now().AddDate(0,0,-1).Format(TIME_FORMAT)
	var newlist []string
	for _, value := range core.MultiRelationSavingTableNames{
		index := strings.Index(value,"_")
		newlist = append(newlist,fmt.Sprintf("%s.%s.%s",value[index+1:],"log",ymd))
	}
	return strings.Join(newlist,",")
}

func main(){
	//pullMQ()
	//return
	//pushMQ()
	//return
	//testInc()
	//
	//return
	outputNotify()

	filterFiles := filter(*backuperPath,strings.Split(*filterSuffix,","))
	log.Printf("current filter list is : [%s]",filterFiles)

	go channelHandler()
	goroutineReader(filterFiles)
	<- haveDone
	log.Println("have finish reading and pushing the data!")
	// 为了防止 rpc 请求过多，导致很多信息丢失
	time.Sleep(time.Minute*3)
}

// 此队列可以用于备份恢复 或者是 多元关系模型存储
func outputNotify() {
	if *savingTag {
		log.Println("current backup mode : [dual relation model saving]")
	}else {
		log.Println("current backup mode : [recover the backup data....]")
	}
}


func goroutineReader(files []string) {
	for i:=0;i<len(files);i++{
		barrier.Add(1)
	}
	for _, filePath := range files{
		go func() {
			queueName,YMD := parse(filePath)
			if *savingTag {
				queueName = fmt.Sprintf("%s_%s",core.SAVING_RELATION_PREFIX,queueName)
			}
			startTime := time.Now()
			_ = YMD
			file, err := os.Open(filePath)
			// 恢复机制，记录点
			util.CheckPanic(err)

			defer func() {
				file.Close()
				barrier.Done()
			}()

			reader := bufio.NewReader(file)
			for {
				line, err := readLine(reader)
				if err!=nil {
					break
				}
				//infoTag := line[0:22]
				info := line[22:]
				//fmt.Println(infoTag,info)

				infoChannel <- backuperStruct{
					// todo 是二元关系存储还是备份的。更改queue的前缀
					queueName:queueName,
					info:info,
				}
			}
			endTime := time.Now()
			fmt.Printf("[%s] reading cost : [%s]\n",queueName,endTime.Sub(startTime))
		}()
	}
	log.Println("waiting reading the files")
	barrier.Wait()
	done <- false
}

func channelHandler(){
	log.Println("start consuming the data....")
	//for infoStruct := range infoChannel{
	//	queueName := infoStruct.queueName
	//	info := infoStruct.info
	//	backupPusher.Publish(queueName,info)
	//}

	for {
		select {
		case infoStruct := <- infoChannel :
			queueName := infoStruct.queueName
			info := infoStruct.info
			backupPusher.Publish(queueName,info)
		case <- done :
			haveDone <- true
			return
		}
	}
}

/**
功能性函数
 */
func parse(filePath string) (string, string) {
	removePrefixPathIndex := strings.LastIndex(filePath,"/")
	filePath = filePath[removePrefixPathIndex+1:]
	arrs := strings.Split(filePath,".")
	return arrs[0],arrs[2]
}


func filter(path string, filterSuffixs []string) []string {
	var files []string
	for _,filterSuffix := range filterSuffixs{
		f := util.WalkDir(path,filterSuffix)
		for _,key := range f{
			files = append(files,key)
		}
	}
	return files
}

type backuperStruct struct {
	queueName string
	info string
}

func readLine(r *bufio.Reader) (string, error) {
	line, isprefix, err := r.ReadLine()
	for isprefix && err == nil {
		var bs []byte
		bs, isprefix, err = r.ReadLine()
		line = append(line, bs...)
	}
	return string(line), err
}


/**
	测试程序
  */
func testInc(){
	files := filter("/temp",[]string{"2018-06-19"})
	log.Println(files)
}

func pushMQ(){
	pf, _ := rabbitmq.NewProducerFactory("amqp://ane:ane1106@202.120.117.179:5672/its-pds","test_exchange","direct",false)
	for i:=1;i<1000;i++{
		pf.Publish("test_queue","hello world")
	}
	select {

	}
}

func pullMQ(){
	consumer, _ := rabbitmq.NewConsumerFactory("amqp://ane:ane1106@202.120.117.179:5672/its-pds","test_exchange","direct",false)
	consumer.Register("test_queue", func(qn string, msgChan <- chan amqp.Delivery) {
		count := 0
		for msg := range msgChan{
			msg.Ack(true)
			log.Println(string(msg.Body))
			count++
		}
		log.Println(count)
	})
	select {

	}
}


