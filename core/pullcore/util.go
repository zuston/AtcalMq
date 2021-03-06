package pullcore

import (
	"strings"
	"fmt"
	"github.com/tsuna/gohbase/hrpc"
	"context"
	"sync"
	"time"
	"math/rand"
	"encoding/json"
	"go/types"
	"github.com/zuston/AtcalMq/util"
	"github.com/zuston/AtcalMq/core"
)

const ERROR_INSERT_PATH = "/extdata/aneError/"

var errorLogger *util.LLogger

var lock *sync.Mutex
var randSeek = int64(1)

var clientLock *sync.Mutex

// 用于测试的cache
var CacheMapper map[string]int64

func init(){
	lock = &sync.Mutex{}
	clientLock = &sync.Mutex{}

	CacheMapper = make(map[string]int64,100)

	errorLogger = util.NewLog4Go()
	errorLogger.SetType(1)
	errorLogger.SetRollingDaily(ERROR_INSERT_PATH,"errorline")
}

func UidGen() (string,int64){
	r := rand.New(rand.NewSource(randomSource()))

	var f = func(min, max int64) int64{
		if min >= max || min == 0 || max == 0 {
			return max
		}
		return r.Int63n(max-min) + min
	}

	randomNumber := f(100000,999999)


	return fmt.Sprintf("%s_%d",time.Now().Format("2006-01-02 15:04:05"),randomNumber),randomNumber
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
	CAR_NO = "carNo"
	OPERATOR_CODE =	"OperatorCode"
	NEXT_SITE_ID = "NextSiteId"
)

// 生成存储model
func SaveModelGen(object map[string]string, queueName string, hconnName string) func(){

	hconn := HconnContainersMapper[hconnName]
	currentHbaseConnection_Status := hconn.Name

	convertInfo,err := json.Marshal(object)
	util.CheckErr(err)

	// 人，站点，车，货
	// 站点有进出
	// operatorCode 以 ","为分割
	linkList := []string{
		EWB_NO,
		SITE_ID,
		VEHICLE_NO,
		CAR_NO,
		OPERATOR_CODE,

		NEXT_SITE_ID,
	}

	// linkTableNameList
	linkTnList := map[string]string{
		EWB_NO:"Link_Ewb",
		SITE_ID:"Link_Site",
		VEHICLE_NO:"Link_Vehicle",
		CAR_NO:"Link_Vehicle",

		OPERATOR_CODE:"Link_Operator",

		NEXT_SITE_ID : "Link_Site",
	}

	// history linktableName
	historyLinkTnList := map[string]string{
		EWB_NO:"History_Link_Ewb",
		SITE_ID:"History_Link_Site",
		VEHICLE_NO:"History_Link_Vehicle",
		CAR_NO:"History_Link_Vehicle",
		OPERATOR_CODE:"History_Link_Operator",
		NEXT_SITE_ID : "History_Link_Site",
	}

	// 基础信息存储 model map,似乎是多余的
	//basicInfoMapper := ModelMapperGen(object)
	basicInfoMapper := V2Byte(object)
	for key,value := range basicInfoMapper{
		hlogger.Debug("change [%s]=[%s]",key,string(value))
	}
	// 基础信息--columnFamily
	basicInfoCf := "basic"
	// 基础信息--tableName
	basicInfoTn := queueName
	// 基础信息--cf 和 model
	basicInfoCfMapper := map[string]map[string][]byte{basicInfoCf:basicInfoMapper}
	// 基础信息--uid----rowkey
	basicInfoUid,randomNumber := UidGen()

	//CacheMapper[string(basicInfoMapper["hewbno"])] = randomNumber

	return func() {
		hlogger.Debug("[%s] metaTable rowkey : [%s]",queueName,basicInfoUid)
		// 基础信息--putrequest
		biPutRequest, err := hrpc.NewPutStr(context.Background(),basicInfoTn,basicInfoUid,basicInfoCfMapper)
		if err!=nil {
			hlogger.Error("[%s] bi build [%s] hrpc error : %s",queueName,currentHbaseConnection_Status,err)
			return
		}
		_, err = singleConn.Client.Put(biPutRequest)
		if err!=nil {
			hlogger.Error("[%s] bi hbase [%s] put error : %s",queueName,currentHbaseConnection_Status,err)
			errorBackup(queueName,"basic",string(convertInfo))
			return
		}

		var saveLinkMapper map[string]string
		if isHistoryData(queueName) {
			saveLinkMapper = historyLinkTnList
		}else {
			saveLinkMapper = linkTnList
		}

		// 关联信息
		for _,v := range linkList{
			hlogger.Debug("[%s] find the mapper value : [%s]",strings.ToLower(v),basicInfoMapper[strings.ToLower(v)])

			if linkrowKey,ok := basicInfoMapper[strings.ToLower(v)]; ok && filterLinkRowKey(linkrowKey) {
				linkrowKeys := []string{string(linkrowKey)}
				// 操作者有多个
				if v==OPERATOR_CODE {
					linkrowKeys = strings.Split(string(linkrowKey),",")
				}

				for _,lrk := range linkrowKeys{
					// 去除例如 21850.0000 后缀的rowkey,变成21850
					lrk := FixRowKey(lrk)
					hlogger.Debug("[%s] LinkTable rowkey : [%s]",saveLinkMapper[v],lrk)
					saveKey := fmt.Sprintf("%s_%d",genCfColumnKeyName(basicInfoMapper,queueName),randomNumber)
					hlogger.Debug("[%s] saveKey",saveKey)
					// linkInfoMapper
					liMapper := map[string][]byte{saveKey:[]byte(basicInfoUid)}
					linkCfName := queueName
					if v==NEXT_SITE_ID {
						linkCfName = fmt.Sprintf("%s_%s","nextSite",linkCfName)
					}
					liCfMapper := map[string]map[string][]byte{linkCfName:liMapper}
					linkPutReq, err := hrpc.NewPutStr(context.Background(),saveLinkMapper[v],lrk,liCfMapper)
					if err!=nil {
						hlogger.Error("[%s] link [%s] build [%s] hrpc error : %s",queueName,v,currentHbaseConnection_Status,err)
						continue
					}
					_, err = singleConn.Client.Put(linkPutReq)
					if err!=nil {
						hlogger.Error("[%s] link [%s] hbase [%s] put error : %s",queueName,v,currentHbaseConnection_Status,err)
						errorBackup(queueName,v,string(convertInfo))
						continue
					}
				}
			}
		}
	}
}

