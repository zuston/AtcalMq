package main

import (
	"github.com/zuston/AtcalMq/util"
	"github.com/zuston/AtcalMq/rabbitmq"
	"github.com/streadway/amqp"
	"fmt"
	"github.com/ivpusic/grpool"
	"strings"
	"github.com/tsuna/gohbase"
	"github.com/tsuna/gohbase/hrpc"
	"context"
	"log"
	"time"
	"github.com/zuston/AtcalMq/core/pullcore"
	"github.com/zuston/AtcalMq/core"
)

var verificationChan chan []byte

var handlerCount int

const(
	// 验证的数目
	N = 10
)

var linkList = []string{
	pullcore.EWB_NO,
	pullcore.SITE_ID,
	pullcore.VEHICLE_NO,
	pullcore.OPERATOR_CODE,
	pullcore.NEXT_SITE_ID,
}

var linkTnList = map[string]string{
	pullcore.EWB_NO:"Link_Ewb",
	pullcore.SITE_ID:"Link_Site",
	pullcore.VEHICLE_NO:"Link_Vehicle",
	pullcore.OPERATOR_CODE:"Link_Operator",
	pullcore.NEXT_SITE_ID : "Link_Site",
}

// 验证数据从mq读取到到是否准确存储进db中,防止并发下出现问题
//func TestModelSaver(t *testing.T){
func main(){

	verificationChan = make(chan []byte,N)

	testQueueName := "ane_its_ai_data_centerLoad_queue"

	hclient := gohbase.NewClient("slave4,slave2,slave3")

	// 启动消费
	configMapper, _ := util.NewConfigReader("/opt/mq.cfg","rabbitmq")
	mq_uri := configMapper["mq_uri"]
	exchange := configMapper["exchange"]
	exchange_type := configMapper["exchange_type"]

	cf, err := rabbitmq.NewConsumerFactory(mq_uri,exchange,exchange_type,false)

	if err!=nil {
		panic("fail to connect to the message queue of rabbitmq")
	}

	cf.Register(testQueueName,handler)

	go cf.Handle()

	count := 1

	time.Sleep(60*time.Second)

	// 验证数据
	for i:=0;i<N;i++{
		msg := <-verificationChan
		list := pullcore.ModelGen(msg)
		for _, v := range list{
			log.Println(v)
			verifyContainer := make(map[string]string,5)
			// 查看是否有关联的字段,且将关联字段放进map中
			for _,lname := range linkList{
				linkName := strings.ToLower(lname)
				linkRowKey, ok := v[linkName]
				linkRowKey = pullcore.FixRowKey(linkRowKey)
				if ok {
					verifyContainer[lname] = linkRowKey
				}
			}

			log.Println(verifyContainer)

			infoOneKey := v[strings.ToLower(core.LinkKey[testQueueName])]

			randomNumber := pullcore.CacheMapper[v["hewbno"]]

			infoOneKey = fmt.Sprintf("%s_%d",infoOneKey,randomNumber)

			var uidVerifyContainer []string

			slog := "tableName:[%s],rowkey:[%s],cfName:[%s],key:[%s],error:[%s]"

			// 看基础信息是否放进对应的表中
			for linkName,linkRowKey := range verifyContainer{
				tableName := linkTnList[linkName]
				cfName := testQueueName
				if linkName==pullcore.NEXT_SITE_ID {
					cfName = fmt.Sprintf("nextSite_%s",testQueueName)
				}
				rowkeys := []string{linkRowKey}

				operatorSplitTag := linkName==pullcore.OPERATOR_CODE
				if  operatorSplitTag {
					rowkeys = strings.Split(linkRowKey,",")
				}

				for _,rowkey := range rowkeys{
					cf := map[string][]string{cfName:[]string{infoOneKey}}
					getRequest, _ := hrpc.NewGetStr(context.Background(), tableName, rowkey,
						hrpc.Families(cf))
					getRsp, err := hclient.Get(getRequest)
					if err!=nil {
						log.Println(fmt.Sprintf(slog,tableName,rowkey,cfName,infoOneKey,err))
					}
					if len(getRsp.Cells)!=1 {
						log.Println(fmt.Sprintf(slog,tableName,rowkey,cfName,infoOneKey,"overflow 1"))
					}
					for _,cellV := range getRsp.Cells{
						uid := string(cellV.Value)
						log.Printf("[%s] uid:[%s]",linkName,uid)
						uidVerifyContainer = append(uidVerifyContainer,uid)
					}
				}

			}
			if modelCheck(verifyContainer,uidVerifyContainer) && uidCheck(uidVerifyContainer) {
				log.Println(fmt.Sprintf("[%d] 验证通过",count))
			}else{
				log.Println(fmt.Sprintf("[%d] 验证失败",count))
			}

			count++
		}

	}

	log.Printf("接收到[%d],验证了[%d]",handlerCount,count)
}


// 检测存储的表和直接获取到的是否一致
func modelCheck(i map[string]string, i2 []string) bool{
	//fmt.Println(len(i),len(i2))
	op,ok := i[pullcore.OPERATOR_CODE]
	if ok {
		count := len(strings.Split(op,","))
		fmt.Println(len(i)+count-1)
		return len(i2)==(len(i)+count-1)
	}
	return len(i2)==len(i)
}

// 看几个表中取出的uid是否全部相同
func uidCheck(i []string) bool{
	checkV := i[0]
	for _,v := range i{
		if checkV!=v {
			return false
		}
	}
	return true
}

func handler(queue string, msgChan <-chan amqp.Delivery){
	handlerCount = 0
	logPath := fmt.Sprintf("/tmp/backup/%s.backup",queue)
	backuper,_ := util.NewLogger(util.DEBUG_LEVEL,logPath)

	pool := grpool.NewPool(10, 100)
	defer pool.Release()

	count := 0
	for msg := range msgChan{
		if count==N {
			return
		}
		msg.Ack(true)

		backuper.Info(string(msg.Body))
		// 进入验证通道
		verificationChan <- msg.Body
		list := pullcore.ModelGen(msg.Body)
		for _,v := range list{
			handlerCount++
			pool.JobQueue <- pullcore.SaveModelGen(v,queue)
		}
	}
}
