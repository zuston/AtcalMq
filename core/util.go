package core

import (
	"strings"
	"fmt"
	"github.com/json-iterator/go"
	"github.com/tsuna/gohbase/hrpc"
	"context"
	"sync"
	"time"
	"math/rand"
)

const (
	CONSUMER_TAG_PREFIX = "CONSUMER-"
)

var lock *sync.Mutex
var randSeek = int64(1)


func init(){
	lock = &sync.Mutex{}
}
func GetQueueNameFromConsumerTag(s string) string {
	index := strings.Index(s,"-")
	return s[index+1:]
}
func SetConsumerTag(queueName string) string{
	return fmt.Sprintf("%s%s",CONSUMER_TAG_PREFIX,queueName)
}


func UidGen()string{
	r := rand.New(rand.NewSource(randomSource()))

	var f = func(min, max int64) int64{
		if min >= max || min == 0 || max == 0 {
			return max
		}
		return r.Int63n(max-min) + min
	}

	randomNumber := f(10000,99999)

	return fmt.Sprintf("%d%d",randomNumber,time.Now().UnixNano())
}

func randomSource()int64{
	lock.Lock()
	defer lock.Unlock()
	if randSeek >= 100000000 {
		randSeek = 1
	}
	randSeek++
	return time.Now().UnixNano()+randSeek
}

const(
	EWB_NO = "EwbNo"
	SITE_ID = "SiteId"
	VEHICLE_NO =	"VehicleNo"
	OPERATOR_CODE =	"OperatorCode"
	NEXT_SITE_ID = "NextSiteId"
)

// 生成存储model
func SaveModelGen(object interface{}, queue string) func(){
	
	// 人，站点，车，货
	// 站点有进出
	// operatorCode 以 ","为分割
	linkList := []string{
		EWB_NO,
		SITE_ID,
		VEHICLE_NO,
		OPERATOR_CODE,

		NEXT_SITE_ID,
	}

	// linkTableNameList
	linkTnList := map[string]string{
		EWB_NO:"Link_Ewb",
		SITE_ID:"Link_Site",
		VEHICLE_NO:"Link_Vehicle",
		OPERATOR_CODE:"Link_Operator",

		NEXT_SITE_ID : "Link_Site",
	}

	// 基础信息存储 model map
	basicInfoMapper := ModelMapperGen(object)
	// 基础信息--columnFamily
	basicInfoCf := "basic"
	// 基础信息--tableName
	basicInfoTn := queue
	// 基础信息--cf 和 model
	basicInfoCfMapper := map[string]map[string][]byte{basicInfoCf:basicInfoMapper}
	// 基础信息--uid----rowkey
	basicInfoUid := UidGen()

	return func() {
		// 基础信息--putrequest
		biPutRequest, err := hrpc.NewPutStr(context.Background(),basicInfoTn,basicInfoUid,basicInfoCfMapper)
		if err!=nil {
			hlogger.Error("[%s] bi build hrpc error : %s",queue,err)
			return
		}
		_, err = Hconn.Client.Put(biPutRequest)
		if err!=nil {
			hlogger.Error("[%s] bi hbase put error : %s",queue,err)
			return
		}

		// 关联信息
		for _,v := range linkList{
			if linkrowKey,ok := basicInfoMapper[v]; ok {
				linkrowKeys := []string{string(linkrowKey)}
				// 操作者有多个
				if v==OPERATOR_CODE {
					linkrowKeys = strings.Split(string(linkrowKey),",")
				}

				for _,lrk := range linkrowKeys{
					// linkInfoMapper
					liMapper := map[string][]byte{"uid":[]byte(basicInfoUid)}
					linkCfName := queue
					if v==NEXT_SITE_ID {
						linkCfName = fmt.Sprintf("%s_%s","nextSite",linkCfName)
					}
					liCfMapper := map[string]map[string][]byte{linkCfName:liMapper}
					linkPutReq, err := hrpc.NewPutStr(context.Background(),linkTnList[v],lrk,liCfMapper)
					if err!=nil {
						hlogger.Error("[%s] link [%s] build hrpc error : %s",queue,v,err)
						continue
					}
					_, err = Hconn.Client.Put(linkPutReq)
					if err!=nil {
						hlogger.Error("[%s] link [%s] hbase put error : %s",queue,v,err)
						continue
					}
				}


			}
		}
	}
}

func ModelMapperGen(object interface{}) map[string][]byte{
	marshal, _ := jsoniter.Marshal(object)
	mlen := len(string(marshal))
	jsonLine := string(marshal)[1:mlen-1]
	basicDataMapper := make(map[string][]byte)
	kvArr := strings.Split(string(jsonLine),`,"`)
	for i,v := range kvArr{
		index := strings.Index(v,":")
		// 去除前缀
		key := v[:index-1]
		if i==0 {
			key = v[1:index-1]
		}
		// 判断是否有前缀
		value := v[index+1:]
		if strings.Contains(value,`"`) {
			value = value[1:len(value)-1]
		}
		basicDataMapper[key] = []byte(value)
	}

	return basicDataMapper
}