func isHistoryData(queueName string) bool {
	return strings.Contains(strings.ToLower(queueName), "history");
}

func errorBackup(queueName string,savePeriod string,info string){
	errorLogger.Println(
		fmt.Sprintf(
			"%s\t%s\t%s",
			queueName,
			savePeriod,
			string(info),
		),
	)
}

func filterLinkRowKey(bytes []byte) bool{
	return (fmt.Sprintf("%s",bytes)!="/") && len(string(bytes))>=1
}

// 转成 value 为字节数组的形式
func V2Byte(i map[string]string) map[string][]byte {
	transferMapper := make(map[string][]byte,len(i))
	for k,v := range i{
		transferMapper[k] = []byte(v)
	}
	return transferMapper
}


func FixRowKey(str string) string {
	index := strings.Index(str,".")
	if str[index+1:]=="0000" {
		return str[:index]
	}
	return str
}

func genCfColumnKeyName(mapper map[string][]byte, queueName string) string {
	correspondingColumnName,ok := core.LinkKey[queueName]
	hlogger.Debug("[%s] columnName",correspondingColumnName)
	// 以多版本来记录

	if strings.ToLower(correspondingColumnName) == "scantime" && string(mapper[strings.ToLower(correspondingColumnName)])=="" {
		util.BarkNotify(fmt.Sprintf("[%s] lack the [%s]'s value",queueName,correspondingColumnName))
	}
	if !ok || correspondingColumnName=="uid" || string(mapper[strings.ToLower(correspondingColumnName)])==""{
		// 设置接收的当前时间
		return time.Now().Format("2006-01-02 15:04:05")
		// 备份数据固定时间
		//return ("2018-06-04 23:00:00")
	}
	return string(mapper[strings.ToLower(correspondingColumnName)])
}

func ModelGen(jsonByte []byte) []map[string]string{
	var tempInterface interface{}
	err := json.Unmarshal(jsonByte,&tempInterface)
	if err!=nil {
		return nil
	}
	var retArr []map[string]string
	for _,v := range tempInterface.([]interface{}){

		jsonMapper := make(map[string]string,16)
		jsonOne := v.(map[string]interface{})
		for key,value := range jsonOne{
			var mapperValue string
			switch value.(type) {
			case []interface{}:
				var sonList []map[string]string
				for _,sonv := range value.([]interface{}){
					sonMapper := make(map[string]string,4)
					for key,value := range sonv.(map[string]interface{}){
						sonMapper[key] = fmt.Sprintf("%v",value)
					}
					sonList = append(sonList,sonMapper)
				}
				s,_ := json.Marshal(sonList)
				mapperValue = string(s)
			case bool:
				mapperValue = fmt.Sprintf("%v",value)
			case float64:
				mapperValue = fmt.Sprintf("%.4f",value)
			case string:
				mapperValue = value.(string)
			case types.Nil:
			}
			jsonMapper[strings.ToLower(key)] = mapperValue
		}
		retArr = append(retArr,jsonMapper)
	}
	return retArr
}

